package decimal

import (
	"math"
	"strconv"
	"strings"
)

// IntWrapper is a lightweight container to represent a decimal
// value in int value.
type IntWrapper struct {
	// V contains all digits of a decimal value
	V int
	// E is an exponent value the real decimal value is V*10^E. So
	// if E is less than zero, it means that least significant E digits of V
	// goes to the decimal fraction. For example, IntWrapper{V:1234, E:-2} means 12.34
	// if E is greater than 0, than it is a decimal multiplier, so IntWrapper{V:1234, E:2}
	// means 123400 etc.
	E int
}

var pow = []int{0, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000, 10000000000}

// FromFloat64 returns the wrapper for the float64 value v and the decimal exponent e
func FromFloat64(v float64, e int) IntWrapper {
	if e == 0 {
		return IntWrapper{V: int(math.Round(v))}
	}
	if e < 0 {
		v *= float64(pow[-e])
		return IntWrapper{V: int(math.Round(v)), E: e}
	}
	return IntWrapper{V: int(v / float64(pow[e])), E: e}
}

// String formats d as decimal value
func (d IntWrapper) String() string {
	if d.E < 0 {
		e := -d.E
		div := pow[-d.E]
		var sb strings.Builder
		v := d.V
		if v < 0 {
			sb.WriteString("-")
			v = -v
		}
		sb.WriteString(strconv.Itoa(v / div))
		sb.WriteString(".")
		rest := v % div
		rst := strconv.Itoa(rest)
		if len(rst) < e {
			sb.WriteString(strings.Repeat("0", e-len(rst)))
		}
		sb.WriteString(rst)
		return sb.String()
	}
	if d.E == 0 {
		return strconv.Itoa(d.V)
	}
	return strconv.Itoa(d.V * pow[d.E])
}

// Float64 returs float64 for the d
func (d IntWrapper) Float64() float64 {
	if d.E == 0 {
		return float64(d.V)
	}
	if d.E < 0 {
		return float64(d.V) / float64(pow[-d.E])
	}
	return float64(d.V) * float64(pow[-d.E])
}
