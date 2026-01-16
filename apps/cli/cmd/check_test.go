// ATDD: Story 1-4, AC6 - Placeholder Commands Work

package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// TestCheckCommandRegistered verifies check command is registered
// AC3: Available commands include "check"
func TestCheckCommandRegistered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "check" {
			found = true
			break
		}
	}

	if !found {
		t.Error("check command not registered on rootCmd")
	}
}

// TestCheckCommandExitCode verifies exit code 0 for placeholder success
// AC6: ./monoguard check exits with code 0 (placeholder success)
func TestCheckCommandExitCode(t *testing.T) {
	ResetForTesting()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"check"})

	err := rootCmd.Execute()

	// Should not return error (exit code 0)
	if err != nil {
		t.Errorf("Execute() error = %v, want nil (exit code 0)", err)
	}
}

// TestCheckCommandTextOutput verifies text output format
// AC6: ./monoguard check outputs placeholder message
func TestCheckCommandTextOutput(t *testing.T) {
	ResetForTesting()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"check", "--format", "text"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()

	// Should contain check or placeholder indicator
	if !strings.Contains(strings.ToLower(output), "check") &&
		!strings.Contains(strings.ToLower(output), "placeholder") {
		t.Errorf("Text output should mention check or placeholder: %q", output)
	}

	// Should indicate success (passed)
	if !strings.Contains(strings.ToLower(output), "pass") &&
		!strings.Contains(output, "✅") {
		t.Errorf("Text output should indicate success: %q", output)
	}
}

// TestCheckCommandJSONOutput verifies JSON output format
// AC6: ./monoguard check --format json outputs valid JSON
func TestCheckCommandJSONOutput(t *testing.T) {
	ResetForTesting()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"check", "--format", "json"})

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

	// Must contain passed field (boolean)
	if _, ok := parsed["passed"]; !ok {
		t.Error("JSON output should contain 'passed' field")
	}

	// passed should be true for placeholder
	if passed, ok := parsed["passed"].(bool); ok {
		if !passed {
			t.Error("JSON output 'passed' should be true for placeholder")
		}
	}
}

// TestCheckCommandFlags verifies check-specific flags
// AC6: Check command supports --fail-on and --threshold flags
func TestCheckCommandFlags(t *testing.T) {
	cmd := findCommand("check")
	if cmd == nil {
		t.Fatal("check command not found")
	}

	// --fail-on flag
	failOnFlag := cmd.Flags().Lookup("fail-on")
	if failOnFlag == nil {
		t.Error("--fail-on flag not registered")
	}

	// --threshold flag
	thresholdFlag := cmd.Flags().Lookup("threshold")
	if thresholdFlag == nil {
		t.Error("--threshold flag not registered")
	}
}

// TestCheckCommandCIMode verifies behavior suitable for CI/CD
// AC6: Check command designed for CI/CD integration
func TestCheckCommandCIMode(t *testing.T) {
	// In CI mode, check should:
	// 1. Exit 0 if all checks pass
	// 2. Exit 1 if any checks fail
	// 3. Output machine-readable format

	ResetForTesting()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"check", "--format", "json"})

	err := rootCmd.Execute()

	// Placeholder always passes
	if err != nil {
		t.Errorf("CI mode should return exit code 0 for placeholder: %v", err)
	}

	output := buf.String()

	// Output should be parseable by CI tools
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("CI output must be valid JSON: %v", err)
	}
}
