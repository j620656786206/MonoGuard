// Package analyzer provides dependency graph analysis for monorepo workspaces.
package analyzer

import (
	"sort"
	"strings"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// CycleDetector finds circular dependencies in a dependency graph.
// Uses Tarjan's strongly connected components algorithm for O(V+E) detection.
type CycleDetector struct {
	graph      *types.DependencyGraph
	index      int                // Global index counter for Tarjan's algorithm
	stack      []string           // Stack for Tarjan's algorithm
	onStack    map[string]bool    // Track which nodes are on stack
	indices    map[string]int     // Discovery index for each node
	lowLinks   map[string]int     // Lowest reachable index for each node
	components [][]string         // Strongly connected components found
	adjacency  map[string][]string // Adjacency list built from graph
}

// NewCycleDetector creates a new detector for the given graph.
func NewCycleDetector(graph *types.DependencyGraph) *CycleDetector {
	return &CycleDetector{
		graph:      graph,
		index:      0,
		stack:      []string{},
		onStack:    make(map[string]bool),
		indices:    make(map[string]int),
		lowLinks:   make(map[string]int),
		components: [][]string{},
		adjacency:  buildAdjacencyList(graph),
	}
}

// buildAdjacencyList creates an adjacency list from the dependency graph.
// Uses the Nodes' Dependencies field for traversal.
func buildAdjacencyList(graph *types.DependencyGraph) map[string][]string {
	adj := make(map[string][]string)

	for name, node := range graph.Nodes {
		// Combine all dependency types for cycle detection
		deps := make([]string, 0)
		deps = append(deps, node.Dependencies...)
		deps = append(deps, node.DevDependencies...)
		deps = append(deps, node.PeerDependencies...)
		deps = append(deps, node.OptionalDependencies...)

		// Sort for deterministic ordering
		sort.Strings(deps)
		adj[name] = deps
	}

	return adj
}

// DetectCycles finds all circular dependencies in the graph.
// Returns a slice of CircularDependencyInfo, sorted by severity then by cycle length.
func (cd *CycleDetector) DetectCycles() []*types.CircularDependencyInfo {
	// Run Tarjan's algorithm
	cd.tarjanSCC()

	// Extract cycles from SCCs
	cycles := cd.extractCycles()

	// Sort cycles: critical first, then warning, then info, then by depth
	sort.Slice(cycles, func(i, j int) bool {
		// Severity order: critical < warning < info
		severityOrder := map[types.CircularSeverity]int{
			types.CircularSeverityCritical: 0,
			types.CircularSeverityWarning:  1,
			types.CircularSeverityInfo:     2,
		}
		if severityOrder[cycles[i].Severity] != severityOrder[cycles[j].Severity] {
			return severityOrder[cycles[i].Severity] < severityOrder[cycles[j].Severity]
		}
		// Then by depth (shorter cycles first)
		if cycles[i].Depth != cycles[j].Depth {
			return cycles[i].Depth < cycles[j].Depth
		}
		// Then alphabetically by first node
		if len(cycles[i].Cycle) > 0 && len(cycles[j].Cycle) > 0 {
			return cycles[i].Cycle[0] < cycles[j].Cycle[0]
		}
		return false
	})

	return cycles
}

// tarjanSCC implements Tarjan's strongly connected components algorithm.
func (cd *CycleDetector) tarjanSCC() {
	// Get all node names sorted for deterministic order
	nodes := make([]string, 0, len(cd.graph.Nodes))
	for name := range cd.graph.Nodes {
		nodes = append(nodes, name)
	}
	sort.Strings(nodes)

	// Run strongConnect on each unvisited node
	for _, node := range nodes {
		if _, visited := cd.indices[node]; !visited {
			cd.strongConnect(node)
		}
	}
}

// strongConnect is the recursive part of Tarjan's algorithm.
func (cd *CycleDetector) strongConnect(v string) {
	// Set the depth index for v to the smallest unused index
	cd.indices[v] = cd.index
	cd.lowLinks[v] = cd.index
	cd.index++
	cd.stack = append(cd.stack, v)
	cd.onStack[v] = true

	// Consider successors of v
	for _, w := range cd.adjacency[v] {
		if _, visited := cd.indices[w]; !visited {
			// Successor w has not yet been visited; recurse on it
			cd.strongConnect(w)
			// After recursion, update lowLink
			if cd.lowLinks[w] < cd.lowLinks[v] {
				cd.lowLinks[v] = cd.lowLinks[w]
			}
		} else if cd.onStack[w] {
			// Successor w is in stack and hence in the current SCC
			// Use w's index (not lowLink) to avoid issues with cross-edges
			if cd.indices[w] < cd.lowLinks[v] {
				cd.lowLinks[v] = cd.indices[w]
			}
		}
	}

	// If v is a root node, pop the stack and generate an SCC
	if cd.lowLinks[v] == cd.indices[v] {
		scc := []string{}
		for {
			w := cd.stack[len(cd.stack)-1]
			cd.stack = cd.stack[:len(cd.stack)-1]
			cd.onStack[w] = false
			scc = append(scc, w)
			if w == v {
				break
			}
		}
		cd.components = append(cd.components, scc)
	}
}

// extractCycles converts SCCs to CircularDependencyInfo.
// Only SCCs with 2+ nodes are cycles (or 1 node with self-loop).
func (cd *CycleDetector) extractCycles() []*types.CircularDependencyInfo {
	var cycles []*types.CircularDependencyInfo
	seenCycles := make(map[string]bool) // Track normalized cycles to avoid duplicates

	for _, scc := range cd.components {
		if len(scc) > 1 {
			// Multi-node SCC = cycle
			cyclePath := cd.buildCyclePath(scc)
			normalized := normalizeCycle(cyclePath)
			cycleKey := cyclesToKey(normalized)

			if !seenCycles[cycleKey] {
				seenCycles[cycleKey] = true
				info := types.NewCircularDependencyInfo(normalized)
				if info != nil {
					cycles = append(cycles, info)
				}
			}
		} else if len(scc) == 1 {
			// Check for self-loop
			node := scc[0]
			if cd.hasSelfLoop(node) {
				cycle := []string{node, node}
				cycleKey := cyclesToKey(cycle)

				if !seenCycles[cycleKey] {
					seenCycles[cycleKey] = true
					info := types.NewCircularDependencyInfo(cycle)
					if info != nil {
						cycles = append(cycles, info)
					}
				}
			}
		}
	}

	return cycles
}

// buildCyclePath constructs the cycle path from SCC nodes.
// The path starts with the lexicographically smallest node and ends with the same node.
func (cd *CycleDetector) buildCyclePath(scc []string) []string {
	if len(scc) == 0 {
		return nil
	}

	// Sort SCC to find starting node (smallest lexicographically)
	sorted := make([]string, len(scc))
	copy(sorted, scc)
	sort.Strings(sorted)

	// Build adjacency within SCC
	sccSet := make(map[string]bool)
	for _, n := range scc {
		sccSet[n] = true
	}

	// Use DFS to find actual cycle path
	start := sorted[0]
	path := cd.findCyclePath(start, sccSet)

	if path == nil {
		// Fallback: just return sorted nodes with closing node
		result := make([]string, len(sorted)+1)
		copy(result, sorted)
		result[len(sorted)] = sorted[0]
		return result
	}

	return path
}

// findCyclePath uses DFS to find a path through all SCC nodes back to start.
func (cd *CycleDetector) findCyclePath(start string, sccSet map[string]bool) []string {
	visited := make(map[string]bool)
	path := []string{start}
	visited[start] = true

	current := start
	for {
		// Find next unvisited neighbor in SCC
		var next string
		for _, neighbor := range cd.adjacency[current] {
			if sccSet[neighbor] {
				if neighbor == start && len(visited) == len(sccSet) {
					// Found complete cycle back to start
					return append(path, start)
				}
				if !visited[neighbor] {
					next = neighbor
					break
				}
			}
		}

		if next == "" {
			// Check if we can close the cycle
			for _, neighbor := range cd.adjacency[current] {
				if neighbor == start && len(visited) == len(sccSet) {
					return append(path, start)
				}
			}
			// Dead end - can't complete cycle from this path
			// This shouldn't happen in a true SCC, but handle gracefully
			break
		}

		path = append(path, next)
		visited[next] = true
		current = next
	}

	// Fallback: return what we have
	if len(path) > 0 && path[len(path)-1] != start {
		path = append(path, start)
	}
	return path
}

// hasSelfLoop checks if a node depends on itself.
func (cd *CycleDetector) hasSelfLoop(node string) bool {
	pkg, exists := cd.graph.Nodes[node]
	if !exists {
		return false
	}

	// Check all dependency types
	for _, dep := range pkg.Dependencies {
		if dep == node {
			return true
		}
	}
	for _, dep := range pkg.DevDependencies {
		if dep == node {
			return true
		}
	}
	for _, dep := range pkg.PeerDependencies {
		if dep == node {
			return true
		}
	}
	for _, dep := range pkg.OptionalDependencies {
		if dep == node {
			return true
		}
	}

	return false
}

// normalizeCycle converts a cycle to its canonical form.
// The canonical form starts with the lexicographically smallest node.
func normalizeCycle(cycle []string) []string {
	if len(cycle) <= 1 {
		return cycle
	}

	// Remove trailing duplicate if present
	nodes := cycle
	if len(cycle) > 1 && cycle[0] == cycle[len(cycle)-1] {
		nodes = cycle[:len(cycle)-1]
	}

	if len(nodes) == 0 {
		return cycle
	}

	// Find index of smallest node
	minIdx := 0
	for i, node := range nodes {
		if node < nodes[minIdx] {
			minIdx = i
		}
	}

	// Rotate to start with smallest
	result := make([]string, len(nodes)+1)
	for i := 0; i < len(nodes); i++ {
		result[i] = nodes[(minIdx+i)%len(nodes)]
	}
	result[len(nodes)] = result[0] // Close the cycle

	return result
}

// cyclesToKey creates a unique string key for a cycle.
// Uses strings.Join for better performance than string concatenation.
func cyclesToKey(cycle []string) string {
	if len(cycle) == 0 {
		return ""
	}
	return strings.Join(cycle, "|")
}
