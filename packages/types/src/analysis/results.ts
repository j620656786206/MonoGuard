import type { DependencyGraph, DependencyType, WorkspaceType } from './graph'

/**
 * AnalysisResult - Complete analysis output
 *
 * Matches Go: pkg/types/analysis_result.go
 * All date fields use ISO 8601 format (e.g., "2026-01-15T10:30:00Z")
 */
export interface AnalysisResult {
  /** Architecture health score (0-100) */
  healthScore: number
  /** Total packages analyzed (matches Go: packages) */
  packages: number
  /** Number of excluded packages (Story 2.6) */
  excludedPackages?: number
  /** Detected circular dependencies (Story 2.3) */
  circularDependencies?: CircularDependencyInfo[]
  /** Detected version conflicts (Story 2.4) */
  versionConflicts?: VersionConflictInfo[]
  /** Health score breakdown (Story 2.5) */
  healthScoreDetails?: HealthScoreDetails
  /** Dependency graph data */
  graph?: DependencyGraph
  /** Analysis metadata */
  metadata?: AnalysisMetadata
  /** ISO 8601 timestamp */
  createdAt?: string
}

/**
 * CircularDependencyInfo - Enhanced circular dependency with fix suggestions
 *
 * Matches Go: pkg/types/circular.go
 */
export interface CircularDependencyInfo {
  /** Packages involved in the cycle (in order, ends with first package) */
  cycle: string[]
  /** Type of circular dependency */
  type: 'direct' | 'indirect'
  /** Severity level */
  severity: 'critical' | 'warning' | 'info'
  /** Number of unique packages in the cycle */
  depth: number
  /** Impact description */
  impact: string
  /** Legacy: Basic refactoring complexity (1-10) */
  complexity: number
  /** Detailed refactoring complexity (Story 3.5) */
  refactoringComplexity?: RefactoringComplexity
  /** Root cause analysis (Story 3.1) */
  rootCause?: RootCauseAnalysis
  /** Import statements forming the cycle (Story 3.2) */
  importTraces?: ImportTrace[]
  /** Recommended fix strategies, sorted by suitability (Story 3.3) */
  fixStrategies?: FixStrategy[]
  /** Impact assessment with blast radius analysis (Story 3.6) */
  impactAssessment?: ImpactAssessment
}

/**
 * ImportTrace - Single import statement that contributes to a cycle
 *
 * Matches Go: pkg/types/import_trace.go (Story 3.2)
 */
export interface ImportTrace {
  /** Package containing the import */
  fromPackage: string
  /** Package being imported */
  toPackage: string
  /** Relative path to the file containing the import */
  filePath: string
  /** 1-based line number of the import statement */
  lineNumber: number
  /** Actual import/require statement text */
  statement: string
  /** Import style classification */
  importType: ImportType
  /** Specific imports (empty for namespace/side-effect imports) */
  symbols?: string[]
}

/**
 * ImportType - Classification of import style
 *
 * Matches Go: pkg/types/import_trace.go ImportType constants (Story 3.2)
 */
export type ImportType =
  | 'esm-named' // import { foo } from 'bar'
  | 'esm-default' // import foo from 'bar'
  | 'esm-namespace' // import * as foo from 'bar'
  | 'esm-side-effect' // import 'bar'
  | 'esm-dynamic' // import('bar')
  | 'cjs-require' // require('bar')

/**
 * RootCauseAnalysis - Analysis of why a circular dependency exists
 *
 * Matches Go: pkg/types/root_cause.go (Story 3.1)
 */
export interface RootCauseAnalysis {
  /** Package most likely responsible for the cycle */
  originatingPackage: string
  /** The specific dependency creating the cycle */
  problematicDependency: RootCauseEdge
  /** Confidence score (0-100) indicating analysis certainty */
  confidence: number
  /** Human-readable description of the root cause */
  explanation: string
  /** Ordered dependency chain forming the cycle */
  chain: RootCauseEdge[]
  /** The edge most likely to break if removed (optional) */
  criticalEdge?: RootCauseEdge
}

/**
 * RootCauseEdge - Single dependency relationship in root cause analysis
 *
 * Matches Go: pkg/types/root_cause.go (Story 3.1)
 */
export interface RootCauseEdge {
  /** Source package */
  from: string
  /** Target package */
  to: string
  /** Dependency type */
  type: DependencyType
  /** If true, this edge is key to breaking the cycle */
  critical: boolean
}

/**
 * FixStrategy - Suggested fix for circular dependency
 *
 * Matches Go: pkg/types/fix_strategy.go (Story 3.3)
 */
export interface FixStrategy {
  /** Strategy type identifier */
  type: FixStrategyType
  /** Human-readable strategy name */
  name: string
  /** Detailed description of what this strategy does */
  description: string
  /** Suitability score (1-10) indicating how well this strategy fits */
  suitability: number
  /** Estimated implementation effort */
  effort: EffortLevel
  /** Advantages of this strategy for this specific cycle */
  pros: string[]
  /** Disadvantages of this strategy for this specific cycle */
  cons: string[]
  /** True if this is the top recommendation */
  recommended: boolean
  /** Packages that would need modification */
  targetPackages: string[]
  /** Suggested name for extracted module (if applicable) */
  newPackageName?: string
  /** Step-by-step fix guide (Story 3.4) */
  guide?: FixGuide
  /** Detailed refactoring complexity for this strategy (Story 3.5) */
  complexity?: RefactoringComplexity
  /** Before/after comparison data for visualization (Story 3.7) */
  beforeAfterExplanation?: BeforeAfterExplanation
}

/**
 * FixGuide - Step-by-step instructions for implementing a fix strategy
 *
 * Matches Go: pkg/types/fix_guide.go (Story 3.4)
 */
export interface FixGuide {
  /** Links this guide to a specific strategy */
  strategyType: FixStrategyType
  /** Guide headline */
  title: string
  /** Brief overview of what this guide accomplishes */
  summary: string
  /** Ordered implementation instructions */
  steps: FixStep[]
  /** Steps to confirm the fix worked */
  verification: FixStep[]
  /** Instructions to undo the changes */
  rollback?: RollbackInstructions
  /** Approximate time to complete (e.g., "15-30 minutes") */
  estimatedTime: string
}

/**
 * FixStep - Single step in a fix guide
 *
 * Matches Go: pkg/types/fix_guide.go (Story 3.4)
 */
export interface FixStep {
  /** Step number (1-based) */
  number: number
  /** Short description of this step */
  title: string
  /** Detailed instructions */
  description: string
  /** File to modify (if applicable) */
  filePath?: string
  /** Current code (if applicable) */
  codeBefore?: CodeSnippet
  /** Desired code (if applicable) */
  codeAfter?: CodeSnippet
  /** Terminal command to run (if applicable) */
  command?: CommandStep
  /** What should happen after this step */
  expectedOutcome?: string
}

/**
 * CodeSnippet - Code example in a fix step
 *
 * Matches Go: pkg/types/fix_guide.go (Story 3.4)
 */
export interface CodeSnippet {
  /** Syntax highlighting hint (e.g., "typescript", "json") */
  language: string
  /** Actual code content */
  code: string
  /** Approximate line number (for context) */
  startLine?: number
}

/**
 * CommandStep - Terminal command in a fix step
 *
 * Matches Go: pkg/types/fix_guide.go (Story 3.4)
 */
export interface CommandStep {
  /** Exact command to run */
  command: string
  /** Where to run the command (relative to workspace root) */
  workingDirectory?: string
  /** Explanation of what this command does */
  description?: string
}

/**
 * RollbackInstructions - Steps to undo fix changes
 *
 * Matches Go: pkg/types/fix_guide.go (Story 3.4)
 */
export interface RollbackInstructions {
  /** Git commands to revert (if in a git repo) */
  gitCommands?: string[]
  /** Non-git rollback instructions */
  manualSteps?: string[]
  /** Caution message about rollback */
  warning?: string
}

/**
 * FixStrategyType - Identifies the approach to resolve the cycle
 *
 * Matches Go: pkg/types/fix_strategy.go FixStrategyType constants
 */
export type FixStrategyType =
  | 'extract-module' // Move shared code to new package
  | 'dependency-injection' // Invert the dependency relationship
  | 'boundary-refactoring' // Restructure module boundaries

/**
 * EffortLevel - Estimates implementation difficulty
 *
 * Matches Go: pkg/types/fix_strategy.go EffortLevel constants
 */
export type EffortLevel =
  | 'low' // Simple changes, < 1 hour
  | 'medium' // Moderate changes, 1-4 hours
  | 'high' // Significant refactoring, > 4 hours

/**
 * CheckResult - Validation-only output for CI/CD
 *
 * Matches Go: pkg/types/check_result.go
 */
export interface CheckResult {
  /** Overall pass/fail status */
  passed: boolean
  /** List of errors found */
  errors: ValidationError[]
  /** List of warnings */
  warnings: ValidationWarning[]
  /** Health score (0-100) */
  healthScore: number
}

/**
 * ValidationError - Error found during validation check
 *
 * Matches Go: pkg/types/validation_error.go
 */
export interface ValidationError {
  /** Error code */
  code: string
  /** Error message */
  message: string
  /** Related file path (optional) */
  file?: string
  /** Line number (optional) */
  line?: number
}

/**
 * ValidationWarning - Warning found during validation check
 *
 * Matches Go: pkg/types/validation_warning.go
 */
export interface ValidationWarning {
  /** Warning code */
  code: string
  /** Warning message */
  message: string
  /** Related file path (optional) */
  file?: string
}

/**
 * AnalysisMetadata - Metadata about the analysis execution
 *
 * Matches Go: pkg/types/analysis_metadata.go
 */
export interface AnalysisMetadata {
  /** MonoGuard version */
  version: string
  /** Analysis duration in milliseconds */
  durationMs: number
  /** Number of files processed */
  filesProcessed: number
  /** Workspace type detected */
  workspaceType: WorkspaceType
}

/**
 * VersionConflictInfo - Dependency with multiple versions across packages
 *
 * Matches Go: pkg/types/version_conflict.go (Story 2.4)
 */
export interface VersionConflictInfo {
  /** External package name with conflicting versions */
  packageName: string
  /** List of conflicting versions and their consumers */
  conflictingVersions: VersionConflictVersion[]
  /** Severity based on semver difference */
  severity: 'critical' | 'warning' | 'info'
  /** Suggested resolution action */
  resolution: string
  /** Impact description */
  impact: string
}

/**
 * VersionConflictVersion - One version and which packages use it
 * Named differently from domain.ts ConflictingVersion to match Go struct
 */
export interface VersionConflictVersion {
  /** The version string */
  version: string
  /** Workspace packages using this version */
  packages: string[]
  /** True if major version differs from others */
  isBreaking: boolean
  /** Dependency type: production, development, peer */
  depType: 'production' | 'development' | 'peer'
}

/**
 * HealthScoreDetails - Complete health score with breakdown
 *
 * Matches Go: pkg/types/health_score.go (Story 2.5)
 */
export interface HealthScoreDetails {
  /** Overall score 0-100 */
  overall: number
  /** Rating classification */
  rating: 'excellent' | 'good' | 'fair' | 'poor' | 'critical'
  /** Individual factor scores */
  breakdown: ScoreBreakdown
  /** Detailed factor information */
  factors: HealthScoreFactor[]
  /** ISO 8601 timestamp */
  updatedAt: string
}

/**
 * ScoreBreakdown - Individual factor scores
 */
export interface ScoreBreakdown {
  /** Score from circular dependency analysis */
  circularScore: number
  /** Score from version conflict analysis */
  conflictScore: number
  /** Score from dependency depth analysis */
  depthScore: number
  /** Score from package coupling analysis */
  couplingScore: number
}

/**
 * HealthScoreFactor - Single factor in health calculation
 * Named differently from domain.ts HealthFactor to match Go struct with weightedScore
 */
export interface HealthScoreFactor {
  /** Factor name */
  name: string
  /** Raw score 0-100 */
  score: number
  /** Weight 0.0-1.0 */
  weight: number
  /** Contribution to overall (score * weight) */
  weightedScore: number
  /** Human-readable description */
  description: string
  /** Suggested improvements */
  recommendations: string[]
}

/**
 * RefactoringComplexity - Detailed breakdown of fix complexity
 *
 * Matches Go: pkg/types/refactoring_complexity.go (Story 3.5)
 */
export interface RefactoringComplexity {
  /** Overall complexity score (1-10) */
  score: number
  /** Human-readable time range (e.g., "30-60 minutes") */
  estimatedTime: string
  /** Individual factor contributions */
  breakdown: ComplexityBreakdown
  /** Human-readable complexity summary */
  explanation: string
}

/**
 * ComplexityBreakdown - How each factor contributes to the score
 *
 * Matches Go: pkg/types/refactoring_complexity.go (Story 3.5)
 */
export interface ComplexityBreakdown {
  /** Number of source files that need changes */
  filesAffected: ComplexityFactor
  /** Number of import statements to modify */
  importsToChange: ComplexityFactor
  /** Dependency chain depth */
  chainDepth: ComplexityFactor
  /** Number of packages in the cycle */
  packagesInvolved: ComplexityFactor
  /** Whether external dependencies are involved */
  externalDependencies: ComplexityFactor
}

/**
 * ComplexityFactor - Single factor in complexity calculation
 *
 * Matches Go: pkg/types/refactoring_complexity.go (Story 3.5)
 */
export interface ComplexityFactor {
  /** Raw value for this factor */
  value: number
  /** Factor weight (0.0-1.0) */
  weight: number
  /** Weighted score contribution */
  contribution: number
  /** What this factor measures */
  description: string
}

/**
 * ImpactAssessment - Blast radius analysis for a circular dependency
 *
 * Matches Go: pkg/types/impact_assessment.go (Story 3.6)
 */
export interface ImpactAssessment {
  /** Packages directly in the cycle */
  directParticipants: string[]
  /** Packages that depend on cycle participants */
  indirectDependents: IndirectDependent[]
  /** Count of all affected packages (direct + indirect) */
  totalAffected: number
  /** Proportion of workspace affected (0.0-1.0) */
  affectedPercentage: number
  /** Human-readable percentage (e.g., "25%") */
  affectedPercentageDisplay: string
  /** Impact severity classification: critical, high, medium, low */
  riskLevel: 'critical' | 'high' | 'medium' | 'low'
  /** Explanation of risk classification */
  riskExplanation: string
  /** Visualization-ready data */
  rippleEffect?: RippleEffect
}

/**
 * IndirectDependent - Package that depends on a cycle participant
 *
 * Matches Go: pkg/types/impact_assessment.go (Story 3.6)
 */
export interface IndirectDependent {
  /** The affected package */
  packageName: string
  /** Which cycle participant this package depends on */
  dependsOn: string
  /** Hops from the cycle (1 = direct dependent) */
  distance: number
  /** Full path from cycle to this package */
  dependencyPath: string[]
}

/**
 * RippleEffect - Data for visualization
 *
 * Matches Go: pkg/types/impact_assessment.go (Story 3.6)
 */
export interface RippleEffect {
  /** Packages grouped by distance from cycle */
  layers: RippleLayer[]
  /** Maximum distance from cycle */
  totalLayers: number
}

/**
 * RippleLayer - Packages at a specific distance from the cycle
 *
 * Matches Go: pkg/types/impact_assessment.go (Story 3.6)
 */
export interface RippleLayer {
  /** Distance from cycle (0 = direct participants) */
  distance: number
  /** Packages at this distance */
  packages: string[]
  /** Count of packages */
  count: number
}

/**
 * BeforeAfterExplanation - Visual comparison data for fix strategies
 *
 * Matches Go: pkg/types/before_after_explanation.go (Story 3.7)
 */
export interface BeforeAfterExplanation {
  /** Dependency graph before the fix */
  currentState: StateDiagram
  /** Dependency graph after the fix */
  proposedState: StateDiagram
  /** Changes required to package.json files */
  packageJsonDiffs: PackageJsonDiff[]
  /** Changes required to import statements */
  importDiffs: ImportDiff[]
  /** Human-readable explanation */
  explanation: FixExplanation
  /** Potential side effects */
  warnings: SideEffectWarning[]
}

/**
 * StateDiagram - D3.js-compatible visualization data
 *
 * Matches Go: pkg/types/before_after_explanation.go (Story 3.7)
 */
export interface StateDiagram {
  /** Packages in the diagram */
  nodes: DiagramNode[]
  /** Dependency relationships */
  edges: DiagramEdge[]
  /** The cycle path (only in currentState) */
  highlightedPath?: string[]
  /** Whether this state has no cycle */
  cycleResolved: boolean
}

/**
 * DiagramNode - Package in the visualization
 *
 * Matches Go: pkg/types/before_after_explanation.go (Story 3.7)
 */
export interface DiagramNode {
  /** Package name (used for edge references) */
  id: string
  /** Display name */
  label: string
  /** Whether this package is part of the cycle */
  isInCycle: boolean
  /** Whether this package is newly created by the fix */
  isNew: boolean
  /** Node category for visualization styling */
  nodeType: DiagramNodeType
}

/**
 * DiagramNodeType - Categorizes nodes for visualization styling
 *
 * Matches Go: pkg/types/before_after_explanation.go (Story 3.7)
 */
export type DiagramNodeType = 'cycle' | 'affected' | 'new' | 'unchanged'

/**
 * DiagramEdge - Dependency relationship
 *
 * Matches Go: pkg/types/before_after_explanation.go (Story 3.7)
 */
export interface DiagramEdge {
  /** Dependent package */
  from: string
  /** Dependency */
  to: string
  /** Whether this edge is part of the cycle */
  isInCycle: boolean
  /** Whether this edge will be removed by the fix */
  isRemoved: boolean
  /** Whether this edge is added by the fix */
  isNew: boolean
  /** Edge category for visualization styling */
  edgeType: DiagramEdgeType
}

/**
 * DiagramEdgeType - Categorizes edges for visualization styling
 *
 * Matches Go: pkg/types/before_after_explanation.go (Story 3.7)
 */
export type DiagramEdgeType = 'cycle' | 'removed' | 'new' | 'unchanged'

/**
 * PackageJsonDiff - Changes to a package.json file
 *
 * Matches Go: pkg/types/before_after_explanation.go (Story 3.7)
 */
export interface PackageJsonDiff {
  /** Package being modified */
  packageName: string
  /** Relative path to package.json */
  filePath: string
  /** Dependencies to add */
  dependenciesToAdd: DependencyChange[]
  /** Dependencies to remove */
  dependenciesToRemove: DependencyChange[]
  /** Human-readable change description */
  summary: string
}

/**
 * DependencyChange - Dependency addition or removal
 *
 * Matches Go: pkg/types/before_after_explanation.go (Story 3.7)
 */
export interface DependencyChange {
  /** Dependency package name */
  name: string
  /** Version specifier (e.g., "workspace:*", "^1.0.0") */
  version?: string
  /** dependencies vs devDependencies */
  dependencyType: string
}

/**
 * ImportDiff - Changes to import statements in a file
 *
 * Matches Go: pkg/types/before_after_explanation.go (Story 3.7)
 */
export interface ImportDiff {
  /** File containing imports */
  filePath: string
  /** Package containing this file */
  packageName: string
  /** Import statements to remove */
  importsToRemove: ImportChange[]
  /** Import statements to add */
  importsToAdd: ImportChange[]
  /** Location hint (if available from ImportTraces) */
  lineNumber?: number
}

/**
 * ImportChange - Import statement change
 *
 * Matches Go: pkg/types/before_after_explanation.go (Story 3.7)
 */
export interface ImportChange {
  /** Full import statement */
  statement: string
  /** Package being imported from */
  fromPackage: string
  /** What is being imported */
  importedNames?: string[]
}

/**
 * FixExplanation - Human-readable explanation of the fix
 *
 * Matches Go: pkg/types/before_after_explanation.go (Story 3.7)
 */
export interface FixExplanation {
  /** 1-2 sentence overview */
  summary: string
  /** How this resolves the cycle */
  whyItWorks: string
  /** What code changes are required */
  highLevelChanges: string[]
  /** Confidence in the fix (0.0-1.0) */
  confidence: number
}

/**
 * SideEffectWarning - Potential side effect of the fix
 *
 * Matches Go: pkg/types/before_after_explanation.go (Story 3.7)
 */
export interface SideEffectWarning {
  /** Importance level: info, warning, critical */
  severity: WarningSeverity
  /** Short description */
  title: string
  /** Detailed description */
  description: string
  /** Packages that may be affected */
  affectedPackages?: string[]
}

/**
 * WarningSeverity - Importance of a warning
 *
 * Matches Go: pkg/types/before_after_explanation.go (Story 3.7)
 */
export type WarningSeverity = 'info' | 'warning' | 'critical'

export { RiskLevel, Severity } from '../common'
// Re-export for convenience
export type { WorkspaceType } from './graph'
