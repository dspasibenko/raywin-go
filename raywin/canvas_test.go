package raywin

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCanvas(t *testing.T) {
	cc := newCanvas(10, 20)
	assert.Equal(t, 1, len(cc.stack))
	assert.Equal(t, ctxStackElem{r: rl.RectangleInt32{0, 0, 10, 20}}, cc.stack[0])
}

func TestPushRelativeRegion(t *testing.T) {
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

	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{0, 0, 1000, 2000})
	cc.pushRelativeRegion(Vector2Int32{}, rl.RectangleInt32{10, 10, 1000, 2000})
	assert.Equal(t, ctxStackElem{p: Vector2Int32{X: 10, Y: 10}, r: rl.RectangleInt32{0, 0, 100, 100}}, cc.stack[1])
	cc.pop()
}
