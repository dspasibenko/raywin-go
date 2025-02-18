package main

import (
	"github.com/dspasibenko/raywin-go/pkg/golibs/context"
	"github.com/dspasibenko/raywin-go/raywin"
	"github.com/dspasibenko/raywin-go/raywin/components"
	rl "github.com/gen2brain/raylib-go/raylib"
	"os"
	"syscall"
)

func main() {
	cfg := raywin.DefaultConfig()
	cfg.IconsDir = "resources/icons"
	cfg.ResourceDir = "."
	cfg.RegularFontFileName = "resources/fonts/Roboto/Roboto-Medium.ttf"
	cfg.ItalicFontFileName = "resources/fonts/Roboto/Roboto-MediumItalic.ttf"
	// to use components with their style, register its outlet in the config
	cfg.FrameListener = components.DefaultStyleOutlet(cfg.DisplayConfig)
	raywin.Init(cfg)

	components.NewLabel(raywin.RootContainer(), "LeftTop",
		components.DefaultLabelConfig().
			Rectangle(rl.RectangleInt32{X: 10, Y: 10, Width: 200, Height: 70}).
			Alignment(components.AlignLeft|components.AlignTop).
			BackgroundColor(rl.Blue))

	components.NewLabel(raywin.RootContainer(), "CenterTop",
		components.DefaultLabelConfig().
			Rectangle(rl.RectangleInt32{X: 220, Y: 10, Width: 200, Height: 70}).
			Alignment(components.AlignHCenter|components.AlignTop).
			BackgroundColor(rl.Blue))

	components.NewLabel(raywin.RootContainer(), "RightTop",
		components.DefaultLabelConfig().
			Rectangle(rl.RectangleInt32{X: 430, Y: 10, Width: 200, Height: 70}).
			Alignment(components.AlignRight|components.AlignTop).
			BackgroundColor(rl.Blue))

	components.NewLabel(raywin.RootContainer(), "LeftVCenter",
		components.DefaultLabelConfig().
			Rectangle(rl.RectangleInt32{X: 10, Y: 100, Width: 200, Height: 70}).
			Alignment(components.AlignLeft|components.AlignVCenter).
			BackgroundColor(rl.Blue))

	components.NewLabel(raywin.RootContainer(), "CenterCenter",
		components.DefaultLabelConfig().
			Rectangle(rl.RectangleInt32{X: 220, Y: 100, Width: 200, Height: 70}).
			Alignment(components.AlignHCenter|components.AlignVCenter).
			BackgroundColor(rl.Blue))

	components.NewLabel(raywin.RootContainer(), "RightVCenter",
		components.DefaultLabelConfig().
			Rectangle(rl.RectangleInt32{X: 430, Y: 100, Width: 200, Height: 70}).
			Alignment(components.AlignRight|components.AlignVCenter).
			BackgroundColor(rl.Blue))

	components.NewLabel(raywin.RootContainer(), "LeftBottom",
		components.DefaultLabelConfig().
			Rectangle(rl.RectangleInt32{X: 10, Y: 200, Width: 200, Height: 70}).
			Alignment(components.AlignLeft|components.AlignBottom).
			BackgroundColor(rl.Blue))

	components.NewLabel(raywin.RootContainer(), "CenterBottom",
		components.DefaultLabelConfig().
			Rectangle(rl.RectangleInt32{X: 220, Y: 200, Width: 200, Height: 70}).
			Alignment(components.AlignHCenter|components.AlignBottom).
			BackgroundColor(rl.Blue))

	components.NewLabel(raywin.RootContainer(), "RightBottom",
		components.DefaultLabelConfig().
			Rectangle(rl.RectangleInt32{X: 430, Y: 200, Width: 200, Height: 70}).
			Alignment(components.AlignRight|components.AlignBottom).
			BackgroundColor(rl.Blue))

	ctx := context.NewSignalsContext(os.Interrupt, syscall.SIGTERM) // allow to close the window by Ctrl+C in terminal
	raywin.Run(ctx)
}
