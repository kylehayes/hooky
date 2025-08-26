# Hooky Makefile

VERSION ?= 1.0.0
BUILD_DIR = dist
BINARY_NAME = hooky

# Default target
.PHONY: all
all: build

# Build for current platform
.PHONY: build
build:
	@echo "🏗️  Building $(BINARY_NAME) for current platform..."
	go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME) .
	@echo "✅ Build completed: ./$(BINARY_NAME)"

# Build for all platforms
.PHONY: build-all
build-all:
	@echo "🏗️  Building for all platforms..."
	./build.sh

# Clean build artifacts
.PHONY: clean
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)
	@echo "✅ Clean completed"

# Run tests
.PHONY: test
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# Install dependencies
.PHONY: deps
deps:
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy

# Format code
.PHONY: fmt
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	@echo "🔍 Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, running go vet instead..."; \
		go vet ./...; \
	fi

# Development workflow
.PHONY: dev
dev: deps fmt lint test build

# Release workflow
.PHONY: release
release: clean fmt lint test build-all
	@echo "🎉 Release build completed!"
	@echo "📦 Archives available in $(BUILD_DIR)/"

# Install locally (for development)
.PHONY: install
install: build
	@echo "📦 Installing $(BINARY_NAME) locally..."
	cp $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "✅ $(BINARY_NAME) installed to /usr/local/bin/"

# Uninstall local installation
.PHONY: uninstall
uninstall:
	@echo "🗑️  Uninstalling $(BINARY_NAME)..."
	rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✅ $(BINARY_NAME) uninstalled"

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build      - Build for current platform"
	@echo "  build-all  - Build for all platforms"
	@echo "  clean      - Clean build artifacts"
	@echo "  test       - Run tests"
	@echo "  deps       - Install dependencies"
	@echo "  fmt        - Format code"
	@echo "  lint       - Run linter"
	@echo "  dev        - Development workflow (deps + fmt + lint + test + build)"
	@echo "  release    - Release workflow (clean + fmt + lint + test + build-all)"
	@echo "  install    - Install locally for development"
	@echo "  uninstall  - Uninstall local installation"
	@echo "  help       - Show this help message"