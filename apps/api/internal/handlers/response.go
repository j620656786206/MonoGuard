package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// APIResponse represents a standard API response
type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Message   string      `json:"message,omitempty"`
	Timestamp string      `json:"timestamp"`
	Error     *APIError   `json:"error,omitempty"`
}

// APIError represents an API error
type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	APIResponse
	Pagination PaginationMeta `json:"pagination"`
}

// PaginationMeta contains pagination metadata
type PaginationMeta struct {
	CurrentPage     int  `json:"currentPage"`
	TotalPages      int  `json:"totalPages"`
	TotalItems      int64 `json:"totalItems"`
	ItemsPerPage    int  `json:"itemsPerPage"`
	HasNextPage     bool `json:"hasNextPage"`
	HasPreviousPage bool `json:"hasPreviousPage"`
}

// Success sends a successful response
func Success(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      data,
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// Created sends a created response
func Created(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusCreated, APIResponse{
		Success:   true,
		Data:      data,
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// SuccessPaginated sends a successful paginated response
func SuccessPaginated(c *gin.Context, data interface{}, pagination PaginationMeta, message string) {
	c.JSON(http.StatusOK, PaginatedResponse{
		APIResponse: APIResponse{
			Success:   true,
			Data:      data,
			Message:   message,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
		Pagination: pagination,
	})
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, errorCode, message string, details interface{}) {
	c.JSON(statusCode, APIResponse{
		Success:   false,
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Error: &APIError{
			Code:    errorCode,
			Message: message,
			Details: details,
		},
	})
}

// BadRequest sends a bad request error response
func BadRequest(c *gin.Context, message string, details interface{}) {
	Error(c, http.StatusBadRequest, "BAD_REQUEST", message, details)
}

// NotFound sends a not found error response
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, "NOT_FOUND", message, nil)
}

// InternalError sends an internal server error response
func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", message, nil)
}

// Unauthorized sends an unauthorized error response
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, "UNAUTHORIZED", message, nil)
}

// Forbidden sends a forbidden error response
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, "FORBIDDEN", message, nil)
}

// Conflict sends a conflict error response
func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, "CONFLICT", message, nil)
}

// ValidationError sends a validation error response
func ValidationError(c *gin.Context, message string, details interface{}) {
	Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", message, details)
}

// CalculatePagination calculates pagination metadata
func CalculatePagination(page, limit int, totalItems int64) PaginationMeta {
	totalPages := int((totalItems + int64(limit) - 1) / int64(limit))
	
	return PaginationMeta{
		CurrentPage:     page,
		TotalPages:      totalPages,
		TotalItems:      totalItems,
		ItemsPerPage:    limit,
		HasNextPage:     page < totalPages,
		HasPreviousPage: page > 1,
	}
}