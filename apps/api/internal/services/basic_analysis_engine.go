package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/monoguard/api/internal/models"
	"github.com/sirupsen/logrus"
)

// BasicAnalysisEngine provides a focused analysis engine for mono-guard
type BasicAnalysisEngine struct {
	logger              *logrus.Logger
	packageParser       *MonoGuardPackageParser
	duplicateDetector   *MonoGuardDuplicateDetector
	conflictAnalyzer    *MonoGuardConflictAnalyzer
	reportGenerator     *MonoGuardReportGenerator
}

// NewBasicAnalysisEngine creates a new basic analysis engine
func NewBasicAnalysisEngine(logger *logrus.Logger) *BasicAnalysisEngine {
	packageParser := NewMonoGuardPackageParser(logger)
	duplicateDetector := NewMonoGuardDuplicateDetector(logger)
	conflictAnalyzer := NewMonoGuardConflictAnalyzer(logger)
	reportGenerator := NewMonoGuardReportGenerator(logger)

	return &BasicAnalysisEngine{
		logger:            logger,
		packageParser:     packageParser,
		duplicateDetector: duplicateDetector,
		conflictAnalyzer:  conflictAnalyzer,
		reportGenerator:   reportGenerator,
	}
}

// AnalyzeRepository performs a complete analysis of a repository
func (bae *BasicAnalysisEngine) AnalyzeRepository(ctx context.Context, rootPath string, projectID string) (*models.DependencyAnalysisResults, error) {
	startTime := time.Now()
	bae.logger.WithFields(logrus.Fields{
		"project_id": projectID,
		"root_path":  rootPath,
	}).Info("Starting basic analysis engine")

	// Step 1: Parse all package.json files
	packages, err := bae.packageParser.ParseRepository(ctx, rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse repository packages: %w", err)
	}

	bae.logger.WithField("package_count", len(packages)).Info("Parsed packages successfully")

	// Step 2: Detect duplicate dependencies
	duplicates, err := bae.duplicateDetector.FindDuplicates(packages)
	if err != nil {
		return nil, fmt.Errorf("failed to detect duplicates: %w", err)
	}

	bae.logger.WithField("duplicate_count", len(duplicates)).Info("Detected duplicate dependencies")

	// Step 3: Analyze version conflicts
	conflicts, err := bae.conflictAnalyzer.FindConflicts(packages)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze version conflicts: %w", err)
	}

	bae.logger.WithField("conflict_count", len(conflicts)).Info("Analyzed version conflicts")

	// Step 4: Generate comprehensive report
	results, err := bae.reportGenerator.GenerateReport(packages, duplicates, conflicts)
	if err != nil {
		return nil, fmt.Errorf("failed to generate analysis report: %w", err)
	}

	duration := time.Since(startTime)
	bae.logger.WithFields(logrus.Fields{
		"project_id": projectID,
		"duration":   duration,
		"packages":   len(packages),
		"duplicates": len(duplicates),
		"conflicts":  len(conflicts),
	}).Info("Basic analysis engine completed successfully")

	return results, nil
}

// MonoGuardPackageParser handles parsing package.json files
type MonoGuardPackageParser struct {
	logger *logrus.Logger
}

// NewMonoGuardPackageParser creates a new package parser
func NewMonoGuardPackageParser(logger *logrus.Logger) *MonoGuardPackageParser {
	return &MonoGuardPackageParser{
		logger: logger,
	}
}

// MonoGuardPackageInfo represents a parsed package.json file
type MonoGuardPackageInfo struct {
	Name             string            `json:"name"`
	Version          string            `json:"version"`
	Path             string            `json:"path"`
	Dependencies     map[string]string `json:"dependencies"`
	DevDependencies  map[string]string `json:"devDependencies"`
	PeerDependencies map[string]string `json:"peerDependencies"`
}

// ParseRepository discovers and parses all package.json files in a repository
func (mgpp *MonoGuardPackageParser) ParseRepository(ctx context.Context, rootPath string) ([]*MonoGuardPackageInfo, error) {
	var packages []*MonoGuardPackageInfo

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip common directories that shouldn't contain workspace packages
		if info.IsDir() {
			dirName := filepath.Base(path)
			if dirName == "node_modules" || dirName == ".git" || dirName == "dist" || 
			   dirName == "build" || dirName == ".next" || dirName == "coverage" {
				return filepath.SkipDir
			}
		}

		// Process package.json files
		if info.Name() == "package.json" {
			pkg, parseErr := mgpp.parsePackageJSON(path)
			if parseErr != nil {
				mgpp.logger.WithError(parseErr).WithField("path", path).Warn("Failed to parse package.json")
				return nil // Continue processing other files
			}
			packages = append(packages, pkg)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking repository: %w", err)
	}

	return packages, nil
}

// parsePackageJSON parses a single package.json file
func (mgpp *MonoGuardPackageParser) parsePackageJSON(filePath string) (*MonoGuardPackageInfo, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var rawPackage map[string]interface{}
	if err := json.Unmarshal(data, &rawPackage); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	pkg := &MonoGuardPackageInfo{
		Path:             filePath,
		Dependencies:     make(map[string]string),
		DevDependencies:  make(map[string]string),
		PeerDependencies: make(map[string]string),
	}

	// Extract basic package information
	if name, ok := rawPackage["name"].(string); ok {
		pkg.Name = name
	}

	if version, ok := rawPackage["version"].(string); ok {
		pkg.Version = version
	}

	// Parse dependencies
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

	return pkg, nil
}

// MonoGuardDuplicateDetector finds duplicate dependencies across packages
type MonoGuardDuplicateDetector struct {
	logger *logrus.Logger
}

// NewMonoGuardDuplicateDetector creates a new duplicate dependency detector
func NewMonoGuardDuplicateDetector(logger *logrus.Logger) *MonoGuardDuplicateDetector {
	return &MonoGuardDuplicateDetector{
		logger: logger,
	}
}

// FindDuplicates identifies duplicate dependencies across packages
func (mgdd *MonoGuardDuplicateDetector) FindDuplicates(packages []*MonoGuardPackageInfo) ([]models.DuplicateDependency, error) {
	// Map to track dependency usage: dependencyName -> version -> []packageNames
	dependencyMap := make(map[string]map[string][]string)

	// Collect all dependencies
	for _, pkg := range packages {
		mgdd.collectDependencies(pkg, dependencyMap, pkg.Dependencies, "production")
		mgdd.collectDependencies(pkg, dependencyMap, pkg.DevDependencies, "development")
	}

	var duplicates []models.DuplicateDependency

	// Find dependencies with multiple versions
	for depName, versions := range dependencyMap {
		if len(versions) > 1 {
			duplicate := mgdd.createDuplicateInfo(depName, versions)
			duplicates = append(duplicates, duplicate)
		}
	}

	// Sort by package name for consistent output
	sort.Slice(duplicates, func(i, j int) bool {
		return duplicates[i].PackageName < duplicates[j].PackageName
	})

	return duplicates, nil
}

// collectDependencies collects dependencies from a package into the dependency map
func (mgdd *MonoGuardDuplicateDetector) collectDependencies(pkg *MonoGuardPackageInfo, depMap map[string]map[string][]string, deps map[string]string, depType string) {
	for depName, version := range deps {
		if depMap[depName] == nil {
			depMap[depName] = make(map[string][]string)
		}
		
		packageName := pkg.Name
		if packageName == "" {
			packageName = filepath.Base(filepath.Dir(pkg.Path))
		}
		
		depMap[depName][version] = append(depMap[depName][version], packageName)
	}
}

// createDuplicateInfo creates duplicate dependency information
func (mgdd *MonoGuardDuplicateDetector) createDuplicateInfo(depName string, versions map[string][]string) models.DuplicateDependency {
	var versionList []string
	var affectedPackages []string

	for version, packages := range versions {
		versionList = append(versionList, version)
		affectedPackages = append(affectedPackages, packages...)
	}

	sort.Strings(versionList)
	sort.Strings(affectedPackages)

	// Remove duplicates from affected packages
	affectedPackages = mgdd.removeDuplicateStrings(affectedPackages)

	riskLevel := mgdd.calculateRiskLevel(len(versionList))
	estimatedWaste := mgdd.estimateWaste(len(versionList))
	recommendation := mgdd.generateRecommendation(depName, versionList)
	migrationSteps := mgdd.generateMigrationSteps(depName, versionList)

	return models.DuplicateDependency{
		PackageName:      depName,
		Versions:         versionList,
		AffectedPackages: affectedPackages,
		EstimatedWaste:   estimatedWaste,
		RiskLevel:        riskLevel,
		Recommendation:   recommendation,
		MigrationSteps:   migrationSteps,
	}
}

// removeDuplicateStrings removes duplicate strings from a slice
func (mgdd *MonoGuardDuplicateDetector) removeDuplicateStrings(strings []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, str := range strings {
		if !seen[str] {
			seen[str] = true
			result = append(result, str)
		}
	}

	return result
}

// calculateRiskLevel calculates the risk level based on the number of versions
func (mgdd *MonoGuardDuplicateDetector) calculateRiskLevel(versionCount int) models.RiskLevel {
	switch {
	case versionCount >= 4:
		return models.RiskLevelCritical
	case versionCount >= 3:
		return models.RiskLevelHigh
	case versionCount >= 2:
		return models.RiskLevelMedium
	default:
		return models.RiskLevelLow
	}
}

// estimateWaste estimates resource waste from duplicate dependencies
func (mgdd *MonoGuardDuplicateDetector) estimateWaste(versionCount int) string {
	// Simple estimation: each additional version adds ~1MB of waste
	wasteKB := (versionCount - 1) * 1024
	if wasteKB >= 1024 {
		return fmt.Sprintf("%.1f MB", float64(wasteKB)/1024)
	}
	return fmt.Sprintf("%d KB", wasteKB)
}

// generateRecommendation generates a recommendation for fixing duplicate dependencies
func (mgdd *MonoGuardDuplicateDetector) generateRecommendation(depName string, versions []string) string {
	latestVersion := mgdd.findLatestVersion(versions)
	return fmt.Sprintf("Consolidate all usage of '%s' to version %s", depName, latestVersion)
}

// generateMigrationSteps generates migration steps for fixing duplicates
func (mgdd *MonoGuardDuplicateDetector) generateMigrationSteps(depName string, versions []string) []string {
	latestVersion := mgdd.findLatestVersion(versions)
	return []string{
		fmt.Sprintf("Review all usage of %s across packages", depName),
		fmt.Sprintf("Update package.json files to use %s version %s", depName, latestVersion),
		"Run dependency installation in all affected packages",
		"Test all affected packages for compatibility",
		"Update lock files and commit changes",
	}
}

// findLatestVersion finds the latest version (simplified semantic version comparison)
func (mgdd *MonoGuardDuplicateDetector) findLatestVersion(versions []string) string {
	if len(versions) == 0 {
		return ""
	}

	// Simple approach: return the last version when sorted lexically
	// In a production system, this would use proper semantic version comparison
	sortedVersions := make([]string, len(versions))
	copy(sortedVersions, versions)
	sort.Strings(sortedVersions)
	
	return sortedVersions[len(sortedVersions)-1]
}

// MonoGuardConflictAnalyzer analyzes version conflicts between dependencies
type MonoGuardConflictAnalyzer struct {
	logger *logrus.Logger
}

// NewMonoGuardConflictAnalyzer creates a new version conflict analyzer
func NewMonoGuardConflictAnalyzer(logger *logrus.Logger) *MonoGuardConflictAnalyzer {
	return &MonoGuardConflictAnalyzer{
		logger: logger,
	}
}

// FindConflicts identifies version conflicts between dependencies
func (mgca *MonoGuardConflictAnalyzer) FindConflicts(packages []*MonoGuardPackageInfo) ([]models.VersionConflict, error) {
	// Map to track dependency usage: dependencyName -> version -> []packageNames
	dependencyMap := make(map[string]map[string][]string)

	// Collect all dependencies (focusing on production dependencies for conflicts)
	for _, pkg := range packages {
		for depName, version := range pkg.Dependencies {
			if dependencyMap[depName] == nil {
				dependencyMap[depName] = make(map[string][]string)
			}
			
			packageName := pkg.Name
			if packageName == "" {
				packageName = filepath.Base(filepath.Dir(pkg.Path))
			}
			
			dependencyMap[depName][version] = append(dependencyMap[depName][version], packageName)
		}
	}

	var conflicts []models.VersionConflict

	// Find dependencies with potential version conflicts
	for depName, versions := range dependencyMap {
		if len(versions) > 1 {
			if mgca.hasVersionConflict(versions) {
				conflict := mgca.createVersionConflict(depName, versions)
				conflicts = append(conflicts, conflict)
			}
		}
	}

	// Sort by package name for consistent output
	sort.Slice(conflicts, func(i, j int) bool {
		return conflicts[i].PackageName < conflicts[j].PackageName
	})

	return conflicts, nil
}

// hasVersionConflict determines if the versions represent a real conflict
func (mgca *MonoGuardConflictAnalyzer) hasVersionConflict(versions map[string][]string) bool {
	// Check if we have incompatible major versions
	majorVersions := make(map[string]bool)
	
	for version := range versions {
		major := mgca.extractMajorVersion(version)
		majorVersions[major] = true
	}

	// Conflict exists if we have different major versions
	return len(majorVersions) > 1
}

// extractMajorVersion extracts the major version from a version string
func (mgca *MonoGuardConflictAnalyzer) extractMajorVersion(version string) string {
	// Clean version string (remove ^, ~, >=, etc.)
	cleanVersion := strings.TrimLeft(version, "^~>=<")
	
	parts := strings.Split(cleanVersion, ".")
	if len(parts) > 0 {
		return parts[0]
	}
	
	return "0"
}

// createVersionConflict creates version conflict information
func (mgca *MonoGuardConflictAnalyzer) createVersionConflict(depName string, versions map[string][]string) models.VersionConflict {
	var conflictingVersions []models.ConflictingVersion

	for version, packages := range versions {
		isBreaking := mgca.isBreakingVersion(version, versions)
		conflictingVersions = append(conflictingVersions, models.ConflictingVersion{
			Version:    version,
			Packages:   packages,
			IsBreaking: isBreaking,
		})
	}

	// Sort conflicting versions for consistency
	sort.Slice(conflictingVersions, func(i, j int) bool {
		return conflictingVersions[i].Version < conflictingVersions[j].Version
	})

	riskLevel := mgca.calculateConflictRiskLevel(conflictingVersions)
	resolution := mgca.generateResolution(depName, conflictingVersions)
	impact := mgca.assessImpact(conflictingVersions)

	return models.VersionConflict{
		PackageName:         depName,
		ConflictingVersions: conflictingVersions,
		RiskLevel:          riskLevel,
		Resolution:         resolution,
		Impact:             impact,
	}
}

// isBreakingVersion determines if a version is likely to be breaking compared to others
func (mgca *MonoGuardConflictAnalyzer) isBreakingVersion(version string, allVersions map[string][]string) bool {
	majorVersion := mgca.extractMajorVersion(version)
	
	// Check if there are other major versions
	for otherVersion := range allVersions {
		if otherVersion != version {
			otherMajor := mgca.extractMajorVersion(otherVersion)
			if majorVersion != otherMajor {
				return true
			}
		}
	}
	
	return false
}

// calculateConflictRiskLevel calculates the risk level for version conflicts
func (mgca *MonoGuardConflictAnalyzer) calculateConflictRiskLevel(conflicts []models.ConflictingVersion) models.RiskLevel {
	hasBreaking := false
	totalPackages := 0
	
	for _, conflict := range conflicts {
		if conflict.IsBreaking {
			hasBreaking = true
		}
		totalPackages += len(conflict.Packages)
	}

	switch {
	case hasBreaking && totalPackages > 5:
		return models.RiskLevelCritical
	case hasBreaking:
		return models.RiskLevelHigh
	case totalPackages > 3:
		return models.RiskLevelMedium
	default:
		return models.RiskLevelLow
	}
}

// generateResolution generates a resolution recommendation
func (mgca *MonoGuardConflictAnalyzer) generateResolution(depName string, conflicts []models.ConflictingVersion) string {
	if len(conflicts) == 0 {
		return "No action needed"
	}

	// Find the most recent version
	latestVersion := conflicts[0].Version
	for _, conflict := range conflicts[1:] {
		if mgca.isVersionNewer(conflict.Version, latestVersion) {
			latestVersion = conflict.Version
		}
	}

	return fmt.Sprintf("Align all packages to use %s version %s and test for compatibility", depName, latestVersion)
}

// isVersionNewer compares two versions (simplified)
func (mgca *MonoGuardConflictAnalyzer) isVersionNewer(v1, v2 string) bool {
	// Simplified comparison - in production would use proper semver parsing
	return strings.Compare(v1, v2) > 0
}

// assessImpact assesses the impact of version conflicts
func (mgca *MonoGuardConflictAnalyzer) assessImpact(conflicts []models.ConflictingVersion) string {
	totalPackages := 0
	breakingCount := 0

	for _, conflict := range conflicts {
		totalPackages += len(conflict.Packages)
		if conflict.IsBreaking {
			breakingCount++
		}
	}

	if breakingCount > 0 {
		return fmt.Sprintf("High impact: affects %d packages with %d potentially breaking changes", totalPackages, breakingCount)
	}

	return fmt.Sprintf("Medium impact: affects %d packages with version alignment needed", totalPackages)
}

// MonoGuardReportGenerator generates comprehensive analysis reports
type MonoGuardReportGenerator struct {
	logger *logrus.Logger
}

// NewMonoGuardReportGenerator creates a new analysis report generator
func NewMonoGuardReportGenerator(logger *logrus.Logger) *MonoGuardReportGenerator {
	return &MonoGuardReportGenerator{
		logger: logger,
	}
}

// GenerateReport generates a comprehensive analysis report
func (mgrg *MonoGuardReportGenerator) GenerateReport(packages []*MonoGuardPackageInfo, duplicates []models.DuplicateDependency, conflicts []models.VersionConflict) (*models.DependencyAnalysisResults, error) {
	// Generate bundle impact report
	bundleImpact := mgrg.generateBundleImpact(packages, duplicates)
	
	// Generate summary
	summary := mgrg.generateSummary(packages, duplicates, conflicts)
	
	results := &models.DependencyAnalysisResults{
		DuplicateDependencies: duplicates,
		VersionConflicts:      conflicts,
		UnusedDependencies:    []models.UnusedDependency{}, // Not implemented in basic engine
		CircularDependencies:  []models.CircularDependency{}, // Not implemented in basic engine
		BundleImpact:          bundleImpact,
		Summary:               summary,
	}

	mgrg.logger.Info("Generated comprehensive analysis report")
	
	return results, nil
}

// generateBundleImpact generates bundle impact analysis
func (mgrg *MonoGuardReportGenerator) generateBundleImpact(packages []*MonoGuardPackageInfo, duplicates []models.DuplicateDependency) models.BundleImpactReport {
	// Simple bundle impact calculation
	totalPackages := len(packages)
	estimatedTotalSize := totalPackages * 2048 // 2MB per package estimate
	
	duplicateWasteKB := len(duplicates) * 1024 // 1MB per duplicate
	totalSizeKB := estimatedTotalSize
	
	var breakdown []models.BundleBreakdown
	dependencyCount := make(map[string]int)
	
	// Count dependency usage
	for _, pkg := range packages {
		for depName := range pkg.Dependencies {
			dependencyCount[depName]++
		}
	}
	
	// Create breakdown for top dependencies
	type depCount struct {
		name  string
		count int
	}
	
	var deps []depCount
	for name, count := range dependencyCount {
		deps = append(deps, depCount{name: name, count: count})
	}
	
	sort.Slice(deps, func(i, j int) bool {
		return deps[i].count > deps[j].count
	})
	
	// Top 10 most used dependencies
	maxBreakdown := 10
	if len(deps) < maxBreakdown {
		maxBreakdown = len(deps)
	}
	
	for i := 0; i < maxBreakdown; i++ {
		dep := deps[i]
		duplicateCount := 0
		
		// Check if this dependency has duplicates
		for _, duplicate := range duplicates {
			if duplicate.PackageName == dep.name {
				duplicateCount = len(duplicate.Versions) - 1
				break
			}
		}
		
		breakdown = append(breakdown, models.BundleBreakdown{
			PackageName: dep.name,
			Size:        "128 KB", // Estimated size
			Percentage:  float64(dep.count) / float64(totalPackages) * 100,
			Duplicates:  duplicateCount,
		})
	}
	
	return models.BundleImpactReport{
		TotalSize:        fmt.Sprintf("%.1f MB", float64(totalSizeKB)/1024),
		DuplicateSize:    fmt.Sprintf("%.1f MB", float64(duplicateWasteKB)/1024),
		UnusedSize:       "0 KB", // Not calculated in basic engine
		PotentialSavings: fmt.Sprintf("%.1f MB", float64(duplicateWasteKB)/1024),
		Breakdown:        breakdown,
	}
}

// generateSummary generates analysis summary
func (mgrg *MonoGuardReportGenerator) generateSummary(packages []*MonoGuardPackageInfo, duplicates []models.DuplicateDependency, conflicts []models.VersionConflict) models.AnalysisSummary {
	totalPackages := len(packages)
	duplicateCount := len(duplicates)
	conflictCount := len(conflicts)
	
	// Calculate health score (0-100)
	healthScore := 100.0
	
	// Penalties for issues
	healthScore -= float64(duplicateCount) * 5.0  // 5 points per duplicate
	healthScore -= float64(conflictCount) * 10.0  // 10 points per conflict
	
	// Ensure health score doesn't go below 0
	if healthScore < 0 {
		healthScore = 0
	}
	
	return models.AnalysisSummary{
		TotalPackages:  totalPackages,
		DuplicateCount: duplicateCount,
		ConflictCount:  conflictCount,
		UnusedCount:    0, // Not implemented in basic engine
		CircularCount:  0, // Not implemented in basic engine
		HealthScore:    healthScore,
	}
}