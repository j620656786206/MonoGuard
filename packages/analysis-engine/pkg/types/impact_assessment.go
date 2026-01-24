// Package types defines Go types that match TypeScript definitions in @monoguard/types.
// This file contains impact assessment types for Story 3.6.
package types

import "fmt"

// ========================================
// Impact Assessment Types (Story 3.6)
// ========================================

// ImpactAssessment represents the blast radius analysis for a circular dependency.
// Matches @monoguard/types ImpactAssessment interface.
type ImpactAssessment struct {
	// DirectParticipants are packages directly in the cycle
	DirectParticipants []string `json:"directParticipants"`

	// IndirectDependents are packages that depend on cycle participants
	IndirectDependents []IndirectDependent `json:"indirectDependents"`

	// TotalAffected is the count of all affected packages (direct + indirect)
	TotalAffected int `json:"totalAffected"`

	// AffectedPercentage is the proportion of workspace affected (0.0-1.0)
	AffectedPercentage float64 `json:"affectedPercentage"`

	// AffectedPercentageDisplay is human-readable (e.g., "25%")
	AffectedPercentageDisplay string `json:"affectedPercentageDisplay"`

	// RiskLevel classifies the impact severity
	RiskLevel RiskLevel `json:"riskLevel"`

	// RiskExplanation describes why this risk level was assigned
	RiskExplanation string `json:"riskExplanation"`

	// RippleEffect contains visualization-ready data
	RippleEffect *RippleEffect `json:"rippleEffect,omitempty"`
}

// IndirectDependent represents a package that depends on a cycle participant.
type IndirectDependent struct {
	// PackageName is the affected package
	PackageName string `json:"packageName"`

	// DependsOn is the cycle participant this package depends on
	DependsOn string `json:"dependsOn"`

	// Distance is the number of hops from the cycle (1 = direct dependent)
	Distance int `json:"distance"`

	// DependencyPath shows the full path from cycle to this package
	DependencyPath []string `json:"dependencyPath"`
}

// RippleEffect contains data for visualization.
type RippleEffect struct {
	// Layers groups affected packages by distance from cycle
	Layers []RippleLayer `json:"layers"`

	// TotalLayers is the maximum distance from cycle
	TotalLayers int `json:"totalLayers"`
}

// RippleLayer represents packages at a specific distance from the cycle.
type RippleLayer struct {
	// Distance from the cycle (0 = direct participants, 1 = first-level dependents)
	Distance int `json:"distance"`

	// Packages at this distance
	Packages []string `json:"packages"`

	// Count of packages at this layer
	Count int `json:"count"`
}

// RiskLevel classifies the impact severity.
type RiskLevel string

const (
	RiskLevelCritical RiskLevel = "critical"
	RiskLevelHigh     RiskLevel = "high"
	RiskLevelMedium   RiskLevel = "medium"
	RiskLevelLow      RiskLevel = "low"
)

// NewImpactAssessment creates a new ImpactAssessment with initialized slices.
func NewImpactAssessment() *ImpactAssessment {
	return &ImpactAssessment{
		DirectParticipants: []string{},
		IndirectDependents: []IndirectDependent{},
	}
}

// CalculatePercentage computes the affected percentage from count and total.
func CalculatePercentage(affected, total int) (float64, string) {
	if total == 0 {
		return 0.0, "0%"
	}

	percentage := float64(affected) / float64(total)

	// Cap at 1.0
	if percentage > 1.0 {
		percentage = 1.0
	}

	// Format display string
	displayPercentage := int(percentage * 100)
	display := fmt.Sprintf("%d%%", displayPercentage)

	return percentage, display
}
