import type { AnalysisResult, CheckResult } from '../analysis/results'
import type { Result } from '../result'

/**
 * Workspace input for analyze and check functions
 * Provides type-safe input instead of raw JSON strings
 */
export interface WorkspaceInput {
  /** Map of filename to file content */
  files: Record<string, string>
  /** Optional analysis configuration */
  config?: AnalysisConfig
}

/**
 * Analysis configuration options
 */
export interface AnalysisConfig {
  /** Patterns to exclude from analysis (exact, glob, or regex:) */
  exclude?: string[]
}

/**
 * MonoGuardAnalyzer - WASM adapter interface
 *
 * This interface defines the contract between TypeScript and Go WASM.
 * All methods return Promise<Result<T>> to handle async WASM calls.
 *
 * Implementation is in packages/analysis-engine (Go WASM).
 * TypeScript adapter is implemented in Story 2.7.
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
 *   console.log(result.data.healthScore);
 * }
 * ```
 */
export interface MonoGuardAnalyzer {
  /**
   * Initialize the WASM module
   * Must be called before analyze/check/getVersion
   */
  init(options?: WasmLoaderOptions): Promise<void>

  /**
   * Check if the analyzer is initialized
   * @returns true if WASM is loaded and ready
   */
  isInitialized(): boolean

  /**
   * Get MonoGuard version
   * @returns Version information
   */
  getVersion(): Promise<Result<VersionInfo>>

  /**
   * Analyze workspace dependencies
   * @param input Workspace files and optional configuration
   * @returns Complete analysis result
   */
  analyze(input: WorkspaceInput): Promise<Result<AnalysisResult>>

  /**
   * Check workspace for CI/CD validation
   * @param input Workspace files and optional configuration
   * @returns Pass/fail result with errors
   */
  check(input: WorkspaceInput): Promise<Result<CheckResult>>
}

/**
 * VersionInfo - Version information returned by getVersion
 *
 * Matches Go: pkg/types/version.go
 */
export interface VersionInfo {
  /** Semantic version string (e.g., "0.1.0") */
  version: string
  /** Git commit hash (optional) */
  commit?: string
  /** Build date in ISO 8601 format (optional) */
  buildDate?: string
}

/**
 * WasmLoaderOptions - Options for loading WASM module
 */
export interface WasmLoaderOptions {
  /** Path to monoguard.wasm file */
  wasmPath?: string
  /** Timeout for WASM initialization (ms) */
  timeout?: number
}

/**
 * WasmLoadResult - Result of loading WASM module
 */
export interface WasmLoadResult {
  /** The loaded analyzer instance */
  analyzer: MonoGuardAnalyzer
  /** Time taken to load WASM in milliseconds */
  loadTimeMs: number
}
