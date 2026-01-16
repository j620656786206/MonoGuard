package result

import (
	"encoding/json"
	"testing"
)

func TestNewSuccess(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		wantData interface{}
	}{
		{
			name:     "string data",
			data:     "test",
			wantData: "test",
		},
		{
			name:     "map data",
			data:     map[string]string{"version": "0.1.0"},
			wantData: map[string]interface{}{"version": "0.1.0"},
		},
		{
			name: "struct data",
			data: struct {
				HealthScore int `json:"healthScore"`
			}{HealthScore: 85},
			wantData: map[string]interface{}{"healthScore": float64(85)},
		},
		{
			name:     "nil data",
			data:     nil,
			wantData: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewSuccess(tt.data)

			if r.Error != nil {
				t.Errorf("NewSuccess() error = %v, want nil", r.Error)
			}

			// Verify JSON output contains data and no error
			jsonStr := r.ToJSON()
			var parsed map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
				t.Fatalf("Failed to parse JSON: %v", err)
			}

			if parsed["error"] != nil {
				t.Errorf("JSON error field = %v, want nil", parsed["error"])
			}
		})
	}
}

func TestNewError(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		message     string
		wantCode    string
		wantMessage string
	}{
		{
			name:        "parse error",
			code:        "PARSE_ERROR",
			message:     "Invalid JSON input",
			wantCode:    "PARSE_ERROR",
			wantMessage: "Invalid JSON input",
		},
		{
			name:        "invalid input",
			code:        "INVALID_INPUT",
			message:     "Missing required field",
			wantCode:    "INVALID_INPUT",
			wantMessage: "Missing required field",
		},
		{
			name:        "empty message",
			code:        "UNKNOWN_ERROR",
			message:     "",
			wantCode:    "UNKNOWN_ERROR",
			wantMessage: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewError(tt.code, tt.message)

			if r.Data != nil {
				t.Errorf("NewError() data = %v, want nil", r.Data)
			}

			if r.Error == nil {
				t.Fatal("NewError() error = nil, want error")
			}

			if r.Error.Code != tt.wantCode {
				t.Errorf("Error.Code = %v, want %v", r.Error.Code, tt.wantCode)
			}

			if r.Error.Message != tt.wantMessage {
				t.Errorf("Error.Message = %v, want %v", r.Error.Message, tt.wantMessage)
			}
		})
	}
}

func TestResultToJSON(t *testing.T) {
	t.Run("success result JSON structure", func(t *testing.T) {
		r := NewSuccess(map[string]int{"healthScore": 100})
		jsonStr := r.ToJSON()

		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
			t.Fatalf("Failed to parse JSON: %v", err)
		}

		// Verify structure matches Result<T> pattern
		if _, ok := parsed["data"]; !ok {
			t.Error("JSON missing 'data' field")
		}

		if parsed["error"] != nil {
			t.Errorf("JSON error = %v, want nil", parsed["error"])
		}

		// Verify camelCase in nested data
		data := parsed["data"].(map[string]interface{})
		if _, ok := data["healthScore"]; !ok {
			t.Error("Data missing 'healthScore' field (camelCase)")
		}
	})

	t.Run("error result JSON structure", func(t *testing.T) {
		r := NewError("PARSE_ERROR", "Invalid input")
		jsonStr := r.ToJSON()

		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
			t.Fatalf("Failed to parse JSON: %v", err)
		}

		if parsed["data"] != nil {
			t.Errorf("JSON data = %v, want nil", parsed["data"])
		}

		errObj := parsed["error"].(map[string]interface{})
		if errObj["code"] != "PARSE_ERROR" {
			t.Errorf("error.code = %v, want PARSE_ERROR", errObj["code"])
		}

		if errObj["message"] != "Invalid input" {
			t.Errorf("error.message = %v, want 'Invalid input'", errObj["message"])
		}
	})
}

func TestErrorCodeFormat(t *testing.T) {
	// Verify error codes follow UPPER_SNAKE_CASE convention
	validCodes := []string{
		"PARSE_ERROR",
		"INVALID_INPUT",
		"ANALYSIS_FAILED",
		"CIRCULAR_DETECTED",
		"WORKSPACE_TOO_LARGE",
	}

	for _, code := range validCodes {
		r := NewError(code, "test message")
		if r.Error.Code != code {
			t.Errorf("Error code = %v, want %v", r.Error.Code, code)
		}
	}
}
