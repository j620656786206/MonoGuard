// Package types defines Go types that match TypeScript definitions in @monoguard/types.
// This file contains tests for fix guide types (Story 3.4).
package types

import (
	"encoding/json"
	"testing"
)

// TestFixGuideJSONSerialization verifies JSON output uses camelCase.
func TestFixGuideJSONSerialization(t *testing.T) {
	guide := &FixGuide{
		StrategyType:  FixStrategyExtractModule,
		Title:         "Extract Shared Module: @mono/shared",
		Summary:       "Create a new shared package to hold common dependencies.",
		Steps:         []FixStep{},
		Verification:  []FixStep{},
		EstimatedTime: "30-60 minutes",
	}

	data, err := json.Marshal(guide)
	if err != nil {
		t.Fatalf("Failed to marshal FixGuide: %v", err)
	}

	jsonStr := string(data)

	// Verify camelCase JSON keys
	requiredKeys := []string{
		`"strategyType"`,
		`"title"`,
		`"summary"`,
		`"steps"`,
		`"verification"`,
		`"estimatedTime"`,
	}

	for _, key := range requiredKeys {
		if !fixGuideContains(jsonStr, key) {
			t.Errorf("Expected JSON to contain %s, got: %s", key, jsonStr)
		}
	}

	// Verify no snake_case
	forbiddenKeys := []string{
		`"strategy_type"`,
		`"estimated_time"`,
	}

	for _, key := range forbiddenKeys {
		if fixGuideContains(jsonStr, key) {
			t.Errorf("JSON should not contain snake_case key %s, got: %s", key, jsonStr)
		}
	}
}

// TestFixStepJSONSerialization verifies FixStep JSON output.
func TestFixStepJSONSerialization(t *testing.T) {
	step := FixStep{
		Number:      1,
		Title:       "Create new shared package",
		Description: "Create a new package directory for shared code",
		FilePath:    "packages/shared/package.json",
		CodeBefore: &CodeSnippet{
			Language:  "typescript",
			Code:      "import { helper } from '@mono/core'",
			StartLine: 5,
		},
		CodeAfter: &CodeSnippet{
			Language: "typescript",
			Code:     "import { helper } from '@mono/shared'",
		},
		Command: &CommandStep{
			Command:          "mkdir -p packages/shared/src",
			WorkingDirectory: ".",
			Description:      "Create directory structure",
		},
		ExpectedOutcome: "New package directory exists",
	}

	data, err := json.Marshal(step)
	if err != nil {
		t.Fatalf("Failed to marshal FixStep: %v", err)
	}

	jsonStr := string(data)

	// Verify camelCase JSON keys
	requiredKeys := []string{
		`"number"`,
		`"title"`,
		`"description"`,
		`"filePath"`,
		`"codeBefore"`,
		`"codeAfter"`,
		`"command"`,
		`"expectedOutcome"`,
		`"language"`,
		`"code"`,
		`"startLine"`,
		`"workingDirectory"`,
	}

	for _, key := range requiredKeys {
		if !fixGuideContains(jsonStr, key) {
			t.Errorf("Expected JSON to contain %s, got: %s", key, jsonStr)
		}
	}

	// Verify no snake_case
	forbiddenKeys := []string{
		`"file_path"`,
		`"code_before"`,
		`"code_after"`,
		`"expected_outcome"`,
		`"start_line"`,
		`"working_directory"`,
	}

	for _, key := range forbiddenKeys {
		if fixGuideContains(jsonStr, key) {
			t.Errorf("JSON should not contain snake_case key %s, got: %s", key, jsonStr)
		}
	}
}

// TestRollbackInstructionsJSONSerialization verifies RollbackInstructions JSON output.
func TestRollbackInstructionsJSONSerialization(t *testing.T) {
	rollback := &RollbackInstructions{
		GitCommands: []string{
			"git stash",
			"git checkout .",
		},
		ManualSteps: []string{
			"Restore original imports",
			"Delete new packages",
		},
		Warning: "Coordinate with team before reverting",
	}

	data, err := json.Marshal(rollback)
	if err != nil {
		t.Fatalf("Failed to marshal RollbackInstructions: %v", err)
	}

	jsonStr := string(data)

	// Verify camelCase JSON keys
	requiredKeys := []string{
		`"gitCommands"`,
		`"manualSteps"`,
		`"warning"`,
	}

	for _, key := range requiredKeys {
		if !fixGuideContains(jsonStr, key) {
			t.Errorf("Expected JSON to contain %s, got: %s", key, jsonStr)
		}
	}

	// Verify no snake_case
	forbiddenKeys := []string{
		`"git_commands"`,
		`"manual_steps"`,
	}

	for _, key := range forbiddenKeys {
		if fixGuideContains(jsonStr, key) {
			t.Errorf("JSON should not contain snake_case key %s, got: %s", key, jsonStr)
		}
	}
}

// TestFixGuideOmitempty verifies optional fields are omitted when empty.
func TestFixGuideOmitempty(t *testing.T) {
	guide := &FixGuide{
		StrategyType:  FixStrategyDependencyInject,
		Title:         "Dependency Injection",
		Summary:       "Invert the dependency",
		Steps:         []FixStep{},
		Verification:  []FixStep{},
		EstimatedTime: "15-30 minutes",
		// Rollback is nil - should be omitted
	}

	data, err := json.Marshal(guide)
	if err != nil {
		t.Fatalf("Failed to marshal FixGuide: %v", err)
	}

	jsonStr := string(data)

	// Rollback should be omitted when nil
	if fixGuideContains(jsonStr, `"rollback"`) {
		t.Errorf("Expected rollback to be omitted when nil, got: %s", jsonStr)
	}
}

// TestFixStepOmitempty verifies optional fields in FixStep are omitted.
func TestFixStepOmitempty(t *testing.T) {
	step := FixStep{
		Number:      1,
		Title:       "Simple step",
		Description: "A step without code snippets",
		// No FilePath, CodeBefore, CodeAfter, Command, ExpectedOutcome
	}

	data, err := json.Marshal(step)
	if err != nil {
		t.Fatalf("Failed to marshal FixStep: %v", err)
	}

	jsonStr := string(data)

	// Optional fields should be omitted
	omittedFields := []string{
		`"filePath"`,
		`"codeBefore"`,
		`"codeAfter"`,
		`"command"`,
		`"expectedOutcome"`,
	}

	for _, field := range omittedFields {
		if fixGuideContains(jsonStr, field) {
			t.Errorf("Expected %s to be omitted when empty, got: %s", field, jsonStr)
		}
	}
}

// TestCodeSnippetOmitempty verifies StartLine is omitted when zero.
func TestCodeSnippetOmitempty(t *testing.T) {
	snippet := &CodeSnippet{
		Language: "typescript",
		Code:     "const x = 1",
		// StartLine is 0 - should be omitted
	}

	data, err := json.Marshal(snippet)
	if err != nil {
		t.Fatalf("Failed to marshal CodeSnippet: %v", err)
	}

	jsonStr := string(data)

	if contains(jsonStr, `"startLine"`) {
		t.Errorf("Expected startLine to be omitted when 0, got: %s", jsonStr)
	}
}

// TestCommandStepOmitempty verifies optional fields are omitted.
func TestCommandStepOmitempty(t *testing.T) {
	cmd := &CommandStep{
		Command: "npm install",
		// WorkingDirectory and Description are empty
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		t.Fatalf("Failed to marshal CommandStep: %v", err)
	}

	jsonStr := string(data)

	omittedFields := []string{
		`"workingDirectory"`,
		`"description"`,
	}

	for _, field := range omittedFields {
		if fixGuideContains(jsonStr, field) {
			t.Errorf("Expected %s to be omitted when empty, got: %s", field, jsonStr)
		}
	}
}

// TestFixGuideJSONRoundTrip verifies JSON serialization and deserialization.
func TestFixGuideJSONRoundTrip(t *testing.T) {
	original := &FixGuide{
		StrategyType: FixStrategyExtractModule,
		Title:        "Extract Shared Module",
		Summary:      "Create new shared package",
		Steps: []FixStep{
			{
				Number:      1,
				Title:       "Create package",
				Description: "Create new package directory",
				Command: &CommandStep{
					Command:          "mkdir -p packages/shared",
					WorkingDirectory: ".",
				},
			},
			{
				Number:      2,
				Title:       "Update imports",
				Description: "Update import statements",
				FilePath:    "packages/ui/src/index.ts",
				CodeBefore: &CodeSnippet{
					Language: "typescript",
					Code:     "import { x } from '@mono/core'",
				},
				CodeAfter: &CodeSnippet{
					Language: "typescript",
					Code:     "import { x } from '@mono/shared'",
				},
			},
		},
		Verification: []FixStep{
			{
				Number:      1,
				Title:       "Verify build",
				Description: "Run build",
				Command: &CommandStep{
					Command: "pnpm build",
				},
			},
		},
		Rollback: &RollbackInstructions{
			GitCommands: []string{"git checkout ."},
			Warning:     "Will discard changes",
		},
		EstimatedTime: "30 minutes",
	}

	// Serialize
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Deserialize
	var restored FixGuide
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Verify key fields
	if restored.StrategyType != original.StrategyType {
		t.Errorf("StrategyType mismatch: got %v, want %v", restored.StrategyType, original.StrategyType)
	}
	if restored.Title != original.Title {
		t.Errorf("Title mismatch: got %v, want %v", restored.Title, original.Title)
	}
	if len(restored.Steps) != len(original.Steps) {
		t.Errorf("Steps count mismatch: got %d, want %d", len(restored.Steps), len(original.Steps))
	}
	if len(restored.Verification) != len(original.Verification) {
		t.Errorf("Verification count mismatch: got %d, want %d", len(restored.Verification), len(original.Verification))
	}
	if restored.Rollback == nil {
		t.Error("Rollback should not be nil")
	}
	if restored.EstimatedTime != original.EstimatedTime {
		t.Errorf("EstimatedTime mismatch: got %v, want %v", restored.EstimatedTime, original.EstimatedTime)
	}
}

// Helper function for fix guide tests
func fixGuideContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
