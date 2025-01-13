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
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsEmpty(t *testing.T) {
	assert.True(t, IsEmpty[Vector2Int32](Vector2Int32{}))
	assert.False(t, IsEmpty[Vector2Int32](Vector2Int32{X: 1}))
}

func Test_hasArea(t *testing.T) {
	assert.True(t, hasArea(rl.RectangleInt32{0, 0, 1, 2}))
	assert.False(t, hasArea(rl.RectangleInt32{0, 0, 1, 0}))
	assert.False(t, hasArea(rl.RectangleInt32{0, 0, 0, 10}))
}

func TestIsPointInRegionInt32(t *testing.T) {
	assert.True(t, IsPointInRegionInt32(1, 1, rl.RectangleInt32{1, 1, 10, 10}))
	assert.True(t, IsPointInRegionInt32(9, 9, rl.RectangleInt32{1, 1, 10, 10}))
	assert.True(t, IsPointInRegionInt32(1, 9, rl.RectangleInt32{1, 1, 10, 10}))
	assert.True(t, IsPointInRegionInt32(9, 1, rl.RectangleInt32{1, 1, 10, 10}))
	assert.False(t, IsPointInRegionInt32(0, 0, rl.RectangleInt32{1, 1, 10, 10}))
	assert.False(t, IsPointInRegionInt32(10, 10, rl.RectangleInt32{1, 1, 10, 10}))
	assert.False(t, IsPointInRegionInt32(5, 10, rl.RectangleInt32{1, 1, 10, 10}))
	assert.False(t, IsPointInRegionInt32(10, 5, rl.RectangleInt32{1, 1, 10, 10}))
}

func TestVectorDiff(t *testing.T) {
	assert.Equal(t, rl.Vector2{-1, 2}, VectorDiff(rl.Vector2{9, -2}, rl.Vector2{10, -4}))
}

func TestToVector2(t *testing.T) {
	assert.Equal(t, rl.Vector2{-4, 2}, Vector2Int32{-4, 2}.ToVector2())
}

func TestToVector2Int32(t *testing.T) {
	assert.Equal(t, Vector2Int32{33, 4}, ToVector2Int32(rl.Vector2{33, 4}))
}

func TestRectangleInt32ToString(t *testing.T) {
	assert.Equal(t, "{X:1, Y:2, Width:3, Height:4}", RectangleInt32ToString(rl.RectangleInt32{1, 2, 3, 4}))
}
