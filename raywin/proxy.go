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
	// RlProxy interface is used by display to making calls to the raylib
	RlProxy interface {
		Init(cfg DisplayConfig)
		CloseWindow()
		WindowShouldClose() bool
		BeginDrawing()
		EndDrawing()
		BeginScissorMode(r rl.RectangleInt32)
		EndScissorMode()
		ClearBackground(color rl.Color)
		DrawTexture(texture rl.Texture2D, pos Vector2Int32, color rl.Color)
		IsMouseButtonDown(mb rl.MouseButton) bool
		GetMouseDelta() rl.Vector2
		GetMousePosition() rl.Vector2
		LoadTextureFromImage(image *rl.Image) rl.Texture2D
		LoadFontEx(fileName string, size int32) rl.Font
		SetTextureFilter(texture rl.Texture2D, filterMode rl.TextureFilterMode)
	}

	realProxy struct{}
	testProxy struct {
		shouldWindowCLose atomic.Bool
		mousePos          rl.Vector2
		mouseDiff         rl.Vector2
	}
)

func (rp *realProxy) Init(cfg DisplayConfig) {
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.EnableEventWaiting()
	rl.InitWindow(int32(cfg.Width), int32(cfg.Height), "")
	rl.SetTargetFPS(int32(cfg.FPS))
}

func (rp *realProxy) CloseWindow() {
	rl.CloseWindow()
}

func (rp *realProxy) WindowShouldClose() bool {
	return rl.WindowShouldClose()
}

func (rp *realProxy) BeginDrawing() {
	rl.BeginDrawing()
}

func (rp *realProxy) EndDrawing() {
	rl.EndDrawing()
}

func (rp *realProxy) BeginScissorMode(r rl.RectangleInt32) {
	rl.BeginScissorMode(r.X, r.Y, r.Width, r.Height)
}

func (rp *realProxy) EndScissorMode() {
	rl.EndScissorMode()
}

func (rp *realProxy) DrawTexture(texture rl.Texture2D, pos Vector2Int32, color rl.Color) {
	rl.DrawTexture(texture, pos.X, pos.Y, color)
}

func (rp *realProxy) ClearBackground(color rl.Color) {
	rl.ClearBackground(color)
}

func (rp *realProxy) IsMouseButtonDown(mb rl.MouseButton) bool {
	return rl.IsMouseButtonDown(mb)
}

func (rp *realProxy) GetMouseDelta() rl.Vector2 {
	return rl.GetMouseDelta()
}

func (rp *realProxy) GetMousePosition() rl.Vector2 {
	return rl.GetMousePosition()
}

func (rp *realProxy) LoadTextureFromImage(image *rl.Image) rl.Texture2D {
	return rl.LoadTextureFromImage(image)
}

func (rp *realProxy) LoadFontEx(fileName string, fontSize int32) rl.Font {
	return rl.LoadFontEx(fileName, fontSize, nil)
}

func (rp *realProxy) SetTextureFilter(texture rl.Texture2D, filterMode rl.TextureFilterMode) {
	rl.SetTextureFilter(texture, filterMode)
}

// ================== testProxy ======================

func (rp *testProxy) Init(cfg DisplayConfig) {
}

func (rp *testProxy) CloseWindow() {
	rp.shouldWindowCLose.Store(true)
}

func (rp *testProxy) WindowShouldClose() bool {
	return rp.shouldWindowCLose.Load()
}

func (rp *testProxy) BeginDrawing() {
}

func (rp *testProxy) EndDrawing() {
}

func (rp *testProxy) BeginScissorMode(r rl.RectangleInt32) {
}

func (rp *testProxy) EndScissorMode() {
}

func (rp *testProxy) DrawTexture(texture rl.Texture2D, pos Vector2Int32, color rl.Color) {
}

func (rp *testProxy) ClearBackground(color rl.Color) {
}

func (rp *testProxy) IsMouseButtonDown(mb rl.MouseButton) bool {
	return !IsEmpty(rp.mousePos)
}

func (rp *testProxy) GetMouseDelta() rl.Vector2 {
	return rp.mouseDiff
}

func (rp *testProxy) GetMousePosition() rl.Vector2 {
	return rp.mousePos
}

func (rp *testProxy) LoadFontEx(fileName string, fontSize int32) rl.Font {
	return rl.Font{BaseSize: fontSize, CharsCount: int32(len(fileName))}
}

func (rp *testProxy) SetTextureFilter(texture rl.Texture2D, filterMode rl.TextureFilterMode) {}

func (rp *testProxy) LoadTextureFromImage(image *rl.Image) rl.Texture2D {
	res := rl.Texture2D{}
	if image != nil {
		// This is fake setting for the testing purposes only
		res.ID = 1
	}
	return res
}
