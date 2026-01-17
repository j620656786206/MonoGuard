// Package parser provides workspace configuration parsing for monorepos.
// Supports npm, yarn, and pnpm workspaces.
package parser

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// Parser handles workspace configuration parsing.
type Parser struct {
	rootPath string
}

// NewParser creates a new workspace parser.
func NewParser(rootPath string) *Parser {
	return &Parser{rootPath: rootPath}
}

// DetectWorkspaceType determines if workspace is npm/yarn/pnpm based on files present.
// Detection priority:
//  1. pnpm-workspace.yaml exists → pnpm
//  2. yarn.lock exists → yarn
//  3. package-lock.json exists → npm
//  4. Otherwise → unknown
func (p *Parser) DetectWorkspaceType(files map[string][]byte) types.WorkspaceType {
	// Priority 1: pnpm
	if _, ok := files["pnpm-workspace.yaml"]; ok {
		return types.WorkspaceTypePnpm
	}

	// Priority 2: yarn
	if _, ok := files["yarn.lock"]; ok {
		return types.WorkspaceTypeYarn
	}

	// Priority 3: npm
	if _, ok := files["package-lock.json"]; ok {
		return types.WorkspaceTypeNpm
	}

	return types.WorkspaceTypeUnknown
}

// Parse detects workspace type and parses all packages.
// Returns a structured WorkspaceData with all parsed packages.
func (p *Parser) Parse(files map[string][]byte) (*types.WorkspaceData, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no files provided")
	}

	// Check for root package.json
	rootPkgData, ok := files["package.json"]
	if !ok {
		return nil, fmt.Errorf("missing root package.json")
	}

	// Parse root package.json
	rootPkg, err := ParsePackageJSON(rootPkgData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse root package.json: %w", err)
	}

	// Detect workspace type
	wsType := p.DetectWorkspaceType(files)

	// Get workspace patterns
	patterns, err := p.getWorkspacePatterns(files, wsType, rootPkg)
	if err != nil {
		return nil, err
	}

	// Expand patterns to find package directories
	packageDirs := ExpandGlobPatternsFromFiles(files, patterns)

	// Parse each package
	packages := make(map[string]*types.PackageInfo)
	for _, dir := range packageDirs {
		pkgPath := filepath.ToSlash(filepath.Join(dir, "package.json"))
		pkgData, ok := files[pkgPath]
		if !ok {
			continue // Skip directories without package.json
		}

		pkg, err := ParsePackageJSON(pkgData)
		if err != nil {
			// Log warning but continue parsing other packages
			continue
		}

		if pkg.Name == "" {
			continue // Skip packages without name
		}

		// Initialize empty maps if nil
		deps := pkg.Dependencies
		if deps == nil {
			deps = make(map[string]string)
		}
		devDeps := pkg.DevDependencies
		if devDeps == nil {
			devDeps = make(map[string]string)
		}
		peerDeps := pkg.PeerDependencies
		if peerDeps == nil {
			peerDeps = make(map[string]string)
		}

		packages[pkg.Name] = &types.PackageInfo{
			Name:             pkg.Name,
			Version:          pkg.Version,
			Path:             dir,
			Dependencies:     deps,
			DevDependencies:  devDeps,
			PeerDependencies: peerDeps,
		}
	}

	return &types.WorkspaceData{
		RootPath:      p.rootPath,
		WorkspaceType: wsType,
		Packages:      packages,
	}, nil
}

// getWorkspacePatterns extracts workspace patterns based on workspace type.
func (p *Parser) getWorkspacePatterns(files map[string][]byte, wsType types.WorkspaceType, rootPkg *PackageJSON) ([]string, error) {
	switch wsType {
	case types.WorkspaceTypePnpm:
		// Parse pnpm-workspace.yaml
		pnpmData, ok := files["pnpm-workspace.yaml"]
		if !ok {
			// Fall back to package.json workspaces
			return ExtractWorkspacePatterns(rootPkg)
		}
		pnpmWs, err := ParsePnpmWorkspace(pnpmData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse pnpm-workspace.yaml: %w", err)
		}
		// Clean patterns (remove quotes that YAML preserves)
		patterns := make([]string, len(pnpmWs.Packages))
		for i, p := range pnpmWs.Packages {
			patterns[i] = strings.Trim(p, "'\"")
		}
		return patterns, nil

	case types.WorkspaceTypeNpm, types.WorkspaceTypeYarn:
		// Parse workspaces from package.json
		return ExtractWorkspacePatterns(rootPkg)

	default:
		// Unknown type - try package.json workspaces first
		patterns, err := ExtractWorkspacePatterns(rootPkg)
		if err != nil || len(patterns) == 0 {
			// No workspaces found, return empty
			return []string{}, nil
		}
		return patterns, nil
	}
}
