// Package types defines Go types that match TypeScript definitions in @monoguard/types.
// All JSON tags use camelCase for cross-language consistency.
package types

// AnalysisResult represents the complete analysis output.
// This matches @monoguard/types AnalysisResult.
type AnalysisResult struct {
	HealthScore int    `json:"healthScore"`
	Packages    int    `json:"packages"`
	CreatedAt   string `json:"createdAt"` // ISO 8601 format
}

// CheckResult represents validation-only output for CI/CD pipelines.
type CheckResult struct {
	Passed bool     `json:"passed"`
	Errors []string `json:"errors"`
}

// Package represents a single package in the workspace.
type Package struct {
	Name         string   `json:"name"`
	Path         string   `json:"path"`
	Dependencies []string `json:"dependencies"`
}

// CircularDependency represents a detected circular dependency chain.
type CircularDependency struct {
	Nodes []string `json:"nodes"` // Package names in the cycle
	Depth int      `json:"depth"` // Length of the cycle
}

// VersionConflict represents a dependency with conflicting versions.
type VersionConflict struct {
	PackageName string            `json:"packageName"`
	Versions    map[string]string `json:"versions"` // dependent -> version
}
