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

	// Scrolling
	ScrollBarDarkColor       rl.Color
	ScrollBarLightColor      rl.Color
	ScrollBarThiknessMm      float32
	ScrollBarOffsetMm        float32
	ScrollBarDisappearMillis int64

	// Toggle
	TogglePressMillis int64
	TogglePressRadius float32
	ToggleHeightMm    float32
	ToggleWidthMm     float32
	ToggleSpaceMm     float32
	ToggleOnColor     rl.Color
	ToggleOffColor    rl.Color

	// Dimensions
	PPcm  float32
	PPI   float32
	Scale float32
}

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

		// Scrolling
		ScrollBarDarkColor:       color.RGBA{0, 0, 0, 90},
		ScrollBarLightColor:      color.RGBA{100, 100, 100, 90},
		ScrollBarThiknessMm:      2.0,
		ScrollBarOffsetMm:        0.7,
		ScrollBarDisappearMillis: 1000,

		// Toggle
		TogglePressMillis: 50,
		TogglePressRadius: 10.0,
		ToggleHeightMm:    8.0,
		ToggleWidthMm:     13.0,
		ToggleSpaceMm:     0.7,
		ToggleOnColor:     color.RGBA{16, 173, 55, 255},
		ToggleOffColor:    rl.DarkGray,

		// Dimensions
		PPcm:  cfg.PPI / 2.54,
		PPI:   cfg.PPI,
		Scale: 100.0,
	}
}
