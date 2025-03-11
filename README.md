# Port Service

A microservice to process and store port data from ports.json using Golang and Docker.

## Prerequisites
- Go 1.21+
- Docker
- A ports.json file in the project root (not included in this repo)

## How to Run
1. *Build the binary locally*:
   ```bash
   go build -o port-service main.go
   ./port-service
2. ```docker build -t port-service .
   ```docker run --rm -v $(pwd)/ports.json:/app/ports.json port-service