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
	// when the state is observed (not when the touchpad switched to the state!)
	TPState struct {
		// State contains the state of the touchpad
		State int
		// Pos defines the position of the first touched point (if any)
		Pos rl.Vector2
		// Millis contains the timestamp when the state is read last time
		// (latest time, not the time when the state is set!)
		Millis int64
		// Sequence is the state unique identifier. Every new state
		// has a new monotonically increasing sequence
		Sequence int64
	}

	// OnTPSResult the result which will be returned by the OnTPState by the Touchpadable
	// component. Please see the results descriptions in the constants below
	OnTPSResult int

	// Touchpadable interface maybe implemented by a component to let raywin-go know
	// that the component wants to react on the touchpad events.
	Touchpadable interface {
		// OnTPState is called every frame with the current touchpad State
		// if any. The method must return OnTPSResult value(see below).
		OnTPState(tps TPState) OnTPSResult
	}
)

const (
	// TPStateNA indicates that no points on the touchpad are currently pressed.
	TPStateNA = iota

	// TPStatePressed indicates that the touchpad is pressed at a specific position,
	// and the position has not changed since the event. If the point is moved and then
	// stops without being released, the state will transition to TPStateMoving and
	// will not return to TPStatePressed.
	TPStatePressed

	// TPStateMoving indicates that the touchpad is pressed and the position has moved
	// (but the touchpad is not yet released). Even if the finger moves (or not), but remains pressed,
	// this state will be reported instead of TPStatePressed.
	TPStateMoving

	// TPStateReleased reports the position where the touchpad was released.
	TPStateReleased
)

const (
	// OnTPSResultNA tells that the TPState is not handled and another
	// component may be notified about the event
	OnTPSResultNA = OnTPSResult(0)

	// OnTPSResultLocked indicates that the component has locked focus
	// and will handle all further touchpad events. These events will
	// be sent exclusively to the component until it returns a different
	// result or is closed. During this time, other components will not
	// receive touchpad event notifications.
	OnTPSResultLocked = OnTPSResult(1)

	// OnTPSResultStop indicates that the component does not hold focus,
	// but the touchpad event should not be passed to other components.
	// This signals raywin to stop the cycle and refrain from notifying
	// other components about the touchpad events BUT ONLY FOR THE CURRENT FRAME!
	//
	// The difference between this state and OnTPSResultLocked is as follows:
	// OnTPSResultStop terminates only the current cycle of notifications,
	// whereas OnTPSResultLocked prevents the cycle from starting in future
	// frames. In the locked state, only the locked component is notified
	// until it indicates otherwise.
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

func (tp *touchPad) onNewFrame(millis int64, proxy rlProxy) TPState {
	tp.millis = millis
	prevState := tp.state
	if proxy.isMouseButtonDown(rl.MouseLeftButton) {
		switch tp.state {
		case tpsInit, tpsReleased:
			tp.state = tpsPressed
		case tpsPressed:
			if !IsEmpty(proxy.getMouseDelta()) {
				tp.state = tpsMoving
			}
		}
		tp.pos = proxy.getMousePosition()
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
