// Package analyzer provides dependency graph analysis for monorepo workspaces.
// This file contains benchmark tests for fix strategy generator for Story 3.3.
package analyzer

import (
	"testing"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// ========================================
// Performance Benchmarks (AC6)
// ========================================

// BenchmarkFixStrategyGenerator_DirectCycle benchmarks generation for A ↔ B cycles.
func BenchmarkFixStrategyGenerator_DirectCycle(b *testing.B) {
	graph := types.NewDependencyGraph("/test", types.WorkspaceTypePnpm)
	graph.Nodes["@mono/a"] = types.NewPackageNode("@mono/a", "1.0.0", "packages/a")
	graph.Nodes["@mono/b"] = types.NewPackageNode("@mono/b", "1.0.0", "packages/b")
	graph.Edges = []*types.DependencyEdge{
		{From: "@mono/a", To: "@mono/b", Type: types.DependencyTypeProduction},
		{From: "@mono/b", To: "@mono/a", Type: types.DependencyTypeProduction},
	}

	workspace := &types.WorkspaceData{
		WorkspaceType: types.WorkspaceTypePnpm,
		RootPath:      "/test",
		Packages: map[string]*types.PackageInfo{
			"@mono/a": {Name: "@mono/a", Version: "1.0.0"},
			"@mono/b": {Name: "@mono/b", Version: "1.0.0"},
		},
	}

	generator := NewFixStrategyGenerator(graph, workspace)
	cycle := types.NewCircularDependencyInfo([]string{"@mono/a", "@mono/b", "@mono/a"})
	cycle.RootCause = &types.RootCauseAnalysis{
		OriginatingPackage: "@mono/a",
		CriticalEdge: &types.RootCauseEdge{
			From:     "@mono/b",
			To:       "@mono/a",
			Type:     types.DependencyTypeProduction,
			Critical: true,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = generator.Generate(cycle)
	}
}

// BenchmarkFixStrategyGenerator_IndirectCycle benchmarks generation for 3-node cycles.
func BenchmarkFixStrategyGenerator_IndirectCycle(b *testing.B) {
	graph := types.NewDependencyGraph("/test", types.WorkspaceTypePnpm)
	packages := []string{"@mono/ui", "@mono/api", "@mono/core"}
	for _, pkg := range packages {
		graph.Nodes[pkg] = types.NewPackageNode(pkg, "1.0.0", "packages/"+pkg[6:])
	}
	graph.Edges = []*types.DependencyEdge{
		{From: "@mono/ui", To: "@mono/api", Type: types.DependencyTypeProduction},
		{From: "@mono/api", To: "@mono/core", Type: types.DependencyTypeProduction},
		{From: "@mono/core", To: "@mono/ui", Type: types.DependencyTypeProduction},
	}

	workspace := &types.WorkspaceData{
		WorkspaceType: types.WorkspaceTypePnpm,
		RootPath:      "/test",
		Packages:      map[string]*types.PackageInfo{},
	}
	for _, pkg := range packages {
		workspace.Packages[pkg] = &types.PackageInfo{Name: pkg, Version: "1.0.0"}
	}

	generator := NewFixStrategyGenerator(graph, workspace)
	cycle := types.NewCircularDependencyInfo([]string{
		"@mono/ui", "@mono/api", "@mono/core", "@mono/ui",
	})
	cycle.RootCause = &types.RootCauseAnalysis{
		OriginatingPackage: "@mono/ui",
		Confidence:         75,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = generator.Generate(cycle)
	}
}

// BenchmarkFixStrategyGenerator_LongCycle benchmarks generation for 10-node cycles.
func BenchmarkFixStrategyGenerator_LongCycle(b *testing.B) {
	graph := types.NewDependencyGraph("/test", types.WorkspaceTypePnpm)
	packages := make([]string, 10)
	for i := 0; i < 10; i++ {
		packages[i] = "@mono/pkg" + string(rune('a'+i))
		graph.Nodes[packages[i]] = types.NewPackageNode(packages[i], "1.0.0", "packages/"+packages[i][6:])
	}
	// Create cycle: a -> b -> c -> ... -> j -> a
	for i := 0; i < 10; i++ {
		next := (i + 1) % 10
		graph.Edges = append(graph.Edges, &types.DependencyEdge{
			From: packages[i],
			To:   packages[next],
			Type: types.DependencyTypeProduction,
		})
	}

	workspace := &types.WorkspaceData{
		WorkspaceType: types.WorkspaceTypePnpm,
		RootPath:      "/test",
		Packages:      map[string]*types.PackageInfo{},
	}
	for _, pkg := range packages {
		workspace.Packages[pkg] = &types.PackageInfo{Name: pkg, Version: "1.0.0"}
	}

	cycleNodes := append(packages, packages[0])
	generator := NewFixStrategyGenerator(graph, workspace)
	cycle := types.NewCircularDependencyInfo(cycleNodes)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = generator.Generate(cycle)
	}
}

// ========================================
// Performance Requirement Tests
// ========================================

// TestFixStrategyGenerator_PerformanceRequirement verifies < 5ms per cycle.
func TestFixStrategyGenerator_PerformanceRequirement(t *testing.T) {
	graph := types.NewDependencyGraph("/test", types.WorkspaceTypePnpm)
	packages := make([]string, 20)
	for i := 0; i < 20; i++ {
		packages[i] = "@mono/pkg" + string(rune('a'+i))
		graph.Nodes[packages[i]] = types.NewPackageNode(packages[i], "1.0.0", "packages/"+packages[i][6:])
	}
	// Create long cycle
	for i := 0; i < 20; i++ {
		next := (i + 1) % 20
		graph.Edges = append(graph.Edges, &types.DependencyEdge{
			From: packages[i],
			To:   packages[next],
			Type: types.DependencyTypeProduction,
		})
	}

	workspace := &types.WorkspaceData{
		WorkspaceType: types.WorkspaceTypePnpm,
		RootPath:      "/test",
		Packages:      map[string]*types.PackageInfo{},
	}
	for _, pkg := range packages {
		workspace.Packages[pkg] = &types.PackageInfo{Name: pkg, Version: "1.0.0"}
	}

	cycleNodes := append(packages, packages[0])
	generator := NewFixStrategyGenerator(graph, workspace)
	cycle := types.NewCircularDependencyInfo(cycleNodes)

	// Run 100 iterations and measure total time
	iterations := 100
	start := testing.Benchmark(func(b *testing.B) {
		for i := 0; i < iterations; i++ {
			_ = generator.Generate(cycle)
		}
	})

	avgNs := start.NsPerOp() / int64(iterations)
	avgMs := float64(avgNs) / 1e6

	t.Logf("Average time per cycle (20 packages): %.3fms", avgMs)

	// AC6: Strategy generation should complete < 5ms per cycle
	if avgMs > 5.0 {
		t.Errorf("Performance requirement not met: %.3fms > 5ms", avgMs)
	}
}

// TestFixStrategyGenerator_MultipleCyclesPerformance verifies performance with many cycles.
func TestFixStrategyGenerator_MultipleCyclesPerformance(t *testing.T) {
	// Create graph with 100 packages
	graph := types.NewDependencyGraph("/test", types.WorkspaceTypePnpm)
	for i := 0; i < 100; i++ {
		name := "@mono/pkg" + string(rune('a'+i/26)) + string(rune('a'+i%26))
		graph.Nodes[name] = types.NewPackageNode(name, "1.0.0", "packages/"+name[6:])
	}

	workspace := &types.WorkspaceData{
		WorkspaceType: types.WorkspaceTypePnpm,
		RootPath:      "/test",
		Packages:      map[string]*types.PackageInfo{},
	}
	for name := range graph.Nodes {
		workspace.Packages[name] = &types.PackageInfo{Name: name, Version: "1.0.0"}
	}

	generator := NewFixStrategyGenerator(graph, workspace)

	// Generate 50 different cycles
	cycles := make([]*types.CircularDependencyInfo, 50)
	for i := 0; i < 50; i++ {
		size := 2 + (i % 5) // Cycles of size 2-6
		cycleNodes := make([]string, size+1)
		for j := 0; j <= size; j++ {
			idx := (i*3 + j) % 100
			name := "@mono/pkg" + string(rune('a'+idx/26)) + string(rune('a'+idx%26))
			cycleNodes[j] = name
		}
		cycleNodes[size] = cycleNodes[0] // Close the cycle
		cycles[i] = types.NewCircularDependencyInfo(cycleNodes)
	}

	// Measure time to generate strategies for all 50 cycles
	start := testing.Benchmark(func(b *testing.B) {
		for _, cycle := range cycles {
			_ = generator.Generate(cycle)
		}
	})

	totalMs := float64(start.NsPerOp()) / 1e6
	avgMs := totalMs / float64(len(cycles))

	t.Logf("50 cycles total: %.3fms, average: %.3fms per cycle", totalMs, avgMs)

	// Total should be < 250ms (50 cycles * 5ms each)
	if totalMs > 250.0 {
		t.Errorf("Performance requirement not met for multiple cycles: %.3fms > 250ms", totalMs)
	}
}
