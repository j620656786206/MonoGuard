// Package types defines Go types that match TypeScript definitions in @monoguard/types.
// This file contains FixSummary types for Story 3.8.
package types

// ========================================
// Fix Summary Types (Story 3.8)
// ========================================

// FixSummary provides aggregated statistics about fix recommendations.
// Added to AnalysisResult for high-level overview.
// Matches @monoguard/types FixSummary interface.
type FixSummary struct {
	// TotalCircularDependencies is the count of all detected cycles
	TotalCircularDependencies int `json:"totalCircularDependencies"`

	// TotalEstimatedFixTime is the sum of all fix times (human-readable)
	TotalEstimatedFixTime string `json:"totalEstimatedFixTime"`

	// QuickWinsCount is the number of low-complexity (1-3) fixes
	QuickWinsCount int `json:"quickWinsCount"`

	// CriticalCyclesCount is the number of critical impact cycles
	CriticalCyclesCount int `json:"criticalCyclesCount"`

	// HighPriorityCycles lists the top 3 cycles to fix first (by priority score)
	HighPriorityCycles []PriorityCycleSummary `json:"highPriorityCycles"`
}

// PriorityCycleSummary provides a brief overview of a prioritized cycle.
// Matches @monoguard/types PriorityCycleSummary interface.
type PriorityCycleSummary struct {
	// CycleID is a unique identifier for the cycle (e.g., first two packages)
	CycleID string `json:"cycleId"`

	// PackagesInvolved lists the packages in the cycle
	PackagesInvolved []string `json:"packagesInvolved"`

	// PriorityScore is the calculated priority (impact × ease)
	PriorityScore float64 `json:"priorityScore"`

	// RecommendedFix is the quick fix strategy type
	RecommendedFix FixStrategyType `json:"recommendedFix"`

	// EstimatedTime is the fix time
	EstimatedTime string `json:"estimatedTime"`
}
