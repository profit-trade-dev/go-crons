.PHONY: help build test clean deps run-example run-advanced-example

# Default target
help:
	@echo "Available targets:"
	@echo "  build              - Build the go-crons library"
	@echo "  test               - Run tests"
	@echo "  clean              - Clean build artifacts"
	@echo "  deps               - Download dependencies"
	@echo "  run-example        - Run basic usage example"
	@echo "  run-advanced-example - Run advanced usage example"
	@echo "  fmt                - Format code"
	@echo "  vet                - Run go vet"
	@echo "  lint               - Run golangci-lint"

# Build the library
build:
	@echo "Building go-crons library..."
	go build ./...

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	go clean ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Run basic example
run-example:
	@echo "Running basic usage example..."
	go run examples/basic_usage.go

# Run advanced example
run-advanced-example:
	@echo "Running advanced usage example..."
	go run examples/advanced_usage.go

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Run linter (if golangci-lint is installed)
lint:
	@echo "Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping..."; \
	fi

# Install dependencies for development
install-deps:
	@echo "Installing development dependencies..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Check if required tools are installed
check-tools:
	@echo "Checking required tools..."
	@command -v go >/dev/null 2>&1 || { echo "Go is not installed"; exit 1; }
	@echo "Go version: $$(go version)"
	@echo "All required tools are installed"

# Setup development environment
setup: check-tools deps install-deps
	@echo "Development environment setup complete"

# Run all checks
check: fmt vet test
	@echo "All checks passed"
