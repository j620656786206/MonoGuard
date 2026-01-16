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
			name:      "valid empty JSON input",
			input:     "{}",
			wantError: false,
		},
		{
			name:      "valid workspace JSON",
			input:     `{"projects": {}}`,
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
					return
				}

				data, ok := parsed["data"].(map[string]interface{})
				if !ok {
					t.Fatal("data is not a map")
				}

				// Verify AnalysisResult fields
				if _, ok := data["healthScore"]; !ok {
					t.Error("Missing healthScore field")
				}
				if _, ok := data["packages"]; !ok {
					t.Error("Missing packages field")
				}
				if data["placeholder"] != true {
					t.Error("Expected placeholder to be true")
				}
			}
		})
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
	handlers := []struct {
		name   string
		result string
	}{
		{"GetVersion", GetVersion()},
		{"Analyze", Analyze("{}")},
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
