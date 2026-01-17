// Package types defines Go types that match TypeScript definitions in @monoguard/types.
// This file contains health score types for Story 2.5.
package types

// ========================================
// Health Score Types (Story 2.5)
// ========================================

// HealthScoreResult represents the complete health score with breakdown.
// Matches @monoguard/types HealthScore interface.
type HealthScoreResult struct {
	Overall   int             `json:"overall"`   // 0-100
	Rating    HealthRating    `json:"rating"`    // excellent, good, fair, poor, critical
	Breakdown *ScoreBreakdown `json:"breakdown"` // Individual factor scores
	Factors   []*HealthFactor `json:"factors"`   // Detailed factor information
	UpdatedAt string          `json:"updatedAt"` // ISO 8601 format
}

// ScoreBreakdown shows individual factor scores (0-100 each).
type ScoreBreakdown struct {
	CircularScore int `json:"circularScore"` // Score from circular dependency analysis
	ConflictScore int `json:"conflictScore"` // Score from version conflict analysis
	DepthScore    int `json:"depthScore"`    // Score from dependency depth analysis
	CouplingScore int `json:"couplingScore"` // Score from package coupling analysis
}

// HealthFactor represents a single factor in the health calculation.
type HealthFactor struct {
	Name            string   `json:"name"`            // Factor name
	Score           int      `json:"score"`           // Raw score 0-100
	Weight          float64  `json:"weight"`          // Weight 0.0-1.0
	WeightedScore   int      `json:"weightedScore"`   // Contribution to overall (score * weight)
	Description     string   `json:"description"`     // Human-readable description
	Recommendations []string `json:"recommendations"` // Suggested improvements
}

// HealthRating classifies the overall health score.
type HealthRating string

const (
	// HealthRatingExcellent indicates score 85-100: Well-maintained architecture.
	HealthRatingExcellent HealthRating = "excellent"
	// HealthRatingGood indicates score 70-84: Minor improvements possible.
	HealthRatingGood HealthRating = "good"
	// HealthRatingFair indicates score 50-69: Attention needed.
	HealthRatingFair HealthRating = "fair"
	// HealthRatingPoor indicates score 30-49: Significant issues.
	HealthRatingPoor HealthRating = "poor"
	// HealthRatingCritical indicates score 0-29: Immediate action required.
	HealthRatingCritical HealthRating = "critical"
)

// GetHealthRating returns the rating for a given score.
// Thresholds: excellent (85-100), good (70-84), fair (50-69), poor (30-49), critical (0-29).
func GetHealthRating(score int) HealthRating {
	switch {
	case score >= 85:
		return HealthRatingExcellent
	case score >= 70:
		return HealthRatingGood
	case score >= 50:
		return HealthRatingFair
	case score >= 30:
		return HealthRatingPoor
	default:
		return HealthRatingCritical
	}
}
