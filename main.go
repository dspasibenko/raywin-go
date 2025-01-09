package main

import (
	"github.com/dspasibenko/raywin-go/raywin"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type myWidget struct {
	raywin.BaseContainer
	col  rl.Color
	ppos raywin.TPState
}

func (mw *myWidget) Draw(cc *raywin.CanvasContext) {
	rl.DrawRectangle(0, 0, 1000, 1000, mw.col)
}

func (mw *myWidget) OnTPState(tps raywin.TPState) raywin.OnTPSResult {
	if tps.State == raywin.TPStatePressed || tps.State == raywin.TPStateMoving {
		if tps.Sequence == mw.ppos.Sequence {
			r := mw.Bounds()
			r.X += int32(tps.Pos.X - mw.ppos.Pos.X)
			r.Y += int32(tps.Pos.Y - mw.ppos.Pos.Y)
			mw.SetBounds(r)
		}
		mw.ppos = tps
		return raywin.OnTPSResultLocked
	}
	return raywin.OnTPSResultNA
}

func main() {
	cfg := raywin.DefaultConfig()
	raywin.Init(cfg)
	mw := &myWidget{col: rl.White}
	mw.Init(raywin.RootContainer(), mw)
	mw.SetBounds(rl.RectangleInt32{10, 10, 300, 300})
	mw2 := &myWidget{col: rl.Blue}
	mw2.Init(mw, mw2)
	mw2.SetBounds(rl.RectangleInt32{10, 10, 100, 100})
	raywin.Run()
}
