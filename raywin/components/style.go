package components

import (
	"github.com/dspasibenko/raywin-go/raywin"
	rl "github.com/gen2brain/raylib-go/raylib"
	"image/color"
)

type Style struct {
	// Common
	FrameColor            rl.Color
	DialogBackgroundLight rl.Color
	DialogBackgroundDark  rl.Color

	// Scrolling
	ScrollBarDarkColor       rl.Color
	ScrollBarLightColor      rl.Color
	ScrollBarThiknessMm      float32
	ScrollBarOffsetMm        float32
	ScrollBarDisappearMillis int64

	// Dimensions
	PPcm  float32
	PPI   float32
	Scale float32
}

var S Style

func InitDefaultStyle(cfg raywin.DisplayConfig) error {
	S = Style{
		FrameColor:            color.RGBA{220, 220, 220, 255},
		DialogBackgroundLight: color.RGBA{7, 83, 97, 255},
		DialogBackgroundDark:  color.RGBA{2, 41, 48, 255},

		// Scrolling
		ScrollBarDarkColor:       color.RGBA{0, 0, 0, 70},
		ScrollBarLightColor:      color.RGBA{100, 100, 100, 70},
		ScrollBarThiknessMm:      2.5,
		ScrollBarOffsetMm:        1.0,
		ScrollBarDisappearMillis: 1000,

		// Dimensions
		PPcm:  cfg.PPI / 2.54,
		PPI:   cfg.PPI,
		Scale: 100.0,
	}
	return nil
}
