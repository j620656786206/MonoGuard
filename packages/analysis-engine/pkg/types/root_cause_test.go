// Package types tests for RootCauseAnalysis types.
package types

import (
	"encoding/json"
	"testing"
)

func TestRootCauseAnalysis_JSONSerialization(t *testing.T) {
	tests := []struct {
		name     string
		input    *RootCauseAnalysis
		wantJSON string
	}{
		{
			name: "full root cause analysis",
			input: &RootCauseAnalysis{
				OriginatingPackage: "pkg-ui",
				ProblematicDependency: RootCauseEdge{
					From:     "pkg-ui",
					To:       "pkg-api",
					Type:     DependencyTypeProduction,
					Critical: false,
				},
				Confidence:  82,
				Explanation: "Package 'pkg-ui' is highly likely the root cause.",
				Chain: []RootCauseEdge{
					{From: "pkg-ui", To: "pkg-api", Type: DependencyTypeProduction, Critical: false},
					{From: "pkg-api", To: "pkg-core", Type: DependencyTypeProduction, Critical: false},
					{From: "pkg-core", To: "pkg-ui", Type: DependencyTypeProduction, Critical: true},
				},
				CriticalEdge: &RootCauseEdge{
					From:     "pkg-core",
					To:       "pkg-ui",
					Type:     DependencyTypeProduction,
					Critical: true,
				},
			},
			wantJSON: `{"originatingPackage":"pkg-ui","problematicDependency":{"from":"pkg-ui","to":"pkg-api","type":"production","critical":false},"confidence":82,"explanation":"Package 'pkg-ui' is highly likely the root cause.","chain":[{"from":"pkg-ui","to":"pkg-api","type":"production","critical":false},{"from":"pkg-api","to":"pkg-core","type":"production","critical":false},{"from":"pkg-core","to":"pkg-ui","type":"production","critical":true}],"criticalEdge":{"from":"pkg-core","to":"pkg-ui","type":"production","critical":true}}`,
		},
		{
			name: "without critical edge",
			input: &RootCauseAnalysis{
				OriginatingPackage: "pkg-a",
				ProblematicDependency: RootCauseEdge{
					From:     "pkg-a",
					To:       "pkg-b",
					Type:     DependencyTypeDevelopment,
					Critical: false,
				},
				Confidence:  65,
				Explanation: "Package 'pkg-a' is likely the root cause.",
				Chain: []RootCauseEdge{
					{From: "pkg-a", To: "pkg-b", Type: DependencyTypeDevelopment, Critical: false},
					{From: "pkg-b", To: "pkg-a", Type: DependencyTypeProduction, Critical: true},
				},
				CriticalEdge: nil, // omitempty should exclude this
			},
			wantJSON: `{"originatingPackage":"pkg-a","problematicDependency":{"from":"pkg-a","to":"pkg-b","type":"development","critical":false},"confidence":65,"explanation":"Package 'pkg-a' is likely the root cause.","chain":[{"from":"pkg-a","to":"pkg-b","type":"development","critical":false},{"from":"pkg-b","to":"pkg-a","type":"production","critical":true}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}

			if string(got) != tt.wantJSON {
				t.Errorf("json.Marshal() =\n%s\nwant:\n%s", string(got), tt.wantJSON)
			}

			// Test round-trip deserialization
			var decoded RootCauseAnalysis
			if err := json.Unmarshal(got, &decoded); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}

			if decoded.OriginatingPackage != tt.input.OriginatingPackage {
				t.Errorf("OriginatingPackage = %v, want %v", decoded.OriginatingPackage, tt.input.OriginatingPackage)
			}
			if decoded.Confidence != tt.input.Confidence {
				t.Errorf("Confidence = %v, want %v", decoded.Confidence, tt.input.Confidence)
			}
			if len(decoded.Chain) != len(tt.input.Chain) {
				t.Errorf("Chain length = %v, want %v", len(decoded.Chain), len(tt.input.Chain))
			}
		})
	}
}

func TestRootCauseEdge_JSONSerialization(t *testing.T) {
	tests := []struct {
		name     string
		input    RootCauseEdge
		wantJSON string
	}{
		{
			name: "production dependency",
			input: RootCauseEdge{
				From:     "pkg-a",
				To:       "pkg-b",
				Type:     DependencyTypeProduction,
				Critical: false,
			},
			wantJSON: `{"from":"pkg-a","to":"pkg-b","type":"production","critical":false}`,
		},
		{
			name: "development dependency critical",
			input: RootCauseEdge{
				From:     "pkg-test",
				To:       "pkg-core",
				Type:     DependencyTypeDevelopment,
				Critical: true,
			},
			wantJSON: `{"from":"pkg-test","to":"pkg-core","type":"development","critical":true}`,
		},
		{
			name: "peer dependency",
			input: RootCauseEdge{
				From:     "pkg-plugin",
				To:       "pkg-host",
				Type:     DependencyTypePeer,
				Critical: false,
			},
			wantJSON: `{"from":"pkg-plugin","to":"pkg-host","type":"peer","critical":false}`,
		},
		{
			name: "optional dependency",
			input: RootCauseEdge{
				From:     "pkg-main",
				To:       "pkg-optional",
				Type:     DependencyTypeOptional,
				Critical: false,
			},
			wantJSON: `{"from":"pkg-main","to":"pkg-optional","type":"optional","critical":false}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}

			if string(got) != tt.wantJSON {
				t.Errorf("json.Marshal() = %s, want %s", string(got), tt.wantJSON)
			}

			// Verify round-trip
			var decoded RootCauseEdge
			if err := json.Unmarshal(got, &decoded); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}

			if decoded != tt.input {
				t.Errorf("Round-trip failed: got %+v, want %+v", decoded, tt.input)
			}
		})
	}
}

func TestRootCauseAnalysis_CamelCaseJSON(t *testing.T) {
	// Verify that all JSON keys are camelCase (project-context.md requirement)
	rca := &RootCauseAnalysis{
		OriginatingPackage: "test",
		ProblematicDependency: RootCauseEdge{
			From:     "a",
			To:       "b",
			Type:     DependencyTypeProduction,
			Critical: false,
		},
		Confidence:  50,
		Explanation: "test",
		Chain:       []RootCauseEdge{},
		CriticalEdge: &RootCauseEdge{
			From:     "c",
			To:       "d",
			Type:     DependencyTypeProduction,
			Critical: true,
		},
	}

	data, err := json.Marshal(rca)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	jsonStr := string(data)

	// Check for camelCase keys (should exist)
	camelCaseKeys := []string{
		`"originatingPackage"`,
		`"problematicDependency"`,
		`"confidence"`,
		`"explanation"`,
		`"chain"`,
		`"criticalEdge"`,
		`"from"`,
		`"to"`,
		`"type"`,
		`"critical"`,
	}
	for _, key := range camelCaseKeys {
		if !contains(jsonStr, key) {
			t.Errorf("Expected camelCase key %s not found in JSON: %s", key, jsonStr)
		}
	}

	// Check that snake_case keys do NOT exist
	snakeCaseKeys := []string{
		`"originating_package"`,
		`"problematic_dependency"`,
		`"critical_edge"`,
	}
	for _, key := range snakeCaseKeys {
		if contains(jsonStr, key) {
			t.Errorf("Unexpected snake_case key %s found in JSON: %s", key, jsonStr)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
