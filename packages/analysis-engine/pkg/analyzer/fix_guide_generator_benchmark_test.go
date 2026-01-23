// Package analyzer provides dependency graph analysis for monorepo workspaces.
// This file contains benchmark tests for fix guide generator (Story 3.4).
package analyzer

import (
	"fmt"
	"testing"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// BenchmarkFixGuideGeneration benchmarks guide generation performance.
// AC8: Generation should complete in < 500ms for 100 packages with 5 cycles.
func BenchmarkFixGuideGeneration(b *testing.B) {
	workspace := generateBenchmarkWorkspace(100)
	cycles := generateBenchmarkCyclesWithStrategies(5, workspace)
	generator := NewFixGuideGenerator(workspace)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, cycle := range cycles {
			for j := range cycle.FixStrategies {
				generator.Generate(cycle, &cycle.FixStrategies[j])
			}
		}
	}
}

// BenchmarkFixGuideGeneration_ExtractModule benchmarks extract module guide.
func BenchmarkFixGuideGeneration_ExtractModule(b *testing.B) {
	workspace := generateBenchmarkWorkspace(50)
	generator := NewFixGuideGenerator(workspace)

	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/core", "@mono/api", "@mono/ui"},
		Depth: 3,
	}

	strategy := &types.FixStrategy{
		Type:           types.FixStrategyExtractModule,
		TargetPackages: []string{"@mono/ui", "@mono/core", "@mono/api"},
		NewPackageName: "@mono/shared",
		Effort:         types.EffortMedium,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generator.Generate(cycle, strategy)
	}
}

// BenchmarkFixGuideGeneration_DI benchmarks dependency injection guide.
func BenchmarkFixGuideGeneration_DI(b *testing.B) {
	workspace := generateBenchmarkWorkspace(50)
	generator := NewFixGuideGenerator(workspace)

	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/core", "@mono/ui"},
		Depth: 2,
		RootCause: &types.RootCauseAnalysis{
			CriticalEdge: &types.RootCauseEdge{
				From:     "@mono/ui",
				To:       "@mono/core",
				Type:     types.DependencyTypeProduction,
				Critical: true,
			},
		},
	}

	strategy := &types.FixStrategy{
		Type:           types.FixStrategyDependencyInject,
		TargetPackages: []string{"@mono/ui", "@mono/core"},
		Effort:         types.EffortLow,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generator.Generate(cycle, strategy)
	}
}

// BenchmarkFixGuideGeneration_BoundaryRefactor benchmarks boundary refactoring guide.
func BenchmarkFixGuideGeneration_BoundaryRefactor(b *testing.B) {
	workspace := generateBenchmarkWorkspace(50)
	generator := NewFixGuideGenerator(workspace)

	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/core", "@mono/api", "@mono/ui"},
		Depth: 3,
	}

	strategy := &types.FixStrategy{
		Type:           types.FixStrategyBoundaryRefactor,
		TargetPackages: []string{"@mono/ui", "@mono/core", "@mono/api"},
		Effort:         types.EffortHigh,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generator.Generate(cycle, strategy)
	}
}

// TestFixGuideGenerationPerformance verifies AC8: < 500ms for 100 packages, 5 cycles.
func TestFixGuideGenerationPerformance(t *testing.T) {
	workspace := generateBenchmarkWorkspace(100)
	cycles := generateBenchmarkCyclesWithStrategies(5, workspace)
	generator := NewFixGuideGenerator(workspace)

	// Measure time
	result := testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, cycle := range cycles {
				for j := range cycle.FixStrategies {
					generator.Generate(cycle, &cycle.FixStrategies[j])
				}
			}
		}
	})

	// Calculate per-operation time in milliseconds
	nsPerOp := result.NsPerOp()
	msPerOp := float64(nsPerOp) / 1e6

	t.Logf("Guide generation time: %.2f ms for 100 packages, 5 cycles", msPerOp)

	// AC8: Must be < 500ms
	if msPerOp > 500 {
		t.Errorf("Guide generation took %.2f ms, exceeds 500ms threshold", msPerOp)
	}
}

// generateBenchmarkWorkspace creates a workspace with n packages.
func generateBenchmarkWorkspace(n int) *types.WorkspaceData {
	packages := make(map[string]*types.PackageInfo)
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("@mono/pkg-%d", i)
		packages[name] = &types.PackageInfo{
			Name:    name,
			Version: "1.0.0",
			Path:    fmt.Sprintf("packages/pkg-%d", i),
		}
	}

	return &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages:      packages,
	}
}

// generateBenchmarkCyclesWithStrategies creates cycles with strategies for benchmarking.
func generateBenchmarkCyclesWithStrategies(n int, workspace *types.WorkspaceData) []*types.CircularDependencyInfo {
	cycles := make([]*types.CircularDependencyInfo, n)

	for i := 0; i < n; i++ {
		// Create cycle of varying depth
		depth := 2 + (i % 3) // 2, 3, 4, 2, 3...
		cyclePkgs := make([]string, depth+1)
		targetPkgs := make([]string, depth)

		for j := 0; j < depth; j++ {
			pkgName := fmt.Sprintf("@mono/pkg-%d", (i*depth+j)%100)
			cyclePkgs[j] = pkgName
			targetPkgs[j] = pkgName
		}
		cyclePkgs[depth] = cyclePkgs[0] // Close the cycle

		cycles[i] = &types.CircularDependencyInfo{
			Cycle: cyclePkgs,
			Depth: depth,
			Type:  types.CircularTypeIndirect,
			FixStrategies: []types.FixStrategy{
				{
					Type:           types.FixStrategyExtractModule,
					TargetPackages: targetPkgs,
					NewPackageName: fmt.Sprintf("@mono/shared-%d", i),
					Effort:         types.EffortMedium,
				},
				{
					Type:           types.FixStrategyDependencyInject,
					TargetPackages: targetPkgs[:2],
					Effort:         types.EffortLow,
				},
				{
					Type:           types.FixStrategyBoundaryRefactor,
					TargetPackages: targetPkgs,
					Effort:         types.EffortHigh,
				},
			},
		}
	}

	return cycles
}
