// Package analyzer provides dependency graph analysis for monorepo workspaces.
// This package implements the core analysis logic including:
//   - Building dependency graphs from workspace configuration
//   - Detecting circular dependencies (Story 2.3)
//   - Calculating architecture health scores (Story 2.5)
//   - Identifying duplicate dependencies with version conflicts (Story 2.4)
package analyzer

import (
	"time"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// Analyzer orchestrates the complete workspace analysis process.
type Analyzer struct {
	graphBuilder *GraphBuilder
}

// NewAnalyzer creates a new analyzer instance.
func NewAnalyzer() *Analyzer {
	return &Analyzer{
		graphBuilder: NewGraphBuilder(),
	}
}

// Analyze performs complete workspace analysis and returns the result.
// This builds the dependency graph, detects circular dependencies,
// identifies version conflicts, and will include health score calculation in future stories.
func (a *Analyzer) Analyze(workspace *types.WorkspaceData) (*types.AnalysisResult, error) {
	// Build dependency graph (Story 2.2)
	graph, err := a.graphBuilder.Build(workspace)
	if err != nil {
		return nil, err
	}

	// Calculate package count
	packageCount := len(graph.Nodes)

	// Detect circular dependencies (Story 2.3)
	cycleDetector := NewCycleDetector(graph)
	cycles := cycleDetector.DetectCycles()

	// Detect version conflicts (Story 2.4)
	conflictDetector := NewConflictDetector(graph)
	conflicts := conflictDetector.DetectConflicts()

	// Return result with graph, cycles, and conflicts
	// Note: HealthScore is placeholder (100) until Story 2.5
	return &types.AnalysisResult{
		HealthScore:          100, // Placeholder until Story 2.5
		Packages:             packageCount,
		Graph:                graph,
		CircularDependencies: cycles,
		VersionConflicts:     conflicts,
		CreatedAt:            time.Now().UTC().Format(time.RFC3339),
	}, nil
}
