// Package analyzer provides dependency graph analysis for monorepo workspaces.
// This file implements root cause analysis for circular dependencies (Story 3.1).
package analyzer

import (
	"fmt"
	"strings"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// RootCauseAnalyzer determines the root cause of circular dependencies.
type RootCauseAnalyzer struct {
	graph     *types.DependencyGraph
	adjacency map[string][]string // Cached adjacency list
}

// NewRootCauseAnalyzer creates a new analyzer for the given graph.
func NewRootCauseAnalyzer(graph *types.DependencyGraph) *RootCauseAnalyzer {
	return &RootCauseAnalyzer{
		graph:     graph,
		adjacency: buildAdjacencyList(graph),
	}
}

// Analyze determines root cause for a circular dependency.
// Returns nil if cycle is invalid (nil or less than 2 nodes).
func (rca *RootCauseAnalyzer) Analyze(cycle *types.CircularDependencyInfo) *types.RootCauseAnalysis {
	if cycle == nil || len(cycle.Cycle) < 2 {
		return nil
	}

	// Build dependency chain from cycle
	chain := rca.buildDependencyChain(cycle.Cycle)
	if len(chain) == 0 {
		return nil
	}

	// Identify root cause package and confidence score
	originatingPackage, confidence := rca.identifyOriginatingPackage(cycle.Cycle, chain)

	// Find critical edge (best edge to break the cycle)
	criticalEdge := rca.findCriticalEdge(chain)

	// Determine problematic dependency (first edge from originating package)
	problematicDep := rca.findProblematicDependency(originatingPackage, chain)

	// Generate human-readable explanation
	explanation := generateExplanation(originatingPackage, criticalEdge, confidence)

	return types.NewRootCauseAnalysis(
		originatingPackage,
		problematicDep,
		confidence,
		explanation,
		chain,
		criticalEdge,
	)
}

// buildDependencyChain creates the ordered edge list for the cycle.
func (rca *RootCauseAnalyzer) buildDependencyChain(cycle []string) []types.RootCauseEdge {
	if len(cycle) < 2 {
		return nil
	}

	// Number of edges = len(cycle) - 1 (since last node == first node)
	numEdges := len(cycle) - 1
	edges := make([]types.RootCauseEdge, numEdges)

	for i := 0; i < numEdges; i++ {
		from := cycle[i]
		to := cycle[i+1]
		depType := rca.getDependencyType(from, to)

		edges[i] = types.RootCauseEdge{
			From:     from,
			To:       to,
			Type:     depType,
			Critical: false, // Will be set by findCriticalEdge
		}
	}

	return edges
}

// getDependencyType determines the type of dependency from one package to another.
func (rca *RootCauseAnalyzer) getDependencyType(from, to string) types.DependencyType {
	node, exists := rca.graph.Nodes[from]
	if !exists {
		return types.DependencyTypeProduction // Default
	}

	// Check each dependency type in priority order
	for _, dep := range node.Dependencies {
		if dep == to {
			return types.DependencyTypeProduction
		}
	}
	for _, dep := range node.DevDependencies {
		if dep == to {
			return types.DependencyTypeDevelopment
		}
	}
	for _, dep := range node.PeerDependencies {
		if dep == to {
			return types.DependencyTypePeer
		}
	}
	for _, dep := range node.OptionalDependencies {
		if dep == to {
			return types.DependencyTypeOptional
		}
	}

	return types.DependencyTypeProduction // Default
}

// identifyOriginatingPackage determines which package is most likely the root cause.
// Returns the package name and confidence score.
func (rca *RootCauseAnalyzer) identifyOriginatingPackage(cycle []string, chain []types.RootCauseEdge) (string, int) {
	// Handle self-loop (A → A)
	if len(cycle) == 2 && cycle[0] == cycle[1] {
		return cycle[0], 100
	}

	// Get unique packages (exclude closing node)
	packages := cycle[:len(cycle)-1]

	var bestPackage string
	var bestScore int

	for i, pkg := range packages {
		score := rca.calculateTotalScore(pkg, packages, i)
		if score > bestScore || (score == bestScore && pkg < bestPackage) {
			bestScore = score
			bestPackage = pkg
		}
	}

	return bestPackage, bestScore
}

// calculateTotalScore combines all heuristics for a package.
// Max score: 30 (incoming) + 20 (outgoing) + 25 (name) + 15 (position) + 10 (edge type) = 100
func (rca *RootCauseAnalyzer) calculateTotalScore(pkg string, cycle []string, position int) int {
	score := 0
	score += rca.calculateIncomingDepsScore(pkg)
	score += rca.calculateOutgoingDepsScore(pkg)
	score += rca.calculateNamePatternScore(pkg)
	score += rca.calculatePositionScore(pkg, cycle, position)
	return score
}

// calculateIncomingDepsScore scores packages based on incoming dependencies.
// Packages with fewer incoming deps are more likely to be high-level
// and thus more likely to be the root cause (they shouldn't depend on lower-level).
func (rca *RootCauseAnalyzer) calculateIncomingDepsScore(pkg string) int {
	incoming := 0
	for _, node := range rca.graph.Nodes {
		for _, dep := range node.Dependencies {
			if dep == pkg {
				incoming++
			}
		}
		for _, dep := range node.DevDependencies {
			if dep == pkg {
				incoming++
			}
		}
		for _, dep := range node.PeerDependencies {
			if dep == pkg {
				incoming++
			}
		}
		for _, dep := range node.OptionalDependencies {
			if dep == pkg {
				incoming++
			}
		}
	}

	// Lower incoming = higher score (max 30 points)
	if incoming == 0 {
		return 30
	}
	score := 30 - incoming*5
	if score < 0 {
		return 0
	}
	return score
}

// calculateOutgoingDepsScore scores packages based on outgoing dependencies.
// Packages with more outgoing deps are more likely to be low-level
// and thus less likely to be the root cause.
func (rca *RootCauseAnalyzer) calculateOutgoingDepsScore(pkg string) int {
	node, exists := rca.graph.Nodes[pkg]
	if !exists {
		return 0
	}

	outgoing := len(node.Dependencies) +
		len(node.DevDependencies) +
		len(node.PeerDependencies) +
		len(node.OptionalDependencies)

	// More outgoing = lower score (max 20 points)
	score := 20 - outgoing*3
	if score < 0 {
		return 0
	}
	return score
}

// calculateNamePatternScore scores packages based on naming patterns.
// "Core", "common", "shared", "utils" packages are less likely to be root cause.
func (rca *RootCauseAnalyzer) calculateNamePatternScore(pkg string) int {
	lowerName := strings.ToLower(pkg)
	lowLevelPatterns := []string{"core", "common", "shared", "utils", "lib", "base", "util"}

	for _, pattern := range lowLevelPatterns {
		if strings.Contains(lowerName, pattern) {
			return 0 // Low-level package, not likely root cause
		}
	}
	return 25 // High-level package, more likely root cause
}

// calculatePositionScore scores packages based on position in cycle.
// First package in cycle (lexicographically sorted) gets bonus for consistency.
func (rca *RootCauseAnalyzer) calculatePositionScore(pkg string, cycle []string, position int) int {
	if position == 0 {
		return 15
	}
	return 0
}

// findCriticalEdge identifies the edge that would best break the cycle.
// Prioritizes: optional > peer > dev > production (easier to break first)
func (rca *RootCauseAnalyzer) findCriticalEdge(chain []types.RootCauseEdge) *types.RootCauseEdge {
	if len(chain) == 0 {
		return nil
	}

	// Edge type priority (lower = easier to break = better critical edge)
	priority := map[types.DependencyType]int{
		types.DependencyTypeOptional:    1,
		types.DependencyTypePeer:        2,
		types.DependencyTypeDevelopment: 3,
		types.DependencyTypeProduction:  4,
	}

	var bestEdge *types.RootCauseEdge
	bestPriority := 100 // Higher than any actual priority

	for i := range chain {
		edgePriority := priority[chain[i].Type]
		if edgePriority < bestPriority {
			bestPriority = edgePriority
			edge := chain[i]
			edge.Critical = true
			bestEdge = &edge
		}
	}

	return bestEdge
}

// findProblematicDependency finds the first edge from the originating package.
func (rca *RootCauseAnalyzer) findProblematicDependency(originatingPackage string, chain []types.RootCauseEdge) types.RootCauseEdge {
	for _, edge := range chain {
		if edge.From == originatingPackage {
			return edge
		}
	}

	// Fallback to first edge if originating package not in chain
	if len(chain) > 0 {
		return chain[0]
	}

	return types.RootCauseEdge{}
}

// generateExplanation creates a human-readable explanation of the root cause.
func generateExplanation(origin string, criticalEdge *types.RootCauseEdge, confidence int) string {
	var sb strings.Builder

	// Confidence level description
	confidenceLevel := "likely"
	if confidence >= 80 {
		confidenceLevel = "highly likely"
	} else if confidence < 50 {
		confidenceLevel = "possibly"
	}

	sb.WriteString(fmt.Sprintf("Package '%s' is %s the root cause of this circular dependency. ", origin, confidenceLevel))

	// Explain why based on critical edge
	if criticalEdge != nil {
		sb.WriteString(fmt.Sprintf("The dependency from '%s' to '%s' creates the problematic relationship. ",
			criticalEdge.From, criticalEdge.To))

		// Suggest action based on dependency type
		switch criticalEdge.Type {
		case types.DependencyTypeDevelopment:
			sb.WriteString("Since this is a dev dependency, it may be easier to break by restructuring test utilities.")
		case types.DependencyTypeOptional:
			sb.WriteString("Since this is an optional dependency, removing it may be the simplest fix.")
		case types.DependencyTypePeer:
			sb.WriteString("Since this is a peer dependency, consider restructuring the plugin architecture.")
		default:
			sb.WriteString("Consider extracting shared code to a new package or using dependency injection.")
		}
	} else {
		sb.WriteString("Consider extracting shared code to a new package or using dependency injection.")
	}

	return sb.String()
}
