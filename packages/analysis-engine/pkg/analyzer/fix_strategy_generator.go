// Package analyzer provides dependency graph analysis for monorepo workspaces.
// This file contains fix strategy generator for Story 3.3.
package analyzer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// ========================================
// Fix Strategy Generator (Story 3.3)
// ========================================

// corePackagePatterns are keywords that indicate a core/shared package.
var corePackagePatterns = []string{"core", "common", "shared", "utils", "lib", "base"}

// FixStrategyGenerator creates fix recommendations for circular dependencies.
type FixStrategyGenerator struct {
	graph     *types.DependencyGraph
	workspace *types.WorkspaceData
}

// NewFixStrategyGenerator creates a new generator.
func NewFixStrategyGenerator(graph *types.DependencyGraph, workspace *types.WorkspaceData) *FixStrategyGenerator {
	return &FixStrategyGenerator{
		graph:     graph,
		workspace: workspace,
	}
}

// Generate creates fix strategies for a circular dependency.
// Requires RootCause to be populated (from Story 3.1) for optimal recommendations.
// Returns up to 3 strategies, sorted by suitability (highest first).
func (fsg *FixStrategyGenerator) Generate(cycle *types.CircularDependencyInfo) []types.FixStrategy {
	// Return empty slice for invalid input (never nil per project rules)
	if cycle == nil || len(cycle.Cycle) < 2 {
		return []types.FixStrategy{}
	}

	strategies := []types.FixStrategy{}

	// Generate all three strategy types
	if extractStrategy := fsg.generateExtractModule(cycle); extractStrategy != nil {
		strategies = append(strategies, *extractStrategy)
	}

	if diStrategy := fsg.generateDependencyInjection(cycle); diStrategy != nil {
		strategies = append(strategies, *diStrategy)
	}

	if boundaryStrategy := fsg.generateBoundaryRefactoring(cycle); boundaryStrategy != nil {
		strategies = append(strategies, *boundaryStrategy)
	}

	// Rank strategies by suitability
	strategies = rankStrategies(strategies)

	return strategies
}

// ========================================
// Extract Shared Module Strategy (Task 3)
// ========================================

// generateExtractModule creates the Extract Shared Module strategy.
func (fsg *FixStrategyGenerator) generateExtractModule(cycle *types.CircularDependencyInfo) *types.FixStrategy {
	// Calculate suitability based on cycle characteristics
	suitability := fsg.calculateExtractModuleSuitability(cycle)

	// Estimate effort based on cycle depth
	effort := fsg.calculateExtractModuleEffort(cycle)

	// Generate contextual pros/cons
	pros, cons := fsg.generateExtractModuleProsCons(cycle)

	// Get target packages (all packages in the cycle, excluding the closing node)
	targetPackages := fsg.getTargetPackages(cycle)

	// Suggest new package name
	newPkgName := fsg.suggestNewPackageName(cycle)

	return &types.FixStrategy{
		Type:           types.FixStrategyExtractModule,
		Name:           "Extract Shared Module",
		Description:    "Create a new shared package to hold common dependencies, breaking the cycle.",
		Suitability:    suitability,
		Effort:         effort,
		Pros:           pros,
		Cons:           cons,
		Recommended:    false, // Will be set by rankStrategies
		TargetPackages: targetPackages,
		NewPackageName: newPkgName,
	}
}

// calculateExtractModuleSuitability calculates suitability score for Extract Module.
// Higher for longer cycles, lower when core package already exists.
func (fsg *FixStrategyGenerator) calculateExtractModuleSuitability(cycle *types.CircularDependencyInfo) int {
	// Base score from cycle depth factor (40% weight)
	depthScore := cycleDepthFactor(cycle.Depth, types.FixStrategyExtractModule)

	// Dependency type factor (30% weight)
	depTypeScore := fsg.dependencyTypeFactor(cycle, types.FixStrategyExtractModule)

	// Package naming pattern factor (30% weight)
	namingScore := namingPatternFactor(cycle.Cycle, types.FixStrategyExtractModule)

	// Weighted average
	suitability := (depthScore*40 + depTypeScore*30 + namingScore*30) / 100

	// Ensure bounds 1-10
	if suitability < 1 {
		suitability = 1
	}
	if suitability > 10 {
		suitability = 10
	}

	return suitability
}

// calculateExtractModuleEffort estimates effort for Extract Module strategy.
func (fsg *FixStrategyGenerator) calculateExtractModuleEffort(cycle *types.CircularDependencyInfo) types.EffortLevel {
	depth := cycle.Depth

	// Effort scales with cycle depth
	if depth <= 2 {
		return types.EffortMedium
	}
	if depth <= 4 {
		return types.EffortMedium
	}
	return types.EffortHigh
}

// generateExtractModuleProsCons generates contextual pros/cons for Extract Module.
func (fsg *FixStrategyGenerator) generateExtractModuleProsCons(cycle *types.CircularDependencyInfo) ([]string, []string) {
	// Get package names for contextual messages
	packageList := formatPackageList(cycle.Cycle[:len(cycle.Cycle)-1])

	pros := []string{
		"Creates clear separation of concerns",
		fmt.Sprintf("Isolates shared code between %s", packageList),
	}

	cons := []string{
		"Introduces a new package to maintain",
	}

	// Add contextual pros/cons based on cycle characteristics
	if cycle.Depth >= 4 {
		pros = append(pros, "Effectively breaks complex multi-package cycle")
	}
	if containsPattern(cycle.Cycle, corePackagePatterns) {
		cons = append(cons, "May require updating existing core package consumers")
	}

	return pros, cons
}

// ========================================
// Dependency Injection Strategy (Task 4)
// ========================================

// generateDependencyInjection creates the Dependency Injection strategy.
func (fsg *FixStrategyGenerator) generateDependencyInjection(cycle *types.CircularDependencyInfo) *types.FixStrategy {
	// Calculate suitability based on cycle characteristics
	suitability := fsg.calculateDISuitability(cycle)

	// Estimate effort
	effort := fsg.calculateDIEffort(cycle)

	// Generate contextual pros/cons
	pros, cons := fsg.generateDIProsCons(cycle)

	// Get target packages (primarily the packages around the critical edge)
	targetPackages := fsg.getDITargetPackages(cycle)

	return &types.FixStrategy{
		Type:           types.FixStrategyDependencyInject,
		Name:           "Dependency Injection",
		Description:    "Invert the problematic dependency by introducing an interface or callback pattern.",
		Suitability:    suitability,
		Effort:         effort,
		Pros:           pros,
		Cons:           cons,
		Recommended:    false, // Will be set by rankStrategies
		TargetPackages: targetPackages,
	}
}

// calculateDISuitability calculates suitability score for Dependency Injection.
// Higher for direct cycles, lower for deeply nested cycles.
func (fsg *FixStrategyGenerator) calculateDISuitability(cycle *types.CircularDependencyInfo) int {
	// Base score from cycle depth factor (40% weight)
	depthScore := cycleDepthFactor(cycle.Depth, types.FixStrategyDependencyInject)

	// Dependency type factor (30% weight)
	depTypeScore := fsg.dependencyTypeFactor(cycle, types.FixStrategyDependencyInject)

	// Package naming pattern factor (30% weight)
	namingScore := namingPatternFactor(cycle.Cycle, types.FixStrategyDependencyInject)

	// Bonus for clear critical edge from RootCause analysis
	criticalEdgeBonus := 0
	if cycle.RootCause != nil && cycle.RootCause.CriticalEdge != nil {
		criticalEdgeBonus = 1
	}

	// Weighted average
	suitability := (depthScore*40 + depTypeScore*30 + namingScore*30) / 100
	suitability += criticalEdgeBonus

	// Ensure bounds 1-10
	if suitability < 1 {
		suitability = 1
	}
	if suitability > 10 {
		suitability = 10
	}

	return suitability
}

// calculateDIEffort estimates effort for Dependency Injection strategy.
func (fsg *FixStrategyGenerator) calculateDIEffort(cycle *types.CircularDependencyInfo) types.EffortLevel {
	depth := cycle.Depth

	// DI is typically lower effort for direct cycles
	if depth <= 2 {
		return types.EffortLow
	}
	if depth <= 3 {
		return types.EffortMedium
	}
	return types.EffortHigh
}

// generateDIProsCons generates contextual pros/cons for Dependency Injection.
func (fsg *FixStrategyGenerator) generateDIProsCons(cycle *types.CircularDependencyInfo) ([]string, []string) {
	pros := []string{
		"Minimal code changes required",
		"Preserves existing package structure",
	}

	cons := []string{
		"Adds indirection to the codebase",
	}

	// Add contextual pros/cons based on critical edge
	if cycle.RootCause != nil && cycle.RootCause.CriticalEdge != nil {
		edge := cycle.RootCause.CriticalEdge
		pros = append(pros, fmt.Sprintf("Clear injection point: %s → %s", edge.From, edge.To))
	}

	if cycle.Depth > 3 {
		cons = append(cons, "May require multiple interfaces for complex cycles")
	}

	return pros, cons
}

// getDITargetPackages returns the packages to modify for DI strategy.
func (fsg *FixStrategyGenerator) getDITargetPackages(cycle *types.CircularDependencyInfo) []string {
	// If we have root cause analysis with critical edge, use those packages
	if cycle.RootCause != nil && cycle.RootCause.CriticalEdge != nil {
		edge := cycle.RootCause.CriticalEdge
		return []string{edge.From, edge.To}
	}

	// Otherwise, use first two packages in cycle
	if len(cycle.Cycle) >= 2 {
		return []string{cycle.Cycle[0], cycle.Cycle[1]}
	}

	return fsg.getTargetPackages(cycle)
}

// ========================================
// Module Boundary Refactoring Strategy (Task 5)
// ========================================

// generateBoundaryRefactoring creates the Module Boundary Refactoring strategy.
func (fsg *FixStrategyGenerator) generateBoundaryRefactoring(cycle *types.CircularDependencyInfo) *types.FixStrategy {
	// Calculate suitability based on cycle characteristics
	suitability := fsg.calculateBoundaryRefactorSuitability(cycle)

	// Estimate effort
	effort := fsg.calculateBoundaryRefactorEffort(cycle)

	// Generate contextual pros/cons
	pros, cons := fsg.generateBoundaryRefactorProsCons(cycle)

	// Get target packages
	targetPackages := fsg.getBoundaryRefactorTargetPackages(cycle)

	return &types.FixStrategy{
		Type:           types.FixStrategyBoundaryRefactor,
		Name:           "Module Boundary Refactoring",
		Description:    "Restructure package boundaries to eliminate the cyclic relationship.",
		Suitability:    suitability,
		Effort:         effort,
		Pros:           pros,
		Cons:           cons,
		Recommended:    false, // Will be set by rankStrategies
		TargetPackages: targetPackages,
	}
}

// calculateBoundaryRefactorSuitability calculates suitability for Boundary Refactoring.
// Higher when packages have overlapping responsibilities or involve core packages.
func (fsg *FixStrategyGenerator) calculateBoundaryRefactorSuitability(cycle *types.CircularDependencyInfo) int {
	// Base score from cycle depth factor (30% weight)
	depthScore := cycleDepthFactor(cycle.Depth, types.FixStrategyBoundaryRefactor)

	// Dependency type factor (20% weight)
	depTypeScore := fsg.dependencyTypeFactor(cycle, types.FixStrategyBoundaryRefactor)

	// Package naming pattern factor (50% weight) - more important for boundary refactoring
	namingScore := namingPatternFactor(cycle.Cycle, types.FixStrategyBoundaryRefactor)

	// Weighted average with emphasis on naming patterns
	suitability := (depthScore*30 + depTypeScore*20 + namingScore*50) / 100

	// Ensure bounds 1-10
	if suitability < 1 {
		suitability = 1
	}
	if suitability > 10 {
		suitability = 10
	}

	return suitability
}

// calculateBoundaryRefactorEffort estimates effort for Boundary Refactoring.
func (fsg *FixStrategyGenerator) calculateBoundaryRefactorEffort(cycle *types.CircularDependencyInfo) types.EffortLevel {
	depth := cycle.Depth

	// Boundary refactoring is usually higher effort
	if depth <= 2 {
		return types.EffortMedium
	}
	return types.EffortHigh
}

// generateBoundaryRefactorProsCons generates contextual pros/cons for Boundary Refactoring.
func (fsg *FixStrategyGenerator) generateBoundaryRefactorProsCons(cycle *types.CircularDependencyInfo) ([]string, []string) {
	pros := []string{
		"Addresses root architectural issue",
		"Results in cleaner long-term design",
	}

	cons := []string{
		"Requires significant refactoring effort",
		"May affect external package consumers",
	}

	// Add contextual pros/cons
	if containsPattern(cycle.Cycle, corePackagePatterns) {
		pros = append(pros, "Opportunity to properly define core package boundaries")
	}

	return pros, cons
}

// getBoundaryRefactorTargetPackages returns all packages in the cycle.
func (fsg *FixStrategyGenerator) getBoundaryRefactorTargetPackages(cycle *types.CircularDependencyInfo) []string {
	return fsg.getTargetPackages(cycle)
}

// ========================================
// Suitability Scoring Algorithm (Task 6)
// ========================================

// cycleDepthFactor returns a 0-10 score based on cycle depth.
// Different strategies prefer different cycle depths.
func cycleDepthFactor(depth int, strategyType types.FixStrategyType) int {
	switch strategyType {
	case types.FixStrategyExtractModule:
		// Better for longer cycles
		if depth >= 4 {
			return 10
		}
		if depth == 3 {
			return 7
		}
		return 4 // Direct cycle (depth 2)

	case types.FixStrategyDependencyInject:
		// Better for shorter cycles
		if depth == 2 {
			return 10
		}
		if depth == 3 {
			return 6
		}
		return 3 // Long cycles harder to DI

	case types.FixStrategyBoundaryRefactor:
		// Always viable, slightly better for medium
		if depth == 3 {
			return 8
		}
		return 6
	}
	return 5
}

// dependencyTypeFactor returns a 0-10 score based on dependency types in the cycle.
func (fsg *FixStrategyGenerator) dependencyTypeFactor(cycle *types.CircularDependencyInfo, strategyType types.FixStrategyType) int {
	// Check if any edge in the cycle is a dev dependency
	hasDevDep := fsg.cycleHasDevDependency(cycle)

	switch strategyType {
	case types.FixStrategyExtractModule:
		// Dev deps are easier to extract
		if hasDevDep {
			return 8
		}
		return 6

	case types.FixStrategyDependencyInject:
		// Works well for production deps
		if !hasDevDep {
			return 8
		}
		return 5

	case types.FixStrategyBoundaryRefactor:
		return 6 // Neutral
	}
	return 5
}

// cycleHasDevDependency checks if the cycle contains any development dependencies.
func (fsg *FixStrategyGenerator) cycleHasDevDependency(cycle *types.CircularDependencyInfo) bool {
	// Check root cause chain if available
	if cycle.RootCause != nil && len(cycle.RootCause.Chain) > 0 {
		for _, edge := range cycle.RootCause.Chain {
			if edge.Type == types.DependencyTypeDevelopment {
				return true
			}
		}
	}

	// Check graph edges
	for i := 0; i < len(cycle.Cycle)-1; i++ {
		from := cycle.Cycle[i]
		to := cycle.Cycle[i+1]

		for _, edge := range fsg.graph.Edges {
			if edge.From == from && edge.To == to {
				if edge.Type == types.DependencyTypeDevelopment {
					return true
				}
			}
		}
	}

	return false
}

// namingPatternFactor returns a 0-10 score based on package naming patterns.
func namingPatternFactor(cycle []string, strategyType types.FixStrategyType) int {
	hasCorePackage := containsPattern(cycle, corePackagePatterns)

	switch strategyType {
	case types.FixStrategyExtractModule:
		// If already has core, extraction less valuable
		if hasCorePackage {
			return 4
		}
		return 8

	case types.FixStrategyBoundaryRefactor:
		// Core package cycles often need boundary rework
		if hasCorePackage {
			return 9
		}
		return 5

	default:
		return 6
	}
}

// ========================================
// Strategy Ranking (Task 9)
// ========================================

// rankStrategies sorts strategies by suitability (descending), then effort (ascending).
// Marks the top strategy as recommended.
func rankStrategies(strategies []types.FixStrategy) []types.FixStrategy {
	if len(strategies) == 0 {
		return strategies
	}

	// Sort by suitability (descending), then effort (ascending)
	sort.Slice(strategies, func(i, j int) bool {
		if strategies[i].Suitability != strategies[j].Suitability {
			return strategies[i].Suitability > strategies[j].Suitability
		}
		// Lower effort wins ties
		effortOrder := map[types.EffortLevel]int{
			types.EffortLow:    0,
			types.EffortMedium: 1,
			types.EffortHigh:   2,
		}
		return effortOrder[strategies[i].Effort] < effortOrder[strategies[j].Effort]
	})

	// Mark top strategy as recommended
	strategies[0].Recommended = true

	return strategies
}

// ========================================
// Helper Functions
// ========================================

// getTargetPackages returns all unique packages in the cycle (excluding closing node).
func (fsg *FixStrategyGenerator) getTargetPackages(cycle *types.CircularDependencyInfo) []string {
	if len(cycle.Cycle) < 2 {
		return []string{}
	}

	// Exclude closing node (last element is same as first)
	return cycle.Cycle[:len(cycle.Cycle)-1]
}

// suggestNewPackageName suggests a name for a new shared package.
func (fsg *FixStrategyGenerator) suggestNewPackageName(cycle *types.CircularDependencyInfo) string {
	if len(cycle.Cycle) < 2 {
		return ""
	}

	// Extract scope from first package name (e.g., "@mono/ui" -> "@mono")
	firstPkg := cycle.Cycle[0]
	scope := ""
	if strings.HasPrefix(firstPkg, "@") {
		parts := strings.SplitN(firstPkg, "/", 2)
		if len(parts) >= 1 {
			scope = parts[0]
		}
	}

	if scope != "" {
		return fmt.Sprintf("%s/shared", scope)
	}

	return "shared"
}

// formatPackageList formats package names as a comma-separated list.
func formatPackageList(packages []string) string {
	if len(packages) == 0 {
		return ""
	}
	if len(packages) == 1 {
		return packages[0]
	}
	if len(packages) == 2 {
		return packages[0] + " and " + packages[1]
	}

	// Multiple packages: "A, B, and C"
	last := packages[len(packages)-1]
	rest := packages[:len(packages)-1]
	return strings.Join(rest, ", ") + ", and " + last
}

// containsPattern checks if any package name contains any of the given patterns.
func containsPattern(packages []string, patterns []string) bool {
	for _, pkg := range packages {
		pkgLower := strings.ToLower(pkg)
		for _, pattern := range patterns {
			if strings.Contains(pkgLower, pattern) {
				return true
			}
		}
	}
	return false
}
