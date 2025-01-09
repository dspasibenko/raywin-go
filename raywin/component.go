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
	"sync"
	"sync/atomic"
)

type (
	Component interface {
		Draw(cc *CanvasContext)
		IsVisible() bool
		SetVisible(b bool)
		Close()
		SetBounds(r rl.RectangleInt32)
		Bounds() rl.RectangleInt32

		baseComponent() *BaseComponent
	}

	Container interface {
		// Children is a list of Components owned by the container. The latest one is on drawing last
		// so it is on top of others
		Children() []Component
		AddChild(c Component) error
		RemoveChild(c Component) bool
	}

	FrameListener interface {
		// OnNewFrame is called for all components that implement the interface.
		// It is called after OnTPState call
		OnNewFrame(millis int64)
	}

	// BaseComponent contains some basic implementation of Component
	// After construction SetPosition() must be called as a first thing!
	BaseComponent struct {
		visible int32
		bounds  atomic.Value
		Lock    sync.Mutex
		// owner contains a reference to the owner of the component
		owner Container
		// this is the reference to the BaseComponent holder. This is because
		// a component may "extend" BaseComponent, we need to store the reference
		// to the holder. See Init()
		this   Component
		closed atomic.Bool
	}

	BaseContainer struct {
		BaseComponent
		children atomic.Value
	}
)

func (bc *BaseContainer) InitUnsafe(owner Container, this Component) error {
	if err := bc.BaseComponent.InitUnsafe(owner, this); err != nil {
		return err
	}
	bc.children.Store([]Component(nil))
	return nil
}

func (bc *BaseContainer) String() string {
	return "fixme"
}

func (bc *BaseContainer) AddChild(c Component) error {
	if !bc.lockIfAlive() {
		return fmt.Errorf("AddChild: failed to add %s to the %s container, which is not initialized: %w", c, bc, errors.ErrClosed)
	}
	defer bc.Lock.Unlock()

	cb := c.baseComponent()
	if err := cb.isInitialized(); err != nil {
		return err
	}
	cb.owner = bc

	v := bc.children.Load().([]Component)
	idx := childIndex(v, c)
	var nv []Component
	if idx < len(v) {
		// c is in the list, change its position then, briniging on top
		nv = make([]Component, 0, len(v))
		nv = append(nv, v[:idx]...)
		nv = append(nv, v[idx+1:]...)
	} else {
		nv = make([]Component, 0, len(v)+1)
		nv = append(nv, v...)
	}
	nv = append(nv, c)
	bc.children.Store(nv)
	return nil
}

func (bc *BaseContainer) RemoveChild(c Component) bool {
	if !bc.lockIfAlive() {
		return false
	}
	defer bc.Lock.Unlock()

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

func (bc *BaseContainer) Children() []Component {
	return bc.children.Load().([]Component)
}

func (bc *BaseContainer) Close() {
	if !bc.lockIfAlive() {
		return
	}
	children := bc.children.Load().([]Component)
	bc.children.Store([]Component(nil))
	bc.close()
	bc.Lock.Unlock()
	for _, c := range children {
		c.Close()
	}
}

func childIndex(children []Component, c Component) int {
	for idx, c1 := range children {
		if c1 == c {
			return idx
		}
	}
	return len(children)
}

func (bc *BaseComponent) String() string {
	return "fixme"
}

func (bc *BaseComponent) lockIfAlive() bool {
	if bc.closed.Load() {
		return false
	}
	bc.Lock.Lock()
	if bc.closed.Load() {
		bc.Lock.Unlock()
		return false
	}
	return true
}

func (bc *BaseComponent) isInitialized() error {
	if bc.this == nil || bc.this.baseComponent() != bc {
		return fmt.Errorf("Init() is not called or the compnent %s is closed", bc)
	}
	return nil
}

func (bc *BaseComponent) baseComponent() *BaseComponent {
	return bc
}

type unsafeInitializer interface {
	InitUnsafe(owner Container, this Component) error
}

func (bc *BaseComponent) Init(owner Container, this Component) error {
	bc.Lock.Lock()
	defer bc.Lock.Unlock()

	if this.baseComponent() != bc {
		return fmt.Errorf("this %s must embed %s: %w", this, bc, errors.ErrInvalid)
	}

	return this.(unsafeInitializer).InitUnsafe(owner, this)
}

// InitLocked must be called for any instance as first thing after its creation.
// owner defines the component which owns this. this is the Component, which
// holds the BaseComponent
func (bc *BaseComponent) InitUnsafe(owner Container, this Component) error {
	if this.baseComponent() != bc {
		return fmt.Errorf("this %s must embed %s: %w", this, bc, errors.ErrInvalid)
	}
	if bc.owner != nil {
		return fmt.Errorf("this %s already has owner %s: %w", this, bc.owner, errors.ErrInvalid)
	}
	if owner.(Component).baseComponent() == bc {
		return fmt.Errorf("this %s cannot be added to itself %s: %w", this, owner, errors.ErrInvalid)
	}
	bc.SetVisible(true)
	if bc.bounds.Load() == nil {
		bc.bounds.Store(rl.RectangleInt32{})
	}
	bc.this = this
	err := owner.AddChild(this)
	if err != nil {
		bc.this = nil
	}
	return err
}

func (bc *BaseComponent) SetBounds(r rl.RectangleInt32) {
	bc.bounds.Store(r)
}

func (bc *BaseComponent) Bounds() rl.RectangleInt32 {
	res := bc.bounds.Load().(rl.RectangleInt32)
	return res
}

// IsVisible returns whether the component is visible or not
func (bc *BaseComponent) IsVisible() bool {
	return atomic.LoadInt32(&bc.visible) != 0
}

// SetVisible allows to specify the component visibility
func (bc *BaseComponent) SetVisible(visible bool) {
	if visible {
		atomic.StoreInt32(&bc.visible, 1)
	} else {
		atomic.StoreInt32(&bc.visible, 0)
	}
}

func (bc *BaseComponent) Close() {
	bc.Lock.Lock()
	defer bc.Lock.Unlock()

	bc.close()
}

func (bc *BaseComponent) close() {
	if bc.closed.Load() {
		return
	}
	bc.closed.Store(true)
	bc.owner.RemoveChild(bc.this)
	bc.owner = nil
	bc.this = nil
}

func (bc *BaseComponent) Draw(cc *CanvasContext) {
}
