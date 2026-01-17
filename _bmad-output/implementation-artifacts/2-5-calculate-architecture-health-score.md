# Story 2.5: Calculate Architecture Health Score

Status: ready-for-dev

## Story

As a **user**,
I want **to see an overall health score (0-100) for my monorepo architecture**,
So that **I can quickly assess the state of my dependency structure**.

## Acceptance Criteria

1. **AC1: Health Score Calculation**
   - Given complete analysis results (graph, cycles, conflicts)
   - When I calculate the health score
   - Then the score is calculated based on weighted factors:
     - Circular dependencies (weight: 40%) - heavily penalized
     - Version conflicts (weight: 25%)
     - Dependency depth (weight: 20%)
     - Package coupling metrics (weight: 15%)
   - And score is 0-100 where higher is better

2. **AC2: Score Breakdown**
   - Given a calculated health score
   - When results are returned
   - Then the breakdown shows contribution of each factor:
     - `circularScore` - penalty from circular dependencies
     - `conflictScore` - penalty from version conflicts
     - `depthScore` - penalty from deep dependency chains
     - `couplingScore` - penalty from high coupling
   - And each factor shows its individual score (0-100)
   - And each factor shows its weighted contribution to overall score

3. **AC3: Rating Classification**
   - Given a health score
   - When rating is determined
   - Then classification follows thresholds:
     - `excellent` (85-100) - Well-maintained architecture
     - `good` (70-84) - Minor improvements possible
     - `fair` (50-69) - Attention needed
     - `poor` (30-49) - Significant issues
     - `critical` (0-29) - Immediate action required

4. **AC4: Circular Dependency Scoring**
   - Given circular dependencies from Story 2.3
   - When circular score is calculated
   - Then:
     - 0 cycles = 100 points
     - Each cycle deducts points based on severity
     - Direct cycles (2 packages) deduct more than indirect
     - Self-loops are heavily penalized
   - And formula: `100 - (directCycles * 15 + indirectCycles * 10 + selfLoops * 25)`

5. **AC5: Version Conflict Scoring**
   - Given version conflicts from Story 2.4
   - When conflict score is calculated
   - Then:
     - 0 conflicts = 100 points
     - Critical conflicts (major version) deduct 10 points each
     - Warning conflicts (minor version) deduct 5 points each
     - Info conflicts (patch version) deduct 2 points each
   - And score is capped at minimum 0

6. **AC6: Dependency Depth Scoring**
   - Given the dependency graph from Story 2.2
   - When depth score is calculated
   - Then:
     - Calculate max dependency depth (longest path)
     - Calculate average dependency depth
     - Optimal depth: 3-4 levels = 100 points
     - Each level above optimal deducts points
   - And formula considers both max and average depth

7. **AC7: Coupling Metrics Scoring**
   - Given the dependency graph
   - When coupling score is calculated
   - Then metrics include:
     - Afferent coupling (Ca) - packages depending on this package
     - Efferent coupling (Ce) - packages this package depends on
     - Instability = Ce / (Ca + Ce)
   - And high coupling packages are identified
   - And overall coupling score reflects average instability

8. **AC8: Performance Requirements**
   - Given a workspace with 100 packages
   - When health score calculation completes
   - Then it finishes in < 100ms

## Tasks / Subtasks

- [ ] **Task 1: Define HealthScore Types in Go** (AC: #2)
  - [ ] 1.1 Create `pkg/types/health_score.go`:
    ```go
    package types

    // HealthScoreResult represents the complete health score with breakdown.
    // Matches @monoguard/types HealthScore.
    type HealthScoreResult struct {
        Overall    int              `json:"overall"`    // 0-100
        Rating     HealthRating     `json:"rating"`     // excellent, good, fair, poor, critical
        Breakdown  *ScoreBreakdown  `json:"breakdown"`
        Factors    []*HealthFactor  `json:"factors"`
        UpdatedAt  string           `json:"updatedAt"`  // ISO 8601
    }

    // ScoreBreakdown shows individual factor scores
    type ScoreBreakdown struct {
        CircularScore int `json:"circularScore"` // 0-100
        ConflictScore int `json:"conflictScore"` // 0-100
        DepthScore    int `json:"depthScore"`    // 0-100
        CouplingScore int `json:"couplingScore"` // 0-100
    }

    // HealthFactor represents a single factor in the health calculation
    type HealthFactor struct {
        Name            string   `json:"name"`
        Score           int      `json:"score"`           // 0-100
        Weight          float64  `json:"weight"`          // 0.0-1.0
        WeightedScore   int      `json:"weightedScore"`   // score * weight
        Description     string   `json:"description"`
        Recommendations []string `json:"recommendations"`
    }

    // HealthRating classifies the overall score
    type HealthRating string

    const (
        HealthRatingExcellent HealthRating = "excellent" // 85-100
        HealthRatingGood      HealthRating = "good"      // 70-84
        HealthRatingFair      HealthRating = "fair"      // 50-69
        HealthRatingPoor      HealthRating = "poor"      // 30-49
        HealthRatingCritical  HealthRating = "critical"  // 0-29
    )

    // GetRating returns the rating for a given score
    func GetRating(score int) HealthRating {
        switch {
        case score >= 85:
            return HealthRatingExcellent
        case score >= 70:
            return HealthRatingGood
        case score >= 50:
            return HealthRatingFair
        case score >= 30:
            return HealthRatingPoor
        default:
            return HealthRatingCritical
        }
    }
    ```
  - [ ] 1.2 Add JSON serialization tests
  - [ ] 1.3 Add GetRating tests for all thresholds

- [ ] **Task 2: Implement Health Calculator** (AC: #1, #2)
  - [ ] 2.1 Create `pkg/analyzer/health_calculator.go`:
    ```go
    package analyzer

    import "github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"

    // Weight constants for health score factors
    const (
        WeightCircular = 0.40
        WeightConflict = 0.25
        WeightDepth    = 0.20
        WeightCoupling = 0.15
    )

    // HealthCalculator computes architecture health scores
    type HealthCalculator struct {
        graph     *types.DependencyGraph
        cycles    []*types.CircularDependencyInfo
        conflicts []*types.VersionConflict
    }

    // NewHealthCalculator creates a new calculator with analysis results
    func NewHealthCalculator(
        graph *types.DependencyGraph,
        cycles []*types.CircularDependencyInfo,
        conflicts []*types.VersionConflict,
    ) *HealthCalculator

    // Calculate computes the complete health score
    func (hc *HealthCalculator) Calculate() *types.HealthScoreResult

    // calculateCircularScore computes score from circular dependencies
    func (hc *HealthCalculator) calculateCircularScore() (int, *types.HealthFactor)

    // calculateConflictScore computes score from version conflicts
    func (hc *HealthCalculator) calculateConflictScore() (int, *types.HealthFactor)

    // calculateDepthScore computes score from dependency depth
    func (hc *HealthCalculator) calculateDepthScore() (int, *types.HealthFactor)

    // calculateCouplingScore computes score from package coupling
    func (hc *HealthCalculator) calculateCouplingScore() (int, *types.HealthFactor)
    ```
  - [ ] 2.2 Implement weighted score aggregation
  - [ ] 2.3 Create comprehensive tests

- [ ] **Task 3: Implement Circular Score Calculation** (AC: #4)
  - [ ] 3.1 Implement `calculateCircularScore`:
    ```go
    func (hc *HealthCalculator) calculateCircularScore() (int, *types.HealthFactor) {
        if len(hc.cycles) == 0 {
            return 100, &types.HealthFactor{
                Name:        "Circular Dependencies",
                Score:       100,
                Weight:      WeightCircular,
                Description: "No circular dependencies detected",
            }
        }

        deductions := 0
        directCount := 0
        indirectCount := 0
        selfLoopCount := 0

        for _, cycle := range hc.cycles {
            switch {
            case cycle.Depth == 1: // Self-loop
                selfLoopCount++
                deductions += 25
            case cycle.Type == types.CircularTypeDirect:
                directCount++
                deductions += 15
            default: // Indirect
                indirectCount++
                deductions += 10
            }
        }

        score := max(0, 100-deductions)
        recommendations := generateCircularRecommendations(directCount, indirectCount, selfLoopCount)

        return score, &types.HealthFactor{
            Name:            "Circular Dependencies",
            Score:           score,
            Weight:          WeightCircular,
            Description:     fmt.Sprintf("%d cycles detected", len(hc.cycles)),
            Recommendations: recommendations,
        }
    }
    ```
  - [ ] 3.2 Add tests for various cycle scenarios

- [ ] **Task 4: Implement Conflict Score Calculation** (AC: #5)
  - [ ] 4.1 Implement `calculateConflictScore`:
    ```go
    func (hc *HealthCalculator) calculateConflictScore() (int, *types.HealthFactor) {
        if len(hc.conflicts) == 0 {
            return 100, &types.HealthFactor{
                Name:        "Version Conflicts",
                Score:       100,
                Weight:      WeightConflict,
                Description: "No version conflicts detected",
            }
        }

        deductions := 0
        criticalCount := 0
        warningCount := 0
        infoCount := 0

        for _, conflict := range hc.conflicts {
            switch conflict.Severity {
            case types.ConflictSeverityCritical:
                criticalCount++
                deductions += 10
            case types.ConflictSeverityWarning:
                warningCount++
                deductions += 5
            case types.ConflictSeverityInfo:
                infoCount++
                deductions += 2
            }
        }

        score := max(0, 100-deductions)
        recommendations := generateConflictRecommendations(criticalCount, warningCount, infoCount)

        return score, &types.HealthFactor{
            Name:            "Version Conflicts",
            Score:           score,
            Weight:          WeightConflict,
            Description:     fmt.Sprintf("%d conflicts detected", len(hc.conflicts)),
            Recommendations: recommendations,
        }
    }
    ```
  - [ ] 4.2 Add tests for various conflict scenarios

- [ ] **Task 5: Implement Depth Score Calculation** (AC: #6)
  - [ ] 5.1 Implement `calculateDepthScore`:
    ```go
    func (hc *HealthCalculator) calculateDepthScore() (int, *types.HealthFactor) {
        maxDepth, avgDepth := hc.calculateDepthMetrics()

        // Optimal depth is 3-4 levels
        const optimalDepth = 4

        score := 100
        if maxDepth > optimalDepth {
            // Deduct 10 points per level above optimal
            score -= (maxDepth - optimalDepth) * 10
        }
        // Also consider average depth
        if avgDepth > float64(optimalDepth) {
            score -= int((avgDepth - float64(optimalDepth)) * 5)
        }

        score = max(0, score)

        return score, &types.HealthFactor{
            Name:        "Dependency Depth",
            Score:       score,
            Weight:      WeightDepth,
            Description: fmt.Sprintf("Max depth: %d, Avg depth: %.1f", maxDepth, avgDepth),
            Recommendations: generateDepthRecommendations(maxDepth, avgDepth),
        }
    }

    // calculateDepthMetrics computes max and average dependency depth
    func (hc *HealthCalculator) calculateDepthMetrics() (maxDepth int, avgDepth float64) {
        // Use BFS/DFS to calculate depth for each package
        // Return max depth and average depth across all packages
    }
    ```
  - [ ] 5.2 Implement BFS-based depth calculation
  - [ ] 5.3 Add tests for various graph structures

- [ ] **Task 6: Implement Coupling Score Calculation** (AC: #7)
  - [ ] 6.1 Implement `calculateCouplingScore`:
    ```go
    func (hc *HealthCalculator) calculateCouplingScore() (int, *types.HealthFactor) {
        couplingMetrics := hc.calculateCouplingMetrics()

        // Average instability across all packages
        avgInstability := couplingMetrics.AverageInstability

        // Score: 100 for 0.5 instability (balanced), decreases as it deviates
        // Very stable (0.0) or very unstable (1.0) packages are concerning
        deviation := math.Abs(avgInstability - 0.5)
        score := int(100 - (deviation * 100))
        score = max(0, min(100, score))

        return score, &types.HealthFactor{
            Name:        "Package Coupling",
            Score:       score,
            Weight:      WeightCoupling,
            Description: fmt.Sprintf("Avg instability: %.2f", avgInstability),
            Recommendations: generateCouplingRecommendations(couplingMetrics),
        }
    }

    // CouplingMetrics holds coupling analysis results
    type CouplingMetrics struct {
        AverageInstability float64
        HighCoupling       []string // Packages with concerning coupling
        PackageMetrics     map[string]*PackageCoupling
    }

    // PackageCoupling holds Ca, Ce, and instability for a package
    type PackageCoupling struct {
        AfferentCoupling int     // Ca - packages depending on this
        EfferentCoupling int     // Ce - packages this depends on
        Instability      float64 // Ce / (Ca + Ce)
    }
    ```
  - [ ] 6.2 Implement coupling metrics calculation
  - [ ] 6.3 Add tests for coupling scenarios

- [ ] **Task 7: Wire to Analyzer** (AC: all)
  - [ ] 7.1 Update `pkg/analyzer/analyzer.go`:
    ```go
    func (a *Analyzer) Analyze(workspace *types.WorkspaceData) (*types.AnalysisResult, error) {
        // Build dependency graph (Story 2.2)
        graph, err := a.graphBuilder.Build(workspace)
        if err != nil {
            return nil, err
        }

        // Detect circular dependencies (Story 2.3)
        cycleDetector := NewCycleDetector(graph)
        cycles := cycleDetector.DetectCycles()

        // Detect version conflicts (Story 2.4)
        conflictDetector := NewConflictDetector(workspace)
        conflicts := conflictDetector.DetectConflicts()

        // Calculate health score (Story 2.5)
        healthCalc := NewHealthCalculator(graph, cycles, conflicts)
        healthScore := healthCalc.Calculate()

        return &types.AnalysisResult{
            HealthScore:          healthScore.Overall,
            HealthScoreResult:    healthScore,
            Packages:             len(graph.Nodes),
            CircularDependencies: cycles,
            VersionConflicts:     conflicts,
            Graph:                graph,
            CreatedAt:            time.Now().UTC().Format(time.RFC3339),
        }, nil
    }
    ```
  - [ ] 7.2 Update AnalysisResult type to include HealthScoreResult
  - [ ] 7.3 Update handler and WASM tests

- [ ] **Task 8: Performance Testing** (AC: #8)
  - [ ] 8.1 Create `pkg/analyzer/health_calculator_benchmark_test.go`:
    ```go
    func BenchmarkCalculateHealth100Packages(b *testing.B) {
        graph := generateGraph(100)
        cycles := generateCycles(5)
        conflicts := generateConflicts(10)

        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            calc := NewHealthCalculator(graph, cycles, conflicts)
            calc.Calculate()
        }
    }
    ```
  - [ ] 8.2 Verify calculation completes in < 100ms

- [ ] **Task 9: Integration Verification** (AC: all)
  - [ ] 9.1 Build WASM: `pnpm nx build @monoguard/analysis-engine`
  - [ ] 9.2 Update smoke test to verify health score
  - [ ] 9.3 Test with various scenarios (healthy, unhealthy, critical)
  - [ ] 9.4 Verify all tests pass: `make test`

## Dev Notes

### Architecture Patterns & Constraints

**From Architecture Document:**
- **Location:** Health calculator in `pkg/analyzer/`
- **Input:** DependencyGraph, CircularDependencyInfo[], VersionConflict[]
- **Output:** HealthScoreResult with breakdown and factors

**Score Formula:**
```
Overall = (CircularScore × 0.40) + (ConflictScore × 0.25) + (DepthScore × 0.20) + (CouplingScore × 0.15)
```

**Critical Constraints:**
- **camelCase JSON:** All struct tags MUST use camelCase
- **Score Range:** All scores MUST be 0-100, never negative or >100
- **Deterministic:** Same input MUST produce same score
- **Fast:** < 100ms for 100 packages

### Critical Don't-Miss Rules

**From project-context.md:**

1. **Score Bounds:**
   ```go
   // ✅ CORRECT: Always bound scores
   score := max(0, min(100, calculatedScore))

   // ❌ WRONG: Unbounded score
   score := 100 - deductions // Could go negative!
   ```

2. **Weighted Calculation:**
   ```go
   // ✅ CORRECT: Use float for weights, convert at end
   weighted := float64(circularScore) * WeightCircular +
               float64(conflictScore) * WeightConflict +
               float64(depthScore) * WeightDepth +
               float64(couplingScore) * WeightCoupling
   overall := int(math.Round(weighted))

   // ❌ WRONG: Integer division loses precision
   overall := (circularScore * 40 + conflictScore * 25) / 100
   ```

3. **Rating Thresholds (must be consistent):**
   ```go
   // Exact thresholds - no gaps, no overlaps
   // 85-100 = excellent
   // 70-84  = good
   // 50-69  = fair
   // 30-49  = poor
   // 0-29   = critical
   ```

### Project Structure Notes

**Target Directory Structure:**
```
packages/analysis-engine/
├── pkg/
│   ├── analyzer/
│   │   ├── analyzer.go                        # UPDATE: Add health calculation
│   │   ├── analyzer_test.go                   # UPDATE
│   │   ├── graph_builder.go                   # From Story 2.2
│   │   ├── cycle_detector.go                  # From Story 2.3
│   │   ├── conflict_detector.go               # From Story 2.4
│   │   ├── health_calculator.go               # NEW: Health score calculation
│   │   ├── health_calculator_test.go          # NEW: Calculator tests
│   │   ├── health_calculator_benchmark_test.go # NEW: Performance tests
│   │   ├── depth_calculator.go                # NEW: Depth metrics
│   │   ├── coupling_calculator.go             # NEW: Coupling metrics
│   │   └── recommendations.go                 # NEW: Recommendation generators
│   └── types/
│       ├── types.go                           # UPDATE: Add HealthScoreResult to AnalysisResult
│       ├── health_score.go                    # NEW: HealthScoreResult types
│       └── health_score_test.go               # NEW: Type tests
└── ...
```

### Input/Output Format

**Input (from previous stories):**
- `graph`: DependencyGraph with nodes and edges
- `cycles`: CircularDependencyInfo[] with severity info
- `conflicts`: VersionConflict[] with severity info

**Output (HealthScoreResult):**
```json
{
  "overall": 72,
  "rating": "good",
  "breakdown": {
    "circularScore": 70,
    "conflictScore": 85,
    "depthScore": 60,
    "couplingScore": 75
  },
  "factors": [
    {
      "name": "Circular Dependencies",
      "score": 70,
      "weight": 0.4,
      "weightedScore": 28,
      "description": "2 cycles detected",
      "recommendations": ["Break cycle between pkg-a and pkg-b by extracting shared code"]
    },
    {
      "name": "Version Conflicts",
      "score": 85,
      "weight": 0.25,
      "weightedScore": 21,
      "description": "3 conflicts detected",
      "recommendations": ["Upgrade lodash to ^4.17.21 in all packages"]
    },
    {
      "name": "Dependency Depth",
      "score": 60,
      "weight": 0.2,
      "weightedScore": 12,
      "description": "Max depth: 7, Avg depth: 4.2",
      "recommendations": ["Consider flattening deep dependency chains"]
    },
    {
      "name": "Package Coupling",
      "score": 75,
      "weight": 0.15,
      "weightedScore": 11,
      "description": "Avg instability: 0.45",
      "recommendations": []
    }
  ],
  "updatedAt": "2026-01-17T10:30:00Z"
}
```

### Scoring Formulas

**Circular Score:**
```
deductions = (selfLoops × 25) + (directCycles × 15) + (indirectCycles × 10)
score = max(0, 100 - deductions)
```

**Conflict Score:**
```
deductions = (critical × 10) + (warning × 5) + (info × 2)
score = max(0, 100 - deductions)
```

**Depth Score:**
```
optimalDepth = 4
maxPenalty = (maxDepth - optimalDepth) × 10  // if maxDepth > optimal
avgPenalty = (avgDepth - optimalDepth) × 5   // if avgDepth > optimal
score = max(0, 100 - maxPenalty - avgPenalty)
```

**Coupling Score (Instability-based):**
```
For each package:
  Ca = afferent coupling (incoming dependencies)
  Ce = efferent coupling (outgoing dependencies)
  Instability = Ce / (Ca + Ce)  // 0.0 = stable, 1.0 = unstable

avgInstability = average of all package instabilities
deviation = |avgInstability - 0.5|  // 0.5 is ideal balance
score = 100 - (deviation × 100)
```

### Test Scenarios

| Scenario | Cycles | Conflicts | Max Depth | Expected Score |
|----------|--------|-----------|-----------|----------------|
| Perfect | 0 | 0 | 3 | ~100 (excellent) |
| Good | 1 indirect | 2 info | 4 | ~80 (good) |
| Fair | 2 direct | 5 warning | 6 | ~55 (fair) |
| Poor | 3 cycles | 3 critical | 8 | ~35 (poor) |
| Critical | 5+ cycles | 5+ critical | 10+ | <30 (critical) |

### Previous Story Intelligence

**From Story 2.2 (ready-for-dev):**
- DependencyGraph provides nodes and edges for depth/coupling calculation
- `graph.Nodes[name].Dependencies` shows internal dependencies

**From Story 2.3 (ready-for-dev):**
- CircularDependencyInfo has `Type` (direct/indirect) and `Depth`
- Cycles are already classified for scoring

**From Story 2.4 (ready-for-dev):**
- VersionConflict has `Severity` (critical/warning/info)
- Conflicts are already classified for scoring

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story 2.5]
- [Source: packages/types/src/domain.ts#HealthScore]
- [Source: _bmad-output/project-context.md#Result Type Pattern]
- [Software Metrics - Coupling](https://en.wikipedia.org/wiki/Coupling_(computer_programming))
- [Martin Metrics - Instability](https://en.wikipedia.org/wiki/Software_package_metrics)

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
