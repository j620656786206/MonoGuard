# Story 3.6: Generate Impact Assessment

Status: ready-for-dev

## Story

As a **user**,
I want **to see how many packages are affected by each circular dependency**,
So that **I can understand the blast radius and prioritize high-impact fixes**.

## Acceptance Criteria

1. **AC1: Direct Participants**
   - Given a circular dependency
   - When I request impact assessment
   - Then I see the direct participants:
     - List of packages directly in the cycle
     - Count of direct participants
   - And packages are ordered as they appear in the cycle

2. **AC2: Indirect Dependents (Ripple Effect)**
   - Given packages in a cycle
   - When calculating indirect dependents
   - Then identify all packages that depend on cycle participants:
     - First-level dependents (directly depend on cycle packages)
     - Transitive dependents (depend on first-level dependents)
   - And each dependent includes which cycle package it depends on
   - And duplicates are counted only once

3. **AC3: Total Affected Package Count**
   - Given direct and indirect dependents
   - When calculating total affected
   - Then return:
     - Direct participant count
     - Indirect dependent count
     - Total affected count (direct + indirect, no duplicates)
   - And total is never greater than total packages in workspace

4. **AC4: Percentage of Monorepo Affected**
   - Given total affected packages and workspace size
   - When calculating percentage
   - Then return:
     - Percentage as decimal (0.0-1.0) and display string ("25%")
     - Percentage relative to total workspace packages
   - And handle edge cases (empty workspace, all packages affected)

5. **AC5: Risk Level Classification**
   - Given impact metrics
   - When determining risk level
   - Then classify as:
     - **Critical**: > 50% of packages affected OR cycle includes core/shared package
     - **High**: 25-50% of packages affected
     - **Medium**: 10-25% of packages affected
     - **Low**: < 10% of packages affected
   - And risk level considers both percentage and package naming patterns

6. **AC6: Ripple Effect Visualization Data**
   - Given impact assessment
   - When preparing visualization data
   - Then include:
     - Direct participants with visual markers
     - Indirect dependents grouped by distance from cycle
     - Dependency paths from cycle to each affected package
   - And data is structured for D3.js force-directed graph (Epic 4)

7. **AC7: Integration with CircularDependencyInfo**
   - Given analysis results
   - When enriching CircularDependencyInfo
   - Then add `impactAssessment` field:
     ```go
     type CircularDependencyInfo struct {
         // ... existing fields ...
         ImpactAssessment *ImpactAssessment `json:"impactAssessment,omitempty"`
     }
     ```
   - And existing consumers continue working (backward compatible)

8. **AC8: Performance**
   - Given a workspace with 100 packages and 5 cycles
   - When calculating impact for all cycles
   - Then calculation completes in < 200ms additional overhead
   - And memory usage increase is < 10MB

## Tasks / Subtasks

- [ ] **Task 1: Define ImpactAssessment Types** (AC: #1, #2, #3, #4, #5, #6)
  - [ ] 1.1 Create `pkg/types/impact_assessment.go`:
    ```go
    package types

    // ImpactAssessment represents the blast radius analysis for a circular dependency.
    // Matches @monoguard/types ImpactAssessment interface.
    type ImpactAssessment struct {
        // DirectParticipants are packages directly in the cycle
        DirectParticipants []string `json:"directParticipants"`

        // IndirectDependents are packages that depend on cycle participants
        IndirectDependents []IndirectDependent `json:"indirectDependents"`

        // TotalAffected is the count of all affected packages (direct + indirect)
        TotalAffected int `json:"totalAffected"`

        // AffectedPercentage is the proportion of workspace affected (0.0-1.0)
        AffectedPercentage float64 `json:"affectedPercentage"`

        // AffectedPercentageDisplay is human-readable (e.g., "25%")
        AffectedPercentageDisplay string `json:"affectedPercentageDisplay"`

        // RiskLevel classifies the impact severity
        RiskLevel RiskLevel `json:"riskLevel"`

        // RiskExplanation describes why this risk level was assigned
        RiskExplanation string `json:"riskExplanation"`

        // RippleEffect contains visualization-ready data
        RippleEffect *RippleEffect `json:"rippleEffect,omitempty"`
    }

    // IndirectDependent represents a package that depends on a cycle participant.
    type IndirectDependent struct {
        // PackageName is the affected package
        PackageName string `json:"packageName"`

        // DependsOn is the cycle participant this package depends on
        DependsOn string `json:"dependsOn"`

        // Distance is the number of hops from the cycle (1 = direct dependent)
        Distance int `json:"distance"`

        // DependencyPath shows the full path from cycle to this package
        DependencyPath []string `json:"dependencyPath"`
    }

    // RippleEffect contains data for visualization.
    type RippleEffect struct {
        // Layers groups affected packages by distance from cycle
        Layers []RippleLayer `json:"layers"`

        // TotalLayers is the maximum distance from cycle
        TotalLayers int `json:"totalLayers"`
    }

    // RippleLayer represents packages at a specific distance from the cycle.
    type RippleLayer struct {
        // Distance from the cycle (0 = direct participants, 1 = first-level dependents)
        Distance int `json:"distance"`

        // Packages at this distance
        Packages []string `json:"packages"`

        // Count of packages at this layer
        Count int `json:"count"`
    }

    // RiskLevel classifies the impact severity.
    type RiskLevel string

    const (
        RiskLevelCritical RiskLevel = "critical"
        RiskLevelHigh     RiskLevel = "high"
        RiskLevelMedium   RiskLevel = "medium"
        RiskLevelLow      RiskLevel = "low"
    )
    ```
  - [ ] 1.2 Add JSON serialization tests in `pkg/types/impact_assessment_test.go`
  - [ ] 1.3 Ensure all JSON tags use camelCase

- [ ] **Task 2: Create ImpactAnalyzer** (AC: #1, #2, #7)
  - [ ] 2.1 Create `pkg/analyzer/impact_analyzer.go`:
    ```go
    package analyzer

    import "github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"

    // ImpactAnalyzer calculates the blast radius of circular dependencies.
    type ImpactAnalyzer struct {
        graph     *types.DependencyGraph
        workspace *types.WorkspaceData
        // reverseDeps maps package -> packages that depend on it
        reverseDeps map[string][]string
    }

    // NewImpactAnalyzer creates a new analyzer.
    func NewImpactAnalyzer(graph *types.DependencyGraph, workspace *types.WorkspaceData) *ImpactAnalyzer

    // Analyze calculates the impact assessment for a cycle.
    func (ia *ImpactAnalyzer) Analyze(cycle *types.CircularDependencyInfo) *types.ImpactAssessment

    // buildReverseDependencies creates a reverse lookup map.
    func (ia *ImpactAnalyzer) buildReverseDependencies()

    // getDirectParticipants extracts unique packages from cycle.
    func (ia *ImpactAnalyzer) getDirectParticipants(cycle *types.CircularDependencyInfo) []string

    // findIndirectDependents finds all packages depending on cycle participants.
    func (ia *ImpactAnalyzer) findIndirectDependents(directParticipants []string) []types.IndirectDependent

    // calculateRiskLevel determines risk based on metrics.
    func (ia *ImpactAnalyzer) calculateRiskLevel(
        affected int,
        total int,
        directParticipants []string,
    ) (types.RiskLevel, string)

    // buildRippleEffect creates visualization data.
    func (ia *ImpactAnalyzer) buildRippleEffect(
        directParticipants []string,
        indirectDependents []types.IndirectDependent,
    ) *types.RippleEffect
    ```
  - [ ] 2.2 Implement impact analysis logic
  - [ ] 2.3 Create comprehensive tests in `pkg/analyzer/impact_analyzer_test.go`

- [ ] **Task 3: Implement Reverse Dependency Building** (AC: #2)
  - [ ] 3.1 Implement `buildReverseDependencies`:
    ```go
    func (ia *ImpactAnalyzer) buildReverseDependencies() {
        ia.reverseDeps = make(map[string][]string)

        // Initialize all packages
        for pkgName := range ia.graph.Nodes {
            ia.reverseDeps[pkgName] = []string{}
        }

        // Build reverse mapping from edges
        for _, edge := range ia.graph.Edges {
            // edge.From depends on edge.To
            // So edge.From is a "dependent" of edge.To
            ia.reverseDeps[edge.To] = append(ia.reverseDeps[edge.To], edge.From)
        }
    }
    ```
  - [ ] 3.2 Add tests for reverse dependency building

- [ ] **Task 4: Implement Direct Participants Extraction** (AC: #1)
  - [ ] 4.1 Implement `getDirectParticipants`:
    ```go
    func (ia *ImpactAnalyzer) getDirectParticipants(cycle *types.CircularDependencyInfo) []string {
        // Cycle array ends with first element, so exclude last
        if len(cycle.Cycle) == 0 {
            return []string{}
        }

        // Use map to dedupe (shouldn't be needed but safe)
        seen := make(map[string]bool)
        participants := []string{}

        for i := 0; i < len(cycle.Cycle)-1; i++ {
            pkg := cycle.Cycle[i]
            if !seen[pkg] {
                seen[pkg] = true
                participants = append(participants, pkg)
            }
        }

        return participants
    }
    ```
  - [ ] 4.2 Add tests for direct participants extraction

- [ ] **Task 5: Implement Indirect Dependent Finding** (AC: #2)
  - [ ] 5.1 Implement `findIndirectDependents` with BFS:
    ```go
    func (ia *ImpactAnalyzer) findIndirectDependents(directParticipants []string) []types.IndirectDependent {
        // Mark direct participants as visited
        visited := make(map[string]bool)
        for _, pkg := range directParticipants {
            visited[pkg] = true
        }

        result := []types.IndirectDependent{}

        // BFS queue: (package, dependsOnCyclePackage, distance, path)
        type queueItem struct {
            pkg       string
            dependsOn string
            distance  int
            path      []string
        }

        queue := []queueItem{}

        // Initialize queue with direct dependents of cycle participants
        for _, cyclePkg := range directParticipants {
            for _, dependent := range ia.reverseDeps[cyclePkg] {
                if !visited[dependent] {
                    queue = append(queue, queueItem{
                        pkg:       dependent,
                        dependsOn: cyclePkg,
                        distance:  1,
                        path:      []string{cyclePkg, dependent},
                    })
                }
            }
        }

        // BFS traversal
        for len(queue) > 0 {
            item := queue[0]
            queue = queue[1:]

            if visited[item.pkg] {
                continue
            }
            visited[item.pkg] = true

            result = append(result, types.IndirectDependent{
                PackageName:    item.pkg,
                DependsOn:      item.dependsOn,
                Distance:       item.distance,
                DependencyPath: item.path,
            })

            // Add this package's dependents to queue
            for _, nextDependent := range ia.reverseDeps[item.pkg] {
                if !visited[nextDependent] {
                    newPath := make([]string, len(item.path)+1)
                    copy(newPath, item.path)
                    newPath[len(item.path)] = nextDependent

                    queue = append(queue, queueItem{
                        pkg:       nextDependent,
                        dependsOn: item.dependsOn, // Original cycle package
                        distance:  item.distance + 1,
                        path:      newPath,
                    })
                }
            }
        }

        return result
    }
    ```
  - [ ] 5.2 Add tests for indirect dependent finding

- [ ] **Task 6: Implement Percentage Calculation** (AC: #4)
  - [ ] 6.1 Implement percentage calculation:
    ```go
    func calculatePercentage(affected, total int) (float64, string) {
        if total == 0 {
            return 0.0, "0%"
        }

        percentage := float64(affected) / float64(total)

        // Cap at 1.0
        if percentage > 1.0 {
            percentage = 1.0
        }

        // Format display string
        displayPercentage := int(percentage * 100)
        display := fmt.Sprintf("%d%%", displayPercentage)

        return percentage, display
    }
    ```
  - [ ] 6.2 Add tests for percentage calculation

- [ ] **Task 7: Implement Risk Level Classification** (AC: #5)
  - [ ] 7.1 Implement `calculateRiskLevel`:
    ```go
    func (ia *ImpactAnalyzer) calculateRiskLevel(
        affected int,
        total int,
        directParticipants []string,
    ) (types.RiskLevel, string) {
        percentage := float64(affected) / float64(total)

        // Check for core/shared package patterns
        hasCorePackage := false
        corePatterns := []string{"core", "common", "shared", "utils", "lib"}
        for _, pkg := range directParticipants {
            pkgLower := strings.ToLower(pkg)
            for _, pattern := range corePatterns {
                if strings.Contains(pkgLower, pattern) {
                    hasCorePackage = true
                    break
                }
            }
            if hasCorePackage {
                break
            }
        }

        // Critical: >50% affected OR core package involved
        if percentage > 0.50 || hasCorePackage {
            explanation := "Critical impact: "
            if hasCorePackage {
                explanation += "cycle includes core/shared package"
            } else {
                explanation += fmt.Sprintf("%.0f%% of packages affected", percentage*100)
            }
            return types.RiskLevelCritical, explanation
        }

        // High: 25-50%
        if percentage > 0.25 {
            return types.RiskLevelHigh, fmt.Sprintf("High impact: %.0f%% of packages affected", percentage*100)
        }

        // Medium: 10-25%
        if percentage > 0.10 {
            return types.RiskLevelMedium, fmt.Sprintf("Medium impact: %.0f%% of packages affected", percentage*100)
        }

        // Low: <10%
        return types.RiskLevelLow, fmt.Sprintf("Low impact: %.0f%% of packages affected", percentage*100)
    }
    ```
  - [ ] 7.2 Add tests for risk level classification

- [ ] **Task 8: Implement Ripple Effect Builder** (AC: #6)
  - [ ] 8.1 Implement `buildRippleEffect`:
    ```go
    func (ia *ImpactAnalyzer) buildRippleEffect(
        directParticipants []string,
        indirectDependents []types.IndirectDependent,
    ) *types.RippleEffect {
        // Group by distance
        layerMap := make(map[int][]string)

        // Layer 0: direct participants
        layerMap[0] = directParticipants

        // Group indirect dependents by distance
        maxDistance := 0
        for _, dep := range indirectDependents {
            layerMap[dep.Distance] = append(layerMap[dep.Distance], dep.PackageName)
            if dep.Distance > maxDistance {
                maxDistance = dep.Distance
            }
        }

        // Build ordered layers
        layers := []types.RippleLayer{}
        for distance := 0; distance <= maxDistance; distance++ {
            packages := layerMap[distance]
            if len(packages) > 0 {
                layers = append(layers, types.RippleLayer{
                    Distance: distance,
                    Packages: packages,
                    Count:    len(packages),
                })
            }
        }

        return &types.RippleEffect{
            Layers:      layers,
            TotalLayers: maxDistance + 1,
        }
    }
    ```
  - [ ] 8.2 Add tests for ripple effect builder

- [ ] **Task 9: Implement Main Analyze Function** (AC: #1, #2, #3, #4, #5, #6)
  - [ ] 9.1 Implement `Analyze`:
    ```go
    func (ia *ImpactAnalyzer) Analyze(cycle *types.CircularDependencyInfo) *types.ImpactAssessment {
        // Build reverse dependencies if not done
        if ia.reverseDeps == nil {
            ia.buildReverseDependencies()
        }

        // Get direct participants
        directParticipants := ia.getDirectParticipants(cycle)

        // Find indirect dependents
        indirectDependents := ia.findIndirectDependents(directParticipants)

        // Calculate totals
        totalAffected := len(directParticipants) + len(indirectDependents)
        totalPackages := len(ia.graph.Nodes)

        // Calculate percentage
        percentage, displayPercentage := calculatePercentage(totalAffected, totalPackages)

        // Determine risk level
        riskLevel, riskExplanation := ia.calculateRiskLevel(
            totalAffected,
            totalPackages,
            directParticipants,
        )

        // Build ripple effect data
        rippleEffect := ia.buildRippleEffect(directParticipants, indirectDependents)

        return &types.ImpactAssessment{
            DirectParticipants:        directParticipants,
            IndirectDependents:        indirectDependents,
            TotalAffected:             totalAffected,
            AffectedPercentage:        percentage,
            AffectedPercentageDisplay: displayPercentage,
            RiskLevel:                 riskLevel,
            RiskExplanation:           riskExplanation,
            RippleEffect:              rippleEffect,
        }
    }
    ```
  - [ ] 9.2 Add comprehensive tests

- [ ] **Task 10: Integrate with CircularDependencyInfo** (AC: #7)
  - [ ] 10.1 Update `pkg/types/circular.go`:
    ```go
    type CircularDependencyInfo struct {
        Cycle                 []string               `json:"cycle"`
        Type                  CircularType           `json:"type"`
        Severity              CircularSeverity       `json:"severity"`
        Depth                 int                    `json:"depth"`
        Impact                string                 `json:"impact"`
        Complexity            int                    `json:"complexity"`
        RefactoringComplexity *RefactoringComplexity `json:"refactoringComplexity,omitempty"` // Story 3.5
        ImpactAssessment      *ImpactAssessment      `json:"impactAssessment,omitempty"`      // NEW Story 3.6
        RootCause             *RootCauseAnalysis     `json:"rootCause,omitempty"`
        ImportTraces          []ImportTrace          `json:"importTraces,omitempty"`
        FixStrategies         []FixStrategy          `json:"fixStrategies,omitempty"`
    }
    ```
  - [ ] 10.2 Verify existing tests still pass

- [ ] **Task 11: Wire to Analyzer Pipeline** (AC: all)
  - [ ] 11.1 Update `pkg/analyzer/analyzer.go`:
    ```go
    func (a *Analyzer) AnalyzeWithSources(...) (*types.AnalysisResult, error) {
        // ... existing analysis ...

        // NEW: Calculate impact assessment (Story 3.6)
        impactAnalyzer := NewImpactAnalyzer(graph, workspace)
        for _, cycle := range cycles {
            cycle.ImpactAssessment = impactAnalyzer.Analyze(cycle)
        }

        return result, nil
    }
    ```
  - [ ] 11.2 Update analyzer tests

- [ ] **Task 12: Update TypeScript Types** (AC: #7)
  - [ ] 12.1 Update `packages/types/src/analysis/results.ts`:
    ```typescript
    export interface ImpactAssessment {
      /** Packages directly in the cycle */
      directParticipants: string[];
      /** Packages that depend on cycle participants */
      indirectDependents: IndirectDependent[];
      /** Count of all affected packages */
      totalAffected: number;
      /** Proportion of workspace affected (0.0-1.0) */
      affectedPercentage: number;
      /** Human-readable percentage (e.g., "25%") */
      affectedPercentageDisplay: string;
      /** Impact severity classification */
      riskLevel: RiskLevel;
      /** Explanation of risk classification */
      riskExplanation: string;
      /** Visualization-ready data */
      rippleEffect?: RippleEffect;
    }

    export interface IndirectDependent {
      /** The affected package */
      packageName: string;
      /** Which cycle participant this package depends on */
      dependsOn: string;
      /** Hops from the cycle (1 = direct dependent) */
      distance: number;
      /** Full path from cycle to this package */
      dependencyPath: string[];
    }

    export interface RippleEffect {
      /** Packages grouped by distance from cycle */
      layers: RippleLayer[];
      /** Maximum distance from cycle */
      totalLayers: number;
    }

    export interface RippleLayer {
      /** Distance from cycle (0 = direct participants) */
      distance: number;
      /** Packages at this distance */
      packages: string[];
      /** Count of packages */
      count: number;
    }

    export type RiskLevel = 'critical' | 'high' | 'medium' | 'low';

    export interface CircularDependencyInfo {
      // ... existing fields ...
      /** Impact assessment with blast radius analysis (Story 3.6) */
      impactAssessment?: ImpactAssessment;
    }
    ```
  - [ ] 12.2 Run `pnpm nx build types` to verify
  - [ ] 12.3 Add type tests for ImpactAssessment

- [ ] **Task 13: Performance Testing** (AC: #8)
  - [ ] 13.1 Create `pkg/analyzer/impact_analyzer_benchmark_test.go`:
    ```go
    func BenchmarkImpactAnalysis(b *testing.B) {
        graph := generateGraphWithDependencies(100)
        workspace := generateWorkspace(100)
        cycles := generateCycles(5)
        analyzer := NewImpactAnalyzer(graph, workspace)

        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            for _, cycle := range cycles {
                analyzer.Analyze(cycle)
            }
        }
    }
    ```
  - [ ] 13.2 Verify < 200ms for 100 packages with 5 cycles
  - [ ] 13.3 Document actual performance in completion notes

- [ ] **Task 14: Integration Verification** (AC: all)
  - [ ] 14.1 Run all tests: `cd packages/analysis-engine && make test`
  - [ ] 14.2 Build WASM: `pnpm nx build @monoguard/analysis-engine`
  - [ ] 14.3 Run affected CI checks: `pnpm nx affected --target=lint,test,type-check --base=main`
  - [ ] 14.4 Verify JSON output includes impactAssessment field

## Dev Notes

### Architecture Patterns & Constraints

**From Architecture Document:**
- **Location:** Impact analyzer in `pkg/analyzer/`
- **Pattern:** Analyzer pattern with graph traversal
- **Integration:** Enriches CircularDependencyInfo with ImpactAssessment
- **Dependency:** Uses DependencyGraph for reverse dependency lookup

**Critical Constraints:**
- **camelCase JSON:** All struct tags MUST use camelCase
- **Optional Fields:** ImpactAssessment is `omitempty` - backward compatible
- **Performance:** Reverse dependency map should be built once and reused
- **Visualization Ready:** RippleEffect data must be compatible with D3.js

### Critical Don't-Miss Rules

**From project-context.md:**

1. **JSON Naming Convention:**
   ```go
   // CORRECT: camelCase JSON tags
   type ImpactAssessment struct {
       DirectParticipants        []string  `json:"directParticipants"`
       AffectedPercentageDisplay string    `json:"affectedPercentageDisplay"`
   }

   // WRONG: snake_case JSON tags
   type ImpactAssessment struct {
       DirectParticipants []string `json:"direct_participants"` // WRONG!
   }
   ```

2. **Enum Constants:**
   ```go
   // CORRECT: kebab-case string values for TypeScript compatibility
   const (
       RiskLevelCritical RiskLevel = "critical"
       RiskLevelHigh     RiskLevel = "high"
   )

   // WRONG: UPPER_CASE
   const (
       RiskLevelCritical RiskLevel = "CRITICAL" // WRONG!
   )
   ```

3. **Slice Initialization:**
   ```go
   // CORRECT: Initialize as empty slice
   directParticipants := []string{}

   // WRONG: Nil slice (serializes as null)
   var directParticipants []string // nil
   ```

### Project Structure Notes

**Target Directory Structure:**
```
packages/analysis-engine/
├── pkg/
│   ├── analyzer/
│   │   ├── analyzer.go                        # UPDATE: Add impact analysis
│   │   ├── impact_analyzer.go                 # NEW: Impact analyzer
│   │   ├── impact_analyzer_test.go            # NEW: Analyzer tests
│   │   └── impact_analyzer_benchmark_test.go  # NEW: Performance
│   └── types/
│       ├── circular.go                        # UPDATE: Add ImpactAssessment field
│       ├── impact_assessment.go               # NEW: ImpactAssessment types
│       └── impact_assessment_test.go          # NEW: Type tests
└── ...

packages/types/src/analysis/
├── results.ts                                 # UPDATE: Add TS types
└── ...
```

### Previous Story Intelligence

**From Story 3.5 (ready-for-dev):**
- RefactoringComplexity uses weighted factors - similar pattern for risk calculation
- ComplexityFactor has value, weight, contribution - consider similar breakdown
- Time estimation pattern - risk level provides similar prioritization signal

**From Story 3.3 (done):**
- Fix strategies use cycle.Depth and cycle.Cycle - reuse for direct participants
- Pattern matching for "core" packages - reuse for risk level classification
- Performance: ~5 microseconds per cycle - target similar performance

**From Story 2.2 (done):**
- DependencyGraph has Nodes and Edges - use Edges for reverse dependency building
- DependencyEdge has From and To - "From depends on To" semantic

### Algorithm Notes

**Reverse Dependency Building:**
- For each edge (A → B), A depends on B
- So B's "dependents" include A
- Build once, cache in analyzer

**BFS for Ripple Effect:**
- Start from cycle participants (distance 0)
- Find their dependents (distance 1)
- Continue until no new packages found
- Track visited to avoid infinite loops and duplicates

**Risk Level Thresholds:**
| Condition | Risk Level |
|-----------|------------|
| > 50% OR core package | Critical |
| 25% - 50% | High |
| 10% - 25% | Medium |
| < 10% | Low |

### Input/Output Format

**Input (CircularDependencyInfo):**
```json
{
  "cycle": ["@mono/ui", "@mono/api", "@mono/core", "@mono/ui"],
  "type": "indirect",
  "depth": 3
}
```

**Output (CircularDependencyInfo with ImpactAssessment):**
```json
{
  "cycle": ["@mono/ui", "@mono/api", "@mono/core", "@mono/ui"],
  "type": "indirect",
  "depth": 3,
  "impactAssessment": {
    "directParticipants": ["@mono/ui", "@mono/api", "@mono/core"],
    "indirectDependents": [
      {
        "packageName": "@mono/app",
        "dependsOn": "@mono/ui",
        "distance": 1,
        "dependencyPath": ["@mono/ui", "@mono/app"]
      },
      {
        "packageName": "@mono/dashboard",
        "dependsOn": "@mono/ui",
        "distance": 2,
        "dependencyPath": ["@mono/ui", "@mono/app", "@mono/dashboard"]
      }
    ],
    "totalAffected": 5,
    "affectedPercentage": 0.5,
    "affectedPercentageDisplay": "50%",
    "riskLevel": "high",
    "riskExplanation": "High impact: 50% of packages affected",
    "rippleEffect": {
      "layers": [
        {"distance": 0, "packages": ["@mono/ui", "@mono/api", "@mono/core"], "count": 3},
        {"distance": 1, "packages": ["@mono/app"], "count": 1},
        {"distance": 2, "packages": ["@mono/dashboard"], "count": 1}
      ],
      "totalLayers": 3
    }
  }
}
```

### Test Scenarios

| Scenario | Direct | Indirect | Total | Percentage | Expected Risk |
|----------|--------|----------|-------|------------|---------------|
| Isolated cycle | 2 | 0 | 2 | 20% (10 pkg) | Medium |
| Core package cycle | 2 | 5 | 7 | 70% (10 pkg) | Critical (core) |
| High impact | 3 | 3 | 6 | 30% (20 pkg) | High |
| Low impact | 2 | 1 | 3 | 5% (60 pkg) | Low |
| Full cascade | 5 | 15 | 20 | 100% (20 pkg) | Critical |

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story 3.6]
- [Source: _bmad-output/planning-artifacts/prd.md#FR8]
- [Source: _bmad-output/project-context.md#Language-Specific Rules]
- [Source: _bmad-output/implementation-artifacts/3-3-generate-fix-strategy-recommendations.md]
- [Source: _bmad-output/implementation-artifacts/3-5-calculate-refactoring-complexity-scores.md]
- [BFS Algorithm - Wikipedia](https://en.wikipedia.org/wiki/Breadth-first_search)

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List

