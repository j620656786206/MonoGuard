package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestHealthScoreResult_JSONSerialization(t *testing.T) {
	result := &HealthScoreResult{
		Overall: 72,
		Rating:  HealthRatingGood,
		Breakdown: &ScoreBreakdown{
			CircularScore: 70,
			ConflictScore: 85,
			DepthScore:    60,
			CouplingScore: 75,
		},
		Factors: []*HealthFactor{
			{
				Name:            "Circular Dependencies",
				Score:           70,
				Weight:          0.40,
				WeightedScore:   28,
				Description:     "2 cycles detected",
				Recommendations: []string{"Break cycle between pkg-a and pkg-b"},
			},
			{
				Name:            "Version Conflicts",
				Score:           85,
				Weight:          0.25,
				WeightedScore:   21,
				Description:     "3 conflicts detected",
				Recommendations: []string{},
			},
		},
		UpdatedAt: "2026-01-17T10:30:00Z",
	}

	// Test serialization
	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal HealthScoreResult: %v", err)
	}

	// Verify camelCase JSON keys
	jsonStr := string(data)
	expectedKeys := []string{
		`"overall"`,
		`"rating"`,
		`"breakdown"`,
		`"factors"`,
		`"updatedAt"`,
		`"circularScore"`,
		`"conflictScore"`,
		`"depthScore"`,
		`"couplingScore"`,
		`"name"`,
		`"score"`,
		`"weight"`,
		`"weightedScore"`,
		`"description"`,
		`"recommendations"`,
	}

	for _, key := range expectedKeys {
		if !strings.Contains(jsonStr, key) {
			t.Errorf("Expected JSON to contain key %s, but it didn't. JSON: %s", key, jsonStr)
		}
	}

	// Test deserialization
	var decoded HealthScoreResult
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal HealthScoreResult: %v", err)
	}

	if decoded.Overall != result.Overall {
		t.Errorf("Overall mismatch: got %d, want %d", decoded.Overall, result.Overall)
	}

	if decoded.Rating != result.Rating {
		t.Errorf("Rating mismatch: got %s, want %s", decoded.Rating, result.Rating)
	}

	if decoded.Breakdown.CircularScore != result.Breakdown.CircularScore {
		t.Errorf("CircularScore mismatch: got %d, want %d",
			decoded.Breakdown.CircularScore, result.Breakdown.CircularScore)
	}

	if len(decoded.Factors) != len(result.Factors) {
		t.Errorf("Factors length mismatch: got %d, want %d",
			len(decoded.Factors), len(result.Factors))
	}
}

func TestGetHealthRating(t *testing.T) {
	tests := []struct {
		score    int
		expected HealthRating
	}{
		// Excellent: 85-100
		{100, HealthRatingExcellent},
		{95, HealthRatingExcellent},
		{85, HealthRatingExcellent},

		// Good: 70-84
		{84, HealthRatingGood},
		{75, HealthRatingGood},
		{70, HealthRatingGood},

		// Fair: 50-69
		{69, HealthRatingFair},
		{60, HealthRatingFair},
		{50, HealthRatingFair},

		// Poor: 30-49
		{49, HealthRatingPoor},
		{40, HealthRatingPoor},
		{30, HealthRatingPoor},

		// Critical: 0-29
		{29, HealthRatingCritical},
		{15, HealthRatingCritical},
		{0, HealthRatingCritical},

		// Edge cases
		{-5, HealthRatingCritical},  // Below 0
		{105, HealthRatingExcellent}, // Above 100
	}

	for _, tt := range tests {
		// Use unique test name with score to avoid collision when same rating appears multiple times
		testName := fmt.Sprintf("score_%d_expects_%s", tt.score, tt.expected)
		t.Run(testName, func(t *testing.T) {
			result := GetHealthRating(tt.score)
			if result != tt.expected {
				t.Errorf("GetHealthRating(%d) = %s, want %s", tt.score, result, tt.expected)
			}
		})
	}
}

func TestHealthRating_Values(t *testing.T) {
	tests := []struct {
		rating   HealthRating
		expected string
	}{
		{HealthRatingExcellent, "excellent"},
		{HealthRatingGood, "good"},
		{HealthRatingFair, "fair"},
		{HealthRatingPoor, "poor"},
		{HealthRatingCritical, "critical"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if string(tt.rating) != tt.expected {
				t.Errorf("HealthRating value mismatch: got %s, want %s", tt.rating, tt.expected)
			}
		})
	}
}

func TestScoreBreakdown_JSONSerialization(t *testing.T) {
	breakdown := &ScoreBreakdown{
		CircularScore: 100,
		ConflictScore: 90,
		DepthScore:    80,
		CouplingScore: 70,
	}

	data, err := json.Marshal(breakdown)
	if err != nil {
		t.Fatalf("Failed to marshal ScoreBreakdown: %v", err)
	}

	var decoded ScoreBreakdown
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal ScoreBreakdown: %v", err)
	}

	if decoded.CircularScore != breakdown.CircularScore {
		t.Errorf("CircularScore mismatch: got %d, want %d", decoded.CircularScore, breakdown.CircularScore)
	}
	if decoded.ConflictScore != breakdown.ConflictScore {
		t.Errorf("ConflictScore mismatch: got %d, want %d", decoded.ConflictScore, breakdown.ConflictScore)
	}
	if decoded.DepthScore != breakdown.DepthScore {
		t.Errorf("DepthScore mismatch: got %d, want %d", decoded.DepthScore, breakdown.DepthScore)
	}
	if decoded.CouplingScore != breakdown.CouplingScore {
		t.Errorf("CouplingScore mismatch: got %d, want %d", decoded.CouplingScore, breakdown.CouplingScore)
	}
}

func TestHealthFactor_JSONSerialization(t *testing.T) {
	factor := &HealthFactor{
		Name:            "Test Factor",
		Score:           75,
		Weight:          0.25,
		WeightedScore:   19,
		Description:     "Test description",
		Recommendations: []string{"Recommendation 1", "Recommendation 2"},
	}

	data, err := json.Marshal(factor)
	if err != nil {
		t.Fatalf("Failed to marshal HealthFactor: %v", err)
	}

	var decoded HealthFactor
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal HealthFactor: %v", err)
	}

	if decoded.Name != factor.Name {
		t.Errorf("Name mismatch: got %s, want %s", decoded.Name, factor.Name)
	}
	if decoded.Score != factor.Score {
		t.Errorf("Score mismatch: got %d, want %d", decoded.Score, factor.Score)
	}
	if decoded.Weight != factor.Weight {
		t.Errorf("Weight mismatch: got %f, want %f", decoded.Weight, factor.Weight)
	}
	if decoded.WeightedScore != factor.WeightedScore {
		t.Errorf("WeightedScore mismatch: got %d, want %d", decoded.WeightedScore, factor.WeightedScore)
	}
	if len(decoded.Recommendations) != len(factor.Recommendations) {
		t.Errorf("Recommendations length mismatch: got %d, want %d",
			len(decoded.Recommendations), len(factor.Recommendations))
	}
}

// Note: Using strings.Contains from standard library instead of custom helper
