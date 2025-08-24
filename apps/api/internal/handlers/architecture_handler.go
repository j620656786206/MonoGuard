package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/monoguard/api/internal/services"
	"github.com/sirupsen/logrus"
)

// ArchitectureHandler handles architecture validation related requests
type ArchitectureHandler struct {
	layerValidator *services.LayerValidatorService
	logger         *logrus.Logger
}

// NewArchitectureHandler creates a new architecture handler
func NewArchitectureHandler(
	layerValidator *services.LayerValidatorService,
	logger *logrus.Logger,
) *ArchitectureHandler {
	return &ArchitectureHandler{
		layerValidator: layerValidator,
		logger:         logger,
	}
}

// ValidateArchitectureRequest represents the request body for architecture validation
type ValidateArchitectureRequest struct {
	ConfigPath        string   `json:"config_path"`
	TargetPackages    []string `json:"target_packages"`
	SeverityThreshold string   `json:"severity_threshold"`
}

// ValidateArchitecture validates project architecture
func (h *ArchitectureHandler) ValidateArchitecture(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		BadRequest(c, "Project ID is required", nil)
		return
	}

	var req ValidateArchitectureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Warning("Invalid request body")
		BadRequest(c, "Invalid request body", err.Error())
		return
	}

	// Set defaults
	if req.ConfigPath == "" {
		req.ConfigPath = ".monoguard.yml"
	}
	if req.SeverityThreshold == "" {
		req.SeverityThreshold = "warning"
	}

	h.logger.WithFields(logrus.Fields{
		"project_id":         projectID,
		"config_path":        req.ConfigPath,
		"target_packages":    req.TargetPackages,
		"severity_threshold": req.SeverityThreshold,
	}).Info("Validating architecture")

	// Trigger architecture validation
	validation, err := h.layerValidator.ValidateArchitecture(
		c.Request.Context(),
		projectID,
		req.ConfigPath,
		req.TargetPackages,
		req.SeverityThreshold,
	)
	if err != nil {
		h.logger.WithError(err).Error("Failed to validate architecture")
		InternalError(c, "Failed to validate architecture")
		return
	}

	Success(c, validation, "Architecture validation completed successfully")
}

// GetArchitectureGraph returns architecture dependency graph data
func (h *ArchitectureHandler) GetArchitectureGraph(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		BadRequest(c, "Project ID is required", nil)
		return
	}

	h.logger.WithField("project_id", projectID).Info("Generating architecture graph")

	// Mock graph data for demo
	graphData := map[string]interface{}{
		"nodes": []map[string]interface{}{
			{
				"id":               "apps/frontend",
				"name":             "Frontend App",
				"layer":            "Application Layer",
				"has_violations":   true,
				"violation_count":  1,
			},
			{
				"id":               "apps/api", 
				"name":             "API Server",
				"layer":            "Application Layer",
				"has_violations":   false,
				"violation_count":  0,
			},
			{
				"id":               "libs/ui",
				"name":             "UI Components",
				"layer":            "UI Component Library",
				"has_violations":   false,
				"violation_count":  0,
			},
			{
				"id":               "libs/shared-types",
				"name":             "Shared Types",
				"layer":            "Shared Libraries", 
				"has_violations":   false,
				"violation_count":  0,
			},
		},
		"edges": []map[string]interface{}{
			{
				"source":      "apps/frontend",
				"target":      "libs/ui",
				"type":        "dependency",
				"is_violation": false,
				"is_circular":  false,
			},
			{
				"source":      "apps/frontend",
				"target":      "libs/shared-types", 
				"type":        "dependency",
				"is_violation": false,
				"is_circular":  false,
			},
			{
				"source":      "apps/frontend",
				"target":      "apps/api", // This is a violation
				"type":        "dependency",
				"is_violation": true,
				"is_circular":  false,
			},
			{
				"source":      "libs/ui",
				"target":      "libs/shared-types",
				"type":        "dependency", 
				"is_violation": false,
				"is_circular":  false,
			},
		},
		"circular_paths": []map[string]interface{}{
			// No circular paths in this example
		},
	}

	Success(c, graphData, "Architecture graph generated successfully")
}

// GetArchitectureViolations returns architecture violations for a project
func (h *ArchitectureHandler) GetArchitectureViolations(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		BadRequest(c, "Project ID is required", nil)
		return
	}

	h.logger.WithField("project_id", projectID).Info("Retrieving architecture violations")

	// This would retrieve violations from the database
	// For now, return a placeholder response
	Error(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Feature not yet implemented", nil)
}

// ValidateConfig validates a configuration file
func (h *ArchitectureHandler) ValidateConfig(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		BadRequest(c, "Project ID is required", nil)
		return
	}

	var req struct {
		ConfigContent string `json:"config_content" validate:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Warning("Invalid request body")
		BadRequest(c, "Invalid request body", err.Error())
		return
	}

	h.logger.WithField("project_id", projectID).Info("Validating configuration")

	// Mock validation response
	validationResult := map[string]interface{}{
		"valid": true,
		"errors": []string{},
		"warnings": []string{},
		"suggestions": []string{
			"Consider adding more specific patterns for better layer separation",
		},
	}

	Success(c, validationResult, "Configuration validated successfully")
}