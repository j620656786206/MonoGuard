// Package analyzer provides dependency graph analysis for monorepo workspaces.
// This file contains tests for fix guide generator (Story 3.4).
package analyzer

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// ========================================
// FixGuideGenerator Tests (Story 3.4)
// ========================================

// TestNewFixGuideGenerator verifies constructor creates valid generator.
func TestNewFixGuideGenerator(t *testing.T) {
	workspace := &types.WorkspaceData{
		RootPath:      "/test/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/ui": {
				Name:    "@mono/ui",
				Version: "1.0.0",
				Path:    "packages/ui",
			},
		},
	}

	generator := NewFixGuideGenerator(workspace)

	if generator == nil {
		t.Fatal("Expected generator to be created")
	}
	if generator.workspace != workspace {
		t.Error("Expected workspace to be set")
	}
	if generator.packageManager != "pnpm" {
		t.Errorf("Expected packageManager to be 'pnpm', got '%s'", generator.packageManager)
	}
}

// TestFixGuideGeneratorGenerate verifies guide generation dispatch.
func TestFixGuideGeneratorGenerate(t *testing.T) {
	workspace := &types.WorkspaceData{
		RootPath:      "/test/workspace",
		WorkspaceType: types.WorkspaceTypeNpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/ui": {
				Name:    "@mono/ui",
				Version: "1.0.0",
				Path:    "packages/ui",
			},
			"@mono/core": {
				Name:    "@mono/core",
				Version: "1.0.0",
				Path:    "packages/core",
			},
		},
	}

	generator := NewFixGuideGenerator(workspace)

	tests := []struct {
		name         string
		strategyType types.FixStrategyType
		wantTitle    string
	}{
		{
			name:         "extract module strategy",
			strategyType: types.FixStrategyExtractModule,
			wantTitle:    "Extract Shared Module",
		},
		{
			name:         "dependency injection strategy",
			strategyType: types.FixStrategyDependencyInject,
			wantTitle:    "Dependency Injection",
		},
		{
			name:         "boundary refactoring strategy",
			strategyType: types.FixStrategyBoundaryRefactor,
			wantTitle:    "Module Boundary Refactoring",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cycle := &types.CircularDependencyInfo{
				Cycle: []string{"@mono/ui", "@mono/core", "@mono/ui"},
				Type:  types.CircularTypeDirect,
				Depth: 2,
			}

			strategy := &types.FixStrategy{
				Type:           tt.strategyType,
				TargetPackages: []string{"@mono/ui", "@mono/core"},
				NewPackageName: "@mono/shared",
			}

			guide := generator.Generate(cycle, strategy)

			if guide == nil {
				t.Fatal("Expected guide to be generated")
			}
			if guide.StrategyType != tt.strategyType {
				t.Errorf("Expected strategyType '%s', got '%s'", tt.strategyType, guide.StrategyType)
			}
			if !strings.Contains(guide.Title, tt.wantTitle) {
				t.Errorf("Expected title to contain '%s', got '%s'", tt.wantTitle, guide.Title)
			}
			if len(guide.Steps) == 0 {
				t.Error("Expected guide to have steps")
			}
			if len(guide.Verification) == 0 {
				t.Error("Expected guide to have verification steps")
			}
		})
	}
}

// TestFixGuideGeneratorNilInputs handles edge cases.
func TestFixGuideGeneratorNilInputs(t *testing.T) {
	workspace := &types.WorkspaceData{
		WorkspaceType: types.WorkspaceTypeNpm,
	}
	generator := NewFixGuideGenerator(workspace)

	// Nil cycle
	guide := generator.Generate(nil, &types.FixStrategy{})
	if guide != nil {
		t.Error("Expected nil guide for nil cycle")
	}

	// Nil strategy
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/core", "@mono/ui"},
	}
	guide = generator.Generate(cycle, nil)
	if guide != nil {
		t.Error("Expected nil guide for nil strategy")
	}
}

// ========================================
// Extract Module Guide Tests (Task 3)
// ========================================

// TestExtractModuleGuideSteps verifies extract module guide has correct steps.
func TestExtractModuleGuideSteps(t *testing.T) {
	workspace := &types.WorkspaceData{
		RootPath:      "/test/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/ui": {
				Name:    "@mono/ui",
				Version: "1.0.0",
				Path:    "packages/ui",
			},
			"@mono/core": {
				Name:    "@mono/core",
				Version: "1.0.0",
				Path:    "packages/core",
			},
		},
	}

	generator := NewFixGuideGenerator(workspace)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/core", "@mono/ui"},
		Type:  types.CircularTypeDirect,
		Depth: 2,
	}

	strategy := &types.FixStrategy{
		Type:           types.FixStrategyExtractModule,
		TargetPackages: []string{"@mono/ui", "@mono/core"},
		NewPackageName: "@mono/shared",
		Effort:         types.EffortMedium,
	}

	guide := generator.Generate(cycle, strategy)

	if guide == nil {
		t.Fatal("Expected guide to be generated")
	}

	// Verify step structure
	stepTitles := make([]string, len(guide.Steps))
	for i, step := range guide.Steps {
		stepTitles[i] = step.Title
		// Verify step numbering
		if step.Number != i+1 {
			t.Errorf("Step %d has wrong number: %d", i+1, step.Number)
		}
		// Verify steps have descriptions
		if step.Description == "" {
			t.Errorf("Step %d missing description", i+1)
		}
	}

	// Should include key steps for extract module
	expectedStepPatterns := []string{
		"create", "package.json", "install",
	}
	foundPatterns := 0
	for _, pattern := range expectedStepPatterns {
		for _, title := range stepTitles {
			if strings.Contains(strings.ToLower(title), pattern) {
				foundPatterns++
				break
			}
		}
	}
	if foundPatterns < 2 {
		t.Errorf("Expected at least 2 key step patterns, found %d in %v", foundPatterns, stepTitles)
	}
}

// TestExtractModuleGuideCommands verifies commands use correct package manager.
func TestExtractModuleGuideCommands(t *testing.T) {
	tests := []struct {
		name            string
		workspaceType   types.WorkspaceType
		expectedInstall string
		expectedBuild   string
	}{
		{
			name:            "pnpm workspace",
			workspaceType:   types.WorkspaceTypePnpm,
			expectedInstall: "pnpm install",
			expectedBuild:   "pnpm run build",
		},
		{
			name:            "yarn workspace",
			workspaceType:   types.WorkspaceTypeYarn,
			expectedInstall: "yarn install",
			expectedBuild:   "yarn build",
		},
		{
			name:            "npm workspace",
			workspaceType:   types.WorkspaceTypeNpm,
			expectedInstall: "npm install",
			expectedBuild:   "npm run build",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspace := &types.WorkspaceData{
				WorkspaceType: tt.workspaceType,
				Packages: map[string]*types.PackageInfo{
					"@mono/ui": {Name: "@mono/ui", Path: "packages/ui"},
				},
			}

			generator := NewFixGuideGenerator(workspace)
			cycle := &types.CircularDependencyInfo{
				Cycle: []string{"@mono/ui", "@mono/core", "@mono/ui"},
				Depth: 2,
			}
			strategy := &types.FixStrategy{
				Type:           types.FixStrategyExtractModule,
				TargetPackages: []string{"@mono/ui", "@mono/core"},
				NewPackageName: "@mono/shared",
			}

			guide := generator.Generate(cycle, strategy)

			// Check install command
			foundInstall := false
			for _, step := range guide.Steps {
				if step.Command != nil && step.Command.Command == tt.expectedInstall {
					foundInstall = true
					break
				}
			}
			if !foundInstall {
				t.Errorf("Expected install command '%s' not found", tt.expectedInstall)
			}

			// Check build command in verification
			foundBuild := false
			for _, step := range guide.Verification {
				if step.Command != nil && strings.Contains(step.Command.Command, tt.expectedBuild) {
					foundBuild = true
					break
				}
			}
			if !foundBuild {
				t.Errorf("Expected build command '%s' not found in verification", tt.expectedBuild)
			}
		})
	}
}

// ========================================
// Dependency Injection Guide Tests (Task 4)
// ========================================

// TestDIGuideWithCriticalEdge verifies DI guide uses root cause analysis.
func TestDIGuideWithCriticalEdge(t *testing.T) {
	workspace := &types.WorkspaceData{
		WorkspaceType: types.WorkspaceTypeNpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/ui":   {Name: "@mono/ui", Path: "packages/ui"},
			"@mono/core": {Name: "@mono/core", Path: "packages/core"},
		},
	}

	generator := NewFixGuideGenerator(workspace)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/core", "@mono/ui"},
		Depth: 2,
		RootCause: &types.RootCauseAnalysis{
			CriticalEdge: &types.RootCauseEdge{
				From:     "@mono/ui",
				To:       "@mono/core",
				Type:     types.DependencyTypeProduction,
				Critical: true,
			},
		},
	}

	strategy := &types.FixStrategy{
		Type:           types.FixStrategyDependencyInject,
		TargetPackages: []string{"@mono/ui", "@mono/core"},
	}

	guide := generator.Generate(cycle, strategy)

	if guide == nil {
		t.Fatal("Expected guide to be generated")
	}

	// Should have interface creation step
	foundInterface := false
	for _, step := range guide.Steps {
		if strings.Contains(strings.ToLower(step.Title), "interface") {
			foundInterface = true
			break
		}
	}
	if !foundInterface {
		t.Error("Expected DI guide to include interface creation step")
	}
}

// ========================================
// Boundary Refactoring Guide Tests (Task 5)
// ========================================

// TestBoundaryRefactorGuide verifies boundary refactoring guide structure.
func TestBoundaryRefactorGuide(t *testing.T) {
	workspace := &types.WorkspaceData{
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/ui":   {Name: "@mono/ui", Path: "packages/ui"},
			"@mono/core": {Name: "@mono/core", Path: "packages/core"},
		},
	}

	generator := NewFixGuideGenerator(workspace)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/core", "@mono/ui"},
		Depth: 2,
	}

	strategy := &types.FixStrategy{
		Type:           types.FixStrategyBoundaryRefactor,
		TargetPackages: []string{"@mono/ui", "@mono/core"},
	}

	guide := generator.Generate(cycle, strategy)

	if guide == nil {
		t.Fatal("Expected guide to be generated")
	}

	// Should have responsibility analysis step
	foundResponsibility := false
	for _, step := range guide.Steps {
		if strings.Contains(strings.ToLower(step.Title), "responsibilit") ||
			strings.Contains(strings.ToLower(step.Description), "responsibilit") {
			foundResponsibility = true
			break
		}
	}
	if !foundResponsibility {
		t.Error("Expected boundary guide to include responsibility analysis")
	}
}

// ========================================
// Verification Steps Tests (Task 6)
// ========================================

// TestVerificationSteps verifies verification steps are generated.
func TestVerificationSteps(t *testing.T) {
	workspace := &types.WorkspaceData{
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages:      map[string]*types.PackageInfo{},
	}

	generator := NewFixGuideGenerator(workspace)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/core", "@mono/ui"},
		Depth: 2,
	}

	strategy := &types.FixStrategy{
		Type: types.FixStrategyExtractModule,
	}

	guide := generator.Generate(cycle, strategy)

	if len(guide.Verification) < 2 {
		t.Errorf("Expected at least 2 verification steps, got %d", len(guide.Verification))
	}

	// Verify step numbering
	for i, step := range guide.Verification {
		if step.Number != i+1 {
			t.Errorf("Verification step %d has wrong number: %d", i+1, step.Number)
		}
	}

	// Should include key verification patterns
	expectedPatterns := []string{"monoguard", "build", "test"}
	foundCount := 0
	for _, pattern := range expectedPatterns {
		for _, step := range guide.Verification {
			if step.Command != nil && strings.Contains(strings.ToLower(step.Command.Command), pattern) {
				foundCount++
				break
			}
		}
	}
	if foundCount < 2 {
		t.Errorf("Expected at least 2 verification patterns, found %d", foundCount)
	}
}

// ========================================
// Rollback Instructions Tests (Task 7)
// ========================================

// TestRollbackInstructions verifies rollback instructions are generated.
func TestRollbackInstructions(t *testing.T) {
	workspace := &types.WorkspaceData{
		WorkspaceType: types.WorkspaceTypeNpm,
		Packages:      map[string]*types.PackageInfo{},
	}

	generator := NewFixGuideGenerator(workspace)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/core", "@mono/ui"},
		Depth: 2,
	}

	strategy := &types.FixStrategy{
		Type: types.FixStrategyExtractModule,
	}

	guide := generator.Generate(cycle, strategy)

	if guide.Rollback == nil {
		t.Fatal("Expected rollback instructions to be present")
	}

	// Should have git commands
	if len(guide.Rollback.GitCommands) == 0 {
		t.Error("Expected git commands in rollback")
	}

	// Should have manual steps
	if len(guide.Rollback.ManualSteps) == 0 {
		t.Error("Expected manual steps in rollback")
	}

	// Should have warning
	if guide.Rollback.Warning == "" {
		t.Error("Expected warning in rollback")
	}
}

// ========================================
// Package Manager Detection Tests (Task 8)
// ========================================

// TestDetectPackageManager verifies package manager detection.
func TestDetectPackageManager(t *testing.T) {
	tests := []struct {
		workspaceType types.WorkspaceType
		expected      string
	}{
		{types.WorkspaceTypePnpm, "pnpm"},
		{types.WorkspaceTypeYarn, "yarn"},
		{types.WorkspaceTypeNpm, "npm"},
		{types.WorkspaceTypeUnknown, "npm"},
	}

	for _, tt := range tests {
		t.Run(string(tt.workspaceType), func(t *testing.T) {
			workspace := &types.WorkspaceData{
				WorkspaceType: tt.workspaceType,
			}

			result := detectPackageManager(workspace)

			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// ========================================
// Estimated Time Tests
// ========================================

// TestEstimatedTime verifies time estimation based on effort level.
func TestEstimatedTime(t *testing.T) {
	workspace := &types.WorkspaceData{
		WorkspaceType: types.WorkspaceTypeNpm,
		Packages:      map[string]*types.PackageInfo{},
	}

	generator := NewFixGuideGenerator(workspace)

	tests := []struct {
		effort  types.EffortLevel
		wantMin string
	}{
		{types.EffortLow, "15"},
		{types.EffortMedium, "30"},
		{types.EffortHigh, "60"},
	}

	for _, tt := range tests {
		t.Run(string(tt.effort), func(t *testing.T) {
			cycle := &types.CircularDependencyInfo{
				Cycle: []string{"@mono/ui", "@mono/core", "@mono/ui"},
				Depth: 2,
			}

			strategy := &types.FixStrategy{
				Type:   types.FixStrategyExtractModule,
				Effort: tt.effort,
			}

			guide := generator.Generate(cycle, strategy)

			if !strings.Contains(guide.EstimatedTime, tt.wantMin) {
				t.Errorf("Expected estimated time to contain '%s', got '%s'", tt.wantMin, guide.EstimatedTime)
			}
		})
	}
}

// ========================================
// JSON Serialization Tests
// ========================================

// TestFixGuideJSONOutput verifies guide serializes to proper JSON.
func TestFixGuideJSONOutput(t *testing.T) {
	workspace := &types.WorkspaceData{
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/ui": {Name: "@mono/ui", Path: "packages/ui"},
		},
	}

	generator := NewFixGuideGenerator(workspace)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/core", "@mono/ui"},
		Depth: 2,
	}

	strategy := &types.FixStrategy{
		Type:           types.FixStrategyExtractModule,
		TargetPackages: []string{"@mono/ui", "@mono/core"},
		NewPackageName: "@mono/shared",
		Effort:         types.EffortMedium,
	}

	guide := generator.Generate(cycle, strategy)

	data, err := json.Marshal(guide)
	if err != nil {
		t.Fatalf("Failed to marshal guide: %v", err)
	}

	jsonStr := string(data)

	// Verify camelCase keys
	requiredKeys := []string{
		`"strategyType"`,
		`"title"`,
		`"summary"`,
		`"steps"`,
		`"verification"`,
		`"estimatedTime"`,
	}

	for _, key := range requiredKeys {
		if !strings.Contains(jsonStr, key) {
			t.Errorf("Expected JSON to contain %s", key)
		}
	}

	// Verify no snake_case
	forbiddenKeys := []string{
		`"strategy_type"`,
		`"estimated_time"`,
	}

	for _, key := range forbiddenKeys {
		if strings.Contains(jsonStr, key) {
			t.Errorf("JSON should not contain snake_case key %s", key)
		}
	}
}
