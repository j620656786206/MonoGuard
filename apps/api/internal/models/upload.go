package models

import (
	"encoding/json"
	"time"
	
	"gorm.io/gorm"
)

// UploadedFile represents an uploaded file in the system
type UploadedFile struct {
	ID                   string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	OriginalName         string    `json:"originalName" gorm:"not null;type:varchar(255)"`
	FileName             string    `json:"fileName" gorm:"not null;type:varchar(255);unique"`
	FileSize             int64     `json:"fileSize" gorm:"not null"`
	MimeType             string    `json:"mimeType" gorm:"not null;type:varchar(100)"`
	FilePath             string    `json:"filePath" gorm:"not null;type:varchar(500)"`
	Status               string    `json:"status" gorm:"not null;type:varchar(20);default:\"uploaded\""`
	UploadedAt           time.Time `json:"uploadedAt" gorm:"not null;autoCreateTime"`
	ProcessedAt          *time.Time `json:"processedAt,omitempty"`
	ProcessingResultID   string    `json:"processingResultId" gorm:"type:varchar(36)"`
	CreatedAt            time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt            time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt            gorm.DeletedAt `json:"-" gorm:"index"`
}

// FileProcessingResult represents the result of file processing
type FileProcessingResult struct {
	ID               string              `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Files            []UploadedFile      `json:"files" gorm:"foreignKey:ProcessingResultID"`
	PackageJsonFiles []PackageJsonFile   `json:"packageJsonFiles" gorm:"foreignKey:ProcessingResultID"`
	ProcessedAt      time.Time           `json:"processedAt" gorm:"not null"`
	Errors           string              `json:"errors,omitempty" gorm:"type:text"`
	CreatedAt        time.Time           `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt        time.Time           `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt        gorm.DeletedAt      `json:"-" gorm:"index"`
}

// PackageJsonFile represents a package.json file found in the uploaded files
type PackageJsonFile struct {
	ID                   string             `json:"id" gorm:"primaryKey;type:varchar(36)"`
	ProcessingResultID   string             `json:"processingResultId" gorm:"type:varchar(36)"`
	Path                 string             `json:"path" gorm:"not null;type:varchar(500)"`
	Content              string             `json:"content" gorm:"type:text"`
	Name                 *string            `json:"name,omitempty" gorm:"type:varchar(255)"`
	Version              *string            `json:"version,omitempty" gorm:"type:varchar(50)"`
	Dependencies         string             `json:"dependencies,omitempty" gorm:"type:text"`
	DevDependencies      string             `json:"devDependencies,omitempty" gorm:"type:text"`
	CreatedAt            time.Time          `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt            time.Time          `json:"updatedAt" gorm:"autoUpdateTime"`
}

// UploadStatus defines the status of uploaded files
type UploadStatus string

const (
	UploadStatusUploaded   UploadStatus = "uploaded"
	UploadStatusProcessing UploadStatus = "processing"
	UploadStatusCompleted  UploadStatus = "completed"
	UploadStatusError      UploadStatus = "error"
)

// TableName overrides the table name for UploadedFile
func (UploadedFile) TableName() string {
	return "uploaded_files"
}

// TableName overrides the table name for FileProcessingResult
func (FileProcessingResult) TableName() string {
	return "file_processing_results"
}

// PackageJsonFileResponse is used for API responses with proper JSON structure
type PackageJsonFileResponse struct {
	ID                   string                 `json:"id"`
	ProcessingResultID   string                 `json:"processingResultId"`
	Path                 string                 `json:"path"`
	Content              map[string]interface{} `json:"content"`
	Name                 *string                `json:"name,omitempty"`
	Version              *string                `json:"version,omitempty"`
	Dependencies         map[string]string      `json:"dependencies,omitempty"`
	DevDependencies      map[string]string      `json:"devDependencies,omitempty"`
	CreatedAt            time.Time              `json:"createdAt"`
	UpdatedAt            time.Time              `json:"updatedAt"`
}

// ToResponse converts PackageJsonFile to response format with parsed JSON fields
func (p PackageJsonFile) ToResponse() PackageJsonFileResponse {
	response := PackageJsonFileResponse{
		ID:                 p.ID,
		ProcessingResultID: p.ProcessingResultID,
		Path:               p.Path,
		Name:               p.Name,
		Version:            p.Version,
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
		Dependencies:       make(map[string]string),
		DevDependencies:    make(map[string]string),
	}

	// Parse content JSON
	if p.Content != "" {
		var content map[string]interface{}
		if err := json.Unmarshal([]byte(p.Content), &content); err == nil {
			response.Content = content
		}
	}

	// Parse dependencies JSON
	if p.Dependencies != "" {
		var deps map[string]string
		if err := json.Unmarshal([]byte(p.Dependencies), &deps); err == nil {
			response.Dependencies = deps
		}
	}

	// Parse devDependencies JSON
	if p.DevDependencies != "" {
		var devDeps map[string]string
		if err := json.Unmarshal([]byte(p.DevDependencies), &devDeps); err == nil {
			response.DevDependencies = devDeps
		}
	}

	return response
}

// FileProcessingResultResponse is used for API responses with properly converted PackageJsonFiles
type FileProcessingResultResponse struct {
	ID               string                     `json:"id"`
	Files            []UploadedFile             `json:"files"`
	PackageJsonFiles []PackageJsonFileResponse  `json:"packageJsonFiles"`
	ProcessedAt      time.Time                  `json:"processedAt"`
	Errors           []string                   `json:"errors,omitempty"`
	CreatedAt        time.Time                  `json:"createdAt"`
	UpdatedAt        time.Time                  `json:"updatedAt"`
}

// ToResponse converts FileProcessingResult to response format
func (f FileProcessingResult) ToResponse() FileProcessingResultResponse {
	response := FileProcessingResultResponse{
		ID:          f.ID,
		Files:       f.Files,
		ProcessedAt: f.ProcessedAt,
		CreatedAt:   f.CreatedAt,
		UpdatedAt:   f.UpdatedAt,
	}

	// Convert PackageJsonFiles to response format
	response.PackageJsonFiles = make([]PackageJsonFileResponse, len(f.PackageJsonFiles))
	for i, pkg := range f.PackageJsonFiles {
		response.PackageJsonFiles[i] = pkg.ToResponse()
	}

	// Parse errors string to array if not empty
	if f.Errors != "" {
		var errors []string
		if err := json.Unmarshal([]byte(f.Errors), &errors); err == nil {
			response.Errors = errors
		} else {
			// If it's not JSON, treat as single error
			response.Errors = []string{f.Errors}
		}
	}

	return response
}

// BeforeCreate generates UUIDs for upload models before creating
// Note: BeforeCreate hook removed to avoid Railway PostgreSQL migration issues
// UUID generation is now handled in service layer

// Note: BeforeCreate hook removed to avoid Railway PostgreSQL migration issues
// UUID generation is now handled in service layer

// Note: BeforeCreate hook removed to avoid Railway PostgreSQL migration issues
// UUID generation is now handled in service layer

// TableName overrides the table name for PackageJsonFile
func (PackageJsonFile) TableName() string {
	return "package_json_files"
}