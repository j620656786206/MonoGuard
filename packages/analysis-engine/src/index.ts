/**
 * MonoGuard Analysis Engine
 *
 * TypeScript adapter for the Go WASM analysis engine.
 * Provides type-safe access to dependency analysis functions.
 *
 * @packageDocumentation
 * @module @monoguard/analysis-engine
 *
 * @example
 * ```typescript
 * import {
 *   MonoGuardAnalyzer,
 *   isSuccess,
 *   isError,
 *   type WorkspaceInput,
 * } from '@monoguard/analysis-engine';
 *
 * // Initialize analyzer
 * const analyzer = new MonoGuardAnalyzer();
 * await analyzer.init({ wasmPath: '/monoguard.wasm' });
 *
 * // Analyze workspace
 * const input: WorkspaceInput = {
 *   files: {
 *     'package.json': JSON.stringify({ workspaces: ['packages/*'] }),
 *     'packages/core/package.json': JSON.stringify({ name: '@mono/core' }),
 *   },
 *   config: {
 *     exclude: ['packages/legacy-*'],
 *   },
 * };
 *
 * const result = await analyzer.analyze(input);
 *
 * if (isSuccess(result)) {
 *   console.log(`Health Score: ${result.data.healthScore}`);
 *   console.log(`Packages: ${result.data.packages}`);
 * } else {
 *   console.error(`Error: ${result.error.code}`);
 * }
 * ```
 */

// Re-export commonly used types from @monoguard/types
export type {
  Result,
  ResultError,
  VersionInfo,
  WasmLoaderOptions,
} from '@monoguard/types'
// Re-export type guards from @monoguard/types
export { isError, isSuccess } from '@monoguard/types'
// Analyzer exports
export {
  type AnalysisConfig,
  type CircularDependencyInfo,
  createAnalyzer,
  type DependencyEdge,
  type DependencyGraph,
  type HealthFactor,
  type HealthScoreDetails,
  MonoGuardAnalyzer,
  type PackageNode,
  type VersionConflictInfo,
  type WasmAnalysisResult,
  type WasmCheckResult,
  type WorkspaceInput,
} from './analyzer'
// Loader exports
export { isWasmInitialized, loadWasm, WasmLoadError } from './loader'
