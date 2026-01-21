// Package analyzer benchmarks for RootCauseAnalyzer.
// Story 3.1 AC: Root cause analysis completes in <50ms for 50+ package cycles.
package analyzer

import (
	"fmt"
	"testing"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// BenchmarkRootCauseAnalyzer_Analyze benchmarks basic analysis performance.
func BenchmarkRootCauseAnalyzer_Analyze(b *testing.B) {
	graph := createBenchmarkGraph(10)
	cycle := &types.CircularDependencyInfo{
		Cycle:    []string{"pkg-0", "pkg-1", "pkg-2", "pkg-0"},
		Type:     types.CircularTypeIndirect,
		Severity: types.CircularSeverityInfo,
		Depth:    3,
	}

	analyzer := NewRootCauseAnalyzer(graph)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzer.Analyze(cycle)
	}
}

// BenchmarkRootCauseAnalyzer_50Packages benchmarks analysis with 50 packages.
// Story 3.1 AC: Must complete in <50ms.
func BenchmarkRootCauseAnalyzer_50Packages(b *testing.B) {
	graph := createBenchmarkGraph(50)
	cycle := createLongCycle(50)

	analyzer := NewRootCauseAnalyzer(graph)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzer.Analyze(cycle)
	}
}

// BenchmarkRootCauseAnalyzer_100Packages benchmarks analysis with 100 packages.
func BenchmarkRootCauseAnalyzer_100Packages(b *testing.B) {
	graph := createBenchmarkGraph(100)
	cycle := createLongCycle(100)

	analyzer := NewRootCauseAnalyzer(graph)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzer.Analyze(cycle)
	}
}

// BenchmarkRootCauseAnalyzer_200Packages benchmarks analysis with 200 packages.
func BenchmarkRootCauseAnalyzer_200Packages(b *testing.B) {
	graph := createBenchmarkGraph(200)
	cycle := createLongCycle(200)

	analyzer := NewRootCauseAnalyzer(graph)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzer.Analyze(cycle)
	}
}

// BenchmarkRootCauseAnalyzer_DirectCycle benchmarks analysis of direct cycle.
func BenchmarkRootCauseAnalyzer_DirectCycle(b *testing.B) {
	graph := &types.DependencyGraph{
		Nodes: map[string]*types.PackageNode{
			"pkg-a": {Name: "pkg-a", Dependencies: []string{"pkg-b"}},
			"pkg-b": {Name: "pkg-b", Dependencies: []string{"pkg-a"}},
		},
	}
	cycle := &types.CircularDependencyInfo{
		Cycle:    []string{"pkg-a", "pkg-b", "pkg-a"},
		Type:     types.CircularTypeDirect,
		Severity: types.CircularSeverityWarning,
		Depth:    2,
	}

	analyzer := NewRootCauseAnalyzer(graph)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzer.Analyze(cycle)
	}
}

// BenchmarkRootCauseAnalyzer_BuildDependencyChain benchmarks chain building.
func BenchmarkRootCauseAnalyzer_BuildDependencyChain(b *testing.B) {
	graph := createBenchmarkGraph(50)
	cycle := createLongCycleSlice(50)
	analyzer := NewRootCauseAnalyzer(graph)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzer.buildDependencyChain(cycle)
	}
}

// BenchmarkRootCauseAnalyzer_CalculateIncomingDepsScore benchmarks incoming deps scoring.
func BenchmarkRootCauseAnalyzer_CalculateIncomingDepsScore(b *testing.B) {
	graph := createBenchmarkGraph(100)
	analyzer := NewRootCauseAnalyzer(graph)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzer.calculateIncomingDepsScore("pkg-50")
	}
}

// TestRootCauseAnalyzer_PerformanceRequirement verifies the 50ms requirement.
func TestRootCauseAnalyzer_PerformanceRequirement(t *testing.T) {
	// Create graph with 50+ packages
	graph := createBenchmarkGraph(60)
	cycle := createLongCycle(60)

	analyzer := NewRootCauseAnalyzer(graph)

	// Run analysis and measure time
	result := testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			analyzer.Analyze(cycle)
		}
	})

	// Calculate average time per operation
	avgNsPerOp := float64(result.T.Nanoseconds()) / float64(result.N)
	avgMsPerOp := avgNsPerOp / 1e6

	t.Logf("Root cause analysis: %.4f ms per operation (N=%d)", avgMsPerOp, result.N)

	// Story 3.1 AC: <50ms for 50+ package cycles
	if avgMsPerOp > 50 {
		t.Errorf("Performance requirement failed: %.4f ms > 50ms", avgMsPerOp)
	}
}

// createBenchmarkGraph creates a graph with n packages, each depending on the next.
func createBenchmarkGraph(n int) *types.DependencyGraph {
	graph := &types.DependencyGraph{
		Nodes: make(map[string]*types.PackageNode),
	}

	for i := 0; i < n; i++ {
		name := fmt.Sprintf("pkg-%d", i)
		deps := []string{}

		// Each package depends on the next (circular)
		if i < n-1 {
			deps = append(deps, fmt.Sprintf("pkg-%d", i+1))
		} else {
			// Last package depends on first (closes the cycle)
			deps = append(deps, "pkg-0")
		}

		// Add some additional deps for complexity
		if i > 0 && i < n-1 {
			devDeps := []string{fmt.Sprintf("pkg-%d", (i+5)%n)}
			graph.Nodes[name] = &types.PackageNode{
				Name:            name,
				Dependencies:    deps,
				DevDependencies: devDeps,
			}
		} else {
			graph.Nodes[name] = &types.PackageNode{
				Name:         name,
				Dependencies: deps,
			}
		}
	}

	return graph
}

// createLongCycle creates a CircularDependencyInfo with n packages.
func createLongCycle(n int) *types.CircularDependencyInfo {
	return &types.CircularDependencyInfo{
		Cycle:    createLongCycleSlice(n),
		Type:     types.CircularTypeIndirect,
		Severity: types.CircularSeverityInfo,
		Depth:    n,
	}
}

// createLongCycleSlice creates a cycle slice with n packages.
func createLongCycleSlice(n int) []string {
	cycle := make([]string, n+1)
	for i := 0; i < n; i++ {
		cycle[i] = fmt.Sprintf("pkg-%d", i)
	}
	cycle[n] = "pkg-0" // Close the cycle
	return cycle
}
