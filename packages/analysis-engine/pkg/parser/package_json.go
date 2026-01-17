// Package parser provides workspace configuration parsing for monorepos.
package parser

import (
	"encoding/json"
	"fmt"
)

// PackageJSON represents the structure of a package.json file.
// This type is used for parsing npm/yarn/pnpm package.json files.
type PackageJSON struct {
	Name             string            `json:"name"`
	Version          string            `json:"version"`
	Dependencies     map[string]string `json:"dependencies"`
	DevDependencies  map[string]string `json:"devDependencies"`
	PeerDependencies map[string]string `json:"peerDependencies"`
	// Workspaces can be either []string or WorkspacesConfig object
	// We use json.RawMessage to handle both formats
	Workspaces json.RawMessage `json:"workspaces"`
}

// WorkspacesConfig represents the extended workspaces format with packages and nohoist.
// Example: { "packages": ["packages/*"], "nohoist": ["**/react-native"] }
type WorkspacesConfig struct {
	Packages []string `json:"packages"`
	Nohoist  []string `json:"nohoist"`
}

// ParsePackageJSON parses a single package.json file from raw bytes.
// Returns an error if the JSON is invalid or cannot be parsed.
func ParsePackageJSON(data []byte) (*PackageJSON, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty input")
	}

	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, fmt.Errorf("failed to parse package.json: %w", err)
	}

	// Initialize nil maps to empty maps for consistency
	if pkg.Dependencies == nil {
		pkg.Dependencies = make(map[string]string)
	}
	if pkg.DevDependencies == nil {
		pkg.DevDependencies = make(map[string]string)
	}
	if pkg.PeerDependencies == nil {
		pkg.PeerDependencies = make(map[string]string)
	}

	return &pkg, nil
}

// ExtractWorkspacePatterns extracts workspace patterns from a parsed package.json.
// Handles both array format ["packages/*"] and object format {packages: [...], nohoist: [...]}.
// Returns nil if no workspaces field exists.
func ExtractWorkspacePatterns(pkg *PackageJSON) ([]string, error) {
	if pkg == nil {
		return nil, fmt.Errorf("nil package")
	}

	if len(pkg.Workspaces) == 0 {
		return nil, nil
	}

	// Try to parse as array format first (most common)
	var arrayFormat []string
	if err := json.Unmarshal(pkg.Workspaces, &arrayFormat); err == nil {
		return arrayFormat, nil
	}

	// Try to parse as object format
	var objectFormat WorkspacesConfig
	if err := json.Unmarshal(pkg.Workspaces, &objectFormat); err == nil {
		return objectFormat.Packages, nil
	}

	// If neither format works, return error
	return nil, fmt.Errorf("workspaces field has unsupported format")
}
