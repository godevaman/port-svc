FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o port-service main.go

# Final stage: lightweight runtime image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/port-service .
COPY ports.json .

# Run as non-root user for security
RUN adduser -D appuser
USER appuser

ENTRYPOINT ["./port-service"]
