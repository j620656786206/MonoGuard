// Package types defines Go types that match TypeScript definitions in @monoguard/types.
// All JSON tags use camelCase for cross-language consistency.
package types

// ========================================
// Workspace Configuration Types (Story 2.1)
// ========================================

// WorkspaceType identifies the package manager workspace format.
// Matches @monoguard/types WorkspaceType.
type WorkspaceType string

const (
	WorkspaceTypeNpm     WorkspaceType = "npm"
	WorkspaceTypeYarn    WorkspaceType = "yarn"
	WorkspaceTypePnpm    WorkspaceType = "pnpm"
	WorkspaceTypeUnknown WorkspaceType = "unknown"
)

// WorkspaceData represents the complete parsed workspace configuration.
// Matches @monoguard/types WorkspaceData.
type WorkspaceData struct {
	RootPath      string                  `json:"rootPath"`
	WorkspaceType WorkspaceType           `json:"workspaceType"`
	Packages      map[string]*PackageInfo `json:"packages"`
}

// PackageInfo represents a single package in the workspace with full dependency information.
// This is the expanded version that includes version strings for all dependencies.
// Matches @monoguard/types Package interface.
type PackageInfo struct {
	Name                 string            `json:"name"`
	Version              string            `json:"version"`
	Path                 string            `json:"path"`
	Dependencies         map[string]string `json:"dependencies"`
	DevDependencies      map[string]string `json:"devDependencies"`
	PeerDependencies     map[string]string `json:"peerDependencies"`
	OptionalDependencies map[string]string `json:"optionalDependencies,omitempty"`
}

// ========================================
// Analysis Result Types
// ========================================

// AnalysisResult represents the complete analysis output.
// This matches @monoguard/types AnalysisResult.
type AnalysisResult struct {
	HealthScore int              `json:"healthScore"`
	Packages    int              `json:"packages"`
	Graph       *DependencyGraph `json:"graph,omitempty"`
	CreatedAt   string           `json:"createdAt,omitempty"` // ISO 8601 format
	Placeholder bool             `json:"placeholder,omitempty"` // True when returning placeholder data
}

// VersionInfo represents the version response.
type VersionInfo struct {
	Version string `json:"version"`
}

// CheckResult represents validation-only output for CI/CD pipelines.
type CheckResult struct {
	Passed      bool     `json:"passed"`
	Errors      []string `json:"errors"`
	Placeholder bool     `json:"placeholder,omitempty"` // True when returning placeholder data
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
