package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"port-service/storage"
	"port-service/service"
	"sync"
	"syscall"
)

func main() {
	// Initialize in-memory storage
	store := storage.NewMemoryStore()

	// Create port service
	portSvc := service.NewPortService(store)

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Handle shutdown or errors
	select {
	case <-sigChan:
		log.Println("Received shutdown signal, gracefully shutting down...")
		cancel()
	case err := <-errChan:
		log.Printf("Error processing ports: %v", err)
		cancel()
	}

	log.Println("Service shutdown complete.")
}