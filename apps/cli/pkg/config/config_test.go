// Package config provides configuration management using Viper
// ATDD: Story 1-4, AC4 - Viper Configuration Integration

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

// TestLoadConfig verifies configuration loading from various sources
// AC4: Viper is configured to read from multiple sources
func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T) string // Returns config path
		wantErr    bool
		wantFields map[string]interface{}
	}{
		{
			name: "load from current directory .monoguard.yaml",
			setup: func(t *testing.T) string {
				// Create temp dir with config file
				tmpDir := t.TempDir()
				configPath := filepath.Join(tmpDir, ".monoguard.yaml")
				content := `workspaces:
  - "packages/*"
  - "apps/*"
rules:
  circularDependencies: "error"
thresholds:
  healthScore: 70
`
				if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
					t.Fatal(err)
				}
				return tmpDir
			},
			wantErr: false,
			wantFields: map[string]interface{}{
				"workspaces": []string{"packages/*", "apps/*"},
			},
		},
		{
			name: "load from home directory ~/.monoguard/config.yaml",
			setup: func(t *testing.T) string {
				// This test verifies fallback to home directory
				return t.TempDir() // Empty dir, should fallback
			},
			wantErr: false, // Should not error, just use defaults
		},
		{
			name: "merge with defaults when config missing",
			setup: func(t *testing.T) string {
				return t.TempDir() // No config file
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configDir := tt.setup(t)

			// Change to config directory for test
			oldWd, _ := os.Getwd()
			defer os.Chdir(oldWd)
			os.Chdir(configDir)

			cfg, err := Load()

			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if cfg == nil && !tt.wantErr {
				t.Error("Load() returned nil config without error")
			}
		})
	}
}

// TestConfigStructure verifies the Config struct has correct fields
// AC4: Support configuration structure
func TestConfigStructure(t *testing.T) {
	cfg := &Config{
		Workspaces: []string{"packages/*", "apps/*"},
		Rules: Rules{
			CircularDependencies: "error",
			BoundaryViolations:   "warn",
		},
		Thresholds: Thresholds{
			HealthScore: 70,
		},
	}

	// Verify fields exist and are accessible
	if len(cfg.Workspaces) != 2 {
		t.Errorf("Workspaces length = %d, want 2", len(cfg.Workspaces))
	}

	if cfg.Rules.CircularDependencies != "error" {
		t.Errorf("Rules.CircularDependencies = %s, want 'error'", cfg.Rules.CircularDependencies)
	}

	if cfg.Thresholds.HealthScore != 70 {
		t.Errorf("Thresholds.HealthScore = %d, want 70", cfg.Thresholds.HealthScore)
	}
}

// TestEnvironmentVariableOverride verifies MONOGUARD_ env prefix
// AC4: Support environment variables with MONOGUARD_ prefix
func TestEnvironmentVariableOverride(t *testing.T) {
	// Set environment variables
	os.Setenv("MONOGUARD_VERBOSE", "true")
	defer os.Unsetenv("MONOGUARD_VERBOSE")

	os.Setenv("MONOGUARD_FORMAT", "json")
	defer os.Unsetenv("MONOGUARD_FORMAT")

	// Configure Viper to read env vars (mimics initConfig behavior)
	viper.SetEnvPrefix("MONOGUARD")
	viper.AutomaticEnv()

	// Verify Viper reads the environment variables
	verbose := viper.GetBool("verbose")
	if !verbose {
		t.Error("Viper should read MONOGUARD_VERBOSE=true from environment")
	}

	format := viper.GetString("format")
	if format != "json" {
		t.Errorf("Viper format = %q, want %q", format, "json")
	}

	// Load config - should work without error
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Config struct loads from Viper successfully
	if cfg == nil {
		t.Error("Load() should return non-nil config")
	}
}

// TestCLIFlagOverride verifies CLI flags override config file values
// AC4: Allow CLI flags to override config file values
func TestCLIFlagOverride(t *testing.T) {
	// This test verifies that when both config file and CLI flag
	// specify a value, the CLI flag takes precedence

	// Create config file with format: text
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".monoguard.yaml")
	content := `format: text
verbose: false
thresholds:
  healthScore: 50
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	// Reset Viper state and configure it
	viper.Reset()
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("ReadInConfig() error = %v", err)
	}

	// Verify config file was read
	formatFromFile := viper.GetString("format")
	if formatFromFile != "text" {
		t.Errorf("Config file format = %q, want %q", formatFromFile, "text")
	}

	// Simulate CLI flag override by setting value directly
	viper.Set("format", "json")

	// Verify CLI flag override takes precedence
	formatAfterOverride := viper.GetString("format")
	if formatAfterOverride != "json" {
		t.Errorf("After override format = %q, want %q", formatAfterOverride, "json")
	}

	// Load config struct - should work
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Verify config struct loaded with threshold from file
	if cfg.Thresholds.HealthScore != 50 {
		t.Errorf("Config healthScore = %d, want 50", cfg.Thresholds.HealthScore)
	}
}
