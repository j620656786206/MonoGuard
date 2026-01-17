// Package analyzer provides dependency graph analysis for monorepo workspaces.
// This file implements version conflict detection for Story 2.4.
package analyzer

import (
	"fmt"
	"sort"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// depInfo tracks which packages use a particular external dependency version.
type depInfo struct {
	version  string
	packages []string
	depType  string // "production", "development", "peer"
}

// ConflictDetector finds version conflicts in workspace dependencies.
type ConflictDetector struct {
	graph *types.DependencyGraph
}

// NewConflictDetector creates a new detector from a dependency graph.
// Uses the graph because it has pre-separated external dependencies.
func NewConflictDetector(graph *types.DependencyGraph) *ConflictDetector {
	return &ConflictDetector{
		graph: graph,
	}
}

// DetectConflicts finds all version conflicts across packages.
// Returns conflicts sorted by package name for deterministic output.
func (cd *ConflictDetector) DetectConflicts() []*types.VersionConflictInfo {
	if cd.graph == nil || len(cd.graph.Nodes) == 0 {
		return nil
	}

	// Collect all external dependencies: map[depName][]depInfo
	depMap := cd.collectDependencies()

	// Filter to only dependencies with 2+ different versions
	conflicts := cd.buildConflicts(depMap)

	// Sort for deterministic output
	sort.Slice(conflicts, func(i, j int) bool {
		return conflicts[i].PackageName < conflicts[j].PackageName
	})

	return conflicts
}

// collectDependencies gathers all external dependencies across workspace packages.
// Returns map[depName][]depInfo where each depInfo represents a unique version.
func (cd *ConflictDetector) collectDependencies() map[string][]depInfo {
	// depName -> version -> depInfo
	depMap := make(map[string]map[string]*depInfo)

	for pkgName, node := range cd.graph.Nodes {
		// Production dependencies
		cd.addDependencies(depMap, node.ExternalDeps, pkgName, types.DepTypeProduction)

		// Dev dependencies
		cd.addDependencies(depMap, node.ExternalDevDeps, pkgName, types.DepTypeDevelopment)

		// Peer dependencies
		cd.addDependencies(depMap, node.ExternalPeerDeps, pkgName, types.DepTypePeer)
	}

	// Convert to slice format
	result := make(map[string][]depInfo)
	for depName, versionMap := range depMap {
		infos := make([]depInfo, 0, len(versionMap))
		for _, info := range versionMap {
			infos = append(infos, *info)
		}
		// Sort by version for deterministic output
		sort.Slice(infos, func(i, j int) bool {
			return infos[i].version < infos[j].version
		})
		result[depName] = infos
	}

	return result
}

// addDependencies adds dependencies from a package to the collection map.
func (cd *ConflictDetector) addDependencies(
	depMap map[string]map[string]*depInfo,
	deps map[string]string,
	pkgName string,
	depType string,
) {
	for depName, version := range deps {
		if _, ok := depMap[depName]; !ok {
			depMap[depName] = make(map[string]*depInfo)
		}

		if existing, ok := depMap[depName][version]; ok {
			// Add package to existing version entry
			existing.packages = append(existing.packages, pkgName)
		} else {
			// Create new version entry
			depMap[depName][version] = &depInfo{
				version:  version,
				packages: []string{pkgName},
				depType:  depType,
			}
		}
	}
}

// buildConflicts creates VersionConflictInfo for dependencies with version mismatches.
func (cd *ConflictDetector) buildConflicts(depMap map[string][]depInfo) []*types.VersionConflictInfo {
	var conflicts []*types.VersionConflictInfo

	for depName, infos := range depMap {
		// Only report conflicts when there are 2+ different versions
		if len(infos) < 2 {
			continue
		}

		conflict := cd.buildConflict(depName, infos)
		if conflict != nil {
			conflicts = append(conflicts, conflict)
		}
	}

	return conflicts
}

// buildConflict creates a single VersionConflictInfo from collected version data.
func (cd *ConflictDetector) buildConflict(depName string, infos []depInfo) *types.VersionConflictInfo {
	// Extract version strings for severity calculation
	versions := make([]string, len(infos))
	for i, info := range infos {
		versions[i] = info.version
	}

	// Determine severity based on version differences
	severity := determineSeverity(versions)

	// Find highest version for resolution suggestion
	highestVersion := FindHighestVersion(versions)

	// Build conflicting versions list
	conflictingVersions := make([]*types.ConflictingVersion, len(infos))
	for i, info := range infos {
		// Sort packages for deterministic output
		pkgs := make([]string, len(info.packages))
		copy(pkgs, info.packages)
		sort.Strings(pkgs)

		// Determine if this version has breaking changes vs others
		isBreaking := isBreakingVersion(info.version, versions)

		conflictingVersions[i] = &types.ConflictingVersion{
			Version:    info.version,
			Packages:   pkgs,
			IsBreaking: isBreaking,
			DepType:    info.depType,
		}
	}

	// Sort conflicting versions by version string for deterministic output
	sort.Slice(conflictingVersions, func(i, j int) bool {
		return conflictingVersions[i].Version < conflictingVersions[j].Version
	})

	return &types.VersionConflictInfo{
		PackageName:         depName,
		ConflictingVersions: conflictingVersions,
		Severity:            severity,
		Resolution:          generateResolution(depName, severity, highestVersion),
		Impact:              generateImpact(depName, severity, len(infos)),
	}
}

// determineSeverity calculates severity based on version differences.
// Returns the highest severity based on the maximum version difference.
func determineSeverity(versions []string) types.ConflictSeverity {
	maxDiff := FindMaxDifference(versions)

	switch maxDiff {
	case VersionDifferenceMajor:
		return types.ConflictSeverityCritical
	case VersionDifferenceMinor:
		return types.ConflictSeverityWarning
	default:
		return types.ConflictSeverityInfo
	}
}

// isBreakingVersion checks if a version has a major difference from any other version.
func isBreakingVersion(version string, allVersions []string) bool {
	v1 := ParseSemVer(version)
	if v1 == nil {
		return false
	}

	for _, other := range allVersions {
		if other == version {
			continue
		}

		v2 := ParseSemVer(other)
		if v2 == nil {
			continue
		}

		if v1.Major != v2.Major {
			return true
		}
	}

	return false
}

// generateResolution suggests how to resolve the conflict based on severity.
func generateResolution(depName string, severity types.ConflictSeverity, highestVersion string) string {
	switch severity {
	case types.ConflictSeverityCritical:
		return fmt.Sprintf(
			"Major version conflict detected. Review breaking changes before upgrading all packages to %s",
			highestVersion,
		)
	case types.ConflictSeverityWarning:
		return fmt.Sprintf(
			"Consider upgrading all packages to %s for consistency",
			highestVersion,
		)
	default:
		return fmt.Sprintf(
			"Patch version difference. Safe to upgrade all packages to %s",
			highestVersion,
		)
	}
}

// generateImpact describes the impact of the conflict.
func generateImpact(depName string, severity types.ConflictSeverity, versionCount int) string {
	switch severity {
	case types.ConflictSeverityCritical:
		return fmt.Sprintf(
			"Breaking changes likely between versions. %d different versions of %s detected across workspace.",
			versionCount, depName,
		)
	case types.ConflictSeverityWarning:
		return fmt.Sprintf(
			"Minor version differences may introduce new features or behavior changes. %d versions of %s in use.",
			versionCount, depName,
		)
	default:
		return fmt.Sprintf(
			"Patch version differences only. Minimal risk. %d versions of %s in use.",
			versionCount, depName,
		)
	}
}
