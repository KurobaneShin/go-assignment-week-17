package main

import (
	"sync"
)

type MemoryObject struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewMemoryObject() *MemoryObject {
	return &MemoryObject{
		data: map[string]string{},
	}
}

func (mo *MemoryObject) Set(key, value string) error {
	mo.mu.Lock()

	defer mo.mu.Unlock()

	mo.data[key] = value

	return nil
}

func (mo *MemoryObject) Get(key string) (string, bool) {
	mo.mu.RLock()

	defer mo.mu.RUnlock()

	val, ok := mo.data[key]

	return val, ok
}
