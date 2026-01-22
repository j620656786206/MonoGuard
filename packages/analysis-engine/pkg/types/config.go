// Package types defines Go types that match TypeScript definitions in @monoguard/types.
// This file contains analysis configuration types for Story 2.6.
package types

// ========================================
// Analysis Configuration Types (Story 2.6)
// ========================================

// AnalysisConfig holds configuration options for analysis.
// Matches @monoguard/types AnalysisConfig interface.
type AnalysisConfig struct {
	Exclude []string `json:"exclude,omitempty"` // Exclusion patterns (exact, glob, or regex:)
}

// AnalysisInput represents the complete input to the analyze function.
// This is the top-level structure for WASM input.
type AnalysisInput struct {
	Files       map[string]string `json:"files"`                 // filename -> content (package.json files)
	SourceFiles map[string]string `json:"sourceFiles,omitempty"` // Story 3.2: Optional source files for import tracing
	Config      *AnalysisConfig   `json:"config,omitempty"`      // Optional configuration
}

// NewAnalysisConfig creates a new empty analysis config.
func NewAnalysisConfig() *AnalysisConfig {
	return &AnalysisConfig{
		Exclude: []string{},
	}
}
