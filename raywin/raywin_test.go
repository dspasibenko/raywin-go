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
	"fmt"
	"github.com/dspasibenko/raywin-go/pkg/golibs/errors"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
	"time"
)

func Test_assertInitialized(t *testing.T) {
	assert.Panics(t, func() {
		assertInitialized()
	})
}

func TestInit(t *testing.T) {
	c = &controller{}
	p = &testProxy{}
	assert.Nil(t, Init(DefaultConfig()))
	assert.True(t, c.valid.Load())
	assert.NotNil(t, c.disp)
}

type testFrameListener struct {
	ms int64
}

// OnNewFrame implements FrameListener
func (tfl *testFrameListener) OnNewFrame(millis int64) {
	tfl.ms = millis
}

func TestRunAndMisc(t *testing.T) {
	c = &controller{}
	p = &testProxy{}
	cfg := DefaultConfig()
	fmt.Println(cfg)
	tfl := &testFrameListener{}
	cfg.FrameListener = tfl
	assert.Nil(t, Init(cfg))
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()
	err := Run(ctx)
	assert.Equal(t, ctx.Err(), err)
	assert.Equal(t, tfl.ms, Millis())

	_, err = GetIcon("lala")
	assert.Equal(t, errors.ErrNotExist, err)
}

func Test_controller_initConfig(t *testing.T) {
	cfg := Config{
		DisplayConfig:       DefaultDisplayConfig(),
		WallpaperFileName:   filepath.FromSlash("testdata/images/wallpaper800x.png"),
		RegularFontFileName: filepath.FromSlash("testdata/fonts/Roboto/Roboto-Medium.ttf"),
		ItalicFontFileName:  filepath.FromSlash("testdata/fonts/Roboto/Roboto-MediumItalic.ttf"),
		IconsDir:            filepath.FromSlash("testdata/icons"),
	}
	c = &controller{}
	defer func() {
		c = &controller{}
	}()
	assert.Nil(t, c.initConfig(cfg, &testProxy{}))
	assert.NotNil(t, c.initConfig(cfg, &testProxy{}))

	assert.Equal(t, uint32(1), c.disp.root.wallpaper.ID)
	ag, err := c.getIcon("airplane-green")
	assert.Nil(t, err)
	assert.Equal(t, uint32(1), ag.ID)

	assert.Equal(t, 3, len(c.resources.Load().(map[string]any)))

	_, err = c.getIcon("airplane-blue")
	assert.Equal(t, errors.ErrNotExist, err)

	assert.Equal(t, &c.disp.root, RootContainer())
	f := SystemFont(10)
	f1 := SystemFont(20)
	assert.Equal(t, f, f1)
	f1 = SystemFont(200)
	assert.NotEqual(t, f, f1)
}

func Test_controller_checkFileName(t *testing.T) {
	cfg := Config{
		DisplayConfig: DefaultDisplayConfig(),
	}
	c := &controller{}
	assert.Nil(t, c.initConfig(cfg, &testProxy{}))
	filename := filepath.FromSlash("testdata/icons/airplane-green.png")
	fn, err := c.checkFileName("", filename)
	assert.Nil(t, err)
	assert.Equal(t, filename, fn)

	fn, err = c.checkFileName("testdata", filepath.FromSlash("icons/airplane-green.png"))
	assert.Nil(t, err)
	assert.Equal(t, filename, fn)

	fn, err = c.checkFileName("", filepath.FromSlash("icons/airplane-green.png"))
	assert.NotNil(t, err)
}
