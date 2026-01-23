// Package types defines Go types that match TypeScript definitions in @monoguard/types.
// This file contains tests for fix strategy types for Story 3.3.
package types

import (
	"encoding/json"
	"testing"
)

func TestFixStrategyType_JSONSerialization(t *testing.T) {
	tests := []struct {
		name     string
		strategy FixStrategyType
		wantJSON string
	}{
		{
			name:     "extract-module type",
			strategy: FixStrategyExtractModule,
			wantJSON: `"extract-module"`,
		},
		{
			name:     "dependency-injection type",
			strategy: FixStrategyDependencyInject,
			wantJSON: `"dependency-injection"`,
		},
		{
			name:     "boundary-refactoring type",
			strategy: FixStrategyBoundaryRefactor,
			wantJSON: `"boundary-refactoring"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.strategy)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}
			if string(got) != tt.wantJSON {
				t.Errorf("json.Marshal() = %s, want %s", string(got), tt.wantJSON)
			}
		})
	}
}

func TestEffortLevel_JSONSerialization(t *testing.T) {
	tests := []struct {
		name     string
		effort   EffortLevel
		wantJSON string
	}{
		{
			name:     "low effort",
			effort:   EffortLow,
			wantJSON: `"low"`,
		},
		{
			name:     "medium effort",
			effort:   EffortMedium,
			wantJSON: `"medium"`,
		},
		{
			name:     "high effort",
			effort:   EffortHigh,
			wantJSON: `"high"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.effort)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}
			if string(got) != tt.wantJSON {
				t.Errorf("json.Marshal() = %s, want %s", string(got), tt.wantJSON)
			}
		})
	}
}

func TestFixStrategy_JSONSerialization(t *testing.T) {
	strategy := FixStrategy{
		Type:           FixStrategyExtractModule,
		Name:           "Extract Shared Module",
		Description:    "Create a new shared package to hold common dependencies.",
		Suitability:    8,
		Effort:         EffortMedium,
		Pros:           []string{"Creates clear separation", "Isolates shared code"},
		Cons:           []string{"New package to maintain"},
		Recommended:    true,
		TargetPackages: []string{"@mono/ui", "@mono/api"},
		NewPackageName: "@mono/shared",
	}

	got, err := json.Marshal(strategy)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	// Verify camelCase JSON keys
	var result map[string]interface{}
	if err := json.Unmarshal(got, &result); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	// Check required fields exist with correct keys (camelCase)
	expectedKeys := []string{
		"type", "name", "description", "suitability", "effort",
		"pros", "cons", "recommended", "targetPackages", "newPackageName",
	}
	for _, key := range expectedKeys {
		if _, ok := result[key]; !ok {
			t.Errorf("Expected JSON key %q not found", key)
		}
	}

	// Verify no snake_case keys
	badKeys := []string{"target_packages", "new_package_name"}
	for _, key := range badKeys {
		if _, ok := result[key]; ok {
			t.Errorf("Unexpected snake_case JSON key %q found", key)
		}
	}
}

func TestFixStrategy_JSONRoundTrip(t *testing.T) {
	original := FixStrategy{
		Type:           FixStrategyDependencyInject,
		Name:           "Dependency Injection",
		Description:    "Invert the problematic dependency.",
		Suitability:    7,
		Effort:         EffortLow,
		Pros:           []string{"Minimal changes", "Preserves structure"},
		Cons:           []string{"Adds indirection"},
		Recommended:    false,
		TargetPackages: []string{"@mono/core", "@mono/ui"},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var decoded FixStrategy
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if decoded.Type != original.Type {
		t.Errorf("Type = %v, want %v", decoded.Type, original.Type)
	}
	if decoded.Name != original.Name {
		t.Errorf("Name = %v, want %v", decoded.Name, original.Name)
	}
	if decoded.Suitability != original.Suitability {
		t.Errorf("Suitability = %v, want %v", decoded.Suitability, original.Suitability)
	}
	if decoded.Effort != original.Effort {
		t.Errorf("Effort = %v, want %v", decoded.Effort, original.Effort)
	}
	if decoded.Recommended != original.Recommended {
		t.Errorf("Recommended = %v, want %v", decoded.Recommended, original.Recommended)
	}
	if len(decoded.Pros) != len(original.Pros) {
		t.Errorf("Pros length = %v, want %v", len(decoded.Pros), len(original.Pros))
	}
	if len(decoded.Cons) != len(original.Cons) {
		t.Errorf("Cons length = %v, want %v", len(decoded.Cons), len(original.Cons))
	}
	if len(decoded.TargetPackages) != len(original.TargetPackages) {
		t.Errorf("TargetPackages length = %v, want %v", len(decoded.TargetPackages), len(original.TargetPackages))
	}
}

func TestFixStrategy_OmitEmptyNewPackageName(t *testing.T) {
	strategy := FixStrategy{
		Type:           FixStrategyDependencyInject,
		Name:           "Dependency Injection",
		Description:    "Invert the dependency.",
		Suitability:    6,
		Effort:         EffortMedium,
		Pros:           []string{"Minimal changes"},
		Cons:           []string{"Adds indirection"},
		Recommended:    false,
		TargetPackages: []string{"@mono/a", "@mono/b"},
		// NewPackageName intentionally empty
	}

	got, err := json.Marshal(strategy)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	// Verify newPackageName is omitted when empty
	var result map[string]interface{}
	if err := json.Unmarshal(got, &result); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if _, ok := result["newPackageName"]; ok {
		t.Error("newPackageName should be omitted when empty")
	}
}

func TestFixStrategy_EmptySlicesSerialize(t *testing.T) {
	strategy := FixStrategy{
		Type:           FixStrategyBoundaryRefactor,
		Name:           "Boundary Refactoring",
		Description:    "Restructure boundaries.",
		Suitability:    5,
		Effort:         EffortHigh,
		Pros:           []string{}, // Empty slice
		Cons:           []string{}, // Empty slice
		Recommended:    false,
		TargetPackages: []string{}, // Empty slice
	}

	got, err := json.Marshal(strategy)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	// Verify empty arrays serialize as [] not null
	var result map[string]interface{}
	if err := json.Unmarshal(got, &result); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	// Pros, cons, targetPackages should be arrays not nil
	if result["pros"] == nil {
		t.Error("pros should serialize as empty array, not null")
	}
	if result["cons"] == nil {
		t.Error("cons should serialize as empty array, not null")
	}
	if result["targetPackages"] == nil {
		t.Error("targetPackages should serialize as empty array, not null")
	}
}
