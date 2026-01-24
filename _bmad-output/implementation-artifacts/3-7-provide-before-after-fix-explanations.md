# Story 3.7: Provide Before/After Fix Explanations

Status: done

## Story

As a **user**,
I want **to see clear before/after comparisons for each fix**,
So that **I understand exactly what will change and can build confidence in the fix**.

## Acceptance Criteria

1. **AC1: Current State Diagram Data**
   - Given a fix recommendation
   - When I request before/after explanation
   - Then I see current state visualization data:
     - Dependency graph snapshot with cycle highlighted
     - Nodes for all packages involved in the cycle
     - Edges showing the problematic dependency path
   - And data is structured for D3.js visualization (Epic 4)

2. **AC2: Proposed State Diagram Data**
   - Given a fix strategy (from Story 3.3)
   - When generating proposed state
   - Then I see visualization data showing:
     - Resolved dependency graph (cycle broken)
     - New packages added (if applicable, e.g., extracted shared module)
     - Modified edges showing new dependency relationships
   - And the cycle is no longer present in the proposed state

3. **AC3: Package.json Diff**
   - Given a fix strategy
   - When generating package.json changes
   - Then I see for each affected package:
     - File path to package.json
     - Dependencies to add (with version)
     - Dependencies to remove (if applicable)
     - Summary of changes
   - And diffs are in a structured format (not raw text)

4. **AC4: Import Statement Diff**
   - Given a fix strategy with import changes
   - When generating import diffs
   - Then I see for each affected file:
     - File path
     - Import statements to remove
     - Import statements to add
     - Line number hints (if available from ImportTraces)
   - And diffs reference actual package names involved

5. **AC5: Plain Language Explanation**
   - Given before/after data
   - When generating explanation
   - Then provide:
     - Summary of what the fix accomplishes (1-2 sentences)
     - Why this resolves the circular dependency
     - What code changes are required (high-level)
   - And explanation avoids technical jargon where possible
   - And explanation is understandable by non-expert developers

6. **AC6: Side Effect Warnings**
   - Given a fix that may have side effects
   - When generating warnings
   - Then include:
     - Breaking changes (API changes, removed exports)
     - Build/test impacts (packages that need rebuilding)
     - Runtime behavior changes (if detectable)
   - And warnings are clearly marked with severity (info/warning/critical)

7. **AC7: Integration with FixStrategy**
   - Given analysis results with fix strategies
   - When enriching FixStrategy
   - Then add `beforeAfterExplanation` field:
     ```go
     type FixStrategy struct {
         // ... existing fields ...
         BeforeAfterExplanation *BeforeAfterExplanation `json:"beforeAfterExplanation,omitempty"`
     }
     ```
   - And explanation generation is optional (can be requested separately)

8. **AC8: Performance**
   - Given a workspace with 100 packages and 5 cycles
   - When generating all before/after explanations
   - Then generation completes in < 300ms additional overhead
   - And memory usage increase is < 15MB

## Tasks / Subtasks

- [x] **Task 1: Define BeforeAfterExplanation Types** (AC: #1, #2, #3, #4, #5, #6)
  - [x] 1.1 Create `pkg/types/before_after_explanation.go`:
    ```go
    package types

    // BeforeAfterExplanation provides visual comparison data for fix strategies.
    // Matches @monoguard/types BeforeAfterExplanation interface.
    type BeforeAfterExplanation struct {
        // CurrentState represents the dependency graph before the fix
        CurrentState *StateDiagram `json:"currentState"`

        // ProposedState represents the dependency graph after the fix
        ProposedState *StateDiagram `json:"proposedState"`

        // PackageJsonDiffs shows changes required to package.json files
        PackageJsonDiffs []PackageJsonDiff `json:"packageJsonDiffs"`

        // ImportDiffs shows changes required to import statements
        ImportDiffs []ImportDiff `json:"importDiffs"`

        // Explanation provides human-readable summary
        Explanation *FixExplanation `json:"explanation"`

        // Warnings about potential side effects
        Warnings []SideEffectWarning `json:"warnings"`
    }

    // StateDiagram contains D3.js-compatible visualization data.
    type StateDiagram struct {
        // Nodes are the packages in the diagram
        Nodes []DiagramNode `json:"nodes"`

        // Edges are the dependency relationships
        Edges []DiagramEdge `json:"edges"`

        // HighlightedPath shows the cycle path (only in CurrentState)
        HighlightedPath []string `json:"highlightedPath,omitempty"`

        // CycleResolved indicates if this state has no cycle
        CycleResolved bool `json:"cycleResolved"`
    }

    // DiagramNode represents a package in the visualization.
    type DiagramNode struct {
        // ID is the package name (used for edge references)
        ID string `json:"id"`

        // Label is the display name
        Label string `json:"label"`

        // IsInCycle indicates if this package is part of the cycle
        IsInCycle bool `json:"isInCycle"`

        // IsNew indicates if this package is newly created by the fix
        IsNew bool `json:"isNew"`

        // NodeType categorizes the package (cycle, affected, new, unchanged)
        NodeType DiagramNodeType `json:"nodeType"`
    }

    // DiagramNodeType categorizes nodes for visualization styling.
    type DiagramNodeType string

    const (
        NodeTypeCycle     DiagramNodeType = "cycle"     // Part of the cycle
        NodeTypeAffected  DiagramNodeType = "affected"  // Indirectly affected
        NodeTypeNew       DiagramNodeType = "new"       // Newly created package
        NodeTypeUnchanged DiagramNodeType = "unchanged" // Not affected by fix
    )

    // DiagramEdge represents a dependency relationship.
    type DiagramEdge struct {
        // From is the dependent package
        From string `json:"from"`

        // To is the dependency
        To string `json:"to"`

        // IsInCycle indicates if this edge is part of the cycle
        IsInCycle bool `json:"isInCycle"`

        // IsRemoved indicates if this edge will be removed by the fix
        IsRemoved bool `json:"isRemoved"`

        // IsNew indicates if this edge is added by the fix
        IsNew bool `json:"isNew"`

        // EdgeType categorizes the edge for visualization styling
        EdgeType DiagramEdgeType `json:"edgeType"`
    }

    // DiagramEdgeType categorizes edges for visualization styling.
    type DiagramEdgeType string

    const (
        EdgeTypeCycle     DiagramEdgeType = "cycle"     // Part of the cycle (red)
        EdgeTypeRemoved   DiagramEdgeType = "removed"   // To be removed (strikethrough)
        EdgeTypeNew       DiagramEdgeType = "new"       // New dependency (green)
        EdgeTypeUnchanged DiagramEdgeType = "unchanged" // Not affected
    )

    // PackageJsonDiff describes changes to a package.json file.
    type PackageJsonDiff struct {
        // PackageName is the package being modified
        PackageName string `json:"packageName"`

        // FilePath is the relative path to package.json
        FilePath string `json:"filePath"`

        // DependenciesToAdd lists dependencies to add
        DependenciesToAdd []DependencyChange `json:"dependenciesToAdd"`

        // DependenciesToRemove lists dependencies to remove
        DependenciesToRemove []DependencyChange `json:"dependenciesToRemove"`

        // Summary is a human-readable change description
        Summary string `json:"summary"`
    }

    // DependencyChange describes a dependency addition or removal.
    type DependencyChange struct {
        // Name is the dependency package name
        Name string `json:"name"`

        // Version is the version specifier (e.g., "workspace:*", "^1.0.0")
        Version string `json:"version,omitempty"`

        // DependencyType indicates dependencies vs devDependencies
        DependencyType string `json:"dependencyType"`
    }

    // ImportDiff describes changes to import statements in a file.
    type ImportDiff struct {
        // FilePath is the file containing imports
        FilePath string `json:"filePath"`

        // PackageName is the package containing this file
        PackageName string `json:"packageName"`

        // ImportsToRemove lists import statements to remove
        ImportsToRemove []ImportChange `json:"importsToRemove"`

        // ImportsToAdd lists import statements to add
        ImportsToAdd []ImportChange `json:"importsToAdd"`

        // LineNumber hints at location (if available from ImportTraces)
        LineNumber int `json:"lineNumber,omitempty"`
    }

    // ImportChange describes an import statement change.
    type ImportChange struct {
        // Statement is the full import statement
        Statement string `json:"statement"`

        // FromPackage is the package being imported from
        FromPackage string `json:"fromPackage"`

        // ImportedNames lists what is being imported
        ImportedNames []string `json:"importedNames,omitempty"`
    }

    // FixExplanation provides human-readable explanation of the fix.
    type FixExplanation struct {
        // Summary is a 1-2 sentence overview
        Summary string `json:"summary"`

        // WhyItWorks explains how this resolves the cycle
        WhyItWorks string `json:"whyItWorks"`

        // HighLevelChanges describes what code changes are required
        HighLevelChanges []string `json:"highLevelChanges"`

        // Confidence indicates how certain we are about the fix (0.0-1.0)
        Confidence float64 `json:"confidence"`
    }

    // SideEffectWarning describes a potential side effect of the fix.
    type SideEffectWarning struct {
        // Severity indicates the importance (info, warning, critical)
        Severity WarningSeverity `json:"severity"`

        // Title is a short description
        Title string `json:"title"`

        // Description provides details
        Description string `json:"description"`

        // AffectedPackages lists packages that may be affected
        AffectedPackages []string `json:"affectedPackages,omitempty"`
    }

    // WarningSeverity indicates the importance of a warning.
    type WarningSeverity string

    const (
        WarningSeverityInfo     WarningSeverity = "info"
        WarningSeverityWarning  WarningSeverity = "warning"
        WarningSeverityCritical WarningSeverity = "critical"
    )
    ```
  - [x] 1.2 Add JSON serialization tests in `pkg/types/before_after_explanation_test.go`
  - [x] 1.3 Ensure all JSON tags use camelCase

- [x] **Task 2: Create BeforeAfterGenerator** (AC: #1, #2, #7)
  - [x] 2.1 Create `pkg/analyzer/before_after_generator.go`:
    ```go
    package analyzer

    import "github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"

    // BeforeAfterGenerator creates before/after comparisons for fix strategies.
    type BeforeAfterGenerator struct {
        graph     *types.DependencyGraph
        workspace *types.WorkspaceData
    }

    // NewBeforeAfterGenerator creates a new generator.
    func NewBeforeAfterGenerator(graph *types.DependencyGraph, workspace *types.WorkspaceData) *BeforeAfterGenerator

    // Generate creates the before/after explanation for a strategy.
    func (bag *BeforeAfterGenerator) Generate(
        cycle *types.CircularDependencyInfo,
        strategy *types.FixStrategy,
    ) *types.BeforeAfterExplanation

    // generateCurrentState creates the "before" diagram.
    func (bag *BeforeAfterGenerator) generateCurrentState(cycle *types.CircularDependencyInfo) *types.StateDiagram

    // generateProposedState creates the "after" diagram.
    func (bag *BeforeAfterGenerator) generateProposedState(
        cycle *types.CircularDependencyInfo,
        strategy *types.FixStrategy,
    ) *types.StateDiagram

    // generatePackageJsonDiffs creates package.json change descriptions.
    func (bag *BeforeAfterGenerator) generatePackageJsonDiffs(strategy *types.FixStrategy) []types.PackageJsonDiff

    // generateImportDiffs creates import statement change descriptions.
    func (bag *BeforeAfterGenerator) generateImportDiffs(
        cycle *types.CircularDependencyInfo,
        strategy *types.FixStrategy,
    ) []types.ImportDiff

    // generateExplanation creates the human-readable explanation.
    func (bag *BeforeAfterGenerator) generateExplanation(
        cycle *types.CircularDependencyInfo,
        strategy *types.FixStrategy,
    ) *types.FixExplanation

    // generateWarnings identifies potential side effects.
    func (bag *BeforeAfterGenerator) generateWarnings(
        cycle *types.CircularDependencyInfo,
        strategy *types.FixStrategy,
    ) []types.SideEffectWarning
    ```
  - [x] 2.2 Implement generator dispatch logic
  - [x] 2.3 Create comprehensive tests in `pkg/analyzer/before_after_generator_test.go`

- [x] **Task 3: Implement Current State Diagram Generation** (AC: #1)
  - [x] 3.1 Implement `generateCurrentState`:
    ```go
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
    ```
  - [x] 3.2 Add tests for current state generation

- [x] **Task 4: Implement Proposed State Diagram Generation** (AC: #2)
  - [x] 4.1 Implement `generateProposedState` for Extract Module strategy:
    ```go
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
    ```
  - [x] 4.2 Implement `generateProposedState` for Dependency Injection strategy
  - [x] 4.3 Implement `generateProposedState` for Boundary Refactoring strategy
  - [x] 4.4 Add tests for proposed state generation

- [x] **Task 5: Implement Package.json Diff Generation** (AC: #3)
  - [x] 5.1 Implement `generatePackageJsonDiffs`:
    ```go
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
            // But may need to add types package
            diffs = append(diffs, types.PackageJsonDiff{
                PackageName:          strategy.TargetPackages[0],
                FilePath:             bag.getPackageJsonPath(strategy.TargetPackages[0]),
                DependenciesToAdd:    []types.DependencyChange{},
                DependenciesToRemove: []types.DependencyChange{},
                Summary:              "No package.json changes required for dependency injection",
            })

        case types.FixStrategyBoundaryRefactor:
            // Boundary refactoring may involve moving code between packages
            for _, pkgName := range strategy.TargetPackages {
                diffs = append(diffs, types.PackageJsonDiff{
                    PackageName: pkgName,
                    FilePath:    bag.getPackageJsonPath(pkgName),
                    Summary:     "Review and update dependencies after boundary refactoring",
                })
            }
        }

        return diffs
    }
    ```
  - [x] 5.2 Add tests for package.json diff generation

- [x] **Task 6: Implement Import Diff Generation** (AC: #4)
  - [x] 6.1 Implement `generateImportDiffs`:
    ```go
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
                            Statement:   trace.Statement,
                            FromPackage: trace.ToPackage,
                        },
                    }
                    toAdd = []types.ImportChange{
                        {
                            Statement:   bag.generateNewImportStatement(trace, strategy.NewPackageName),
                            FromPackage: strategy.NewPackageName,
                        },
                    }
                }

                diffs = append(diffs, types.ImportDiff{
                    FilePath:        trace.FilePath,
                    PackageName:     trace.FromPackage,
                    ImportsToRemove: toRemove,
                    ImportsToAdd:    toAdd,
                    LineNumber:      trace.LineNumber,
                })
            }
        } else {
            // Estimate import diffs based on cycle structure
            for i := 0; i < len(cycle.Cycle)-1; i++ {
                fromPkg := cycle.Cycle[i]
                toPkg := cycle.Cycle[i+1]

                diffs = append(diffs, types.ImportDiff{
                    FilePath:    fmt.Sprintf("packages/%s/src/index.ts", extractPkgDirName(fromPkg)),
                    PackageName: fromPkg,
                    ImportsToRemove: []types.ImportChange{
                        {
                            Statement:   fmt.Sprintf("import { ... } from '%s'", toPkg),
                            FromPackage: toPkg,
                        },
                    },
                    ImportsToAdd: bag.generateEstimatedImportAdd(strategy, toPkg),
                })
            }
        }

        return diffs
    }
    ```
  - [x] 6.2 Add tests for import diff generation

- [x] **Task 7: Implement Plain Language Explanation** (AC: #5)
  - [x] 7.1 Implement `generateExplanation`:
    ```go
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
    ```
  - [x] 7.2 Add tests for explanation generation

- [x] **Task 8: Implement Side Effect Warnings** (AC: #6)
  - [x] 8.1 Implement `generateWarnings`:
    ```go
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
    ```
  - [x] 8.2 Add tests for warning generation

- [x] **Task 9: Implement Main Generate Function** (AC: #1, #2, #3, #4, #5, #6)
  - [x] 9.1 Implement `Generate`:
    ```go
    func (bag *BeforeAfterGenerator) Generate(
        cycle *types.CircularDependencyInfo,
        strategy *types.FixStrategy,
    ) *types.BeforeAfterExplanation {
        return &types.BeforeAfterExplanation{
            CurrentState:     bag.generateCurrentState(cycle),
            ProposedState:    bag.generateProposedState(cycle, strategy),
            PackageJsonDiffs: bag.generatePackageJsonDiffs(strategy),
            ImportDiffs:      bag.generateImportDiffs(cycle, strategy),
            Explanation:      bag.generateExplanation(cycle, strategy),
            Warnings:         bag.generateWarnings(cycle, strategy),
        }
    }
    ```
  - [x] 9.2 Add comprehensive integration tests

- [x] **Task 10: Integrate with FixStrategy** (AC: #7)
  - [x] 10.1 Update `pkg/types/fix_strategy.go`:
    ```go
    type FixStrategy struct {
        Type                   FixStrategyType         `json:"type"`
        Name                   string                  `json:"name"`
        Description            string                  `json:"description"`
        Suitability            int                     `json:"suitability"`
        Effort                 EffortLevel             `json:"effort"`
        Pros                   []string                `json:"pros"`
        Cons                   []string                `json:"cons"`
        TargetPackages         []string                `json:"targetPackages"`
        NewPackageName         string                  `json:"newPackageName,omitempty"`
        Guide                  *FixGuide               `json:"guide,omitempty"`               // Story 3.4
        Complexity             *RefactoringComplexity  `json:"complexity,omitempty"`          // Story 3.5
        BeforeAfterExplanation *BeforeAfterExplanation `json:"beforeAfterExplanation,omitempty"` // NEW Story 3.7
    }
    ```
  - [x] 10.2 Verify existing tests still pass

- [x] **Task 11: Wire to Analyzer Pipeline** (AC: all)
  - [x] 11.1 Update `pkg/analyzer/analyzer.go`:
    ```go
    func (a *Analyzer) AnalyzeWithSources(...) (*types.AnalysisResult, error) {
        // ... existing analysis ...

        // Generate fix strategies (Story 3.3)
        fixGenerator := NewFixStrategyGenerator(graph, workspace)
        guideGenerator := NewFixGuideGenerator(workspace)          // Story 3.4
        beforeAfterGenerator := NewBeforeAfterGenerator(graph, workspace) // NEW Story 3.7

        for _, cycle := range cycles {
            strategies := fixGenerator.Generate(cycle)

            for i := range strategies {
                // Story 3.4: Generate guides
                strategies[i].Guide = guideGenerator.Generate(cycle, &strategies[i])

                // Story 3.7: Generate before/after explanations
                strategies[i].BeforeAfterExplanation = beforeAfterGenerator.Generate(cycle, &strategies[i])
            }

            cycle.FixStrategies = strategies
        }

        return result, nil
    }
    ```
  - [x] 11.2 Update analyzer tests

- [x] **Task 12: Update TypeScript Types** (AC: #7)
  - [x] 12.1 Update `packages/types/src/analysis/results.ts`:
    ```typescript
    export interface BeforeAfterExplanation {
      /** Dependency graph before the fix */
      currentState: StateDiagram;
      /** Dependency graph after the fix */
      proposedState: StateDiagram;
      /** Changes required to package.json files */
      packageJsonDiffs: PackageJsonDiff[];
      /** Changes required to import statements */
      importDiffs: ImportDiff[];
      /** Human-readable explanation */
      explanation: FixExplanation;
      /** Potential side effects */
      warnings: SideEffectWarning[];
    }

    export interface StateDiagram {
      /** Packages in the diagram */
      nodes: DiagramNode[];
      /** Dependency relationships */
      edges: DiagramEdge[];
      /** The cycle path (only in currentState) */
      highlightedPath?: string[];
      /** Whether this state has no cycle */
      cycleResolved: boolean;
    }

    export interface DiagramNode {
      /** Package name (used for edge references) */
      id: string;
      /** Display name */
      label: string;
      /** Whether this package is part of the cycle */
      isInCycle: boolean;
      /** Whether this package is newly created by the fix */
      isNew: boolean;
      /** Node category for visualization styling */
      nodeType: DiagramNodeType;
    }

    export type DiagramNodeType = 'cycle' | 'affected' | 'new' | 'unchanged';

    export interface DiagramEdge {
      /** Dependent package */
      from: string;
      /** Dependency */
      to: string;
      /** Whether this edge is part of the cycle */
      isInCycle: boolean;
      /** Whether this edge will be removed by the fix */
      isRemoved: boolean;
      /** Whether this edge is added by the fix */
      isNew: boolean;
      /** Edge category for visualization styling */
      edgeType: DiagramEdgeType;
    }

    export type DiagramEdgeType = 'cycle' | 'removed' | 'new' | 'unchanged';

    export interface PackageJsonDiff {
      /** Package being modified */
      packageName: string;
      /** Relative path to package.json */
      filePath: string;
      /** Dependencies to add */
      dependenciesToAdd: DependencyChange[];
      /** Dependencies to remove */
      dependenciesToRemove: DependencyChange[];
      /** Human-readable change description */
      summary: string;
    }

    export interface DependencyChange {
      /** Dependency package name */
      name: string;
      /** Version specifier */
      version?: string;
      /** dependencies vs devDependencies */
      dependencyType: string;
    }

    export interface ImportDiff {
      /** File containing imports */
      filePath: string;
      /** Package containing this file */
      packageName: string;
      /** Import statements to remove */
      importsToRemove: ImportChange[];
      /** Import statements to add */
      importsToAdd: ImportChange[];
      /** Location hint (if available) */
      lineNumber?: number;
    }

    export interface ImportChange {
      /** Full import statement */
      statement: string;
      /** Package being imported from */
      fromPackage: string;
      /** What is being imported */
      importedNames?: string[];
    }

    export interface FixExplanation {
      /** 1-2 sentence overview */
      summary: string;
      /** How this resolves the cycle */
      whyItWorks: string;
      /** What code changes are required */
      highLevelChanges: string[];
      /** Confidence in the fix (0.0-1.0) */
      confidence: number;
    }

    export interface SideEffectWarning {
      /** Importance level */
      severity: WarningSeverity;
      /** Short description */
      title: string;
      /** Detailed description */
      description: string;
      /** Packages that may be affected */
      affectedPackages?: string[];
    }

    export type WarningSeverity = 'info' | 'warning' | 'critical';

    export interface FixStrategy {
      // ... existing fields ...
      /** Before/after comparison (Story 3.7) */
      beforeAfterExplanation?: BeforeAfterExplanation;
    }
    ```
  - [x] 12.2 Run `pnpm nx build types` to verify
  - [x] 12.3 Add type tests for BeforeAfterExplanation

- [x] **Task 13: Performance Testing** (AC: #8)
  - [x] 13.1 Create `pkg/analyzer/before_after_generator_benchmark_test.go`:
    ```go
    func BenchmarkBeforeAfterGeneration(b *testing.B) {
        graph := generateGraph(100)
        workspace := generateWorkspace(100)
        cycles := generateCyclesWithStrategies(5)
        generator := NewBeforeAfterGenerator(graph, workspace)

        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            for _, cycle := range cycles {
                for _, strategy := range cycle.FixStrategies {
                    generator.Generate(cycle, &strategy)
                }
            }
        }
    }
    ```
  - [x] 13.2 Verify < 300ms for 100 packages with 5 cycles
  - [x] 13.3 Document actual performance in completion notes

- [x] **Task 14: Integration Verification** (AC: all)
  - [x] 14.1 Run all tests: `cd packages/analysis-engine && make test`
  - [x] 14.2 Build WASM: `pnpm nx build @monoguard/analysis-engine`
  - [x] 14.3 Run affected CI checks: `pnpm nx affected --target=lint,test,type-check --base=main`
  - [x] 14.4 Verify JSON output includes beforeAfterExplanation field

## Dev Notes

### Architecture Patterns & Constraints

**From Architecture Document:**
- **Location:** Before/after generator in `pkg/analyzer/`
- **Pattern:** Generator pattern with graph + workspace + cycle + strategy input
- **Integration:** Enriches FixStrategy with BeforeAfterExplanation
- **Dependency:** Uses data from Story 3.2 (ImportTraces), 3.3 (FixStrategy), 3.6 (ImpactAssessment)

**Critical Constraints:**
- **camelCase JSON:** All struct tags MUST use camelCase
- **Optional Fields:** BeforeAfterExplanation is `omitempty` - backward compatible
- **D3.js Compatible:** StateDiagram must be usable by D3.js force-directed graph (Epic 4)
- **Plain Language:** Explanations should be understandable by non-expert developers

### Critical Don't-Miss Rules

**From project-context.md:**

1. **JSON Naming Convention:**
   ```go
   // ✅ CORRECT: camelCase JSON tags
   type BeforeAfterExplanation struct {
       CurrentState  *StateDiagram     `json:"currentState"`
       ProposedState *StateDiagram     `json:"proposedState"`
       PackageJsonDiffs []PackageJsonDiff `json:"packageJsonDiffs"`
   }

   // ❌ WRONG: snake_case JSON tags
   type BeforeAfterExplanation struct {
       CurrentState *StateDiagram `json:"current_state"` // WRONG!
   }
   ```

2. **Enum Constants:**
   ```go
   // ✅ CORRECT: kebab-case string values for TypeScript compatibility
   const (
       NodeTypeCycle    DiagramNodeType = "cycle"
       EdgeTypeNew      DiagramEdgeType = "new"
       WarningSeverityCritical WarningSeverity = "critical"
   )

   // ❌ WRONG: UPPER_CASE
   const (
       NodeTypeCycle DiagramNodeType = "CYCLE" // WRONG!
   )
   ```

3. **Slice Initialization:**
   ```go
   // ✅ CORRECT: Initialize as empty slice
   nodes := []types.DiagramNode{}
   warnings := []types.SideEffectWarning{}

   // ❌ WRONG: Nil slice (serializes as null)
   var nodes []types.DiagramNode // nil
   ```

4. **Pointer vs Value for Optional Structs:**
   ```go
   // ✅ CORRECT: Use pointer for optional nested structs
   type FixStrategy struct {
       BeforeAfterExplanation *BeforeAfterExplanation `json:"beforeAfterExplanation,omitempty"`
   }

   // ❌ WRONG: Value type with omitempty doesn't work for structs
   type FixStrategy struct {
       BeforeAfterExplanation BeforeAfterExplanation `json:"beforeAfterExplanation,omitempty"` // Won't omit!
   }
   ```

### Project Structure Notes

**Target Directory Structure:**
```
packages/analysis-engine/
├── pkg/
│   ├── analyzer/
│   │   ├── analyzer.go                          # UPDATE: Add before/after generation
│   │   ├── before_after_generator.go            # NEW: Generator implementation
│   │   ├── before_after_generator_test.go       # NEW: Generator tests
│   │   └── before_after_generator_benchmark_test.go # NEW: Performance
│   └── types/
│       ├── fix_strategy.go                      # UPDATE: Add BeforeAfterExplanation field
│       ├── before_after_explanation.go          # NEW: All types
│       └── before_after_explanation_test.go     # NEW: Type tests
└── ...

packages/types/src/analysis/
├── results.ts                                   # UPDATE: Add TS types
└── ...
```

### Previous Story Intelligence

**From Story 3.4 (done):**
- FixGuide has codeBefore/codeAfter for code snippets - **Key Insight:** Similar pattern for import diffs
- CodeSnippet has language field - **Key Insight:** Reuse for import statements
- Verification steps pattern - **Key Insight:** Warnings serve similar user-reassurance purpose
- Performance: ~0.12ms for guide generation - target similar performance

**From Story 3.6 (ready-for-dev):**
- ImpactAssessment has RiskLevel - **Key Insight:** Use for warning severity decisions
- RippleEffect has layers for visualization - **Key Insight:** Similar D3.js-ready structure
- IndirectDependent has dependencyPath - **Key Insight:** Similar to DiagramEdge chain

**From Story 3.3 (done):**
- FixStrategy has type, targetPackages, newPackageName - **Key Insight:** Use for generating diffs
- FixStrategyType enum - **Key Insight:** Drive different generation logic per type
- Suitability scoring - **Key Insight:** Confidence field serves similar purpose

**From Story 3.2 (done):**
- ImportTrace has filePath, lineNumber, statement - **Key Insight:** Use for accurate import diffs
- ImportType enum - **Key Insight:** Could inform warning generation

### D3.js Visualization Compatibility

**StateDiagram is designed for D3.js force-directed graph:**
- Nodes have `id` and `label` for D3.js data binding
- Edges have `from` and `to` matching node IDs
- Boolean flags (`isInCycle`, `isNew`, `isRemoved`) for conditional styling
- `nodeType` and `edgeType` enums for CSS class mapping

**Example D3.js usage:**
```javascript
// Nodes
const nodes = stateDiagram.nodes.map(n => ({
  id: n.id,
  label: n.label,
  class: n.nodeType // 'cycle', 'affected', 'new', 'unchanged'
}));

// Edges
const edges = stateDiagram.edges.map(e => ({
  source: e.from,
  target: e.to,
  class: e.edgeType // 'cycle', 'removed', 'new', 'unchanged'
}));
```

### Input/Output Format

**Input (CircularDependencyInfo with FixStrategy):**
```json
{
  "cycle": ["@mono/ui", "@mono/api", "@mono/core", "@mono/ui"],
  "fixStrategies": [
    {
      "type": "extract-module",
      "name": "Extract Shared Module",
      "targetPackages": ["@mono/ui", "@mono/api", "@mono/core"],
      "newPackageName": "@mono/shared"
    }
  ]
}
```

**Output (FixStrategy with BeforeAfterExplanation):**
```json
{
  "type": "extract-module",
  "name": "Extract Shared Module",
  "targetPackages": ["@mono/ui", "@mono/api", "@mono/core"],
  "newPackageName": "@mono/shared",
  "beforeAfterExplanation": {
    "currentState": {
      "nodes": [
        {"id": "@mono/ui", "label": "ui", "isInCycle": true, "nodeType": "cycle"},
        {"id": "@mono/api", "label": "api", "isInCycle": true, "nodeType": "cycle"},
        {"id": "@mono/core", "label": "core", "isInCycle": true, "nodeType": "cycle"}
      ],
      "edges": [
        {"from": "@mono/ui", "to": "@mono/api", "isInCycle": true, "edgeType": "cycle"},
        {"from": "@mono/api", "to": "@mono/core", "isInCycle": true, "edgeType": "cycle"},
        {"from": "@mono/core", "to": "@mono/ui", "isInCycle": true, "edgeType": "cycle"}
      ],
      "highlightedPath": ["@mono/ui", "@mono/api", "@mono/core", "@mono/ui"],
      "cycleResolved": false
    },
    "proposedState": {
      "nodes": [
        {"id": "@mono/ui", "label": "ui", "isInCycle": false, "nodeType": "affected"},
        {"id": "@mono/api", "label": "api", "isInCycle": false, "nodeType": "affected"},
        {"id": "@mono/core", "label": "core", "isInCycle": false, "nodeType": "affected"},
        {"id": "@mono/shared", "label": "shared", "isNew": true, "nodeType": "new"}
      ],
      "edges": [
        {"from": "@mono/ui", "to": "@mono/shared", "isNew": true, "edgeType": "new"},
        {"from": "@mono/api", "to": "@mono/shared", "isNew": true, "edgeType": "new"},
        {"from": "@mono/core", "to": "@mono/shared", "isNew": true, "edgeType": "new"}
      ],
      "cycleResolved": true
    },
    "packageJsonDiffs": [
      {
        "packageName": "@mono/ui",
        "filePath": "packages/ui/package.json",
        "dependenciesToAdd": [
          {"name": "@mono/shared", "version": "workspace:*", "dependencyType": "dependencies"}
        ],
        "summary": "Add dependency on @mono/shared"
      }
    ],
    "importDiffs": [
      {
        "filePath": "packages/ui/src/client.ts",
        "packageName": "@mono/ui",
        "importsToRemove": [
          {"statement": "import { helper } from '@mono/api'", "fromPackage": "@mono/api"}
        ],
        "importsToAdd": [
          {"statement": "import { helper } from '@mono/shared'", "fromPackage": "@mono/shared"}
        ],
        "lineNumber": 5
      }
    ],
    "explanation": {
      "summary": "Create a new shared package '@mono/shared' to hold the common code that causes the circular dependency between @mono/ui, @mono/api, and @mono/core.",
      "whyItWorks": "The cycle exists because these packages each import something from each other. By moving the shared functionality to '@mono/shared', all packages can import from the new package instead of from each other, breaking the cycle.",
      "highLevelChanges": [
        "Create new package: @mono/shared",
        "Move shared types, functions, or constants to the new package",
        "Update imports in: @mono/ui, @mono/api, @mono/core",
        "Add the new package as a dependency in affected package.json files"
      ],
      "confidence": 0.9
    },
    "warnings": [
      {
        "severity": "info",
        "title": "New package requires installation",
        "description": "After creating '@mono/shared', run your package manager's install command to link it."
      }
    ]
  }
}
```

### Test Scenarios

| Scenario | Strategy | Expected Nodes | Expected Edges | Warnings |
|----------|----------|----------------|----------------|----------|
| 2-pkg extract | extract-module | 3 (2 affected + 1 new) | 2 new edges | info: install |
| 3-pkg extract | extract-module | 4 (3 affected + 1 new) | 3 new edges | info: install |
| DI simple | dependency-injection | 2 affected | 1 removed, 1 new | warning: API change |
| Boundary refactor | boundary-refactoring | 2+ affected | varies | critical: restructuring |
| Core package | any | varies | varies | critical: core package |

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story 3.7]
- [Source: _bmad-output/planning-artifacts/prd.md#FR14]
- [Source: _bmad-output/project-context.md#Language-Specific Rules]
- [Source: _bmad-output/implementation-artifacts/3-4-create-step-by-step-fix-guides.md]
- [Source: _bmad-output/implementation-artifacts/3-6-generate-impact-assessment.md]
- [D3.js Force-Directed Graph](https://d3js.org/d3-force)

## Dev Agent Record

### Agent Model Used

Claude Opus 4.5 (claude-opus-4-5-20251101)

### Debug Log References

N/A

### Completion Notes List

1. **All ACs Implemented:**
   - AC1: Current state diagram generation with cycle highlighting
   - AC2: Proposed state diagram generation for all 3 strategy types
   - AC3: Package.json diff generation with structured format
   - AC4: Import diff generation using ImportTraces (when available) with line numbers
   - AC5: Plain language explanations with confidence scores
   - AC6: Side effect warnings with severity levels (info/warning/critical)
   - AC7: BeforeAfterExplanation integrated into FixStrategy
   - AC8: Performance verified at ~0.07ms (well under 300ms requirement)

2. **Key Implementation Details:**
   - BeforeAfterGenerator follows the same pattern as FixGuideGenerator and ComplexityCalculator
   - Impact assessment is now calculated before fix strategies to enable warnings to use ImpactAssessment data
   - D3.js-compatible StateDiagram with nodes (id, label, nodeType) and edges (from, to, edgeType)
   - TypeScript types added to @monoguard/types for full cross-language compatibility

3. **Performance Results:**
   - Standard benchmark (5 cycles × 3 strategies): ~70,887 ns/op (~0.07ms)
   - Large cycles benchmark (10 packages per cycle): ~163,298 ns/op (~0.16ms)
   - Both well under the 300ms AC#8 requirement

### File List

**New Files:**
- `packages/analysis-engine/pkg/types/before_after_explanation.go` - Type definitions
- `packages/analysis-engine/pkg/types/before_after_explanation_test.go` - Type tests
- `packages/analysis-engine/pkg/analyzer/before_after_generator.go` - Generator implementation
- `packages/analysis-engine/pkg/analyzer/before_after_generator_test.go` - Generator tests
- `packages/analysis-engine/pkg/analyzer/before_after_generator_benchmark_test.go` - Performance benchmarks

**Modified Files:**
- `packages/analysis-engine/pkg/types/fix_strategy.go` - Added BeforeAfterExplanation field
- `packages/analysis-engine/pkg/analyzer/analyzer.go` - Wired generator into pipeline
- `packages/types/src/analysis/results.ts` - Added TypeScript types for Story 3.7

### Code Review (AI)

**Reviewer:** Amelia (Developer Agent)
**Date:** 2026-01-24
**Outcome:** Approved with fixes applied

**Issues Found and Fixed:**
1. **[FIXED]** Boundary Refactoring import diff was returning empty - added handling for `FixStrategyBoundaryRefactor` in `generateEstimatedImportAdd`
2. **[FIXED]** Added empty `TargetPackages` guard in `generatePackageJsonDiffs`
3. **[FIXED]** Improved DI and Boundary proposed state edge generation to handle multi-package cycles (now shows full chain instead of just first edge)
4. **[FIXED]** Added RiskLevel and Severity re-export in TypeScript types
5. **[FIXED]** Updated test for DI proposed state to expect 2 edges (1 removed + 1 unchanged)

**Files Modified by Code Review:**
- `packages/analysis-engine/pkg/analyzer/before_after_generator.go`
- `packages/analysis-engine/pkg/analyzer/before_after_generator_test.go`
- `packages/types/src/analysis/results.ts`
