package decimal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromFloat64(t *testing.T) {
	d := FromFloat64(0, 0)
	assert.Equal(t, IntWrapper{}, d)
	d = FromFloat64(-0, 0)
	assert.Equal(t, IntWrapper{}, d)
	d = FromFloat64(-0.3, 0)
	assert.Equal(t, IntWrapper{}, d)
	d = FromFloat64(-0.7, 0)
	assert.Equal(t, IntWrapper{V: -1}, d)
	d = FromFloat64(-0.7, -4)
	assert.Equal(t, IntWrapper{V: -7000, E: -4}, d)
	d = FromFloat64(-7000, 3)
	assert.Equal(t, IntWrapper{V: -7, E: 3}, d)
	d = FromFloat64(57.34523, -3)
	assert.Equal(t, IntWrapper{V: 57345, E: -3}, d)
	d = FromFloat64(-57.34523, -3)
	assert.Equal(t, IntWrapper{V: -57345, E: -3}, d)
}

func TestIntWrapper_String(t *testing.T) {
	assert.Equal(t, IntWrapper{V: 57345, E: -3}.String(), "57.345")
	assert.Equal(t, IntWrapper{V: -1, E: -3}.String(), "-0.001")
	assert.Equal(t, IntWrapper{V: 0, E: -3}.String(), "0.000")
	assert.Equal(t, IntWrapper{V: -123, E: -3}.String(), "-0.123")
	assert.Equal(t, IntWrapper{V: -123123, E: -3}.String(), "-123.123")
	assert.Equal(t, IntWrapper{V: 1, E: -3}.String(), "0.001")
	assert.Equal(t, IntWrapper{V: 123, E: -3}.String(), "0.123")
	assert.Equal(t, IntWrapper{V: 123123, E: -3}.String(), "123.123")
	assert.Equal(t, IntWrapper{V: 1, E: 3}.String(), "1000")
	assert.Equal(t, IntWrapper{V: 123, E: 3}.String(), "123000")
	assert.Equal(t, IntWrapper{V: 123123, E: 3}.String(), "123123000")
	assert.Equal(t, "0.000", IntWrapper{V: 0, E: -3}.String())
	assert.Equal(t, IntWrapper{V: 0, E: 3}.String(), "0")
}

func TestIntWrapper_Float64(t *testing.T) {
	assert.Equal(t, float64(3.22), IntWrapper{322000, -5}.Float64())
	assert.Equal(t, float64(-3.22), IntWrapper{-322, -2}.Float64())
	assert.Equal(t, float64(3.123), IntWrapper{3123, -3}.Float64())
	assert.Equal(t, float64(0.001), IntWrapper{1, -3}.Float64())
}
