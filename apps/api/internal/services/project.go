package services

import (
	"context"
	"fmt"
	"time"

	"github.com/monoguard/api/internal/models"
	"github.com/monoguard/api/internal/repository"
	"github.com/monoguard/api/internal/utils"
	"github.com/sirupsen/logrus"
)

// ProjectService handles project business logic
type ProjectService struct {
	projectRepo   *repository.ProjectRepository
	analysisRepo  *repository.AnalysisRepository
	analyzer      *DependencyAnalyzer
	logger        *logrus.Logger
}

// NewProjectService creates a new project service
func NewProjectService(
	projectRepo *repository.ProjectRepository,
	analysisRepo *repository.AnalysisRepository,
	analyzer *DependencyAnalyzer,
	logger *logrus.Logger,
) *ProjectService {
	return &ProjectService{
		projectRepo:  projectRepo,
		analysisRepo: analysisRepo,
		analyzer:     analyzer,
		logger:       logger,
	}
}

// CreateProjectRequest represents a request to create a project
type CreateProjectRequest struct {
	Name          string                  `json:"name" validate:"required,min=1,max=100"`
	Description   *string                 `json:"description,omitempty"`
	RepositoryURL string                  `json:"repositoryUrl" validate:"required,url"`
	Branch        string                  `json:"branch" validate:"required"`
	OwnerID       string                  `json:"ownerId" validate:"required"`
	Settings      *models.ProjectSettings `json:"settings,omitempty"`
}

// UpdateProjectRequest represents a request to update a project
type UpdateProjectRequest struct {
	Name        *string                 `json:"name,omitempty"`
	Description *string                 `json:"description,omitempty"`
	Branch      *string                 `json:"branch,omitempty"`
	Settings    *models.ProjectSettings `json:"settings,omitempty"`
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(ctx context.Context, req *CreateProjectRequest) (*models.Project, error) {
	s.logger.WithFields(logrus.Fields{
		"name":           req.Name,
		"repository_url": req.RepositoryURL,
		"owner_id":       req.OwnerID,
	}).Info("Creating new project")

	// Set default settings if not provided
	settings := req.Settings
	if settings == nil {
		settings = &models.ProjectSettings{
			AutoAnalysis: true,
			NotificationSettings: models.NotificationSettings{
				Email:    false,
				Severity: []string{"high", "critical"}, // Temporarily using string values
			},
			ExcludePatterns: []string{"node_modules/**", "dist/**", "build/**"},
			IncludePatterns: []string{"**/*.json"},
			ArchitectureRules: models.ArchitectureRules{
				Layers: []models.ArchitectureLayer{},
				Rules:  []models.ArchitectureRule{},
			},
		}
	}

	description := ""
	if req.Description != nil {
		description = *req.Description
	}
	
	project := &models.Project{
		Name:          req.Name,
		Description:   description,                    // Now a string, not pointer
		RepositoryURL: req.RepositoryURL,
		Branch:        req.Branch,
		Status:        "pending",
		HealthScore:   0,
		OwnerID:       req.OwnerID,
	}

	// Generate UUID for the project
	utils.GenerateProjectID(project)

	if err := s.projectRepo.Create(ctx, project); err != nil {
		s.logger.WithError(err).Error("Failed to create project")
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	s.logger.WithField("project_id", project.ID).Info("Project created successfully")
	return project, nil
}

// GetProject gets a project by ID
func (s *ProjectService) GetProject(ctx context.Context, id string) (*models.Project, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrProjectNotFound
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return project, nil
}

// GetProjects gets all projects with pagination
func (s *ProjectService) GetProjects(ctx context.Context, params *repository.QueryParams) ([]*models.Project, int64, error) {
	params.Validate()
	
	projects, total, err := s.projectRepo.GetAll(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get projects: %w", err)
	}

	return projects, total, nil
}

// GetProjectsByOwner gets projects by owner ID
func (s *ProjectService) GetProjectsByOwner(ctx context.Context, ownerID string, params *repository.QueryParams) ([]*models.Project, int64, error) {
	params.Validate()
	
	projects, total, err := s.projectRepo.GetByOwnerID(ctx, ownerID, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get projects by owner: %w", err)
	}

	return projects, total, nil
}

// UpdateProject updates a project
func (s *ProjectService) UpdateProject(ctx context.Context, id string, req *UpdateProjectRequest) (*models.Project, error) {
	s.logger.WithField("project_id", id).Info("Updating project")

	updates := make(map[string]interface{})
	
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	
	if req.Branch != nil {
		updates["branch"] = *req.Branch
	}
	
	// Settings field temporarily removed for Railway compatibility
	// if req.Settings != nil {
	//	updates["settings"] = *req.Settings
	// }
	
	if len(updates) == 0 {
		return s.GetProject(ctx, id)
	}

	updates["updated_at"] = time.Now().UTC()

	if err := s.projectRepo.Update(ctx, id, updates); err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrProjectNotFound
		}
		s.logger.WithError(err).WithField("project_id", id).Error("Failed to update project")
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	s.logger.WithField("project_id", id).Info("Project updated successfully")
	return s.GetProject(ctx, id)
}

// DeleteProject deletes a project
func (s *ProjectService) DeleteProject(ctx context.Context, id string) error {
	s.logger.WithField("project_id", id).Info("Deleting project")

	if err := s.projectRepo.Delete(ctx, id); err != nil {
		if err == repository.ErrNotFound {
			return ErrProjectNotFound
		}
		s.logger.WithError(err).WithField("project_id", id).Error("Failed to delete project")
		return fmt.Errorf("failed to delete project: %w", err)
	}

	s.logger.WithField("project_id", id).Info("Project deleted successfully")
	return nil
}

// TriggerAnalysis triggers a dependency analysis for a project
func (s *ProjectService) TriggerAnalysis(ctx context.Context, projectID string, repoPath string) (*models.DependencyAnalysis, error) {
	s.logger.WithFields(logrus.Fields{
		"project_id": projectID,
		"repo_path":  repoPath,
	}).Info("Triggering dependency analysis")

	// Get project to validate it exists - temporarily not using project.Settings
	_, err := s.GetProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// Create new analysis record
	analysis := &models.DependencyAnalysis{
		ProjectID: projectID,
		Status:    "in_progress",
		StartedAt: time.Now().UTC(),
		Metadata: &models.AnalysisMetadata{
			Version:          "1.0.0",
			FilesProcessed:   0,
			PackagesAnalyzed: 0,
			Configuration:    map[string]interface{}{
				"excludePatterns": []string{}, // project.Settings.ExcludePatterns - temporarily disabled
				"includePatterns": []string{}, // project.Settings.IncludePatterns - temporarily disabled
			},
			Environment: models.AnalysisEnvironment{
				Platform: "linux",
			},
		},
	}

	if err := s.analysisRepo.CreateDependencyAnalysis(ctx, analysis); err != nil {
		return nil, fmt.Errorf("failed to create analysis record: %w", err)
	}

	// Update project status
	_ = s.projectRepo.Update(ctx, projectID, map[string]interface{}{
		"status":           "in_progress",
		"last_analysis_at": time.Now().UTC(),
	})

	// Run analysis in background (in a real implementation, this would be async)
	go s.runAnalysis(context.Background(), analysis, repoPath)

	return analysis, nil
}

// runAnalysis runs the actual dependency analysis
func (s *ProjectService) runAnalysis(ctx context.Context, analysis *models.DependencyAnalysis, repoPath string) {
	startTime := time.Now()
	
	defer func() {
		duration := time.Since(startTime)
		s.logger.WithFields(logrus.Fields{
			"analysis_id": analysis.ID,
			"duration":    duration,
		}).Info("Analysis completed")
	}()

	// Perform the analysis
	results, err := s.analyzer.AnalyzeMonorepo(ctx, repoPath, analysis.ProjectID)
	if err != nil {
		s.logger.WithError(err).WithField("analysis_id", analysis.ID).Error("Analysis failed")
		
		// Update analysis status to failed
		_ = s.analysisRepo.UpdateDependencyAnalysis(ctx, analysis.ID, map[string]interface{}{
			"status":       "failed",
			"completed_at": time.Now().UTC(),
		})
		
		// Update project status to failed
		_ = s.projectRepo.Update(ctx, analysis.ProjectID, map[string]interface{}{
			"status": "failed",
		})
		
		return
	}

	duration := time.Since(startTime)
	
	// Update analysis with results
	updates := map[string]interface{}{
		"status":       "completed",
		"completed_at": time.Now().UTC(),
		"results":      *results,
		"metadata": models.AnalysisMetadata{
			Version:          "1.0.0",
			Duration:         duration.Milliseconds(),
			FilesProcessed:   100, // Would be actual count
			PackagesAnalyzed: results.Summary.TotalPackages,
			Configuration:    analysis.Metadata.Configuration,
			Environment:      analysis.Metadata.Environment,
		},
	}

	if err := s.analysisRepo.UpdateDependencyAnalysis(ctx, analysis.ID, updates); err != nil {
		s.logger.WithError(err).WithField("analysis_id", analysis.ID).Error("Failed to update analysis")
		return
	}

	// Update project status and health score
	_ = s.projectRepo.Update(ctx, analysis.ProjectID, map[string]interface{}{
		"status":       "completed",
		"health_score": int(results.Summary.HealthScore),
	})

	s.logger.WithFields(logrus.Fields{
		"analysis_id":  analysis.ID,
		"project_id":   analysis.ProjectID,
		"health_score": results.Summary.HealthScore,
		"duplicates":   results.Summary.DuplicateCount,
		"conflicts":    results.Summary.ConflictCount,
		"unused":       results.Summary.UnusedCount,
	}).Info("Analysis completed successfully")
}