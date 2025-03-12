package storage

import (
	"context"
	"errors"
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

func (m *MemoryStore) Read(_ context.Context, id string) (domain.Port, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	port, exists := m.ports[id]
	if !exists {
		return domain.Port{}, errors.New("not found")
	}

	return port, nil
}
