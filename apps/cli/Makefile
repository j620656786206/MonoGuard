.PHONY: build build-all clean test install dev

BINARY_NAME := monoguard
DIST_DIR := dist
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "0.1.0")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS := -ldflags "-X github.com/j620656786206/MonoGuard/apps/cli/cmd.version=$(VERSION) \
	-X github.com/j620656786206/MonoGuard/apps/cli/cmd.commit=$(COMMIT) \
	-X github.com/j620656786206/MonoGuard/apps/cli/cmd.buildDate=$(BUILD_DATE)"

# Build for current platform
build: $(DIST_DIR)
	go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME) .
	@echo "Built $(DIST_DIR)/$(BINARY_NAME)"
	@ls -lh $(DIST_DIR)/$(BINARY_NAME)

# Build for all platforms
build-all: $(DIST_DIR)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 .
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "Built binaries for all platforms"

$(DIST_DIR):
	mkdir -p $(DIST_DIR)

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf $(DIST_DIR)

# Install to GOPATH/bin
install:
	go install $(LDFLAGS) .

# Development: build and run help
dev: build
	$(DIST_DIR)/$(BINARY_NAME) --help
