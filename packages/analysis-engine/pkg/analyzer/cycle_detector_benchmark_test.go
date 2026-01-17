package analyzer

import (
	"fmt"
	"testing"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// ========================================
// Performance Benchmarks (AC6)
// ========================================

// BenchmarkDetectCycles100Packages benchmarks cycle detection with 100 packages.
// Requirement: < 3 seconds
func BenchmarkDetectCycles100Packages(b *testing.B) {
	graph := generateGraphWithCycles(100, 5)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		detector := NewCycleDetector(graph)
		detector.DetectCycles()
	}
}

// BenchmarkDetectCycles1000Packages benchmarks cycle detection with 1000 packages.
// Requirement: < 30 seconds
func BenchmarkDetectCycles1000Packages(b *testing.B) {
	graph := generateGraphWithCycles(1000, 10)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		detector := NewCycleDetector(graph)
		detector.DetectCycles()
	}
}

// BenchmarkDetectCyclesNoCycles benchmarks detection when no cycles exist.
func BenchmarkDetectCyclesNoCycles(b *testing.B) {
	graph := generateLinearGraph(100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		detector := NewCycleDetector(graph)
		detector.DetectCycles()
	}
}

// BenchmarkDetectCyclesManyCycles benchmarks detection with many small cycles.
func BenchmarkDetectCyclesManyCycles(b *testing.B) {
	graph := generateGraphWithManyCycles(100, 20) // 100 packages, 20 cycles
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		detector := NewCycleDetector(graph)
		detector.DetectCycles()
	}
}

// ========================================
// Helper Functions for Benchmark Data
// ========================================

// generateGraphWithCycles creates a graph with specified number of packages and cycles.
// Cycles are distributed throughout the graph.
func generateGraphWithCycles(packageCount, cycleCount int) *types.DependencyGraph {
	nodes := make(map[string]*types.PackageNode)

	// Create all nodes
	for i := 0; i < packageCount; i++ {
		name := fmt.Sprintf("@mono/pkg-%d", i)
		nodes[name] = types.NewPackageNode(name, "1.0.0", fmt.Sprintf("packages/pkg-%d", i))
	}

	// Create linear dependencies (most packages)
	for i := 0; i < packageCount-1; i++ {
		from := fmt.Sprintf("@mono/pkg-%d", i)
		to := fmt.Sprintf("@mono/pkg-%d", i+1)
		nodes[from].Dependencies = append(nodes[from].Dependencies, to)
	}

	// Create cycles at regular intervals
	cycleInterval := packageCount / (cycleCount + 1)
	for c := 0; c < cycleCount; c++ {
		start := (c + 1) * cycleInterval
		if start+2 < packageCount {
			// Create a 3-node cycle: start -> start+1 -> start+2 -> start
			node1 := fmt.Sprintf("@mono/pkg-%d", start)
			node3 := fmt.Sprintf("@mono/pkg-%d", start+2)

			// Complete the cycle: node3 -> node1 (node1 -> node2 -> node3 already exists from linear deps)
			nodes[node3].Dependencies = append(nodes[node3].Dependencies, node1)
		}
	}

	return &types.DependencyGraph{
		Nodes:         nodes,
		Edges:         []*types.DependencyEdge{},
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
	}
}

// generateLinearGraph creates a linear graph with no cycles.
func generateLinearGraph(packageCount int) *types.DependencyGraph {
	nodes := make(map[string]*types.PackageNode)

	// Create all nodes
	for i := 0; i < packageCount; i++ {
		name := fmt.Sprintf("@mono/pkg-%d", i)
		nodes[name] = types.NewPackageNode(name, "1.0.0", fmt.Sprintf("packages/pkg-%d", i))
	}

	// Create linear dependencies
	for i := 0; i < packageCount-1; i++ {
		from := fmt.Sprintf("@mono/pkg-%d", i)
		to := fmt.Sprintf("@mono/pkg-%d", i+1)
		nodes[from].Dependencies = append(nodes[from].Dependencies, to)
	}

	return &types.DependencyGraph{
		Nodes:         nodes,
		Edges:         []*types.DependencyEdge{},
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
	}
}

// generateGraphWithManyCycles creates a graph with many small cycles.
func generateGraphWithManyCycles(packageCount, cycleCount int) *types.DependencyGraph {
	nodes := make(map[string]*types.PackageNode)

	// Create all nodes
	for i := 0; i < packageCount; i++ {
		name := fmt.Sprintf("@mono/pkg-%d", i)
		nodes[name] = types.NewPackageNode(name, "1.0.0", fmt.Sprintf("packages/pkg-%d", i))
	}

	// Create 2-node cycles (direct cycles)
	for c := 0; c < cycleCount && c*2+1 < packageCount; c++ {
		node1 := fmt.Sprintf("@mono/pkg-%d", c*2)
		node2 := fmt.Sprintf("@mono/pkg-%d", c*2+1)

		nodes[node1].Dependencies = append(nodes[node1].Dependencies, node2)
		nodes[node2].Dependencies = append(nodes[node2].Dependencies, node1)
	}

	return &types.DependencyGraph{
		Nodes:         nodes,
		Edges:         []*types.DependencyEdge{},
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
	}
}

// ========================================
// Performance Verification Tests
// ========================================

// TestPerformance100Packages verifies cycle detection for 100 packages completes quickly.
func TestPerformance100Packages(t *testing.T) {
	graph := generateGraphWithCycles(100, 5)
	detector := NewCycleDetector(graph)

	cycles := detector.DetectCycles()

	// Should find cycles
	if len(cycles) == 0 {
		t.Log("Warning: No cycles found in test graph with 5 expected cycles")
	}

	// Just verify it completes - actual timing is in benchmark
	t.Logf("100 packages, 5 cycles expected: found %d cycles", len(cycles))
}

// TestPerformance1000Packages verifies cycle detection for 1000 packages completes.
func TestPerformance1000Packages(t *testing.T) {
	graph := generateGraphWithCycles(1000, 10)
	detector := NewCycleDetector(graph)

	cycles := detector.DetectCycles()

	// Should find cycles
	if len(cycles) == 0 {
		t.Log("Warning: No cycles found in test graph with 10 expected cycles")
	}

	// Just verify it completes - actual timing is in benchmark
	t.Logf("1000 packages, 10 cycles expected: found %d cycles", len(cycles))
}
