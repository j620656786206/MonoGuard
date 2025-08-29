package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/monoguard/api/internal/repository"
	"github.com/monoguard/api/internal/services"
	"github.com/sirupsen/logrus"
)

// AnalysisHandler handles analysis-related HTTP requests
type AnalysisHandler struct {
	analysisRepo          *repository.AnalysisRepository
	integratedAnalysis    *services.IntegratedAnalysisService
	logger                *logrus.Logger
}

// NewAnalysisHandler creates a new analysis handler
func NewAnalysisHandler(analysisRepo *repository.AnalysisRepository, integratedAnalysis *services.IntegratedAnalysisService, logger *logrus.Logger) *AnalysisHandler {
	return &AnalysisHandler{
		analysisRepo:       analysisRepo,
		integratedAnalysis: integratedAnalysis,
		logger:             logger,
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

// AnalyzeUploadRequest represents the request body for analyzing uploaded files
type AnalyzeUploadRequest struct {
	ProcessingResultID  string `json:"processing_result_id" binding:"required"`
	IncludeCircular     bool   `json:"include_circular"`
	IncludeArchitecture bool   `json:"include_architecture"`
	SeverityThreshold   string `json:"severity_threshold"`
	ConfigPath          string `json:"config_path"`
}

// StartComprehensiveAnalysis triggers comprehensive analysis on uploaded files by processing result ID
// @Summary Start comprehensive analysis
// @Description Triggers comprehensive dependency analysis on uploaded files
// @Tags analysis
// @Accept json
// @Produce json
// @Param uploadId path string true "Processing Result ID"
// @Success 202 {object} Response{data=models.DependencyAnalysis}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/analysis/comprehensive/{uploadId} [post]
func (h *AnalysisHandler) StartComprehensiveAnalysis(c *gin.Context) {
	uploadId := c.Param("uploadId")
	if uploadId == "" {
		BadRequest(c, "Upload ID is required", nil)
		return
	}

	h.logger.WithFields(logrus.Fields{
		"upload_id": uploadId,
	}).Info("Starting comprehensive analysis")

	// Use default analysis options
	options := services.AnalysisOptions{
		IncludeCircular:     true,
		IncludeArchitecture: true,
		SeverityThreshold:   "warning",
	}

	// Start analysis (this is async but we return immediately)
	analysis, err := h.integratedAnalysis.AnalyzeProcessingResult(
		c.Request.Context(),
		uploadId,
		options,
	)
	if err != nil {
		h.logger.WithError(err).Error("Failed to start comprehensive analysis")
		InternalError(c, "Failed to start analysis")
		return
	}

	// Return 202 Accepted with analysis ID
	c.JSON(202, gin.H{
		"success": true,
		"message": "Analysis started successfully",
		"data": gin.H{
			"id": analysis.ID,
		},
	})
}

// AnalyzeUploadedFiles triggers comprehensive analysis on uploaded files
// @Summary Analyze uploaded files
// @Description Triggers comprehensive dependency analysis on uploaded files
// @Tags analysis
// @Accept json
// @Produce json
// @Param request body AnalyzeUploadRequest true "Analysis request"
// @Success 202 {object} Response{data=models.DependencyAnalysis}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/analysis/upload [post]
func (h *AnalysisHandler) AnalyzeUploadedFiles(c *gin.Context) {
	var req AnalyzeUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Warning("Invalid request body")
		BadRequest(c, "Invalid request body", err.Error())
		return
	}

	// Set defaults
	if req.SeverityThreshold == "" {
		req.SeverityThreshold = "warning"
	}

	h.logger.WithFields(logrus.Fields{
		"processing_result_id": req.ProcessingResultID,
		"include_circular":     req.IncludeCircular,
		"include_architecture": req.IncludeArchitecture,
		"severity_threshold":   req.SeverityThreshold,
	}).Info("Starting analysis of uploaded files")

	// Convert to analysis options
	options := services.AnalysisOptions{
		IncludeCircular:     req.IncludeCircular,
		IncludeArchitecture: req.IncludeArchitecture,
		SeverityThreshold:   req.SeverityThreshold,
		ConfigPath:          req.ConfigPath,
	}

	// Start analysis (this is async but we return immediately)
	analysis, err := h.integratedAnalysis.AnalyzeProcessingResult(
		c.Request.Context(),
		req.ProcessingResultID,
		options,
	)
	if err != nil {
		h.logger.WithError(err).Error("Failed to analyze uploaded files")
		InternalError(c, "Failed to start analysis")
		return
	}

	// Return 202 Accepted with analysis ID
	c.JSON(202, gin.H{
		"success": true,
		"message": "Analysis started successfully",
		"data":    analysis,
	})
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