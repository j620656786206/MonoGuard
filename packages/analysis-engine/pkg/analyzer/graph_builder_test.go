package analyzer

import (
	"testing"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// TestNewGraphBuilder verifies graph builder creation.
func TestNewGraphBuilder(t *testing.T) {
	gb := NewGraphBuilder()
	if gb == nil {
		t.Fatal("NewGraphBuilder returned nil")
	}
}

// TestBuildEmptyWorkspace verifies graph construction with no packages.
func TestBuildEmptyWorkspace(t *testing.T) {
	gb := NewGraphBuilder()
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages:      map[string]*types.PackageInfo{},
	}

	graph, err := gb.Build(workspace)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if len(graph.Nodes) != 0 {
		t.Errorf("Expected 0 nodes, got %d", len(graph.Nodes))
	}
	if len(graph.Edges) != 0 {
		t.Errorf("Expected 0 edges, got %d", len(graph.Edges))
	}
	if graph.RootPath != "/workspace" {
		t.Errorf("RootPath = %q, want %q", graph.RootPath, "/workspace")
	}
	if graph.WorkspaceType != types.WorkspaceTypePnpm {
		t.Errorf("WorkspaceType = %q, want %q", graph.WorkspaceType, types.WorkspaceTypePnpm)
	}
}

// TestBuildSinglePackageNoDepends verifies graph with isolated package.
func TestBuildSinglePackageNoDepends(t *testing.T) {
	gb := NewGraphBuilder()
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/utils": {
				Name:             "@mono/utils",
				Version:          "1.0.0",
				Path:             "packages/utils",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	graph, err := gb.Build(workspace)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if len(graph.Nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(graph.Nodes))
	}
	if len(graph.Edges) != 0 {
		t.Errorf("Expected 0 edges, got %d", len(graph.Edges))
	}

	node := graph.Nodes["@mono/utils"]
	if node == nil {
		t.Fatal("Missing @mono/utils node")
	}
	if node.Name != "@mono/utils" {
		t.Errorf("Node name = %q, want %q", node.Name, "@mono/utils")
	}
	if node.Version != "1.0.0" {
		t.Errorf("Node version = %q, want %q", node.Version, "1.0.0")
	}
	if node.Path != "packages/utils" {
		t.Errorf("Node path = %q, want %q", node.Path, "packages/utils")
	}
}

// TestBuildInternalDependencies verifies edges created for internal workspace deps.
func TestBuildInternalDependencies(t *testing.T) {
	gb := NewGraphBuilder()
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
				Name:             "@mono/ui",
				Version:          "1.0.0",
				Path:             "packages/ui",
				Dependencies:     map[string]string{},
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

	graph, err := gb.Build(workspace)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Verify nodes
	if len(graph.Nodes) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(graph.Nodes))
	}

	// Verify edges (should have 2: app->ui (production), app->types (development))
	if len(graph.Edges) != 2 {
		t.Errorf("Expected 2 edges, got %d", len(graph.Edges))
	}

	// Find edges
	var productionEdge, devEdge *types.DependencyEdge
	for _, edge := range graph.Edges {
		if edge.From == "@mono/app" && edge.To == "@mono/ui" {
			productionEdge = edge
		}
		if edge.From == "@mono/app" && edge.To == "@mono/types" {
			devEdge = edge
		}
	}

	if productionEdge == nil {
		t.Error("Missing production edge @mono/app -> @mono/ui")
	} else {
		if productionEdge.Type != types.DependencyTypeProduction {
			t.Errorf("Production edge type = %q, want %q", productionEdge.Type, types.DependencyTypeProduction)
		}
		if productionEdge.VersionRange != "^1.0.0" {
			t.Errorf("Production edge versionRange = %q, want %q", productionEdge.VersionRange, "^1.0.0")
		}
	}

	if devEdge == nil {
		t.Error("Missing development edge @mono/app -> @mono/types")
	} else {
		if devEdge.Type != types.DependencyTypeDevelopment {
			t.Errorf("Development edge type = %q, want %q", devEdge.Type, types.DependencyTypeDevelopment)
		}
	}

	// Verify node dependencies lists
	appNode := graph.Nodes["@mono/app"]
	if len(appNode.Dependencies) != 1 || appNode.Dependencies[0] != "@mono/ui" {
		t.Errorf("App dependencies = %v, want [\"@mono/ui\"]", appNode.Dependencies)
	}
	if len(appNode.DevDependencies) != 1 || appNode.DevDependencies[0] != "@mono/types" {
		t.Errorf("App devDependencies = %v, want [\"@mono/types\"]", appNode.DevDependencies)
	}
}

// TestBuildExternalDependencies verifies external deps stored in metadata, not edges.
func TestBuildExternalDependencies(t *testing.T) {
	gb := NewGraphBuilder()
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/app": {
				Name:    "@mono/app",
				Version: "1.0.0",
				Path:    "apps/web",
				Dependencies: map[string]string{
					"react":     "^18.0.0",
					"lodash":    "^4.17.0",
					"@mono/ui":  "^1.0.0", // Internal
				},
				DevDependencies: map[string]string{
					"typescript": "^5.0.0",
					"jest":       "^29.0.0",
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
		},
	}

	graph, err := gb.Build(workspace)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Should only have 1 edge (internal: app -> ui)
	if len(graph.Edges) != 1 {
		t.Errorf("Expected 1 edge (internal only), got %d", len(graph.Edges))
	}

	// Verify external deps stored in node metadata
	appNode := graph.Nodes["@mono/app"]
	if appNode.ExternalDeps["react"] != "^18.0.0" {
		t.Errorf("App externalDeps[react] = %q, want ^18.0.0", appNode.ExternalDeps["react"])
	}
	if appNode.ExternalDeps["lodash"] != "^4.17.0" {
		t.Errorf("App externalDeps[lodash] = %q, want ^4.17.0", appNode.ExternalDeps["lodash"])
	}
	if appNode.ExternalDevDeps["typescript"] != "^5.0.0" {
		t.Errorf("App externalDevDeps[typescript] = %q, want ^5.0.0", appNode.ExternalDevDeps["typescript"])
	}

	// Internal dep should be in Dependencies list, NOT ExternalDeps
	if len(appNode.Dependencies) != 1 || appNode.Dependencies[0] != "@mono/ui" {
		t.Errorf("App dependencies = %v, want [\"@mono/ui\"]", appNode.Dependencies)
	}
	if _, exists := appNode.ExternalDeps["@mono/ui"]; exists {
		t.Error("Internal dep @mono/ui should NOT be in externalDeps")
	}
}

// TestBuildPeerDependencies verifies peer dependency edge creation.
func TestBuildPeerDependencies(t *testing.T) {
	gb := NewGraphBuilder()
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/plugin": {
				Name:             "@mono/plugin",
				Version:          "1.0.0",
				Path:             "packages/plugin",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{
					"@mono/core": "^1.0.0", // Internal peer dep
					"react":      "^18.0.0", // External peer dep
				},
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

	graph, err := gb.Build(workspace)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Should have 1 edge (internal peer: plugin -> core)
	if len(graph.Edges) != 1 {
		t.Errorf("Expected 1 edge, got %d", len(graph.Edges))
	}

	edge := graph.Edges[0]
	if edge.Type != types.DependencyTypePeer {
		t.Errorf("Edge type = %q, want %q", edge.Type, types.DependencyTypePeer)
	}
	if edge.From != "@mono/plugin" || edge.To != "@mono/core" {
		t.Errorf("Edge = %s -> %s, want @mono/plugin -> @mono/core", edge.From, edge.To)
	}

	// Verify node peerDependencies list
	pluginNode := graph.Nodes["@mono/plugin"]
	if len(pluginNode.PeerDependencies) != 1 || pluginNode.PeerDependencies[0] != "@mono/core" {
		t.Errorf("Plugin peerDependencies = %v, want [\"@mono/core\"]", pluginNode.PeerDependencies)
	}
}

// TestBuildDependencyTypeClassification verifies correct type classification.
func TestBuildDependencyTypeClassification(t *testing.T) {
	gb := NewGraphBuilder()
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/app": {
				Name:    "@mono/app",
				Version: "1.0.0",
				Path:    "apps/web",
				Dependencies: map[string]string{
					"@mono/prod": "^1.0.0",
				},
				DevDependencies: map[string]string{
					"@mono/dev": "^1.0.0",
				},
				PeerDependencies: map[string]string{
					"@mono/peer": "^1.0.0",
				},
			},
			"@mono/prod": {
				Name:             "@mono/prod",
				Version:          "1.0.0",
				Path:             "packages/prod",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/dev": {
				Name:             "@mono/dev",
				Version:          "1.0.0",
				Path:             "packages/dev",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/peer": {
				Name:             "@mono/peer",
				Version:          "1.0.0",
				Path:             "packages/peer",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	graph, err := gb.Build(workspace)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Should have 3 edges with different types
	if len(graph.Edges) != 3 {
		t.Errorf("Expected 3 edges, got %d", len(graph.Edges))
	}

	typeCount := map[types.DependencyType]int{}
	for _, edge := range graph.Edges {
		typeCount[edge.Type]++
	}

	if typeCount[types.DependencyTypeProduction] != 1 {
		t.Errorf("Expected 1 production edge, got %d", typeCount[types.DependencyTypeProduction])
	}
	if typeCount[types.DependencyTypeDevelopment] != 1 {
		t.Errorf("Expected 1 development edge, got %d", typeCount[types.DependencyTypeDevelopment])
	}
	if typeCount[types.DependencyTypePeer] != 1 {
		t.Errorf("Expected 1 peer edge, got %d", typeCount[types.DependencyTypePeer])
	}
}

// ========================================
// Edge Case Tests (Task 3)
// ========================================

// TestBuildSelfReferencingPackage verifies A->A does NOT create an edge.
func TestBuildSelfReferencingPackage(t *testing.T) {
	gb := NewGraphBuilder()
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/recursive": {
				Name:    "@mono/recursive",
				Version: "1.0.0",
				Path:    "packages/recursive",
				Dependencies: map[string]string{
					"@mono/recursive": "^1.0.0", // Self-reference - should NOT create edge
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	graph, err := gb.Build(workspace)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Should have 0 edges (self-reference skipped)
	if len(graph.Edges) != 0 {
		t.Errorf("Expected 0 edges (self-reference skipped), got %d", len(graph.Edges))
		for _, edge := range graph.Edges {
			t.Errorf("  Unexpected edge: %s -> %s", edge.From, edge.To)
		}
	}

	// Node should exist but have empty internal dependencies
	node := graph.Nodes["@mono/recursive"]
	if node == nil {
		t.Fatal("Missing @mono/recursive node")
	}
	if len(node.Dependencies) != 0 {
		t.Errorf("Expected 0 internal dependencies, got %d", len(node.Dependencies))
	}
}

// TestBuildMissingDependency verifies deps not in workspace are treated as external.
func TestBuildMissingDependency(t *testing.T) {
	gb := NewGraphBuilder()
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/app": {
				Name:    "@mono/app",
				Version: "1.0.0",
				Path:    "apps/web",
				Dependencies: map[string]string{
					"@mono/nonexistent": "^1.0.0", // Not in workspace
					"lodash":            "^4.17.0", // External npm package
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	graph, err := gb.Build(workspace)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Should have 0 edges (no internal deps found)
	if len(graph.Edges) != 0 {
		t.Errorf("Expected 0 edges, got %d", len(graph.Edges))
	}

	// Missing internal dep should be treated as external
	node := graph.Nodes["@mono/app"]
	if len(node.Dependencies) != 0 {
		t.Errorf("Expected 0 internal dependencies, got %d: %v", len(node.Dependencies), node.Dependencies)
	}
	if node.ExternalDeps["@mono/nonexistent"] != "^1.0.0" {
		t.Errorf("Missing dep should be in externalDeps: %v", node.ExternalDeps)
	}
	if node.ExternalDeps["lodash"] != "^4.17.0" {
		t.Errorf("External npm package should be in externalDeps: %v", node.ExternalDeps)
	}
}

// TestBuildIsolatedNodes verifies packages with no deps are included as nodes.
func TestBuildIsolatedNodes(t *testing.T) {
	gb := NewGraphBuilder()
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/isolated-a": {
				Name:             "@mono/isolated-a",
				Version:          "1.0.0",
				Path:             "packages/isolated-a",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/isolated-b": {
				Name:             "@mono/isolated-b",
				Version:          "2.0.0",
				Path:             "packages/isolated-b",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/isolated-c": {
				Name:             "@mono/isolated-c",
				Version:          "3.0.0",
				Path:             "packages/isolated-c",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	graph, err := gb.Build(workspace)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// All 3 isolated nodes should exist
	if len(graph.Nodes) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(graph.Nodes))
	}

	// No edges since no dependencies
	if len(graph.Edges) != 0 {
		t.Errorf("Expected 0 edges, got %d", len(graph.Edges))
	}

	// Verify each node has correct data
	expectedVersions := map[string]string{
		"@mono/isolated-a": "1.0.0",
		"@mono/isolated-b": "2.0.0",
		"@mono/isolated-c": "3.0.0",
	}

	for name, version := range expectedVersions {
		node := graph.Nodes[name]
		if node == nil {
			t.Errorf("Missing node: %s", name)
			continue
		}
		if node.Version != version {
			t.Errorf("Node %s version = %q, want %q", name, node.Version, version)
		}
	}
}

// TestBuildDuplicateDependencyEntries verifies last value wins for duplicates.
// Note: In practice, JSON parsing handles this - Go maps don't allow duplicates.
// This test verifies the behavior is predictable.
func TestBuildDuplicateDependencyEntries(t *testing.T) {
	gb := NewGraphBuilder()

	// Simulate what would happen if a package had "duplicate" entries
	// In practice, JSON maps can't have duplicate keys, so this tests
	// that our code handles the normal case correctly.
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/app": {
				Name:    "@mono/app",
				Version: "1.0.0",
				Path:    "apps/web",
				Dependencies: map[string]string{
					"@mono/lib": "^2.0.0", // Only one version can exist in a map
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/lib": {
				Name:             "@mono/lib",
				Version:          "2.0.0",
				Path:             "packages/lib",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	graph, err := gb.Build(workspace)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Should have exactly 1 edge
	if len(graph.Edges) != 1 {
		t.Errorf("Expected 1 edge, got %d", len(graph.Edges))
	}

	// Verify the edge has the expected version
	if len(graph.Edges) > 0 {
		edge := graph.Edges[0]
		if edge.VersionRange != "^2.0.0" {
			t.Errorf("Edge versionRange = %q, want ^2.0.0", edge.VersionRange)
		}
	}
}

// TestBuildMixedInternalExternalInAllDeptypes verifies correct classification across all types.
func TestBuildMixedInternalExternalInAllDeptypes(t *testing.T) {
	gb := NewGraphBuilder()
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/app": {
				Name:    "@mono/app",
				Version: "1.0.0",
				Path:    "apps/web",
				Dependencies: map[string]string{
					"@mono/internal-prod": "^1.0.0",
					"external-prod":       "^1.0.0",
				},
				DevDependencies: map[string]string{
					"@mono/internal-dev": "^1.0.0",
					"external-dev":       "^1.0.0",
				},
				PeerDependencies: map[string]string{
					"@mono/internal-peer": "^1.0.0",
					"external-peer":       "^1.0.0",
				},
			},
			"@mono/internal-prod": {
				Name:             "@mono/internal-prod",
				Version:          "1.0.0",
				Path:             "packages/internal-prod",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/internal-dev": {
				Name:             "@mono/internal-dev",
				Version:          "1.0.0",
				Path:             "packages/internal-dev",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/internal-peer": {
				Name:             "@mono/internal-peer",
				Version:          "1.0.0",
				Path:             "packages/internal-peer",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	graph, err := gb.Build(workspace)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Should have 3 internal edges (one of each type)
	if len(graph.Edges) != 3 {
		t.Errorf("Expected 3 edges, got %d", len(graph.Edges))
	}

	appNode := graph.Nodes["@mono/app"]

	// Verify internal deps in arrays
	if len(appNode.Dependencies) != 1 || appNode.Dependencies[0] != "@mono/internal-prod" {
		t.Errorf("Dependencies = %v, want [\"@mono/internal-prod\"]", appNode.Dependencies)
	}
	if len(appNode.DevDependencies) != 1 || appNode.DevDependencies[0] != "@mono/internal-dev" {
		t.Errorf("DevDependencies = %v, want [\"@mono/internal-dev\"]", appNode.DevDependencies)
	}
	if len(appNode.PeerDependencies) != 1 || appNode.PeerDependencies[0] != "@mono/internal-peer" {
		t.Errorf("PeerDependencies = %v, want [\"@mono/internal-peer\"]", appNode.PeerDependencies)
	}

	// Verify external deps in maps
	if appNode.ExternalDeps["external-prod"] != "^1.0.0" {
		t.Errorf("ExternalDeps missing external-prod: %v", appNode.ExternalDeps)
	}
	if appNode.ExternalDevDeps["external-dev"] != "^1.0.0" {
		t.Errorf("ExternalDevDeps missing external-dev: %v", appNode.ExternalDevDeps)
	}
	// Note: External peer deps are not stored separately per design
}

// ========================================
// Performance Tests
// ========================================

// ========================================
// Bidirectional Edge Tests (H3 fix)
// ========================================

// TestBuildBidirectionalDependencies verifies A↔B creates two edges (for cycle detection in Story 2.3).
func TestBuildBidirectionalDependencies(t *testing.T) {
	gb := NewGraphBuilder()
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
					"@mono/a": "^1.0.0", // Circular dependency
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	graph, err := gb.Build(workspace)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Should have 2 nodes
	if len(graph.Nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(graph.Nodes))
	}

	// Should have 2 edges (A->B and B->A)
	if len(graph.Edges) != 2 {
		t.Errorf("Expected 2 edges for bidirectional deps, got %d", len(graph.Edges))
	}

	// Verify both edges exist
	edgeMap := make(map[string]bool)
	for _, edge := range graph.Edges {
		key := edge.From + "->" + edge.To
		edgeMap[key] = true
	}

	if !edgeMap["@mono/a->@mono/b"] {
		t.Error("Missing edge @mono/a -> @mono/b")
	}
	if !edgeMap["@mono/b->@mono/a"] {
		t.Error("Missing edge @mono/b -> @mono/a")
	}
}

// ========================================
// Optional Dependencies Tests (H1 fix)
// ========================================

// TestBuildOptionalDependencies verifies optional dependencies create edges with correct type.
func TestBuildOptionalDependencies(t *testing.T) {
	gb := NewGraphBuilder()
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
				OptionalDependencies: map[string]string{
					"@mono/optional-plugin": "^1.0.0", // Internal optional
					"optional-external":     "^2.0.0", // External optional
				},
			},
			"@mono/optional-plugin": {
				Name:             "@mono/optional-plugin",
				Version:          "1.0.0",
				Path:             "packages/optional-plugin",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	graph, err := gb.Build(workspace)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Should have 1 edge with type "optional"
	if len(graph.Edges) != 1 {
		t.Errorf("Expected 1 edge, got %d", len(graph.Edges))
	}

	if len(graph.Edges) > 0 {
		edge := graph.Edges[0]
		if edge.Type != types.DependencyTypeOptional {
			t.Errorf("Edge type = %q, want %q", edge.Type, types.DependencyTypeOptional)
		}
		if edge.From != "@mono/app" || edge.To != "@mono/optional-plugin" {
			t.Errorf("Edge = %s -> %s, want @mono/app -> @mono/optional-plugin", edge.From, edge.To)
		}
	}

	// Verify optional dependencies stored in node
	appNode := graph.Nodes["@mono/app"]
	if len(appNode.OptionalDependencies) != 1 || appNode.OptionalDependencies[0] != "@mono/optional-plugin" {
		t.Errorf("OptionalDependencies = %v, want [\"@mono/optional-plugin\"]", appNode.OptionalDependencies)
	}

	// Verify external optional deps stored
	if appNode.ExternalOptionalDeps["optional-external"] != "^2.0.0" {
		t.Errorf("ExternalOptionalDeps missing optional-external: %v", appNode.ExternalOptionalDeps)
	}
}

// ========================================
// Empty Package Name Validation (M1 fix)
// ========================================

// TestBuildEmptyPackageName verifies empty package names return error.
func TestBuildEmptyPackageName(t *testing.T) {
	gb := NewGraphBuilder()
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"": { // Empty package name
				Name:             "",
				Version:          "1.0.0",
				Path:             "packages/invalid",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	_, err := gb.Build(workspace)
	if err == nil {
		t.Error("Expected error for empty package name, got nil")
	}
}

// ========================================
// Complex Version Range Tests (M4 fix)
// ========================================

// TestBuildComplexVersionRanges verifies various version range formats are preserved.
func TestBuildComplexVersionRanges(t *testing.T) {
	gb := NewGraphBuilder()
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/app": {
				Name:    "@mono/app",
				Version: "1.0.0",
				Path:    "apps/web",
				Dependencies: map[string]string{
					"@mono/semver":    ">=1.0.0 <2.0.0",  // Range
					"@mono/workspace": "workspace:*",     // pnpm workspace protocol
					"@mono/exact":     "1.2.3",           // Exact version
					"@mono/tilde":     "~1.2.0",          // Tilde range
					"@mono/caret":     "^1.0.0",          // Caret range
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/semver": {
				Name:             "@mono/semver",
				Version:          "1.5.0",
				Path:             "packages/semver",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/workspace": {
				Name:             "@mono/workspace",
				Version:          "2.0.0",
				Path:             "packages/workspace",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/exact": {
				Name:             "@mono/exact",
				Version:          "1.2.3",
				Path:             "packages/exact",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/tilde": {
				Name:             "@mono/tilde",
				Version:          "1.2.5",
				Path:             "packages/tilde",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/caret": {
				Name:             "@mono/caret",
				Version:          "1.9.0",
				Path:             "packages/caret",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	graph, err := gb.Build(workspace)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Verify all edges created with correct version ranges
	expectedVersionRanges := map[string]string{
		"@mono/semver":    ">=1.0.0 <2.0.0",
		"@mono/workspace": "workspace:*",
		"@mono/exact":     "1.2.3",
		"@mono/tilde":     "~1.2.0",
		"@mono/caret":     "^1.0.0",
	}

	if len(graph.Edges) != 5 {
		t.Errorf("Expected 5 edges, got %d", len(graph.Edges))
	}

	for _, edge := range graph.Edges {
		if edge.From != "@mono/app" {
			continue
		}
		expected, ok := expectedVersionRanges[edge.To]
		if !ok {
			t.Errorf("Unexpected edge to %s", edge.To)
			continue
		}
		if edge.VersionRange != expected {
			t.Errorf("Edge to %s versionRange = %q, want %q", edge.To, edge.VersionRange, expected)
		}
	}
}

// ========================================
// External Peer Dependencies Test (H2 fix)
// ========================================

// TestBuildExternalPeerDependencies verifies external peer deps are stored (not discarded).
func TestBuildExternalPeerDependencies(t *testing.T) {
	gb := NewGraphBuilder()
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/plugin": {
				Name:            "@mono/plugin",
				Version:         "1.0.0",
				Path:            "packages/plugin",
				Dependencies:    map[string]string{},
				DevDependencies: map[string]string{},
				PeerDependencies: map[string]string{
					"react":      "^18.0.0", // External peer
					"@mono/core": "^1.0.0",  // Internal peer
				},
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

	graph, err := gb.Build(workspace)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	pluginNode := graph.Nodes["@mono/plugin"]

	// Verify internal peer dep in array
	if len(pluginNode.PeerDependencies) != 1 || pluginNode.PeerDependencies[0] != "@mono/core" {
		t.Errorf("PeerDependencies = %v, want [\"@mono/core\"]", pluginNode.PeerDependencies)
	}

	// Verify external peer dep stored (not discarded)
	if pluginNode.ExternalPeerDeps["react"] != "^18.0.0" {
		t.Errorf("ExternalPeerDeps missing react: %v", pluginNode.ExternalPeerDeps)
	}
}

// ========================================
// Deterministic Output Tests (H4-H5 fix)
// ========================================

// TestBuildDeterministicEdgeOrder verifies edges are always in same order.
func TestBuildDeterministicEdgeOrder(t *testing.T) {
	gb := NewGraphBuilder()
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/z-app": {
				Name:    "@mono/z-app",
				Version: "1.0.0",
				Path:    "apps/z-app",
				Dependencies: map[string]string{
					"@mono/b-lib": "^1.0.0",
					"@mono/a-lib": "^1.0.0",
					"@mono/c-lib": "^1.0.0",
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/a-lib": {
				Name:             "@mono/a-lib",
				Version:          "1.0.0",
				Path:             "packages/a-lib",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/b-lib": {
				Name:             "@mono/b-lib",
				Version:          "1.0.0",
				Path:             "packages/b-lib",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/c-lib": {
				Name:             "@mono/c-lib",
				Version:          "1.0.0",
				Path:             "packages/c-lib",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	// Build multiple times and verify same order
	for i := 0; i < 5; i++ {
		graph, err := gb.Build(workspace)
		if err != nil {
			t.Fatalf("Build failed: %v", err)
		}

		// Edges should always be sorted: z-app->a-lib, z-app->b-lib, z-app->c-lib
		expectedOrder := []string{"@mono/a-lib", "@mono/b-lib", "@mono/c-lib"}
		for j, edge := range graph.Edges {
			if edge.To != expectedOrder[j] {
				t.Errorf("Iteration %d: Edge %d to = %s, want %s", i, j, edge.To, expectedOrder[j])
			}
		}
	}
}

// TestBuildDeterministicDependencyOrder verifies node.Dependencies are always in same order.
func TestBuildDeterministicDependencyOrder(t *testing.T) {
	gb := NewGraphBuilder()
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/app": {
				Name:    "@mono/app",
				Version: "1.0.0",
				Path:    "apps/web",
				Dependencies: map[string]string{
					"@mono/zebra": "^1.0.0",
					"@mono/apple": "^1.0.0",
					"@mono/mango": "^1.0.0",
				},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			"@mono/zebra": {Name: "@mono/zebra", Version: "1.0.0", Path: "packages/zebra", Dependencies: map[string]string{}, DevDependencies: map[string]string{}, PeerDependencies: map[string]string{}},
			"@mono/apple": {Name: "@mono/apple", Version: "1.0.0", Path: "packages/apple", Dependencies: map[string]string{}, DevDependencies: map[string]string{}, PeerDependencies: map[string]string{}},
			"@mono/mango": {Name: "@mono/mango", Version: "1.0.0", Path: "packages/mango", Dependencies: map[string]string{}, DevDependencies: map[string]string{}, PeerDependencies: map[string]string{}},
		},
	}

	// Build multiple times and verify same order
	for i := 0; i < 5; i++ {
		graph, err := gb.Build(workspace)
		if err != nil {
			t.Fatalf("Build failed: %v", err)
		}

		appNode := graph.Nodes["@mono/app"]
		// Dependencies should be sorted alphabetically
		expectedOrder := []string{"@mono/apple", "@mono/mango", "@mono/zebra"}
		for j, dep := range appNode.Dependencies {
			if dep != expectedOrder[j] {
				t.Errorf("Iteration %d: Dependency %d = %s, want %s", i, j, dep, expectedOrder[j])
			}
		}
	}
}

// ========================================
// Performance Tests
// ========================================

// TestBuildNodeLookupPerformance verifies O(1) node lookup.
func TestBuildNodeLookupPerformance(t *testing.T) {
	gb := NewGraphBuilder()
	packages := make(map[string]*types.PackageInfo)

	// Create 100 packages
	for i := 0; i < 100; i++ {
		name := "@mono/pkg-" + string(rune('a'+i%26)) + string(rune('0'+i/26))
		packages[name] = &types.PackageInfo{
			Name:             name,
			Version:          "1.0.0",
			Path:             "packages/" + name,
			Dependencies:     map[string]string{},
			DevDependencies:  map[string]string{},
			PeerDependencies: map[string]string{},
		}
	}

	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages:      packages,
	}

	graph, err := gb.Build(workspace)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Verify all nodes accessible by name (O(1) lookup via map)
	for name := range packages {
		node := graph.Nodes[name]
		if node == nil {
			t.Errorf("Node %s not found in graph", name)
		}
	}
}

// ========================================
// Story 2.6: Exclusion Tests
// ========================================

// TestNewGraphBuilderWithExclusions verifies exclusion matcher integration.
func TestNewGraphBuilderWithExclusions(t *testing.T) {
	patterns := []string{"@mono/legacy", "@mono/deprecated-*", "regex:.*-test$"}

	gb, err := NewGraphBuilderWithExclusions(patterns)
	if err != nil {
		t.Fatalf("NewGraphBuilderWithExclusions failed: %v", err)
	}
	if gb == nil {
		t.Fatal("GraphBuilder is nil")
	}
	if gb.exclusionMatcher == nil {
		t.Error("exclusionMatcher is nil")
	}
	if gb.exclusionMatcher.PatternCount() != 3 {
		t.Errorf("PatternCount = %d, want 3", gb.exclusionMatcher.PatternCount())
	}
}

// TestNewGraphBuilderWithExclusionsInvalidRegex verifies error on invalid regex.
func TestNewGraphBuilderWithExclusionsInvalidRegex(t *testing.T) {
	patterns := []string{"regex:[invalid"}

	_, err := NewGraphBuilderWithExclusions(patterns)
	if err == nil {
		t.Error("Expected error for invalid regex, got nil")
	}
}

// TestNewGraphBuilderWithExclusionsEmpty verifies empty patterns work.
func TestNewGraphBuilderWithExclusionsEmpty(t *testing.T) {
	gb, err := NewGraphBuilderWithExclusions([]string{})
	if err != nil {
		t.Fatalf("NewGraphBuilderWithExclusions failed: %v", err)
	}
	if gb == nil {
		t.Fatal("GraphBuilder is nil")
	}
	if gb.exclusionMatcher == nil {
		t.Error("exclusionMatcher should not be nil")
	}
}

// TestBuildWithExclusionFlag verifies excluded packages are marked.
func TestBuildWithExclusionFlag(t *testing.T) {
	gb, err := NewGraphBuilderWithExclusions([]string{"@mono/legacy"})
	if err != nil {
		t.Fatalf("NewGraphBuilderWithExclusions failed: %v", err)
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
			"@mono/legacy": {
				Name:             "@mono/legacy",
				Version:          "1.0.0",
				Path:             "packages/legacy",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
		},
	}

	graph, err := gb.Build(workspace)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Verify legacy is marked excluded
	legacyNode := graph.Nodes["@mono/legacy"]
	if legacyNode == nil {
		t.Fatal("Legacy node missing")
	}
	if !legacyNode.Excluded {
		t.Error("Legacy node should be marked as excluded")
	}

	// Verify app is not marked excluded
	appNode := graph.Nodes["@mono/app"]
	if appNode == nil {
		t.Fatal("App node missing")
	}
	if appNode.Excluded {
		t.Error("App node should not be marked as excluded")
	}
}

// TestBuildWithGlobExclusionFlag verifies glob patterns mark packages.
func TestBuildWithGlobExclusionFlag(t *testing.T) {
	gb, err := NewGraphBuilderWithExclusions([]string{"@mono/deprecated-*"})
	if err != nil {
		t.Fatalf("NewGraphBuilderWithExclusions failed: %v", err)
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

	graph, err := gb.Build(workspace)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Count excluded packages
	excludedCount := 0
	for _, node := range graph.Nodes {
		if node.Excluded {
			excludedCount++
		}
	}

	if excludedCount != 2 {
		t.Errorf("Excluded count = %d, want 2", excludedCount)
	}
}

// TestBuildWithNoExclusion verifies normal builder has no exclusions.
func TestBuildWithNoExclusion(t *testing.T) {
	gb := NewGraphBuilder() // No exclusions

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

	graph, err := gb.Build(workspace)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// No packages should be excluded
	for name, node := range graph.Nodes {
		if node.Excluded {
			t.Errorf("Node %s should not be excluded (no exclusion matcher)", name)
		}
	}
}

// TestBuildExcludedPackagesStillInGraph verifies AC5: excluded packages ARE in graph.
func TestBuildExcludedPackagesStillInGraph(t *testing.T) {
	gb, err := NewGraphBuilderWithExclusions([]string{"@mono/legacy"})
	if err != nil {
		t.Fatalf("NewGraphBuilderWithExclusions failed: %v", err)
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
		},
	}

	graph, err := gb.Build(workspace)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// AC5: Excluded packages ARE included in graph nodes (for visualization)
	if len(graph.Nodes) != 2 {
		t.Errorf("Graph.Nodes count = %d, want 2 (including excluded)", len(graph.Nodes))
	}

	// AC5: Edges to/from excluded packages are still included
	if len(graph.Edges) != 1 {
		t.Errorf("Graph.Edges count = %d, want 1 (edge to excluded package)", len(graph.Edges))
	}

	// Verify the excluded node exists with flag
	legacyNode := graph.Nodes["@mono/legacy"]
	if legacyNode == nil {
		t.Fatal("Excluded node should still be in graph")
	}
	if !legacyNode.Excluded {
		t.Error("Excluded node should have Excluded=true")
	}
}
