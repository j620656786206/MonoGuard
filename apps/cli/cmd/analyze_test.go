// ATDD: Story 1-4, AC6 - Placeholder Commands Work

package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// TestAnalyzeCommandRegistered verifies analyze command is registered
// AC3: Available commands include "analyze"
func TestAnalyzeCommandRegistered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "analyze" {
			found = true
			break
		}
	}

	if !found {
		t.Error("analyze command not registered on rootCmd")
	}
}

// TestAnalyzeCommandUsage verifies command usage string
// AC6: monoguard analyze [path]
func TestAnalyzeCommandUsage(t *testing.T) {
	cmd := findCommand("analyze")
	if cmd == nil {
		t.Fatal("analyze command not found")
	}

	if !strings.Contains(cmd.Use, "analyze") {
		t.Errorf("Use = %q, should contain 'analyze'", cmd.Use)
	}
}

// TestAnalyzeCommandTextOutput verifies text output format
// AC6: ./monoguard analyze outputs placeholder message
func TestAnalyzeCommandTextOutput(t *testing.T) {
	ResetForTesting()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"analyze", "--format", "text"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()

	// Should contain placeholder indicator
	if !strings.Contains(strings.ToLower(output), "placeholder") &&
		!strings.Contains(strings.ToLower(output), "epic 2") {
		t.Errorf("Text output should mention placeholder or Epic 2: %q", output)
	}

	// Should mention analysis
	if !strings.Contains(strings.ToLower(output), "analysis") &&
		!strings.Contains(strings.ToLower(output), "analyze") {
		t.Errorf("Text output should mention analysis: %q", output)
	}
}

// TestAnalyzeCommandJSONOutput verifies JSON output format
// AC6: ./monoguard analyze --format json outputs valid JSON
func TestAnalyzeCommandJSONOutput(t *testing.T) {
	ResetForTesting()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"analyze", "--format", "json"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()

	// Must be valid JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("Output is not valid JSON: %v\nOutput: %s", err, output)
	}

	// Must contain status field
	if _, ok := parsed["status"]; !ok {
		t.Error("JSON output should contain 'status' field")
	}

	// Must contain message field
	if _, ok := parsed["message"]; !ok {
		t.Error("JSON output should contain 'message' field")
	}
}

// TestAnalyzeCommandWithPath verifies path argument handling
// AC6: analyze accepts optional path argument
func TestAnalyzeCommandWithPath(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantPath string
	}{
		{
			name:     "no path defaults to current dir",
			args:     []string{"analyze"},
			wantPath: ".",
		},
		{
			name:     "explicit path",
			args:     []string{"analyze", "/some/path"},
			wantPath: "/some/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ResetForTesting()

			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)
			rootCmd.SetArgs(append(tt.args, "--format", "json"))

			err := rootCmd.Execute()
			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}

			output := buf.String()

			var parsed map[string]interface{}
			if err := json.Unmarshal([]byte(output), &parsed); err != nil {
				t.Fatalf("Output is not valid JSON: %v", err)
			}

			if path, ok := parsed["path"].(string); ok {
				if path != tt.wantPath {
					t.Errorf("path = %q, want %q", path, tt.wantPath)
				}
			}
		})
	}
}

// Helper function to find a command by name
func findCommand(name string) *cobra.Command {
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == name {
			return cmd
		}
	}
	return nil
}
