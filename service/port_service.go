package service

import (
	"context"
	"port-service/domain"
)

// Store defines the storage interface for ports.
type Store interface {
	Upsert(ctx context.Context, port domain.Port) error
}

// PortService handles port-related business logic.
type PortService struct {
	store Store
}

// NewPortService creates a new PortService instance.
func NewPortService(store Store) *PortService {
	return &PortService{store: store}
}

// Upsert creates or updates a port in the store.
func (s *PortService) Upsert(ctx context.Context, port domain.Port) error {
	if err := port.Validate(); err != nil {
		return err
	}
	return s.store.Upsert(ctx, port)
}