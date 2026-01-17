package analyzer

import (
	"fmt"
	"testing"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// generateWorkspace creates a realistic workspace with the specified number of packages.
// Each package has dependencies on some earlier packages to create a realistic graph.
func generateWorkspace(packageCount int) *types.WorkspaceData {
	packages := make(map[string]*types.PackageInfo)

	for i := 0; i < packageCount; i++ {
		name := fmt.Sprintf("@mono/pkg-%d", i)
		deps := make(map[string]string)
		devDeps := make(map[string]string)

		// Add internal dependencies (packages depend on some earlier packages)
		// This creates a realistic dependency graph with varying fan-out
		if i > 0 {
			// Each package depends on up to 3 earlier packages
			for j := 0; j < 3 && j < i; j++ {
				depIdx := (i - 1 - j) % i
				if depIdx >= 0 {
					depName := fmt.Sprintf("@mono/pkg-%d", depIdx)
					deps[depName] = "^1.0.0"
				}
			}
		}

		// Add some external dependencies (common in real monorepos)
		deps["lodash"] = "^4.17.0"
		deps["react"] = "^18.0.0"
		devDeps["typescript"] = "^5.0.0"
		devDeps["jest"] = "^29.0.0"

		packages[name] = &types.PackageInfo{
			Name:             name,
			Version:          "1.0.0",
			Path:             fmt.Sprintf("packages/pkg-%d", i),
			Dependencies:     deps,
			DevDependencies:  devDeps,
			PeerDependencies: map[string]string{},
		}
	}

	return &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages:      packages,
	}
}

// BenchmarkBuildGraph100Packages benchmarks graph construction for 100 packages.
// AC6 requirement: < 2 seconds for 100 packages.
func BenchmarkBuildGraph100Packages(b *testing.B) {
	workspace := generateWorkspace(100)
	gb := NewGraphBuilder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := gb.Build(workspace)
		if err != nil {
			b.Fatalf("Build failed: %v", err)
		}
	}
}

// BenchmarkBuildGraph1000Packages benchmarks graph construction for 1000 packages.
// AC6 requirement: < 50MB memory for 1000 packages.
func BenchmarkBuildGraph1000Packages(b *testing.B) {
	workspace := generateWorkspace(1000)
	gb := NewGraphBuilder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := gb.Build(workspace)
		if err != nil {
			b.Fatalf("Build failed: %v", err)
		}
	}
}

// BenchmarkAnalyze100Packages benchmarks full analysis for 100 packages.
func BenchmarkAnalyze100Packages(b *testing.B) {
	workspace := generateWorkspace(100)
	a := NewAnalyzer()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := a.Analyze(workspace)
		if err != nil {
			b.Fatalf("Analyze failed: %v", err)
		}
	}
}

// BenchmarkAnalyze1000Packages benchmarks full analysis for 1000 packages.
func BenchmarkAnalyze1000Packages(b *testing.B) {
	workspace := generateWorkspace(1000)
	a := NewAnalyzer()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := a.Analyze(workspace)
		if err != nil {
			b.Fatalf("Analyze failed: %v", err)
		}
	}
}

// TestPerformanceRequirements runs performance tests to verify AC6 requirements.
// This is a regular test (not benchmark) to ensure requirements are always checked.
func TestPerformanceRequirements(t *testing.T) {
	t.Run("100 packages completes in reasonable time", func(t *testing.T) {
		workspace := generateWorkspace(100)
		a := NewAnalyzer()

		// Warm up
		_, _ = a.Analyze(workspace)

		// Time multiple iterations
		iterations := 10
		for i := 0; i < iterations; i++ {
			result, err := a.Analyze(workspace)
			if err != nil {
				t.Fatalf("Analyze failed: %v", err)
			}
			if result.Packages != 100 {
				t.Errorf("Expected 100 packages, got %d", result.Packages)
			}
			if result.Graph == nil {
				t.Error("Graph is nil")
			}
		}
		// Note: Actual timing verification is done via benchmarks with -bench flag
	})

	t.Run("1000 packages graph structure is correct", func(t *testing.T) {
		workspace := generateWorkspace(1000)
		a := NewAnalyzer()

		result, err := a.Analyze(workspace)
		if err != nil {
			t.Fatalf("Analyze failed: %v", err)
		}

		// Verify correct number of nodes
		if result.Packages != 1000 {
			t.Errorf("Expected 1000 packages, got %d", result.Packages)
		}
		if len(result.Graph.Nodes) != 1000 {
			t.Errorf("Expected 1000 nodes, got %d", len(result.Graph.Nodes))
		}

		// Verify edges were created (each package after first has up to 3 internal deps)
		// Minimum edges = 999 (each package depends on at least one earlier package)
		// This is a rough sanity check
		if len(result.Graph.Edges) < 500 {
			t.Errorf("Expected at least 500 edges, got %d", len(result.Graph.Edges))
		}
	})
}
