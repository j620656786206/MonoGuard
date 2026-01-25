// Package types defines Go types that match TypeScript definitions in @monoguard/types.
// This file contains circular dependency types for Story 2.3 and Story 3.2.
package types

import "fmt"

// ========================================
// Circular Dependency Types (Story 2.3)
// ========================================

// CircularDependencyInfo represents a detected circular dependency.
// Matches @monoguard/types CircularDependencyInfo interface.
type CircularDependencyInfo struct {
	Cycle                 []string               `json:"cycle"`                           // Package names in order, ends with first
	Type                  CircularType           `json:"type"`                            // direct or indirect
	Severity              CircularSeverity       `json:"severity"`                        // critical, warning, or info
	Depth                 int                    `json:"depth"`                           // Number of unique packages in cycle
	Impact                string                 `json:"impact"`                          // Human-readable impact description
	Complexity            int                    `json:"complexity"`                      // Legacy: Basic refactoring complexity (1-10)
	RefactoringComplexity *RefactoringComplexity `json:"refactoringComplexity,omitempty"` // Story 3.5: Detailed refactoring complexity
	RootCause             *RootCauseAnalysis     `json:"rootCause,omitempty"`             // Story 3.1: Optional root cause analysis
	ImportTraces          []ImportTrace          `json:"importTraces,omitempty"`          // Story 3.2: Import statements forming the cycle
	FixStrategies         []FixStrategy          `json:"fixStrategies,omitempty"`         // Story 3.3: Recommended fix strategies
	ImpactAssessment      *ImpactAssessment      `json:"impactAssessment,omitempty"`      // Story 3.6: Impact assessment
	QuickFix              *QuickFixSummary       `json:"quickFix,omitempty"`              // Story 3.8: Quick access to best fix
	PriorityScore         float64                `json:"priorityScore"`                   // Story 3.8: Priority for sorting (higher = fix first)
}

// CircularType classifies the cycle length.
// Matches @monoguard/types CircularType union type.
type CircularType string

const (
	CircularTypeDirect   CircularType = "direct"   // 2 packages: A ↔ B
	CircularTypeIndirect CircularType = "indirect" // 3+ packages: A → B → C → A
)

// CircularSeverity indicates how problematic the cycle is.
// Matches @monoguard/types CircularSeverity union type.
type CircularSeverity string

const (
	CircularSeverityCritical CircularSeverity = "critical" // Self-loop or blocking build
	CircularSeverityWarning  CircularSeverity = "warning"  // Should be fixed
	CircularSeverityInfo     CircularSeverity = "info"     // Nice to fix
)

// NewCircularDependencyInfo creates a new CircularDependencyInfo with calculated fields.
func NewCircularDependencyInfo(cycle []string) *CircularDependencyInfo {
	if len(cycle) == 0 {
		return nil
	}

	// Calculate depth (unique packages, excluding the closing node)
	depth := len(cycle) - 1
	if depth < 1 {
		depth = 1
	}

	// Determine type
	cycleType := CircularTypeIndirect
	if depth <= 2 {
		cycleType = CircularTypeDirect
	}

	// Determine severity
	severity := classifySeverity(cycle, cycleType)

	// Generate impact description
	impact := generateImpactDescription(cycle)

	// Calculate complexity
	complexity := calculateBaseComplexity(depth)

	return &CircularDependencyInfo{
		Cycle:      cycle,
		Type:       cycleType,
		Severity:   severity,
		Depth:      depth,
		Impact:     impact,
		Complexity: complexity,
	}
}

// classifySeverity determines the severity based on cycle characteristics.
func classifySeverity(cycle []string, cycleType CircularType) CircularSeverity {
	// Self-loop is always critical
	if len(cycle) == 2 && cycle[0] == cycle[1] {
		return CircularSeverityCritical
	}

	// Direct cycles (A ↔ B) are warnings
	if cycleType == CircularTypeDirect {
		return CircularSeverityWarning
	}

	// Indirect cycles (3+ packages) are info
	return CircularSeverityInfo
}

// generateImpactDescription creates a human-readable description of the cycle.
func generateImpactDescription(cycle []string) string {
	if len(cycle) == 0 {
		return ""
	}

	// Self-loop
	if len(cycle) == 2 && cycle[0] == cycle[1] {
		return fmt.Sprintf("Self-referencing package: %s", cycle[0])
	}

	// Direct cycle (2 unique packages)
	if len(cycle) == 3 {
		return fmt.Sprintf("Direct circular dependency between %s and %s", cycle[0], cycle[1])
	}

	// Indirect cycle (3+ unique packages)
	depth := len(cycle) - 1
	return fmt.Sprintf("Indirect circular dependency involving %d packages", depth)
}

// Complexity constants for refactoring effort estimation.
// These values represent the relative difficulty of breaking cycles.
const (
	// ComplexitySelfLoop is for self-referencing packages (trivial to fix).
	ComplexitySelfLoop = 1
	// ComplexityDirect is for A ↔ B cycles (requires interface extraction).
	ComplexityDirect = 3
	// ComplexityShortIndirect is for 3-4 package cycles (moderate effort).
	ComplexityShortIndirect = 5
	// ComplexityMediumIndirect is for 5-6 package cycles (significant refactoring).
	ComplexityMediumIndirect = 7
	// ComplexityMax is the maximum complexity score.
	ComplexityMax = 10
)

// calculateBaseComplexity estimates refactoring effort (1-10).
// Complexity scales with cycle depth:
// - Self-loop (depth 1): Trivial fix, just remove self-reference
// - Direct (depth 2): Extract shared interface or dependency injection
// - Short indirect (depth 3-4): Moderate refactoring, may need new module
// - Long indirect (depth 5+): Significant architectural changes needed
func calculateBaseComplexity(depth int) int {
	switch {
	case depth <= 1:
		return ComplexitySelfLoop
	case depth == 2:
		return ComplexityDirect
	case depth <= 4:
		return ComplexityShortIndirect
	case depth <= 6:
		return ComplexityMediumIndirect
	default:
		if depth > ComplexityMax {
			return ComplexityMax
		}
		return depth
	}
}
