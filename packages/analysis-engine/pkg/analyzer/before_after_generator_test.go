// Package analyzer contains tests for before/after explanation generator.
package analyzer

import (
	"encoding/json"
	"testing"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

func TestNewBeforeAfterGenerator(t *testing.T) {
	graph := types.NewDependencyGraph("/root", types.WorkspaceTypePnpm)
	workspace := &types.WorkspaceData{
		RootPath:      "/root",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages:      map[string]*types.PackageInfo{},
	}

	generator := NewBeforeAfterGenerator(graph, workspace)

	if generator == nil {
		t.Fatal("Expected NewBeforeAfterGenerator to return non-nil")
	}
	if generator.graph != graph {
		t.Error("Expected graph to be set")
	}
	if generator.workspace != workspace {
		t.Error("Expected workspace to be set")
	}
}

func TestGenerate_NilInputs(t *testing.T) {
	generator := NewBeforeAfterGenerator(nil, nil)

	// Test with nil cycle
	result := generator.Generate(nil, &types.FixStrategy{})
	if result != nil {
		t.Error("Expected nil result for nil cycle")
	}

	// Test with nil strategy
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/a", "@mono/b", "@mono/a"},
	}
	result = generator.Generate(cycle, nil)
	if result != nil {
		t.Error("Expected nil result for nil strategy")
	}
}

func TestGenerateCurrentState(t *testing.T) {
	generator := NewBeforeAfterGenerator(nil, nil)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/api", "@mono/core", "@mono/ui"},
	}

	state := generator.generateCurrentState(cycle)

	// Verify nodes
	if len(state.Nodes) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(state.Nodes))
	}

	// Verify all nodes are marked as in-cycle
	for _, node := range state.Nodes {
		if !node.IsInCycle {
			t.Errorf("Expected node %s to be in cycle", node.ID)
		}
		if node.NodeType != types.NodeTypeCycle {
			t.Errorf("Expected node %s to have type 'cycle', got %s", node.ID, node.NodeType)
		}
	}

	// Verify edges
	if len(state.Edges) != 3 {
		t.Errorf("Expected 3 edges, got %d", len(state.Edges))
	}

	// Verify all edges are marked as in-cycle
	for _, edge := range state.Edges {
		if !edge.IsInCycle {
			t.Errorf("Expected edge %s -> %s to be in cycle", edge.From, edge.To)
		}
		if edge.EdgeType != types.EdgeTypeCycle {
			t.Errorf("Expected edge to have type 'cycle', got %s", edge.EdgeType)
		}
	}

	// Verify highlighted path
	if len(state.HighlightedPath) != 4 {
		t.Errorf("Expected 4 elements in highlighted path, got %d", len(state.HighlightedPath))
	}

	// Verify cycle not resolved
	if state.CycleResolved {
		t.Error("Expected CycleResolved to be false")
	}
}

func TestGenerateProposedState_ExtractModule(t *testing.T) {
	generator := NewBeforeAfterGenerator(nil, nil)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/api", "@mono/ui"},
	}
	strategy := &types.FixStrategy{
		Type:           types.FixStrategyExtractModule,
		TargetPackages: []string{"@mono/ui", "@mono/api"},
		NewPackageName: "@mono/shared",
	}

	state := generator.generateProposedState(cycle, strategy)

	// Verify nodes: 2 affected + 1 new = 3
	if len(state.Nodes) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(state.Nodes))
	}

	// Check for new package node
	hasNewPkg := false
	for _, node := range state.Nodes {
		if node.ID == "@mono/shared" {
			hasNewPkg = true
			if !node.IsNew {
				t.Error("Expected @mono/shared to be marked as new")
			}
			if node.NodeType != types.NodeTypeNew {
				t.Errorf("Expected @mono/shared to have type 'new', got %s", node.NodeType)
			}
		}
	}
	if !hasNewPkg {
		t.Error("Expected @mono/shared node in proposed state")
	}

	// Verify edges point to new package
	if len(state.Edges) != 2 {
		t.Errorf("Expected 2 edges, got %d", len(state.Edges))
	}
	for _, edge := range state.Edges {
		if edge.To != "@mono/shared" {
			t.Errorf("Expected edge to point to @mono/shared, got %s", edge.To)
		}
		if !edge.IsNew {
			t.Error("Expected edge to be marked as new")
		}
		if edge.EdgeType != types.EdgeTypeNew {
			t.Errorf("Expected edge type 'new', got %s", edge.EdgeType)
		}
	}

	// Verify cycle resolved
	if !state.CycleResolved {
		t.Error("Expected CycleResolved to be true")
	}
}

func TestGenerateProposedState_DependencyInjection(t *testing.T) {
	generator := NewBeforeAfterGenerator(nil, nil)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/api", "@mono/ui"},
	}
	strategy := &types.FixStrategy{
		Type:           types.FixStrategyDependencyInject,
		TargetPackages: []string{"@mono/ui", "@mono/api"},
	}

	state := generator.generateProposedState(cycle, strategy)

	// Verify nodes
	if len(state.Nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(state.Nodes))
	}

	// Verify no nodes are in cycle
	for _, node := range state.Nodes {
		if node.IsInCycle {
			t.Errorf("Expected node %s to not be in cycle", node.ID)
		}
	}

	// Verify removed edge
	if len(state.Edges) != 1 {
		t.Errorf("Expected 1 edge, got %d", len(state.Edges))
	}
	if !state.Edges[0].IsRemoved {
		t.Error("Expected edge to be marked as removed")
	}

	// Verify cycle resolved
	if !state.CycleResolved {
		t.Error("Expected CycleResolved to be true")
	}
}

func TestGenerateProposedState_BoundaryRefactor(t *testing.T) {
	generator := NewBeforeAfterGenerator(nil, nil)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/api", "@mono/ui"},
	}
	strategy := &types.FixStrategy{
		Type:           types.FixStrategyBoundaryRefactor,
		TargetPackages: []string{"@mono/ui", "@mono/api"},
	}

	state := generator.generateProposedState(cycle, strategy)

	// Verify nodes
	if len(state.Nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(state.Nodes))
	}

	// Verify edges (one-directional after refactoring)
	if len(state.Edges) != 1 {
		t.Errorf("Expected 1 edge, got %d", len(state.Edges))
	}

	// Verify cycle resolved
	if !state.CycleResolved {
		t.Error("Expected CycleResolved to be true")
	}
}

func TestGeneratePackageJsonDiffs_ExtractModule(t *testing.T) {
	workspace := &types.WorkspaceData{
		RootPath:      "/root",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/ui":  {Name: "@mono/ui", Path: "packages/ui"},
			"@mono/api": {Name: "@mono/api", Path: "packages/api"},
		},
	}
	generator := NewBeforeAfterGenerator(nil, workspace)
	strategy := &types.FixStrategy{
		Type:           types.FixStrategyExtractModule,
		TargetPackages: []string{"@mono/ui", "@mono/api"},
		NewPackageName: "@mono/shared",
	}

	diffs := generator.generatePackageJsonDiffs(strategy)

	if len(diffs) != 2 {
		t.Errorf("Expected 2 diffs, got %d", len(diffs))
	}

	for _, diff := range diffs {
		if len(diff.DependenciesToAdd) != 1 {
			t.Errorf("Expected 1 dependency to add for %s", diff.PackageName)
		}
		if diff.DependenciesToAdd[0].Name != "@mono/shared" {
			t.Errorf("Expected dependency @mono/shared, got %s", diff.DependenciesToAdd[0].Name)
		}
		if diff.DependenciesToAdd[0].Version != "workspace:*" {
			t.Errorf("Expected version workspace:*, got %s", diff.DependenciesToAdd[0].Version)
		}
	}
}

func TestGeneratePackageJsonDiffs_DI(t *testing.T) {
	generator := NewBeforeAfterGenerator(nil, nil)
	strategy := &types.FixStrategy{
		Type:           types.FixStrategyDependencyInject,
		TargetPackages: []string{"@mono/ui"},
	}

	diffs := generator.generatePackageJsonDiffs(strategy)

	if len(diffs) != 1 {
		t.Errorf("Expected 1 diff, got %d", len(diffs))
	}
	if len(diffs[0].DependenciesToAdd) != 0 {
		t.Error("Expected no dependencies to add for DI")
	}
}

func TestGenerateImportDiffs_WithTraces(t *testing.T) {
	generator := NewBeforeAfterGenerator(nil, nil)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/api", "@mono/ui"},
		ImportTraces: []types.ImportTrace{
			{
				FromPackage: "@mono/ui",
				ToPackage:   "@mono/api",
				FilePath:    "packages/ui/src/client.ts",
				LineNumber:  5,
				Statement:   "import { helper } from '@mono/api';",
				ImportType:  types.ImportTypeESMNamed,
				Symbols:     []string{"helper"},
			},
		},
	}
	strategy := &types.FixStrategy{
		Type:           types.FixStrategyExtractModule,
		TargetPackages: []string{"@mono/ui", "@mono/api"},
		NewPackageName: "@mono/shared",
	}

	diffs := generator.generateImportDiffs(cycle, strategy)

	if len(diffs) != 1 {
		t.Errorf("Expected 1 diff, got %d", len(diffs))
	}

	diff := diffs[0]
	if diff.FilePath != "packages/ui/src/client.ts" {
		t.Errorf("Expected file path packages/ui/src/client.ts, got %s", diff.FilePath)
	}
	if diff.LineNumber != 5 {
		t.Errorf("Expected line number 5, got %d", diff.LineNumber)
	}
	if len(diff.ImportsToRemove) != 1 {
		t.Errorf("Expected 1 import to remove, got %d", len(diff.ImportsToRemove))
	}
	if len(diff.ImportsToAdd) != 1 {
		t.Errorf("Expected 1 import to add, got %d", len(diff.ImportsToAdd))
	}
	if diff.ImportsToAdd[0].FromPackage != "@mono/shared" {
		t.Errorf("Expected new import from @mono/shared, got %s", diff.ImportsToAdd[0].FromPackage)
	}
}

func TestGenerateImportDiffs_WithoutTraces(t *testing.T) {
	generator := NewBeforeAfterGenerator(nil, nil)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/api", "@mono/ui"},
		// No ImportTraces
	}
	strategy := &types.FixStrategy{
		Type:           types.FixStrategyExtractModule,
		TargetPackages: []string{"@mono/ui", "@mono/api"},
		NewPackageName: "@mono/shared",
	}

	diffs := generator.generateImportDiffs(cycle, strategy)

	if len(diffs) != 2 {
		t.Errorf("Expected 2 diffs (one per cycle edge), got %d", len(diffs))
	}
}

func TestGenerateExplanation_ExtractModule(t *testing.T) {
	generator := NewBeforeAfterGenerator(nil, nil)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/api", "@mono/ui"},
	}
	strategy := &types.FixStrategy{
		Type:           types.FixStrategyExtractModule,
		TargetPackages: []string{"@mono/ui", "@mono/api"},
		NewPackageName: "@mono/shared",
	}

	explanation := generator.generateExplanation(cycle, strategy)

	if explanation == nil {
		t.Fatal("Expected non-nil explanation")
	}
	if explanation.Summary == "" {
		t.Error("Expected non-empty summary")
	}
	if explanation.WhyItWorks == "" {
		t.Error("Expected non-empty whyItWorks")
	}
	if len(explanation.HighLevelChanges) == 0 {
		t.Error("Expected non-empty highLevelChanges")
	}
	if explanation.Confidence != 0.9 {
		t.Errorf("Expected confidence 0.9 for extract module, got %f", explanation.Confidence)
	}
}

func TestGenerateExplanation_DI(t *testing.T) {
	generator := NewBeforeAfterGenerator(nil, nil)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/api", "@mono/ui"},
	}
	strategy := &types.FixStrategy{
		Type:           types.FixStrategyDependencyInject,
		TargetPackages: []string{"@mono/ui", "@mono/api"},
	}

	explanation := generator.generateExplanation(cycle, strategy)

	if explanation.Confidence != 0.75 {
		t.Errorf("Expected confidence 0.75 for DI, got %f", explanation.Confidence)
	}
}

func TestGenerateExplanation_BoundaryRefactor(t *testing.T) {
	generator := NewBeforeAfterGenerator(nil, nil)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/api", "@mono/ui"},
	}
	strategy := &types.FixStrategy{
		Type:           types.FixStrategyBoundaryRefactor,
		TargetPackages: []string{"@mono/ui", "@mono/api"},
	}

	explanation := generator.generateExplanation(cycle, strategy)

	if explanation.Confidence != 0.7 {
		t.Errorf("Expected confidence 0.7 for boundary refactor, got %f", explanation.Confidence)
	}
}

func TestGenerateWarnings_ExtractModule(t *testing.T) {
	generator := NewBeforeAfterGenerator(nil, nil)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/api", "@mono/ui"},
	}
	strategy := &types.FixStrategy{
		Type:           types.FixStrategyExtractModule,
		TargetPackages: []string{"@mono/ui", "@mono/api"},
		NewPackageName: "@mono/shared",
	}

	warnings := generator.generateWarnings(cycle, strategy)

	// Should have at least the install warning
	if len(warnings) == 0 {
		t.Error("Expected at least one warning")
	}

	hasInstallWarning := false
	for _, w := range warnings {
		if w.Title == "New package requires installation" {
			hasInstallWarning = true
			if w.Severity != types.WarningSeverityInfo {
				t.Errorf("Expected info severity, got %s", w.Severity)
			}
		}
	}
	if !hasInstallWarning {
		t.Error("Expected install warning for extract module")
	}
}

func TestGenerateWarnings_ManyPackages(t *testing.T) {
	generator := NewBeforeAfterGenerator(nil, nil)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/a", "@mono/b", "@mono/c", "@mono/d", "@mono/e", "@mono/a"},
	}
	strategy := &types.FixStrategy{
		Type:           types.FixStrategyExtractModule,
		TargetPackages: []string{"@mono/a", "@mono/b", "@mono/c", "@mono/d"},
		NewPackageName: "@mono/shared",
	}

	warnings := generator.generateWarnings(cycle, strategy)

	hasMultiplePackagesWarning := false
	for _, w := range warnings {
		if w.Title == "Multiple packages affected" {
			hasMultiplePackagesWarning = true
			if w.Severity != types.WarningSeverityWarning {
				t.Errorf("Expected warning severity, got %s", w.Severity)
			}
		}
	}
	if !hasMultiplePackagesWarning {
		t.Error("Expected multiple packages warning when > 3 packages")
	}
}

func TestGenerateWarnings_DI(t *testing.T) {
	generator := NewBeforeAfterGenerator(nil, nil)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/api", "@mono/ui"},
	}
	strategy := &types.FixStrategy{
		Type:           types.FixStrategyDependencyInject,
		TargetPackages: []string{"@mono/ui", "@mono/api"},
	}

	warnings := generator.generateWarnings(cycle, strategy)

	hasAPIWarning := false
	hasWiringWarning := false
	for _, w := range warnings {
		if w.Title == "API signature changes" {
			hasAPIWarning = true
		}
		if w.Title == "Runtime wiring required" {
			hasWiringWarning = true
		}
	}
	if !hasAPIWarning {
		t.Error("Expected API signature warning for DI")
	}
	if !hasWiringWarning {
		t.Error("Expected runtime wiring warning for DI")
	}
}

func TestGenerateWarnings_BoundaryRefactor(t *testing.T) {
	generator := NewBeforeAfterGenerator(nil, nil)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/api", "@mono/ui"},
	}
	strategy := &types.FixStrategy{
		Type:           types.FixStrategyBoundaryRefactor,
		TargetPackages: []string{"@mono/ui", "@mono/api"},
	}

	warnings := generator.generateWarnings(cycle, strategy)

	hasCriticalWarning := false
	for _, w := range warnings {
		if w.Title == "Significant code restructuring" && w.Severity == types.WarningSeverityCritical {
			hasCriticalWarning = true
		}
	}
	if !hasCriticalWarning {
		t.Error("Expected critical warning for boundary refactoring")
	}
}

func TestGenerateWarnings_CorePackage(t *testing.T) {
	generator := NewBeforeAfterGenerator(nil, nil)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/core", "@mono/api", "@mono/core"},
	}
	strategy := &types.FixStrategy{
		Type:           types.FixStrategyExtractModule,
		TargetPackages: []string{"@mono/core", "@mono/api"},
		NewPackageName: "@mono/shared",
	}

	warnings := generator.generateWarnings(cycle, strategy)

	hasCoreWarning := false
	for _, w := range warnings {
		if w.Title == "Core package affected" {
			hasCoreWarning = true
			if w.Severity != types.WarningSeverityCritical {
				t.Errorf("Expected critical severity for core package, got %s", w.Severity)
			}
		}
	}
	if !hasCoreWarning {
		t.Error("Expected core package warning")
	}
}

func TestGenerateWarnings_HighImpact(t *testing.T) {
	generator := NewBeforeAfterGenerator(nil, nil)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/api", "@mono/ui"},
		ImpactAssessment: &types.ImpactAssessment{
			RiskLevel:          types.RiskLevelCritical,
			AffectedPercentage: 0.75,
		},
	}
	strategy := &types.FixStrategy{
		Type:           types.FixStrategyExtractModule,
		TargetPackages: []string{"@mono/ui", "@mono/api"},
		NewPackageName: "@mono/shared",
	}

	warnings := generator.generateWarnings(cycle, strategy)

	hasHighImpactWarning := false
	for _, w := range warnings {
		if w.Title == "High-impact cycle" {
			hasHighImpactWarning = true
		}
	}
	if !hasHighImpactWarning {
		t.Error("Expected high-impact warning when ImpactAssessment is critical")
	}
}

func TestGenerate_FullIntegration(t *testing.T) {
	workspace := &types.WorkspaceData{
		RootPath:      "/root",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/ui":  {Name: "@mono/ui", Path: "packages/ui"},
			"@mono/api": {Name: "@mono/api", Path: "packages/api"},
		},
	}
	graph := types.NewDependencyGraph("/root", types.WorkspaceTypePnpm)
	generator := NewBeforeAfterGenerator(graph, workspace)

	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/api", "@mono/ui"},
		ImportTraces: []types.ImportTrace{
			{
				FromPackage: "@mono/ui",
				ToPackage:   "@mono/api",
				FilePath:    "packages/ui/src/client.ts",
				LineNumber:  5,
				Statement:   "import { helper } from '@mono/api';",
				ImportType:  types.ImportTypeESMNamed,
				Symbols:     []string{"helper"},
			},
		},
	}
	strategy := &types.FixStrategy{
		Type:           types.FixStrategyExtractModule,
		TargetPackages: []string{"@mono/ui", "@mono/api"},
		NewPackageName: "@mono/shared",
	}

	explanation := generator.Generate(cycle, strategy)

	if explanation == nil {
		t.Fatal("Expected non-nil explanation")
	}

	// Verify all fields are populated
	if explanation.CurrentState == nil {
		t.Error("Expected CurrentState to be set")
	}
	if explanation.ProposedState == nil {
		t.Error("Expected ProposedState to be set")
	}
	if len(explanation.PackageJsonDiffs) == 0 {
		t.Error("Expected PackageJsonDiffs to be populated")
	}
	if len(explanation.ImportDiffs) == 0 {
		t.Error("Expected ImportDiffs to be populated")
	}
	if explanation.Explanation == nil {
		t.Error("Expected Explanation to be set")
	}
	if len(explanation.Warnings) == 0 {
		t.Error("Expected Warnings to be populated")
	}

	// Verify JSON serialization works
	jsonData, err := json.Marshal(explanation)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Verify camelCase keys in JSON
	jsonStr := string(jsonData)
	expectedKeys := []string{
		`"currentState"`,
		`"proposedState"`,
		`"packageJsonDiffs"`,
		`"importDiffs"`,
		`"explanation"`,
		`"warnings"`,
	}
	for _, key := range expectedKeys {
		if !containsSubstring(jsonStr, key) {
			t.Errorf("Expected JSON to contain key %s", key)
		}
	}
}

func TestExtractShortName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"@mono/shared", "shared"},
		{"@mono/ui", "ui"},
		{"simple-pkg", "simple-pkg"},
		{"@scope/nested/pkg", "pkg"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := extractShortName(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestFormatPackageList(t *testing.T) {
	tests := []struct {
		input    []string
		expected string
	}{
		{[]string{}, ""},
		{[]string{"@mono/a"}, "@mono/a"},
		{[]string{"@mono/a", "@mono/b"}, "@mono/a and @mono/b"},
		{[]string{"@mono/a", "@mono/b", "@mono/c"}, "@mono/a, @mono/b, and @mono/c"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatPackageList(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestGenerateInterfaceNameForPkg(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"@mono/api", "ApiHandler"},
		{"@mono/ui", "UiHandler"},
		{"simple", "SimpleHandler"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := generateInterfaceNameForPkg(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// Helper function
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
