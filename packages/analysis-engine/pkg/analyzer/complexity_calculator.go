// Package analyzer provides dependency graph analysis for monorepo workspaces.
// This file contains complexity calculator for Story 3.5.
package analyzer

import (
	"fmt"
	"math"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// ========================================
// Complexity Calculator (Story 3.5)
// ========================================

// ComplexityCalculator computes refactoring complexity scores for circular dependencies.
type ComplexityCalculator struct {
	graph     *types.DependencyGraph
	workspace *types.WorkspaceData
}

// NewComplexityCalculator creates a new calculator.
func NewComplexityCalculator(graph *types.DependencyGraph, workspace *types.WorkspaceData) *ComplexityCalculator {
	return &ComplexityCalculator{
		graph:     graph,
		workspace: workspace,
	}
}

// Calculate computes the refactoring complexity for a cycle.
// Returns nil for invalid input (nil or empty cycle).
func (cc *ComplexityCalculator) Calculate(cycle *types.CircularDependencyInfo) *types.RefactoringComplexity {
	// Validate input
	if cycle == nil || len(cycle.Cycle) < 2 {
		return nil
	}

	// Calculate all factors
	breakdown := types.ComplexityBreakdown{
		FilesAffected:        cc.calculateFilesAffected(cycle),
		ImportsToChange:      cc.calculateImportsToChange(cycle),
		ChainDepth:           cc.calculateChainDepth(cycle),
		PackagesInvolved:     cc.calculatePackagesInvolved(cycle),
		ExternalDependencies: cc.calculateExternalDependencies(cycle),
	}

	// Sum all contributions
	totalContribution := breakdown.FilesAffected.Contribution +
		breakdown.ImportsToChange.Contribution +
		breakdown.ChainDepth.Contribution +
		breakdown.PackagesInvolved.Contribution +
		breakdown.ExternalDependencies.Contribution

	// Convert to 1-10 scale (max contribution is ~10)
	score := int(math.Round(totalContribution))
	if score < 1 {
		score = 1
	}
	if score > 10 {
		score = 10
	}

	return &types.RefactoringComplexity{
		Score:         score,
		EstimatedTime: estimateTime(score),
		Breakdown:     breakdown,
		Explanation:   generateComplexityExplanation(score, &breakdown),
	}
}

// ========================================
// Factor Calculations (Tasks 3-7)
// ========================================

// calculateFilesAffected estimates files needing changes.
// Uses ImportTraces if available, otherwise estimates based on depth.
func (cc *ComplexityCalculator) calculateFilesAffected(cycle *types.CircularDependencyInfo) types.ComplexityFactor {
	var fileCount int

	// Use ImportTraces if available (Story 3.2 provides these)
	if len(cycle.ImportTraces) > 0 {
		uniqueFiles := make(map[string]bool)
		for _, trace := range cycle.ImportTraces {
			uniqueFiles[trace.FilePath] = true
		}
		fileCount = len(uniqueFiles)
	} else {
		// Estimate: ~2 files per package in cycle (entry point + consumers)
		fileCount = cycle.Depth * 2
	}

	// Calculate contribution (weight: 0.25)
	// Scale: 1-2 files = low, 3-5 = medium, 6+ = high
	var contribution float64
	switch {
	case fileCount <= 2:
		contribution = 0.5
	case fileCount <= 5:
		contribution = 1.5
	case fileCount <= 10:
		contribution = 2.0
	default:
		contribution = 2.5
	}

	return types.ComplexityFactor{
		Value:        fileCount,
		Weight:       0.25,
		Contribution: contribution,
		Description:  fmt.Sprintf("%d source files need modification", fileCount),
	}
}

// calculateImportsToChange counts import statements to modify.
func (cc *ComplexityCalculator) calculateImportsToChange(cycle *types.CircularDependencyInfo) types.ComplexityFactor {
	var importCount int

	// Use ImportTraces if available
	if len(cycle.ImportTraces) > 0 {
		importCount = len(cycle.ImportTraces)
	} else {
		// Estimate: ~1-2 imports per edge in cycle
		importCount = cycle.Depth
	}

	// Calculate contribution (weight: 0.20)
	// Scale: 1-2 imports = low, 3-6 = medium, 7+ = high
	var contribution float64
	switch {
	case importCount <= 2:
		contribution = 0.4
	case importCount <= 6:
		contribution = 1.2
	case importCount <= 10:
		contribution = 1.6
	default:
		contribution = 2.0
	}

	return types.ComplexityFactor{
		Value:        importCount,
		Weight:       0.20,
		Contribution: contribution,
		Description:  fmt.Sprintf("%d import statements need updating", importCount),
	}
}

// calculateChainDepth evaluates dependency chain depth.
func (cc *ComplexityCalculator) calculateChainDepth(cycle *types.CircularDependencyInfo) types.ComplexityFactor {
	depth := cycle.Depth

	// Calculate contribution (weight: 0.25)
	// Scale: depth 2 = low, 3-4 = medium, 5+ = high
	var contribution float64
	switch {
	case depth <= 2:
		contribution = 0.5
	case depth <= 4:
		contribution = 1.5
	case depth <= 6:
		contribution = 2.0
	default:
		contribution = 2.5
	}

	return types.ComplexityFactor{
		Value:        depth,
		Weight:       0.25,
		Contribution: contribution,
		Description:  fmt.Sprintf("Dependency chain has %d levels", depth),
	}
}

// calculatePackagesInvolved counts packages in cycle.
func (cc *ComplexityCalculator) calculatePackagesInvolved(cycle *types.CircularDependencyInfo) types.ComplexityFactor {
	// Unique packages (cycle array ends with first, so -1)
	packageCount := len(cycle.Cycle) - 1
	if packageCount < 1 {
		packageCount = 1
	}

	// Calculate contribution (weight: 0.15)
	// Scale: 2 packages = low, 3-4 = medium, 5+ = high
	var contribution float64
	switch {
	case packageCount <= 2:
		contribution = 0.3
	case packageCount <= 4:
		contribution = 0.9
	case packageCount <= 6:
		contribution = 1.2
	default:
		contribution = 1.5
	}

	return types.ComplexityFactor{
		Value:        packageCount,
		Weight:       0.15,
		Contribution: contribution,
		Description:  fmt.Sprintf("%d packages involved in cycle", packageCount),
	}
}

// calculateExternalDependencies checks for external deps in cycle.
func (cc *ComplexityCalculator) calculateExternalDependencies(cycle *types.CircularDependencyInfo) types.ComplexityFactor {
	hasExternal := false

	// Check if any edge in root cause chain involves external deps
	if cycle.RootCause != nil {
		for _, edge := range cycle.RootCause.Chain {
			// Check if target is external (not in workspace packages)
			if cc.isExternalPackage(edge.To) {
				hasExternal = true
				break
			}
		}
	}

	// Calculate contribution (weight: 0.15)
	var contribution float64
	var value int
	if hasExternal {
		contribution = 1.5
		value = 1
	} else {
		contribution = 0.15
		value = 0
	}

	description := "No external dependencies in cycle"
	if hasExternal {
		description = "External dependencies increase complexity"
	}

	return types.ComplexityFactor{
		Value:        value,
		Weight:       0.15,
		Contribution: contribution,
		Description:  description,
	}
}

// isExternalPackage checks if a package is external (not in workspace).
func (cc *ComplexityCalculator) isExternalPackage(pkgName string) bool {
	if cc.workspace == nil || cc.workspace.Packages == nil {
		return false
	}
	_, exists := cc.workspace.Packages[pkgName]
	return !exists
}

// ========================================
// Score Aggregation and Time Estimation (Task 8)
// ========================================

// estimateTime converts score to human-readable time range.
func estimateTime(score int) string {
	switch {
	case score <= 2:
		return "5-15 minutes"
	case score <= 4:
		return "15-30 minutes"
	case score <= 6:
		return "30-60 minutes"
	case score <= 8:
		return "1-2 hours"
	default:
		return "2-4 hours"
	}
}

// generateComplexityExplanation creates human-readable complexity summary.
func generateComplexityExplanation(score int, breakdown *types.ComplexityBreakdown) string {
	var level string
	switch {
	case score <= 3:
		level = "Straightforward"
	case score <= 6:
		level = "Moderate"
	case score <= 8:
		level = "Significant"
	default:
		level = "Complex"
	}

	return fmt.Sprintf(
		"%s refactoring: %d files, %d imports, %d-level chain",
		level,
		breakdown.FilesAffected.Value,
		breakdown.ImportsToChange.Value,
		breakdown.ChainDepth.Value,
	)
}

// ========================================
// Strategy-Specific Complexity (AC6)
// ========================================

// CalculateForStrategy computes complexity adjusted for a specific fix strategy.
// Different strategies have different inherent complexity multipliers:
// - extract-module: 1.0x (baseline - cleanest approach)
// - dependency-injection: 1.1x (requires interface design)
// - boundary-refactoring: 1.2x (most invasive, architectural changes)
func (cc *ComplexityCalculator) CalculateForStrategy(cycle *types.CircularDependencyInfo, strategyType types.FixStrategyType) *types.RefactoringComplexity {
	base := cc.Calculate(cycle)
	if base == nil {
		return nil
	}

	// Apply strategy-specific multiplier
	multiplier := getStrategyMultiplier(strategyType)

	// Adjust score with multiplier (keep within 1-10 bounds)
	adjustedScore := int(math.Round(float64(base.Score) * multiplier))
	if adjustedScore < 1 {
		adjustedScore = 1
	}
	if adjustedScore > 10 {
		adjustedScore = 10
	}

	// Return new complexity with adjusted score
	return &types.RefactoringComplexity{
		Score:         adjustedScore,
		EstimatedTime: estimateTime(adjustedScore),
		Breakdown:     base.Breakdown,
		Explanation:   generateStrategyExplanation(adjustedScore, &base.Breakdown, strategyType),
	}
}

// getStrategyMultiplier returns the complexity multiplier for a strategy type.
func getStrategyMultiplier(strategyType types.FixStrategyType) float64 {
	switch strategyType {
	case types.FixStrategyExtractModule:
		return 1.0 // Baseline - cleanest refactoring approach
	case types.FixStrategyDependencyInject:
		return 1.1 // Requires interface design and abstraction
	case types.FixStrategyBoundaryRefactor:
		return 1.2 // Most invasive, affects module boundaries
	default:
		return 1.0
	}
}

// generateStrategyExplanation creates explanation including strategy context.
func generateStrategyExplanation(score int, breakdown *types.ComplexityBreakdown, strategyType types.FixStrategyType) string {
	var level string
	switch {
	case score <= 3:
		level = "Straightforward"
	case score <= 6:
		level = "Moderate"
	case score <= 8:
		level = "Significant"
	default:
		level = "Complex"
	}

	strategyName := string(strategyType)
	return fmt.Sprintf(
		"%s refactoring via %s: %d files, %d imports, %d-level chain",
		level,
		strategyName,
		breakdown.FilesAffected.Value,
		breakdown.ImportsToChange.Value,
		breakdown.ChainDepth.Value,
	)
}
