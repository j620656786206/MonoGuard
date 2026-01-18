// Package handlers contains the business logic for WASM exported functions.
// This package is separated from cmd/wasm to allow unit testing without
// WASM build constraints.
package handlers

import (
	"encoding/json"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/internal/result"
	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/analyzer"
	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/parser"
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
// Input can be either:
// 1. Legacy format - JSON object mapping filenames to file contents:
//
//	{
//	  "package.json": "{ \"name\": \"root\", \"workspaces\": [\"packages/*\"] }",
//	  "packages/pkg-a/package.json": "{ \"name\": \"@mono/pkg-a\", ... }"
//	}
//
// 2. New format (Story 2.6) - AnalysisInput with optional config:
//
//	{
//	  "files": { "package.json": "...", ... },
//	  "config": { "exclude": ["packages/legacy", "regex:.*-test$"] }
//	}
//
// Returns a Result JSON string with AnalysisResult (including dependency graph) or error.
func Analyze(input string) string {
	if input == "" {
		r := result.NewError(result.ErrInvalidInput, "Missing JSON input")
		return r.ToJSON()
	}

	// Try to parse as AnalysisInput first (Story 2.6 format)
	var analysisInput types.AnalysisInput
	var filesInput map[string]string
	var config *types.AnalysisConfig

	if err := json.Unmarshal([]byte(input), &analysisInput); err == nil && analysisInput.Files != nil {
		// New format with optional config
		filesInput = analysisInput.Files
		config = analysisInput.Config
	} else {
		// Legacy format - just files map
		if err := json.Unmarshal([]byte(input), &filesInput); err != nil {
			r := result.NewError(result.ErrInvalidInput, "Failed to parse input JSON: "+err.Error())
			return r.ToJSON()
		}
	}

	// Convert string content to []byte
	files := make(map[string][]byte)
	for name, content := range filesInput {
		files[name] = []byte(content)
	}

	// Parse workspace using the real parser
	p := parser.NewParser("/workspace")
	workspaceData, err := p.Parse(files)
	if err != nil {
		r := result.NewError(result.ErrAnalysisFailed, err.Error())
		return r.ToJSON()
	}

	// Run analysis with config (Story 2.6: exclusion patterns)
	a, err := analyzer.NewAnalyzerWithConfig(config)
	if err != nil {
		r := result.NewError(result.ErrInvalidInput, "Invalid exclusion pattern: "+err.Error())
		return r.ToJSON()
	}

	analysisResult, err := a.Analyze(workspaceData)
	if err != nil {
		r := result.NewError(result.ErrAnalysisFailed, err.Error())
		return r.ToJSON()
	}

	r := result.NewSuccess(analysisResult)
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
