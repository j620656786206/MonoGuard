package services

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

// WorkspaceType represents different types of workspace configurations
type WorkspaceType string

const (
	WorkspaceTypeNpm   WorkspaceType = "npm"
	WorkspaceTypePnpm  WorkspaceType = "pnpm"
	WorkspaceTypeLerna WorkspaceType = "lerna"
	WorkspaceTypeNx    WorkspaceType = "nx"
	WorkspaceTypeYarn  WorkspaceType = "yarn"
)

// WorkspaceConfiguration represents a workspace configuration
type WorkspaceConfiguration struct {
	Type           WorkspaceType `json:"type"`
	RootPath       string        `json:"rootPath"`
	ConfigPath     string        `json:"configPath"`
	Packages       []string      `json:"packages"`
	PackageManager string        `json:"packageManager"`
	Priority       int           `json:"priority"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// WorkspacePackage represents a package within a workspace
type WorkspacePackage struct {
	Name         string            `json:"name"`
	Path         string            `json:"path"`
	Version      string            `json:"version"`
	Dependencies map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	PeerDependencies map[string]string `json:"peerDependencies"`
	Private      bool              `json:"private"`
	Workspace    *WorkspaceConfiguration `json:"workspace,omitempty"`
}

// PnpmWorkspace represents pnpm-workspace.yaml structure
type PnpmWorkspace struct {
	Packages []string `yaml:"packages"`
}

// LernaConfig represents lerna.json structure
type LernaConfig struct {
	Version      string   `json:"version"`
	Packages     []string `json:"packages"`
	Command      map[string]interface{} `json:"command"`
	NpmClient    string   `json:"npmClient"`
	UseWorkspaces bool    `json:"useWorkspaces"`
}

// NxWorkspace represents nx.json structure
type NxWorkspace struct {
	Version          int                    `json:"version"`
	Projects         map[string]interface{} `json:"projects"`
	WorkspaceLayout  map[string]string     `json:"workspaceLayout"`
	TargetDefaults   map[string]interface{} `json:"targetDefaults"`
	TasksRunnerOptions map[string]interface{} `json:"tasksRunnerOptions"`
}

// WorkspaceParser handles parsing of different workspace configurations
type WorkspaceParser struct {
	cache       map[string]*WorkspaceConfiguration
	packageCache map[string]*WorkspacePackage
	cacheMutex  sync.RWMutex
	stringInternMap map[string]string
	stringInternMutex sync.RWMutex
}

// NewWorkspaceParser creates a new workspace parser with optimizations
func NewWorkspaceParser() *WorkspaceParser {
	return &WorkspaceParser{
		cache:           make(map[string]*WorkspaceConfiguration),
		packageCache:    make(map[string]*WorkspacePackage),
		stringInternMap: make(map[string]string),
	}
}

// intern implements string interning for memory optimization
func (wp *WorkspaceParser) intern(s string) string {
	wp.stringInternMutex.RLock()
	if interned, exists := wp.stringInternMap[s]; exists {
		wp.stringInternMutex.RUnlock()
		return interned
	}
	wp.stringInternMutex.RUnlock()

	wp.stringInternMutex.Lock()
	defer wp.stringInternMutex.Unlock()
	
	// Double-check after acquiring write lock
	if interned, exists := wp.stringInternMap[s]; exists {
		return interned
	}
	
	wp.stringInternMap[s] = s
	return s
}

// DiscoverWorkspaces discovers all workspace configurations in a repository
func (wp *WorkspaceParser) DiscoverWorkspaces(rootPath string) ([]*WorkspaceConfiguration, error) {
	wp.cacheMutex.RLock()
	if cached, exists := wp.cache[rootPath]; exists {
		wp.cacheMutex.RUnlock()
		return []*WorkspaceConfiguration{cached}, nil
	}
	wp.cacheMutex.RUnlock()

	var workspaces []*WorkspaceConfiguration
	var discoveredPaths []string

	// Walk the directory tree to find workspace configuration files
	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip node_modules and common build directories
		if d.IsDir() {
			name := d.Name()
			if name == "node_modules" || name == ".git" || name == "dist" || 
			   name == "build" || name == ".next" || name == "coverage" {
				return filepath.SkipDir
			}
		}

		if d.IsDir() {
			return nil
		}

		// Check for workspace configuration files
		filename := d.Name()
		switch filename {
		case "pnpm-workspace.yaml", "pnpm-workspace.yml":
			if workspace, err := wp.parsePnpmWorkspace(path, rootPath); err == nil {
				workspaces = append(workspaces, workspace)
				discoveredPaths = append(discoveredPaths, path)
			}
		case "lerna.json":
			if workspace, err := wp.parseLernaWorkspace(path, rootPath); err == nil {
				workspaces = append(workspaces, workspace)
				discoveredPaths = append(discoveredPaths, path)
			}
		case "nx.json":
			if workspace, err := wp.parseNxWorkspace(path, rootPath); err == nil {
				workspaces = append(workspaces, workspace)
				discoveredPaths = append(discoveredPaths, path)
			}
		case "package.json":
			if workspace, err := wp.parseNpmWorkspace(path, rootPath); err == nil && workspace != nil {
				workspaces = append(workspaces, workspace)
				discoveredPaths = append(discoveredPaths, path)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to discover workspaces: %w", err)
	}

	// Apply priority-based resolution for conflicting workspace definitions
	resolvedWorkspaces := wp.resolveWorkspaceConflicts(workspaces)

	// Cache the results
	wp.cacheMutex.Lock()
	for _, workspace := range resolvedWorkspaces {
		wp.cache[workspace.RootPath] = workspace
	}
	wp.cacheMutex.Unlock()

	return resolvedWorkspaces, nil
}

// parsePnpmWorkspace parses pnpm-workspace.yaml files
func (wp *WorkspaceParser) parsePnpmWorkspace(configPath, rootPath string) (*WorkspaceConfiguration, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read pnpm workspace file: %w", err)
	}

	var pnpmConfig PnpmWorkspace
	if err := yaml.Unmarshal(data, &pnpmConfig); err != nil {
		return nil, fmt.Errorf("failed to parse pnpm workspace YAML: %w", err)
	}

	return &WorkspaceConfiguration{
		Type:           WorkspaceTypePnpm,
		RootPath:       wp.intern(rootPath),
		ConfigPath:     wp.intern(configPath),
		Packages:       wp.internSlice(pnpmConfig.Packages),
		PackageManager: wp.intern("pnpm"),
		Priority:       1, // High priority for pnpm
		Metadata: map[string]interface{}{
			"configType": "pnpm-workspace.yaml",
		},
	}, nil
}

// parseLernaWorkspace parses lerna.json files
func (wp *WorkspaceParser) parseLernaWorkspace(configPath, rootPath string) (*WorkspaceConfiguration, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read lerna config file: %w", err)
	}

	var lernaConfig LernaConfig
	if err := json.Unmarshal(data, &lernaConfig); err != nil {
		return nil, fmt.Errorf("failed to parse lerna JSON: %w", err)
	}

	packages := lernaConfig.Packages
	if len(packages) == 0 {
		packages = []string{"packages/*"} // Default Lerna pattern
	}

	return &WorkspaceConfiguration{
		Type:           WorkspaceTypeLerna,
		RootPath:       wp.intern(rootPath),
		ConfigPath:     wp.intern(configPath),
		Packages:       wp.internSlice(packages),
		PackageManager: wp.intern(lernaConfig.NpmClient),
		Priority:       2, // Medium-high priority for Lerna
		Metadata: map[string]interface{}{
			"version":       lernaConfig.Version,
			"useWorkspaces": lernaConfig.UseWorkspaces,
			"command":       lernaConfig.Command,
		},
	}, nil
}

// parseNxWorkspace parses nx.json files
func (wp *WorkspaceParser) parseNxWorkspace(configPath, rootPath string) (*WorkspaceConfiguration, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read nx config file: %w", err)
	}

	var nxConfig NxWorkspace
	if err := json.Unmarshal(data, &nxConfig); err != nil {
		return nil, fmt.Errorf("failed to parse nx JSON: %w", err)
	}

	// Extract project patterns from Nx configuration
	var packages []string
	for projectName := range nxConfig.Projects {
		packages = append(packages, projectName)
	}

	// Check for workspace layout
	if libsDir, exists := nxConfig.WorkspaceLayout["libsDir"]; exists {
		packages = append(packages, filepath.Join(libsDir, "*"))
	}
	if appsDir, exists := nxConfig.WorkspaceLayout["appsDir"]; exists {
		packages = append(packages, filepath.Join(appsDir, "*"))
	}

	return &WorkspaceConfiguration{
		Type:           WorkspaceTypeNx,
		RootPath:       wp.intern(rootPath),
		ConfigPath:     wp.intern(configPath),
		Packages:       wp.internSlice(packages),
		PackageManager: wp.intern("npm"), // Default for Nx
		Priority:       3, // Medium priority for Nx
		Metadata: map[string]interface{}{
			"version":        nxConfig.Version,
			"projects":       nxConfig.Projects,
			"workspaceLayout": nxConfig.WorkspaceLayout,
		},
	}, nil
}

// parseNpmWorkspace parses package.json files for npm workspaces
func (wp *WorkspaceParser) parseNpmWorkspace(configPath, rootPath string) (*WorkspaceConfiguration, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read package.json file: %w", err)
	}

	var packageJSON struct {
		Name       string   `json:"name"`
		Workspaces interface{} `json:"workspaces"`
	}
	if err := json.Unmarshal(data, &packageJSON); err != nil {
		return nil, fmt.Errorf("failed to parse package.json: %w", err)
	}

	// Extract workspaces configuration
	var packages []string
	switch workspaces := packageJSON.Workspaces.(type) {
	case []interface{}:
		for _, ws := range workspaces {
			if wsStr, ok := ws.(string); ok {
				packages = append(packages, wsStr)
			}
		}
	case map[string]interface{}:
		if packagesArray, exists := workspaces["packages"]; exists {
			if packagesSlice, ok := packagesArray.([]interface{}); ok {
				for _, pkg := range packagesSlice {
					if pkgStr, ok := pkg.(string); ok {
						packages = append(packages, pkgStr)
					}
				}
			}
		}
	case nil:
		return nil, nil // No workspaces configuration
	}

	if len(packages) == 0 {
		return nil, nil // No workspace packages defined
	}

	workspaceType := WorkspaceTypeNpm
	if wp.hasYarnLock(filepath.Dir(configPath)) {
		workspaceType = WorkspaceTypeYarn
	}

	return &WorkspaceConfiguration{
		Type:           workspaceType,
		RootPath:       wp.intern(rootPath),
		ConfigPath:     wp.intern(configPath),
		Packages:       wp.internSlice(packages),
		PackageManager: wp.intern(string(workspaceType)),
		Priority:       4, // Lower priority for npm/yarn workspaces
		Metadata: map[string]interface{}{
			"name": packageJSON.Name,
		},
	}, nil
}

// hasYarnLock checks if yarn.lock exists in the directory
func (wp *WorkspaceParser) hasYarnLock(dir string) bool {
	yarnLockPath := filepath.Join(dir, "yarn.lock")
	_, err := os.Stat(yarnLockPath)
	return err == nil
}

// resolveWorkspaceConflicts applies priority-based resolution for conflicting workspace definitions
func (wp *WorkspaceParser) resolveWorkspaceConflicts(workspaces []*WorkspaceConfiguration) []*WorkspaceConfiguration {
	if len(workspaces) <= 1 {
		return workspaces
	}

	// Group by root path
	rootGroups := make(map[string][]*WorkspaceConfiguration)
	for _, workspace := range workspaces {
		rootGroups[workspace.RootPath] = append(rootGroups[workspace.RootPath], workspace)
	}

	var resolved []*WorkspaceConfiguration
	for _, group := range rootGroups {
		if len(group) == 1 {
			resolved = append(resolved, group[0])
			continue
		}

		// Sort by priority (lower number = higher priority)
		sort.Slice(group, func(i, j int) bool {
			return group[i].Priority < group[j].Priority
		})

		// Take the highest priority workspace
		resolved = append(resolved, group[0])
	}

	return resolved
}

// DiscoverPackages discovers all packages within workspace configurations
func (wp *WorkspaceParser) DiscoverPackages(workspaces []*WorkspaceConfiguration) ([]*WorkspacePackage, error) {
	var allPackages []*WorkspacePackage
	var mu sync.Mutex

	// Process workspaces concurrently
	var wg sync.WaitGroup
	packageChan := make(chan *WorkspacePackage, 100)
	errorChan := make(chan error, len(workspaces))

	for _, workspace := range workspaces {
		wg.Add(1)
		go func(ws *WorkspaceConfiguration) {
			defer wg.Done()
			packages, err := wp.discoverPackagesInWorkspace(ws)
			if err != nil {
				errorChan <- err
				return
			}
			for _, pkg := range packages {
				packageChan <- pkg
			}
		}(workspace)
	}

	// Collect results
	go func() {
		wg.Wait()
		close(packageChan)
		close(errorChan)
	}()

	// Collect packages
	for pkg := range packageChan {
		mu.Lock()
		allPackages = append(allPackages, pkg)
		mu.Unlock()
	}

	// Check for errors
	if len(errorChan) > 0 {
		return allPackages, <-errorChan
	}

	// Remove duplicates and cache results
	uniquePackages := wp.deduplicatePackages(allPackages)
	
	wp.cacheMutex.Lock()
	for _, pkg := range uniquePackages {
		wp.packageCache[pkg.Path] = pkg
	}
	wp.cacheMutex.Unlock()

	return uniquePackages, nil
}

// discoverPackagesInWorkspace discovers packages within a single workspace
func (wp *WorkspaceParser) discoverPackagesInWorkspace(workspace *WorkspaceConfiguration) ([]*WorkspacePackage, error) {
	var packages []*WorkspacePackage
	rootDir := filepath.Dir(workspace.ConfigPath)

	for _, pattern := range workspace.Packages {
		matches, err := wp.expandPackagePattern(rootDir, pattern)
		if err != nil {
			continue // Skip invalid patterns
		}

		for _, match := range matches {
			packageJSONPath := filepath.Join(match, "package.json")
			if _, err := os.Stat(packageJSONPath); err != nil {
				continue // Skip directories without package.json
			}

			// Check cache first
			wp.cacheMutex.RLock()
			if cachedPkg, exists := wp.packageCache[packageJSONPath]; exists {
				wp.cacheMutex.RUnlock()
				packages = append(packages, cachedPkg)
				continue
			}
			wp.cacheMutex.RUnlock()

			pkg, err := wp.parseWorkspacePackage(packageJSONPath, workspace)
			if err != nil {
				continue // Skip malformed packages with graceful degradation
			}

			packages = append(packages, pkg)
		}
	}

	return packages, nil
}

// expandPackagePattern expands glob patterns to actual directory paths
func (wp *WorkspaceParser) expandPackagePattern(rootDir, pattern string) ([]string, error) {
	// Handle absolute vs relative patterns
	searchPattern := pattern
	if !filepath.IsAbs(pattern) {
		searchPattern = filepath.Join(rootDir, pattern)
	}

	// Use filepath.Glob for basic glob support
	matches, err := filepath.Glob(searchPattern)
	if err != nil {
		return nil, err
	}

	// Filter to only directories
	var dirs []string
	for _, match := range matches {
		if stat, err := os.Stat(match); err == nil && stat.IsDir() {
			dirs = append(dirs, match)
		}
	}

	// Handle nested patterns like "packages/**"
	if strings.Contains(pattern, "**") {
		nestedDirs, err := wp.expandRecursivePattern(rootDir, pattern)
		if err == nil {
			dirs = append(dirs, nestedDirs...)
		}
	}

	return dirs, nil
}

// expandRecursivePattern handles recursive glob patterns with **
func (wp *WorkspaceParser) expandRecursivePattern(rootDir, pattern string) ([]string, error) {
	var dirs []string
	
	// Convert ** pattern to regex
	regexPattern := strings.ReplaceAll(pattern, "**", ".*")
	regexPattern = strings.ReplaceAll(regexPattern, "*", "[^/]*")
	regex, err := regexp.Compile("^" + regexPattern + "$")
	if err != nil {
		return nil, err
	}

	err = filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			return nil
		}

		// Skip common directories that shouldn't contain packages
		name := d.Name()
		if name == "node_modules" || name == ".git" || name == "dist" || 
		   name == "build" || name == ".next" || name == "coverage" {
			return filepath.SkipDir
		}

		// Get relative path for pattern matching
		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return nil
		}

		if regex.MatchString(relPath) {
			// Check if this directory has a package.json
			packageJSONPath := filepath.Join(path, "package.json")
			if _, err := os.Stat(packageJSONPath); err == nil {
				dirs = append(dirs, path)
			}
		}

		return nil
	})

	return dirs, err
}

// parseWorkspacePackage parses a package.json file into a WorkspacePackage
func (wp *WorkspaceParser) parseWorkspacePackage(packageJSONPath string, workspace *WorkspaceConfiguration) (*WorkspacePackage, error) {
	data, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read package.json: %w", err)
	}

	var packageJSON struct {
		Name             string            `json:"name"`
		Version          string            `json:"version"`
		Dependencies     map[string]string `json:"dependencies"`
		DevDependencies  map[string]string `json:"devDependencies"`
		PeerDependencies map[string]string `json:"peerDependencies"`
		Private          bool              `json:"private"`
	}

	if err := json.Unmarshal(data, &packageJSON); err != nil {
		return nil, fmt.Errorf("failed to parse package.json: %w", err)
	}

	return &WorkspacePackage{
		Name:             wp.intern(packageJSON.Name),
		Path:             wp.intern(packageJSONPath),
		Version:          wp.intern(packageJSON.Version),
		Dependencies:     wp.internMap(packageJSON.Dependencies),
		DevDependencies:  wp.internMap(packageJSON.DevDependencies),
		PeerDependencies: wp.internMap(packageJSON.PeerDependencies),
		Private:          packageJSON.Private,
		Workspace:        workspace,
	}, nil
}

// deduplicatePackages removes duplicate packages based on path
func (wp *WorkspaceParser) deduplicatePackages(packages []*WorkspacePackage) []*WorkspacePackage {
	seen := make(map[string]bool)
	var unique []*WorkspacePackage

	for _, pkg := range packages {
		if !seen[pkg.Path] {
			seen[pkg.Path] = true
			unique = append(unique, pkg)
		}
	}

	return unique
}

// internSlice applies string interning to a slice of strings
func (wp *WorkspaceParser) internSlice(slice []string) []string {
	interned := make([]string, len(slice))
	for i, s := range slice {
		interned[i] = wp.intern(s)
	}
	return interned
}

// internMap applies string interning to a map[string]string
func (wp *WorkspaceParser) internMap(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}
	
	interned := make(map[string]string, len(m))
	for k, v := range m {
		interned[wp.intern(k)] = wp.intern(v)
	}
	return interned
}

// ClearCache clears the internal caches
func (wp *WorkspaceParser) ClearCache() {
	wp.cacheMutex.Lock()
	wp.cache = make(map[string]*WorkspaceConfiguration)
	wp.packageCache = make(map[string]*WorkspacePackage)
	wp.cacheMutex.Unlock()

	wp.stringInternMutex.Lock()
	wp.stringInternMap = make(map[string]string)
	wp.stringInternMutex.Unlock()
}

// GetCacheStats returns cache statistics for monitoring
func (wp *WorkspaceParser) GetCacheStats() map[string]int {
	wp.cacheMutex.RLock()
	workspaceCount := len(wp.cache)
	packageCount := len(wp.packageCache)
	wp.cacheMutex.RUnlock()

	wp.stringInternMutex.RLock()
	stringCount := len(wp.stringInternMap)
	wp.stringInternMutex.RUnlock()

	return map[string]int{
		"workspaces":     workspaceCount,
		"packages":       packageCount,
		"internedStrings": stringCount,
	}
}