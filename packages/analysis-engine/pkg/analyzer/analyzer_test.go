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

// TestAnalyzeResultHealthScorePlaceholder verifies health score is placeholder.
func TestAnalyzeResultHealthScorePlaceholder(t *testing.T) {
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

	// Health score should be 100 (placeholder until Story 2.5)
	if result.HealthScore != 100 {
		t.Errorf("HealthScore = %d, want 100 (placeholder)", result.HealthScore)
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
