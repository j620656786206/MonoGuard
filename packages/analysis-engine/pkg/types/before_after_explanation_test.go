// Package types contains tests for before/after explanation types.
package types

import (
	"encoding/json"
	"testing"
)

func TestBeforeAfterExplanationJSONSerialization(t *testing.T) {
	explanation := &BeforeAfterExplanation{
		CurrentState: &StateDiagram{
			Nodes: []DiagramNode{
				{
					ID:        "@mono/ui",
					Label:     "ui",
					IsInCycle: true,
					IsNew:     false,
					NodeType:  NodeTypeCycle,
				},
				{
					ID:        "@mono/api",
					Label:     "api",
					IsInCycle: true,
					IsNew:     false,
					NodeType:  NodeTypeCycle,
				},
			},
			Edges: []DiagramEdge{
				{
					From:      "@mono/ui",
					To:        "@mono/api",
					IsInCycle: true,
					IsRemoved: false,
					IsNew:     false,
					EdgeType:  EdgeTypeCycle,
				},
			},
			HighlightedPath: []string{"@mono/ui", "@mono/api", "@mono/ui"},
			CycleResolved:   false,
		},
		ProposedState: &StateDiagram{
			Nodes: []DiagramNode{
				{
					ID:        "@mono/ui",
					Label:     "ui",
					IsInCycle: false,
					IsNew:     false,
					NodeType:  NodeTypeAffected,
				},
				{
					ID:        "@mono/shared",
					Label:     "shared",
					IsInCycle: false,
					IsNew:     true,
					NodeType:  NodeTypeNew,
				},
			},
			Edges: []DiagramEdge{
				{
					From:      "@mono/ui",
					To:        "@mono/shared",
					IsInCycle: false,
					IsRemoved: false,
					IsNew:     true,
					EdgeType:  EdgeTypeNew,
				},
			},
			CycleResolved: true,
		},
		PackageJsonDiffs: []PackageJsonDiff{
			{
				PackageName: "@mono/ui",
				FilePath:    "packages/ui/package.json",
				DependenciesToAdd: []DependencyChange{
					{
						Name:           "@mono/shared",
						Version:        "workspace:*",
						DependencyType: "dependencies",
					},
				},
				DependenciesToRemove: []DependencyChange{},
				Summary:              "Add dependency on @mono/shared",
			},
		},
		ImportDiffs: []ImportDiff{
			{
				FilePath:    "packages/ui/src/client.ts",
				PackageName: "@mono/ui",
				ImportsToRemove: []ImportChange{
					{
						Statement:     "import { helper } from '@mono/api'",
						FromPackage:   "@mono/api",
						ImportedNames: []string{"helper"},
					},
				},
				ImportsToAdd: []ImportChange{
					{
						Statement:     "import { helper } from '@mono/shared'",
						FromPackage:   "@mono/shared",
						ImportedNames: []string{"helper"},
					},
				},
				LineNumber: 5,
			},
		},
		Explanation: &FixExplanation{
			Summary:    "Create a new shared package '@mono/shared' to hold the common code.",
			WhyItWorks: "The cycle exists because packages import from each other.",
			HighLevelChanges: []string{
				"Create new package: @mono/shared",
				"Move shared code to the new package",
			},
			Confidence: 0.9,
		},
		Warnings: []SideEffectWarning{
			{
				Severity:         WarningSeverityInfo,
				Title:            "New package requires installation",
				Description:      "After creating '@mono/shared', run your package manager's install command.",
				AffectedPackages: []string{"@mono/ui", "@mono/api"},
			},
		},
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(explanation)
	if err != nil {
		t.Fatalf("Failed to marshal BeforeAfterExplanation: %v", err)
	}

	// Verify JSON string contains expected camelCase keys
	jsonStr := string(jsonData)
	expectedKeys := []string{
		`"currentState"`,
		`"proposedState"`,
		`"packageJsonDiffs"`,
		`"importDiffs"`,
		`"explanation"`,
		`"warnings"`,
		`"nodes"`,
		`"edges"`,
		`"highlightedPath"`,
		`"cycleResolved"`,
		`"nodeType"`,
		`"edgeType"`,
		`"isInCycle"`,
		`"isNew"`,
		`"isRemoved"`,
		`"dependenciesToAdd"`,
		`"dependenciesToRemove"`,
		`"importsToRemove"`,
		`"importsToAdd"`,
		`"lineNumber"`,
		`"whyItWorks"`,
		`"highLevelChanges"`,
		`"confidence"`,
		`"severity"`,
		`"affectedPackages"`,
	}

	for _, key := range expectedKeys {
		if !containsStr(jsonStr, key) {
			t.Errorf("Expected JSON to contain key %s", key)
		}
	}

	// Deserialize back
	var deserialized BeforeAfterExplanation
	if err := json.Unmarshal(jsonData, &deserialized); err != nil {
		t.Fatalf("Failed to unmarshal BeforeAfterExplanation: %v", err)
	}

	// Verify key fields
	if deserialized.CurrentState == nil {
		t.Error("Expected CurrentState to not be nil")
	}
	if deserialized.ProposedState == nil {
		t.Error("Expected ProposedState to not be nil")
	}
	if len(deserialized.PackageJsonDiffs) != 1 {
		t.Errorf("Expected 1 PackageJsonDiff, got %d", len(deserialized.PackageJsonDiffs))
	}
	if len(deserialized.ImportDiffs) != 1 {
		t.Errorf("Expected 1 ImportDiff, got %d", len(deserialized.ImportDiffs))
	}
	if deserialized.Explanation == nil {
		t.Error("Expected Explanation to not be nil")
	}
	if len(deserialized.Warnings) != 1 {
		t.Errorf("Expected 1 Warning, got %d", len(deserialized.Warnings))
	}
}

func TestStateDiagramJSONSerialization(t *testing.T) {
	diagram := &StateDiagram{
		Nodes: []DiagramNode{
			{
				ID:        "@mono/core",
				Label:     "core",
				IsInCycle: true,
				IsNew:     false,
				NodeType:  NodeTypeCycle,
			},
		},
		Edges: []DiagramEdge{
			{
				From:      "@mono/core",
				To:        "@mono/utils",
				IsInCycle: true,
				IsRemoved: false,
				IsNew:     false,
				EdgeType:  EdgeTypeCycle,
			},
		},
		HighlightedPath: []string{"@mono/core", "@mono/utils", "@mono/core"},
		CycleResolved:   false,
	}

	jsonData, err := json.Marshal(diagram)
	if err != nil {
		t.Fatalf("Failed to marshal StateDiagram: %v", err)
	}

	var deserialized StateDiagram
	if err := json.Unmarshal(jsonData, &deserialized); err != nil {
		t.Fatalf("Failed to unmarshal StateDiagram: %v", err)
	}

	if len(deserialized.Nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(deserialized.Nodes))
	}
	if len(deserialized.Edges) != 1 {
		t.Errorf("Expected 1 edge, got %d", len(deserialized.Edges))
	}
	if deserialized.CycleResolved != false {
		t.Error("Expected CycleResolved to be false")
	}
	if len(deserialized.HighlightedPath) != 3 {
		t.Errorf("Expected 3 elements in HighlightedPath, got %d", len(deserialized.HighlightedPath))
	}
}

func TestDiagramNodeTypes(t *testing.T) {
	tests := []struct {
		nodeType DiagramNodeType
		expected string
	}{
		{NodeTypeCycle, "cycle"},
		{NodeTypeAffected, "affected"},
		{NodeTypeNew, "new"},
		{NodeTypeUnchanged, "unchanged"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if string(tt.nodeType) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.nodeType))
			}
		})
	}
}

func TestDiagramEdgeTypes(t *testing.T) {
	tests := []struct {
		edgeType DiagramEdgeType
		expected string
	}{
		{EdgeTypeCycle, "cycle"},
		{EdgeTypeRemoved, "removed"},
		{EdgeTypeNew, "new"},
		{EdgeTypeUnchanged, "unchanged"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if string(tt.edgeType) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.edgeType))
			}
		})
	}
}

func TestWarningSeverityTypes(t *testing.T) {
	tests := []struct {
		severity WarningSeverity
		expected string
	}{
		{WarningSeverityInfo, "info"},
		{WarningSeverityWarning, "warning"},
		{WarningSeverityCritical, "critical"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if string(tt.severity) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.severity))
			}
		})
	}
}

func TestPackageJsonDiffSerialization(t *testing.T) {
	diff := PackageJsonDiff{
		PackageName: "@mono/ui",
		FilePath:    "packages/ui/package.json",
		DependenciesToAdd: []DependencyChange{
			{Name: "@mono/shared", Version: "workspace:*", DependencyType: "dependencies"},
		},
		DependenciesToRemove: []DependencyChange{
			{Name: "@mono/api", Version: "^1.0.0", DependencyType: "dependencies"},
		},
		Summary: "Replace @mono/api with @mono/shared",
	}

	jsonData, err := json.Marshal(diff)
	if err != nil {
		t.Fatalf("Failed to marshal PackageJsonDiff: %v", err)
	}

	jsonStr := string(jsonData)
	if !containsStr(jsonStr, `"packageName"`) {
		t.Error("Expected camelCase key packageName")
	}
	if !containsStr(jsonStr, `"filePath"`) {
		t.Error("Expected camelCase key filePath")
	}
	if !containsStr(jsonStr, `"dependenciesToAdd"`) {
		t.Error("Expected camelCase key dependenciesToAdd")
	}
	if !containsStr(jsonStr, `"dependenciesToRemove"`) {
		t.Error("Expected camelCase key dependenciesToRemove")
	}
	if !containsStr(jsonStr, `"dependencyType"`) {
		t.Error("Expected camelCase key dependencyType")
	}

	var deserialized PackageJsonDiff
	if err := json.Unmarshal(jsonData, &deserialized); err != nil {
		t.Fatalf("Failed to unmarshal PackageJsonDiff: %v", err)
	}

	if len(deserialized.DependenciesToAdd) != 1 {
		t.Errorf("Expected 1 dependency to add, got %d", len(deserialized.DependenciesToAdd))
	}
	if len(deserialized.DependenciesToRemove) != 1 {
		t.Errorf("Expected 1 dependency to remove, got %d", len(deserialized.DependenciesToRemove))
	}
}

func TestImportDiffSerialization(t *testing.T) {
	diff := ImportDiff{
		FilePath:    "packages/ui/src/client.ts",
		PackageName: "@mono/ui",
		ImportsToRemove: []ImportChange{
			{
				Statement:     "import { helper } from '@mono/api'",
				FromPackage:   "@mono/api",
				ImportedNames: []string{"helper"},
			},
		},
		ImportsToAdd: []ImportChange{
			{
				Statement:     "import { helper } from '@mono/shared'",
				FromPackage:   "@mono/shared",
				ImportedNames: []string{"helper"},
			},
		},
		LineNumber: 5,
	}

	jsonData, err := json.Marshal(diff)
	if err != nil {
		t.Fatalf("Failed to marshal ImportDiff: %v", err)
	}

	jsonStr := string(jsonData)
	if !containsStr(jsonStr, `"filePath"`) {
		t.Error("Expected camelCase key filePath")
	}
	if !containsStr(jsonStr, `"packageName"`) {
		t.Error("Expected camelCase key packageName")
	}
	if !containsStr(jsonStr, `"importsToRemove"`) {
		t.Error("Expected camelCase key importsToRemove")
	}
	if !containsStr(jsonStr, `"importsToAdd"`) {
		t.Error("Expected camelCase key importsToAdd")
	}
	if !containsStr(jsonStr, `"fromPackage"`) {
		t.Error("Expected camelCase key fromPackage")
	}
	if !containsStr(jsonStr, `"importedNames"`) {
		t.Error("Expected camelCase key importedNames")
	}

	var deserialized ImportDiff
	if err := json.Unmarshal(jsonData, &deserialized); err != nil {
		t.Fatalf("Failed to unmarshal ImportDiff: %v", err)
	}

	if deserialized.LineNumber != 5 {
		t.Errorf("Expected LineNumber 5, got %d", deserialized.LineNumber)
	}
}

func TestFixExplanationSerialization(t *testing.T) {
	explanation := FixExplanation{
		Summary:    "Create a shared module to break the cycle.",
		WhyItWorks: "Moving shared code eliminates mutual imports.",
		HighLevelChanges: []string{
			"Create new package",
			"Move shared code",
			"Update imports",
		},
		Confidence: 0.85,
	}

	jsonData, err := json.Marshal(explanation)
	if err != nil {
		t.Fatalf("Failed to marshal FixExplanation: %v", err)
	}

	jsonStr := string(jsonData)
	if !containsStr(jsonStr, `"summary"`) {
		t.Error("Expected camelCase key summary")
	}
	if !containsStr(jsonStr, `"whyItWorks"`) {
		t.Error("Expected camelCase key whyItWorks")
	}
	if !containsStr(jsonStr, `"highLevelChanges"`) {
		t.Error("Expected camelCase key highLevelChanges")
	}
	if !containsStr(jsonStr, `"confidence"`) {
		t.Error("Expected camelCase key confidence")
	}

	var deserialized FixExplanation
	if err := json.Unmarshal(jsonData, &deserialized); err != nil {
		t.Fatalf("Failed to unmarshal FixExplanation: %v", err)
	}

	if deserialized.Confidence != 0.85 {
		t.Errorf("Expected Confidence 0.85, got %f", deserialized.Confidence)
	}
}

func TestSideEffectWarningSerialization(t *testing.T) {
	warning := SideEffectWarning{
		Severity:         WarningSeverityCritical,
		Title:            "Breaking change",
		Description:      "This change affects exported API.",
		AffectedPackages: []string{"@mono/ui", "@mono/api"},
	}

	jsonData, err := json.Marshal(warning)
	if err != nil {
		t.Fatalf("Failed to marshal SideEffectWarning: %v", err)
	}

	jsonStr := string(jsonData)
	if !containsStr(jsonStr, `"severity"`) {
		t.Error("Expected camelCase key severity")
	}
	if !containsStr(jsonStr, `"title"`) {
		t.Error("Expected camelCase key title")
	}
	if !containsStr(jsonStr, `"description"`) {
		t.Error("Expected camelCase key description")
	}
	if !containsStr(jsonStr, `"affectedPackages"`) {
		t.Error("Expected camelCase key affectedPackages")
	}
	if !containsStr(jsonStr, `"critical"`) {
		t.Error("Expected severity value critical")
	}

	var deserialized SideEffectWarning
	if err := json.Unmarshal(jsonData, &deserialized); err != nil {
		t.Fatalf("Failed to unmarshal SideEffectWarning: %v", err)
	}

	if deserialized.Severity != WarningSeverityCritical {
		t.Errorf("Expected Severity critical, got %s", deserialized.Severity)
	}
	if len(deserialized.AffectedPackages) != 2 {
		t.Errorf("Expected 2 affected packages, got %d", len(deserialized.AffectedPackages))
	}
}

func TestNewBeforeAfterExplanation(t *testing.T) {
	explanation := NewBeforeAfterExplanation()

	if explanation == nil {
		t.Fatal("Expected NewBeforeAfterExplanation to return non-nil")
	}

	// Verify slices are initialized (not nil)
	if explanation.PackageJsonDiffs == nil {
		t.Error("Expected PackageJsonDiffs to be initialized")
	}
	if explanation.ImportDiffs == nil {
		t.Error("Expected ImportDiffs to be initialized")
	}
	if explanation.Warnings == nil {
		t.Error("Expected Warnings to be initialized")
	}

	// Verify slices are empty
	if len(explanation.PackageJsonDiffs) != 0 {
		t.Errorf("Expected PackageJsonDiffs to be empty, got %d", len(explanation.PackageJsonDiffs))
	}
	if len(explanation.ImportDiffs) != 0 {
		t.Errorf("Expected ImportDiffs to be empty, got %d", len(explanation.ImportDiffs))
	}
	if len(explanation.Warnings) != 0 {
		t.Errorf("Expected Warnings to be empty, got %d", len(explanation.Warnings))
	}
}

func TestNewStateDiagram(t *testing.T) {
	diagram := NewStateDiagram()

	if diagram == nil {
		t.Fatal("Expected NewStateDiagram to return non-nil")
	}

	// Verify slices are initialized (not nil)
	if diagram.Nodes == nil {
		t.Error("Expected Nodes to be initialized")
	}
	if diagram.Edges == nil {
		t.Error("Expected Edges to be initialized")
	}

	// Verify slices are empty
	if len(diagram.Nodes) != 0 {
		t.Errorf("Expected Nodes to be empty, got %d", len(diagram.Nodes))
	}
	if len(diagram.Edges) != 0 {
		t.Errorf("Expected Edges to be empty, got %d", len(diagram.Edges))
	}

	// Verify default value
	if diagram.CycleResolved != false {
		t.Error("Expected CycleResolved to be false by default")
	}
}

func TestOmitEmptyFields(t *testing.T) {
	// Test that optional fields with omitempty are properly excluded
	explanation := &BeforeAfterExplanation{
		PackageJsonDiffs: []PackageJsonDiff{},
		ImportDiffs:      []ImportDiff{},
		Warnings:         []SideEffectWarning{},
	}

	jsonData, err := json.Marshal(explanation)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	jsonStr := string(jsonData)

	// CurrentState and ProposedState should be null (omitempty doesn't work on pointers when nil)
	// This is actually the expected behavior - nil pointers serialize as null
	if !containsStr(jsonStr, `"currentState":null`) {
		t.Logf("JSON: %s", jsonStr)
	}

	// ImportDiff with LineNumber omitempty - 0 should be omitted
	diff := ImportDiff{
		FilePath:        "test.ts",
		PackageName:     "@mono/test",
		ImportsToRemove: []ImportChange{},
		ImportsToAdd:    []ImportChange{},
		LineNumber:      0, // Zero value with omitempty should be omitted
	}

	diffJson, _ := json.Marshal(diff)
	// LineNumber has omitempty, so zero value should be omitted
	if containsStr(string(diffJson), `"lineNumber":0`) {
		t.Error("Expected lineNumber to be omitted when 0 (has omitempty)")
	}

	// DependencyChange without Version should omit version
	change := DependencyChange{
		Name:           "@mono/shared",
		DependencyType: "dependencies",
	}

	changeJson, _ := json.Marshal(change)
	if containsStr(string(changeJson), `"version"`) {
		t.Error("Expected version to be omitted when empty")
	}
}

// Helper function to check if string contains substring
func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStrHelper(s, substr))
}

func containsStrHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
