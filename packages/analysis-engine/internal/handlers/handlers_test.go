package handlers

import (
	"encoding/json"
	"testing"
)

func TestGetVersion(t *testing.T) {
	result := GetVersion()

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Failed to parse JSON result: %v", err)
	}

	// Verify Result<T> structure
	if parsed["error"] != nil {
		t.Errorf("GetVersion returned error: %v", parsed["error"])
	}

	data, ok := parsed["data"].(map[string]interface{})
	if !ok {
		t.Fatal("GetVersion data is not a map")
	}

	version, ok := data["version"].(string)
	if !ok {
		t.Fatal("version is not a string")
	}

	if version != Version {
		t.Errorf("version = %q, want %q", version, Version)
	}
}

func TestAnalyze(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
		errorCode string
	}{
		{
			name:      "empty string input returns error",
			input:     "",
			wantError: true,
			errorCode: "INVALID_INPUT",
		},
		{
			name:      "invalid JSON input returns error",
			input:     "{invalid json",
			wantError: true,
			errorCode: "INVALID_INPUT",
		},
		{
			name:      "empty files object returns error (missing package.json)",
			input:     "{}",
			wantError: true,
			errorCode: "ANALYSIS_FAILED",
		},
		{
			name:      "missing root package.json returns error",
			input:     `{"packages/pkg-a/package.json": "{\"name\": \"pkg-a\"}"}`,
			wantError: true,
			errorCode: "ANALYSIS_FAILED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Analyze(tt.input)

			var parsed map[string]interface{}
			if err := json.Unmarshal([]byte(result), &parsed); err != nil {
				t.Fatalf("Failed to parse JSON result: %v", err)
			}

			if tt.wantError {
				if parsed["error"] == nil {
					t.Error("Expected error but got success")
					return
				}
				errObj := parsed["error"].(map[string]interface{})
				if errObj["code"] != tt.errorCode {
					t.Errorf("error code = %v, want %v", errObj["code"], tt.errorCode)
				}
			} else {
				if parsed["error"] != nil {
					t.Errorf("Unexpected error: %v", parsed["error"])
				}
			}
		})
	}
}

func TestAnalyzeWithRealWorkspace(t *testing.T) {
	// Test with a real npm workspace configuration
	input := `{
		"package.json": "{\"name\": \"monorepo-root\", \"workspaces\": [\"packages/*\"]}",
		"package-lock.json": "{}",
		"packages/pkg-a/package.json": "{\"name\": \"@mono/pkg-a\", \"version\": \"1.0.0\", \"dependencies\": {\"@mono/pkg-b\": \"^1.0.0\"}}",
		"packages/pkg-b/package.json": "{\"name\": \"@mono/pkg-b\", \"version\": \"1.0.0\"}"
	}`

	result := Analyze(input)

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Failed to parse JSON result: %v", err)
	}

	// Verify no error
	if parsed["error"] != nil {
		t.Fatalf("Unexpected error: %v", parsed["error"])
	}

	// Verify WorkspaceData structure
	data, ok := parsed["data"].(map[string]interface{})
	if !ok {
		t.Fatal("data is not a map")
	}

	// Verify rootPath
	if data["rootPath"] != "/workspace" {
		t.Errorf("rootPath = %v, want /workspace", data["rootPath"])
	}

	// Verify workspaceType
	if data["workspaceType"] != "npm" {
		t.Errorf("workspaceType = %v, want npm", data["workspaceType"])
	}

	// Verify packages
	packages, ok := data["packages"].(map[string]interface{})
	if !ok {
		t.Fatal("packages is not a map")
	}

	if len(packages) != 2 {
		t.Errorf("packages count = %d, want 2", len(packages))
	}

	// Verify specific package
	pkgA, ok := packages["@mono/pkg-a"].(map[string]interface{})
	if !ok {
		t.Fatal("@mono/pkg-a not found or not a map")
	}

	if pkgA["name"] != "@mono/pkg-a" {
		t.Errorf("pkg-a name = %v, want @mono/pkg-a", pkgA["name"])
	}
	if pkgA["version"] != "1.0.0" {
		t.Errorf("pkg-a version = %v, want 1.0.0", pkgA["version"])
	}

	// Verify dependencies are parsed
	deps, ok := pkgA["dependencies"].(map[string]interface{})
	if !ok {
		t.Fatal("pkg-a dependencies not found or not a map")
	}
	if deps["@mono/pkg-b"] != "^1.0.0" {
		t.Errorf("pkg-a dependency @mono/pkg-b = %v, want ^1.0.0", deps["@mono/pkg-b"])
	}
}

func TestAnalyzeWithPnpmWorkspace(t *testing.T) {
	// Test with a pnpm workspace configuration
	input := `{
		"pnpm-workspace.yaml": "packages:\n  - 'packages/*'",
		"package.json": "{\"name\": \"monorepo-root\"}",
		"packages/core/package.json": "{\"name\": \"@mono/core\", \"version\": \"2.0.0\"}"
	}`

	result := Analyze(input)

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Failed to parse JSON result: %v", err)
	}

	if parsed["error"] != nil {
		t.Fatalf("Unexpected error: %v", parsed["error"])
	}

	data := parsed["data"].(map[string]interface{})

	// Verify pnpm workspace type
	if data["workspaceType"] != "pnpm" {
		t.Errorf("workspaceType = %v, want pnpm", data["workspaceType"])
	}

	// Verify package was parsed
	packages := data["packages"].(map[string]interface{})
	if len(packages) != 1 {
		t.Errorf("packages count = %d, want 1", len(packages))
	}
}

func TestCheck(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
		errorCode string
	}{
		{
			name:      "valid empty JSON input",
			input:     "{}",
			wantError: false,
		},
		{
			name:      "valid config JSON",
			input:     `{"rules": []}`,
			wantError: false,
		},
		{
			name:      "empty string input returns error",
			input:     "",
			wantError: true,
			errorCode: "INVALID_INPUT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Check(tt.input)

			var parsed map[string]interface{}
			if err := json.Unmarshal([]byte(result), &parsed); err != nil {
				t.Fatalf("Failed to parse JSON result: %v", err)
			}

			if tt.wantError {
				if parsed["error"] == nil {
					t.Error("Expected error but got success")
					return
				}
				errObj := parsed["error"].(map[string]interface{})
				if errObj["code"] != tt.errorCode {
					t.Errorf("error code = %v, want %v", errObj["code"], tt.errorCode)
				}
			} else {
				if parsed["error"] != nil {
					t.Errorf("Unexpected error: %v", parsed["error"])
					return
				}

				data, ok := parsed["data"].(map[string]interface{})
				if !ok {
					t.Fatal("data is not a map")
				}

				// Verify CheckResult fields
				if data["passed"] != true {
					t.Error("Expected passed to be true")
				}
				if _, ok := data["errors"]; !ok {
					t.Error("Missing errors field")
				}
				if data["placeholder"] != true {
					t.Error("Expected placeholder to be true")
				}
			}
		})
	}
}

func TestResultStructure(t *testing.T) {
	// Verify all handlers return proper Result<T> structure
	// Note: Analyze needs a valid workspace to succeed
	validWorkspace := `{
		"package.json": "{\"name\": \"root\", \"workspaces\": [\"packages/*\"]}",
		"package-lock.json": "{}",
		"packages/a/package.json": "{\"name\": \"a\", \"version\": \"1.0.0\"}"
	}`

	handlers := []struct {
		name   string
		result string
	}{
		{"GetVersion", GetVersion()},
		{"Analyze", Analyze(validWorkspace)},
		{"Check", Check("{}")},
	}

	for _, h := range handlers {
		t.Run(h.name+" returns Result structure", func(t *testing.T) {
			var parsed map[string]interface{}
			if err := json.Unmarshal([]byte(h.result), &parsed); err != nil {
				t.Fatalf("Failed to parse JSON: %v", err)
			}

			// Result<T> must have both "data" and "error" keys
			if _, hasData := parsed["data"]; !hasData {
				t.Error("Missing 'data' key in Result")
			}
			if _, hasError := parsed["error"]; !hasError {
				t.Error("Missing 'error' key in Result")
			}
		})
	}
}

func TestErrorResultStructure(t *testing.T) {
	// Verify error results have proper structure
	result := Analyze("") // Empty input triggers error

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if parsed["data"] != nil {
		t.Error("Error result should have null data")
	}

	errObj, ok := parsed["error"].(map[string]interface{})
	if !ok {
		t.Fatal("Error result missing error object")
	}

	// Error must have code and message
	if _, ok := errObj["code"]; !ok {
		t.Error("Error missing 'code' field")
	}
	if _, ok := errObj["message"]; !ok {
		t.Error("Error missing 'message' field")
	}

	// Verify UPPER_SNAKE_CASE error code
	code := errObj["code"].(string)
	if code != "INVALID_INPUT" {
		t.Errorf("error code = %v, want INVALID_INPUT", code)
	}
}
