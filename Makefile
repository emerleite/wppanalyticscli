# Variables
BINARY_NAME=wppanalyticscli
VERSION?=dev
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}"

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	@echo "Building ${BINARY_NAME}..."
	go build ${LDFLAGS} -o ${BINARY_NAME} .

# Build for specific platforms
.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY_NAME}-linux-amd64 .

.PHONY: build-windows
build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY_NAME}-windows-amd64.exe .

.PHONY: build-darwin
build-darwin:
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY_NAME}-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ${BINARY_NAME}-darwin-arm64 .

# Build all platforms
.PHONY: build-all
build-all: build-linux build-windows build-darwin

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -f ${BINARY_NAME}
	rm -f ${BINARY_NAME}-*

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping lint"; \
	fi

# Tidy dependencies
.PHONY: tidy
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# Install binary to GOPATH/bin
.PHONY: install
install:
	@echo "Installing ${BINARY_NAME}..."
	go install ${LDFLAGS} .

# Development build with race detection
.PHONY: dev
dev:
	@echo "Building development version..."
	go build -race ${LDFLAGS} -o ${BINARY_NAME}-dev .

# Run the binary (requires FB_ACCESS_TOKEN to be set)
.PHONY: run
run: build
	@echo "Running ${BINARY_NAME}..."
	@if [ -z "$$FB_ACCESS_TOKEN" ]; then \
		echo "Error: FB_ACCESS_TOKEN environment variable is required"; \
		exit 1; \
	fi
	./${BINARY_NAME} $(ARGS)

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary for current platform"
	@echo "  build-linux   - Build for Linux (amd64)"
	@echo "  build-windows - Build for Windows (amd64)"
	@echo "  build-darwin  - Build for macOS (amd64 and arm64)"
	@echo "  build-all     - Build for all platforms"
	@echo "  test          - Run tests"
	@echo "  clean         - Remove build artifacts"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code (requires golangci-lint)"
	@echo "  tidy          - Tidy dependencies"
	@echo "  install       - Install binary to GOPATH/bin"
	@echo "  dev           - Build development version with race detection"
	@echo "  run           - Build and run (set ARGS for arguments)"
	@echo "  help          - Show this help"
	@echo ""
	@echo "Examples:"
	@echo "  make build"
	@echo "  make run ARGS='-wbaid=123 -start=2025-06-20T00:00:00Z -end=2025-06-24T00:00:00Z'"