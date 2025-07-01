.PHONY: build run clean test help

# Default target
all: build

# Build the executable
build:
	@echo "Building randomfs-web..."
	go build -o randomfs-web main.go
	@echo "Build complete! Executable: ./randomfs-web"

# Run the server
run: build
	@echo "Starting RandomFS Web Server..."
	./randomfs-web

# Run with custom port
run-port: build
	@echo "Starting RandomFS Web Server on port $(PORT)..."
	RANDOMFS_PORT=$(PORT) ./randomfs-web

# Run without IPFS (for testing)
run-no-ipfs: build
	@echo "Starting RandomFS Web Server without IPFS..."
	RANDOMFS_IPFS_API="" ./randomfs-web

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f randomfs-web
	@echo "Clean complete!"

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Show help
help:
	@echo "Available targets:"
	@echo "  build      - Build the executable"
	@echo "  run        - Build and run the server (port 8080)"
	@echo "  run-port   - Run with custom port (make run-port PORT=3000)"
	@echo "  run-no-ipfs- Run without IPFS requirement"
	@echo "  clean      - Remove build artifacts"
	@echo "  test       - Run tests"
	@echo "  deps       - Install dependencies"
	@echo "  help       - Show this help message"
	@echo ""
	@echo "Environment variables:"
	@echo "  RANDOMFS_PORT     - Server port (default: 8080)"
	@echo "  RANDOMFS_DATA_DIR - Data directory (default: ./data)"
	@echo "  RANDOMFS_IPFS_API - IPFS API endpoint (default: http://localhost:5001)" 