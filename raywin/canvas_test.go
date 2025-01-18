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
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCanvasContext_NewCanvas(t *testing.T) {
	cc := newCanvas(10, 20)
	assert.Equal(t, 1, len(cc.stack))
	assert.Equal(t, ctxStackElem{r: rl.RectangleInt32{X: 0, Y: 0, Width: 10, Height: 20}}, cc.stack[0])
}

func TestCanvasContext_PushRelativeRegion(t *testing.T) {
	cc := newCanvas(100, 100)
	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{X: 10, Y: 10, Width: 1000, Height: 2000})
	assert.Equal(t, 2, len(cc.stack))
	assert.Equal(t, ctxStackElem{r: rl.RectangleInt32{X: 10, Y: 10, Width: 90, Height: 90}}, cc.stack[1])
	cc.pop()
	assert.Equal(t, 1, len(cc.stack))

	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{X: -10, Y: 10, Width: 1000, Height: 2000})
	assert.Equal(t, ctxStackElem{p: Vector2Int32{X: 10}, r: rl.RectangleInt32{X: 0, Y: 10, Width: 100, Height: 90}}, cc.stack[1])
	cc.pop()

	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{X: -10, Y: -10, Width: 1000, Height: 2000})
	assert.Equal(t, ctxStackElem{p: Vector2Int32{X: 10, Y: 10}, r: rl.RectangleInt32{X: 0, Y: 0, Width: 100, Height: 100}}, cc.stack[1])
	cc.pop()

	cc.pushRelativeRegion(Vector2Int32{X: 10, Y: 5}, rl.RectangleInt32{X: 0, Y: 0, Width: 1000, Height: 2000})
	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{X: 10, Y: 10, Width: 1000, Height: 2000})
	assert.Equal(t, ctxStackElem{p: Vector2Int32{X: 0, Y: 0}, r: rl.RectangleInt32{X: 0, Y: 5, Width: 100, Height: 95}}, cc.stack[2])
	cc.pop()
	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{X: 0, Y: 0, Width: 50, Height: 50})
	assert.Equal(t, ctxStackElem{p: Vector2Int32{X: 10, Y: 5}, r: rl.RectangleInt32{X: 0, Y: 0, Width: 40, Height: 45}}, cc.stack[2])
	cc.pop()
	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{X: 0, Y: 0, Width: 5, Height: 5})
	assert.Equal(t, ctxStackElem{p: Vector2Int32{X: 10, Y: 5}, r: rl.RectangleInt32{}}, cc.stack[2])
	cc.pop()
	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{X: 109, Y: 104, Width: 50, Height: 50})
	assert.Equal(t, ctxStackElem{p: Vector2Int32{X: 0, Y: 0}, r: rl.RectangleInt32{X: 99, Y: 99, Width: 1, Height: 1}}, cc.stack[2])
	cc.pop()
	cc.pop()
}

func TestCanvasContext_Pop(t *testing.T) {
	cc := newCanvas(100, 100)
	assert.True(t, cc.isEmpty())
	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{X: 109, Y: 104, Width: 50, Height: 50})
	assert.False(t, cc.isEmpty())
	cc.pop()
	assert.True(t, cc.isEmpty())
	assert.Panics(t, func() {
		cc.pop()
	})
	assert.True(t, cc.isEmpty())
}

func TestCanvasContext_RelativePointXY(t *testing.T) {
	cc := newCanvas(100, 100)
	x, y := cc.relativePointXY(10, 10)
	assert.True(t, x == 10 && y == 10)

	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{X: 10, Y: 10, Width: 50, Height: 50})
	x, y = cc.relativePointXY(10, 10)
	assert.True(t, x == 0 && y == 0)
	x, y = cc.relativePointXY(0, 0)
	assert.True(t, x == -10 && y == -10)
	cc.pop()

	cc.pushRelativeRegion(Vector2Int32{X: 10, Y: 10}, rl.RectangleInt32{X: 10, Y: 10, Width: 50, Height: 50})
	x, y = cc.relativePointXY(10, 10)
	assert.True(t, x == 10 && y == 10)
	cc.pop()
}

func TestCanvasContext_PhysicalPointXY(t *testing.T) {
	cc := newCanvas(100, 100)
	x, y := cc.PhysicalPointXY(10, 10)
	assert.True(t, x == 10 && y == 10)

	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{X: 10, Y: 10, Width: 50, Height: 50})
	x, y = cc.PhysicalPointXY(10, 10)
	assert.True(t, x == 20 && y == 20)
	x, y = cc.PhysicalPointXY(-5, -5)
	assert.True(t, x == 5 && y == 5)
	cc.pop()

	cc.pushRelativeRegion(Vector2Int32{5, 5}, rl.RectangleInt32{X: 10, Y: 10, Width: 50, Height: 50})
	x, y = cc.PhysicalPointXY(10, 10)
	assert.True(t, x == 15 && y == 15)
	cc.pop()
}

func TestCanvasContext_IsVisible(t *testing.T) {
	cc := newCanvas(100, 100)
	assert.True(t, cc.IsVisible(rl.RectangleInt32{X: 10, Y: 10, Width: 20, Height: 20}))
	assert.True(t, cc.IsVisible(rl.RectangleInt32{X: -10, Y: -110, Width: 200, Height: 200}))
	assert.True(t, cc.IsVisible(rl.RectangleInt32{X: 10, Y: 10, Width: 20, Height: 200}))
	assert.True(t, cc.IsVisible(rl.RectangleInt32{X: 10, Y: 10, Width: 200, Height: 20}))
	assert.True(t, cc.IsVisible(rl.RectangleInt32{X: -10, Y: 10, Width: 20, Height: 20}))
	assert.True(t, cc.IsVisible(rl.RectangleInt32{X: 10, Y: -10, Width: 20, Height: 20}))

	assert.False(t, cc.IsVisible(rl.RectangleInt32{X: -10, Y: -10, Width: 5, Height: 5}))
	assert.False(t, cc.IsVisible(rl.RectangleInt32{X: 110, Y: 110, Width: 5, Height: 5}))

	cc.pushRelativeRegion(Vector2Int32{20, 20}, rl.RectangleInt32{X: 0, Y: 0, Width: 50, Height: 50})
	assert.False(t, cc.IsVisible(rl.RectangleInt32{X: 10, Y: 10, Width: 5, Height: 5}))
	assert.True(t, cc.IsVisible(rl.RectangleInt32{X: 60, Y: 60, Width: 5, Height: 5}))
}

func TestCanvasContext_PhysicalRegion(t *testing.T) {
	cc := newCanvas(100, 100)
	assert.Equal(t, rl.RectangleInt32{X: 0, Y: 0, Width: 100, Height: 100}, cc.PhysicalRegion())
	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{X: 10, Y: 10, Width: 50, Height: 50})
	assert.Equal(t, rl.RectangleInt32{X: 10, Y: 10, Width: 50, Height: 50}, cc.PhysicalRegion())
	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{X: 10, Y: 10, Width: 50, Height: 50})
	assert.Equal(t, rl.RectangleInt32{X: 20, Y: 20, Width: 40, Height: 40}, cc.PhysicalRegion())
	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{X: 10, Y: 10, Width: 5, Height: 5})
	assert.Equal(t, rl.RectangleInt32{X: 30, Y: 30, Width: 5, Height: 5}, cc.PhysicalRegion())
	cc.pop()
	cc.pop()
	cc.pop()

	cc.pushRelativeRegion(Vector2Int32{X: 10, Y: 10}, rl.RectangleInt32{X: 10, Y: 10, Width: 50, Height: 50})
	assert.Equal(t, rl.RectangleInt32{X: 10, Y: 10, Width: 50, Height: 50}, cc.PhysicalRegion())
	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{X: 10, Y: 10, Width: 40, Height: 40})
	assert.Equal(t, rl.RectangleInt32{X: 10, Y: 10, Width: 40, Height: 40}, cc.PhysicalRegion())
	cc.pop()
	cc.pop()

	cc.pushRelativeRegion(Vector2Int32{X: -10, Y: -10}, rl.RectangleInt32{X: 10, Y: 10, Width: 50, Height: 50})
	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{X: -5, Y: -5, Width: 50, Height: 50})
	assert.Equal(t, rl.RectangleInt32{X: 15, Y: 15, Width: 45, Height: 45}, cc.PhysicalRegion())
}
