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
	"fmt"
	"github.com/dspasibenko/raywin-go/pkg/golibs/container"
	"github.com/dspasibenko/raywin-go/pkg/golibs/errors"
	"github.com/dspasibenko/raywin-go/pkg/golibs/files"
	"github.com/dspasibenko/raywin-go/pkg/golibs/logging"
	rl "github.com/gen2brain/raylib-go/raylib"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

func Init(cfg Config) error {
	return c.initConfig(cfg)
}

func Run() error {
	return c.disp.run()
}

func RootContainer() Container {
	return &c.disp.root
}

func SystemFont() rl.Font {
	return c.sysFont
}

func SystemItalicFont() rl.Font {
	return c.sysItalicFont
}

func GetIcon(in string) (rl.Texture2D, error) {
	return c.getIcon(in)
}

type controller struct {
	logger        logging.Logger
	lock          sync.Mutex
	resources     atomic.Value
	cfg           Config
	sysFont       rl.Font
	sysItalicFont rl.Font
	disp          *display
}

var c *controller = &controller{}

func (c *controller) initConfig(cfg Config) error {
	c.logger = logging.NewLogger("raywin")
	c.disp = newDisplay(cfg.DisplayConfig)
	c.resources.Store(map[string]any{})
	c.cfg = cfg
	img, err := c.loadImage("wallpaper", cfg.ResourceDir, cfg.WallpaperFileName)
	if err != nil {
		return err
	}
	if img != nil {
		c.logger.Infof("using wallpaper from the config file %s", cfg.WallpaperFileName)
		c.disp.root.wallpaper = rl.LoadTextureFromImage(img)
	}
	c.sysFont, err = c.loadFont("system font", cfg.ResourceDir, cfg.RegularFontFileName)
	if err != nil {
		return err
	}
	c.sysItalicFont, err = c.loadFont("system italic font", cfg.ResourceDir, cfg.ItalicFontFileName)
	if err != nil {
		return err
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
		tx := rl.LoadTextureFromImage(img)
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

func (c *controller) loadFont(comment, dir, fn string) (rl.Font, error) {
	if fn == "" {
		c.logger.Infof("%s font is not specified, skip it", comment)
		return rl.Font{}, nil
	}
	fn, err := c.checkFileName(dir, fn)
	if err != nil {
		return rl.Font{}, fmt.Errorf("%s file %s file could not be opened: %w", comment, fn, err)
	}
	c.logger.Infof("loading %s from %s", comment, fn)
	return rl.LoadFontEx(fn, 320, nil), nil
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
