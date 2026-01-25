// Package types defines Go types that match TypeScript definitions in @monoguard/types.
// This file contains tests for QuickFixSummary types for Story 3.8.
package types

import (
	"encoding/json"
	"testing"
)

func TestQuickFixSummary_JSONSerialization(t *testing.T) {
	tests := []struct {
		name     string
		summary  QuickFixSummary
		expected map[string]interface{}
	}{
		{
			name: "complete quick fix summary",
			summary: QuickFixSummary{
				StrategyType:  FixStrategyExtractModule,
				StrategyName:  "Extract Shared Module",
				Summary:       "Create new shared package '@mono/shared' to break the cycle",
				Suitability:   8,
				Effort:        EffortMedium,
				EstimatedTime: "30-60 minutes",
				Guide: &FixGuide{
					StrategyType:  FixStrategyExtractModule,
					Title:         "Extract Shared Module Guide",
					Summary:       "Step-by-step guide",
					Steps:         []FixStep{},
					Verification:  []FixStep{},
					EstimatedTime: "30-60 minutes",
				},
				StrategyIndex: 0,
			},
			expected: map[string]interface{}{
				"strategyType":  "extract-module",
				"strategyName":  "Extract Shared Module",
				"summary":       "Create new shared package '@mono/shared' to break the cycle",
				"suitability":   float64(8),
				"effort":        "medium",
				"estimatedTime": "30-60 minutes",
				"strategyIndex": float64(0),
			},
		},
		{
			name: "quick fix without guide",
			summary: QuickFixSummary{
				StrategyType:  FixStrategyDependencyInject,
				StrategyName:  "Dependency Injection",
				Summary:       "Invert dependency using dependency injection pattern",
				Suitability:   6,
				Effort:        EffortHigh,
				EstimatedTime: "1-2 hours",
				Guide:         nil,
				StrategyIndex: 1,
			},
			expected: map[string]interface{}{
				"strategyType":  "dependency-injection",
				"strategyName":  "Dependency Injection",
				"summary":       "Invert dependency using dependency injection pattern",
				"suitability":   float64(6),
				"effort":        "high",
				"estimatedTime": "1-2 hours",
				"strategyIndex": float64(1),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Serialize to JSON
			data, err := json.Marshal(tt.summary)
			if err != nil {
				t.Fatalf("Failed to marshal QuickFixSummary: %v", err)
			}

			// Parse back to map for verification
			var result map[string]interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			// Verify key fields use camelCase
			for key, expectedValue := range tt.expected {
				if actualValue, ok := result[key]; !ok {
					t.Errorf("Expected key %q not found in JSON output", key)
				} else if actualValue != expectedValue {
					t.Errorf("Key %q: expected %v, got %v", key, expectedValue, actualValue)
				}
			}

			// Verify guide is omitted when nil
			if tt.summary.Guide == nil {
				if _, ok := result["guide"]; ok {
					t.Error("Expected guide to be omitted when nil")
				}
			}
		})
	}
}

func TestQuickFixSummary_JSONDeserialization(t *testing.T) {
	jsonData := `{
		"strategyType": "extract-module",
		"strategyName": "Extract Shared Module",
		"summary": "Create new shared package to break the cycle",
		"suitability": 8,
		"effort": "medium",
		"estimatedTime": "30-60 minutes",
		"strategyIndex": 0
	}`

	var summary QuickFixSummary
	if err := json.Unmarshal([]byte(jsonData), &summary); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if summary.StrategyType != FixStrategyExtractModule {
		t.Errorf("Expected strategyType %q, got %q", FixStrategyExtractModule, summary.StrategyType)
	}
	if summary.StrategyName != "Extract Shared Module" {
		t.Errorf("Expected strategyName %q, got %q", "Extract Shared Module", summary.StrategyName)
	}
	if summary.Suitability != 8 {
		t.Errorf("Expected suitability 8, got %d", summary.Suitability)
	}
	if summary.Effort != EffortMedium {
		t.Errorf("Expected effort %q, got %q", EffortMedium, summary.Effort)
	}
	if summary.EstimatedTime != "30-60 minutes" {
		t.Errorf("Expected estimatedTime %q, got %q", "30-60 minutes", summary.EstimatedTime)
	}
	if summary.StrategyIndex != 0 {
		t.Errorf("Expected strategyIndex 0, got %d", summary.StrategyIndex)
	}
}

func TestQuickFixSummary_CamelCaseJSONTags(t *testing.T) {
	summary := QuickFixSummary{
		StrategyType:  FixStrategyExtractModule,
		StrategyName:  "Test",
		Summary:       "Test summary",
		Suitability:   5,
		Effort:        EffortLow,
		EstimatedTime: "5-10 minutes",
		StrategyIndex: 0,
	}

	data, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	jsonStr := string(data)

	// Verify camelCase (correct)
	camelCaseKeys := []string{
		`"strategyType"`,
		`"strategyName"`,
		`"suitability"`,
		`"effort"`,
		`"estimatedTime"`,
		`"strategyIndex"`,
	}

	for _, key := range camelCaseKeys {
		if !containsQFS(jsonStr, key) {
			t.Errorf("Expected camelCase key %s in JSON output", key)
		}
	}

	// Verify snake_case is NOT used (wrong)
	snakeCaseKeys := []string{
		`"strategy_type"`,
		`"strategy_name"`,
		`"estimated_time"`,
		`"strategy_index"`,
	}

	for _, key := range snakeCaseKeys {
		if containsQFS(jsonStr, key) {
			t.Errorf("Unexpected snake_case key %s in JSON output", key)
		}
	}
}

func containsQFS(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsQFSHelper(s, substr))
}

func containsQFSHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
