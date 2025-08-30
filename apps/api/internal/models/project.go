package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/monoguard/api/internal/constants"
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
	// Skip if ID already exists
	if p.ID != "" {
		return nil
	}
	
	// Ultra-safe approach: ONLY generate UUID if we're absolutely sure it's a real record creation
	// First line of defense: Check if hooks are explicitly disabled
	if tx.Statement != nil && tx.Statement.SkipHooks {
		return nil
	}
	
	// Second line of defense: Check global migration mode
	if constants.IsMigrationMode() {
		return nil
	}
	
	// Third line of defense: Multiple migration context checks
	stmt := tx.Statement
	if stmt != nil {
		// Skip if this is any kind of migration context
		if stmt.Context != nil {
			if stmt.Context.Value("gorm:auto_migrate") != nil ||
			   stmt.Context.Value("migration") != nil ||
			   stmt.Context.Value("gorm:migration") != nil {
				return nil
			}
		}
		
		// Skip if SQL is a SELECT query (inspection queries during migration)
		sql := stmt.SQL.String()
		if sql != "" && (
			strings.Contains(strings.ToUpper(sql), "SELECT") ||
			strings.Contains(strings.ToUpper(sql), "PRAGMA") ||
			strings.Contains(strings.ToUpper(sql), "SHOW") ||
			strings.Contains(strings.ToUpper(sql), "DESCRIBE") ||
			strings.Contains(strings.ToUpper(sql), "INFORMATION_SCHEMA")) {
			return nil
		}
	} else {
		// If no statement context, it's likely a migration
		return nil
	}
	
	// Only generate UUID for confirmed legitimate record creation with INSERT statement
	p.ID = uuid.New().String()
	return nil
}

// TableName returns the table name for the Project model
func (Project) TableName() string {
	return "projects"
}