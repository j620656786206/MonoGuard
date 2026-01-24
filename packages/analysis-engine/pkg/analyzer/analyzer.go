// Package analyzer provides dependency graph analysis for monorepo workspaces.
// This package implements the core analysis logic including:
//   - Building dependency graphs from workspace configuration
//   - Detecting circular dependencies (Story 2.3)
//   - Calculating architecture health scores (Story 2.5)
//   - Identifying duplicate dependencies with version conflicts (Story 2.4)
//   - Package exclusion patterns (Story 2.6)
//   - Root cause analysis for circular dependencies (Story 3.1)
//   - Import statement tracing for circular dependencies (Story 3.2)
//   - Fix strategy recommendations for circular dependencies (Story 3.3)
//   - Refactoring complexity calculation (Story 3.5)
package analyzer

import (
	"time"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// Analyzer orchestrates the complete workspace analysis process.
type Analyzer struct {
	graphBuilder *GraphBuilder
	config       *types.AnalysisConfig
}

// NewAnalyzer creates a new analyzer instance.
func NewAnalyzer() *Analyzer {
	return &Analyzer{
		graphBuilder: NewGraphBuilder(),
		config:       nil,
	}
}

// NewAnalyzerWithConfig creates an analyzer with the specified configuration.
// Returns an error if any exclusion regex pattern is invalid.
func NewAnalyzerWithConfig(config *types.AnalysisConfig) (*Analyzer, error) {
	if config == nil {
		return NewAnalyzer(), nil
	}

	graphBuilder, err := NewGraphBuilderWithExclusions(config.Exclude)
	if err != nil {
		return nil, err
	}

	return &Analyzer{
		graphBuilder: graphBuilder,
		config:       config,
	}, nil
}

// Analyze performs complete workspace analysis and returns the result.
// This builds the dependency graph, detects circular dependencies,
// identifies version conflicts, and calculates the architecture health score.
// Story 2.6: Excluded packages are marked in the graph but filtered from metrics.
func (a *Analyzer) Analyze(workspace *types.WorkspaceData) (*types.AnalysisResult, error) {
	// Build dependency graph (Story 2.2)
	// Story 2.6: Excluded packages are marked with Excluded=true
	graph, err := a.graphBuilder.Build(workspace)
	if err != nil {
		return nil, err
	}

	// Story 2.6: Count excluded and non-excluded packages
	excludedCount := 0
	for _, node := range graph.Nodes {
		if node.Excluded {
			excludedCount++
		}
	}
	packageCount := len(graph.Nodes) - excludedCount

	// Story 2.6: Create filtered graph for detectors (excludes excluded packages)
	filteredGraph := filterExcludedPackages(graph)

	// Detect circular dependencies (Story 2.3)
	// Story 2.6: Use filtered graph to exclude excluded packages
	cycleDetector := NewCycleDetector(filteredGraph)
	cycles := cycleDetector.DetectCycles()

	// Story 3.1: Enrich cycles with root cause analysis
	rootCauseAnalyzer := NewRootCauseAnalyzer(filteredGraph)
	for _, cycle := range cycles {
		cycle.RootCause = rootCauseAnalyzer.Analyze(cycle)
	}

	// Story 3.3: Generate fix strategies for cycles
	// Story 3.4: Generate step-by-step guides for each strategy
	// Story 3.5: Calculate refactoring complexity for cycles and strategies
	fixGenerator := NewFixStrategyGenerator(filteredGraph, workspace)
	guideGenerator := NewFixGuideGenerator(workspace)
	complexityCalc := NewComplexityCalculator(filteredGraph, workspace)
	for _, cycle := range cycles {
		// Calculate cycle-level complexity
		cycle.RefactoringComplexity = complexityCalc.Calculate(cycle)

		// Generate fix strategies
		strategies := fixGenerator.Generate(cycle)
		for i := range strategies {
			strategies[i].Guide = guideGenerator.Generate(cycle, &strategies[i])
			// Calculate per-strategy complexity (reuse cycle complexity as base)
			strategies[i].Complexity = complexityCalc.Calculate(cycle)
		}
		cycle.FixStrategies = strategies
	}

	// Detect version conflicts (Story 2.4)
	// Story 2.6: Use filtered graph to exclude excluded packages
	conflictDetector := NewConflictDetector(filteredGraph)
	conflicts := conflictDetector.DetectConflicts()

	// Calculate health score (Story 2.5)
	// Story 2.6: Use filtered graph to exclude excluded packages from metrics
	healthCalc := NewHealthCalculator(filteredGraph, cycles, conflicts)
	healthScore := healthCalc.Calculate()

	return &types.AnalysisResult{
		HealthScore:          healthScore.Overall,
		HealthScoreDetails:   healthScore,
		Packages:             packageCount,
		ExcludedPackages:     excludedCount,
		Graph:                graph, // Full graph with excluded flag for visualization
		CircularDependencies: cycles,
		VersionConflicts:     conflicts,
		CreatedAt:            time.Now().UTC().Format(time.RFC3339),
	}, nil
}

// AnalyzeWithSources performs complete workspace analysis with optional import tracing.
// Story 3.2: When sourceFiles are provided, import statements forming circular dependencies
// are traced and added to each CircularDependencyInfo.
// sourceFiles is optional - if nil or empty, analysis proceeds without import tracing.
func (a *Analyzer) AnalyzeWithSources(
	workspace *types.WorkspaceData,
	sourceFiles map[string][]byte,
) (*types.AnalysisResult, error) {
	// Build dependency graph (Story 2.2)
	// Story 2.6: Excluded packages are marked with Excluded=true
	graph, err := a.graphBuilder.Build(workspace)
	if err != nil {
		return nil, err
	}

	// Story 2.6: Count excluded and non-excluded packages
	excludedCount := 0
	for _, node := range graph.Nodes {
		if node.Excluded {
			excludedCount++
		}
	}
	packageCount := len(graph.Nodes) - excludedCount

	// Story 2.6: Create filtered graph for detectors (excludes excluded packages)
	filteredGraph := filterExcludedPackages(graph)

	// Detect circular dependencies (Story 2.3)
	// Story 2.6: Use filtered graph to exclude excluded packages
	cycleDetector := NewCycleDetector(filteredGraph)
	cycles := cycleDetector.DetectCycles()

	// Story 3.1: Enrich cycles with root cause analysis
	rootCauseAnalyzer := NewRootCauseAnalyzer(filteredGraph)
	for _, cycle := range cycles {
		cycle.RootCause = rootCauseAnalyzer.Analyze(cycle)
	}

	// Story 3.2: Enrich cycles with import traces
	// Always set ImportTraces (empty slice for graceful degradation per AC6)
	importTracer := NewImportTracer(workspace, sourceFiles)
	for _, cycle := range cycles {
		cycle.ImportTraces = importTracer.Trace(cycle)
	}

	// Story 3.3: Generate fix strategies for cycles
	// Story 3.4: Generate step-by-step guides for each strategy
	// Story 3.5: Calculate refactoring complexity for cycles and strategies
	fixGenerator := NewFixStrategyGenerator(filteredGraph, workspace)
	guideGenerator := NewFixGuideGenerator(workspace)
	complexityCalc := NewComplexityCalculator(filteredGraph, workspace)
	for _, cycle := range cycles {
		// Calculate cycle-level complexity
		cycle.RefactoringComplexity = complexityCalc.Calculate(cycle)

		// Generate fix strategies
		strategies := fixGenerator.Generate(cycle)
		for i := range strategies {
			strategies[i].Guide = guideGenerator.Generate(cycle, &strategies[i])
			// Calculate per-strategy complexity (reuse cycle complexity as base)
			strategies[i].Complexity = complexityCalc.Calculate(cycle)
		}
		cycle.FixStrategies = strategies
	}

	// Detect version conflicts (Story 2.4)
	// Story 2.6: Use filtered graph to exclude excluded packages
	conflictDetector := NewConflictDetector(filteredGraph)
	conflicts := conflictDetector.DetectConflicts()

	// Calculate health score (Story 2.5)
	// Story 2.6: Use filtered graph to exclude excluded packages from metrics
	healthCalc := NewHealthCalculator(filteredGraph, cycles, conflicts)
	healthScore := healthCalc.Calculate()

	return &types.AnalysisResult{
		HealthScore:          healthScore.Overall,
		HealthScoreDetails:   healthScore,
		Packages:             packageCount,
		ExcludedPackages:     excludedCount,
		Graph:                graph, // Full graph with excluded flag for visualization
		CircularDependencies: cycles,
		VersionConflicts:     conflicts,
		CreatedAt:            time.Now().UTC().Format(time.RFC3339),
	}, nil
}

// filterExcludedPackages creates a new graph with only non-excluded packages.
// This is used for metrics calculation while preserving the full graph for visualization.
func filterExcludedPackages(graph *types.DependencyGraph) *types.DependencyGraph {
	filtered := types.NewDependencyGraph(graph.RootPath, graph.WorkspaceType)

	// Copy only non-excluded nodes
	for name, node := range graph.Nodes {
		if !node.Excluded {
			// Create a new node with filtered dependencies
			newNode := types.NewPackageNode(node.Name, node.Version, node.Path)
			newNode.ExternalDeps = node.ExternalDeps
			newNode.ExternalDevDeps = node.ExternalDevDeps
			newNode.ExternalPeerDeps = node.ExternalPeerDeps
			newNode.ExternalOptionalDeps = node.ExternalOptionalDeps

			// Filter internal dependencies to exclude excluded packages
			for _, dep := range node.Dependencies {
				if depNode, ok := graph.Nodes[dep]; ok && !depNode.Excluded {
					newNode.Dependencies = append(newNode.Dependencies, dep)
				}
			}
			for _, dep := range node.DevDependencies {
				if depNode, ok := graph.Nodes[dep]; ok && !depNode.Excluded {
					newNode.DevDependencies = append(newNode.DevDependencies, dep)
				}
			}
			for _, dep := range node.PeerDependencies {
				if depNode, ok := graph.Nodes[dep]; ok && !depNode.Excluded {
					newNode.PeerDependencies = append(newNode.PeerDependencies, dep)
				}
			}
			for _, dep := range node.OptionalDependencies {
				if depNode, ok := graph.Nodes[dep]; ok && !depNode.Excluded {
					newNode.OptionalDependencies = append(newNode.OptionalDependencies, dep)
				}
			}

			filtered.Nodes[name] = newNode
		}
	}

	// Copy only edges between non-excluded packages
	for _, edge := range graph.Edges {
		fromNode, fromOk := graph.Nodes[edge.From]
		toNode, toOk := graph.Nodes[edge.To]
		if fromOk && toOk && !fromNode.Excluded && !toNode.Excluded {
			filtered.Edges = append(filtered.Edges, edge)
		}
	}

	return filtered
}
