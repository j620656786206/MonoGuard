// Package types defines Go types that match TypeScript definitions in @monoguard/types.
// This file contains dependency graph types for Story 2.2.
package types

// ========================================
// Dependency Graph Types (Story 2.2)
// ========================================

// DependencyGraph represents the complete dependency structure of a workspace.
// Matches @monoguard/types DependencyGraph interface.
type DependencyGraph struct {
	Nodes         map[string]*PackageNode `json:"nodes"`
	Edges         []*DependencyEdge       `json:"edges"`
	RootPath      string                  `json:"rootPath"`
	WorkspaceType WorkspaceType           `json:"workspaceType"`
}

// PackageNode represents a single package in the dependency graph.
// This differs from PackageInfo:
//   - Dependencies/DevDependencies/PeerDependencies/OptionalDependencies contain ONLY internal workspace package names ([]string)
//   - ExternalDeps/ExternalDevDeps/ExternalPeerDeps/ExternalOptionalDeps contain external npm packages with versions (map[string]string)
//
// Matches @monoguard/types PackageNode interface.
type PackageNode struct {
	Name                 string            `json:"name"`
	Version              string            `json:"version"`
	Path                 string            `json:"path"`
	Dependencies         []string          `json:"dependencies"`         // Internal workspace package names only
	DevDependencies      []string          `json:"devDependencies"`      // Internal workspace package names only
	PeerDependencies     []string          `json:"peerDependencies"`     // Internal workspace package names only
	OptionalDependencies []string          `json:"optionalDependencies"` // Internal workspace package names only
	ExternalDeps         map[string]string `json:"externalDeps,omitempty"`
	ExternalDevDeps      map[string]string `json:"externalDevDeps,omitempty"`
	ExternalPeerDeps     map[string]string `json:"externalPeerDeps,omitempty"`
	ExternalOptionalDeps map[string]string `json:"externalOptionalDeps,omitempty"`
	Excluded             bool              `json:"excluded,omitempty"` // Story 2.6: True if excluded from analysis
}

// DependencyEdge represents a directed edge between packages in the dependency graph.
// Edges only exist between internal workspace packages (not external npm packages).
// Matches @monoguard/types DependencyEdge interface.
type DependencyEdge struct {
	From         string         `json:"from"`
	To           string         `json:"to"`
	Type         DependencyType `json:"type"`
	VersionRange string         `json:"versionRange"`
}

// DependencyType classifies the type of dependency relationship.
// Matches @monoguard/types DependencyType union type.
type DependencyType string

const (
	DependencyTypeProduction  DependencyType = "production"
	DependencyTypeDevelopment DependencyType = "development"
	DependencyTypePeer        DependencyType = "peer"
	DependencyTypeOptional    DependencyType = "optional"
)

// NewDependencyGraph creates a new empty dependency graph.
func NewDependencyGraph(rootPath string, workspaceType WorkspaceType) *DependencyGraph {
	return &DependencyGraph{
		Nodes:         make(map[string]*PackageNode),
		Edges:         []*DependencyEdge{},
		RootPath:      rootPath,
		WorkspaceType: workspaceType,
	}
}

// NewPackageNode creates a new package node with initialized slices.
func NewPackageNode(name, version, path string) *PackageNode {
	return &PackageNode{
		Name:                 name,
		Version:              version,
		Path:                 path,
		Dependencies:         []string{},
		DevDependencies:      []string{},
		PeerDependencies:     []string{},
		OptionalDependencies: []string{},
		ExternalDeps:         make(map[string]string),
		ExternalDevDeps:      make(map[string]string),
		ExternalPeerDeps:     make(map[string]string),
		ExternalOptionalDeps: make(map[string]string),
	}
}
