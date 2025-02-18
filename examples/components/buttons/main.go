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

	components.NewButton(raywin.RootContainer(), rl.RectangleInt32{X: 50, Y: 100, Width: 100, Height: 60}, "hello", components.DialogButtonStyle(), nil)
	components.NewButton(raywin.RootContainer(), rl.RectangleInt32{X: 200, Y: 100, Width: 100, Height: 60}, "cancel", components.DialogButtonCancelStyle(), nil)
	components.NewButton(raywin.RootContainer(), rl.RectangleInt32{X: 350, Y: 100, Width: 100, Height: 60}, "Ok", components.DialogButtonOkStyle(), nil)
	components.NewButton(raywin.RootContainer(), rl.RectangleInt32{X: 500, Y: 100, Width: 100, Height: 60}, "Control", components.DialogButtonControlStyle(), nil)
	components.NewButton(raywin.RootContainer(), rl.RectangleInt32{X: 650, Y: 100, Width: 60, Height: 60}, "", components.DialogButtonCloseStyle(), nil)
	components.NewButton(raywin.RootContainer(), rl.RectangleInt32{X: 50, Y: 200, Width: 100, Height: 100}, "", components.DialogButtonTransparrentStyle(), nil)

	ctx := context.NewSignalsContext(os.Interrupt, syscall.SIGTERM) // allow to close the window by Ctrl+C in terminal
	raywin.Run(ctx)
}
