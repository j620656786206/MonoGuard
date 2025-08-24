package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/monoguard/api/internal/models"
	"github.com/sirupsen/logrus"
)

// DependencyAnalyzer analyzes dependencies in monorepos
type DependencyAnalyzer struct {
	logger       *logrus.Logger
	treeResolver *DependencyTreeResolver
}

// NewDependencyAnalyzer creates a new dependency analyzer
func NewDependencyAnalyzer(logger *logrus.Logger) *DependencyAnalyzer {
	return &DependencyAnalyzer{
		logger:       logger,
		treeResolver: NewDependencyTreeResolver(logger),
	}
}

// PackageJSON represents a package.json structure
type PackageJSON struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	PeerDependencies map[string]string `json:"peerDependencies"`
	Scripts         map[string]string `json:"scripts"`
	Workspaces      []string          `json:"workspaces"`
}

// PackageInfo contains information about a discovered package
type PackageInfo struct {
	Path        string
	PackageJSON PackageJSON
}

// AnalyzeMonorepo analyzes a monorepo for dependency issues
func (da *DependencyAnalyzer) AnalyzeMonorepo(ctx context.Context, repoPath string, projectID string) (*models.DependencyAnalysisResults, error) {
	startTime := time.Now()
	da.logger.WithFields(logrus.Fields{
		"project_id": projectID,
		"repo_path":  repoPath,
	}).Info("Starting dependency analysis")

	// Discover all package.json files
	packages, err := da.discoverPackages(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to discover packages: %w", err)
	}

	da.logger.WithField("package_count", len(packages)).Info("Discovered packages")

	// Analyze duplicates
	duplicates := da.findDuplicateDependencies(packages)
	
	// Analyze version conflicts
	conflicts := da.findVersionConflicts(packages)
	
	// Analyze unused dependencies
	unused := da.findUnusedDependencies(packages, repoPath)
	
	// Find circular dependencies
	circular := da.findCircularDependencies(packages)
	
	// Calculate bundle impact
	bundleImpact := da.calculateBundleImpact(packages, duplicates, unused)
	
	// Generate summary
	summary := da.generateSummary(packages, duplicates, conflicts, unused, circular)

	duration := time.Since(startTime)
	da.logger.WithFields(logrus.Fields{
		"project_id": projectID,
		"duration":   duration,
		"duplicates": len(duplicates),
		"conflicts":  len(conflicts),
		"unused":     len(unused),
		"circular":   len(circular),
	}).Info("Dependency analysis completed")

	return &models.DependencyAnalysisResults{
		DuplicateDependencies: duplicates,
		VersionConflicts:      conflicts,
		UnusedDependencies:    unused,
		CircularDependencies:  circular,
		BundleImpact:          bundleImpact,
		Summary:               summary,
	}, nil
}

// AnalyzeMonorepoWithTreeResolver performs enhanced analysis using the tree resolver
func (da *DependencyAnalyzer) AnalyzeMonorepoWithTreeResolver(ctx context.Context, repoPath string, projectID string) (*models.DependencyAnalysisResults, *DependencyTree, error) {
	startTime := time.Now()
	da.logger.WithFields(logrus.Fields{
		"project_id": projectID,
		"repo_path":  repoPath,
	}).Info("Starting enhanced dependency analysis with tree resolver")

	// Discover all package.json files
	packages, err := da.discoverPackages(repoPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to discover packages: %w", err)
	}

	da.logger.WithField("package_count", len(packages)).Info("Discovered packages")

	// Build dependency tree using the resolver
	buildOptions := BuildOptions{
		MaxDepth:             10,
		IncludeDevDeps:       true,
		IncludePeerDeps:      true,
		IncludeOptional:      false,
		Strategy:             StrategyUpgradeAll,
		PreferWorkspace:      true,
		AllowPreRelease:      false,
		EnableCaching:        true,
		ConcurrencyLevel:     5,
		TimeoutPerPackage:    30 * time.Second,
		UseNpmRegistry:       true,
		UseLocalCache:        true,
		ConflictThreshold:    models.SeverityLow,
		AutoResolveConflicts: false,
	}

	dependencyTree, err := da.treeResolver.BuildTree(ctx, packages, buildOptions)
	if err != nil {
		da.logger.WithError(err).Warn("Failed to build dependency tree, falling back to basic analysis")
		// Fall back to basic analysis
		basicResults, fallbackErr := da.AnalyzeMonorepo(ctx, repoPath, projectID)
		return basicResults, nil, fallbackErr
	}

	da.logger.WithFields(logrus.Fields{
		"total_nodes":   dependencyTree.Metadata.TotalNodes,
		"max_depth":     dependencyTree.Metadata.MaxDepth,
		"conflicts":     len(dependencyTree.Conflicts),
	}).Info("Built dependency tree successfully")

	// Convert enhanced conflicts to basic format for compatibility
	var basicConflicts []models.VersionConflict
	for _, enhancedConflict := range dependencyTree.Conflicts {
		basicConflicts = append(basicConflicts, enhancedConflict.VersionConflict)
	}

	// Analyze other issues using existing methods
	duplicates := da.findDuplicateDependencies(packages)
	unused := da.findUnusedDependencies(packages, repoPath)
	circular := da.findCircularDependencies(packages)
	
	// Calculate bundle impact
	bundleImpact := da.calculateBundleImpact(packages, duplicates, unused)
	
	// Generate enhanced summary
	summary := da.generateEnhancedSummary(packages, duplicates, basicConflicts, unused, circular, dependencyTree)

	duration := time.Since(startTime)
	da.logger.WithFields(logrus.Fields{
		"project_id":      projectID,
		"duration":        duration,
		"duplicates":      len(duplicates),
		"tree_conflicts":  len(dependencyTree.Conflicts),
		"basic_conflicts": len(basicConflicts),
		"unused":          len(unused),
		"circular":        len(circular),
	}).Info("Enhanced dependency analysis completed")

	results := &models.DependencyAnalysisResults{
		DuplicateDependencies: duplicates,
		VersionConflicts:      basicConflicts,
		UnusedDependencies:    unused,
		CircularDependencies:  circular,
		BundleImpact:          bundleImpact,
		Summary:               summary,
	}

	return results, dependencyTree, nil
}

// discoverPackages finds all package.json files in the repository
func (da *DependencyAnalyzer) discoverPackages(rootPath string) ([]*PackageInfo, error) {
	var packages []*PackageInfo

	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip node_modules and other common directories
		if d.IsDir() {
			name := d.Name()
			if name == "node_modules" || name == ".git" || name == "dist" || name == "build" {
				return filepath.SkipDir
			}
		}

		if d.Name() == "package.json" {
			packageInfo, err := da.parsePackageJSON(path)
			if err != nil {
				da.logger.WithError(err).WithField("path", path).Warn("Failed to parse package.json")
				return nil // Continue walking
			}
			packages = append(packages, packageInfo)
		}

		return nil
	})

	return packages, err
}

// parsePackageJSON parses a package.json file
func (da *DependencyAnalyzer) parsePackageJSON(filePath string) (*PackageInfo, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}

	return &PackageInfo{
		Path:        filePath,
		PackageJSON: pkg,
	}, nil
}

// findDuplicateDependencies finds duplicate dependencies across packages
func (da *DependencyAnalyzer) findDuplicateDependencies(packages []*PackageInfo) []models.DuplicateDependency {
	dependencyMap := make(map[string]map[string][]string) // pkg -> version -> packages

	// Collect all dependencies
	for _, pkg := range packages {
		for depName, version := range pkg.PackageJSON.Dependencies {
			if dependencyMap[depName] == nil {
				dependencyMap[depName] = make(map[string][]string)
			}
			dependencyMap[depName][version] = append(dependencyMap[depName][version], pkg.PackageJSON.Name)
		}
		for depName, version := range pkg.PackageJSON.DevDependencies {
			if dependencyMap[depName] == nil {
				dependencyMap[depName] = make(map[string][]string)
			}
			dependencyMap[depName][version] = append(dependencyMap[depName][version], pkg.PackageJSON.Name)
		}
	}

	var duplicates []models.DuplicateDependency

	for depName, versions := range dependencyMap {
		if len(versions) > 1 {
			var versionList []string
			var affectedPackages []string
			
			for version, pkgs := range versions {
				versionList = append(versionList, version)
				affectedPackages = append(affectedPackages, pkgs...)
			}

			sort.Strings(versionList)
			sort.Strings(affectedPackages)

			riskLevel := da.calculateDuplicateRisk(versionList)
			
			duplicates = append(duplicates, models.DuplicateDependency{
				PackageName:      depName,
				Versions:         versionList,
				AffectedPackages: affectedPackages,
				EstimatedWaste:   da.estimateWasteSize(depName, len(versionList)),
				RiskLevel:        riskLevel,
				Recommendation:   da.generateDuplicateRecommendation(depName, versionList),
				MigrationSteps:   da.generateMigrationSteps(depName, versionList),
			})
		}
	}

	return duplicates
}

// findVersionConflicts identifies version conflicts between dependencies
func (da *DependencyAnalyzer) findVersionConflicts(packages []*PackageInfo) []models.VersionConflict {
	dependencyMap := make(map[string]map[string][]string)

	// Collect dependencies with semantic version analysis
	for _, pkg := range packages {
		for depName, version := range pkg.PackageJSON.Dependencies {
			if dependencyMap[depName] == nil {
				dependencyMap[depName] = make(map[string][]string)
			}
			dependencyMap[depName][version] = append(dependencyMap[depName][version], pkg.PackageJSON.Name)
		}
	}

	var conflicts []models.VersionConflict

	for depName, versions := range dependencyMap {
		if len(versions) > 1 {
			conflictingVersions := da.analyzeVersionConflicts(versions)
			if len(conflictingVersions) > 0 {
				conflicts = append(conflicts, models.VersionConflict{
					PackageName:         depName,
					ConflictingVersions: conflictingVersions,
					RiskLevel:           da.calculateConflictRisk(conflictingVersions),
					Resolution:          da.suggestResolution(depName, conflictingVersions),
					Impact:              da.assessConflictImpact(depName, conflictingVersions),
				})
			}
		}
	}

	return conflicts
}

// findUnusedDependencies identifies dependencies that are not being used
func (da *DependencyAnalyzer) findUnusedDependencies(packages []*PackageInfo, repoPath string) []models.UnusedDependency {
	var unused []models.UnusedDependency

	for _, pkg := range packages {
		packageDir := filepath.Dir(pkg.Path)
		
		// Check each dependency
		for depName := range pkg.PackageJSON.Dependencies {
			if !da.isDependencyUsed(packageDir, depName) {
				unused = append(unused, models.UnusedDependency{
					PackageName: depName,
					Version:     pkg.PackageJSON.Dependencies[depName],
					PackagePath: pkg.Path,
					SizeImpact:  da.estimatePackageSize(depName),
					Confidence:  da.calculateUnusedConfidence(packageDir, depName),
				})
			}
		}
	}

	return unused
}

// findCircularDependencies identifies circular dependencies
func (da *DependencyAnalyzer) findCircularDependencies(packages []*PackageInfo) []models.CircularDependency {
	var circular []models.CircularDependency
	
	// Build dependency graph
	graph := make(map[string][]string)
	for _, pkg := range packages {
		for dep := range pkg.PackageJSON.Dependencies {
			graph[pkg.PackageJSON.Name] = append(graph[pkg.PackageJSON.Name], dep)
		}
	}

	// Detect cycles using DFS
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	
	var dfs func(string, []string) []string
	dfs = func(node string, path []string) []string {
		visited[node] = true
		recStack[node] = true
		newPath := append(path, node)
		
		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				if cycle := dfs(neighbor, newPath); cycle != nil {
					return cycle
				}
			} else if recStack[neighbor] {
				// Found cycle
				cycleStart := -1
				for i, p := range newPath {
					if p == neighbor {
						cycleStart = i
						break
					}
				}
				if cycleStart != -1 {
					return newPath[cycleStart:]
				}
			}
		}
		
		recStack[node] = false
		return nil
	}

	for node := range graph {
		if !visited[node] {
			if cycle := dfs(node, []string{}); cycle != nil {
				circular = append(circular, models.CircularDependency{
					Cycle:    cycle,
					Type:     da.determineCycleType(cycle),
					Severity: da.calculateCycleSeverity(cycle),
					Impact:   da.assessCycleImpact(cycle),
				})
			}
		}
	}

	return circular
}

// Helper methods for analysis

func (da *DependencyAnalyzer) calculateDuplicateRisk(versions []string) models.RiskLevel {
	if len(versions) > 3 {
		return models.RiskLevelHigh
	} else if len(versions) > 2 {
		return models.RiskLevelMedium
	}
	return models.RiskLevelLow
}

func (da *DependencyAnalyzer) estimateWasteSize(packageName string, versionCount int) string {
	// Simplified estimation - would need real package size data
	baseSize := 100 // KB
	wasteSize := baseSize * (versionCount - 1)
	return fmt.Sprintf("%d KB", wasteSize)
}

func (da *DependencyAnalyzer) generateDuplicateRecommendation(packageName string, versions []string) string {
	latest := da.findLatestVersion(versions)
	return fmt.Sprintf("Consolidate to version %s across all packages", latest)
}

func (da *DependencyAnalyzer) generateMigrationSteps(packageName string, versions []string) []string {
	latest := da.findLatestVersion(versions)
	return []string{
		fmt.Sprintf("Update all references of %s to version %s", packageName, latest),
		"Test compatibility across all affected packages",
		"Update lock files and reinstall dependencies",
	}
}

func (da *DependencyAnalyzer) findLatestVersion(versions []string) string {
	// Simplified version comparison
	latest := versions[0]
	for _, v := range versions[1:] {
		if da.compareVersions(v, latest) > 0 {
			latest = v
		}
	}
	return latest
}

func (da *DependencyAnalyzer) compareVersions(v1, v2 string) int {
	// Simplified semantic version comparison
	v1Clean := strings.TrimPrefix(strings.TrimPrefix(v1, "^"), "~")
	v2Clean := strings.TrimPrefix(strings.TrimPrefix(v2, "^"), "~")
	
	v1Parts := strings.Split(v1Clean, ".")
	v2Parts := strings.Split(v2Clean, ".")
	
	maxLen := len(v1Parts)
	if len(v2Parts) > maxLen {
		maxLen = len(v2Parts)
	}
	
	for i := 0; i < maxLen; i++ {
		var p1, p2 int
		if i < len(v1Parts) {
			p1, _ = strconv.Atoi(v1Parts[i])
		}
		if i < len(v2Parts) {
			p2, _ = strconv.Atoi(v2Parts[i])
		}
		
		if p1 > p2 {
			return 1
		} else if p1 < p2 {
			return -1
		}
	}
	
	return 0
}

func (da *DependencyAnalyzer) analyzeVersionConflicts(versions map[string][]string) []models.ConflictingVersion {
	var conflicting []models.ConflictingVersion
	
	for version, packages := range versions {
		// Check if this version conflicts with others
		isBreaking := da.isBreakingVersionChange(version, versions)
		
		conflicting = append(conflicting, models.ConflictingVersion{
			Version:    version,
			Packages:   packages,
			IsBreaking: isBreaking,
		})
	}
	
	return conflicting
}

func (da *DependencyAnalyzer) isBreakingVersionChange(version string, allVersions map[string][]string) bool {
	// Simplified breaking change detection based on major version
	currentMajor := da.extractMajorVersion(version)
	
	for otherVersion := range allVersions {
		if otherVersion != version {
			otherMajor := da.extractMajorVersion(otherVersion)
			if currentMajor != otherMajor {
				return true
			}
		}
	}
	
	return false
}

func (da *DependencyAnalyzer) extractMajorVersion(version string) int {
	clean := strings.TrimPrefix(strings.TrimPrefix(version, "^"), "~")
	parts := strings.Split(clean, ".")
	if len(parts) > 0 {
		major, _ := strconv.Atoi(parts[0])
		return major
	}
	return 0
}

func (da *DependencyAnalyzer) calculateConflictRisk(conflicts []models.ConflictingVersion) models.RiskLevel {
	hasBreaking := false
	for _, conflict := range conflicts {
		if conflict.IsBreaking {
			hasBreaking = true
			break
		}
	}
	
	if hasBreaking {
		return models.RiskLevelHigh
	}
	return models.RiskLevelMedium
}

func (da *DependencyAnalyzer) suggestResolution(packageName string, conflicts []models.ConflictingVersion) string {
	return fmt.Sprintf("Align all packages to use the same version of %s", packageName)
}

func (da *DependencyAnalyzer) assessConflictImpact(packageName string, conflicts []models.ConflictingVersion) string {
	totalPackages := 0
	for _, conflict := range conflicts {
		totalPackages += len(conflict.Packages)
	}
	return fmt.Sprintf("Affects %d packages in the monorepo", totalPackages)
}

func (da *DependencyAnalyzer) isDependencyUsed(packageDir, depName string) bool {
	// Search for import/require statements
	found := false
	
	err := filepath.WalkDir(packageDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		if d.IsDir() && (d.Name() == "node_modules" || d.Name() == "dist" || d.Name() == "build") {
			return filepath.SkipDir
		}
		
		if strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".ts") || 
		   strings.HasSuffix(path, ".jsx") || strings.HasSuffix(path, ".tsx") {
			content, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			
			// Check for import/require statements
			patterns := []string{
				fmt.Sprintf(`import.*['"]%s['"]`, depName),
				fmt.Sprintf(`require\(['"]%s['"]\)`, depName),
				fmt.Sprintf(`from ['"]%s['"]`, depName),
			}
			
			for _, pattern := range patterns {
				matched, _ := regexp.MatchString(pattern, string(content))
				if matched {
					found = true
					return filepath.SkipAll
				}
			}
		}
		
		return nil
	})
	
	if err != nil {
		da.logger.WithError(err).Warn("Error checking dependency usage")
	}
	
	return found
}

func (da *DependencyAnalyzer) estimatePackageSize(packageName string) string {
	// Simplified size estimation
	return "50 KB" // Would need real package registry data
}

func (da *DependencyAnalyzer) calculateUnusedConfidence(packageDir, depName string) float64 {
	// Simplified confidence calculation
	return 0.8 // 80% confidence
}

func (da *DependencyAnalyzer) determineCycleType(cycle []string) string {
	if len(cycle) == 2 {
		return "direct"
	}
	return "indirect"
}

func (da *DependencyAnalyzer) calculateCycleSeverity(cycle []string) models.Severity {
	if len(cycle) <= 2 {
		return models.SeverityHigh
	}
	return models.SeverityMedium
}

func (da *DependencyAnalyzer) assessCycleImpact(cycle []string) string {
	return fmt.Sprintf("Circular dependency involving %d packages", len(cycle))
}

func (da *DependencyAnalyzer) calculateBundleImpact(packages []*PackageInfo, duplicates []models.DuplicateDependency, unused []models.UnusedDependency) models.BundleImpactReport {
	// Simplified bundle impact calculation
	totalSize := len(packages) * 100 // KB
	duplicateSize := len(duplicates) * 50
	unusedSize := len(unused) * 30
	potentialSavings := duplicateSize + unusedSize
	
	var breakdown []models.BundleBreakdown
	depCount := make(map[string]int)
	
	for _, pkg := range packages {
		for dep := range pkg.PackageJSON.Dependencies {
			depCount[dep]++
		}
	}
	
	for dep, count := range depCount {
		breakdown = append(breakdown, models.BundleBreakdown{
			PackageName: dep,
			Size:        "25 KB",
			Percentage:  float64(count) / float64(len(packages)) * 100,
			Duplicates:  count - 1,
		})
	}
	
	return models.BundleImpactReport{
		TotalSize:        fmt.Sprintf("%d KB", totalSize),
		DuplicateSize:    fmt.Sprintf("%d KB", duplicateSize),
		UnusedSize:       fmt.Sprintf("%d KB", unusedSize),
		PotentialSavings: fmt.Sprintf("%d KB", potentialSavings),
		Breakdown:        breakdown,
	}
}

func (da *DependencyAnalyzer) generateSummary(packages []*PackageInfo, duplicates []models.DuplicateDependency, conflicts []models.VersionConflict, unused []models.UnusedDependency, circular []models.CircularDependency) models.AnalysisSummary {
	totalDeps := 0
	for _, pkg := range packages {
		totalDeps += len(pkg.PackageJSON.Dependencies) + len(pkg.PackageJSON.DevDependencies)
	}
	
	// Calculate health score (0-100)
	healthScore := 100.0
	healthScore -= float64(len(duplicates)) * 5
	healthScore -= float64(len(conflicts)) * 10
	healthScore -= float64(len(unused)) * 2
	healthScore -= float64(len(circular)) * 15
	
	if healthScore < 0 {
		healthScore = 0
	}
	
	return models.AnalysisSummary{
		TotalPackages:  len(packages),
		DuplicateCount: len(duplicates),
		ConflictCount:  len(conflicts),
		UnusedCount:    len(unused),
		CircularCount:  len(circular),
		HealthScore:    healthScore,
	}
}

// generateEnhancedSummary generates an enhanced analysis summary using tree data
func (da *DependencyAnalyzer) generateEnhancedSummary(packages []*PackageInfo, duplicates []models.DuplicateDependency, conflicts []models.VersionConflict, unused []models.UnusedDependency, circular []models.CircularDependency, tree *DependencyTree) models.AnalysisSummary {
	totalDeps := 0
	for _, pkg := range packages {
		totalDeps += len(pkg.PackageJSON.Dependencies) + len(pkg.PackageJSON.DevDependencies)
	}
	
	// Enhanced health score calculation using tree metrics
	healthScore := 100.0
	
	// Penalize based on tree complexity
	if tree != nil && tree.Metadata != nil {
		// Penalty for deep dependency trees
		if tree.Metadata.MaxDepth > 5 {
			healthScore -= float64(tree.Metadata.MaxDepth-5) * 2
		}
		
		// Penalty for too many external dependencies
		externalRatio := float64(tree.Metadata.ExternalPackages) / float64(tree.Metadata.TotalNodes)
		if externalRatio > 0.7 {
			healthScore -= (externalRatio - 0.7) * 20
		}
	}
	
	// Standard penalties
	healthScore -= float64(len(duplicates)) * 5
	healthScore -= float64(len(conflicts)) * 10
	healthScore -= float64(len(unused)) * 2
	healthScore -= float64(len(circular)) * 15
	
	if healthScore < 0 {
		healthScore = 0
	}
	
	return models.AnalysisSummary{
		TotalPackages:  len(packages),
		DuplicateCount: len(duplicates),
		ConflictCount:  len(conflicts),
		UnusedCount:    len(unused),
		CircularCount:  len(circular),
		HealthScore:    healthScore,
	}
}