/**
 * JSON Bridge Module for WASM Communication
 *
 * Provides type-safe wrappers for calling WASM functions with JSON serialization.
 *
 * @module bridge
 */

import type { Result } from '@monoguard/types'

/**
 * MonoGuard global interface available after WASM initialization
 */
interface MonoGuardGlobal {
  getVersion(): string
  analyze(input: string): string
  check(input: string): string
}

declare global {
  interface Window {
    MonoGuard?: MonoGuardGlobal
  }
}

/**
 * Create an error result
 */
function createError<T>(code: string, message: string): Result<T> {
  return {
    data: null,
    error: { code, message },
  }
}

/**
 * Call a WASM function with JSON input and parse the result
 *
 * @param funcName - Name of the MonoGuard function to call
 * @param input - Input data to serialize as JSON
 * @returns Parsed Result<T> from WASM
 *
 * @example
 * ```typescript
 * const result = callWasm<AnalysisResult>('analyze', { files: {...} });
 * if (isSuccess(result)) {
 *   console.log(result.data.healthScore);
 * }
 * ```
 */
export function callWasm<T>(funcName: keyof MonoGuardGlobal, input: unknown): Result<T> {
  if (typeof window === 'undefined') {
    return createError('WASM_NOT_AVAILABLE', 'WASM is only available in browser environment')
  }

  if (typeof window.MonoGuard === 'undefined') {
    return createError('WASM_NOT_INITIALIZED', 'WASM not loaded. Call loadWasm() first.')
  }

  const func = window.MonoGuard[funcName]
  if (typeof func !== 'function') {
    return createError('WASM_FUNCTION_NOT_FOUND', `MonoGuard.${funcName} is not a function`)
  }

  try {
    const inputJson = JSON.stringify(input)
    const resultJson = (func as (input: string) => string).call(window.MonoGuard, inputJson)
    return parseWasmResult<T>(resultJson)
  } catch (err) {
    return createError(
      'WASM_CALL_FAILED',
      err instanceof Error ? err.message : 'Unknown error calling WASM function'
    )
  }
}

/**
 * Call a WASM function that takes no input
 *
 * @param funcName - Name of the MonoGuard function to call
 * @returns Parsed Result<T> from WASM
 *
 * @example
 * ```typescript
 * const result = callWasmNoInput<VersionInfo>('getVersion');
 * if (isSuccess(result)) {
 *   console.log(result.data.version);
 * }
 * ```
 */
export function callWasmNoInput<T>(funcName: keyof MonoGuardGlobal): Result<T> {
  if (typeof window === 'undefined') {
    return createError('WASM_NOT_AVAILABLE', 'WASM is only available in browser environment')
  }

  if (typeof window.MonoGuard === 'undefined') {
    return createError('WASM_NOT_INITIALIZED', 'WASM not loaded. Call loadWasm() first.')
  }

  const func = window.MonoGuard[funcName]
  if (typeof func !== 'function') {
    return createError('WASM_FUNCTION_NOT_FOUND', `MonoGuard.${funcName} is not a function`)
  }

  try {
    const resultJson = (func as () => string).call(window.MonoGuard)
    return parseWasmResult<T>(resultJson)
  } catch (err) {
    return createError(
      'WASM_CALL_FAILED',
      err instanceof Error ? err.message : 'Unknown error calling WASM function'
    )
  }
}

/**
 * Parse WASM result JSON string into Result<T>
 */
function parseWasmResult<T>(resultJson: string): Result<T> {
  try {
    const parsed = JSON.parse(resultJson) as Result<T>

    // Validate result structure
    if (typeof parsed !== 'object' || parsed === null) {
      return createError('WASM_INVALID_RESULT', 'WASM returned invalid result structure')
    }

    // Ensure proper Result structure
    if (!('data' in parsed) || !('error' in parsed)) {
      return createError('WASM_INVALID_RESULT', 'WASM result missing data or error field')
    }

    return parsed
  } catch (err) {
    return createError(
      'WASM_PARSE_ERROR',
      `Failed to parse WASM result: ${err instanceof Error ? err.message : 'Unknown error'}`
    )
  }
}
