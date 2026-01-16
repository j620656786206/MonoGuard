// Package output provides formatted output utilities
// ATDD: Story 1-4, AC6 - All commands accept --format json|text flag

package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// TestFormatterText verifies text output formatting
// AC6: Commands output human-readable text by default
func TestFormatterText(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		wantText string
	}{
		{
			name: "simple string",
			data: "Hello, World!",
			wantText: "Hello, World!",
		},
		{
			name: "struct with fields",
			data: struct {
				Status  string `json:"status"`
				Message string `json:"message"`
			}{
				Status:  "placeholder",
				Message: "Analysis will be implemented in Epic 2",
			},
			wantText: "Status: placeholder",
		},
		{
			name: "analysis result placeholder",
			data: map[string]interface{}{
				"status":  "placeholder",
				"path":    ".",
				"message": "Will be implemented in Epic 2",
			},
			wantText: "placeholder",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFormatter("text")
			var buf bytes.Buffer

			err := f.PrintTo(&buf, tt.data)
			if err != nil {
				t.Fatalf("PrintTo() error = %v", err)
			}

			output := buf.String()
			if !strings.Contains(output, tt.wantText) {
				t.Errorf("PrintTo() output = %q, want to contain %q", output, tt.wantText)
			}
		})
	}
}

// TestFormatterJSON verifies JSON output formatting
// AC6: Commands accept --format json flag
func TestFormatterJSON(t *testing.T) {
	tests := []struct {
		name       string
		data       interface{}
		wantFields []string // Fields that must exist in JSON output
	}{
		{
			name: "analyze command JSON",
			data: map[string]interface{}{
				"status":  "placeholder",
				"path":    ".",
				"message": "Analysis will be implemented in Epic 2",
			},
			wantFields: []string{"status", "path", "message"},
		},
		{
			name: "check command JSON",
			data: map[string]interface{}{
				"status":  "placeholder",
				"path":    ".",
				"passed":  true,
				"message": "Check will be implemented in Epic 2",
			},
			wantFields: []string{"status", "passed", "message"},
		},
		{
			name: "fix command JSON",
			data: map[string]interface{}{
				"status":  "placeholder",
				"path":    ".",
				"dryRun":  true,
				"message": "Fix will be implemented in Epic 3",
			},
			wantFields: []string{"status", "dryRun", "message"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFormatter("json")
			var buf bytes.Buffer

			err := f.PrintTo(&buf, tt.data)
			if err != nil {
				t.Fatalf("PrintTo() error = %v", err)
			}

			// Verify valid JSON
			var parsed map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
				t.Fatalf("Output is not valid JSON: %v\nOutput: %s", err, buf.String())
			}

			// Verify required fields exist
			for _, field := range tt.wantFields {
				if _, ok := parsed[field]; !ok {
					t.Errorf("JSON output missing field %q", field)
				}
			}
		})
	}
}

// TestFormatterPrettyJSON verifies indented JSON output
// AC6: JSON output should be readable (pretty printed)
func TestFormatterPrettyJSON(t *testing.T) {
	f := NewFormatter("json")
	var buf bytes.Buffer

	data := map[string]interface{}{
		"status": "placeholder",
		"nested": map[string]string{
			"key": "value",
		},
	}

	err := f.PrintTo(&buf, data)
	if err != nil {
		t.Fatalf("PrintTo() error = %v", err)
	}

	output := buf.String()

	// Pretty JSON should contain newlines and indentation
	if !strings.Contains(output, "\n") {
		t.Error("JSON output should be pretty-printed with newlines")
	}
}

// TestNewFormatter verifies formatter creation
func TestNewFormatter(t *testing.T) {
	tests := []struct {
		name   string
		format string
		want   string
	}{
		{"text format", "text", "text"},
		{"json format", "json", "json"},
		{"default to text", "", "text"},
		{"uppercase JSON", "JSON", "json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFormatter(tt.format)
			if f == nil {
				t.Fatal("NewFormatter() returned nil")
			}
			if f.Format != tt.want {
				t.Errorf("Formatter.Format = %q, want %q", f.Format, tt.want)
			}
		})
	}
}

// TestFormatterPrint verifies Print to stdout
func TestFormatterPrint(t *testing.T) {
	f := NewFormatter("text")

	// Print should not error
	err := f.Print(map[string]string{"test": "data"})
	if err != nil {
		t.Errorf("Print() error = %v", err)
	}
}
