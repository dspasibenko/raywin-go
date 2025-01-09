package golibs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOnce_Do(t *testing.T) {
	cnt := 0
	var o Once
	o.Do(func() { cnt++ })
	assert.Equal(t, 1, cnt)
	o.Do(func() { cnt++ })
	assert.Equal(t, 1, cnt)
}
