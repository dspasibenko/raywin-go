package components

import (
	"github.com/dspasibenko/raywin-go/pkg/golibs"
	"github.com/dspasibenko/raywin-go/raywin"
	rl "github.com/gen2brain/raylib-go/raylib"
	"image/color"
	"sync/atomic"
)

// Button is the Component which represents a pressable component with different design styles
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

// ButtonStyle allows to specify the style of the button
type ButtonStyle struct {
	textFont     rl.Font
	textFontSize float32
	color        color.RGBA
	outlineColor color.RGBA
	textColor    color.RGBA
	selectColor  color.RGBA
	icon         string
	flags        int // See constants below (ButtonSelectStyleJumpOut etc.)
}

// TextFont specifies the font which will be used for the button text
func (bs ButtonStyle) TextFont(textFont rl.Font) ButtonStyle {
	bs.textFont = textFont
	return bs
}

// TextFontSize specifies the font size used for the button text
func (bs ButtonStyle) TextFontSize(textFontSize float32) ButtonStyle {
	bs.textFontSize = textFontSize
	return bs
}

// Color specifies the button color background
func (bs ButtonStyle) Color(color rl.Color) ButtonStyle {
	bs.color = color
	return bs
}

// OutlineColor specifies the frame color
func (bs ButtonStyle) OutlineColor(color rl.Color) ButtonStyle {
	bs.outlineColor = color
	return bs
}

// TextColor specifies the text color
func (bs ButtonStyle) TextColor(color rl.Color) ButtonStyle {
	bs.textColor = color
	return bs
}

// SelectColor specifies the button color when pressed
func (bs ButtonStyle) SelectColor(color rl.Color) ButtonStyle {
	bs.selectColor = color
	return bs
}

// Icon specifies the icon name which will be on the button
func (bs ButtonStyle) Icon(iconName string) ButtonStyle {
	bs.icon = iconName
	return bs
}

// Flags provides the button flags (see below in the file)
func (bs ButtonStyle) Flags(flags int) ButtonStyle {
	bs.flags = flags
	return bs
}

// DialogButtonStyle - just standard button style which jumps out when it is pressed
func DialogButtonStyle() ButtonStyle {
	return ButtonStyle{
		textFont:     raywin.SystemItalicFont(),
		textFontSize: 50,
		color:        S.DialogBackgroundDark,
		outlineColor: S.OutlineColor,
		selectColor:  S.DialogBackgroundLight,
		textColor:    rl.White,
		flags:        ButtonSelectStyleJumpOut | ButtonFrameOutlined,
	}
}

// DialogButtonCancelStyle offers "cancel" button for dialogs
func DialogButtonCancelStyle() ButtonStyle {
	return ButtonStyle{
		textFont:     raywin.SystemFont(),
		textFontSize: 30,
		color:        color.RGBA{82, 2, 2, 255},
		selectColor:  color.RGBA{107, 2, 2, 255},
		textColor:    rl.White,
		flags:        ButtonSelectStyleSwell | ButtonFrameRounded,
	}
}

// DialogButtonOkStyle offers "ok" button for dialogs
func DialogButtonOkStyle() ButtonStyle {
	return ButtonStyle{
		textFont:     raywin.SystemFont(),
		textFontSize: 30,
		color:        color.RGBA{4, 51, 38, 255},
		selectColor:  color.RGBA{6, 71, 53, 255},
		textColor:    rl.White,
		flags:        ButtonSelectStyleSwell | ButtonFrameRounded,
	}
}

// DialogButtonControlStyle offers a button for controls button style
func DialogButtonControlStyle() ButtonStyle {
	return ButtonStyle{
		textFont:     raywin.SystemItalicFont(),
		textFontSize: 25,
		color:        S.DialogBackgroundDark,
		outlineColor: S.OutlineColor,
		selectColor:  S.DialogBackgroundLight,
		textColor:    rl.White,
		flags:        ButtonSelectStyleJumpOut | ButtonFrameOutlined,
	}
}

// DialogButtonCloseStyle style for the close dialog button style
func DialogButtonCloseStyle() ButtonStyle {
	return ButtonStyle{
		textFont:     raywin.SystemFont(),
		textFontSize: 0,
		color:        S.DialogBackgroundDark,
		outlineColor: S.OutlineColor,
		selectColor:  S.DialogBackgroundLight,
		textColor:    rl.White,
		icon:         "x-white",
		flags:        ButtonSelectStyleHighlighted | ButtonFrameOutlined,
	}
}

// DialogButtonTransparrentStyle offers a rounded when pressed an icon button style
func DialogButtonTransparrentStyle() ButtonStyle {
	return ButtonStyle{
		textFont:     raywin.SystemFont(),
		textFontSize: 0,
		color:        S.TransparentColor,
		outlineColor: S.TransparentColor,
		selectColor:  rl.Fade(rl.Orange, 0.4),
		textColor:    rl.White,
		icon:         "airplane-yellow",
		flags:        ButtonSelectStyleHighlighted | ButtonFrameRound,
	}
}

const (
	// ButtonSelectStyleSwell flag allows to change the button size when it is pressed
	ButtonSelectStyleSwell = 0
	// ButtonSelectStyleJumpOut flag makes the button pops up when it is pressed
	ButtonSelectStyleJumpOut = 1
	// ButtonSelectStyleHighlighted flag makes the button be highlighted when pressed
	ButtonSelectStyleHighlighted = 2

	// ButtonFrameRounded flag says the frame is rectangular with rounded corners
	ButtonFrameRounded = 0
	// ButtonFrameRound flag makes the button be round
	ButtonFrameRound = 2 << 3
	// ButtonFrameSquare flag makes the button squared
	ButtonFrameSquare = 3 << 3
	// ButtonFrameOutlined flag allows to draw frame for the button
	ButtonFrameOutlined = 4 << 3

	// ButtonSmallNoiseRadius defines the noise reduction zone. When the button is pressed
	// the position of the touchpad maybe moved a bit. We may reduce the noise by the
	// setting a radius of the finger move within initial touch point and don't consider
	// this move as a scroll signal. The flag ButtonSmallNoiseRadius defines small radius (10)
	// which may increase the move sensitivity and the noise, but improve the scrolling experience
	ButtonSmallNoiseRadius = 1 << 6

	// ButtonPresDelay sets up the press reaction to 70ms on the button click. It is useful
	// to use the feature when a button is placed on some scrolling group, so the button will
	// be pressed with some delay, not instantly
	ButtonPresDelay = 1 << 7
)

// NewButton returns new Button. After the creation, no Init() calls are needed.
func NewButton(owner raywin.Container, r rl.RectangleInt32, text string, bs ButtonStyle, clickFn func()) (*Button, error) {
	b := &Button{}
	if err := b.BaseComponent.Init(owner, b); err != nil {
		return nil, err
	}
	b.SetBounds(r)
	b.text = text
	b.SetStyle(bs)
	delay := int64(0)
	if bs.flags&ButtonPresDelay != 0 {
		delay = 70
	}
	if bs.flags&ButtonSmallNoiseRadius != 0 {
		b.InitPressor(10.0, delay, clickFn)
	} else {
		b.InitPressor(50.0, delay, clickFn)
	}
	return b, nil
}

// SetStyle set button style
func (b *Button) SetStyle(bs ButtonStyle) {
	b.bs.Store(bs)
}

// Style returns ButtonStyle
func (b *Button) Style() ButtonStyle {
	return b.bs.Load().(ButtonStyle)
}

func (b *Button) onFirstDraw(cc *raywin.CanvasContext) {
	bs := b.Style()
	b.textSize = rl.MeasureTextEx(bs.textFont, b.text, bs.textFontSize, 0)
}

// OnTPState the TouchPad notification
func (b *Button) OnTPState(tps raywin.TPState) raywin.OnTPSResult {
	b.Pressor.OnTPState(tps)
	if b.Pressed() {
		b.pressedAt = tps.Millis
		b.fadeK = 1.0
		return raywin.OnTPSResultLocked
	}
	return raywin.OnTPSResultNA
}

// OnNewFrame is the on new notification
func (b *Button) OnNewFrame(millis int64) {
	if !b.Pressed() && b.fadeK > 0.0 {
		b.fadeK = max(0.0, 1.0-float32(millis-b.pressedAt)/500)
	}
}

// Draw the drawing of ToggleButton notification
func (b *Button) Draw(cc *raywin.CanvasContext) {
	b.once.Do(func() { b.onFirstDraw(cc) })
	bs := b.Style()
	switch bs.flags & 0x7 {
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
		dx := int32(float32(phr.Width) * S.ButtonJumpOutCoef)
		dy := int32(float32(phr.Height) * S.ButtonJumpOutCoef)
		rl.BeginScissorMode(phr.X-(dx-phr.Width)/2, phr.Y-phr.Height, phr.Width+dx, phr.Height*2)
		defer rl.BeginScissorMode(phr.X, phr.Y, phr.Width, phr.Height)

		b.drawFrame(pr.ToFloat32(), bs.color)
		pr.X -= (dx - phr.Width) / 2
		pr.Y -= phr.Height
		pr.Width = dx
		pr.Height = dy

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
		b.drawFrame(pr.ToFloat32(), bs.selectColor)
		dy2 := float32(pr.Y + pr.Height/2)
		center := rl.Vector2{X: float32(pr.X+pr.Width/2) - b.textSize.X/2, Y: dy2 - b.textSize.Y/2}
		rl.DrawTextEx(bs.textFont, b.text, center, bs.textFontSize, 0, bs.textColor)
		return
	}
	b.drawFrame(pr.ToFloat32(), bs.color)
	b.drawIcon(cc)
	center := rl.Vector2{X: float32(pr.X+pr.Width/2) - b.textSize.X/2, Y: dy - b.textSize.Y/2}
	rl.DrawTextEx(bs.textFont, b.text, center, bs.textFontSize, 0, bs.textColor)
}

func (b *Button) drawFaded(cc *raywin.CanvasContext) {
	bs := b.Style()
	pr := b.Bounds()
	x, y := cc.PhysicalPointXY(0, 0)
	pr.X = x
	pr.Y = y
	dy := float32(pr.Y + pr.Height/2)
	col := bs.color
	if b.Pressed() {
		col = bs.selectColor
	} else {
		from := bs.selectColor
		col.R = uint8(float32(from.R)*b.fadeK + float32(1.0-b.fadeK)*float32(col.R))
		col.G = uint8(float32(from.G)*b.fadeK + float32(1.0-b.fadeK)*float32(col.G))
		col.B = uint8(float32(from.B)*b.fadeK + float32(1.0-b.fadeK)*float32(col.B))
		col.A = uint8(float32(from.A)*b.fadeK + float32(1.0-b.fadeK)*float32(col.A))
	}
	b.drawFrame(pr.ToFloat32(), col)
	b.drawIcon(cc)
	center := rl.Vector2{X: float32(pr.X+pr.Width/2) - b.textSize.X/2, Y: dy - b.textSize.Y/2}
	rl.DrawTextEx(bs.textFont, b.text, center, bs.textFontSize, 0, bs.textColor)
}

func (b *Button) drawSwallen(cc *raywin.CanvasContext) {
	bs := b.Style()
	pr := b.Bounds()
	x, y := cc.PhysicalPointXY(0, 0)
	pr.X = x
	pr.Y = y
	fs := bs.textFontSize
	d := float32(0)
	if b.Pressed() {
		phr := cc.PhysicalRegion()
		diff := int32(float32(phr.Width) * S.ButtonSwallenCoef)
		rl.BeginScissorMode(phr.X-diff/2, phr.Y-diff/2, phr.Width+diff, phr.Height+diff)
		defer rl.BeginScissorMode(phr.X, phr.Y, phr.Width, phr.Height)

		pr.X -= diff / 2
		pr.Y -= diff / 2
		pr.Width += diff
		pr.Height += diff

		for i := 0; i < 4; i++ {
			b.drawFrame(pr.ToFloat32(), rl.Fade(S.FrameSelectToneColor, 0.3))
			pr.X += 1
			pr.Y += 1
			pr.Height -= 2
			pr.Width -= 2
		}
		b.drawFrame(pr.ToFloat32(), S.FrameSelectToneColor)
		pr.X += 1
		pr.Y += 1
		pr.Height -= 2
		pr.Width -= 2
		fs *= float32(pr.Width) / float32(pr.Width-diff/2)
		d = -5.0
		b.drawFrame(pr.ToFloat32(), bs.selectColor)
	} else {
		b.drawFrame(pr.ToFloat32(), bs.color)
	}
	b.drawIcon(cc)
	center := rl.Vector2{X: float32(pr.X+pr.Width/2) - b.textSize.X/2 + d, Y: float32(pr.Y+pr.Height/2) - b.textSize.Y/2 + d}
	rl.DrawTextEx(bs.textFont, b.text, center, fs, 0, bs.textColor)
}

func (b *Button) drawFrame(r rl.Rectangle, col color.RGBA) {
	bs := b.Style()
	switch bs.flags & 0x38 {
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
			rl.DrawRectangleRounded(r, 0.2, 5, bs.outlineColor)
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
	if bs.icon == "" {
		return
	}
	r := b.Bounds()
	x, y := cc.PhysicalPointXY(0, 0)
	tx, _ := raywin.GetIcon(bs.icon)
	rl.DrawTexture(tx, x+r.Width/2-tx.Width/2, y+r.Height/2-tx.Height/2, rl.White)
}
