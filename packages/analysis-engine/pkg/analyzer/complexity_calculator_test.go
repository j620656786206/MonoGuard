// Package analyzer provides dependency graph analysis for monorepo workspaces.
// This file contains tests for complexity calculator for Story 3.5.
package analyzer

import (
	"testing"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// TestNewComplexityCalculator verifies constructor.
func TestNewComplexityCalculator(t *testing.T) {
	graph := types.NewDependencyGraph("@mono/root", types.WorkspaceTypeNpm)
	workspace := &types.WorkspaceData{
		RootPath:      "@mono/root",
		WorkspaceType: types.WorkspaceTypeNpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/ui":   {Name: "@mono/ui", Version: "1.0.0"},
			"@mono/api":  {Name: "@mono/api", Version: "1.0.0"},
			"@mono/core": {Name: "@mono/core", Version: "1.0.0"},
		},
	}

	calc := NewComplexityCalculator(graph, workspace)
	if calc == nil {
		t.Fatal("NewComplexityCalculator returned nil")
	}
	if calc.graph != graph {
		t.Error("graph not set correctly")
	}
	if calc.workspace != workspace {
		t.Error("workspace not set correctly")
	}
}

// TestCalculate_SimpleDirect tests a simple direct cycle (A ↔ B).
func TestCalculate_SimpleDirect(t *testing.T) {
	graph := types.NewDependencyGraph("@mono/root", types.WorkspaceTypeNpm)
	workspace := &types.WorkspaceData{
		RootPath:      "@mono/root",
		WorkspaceType: types.WorkspaceTypeNpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/a": {Name: "@mono/a", Version: "1.0.0"},
			"@mono/b": {Name: "@mono/b", Version: "1.0.0"},
		},
	}

	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/a", "@mono/b", "@mono/a"},
		Type:  types.CircularTypeDirect,
		Depth: 2,
	}

	calc := NewComplexityCalculator(graph, workspace)
	result := calc.Calculate(cycle)

	if result == nil {
		t.Fatal("Calculate returned nil")
	}

	// Direct cycle with 2 packages should have low-medium complexity
	if result.Score < 1 || result.Score > 10 {
		t.Errorf("Score = %d, want 1-10", result.Score)
	}

	// Score should be relatively low for simple direct cycle
	if result.Score > 5 {
		t.Errorf("Score = %d, want <= 5 for simple direct cycle", result.Score)
	}

	// Verify estimated time
	if result.EstimatedTime == "" {
		t.Error("EstimatedTime is empty")
	}

	// Verify breakdown exists
	if result.Breakdown.FilesAffected.Value < 0 {
		t.Error("FilesAffected.Value should be >= 0")
	}
	if result.Breakdown.ChainDepth.Value != 2 {
		t.Errorf("ChainDepth.Value = %d, want 2", result.Breakdown.ChainDepth.Value)
	}
	if result.Breakdown.PackagesInvolved.Value != 2 {
		t.Errorf("PackagesInvolved.Value = %d, want 2", result.Breakdown.PackagesInvolved.Value)
	}
}

// TestCalculate_Medium3Package tests a 3-package indirect cycle.
func TestCalculate_Medium3Package(t *testing.T) {
	graph := types.NewDependencyGraph("@mono/root", types.WorkspaceTypeNpm)
	workspace := &types.WorkspaceData{
		RootPath:      "@mono/root",
		WorkspaceType: types.WorkspaceTypeNpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/ui":   {Name: "@mono/ui", Version: "1.0.0"},
			"@mono/api":  {Name: "@mono/api", Version: "1.0.0"},
			"@mono/core": {Name: "@mono/core", Version: "1.0.0"},
		},
	}

	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/api", "@mono/core", "@mono/ui"},
		Type:  types.CircularTypeIndirect,
		Depth: 3,
	}

	calc := NewComplexityCalculator(graph, workspace)
	result := calc.Calculate(cycle)

	if result == nil {
		t.Fatal("Calculate returned nil")
	}

	// 3-package cycle should have medium complexity
	if result.Score < 3 || result.Score > 7 {
		t.Errorf("Score = %d, want 3-7 for medium cycle", result.Score)
	}

	// Verify breakdown
	if result.Breakdown.ChainDepth.Value != 3 {
		t.Errorf("ChainDepth.Value = %d, want 3", result.Breakdown.ChainDepth.Value)
	}
	if result.Breakdown.PackagesInvolved.Value != 3 {
		t.Errorf("PackagesInvolved.Value = %d, want 3", result.Breakdown.PackagesInvolved.Value)
	}
}

// TestCalculate_Complex5Package tests a complex 5-package cycle.
func TestCalculate_Complex5Package(t *testing.T) {
	graph := types.NewDependencyGraph("@mono/root", types.WorkspaceTypeNpm)
	workspace := &types.WorkspaceData{
		RootPath:      "@mono/root",
		WorkspaceType: types.WorkspaceTypeNpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/a": {Name: "@mono/a", Version: "1.0.0"},
			"@mono/b": {Name: "@mono/b", Version: "1.0.0"},
			"@mono/c": {Name: "@mono/c", Version: "1.0.0"},
			"@mono/d": {Name: "@mono/d", Version: "1.0.0"},
			"@mono/e": {Name: "@mono/e", Version: "1.0.0"},
		},
	}

	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/a", "@mono/b", "@mono/c", "@mono/d", "@mono/e", "@mono/a"},
		Type:  types.CircularTypeIndirect,
		Depth: 5,
	}

	calc := NewComplexityCalculator(graph, workspace)
	result := calc.Calculate(cycle)

	if result == nil {
		t.Fatal("Calculate returned nil")
	}

	// 5-package cycle should have higher complexity
	if result.Score < 5 {
		t.Errorf("Score = %d, want >= 5 for complex cycle", result.Score)
	}
}

// TestCalculate_WithImportTraces tests using ImportTraces for accurate file counting.
func TestCalculate_WithImportTraces(t *testing.T) {
	graph := types.NewDependencyGraph("@mono/root", types.WorkspaceTypeNpm)
	workspace := &types.WorkspaceData{
		RootPath:      "@mono/root",
		WorkspaceType: types.WorkspaceTypeNpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/ui":   {Name: "@mono/ui", Version: "1.0.0"},
			"@mono/api":  {Name: "@mono/api", Version: "1.0.0"},
			"@mono/core": {Name: "@mono/core", Version: "1.0.0"},
		},
	}

	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/api", "@mono/core", "@mono/ui"},
		Type:  types.CircularTypeIndirect,
		Depth: 3,
		ImportTraces: []types.ImportTrace{
			{FromPackage: "@mono/ui", ToPackage: "@mono/api", FilePath: "packages/ui/src/client.ts", LineNumber: 5, ImportType: types.ImportTypeESMNamed},
			{FromPackage: "@mono/api", ToPackage: "@mono/core", FilePath: "packages/api/src/service.ts", LineNumber: 10, ImportType: types.ImportTypeESMNamed},
			{FromPackage: "@mono/core", ToPackage: "@mono/ui", FilePath: "packages/core/src/render.ts", LineNumber: 3, ImportType: types.ImportTypeESMDefault},
		},
	}

	calc := NewComplexityCalculator(graph, workspace)
	result := calc.Calculate(cycle)

	if result == nil {
		t.Fatal("Calculate returned nil")
	}

	// Should use ImportTraces for file count
	if result.Breakdown.FilesAffected.Value != 3 {
		t.Errorf("FilesAffected.Value = %d, want 3 (from ImportTraces)", result.Breakdown.FilesAffected.Value)
	}

	if result.Breakdown.ImportsToChange.Value != 3 {
		t.Errorf("ImportsToChange.Value = %d, want 3 (from ImportTraces)", result.Breakdown.ImportsToChange.Value)
	}
}

// TestCalculate_WithExternalDeps tests detection of external dependencies.
func TestCalculate_WithExternalDeps(t *testing.T) {
	graph := types.NewDependencyGraph("@mono/root", types.WorkspaceTypeNpm)
	workspace := &types.WorkspaceData{
		RootPath:      "@mono/root",
		WorkspaceType: types.WorkspaceTypeNpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/ui":  {Name: "@mono/ui", Version: "1.0.0"},
			"@mono/api": {Name: "@mono/api", Version: "1.0.0"},
			// Note: "@external/lib" is NOT in workspace
		},
	}

	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/api", "@mono/ui"},
		Type:  types.CircularTypeDirect,
		Depth: 2,
		RootCause: &types.RootCauseAnalysis{
			Chain: []types.RootCauseEdge{
				{From: "@mono/ui", To: "@mono/api", Type: types.DependencyTypeProduction},
				{From: "@mono/api", To: "@external/lib", Type: types.DependencyTypeProduction},
				{From: "@external/lib", To: "@mono/ui", Type: types.DependencyTypeProduction},
			},
		},
	}

	calc := NewComplexityCalculator(graph, workspace)
	result := calc.Calculate(cycle)

	if result == nil {
		t.Fatal("Calculate returned nil")
	}

	// Should detect external dependency
	if result.Breakdown.ExternalDependencies.Value != 1 {
		t.Errorf("ExternalDependencies.Value = %d, want 1", result.Breakdown.ExternalDependencies.Value)
	}
}

// TestCalculate_NilCycle tests nil input handling.
func TestCalculate_NilCycle(t *testing.T) {
	graph := types.NewDependencyGraph("@mono/root", types.WorkspaceTypeNpm)
	workspace := &types.WorkspaceData{}

	calc := NewComplexityCalculator(graph, workspace)
	result := calc.Calculate(nil)

	if result != nil {
		t.Error("Calculate should return nil for nil cycle")
	}
}

// TestCalculate_EmptyCycle tests empty cycle handling.
func TestCalculate_EmptyCycle(t *testing.T) {
	graph := types.NewDependencyGraph("@mono/root", types.WorkspaceTypeNpm)
	workspace := &types.WorkspaceData{}

	cycle := &types.CircularDependencyInfo{
		Cycle: []string{},
		Depth: 0,
	}

	calc := NewComplexityCalculator(graph, workspace)
	result := calc.Calculate(cycle)

	if result != nil {
		t.Error("Calculate should return nil for empty cycle")
	}
}

// TestEstimateTime tests time estimation for different scores.
func TestEstimateTime(t *testing.T) {
	tests := []struct {
		score int
		want  string
	}{
		{1, "5-15 minutes"},
		{2, "5-15 minutes"},
		{3, "15-30 minutes"},
		{4, "15-30 minutes"},
		{5, "30-60 minutes"},
		{6, "30-60 minutes"},
		{7, "1-2 hours"},
		{8, "1-2 hours"},
		{9, "2-4 hours"},
		{10, "2-4 hours"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := estimateTime(tt.score)
			if got != tt.want {
				t.Errorf("estimateTime(%d) = %q, want %q", tt.score, got, tt.want)
			}
		})
	}
}

// TestGenerateComplexityExplanation tests explanation generation.
func TestGenerateComplexityExplanation(t *testing.T) {
	breakdown := &types.ComplexityBreakdown{
		FilesAffected:   types.ComplexityFactor{Value: 3},
		ImportsToChange: types.ComplexityFactor{Value: 4},
		ChainDepth:      types.ComplexityFactor{Value: 3},
	}

	tests := []struct {
		score       int
		wantContain string
	}{
		{2, "Straightforward"},
		{5, "Moderate"},
		{7, "Significant"},
		{9, "Complex"},
	}

	for _, tt := range tests {
		t.Run(tt.wantContain, func(t *testing.T) {
			got := generateComplexityExplanation(tt.score, breakdown)
			if got == "" {
				t.Error("generateComplexityExplanation returned empty string")
			}
			// The explanation should contain the level word
			if !containsIgnoreCase(got, tt.wantContain) {
				t.Errorf("generateComplexityExplanation(%d) = %q, want to contain %q", tt.score, got, tt.wantContain)
			}
		})
	}
}

// TestCalculate_WeightsSumToOne verifies all weights sum to 1.0.
func TestCalculate_WeightsSumToOne(t *testing.T) {
	graph := types.NewDependencyGraph("@mono/root", types.WorkspaceTypeNpm)
	workspace := &types.WorkspaceData{
		RootPath:      "@mono/root",
		WorkspaceType: types.WorkspaceTypeNpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/a": {Name: "@mono/a", Version: "1.0.0"},
			"@mono/b": {Name: "@mono/b", Version: "1.0.0"},
		},
	}

	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/a", "@mono/b", "@mono/a"},
		Depth: 2,
	}

	calc := NewComplexityCalculator(graph, workspace)
	result := calc.Calculate(cycle)

	if result == nil {
		t.Fatal("Calculate returned nil")
	}

	totalWeight := result.Breakdown.FilesAffected.Weight +
		result.Breakdown.ImportsToChange.Weight +
		result.Breakdown.ChainDepth.Weight +
		result.Breakdown.PackagesInvolved.Weight +
		result.Breakdown.ExternalDependencies.Weight

	if totalWeight != 1.0 {
		t.Errorf("Total weight = %v, want 1.0", totalWeight)
	}
}

// TestCalculate_ScoreBounds verifies score is always 1-10.
func TestCalculate_ScoreBounds(t *testing.T) {
	tests := []struct {
		name  string
		depth int
	}{
		{"depth 1", 1},
		{"depth 2", 2},
		{"depth 5", 5},
		{"depth 10", 10},
		{"depth 20", 20}, // Very large cycle
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graph := types.NewDependencyGraph("@mono/root", types.WorkspaceTypeNpm)
			workspace := &types.WorkspaceData{Packages: map[string]*types.PackageInfo{}}

			cycle := make([]string, tt.depth+1)
			for i := 0; i < tt.depth; i++ {
				cycle[i] = "@mono/pkg" + string(rune('a'+i))
			}
			cycle[tt.depth] = cycle[0] // Close the cycle

			cycleInfo := &types.CircularDependencyInfo{
				Cycle: cycle,
				Depth: tt.depth,
			}

			calc := NewComplexityCalculator(graph, workspace)
			result := calc.Calculate(cycleInfo)

			if result == nil {
				return // Short cycles may return nil
			}

			if result.Score < 1 || result.Score > 10 {
				t.Errorf("Score = %d, want 1-10", result.Score)
			}
		})
	}
}

// helper function for case-insensitive contains
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(substr) == 0 ||
			(len(s) > 0 && containsIgnoreCaseHelper(s, substr)))
}

func containsIgnoreCaseHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			sc := s[i+j]
			subc := substr[j]
			// Simple case-insensitive compare for ASCII
			if sc != subc && sc != subc+32 && sc != subc-32 {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
