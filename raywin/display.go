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
	"fmt"
	"github.com/dspasibenko/raywin-go/pkg/golibs/errors"
	"github.com/dspasibenko/raywin-go/pkg/golibs/logging"
	rl "github.com/gen2brain/raylib-go/raylib"
	"sync/atomic"
	"time"
)

type display struct {
	// fc frame counter
	fc  uint64
	cfg DisplayConfig

	running int32
	logger  logging.Logger

	cc     *CanvasContext
	tp     *touchPad
	millis atomic.Int64

	root        rootContainer
	tpsAcceptor Component
}

type rootContainer struct {
	BaseContainer

	backgroundColor rl.Color
	wallpaper       rl.Texture2D
}

func (r *rootContainer) init() {
	r.children.Store([]Component(nil))
	r.SetVisible(true)
	r.this = r
}

func (r *rootContainer) IsVisible() bool {
	return true
}

func (r *rootContainer) Close() {
	if !r.lockIfAlive() {
		return
	}
	children := r.children.Load().([]Component)
	r.children.Store([]Component(nil))
	r.closed.Store(true)
	r.lock.Unlock()
	for _, c := range children {
		c.Close()
	}
}

func (r *rootContainer) Draw(cc *CanvasContext) {
	if r.wallpaper.Width == 0 {
		rl.ClearBackground(r.backgroundColor)
	} else {
		rl.DrawTexture(r.wallpaper, 0, 0, rl.White)
	}
}

func newDisplay(cfg DisplayConfig) *display {
	d := &display{cfg: cfg, logger: logging.NewLogger("raywin.display")}
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.EnableEventWaiting()
	rl.InitWindow(int32(d.cfg.Width), int32(d.cfg.Height), "")
	rl.SetTargetFPS(int32(d.cfg.FPS))
	d.root.init()
	d.root.SetBounds(rl.RectangleInt32{0, 0, int32(cfg.Width), int32(cfg.Height)})
	return d
}

func (d *display) run() error {
	if !atomic.CompareAndSwapInt32(&d.running, 0, 1) {
		return fmt.Errorf("Run() is already runnning: %w", errors.ErrExist)
	}
	d.logger.Infof("Run() starting with %s", d.cfg)
	defer rl.CloseWindow()
	defer func() {
		atomic.StoreInt32(&d.running, 0)
		d.logger.Infof("Run() finishing")
	}()

	d.cc = newCanvas(d.cfg.Width, d.cfg.Height)
	d.tp = &touchPad{}

	startTime := time.Now()
	for !rl.WindowShouldClose() {
		millis := time.Now().Sub(startTime).Milliseconds()
		d.millis.Store(millis)
		d.formFrame(millis)
	}
	return nil
}

func (d *display) formFrame(millis int64) {
	tps := d.tp.onNewFrame(millis)
	if d.tpsAcceptor == nil || d.tpsAcceptor.(Touchpadable).OnTPState(tps) != OnTPSResultLocked {
		d.tpsAcceptor = nil
		// the root is passive, so skip it and start from its children
		d.walkForTouchPadChildren(&d.root)
	}

	d.walkForFC(&d.root, millis)

	rl.BeginDrawing()
	defer rl.EndDrawing()

	d.walkForDrawComp(&d.root, true)
}

func (d *display) walkForFC(c Component, millis int64) {
	if fl, ok := c.(FrameListener); ok {
		fl.OnNewFrame(millis)
	}
	if cont, ok := c.(Container); ok {
		for _, chld := range cont.Children() {
			d.walkForFC(chld, millis)
		}
	}
}

func (d *display) walkForDrawChildren(root Container) {
	var active Component
	for _, chld := range root.Children() {
		r := chld.Bounds()
		if !d.cc.IsVisible(r) || !chld.IsVisible() {
			continue
		}
		if !d.walkForDrawComp(chld, false) {
			active = chld
		}
	}
	if active != nil {
		d.walkForDrawComp(active, true)
	}
}

// walkForDrawComp draws c if it is visble and its children if c is a Container.
// The force flag indicates weather the active (tpsAcceptor) should be drawn. The flag
// is provided to make sure that the active component will be drawn last (see walkForDrawChildren)
func (d *display) walkForDrawComp(c Component, force bool) bool {
	if c != nil && c == d.tpsAcceptor && !force {
		return false
	}

	prevPR := d.cc.PhysicalRegion()
	//TODO fix the Vector2Int32{} param value below
	d.cc.pushRelativeRegion(Vector2Int32{}, c.Bounds())
	scissors := false
	defer func() {
		d.cc.pop()
		if d.cc.isEmpty() {
			rl.EndScissorMode()
		} else if scissors {
			rl.BeginScissorMode(prevPR.X, prevPR.Y, prevPR.Width, prevPR.Height)
		}
	}()

	curPR := d.cc.PhysicalRegion()
	if !hasArea(curPR) {
		// the box is not visible, repoerts true, like it was drawn
		return true
	}
	if curPR != prevPR {
		scissors = true
		rl.BeginScissorMode(curPR.X, curPR.Y, curPR.Width, curPR.Height)
	}

	c.Draw(d.cc)
	if cont, ok := c.(Container); ok {
		d.walkForDrawChildren(cont)
	}
	return true
}

func (d *display) walkForTouchPadChildren(root Container) OnTPSResult {
	if d.tpsAcceptor != nil {
		return OnTPSResultLocked
	}
	x, y := d.cc.relativePointXY(d.tp.tpState().PosXY())
	children := root.Children()
	// walk in backward order, cause the lastest component is the toppest (higher priority) one
	for i := len(children) - 1; i >= 0; i-- {
		c := children[i]
		rect := c.Bounds()
		if !IsPointInRegionInt32(x, y, rect) || !c.IsVisible() {
			continue
		}
		res := d.walkForTouchPadComp(c)
		if res != OnTPSResultNA {
			return res
		}
	}
	return OnTPSResultNA
}

func (d *display) walkForTouchPadComp(root Component) OnTPSResult {
	// TODO: fix the first param value below
	d.cc.pushRelativeRegion(Vector2Int32{}, root.Bounds())
	defer d.cc.pop()

	if !hasArea(d.cc.PhysicalRegion()) {
		// the box is not visible
		return OnTPSResultNA
	}

	if cont, ok := root.(Container); ok {
		res := d.walkForTouchPadChildren(cont)
		if res != OnTPSResultNA {
			return res
		}
	}

	if tp, ok := root.(Touchpadable); ok {
		res := tp.OnTPState(d.tp.tpState())
		if res == OnTPSResultLocked {
			d.tpsAcceptor = root
		}
		return res
	}

	return OnTPSResultNA
}
