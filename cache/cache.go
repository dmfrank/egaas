package cache

import (
	"sync"
)

// Auth stores sync auth information
type Auth struct {
	sync.Mutex
	Values map[string]string
}

// Work keeps
type Work struct {
	sync.Mutex
	Values map[string]int32
}

// IsExist checks availability of item
func (a *Auth) IsExist(k, v string) bool {
	a.Lock()
	defer a.Unlock()

	if val, ok := a.Values[k]; ok {
		if val == v {
			return true
		}
	}
	return false
}

// Push keeps item in Auth map
func (a *Auth) Push(k, v string) {
	a.Lock()
	defer a.Unlock()

	a.Values[k] = v
}

// Push keeps item in Work map
func (w *Work) Push(k string, v int32) {
	w.Lock()
	defer w.Unlock()

	w.Values[k] = v
}
