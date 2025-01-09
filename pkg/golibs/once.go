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
