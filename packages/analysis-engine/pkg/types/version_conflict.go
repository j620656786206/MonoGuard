// Package types defines Go types that match TypeScript definitions in @monoguard/types.
// This file contains version conflict types for Story 2.4.
package types

// ========================================
// Version Conflict Types (Story 2.4)
// ========================================

// VersionConflictInfo represents a dependency with multiple versions across packages.
// Matches @monoguard/types VersionConflict interface with extended severity classification.
type VersionConflictInfo struct {
	PackageName         string                `json:"packageName"`
	ConflictingVersions []*ConflictingVersion `json:"conflictingVersions"`
	Severity            ConflictSeverity      `json:"severity"`
	Resolution          string                `json:"resolution"`
	Impact              string                `json:"impact"`
}

// ConflictingVersion represents one version and which packages use it.
type ConflictingVersion struct {
	Version    string   `json:"version"`
	Packages   []string `json:"packages"`   // Workspace packages using this version
	IsBreaking bool     `json:"isBreaking"` // True if major version differs from others
	DepType    string   `json:"depType"`    // "production", "development", "peer"
}

// ConflictSeverity indicates how serious the version mismatch is.
// Based on semver differences as specified in Story 2.4 AC3.
type ConflictSeverity string

const (
	// ConflictSeverityCritical indicates major version differences (e.g., v3 vs v4).
	// Breaking changes are likely.
	ConflictSeverityCritical ConflictSeverity = "critical"

	// ConflictSeverityWarning indicates minor version differences (e.g., v4.17 vs v4.18).
	// New features may have been added, possible issues.
	ConflictSeverityWarning ConflictSeverity = "warning"

	// ConflictSeverityInfo indicates patch version differences (e.g., v4.17.21 vs v4.17.19).
	// Bug fixes only, generally safe.
	ConflictSeverityInfo ConflictSeverity = "info"
)

// DepType constants for dependency classification.
const (
	DepTypeProduction  = "production"
	DepTypeDevelopment = "development"
	DepTypePeer        = "peer"
)
