package components

import (
	"github.com/dspasibenko/raywin-go/raywin"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Toggle struct {
	raywin.BaseComponent
	raywin.Pressor

	on        bool
	pressedAt int64
}

func (t *Toggle) InitToggle(owner raywin.Container, onToggle func(newState bool) bool) error {
	t.InitPressor(S.TogglePressRadius, S.TogglePressMillis, func() {
		t.on = !t.on
		if onToggle != nil {
			t.on = onToggle(t.on)
		}
		t.pressedAt = raywin.Millis()
	})
	t.SetBounds(rl.RectangleInt32{})
	return t.Init(owner, t)
}

func (t *Toggle) SetBounds(b rl.RectangleInt32) {
	b.Width = int32(S.ToggleWidthMm * S.PPcm / 10)
	b.Height = int32(S.ToggleHeightMm * S.PPcm / 10)
	t.BaseComponent.SetBounds(b)
}

func (t *Toggle) Draw(cc *raywin.CanvasContext) {
	b := t.Bounds()
	s := S.ToggleSpaceMm * S.PPcm / 10.0
	rad := (float32(b.Height)) / 2.0
	x, y := cc.PhysicalPointXY(0, 0)
	r := rl.Rectangle{X: float32(x) + rad, Y: float32(y), Width: float32(b.Width) - 2*rad, Height: float32(b.Height)}
	if t.Pressed() {
		r1 := r
		st := s / 4
		for i := 0; i < 4; i++ {
			t.drawFrame(r1, rl.Fade(S.FrameSelectToneColor, 0.3))
			r1.X += st
			r1.Y += st
			r1.Width -= 2 * st
			r1.Height -= 2 * st
		}
		t.drawFrame(r1, S.FrameSelectColor)
	}
	r.X += s
	r.Y += s
	r.Width -= 2 * s
	r.Height -= 2 * s
	t.drawFrame(r, S.FrameColor)
	r.Y += 2
	r.Height -= 4
	ballOffset := max(0.0, 1.0+float32(t.pressedAt-raywin.Millis())/100)
	if t.on {
		x := r.X + r.Width - float32(r.Width)*ballOffset
		t.drawFrame(r, S.ToggleOnColor)
		v := rl.Vector2{X: x, Y: r.Y + r.Height/2}
		rl.DrawCircleV(v, float32(r.Height)/2-2.0, S.FrameSelectColor)
		rl.DrawCircleLinesV(v, float32(r.Height)/2-2.0, S.FrameShadeColor)
	} else {
		x := r.X + float32(r.Width)*ballOffset
		t.drawFrame(r, S.ToggleOffColor)
		v := rl.Vector2{X: x, Y: r.Y + r.Height/2}
		rl.DrawCircleV(v, float32(r.Height)/2-2.0, S.FrameShadeColor)
	}
}

func (t *Toggle) drawFrame(r rl.Rectangle, color rl.Color) {
	rad := r.Height / 2.0
	rl.DrawRectangleV(rl.Vector2{X: r.X, Y: r.Y}, rl.Vector2{X: r.Width, Y: r.Height}, color)
	rl.DrawCircleSector(rl.Vector2{X: r.X, Y: r.Y + rad}, rad, 90, 270, int32(rad), color)
	rl.DrawCircleSector(rl.Vector2{X: r.X + r.Width, Y: r.Y + rad}, rad, 270, 450, int32(rad), color)
}
