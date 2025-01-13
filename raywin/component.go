// Copyright 2025 Dmitry Spasibenko
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package raywin

import (
	"fmt"
	"github.com/dspasibenko/raywin-go/pkg/golibs/errors"
	rl "github.com/gen2brain/raylib-go/raylib"
	"reflect"
	"sync"
	"sync/atomic"
)

type (
	// Component defines an interface for a drawable box. All structs that implement the interface
	// must have BaseComponent embedded into the structs (please see below)
	Component interface {
		// Draw renders the component within the specified physical region, as defined by the cc parameter.
		// The implementation utilizes Raylib functions, such as rl.Rectangle(), to draw the component on
		// the display. The cc parameter specifies the position of the component on the physical display.
		//
		// By default, Draw() is invoked for the physical region where the component is defined. Raywin
		// uses scissors to constrain the drawing area. The implementation can adjust the drawing area
		// by calling rl.BeginScissorMode() if the region need to be changed.
		//
		// Raywin invokes Draw() for all visible components in each frame. A component is considered visible
		// if IsVisible() returns true and its Bounds() intersect with the visible region defined by its
		// parent Component (see Container).
		Draw(cc *CanvasContext)

		// IsVisible returns whether the component is visible or not
		IsVisible() bool

		// SetVisible sets the component visibility
		SetVisible(b bool)

		// Close closes the component and frees all resources
		Close()

		// SetBounds defines the Component position on the parent's component region and its size.
		SetBounds(r rl.RectangleInt32)

		// Bounds returns the position and size of the component. The position is defined relative to the region
		// of the parent component.
		Bounds() rl.RectangleInt32

		// baseComponent() is the private function to make the interface be implemented by the BaseComponent
		// defined in the package.
		baseComponent() *BaseComponent
	}

	// Container is a Component which may contain children components. The drawing area for the
	// Container's children is limited by the Container's region
	Container interface {
		// Children returns the list of components owned by the container. The components in the list
		// are drawn in the order they appear. The first component is drawn first, followed by the second,
		// third, and so on. As a result, the last component is drawn on top and will cover all previous
		// components if they overlap.
		//
		// The Draw function for the Container is called before the Draw functions of its children,
		// ensuring that the children are drawn on top of the container's drawings.
		Children() []Component

		// OnAddChild is called to update the list of children by adding the child `c` to it.
		// It must return the updated collection of children with `c` added, or an error if the operation is not possible.
		//
		// OnAddChild is a special public function that is called while holding the internal lock of `bc`.
		// Therefore, if this function is overridden, no calls to `bc` should be made within the override,
		// as it may lead to deadlock.
		//
		// Even the function is implemented in BaseContainer and the default implementation may be good enough,
		// users may override this function to customize the order in which children are stored.
		// The default implementation simply adds `c` to the end of the children slice, or adds to the end
		// it if it already exists.
		OnAddChild(c Component, children []Component) ([]Component, error)

		// baseContainer() makes all containers implemented on top of BaseContainer
		baseContainer() *BaseContainer
	}

	// FrameListener interface allows components to be notified about each new frame
	FrameListener interface {
		// OnNewFrame is called for each component that implements the interface on every new frame.
		// The `millis` timestamp is monotonically increasing and can be used to measure the time elapsed
		// between different frames renderings. It is also used in various notification calls to identify the frame's
		// timestamp, effectively serving as a unique identifier for the frame.
		//
		// the `millis` may be related to the clock and may be not, so it cannot be used to identify
		// the current time, but for measuring the time intervals between frames only.
		OnNewFrame(millis int64)
	}

	// BaseComponent provides the fundamental implementation of all components. It means that any component
	// should embed the struct. So as the Component interface has the package private baseComponent() function,
	// it is not possible to implement a Component without embedding the BaseComponent.
	//
	// BaseComponent supports any component life-cycle:
	// - Init() - the function initializes BaseComponent
	// - Close() - the Component termination
	// In addition, the implementation also provides implementation for the component visibility, to avoid
	// some boilerplate code in the derived components
	BaseComponent struct {
		visible atomic.Bool
		bounds  atomic.Value
		lock    sync.Mutex
		// owner contains a reference to the owner of the component
		owner *BaseContainer
		// this is the reference to the BaseComponent holder. This is because
		// a component may "extend" BaseComponent, we need to store the reference
		// to the holder. See Init()
		this   Component
		tpName atomic.Value
		closed atomic.Bool
	}

	// BaseContainer struct offers a basic implementation of Container interface. Complex
	// components, that suppose to own other components, may use the BaseContainer for the
	// basic Container implementation.
	BaseContainer struct {
		BaseComponent

		children atomic.Value
	}
)

// Init initializes BaseContainer, returns nil if the component initialized successfully
func (bc *BaseContainer) Init(owner Container, this Component) error {
	bc.lock.Lock()
	defer bc.lock.Unlock()
	if err := bc.init(owner, this); err != nil {
		return err
	}
	bc.children.Store([]Component(nil))
	return nil
}

// OnAddChild is the default implementation, please see Container interface
func (bc *BaseContainer) OnAddChild(c Component, children []Component) ([]Component, error) {
	idx := childIndex(children, c)
	var nv []Component
	if idx < len(children) {
		// c is in the list, change its position then, briniging on top
		nv = make([]Component, 0, len(children))
		nv = append(nv, children[:idx]...)
		nv = append(nv, children[idx+1:]...)
	} else {
		nv = make([]Component, 0, len(children)+1)
		nv = append(nv, children...)
	}
	return append(nv, c), nil
}

// addChild adds the new comopnent c to the container
func (bc *BaseContainer) addChild(c Component) error {
	if !bc.lockIfAlive() {
		return fmt.Errorf("AddChild: failed to add %s to the %s container, which is not initialized: %w", c, bc, errors.ErrClosed)
	}
	defer bc.lock.Unlock()

	cb := c.baseComponent()
	if err := cb.AssertInitialized(); err != nil {
		return err
	}
	if cb.owner != nil && cb.owner != bc.this.(Container).baseContainer() {
		return fmt.Errorf("the component %s, already has an owner: %w", cb, errors.ErrInvalid)
	}

	v := bc.children.Load().([]Component)
	nv, err := bc.this.(Container).OnAddChild(c, v)
	if err != nil {
		return err
	}
	bc.children.Store(nv)
	return nil
}

func (bc *BaseContainer) baseContainer() *BaseContainer {
	return bc
}

// removeChild removes the component c from the container
func (bc *BaseContainer) removeChild(c Component) bool {
	if !bc.lockIfAlive() {
		return false
	}
	defer bc.lock.Unlock()

	v := bc.children.Load().([]Component)
	idx := childIndex(v, c)
	if idx < len(v) {
		nv := make([]Component, 0, len(v)-1)
		nv = append(nv, v[:idx]...)
		nv = append(nv, v[idx+1:]...)
		bc.children.Store(nv)
		return true
	}
	return false
}

// Children returns list of owned components
func (bc *BaseContainer) Children() []Component {
	return bc.children.Load().([]Component)
}

// Close terminates the Container and all its children
func (bc *BaseContainer) Close() {
	if !bc.lockIfAlive() {
		return
	}
	children := bc.children.Load().([]Component)
	bc.children.Store([]Component(nil))
	bc.close()
	bc.lock.Unlock()
	for _, c := range children {
		c.Close()
	}
}

// String returns the BaseContainer string representation
func (bc *BaseContainer) String() string {
	return fmt.Sprintf("{BC: %s, children: %d}", bc.baseComponent(), len(bc.Children()))
}

func childIndex(children []Component, c Component) int {
	for idx, c1 := range children {
		if c1 == c {
			return idx
		}
	}
	return len(children)
}

// Init initializes BaseComponent. owner should be non-nil the owner of the Component,
// `this` contains the final struct, which implements Component, but which embed the bc
//
// Init must be called for any instance as first thing after its creation.
func (bc *BaseComponent) Init(owner Container, this Component) error {
	bc.lock.Lock()
	defer bc.lock.Unlock()

	return bc.init(owner, this)
}

func (bc *BaseComponent) init(owner Container, this Component) error {
	if this.baseComponent() != bc {
		return fmt.Errorf("this %s must embed %s: %w", this, bc, errors.ErrInvalid)
	}
	if bc.owner != nil {
		return fmt.Errorf("this %s already has owner %s: %w", this, bc.owner, errors.ErrInvalid)
	}
	o := owner.baseContainer()
	if owner.(Component).baseComponent() == bc {
		return fmt.Errorf("this %s cannot be added to itself %s: %w", this, owner, errors.ErrInvalid)
	}
	bc.SetVisible(true)
	if bc.bounds.Load() == nil {
		bc.bounds.Store(rl.RectangleInt32{})
	}
	bc.tpName.Store(reflect.TypeOf(this).String()) // to be sure that AssertInitialized is nil
	bc.this = this
	err := o.addChild(this)
	if err != nil {
		bc.this = nil
	} else {
		bc.owner = o
	}
	return err
}

// SetBounds allows to assing the comonent position and dimensions by the `r`
func (bc *BaseComponent) SetBounds(r rl.RectangleInt32) {
	bc.bounds.Store(r)
}

// Bounds returns the component position on its owner coordinates, and its size as rl.RectangleInt32
func (bc *BaseComponent) Bounds() rl.RectangleInt32 {
	v := bc.bounds.Load()
	if v == nil {
		return rl.RectangleInt32{}
	}
	return v.(rl.RectangleInt32)
}

// IsVisible returns whether the component is visible or not
func (bc *BaseComponent) IsVisible() bool {
	return bc.visible.Load()
}

// SetVisible allows to specify the component visibility
func (bc *BaseComponent) SetVisible(visible bool) {
	bc.visible.Store(visible)
}

// Draw is the BaseComponent drawing procedure which does nothing. It is here to support
// the Component interface, should be re-defined in the derived structure
func (bc *BaseComponent) Draw(cc *CanvasContext) {
}

// Close allows to close the BaseComponent
func (bc *BaseComponent) Close() {
	bc.lock.Lock()
	defer bc.lock.Unlock()

	bc.close()
}

// AssertInitialized returns an error if the component is not initialized
func (bc *BaseComponent) AssertInitialized() error {
	if bc.tpName.Load() == nil || bc.closed.Load() {
		return fmt.Errorf("Init() is not called or the compnent %s is closed: %w", bc.String(), errors.ErrInvalid)
	}
	return nil
}

func (bc *BaseComponent) isClosed() bool {
	return bc.closed.Load()
}

// String returns the `bc` description
func (bc *BaseComponent) String() string {
	v := bc.tpName.Load()
	tp := "N/A"
	if v != nil {
		tp = v.(string)
	}
	return fmt.Sprintf("{Type:%s, Bounds:%s, visible:%t, closed:%t}", tp, RectangleInt32ToString(bc.Bounds()),
		bc.IsVisible(), bc.closed.Load())
}

func (bc *BaseComponent) lockIfAlive() bool {
	if bc.closed.Load() {
		return false
	}
	bc.lock.Lock()
	if bc.closed.Load() {
		bc.lock.Unlock()
		return false
	}
	return true
}

func (bc *BaseComponent) baseComponent() *BaseComponent {
	return bc
}

func (bc *BaseComponent) close() {
	if bc.closed.Load() {
		return
	}
	bc.closed.Store(true)
	if bc.owner != nil {
		bc.owner.removeChild(bc.this)
	}
	bc.owner = nil
	bc.this = nil
}
