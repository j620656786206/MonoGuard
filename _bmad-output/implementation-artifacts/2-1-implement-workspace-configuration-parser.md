# Story 2.1: Implement Workspace Configuration Parser

Status: done

## Story

As a **developer**,
I want **the analysis engine to parse workspace configuration files**,
So that **I can extract package information from any supported monorepo format**.

## Acceptance Criteria

1. **AC1: npm/yarn Workspaces Parsing**
   - Given a monorepo with `package.json` containing `workspaces` field
   - When I provide the workspace files to the parser
   - Then the parser correctly extracts:
     - All package names from each package's package.json
     - Package paths relative to workspace root
     - Workspace root configuration
   - And supports both npm and yarn workspaces formats (they use the same syntax)
   - And handles glob patterns in workspaces field (e.g., `["packages/*", "apps/*"]`)

2. **AC2: pnpm Workspaces Parsing**
   - Given a monorepo with `pnpm-workspace.yaml`
   - When I provide the workspace files to the parser
   - Then the parser correctly extracts:
     - All package names and paths
     - Handles pnpm-specific `packages:` array format
   - And correctly interprets glob patterns with negation (e.g., `!packages/experimental-*`)

3. **AC3: Package Dependencies Extraction**
   - Given parsed workspace packages
   - When I analyze each package's package.json
   - Then the parser extracts for each package:
     - `dependencies` - production dependencies
     - `devDependencies` - development dependencies
     - `peerDependencies` - peer dependencies
   - And distinguishes between internal (workspace) and external dependencies
   - And captures version ranges for each dependency

4. **AC4: WorkspaceData Output Type**
   - Given successful parsing
   - When the parser returns results
   - Then it returns a structured `WorkspaceData` type containing:
     - `rootPath` - absolute path to workspace root
     - `workspaceType` - "npm" | "yarn" | "pnpm"
     - `packages` - map of package name to Package info
   - And all JSON serialization uses camelCase per project conventions
   - And matches the TypeScript `WorkspaceData` type definition

5. **AC5: Error Handling**
   - Given invalid or malformed workspace configuration
   - When parsing fails
   - Then the parser returns a descriptive error with:
     - Error code in UPPER_SNAKE_CASE (e.g., `INVALID_WORKSPACE`, `MISSING_PACKAGE_JSON`)
     - Human-readable error message
   - And uses the unified `Result<T>` pattern from `internal/result`

6. **AC6: Performance**
   - Given a workspace with 100 packages
   - When parsing completes
   - Then it finishes in < 1 second
   - And memory usage is reasonable (no excessive allocations)

## Tasks / Subtasks

- [x] **Task 1: Define WorkspaceData Types in Go** (AC: #4)
  - [x] 1.1 Update `pkg/types/types.go` to add WorkspaceData and related types:
    ```go
    // WorkspaceData represents the complete parsed workspace configuration.
    // Matches @monoguard/types WorkspaceData.
    type WorkspaceData struct {
        RootPath      string              `json:"rootPath"`
        WorkspaceType WorkspaceType       `json:"workspaceType"`
        Packages      map[string]*Package `json:"packages"`
    }

    // WorkspaceType identifies the package manager workspace format.
    type WorkspaceType string

    const (
        WorkspaceTypeNpm     WorkspaceType = "npm"
        WorkspaceTypeYarn    WorkspaceType = "yarn"
        WorkspaceTypePnpm    WorkspaceType = "pnpm"
        WorkspaceTypeUnknown WorkspaceType = "unknown"
    )

    // Package represents a single package in the workspace (expanded).
    type Package struct {
        Name             string            `json:"name"`
        Version          string            `json:"version"`
        Path             string            `json:"path"`
        Dependencies     map[string]string `json:"dependencies"`
        DevDependencies  map[string]string `json:"devDependencies"`
        PeerDependencies map[string]string `json:"peerDependencies"`
    }
    ```
  - [x] 1.2 Add corresponding tests in `pkg/types/types_test.go`
  - [x] 1.3 Verify JSON serialization produces camelCase output

- [x] **Task 2: Implement Package.json Parser** (AC: #1, #3)
  - [x] 2.1 Create `pkg/parser/package_json.go`:
    ```go
    // PackageJSON represents the structure of a package.json file
    type PackageJSON struct {
        Name             string            `json:"name"`
        Version          string            `json:"version"`
        Dependencies     map[string]string `json:"dependencies"`
        DevDependencies  map[string]string `json:"devDependencies"`
        PeerDependencies map[string]string `json:"peerDependencies"`
        Workspaces       interface{}       `json:"workspaces"` // Can be []string or WorkspacesConfig
    }

    // WorkspacesConfig for extended workspaces format
    type WorkspacesConfig struct {
        Packages []string `json:"packages"`
        Nohoist  []string `json:"nohoist"`
    }

    // ParsePackageJSON parses a single package.json file
    func ParsePackageJSON(data []byte) (*PackageJSON, error)

    // ExtractWorkspacePatterns extracts workspace patterns from package.json
    func ExtractWorkspacePatterns(pkg *PackageJSON) ([]string, error)
    ```
  - [x] 2.2 Handle both array format `["packages/*"]` and object format `{packages: [...], nohoist: [...]}`
  - [x] 2.3 Create tests in `pkg/parser/package_json_test.go`

- [x] **Task 3: Implement Glob Pattern Matching** (AC: #1, #2)
  - [x] 3.1 Create `pkg/parser/glob.go`:
    ```go
    // ExpandGlobPatterns expands workspace glob patterns to actual package paths
    // Supports: *, **, ? wildcards and negation patterns (!prefix)
    func ExpandGlobPatterns(rootPath string, patterns []string) ([]string, error)

    // MatchPattern checks if a path matches a glob pattern
    func MatchPattern(pattern, path string) bool
    ```
  - [x] 3.2 Support common patterns: `packages/*`, `apps/*`, `packages/**/*`
  - [x] 3.3 Support negation: `!packages/deprecated-*`
  - [x] 3.4 Use Go's `path/filepath.Glob` as base, extend for ** support
  - [x] 3.5 Create tests in `pkg/parser/glob_test.go` with edge cases

- [x] **Task 4: Implement pnpm-workspace.yaml Parser** (AC: #2)
  - [x] 4.1 Create `pkg/parser/pnpm_workspace.go`:
    ```go
    // PnpmWorkspace represents pnpm-workspace.yaml structure
    type PnpmWorkspace struct {
        Packages []string `yaml:"packages"`
    }

    // ParsePnpmWorkspace parses pnpm-workspace.yaml
    func ParsePnpmWorkspace(data []byte) (*PnpmWorkspace, error)
    ```
  - [x] 4.2 Add YAML parsing (use `gopkg.in/yaml.v3` - only external dep needed)
  - [x] 4.3 Create tests in `pkg/parser/pnpm_workspace_test.go`

- [x] **Task 5: Implement Main Parser Interface** (AC: #4, #5)
  - [x] 5.1 Update `pkg/parser/parser.go` with main parsing logic:
    ```go
    // Parser handles workspace configuration parsing
    type Parser struct {
        rootPath string
    }

    // NewParser creates a new workspace parser
    func NewParser(rootPath string) *Parser

    // Parse detects workspace type and parses all packages
    func (p *Parser) Parse(files map[string][]byte) (*types.WorkspaceData, error)

    // DetectWorkspaceType determines if workspace is npm/yarn/pnpm
    func (p *Parser) DetectWorkspaceType(files map[string][]byte) types.WorkspaceType
    ```
  - [x] 5.2 Implement workspace type auto-detection:
    - If `pnpm-workspace.yaml` exists → pnpm
    - If `yarn.lock` exists → yarn
    - If `package-lock.json` exists → npm
    - Otherwise → unknown (fallback to npm-style parsing)
  - [x] 5.3 Create comprehensive tests in `pkg/parser/parser_test.go`

- [x] **Task 6: Add Error Codes** (AC: #5)
  - [x] 6.1 Add parser-specific error codes to `internal/result/result.go`:
    ```go
    const (
        ErrInvalidWorkspace     = "INVALID_WORKSPACE"
        ErrMissingPackageJSON   = "MISSING_PACKAGE_JSON"
        ErrInvalidPackageJSON   = "INVALID_PACKAGE_JSON"
        ErrInvalidPnpmWorkspace = "INVALID_PNPM_WORKSPACE"
        ErrGlobPatternFailed    = "GLOB_PATTERN_FAILED"
    )
    ```
  - [x] 6.2 Ensure all parser errors use these codes

- [x] **Task 7: Wire Parser to WASM Entry Point** (AC: #4)
  - [x] 7.1 Update `internal/handlers/handlers.go` to use real parser:
    ```go
    func HandleAnalyze(input string) *result.Result {
        var filesInput map[string]string // filename -> content
        if err := json.Unmarshal([]byte(input), &filesInput); err != nil {
            return result.NewError(ErrInvalidInput, "Failed to parse input JSON")
        }

        // Convert to []byte map
        files := make(map[string][]byte)
        for name, content := range filesInput {
            files[name] = []byte(content)
        }

        parser := parser.NewParser("/workspace") // Virtual root
        workspaceData, err := parser.Parse(files)
        if err != nil {
            return result.NewError(ErrAnalysisFailed, err.Error())
        }

        // For now, return workspace data as analysis result
        // Full analysis logic comes in Story 2.2+
        return result.NewSuccess(workspaceData)
    }
    ```
  - [x] 7.2 Update handler tests to use real parser

- [x] **Task 8: Performance Testing** (AC: #6)
  - [x] 8.1 Create benchmark tests in `pkg/parser/benchmark_test.go`:
    ```go
    func BenchmarkParse100Packages(b *testing.B) {
        // Generate 100 package workspace configuration
        // Verify parsing completes in < 1 second
    }
    ```
  - [x] 8.2 Verify memory usage with `go test -bench=. -benchmem`
    - 100 packages: 414µs, 156KB memory
  - [x] 8.3 Optimize if any parsing step exceeds performance targets
    - No optimization needed - performance excellent

- [x] **Task 9: Integration Verification** (AC: all)
  - [x] 9.1 Build WASM: `pnpm nx build @monoguard/analysis-engine` (4.2MB)
  - [x] 9.2 Integration verified via handler tests (TestAnalyzeWithRealWorkspace, TestAnalyzeWithPnpmWorkspace)
  - [x] 9.3 Tested npm, yarn, pnpm workspace structures in parser tests
  - [x] 9.4 All tests pass with race detection: `make test`

## Dev Notes

### Architecture Patterns & Constraints

**From Architecture Document:**
- **Technology:** Go 1.21+ compiled to WASM (`GOOS=js GOARCH=wasm`)
- **Location:** `packages/analysis-engine/pkg/parser/`
- **Dependency:** Only YAML parsing requires external dependency (`gopkg.in/yaml.v3`)
- **Pattern:** Use Result<T> for all function returns

**Critical Constraints:**
- **camelCase JSON:** All struct tags MUST use camelCase
- **Result Pattern:** All errors wrapped in Result type
- **UPPER_SNAKE_CASE Errors:** Error codes MUST be UPPER_SNAKE_CASE
- **Zero Backend:** Parser runs entirely in browser via WASM

### Critical Don't-Miss Rules

**From project-context.md:**

1. **Go Naming Conventions:**
   - PascalCase for exported functions/types (e.g., `ParseWorkspace`)
   - camelCase for unexported functions (e.g., `extractPatterns`)
   - snake_case for file names (e.g., `package_json.go`, `pnpm_workspace.go`)
   - Test files: `*_test.go` in same directory

2. **JSON Serialization - CRITICAL:**
   ```go
   // ✅ CORRECT: camelCase JSON tags
   type WorkspaceData struct {
       RootPath      string `json:"rootPath"`
       WorkspaceType string `json:"workspaceType"`
   }

   // ❌ WRONG: snake_case breaks TypeScript
   type WorkspaceData struct {
       RootPath      string `json:"root_path"` // BREAKS FRONTEND
   }
   ```

3. **Result Type Pattern - MANDATORY:**
   ```go
   // ✅ CORRECT: Always wrap returns
   func Parse(files map[string][]byte) (*types.WorkspaceData, error) {
       // Implementation
   }

   // WASM handler wraps in Result:
   func HandleAnalyze(input string) *result.Result {
       data, err := parser.Parse(files)
       if err != nil {
           return result.NewError("PARSE_ERROR", err.Error())
       }
       return result.NewSuccess(data)
   }
   ```

4. **Error Codes - UPPER_SNAKE_CASE:**
   - `INVALID_WORKSPACE` (not `InvalidWorkspace` or `invalid-workspace`)
   - `MISSING_PACKAGE_JSON` (not `MissingPackageJson`)

### Project Structure Notes

**Target Directory Structure:**
```
packages/analysis-engine/
├── cmd/wasm/
│   └── main.go                    # WASM entry point (existing)
├── pkg/
│   ├── parser/
│   │   ├── parser.go              # Main parser interface (UPDATE)
│   │   ├── parser_test.go         # Parser tests (UPDATE)
│   │   ├── package_json.go        # NEW: package.json parsing
│   │   ├── package_json_test.go   # NEW: package.json tests
│   │   ├── pnpm_workspace.go      # NEW: pnpm-workspace.yaml parsing
│   │   ├── pnpm_workspace_test.go # NEW: pnpm tests
│   │   ├── glob.go                # NEW: glob pattern matching
│   │   ├── glob_test.go           # NEW: glob tests
│   │   └── benchmark_test.go      # NEW: performance benchmarks
│   ├── types/
│   │   ├── types.go               # UPDATE: Add WorkspaceData
│   │   └── types_test.go          # UPDATE: Add WorkspaceData tests
│   └── analyzer/
│       └── analyzer.go            # Placeholder (future stories)
├── internal/
│   ├── result/
│   │   ├── result.go              # UPDATE: Add error codes
│   │   └── result_test.go
│   └── handlers/
│       ├── handlers.go            # UPDATE: Use real parser
│       └── handlers_test.go       # UPDATE: Real parser tests
├── go.mod                         # UPDATE: Add yaml.v3 dependency
├── go.sum
└── Makefile
```

### Input/Output Format

**Input to WASM analyze function:**
```json
{
  "package.json": "{ \"name\": \"root\", \"workspaces\": [\"packages/*\"] }",
  "packages/pkg-a/package.json": "{ \"name\": \"@mono/pkg-a\", \"version\": \"1.0.0\", \"dependencies\": { \"@mono/pkg-b\": \"^1.0.0\" } }",
  "packages/pkg-b/package.json": "{ \"name\": \"@mono/pkg-b\", \"version\": \"1.0.0\" }"
}
```

**Output WorkspaceData:**
```json
{
  "data": {
    "rootPath": "/workspace",
    "workspaceType": "npm",
    "packages": {
      "@mono/pkg-a": {
        "name": "@mono/pkg-a",
        "version": "1.0.0",
        "path": "packages/pkg-a",
        "dependencies": { "@mono/pkg-b": "^1.0.0" },
        "devDependencies": {},
        "peerDependencies": {}
      },
      "@mono/pkg-b": {
        "name": "@mono/pkg-b",
        "version": "1.0.0",
        "path": "packages/pkg-b",
        "dependencies": {},
        "devDependencies": {},
        "peerDependencies": {}
      }
    }
  },
  "error": null
}
```

### Previous Story Intelligence

**From Story 1.3 (Go WASM Setup):**
- Go module: `github.com/j620656786206/MonoGuard/packages/analysis-engine`
- Result type in `internal/result/result.go` with `NewSuccess()`, `NewError()`, `ToJSON()`
- Handlers in `internal/handlers/handlers.go` - refactored for testability
- WASM exports: `MonoGuard.analyze()`, `MonoGuard.check()`, `MonoGuard.getVersion()`
- Current analyze returns placeholder data - needs real parser

**From Story 1.5 (TypeScript Types):**
- `@monoguard/types` defines TypeScript equivalents
- `WorkspaceType`: 'npm' | 'yarn' | 'pnpm' | 'unknown'
- `DependencyType`: 'production' | 'development' | 'peer' | 'optional'
- Types must match exactly for JSON serialization to work

**Key Learnings from Epic 1:**
- Always test JSON serialization to ensure camelCase
- Use table-driven tests for multiple scenarios
- Handler pattern separates WASM interop from business logic
- Coverage target: > 80%

### Git Intelligence

**Recent commits:**
- `a71f05f` fix(review): address code review issues for story 1-8
- `aa80f65` fix(wasm): improve escapeJSON and increase test coverage
- Story 1-8 set up testing framework and code quality tools

**Patterns established:**
- Commit format: `type(scope): description`
- Code review fixes in separate commits
- Test coverage improvements tracked

### Testing Requirements

**Unit Tests:**
- Each parser function needs tests
- Test valid inputs, invalid inputs, edge cases
- Table-driven tests preferred
- Target > 80% coverage

**Test Cases to Cover:**
- npm workspaces with array format
- npm workspaces with object format
- yarn workspaces (same as npm)
- pnpm workspaces with pnpm-workspace.yaml
- Glob patterns: `*`, `**`, negation `!`
- Missing package.json in workspace
- Invalid JSON in package.json
- Empty workspaces array
- Nested packages
- Packages without name field

**Performance Benchmarks:**
- Parse 10 packages: < 100ms
- Parse 100 packages: < 1 second
- Memory: No excessive allocations

### TypeScript Type Reference

The Go types must produce JSON that matches these TypeScript types from `@monoguard/types`:

```typescript
// From packages/types/src/analysis/graph.ts
export type WorkspaceType = 'npm' | 'yarn' | 'pnpm' | 'unknown';

export interface PackageNode {
  name: string;
  version: string;
  path: string;
  dependencies: string[];
  devDependencies: string[];
  peerDependencies: string[];
}
```

Note: TypeScript `PackageNode.dependencies` is `string[]` but for parsing we need `map[string]string` to capture version ranges. Update TypeScript types if needed in a follow-up, or adjust Go output to match.

### External Dependencies

**YAML Parsing:**
- Need to add `gopkg.in/yaml.v3` for pnpm-workspace.yaml parsing
- Run: `go get gopkg.in/yaml.v3`
- This is the only external dependency for Epic 2

**Built-in Go packages used:**
- `encoding/json` - JSON parsing
- `path/filepath` - Path manipulation and glob
- `strings` - String manipulation
- `regexp` - Pattern matching for complex globs

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story 2.1]
- [Source: _bmad-output/project-context.md#Go Naming Conventions]
- [Source: _bmad-output/project-context.md#JSON Serialization]
- [Source: _bmad-output/implementation-artifacts/1-3-setup-go-wasm-analysis-engine-project.md]
- [Source: _bmad-output/implementation-artifacts/1-5-setup-shared-typescript-types-package.md]
- [Go filepath.Glob Documentation](https://pkg.go.dev/path/filepath#Glob)
- [YAML v3 Package](https://pkg.go.dev/gopkg.in/yaml.v3)

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
