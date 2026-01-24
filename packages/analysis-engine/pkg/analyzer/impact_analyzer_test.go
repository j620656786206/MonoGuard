// Package analyzer provides tests for the impact analyzer (Story 3.6).
package analyzer

import (
	"fmt"
	"testing"
	"time"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// ========================================
// Test Helpers
// ========================================

// createTestGraph creates a dependency graph for testing.
func createImpactTestGraph(nodes []string, edges [][2]string) *types.DependencyGraph {
	graph := types.NewDependencyGraph("/test", types.WorkspaceTypePnpm)

	for _, name := range nodes {
		graph.Nodes[name] = types.NewPackageNode(name, "1.0.0", "/test/"+name)
	}

	for _, edge := range edges {
		graph.Edges = append(graph.Edges, &types.DependencyEdge{
			From:         edge[0],
			To:           edge[1],
			Type:         types.DependencyTypeProduction,
			VersionRange: "^1.0.0",
		})
	}

	return graph
}

// createTestWorkspace creates a workspace for testing.
func createImpactTestWorkspace(packages []string) *types.WorkspaceData {
	ws := &types.WorkspaceData{
		RootPath:      "/test",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages:      make(map[string]*types.PackageInfo),
	}

	for _, pkg := range packages {
		ws.Packages[pkg] = &types.PackageInfo{
			Name:    pkg,
			Version: "1.0.0",
			Path:    "/test/" + pkg,
		}
	}

	return ws
}

// ========================================
// Tests for NewImpactAnalyzer
// ========================================

func TestNewImpactAnalyzer(t *testing.T) {
	graph := createImpactTestGraph([]string{"@mono/a", "@mono/b"}, nil)
	workspace := createImpactTestWorkspace([]string{"@mono/a", "@mono/b"})

	analyzer := NewImpactAnalyzer(graph, workspace)

	if analyzer == nil {
		t.Fatal("NewImpactAnalyzer() returned nil")
	}

	if analyzer.graph != graph {
		t.Error("analyzer.graph should be the provided graph")
	}

	if analyzer.workspace != workspace {
		t.Error("analyzer.workspace should be the provided workspace")
	}

	if analyzer.reverseDeps != nil {
		t.Error("reverseDeps should be nil until Analyze is called")
	}
}

// ========================================
// Tests for buildReverseDependencies
// ========================================

func TestBuildReverseDependencies(t *testing.T) {
	// Graph: A -> B -> C
	//        D -> B
	graph := createImpactTestGraph(
		[]string{"@mono/a", "@mono/b", "@mono/c", "@mono/d"},
		[][2]string{
			{"@mono/a", "@mono/b"}, // A depends on B
			{"@mono/b", "@mono/c"}, // B depends on C
			{"@mono/d", "@mono/b"}, // D depends on B
		},
	)
	workspace := createImpactTestWorkspace([]string{"@mono/a", "@mono/b", "@mono/c", "@mono/d"})

	analyzer := NewImpactAnalyzer(graph, workspace)
	analyzer.buildReverseDependencies()

	tests := []struct {
		pkg        string
		dependents []string
	}{
		{"@mono/b", []string{"@mono/a", "@mono/d"}}, // B is depended on by A and D
		{"@mono/c", []string{"@mono/b"}},            // C is depended on by B
		{"@mono/a", []string{}},                     // A has no dependents
		{"@mono/d", []string{}},                     // D has no dependents
	}

	for _, tt := range tests {
		t.Run(tt.pkg, func(t *testing.T) {
			dependents := analyzer.reverseDeps[tt.pkg]
			if len(dependents) != len(tt.dependents) {
				t.Errorf("reverseDeps[%s] = %v, want %v", tt.pkg, dependents, tt.dependents)
			}
		})
	}
}

// ========================================
// Tests for getDirectParticipants
// ========================================

func TestGetDirectParticipants(t *testing.T) {
	graph := createImpactTestGraph([]string{"@mono/a", "@mono/b", "@mono/c"}, nil)
	workspace := createImpactTestWorkspace([]string{"@mono/a", "@mono/b", "@mono/c"})
	analyzer := NewImpactAnalyzer(graph, workspace)

	tests := []struct {
		name  string
		cycle []string
		want  []string
	}{
		{
			name:  "direct cycle A->B->A",
			cycle: []string{"@mono/a", "@mono/b", "@mono/a"},
			want:  []string{"@mono/a", "@mono/b"},
		},
		{
			name:  "indirect cycle A->B->C->A",
			cycle: []string{"@mono/a", "@mono/b", "@mono/c", "@mono/a"},
			want:  []string{"@mono/a", "@mono/b", "@mono/c"},
		},
		{
			name:  "self-loop A->A",
			cycle: []string{"@mono/a", "@mono/a"},
			want:  []string{"@mono/a"},
		},
		{
			name:  "empty cycle",
			cycle: []string{},
			want:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cycleInfo := &types.CircularDependencyInfo{Cycle: tt.cycle}
			got := analyzer.getDirectParticipants(cycleInfo)

			if len(got) != len(tt.want) {
				t.Errorf("getDirectParticipants() = %v, want %v", got, tt.want)
				return
			}

			for i, pkg := range got {
				if pkg != tt.want[i] {
					t.Errorf("getDirectParticipants()[%d] = %s, want %s", i, pkg, tt.want[i])
				}
			}
		})
	}
}

// ========================================
// Tests for findIndirectDependents
// ========================================

func TestFindIndirectDependents(t *testing.T) {
	// Graph structure:
	// A -> B (A depends on B)
	// B -> C (B depends on C, C is in cycle)
	// D -> C (D depends on C, C is in cycle)
	// E -> A (E depends on A)
	// F -> E (F depends on E) - transitive dependent of cycle via E->A->...
	graph := createImpactTestGraph(
		[]string{"@mono/a", "@mono/b", "@mono/c", "@mono/d", "@mono/e", "@mono/f"},
		[][2]string{
			{"@mono/a", "@mono/b"},
			{"@mono/b", "@mono/c"},
			{"@mono/d", "@mono/c"},
			{"@mono/e", "@mono/a"},
			{"@mono/f", "@mono/e"},
		},
	)
	workspace := createImpactTestWorkspace([]string{"@mono/a", "@mono/b", "@mono/c", "@mono/d", "@mono/e", "@mono/f"})

	analyzer := NewImpactAnalyzer(graph, workspace)
	analyzer.buildReverseDependencies()

	// Cycle: B -> C -> B (direct participants: B, C)
	directParticipants := []string{"@mono/b", "@mono/c"}

	indirectDependents := analyzer.findIndirectDependents(directParticipants)

	// Expected: A (depends on B, distance 1), D (depends on C, distance 1),
	// E (depends on A, distance 2), F (depends on E, distance 3)
	if len(indirectDependents) != 4 {
		t.Errorf("findIndirectDependents() returned %d dependents, want 4", len(indirectDependents))
	}

	// Verify distance 1 dependents
	distance1Count := 0
	for _, dep := range indirectDependents {
		if dep.Distance == 1 {
			distance1Count++
			if dep.PackageName != "@mono/a" && dep.PackageName != "@mono/d" {
				t.Errorf("Unexpected distance 1 dependent: %s", dep.PackageName)
			}
		}
	}
	if distance1Count != 2 {
		t.Errorf("Expected 2 distance-1 dependents, got %d", distance1Count)
	}
}

func TestFindIndirectDependentsNoDependents(t *testing.T) {
	// Isolated cycle with no external dependents
	graph := createImpactTestGraph(
		[]string{"@mono/a", "@mono/b", "@mono/c"},
		[][2]string{
			{"@mono/a", "@mono/b"},
			{"@mono/b", "@mono/a"},
		},
	)
	workspace := createImpactTestWorkspace([]string{"@mono/a", "@mono/b", "@mono/c"})

	analyzer := NewImpactAnalyzer(graph, workspace)
	analyzer.buildReverseDependencies()

	// Cycle: A <-> B
	directParticipants := []string{"@mono/a", "@mono/b"}

	indirectDependents := analyzer.findIndirectDependents(directParticipants)

	// C doesn't depend on A or B, so no indirect dependents
	if len(indirectDependents) != 0 {
		t.Errorf("findIndirectDependents() = %v, want empty", indirectDependents)
	}
}

// ========================================
// Tests for calculateRiskLevel
// ========================================

func TestCalculateRiskLevel(t *testing.T) {
	graph := createImpactTestGraph([]string{"@mono/a", "@mono/b"}, nil)
	workspace := createImpactTestWorkspace([]string{"@mono/a", "@mono/b"})
	analyzer := NewImpactAnalyzer(graph, workspace)

	tests := []struct {
		name        string
		affected    int
		total       int
		participants []string
		wantLevel   types.RiskLevel
	}{
		{
			name:        "critical - over 50%",
			affected:    6,
			total:       10,
			participants: []string{"@mono/ui", "@mono/api"},
			wantLevel:   types.RiskLevelCritical,
		},
		{
			name:        "critical - core package",
			affected:    2,
			total:       10,
			participants: []string{"@mono/core", "@mono/api"},
			wantLevel:   types.RiskLevelCritical,
		},
		{
			name:        "critical - shared package",
			affected:    1,
			total:       10,
			participants: []string{"@mono/shared"},
			wantLevel:   types.RiskLevelCritical,
		},
		{
			name:        "critical - common package",
			affected:    1,
			total:       10,
			participants: []string{"@mono/common"},
			wantLevel:   types.RiskLevelCritical,
		},
		{
			name:        "critical - utils package",
			affected:    1,
			total:       10,
			participants: []string{"@mono/utils"},
			wantLevel:   types.RiskLevelCritical,
		},
		{
			name:        "critical - lib package",
			affected:    1,
			total:       10,
			participants: []string{"@mono/lib"},
			wantLevel:   types.RiskLevelCritical,
		},
		{
			name:        "high - 25-50%",
			affected:    3,
			total:       10,
			participants: []string{"@mono/ui", "@mono/api"},
			wantLevel:   types.RiskLevelHigh,
		},
		{
			name:        "medium - 10-25%",
			affected:    15,
			total:       100,
			participants: []string{"@mono/ui", "@mono/api"},
			wantLevel:   types.RiskLevelMedium,
		},
		{
			name:        "low - under 10%",
			affected:    5,
			total:       100,
			participants: []string{"@mono/ui", "@mono/api"},
			wantLevel:   types.RiskLevelLow,
		},
		{
			name:        "low - empty workspace",
			affected:    0,
			total:       0,
			participants: []string{},
			wantLevel:   types.RiskLevelLow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, explanation := analyzer.calculateRiskLevel(tt.affected, tt.total, tt.participants)

			if level != tt.wantLevel {
				t.Errorf("calculateRiskLevel() level = %s, want %s (explanation: %s)",
					level, tt.wantLevel, explanation)
			}

			if explanation == "" {
				t.Error("calculateRiskLevel() should return non-empty explanation")
			}
		})
	}
}

// ========================================
// Tests for buildRippleEffect
// ========================================

func TestBuildRippleEffect(t *testing.T) {
	graph := createImpactTestGraph([]string{"@mono/a", "@mono/b", "@mono/c"}, nil)
	workspace := createImpactTestWorkspace([]string{"@mono/a", "@mono/b", "@mono/c"})
	analyzer := NewImpactAnalyzer(graph, workspace)

	directParticipants := []string{"@mono/a", "@mono/b"}
	indirectDependents := []types.IndirectDependent{
		{PackageName: "@mono/c", DependsOn: "@mono/a", Distance: 1, DependencyPath: []string{"@mono/a", "@mono/c"}},
		{PackageName: "@mono/d", DependsOn: "@mono/b", Distance: 1, DependencyPath: []string{"@mono/b", "@mono/d"}},
		{PackageName: "@mono/e", DependsOn: "@mono/a", Distance: 2, DependencyPath: []string{"@mono/a", "@mono/c", "@mono/e"}},
	}

	ripple := analyzer.buildRippleEffect(directParticipants, indirectDependents)

	if ripple == nil {
		t.Fatal("buildRippleEffect() returned nil")
	}

	if ripple.TotalLayers != 3 {
		t.Errorf("TotalLayers = %d, want 3", ripple.TotalLayers)
	}

	if len(ripple.Layers) != 3 {
		t.Errorf("len(Layers) = %d, want 3", len(ripple.Layers))
	}

	// Layer 0: direct participants
	if ripple.Layers[0].Distance != 0 {
		t.Errorf("Layer 0 distance = %d, want 0", ripple.Layers[0].Distance)
	}
	if ripple.Layers[0].Count != 2 {
		t.Errorf("Layer 0 count = %d, want 2", ripple.Layers[0].Count)
	}

	// Layer 1: distance 1 dependents
	if ripple.Layers[1].Distance != 1 {
		t.Errorf("Layer 1 distance = %d, want 1", ripple.Layers[1].Distance)
	}
	if ripple.Layers[1].Count != 2 {
		t.Errorf("Layer 1 count = %d, want 2", ripple.Layers[1].Count)
	}

	// Layer 2: distance 2 dependents
	if ripple.Layers[2].Distance != 2 {
		t.Errorf("Layer 2 distance = %d, want 2", ripple.Layers[2].Distance)
	}
	if ripple.Layers[2].Count != 1 {
		t.Errorf("Layer 2 count = %d, want 1", ripple.Layers[2].Count)
	}
}

func TestBuildRippleEffectNoIndirectDependents(t *testing.T) {
	graph := createImpactTestGraph([]string{"@mono/a", "@mono/b"}, nil)
	workspace := createImpactTestWorkspace([]string{"@mono/a", "@mono/b"})
	analyzer := NewImpactAnalyzer(graph, workspace)

	directParticipants := []string{"@mono/a", "@mono/b"}
	indirectDependents := []types.IndirectDependent{}

	ripple := analyzer.buildRippleEffect(directParticipants, indirectDependents)

	if ripple.TotalLayers != 1 {
		t.Errorf("TotalLayers = %d, want 1", ripple.TotalLayers)
	}

	if len(ripple.Layers) != 1 {
		t.Errorf("len(Layers) = %d, want 1", len(ripple.Layers))
	}

	if ripple.Layers[0].Count != 2 {
		t.Errorf("Layer 0 count = %d, want 2", ripple.Layers[0].Count)
	}
}

// ========================================
// Tests for Analyze (Integration)
// ========================================

func TestAnalyze(t *testing.T) {
	// Graph: ui -> api -> core -> ui (cycle)
	//        app -> ui (depends on cycle)
	//        dashboard -> app (transitive dependent)
	graph := createImpactTestGraph(
		[]string{"@mono/ui", "@mono/api", "@mono/core", "@mono/app", "@mono/dashboard"},
		[][2]string{
			{"@mono/ui", "@mono/api"},
			{"@mono/api", "@mono/core"},
			{"@mono/core", "@mono/ui"},
			{"@mono/app", "@mono/ui"},
			{"@mono/dashboard", "@mono/app"},
		},
	)
	workspace := createImpactTestWorkspace([]string{"@mono/ui", "@mono/api", "@mono/core", "@mono/app", "@mono/dashboard"})

	analyzer := NewImpactAnalyzer(graph, workspace)

	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/api", "@mono/core", "@mono/ui"},
	}

	result := analyzer.Analyze(cycle)

	// Verify direct participants
	if len(result.DirectParticipants) != 3 {
		t.Errorf("DirectParticipants = %d, want 3", len(result.DirectParticipants))
	}

	// Verify indirect dependents (app, dashboard)
	if len(result.IndirectDependents) != 2 {
		t.Errorf("IndirectDependents = %d, want 2", len(result.IndirectDependents))
	}

	// Verify total affected (3 direct + 2 indirect = 5)
	if result.TotalAffected != 5 {
		t.Errorf("TotalAffected = %d, want 5", result.TotalAffected)
	}

	// Verify percentage (5/5 = 100%)
	if result.AffectedPercentage != 1.0 {
		t.Errorf("AffectedPercentage = %f, want 1.0", result.AffectedPercentage)
	}

	if result.AffectedPercentageDisplay != "100%" {
		t.Errorf("AffectedPercentageDisplay = %s, want 100%%", result.AffectedPercentageDisplay)
	}

	// Risk level should be critical (100% affected)
	if result.RiskLevel != types.RiskLevelCritical {
		t.Errorf("RiskLevel = %s, want critical", result.RiskLevel)
	}

	// Verify ripple effect
	if result.RippleEffect == nil {
		t.Fatal("RippleEffect should not be nil")
	}
}

func TestAnalyzeNilCycle(t *testing.T) {
	graph := createImpactTestGraph([]string{"@mono/a"}, nil)
	workspace := createImpactTestWorkspace([]string{"@mono/a"})
	analyzer := NewImpactAnalyzer(graph, workspace)

	result := analyzer.Analyze(nil)

	if result == nil {
		t.Fatal("Analyze(nil) should return empty assessment, not nil")
	}

	if len(result.DirectParticipants) != 0 {
		t.Errorf("DirectParticipants should be empty for nil cycle")
	}
}

func TestAnalyzeEmptyCycle(t *testing.T) {
	graph := createImpactTestGraph([]string{"@mono/a"}, nil)
	workspace := createImpactTestWorkspace([]string{"@mono/a"})
	analyzer := NewImpactAnalyzer(graph, workspace)

	cycle := &types.CircularDependencyInfo{Cycle: []string{}}
	result := analyzer.Analyze(cycle)

	if result == nil {
		t.Fatal("Analyze() should return empty assessment, not nil")
	}

	if len(result.DirectParticipants) != 0 {
		t.Errorf("DirectParticipants should be empty for empty cycle")
	}
}

func TestAnalyzeIsolatedCycle(t *testing.T) {
	// Cycle with no external dependents
	graph := createImpactTestGraph(
		[]string{"@mono/a", "@mono/b", "@mono/c", "@mono/d"},
		[][2]string{
			{"@mono/a", "@mono/b"},
			{"@mono/b", "@mono/a"},
			// c and d are isolated
		},
	)
	workspace := createImpactTestWorkspace([]string{"@mono/a", "@mono/b", "@mono/c", "@mono/d"})

	analyzer := NewImpactAnalyzer(graph, workspace)

	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/a", "@mono/b", "@mono/a"},
	}

	result := analyzer.Analyze(cycle)

	if len(result.DirectParticipants) != 2 {
		t.Errorf("DirectParticipants = %d, want 2", len(result.DirectParticipants))
	}

	if len(result.IndirectDependents) != 0 {
		t.Errorf("IndirectDependents = %d, want 0 (isolated cycle)", len(result.IndirectDependents))
	}

	if result.TotalAffected != 2 {
		t.Errorf("TotalAffected = %d, want 2", result.TotalAffected)
	}

	// 2/4 = 50%, which is exactly at the 50% threshold, so should be high (not critical)
	// Actually > 0.50 is critical, so 50% exact is high
	if result.RiskLevel != types.RiskLevelHigh {
		t.Errorf("RiskLevel = %s, want high (50%% exactly)", result.RiskLevel)
	}
}

// ========================================
// Test Scenarios from Story
// ========================================

func TestScenarioIsolatedCycle(t *testing.T) {
	// Scenario: Isolated cycle
	// Direct: 2, Indirect: 0, Total: 2, Percentage: 20% (10 pkg), Expected: Medium
	graph := createImpactTestGraph(
		[]string{"@mono/a", "@mono/b", "@mono/c", "@mono/d", "@mono/e", "@mono/f", "@mono/g", "@mono/h", "@mono/i", "@mono/j"},
		[][2]string{
			{"@mono/a", "@mono/b"},
			{"@mono/b", "@mono/a"},
		},
	)
	workspace := createImpactTestWorkspace([]string{"@mono/a", "@mono/b", "@mono/c", "@mono/d", "@mono/e", "@mono/f", "@mono/g", "@mono/h", "@mono/i", "@mono/j"})

	analyzer := NewImpactAnalyzer(graph, workspace)
	cycle := &types.CircularDependencyInfo{Cycle: []string{"@mono/a", "@mono/b", "@mono/a"}}

	result := analyzer.Analyze(cycle)

	if result.TotalAffected != 2 {
		t.Errorf("TotalAffected = %d, want 2", result.TotalAffected)
	}

	if result.RiskLevel != types.RiskLevelMedium {
		t.Errorf("RiskLevel = %s, want medium (20%%)", result.RiskLevel)
	}
}

func TestScenarioCorePackageCycle(t *testing.T) {
	// Scenario: Core package cycle
	// Direct: 2 (includes core), Indirect: 5, Total: 7, Percentage: 70% (10 pkg), Expected: Critical (core)
	packages := []string{"@mono/core", "@mono/utils", "@mono/a", "@mono/b", "@mono/c", "@mono/d", "@mono/e", "@mono/f", "@mono/g", "@mono/h"}
	edges := [][2]string{
		{"@mono/core", "@mono/utils"},
		{"@mono/utils", "@mono/core"},
		{"@mono/a", "@mono/core"},
		{"@mono/b", "@mono/core"},
		{"@mono/c", "@mono/a"},
		{"@mono/d", "@mono/b"},
		{"@mono/e", "@mono/c"},
	}

	graph := createImpactTestGraph(packages, edges)
	workspace := createImpactTestWorkspace(packages)

	analyzer := NewImpactAnalyzer(graph, workspace)
	cycle := &types.CircularDependencyInfo{Cycle: []string{"@mono/core", "@mono/utils", "@mono/core"}}

	result := analyzer.Analyze(cycle)

	if result.RiskLevel != types.RiskLevelCritical {
		t.Errorf("RiskLevel = %s, want critical (core package)", result.RiskLevel)
	}

	if result.RiskExplanation != "Critical impact: cycle includes core/shared package" {
		t.Errorf("RiskExplanation = %s, want core/shared explanation", result.RiskExplanation)
	}
}

func TestScenarioLowImpact(t *testing.T) {
	// Scenario: Low impact
	// Direct: 2, Indirect: 1, Total: 3, Percentage: 5% (60 pkg), Expected: Low
	packages := make([]string, 60)
	for i := 0; i < 60; i++ {
		packages[i] = "@mono/" + string(rune('a'+i%26)) + string(rune('0'+i/26))
	}

	// Create cycle between first two packages
	edges := [][2]string{
		{packages[0], packages[1]},
		{packages[1], packages[0]},
		{packages[2], packages[0]}, // One dependent
	}

	graph := createImpactTestGraph(packages, edges)
	workspace := createImpactTestWorkspace(packages)

	analyzer := NewImpactAnalyzer(graph, workspace)
	cycle := &types.CircularDependencyInfo{Cycle: []string{packages[0], packages[1], packages[0]}}

	result := analyzer.Analyze(cycle)

	// 3/60 = 5%
	if result.RiskLevel != types.RiskLevelLow {
		t.Errorf("RiskLevel = %s, want low (5%%)", result.RiskLevel)
	}
}

// ========================================
// Test JSON Serialization
// ========================================

func TestAnalyzeResultJSONSerialization(t *testing.T) {
	graph := createImpactTestGraph(
		[]string{"@mono/ui", "@mono/api"},
		[][2]string{
			{"@mono/ui", "@mono/api"},
			{"@mono/api", "@mono/ui"},
		},
	)
	workspace := createImpactTestWorkspace([]string{"@mono/ui", "@mono/api"})

	analyzer := NewImpactAnalyzer(graph, workspace)
	cycle := &types.CircularDependencyInfo{Cycle: []string{"@mono/ui", "@mono/api", "@mono/ui"}}

	result := analyzer.Analyze(cycle)

	// Verify result can be used in CircularDependencyInfo
	cycle.ImpactAssessment = result

	if cycle.ImpactAssessment == nil {
		t.Error("ImpactAssessment should be set on cycle")
	}

	if cycle.ImpactAssessment.RiskLevel == "" {
		t.Error("ImpactAssessment.RiskLevel should not be empty")
	}
}

// ========================================
// Performance Tests (AC8)
// ========================================

func TestPerformance100PackagesWithCycles(t *testing.T) {
	// Create a workspace with 100 packages
	packages := make([]string, 100)
	for i := 0; i < 100; i++ {
		packages[i] = fmt.Sprintf("@mono/pkg-%03d", i)
	}

	// Create dependency edges forming a chain with some cycles
	edges := [][2]string{}

	// Chain dependencies (0->1->2->...->99)
	for i := 0; i < 99; i++ {
		edges = append(edges, [2]string{packages[i], packages[i+1]})
	}

	// Add 5 cycles at various points
	// Cycle 1: 0 <-> 1
	edges = append(edges, [2]string{packages[1], packages[0]})
	// Cycle 2: 20 <-> 21
	edges = append(edges, [2]string{packages[21], packages[20]})
	// Cycle 3: 40 <-> 41
	edges = append(edges, [2]string{packages[41], packages[40]})
	// Cycle 4: 60 <-> 61
	edges = append(edges, [2]string{packages[61], packages[60]})
	// Cycle 5: 80 -> 81 -> 82 -> 80 (indirect)
	edges = append(edges, [2]string{packages[82], packages[80]})

	graph := createImpactTestGraph(packages, edges)
	workspace := createImpactTestWorkspace(packages)

	// Create cycles to analyze
	cycles := []*types.CircularDependencyInfo{
		{Cycle: []string{packages[0], packages[1], packages[0]}},
		{Cycle: []string{packages[20], packages[21], packages[20]}},
		{Cycle: []string{packages[40], packages[41], packages[40]}},
		{Cycle: []string{packages[60], packages[61], packages[60]}},
		{Cycle: []string{packages[80], packages[81], packages[82], packages[80]}},
	}

	// Measure performance
	start := time.Now()

	analyzer := NewImpactAnalyzer(graph, workspace)
	for _, cycle := range cycles {
		result := analyzer.Analyze(cycle)
		if result == nil {
			t.Fatal("Analyze returned nil")
		}
		if len(result.DirectParticipants) == 0 {
			t.Error("DirectParticipants should not be empty")
		}
	}

	elapsed := time.Since(start)

	// AC8: < 200ms for 100 packages with 5 cycles
	maxDuration := 200 * time.Millisecond
	if elapsed > maxDuration {
		t.Errorf("Performance test failed: took %v, want < %v", elapsed, maxDuration)
	}

	t.Logf("Performance: 100 packages, 5 cycles analyzed in %v", elapsed)
}

func BenchmarkImpactAnalyzer100Packages(b *testing.B) {
	// Create a workspace with 100 packages
	packages := make([]string, 100)
	for i := 0; i < 100; i++ {
		packages[i] = fmt.Sprintf("@mono/pkg-%03d", i)
	}

	edges := [][2]string{}
	for i := 0; i < 99; i++ {
		edges = append(edges, [2]string{packages[i], packages[i+1]})
	}
	// Add cycle
	edges = append(edges, [2]string{packages[50], packages[0]})

	graph := createImpactTestGraph(packages, edges)
	workspace := createImpactTestWorkspace(packages)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{packages[0], packages[50], packages[0]},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzer := NewImpactAnalyzer(graph, workspace)
		analyzer.Analyze(cycle)
	}
}

func BenchmarkImpactAnalyzer500Packages(b *testing.B) {
	// Create a workspace with 500 packages (larger workspace)
	packages := make([]string, 500)
	for i := 0; i < 500; i++ {
		packages[i] = fmt.Sprintf("@mono/pkg-%03d", i)
	}

	edges := [][2]string{}
	// Create complex dependency structure
	for i := 0; i < 499; i++ {
		edges = append(edges, [2]string{packages[i], packages[i+1]})
		// Add some cross-dependencies
		if i%10 == 0 && i+50 < 500 {
			edges = append(edges, [2]string{packages[i], packages[i+50]})
		}
	}
	// Add cycle
	edges = append(edges, [2]string{packages[100], packages[0]})

	graph := createImpactTestGraph(packages, edges)
	workspace := createImpactTestWorkspace(packages)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{packages[0], packages[100], packages[0]},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzer := NewImpactAnalyzer(graph, workspace)
		analyzer.Analyze(cycle)
	}
}
