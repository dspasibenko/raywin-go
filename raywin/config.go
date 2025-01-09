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
	"encoding/json"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type (
	Config struct {
		ResourceDir         string
		WallpaperFileName   string
		RegularFontFileName string
		ItalicFontFileName  string
		DisplayConfig       DisplayConfig
		IconsDir            string
	}

	DisplayConfig struct {
		Width           uint32
		Height          uint32
		FPS             int
		BackgroundColor rl.Color
	}
)

func DefaultDisplayConfig() DisplayConfig {
	return DisplayConfig{
		Width:           800,
		Height:          480,
		FPS:             60,
		BackgroundColor: rl.Black,
	}
}

func (dc DisplayConfig) String() string {
	b, _ := json.MarshalIndent(dc, "", "  ")
	return string(b)
}

func DefaultConfig() Config {
	return Config{
		DisplayConfig: DefaultDisplayConfig(),
	}
}

func (c Config) String() string {
	b, _ := json.MarshalIndent(c, "", "  ")
	return string(b)
}
