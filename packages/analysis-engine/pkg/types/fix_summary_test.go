// Package types defines Go types that match TypeScript definitions in @monoguard/types.
// This file contains tests for FixSummary types for Story 3.8.
package types

import (
	"encoding/json"
	"testing"
)

func TestFixSummary_JSONSerialization(t *testing.T) {
	tests := []struct {
		name     string
		summary  FixSummary
		expected map[string]interface{}
	}{
		{
			name: "complete fix summary",
			summary: FixSummary{
				TotalCircularDependencies: 5,
				TotalEstimatedFixTime:     "2 hours 30 minutes",
				QuickWinsCount:            3,
				CriticalCyclesCount:       1,
				HighPriorityCycles: []PriorityCycleSummary{
					{
						CycleID:          "core→ui",
						PackagesInvolved: []string{"@mono/core", "@mono/ui"},
						PriorityScore:    80.0,
						RecommendedFix:   FixStrategyExtractModule,
						EstimatedTime:    "30-60 minutes",
					},
				},
			},
			expected: map[string]interface{}{
				"totalCircularDependencies": float64(5),
				"totalEstimatedFixTime":     "2 hours 30 minutes",
				"quickWinsCount":            float64(3),
				"criticalCyclesCount":       float64(1),
			},
		},
		{
			name: "empty high priority cycles",
			summary: FixSummary{
				TotalCircularDependencies: 0,
				TotalEstimatedFixTime:     "0 minutes",
				QuickWinsCount:            0,
				CriticalCyclesCount:       0,
				HighPriorityCycles:        []PriorityCycleSummary{},
			},
			expected: map[string]interface{}{
				"totalCircularDependencies": float64(0),
				"totalEstimatedFixTime":     "0 minutes",
				"quickWinsCount":            float64(0),
				"criticalCyclesCount":       float64(0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.summary)
			if err != nil {
				t.Fatalf("Failed to marshal FixSummary: %v", err)
			}

			var result map[string]interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			for key, expectedValue := range tt.expected {
				if actualValue, ok := result[key]; !ok {
					t.Errorf("Expected key %q not found in JSON output", key)
				} else if actualValue != expectedValue {
					t.Errorf("Key %q: expected %v, got %v", key, expectedValue, actualValue)
				}
			}
		})
	}
}

func TestPriorityCycleSummary_JSONSerialization(t *testing.T) {
	summary := PriorityCycleSummary{
		CycleID:          "core→ui",
		PackagesInvolved: []string{"@mono/core", "@mono/ui"},
		PriorityScore:    80.0,
		RecommendedFix:   FixStrategyExtractModule,
		EstimatedTime:    "30-60 minutes",
	}

	data, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("Failed to marshal PriorityCycleSummary: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify camelCase keys
	expectedKeys := []string{
		"cycleId",
		"packagesInvolved",
		"priorityScore",
		"recommendedFix",
		"estimatedTime",
	}

	for _, key := range expectedKeys {
		if _, ok := result[key]; !ok {
			t.Errorf("Expected camelCase key %q not found", key)
		}
	}

	// Verify snake_case is NOT used
	snakeCaseKeys := []string{
		"cycle_id",
		"packages_involved",
		"priority_score",
		"recommended_fix",
		"estimated_time",
	}

	jsonStr := string(data)
	for _, key := range snakeCaseKeys {
		if containsFS(jsonStr, `"`+key+`"`) {
			t.Errorf("Unexpected snake_case key %q in JSON output", key)
		}
	}
}

func TestFixSummary_JSONDeserialization(t *testing.T) {
	jsonData := `{
		"totalCircularDependencies": 3,
		"totalEstimatedFixTime": "1 hour 15 minutes",
		"quickWinsCount": 2,
		"criticalCyclesCount": 1,
		"highPriorityCycles": [
			{
				"cycleId": "auth→user",
				"packagesInvolved": ["@app/auth", "@app/user"],
				"priorityScore": 75.5,
				"recommendedFix": "dependency-injection",
				"estimatedTime": "1-2 hours"
			}
		]
	}`

	var summary FixSummary
	if err := json.Unmarshal([]byte(jsonData), &summary); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if summary.TotalCircularDependencies != 3 {
		t.Errorf("Expected totalCircularDependencies 3, got %d", summary.TotalCircularDependencies)
	}
	if summary.TotalEstimatedFixTime != "1 hour 15 minutes" {
		t.Errorf("Expected totalEstimatedFixTime %q, got %q", "1 hour 15 minutes", summary.TotalEstimatedFixTime)
	}
	if summary.QuickWinsCount != 2 {
		t.Errorf("Expected quickWinsCount 2, got %d", summary.QuickWinsCount)
	}
	if summary.CriticalCyclesCount != 1 {
		t.Errorf("Expected criticalCyclesCount 1, got %d", summary.CriticalCyclesCount)
	}
	if len(summary.HighPriorityCycles) != 1 {
		t.Fatalf("Expected 1 high priority cycle, got %d", len(summary.HighPriorityCycles))
	}

	cycle := summary.HighPriorityCycles[0]
	if cycle.CycleID != "auth→user" {
		t.Errorf("Expected cycleId %q, got %q", "auth→user", cycle.CycleID)
	}
	if cycle.PriorityScore != 75.5 {
		t.Errorf("Expected priorityScore 75.5, got %f", cycle.PriorityScore)
	}
	if cycle.RecommendedFix != FixStrategyDependencyInject {
		t.Errorf("Expected recommendedFix %q, got %q", FixStrategyDependencyInject, cycle.RecommendedFix)
	}
}

func TestFixSummary_EmptySliceNotNil(t *testing.T) {
	summary := FixSummary{
		TotalCircularDependencies: 0,
		TotalEstimatedFixTime:     "0 minutes",
		QuickWinsCount:            0,
		CriticalCyclesCount:       0,
		HighPriorityCycles:        []PriorityCycleSummary{},
	}

	data, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Empty slice should serialize as [] not null
	jsonStr := string(data)
	if containsFS(jsonStr, `"highPriorityCycles":null`) {
		t.Error("Empty HighPriorityCycles should serialize as [] not null")
	}
	if !containsFS(jsonStr, `"highPriorityCycles":[]`) {
		t.Error("Expected HighPriorityCycles to serialize as []")
	}
}

func containsFS(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
