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
package context

import (
	ctx "context"
	"fmt"
	errors2 "github.com/dspasibenko/raywin-go/pkg/golibs/errors"
	"time"
)

type (
	closingCtx struct {
		ch <-chan struct{}
	}
)

var _ ctx.Context = (*closingCtx)(nil)

// WrapChannel receives a channel and returns a context which wraps the channel
// The context will be closed when the channel is closed.
func WrapChannel(ch <-chan struct{}) ctx.Context {
	return &closingCtx{ch}
}

// Deadline is a part of context.Context interface
func (cc *closingCtx) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

// Done is a part of context.Context interface
func (cc *closingCtx) Done() <-chan struct{} {
	return cc.ch
}

// Err is a part of context.Context interface
func (cc *closingCtx) Err() error {
	select {
	case _, ok := <-cc.ch:
		if ok {
			panic("Improper use of the the context wrapper")
		}
		return fmt.Errorf("The underlying channel was closed: %w", errors2.ErrClosed)
	default:
		return nil
	}
}

// Value is a part of context.Context interface
func (cc *closingCtx) Value(key interface{}) interface{} {
	return nil
}
