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

import rl "github.com/gen2brain/raylib-go/raylib"

type (
	Number interface {
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
			~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
			~float32 | ~float64
	}

	Vector2Int32 struct {
		X int32
		Y int32
	}
)

func IsEmpty[T rl.Vector2 | Number](v T) bool {
	return v == *new(T)
}

func hasArea(r rl.RectangleInt32) bool {
	return r.Width != 0 && r.Height != 0
}

func IsPointInRegionInt32(x, y int32, r rl.RectangleInt32) bool {
	return x < r.X+r.Width && x >= r.X && y < r.Y+r.Height && y >= r.Y
}

func VectorDiff(v1, v2 rl.Vector2) rl.Vector2 {
	return rl.Vector2{X: v1.X - v2.X, Y: v1.Y - v2.Y}
}

func (v Vector2Int32) ToVector2() rl.Vector2 {
	return rl.Vector2{X: float32(v.X), Y: float32(v.Y)}
}

func ToVector2Int32(v rl.Vector2) Vector2Int32 {
	return Vector2Int32{X: int32(v.X), Y: int32(v.Y)}
}
