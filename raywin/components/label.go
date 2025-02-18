package components

import (
	raywin "github.com/dspasibenko/raywin-go/raywin"
	rl "github.com/gen2brain/raylib-go/raylib"
	"image/color"
	"sync"
)

// Label Component allows to draw a text on the screen with the specified
// font, color, and alignment
type Label struct {
	raywin.BaseComponent

	lock   sync.Mutex
	text   string
	cacheV *rl.Vector2
	cfg    LabelConfig
}

// LabelConfig allows to specify a label settings
type LabelConfig struct {
	alignment       int
	font            rl.Font
	fontSize        float32
	textColor       color.RGBA
	rect            rl.RectangleInt32
	backgroundColor color.RGBA
}

// DefaultLabelConfig returns the config with bottom left text alignment. The
// text will have the white color and size 32ppt. Default region is {0, 0, 100, 100}
func DefaultLabelConfig() LabelConfig {
	return LabelConfig{
		font:      raywin.SystemFont(),
		fontSize:  32,
		rect:      rl.RectangleInt32{X: 0, Y: 0, Width: 100, Height: 100},
		textColor: rl.White,
	}
}

// Alignment allows to set the text allignments (see Alignment flags in style)
func (lcfg LabelConfig) Alignment(flags int) LabelConfig {
	lcfg.alignment = flags
	return lcfg
}

// Font allows to change the label font
func (lcfg LabelConfig) Font(font rl.Font) LabelConfig {
	lcfg.font = font
	return lcfg
}

// FontSize specifies the font size
func (lcfg LabelConfig) FontSize(fontSize float32) LabelConfig {
	lcfg.fontSize = fontSize
	return lcfg
}

// Rectangle specifies the label bounds
func (lcfg LabelConfig) Rectangle(r rl.RectangleInt32) LabelConfig {
	lcfg.rect = r
	return lcfg
}

// Color specifies the text color
func (lcfg LabelConfig) Color(col color.RGBA) LabelConfig {
	lcfg.textColor = col
	return lcfg
}

// BackgroundColor specifies the backgour (transparent by default)
func (lcfg LabelConfig) BackgroundColor(col color.RGBA) LabelConfig {
	lcfg.backgroundColor = col
	return lcfg
}

// NewLabel creates a new label owned by `owner` with the text and lablel `cfg`
func NewLabel(owner raywin.Container, text string, cfg LabelConfig) (*Label, error) {
	l := &Label{}
	l.text = text
	l.cfg = cfg
	err := l.Init(owner, l)
	l.SetBounds(cfg.rect)
	return l, err
}

// SetText specifies the label text
func (l *Label) SetText(text string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.text = text
	l.cacheV = nil
}

// SetBounds allows to change the label region
func (l *Label) SetBounds(rect rl.RectangleInt32) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.BaseComponent.SetBounds(rect)
	l.cacheV = nil
}

// Draw draws the label on the screen
func (l *Label) Draw(cc *raywin.CanvasContext) {
	l.lock.Lock()
	defer l.lock.Unlock()
	txt := l.text
	if l.cacheV == nil {
		l.cacheV = &rl.Vector2{}
		r := l.Bounds()
		v := rl.MeasureTextEx(l.cfg.font, txt, l.cfg.fontSize, 0)
		switch l.cfg.alignment & 3 {
		case AlignBottom:
			l.cacheV.Y = float32(r.Height) - v.Y
		case AlignVCenter:
			l.cacheV.Y = (float32(r.Height) - v.Y) / 2
		}
		switch l.cfg.alignment & 12 {
		case AlignRight:
			l.cacheV.X = float32(r.Width) - v.X
		case AlignHCenter:
			l.cacheV.X = (float32(r.Width) - v.X) / 2
		}
	}
	x, y := cc.PhysicalPointXY(0, 0)
	if l.cfg.backgroundColor.A != 0 {
		b := l.Bounds()
		rl.DrawRectangle(x, y, b.Width, b.Height, l.cfg.backgroundColor)
	}
	rl.DrawTextEx(l.cfg.font, txt, rl.Vector2{X: float32(x) + l.cacheV.X, Y: float32(y) + l.cacheV.Y}, l.cfg.fontSize, 0, l.cfg.textColor)
}
