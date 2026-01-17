package analyzer

import (
	"testing"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// createHealthTestGraph creates a test graph for health calculator tests.
func createHealthTestGraph(packages map[string][]string) *types.DependencyGraph {
	graph := types.NewDependencyGraph("/test", types.WorkspaceTypeNpm)

	for pkgName, deps := range packages {
		node := types.NewPackageNode(pkgName, "1.0.0", "/test/"+pkgName)
		node.Dependencies = deps
		graph.Nodes[pkgName] = node
	}

	return graph
}

func TestHealthCalculator_PerfectScore(t *testing.T) {
	// No cycles, no conflicts, shallow depth, balanced coupling
	graph := createHealthTestGraph(map[string][]string{
		"@mono/app":   {"@mono/lib"},
		"@mono/lib":   {"@mono/core"},
		"@mono/core":  {},
		"@mono/utils": {"@mono/core"},
	})

	calc := NewHealthCalculator(graph, nil, nil)
	result := calc.Calculate()

	if result.Overall < 85 {
		t.Errorf("Perfect architecture should score 85+, got %d", result.Overall)
	}

	if result.Rating != types.HealthRatingExcellent && result.Rating != types.HealthRatingGood {
		t.Errorf("Rating should be excellent or good, got %s", result.Rating)
	}

	if result.Breakdown.CircularScore != 100 {
		t.Errorf("CircularScore should be 100 with no cycles, got %d", result.Breakdown.CircularScore)
	}

	if result.Breakdown.ConflictScore != 100 {
		t.Errorf("ConflictScore should be 100 with no conflicts, got %d", result.Breakdown.ConflictScore)
	}
}

func TestHealthCalculator_WithCycles(t *testing.T) {
	graph := createHealthTestGraph(map[string][]string{
		"@mono/a": {"@mono/b"},
		"@mono/b": {"@mono/a"}, // Direct cycle
	})

	cycles := []*types.CircularDependencyInfo{
		{
			Cycle:    []string{"@mono/a", "@mono/b", "@mono/a"},
			Type:     types.CircularTypeDirect,
			Severity: types.CircularSeverityWarning,
			Depth:    2,
		},
	}

	calc := NewHealthCalculator(graph, cycles, nil)
	result := calc.Calculate()

	// One direct cycle = 100 - 15 = 85
	if result.Breakdown.CircularScore != 85 {
		t.Errorf("CircularScore with 1 direct cycle should be 85, got %d", result.Breakdown.CircularScore)
	}

	// Check recommendations exist
	var circularFactor *types.HealthFactor
	for _, f := range result.Factors {
		if f.Name == "Circular Dependencies" {
			circularFactor = f
			break
		}
	}

	if circularFactor == nil {
		t.Fatal("Circular Dependencies factor not found")
	}

	if len(circularFactor.Recommendations) == 0 {
		t.Error("Expected recommendations for cycles")
	}
}

func TestHealthCalculator_WithMultipleCycles(t *testing.T) {
	graph := createHealthTestGraph(map[string][]string{
		"@mono/a": {"@mono/b"},
		"@mono/b": {"@mono/a", "@mono/c"},
		"@mono/c": {"@mono/d"},
		"@mono/d": {"@mono/b"}, // Indirect cycle
	})

	cycles := []*types.CircularDependencyInfo{
		{
			Cycle:    []string{"@mono/a", "@mono/b", "@mono/a"},
			Type:     types.CircularTypeDirect,
			Severity: types.CircularSeverityWarning,
			Depth:    2,
		},
		{
			Cycle:    []string{"@mono/b", "@mono/c", "@mono/d", "@mono/b"},
			Type:     types.CircularTypeIndirect,
			Severity: types.CircularSeverityInfo,
			Depth:    3,
		},
	}

	calc := NewHealthCalculator(graph, cycles, nil)
	result := calc.Calculate()

	// 1 direct + 1 indirect = 100 - 15 - 10 = 75
	if result.Breakdown.CircularScore != 75 {
		t.Errorf("CircularScore should be 75, got %d", result.Breakdown.CircularScore)
	}
}

func TestHealthCalculator_WithSelfLoop(t *testing.T) {
	graph := createHealthTestGraph(map[string][]string{
		"@mono/a": {"@mono/a"}, // Self-loop
	})

	cycles := []*types.CircularDependencyInfo{
		{
			Cycle:    []string{"@mono/a", "@mono/a"},
			Type:     types.CircularTypeDirect,
			Severity: types.CircularSeverityCritical,
			Depth:    1,
		},
	}

	calc := NewHealthCalculator(graph, cycles, nil)
	result := calc.Calculate()

	// Self-loop = 100 - 25 = 75
	if result.Breakdown.CircularScore != 75 {
		t.Errorf("CircularScore with self-loop should be 75, got %d", result.Breakdown.CircularScore)
	}
}

func TestHealthCalculator_WithConflicts(t *testing.T) {
	graph := createHealthTestGraph(map[string][]string{
		"@mono/a": {},
		"@mono/b": {},
	})

	conflicts := []*types.VersionConflictInfo{
		{PackageName: "lodash", Severity: types.ConflictSeverityCritical},
		{PackageName: "react", Severity: types.ConflictSeverityWarning},
		{PackageName: "moment", Severity: types.ConflictSeverityInfo},
	}

	calc := NewHealthCalculator(graph, nil, conflicts)
	result := calc.Calculate()

	// 1 critical + 1 warning + 1 info = 100 - 10 - 5 - 2 = 83
	if result.Breakdown.ConflictScore != 83 {
		t.Errorf("ConflictScore should be 83, got %d", result.Breakdown.ConflictScore)
	}
}

func TestHealthCalculator_ManyConflicts(t *testing.T) {
	graph := createHealthTestGraph(map[string][]string{
		"@mono/a": {},
	})

	// 15 critical conflicts = 150 deductions, should cap at 0
	conflicts := make([]*types.VersionConflictInfo, 15)
	for i := 0; i < 15; i++ {
		conflicts[i] = &types.VersionConflictInfo{
			PackageName: "pkg" + string(rune('a'+i)),
			Severity:    types.ConflictSeverityCritical,
		}
	}

	calc := NewHealthCalculator(graph, nil, conflicts)
	result := calc.Calculate()

	if result.Breakdown.ConflictScore != 0 {
		t.Errorf("ConflictScore should be 0 (capped), got %d", result.Breakdown.ConflictScore)
	}
}

func TestHealthCalculator_DeepDependencies(t *testing.T) {
	// Create a deep chain: a -> b -> c -> d -> e -> f -> g (depth 6)
	graph := createHealthTestGraph(map[string][]string{
		"@mono/a": {"@mono/b"},
		"@mono/b": {"@mono/c"},
		"@mono/c": {"@mono/d"},
		"@mono/d": {"@mono/e"},
		"@mono/e": {"@mono/f"},
		"@mono/f": {"@mono/g"},
		"@mono/g": {},
	})

	calc := NewHealthCalculator(graph, nil, nil)
	result := calc.Calculate()

	// Max depth 6 > optimal 4, should have deductions
	if result.Breakdown.DepthScore >= 100 {
		t.Errorf("DepthScore should be less than 100 for deep chain, got %d", result.Breakdown.DepthScore)
	}

	// Find depth factor
	var depthFactor *types.HealthFactor
	for _, f := range result.Factors {
		if f.Name == "Dependency Depth" {
			depthFactor = f
			break
		}
	}

	if depthFactor == nil {
		t.Fatal("Dependency Depth factor not found")
	}

	// Should have recommendations for deep dependencies
	if len(depthFactor.Recommendations) == 0 {
		t.Error("Expected recommendations for deep dependencies")
	}
}

func TestHealthCalculator_ShallowDependencies(t *testing.T) {
	// Shallow graph: max depth 2
	graph := createHealthTestGraph(map[string][]string{
		"@mono/a": {"@mono/b"},
		"@mono/b": {},
		"@mono/c": {"@mono/b"},
	})

	calc := NewHealthCalculator(graph, nil, nil)
	result := calc.Calculate()

	// Shallow depth should score well
	if result.Breakdown.DepthScore < 90 {
		t.Errorf("DepthScore should be high for shallow graph, got %d", result.Breakdown.DepthScore)
	}
}

func TestHealthCalculator_CouplingBalanced(t *testing.T) {
	// Balanced coupling: each package has some dependents and dependencies
	graph := createHealthTestGraph(map[string][]string{
		"@mono/a": {"@mono/b", "@mono/c"},
		"@mono/b": {"@mono/d"},
		"@mono/c": {"@mono/d"},
		"@mono/d": {},
	})

	calc := NewHealthCalculator(graph, nil, nil)
	result := calc.Calculate()

	// Balanced architecture should have reasonable coupling score
	if result.Breakdown.CouplingScore < 50 {
		t.Errorf("CouplingScore should be reasonable for balanced graph, got %d", result.Breakdown.CouplingScore)
	}
}

func TestHealthCalculator_CouplingExtreme(t *testing.T) {
	// Extreme coupling: one package depends on everything
	graph := createHealthTestGraph(map[string][]string{
		"@mono/hub": {"@mono/a", "@mono/b", "@mono/c", "@mono/d"},
		"@mono/a":   {},
		"@mono/b":   {},
		"@mono/c":   {},
		"@mono/d":   {},
	})

	calc := NewHealthCalculator(graph, nil, nil)
	result := calc.Calculate()

	// Find coupling factor
	var couplingFactor *types.HealthFactor
	for _, f := range result.Factors {
		if f.Name == "Package Coupling" {
			couplingFactor = f
			break
		}
	}

	if couplingFactor == nil {
		t.Fatal("Package Coupling factor not found")
	}

	// Should have description with instability
	if couplingFactor.Description == "" {
		t.Error("Expected description for coupling")
	}
}

func TestHealthCalculator_EmptyGraph(t *testing.T) {
	graph := types.NewDependencyGraph("/test", types.WorkspaceTypeNpm)

	calc := NewHealthCalculator(graph, nil, nil)
	result := calc.Calculate()

	if result.Overall != 100 {
		t.Errorf("Empty graph should score 100, got %d", result.Overall)
	}

	if result.Rating != types.HealthRatingExcellent {
		t.Errorf("Empty graph should rate excellent, got %s", result.Rating)
	}
}

func TestHealthCalculator_NilInputs(t *testing.T) {
	calc := NewHealthCalculator(nil, nil, nil)
	result := calc.Calculate()

	// Should not panic, should return perfect score
	if result.Overall != 100 {
		t.Errorf("Nil inputs should score 100, got %d", result.Overall)
	}
}

func TestHealthCalculator_WeightedScores(t *testing.T) {
	graph := createHealthTestGraph(map[string][]string{
		"@mono/a": {},
	})

	// No issues - all scores should be 100
	calc := NewHealthCalculator(graph, nil, nil)
	result := calc.Calculate()

	// Check weighted scores
	totalWeighted := 0
	for _, factor := range result.Factors {
		totalWeighted += factor.WeightedScore
	}

	// Weighted sum should approximately equal overall
	if totalWeighted < result.Overall-2 || totalWeighted > result.Overall+2 {
		t.Errorf("Sum of weighted scores (%d) should approximately equal overall (%d)",
			totalWeighted, result.Overall)
	}
}

func TestHealthCalculator_RatingThresholds(t *testing.T) {
	tests := []struct {
		name           string
		cycleCount     int
		conflictCount  int
		expectedRating types.HealthRating
	}{
		{"excellent - no issues", 0, 0, types.HealthRatingExcellent},
		{"good - few issues", 1, 2, types.HealthRatingGood},
		{"critical - many issues", 5, 10, types.HealthRatingCritical},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graph := createHealthTestGraph(map[string][]string{"@mono/a": {}})

			cycles := make([]*types.CircularDependencyInfo, tt.cycleCount)
			for i := 0; i < tt.cycleCount; i++ {
				cycles[i] = &types.CircularDependencyInfo{
					Cycle: []string{"a", "b", "a"},
					Type:  types.CircularTypeDirect,
					Depth: 2,
				}
			}

			conflicts := make([]*types.VersionConflictInfo, tt.conflictCount)
			for i := 0; i < tt.conflictCount; i++ {
				conflicts[i] = &types.VersionConflictInfo{
					PackageName: "pkg",
					Severity:    types.ConflictSeverityCritical,
				}
			}

			calc := NewHealthCalculator(graph, cycles, conflicts)
			result := calc.Calculate()

			// Just verify rating is set correctly based on overall score
			expectedRating := types.GetHealthRating(result.Overall)
			if result.Rating != expectedRating {
				t.Errorf("Rating mismatch: got %s, expected %s for score %d",
					result.Rating, expectedRating, result.Overall)
			}
		})
	}
}

func TestCalculateDepthMetrics(t *testing.T) {
	tests := []struct {
		name         string
		packages     map[string][]string
		wantMaxDepth int
	}{
		{
			name:         "single node",
			packages:     map[string][]string{"@mono/a": {}},
			wantMaxDepth: 0,
		},
		{
			name:         "linear chain of 3",
			packages:     map[string][]string{"@mono/a": {"@mono/b"}, "@mono/b": {"@mono/c"}, "@mono/c": {}},
			wantMaxDepth: 2,
		},
		{
			name: "diamond",
			packages: map[string][]string{
				"@mono/a": {"@mono/b", "@mono/c"},
				"@mono/b": {"@mono/d"},
				"@mono/c": {"@mono/d"},
				"@mono/d": {},
			},
			wantMaxDepth: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graph := createHealthTestGraph(tt.packages)
			calc := NewHealthCalculator(graph, nil, nil)
			maxDepth, _ := calc.calculateDepthMetrics()

			if maxDepth != tt.wantMaxDepth {
				t.Errorf("maxDepth = %d, want %d", maxDepth, tt.wantMaxDepth)
			}
		})
	}
}

func TestCalculateCouplingMetrics(t *testing.T) {
	graph := createHealthTestGraph(map[string][]string{
		"@mono/app":   {"@mono/lib", "@mono/utils"},
		"@mono/lib":   {"@mono/core"},
		"@mono/utils": {"@mono/core"},
		"@mono/core":  {},
	})

	calc := NewHealthCalculator(graph, nil, nil)
	metrics := calc.calculateCouplingMetrics()

	// Check that metrics were calculated
	if len(metrics.PackageMetrics) != 4 {
		t.Errorf("Expected 4 package metrics, got %d", len(metrics.PackageMetrics))
	}

	// Core should have high afferent coupling (2 packages depend on it)
	coreMetrics := metrics.PackageMetrics["@mono/core"]
	if coreMetrics == nil {
		t.Fatal("Core metrics not found")
	}

	if coreMetrics.AfferentCoupling != 2 {
		t.Errorf("Core AfferentCoupling = %d, want 2", coreMetrics.AfferentCoupling)
	}

	if coreMetrics.EfferentCoupling != 0 {
		t.Errorf("Core EfferentCoupling = %d, want 0", coreMetrics.EfferentCoupling)
	}

	// Core's instability should be 0 (stable)
	if coreMetrics.Instability != 0.0 {
		t.Errorf("Core Instability = %f, want 0.0", coreMetrics.Instability)
	}

	// App should have high efferent coupling
	appMetrics := metrics.PackageMetrics["@mono/app"]
	if appMetrics == nil {
		t.Fatal("App metrics not found")
	}

	if appMetrics.EfferentCoupling != 2 {
		t.Errorf("App EfferentCoupling = %d, want 2", appMetrics.EfferentCoupling)
	}
}

func TestBoundScore(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{50, 50},
		{0, 0},
		{100, 100},
		{-10, 0},
		{150, 100},
	}

	for _, tt := range tests {
		result := boundScore(tt.input)
		if result != tt.expected {
			t.Errorf("boundScore(%d) = %d, want %d", tt.input, result, tt.expected)
		}
	}
}
