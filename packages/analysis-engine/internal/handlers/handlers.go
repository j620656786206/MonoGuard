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
//
// # Input Formats (Backward Compatible)
//
// The function accepts two JSON input formats:
//
// ## Format 1: Legacy (pre-Story 2.6)
//
// A flat JSON object where keys are file paths and values are file contents.
// This format does NOT support exclusion patterns.
//
//	{
//	  "package.json": "{ \"name\": \"root\", \"workspaces\": [\"packages/*\"] }",
//	  "packages/pkg-a/package.json": "{ \"name\": \"@mono/pkg-a\", ... }"
//	}
//
// ## Format 2: AnalysisInput (Story 2.6+)
//
// A structured object with "files" and optional "config" fields.
// Use this format to specify exclusion patterns.
//
//	{
//	  "files": { "package.json": "...", ... },
//	  "config": { "exclude": ["packages/legacy", "regex:.*-test$"] }
//	}
//
// # Detection Logic
//
// The function first attempts to parse as AnalysisInput (Format 2).
// If the parsed object has a non-nil "files" field, it uses Format 2.
// Otherwise, it falls back to Format 1 (legacy flat map).
//
// Returns a Result JSON string with AnalysisResult (including dependency graph) or error.
func Analyze(input string) string {
	if input == "" {
		r := result.NewError(result.ErrInvalidInput, "Missing JSON input")
		return r.ToJSON()
	}

	// Parse input: try AnalysisInput format first, fallback to legacy format
	var analysisInput types.AnalysisInput
	var filesInput map[string]string
	var sourceFilesInput map[string]string
	var config *types.AnalysisConfig

	if err := json.Unmarshal([]byte(input), &analysisInput); err == nil && analysisInput.Files != nil {
		// Format 2: AnalysisInput with "files" field (and optional "config" and "sourceFiles")
		filesInput = analysisInput.Files
		sourceFilesInput = analysisInput.SourceFiles // Story 3.2: Optional source files
		config = analysisInput.Config
	} else {
		// Format 1: Legacy flat map (keys=paths, values=contents)
		// Note: This path also handles Format 2 parse failures gracefully
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

	// Story 3.2: Convert source files to []byte (if provided)
	var sourceFiles map[string][]byte
	if len(sourceFilesInput) > 0 {
		sourceFiles = make(map[string][]byte)
		for name, content := range sourceFilesInput {
			sourceFiles[name] = []byte(content)
		}
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

	// Story 3.2: Use AnalyzeWithSources to enable import tracing when source files provided
	analysisResult, err := a.AnalyzeWithSources(workspaceData, sourceFiles)
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
