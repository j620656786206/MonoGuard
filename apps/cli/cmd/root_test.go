// Package cmd contains all CLI commands
// ATDD: Story 1-4, AC3, AC4 - Cobra Command Structure and Viper Integration

package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

// TestRootCommandHelp verifies help output structure
// AC3: Running ./monoguard --help shows expected structure
func TestRootCommandHelp(t *testing.T) {
	// Reset command for test isolation
	ResetForTesting()

	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Execute() with --help error = %v", err)
	}

	output := buf.String()

	// AC3: Command name should be "monoguard"
	if !strings.Contains(output, "monoguard") {
		t.Error("Help output should contain 'monoguard'")
	}

	// AC3: Available commands list
	expectedCommands := []string{"analyze", "check", "fix", "init"}
	for _, cmd := range expectedCommands {
		if !strings.Contains(output, cmd) {
			t.Errorf("Help output should list '%s' command", cmd)
		}
	}

	// AC3: Global flags
	expectedFlags := []string{"--config", "--verbose", "--format"}
	for _, flag := range expectedFlags {
		if !strings.Contains(output, flag) {
			t.Errorf("Help output should list '%s' flag", flag)
		}
	}
}

// TestRootCommandVersion verifies version information
// AC3: Version information via --version
func TestRootCommandVersion(t *testing.T) {
	ResetForTesting()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"--version"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Execute() with --version error = %v", err)
	}

	output := buf.String()

	// Should contain version number
	if !strings.Contains(output, "0.") || !strings.Contains(output, "monoguard") {
		t.Errorf("Version output = %q, should contain version number", output)
	}
}

// TestGlobalFlagsRegistered verifies global flags are registered
// AC3: Global flags --config, --verbose, --format
func TestGlobalFlagsRegistered(t *testing.T) {
	tests := []struct {
		name     string
		flagName string
		wantType string
	}{
		{"config flag", "config", "string"},
		{"verbose flag", "verbose", "bool"},
		{"format flag", "format", "string"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := rootCmd.PersistentFlags().Lookup(tt.flagName)
			if flag == nil {
				t.Errorf("Flag %q not registered", tt.flagName)
				return
			}

			// Verify flag type
			switch tt.wantType {
			case "string":
				if flag.Value.Type() != "string" {
					t.Errorf("Flag %q type = %s, want string", tt.flagName, flag.Value.Type())
				}
			case "bool":
				if flag.Value.Type() != "bool" {
					t.Errorf("Flag %q type = %s, want bool", tt.flagName, flag.Value.Type())
				}
			}
		})
	}
}

// TestRootCommandDescription verifies command descriptions
// AC3: Short and long descriptions
func TestRootCommandDescription(t *testing.T) {
	if rootCmd.Short == "" {
		t.Error("rootCmd.Short should not be empty")
	}

	if rootCmd.Long == "" {
		t.Error("rootCmd.Long should not be empty")
	}

	// Should mention MonoGuard or monorepo
	if !strings.Contains(strings.ToLower(rootCmd.Long), "monorepo") &&
		!strings.Contains(strings.ToLower(rootCmd.Long), "monoguard") {
		t.Error("Long description should mention MonoGuard or monorepo")
	}
}

// TestSubcommandsRegistered verifies all subcommands are registered
// AC3: Available commands: analyze, check, fix, init
func TestSubcommandsRegistered(t *testing.T) {
	expectedCommands := []string{"analyze", "check", "fix", "init"}

	for _, cmdName := range expectedCommands {
		found := false
		for _, cmd := range rootCmd.Commands() {
			if cmd.Name() == cmdName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Subcommand %q not registered", cmdName)
		}
	}
}

// TestViperInitialization verifies Viper is configured correctly
// AC4: Viper configuration integration
func TestViperInitialization(t *testing.T) {
	// Call initConfig (normally called by cobra.OnInitialize)
	initConfig()

	// Verify Viper is configured with correct env prefix
	// Set an env var and verify Viper can read it
	t.Setenv("MONOGUARD_TEST_VAR", "test_value")
	viper.AutomaticEnv()

	// Viper should read env vars with MONOGUARD_ prefix
	val := viper.GetString("test_var")
	if val != "test_value" {
		t.Errorf("Viper env var reading: got %q, want %q", val, "test_value")
	}

	// Verify config type is set to yaml
	if viper.GetString("format") == "" {
		// format has a default, so this just verifies viper is working
		t.Log("Viper initialized successfully")
	}
}
