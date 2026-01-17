package analyzer

import (
	"fmt"
	"testing"
	"time"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// generateHealthTestWorkspace creates a workspace with specified number of packages
// and various issues for health score benchmarking.
func generateHealthTestWorkspace(packageCount int, cycleCount int, conflictCount int) (
	*types.DependencyGraph,
	[]*types.CircularDependencyInfo,
	[]*types.VersionConflictInfo,
) {
	graph := types.NewDependencyGraph("/benchmark", types.WorkspaceTypeNpm)

	// Create packages with varying dependency depths
	for i := 0; i < packageCount; i++ {
		pkgName := fmt.Sprintf("@mono/pkg-%d", i)
		node := types.NewPackageNode(pkgName, "1.0.0", fmt.Sprintf("/benchmark/packages/pkg-%d", i))

		// Create dependencies to create depth
		if i > 0 {
			// Depend on previous package to create a chain
			prevPkg := fmt.Sprintf("@mono/pkg-%d", i-1)
			node.Dependencies = append(node.Dependencies, prevPkg)
		}

		// Add some external dependencies
		node.ExternalDeps["lodash"] = fmt.Sprintf("^4.17.%d", i%10)
		node.ExternalDevDeps["jest"] = "^29.0.0"

		graph.Nodes[pkgName] = node
	}

	// Generate cycles
	cycles := make([]*types.CircularDependencyInfo, cycleCount)
	for i := 0; i < cycleCount; i++ {
		cycleType := types.CircularTypeDirect
		if i%2 == 0 {
			cycleType = types.CircularTypeIndirect
		}
		cycles[i] = &types.CircularDependencyInfo{
			Cycle:    []string{fmt.Sprintf("pkg-%d", i), fmt.Sprintf("pkg-%d", i+1), fmt.Sprintf("pkg-%d", i)},
			Type:     cycleType,
			Severity: types.CircularSeverityWarning,
			Depth:    2,
		}
	}

	// Generate conflicts
	conflicts := make([]*types.VersionConflictInfo, conflictCount)
	for i := 0; i < conflictCount; i++ {
		severity := types.ConflictSeverityInfo
		if i%3 == 0 {
			severity = types.ConflictSeverityCritical
		} else if i%3 == 1 {
			severity = types.ConflictSeverityWarning
		}
		conflicts[i] = &types.VersionConflictInfo{
			PackageName: fmt.Sprintf("external-dep-%d", i),
			Severity:    severity,
		}
	}

	return graph, cycles, conflicts
}

// BenchmarkCalculateHealth100Packages benchmarks health calculation with 100 packages.
// AC8: Must complete in < 100ms.
func BenchmarkCalculateHealth100Packages(b *testing.B) {
	graph, cycles, conflicts := generateHealthTestWorkspace(100, 5, 10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc := NewHealthCalculator(graph, cycles, conflicts)
		_ = calc.Calculate()
	}
}

// BenchmarkCalculateHealth100PackagesNoCycles tests without cycles.
func BenchmarkCalculateHealth100PackagesNoCycles(b *testing.B) {
	graph, _, conflicts := generateHealthTestWorkspace(100, 0, 10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc := NewHealthCalculator(graph, nil, conflicts)
		_ = calc.Calculate()
	}
}

// BenchmarkCalculateHealth200Packages tests scalability.
func BenchmarkCalculateHealth200Packages(b *testing.B) {
	graph, cycles, conflicts := generateHealthTestWorkspace(200, 10, 20)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc := NewHealthCalculator(graph, cycles, conflicts)
		_ = calc.Calculate()
	}
}

// BenchmarkCalculateHealth500Packages tests with larger workspace.
func BenchmarkCalculateHealth500Packages(b *testing.B) {
	graph, cycles, conflicts := generateHealthTestWorkspace(500, 15, 30)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc := NewHealthCalculator(graph, cycles, conflicts)
		_ = calc.Calculate()
	}
}

// TestPerformanceRequirement_HealthScore100Packages explicitly tests AC8:
// "Given a workspace with 100 packages, When health score calculation completes, Then it finishes in < 100ms"
func TestPerformanceRequirement_HealthScore100Packages(t *testing.T) {
	graph, cycles, conflicts := generateHealthTestWorkspace(100, 5, 10)

	// Run health calculation and measure time
	start := time.Now()
	calc := NewHealthCalculator(graph, cycles, conflicts)
	result := calc.Calculate()
	duration := time.Since(start)

	// Verify it completes in < 100ms (AC8)
	if duration >= 100*time.Millisecond {
		t.Errorf("Health calculation took %v, want < 100ms", duration)
	}

	// Log performance info
	t.Logf("100 packages with 5 cycles, 10 conflicts: %v, score=%d (%s)",
		duration, result.Overall, result.Rating)

	// Sanity checks
	if result.Overall < 0 || result.Overall > 100 {
		t.Errorf("Score %d out of range 0-100", result.Overall)
	}

	if result.Breakdown == nil {
		t.Error("Breakdown should not be nil")
	}

	if len(result.Factors) != 4 {
		t.Errorf("Expected 4 factors, got %d", len(result.Factors))
	}
}

// TestPerformanceRequirement_HealthScore500Packages tests larger scale.
func TestPerformanceRequirement_HealthScore500Packages(t *testing.T) {
	graph, cycles, conflicts := generateHealthTestWorkspace(500, 20, 50)

	start := time.Now()
	calc := NewHealthCalculator(graph, cycles, conflicts)
	result := calc.Calculate()
	duration := time.Since(start)

	// Even 500 packages should be fast (< 500ms)
	if duration >= 500*time.Millisecond {
		t.Errorf("Health calculation took %v, want < 500ms for 500 packages", duration)
	}

	t.Logf("500 packages with 20 cycles, 50 conflicts: %v, score=%d (%s)",
		duration, result.Overall, result.Rating)
}

// TestPerformanceRequirement_DepthCalculation tests depth calculation performance.
func TestPerformanceRequirement_DepthCalculation(t *testing.T) {
	// Create a deep chain
	graph := types.NewDependencyGraph("/test", types.WorkspaceTypeNpm)

	for i := 0; i < 100; i++ {
		pkgName := fmt.Sprintf("@mono/pkg-%d", i)
		node := types.NewPackageNode(pkgName, "1.0.0", fmt.Sprintf("/packages/pkg-%d", i))

		if i > 0 {
			node.Dependencies = []string{fmt.Sprintf("@mono/pkg-%d", i-1)}
		}

		graph.Nodes[pkgName] = node
	}

	calc := NewHealthCalculator(graph, nil, nil)

	start := time.Now()
	maxDepth, avgDepth := calc.calculateDepthMetrics()
	duration := time.Since(start)

	t.Logf("Depth calculation for 100-package chain: %v, maxDepth=%d, avgDepth=%.1f",
		duration, maxDepth, avgDepth)

	if duration >= 50*time.Millisecond {
		t.Errorf("Depth calculation took %v, want < 50ms", duration)
	}
}

// TestPerformanceRequirement_CouplingCalculation tests coupling calculation performance.
func TestPerformanceRequirement_CouplingCalculation(t *testing.T) {
	graph := types.NewDependencyGraph("/test", types.WorkspaceTypeNpm)

	// Create a complex graph with many connections
	for i := 0; i < 100; i++ {
		pkgName := fmt.Sprintf("@mono/pkg-%d", i)
		node := types.NewPackageNode(pkgName, "1.0.0", fmt.Sprintf("/packages/pkg-%d", i))

		// Each package depends on several others
		for j := 0; j < 5 && i-j-1 >= 0; j++ {
			node.Dependencies = append(node.Dependencies, fmt.Sprintf("@mono/pkg-%d", i-j-1))
		}

		graph.Nodes[pkgName] = node
	}

	calc := NewHealthCalculator(graph, nil, nil)

	start := time.Now()
	metrics := calc.calculateCouplingMetrics()
	duration := time.Since(start)

	t.Logf("Coupling calculation for 100 packages: %v, avgInstability=%.2f",
		duration, metrics.AverageInstability)

	if duration >= 50*time.Millisecond {
		t.Errorf("Coupling calculation took %v, want < 50ms", duration)
	}
}
