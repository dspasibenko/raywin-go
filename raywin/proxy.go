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
	"sync/atomic"
)

type (
	// rlProxy interface is used by display to making calls to the raylib
	rlProxy interface {
		init(cfg DisplayConfig)
		closeWindow()
		windowShouldClose() bool
		beginDrawing()
		endDrawing()
		beginScissorMode(r rl.RectangleInt32)
		endScissorMode()
		clearBackground(color rl.Color)
		drawTexture(texture rl.Texture2D, pos Vector2Int32, color rl.Color)

		isMouseButtonDown(mb rl.MouseButton) bool
		getMouseDelta() rl.Vector2
		getMousePosition() rl.Vector2

		loadTextureFromImage(image *rl.Image) rl.Texture2D
	}

	realProxy struct{}
	testProxy struct {
		shouldWindowCLose atomic.Bool
		mousePos          rl.Vector2
		mouseDiff         rl.Vector2
	}
)

func (rp *realProxy) init(cfg DisplayConfig) {
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.EnableEventWaiting()
	rl.InitWindow(int32(cfg.Width), int32(cfg.Height), "")
	rl.SetTargetFPS(int32(cfg.FPS))
}

func (rp *realProxy) closeWindow() {
	rl.CloseWindow()
}

func (rp *realProxy) windowShouldClose() bool {
	return rl.WindowShouldClose()
}

func (rp *realProxy) beginDrawing() {
	rl.BeginDrawing()
}

func (rp *realProxy) endDrawing() {
	rl.EndDrawing()
}

func (rp *realProxy) beginScissorMode(r rl.RectangleInt32) {
	rl.BeginScissorMode(r.X, r.Y, r.Width, r.Height)
}

func (rp *realProxy) endScissorMode() {

}

func (rp *realProxy) drawTexture(texture rl.Texture2D, pos Vector2Int32, color rl.Color) {
	rl.DrawTexture(texture, pos.X, pos.Y, color)
}

func (rp *realProxy) clearBackground(color rl.Color) {
	rl.ClearBackground(color)
}

func (rp *realProxy) isMouseButtonDown(mb rl.MouseButton) bool {
	return rl.IsMouseButtonDown(rl.MouseLeftButton)
}

func (rp *realProxy) getMouseDelta() rl.Vector2 {
	return rl.GetMouseDelta()
}

func (rp *realProxy) getMousePosition() rl.Vector2 {
	return rl.GetMousePosition()
}

func (rp *realProxy) loadTextureFromImage(image *rl.Image) rl.Texture2D {
	return rl.LoadTextureFromImage(image)
}

// ================== testProxy ======================

func (rp *testProxy) init(cfg DisplayConfig) {
}

func (rp *testProxy) closeWindow() {
	rp.shouldWindowCLose.Store(true)
}

func (rp *testProxy) windowShouldClose() bool {
	return rp.shouldWindowCLose.Load()
}

func (rp *testProxy) beginDrawing() {
}

func (rp *testProxy) endDrawing() {
}

func (rp *testProxy) beginScissorMode(r rl.RectangleInt32) {
}

func (rp *testProxy) endScissorMode() {
}

func (rp *testProxy) drawTexture(texture rl.Texture2D, pos Vector2Int32, color rl.Color) {
}

func (rp *testProxy) clearBackground(color rl.Color) {
}

func (rp *testProxy) isMouseButtonDown(mb rl.MouseButton) bool {
	return !IsEmpty(rp.mousePos)
}

func (rp *testProxy) getMouseDelta() rl.Vector2 {
	return rp.mouseDiff
}

func (rp *testProxy) getMousePosition() rl.Vector2 {
	return rp.mousePos
}

func (rp *testProxy) loadTextureFromImage(image *rl.Image) rl.Texture2D {
	res := rl.Texture2D{}
	if image != nil {
		// This is fake setting for the testing purposes only
		res.ID = 1
	}
	return res
}
