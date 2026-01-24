// Package analyzer contains benchmarks for before/after explanation generator.
package analyzer

import (
	"fmt"
	"testing"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// generateBeforeAfterTestGraph creates a test graph with n packages.
func generateBeforeAfterTestGraph(n int) *types.DependencyGraph {
	graph := types.NewDependencyGraph("/root", types.WorkspaceTypePnpm)

	for i := 0; i < n; i++ {
		pkgName := fmt.Sprintf("@mono/pkg%d", i)
		node := types.NewPackageNode(pkgName, "1.0.0", fmt.Sprintf("packages/pkg%d", i))

		// Add some dependencies to make it realistic
		if i > 0 {
			depIdx := i - 1
			node.Dependencies = append(node.Dependencies, fmt.Sprintf("@mono/pkg%d", depIdx))
		}

		graph.Nodes[pkgName] = node
	}

	return graph
}

// generateBeforeAfterTestWorkspace creates a test workspace with n packages.
func generateBeforeAfterTestWorkspace(n int) *types.WorkspaceData {
	workspace := &types.WorkspaceData{
		RootPath:      "/root",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages:      make(map[string]*types.PackageInfo),
	}

	for i := 0; i < n; i++ {
		pkgName := fmt.Sprintf("@mono/pkg%d", i)
		workspace.Packages[pkgName] = &types.PackageInfo{
			Name:    pkgName,
			Version: "1.0.0",
			Path:    fmt.Sprintf("packages/pkg%d", i),
		}
	}

	return workspace
}

// generateBeforeAfterTestCycles creates test cycles with fix strategies.
func generateBeforeAfterTestCycles(numCycles int, packagesPerCycle int) []*types.CircularDependencyInfo {
	cycles := make([]*types.CircularDependencyInfo, numCycles)

	for i := 0; i < numCycles; i++ {
		// Create a cycle of packagesPerCycle packages
		cyclePath := make([]string, packagesPerCycle+1)
		for j := 0; j < packagesPerCycle; j++ {
			cyclePath[j] = fmt.Sprintf("@mono/cycle%d-pkg%d", i, j)
		}
		cyclePath[packagesPerCycle] = cyclePath[0] // Close the cycle

		cycle := &types.CircularDependencyInfo{
			Cycle:    cyclePath,
			Type:     types.CircularTypeIndirect,
			Severity: types.CircularSeverityWarning,
			Depth:    packagesPerCycle,
			Impact:   fmt.Sprintf("Cycle %d affects %d packages", i, packagesPerCycle),
			ImportTraces: []types.ImportTrace{
				{
					FromPackage: cyclePath[0],
					ToPackage:   cyclePath[1],
					FilePath:    fmt.Sprintf("packages/cycle%d-pkg0/src/index.ts", i),
					LineNumber:  5,
					Statement:   fmt.Sprintf("import { helper } from '%s';", cyclePath[1]),
					ImportType:  types.ImportTypeESMNamed,
					Symbols:     []string{"helper"},
				},
			},
			FixStrategies: []types.FixStrategy{
				{
					Type:           types.FixStrategyExtractModule,
					Name:           "Extract Shared Module",
					Description:    "Create a new shared package",
					Suitability:    9,
					Effort:         types.EffortMedium,
					Pros:           []string{"Clean separation"},
					Cons:           []string{"New package to maintain"},
					Recommended:    true,
					TargetPackages: cyclePath[:packagesPerCycle],
					NewPackageName: fmt.Sprintf("@mono/cycle%d-shared", i),
				},
				{
					Type:           types.FixStrategyDependencyInject,
					Name:           "Dependency Injection",
					Description:    "Invert dependencies",
					Suitability:    7,
					Effort:         types.EffortMedium,
					Pros:           []string{"No new packages"},
					Cons:           []string{"More complex code"},
					Recommended:    false,
					TargetPackages: cyclePath[:packagesPerCycle],
				},
				{
					Type:           types.FixStrategyBoundaryRefactor,
					Name:           "Boundary Refactoring",
					Description:    "Restructure boundaries",
					Suitability:    5,
					Effort:         types.EffortHigh,
					Pros:           []string{"Better architecture"},
					Cons:           []string{"Significant effort"},
					Recommended:    false,
					TargetPackages: cyclePath[:packagesPerCycle],
				},
			},
		}

		cycles[i] = cycle
	}

	return cycles
}

// BenchmarkBeforeAfterGeneration benchmarks the full generation.
// AC#8: 100 packages, 5 cycles, should complete in <300ms
func BenchmarkBeforeAfterGeneration(b *testing.B) {
	graph := generateBeforeAfterTestGraph(100)
	workspace := generateBeforeAfterTestWorkspace(100)
	cycles := generateBeforeAfterTestCycles(5, 3)
	generator := NewBeforeAfterGenerator(graph, workspace)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, cycle := range cycles {
			for j := range cycle.FixStrategies {
				generator.Generate(cycle, &cycle.FixStrategies[j])
			}
		}
	}
}

// BenchmarkBeforeAfterGeneration_LargeCycles benchmarks with larger cycles.
func BenchmarkBeforeAfterGeneration_LargeCycles(b *testing.B) {
	graph := generateBeforeAfterTestGraph(100)
	workspace := generateBeforeAfterTestWorkspace(100)
	cycles := generateBeforeAfterTestCycles(5, 10) // 10 packages per cycle
	generator := NewBeforeAfterGenerator(graph, workspace)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, cycle := range cycles {
			for j := range cycle.FixStrategies {
				generator.Generate(cycle, &cycle.FixStrategies[j])
			}
		}
	}
}

// BenchmarkCurrentStateGeneration benchmarks just the current state generation.
func BenchmarkCurrentStateGeneration(b *testing.B) {
	generator := NewBeforeAfterGenerator(nil, nil)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/a", "@mono/b", "@mono/c", "@mono/d", "@mono/e", "@mono/a"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generator.generateCurrentState(cycle)
	}
}

// BenchmarkProposedStateGeneration benchmarks proposed state generation.
func BenchmarkProposedStateGeneration(b *testing.B) {
	graph := generateBeforeAfterTestGraph(100)
	workspace := generateBeforeAfterTestWorkspace(100)
	generator := NewBeforeAfterGenerator(graph, workspace)

	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/a", "@mono/b", "@mono/c", "@mono/d", "@mono/e", "@mono/a"},
	}
	strategy := &types.FixStrategy{
		Type:           types.FixStrategyExtractModule,
		TargetPackages: []string{"@mono/a", "@mono/b", "@mono/c", "@mono/d", "@mono/e"},
		NewPackageName: "@mono/shared",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generator.generateProposedState(cycle, strategy)
	}
}

// BenchmarkWarningsGeneration benchmarks warnings generation.
func BenchmarkWarningsGeneration(b *testing.B) {
	graph := generateBeforeAfterTestGraph(100)
	workspace := generateBeforeAfterTestWorkspace(100)
	generator := NewBeforeAfterGenerator(graph, workspace)

	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/core", "@mono/api", "@mono/ui", "@mono/core"},
		ImpactAssessment: &types.ImpactAssessment{
			RiskLevel:          types.RiskLevelCritical,
			AffectedPercentage: 0.75,
		},
	}
	strategy := &types.FixStrategy{
		Type:           types.FixStrategyExtractModule,
		TargetPackages: []string{"@mono/core", "@mono/api", "@mono/ui"},
		NewPackageName: "@mono/shared",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generator.generateWarnings(cycle, strategy)
	}
}

// TestBeforeAfterPerformance verifies the 300ms requirement.
func TestBeforeAfterPerformance(t *testing.T) {
	graph := generateBeforeAfterTestGraph(100)
	workspace := generateBeforeAfterTestWorkspace(100)
	cycles := generateBeforeAfterTestCycles(5, 3)
	generator := NewBeforeAfterGenerator(graph, workspace)

	// Run the benchmark multiple times and take average
	iterations := 100
	totalDuration := int64(0)

	for iter := 0; iter < iterations; iter++ {
		start := testing.Benchmark(func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, cycle := range cycles {
					for j := range cycle.FixStrategies {
						generator.Generate(cycle, &cycle.FixStrategies[j])
					}
				}
			}
		})
		totalDuration += start.NsPerOp()
	}

	avgNsPerOp := totalDuration / int64(iterations)
	avgMsPerOp := float64(avgNsPerOp) / 1_000_000

	t.Logf("Average time per operation: %.3f ms", avgMsPerOp)

	// AC#8: Should complete in <300ms
	// Note: This is for the full operation (5 cycles x 3 strategies = 15 generations)
	if avgMsPerOp > 300 {
		t.Errorf("Performance requirement not met: %.3f ms > 300ms", avgMsPerOp)
	}
}
