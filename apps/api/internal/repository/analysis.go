package repository

import (
	"context"
	"fmt"

	"github.com/monoguard/api/internal/models"
	"github.com/monoguard/api/pkg/database"
	"gorm.io/gorm"
)

// AnalysisRepository handles analysis data operations
type AnalysisRepository struct {
	db *database.DB
}

// NewAnalysisRepository creates a new analysis repository
func NewAnalysisRepository(db *database.DB) *AnalysisRepository {
	return &AnalysisRepository{db: db}
}

// CreateDependencyAnalysis creates a new dependency analysis
func (r *AnalysisRepository) CreateDependencyAnalysis(ctx context.Context, analysis *models.DependencyAnalysis) error {
	if err := r.db.WithContext(ctx).Create(analysis).Error; err != nil {
		return fmt.Errorf("failed to create dependency analysis: %w", err)
	}
	return nil
}

// GetDependencyAnalysisByID gets a dependency analysis by ID
func (r *AnalysisRepository) GetDependencyAnalysisByID(ctx context.Context, id string) (*models.DependencyAnalysis, error) {
	var analysis models.DependencyAnalysis
	err := r.db.WithContext(ctx).
		Preload("Project").
		First(&analysis, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get dependency analysis: %w", err)
	}

	return &analysis, nil
}

// GetDependencyAnalysesByProjectID gets dependency analyses for a project
func (r *AnalysisRepository) GetDependencyAnalysesByProjectID(ctx context.Context, projectID string, params *QueryParams) ([]*models.DependencyAnalysis, int64, error) {
	var analyses []*models.DependencyAnalysis
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.DependencyAnalysis{}).
		Where("project_id = ?", projectID)

	// Apply status filter
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count dependency analyses: %w", err)
	}

	// Apply pagination and sorting
	offset := (params.Page - 1) * params.Limit
	orderBy := params.Sort + " " + params.Order

	err := query.
		Order(orderBy).
		Offset(offset).
		Limit(params.Limit).
		Find(&analyses).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get dependency analyses: %w", err)
	}

	return analyses, total, nil
}

// UpdateDependencyAnalysis updates a dependency analysis
func (r *AnalysisRepository) UpdateDependencyAnalysis(ctx context.Context, id string, updates map[string]interface{}) error {
	result := r.db.WithContext(ctx).
		Model(&models.DependencyAnalysis{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update dependency analysis: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// SaveDependencyAnalysis saves a dependency analysis using GORM Save method
func (r *AnalysisRepository) SaveDependencyAnalysis(ctx context.Context, analysis *models.DependencyAnalysis) error {
	if err := r.db.WithContext(ctx).Save(analysis).Error; err != nil {
		return fmt.Errorf("failed to save dependency analysis: %w", err)
	}
	return nil
}

// GetLatestDependencyAnalysis gets the latest dependency analysis for a project
func (r *AnalysisRepository) GetLatestDependencyAnalysis(ctx context.Context, projectID string) (*models.DependencyAnalysis, error) {
	var analysis models.DependencyAnalysis
	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Order("started_at DESC").
		First(&analysis).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get latest dependency analysis: %w", err)
	}

	return &analysis, nil
}

// CreateArchitectureValidation creates a new architecture validation
func (r *AnalysisRepository) CreateArchitectureValidation(ctx context.Context, validation *models.ArchitectureValidation) error {
	if err := r.db.WithContext(ctx).Create(validation).Error; err != nil {
		return fmt.Errorf("failed to create architecture validation: %w", err)
	}
	return nil
}

// UpdateArchitectureValidation updates an architecture validation
func (r *AnalysisRepository) UpdateArchitectureValidation(ctx context.Context, validation *models.ArchitectureValidation) error {
	if err := r.db.WithContext(ctx).Save(validation).Error; err != nil {
		return fmt.Errorf("failed to update architecture validation: %w", err)
	}
	return nil
}

// GetArchitectureValidationByID gets an architecture validation by ID
func (r *AnalysisRepository) GetArchitectureValidationByID(ctx context.Context, id string) (*models.ArchitectureValidation, error) {
	var validation models.ArchitectureValidation
	err := r.db.WithContext(ctx).
		Preload("Project").
		First(&validation, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get architecture validation: %w", err)
	}

	return &validation, nil
}

// CreateHealthScore creates a new health score
func (r *AnalysisRepository) CreateHealthScore(ctx context.Context, healthScore *models.HealthScore) error {
	if err := r.db.WithContext(ctx).Create(healthScore).Error; err != nil {
		return fmt.Errorf("failed to create health score: %w", err)
	}
	return nil
}

// GetLatestHealthScore gets the latest health score for a project
func (r *AnalysisRepository) GetLatestHealthScore(ctx context.Context, projectID string) (*models.HealthScore, error) {
	var healthScore models.HealthScore
	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Order("last_updated DESC").
		First(&healthScore).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get latest health score: %w", err)
	}

	return &healthScore, nil
}

// UpdateHealthScore updates a health score
func (r *AnalysisRepository) UpdateHealthScore(ctx context.Context, id string, updates map[string]interface{}) error {
	result := r.db.WithContext(ctx).
		Model(&models.HealthScore{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update health score: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}