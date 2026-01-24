// Package types defines Go types that match TypeScript definitions in @monoguard/types.
// This file contains refactoring complexity types for Story 3.5.
package types

// ========================================
// Refactoring Complexity Types (Story 3.5)
// ========================================

// RefactoringComplexity provides detailed breakdown of fix complexity.
// Matches @monoguard/types RefactoringComplexity interface.
type RefactoringComplexity struct {
	// Score is the overall complexity (1-10)
	Score int `json:"score"`

	// EstimatedTime is human-readable time range (e.g., "15-30 minutes")
	EstimatedTime string `json:"estimatedTime"`

	// Breakdown shows individual factor contributions
	Breakdown ComplexityBreakdown `json:"breakdown"`

	// Explanation provides human-readable summary
	Explanation string `json:"explanation"`
}

// ComplexityBreakdown shows how each factor contributes to the score.
// Matches @monoguard/types ComplexityBreakdown interface.
type ComplexityBreakdown struct {
	// FilesAffected is number of source files that need changes
	FilesAffected ComplexityFactor `json:"filesAffected"`

	// ImportsToChange is number of import statements to modify
	ImportsToChange ComplexityFactor `json:"importsToChange"`

	// ChainDepth is the dependency chain depth
	ChainDepth ComplexityFactor `json:"chainDepth"`

	// PackagesInvolved is number of packages in the cycle
	PackagesInvolved ComplexityFactor `json:"packagesInvolved"`

	// ExternalDependencies indicates if external deps are involved
	ExternalDependencies ComplexityFactor `json:"externalDependencies"`
}

// ComplexityFactor represents a single factor in complexity calculation.
// Matches @monoguard/types ComplexityFactor interface.
type ComplexityFactor struct {
	// Value is the raw value for this factor
	Value int `json:"value"`

	// Weight is the factor weight (0.0-1.0)
	Weight float64 `json:"weight"`

	// Contribution is the weighted score contribution
	Contribution float64 `json:"contribution"`

	// Description explains what this factor measures
	Description string `json:"description"`
}
