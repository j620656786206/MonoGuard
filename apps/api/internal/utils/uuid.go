package utils

import (
	"time"

	"github.com/google/uuid"
	"github.com/monoguard/api/internal/models"
)

// GenerateProjectID generates a UUID for a new project
func GenerateProjectID(project *models.Project) {
	if project.ID == "" {
		project.ID = uuid.New().String()
	}
}

// GenerateDependencyAnalysisID generates a UUID for a new dependency analysis
func GenerateDependencyAnalysisID(analysis *models.DependencyAnalysis) {
	if analysis.ID == "" {
		analysis.ID = uuid.New().String()
		if analysis.StartedAt.IsZero() {
			analysis.StartedAt = time.Now()
		}
	}
}

// GenerateArchitectureValidationID generates a UUID for a new architecture validation
func GenerateArchitectureValidationID(validation *models.ArchitectureValidation) {
	if validation.ID == "" {
		validation.ID = uuid.New().String()
		if validation.StartedAt.IsZero() {
			validation.StartedAt = time.Now()
		}
	}
}

// GenerateHealthScoreID generates a UUID for a new health score
func GenerateHealthScoreID(score *models.HealthScore) {
	if score.ID == "" {
		score.ID = uuid.New().String()
		if score.LastUpdated.IsZero() {
			score.LastUpdated = time.Now()
		}
	}
}

// GeneratePackageJSONAnalysisID generates a UUID for a new package.json analysis
func GeneratePackageJSONAnalysisID(analysis *models.PackageJSONAnalysis) {
	if analysis.ID == "" {
		analysis.ID = uuid.New().String()
		if analysis.StartedAt.IsZero() {
			analysis.StartedAt = time.Now()
		}
	}
}

// GenerateUploadedFileID generates a UUID for a new uploaded file
func GenerateUploadedFileID(file *models.UploadedFile) {
	if file.ID == "" {
		file.ID = uuid.New().String()
	}
}

// GenerateFileProcessingResultID generates a UUID for a new file processing result
func GenerateFileProcessingResultID(result *models.FileProcessingResult) {
	if result.ID == "" {
		result.ID = uuid.New().String()
	}
}

// GeneratePackageJsonFileID generates a UUID for a new package.json file
func GeneratePackageJsonFileID(file *models.PackageJsonFile) {
	if file.ID == "" {
		file.ID = uuid.New().String()
	}
}