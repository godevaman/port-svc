package storage

import (
	"context"
	"port-service/domain"
	"sync"
)

// MemoryStore is an in-memory implementation of the Store interface.
type MemoryStore struct {
	ports map[string]domain.Port
	mu    sync.RWMutex
}

// NewMemoryStore initializes a new MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		ports: make(map[string]domain.Port),
	}
}

func (m *MemoryStore) Upsert(_ context.Context, port domain.Port) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ports[port.ID] = port
	return nil
}