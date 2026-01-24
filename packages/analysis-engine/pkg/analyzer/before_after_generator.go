// Package analyzer provides dependency graph analysis for monorepo workspaces.
// This file contains before/after explanation generator for Story 3.7.
package analyzer

import (
	"fmt"
	"strings"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// ========================================
// Before/After Generator (Story 3.7)
// ========================================

// BeforeAfterGenerator creates before/after comparisons for fix strategies.
type BeforeAfterGenerator struct {
	graph     *types.DependencyGraph
	workspace *types.WorkspaceData
}

// NewBeforeAfterGenerator creates a new generator.
func NewBeforeAfterGenerator(graph *types.DependencyGraph, workspace *types.WorkspaceData) *BeforeAfterGenerator {
	return &BeforeAfterGenerator{
		graph:     graph,
		workspace: workspace,
	}
}

// Generate creates the before/after explanation for a strategy.
func (bag *BeforeAfterGenerator) Generate(
	cycle *types.CircularDependencyInfo,
	strategy *types.FixStrategy,
) *types.BeforeAfterExplanation {
	if cycle == nil || strategy == nil {
		return nil
	}

	return &types.BeforeAfterExplanation{
		CurrentState:     bag.generateCurrentState(cycle),
		ProposedState:    bag.generateProposedState(cycle, strategy),
		PackageJsonDiffs: bag.generatePackageJsonDiffs(strategy),
		ImportDiffs:      bag.generateImportDiffs(cycle, strategy),
		Explanation:      bag.generateExplanation(cycle, strategy),
		Warnings:         bag.generateWarnings(cycle, strategy),
	}
}

// ========================================
// Current State Generation (Task 3, AC #1)
// ========================================

// generateCurrentState creates the "before" diagram.
func (bag *BeforeAfterGenerator) generateCurrentState(cycle *types.CircularDependencyInfo) *types.StateDiagram {
	nodes := []types.DiagramNode{}
	edges := []types.DiagramEdge{}

	// Create set of packages in cycle for quick lookup
	cycleSet := make(map[string]bool)
	for i := 0; i < len(cycle.Cycle)-1; i++ {
		cycleSet[cycle.Cycle[i]] = true
	}

	// Add nodes for all packages involved
	for pkgName := range cycleSet {
		nodes = append(nodes, types.DiagramNode{
			ID:        pkgName,
			Label:     extractShortName(pkgName),
			IsInCycle: true,
			IsNew:     false,
			NodeType:  types.NodeTypeCycle,
		})
	}

	// Add edges for the cycle path
	for i := 0; i < len(cycle.Cycle)-1; i++ {
		from := cycle.Cycle[i]
		to := cycle.Cycle[i+1]
		edges = append(edges, types.DiagramEdge{
			From:      from,
			To:        to,
			IsInCycle: true,
			IsRemoved: false,
			IsNew:     false,
			EdgeType:  types.EdgeTypeCycle,
		})
	}

	return &types.StateDiagram{
		Nodes:           nodes,
		Edges:           edges,
		HighlightedPath: cycle.Cycle,
		CycleResolved:   false,
	}
}

// ========================================
// Proposed State Generation (Task 4, AC #2)
// ========================================

// generateProposedState creates the "after" diagram.
func (bag *BeforeAfterGenerator) generateProposedState(
	cycle *types.CircularDependencyInfo,
	strategy *types.FixStrategy,
) *types.StateDiagram {
	switch strategy.Type {
	case types.FixStrategyExtractModule:
		return bag.generateProposedStateExtractModule(cycle, strategy)
	case types.FixStrategyDependencyInject:
		return bag.generateProposedStateDI(cycle, strategy)
	case types.FixStrategyBoundaryRefactor:
		return bag.generateProposedStateBoundary(cycle, strategy)
	default:
		return types.NewStateDiagram()
	}
}

// generateProposedStateExtractModule creates proposed state for Extract Module strategy.
func (bag *BeforeAfterGenerator) generateProposedStateExtractModule(
	cycle *types.CircularDependencyInfo,
	strategy *types.FixStrategy,
) *types.StateDiagram {
	nodes := []types.DiagramNode{}
	edges := []types.DiagramEdge{}

	// Add existing packages (no longer in cycle)
	for _, pkgName := range strategy.TargetPackages {
		nodes = append(nodes, types.DiagramNode{
			ID:        pkgName,
			Label:     extractShortName(pkgName),
			IsInCycle: false,
			IsNew:     false,
			NodeType:  types.NodeTypeAffected,
		})
	}

	// Add new shared package
	if strategy.NewPackageName != "" {
		nodes = append(nodes, types.DiagramNode{
			ID:        strategy.NewPackageName,
			Label:     extractShortName(strategy.NewPackageName),
			IsInCycle: false,
			IsNew:     true,
			NodeType:  types.NodeTypeNew,
		})

		// Add new edges to shared package
		for _, pkgName := range strategy.TargetPackages {
			edges = append(edges, types.DiagramEdge{
				From:      pkgName,
				To:        strategy.NewPackageName,
				IsInCycle: false,
				IsRemoved: false,
				IsNew:     true,
				EdgeType:  types.EdgeTypeNew,
			})
		}
	}

	return &types.StateDiagram{
		Nodes:         nodes,
		Edges:         edges,
		CycleResolved: true,
	}
}

// generateProposedStateDI creates proposed state for Dependency Injection strategy.
func (bag *BeforeAfterGenerator) generateProposedStateDI(
	cycle *types.CircularDependencyInfo,
	strategy *types.FixStrategy,
) *types.StateDiagram {
	nodes := []types.DiagramNode{}
	edges := []types.DiagramEdge{}

	// Determine which packages are involved
	targetPkgs := strategy.TargetPackages
	if len(targetPkgs) < 2 && len(cycle.Cycle) >= 2 {
		targetPkgs = cycle.Cycle[:len(cycle.Cycle)-1]
	}

	// Add nodes for all packages (affected but no longer in cycle)
	for _, pkgName := range targetPkgs {
		nodes = append(nodes, types.DiagramNode{
			ID:        pkgName,
			Label:     extractShortName(pkgName),
			IsInCycle: false,
			IsNew:     false,
			NodeType:  types.NodeTypeAffected,
		})
	}

	// In DI, the dependency is inverted - edges change direction
	// The problematic direct import is removed
	// Show remaining edges (dependencies that don't form cycle)
	if len(targetPkgs) >= 2 {
		// First package no longer directly imports second
		// Second may still be used but through injection
		// Add edge showing interface dependency
		edges = append(edges, types.DiagramEdge{
			From:      targetPkgs[0],
			To:        targetPkgs[1],
			IsInCycle: false,
			IsRemoved: true, // Direct import removed
			IsNew:     false,
			EdgeType:  types.EdgeTypeRemoved,
		})
	}

	return &types.StateDiagram{
		Nodes:         nodes,
		Edges:         edges,
		CycleResolved: true,
	}
}

// generateProposedStateBoundary creates proposed state for Boundary Refactoring strategy.
func (bag *BeforeAfterGenerator) generateProposedStateBoundary(
	cycle *types.CircularDependencyInfo,
	strategy *types.FixStrategy,
) *types.StateDiagram {
	nodes := []types.DiagramNode{}
	edges := []types.DiagramEdge{}

	// Determine which packages are involved
	targetPkgs := strategy.TargetPackages
	if len(targetPkgs) < 2 && len(cycle.Cycle) >= 2 {
		targetPkgs = cycle.Cycle[:len(cycle.Cycle)-1]
	}

	// Add nodes for all packages (affected, restructured)
	for _, pkgName := range targetPkgs {
		nodes = append(nodes, types.DiagramNode{
			ID:        pkgName,
			Label:     extractShortName(pkgName),
			IsInCycle: false,
			IsNew:     false,
			NodeType:  types.NodeTypeAffected,
		})
	}

	// In boundary refactoring, dependencies become one-directional
	// Only show edges that don't create a cycle
	if len(targetPkgs) >= 2 {
		// After refactoring, only one direction of dependency remains
		edges = append(edges, types.DiagramEdge{
			From:      targetPkgs[0],
			To:        targetPkgs[1],
			IsInCycle: false,
			IsRemoved: false,
			IsNew:     false,
			EdgeType:  types.EdgeTypeUnchanged,
		})
	}

	return &types.StateDiagram{
		Nodes:         nodes,
		Edges:         edges,
		CycleResolved: true,
	}
}

// ========================================
// Package.json Diff Generation (Task 5, AC #3)
// ========================================

// generatePackageJsonDiffs creates package.json change descriptions.
func (bag *BeforeAfterGenerator) generatePackageJsonDiffs(strategy *types.FixStrategy) []types.PackageJsonDiff {
	diffs := []types.PackageJsonDiff{}

	switch strategy.Type {
	case types.FixStrategyExtractModule:
		// Each target package needs to add dependency on new shared package
		for _, pkgName := range strategy.TargetPackages {
			diffs = append(diffs, types.PackageJsonDiff{
				PackageName: pkgName,
				FilePath:    bag.getPackageJsonPath(pkgName),
				DependenciesToAdd: []types.DependencyChange{
					{
						Name:           strategy.NewPackageName,
						Version:        "workspace:*",
						DependencyType: "dependencies",
					},
				},
				DependenciesToRemove: []types.DependencyChange{},
				Summary:              fmt.Sprintf("Add dependency on %s", strategy.NewPackageName),
			})
		}

	case types.FixStrategyDependencyInject:
		// Dependency injection typically removes direct imports, not package.json deps
		if len(strategy.TargetPackages) > 0 {
			diffs = append(diffs, types.PackageJsonDiff{
				PackageName:          strategy.TargetPackages[0],
				FilePath:             bag.getPackageJsonPath(strategy.TargetPackages[0]),
				DependenciesToAdd:    []types.DependencyChange{},
				DependenciesToRemove: []types.DependencyChange{},
				Summary:              "No package.json changes required for dependency injection",
			})
		}

	case types.FixStrategyBoundaryRefactor:
		// Boundary refactoring may involve moving code between packages
		for _, pkgName := range strategy.TargetPackages {
			diffs = append(diffs, types.PackageJsonDiff{
				PackageName:          pkgName,
				FilePath:             bag.getPackageJsonPath(pkgName),
				DependenciesToAdd:    []types.DependencyChange{},
				DependenciesToRemove: []types.DependencyChange{},
				Summary:              "Review and update dependencies after boundary refactoring",
			})
		}
	}

	return diffs
}

// getPackageJsonPath returns the path to a package's package.json.
func (bag *BeforeAfterGenerator) getPackageJsonPath(pkgName string) string {
	if bag.workspace != nil && bag.workspace.Packages != nil {
		if pkg, ok := bag.workspace.Packages[pkgName]; ok && pkg.Path != "" {
			return fmt.Sprintf("%s/package.json", pkg.Path)
		}
	}
	// Fallback: derive from package name
	dirName := extractShortName(pkgName)
	return fmt.Sprintf("packages/%s/package.json", dirName)
}

// ========================================
// Import Diff Generation (Task 6, AC #4)
// ========================================

// generateImportDiffs creates import statement change descriptions.
func (bag *BeforeAfterGenerator) generateImportDiffs(
	cycle *types.CircularDependencyInfo,
	strategy *types.FixStrategy,
) []types.ImportDiff {
	diffs := []types.ImportDiff{}

	// Use ImportTraces from Story 3.2 if available
	if len(cycle.ImportTraces) > 0 {
		for _, trace := range cycle.ImportTraces {
			// Find the import that needs to change based on strategy
			var toRemove []types.ImportChange
			var toAdd []types.ImportChange

			if strategy.Type == types.FixStrategyExtractModule && strategy.NewPackageName != "" {
				// Replace import from cycle package with import from shared package
				toRemove = []types.ImportChange{
					{
						Statement:     trace.Statement,
						FromPackage:   trace.ToPackage,
						ImportedNames: trace.Symbols,
					},
				}
				toAdd = []types.ImportChange{
					{
						Statement:     bag.generateNewImportStatement(trace, strategy.NewPackageName),
						FromPackage:   strategy.NewPackageName,
						ImportedNames: trace.Symbols,
					},
				}
			} else if strategy.Type == types.FixStrategyDependencyInject {
				// For DI, the import is removed entirely
				toRemove = []types.ImportChange{
					{
						Statement:     trace.Statement,
						FromPackage:   trace.ToPackage,
						ImportedNames: trace.Symbols,
					},
				}
				// Add interface import instead
				interfaceName := generateInterfaceNameForPkg(trace.ToPackage)
				toAdd = []types.ImportChange{
					{
						Statement:     fmt.Sprintf("import type { %s } from '%s';", interfaceName, trace.ToPackage),
						FromPackage:   trace.ToPackage,
						ImportedNames: []string{interfaceName},
					},
				}
			}

			if len(toRemove) > 0 || len(toAdd) > 0 {
				diffs = append(diffs, types.ImportDiff{
					FilePath:        trace.FilePath,
					PackageName:     trace.FromPackage,
					ImportsToRemove: toRemove,
					ImportsToAdd:    toAdd,
					LineNumber:      trace.LineNumber,
				})
			}
		}
	} else {
		// Estimate import diffs based on cycle structure
		for i := 0; i < len(cycle.Cycle)-1; i++ {
			fromPkg := cycle.Cycle[i]
			toPkg := cycle.Cycle[i+1]

			toRemove := []types.ImportChange{
				{
					Statement:   fmt.Sprintf("import { ... } from '%s'", toPkg),
					FromPackage: toPkg,
				},
			}
			toAdd := bag.generateEstimatedImportAdd(strategy, toPkg)

			diffs = append(diffs, types.ImportDiff{
				FilePath:        fmt.Sprintf("packages/%s/src/index.ts", extractShortName(fromPkg)),
				PackageName:     fromPkg,
				ImportsToRemove: toRemove,
				ImportsToAdd:    toAdd,
			})
		}
	}

	return diffs
}

// generateNewImportStatement creates a new import statement for the shared package.
func (bag *BeforeAfterGenerator) generateNewImportStatement(trace types.ImportTrace, newPkgName string) string {
	switch trace.ImportType {
	case types.ImportTypeESMNamed:
		if len(trace.Symbols) > 0 {
			return fmt.Sprintf("import { %s } from '%s';", strings.Join(trace.Symbols, ", "), newPkgName)
		}
		return fmt.Sprintf("import { ... } from '%s';", newPkgName)
	case types.ImportTypeESMDefault:
		return fmt.Sprintf("import defaultExport from '%s';", newPkgName)
	case types.ImportTypeESMNamespace:
		return fmt.Sprintf("import * as shared from '%s';", newPkgName)
	default:
		return fmt.Sprintf("import { ... } from '%s';", newPkgName)
	}
}

// generateEstimatedImportAdd creates estimated import additions based on strategy.
func (bag *BeforeAfterGenerator) generateEstimatedImportAdd(strategy *types.FixStrategy, toPkg string) []types.ImportChange {
	switch strategy.Type {
	case types.FixStrategyExtractModule:
		if strategy.NewPackageName != "" {
			return []types.ImportChange{
				{
					Statement:   fmt.Sprintf("import { ... } from '%s'", strategy.NewPackageName),
					FromPackage: strategy.NewPackageName,
				},
			}
		}
	case types.FixStrategyDependencyInject:
		interfaceName := generateInterfaceNameForPkg(toPkg)
		return []types.ImportChange{
			{
				Statement:     fmt.Sprintf("import type { %s } from '%s'", interfaceName, toPkg),
				FromPackage:   toPkg,
				ImportedNames: []string{interfaceName},
			},
		}
	}
	return []types.ImportChange{}
}

// generateInterfaceNameForPkg creates an interface name from package name.
func generateInterfaceNameForPkg(pkgName string) string {
	dirName := extractShortName(pkgName)
	if len(dirName) > 0 {
		return strings.ToUpper(dirName[:1]) + dirName[1:] + "Handler"
	}
	return "Handler"
}

// ========================================
// Plain Language Explanation (Task 7, AC #5)
// ========================================

// generateExplanation creates the human-readable explanation.
func (bag *BeforeAfterGenerator) generateExplanation(
	cycle *types.CircularDependencyInfo,
	strategy *types.FixStrategy,
) *types.FixExplanation {
	var summary, whyItWorks string
	var highLevelChanges []string
	confidence := 0.8 // Default confidence

	switch strategy.Type {
	case types.FixStrategyExtractModule:
		summary = fmt.Sprintf(
			"Create a new shared package '%s' to hold the common code that causes the circular dependency between %s.",
			strategy.NewPackageName,
			formatPackageList(strategy.TargetPackages),
		)
		whyItWorks = fmt.Sprintf(
			"The cycle exists because %s each import something from each other. "+
				"By moving the shared functionality to '%s', all packages can import from the "+
				"new package instead of from each other, breaking the cycle.",
			formatPackageList(strategy.TargetPackages),
			strategy.NewPackageName,
		)
		highLevelChanges = []string{
			fmt.Sprintf("Create new package: %s", strategy.NewPackageName),
			"Move shared types, functions, or constants to the new package",
			fmt.Sprintf("Update imports in: %s", strings.Join(strategy.TargetPackages, ", ")),
			"Add the new package as a dependency in affected package.json files",
		}
		confidence = 0.9

	case types.FixStrategyDependencyInject:
		summary = "Invert the problematic dependency by using dependency injection."
		whyItWorks = "Instead of having packages directly import each other (creating a cycle), " +
			"one package defines an interface that the other implements. The dependency is then " +
			"provided at runtime, breaking the compile-time circular dependency."
		highLevelChanges = []string{
			"Create an interface that abstracts the dependency",
			"Update the dependent package to accept the dependency as a parameter",
			"Wire up the dependency at the application's composition root",
			"Remove the direct circular import",
		}
		confidence = 0.75

	case types.FixStrategyBoundaryRefactor:
		summary = "Restructure package boundaries to eliminate overlapping responsibilities."
		whyItWorks = "The cycle often indicates that package responsibilities are not clearly defined. " +
			"By identifying code that belongs in the wrong package and moving it, the dependency " +
			"relationship becomes one-directional instead of circular."
		highLevelChanges = []string{
			"Analyze which code is causing the cycle",
			"Identify the correct package for each piece of functionality",
			"Move code to appropriate packages",
			"Update imports throughout the codebase",
		}
		confidence = 0.7
	}

	return &types.FixExplanation{
		Summary:          summary,
		WhyItWorks:       whyItWorks,
		HighLevelChanges: highLevelChanges,
		Confidence:       confidence,
	}
}

// ========================================
// Side Effect Warnings (Task 8, AC #6)
// ========================================

// generateWarnings identifies potential side effects.
func (bag *BeforeAfterGenerator) generateWarnings(
	cycle *types.CircularDependencyInfo,
	strategy *types.FixStrategy,
) []types.SideEffectWarning {
	warnings := []types.SideEffectWarning{}

	// Check for breaking changes based on strategy type
	switch strategy.Type {
	case types.FixStrategyExtractModule:
		// Warn about potential API changes
		warnings = append(warnings, types.SideEffectWarning{
			Severity:    types.WarningSeverityInfo,
			Title:       "New package requires installation",
			Description: fmt.Sprintf("After creating '%s', run your package manager's install command to link it.", strategy.NewPackageName),
		})

		if len(strategy.TargetPackages) > 3 {
			warnings = append(warnings, types.SideEffectWarning{
				Severity:         types.WarningSeverityWarning,
				Title:            "Multiple packages affected",
				Description:      "This fix affects many packages. Consider making changes incrementally and testing after each change.",
				AffectedPackages: strategy.TargetPackages,
			})
		}

	case types.FixStrategyDependencyInject:
		warnings = append(warnings, types.SideEffectWarning{
			Severity:    types.WarningSeverityWarning,
			Title:       "API signature changes",
			Description: "Functions or classes may require additional parameters for dependency injection. Update all call sites.",
		})

		warnings = append(warnings, types.SideEffectWarning{
			Severity:    types.WarningSeverityInfo,
			Title:       "Runtime wiring required",
			Description: "Dependencies must be wired up at the application entry point. Ensure proper initialization order.",
		})

	case types.FixStrategyBoundaryRefactor:
		warnings = append(warnings, types.SideEffectWarning{
			Severity:    types.WarningSeverityCritical,
			Title:       "Significant code restructuring",
			Description: "This strategy involves moving code between packages. Thoroughly test affected functionality.",
		})

		warnings = append(warnings, types.SideEffectWarning{
			Severity:    types.WarningSeverityWarning,
			Title:       "Potential export changes",
			Description: "Moved code may change which exports are available from each package. Update consumers accordingly.",
		})
	}

	// Check if cycle involves core packages (higher risk)
	corePatterns := []string{"core", "common", "shared", "utils", "lib"}
	for _, pkg := range strategy.TargetPackages {
		pkgLower := strings.ToLower(pkg)
		for _, pattern := range corePatterns {
			if strings.Contains(pkgLower, pattern) {
				warnings = append(warnings, types.SideEffectWarning{
					Severity:         types.WarningSeverityCritical,
					Title:            "Core package affected",
					Description:      fmt.Sprintf("Package '%s' appears to be a core/shared package. Changes may have widespread impact.", pkg),
					AffectedPackages: []string{pkg},
				})
				break
			}
		}
	}

	// Use ImpactAssessment if available
	if cycle.ImpactAssessment != nil && cycle.ImpactAssessment.RiskLevel == types.RiskLevelCritical {
		warnings = append(warnings, types.SideEffectWarning{
			Severity:    types.WarningSeverityCritical,
			Title:       "High-impact cycle",
			Description: fmt.Sprintf("This cycle affects %d%% of the monorepo. Proceed carefully.", int(cycle.ImpactAssessment.AffectedPercentage*100)),
		})
	}

	return warnings
}

// ========================================
// Helper Functions
// ========================================

// extractShortName extracts the short name from a package name.
// e.g., "@mono/shared" -> "shared"
func extractShortName(pkgName string) string {
	parts := strings.Split(pkgName, "/")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return pkgName
}
