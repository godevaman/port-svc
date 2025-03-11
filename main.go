package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"port-service/storage"
	"sync"
	"syscall"
)

func main() {
	// Initialize in-memory storage
	store := storage.NewMemoryStore()
	log.Println("Service shutdown complete.")
}