// Package analyzer provides dependency graph analysis for monorepo workspaces.
// This file contains benchmark tests for complexity calculator (Story 3.5).
package analyzer

import (
	"fmt"
	"testing"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// BenchmarkComplexityCalculator_Calculate benchmarks Calculate for various cycle sizes.
func BenchmarkComplexityCalculator_Calculate(b *testing.B) {
	sizes := []int{2, 5, 10, 20}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("depth_%d", size), func(b *testing.B) {
			graph := types.NewDependencyGraph("@mono/root", types.WorkspaceTypeNpm)
			workspace := &types.WorkspaceData{
				RootPath:      "@mono/root",
				WorkspaceType: types.WorkspaceTypeNpm,
				Packages:      make(map[string]*types.PackageInfo),
			}

			// Create packages
			cycle := make([]string, size+1)
			for i := 0; i < size; i++ {
				pkgName := fmt.Sprintf("@mono/pkg%d", i)
				cycle[i] = pkgName
				workspace.Packages[pkgName] = &types.PackageInfo{
					Name:    pkgName,
					Version: "1.0.0",
				}
			}
			cycle[size] = cycle[0] // Close the cycle

			cycleInfo := &types.CircularDependencyInfo{
				Cycle: cycle,
				Depth: size,
			}

			calc := NewComplexityCalculator(graph, workspace)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = calc.Calculate(cycleInfo)
			}
		})
	}
}

// BenchmarkComplexityCalculator_WithImportTraces benchmarks with ImportTraces.
func BenchmarkComplexityCalculator_WithImportTraces(b *testing.B) {
	sizes := []int{3, 10, 50}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("traces_%d", size), func(b *testing.B) {
			graph := types.NewDependencyGraph("@mono/root", types.WorkspaceTypeNpm)
			workspace := &types.WorkspaceData{
				RootPath:      "@mono/root",
				WorkspaceType: types.WorkspaceTypeNpm,
				Packages: map[string]*types.PackageInfo{
					"@mono/ui":   {Name: "@mono/ui", Version: "1.0.0"},
					"@mono/api":  {Name: "@mono/api", Version: "1.0.0"},
					"@mono/core": {Name: "@mono/core", Version: "1.0.0"},
				},
			}

			// Create import traces
			traces := make([]types.ImportTrace, size)
			for i := 0; i < size; i++ {
				traces[i] = types.ImportTrace{
					FromPackage: "@mono/ui",
					ToPackage:   "@mono/api",
					FilePath:    fmt.Sprintf("packages/ui/src/file%d.ts", i),
					LineNumber:  i + 1,
					ImportType:  types.ImportTypeESMNamed,
				}
			}

			cycleInfo := &types.CircularDependencyInfo{
				Cycle:        []string{"@mono/ui", "@mono/api", "@mono/core", "@mono/ui"},
				Depth:        3,
				ImportTraces: traces,
			}

			calc := NewComplexityCalculator(graph, workspace)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = calc.Calculate(cycleInfo)
			}
		})
	}
}

// BenchmarkComplexityCalculator_ExternalDeps benchmarks external dependency detection.
func BenchmarkComplexityCalculator_ExternalDeps(b *testing.B) {
	graph := types.NewDependencyGraph("@mono/root", types.WorkspaceTypeNpm)
	workspace := &types.WorkspaceData{
		RootPath:      "@mono/root",
		WorkspaceType: types.WorkspaceTypeNpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/ui":  {Name: "@mono/ui", Version: "1.0.0"},
			"@mono/api": {Name: "@mono/api", Version: "1.0.0"},
		},
	}

	// Create root cause chain with external deps
	chain := make([]types.RootCauseEdge, 20)
	for i := 0; i < 20; i++ {
		chain[i] = types.RootCauseEdge{
			From: fmt.Sprintf("@mono/pkg%d", i),
			To:   fmt.Sprintf("@external/lib%d", i),
			Type: types.DependencyTypeProduction,
		}
	}

	cycleInfo := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/api", "@mono/ui"},
		Depth: 2,
		RootCause: &types.RootCauseAnalysis{
			Chain: chain,
		},
	}

	calc := NewComplexityCalculator(graph, workspace)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = calc.Calculate(cycleInfo)
	}
}

// BenchmarkComplexityCalculator_200Packages benchmarks 200 package scenario.
// Target: < 100ms as per AC7.
func BenchmarkComplexityCalculator_200Packages(b *testing.B) {
	graph := types.NewDependencyGraph("@mono/root", types.WorkspaceTypeNpm)
	workspace := &types.WorkspaceData{
		RootPath:      "@mono/root",
		WorkspaceType: types.WorkspaceTypeNpm,
		Packages:      make(map[string]*types.PackageInfo),
	}

	// Create 200 packages
	cycle := make([]string, 11) // 10 packages in cycle
	for i := 0; i < 200; i++ {
		pkgName := fmt.Sprintf("@mono/pkg%d", i)
		workspace.Packages[pkgName] = &types.PackageInfo{
			Name:    pkgName,
			Version: "1.0.0",
		}
		if i < 10 {
			cycle[i] = pkgName
		}
	}
	cycle[10] = cycle[0] // Close the cycle

	cycleInfo := &types.CircularDependencyInfo{
		Cycle: cycle,
		Depth: 10,
	}

	calc := NewComplexityCalculator(graph, workspace)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = calc.Calculate(cycleInfo)
	}
}

// BenchmarkComplexityCalculator_AC7Scenario benchmarks exact AC7 requirement:
// "Given a workspace with 100 packages and 5 cycles, calculation completes in < 100ms"
func BenchmarkComplexityCalculator_AC7Scenario(b *testing.B) {
	graph := types.NewDependencyGraph("@mono/root", types.WorkspaceTypeNpm)
	workspace := &types.WorkspaceData{
		RootPath:      "@mono/root",
		WorkspaceType: types.WorkspaceTypeNpm,
		Packages:      make(map[string]*types.PackageInfo),
	}

	// Create exactly 100 packages
	for i := 0; i < 100; i++ {
		pkgName := fmt.Sprintf("@mono/pkg%d", i)
		workspace.Packages[pkgName] = &types.PackageInfo{
			Name:    pkgName,
			Version: "1.0.0",
		}
	}

	// Create 5 cycles of varying sizes (as per AC7)
	cycles := []*types.CircularDependencyInfo{
		{Cycle: []string{"@mono/pkg0", "@mono/pkg1", "@mono/pkg0"}, Depth: 2},
		{Cycle: []string{"@mono/pkg10", "@mono/pkg11", "@mono/pkg12", "@mono/pkg10"}, Depth: 3},
		{Cycle: []string{"@mono/pkg20", "@mono/pkg21", "@mono/pkg22", "@mono/pkg23", "@mono/pkg20"}, Depth: 4},
		{Cycle: []string{"@mono/pkg30", "@mono/pkg31", "@mono/pkg32", "@mono/pkg30"}, Depth: 3},
		{Cycle: []string{"@mono/pkg40", "@mono/pkg41", "@mono/pkg40"}, Depth: 2},
	}

	calc := NewComplexityCalculator(graph, workspace)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, cycle := range cycles {
			_ = calc.Calculate(cycle)
		}
	}
}
