# Makefile for daemon-control

# Variables
BINARY_NAME=daemon-control
DIST_DIR=dist
GO_FILES=$(shell find . -name '*.go' -type f)
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Platforms
PLATFORMS=darwin/amd64 darwin/arm64 linux/amd64 linux/arm64

.PHONY: all build clean test fmt vet lint deps run install help

# Default target
all: clean build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(DIST_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME) .
	@echo "Build complete: $(DIST_DIR)/$(BINARY_NAME)"

# Build for multiple platforms
build-all:
	@echo "Building for all platforms..."
	@mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		GOOS=$$(echo $$platform | cut -d/ -f1) \
		GOARCH=$$(echo $$platform | cut -d/ -f2) \
		output_name=$(DIST_DIR)/$(BINARY_NAME)-$$(echo $$platform | tr '/' '-'); \
		echo "Building $$output_name..."; \
		GOOS=$$(echo $$platform | cut -d/ -f1) GOARCH=$$(echo $$platform | cut -d/ -f2) \
			$(GOBUILD) $(LDFLAGS) -o $$output_name .; \
	done
	@echo "Multi-platform build complete"

# Build for current platform with debug info
build-debug:
	@echo "Building $(BINARY_NAME) with debug info..."
	@mkdir -p $(DIST_DIR)
	$(GOBUILD) -gcflags="all=-N -l" -o $(DIST_DIR)/$(BINARY_NAME)-debug .
	@echo "Debug build complete: $(DIST_DIR)/$(BINARY_NAME)-debug"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(DIST_DIR)
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@mkdir -p $(DIST_DIR)
	$(GOTEST) -v -coverprofile=$(DIST_DIR)/coverage.out ./...
	$(GOCMD) tool cover -html=$(DIST_DIR)/coverage.out -o $(DIST_DIR)/coverage.html
	@echo "Coverage report: $(DIST_DIR)/coverage.html"

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...
	@echo "Format complete"

# Run go vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...
	@echo "Vet complete"

# Run linters (requires golangci-lint)
lint:
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: brew install golangci-lint"; \
	fi

# Update dependencies
deps:
	@echo "Updating dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Dependencies updated"

# Verify dependencies
deps-verify:
	@echo "Verifying dependencies..."
	$(GOMOD) verify
	@echo "Dependencies verified"

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	$(DIST_DIR)/$(BINARY_NAME)

# Install the binary to /usr/local/bin
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(DIST_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installation complete"

# Uninstall the binary from /usr/local/bin
uninstall:
	@echo "Uninstalling $(BINARY_NAME) from /usr/local/bin..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "Uninstall complete"

# Development mode - build and run with live reload (requires air)
dev:
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "air not installed. Install with: go install github.com/cosmtrek/air@latest"; \
	fi

# Create a release tarball
release: build-all
	@echo "Creating release tarballs..."
	@mkdir -p $(DIST_DIR)/releases
	@for platform in $(PLATFORMS); do \
		platform_name=$$(echo $$platform | tr '/' '-'); \
		tar -czf $(DIST_DIR)/releases/$(BINARY_NAME)-$(VERSION)-$$platform_name.tar.gz \
			-C $(DIST_DIR) $(BINARY_NAME)-$$platform_name; \
	done
	@echo "Release tarballs created in $(DIST_DIR)/releases/"

# Generate plist files from config
generate: build
	@echo "Generating plist files..."
	$(DIST_DIR)/$(BINARY_NAME) generate
	@echo "Plist files generated"

# Show help
help:
	@echo "Available targets:"
	@echo "  make              - Clean and build for current platform"
	@echo "  make build        - Build binary for current platform"
	@echo "  make build-all    - Build binaries for all platforms"
	@echo "  make build-debug  - Build with debug symbols"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make test         - Run tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make fmt          - Format code"
	@echo "  make vet          - Run go vet"
	@echo "  make lint         - Run linters (requires golangci-lint)"
	@echo "  make deps         - Update dependencies"
	@echo "  make deps-verify  - Verify dependencies"
	@echo "  make run          - Build and run the application"
	@echo "  make generate     - Generate plist files from config"
	@echo "  make install      - Install binary to /usr/local/bin"
	@echo "  make uninstall    - Remove binary from /usr/local/bin"
	@echo "  make dev          - Run in development mode (requires air)"
	@echo "  make release      - Create release tarballs for all platforms"
	@echo "  make help         - Show this help message"

# Ensure dist directory exists for any target that needs it
$(DIST_DIR):
	@mkdir -p $(DIST_DIR)