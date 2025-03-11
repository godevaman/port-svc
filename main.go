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

	log.Println("Service shutdown complete.")
}