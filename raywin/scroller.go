package raywin

import "C"
import (
	"fmt"
	"github.com/dspasibenko/raywin-go/pkg/golibs/container"
	"github.com/dspasibenko/raywin-go/pkg/golibs/errors"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
	"sync/atomic"
)

type (

	// Scrollable interface applied to the components that support scrolling of their
	// drawing area.
	Scrollable interface {
		// Offset defines the position of the component's top-left corner relative to the
		// drawing grid's origin (0,0). By default, the top-left corner is at position (0,0).
		// Negative offsets indicate that the grid's origin has shifted to the right and downward.
		// For example, an offset of (-5, -10) moves the grid's origin 5 pixels to the right and
		// 10 pixels down from the top-left corner of the component's drawing area.
		// Conversely, positive offset values move the grid's origin upward and to the left
		// relative to the component's top-left corner.
		//
		// Components that support scrolling of their drawing area expose the interface to let
		// the raywin properly define the drawing context.
		Offset() Vector2Int32
	}

	// InertialScroller is a helper structure that facilitates scrolling with an inertia effect over a virtual area.
	//
	// The struct implements three interfaces: FrameListener, Scrollable, and Touchpadable.
	// A Component can either embed this structure to add scrolling functionality or decorate the interfaces
	// to include additional processing within the Component.
	//
	// Refer to the inertial_scroller example for usage instructions.
	InertialScroller struct {
		// flags contains settings whether the InertialScroller will work Horizontally, Vertically, or Both
		flags uint8

		// 'decel' represents the deceleration applied when the touchpad is released, controlling the inertia of movement.
		// Deceleration values are negative, where a lower (more negative) value results in faster stopping of movement.
		// Positive values are not allowed.
		decel   rl.Vector2
		owner   Component
		locked  bool
		prevPos rl.Vector2

		sinceMillis int64
		diff        rl.Vector2
		samples     container.RingBuffer[rl.Vector2]
		// post release scrolling
		velo rl.Vector2
		dir  rl.Vector2

		virtBounds atomic.Value // virtual Bounds
	}
)

var _ FrameListener = (*InertialScroller)(nil)
var _ Scrollable = (*InertialScroller)(nil)
var _ Touchpadable = (*InertialScroller)(nil)

const (
	ScrollBoth       = ScrollHorizontal | ScrollVertical
	ScrollHorizontal = 1
	ScrollVertical   = 2
)

// DefaultInternalScrollerDeceleration returns the default InternalScroller deceleration. The parameters
// are adjusted for 60FPS and the 800x600 screen size. To decelerate faster, put lower (bigger absolute)
// values
func DefaultInternalScrollerDeceleration() rl.Vector2 {
	assertInitialized()
	fps := uint(c.disp.cfg.FPS)
	return rl.Vector2{X: -float32(8) / float32(fps), Y: -float32(8) / float32(fps)}
}

func (s *InertialScroller) InitScroller(owner Component, virtBounds rl.RectangleInt32, decel rl.Vector2, flags uint8) error {
	if decel.Y >= 0 || decel.X >= 0 {
		return fmt.Errorf("InitScroller: decel.X=%f, decel.Y=%f cannot be positive: %w", decel.X, decel.Y, errors.ErrInvalid)
	}
	if owner == nil {
		return fmt.Errorf("InitScroller: owner is nil: %w", errors.ErrInvalid)
	}
	fps := uint(c.disp.cfg.FPS)
	s.samples = container.NewRingBuffer[rl.Vector2](fps / 3)
	s.flags = flags
	s.decel = decel
	s.owner = owner
	s.virtBounds.Store(virtBounds)
	return nil
}

// OnTPState implements the Touchpadable interface. The function must be called
// by raywin only
func (s *InertialScroller) OnTPState(tps TPState) OnTPSResult {
	s.diff = rl.Vector2{}
	if tps.State == TPStateReleased && s.locked {
		mn := s.samples.At(0)
		mx := mn
		ln := s.samples.Len()
		mnx, mny, mxx, mxy := ln, ln, ln, ln
		for s.samples.Len() > 0 {
			v, _ := s.samples.Read()
			if mn.X > v.X {
				mn.X = v.X
				mnx = s.samples.Len()
			}
			if mn.Y > v.Y {
				mn.Y = v.Y
				mny = s.samples.Len()
			}
			if mx.X < v.X {
				mx.X = v.X
				mxx = s.samples.Len()
			}
			if mx.Y < v.Y {
				mx.Y = v.Y
				mxy = s.samples.Len()
			}
		}
		dir := rl.Vector2{X: 1.0, Y: 1.0}
		if mnx > mxx {
			dir.X = -1.0
			mxx, mnx = mnx, mxx
		}
		if mny > mxy {
			dir.Y = -1.0
			mxy, mny = mny, mxy
		}
		velo := rl.Vector2{}
		if mx.X-mn.X >= 0.1 {
			velo.X = 0.5 * (mx.X - mn.X) / float32(mxx-mnx)
		}
		if mx.Y-mn.Y >= 0.1 {
			velo.Y = 0.5 * (mx.Y - mn.Y) / float32(mxy-mny)
		}
		s.sinceMillis = tps.Millis
		s.velo = velo
		s.dir = dir
	}
	if tps.State == TPStateMoving {
		if s.samples.Len() == s.samples.Cap() {
			s.samples.Skip(1)
		}
		s.samples.Write(tps.Pos)
		if s.locked {
			s.diff = VectorDiff(s.prevPos, tps.Pos)
		} else {
			s.diff = rl.Vector2{}
		}
	}
	// if the component has not locked the touchpad and we have scrolling in one
	// direction only, will lock the touchpad only if the movement was made in the dirrection then
	if !s.locked && tps.State == TPStateMoving && s.flags&ScrollBoth != ScrollBoth {
		s.locked = (s.flags&ScrollHorizontal != 0 && math.Abs(float64(s.prevPos.X-tps.Pos.X)) > 3*math.Abs(float64(s.prevPos.Y-tps.Pos.Y))) ||
			(s.flags&ScrollVertical != 0 && math.Abs(float64(s.prevPos.Y-tps.Pos.Y)) > 3*math.Abs(float64(s.prevPos.X-tps.Pos.X)))
	} else {
		s.locked = tps.State == TPStateMoving
	}
	s.prevPos = tps.Pos
	if s.locked {
		return OnTPSResultLocked
	}
	return OnTPSResultNA
}

// SetVirtualBounds allows to specify the virtual size for the scroller
func (s *InertialScroller) SetVirtualBounds(bounds rl.RectangleInt32) {
	s.virtBounds.Store(bounds)
}

// VirtualBounds returns the VirtualBounds (the offset point and the area size) as rl.RectangleInt32
func (s *InertialScroller) VirtualBounds() rl.RectangleInt32 {
	return s.virtBounds.Load().(rl.RectangleInt32)
}

// IsLocked returns whether the InertialScroller holds touchpad control or not
func (s *InertialScroller) IsTPLocked() bool {
	return s.locked
}

// Offset provides an implementation of Scrollable interface.
func (s *InertialScroller) Offset() Vector2Int32 {
	r := s.virtBounds.Load().(rl.RectangleInt32)
	return Vector2Int32{X: r.X, Y: r.Y}
}

// OnNewFrame provides the FrameListener interface. The function must be called
// by raywin only
func (s *InertialScroller) OnNewFrame(millis int64) {
	if !s.locked && !IsEmpty(s.dir) {
		s.diff.X = max(0.0, s.velo.X+s.decel.X*float32(millis-s.sinceMillis)/15.0-s.decel.X/2)
		s.diff.Y = max(0.0, s.velo.Y+s.decel.Y*float32(millis-s.sinceMillis)/15.0-s.decel.Y/2)
		if IsEmpty(s.diff) {
			s.dir = rl.Vector2{}
		}
		s.diff.X *= s.dir.X
		s.diff.Y *= s.dir.Y
	}

	diff := s.getDiffForLastFrame()
	p := s.virtBounds.Load().(rl.RectangleInt32)
	p.X = int32(float32(p.X) + diff.X)
	p.Y = int32(float32(p.Y) + diff.Y)
	if !s.locked {
		if p.X < 0 {
			p.X = min(0, p.X-p.X/3+1)
		}
		if p.Y < 0 {
			p.Y = min(0, p.Y-p.Y/3+1)
		}
		r := s.owner.Bounds()
		if p.X > 0 && r.Width > (p.Width-p.X) {
			p.X = max(0, p.X-(r.Width-p.Width+p.X)/3-1)
		}
		if p.Y > 0 && r.Height > (p.Height-p.Y) {
			p.Y = max(0, p.Y-(r.Height-p.Height+p.Y)/3-1)
		}
	}
	s.virtBounds.Store(p)
}

func (s *InertialScroller) getDiffForLastFrame() rl.Vector2 {
	if s.flags&ScrollHorizontal == 0 {
		s.diff.X = 0
	}
	if s.flags&ScrollVertical == 0 {
		s.diff.Y = 0
	}
	return s.diff
}
