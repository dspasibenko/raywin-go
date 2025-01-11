package raywin

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCanvasContext_NewCanvas(t *testing.T) {
	cc := newCanvas(10, 20)
	assert.Equal(t, 1, len(cc.stack))
	assert.Equal(t, ctxStackElem{r: rl.RectangleInt32{0, 0, 10, 20}}, cc.stack[0])
}

func TestCanvasContext_PushRelativeRegion(t *testing.T) {
	cc := newCanvas(100, 100)
	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{10, 10, 1000, 2000})
	assert.Equal(t, 2, len(cc.stack))
	assert.Equal(t, ctxStackElem{r: rl.RectangleInt32{10, 10, 90, 90}}, cc.stack[1])
	cc.pop()
	assert.Equal(t, 1, len(cc.stack))

	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{-10, 10, 1000, 2000})
	assert.Equal(t, ctxStackElem{p: Vector2Int32{X: 10}, r: rl.RectangleInt32{0, 10, 100, 90}}, cc.stack[1])
	cc.pop()

	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{-10, -10, 1000, 2000})
	assert.Equal(t, ctxStackElem{p: Vector2Int32{X: 10, Y: 10}, r: rl.RectangleInt32{0, 0, 100, 100}}, cc.stack[1])
	cc.pop()

	cc.pushRelativeRegion(Vector2Int32{10, 5}, rl.RectangleInt32{0, 0, 1000, 2000})
	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{10, 10, 1000, 2000})
	assert.Equal(t, ctxStackElem{p: Vector2Int32{X: 0, Y: 0}, r: rl.RectangleInt32{0, 5, 100, 95}}, cc.stack[2])
	cc.pop()
	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{0, 0, 50, 50})
	assert.Equal(t, ctxStackElem{p: Vector2Int32{X: 10, Y: 5}, r: rl.RectangleInt32{0, 0, 40, 45}}, cc.stack[2])
	cc.pop()
	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{0, 0, 5, 5})
	assert.Equal(t, ctxStackElem{p: Vector2Int32{X: 10, Y: 5}, r: rl.RectangleInt32{0, 0, 0, 0}}, cc.stack[2])
	cc.pop()
	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{109, 104, 50, 50})
	assert.Equal(t, ctxStackElem{p: Vector2Int32{X: 0, Y: 0}, r: rl.RectangleInt32{99, 99, 1, 1}}, cc.stack[2])
	cc.pop()
	cc.pop()
}

func TestCanvasContext_Pop(t *testing.T) {
	cc := newCanvas(100, 100)
	assert.True(t, cc.isEmpty())
	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{109, 104, 50, 50})
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

	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{10, 10, 50, 50})
	x, y = cc.relativePointXY(10, 10)
	assert.True(t, x == 0 && y == 0)
	x, y = cc.relativePointXY(0, 0)
	assert.True(t, x == -10 && y == -10)
	cc.pop()

	cc.pushRelativeRegion(Vector2Int32{10, 10}, rl.RectangleInt32{10, 10, 50, 50})
	x, y = cc.relativePointXY(10, 10)
	assert.True(t, x == 10 && y == 10)
	cc.pop()
}

func TestCanvasContext_PhysicalPointXY(t *testing.T) {
	cc := newCanvas(100, 100)
	x, y := cc.PhysicalPointXY(10, 10)
	assert.True(t, x == 10 && y == 10)

	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{10, 10, 50, 50})
	x, y = cc.PhysicalPointXY(10, 10)
	assert.True(t, x == 20 && y == 20)
	x, y = cc.PhysicalPointXY(-5, -5)
	assert.True(t, x == 5 && y == 5)
	cc.pop()

	cc.pushRelativeRegion(Vector2Int32{5, 5}, rl.RectangleInt32{10, 10, 50, 50})
	x, y = cc.PhysicalPointXY(10, 10)
	assert.True(t, x == 15 && y == 15)
	cc.pop()
}

func TestCanvasContext_IsVisible(t *testing.T) {
	cc := newCanvas(100, 100)
	assert.True(t, cc.IsVisible(rl.RectangleInt32{10, 10, 20, 20}))
	assert.True(t, cc.IsVisible(rl.RectangleInt32{-10, -110, 200, 200}))
	assert.True(t, cc.IsVisible(rl.RectangleInt32{10, 10, 20, 200}))
	assert.True(t, cc.IsVisible(rl.RectangleInt32{10, 10, 200, 20}))
	assert.True(t, cc.IsVisible(rl.RectangleInt32{-10, 10, 20, 20}))
	assert.True(t, cc.IsVisible(rl.RectangleInt32{10, -10, 20, 20}))

	assert.False(t, cc.IsVisible(rl.RectangleInt32{-10, -10, 5, 5}))
	assert.False(t, cc.IsVisible(rl.RectangleInt32{110, 110, 5, 5}))

	cc.pushRelativeRegion(Vector2Int32{20, 20}, rl.RectangleInt32{0, 0, 50, 50})
	assert.False(t, cc.IsVisible(rl.RectangleInt32{10, 10, 5, 5}))
	assert.True(t, cc.IsVisible(rl.RectangleInt32{60, 60, 5, 5}))
}

func TestCanvasContext_PhysicalRegion(t *testing.T) {
	cc := newCanvas(100, 100)
	assert.Equal(t, rl.RectangleInt32{0, 0, 100, 100}, cc.PhysicalRegion())
	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{10, 10, 50, 50})
	assert.Equal(t, rl.RectangleInt32{10, 10, 50, 50}, cc.PhysicalRegion())
	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{10, 10, 50, 50})
	assert.Equal(t, rl.RectangleInt32{20, 20, 40, 40}, cc.PhysicalRegion())
	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{10, 10, 5, 5})
	assert.Equal(t, rl.RectangleInt32{30, 30, 5, 5}, cc.PhysicalRegion())
	cc.pop()
	cc.pop()
	cc.pop()

	cc.pushRelativeRegion(Vector2Int32{10, 10}, rl.RectangleInt32{10, 10, 50, 50})
	assert.Equal(t, rl.RectangleInt32{10, 10, 50, 50}, cc.PhysicalRegion())
	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{10, 10, 40, 40})
	assert.Equal(t, rl.RectangleInt32{10, 10, 40, 40}, cc.PhysicalRegion())
	cc.pop()
	cc.pop()

	cc.pushRelativeRegion(Vector2Int32{-10, -10}, rl.RectangleInt32{10, 10, 50, 50})
	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{-5, -5, 50, 50})
	assert.Equal(t, rl.RectangleInt32{15, 15, 45, 45}, cc.PhysicalRegion())
}
