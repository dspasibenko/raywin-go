package logging

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAll(t *testing.T) {
	SetLevel(DEBUG)
	assert.Equal(t, DEBUG, GetLevel())
	l := NewLogger("test")
	l.Infof("hello world %d", 1)
	l.Warnf("hello world %d", 2)
	l.Errorf("hello world %d", 2)
	l.Debugf("hello world %d", 2)
	l.Tracef("hello world %d", 3)
}
