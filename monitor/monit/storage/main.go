package storage

import (
	"sync"
)

func NewManager() *Manager {
	return &Manager{
		items:     make(map[string]*item),
		threshold: 3,
		mutex:     sync.Mutex{},
	}
}
