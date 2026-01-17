import type { AnalysisResult, CheckResult } from '../analysis/results'
import type { Result } from '../result'

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
 * const analyzer: MonoGuardAnalyzer = await loadWasm();
 * const result = await analyzer.analyze(JSON.stringify(workspaceConfig));
 * if (isSuccess(result)) {
 *   console.log(result.data.healthScore);
 * }
 * ```
 */
export interface MonoGuardAnalyzer {
  /**
   * Get MonoGuard version
   * @returns Version information
   */
  getVersion(): Promise<Result<VersionInfo>>

  /**
   * Analyze workspace dependencies
   * @param input JSON string of workspace configuration
   * @returns Complete analysis result
   */
  analyze(input: string): Promise<Result<AnalysisResult>>

  /**
   * Check workspace for CI/CD validation
   * @param input JSON string of workspace configuration
   * @returns Pass/fail result with errors
   */
  check(input: string): Promise<Result<CheckResult>>
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
