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

func TestCircularDependencyJSON(t *testing.T) {
	c := CircularDependency{
		Nodes: []string{"A", "B", "C", "A"},
		Depth: 3,
	}

	jsonBytes, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("Failed to marshal CircularDependency: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if _, ok := parsed["nodes"]; !ok {
		t.Error("Missing 'nodes' key")
	}
	if _, ok := parsed["depth"]; !ok {
		t.Error("Missing 'depth' key")
	}
}

func TestVersionConflictJSON(t *testing.T) {
	v := VersionConflict{
		PackageName: "lodash",
		Versions: map[string]string{
			"@app/web":    "4.17.21",
			"@app/shared": "4.17.15",
		},
	}

	jsonBytes, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Failed to marshal VersionConflict: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify camelCase - packageName not package_name
	if _, ok := parsed["packageName"]; !ok {
		t.Error("Missing 'packageName' key (should be camelCase)")
	}
	if _, ok := parsed["package_name"]; ok {
		t.Error("Found 'package_name' key - should be camelCase 'packageName'")
	}
}
