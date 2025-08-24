package repository

import (
	"context"
	"fmt"

	"github.com/monoguard/api/internal/models"
	"github.com/monoguard/api/pkg/database"
	"gorm.io/gorm"
)

// ProjectRepository handles project data operations
type ProjectRepository struct {
	db *database.DB
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(db *database.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// Create creates a new project
func (r *ProjectRepository) Create(ctx context.Context, project *models.Project) error {
	if err := r.db.WithContext(ctx).Create(project).Error; err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}
	return nil
}

// GetByID gets a project by ID
func (r *ProjectRepository) GetByID(ctx context.Context, id string) (*models.Project, error) {
	var project models.Project
	err := r.db.WithContext(ctx).
		Preload("DependencyAnalyses").
		Preload("ArchitectureValidations").
		First(&project, "id = ?", id).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return &project, nil
}

// GetAll gets all projects with pagination
func (r *ProjectRepository) GetAll(ctx context.Context, params *QueryParams) ([]*models.Project, int64, error) {
	var projects []*models.Project
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Project{})

	// Apply filters
	if params.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count projects: %w", err)
	}

	// Apply pagination and sorting
	offset := (params.Page - 1) * params.Limit
	orderBy := params.Sort + " " + params.Order
	
	err := query.
		Preload("DependencyAnalyses", func(db *gorm.DB) *gorm.DB {
			return db.Order("started_at DESC").Limit(5)
		}).
		Preload("ArchitectureValidations", func(db *gorm.DB) *gorm.DB {
			return db.Order("started_at DESC").Limit(5)
		}).
		Order(orderBy).
		Offset(offset).
		Limit(params.Limit).
		Find(&projects).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get projects: %w", err)
	}

	return projects, total, nil
}

// Update updates a project
func (r *ProjectRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	result := r.db.WithContext(ctx).
		Model(&models.Project{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update project: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// Delete deletes a project
func (r *ProjectRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.Project{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete project: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// GetByOwnerID gets projects by owner ID
func (r *ProjectRepository) GetByOwnerID(ctx context.Context, ownerID string, params *QueryParams) ([]*models.Project, int64, error) {
	var projects []*models.Project
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.Project{}).
		Where("owner_id = ?", ownerID)

	// Apply filters
	if params.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count projects: %w", err)
	}

	// Apply pagination and sorting
	offset := (params.Page - 1) * params.Limit
	orderBy := params.Sort + " " + params.Order

	err := query.
		Preload("DependencyAnalyses", func(db *gorm.DB) *gorm.DB {
			return db.Order("started_at DESC").Limit(5)
		}).
		Preload("ArchitectureValidations", func(db *gorm.DB) *gorm.DB {
			return db.Order("started_at DESC").Limit(5)
		}).
		Order(orderBy).
		Offset(offset).
		Limit(params.Limit).
		Find(&projects).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get projects: %w", err)
	}

	return projects, total, nil
}

// UpdateHealthScore updates the health score of a project
func (r *ProjectRepository) UpdateHealthScore(ctx context.Context, projectID string, score int) error {
	err := r.db.WithContext(ctx).
		Model(&models.Project{}).
		Where("id = ?", projectID).
		Update("health_score", score).Error

	if err != nil {
		return fmt.Errorf("failed to update health score: %w", err)
	}

	return nil
}