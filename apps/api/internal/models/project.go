package models

import (
	"time"
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

// Project represents a project in the system - simplified for Railway compatibility
type Project struct {
	ID             string    `json:"id" gorm:"primaryKey"`
	Name           string    `json:"name" gorm:"not null"`
	Description    string    `json:"description"`                     // Changed from *string to string
	RepositoryURL  string    `json:"repositoryUrl"`                   // Removed custom column name
	Branch         string    `json:"branch"`                          // Removed default value
	Status         string    `json:"status"`                          // Simplified, no type constraint or default
	HealthScore    int       `json:"healthScore"`                     // Removed custom column name
	OwnerID        string    `json:"ownerId"`                         // Removed custom column name
	CreatedAt      time.Time `json:"createdAt"`                       // Removed custom column name
	UpdatedAt      time.Time `json:"updatedAt"`                       // Removed custom column name
	
	// All complex fields temporarily removed for Railway compatibility
	// LastAnalysisAt *time.Time            `json:"lastAnalysisAt,omitempty" gorm:"column:last_analysis_at"`
	// Settings       *ProjectSettings      `json:"settings"`
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
	Severity     []string   `json:"severity"` // Temporarily changed to string slice
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
	Severity    string   `json:"severity"` // Temporarily changed to string
	Description string   `json:"description"`
	Pattern     *string  `json:"pattern,omitempty"`
	Enabled     bool     `json:"enabled"`
}

// Note: BeforeCreate hook removed to avoid Railway PostgreSQL migration issues
// UUID generation is now handled in service layer

// TableName returns the table name for the Project model
func (Project) TableName() string {
	return "projects"
}