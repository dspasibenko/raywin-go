// Copyright 2018-2019 The logrange Authors
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
package strutil

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveDups(t *testing.T) {
	e := []string{"a", "b", "c", "d"}
	a := RemoveDups([]string{"a", "a", "b", "c", "b", "c", "d"})

	assert.NotNil(t, a)
	assert.ElementsMatch(t, a, e)
}

func TestSwapEvenOdd(t *testing.T) {
	assert.Equal(t, []string{}, SwapEvenOdd([]string{}))
	assert.Equal(t, []string{"a"}, SwapEvenOdd([]string{"a"}))
	assert.Equal(t, []string{"a", "b"}, SwapEvenOdd([]string{"b", "a"}))
	assert.Equal(t, []string{"a", "b", "c"}, SwapEvenOdd([]string{"b", "a", "c"}))
	assert.Equal(t, []string{"a", "b", "c", "d"}, SwapEvenOdd([]string{"b", "a", "d", "c"}))
}

func TestTruncateWithEllipses(t *testing.T) {
	assert.Equal(t, "zhuk", TruncateWithEllipses("zhuk", 100))
	assert.Equal(t, "...", TruncateWithEllipses("zhuk", 2))
	assert.Equal(t, "zhuk...", TruncateWithEllipses("zhukzhukzhuk", 7))
}

func TestGetRandomString(t *testing.T) {
	s := GetRandomString(1, "a")
	assert.Equal(t, "a", s)
	s = GetRandomString(100, "a")
	assert.Equal(t, strings.Repeat("a", 100), s)
	s = RandomString(100)
	assert.Equal(t, 100, len(s))
	assert.Equal(t, "", GetRandomString(-5, "a"))
	assert.Equal(t, "", GetRandomString(0, "a"))
}

func TestBytes2String(t *testing.T) {
	assert.Equal(t, "", Bytes2String(nil, "a", 8))
	assert.Equal(t, "", Bytes2String([]byte{1}, "a", 0))
	assert.Equal(t, "aaaaaaaa", Bytes2String([]byte{1}, "a", 1))
	assert.Equal(t, "baaaaaaa", Bytes2String([]byte{1}, "ab", 1))
	assert.Equal(t, "1F", Bytes2String([]byte{0xF1}, "0123456789ABCDEF", 4))
	assert.Equal(t, "123450", Bytes2String([]byte{0b11010001, 0b01011000}, "0123456789ABCDEF", 3))
}
