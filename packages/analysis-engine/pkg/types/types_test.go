package types

import (
	"encoding/json"
	"testing"
)

func TestAnalysisResultJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    AnalysisResult
		wantKeys []string
	}{
		{
			name: "full result with camelCase JSON",
			input: AnalysisResult{
				HealthScore: 85,
				Packages:    10,
				CreatedAt:   "2026-01-16T10:30:00Z",
				Placeholder: false,
			},
			wantKeys: []string{"healthScore", "packages", "createdAt"},
		},
		{
			name: "placeholder result omits empty fields",
			input: AnalysisResult{
				HealthScore: 100,
				Packages:    0,
				Placeholder: true,
			},
			wantKeys: []string{"healthScore", "packages", "placeholder"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("Failed to marshal AnalysisResult: %v", err)
			}

			var parsed map[string]interface{}
			if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			// Verify expected keys exist
			for _, key := range tt.wantKeys {
				if _, ok := parsed[key]; !ok {
					t.Errorf("Missing expected key %q in JSON output", key)
				}
			}

			// Verify NO snake_case keys exist
			snakeCaseKeys := []string{"health_score", "created_at"}
			for _, key := range snakeCaseKeys {
				if _, ok := parsed[key]; ok {
					t.Errorf("Found snake_case key %q - should be camelCase", key)
				}
			}
		})
	}
}

func TestCheckResultJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    CheckResult
		wantKeys []string
	}{
		{
			name: "passed check with camelCase JSON",
			input: CheckResult{
				Passed: true,
				Errors: []string{},
			},
			wantKeys: []string{"passed", "errors"},
		},
		{
			name: "failed check with errors",
			input: CheckResult{
				Passed: false,
				Errors: []string{"circular dependency detected", "missing package"},
			},
			wantKeys: []string{"passed", "errors"},
		},
		{
			name: "placeholder check",
			input: CheckResult{
				Passed:      true,
				Errors:      []string{},
				Placeholder: true,
			},
			wantKeys: []string{"passed", "errors", "placeholder"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("Failed to marshal CheckResult: %v", err)
			}

			var parsed map[string]interface{}
			if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			for _, key := range tt.wantKeys {
				if _, ok := parsed[key]; !ok {
					t.Errorf("Missing expected key %q in JSON output", key)
				}
			}
		})
	}
}

func TestVersionInfoJSON(t *testing.T) {
	v := VersionInfo{Version: "0.1.0"}

	jsonBytes, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Failed to marshal VersionInfo: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if parsed["version"] != "0.1.0" {
		t.Errorf("version = %v, want 0.1.0", parsed["version"])
	}
}

func TestPackageJSON(t *testing.T) {
	p := Package{
		Name:         "@monoguard/types",
		Path:         "packages/types",
		Dependencies: []string{"typescript", "zod"},
	}

	jsonBytes, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Failed to marshal Package: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify camelCase keys
	expectedKeys := []string{"name", "path", "dependencies"}
	for _, key := range expectedKeys {
		if _, ok := parsed[key]; !ok {
			t.Errorf("Missing expected key %q", key)
		}
	}
}

// Note: Tests for CircularDependency and VersionConflict have been moved to:
// - circular_test.go (CircularDependencyInfo tests)
// - version_conflict_test.go (VersionConflictInfo tests)

// ========================================
// WorkspaceData Types Tests (Story 2.1)
// ========================================

func TestWorkspaceTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		wt       WorkspaceType
		expected string
	}{
		{"npm workspace type", WorkspaceTypeNpm, "npm"},
		{"yarn workspace type", WorkspaceTypeYarn, "yarn"},
		{"pnpm workspace type", WorkspaceTypePnpm, "pnpm"},
		{"unknown workspace type", WorkspaceTypeUnknown, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.wt) != tt.expected {
				t.Errorf("WorkspaceType = %q, want %q", tt.wt, tt.expected)
			}
		})
	}
}

func TestWorkspaceDataJSON(t *testing.T) {
	tests := []struct {
		name            string
		input           WorkspaceData
		wantKeys        []string
		wantSnakeCases  []string
		wantType        string
	}{
		{
			name: "npm workspace with packages",
			input: WorkspaceData{
				RootPath:      "/workspace",
				WorkspaceType: WorkspaceTypeNpm,
				Packages: map[string]*PackageInfo{
					"@mono/pkg-a": {
						Name:             "@mono/pkg-a",
						Version:          "1.0.0",
						Path:             "packages/pkg-a",
						Dependencies:     map[string]string{"@mono/pkg-b": "^1.0.0"},
						DevDependencies:  map[string]string{"typescript": "^5.0.0"},
						PeerDependencies: map[string]string{},
					},
				},
			},
			wantKeys:       []string{"rootPath", "workspaceType", "packages"},
			wantSnakeCases: []string{"root_path", "workspace_type"},
			wantType:       "npm",
		},
		{
			name: "pnpm workspace",
			input: WorkspaceData{
				RootPath:      "/workspace",
				WorkspaceType: WorkspaceTypePnpm,
				Packages:      map[string]*PackageInfo{},
			},
			wantKeys:       []string{"rootPath", "workspaceType", "packages"},
			wantSnakeCases: []string{"root_path", "workspace_type"},
			wantType:       "pnpm",
		},
		{
			name: "yarn workspace",
			input: WorkspaceData{
				RootPath:      "/workspace",
				WorkspaceType: WorkspaceTypeYarn,
				Packages:      map[string]*PackageInfo{},
			},
			wantKeys:       []string{"rootPath", "workspaceType", "packages"},
			wantSnakeCases: []string{"root_path", "workspace_type"},
			wantType:       "yarn",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("Failed to marshal WorkspaceData: %v", err)
			}

			var parsed map[string]interface{}
			if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			// Verify expected camelCase keys exist
			for _, key := range tt.wantKeys {
				if _, ok := parsed[key]; !ok {
					t.Errorf("Missing expected key %q in JSON output", key)
				}
			}

			// Verify NO snake_case keys exist
			for _, key := range tt.wantSnakeCases {
				if _, ok := parsed[key]; ok {
					t.Errorf("Found snake_case key %q - should be camelCase", key)
				}
			}

			// Verify workspaceType value
			if parsed["workspaceType"] != tt.wantType {
				t.Errorf("workspaceType = %v, want %v", parsed["workspaceType"], tt.wantType)
			}
		})
	}
}

func TestPackageInfoJSON(t *testing.T) {
	tests := []struct {
		name           string
		input          PackageInfo
		wantKeys       []string
		wantSnakeCases []string
	}{
		{
			name: "full package info with all dependency types",
			input: PackageInfo{
				Name:             "@mono/pkg-a",
				Version:          "1.0.0",
				Path:             "packages/pkg-a",
				Dependencies:     map[string]string{"lodash": "^4.17.21"},
				DevDependencies:  map[string]string{"typescript": "^5.0.0", "vitest": "^1.0.0"},
				PeerDependencies: map[string]string{"react": "^18.0.0"},
			},
			wantKeys:       []string{"name", "version", "path", "dependencies", "devDependencies", "peerDependencies"},
			wantSnakeCases: []string{"dev_dependencies", "peer_dependencies"},
		},
		{
			name: "minimal package info",
			input: PackageInfo{
				Name:             "@mono/simple",
				Version:          "0.0.1",
				Path:             "packages/simple",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{},
				PeerDependencies: map[string]string{},
			},
			wantKeys:       []string{"name", "version", "path", "dependencies", "devDependencies", "peerDependencies"},
			wantSnakeCases: []string{"dev_dependencies", "peer_dependencies"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("Failed to marshal PackageInfo: %v", err)
			}

			var parsed map[string]interface{}
			if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			// Verify expected camelCase keys exist
			for _, key := range tt.wantKeys {
				if _, ok := parsed[key]; !ok {
					t.Errorf("Missing expected key %q in JSON output", key)
				}
			}

			// Verify NO snake_case keys exist
			for _, key := range tt.wantSnakeCases {
				if _, ok := parsed[key]; ok {
					t.Errorf("Found snake_case key %q - should be camelCase", key)
				}
			}
		})
	}
}

func TestWorkspaceDataWithNestedPackages(t *testing.T) {
	// Test that nested package data serializes correctly
	ws := WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: WorkspaceTypeNpm,
		Packages: map[string]*PackageInfo{
			"@mono/pkg-a": {
				Name:    "@mono/pkg-a",
				Version: "1.0.0",
				Path:    "packages/pkg-a",
				Dependencies: map[string]string{
					"@mono/pkg-b": "^1.0.0",
					"lodash":      "^4.17.21",
				},
				DevDependencies:  map[string]string{"typescript": "^5.0.0"},
				PeerDependencies: map[string]string{},
			},
			"@mono/pkg-b": {
				Name:             "@mono/pkg-b",
				Version:          "1.0.0",
				Path:             "packages/pkg-b",
				Dependencies:     map[string]string{},
				DevDependencies:  map[string]string{"typescript": "^5.0.0"},
				PeerDependencies: map[string]string{},
			},
		},
	}

	jsonBytes, err := json.Marshal(ws)
	if err != nil {
		t.Fatalf("Failed to marshal WorkspaceData: %v", err)
	}

	// Unmarshal back to verify round-trip
	var parsed WorkspaceData
	if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal WorkspaceData: %v", err)
	}

	// Verify package count
	if len(parsed.Packages) != 2 {
		t.Errorf("Expected 2 packages, got %d", len(parsed.Packages))
	}

	// Verify specific package data
	pkgA, ok := parsed.Packages["@mono/pkg-a"]
	if !ok {
		t.Fatal("Missing @mono/pkg-a package")
	}
	if pkgA.Name != "@mono/pkg-a" {
		t.Errorf("Package name = %q, want %q", pkgA.Name, "@mono/pkg-a")
	}
	if len(pkgA.Dependencies) != 2 {
		t.Errorf("Expected 2 dependencies, got %d", len(pkgA.Dependencies))
	}
	if pkgA.Dependencies["@mono/pkg-b"] != "^1.0.0" {
		t.Errorf("Dependency version = %q, want %q", pkgA.Dependencies["@mono/pkg-b"], "^1.0.0")
	}
}

// ========================================
// FixSummary Integration Tests (Story 3.8)
// ========================================

func TestAnalysisResult_WithFixSummary(t *testing.T) {
	// Test that FixSummary field is optional and omitted when nil
	result := AnalysisResult{
		HealthScore: 85,
		Packages:    10,
		FixSummary:  nil, // Should be omitted in JSON
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal AnalysisResult: %v", err)
	}

	jsonStr := string(jsonBytes)

	// FixSummary should NOT be in JSON when nil (omitempty)
	if containsTypesStr(jsonStr, `"fixSummary"`) {
		t.Errorf("Expected fixSummary to be omitted when nil, got: %s", jsonStr)
	}
}

func TestAnalysisResult_WithFixSummaryPresent(t *testing.T) {
	// Test that FixSummary field is included when present
	result := AnalysisResult{
		HealthScore: 75,
		Packages:    15,
		FixSummary: &FixSummary{
			TotalCircularDependencies: 3,
			TotalEstimatedFixTime:     "2 hours",
			QuickWinsCount:            2,
			CriticalCyclesCount:       1,
			HighPriorityCycles: []PriorityCycleSummary{
				{
					CycleID:          "core→ui",
					PackagesInvolved: []string{"@mono/core", "@mono/ui"},
					PriorityScore:    80.0,
					RecommendedFix:   FixStrategyExtractModule,
					EstimatedTime:    "30-60 minutes",
				},
			},
		},
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal AnalysisResult: %v", err)
	}

	jsonStr := string(jsonBytes)

	// FixSummary should be in JSON when present
	expectedFields := []string{
		`"fixSummary"`,
		`"totalCircularDependencies"`,
		`"totalEstimatedFixTime"`,
		`"quickWinsCount"`,
		`"criticalCyclesCount"`,
		`"highPriorityCycles"`,
	}

	for _, field := range expectedFields {
		if !containsTypesStr(jsonStr, field) {
			t.Errorf("Expected JSON to contain %s, got: %s", field, jsonStr)
		}
	}

	// Verify round-trip
	var decoded AnalysisResult
	if err := json.Unmarshal(jsonBytes, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if decoded.FixSummary == nil {
		t.Fatal("FixSummary should not be nil after round-trip")
	}
	if decoded.FixSummary.TotalCircularDependencies != 3 {
		t.Errorf("TotalCircularDependencies = %d, want 3", decoded.FixSummary.TotalCircularDependencies)
	}
	if decoded.FixSummary.QuickWinsCount != 2 {
		t.Errorf("QuickWinsCount = %d, want 2", decoded.FixSummary.QuickWinsCount)
	}
}

func TestAnalysisResult_FixSummaryBackwardCompatibility(t *testing.T) {
	// Test that existing JSON without fixSummary still deserializes correctly
	jsonStr := `{
		"healthScore": 85,
		"packages": 10
	}`

	var decoded AnalysisResult
	if err := json.Unmarshal([]byte(jsonStr), &decoded); err != nil {
		t.Fatalf("Failed to unmarshal legacy JSON: %v", err)
	}

	// FixSummary should be nil for legacy JSON
	if decoded.FixSummary != nil {
		t.Error("FixSummary should be nil for legacy JSON without fixSummary field")
	}

	// Other fields should be correct
	if decoded.HealthScore != 85 {
		t.Errorf("HealthScore = %d, want 85", decoded.HealthScore)
	}
	if decoded.Packages != 10 {
		t.Errorf("Packages = %d, want 10", decoded.Packages)
	}
}

func containsTypesStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
