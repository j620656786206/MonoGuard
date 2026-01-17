package types

import (
	"encoding/json"
	"testing"
)

// TestAnalysisConfig_JSONSerialization verifies config JSON marshaling.
func TestAnalysisConfig_JSONSerialization(t *testing.T) {
	config := &AnalysisConfig{
		Exclude: []string{
			"packages/legacy",
			"packages/deprecated-*",
			"regex:^@mono/test-.*$",
		},
	}

	// Marshal
	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Verify JSON structure
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Unmarshal to map failed: %v", err)
	}

	exclude, ok := parsed["exclude"].([]interface{})
	if !ok {
		t.Error("exclude field missing or wrong type")
	}

	if len(exclude) != 3 {
		t.Errorf("exclude length = %d, want 3", len(exclude))
	}

	// Unmarshal back
	var unmarshaled AnalysisConfig
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Unmarshal back failed: %v", err)
	}

	if len(unmarshaled.Exclude) != 3 {
		t.Errorf("Exclude length = %d, want 3", len(unmarshaled.Exclude))
	}
}

// TestAnalysisConfig_EmptyExclude verifies empty exclusions.
func TestAnalysisConfig_EmptyExclude(t *testing.T) {
	config := &AnalysisConfig{
		Exclude: []string{},
	}

	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var unmarshaled AnalysisConfig
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(unmarshaled.Exclude) != 0 {
		t.Errorf("Exclude length = %d, want 0", len(unmarshaled.Exclude))
	}
}

// TestAnalysisConfig_OmitEmpty verifies omitempty behavior.
func TestAnalysisConfig_OmitEmpty(t *testing.T) {
	config := &AnalysisConfig{} // nil Exclude

	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// With omitempty, empty slice may be omitted
	jsonStr := string(data)
	t.Logf("Empty config JSON: %s", jsonStr)

	var unmarshaled AnalysisConfig
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
}

// TestAnalysisInput_JSONSerialization verifies input JSON structure.
func TestAnalysisInput_JSONSerialization(t *testing.T) {
	input := &AnalysisInput{
		Files: map[string]string{
			"package.json":              `{"name": "root"}`,
			"packages/app/package.json": `{"name": "@mono/app"}`,
		},
		Config: &AnalysisConfig{
			Exclude: []string{"packages/legacy"},
		},
	}

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Unmarshal to map failed: %v", err)
	}

	// Verify files field
	files, ok := parsed["files"].(map[string]interface{})
	if !ok {
		t.Error("files field missing or wrong type")
	}
	if len(files) != 2 {
		t.Errorf("files count = %d, want 2", len(files))
	}

	// Verify config field
	config, ok := parsed["config"].(map[string]interface{})
	if !ok {
		t.Error("config field missing or wrong type")
	}
	if config["exclude"] == nil {
		t.Error("config.exclude missing")
	}

	// Unmarshal back
	var unmarshaled AnalysisInput
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Unmarshal back failed: %v", err)
	}

	if len(unmarshaled.Files) != 2 {
		t.Errorf("Files count = %d, want 2", len(unmarshaled.Files))
	}
	if unmarshaled.Config == nil {
		t.Error("Config is nil")
	}
	if len(unmarshaled.Config.Exclude) != 1 {
		t.Errorf("Config.Exclude length = %d, want 1", len(unmarshaled.Config.Exclude))
	}
}

// TestAnalysisInput_NoConfig verifies optional config.
func TestAnalysisInput_NoConfig(t *testing.T) {
	input := &AnalysisInput{
		Files: map[string]string{
			"package.json": `{"name": "root"}`,
		},
		// No config
	}

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var unmarshaled AnalysisInput
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(unmarshaled.Files) != 1 {
		t.Errorf("Files count = %d, want 1", len(unmarshaled.Files))
	}
	if unmarshaled.Config != nil {
		t.Error("Config should be nil when not provided")
	}
}

// TestNewAnalysisConfig verifies constructor.
func TestNewAnalysisConfig(t *testing.T) {
	config := NewAnalysisConfig()
	if config == nil {
		t.Fatal("NewAnalysisConfig returned nil")
	}
	if config.Exclude == nil {
		t.Error("Exclude is nil, should be empty slice")
	}
	if len(config.Exclude) != 0 {
		t.Errorf("Exclude length = %d, want 0", len(config.Exclude))
	}
}
