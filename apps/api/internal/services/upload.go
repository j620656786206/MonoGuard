package services

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/monoguard/api/internal/models"
	"github.com/monoguard/api/pkg/database"
	"github.com/sirupsen/logrus"
)

// UploadService handles file upload operations
type UploadService struct {
	db     *database.DB
	logger *logrus.Logger
}

// NewUploadService creates a new upload service
func NewUploadService(db *database.DB, logger *logrus.Logger) *UploadService {
	return &UploadService{
		db:     db,
		logger: logger,
	}
}

// UploadFiles handles multiple file uploads
func (s *UploadService) UploadFiles(files []*multipart.FileHeader, uploadDir string) (*models.FileProcessingResult, error) {
	// Create processing result record
	processingResult := &models.FileProcessingResult{
		ID:          uuid.New().String(),
		ProcessedAt: time.Now(),
		Errors:      "",
	}

	var errors []string

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	var uploadedFiles []models.UploadedFile
	var packageJsonFiles []models.PackageJsonFile

	for _, fileHeader := range files {
		// Validate file
		if err := s.validateFile(fileHeader); err != nil {
			errors = append(errors, fmt.Sprintf("Invalid file %s: %v", fileHeader.Filename, err))
			continue
		}

		// Save file to disk
		uploadedFile, err := s.saveFile(fileHeader, uploadDir, processingResult.ID)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to save file %s: %v", fileHeader.Filename, err))
			continue
		}

		uploadedFiles = append(uploadedFiles, *uploadedFile)

		// Process file based on type
		if strings.ToLower(filepath.Ext(fileHeader.Filename)) == ".zip" {
			pkgFiles, err := s.processZipFile(uploadedFile.FilePath, processingResult.ID)
			if err != nil {
				errors = append(errors, fmt.Sprintf("Failed to process zip file %s: %v", fileHeader.Filename, err))
				continue
			}
			packageJsonFiles = append(packageJsonFiles, pkgFiles...)
		} else if strings.ToLower(filepath.Base(fileHeader.Filename)) == "package.json" || strings.ToLower(filepath.Ext(fileHeader.Filename)) == ".json" {
			pkgFile, err := s.processPackageJsonFile(uploadedFile.FilePath, processingResult.ID)
			if err != nil {
				errors = append(errors, fmt.Sprintf("Failed to process package.json file %s: %v", fileHeader.Filename, err))
				continue
			}
			packageJsonFiles = append(packageJsonFiles, *pkgFile)
		}
	}

	processingResult.Files = uploadedFiles
	processingResult.PackageJsonFiles = packageJsonFiles

	// Convert errors to JSON string
	if len(errors) > 0 {
		if errorsBytes, err := json.Marshal(errors); err == nil {
			processingResult.Errors = string(errorsBytes)
		}
	}

	// Save to database
	if err := s.db.Create(processingResult).Error; err != nil {
		return nil, fmt.Errorf("failed to save processing result to database: %w", err)
	}

	return processingResult, nil
}

// validateFile validates the uploaded file
func (s *UploadService) validateFile(fileHeader *multipart.FileHeader) error {
	// Check file size (50MB limit)
	const maxSize = 50 * 1024 * 1024 // 50MB
	if fileHeader.Size > maxSize {
		return fmt.Errorf("file size %d exceeds maximum allowed size %d", fileHeader.Size, maxSize)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	fileName := strings.ToLower(filepath.Base(fileHeader.Filename))
	
	if ext != ".zip" && fileName != "package.json" && ext != ".json" {
		return fmt.Errorf("unsupported file type: %s", ext)
	}

	return nil
}

// saveFile saves the uploaded file to disk
func (s *UploadService) saveFile(fileHeader *multipart.FileHeader, uploadDir, processingResultID string) (*models.UploadedFile, error) {
	// Open uploaded file
	src, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Generate unique filename
	fileID := uuid.New().String()
	fileName := fmt.Sprintf("%s_%s", fileID, fileHeader.Filename)
	filePath := filepath.Join(uploadDir, fileName)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	// Create uploaded file model
	uploadedFile := &models.UploadedFile{
		ID:                 fileID,
		OriginalName:       fileHeader.Filename,
		FileName:           fileName,
		FileSize:           fileHeader.Size,
		MimeType:           fileHeader.Header.Get("Content-Type"),
		FilePath:           filePath,
		Status:             string(models.UploadStatusUploaded),
		UploadedAt:         time.Now(),
		ProcessingResultID: processingResultID,
	}

	return uploadedFile, nil
}

// processZipFile extracts and processes a zip file
func (s *UploadService) processZipFile(zipPath, processingResultID string) ([]models.PackageJsonFile, error) {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open zip file: %w", err)
	}
	defer reader.Close()

	var packageJsonFiles []models.PackageJsonFile

	for _, file := range reader.File {
		if strings.ToLower(filepath.Base(file.Name)) == "package.json" {
			pkgFile, err := s.extractAndProcessPackageJson(file, processingResultID)
			if err != nil {
				s.logger.WithError(err).Warnf("Failed to process package.json from zip: %s", file.Name)
				continue
			}
			packageJsonFiles = append(packageJsonFiles, *pkgFile)
		}
	}

	return packageJsonFiles, nil
}

// extractAndProcessPackageJson extracts and processes a package.json file from zip
func (s *UploadService) extractAndProcessPackageJson(file *zip.File, processingResultID string) (*models.PackageJsonFile, error) {
	rc, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file in zip: %w", err)
	}
	defer rc.Close()

	content, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	return s.parsePackageJson(content, file.Name, processingResultID)
}

// processPackageJsonFile processes a standalone package.json file
func (s *UploadService) processPackageJsonFile(filePath, processingResultID string) (*models.PackageJsonFile, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read package.json file: %w", err)
	}

	return s.parsePackageJson(content, filepath.Base(filePath), processingResultID)
}

// parsePackageJson parses package.json content
func (s *UploadService) parsePackageJson(content []byte, path, processingResultID string) (*models.PackageJsonFile, error) {
	var pkgContent map[string]interface{}
	if err := json.Unmarshal(content, &pkgContent); err != nil {
		return nil, fmt.Errorf("failed to parse package.json: %w", err)
	}

	pkgFile := &models.PackageJsonFile{
		ID:                 uuid.New().String(),
		ProcessingResultID: processingResultID,
		Path:               path,
		Content:            string(content),
	}

	// Extract common fields
	if name, ok := pkgContent["name"].(string); ok {
		pkgFile.Name = &name
	}

	if version, ok := pkgContent["version"].(string); ok {
		pkgFile.Version = &version
	}

	if deps, ok := pkgContent["dependencies"].(map[string]interface{}); ok {
		dependencies := make(map[string]string)
		for k, v := range deps {
			if vStr, ok := v.(string); ok {
				dependencies[k] = vStr
			}
		}
		if depBytes, err := json.Marshal(dependencies); err == nil {
			pkgFile.Dependencies = string(depBytes)
		}
	}

	if devDeps, ok := pkgContent["devDependencies"].(map[string]interface{}); ok {
		devDependencies := make(map[string]string)
		for k, v := range devDeps {
			if vStr, ok := v.(string); ok {
				devDependencies[k] = vStr
			}
		}
		if devDepBytes, err := json.Marshal(devDependencies); err == nil {
			pkgFile.DevDependencies = string(devDepBytes)
		}
	}

	return pkgFile, nil
}

// GetProcessingResult retrieves a processing result by ID
func (s *UploadService) GetProcessingResult(id string) (*models.FileProcessingResult, error) {
	var result models.FileProcessingResult
	if err := s.db.Preload("Files").Preload("PackageJsonFiles").First(&result, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("failed to get processing result: %w", err)
	}

	return &result, nil
}

// CleanupOldFiles removes old uploaded files and their database records
func (s *UploadService) CleanupOldFiles(olderThan time.Duration) error {
	cutoffTime := time.Now().Add(-olderThan)

	// Get old files
	var oldFiles []models.UploadedFile
	if err := s.db.Where("created_at < ?", cutoffTime).Find(&oldFiles).Error; err != nil {
		return fmt.Errorf("failed to query old files: %w", err)
	}

	// Delete files from disk and database
	for _, file := range oldFiles {
		if err := os.Remove(file.FilePath); err != nil {
			s.logger.WithError(err).Warnf("Failed to delete file from disk: %s", file.FilePath)
		}

		if err := s.db.Delete(&file).Error; err != nil {
			s.logger.WithError(err).Warnf("Failed to delete file record from database: %s", file.ID)
		}
	}

	// Clean up old processing results
	if err := s.db.Where("created_at < ?", cutoffTime).Delete(&models.FileProcessingResult{}).Error; err != nil {
		return fmt.Errorf("failed to delete old processing results: %w", err)
	}

	s.logger.Infof("Cleaned up %d old files", len(oldFiles))
	return nil
}

// CreateProcessingResult creates a processing result from package.json files
func (s *UploadService) CreateProcessingResult(packageJsonFiles []models.PackageJsonFile) (*models.FileProcessingResult, error) {
	// Create processing result record
	processingResult := &models.FileProcessingResult{
		ID:          uuid.New().String(),
		ProcessedAt: time.Now(),
		Errors:      "",
		PackageJsonFiles: packageJsonFiles,
	}

	// Update processing result ID for all package.json files
	for i := range packageJsonFiles {
		packageJsonFiles[i].ID = uuid.New().String()
		packageJsonFiles[i].ProcessingResultID = processingResult.ID
	}

	// Save to database
	if err := s.db.Create(processingResult).Error; err != nil {
		return nil, fmt.Errorf("failed to save processing result to database: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"processing_result_id": processingResult.ID,
		"package_json_count":   len(packageJsonFiles),
	}).Info("Created processing result from GitHub")

	return processingResult, nil
}