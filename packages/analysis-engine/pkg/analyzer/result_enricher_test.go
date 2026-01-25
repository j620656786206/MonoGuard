// Package analyzer provides dependency graph analysis for monorepo workspaces.
// This file contains tests for the result enricher (Story 3.8).
package analyzer

import (
	"testing"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

func TestSortStrategies(t *testing.T) {
	tests := []struct {
		name       string
		strategies []types.FixStrategy
		wantOrder  []int // Expected suitability order
	}{
		{
			name: "sorts by suitability descending",
			strategies: []types.FixStrategy{
				{Name: "Low", Suitability: 3},
				{Name: "High", Suitability: 9},
				{Name: "Medium", Suitability: 6},
			},
			wantOrder: []int{9, 6, 3},
		},
		{
			name: "stable sort preserves order for equal values",
			strategies: []types.FixStrategy{
				{Name: "First", Suitability: 5},
				{Name: "Second", Suitability: 5},
				{Name: "Third", Suitability: 5},
			},
			wantOrder: []int{5, 5, 5},
		},
		{
			name:       "empty slice",
			strategies: []types.FixStrategy{},
			wantOrder:  []int{},
		},
		{
			name: "single element",
			strategies: []types.FixStrategy{
				{Name: "Only", Suitability: 7},
			},
			wantOrder: []int{7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortStrategies(tt.strategies)

			for i, want := range tt.wantOrder {
				if tt.strategies[i].Suitability != want {
					t.Errorf("Position %d: expected suitability %d, got %d",
						i, want, tt.strategies[i].Suitability)
				}
			}
		})
	}
}

func TestCalculatePriorityScore(t *testing.T) {
	tests := []struct {
		name      string
		cycle     *types.CircularDependencyInfo
		wantScore float64
	}{
		{
			name: "critical impact + low complexity = highest priority",
			cycle: &types.CircularDependencyInfo{
				Cycle: []string{"a", "b", "a"},
				ImpactAssessment: &types.ImpactAssessment{
					RiskLevel: types.RiskLevelCritical,
				},
				RefactoringComplexity: &types.RefactoringComplexity{
					Score: 2, // Ease = 11 - 2 = 9
				},
			},
			wantScore: 90.0, // 10 * 9
		},
		{
			name: "low impact + high complexity = lowest priority",
			cycle: &types.CircularDependencyInfo{
				Cycle: []string{"a", "b", "a"},
				ImpactAssessment: &types.ImpactAssessment{
					RiskLevel: types.RiskLevelLow,
				},
				RefactoringComplexity: &types.RefactoringComplexity{
					Score: 9, // Ease = 11 - 9 = 2
				},
			},
			wantScore: 5.0, // 2.5 * 2
		},
		{
			name: "high impact + medium complexity",
			cycle: &types.CircularDependencyInfo{
				Cycle: []string{"a", "b", "a"},
				ImpactAssessment: &types.ImpactAssessment{
					RiskLevel: types.RiskLevelHigh,
				},
				RefactoringComplexity: &types.RefactoringComplexity{
					Score: 5, // Ease = 11 - 5 = 6
				},
			},
			wantScore: 45.0, // 7.5 * 6
		},
		{
			name: "medium impact + no complexity (default)",
			cycle: &types.CircularDependencyInfo{
				Cycle: []string{"a", "b", "a"},
				ImpactAssessment: &types.ImpactAssessment{
					RiskLevel: types.RiskLevelMedium,
				},
				RefactoringComplexity: nil,
			},
			wantScore: 30.0, // 5.0 * 6.0 (default ease)
		},
		{
			name: "no impact assessment (default)",
			cycle: &types.CircularDependencyInfo{
				Cycle:            []string{"a", "b", "a"},
				ImpactAssessment: nil,
				RefactoringComplexity: &types.RefactoringComplexity{
					Score: 4, // Ease = 11 - 4 = 7
				},
			},
			wantScore: 35.0, // 5.0 (default impact) * 7
		},
		{
			name: "no assessment data (all defaults)",
			cycle: &types.CircularDependencyInfo{
				Cycle:                 []string{"a", "b", "a"},
				ImpactAssessment:      nil,
				RefactoringComplexity: nil,
			},
			wantScore: 30.0, // 5.0 * 6.0
		},
		{
			name: "complexity 10 = minimum ease 1",
			cycle: &types.CircularDependencyInfo{
				Cycle: []string{"a", "b", "a"},
				ImpactAssessment: &types.ImpactAssessment{
					RiskLevel: types.RiskLevelMedium,
				},
				RefactoringComplexity: &types.RefactoringComplexity{
					Score: 10, // Ease = max(11 - 10, 1) = 1
				},
			},
			wantScore: 5.0, // 5.0 * 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculatePriorityScore(tt.cycle)
			if got != tt.wantScore {
				t.Errorf("calculatePriorityScore() = %v, want %v", got, tt.wantScore)
			}
		})
	}
}

func TestSortCircularDependencies(t *testing.T) {
	cycles := []*types.CircularDependencyInfo{
		{Cycle: []string{"low", "x", "low"}, PriorityScore: 10.0},
		{Cycle: []string{"high", "y", "high"}, PriorityScore: 90.0},
		{Cycle: []string{"med", "z", "med"}, PriorityScore: 45.0},
	}

	sortCircularDependencies(cycles)

	expectedOrder := []float64{90.0, 45.0, 10.0}
	for i, want := range expectedOrder {
		if cycles[i].PriorityScore != want {
			t.Errorf("Position %d: expected priority %v, got %v",
				i, want, cycles[i].PriorityScore)
		}
	}
}

func TestCreateQuickFix(t *testing.T) {
	tests := []struct {
		name       string
		strategies []types.FixStrategy
		want       *types.QuickFixSummary
	}{
		{
			name:       "empty strategies returns nil",
			strategies: []types.FixStrategy{},
			want:       nil,
		},
		{
			name: "extract module strategy",
			strategies: []types.FixStrategy{
				{
					Type:           types.FixStrategyExtractModule,
					Name:           "Extract Shared Module",
					Description:    "Move shared code to new package",
					Suitability:    8,
					Effort:         types.EffortMedium,
					NewPackageName: "@mono/shared",
					Guide: &types.FixGuide{
						EstimatedTime: "30-60 minutes",
					},
				},
			},
			want: &types.QuickFixSummary{
				StrategyType:  types.FixStrategyExtractModule,
				StrategyName:  "Extract Shared Module",
				Summary:       "Create new shared package '@mono/shared' to break the cycle",
				Suitability:   8,
				Effort:        types.EffortMedium,
				EstimatedTime: "30-60 minutes",
				StrategyIndex: 0,
			},
		},
		{
			name: "dependency injection strategy",
			strategies: []types.FixStrategy{
				{
					Type:        types.FixStrategyDependencyInject,
					Name:        "Dependency Injection",
					Description: "Invert using DI",
					Suitability: 7,
					Effort:      types.EffortHigh,
				},
			},
			want: &types.QuickFixSummary{
				StrategyType:  types.FixStrategyDependencyInject,
				StrategyName:  "Dependency Injection",
				Summary:       "Invert dependency using dependency injection pattern",
				Suitability:   7,
				Effort:        types.EffortHigh,
				EstimatedTime: "15-30 minutes",
				StrategyIndex: 0,
			},
		},
		{
			name: "boundary refactoring strategy",
			strategies: []types.FixStrategy{
				{
					Type:        types.FixStrategyBoundaryRefactor,
					Name:        "Boundary Refactoring",
					Suitability: 6,
					Effort:      types.EffortHigh,
				},
			},
			want: &types.QuickFixSummary{
				StrategyType:  types.FixStrategyBoundaryRefactor,
				StrategyName:  "Boundary Refactoring",
				Summary:       "Restructure package boundaries to eliminate overlap",
				Suitability:   6,
				Effort:        types.EffortHigh,
				EstimatedTime: "15-30 minutes",
				StrategyIndex: 0,
			},
		},
		{
			name: "uses complexity estimated time as fallback",
			strategies: []types.FixStrategy{
				{
					Type:        types.FixStrategyExtractModule,
					Name:        "Extract",
					Suitability: 8,
					Effort:      types.EffortMedium,
					Guide:       nil,
					Complexity: &types.RefactoringComplexity{
						EstimatedTime: "1-2 hours",
					},
				},
			},
			want: &types.QuickFixSummary{
				StrategyType:  types.FixStrategyExtractModule,
				StrategyName:  "Extract",
				Summary:       "Extract shared code into a new package to break the cycle",
				Suitability:   8,
				Effort:        types.EffortMedium,
				EstimatedTime: "1-2 hours",
				StrategyIndex: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createQuickFix(tt.strategies)

			if tt.want == nil {
				if got != nil {
					t.Error("Expected nil, got non-nil")
				}
				return
			}

			if got == nil {
				t.Fatal("Expected non-nil, got nil")
			}

			if got.StrategyType != tt.want.StrategyType {
				t.Errorf("StrategyType = %q, want %q", got.StrategyType, tt.want.StrategyType)
			}
			if got.StrategyName != tt.want.StrategyName {
				t.Errorf("StrategyName = %q, want %q", got.StrategyName, tt.want.StrategyName)
			}
			if got.Summary != tt.want.Summary {
				t.Errorf("Summary = %q, want %q", got.Summary, tt.want.Summary)
			}
			if got.Suitability != tt.want.Suitability {
				t.Errorf("Suitability = %d, want %d", got.Suitability, tt.want.Suitability)
			}
			if got.Effort != tt.want.Effort {
				t.Errorf("Effort = %q, want %q", got.Effort, tt.want.Effort)
			}
			if got.EstimatedTime != tt.want.EstimatedTime {
				t.Errorf("EstimatedTime = %q, want %q", got.EstimatedTime, tt.want.EstimatedTime)
			}
			if got.StrategyIndex != tt.want.StrategyIndex {
				t.Errorf("StrategyIndex = %d, want %d", got.StrategyIndex, tt.want.StrategyIndex)
			}
		})
	}
}

func TestParseEstimatedMinutes(t *testing.T) {
	tests := []struct {
		timeStr string
		want    int
	}{
		{"15-30 minutes", 22},
		{"30-60 minutes", 45},
		{"5-15 minutes", 10},
		{"1-2 hours", 90},
		{"2-4 hours", 180},
		{"30 minutes", 30},
		{"1 hour", 60},
		{"2 hours", 120},
		{"", 30},                // Default
		{"unknown", 30},        // Default
		{"15–30 minutes", 22},  // En-dash
		{"1-2 HOURS", 90},      // Case insensitive
	}

	for _, tt := range tests {
		t.Run(tt.timeStr, func(t *testing.T) {
			got := parseEstimatedMinutes(tt.timeStr)
			if got != tt.want {
				t.Errorf("parseEstimatedMinutes(%q) = %d, want %d", tt.timeStr, got, tt.want)
			}
		})
	}
}

func TestFormatTotalTime(t *testing.T) {
	tests := []struct {
		minutes int
		want    string
	}{
		{30, "30 minutes"},
		{45, "45 minutes"},
		{60, "1 hour"},
		{90, "1 hour 30 minutes"},
		{120, "2 hours"},
		{150, "2 hours 30 minutes"},
		{0, "0 minutes"},
		{1, "1 minutes"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := formatTotalTime(tt.minutes)
			if got != tt.want {
				t.Errorf("formatTotalTime(%d) = %q, want %q", tt.minutes, got, tt.want)
			}
		})
	}
}

func TestGenerateCycleID(t *testing.T) {
	tests := []struct {
		cycle []string
		want  string
	}{
		{[]string{"@mono/core", "@mono/ui", "@mono/core"}, "core→ui"},
		{[]string{"packages/auth", "packages/user", "packages/auth"}, "auth→user"},
		{[]string{"a", "b", "a"}, "a→b"},
		{[]string{"a"}, "unknown"},
		{[]string{}, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := generateCycleID(tt.cycle)
			if got != tt.want {
				t.Errorf("generateCycleID(%v) = %q, want %q", tt.cycle, got, tt.want)
			}
		})
	}
}

func TestGetUniquePackages(t *testing.T) {
	tests := []struct {
		name  string
		cycle []string
		want  []string
	}{
		{
			name:  "closed cycle excludes duplicate",
			cycle: []string{"a", "b", "c", "a"},
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "open cycle returns as-is",
			cycle: []string{"a", "b", "c"},
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "empty returns empty",
			cycle: []string{},
			want:  []string{},
		},
		{
			name:  "self-loop",
			cycle: []string{"a", "a"},
			want:  []string{"a"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getUniquePackages(tt.cycle)
			if len(got) != len(tt.want) {
				t.Errorf("getUniquePackages() length = %d, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("getUniquePackages()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestResultEnricher_Enrich(t *testing.T) {
	t.Run("nil result returns nil", func(t *testing.T) {
		enricher := NewResultEnricher(nil, nil)
		got := enricher.Enrich(nil)
		if got != nil {
			t.Error("Expected nil for nil input")
		}
	})

	t.Run("empty cycles returns unchanged", func(t *testing.T) {
		enricher := NewResultEnricher(nil, nil)
		result := &types.AnalysisResult{
			CircularDependencies: []*types.CircularDependencyInfo{},
		}
		got := enricher.Enrich(result)
		if got.FixSummary != nil {
			t.Error("Expected nil FixSummary for empty cycles")
		}
	})

	t.Run("enriches cycles with QuickFix and sorts by priority", func(t *testing.T) {
		enricher := NewResultEnricher(nil, nil)

		result := &types.AnalysisResult{
			CircularDependencies: []*types.CircularDependencyInfo{
				{
					Cycle: []string{"low", "x", "low"},
					ImpactAssessment: &types.ImpactAssessment{
						RiskLevel: types.RiskLevelLow,
					},
					RefactoringComplexity: &types.RefactoringComplexity{
						Score: 8,
					},
					FixStrategies: []types.FixStrategy{
						{Type: types.FixStrategyExtractModule, Suitability: 5},
					},
				},
				{
					Cycle: []string{"high", "y", "high"},
					ImpactAssessment: &types.ImpactAssessment{
						RiskLevel: types.RiskLevelCritical,
					},
					RefactoringComplexity: &types.RefactoringComplexity{
						Score: 2,
					},
					FixStrategies: []types.FixStrategy{
						{Type: types.FixStrategyDependencyInject, Suitability: 8},
					},
				},
			},
		}

		got := enricher.Enrich(result)

		// Verify sorted by priority (critical + low complexity first)
		if got.CircularDependencies[0].Cycle[0] != "high" {
			t.Errorf("Expected high-priority cycle first, got %s", got.CircularDependencies[0].Cycle[0])
		}

		// Verify QuickFix was created
		if got.CircularDependencies[0].QuickFix == nil {
			t.Error("Expected QuickFix to be created")
		}
		if got.CircularDependencies[0].QuickFix.StrategyType != types.FixStrategyDependencyInject {
			t.Errorf("Expected best strategy type, got %s", got.CircularDependencies[0].QuickFix.StrategyType)
		}

		// Verify FixSummary was generated
		if got.FixSummary == nil {
			t.Fatal("Expected FixSummary to be created")
		}
		if got.FixSummary.TotalCircularDependencies != 2 {
			t.Errorf("Expected 2 total cycles, got %d", got.FixSummary.TotalCircularDependencies)
		}
		if got.FixSummary.CriticalCyclesCount != 1 {
			t.Errorf("Expected 1 critical cycle, got %d", got.FixSummary.CriticalCyclesCount)
		}
		if got.FixSummary.QuickWinsCount != 1 {
			t.Errorf("Expected 1 quick win (complexity <= 3), got %d", got.FixSummary.QuickWinsCount)
		}
	})

	t.Run("sorts strategies by suitability within each cycle", func(t *testing.T) {
		enricher := NewResultEnricher(nil, nil)

		result := &types.AnalysisResult{
			CircularDependencies: []*types.CircularDependencyInfo{
				{
					Cycle: []string{"a", "b", "a"},
					FixStrategies: []types.FixStrategy{
						{Name: "Low", Suitability: 3},
						{Name: "High", Suitability: 9},
						{Name: "Medium", Suitability: 6},
					},
				},
			},
		}

		got := enricher.Enrich(result)

		strategies := got.CircularDependencies[0].FixStrategies
		if strategies[0].Suitability != 9 {
			t.Errorf("Expected highest suitability first, got %d", strategies[0].Suitability)
		}
		if strategies[1].Suitability != 6 {
			t.Errorf("Expected medium suitability second, got %d", strategies[1].Suitability)
		}
		if strategies[2].Suitability != 3 {
			t.Errorf("Expected lowest suitability last, got %d", strategies[2].Suitability)
		}
	})
}

func TestResultEnricher_GenerateFixSummary(t *testing.T) {
	enricher := NewResultEnricher(nil, nil)

	t.Run("counts quick wins correctly", func(t *testing.T) {
		cycles := []*types.CircularDependencyInfo{
			{
				Cycle: []string{"a", "b", "a"},
				RefactoringComplexity: &types.RefactoringComplexity{
					Score: 2, // Quick win
				},
			},
			{
				Cycle: []string{"c", "d", "c"},
				RefactoringComplexity: &types.RefactoringComplexity{
					Score: 3, // Quick win
				},
			},
			{
				Cycle: []string{"e", "f", "e"},
				RefactoringComplexity: &types.RefactoringComplexity{
					Score: 4, // Not a quick win
				},
			},
		}

		summary := enricher.generateFixSummary(cycles)
		if summary.QuickWinsCount != 2 {
			t.Errorf("Expected 2 quick wins, got %d", summary.QuickWinsCount)
		}
	})

	t.Run("counts critical cycles correctly", func(t *testing.T) {
		cycles := []*types.CircularDependencyInfo{
			{
				Cycle: []string{"a", "b", "a"},
				ImpactAssessment: &types.ImpactAssessment{
					RiskLevel: types.RiskLevelCritical,
				},
			},
			{
				Cycle: []string{"c", "d", "c"},
				ImpactAssessment: &types.ImpactAssessment{
					RiskLevel: types.RiskLevelHigh,
				},
			},
		}

		summary := enricher.generateFixSummary(cycles)
		if summary.CriticalCyclesCount != 1 {
			t.Errorf("Expected 1 critical cycle, got %d", summary.CriticalCyclesCount)
		}
	})

	t.Run("limits high priority cycles to 3", func(t *testing.T) {
		cycles := make([]*types.CircularDependencyInfo, 5)
		for i := 0; i < 5; i++ {
			cycles[i] = &types.CircularDependencyInfo{
				Cycle: []string{"a", "b", "a"},
			}
		}

		summary := enricher.generateFixSummary(cycles)
		if len(summary.HighPriorityCycles) != 3 {
			t.Errorf("Expected 3 high priority cycles, got %d", len(summary.HighPriorityCycles))
		}
	})

	t.Run("high priority cycles have correct data", func(t *testing.T) {
		cycles := []*types.CircularDependencyInfo{
			{
				Cycle:         []string{"@mono/core", "@mono/ui", "@mono/core"},
				PriorityScore: 80.0,
				QuickFix: &types.QuickFixSummary{
					StrategyType:  types.FixStrategyExtractModule,
					EstimatedTime: "30-60 minutes",
				},
			},
		}

		summary := enricher.generateFixSummary(cycles)
		if len(summary.HighPriorityCycles) != 1 {
			t.Fatalf("Expected 1 high priority cycle, got %d", len(summary.HighPriorityCycles))
		}

		hpc := summary.HighPriorityCycles[0]
		if hpc.CycleID != "core→ui" {
			t.Errorf("Expected cycleId 'core→ui', got %q", hpc.CycleID)
		}
		if hpc.PriorityScore != 80.0 {
			t.Errorf("Expected priority 80.0, got %f", hpc.PriorityScore)
		}
		if hpc.RecommendedFix != types.FixStrategyExtractModule {
			t.Errorf("Expected extract-module, got %s", hpc.RecommendedFix)
		}
	})
}
