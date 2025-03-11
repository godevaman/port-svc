package tests

import (
	"context"
	"port-service/domain"
	"port-service/service"
	"port-service/storage"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPortService_Upsert(t *testing.T) {
	store := storage.NewMemoryStore()
	srv := service.NewPortService(store)
	ctx := context.Background()

	t.Run("valid port", func(t *testing.T) {
		port := domain.Port{ID: "TEST1", Name: "Test Port", City: "Test City"}
		err := srv.Upsert(ctx, port)
		assert.NoError(t, err)
	})

	t.Run("invalid port", func(t *testing.T) {
		port := domain.Port{ID: "", Name: "Invalid Port"}
		err := srv.Upsert(ctx, port)
		assert.ErrorIs(t, err, domain.ErrInvalidID)
	})
}