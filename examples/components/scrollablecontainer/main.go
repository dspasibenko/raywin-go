package main

import (
	"github.com/dspasibenko/raywin-go/pkg/golibs/context"
	"github.com/dspasibenko/raywin-go/raywin"
	"github.com/dspasibenko/raywin-go/raywin/components"
	rl "github.com/gen2brain/raylib-go/raylib"
	"os"
	"syscall"
)

// myScrollable is the Container with scrolling bars
type myScrollable struct {
	components.ScrollableContainer
}

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

// Draw allows to draw a white background of myScrollable
func (ms *myScrollable) Draw(cc *raywin.CanvasContext) {
	b := ms.Bounds()
	off := ms.Offset()
	x, y := cc.PhysicalPointXY(off.X, off.Y)
	rl.DrawRectangle(x, y, b.Width, b.Height, rl.White)
	rl.DrawText("Press mouse button and move it", x, y+200, 30, rl.Black)
}

func main() {
	cfg := raywin.DefaultConfig()
	// to use components with their style, register its outlet in the config
	cfg.FrameListener = components.DefaultStyleOutlet(cfg.DisplayConfig)
	raywin.Init(cfg)

	mw := &myScrollable{}
	mw.InitScrollableContainer(raywin.RootContainer(), mw, components.ShowBothScrollBar|raywin.ScrollBoth)
	mw.SetBounds(rl.RectangleInt32{X: 50, Y: 50, Width: 500, Height: 500})
	mw.SetVirtualBounds(rl.RectangleInt32{X: 0, Y: 0, Width: 1280, Height: 720})

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
