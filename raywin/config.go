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
	"encoding/json"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type (
	// Config struct describes the raywin-go configuration
	Config struct {
		// ResourceDir specify the path to dir with the library resources(wallpaper, icons, fonts etc.)
		ResourceDir string

		// WallpaperFileName provides the name to .png file the file should be in the current dir or in the
		// ResourceDir, if it is not find in the current directory. The field may be empty.
		WallpaperFileName string

		// RegularFontFileName provides the name to .ttf file with the default system font.
		// The file should be in the current dir or in the ResourceDir, if it is not find
		// in the current directory. The field may be empty.
		RegularFontFileName string

		// ItalicFontFileName provides the name to .ttf file containing italic system font.
		// The file should be in the current dir or in the ResourceDir, if it is not find
		// in the current directory. The field may be empty.
		ItalicFontFileName string

		// DisplayConfig contains the display physical characteristics
		DisplayConfig DisplayConfig

		// IconsDir provides the name of the directory where the icon images (in PNG) are stored.
		// The directory should be in the current dir or in the ResourceDir, if it is not find
		// in the current directory. The field may be empty.
		//
		// All icons will be read into memory during Init() and they will be available via
		// GetIcon() call. The file name (without the extension) is used as the icon name.
		IconsDir string

		// FrameListener allows to specify an external frame listener which will be called
		// on each new frame. It can be nil
		FrameListener FrameListener
	}

	// DisplayConfig contain the basic display configuration
	DisplayConfig struct {
		// Width of the display area in number of pixels
		Width uint32
		// Height of the display area in number of pixels
		Height uint32
		// PPI is the Pixels per Inch, depends on the display physical size
		PPI float32
		// FPS - frames per second. The number of the display updates the library
		// will try to support.
		FPS int
		// BackgroundColor contains the default color for display area
		BackgroundColor rl.Color
	}
)

// DefaultDisplayConfig returns the default DisplayConfig
func DefaultDisplayConfig() DisplayConfig {
	return DisplayConfig{
		Width:           1280,
		Height:          720,
		PPI:             209.8,
		FPS:             60,
		BackgroundColor: rl.Black,
	}
}

// String - DisplayConfig's implementation of fmt.Stringer
func (dc DisplayConfig) String() string {
	b, _ := json.MarshalIndent(dc, "", "  ")
	return string(b)
}

// DefaultConfig returns the default Config
func DefaultConfig() Config {
	return Config{
		DisplayConfig: DefaultDisplayConfig(),
	}
}

// String - Config's implementation of fmt.Stringer
func (c Config) String() string {
	b, _ := json.MarshalIndent(c, "", "  ")
	return string(b)
}
