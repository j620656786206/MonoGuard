// Package types defines Go types that match TypeScript definitions in @monoguard/types.
// This file contains root cause analysis types for Story 3.1.
package types

// ========================================
// Root Cause Analysis Types (Story 3.1)
// ========================================

// RootCauseAnalysis provides insight into why a circular dependency exists.
// Matches @monoguard/types RootCauseAnalysis interface.
type RootCauseAnalysis struct {
	// OriginatingPackage is the package most likely responsible for the cycle
	OriginatingPackage string `json:"originatingPackage"`

	// ProblematicDependency is the specific dependency creating the cycle
	ProblematicDependency RootCauseEdge `json:"problematicDependency"`

	// Confidence is a score (0-100) indicating analysis certainty
	Confidence int `json:"confidence"`

	// Explanation is a human-readable description of the root cause
	Explanation string `json:"explanation"`

	// Chain is the ordered dependency chain forming the cycle
	Chain []RootCauseEdge `json:"chain"`

	// CriticalEdge is the edge most likely to break if removed
	CriticalEdge *RootCauseEdge `json:"criticalEdge,omitempty"`
}

// RootCauseEdge represents a single dependency relationship in root cause analysis.
// This is separate from DependencyEdge in graph.go as it has different fields
// (Critical instead of VersionRange) and serves a different purpose.
type RootCauseEdge struct {
	From     string         `json:"from"`     // Source package
	To       string         `json:"to"`       // Target package
	Type     DependencyType `json:"type"`     // production, development, peer, optional
	Critical bool           `json:"critical"` // If true, this edge is key to breaking cycle
}

// NewRootCauseAnalysis creates a new RootCauseAnalysis with validated fields.
// Returns nil if the cycle is invalid (nil or less than 2 nodes).
func NewRootCauseAnalysis(
	originatingPackage string,
	problematicDependency RootCauseEdge,
	confidence int,
	explanation string,
	chain []RootCauseEdge,
	criticalEdge *RootCauseEdge,
) *RootCauseAnalysis {
	// Validate confidence score bounds
	if confidence < 0 {
		confidence = 0
	}
	if confidence > 100 {
		confidence = 100
	}

	return &RootCauseAnalysis{
		OriginatingPackage:    originatingPackage,
		ProblematicDependency: problematicDependency,
		Confidence:            confidence,
		Explanation:           explanation,
		Chain:                 chain,
		CriticalEdge:          criticalEdge,
	}
}
