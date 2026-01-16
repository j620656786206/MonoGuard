// ATDD: Story 1-4, AC6 - Placeholder Commands Work

package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// TestInitCommandRegistered verifies init command is registered
// AC3: Available commands include "init"
func TestInitCommandRegistered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "init" {
			found = true
			break
		}
	}

	if !found {
		t.Error("init command not registered on rootCmd")
	}
}

// TestInitCommandTextOutput verifies text output format
// AC6: ./monoguard init outputs placeholder message
func TestInitCommandTextOutput(t *testing.T) {
	ResetForTesting()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"init", "--format", "text"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()

	// Should contain init or placeholder indicator
	if !strings.Contains(strings.ToLower(output), "init") &&
		!strings.Contains(strings.ToLower(output), "placeholder") {
		t.Errorf("Text output should mention init or placeholder: %q", output)
	}
}

// TestInitCommandJSONOutput verifies JSON output format
// AC6: ./monoguard init --format json outputs valid JSON
func TestInitCommandJSONOutput(t *testing.T) {
	ResetForTesting()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"init", "--format", "json"})

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

// TestInitCommandDescription verifies command has description
func TestInitCommandDescription(t *testing.T) {
	cmd := findCommand("init")
	if cmd == nil {
		t.Fatal("init command not found")
	}

	if cmd.Short == "" {
		t.Error("init command should have Short description")
	}

	if cmd.Long == "" {
		t.Error("init command should have Long description")
	}

	// Should mention configuration
	if !strings.Contains(strings.ToLower(cmd.Long), "config") {
		t.Error("init Long description should mention configuration")
	}
}
