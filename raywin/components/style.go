package components

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
	"github.com/dspasibenko/raywin-go/raywin"
	rl "github.com/gen2brain/raylib-go/raylib"
	"image/color"
	"sync/atomic"
)

// Style struct enables to specify different configurations for the components. It is used
// when drawing components from the display drawing loop, so the structure MUST be used for
// reading only in the callbacks called by raywin (Draw, OnNewFrame, DrawAfter etc.)
type Style struct {
	// Common
	FrameColor            rl.Color
	DialogBackgroundLight rl.Color
	DialogBackgroundDark  rl.Color
	FrameSelectColor      rl.Color
	FrameSelectToneColor  rl.Color
	FrameShadeColor       rl.Color
	TransparentColor      rl.Color
	OutlineColor          rl.Color

	// Scrolling
	ScrollBarDarkColor       rl.Color
	ScrollBarLightColor      rl.Color
	ScrollBarThiknessMm      float32
	ScrollBarOffsetMm        float32
	ScrollBarDisappearMillis int64

	// Buttons
	ButtonJumpOutCoef float32
	ButtonSwallenCoef float32

	// Toggle
	TogglePressMillis int64
	TogglePressRadius float32
	ToggleHeightMm    float32
	ToggleWidthMm     float32
	ToggleSpaceMm     float32
	ToggleOnColor     rl.Color
	ToggleOffColor    rl.Color

	CurorWidth float32

	// EditBox
	EditBoxHeightMm       float32
	EditBoxSpacerMm       float32
	EditBoxFontSize       float32
	EditBoxCursorWidthMm  float32
	EditBoxTextColor      rl.Color
	EditBoxBackgoundColor rl.Color
	EditBoxOutlineColor   rl.Color

	// Dimensions
	PPcm  float32
	PPI   float32
	Scale float32
}

const (
	AlignBottom  = 0
	AlignVCenter = 1 << 0
	AlignTop     = 2 << 0
	AlignLeft    = 0
	AlignHCenter = 1 << 2
	AlignRight   = 2 << 2
)

// S is the current Style. It MUST NOT be accessed or modified outside the drawing loop goroutine,
// which invokes raywin callbacks. To change the current style, use SetStyle() instead.
var S Style
var s atomic.Value

type defaultOutlet struct{}

// OnNewFrame implements FrameListener
func (do defaultOutlet) OnNewFrame(millis int64) {
	S = s.Load().(Style)
}

// DefaultStyleOutlet returns the FrameListener for managing Style settings within raywin. Use this
// function to receive FrameListener when setting up Config.FrameListener for raywin.Init() if needed
func DefaultStyleOutlet(cfg raywin.DisplayConfig) raywin.FrameListener {
	SetStyle(initDefaultStyle(cfg))
	do := defaultOutlet{}
	do.OnNewFrame(0)
	return do
}

// SetStyle allows to change the style dyncamically
func SetStyle(newStyle Style) error {
	s.Store(newStyle)
	return nil
}

func initDefaultStyle(cfg raywin.DisplayConfig) Style {
	return Style{
		FrameColor:            color.RGBA{220, 220, 220, 255},
		DialogBackgroundLight: color.RGBA{7, 83, 97, 255},
		DialogBackgroundDark:  color.RGBA{2, 41, 48, 255},
		FrameSelectColor:      color.RGBA{220, 220, 220, 255},
		FrameSelectToneColor:  color.RGBA{189, 241, 252, 255},
		FrameShadeColor:       rl.Gray,
		TransparentColor:      color.RGBA{0, 0, 0, 0},
		OutlineColor:          color.RGBA{14, 110, 138, 255},

		// Scrolling
		ScrollBarDarkColor:       color.RGBA{0, 0, 0, 90},
		ScrollBarLightColor:      color.RGBA{100, 100, 100, 90},
		ScrollBarThiknessMm:      2.0,
		ScrollBarOffsetMm:        0.7,
		ScrollBarDisappearMillis: 1000,

		// Buttons
		ButtonJumpOutCoef: 1.7,
		ButtonSwallenCoef: 0.3, // 30%

		// Toggle
		TogglePressMillis: 50,
		TogglePressRadius: 10.0,
		ToggleHeightMm:    12.0,
		ToggleWidthMm:     20.0,
		ToggleSpaceMm:     0.7,
		ToggleOnColor:     color.RGBA{16, 173, 55, 255},
		ToggleOffColor:    rl.DarkGray,

		CurorWidth: 6.0,

		// EditBox
		EditBoxHeightMm:       12.0,
		EditBoxSpacerMm:       1.4,
		EditBoxFontSize:       60.0,
		EditBoxCursorWidthMm:  0.8,
		EditBoxTextColor:      rl.Color{R: 240, G: 252, B: 255, A: 255},
		EditBoxBackgoundColor: rl.Color{R: 73, G: 85, B: 79, A: 255},
		EditBoxOutlineColor:   rl.Color{R: 147, G: 169, B: 158, A: 255},

		// Dimensions
		PPcm:  cfg.PPI / 2.54,
		PPI:   cfg.PPI,
		Scale: 100.0,
	}
}
