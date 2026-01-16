import type { DependencyGraph, WorkspaceType } from './graph';

/**
 * AnalysisResult - Complete analysis output
 *
 * Matches Go: pkg/types/analysis_result.go
 * All date fields use ISO 8601 format (e.g., "2026-01-15T10:30:00Z")
 */
export interface AnalysisResult {
  /** Architecture health score (0-100) */
  healthScore: number;
  /** Total packages analyzed */
  packageCount: number;
  /** Detected circular dependencies */
  circularDependencies: CircularDependencyInfo[];
  /** Dependency graph data */
  graph: DependencyGraph;
  /** Analysis metadata */
  metadata: AnalysisMetadata;
  /** ISO 8601 timestamp */
  createdAt: string;
}

/**
 * CircularDependencyInfo - Enhanced circular dependency with fix suggestions
 *
 * Matches Go: pkg/types/circular.go
 */
export interface CircularDependencyInfo {
  /** Packages involved in the cycle (in order) */
  cycle: string[];
  /** Type of circular dependency */
  type: 'direct' | 'indirect';
  /** Severity level */
  severity: 'critical' | 'warning' | 'info';
  /** Impact description */
  impact: string;
  /** Suggested fix strategy */
  fixStrategy?: FixStrategy;
  /** Refactoring complexity (1-10) */
  complexity: number;
}

/**
 * FixStrategy - Suggested fix for circular dependency
 *
 * Matches Go: pkg/types/fix_strategy.go
 */
export interface FixStrategy {
  /** Strategy type */
  type: 'extract_module' | 'dependency_injection' | 'boundary_refactor';
  /** Human-readable description */
  description: string;
  /** Step-by-step instructions */
  steps: string[];
  /** Files that need modification */
  affectedFiles: string[];
}

/**
 * CheckResult - Validation-only output for CI/CD
 *
 * Matches Go: pkg/types/check_result.go
 */
export interface CheckResult {
  /** Overall pass/fail status */
  passed: boolean;
  /** List of errors found */
  errors: ValidationError[];
  /** List of warnings */
  warnings: ValidationWarning[];
  /** Health score (0-100) */
  healthScore: number;
}

/**
 * ValidationError - Error found during validation check
 *
 * Matches Go: pkg/types/validation_error.go
 */
export interface ValidationError {
  /** Error code */
  code: string;
  /** Error message */
  message: string;
  /** Related file path (optional) */
  file?: string;
  /** Line number (optional) */
  line?: number;
}

/**
 * ValidationWarning - Warning found during validation check
 *
 * Matches Go: pkg/types/validation_warning.go
 */
export interface ValidationWarning {
  /** Warning code */
  code: string;
  /** Warning message */
  message: string;
  /** Related file path (optional) */
  file?: string;
}

/**
 * AnalysisMetadata - Metadata about the analysis execution
 *
 * Matches Go: pkg/types/analysis_metadata.go
 */
export interface AnalysisMetadata {
  /** MonoGuard version */
  version: string;
  /** Analysis duration in milliseconds */
  durationMs: number;
  /** Number of files processed */
  filesProcessed: number;
  /** Workspace type detected */
  workspaceType: WorkspaceType;
}

// Re-export for convenience
export type { WorkspaceType } from './graph';
