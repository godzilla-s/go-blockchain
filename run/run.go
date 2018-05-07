package run

import (
	"sync"
)

// for test
type handle = func()

var functionTest = make(map[string]handle)
var mux sync.Mutex

func Register(name string, f handle) {
	mux.Lock()
	defer mux.Unlock()

	functionTest[name] = f
}

func GetFunctions() map[string]handle {
	return functionTest
}
