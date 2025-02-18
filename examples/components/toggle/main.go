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
	// to use components with their style, register its outlet in the config
	cfg.FrameListener = components.DefaultStyleOutlet(cfg.DisplayConfig)
	raywin.Init(cfg)

	t, _ := components.NewToggle(raywin.RootContainer(), nil)
	t.SetBounds(rl.RectangleInt32{X: 100, Y: 100})

	ctx := context.NewSignalsContext(os.Interrupt, syscall.SIGTERM) // allow to close the window by Ctrl+C in terminal
	raywin.Run(ctx)
}
