package components

import (
	"github.com/dspasibenko/raywin-go/raywin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultStyleOutlet(t *testing.T) {
	cfg := raywin.DefaultDisplayConfig()
	fl := DefaultStyleOutlet(cfg)
	assert.Equal(t, float32(100), S.Scale)
	fl.OnNewFrame(0)
	assert.Equal(t, initDefaultStyle(cfg), S)
}
