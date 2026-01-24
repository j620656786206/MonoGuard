// Package types provides tests for impact assessment types (Story 3.6).
package types

import (
	"encoding/json"
	"testing"
)

func TestImpactAssessmentJSONSerialization(t *testing.T) {
	tests := []struct {
		name     string
		input    *ImpactAssessment
		wantJSON string
	}{
		{
			name: "full impact assessment",
			input: &ImpactAssessment{
				DirectParticipants: []string{"@mono/ui", "@mono/api", "@mono/core"},
				IndirectDependents: []IndirectDependent{
					{
						PackageName:    "@mono/app",
						DependsOn:      "@mono/ui",
						Distance:       1,
						DependencyPath: []string{"@mono/ui", "@mono/app"},
					},
				},
				TotalAffected:             4,
				AffectedPercentage:        0.4,
				AffectedPercentageDisplay: "40%",
				RiskLevel:                 RiskLevelHigh,
				RiskExplanation:           "High impact: 40% of packages affected",
				RippleEffect: &RippleEffect{
					Layers: []RippleLayer{
						{Distance: 0, Packages: []string{"@mono/ui", "@mono/api", "@mono/core"}, Count: 3},
						{Distance: 1, Packages: []string{"@mono/app"}, Count: 1},
					},
					TotalLayers: 2,
				},
			},
			wantJSON: `{"directParticipants":["@mono/ui","@mono/api","@mono/core"],"indirectDependents":[{"packageName":"@mono/app","dependsOn":"@mono/ui","distance":1,"dependencyPath":["@mono/ui","@mono/app"]}],"totalAffected":4,"affectedPercentage":0.4,"affectedPercentageDisplay":"40%","riskLevel":"high","riskExplanation":"High impact: 40% of packages affected","rippleEffect":{"layers":[{"distance":0,"packages":["@mono/ui","@mono/api","@mono/core"],"count":3},{"distance":1,"packages":["@mono/app"],"count":1}],"totalLayers":2}}`,
		},
		{
			name: "minimal impact assessment without ripple effect",
			input: &ImpactAssessment{
				DirectParticipants:        []string{"@mono/a", "@mono/b"},
				IndirectDependents:        []IndirectDependent{},
				TotalAffected:             2,
				AffectedPercentage:        0.1,
				AffectedPercentageDisplay: "10%",
				RiskLevel:                 RiskLevelLow,
				RiskExplanation:           "Low impact: 10% of packages affected",
			},
			wantJSON: `{"directParticipants":["@mono/a","@mono/b"],"indirectDependents":[],"totalAffected":2,"affectedPercentage":0.1,"affectedPercentageDisplay":"10%","riskLevel":"low","riskExplanation":"Low impact: 10% of packages affected"}`,
		},
		{
			name: "critical risk with core package",
			input: &ImpactAssessment{
				DirectParticipants:        []string{"@mono/core", "@mono/utils"},
				IndirectDependents:        []IndirectDependent{},
				TotalAffected:             2,
				AffectedPercentage:        0.2,
				AffectedPercentageDisplay: "20%",
				RiskLevel:                 RiskLevelCritical,
				RiskExplanation:           "Critical impact: cycle includes core/shared package",
			},
			wantJSON: `{"directParticipants":["@mono/core","@mono/utils"],"indirectDependents":[],"totalAffected":2,"affectedPercentage":0.2,"affectedPercentageDisplay":"20%","riskLevel":"critical","riskExplanation":"Critical impact: cycle includes core/shared package"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshalling
			gotJSON, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}

			if string(gotJSON) != tt.wantJSON {
				t.Errorf("json.Marshal() = %s, want %s", string(gotJSON), tt.wantJSON)
			}

			// Test unmarshalling (round trip)
			var unmarshalled ImpactAssessment
			if err := json.Unmarshal(gotJSON, &unmarshalled); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}

			// Verify key fields
			if len(unmarshalled.DirectParticipants) != len(tt.input.DirectParticipants) {
				t.Errorf("DirectParticipants length = %d, want %d",
					len(unmarshalled.DirectParticipants), len(tt.input.DirectParticipants))
			}
			if unmarshalled.RiskLevel != tt.input.RiskLevel {
				t.Errorf("RiskLevel = %s, want %s", unmarshalled.RiskLevel, tt.input.RiskLevel)
			}
		})
	}
}

func TestIndirectDependentJSONSerialization(t *testing.T) {
	dep := IndirectDependent{
		PackageName:    "@mono/dashboard",
		DependsOn:      "@mono/ui",
		Distance:       2,
		DependencyPath: []string{"@mono/ui", "@mono/app", "@mono/dashboard"},
	}

	gotJSON, err := json.Marshal(dep)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	wantJSON := `{"packageName":"@mono/dashboard","dependsOn":"@mono/ui","distance":2,"dependencyPath":["@mono/ui","@mono/app","@mono/dashboard"]}`
	if string(gotJSON) != wantJSON {
		t.Errorf("json.Marshal() = %s, want %s", string(gotJSON), wantJSON)
	}
}

func TestRippleLayerJSONSerialization(t *testing.T) {
	layer := RippleLayer{
		Distance: 1,
		Packages: []string{"@mono/app", "@mono/web"},
		Count:    2,
	}

	gotJSON, err := json.Marshal(layer)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	wantJSON := `{"distance":1,"packages":["@mono/app","@mono/web"],"count":2}`
	if string(gotJSON) != wantJSON {
		t.Errorf("json.Marshal() = %s, want %s", string(gotJSON), wantJSON)
	}
}

func TestRiskLevelConstants(t *testing.T) {
	tests := []struct {
		level RiskLevel
		want  string
	}{
		{RiskLevelCritical, "critical"},
		{RiskLevelHigh, "high"},
		{RiskLevelMedium, "medium"},
		{RiskLevelLow, "low"},
	}

	for _, tt := range tests {
		t.Run(string(tt.level), func(t *testing.T) {
			if string(tt.level) != tt.want {
				t.Errorf("RiskLevel = %s, want %s", tt.level, tt.want)
			}
		})
	}
}

func TestNewImpactAssessment(t *testing.T) {
	ia := NewImpactAssessment()

	if ia == nil {
		t.Fatal("NewImpactAssessment() returned nil")
	}

	if ia.DirectParticipants == nil {
		t.Error("DirectParticipants should be initialized as empty slice, not nil")
	}

	if ia.IndirectDependents == nil {
		t.Error("IndirectDependents should be initialized as empty slice, not nil")
	}

	// Verify empty slices serialize as [] not null
	gotJSON, err := json.Marshal(ia)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	if string(gotJSON) == `null` {
		t.Error("NewImpactAssessment() should not serialize to null")
	}
}

func TestCalculatePercentage(t *testing.T) {
	tests := []struct {
		name           string
		affected       int
		total          int
		wantPercentage float64
		wantDisplay    string
	}{
		{
			name:           "zero total",
			affected:       0,
			total:          0,
			wantPercentage: 0.0,
			wantDisplay:    "0%",
		},
		{
			name:           "zero affected",
			affected:       0,
			total:          10,
			wantPercentage: 0.0,
			wantDisplay:    "0%",
		},
		{
			name:           "half affected",
			affected:       5,
			total:          10,
			wantPercentage: 0.5,
			wantDisplay:    "50%",
		},
		{
			name:           "all affected",
			affected:       10,
			total:          10,
			wantPercentage: 1.0,
			wantDisplay:    "100%",
		},
		{
			name:           "quarter affected",
			affected:       25,
			total:          100,
			wantPercentage: 0.25,
			wantDisplay:    "25%",
		},
		{
			name:           "more than total (capped at 1.0)",
			affected:       15,
			total:          10,
			wantPercentage: 1.0,
			wantDisplay:    "100%",
		},
		{
			name:           "low percentage",
			affected:       1,
			total:          100,
			wantPercentage: 0.01,
			wantDisplay:    "1%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPercentage, gotDisplay := CalculatePercentage(tt.affected, tt.total)

			if gotPercentage != tt.wantPercentage {
				t.Errorf("CalculatePercentage() percentage = %v, want %v", gotPercentage, tt.wantPercentage)
			}
			if gotDisplay != tt.wantDisplay {
				t.Errorf("CalculatePercentage() display = %v, want %v", gotDisplay, tt.wantDisplay)
			}
		})
	}
}

func TestCamelCaseJSONTags(t *testing.T) {
	// Verify all JSON tags use camelCase (not snake_case)
	ia := &ImpactAssessment{
		DirectParticipants:        []string{"a"},
		IndirectDependents:        []IndirectDependent{},
		TotalAffected:             1,
		AffectedPercentage:        0.1,
		AffectedPercentageDisplay: "10%",
		RiskLevel:                 RiskLevelLow,
		RiskExplanation:           "test",
	}

	gotJSON, err := json.Marshal(ia)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	jsonStr := string(gotJSON)

	// Check for correct camelCase keys
	camelCaseKeys := []string{
		"directParticipants",
		"indirectDependents",
		"totalAffected",
		"affectedPercentage",
		"affectedPercentageDisplay",
		"riskLevel",
		"riskExplanation",
	}

	for _, key := range camelCaseKeys {
		if !containsImpact(jsonStr, `"`+key+`"`) {
			t.Errorf("JSON should contain camelCase key %q", key)
		}
	}

	// Check that snake_case keys are NOT present
	snakeCaseKeys := []string{
		"direct_participants",
		"indirect_dependents",
		"total_affected",
		"affected_percentage",
		"risk_level",
		"risk_explanation",
	}

	for _, key := range snakeCaseKeys {
		if containsImpact(jsonStr, `"`+key+`"`) {
			t.Errorf("JSON should NOT contain snake_case key %q", key)
		}
	}
}

func containsImpact(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsImpactHelper(s, substr))
}

func containsImpactHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
