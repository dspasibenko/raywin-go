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
	"github.com/dspasibenko/raywin-go/pkg/golibs/container"
	"github.com/dspasibenko/raywin-go/pkg/golibs/container/lru"
	"github.com/dspasibenko/raywin-go/pkg/golibs/errors"
	"github.com/dspasibenko/raywin-go/pkg/golibs/files"
	"github.com/dspasibenko/raywin-go/pkg/golibs/logging"
	rl "github.com/gen2brain/raylib-go/raylib"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

// Init should be called before the Run() to initialize the raywin-go
func Init(cfg Config) error {
	return c.initConfig(cfg, p)
}

// Run runs the drawing cycle and rendering the main window. It will be stopped when
// the context ctx is closed
func Run(ctx context.Context) error {
	return c.disp.run(ctx)
}

// RootContainer returns the container for the display
func RootContainer() Container {
	return &c.disp.root
}

// SystmeFont returns the default system font
func SystemFont(size int) rl.Font {
	return Font(c.cfg.RegularFontFileName, size)
}

// Millis returns the current raywin timestamp. This is not the clock time,
// but a reference time used in raywin functions
func Millis() int64 {
	return c.disp.millis.Load()
}

// SystemItalicFont returns the Italic version of the system font
func SystemItalicFont(size int) rl.Font {
	return Font(c.cfg.ItalicFontFileName, size)
}

const fontCacheScaleFactor = 97

// Font returns the rl.Font for the requested size points (1/72")
func Font(fontFile string, size int) rl.Font {
	f := fmt.Sprintf("%s%%%d", fontFile, size/fontCacheScaleFactor)
	font, _ := c.fontsCache.GetOrCreate(f)
	return font
}

// GetIcon returns the icon by its name without the extension. If the file name is
// "picture.png" it can be obtained by "picture". See Config
func GetIcon(in string) (rl.Texture2D, error) {
	return c.getIcon(in)
}

type controller struct {
	logger     logging.Logger
	lock       sync.Mutex
	resources  atomic.Value
	cfg        Config
	disp       *display
	valid      atomic.Bool
	fontsCache *lru.Cache[string, rl.Font]
}

var p RlProxy = &realProxy{}
var c = &controller{}

func assertInitialized() {
	if !c.valid.Load() {
		panic("raywin is not initialized (call raywin.Init())")
	}
}

func (c *controller) initConfig(cfg Config, proxy RlProxy) error {
	if !c.valid.CompareAndSwap(false, true) {
		return fmt.Errorf("initConfig: already initialized: %w", errors.ErrInvalid)
	}
	c.logger = logging.NewLogger("raywin")
	c.disp = newDisplay(cfg.DisplayConfig, proxy)
	c.disp.frmListener = cfg.FrameListener
	c.resources.Store(map[string]any{})
	c.cfg = cfg
	c.fontsCache, _ = lru.NewCache[string, rl.Font](20, func(cacheKey string) (rl.Font, error) {
		s := strings.Split(cacheKey, "%")
		if len(s) != 2 {
			return rl.Font{}, fmt.Errorf("invalid cache key: %s, expecting \"fileName%%size\"", cacheKey)
		}
		sz, err := strconv.Atoi(s[1])
		if err != nil {
			return rl.Font{}, fmt.Errorf("invalid cache key: %s, expecting \"fileName%%size\", size=%s cannot be parsed as int", cacheKey, s[1])
		}
		sz = max(1, sz)
		return c.loadFont("system italic font", cfg.ResourceDir, s[0], int32(sz*fontCacheScaleFactor))
	}, nil)
	img, err := c.loadImage("wallpaper", cfg.ResourceDir, cfg.WallpaperFileName)
	if err != nil {
		return err
	}
	if img != nil {
		c.logger.Infof("using wallpaper from the config file %s", cfg.WallpaperFileName)
		c.disp.root.wallpaper = c.disp.proxy.LoadTextureFromImage(img)
	}
	if err := c.loadIcons(cfg.ResourceDir, cfg.IconsDir); err != nil {
		return err
	}
	return nil
}

func (c *controller) loadImage(comment, dir, fn string) (*rl.Image, error) {
	if fn == "" {
		c.logger.Infof("%s image file is not specified, skip it", comment)
		return nil, nil
	}
	fn, err := c.checkFileName(dir, fn)
	if err != nil {
		return nil, fmt.Errorf("%s file could not be opened: %w", fn, err)
	}
	img := rl.LoadImage(fn)
	if img == nil {
		return nil, fmt.Errorf("could not load image from file %s", fn)
	}
	c.logger.Infof("read %s image file successfully", fn)
	return img, nil
}

func (c *controller) loadIcons(dir, fn string) error {
	if fn == "" {
		c.logger.Warnf("no icons to load, the file dir name is not provided")
		return nil
	}
	fn, err := c.checkFileName(dir, fn)
	if err != nil {
		return fmt.Errorf("could not open icons dir: %w", err)
	}
	files := files.ListDir(fn)
	for _, f := range files {
		ext := filepath.Ext(f.Name())
		if ext != ".png" {
			c.logger.Warnf("don't support icon format %s, skipping it", f.Name())
			continue
		}
		img := rl.LoadImage(filepath.Join(fn, f.Name()))
		icoName := f.Name()[:len(f.Name())-len(ext)]
		tx := c.disp.proxy.LoadTextureFromImage(img)
		c.addResouce("ico_"+icoName, tx)
	}
	return nil
}

func (c *controller) getIcon(in string) (rl.Texture2D, error) {
	r := c.resource("ico_" + in)
	if r == nil {
		return rl.Texture2D{}, errors.ErrNotExist
	}
	return r.(rl.Texture2D), nil
}

func (c *controller) loadFont(comment, dir, fn string, fontSize int32) (rl.Font, error) {
	if fn == "" {
		c.logger.Infof("%s font is not specified, skip it", comment)
		return rl.Font{}, nil
	}
	fn, err := c.checkFileName(dir, fn)
	if err != nil {
		return rl.Font{}, fmt.Errorf("%s file %s file could not be opened: %w", comment, fn, err)
	}
	c.logger.Infof("loading %s from %s", comment, fn)
	f := c.disp.proxy.LoadFontEx(fn, fontSize)
	c.disp.proxy.SetTextureFilter(f.Texture, rl.FilterBilinear)
	return f, nil
}

func (c *controller) checkFileName(dir, fn string) (string, error) {
	if _, err := os.Stat(fn); err != nil {
		fn1 := filepath.Join(dir, fn)
		c.logger.Warnf("could not open file %s, will check %s: %v", fn, fn1, err)
		_, err := os.Stat(fn1)
		if err != nil {
			c.logger.Warnf("could not open file %s either: %v", fn1, err)
		}
		return fn1, err
	}
	return fn, nil
}

func (c *controller) resource(name string) any {
	m := c.resources.Load().(map[string]any)
	return m[name]
}

func (c *controller) addResouce(name string, v any) {
	c.lock.Lock()
	defer c.lock.Unlock()

	m := c.resources.Load().(map[string]any)
	m = container.CopyMap(m)
	m[name] = v
	c.resources.Store(m)
}
