.PHONY: build run stop test lint clean logs help

build:
	@echo "Building Raft KV Store..."
	@go build -o bin/raft-node ./cmd/raft-node
	@go build -o bin/kv-gateway ./cmd/kv-gateway

run:
	@echo "Starting Raft cluster and gateway..."
	@docker-compose up -d

stop:
	@echo "Stopping services..."
	@docker-compose down

test:
	@echo "Running tests..."
	@go test ./...

lint:
	@echo "Running linter..."
	@golangci-lint run ./...

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/

logs:
	@docker-compose logs -f

help:
	@echo "Available targets:"
	@echo "  make build   - Build all binaries"
	@echo "  make run     - Start services with Docker Compose"
	@echo "  make stop    - Stop all services"
	@echo "  make test    - Run tests"
	@echo "  make lint    - Run linter"
	@echo "  make clean   - Clean build artifacts"
	@echo "  make logs    - View service logs"
