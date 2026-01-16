// ATDD: Story 1-4, AC6 - Placeholder Commands Work

package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// TestFixCommandRegistered verifies fix command is registered
// AC3: Available commands include "fix"
func TestFixCommandRegistered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "fix" {
			found = true
			break
		}
	}

	if !found {
		t.Error("fix command not registered on rootCmd")
	}
}

// TestFixCommandTextOutput verifies text output format
// AC6: ./monoguard fix --dry-run outputs placeholder message
func TestFixCommandTextOutput(t *testing.T) {
	ResetForTesting()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"fix", "--dry-run", "--format", "text"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()

	// Should contain fix or placeholder indicator
	if !strings.Contains(strings.ToLower(output), "fix") &&
		!strings.Contains(strings.ToLower(output), "placeholder") {
		t.Errorf("Text output should mention fix or placeholder: %q", output)
	}

	// Should indicate dry run mode
	if !strings.Contains(strings.ToLower(output), "dry") {
		t.Errorf("Text output should mention dry run: %q", output)
	}
}

// TestFixCommandJSONOutput verifies JSON output format
// AC6: ./monoguard fix --format json outputs valid JSON
func TestFixCommandJSONOutput(t *testing.T) {
	ResetForTesting()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"fix", "--dry-run", "--format", "json"})

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

	// Must contain dryRun field
	if _, ok := parsed["dryRun"]; !ok {
		t.Error("JSON output should contain 'dryRun' field")
	}

	// dryRun should be true when --dry-run flag is passed
	if dryRun, ok := parsed["dryRun"].(bool); ok {
		if !dryRun {
			t.Error("JSON output 'dryRun' should be true when --dry-run flag passed")
		}
	}
}

// TestFixCommandDryRunFlag verifies --dry-run flag
// AC6: Fix command supports --dry-run flag
func TestFixCommandDryRunFlag(t *testing.T) {
	cmd := findCommand("fix")
	if cmd == nil {
		t.Fatal("fix command not found")
	}

	// --dry-run flag
	dryRunFlag := cmd.Flags().Lookup("dry-run")
	if dryRunFlag == nil {
		t.Error("--dry-run flag not registered")
	}

	// Default should be false
	if dryRunFlag.DefValue != "false" {
		t.Errorf("--dry-run default = %q, want 'false'", dryRunFlag.DefValue)
	}
}

// TestFixCommandWithoutDryRun verifies fix without dry-run
func TestFixCommandWithoutDryRun(t *testing.T) {
	ResetForTesting()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"fix", "--format", "json"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	// dryRun should be false when not specified
	if dryRun, ok := parsed["dryRun"].(bool); ok {
		if dryRun {
			t.Error("JSON output 'dryRun' should be false when flag not specified")
		}
	}
}
