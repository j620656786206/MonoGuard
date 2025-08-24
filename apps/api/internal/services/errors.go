package services

import "errors"

// Service errors
var (
	ErrProjectNotFound   = errors.New("project not found")
	ErrAnalysisNotFound  = errors.New("analysis not found")
	ErrInvalidInput      = errors.New("invalid input")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrDuplicateProject  = errors.New("project already exists")
)