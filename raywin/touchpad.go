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
package raywin

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type (
	// TPState describes the current touchpad state. The event describes
	// most relevant state for the touchpad. The Millis field contains the timestamp
	// whe the state is observed (not when the touchpad switched to the state!)
	TPState struct {
		// State contains the state of the touchpad
		State int
		// Pos defines the position of the first touched point (if any)
		Pos rl.Vector2
		// Millis contains the timestamp when the state is read last time
		// (latest time, not the time when the state is set!)
		Millis int64
		// Sequence is the state unique identifier. Every new state
		// has a new montonically increasing sequence
		Sequence int64
	}
	OnTPSResult int

	Touchpadable interface {
		// OnTPState is called every frame with the current touchpad State
		// if any. The method must return true, if the control locking the events.
		// This case all following touchpad events will be sent to the component
		// only, and not to other ones.
		OnTPState(tps TPState) OnTPSResult
	}
)

const (
	TPStateNA = iota
	// TPStatePressed indicates that the touchpad is pressed in the position and
	// the position is not changed after the event
	TPStatePressed
	// TPStateMoving indicates the fact that the touchpad position is moved (but not released)
	TPStateMoving
	// TPStateReleased reports the position when the touchpad was released
	TPStateReleased
)

const (
	// OnTPSResultNA tells display that the TPState is not handled and another
	// component may be notified
	OnTPSResultNA = OnTPSResult(0)
	// OnTPSResultLocked tells display that the component locked focus and it
	// will handle further events
	OnTPSResultLocked = OnTPSResult(1)
	// OnTPSResultStop tells display that the component doesn't lock the focus,
	// but requests not processing with other components
	OnTPSResultStop = OnTPSResult(2)
)

type touchPad struct {
	state  int
	pos    rl.Vector2
	millis int64
	seq    int64
}

const (
	tpsInit = iota
	tpsPressed
	tpsMoving
	tpsReleased
)

func (tps TPState) PosXY() (int32, int32) {
	return int32(tps.Pos.X), int32(tps.Pos.Y)
}

func (tp *touchPad) tpState() TPState {
	res := TPState{Pos: tp.pos, State: TPStateNA, Millis: tp.millis, Sequence: tp.seq}
	switch tp.state {
	case tpsPressed:
		res.State = TPStatePressed
	case tpsReleased:
		res.State = TPStateReleased
	case tpsMoving:
		res.State = TPStateMoving
	}
	return res
}

func (tp *touchPad) onNewFrame(millis int64) TPState {
	tp.millis = millis
	prevState := tp.state
	if rl.IsMouseButtonDown(rl.MouseLeftButton) {
		switch tp.state {
		case tpsInit, tpsReleased:
			tp.state = tpsPressed
		case tpsPressed:
			if !IsEmpty(rl.GetMouseDelta()) {
				tp.state = tpsMoving
			}
		}
		tp.pos = rl.GetMousePosition()
	} else {
		switch tp.state {
		case tpsMoving, tpsPressed:
			tp.state = tpsReleased
		case tpsReleased:
			tp.state = tpsInit
		}
	}
	if prevState != tp.state {
		tp.seq++
	}
	return tp.tpState()
}
