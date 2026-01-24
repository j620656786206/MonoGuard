// Package analyzer provides analysis functionality for mono-guard.
// This file contains the impact analyzer for Story 3.6.
package analyzer

import (
	"fmt"
	"strings"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// ========================================
// Impact Analyzer (Story 3.6)
// ========================================

// ImpactAnalyzer calculates the blast radius of circular dependencies.
type ImpactAnalyzer struct {
	graph     *types.DependencyGraph
	workspace *types.WorkspaceData
	// reverseDeps maps package -> packages that depend on it
	reverseDeps map[string][]string
}

// NewImpactAnalyzer creates a new analyzer with the given graph and workspace.
func NewImpactAnalyzer(graph *types.DependencyGraph, workspace *types.WorkspaceData) *ImpactAnalyzer {
	return &ImpactAnalyzer{
		graph:     graph,
		workspace: workspace,
	}
}

// Analyze calculates the impact assessment for a cycle.
func (ia *ImpactAnalyzer) Analyze(cycle *types.CircularDependencyInfo) *types.ImpactAssessment {
	if cycle == nil || len(cycle.Cycle) == 0 {
		return types.NewImpactAssessment()
	}

	// Build reverse dependencies if not done
	if ia.reverseDeps == nil {
		ia.buildReverseDependencies()
	}

	// Get direct participants
	directParticipants := ia.getDirectParticipants(cycle)

	// Find indirect dependents
	indirectDependents := ia.findIndirectDependents(directParticipants)

	// Calculate totals
	totalAffected := len(directParticipants) + len(indirectDependents)
	totalPackages := len(ia.graph.Nodes)

	// Calculate percentage
	percentage, displayPercentage := types.CalculatePercentage(totalAffected, totalPackages)

	// Determine risk level
	riskLevel, riskExplanation := ia.calculateRiskLevel(
		totalAffected,
		totalPackages,
		directParticipants,
	)

	// Build ripple effect data
	rippleEffect := ia.buildRippleEffect(directParticipants, indirectDependents)

	return &types.ImpactAssessment{
		DirectParticipants:        directParticipants,
		IndirectDependents:        indirectDependents,
		TotalAffected:             totalAffected,
		AffectedPercentage:        percentage,
		AffectedPercentageDisplay: displayPercentage,
		RiskLevel:                 riskLevel,
		RiskExplanation:           riskExplanation,
		RippleEffect:              rippleEffect,
	}
}

// buildReverseDependencies creates a reverse lookup map.
// For each edge (A → B), A depends on B, so B's "dependents" include A.
func (ia *ImpactAnalyzer) buildReverseDependencies() {
	ia.reverseDeps = make(map[string][]string)

	// Initialize all packages
	for pkgName := range ia.graph.Nodes {
		ia.reverseDeps[pkgName] = []string{}
	}

	// Build reverse mapping from edges
	for _, edge := range ia.graph.Edges {
		// edge.From depends on edge.To
		// So edge.From is a "dependent" of edge.To
		ia.reverseDeps[edge.To] = append(ia.reverseDeps[edge.To], edge.From)
	}
}

// getDirectParticipants extracts unique packages from cycle.
// The cycle array ends with the first element repeated, so we exclude the last element.
func (ia *ImpactAnalyzer) getDirectParticipants(cycle *types.CircularDependencyInfo) []string {
	if len(cycle.Cycle) == 0 {
		return []string{}
	}

	// Use map to dedupe (shouldn't be needed but safe)
	seen := make(map[string]bool)
	participants := []string{}

	// Cycle array ends with first element, so exclude last
	for i := 0; i < len(cycle.Cycle)-1; i++ {
		pkg := cycle.Cycle[i]
		if !seen[pkg] {
			seen[pkg] = true
			participants = append(participants, pkg)
		}
	}

	return participants
}

// queueItem represents an item in the BFS queue for finding indirect dependents.
type queueItem struct {
	pkg       string
	dependsOn string
	distance  int
	path      []string
}

// findIndirectDependents finds all packages depending on cycle participants using BFS.
func (ia *ImpactAnalyzer) findIndirectDependents(directParticipants []string) []types.IndirectDependent {
	// Mark direct participants as visited
	visited := make(map[string]bool)
	for _, pkg := range directParticipants {
		visited[pkg] = true
	}

	result := []types.IndirectDependent{}

	// BFS queue
	queue := []queueItem{}

	// Initialize queue with direct dependents of cycle participants
	for _, cyclePkg := range directParticipants {
		for _, dependent := range ia.reverseDeps[cyclePkg] {
			if !visited[dependent] {
				queue = append(queue, queueItem{
					pkg:       dependent,
					dependsOn: cyclePkg,
					distance:  1,
					path:      []string{cyclePkg, dependent},
				})
			}
		}
	}

	// BFS traversal
	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]

		if visited[item.pkg] {
			continue
		}
		visited[item.pkg] = true

		result = append(result, types.IndirectDependent{
			PackageName:    item.pkg,
			DependsOn:      item.dependsOn,
			Distance:       item.distance,
			DependencyPath: item.path,
		})

		// Add this package's dependents to queue
		for _, nextDependent := range ia.reverseDeps[item.pkg] {
			if !visited[nextDependent] {
				newPath := make([]string, len(item.path)+1)
				copy(newPath, item.path)
				newPath[len(item.path)] = nextDependent

				queue = append(queue, queueItem{
					pkg:       nextDependent,
					dependsOn: item.dependsOn, // Original cycle package
					distance:  item.distance + 1,
					path:      newPath,
				})
			}
		}
	}

	return result
}

// corePatterns are naming patterns that indicate a core/shared package.
var corePatterns = []string{"core", "common", "shared", "utils", "lib"}

// calculateRiskLevel determines risk based on metrics.
// Risk levels:
// - Critical: >50% affected OR core package involved
// - High: 25-50%
// - Medium: 10-25%
// - Low: <10%
func (ia *ImpactAnalyzer) calculateRiskLevel(
	affected int,
	total int,
	directParticipants []string,
) (types.RiskLevel, string) {
	if total == 0 {
		return types.RiskLevelLow, "Low impact: no packages in workspace"
	}

	percentage := float64(affected) / float64(total)

	// Check for core/shared package patterns
	hasCorePackage := false
	for _, pkg := range directParticipants {
		pkgLower := strings.ToLower(pkg)
		for _, pattern := range corePatterns {
			if strings.Contains(pkgLower, pattern) {
				hasCorePackage = true
				break
			}
		}
		if hasCorePackage {
			break
		}
	}

	// Critical: >50% affected OR core package involved
	if percentage > 0.50 || hasCorePackage {
		explanation := "Critical impact: "
		if hasCorePackage {
			explanation += "cycle includes core/shared package"
		} else {
			explanation += fmt.Sprintf("%.0f%% of packages affected", percentage*100)
		}
		return types.RiskLevelCritical, explanation
	}

	// High: 25-50%
	if percentage > 0.25 {
		return types.RiskLevelHigh, fmt.Sprintf("High impact: %.0f%% of packages affected", percentage*100)
	}

	// Medium: 10-25%
	if percentage > 0.10 {
		return types.RiskLevelMedium, fmt.Sprintf("Medium impact: %.0f%% of packages affected", percentage*100)
	}

	// Low: <10%
	return types.RiskLevelLow, fmt.Sprintf("Low impact: %.0f%% of packages affected", percentage*100)
}

// buildRippleEffect creates visualization data for D3.js.
func (ia *ImpactAnalyzer) buildRippleEffect(
	directParticipants []string,
	indirectDependents []types.IndirectDependent,
) *types.RippleEffect {
	// Group by distance
	layerMap := make(map[int][]string)

	// Layer 0: direct participants
	layerMap[0] = directParticipants

	// Group indirect dependents by distance
	maxDistance := 0
	for _, dep := range indirectDependents {
		layerMap[dep.Distance] = append(layerMap[dep.Distance], dep.PackageName)
		if dep.Distance > maxDistance {
			maxDistance = dep.Distance
		}
	}

	// Build ordered layers
	layers := []types.RippleLayer{}
	for distance := 0; distance <= maxDistance; distance++ {
		packages := layerMap[distance]
		if len(packages) > 0 {
			layers = append(layers, types.RippleLayer{
				Distance: distance,
				Packages: packages,
				Count:    len(packages),
			})
		}
	}

	// Handle case with no indirect dependents
	if len(layers) == 0 && len(directParticipants) > 0 {
		layers = append(layers, types.RippleLayer{
			Distance: 0,
			Packages: directParticipants,
			Count:    len(directParticipants),
		})
		maxDistance = 0
	}

	return &types.RippleEffect{
		Layers:      layers,
		TotalLayers: maxDistance + 1,
	}
}
