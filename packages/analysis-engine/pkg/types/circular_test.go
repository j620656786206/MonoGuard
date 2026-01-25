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

// ========================================
// FixStrategies Integration Tests (Story 3.3)
// ========================================

func TestCircularDependencyInfo_WithFixStrategies(t *testing.T) {
	// Test that FixStrategies field is optional and omitted when nil/empty
	info := &CircularDependencyInfo{
		Cycle:         []string{"A", "B", "A"},
		Type:          CircularTypeDirect,
		Severity:      CircularSeverityWarning,
		Depth:         2,
		Impact:        "Direct circular dependency between A and B",
		Complexity:    3,
		FixStrategies: nil, // Should be omitted in JSON
	}

	data, err := json.Marshal(info)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	jsonStr := string(data)

	// FixStrategies should NOT be in JSON when nil (omitempty)
	if strings.Contains(jsonStr, "fixStrategies") {
		t.Errorf("Expected fixStrategies to be omitted when nil, got: %s", jsonStr)
	}
}

func TestCircularDependencyInfo_WithFixStrategiesPresent(t *testing.T) {
	// Test that FixStrategies field is included when present
	info := &CircularDependencyInfo{
		Cycle:      []string{"@mono/ui", "@mono/api", "@mono/core", "@mono/ui"},
		Type:       CircularTypeIndirect,
		Severity:   CircularSeverityInfo,
		Depth:      3,
		Impact:     "Indirect circular dependency involving 3 packages",
		Complexity: 5,
		FixStrategies: []FixStrategy{
			{
				Type:           FixStrategyExtractModule,
				Name:           "Extract Shared Module",
				Description:    "Create a new shared package.",
				Suitability:    8,
				Effort:         EffortMedium,
				Pros:           []string{"Clean separation"},
				Cons:           []string{"New package"},
				Recommended:    true,
				TargetPackages: []string{"@mono/ui", "@mono/api", "@mono/core"},
				NewPackageName: "@mono/shared",
			},
			{
				Type:           FixStrategyDependencyInject,
				Name:           "Dependency Injection",
				Description:    "Invert the dependency.",
				Suitability:    6,
				Effort:         EffortLow,
				Pros:           []string{"Minimal changes"},
				Cons:           []string{"Adds indirection"},
				Recommended:    false,
				TargetPackages: []string{"@mono/core", "@mono/ui"},
			},
		},
	}

	data, err := json.Marshal(info)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	jsonStr := string(data)

	// FixStrategies should be in JSON when present
	expectedFields := []string{
		`"fixStrategies"`,
		`"extract-module"`,
		`"dependency-injection"`,
		`"suitability"`,
		`"effort"`,
		`"pros"`,
		`"cons"`,
		`"recommended"`,
		`"targetPackages"`,
		`"newPackageName"`,
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

	if len(decoded.FixStrategies) != 2 {
		t.Fatalf("FixStrategies length = %d, want 2", len(decoded.FixStrategies))
	}
	if decoded.FixStrategies[0].Type != FixStrategyExtractModule {
		t.Errorf("First strategy type = %s, want extract-module", decoded.FixStrategies[0].Type)
	}
	if !decoded.FixStrategies[0].Recommended {
		t.Error("First strategy should be recommended")
	}
	if decoded.FixStrategies[0].NewPackageName != "@mono/shared" {
		t.Errorf("NewPackageName = %s, want @mono/shared", decoded.FixStrategies[0].NewPackageName)
	}
}

// ========================================
// RefactoringComplexity Integration Tests (Story 3.5)
// ========================================

func TestCircularDependencyInfo_WithRefactoringComplexity(t *testing.T) {
	// Test that RefactoringComplexity field is optional and omitted when nil
	info := &CircularDependencyInfo{
		Cycle:                 []string{"A", "B", "A"},
		Type:                  CircularTypeDirect,
		Severity:              CircularSeverityWarning,
		Depth:                 2,
		Impact:                "Direct circular dependency between A and B",
		Complexity:            3,
		RefactoringComplexity: nil, // Should be omitted in JSON
	}

	data, err := json.Marshal(info)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	jsonStr := string(data)

	// RefactoringComplexity should NOT be in JSON when nil (omitempty)
	if strings.Contains(jsonStr, "refactoringComplexity") {
		t.Errorf("Expected refactoringComplexity to be omitted when nil, got: %s", jsonStr)
	}
}

func TestCircularDependencyInfo_WithRefactoringComplexityPresent(t *testing.T) {
	// Test that RefactoringComplexity field is included when present
	info := &CircularDependencyInfo{
		Cycle:      []string{"@mono/ui", "@mono/api", "@mono/core", "@mono/ui"},
		Type:       CircularTypeIndirect,
		Severity:   CircularSeverityInfo,
		Depth:      3,
		Impact:     "Indirect circular dependency involving 3 packages",
		Complexity: 5,
		RefactoringComplexity: &RefactoringComplexity{
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
		},
	}

	data, err := json.Marshal(info)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	jsonStr := string(data)

	// RefactoringComplexity should be in JSON when present
	expectedFields := []string{
		`"refactoringComplexity"`,
		`"score"`,
		`"estimatedTime"`,
		`"breakdown"`,
		`"filesAffected"`,
		`"importsToChange"`,
		`"chainDepth"`,
		`"packagesInvolved"`,
		`"externalDependencies"`,
		`"explanation"`,
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

	if decoded.RefactoringComplexity == nil {
		t.Fatal("RefactoringComplexity should not be nil after round-trip")
	}
	if decoded.RefactoringComplexity.Score != 5 {
		t.Errorf("Score = %d, want 5", decoded.RefactoringComplexity.Score)
	}
	if decoded.RefactoringComplexity.EstimatedTime != "30-60 minutes" {
		t.Errorf("EstimatedTime = %s, want 30-60 minutes", decoded.RefactoringComplexity.EstimatedTime)
	}
	if decoded.RefactoringComplexity.Breakdown.FilesAffected.Value != 3 {
		t.Errorf("FilesAffected.Value = %d, want 3", decoded.RefactoringComplexity.Breakdown.FilesAffected.Value)
	}
}

func TestCircularDependencyInfo_RefactoringComplexityBackwardCompatibility(t *testing.T) {
	// Test that existing JSON without refactoringComplexity still deserializes correctly
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

	// RefactoringComplexity should be nil for legacy JSON
	if decoded.RefactoringComplexity != nil {
		t.Error("RefactoringComplexity should be nil for legacy JSON without refactoringComplexity field")
	}

	// Legacy complexity field should still work
	if decoded.Complexity != 3 {
		t.Errorf("Complexity = %d, want 3", decoded.Complexity)
	}
}

// ========================================
// QuickFix and PriorityScore Tests (Story 3.8)
// ========================================

func TestCircularDependencyInfo_WithQuickFix(t *testing.T) {
	// Test that QuickFix field is optional and omitted when nil
	info := &CircularDependencyInfo{
		Cycle:         []string{"A", "B", "A"},
		Type:          CircularTypeDirect,
		Severity:      CircularSeverityWarning,
		Depth:         2,
		Impact:        "Direct circular dependency between A and B",
		Complexity:    3,
		QuickFix:      nil, // Should be omitted in JSON
		PriorityScore: 0,
	}

	data, err := json.Marshal(info)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	jsonStr := string(data)

	// QuickFix should NOT be in JSON when nil (omitempty)
	if strings.Contains(jsonStr, `"quickFix"`) {
		t.Errorf("Expected quickFix to be omitted when nil, got: %s", jsonStr)
	}

	// PriorityScore should always be present (not omitempty)
	if !strings.Contains(jsonStr, `"priorityScore"`) {
		t.Errorf("Expected priorityScore to be present, got: %s", jsonStr)
	}
}

func TestCircularDependencyInfo_WithQuickFixPresent(t *testing.T) {
	// Test that QuickFix field is included when present
	info := &CircularDependencyInfo{
		Cycle:      []string{"@mono/ui", "@mono/api", "@mono/core", "@mono/ui"},
		Type:       CircularTypeIndirect,
		Severity:   CircularSeverityInfo,
		Depth:      3,
		Impact:     "Indirect circular dependency involving 3 packages",
		Complexity: 5,
		QuickFix: &QuickFixSummary{
			StrategyType:  FixStrategyExtractModule,
			StrategyName:  "Extract Shared Module",
			Summary:       "Create new shared package '@mono/shared' to break the cycle",
			Suitability:   8,
			Effort:        EffortMedium,
			EstimatedTime: "30-60 minutes",
			StrategyIndex: 0,
		},
		PriorityScore: 80.0,
	}

	data, err := json.Marshal(info)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	jsonStr := string(data)

	// QuickFix should be in JSON when present
	expectedFields := []string{
		`"quickFix"`,
		`"strategyType"`,
		`"strategyName"`,
		`"summary"`,
		`"suitability"`,
		`"effort"`,
		`"estimatedTime"`,
		`"strategyIndex"`,
		`"priorityScore"`,
	}

	for _, field := range expectedFields {
		if !strings.Contains(jsonStr, field) {
			t.Errorf("Expected JSON to contain %s, got: %s", field, jsonStr)
		}
	}

	// Verify priorityScore value
	if !strings.Contains(jsonStr, `"priorityScore":80`) {
		t.Errorf("Expected priorityScore:80, got: %s", jsonStr)
	}

	// Verify round-trip
	var decoded CircularDependencyInfo
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if decoded.QuickFix == nil {
		t.Fatal("QuickFix should not be nil after round-trip")
	}
	if decoded.QuickFix.StrategyType != FixStrategyExtractModule {
		t.Errorf("StrategyType = %s, want extract-module", decoded.QuickFix.StrategyType)
	}
	if decoded.QuickFix.Suitability != 8 {
		t.Errorf("Suitability = %d, want 8", decoded.QuickFix.Suitability)
	}
	if decoded.PriorityScore != 80.0 {
		t.Errorf("PriorityScore = %f, want 80.0", decoded.PriorityScore)
	}
}

func TestCircularDependencyInfo_QuickFixBackwardCompatibility(t *testing.T) {
	// Test that existing JSON without quickFix/priorityScore still deserializes correctly
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

	// QuickFix should be nil for legacy JSON
	if decoded.QuickFix != nil {
		t.Error("QuickFix should be nil for legacy JSON without quickFix field")
	}

	// PriorityScore should default to 0
	if decoded.PriorityScore != 0 {
		t.Errorf("PriorityScore should be 0 for legacy JSON, got %f", decoded.PriorityScore)
	}

	// Other fields should be correct
	if decoded.Type != CircularTypeDirect {
		t.Errorf("Type = %q, want %q", decoded.Type, CircularTypeDirect)
	}
}

func TestCircularDependencyInfo_PriorityScoreValues(t *testing.T) {
	tests := []struct {
		name          string
		priorityScore float64
	}{
		{"zero", 0},
		{"low", 10.5},
		{"medium", 45.0},
		{"high", 80.0},
		{"maximum", 100.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &CircularDependencyInfo{
				Cycle:         []string{"A", "B", "A"},
				Type:          CircularTypeDirect,
				Severity:      CircularSeverityWarning,
				Depth:         2,
				Impact:        "Test",
				Complexity:    3,
				PriorityScore: tt.priorityScore,
			}

			data, err := json.Marshal(info)
			if err != nil {
				t.Fatalf("Failed to marshal: %v", err)
			}

			var decoded CircularDependencyInfo
			if err := json.Unmarshal(data, &decoded); err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			if decoded.PriorityScore != tt.priorityScore {
				t.Errorf("PriorityScore = %f, want %f", decoded.PriorityScore, tt.priorityScore)
			}
		})
	}
}
