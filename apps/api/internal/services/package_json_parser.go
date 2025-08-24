package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// PackageJSONParser is the main parser that orchestrates workspace discovery and version analysis
type PackageJSONParser struct {
	workspaceParser *WorkspaceParser
	versionParser   *VersionParser
	logger          *logrus.Logger
	
	// Performance optimization fields
	parseCache      map[string]*ParsedRepository
	cacheMutex      sync.RWMutex
	cacheExpiry     time.Duration
	
	// Configuration
	config *ParserConfig
}

// ParsedRepository represents a fully parsed monorepo with all workspace and version information
type ParsedRepository struct {
	RootPath             string                         `json:"rootPath"`
	WorkspaceConfigs     []*WorkspaceConfiguration      `json:"workspaceConfigs"`
	Packages             []*WorkspacePackage            `json:"packages"`
	DependencyGraph      *DependencyGraph               `json:"dependencyGraph"`
	VersionConflicts     []*VersionConflictInfo         `json:"versionConflicts"`
	DuplicateDependencies map[string]*DuplicateInfo     `json:"duplicateDependencies"`
	UnusedDependencies   []*UnusedDependencyInfo        `json:"unusedDependencies"`
	PackageManagerInfo   *PackageManagerInfo            `json:"packageManagerInfo"`
	ParseMetadata        *ParseMetadata                 `json:"parseMetadata"`
	CachedAt             time.Time                      `json:"cachedAt"`
}

// DependencyGraph represents the dependency relationships between packages
type DependencyGraph struct {
	Nodes map[string]*DependencyNode `json:"nodes"`
	Edges []*DependencyEdge          `json:"edges"`
}

// DependencyNode represents a node in the dependency graph
type DependencyNode struct {
	PackageName     string            `json:"packageName"`
	Version         string            `json:"version"`
	Path            string            `json:"path"`
	IsWorkspace     bool              `json:"isWorkspace"`
	Dependencies    []string          `json:"dependencies"`
	DevDependencies []string          `json:"devDependencies"`
	Dependents      []string          `json:"dependents"`
}

// DependencyEdge represents an edge in the dependency graph
type DependencyEdge struct {
	From         string `json:"from"`
	To           string `json:"to"`
	VersionRange string `json:"versionRange"`
	Type         string `json:"type"` // "production", "development", "peer"
}

// DuplicateInfo contains information about duplicate dependencies
type DuplicateInfo struct {
	PackageName      string                    `json:"packageName"`
	Versions         []*SemanticVersion        `json:"versions"`
	VersionRanges    []*VersionRange           `json:"versionRanges"`
	AffectedPackages []string                  `json:"affectedPackages"`
	EstimatedWaste   *ResourceWasteInfo        `json:"estimatedWaste"`
	ConsolidationPlan *ConsolidationPlan       `json:"consolidationPlan"`
}

// UnusedDependencyInfo contains information about unused dependencies
type UnusedDependencyInfo struct {
	PackageName    string             `json:"packageName"`
	Version        string             `json:"version"`
	PackagePath    string             `json:"packagePath"`
	UsageAnalysis  *UsageAnalysis     `json:"usageAnalysis"`
	RemovalImpact  *RemovalImpact     `json:"removalImpact"`
	Confidence     float64            `json:"confidence"`
}

// ResourceWasteInfo contains information about resource waste
type ResourceWasteInfo struct {
	DiskSpaceWaste    string  `json:"diskSpaceWaste"`
	BundleSizeWaste   string  `json:"bundleSizeWaste"`
	InstallTimeWaste  string  `json:"installTimeWaste"`
	MemoryWaste       string  `json:"memoryWaste"`
	WastePercentage   float64 `json:"wastePercentage"`
}

// ConsolidationPlan contains a plan for consolidating duplicate dependencies
type ConsolidationPlan struct {
	TargetVersion     *SemanticVersion        `json:"targetVersion"`
	MigrationSteps    []*MigrationStep        `json:"migrationSteps"`
	RiskAssessment    *ConflictRiskAssessment `json:"riskAssessment"`
	EstimatedEffort   string                  `json:"estimatedEffort"`
	ExpectedSavings   *ResourceWasteInfo      `json:"expectedSavings"`
}

// MigrationStep represents a step in the consolidation migration
type MigrationStep struct {
	Order       int    `json:"order"`
	Description string `json:"description"`
	PackagePath string `json:"packagePath"`
	OldVersion  string `json:"oldVersion"`
	NewVersion  string `json:"newVersion"`
	Command     string `json:"command"`
	Validation  string `json:"validation"`
}

// UsageAnalysis contains analysis of dependency usage
type UsageAnalysis struct {
	ImportStatements []string `json:"importStatements"`
	UsageLocations   []string `json:"usageLocations"`
	LastUsed         *time.Time `json:"lastUsed,omitempty"`
	UsageFrequency   int      `json:"usageFrequency"`
	IsTransitive     bool     `json:"isTransitive"`
}

// RemovalImpact contains information about the impact of removing a dependency
type RemovalImpact struct {
	AffectedFiles     []string `json:"affectedFiles"`
	BreakingChanges   []string `json:"breakingChanges"`
	AlternativePackages []string `json:"alternativePackages,omitempty"`
	SafeToRemove      bool     `json:"safeToRemove"`
}

// PackageManagerInfo contains information about the package manager configuration
type PackageManagerInfo struct {
	Type              string            `json:"type"`
	Version           string            `json:"version"`
	LockFilePresent   bool              `json:"lockFilePresent"`
	LockFilePath      string            `json:"lockFilePath"`
	WorkspaceFeatures []string          `json:"workspaceFeatures"`
	Configuration     map[string]interface{} `json:"configuration"`
}

// ParseMetadata contains metadata about the parsing operation
type ParseMetadata struct {
	ParsedAt           time.Time `json:"parsedAt"`
	Duration           time.Duration `json:"duration"`
	PackagesProcessed  int       `json:"packagesProcessed"`
	FilesScanned       int       `json:"filesScanned"`
	ErrorsEncountered  int       `json:"errorsEncountered"`
	WarningsGenerated  int       `json:"warningsGenerated"`
	ParserVersion      string    `json:"parserVersion"`
	RecoveryActions    []string  `json:"recoveryActions"`
}

// ParserConfig contains configuration for the parser
type ParserConfig struct {
	EnableCaching        bool          `json:"enableCaching"`
	CacheExpiry         time.Duration `json:"cacheExpiry"`
	MaxConcurrency      int           `json:"maxConcurrency"`
	EnableUsageAnalysis bool          `json:"enableUsageAnalysis"`
	SkipNodeModules     bool          `json:"skipNodeModules"`
	ExcludePatterns     []string      `json:"excludePatterns"`
	IncludeDevDeps      bool          `json:"includeDevDeps"`
	DeepAnalysis        bool          `json:"deepAnalysis"`
	MemoryOptimized     bool          `json:"memoryOptimized"`
}

// DefaultParserConfig returns a default parser configuration
func DefaultParserConfig() *ParserConfig {
	return &ParserConfig{
		EnableCaching:        true,
		CacheExpiry:         time.Hour * 1,
		MaxConcurrency:      10,
		EnableUsageAnalysis: true,
		SkipNodeModules:     true,
		IncludeDevDeps:      true,
		DeepAnalysis:        false,
		MemoryOptimized:     true,
		ExcludePatterns: []string{
			"node_modules",
			".git",
			"dist",
			"build",
			".next",
			"coverage",
			".nyc_output",
			"tmp",
			"temp",
		},
	}
}

// NewPackageJSONParser creates a new package.json parser
func NewPackageJSONParser(logger *logrus.Logger, config *ParserConfig) *PackageJSONParser {
	if config == nil {
		config = DefaultParserConfig()
	}

	return &PackageJSONParser{
		workspaceParser: NewWorkspaceParser(),
		versionParser:   NewVersionParser(),
		logger:          logger,
		parseCache:      make(map[string]*ParsedRepository),
		cacheExpiry:     config.CacheExpiry,
		config:          config,
	}
}

// ParseRepository parses a repository and returns comprehensive dependency analysis
func (pjp *PackageJSONParser) ParseRepository(ctx context.Context, rootPath string) (*ParsedRepository, error) {
	startTime := time.Now()
	pjp.logger.WithField("root_path", rootPath).Info("Starting repository parsing")

	// Check cache first if enabled
	if pjp.config.EnableCaching {
		if cached := pjp.getCachedResult(rootPath); cached != nil {
			pjp.logger.WithField("root_path", rootPath).Debug("Returning cached parse result")
			return cached, nil
		}
	}

	// Initialize result structure
	result := &ParsedRepository{
		RootPath:             rootPath,
		DuplicateDependencies: make(map[string]*DuplicateInfo),
		ParseMetadata: &ParseMetadata{
			ParsedAt:      time.Now(),
			ParserVersion: "1.0.0",
		},
	}

	var errors []error
	metadata := result.ParseMetadata

	// Step 1: Discover workspace configurations
	pjp.logger.Debug("Discovering workspace configurations")
	workspaceConfigs, err := pjp.workspaceParser.DiscoverWorkspaces(rootPath)
	if err != nil {
		errors = append(errors, fmt.Errorf("workspace discovery failed: %w", err))
		metadata.ErrorsEncountered++
	} else {
		result.WorkspaceConfigs = workspaceConfigs
		pjp.logger.WithField("workspace_count", len(workspaceConfigs)).Info("Discovered workspaces")
	}

	// Step 2: Discover all packages
	pjp.logger.Debug("Discovering packages")
	packages, err := pjp.workspaceParser.DiscoverPackages(workspaceConfigs)
	if err != nil {
		errors = append(errors, fmt.Errorf("package discovery failed: %w", err))
		metadata.ErrorsEncountered++
	} else {
		result.Packages = packages
		metadata.PackagesProcessed = len(packages)
		pjp.logger.WithField("package_count", len(packages)).Info("Discovered packages")
	}

	// Step 3: Build dependency graph
	pjp.logger.Debug("Building dependency graph")
	if len(packages) > 0 {
		dependencyGraph, err := pjp.buildDependencyGraph(packages)
		if err != nil {
			errors = append(errors, fmt.Errorf("dependency graph construction failed: %w", err))
			metadata.ErrorsEncountered++
		} else {
			result.DependencyGraph = dependencyGraph
		}
	}

	// Step 4: Analyze version conflicts
	pjp.logger.Debug("Analyzing version conflicts")
	versionConflicts, err := pjp.analyzeVersionConflicts(packages)
	if err != nil {
		errors = append(errors, fmt.Errorf("version conflict analysis failed: %w", err))
		metadata.ErrorsEncountered++
	} else {
		result.VersionConflicts = versionConflicts
	}

	// Step 5: Find duplicate dependencies
	pjp.logger.Debug("Finding duplicate dependencies")
	duplicates, err := pjp.findDuplicateDependencies(packages)
	if err != nil {
		errors = append(errors, fmt.Errorf("duplicate analysis failed: %w", err))
		metadata.ErrorsEncountered++
	} else {
		result.DuplicateDependencies = duplicates
	}

	// Step 6: Analyze unused dependencies (if enabled)
	if pjp.config.EnableUsageAnalysis {
		pjp.logger.Debug("Analyzing unused dependencies")
		unused, err := pjp.analyzeUnusedDependencies(ctx, packages, rootPath)
		if err != nil {
			errors = append(errors, fmt.Errorf("unused dependency analysis failed: %w", err))
			metadata.ErrorsEncountered++
		} else {
			result.UnusedDependencies = unused
		}
	}

	// Step 7: Analyze package manager configuration
	pjp.logger.Debug("Analyzing package manager configuration")
	packageManagerInfo, err := pjp.analyzePackageManagerInfo(rootPath)
	if err != nil {
		pjp.logger.WithError(err).Warn("Package manager analysis failed")
		metadata.WarningsGenerated++
	} else {
		result.PackageManagerInfo = packageManagerInfo
	}

	// Finalize metadata
	metadata.Duration = time.Since(startTime)
	result.CachedAt = time.Now()

	// Handle errors with graceful degradation
	if len(errors) > 0 {
		pjp.logger.WithField("error_count", len(errors)).Warn("Parsing completed with errors")
		for _, err := range errors {
			pjp.logger.WithError(err).Debug("Parse error details")
		}
		
		// Add recovery actions to metadata
		metadata.RecoveryActions = pjp.generateRecoveryActions(errors)
		
		// Only fail if we have no useful data
		if result.Packages == nil || len(result.Packages) == 0 {
			return nil, fmt.Errorf("parsing failed with critical errors: %v", errors)
		}
	}

	// Cache the result if enabled
	if pjp.config.EnableCaching {
		pjp.cacheResult(rootPath, result)
	}

	pjp.logger.WithFields(logrus.Fields{
		"duration":           metadata.Duration,
		"packages_processed": metadata.PackagesProcessed,
		"errors":             metadata.ErrorsEncountered,
		"warnings":           metadata.WarningsGenerated,
	}).Info("Repository parsing completed")

	return result, nil
}

// buildDependencyGraph constructs a dependency graph from packages
func (pjp *PackageJSONParser) buildDependencyGraph(packages []*WorkspacePackage) (*DependencyGraph, error) {
	graph := &DependencyGraph{
		Nodes: make(map[string]*DependencyNode),
		Edges: make([]*DependencyEdge, 0),
	}

	// Create nodes
	for _, pkg := range packages {
		node := &DependencyNode{
			PackageName:     pkg.Name,
			Version:         pkg.Version,
			Path:            pkg.Path,
			IsWorkspace:     pkg.Workspace != nil,
			Dependencies:    make([]string, 0),
			DevDependencies: make([]string, 0),
			Dependents:      make([]string, 0),
		}

		// Collect dependencies
		for depName := range pkg.Dependencies {
			node.Dependencies = append(node.Dependencies, depName)
		}
		for depName := range pkg.DevDependencies {
			node.DevDependencies = append(node.DevDependencies, depName)
		}

		graph.Nodes[pkg.Name] = node
	}

	// Create edges and populate dependents
	for _, pkg := range packages {
		// Production dependencies
		for depName, versionRange := range pkg.Dependencies {
			edge := &DependencyEdge{
				From:         pkg.Name,
				To:           depName,
				VersionRange: versionRange,
				Type:         "production",
			}
			graph.Edges = append(graph.Edges, edge)

			// Add to dependents
			if depNode, exists := graph.Nodes[depName]; exists {
				depNode.Dependents = append(depNode.Dependents, pkg.Name)
			}
		}

		// Development dependencies
		if pjp.config.IncludeDevDeps {
			for depName, versionRange := range pkg.DevDependencies {
				edge := &DependencyEdge{
					From:         pkg.Name,
					To:           depName,
					VersionRange: versionRange,
					Type:         "development",
				}
				graph.Edges = append(graph.Edges, edge)

				if depNode, exists := graph.Nodes[depName]; exists {
					depNode.Dependents = append(depNode.Dependents, pkg.Name)
				}
			}
		}

		// Peer dependencies
		for depName, versionRange := range pkg.PeerDependencies {
			edge := &DependencyEdge{
				From:         pkg.Name,
				To:           depName,
				VersionRange: versionRange,
				Type:         "peer",
			}
			graph.Edges = append(graph.Edges, edge)
		}
	}

	return graph, nil
}

// analyzeVersionConflicts analyzes version conflicts across packages
func (pjp *PackageJSONParser) analyzeVersionConflicts(packages []*WorkspacePackage) ([]*VersionConflictInfo, error) {
	// Collect all dependencies and their versions
	dependencyVersions := make(map[string][]string)

	for _, pkg := range packages {
		for depName, version := range pkg.Dependencies {
			dependencyVersions[depName] = append(dependencyVersions[depName], version)
		}

		if pjp.config.IncludeDevDeps {
			for depName, version := range pkg.DevDependencies {
				dependencyVersions[depName] = append(dependencyVersions[depName], version)
			}
		}
	}

	return pjp.versionParser.AnalyzeVersionConflicts(dependencyVersions)
}

// findDuplicateDependencies finds and analyzes duplicate dependencies
func (pjp *PackageJSONParser) findDuplicateDependencies(packages []*WorkspacePackage) (map[string]*DuplicateInfo, error) {
	duplicates := make(map[string]*DuplicateInfo)
	
	// Collect dependency usage
	dependencyUsage := make(map[string]map[string][]string) // dep -> version -> packages

	for _, pkg := range packages {
		for depName, version := range pkg.Dependencies {
			if dependencyUsage[depName] == nil {
				dependencyUsage[depName] = make(map[string][]string)
			}
			dependencyUsage[depName][version] = append(dependencyUsage[depName][version], pkg.Name)
		}

		if pjp.config.IncludeDevDeps {
			for depName, version := range pkg.DevDependencies {
				if dependencyUsage[depName] == nil {
					dependencyUsage[depName] = make(map[string][]string)
				}
				dependencyUsage[depName][version] = append(dependencyUsage[depName][version], pkg.Name)
			}
		}
	}

	// Identify duplicates
	for depName, versions := range dependencyUsage {
		if len(versions) > 1 {
			duplicateInfo, err := pjp.analyzeDuplicateInfo(depName, versions)
			if err != nil {
				pjp.logger.WithError(err).WithField("package", depName).Warn("Failed to analyze duplicate")
				continue
			}
			duplicates[depName] = duplicateInfo
		}
	}

	return duplicates, nil
}

// analyzeDuplicateInfo creates detailed duplicate information
func (pjp *PackageJSONParser) analyzeDuplicateInfo(packageName string, versionUsage map[string][]string) (*DuplicateInfo, error) {
	var versions []*SemanticVersion
	var versionRanges []*VersionRange
	var affectedPackages []string

	// Parse versions and ranges
	for versionStr, pkgs := range versionUsage {
		if version, err := pjp.versionParser.ParseVersion(versionStr); err == nil {
			versions = append(versions, version)
		}
		if versionRange, err := pjp.versionParser.ParseVersionRange(versionStr); err == nil {
			versionRanges = append(versionRanges, versionRange)
		}
		affectedPackages = append(affectedPackages, pkgs...)
	}

	// Sort versions for consistency
	sort.Slice(versions, func(i, j int) bool {
		return pjp.versionParser.compareVersions(versions[i], versions[j]) < 0
	})

	// Estimate waste
	estimatedWaste := pjp.estimateResourceWaste(packageName, len(versions))

	// Create consolidation plan
	consolidationPlan := pjp.createConsolidationPlan(packageName, versions, versionRanges)

	return &DuplicateInfo{
		PackageName:       packageName,
		Versions:          versions,
		VersionRanges:     versionRanges,
		AffectedPackages:  affectedPackages,
		EstimatedWaste:    estimatedWaste,
		ConsolidationPlan: consolidationPlan,
	}, nil
}

// estimateResourceWaste estimates the resource waste from duplicate dependencies
func (pjp *PackageJSONParser) estimateResourceWaste(packageName string, versionCount int) *ResourceWasteInfo {
	// These are simplified estimates - would be enhanced with real package registry data
	baseDiskSize := 500 * 1024 // 500KB base estimate
	baseBundleSize := 100 * 1024 // 100KB bundle estimate
	
	wasteCount := versionCount - 1
	diskWaste := baseDiskSize * wasteCount
	bundleWaste := baseBundleSize * wasteCount

	return &ResourceWasteInfo{
		DiskSpaceWaste:   formatBytes(diskWaste),
		BundleSizeWaste:  formatBytes(bundleWaste),
		InstallTimeWaste: fmt.Sprintf("%ds", wasteCount*2),
		MemoryWaste:      formatBytes(bundleWaste / 2),
		WastePercentage:  float64(wasteCount) / float64(versionCount) * 100,
	}
}

// createConsolidationPlan creates a plan to consolidate duplicate dependencies
func (pjp *PackageJSONParser) createConsolidationPlan(packageName string, versions []*SemanticVersion, ranges []*VersionRange) *ConsolidationPlan {
	targetVersion := pjp.versionParser.suggestVersionFix(versions, ranges)
	
	var migrationSteps []*MigrationStep
	stepOrder := 1

	// Create migration steps for each version that needs to change
	for _, version := range versions {
		if pjp.versionParser.compareVersions(version, targetVersion) != 0 {
			step := &MigrationStep{
				Order:       stepOrder,
				Description: fmt.Sprintf("Update %s from %s to %s", packageName, version.String(), targetVersion.String()),
				OldVersion:  version.String(),
				NewVersion:  targetVersion.String(),
				Command:     fmt.Sprintf("npm install %s@%s", packageName, targetVersion.String()),
				Validation:  fmt.Sprintf("npm list %s", packageName),
			}
			migrationSteps = append(migrationSteps, step)
			stepOrder++
		}
	}

	// Assess risk
	conflictType := "minor" // Default assumption
	if len(versions) > 0 && targetVersion != nil {
		if versions[0].Major != targetVersion.Major {
			conflictType = "major"
		}
	}

	riskAssessment := pjp.versionParser.assessConflictRisk(packageName, versions, conflictType)
	
	// Estimate effort
	effortLevel := "Low"
	if len(migrationSteps) > 5 {
		effortLevel = "Medium"
	}
	if riskAssessment.Level == "critical" || riskAssessment.Level == "high" {
		effortLevel = "High"
	}

	return &ConsolidationPlan{
		TargetVersion:   targetVersion,
		MigrationSteps:  migrationSteps,
		RiskAssessment:  riskAssessment,
		EstimatedEffort: effortLevel,
		ExpectedSavings: pjp.estimateResourceWaste(packageName, len(versions)),
	}
}

// analyzeUnusedDependencies analyzes unused dependencies (simplified implementation)
func (pjp *PackageJSONParser) analyzeUnusedDependencies(ctx context.Context, packages []*WorkspacePackage, rootPath string) ([]*UnusedDependencyInfo, error) {
	var unused []*UnusedDependencyInfo

	// This would be enhanced with actual usage analysis
	// For now, provide a basic structure
	for _, pkg := range packages {
		packageDir := filepath.Dir(pkg.Path)
		
		for depName, version := range pkg.Dependencies {
			// Simplified usage check (would be enhanced)
			if !pjp.isDependencyUsed(packageDir, depName) {
				unusedInfo := &UnusedDependencyInfo{
					PackageName: depName,
					Version:     version,
					PackagePath: pkg.Path,
					UsageAnalysis: &UsageAnalysis{
						ImportStatements: []string{},
						UsageLocations:   []string{},
						UsageFrequency:   0,
						IsTransitive:     false,
					},
					RemovalImpact: &RemovalImpact{
						AffectedFiles:   []string{},
						BreakingChanges: []string{},
						SafeToRemove:    true,
					},
					Confidence: 0.8,
				}
				unused = append(unused, unusedInfo)
			}
		}
	}

	return unused, nil
}

// isDependencyUsed checks if a dependency is used (simplified implementation)
func (pjp *PackageJSONParser) isDependencyUsed(packageDir, depName string) bool {
	// This is a simplified implementation
	// In a full implementation, this would scan source files for import/require statements
	return true // Always return true for now to avoid false positives
}

// analyzePackageManagerInfo analyzes package manager configuration
func (pjp *PackageJSONParser) analyzePackageManagerInfo(rootPath string) (*PackageManagerInfo, error) {
	info := &PackageManagerInfo{
		Configuration:     make(map[string]interface{}),
		WorkspaceFeatures: make([]string, 0),
	}

	// Check for different lock files to determine package manager
	lockFiles := map[string]string{
		"package-lock.json": "npm",
		"yarn.lock":         "yarn",
		"pnpm-lock.yaml":    "pnpm",
	}

	for lockFile, pmType := range lockFiles {
		lockPath := filepath.Join(rootPath, lockFile)
		if _, err := os.Stat(lockPath); err == nil {
			info.Type = pmType
			info.LockFilePresent = true
			info.LockFilePath = lockPath
			break
		}
	}

	// Default to npm if no lock file found
	if info.Type == "" {
		info.Type = "npm"
	}

	return info, nil
}

// generateRecoveryActions generates recovery actions for errors
func (pjp *PackageJSONParser) generateRecoveryActions(errors []error) []string {
	var actions []string

	for _, err := range errors {
		errStr := err.Error()
		switch {
		case strings.Contains(errStr, "workspace discovery"):
			actions = append(actions, "Manual workspace configuration validation recommended")
		case strings.Contains(errStr, "package discovery"):
			actions = append(actions, "Check package.json file permissions and syntax")
		case strings.Contains(errStr, "version conflict"):
			actions = append(actions, "Manual version range validation needed")
		default:
			actions = append(actions, "Review error logs for specific issues")
		}
	}

	return actions
}

// Cache management methods
func (pjp *PackageJSONParser) getCachedResult(rootPath string) *ParsedRepository {
	pjp.cacheMutex.RLock()
	defer pjp.cacheMutex.RUnlock()

	if cached, exists := pjp.parseCache[rootPath]; exists {
		if time.Since(cached.CachedAt) < pjp.cacheExpiry {
			return cached
		}
		// Cache expired, remove it
		delete(pjp.parseCache, rootPath)
	}

	return nil
}

func (pjp *PackageJSONParser) cacheResult(rootPath string, result *ParsedRepository) {
	pjp.cacheMutex.Lock()
	defer pjp.cacheMutex.Unlock()

	pjp.parseCache[rootPath] = result
}

// ClearCache clears all caches
func (pjp *PackageJSONParser) ClearCache() {
	pjp.cacheMutex.Lock()
	defer pjp.cacheMutex.Unlock()

	pjp.parseCache = make(map[string]*ParsedRepository)
	pjp.workspaceParser.ClearCache()
}

// formatBytes formats bytes as human readable string
func formatBytes(bytes int) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}