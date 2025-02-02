// Copyright 2023 The acquirecloud Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package files

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func TestOpenCloseMMFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "TestOpenCloseMMFile")
	assert.Nil(t, err)
	defer os.RemoveAll(dir) // clean up

	fsz := int64(23451 * 4096)
	fn := path.Join(dir, "testFile")
	assert.Nil(t, EnsureFileExists(fn))
	mmf, err := NewMMFile(fn, fsz)
	assert.Nil(t, err)
	defer mmf.Close()

	buf := []byte{1, 2, 3, 4, 5}
	res, err := mmf.Buffer(12345, len(buf))
	assert.Nil(t, err)
	copy(res, buf)

	res, err = mmf.Buffer(mmf.Size()-2, len(buf))
	assert.Nil(t, err)
	assert.Equal(t, 2, len(res))
	copy(res, buf)

	mmf.Close()
	assert.NotEqual(t, -1, mmf.Size())

	mmf, err = NewMMFile(fn, -1)
	assert.Nil(t, err)
	defer mmf.Close()
	assert.Equal(t, fsz, mmf.Size())

	res, err = mmf.Buffer(12345, len(buf))
	assert.Nil(t, err)
	assert.Equal(t, buf, res)

	res, err = mmf.Buffer(mmf.Size()-3, len(buf))
	assert.Nil(t, err)
	assert.Equal(t, buf[:2], res[1:3])
}

func TestCreateNotExistingMMFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "TestCreateNotExistingMMFile")
	assert.Nil(t, err)
	defer os.RemoveAll(dir) // clean up

	_, err = NewMMFile(path.Join(dir, "testFile"), -1)
	assert.NotNil(t, err)
}

func TestGrowMMFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "TestGrowMMFile")
	assert.Nil(t, err)
	defer os.RemoveAll(dir) // clean up

	fsz := int64(8 * 4096)
	fn := path.Join(dir, "testFile")
	assert.Nil(t, EnsureFileExists(fn))
	mmf, err := NewMMFile(fn, fsz)
	assert.Nil(t, err)
	defer mmf.Close()

	buf := []byte{1, 2, 3, 4, 5}
	res, err := mmf.Buffer(4093, len(buf))
	assert.Nil(t, err)
	n := copy(res, buf)
	assert.Nil(t, err)
	assert.Equal(t, n, len(buf))

	err = mmf.Grow(4 * 4096)
	assert.NotNil(t, err)

	err = mmf.Grow(20*4096 - 1)
	assert.NotNil(t, err)
	assert.Equal(t, fsz, mmf.Size())

	err = mmf.Grow(20 * 4096)
	assert.Nil(t, err)
	assert.Equal(t, int64(20*4096), mmf.Size())

	res, err = mmf.Buffer(4093, len(buf))
	assert.Nil(t, err)
	assert.Equal(t, buf, res)
}

func TestSyncMMFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "TestSyncMMFile")
	assert.Nil(t, err)
	defer os.RemoveAll(dir) // clean up

	fsz := int64(2 * 4096)
	fn := path.Join(dir, "testFile")
	assert.Nil(t, EnsureFileExists(fn))
	mmf, err := NewMMFile(fn, fsz)
	assert.Nil(t, err)
	defer mmf.Close()

	fi, err := os.Stat(fn)
	assert.Nil(t, err)

	fi1, err := os.Stat(fn)
	assert.Nil(t, err)
	assert.Equal(t, fi1.ModTime(), fi.ModTime())

	mmf.Flush()
	fi1, err = os.Stat(fn)
	assert.Nil(t, err)
	assert.Equal(t, fi1.ModTime(), fi.ModTime())

	buf := []byte{1, 2, 3, 4, 5}
	res, err := mmf.Buffer(0, len(buf))
	copy(res, buf)

	mmf.Flush()
	fi1, err = os.Stat(fn)
	assert.Nil(t, err)
	assert.True(t, fi1.ModTime().After(fi.ModTime()))
}

func TestParallelMMFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "TestParrallelMMFile")
	assert.Nil(t, err)
	defer os.RemoveAll(dir) // clean up

	ps := os.Getpagesize()

	fsz := int64(ps * 10)
	fn := path.Join(dir, "testFile2")
	EnsureFileExists(fn)
	mmf, err := NewMMFile(fn, fsz)
	assert.Nil(t, err)
	defer mmf.Close()

	var wg sync.WaitGroup
	var errs int32

	for i := 0; i < 640; i++ {
		wg.Add(1)
		go func(pid int) {
			buf := make([]byte, 64)
			for j, _ := range buf {
				buf[j] = byte(pid)
			}
			for i := 0; i < 500; i++ {
				res, err := mmf.Buffer(int64(pid*64), len(buf))
				if len(res) != len(buf) || err != nil {
					fmt.Println("Error when write n=", len(res), ", err=", err)
					atomic.AddInt32(&errs, 1)
				}
				copy(res, buf)

				res, err = mmf.Buffer(int64(pid*64), len(buf))
				if len(res) != len(buf) || err != nil {
					fmt.Println("Error when read n=", len(res), ", err=", err)
					atomic.AddInt32(&errs, 1)
				}

				if !reflect.DeepEqual(buf, res) {
					fmt.Println("Error buf=", buf, ", res=", res)
					atomic.AddInt32(&errs, 1)
					break
				}
			}

			wg.Done()
		}(i)
	}
	wg.Wait()
	if atomic.LoadInt32(&errs) != 0 {
		t.Fatal(" errs=", atomic.LoadInt32(&errs))
	}
}

func TestMMFileTwoReaders(t *testing.T) {
	dir, err := os.MkdirTemp("", "TestParrallelMMFile")
	assert.Nil(t, err)
	defer os.RemoveAll(dir) // clean up

	ps := os.Getpagesize()

	fsz := int64(ps * 10)
	fn := path.Join(dir, "testFile")
	EnsureFileExists(fn)
	mmf1, err := NewMMFile(fn, fsz)
	assert.Nil(t, err)
	defer mmf1.Close()
	b1, err := mmf1.Buffer(0, 512)
	assert.Nil(t, err)

	mmf2, err := NewMMFile(fn, fsz)
	assert.Nil(t, err)
	defer mmf2.Close()
	b2, err := mmf2.Buffer(0, 512)
	assert.Nil(t, err)

	assert.Nil(t, os.Remove(fn))
	go func() {
		idx := 0
		for idx < 256 {
			if b2[idx] == byte(idx) {
				idx++
			} else {
				runtime.Gosched()
			}
		}
		fmt.Println("mmf2: read 256 bytes from mmf2")
		for idx < 512 {
			b2[idx] = byte(idx - 256)
			idx++
		}
		mmf2.Flush()
		fmt.Println("mmf2: 256 bytes written to mmf2")
	}()

	idx := 0
	for idx < 256 {
		b1[idx] = byte(idx)
		idx++
	}
	mmf1.Flush()
	fmt.Println("mmf1: 256 bytes written to mmf2")

	for idx < 512 {
		if b1[idx] == byte(idx-256) {
			idx++
		} else {
			runtime.Gosched()
		}
	}
	fmt.Println("mmf1: read 256 bytes from mmf2")
}
