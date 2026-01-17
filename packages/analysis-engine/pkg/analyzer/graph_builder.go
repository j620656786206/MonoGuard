// Package analyzer provides dependency graph analysis for monorepo workspaces.
package analyzer

import (
	"errors"
	"sort"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// GraphBuilder constructs dependency graphs from workspace data.
// It separates internal workspace packages from external npm dependencies.
type GraphBuilder struct {
	workspacePackages map[string]bool    // Set of internal package names for O(1) lookup
	exclusionMatcher  *ExclusionMatcher  // Story 2.6: Pattern matcher for exclusions
}

// NewGraphBuilder creates a new graph builder instance.
func NewGraphBuilder() *GraphBuilder {
	return &GraphBuilder{
		workspacePackages: make(map[string]bool),
		exclusionMatcher:  nil,
	}
}

// NewGraphBuilderWithExclusions creates a graph builder with exclusion patterns.
// Returns an error if any regex pattern is invalid.
func NewGraphBuilderWithExclusions(excludePatterns []string) (*GraphBuilder, error) {
	matcher, err := NewExclusionMatcher(excludePatterns)
	if err != nil {
		return nil, err
	}
	return &GraphBuilder{
		workspacePackages: make(map[string]bool),
		exclusionMatcher:  matcher,
	}, nil
}

// Build constructs a DependencyGraph from WorkspaceData.
// Only internal workspace packages create edges; external deps are stored in node metadata.
// Returns an error if any package has an empty name.
func (gb *GraphBuilder) Build(workspace *types.WorkspaceData) (*types.DependencyGraph, error) {
	// Validate and initialize the set of workspace package names for O(1) lookup
	gb.workspacePackages = make(map[string]bool)
	for name := range workspace.Packages {
		if name == "" {
			return nil, errors.New("package name cannot be empty")
		}
		gb.workspacePackages[name] = true
	}

	// Build nodes
	nodes := gb.buildNodes(workspace)

	// Build edges (sorted for deterministic output)
	edges := gb.buildEdges(workspace, nodes)

	return &types.DependencyGraph{
		Nodes:         nodes,
		Edges:         edges,
		RootPath:      workspace.RootPath,
		WorkspaceType: workspace.WorkspaceType,
	}, nil
}

// buildNodes creates PackageNode entries for each package in the workspace.
// Internal dependencies are sorted alphabetically for deterministic output.
// Story 2.6: Excluded packages are marked with Excluded=true.
func (gb *GraphBuilder) buildNodes(workspace *types.WorkspaceData) map[string]*types.PackageNode {
	nodes := make(map[string]*types.PackageNode)

	for name, pkg := range workspace.Packages {
		node := types.NewPackageNode(pkg.Name, pkg.Version, pkg.Path)

		// Classify dependencies as internal or external, excluding self-references
		// Results are sorted alphabetically for deterministic output
		node.Dependencies, node.ExternalDeps = gb.classifyDependenciesExcludingSelf(pkg.Dependencies, name)
		node.DevDependencies, node.ExternalDevDeps = gb.classifyDependenciesExcludingSelf(pkg.DevDependencies, name)
		node.PeerDependencies, node.ExternalPeerDeps = gb.classifyDependenciesExcludingSelf(pkg.PeerDependencies, name)
		node.OptionalDependencies, node.ExternalOptionalDeps = gb.classifyDependenciesExcludingSelf(pkg.OptionalDependencies, name)

		// Story 2.6: Mark excluded packages
		node.Excluded = gb.isExcluded(name)

		nodes[name] = node
	}

	return nodes
}

// isExcluded checks if a package matches any exclusion pattern.
// Returns false if no exclusion matcher is configured.
func (gb *GraphBuilder) isExcluded(name string) bool {
	if gb.exclusionMatcher == nil {
		return false
	}
	return gb.exclusionMatcher.IsExcluded(name)
}

// buildEdges creates DependencyEdge entries for internal dependencies only.
// Edges are sorted by (From, To) for deterministic output.
func (gb *GraphBuilder) buildEdges(workspace *types.WorkspaceData, nodes map[string]*types.PackageNode) []*types.DependencyEdge {
	var edges []*types.DependencyEdge

	for name, pkg := range workspace.Packages {
		// Production dependencies
		edges = gb.addEdgesForDependencyType(edges, name, pkg.Dependencies, types.DependencyTypeProduction)

		// Development dependencies
		edges = gb.addEdgesForDependencyType(edges, name, pkg.DevDependencies, types.DependencyTypeDevelopment)

		// Peer dependencies
		edges = gb.addEdgesForDependencyType(edges, name, pkg.PeerDependencies, types.DependencyTypePeer)

		// Optional dependencies
		edges = gb.addEdgesForDependencyType(edges, name, pkg.OptionalDependencies, types.DependencyTypeOptional)
	}

	// Sort edges for deterministic output (by From, then To)
	sort.Slice(edges, func(i, j int) bool {
		if edges[i].From != edges[j].From {
			return edges[i].From < edges[j].From
		}
		return edges[i].To < edges[j].To
	})

	return edges
}

// addEdgesForDependencyType adds edges for a specific dependency type.
// Only creates edges for internal workspace packages, excluding self-references.
func (gb *GraphBuilder) addEdgesForDependencyType(
	edges []*types.DependencyEdge,
	fromPkg string,
	deps map[string]string,
	depType types.DependencyType,
) []*types.DependencyEdge {
	for depName, versionRange := range deps {
		if gb.isInternalPackage(depName) && depName != fromPkg {
			edges = append(edges, &types.DependencyEdge{
				From:         fromPkg,
				To:           depName,
				Type:         depType,
				VersionRange: versionRange,
			})
		}
	}
	return edges
}

// isInternalPackage checks if a dependency is a workspace package.
func (gb *GraphBuilder) isInternalPackage(name string) bool {
	return gb.workspacePackages[name]
}

// classifyDependenciesExcludingSelf separates internal and external dependencies,
// excluding self-references (when a package depends on itself).
// Returns (internal package names sorted alphabetically, external deps with versions).
func (gb *GraphBuilder) classifyDependenciesExcludingSelf(allDeps map[string]string, selfName string) (internal []string, external map[string]string) {
	internal = []string{}
	external = make(map[string]string)

	for name, version := range allDeps {
		// Skip self-references
		if name == selfName {
			continue
		}

		if gb.isInternalPackage(name) {
			internal = append(internal, name)
		} else {
			external[name] = version
		}
	}

	// Sort internal dependencies for deterministic output
	sort.Strings(internal)

	return internal, external
}
