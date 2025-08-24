package services

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/monoguard/api/internal/models"
	"github.com/monoguard/api/internal/repository"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// LayerValidatorService handles architecture layer validation
type LayerValidatorService struct {
	projectRepo  repository.ProjectRepository
	analysisRepo repository.AnalysisRepository
	logger       *logrus.Logger
}

// NewLayerValidatorService creates a new layer validator service
func NewLayerValidatorService(
	projectRepo repository.ProjectRepository,
	analysisRepo repository.AnalysisRepository,
	logger *logrus.Logger,
) *LayerValidatorService {
	return &LayerValidatorService{
		projectRepo:  projectRepo,
		analysisRepo: analysisRepo,
		logger:       logger,
	}
}

// ArchitectureConfig represents the .monoguard.yml configuration
type ArchitectureConfig struct {
	Architecture struct {
		Layers []LayerDefinition `yaml:"layers"`
		Rules  []ArchitectureRule `yaml:"rules"`
	} `yaml:"architecture"`
}

// LayerDefinition represents a layer in the architecture
type LayerDefinition struct {
	Name         string   `yaml:"name"`
	Pattern      string   `yaml:"pattern"`
	Description  string   `yaml:"description"`
	CanImport    []string `yaml:"can_import"`
	CannotImport []string `yaml:"cannot_import"`
}

// ArchitectureRule represents an architecture rule
type ArchitectureRule struct {
	Name        string `yaml:"name"`
	Severity    string `yaml:"severity"`
	Description string `yaml:"description"`
	AutoFix     bool   `yaml:"auto_fix"`
}

// ImportAnalysis represents an import statement analysis
type ImportAnalysis struct {
	SourcePackage   string `json:"source_package"`
	ImportedPackage string `json:"imported_package"`
	ImportType      string `json:"import_type"`
	ImportPath      string `json:"import_path"`
	LineNumber      int    `json:"line_number"`
	IsViolation     bool   `json:"is_violation"`
	ViolationReason string `json:"violation_reason,omitempty"`
}

// ViolationReport represents an architecture violation report
type ViolationReport struct {
	PackageName    string            `json:"package_name"`
	LayerName      string            `json:"layer_name"`
	ViolatedRule   string            `json:"violated_rule"`
	Violations     []ImportAnalysis  `json:"violations"`
	Severity       string            `json:"severity"`
	Recommendation string            `json:"recommendation"`
	FixSuggestion  string            `json:"fix_suggestion"`
}

// ValidateArchitecture validates project architecture against configuration
func (s *LayerValidatorService) ValidateArchitecture(ctx context.Context, projectID string, configPath string, targetPackages []string, severityThreshold string) (*models.ArchitectureValidation, error) {
	s.logger.WithFields(logrus.Fields{
		"project_id":         projectID,
		"config_path":        configPath,
		"target_packages":    targetPackages,
		"severity_threshold": severityThreshold,
	}).Info("Starting architecture validation")

	// Get project
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	// Create validation record
	validation := &models.ArchitectureValidation{
		ID:        uuid.New().String(),
		ProjectID: projectID,
		Status:    models.StatusInProgress,
		StartedAt: time.Now().UTC(),
		Results:   &models.ArchitectureValidationResults{},
		Metadata: &models.AnalysisMetadata{
			Version:          "1.0.0",
			FilesProcessed:   0,
			PackagesAnalyzed: 0,
			Configuration: map[string]interface{}{
				"config_path":        configPath,
				"target_packages":    targetPackages,
				"severity_threshold": severityThreshold,
			},
			Environment: models.AnalysisEnvironment{
				Platform: "linux",
			},
		},
	}

	// Save initial validation
	if err := s.analysisRepo.CreateArchitectureValidation(ctx, validation); err != nil {
		return nil, fmt.Errorf("failed to create validation: %w", err)
	}

	// Load configuration
	config, err := s.loadArchitectureConfig(configPath)
	if err != nil {
		s.logger.WithError(err).Error("Failed to load architecture configuration")
		validation.Status = models.StatusFailed
		s.analysisRepo.UpdateArchitectureValidation(ctx, validation)
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Discover packages
	packages, err := s.discoverPackages(project, targetPackages)
	if err != nil {
		return nil, fmt.Errorf("failed to discover packages: %w", err)
	}

	// Validate each package against layers
	var violations []ViolationReport
	filesProcessed := 0

	for _, pkg := range packages {
		pkgViolations, files, err := s.validatePackage(pkg, config)
		if err != nil {
			s.logger.WithError(err).Warnf("Failed to validate package %s", pkg)
			continue
		}
		violations = append(violations, pkgViolations...)
		filesProcessed += files
	}

	// Filter by severity threshold
	filteredViolations := s.filterBySeverity(violations, severityThreshold)

	// Generate summary
	summary := s.generateValidationSummary(filteredViolations, len(packages))

	// Update validation results
	validation.Results.Violations = s.convertToModelViolations(filteredViolations)
	validation.Results.Summary = summary
	
	// Update metadata
	validation.Metadata.PackagesAnalyzed = len(packages)
	validation.Metadata.FilesProcessed = filesProcessed
	validation.Metadata.Duration = int64(time.Since(validation.StartedAt).Milliseconds())
	
	// Mark as completed
	now := time.Now().UTC()
	validation.CompletedAt = &now
	validation.Status = models.StatusCompleted

	// Update validation
	if err := s.analysisRepo.UpdateArchitectureValidation(ctx, validation); err != nil {
		return nil, fmt.Errorf("failed to update validation: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"project_id":      projectID,
		"violations":      len(filteredViolations),
		"packages":        len(packages),
		"files_processed": filesProcessed,
		"duration_ms":     validation.Metadata.Duration,
	}).Info("Architecture validation completed")

	return validation, nil
}

// loadArchitectureConfig loads the .monoguard.yml configuration
func (s *LayerValidatorService) loadArchitectureConfig(configPath string) (*ArchitectureConfig, error) {
	// If no config path provided, look for default
	if configPath == "" || configPath == ".monoguard.yml" {
		// Return a default configuration for demo purposes
		return s.getDefaultConfiguration(), nil
	}

	// Load configuration file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config ArchitectureConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config YAML: %w", err)
	}

	return &config, nil
}

// getDefaultConfiguration returns a default architecture configuration
func (s *LayerValidatorService) getDefaultConfiguration() *ArchitectureConfig {
	return &ArchitectureConfig{
		Architecture: struct {
			Layers []LayerDefinition  `yaml:"layers"`
			Rules  []ArchitectureRule `yaml:"rules"`
		}{
			Layers: []LayerDefinition{
				{
					Name:         "Application Layer",
					Pattern:      "apps/*",
					Description:  "Frontend and backend applications",
					CanImport:    []string{"libs/*"},
					CannotImport: []string{"apps/*"},
				},
				{
					Name:         "UI Component Library",
					Pattern:      "libs/ui/*",
					Description:  "Reusable UI components",
					CanImport:    []string{"libs/shared/*"},
					CannotImport: []string{"libs/business/*", "apps/*"},
				},
				{
					Name:         "Shared Libraries",
					Pattern:      "libs/shared/*",
					Description:  "Shared utilities and types",
					CanImport:    []string{},
					CannotImport: []string{"apps/*", "libs/ui/*", "libs/business/*"},
				},
				{
					Name:         "Business Logic",
					Pattern:      "libs/business/*",
					Description:  "Core business logic",
					CanImport:    []string{"libs/shared/*"},
					CannotImport: []string{"libs/ui/*", "apps/*"},
				},
			},
			Rules: []ArchitectureRule{
				{
					Name:        "Layer Architecture Violation",
					Severity:    "error",
					Description: "Package violates layer architecture rules",
					AutoFix:     false,
				},
				{
					Name:        "Cross-Application Dependencies",
					Severity:    "error", 
					Description: "Applications cannot depend on other applications",
					AutoFix:     false,
				},
				{
					Name:        "Upward Dependencies",
					Severity:    "warning",
					Description: "Lower layers should not depend on higher layers",
					AutoFix:     false,
				},
			},
		},
	}
}

// discoverPackages discovers packages to validate
func (s *LayerValidatorService) discoverPackages(project *models.Project, targetPackages []string) ([]string, error) {
	// If target packages specified, use them
	if len(targetPackages) > 0 {
		return targetPackages, nil
	}

	// Otherwise, discover from project structure
	// For demo purposes, return common monorepo packages
	return []string{
		"apps/frontend",
		"apps/api",
		"apps/cli", 
		"libs/shared-types",
		"libs/ui",
	}, nil
}

// validatePackage validates a package against architecture layers
func (s *LayerValidatorService) validatePackage(packagePath string, config *ArchitectureConfig) ([]ViolationReport, int, error) {
	var violations []ViolationReport
	filesProcessed := 0

	// Find which layer this package belongs to
	layer := s.findLayerForPackage(packagePath, config)
	if layer == nil {
		// Package doesn't match any layer - this could be a violation
		violations = append(violations, ViolationReport{
			PackageName:    packagePath,
			LayerName:      "Unknown",
			ViolatedRule:   "Unclassified Package",
			Severity:       "warning",
			Recommendation: "Add package pattern to architecture configuration",
			FixSuggestion:  fmt.Sprintf("Add pattern matching '%s' to .monoguard.yml", packagePath),
		})
		return violations, 0, nil
	}

	// Analyze imports for this package (mock analysis for demo)
	imports := s.analyzePackageImports(packagePath)
	filesProcessed = len(imports)

	// Check each import against layer rules
	for _, imp := range imports {
		violation := s.checkImportViolation(packagePath, imp, layer, config)
		if violation != nil {
			violations = append(violations, *violation)
		}
	}

	return violations, filesProcessed, nil
}

// findLayerForPackage finds which architecture layer a package belongs to
func (s *LayerValidatorService) findLayerForPackage(packagePath string, config *ArchitectureConfig) *LayerDefinition {
	for _, layer := range config.Architecture.Layers {
		if s.matchesPattern(packagePath, layer.Pattern) {
			return &layer
		}
	}
	return nil
}

// matchesPattern checks if a package path matches a glob pattern
func (s *LayerValidatorService) matchesPattern(packagePath, pattern string) bool {
	// Convert glob pattern to regex
	regexPattern := strings.Replace(pattern, "*", "[^/]*", -1)
	regexPattern = strings.Replace(regexPattern, "/**", "/.*", -1) 
	regexPattern = "^" + regexPattern + "$"
	
	matched, err := regexp.MatchString(regexPattern, packagePath)
	if err != nil {
		s.logger.WithError(err).Warnf("Invalid pattern: %s", pattern)
		return false
	}
	
	return matched
}

// analyzePackageImports analyzes imports for a package (mock implementation)
func (s *LayerValidatorService) analyzePackageImports(packagePath string) []ImportAnalysis {
	// Mock import analysis - in real implementation would parse TypeScript/JavaScript files
	mockImports := []ImportAnalysis{}

	switch packagePath {
	case "apps/frontend":
		mockImports = []ImportAnalysis{
			{
				SourcePackage:   "apps/frontend",
				ImportedPackage: "libs/ui",
				ImportType:      "named",
				ImportPath:      "@monoguard/ui",
				LineNumber:      1,
			},
			{
				SourcePackage:   "apps/frontend", 
				ImportedPackage: "libs/shared-types",
				ImportType:      "named",
				ImportPath:      "@monoguard/shared-types",
				LineNumber:      2,
			},
			{
				SourcePackage:   "apps/frontend",
				ImportedPackage: "apps/api", // This should be a violation
				ImportType:      "named", 
				ImportPath:      "../api/shared-utils",
				LineNumber:      3,
			},
		}
	case "libs/ui":
		mockImports = []ImportAnalysis{
			{
				SourcePackage:   "libs/ui",
				ImportedPackage: "libs/shared-types",
				ImportType:      "named",
				ImportPath:      "@monoguard/shared-types", 
				LineNumber:      1,
			},
		}
	default:
		mockImports = []ImportAnalysis{
			{
				SourcePackage:   packagePath,
				ImportedPackage: "libs/shared-types",
				ImportType:      "named",
				ImportPath:      "@monoguard/shared-types",
				LineNumber:      1,
			},
		}
	}

	return mockImports
}

// checkImportViolation checks if an import violates layer rules
func (s *LayerValidatorService) checkImportViolation(packagePath string, imp ImportAnalysis, layer *LayerDefinition, config *ArchitectureConfig) *ViolationReport {
	importedPackage := imp.ImportedPackage
	
	// Check cannot_import rules
	for _, pattern := range layer.CannotImport {
		if s.matchesPattern(importedPackage, pattern) {
			return &ViolationReport{
				PackageName:  packagePath,
				LayerName:    layer.Name,
				ViolatedRule: "Layer Architecture Violation",
				Violations: []ImportAnalysis{
					{
						SourcePackage:   imp.SourcePackage,
						ImportedPackage: imp.ImportedPackage,
						ImportPath:      imp.ImportPath,
						LineNumber:      imp.LineNumber,
						IsViolation:     true,
						ViolationReason: fmt.Sprintf("Package '%s' in layer '%s' cannot import from '%s'", packagePath, layer.Name, importedPackage),
					},
				},
				Severity:       "error",
				Recommendation: s.generateRecommendation(packagePath, importedPackage, layer),
				FixSuggestion:  s.generateFixSuggestion(packagePath, importedPackage, layer),
			}
		}
	}

	// Check can_import rules (if specified and not empty)
	if len(layer.CanImport) > 0 {
		allowed := false
		for _, pattern := range layer.CanImport {
			if s.matchesPattern(importedPackage, pattern) {
				allowed = true
				break
			}
		}
		
		if !allowed {
			return &ViolationReport{
				PackageName:  packagePath,
				LayerName:    layer.Name,
				ViolatedRule: "Layer Architecture Violation",
				Violations: []ImportAnalysis{
					{
						SourcePackage:   imp.SourcePackage,
						ImportedPackage: imp.ImportedPackage,
						ImportPath:      imp.ImportPath,
						LineNumber:      imp.LineNumber,
						IsViolation:     true,
						ViolationReason: fmt.Sprintf("Package '%s' in layer '%s' can only import from %v", packagePath, layer.Name, layer.CanImport),
					},
				},
				Severity:       "warning",
				Recommendation: s.generateRecommendation(packagePath, importedPackage, layer),
				FixSuggestion:  s.generateFixSuggestion(packagePath, importedPackage, layer),
			}
		}
	}

	return nil
}

// generateRecommendation generates a recommendation for fixing violation
func (s *LayerValidatorService) generateRecommendation(packagePath, importedPackage string, layer *LayerDefinition) string {
	if strings.Contains(importedPackage, "apps/") {
		return "Move shared functionality to a library package"
	}
	if strings.Contains(packagePath, "libs/ui") && strings.Contains(importedPackage, "libs/business") {
		return "UI components should not directly depend on business logic"
	}
	return "Refactor to follow layer architecture principles"
}

// generateFixSuggestion generates a specific fix suggestion
func (s *LayerValidatorService) generateFixSuggestion(packagePath, importedPackage string, layer *LayerDefinition) string {
	if strings.Contains(importedPackage, "apps/") {
		return "Create a shared library in libs/shared and move the common code there"
	}
	return "Consider using dependency injection or extracting interfaces to a shared package"
}

// filterBySeverity filters violations by severity threshold
func (s *LayerValidatorService) filterBySeverity(violations []ViolationReport, threshold string) []ViolationReport {
	if threshold == "" {
		return violations
	}

	severityOrder := map[string]int{
		"info":    1,
		"warning": 2,
		"error":   3,
	}

	thresholdLevel, exists := severityOrder[threshold]
	if !exists {
		return violations
	}

	var filtered []ViolationReport
	for _, violation := range violations {
		if severityOrder[violation.Severity] >= thresholdLevel {
			filtered = append(filtered, violation)
		}
	}

	return filtered
}

// generateValidationSummary generates a summary of validation results
func (s *LayerValidatorService) generateValidationSummary(violations []ViolationReport, totalPackages int) models.ValidationSummary {
	summary := models.ValidationSummary{
		LayersAnalyzed:     totalPackages,
		TotalViolations:    len(violations),
	}

	for _, violation := range violations {
		switch violation.Severity {
		case "error":
			summary.CriticalViolations++
		case "warning":
			summary.WarningViolations++
		}
	}

	// Calculate overall compliance
	if totalPackages > 0 {
		violationRatio := float64(len(violations)) / float64(totalPackages)
		summary.OverallCompliance = (100 - (violationRatio * 50)) // Up to 50 points penalty
		if summary.OverallCompliance < 0 {
			summary.OverallCompliance = 0
		}
	} else {
		summary.OverallCompliance = 100
	}

	return summary
}

// convertToModelViolations converts internal violations to model format
func (s *LayerValidatorService) convertToModelViolations(violations []ViolationReport) []models.ArchitectureViolation {
	var result []models.ArchitectureViolation

	for _, violation := range violations {
		modelViolation := models.ArchitectureViolation{
			RuleName:        violation.ViolatedRule,
			Severity:        models.Severity(violation.Severity),
			Description:     violation.ViolatedRule,
			ViolatingFile:   violation.PackageName,
			ViolatingImport: "",
			ExpectedLayer:   violation.LayerName,
			ActualLayer:     "",
			Suggestion:      violation.FixSuggestion,
		}
		result = append(result, modelViolation)
	}

	return result
}