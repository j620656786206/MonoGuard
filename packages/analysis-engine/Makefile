.PHONY: build-wasm clean test copy-wasm-exec dev

DIST_DIR := dist
WASM_OUTPUT := $(DIST_DIR)/monoguard.wasm
GO_ROOT := $(shell go env GOROOT)

# Build WASM module
build-wasm: clean $(DIST_DIR) copy-wasm-exec
	GOOS=js GOARCH=wasm go build -o $(WASM_OUTPUT) ./cmd/wasm/main.go
	@echo "Built $(WASM_OUTPUT)"
	@ls -lh $(WASM_OUTPUT)

# Create dist directory
$(DIST_DIR):
	mkdir -p $(DIST_DIR)

# Copy wasm_exec.js from Go installation
copy-wasm-exec: $(DIST_DIR)
	cp "$(GO_ROOT)/lib/wasm/wasm_exec.js" $(DIST_DIR)/

# Run Go tests (excludes WASM-specific code)
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf $(DIST_DIR)

# Development: build and show size
dev: build-wasm
	@echo "\nWASM Size Analysis:"
	@du -h $(WASM_OUTPUT)
