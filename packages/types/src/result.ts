/**
 * Result<T> - Unified response type for WASM functions
 *
 * This type MUST match the Go Result struct exactly:
 * ```go
 * type Result struct {
 *     Data  interface{} `json:"data"`
 *     Error *Error      `json:"error"`
 * }
 * ```
 *
 * @example
 * // Success case
 * { data: { healthScore: 85 }, error: null }
 *
 * // Error case
 * { data: null, error: { code: "PARSE_ERROR", message: "Invalid JSON" } }
 */
export interface Result<T> {
  data: T | null
  error: ResultError | null
}

/**
 * ResultError - Error structure matching Go Error type
 *
 * Error codes use UPPER_SNAKE_CASE convention.
 */
export interface ResultError {
  /** Error code in UPPER_SNAKE_CASE (e.g., PARSE_ERROR, CIRCULAR_DETECTED) */
  code: string
  /** Human-readable error message */
  message: string
}

/**
 * Standard error codes used by MonoGuard
 * Must match Go constants in internal/result/result.go
 */
export const ErrorCodes = {
  PARSE_ERROR: 'PARSE_ERROR',
  INVALID_INPUT: 'INVALID_INPUT',
  CIRCULAR_DETECTED: 'CIRCULAR_DETECTED',
  ANALYSIS_FAILED: 'ANALYSIS_FAILED',
  WASM_ERROR: 'WASM_ERROR',
  TIMEOUT: 'TIMEOUT',
} as const

export type ErrorCode = (typeof ErrorCodes)[keyof typeof ErrorCodes]

/**
 * Type guard to check if a Result is successful
 * @param result - The Result to check
 * @returns true if the result has data and no error
 */
export function isSuccess<T>(result: Result<T>): result is Result<T> & { data: T; error: null } {
  return result.error === null && result.data !== null
}

/**
 * Type guard to check if a Result is an error
 * @param result - The Result to check
 * @returns true if the result has an error
 */
export function isError<T>(
  result: Result<T>
): result is Result<T> & { data: null; error: ResultError } {
  return result.error !== null
}
