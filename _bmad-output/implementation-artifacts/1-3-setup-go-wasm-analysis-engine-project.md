# Story 1.3: Setup Go WASM Analysis Engine Project

Status: done

## Story

As a **developer**,
I want **a Go project structure configured for WASM compilation**,
So that **I can build the analysis engine that runs in the browser with zero backend dependency**.

## Acceptance Criteria

1. **AC1: Go Module Initialization**
   - Given the Nx monorepo from Story 1.1
   - When I initialize the Go WASM project in packages/analysis-engine
   - Then I have:
     - `go.mod` with module path `github.com/j620656786206/MonoGuard/packages/analysis-engine`
     - Go version 1.21+ specified in go.mod
     - No external dependencies (pure Go standard library for MVP)

2. **AC2: Project Directory Structure**
   - Given the Go module is initialized
   - When I verify the project structure
   - Then I have:
     ```
     packages/analysis-engine/
     ├── cmd/
     │   └── wasm/
     │       └── main.go         # WASM entry point
     ├── pkg/
     │   ├── analyzer/           # Core analysis logic (placeholder)
     │   ├── parser/             # Workspace parser (placeholder)
     │   └── types/              # Go types matching TypeScript
     ├── internal/
     │   └── result/             # Result<T> type implementation
     ├── dist/                   # Build output (gitignored)
     ├── go.mod
     ├── go.sum
     └── Makefile
     ```

3. **AC3: WASM Build Configuration**
   - Given the project structure
   - When I run `make build-wasm`
   - Then:
     - Build uses `GOOS=js GOARCH=wasm`
     - Output produces `dist/monoguard.wasm` file
     - `wasm_exec.js` is copied from Go installation to `dist/`
     - WASM file size is reasonable (< 5MB uncompressed for MVP)

4. **AC4: Basic WASM Exports**
   - Given the WASM build succeeds
   - When I examine the exported functions
   - Then the following functions are exported to JavaScript:
     - `MonoGuard.getVersion()` - Returns version string
     - `MonoGuard.analyze(jsonInput)` - Placeholder analysis function
     - `MonoGuard.check(jsonInput)` - Placeholder check function
   - And all functions return `Result<T>` JSON structure

5. **AC5: Result Type Implementation**
   - Given the WASM exports
   - When functions return data
   - Then they follow the unified Result pattern:
     ```json
     {
       "data": { ... } | null,
       "error": { "code": "ERROR_CODE", "message": "..." } | null
     }
     ```
   - And error codes use UPPER_SNAKE_CASE
   - And all JSON uses camelCase field names

6. **AC6: Browser Smoke Test**
   - Given the compiled WASM file
   - When I load it in a browser environment
   - Then:
     - WASM initializes without errors
     - `MonoGuard.getVersion()` returns expected version
     - `MonoGuard.analyze("{}")` returns valid Result JSON
     - Console shows no errors

7. **AC7: Nx Integration**
   - Given the Go project
   - When I run Nx commands
   - Then:
     - `pnpm nx build analysis-engine` runs the Makefile
     - Build output goes to `packages/analysis-engine/dist/`
     - Project appears correctly in `pnpm nx graph`

## Tasks / Subtasks

- [x] **Task 1: Initialize Go Module** (AC: #1)
  - [x] 1.1 Remove existing TypeScript placeholder files:
    ```bash
    rm -rf packages/analysis-engine/src
    rm packages/analysis-engine/tsconfig.json
    ```
  - [x] 1.2 Initialize Go module:
    ```bash
    cd packages/analysis-engine
    go mod init github.com/j620656786206/MonoGuard/packages/analysis-engine
    ```
  - [x] 1.3 Verify Go version is 1.21+ in go.mod

- [x] **Task 2: Create Directory Structure** (AC: #2)
  - [x] 2.1 Create cmd/wasm directory for WASM entry point
  - [x] 2.2 Create pkg directories for future analysis code:
    ```bash
    mkdir -p cmd/wasm pkg/{analyzer,parser,types} internal/result
    ```
  - [x] 2.3 Create .gitignore for dist/ directory
  - [x] 2.4 Add placeholder README.md explaining Go structure

- [x] **Task 3: Implement Result Type** (AC: #5)
  - [x] 3.1 Create `internal/result/result.go`:

    ```go
    package result

    import "encoding/json"

    // Result represents the unified response type for all WASM functions
    // This matches the TypeScript Result<T> type
    type Result struct {
        Data  interface{} `json:"data"`
        Error *Error      `json:"error"`
    }

    // Error represents an error response
    type Error struct {
        Code    string `json:"code"`
        Message string `json:"message"`
    }

    // NewSuccess creates a successful Result
    func NewSuccess(data interface{}) *Result {
        return &Result{Data: data, Error: nil}
    }

    // NewError creates an error Result
    func NewError(code, message string) *Result {
        return &Result{Data: nil, Error: &Error{Code: code, Message: message}}
    }

    // ToJSON serializes Result to JSON string
    func (r *Result) ToJSON() string {
        b, _ := json.Marshal(r)
        return string(b)
    }
    ```

  - [x] 3.2 Create `internal/result/result_test.go` with unit tests

- [x] **Task 4: Create WASM Entry Point** (AC: #4)
  - [x] 4.1 Create `cmd/wasm/main.go`:

    ```go
    //go:build js && wasm

    package main

    import (
        "syscall/js"

        "github.com/j620656786206/MonoGuard/packages/analysis-engine/internal/result"
    )

    const version = "0.1.0"

    func getVersion(this js.Value, args []js.Value) interface{} {
        r := result.NewSuccess(map[string]string{"version": version})
        return r.ToJSON()
    }

    func analyze(this js.Value, args []js.Value) interface{} {
        if len(args) < 1 {
            r := result.NewError("INVALID_INPUT", "Missing JSON input")
            return r.ToJSON()
        }

        // Placeholder - will be implemented in Epic 2
        input := args[0].String()
        _ = input // Suppress unused warning

        r := result.NewSuccess(map[string]interface{}{
            "healthScore": 100,
            "packages":    0,
            "placeholder": true,
        })
        return r.ToJSON()
    }

    func check(this js.Value, args []js.Value) interface{} {
        if len(args) < 1 {
            r := result.NewError("INVALID_INPUT", "Missing JSON input")
            return r.ToJSON()
        }

        // Placeholder - will be implemented in Epic 2
        r := result.NewSuccess(map[string]interface{}{
            "passed":      true,
            "errors":      []string{},
            "placeholder": true,
        })
        return r.ToJSON()
    }

    func main() {
        // Create MonoGuard namespace
        monoguard := make(map[string]interface{})
        monoguard["getVersion"] = js.FuncOf(getVersion)
        monoguard["analyze"] = js.FuncOf(analyze)
        monoguard["check"] = js.FuncOf(check)

        js.Global().Set("MonoGuard", monoguard)

        // Keep the Go program running
        <-make(chan bool)
    }
    ```

- [x] **Task 5: Create Makefile** (AC: #3)
  - [x] 5.1 Create `Makefile`:

    ```makefile
    .PHONY: build-wasm clean test copy-wasm-exec

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
    	cp "$(GO_ROOT)/misc/wasm/wasm_exec.js" $(DIST_DIR)/

    # Run Go tests
    test:
    	go test -v ./...

    # Clean build artifacts
    clean:
    	rm -rf $(DIST_DIR)

    # Development: build and show size
    dev: build-wasm
    	@echo "\nWASM Size Analysis:"
    	@du -h $(WASM_OUTPUT)
    ```

- [x] **Task 6: Update Package Configuration** (AC: #7)
  - [x] 6.1 Update `packages/analysis-engine/package.json`:
    ```json
    {
      "name": "@monoguard/analysis-engine",
      "version": "0.1.0",
      "private": true,
      "description": "Go WASM analysis engine for MonoGuard",
      "type": "module",
      "main": "./dist/wasm_exec.js",
      "files": ["dist/"],
      "scripts": {
        "build": "make build-wasm",
        "test": "make test",
        "clean": "make clean"
      }
    }
    ```
  - [x] 6.2 Create or update `packages/analysis-engine/project.json`:
    ```json
    {
      "name": "@monoguard/analysis-engine",
      "projectType": "library",
      "sourceRoot": "packages/analysis-engine",
      "targets": {
        "build": {
          "executor": "nx:run-commands",
          "options": {
            "command": "make build-wasm",
            "cwd": "packages/analysis-engine"
          },
          "outputs": ["{projectRoot}/dist"]
        },
        "test": {
          "executor": "nx:run-commands",
          "options": {
            "command": "make test",
            "cwd": "packages/analysis-engine"
          }
        },
        "clean": {
          "executor": "nx:run-commands",
          "options": {
            "command": "make clean",
            "cwd": "packages/analysis-engine"
          }
        }
      }
    }
    ```

- [x] **Task 7: Create Browser Smoke Test** (AC: #6)
  - [x] 7.1 Create `packages/analysis-engine/test/smoke-test.html`:

    ```html
    <!DOCTYPE html>
    <html>
      <head>
        <title>MonoGuard WASM Smoke Test</title>
      </head>
      <body>
        <h1>MonoGuard WASM Smoke Test</h1>
        <pre id="output"></pre>
        <script src="../dist/wasm_exec.js"></script>
        <script>
          const output = document.getElementById('output');
          function log(msg) {
            output.textContent += msg + '\n';
            console.log(msg);
          }

          async function runTests() {
            log('Loading WASM...');
            const go = new Go();
            const result = await WebAssembly.instantiateStreaming(
              fetch('../dist/monoguard.wasm'),
              go.importObject
            );
            go.run(result.instance);
            log('WASM loaded successfully!\n');

            // Test getVersion
            log('Testing MonoGuard.getVersion():');
            const version = MonoGuard.getVersion();
            log(version);

            // Test analyze
            log('\nTesting MonoGuard.analyze("{}"):');
            const analyzeResult = MonoGuard.analyze('{}');
            log(analyzeResult);

            // Test check
            log('\nTesting MonoGuard.check("{}"):');
            const checkResult = MonoGuard.check('{}');
            log(checkResult);

            log('\n✅ All smoke tests passed!');
          }

          runTests().catch((err) => {
            log('❌ Error: ' + err.message);
            console.error(err);
          });
        </script>
      </body>
    </html>
    ```

  - [x] 7.2 Add instructions for running smoke test in README

- [x] **Task 8: Create pkg Placeholder Files** (AC: #2)
  - [x] 8.1 Create `pkg/analyzer/analyzer.go`:

    ```go
    // Package analyzer provides dependency graph analysis
    // This package will be implemented in Epic 2
    package analyzer

    // Placeholder for Epic 2 implementation
    ```

  - [x] 8.2 Create `pkg/parser/parser.go`:

    ```go
    // Package parser provides workspace configuration parsing
    // Supports npm, yarn, and pnpm workspaces
    // This package will be implemented in Epic 2
    package parser

    // Placeholder for Epic 2 implementation
    ```

  - [x] 8.3 Create `pkg/types/types.go`:

    ```go
    // Package types defines Go types that match TypeScript definitions
    // All JSON tags use camelCase for cross-language consistency
    package types

    // AnalysisResult represents the complete analysis output
    // This matches @monoguard/types AnalysisResult
    type AnalysisResult struct {
        HealthScore int    `json:"healthScore"`
        Packages    int    `json:"packages"`
        CreatedAt   string `json:"createdAt"` // ISO 8601 format
    }

    // CheckResult represents validation-only output
    type CheckResult struct {
        Passed bool     `json:"passed"`
        Errors []string `json:"errors"`
    }
    ```

- [x] **Task 9: Verification** (AC: #3, #6, #7)
  - [x] 9.1 Run `make build-wasm` - verify WASM builds successfully
  - [x] 9.2 Check WASM file size: `du -h dist/monoguard.wasm`
  - [x] 9.3 Run `make test` - verify Go tests pass
  - [x] 9.4 Serve smoke test and verify in browser:
    ```bash
    cd packages/analysis-engine
    npx serve .
    # Open http://localhost:3000/test/smoke-test.html
    ```
  - [x] 9.5 Run `pnpm nx build analysis-engine` - verify Nx integration
  - [x] 9.6 Run `pnpm nx graph` - verify project appears correctly

## Dev Notes

### Architecture Patterns & Constraints

**From Architecture Document:**

- **Technology:** Go 1.21+ compiled to WASM (`GOOS=js GOARCH=wasm`)
- **Entry Point:** `cmd/wasm/main.go` with `syscall/js` for JavaScript interop
- **Build Tool:** Makefile with standard Go toolchain
- **Output:** `monoguard.wasm` + `wasm_exec.js` (Go runtime for WASM)

**Critical Constraints:**

- **Zero External Dependencies:** Pure Go standard library only for MVP
- **Result Type Mandatory:** ALL exported functions MUST return `Result<T>` JSON
- **camelCase JSON:** All struct tags MUST use camelCase (NOT snake_case)
- **UPPER_SNAKE_CASE Errors:** Error codes MUST be UPPER_SNAKE_CASE

### Critical Don't-Miss Rules

**From project-context.md:**

1. **Go Naming Conventions:**
   - PascalCase for exported functions/types (e.g., `AnalyzeWorkspace`)
   - camelCase for unexported functions/types (e.g., `parseInput`)
   - snake_case for file names (e.g., `workspace_parser.go`)
   - Test files: `*_test.go` in same directory

2. **JSON Serialization:**

   ```go
   // ✅ CORRECT: camelCase JSON tags
   type AnalysisResult struct {
       HealthScore int `json:"healthScore"`
       CreatedAt   string `json:"createdAt"`
   }

   // ❌ WRONG: snake_case JSON tags
   type AnalysisResult struct {
       HealthScore int `json:"health_score"` // BREAKS TypeScript
   }
   ```

3. **Result Type Pattern:**

   ```go
   // ✅ CORRECT: Always wrap returns in Result
   func AnalyzeWorkspace(input string) string {
       if err != nil {
           return result.NewError("PARSE_ERROR", err.Error()).ToJSON()
       }
       return result.NewSuccess(data).ToJSON()
   }

   // ❌ WRONG: Raw error returns
   func AnalyzeWorkspace(input string) string {
       if err != nil {
           return err.Error() // Frontend can't parse this!
       }
   }
   ```

4. **Date Format:**
   - Always ISO 8601 strings: `"2026-01-15T10:30:00Z"`
   - NEVER Unix timestamps

### Project Structure Notes

**Target Directory Structure:**

```
packages/analysis-engine/
├── cmd/
│   └── wasm/
│       └── main.go              # WASM entry point, js.FuncOf exports
├── pkg/
│   ├── analyzer/                # Dependency graph analysis (Epic 2)
│   │   ├── analyzer.go
│   │   └── analyzer_test.go
│   ├── parser/                  # Workspace parsing (Epic 2)
│   │   ├── parser.go
│   │   └── parser_test.go
│   └── types/                   # Shared Go types
│       └── types.go
├── internal/
│   └── result/                  # Result<T> implementation
│       ├── result.go
│       └── result_test.go
├── test/
│   └── smoke-test.html          # Browser smoke test
├── dist/                        # Build output (gitignored)
│   ├── monoguard.wasm
│   └── wasm_exec.js
├── go.mod
├── go.sum
├── Makefile
├── package.json                 # npm package config
├── project.json                 # Nx project config
└── README.md
```

**Integration with Web App (Future):**

- WASM file will be copied to `apps/web/public/` during build
- TypeScript adapter (Story 2.7) will wrap WASM calls with proper types
- Story 5.3 will implement actual WASM loading in browser

### WASM JavaScript Interop

**Exported Functions Pattern:**

```go
// Register function for JavaScript
js.Global().Set("MonoGuard", map[string]interface{}{
    "analyze": js.FuncOf(analyzeFunc),
})

// Function signature for js.FuncOf
func analyzeFunc(this js.Value, args []js.Value) interface{} {
    // args[0].String() to get string input
    // Return string (JSON) or js.Value
    return resultJSON
}
```

**Browser Usage:**

```javascript
// After loading WASM
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

### Testing Requirements

**Go Unit Tests:**

- Use standard `testing` package
- Table-driven tests for multiple scenarios
- Test files: `*_test.go` in same directory as source

**Browser Smoke Test:**

- Must verify WASM loads without errors
- Must verify all exported functions return valid JSON
- Must verify Result type structure is correct

### Previous Story Intelligence

**From Story 1.1 (done):**

- `packages/analysis-engine/` directory exists with TypeScript placeholder
- Path mapping `@monoguard/analysis-engine` configured in tsconfig.base.json
- Package named `@monoguard/analysis-engine` in package.json

**From Story 1.2 (ready-for-dev):**

- Web app will be TanStack Start (not Next.js)
- WASM files should go in `apps/web/public/` for serving
- Bundle size concern: WASM < 2MB compressed target

### Performance Targets

**From NFR Requirements:**

- Analyze 100 packages: < 5 seconds
- Analyze 1000 packages: < 30 seconds
- Memory usage: < 100MB in-browser
- WASM file size: < 2MB (compressed)

**This Story:** Focus on build infrastructure, not performance optimization.

### References

- [Source: _bmad-output/planning-artifacts/architecture.md#Analysis Engine (Go WASM)]
- [Source: _bmad-output/planning-artifacts/architecture.md#Starter Options Considered]
- [Source: _bmad-output/project-context.md#Go Naming Conventions]
- [Source: _bmad-output/project-context.md#Result Type Pattern]
- [Source: _bmad-output/planning-artifacts/epics.md#Story 1.3]
- [Go WASM Wiki](https://github.com/golang/go/wiki/WebAssembly)
- [syscall/js Documentation](https://pkg.go.dev/syscall/js)

## Dev Agent Record

### Agent Model Used

Claude Opus 4.5 (claude-opus-4-5-20251101)

### Debug Log References

N/A - No issues encountered during implementation.

### Completion Notes List

- ✅ **Task 1 Complete:** Go module initialized with Go 1.25.5 (exceeds 1.21+ requirement)
- ✅ **Task 2 Complete:** Directory structure created: cmd/wasm, pkg/{analyzer,parser,types}, internal/result, test/
- ✅ **Task 3 Complete:** Result<T> type implemented with 4 table-driven unit tests (all passing)
  - Follows red-green-refactor cycle: tests written first, then implementation
  - Supports camelCase JSON serialization and UPPER_SNAKE_CASE error codes
- ✅ **Task 4 Complete:** WASM entry point with MonoGuard.{getVersion, analyze, check} exports
- ✅ **Task 5 Complete:** Makefile with build-wasm, test, clean, dev targets
- ✅ **Task 6 Complete:** package.json and project.json updated for Nx integration
- ✅ **Task 7 Complete:** Browser smoke test HTML page with comprehensive test cases
- ✅ **Task 8 Complete:** Placeholder files created for pkg/{analyzer,parser,types}
- ✅ **Task 9 Complete:** All verifications passed
  - WASM builds successfully (2.8MB < 5MB target)
  - Go tests pass (4 tests, 100%)
  - Nx build integration works (`pnpm nx build @monoguard/analysis-engine`)
  - Project appears correctly in Nx graph

**Key Implementation Decisions:**

1. Used Go 1.25.5 `lib/wasm/wasm_exec.js` path (newer Go versions)
2. Result type includes fallback error JSON if marshal fails
3. Smoke test includes both success and error case testing
4. Types in pkg/types follow camelCase JSON convention per project-context.md

### File List

**New Files:**

- packages/analysis-engine/go.mod
- packages/analysis-engine/Makefile
- packages/analysis-engine/README.md
- packages/analysis-engine/.gitignore
- packages/analysis-engine/project.json
- packages/analysis-engine/cmd/wasm/main.go
- packages/analysis-engine/internal/result/result.go
- packages/analysis-engine/internal/result/result_test.go
- packages/analysis-engine/internal/handlers/handlers.go (added in code review)
- packages/analysis-engine/internal/handlers/handlers_test.go (added in code review)
- packages/analysis-engine/pkg/analyzer/analyzer.go
- packages/analysis-engine/pkg/analyzer/analyzer_test.go (added in code review)
- packages/analysis-engine/pkg/parser/parser.go
- packages/analysis-engine/pkg/parser/parser_test.go (added in code review)
- packages/analysis-engine/pkg/types/types.go
- packages/analysis-engine/pkg/types/types_test.go (added in code review)
- packages/analysis-engine/test/smoke-test.html

**Modified Files:**

- packages/analysis-engine/package.json
- packages/analysis-engine/cmd/wasm/main.go (refactored to use handlers package)
- packages/analysis-engine/internal/result/result.go (improved error handling, added control char escaping)
- packages/analysis-engine/internal/result/result_test.go (added escapeJSON and fallback tests, coverage 88.5%)
- packages/analysis-engine/pkg/types/types.go (added VersionInfo, Placeholder fields)
- packages/analysis-engine/README.md (added handlers directory to structure)

**Removed Files:**

- packages/analysis-engine/src/ (TypeScript placeholder directory)
- packages/analysis-engine/tsconfig.json

**Build Outputs (gitignored):**

- packages/analysis-engine/dist/monoguard.wasm (2.8MB)
- packages/analysis-engine/dist/wasm_exec.js (17KB)

### Senior Developer Review (AI)

**Reviewer:** Amelia (Dev Agent) | **Date:** 2026-01-16 | **Outcome:** APPROVED with fixes applied

**Issues Found & Fixed:**

| ID     | Severity | Issue                                                   | Resolution                                                                            |
| ------ | -------- | ------------------------------------------------------- | ------------------------------------------------------------------------------------- |
| HIGH-1 | HIGH     | pkg/types unused - inline maps instead of typed structs | Refactored main.go to use types.VersionInfo, types.AnalysisResult, types.CheckResult  |
| MED-1  | MEDIUM   | No tests for pkg/types JSON serialization               | Added types_test.go with 6 test functions verifying camelCase JSON                    |
| MED-2  | MEDIUM   | No tests for pkg/analyzer and pkg/parser                | Added placeholder tests documenting expected Epic 2 interface                         |
| MED-3  | MEDIUM   | WASM entry point logic untested                         | Created internal/handlers package with extracted testable logic + comprehensive tests |
| LOW-2  | LOW      | Result.ToJSON silently swallows marshal errors          | Updated to preserve original error message in fallback JSON                           |

**Test Results After Fixes:**

- Total tests: 25+ (all passing)
- Packages tested: 5/5
- WASM build: Success (2.8MB)

### Senior Developer Review #2 (AI)

**Reviewer:** Amelia (Dev Agent) | **Date:** 2026-01-16 | **Outcome:** APPROVED with fixes applied

**Issues Found & Fixed:**

| ID    | Severity | Issue                                          | Resolution                                                           |
| ----- | -------- | ---------------------------------------------- | -------------------------------------------------------------------- |
| MED-1 | MEDIUM   | result package test coverage only 27.8%        | Added TestEscapeJSON (10 cases) + TestToJSONWithUnmarshalableData    |
| MED-2 | MEDIUM   | escapeJSON missing control character handling  | Updated to handle all control chars (0x00-0x1F) with \uXXXX encoding |
| LOW-1 | LOW      | README missing handlers directory in structure | Updated README.md to include internal/handlers/                      |

**Test Coverage After Fixes:**

- handlers: 100.0%
- result: **27.8% → 88.5%** ✅
- WASM build: Success (2.8MB)

### Change Log

- 2026-01-16: Story 1.3 implemented - Go WASM analysis engine project setup complete
- 2026-01-16: Code review #1 completed - 1 HIGH, 4 MEDIUM issues fixed; refactored to handlers pattern for testability
- 2026-01-16: Code review #2 completed - 2 MEDIUM, 1 LOW issues fixed; result coverage improved to 88.5%
