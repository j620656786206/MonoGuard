# Story 3.3: Generate Fix Strategy Recommendations

Status: done

## Story

As a **user**,
I want **to receive recommended fix strategies for each circular dependency**,
So that **I have actionable options to resolve the problem**.

## Acceptance Criteria

1. **AC1: Three Fix Strategy Types**
   - Given a circular dependency with root cause analysis (from Story 3.1)
   - When I request fix recommendations
   - Then I receive up to 3 strategies:
     1. **Extract Shared Module** - Move shared code to new package
     2. **Dependency Injection** - Invert the dependency relationship
     3. **Module Boundary Refactoring** - Restructure module boundaries
   - And strategies are returned as `FixStrategy[]`

2. **AC2: Suitability Scoring**
   - Given each fix strategy
   - When calculating suitability
   - Then the score (1-10) is based on:
     - Cycle characteristics (direct vs indirect, depth)
     - Dependency types involved (production vs dev)
     - Package naming patterns (core/shared packages)
     - Critical edge location (from Story 3.1)
   - And higher scores indicate better fit for this specific cycle

3. **AC3: Effort Estimation**
   - Given each fix strategy
   - When estimating effort
   - Then classify as:
     - **Low** - Simple changes, < 1 hour
     - **Medium** - Moderate changes, 1-4 hours
     - **High** - Significant refactoring, > 4 hours
   - And effort is based on cycle depth and strategy complexity

4. **AC4: Pros and Cons**
   - Given each fix strategy for a specific cycle
   - When generating recommendations
   - Then include:
     - 2-3 pros specific to this cycle
     - 2-3 cons specific to this cycle
   - And pros/cons are contextual (not generic)

5. **AC5: Strategy Ranking**
   - Given multiple applicable strategies
   - When presenting recommendations
   - Then strategies are ranked by:
     1. Suitability score (highest first)
     2. Effort (lower effort wins ties)
   - And the top recommendation is clearly marked

6. **AC6: Integration with CircularDependencyInfo**
   - Given analysis results
   - When enriching CircularDependencyInfo
   - Then add optional `fixStrategies` field:
     ```go
     type CircularDependencyInfo struct {
         // ... existing fields ...
         FixStrategies []FixStrategy `json:"fixStrategies,omitempty"`
     }
     ```
   - And existing consumers continue working (backward compatible)

7. **AC7: Performance**
   - Given a workspace with 100 packages and 5 cycles
   - When generating fix strategies
   - Then generation completes in < 200ms additional overhead
   - And memory usage increase is < 5MB

## Tasks / Subtasks

- [x] **Task 1: Define FixStrategy Types** (AC: #1, #3, #4)
  - [x] 1.1 Create `pkg/types/fix_strategy.go`:
    ```go
    package types

    // FixStrategy represents a recommended approach to resolve a circular dependency.
    // Matches @monoguard/types FixStrategy interface.
    type FixStrategy struct {
        // Type identifies the strategy approach
        Type FixStrategyType `json:"type"`

        // Name is the human-readable strategy name
        Name string `json:"name"`

        // Description explains what this strategy does
        Description string `json:"description"`

        // Suitability is a score (1-10) indicating how well this strategy fits
        Suitability int `json:"suitability"`

        // Effort estimates the implementation difficulty
        Effort EffortLevel `json:"effort"`

        // Pros are advantages of this strategy for this specific cycle
        Pros []string `json:"pros"`

        // Cons are disadvantages of this strategy for this specific cycle
        Cons []string `json:"cons"`

        // Recommended indicates this is the top recommendation
        Recommended bool `json:"recommended"`

        // TargetPackages are the packages that would need modification
        TargetPackages []string `json:"targetPackages"`

        // NewPackageName is suggested name for extracted module (if applicable)
        NewPackageName string `json:"newPackageName,omitempty"`
    }

    // FixStrategyType identifies the approach to resolve the cycle.
    type FixStrategyType string

    const (
        FixStrategyExtractModule    FixStrategyType = "extract-module"
        FixStrategyDependencyInject FixStrategyType = "dependency-injection"
        FixStrategyBoundaryRefactor FixStrategyType = "boundary-refactoring"
    )

    // EffortLevel estimates implementation difficulty.
    type EffortLevel string

    const (
        EffortLow    EffortLevel = "low"    // < 1 hour
        EffortMedium EffortLevel = "medium" // 1-4 hours
        EffortHigh   EffortLevel = "high"   // > 4 hours
    )
    ```
  - [x] 1.2 Add JSON serialization tests in `pkg/types/fix_strategy_test.go`
  - [x] 1.3 Ensure all JSON tags use camelCase

- [x] **Task 2: Create FixStrategyGenerator** (AC: #1, #2)
  - [x] 2.1 Create `pkg/analyzer/fix_strategy_generator.go`:
    ```go
    package analyzer

    import "github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"

    // FixStrategyGenerator creates fix recommendations for circular dependencies.
    type FixStrategyGenerator struct {
        graph     *types.DependencyGraph
        workspace *types.WorkspaceData
    }

    // NewFixStrategyGenerator creates a new generator.
    func NewFixStrategyGenerator(graph *types.DependencyGraph, workspace *types.WorkspaceData) *FixStrategyGenerator

    // Generate creates fix strategies for a circular dependency.
    // Requires RootCause to be populated (from Story 3.1).
    func (fsg *FixStrategyGenerator) Generate(cycle *types.CircularDependencyInfo) []types.FixStrategy

    // generateExtractModule creates the Extract Shared Module strategy.
    func (fsg *FixStrategyGenerator) generateExtractModule(cycle *types.CircularDependencyInfo) *types.FixStrategy

    // generateDependencyInjection creates the Dependency Injection strategy.
    func (fsg *FixStrategyGenerator) generateDependencyInjection(cycle *types.CircularDependencyInfo) *types.FixStrategy

    // generateBoundaryRefactoring creates the Module Boundary Refactoring strategy.
    func (fsg *FixStrategyGenerator) generateBoundaryRefactoring(cycle *types.CircularDependencyInfo) *types.FixStrategy

    // rankStrategies sorts strategies by suitability and effort.
    func rankStrategies(strategies []types.FixStrategy) []types.FixStrategy
    ```
  - [x] 2.2 Implement strategy generation logic
  - [x] 2.3 Create comprehensive tests in `pkg/analyzer/fix_strategy_generator_test.go`

- [x] **Task 3: Implement Extract Shared Module Strategy** (AC: #1, #2, #3, #4)
  - [x] 3.1 Implement `generateExtractModule`:
    ```go
    func (fsg *FixStrategyGenerator) generateExtractModule(cycle *types.CircularDependencyInfo) *types.FixStrategy {
        // Calculate suitability based on:
        // - Higher for cycles with shared dependencies
        // - Higher for indirect cycles (3+ packages)
        // - Lower for direct cycles (DI often better)
        suitability := fsg.calculateExtractModuleSuitability(cycle)

        // Estimate effort based on:
        // - Cycle depth (more packages = more effort)
        // - Number of shared imports (more = more effort)
        effort := fsg.calculateExtractModuleEffort(cycle)

        // Generate contextual pros/cons
        pros, cons := fsg.generateExtractModuleProsCons(cycle)

        // Suggest new package name
        newPkgName := fsg.suggestNewPackageName(cycle)

        return &types.FixStrategy{
            Type:           types.FixStrategyExtractModule,
            Name:           "Extract Shared Module",
            Description:    "Create a new shared package to hold common dependencies, breaking the cycle.",
            Suitability:    suitability,
            Effort:         effort,
            Pros:           pros,
            Cons:           cons,
            TargetPackages: fsg.getTargetPackages(cycle),
            NewPackageName: newPkgName,
        }
    }
    ```
  - [x] 3.2 Implement suitability calculation
  - [x] 3.3 Implement effort estimation
  - [x] 3.4 Implement contextual pros/cons generation
  - [x] 3.5 Add tests for Extract Module strategy

- [x] **Task 4: Implement Dependency Injection Strategy** (AC: #1, #2, #3, #4)
  - [x] 4.1 Implement `generateDependencyInjection`:
    ```go
    func (fsg *FixStrategyGenerator) generateDependencyInjection(cycle *types.CircularDependencyInfo) *types.FixStrategy {
        // Calculate suitability based on:
        // - Higher for direct cycles (A ↔ B)
        // - Higher when critical edge is clear
        // - Lower for deeply nested cycles
        suitability := fsg.calculateDISuitability(cycle)

        // Effort is typically medium for DI
        effort := fsg.calculateDIEffort(cycle)

        // Generate contextual pros/cons
        pros, cons := fsg.generateDIProsCons(cycle)

        return &types.FixStrategy{
            Type:           types.FixStrategyDependencyInject,
            Name:           "Dependency Injection",
            Description:    "Invert the problematic dependency by introducing an interface or callback pattern.",
            Suitability:    suitability,
            Effort:         effort,
            Pros:           pros,
            Cons:           cons,
            TargetPackages: fsg.getDITargetPackages(cycle),
        }
    }
    ```
  - [x] 4.2 Implement suitability calculation
  - [x] 4.3 Implement effort estimation
  - [x] 4.4 Implement contextual pros/cons generation
  - [x] 4.5 Add tests for Dependency Injection strategy

- [x] **Task 5: Implement Module Boundary Refactoring Strategy** (AC: #1, #2, #3, #4)
  - [x] 5.1 Implement `generateBoundaryRefactoring`:
    ```go
    func (fsg *FixStrategyGenerator) generateBoundaryRefactoring(cycle *types.CircularDependencyInfo) *types.FixStrategy {
        // Calculate suitability based on:
        // - Higher when packages have overlapping responsibilities
        // - Higher for cycles involving "core" or "common" packages
        // - Often the most comprehensive fix
        suitability := fsg.calculateBoundaryRefactorSuitability(cycle)

        // Effort is typically high for boundary refactoring
        effort := fsg.calculateBoundaryRefactorEffort(cycle)

        // Generate contextual pros/cons
        pros, cons := fsg.generateBoundaryRefactorProsCons(cycle)

        return &types.FixStrategy{
            Type:           types.FixStrategyBoundaryRefactor,
            Name:           "Module Boundary Refactoring",
            Description:    "Restructure package boundaries to eliminate the cyclic relationship.",
            Suitability:    suitability,
            Effort:         effort,
            Pros:           pros,
            Cons:           cons,
            TargetPackages: fsg.getBoundaryRefactorTargetPackages(cycle),
        }
    }
    ```
  - [x] 5.2 Implement suitability calculation
  - [x] 5.3 Implement effort estimation
  - [x] 5.4 Implement contextual pros/cons generation
  - [x] 5.5 Add tests for Boundary Refactoring strategy

- [x] **Task 6: Implement Suitability Scoring Algorithm** (AC: #2)
  - [x] 6.1 Create scoring helper functions:
    ```go
    // Suitability factors (each 0-10, combined with weights)

    // Factor 1: Cycle depth impact
    func cycleDepthFactor(depth int, strategyType FixStrategyType) int {
        switch strategyType {
        case FixStrategyExtractModule:
            // Better for longer cycles
            if depth >= 4 { return 10 }
            if depth == 3 { return 7 }
            return 4 // Direct cycle
        case FixStrategyDependencyInject:
            // Better for shorter cycles
            if depth == 2 { return 10 }
            if depth == 3 { return 6 }
            return 3 // Long cycles harder to DI
        case FixStrategyBoundaryRefactor:
            // Always viable, slightly better for medium
            if depth == 3 { return 8 }
            return 6
        }
        return 5
    }

    // Factor 2: Dependency type impact
    func dependencyTypeFactor(chain []types.DependencyEdge, strategyType FixStrategyType) int {
        hasDevDep := false
        for _, edge := range chain {
            if edge.Type == types.DependencyTypeDev {
                hasDevDep = true
                break
            }
        }

        switch strategyType {
        case FixStrategyExtractModule:
            // Dev deps are easier to extract
            if hasDevDep { return 8 }
            return 6
        case FixStrategyDependencyInject:
            // Works well for production deps
            if !hasDevDep { return 8 }
            return 5
        case FixStrategyBoundaryRefactor:
            return 6 // Neutral
        }
        return 5
    }

    // Factor 3: Package naming patterns
    func namingPatternFactor(cycle []string, strategyType FixStrategyType) int {
        hasCorePackage := containsPattern(cycle, []string{"core", "common", "shared", "utils"})

        switch strategyType {
        case FixStrategyExtractModule:
            // If already has core, extraction less valuable
            if hasCorePackage { return 4 }
            return 8
        case FixStrategyBoundaryRefactor:
            // Core package cycles often need boundary rework
            if hasCorePackage { return 9 }
            return 5
        default:
            return 6
        }
    }
    ```
  - [x] 6.2 Combine factors with weights
  - [x] 6.3 Add tests for scoring algorithm

- [x] **Task 7: Implement Effort Estimation** (AC: #3)
  - [x] 7.1 Create effort estimation logic:
    ```go
    func estimateEffort(cycle *types.CircularDependencyInfo, strategyType FixStrategyType) types.EffortLevel {
        depth := cycle.Depth

        switch strategyType {
        case FixStrategyExtractModule:
            // Effort scales with cycle depth
            if depth <= 2 { return types.EffortMedium }
            if depth <= 4 { return types.EffortMedium }
            return types.EffortHigh

        case FixStrategyDependencyInject:
            // DI is typically medium effort
            if depth <= 2 { return types.EffortLow }
            if depth <= 3 { return types.EffortMedium }
            return types.EffortHigh

        case FixStrategyBoundaryRefactor:
            // Boundary refactoring is usually high effort
            if depth <= 2 { return types.EffortMedium }
            return types.EffortHigh
        }

        return types.EffortMedium
    }
    ```
  - [x] 7.2 Add tests for effort estimation

- [x] **Task 8: Implement Contextual Pros/Cons Generation** (AC: #4)
  - [x] 8.1 Create pros/cons generators:
    ```go
    func (fsg *FixStrategyGenerator) generateExtractModuleProsCons(cycle *types.CircularDependencyInfo) ([]string, []string) {
        pros := []string{
            "Creates clear separation of concerns",
            fmt.Sprintf("Isolates shared code between %s", formatPackageList(cycle.Cycle[:len(cycle.Cycle)-1])),
        }

        cons := []string{
            "Introduces a new package to maintain",
        }

        // Add contextual pros/cons
        if cycle.Depth >= 4 {
            pros = append(pros, "Effectively breaks complex multi-package cycle")
        }
        if containsPattern(cycle.Cycle, []string{"core", "common"}) {
            cons = append(cons, "May require updating existing core package consumers")
        }

        return pros, cons
    }

    func (fsg *FixStrategyGenerator) generateDIProsCons(cycle *types.CircularDependencyInfo) ([]string, []string) {
        pros := []string{
            "Minimal code changes required",
            "Preserves existing package structure",
        }

        cons := []string{
            "Adds indirection to the codebase",
        }

        // Add contextual pros/cons based on critical edge
        if cycle.RootCause != nil && cycle.RootCause.CriticalEdge != nil {
            edge := cycle.RootCause.CriticalEdge
            pros = append(pros, fmt.Sprintf("Clear injection point: %s → %s", edge.From, edge.To))
        }

        if cycle.Depth > 3 {
            cons = append(cons, "May require multiple interfaces for complex cycles")
        }

        return pros, cons
    }

    func (fsg *FixStrategyGenerator) generateBoundaryRefactorProsCons(cycle *types.CircularDependencyInfo) ([]string, []string) {
        pros := []string{
            "Addresses root architectural issue",
            "Results in cleaner long-term design",
        }

        cons := []string{
            "Requires significant refactoring effort",
            "May affect external package consumers",
        }

        // Add contextual pros/cons
        if containsPattern(cycle.Cycle, []string{"core", "common", "shared"}) {
            pros = append(pros, "Opportunity to properly define core package boundaries")
        }

        return pros, cons
    }
    ```
  - [x] 8.2 Add tests for pros/cons generation

- [x] **Task 9: Implement Strategy Ranking** (AC: #5)
  - [x] 9.1 Implement `rankStrategies`:
    ```go
    func rankStrategies(strategies []types.FixStrategy) []types.FixStrategy {
        // Sort by suitability (descending), then effort (ascending)
        sort.Slice(strategies, func(i, j int) bool {
            if strategies[i].Suitability != strategies[j].Suitability {
                return strategies[i].Suitability > strategies[j].Suitability
            }
            // Lower effort wins ties
            effortOrder := map[types.EffortLevel]int{
                types.EffortLow:    0,
                types.EffortMedium: 1,
                types.EffortHigh:   2,
            }
            return effortOrder[strategies[i].Effort] < effortOrder[strategies[j].Effort]
        })

        // Mark top strategy as recommended
        if len(strategies) > 0 {
            strategies[0].Recommended = true
        }

        return strategies
    }
    ```
  - [x] 9.2 Add tests for ranking logic

- [x] **Task 10: Integrate with CircularDependencyInfo** (AC: #6)
  - [x] 10.1 Update `pkg/types/circular.go`:
    ```go
    type CircularDependencyInfo struct {
        Cycle        []string            `json:"cycle"`
        Type         CircularType        `json:"type"`
        Severity     CircularSeverity    `json:"severity"`
        Depth        int                 `json:"depth"`
        Impact       string              `json:"impact"`
        Complexity   int                 `json:"complexity"`
        RootCause    *RootCauseAnalysis  `json:"rootCause,omitempty"`    // Story 3.1
        ImportTraces []ImportTrace       `json:"importTraces,omitempty"` // Story 3.2
        FixStrategies []FixStrategy      `json:"fixStrategies,omitempty"` // Story 3.3 NEW
    }
    ```
  - [x] 10.2 Verify existing tests still pass

- [x] **Task 11: Wire to Analyzer Pipeline** (AC: all)
  - [x] 11.1 Update `pkg/analyzer/analyzer.go`:
    ```go
    func (a *Analyzer) AnalyzeWithSources(...) (*types.AnalysisResult, error) {
        // ... existing analysis ...

        // Enrich cycles with root cause (Story 3.1)
        rootCauseAnalyzer := NewRootCauseAnalyzer(graph)
        for _, cycle := range cycles {
            cycle.RootCause = rootCauseAnalyzer.Analyze(cycle)
        }

        // Enrich cycles with import traces (Story 3.2)
        if len(sourceFiles) > 0 {
            importTracer := NewImportTracer(workspace, sourceFiles)
            for _, cycle := range cycles {
                cycle.ImportTraces = importTracer.Trace(cycle)
            }
        }

        // NEW: Generate fix strategies (Story 3.3)
        fixGenerator := NewFixStrategyGenerator(graph, workspace)
        for _, cycle := range cycles {
            cycle.FixStrategies = fixGenerator.Generate(cycle)
        }

        return result, nil
    }
    ```
  - [x] 11.2 Update analyzer tests

- [x] **Task 12: Update TypeScript Types** (AC: #6)
  - [x] 12.1 Update `packages/types/src/analysis/results.ts`:
    ```typescript
    export type FixStrategyType = 'extract-module' | 'dependency-injection' | 'boundary-refactoring';
    export type EffortLevel = 'low' | 'medium' | 'high';

    export interface FixStrategy {
      type: FixStrategyType;
      name: string;
      description: string;
      suitability: number;
      effort: EffortLevel;
      pros: string[];
      cons: string[];
      recommended: boolean;
      targetPackages: string[];
      newPackageName?: string;
    }

    export interface CircularDependencyInfo {
      cycle: string[];
      type: 'direct' | 'indirect';
      severity: 'critical' | 'warning' | 'info';
      depth: number;
      impact: string;
      complexity: number;
      rootCause?: RootCauseAnalysis;    // Story 3.1
      importTraces?: ImportTrace[];      // Story 3.2
      fixStrategies?: FixStrategy[];     // Story 3.3 NEW
    }
    ```
  - [x] 12.2 Run `pnpm nx build types` to verify
  - [x] 12.3 Update type tests if needed

- [x] **Task 13: Performance Testing** (AC: #7)
  - [x] 13.1 Create `pkg/analyzer/fix_strategy_generator_benchmark_test.go`:
    ```go
    func BenchmarkFixStrategyGeneration(b *testing.B) {
        graph := generateGraphWithCycles(100, 5)
        workspace := generateWorkspace(100)
        detector := NewCycleDetector(graph)
        cycles := detector.DetectCycles()

        // Enrich with root cause first
        rca := NewRootCauseAnalyzer(graph)
        for _, cycle := range cycles {
            cycle.RootCause = rca.Analyze(cycle)
        }

        generator := NewFixStrategyGenerator(graph, workspace)

        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            for _, cycle := range cycles {
                generator.Generate(cycle)
            }
        }
    }
    ```
  - [x] 13.2 Verify < 200ms for 100 packages with 5 cycles
  - [x] 13.3 Document actual performance in completion notes

- [x] **Task 14: Integration Verification** (AC: all)
  - [x] 14.1 Run all tests: `cd packages/analysis-engine && make test`
  - [x] 14.2 Build WASM: `pnpm nx build @monoguard/analysis-engine`
  - [x] 14.3 Run affected CI checks: `pnpm nx affected --target=lint,test,type-check --base=main`
  - [x] 14.4 Verify JSON output includes fixStrategies field

## Dev Notes

### Architecture Patterns & Constraints

**From Architecture Document:**
- **Location:** Fix strategy generator in `pkg/analyzer/`
- **Pattern:** Generator pattern with graph + workspace input
- **Integration:** Enriches CircularDependencyInfo with optional FixStrategies
- **Dependency:** Requires Story 3.1 (RootCause) to be implemented first

**Critical Constraints:**
- **camelCase JSON:** All struct tags MUST use camelCase
- **Optional Fields:** FixStrategies is `omitempty` - backward compatible
- **Pattern Matching:** Use rule-based logic (not ML) for Phase 0
- **Contextual Output:** Pros/cons must be specific to the cycle, not generic

### Critical Don't-Miss Rules

**From project-context.md:**

1. **JSON Naming Convention:**
   ```go
   // ✅ CORRECT: camelCase JSON tags
   type FixStrategy struct {
       TargetPackages []string `json:"targetPackages"`
       NewPackageName string   `json:"newPackageName,omitempty"`
   }

   // ❌ WRONG: snake_case or PascalCase JSON tags
   type FixStrategy struct {
       TargetPackages []string `json:"target_packages"` // WRONG!
   }
   ```

2. **Enum Constants:**
   ```go
   // ✅ CORRECT: kebab-case string values for TypeScript compatibility
   const (
       FixStrategyExtractModule    FixStrategyType = "extract-module"
       FixStrategyDependencyInject FixStrategyType = "dependency-injection"
       FixStrategyBoundaryRefactor FixStrategyType = "boundary-refactoring"
   )

   // ❌ WRONG: camelCase or UPPER_CASE
   const (
       FixStrategyExtractModule FixStrategyType = "extractModule" // WRONG!
   )
   ```

3. **Slice vs Nil:**
   ```go
   // ✅ CORRECT: Return empty slice for JSON serialization
   func (fsg *FixStrategyGenerator) Generate(...) []types.FixStrategy {
       strategies := []types.FixStrategy{} // Never nil
       // ...
       return strategies
   }
   ```

### Project Structure Notes

**Target Directory Structure:**
```
packages/analysis-engine/
├── pkg/
│   ├── analyzer/
│   │   ├── analyzer.go                        # UPDATE: Add fix strategy generation
│   │   ├── root_cause_analyzer.go             # From Story 3.1
│   │   ├── import_tracer.go                   # From Story 3.2
│   │   ├── fix_strategy_generator.go          # NEW: Fix strategy generator
│   │   ├── fix_strategy_generator_test.go     # NEW: Generator tests
│   │   └── fix_strategy_generator_benchmark_test.go # NEW: Performance
│   └── types/
│       ├── circular.go                        # UPDATE: Add FixStrategies field
│       ├── fix_strategy.go                    # NEW: FixStrategy types
│       └── fix_strategy_test.go               # NEW: Type tests
└── ...

packages/types/src/analysis/
├── results.ts                                 # UPDATE: Add TS types
└── ...
```

### Previous Story Intelligence

**From Story 3.1 (ready-for-dev):**
- RootCauseAnalysis provides: originatingPackage, criticalEdge, confidence, chain
- DependencyEdge has: from, to, type (production/dev/peer/optional), critical
- **Key Insight:** Use criticalEdge to determine best injection point for DI strategy
- **Key Insight:** Use originatingPackage to determine primary target for boundary refactoring

**From Story 3.2 (ready-for-dev):**
- ImportTraces show actual code locations
- **Key Insight:** Can use import count to estimate effort more accurately (future enhancement)

**From Story 2.3 (done):**
- CircularDependencyInfo has: cycle, type (direct/indirect), depth, severity, complexity
- **Key Insight:** Use depth and type for suitability scoring

### Strategy Selection Heuristics

**Extract Shared Module - Best When:**
- Cycle involves 3+ packages (indirect)
- Multiple packages depend on the same functionality
- No existing "core" or "shared" package
- **Suitability Score:** Higher for longer cycles, lower when core exists

**Dependency Injection - Best When:**
- Direct cycle (A ↔ B)
- Clear critical edge identified
- One direction is clearly the "right" dependency
- **Suitability Score:** Higher for direct cycles, lower for complex cycles

**Module Boundary Refactoring - Best When:**
- Cycle involves "core", "common", or "shared" packages
- Architectural smell detected (responsibilities overlap)
- Other strategies don't fit well
- **Suitability Score:** Higher when core packages involved

### Suitability Scoring Matrix

| Factor | Extract Module | DI | Boundary Refactor |
|--------|---------------|----|--------------------|
| Direct cycle (depth 2) | 4 | 10 | 6 |
| Short indirect (depth 3) | 7 | 6 | 8 |
| Long indirect (depth 4+) | 10 | 3 | 6 |
| Has dev dependency | 8 | 5 | 6 |
| All production deps | 6 | 8 | 6 |
| Has "core" package | 4 | 6 | 9 |
| No "core" package | 8 | 6 | 5 |

**Final Score:** Weighted average (depth: 40%, dep type: 30%, naming: 30%)

### Input/Output Format

**Input (CircularDependencyInfo with RootCause):**
```json
{
  "cycle": ["@mono/ui", "@mono/api", "@mono/core", "@mono/ui"],
  "type": "indirect",
  "depth": 3,
  "rootCause": {
    "originatingPackage": "@mono/ui",
    "criticalEdge": { "from": "@mono/core", "to": "@mono/ui", "type": "production" },
    "confidence": 82
  }
}
```

**Output (CircularDependencyInfo with FixStrategies):**
```json
{
  "cycle": ["@mono/ui", "@mono/api", "@mono/core", "@mono/ui"],
  "type": "indirect",
  "depth": 3,
  "rootCause": { ... },
  "fixStrategies": [
    {
      "type": "extract-module",
      "name": "Extract Shared Module",
      "description": "Create a new shared package to hold common dependencies, breaking the cycle.",
      "suitability": 8,
      "effort": "medium",
      "pros": [
        "Creates clear separation of concerns",
        "Isolates shared code between @mono/ui, @mono/api, @mono/core"
      ],
      "cons": [
        "Introduces a new package to maintain"
      ],
      "recommended": true,
      "targetPackages": ["@mono/ui", "@mono/api", "@mono/core"],
      "newPackageName": "@mono/shared"
    },
    {
      "type": "dependency-injection",
      "name": "Dependency Injection",
      "description": "Invert the problematic dependency by introducing an interface or callback pattern.",
      "suitability": 6,
      "effort": "medium",
      "pros": [
        "Minimal code changes required",
        "Clear injection point: @mono/core → @mono/ui"
      ],
      "cons": [
        "Adds indirection to the codebase"
      ],
      "recommended": false,
      "targetPackages": ["@mono/core", "@mono/ui"]
    },
    {
      "type": "boundary-refactoring",
      "name": "Module Boundary Refactoring",
      "description": "Restructure package boundaries to eliminate the cyclic relationship.",
      "suitability": 7,
      "effort": "high",
      "pros": [
        "Addresses root architectural issue",
        "Opportunity to properly define core package boundaries"
      ],
      "cons": [
        "Requires significant refactoring effort",
        "May affect external package consumers"
      ],
      "recommended": false,
      "targetPackages": ["@mono/ui", "@mono/api", "@mono/core"]
    }
  ]
}
```

### Test Scenarios

| Scenario | Cycle | Expected Top Strategy | Suitability |
|----------|-------|----------------------|-------------|
| Direct cycle | A ↔ B | Dependency Injection | 8-10 |
| 3-package cycle | A → B → C → A | Extract Module | 7-8 |
| 4+ package cycle | A → B → C → D → A | Extract Module | 9-10 |
| Has "core" package | ui → api → core → ui | Boundary Refactoring | 8-9 |
| Dev dependency cycle | A ↔ B (dev) | Extract Module | 7-8 |

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story 3.3]
- [Source: _bmad-output/planning-artifacts/prd.md#FR9-FR11]
- [Source: _bmad-output/project-context.md#Language-Specific Rules]
- [Source: _bmad-output/implementation-artifacts/3-1-implement-root-cause-analysis-for-circular-dependencies.md]
- [Refactoring Patterns - Martin Fowler](https://refactoring.com/)
- [Dependency Injection Explained](https://martinfowler.com/articles/injection.html)

## Dev Agent Record

### Agent Model Used

Claude Opus 4.5 (claude-opus-4-5-20251101)

### Debug Log References

None required - all tests passed on first run.

### Completion Notes List

1. **All 14 tasks completed successfully** - Story 3.3 fix strategy recommendations fully implemented
2. **Performance exceeds requirements**:
   - Direct cycle (2 nodes): ~4.3 microseconds (0.004ms)
   - Indirect cycle (3 nodes): ~5.5 microseconds (0.005ms)
   - Long cycle (10 nodes): ~12 microseconds (0.012ms)
   - All well under the 5ms per cycle requirement (AC7)
3. **Strategy heuristics implemented**:
   - Extract Module favored for longer cycles (4+ packages)
   - Dependency Injection favored for direct cycles (2 packages)
   - Boundary Refactoring favored when core/shared packages involved
4. **Full CI verification passed**:
   - All Go tests pass
   - Nx affected lint passes
   - Nx affected type-check passes
5. **TypeScript types updated** with FixStrategyType and EffortLevel union types

### File List

**New Files Created:**
- `packages/analysis-engine/pkg/types/fix_strategy.go` - FixStrategy types
- `packages/analysis-engine/pkg/types/fix_strategy_test.go` - Type tests
- `packages/analysis-engine/pkg/analyzer/fix_strategy_generator.go` - Generator implementation
- `packages/analysis-engine/pkg/analyzer/fix_strategy_generator_test.go` - Generator tests
- `packages/analysis-engine/pkg/analyzer/fix_strategy_generator_benchmark_test.go` - Performance benchmarks

**Modified Files:**
- `packages/analysis-engine/pkg/types/circular.go` - Added FixStrategies field
- `packages/analysis-engine/pkg/types/circular_test.go` - Added FixStrategies tests
- `packages/analysis-engine/pkg/analyzer/analyzer.go` - Wired fix strategy generation
- `packages/analysis-engine/pkg/analyzer/analyzer_test.go` - Added integration tests
- `packages/types/src/analysis/results.ts` - Updated TypeScript types

### Change Log

| Date | Change | Files |
|------|--------|-------|
| 2026-01-23 | Created FixStrategy and EffortLevel types | pkg/types/fix_strategy.go |
| 2026-01-23 | Implemented FixStrategyGenerator with 3 strategies | pkg/analyzer/fix_strategy_generator.go |
| 2026-01-23 | Added suitability scoring algorithm | pkg/analyzer/fix_strategy_generator.go |
| 2026-01-23 | Integrated FixStrategies into CircularDependencyInfo | pkg/types/circular.go |
| 2026-01-23 | Wired generator to Analyze and AnalyzeWithSources | pkg/analyzer/analyzer.go |
| 2026-01-23 | Updated TypeScript types | packages/types/src/analysis/results.ts |
| 2026-01-23 | Added performance benchmarks | pkg/analyzer/fix_strategy_generator_benchmark_test.go |
