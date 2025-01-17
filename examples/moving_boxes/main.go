package main

import (
	"github.com/dspasibenko/raywin-go/pkg/golibs/context"
	"github.com/dspasibenko/raywin-go/raywin"
	rl "github.com/gen2brain/raylib-go/raylib"
	"os"
	"syscall"
)

// myMovableWidget a component filled by the specified color. The component maybe moved
// by touchpad (or mouse). The Component is the Container, so it may contain other components
type myMovableWidget struct {
	raywin.BaseContainer
	col  rl.Color
	ppos raywin.TPState
}

// Draw just fills the whole drawing area by the component color
func (mw *myMovableWidget) Draw(cc *raywin.CanvasContext) {
	rl.DrawRectangle(0, 0, 1000, 1000, mw.col)
}

// OnTPState makes myMovableWidget implements Touchpadable interface. This is the component
// reaction on the touchpad (mouse) events
func (mw *myMovableWidget) OnTPState(tps raywin.TPState) raywin.OnTPSResult {
	if tps.State == raywin.TPStatePressed || tps.State == raywin.TPStateMoving {
		// the touchpad is pressed or the point moving...
		if tps.Sequence == mw.ppos.Sequence {
			// the previous notification was the same State
			r := mw.Bounds()
			r.X += int32(tps.Pos.X - mw.ppos.Pos.X)
			r.Y += int32(tps.Pos.Y - mw.ppos.Pos.Y)
			mw.SetBounds(r)
		}
		mw.ppos = tps
		// keeps the container in the focus
		return raywin.OnTPSResultLocked
	}
	return raywin.OnTPSResultNA
}

func main() {
	cfg := raywin.DefaultConfig()
	raywin.Init(cfg)

	// the white box owned by the display
	mw := &myMovableWidget{col: rl.White}
	mw.Init(raywin.RootContainer(), mw)
	mw.SetBounds(rl.RectangleInt32{X: 10, Y: 10, Width: 300, Height: 300})

	// the blue box in the white one
	mw2 := &myMovableWidget{col: rl.Blue}
	mw2.Init(mw, mw2)
	mw2.SetBounds(rl.RectangleInt32{X: 10, Y: 10, Width: 100, Height: 100})

	// the red box in the white one as well
	mw3 := &myMovableWidget{col: rl.Red}
	mw3.Init(mw, mw3)
	mw3.SetBounds(rl.RectangleInt32{X: 50, Y: 50, Width: 100, Height: 100})

	ctx := context.NewSignalsContext(os.Interrupt, syscall.SIGTERM) // allow to close the window by Ctrl+C in terminal
	raywin.Run(ctx)
}
