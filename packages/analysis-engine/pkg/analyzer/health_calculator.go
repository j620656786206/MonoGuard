// Package analyzer provides dependency graph analysis for monorepo workspaces.
// This file implements the health score calculator for Story 2.5.
package analyzer

import (
	"fmt"
	"math"
	"time"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// Weight constants for health score factors (must sum to 1.0).
const (
	WeightCircular = 0.40 // Circular dependencies - most impactful
	WeightConflict = 0.25 // Version conflicts
	WeightDepth    = 0.20 // Dependency depth
	WeightCoupling = 0.15 // Package coupling
)

// Deduction constants for circular dependencies.
const (
	DeductionSelfLoop   = 25 // Self-referencing package
	DeductionDirect     = 15 // Direct cycle (2 packages)
	DeductionIndirect   = 10 // Indirect cycle (3+ packages)
)

// Deduction constants for version conflicts.
const (
	DeductionCritical = 10 // Major version difference
	DeductionWarning  = 5  // Minor version difference
	DeductionInfo     = 2  // Patch version difference
)

// Depth scoring constants.
// Note: AC6 mentions "3-4 levels" as optimal range. We use 4 as the threshold,
// meaning depths of 1-4 get full score (100), and deductions start at depth 5+.
// This is a design choice that treats 4 as the upper bound of "optimal".
const (
	OptimalDepth       = 4  // Ideal max depth (upper bound of 3-4 range per AC6)
	DeductionPerLevel  = 10 // Points deducted per level above optimal
	AvgDepthMultiplier = 5  // Multiplier for average depth penalty
)

// HealthCalculator computes architecture health scores.
type HealthCalculator struct {
	graph     *types.DependencyGraph
	cycles    []*types.CircularDependencyInfo
	conflicts []*types.VersionConflictInfo
}

// NewHealthCalculator creates a new calculator with analysis results.
func NewHealthCalculator(
	graph *types.DependencyGraph,
	cycles []*types.CircularDependencyInfo,
	conflicts []*types.VersionConflictInfo,
) *HealthCalculator {
	return &HealthCalculator{
		graph:     graph,
		cycles:    cycles,
		conflicts: conflicts,
	}
}

// Calculate computes the complete health score with breakdown.
func (hc *HealthCalculator) Calculate() *types.HealthScoreResult {
	// Calculate individual factor scores
	circularScore, circularFactor := hc.calculateCircularScore()
	conflictScore, conflictFactor := hc.calculateConflictScore()
	depthScore, depthFactor := hc.calculateDepthScore()
	couplingScore, couplingFactor := hc.calculateCouplingScore()

	// Calculate weighted overall score
	weighted := float64(circularScore)*WeightCircular +
		float64(conflictScore)*WeightConflict +
		float64(depthScore)*WeightDepth +
		float64(couplingScore)*WeightCoupling

	overall := int(math.Round(weighted))
	overall = boundScore(overall)

	// Set weighted scores on factors
	circularFactor.WeightedScore = int(math.Round(float64(circularScore) * WeightCircular))
	conflictFactor.WeightedScore = int(math.Round(float64(conflictScore) * WeightConflict))
	depthFactor.WeightedScore = int(math.Round(float64(depthScore) * WeightDepth))
	couplingFactor.WeightedScore = int(math.Round(float64(couplingScore) * WeightCoupling))

	return &types.HealthScoreResult{
		Overall: overall,
		Rating:  types.GetHealthRating(overall),
		Breakdown: &types.ScoreBreakdown{
			CircularScore: circularScore,
			ConflictScore: conflictScore,
			DepthScore:    depthScore,
			CouplingScore: couplingScore,
		},
		Factors: []*types.HealthFactor{
			circularFactor,
			conflictFactor,
			depthFactor,
			couplingFactor,
		},
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// calculateCircularScore computes the score from circular dependencies.
// Formula: 100 - (selfLoops * 25 + directCycles * 15 + indirectCycles * 10)
func (hc *HealthCalculator) calculateCircularScore() (int, *types.HealthFactor) {
	if len(hc.cycles) == 0 {
		return 100, &types.HealthFactor{
			Name:            "Circular Dependencies",
			Score:           100,
			Weight:          WeightCircular,
			Description:     "No circular dependencies detected",
			Recommendations: []string{},
		}
	}

	deductions := 0
	selfLoopCount := 0
	directCount := 0
	indirectCount := 0

	for _, cycle := range hc.cycles {
		// Self-loop detection: A package that depends on itself.
		// We check two conditions for robustness:
		//   1. cycle.Depth == 1: Standard self-loop indicator from cycle detector
		//   2. Cycle array has 2 identical elements: Handles edge case where cycle
		//      representation is ["A", "A"] (package A depending on itself)
		// Both conditions are checked to handle different cycle detector implementations.
		if cycle.Depth == 1 || (len(cycle.Cycle) == 2 && cycle.Cycle[0] == cycle.Cycle[1]) {
			selfLoopCount++
			deductions += DeductionSelfLoop
		} else if cycle.Type == types.CircularTypeDirect {
			directCount++
			deductions += DeductionDirect
		} else {
			indirectCount++
			deductions += DeductionIndirect
		}
	}

	score := boundScore(100 - deductions)
	recommendations := generateCircularRecommendations(selfLoopCount, directCount, indirectCount)

	return score, &types.HealthFactor{
		Name:            "Circular Dependencies",
		Score:           score,
		Weight:          WeightCircular,
		Description:     fmt.Sprintf("%d cycles detected", len(hc.cycles)),
		Recommendations: recommendations,
	}
}

// calculateConflictScore computes the score from version conflicts.
// Formula: 100 - (critical * 10 + warning * 5 + info * 2)
func (hc *HealthCalculator) calculateConflictScore() (int, *types.HealthFactor) {
	if len(hc.conflicts) == 0 {
		return 100, &types.HealthFactor{
			Name:            "Version Conflicts",
			Score:           100,
			Weight:          WeightConflict,
			Description:     "No version conflicts detected",
			Recommendations: []string{},
		}
	}

	deductions := 0
	criticalCount := 0
	warningCount := 0
	infoCount := 0

	for _, conflict := range hc.conflicts {
		switch conflict.Severity {
		case types.ConflictSeverityCritical:
			criticalCount++
			deductions += DeductionCritical
		case types.ConflictSeverityWarning:
			warningCount++
			deductions += DeductionWarning
		case types.ConflictSeverityInfo:
			infoCount++
			deductions += DeductionInfo
		}
	}

	score := boundScore(100 - deductions)
	recommendations := generateConflictRecommendations(criticalCount, warningCount, infoCount)

	return score, &types.HealthFactor{
		Name:            "Version Conflicts",
		Score:           score,
		Weight:          WeightConflict,
		Description:     fmt.Sprintf("%d conflicts detected", len(hc.conflicts)),
		Recommendations: recommendations,
	}
}

// calculateDepthScore computes the score from dependency depth.
// Optimal depth: 3-4 levels = 100 points.
// Deduct 10 points per level above optimal for max depth.
// Deduct 5 points per level above optimal for average depth.
func (hc *HealthCalculator) calculateDepthScore() (int, *types.HealthFactor) {
	if hc.graph == nil || len(hc.graph.Nodes) == 0 {
		return 100, &types.HealthFactor{
			Name:            "Dependency Depth",
			Score:           100,
			Weight:          WeightDepth,
			Description:     "No packages to analyze",
			Recommendations: []string{},
		}
	}

	maxDepth, avgDepth := hc.calculateDepthMetrics()

	score := 100
	if maxDepth > OptimalDepth {
		score -= (maxDepth - OptimalDepth) * DeductionPerLevel
	}
	if avgDepth > float64(OptimalDepth) {
		score -= int((avgDepth - float64(OptimalDepth)) * AvgDepthMultiplier)
	}

	score = boundScore(score)
	recommendations := generateDepthRecommendations(maxDepth, avgDepth)

	return score, &types.HealthFactor{
		Name:            "Dependency Depth",
		Score:           score,
		Weight:          WeightDepth,
		Description:     fmt.Sprintf("Max depth: %d, Avg depth: %.1f", maxDepth, avgDepth),
		Recommendations: recommendations,
	}
}

// calculateDepthMetrics computes max and average dependency depth using memoized DFS.
func (hc *HealthCalculator) calculateDepthMetrics() (maxDepth int, avgDepth float64) {
	if hc.graph == nil || len(hc.graph.Nodes) == 0 {
		return 0, 0.0
	}

	// Build adjacency list for efficient traversal
	adjList := make(map[string][]string)
	for _, node := range hc.graph.Nodes {
		adjList[node.Name] = node.Dependencies
	}

	// Memoization cache for computed depths
	memo := make(map[string]int)

	// Calculate depth for each package using memoized DFS
	depths := make([]int, 0, len(hc.graph.Nodes))

	for pkgName := range hc.graph.Nodes {
		inPath := make(map[string]bool)
		depth := hc.dfsDepthMemo(pkgName, adjList, memo, inPath)
		depths = append(depths, depth)
		if depth > maxDepth {
			maxDepth = depth
		}
	}

	// Calculate average depth
	if len(depths) > 0 {
		total := 0
		for _, d := range depths {
			total += d
		}
		avgDepth = float64(total) / float64(len(depths))
	}

	return maxDepth, avgDepth
}

// dfsDepthMemo performs memoized DFS to find the maximum depth from a node.
// Uses memo for caching computed depths and inPath for cycle detection.
func (hc *HealthCalculator) dfsDepthMemo(pkg string, adjList map[string][]string, memo map[string]int, inPath map[string]bool) int {
	// Check memo first
	if depth, ok := memo[pkg]; ok {
		return depth
	}

	// Cycle detection - if we're already in the current path, it's a cycle
	if inPath[pkg] {
		return 0
	}

	deps := adjList[pkg]
	if len(deps) == 0 {
		memo[pkg] = 0
		return 0 // Leaf node
	}

	inPath[pkg] = true
	maxChildDepth := 0

	for _, dep := range deps {
		childDepth := hc.dfsDepthMemo(dep, adjList, memo, inPath)
		if childDepth > maxChildDepth {
			maxChildDepth = childDepth
		}
	}

	inPath[pkg] = false
	depth := maxChildDepth + 1
	memo[pkg] = depth
	return depth
}

// calculateCouplingScore computes the score from package coupling metrics.
// Uses instability metric: I = Ce / (Ca + Ce)
// Ideal average instability is 0.5 (balanced between stable and unstable).
func (hc *HealthCalculator) calculateCouplingScore() (int, *types.HealthFactor) {
	if hc.graph == nil || len(hc.graph.Nodes) == 0 {
		return 100, &types.HealthFactor{
			Name:            "Package Coupling",
			Score:           100,
			Weight:          WeightCoupling,
			Description:     "No packages to analyze",
			Recommendations: []string{},
		}
	}

	metrics := hc.calculateCouplingMetrics()

	// Score: 100 for 0.5 instability (balanced), decreases as it deviates
	deviation := math.Abs(metrics.AverageInstability - 0.5)
	score := int(100 - (deviation * 100))
	score = boundScore(score)

	recommendations := generateCouplingRecommendations(metrics)

	return score, &types.HealthFactor{
		Name:            "Package Coupling",
		Score:           score,
		Weight:          WeightCoupling,
		Description:     fmt.Sprintf("Avg instability: %.2f", metrics.AverageInstability),
		Recommendations: recommendations,
	}
}

// CouplingMetrics holds coupling analysis results.
type CouplingMetrics struct {
	AverageInstability float64
	HighCoupling       []string                    // Packages with concerning coupling
	PackageMetrics     map[string]*PackageCoupling // Per-package metrics
}

// PackageCoupling holds Ca, Ce, and instability for a package.
type PackageCoupling struct {
	AfferentCoupling int     // Ca - packages depending on this
	EfferentCoupling int     // Ce - packages this depends on
	Instability      float64 // Ce / (Ca + Ce)
}

// calculateCouplingMetrics computes coupling metrics for all packages.
func (hc *HealthCalculator) calculateCouplingMetrics() *CouplingMetrics {
	metrics := &CouplingMetrics{
		PackageMetrics: make(map[string]*PackageCoupling),
		HighCoupling:   []string{},
	}

	if hc.graph == nil || len(hc.graph.Nodes) == 0 {
		metrics.AverageInstability = 0.5 // Default to balanced
		return metrics
	}

	// Count afferent coupling (Ca) for each package
	afferentCount := make(map[string]int)
	for _, node := range hc.graph.Nodes {
		for _, dep := range node.Dependencies {
			afferentCount[dep]++
		}
	}

	// Calculate metrics for each package
	totalInstability := 0.0
	packageCount := 0

	for name, node := range hc.graph.Nodes {
		ca := afferentCount[name]              // Packages depending on this
		ce := len(node.Dependencies)           // Packages this depends on

		var instability float64
		if ca+ce > 0 {
			instability = float64(ce) / float64(ca+ce)
		} else {
			instability = 0.5 // No dependencies = neutral
		}

		metrics.PackageMetrics[name] = &PackageCoupling{
			AfferentCoupling: ca,
			EfferentCoupling: ce,
			Instability:      instability,
		}

		totalInstability += instability
		packageCount++

		// Flag packages with extreme instability
		if instability < 0.2 || instability > 0.8 {
			metrics.HighCoupling = append(metrics.HighCoupling, name)
		}
	}

	if packageCount > 0 {
		metrics.AverageInstability = totalInstability / float64(packageCount)
	} else {
		metrics.AverageInstability = 0.5
	}

	return metrics
}

// boundScore ensures score is within 0-100 range.
func boundScore(score int) int {
	if score < 0 {
		return 0
	}
	if score > 100 {
		return 100
	}
	return score
}

// generateCircularRecommendations creates recommendations for circular dependencies.
func generateCircularRecommendations(selfLoops, direct, indirect int) []string {
	recommendations := []string{}

	if selfLoops > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("Fix %d self-referencing package(s) - these are critical issues", selfLoops))
	}

	if direct > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("Break %d direct cycle(s) by extracting shared code into separate packages", direct))
	}

	if indirect > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("Refactor %d indirect cycle(s) - consider dependency inversion", indirect))
	}

	return recommendations
}

// generateConflictRecommendations creates recommendations for version conflicts.
func generateConflictRecommendations(critical, warning, info int) []string {
	recommendations := []string{}

	if critical > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("Resolve %d critical version conflict(s) - major version differences may cause breaking changes", critical))
	}

	if warning > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("Consider aligning %d minor version conflict(s) for consistency", warning))
	}

	if info > 0 && critical == 0 && warning == 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("%d patch-level conflict(s) detected - low priority but can be aligned", info))
	}

	return recommendations
}

// generateDepthRecommendations creates recommendations for dependency depth.
func generateDepthRecommendations(maxDepth int, avgDepth float64) []string {
	recommendations := []string{}

	if maxDepth > OptimalDepth+2 {
		recommendations = append(recommendations,
			fmt.Sprintf("Dependency chain is too deep (max: %d). Consider flattening the architecture", maxDepth))
	} else if maxDepth > OptimalDepth {
		recommendations = append(recommendations,
			fmt.Sprintf("Max depth of %d is slightly above optimal (%d). Review if simplification is possible", maxDepth, OptimalDepth))
	}

	if avgDepth > float64(OptimalDepth) {
		recommendations = append(recommendations,
			fmt.Sprintf("Average depth of %.1f suggests many packages have deep dependencies", avgDepth))
	}

	return recommendations
}

// generateCouplingRecommendations creates recommendations for coupling metrics.
func generateCouplingRecommendations(metrics *CouplingMetrics) []string {
	recommendations := []string{}

	if metrics.AverageInstability < 0.3 {
		recommendations = append(recommendations,
			"Architecture is overly stable - consider if this limits flexibility")
	} else if metrics.AverageInstability > 0.7 {
		recommendations = append(recommendations,
			"Architecture is highly unstable - consider adding stable foundational packages")
	}

	if len(metrics.HighCoupling) > 0 && len(metrics.HighCoupling) <= 3 {
		recommendations = append(recommendations,
			fmt.Sprintf("Review coupling in: %v", metrics.HighCoupling))
	} else if len(metrics.HighCoupling) > 3 {
		recommendations = append(recommendations,
			fmt.Sprintf("%d packages have extreme coupling - architectural review recommended", len(metrics.HighCoupling)))
	}

	return recommendations
}
