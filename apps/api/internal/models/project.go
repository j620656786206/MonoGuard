package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Status represents the status of an entity
type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

// Severity represents the severity level
type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

// RiskLevel represents the risk level
type RiskLevel string

const (
	RiskLevelLow      RiskLevel = "low"
	RiskLevelMedium   RiskLevel = "medium"
	RiskLevelHigh     RiskLevel = "high"
	RiskLevelCritical RiskLevel = "critical"
)

// Project represents a project in the system
type Project struct {
	ID             string                `json:"id" gorm:"primaryKey"`
	Name           string                `json:"name" gorm:"not null"`
	Description    *string               `json:"description,omitempty"`
	RepositoryURL  string                `json:"repositoryUrl" gorm:"column:repository_url;not null"`
	Branch         string                `json:"branch" gorm:"not null;default:\"main\""`
	Status         Status                `json:"status" gorm:"type:varchar(20);not null;default:\"pending\""`
	HealthScore    int                   `json:"healthScore" gorm:"column:health_score;default:0"`
	LastAnalysisAt *time.Time            `json:"lastAnalysisAt,omitempty" gorm:"column:last_analysis_at"`
	OwnerID        string                `json:"ownerId" gorm:"column:owner_id;not null"`
	Settings       *ProjectSettings      `json:"settings" gorm:"type:jsonb"`
	CreatedAt      time.Time             `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt      time.Time             `json:"updatedAt" gorm:"column:updated_at"`
	
	// Associations
	DependencyAnalyses     []DependencyAnalysis     `json:"dependencyAnalyses,omitempty" gorm:"foreignKey:ProjectID"`
	ArchitectureValidations []ArchitectureValidation `json:"architectureValidations,omitempty" gorm:"foreignKey:ProjectID"`
}

// ProjectSettings contains project-specific settings
type ProjectSettings struct {
	AutoAnalysis         bool                     `json:"autoAnalysis"`
	AnalysisSchedule     *string                  `json:"analysisSchedule,omitempty"`
	NotificationSettings NotificationSettings     `json:"notificationSettings"`
	ExcludePatterns      []string                 `json:"excludePatterns"`
	IncludePatterns      []string                 `json:"includePatterns"`
	ArchitectureRules    ArchitectureRules        `json:"architectureRules"`
}

// NotificationSettings contains notification configuration
type NotificationSettings struct {
	Email        bool       `json:"email"`
	Webhook      *string    `json:"webhook,omitempty"`
	SlackWebhook *string    `json:"slackWebhook,omitempty"`
	Severity     []Severity `json:"severity"`
}

// ArchitectureRules contains architecture validation rules
type ArchitectureRules struct {
	Layers []ArchitectureLayer `json:"layers"`
	Rules  []ArchitectureRule  `json:"rules"`
}

// ArchitectureLayer defines an architectural layer
type ArchitectureLayer struct {
	Name        string   `json:"name"`
	Pattern     string   `json:"pattern"`
	Description string   `json:"description"`
	CanImport   []string `json:"canImport"`
	CannotImport []string `json:"cannotImport"`
}

// ArchitectureRule defines an architectural rule
type ArchitectureRule struct {
	Name        string   `json:"name"`
	Severity    Severity `json:"severity"`
	Description string   `json:"description"`
	Pattern     *string  `json:"pattern,omitempty"`
	Enabled     bool     `json:"enabled"`
}

// BeforeCreate generates a UUID for the project before creating
func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

// TableName returns the table name for the Project model
func (Project) TableName() string {
	return "projects"
}