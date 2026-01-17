package types

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestDependencyTypeConstants verifies the dependency type constants.
func TestDependencyTypeConstants(t *testing.T) {
	tests := []struct {
		depType  DependencyType
		expected string
	}{
		{DependencyTypeProduction, "production"},
		{DependencyTypeDevelopment, "development"},
		{DependencyTypePeer, "peer"},
		{DependencyTypeOptional, "optional"},
	}

	for _, tt := range tests {
		t.Run(string(tt.depType), func(t *testing.T) {
			if string(tt.depType) != tt.expected {
				t.Errorf("DependencyType = %q, want %q", tt.depType, tt.expected)
			}
		})
	}
}

// TestPackageNodeJSONSerialization verifies camelCase JSON output for PackageNode.
func TestPackageNodeJSONSerialization(t *testing.T) {
	node := &PackageNode{
		Name:             "@mono/app",
		Version:          "1.0.0",
		Path:             "apps/web",
		Dependencies:     []string{"@mono/ui"},
		DevDependencies:  []string{"@mono/types"},
		PeerDependencies: []string{},
		ExternalDeps:     map[string]string{"react": "^18.0.0"},
		ExternalDevDeps:  map[string]string{"typescript": "^5.0.0"},
	}

	data, err := json.Marshal(node)
	if err != nil {
		t.Fatalf("Failed to marshal PackageNode: %v", err)
	}

	jsonStr := string(data)

	// Verify camelCase field names
	expectedFields := []string{
		`"name"`,
		`"version"`,
		`"path"`,
		`"dependencies"`,
		`"devDependencies"`,
		`"peerDependencies"`,
		`"externalDeps"`,
		`"externalDevDeps"`,
	}

	for _, field := range expectedFields {
		if !strings.Contains(jsonStr, field) {
			t.Errorf("JSON missing camelCase field %s, got: %s", field, jsonStr)
		}
	}

	// Verify NO snake_case
	snakeCaseFields := []string{
		`"dev_dependencies"`,
		`"peer_dependencies"`,
		`"external_deps"`,
		`"external_dev_deps"`,
	}

	for _, field := range snakeCaseFields {
		if strings.Contains(jsonStr, field) {
			t.Errorf("JSON contains snake_case field %s (should be camelCase)", field)
		}
	}
}

// TestDependencyEdgeJSONSerialization verifies camelCase JSON output for DependencyEdge.
func TestDependencyEdgeJSONSerialization(t *testing.T) {
	edge := &DependencyEdge{
		From:         "@mono/app",
		To:           "@mono/ui",
		Type:         DependencyTypeProduction,
		VersionRange: "^1.0.0",
	}

	data, err := json.Marshal(edge)
	if err != nil {
		t.Fatalf("Failed to marshal DependencyEdge: %v", err)
	}

	jsonStr := string(data)

	// Verify camelCase field names
	expectedFields := []string{
		`"from"`,
		`"to"`,
		`"type"`,
		`"versionRange"`,
	}

	for _, field := range expectedFields {
		if !strings.Contains(jsonStr, field) {
			t.Errorf("JSON missing camelCase field %s, got: %s", field, jsonStr)
		}
	}

	// Verify type value is correct
	if !strings.Contains(jsonStr, `"type":"production"`) {
		t.Errorf("Expected type:production, got: %s", jsonStr)
	}

	// Verify NO snake_case
	if strings.Contains(jsonStr, `"version_range"`) {
		t.Errorf("JSON contains snake_case version_range (should be versionRange)")
	}
}

// TestDependencyGraphJSONSerialization verifies the complete graph serialization.
func TestDependencyGraphJSONSerialization(t *testing.T) {
	graph := &DependencyGraph{
		Nodes: map[string]*PackageNode{
			"@mono/app": {
				Name:             "@mono/app",
				Version:          "1.0.0",
				Path:             "apps/web",
				Dependencies:     []string{"@mono/ui"},
				DevDependencies:  []string{"@mono/types"},
				PeerDependencies: []string{},
				ExternalDeps:     map[string]string{"react": "^18.0.0"},
			},
			"@mono/ui": {
				Name:             "@mono/ui",
				Version:          "1.0.0",
				Path:             "packages/ui",
				Dependencies:     []string{},
				DevDependencies:  []string{},
				PeerDependencies: []string{},
				ExternalDeps:     map[string]string{"react": "^18.0.0"},
			},
		},
		Edges: []*DependencyEdge{
			{
				From:         "@mono/app",
				To:           "@mono/ui",
				Type:         DependencyTypeProduction,
				VersionRange: "^1.0.0",
			},
		},
		RootPath:      "/workspace",
		WorkspaceType: WorkspaceTypePnpm,
	}

	data, err := json.Marshal(graph)
	if err != nil {
		t.Fatalf("Failed to marshal DependencyGraph: %v", err)
	}

	jsonStr := string(data)

	// Verify top-level camelCase field names
	expectedFields := []string{
		`"nodes"`,
		`"edges"`,
		`"rootPath"`,
		`"workspaceType"`,
	}

	for _, field := range expectedFields {
		if !strings.Contains(jsonStr, field) {
			t.Errorf("JSON missing camelCase field %s, got: %s", field, jsonStr)
		}
	}

	// Verify NO snake_case
	snakeCaseFields := []string{
		`"root_path"`,
		`"workspace_type"`,
	}

	for _, field := range snakeCaseFields {
		if strings.Contains(jsonStr, field) {
			t.Errorf("JSON contains snake_case field %s (should be camelCase)", field)
		}
	}
}

// TestDependencyGraphRoundTrip verifies JSON marshal/unmarshal preserves data.
func TestDependencyGraphRoundTrip(t *testing.T) {
	original := &DependencyGraph{
		Nodes: map[string]*PackageNode{
			"@mono/app": {
				Name:             "@mono/app",
				Version:          "1.0.0",
				Path:             "apps/web",
				Dependencies:     []string{"@mono/ui", "@mono/utils"},
				DevDependencies:  []string{"@mono/types"},
				PeerDependencies: []string{},
				ExternalDeps:     map[string]string{"react": "^18.0.0", "lodash": "^4.17.0"},
				ExternalDevDeps:  map[string]string{"typescript": "^5.0.0"},
			},
		},
		Edges: []*DependencyEdge{
			{From: "@mono/app", To: "@mono/ui", Type: DependencyTypeProduction, VersionRange: "^1.0.0"},
			{From: "@mono/app", To: "@mono/utils", Type: DependencyTypeProduction, VersionRange: "workspace:*"},
			{From: "@mono/app", To: "@mono/types", Type: DependencyTypeDevelopment, VersionRange: "^1.0.0"},
		},
		RootPath:      "/workspace",
		WorkspaceType: WorkspaceTypePnpm,
	}

	// Marshal
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Unmarshal
	var decoded DependencyGraph
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Verify key fields
	if decoded.RootPath != original.RootPath {
		t.Errorf("RootPath = %q, want %q", decoded.RootPath, original.RootPath)
	}
	if decoded.WorkspaceType != original.WorkspaceType {
		t.Errorf("WorkspaceType = %q, want %q", decoded.WorkspaceType, original.WorkspaceType)
	}
	if len(decoded.Nodes) != len(original.Nodes) {
		t.Errorf("Nodes count = %d, want %d", len(decoded.Nodes), len(original.Nodes))
	}
	if len(decoded.Edges) != len(original.Edges) {
		t.Errorf("Edges count = %d, want %d", len(decoded.Edges), len(original.Edges))
	}

	// Verify node data
	appNode := decoded.Nodes["@mono/app"]
	if appNode == nil {
		t.Fatal("Missing @mono/app node")
	}
	if appNode.Name != "@mono/app" {
		t.Errorf("Node name = %q, want %q", appNode.Name, "@mono/app")
	}
	if len(appNode.Dependencies) != 2 {
		t.Errorf("Dependencies count = %d, want 2", len(appNode.Dependencies))
	}
	if len(appNode.ExternalDeps) != 2 {
		t.Errorf("ExternalDeps count = %d, want 2", len(appNode.ExternalDeps))
	}
}

// TestEmptyGraphSerialization verifies empty graphs serialize correctly.
func TestEmptyGraphSerialization(t *testing.T) {
	graph := &DependencyGraph{
		Nodes:         map[string]*PackageNode{},
		Edges:         []*DependencyEdge{},
		RootPath:      "/workspace",
		WorkspaceType: WorkspaceTypeUnknown,
	}

	data, err := json.Marshal(graph)
	if err != nil {
		t.Fatalf("Failed to marshal empty graph: %v", err)
	}

	var decoded DependencyGraph
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal empty graph: %v", err)
	}

	if len(decoded.Nodes) != 0 {
		t.Errorf("Expected empty nodes, got %d", len(decoded.Nodes))
	}
	if len(decoded.Edges) != 0 {
		t.Errorf("Expected empty edges, got %d", len(decoded.Edges))
	}
}

// TestPackageNodeEmptySlices verifies empty slices serialize as arrays, not null.
func TestPackageNodeEmptySlices(t *testing.T) {
	node := &PackageNode{
		Name:             "@mono/lib",
		Version:          "1.0.0",
		Path:             "packages/lib",
		Dependencies:     []string{},
		DevDependencies:  []string{},
		PeerDependencies: []string{},
	}

	data, err := json.Marshal(node)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	jsonStr := string(data)

	// Verify empty arrays (not null)
	if strings.Contains(jsonStr, `"dependencies":null`) {
		t.Error("dependencies should be [] not null")
	}
	if strings.Contains(jsonStr, `"devDependencies":null`) {
		t.Error("devDependencies should be [] not null")
	}
	if strings.Contains(jsonStr, `"peerDependencies":null`) {
		t.Error("peerDependencies should be [] not null")
	}
}

