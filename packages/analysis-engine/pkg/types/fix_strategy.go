// Package types defines Go types that match TypeScript definitions in @monoguard/types.
// This file contains fix strategy types for Story 3.3.
package types

// ========================================
// Fix Strategy Types (Story 3.3)
// ========================================

// FixStrategy represents a recommended approach to resolve a circular dependency.
// Matches @monoguard/types FixStrategy interface.
type FixStrategy struct {
	// Type identifies the strategy approach
	Type FixStrategyType `json:"type"`

	// Name is the human-readable strategy name
	Name string `json:"name"`

	// Description explains what this strategy does
	Description string `json:"description"`

	// Suitability is a score (1-10) indicating how well this strategy fits
	Suitability int `json:"suitability"`

	// Effort estimates the implementation difficulty
	Effort EffortLevel `json:"effort"`

	// Pros are advantages of this strategy for this specific cycle
	Pros []string `json:"pros"`

	// Cons are disadvantages of this strategy for this specific cycle
	Cons []string `json:"cons"`

	// Recommended indicates this is the top recommendation
	Recommended bool `json:"recommended"`

	// TargetPackages are the packages that would need modification
	TargetPackages []string `json:"targetPackages"`

	// NewPackageName is suggested name for extracted module (if applicable)
	NewPackageName string `json:"newPackageName,omitempty"`

	// Guide is the step-by-step fix guide (Story 3.4)
	Guide *FixGuide `json:"guide,omitempty"`
}

// FixStrategyType identifies the approach to resolve the cycle.
// Matches @monoguard/types FixStrategyType union type.
type FixStrategyType string

const (
	// FixStrategyExtractModule - Move shared code to new package
	FixStrategyExtractModule FixStrategyType = "extract-module"
	// FixStrategyDependencyInject - Invert the dependency relationship
	FixStrategyDependencyInject FixStrategyType = "dependency-injection"
	// FixStrategyBoundaryRefactor - Restructure module boundaries
	FixStrategyBoundaryRefactor FixStrategyType = "boundary-refactoring"
)

// EffortLevel estimates implementation difficulty.
// Matches @monoguard/types EffortLevel union type.
type EffortLevel string

const (
	// EffortLow - Simple changes, < 1 hour
	EffortLow EffortLevel = "low"
	// EffortMedium - Moderate changes, 1-4 hours
	EffortMedium EffortLevel = "medium"
	// EffortHigh - Significant refactoring, > 4 hours
	EffortHigh EffortLevel = "high"
)
