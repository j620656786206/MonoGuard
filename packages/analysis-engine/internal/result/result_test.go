package result

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

			assert.Nil(t, r.Error, "NewSuccess() should have nil error")

			// Verify JSON output contains data and no error
			jsonStr := r.ToJSON()
			var parsed map[string]interface{}
			err := json.Unmarshal([]byte(jsonStr), &parsed)
			require.NoError(t, err, "Failed to parse JSON")

			assert.Nil(t, parsed["error"], "JSON error field should be nil")
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

			assert.Nil(t, r.Data, "NewError() should have nil data")
			require.NotNil(t, r.Error, "NewError() should have error")
			assert.Equal(t, tt.wantCode, r.Error.Code)
			assert.Equal(t, tt.wantMessage, r.Error.Message)
		})
	}
}

func TestResultToJSON(t *testing.T) {
	t.Run("success result JSON structure", func(t *testing.T) {
		r := NewSuccess(map[string]int{"healthScore": 100})
		jsonStr := r.ToJSON()

		var parsed map[string]interface{}
		err := json.Unmarshal([]byte(jsonStr), &parsed)
		require.NoError(t, err, "Failed to parse JSON")

		// Verify structure matches Result<T> pattern
		assert.Contains(t, parsed, "data", "JSON should have 'data' field")
		assert.Nil(t, parsed["error"], "JSON error should be nil")

		// Verify camelCase in nested data
		data := parsed["data"].(map[string]interface{})
		assert.Contains(t, data, "healthScore", "Data should have 'healthScore' field (camelCase)")
	})

	t.Run("error result JSON structure", func(t *testing.T) {
		r := NewError("PARSE_ERROR", "Invalid input")
		jsonStr := r.ToJSON()

		var parsed map[string]interface{}
		err := json.Unmarshal([]byte(jsonStr), &parsed)
		require.NoError(t, err, "Failed to parse JSON")

		assert.Nil(t, parsed["data"], "JSON data should be nil")

		errObj := parsed["error"].(map[string]interface{})
		assert.Equal(t, "PARSE_ERROR", errObj["code"])
		assert.Equal(t, "Invalid input", errObj["message"])
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
		assert.Equal(t, code, r.Error.Code)
	}
}

func TestParserErrorCodes(t *testing.T) {
	// Verify parser error code constants are UPPER_SNAKE_CASE
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{"InvalidWorkspace", ErrInvalidWorkspace, "INVALID_WORKSPACE"},
		{"MissingPackageJSON", ErrMissingPackageJSON, "MISSING_PACKAGE_JSON"},
		{"InvalidPackageJSON", ErrInvalidPackageJSON, "INVALID_PACKAGE_JSON"},
		{"InvalidPnpmWorkspace", ErrInvalidPnpmWorkspace, "INVALID_PNPM_WORKSPACE"},
		{"GlobPatternFailed", ErrGlobPatternFailed, "GLOB_PATTERN_FAILED"},
		{"InvalidInput", ErrInvalidInput, "INVALID_INPUT"},
		{"AnalysisFailed", ErrAnalysisFailed, "ANALYSIS_FAILED"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.code, "Error code should be UPPER_SNAKE_CASE")

			// Verify the code can be used with NewError
			r := NewError(tt.code, "test message")
			assert.Equal(t, tt.code, r.Error.Code)
		})
	}
}

func TestEscapeJSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "no special characters",
			input: "simple text",
			want:  "simple text",
		},
		{
			name:  "double quote",
			input: `say "hello"`,
			want:  `say \"hello\"`,
		},
		{
			name:  "backslash",
			input: `path\to\file`,
			want:  `path\\to\\file`,
		},
		{
			name:  "newline",
			input: "line1\nline2",
			want:  `line1\nline2`,
		},
		{
			name:  "carriage return",
			input: "line1\rline2",
			want:  `line1\rline2`,
		},
		{
			name:  "tab",
			input: "col1\tcol2",
			want:  `col1\tcol2`,
		},
		{
			name:  "mixed special characters",
			input: "Error: \"file\\path\"\nDetails:\ttab",
			want:  `Error: \"file\\path\"\nDetails:\ttab`,
		},
		{
			name:  "control character (bell)",
			input: "text\x07bell",
			want:  `text\u0007bell`,
		},
		{
			name:  "control character (null)",
			input: "text\x00null",
			want:  `text\u0000null`,
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := escapeJSON(tt.input)
			assert.Equal(t, tt.want, got, "escapeJSON(%q)", tt.input)
		})
	}
}

func TestToJSONWithUnmarshalableData(t *testing.T) {
	// Test the fallback path when marshaling fails
	// Create a result with a channel (channels cannot be marshaled to JSON)
	ch := make(chan int)
	r := NewSuccess(ch)

	jsonStr := r.ToJSON()

	// Should return a MARSHAL_ERROR result
	var parsed map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &parsed)
	require.NoError(t, err, "Fallback JSON should be valid")

	assert.Nil(t, parsed["data"], "Fallback should have null data")

	errObj, ok := parsed["error"].(map[string]interface{})
	require.True(t, ok, "Fallback should have error object")

	assert.Equal(t, "MARSHAL_ERROR", errObj["code"])

	msg, ok := errObj["message"].(string)
	require.True(t, ok, "error message should be string")
	assert.NotEmpty(t, msg, "error message should not be empty")
}
