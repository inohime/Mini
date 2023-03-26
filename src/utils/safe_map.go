package utils

import "sync"

// A thread-safe map
type SafeMap struct {
	sync.RWMutex
	_map map[string]interface{}
}

// Creates a new thread-safe map
func NewSafeMap() SafeMap {
	return SafeMap{
		_map: make(map[string]interface{}),
	}
}

// Obtain the value for a given key in the map
//
// Returns a untyped value and a check
func (sm *SafeMap) Read(key string) (value interface{}, ok bool) {
	sm.RLock()
	defer sm.RUnlock()

	value, ok = sm._map[key]

	return
}

// Create a new <key, value> pair or modify an existing pair
func (sm *SafeMap) Write(key string, value interface{}) {
	sm.Lock()
	defer sm.Unlock()

	sm._map[key] = value
}

// Delete a <key, value> pair from the map
func (sm *SafeMap) Release(key string) {
	sm.Lock()
	defer sm.Unlock()

	delete(sm._map, key)
}
