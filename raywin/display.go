package raywin

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
	"context"
	"fmt"
	"github.com/dspasibenko/raywin-go/pkg/golibs/errors"
	"github.com/dspasibenko/raywin-go/pkg/golibs/logging"
	rl "github.com/gen2brain/raylib-go/raylib"
	"sync/atomic"
	"time"
)

type display struct {
	proxy RlProxy
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
	frmListener FrameListener
}

type rootContainer struct {
	BaseContainer

	proxy           RlProxy
	backgroundColor rl.Color
	wallpaper       rl.Texture2D
}

func (r *rootContainer) init() {
	r.children.Store([]Component(nil))
	r.SetVisible(true)
	r.this = r
}

// IsVisible for rootContainer is always true
func (r *rootContainer) IsVisible() bool {
	return true
}

// Close is overwritten for rootContainer due to no owner notification about the close
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

// Draw for the display - either the background color or a wallpaper picture
func (r *rootContainer) Draw(cc *CanvasContext) {
	if r.wallpaper.Width == 0 {
		r.proxy.ClearBackground(r.backgroundColor)
	} else {
		r.proxy.DrawTexture(r.wallpaper, Vector2Int32{0, 0}, rl.White)
	}
}

func newDisplay(cfg DisplayConfig, rp RlProxy) *display {
	d := &display{cfg: cfg, logger: logging.NewLogger("raywin.display")}
	d.proxy = rp
	d.proxy.Init(cfg)
	d.root.proxy = rp
	d.root.init()
	d.root.SetBounds(rl.RectangleInt32{X: 0, Y: 0, Width: int32(cfg.Width), Height: int32(cfg.Height)})
	d.cc = newCanvas(d.cfg.Width, d.cfg.Height)
	d.tp = &touchPad{}
	return d
}

func (d *display) run(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&d.running, 0, 1) {
		return fmt.Errorf("Run() is already runnning: %w", errors.ErrExist)
	}
	d.logger.Infof("Run() starting with %s", d.cfg)
	defer d.proxy.CloseWindow()
	defer d.root.Close()
	defer func() {
		atomic.StoreInt32(&d.running, 0)
		d.logger.Infof("Run() finishing")
	}()

	startTime := time.Now()
	for !d.proxy.WindowShouldClose() && ctx.Err() == nil {
		millis := time.Now().Sub(startTime).Milliseconds()
		d.millis.Store(millis)
		if d.frmListener != nil {
			d.frmListener.OnNewFrame(millis)
		}
		d.formFrame(millis)
	}
	return ctx.Err()
}

func (d *display) formFrame(millis int64) {
	tps := d.tp.onNewFrame(millis, d.proxy)
	if d.tpsAcceptor == nil || d.tpsAcceptor.baseComponent().isClosed() || d.tpsAcceptor.(Touchpadable).OnTPState(tps) != OnTPSResultLocked {
		d.tpsAcceptor = nil
		// the root is passive, so skip it and start from its children
		d.walkForTouchPadChildren(&d.root)
	}

	d.walkForFC(&d.root, millis)

	d.proxy.BeginDrawing()
	defer d.proxy.EndDrawing()

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
	var offs Vector2Int32
	if s, ok := c.(Scrollable); ok {
		offs = s.Offset()
	}
	d.cc.pushRelativeRegion(offs, c.Bounds())
	scissors := false
	defer func() {
		d.cc.pop()
		if d.cc.isEmpty() {
			d.proxy.EndScissorMode()
		} else if scissors {
			d.proxy.BeginScissorMode(prevPR)
		}
	}()

	curPR := d.cc.PhysicalRegion()
	if !hasArea(curPR) {
		// the box is not visible, repoerts true, like it was drawn
		return true
	}
	if curPR != prevPR {
		scissors = true
		d.proxy.BeginScissorMode(curPR)
	}

	c.Draw(d.cc)
	if cont, ok := c.(Container); ok {
		d.walkForDrawChildren(cont)
	}
	if pd, ok := c.(PostDrawer); ok {
		pd.DrawAfter(d.cc)
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
	var offs Vector2Int32
	if s, ok := root.(Scrollable); ok {
		offs = s.Offset()
	}
	d.cc.pushRelativeRegion(offs, root.Bounds())
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
