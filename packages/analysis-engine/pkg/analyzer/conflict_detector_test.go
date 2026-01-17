package analyzer

import (
	"strings"
	"testing"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// createConflictTestGraph creates a test dependency graph with specified external dependencies.
func createConflictTestGraph(packages map[string]map[string]map[string]string) *types.DependencyGraph {
	graph := types.NewDependencyGraph("/test", types.WorkspaceTypeNpm)

	for pkgName, deps := range packages {
		node := types.NewPackageNode(pkgName, "1.0.0", "/test/"+pkgName)

		if prodDeps, ok := deps["production"]; ok {
			node.ExternalDeps = prodDeps
		}
		if devDeps, ok := deps["development"]; ok {
			node.ExternalDevDeps = devDeps
		}
		if peerDeps, ok := deps["peer"]; ok {
			node.ExternalPeerDeps = peerDeps
		}

		graph.Nodes[pkgName] = node
	}

	return graph
}

func TestConflictDetector_NoConflicts(t *testing.T) {
	// All packages use the same version of lodash
	graph := createConflictTestGraph(map[string]map[string]map[string]string{
		"@mono/app": {
			"production": {"lodash": "^4.17.21"},
		},
		"@mono/lib": {
			"production": {"lodash": "^4.17.21"},
		},
	})

	detector := NewConflictDetector(graph)
	conflicts := detector.DetectConflicts()

	if len(conflicts) != 0 {
		t.Errorf("Expected no conflicts, got %d", len(conflicts))
	}
}

func TestConflictDetector_PatchVersionConflict(t *testing.T) {
	// Patch version difference: 4.17.19 vs 4.17.21
	graph := createConflictTestGraph(map[string]map[string]map[string]string{
		"@mono/app": {
			"production": {"lodash": "^4.17.21"},
		},
		"@mono/lib": {
			"production": {"lodash": "^4.17.19"},
		},
	})

	detector := NewConflictDetector(graph)
	conflicts := detector.DetectConflicts()

	if len(conflicts) != 1 {
		t.Fatalf("Expected 1 conflict, got %d", len(conflicts))
	}

	conflict := conflicts[0]
	if conflict.PackageName != "lodash" {
		t.Errorf("Expected package name 'lodash', got '%s'", conflict.PackageName)
	}

	if conflict.Severity != types.ConflictSeverityInfo {
		t.Errorf("Expected severity 'info' for patch difference, got '%s'", conflict.Severity)
	}

	if len(conflict.ConflictingVersions) != 2 {
		t.Errorf("Expected 2 conflicting versions, got %d", len(conflict.ConflictingVersions))
	}
}

func TestConflictDetector_MinorVersionConflict(t *testing.T) {
	// Minor version difference: 4.17.x vs 4.18.x
	graph := createConflictTestGraph(map[string]map[string]map[string]string{
		"@mono/app": {
			"production": {"lodash": "^4.18.0"},
		},
		"@mono/lib": {
			"production": {"lodash": "^4.17.0"},
		},
	})

	detector := NewConflictDetector(graph)
	conflicts := detector.DetectConflicts()

	if len(conflicts) != 1 {
		t.Fatalf("Expected 1 conflict, got %d", len(conflicts))
	}

	if conflicts[0].Severity != types.ConflictSeverityWarning {
		t.Errorf("Expected severity 'warning' for minor difference, got '%s'", conflicts[0].Severity)
	}
}

func TestConflictDetector_MajorVersionConflict(t *testing.T) {
	// Major version difference: 4.x vs 5.x
	graph := createConflictTestGraph(map[string]map[string]map[string]string{
		"@mono/app": {
			"production": {"typescript": "^5.0.0"},
		},
		"@mono/lib": {
			"production": {"typescript": "^4.9.0"},
		},
	})

	detector := NewConflictDetector(graph)
	conflicts := detector.DetectConflicts()

	if len(conflicts) != 1 {
		t.Fatalf("Expected 1 conflict, got %d", len(conflicts))
	}

	conflict := conflicts[0]
	if conflict.Severity != types.ConflictSeverityCritical {
		t.Errorf("Expected severity 'critical' for major difference, got '%s'", conflict.Severity)
	}

	// Check that one version is marked as breaking
	hasBreaking := false
	for _, cv := range conflict.ConflictingVersions {
		if cv.IsBreaking {
			hasBreaking = true
			break
		}
	}
	if !hasBreaking {
		t.Error("Expected at least one version to be marked as breaking")
	}
}

func TestConflictDetector_MultipleConflicts(t *testing.T) {
	// Multiple packages with conflicts
	graph := createConflictTestGraph(map[string]map[string]map[string]string{
		"@mono/app": {
			"production":  {"lodash": "^4.17.21", "react": "^18.2.0"},
			"development": {"typescript": "^5.0.0"},
		},
		"@mono/lib": {
			"production":  {"lodash": "^4.17.19", "react": "^17.0.0"},
			"development": {"typescript": "^4.9.0"},
		},
	})

	detector := NewConflictDetector(graph)
	conflicts := detector.DetectConflicts()

	if len(conflicts) != 3 {
		t.Errorf("Expected 3 conflicts (lodash, react, typescript), got %d", len(conflicts))
	}

	// Verify conflicts are sorted by package name
	for i := 0; i < len(conflicts)-1; i++ {
		if conflicts[i].PackageName > conflicts[i+1].PackageName {
			t.Errorf("Conflicts not sorted: %s > %s", conflicts[i].PackageName, conflicts[i+1].PackageName)
		}
	}
}

func TestConflictDetector_DevDependencyConflict(t *testing.T) {
	// Conflict in dev dependencies
	graph := createConflictTestGraph(map[string]map[string]map[string]string{
		"@mono/app": {
			"development": {"eslint": "^8.0.0"},
		},
		"@mono/lib": {
			"development": {"eslint": "^7.0.0"},
		},
	})

	detector := NewConflictDetector(graph)
	conflicts := detector.DetectConflicts()

	if len(conflicts) != 1 {
		t.Fatalf("Expected 1 conflict, got %d", len(conflicts))
	}

	// Check depType
	conflict := conflicts[0]
	for _, cv := range conflict.ConflictingVersions {
		if cv.DepType != types.DepTypeDevelopment {
			t.Errorf("Expected depType 'development', got '%s'", cv.DepType)
		}
	}
}

func TestConflictDetector_PeerDependencyConflict(t *testing.T) {
	// Conflict in peer dependencies
	graph := createConflictTestGraph(map[string]map[string]map[string]string{
		"@mono/app": {
			"peer": {"react": "^18.0.0"},
		},
		"@mono/lib": {
			"peer": {"react": "^17.0.0"},
		},
	})

	detector := NewConflictDetector(graph)
	conflicts := detector.DetectConflicts()

	if len(conflicts) != 1 {
		t.Fatalf("Expected 1 conflict, got %d", len(conflicts))
	}

	// Check severity (major version diff = critical)
	conflict := conflicts[0]
	if conflict.Severity != types.ConflictSeverityCritical {
		t.Errorf("Expected severity 'critical' for major peer dep difference, got '%s'", conflict.Severity)
	}

	// Check depType
	for _, cv := range conflict.ConflictingVersions {
		if cv.DepType != types.DepTypePeer {
			t.Errorf("Expected depType 'peer', got '%s'", cv.DepType)
		}
	}
}

func TestConflictDetector_MixedDepTypes(t *testing.T) {
	// Same dependency used as prod in one package and dev in another
	graph := createConflictTestGraph(map[string]map[string]map[string]string{
		"@mono/app": {
			"production": {"lodash": "^4.17.21"},
		},
		"@mono/lib": {
			"development": {"lodash": "^4.17.19"},
		},
	})

	detector := NewConflictDetector(graph)
	conflicts := detector.DetectConflicts()

	if len(conflicts) != 1 {
		t.Fatalf("Expected 1 conflict, got %d", len(conflicts))
	}

	// Verify different depTypes are captured
	depTypes := make(map[string]bool)
	for _, cv := range conflicts[0].ConflictingVersions {
		depTypes[cv.DepType] = true
	}

	if !depTypes[types.DepTypeProduction] || !depTypes[types.DepTypeDevelopment] {
		t.Error("Expected both production and development depTypes")
	}
}

func TestConflictDetector_ThreeVersions(t *testing.T) {
	// Three different versions
	graph := createConflictTestGraph(map[string]map[string]map[string]string{
		"@mono/app": {
			"production": {"lodash": "^4.17.21"},
		},
		"@mono/lib": {
			"production": {"lodash": "^4.17.19"},
		},
		"@mono/utils": {
			"production": {"lodash": "^4.17.20"},
		},
	})

	detector := NewConflictDetector(graph)
	conflicts := detector.DetectConflicts()

	if len(conflicts) != 1 {
		t.Fatalf("Expected 1 conflict, got %d", len(conflicts))
	}

	if len(conflicts[0].ConflictingVersions) != 3 {
		t.Errorf("Expected 3 conflicting versions, got %d", len(conflicts[0].ConflictingVersions))
	}
}

func TestConflictDetector_MultiplePackagesSameVersion(t *testing.T) {
	// Multiple packages using the same version should be grouped
	graph := createConflictTestGraph(map[string]map[string]map[string]string{
		"@mono/app": {
			"production": {"lodash": "^4.17.21"},
		},
		"@mono/lib": {
			"production": {"lodash": "^4.17.21"},
		},
		"@mono/utils": {
			"production": {"lodash": "^4.17.19"},
		},
	})

	detector := NewConflictDetector(graph)
	conflicts := detector.DetectConflicts()

	if len(conflicts) != 1 {
		t.Fatalf("Expected 1 conflict, got %d", len(conflicts))
	}

	conflict := conflicts[0]
	if len(conflict.ConflictingVersions) != 2 {
		t.Errorf("Expected 2 unique versions, got %d", len(conflict.ConflictingVersions))
	}

	// Find the version with 2 packages
	for _, cv := range conflict.ConflictingVersions {
		if cv.Version == "^4.17.21" && len(cv.Packages) != 2 {
			t.Errorf("Expected 2 packages for version ^4.17.21, got %d", len(cv.Packages))
		}
	}
}

func TestConflictDetector_EmptyGraph(t *testing.T) {
	graph := types.NewDependencyGraph("/test", types.WorkspaceTypeNpm)

	detector := NewConflictDetector(graph)
	conflicts := detector.DetectConflicts()

	if conflicts != nil && len(conflicts) > 0 {
		t.Errorf("Expected no conflicts for empty graph, got %d", len(conflicts))
	}
}

func TestConflictDetector_NilGraph(t *testing.T) {
	detector := NewConflictDetector(nil)
	conflicts := detector.DetectConflicts()

	if conflicts != nil && len(conflicts) > 0 {
		t.Errorf("Expected no conflicts for nil graph, got %d", len(conflicts))
	}
}

func TestConflictDetector_UnparseableVersions(t *testing.T) {
	// Test with special version strings that cannot be parsed as semver
	graph := createConflictTestGraph(map[string]map[string]map[string]string{
		"@mono/app": {
			"production": {"some-pkg": "latest"},
		},
		"@mono/lib": {
			"production": {"some-pkg": "*"},
		},
		"@mono/utils": {
			"production": {"some-pkg": "^1.0.0"}, // Only parseable version
		},
	})

	detector := NewConflictDetector(graph)
	conflicts := detector.DetectConflicts()

	// Should still detect conflict (3 different version strings)
	if len(conflicts) != 1 {
		t.Fatalf("Expected 1 conflict, got %d", len(conflicts))
	}

	conflict := conflicts[0]
	if conflict.PackageName != "some-pkg" {
		t.Errorf("Expected package name 'some-pkg', got '%s'", conflict.PackageName)
	}

	// With unparseable versions, severity comparison may default to major
	// The behavior should be graceful, not panic
	if len(conflict.ConflictingVersions) != 3 {
		t.Errorf("Expected 3 conflicting versions, got %d", len(conflict.ConflictingVersions))
	}
}

func TestConflictDetector_ResolutionMessage(t *testing.T) {
	tests := []struct {
		name           string
		versions       map[string]string
		wantSeverity   types.ConflictSeverity
		containsInRes  string
	}{
		{
			name:          "critical resolution",
			versions:      map[string]string{"@mono/app": "^5.0.0", "@mono/lib": "^4.0.0"},
			wantSeverity:  types.ConflictSeverityCritical,
			containsInRes: "Major version conflict",
		},
		{
			name:          "warning resolution",
			versions:      map[string]string{"@mono/app": "^4.18.0", "@mono/lib": "^4.17.0"},
			wantSeverity:  types.ConflictSeverityWarning,
			containsInRes: "Consider upgrading",
		},
		{
			name:          "info resolution",
			versions:      map[string]string{"@mono/app": "^4.17.21", "@mono/lib": "^4.17.19"},
			wantSeverity:  types.ConflictSeverityInfo,
			containsInRes: "Patch version difference",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packages := make(map[string]map[string]map[string]string)
			for pkg, ver := range tt.versions {
				packages[pkg] = map[string]map[string]string{
					"production": {"testpkg": ver},
				}
			}

			graph := createConflictTestGraph(packages)
			detector := NewConflictDetector(graph)
			conflicts := detector.DetectConflicts()

			if len(conflicts) != 1 {
				t.Fatalf("Expected 1 conflict, got %d", len(conflicts))
			}

			conflict := conflicts[0]
			if conflict.Severity != tt.wantSeverity {
				t.Errorf("Severity = %s, want %s", conflict.Severity, tt.wantSeverity)
			}

			if !strings.Contains(conflict.Resolution, tt.containsInRes) {
				t.Errorf("Resolution %q should contain %q", conflict.Resolution, tt.containsInRes)
			}
		})
	}
}

func TestConflictDetector_ImpactMessage(t *testing.T) {
	graph := createConflictTestGraph(map[string]map[string]map[string]string{
		"@mono/app": {
			"production": {"lodash": "^5.0.0"},
		},
		"@mono/lib": {
			"production": {"lodash": "^4.0.0"},
		},
	})

	detector := NewConflictDetector(graph)
	conflicts := detector.DetectConflicts()

	if len(conflicts) != 1 {
		t.Fatalf("Expected 1 conflict, got %d", len(conflicts))
	}

	impact := conflicts[0].Impact
	if !strings.Contains(impact, "lodash") {
		t.Errorf("Impact should mention the package name")
	}
	if !strings.Contains(impact, "Breaking changes") {
		t.Errorf("Impact should mention breaking changes for critical severity")
	}
}

func TestDetermineSeverity(t *testing.T) {
	tests := []struct {
		name     string
		versions []string
		want     types.ConflictSeverity
	}{
		{"major diff", []string{"3.0.0", "4.0.0"}, types.ConflictSeverityCritical},
		{"minor diff", []string{"4.17.0", "4.18.0"}, types.ConflictSeverityWarning},
		{"patch diff", []string{"4.17.19", "4.17.21"}, types.ConflictSeverityInfo},
		{"mixed with major", []string{"3.0.0", "4.17.0", "4.18.0"}, types.ConflictSeverityCritical},
		{"single version", []string{"4.17.21"}, types.ConflictSeverityInfo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := determineSeverity(tt.versions)
			if result != tt.want {
				t.Errorf("determineSeverity(%v) = %s, want %s", tt.versions, result, tt.want)
			}
		})
	}
}

func TestIsBreakingVersion(t *testing.T) {
	tests := []struct {
		name        string
		version     string
		allVersions []string
		want        bool
	}{
		{"major diff exists", "5.0.0", []string{"4.0.0", "5.0.0"}, true},
		{"no major diff", "4.17.21", []string{"4.17.19", "4.17.21"}, false},
		{"lower major is breaking", "4.0.0", []string{"4.0.0", "5.0.0"}, true},
		{"single version", "4.0.0", []string{"4.0.0"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isBreakingVersion(tt.version, tt.allVersions)
			if result != tt.want {
				t.Errorf("isBreakingVersion(%s, %v) = %v, want %v",
					tt.version, tt.allVersions, result, tt.want)
			}
		})
	}
}

