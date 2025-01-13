package raywin

import (
	"github.com/dspasibenko/raywin-go/pkg/golibs/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_controller_initConfig(t *testing.T) {
	cfg := Config{
		DisplayConfig:       DefaultDisplayConfig(),
		WallpaperFileName:   "testdata/images/wallpaper800x.png",
		RegularFontFileName: "testdata/fonts/Roboto/Roboto-Medium.ttf",
		ItalicFontFileName:  "testdata/fonts/Roboto/Roboto-MediumItalic.ttf",
		IconsDir:            "testdata/icons",
	}
	c := &controller{}
	assert.Nil(t, c.initConfig(cfg, &testProxy{}))

	assert.Equal(t, uint32(1), c.disp.root.wallpaper.ID)
	ag, err := c.getIcon("airplane-green")
	assert.Nil(t, err)
	assert.Equal(t, uint32(1), ag.ID)

	assert.Equal(t, 3, len(c.resources.Load().(map[string]any)))

	_, err = c.getIcon("airplane-blue")
	assert.Equal(t, errors.ErrNotExist, err)
}

func Test_controller_checkFileName(t *testing.T) {
	cfg := Config{
		DisplayConfig: DefaultDisplayConfig(),
	}
	c := &controller{}
	assert.Nil(t, c.initConfig(cfg, &testProxy{}))
	filename := "testdata/icons/airplane-green.png"
	fn, err := c.checkFileName("", filename)
	assert.Nil(t, err)
	assert.Equal(t, filename, fn)

	fn, err = c.checkFileName("testdata", "icons/airplane-green.png")
	assert.Nil(t, err)
	assert.Equal(t, filename, fn)

	fn, err = c.checkFileName("", "icons/airplane-green.png")
	assert.NotNil(t, err)
}
