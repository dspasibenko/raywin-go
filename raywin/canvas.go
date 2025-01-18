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

import rl "github.com/gen2brain/raylib-go/raylib"

// CanvasContext struct allows to track the stack of regions, so that one includes another.
// The struct allows to find the physical basis (the position of (0,0) of
// a Component pixel in the physical display coordinates)
type CanvasContext struct {
	stack []ctxStackElem
}

type ctxStackElem struct {
	p Vector2Int32
	r rl.RectangleInt32
}

// PhysicalPointXY returns the coordinates for a Component's point (x,y) on the
// physical display.
func (cc *CanvasContext) PhysicalPointXY(x, y int32) (int32, int32) {
	cse := cc.stack[len(cc.stack)-1]
	return x - cse.p.X + cse.r.X, y - cse.p.Y + cse.r.Y
}

// IsVisible returns whether the region r is visible on the stack of regions
func (cc *CanvasContext) IsVisible(r rl.RectangleInt32) bool {
	px, py := cc.PhysicalPointXY(r.X, r.Y)
	pr := cc.PhysicalRegion()
	return !(px+r.Width < pr.X || pr.X+pr.Width < px ||
		py+r.Height < pr.Y || pr.Y+pr.Height < py)
}

// PhysicalRegion returns the region for the physical screen
func (cc *CanvasContext) PhysicalRegion() rl.RectangleInt32 {
	return cc.stack[len(cc.stack)-1].r
}

// newCanvas constructs the new instance of CanvasContext with the physical dimensions
func newCanvas(width, height uint32) *CanvasContext {
	cc := &CanvasContext{}
	disp := ctxStackElem{r: rl.RectangleInt32{X: 0, Y: 0, Width: int32(width), Height: int32(height)}}
	cc.stack = append(cc.stack, disp) // the cc.stack[0] is always the display resolution
	return cc
}

// pushRelativeRegion adds the physical region (the display coordinates) for the region r, which
// is defined in its parent coordinates stored on top of the stack. vp defines the virtual offset
// in the r. After the call the top of the stack will contain the physical region for r and its
// virtual point for calculation of the region r children, if any...
func (cc *CanvasContext) pushRelativeRegion(vp Vector2Int32, r rl.RectangleInt32) {
	cse := cc.stack[len(cc.stack)-1]
	r.X += cse.r.X // physical X
	r.Y += cse.r.Y // physical Y
	r.X -= cse.p.X // make a correction to the virtual offset for X
	r.Y -= cse.p.Y // and Y
	if r.X < cse.r.X {
		r.Width = max(0, r.Width-(cse.r.X-r.X))
		vp.X += cse.r.X - r.X
		r.X = cse.r.X
	}
	if r.Y < cse.r.Y {
		r.Height = max(0, r.Height-(cse.r.Y-r.Y))
		vp.Y += cse.r.Y - r.Y
		r.Y = cse.r.Y
	}
	r.Width = max(0, min(r.Width, cse.r.Width-(r.X-cse.r.X)))
	r.Height = max(0, min(r.Height, cse.r.Height-(r.Y-cse.r.Y)))
	cc.stack = append(cc.stack, ctxStackElem{p: vp, r: r})
}

func (cc *CanvasContext) pop() {
	if len(cc.stack) < 2 {
		panic("pop() for empty stack called")
	}
	cc.stack = cc.stack[:len(cc.stack)-1]
}

// relativePoingXY gets the physical point {px, py} and turns it to the canvas relative point {x, y}
func (cc *CanvasContext) relativePointXY(px, py int32) (int32, int32) {
	cse := cc.stack[len(cc.stack)-1]
	return px - cse.r.X + cse.p.X, py - cse.r.Y + cse.p.Y
}

func (cc *CanvasContext) isEmpty() bool {
	return len(cc.stack) == 1
}
