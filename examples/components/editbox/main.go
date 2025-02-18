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

	eb, _ := components.NewEditBox(raywin.RootContainer())
	eb.SetBounds(rl.RectangleInt32{X: 50, Y: 50, Width: 200, Height: 100})

	components.NewButton(raywin.RootContainer(), rl.RectangleInt32{X: 50, Y: 200, Width: 60, Height: 60}, "H", components.DialogButtonStyle(), func() {
		eb.SetText(eb.Text() + "H")
	})

	components.NewButton(raywin.RootContainer(), rl.RectangleInt32{X: 150, Y: 200, Width: 60, Height: 60}, "<-", components.DialogButtonStyle(), func() {
		s := eb.Text()
		if len(s) > 0 {
			s = s[:len(s)-1]
		}
		eb.SetText(s)
	})

	ctx := context.NewSignalsContext(os.Interrupt, syscall.SIGTERM) // allow to close the window by Ctrl+C in terminal
	raywin.Run(ctx)
}
