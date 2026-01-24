// Package types defines Go types that match TypeScript definitions in @monoguard/types.
// This file contains before/after explanation types for Story 3.7.
package types

// ========================================
// Before/After Explanation Types (Story 3.7)
// ========================================

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

// NewBeforeAfterExplanation creates a new BeforeAfterExplanation with initialized slices.
func NewBeforeAfterExplanation() *BeforeAfterExplanation {
	return &BeforeAfterExplanation{
		PackageJsonDiffs: []PackageJsonDiff{},
		ImportDiffs:      []ImportDiff{},
		Warnings:         []SideEffectWarning{},
	}
}

// NewStateDiagram creates a new StateDiagram with initialized slices.
func NewStateDiagram() *StateDiagram {
	return &StateDiagram{
		Nodes:         []DiagramNode{},
		Edges:         []DiagramEdge{},
		CycleResolved: false,
	}
}
