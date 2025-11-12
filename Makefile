.PHONY: help dev build run test clean docker-up docker-down migrate lint format

# Default target
help:
	@echo "Available commands:"
	@echo "  make dev          - Run with hot reload (Air)"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make test         - Run tests"
	@echo "  make test-cover   - Run tests with coverage"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make docker-up    - Start Docker services"
	@echo "  make docker-down  - Stop Docker services"
	@echo "  make lint         - Run linters"
	@echo "  make format       - Format code"
	@echo "  make deps         - Download dependencies"

# Development with hot reload
dev:
	air

# Build the application
build:
	go build -o bin/api .

# Run the application
run:
	go run main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-cover:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin tmp coverage.out coverage.html

# Start Docker services (PostgreSQL, MongoDB, Redis)
docker-up:
	docker-compose up -d

# Stop Docker services
docker-down:
	docker-compose down

# Run linters
lint:
	golangci-lint run

# Format code
format:
	goimports -w .
	gofmt -s -w .

# Download dependencies
deps:
	go mod download
	go mod tidy

# Initialize Go module (if not exists)
init:
	go mod init github.com/kaiqueyamamoto/sr_robot/api

# Install development tools
tools:
	go install github.com/cosmtrek/air@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest

# Database migrations (placeholder)
migrate-up:
	@echo "Running migrations..."
	# Add your migration tool here (e.g., golang-migrate, goose, etc.)

migrate-down:
	@echo "Rolling back migrations..."
	# Add your migration tool here

# Generate Swagger docs
swagger:
	swag init

# Run the API in production mode
prod:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/api .
	./bin/api

