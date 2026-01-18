# Story 2.6: Implement Package Exclusion Patterns

Status: review

## Story

As a **user**,
I want **to exclude specific packages or patterns from analysis**,
So that **I can focus on relevant parts of my monorepo**.

## Acceptance Criteria

1. **AC1: Exact Package Name Exclusion**
   - Given an exclusion list with exact package names
   - When I run analysis with exclusions
   - Then packages matching exact names are excluded
   - Example: `["packages/legacy", "@mono/deprecated"]`

2. **AC2: Glob Pattern Exclusion**
   - Given an exclusion list with glob patterns
   - When I run analysis with exclusions
   - Then packages matching glob patterns are excluded
   - Supports: `*` (any), `**` (recursive), `?` (single char)
   - Example: `["packages/deprecated-*", "**/test-*"]`

3. **AC3: Regex Pattern Exclusion**
   - Given an exclusion list with regex patterns (prefixed with `regex:`)
   - When I run analysis with exclusions
   - Then packages matching regex are excluded
   - Example: `["regex:^packages/legacy-.*$", "regex:.*-deprecated$"]`

4. **AC4: Exclusion from Metrics**
   - Given excluded packages
   - When metrics are calculated
   - Then excluded packages:
     - Do NOT contribute to health score
     - Do NOT appear in circular dependency detection
     - Do NOT appear in version conflict detection
     - Do NOT count toward package total
   - And exclusions are applied BEFORE all calculations

5. **AC5: Graph Representation**
   - Given excluded packages
   - When the dependency graph is returned
   - Then excluded packages:
     - ARE included in the graph nodes (for visualization)
     - Are marked with `excluded: true` flag
     - Have edges to/from them included
   - So the frontend can display them grayed out

6. **AC6: Configuration via API**
   - Given the analyze function
   - When I provide exclusion configuration
   - Then exclusions are accepted via:
     - `AnalysisConfig.exclude` array in input JSON
     - Patterns can be mixed (exact, glob, regex)
   - Example input:
     ```json
     {
       "files": { ... },
       "config": {
         "exclude": ["packages/legacy", "packages/deprecated-*", "regex:.*-test$"]
       }
     }
     ```

7. **AC7: Performance**
   - Given 100 packages with 20 exclusion patterns
   - When exclusion matching runs
   - Then it completes in < 50ms

## Tasks / Subtasks

- [x] **Task 1: Define AnalysisConfig Types** (AC: #6)
  - [x] 1.1 Create `pkg/types/config.go`:
    ```go
    package types

    // AnalysisConfig holds configuration options for analysis
    type AnalysisConfig struct {
        Exclude []string `json:"exclude"` // Exclusion patterns
    }

    // AnalysisInput represents the complete input to analyze function
    type AnalysisInput struct {
        Files  map[string]string `json:"files"`  // filename -> content
        Config *AnalysisConfig   `json:"config"` // Optional config
    }
    ```
  - [x] 1.2 Add JSON serialization tests
  - [x] 1.3 Update WASM handler to accept AnalysisInput

- [x] **Task 2: Implement Pattern Matcher** (AC: #1, #2, #3)
  - [x] 2.1 Create `pkg/analyzer/exclusion_matcher.go`:
    ```go
    package analyzer

    import (
        "regexp"
        "strings"
    )

    // ExclusionMatcher handles package exclusion pattern matching
    type ExclusionMatcher struct {
        exactMatches  map[string]bool
        globPatterns  []string
        regexPatterns []*regexp.Regexp
    }

    // NewExclusionMatcher creates a matcher from exclusion patterns
    func NewExclusionMatcher(patterns []string) (*ExclusionMatcher, error)

    // IsExcluded checks if a package name matches any exclusion pattern
    func (em *ExclusionMatcher) IsExcluded(packageName string) bool

    // parsePattern categorizes a pattern as exact, glob, or regex
    func parsePattern(pattern string) (patternType, cleanPattern string)

    // matchGlob matches a package name against a glob pattern
    func matchGlob(pattern, name string) bool

    // Constants for pattern types
    const (
        PatternTypeExact = "exact"
        PatternTypeGlob  = "glob"
        PatternTypeRegex = "regex"
    )
    ```
  - [x] 2.2 Implement exact matching (simple map lookup)
  - [x] 2.3 Implement glob matching (support *, **, ?)
  - [x] 2.4 Implement regex matching (patterns prefixed with `regex:`)
  - [x] 2.5 Create comprehensive tests in `pkg/analyzer/exclusion_matcher_test.go`

- [x] **Task 3: Update PackageNode for Exclusion Flag** (AC: #5)
  - [x] 3.1 Update `pkg/types/graph.go`:
    ```go
    type PackageNode struct {
        Name             string   `json:"name"`
        Version          string   `json:"version"`
        Path             string   `json:"path"`
        Dependencies     []string `json:"dependencies"`
        DevDependencies  []string `json:"devDependencies"`
        PeerDependencies []string `json:"peerDependencies"`
        Excluded         bool     `json:"excluded,omitempty"` // NEW: True if excluded
    }
    ```
  - [x] 3.2 Update graph builder tests

- [x] **Task 4: Integrate Exclusion into Graph Builder** (AC: #4, #5)
  - [x] 4.1 Update `pkg/analyzer/graph_builder.go`:
    ```go
    type GraphBuilder struct {
        workspacePackages map[string]bool
        exclusionMatcher  *ExclusionMatcher // NEW
    }

    func NewGraphBuilder(excludePatterns []string) (*GraphBuilder, error) {
        matcher, err := NewExclusionMatcher(excludePatterns)
        if err != nil {
            return nil, err
        }
        return &GraphBuilder{
            workspacePackages: make(map[string]bool),
            exclusionMatcher:  matcher,
        }, nil
    }

    func (gb *GraphBuilder) buildNodes(workspace *types.WorkspaceData) map[string]*types.PackageNode {
        nodes := make(map[string]*types.PackageNode)
        for name, pkg := range workspace.Packages {
            node := &types.PackageNode{
                Name:     name,
                // ... other fields
                Excluded: gb.exclusionMatcher.IsExcluded(name), // Mark excluded
            }
            nodes[name] = node
        }
        return nodes
    }
    ```
  - [x] 4.2 Update tests

- [x] **Task 5: Filter Excluded from Detectors** (AC: #4)
  - [x] 5.1 Update `pkg/analyzer/cycle_detector.go`:
    ```go
    func (cd *CycleDetector) DetectCycles() []*types.CircularDependencyInfo {
        // Filter out excluded packages before detection
        filteredGraph := cd.filterExcluded()
        // Run detection on filtered graph
        return cd.detectOnGraph(filteredGraph)
    }

    func (cd *CycleDetector) filterExcluded() *types.DependencyGraph {
        // Create new graph with only non-excluded packages
    }
    ```
  - [x] 5.2 Update `pkg/analyzer/conflict_detector.go`:
    ```go
    func (cd *ConflictDetector) DetectConflicts() []*types.VersionConflict {
        // Skip excluded packages when collecting dependencies
        for pkgName, pkg := range cd.workspace.Packages {
            if cd.isExcluded(pkgName) {
                continue
            }
            // ... collect dependencies
        }
    }
    ```
  - [x] 5.3 Update `pkg/analyzer/health_calculator.go` to exclude from depth/coupling
  - [x] 5.4 Update tests for all detectors

- [x] **Task 6: Update Analyzer to Accept Config** (AC: #6)
  - [x] 6.1 Update `pkg/analyzer/analyzer.go`:
    ```go
    type Analyzer struct {
        config *types.AnalysisConfig
    }

    func NewAnalyzer(config *types.AnalysisConfig) *Analyzer {
        if config == nil {
            config = &types.AnalysisConfig{}
        }
        return &Analyzer{config: config}
    }

    func (a *Analyzer) Analyze(workspace *types.WorkspaceData) (*types.AnalysisResult, error) {
        // Create graph builder with exclusions
        graphBuilder, err := NewGraphBuilder(a.config.Exclude)
        if err != nil {
            return nil, err
        }
        // ... rest of analysis
    }
    ```
  - [x] 6.2 Update handler to parse AnalysisInput
  - [x] 6.3 Update tests

- [x] **Task 7: Update AnalysisResult** (AC: #4)
  - [x] 7.1 Add excluded count to result:
    ```go
    type AnalysisResult struct {
        // ... existing fields
        ExcludedPackages int `json:"excludedPackages"` // Count of excluded
    }
    ```
  - [x] 7.2 Update tests

- [x] **Task 8: Performance Testing** (AC: #7)
  - [x] 8.1 Create `pkg/analyzer/exclusion_matcher_benchmark_test.go`:
    ```go
    func BenchmarkExclusionMatcher(b *testing.B) {
        patterns := generatePatterns(20) // 20 mixed patterns
        matcher, _ := NewExclusionMatcher(patterns)
        packages := generatePackageNames(100)

        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            for _, pkg := range packages {
                matcher.IsExcluded(pkg)
            }
        }
    }
    ```
  - [x] 8.2 Verify 100 packages × 20 patterns < 50ms

- [x] **Task 9: Integration Verification** (AC: all)
  - [x] 9.1 Build WASM: `pnpm nx build @monoguard/analysis-engine`
  - [x] 9.2 Update smoke test with exclusion examples
  - [x] 9.3 Verify excluded packages appear in graph with flag
  - [x] 9.4 Verify excluded packages don't affect metrics
  - [x] 9.5 Verify all tests pass: `make test`

## Dev Notes

### Architecture Patterns & Constraints

**From Architecture Document:**
- **Location:** Exclusion matcher in `pkg/analyzer/`
- **Pattern:** Strategy pattern for different match types
- **Integration:** Applied before all metric calculations

**Exclusion Flow:**
```
AnalysisInput.Config.Exclude
    ↓
ExclusionMatcher (parse patterns)
    ↓
GraphBuilder.buildNodes() → mark excluded=true
    ↓
CycleDetector.filterExcluded() → exclude from detection
ConflictDetector → skip excluded packages
HealthCalculator → exclude from metrics
    ↓
Result: graph has all nodes (excluded flagged), metrics exclude them
```

**Critical Constraints:**
- **camelCase JSON:** All struct tags MUST use camelCase
- **Graph Completeness:** Excluded packages MUST appear in graph (for visualization)
- **Metrics Purity:** Excluded packages MUST NOT affect any scores
- **Pattern Priority:** Exact > Glob > Regex (for performance)

### Critical Don't-Miss Rules

**From project-context.md:**

1. **Pattern Syntax:**
   ```go
   // ✅ CORRECT: Clear pattern types
   "packages/legacy"           // Exact match
   "packages/deprecated-*"     // Glob pattern
   "regex:^packages/test-.*$"  // Regex (explicit prefix)

   // ❌ WRONG: Ambiguous patterns
   "packages/*legacy*"         // Is this glob or exact? Use explicit prefix
   ```

2. **Exclusion Before Calculation:**
   ```go
   // ✅ CORRECT: Filter first, then calculate
   filteredPackages := filterExcluded(allPackages)
   cycles := detectCycles(filteredPackages)
   score := calculateScore(filteredPackages, cycles)

   // ❌ WRONG: Calculate then filter
   allCycles := detectCycles(allPackages)
   filteredCycles := filterExcludedCycles(allCycles) // May miss cycles!
   ```

3. **Graph Includes All:**
   ```go
   // ✅ CORRECT: Graph has all, flagged
   for _, pkg := range workspace.Packages {
       node := &PackageNode{
           Name:     pkg.Name,
           Excluded: matcher.IsExcluded(pkg.Name),
       }
       graph.Nodes[pkg.Name] = node
   }

   // ❌ WRONG: Graph excludes packages
   if matcher.IsExcluded(pkg.Name) {
       continue // DON'T skip - frontend needs to show grayed out
   }
   ```

### Project Structure Notes

**Target Directory Structure:**
```
packages/analysis-engine/
├── pkg/
│   ├── analyzer/
│   │   ├── analyzer.go                        # UPDATE: Accept config
│   │   ├── graph_builder.go                   # UPDATE: Use exclusion matcher
│   │   ├── cycle_detector.go                  # UPDATE: Filter excluded
│   │   ├── conflict_detector.go               # UPDATE: Skip excluded
│   │   ├── health_calculator.go               # UPDATE: Exclude from metrics
│   │   ├── exclusion_matcher.go               # NEW: Pattern matching
│   │   ├── exclusion_matcher_test.go          # NEW: Matcher tests
│   │   └── exclusion_matcher_benchmark_test.go # NEW: Performance tests
│   └── types/
│       ├── config.go                          # NEW: AnalysisConfig types
│       ├── config_test.go                     # NEW: Config tests
│       └── graph.go                           # UPDATE: Add Excluded field
├── internal/
│   └── handlers/
│       └── handlers.go                        # UPDATE: Parse AnalysisInput
└── ...
```

### Input/Output Format

**Input (AnalysisInput):**
```json
{
  "files": {
    "package.json": "{ \"workspaces\": [\"packages/*\"] }",
    "packages/core/package.json": "{ \"name\": \"@mono/core\" }",
    "packages/legacy/package.json": "{ \"name\": \"@mono/legacy\" }",
    "packages/deprecated-utils/package.json": "{ \"name\": \"@mono/deprecated-utils\" }"
  },
  "config": {
    "exclude": [
      "packages/legacy",
      "packages/deprecated-*",
      "regex:.*-test$"
    ]
  }
}
```

**Output (DependencyGraph with excluded flag):**
```json
{
  "nodes": {
    "@mono/core": {
      "name": "@mono/core",
      "excluded": false,
      "dependencies": ["@mono/legacy"]
    },
    "@mono/legacy": {
      "name": "@mono/legacy",
      "excluded": true,
      "dependencies": []
    },
    "@mono/deprecated-utils": {
      "name": "@mono/deprecated-utils",
      "excluded": true,
      "dependencies": []
    }
  },
  "edges": [
    { "from": "@mono/core", "to": "@mono/legacy", "type": "production" }
  ]
}
```

### Test Scenarios

| Pattern | Package Name | Should Match |
|---------|--------------|--------------|
| `packages/legacy` | `packages/legacy` | Yes (exact) |
| `packages/legacy` | `packages/legacy-v2` | No |
| `packages/deprecated-*` | `packages/deprecated-utils` | Yes (glob) |
| `packages/deprecated-*` | `packages/deprecated` | No (needs suffix) |
| `**/test-*` | `packages/test-utils` | Yes |
| `**/test-*` | `apps/web/test-helpers` | Yes |
| `regex:^@mono/legacy-.*$` | `@mono/legacy-v1` | Yes (regex) |
| `regex:.*-deprecated$` | `@mono/utils-deprecated` | Yes |

### Previous Story Intelligence

**From Story 2.1 (ready-for-dev):**
- Parser already has glob matching in `pkg/parser/glob.go`
- Can reuse/extend for exclusion matching

**From Story 2.2 (ready-for-dev):**
- GraphBuilder creates nodes - add exclusion flag here
- Edges connect all packages - keep edges to/from excluded

**From Story 2.3, 2.4, 2.5 (ready-for-dev):**
- All detectors need to filter excluded packages
- Health calculator needs to exclude from all metrics

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story 2.6]
- [Source: _bmad-output/implementation-artifacts/2-1-implement-workspace-configuration-parser.md#Task 3]
- [Go regexp Package](https://pkg.go.dev/regexp)
- [Glob Pattern Syntax](https://en.wikipedia.org/wiki/Glob_(programming))

## Dev Agent Record

### Agent Model Used

claude-opus-4-5-20251101

### Debug Log References

### Completion Notes List

1. **Task 1: Define AnalysisConfig Types** - Created `pkg/types/config.go` with AnalysisConfig and AnalysisInput types. Supports exclusion patterns via `exclude` array in JSON config.

2. **Task 2: Implement Pattern Matcher** - Created `pkg/analyzer/exclusion_matcher.go` with support for:
   - Exact matches: simple string comparison
   - Glob patterns: `*` (within segment), `**` (recursive), `?` (single char)
   - Regex patterns: prefixed with `regex:` for explicit identification
   - Performance: 7µs for 100 packages × 20 patterns (AC7 requires <50ms)

3. **Task 3: Update PackageNode for Exclusion Flag** - Added `Excluded bool` field to PackageNode in `pkg/types/graph.go`. Uses omitempty for clean JSON when false.

4. **Task 4: Integrate Exclusion into Graph Builder** - Updated GraphBuilder with:
   - `NewGraphBuilderWithExclusions()` constructor for exclusion patterns
   - Marks excluded packages in `buildNodes()` with `Excluded=true`

5. **Task 5: Filter Excluded from Detectors** - Created `filterExcludedPackages()` function in analyzer.go that creates a filtered graph for cycle/conflict detection and health calculation. Excluded packages are removed from nodes and their references are removed from dependency lists.

6. **Task 6: Update Analyzer to Accept Config** - Added:
   - `NewAnalyzerWithConfig()` constructor
   - `ExcludedPackages` count in AnalysisResult
   - Full graph returned for visualization (excluded packages marked)
   - Filtered graph used for all metrics

7. **Tasks 7-9: Testing and Verification** - All tests pass. Performance verified at 7µs for 100×20 (well under 50ms requirement). WASM builds successfully (4.5MB).

### File List

**New Files:**
- `packages/analysis-engine/pkg/types/config.go` - AnalysisConfig and AnalysisInput types
- `packages/analysis-engine/pkg/types/config_test.go` - Config type tests
- `packages/analysis-engine/pkg/analyzer/exclusion_matcher.go` - Pattern matching implementation
- `packages/analysis-engine/pkg/analyzer/exclusion_matcher_test.go` - Matcher unit tests
- `packages/analysis-engine/pkg/analyzer/exclusion_matcher_benchmark_test.go` - Performance benchmarks

**Modified Files:**
- `packages/analysis-engine/pkg/types/graph.go` - Added Excluded field to PackageNode
- `packages/analysis-engine/pkg/types/types.go` - Added ExcludedPackages field to AnalysisResult
- `packages/analysis-engine/pkg/analyzer/graph_builder.go` - Added exclusion matcher integration
- `packages/analysis-engine/pkg/analyzer/analyzer.go` - Added config support and filterExcludedPackages
