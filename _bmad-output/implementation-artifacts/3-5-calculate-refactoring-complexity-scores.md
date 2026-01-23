# Story 3.5: Calculate Refactoring Complexity Scores

Status: ready-for-dev

## Story

As a **user**,
I want **to see complexity scores for each circular dependency fix**,
So that **I can prioritize which cycles to fix based on effort required**.

## Acceptance Criteria

1. **AC1: Enhanced Complexity Score (1-10)**
   - Given a circular dependency
   - When I calculate refactoring complexity
   - Then the score (1-10) considers:
     - Number of files affected
     - Number of import statements to change
     - Depth of the dependency chain
     - Number of packages involved
     - Presence of external dependencies in the cycle
   - And score is 1-10 where 1 is simple, 10 is complex

2. **AC2: Complexity Score Breakdown**
   - Given a complexity score
   - When viewing the details
   - Then I see breakdown of contributing factors:
     - `filesAffected`: number of source files that need changes
     - `importsToChange`: number of import statements to modify
     - `chainDepth`: depth of the dependency chain
     - `packagesInvolved`: number of packages in the cycle
     - `hasExternalDeps`: whether external dependencies are involved
   - And each factor has a sub-score contribution

3. **AC3: Estimated Time Range**
   - Given a complexity score
   - When generating the estimate
   - Then provide human-readable time range:
     - Score 1-2: "5-15 minutes"
     - Score 3-4: "15-30 minutes"
     - Score 5-6: "30-60 minutes"
     - Score 7-8: "1-2 hours"
     - Score 9-10: "2-4 hours"
   - And time range is included in `RefactoringComplexity` struct

4. **AC4: Integration with CircularDependencyInfo**
   - Given analysis results with circular dependencies
   - When enriching CircularDependencyInfo
   - Then add `refactoringComplexity` field:
     ```go
     type CircularDependencyInfo struct {
         // ... existing fields ...
         RefactoringComplexity *RefactoringComplexity `json:"refactoringComplexity,omitempty"`
     }
     ```
   - And existing `complexity` field (int) remains for backward compatibility
   - And `refactoringComplexity` provides detailed breakdown

5. **AC5: Import Trace Integration**
   - Given ImportTraces from Story 3.2
   - When calculating complexity
   - Then use actual import data:
     - Count unique files from ImportTraces
     - Count total import statements
     - Analyze import types (named, default, namespace)
   - And if ImportTraces unavailable, estimate based on cycle depth

6. **AC6: Fix Strategy Integration**
   - Given FixStrategies from Story 3.3
   - When calculating complexity
   - Then complexity is per-strategy:
     - Each FixStrategy includes its own complexity
     - Different strategies may have different complexities
   - And complexity influences strategy ranking

7. **AC7: Performance**
   - Given a workspace with 100 packages and 5 cycles
   - When calculating complexity for all cycles
   - Then calculation completes in < 100ms additional overhead
   - And memory usage increase is < 5MB

## Tasks / Subtasks

- [ ] **Task 1: Define RefactoringComplexity Types** (AC: #1, #2, #3)
  - [ ] 1.1 Create `pkg/types/refactoring_complexity.go`:
    ```go
    package types

    // RefactoringComplexity provides detailed breakdown of fix complexity.
    // Matches @monoguard/types RefactoringComplexity interface.
    type RefactoringComplexity struct {
        // Score is the overall complexity (1-10)
        Score int `json:"score"`

        // EstimatedTime is human-readable time range (e.g., "15-30 minutes")
        EstimatedTime string `json:"estimatedTime"`

        // Breakdown shows individual factor contributions
        Breakdown ComplexityBreakdown `json:"breakdown"`

        // Explanation provides human-readable summary
        Explanation string `json:"explanation"`
    }

    // ComplexityBreakdown shows how each factor contributes to the score.
    type ComplexityBreakdown struct {
        // FilesAffected is number of source files that need changes
        FilesAffected ComplexityFactor `json:"filesAffected"`

        // ImportsToChange is number of import statements to modify
        ImportsToChange ComplexityFactor `json:"importsToChange"`

        // ChainDepth is the dependency chain depth
        ChainDepth ComplexityFactor `json:"chainDepth"`

        // PackagesInvolved is number of packages in the cycle
        PackagesInvolved ComplexityFactor `json:"packagesInvolved"`

        // ExternalDependencies indicates if external deps are involved
        ExternalDependencies ComplexityFactor `json:"externalDependencies"`
    }

    // ComplexityFactor represents a single factor in complexity calculation.
    type ComplexityFactor struct {
        // Value is the raw value for this factor
        Value int `json:"value"`

        // Weight is the factor weight (0.0-1.0)
        Weight float64 `json:"weight"`

        // Contribution is the weighted score contribution
        Contribution float64 `json:"contribution"`

        // Description explains what this factor measures
        Description string `json:"description"`
    }
    ```
  - [ ] 1.2 Add JSON serialization tests in `pkg/types/refactoring_complexity_test.go`
  - [ ] 1.3 Ensure all JSON tags use camelCase

- [ ] **Task 2: Create ComplexityCalculator** (AC: #1, #5)
  - [ ] 2.1 Create `pkg/analyzer/complexity_calculator.go`:
    ```go
    package analyzer

    import "github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"

    // ComplexityCalculator computes refactoring complexity scores.
    type ComplexityCalculator struct {
        graph     *types.DependencyGraph
        workspace *types.WorkspaceData
    }

    // NewComplexityCalculator creates a new calculator.
    func NewComplexityCalculator(graph *types.DependencyGraph, workspace *types.WorkspaceData) *ComplexityCalculator

    // Calculate computes the refactoring complexity for a cycle.
    func (cc *ComplexityCalculator) Calculate(cycle *types.CircularDependencyInfo) *types.RefactoringComplexity

    // calculateFilesAffected estimates files needing changes.
    func (cc *ComplexityCalculator) calculateFilesAffected(cycle *types.CircularDependencyInfo) types.ComplexityFactor

    // calculateImportsToChange counts import statements to modify.
    func (cc *ComplexityCalculator) calculateImportsToChange(cycle *types.CircularDependencyInfo) types.ComplexityFactor

    // calculateChainDepth evaluates dependency chain depth.
    func (cc *ComplexityCalculator) calculateChainDepth(cycle *types.CircularDependencyInfo) types.ComplexityFactor

    // calculatePackagesInvolved counts packages in cycle.
    func (cc *ComplexityCalculator) calculatePackagesInvolved(cycle *types.CircularDependencyInfo) types.ComplexityFactor

    // calculateExternalDependencies checks for external deps in cycle.
    func (cc *ComplexityCalculator) calculateExternalDependencies(cycle *types.CircularDependencyInfo) types.ComplexityFactor

    // estimateTime converts score to human-readable time range.
    func estimateTime(score int) string

    // generateExplanation creates human-readable complexity summary.
    func generateExplanation(score int, breakdown *types.ComplexityBreakdown) string
    ```
  - [ ] 2.2 Implement complexity calculation logic
  - [ ] 2.3 Create comprehensive tests in `pkg/analyzer/complexity_calculator_test.go`

- [ ] **Task 3: Implement Files Affected Calculation** (AC: #1, #5)
  - [ ] 3.1 Implement `calculateFilesAffected`:
    ```go
    func (cc *ComplexityCalculator) calculateFilesAffected(cycle *types.CircularDependencyInfo) types.ComplexityFactor {
        var fileCount int

        // Use ImportTraces if available (Story 3.2 provides these)
        if len(cycle.ImportTraces) > 0 {
            uniqueFiles := make(map[string]bool)
            for _, trace := range cycle.ImportTraces {
                uniqueFiles[trace.FilePath] = true
            }
            fileCount = len(uniqueFiles)
        } else {
            // Estimate: ~2 files per package in cycle (entry point + consumers)
            fileCount = (cycle.Depth) * 2
        }

        // Calculate contribution (weight: 0.25)
        // Scale: 1-2 files = low, 3-5 = medium, 6+ = high
        var contribution float64
        switch {
        case fileCount <= 2:
            contribution = 0.5
        case fileCount <= 5:
            contribution = 1.5
        case fileCount <= 10:
            contribution = 2.0
        default:
            contribution = 2.5
        }

        return types.ComplexityFactor{
            Value:        fileCount,
            Weight:       0.25,
            Contribution: contribution,
            Description:  fmt.Sprintf("%d source files need modification", fileCount),
        }
    }
    ```
  - [ ] 3.2 Add tests for files affected calculation

- [ ] **Task 4: Implement Imports To Change Calculation** (AC: #1, #5)
  - [ ] 4.1 Implement `calculateImportsToChange`:
    ```go
    func (cc *ComplexityCalculator) calculateImportsToChange(cycle *types.CircularDependencyInfo) types.ComplexityFactor {
        var importCount int

        // Use ImportTraces if available
        if len(cycle.ImportTraces) > 0 {
            importCount = len(cycle.ImportTraces)
        } else {
            // Estimate: ~1-2 imports per edge in cycle
            importCount = cycle.Depth
        }

        // Calculate contribution (weight: 0.20)
        // Scale: 1-2 imports = low, 3-6 = medium, 7+ = high
        var contribution float64
        switch {
        case importCount <= 2:
            contribution = 0.4
        case importCount <= 6:
            contribution = 1.2
        case importCount <= 10:
            contribution = 1.6
        default:
            contribution = 2.0
        }

        return types.ComplexityFactor{
            Value:        importCount,
            Weight:       0.20,
            Contribution: contribution,
            Description:  fmt.Sprintf("%d import statements need updating", importCount),
        }
    }
    ```
  - [ ] 4.2 Add tests for imports calculation

- [ ] **Task 5: Implement Chain Depth Calculation** (AC: #1)
  - [ ] 5.1 Implement `calculateChainDepth`:
    ```go
    func (cc *ComplexityCalculator) calculateChainDepth(cycle *types.CircularDependencyInfo) types.ComplexityFactor {
        depth := cycle.Depth

        // Calculate contribution (weight: 0.25)
        // Scale: depth 2 = low, 3-4 = medium, 5+ = high
        var contribution float64
        switch {
        case depth <= 2:
            contribution = 0.5
        case depth <= 4:
            contribution = 1.5
        case depth <= 6:
            contribution = 2.0
        default:
            contribution = 2.5
        }

        return types.ComplexityFactor{
            Value:        depth,
            Weight:       0.25,
            Contribution: contribution,
            Description:  fmt.Sprintf("Dependency chain has %d levels", depth),
        }
    }
    ```
  - [ ] 5.2 Add tests for chain depth calculation

- [ ] **Task 6: Implement Packages Involved Calculation** (AC: #1)
  - [ ] 6.1 Implement `calculatePackagesInvolved`:
    ```go
    func (cc *ComplexityCalculator) calculatePackagesInvolved(cycle *types.CircularDependencyInfo) types.ComplexityFactor {
        // Unique packages (cycle array ends with first, so -1)
        packageCount := len(cycle.Cycle) - 1
        if packageCount < 1 {
            packageCount = 1
        }

        // Calculate contribution (weight: 0.15)
        // Scale: 2 packages = low, 3-4 = medium, 5+ = high
        var contribution float64
        switch {
        case packageCount <= 2:
            contribution = 0.3
        case packageCount <= 4:
            contribution = 0.9
        case packageCount <= 6:
            contribution = 1.2
        default:
            contribution = 1.5
        }

        return types.ComplexityFactor{
            Value:        packageCount,
            Weight:       0.15,
            Contribution: contribution,
            Description:  fmt.Sprintf("%d packages involved in cycle", packageCount),
        }
    }
    ```
  - [ ] 6.2 Add tests for packages calculation

- [ ] **Task 7: Implement External Dependencies Check** (AC: #1)
  - [ ] 7.1 Implement `calculateExternalDependencies`:
    ```go
    func (cc *ComplexityCalculator) calculateExternalDependencies(cycle *types.CircularDependencyInfo) types.ComplexityFactor {
        hasExternal := false

        // Check if any edge in root cause chain involves external deps
        if cycle.RootCause != nil {
            for _, edge := range cycle.RootCause.Chain {
                // Check if target is external (not in workspace packages)
                if cc.isExternalPackage(edge.To) {
                    hasExternal = true
                    break
                }
            }
        }

        // Calculate contribution (weight: 0.15)
        var contribution float64
        var value int
        if hasExternal {
            contribution = 1.5
            value = 1
        } else {
            contribution = 0.15
            value = 0
        }

        description := "No external dependencies in cycle"
        if hasExternal {
            description = "External dependencies increase complexity"
        }

        return types.ComplexityFactor{
            Value:        value,
            Weight:       0.15,
            Contribution: contribution,
            Description:  description,
        }
    }

    // isExternalPackage checks if a package is external (not in workspace).
    func (cc *ComplexityCalculator) isExternalPackage(pkgName string) bool {
        if cc.workspace == nil {
            return false
        }
        _, exists := cc.workspace.Packages[pkgName]
        return !exists
    }
    ```
  - [ ] 7.2 Add tests for external dependencies check

- [ ] **Task 8: Implement Score Aggregation and Time Estimation** (AC: #1, #3)
  - [ ] 8.1 Implement main `Calculate` function:
    ```go
    func (cc *ComplexityCalculator) Calculate(cycle *types.CircularDependencyInfo) *types.RefactoringComplexity {
        breakdown := types.ComplexityBreakdown{
            FilesAffected:        cc.calculateFilesAffected(cycle),
            ImportsToChange:      cc.calculateImportsToChange(cycle),
            ChainDepth:           cc.calculateChainDepth(cycle),
            PackagesInvolved:     cc.calculatePackagesInvolved(cycle),
            ExternalDependencies: cc.calculateExternalDependencies(cycle),
        }

        // Sum all contributions
        totalContribution := breakdown.FilesAffected.Contribution +
            breakdown.ImportsToChange.Contribution +
            breakdown.ChainDepth.Contribution +
            breakdown.PackagesInvolved.Contribution +
            breakdown.ExternalDependencies.Contribution

        // Convert to 1-10 scale (max contribution is ~10)
        score := int(math.Round(totalContribution))
        if score < 1 {
            score = 1
        }
        if score > 10 {
            score = 10
        }

        return &types.RefactoringComplexity{
            Score:         score,
            EstimatedTime: estimateTime(score),
            Breakdown:     breakdown,
            Explanation:   generateExplanation(score, &breakdown),
        }
    }
    ```
  - [ ] 8.2 Implement `estimateTime`:
    ```go
    func estimateTime(score int) string {
        switch {
        case score <= 2:
            return "5-15 minutes"
        case score <= 4:
            return "15-30 minutes"
        case score <= 6:
            return "30-60 minutes"
        case score <= 8:
            return "1-2 hours"
        default:
            return "2-4 hours"
        }
    }
    ```
  - [ ] 8.3 Implement `generateExplanation`:
    ```go
    func generateExplanation(score int, breakdown *types.ComplexityBreakdown) string {
        var level string
        switch {
        case score <= 3:
            level = "straightforward"
        case score <= 6:
            level = "moderate"
        case score <= 8:
            level = "significant"
        default:
            level = "complex"
        }

        return fmt.Sprintf(
            "%s refactoring: %d files, %d imports, %d-level chain",
            strings.Title(level),
            breakdown.FilesAffected.Value,
            breakdown.ImportsToChange.Value,
            breakdown.ChainDepth.Value,
        )
    }
    ```
  - [ ] 8.4 Add tests for score aggregation

- [ ] **Task 9: Integrate with CircularDependencyInfo** (AC: #4)
  - [ ] 9.1 Update `pkg/types/circular.go`:
    ```go
    type CircularDependencyInfo struct {
        Cycle                 []string               `json:"cycle"`
        Type                  CircularType           `json:"type"`
        Severity              CircularSeverity       `json:"severity"`
        Depth                 int                    `json:"depth"`
        Impact                string                 `json:"impact"`
        Complexity            int                    `json:"complexity"`              // Legacy field
        RefactoringComplexity *RefactoringComplexity `json:"refactoringComplexity,omitempty"` // NEW Story 3.5
        RootCause             *RootCauseAnalysis     `json:"rootCause,omitempty"`
        ImportTraces          []ImportTrace          `json:"importTraces,omitempty"`
        FixStrategies         []FixStrategy          `json:"fixStrategies,omitempty"`
    }
    ```
  - [ ] 9.2 Verify existing tests still pass

- [ ] **Task 10: Integrate with FixStrategy** (AC: #6)
  - [ ] 10.1 Update `pkg/types/fix_strategy.go`:
    ```go
    type FixStrategy struct {
        // ... existing fields ...
        Complexity *RefactoringComplexity `json:"complexity,omitempty"` // NEW Story 3.5
    }
    ```
  - [ ] 10.2 Update FixStrategyGenerator to include complexity:
    ```go
    // In fix_strategy_generator.go
    func (fsg *FixStrategyGenerator) Generate(cycle *types.CircularDependencyInfo) []types.FixStrategy {
        strategies := []types.FixStrategy{}

        complexityCalc := NewComplexityCalculator(fsg.graph, fsg.workspace)

        // Generate each strategy with its specific complexity
        extractModule := fsg.generateExtractModule(cycle)
        extractModule.Complexity = fsg.calculateStrategyComplexity(complexityCalc, cycle, extractModule)
        strategies = append(strategies, *extractModule)

        // ... similar for other strategies
    }
    ```
  - [ ] 10.3 Add tests for strategy complexity integration

- [ ] **Task 11: Wire to Analyzer Pipeline** (AC: all)
  - [ ] 11.1 Update `pkg/analyzer/analyzer.go`:
    ```go
    func (a *Analyzer) AnalyzeWithSources(...) (*types.AnalysisResult, error) {
        // ... existing analysis ...

        // NEW: Calculate refactoring complexity (Story 3.5)
        complexityCalc := NewComplexityCalculator(graph, workspace)
        for _, cycle := range cycles {
            cycle.RefactoringComplexity = complexityCalc.Calculate(cycle)
        }

        return result, nil
    }
    ```
  - [ ] 11.2 Update analyzer tests

- [ ] **Task 12: Update TypeScript Types** (AC: #4)
  - [ ] 12.1 Update `packages/types/src/analysis/results.ts`:
    ```typescript
    export interface RefactoringComplexity {
      /** Overall complexity score (1-10) */
      score: number;
      /** Human-readable time estimate (e.g., "15-30 minutes") */
      estimatedTime: string;
      /** Breakdown of contributing factors */
      breakdown: ComplexityBreakdown;
      /** Human-readable explanation */
      explanation: string;
    }

    export interface ComplexityBreakdown {
      filesAffected: ComplexityFactor;
      importsToChange: ComplexityFactor;
      chainDepth: ComplexityFactor;
      packagesInvolved: ComplexityFactor;
      externalDependencies: ComplexityFactor;
    }

    export interface ComplexityFactor {
      /** Raw value for this factor */
      value: number;
      /** Factor weight (0.0-1.0) */
      weight: number;
      /** Weighted score contribution */
      contribution: number;
      /** Human-readable description */
      description: string;
    }

    export interface CircularDependencyInfo {
      // ... existing fields ...
      /** Detailed refactoring complexity (Story 3.5) */
      refactoringComplexity?: RefactoringComplexity;
    }

    export interface FixStrategy {
      // ... existing fields ...
      /** Strategy-specific complexity (Story 3.5) */
      complexity?: RefactoringComplexity;
    }
    ```
  - [ ] 12.2 Run `pnpm nx build types` to verify
  - [ ] 12.3 Add type tests for RefactoringComplexity

- [ ] **Task 13: Performance Testing** (AC: #7)
  - [ ] 13.1 Create `pkg/analyzer/complexity_calculator_benchmark_test.go`:
    ```go
    func BenchmarkComplexityCalculation(b *testing.B) {
        workspace := generateWorkspace(100)
        graph := generateGraph(100)
        cycles := generateCyclesWithImportTraces(5)
        calculator := NewComplexityCalculator(graph, workspace)

        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            for _, cycle := range cycles {
                calculator.Calculate(cycle)
            }
        }
    }
    ```
  - [ ] 13.2 Verify < 100ms for 100 packages with 5 cycles
  - [ ] 13.3 Document actual performance in completion notes

- [ ] **Task 14: Integration Verification** (AC: all)
  - [ ] 14.1 Run all tests: `cd packages/analysis-engine && make test`
  - [ ] 14.2 Build WASM: `pnpm nx build @monoguard/analysis-engine`
  - [ ] 14.3 Run affected CI checks: `pnpm nx affected --target=lint,test,type-check --base=main`
  - [ ] 14.4 Verify JSON output includes refactoringComplexity field

## Dev Notes

### Architecture Patterns & Constraints

**From Architecture Document:**
- **Location:** Complexity calculator in `pkg/analyzer/`
- **Pattern:** Calculator pattern with graph + workspace input
- **Integration:** Enriches CircularDependencyInfo with RefactoringComplexity
- **Dependency:** Uses ImportTraces from Story 3.2 (optional enhancement)

**Critical Constraints:**
- **camelCase JSON:** All struct tags MUST use camelCase
- **Optional Fields:** RefactoringComplexity is `omitempty` - backward compatible
- **Legacy Support:** Keep existing `complexity` int field for backward compatibility
- **Performance:** Must not significantly slow down analysis pipeline

### Critical Don't-Miss Rules

**From project-context.md:**

1. **JSON Naming Convention:**
   ```go
   // ✅ CORRECT: camelCase JSON tags
   type RefactoringComplexity struct {
       EstimatedTime string            `json:"estimatedTime"`
       Breakdown     ComplexityBreakdown `json:"breakdown"`
   }

   // ❌ WRONG: snake_case JSON tags
   type RefactoringComplexity struct {
       EstimatedTime string `json:"estimated_time"` // WRONG!
   }
   ```

2. **Float vs Int for Weights:**
   ```go
   // ✅ CORRECT: Use float64 for weights and contributions
   type ComplexityFactor struct {
       Weight       float64 `json:"weight"`
       Contribution float64 `json:"contribution"`
   }

   // ❌ WRONG: Using int for weights
   type ComplexityFactor struct {
       Weight int `json:"weight"` // WRONG - loses precision
   }
   ```

3. **Pointer vs Value for Optional Structs:**
   ```go
   // ✅ CORRECT: Use pointer for optional nested structs
   type CircularDependencyInfo struct {
       RefactoringComplexity *RefactoringComplexity `json:"refactoringComplexity,omitempty"`
   }

   // ❌ WRONG: Value type with omitempty doesn't work for structs
   type CircularDependencyInfo struct {
       RefactoringComplexity RefactoringComplexity `json:"refactoringComplexity,omitempty"` // Won't omit!
   }
   ```

### Project Structure Notes

**Target Directory Structure:**
```
packages/analysis-engine/
├── pkg/
│   ├── analyzer/
│   │   ├── analyzer.go                        # UPDATE: Add complexity calculation
│   │   ├── fix_strategy_generator.go          # UPDATE: Add strategy complexity
│   │   ├── complexity_calculator.go           # NEW: Complexity calculator
│   │   ├── complexity_calculator_test.go      # NEW: Calculator tests
│   │   └── complexity_calculator_benchmark_test.go # NEW: Performance
│   └── types/
│       ├── circular.go                        # UPDATE: Add RefactoringComplexity field
│       ├── fix_strategy.go                    # UPDATE: Add Complexity field
│       ├── refactoring_complexity.go          # NEW: RefactoringComplexity types
│       └── refactoring_complexity_test.go     # NEW: Type tests
└── ...

packages/types/src/analysis/
├── results.ts                                 # UPDATE: Add TS types
└── ...
```

### Previous Story Intelligence

**From Story 3.4 (done):**
- FixGuide has `estimatedTime` field - **Key Insight:** Reuse time estimation logic
- FixStep has clear structure - follow similar patterns
- Performance: ~0.12ms for guide generation - maintain similar performance

**From Story 3.3 (done):**
- FixStrategy has `effort` field (low/medium/high) - different from detailed complexity
- Suitability scoring uses weighted factors - **Key Insight:** Use similar factor-based approach
- Performance: ~5 microseconds per cycle - target similar performance

**From Story 3.2 (done):**
- ImportTrace has filePath, lineNumber, statement - **Key Insight:** Use for accurate file counting
- ImportType enum for classification - can use for complexity weighting

**From Story 3.1 (done):**
- RootCauseAnalysis has chain of edges - **Key Insight:** Use to check for external deps
- CriticalEdge identifies key breaking point

### Complexity Scoring Algorithm

**Factor Weights:**
| Factor | Weight | Rationale |
|--------|--------|-----------|
| Files Affected | 0.25 | Most direct impact on work |
| Imports to Change | 0.20 | Line-by-line changes |
| Chain Depth | 0.25 | Architectural complexity |
| Packages Involved | 0.15 | Coordination overhead |
| External Deps | 0.15 | Additional constraints |

**Score Mapping:**
| Total Contribution | Score | Time Estimate |
|-------------------|-------|---------------|
| 0.0 - 2.0 | 1-2 | 5-15 minutes |
| 2.1 - 4.0 | 3-4 | 15-30 minutes |
| 4.1 - 6.0 | 5-6 | 30-60 minutes |
| 6.1 - 8.0 | 7-8 | 1-2 hours |
| 8.1 - 10.0 | 9-10 | 2-4 hours |

### Input/Output Format

**Input (CircularDependencyInfo with ImportTraces):**
```json
{
  "cycle": ["@mono/ui", "@mono/api", "@mono/core", "@mono/ui"],
  "type": "indirect",
  "depth": 3,
  "complexity": 5,
  "importTraces": [
    {"fromPackage": "@mono/ui", "toPackage": "@mono/api", "filePath": "packages/ui/src/client.ts", "lineNumber": 5},
    {"fromPackage": "@mono/api", "toPackage": "@mono/core", "filePath": "packages/api/src/service.ts", "lineNumber": 10},
    {"fromPackage": "@mono/core", "toPackage": "@mono/ui", "filePath": "packages/core/src/render.ts", "lineNumber": 3}
  ]
}
```

**Output (CircularDependencyInfo with RefactoringComplexity):**
```json
{
  "cycle": ["@mono/ui", "@mono/api", "@mono/core", "@mono/ui"],
  "type": "indirect",
  "depth": 3,
  "complexity": 5,
  "refactoringComplexity": {
    "score": 5,
    "estimatedTime": "30-60 minutes",
    "breakdown": {
      "filesAffected": {
        "value": 3,
        "weight": 0.25,
        "contribution": 1.5,
        "description": "3 source files need modification"
      },
      "importsToChange": {
        "value": 3,
        "weight": 0.20,
        "contribution": 1.2,
        "description": "3 import statements need updating"
      },
      "chainDepth": {
        "value": 3,
        "weight": 0.25,
        "contribution": 1.5,
        "description": "Dependency chain has 3 levels"
      },
      "packagesInvolved": {
        "value": 3,
        "weight": 0.15,
        "contribution": 0.9,
        "description": "3 packages involved in cycle"
      },
      "externalDependencies": {
        "value": 0,
        "weight": 0.15,
        "contribution": 0.15,
        "description": "No external dependencies in cycle"
      }
    },
    "explanation": "Moderate refactoring: 3 files, 3 imports, 3-level chain"
  },
  "importTraces": [...],
  "fixStrategies": [...]
}
```

### Test Scenarios

| Scenario | Files | Imports | Depth | Packages | External | Expected Score |
|----------|-------|---------|-------|----------|----------|----------------|
| Simple direct | 2 | 2 | 2 | 2 | No | 2-3 |
| Medium 3-pkg | 4 | 3 | 3 | 3 | No | 4-5 |
| Complex 5-pkg | 8 | 6 | 5 | 5 | No | 6-7 |
| External deps | 4 | 3 | 3 | 3 | Yes | 5-6 |
| Large cycle | 15 | 12 | 8 | 8 | Yes | 9-10 |

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story 3.5]
- [Source: _bmad-output/planning-artifacts/prd.md#FR12]
- [Source: _bmad-output/project-context.md#Language-Specific Rules]
- [Source: _bmad-output/implementation-artifacts/3-3-generate-fix-strategy-recommendations.md]
- [Source: _bmad-output/implementation-artifacts/3-4-create-step-by-step-fix-guides.md]
- [Technical Debt Metrics - SonarQube](https://docs.sonarqube.org/latest/user-guide/metric-definitions/)

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List

