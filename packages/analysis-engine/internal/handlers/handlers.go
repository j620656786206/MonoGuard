// Package handlers contains the business logic for WASM exported functions.
// This package is separated from cmd/wasm to allow unit testing without
// WASM build constraints.
package handlers

import (
	"github.com/j620656786206/MonoGuard/packages/analysis-engine/internal/result"
	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// Version is the current version of the analysis engine.
const Version = "0.1.0"

// GetVersion returns the analysis engine version as a Result JSON string.
func GetVersion() string {
	r := result.NewSuccess(types.VersionInfo{Version: Version})
	return r.ToJSON()
}

// Analyze performs dependency analysis on the provided workspace data.
// Returns a Result JSON string with AnalysisResult or error.
func Analyze(input string) string {
	if input == "" {
		r := result.NewError("INVALID_INPUT", "Missing JSON input")
		return r.ToJSON()
	}

	// Placeholder - will be implemented in Epic 2
	// TODO: Parse input JSON and perform actual analysis
	r := result.NewSuccess(types.AnalysisResult{
		HealthScore: 100,
		Packages:    0,
		Placeholder: true,
	})
	return r.ToJSON()
}

// Check validates the workspace configuration against configured rules.
// Returns a Result JSON string with CheckResult or error.
func Check(input string) string {
	if input == "" {
		r := result.NewError("INVALID_INPUT", "Missing JSON input")
		return r.ToJSON()
	}

	// Placeholder - will be implemented in Epic 2
	// TODO: Parse input JSON and perform actual validation
	r := result.NewSuccess(types.CheckResult{
		Passed:      true,
		Errors:      []string{},
		Placeholder: true,
	})
	return r.ToJSON()
}
