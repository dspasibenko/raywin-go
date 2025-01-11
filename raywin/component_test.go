package raywin

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBaseComponent_AssertInitialized(t *testing.T) {
	var bc BaseComponent
	assert.NotNil(t, bc.AssertInitialized())

	var owner rootContainer
	owner.init()
	assert.Nil(t, bc.Init(&owner, &bc))

	assert.Nil(t, bc.AssertInitialized())
	bc.Close()
	assert.NotNil(t, bc.AssertInitialized())
}

func TestBaseComponent_Init(t *testing.T) {
	var bc1, bc2 BaseComponent
	var owner rootContainer

	assert.NotNil(t, owner.Init(&owner, &owner)) // owner cannot be the owner by itself

	owner.init()

	assert.NotNil(t, bc1.Init(&owner, &bc2)) // bc2 is not embedded to bc1 !
	assert.Nil(t, bc1.Init(&owner, &bc1))
	assert.NotNil(t, bc1.Init(&owner, &bc1)) // bc1 already has owner

	owner.Close()
	assert.NotNil(t, bc2.Init(&owner, &bc2)) // AddChild is failed here
}

func TestBaseComponent_Bounds(t *testing.T) {
	var bc BaseComponent
	assert.Equal(t, rl.RectangleInt32{}, bc.Bounds())
	r := rl.RectangleInt32{1, 2, 3, 4}
	bc.SetBounds(r)
	assert.Equal(t, r, bc.Bounds())
}

func TestBaseComponent_IsVisible(t *testing.T) {
	var bc BaseComponent
	assert.False(t, bc.IsVisible())
	bc.SetVisible(true)
	assert.True(t, bc.IsVisible())
}

func TestBaseComponent_Close(t *testing.T) {
	var bc BaseComponent
	var owner rootContainer
	owner.init()
	assert.Nil(t, bc.Init(&owner, &bc))
	assert.True(t, len(owner.Children()) == 1)
	assert.Equal(t, &bc, owner.Children()[0])
	assert.Nil(t, bc.AssertInitialized())
	bc.Close()
	assert.NotNil(t, bc.AssertInitialized())
	assert.True(t, len(owner.Children()) == 0)
	bc.Close()
}

func TestBaseComponent_lockIfAlive(t *testing.T) {
	var bc BaseComponent
	assert.True(t, bc.lockIfAlive())
	go func() {
		time.Sleep(time.Millisecond * 10)
		bc.lock.Unlock()
	}()
	assert.True(t, bc.lockIfAlive())
	go func() {
		time.Sleep(time.Millisecond * 10)
		bc.close()
		bc.lock.Unlock()
	}()
	assert.False(t, bc.lockIfAlive())
}

func TestBaseContainer_InitClose(t *testing.T) {
	var bc BaseContainer
	assert.Panics(t, func() {
		bc.Init(nil, &bc)
	})
	var owner rootContainer
	owner.init()

	assert.Nil(t, bc.Init(&owner, &bc))
	assert.True(t, len(owner.Children()) == 1)
	assert.Equal(t, &bc, owner.Children()[0])
	assert.Nil(t, bc.AssertInitialized())
	bc.Close()
	assert.NotNil(t, bc.AssertInitialized())
	assert.True(t, len(owner.Children()) == 0)
	bc.Close()
}

func TestBaseContainer_AddChild(t *testing.T) {
	var bc1, bc2 BaseComponent
	var owner rootContainer
	owner.init()

	assert.Nil(t, bc1.Init(&owner, &bc1))
	assert.Nil(t, bc2.Init(&owner, &bc2))

	assert.Equal(t, []Component{&bc1, &bc2}, owner.Children())
	assert.Nil(t, owner.AddChild(&bc1))
	assert.Equal(t, []Component{&bc2, &bc1}, owner.Children())

	var bc3 BaseComponent
	var owner2 rootContainer
	owner2.init()
	assert.Nil(t, bc3.Init(&owner2, &bc3))

	assert.Nil(t, owner2.AddChild(&bc3))
	assert.NotNil(t, owner2.AddChild(&bc1))
}

func TestBaseContainer_RemoveChild(t *testing.T) {
	var bc1, bc2 BaseComponent
	var owner rootContainer
	owner.init()

	assert.Nil(t, bc1.Init(&owner, &bc1))
	assert.Nil(t, bc2.Init(&owner, &bc2))

	assert.Equal(t, []Component{&bc1, &bc2}, owner.Children())
	assert.True(t, owner.RemoveChild(&bc1))
	assert.False(t, owner.RemoveChild(&bc1))
	assert.Nil(t, owner.AddChild(&bc1))
	assert.True(t, owner.RemoveChild(&bc1))
}
