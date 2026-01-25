// Package analyzer provides dependency graph analysis for monorepo workspaces.
// This file contains benchmark tests for the result enricher (Story 3.8).
package analyzer

import (
	"fmt"
	"testing"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// generateEnricherTestGraph creates a dependency graph with n packages for enricher benchmarks.
func generateEnricherTestGraph(n int) *types.DependencyGraph {
	graph := types.NewDependencyGraph("/test/workspace", types.WorkspaceTypePnpm)

	for i := 0; i < n; i++ {
		name := fmt.Sprintf("@test/pkg-%d", i)
		node := types.NewPackageNode(name, "1.0.0", fmt.Sprintf("/test/packages/pkg-%d", i))

		// Add some internal dependencies
		if i > 0 {
			node.Dependencies = append(node.Dependencies, fmt.Sprintf("@test/pkg-%d", i-1))
		}

		graph.Nodes[name] = node
	}

	return graph
}

// generateEnricherTestWorkspace creates a workspace with n packages for enricher benchmarks.
func generateEnricherTestWorkspace(n int) *types.WorkspaceData {
	workspace := &types.WorkspaceData{
		RootPath:      "/test/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages:      make(map[string]*types.PackageInfo),
	}

	for i := 0; i < n; i++ {
		name := fmt.Sprintf("@test/pkg-%d", i)
		workspace.Packages[name] = &types.PackageInfo{
			Name:         name,
			Version:      "1.0.0",
			Path:         fmt.Sprintf("/test/packages/pkg-%d", i),
			Dependencies: make(map[string]string),
		}
	}

	return workspace
}

// generateTestCyclesWithEnrichments creates n CircularDependencyInfo with full enrichments.
func generateTestCyclesWithEnrichments(n int) []*types.CircularDependencyInfo {
	cycles := make([]*types.CircularDependencyInfo, n)

	for i := 0; i < n; i++ {
		cycles[i] = &types.CircularDependencyInfo{
			Cycle:    []string{fmt.Sprintf("@test/pkg-%d", i), fmt.Sprintf("@test/pkg-%d", i+1), fmt.Sprintf("@test/pkg-%d", i)},
			Type:     types.CircularTypeDirect,
			Severity: types.CircularSeverityWarning,
			Depth:    2,
			Impact:   "Test cycle",
			ImpactAssessment: &types.ImpactAssessment{
				RiskLevel:      types.RiskLevelMedium,
				TotalAffected:  5,
				RiskExplanation: "Test risk",
			},
			RefactoringComplexity: &types.RefactoringComplexity{
				Score:         5,
				EstimatedTime: "30-60 minutes",
			},
			FixStrategies: []types.FixStrategy{
				{
					Type:        types.FixStrategyExtractModule,
					Name:        "Extract Module",
					Suitability: 8,
					Effort:      types.EffortMedium,
					Guide: &types.FixGuide{
						EstimatedTime: "30-60 minutes",
					},
				},
				{
					Type:        types.FixStrategyDependencyInject,
					Name:        "Dependency Injection",
					Suitability: 6,
					Effort:      types.EffortHigh,
				},
				{
					Type:        types.FixStrategyBoundaryRefactor,
					Name:        "Boundary Refactoring",
					Suitability: 4,
					Effort:      types.EffortHigh,
				},
			},
		}
	}

	return cycles
}

// BenchmarkResultEnrichment benchmarks the Enrich function.
// AC9: Enrichment should add < 50ms overhead.
func BenchmarkResultEnrichment(b *testing.B) {
	graph := generateEnricherTestGraph(100)
	workspace := generateEnricherTestWorkspace(100)
	enricher := NewResultEnricher(graph, workspace)

	// Create result with 5 cycles (typical case)
	cycles := generateTestCyclesWithEnrichments(5)
	result := &types.AnalysisResult{
		HealthScore:          85,
		Packages:             100,
		CircularDependencies: cycles,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create a fresh copy to avoid mutation issues
		resultCopy := &types.AnalysisResult{
			HealthScore:          result.HealthScore,
			Packages:             result.Packages,
			CircularDependencies: generateTestCyclesWithEnrichments(5),
		}
		enricher.Enrich(resultCopy)
	}
}

// BenchmarkStrategySorting benchmarks strategy sorting.
func BenchmarkStrategySorting(b *testing.B) {
	strategies := []types.FixStrategy{
		{Name: "A", Suitability: 3},
		{Name: "B", Suitability: 9},
		{Name: "C", Suitability: 6},
		{Name: "D", Suitability: 7},
		{Name: "E", Suitability: 5},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create copy to avoid already-sorted slice
		strategiesCopy := make([]types.FixStrategy, len(strategies))
		copy(strategiesCopy, strategies)
		sortStrategies(strategiesCopy)
	}
}

// BenchmarkPriorityCalculation benchmarks priority score calculation.
func BenchmarkPriorityCalculation(b *testing.B) {
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"a", "b", "a"},
		ImpactAssessment: &types.ImpactAssessment{
			RiskLevel: types.RiskLevelHigh,
		},
		RefactoringComplexity: &types.RefactoringComplexity{
			Score: 5,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calculatePriorityScore(cycle)
	}
}

// BenchmarkCircularDependencySorting benchmarks cycle sorting.
func BenchmarkCircularDependencySorting(b *testing.B) {
	cycles := generateTestCyclesWithEnrichments(10)
	// Pre-calculate priority scores
	for _, c := range cycles {
		c.PriorityScore = calculatePriorityScore(c)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cyclesCopy := make([]*types.CircularDependencyInfo, len(cycles))
		copy(cyclesCopy, cycles)
		sortCircularDependencies(cyclesCopy)
	}
}

// BenchmarkFixSummaryGeneration benchmarks fix summary generation.
func BenchmarkFixSummaryGeneration(b *testing.B) {
	graph := generateEnricherTestGraph(100)
	workspace := generateEnricherTestWorkspace(100)
	enricher := NewResultEnricher(graph, workspace)
	cycles := generateTestCyclesWithEnrichments(5)

	// Pre-enrich cycles to have QuickFix
	for _, c := range cycles {
		sortStrategies(c.FixStrategies)
		c.QuickFix = createQuickFix(c.FixStrategies)
		c.PriorityScore = calculatePriorityScore(c)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		enricher.generateFixSummary(cycles)
	}
}

// BenchmarkCreateQuickFix benchmarks quick fix creation.
func BenchmarkCreateQuickFix(b *testing.B) {
	strategies := []types.FixStrategy{
		{
			Type:        types.FixStrategyExtractModule,
			Name:        "Extract Module",
			Suitability: 8,
			Effort:      types.EffortMedium,
			Guide: &types.FixGuide{
				EstimatedTime: "30-60 minutes",
			},
		},
		{
			Type:        types.FixStrategyDependencyInject,
			Name:        "DI",
			Suitability: 6,
			Effort:      types.EffortHigh,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		createQuickFix(strategies)
	}
}

// BenchmarkParseEstimatedMinutes benchmarks time parsing.
func BenchmarkParseEstimatedMinutes(b *testing.B) {
	timeStrs := []string{
		"15-30 minutes",
		"1-2 hours",
		"30 minutes",
		"",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, ts := range timeStrs {
			parseEstimatedMinutes(ts)
		}
	}
}

// TestEnrichmentOverhead verifies enrichment adds < 50ms overhead.
// This is not a benchmark but a verification test (AC9).
func TestEnrichmentOverhead(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping overhead test in short mode")
	}

	graph := generateEnricherTestGraph(100)
	workspace := generateEnricherTestWorkspace(100)
	enricher := NewResultEnricher(graph, workspace)

	// Run 100 enrichments and check average time
	iterations := 100
	totalDuration := int64(0)

	for i := 0; i < iterations; i++ {
		cycles := generateTestCyclesWithEnrichments(5)
		result := &types.AnalysisResult{
			HealthScore:          85,
			Packages:             100,
			CircularDependencies: cycles,
		}

		start := testing.Benchmark(func(b *testing.B) {
			enricher.Enrich(result)
		})

		totalDuration += start.T.Nanoseconds()
	}

	avgMs := float64(totalDuration) / float64(iterations) / 1e6
	t.Logf("Average enrichment time: %.2f ms", avgMs)

	// AC9: Enrichment should add < 50ms overhead
	if avgMs > 50 {
		t.Errorf("Enrichment overhead %.2f ms exceeds 50ms limit", avgMs)
	}
}
