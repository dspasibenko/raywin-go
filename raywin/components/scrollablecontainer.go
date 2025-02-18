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
)

// ScrollableContainer struct offers the BascContainer functionality with scroll
// bars on it. It inherits raywin.BaseContainer and embeds raywin.InertialScroller to
// provide the scrolling functionality.
//
// # The component using the current Style to draw the scrollbars
//
// The ScrollableContainer implements PostDrawer (see DrawAfter) to draw the scrollbars
// on top of the ScrollableContainer area and its children. The Draw() function is not
// implemented, so if the ScrollableContainer is embedded the embedded component may draw
// a background for the container. Please take a look at `scrollablecontainer` example.
type ScrollableContainer struct {
	raywin.BaseContainer
	raywin.InertialScroller

	showFlags     int
	releaseMillis int64
}

const (
	// ShowHorizontalScrollBar enables showing the horizontal scroll bar
	ShowHorizontalScrollBar = 0b100
	// ShowVerticalScrollBar enables showint the vertical scroll bar
	ShowVerticalScrollBar = 0b1000
	// ShowBothScrollBar enables to show vertical and horizontal bars both
	ShowBothScrollBar = 0b1100
	// ScrollBarLightColor use the light scroll bar color
	ScrollBarLightColor = 0b10000
	// ScrollableContainerAutoVirtualSize automatically adjusting the scrolling area to the
	// container childs. If the flag is provided the virtual bounds are automatically calculated
	// based on the children dimenstions
	ScrollableContainerAutoVirtualSize = 0b100000
)

// InitScrollableContainer initializes the ScrollableContainer:
// `owner` - the container, which owns the ScrollableContainer
// `this` - the instance which embeds (if any) of the ScrollableContainer or the ScrollableContainer instance itself
// `flags` - the scrollbar settings. The flags are mix of the flags above and the InertialScroller flags
//
// The ScrollableContainer exposes the InitScrollableContainer cause it can be embedded into another component
func (sc *ScrollableContainer) InitScrollableContainer(owner, this raywin.Container, flags int) error {
	sc.showFlags = flags
	o := owner.(raywin.Component)
	c := this.(raywin.Component)
	if err := sc.InitInertialScroller(c, o.Bounds(), raywin.DefaultInternalScrollerDeceleration(), uint8(flags)&raywin.ScrollBoth); err != nil {
		return err
	}
	return sc.Init(owner, c)
}

// OnNewFrame provides the FrameListener implementation
func (sc *ScrollableContainer) OnNewFrame(millis int64) {
	if sc.showFlags&ScrollableContainerAutoVirtualSize != 0 {
		sc.autoResize()
	}
	if sc.IsTPLocked() {
		sc.releaseMillis = -1
	} else if sc.releaseMillis == -1 {
		sc.releaseMillis = millis
	}
	sc.InertialScroller.OnNewFrame(millis)
}

// DrawAfter provides the PostDrawer implementation
func (sc *ScrollableContainer) DrawAfter(cc *raywin.CanvasContext) {
	if !sc.shouldDraw() {
		return
	}
	bi := sc.Bounds()
	b0 := bi.ToFloat32()
	vbi := sc.VirtualBounds()
	vb := vbi.ToFloat32()

	showHorizontal := b0.Width < vb.Width && (sc.showFlags&ShowHorizontalScrollBar != 0)
	showVertical := b0.Height < vb.Height && (sc.showFlags&ShowVerticalScrollBar != 0)
	col := S.ScrollBarDarkColor
	if sc.showFlags&ScrollBarLightColor != 0 {
		col = S.ScrollBarLightColor
	}
	w := S.ScrollBarThiknessMm * S.PPcm / 10.0
	space := S.ScrollBarOffsetMm * S.PPcm / 10.0
	if showHorizontal {
		b := b0
		if showVertical {
			b.Width -= space + w
		}
		ln := b.Width * b.Width / vb.Width
		if vb.X < 0 {
			ln = b.Width * b.Width / (vb.Width - vb.X)
		}
		if vb.X > vb.Width-b.Width {
			ln = b.Width * b.Width / (vb.X + b.Width)
		}
		ln = min(b.Width, max(ln, w))

		offs := float32(0.0)
		if vb.X > 0 {
			offs = (b.Width - ln) * min(1.0, vb.X/(vb.Width-b.Width))
		}

		px, py := cc.PhysicalPointXY(vbi.X, vbi.Y)
		x := float32(px) + offs
		y := float32(py) + b.Height - w - space

		r := w / 2.0
		c := rl.Vector2{X: x + r, Y: y + r}
		rl.DrawCircleSector(c, r, 90, 270, int32(r), col)
		rl.DrawRectangleV(rl.Vector2{X: x + r, Y: y}, rl.Vector2{X: ln - w, Y: w}, col)
		c.X += ln - w
		rl.DrawCircleSector(c, r, 270, 450, int32(r), col)
	}
	if showVertical {
		b := b0
		if showHorizontal {
			b.Height -= space + w
		}
		ln := b.Height * b.Height / vb.Height
		if vb.Y < 0 {
			ln = b.Height * b.Height / (vb.Height - vb.Y)
		}
		if vb.Y > vb.Height-b.Height {
			ln = b.Height * b.Height / (vb.Y + b.Height)
		}
		ln = min(b.Height, max(ln, w))

		offs := float32(0.0)
		if vb.Y > 0 {
			offs = (b.Height - ln) * min(1.0, vb.Y/(vb.Height-b.Height))
		}

		px, py := cc.PhysicalPointXY(vbi.X, vbi.Y)
		x := float32(px) + b.Width - w - space
		y := float32(py) + offs

		r := w / 2.0
		c := rl.Vector2{X: x + r, Y: y + r}
		rl.DrawCircleSector(c, r, 180, 360, int32(r), col)
		rl.DrawRectangleV(rl.Vector2{X: x, Y: y + r}, rl.Vector2{X: w, Y: ln - w}, col)
		c.Y += ln - w
		rl.DrawCircleSector(c, r, 180, 0, int32(r), col)
	}
}

func (sc *ScrollableContainer) shouldDraw() bool {
	if sc.showFlags&ShowBothScrollBar == 0 {
		return false
	}
	if sc.IsTPLocked() {
		return true
	}
	if sc.releaseMillis == 0 {
		return false
	}
	return raywin.Millis()-sc.releaseMillis < int64(S.ScrollBarDisappearMillis)
}

func (sc *ScrollableContainer) autoResize() {
	left := int32(-1)
	top := int32(-1)
	bnds := sc.Bounds()
	width := bnds.Width
	height := bnds.Height
	for _, c := range sc.Children() {
		b := c.Bounds()
		if b.X >= 0 {
			if left < 0 {
				left = b.X
			} else {
				left = min(b.X, left)
			}
		}
		if b.Y >= 0 {
			if top < 0 {
				top = b.Y
			} else {
				top = min(b.Y, top)
			}
		}
		width = max(b.X+b.Width, width)
		height = max(b.Y+b.Height, height)
	}
	width += left
	height += top
	vb := sc.VirtualBounds()
	if width != vb.Width || height != vb.Height {
		vb.Height = height
		vb.Width = width
		sc.SetVirtualBounds(vb)
	}
}
