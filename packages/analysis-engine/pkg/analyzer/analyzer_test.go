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
