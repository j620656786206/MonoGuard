package types

import (
	"encoding/json"
	"testing"
)

func TestVersionConflictInfo_JSONSerialization(t *testing.T) {
	conflict := &VersionConflictInfo{
		PackageName: "lodash",
		ConflictingVersions: []*ConflictingVersion{
			{
				Version:    "^4.17.21",
				Packages:   []string{"@mono/app"},
				IsBreaking: false,
				DepType:    DepTypeProduction,
			},
			{
				Version:    "^4.17.19",
				Packages:   []string{"@mono/utils"},
				IsBreaking: false,
				DepType:    DepTypeProduction,
			},
		},
		Severity:   ConflictSeverityInfo,
		Resolution: "Patch version difference. Safe to upgrade all packages to ^4.17.21",
		Impact:     "Minor bundle size increase. No breaking changes expected.",
	}

	// Test serialization
	data, err := json.Marshal(conflict)
	if err != nil {
		t.Fatalf("Failed to marshal VersionConflictInfo: %v", err)
	}

	// Verify camelCase JSON keys
	jsonStr := string(data)
	expectedKeys := []string{
		`"packageName"`,
		`"conflictingVersions"`,
		`"severity"`,
		`"resolution"`,
		`"impact"`,
		`"version"`,
		`"packages"`,
		`"isBreaking"`,
		`"depType"`,
	}

	for _, key := range expectedKeys {
		if !contains(jsonStr, key) {
			t.Errorf("Expected JSON to contain key %s, but it didn't. JSON: %s", key, jsonStr)
		}
	}

	// Test deserialization
	var decoded VersionConflictInfo
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal VersionConflictInfo: %v", err)
	}

	if decoded.PackageName != conflict.PackageName {
		t.Errorf("PackageName mismatch: got %s, want %s", decoded.PackageName, conflict.PackageName)
	}

	if len(decoded.ConflictingVersions) != len(conflict.ConflictingVersions) {
		t.Errorf("ConflictingVersions length mismatch: got %d, want %d",
			len(decoded.ConflictingVersions), len(conflict.ConflictingVersions))
	}

	if decoded.Severity != conflict.Severity {
		t.Errorf("Severity mismatch: got %s, want %s", decoded.Severity, conflict.Severity)
	}
}

func TestConflictSeverity_Values(t *testing.T) {
	tests := []struct {
		severity ConflictSeverity
		expected string
	}{
		{ConflictSeverityCritical, "critical"},
		{ConflictSeverityWarning, "warning"},
		{ConflictSeverityInfo, "info"},
	}

	for _, tt := range tests {
		t.Run(string(tt.severity), func(t *testing.T) {
			if string(tt.severity) != tt.expected {
				t.Errorf("ConflictSeverity value mismatch: got %s, want %s", tt.severity, tt.expected)
			}
		})
	}
}

func TestConflictingVersion_JSONSerialization(t *testing.T) {
	cv := &ConflictingVersion{
		Version:    "^5.0.0",
		Packages:   []string{"@mono/app", "@mono/lib"},
		IsBreaking: true,
		DepType:    DepTypeDevelopment,
	}

	data, err := json.Marshal(cv)
	if err != nil {
		t.Fatalf("Failed to marshal ConflictingVersion: %v", err)
	}

	var decoded ConflictingVersion
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal ConflictingVersion: %v", err)
	}

	if decoded.Version != cv.Version {
		t.Errorf("Version mismatch: got %s, want %s", decoded.Version, cv.Version)
	}

	if decoded.IsBreaking != cv.IsBreaking {
		t.Errorf("IsBreaking mismatch: got %v, want %v", decoded.IsBreaking, cv.IsBreaking)
	}

	if decoded.DepType != cv.DepType {
		t.Errorf("DepType mismatch: got %s, want %s", decoded.DepType, cv.DepType)
	}

	if len(decoded.Packages) != len(cv.Packages) {
		t.Errorf("Packages length mismatch: got %d, want %d", len(decoded.Packages), len(cv.Packages))
	}
}

func TestDepTypeConstants(t *testing.T) {
	if DepTypeProduction != "production" {
		t.Errorf("DepTypeProduction: got %s, want production", DepTypeProduction)
	}
	if DepTypeDevelopment != "development" {
		t.Errorf("DepTypeDevelopment: got %s, want development", DepTypeDevelopment)
	}
	if DepTypePeer != "peer" {
		t.Errorf("DepTypePeer: got %s, want peer", DepTypePeer)
	}
}

func TestVersionConflictInfo_CriticalSeverity(t *testing.T) {
	conflict := &VersionConflictInfo{
		PackageName: "typescript",
		ConflictingVersions: []*ConflictingVersion{
			{
				Version:    "^5.0.0",
				Packages:   []string{"@mono/app"},
				IsBreaking: true,
				DepType:    DepTypeDevelopment,
			},
			{
				Version:    "^4.9.0",
				Packages:   []string{"@mono/utils"},
				IsBreaking: false,
				DepType:    DepTypeDevelopment,
			},
		},
		Severity:   ConflictSeverityCritical,
		Resolution: "Major version conflict detected. Review breaking changes before upgrading all packages to ^5.0.0",
		Impact:     "TypeScript 5.x has breaking changes. Check compatibility before upgrading.",
	}

	data, err := json.Marshal(conflict)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Verify severity is "critical"
	if !contains(string(data), `"severity":"critical"`) {
		t.Errorf("Expected severity to be 'critical' in JSON: %s", string(data))
	}

	// Verify isBreaking is true for the first version
	if !contains(string(data), `"isBreaking":true`) {
		t.Errorf("Expected isBreaking:true in JSON: %s", string(data))
	}
}

// Helper function
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
