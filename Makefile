# NannyTracker Makefile
# Single source of truth for build configuration and release management

.PHONY: help build test clean release version verify-release deps lint fmt security deps-check deps-update

# Build configuration - single source of truth
VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT := $(shell git rev-parse --short HEAD)

# Platform definitions - DRY principle
PLATFORMS := linux-amd64 linux-arm64 darwin-amd64 darwin-arm64 windows-amd64
BINARY_NAMES := nannytracker nannytracker-web

# Build flags
LDFLAGS := -X github.com/laurendc/nannytracker/pkg/version.Version=$(VERSION) \
           -X github.com/laurendc/nannytracker/pkg/version.BuildTime=$(BUILD_TIME) \
           -X github.com/laurendc/nannytracker/pkg/version.GitCommit=$(GIT_COMMIT)

# Default target
help:
	@echo "NannyTracker Build System"
	@echo "========================="
	@echo "Build Targets:"
	@echo "  build         - Build for current platform"
	@echo "  build-all     - Build for all supported platforms"
	@echo "  clean         - Clean build artifacts"
	@echo ""
	@echo "Test Targets:"
	@echo "  test          - Run all backend and frontend tests"
	@echo "  test-backend  - Run Go backend tests only"
	@echo "  test-frontend - Run frontend (web) tests only"
	@echo "  test-race     - Run backend tests with race detection"
	@echo "  test-coverage - Run backend tests with coverage"
	@echo ""
	@echo "Release Targets:"
	@echo "  release       - Create a new release (requires VERSION=)"
	@echo "  verify-release- Verify release artifacts (requires VERSION=)"
	@echo "  version       - Show current version information"
	@echo ""
	@echo "Development Targets:"
	@echo "  deps          - Install development dependencies"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
	@echo "  security      - Run security scan"
	@echo "  deps-check    - Check for outdated dependencies"
	@echo "  deps-update   - Update dependencies"

# Build for current platform
build:
	@echo "Building for current platform..."
	go build -ldflags="$(LDFLAGS)" -o nannytracker ./cmd/tui
	go build -ldflags="$(LDFLAGS)" -o nannytracker-web ./cmd/web
	@echo "Build complete."

# Build for all supported platforms
build-all:
	@echo "Building for all platforms..."
	@mkdir -p dist
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%-*} GOARCH=$${platform#*-} go build -ldflags="$(LDFLAGS)" \
			-o dist/nannytracker-$$platform$$(if [ "$${platform#*-}" = "amd64" ] && [ "$${platform%-*}" = "windows" ]; then echo .exe; fi) ./cmd/tui; \
		GOOS=$${platform%-*} GOARCH=$${platform#*-} go build -ldflags="$(LDFLAGS)" \
			-o dist/nannytracker-web-$$platform$$(if [ "$${platform#*-}" = "amd64" ] && [ "$${platform%-*}" = "windows" ]; then echo .exe; fi) ./cmd/web; \
	done
	@echo "Build complete. Binaries are in the dist/ directory."

# Run all tests (backend and frontend)
test: test-backend test-frontend

# Run Go backend tests only
test-backend:
	@echo "Running Go backend tests..."
	go test ./... -v

# Run frontend (web) tests only
test-frontend:
	@echo "Running frontend (web) tests..."
	cd web && npm install --no-audit --no-fund && npm test -- --watchAll=false

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf dist/
	rm -f nannytracker nannytracker-web
	go clean -cache
	@echo "Clean complete."

# Create a new release
release:
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Usage: make release VERSION=v1.0.0"; \
		exit 1; \
	fi
	@echo "Creating release $(VERSION)..."
	
	# Validate version format
	@if ! echo "$(VERSION)" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+'; then \
		echo "Error: VERSION must be in format vX.Y.Z (e.g., v1.0.0)"; \
		exit 1; \
	fi
	
	# Check if tag already exists
	@if git tag -l | grep -q "$(VERSION)"; then \
		echo "Error: Tag $(VERSION) already exists"; \
		exit 1; \
	fi
	
	# Check for uncommitted changes
	@if [ -n "$(shell git status --porcelain)" ]; then \
		echo "Error: Working directory is not clean. Please commit or stash changes."; \
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
	@echo "Version Information:"
	@echo "  Version: $(VERSION)"
	@echo "  Build Time: $(BUILD_TIME)"
	@echo "  Git Commit: $(GIT_COMMIT)"
	@echo ""
	@echo "Binary Version Output:"
	@if [ -f ./nannytracker ]; then \
		./nannytracker --version 2>/dev/null || echo "  TUI: Version information not available"; \
	else \
		echo "  TUI: Binary not built"; \
	fi
	@if [ -f ./nannytracker-web ]; then \
		./nannytracker-web --version 2>/dev/null || echo "  Web: Version information not available"; \
	else \
		echo "  Web: Binary not built"; \
	fi

# Verify release artifacts for a specific version
verify-release:
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Usage: make verify-release VERSION=v1.0.0"; \
		exit 1; \
	fi
	@echo "Verifying release artifacts for $(VERSION)..."
	@./scripts/verify-release.sh $(VERSION)

# Development dependencies
deps:
	@echo "Installing development dependencies..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	@echo "Dependencies installed."

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Security scan
security:
	@echo "Running security scan..."
	gosec ./...

# Check for outdated dependencies
deps-check:
	@echo "Checking for outdated dependencies..."
	go list -u -m all

# Update dependencies
deps-update:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy
	@echo "Dependencies updated." 