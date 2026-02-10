.PHONY: build install clean test lint run help

# Variables
BINARY_NAME=shadow
INSTALL_PATH=/usr/local/bin
GO=go
GOFLAGS=-v

# Build the binary
build:
	@echo "🔨 Building Shadow..."
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) cmd/shadow/main.go
	@echo "✅ Build complete: ./$(BINARY_NAME)"

# Install to system
install: build
	@echo "📦 Installing Shadow to $(INSTALL_PATH)..."
	sudo mv $(BINARY_NAME) $(INSTALL_PATH)/
	@echo "✅ Installed: $(INSTALL_PATH)/$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -f $(BINARY_NAME)
	$(GO) clean
	@echo "✅ Clean complete"

# Run tests
test:
	@echo "🧪 Running tests..."
	$(GO) test -v ./...

# Run linter
lint:
	@echo "🔍 Running linter..."
	golangci-lint run

# Download dependencies
deps:
	@echo "📥 Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# Run the application
run: build
	./$(BINARY_NAME)

# Build for multiple platforms
build-all:
	@echo "🔨 Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BINARY_NAME)-linux-amd64 cmd/shadow/main.go
	GOOS=darwin GOARCH=amd64 $(GO) build -o $(BINARY_NAME)-darwin-amd64 cmd/shadow/main.go
	GOOS=darwin GOARCH=arm64 $(GO) build -o $(BINARY_NAME)-darwin-arm64 cmd/shadow/main.go
	GOOS=windows GOARCH=amd64 $(GO) build -o $(BINARY_NAME)-windows-amd64.exe cmd/shadow/main.go
	@echo "✅ Multi-platform build complete"

# Development mode - build and run with example
dev: build
	./$(BINARY_NAME) --help

# Help
help:
	@echo "Shadow - Makefile Commands"
	@echo ""
	@echo "  make build       - Build the binary"
	@echo "  make install     - Install to /usr/local/bin"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make test        - Run tests"
	@echo "  make lint        - Run linter"
	@echo "  make deps        - Download dependencies"
	@echo "  make run         - Build and run"
	@echo "  make build-all   - Build for all platforms"
	@echo "  make dev         - Development mode"
	@echo "  make help        - Show this help"
