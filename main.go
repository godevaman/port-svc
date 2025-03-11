package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"port-service/domain"
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


	// Open ports.json file
	file, err := os.Open("ports.json")
	if err != nil {
		log.Fatalf("failed to open ports.json: %v", err)
	}
	defer file.Close()

	// Process file in a streaming manner
	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		decoder := json.NewDecoder(file)
		_, err := decoder.Token() // Skip opening brace
		if err != nil {
			errChan <- err
			return
		}

		for decoder.More() {
			var port domain.Port
			var key string
			if err := decoder.Decode(&key); err != nil {
				errChan <- err
				return
			}
			if err := decoder.Decode(&port); err != nil {
				errChan <- err
				return
			}
			port.ID = key
			if err := portSvc.Upsert(ctx, port); err != nil {
				errChan <- err
				return
			}
		}
	}()

	// Handle shutdown or errors
	select {
	case <-sigChan:
		log.Println("Received shutdown signal, gracefully shutting down...")
		cancel()
	case err := <-errChan:
		log.Printf("Error processing ports: %v", err)
		cancel()
	}

	wg.Wait()
	log.Println("Service shutdown complete.")
}