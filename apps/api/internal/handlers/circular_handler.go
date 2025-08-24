package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/monoguard/api/internal/services"
	"github.com/sirupsen/logrus"
)

// CircularHandler handles circular dependency related requests
type CircularHandler struct {
	circularService *services.CircularDetectorService
	logger          *logrus.Logger
}

// NewCircularHandler creates a new circular handler
func NewCircularHandler(
	circularService *services.CircularDetectorService,
	logger *logrus.Logger,
) *CircularHandler {
	return &CircularHandler{
		circularService: circularService,
		logger:          logger,
	}
}

// DetectCircularDependenciesRequest represents the request body for circular dependency detection
type DetectCircularDependenciesRequest struct {
	RepoPath string `json:"repo_path" validate:"required"`
}

// DetectCircularDependencies detects circular dependencies for a project
func (h *CircularHandler) DetectCircularDependencies(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		BadRequest(c, "Project ID is required", nil)
		return
	}

	var req DetectCircularDependenciesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Warning("Invalid request body")
		BadRequest(c, "Invalid request body", err.Error())
		return
	}

	// Basic validation
	if req.RepoPath == "" {
		h.logger.Warning("Validation failed: repo_path is required")
		ValidationError(c, "Validation failed", "repo_path is required")
		return
	}

	h.logger.WithFields(logrus.Fields{
		"project_id": projectID,
		"repo_path":  req.RepoPath,
	}).Info("Detecting circular dependencies")

	// Trigger circular dependency detection
	analysis, err := h.circularService.DetectCircularDependencies(c.Request.Context(), projectID, req.RepoPath)
	if err != nil {
		h.logger.WithError(err).Error("Failed to detect circular dependencies")
		InternalError(c, "Failed to detect circular dependencies")
		return
	}

	Success(c, analysis, "Circular dependency detection completed successfully")
}

// GetCircularDependencyAnalysis retrieves a specific circular dependency analysis
func (h *CircularHandler) GetCircularDependencyAnalysis(c *gin.Context) {
	projectID := c.Param("id")
	analysisID := c.Param("analysis_id")

	if projectID == "" || analysisID == "" {
		BadRequest(c, "Project ID and Analysis ID are required", nil)
		return
	}

	h.logger.WithFields(logrus.Fields{
		"project_id":  projectID,
		"analysis_id": analysisID,
	}).Info("Retrieving circular dependency analysis")

	// This would retrieve the analysis from repository
	// For now, return a placeholder response
	Error(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Feature not yet implemented", nil)
}