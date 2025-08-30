package services

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/monoguard/api/internal/models"
	"github.com/monoguard/api/internal/repository"
	"github.com/sirupsen/logrus"
)

// IntegratedAnalysisService provides comprehensive analysis capabilities
type IntegratedAnalysisService struct {
	basicEngine          *BasicAnalysisEngine
	circularDetector     *CircularDetectorService
	layerValidator       *LayerValidatorService
	uploadService        *UploadService
	projectRepo          repository.ProjectRepositoryInterface
	analysisRepo         repository.AnalysisRepositoryInterface
	logger               *logrus.Logger
}

// NewIntegratedAnalysisService creates a new integrated analysis service
func NewIntegratedAnalysisService(
	basicEngine *BasicAnalysisEngine,
	circularDetector *CircularDetectorService,
	layerValidator *LayerValidatorService,
	uploadService *UploadService,
	projectRepo repository.ProjectRepositoryInterface,
	analysisRepo repository.AnalysisRepositoryInterface,
	logger *logrus.Logger,
) *IntegratedAnalysisService {
	return &IntegratedAnalysisService{
		basicEngine:      basicEngine,
		circularDetector: circularDetector,
		layerValidator:   layerValidator,
		uploadService:    uploadService,
		projectRepo:      projectRepo,
		analysisRepo:     analysisRepo,
		logger:           logger,
	}
}

// AnalyzeProcessingResult performs comprehensive analysis on uploaded files
func (s *IntegratedAnalysisService) AnalyzeProcessingResult(ctx context.Context, processingResultID string, analysisOptions AnalysisOptions) (*models.DependencyAnalysis, error) {
	s.logger.WithField("processing_result_id", processingResultID).Info("Starting integrated analysis")

	// Get processing result
	processingResult, err := s.uploadService.GetProcessingResult(processingResultID)
	if err != nil {
		return nil, fmt.Errorf("failed to get processing result: %w", err)
	}

	if len(processingResult.PackageJsonFiles) == 0 {
		return nil, fmt.Errorf("no package.json files found in processing result")
	}

	// Create temporary directory for analysis
	tempDir, err := s.createTempAnalysisDir(processingResult)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp analysis directory: %w", err)
	}
	defer s.cleanupTempDir(tempDir)

	// Create or get project
	project, err := s.getOrCreateProject(processingResult)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create project: %w", err)
	}

	// Run comprehensive analysis
	return s.runComprehensiveAnalysis(ctx, tempDir, project.ID, analysisOptions)
}

// AnalysisOptions defines options for analysis
type AnalysisOptions struct {
	IncludeCircular     bool   `json:"include_circular"`
	IncludeArchitecture bool   `json:"include_architecture"`
	SeverityThreshold   string `json:"severity_threshold"`
	ConfigPath          string `json:"config_path"`
}

// createTempAnalysisDir creates a temporary directory structure for analysis
func (s *IntegratedAnalysisService) createTempAnalysisDir(processingResult *models.FileProcessingResult) (string, error) {
	tempDir, err := ioutil.TempDir("", fmt.Sprintf("monoguard-analysis-%s", processingResult.ID))
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	s.logger.WithField("temp_dir", tempDir).Debug("Created temporary analysis directory")

	// Write package.json files to temp directory
	for i, pkgFile := range processingResult.PackageJsonFiles {
		// Create directory structure if path contains subdirectories
		dir := filepath.Dir(pkgFile.Path)
		if dir != "." {
			fullDir := filepath.Join(tempDir, dir)
			if err := os.MkdirAll(fullDir, 0755); err != nil {
				return "", fmt.Errorf("failed to create directory %s: %w", fullDir, err)
			}
		}

		filePath := filepath.Join(tempDir, pkgFile.Path)
		if err := ioutil.WriteFile(filePath, []byte(pkgFile.Content), 0644); err != nil {
			return "", fmt.Errorf("failed to write package.json file %d: %w", i, err)
		}

		s.logger.WithFields(logrus.Fields{
			"file_path": filePath,
			"file_size": len(pkgFile.Content),
		}).Debug("Wrote package.json to temp directory")
	}

	return tempDir, nil
}

// cleanupTempDir removes the temporary directory
func (s *IntegratedAnalysisService) cleanupTempDir(tempDir string) {
	if err := os.RemoveAll(tempDir); err != nil {
		s.logger.WithError(err).WithField("temp_dir", tempDir).Warn("Failed to cleanup temp directory")
	} else {
		s.logger.WithField("temp_dir", tempDir).Debug("Cleaned up temporary directory")
	}
}

// getOrCreateProject gets existing project or creates a new one
func (s *IntegratedAnalysisService) getOrCreateProject(processingResult *models.FileProcessingResult) (*models.Project, error) {
	// For now, create a temporary project for analysis
	// In a real implementation, this would be associated with a user account
	description := "Temporary project for file upload analysis"
	project := &models.Project{
		ID:          uuid.New().String(),
		Name:        fmt.Sprintf("Analysis-%s", processingResult.ID[:8]),
		Description: &description,
		RepositoryURL: "",
		// Settings: &models.ProjectSettings{
		//	ExcludePatterns: []string{"node_modules/**", "dist/**", "build/**"},
		//	IncludePatterns: []string{"**/*.json"},
		// }, // Temporarily commented out
		Status:    models.StatusInProgress,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.projectRepo.Create(context.Background(), project); err != nil {
		return nil, fmt.Errorf("failed to create temporary project: %w", err)
	}

	s.logger.WithField("project_id", project.ID).Info("Created temporary project for analysis")

	return project, nil
}

// runComprehensiveAnalysis runs all available analyses
func (s *IntegratedAnalysisService) runComprehensiveAnalysis(ctx context.Context, repoPath string, projectID string, options AnalysisOptions) (*models.DependencyAnalysis, error) {
	s.logger.WithFields(logrus.Fields{
		"project_id": projectID,
		"repo_path":  repoPath,
		"options":    options,
	}).Info("Running comprehensive analysis")

	// Create analysis record
	analysis := &models.DependencyAnalysis{
		ID:        uuid.New().String(),
		ProjectID: projectID,
		Status:    models.StatusInProgress,
		StartedAt: time.Now().UTC(),
		Results:   &models.DependencyAnalysisResults{},
		Metadata: &models.AnalysisMetadata{
			Version:          "1.0.0",
			FilesProcessed:   0,
			PackagesAnalyzed: 0,
			Configuration: map[string]interface{}{
				"include_circular":     options.IncludeCircular,
				"include_architecture": options.IncludeArchitecture,
				"severity_threshold":   options.SeverityThreshold,
			},
			Environment: models.AnalysisEnvironment{
				Platform: "linux",
			},
		},
	}

	// Save initial analysis
	if err := s.analysisRepo.CreateDependencyAnalysis(ctx, analysis); err != nil {
		return nil, fmt.Errorf("failed to create analysis: %w", err)
	}

	// Run basic analysis (always included)
	s.logger.Debug("Running basic dependency analysis")
	basicResults, err := s.basicEngine.AnalyzeRepository(ctx, repoPath, projectID)
	if err != nil {
		s.logger.WithError(err).Error("Basic analysis failed")
		analysis.Status = models.StatusFailed
		s.analysisRepo.UpdateDependencyAnalysis(ctx, analysis.ID, map[string]interface{}{
			"status": models.StatusFailed,
		})
		return nil, fmt.Errorf("basic analysis failed: %w", err)
	}

	// Merge basic results
	analysis.Results = basicResults

	// Run circular dependency detection if requested
	if options.IncludeCircular {
		s.logger.Debug("Running circular dependency analysis")
		circularResults, err := s.circularDetector.DetectCircularDependencies(ctx, projectID, repoPath)
		if err != nil {
			s.logger.WithError(err).Warn("Circular dependency analysis failed, continuing without it")
		} else {
			// Merge circular dependency results
			analysis.Results.CircularDependencies = circularResults.Results.CircularDependencies
			// Update summary
			analysis.Results.Summary.CircularCount = len(circularResults.Results.CircularDependencies)
			// Adjust health score
			if len(circularResults.Results.CircularDependencies) > 0 {
				analysis.Results.Summary.HealthScore -= float64(len(circularResults.Results.CircularDependencies)) * 15
				if analysis.Results.Summary.HealthScore < 0 {
					analysis.Results.Summary.HealthScore = 0
				}
			}
		}
	}

	// Run architecture validation if requested
	if options.IncludeArchitecture {
		s.logger.Debug("Running architecture validation")
		architectureResults, err := s.layerValidator.ValidateArchitecture(
			ctx, 
			projectID, 
			options.ConfigPath, 
			[]string{}, 
			options.SeverityThreshold,
		)
		if err != nil {
			s.logger.WithError(err).Warn("Architecture validation failed, continuing without it")
		} else {
			// Store architecture results separately (they have their own model)
			s.logger.WithField("violations", len(architectureResults.Results.Violations)).Info("Architecture validation completed")
		}
	}

	// Calculate final metadata
	analysis.Metadata.Duration = int64(time.Since(analysis.StartedAt).Milliseconds())

	// Mark as completed
	now := time.Now().UTC()
	analysis.CompletedAt = &now
	analysis.Status = models.StatusCompleted

	// Update analysis using Save instead of Updates to handle complex types
	if err := s.analysisRepo.SaveDependencyAnalysis(ctx, analysis); err != nil {
		return nil, fmt.Errorf("failed to update analysis: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"project_id":        projectID,
		"analysis_id":       analysis.ID,
		"duplicates":        analysis.Results.Summary.DuplicateCount,
		"conflicts":         analysis.Results.Summary.ConflictCount,
		"circular":          analysis.Results.Summary.CircularCount,
		"health_score":      analysis.Results.Summary.HealthScore,
		"duration_ms":       analysis.Metadata.Duration,
	}).Info("Comprehensive analysis completed successfully")

	return analysis, nil
}

// GetAnalysisByID retrieves an analysis by ID
func (s *IntegratedAnalysisService) GetAnalysisByID(ctx context.Context, analysisID string) (*models.DependencyAnalysis, error) {
	analysis, err := s.analysisRepo.GetDependencyAnalysisByID(ctx, analysisID)
	if err != nil {
		return nil, fmt.Errorf("failed to get analysis: %w", err)
	}

	return analysis, nil
}

// ListAnalysesForProject lists all analyses for a project
func (s *IntegratedAnalysisService) ListAnalysesForProject(ctx context.Context, projectID string, limit, offset int) ([]*models.DependencyAnalysis, int64, error) {
	params := &repository.QueryParams{
		Page:  offset/limit + 1,
		Limit: limit,
	}
	analyses, total, err := s.analysisRepo.GetDependencyAnalysesByProjectID(ctx, projectID, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list analyses: %w", err)
	}

	return analyses, total, nil
}