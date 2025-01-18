package raywin

// Copyright 2025 Dmitry Spasibenko
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_touchpad_onNewFrame(t *testing.T) {
	tp := &touchPad{}
	pxy := &testProxy{}
	s := tp.onNewFrame(1, pxy)
	assert.Equal(t, TPState{Pos: pxy.mousePos, State: TPStateNA, Millis: 1, Sequence: 0}, s)
	s = tp.tpState()
	assert.Equal(t, TPState{Pos: pxy.mousePos, State: TPStateNA, Millis: 1, Sequence: 0}, s)

	s = tp.onNewFrame(2, pxy)
	assert.Equal(t, TPState{Pos: pxy.mousePos, State: TPStateNA, Millis: 2, Sequence: 0}, s)

	pxy.mousePos = rl.Vector2{X: 1, Y: 2}
	s = tp.onNewFrame(3, pxy)
	assert.Equal(t, TPState{Pos: pxy.mousePos, State: TPStatePressed, Millis: 3, Sequence: 1}, s)
	s = tp.onNewFrame(4, pxy)
	assert.Equal(t, TPState{Pos: pxy.mousePos, State: TPStatePressed, Millis: 4, Sequence: 1}, s)

	pxy.mouseDiff = rl.Vector2{X: 1, Y: 1}
	s = tp.onNewFrame(5, pxy)
	assert.Equal(t, TPState{Pos: pxy.mousePos, State: TPStateMoving, Millis: 5, Sequence: 2}, s)
	s = tp.onNewFrame(6, pxy)
	assert.Equal(t, TPState{Pos: pxy.mousePos, State: TPStateMoving, Millis: 6, Sequence: 2}, s)
	pxy.mouseDiff = rl.Vector2{}
	s = tp.onNewFrame(7, pxy)
	// if we switched to moving, no changes like pressed after that until it is released
	assert.Equal(t, TPState{Pos: pxy.mousePos, State: TPStateMoving, Millis: 7, Sequence: 2}, s)

	// Now release it!
	prevPos := pxy.mousePos
	pxy.mousePos = rl.Vector2{}
	s = tp.onNewFrame(8, pxy)
	assert.Equal(t, TPState{Pos: prevPos, State: TPStateReleased, Millis: 8, Sequence: 3}, s)
	s = tp.onNewFrame(9, pxy)
	assert.Equal(t, TPState{Pos: prevPos, State: TPStateNA, Millis: 9, Sequence: 4}, s)
	s = tp.onNewFrame(10, pxy)
	assert.Equal(t, TPState{Pos: prevPos, State: TPStateNA, Millis: 10, Sequence: 4}, s)

	// Ok Press and release:
	pxy.mousePos = prevPos
	s = tp.onNewFrame(11, pxy)
	assert.Equal(t, TPState{Pos: pxy.mousePos, State: TPStatePressed, Millis: 11, Sequence: 5}, s)
	pxy.mousePos = rl.Vector2{}
	s = tp.onNewFrame(12, pxy)
	assert.Equal(t, TPState{Pos: prevPos, State: TPStateReleased, Millis: 12, Sequence: 6}, s)
	s = tp.onNewFrame(13, pxy)
	assert.Equal(t, TPState{Pos: prevPos, State: TPStateNA, Millis: 13, Sequence: 7}, s)
}
