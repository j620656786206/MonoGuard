package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/monoguard/api/internal/repository"
	"github.com/sirupsen/logrus"
)

// AnalysisHandler handles analysis-related HTTP requests
type AnalysisHandler struct {
	analysisRepo *repository.AnalysisRepository
	logger       *logrus.Logger
}

// NewAnalysisHandler creates a new analysis handler
func NewAnalysisHandler(analysisRepo *repository.AnalysisRepository, logger *logrus.Logger) *AnalysisHandler {
	return &AnalysisHandler{
		analysisRepo: analysisRepo,
		logger:       logger,
	}
}

// GetDependencyAnalysis handles GET /analysis/dependencies/:id
func (h *AnalysisHandler) GetDependencyAnalysis(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "Analysis ID is required", nil)
		return
	}

	analysis, err := h.analysisRepo.GetDependencyAnalysisByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			NotFound(c, "Analysis not found")
			return
		}
		h.logger.WithError(err).Error("Failed to get dependency analysis")
		InternalError(c, "Failed to get dependency analysis")
		return
	}

	Success(c, analysis, "Dependency analysis retrieved successfully")
}

// GetProjectDependencyAnalyses handles GET /projects/:projectId/analyses/dependencies
func (h *AnalysisHandler) GetProjectDependencyAnalyses(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		BadRequest(c, "Project ID is required", nil)
		return
	}

	params := h.parseQueryParams(c)

	analyses, total, err := h.analysisRepo.GetDependencyAnalysesByProjectID(c.Request.Context(), projectID, params)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get dependency analyses")
		InternalError(c, "Failed to get dependency analyses")
		return
	}

	pagination := CalculatePagination(params.Page, params.Limit, total)
	SuccessPaginated(c, analyses, pagination, "Dependency analyses retrieved successfully")
}

// GetLatestDependencyAnalysis handles GET /projects/:projectId/analyses/dependencies/latest
func (h *AnalysisHandler) GetLatestDependencyAnalysis(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		BadRequest(c, "Project ID is required", nil)
		return
	}

	analysis, err := h.analysisRepo.GetLatestDependencyAnalysis(c.Request.Context(), projectID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			NotFound(c, "No dependency analysis found for this project")
			return
		}
		h.logger.WithError(err).Error("Failed to get latest dependency analysis")
		InternalError(c, "Failed to get latest dependency analysis")
		return
	}

	Success(c, analysis, "Latest dependency analysis retrieved successfully")
}

// GetArchitectureValidation handles GET /analysis/architecture/:id
func (h *AnalysisHandler) GetArchitectureValidation(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "Analysis ID is required", nil)
		return
	}

	validation, err := h.analysisRepo.GetArchitectureValidationByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			NotFound(c, "Architecture validation not found")
			return
		}
		h.logger.WithError(err).Error("Failed to get architecture validation")
		InternalError(c, "Failed to get architecture validation")
		return
	}

	Success(c, validation, "Architecture validation retrieved successfully")
}

// GetLatestHealthScore handles GET /projects/:projectId/health-score/latest
func (h *AnalysisHandler) GetLatestHealthScore(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		BadRequest(c, "Project ID is required", nil)
		return
	}

	healthScore, err := h.analysisRepo.GetLatestHealthScore(c.Request.Context(), projectID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			NotFound(c, "No health score found for this project")
			return
		}
		h.logger.WithError(err).Error("Failed to get latest health score")
		InternalError(c, "Failed to get latest health score")
		return
	}

	Success(c, healthScore, "Latest health score retrieved successfully")
}

// parseQueryParams parses common query parameters for analysis endpoints
func (h *AnalysisHandler) parseQueryParams(c *gin.Context) *repository.QueryParams {
	// Reuse the same logic from project handler
	params := repository.DefaultQueryParams()

	// Parse page
	if page := c.Query("page"); page != "" {
		if p, err := parseIntDefault(page, 1); err == nil {
			params.Page = p
		}
	}

	// Parse limit
	if limit := c.Query("limit"); limit != "" {
		if l, err := parseIntDefault(limit, 10); err == nil {
			params.Limit = l
		}
	}

	// Parse sort
	if sort := c.Query("sort"); sort != "" {
		params.Sort = sort
	}

	// Parse order
	if order := c.Query("order"); order != "" {
		params.Order = order
	}

	// Parse status filter
	if status := c.Query("status"); status != "" {
		params.Status = status
	}

	return params
}