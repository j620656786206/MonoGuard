package models

import (
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
	Status               string    `json:"status" gorm:"not null;type:varchar(20);default:'uploaded'"`
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

// TableName overrides the table name for PackageJsonFile
func (PackageJsonFile) TableName() string {
	return "package_json_files"
}