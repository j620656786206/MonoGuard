# Story 3.2: Trace Import Statement Paths

Status: ready-for-dev

## Story

As a **user**,
I want **to see exactly which import statements create the circular dependency**,
So that **I know the specific code locations that need to be modified**.

## Acceptance Criteria

1. **AC1: Import Statement Detection**
   - Given source files provided in the workspace upload
   - When I analyze a circular dependency
   - Then the trace identifies:
     - Specific import statements (e.g., `import { foo } from '@pkg/bar'`)
     - File paths containing the imports
     - Line numbers where imports occur
   - And results are returned as `ImportTrace[]`

2. **AC2: ESM Import Support**
   - Given TypeScript/JavaScript files with ES Module imports
   - When parsing imports
   - Then the parser detects:
     - Named imports: `import { foo, bar } from '@pkg/core'`
     - Default imports: `import foo from '@pkg/core'`
     - Namespace imports: `import * as core from '@pkg/core'`
     - Side-effect imports: `import '@pkg/core'`
     - Dynamic imports: `import('@pkg/core')`
   - And extracts the target package name from each

3. **AC3: CommonJS Require Support**
   - Given JavaScript files with CommonJS requires
   - When parsing requires
   - Then the parser detects:
     - Standard require: `const foo = require('@pkg/core')`
     - Destructured require: `const { foo } = require('@pkg/core')`
     - Direct require: `require('@pkg/core')`
   - And handles both single and double quoted strings

4. **AC4: Import Chain Construction**
   - Given a circular dependency cycle A → B → C → A
   - When tracing imports
   - Then construct complete import chain:
     - A imports B: file, line, statement
     - B imports C: file, line, statement
     - C imports A: file, line, statement (completes the cycle)
   - And chain order matches cycle order

5. **AC5: Integration with CircularDependencyInfo**
   - Given analysis results
   - When enriching CircularDependencyInfo
   - Then add optional `importTraces` field:
     ```go
     type CircularDependencyInfo struct {
         // ... existing fields ...
         ImportTraces []ImportTrace `json:"importTraces,omitempty"`
     }
     ```
   - And traces are only populated when source files are provided

6. **AC6: Graceful Degradation**
   - Given a workspace upload without source files
   - When import tracing is requested
   - Then:
     - Return empty `importTraces` array (not null)
     - Analysis continues without error
     - Other features work normally
   - And no error is thrown for missing source files

7. **AC7: Performance**
   - Given a workspace with 100 packages and 500 source files
   - When import tracing runs
   - Then tracing completes in < 2 seconds additional overhead
   - And memory usage increase is < 50MB

## Tasks / Subtasks

- [ ] **Task 1: Define ImportTrace Type** (AC: #1, #4)
  - [ ] 1.1 Create `pkg/types/import_trace.go`:
    ```go
    package types

    // ImportTrace represents a single import statement that contributes to a cycle.
    // Matches @monoguard/types ImportTrace interface.
    type ImportTrace struct {
        // FromPackage is the package containing the import
        FromPackage string `json:"fromPackage"`

        // ToPackage is the package being imported
        ToPackage string `json:"toPackage"`

        // FilePath is the relative path to the file containing the import
        FilePath string `json:"filePath"`

        // LineNumber is the 1-based line number of the import statement
        LineNumber int `json:"lineNumber"`

        // Statement is the actual import/require statement text
        Statement string `json:"statement"`

        // ImportType classifies the import style
        ImportType ImportType `json:"importType"`

        // Symbols are the specific imports (empty for namespace/side-effect imports)
        Symbols []string `json:"symbols,omitempty"`
    }

    // ImportType classifies the import style.
    type ImportType string

    const (
        ImportTypeESMNamed     ImportType = "esm-named"     // import { foo } from 'bar'
        ImportTypeESMDefault   ImportType = "esm-default"   // import foo from 'bar'
        ImportTypeESMNamespace ImportType = "esm-namespace" // import * as foo from 'bar'
        ImportTypeESMSideEffect ImportType = "esm-side-effect" // import 'bar'
        ImportTypeESMDynamic   ImportType = "esm-dynamic"   // import('bar')
        ImportTypeCJSRequire   ImportType = "cjs-require"   // require('bar')
    )
    ```
  - [ ] 1.2 Add JSON serialization tests in `pkg/types/import_trace_test.go`
  - [ ] 1.3 Ensure all JSON tags use camelCase

- [ ] **Task 2: Create Import Parser** (AC: #2, #3)
  - [ ] 2.1 Create `pkg/parser/import_parser.go`:
    ```go
    package parser

    import (
        "regexp"
        "strings"

        "github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
    )

    // ImportParser extracts import statements from source files.
    type ImportParser struct {
        // Regex patterns for different import types
        esmNamedPattern     *regexp.Regexp
        esmDefaultPattern   *regexp.Regexp
        esmNamespacePattern *regexp.Regexp
        esmSideEffectPattern *regexp.Regexp
        esmDynamicPattern   *regexp.Regexp
        cjsRequirePattern   *regexp.Regexp
    }

    // NewImportParser creates a new parser with compiled regex patterns.
    func NewImportParser() *ImportParser

    // ParseFile extracts all imports from a source file.
    // Returns imports that reference the specified target packages.
    func (ip *ImportParser) ParseFile(content []byte, filePath string, targetPackages map[string]bool) []types.ImportTrace

    // parseESMImports extracts ES Module import statements.
    func (ip *ImportParser) parseESMImports(content string, filePath string, targets map[string]bool) []types.ImportTrace

    // parseCJSRequires extracts CommonJS require statements.
    func (ip *ImportParser) parseCJSRequires(content string, filePath string, targets map[string]bool) []types.ImportTrace

    // extractPackageName extracts the package name from an import path.
    // Handles scoped packages (@scope/pkg) and subpath imports (pkg/submodule).
    func extractPackageName(importPath string) string
    ```
  - [ ] 2.2 Implement ESM import regex patterns:
    ```go
    // ESM Named: import { foo, bar } from 'package'
    esmNamedPattern := regexp.MustCompile(`import\s*\{([^}]+)\}\s*from\s*['"]([^'"]+)['"]`)

    // ESM Default: import foo from 'package'
    esmDefaultPattern := regexp.MustCompile(`import\s+(\w+)\s+from\s*['"]([^'"]+)['"]`)

    // ESM Namespace: import * as foo from 'package'
    esmNamespacePattern := regexp.MustCompile(`import\s*\*\s*as\s+(\w+)\s+from\s*['"]([^'"]+)['"]`)

    // ESM Side-effect: import 'package'
    esmSideEffectPattern := regexp.MustCompile(`import\s*['"]([^'"]+)['"]`)

    // ESM Dynamic: import('package')
    esmDynamicPattern := regexp.MustCompile(`import\s*\(\s*['"]([^'"]+)['"]\s*\)`)
    ```
  - [ ] 2.3 Implement CJS require regex patterns:
    ```go
    // CJS Require: require('package') - matches various forms
    cjsRequirePattern := regexp.MustCompile(`require\s*\(\s*['"]([^'"]+)['"]\s*\)`)
    ```
  - [ ] 2.4 Create comprehensive tests in `pkg/parser/import_parser_test.go`

- [ ] **Task 3: Implement Line Number Tracking** (AC: #1)
  - [ ] 3.1 Implement line number calculation:
    ```go
    // getLineNumber returns the 1-based line number for a byte offset in content.
    func getLineNumber(content string, offset int) int {
        lines := 1
        for i := 0; i < offset && i < len(content); i++ {
            if content[i] == '\n' {
                lines++
            }
        }
        return lines
    }
    ```
  - [ ] 3.2 Track match positions when parsing
  - [ ] 3.3 Add line number tests

- [ ] **Task 4: Create Import Tracer** (AC: #4, #5)
  - [ ] 4.1 Create `pkg/analyzer/import_tracer.go`:
    ```go
    package analyzer

    import (
        "path/filepath"
        "strings"

        "github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/parser"
        "github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
    )

    // ImportTracer traces import statements that create circular dependencies.
    type ImportTracer struct {
        workspace *types.WorkspaceData
        files     map[string][]byte // Source files (*.ts, *.js, *.tsx, *.jsx)
        parser    *parser.ImportParser
    }

    // NewImportTracer creates a new tracer for the given workspace and files.
    func NewImportTracer(workspace *types.WorkspaceData, files map[string][]byte) *ImportTracer

    // Trace finds import statements that form the circular dependency.
    func (it *ImportTracer) Trace(cycle *types.CircularDependencyInfo) []types.ImportTrace

    // traceEdge finds imports from one package to another.
    func (it *ImportTracer) traceEdge(fromPkg, toPkg string) []types.ImportTrace

    // getSourceFilesForPackage returns source files belonging to a package.
    func (it *ImportTracer) getSourceFilesForPackage(pkgName string) map[string][]byte

    // isSourceFile checks if a file is a parseable source file.
    func isSourceFile(path string) bool {
        ext := strings.ToLower(filepath.Ext(path))
        return ext == ".ts" || ext == ".tsx" || ext == ".js" || ext == ".jsx" || ext == ".mjs" || ext == ".cjs"
    }
    ```
  - [ ] 4.2 Implement trace logic to find imports between packages
  - [ ] 4.3 Create comprehensive tests in `pkg/analyzer/import_tracer_test.go`

- [ ] **Task 5: Handle Package Name Extraction** (AC: #2, #3)
  - [ ] 5.1 Implement `extractPackageName`:
    ```go
    // extractPackageName extracts the package name from an import path.
    // Examples:
    //   '@scope/pkg'       → '@scope/pkg'
    //   '@scope/pkg/sub'   → '@scope/pkg'
    //   'lodash'           → 'lodash'
    //   'lodash/debounce'  → 'lodash'
    //   './local'          → '' (relative import, skip)
    //   '../parent'        → '' (relative import, skip)
    func extractPackageName(importPath string) string {
        // Skip relative imports
        if strings.HasPrefix(importPath, ".") {
            return ""
        }

        // Handle scoped packages (@scope/pkg)
        if strings.HasPrefix(importPath, "@") {
            parts := strings.SplitN(importPath, "/", 3)
            if len(parts) >= 2 {
                return parts[0] + "/" + parts[1]
            }
            return importPath
        }

        // Handle regular packages
        parts := strings.SplitN(importPath, "/", 2)
        return parts[0]
    }
    ```
  - [ ] 5.2 Add tests for scoped packages, subpaths, and edge cases

- [ ] **Task 6: Integrate with CircularDependencyInfo** (AC: #5)
  - [ ] 6.1 Update `pkg/types/circular.go`:
    ```go
    type CircularDependencyInfo struct {
        Cycle        []string            `json:"cycle"`
        Type         CircularType        `json:"type"`
        Severity     CircularSeverity    `json:"severity"`
        Depth        int                 `json:"depth"`
        Impact       string              `json:"impact"`
        Complexity   int                 `json:"complexity"`
        RootCause    *RootCauseAnalysis  `json:"rootCause,omitempty"`    // Story 3.1
        ImportTraces []ImportTrace       `json:"importTraces,omitempty"` // Story 3.2 NEW
    }
    ```
  - [ ] 6.2 Verify existing tests still pass (backward compatible)

- [ ] **Task 7: Wire to Analyzer Pipeline** (AC: all)
  - [ ] 7.1 Update `pkg/analyzer/analyzer.go`:
    ```go
    // Analyze performs complete workspace analysis.
    // sourceFiles is optional - if provided, enables import tracing.
    func (a *Analyzer) AnalyzeWithSources(
        workspace *types.WorkspaceData,
        sourceFiles map[string][]byte,
    ) (*types.AnalysisResult, error) {
        // ... existing analysis ...

        // Enrich cycles with root cause analysis (Story 3.1)
        rootCauseAnalyzer := NewRootCauseAnalyzer(graph)
        for _, cycle := range cycles {
            cycle.RootCause = rootCauseAnalyzer.Analyze(cycle)
        }

        // NEW: Enrich cycles with import traces (Story 3.2)
        if len(sourceFiles) > 0 {
            importTracer := NewImportTracer(workspace, sourceFiles)
            for _, cycle := range cycles {
                cycle.ImportTraces = importTracer.Trace(cycle)
            }
        }

        return result, nil
    }
    ```
  - [ ] 7.2 Update WASM handler to accept optional source files
  - [ ] 7.3 Maintain backward compatibility (analyze without sources still works)

- [ ] **Task 8: Implement Graceful Degradation** (AC: #6)
  - [ ] 8.1 Ensure empty sourceFiles returns empty traces (not nil):
    ```go
    func (it *ImportTracer) Trace(cycle *types.CircularDependencyInfo) []types.ImportTrace {
        if it.files == nil || len(it.files) == 0 {
            return []types.ImportTrace{} // Empty slice, not nil
        }
        // ... trace logic
    }
    ```
  - [ ] 8.2 Add tests for graceful degradation scenarios
  - [ ] 8.3 Verify no errors when source files missing

- [ ] **Task 9: Update TypeScript Types** (AC: #5)
  - [ ] 9.1 Update `packages/types/src/analysis/results.ts`:
    ```typescript
    export interface ImportTrace {
      fromPackage: string;
      toPackage: string;
      filePath: string;
      lineNumber: number;
      statement: string;
      importType: 'esm-named' | 'esm-default' | 'esm-namespace' | 'esm-side-effect' | 'esm-dynamic' | 'cjs-require';
      symbols?: string[];
    }

    export interface CircularDependencyInfo {
      cycle: string[];
      type: 'direct' | 'indirect';
      severity: 'critical' | 'warning' | 'info';
      depth: number;
      impact: string;
      complexity: number;
      rootCause?: RootCauseAnalysis;   // Story 3.1
      importTraces?: ImportTrace[];     // Story 3.2 NEW
    }
    ```
  - [ ] 9.2 Run `pnpm nx build types` to verify
  - [ ] 9.3 Update type tests if needed

- [ ] **Task 10: Performance Testing** (AC: #7)
  - [ ] 10.1 Create `pkg/analyzer/import_tracer_benchmark_test.go`:
    ```go
    func BenchmarkImportTracing(b *testing.B) {
        workspace := generateWorkspace(100)
        sourceFiles := generateSourceFiles(500) // 500 files, ~100 lines each
        cycles := generateCycles(5)
        tracer := NewImportTracer(workspace, sourceFiles)

        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            for _, cycle := range cycles {
                tracer.Trace(cycle)
            }
        }
    }
    ```
  - [ ] 10.2 Verify < 2 seconds for 100 packages with 500 source files
  - [ ] 10.3 Document actual performance in completion notes

- [ ] **Task 11: Integration Verification** (AC: all)
  - [ ] 11.1 Run all tests: `cd packages/analysis-engine && make test`
  - [ ] 11.2 Build WASM: `pnpm nx build @monoguard/analysis-engine`
  - [ ] 11.3 Run affected CI checks: `pnpm nx affected --target=lint,test,type-check --base=main`
  - [ ] 11.4 Test with real monorepo source files
  - [ ] 11.5 Verify JSON output includes importTraces field

## Dev Notes

### Architecture Patterns & Constraints

**From Architecture Document:**
- **Location:** Import parser in `pkg/parser/`, Import tracer in `pkg/analyzer/`
- **Pattern:** Parser + Tracer separation (single responsibility)
- **Integration:** Enriches existing CircularDependencyInfo with optional ImportTraces
- **Privacy:** Users explicitly choose to upload source files; no automatic file access

**Critical Constraints:**
- **camelCase JSON:** All struct tags MUST use camelCase
- **Optional Fields:** ImportTraces is `omitempty` - empty when no source files
- **Regex-Based Parsing:** Use regex for Go WASM (no external AST parser available)
- **Performance:** Must not significantly slow down analysis

**Why Regex Instead of AST:**
- Go WASM has limited ecosystem for TypeScript/JavaScript AST parsing
- Regex is sufficient for import statement extraction (imports have predictable syntax)
- Keeps WASM bundle size small
- Performance is excellent for pattern matching

### Critical Don't-Miss Rules

**From project-context.md:**

1. **JSON Naming Convention:**
   ```go
   // ✅ CORRECT: camelCase JSON tags
   type ImportTrace struct {
       FromPackage string     `json:"fromPackage"`
       ToPackage   string     `json:"toPackage"`
       LineNumber  int        `json:"lineNumber"`
       ImportType  ImportType `json:"importType"`
   }

   // ❌ WRONG: snake_case JSON tags
   type ImportTrace struct {
       FromPackage string `json:"from_package"` // WRONG!
   }
   ```

2. **Empty Slice vs Nil:**
   ```go
   // ✅ CORRECT: Return empty slice for JSON serialization
   func (it *ImportTracer) Trace(cycle *types.CircularDependencyInfo) []types.ImportTrace {
       if it.files == nil {
           return []types.ImportTrace{} // Serializes as []
       }
   }

   // ❌ WRONG: Return nil (serializes as null, may break consumers)
   func (it *ImportTracer) Trace(...) []types.ImportTrace {
       if it.files == nil {
           return nil // Serializes as null
       }
   }
   ```

3. **Test File Naming:**
   ```
   ✅ CORRECT:
   pkg/parser/import_parser.go
   pkg/parser/import_parser_test.go
   pkg/analyzer/import_tracer.go
   pkg/analyzer/import_tracer_test.go

   ❌ WRONG:
   pkg/parser/__tests__/import_parser.test.go
   ```

### Project Structure Notes

**Target Directory Structure:**
```
packages/analysis-engine/
├── pkg/
│   ├── parser/
│   │   ├── parser.go                     # Existing workspace parser
│   │   ├── import_parser.go              # NEW: Import statement parser
│   │   └── import_parser_test.go         # NEW: Parser tests
│   ├── analyzer/
│   │   ├── analyzer.go                   # UPDATE: Add import tracing call
│   │   ├── root_cause_analyzer.go        # From Story 3.1
│   │   ├── import_tracer.go              # NEW: Import tracer
│   │   ├── import_tracer_test.go         # NEW: Tracer tests
│   │   └── import_tracer_benchmark_test.go # NEW: Performance
│   └── types/
│       ├── circular.go                   # UPDATE: Add ImportTraces field
│       ├── import_trace.go               # NEW: ImportTrace type
│       └── import_trace_test.go          # NEW: Type tests
└── ...

packages/types/src/analysis/
├── results.ts                            # UPDATE: Add TS types
└── ...
```

### Previous Story Intelligence

**From Story 3.1 (ready-for-dev):**
- CircularDependencyInfo will have RootCause field added
- DependencyEdge type defined (from, to, type, critical)
- Analyzer pipeline pattern: detect cycles → enrich with analysis
- **Key Insight:** Follow same pattern for import traces

**From Story 2.3 (done):**
- Cycle format: `["A", "B", "C", "A"]` (starts and ends with same)
- Cycles sorted by severity, then depth
- **Key Insight:** Import traces should follow cycle order

**From Parser (parser.go):**
- Files provided as `map[string][]byte` (path → content)
- Package paths stored in WorkspaceData.Packages[name].Path
- **Key Insight:** Can map source files to packages via path prefix

### Regex Patterns for Import Parsing

**ESM Import Patterns:**
```go
// Named imports: import { foo, bar as baz } from 'package'
// Captures: group(1) = imports, group(2) = package
esmNamedPattern := `import\s*\{([^}]+)\}\s*from\s*['"]([^'"]+)['"]`

// Default imports: import foo from 'package'
// Captures: group(1) = name, group(2) = package
esmDefaultPattern := `import\s+(\w+)\s+from\s*['"]([^'"]+)['"]`

// Namespace imports: import * as foo from 'package'
// Captures: group(1) = alias, group(2) = package
esmNamespacePattern := `import\s*\*\s*as\s+(\w+)\s+from\s*['"]([^'"]+)['"]`

// Side-effect imports: import 'package'
// Captures: group(1) = package
esmSideEffectPattern := `^\s*import\s*['"]([^'"]+)['"]\s*;?\s*$`

// Dynamic imports: import('package') or await import('package')
// Captures: group(1) = package
esmDynamicPattern := `import\s*\(\s*['"]([^'"]+)['"]\s*\)`
```

**CommonJS Patterns:**
```go
// Standard require: require('package')
// Destructured: const { foo } = require('package')
// Captures: group(1) = package
cjsRequirePattern := `require\s*\(\s*['"]([^'"]+)['"]\s*\)`
```

**Pattern Priority (to avoid false matches):**
1. ESM Named (most specific with braces)
2. ESM Namespace (has `* as`)
3. ESM Default (has identifier before `from`)
4. ESM Side-effect (just import + string)
5. ESM Dynamic (import function call)
6. CJS Require (require function call)

### Input/Output Format

**Input (Files Map with Source Files):**
```go
files := map[string][]byte{
    "package.json":                     []byte(`{...}`),
    "packages/ui/package.json":         []byte(`{"name": "@mono/ui"}`),
    "packages/ui/src/index.ts":         []byte(`import { api } from '@mono/api'`),
    "packages/api/package.json":        []byte(`{"name": "@mono/api"}`),
    "packages/api/src/client.ts":       []byte(`import { core } from '@mono/core'`),
    "packages/core/package.json":       []byte(`{"name": "@mono/core"}`),
    "packages/core/src/utils.ts":       []byte(`import { ui } from '@mono/ui'`), // Creates cycle!
}
```

**Output (CircularDependencyInfo with ImportTraces):**
```json
{
  "cycle": ["@mono/ui", "@mono/api", "@mono/core", "@mono/ui"],
  "type": "indirect",
  "severity": "info",
  "depth": 3,
  "impact": "Indirect circular dependency involving 3 packages",
  "complexity": 5,
  "rootCause": { ... },
  "importTraces": [
    {
      "fromPackage": "@mono/ui",
      "toPackage": "@mono/api",
      "filePath": "packages/ui/src/index.ts",
      "lineNumber": 1,
      "statement": "import { api } from '@mono/api'",
      "importType": "esm-named",
      "symbols": ["api"]
    },
    {
      "fromPackage": "@mono/api",
      "toPackage": "@mono/core",
      "filePath": "packages/api/src/client.ts",
      "lineNumber": 1,
      "statement": "import { core } from '@mono/core'",
      "importType": "esm-named",
      "symbols": ["core"]
    },
    {
      "fromPackage": "@mono/core",
      "toPackage": "@mono/ui",
      "filePath": "packages/core/src/utils.ts",
      "lineNumber": 1,
      "statement": "import { ui } from '@mono/ui'",
      "importType": "esm-named",
      "symbols": ["ui"]
    }
  ]
}
```

### Test Scenarios

| Scenario | Input | Expected Output |
|----------|-------|-----------------|
| ESM named import | `import { foo } from '@pkg/bar'` | `esm-named`, symbols: ["foo"] |
| ESM default import | `import foo from '@pkg/bar'` | `esm-default`, symbols: [] |
| ESM namespace import | `import * as bar from '@pkg/bar'` | `esm-namespace`, symbols: [] |
| ESM side-effect | `import '@pkg/bar'` | `esm-side-effect`, symbols: [] |
| ESM dynamic | `import('@pkg/bar')` | `esm-dynamic`, symbols: [] |
| CJS require | `require('@pkg/bar')` | `cjs-require`, symbols: [] |
| Scoped package | `import x from '@scope/pkg/sub'` | package: `@scope/pkg` |
| No source files | Empty files map | Empty importTraces: [] |
| Relative import | `import './local'` | Skip (not a package) |

### Edge Cases to Handle

1. **Multi-line imports:**
   ```typescript
   import {
     foo,
     bar,
     baz
   } from '@pkg/core';
   ```
   → Use `(?s)` flag or `[\s\S]` for multi-line matching

2. **Comments containing imports:**
   ```typescript
   // import { foo } from '@pkg/bar'
   /* import { foo } from '@pkg/bar' */
   ```
   → Consider stripping comments first, or accept false positives (low impact)

3. **String literals containing import text:**
   ```typescript
   const code = "import { foo } from '@pkg/bar'";
   ```
   → Accept as import (conservative approach, user can verify)

4. **Re-exports:**
   ```typescript
   export { foo } from '@pkg/bar';
   export * from '@pkg/bar';
   ```
   → Treat as imports (they do create dependencies)

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story 3.2]
- [Source: _bmad-output/project-context.md#Language-Specific Rules]
- [Source: packages/analysis-engine/pkg/parser/parser.go]
- [Source: packages/analysis-engine/pkg/types/circular.go]
- [Source: _bmad-output/implementation-artifacts/3-1-implement-root-cause-analysis-for-circular-dependencies.md]
- [MDN: import statement](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/import)
- [Node.js: require](https://nodejs.org/api/modules.html#requireid)

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List

### Change Log
