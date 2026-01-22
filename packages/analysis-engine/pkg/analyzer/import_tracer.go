// Package analyzer provides dependency graph analysis for monorepo workspaces.
// This file implements import tracing for circular dependencies (Story 3.2).
package analyzer

import (
	"path/filepath"
	"strings"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/parser"
	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// ========================================
// Import Tracer (Story 3.2)
// ========================================

// ImportTracer traces import statements that create circular dependencies.
type ImportTracer struct {
	workspace *types.WorkspaceData
	files     map[string][]byte // Source files (*.ts, *.js, *.tsx, *.jsx)
	parser    *parser.ImportParser
}

// NewImportTracer creates a new tracer for the given workspace and files.
func NewImportTracer(workspace *types.WorkspaceData, files map[string][]byte) *ImportTracer {
	return &ImportTracer{
		workspace: workspace,
		files:     files,
		parser:    parser.NewImportParser(),
	}
}

// Trace finds import statements that form the circular dependency.
// Returns empty slice (not nil) if no traces found or files are empty.
func (it *ImportTracer) Trace(cycle *types.CircularDependencyInfo) []types.ImportTrace {
	// Graceful degradation: return empty slice for invalid input
	if cycle == nil || len(cycle.Cycle) < 2 {
		return []types.ImportTrace{}
	}

	// Graceful degradation: return empty slice for no source files
	if it.files == nil || len(it.files) == 0 {
		return []types.ImportTrace{}
	}

	var traces []types.ImportTrace

	// Trace each edge in the cycle
	// Cycle format: [A, B, C, A] - has len-1 edges
	numEdges := len(cycle.Cycle) - 1
	for i := 0; i < numEdges; i++ {
		fromPkg := cycle.Cycle[i]
		toPkg := cycle.Cycle[i+1]

		edgeTraces := it.traceEdge(fromPkg, toPkg)
		traces = append(traces, edgeTraces...)
	}

	return traces
}

// traceEdge finds imports from one package to another.
func (it *ImportTracer) traceEdge(fromPkg, toPkg string) []types.ImportTrace {
	var traces []types.ImportTrace

	// Get source files for the "from" package
	sourceFiles := it.getSourceFilesForPackage(fromPkg)
	if len(sourceFiles) == 0 {
		return traces
	}

	// Create target packages map (just the "to" package)
	targets := map[string]bool{toPkg: true}

	// Parse each source file
	for filePath, content := range sourceFiles {
		fileTraces := it.parser.ParseFile(content, filePath, targets)

		// Set the FromPackage for each trace
		for i := range fileTraces {
			fileTraces[i].FromPackage = fromPkg
		}

		traces = append(traces, fileTraces...)
	}

	return traces
}

// getSourceFilesForPackage returns source files belonging to a package.
func (it *ImportTracer) getSourceFilesForPackage(pkgName string) map[string][]byte {
	sourceFiles := make(map[string][]byte)

	// Find the package path from workspace
	pkgInfo, exists := it.workspace.Packages[pkgName]
	if !exists {
		return sourceFiles
	}

	pkgPath := pkgInfo.Path

	// Find all source files that belong to this package
	for filePath, content := range it.files {
		// Check if file path starts with package path
		if !strings.HasPrefix(filePath, pkgPath+"/") && filePath != pkgPath {
			continue
		}

		// Check if it's a source file
		if !IsSourceFile(filePath) {
			continue
		}

		sourceFiles[filePath] = content
	}

	return sourceFiles
}

// IsSourceFile checks if a file is a parseable source file.
func IsSourceFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".ts", ".tsx", ".js", ".jsx", ".mjs", ".cjs":
		return true
	default:
		return false
	}
}
