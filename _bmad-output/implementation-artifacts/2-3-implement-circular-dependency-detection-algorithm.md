# Story 2.3: Implement Circular Dependency Detection Algorithm

Status: done

## Story

As a **user**,
I want **to identify all circular dependencies in my monorepo**,
So that **I know which packages have problematic dependency relationships**.

## Acceptance Criteria

1. **AC1: Cycle Detection Algorithm**
   - Given a dependency graph from Story 2.2
   - When I run circular dependency detection
   - Then the algorithm detects all cycles using Tarjan's SCC algorithm or DFS-based approach
   - And handles both simple cycles (A → B → A) and complex multi-package cycles (A → B → C → D → A)

2. **AC2: CircularDependency Output**
   - Given detected cycles
   - When results are returned
   - Then each cycle includes:
     - `cycle` - array of package names in order (e.g., ["A", "B", "C", "A"])
     - `type` - "direct" (2 packages) or "indirect" (3+ packages)
     - `severity` - classification based on cycle characteristics
     - `depth` - number of packages in the cycle
   - And the cycle representation starts and ends with the same package

3. **AC3: Shortest Cycle Representation**
   - Given overlapping cycles
   - When cycles are reported
   - Then each cycle is reported in its shortest/canonical form
   - And duplicate cycles are not reported (A→B→C→A is same as B→C→A→B)

4. **AC4: Nested Cycle Handling**
   - Given complex dependency structures with nested cycles
   - When detection runs
   - Then all distinct cycles are identified:
     - Inner cycles (A → B → A)
     - Outer cycles that contain inner cycles (A → B → C → A where B → D → B exists)
   - And cycles are reported independently

5. **AC5: Self-Loop Detection**
   - Given a package that depends on itself
   - When detection runs
   - Then self-loops are detected and reported as cycles with depth 1
   - And they are marked with severity "critical"

6. **AC6: Performance Requirements**
   - Given a workspace with 100 packages and up to 5 cycles
   - When detection completes
   - Then it finishes in < 3 seconds
   - And given 1000 packages, completes in < 30 seconds

## Tasks / Subtasks

- [x] **Task 1: Expand CircularDependency Type** (AC: #2)
  - [x] 1.1 Update `pkg/types/circular.go` (create new file):
    ```go
    package types

    // CircularDependencyInfo represents a detected circular dependency.
    // Matches @monoguard/types CircularDependencyInfo.
    type CircularDependencyInfo struct {
        Cycle      []string           `json:"cycle"`      // Package names in order, ends with first
        Type       CircularType       `json:"type"`       // direct or indirect
        Severity   CircularSeverity   `json:"severity"`   // critical, warning, or info
        Depth      int                `json:"depth"`      // Number of unique packages
        Impact     string             `json:"impact"`     // Human-readable impact description
        Complexity int                `json:"complexity"` // Refactoring complexity (1-10)
    }

    // CircularType classifies the cycle length
    type CircularType string

    const (
        CircularTypeDirect   CircularType = "direct"   // 2 packages: A ↔ B
        CircularTypeIndirect CircularType = "indirect" // 3+ packages: A → B → C → A
    )

    // CircularSeverity indicates how problematic the cycle is
    type CircularSeverity string

    const (
        CircularSeverityCritical CircularSeverity = "critical" // Self-loop or blocking build
        CircularSeverityWarning  CircularSeverity = "warning"  // Should be fixed
        CircularSeverityInfo     CircularSeverity = "info"     // Nice to fix
    )
    ```
  - [x] 1.2 Add JSON serialization tests
  - [x] 1.3 Remove old CircularDependency from types.go (migrate to new type)

- [x] **Task 2: Implement Tarjan's SCC Algorithm** (AC: #1, #3, #4)
  - [x] 2.1 Create `pkg/analyzer/cycle_detector.go`:
    ```go
    package analyzer

    import "github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"

    // CycleDetector finds circular dependencies in a graph
    type CycleDetector struct {
        graph      *types.DependencyGraph
        index      int
        stack      []string
        onStack    map[string]bool
        indices    map[string]int
        lowLinks   map[string]int
        components [][]string // Strongly connected components
    }

    // NewCycleDetector creates a new detector for the given graph
    func NewCycleDetector(graph *types.DependencyGraph) *CycleDetector

    // DetectCycles finds all circular dependencies
    func (cd *CycleDetector) DetectCycles() []*types.CircularDependencyInfo

    // tarjanSCC implements Tarjan's strongly connected components algorithm
    func (cd *CycleDetector) tarjanSCC()

    // strongConnect is the recursive part of Tarjan's algorithm
    func (cd *CycleDetector) strongConnect(node string)

    // extractCycles converts SCCs to CircularDependencyInfo
    func (cd *CycleDetector) extractCycles(sccs [][]string) []*types.CircularDependencyInfo

    // normalizeCycle converts a cycle to its canonical form
    func normalizeCycle(cycle []string) []string

    // classifySeverity determines the severity of a cycle
    func classifySeverity(cycle []string, cycleType CircularType) CircularSeverity

    // calculateComplexity estimates refactoring effort (1-10)
    func calculateComplexity(cycle []string, graph *types.DependencyGraph) int
    ```
  - [x] 2.2 Implement Tarjan's algorithm with proper index tracking
  - [x] 2.3 Handle disconnected graph components
  - [x] 2.4 Create comprehensive tests in `pkg/analyzer/cycle_detector_test.go`

- [x] **Task 3: Implement Cycle Extraction and Classification** (AC: #2, #3)
  - [x] 3.1 Implement `extractCycles`:
    ```go
    func (cd *CycleDetector) extractCycles(sccs [][]string) []*types.CircularDependencyInfo {
        var cycles []*types.CircularDependencyInfo

        for _, scc := range sccs {
            // Only SCCs with 2+ nodes are cycles
            // (SCC with 1 node is a cycle only if it has self-edge)
            if len(scc) > 1 {
                cycle := cd.buildCyclePath(scc)
                normalized := normalizeCycle(cycle)

                cycleType := types.CircularTypeIndirect
                if len(scc) == 2 {
                    cycleType = types.CircularTypeDirect
                }

                cycles = append(cycles, &types.CircularDependencyInfo{
                    Cycle:      normalized,
                    Type:       cycleType,
                    Severity:   classifySeverity(normalized, cycleType),
                    Depth:      len(scc),
                    Impact:     generateImpactDescription(scc),
                    Complexity: calculateComplexity(scc, cd.graph),
                })
            } else if len(scc) == 1 {
                // Check for self-loop
                if cd.hasSelfLoop(scc[0]) {
                    cycles = append(cycles, &types.CircularDependencyInfo{
                        Cycle:      []string{scc[0], scc[0]},
                        Type:       types.CircularTypeDirect,
                        Severity:   types.CircularSeverityCritical,
                        Depth:      1,
                        Impact:     fmt.Sprintf("Package %s depends on itself", scc[0]),
                        Complexity: 1,
                    })
                }
            }
        }

        return cycles
    }
    ```
  - [x] 3.2 Implement cycle path reconstruction from SCC
  - [x] 3.3 Implement `normalizeCycle` for canonical form (rotate to start with smallest name)
  - [x] 3.4 Add tests for cycle extraction

- [x] **Task 4: Handle Self-Loops** (AC: #5)
  - [x] 4.1 Add self-loop detection in graph building (from 2.2) or detection phase:
    ```go
    func (cd *CycleDetector) hasSelfLoop(node string) bool {
        pkg, exists := cd.graph.Nodes[node]
        if !exists {
            return false
        }
        for _, dep := range pkg.Dependencies {
            if dep == node {
                return true
            }
        }
        // Also check devDependencies and peerDependencies
        return false
    }
    ```
  - [x] 4.2 Add tests for self-loop scenarios

- [x] **Task 5: Implement Impact and Complexity Scoring** (AC: #2)
  - [x] 5.1 Implement `generateImpactDescription`:
    ```go
    func generateImpactDescription(cycle []string) string {
        if len(cycle) == 1 {
            return fmt.Sprintf("Self-referencing package: %s", cycle[0])
        }
        if len(cycle) == 2 {
            return fmt.Sprintf("Direct circular dependency between %s and %s", cycle[0], cycle[1])
        }
        return fmt.Sprintf("Indirect circular dependency involving %d packages: %s",
            len(cycle), strings.Join(cycle, " → "))
    }
    ```
  - [x] 5.2 Implement `calculateComplexity` based on:
    - Cycle length (longer = more complex)
    - Number of edges in cycle
    - Whether packages are heavily depended upon
  - [x] 5.3 Add tests for impact and complexity

- [x] **Task 6: Wire to Analyzer** (AC: all)
  - [x] 6.1 Update `pkg/analyzer/analyzer.go`:
    ```go
    func (a *Analyzer) Analyze(workspace *types.WorkspaceData) (*types.AnalysisResult, error) {
        // Build dependency graph (Story 2.2)
        graph, err := a.graphBuilder.Build(workspace)
        if err != nil {
            return nil, err
        }

        // Detect circular dependencies (Story 2.3)
        detector := NewCycleDetector(graph)
        cycles := detector.DetectCycles()

        return &types.AnalysisResult{
            HealthScore:          100, // Will be calculated in Story 2.5
            Packages:             len(graph.Nodes),
            CircularDependencies: cycles,
            Graph:                graph,
            CreatedAt:            time.Now().UTC().Format(time.RFC3339),
        }, nil
    }
    ```
  - [x] 6.2 Update AnalysisResult type to include CircularDependencies field
  - [x] 6.3 Update handler and WASM tests

- [x] **Task 7: Performance Testing** (AC: #6)
  - [x] 7.1 Create `pkg/analyzer/cycle_detector_benchmark_test.go`:
    ```go
    func BenchmarkDetectCycles100Packages(b *testing.B) {
        graph := generateGraphWithCycles(100, 5)
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            detector := NewCycleDetector(graph)
            detector.DetectCycles()
        }
    }

    func BenchmarkDetectCycles1000Packages(b *testing.B) {
        graph := generateGraphWithCycles(1000, 10)
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            detector := NewCycleDetector(graph)
            detector.DetectCycles()
        }
    }

    func generateGraphWithCycles(packageCount, cycleCount int) *types.DependencyGraph {
        // Generate realistic graph with specified number of cycles
    }
    ```
  - [x] 7.2 Verify 100 packages < 3 seconds (Actual: 0.1ms)
  - [x] 7.3 Verify 1000 packages < 30 seconds (Actual: 1.4ms)

- [x] **Task 8: Integration Verification** (AC: all)
  - [x] 8.1 Build WASM: `pnpm nx build @monoguard/analysis-engine` (4.4MB)
  - [x] 8.2 Update smoke test to verify cycle detection
  - [x] 8.3 Test with known cycle scenarios
  - [x] 8.4 Verify all tests pass: `make test`

## Dev Notes

### Architecture Patterns & Constraints

**From Architecture Document:**
- **Algorithm:** Tarjan's Strongly Connected Components (O(V+E) complexity)
- **Location:** Cycle detector in `pkg/analyzer/`
- **Pattern:** Detector pattern with graph input
- **Output:** List of CircularDependencyInfo matching TypeScript types

**Why Tarjan's Algorithm:**
- Finds ALL strongly connected components in one pass
- O(V+E) time complexity - optimal for graph traversal
- Each SCC with 2+ nodes represents a cycle
- Well-documented and proven algorithm

**Critical Constraints:**
- **camelCase JSON:** All struct tags MUST use camelCase
- **Cycle Format:** Array starts and ends with same package
- **No Duplicates:** Same cycle rotated should not appear twice
- **Self-Loops:** Must be detected and reported

### Critical Don't-Miss Rules

**From project-context.md:**

1. **Cycle Representation:**
   ```go
   // ✅ CORRECT: Cycle ends with starting node
   cycle := []string{"A", "B", "C", "A"}

   // ❌ WRONG: Doesn't show the loop back
   cycle := []string{"A", "B", "C"}
   ```

2. **Canonical Form:**
   ```go
   // All these represent the same cycle:
   // A → B → C → A
   // B → C → A → B
   // C → A → B → C

   // ✅ CORRECT: Normalize to start with lexicographically smallest
   normalized := []string{"A", "B", "C", "A"} // A is smallest

   // Function to normalize:
   func normalizeCycle(cycle []string) []string {
       if len(cycle) <= 1 {
           return cycle
       }
       // Remove trailing duplicate
       nodes := cycle[:len(cycle)-1]

       // Find index of smallest node
       minIdx := 0
       for i, node := range nodes {
           if node < nodes[minIdx] {
               minIdx = i
           }
       }

       // Rotate to start with smallest
       result := make([]string, len(nodes)+1)
       for i := 0; i < len(nodes); i++ {
           result[i] = nodes[(minIdx+i)%len(nodes)]
       }
       result[len(nodes)] = result[0] // Close the cycle

       return result
   }
   ```

3. **Severity Classification:**
   - `critical`: Self-loops, cycles blocking builds
   - `warning`: Direct cycles (A ↔ B) - should fix
   - `info`: Indirect cycles with 3+ packages - nice to fix

### Project Structure Notes

**Target Directory Structure:**
```
packages/analysis-engine/
├── pkg/
│   ├── analyzer/
│   │   ├── analyzer.go                      # UPDATE: Add cycle detection call
│   │   ├── analyzer_test.go                 # UPDATE
│   │   ├── graph_builder.go                 # From Story 2.2
│   │   ├── graph_builder_test.go
│   │   ├── cycle_detector.go                # NEW: Tarjan's algorithm
│   │   ├── cycle_detector_test.go           # NEW: Detection tests
│   │   └── cycle_detector_benchmark_test.go # NEW: Performance tests
│   └── types/
│       ├── types.go                         # UPDATE: Add CircularDependencies to AnalysisResult
│       ├── circular.go                      # NEW: CircularDependencyInfo types
│       ├── circular_test.go                 # NEW: Type tests
│       └── graph.go                         # From Story 2.2
└── ...
```

### Algorithm: Tarjan's SCC

**Pseudocode:**
```
algorithm tarjan is
    input: graph G = (V, E)
    output: set of strongly connected components

    index := 0
    S := empty stack
    for each v in V do
        if v.index is undefined then
            strongconnect(v)

    function strongconnect(v)
        v.index := index
        v.lowlink := index
        index := index + 1
        S.push(v)
        v.onStack := true

        for each (v, w) in E do
            if w.index is undefined then
                strongconnect(w)
                v.lowlink := min(v.lowlink, w.lowlink)
            else if w.onStack then
                v.lowlink := min(v.lowlink, w.index)

        if v.lowlink = v.index then
            start a new SCC
            repeat
                w := S.pop()
                w.onStack := false
                add w to current SCC
            while w ≠ v
            output the current SCC
```

**Key Insight:** A strongly connected component with more than one node is a cycle. The algorithm finds ALL SCCs efficiently in O(V+E) time.

### Input/Output Format

**Input (DependencyGraph from Story 2.2):**
```json
{
  "nodes": {
    "A": { "name": "A", "dependencies": ["B"] },
    "B": { "name": "B", "dependencies": ["C"] },
    "C": { "name": "C", "dependencies": ["A"] }
  },
  "edges": [
    { "from": "A", "to": "B", "type": "production" },
    { "from": "B", "to": "C", "type": "production" },
    { "from": "C", "to": "A", "type": "production" }
  ]
}
```

**Output (CircularDependencyInfo[]):**
```json
[
  {
    "cycle": ["A", "B", "C", "A"],
    "type": "indirect",
    "severity": "warning",
    "depth": 3,
    "impact": "Indirect circular dependency involving 3 packages: A → B → C → A",
    "complexity": 5
  }
]
```

### Test Scenarios

| Scenario | Graph | Expected Cycles |
|----------|-------|-----------------|
| No cycles | A→B→C | [] |
| Direct cycle | A↔B | [["A", "B", "A"]] |
| Indirect cycle | A→B→C→A | [["A", "B", "C", "A"]] |
| Self-loop | A→A | [["A", "A"]] |
| Multiple cycles | A↔B, C→D→E→C | 2 cycles |
| Nested cycles | A→B→A, A→B→C→A | 2 distinct cycles |
| Disconnected | A→B, C→D→C | 1 cycle (C→D→C) |
| Complex graph | Many interconnections | All unique cycles |

### Previous Story Intelligence

**From Story 2.2 (ready-for-dev):**
- DependencyGraph has `Nodes` map and `Edges` array
- PackageNode has `Dependencies []string` (internal only)
- Graph builder separates internal/external dependencies
- Edges only exist for internal workspace dependencies

**Key Data Available:**
- `graph.Nodes[name].Dependencies` - internal deps for traversal
- `graph.Edges` - can also be used but Nodes is sufficient

### TypeScript Type Reference

Must match `@monoguard/types`:
```typescript
interface CircularDependencyInfo {
  cycle: string[]           // ["A", "B", "C", "A"]
  type: 'direct' | 'indirect'
  severity: 'critical' | 'warning' | 'info'
  impact: string
  fixStrategy?: FixStrategy  // NOT in this story - Epic 3
  complexity: number         // 1-10
}
```

Note: `fixStrategy` is optional and will be implemented in Epic 3 (Stories 3.1-3.8).

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story 2.3]
- [Source: _bmad-output/project-context.md#Result Type Pattern]
- [Source: packages/types/src/analysis/results.ts#CircularDependencyInfo]
- [Source: _bmad-output/implementation-artifacts/2-2-build-dependency-graph-data-structure.md]
- [Tarjan's Algorithm - Wikipedia](https://en.wikipedia.org/wiki/Tarjan%27s_strongly_connected_components_algorithm)
- [SCC Applications in Dependency Analysis](https://cs.stackexchange.com/questions/tagged/strongly-connected-components)

## Dev Agent Record

### Agent Model Used

Claude Opus 4.5 (claude-opus-4-5-20251101)

### Debug Log References

N/A

### Completion Notes List

1. **Task 1:** Created `pkg/types/circular.go` with CircularDependencyInfo, CircularType, CircularSeverity types matching TypeScript interfaces exactly. All JSON tags use camelCase. Added comprehensive tests in `circular_test.go`.

2. **Task 2:** Implemented Tarjan's strongly connected components algorithm in `pkg/analyzer/cycle_detector.go`:
   - Uses O(V+E) time complexity for optimal performance
   - Handles disconnected graph components correctly
   - Added adjacency list builder for graph traversal

3. **Task 3:** Implemented cycle extraction and classification:
   - `extractCycles()` converts SCCs to CircularDependencyInfo
   - `normalizeCycle()` ensures canonical form (starts with lexicographically smallest node)
   - Deduplication prevents same cycle reported multiple times

4. **Task 4:** Implemented self-loop detection:
   - `hasSelfLoop()` checks all dependency types (deps, devDeps, peerDeps, optionalDeps)
   - Self-loops marked as severity "critical"

5. **Task 5:** Implemented impact and complexity scoring:
   - `generateImpactDescription()` creates human-readable descriptions
   - `calculateBaseComplexity()` scores 1-10 based on cycle depth

6. **Task 6:** Wired CycleDetector to Analyzer:
   - Updated `analyzer.go` to call `NewCycleDetector().DetectCycles()`
   - Added `CircularDependencies` field to AnalysisResult
   - Added `CreatedAt` ISO 8601 timestamp

7. **Task 7:** Performance benchmarks:
   - 100 packages with 5 cycles: 0.1ms (requirement: < 3s) ✅
   - 1000 packages with 10 cycles: 1.4ms (requirement: < 30s) ✅

8. **Task 8:** Integration verification:
   - WASM builds successfully (4.4MB)
   - All tests pass (100+)

### File List

**New Files:**
- `packages/analysis-engine/pkg/types/circular.go` - CircularDependencyInfo types
- `packages/analysis-engine/pkg/types/circular_test.go` - Type tests (15 tests)
- `packages/analysis-engine/pkg/analyzer/cycle_detector.go` - Tarjan's SCC algorithm
- `packages/analysis-engine/pkg/analyzer/cycle_detector_test.go` - Detection tests (20+ tests)
- `packages/analysis-engine/pkg/analyzer/cycle_detector_benchmark_test.go` - Performance benchmarks

**Modified Files:**
- `packages/analysis-engine/pkg/types/types.go` - Added CircularDependencies to AnalysisResult
- `packages/analysis-engine/pkg/analyzer/analyzer.go` - Added cycle detection call

### Change Log

- 2026-01-17: Implemented circular dependency detection using Tarjan's SCC algorithm (Story 2.3)
- 2026-01-17: Code Review Fixes Applied (7 issues resolved):
  - HIGH: Added missing `depth` field to TypeScript CircularDependencyInfo
  - MEDIUM: Fixed TypeScript AnalysisResult field name (`packageCount` → `packages`)
  - MEDIUM: Added TypeScript types for VersionConflictInfo, HealthScoreDetails matching Go
  - MEDIUM: Added optionalDependencies self-loop test coverage
  - LOW: Removed deprecated CircularDependency and VersionConflict types from types.go
  - LOW: Added named constants for complexity calculation (ComplexitySelfLoop, etc.)
  - LOW: Optimized cyclesToKey() using strings.Join for better performance
