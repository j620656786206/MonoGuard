/**
 * MonoGuard Analyzer Module
 *
 * Provides a type-safe wrapper for the WASM analysis engine.
 *
 * @module analyzer
 */

import type { Result, VersionInfo, WasmLoaderOptions } from '@monoguard/types'
import { callWasm, callWasmNoInput } from './bridge'
import { isWasmInitialized, loadWasm } from './loader'

/**
 * Input for analyze and check functions
 */
export interface WorkspaceInput {
  /** Map of filename to file content */
  files: Record<string, string>
  /** Optional analysis configuration */
  config?: AnalysisConfig
}

/**
 * Analysis configuration options (Story 2.6)
 */
export interface AnalysisConfig {
  /** Patterns to exclude from analysis (exact, glob, or regex:) */
  exclude?: string[]
}

/**
 * Analysis result from WASM (matches Go types.AnalysisResult)
 */
export interface WasmAnalysisResult {
  /** Architecture health score (0-100) */
  healthScore: number
  /** Detailed health score breakdown (Story 2.5) */
  healthScoreDetails?: HealthScoreDetails
  /** Total packages analyzed */
  packages: number
  /** Excluded packages count (Story 2.6) */
  excludedPackages?: number
  /** Dependency graph */
  graph?: DependencyGraph
  /** Detected circular dependencies */
  circularDependencies?: CircularDependencyInfo[]
  /** Version conflicts (Story 2.4) */
  versionConflicts?: VersionConflictInfo[]
  /** ISO 8601 timestamp */
  createdAt?: string
}

/**
 * Health score breakdown (Story 2.5)
 */
export interface HealthScoreDetails {
  overall: number
  rating: 'excellent' | 'good' | 'fair' | 'poor' | 'critical'
  breakdown: {
    circularScore: number
    conflictScore: number
    depthScore: number
    couplingScore: number
  }
  factors: HealthFactor[]
  updatedAt: string
}

export interface HealthFactor {
  name: string
  score: number
  weight: number
  weightedScore: number
  description: string
  recommendations: string[]
}

/**
 * Dependency graph structure
 */
export interface DependencyGraph {
  nodes: Record<string, PackageNode>
  edges: DependencyEdge[]
  rootPath: string
  workspaceType: string
}

export interface PackageNode {
  name: string
  version: string
  path: string
  dependencies: string[]
  devDependencies: string[]
  peerDependencies: string[]
  optionalDependencies: string[]
  externalDeps?: Record<string, string>
  externalDevDeps?: Record<string, string>
  excluded?: boolean
}

export interface DependencyEdge {
  from: string
  to: string
  type: 'production' | 'development' | 'peer' | 'optional'
  versionRange: string
}

/**
 * Circular dependency info (Story 2.3)
 */
export interface CircularDependencyInfo {
  cycle: string[]
  type: 'direct' | 'indirect'
  severity: 'critical' | 'warning' | 'info'
  depth: number
  impact?: string
}

/**
 * Version conflict info (Story 2.4)
 */
export interface VersionConflictInfo {
  packageName: string
  conflictingVersions: Array<{
    version: string
    usedBy: string
    depType: string
  }>
  severity: 'critical' | 'warning' | 'info'
  resolution?: string
  impact?: string
}

/**
 * Check result from WASM
 */
export interface WasmCheckResult {
  passed: boolean
  errors: string[]
  placeholder?: boolean
}

/**
 * MonoGuard analyzer wrapper with full type safety
 *
 * @example
 * ```typescript
 * import { MonoGuardAnalyzer, isSuccess } from '@monoguard/analysis-engine';
 *
 * const analyzer = new MonoGuardAnalyzer();
 * await analyzer.init();
 *
 * const result = await analyzer.analyze({
 *   files: { 'package.json': '{"workspaces": ["packages/*"]}' }
 * });
 *
 * if (isSuccess(result)) {
 *   console.log(`Health Score: ${result.data.healthScore}`);
 * }
 * ```
 */
export class MonoGuardAnalyzer {
  private initialized = false

  /**
   * Initialize the WASM module
   *
   * Must be called before analyze/check/getVersion.
   * Safe to call multiple times - subsequent calls are no-ops.
   *
   * @param options - WASM loading options
   * @throws {WasmLoadError} If loading fails or times out
   *
   * @example
   * ```typescript
   * const analyzer = new MonoGuardAnalyzer();
   * await analyzer.init({ wasmPath: '/wasm/monoguard.wasm', timeout: 5000 });
   * ```
   */
  async init(options?: WasmLoaderOptions): Promise<void> {
    if (this.initialized && isWasmInitialized()) {
      return
    }

    await loadWasm(options)
    this.initialized = true
  }

  /**
   * Check if the analyzer is initialized
   *
   * @returns true if WASM is loaded and ready
   */
  isInitialized(): boolean {
    return this.initialized && isWasmInitialized()
  }

  /**
   * Analyze workspace dependencies
   *
   * Performs comprehensive analysis including:
   * - Dependency graph construction
   * - Circular dependency detection
   * - Version conflict detection
   * - Health score calculation
   *
   * @param input - Workspace files and optional configuration
   * @returns Analysis result with health score, cycles, conflicts, and graph
   *
   * @example
   * ```typescript
   * const result = await analyzer.analyze({
   *   files: {
   *     'package.json': '{"workspaces": ["packages/*"]}',
   *     'packages/core/package.json': '{"name": "@mono/core"}'
   *   },
   *   config: { exclude: ['packages/legacy-*'] }
   * });
   *
   * if (isSuccess(result)) {
   *   console.log(`Health Score: ${result.data.healthScore}`);
   *   console.log(`Packages: ${result.data.packages}`);
   *   console.log(`Cycles: ${result.data.circularDependencies?.length ?? 0}`);
   * }
   * ```
   */
  async analyze(input: WorkspaceInput): Promise<Result<WasmAnalysisResult>> {
    if (!this.isInitialized()) {
      return {
        data: null,
        error: {
          code: 'NOT_INITIALIZED',
          message: 'Analyzer not initialized. Call init() first.',
        },
      }
    }

    return callWasm<WasmAnalysisResult>('analyze', input)
  }

  /**
   * Check workspace for CI/CD validation
   *
   * Lightweight validation suitable for CI/CD pipelines.
   * Returns pass/fail with error details.
   *
   * @param input - Workspace files and optional configuration
   * @returns Pass/fail result with errors
   *
   * @example
   * ```typescript
   * const result = await analyzer.check({
   *   files: { 'package.json': '...' }
   * });
   *
   * if (isSuccess(result) && result.data.passed) {
   *   console.log('All checks passed!');
   * } else if (isSuccess(result)) {
   *   console.log('Check failed:', result.data.errors);
   * }
   * ```
   */
  async check(input: WorkspaceInput): Promise<Result<WasmCheckResult>> {
    if (!this.isInitialized()) {
      return {
        data: null,
        error: {
          code: 'NOT_INITIALIZED',
          message: 'Analyzer not initialized. Call init() first.',
        },
      }
    }

    return callWasm<WasmCheckResult>('check', input)
  }

  /**
   * Get MonoGuard version information
   *
   * Useful for verifying WASM loaded correctly and checking version compatibility.
   *
   * @returns Version info including version string
   *
   * @example
   * ```typescript
   * const result = await analyzer.getVersion();
   * if (isSuccess(result)) {
   *   console.log(`MonoGuard v${result.data.version}`);
   * }
   * ```
   */
  async getVersion(): Promise<Result<VersionInfo>> {
    if (!this.isInitialized()) {
      return {
        data: null,
        error: {
          code: 'NOT_INITIALIZED',
          message: 'Analyzer not initialized. Call init() first.',
        },
      }
    }

    return callWasmNoInput<VersionInfo>('getVersion')
  }
}

/**
 * Create a new analyzer instance
 *
 * Convenience function for creating an analyzer.
 *
 * @returns New MonoGuardAnalyzer instance
 */
export function createAnalyzer(): MonoGuardAnalyzer {
  return new MonoGuardAnalyzer()
}
