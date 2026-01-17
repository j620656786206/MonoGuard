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
		{1, 1},  // self-loop
		{2, 3},  // direct
		{3, 5},  // short indirect
		{4, 5},  // short indirect
		{5, 7},  // medium indirect
		{6, 7},  // medium indirect
		{7, 7},  // long indirect
		{8, 8},  // long indirect
		{15, 10}, // max complexity
	}

	for _, tt := range tests {
		result := calculateBaseComplexity(tt.depth)
		if result != tt.expected {
			t.Errorf("calculateBaseComplexity(%d) = %d, want %d", tt.depth, result, tt.expected)
		}
	}
}
