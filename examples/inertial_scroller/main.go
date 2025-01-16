package main

import (
	"fmt"
	"github.com/dspasibenko/raywin-go/pkg/golibs/context"
	"github.com/dspasibenko/raywin-go/raywin"
	rl "github.com/gen2brain/raylib-go/raylib"
	"os"
	"syscall"
)

type topPanel struct {
	raywin.BaseComponent
	s raywin.Scrollable
}

type myScrollable struct {
	raywin.BaseContainer
	raywin.InertialScroller
}

// myMovableWidget a component filled by the specified color. The component maybe moved
// by touchpad (or mouse). The Component is the Container, so it may contain other components
type myMovableWidget struct {
	raywin.BaseComponent
	col rl.Color
	txt string
}

func (tp *topPanel) Init(owner raywin.Container, s raywin.Scrollable) error {
	if err := tp.BaseComponent.Init(owner, tp); err != nil {
		return err
	}
	r := owner.(raywin.Component).Bounds()
	r.Height = 50
	tp.SetBounds(r)
	tp.s = s
	return nil
}

func (tp *topPanel) Draw(cc *raywin.CanvasContext) {
	r := tp.Bounds()
	rl.DrawRectangle(0, 0, r.Width, r.Height, rl.White)
	p := tp.s.Offset()
	rl.DrawText(fmt.Sprintf("Use mouse to move the grid below. Offset X=%d, Y=%d", p.X, p.Y), 15, 15, 20, rl.Color{255, 0, 0, 255})
}

func (ms *myScrollable) Init(owner raywin.Container) error {
	err := ms.BaseContainer.Init(owner, ms)
	if err != nil {
		return err
	}
	r := owner.(raywin.Component).Bounds()
	r.Y = 50
	r.Height -= 50
	ms.SetBounds(r)
	return ms.InitScroller(ms, rl.RectangleInt32{0, 0, 2000, 2000}, raywin.DefaultInternalScrollerDeceleration(), raywin.ScrollBoth)
}

// Draw just fills the whole drawing area by the component color
func (mw *myMovableWidget) Draw(cc *raywin.CanvasContext) {
	r := mw.Bounds()
	x, y := cc.PhysicalPointXY(0, 0)
	rl.DrawRectangle(x, y, r.Width, r.Height, mw.col)
}

func main() {
	cfg := raywin.DefaultConfig()
	raywin.Init(cfg)

	scrlbl := &myScrollable{}
	if err := scrlbl.Init(raywin.RootContainer()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	panel := &topPanel{}
	if err := panel.Init(raywin.RootContainer(), scrlbl); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for i := int32(0); i < 10; i++ {
		for j := int32(0); j < 10; j++ {
			mw := &myMovableWidget{col: rl.Color{R: 150 + uint8(i*10+j), G: 100 + uint8(i)*10, B: 100 + uint8(j)*10, A: 255}}
			if err := mw.Init(scrlbl, mw); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			mw.txt = fmt.Sprintf("%d", i*10+j)
			mw.SetBounds(rl.RectangleInt32{10 + 200*j, 10 + 200*i, 190, 190})
		}
	}

	ctx := context.NewSignalsContext(os.Interrupt, syscall.SIGTERM) // allow to close the window by Ctrl+C in terminal
	raywin.Run(ctx)
}
