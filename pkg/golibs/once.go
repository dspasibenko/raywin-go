// Copyright 2025 Dmitry Spasibenko
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
package golibs

// Once is a single-thread version of sync.Once
type Once struct {
	called bool
}

// Do calls the func f() only once. Must be NOT used in mutliple go-routines
func (o *Once) Do(f func()) {
	if !o.called {
		f()
		o.called = true
	}
}
