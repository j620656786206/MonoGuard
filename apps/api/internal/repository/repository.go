package repository

import (
	"context"
	"errors"

	"github.com/monoguard/api/internal/models"
)

// Common errors
var (
	ErrNotFound      = errors.New("record not found")
	ErrDuplicate     = errors.New("record already exists")
	ErrInvalidInput  = errors.New("invalid input")
)

// QueryParams represents common query parameters
type QueryParams struct {
	Page   int    `json:"page" form:"page"`
	Limit  int    `json:"limit" form:"limit"`
	Sort   string `json:"sort" form:"sort"`
	Order  string `json:"order" form:"order"`
	Search string `json:"search" form:"search"`
	Status string `json:"status" form:"status"`
}

// ProjectRepositoryInterface defines the interface for project data operations
type ProjectRepositoryInterface interface {
	GetByID(ctx context.Context, id string) (*models.Project, error)
	GetAll(ctx context.Context, params *QueryParams) ([]*models.Project, int64, error)
	Create(ctx context.Context, project *models.Project) error
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	GetByOwnerID(ctx context.Context, ownerID string, params *QueryParams) ([]*models.Project, int64, error)
	UpdateHealthScore(ctx context.Context, projectID string, score int) error
}

// AnalysisRepositoryInterface defines the interface for analysis data operations
type AnalysisRepositoryInterface interface {
	CreateDependencyAnalysis(ctx context.Context, analysis *models.DependencyAnalysis) error
	GetDependencyAnalysisByID(ctx context.Context, id string) (*models.DependencyAnalysis, error)
	GetDependencyAnalysesByProjectID(ctx context.Context, projectID string, params *QueryParams) ([]*models.DependencyAnalysis, int64, error)
	UpdateDependencyAnalysis(ctx context.Context, id string, updates map[string]interface{}) error
	GetLatestDependencyAnalysis(ctx context.Context, projectID string) (*models.DependencyAnalysis, error)
	CreateArchitectureValidation(ctx context.Context, validation *models.ArchitectureValidation) error
	UpdateArchitectureValidation(ctx context.Context, validation *models.ArchitectureValidation) error
	GetArchitectureValidationByID(ctx context.Context, id string) (*models.ArchitectureValidation, error)
	CreateHealthScore(ctx context.Context, healthScore *models.HealthScore) error
	GetLatestHealthScore(ctx context.Context, projectID string) (*models.HealthScore, error)
	UpdateHealthScore(ctx context.Context, id string, updates map[string]interface{}) error
}

// DefaultQueryParams returns default query parameters
func DefaultQueryParams() *QueryParams {
	return &QueryParams{
		Page:  1,
		Limit: 10,
		Sort:  "created_at",
		Order: "desc",
	}
}

// Validate validates query parameters
func (q *QueryParams) Validate() {
	if q.Page < 1 {
		q.Page = 1
	}
	
	if q.Limit < 1 {
		q.Limit = 10
	}
	
	if q.Limit > 100 {
		q.Limit = 100
	}

	if q.Sort == "" {
		q.Sort = "created_at"
	}

	if q.Order != "asc" && q.Order != "desc" {
		q.Order = "desc"
	}
}