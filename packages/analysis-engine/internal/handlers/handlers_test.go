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

	// Verify AnalysisResult structure
	data, ok := parsed["data"].(map[string]interface{})
	if !ok {
		t.Fatal("data is not a map")
	}

	// Verify healthScore
	healthScore, ok := data["healthScore"].(float64)
	if !ok {
		t.Fatal("healthScore is not a number")
	}
	if healthScore != 100 {
		t.Errorf("healthScore = %v, want 100 (placeholder)", healthScore)
	}

	// Verify packages count
	packages, ok := data["packages"].(float64)
	if !ok {
		t.Fatal("packages is not a number")
	}
	if packages != 2 {
		t.Errorf("packages = %v, want 2", packages)
	}

	// Verify graph exists and has correct structure
	graph, ok := data["graph"].(map[string]interface{})
	if !ok {
		t.Fatal("graph is not a map")
	}

	// Verify graph rootPath
	if graph["rootPath"] != "/workspace" {
		t.Errorf("graph.rootPath = %v, want /workspace", graph["rootPath"])
	}

	// Verify graph workspaceType
	if graph["workspaceType"] != "npm" {
		t.Errorf("graph.workspaceType = %v, want npm", graph["workspaceType"])
	}

	// Verify graph nodes
	nodes, ok := graph["nodes"].(map[string]interface{})
	if !ok {
		t.Fatal("graph.nodes is not a map")
	}
	if len(nodes) != 2 {
		t.Errorf("graph.nodes count = %d, want 2", len(nodes))
	}

	// Verify specific node
	pkgA, ok := nodes["@mono/pkg-a"].(map[string]interface{})
	if !ok {
		t.Fatal("@mono/pkg-a node not found or not a map")
	}
	if pkgA["name"] != "@mono/pkg-a" {
		t.Errorf("pkg-a name = %v, want @mono/pkg-a", pkgA["name"])
	}

	// Verify graph edges (should have 1 edge: pkg-a -> pkg-b)
	edges, ok := graph["edges"].([]interface{})
	if !ok {
		t.Fatal("graph.edges is not an array")
	}
	if len(edges) != 1 {
		t.Errorf("graph.edges count = %d, want 1", len(edges))
	}

	// Verify edge details
	if len(edges) > 0 {
		edge := edges[0].(map[string]interface{})
		if edge["from"] != "@mono/pkg-a" {
			t.Errorf("edge.from = %v, want @mono/pkg-a", edge["from"])
		}
		if edge["to"] != "@mono/pkg-b" {
			t.Errorf("edge.to = %v, want @mono/pkg-b", edge["to"])
		}
		if edge["type"] != "production" {
			t.Errorf("edge.type = %v, want production", edge["type"])
		}
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

	// Verify packages count (AnalysisResult)
	packages, ok := data["packages"].(float64)
	if !ok {
		t.Fatal("packages is not a number")
	}
	if packages != 1 {
		t.Errorf("packages = %v, want 1", packages)
	}

	// Verify graph has pnpm workspace type
	graph := data["graph"].(map[string]interface{})
	if graph["workspaceType"] != "pnpm" {
		t.Errorf("graph.workspaceType = %v, want pnpm", graph["workspaceType"])
	}

	// Verify node was created
	nodes := graph["nodes"].(map[string]interface{})
	if len(nodes) != 1 {
		t.Errorf("graph.nodes count = %d, want 1", len(nodes))
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

// TestAnalyzeWithVersionConflicts verifies version conflict detection (Story 2.4).
func TestAnalyzeWithVersionConflicts(t *testing.T) {
	// Test with external dependencies that have version conflicts
	input := `{
		"package.json": "{\"name\": \"monorepo-root\", \"workspaces\": [\"packages/*\"]}",
		"package-lock.json": "{}",
		"packages/app/package.json": "{\"name\": \"@mono/app\", \"version\": \"1.0.0\", \"dependencies\": {\"lodash\": \"^4.17.21\", \"typescript\": \"^5.0.0\"}}",
		"packages/lib/package.json": "{\"name\": \"@mono/lib\", \"version\": \"1.0.0\", \"dependencies\": {\"lodash\": \"^4.17.19\", \"typescript\": \"^4.9.0\"}}"
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

	data, ok := parsed["data"].(map[string]interface{})
	if !ok {
		t.Fatal("data is not a map")
	}

	// Verify versionConflicts field exists
	versionConflicts, ok := data["versionConflicts"].([]interface{})
	if !ok {
		t.Fatal("versionConflicts is not an array or missing")
	}

	// Should have 2 conflicts: lodash (patch diff) and typescript (major diff)
	if len(versionConflicts) != 2 {
		t.Errorf("versionConflicts count = %d, want 2", len(versionConflicts))
	}

	// Find typescript conflict and verify critical severity
	for _, c := range versionConflicts {
		conflict := c.(map[string]interface{})
		if conflict["packageName"] == "typescript" {
			if conflict["severity"] != "critical" {
				t.Errorf("typescript severity = %v, want critical", conflict["severity"])
			}

			// Verify conflictingVersions structure
			versions, ok := conflict["conflictingVersions"].([]interface{})
			if !ok {
				t.Fatal("conflictingVersions is not an array")
			}
			if len(versions) != 2 {
				t.Errorf("typescript conflictingVersions count = %d, want 2", len(versions))
			}

			// Verify resolution and impact exist
			if conflict["resolution"] == nil || conflict["resolution"] == "" {
				t.Error("typescript conflict missing resolution")
			}
			if conflict["impact"] == nil || conflict["impact"] == "" {
				t.Error("typescript conflict missing impact")
			}
		}

		if conflict["packageName"] == "lodash" {
			if conflict["severity"] != "info" {
				t.Errorf("lodash severity = %v, want info", conflict["severity"])
			}
		}
	}
}

// TestAnalyzeNoVersionConflictsInResult verifies no conflicts when versions match.
func TestAnalyzeNoVersionConflictsInResult(t *testing.T) {
	input := `{
		"package.json": "{\"name\": \"monorepo-root\", \"workspaces\": [\"packages/*\"]}",
		"package-lock.json": "{}",
		"packages/app/package.json": "{\"name\": \"@mono/app\", \"version\": \"1.0.0\", \"dependencies\": {\"lodash\": \"^4.17.21\"}}",
		"packages/lib/package.json": "{\"name\": \"@mono/lib\", \"version\": \"1.0.0\", \"dependencies\": {\"lodash\": \"^4.17.21\"}}"
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

	// versionConflicts should be nil or empty when no conflicts
	versionConflicts := data["versionConflicts"]
	if versionConflicts != nil {
		conflicts := versionConflicts.([]interface{})
		if len(conflicts) != 0 {
			t.Errorf("Expected no version conflicts, got %d", len(conflicts))
		}
	}
}

// ========================================
// Story 2.6: Exclusion Pattern Tests
// ========================================

// TestAnalyzeWithExclusionConfig verifies exclusion config is accepted.
func TestAnalyzeWithExclusionConfig(t *testing.T) {
	// New AnalysisInput format with config
	input := `{
		"files": {
			"package.json": "{\"name\": \"monorepo-root\", \"workspaces\": [\"packages/*\"]}",
			"package-lock.json": "{}",
			"packages/app/package.json": "{\"name\": \"@mono/app\", \"version\": \"1.0.0\"}",
			"packages/legacy/package.json": "{\"name\": \"@mono/legacy\", \"version\": \"1.0.0\"}"
		},
		"config": {
			"exclude": ["@mono/legacy"]
		}
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

	data, ok := parsed["data"].(map[string]interface{})
	if !ok {
		t.Fatal("data is not a map")
	}

	// Verify excludedPackages count
	excludedPackages, ok := data["excludedPackages"].(float64)
	if !ok {
		t.Fatal("excludedPackages is not a number")
	}
	if excludedPackages != 1 {
		t.Errorf("excludedPackages = %v, want 1", excludedPackages)
	}

	// Verify packages count (excludes legacy)
	packages, ok := data["packages"].(float64)
	if !ok {
		t.Fatal("packages is not a number")
	}
	if packages != 1 {
		t.Errorf("packages = %v, want 1 (excluding legacy)", packages)
	}

	// Verify legacy is in graph with excluded flag
	graph := data["graph"].(map[string]interface{})
	nodes := graph["nodes"].(map[string]interface{})

	legacyNode, ok := nodes["@mono/legacy"].(map[string]interface{})
	if !ok {
		t.Fatal("legacy node not found in graph")
	}
	if legacyNode["excluded"] != true {
		t.Error("legacy node should have excluded=true")
	}

	// Verify app is not excluded
	appNode, ok := nodes["@mono/app"].(map[string]interface{})
	if !ok {
		t.Fatal("app node not found in graph")
	}
	// excluded should be false or omitted (omitempty)
	if appNode["excluded"] == true {
		t.Error("app node should not be excluded")
	}
}

// TestAnalyzeLegacyFormatStillWorks verifies old format still works.
func TestAnalyzeLegacyFormatStillWorks(t *testing.T) {
	// Old format: just map of files (no config wrapper)
	input := `{
		"package.json": "{\"name\": \"monorepo-root\", \"workspaces\": [\"packages/*\"]}",
		"package-lock.json": "{}",
		"packages/app/package.json": "{\"name\": \"@mono/app\", \"version\": \"1.0.0\"}"
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

	data, ok := parsed["data"].(map[string]interface{})
	if !ok {
		t.Fatal("data is not a map")
	}

	// Verify packages count
	packages, ok := data["packages"].(float64)
	if !ok {
		t.Fatal("packages is not a number")
	}
	if packages != 1 {
		t.Errorf("packages = %v, want 1", packages)
	}
}

// TestAnalyzeWithGlobExclusion verifies glob pattern exclusion.
func TestAnalyzeWithGlobExclusion(t *testing.T) {
	input := `{
		"files": {
			"package.json": "{\"name\": \"monorepo-root\", \"workspaces\": [\"packages/*\"]}",
			"package-lock.json": "{}",
			"packages/app/package.json": "{\"name\": \"@mono/app\", \"version\": \"1.0.0\"}",
			"packages/deprecated-utils/package.json": "{\"name\": \"@mono/deprecated-utils\", \"version\": \"1.0.0\"}",
			"packages/deprecated-api/package.json": "{\"name\": \"@mono/deprecated-api\", \"version\": \"1.0.0\"}"
		},
		"config": {
			"exclude": ["@mono/deprecated-*"]
		}
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

	// 2 packages should be excluded
	excludedPackages := data["excludedPackages"].(float64)
	if excludedPackages != 2 {
		t.Errorf("excludedPackages = %v, want 2", excludedPackages)
	}

	// Only app should be counted
	packages := data["packages"].(float64)
	if packages != 1 {
		t.Errorf("packages = %v, want 1", packages)
	}
}

// TestAnalyzeWithRegexExclusion verifies regex pattern exclusion.
func TestAnalyzeWithRegexExclusion(t *testing.T) {
	input := `{
		"files": {
			"package.json": "{\"name\": \"monorepo-root\", \"workspaces\": [\"packages/*\"]}",
			"package-lock.json": "{}",
			"packages/app/package.json": "{\"name\": \"@mono/app\", \"version\": \"1.0.0\"}",
			"packages/core-test/package.json": "{\"name\": \"@mono/core-test\", \"version\": \"1.0.0\"}"
		},
		"config": {
			"exclude": ["regex:.*-test$"]
		}
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

	// core-test should be excluded
	excludedPackages := data["excludedPackages"].(float64)
	if excludedPackages != 1 {
		t.Errorf("excludedPackages = %v, want 1", excludedPackages)
	}
}

// TestAnalyzeWithInvalidRegexPattern verifies error on invalid regex.
func TestAnalyzeWithInvalidRegexPattern(t *testing.T) {
	input := `{
		"files": {
			"package.json": "{\"name\": \"root\", \"workspaces\": [\"packages/*\"]}",
			"package-lock.json": "{}"
		},
		"config": {
			"exclude": ["regex:[invalid"]
		}
	}`

	result := Analyze(input)

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("Failed to parse JSON result: %v", err)
	}

	// Should return error for invalid regex
	if parsed["error"] == nil {
		t.Error("Expected error for invalid regex pattern")
		return
	}

	errObj := parsed["error"].(map[string]interface{})
	if errObj["code"] != "INVALID_INPUT" {
		t.Errorf("error code = %v, want INVALID_INPUT", errObj["code"])
	}
}

// TestAnalyzeExcludedFromConflictDetection verifies excluded packages don't affect conflicts.
func TestAnalyzeExcludedFromConflictDetection(t *testing.T) {
	input := `{
		"files": {
			"package.json": "{\"name\": \"monorepo-root\", \"workspaces\": [\"packages/*\"]}",
			"package-lock.json": "{}",
			"packages/app/package.json": "{\"name\": \"@mono/app\", \"version\": \"1.0.0\", \"dependencies\": {\"typescript\": \"^5.0.0\"}}",
			"packages/legacy/package.json": "{\"name\": \"@mono/legacy\", \"version\": \"1.0.0\", \"dependencies\": {\"typescript\": \"^4.0.0\"}}"
		},
		"config": {
			"exclude": ["@mono/legacy"]
		}
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

	// No conflicts should be detected because legacy is excluded
	versionConflicts := data["versionConflicts"]
	if versionConflicts != nil {
		conflicts := versionConflicts.([]interface{})
		if len(conflicts) != 0 {
			t.Errorf("Expected no version conflicts (legacy excluded), got %d", len(conflicts))
		}
	}
}

// ========================================
// Story 3.2: Import Tracing Tests
// ========================================

// TestAnalyzeWithSourceFiles verifies import tracing when source files provided.
func TestAnalyzeWithSourceFiles(t *testing.T) {
	// Create workspace with cycle: ui → api → ui
	input := `{
		"files": {
			"package.json": "{\"name\": \"monorepo-root\", \"workspaces\": [\"packages/*\"]}",
			"package-lock.json": "{}",
			"packages/ui/package.json": "{\"name\": \"@mono/ui\", \"version\": \"1.0.0\", \"dependencies\": {\"@mono/api\": \"^1.0.0\"}}",
			"packages/api/package.json": "{\"name\": \"@mono/api\", \"version\": \"1.0.0\", \"dependencies\": {\"@mono/ui\": \"^1.0.0\"}}"
		},
		"sourceFiles": {
			"packages/ui/src/index.ts": "import { api } from '@mono/api';",
			"packages/api/src/index.ts": "import { ui } from '@mono/ui';"
		}
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

	data, ok := parsed["data"].(map[string]interface{})
	if !ok {
		t.Fatal("data is not a map")
	}

	// Verify circular dependency detected
	circularDeps, ok := data["circularDependencies"].([]interface{})
	if !ok || len(circularDeps) == 0 {
		t.Fatal("circularDependencies not found or empty")
	}

	cycle := circularDeps[0].(map[string]interface{})

	// Verify importTraces is populated (Story 3.2)
	importTraces, ok := cycle["importTraces"].([]interface{})
	if !ok {
		t.Fatal("importTraces not found or not an array")
	}
	if len(importTraces) == 0 {
		t.Error("importTraces should be populated when source files provided")
	}

	// Verify trace structure
	for _, trace := range importTraces {
		traceMap := trace.(map[string]interface{})
		if traceMap["fromPackage"] == nil || traceMap["fromPackage"] == "" {
			t.Error("importTrace.fromPackage should not be empty")
		}
		if traceMap["toPackage"] == nil || traceMap["toPackage"] == "" {
			t.Error("importTrace.toPackage should not be empty")
		}
		if traceMap["filePath"] == nil || traceMap["filePath"] == "" {
			t.Error("importTrace.filePath should not be empty")
		}
		if traceMap["lineNumber"] == nil {
			t.Error("importTrace.lineNumber should not be nil")
		}
		if traceMap["statement"] == nil || traceMap["statement"] == "" {
			t.Error("importTrace.statement should not be empty")
		}
		if traceMap["importType"] == nil || traceMap["importType"] == "" {
			t.Error("importTrace.importType should not be empty")
		}
	}
}

// TestAnalyzeWithoutSourceFiles verifies graceful degradation without source files.
func TestAnalyzeWithoutSourceFiles(t *testing.T) {
	// Create workspace with cycle but NO source files
	input := `{
		"files": {
			"package.json": "{\"name\": \"monorepo-root\", \"workspaces\": [\"packages/*\"]}",
			"package-lock.json": "{}",
			"packages/ui/package.json": "{\"name\": \"@mono/ui\", \"version\": \"1.0.0\", \"dependencies\": {\"@mono/api\": \"^1.0.0\"}}",
			"packages/api/package.json": "{\"name\": \"@mono/api\", \"version\": \"1.0.0\", \"dependencies\": {\"@mono/ui\": \"^1.0.0\"}}"
		}
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

	data := parsed["data"].(map[string]interface{})

	// Verify circular dependency detected
	circularDeps, ok := data["circularDependencies"].([]interface{})
	if !ok || len(circularDeps) == 0 {
		t.Fatal("circularDependencies not found or empty")
	}

	cycle := circularDeps[0].(map[string]interface{})

	// importTraces should be empty (graceful degradation) when no source files
	importTraces := cycle["importTraces"]
	if importTraces != nil {
		traces := importTraces.([]interface{})
		if len(traces) != 0 {
			t.Errorf("importTraces should be empty without source files, got %d", len(traces))
		}
	}
}

// TestAnalyzeWithSourceFilesBackwardCompatible verifies legacy format still works.
func TestAnalyzeWithSourceFilesBackwardCompatible(t *testing.T) {
	// Legacy format (flat map) should still work without source files
	input := `{
		"package.json": "{\"name\": \"monorepo-root\", \"workspaces\": [\"packages/*\"]}",
		"package-lock.json": "{}",
		"packages/app/package.json": "{\"name\": \"@mono/app\", \"version\": \"1.0.0\"}"
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

	data := parsed["data"].(map[string]interface{})

	// Basic analysis should succeed
	packages := data["packages"].(float64)
	if packages != 1 {
		t.Errorf("packages = %v, want 1", packages)
	}
}
