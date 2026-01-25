// Package types defines Go types that match TypeScript definitions in @monoguard/types.
// This file contains QuickFixSummary types for Story 3.8.
package types

// ========================================
// Quick Fix Summary Types (Story 3.8)
// ========================================

// QuickFixSummary provides quick access to the best fix recommendation.
// This is a convenience wrapper around the best FixStrategy.
// Matches @monoguard/types QuickFixSummary interface.
type QuickFixSummary struct {
	// StrategyType is the type of the recommended fix
	StrategyType FixStrategyType `json:"strategyType"`

	// StrategyName is the human-readable name
	StrategyName string `json:"strategyName"`

	// Summary is a one-line description of what the fix accomplishes
	Summary string `json:"summary"`

	// Suitability is the strategy's suitability score (1-10)
	Suitability int `json:"suitability"`

	// Effort is the estimated effort level
	Effort EffortLevel `json:"effort"`

	// EstimatedTime is the time to implement (e.g., "15-30 minutes")
	EstimatedTime string `json:"estimatedTime"`

	// Guide is the full step-by-step guide (embedded for one-click access)
	Guide *FixGuide `json:"guide,omitempty"`

	// StrategyIndex is the index into fixStrategies[] for full details
	StrategyIndex int `json:"strategyIndex"`
}
