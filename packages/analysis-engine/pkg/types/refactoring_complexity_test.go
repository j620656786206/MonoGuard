// Package types defines Go types that match TypeScript definitions in @monoguard/types.
// This file contains tests for refactoring complexity types for Story 3.5.
package types

import (
	"encoding/json"
	"testing"
)

func TestComplexityFactor_JSONSerialization(t *testing.T) {
	factor := ComplexityFactor{
		Value:        5,
		Weight:       0.25,
		Contribution: 1.5,
		Description:  "5 source files need modification",
	}

	got, err := json.Marshal(factor)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	// Verify camelCase JSON keys
	var result map[string]interface{}
	if err := json.Unmarshal(got, &result); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	// Check required fields exist with correct keys (camelCase)
	expectedKeys := []string{"value", "weight", "contribution", "description"}
	for _, key := range expectedKeys {
		if _, ok := result[key]; !ok {
			t.Errorf("Expected JSON key %q not found", key)
		}
	}

	// Verify no snake_case keys
	badKeys := []string{"Value", "Weight", "Contribution", "Description"}
	for _, key := range badKeys {
		if _, ok := result[key]; ok {
			t.Errorf("Unexpected PascalCase JSON key %q found", key)
		}
	}
}

func TestComplexityFactor_JSONRoundTrip(t *testing.T) {
	original := ComplexityFactor{
		Value:        8,
		Weight:       0.20,
		Contribution: 1.6,
		Description:  "8 import statements need updating",
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var decoded ComplexityFactor
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if decoded.Value != original.Value {
		t.Errorf("Value = %v, want %v", decoded.Value, original.Value)
	}
	if decoded.Weight != original.Weight {
		t.Errorf("Weight = %v, want %v", decoded.Weight, original.Weight)
	}
	if decoded.Contribution != original.Contribution {
		t.Errorf("Contribution = %v, want %v", decoded.Contribution, original.Contribution)
	}
	if decoded.Description != original.Description {
		t.Errorf("Description = %v, want %v", decoded.Description, original.Description)
	}
}

func TestComplexityBreakdown_JSONSerialization(t *testing.T) {
	breakdown := ComplexityBreakdown{
		FilesAffected: ComplexityFactor{
			Value:        3,
			Weight:       0.25,
			Contribution: 1.5,
			Description:  "3 source files need modification",
		},
		ImportsToChange: ComplexityFactor{
			Value:        4,
			Weight:       0.20,
			Contribution: 1.2,
			Description:  "4 import statements need updating",
		},
		ChainDepth: ComplexityFactor{
			Value:        3,
			Weight:       0.25,
			Contribution: 1.5,
			Description:  "Dependency chain has 3 levels",
		},
		PackagesInvolved: ComplexityFactor{
			Value:        3,
			Weight:       0.15,
			Contribution: 0.9,
			Description:  "3 packages involved in cycle",
		},
		ExternalDependencies: ComplexityFactor{
			Value:        0,
			Weight:       0.15,
			Contribution: 0.15,
			Description:  "No external dependencies in cycle",
		},
	}

	got, err := json.Marshal(breakdown)
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
		"filesAffected", "importsToChange", "chainDepth",
		"packagesInvolved", "externalDependencies",
	}
	for _, key := range expectedKeys {
		if _, ok := result[key]; !ok {
			t.Errorf("Expected JSON key %q not found", key)
		}
	}

	// Verify no snake_case keys
	badKeys := []string{
		"files_affected", "imports_to_change", "chain_depth",
		"packages_involved", "external_dependencies",
	}
	for _, key := range badKeys {
		if _, ok := result[key]; ok {
			t.Errorf("Unexpected snake_case JSON key %q found", key)
		}
	}
}

func TestRefactoringComplexity_JSONSerialization(t *testing.T) {
	complexity := RefactoringComplexity{
		Score:         5,
		EstimatedTime: "30-60 minutes",
		Breakdown: ComplexityBreakdown{
			FilesAffected: ComplexityFactor{
				Value:        3,
				Weight:       0.25,
				Contribution: 1.5,
				Description:  "3 source files need modification",
			},
			ImportsToChange: ComplexityFactor{
				Value:        3,
				Weight:       0.20,
				Contribution: 1.2,
				Description:  "3 import statements need updating",
			},
			ChainDepth: ComplexityFactor{
				Value:        3,
				Weight:       0.25,
				Contribution: 1.5,
				Description:  "Dependency chain has 3 levels",
			},
			PackagesInvolved: ComplexityFactor{
				Value:        3,
				Weight:       0.15,
				Contribution: 0.9,
				Description:  "3 packages involved in cycle",
			},
			ExternalDependencies: ComplexityFactor{
				Value:        0,
				Weight:       0.15,
				Contribution: 0.15,
				Description:  "No external dependencies in cycle",
			},
		},
		Explanation: "Moderate refactoring: 3 files, 3 imports, 3-level chain",
	}

	got, err := json.Marshal(complexity)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	// Verify camelCase JSON keys
	var result map[string]interface{}
	if err := json.Unmarshal(got, &result); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	// Check required fields exist with correct keys (camelCase)
	expectedKeys := []string{"score", "estimatedTime", "breakdown", "explanation"}
	for _, key := range expectedKeys {
		if _, ok := result[key]; !ok {
			t.Errorf("Expected JSON key %q not found", key)
		}
	}

	// Verify no snake_case keys
	badKeys := []string{"estimated_time"}
	for _, key := range badKeys {
		if _, ok := result[key]; ok {
			t.Errorf("Unexpected snake_case JSON key %q found", key)
		}
	}
}

func TestRefactoringComplexity_JSONRoundTrip(t *testing.T) {
	original := RefactoringComplexity{
		Score:         7,
		EstimatedTime: "1-2 hours",
		Breakdown: ComplexityBreakdown{
			FilesAffected: ComplexityFactor{
				Value:        8,
				Weight:       0.25,
				Contribution: 2.0,
				Description:  "8 source files need modification",
			},
			ImportsToChange: ComplexityFactor{
				Value:        10,
				Weight:       0.20,
				Contribution: 1.6,
				Description:  "10 import statements need updating",
			},
			ChainDepth: ComplexityFactor{
				Value:        5,
				Weight:       0.25,
				Contribution: 2.0,
				Description:  "Dependency chain has 5 levels",
			},
			PackagesInvolved: ComplexityFactor{
				Value:        5,
				Weight:       0.15,
				Contribution: 1.2,
				Description:  "5 packages involved in cycle",
			},
			ExternalDependencies: ComplexityFactor{
				Value:        1,
				Weight:       0.15,
				Contribution: 1.5,
				Description:  "External dependencies increase complexity",
			},
		},
		Explanation: "Significant refactoring: 8 files, 10 imports, 5-level chain",
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var decoded RefactoringComplexity
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if decoded.Score != original.Score {
		t.Errorf("Score = %v, want %v", decoded.Score, original.Score)
	}
	if decoded.EstimatedTime != original.EstimatedTime {
		t.Errorf("EstimatedTime = %v, want %v", decoded.EstimatedTime, original.EstimatedTime)
	}
	if decoded.Explanation != original.Explanation {
		t.Errorf("Explanation = %v, want %v", decoded.Explanation, original.Explanation)
	}
	if decoded.Breakdown.FilesAffected.Value != original.Breakdown.FilesAffected.Value {
		t.Errorf("Breakdown.FilesAffected.Value = %v, want %v",
			decoded.Breakdown.FilesAffected.Value, original.Breakdown.FilesAffected.Value)
	}
	if decoded.Breakdown.ExternalDependencies.Contribution != original.Breakdown.ExternalDependencies.Contribution {
		t.Errorf("Breakdown.ExternalDependencies.Contribution = %v, want %v",
			decoded.Breakdown.ExternalDependencies.Contribution, original.Breakdown.ExternalDependencies.Contribution)
	}
}

func TestRefactoringComplexity_ScoreRange(t *testing.T) {
	tests := []struct {
		name  string
		score int
		valid bool
	}{
		{"score 0 invalid", 0, false},
		{"score 1 valid", 1, true},
		{"score 5 valid", 5, true},
		{"score 10 valid", 10, true},
		{"score 11 invalid", 11, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.score >= 1 && tt.score <= 10
			if valid != tt.valid {
				t.Errorf("score %d validity = %v, want %v", tt.score, valid, tt.valid)
			}
		})
	}
}

func TestRefactoringComplexity_WeightSum(t *testing.T) {
	// All weights should sum to 1.0
	breakdown := ComplexityBreakdown{
		FilesAffected:        ComplexityFactor{Weight: 0.25},
		ImportsToChange:      ComplexityFactor{Weight: 0.20},
		ChainDepth:           ComplexityFactor{Weight: 0.25},
		PackagesInvolved:     ComplexityFactor{Weight: 0.15},
		ExternalDependencies: ComplexityFactor{Weight: 0.15},
	}

	totalWeight := breakdown.FilesAffected.Weight +
		breakdown.ImportsToChange.Weight +
		breakdown.ChainDepth.Weight +
		breakdown.PackagesInvolved.Weight +
		breakdown.ExternalDependencies.Weight

	expected := 1.0
	if totalWeight != expected {
		t.Errorf("Total weight = %v, want %v", totalWeight, expected)
	}
}

func TestRefactoringComplexity_EstimatedTimeRanges(t *testing.T) {
	validTimeRanges := []string{
		"5-15 minutes",
		"15-30 minutes",
		"30-60 minutes",
		"1-2 hours",
		"2-4 hours",
	}

	for _, timeRange := range validTimeRanges {
		complexity := RefactoringComplexity{
			Score:         5,
			EstimatedTime: timeRange,
		}

		data, err := json.Marshal(complexity)
		if err != nil {
			t.Fatalf("json.Marshal() error = %v for time range %q", err, timeRange)
		}

		var decoded RefactoringComplexity
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("json.Unmarshal() error = %v for time range %q", err, timeRange)
		}

		if decoded.EstimatedTime != timeRange {
			t.Errorf("EstimatedTime = %v, want %v", decoded.EstimatedTime, timeRange)
		}
	}
}
