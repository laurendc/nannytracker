# NannyTracker Makefile
# Common tasks for development and release management

.PHONY: help build test clean release version

# Default target
help:
	@echo "Available targets:"
	@echo "  build     - Build the application for current platform"
	@echo "  build-all - Build for all supported platforms"
	@echo "  test      - Run all tests"
	@echo "  test-race - Run tests with race detection"
	@echo "  clean     - Clean build artifacts"
	@echo "  release   - Create a new release (requires VERSION=)"
	@echo "  version   - Show current version information"
	@echo "  lint      - Run linter"
	@echo "  fmt       - Format code"

# Build for current platform
build:
	go build -o nannytracker ./cmd/tui
	go build -o nannytracker-web ./cmd/web

# Build for all supported platforms
build-all:
	@echo "Building for all platforms..."
	mkdir -p dist
	
	# TUI Application
	GOOS=linux GOARCH=amd64 go build -o dist/nannytracker-linux-amd64 ./cmd/tui
	GOOS=linux GOARCH=arm64 go build -o dist/nannytracker-linux-arm64 ./cmd/tui
	GOOS=darwin GOARCH=amd64 go build -o dist/nannytracker-darwin-amd64 ./cmd/tui
	GOOS=darwin GOARCH=arm64 go build -o dist/nannytracker-darwin-arm64 ./cmd/tui
	GOOS=windows GOARCH=amd64 go build -o dist/nannytracker-windows-amd64.exe ./cmd/tui
	
	# Web Server
	GOOS=linux GOARCH=amd64 go build -o dist/nannytracker-web-linux-amd64 ./cmd/web
	GOOS=linux GOARCH=arm64 go build -o dist/nannytracker-web-linux-arm64 ./cmd/web
	GOOS=darwin GOARCH=amd64 go build -o dist/nannytracker-web-darwin-amd64 ./cmd/web
	GOOS=darwin GOARCH=arm64 go build -o dist/nannytracker-web-darwin-arm64 ./cmd/web
	GOOS=windows GOARCH=amd64 go build -o dist/nannytracker-web-windows-amd64.exe ./cmd/web
	
	@echo "Build complete. Binaries are in the dist/ directory."

# Run tests
test:
	go test ./... -v

# Run tests with race detection
test-race:
	go test -race ./...

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Clean build artifacts
clean:
	rm -rf dist/
	rm -f nannytracker
	rm -f nannytracker-web
	go clean -cache

# Create a new release
release:
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Usage: make release VERSION=v1.0.0"; \
		exit 1; \
	fi
	@echo "Creating release $(VERSION)..."
	
	# Check if tag already exists
	@if git tag -l | grep -q "$(VERSION)"; then \
		echo "Error: Tag $(VERSION) already exists"; \
		exit 1; \
	fi
	
	# Run tests
	@echo "Running tests..."
	@make test
	
	# Build for all platforms
	@echo "Building for all platforms..."
	@make build-all
	
	# Create git tag
	@echo "Creating git tag..."
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)
	
	@echo "Release $(VERSION) created successfully!"
	@echo "GitHub Actions will automatically create a release with binaries."

# Show version information
version:
	@echo "Current version information:"
	@go run -ldflags="-X github.com/laurendc/nannytracker/pkg/version.Version=dev -X github.com/laurendc/nannytracker/pkg/version.BuildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ) -X github.com/laurendc/nannytracker/pkg/version.GitCommit=$(shell git rev-parse --short HEAD)" ./cmd/tui --version 2>/dev/null || echo "Version information not available"

# Run linter
lint:
	golangci-lint run ./...

# Format code
fmt:
	go fmt ./...

# Install development dependencies
deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest

# Security scan
security:
	gosec ./...

# Check for outdated dependencies
deps-check:
	go list -u -m all

# Update dependencies
deps-update:
	go get -u ./...
	go mod tidy 