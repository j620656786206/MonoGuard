// Package tests contains E2E tests for the CLI binary
// ATDD: Story 1-4, AC3, AC5, AC6 - End-to-end CLI verification
//
// These tests run the actual compiled binary and verify output/exit codes.
// Prerequisites:
// - Run `make build` first to compile the binary
// - Binary should be at dist/monoguard

package tests

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// binaryPath returns the path to the compiled CLI binary
func binaryPath() string {
	// Get the directory of this test file
	_, err := os.Getwd()
	if err != nil {
		return "dist/monoguard"
	}
	return filepath.Join("..", "dist", "monoguard")
}

// runCLI executes the CLI binary with given arguments
func runCLI(t *testing.T, args ...string) (string, int, error) {
	t.Helper()

	binary := binaryPath()
	cmd := exec.Command(binary, args...)

	output, err := cmd.CombinedOutput()

	exitCode := 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	} else if err != nil {
		return "", -1, err
	}

	return string(output), exitCode, nil
}

// TestBinaryExists verifies the binary was built
// AC5: make build produces executable binary
func TestBinaryExists(t *testing.T) {
	binary := binaryPath()

	info, err := os.Stat(binary)
	if os.IsNotExist(err) {
		t.Fatalf("Binary not found at %s. Run 'make build' first.", binary)
	}
	if err != nil {
		t.Fatalf("Error checking binary: %v", err)
	}

	// Should be executable
	if info.Mode()&0111 == 0 {
		t.Error("Binary is not executable")
	}

	// AC5: Binary size should be reasonable (< 15MB)
	maxSize := int64(15 * 1024 * 1024) // 15MB
	if info.Size() > maxSize {
		t.Errorf("Binary size = %d bytes, want < %d bytes", info.Size(), maxSize)
	}
}

// TestHelpOutput verifies --help output
// AC3: ./monoguard --help shows expected structure
func TestHelpOutput(t *testing.T) {
	output, exitCode, err := runCLI(t, "--help")
	if err != nil {
		t.Fatalf("Failed to run CLI: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Exit code = %d, want 0", exitCode)
	}

	// AC3: Command name
	if !strings.Contains(output, "monoguard") {
		t.Error("Help should contain 'monoguard'")
	}

	// AC3: Available commands
	expectedCommands := []string{"analyze", "check", "fix", "init"}
	for _, cmd := range expectedCommands {
		if !strings.Contains(output, cmd) {
			t.Errorf("Help should list '%s' command", cmd)
		}
	}

	// AC3: Global flags
	expectedFlags := []string{"--config", "--verbose", "--format", "--version"}
	for _, flag := range expectedFlags {
		if !strings.Contains(output, flag) {
			t.Errorf("Help should list '%s' flag", flag)
		}
	}
}

// TestVersionOutput verifies --version output
// AC3: Version information via --version
func TestVersionOutput(t *testing.T) {
	output, exitCode, err := runCLI(t, "--version")
	if err != nil {
		t.Fatalf("Failed to run CLI: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Exit code = %d, want 0", exitCode)
	}

	// Should contain version number
	if !strings.Contains(output, "0.") && !strings.Contains(output, "monoguard") {
		t.Errorf("Version output = %q, should contain version number", output)
	}
}

// TestAnalyzeCommand verifies analyze command
// AC6: ./monoguard analyze outputs placeholder message
func TestAnalyzeCommand(t *testing.T) {
	t.Run("text output", func(t *testing.T) {
		output, exitCode, err := runCLI(t, "analyze", "--format", "text")
		if err != nil {
			t.Fatalf("Failed to run CLI: %v", err)
		}

		if exitCode != 0 {
			t.Errorf("Exit code = %d, want 0", exitCode)
		}

		if !strings.Contains(strings.ToLower(output), "placeholder") &&
			!strings.Contains(strings.ToLower(output), "analysis") {
			t.Errorf("Output should contain placeholder or analysis: %q", output)
		}
	})

	t.Run("json output", func(t *testing.T) {
		output, exitCode, err := runCLI(t, "analyze", "--format", "json")
		if err != nil {
			t.Fatalf("Failed to run CLI: %v", err)
		}

		if exitCode != 0 {
			t.Errorf("Exit code = %d, want 0", exitCode)
		}

		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(output), &parsed); err != nil {
			t.Fatalf("Output is not valid JSON: %v\nOutput: %s", err, output)
		}

		if _, ok := parsed["status"]; !ok {
			t.Error("JSON should contain 'status' field")
		}
	})
}

// TestCheckCommand verifies check command
// AC6: ./monoguard check exits with code 0
func TestCheckCommand(t *testing.T) {
	t.Run("exit code 0 on success", func(t *testing.T) {
		_, exitCode, err := runCLI(t, "check")
		if err != nil && exitCode != 0 && exitCode != 1 {
			t.Fatalf("Failed to run CLI: %v", err)
		}

		// AC6: Placeholder should return exit code 0
		if exitCode != 0 {
			t.Errorf("Exit code = %d, want 0 (placeholder success)", exitCode)
		}
	})

	t.Run("json output with passed field", func(t *testing.T) {
		output, _, err := runCLI(t, "check", "--format", "json")
		if err != nil {
			t.Logf("CLI error (may be expected): %v", err)
		}

		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(output), &parsed); err != nil {
			t.Fatalf("Output is not valid JSON: %v\nOutput: %s", err, output)
		}

		if passed, ok := parsed["passed"].(bool); ok {
			if !passed {
				t.Error("Placeholder should return passed=true")
			}
		} else {
			t.Error("JSON should contain boolean 'passed' field")
		}
	})
}

// TestFixCommand verifies fix command
// AC6: ./monoguard fix --dry-run outputs placeholder message
func TestFixCommand(t *testing.T) {
	t.Run("dry-run flag", func(t *testing.T) {
		output, exitCode, err := runCLI(t, "fix", "--dry-run", "--format", "json")
		if err != nil {
			t.Fatalf("Failed to run CLI: %v", err)
		}

		if exitCode != 0 {
			t.Errorf("Exit code = %d, want 0", exitCode)
		}

		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(output), &parsed); err != nil {
			t.Fatalf("Output is not valid JSON: %v\nOutput: %s", err, output)
		}

		if dryRun, ok := parsed["dryRun"].(bool); ok {
			if !dryRun {
				t.Error("dryRun should be true when --dry-run flag passed")
			}
		} else {
			t.Error("JSON should contain boolean 'dryRun' field")
		}
	})
}

// TestInitCommand verifies init command
// AC6: ./monoguard init outputs placeholder message
func TestInitCommand(t *testing.T) {
	output, exitCode, err := runCLI(t, "init", "--format", "json")
	if err != nil {
		t.Fatalf("Failed to run CLI: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Exit code = %d, want 0", exitCode)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("Output is not valid JSON: %v\nOutput: %s", err, output)
	}

	if _, ok := parsed["status"]; !ok {
		t.Error("JSON should contain 'status' field")
	}
}

// TestFormatFlag verifies --format flag works globally
// AC6: All commands accept --format json|text flag
func TestFormatFlag(t *testing.T) {
	commands := []string{"analyze", "check", "init"}

	for _, cmd := range commands {
		t.Run(cmd+" with json format", func(t *testing.T) {
			output, _, err := runCLI(t, cmd, "--format", "json")
			if err != nil {
				t.Logf("CLI error (may be expected): %v", err)
			}

			// Output should be valid JSON
			var parsed interface{}
			if err := json.Unmarshal([]byte(output), &parsed); err != nil {
				t.Errorf("%s --format json should output valid JSON: %v\nOutput: %s", cmd, err, output)
			}
		})

		t.Run(cmd+" with text format", func(t *testing.T) {
			output, _, err := runCLI(t, cmd, "--format", "text")
			if err != nil {
				t.Logf("CLI error (may be expected): %v", err)
			}

			// Output should NOT be JSON (text format)
			var parsed interface{}
			err = json.Unmarshal([]byte(output), &parsed)
			// Text format should fail JSON parsing or at least not be the default format
			if err == nil {
				// If it parses as JSON, make sure it's because it's a simple string
				if _, isString := parsed.(string); !isString {
					t.Logf("Note: %s --format text outputs JSON-like content", cmd)
				}
			}
		})
	}
}
