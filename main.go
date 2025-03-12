package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"port-service/domain"
	"port-service/service"
	"port-service/storage"
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
		decoder := json.NewDecoder(file) // No need for bufio.NewReader

		// Read the opening `{`
		tok, err := decoder.Token()
		if err != nil {
			errChan <- fmt.Errorf("error reading opening brace: %w", err)
			return
		}
		if delim, ok := tok.(json.Delim); !ok || delim != '{' {
			errChan <- fmt.Errorf("expected { at start of JSON, got %v", tok)
			return
		}

		// Read key-value pairs
		for decoder.More() {
			// Read the key (port code)
			tok, err := decoder.Token()
			if err != nil {
				errChan <- fmt.Errorf("error reading port code token: %w", err)
				return
			}

			portCode, ok := tok.(string)
			if !ok {
				errChan <- fmt.Errorf("expected string port code, got %v", tok)
				return
			}

			// Decode the value (port details)
			var port domain.Port
			if err := decoder.Decode(&port); err != nil {
				errChan <- fmt.Errorf("error decoding port details for %s: %w", portCode, err)
				return
			}

			// Assign and upsert into storage
			port.ID = portCode
			if err := portSvc.Upsert(ctx, port); err != nil {
				errChan <- fmt.Errorf("error upserting port: %w", err)
				return
			}

			fmt.Printf("Processed: %s - %s, %s\n", portCode, port.Name, port.Country)
		}

		port, err := portSvc.Read(ctx, "AEAJM")
		if err != nil {
			fmt.Println("Port not found!")
		}
		fmt.Println("Reading port", port.Name)

		// Read the closing `}`
		tok, err = decoder.Token()
		if err != nil {
			errChan <- fmt.Errorf("error reading closing brace: %w", err)
			return
		}
		if delim, ok := tok.(json.Delim); !ok || delim != '}' {
			errChan <- fmt.Errorf("expected } at end of JSON, got %v", tok)
			return
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
