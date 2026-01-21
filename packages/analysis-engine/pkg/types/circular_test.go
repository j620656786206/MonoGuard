package types

import (
	"encoding/json"
	"strings"
	"testing"
)

// ========================================
// Circular Type Constants Tests
// ========================================

func TestCircularTypeConstants(t *testing.T) {
	tests := []struct {
		circularType CircularType
		expected     string
	}{
		{CircularTypeDirect, "direct"},
		{CircularTypeIndirect, "indirect"},
	}

	for _, tt := range tests {
		t.Run(string(tt.circularType), func(t *testing.T) {
			if string(tt.circularType) != tt.expected {
				t.Errorf("CircularType = %q, want %q", tt.circularType, tt.expected)
			}
		})
	}
}

func TestCircularSeverityConstants(t *testing.T) {
	tests := []struct {
		severity CircularSeverity
		expected string
	}{
		{CircularSeverityCritical, "critical"},
		{CircularSeverityWarning, "warning"},
		{CircularSeverityInfo, "info"},
	}

	for _, tt := range tests {
		t.Run(string(tt.severity), func(t *testing.T) {
			if string(tt.severity) != tt.expected {
				t.Errorf("CircularSeverity = %q, want %q", tt.severity, tt.expected)
			}
		})
	}
}

// ========================================
// CircularDependencyInfo JSON Tests
// ========================================

func TestCircularDependencyInfoJSONSerialization(t *testing.T) {
	info := &CircularDependencyInfo{
		Cycle:      []string{"A", "B", "C", "A"},
		Type:       CircularTypeIndirect,
		Severity:   CircularSeverityInfo,
		Depth:      3,
		Impact:     "Indirect circular dependency involving 3 packages",
		Complexity: 5,
	}

	data, err := json.Marshal(info)
	if err != nil {
		t.Fatalf("Failed to marshal CircularDependencyInfo: %v", err)
	}

	jsonStr := string(data)

	// Verify camelCase field names
	expectedFields := []string{
		`"cycle"`,
		`"type"`,
		`"severity"`,
		`"depth"`,
		`"impact"`,
		`"complexity"`,
	}

	for _, field := range expectedFields {
		if !strings.Contains(jsonStr, field) {
			t.Errorf("JSON missing field %s, got: %s", field, jsonStr)
		}
	}

	// Verify NO snake_case
	snakeCaseFields := []string{
		`"circular_type"`,
		`"circular_severity"`,
	}

	for _, field := range snakeCaseFields {
		if strings.Contains(jsonStr, field) {
			t.Errorf("JSON contains snake_case field %s (should be camelCase)", field)
		}
	}
}

func TestCircularDependencyInfoRoundTrip(t *testing.T) {
	original := &CircularDependencyInfo{
		Cycle:      []string{"@mono/a", "@mono/b", "@mono/c", "@mono/a"},
		Type:       CircularTypeIndirect,
		Severity:   CircularSeverityWarning,
		Depth:      3,
		Impact:     "Indirect circular dependency involving 3 packages",
		Complexity: 5,
	}

	// Marshal
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Unmarshal
	var decoded CircularDependencyInfo
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Verify fields
	if len(decoded.Cycle) != len(original.Cycle) {
		t.Errorf("Cycle length = %d, want %d", len(decoded.Cycle), len(original.Cycle))
	}
	if decoded.Type != original.Type {
		t.Errorf("Type = %q, want %q", decoded.Type, original.Type)
	}
	if decoded.Severity != original.Severity {
		t.Errorf("Severity = %q, want %q", decoded.Severity, original.Severity)
	}
	if decoded.Depth != original.Depth {
		t.Errorf("Depth = %d, want %d", decoded.Depth, original.Depth)
	}
	if decoded.Complexity != original.Complexity {
		t.Errorf("Complexity = %d, want %d", decoded.Complexity, original.Complexity)
	}
}

// ========================================
// NewCircularDependencyInfo Tests
// ========================================

func TestNewCircularDependencyInfo_SelfLoop(t *testing.T) {
	info := NewCircularDependencyInfo([]string{"A", "A"})

	if info == nil {
		t.Fatal("Expected non-nil result for self-loop")
	}
	if info.Depth != 1 {
		t.Errorf("Depth = %d, want 1", info.Depth)
	}
	if info.Type != CircularTypeDirect {
		t.Errorf("Type = %q, want %q", info.Type, CircularTypeDirect)
	}
	if info.Severity != CircularSeverityCritical {
		t.Errorf("Severity = %q, want %q", info.Severity, CircularSeverityCritical)
	}
	if info.Complexity != 1 {
		t.Errorf("Complexity = %d, want 1", info.Complexity)
	}
	if !strings.Contains(info.Impact, "Self-referencing") {
		t.Errorf("Impact should mention self-referencing, got: %s", info.Impact)
	}
}

func TestNewCircularDependencyInfo_DirectCycle(t *testing.T) {
	info := NewCircularDependencyInfo([]string{"A", "B", "A"})

	if info == nil {
		t.Fatal("Expected non-nil result for direct cycle")
	}
	if info.Depth != 2 {
		t.Errorf("Depth = %d, want 2", info.Depth)
	}
	if info.Type != CircularTypeDirect {
		t.Errorf("Type = %q, want %q", info.Type, CircularTypeDirect)
	}
	if info.Severity != CircularSeverityWarning {
		t.Errorf("Severity = %q, want %q", info.Severity, CircularSeverityWarning)
	}
	if info.Complexity != 3 {
		t.Errorf("Complexity = %d, want 3", info.Complexity)
	}
	if !strings.Contains(info.Impact, "Direct circular dependency") {
		t.Errorf("Impact should mention direct, got: %s", info.Impact)
	}
}

func TestNewCircularDependencyInfo_IndirectCycle(t *testing.T) {
	info := NewCircularDependencyInfo([]string{"A", "B", "C", "A"})

	if info == nil {
		t.Fatal("Expected non-nil result for indirect cycle")
	}
	if info.Depth != 3 {
		t.Errorf("Depth = %d, want 3", info.Depth)
	}
	if info.Type != CircularTypeIndirect {
		t.Errorf("Type = %q, want %q", info.Type, CircularTypeIndirect)
	}
	if info.Severity != CircularSeverityInfo {
		t.Errorf("Severity = %q, want %q", info.Severity, CircularSeverityInfo)
	}
	if info.Complexity != 5 {
		t.Errorf("Complexity = %d, want 5", info.Complexity)
	}
	if !strings.Contains(info.Impact, "Indirect") && !strings.Contains(info.Impact, "3 packages") {
		t.Errorf("Impact should mention indirect and 3 packages, got: %s", info.Impact)
	}
}

func TestNewCircularDependencyInfo_LongCycle(t *testing.T) {
	// 7 unique packages
	info := NewCircularDependencyInfo([]string{"A", "B", "C", "D", "E", "F", "G", "A"})

	if info == nil {
		t.Fatal("Expected non-nil result for long cycle")
	}
	if info.Depth != 7 {
		t.Errorf("Depth = %d, want 7", info.Depth)
	}
	if info.Complexity != 7 {
		t.Errorf("Complexity = %d, want 7", info.Complexity)
	}
}

func TestNewCircularDependencyInfo_EmptyCycle(t *testing.T) {
	info := NewCircularDependencyInfo([]string{})

	if info != nil {
		t.Error("Expected nil result for empty cycle")
	}
}

// ========================================
// classifySeverity Tests
// ========================================

func TestClassifySeverity(t *testing.T) {
	tests := []struct {
		name     string
		cycle    []string
		expected CircularSeverity
	}{
		{"self-loop", []string{"A", "A"}, CircularSeverityCritical},
		{"direct", []string{"A", "B", "A"}, CircularSeverityWarning},
		{"indirect 3", []string{"A", "B", "C", "A"}, CircularSeverityInfo},
		{"indirect 5", []string{"A", "B", "C", "D", "E", "A"}, CircularSeverityInfo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cycleType := CircularTypeIndirect
			if len(tt.cycle) <= 3 {
				cycleType = CircularTypeDirect
			}
			result := classifySeverity(tt.cycle, cycleType)
			if result != tt.expected {
				t.Errorf("classifySeverity() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// ========================================
// generateImpactDescription Tests
// ========================================

func TestGenerateImpactDescription(t *testing.T) {
	tests := []struct {
		name     string
		cycle    []string
		contains string
	}{
		{"empty", []string{}, ""},
		{"self-loop", []string{"A", "A"}, "Self-referencing"},
		{"direct", []string{"A", "B", "A"}, "Direct circular"},
		{"indirect", []string{"A", "B", "C", "A"}, "Indirect"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateImpactDescription(tt.cycle)
			if tt.contains == "" && result != "" {
				t.Errorf("Expected empty string, got %q", result)
			} else if tt.contains != "" && !strings.Contains(result, tt.contains) {
				t.Errorf("Expected result to contain %q, got %q", tt.contains, result)
			}
		})
	}
}

// ========================================
// calculateBaseComplexity Tests
// ========================================

func TestCalculateBaseComplexity(t *testing.T) {
	tests := []struct {
		depth    int
		expected int
	}{
		{1, ComplexitySelfLoop},        // self-loop
		{2, ComplexityDirect},          // direct
		{3, ComplexityShortIndirect},   // short indirect
		{4, ComplexityShortIndirect},   // short indirect
		{5, ComplexityMediumIndirect},  // medium indirect
		{6, ComplexityMediumIndirect},  // medium indirect
		{7, 7},                         // long indirect (equals depth)
		{8, 8},                         // long indirect (equals depth)
		{15, ComplexityMax},            // max complexity (capped)
	}

	for _, tt := range tests {
		result := calculateBaseComplexity(tt.depth)
		if result != tt.expected {
			t.Errorf("calculateBaseComplexity(%d) = %d, want %d", tt.depth, result, tt.expected)
		}
	}
}

func TestComplexityConstants(t *testing.T) {
	// Verify constants have expected values
	if ComplexitySelfLoop != 1 {
		t.Errorf("ComplexitySelfLoop = %d, want 1", ComplexitySelfLoop)
	}
	if ComplexityDirect != 3 {
		t.Errorf("ComplexityDirect = %d, want 3", ComplexityDirect)
	}
	if ComplexityShortIndirect != 5 {
		t.Errorf("ComplexityShortIndirect = %d, want 5", ComplexityShortIndirect)
	}
	if ComplexityMediumIndirect != 7 {
		t.Errorf("ComplexityMediumIndirect = %d, want 7", ComplexityMediumIndirect)
	}
	if ComplexityMax != 10 {
		t.Errorf("ComplexityMax = %d, want 10", ComplexityMax)
	}
}

// ========================================
// RootCause Integration Tests (Story 3.1)
// ========================================

func TestCircularDependencyInfo_WithRootCause(t *testing.T) {
	// Test that RootCause field is optional and omitted when nil
	info := &CircularDependencyInfo{
		Cycle:      []string{"A", "B", "A"},
		Type:       CircularTypeDirect,
		Severity:   CircularSeverityWarning,
		Depth:      2,
		Impact:     "Direct circular dependency between A and B",
		Complexity: 3,
		RootCause:  nil, // Should be omitted in JSON
	}

	data, err := json.Marshal(info)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	jsonStr := string(data)

	// RootCause should NOT be in JSON when nil (omitempty)
	if strings.Contains(jsonStr, "rootCause") {
		t.Errorf("Expected rootCause to be omitted when nil, got: %s", jsonStr)
	}
}

func TestCircularDependencyInfo_WithRootCausePresent(t *testing.T) {
	// Test that RootCause field is included when present
	info := &CircularDependencyInfo{
		Cycle:      []string{"pkg-ui", "pkg-api", "pkg-core", "pkg-ui"},
		Type:       CircularTypeIndirect,
		Severity:   CircularSeverityInfo,
		Depth:      3,
		Impact:     "Indirect circular dependency involving 3 packages",
		Complexity: 5,
		RootCause: &RootCauseAnalysis{
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
	}

	data, err := json.Marshal(info)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	jsonStr := string(data)

	// RootCause should be in JSON when present
	expectedFields := []string{
		`"rootCause"`,
		`"originatingPackage"`,
		`"problematicDependency"`,
		`"confidence"`,
		`"explanation"`,
		`"chain"`,
		`"criticalEdge"`,
	}

	for _, field := range expectedFields {
		if !strings.Contains(jsonStr, field) {
			t.Errorf("Expected JSON to contain %s, got: %s", field, jsonStr)
		}
	}

	// Verify round-trip
	var decoded CircularDependencyInfo
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if decoded.RootCause == nil {
		t.Fatal("RootCause should not be nil after round-trip")
	}
	if decoded.RootCause.OriginatingPackage != "pkg-ui" {
		t.Errorf("OriginatingPackage = %s, want pkg-ui", decoded.RootCause.OriginatingPackage)
	}
	if decoded.RootCause.Confidence != 82 {
		t.Errorf("Confidence = %d, want 82", decoded.RootCause.Confidence)
	}
}

func TestCircularDependencyInfo_BackwardCompatibility(t *testing.T) {
	// Test that existing JSON without rootCause still deserializes correctly
	jsonStr := `{
		"cycle": ["A", "B", "A"],
		"type": "direct",
		"severity": "warning",
		"depth": 2,
		"impact": "Direct circular dependency between A and B",
		"complexity": 3
	}`

	var decoded CircularDependencyInfo
	if err := json.Unmarshal([]byte(jsonStr), &decoded); err != nil {
		t.Fatalf("Failed to unmarshal legacy JSON: %v", err)
	}

	// RootCause should be nil for legacy JSON
	if decoded.RootCause != nil {
		t.Error("RootCause should be nil for legacy JSON without rootCause field")
	}

	// Other fields should be correct
	if decoded.Type != CircularTypeDirect {
		t.Errorf("Type = %q, want %q", decoded.Type, CircularTypeDirect)
	}
	if decoded.Severity != CircularSeverityWarning {
		t.Errorf("Severity = %q, want %q", decoded.Severity, CircularSeverityWarning)
	}
}
