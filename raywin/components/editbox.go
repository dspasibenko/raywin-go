package components

import (
	"github.com/dspasibenko/raywin-go/raywin"
	rl "github.com/gen2brain/raylib-go/raylib"
	"sync/atomic"
)

type EditBox struct {
	raywin.BaseComponent

	text atomic.Value
}

func NewEditBox(owner raywin.Container) (*EditBox, error) {
	eb := &EditBox{}
	eb.SetText("ghpy")
	eb.SetBounds(rl.RectangleInt32{X: 0, Y: 0, Width: 100, Height: 100})
	return eb, eb.Init(owner, eb)
}

func (eb *EditBox) SetBounds(r rl.RectangleInt32) {
	r.Height = int32(S.PPcm * S.EditBoxHeightMm / 10.0)
	r.Width = max(r.Width, r.Height)
	eb.BaseComponent.SetBounds(r)
}

func (eb *EditBox) SetText(s string) {
	eb.text.Store(s)
}

func (eb *EditBox) Text() string {
	return eb.text.Load().(string)
}

func (eb *EditBox) Draw(cc *raywin.CanvasContext) {
	txt := eb.text.Load().(string)
	bi := eb.Bounds()
	b := bi.ToFloat32()
	x, y := cc.PhysicalPointXY(0, 0)
	b.X, b.Y = float32(x), float32(y)
	spacer := S.PPcm * S.EditBoxSpacerMm / 10.0
	v := rl.MeasureTextEx(raywin.SystemFont(int(S.EditBoxFontSize)), txt, S.EditBoxFontSize, 0)
	e := rl.Rectangle{X: float32(x) + b.Height/4.0, Y: float32(y) + spacer, Width: b.Width - b.Height/2, Height: b.Height - 2*spacer}
	curPos := v.X
	curPos = max(0.0, min(curPos, e.Width))
	curPos += e.X
	rl.DrawRectangleRounded(b, 0.5, 10, S.EditBoxOutlineColor)
	b.X += 2
	b.Y += 2
	b.Width -= 4
	b.Height -= 4
	rl.DrawRectangleRounded(b, 0.5, 10, S.EditBoxBackgoundColor)
	if v.X > 0.0 {
		rl.BeginScissorMode(int32(e.X), int32(e.Y), int32(e.Width), int32(e.Height))
		rl.DrawTextEx(raywin.SystemFont(int(S.EditBoxFontSize)), txt, rl.Vector2{X: curPos - v.X, Y: e.Y}, S.EditBoxFontSize, 0.0, S.EditBoxTextColor)
		rl.BeginScissorMode(int32(b.X), int32(b.Y), int32(b.Width), int32(b.Height))
	}
	m := raywin.Millis() % 1000
	col := S.EditBoxTextColor
	if m > 655 {
		col.A = 0
	} else if m > 400 {
		col.A = uint8(655 - int(m))
	}
	if col.A > 0 {
		rl.DrawRectangleV(rl.Vector2{X: curPos, Y: e.Y}, rl.Vector2{X: S.CurorWidth, Y: e.Height}, col)
	}
}
