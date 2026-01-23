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
  /** Refactoring complexity (1-10) */
  complexity: number
  /** Root cause analysis (Story 3.1) */
  rootCause?: RootCauseAnalysis
  /** Import statements forming the cycle (Story 3.2) */
  importTraces?: ImportTrace[]
  /** Recommended fix strategies, sorted by suitability (Story 3.3) */
  fixStrategies?: FixStrategy[]
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

// Re-export for convenience
export type { WorkspaceType } from './graph'
