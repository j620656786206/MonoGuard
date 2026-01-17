package analyzer

import (
	"fmt"
	"testing"
	"time"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// generateWorkspaceWithConflicts creates a workspace with specified number of packages
// and conflicts for benchmarking purposes.
func generateWorkspaceWithConflicts(packageCount, conflictCount int) *types.DependencyGraph {
	graph := types.NewDependencyGraph("/benchmark", types.WorkspaceTypeNpm)

	// Common external dependencies that will have conflicts
	conflictDeps := make([]string, conflictCount)
	for i := 0; i < conflictCount; i++ {
		conflictDeps[i] = fmt.Sprintf("external-dep-%d", i)
	}

	// Generate packages
	for i := 0; i < packageCount; i++ {
		pkgName := fmt.Sprintf("@mono/pkg-%d", i)
		node := types.NewPackageNode(pkgName, "1.0.0", fmt.Sprintf("/benchmark/packages/pkg-%d", i))

		// Add external dependencies with version conflicts
		for j, depName := range conflictDeps {
			// Create version variance: packages 0-49 use one version, 50-99 use another
			var version string
			if i%2 == 0 {
				version = fmt.Sprintf("^4.%d.0", j)
			} else {
				version = fmt.Sprintf("^4.%d.1", j) // Patch difference
			}
			node.ExternalDeps[depName] = version
		}

		// Add some non-conflicting deps
		node.ExternalDeps["unique-dep-"+pkgName] = "^1.0.0"

		graph.Nodes[pkgName] = node
	}

	return graph
}

// generateWorkspaceWithMajorConflicts creates a workspace with major version conflicts.
func generateWorkspaceWithMajorConflicts(packageCount, conflictCount int) *types.DependencyGraph {
	graph := types.NewDependencyGraph("/benchmark", types.WorkspaceTypeNpm)

	conflictDeps := make([]string, conflictCount)
	for i := 0; i < conflictCount; i++ {
		conflictDeps[i] = fmt.Sprintf("external-dep-%d", i)
	}

	for i := 0; i < packageCount; i++ {
		pkgName := fmt.Sprintf("@mono/pkg-%d", i)
		node := types.NewPackageNode(pkgName, "1.0.0", fmt.Sprintf("/benchmark/packages/pkg-%d", i))

		for j, depName := range conflictDeps {
			// Create major version variance
			majorVersion := (i % 3) + 1 // versions 1, 2, or 3
			version := fmt.Sprintf("^%d.%d.0", majorVersion, j)
			node.ExternalDeps[depName] = version
		}

		graph.Nodes[pkgName] = node
	}

	return graph
}

// BenchmarkDetectConflicts100Packages benchmarks conflict detection with 100 packages.
// AC6: Must complete in < 1 second.
func BenchmarkDetectConflicts100Packages(b *testing.B) {
	graph := generateWorkspaceWithConflicts(100, 20)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		detector := NewConflictDetector(graph)
		_ = detector.DetectConflicts()
	}
}

// BenchmarkDetectConflicts100PackagesWithMajorConflicts tests with critical severity conflicts.
func BenchmarkDetectConflicts100PackagesWithMajorConflicts(b *testing.B) {
	graph := generateWorkspaceWithMajorConflicts(100, 20)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		detector := NewConflictDetector(graph)
		_ = detector.DetectConflicts()
	}
}

// BenchmarkDetectConflicts200Packages tests scalability with larger workspaces.
func BenchmarkDetectConflicts200Packages(b *testing.B) {
	graph := generateWorkspaceWithConflicts(200, 30)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		detector := NewConflictDetector(graph)
		_ = detector.DetectConflicts()
	}
}

// BenchmarkDetectConflicts500Packages tests with very large workspaces.
func BenchmarkDetectConflicts500Packages(b *testing.B) {
	graph := generateWorkspaceWithConflicts(500, 50)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		detector := NewConflictDetector(graph)
		_ = detector.DetectConflicts()
	}
}

// TestPerformanceRequirement_100Packages explicitly tests the AC6 requirement:
// "Given a workspace with 100 packages, When conflict detection completes, Then it finishes in < 1 second"
func TestPerformanceRequirement_100Packages(t *testing.T) {
	graph := generateWorkspaceWithConflicts(100, 20)

	// Run conflict detection and measure time
	start := time.Now()
	detector := NewConflictDetector(graph)
	conflicts := detector.DetectConflicts()
	duration := time.Since(start)

	// Verify it completes in < 1 second (AC6)
	if duration >= time.Second {
		t.Errorf("Conflict detection took %v, want < 1 second", duration)
	}

	// Log performance info
	t.Logf("100 packages with 20 conflict-prone deps: %v, found %d conflicts", duration, len(conflicts))

	// Sanity check: should find conflicts
	if len(conflicts) == 0 {
		t.Error("Expected to find conflicts in test data")
	}
}

// TestPerformanceRequirement_100PackagesRealistic tests with more realistic conflict patterns.
func TestPerformanceRequirement_100PackagesRealistic(t *testing.T) {
	// Create realistic workspace with mixed dependency types
	graph := types.NewDependencyGraph("/benchmark", types.WorkspaceTypeNpm)

	commonDeps := []struct {
		name     string
		versions []string
	}{
		{"lodash", []string{"^4.17.19", "^4.17.21"}},
		{"react", []string{"^17.0.0", "^18.0.0"}},
		{"typescript", []string{"^4.9.0", "^5.0.0"}},
		{"axios", []string{"^0.27.0", "^1.0.0"}},
		{"moment", []string{"^2.29.1", "^2.29.4"}},
	}

	for i := 0; i < 100; i++ {
		pkgName := fmt.Sprintf("@mono/pkg-%d", i)
		node := types.NewPackageNode(pkgName, "1.0.0", fmt.Sprintf("/packages/pkg-%d", i))

		// Add common deps with version variance
		for _, dep := range commonDeps {
			versionIdx := i % len(dep.versions)
			node.ExternalDeps[dep.name] = dep.versions[versionIdx]
		}

		// Add some unique deps
		node.ExternalDeps[fmt.Sprintf("unique-%d", i)] = "^1.0.0"
		node.ExternalDevDeps["jest"] = "^29.0.0"

		graph.Nodes[pkgName] = node
	}

	// Measure performance
	start := time.Now()
	detector := NewConflictDetector(graph)
	conflicts := detector.DetectConflicts()
	duration := time.Since(start)

	// AC6: < 1 second
	if duration >= time.Second {
		t.Errorf("Conflict detection took %v, want < 1 second", duration)
	}

	t.Logf("Realistic 100 packages: %v, found %d conflicts", duration, len(conflicts))

	// Should find conflicts for lodash, react, typescript, axios, moment
	if len(conflicts) < 5 {
		t.Errorf("Expected at least 5 conflicts, found %d", len(conflicts))
	}
}
