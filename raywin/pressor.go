package raywin

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

// Pressor is the component which allows to handle touchpad state notifications.
// The component allows to hold the touchpad state in case of the touch point
// is moved within a radius after the initial press point. This allows to make
// the component be more responsive and user-friendly
//
// It also provides some press delay capabilities. If the setting is not 0, the
// pressor will not switch to the pressed state immediately, but wait, if the
// point released or start moving, to be sure it is a press action, but not something
// else.
type Pressor struct {
	presPos rl.Vector2
	radius  float32

	// press delay settings
	presDelayMillis int64
	presMillis      int64
	presSeq         int64

	pressed    bool
	onReleaseF func()
}

// InitPressor allows to set up the press radius sensitivity, the pressDelay (to reduce the noise)
// and the onReleaseF function for the notification when the pressor is released
func (p *Pressor) InitPressor(radius float32, pressDelay int64, onReleaseF func()) {
	p.radius = radius
	p.presDelayMillis = pressDelay
	p.onReleaseF = onReleaseF
}

// OnTPState implements Touchpadable
func (p *Pressor) OnTPState(tps TPState) OnTPSResult {
	switch tps.State {
	case TPStatePressed:
		if tps.Sequence != p.presSeq {
			p.presSeq = tps.Sequence
			p.presMillis = tps.Millis
		}
		p.pressed = tps.Millis-p.presMillis >= p.presDelayMillis
		p.presPos = tps.Pos
	case TPStateMoving:
		if !p.pressed {
			break
		}
		dx := p.presPos.X - tps.Pos.X
		dy := p.presPos.Y - tps.Pos.Y
		dist := float32(math.Sqrt(float64(dx*dx) + float64(dy*dy)))
		p.pressed = dist < p.radius
	case TPStateReleased:
		if p.pressed && p.onReleaseF != nil {
			p.onReleaseF()
		}
		p.pressed = false
	}
	if p.pressed {
		return OnTPSResultLocked
	}
	return OnTPSResultNA
}

// Pressed returns whether the Pressor is pressed or not
func (p *Pressor) Pressed() bool {
	return p.pressed
}
