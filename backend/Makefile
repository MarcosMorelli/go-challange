# Makefile for Jobsity Backend

.PHONY: test test-unit test-integration test-all build run clean

# Default target
all: test-all

# Run all tests
test-all:
	@echo "Running all tests..."
	go test -v ./tests/...

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	go test -v ./tests/unit/...

# Run integration tests only
test-integration:
	@echo "Running integration tests..."
	go test -v ./tests/integration/...

# Run integration tests with Docker
test-integration-docker:
	@echo "Running integration tests with Docker..."
	./scripts/run-integration-tests.sh

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./tests/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/server cmd/server/main.go

# Run the application
run:
	@echo "Running application..."
	go run cmd/server/main.go

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Help
help:
	@echo "Available targets:"
	@echo "  test-all       - Run all tests"
	@echo "  test-unit      - Run unit tests only"
	@echo "  test-integration - Run integration tests only"
	@echo "  test-integration-docker - Run integration tests with Docker"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  clean          - Clean build artifacts"
	@echo "  deps           - Install dependencies"
	@echo "  fmt            - Format code"
	@echo "  lint           - Lint code"
	@echo "  help           - Show this help"
