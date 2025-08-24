package handlers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/monoguard/api/internal/repository"
	"github.com/monoguard/api/internal/services"
	"github.com/sirupsen/logrus"
)

// ProjectHandler handles project-related HTTP requests
type ProjectHandler struct {
	projectService *services.ProjectService
	logger         *logrus.Logger
	validator      *validator.Validate
}

// NewProjectHandler creates a new project handler
func NewProjectHandler(projectService *services.ProjectService, logger *logrus.Logger) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
		logger:         logger,
		validator:      validator.New(),
	}
}

// CreateProject handles POST /projects
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req services.CreateProjectRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Warn("Invalid request body")
		BadRequest(c, "Invalid request body", err.Error())
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		h.logger.WithError(err).Warn("Validation failed")
		ValidationError(c, "Validation failed", err.Error())
		return
	}

	project, err := h.projectService.CreateProject(c.Request.Context(), &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create project")
		InternalError(c, "Failed to create project")
		return
	}

	Created(c, project, "Project created successfully")
}

// GetProject handles GET /projects/:id
func (h *ProjectHandler) GetProject(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "Project ID is required", nil)
		return
	}

	project, err := h.projectService.GetProject(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrProjectNotFound) {
			NotFound(c, "Project not found")
			return
		}
		h.logger.WithError(err).Error("Failed to get project")
		InternalError(c, "Failed to get project")
		return
	}

	Success(c, project, "Project retrieved successfully")
}

// GetProjects handles GET /projects
func (h *ProjectHandler) GetProjects(c *gin.Context) {
	params := h.parseQueryParams(c)

	projects, total, err := h.projectService.GetProjects(c.Request.Context(), params)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get projects")
		InternalError(c, "Failed to get projects")
		return
	}

	pagination := CalculatePagination(params.Page, params.Limit, total)
	SuccessPaginated(c, projects, pagination, "Projects retrieved successfully")
}

// GetProjectsByOwner handles GET /projects/owner/:ownerId
func (h *ProjectHandler) GetProjectsByOwner(c *gin.Context) {
	ownerID := c.Param("ownerId")
	if ownerID == "" {
		BadRequest(c, "Owner ID is required", nil)
		return
	}

	params := h.parseQueryParams(c)

	projects, total, err := h.projectService.GetProjectsByOwner(c.Request.Context(), ownerID, params)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get projects by owner")
		InternalError(c, "Failed to get projects by owner")
		return
	}

	pagination := CalculatePagination(params.Page, params.Limit, total)
	SuccessPaginated(c, projects, pagination, "Projects retrieved successfully")
}

// UpdateProject handles PUT /projects/:id
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "Project ID is required", nil)
		return
	}

	var req services.UpdateProjectRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Warn("Invalid request body")
		BadRequest(c, "Invalid request body", err.Error())
		return
	}

	project, err := h.projectService.UpdateProject(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, services.ErrProjectNotFound) {
			NotFound(c, "Project not found")
			return
		}
		h.logger.WithError(err).Error("Failed to update project")
		InternalError(c, "Failed to update project")
		return
	}

	Success(c, project, "Project updated successfully")
}

// DeleteProject handles DELETE /projects/:id
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "Project ID is required", nil)
		return
	}

	if err := h.projectService.DeleteProject(c.Request.Context(), id); err != nil {
		if errors.Is(err, services.ErrProjectNotFound) {
			NotFound(c, "Project not found")
			return
		}
		h.logger.WithError(err).Error("Failed to delete project")
		InternalError(c, "Failed to delete project")
		return
	}

	Success(c, nil, "Project deleted successfully")
}

// TriggerAnalysis handles POST /projects/:id/analyze
func (h *ProjectHandler) TriggerAnalysis(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "Project ID is required", nil)
		return
	}

	var req struct {
		RepoPath string `json:"repoPath" validate:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Warn("Invalid request body")
		BadRequest(c, "Invalid request body", err.Error())
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		h.logger.WithError(err).Warn("Validation failed")
		ValidationError(c, "Validation failed", err.Error())
		return
	}

	analysis, err := h.projectService.TriggerAnalysis(c.Request.Context(), id, req.RepoPath)
	if err != nil {
		if errors.Is(err, services.ErrProjectNotFound) {
			NotFound(c, "Project not found")
			return
		}
		h.logger.WithError(err).Error("Failed to trigger analysis")
		InternalError(c, "Failed to trigger analysis")
		return
	}

	Success(c, analysis, "Analysis triggered successfully")
}

// parseQueryParams parses common query parameters
func (h *ProjectHandler) parseQueryParams(c *gin.Context) *repository.QueryParams {
	params := repository.DefaultQueryParams()

	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			params.Page = p
		}
	}

	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			params.Limit = l
		}
	}

	if sort := c.Query("sort"); sort != "" {
		params.Sort = sort
	}

	if order := c.Query("order"); order != "" {
		params.Order = order
	}

	if search := c.Query("search"); search != "" {
		params.Search = search
	}

	if status := c.Query("status"); status != "" {
		params.Status = status
	}

	return params
}