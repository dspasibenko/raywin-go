// Copyright 2023 The acquirecloud Authors
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
package xbinary

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func BenchmarkMarshalUint(b *testing.B) {
	var bb [30]byte
	buf := bb[:]
	ui := uint(1347598723405981734)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sz, _ := MarshalUint(ui, buf)
		copy(bb[10:], bb[:sz])
	}
}

func BenchmarkMarshalInt64(b *testing.B) {
	var bb [30]byte
	buf := bb[:]
	ui := uint64(1347598723405981734)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MarshalUint64(ui, buf)
	}
}

func BenchmarkMarshalBytes(b *testing.B) {
	var bb [30]byte
	buf := bb[:]
	bbb := []byte("test string")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MarshalBytes(bbb, buf)
	}
}

func BenchmarkWriteBytes(b *testing.B) {
	var bb [30]byte
	buf := bb[:]
	bbb := []byte("test string")
	btb := bytes.NewBuffer(buf)
	ow := ObjectsWriter{Writer: btb}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		btb.Reset()
		ow.WritePureBytes(bbb)
	}
}

func BenchmarkWriteUInt(b *testing.B) {
	var bb [30]byte
	buf := bb[:]
	btb := bytes.NewBuffer(buf)
	ow := ObjectsWriter{Writer: btb}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		btb.Reset()
		ow.WriteUint(1234719238471923749)
	}
}

func BenchmarkUnmarshalBytes(b *testing.B) {
	var bb [30]byte
	buf := bb[:]
	bbb := []byte("test string")
	MarshalBytes(bbb, buf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		UnmarshalBytes(buf, true)
	}
}

func BenchmarkUnmarshalBytesNoCopy(b *testing.B) {
	var bb [30]byte
	buf := bb[:]
	bbb := []byte("test string")
	MarshalBytes(bbb, buf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		UnmarshalBytes(buf, false)
	}
}

func TestMarshalUInt(t *testing.T) {
	var b [3]byte
	buf := b[:]
	idx, err := MarshalUint(0, buf)
	if idx != 1 || buf[0] != 0 || err != nil {
		t.Fatal("Unexpected idx=", idx, " buf=", buf, ", err=", err)
	}

	idx, err = MarshalUint(25, buf)
	if idx != 1 || buf[0] != 25 || err != nil {
		t.Fatal("Unexpected idx=", idx, " buf=", buf, ", err=", err)
	}

	idx, v, err := UnmarshalUint(buf)
	if v != 25 || idx != 1 || err != nil {
		t.Fatal("Unexpected idx=", idx, " v=", v, ", err=", err)
	}

	idx, err = MarshalUint(129, buf)
	if idx != 2 || buf[0] != 129 || buf[1] != 1 || err != nil {
		t.Fatal("Unexpected idx=", idx, " buf=", buf, ", err=", err)
	}

	idx, err = MarshalUint(32565, buf)
	if idx != 3 || err != nil {
		t.Fatal("Unexpected idx=", idx, " buf=", buf, ", err=", err)
	}

	idx, v, err = UnmarshalUint(buf)
	if v != 32565 || idx != 3 || err != nil {
		t.Fatal("Unexpected idx=", idx, " v=", v, ", err=", err)
	}

	idx, err = MarshalUint(3256512341234, buf)
	if err == nil {
		t.Fatal("Expecting err != nil, but idx=", idx, " buf=", buf)
	}
}

func TestMarshalBytes(t *testing.T) {
	const str = "abcasdfadfasd"
	bstr := []byte(str)
	var b [20]byte
	buf := b[:]

	idx, err := MarshalBytes(bstr, buf)
	if idx != 1+len(str) || err != nil {
		t.Fatal("idx=", idx, ", err=", err)
	}

	idx2, bts, err := UnmarshalBytes(buf, false)
	if idx != idx2 || string(bts) != str || err != nil {
		t.Fatal("idx2=", idx2, " bts=", string(bts), ", err=", err)
	}

	idx, err = MarshalBytes(bstr[:0], buf)
	if idx != 1 {
		t.Fatal("empty bytes should be 1 byte in result length")
	}
}

func TestMarshalByte(t *testing.T) {
	testMarshalInts(t, 254, 1, func(v int, b []byte) (int, error) {
		return MarshalByte(byte(v), b)
	}, func(b []byte) (int, int, error) {
		i, v, err := UnmarshalByte(b)
		return i, int(v), err
	})
}

func TestMarshalUint16(t *testing.T) {
	testMarshalInts(t, 254, 2, func(v int, b []byte) (int, error) {
		return MarshalUint16(uint16(v), b)
	}, func(b []byte) (int, int, error) {
		i, v, err := UnmarshalUint16(b)
		return i, int(v), err
	})
}

func TestMarshalUint32(t *testing.T) {
	testMarshalInts(t, 223454, 4, func(v int, b []byte) (int, error) {
		return MarshalUint32(uint32(v), b)
	}, func(b []byte) (int, int, error) {
		i, v, err := UnmarshalUint32(b)
		return i, int(v), err
	})
}

func TestMarshalInt64(t *testing.T) {
	testMarshalInts(t, 2542341, 8, func(v int, b []byte) (int, error) {
		return MarshalUint64(uint64(v), b)
	}, func(b []byte) (int, int, error) {
		i, v, err := UnmarshalUint64(b)
		return i, int(v), err
	})
}

func TestSizeUint(t *testing.T) {
	testSizeUint(t, bit7, 2)
	testSizeUint(t, bit14, 3)
	testSizeUint(t, bit21, 4)
	testSizeUint(t, bit28, 5)
	testSizeUint(t, bit35, 6)
	testSizeUint(t, bit42, 7)
	testSizeUint(t, bit49, 8)
	testSizeUint(t, bit56, 9)
	testSizeUint(t, bit63, 10)
}

func testSizeUint(t *testing.T, val uint64, sz int) {
	if WritableUintSize(val) != sz || WritableUintSize(val-1) != sz-1 {
		t.Fatal("for ", val, " expecting size ", sz, ", but it is ", WritableUintSize(val))
	}

	var b [20]byte
	buf := b[:]
	idx, _ := MarshalUint(uint(val), buf)
	if idx != sz {
		t.Fatal("MarshalUInt returns another value for val=", val, " idx=", idx, ", but sz=", sz)
	}

	idx, _ = MarshalUint(uint(val-1), buf)
	if idx != sz-1 {
		t.Fatal("MarshalUInt returns another value for val=", val, " idx=", idx, ", but sz=", sz-1)
	}
}

func testMarshalInts(t *testing.T, v, elen int, mf func(v int, b []byte) (int, error), uf func(b []byte) (int, int, error)) {
	var b [20]byte
	buf := b[:]

	idx, err := mf(v, buf)
	if idx != elen || err != nil {
		t.Fatal("expected len=", elen, ", but idx=", idx, ", buf=", buf, ", err=", err)
	}

	idx2, val, err := uf(buf)
	if idx != idx2 || v != val || err != nil {
		t.Fatal("unmarshal idx=", idx2, " val=", val, ", but expected=", v, ", err=", err)
	}
}

func TestObjectWriter(t *testing.T) {
	btb := bytes.NewBuffer(nil)
	ow := ObjectsWriter{Writer: btb}
	i1, _ := ow.WriteUint(12341341234134)
	i2, v, err := UnmarshalUint(btb.Bytes())
	if v != 12341341234134 || err != nil || i2 != i1 {
		t.Fatal("err=", err, " v=", v, "i2=", i2, ", i1=", i1, " buf=", btb.Bytes())
	}

	btb.Reset()
	buf := []byte{1, 2, 3}
	ow.WriteBytes(buf)
	_, b, _ := UnmarshalBytes(btb.Bytes(), false)
	if !reflect.DeepEqual(b, buf) {
		t.Fatal("Expected", buf, " but read ", b)
	}

	btb.Reset()
	ow.WriteUint64(uint64(4857293487592347598))
	_, i64, _ := UnmarshalUint64(btb.Bytes())
	if i64 != 4857293487592347598 {
		t.Fatal("Unexpected i64=", i64)
	}

	btb.Reset()
	n, err := ow.WriteByteWithSize(1)
	assert.Nil(t, err)
	assert.Equal(t, 1, n)

	btb.Reset()
	n, err = ow.WriteUint16(1234)
	assert.Nil(t, err)
	assert.Equal(t, 2, n)

	btb.Reset()
	n, err = ow.WriteUint32(1234)
	assert.Nil(t, err)
	assert.Equal(t, 4, n)

}
