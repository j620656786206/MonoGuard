package services

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// BasicPackageParser provides focused package.json parsing functionality
type BasicPackageParser struct {
	excludeDirs []string
}

// PackageData represents essential package.json information
type PackageData struct {
	Name             string            `json:"name"`
	Version          string            `json:"version"`
	Path             string            `json:"path"`
	Dependencies     map[string]string `json:"dependencies"`
	DevDependencies  map[string]string `json:"devDependencies"`
	PeerDependencies map[string]string `json:"peerDependencies"`
	Scripts          map[string]string `json:"scripts"`
	Workspaces       []string          `json:"workspaces"`
}

// ParseResult contains the results of parsing a repository
type ParseResult struct {
	Packages     []*PackageData `json:"packages"`
	RootPackage  *PackageData   `json:"rootPackage,omitempty"`
	PackageCount int            `json:"packageCount"`
	Errors       []ParseError   `json:"errors"`
}

// ParseError represents parsing errors
type ParseError struct {
	Path    string `json:"path"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

// NewBasicPackageParser creates a new basic package parser
func NewBasicPackageParser() *BasicPackageParser {
	return &BasicPackageParser{
		excludeDirs: []string{
			"node_modules",
			".git",
			"dist",
			"build",
			".next",
			"coverage",
			".nyc_output",
			"tmp",
			"temp",
			".cache",
		},
	}
}

// ParseRepository discovers and parses all package.json files in a repository
func (parser *BasicPackageParser) ParseRepository(rootPath string) (*ParseResult, error) {
	result := &ParseResult{
		Packages: make([]*PackageData, 0),
		Errors:   make([]ParseError, 0),
	}

	// Parse root package.json first
	rootPackagePath := filepath.Join(rootPath, "package.json")
	if rootPackage, err := parser.parsePackageFile(rootPackagePath); err == nil {
		result.RootPackage = rootPackage
		result.Packages = append(result.Packages, rootPackage)
	} else if !os.IsNotExist(err) {
		result.Errors = append(result.Errors, ParseError{
			Path:    rootPackagePath,
			Message: err.Error(),
			Type:    "root_parse_error",
		})
	}

	// Walk directory tree to find all package.json files
	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			result.Errors = append(result.Errors, ParseError{
				Path:    path,
				Message: err.Error(),
				Type:    "walk_error",
			})
			return nil // Continue walking
		}

		// Skip excluded directories
		if d.IsDir() && parser.shouldExcludeDir(d.Name()) {
			return filepath.SkipDir
		}

		// Parse package.json files (excluding root which we already processed)
		if d.Name() == "package.json" && path != rootPackagePath {
			if pkg, parseErr := parser.parsePackageFile(path); parseErr == nil {
				result.Packages = append(result.Packages, pkg)
			} else {
				result.Errors = append(result.Errors, ParseError{
					Path:    path,
					Message: parseErr.Error(),
					Type:    "package_parse_error",
				})
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	result.PackageCount = len(result.Packages)
	return result, nil
}

// parsePackageFile parses a single package.json file
func (parser *BasicPackageParser) parsePackageFile(filePath string) (*PackageData, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var rawPackage map[string]interface{}
	if err := json.Unmarshal(data, &rawPackage); err != nil {
		return nil, fmt.Errorf("invalid JSON in %s: %w", filePath, err)
	}

	pkg := &PackageData{
		Path:             filePath,
		Dependencies:     make(map[string]string),
		DevDependencies:  make(map[string]string),
		PeerDependencies: make(map[string]string),
		Scripts:          make(map[string]string),
		Workspaces:       make([]string, 0),
	}

	// Extract basic fields
	if name, ok := rawPackage["name"].(string); ok {
		pkg.Name = name
	}

	if version, ok := rawPackage["version"].(string); ok {
		pkg.Version = version
	}

	// Extract dependencies
	if deps, ok := rawPackage["dependencies"].(map[string]interface{}); ok {
		for name, version := range deps {
			if versionStr, ok := version.(string); ok {
				pkg.Dependencies[name] = versionStr
			}
		}
	}

	if devDeps, ok := rawPackage["devDependencies"].(map[string]interface{}); ok {
		for name, version := range devDeps {
			if versionStr, ok := version.(string); ok {
				pkg.DevDependencies[name] = versionStr
			}
		}
	}

	if peerDeps, ok := rawPackage["peerDependencies"].(map[string]interface{}); ok {
		for name, version := range peerDeps {
			if versionStr, ok := version.(string); ok {
				pkg.PeerDependencies[name] = versionStr
			}
		}
	}

	// Extract scripts
	if scripts, ok := rawPackage["scripts"].(map[string]interface{}); ok {
		for name, script := range scripts {
			if scriptStr, ok := script.(string); ok {
				pkg.Scripts[name] = scriptStr
			}
		}
	}

	// Extract workspaces
	if workspaces, ok := rawPackage["workspaces"]; ok {
		pkg.Workspaces = parser.extractWorkspaces(workspaces)
	}

	return pkg, nil
}

// extractWorkspaces handles different workspace configuration formats
func (parser *BasicPackageParser) extractWorkspaces(workspaces interface{}) []string {
	switch ws := workspaces.(type) {
	case []interface{}:
		// Array format: ["packages/*", "apps/*"]
		result := make([]string, 0, len(ws))
		for _, w := range ws {
			if str, ok := w.(string); ok {
				result = append(result, str)
			}
		}
		return result

	case map[string]interface{}:
		// Object format: {"packages": ["packages/*", "apps/*"]}
		if packages, ok := ws["packages"].([]interface{}); ok {
			result := make([]string, 0, len(packages))
			for _, p := range packages {
				if str, ok := p.(string); ok {
					result = append(result, str)
				}
			}
			return result
		}
	}

	return []string{}
}

// shouldExcludeDir checks if a directory should be excluded from parsing
func (parser *BasicPackageParser) shouldExcludeDir(dirName string) bool {
	for _, exclude := range parser.excludeDirs {
		if dirName == exclude {
			return true
		}
	}
	return false
}

// GetAllDependencies returns all dependencies from all packages
func (parser *BasicPackageParser) GetAllDependencies(result *ParseResult) map[string][]string {
	deps := make(map[string][]string)

	for _, pkg := range result.Packages {
		// Regular dependencies
		for name, version := range pkg.Dependencies {
			deps[name] = append(deps[name], version)
		}

		// Dev dependencies
		for name, version := range pkg.DevDependencies {
			deps[name] = append(deps[name], version)
		}

		// Peer dependencies
		for name, version := range pkg.PeerDependencies {
			deps[name] = append(deps[name], version)
		}
	}

	return deps
}

// GetPackagesByDependency returns packages that use a specific dependency
func (parser *BasicPackageParser) GetPackagesByDependency(result *ParseResult, dependencyName string) []*PackageData {
	var packages []*PackageData

	for _, pkg := range result.Packages {
		if parser.packageUsesDependency(pkg, dependencyName) {
			packages = append(packages, pkg)
		}
	}

	return packages
}

// packageUsesDependency checks if a package uses a specific dependency
func (parser *BasicPackageParser) packageUsesDependency(pkg *PackageData, dependencyName string) bool {
	if _, exists := pkg.Dependencies[dependencyName]; exists {
		return true
	}
	if _, exists := pkg.DevDependencies[dependencyName]; exists {
		return true
	}
	if _, exists := pkg.PeerDependencies[dependencyName]; exists {
		return true
	}
	return false
}

// FilterPackages filters packages based on a predicate function
func (parser *BasicPackageParser) FilterPackages(result *ParseResult, predicate func(*PackageData) bool) []*PackageData {
	var filtered []*PackageData
	for _, pkg := range result.Packages {
		if predicate(pkg) {
			filtered = append(filtered, pkg)
		}
	}
	return filtered
}

// GetWorkspacePackages returns only packages that are part of workspaces
func (parser *BasicPackageParser) GetWorkspacePackages(result *ParseResult) []*PackageData {
	if result.RootPackage == nil || len(result.RootPackage.Workspaces) == 0 {
		return []*PackageData{}
	}

	var workspacePackages []*PackageData
	rootDir := filepath.Dir(result.RootPackage.Path)

	for _, pkg := range result.Packages {
		if pkg == result.RootPackage {
			continue // Skip root package
		}

		pkgDir := filepath.Dir(pkg.Path)
		relPath, err := filepath.Rel(rootDir, pkgDir)
		if err != nil {
			continue
		}

		// Check if package path matches any workspace pattern
		for _, pattern := range result.RootPackage.Workspaces {
			if parser.matchesPattern(relPath, pattern) {
				workspacePackages = append(workspacePackages, pkg)
				break
			}
		}
	}

	return workspacePackages
}

// matchesPattern checks if a path matches a workspace pattern
func (parser *BasicPackageParser) matchesPattern(path, pattern string) bool {
	// Simple glob-like matching
	pattern = strings.ReplaceAll(pattern, "*", "")
	pattern = strings.TrimSuffix(pattern, "/")
	
	return strings.HasPrefix(path, pattern)
}

// ValidatePackageStructure performs basic validation on parsed packages
func (parser *BasicPackageParser) ValidatePackageStructure(result *ParseResult) []ValidationIssue {
	var issues []ValidationIssue

	for _, pkg := range result.Packages {
		// Check for missing name
		if pkg.Name == "" {
			issues = append(issues, ValidationIssue{
				Type:        "missing_name",
				Severity:    "error",
				Package:     pkg.Path,
				Message:     "Package missing name field",
				Suggestion:  "Add a name field to package.json",
			})
		}

		// Check for missing version
		if pkg.Version == "" {
			issues = append(issues, ValidationIssue{
				Type:        "missing_version",
				Severity:    "warning",
				Package:     pkg.Path,
				Message:     "Package missing version field",
				Suggestion:  "Add a version field to package.json",
			})
		}

		// Check for duplicate names
		for _, otherPkg := range result.Packages {
			if pkg != otherPkg && pkg.Name != "" && pkg.Name == otherPkg.Name {
				issues = append(issues, ValidationIssue{
					Type:        "duplicate_name",
					Severity:    "error",
					Package:     pkg.Path,
					Message:     fmt.Sprintf("Package name '%s' is duplicated", pkg.Name),
					Suggestion:  "Use unique package names across the monorepo",
					RelatedPath: otherPkg.Path,
				})
			}
		}
	}

	return issues
}

// ValidationIssue represents a package validation issue
type ValidationIssue struct {
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	Package     string `json:"package"`
	Message     string `json:"message"`
	Suggestion  string `json:"suggestion"`
	RelatedPath string `json:"relatedPath,omitempty"`
}