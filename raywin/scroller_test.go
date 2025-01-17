package raywin

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultInternalScrollerDeceleration(t *testing.T) {
	assert.Nil(t, c.initConfig(DefaultConfig(), &testProxy{}))
	v := DefaultInternalScrollerDeceleration()
	assert.True(t, v.X < 0 && v.Y < 0 && v.X > -2.0 && v.Y > -2.0)
}

func TestInertialScroller_InitScroller(t *testing.T) {
	c = &controller{}
	defer func() {
		c = &controller{}
	}()
	assert.Nil(t, c.initConfig(DefaultConfig(), &testProxy{}))
	var is InertialScroller
	assert.NotNil(t, is.InitScroller(&c.disp.root, rl.RectangleInt32{}, rl.Vector2{-1, 2}, ScrollBoth))
	assert.NotNil(t, is.InitScroller(&c.disp.root, rl.RectangleInt32{}, rl.Vector2{1, -1}, ScrollBoth))
	assert.NotNil(t, is.InitScroller(nil, rl.RectangleInt32{}, rl.Vector2{-1, -1}, ScrollBoth))
	b := rl.RectangleInt32{1, 2, 3, 4}
	assert.Nil(t, is.InitScroller(&c.disp.root, b, rl.Vector2{-1, -1}, ScrollBoth))
	assert.Equal(t, uint8(ScrollBoth), is.flags)
	assert.Equal(t, rl.Vector2{-1, -1}, is.decel)
	assert.Equal(t, b, is.virtBounds.Load().(rl.RectangleInt32))
	assert.Equal(t, &c.disp.root, is.owner)
}

func TestInertialScroller_OnTPState(t *testing.T) {
	c = &controller{}
	defer func() {
		c = &controller{}
	}()
	assert.Nil(t, c.initConfig(DefaultConfig(), &testProxy{}))
	var is InertialScroller
	assert.Nil(t, is.InitScroller(&c.disp.root, rl.RectangleInt32{0, 0, 200, 200},
		DefaultInternalScrollerDeceleration(), ScrollBoth))
	for i := 0; i < c.disp.cfg.FPS; i++ {
		assert.Equal(t, OnTPSResultLocked, is.OnTPState(TPState{Pos: rl.Vector2{float32(200 - i), float32(200 - i)}, Millis: int64(i), State: TPStateMoving}))
	}
	assert.Equal(t, OnTPSResultLocked, is.OnTPState(TPState{Pos: rl.Vector2{float32(200), float32(200)}, Millis: int64(100), State: TPStateMoving}))
	assert.Equal(t, is.samples.Cap(), is.samples.Len())
	assert.True(t, is.IsTPLocked())
	assert.Equal(t, OnTPSResultNA, is.OnTPState(TPState{State: TPStateReleased, Millis: int64(101), Pos: rl.Vector2{100, 100}}))
	assert.False(t, is.IsTPLocked())
	assert.Equal(t, Vector2Int32{}, is.Offset())
	assert.Equal(t, rl.Vector2{29.5, 29.5}, is.velo)
	assert.Equal(t, rl.Vector2{-1, -1}, is.dir)
}

func TestInertialScroller_getDiffForLastFrame(t *testing.T) {
	var is InertialScroller
	is.diff = rl.Vector2{-1, -1}
	is.flags = ScrollBoth
	assert.Equal(t, rl.Vector2{-1, -1}, is.getDiffForLastFrame())
	is.flags = ScrollHorizontal
	assert.Equal(t, rl.Vector2{-1, 0}, is.getDiffForLastFrame())
	is.flags = ScrollVertical
	assert.Equal(t, rl.Vector2{0, 0}, is.getDiffForLastFrame())
}

func TestInertialScroller_SetVirtualBounds(t *testing.T) {
	var is InertialScroller
	r := rl.RectangleInt32{1, 2, 3, 4}
	is.SetVirtualBounds(r)
	assert.Equal(t, r, is.VirtualBounds())
}

func TestInertialScroller_OnNewFrame(t *testing.T) {
	c = &controller{}
	defer func() {
		c = &controller{}
	}()
	assert.Nil(t, c.initConfig(DefaultConfig(), &testProxy{}))
	var is InertialScroller
	is.InitScroller(RootContainer().(Component), rl.RectangleInt32{0, 0, 200, 200},
		DefaultInternalScrollerDeceleration(), ScrollBoth)
	r := rl.RectangleInt32{0, 0, 100, 100}
	is.SetVirtualBounds(r)
	is.velo = rl.Vector2{10, 10}
	is.decel = rl.Vector2{-1, 1}
	is.dir = rl.Vector2{-1, -1}
	is.OnNewFrame(10)
	assert.Equal(t, Vector2Int32{-5, -6}, is.Offset())
}
