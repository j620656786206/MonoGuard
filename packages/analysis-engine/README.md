# @monoguard/analysis-engine

Go WASM analysis engine for MonoGuard - analyzes monorepo dependency graphs in the browser.

## Directory Structure

```
packages/analysis-engine/
├── cmd/
│   └── wasm/
│       └── main.go              # WASM entry point with js.FuncOf exports
├── pkg/
│   ├── analyzer/                # Dependency graph analysis (Epic 2)
│   ├── parser/                  # Workspace parsing (Epic 2)
│   └── types/                   # Shared Go types matching TypeScript
├── internal/
│   ├── handlers/                # Business logic for WASM exports (testable)
│   └── result/                  # Result<T> implementation
├── test/
│   └── smoke-test.html          # Browser smoke test
├── dist/                        # Build output (gitignored)
│   ├── monoguard.wasm
│   └── wasm_exec.js
├── go.mod
├── Makefile
└── package.json
```

## Building

```bash
# Build WASM module
make build-wasm

# Run Go tests
make test

# Clean build artifacts
make clean
```

## Usage in Browser

```javascript
// Load WASM
const go = new Go();
const result = await WebAssembly.instantiateStreaming(
  fetch('monoguard.wasm'),
  go.importObject
);
go.run(result.instance);

// Call exported functions
const versionResult = MonoGuard.getVersion();
const analysisResult = MonoGuard.analyze(JSON.stringify(workspaceData));
```

## Exported Functions

All functions return a Result<T> JSON structure:

```json
{
  "data": { ... } | null,
  "error": { "code": "ERROR_CODE", "message": "..." } | null
}
```

- `MonoGuard.getVersion()` - Returns version information
- `MonoGuard.analyze(jsonInput)` - Analyzes workspace dependencies
- `MonoGuard.check(jsonInput)` - Validates workspace configuration

## Development

Run the browser smoke test:

```bash
cd packages/analysis-engine
make build-wasm
npx serve .
# Open http://localhost:3000/test/smoke-test.html
```
