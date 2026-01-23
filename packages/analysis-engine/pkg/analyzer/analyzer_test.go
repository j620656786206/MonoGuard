package analyzer

import (
	"testing"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// TestNewAnalyzer verifies analyzer creation.
func TestNewAnalyzer(t *testing.T) {
	a := NewAnalyzer()
	if a == nil {
		t.Fatal("NewAnalyzer returned nil")
	}
	if a.graphBuilder == nil {
		t.Fatal("Analyzer graphBuilder is nil")
	}
}

// TestAnalyzeEmptyWorkspace verifies analysis of empty workspace.
func TestAnalyzeEmptyWorkspace(t *testing.T) {
	a := NewAnalyzer()
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages:      map[string]*types.PackageInfo{},
	}

	result, err := a.Analyze(workspace)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	if result.Packages != 0 {
		t.Errorf("Packages = %d, want 0", result.Packages)
	}
	if result.HealthScore != 100 {
		t.Errorf("HealthScore = %d, want 100 (placeholder)", result.HealthScore)
	}
	if result.Graph == nil {
		t.Fatal("Graph is nil")
	}
	if len(result.Graph.Nodes) != 0 {
		t.Errorf("Graph.Nodes = %d, want 0", len(result.Graph.Nodes))
	}
	if len(result.Graph.Edges) != 0 {
		t.Errorf("Graph.Edges = %d, want 0", len(result.Graph.Edges))
	}
}

// TestAnalyzeSimpleWorkspace verifies analysis with real packages.
func TestAnalyzeSimpleWorkspace(t *testing.T) {
	a := NewAnalyzer()
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/app": {
				Name:    "@mono/app",
				Version: "1.0.0",
				Path:    "apps/web",
				Dependencies: map[string]string{
					"@mono/ui": "^1.0.0",
				},
				DevDependencies: map[string]string{
					"@mono/types": "^1.0.0",
				},
				PeerDependencies: map[string]string{},
			},
			"@mono/ui": {
				Name:    "@mono/ui",
				Version: "1.0.0",
				Path:    "packages/ui",
				Dependencies: map[string]string{
					"react": "^18.0.0", // External
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/types": {
				Name:             "@mono/types",
				Version:          "1.0.0",
				Path:             "packages/types",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	result, err := a.Analyze(workspace)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// Verify package count
	if result.Packages != 3 {
		t.Errorf("Packages = %d, want 3", result.Packages)
	}

	// Verify graph exists
	if result.Graph == nil {
		t.Fatal("Graph is nil")
	}

	// Verify nodes
	if len(result.Graph.Nodes) != 3 {
		t.Errorf("Graph.Nodes = %d, want 3", len(result.Graph.Nodes))
	}

	// Verify edges (2 internal: app->ui, app->types)
	if len(result.Graph.Edges) != 2 {
		t.Errorf("Graph.Edges = %d, want 2", len(result.Graph.Edges))
	}

	// Verify workspace metadata preserved
	if result.Graph.RootPath != "/workspace" {
		t.Errorf("Graph.RootPath = %q, want %q", result.Graph.RootPath, "/workspace")
	}
	if result.Graph.WorkspaceType != types.WorkspaceTypePnpm {
		t.Errorf("Graph.WorkspaceType = %q, want %q", result.Graph.WorkspaceType, types.WorkspaceTypePnpm)
	}
}

// TestAnalyzeResultHealthScore verifies health score calculation (Story 2.5).
func TestAnalyzeResultHealthScore(t *testing.T) {
	a := NewAnalyzer()
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypeNpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/pkg": {
				Name:             "@mono/pkg",
				Version:          "1.0.0",
				Path:             "packages/pkg",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	result, err := a.Analyze(workspace)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// Perfect architecture should score high
	if result.HealthScore < 85 {
		t.Errorf("HealthScore = %d, want >= 85 for perfect architecture", result.HealthScore)
	}

	// Verify HealthScoreDetails is populated
	if result.HealthScoreDetails == nil {
		t.Fatal("HealthScoreDetails is nil")
	}

	if result.HealthScoreDetails.Overall != result.HealthScore {
		t.Errorf("HealthScoreDetails.Overall (%d) != HealthScore (%d)",
			result.HealthScoreDetails.Overall, result.HealthScore)
	}

	if result.HealthScoreDetails.Rating == "" {
		t.Error("HealthScoreDetails.Rating is empty")
	}

	if result.HealthScoreDetails.Breakdown == nil {
		t.Fatal("HealthScoreDetails.Breakdown is nil")
	}

	if len(result.HealthScoreDetails.Factors) != 4 {
		t.Errorf("HealthScoreDetails.Factors count = %d, want 4", len(result.HealthScoreDetails.Factors))
	}
}

// TestAnalyzeResultHealthScoreWithIssues verifies health score decreases with issues.
func TestAnalyzeResultHealthScoreWithIssues(t *testing.T) {
	a := NewAnalyzer()
	// Create workspace with version conflicts
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypeNpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/app": {
				Name:    "@mono/app",
				Version: "1.0.0",
				Path:    "apps/web",
				Dependencies: map[string]string{
					"typescript": "^5.0.0",
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/lib": {
				Name:    "@mono/lib",
				Version: "1.0.0",
				Path:    "packages/lib",
				Dependencies: map[string]string{
					"typescript": "^4.0.0", // Major version conflict
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	result, err := a.Analyze(workspace)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// Score should be less than 100 due to version conflict
	if result.HealthScore >= 100 {
		t.Errorf("HealthScore = %d, want < 100 with version conflicts", result.HealthScore)
	}

	// Verify conflict score is reflected in breakdown
	if result.HealthScoreDetails.Breakdown.ConflictScore >= 100 {
		t.Errorf("ConflictScore = %d, want < 100 with conflicts",
			result.HealthScoreDetails.Breakdown.ConflictScore)
	}
}

// TestAnalyzeResultGraphIntegrity verifies graph structure matches workspace.
func TestAnalyzeResultGraphIntegrity(t *testing.T) {
	a := NewAnalyzer()
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypeYarn,
		Packages: map[string]*types.PackageInfo{
			"@mono/a": {
				Name:    "@mono/a",
				Version: "1.0.0",
				Path:    "packages/a",
				Dependencies: map[string]string{
					"@mono/b": "^1.0.0",
					"@mono/c": "^1.0.0",
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/b": {
				Name:    "@mono/b",
				Version: "1.0.0",
				Path:    "packages/b",
				Dependencies: map[string]string{
					"@mono/c": "^1.0.0",
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/c": {
				Name:             "@mono/c",
				Version:          "1.0.0",
				Path:             "packages/c",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	result, err := a.Analyze(workspace)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// Verify all nodes present
	graph := result.Graph
	for name := range workspace.Packages {
		if graph.Nodes[name] == nil {
			t.Errorf("Missing node: %s", name)
		}
	}

	// Verify edges: a->b, a->c, b->c (3 edges)
	if len(graph.Edges) != 3 {
		t.Errorf("Graph.Edges = %d, want 3", len(graph.Edges))
	}

	// Verify edge details
	edgeMap := make(map[string]bool)
	for _, edge := range graph.Edges {
		key := edge.From + "->" + edge.To
		edgeMap[key] = true
	}

	expectedEdges := []string{
		"@mono/a->@mono/b",
		"@mono/a->@mono/c",
		"@mono/b->@mono/c",
	}

	for _, expected := range expectedEdges {
		if !edgeMap[expected] {
			t.Errorf("Missing edge: %s", expected)
		}
	}
}

// TestAnalyzeDetectsVersionConflicts verifies version conflict detection (Story 2.4).
func TestAnalyzeDetectsVersionConflicts(t *testing.T) {
	a := NewAnalyzer()
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypeNpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/app": {
				Name:    "@mono/app",
				Version: "1.0.0",
				Path:    "apps/web",
				Dependencies: map[string]string{
					"lodash":     "^4.17.21",
					"typescript": "^5.0.0",
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/lib": {
				Name:    "@mono/lib",
				Version: "1.0.0",
				Path:    "packages/lib",
				Dependencies: map[string]string{
					"lodash":     "^4.17.19", // Patch version conflict
					"typescript": "^4.9.0",   // Major version conflict
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	result, err := a.Analyze(workspace)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// Verify version conflicts detected
	if result.VersionConflicts == nil {
		t.Fatal("VersionConflicts is nil")
	}

	if len(result.VersionConflicts) != 2 {
		t.Errorf("VersionConflicts count = %d, want 2 (lodash and typescript)", len(result.VersionConflicts))
	}

	// Find specific conflicts
	var lodashConflict, tsConflict *types.VersionConflictInfo
	for _, c := range result.VersionConflicts {
		if c.PackageName == "lodash" {
			lodashConflict = c
		}
		if c.PackageName == "typescript" {
			tsConflict = c
		}
	}

	// Verify lodash conflict (patch difference = info severity)
	if lodashConflict == nil {
		t.Fatal("lodash conflict not detected")
	}
	if lodashConflict.Severity != types.ConflictSeverityInfo {
		t.Errorf("lodash Severity = %s, want info", lodashConflict.Severity)
	}

	// Verify typescript conflict (major difference = critical severity)
	if tsConflict == nil {
		t.Fatal("typescript conflict not detected")
	}
	if tsConflict.Severity != types.ConflictSeverityCritical {
		t.Errorf("typescript Severity = %s, want critical", tsConflict.Severity)
	}
}

// TestAnalyzeNoVersionConflicts verifies no conflicts when versions match.
func TestAnalyzeNoVersionConflicts(t *testing.T) {
	a := NewAnalyzer()
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypeNpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/app": {
				Name:    "@mono/app",
				Version: "1.0.0",
				Path:    "apps/web",
				Dependencies: map[string]string{
					"lodash": "^4.17.21",
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/lib": {
				Name:    "@mono/lib",
				Version: "1.0.0",
				Path:    "packages/lib",
				Dependencies: map[string]string{
					"lodash": "^4.17.21", // Same version
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	result, err := a.Analyze(workspace)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// No conflicts when all versions match
	if len(result.VersionConflicts) != 0 {
		t.Errorf("VersionConflicts count = %d, want 0", len(result.VersionConflicts))
	}
}

// ========================================
// Story 2.6: Exclusion Pattern Tests
// ========================================

// TestNewAnalyzerWithConfig verifies analyzer creation with config.
func TestNewAnalyzerWithConfig(t *testing.T) {
	config := &types.AnalysisConfig{
		Exclude: []string{"packages/legacy", "packages/deprecated-*"},
	}

	a, err := NewAnalyzerWithConfig(config)
	if err != nil {
		t.Fatalf("NewAnalyzerWithConfig failed: %v", err)
	}
	if a == nil {
		t.Fatal("Analyzer is nil")
	}
	if a.config != config {
		t.Error("Analyzer config not set correctly")
	}
}

// TestNewAnalyzerWithConfigNil verifies nil config creates default analyzer.
func TestNewAnalyzerWithConfigNil(t *testing.T) {
	a, err := NewAnalyzerWithConfig(nil)
	if err != nil {
		t.Fatalf("NewAnalyzerWithConfig(nil) failed: %v", err)
	}
	if a == nil {
		t.Fatal("Analyzer is nil")
	}
}

// TestNewAnalyzerWithConfigInvalidRegex verifies error on invalid regex.
func TestNewAnalyzerWithConfigInvalidRegex(t *testing.T) {
	config := &types.AnalysisConfig{
		Exclude: []string{"regex:[invalid"},
	}

	_, err := NewAnalyzerWithConfig(config)
	if err == nil {
		t.Error("Expected error for invalid regex, got nil")
	}
}

// TestAnalyzeWithExclusions verifies exclusion patterns are applied.
func TestAnalyzeWithExclusions(t *testing.T) {
	config := &types.AnalysisConfig{
		Exclude: []string{"@mono/legacy"},
	}

	a, err := NewAnalyzerWithConfig(config)
	if err != nil {
		t.Fatalf("NewAnalyzerWithConfig failed: %v", err)
	}

	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/app": {
				Name:    "@mono/app",
				Version: "1.0.0",
				Path:    "apps/web",
				Dependencies: map[string]string{
					"@mono/legacy": "^1.0.0",
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/legacy": {
				Name:             "@mono/legacy",
				Version:          "1.0.0",
				Path:             "packages/legacy",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/core": {
				Name:             "@mono/core",
				Version:          "1.0.0",
				Path:             "packages/core",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	result, err := a.Analyze(workspace)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// Verify excluded count
	if result.ExcludedPackages != 1 {
		t.Errorf("ExcludedPackages = %d, want 1", result.ExcludedPackages)
	}

	// Verify non-excluded package count
	if result.Packages != 2 {
		t.Errorf("Packages = %d, want 2 (excluding legacy)", result.Packages)
	}

	// Verify excluded package is marked in graph
	legacyNode := result.Graph.Nodes["@mono/legacy"]
	if legacyNode == nil {
		t.Fatal("Legacy node missing from graph")
	}
	if !legacyNode.Excluded {
		t.Error("Legacy node should be marked as excluded")
	}

	// Verify non-excluded packages are not marked
	appNode := result.Graph.Nodes["@mono/app"]
	if appNode == nil {
		t.Fatal("App node missing from graph")
	}
	if appNode.Excluded {
		t.Error("App node should not be marked as excluded")
	}
}

// TestAnalyzeExcludedFromCycles verifies excluded packages don't affect cycle detection.
func TestAnalyzeExcludedFromCycles(t *testing.T) {
	config := &types.AnalysisConfig{
		Exclude: []string{"@mono/legacy"},
	}

	a, err := NewAnalyzerWithConfig(config)
	if err != nil {
		t.Fatalf("NewAnalyzerWithConfig failed: %v", err)
	}

	// Create workspace where legacy creates a cycle, but it should be excluded
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/app": {
				Name:    "@mono/app",
				Version: "1.0.0",
				Path:    "apps/web",
				Dependencies: map[string]string{
					"@mono/legacy": "^1.0.0",
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/legacy": {
				Name:    "@mono/legacy",
				Version: "1.0.0",
				Path:    "packages/legacy",
				Dependencies: map[string]string{
					"@mono/app": "^1.0.0", // Creates cycle with app
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	result, err := a.Analyze(workspace)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// No cycles should be detected because legacy is excluded
	if len(result.CircularDependencies) != 0 {
		t.Errorf("CircularDependencies = %d, want 0 (legacy excluded)", len(result.CircularDependencies))
	}
}

// TestAnalyzeExcludedFromConflicts verifies excluded packages don't affect conflict detection.
func TestAnalyzeExcludedFromConflicts(t *testing.T) {
	config := &types.AnalysisConfig{
		Exclude: []string{"@mono/legacy"},
	}

	a, err := NewAnalyzerWithConfig(config)
	if err != nil {
		t.Fatalf("NewAnalyzerWithConfig failed: %v", err)
	}

	// Create workspace where legacy has conflicting version
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/app": {
				Name:    "@mono/app",
				Version: "1.0.0",
				Path:    "apps/web",
				Dependencies: map[string]string{
					"typescript": "^5.0.0",
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/legacy": {
				Name:    "@mono/legacy",
				Version: "1.0.0",
				Path:    "packages/legacy",
				Dependencies: map[string]string{
					"typescript": "^4.0.0", // Conflicts with app
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	result, err := a.Analyze(workspace)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// No conflicts should be detected because legacy is excluded
	if len(result.VersionConflicts) != 0 {
		t.Errorf("VersionConflicts = %d, want 0 (legacy excluded)", len(result.VersionConflicts))
	}
}

// TestFilterExcludedPackages verifies filtering logic.
func TestFilterExcludedPackages(t *testing.T) {
	// Create a graph with some excluded packages
	graph := types.NewDependencyGraph("/workspace", types.WorkspaceTypePnpm)

	// Add nodes
	appNode := types.NewPackageNode("@mono/app", "1.0.0", "apps/web")
	appNode.Dependencies = []string{"@mono/legacy", "@mono/core"}
	graph.Nodes["@mono/app"] = appNode

	legacyNode := types.NewPackageNode("@mono/legacy", "1.0.0", "packages/legacy")
	legacyNode.Excluded = true
	graph.Nodes["@mono/legacy"] = legacyNode

	coreNode := types.NewPackageNode("@mono/core", "1.0.0", "packages/core")
	graph.Nodes["@mono/core"] = coreNode

	// Add edges
	graph.Edges = []*types.DependencyEdge{
		{From: "@mono/app", To: "@mono/legacy", Type: types.DependencyTypeProduction},
		{From: "@mono/app", To: "@mono/core", Type: types.DependencyTypeProduction},
	}

	// Filter
	filtered := filterExcludedPackages(graph)

	// Verify excluded node is removed
	if _, ok := filtered.Nodes["@mono/legacy"]; ok {
		t.Error("Excluded node should be removed from filtered graph")
	}

	// Verify non-excluded nodes are present
	if _, ok := filtered.Nodes["@mono/app"]; !ok {
		t.Error("App node should be in filtered graph")
	}
	if _, ok := filtered.Nodes["@mono/core"]; !ok {
		t.Error("Core node should be in filtered graph")
	}

	// Verify edges to excluded packages are removed
	if len(filtered.Edges) != 1 {
		t.Errorf("Filtered edges count = %d, want 1", len(filtered.Edges))
	}

	// Verify dependency list is filtered
	filteredApp := filtered.Nodes["@mono/app"]
	if len(filteredApp.Dependencies) != 1 {
		t.Errorf("Filtered app dependencies = %d, want 1", len(filteredApp.Dependencies))
	}
	if filteredApp.Dependencies[0] != "@mono/core" {
		t.Errorf("Filtered app dependency = %s, want @mono/core", filteredApp.Dependencies[0])
	}
}

// TestAnalyzeWithGlobExclusion verifies glob pattern exclusion.
func TestAnalyzeWithGlobExclusion(t *testing.T) {
	config := &types.AnalysisConfig{
		Exclude: []string{"@mono/deprecated-*"},
	}

	a, err := NewAnalyzerWithConfig(config)
	if err != nil {
		t.Fatalf("NewAnalyzerWithConfig failed: %v", err)
	}

	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/app": {
				Name:             "@mono/app",
				Version:          "1.0.0",
				Path:             "apps/web",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/deprecated-utils": {
				Name:             "@mono/deprecated-utils",
				Version:          "1.0.0",
				Path:             "packages/deprecated-utils",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/deprecated-api": {
				Name:             "@mono/deprecated-api",
				Version:          "1.0.0",
				Path:             "packages/deprecated-api",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	result, err := a.Analyze(workspace)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// 2 packages should be excluded (deprecated-*)
	if result.ExcludedPackages != 2 {
		t.Errorf("ExcludedPackages = %d, want 2", result.ExcludedPackages)
	}

	// Only app should be counted
	if result.Packages != 1 {
		t.Errorf("Packages = %d, want 1", result.Packages)
	}
}

// TestAnalyzeWithRegexExclusion verifies regex pattern exclusion.
func TestAnalyzeWithRegexExclusion(t *testing.T) {
	config := &types.AnalysisConfig{
		Exclude: []string{"regex:.*-test$"},
	}

	a, err := NewAnalyzerWithConfig(config)
	if err != nil {
		t.Fatalf("NewAnalyzerWithConfig failed: %v", err)
	}

	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/app": {
				Name:             "@mono/app",
				Version:          "1.0.0",
				Path:             "apps/web",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/core-test": {
				Name:             "@mono/core-test",
				Version:          "1.0.0",
				Path:             "packages/core-test",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	result, err := a.Analyze(workspace)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// core-test should be excluded
	if result.ExcludedPackages != 1 {
		t.Errorf("ExcludedPackages = %d, want 1", result.ExcludedPackages)
	}

	// Verify the test package is marked excluded
	testNode := result.Graph.Nodes["@mono/core-test"]
	if testNode == nil || !testNode.Excluded {
		t.Error("core-test should be marked as excluded")
	}
}

// ========================================
// Story 2.6: filterExcludedPackages Edge Cases
// ========================================

// TestFilterExcludedPackages_AllExcluded verifies behavior when all packages are excluded.
func TestFilterExcludedPackages_AllExcluded(t *testing.T) {
	graph := types.NewDependencyGraph("/workspace", types.WorkspaceTypePnpm)

	// Add nodes - all excluded
	node1 := types.NewPackageNode("@mono/legacy-a", "1.0.0", "packages/legacy-a")
	node1.Excluded = true
	graph.Nodes["@mono/legacy-a"] = node1

	node2 := types.NewPackageNode("@mono/legacy-b", "1.0.0", "packages/legacy-b")
	node2.Excluded = true
	graph.Nodes["@mono/legacy-b"] = node2

	// Add edge between excluded packages
	graph.Edges = []*types.DependencyEdge{
		{From: "@mono/legacy-a", To: "@mono/legacy-b", Type: types.DependencyTypeProduction},
	}

	// Filter
	filtered := filterExcludedPackages(graph)

	// All nodes should be removed
	if len(filtered.Nodes) != 0 {
		t.Errorf("Filtered nodes = %d, want 0 (all excluded)", len(filtered.Nodes))
	}

	// All edges should be removed
	if len(filtered.Edges) != 0 {
		t.Errorf("Filtered edges = %d, want 0 (all excluded)", len(filtered.Edges))
	}
}

// TestFilterExcludedPackages_NoneExcluded verifies behavior when no packages are excluded.
func TestFilterExcludedPackages_NoneExcluded(t *testing.T) {
	graph := types.NewDependencyGraph("/workspace", types.WorkspaceTypePnpm)

	// Add nodes - none excluded
	node1 := types.NewPackageNode("@mono/app", "1.0.0", "apps/web")
	node1.Dependencies = []string{"@mono/core"}
	graph.Nodes["@mono/app"] = node1

	node2 := types.NewPackageNode("@mono/core", "1.0.0", "packages/core")
	graph.Nodes["@mono/core"] = node2

	// Add edge
	graph.Edges = []*types.DependencyEdge{
		{From: "@mono/app", To: "@mono/core", Type: types.DependencyTypeProduction},
	}

	// Filter
	filtered := filterExcludedPackages(graph)

	// All nodes should be preserved
	if len(filtered.Nodes) != 2 {
		t.Errorf("Filtered nodes = %d, want 2", len(filtered.Nodes))
	}

	// All edges should be preserved
	if len(filtered.Edges) != 1 {
		t.Errorf("Filtered edges = %d, want 1", len(filtered.Edges))
	}

	// Dependency list should be preserved
	appNode := filtered.Nodes["@mono/app"]
	if len(appNode.Dependencies) != 1 || appNode.Dependencies[0] != "@mono/core" {
		t.Errorf("App dependencies = %v, want [@mono/core]", appNode.Dependencies)
	}
}

// TestFilterExcludedPackages_AllDependencyTypes verifies filtering across all dependency types.
func TestFilterExcludedPackages_AllDependencyTypes(t *testing.T) {
	graph := types.NewDependencyGraph("/workspace", types.WorkspaceTypePnpm)

	// Add non-excluded node with all dependency types pointing to excluded
	appNode := types.NewPackageNode("@mono/app", "1.0.0", "apps/web")
	appNode.Dependencies = []string{"@mono/prod-excluded", "@mono/prod-ok"}
	appNode.DevDependencies = []string{"@mono/dev-excluded", "@mono/dev-ok"}
	appNode.PeerDependencies = []string{"@mono/peer-excluded", "@mono/peer-ok"}
	appNode.OptionalDependencies = []string{"@mono/opt-excluded", "@mono/opt-ok"}
	graph.Nodes["@mono/app"] = appNode

	// Add excluded packages
	for _, name := range []string{"@mono/prod-excluded", "@mono/dev-excluded", "@mono/peer-excluded", "@mono/opt-excluded"} {
		node := types.NewPackageNode(name, "1.0.0", "packages/"+name)
		node.Excluded = true
		graph.Nodes[name] = node
	}

	// Add non-excluded packages
	for _, name := range []string{"@mono/prod-ok", "@mono/dev-ok", "@mono/peer-ok", "@mono/opt-ok"} {
		node := types.NewPackageNode(name, "1.0.0", "packages/"+name)
		graph.Nodes[name] = node
	}

	// Filter
	filtered := filterExcludedPackages(graph)

	// Should have 5 nodes (app + 4 non-excluded)
	if len(filtered.Nodes) != 5 {
		t.Errorf("Filtered nodes = %d, want 5", len(filtered.Nodes))
	}

	// Verify dependency lists are filtered
	filteredApp := filtered.Nodes["@mono/app"]

	if len(filteredApp.Dependencies) != 1 || filteredApp.Dependencies[0] != "@mono/prod-ok" {
		t.Errorf("Filtered Dependencies = %v, want [@mono/prod-ok]", filteredApp.Dependencies)
	}
	if len(filteredApp.DevDependencies) != 1 || filteredApp.DevDependencies[0] != "@mono/dev-ok" {
		t.Errorf("Filtered DevDependencies = %v, want [@mono/dev-ok]", filteredApp.DevDependencies)
	}
	if len(filteredApp.PeerDependencies) != 1 || filteredApp.PeerDependencies[0] != "@mono/peer-ok" {
		t.Errorf("Filtered PeerDependencies = %v, want [@mono/peer-ok]", filteredApp.PeerDependencies)
	}
	if len(filteredApp.OptionalDependencies) != 1 || filteredApp.OptionalDependencies[0] != "@mono/opt-ok" {
		t.Errorf("Filtered OptionalDependencies = %v, want [@mono/opt-ok]", filteredApp.OptionalDependencies)
	}
}

// TestFilterExcludedPackages_PreservesExternalDeps verifies external deps are preserved.
func TestFilterExcludedPackages_PreservesExternalDeps(t *testing.T) {
	graph := types.NewDependencyGraph("/workspace", types.WorkspaceTypePnpm)

	// Add node with external dependencies
	appNode := types.NewPackageNode("@mono/app", "1.0.0", "apps/web")
	appNode.ExternalDeps = map[string]string{"react": "^18.0.0", "lodash": "^4.17.21"}
	appNode.ExternalDevDeps = map[string]string{"typescript": "^5.0.0"}
	appNode.ExternalPeerDeps = map[string]string{"react-dom": "^18.0.0"}
	appNode.ExternalOptionalDeps = map[string]string{"fsevents": "^2.3.0"}
	graph.Nodes["@mono/app"] = appNode

	// Filter
	filtered := filterExcludedPackages(graph)

	// Verify external deps are preserved
	filteredApp := filtered.Nodes["@mono/app"]

	if filteredApp.ExternalDeps["react"] != "^18.0.0" {
		t.Errorf("ExternalDeps[react] = %s, want ^18.0.0", filteredApp.ExternalDeps["react"])
	}
	if filteredApp.ExternalDevDeps["typescript"] != "^5.0.0" {
		t.Errorf("ExternalDevDeps[typescript] = %s, want ^5.0.0", filteredApp.ExternalDevDeps["typescript"])
	}
	if filteredApp.ExternalPeerDeps["react-dom"] != "^18.0.0" {
		t.Errorf("ExternalPeerDeps[react-dom] = %s, want ^18.0.0", filteredApp.ExternalPeerDeps["react-dom"])
	}
	if filteredApp.ExternalOptionalDeps["fsevents"] != "^2.3.0" {
		t.Errorf("ExternalOptionalDeps[fsevents] = %s, want ^2.3.0", filteredApp.ExternalOptionalDeps["fsevents"])
	}
}

// ========================================
// Story 3.1: Root Cause Analysis Integration Tests
// ========================================

// TestAnalyzeCyclesHaveRootCause verifies cycles are enriched with root cause analysis.
func TestAnalyzeCyclesHaveRootCause(t *testing.T) {
	a := NewAnalyzer()
	// Create workspace with a simple cycle: A → B → A
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/app": {
				Name:    "@mono/app",
				Version: "1.0.0",
				Path:    "apps/web",
				Dependencies: map[string]string{
					"@mono/lib": "^1.0.0",
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/lib": {
				Name:    "@mono/lib",
				Version: "1.0.0",
				Path:    "packages/lib",
				Dependencies: map[string]string{
					"@mono/app": "^1.0.0", // Creates cycle
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	result, err := a.Analyze(workspace)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// Verify cycle detected
	if len(result.CircularDependencies) != 1 {
		t.Fatalf("CircularDependencies = %d, want 1", len(result.CircularDependencies))
	}

	cycle := result.CircularDependencies[0]

	// Verify root cause is populated (Story 3.1)
	if cycle.RootCause == nil {
		t.Fatal("RootCause should be populated for detected cycles")
	}

	// Verify root cause fields
	if cycle.RootCause.OriginatingPackage == "" {
		t.Error("RootCause.OriginatingPackage should not be empty")
	}
	if cycle.RootCause.Confidence == 0 {
		t.Error("RootCause.Confidence should be > 0")
	}
	if cycle.RootCause.Explanation == "" {
		t.Error("RootCause.Explanation should not be empty")
	}
	if len(cycle.RootCause.Chain) == 0 {
		t.Error("RootCause.Chain should not be empty")
	}
}

// TestAnalyzeCyclesRootCauseIndirectCycle verifies root cause for indirect cycles.
func TestAnalyzeCyclesRootCauseIndirectCycle(t *testing.T) {
	a := NewAnalyzer()
	// Create workspace with indirect cycle: A → B → C → A
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/app": {
				Name:    "@mono/app",
				Version: "1.0.0",
				Path:    "apps/web",
				Dependencies: map[string]string{
					"@mono/service": "^1.0.0",
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/service": {
				Name:    "@mono/service",
				Version: "1.0.0",
				Path:    "packages/service",
				Dependencies: map[string]string{
					"@mono/core": "^1.0.0",
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/core": {
				Name:    "@mono/core",
				Version: "1.0.0",
				Path:    "packages/core",
				Dependencies: map[string]string{
					"@mono/app": "^1.0.0", // Creates indirect cycle
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	result, err := a.Analyze(workspace)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// Verify cycle detected
	if len(result.CircularDependencies) != 1 {
		t.Fatalf("CircularDependencies = %d, want 1", len(result.CircularDependencies))
	}

	cycle := result.CircularDependencies[0]

	// Verify root cause is populated
	if cycle.RootCause == nil {
		t.Fatal("RootCause should be populated for indirect cycles")
	}

	// For indirect cycle with "core" package, core should NOT be the root cause
	if cycle.RootCause.OriginatingPackage == "@mono/core" {
		t.Error("RootCause should not be 'core' package (name pattern heuristic)")
	}

	// Verify chain has 3 edges for 3-node cycle
	if len(cycle.RootCause.Chain) != 3 {
		t.Errorf("RootCause.Chain length = %d, want 3", len(cycle.RootCause.Chain))
	}
}

// TestAnalyzeCyclesNoCycleNoRootCause verifies no root cause when no cycles.
func TestAnalyzeCyclesNoCycleNoRootCause(t *testing.T) {
	a := NewAnalyzer()
	// Create workspace without cycles
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/app": {
				Name:    "@mono/app",
				Version: "1.0.0",
				Path:    "apps/web",
				Dependencies: map[string]string{
					"@mono/lib": "^1.0.0",
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/lib": {
				Name:             "@mono/lib",
				Version:          "1.0.0",
				Path:             "packages/lib",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	result, err := a.Analyze(workspace)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// Verify no cycles
	if len(result.CircularDependencies) != 0 {
		t.Errorf("CircularDependencies = %d, want 0", len(result.CircularDependencies))
	}
}

// ========================================
// Story 3.2: Import Tracing Integration Tests
// ========================================

// TestAnalyzeWithSources verifies import tracing when source files are provided.
func TestAnalyzeWithSources(t *testing.T) {
	a := NewAnalyzer()
	// Create workspace with a cycle: ui → api → ui
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/ui": {
				Name:    "@mono/ui",
				Version: "1.0.0",
				Path:    "packages/ui",
				Dependencies: map[string]string{
					"@mono/api": "^1.0.0",
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/api": {
				Name:    "@mono/api",
				Version: "1.0.0",
				Path:    "packages/api",
				Dependencies: map[string]string{
					"@mono/ui": "^1.0.0", // Creates cycle
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	// Provide source files
	sourceFiles := map[string][]byte{
		"packages/ui/src/index.ts":  []byte(`import { api } from '@mono/api';`),
		"packages/api/src/index.ts": []byte(`import { ui } from '@mono/ui';`),
	}

	result, err := a.AnalyzeWithSources(workspace, sourceFiles)
	if err != nil {
		t.Fatalf("AnalyzeWithSources failed: %v", err)
	}

	// Verify cycle detected
	if len(result.CircularDependencies) != 1 {
		t.Fatalf("CircularDependencies = %d, want 1", len(result.CircularDependencies))
	}

	cycle := result.CircularDependencies[0]

	// Verify import traces are populated (Story 3.2)
	if cycle.ImportTraces == nil {
		t.Fatal("ImportTraces should not be nil when source files provided")
	}
	if len(cycle.ImportTraces) == 0 {
		t.Fatal("ImportTraces should be populated when source files provided")
	}

	// Verify traces have expected structure
	for _, trace := range cycle.ImportTraces {
		if trace.FromPackage == "" {
			t.Error("ImportTrace.FromPackage should not be empty")
		}
		if trace.ToPackage == "" {
			t.Error("ImportTrace.ToPackage should not be empty")
		}
		if trace.FilePath == "" {
			t.Error("ImportTrace.FilePath should not be empty")
		}
		if trace.LineNumber <= 0 {
			t.Error("ImportTrace.LineNumber should be positive")
		}
		if trace.Statement == "" {
			t.Error("ImportTrace.Statement should not be empty")
		}
	}
}

// TestAnalyzeWithSourcesEmptyFiles verifies graceful degradation with empty source files.
func TestAnalyzeWithSourcesEmptyFiles(t *testing.T) {
	a := NewAnalyzer()
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/ui": {
				Name:    "@mono/ui",
				Version: "1.0.0",
				Path:    "packages/ui",
				Dependencies: map[string]string{
					"@mono/api": "^1.0.0",
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/api": {
				Name:    "@mono/api",
				Version: "1.0.0",
				Path:    "packages/api",
				Dependencies: map[string]string{
					"@mono/ui": "^1.0.0",
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	// Empty source files - should still work
	result, err := a.AnalyzeWithSources(workspace, map[string][]byte{})
	if err != nil {
		t.Fatalf("AnalyzeWithSources failed: %v", err)
	}

	// Cycle should still be detected
	if len(result.CircularDependencies) != 1 {
		t.Fatalf("CircularDependencies = %d, want 1", len(result.CircularDependencies))
	}

	cycle := result.CircularDependencies[0]

	// Import traces should be empty (not nil) for graceful degradation
	if cycle.ImportTraces == nil {
		t.Error("ImportTraces should not be nil (graceful degradation)")
	}
}

// TestAnalyzeWithSourcesNilFiles verifies behavior with nil source files.
func TestAnalyzeWithSourcesNilFiles(t *testing.T) {
	a := NewAnalyzer()
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/app": {
				Name:             "@mono/app",
				Version:          "1.0.0",
				Path:             "apps/web",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	// Nil source files - should still work
	result, err := a.AnalyzeWithSources(workspace, nil)
	if err != nil {
		t.Fatalf("AnalyzeWithSources failed: %v", err)
	}

	// Basic analysis should succeed
	if result.Packages != 1 {
		t.Errorf("Packages = %d, want 1", result.Packages)
	}
}

// TestAnalyzeWithSourcesBackwardCompatibility verifies Analyze still works without sources.
func TestAnalyzeWithSourcesBackwardCompatibility(t *testing.T) {
	a := NewAnalyzer()
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/ui": {
				Name:    "@mono/ui",
				Version: "1.0.0",
				Path:    "packages/ui",
				Dependencies: map[string]string{
					"@mono/api": "^1.0.0",
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/api": {
				Name:    "@mono/api",
				Version: "1.0.0",
				Path:    "packages/api",
				Dependencies: map[string]string{
					"@mono/ui": "^1.0.0",
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	// Use original Analyze method (backward compatible)
	result, err := a.Analyze(workspace)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// Cycle should be detected
	if len(result.CircularDependencies) != 1 {
		t.Fatalf("CircularDependencies = %d, want 1", len(result.CircularDependencies))
	}

	cycle := result.CircularDependencies[0]

	// ImportTraces should not be populated when using original Analyze
	// (since no source files were provided)
	if len(cycle.ImportTraces) != 0 {
		t.Errorf("ImportTraces should be empty for Analyze() without sources")
	}

	// RootCause should still be populated (Story 3.1)
	if cycle.RootCause == nil {
		t.Error("RootCause should still be populated for backward compatibility")
	}
}

// ========================================
// Story 3.3: Fix Strategy Integration Tests
// ========================================

// TestAnalyzeCyclesHaveFixStrategies verifies cycles are enriched with fix strategies.
func TestAnalyzeCyclesHaveFixStrategies(t *testing.T) {
	a := NewAnalyzer()
	// Create workspace with a simple cycle: A → B → A
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/app": {
				Name:    "@mono/app",
				Version: "1.0.0",
				Path:    "apps/web",
				Dependencies: map[string]string{
					"@mono/lib": "^1.0.0",
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/lib": {
				Name:    "@mono/lib",
				Version: "1.0.0",
				Path:    "packages/lib",
				Dependencies: map[string]string{
					"@mono/app": "^1.0.0", // Creates cycle
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	result, err := a.Analyze(workspace)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// Verify cycle detected
	if len(result.CircularDependencies) != 1 {
		t.Fatalf("CircularDependencies = %d, want 1", len(result.CircularDependencies))
	}

	cycle := result.CircularDependencies[0]

	// Verify fix strategies are populated (Story 3.3)
	if len(cycle.FixStrategies) == 0 {
		t.Fatal("FixStrategies should be populated for detected cycles")
	}

	// Should have 3 strategies
	if len(cycle.FixStrategies) != 3 {
		t.Errorf("FixStrategies = %d, want 3", len(cycle.FixStrategies))
	}

	// Verify first strategy is marked recommended
	if !cycle.FixStrategies[0].Recommended {
		t.Error("First strategy should be marked as recommended")
	}

	// Verify strategies are sorted by suitability
	for i := 1; i < len(cycle.FixStrategies); i++ {
		if cycle.FixStrategies[i].Suitability > cycle.FixStrategies[i-1].Suitability {
			t.Errorf("Strategies not sorted: index %d has higher suitability than %d", i, i-1)
		}
	}
}

// TestAnalyzeCyclesFixStrategiesForDirectCycle verifies DI is preferred for direct cycles.
func TestAnalyzeCyclesFixStrategiesForDirectCycle(t *testing.T) {
	a := NewAnalyzer()
	// Create workspace with direct cycle: A ↔ B
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/a": {
				Name:    "@mono/a",
				Version: "1.0.0",
				Path:    "packages/a",
				Dependencies: map[string]string{
					"@mono/b": "^1.0.0",
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/b": {
				Name:    "@mono/b",
				Version: "1.0.0",
				Path:    "packages/b",
				Dependencies: map[string]string{
					"@mono/a": "^1.0.0", // Creates cycle
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	result, err := a.Analyze(workspace)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	if len(result.CircularDependencies) != 1 {
		t.Fatalf("CircularDependencies = %d, want 1", len(result.CircularDependencies))
	}

	cycle := result.CircularDependencies[0]

	// Find DI strategy
	var diStrategy *types.FixStrategy
	for i := range cycle.FixStrategies {
		if cycle.FixStrategies[i].Type == types.FixStrategyDependencyInject {
			diStrategy = &cycle.FixStrategies[i]
			break
		}
	}

	if diStrategy == nil {
		t.Fatal("Dependency Injection strategy not found")
	}

	// For direct cycles, DI should have high suitability
	if diStrategy.Suitability < 8 {
		t.Errorf("DI suitability for direct cycle = %d, want >= 8", diStrategy.Suitability)
	}
}

// TestAnalyzeCyclesNoCycleNoFixStrategies verifies no strategies when no cycles.
func TestAnalyzeCyclesNoCycleNoFixStrategies(t *testing.T) {
	a := NewAnalyzer()
	// Create workspace without cycles
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/app": {
				Name:    "@mono/app",
				Version: "1.0.0",
				Path:    "apps/web",
				Dependencies: map[string]string{
					"@mono/lib": "^1.0.0",
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/lib": {
				Name:             "@mono/lib",
				Version:          "1.0.0",
				Path:             "packages/lib",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	result, err := a.Analyze(workspace)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// Verify no cycles
	if len(result.CircularDependencies) != 0 {
		t.Errorf("CircularDependencies = %d, want 0", len(result.CircularDependencies))
	}
}

// TestAnalyzeWithSourcesIncludesFixStrategies verifies AnalyzeWithSources also generates fix strategies.
func TestAnalyzeWithSourcesIncludesFixStrategies(t *testing.T) {
	a := NewAnalyzer()
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/ui": {
				Name:    "@mono/ui",
				Version: "1.0.0",
				Path:    "packages/ui",
				Dependencies: map[string]string{
					"@mono/api": "^1.0.0",
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/api": {
				Name:    "@mono/api",
				Version: "1.0.0",
				Path:    "packages/api",
				Dependencies: map[string]string{
					"@mono/ui": "^1.0.0", // Creates cycle
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	// Provide source files
	sourceFiles := map[string][]byte{
		"packages/ui/src/index.ts":  []byte(`import { api } from '@mono/api';`),
		"packages/api/src/index.ts": []byte(`import { ui } from '@mono/ui';`),
	}

	result, err := a.AnalyzeWithSources(workspace, sourceFiles)
	if err != nil {
		t.Fatalf("AnalyzeWithSources failed: %v", err)
	}

	if len(result.CircularDependencies) != 1 {
		t.Fatalf("CircularDependencies = %d, want 1", len(result.CircularDependencies))
	}

	cycle := result.CircularDependencies[0]

	// FixStrategies should be populated (Story 3.3)
	if len(cycle.FixStrategies) == 0 {
		t.Error("FixStrategies should be populated for AnalyzeWithSources")
	}

	// ImportTraces should also be populated (Story 3.2)
	if len(cycle.ImportTraces) == 0 {
		t.Error("ImportTraces should be populated when source files provided")
	}

	// RootCause should also be populated (Story 3.1)
	if cycle.RootCause == nil {
		t.Error("RootCause should be populated")
	}
}
