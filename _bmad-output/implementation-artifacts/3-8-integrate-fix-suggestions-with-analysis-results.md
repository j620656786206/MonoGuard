# Story 3.8: Integrate Fix Suggestions with Analysis Results

Status: done

## Story

As a **user**,
I want **fix suggestions integrated into the main analysis results**,
So that **I can see problems and solutions together**.

## Acceptance Criteria

1. **AC1: Quick Fix Recommendation**
   - Given a circular dependency with multiple fix strategies
   - When I view analysis results
   - Then each circular dependency includes:
     - A `quickFix` field with the best strategy (highest suitability)
     - Quick access to the recommended fix's guide
     - One-line summary of what the fix accomplishes
   - And the quick fix is always the strategy with highest suitability score

2. **AC2: All Strategies Available**
   - Given analysis results with fix strategies
   - When I want to view all options
   - Then I can access:
     - Full list of fix strategies (already in `fixStrategies[]`)
     - Each strategy with guide, complexity, and before/after
     - Strategies sorted by suitability (best first)
   - And sorting is consistent across WASM and CLI output

3. **AC3: Prioritized Circular Dependencies**
   - Given multiple circular dependencies
   - When viewing analysis results
   - Then circular dependencies are sorted by:
     - Priority score = Impact × (11 - Complexity) (quick wins first)
     - Higher priority scores first
   - And each CircularDependencyInfo includes a `priorityScore` field
   - And results ordering is deterministic (stable sort)

4. **AC4: Aggregated Fix Summary**
   - Given analysis results with circular dependencies
   - When generating summary statistics
   - Then include in AnalysisResult:
     - `totalCircularDependencies`: count of all cycles
     - `totalEstimatedFixTime`: sum of all fix times (human-readable)
     - `quickWinsCount`: number of low-complexity fixes available
     - `criticalCyclesCount`: number of critical impact cycles
   - And summary is at the AnalysisResult level (not per-cycle)

5. **AC5: One-Click Guide Access**
   - Given a circular dependency with quick fix
   - When I want the step-by-step guide
   - Then I can access:
     - Direct reference to the guide in quickFix
     - Full guide content without additional API calls
   - And guide is embedded (not requiring separate lookup)

6. **AC6: Integrated Complexity and Impact**
   - Given a circular dependency
   - When viewing its summary
   - Then I see together:
     - Refactoring complexity score (1-10) from Story 3.5
     - Impact assessment (risk level, affected %) from Story 3.6
     - Estimated fix time from the recommended strategy
   - And these are accessible at CircularDependencyInfo level

7. **AC7: Complete CircularDependencyInfo Structure**
   - Given analysis is complete
   - When circular dependencies are detected
   - Then each CircularDependencyInfo includes ALL enrichments:
     - `rootCause`: RootCauseAnalysis (Story 3.1)
     - `importTraces`: ImportTrace[] (Story 3.2)
     - `fixStrategies`: FixStrategy[] with guides (Stories 3.3, 3.4)
     - `refactoringComplexity`: RefactoringComplexity (Story 3.5)
     - `impactAssessment`: ImpactAssessment (Story 3.6)
     - `quickFix`: QuickFixSummary (NEW)
     - `priorityScore`: number (NEW)
   - And strategies include `beforeAfterExplanation` (Story 3.7)

8. **AC8: Backward Compatibility**
   - Given existing consumers of AnalysisResult
   - When upgrading to include fix suggestions
   - Then:
     - All new fields are optional (`omitempty`)
     - Existing fields remain unchanged
     - JSON structure is additive only
   - And existing integrations continue working without changes

9. **AC9: Performance**
   - Given a workspace with 100 packages and 5 cycles
   - When generating complete analysis with all enrichments
   - Then total analysis completes in < 6 seconds (within NFR1)
   - And sorting/summary calculation adds < 50ms overhead
   - And memory usage remains < 100MB (within NFR4)

## Tasks / Subtasks

- [x] **Task 1: Define QuickFixSummary Type** (AC: #1, #5)
  - [x] 1.1 Create `pkg/types/quick_fix_summary.go`:
    ```go
    package types

    // QuickFixSummary provides quick access to the best fix recommendation.
    // This is a convenience wrapper around the best FixStrategy.
    type QuickFixSummary struct {
        // StrategyType is the type of the recommended fix
        StrategyType FixStrategyType `json:"strategyType"`

        // StrategyName is the human-readable name
        StrategyName string `json:"strategyName"`

        // Summary is a one-line description of what the fix accomplishes
        Summary string `json:"summary"`

        // Suitability is the strategy's suitability score (1-10)
        Suitability int `json:"suitability"`

        // Effort is the estimated effort level
        Effort EffortLevel `json:"effort"`

        // EstimatedTime is the time to implement (e.g., "15-30 minutes")
        EstimatedTime string `json:"estimatedTime"`

        // Guide is the full step-by-step guide (embedded for one-click access)
        Guide *FixGuide `json:"guide,omitempty"`

        // StrategyIndex is the index into fixStrategies[] for full details
        StrategyIndex int `json:"strategyIndex"`
    }
    ```
  - [x] 1.2 Add JSON serialization tests in `pkg/types/quick_fix_summary_test.go`
  - [x] 1.3 Ensure all JSON tags use camelCase

- [x] **Task 2: Define FixSummary Type for AnalysisResult** (AC: #4)
  - [x] 2.1 Create `pkg/types/fix_summary.go`:
    ```go
    package types

    // FixSummary provides aggregated statistics about fix recommendations.
    // Added to AnalysisResult for high-level overview.
    type FixSummary struct {
        // TotalCircularDependencies is the count of all detected cycles
        TotalCircularDependencies int `json:"totalCircularDependencies"`

        // TotalEstimatedFixTime is the sum of all fix times (human-readable)
        TotalEstimatedFixTime string `json:"totalEstimatedFixTime"`

        // QuickWinsCount is the number of low-complexity (1-3) fixes
        QuickWinsCount int `json:"quickWinsCount"`

        // CriticalCyclesCount is the number of critical impact cycles
        CriticalCyclesCount int `json:"criticalCyclesCount"`

        // HighPriorityCycles lists the top 3 cycles to fix first (by priority score)
        HighPriorityCycles []PriorityCycleSummary `json:"highPriorityCycles"`
    }

    // PriorityCycleSummary provides a brief overview of a prioritized cycle.
    type PriorityCycleSummary struct {
        // CycleID is a unique identifier for the cycle (e.g., first two packages)
        CycleID string `json:"cycleId"`

        // PackagesInvolved lists the packages in the cycle
        PackagesInvolved []string `json:"packagesInvolved"`

        // PriorityScore is the calculated priority (impact × ease)
        PriorityScore float64 `json:"priorityScore"`

        // RecommendedFix is the quick fix strategy type
        RecommendedFix FixStrategyType `json:"recommendedFix"`

        // EstimatedTime is the fix time
        EstimatedTime string `json:"estimatedTime"`
    }
    ```
  - [x] 2.2 Add JSON serialization tests
  - [x] 2.3 Ensure all JSON tags use camelCase

- [x] **Task 3: Update CircularDependencyInfo** (AC: #1, #3, #6, #7)
  - [x] 3.1 Update `pkg/types/circular.go`:
    ```go
    type CircularDependencyInfo struct {
        // Existing fields
        Cycle     []string         `json:"cycle"`
        Type      CircularType     `json:"type"`
        Severity  CircularSeverity `json:"severity"`
        Depth     int              `json:"depth"`
        Impact    string           `json:"impact"`
        Complexity int             `json:"complexity"` // Legacy field

        // Story 3.1
        RootCause *RootCauseAnalysis `json:"rootCause,omitempty"`

        // Story 3.2
        ImportTraces []ImportTrace `json:"importTraces,omitempty"`

        // Story 3.3, 3.4, 3.7 (enriched)
        FixStrategies []FixStrategy `json:"fixStrategies,omitempty"`

        // Story 3.5
        RefactoringComplexity *RefactoringComplexity `json:"refactoringComplexity,omitempty"`

        // Story 3.6
        ImpactAssessment *ImpactAssessment `json:"impactAssessment,omitempty"`

        // NEW Story 3.8: Quick access to best fix
        QuickFix *QuickFixSummary `json:"quickFix,omitempty"`

        // NEW Story 3.8: Priority for sorting (higher = fix first)
        PriorityScore float64 `json:"priorityScore"`
    }
    ```
  - [x] 3.2 Verify existing tests still pass

- [x] **Task 4: Update AnalysisResult** (AC: #4)
  - [x] 4.1 Update `pkg/types/analysis.go`:
    ```go
    type AnalysisResult struct {
        // Existing fields
        Workspace           *WorkspaceData           `json:"workspace"`
        DependencyGraph     *DependencyGraph         `json:"dependencyGraph"`
        CircularDependencies []CircularDependencyInfo `json:"circularDependencies"`
        DuplicateDependencies []DuplicateDependency   `json:"duplicateDependencies"`
        HealthScore         *HealthScore             `json:"healthScore"`
        Errors              []AnalysisError          `json:"errors"`

        // NEW Story 3.8: Aggregated fix summary
        FixSummary *FixSummary `json:"fixSummary,omitempty"`
    }
    ```
  - [x] 4.2 Verify existing tests still pass

- [x] **Task 5: Create ResultEnricher** (AC: #1, #2, #3, #4, #7)
  - [x] 5.1 Create `pkg/analyzer/result_enricher.go`:
    ```go
    package analyzer

    import (
        "sort"
        "github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
    )

    // ResultEnricher integrates all fix suggestions into analysis results.
    type ResultEnricher struct {
        graph     *types.DependencyGraph
        workspace *types.WorkspaceData
    }

    // NewResultEnricher creates a new enricher.
    func NewResultEnricher(graph *types.DependencyGraph, workspace *types.WorkspaceData) *ResultEnricher

    // Enrich adds all fix-related data to the analysis result.
    func (re *ResultEnricher) Enrich(result *types.AnalysisResult) *types.AnalysisResult

    // enrichCircularDependency adds quick fix and priority to a single cycle.
    func (re *ResultEnricher) enrichCircularDependency(cycle *types.CircularDependencyInfo)

    // sortStrategies sorts fix strategies by suitability (highest first).
    func sortStrategies(strategies []types.FixStrategy)

    // calculatePriorityScore computes impact × ease for sorting.
    func calculatePriorityScore(cycle *types.CircularDependencyInfo) float64

    // sortCircularDependencies sorts cycles by priority score (highest first).
    func sortCircularDependencies(cycles []types.CircularDependencyInfo)

    // generateFixSummary creates aggregated fix statistics.
    func (re *ResultEnricher) generateFixSummary(cycles []types.CircularDependencyInfo) *types.FixSummary

    // createQuickFix extracts the best strategy as QuickFixSummary.
    func createQuickFix(strategies []types.FixStrategy) *types.QuickFixSummary
    ```
  - [x] 5.2 Implement enrichment logic
  - [x] 5.3 Create comprehensive tests in `pkg/analyzer/result_enricher_test.go`

- [x] **Task 6: Implement Strategy Sorting** (AC: #2)
  - [x] 6.1 Implement `sortStrategies`:
    ```go
    func sortStrategies(strategies []types.FixStrategy) {
        sort.SliceStable(strategies, func(i, j int) bool {
            return strategies[i].Suitability > strategies[j].Suitability
        })
    }
    ```
  - [x] 6.2 Add tests for strategy sorting

- [x] **Task 7: Implement Priority Score Calculation** (AC: #3)
  - [x] 7.1 Implement `calculatePriorityScore`:
    ```go
    func calculatePriorityScore(cycle *types.CircularDependencyInfo) float64 {
        // Impact factor: use ImpactAssessment if available
        var impactFactor float64 = 5.0 // Default medium impact
        if cycle.ImpactAssessment != nil {
            switch cycle.ImpactAssessment.RiskLevel {
            case types.RiskLevelCritical:
                impactFactor = 10.0
            case types.RiskLevelHigh:
                impactFactor = 7.5
            case types.RiskLevelMedium:
                impactFactor = 5.0
            case types.RiskLevelLow:
                impactFactor = 2.5
            }
        }

        // Ease factor: inverse of complexity (11 - complexity)
        // Higher ease = lower complexity = easier to fix
        var easeFactor float64 = 6.0 // Default medium ease
        if cycle.RefactoringComplexity != nil {
            easeFactor = float64(11 - cycle.RefactoringComplexity.Score)
        }

        // Priority = Impact × Ease (quick wins = high impact, low complexity)
        return impactFactor * easeFactor
    }
    ```
  - [x] 7.2 Add tests for priority score calculation

- [x] **Task 8: Implement Circular Dependency Sorting** (AC: #3)
  - [x] 8.1 Implement `sortCircularDependencies`:
    ```go
    func sortCircularDependencies(cycles []types.CircularDependencyInfo) {
        sort.SliceStable(cycles, func(i, j int) bool {
            return cycles[i].PriorityScore > cycles[j].PriorityScore
        })
    }
    ```
  - [x] 8.2 Add tests for cycle sorting

- [x] **Task 9: Implement Quick Fix Creation** (AC: #1, #5)
  - [x] 9.1 Implement `createQuickFix`:
    ```go
    func createQuickFix(strategies []types.FixStrategy) *types.QuickFixSummary {
        if len(strategies) == 0 {
            return nil
        }

        // First strategy is best (already sorted by suitability)
        best := strategies[0]

        // Generate one-line summary
        var summary string
        switch best.Type {
        case types.FixStrategyExtractModule:
            summary = fmt.Sprintf("Create new shared package '%s' to break the cycle", best.NewPackageName)
        case types.FixStrategyDependencyInject:
            summary = "Invert dependency using dependency injection pattern"
        case types.FixStrategyBoundaryRefactor:
            summary = "Restructure package boundaries to eliminate overlap"
        default:
            summary = best.Description
        }

        // Get estimated time from guide or complexity
        estimatedTime := "15-30 minutes" // Default
        if best.Guide != nil && best.Guide.EstimatedTime != "" {
            estimatedTime = best.Guide.EstimatedTime
        }

        return &types.QuickFixSummary{
            StrategyType:  best.Type,
            StrategyName:  best.Name,
            Summary:       summary,
            Suitability:   best.Suitability,
            Effort:        best.Effort,
            EstimatedTime: estimatedTime,
            Guide:         best.Guide,
            StrategyIndex: 0,
        }
    }
    ```
  - [x] 9.2 Add tests for quick fix creation

- [x] **Task 10: Implement Fix Summary Generation** (AC: #4)
  - [x] 10.1 Implement `generateFixSummary`:
    ```go
    func (re *ResultEnricher) generateFixSummary(cycles []types.CircularDependencyInfo) *types.FixSummary {
        if len(cycles) == 0 {
            return nil
        }

        totalMinutes := 0
        quickWinsCount := 0
        criticalCount := 0
        highPriorityCycles := []types.PriorityCycleSummary{}

        for i, cycle := range cycles {
            // Count quick wins (complexity <= 3)
            if cycle.RefactoringComplexity != nil && cycle.RefactoringComplexity.Score <= 3 {
                quickWinsCount++
            }

            // Count critical cycles
            if cycle.ImpactAssessment != nil && cycle.ImpactAssessment.RiskLevel == types.RiskLevelCritical {
                criticalCount++
            }

            // Sum estimated times (parse from string)
            if cycle.QuickFix != nil {
                totalMinutes += parseEstimatedMinutes(cycle.QuickFix.EstimatedTime)
            }

            // Collect top 3 high priority
            if i < 3 {
                cycleID := generateCycleID(cycle.Cycle)
                var recommendedFix types.FixStrategyType
                var estTime string
                if cycle.QuickFix != nil {
                    recommendedFix = cycle.QuickFix.StrategyType
                    estTime = cycle.QuickFix.EstimatedTime
                }

                highPriorityCycles = append(highPriorityCycles, types.PriorityCycleSummary{
                    CycleID:          cycleID,
                    PackagesInvolved: getUniquePackages(cycle.Cycle),
                    PriorityScore:    cycle.PriorityScore,
                    RecommendedFix:   recommendedFix,
                    EstimatedTime:    estTime,
                })
            }
        }

        return &types.FixSummary{
            TotalCircularDependencies: len(cycles),
            TotalEstimatedFixTime:     formatTotalTime(totalMinutes),
            QuickWinsCount:            quickWinsCount,
            CriticalCyclesCount:       criticalCount,
            HighPriorityCycles:        highPriorityCycles,
        }
    }
    ```
  - [x] 10.2 Implement helper functions:
    ```go
    func parseEstimatedMinutes(timeStr string) int {
        // Parse strings like "15-30 minutes", "1-2 hours"
        // Return midpoint in minutes
        if strings.Contains(timeStr, "hour") {
            // Extract hours, convert to minutes
            // "1-2 hours" -> 90 minutes (midpoint)
            return 90 // Default for hour-range estimates
        }
        // "15-30 minutes" -> 22 minutes (midpoint)
        // "5-15 minutes" -> 10 minutes
        if strings.Contains(timeStr, "5-15") {
            return 10
        }
        if strings.Contains(timeStr, "15-30") {
            return 22
        }
        if strings.Contains(timeStr, "30-60") {
            return 45
        }
        return 30 // Default
    }

    func formatTotalTime(minutes int) string {
        if minutes < 60 {
            return fmt.Sprintf("%d minutes", minutes)
        }
        hours := minutes / 60
        remainingMinutes := minutes % 60
        if remainingMinutes == 0 {
            return fmt.Sprintf("%d hours", hours)
        }
        return fmt.Sprintf("%d hours %d minutes", hours, remainingMinutes)
    }

    func generateCycleID(cycle []string) string {
        if len(cycle) < 2 {
            return "unknown"
        }
        return fmt.Sprintf("%s→%s", extractShortName(cycle[0]), extractShortName(cycle[1]))
    }

    func getUniquePackages(cycle []string) []string {
        if len(cycle) == 0 {
            return []string{}
        }
        // Exclude last element (duplicate of first)
        return cycle[:len(cycle)-1]
    }
    ```
  - [x] 10.3 Add tests for fix summary generation

- [x] **Task 11: Implement Main Enrich Function** (AC: #7)
  - [x] 11.1 Implement `Enrich`:
    ```go
    func (re *ResultEnricher) Enrich(result *types.AnalysisResult) *types.AnalysisResult {
        if result == nil || len(result.CircularDependencies) == 0 {
            return result
        }

        // Step 1: Enrich each circular dependency
        for i := range result.CircularDependencies {
            re.enrichCircularDependency(&result.CircularDependencies[i])
        }

        // Step 2: Sort circular dependencies by priority
        sortCircularDependencies(result.CircularDependencies)

        // Step 3: Generate aggregated fix summary
        result.FixSummary = re.generateFixSummary(result.CircularDependencies)

        return result
    }

    func (re *ResultEnricher) enrichCircularDependency(cycle *types.CircularDependencyInfo) {
        // Sort strategies by suitability (if present)
        if len(cycle.FixStrategies) > 0 {
            sortStrategies(cycle.FixStrategies)
            cycle.QuickFix = createQuickFix(cycle.FixStrategies)
        }

        // Calculate priority score
        cycle.PriorityScore = calculatePriorityScore(cycle)
    }
    ```
  - [x] 11.2 Add comprehensive integration tests

- [x] **Task 12: Wire to Analyzer Pipeline** (AC: #7, #8)
  - [x] 12.1 Update `pkg/analyzer/analyzer.go`:
    ```go
    func (a *Analyzer) AnalyzeWithSources(...) (*types.AnalysisResult, error) {
        // ... existing analysis (workspace, graph, cycles, duplicates, health) ...

        // Generate all enrichments (Stories 3.1-3.7)
        rootCauseAnalyzer := NewRootCauseAnalyzer(graph)
        importTracer := NewImportTracer(graph, workspace)
        fixGenerator := NewFixStrategyGenerator(graph, workspace)
        guideGenerator := NewFixGuideGenerator(workspace)
        complexityCalc := NewComplexityCalculator(graph, workspace)
        impactAnalyzer := NewImpactAnalyzer(graph, workspace)
        beforeAfterGen := NewBeforeAfterGenerator(graph, workspace)

        for _, cycle := range result.CircularDependencies {
            // Story 3.1: Root cause
            cycle.RootCause = rootCauseAnalyzer.Analyze(cycle)

            // Story 3.2: Import traces
            cycle.ImportTraces = importTracer.Trace(cycle)

            // Story 3.3: Fix strategies
            cycle.FixStrategies = fixGenerator.Generate(cycle)

            // Story 3.4: Fix guides + Story 3.7: Before/after
            for i := range cycle.FixStrategies {
                cycle.FixStrategies[i].Guide = guideGenerator.Generate(cycle, &cycle.FixStrategies[i])
                cycle.FixStrategies[i].BeforeAfterExplanation = beforeAfterGen.Generate(cycle, &cycle.FixStrategies[i])
            }

            // Story 3.5: Complexity
            cycle.RefactoringComplexity = complexityCalc.Calculate(cycle)

            // Story 3.6: Impact
            cycle.ImpactAssessment = impactAnalyzer.Analyze(cycle)
        }

        // NEW Story 3.8: Integrate, sort, and summarize
        enricher := NewResultEnricher(graph, workspace)
        result = enricher.Enrich(result)

        return result, nil
    }
    ```
  - [x] 12.2 Update analyzer tests

- [x] **Task 13: Update TypeScript Types** (AC: #7, #8)
  - [x] 13.1 Update `packages/types/src/analysis/results.ts`:
    ```typescript
    export interface QuickFixSummary {
      /** Type of the recommended fix */
      strategyType: FixStrategyType;
      /** Human-readable strategy name */
      strategyName: string;
      /** One-line description */
      summary: string;
      /** Suitability score (1-10) */
      suitability: number;
      /** Estimated effort level */
      effort: EffortLevel;
      /** Time to implement */
      estimatedTime: string;
      /** Full step-by-step guide */
      guide?: FixGuide;
      /** Index into fixStrategies[] for full details */
      strategyIndex: number;
    }

    export interface FixSummary {
      /** Count of all detected cycles */
      totalCircularDependencies: number;
      /** Sum of all fix times */
      totalEstimatedFixTime: string;
      /** Number of low-complexity fixes */
      quickWinsCount: number;
      /** Number of critical impact cycles */
      criticalCyclesCount: number;
      /** Top 3 cycles to fix first */
      highPriorityCycles: PriorityCycleSummary[];
    }

    export interface PriorityCycleSummary {
      /** Unique identifier for the cycle */
      cycleId: string;
      /** Packages in the cycle */
      packagesInvolved: string[];
      /** Priority score (impact × ease) */
      priorityScore: number;
      /** Recommended fix strategy */
      recommendedFix: FixStrategyType;
      /** Estimated fix time */
      estimatedTime: string;
    }

    export interface CircularDependencyInfo {
      // ... existing fields ...
      /** Quick access to best fix (Story 3.8) */
      quickFix?: QuickFixSummary;
      /** Priority for sorting - higher = fix first (Story 3.8) */
      priorityScore: number;
    }

    export interface AnalysisResult {
      // ... existing fields ...
      /** Aggregated fix summary (Story 3.8) */
      fixSummary?: FixSummary;
    }
    ```
  - [x] 13.2 Run `pnpm nx build types` to verify
  - [x] 13.3 Add type tests for new interfaces

- [x] **Task 14: Performance Testing** (AC: #9)
  - [x] 14.1 Create `pkg/analyzer/result_enricher_benchmark_test.go`:
    ```go
    func BenchmarkResultEnrichment(b *testing.B) {
        graph := generateGraph(100)
        workspace := generateWorkspace(100)
        result := generateAnalysisResultWithCycles(5)
        enricher := NewResultEnricher(graph, workspace)

        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            enricher.Enrich(result)
        }
    }

    func BenchmarkCompleteAnalysis(b *testing.B) {
        // Test complete analysis pipeline with all enrichments
        workspace := generateWorkspace(100)
        analyzer := NewAnalyzer()

        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            analyzer.AnalyzeWorkspace(workspace)
        }
    }
    ```
  - [x] 14.2 Verify enrichment adds < 50ms overhead
  - [x] 14.3 Verify complete analysis < 6 seconds for 100 packages, 5 cycles
  - [x] 14.4 Document actual performance in completion notes

- [x] **Task 15: Integration Verification** (AC: all)
  - [x] 15.1 Run all tests: `cd packages/analysis-engine && make test`
  - [x] 15.2 Build WASM: `pnpm nx build @monoguard/analysis-engine`
  - [x] 15.3 Run affected CI checks: `pnpm nx affected --target=lint,test,type-check --base=main`
  - [x] 15.4 Verify JSON output includes:
    - CircularDependencyInfo.quickFix
    - CircularDependencyInfo.priorityScore
    - AnalysisResult.fixSummary
  - [x] 15.5 Verify circular dependencies are sorted by priority
  - [x] 15.6 Verify strategies are sorted by suitability

## Dev Notes

### Architecture Patterns & Constraints

**From Architecture Document:**
- **Location:** Result enricher in `pkg/analyzer/`
- **Pattern:** Enricher pattern that transforms AnalysisResult
- **Integration:** Final step in analysis pipeline, after all story generators
- **Dependency:** Depends on outputs from Stories 3.1-3.7

**Critical Constraints:**
- **camelCase JSON:** All struct tags MUST use camelCase
- **Optional Fields:** All new fields are `omitempty` - backward compatible
- **Deterministic:** Sorting must be stable and reproducible
- **Performance:** Enrichment overhead < 50ms

### Critical Don't-Miss Rules

**From project-context.md:**

1. **JSON Naming Convention:**
   ```go
   // ✅ CORRECT: camelCase JSON tags
   type QuickFixSummary struct {
       StrategyType  FixStrategyType `json:"strategyType"`
       StrategyName  string          `json:"strategyName"`
       StrategyIndex int             `json:"strategyIndex"`
   }

   // ❌ WRONG: snake_case JSON tags
   type QuickFixSummary struct {
       StrategyType FixStrategyType `json:"strategy_type"` // WRONG!
   }
   ```

2. **Stable Sorting:**
   ```go
   // ✅ CORRECT: Use SliceStable for deterministic ordering
   sort.SliceStable(strategies, func(i, j int) bool {
       return strategies[i].Suitability > strategies[j].Suitability
   })

   // ❌ WRONG: Unstable sort leads to non-deterministic results
   sort.Slice(strategies, func(i, j int) bool { ... })
   ```

3. **Slice Initialization:**
   ```go
   // ✅ CORRECT: Initialize as empty slice
   highPriorityCycles := []types.PriorityCycleSummary{}

   // ❌ WRONG: Nil slice (serializes as null)
   var highPriorityCycles []types.PriorityCycleSummary // nil
   ```

4. **Pointer vs Value for Optional Structs:**
   ```go
   // ✅ CORRECT: Use pointer for optional nested structs
   type AnalysisResult struct {
       FixSummary *FixSummary `json:"fixSummary,omitempty"`
   }

   // ❌ WRONG: Value type with omitempty doesn't work for structs
   type AnalysisResult struct {
       FixSummary FixSummary `json:"fixSummary,omitempty"` // Won't omit!
   }
   ```

### Project Structure Notes

**Target Directory Structure:**
```
packages/analysis-engine/
├── pkg/
│   ├── analyzer/
│   │   ├── analyzer.go                        # UPDATE: Wire enricher
│   │   ├── result_enricher.go                 # NEW: Result enricher
│   │   ├── result_enricher_test.go            # NEW: Enricher tests
│   │   └── result_enricher_benchmark_test.go  # NEW: Performance
│   └── types/
│       ├── analysis.go                        # UPDATE: Add FixSummary
│       ├── circular.go                        # UPDATE: Add QuickFix, PriorityScore
│       ├── quick_fix_summary.go               # NEW: QuickFixSummary type
│       ├── quick_fix_summary_test.go          # NEW: Type tests
│       ├── fix_summary.go                     # NEW: FixSummary type
│       └── fix_summary_test.go                # NEW: Type tests
└── ...

packages/types/src/analysis/
├── results.ts                                 # UPDATE: Add TS types
└── ...
```

### Previous Story Intelligence

**From Story 3.7 (ready-for-dev):**
- BeforeAfterExplanation enriches FixStrategy - now integrated
- Confidence score pattern - QuickFix uses similar summary pattern
- Warnings pattern - FixSummary.criticalCyclesCount serves similar purpose

**From Story 3.6 (ready-for-dev):**
- ImpactAssessment has RiskLevel - **Key Input:** Used for priority calculation
- RiskLevel: critical/high/medium/low - map to impact factors 10/7.5/5/2.5

**From Story 3.5 (ready-for-dev):**
- RefactoringComplexity has Score (1-10) - **Key Input:** Used for priority calculation
- EstimatedTime string format - reuse for QuickFix and FixSummary
- Ease factor = 11 - Complexity (invert for priority)

**From Story 3.4 (done):**
- FixGuide is attached to FixStrategy - embedded in QuickFix for one-click access
- EstimatedTime format patterns - reuse parsing logic

**From Story 3.3 (done):**
- FixStrategy has Suitability score - **Key Input:** Used for sorting and quick fix selection
- Strategies generated per cycle - sorting happens after generation

### Priority Score Algorithm

**Formula:** `Priority = Impact × Ease`

| Impact Factor | RiskLevel | Value |
|---------------|-----------|-------|
| Critical | > 50% affected OR core package | 10.0 |
| High | 25-50% affected | 7.5 |
| Medium | 10-25% affected | 5.0 |
| Low | < 10% affected | 2.5 |
| Default (no assessment) | - | 5.0 |

| Ease Factor | Complexity Score | Value |
|-------------|------------------|-------|
| Very Easy | 1-2 | 9-10 |
| Easy | 3-4 | 7-8 |
| Medium | 5-6 | 5-6 |
| Hard | 7-8 | 3-4 |
| Very Hard | 9-10 | 1-2 |
| Default (no complexity) | - | 6.0 |

**Examples:**
| Cycle | Impact | Complexity | Priority | Rank |
|-------|--------|------------|----------|------|
| A → B (critical, easy) | 10.0 | 9.0 (score 2) | 90.0 | 1st (quick win!) |
| C → D (high, medium) | 7.5 | 6.0 (score 5) | 45.0 | 2nd |
| E → F (low, hard) | 2.5 | 3.0 (score 8) | 7.5 | 3rd |

### Input/Output Format

**Input (AnalysisResult with enriched cycles):**
```json
{
  "circularDependencies": [
    {
      "cycle": ["@mono/core", "@mono/ui", "@mono/core"],
      "refactoringComplexity": {"score": 3},
      "impactAssessment": {"riskLevel": "critical"},
      "fixStrategies": [
        {"type": "extract-module", "suitability": 8, "guide": {...}},
        {"type": "dependency-injection", "suitability": 6, "guide": {...}}
      ]
    }
  ]
}
```

**Output (Enriched AnalysisResult):**
```json
{
  "circularDependencies": [
    {
      "cycle": ["@mono/core", "@mono/ui", "@mono/core"],
      "refactoringComplexity": {"score": 3},
      "impactAssessment": {"riskLevel": "critical"},
      "fixStrategies": [
        {"type": "extract-module", "suitability": 8, ...},
        {"type": "dependency-injection", "suitability": 6, ...}
      ],
      "quickFix": {
        "strategyType": "extract-module",
        "strategyName": "Extract Shared Module",
        "summary": "Create new shared package '@mono/shared' to break the cycle",
        "suitability": 8,
        "effort": "medium",
        "estimatedTime": "30-60 minutes",
        "guide": {...},
        "strategyIndex": 0
      },
      "priorityScore": 80.0
    }
  ],
  "fixSummary": {
    "totalCircularDependencies": 1,
    "totalEstimatedFixTime": "45 minutes",
    "quickWinsCount": 1,
    "criticalCyclesCount": 1,
    "highPriorityCycles": [
      {
        "cycleId": "core→ui",
        "packagesInvolved": ["@mono/core", "@mono/ui"],
        "priorityScore": 80.0,
        "recommendedFix": "extract-module",
        "estimatedTime": "30-60 minutes"
      }
    ]
  }
}
```

### Test Scenarios

| Scenario | Cycles | Expected QuickWins | Expected Critical | Top Priority |
|----------|--------|-------------------|-------------------|--------------|
| All easy | 3 (all complexity ≤ 3) | 3 | 0 | Highest impact first |
| All hard | 3 (all complexity > 7) | 0 | varies | Highest impact first |
| Mixed | 5 (2 easy, 3 hard) | 2 | varies | Easy + high impact first |
| No cycles | 0 | 0 | 0 | N/A (no FixSummary) |
| One critical | 1 (critical + easy) | 1 | 1 | That one cycle |

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story 3.8]
- [Source: _bmad-output/planning-artifacts/prd.md#FR7-FR14]
- [Source: _bmad-output/project-context.md#Language-Specific Rules]
- [Source: _bmad-output/implementation-artifacts/3-5-calculate-refactoring-complexity-scores.md]
- [Source: _bmad-output/implementation-artifacts/3-6-generate-impact-assessment.md]
- [Source: _bmad-output/implementation-artifacts/3-7-provide-before-after-fix-explanations.md]
- [Sorting Algorithms - Go Documentation](https://pkg.go.dev/sort)

## Dev Agent Record

### Agent Model Used

Claude Opus 4.5 (claude-opus-4-5-20251101)

### Debug Log References

N/A

### Completion Notes List

- **Performance verified**: Enrichment overhead is ~0.044ms (44μs), well under 50ms limit (AC9)
- **All tests pass**: Go tests (pkg/analyzer, pkg/types) and TypeScript tests (@monoguard/types)
- **CI checks pass**: lint, test, type-check all pass for affected projects
- **Types defined**: QuickFixSummary, FixSummary, PriorityCycleSummary in Go and TypeScript
- **Enricher integrated**: ResultEnricher wired into Analyze and AnalyzeWithSources pipelines
- **Sorting implemented**:
  - Strategies sorted by suitability (highest first) using stable sort
  - Circular dependencies sorted by priority score (highest first) using stable sort
- **Priority algorithm**: Priority = ImpactFactor × EaseFactor where:
  - ImpactFactor: critical=10, high=7.5, medium=5, low=2.5
  - EaseFactor: 11 - complexity score (min 1)
- **Quick fix**: Extracts best strategy with embedded guide for one-click access
- **Fix summary**: Aggregates total cycles, estimated time, quick wins count, critical count, top 3 priorities

### File List

**New Files Created:**
- `packages/analysis-engine/pkg/types/quick_fix_summary.go` - QuickFixSummary type
- `packages/analysis-engine/pkg/types/quick_fix_summary_test.go` - Tests
- `packages/analysis-engine/pkg/types/fix_summary.go` - FixSummary and PriorityCycleSummary types
- `packages/analysis-engine/pkg/types/fix_summary_test.go` - Tests
- `packages/analysis-engine/pkg/analyzer/result_enricher.go` - ResultEnricher implementation
- `packages/analysis-engine/pkg/analyzer/result_enricher_test.go` - Unit tests
- `packages/analysis-engine/pkg/analyzer/result_enricher_benchmark_test.go` - Performance benchmarks

**Modified Files:**
- `packages/analysis-engine/pkg/types/circular.go` - Added QuickFix and PriorityScore fields
- `packages/analysis-engine/pkg/types/circular_test.go` - Updated tests
- `packages/analysis-engine/pkg/types/types.go` - Added FixSummary to AnalysisResult
- `packages/analysis-engine/pkg/types/types_test.go` - Updated tests
- `packages/analysis-engine/pkg/analyzer/analyzer.go` - Wired ResultEnricher to pipeline
- `packages/types/src/analysis/results.ts` - Added TypeScript types
