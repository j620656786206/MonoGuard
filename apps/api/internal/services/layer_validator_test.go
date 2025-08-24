package services

import (
	"context"
	"testing"

	"github.com/monoguard/api/internal/models"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLayerValidatorService_ValidateArchitecture(t *testing.T) {
	// Setup
	mockProjectRepo := new(MockProjectRepository)
	mockAnalysisRepo := new(MockAnalysisRepository)
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	
	service := NewLayerValidatorService(mockProjectRepo, mockAnalysisRepo, logger)
	
	ctx := context.Background()
	projectID := "test-project-id"
	configPath := ".monoguard.yml"
	targetPackages := []string{"apps/frontend", "libs/ui"}
	severityThreshold := "warning"
	
	// Mock project
	project := &models.Project{
		ID:   projectID,
		Name: "Test Project",
		Settings: &models.ProjectSettings{
			ExcludePatterns: []string{"node_modules/**"},
			IncludePatterns: []string{"**/*.ts", "**/*.tsx"},
		},
	}
	
	// Setup mocks
	mockProjectRepo.On("GetByID", ctx, projectID).Return(project, nil)
	mockAnalysisRepo.On("CreateArchitectureValidation", ctx, mock.AnythingOfType("*models.ArchitectureValidation")).Return(nil)
	mockAnalysisRepo.On("UpdateArchitectureValidation", ctx, mock.AnythingOfType("*models.ArchitectureValidation")).Return(nil)
	
	// Execute
	result, err := service.ValidateArchitecture(ctx, projectID, configPath, targetPackages, severityThreshold)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, projectID, result.ProjectID)
	assert.Equal(t, models.StatusCompleted, result.Status)
	assert.NotNil(t, result.Results)
	
	// Verify mocks were called
	mockProjectRepo.AssertExpectations(t)
	mockAnalysisRepo.AssertExpectations(t)
}

func TestLayerValidatorService_matchesPattern(t *testing.T) {
	// Setup
	logger := logrus.New()
	service := &LayerValidatorService{logger: logger}
	
	// Test cases
	tests := []struct {
		name        string
		packagePath string
		pattern     string
		expected    bool
	}{
		{"Simple wildcard match", "apps/frontend", "apps/*", true},
		{"Simple wildcard no match", "libs/ui", "apps/*", false},
		{"Double wildcard match", "libs/ui/components", "libs/**", true},
		{"Exact match", "libs/shared-types", "libs/shared-types", true},
		{"Exact no match", "libs/shared-types", "libs/shared-utils", false},
		{"Complex pattern", "apps/frontend/src", "apps/*/src", true},
		{"Complex pattern no match", "apps/frontend/components", "apps/*/src", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.matchesPattern(tt.packagePath, tt.pattern)
			assert.Equal(t, tt.expected, result, "Pattern: %s, Path: %s", tt.pattern, tt.packagePath)
		})
	}
}

func TestLayerValidatorService_findLayerForPackage(t *testing.T) {
	// Setup
	logger := logrus.New()
	service := &LayerValidatorService{logger: logger}
	
	config := &ArchitectureConfig{
		Architecture: struct {
			Layers []LayerDefinition  `yaml:"layers"`
			Rules  []ArchitectureRule `yaml:"rules"`
		}{
			Layers: []LayerDefinition{
				{
					Name:    "Application Layer",
					Pattern: "apps/*",
				},
				{
					Name:    "UI Library",
					Pattern: "libs/ui/*",
				},
				{
					Name:    "Shared Libraries",
					Pattern: "libs/shared/*",
				},
			},
		},
	}
	
	// Test cases
	tests := []struct {
		name         string
		packagePath  string
		expectedName string
		shouldFind   bool
	}{
		{"App package", "apps/frontend", "Application Layer", true},
		{"UI lib package", "libs/ui/button", "UI Library", true},
		{"Shared lib package", "libs/shared/utils", "Shared Libraries", true},
		{"Unknown package", "external/package", "", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			layer := service.findLayerForPackage(tt.packagePath, config)
			if tt.shouldFind {
				assert.NotNil(t, layer)
				assert.Equal(t, tt.expectedName, layer.Name)
			} else {
				assert.Nil(t, layer)
			}
		})
	}
}

func TestLayerValidatorService_checkImportViolation(t *testing.T) {
	// Setup
	logger := logrus.New()
	service := &LayerValidatorService{logger: logger}
	
	config := &ArchitectureConfig{
		Architecture: struct {
			Layers []LayerDefinition  `yaml:"layers"`
			Rules  []ArchitectureRule `yaml:"rules"`
		}{
			Layers: []LayerDefinition{
				{
					Name:         "Application Layer",
					Pattern:      "apps/*",
					CanImport:    []string{"libs/*"},
					CannotImport: []string{"apps/*"},
				},
			},
		},
	}
	
	layer := &LayerDefinition{
		Name:         "Application Layer",
		Pattern:      "apps/*",
		CanImport:    []string{"libs/*"},
		CannotImport: []string{"apps/*"},
	}
	
	// Test violation case
	violatingImport := ImportAnalysis{
		SourcePackage:   "apps/frontend",
		ImportedPackage: "apps/backend",
		ImportPath:      "../backend/shared",
		LineNumber:      10,
	}
	
	violation := service.checkImportViolation("apps/frontend", violatingImport, layer, config)
	assert.NotNil(t, violation)
	assert.Equal(t, "apps/frontend", violation.PackageName)
	assert.Equal(t, "error", violation.Severity)
	
	// Test allowed case
	allowedImport := ImportAnalysis{
		SourcePackage:   "apps/frontend",
		ImportedPackage: "libs/ui",
		ImportPath:      "@monoguard/ui",
		LineNumber:      5,
	}
	
	noViolation := service.checkImportViolation("apps/frontend", allowedImport, layer, config)
	assert.Nil(t, noViolation)
}

func TestLayerValidatorService_filterBySeverity(t *testing.T) {
	// Setup
	logger := logrus.New()
	service := &LayerValidatorService{logger: logger}
	
	violations := []ViolationReport{
		{Severity: "info"},
		{Severity: "warning"},
		{Severity: "error"},
		{Severity: "warning"},
	}
	
	// Test cases
	tests := []struct {
		name      string
		threshold string
		expected  int
	}{
		{"Filter by error", "error", 1},
		{"Filter by warning", "warning", 3},
		{"Filter by info", "info", 4},
		{"No filter", "", 4},
		{"Invalid threshold", "invalid", 4},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := service.filterBySeverity(violations, tt.threshold)
			assert.Len(t, filtered, tt.expected)
		})
	}
}

func TestLayerValidatorService_generateValidationSummary(t *testing.T) {
	// Setup
	logger := logrus.New()
	service := &LayerValidatorService{logger: logger}
	
	violations := []ViolationReport{
		{Severity: "error"},
		{Severity: "error"},
		{Severity: "warning"},
		{Severity: "info"},
	}
	
	summary := service.generateValidationSummary(violations, 10)
	
	assert.Equal(t, 10, summary.TotalPackages)
	assert.Equal(t, 4, summary.TotalViolations)
	assert.Equal(t, 2, summary.ErrorCount)
	assert.Equal(t, 1, summary.WarningCount)
	assert.Equal(t, 1, summary.InfoCount)
	assert.Equal(t, 80, summary.HealthScore) // 100 - (4/10)*50 = 80
}

func TestLayerValidatorService_getDefaultConfiguration(t *testing.T) {
	// Setup
	logger := logrus.New()
	service := &LayerValidatorService{logger: logger}
	
	config := service.getDefaultConfiguration()
	
	assert.NotNil(t, config)
	assert.Len(t, config.Architecture.Layers, 4) // Should have 4 default layers
	assert.Len(t, config.Architecture.Rules, 3)  // Should have 3 default rules
	
	// Check specific layers exist
	layerNames := make([]string, len(config.Architecture.Layers))
	for i, layer := range config.Architecture.Layers {
		layerNames[i] = layer.Name
	}
	
	assert.Contains(t, layerNames, "Application Layer")
	assert.Contains(t, layerNames, "UI Component Library")
	assert.Contains(t, layerNames, "Shared Libraries")
	assert.Contains(t, layerNames, "Business Logic")
	
	// Check rule severities
	for _, rule := range config.Architecture.Rules {
		assert.NotEmpty(t, rule.Name)
		assert.Contains(t, []string{"error", "warning", "info"}, rule.Severity)
	}
}

func TestLayerValidatorService_analyzePackageImports(t *testing.T) {
	// Setup
	logger := logrus.New()
	service := &LayerValidatorService{logger: logger}
	
	// Test frontend package imports
	imports := service.analyzePackageImports("apps/frontend")
	assert.Len(t, imports, 3)
	
	// Check that one import is a violation (apps/api)
	hasViolation := false
	for _, imp := range imports {
		if imp.ImportedPackage == "apps/api" {
			hasViolation = true
			break
		}
	}
	assert.True(t, hasViolation, "Should detect cross-app dependency violation")
	
	// Test UI library imports
	uiImports := service.analyzePackageImports("libs/ui")
	assert.Len(t, uiImports, 1)
	assert.Equal(t, "libs/shared-types", uiImports[0].ImportedPackage)
}

func TestLayerValidatorService_generateRecommendation(t *testing.T) {
	// Setup
	logger := logrus.New()
	service := &LayerValidatorService{logger: logger}
	
	layer := &LayerDefinition{
		Name: "Application Layer",
	}
	
	// Test app to app dependency
	rec1 := service.generateRecommendation("apps/frontend", "apps/backend", layer)
	assert.Contains(t, rec1, "Move shared functionality to a library")
	
	// Test UI to business dependency
	rec2 := service.generateRecommendation("libs/ui/button", "libs/business/orders", layer)
	assert.Contains(t, rec2, "UI components should not directly depend on business logic")
	
	// Test general case
	rec3 := service.generateRecommendation("libs/a", "libs/b", layer)
	assert.Contains(t, rec3, "layer architecture principles")
}