package analyzer

import (
	"testing"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// ========================================
// Helper Functions
// ========================================

// createTestGraph creates a DependencyGraph with the given edges.
// edges format: [][]string{{"from", "to"}, ...}
func createTestGraph(edges [][]string) *types.DependencyGraph {
	nodes := make(map[string]*types.PackageNode)

	// Collect all unique node names
	nodeNames := make(map[string]bool)
	for _, edge := range edges {
		nodeNames[edge[0]] = true
		if len(edge) > 1 {
			nodeNames[edge[1]] = true
		}
	}

	// Create nodes
	for name := range nodeNames {
		nodes[name] = types.NewPackageNode(name, "1.0.0", "packages/"+name)
	}

	// Add dependencies
	for _, edge := range edges {
		from := edge[0]
		if len(edge) > 1 {
			to := edge[1]
			nodes[from].Dependencies = append(nodes[from].Dependencies, to)
		}
	}

	return &types.DependencyGraph{
		Nodes:         nodes,
		Edges:         []*types.DependencyEdge{},
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
	}
}

// findCycleWithNodes checks if a cycle contains all specified nodes.
func findCycleWithNodes(cycles []*types.CircularDependencyInfo, expectedNodes []string) *types.CircularDependencyInfo {
	for _, cycle := range cycles {
		// Check if cycle contains all expected nodes (excluding closing node)
		cycleNodes := make(map[string]bool)
		for i := 0; i < len(cycle.Cycle)-1; i++ {
			cycleNodes[cycle.Cycle[i]] = true
		}

		allFound := true
		for _, node := range expectedNodes {
			if !cycleNodes[node] {
				allFound = false
				break
			}
		}

		if allFound && len(cycleNodes) == len(expectedNodes) {
			return cycle
		}
	}
	return nil
}

// ========================================
// No Cycles Tests
// ========================================

func TestDetectCycles_NoCycles(t *testing.T) {
	// Linear: A → B → C
	graph := createTestGraph([][]string{
		{"A", "B"},
		{"B", "C"},
	})

	detector := NewCycleDetector(graph)
	cycles := detector.DetectCycles()

	if len(cycles) != 0 {
		t.Errorf("Expected 0 cycles, got %d", len(cycles))
	}
}

func TestDetectCycles_EmptyGraph(t *testing.T) {
	graph := &types.DependencyGraph{
		Nodes:         map[string]*types.PackageNode{},
		Edges:         []*types.DependencyEdge{},
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
	}

	detector := NewCycleDetector(graph)
	cycles := detector.DetectCycles()

	if len(cycles) != 0 {
		t.Errorf("Expected 0 cycles for empty graph, got %d", len(cycles))
	}
}

func TestDetectCycles_SingleNode(t *testing.T) {
	graph := &types.DependencyGraph{
		Nodes: map[string]*types.PackageNode{
			"A": types.NewPackageNode("A", "1.0.0", "packages/a"),
		},
		Edges:         []*types.DependencyEdge{},
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
	}

	detector := NewCycleDetector(graph)
	cycles := detector.DetectCycles()

	if len(cycles) != 0 {
		t.Errorf("Expected 0 cycles for single node without self-loop, got %d", len(cycles))
	}
}

// ========================================
// Self-Loop Tests (AC5)
// ========================================

func TestDetectCycles_SelfLoop(t *testing.T) {
	graph := &types.DependencyGraph{
		Nodes: map[string]*types.PackageNode{
			"A": {
				Name:         "A",
				Version:      "1.0.0",
				Path:         "packages/a",
				Dependencies: []string{"A"}, // Self-loop
			},
		},
		Edges:         []*types.DependencyEdge{},
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
	}

	detector := NewCycleDetector(graph)
	cycles := detector.DetectCycles()

	if len(cycles) != 1 {
		t.Fatalf("Expected 1 cycle (self-loop), got %d", len(cycles))
	}

	cycle := cycles[0]
	if cycle.Depth != 1 {
		t.Errorf("Self-loop depth = %d, want 1", cycle.Depth)
	}
	if cycle.Severity != types.CircularSeverityCritical {
		t.Errorf("Self-loop severity = %q, want %q", cycle.Severity, types.CircularSeverityCritical)
	}
	if cycle.Type != types.CircularTypeDirect {
		t.Errorf("Self-loop type = %q, want %q", cycle.Type, types.CircularTypeDirect)
	}
}

// ========================================
// Direct Cycle Tests (AC1, AC2)
// ========================================

func TestDetectCycles_DirectCycle(t *testing.T) {
	// A ↔ B
	graph := createTestGraph([][]string{
		{"A", "B"},
		{"B", "A"},
	})

	detector := NewCycleDetector(graph)
	cycles := detector.DetectCycles()

	if len(cycles) != 1 {
		t.Fatalf("Expected 1 cycle (A ↔ B), got %d", len(cycles))
	}

	cycle := cycles[0]
	if cycle.Depth != 2 {
		t.Errorf("Direct cycle depth = %d, want 2", cycle.Depth)
	}
	if cycle.Type != types.CircularTypeDirect {
		t.Errorf("Direct cycle type = %q, want %q", cycle.Type, types.CircularTypeDirect)
	}
	if cycle.Severity != types.CircularSeverityWarning {
		t.Errorf("Direct cycle severity = %q, want %q", cycle.Severity, types.CircularSeverityWarning)
	}

	// Verify cycle format: starts and ends with same node
	if len(cycle.Cycle) < 2 || cycle.Cycle[0] != cycle.Cycle[len(cycle.Cycle)-1] {
		t.Errorf("Cycle should start and end with same node, got %v", cycle.Cycle)
	}
}

// ========================================
// Indirect Cycle Tests (AC1, AC2)
// ========================================

func TestDetectCycles_IndirectCycle3Nodes(t *testing.T) {
	// A → B → C → A
	graph := createTestGraph([][]string{
		{"A", "B"},
		{"B", "C"},
		{"C", "A"},
	})

	detector := NewCycleDetector(graph)
	cycles := detector.DetectCycles()

	if len(cycles) != 1 {
		t.Fatalf("Expected 1 cycle (A → B → C → A), got %d", len(cycles))
	}

	cycle := cycles[0]
	if cycle.Depth != 3 {
		t.Errorf("Indirect cycle depth = %d, want 3", cycle.Depth)
	}
	if cycle.Type != types.CircularTypeIndirect {
		t.Errorf("Indirect cycle type = %q, want %q", cycle.Type, types.CircularTypeIndirect)
	}
	if cycle.Severity != types.CircularSeverityInfo {
		t.Errorf("Indirect cycle severity = %q, want %q", cycle.Severity, types.CircularSeverityInfo)
	}
}

func TestDetectCycles_IndirectCycle5Nodes(t *testing.T) {
	// A → B → C → D → E → A
	graph := createTestGraph([][]string{
		{"A", "B"},
		{"B", "C"},
		{"C", "D"},
		{"D", "E"},
		{"E", "A"},
	})

	detector := NewCycleDetector(graph)
	cycles := detector.DetectCycles()

	if len(cycles) != 1 {
		t.Fatalf("Expected 1 cycle, got %d", len(cycles))
	}

	cycle := cycles[0]
	if cycle.Depth != 5 {
		t.Errorf("Cycle depth = %d, want 5", cycle.Depth)
	}
}

// ========================================
// Multiple Cycles Tests (AC4)
// ========================================

func TestDetectCycles_MultipleSeparateCycles(t *testing.T) {
	// Two separate cycles: A ↔ B and C → D → E → C
	graph := createTestGraph([][]string{
		{"A", "B"},
		{"B", "A"},
		{"C", "D"},
		{"D", "E"},
		{"E", "C"},
	})

	detector := NewCycleDetector(graph)
	cycles := detector.DetectCycles()

	if len(cycles) != 2 {
		t.Fatalf("Expected 2 cycles, got %d", len(cycles))
	}

	// Find each cycle
	directCycle := findCycleWithNodes(cycles, []string{"A", "B"})
	indirectCycle := findCycleWithNodes(cycles, []string{"C", "D", "E"})

	if directCycle == nil {
		t.Error("Missing direct cycle A ↔ B")
	}
	if indirectCycle == nil {
		t.Error("Missing indirect cycle C → D → E → C")
	}
}

func TestDetectCycles_DisconnectedWithOneCycle(t *testing.T) {
	// Disconnected: A → B (no cycle), C → D → C (cycle)
	graph := createTestGraph([][]string{
		{"A", "B"},
		{"C", "D"},
		{"D", "C"},
	})

	detector := NewCycleDetector(graph)
	cycles := detector.DetectCycles()

	if len(cycles) != 1 {
		t.Fatalf("Expected 1 cycle, got %d", len(cycles))
	}

	cycle := findCycleWithNodes(cycles, []string{"C", "D"})
	if cycle == nil {
		t.Error("Missing cycle C → D → C")
	}
}

// ========================================
// Canonical Form Tests (AC3)
// ========================================

func TestNormalizeCycle(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "already normalized",
			input:    []string{"A", "B", "C", "A"},
			expected: []string{"A", "B", "C", "A"},
		},
		{
			name:     "needs rotation",
			input:    []string{"C", "A", "B", "C"},
			expected: []string{"A", "B", "C", "A"},
		},
		{
			name:     "self-loop",
			input:    []string{"X", "X"},
			expected: []string{"X", "X"},
		},
		{
			name:     "direct cycle needs rotation",
			input:    []string{"Z", "A", "Z"},
			expected: []string{"A", "Z", "A"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeCycle(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("Length = %d, want %d", len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("normalizeCycle(%v) = %v, want %v", tt.input, result, tt.expected)
					break
				}
			}
		})
	}
}

func TestDetectCycles_NoDuplicates(t *testing.T) {
	// Same cycle should not be reported twice
	// B → C → A → B is the same as A → B → C → A
	graph := createTestGraph([][]string{
		{"A", "B"},
		{"B", "C"},
		{"C", "A"},
	})

	detector := NewCycleDetector(graph)
	cycles := detector.DetectCycles()

	if len(cycles) != 1 {
		t.Errorf("Expected exactly 1 unique cycle, got %d", len(cycles))
	}

	// Verify canonical form starts with smallest
	if len(cycles) > 0 && cycles[0].Cycle[0] != "A" {
		t.Errorf("Canonical form should start with 'A', got %v", cycles[0].Cycle)
	}
}

// ========================================
// Complex Graph Tests
// ========================================

func TestDetectCycles_ComplexGraph(t *testing.T) {
	// Complex graph with multiple paths
	// A → B → C → A (cycle 1)
	// B → D (no cycle)
	// E → F → E (cycle 2)
	graph := createTestGraph([][]string{
		{"A", "B"},
		{"B", "C"},
		{"C", "A"},
		{"B", "D"},
		{"E", "F"},
		{"F", "E"},
	})

	detector := NewCycleDetector(graph)
	cycles := detector.DetectCycles()

	if len(cycles) != 2 {
		t.Fatalf("Expected 2 cycles, got %d", len(cycles))
	}
}

// ========================================
// Sorting Tests
// ========================================

func TestDetectCycles_SortedBySeverity(t *testing.T) {
	// Create graph with cycles of different severities
	graph := &types.DependencyGraph{
		Nodes: map[string]*types.PackageNode{
			"self": {
				Name:         "self",
				Version:      "1.0.0",
				Path:         "packages/self",
				Dependencies: []string{"self"}, // Self-loop = critical
			},
			"A": types.NewPackageNode("A", "1.0.0", "packages/a"),
			"B": types.NewPackageNode("B", "1.0.0", "packages/b"),
			"C": types.NewPackageNode("C", "1.0.0", "packages/c"),
			"D": types.NewPackageNode("D", "1.0.0", "packages/d"),
		},
		Edges:         []*types.DependencyEdge{},
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
	}

	// A ↔ B = direct = warning
	graph.Nodes["A"].Dependencies = []string{"B"}
	graph.Nodes["B"].Dependencies = []string{"A"}

	// C → D → ... → C (longer cycle) = indirect = info
	graph.Nodes["C"].Dependencies = []string{"D"}
	graph.Nodes["D"].Dependencies = []string{"C"}

	detector := NewCycleDetector(graph)
	cycles := detector.DetectCycles()

	// Should be sorted: critical first, then warning, then info
	if len(cycles) < 2 {
		t.Fatalf("Expected at least 2 cycles, got %d", len(cycles))
	}

	// First should be critical (self-loop)
	if cycles[0].Severity != types.CircularSeverityCritical {
		t.Errorf("First cycle should be critical, got %q", cycles[0].Severity)
	}
}

// ========================================
// DevDependencies and PeerDependencies Tests
// ========================================

func TestDetectCycles_DevDependencies(t *testing.T) {
	// Cycle through devDependencies
	graph := &types.DependencyGraph{
		Nodes: map[string]*types.PackageNode{
			"A": {
				Name:            "A",
				Version:         "1.0.0",
				Path:            "packages/a",
				Dependencies:    []string{},
				DevDependencies: []string{"B"},
			},
			"B": {
				Name:            "B",
				Version:         "1.0.0",
				Path:            "packages/b",
				Dependencies:    []string{},
				DevDependencies: []string{"A"},
			},
		},
		Edges:         []*types.DependencyEdge{},
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
	}

	detector := NewCycleDetector(graph)
	cycles := detector.DetectCycles()

	if len(cycles) != 1 {
		t.Fatalf("Expected 1 cycle through devDependencies, got %d", len(cycles))
	}
}

func TestDetectCycles_PeerDependencies(t *testing.T) {
	// Cycle through peerDependencies
	graph := &types.DependencyGraph{
		Nodes: map[string]*types.PackageNode{
			"A": {
				Name:             "A",
				Version:          "1.0.0",
				Path:             "packages/a",
				Dependencies:     []string{},
				PeerDependencies: []string{"B"},
			},
			"B": {
				Name:             "B",
				Version:          "1.0.0",
				Path:             "packages/b",
				Dependencies:     []string{},
				PeerDependencies: []string{"A"},
			},
		},
		Edges:         []*types.DependencyEdge{},
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
	}

	detector := NewCycleDetector(graph)
	cycles := detector.DetectCycles()

	if len(cycles) != 1 {
		t.Fatalf("Expected 1 cycle through peerDependencies, got %d", len(cycles))
	}
}

// ========================================
// hasSelfLoop Tests
// ========================================

func TestHasSelfLoop(t *testing.T) {
	graph := &types.DependencyGraph{
		Nodes: map[string]*types.PackageNode{
			"self": {
				Name:         "self",
				Version:      "1.0.0",
				Path:         "packages/self",
				Dependencies: []string{"self"},
			},
			"normal": types.NewPackageNode("normal", "1.0.0", "packages/normal"),
		},
	}

	detector := NewCycleDetector(graph)

	if !detector.hasSelfLoop("self") {
		t.Error("hasSelfLoop('self') should return true")
	}

	if detector.hasSelfLoop("normal") {
		t.Error("hasSelfLoop('normal') should return false")
	}

	if detector.hasSelfLoop("nonexistent") {
		t.Error("hasSelfLoop('nonexistent') should return false")
	}
}

// ========================================
// buildAdjacencyList Tests
// ========================================

func TestBuildAdjacencyList(t *testing.T) {
	graph := &types.DependencyGraph{
		Nodes: map[string]*types.PackageNode{
			"A": {
				Name:                 "A",
				Version:              "1.0.0",
				Path:                 "packages/a",
				Dependencies:         []string{"B"},
				DevDependencies:      []string{"C"},
				PeerDependencies:     []string{"D"},
				OptionalDependencies: []string{"E"},
			},
		},
	}

	adj := buildAdjacencyList(graph)

	if len(adj["A"]) != 4 {
		t.Errorf("Expected 4 adjacencies for A, got %d", len(adj["A"]))
	}

	// Check all deps are included
	deps := make(map[string]bool)
	for _, d := range adj["A"] {
		deps[d] = true
	}

	for _, expected := range []string{"B", "C", "D", "E"} {
		if !deps[expected] {
			t.Errorf("Missing dependency %s in adjacency list", expected)
		}
	}
}
