package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/monoguard/api/internal/services"
	"github.com/sirupsen/logrus"
)

// UploadHandler handles file upload requests
type UploadHandler struct {
	uploadService *services.UploadService
	logger        *logrus.Logger
	uploadDir     string
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler(uploadService *services.UploadService, logger *logrus.Logger, uploadDir string) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
		logger:        logger,
		uploadDir:     uploadDir,
	}
}

// UploadFiles handles file upload requests
// @Summary Upload files
// @Description Upload multiple files (zip or package.json)
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "Files to upload"
// @Success 200 {object} Response{data=models.FileProcessingResult}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/upload [post]
func (h *UploadHandler) UploadFiles(c *gin.Context) {
	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		h.logger.WithError(err).Error("Failed to parse multipart form")
		Error(c, http.StatusBadRequest, "invalid_request", "Failed to parse multipart form", nil)
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		BadRequest(c, "No files provided", nil)
		return
	}

	// Upload files
	result, err := h.uploadService.UploadFiles(files, h.uploadDir)
	if err != nil {
		h.logger.WithError(err).Error("Failed to upload files")
		InternalError(c, fmt.Sprintf("Failed to upload files: %v", err))
		return
	}

	Success(c, result, fmt.Sprintf("Successfully uploaded %d files", len(result.Files)))
}

// GetUploadResult retrieves an upload result by ID
// @Summary Get upload result
// @Description Retrieve upload result and processing details by ID
// @Tags upload
// @Produce json
// @Param id path string true "Processing Result ID"
// @Success 200 {object} Response{data=models.FileProcessingResult}
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/upload/{id} [get]
func (h *UploadHandler) GetUploadResult(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "Processing result ID is required", nil)
		return
	}

	result, err := h.uploadService.GetProcessingResult(id)
	if err != nil {
		h.logger.WithError(err).Errorf("Failed to get processing result: %s", id)
		NotFound(c, "Processing result not found")
		return
	}

	Success(c, result, "")
}

// GetUploadedFile serves an uploaded file
// @Summary Download uploaded file
// @Description Download an uploaded file by its filename
// @Tags upload
// @Produce octet-stream
// @Param filename path string true "File name"
// @Success 200 {file} file
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/upload/files/{filename} [get]
func (h *UploadHandler) GetUploadedFile(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		BadRequest(c, "Filename is required", nil)
		return
	}

	filePath := filepath.Join(h.uploadDir, filename)
	
	// Security check - ensure file is within upload directory
	if !filepath.HasPrefix(filePath, h.uploadDir) {
		BadRequest(c, "Invalid file path", nil)
		return
	}

	c.File(filePath)
}

// CleanupOldFiles handles cleanup of old uploaded files
// @Summary Cleanup old files
// @Description Remove old uploaded files (older than specified duration)
// @Tags upload
// @Produce json
// @Param days query int false "Days to keep files (default: 7)"
// @Success 200 {object} Response{data=string}
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/upload/cleanup [post]
func (h *UploadHandler) CleanupOldFiles(c *gin.Context) {
	days := 7 // default
	if daysParam := c.Query("days"); daysParam != "" {
		if parsedDays, err := time.ParseDuration(daysParam + "d"); err == nil {
			days = int(parsedDays.Hours() / 24)
		}
	}

	duration := time.Duration(days) * 24 * time.Hour
	
	if err := h.uploadService.CleanupOldFiles(duration); err != nil {
		h.logger.WithError(err).Error("Failed to cleanup old files")
		InternalError(c, fmt.Sprintf("Failed to cleanup old files: %v", err))
		return
	}

	Success(c, fmt.Sprintf("Successfully cleaned up files older than %d days", days), "")
}