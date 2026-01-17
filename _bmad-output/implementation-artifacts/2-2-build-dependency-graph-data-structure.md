# Story 2.2: Build Dependency Graph Data Structure

Status: ready-for-dev

## Story

As a **developer**,
I want **the analysis engine to construct a complete dependency graph**,
So that **I can analyze relationships between all packages in the monorepo**.

## Acceptance Criteria

1. **AC1: Graph Node Creation**
   - Given parsed workspace data from Story 2.1
   - When I build the dependency graph
   - Then a node is created for each package containing:
     - Package name (unique identifier)
     - Package version
     - Relative path from workspace root
     - List of direct dependencies
     - List of dev dependencies
     - List of peer dependencies
   - And nodes are stored in a map keyed by package name for O(1) lookup

2. **AC2: Directed Edge Creation**
   - Given workspace packages with dependencies
   - When edges are constructed
   - Then an edge is created for each dependency relationship with:
     - `from` - source package name
     - `to` - target package name (dependency)
     - `type` - dependency type (production/development/peer/optional)
     - `versionRange` - version range specified (e.g., "^1.0.0")
   - And edges only connect internal workspace packages (external deps are metadata only)

3. **AC3: Dependency Type Classification**
   - Given different dependency types in package.json
   - When edges are created
   - Then types are correctly classified:
     - `dependencies` → type: "production"
     - `devDependencies` → type: "development"
     - `peerDependencies` → type: "peer"
     - `optionalDependencies` → type: "optional" (if present)

4. **AC4: Internal vs External Classification**
   - Given a mix of workspace and external dependencies
   - When the graph is built
   - Then:
     - Internal dependencies (workspace packages) create edges
     - External dependencies (npm packages) are stored in node metadata only
     - External dependencies do NOT create edges in the graph
   - And this distinction is clear in the data structure

5. **AC5: DependencyGraph Type Compliance**
   - Given the built graph
   - When serialized to JSON
   - Then the output matches the TypeScript `DependencyGraph` interface:
     ```typescript
     interface DependencyGraph {
       nodes: Record<string, PackageNode>;
       edges: DependencyEdge[];
       rootPath: string;
       workspaceType: WorkspaceType;
     }
     ```
   - And all JSON uses camelCase field names

6. **AC6: Performance Requirements**
   - Given a workspace with 100 packages
   - When graph construction completes
   - Then it finishes in < 2 seconds
   - And given a workspace with 1000 packages
   - Then memory usage is < 50MB

## Tasks / Subtasks

- [ ] **Task 1: Define Graph Types in Go** (AC: #1, #2, #3, #5)
  - [ ] 1.1 Create `pkg/types/graph.go`:
    ```go
    package types

    // DependencyGraph represents the complete dependency structure.
    // Matches @monoguard/types DependencyGraph.
    type DependencyGraph struct {
        Nodes         map[string]*PackageNode `json:"nodes"`
        Edges         []*DependencyEdge       `json:"edges"`
        RootPath      string                  `json:"rootPath"`
        WorkspaceType WorkspaceType           `json:"workspaceType"`
    }

    // PackageNode represents a single package in the graph.
    // Matches @monoguard/types PackageNode.
    type PackageNode struct {
        Name             string   `json:"name"`
        Version          string   `json:"version"`
        Path             string   `json:"path"`
        Dependencies     []string `json:"dependencies"`     // Internal workspace deps
        DevDependencies  []string `json:"devDependencies"`  // Internal workspace deps
        PeerDependencies []string `json:"peerDependencies"` // Internal workspace deps
        // External dependencies stored separately for reference
        ExternalDeps    map[string]string `json:"externalDeps,omitempty"`
        ExternalDevDeps map[string]string `json:"externalDevDeps,omitempty"`
    }

    // DependencyEdge represents a directed edge between packages.
    // Matches @monoguard/types DependencyEdge.
    type DependencyEdge struct {
        From         string         `json:"from"`
        To           string         `json:"to"`
        Type         DependencyType `json:"type"`
        VersionRange string         `json:"versionRange"`
    }

    // DependencyType classifies the dependency relationship.
    type DependencyType string

    const (
        DependencyTypeProduction  DependencyType = "production"
        DependencyTypeDevelopment DependencyType = "development"
        DependencyTypePeer        DependencyType = "peer"
        DependencyTypeOptional    DependencyType = "optional"
    )
    ```
  - [ ] 1.2 Add comprehensive tests in `pkg/types/graph_test.go`
  - [ ] 1.3 Verify JSON serialization produces exact camelCase output

- [ ] **Task 2: Implement Graph Builder** (AC: #1, #2, #4)
  - [ ] 2.1 Create `pkg/analyzer/graph_builder.go`:
    ```go
    package analyzer

    import "github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"

    // GraphBuilder constructs dependency graphs from workspace data
    type GraphBuilder struct {
        workspacePackages map[string]bool // Set of internal package names
    }

    // NewGraphBuilder creates a new graph builder
    func NewGraphBuilder() *GraphBuilder

    // Build constructs a DependencyGraph from WorkspaceData
    func (gb *GraphBuilder) Build(workspace *types.WorkspaceData) (*types.DependencyGraph, error)

    // buildNodes creates PackageNode entries for each package
    func (gb *GraphBuilder) buildNodes(workspace *types.WorkspaceData) map[string]*types.PackageNode

    // buildEdges creates DependencyEdge entries for internal dependencies
    func (gb *GraphBuilder) buildEdges(nodes map[string]*types.PackageNode) []*types.DependencyEdge

    // isInternalPackage checks if a dependency is a workspace package
    func (gb *GraphBuilder) isInternalPackage(name string) bool

    // classifyDependencies separates internal and external dependencies
    func (gb *GraphBuilder) classifyDependencies(
        allDeps map[string]string,
    ) (internal []string, external map[string]string)
    ```
  - [ ] 2.2 Implement node creation logic
  - [ ] 2.3 Implement edge creation logic with type classification
  - [ ] 2.4 Create tests in `pkg/analyzer/graph_builder_test.go`

- [ ] **Task 3: Handle Edge Cases** (AC: #4)
  - [ ] 3.1 Handle self-referencing packages (A depends on A) - should not create edge
  - [ ] 3.2 Handle missing dependencies (package referenced but not in workspace)
  - [ ] 3.3 Handle packages with no dependencies (isolated nodes)
  - [ ] 3.4 Handle duplicate dependency entries in package.json
  - [ ] 3.5 Add tests for all edge cases

- [ ] **Task 4: Wire Graph Builder to Analyzer** (AC: #5)
  - [ ] 4.1 Update `pkg/analyzer/analyzer.go`:
    ```go
    package analyzer

    import "github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"

    // Analyzer orchestrates the analysis process
    type Analyzer struct {
        graphBuilder *GraphBuilder
    }

    // NewAnalyzer creates a new analyzer
    func NewAnalyzer() *Analyzer {
        return &Analyzer{
            graphBuilder: NewGraphBuilder(),
        }
    }

    // Analyze performs complete workspace analysis
    func (a *Analyzer) Analyze(workspace *types.WorkspaceData) (*types.AnalysisResult, error) {
        // Build dependency graph
        graph, err := a.graphBuilder.Build(workspace)
        if err != nil {
            return nil, err
        }

        // Return result with graph (cycle detection comes in Story 2.3)
        return &types.AnalysisResult{
            HealthScore: 100, // Placeholder until Story 2.5
            Packages:    len(graph.Nodes),
            Graph:       graph,
        }, nil
    }
    ```
  - [ ] 4.2 Update AnalysisResult type to include Graph field
  - [ ] 4.3 Create tests in `pkg/analyzer/analyzer_test.go`

- [ ] **Task 5: Update WASM Handler** (AC: #5)
  - [ ] 5.1 Update `internal/handlers/handlers.go` to use Analyzer:
    ```go
    func HandleAnalyze(input string) *result.Result {
        // Parse input to WorkspaceData (from Story 2.1)
        var filesInput map[string]string
        if err := json.Unmarshal([]byte(input), &filesInput); err != nil {
            return result.NewError("INVALID_INPUT", err.Error())
        }

        // Convert to bytes map and parse workspace
        files := make(map[string][]byte)
        for name, content := range filesInput {
            files[name] = []byte(content)
        }

        parser := parser.NewParser("/workspace")
        workspaceData, err := parser.Parse(files)
        if err != nil {
            return result.NewError("PARSE_ERROR", err.Error())
        }

        // Run analysis (builds graph)
        analyzer := analyzer.NewAnalyzer()
        analysisResult, err := analyzer.Analyze(workspaceData)
        if err != nil {
            return result.NewError("ANALYSIS_FAILED", err.Error())
        }

        return result.NewSuccess(analysisResult)
    }
    ```
  - [ ] 5.2 Update handler tests

- [ ] **Task 6: Performance Testing** (AC: #6)
  - [ ] 6.1 Create `pkg/analyzer/benchmark_test.go`:
    ```go
    func BenchmarkBuildGraph100Packages(b *testing.B) {
        workspace := generateWorkspace(100)
        gb := NewGraphBuilder()
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            gb.Build(workspace)
        }
    }

    func BenchmarkBuildGraph1000Packages(b *testing.B) {
        workspace := generateWorkspace(1000)
        gb := NewGraphBuilder()
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            gb.Build(workspace)
        }
    }

    func generateWorkspace(packageCount int) *types.WorkspaceData {
        // Generate realistic workspace with dependencies
    }
    ```
  - [ ] 6.2 Verify 100 packages < 2 seconds
  - [ ] 6.3 Verify 1000 packages memory < 50MB with `go test -bench=. -benchmem`

- [ ] **Task 7: Integration Verification** (AC: all)
  - [ ] 7.1 Build WASM: `pnpm nx build @monoguard/analysis-engine`
  - [ ] 7.2 Update smoke test to verify graph output
  - [ ] 7.3 Test with realistic monorepo structures
  - [ ] 7.4 Verify all tests pass: `make test`

## Dev Notes

### Architecture Patterns & Constraints

**From Architecture Document:**
- **Technology:** Go 1.21+ compiled to WASM
- **Location:** Graph builder in `pkg/analyzer/`, types in `pkg/types/`
- **Pattern:** Builder pattern for graph construction
- **Constraint:** All types must match TypeScript definitions exactly

**Critical Constraints:**
- **camelCase JSON:** All struct tags MUST use camelCase
- **Internal Only Edges:** Edges connect only workspace packages, not external npm deps
- **Result Pattern:** Errors wrapped in Result type from handlers
- **Memory Efficient:** < 50MB for 1000 packages

### Critical Don't-Miss Rules

**From project-context.md:**

1. **Type Mapping to TypeScript - CRITICAL:**
   ```go
   // ✅ CORRECT: Matches TypeScript exactly
   type DependencyGraph struct {
       Nodes         map[string]*PackageNode `json:"nodes"`
       Edges         []*DependencyEdge       `json:"edges"`
       RootPath      string                  `json:"rootPath"`
       WorkspaceType WorkspaceType           `json:"workspaceType"`
   }

   // ❌ WRONG: Different field names/structure
   type DependencyGraph struct {
       Packages []Package `json:"packages"` // TypeScript expects "nodes"
   }
   ```

2. **Internal vs External Dependencies:**
   ```go
   // ✅ CORRECT: Only internal deps create edges
   if gb.isInternalPackage(depName) {
       edges = append(edges, &DependencyEdge{
           From: pkg.Name,
           To:   depName,
           Type: DependencyTypeProduction,
       })
       node.Dependencies = append(node.Dependencies, depName)
   } else {
       node.ExternalDeps[depName] = version
   }

   // ❌ WRONG: Creating edges for external packages
   edges = append(edges, &DependencyEdge{
       From: pkg.Name,
       To:   "lodash", // External package - should NOT be an edge
   })
   ```

3. **Go Naming Conventions:**
   - PascalCase for exported: `GraphBuilder`, `Build`, `DependencyGraph`
   - camelCase for unexported: `buildNodes`, `isInternalPackage`
   - snake_case for files: `graph_builder.go`, `graph_builder_test.go`

### Project Structure Notes

**Target Directory Structure:**
```
packages/analysis-engine/
├── pkg/
│   ├── analyzer/
│   │   ├── analyzer.go            # UPDATE: Orchestrates analysis
│   │   ├── analyzer_test.go       # UPDATE: Analyzer tests
│   │   ├── graph_builder.go       # NEW: Graph construction
│   │   ├── graph_builder_test.go  # NEW: Graph builder tests
│   │   └── benchmark_test.go      # NEW: Performance benchmarks
│   ├── parser/
│   │   └── ... (from Story 2.1)
│   └── types/
│       ├── types.go               # UPDATE: Add Graph field to AnalysisResult
│       ├── types_test.go
│       ├── graph.go               # NEW: DependencyGraph types
│       └── graph_test.go          # NEW: Graph type tests
├── internal/
│   └── handlers/
│       ├── handlers.go            # UPDATE: Use Analyzer
│       └── handlers_test.go       # UPDATE: Tests
└── ...
```

### Input/Output Format

**Input (WorkspaceData from Story 2.1):**
```json
{
  "rootPath": "/workspace",
  "workspaceType": "pnpm",
  "packages": {
    "@mono/app": {
      "name": "@mono/app",
      "version": "1.0.0",
      "path": "apps/web",
      "dependencies": { "@mono/ui": "^1.0.0", "react": "^18.0.0" },
      "devDependencies": { "@mono/types": "^1.0.0" },
      "peerDependencies": {}
    },
    "@mono/ui": {
      "name": "@mono/ui",
      "version": "1.0.0",
      "path": "packages/ui",
      "dependencies": { "react": "^18.0.0" },
      "devDependencies": {},
      "peerDependencies": {}
    },
    "@mono/types": {
      "name": "@mono/types",
      "version": "1.0.0",
      "path": "packages/types",
      "dependencies": {},
      "devDependencies": {},
      "peerDependencies": {}
    }
  }
}
```

**Output (DependencyGraph):**
```json
{
  "nodes": {
    "@mono/app": {
      "name": "@mono/app",
      "version": "1.0.0",
      "path": "apps/web",
      "dependencies": ["@mono/ui"],
      "devDependencies": ["@mono/types"],
      "peerDependencies": [],
      "externalDeps": { "react": "^18.0.0" }
    },
    "@mono/ui": {
      "name": "@mono/ui",
      "version": "1.0.0",
      "path": "packages/ui",
      "dependencies": [],
      "devDependencies": [],
      "peerDependencies": [],
      "externalDeps": { "react": "^18.0.0" }
    },
    "@mono/types": {
      "name": "@mono/types",
      "version": "1.0.0",
      "path": "packages/types",
      "dependencies": [],
      "devDependencies": [],
      "peerDependencies": []
    }
  },
  "edges": [
    { "from": "@mono/app", "to": "@mono/ui", "type": "production", "versionRange": "^1.0.0" },
    { "from": "@mono/app", "to": "@mono/types", "type": "development", "versionRange": "^1.0.0" }
  ],
  "rootPath": "/workspace",
  "workspaceType": "pnpm"
}
```

### Previous Story Intelligence

**From Story 2.1 (ready-for-dev):**
- WorkspaceData type defined with packages map
- Package has `Dependencies map[string]string` (name → version range)
- Parser returns WorkspaceData which is input to this story
- Note: 2.1 Package.Dependencies is `map[string]string`, but graph PackageNode.Dependencies is `[]string` (names only, internal only)

**Key Difference:**
- Story 2.1 `Package.Dependencies`: ALL deps with versions → `map[string]string`
- Story 2.2 `PackageNode.Dependencies`: Internal deps only, names → `[]string`
- Story 2.2 `PackageNode.ExternalDeps`: External deps with versions → `map[string]string`

**From Story 1.5 (TypeScript Types):**
- `DependencyGraph` interface defined in `@monoguard/types`
- `PackageNode.dependencies` is `string[]` (array of names)
- This is a deliberate design: edges show relationships, versions are metadata

### Graph Theory Concepts

**Directed Graph (Digraph):**
- Nodes: Packages in the workspace
- Edges: Dependency relationships (A depends on B = edge from A to B)
- Edge direction: FROM dependent TO dependency

**Why Only Internal Edges:**
- External packages (lodash, react) are not in our workspace
- We can't analyze their internal dependencies
- They represent "leaves" in our analysis scope
- Circular detection only matters for code we control

### Testing Requirements

**Unit Tests:**
- Graph builder with various workspace configs
- Node creation with all dependency types
- Edge classification (production/development/peer)
- Internal vs external separation
- Edge cases: self-reference, missing deps, empty workspace

**Test Cases to Cover:**
| Scenario | Expected Behavior |
|----------|-------------------|
| Empty workspace | Empty graph, no edges |
| Single package, no deps | One node, no edges |
| Two packages, A→B | Two nodes, one edge |
| Circular A↔B | Two nodes, two edges (cycle detection is Story 2.3) |
| External deps only | Nodes created, no edges, external in metadata |
| Mixed internal/external | Edges for internal, metadata for external |
| Self-reference A→A | Node created, NO edge (skip self-loops) |
| Dev dependency | Edge with type "development" |
| Peer dependency | Edge with type "peer" |

**Performance Benchmarks:**
- Build graph 100 packages: < 2 seconds
- Build graph 1000 packages: < 10 seconds, < 50MB memory

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story 2.2]
- [Source: _bmad-output/project-context.md#Go Naming Conventions]
- [Source: packages/types/src/analysis/graph.ts]
- [Source: _bmad-output/implementation-artifacts/2-1-implement-workspace-configuration-parser.md]
- [Graph Theory - Directed Graphs](https://en.wikipedia.org/wiki/Directed_graph)

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
