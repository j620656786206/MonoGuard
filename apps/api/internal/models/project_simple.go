package models

import "time"

// ProjectSimple - 簡化版本用於測試遷移問題
type ProjectSimple struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// TableName returns the table name
func (ProjectSimple) TableName() string {
	return "projects_simple"
}