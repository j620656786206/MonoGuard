// Package parser provides pnpm-workspace.yaml parsing for monorepos.
package parser

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// PnpmWorkspace represents the structure of pnpm-workspace.yaml file.
// This type is used for parsing pnpm workspace configurations.
type PnpmWorkspace struct {
	// Packages is the list of glob patterns defining workspace packages.
	// Supports negation patterns (e.g., "!packages/deprecated-*").
	Packages []string `yaml:"packages"`
}

// ParsePnpmWorkspace parses a pnpm-workspace.yaml file from raw bytes.
// Returns an error if the YAML is invalid or cannot be parsed.
func ParsePnpmWorkspace(data []byte) (*PnpmWorkspace, error) {
	// Handle empty input - return empty struct, not error
	if len(data) == 0 {
		return &PnpmWorkspace{}, nil
	}

	var ws PnpmWorkspace
	if err := yaml.Unmarshal(data, &ws); err != nil {
		return nil, fmt.Errorf("failed to parse pnpm-workspace.yaml: %w", err)
	}

	return &ws, nil
}
