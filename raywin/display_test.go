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
	"github.com/dspasibenko/raywin-go/pkg/golibs/errors"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type _display_test_container struct {
	BaseContainer
	onnewframes int
	ontpsstate  int
	drawings    int
	onTPSResult OnTPSResult
}

func (dc *_display_test_container) OnNewFrame(_ int64) {
	dc.onnewframes++
}

func (dc *_display_test_container) OnTPState(tps TPState) OnTPSResult {
	dc.ontpsstate++
	return dc.onTPSResult
}

func (dc *_display_test_container) Draw(cc *CanvasContext) {
	dc.drawings++
}

func TestRootContainer(t *testing.T) {
	r := rootContainer{}
	r.init()
	assert.True(t, r.IsVisible())
}

func TestRootContainer_Close(t *testing.T) {
	r := rootContainer{}
	r.init()
	assert.False(t, r.isClosed())
	r.Close()
	assert.True(t, r.isClosed())
	r.Close()
	assert.True(t, r.isClosed())
}

func Test_display_run(t *testing.T) {
	tp := &testProxy{}
	d := newDisplay(DefaultDisplayConfig(), tp)
	var err error
	go func() {
		time.Sleep(10 * time.Millisecond)
		err = d.run(context.Background())
		tp.closeWindow()
	}()
	d.run(context.Background())
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, errors.ErrExist))
}

func Test_display_walkForFC(t *testing.T) {
	d := newDisplay(DefaultDisplayConfig(), &testProxy{}) // no init here is needed

	var c1, c2 _display_test_container
	var owner rootContainer
	owner.init()

	assert.Nil(t, c1.Init(&owner, &c1))
	assert.Nil(t, c2.Init(&c1, &c2))
	d.walkForFC(&c1, 0)

	assert.Equal(t, 1, c1.onnewframes)
	assert.Equal(t, 1, c2.onnewframes)

	d.walkForFC(&c2, 0)
	assert.Equal(t, 1, c1.onnewframes)
	assert.Equal(t, 2, c2.onnewframes)

	d.walkForFC(&owner, 0)
	assert.Equal(t, 2, c1.onnewframes)
	assert.Equal(t, 3, c2.onnewframes)
}

func Test_display_walkForTouchPadComp(t *testing.T) {
	var c _display_test_container
	d := newDisplay(DefaultDisplayConfig(), &testProxy{})

	assert.Nil(t, c.Init(&d.root, &c))
	c.SetBounds(rl.RectangleInt32{X: 0, Y: 0, Width: 0, Height: 10}) // not visible
	assert.Equal(t, OnTPSResultNA, d.walkForTouchPadComp(&c))
	assert.Equal(t, 0, c.ontpsstate)

	c.onTPSResult = OnTPSResultLocked
	c.SetBounds(rl.RectangleInt32{X: 0, Y: 0, Width: 10, Height: 10}) // not visible
	assert.Equal(t, OnTPSResultLocked, d.walkForTouchPadComp(&c))
	assert.Equal(t, 1, c.ontpsstate)
	assert.Equal(t, &c, d.tpsAcceptor)
}

func Test_display_walkForTouchPadChildren(t *testing.T) {
	var c _display_test_container
	d := newDisplay(DefaultDisplayConfig(), &testProxy{})

	assert.Nil(t, c.Init(&d.root, &c))
	c.SetBounds(rl.RectangleInt32{X: 0, Y: 0, Width: 10, Height: 10})
	d.tp.pos = rl.Vector2{X: 4, Y: 5} // within rectangle above
	d.tpsAcceptor = &c
	assert.Equal(t, OnTPSResultLocked, d.walkForTouchPadChildren(&d.root))
	assert.Equal(t, 0, c.ontpsstate)

	d.tpsAcceptor = nil
	c.onTPSResult = OnTPSResultNA
	assert.Equal(t, OnTPSResultNA, d.walkForTouchPadChildren(&d.root))
	assert.Equal(t, 1, c.ontpsstate)

	c.onTPSResult = OnTPSResultStop
	assert.Equal(t, OnTPSResultStop, d.walkForTouchPadChildren(&d.root))
	assert.Equal(t, 2, c.ontpsstate)

	c.SetVisible(false)
	assert.Equal(t, OnTPSResultNA, d.walkForTouchPadChildren(&d.root))
	assert.Equal(t, 2, c.ontpsstate)

	c.SetVisible(true)
	d.tp.pos = rl.Vector2{X: 40, Y: 5} // don't hit the child
	assert.Equal(t, OnTPSResultNA, d.walkForTouchPadChildren(&d.root))
	assert.Equal(t, 2, c.ontpsstate)

	d.tp.pos = rl.Vector2{X: 4, Y: 5} // don't hit the child
	assert.Equal(t, OnTPSResultStop, d.walkForTouchPadChildren(&d.root))
	assert.Equal(t, 3, c.ontpsstate)
}

func Test_display_walkForDrawComp(t *testing.T) {
	var c _display_test_container
	d := newDisplay(DefaultDisplayConfig(), &testProxy{})

	assert.Nil(t, c.Init(&d.root, &c))
	d.tpsAcceptor = &c
	assert.False(t, d.walkForDrawComp(&c, false)) // not drawing active
	assert.Equal(t, 0, c.drawings)

	assert.True(t, d.walkForDrawComp(&c, true))
	assert.Equal(t, 0, c.drawings) // no bounds

	c.SetBounds(rl.RectangleInt32{X: 0, Y: 0, Width: 10, Height: 10})
	assert.True(t, d.walkForDrawComp(&c, true))
	assert.Equal(t, 1, c.drawings) // no bounds

	d.tpsAcceptor = nil
	var c2 _display_test_container
	assert.Nil(t, c2.Init(&c, &c2))
	c2.SetBounds(rl.RectangleInt32{X: 1, Y: 1, Width: 20, Height: 20})
	d.walkForDrawChildren(&d.root)
	assert.Equal(t, 2, c.drawings)
	assert.Equal(t, 1, c2.drawings)

	c2.SetVisible(false) // c2 is not visible
	assert.True(t, d.walkForDrawComp(&c, false))
	assert.Equal(t, 3, c.drawings)
	assert.Equal(t, 1, c2.drawings)
	c2.SetVisible(true)

	c2.SetBounds(rl.RectangleInt32{X: 100, Y: 1, Width: 20, Height: 20}) // c2 is not in the visible area
	assert.True(t, d.walkForDrawComp(&c, false))
	assert.Equal(t, 4, c.drawings)
	assert.Equal(t, 1, c2.drawings)
}
