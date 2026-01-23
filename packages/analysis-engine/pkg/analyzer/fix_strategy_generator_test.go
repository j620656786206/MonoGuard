// Package analyzer provides dependency graph analysis for monorepo workspaces.
// This file contains tests for fix strategy generator for Story 3.3.
package analyzer

import (
	"testing"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// ========================================
// Test Setup Helpers for Fix Strategy
// ========================================

// createFixStrategyTestGraph creates a simple dependency graph for fix strategy testing.
func createFixStrategyTestGraph() *types.DependencyGraph {
	graph := types.NewDependencyGraph("/test", types.WorkspaceTypePnpm)

	// Add nodes
	graph.Nodes["@mono/ui"] = types.NewPackageNode("@mono/ui", "1.0.0", "packages/ui")
	graph.Nodes["@mono/api"] = types.NewPackageNode("@mono/api", "1.0.0", "packages/api")
	graph.Nodes["@mono/core"] = types.NewPackageNode("@mono/core", "1.0.0", "packages/core")

	// Add edges: ui -> api -> core -> ui (circular)
	graph.Edges = []*types.DependencyEdge{
		{From: "@mono/ui", To: "@mono/api", Type: types.DependencyTypeProduction},
		{From: "@mono/api", To: "@mono/core", Type: types.DependencyTypeProduction},
		{From: "@mono/core", To: "@mono/ui", Type: types.DependencyTypeProduction},
	}

	return graph
}

// createFixStrategyTestWorkspace creates a simple workspace for fix strategy testing.
func createFixStrategyTestWorkspace() *types.WorkspaceData {
	return &types.WorkspaceData{
		WorkspaceType: types.WorkspaceTypePnpm,
		RootPath:      "/test",
		Packages: map[string]*types.PackageInfo{
			"@mono/ui":   {Name: "@mono/ui", Version: "1.0.0", Path: "packages/ui"},
			"@mono/api":  {Name: "@mono/api", Version: "1.0.0", Path: "packages/api"},
			"@mono/core": {Name: "@mono/core", Version: "1.0.0", Path: "packages/core"},
		},
	}
}

// ========================================
// Generator Tests
// ========================================

func TestNewFixStrategyGenerator(t *testing.T) {
	graph := createFixStrategyTestGraph()
	workspace := createFixStrategyTestWorkspace()

	generator := NewFixStrategyGenerator(graph, workspace)

	if generator == nil {
		t.Fatal("NewFixStrategyGenerator() returned nil")
	}
	if generator.graph != graph {
		t.Error("generator.graph not set correctly")
	}
	if generator.workspace != workspace {
		t.Error("generator.workspace not set correctly")
	}
}

func TestFixStrategyGenerator_Generate_NilCycle(t *testing.T) {
	graph := createFixStrategyTestGraph()
	workspace := createFixStrategyTestWorkspace()
	generator := NewFixStrategyGenerator(graph, workspace)

	strategies := generator.Generate(nil)

	// Should return empty slice, not nil
	if strategies == nil {
		t.Fatal("Generate(nil) should return empty slice, not nil")
	}
	if len(strategies) != 0 {
		t.Errorf("Generate(nil) returned %d strategies, want 0", len(strategies))
	}
}

func TestFixStrategyGenerator_Generate_EmptyCycle(t *testing.T) {
	graph := createFixStrategyTestGraph()
	workspace := createFixStrategyTestWorkspace()
	generator := NewFixStrategyGenerator(graph, workspace)

	cycle := &types.CircularDependencyInfo{
		Cycle: []string{},
		Depth: 0,
	}

	strategies := generator.Generate(cycle)

	if strategies == nil {
		t.Fatal("Generate() should return empty slice, not nil")
	}
	if len(strategies) != 0 {
		t.Errorf("Generate() returned %d strategies for empty cycle, want 0", len(strategies))
	}
}

func TestFixStrategyGenerator_Generate_ReturnsThreeStrategies(t *testing.T) {
	graph := createFixStrategyTestGraph()
	workspace := createFixStrategyTestWorkspace()
	generator := NewFixStrategyGenerator(graph, workspace)

	// Create a 3-package cycle: ui -> api -> core -> ui
	cycle := types.NewCircularDependencyInfo([]string{
		"@mono/ui", "@mono/api", "@mono/core", "@mono/ui",
	})
	// Add root cause analysis (from Story 3.1)
	cycle.RootCause = &types.RootCauseAnalysis{
		OriginatingPackage: "@mono/ui",
		CriticalEdge: &types.RootCauseEdge{
			From:     "@mono/core",
			To:       "@mono/ui",
			Type:     types.DependencyTypeProduction,
			Critical: true,
		},
		Confidence: 75,
	}

	strategies := generator.Generate(cycle)

	// Should return up to 3 strategies (AC1)
	if len(strategies) == 0 {
		t.Fatal("Generate() should return at least one strategy")
	}
	if len(strategies) > 3 {
		t.Errorf("Generate() returned %d strategies, max should be 3", len(strategies))
	}

	// Verify all strategy types are present
	typesSeen := make(map[types.FixStrategyType]bool)
	for _, s := range strategies {
		typesSeen[s.Type] = true
	}

	expectedTypes := []types.FixStrategyType{
		types.FixStrategyExtractModule,
		types.FixStrategyDependencyInject,
		types.FixStrategyBoundaryRefactor,
	}
	for _, expected := range expectedTypes {
		if !typesSeen[expected] {
			t.Errorf("Strategy type %s not found in results", expected)
		}
	}
}

func TestFixStrategyGenerator_Generate_SuitabilityScores(t *testing.T) {
	graph := createFixStrategyTestGraph()
	workspace := createFixStrategyTestWorkspace()
	generator := NewFixStrategyGenerator(graph, workspace)

	cycle := types.NewCircularDependencyInfo([]string{
		"@mono/ui", "@mono/api", "@mono/core", "@mono/ui",
	})
	cycle.RootCause = &types.RootCauseAnalysis{
		OriginatingPackage: "@mono/ui",
		Confidence:         80,
	}

	strategies := generator.Generate(cycle)

	// All suitability scores should be 1-10 (AC2)
	for _, s := range strategies {
		if s.Suitability < 1 || s.Suitability > 10 {
			t.Errorf("Strategy %s has invalid suitability %d, want 1-10",
				s.Type, s.Suitability)
		}
	}
}

func TestFixStrategyGenerator_Generate_EffortLevels(t *testing.T) {
	graph := createFixStrategyTestGraph()
	workspace := createFixStrategyTestWorkspace()
	generator := NewFixStrategyGenerator(graph, workspace)

	cycle := types.NewCircularDependencyInfo([]string{
		"@mono/ui", "@mono/api", "@mono/core", "@mono/ui",
	})

	strategies := generator.Generate(cycle)

	validEfforts := map[types.EffortLevel]bool{
		types.EffortLow:    true,
		types.EffortMedium: true,
		types.EffortHigh:   true,
	}

	// All effort levels should be valid (AC3)
	for _, s := range strategies {
		if !validEfforts[s.Effort] {
			t.Errorf("Strategy %s has invalid effort %s", s.Type, s.Effort)
		}
	}
}

func TestFixStrategyGenerator_Generate_ProsCons(t *testing.T) {
	graph := createFixStrategyTestGraph()
	workspace := createFixStrategyTestWorkspace()
	generator := NewFixStrategyGenerator(graph, workspace)

	cycle := types.NewCircularDependencyInfo([]string{
		"@mono/ui", "@mono/api", "@mono/core", "@mono/ui",
	})

	strategies := generator.Generate(cycle)

	// Each strategy should have pros and cons (AC4)
	for _, s := range strategies {
		if len(s.Pros) < 2 {
			t.Errorf("Strategy %s has %d pros, want at least 2",
				s.Type, len(s.Pros))
		}
		if len(s.Cons) < 1 {
			t.Errorf("Strategy %s has %d cons, want at least 1",
				s.Type, len(s.Cons))
		}
	}
}

func TestFixStrategyGenerator_Generate_StrategyRanking(t *testing.T) {
	graph := createFixStrategyTestGraph()
	workspace := createFixStrategyTestWorkspace()
	generator := NewFixStrategyGenerator(graph, workspace)

	cycle := types.NewCircularDependencyInfo([]string{
		"@mono/ui", "@mono/api", "@mono/core", "@mono/ui",
	})

	strategies := generator.Generate(cycle)

	if len(strategies) < 2 {
		t.Skip("Need at least 2 strategies to test ranking")
	}

	// Strategies should be sorted by suitability descending (AC5)
	for i := 1; i < len(strategies); i++ {
		if strategies[i].Suitability > strategies[i-1].Suitability {
			t.Errorf("Strategies not sorted by suitability: %d > %d at indices %d, %d",
				strategies[i].Suitability, strategies[i-1].Suitability, i, i-1)
		}
	}

	// First strategy should be marked as recommended (AC5)
	if !strategies[0].Recommended {
		t.Error("First (top) strategy should be marked as recommended")
	}

	// Other strategies should not be marked recommended
	for i := 1; i < len(strategies); i++ {
		if strategies[i].Recommended {
			t.Errorf("Strategy at index %d should not be marked as recommended", i)
		}
	}
}

func TestFixStrategyGenerator_Generate_TargetPackages(t *testing.T) {
	graph := createFixStrategyTestGraph()
	workspace := createFixStrategyTestWorkspace()
	generator := NewFixStrategyGenerator(graph, workspace)

	cycle := types.NewCircularDependencyInfo([]string{
		"@mono/ui", "@mono/api", "@mono/core", "@mono/ui",
	})

	strategies := generator.Generate(cycle)

	// Each strategy should have target packages
	for _, s := range strategies {
		if len(s.TargetPackages) == 0 {
			t.Errorf("Strategy %s has no target packages", s.Type)
		}
	}
}

// ========================================
// Direct Cycle Tests (A ↔ B)
// ========================================

func TestFixStrategyGenerator_DirectCycle_PrefersDI(t *testing.T) {
	graph := types.NewDependencyGraph("/test", types.WorkspaceTypePnpm)
	graph.Nodes["@mono/a"] = types.NewPackageNode("@mono/a", "1.0.0", "packages/a")
	graph.Nodes["@mono/b"] = types.NewPackageNode("@mono/b", "1.0.0", "packages/b")
	graph.Edges = []*types.DependencyEdge{
		{From: "@mono/a", To: "@mono/b", Type: types.DependencyTypeProduction},
		{From: "@mono/b", To: "@mono/a", Type: types.DependencyTypeProduction},
	}

	workspace := &types.WorkspaceData{
		WorkspaceType: types.WorkspaceTypePnpm,
		RootPath:      "/test",
		Packages: map[string]*types.PackageInfo{
			"@mono/a": {Name: "@mono/a", Version: "1.0.0"},
			"@mono/b": {Name: "@mono/b", Version: "1.0.0"},
		},
	}

	generator := NewFixStrategyGenerator(graph, workspace)

	// Direct cycle: A -> B -> A
	cycle := types.NewCircularDependencyInfo([]string{"@mono/a", "@mono/b", "@mono/a"})
	cycle.RootCause = &types.RootCauseAnalysis{
		OriginatingPackage: "@mono/a",
		CriticalEdge: &types.RootCauseEdge{
			From:     "@mono/b",
			To:       "@mono/a",
			Type:     types.DependencyTypeProduction,
			Critical: true,
		},
	}

	strategies := generator.Generate(cycle)

	if len(strategies) == 0 {
		t.Fatal("No strategies generated for direct cycle")
	}

	// For direct cycles, DI should score high
	var diStrategy *types.FixStrategy
	for i := range strategies {
		if strategies[i].Type == types.FixStrategyDependencyInject {
			diStrategy = &strategies[i]
			break
		}
	}

	if diStrategy == nil {
		t.Fatal("Dependency Injection strategy not found")
	}

	// DI suitability should be high (8-10) for direct cycles
	if diStrategy.Suitability < 8 {
		t.Errorf("DI suitability for direct cycle = %d, want >= 8", diStrategy.Suitability)
	}
}

// ========================================
// Core Package Cycle Tests
// ========================================

func TestFixStrategyGenerator_CorePackageCycle_PrefersBoundaryRefactor(t *testing.T) {
	graph := types.NewDependencyGraph("/test", types.WorkspaceTypePnpm)
	graph.Nodes["@mono/ui"] = types.NewPackageNode("@mono/ui", "1.0.0", "packages/ui")
	graph.Nodes["@mono/core"] = types.NewPackageNode("@mono/core", "1.0.0", "packages/core")
	graph.Edges = []*types.DependencyEdge{
		{From: "@mono/ui", To: "@mono/core", Type: types.DependencyTypeProduction},
		{From: "@mono/core", To: "@mono/ui", Type: types.DependencyTypeProduction},
	}

	workspace := &types.WorkspaceData{
		WorkspaceType: types.WorkspaceTypePnpm,
		RootPath:      "/test",
		Packages: map[string]*types.PackageInfo{
			"@mono/ui":   {Name: "@mono/ui", Version: "1.0.0"},
			"@mono/core": {Name: "@mono/core", Version: "1.0.0"},
		},
	}

	generator := NewFixStrategyGenerator(graph, workspace)

	cycle := types.NewCircularDependencyInfo([]string{"@mono/ui", "@mono/core", "@mono/ui"})

	strategies := generator.Generate(cycle)

	// Find boundary refactor strategy
	var boundaryStrategy *types.FixStrategy
	for i := range strategies {
		if strategies[i].Type == types.FixStrategyBoundaryRefactor {
			boundaryStrategy = &strategies[i]
			break
		}
	}

	if boundaryStrategy == nil {
		t.Fatal("Boundary Refactor strategy not found")
	}

	// Boundary refactoring suitability should be high for core packages
	if boundaryStrategy.Suitability < 7 {
		t.Errorf("Boundary refactor suitability for core package cycle = %d, want >= 7",
			boundaryStrategy.Suitability)
	}
}

// ========================================
// Long Cycle Tests (4+ packages)
// ========================================

func TestFixStrategyGenerator_LongCycle_PrefersExtractModule(t *testing.T) {
	graph := types.NewDependencyGraph("/test", types.WorkspaceTypePnpm)
	packages := []string{"@mono/a", "@mono/b", "@mono/c", "@mono/d"}
	for _, pkg := range packages {
		graph.Nodes[pkg] = types.NewPackageNode(pkg, "1.0.0", "packages/"+pkg[6:])
	}
	// A -> B -> C -> D -> A
	graph.Edges = []*types.DependencyEdge{
		{From: "@mono/a", To: "@mono/b", Type: types.DependencyTypeProduction},
		{From: "@mono/b", To: "@mono/c", Type: types.DependencyTypeProduction},
		{From: "@mono/c", To: "@mono/d", Type: types.DependencyTypeProduction},
		{From: "@mono/d", To: "@mono/a", Type: types.DependencyTypeProduction},
	}

	workspace := &types.WorkspaceData{
		WorkspaceType: types.WorkspaceTypePnpm,
		RootPath:      "/test",
		Packages:      map[string]*types.PackageInfo{},
	}
	for _, pkg := range packages {
		workspace.Packages[pkg] = &types.PackageInfo{Name: pkg, Version: "1.0.0"}
	}

	generator := NewFixStrategyGenerator(graph, workspace)

	cycle := types.NewCircularDependencyInfo([]string{
		"@mono/a", "@mono/b", "@mono/c", "@mono/d", "@mono/a",
	})

	strategies := generator.Generate(cycle)

	// Find extract module strategy
	var extractStrategy *types.FixStrategy
	for i := range strategies {
		if strategies[i].Type == types.FixStrategyExtractModule {
			extractStrategy = &strategies[i]
			break
		}
	}

	if extractStrategy == nil {
		t.Fatal("Extract Module strategy not found")
	}

	// Extract module suitability should be high for long cycles
	if extractStrategy.Suitability < 8 {
		t.Errorf("Extract module suitability for 4-package cycle = %d, want >= 8",
			extractStrategy.Suitability)
	}

	// Extract module should suggest a new package name
	if extractStrategy.NewPackageName == "" {
		t.Error("Extract module should suggest a new package name")
	}
}
