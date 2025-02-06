package raywin

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPressor_OnTPState(t *testing.T) {
	p := &Pressor{}
	assert.Equal(t, OnTPSResultNA, p.OnTPState(TPState{State: TPStateMoving}))
	assert.False(t, p.Pressed())
	assert.Equal(t, OnTPSResultNA, p.OnTPState(TPState{State: TPStateReleased}))
	assert.False(t, p.Pressed())
	assert.Equal(t, OnTPSResultLocked, p.OnTPState(TPState{State: TPStatePressed}))
	assert.True(t, p.Pressed())
	assert.Equal(t, OnTPSResultNA, p.OnTPState(TPState{State: TPStateReleased}))
	assert.False(t, p.Pressed())

	released := false
	p.InitPressor(10.0, 100, func() { released = true })
	assert.Equal(t, OnTPSResultNA, p.OnTPState(TPState{State: TPStatePressed}))
	assert.Equal(t, OnTPSResultNA, p.OnTPState(TPState{State: TPStatePressed, Millis: 99}))
	assert.Equal(t, OnTPSResultNA, p.OnTPState(TPState{State: TPStateMoving, Millis: 100, Pos: rl.Vector2{X: 100, Y: 100}}))
	assert.Equal(t, OnTPSResultNA, p.OnTPState(TPState{State: TPStatePressed, Millis: 110, Sequence: 1}))
	assert.Equal(t, OnTPSResultLocked, p.OnTPState(TPState{State: TPStatePressed, Millis: 220, Sequence: 1}))
	assert.Equal(t, OnTPSResultNA, p.OnTPState(TPState{State: TPStateMoving, Millis: 250, Pos: rl.Vector2{X: 15, Y: 15}}))
	assert.False(t, released)
	assert.Equal(t, OnTPSResultNA, p.OnTPState(TPState{State: TPStatePressed, Millis: 300, Sequence: 2}))
	assert.Equal(t, OnTPSResultLocked, p.OnTPState(TPState{State: TPStatePressed, Millis: 410, Sequence: 2}))
	assert.Equal(t, OnTPSResultNA, p.OnTPState(TPState{State: TPStateReleased, Millis: 420, Sequence: 2}))
	assert.True(t, released)
}
