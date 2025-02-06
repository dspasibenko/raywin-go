package components

import (
	"github.com/dspasibenko/raywin-go/pkg/golibs"
	"github.com/dspasibenko/raywin-go/raywin"
	rl "github.com/gen2brain/raylib-go/raylib"
	"image/color"
	"sync/atomic"
)

type Button struct {
	raywin.BaseComponent
	raywin.Pressor

	text      string
	textSize  rl.Vector2
	bs        atomic.Value
	once      golibs.Once
	pressedAt int64
	fadeK     float32
}

type ButtonStyle struct {
	TextFont     rl.Font
	TextFontSize float32
	Color        color.RGBA
	OutlineColor color.RGBA
	TextColor    color.RGBA
	SelectColor  color.RGBA
	Icon         string
	Flags        int // See constants below (ButtonSelectStyleJumpOut etc.)
}

func DialogButtonStyle() ButtonStyle {
	return ButtonStyle{
		TextFont:     raywin.SystemItalicFont(),
		TextFontSize: 60,
		Color:        S.DialogBackgroundDark,
		OutlineColor: color.RGBA{14, 110, 138, 255},
		SelectColor:  S.DialogBackgroundLight,
		TextColor:    rl.White,
		Flags:        ButtonSelectStyleJumpOut | ButtonFrameOutlined,
	}
}

func DialogButtonCancelStyle() ButtonStyle {
	return ButtonStyle{
		TextFont:     raywin.SystemFont(),
		TextFontSize: 30,
		Color:        color.RGBA{82, 2, 2, 255},
		SelectColor:  color.RGBA{107, 2, 2, 255},
		TextColor:    rl.White,
		Flags:        ButtonSelectStyleSwell | ButtonFrameRounded,
	}
}

func DialogButtonOkStyle() ButtonStyle {
	return ButtonStyle{
		TextFont:     raywin.SystemFont(),
		TextFontSize: 30,
		Color:        color.RGBA{4, 51, 38, 255},
		SelectColor:  color.RGBA{6, 71, 53, 255},
		TextColor:    rl.White,
		Flags:        ButtonSelectStyleSwell | ButtonFrameRounded,
	}
}

func DialogButtonControlStyle() ButtonStyle {
	return ButtonStyle{
		TextFont:     raywin.SystemItalicFont(),
		TextFontSize: 25,
		Color:        S.DialogBackgroundDark,
		OutlineColor: color.RGBA{14, 110, 138, 255},
		SelectColor:  S.DialogBackgroundLight,
		TextColor:    rl.White,
		Flags:        ButtonSelectStyleJumpOut | ButtonFrameOutlined,
	}
}

func DialogButtonCloseStyle() ButtonStyle {
	return ButtonStyle{
		TextFont:     raywin.SystemFont(),
		TextFontSize: 0,
		Color:        S.DialogBackgroundDark,
		OutlineColor: color.RGBA{14, 110, 138, 255},
		SelectColor:  S.DialogBackgroundLight,
		TextColor:    rl.White,
		Icon:         "x-white",
		Flags:        ButtonSelectStyleHighlighted | ButtonFrameOutlined,
	}
}

func DialogButtonTransparrentStyle() ButtonStyle {
	return ButtonStyle{
		TextFont:     raywin.SystemFont(),
		TextFontSize: 0,
		Color:        S.TransparentColor,
		OutlineColor: S.TransparentColor,
		SelectColor:  rl.Fade(rl.Orange, 0.4),
		TextColor:    rl.White,
		Icon:         "airplane-yellow",
		Flags:        ButtonSelectStyleHighlighted | ButtonFrameRound,
	}
}

const (
	ButtonSelectStyleSwell       = 0
	ButtonSelectStyleJumpOut     = 1
	ButtonSelectStyleHighlighted = 2
	ButtonFrameRounded           = 0
	ButtonFrameRound             = 2 << 3
	ButtonFrameSquare            = 3 << 3
	ButtonFrameOutlined          = 4 << 3

	// ButtonSmallScrollRadius defines the noise reduction zone. When the button is pressed
	// the position of the touchpad maybe moved a bit. We may reduce the noise by the
	// setting a radius of the finger move within initial touch point and don't consider
	// this move as a scroll signal. The flag ButtonSmallScrollRadius defines small radius (10)
	// which may increase the move sensitivity and the noise, but improve the scrolling experience
	ButtonSmallScrollRadius = 1 << 6
	// ButtonPresDelay sets up the press reaction to 100ms on the button click. It is useful
	// to use the feature when a button is placed on some scrolling group, so the button will
	// be pressed with some delay, not instantly
	ButtonPresDelay = 1 << 7
)

func (b *Button) InitButton(owner raywin.Container, r rl.RectangleInt32, text string, bs ButtonStyle, clickFn func()) {
	b.BaseComponent.Init(owner, b)
	b.SetBounds(r)
	b.text = text
	b.SetStyle(bs)
	delay := int64(0)
	if bs.Flags&ButtonPresDelay != 0 {
		delay = 100
	}
	if bs.Flags&ButtonSmallScrollRadius != 0 {
		b.InitPressor(10.0, delay, clickFn)
	} else {
		b.InitPressor(50.0, delay, clickFn)
	}
}

func (b *Button) SetStyle(bs ButtonStyle) {
	b.bs.Store(bs)
}

func (b *Button) Style() ButtonStyle {
	return b.bs.Load().(ButtonStyle)
}

func (b *Button) onFirstDraw(cc *raywin.CanvasContext) {
	bs := b.Style()
	b.textSize = rl.MeasureTextEx(bs.TextFont, b.text, bs.TextFontSize, 0)
}

func (b *Button) OnTPState(tps raywin.TPState) raywin.OnTPSResult {
	b.Pressor.OnTPState(tps)
	if b.Pressed() {
		b.pressedAt = tps.Millis
		b.fadeK = 1.0
		return raywin.OnTPSResultLocked
	}
	return raywin.OnTPSResultNA
}

func (b *Button) OnNewFrame(millis int64) {
	if !b.Pressed() && b.fadeK > 0.0 {
		b.fadeK = max(0.0, 1.0-float32(millis-b.pressedAt)/500)
	}
}

// Draw the drawing of ToggleButton notification
func (b *Button) Draw(cc *raywin.CanvasContext) {
	b.once.Do(func() { b.onFirstDraw(cc) })
	bs := b.Style()
	switch bs.Flags & 0x7 {
	case ButtonSelectStyleJumpOut:
		b.drawJumpedOut(cc)
	case ButtonSelectStyleSwell:
		b.drawSwallen(cc)
	case ButtonSelectStyleHighlighted:
		b.drawFaded(cc)
	default:
		b.drawSwallen(cc)
	}
}

func (b *Button) drawJumpedOut(cc *raywin.CanvasContext) {
	bs := b.Style()
	pr := b.Bounds()
	x, y := cc.PhysicalPointXY(0, 0)
	pr.X = x
	pr.Y = y
	dy := float32(pr.Y + pr.Height/2)
	if b.Pressed() {
		phr := cc.PhysicalRegion()
		rl.BeginScissorMode(phr.X-25, phr.Y-phr.Height, phr.Width+50, 2*phr.Height)
		defer rl.BeginScissorMode(phr.X, phr.Y, phr.Width, phr.Height)

		b.drawFrame(pr.ToFloat32(), bs.Color)
		pr.X -= 25
		pr.Y -= phr.Height
		pr.Width += 50
		pr.Height += 50

		for i := 0; i < 4; i++ {
			b.drawFrame(pr.ToFloat32(), rl.Fade(S.FrameSelectToneColor, 0.3))
			pr.X += 1
			pr.Y++
			pr.Height -= 2
			pr.Width -= 2
		}
		b.drawFrame(pr.ToFloat32(), S.FrameSelectToneColor)
		pr.X += 1
		pr.Y++
		pr.Height -= 2
		pr.Width -= 2
		b.drawFrame(pr.ToFloat32(), bs.SelectColor)
		dy := float32(pr.Y + pr.Height/2)
		center := rl.Vector2{X: float32(pr.X+pr.Width/2) - b.textSize.X/2, Y: dy - b.textSize.Y/2}
		rl.DrawTextEx(bs.TextFont, b.text, center, bs.TextFontSize, 0, bs.TextColor)
		return
	}
	b.drawFrame(pr.ToFloat32(), bs.Color)
	b.drawIcon(cc)
	center := rl.Vector2{X: float32(pr.X+pr.Width/2) - b.textSize.X/2, Y: dy - b.textSize.Y/2}
	rl.DrawTextEx(bs.TextFont, b.text, center, bs.TextFontSize, 0, bs.TextColor)
}

func (b *Button) drawFaded(cc *raywin.CanvasContext) {
	bs := b.Style()
	pr := b.Bounds()
	x, y := cc.PhysicalPointXY(0, 0)
	pr.X = x
	pr.Y = y
	dy := float32(pr.Y + pr.Height/2)
	col := bs.Color
	if b.Pressed() {
		col = bs.SelectColor
	} else {
		from := bs.SelectColor
		col.R = uint8(float32(from.R)*b.fadeK + float32(1.0-b.fadeK)*float32(col.R))
		col.G = uint8(float32(from.G)*b.fadeK + float32(1.0-b.fadeK)*float32(col.G))
		col.B = uint8(float32(from.B)*b.fadeK + float32(1.0-b.fadeK)*float32(col.B))
		col.A = uint8(float32(from.A)*b.fadeK + float32(1.0-b.fadeK)*float32(col.A))
	}
	b.drawFrame(pr.ToFloat32(), col)
	b.drawIcon(cc)
	center := rl.Vector2{X: float32(pr.X+pr.Width/2) - b.textSize.X/2, Y: dy - b.textSize.Y/2}
	rl.DrawTextEx(bs.TextFont, b.text, center, bs.TextFontSize, 0, bs.TextColor)
}

func (b *Button) drawHinted(cc *raywin.CanvasContext) {
	bs := b.Style()
	pr := b.Bounds()
	x, y := cc.PhysicalPointXY(0, 0)
	pr.X = x
	pr.Y = y
	dy := float32(pr.Y + pr.Height/2)
	col := bs.Color
	if b.Pressed() {
		phr := cc.PhysicalRegion()
		rl.BeginScissorMode(phr.X-5, phr.Y-phr.Height-5, phr.Width+10, 2*phr.Height+10)
		defer rl.BeginScissorMode(phr.X, phr.Y, phr.Width, phr.Height)

		pr.X -= 5
		pr.Y -= 5 + phr.Height
		pr.Width += 10
		pr.Height += 10 + pr.Height

		for i := 0; i < 4; i++ {
			b.drawFrame(pr.ToFloat32(), rl.Fade(S.FrameSelectToneColor, 0.3))
			pr.X += 1
			pr.Y++
			pr.Height -= 2
			pr.Width -= 2
		}
		b.drawFrame(pr.ToFloat32(), S.FrameSelectToneColor)
		pr.X += 1
		pr.Y++
		pr.Height -= 2
		pr.Width -= 2
		dy = float32(pr.Y + pr.Height/4)
		col = bs.SelectColor
	}
	b.drawFrame(pr.ToFloat32(), col)
	b.drawIcon(cc)
	center := rl.Vector2{X: float32(pr.X+pr.Width/2) - b.textSize.X/2, Y: dy - b.textSize.Y/2}
	rl.DrawTextEx(bs.TextFont, b.text, center, bs.TextFontSize, 0, bs.TextColor)
}

func (b *Button) drawSwallen(cc *raywin.CanvasContext) {
	bs := b.Style()
	pr := b.Bounds()
	x, y := cc.PhysicalPointXY(0, 0)
	pr.X = x
	pr.Y = y
	fs := bs.TextFontSize
	d := float32(0)
	if b.Pressed() {
		phr := cc.PhysicalRegion()
		rl.BeginScissorMode(phr.X-15, phr.Y-15, phr.Width+30, phr.Height+30)
		defer rl.BeginScissorMode(phr.X, phr.Y, phr.Width, phr.Height)

		pr.X -= 8
		pr.Y -= 8
		pr.Width += 16
		pr.Height += 16

		for i := 0; i < 4; i++ {
			b.drawFrame(pr.ToFloat32(), rl.Fade(S.FrameSelectToneColor, 0.3))
			pr.X += 1
			pr.Y++
			pr.Height -= 2
			pr.Width -= 2
		}
		b.drawFrame(pr.ToFloat32(), S.FrameSelectToneColor)
		pr.X += 1
		pr.Y++
		pr.Height -= 2
		pr.Width -= 2
		fs *= float32(pr.Width) / float32(pr.Width-16)
		d = -5.0
		b.drawFrame(pr.ToFloat32(), bs.SelectColor)
	} else {
		b.drawFrame(pr.ToFloat32(), bs.Color)
	}
	b.drawIcon(cc)
	center := rl.Vector2{X: float32(pr.X+pr.Width/2) - b.textSize.X/2 + d, Y: float32(pr.Y+pr.Height/2) - b.textSize.Y/2 + d}
	rl.DrawTextEx(bs.TextFont, b.text, center, fs, 0, bs.TextColor)
}

func (b *Button) drawFrame(r rl.Rectangle, col color.RGBA) {
	bs := b.Style()
	switch bs.Flags & 0x38 {
	case ButtonFrameRounded:
		rl.DrawRectangleRounded(r, 0.2, 5, col)
	case ButtonFrameRound:
		rl.DrawCircle(int32(r.X+r.Width/2), int32(r.Y+r.Height/2), r.Width/2, col)
	case ButtonFrameSquare:
		rl.DrawRectangle(int32(r.X), int32(r.Y), int32(r.Width), int32(r.Height), col)
		r.X += 3
		r.Y += 3
		r.Width -= 6
		r.Height -= 6
		rl.DrawRectangleLinesEx(r, 2.0, rl.Black)
	case ButtonFrameOutlined:
		if !b.Pressed() {
			rl.DrawRectangleRounded(r, 0.2, 5, bs.OutlineColor)
			r.X += 1.0
			r.Width -= 2.0
			r.Y += 1.0
			r.Height -= 2.0
		}
		rl.DrawRectangleRounded(r, 0.2, 5, col)
	}
}

func (b *Button) drawIcon(cc *raywin.CanvasContext) {
	bs := b.Style()
	if bs.Icon == "" {
		return
	}
	r := b.Bounds()
	x, y := cc.PhysicalPointXY(0, 0)
	tx, _ := raywin.GetIcon(bs.Icon)
	rl.DrawTexture(tx, x+r.Width/2-tx.Width/2, y+r.Height/2-tx.Height/2, rl.White)
}
