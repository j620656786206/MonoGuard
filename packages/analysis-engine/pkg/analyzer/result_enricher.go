// Package analyzer provides dependency graph analysis for monorepo workspaces.
// This file contains the result enricher for Story 3.8.
package analyzer

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// ResultEnricher integrates all fix suggestions into analysis results.
// Story 3.8: Sorts strategies, calculates priority scores, and generates FixSummary.
type ResultEnricher struct {
	graph     *types.DependencyGraph
	workspace *types.WorkspaceData
}

// NewResultEnricher creates a new enricher instance.
func NewResultEnricher(graph *types.DependencyGraph, workspace *types.WorkspaceData) *ResultEnricher {
	return &ResultEnricher{
		graph:     graph,
		workspace: workspace,
	}
}

// Enrich adds all fix-related data to the analysis result.
// This is the final step in the analysis pipeline that:
// 1. Sorts strategies by suitability (best first) - AC2
// 2. Creates QuickFix from best strategy - AC1, AC5
// 3. Calculates priority scores - AC3
// 4. Sorts circular dependencies by priority - AC3
// 5. Generates aggregated FixSummary - AC4
func (re *ResultEnricher) Enrich(result *types.AnalysisResult) *types.AnalysisResult {
	if result == nil || len(result.CircularDependencies) == 0 {
		return result
	}

	// Step 1: Enrich each circular dependency
	for i := range result.CircularDependencies {
		re.enrichCircularDependency(result.CircularDependencies[i])
	}

	// Step 2: Sort circular dependencies by priority (AC3)
	sortCircularDependencies(result.CircularDependencies)

	// Step 3: Generate aggregated fix summary (AC4)
	result.FixSummary = re.generateFixSummary(result.CircularDependencies)

	return result
}

// enrichCircularDependency adds quick fix and priority score to a single cycle.
func (re *ResultEnricher) enrichCircularDependency(cycle *types.CircularDependencyInfo) {
	// Sort strategies by suitability (highest first) - AC2
	if len(cycle.FixStrategies) > 0 {
		sortStrategies(cycle.FixStrategies)
		// Create QuickFix from the best (first) strategy - AC1, AC5
		cycle.QuickFix = createQuickFix(cycle.FixStrategies)
	}

	// Calculate priority score - AC3
	cycle.PriorityScore = calculatePriorityScore(cycle)
}

// sortStrategies sorts fix strategies by suitability (highest first).
// Uses stable sort for deterministic ordering (AC3).
func sortStrategies(strategies []types.FixStrategy) {
	sort.SliceStable(strategies, func(i, j int) bool {
		return strategies[i].Suitability > strategies[j].Suitability
	})
}

// calculatePriorityScore computes impact × ease for sorting.
// Higher priority = higher impact + lower complexity (quick wins first).
// Formula: Priority = ImpactFactor × EaseFactor
// - ImpactFactor: Based on RiskLevel (critical=10, high=7.5, medium=5, low=2.5)
// - EaseFactor: Based on Complexity (11 - complexity score)
func calculatePriorityScore(cycle *types.CircularDependencyInfo) float64 {
	// Impact factor: use ImpactAssessment if available
	var impactFactor float64 = 5.0 // Default medium impact
	if cycle.ImpactAssessment != nil {
		switch cycle.ImpactAssessment.RiskLevel {
		case types.RiskLevelCritical:
			impactFactor = 10.0
		case types.RiskLevelHigh:
			impactFactor = 7.5
		case types.RiskLevelMedium:
			impactFactor = 5.0
		case types.RiskLevelLow:
			impactFactor = 2.5
		}
	}

	// Ease factor: inverse of complexity (11 - complexity)
	// Higher ease = lower complexity = easier to fix
	var easeFactor float64 = 6.0 // Default medium ease (complexity 5)
	if cycle.RefactoringComplexity != nil {
		easeFactor = float64(11 - cycle.RefactoringComplexity.Score)
		// Ensure ease factor is at least 1
		if easeFactor < 1.0 {
			easeFactor = 1.0
		}
	}

	// Priority = Impact × Ease (quick wins = high impact, low complexity)
	return impactFactor * easeFactor
}

// sortCircularDependencies sorts cycles by priority score (highest first).
// Uses stable sort for deterministic ordering (AC3).
func sortCircularDependencies(cycles []*types.CircularDependencyInfo) {
	sort.SliceStable(cycles, func(i, j int) bool {
		return cycles[i].PriorityScore > cycles[j].PriorityScore
	})
}

// createQuickFix extracts the best strategy as QuickFixSummary.
// Assumes strategies are already sorted by suitability.
func createQuickFix(strategies []types.FixStrategy) *types.QuickFixSummary {
	if len(strategies) == 0 {
		return nil
	}

	// First strategy is best (already sorted by suitability)
	best := strategies[0]

	// Generate one-line summary based on strategy type
	summary := generateQuickFixSummary(&best)

	// Get estimated time from guide or use default
	estimatedTime := "15-30 minutes" // Default
	if best.Guide != nil && best.Guide.EstimatedTime != "" {
		estimatedTime = best.Guide.EstimatedTime
	} else if best.Complexity != nil && best.Complexity.EstimatedTime != "" {
		estimatedTime = best.Complexity.EstimatedTime
	}

	return &types.QuickFixSummary{
		StrategyType:  best.Type,
		StrategyName:  best.Name,
		Summary:       summary,
		Suitability:   best.Suitability,
		Effort:        best.Effort,
		EstimatedTime: estimatedTime,
		Guide:         best.Guide, // Embed full guide for one-click access (AC5)
		StrategyIndex: 0,
	}
}

// generateQuickFixSummary creates a one-line description of what the fix accomplishes.
func generateQuickFixSummary(strategy *types.FixStrategy) string {
	switch strategy.Type {
	case types.FixStrategyExtractModule:
		if strategy.NewPackageName != "" {
			return fmt.Sprintf("Create new shared package '%s' to break the cycle", strategy.NewPackageName)
		}
		return "Extract shared code into a new package to break the cycle"
	case types.FixStrategyDependencyInject:
		return "Invert dependency using dependency injection pattern"
	case types.FixStrategyBoundaryRefactor:
		return "Restructure package boundaries to eliminate overlap"
	default:
		// Fall back to strategy description if available
		if strategy.Description != "" {
			// Truncate if too long
			if len(strategy.Description) > 80 {
				return strategy.Description[:77] + "..."
			}
			return strategy.Description
		}
		return "Apply recommended fix to break the circular dependency"
	}
}

// generateFixSummary creates aggregated fix statistics.
// Includes: total cycles, total time, quick wins count, critical count, top 3 priorities.
func (re *ResultEnricher) generateFixSummary(cycles []*types.CircularDependencyInfo) *types.FixSummary {
	if len(cycles) == 0 {
		return nil
	}

	totalMinutes := 0
	quickWinsCount := 0
	criticalCount := 0
	highPriorityCycles := []types.PriorityCycleSummary{}

	for i, cycle := range cycles {
		// Count quick wins (complexity <= 3)
		if cycle.RefactoringComplexity != nil && cycle.RefactoringComplexity.Score <= 3 {
			quickWinsCount++
		}

		// Count critical cycles
		if cycle.ImpactAssessment != nil && cycle.ImpactAssessment.RiskLevel == types.RiskLevelCritical {
			criticalCount++
		}

		// Sum estimated times (parse from string)
		if cycle.QuickFix != nil {
			totalMinutes += parseEstimatedMinutes(cycle.QuickFix.EstimatedTime)
		} else if cycle.RefactoringComplexity != nil {
			totalMinutes += parseEstimatedMinutes(cycle.RefactoringComplexity.EstimatedTime)
		} else {
			totalMinutes += 30 // Default 30 minutes per cycle
		}

		// Collect top 3 high priority cycles
		if i < 3 {
			cycleID := generateCycleID(cycle.Cycle)
			var recommendedFix types.FixStrategyType
			var estTime string
			if cycle.QuickFix != nil {
				recommendedFix = cycle.QuickFix.StrategyType
				estTime = cycle.QuickFix.EstimatedTime
			}

			highPriorityCycles = append(highPriorityCycles, types.PriorityCycleSummary{
				CycleID:          cycleID,
				PackagesInvolved: getUniquePackages(cycle.Cycle),
				PriorityScore:    cycle.PriorityScore,
				RecommendedFix:   recommendedFix,
				EstimatedTime:    estTime,
			})
		}
	}

	return &types.FixSummary{
		TotalCircularDependencies: len(cycles),
		TotalEstimatedFixTime:     formatTotalTime(totalMinutes),
		QuickWinsCount:            quickWinsCount,
		CriticalCyclesCount:       criticalCount,
		HighPriorityCycles:        highPriorityCycles,
	}
}

// parseEstimatedMinutes extracts minutes from time strings like "15-30 minutes", "1-2 hours".
// Returns the midpoint of the range in minutes.
func parseEstimatedMinutes(timeStr string) int {
	if timeStr == "" {
		return 30 // Default
	}

	timeStr = strings.ToLower(timeStr)

	// Try to extract numeric range
	re := regexp.MustCompile(`(\d+)[-–](\d+)\s*(minute|hour)`)
	matches := re.FindStringSubmatch(timeStr)
	if len(matches) >= 4 {
		low, _ := strconv.Atoi(matches[1])
		high, _ := strconv.Atoi(matches[2])

		// Convert to minutes first, then calculate midpoint
		if strings.Contains(matches[3], "hour") {
			low *= 60
			high *= 60
		}
		midpoint := (low + high) / 2
		return midpoint
	}

	// Try single number
	singleRe := regexp.MustCompile(`(\d+)\s*(minute|hour)`)
	singleMatches := singleRe.FindStringSubmatch(timeStr)
	if len(singleMatches) >= 3 {
		value, _ := strconv.Atoi(singleMatches[1])
		if strings.Contains(singleMatches[2], "hour") {
			value *= 60
		}
		return value
	}

	// Common patterns fallback
	if strings.Contains(timeStr, "hour") {
		return 90 // Default for hour-range estimates
	}
	if strings.Contains(timeStr, "5-15") {
		return 10
	}
	if strings.Contains(timeStr, "15-30") {
		return 22
	}
	if strings.Contains(timeStr, "30-60") {
		return 45
	}

	return 30 // Default
}

// formatTotalTime converts minutes to human-readable format.
func formatTotalTime(minutes int) string {
	if minutes < 60 {
		return fmt.Sprintf("%d minutes", minutes)
	}
	hours := minutes / 60
	remainingMinutes := minutes % 60
	if remainingMinutes == 0 {
		if hours == 1 {
			return "1 hour"
		}
		return fmt.Sprintf("%d hours", hours)
	}
	if hours == 1 {
		return fmt.Sprintf("1 hour %d minutes", remainingMinutes)
	}
	return fmt.Sprintf("%d hours %d minutes", hours, remainingMinutes)
}

// generateCycleID creates a unique identifier for the cycle (e.g., "core→ui").
func generateCycleID(cycle []string) string {
	if len(cycle) < 2 {
		return "unknown"
	}
	return fmt.Sprintf("%s→%s", extractShortName(cycle[0]), extractShortName(cycle[1]))
}

// Note: extractShortName is defined in before_after_generator.go and is reused here.

// getUniquePackages returns unique packages in the cycle (excluding the closing node).
func getUniquePackages(cycle []string) []string {
	if len(cycle) == 0 {
		return []string{}
	}
	// Exclude last element (duplicate of first in closed cycle)
	if len(cycle) > 1 && cycle[0] == cycle[len(cycle)-1] {
		return cycle[:len(cycle)-1]
	}
	return cycle
}
