// Package analyzer tests for RootCauseAnalyzer.
package analyzer

import (
	"strings"
	"testing"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

func TestNewRootCauseAnalyzer(t *testing.T) {
	graph := createRootCauseTestGraph()
	analyzer := NewRootCauseAnalyzer(graph)

	if analyzer == nil {
		t.Fatal("NewRootCauseAnalyzer() returned nil")
	}
	if analyzer.graph != graph {
		t.Error("NewRootCauseAnalyzer() graph reference mismatch")
	}
}

func TestRootCauseAnalyzer_Analyze_SelfLoop(t *testing.T) {
	// Self-loop: A → A
	graph := &types.DependencyGraph{
		Nodes: map[string]*types.PackageNode{
			"pkg-a": {
				Name:         "pkg-a",
				Dependencies: []string{"pkg-a"}, // Self-reference
			},
		},
	}

	cycle := &types.CircularDependencyInfo{
		Cycle:    []string{"pkg-a", "pkg-a"},
		Type:     types.CircularTypeDirect,
		Severity: types.CircularSeverityCritical,
		Depth:    1,
	}

	analyzer := NewRootCauseAnalyzer(graph)
	result := analyzer.Analyze(cycle)

	if result == nil {
		t.Fatal("Analyze() returned nil for self-loop")
	}
	if result.OriginatingPackage != "pkg-a" {
		t.Errorf("OriginatingPackage = %s, want pkg-a", result.OriginatingPackage)
	}
	if result.Confidence != 100 {
		t.Errorf("Confidence = %d, want 100 for self-loop", result.Confidence)
	}
}

func TestRootCauseAnalyzer_Analyze_DirectCycle(t *testing.T) {
	// Direct cycle: A ↔ B
	graph := &types.DependencyGraph{
		Nodes: map[string]*types.PackageNode{
			"pkg-a": {
				Name:         "pkg-a",
				Dependencies: []string{"pkg-b"},
			},
			"pkg-b": {
				Name:         "pkg-b",
				Dependencies: []string{"pkg-a"},
			},
		},
	}

	cycle := &types.CircularDependencyInfo{
		Cycle:    []string{"pkg-a", "pkg-b", "pkg-a"},
		Type:     types.CircularTypeDirect,
		Severity: types.CircularSeverityWarning,
		Depth:    2,
	}

	analyzer := NewRootCauseAnalyzer(graph)
	result := analyzer.Analyze(cycle)

	if result == nil {
		t.Fatal("Analyze() returned nil for direct cycle")
	}
	// First alphabetically should be root cause for equal scoring
	if result.OriginatingPackage != "pkg-a" {
		t.Errorf("OriginatingPackage = %s, want pkg-a", result.OriginatingPackage)
	}
	if result.Confidence < 70 || result.Confidence > 90 {
		t.Errorf("Confidence = %d, want 70-90 for direct cycle", result.Confidence)
	}
	if len(result.Chain) != 2 {
		t.Errorf("Chain length = %d, want 2", len(result.Chain))
	}
}

func TestRootCauseAnalyzer_Analyze_IndirectCycleWithCore(t *testing.T) {
	// Indirect cycle with "core" package: ui → api → core → ui
	// "core" should be demoted as root cause
	graph := &types.DependencyGraph{
		Nodes: map[string]*types.PackageNode{
			"pkg-ui": {
				Name:         "pkg-ui",
				Dependencies: []string{"pkg-api"},
			},
			"pkg-api": {
				Name:         "pkg-api",
				Dependencies: []string{"pkg-core"},
			},
			"pkg-core": {
				Name:         "pkg-core",
				Dependencies: []string{"pkg-ui"},
			},
		},
	}

	cycle := &types.CircularDependencyInfo{
		Cycle:    []string{"pkg-api", "pkg-core", "pkg-ui", "pkg-api"},
		Type:     types.CircularTypeIndirect,
		Severity: types.CircularSeverityInfo,
		Depth:    3,
	}

	analyzer := NewRootCauseAnalyzer(graph)
	result := analyzer.Analyze(cycle)

	if result == nil {
		t.Fatal("Analyze() returned nil for indirect cycle")
	}
	// pkg-core should be demoted due to "core" name pattern
	if result.OriginatingPackage == "pkg-core" {
		t.Errorf("OriginatingPackage = %s, should not be 'core' package", result.OriginatingPackage)
	}
	if result.Confidence < 50 {
		t.Errorf("Confidence = %d, want >= 50", result.Confidence)
	}
	if len(result.Chain) != 3 {
		t.Errorf("Chain length = %d, want 3", len(result.Chain))
	}
}

func TestRootCauseAnalyzer_Analyze_HighLevelToLowLevel(t *testing.T) {
	// High-level to low-level: app → service → util → app
	// "app" should be identified as root cause (high-level package)
	graph := &types.DependencyGraph{
		Nodes: map[string]*types.PackageNode{
			"app": {
				Name:         "app",
				Dependencies: []string{"service"},
			},
			"service": {
				Name:         "service",
				Dependencies: []string{"util"},
			},
			"util": {
				Name:              "util",
				Dependencies:     []string{"app"},
				DevDependencies:  []string{},
				PeerDependencies: []string{},
			},
		},
	}

	cycle := &types.CircularDependencyInfo{
		Cycle:    []string{"app", "service", "util", "app"},
		Type:     types.CircularTypeIndirect,
		Severity: types.CircularSeverityInfo,
		Depth:    3,
	}

	analyzer := NewRootCauseAnalyzer(graph)
	result := analyzer.Analyze(cycle)

	if result == nil {
		t.Fatal("Analyze() returned nil")
	}
	// "app" should be root cause (high-level), "util" demoted (low-level pattern)
	if result.OriginatingPackage == "util" {
		t.Errorf("OriginatingPackage = %s, should not be 'util' (low-level)", result.OriginatingPackage)
	}
	if result.Confidence < 70 {
		t.Errorf("Confidence = %d, want >= 70 for high-level root cause", result.Confidence)
	}
}

func TestRootCauseAnalyzer_Analyze_NilCycle(t *testing.T) {
	graph := createRootCauseTestGraph()
	analyzer := NewRootCauseAnalyzer(graph)

	result := analyzer.Analyze(nil)
	if result != nil {
		t.Error("Analyze(nil) should return nil")
	}
}

func TestRootCauseAnalyzer_Analyze_EmptyCycle(t *testing.T) {
	graph := createRootCauseTestGraph()
	analyzer := NewRootCauseAnalyzer(graph)

	cycle := &types.CircularDependencyInfo{
		Cycle: []string{},
	}

	result := analyzer.Analyze(cycle)
	if result != nil {
		t.Error("Analyze(empty cycle) should return nil")
	}
}

func TestRootCauseAnalyzer_Analyze_SingleNodeCycle(t *testing.T) {
	graph := createRootCauseTestGraph()
	analyzer := NewRootCauseAnalyzer(graph)

	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"pkg-a"},
	}

	result := analyzer.Analyze(cycle)
	if result != nil {
		t.Error("Analyze(single node without closing) should return nil")
	}
}

func TestRootCauseAnalyzer_BuildDependencyChain(t *testing.T) {
	graph := &types.DependencyGraph{
		Nodes: map[string]*types.PackageNode{
			"pkg-a": {
				Name:            "pkg-a",
				Dependencies:    []string{"pkg-b"},
				DevDependencies: []string{},
			},
			"pkg-b": {
				Name:            "pkg-b",
				Dependencies:    []string{},
				DevDependencies: []string{"pkg-c"},
			},
			"pkg-c": {
				Name:         "pkg-c",
				Dependencies: []string{"pkg-a"},
			},
		},
	}

	analyzer := NewRootCauseAnalyzer(graph)
	chain := analyzer.buildDependencyChain([]string{"pkg-a", "pkg-b", "pkg-c", "pkg-a"})

	if len(chain) != 3 {
		t.Fatalf("Chain length = %d, want 3", len(chain))
	}

	// Check first edge: pkg-a → pkg-b (production)
	if chain[0].From != "pkg-a" || chain[0].To != "pkg-b" {
		t.Errorf("Edge 0 = %s→%s, want pkg-a→pkg-b", chain[0].From, chain[0].To)
	}
	if chain[0].Type != types.DependencyTypeProduction {
		t.Errorf("Edge 0 type = %s, want production", chain[0].Type)
	}

	// Check second edge: pkg-b → pkg-c (dev)
	if chain[1].From != "pkg-b" || chain[1].To != "pkg-c" {
		t.Errorf("Edge 1 = %s→%s, want pkg-b→pkg-c", chain[1].From, chain[1].To)
	}
	if chain[1].Type != types.DependencyTypeDevelopment {
		t.Errorf("Edge 1 type = %s, want development", chain[1].Type)
	}

	// Check third edge: pkg-c → pkg-a (production)
	if chain[2].From != "pkg-c" || chain[2].To != "pkg-a" {
		t.Errorf("Edge 2 = %s→%s, want pkg-c→pkg-a", chain[2].From, chain[2].To)
	}
}

func TestRootCauseAnalyzer_FindCriticalEdge(t *testing.T) {
	graph := &types.DependencyGraph{
		Nodes: map[string]*types.PackageNode{
			"pkg-a": {Name: "pkg-a", Dependencies: []string{"pkg-b"}},
			"pkg-b": {Name: "pkg-b", DevDependencies: []string{"pkg-a"}},
		},
	}

	analyzer := NewRootCauseAnalyzer(graph)
	chain := []types.RootCauseEdge{
		{From: "pkg-a", To: "pkg-b", Type: types.DependencyTypeProduction, Critical: false},
		{From: "pkg-b", To: "pkg-a", Type: types.DependencyTypeDevelopment, Critical: false},
	}

	criticalEdge := analyzer.findCriticalEdge(chain)

	if criticalEdge == nil {
		t.Fatal("findCriticalEdge() returned nil")
	}
	// Dev dependencies are easier to break, so should be critical
	if criticalEdge.Type != types.DependencyTypeDevelopment {
		t.Errorf("Critical edge type = %s, want development (easier to break)", criticalEdge.Type)
	}
	if !criticalEdge.Critical {
		t.Error("Critical edge should have Critical=true")
	}
}

func TestRootCauseAnalyzer_GenerateExplanation(t *testing.T) {
	tests := []struct {
		name         string
		origin       string
		confidence   int
		wantContains []string
	}{
		{
			name:         "high confidence",
			origin:       "pkg-ui",
			confidence:   85,
			wantContains: []string{"pkg-ui", "highly likely"},
		},
		{
			name:         "medium confidence",
			origin:       "pkg-api",
			confidence:   65,
			wantContains: []string{"pkg-api", "likely"},
		},
		{
			name:         "low confidence",
			origin:       "pkg-core",
			confidence:   40,
			wantContains: []string{"pkg-core", "possibly"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			explanation := generateExplanation(tt.origin, nil, tt.confidence)

			for _, want := range tt.wantContains {
				if !strings.Contains(explanation, want) {
					t.Errorf("Explanation missing '%s': %s", want, explanation)
				}
			}
		})
	}
}

func TestRootCauseAnalyzer_Heuristics(t *testing.T) {
	// Test heuristic scoring with controlled graph
	graph := &types.DependencyGraph{
		Nodes: map[string]*types.PackageNode{
			"high-level-app": {
				Name:         "high-level-app",
				Dependencies: []string{"mid-service"}, // 1 outgoing
			},
			"mid-service": {
				Name:         "mid-service",
				Dependencies: []string{"low-level-utils"}, // 1 outgoing
			},
			"low-level-utils": {
				Name:         "low-level-utils",
				Dependencies: []string{"high-level-app"}, // 1 outgoing, closes cycle
			},
		},
	}

	analyzer := NewRootCauseAnalyzer(graph)

	// Test incoming deps score
	// high-level-app has 1 incoming (from low-level-utils)
	incomingScore := analyzer.calculateIncomingDepsScore("high-level-app")
	if incomingScore < 20 { // Should be decent score with only 1 incoming
		t.Errorf("Incoming score for high-level-app = %d, want >= 20", incomingScore)
	}

	// Test name pattern score
	// "utils" should be demoted
	utilsPatternScore := analyzer.calculateNamePatternScore("low-level-utils")
	if utilsPatternScore != 0 {
		t.Errorf("Name pattern score for 'utils' = %d, want 0", utilsPatternScore)
	}

	// "app" should not be demoted
	appPatternScore := analyzer.calculateNamePatternScore("high-level-app")
	if appPatternScore == 0 {
		t.Error("Name pattern score for 'app' should not be 0")
	}
}

func createRootCauseTestGraph() *types.DependencyGraph {
	return &types.DependencyGraph{
		Nodes: map[string]*types.PackageNode{
			"pkg-a": {
				Name:         "pkg-a",
				Dependencies: []string{"pkg-b"},
			},
			"pkg-b": {
				Name:         "pkg-b",
				Dependencies: []string{"pkg-c"},
			},
			"pkg-c": {
				Name:         "pkg-c",
				Dependencies: []string{"pkg-a"},
			},
		},
	}
}
