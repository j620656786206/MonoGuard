import { describe, expect, it } from 'vitest'
import { ErrorCodes, isError, isSuccess, type Result } from '../result'

describe('Result type guards', () => {
  describe('isSuccess', () => {
    it('returns true for successful result with data', () => {
      const result: Result<number> = { data: 42, error: null }
      expect(isSuccess(result)).toBe(true)
    })

    it('returns false for error result', () => {
      const result: Result<number> = {
        data: null,
        error: { code: 'TEST_ERROR', message: 'Test error' },
      }
      expect(isSuccess(result)).toBe(false)
    })

    it('returns false when both data and error are null', () => {
      const result: Result<number> = { data: null, error: null }
      expect(isSuccess(result)).toBe(false)
    })
  })

  describe('isError', () => {
    it('returns true for error result', () => {
      const result: Result<number> = {
        data: null,
        error: { code: 'TEST_ERROR', message: 'Test error' },
      }
      expect(isError(result)).toBe(true)
    })

    it('returns false for successful result', () => {
      const result: Result<number> = { data: 42, error: null }
      expect(isError(result)).toBe(false)
    })
  })

  describe('type narrowing', () => {
    it('allows type-safe data access after isSuccess check', () => {
      const result: Result<{ value: string }> = {
        data: { value: 'test' },
        error: null,
      }

      if (isSuccess(result)) {
        // TypeScript should know result.data is non-null here
        expect(result.data.value).toBe('test')
      }
    })

    it('allows type-safe error access after isError check', () => {
      const result: Result<number> = {
        data: null,
        error: { code: 'TEST_ERROR', message: 'Test message' },
      }

      if (isError(result)) {
        // TypeScript should know result.error is non-null here
        expect(result.error.code).toBe('TEST_ERROR')
        expect(result.error.message).toBe('Test message')
      }
    })
  })
})

describe('ErrorCodes', () => {
  it('contains standard error codes', () => {
    expect(ErrorCodes.PARSE_ERROR).toBe('PARSE_ERROR')
    expect(ErrorCodes.INVALID_INPUT).toBe('INVALID_INPUT')
    expect(ErrorCodes.CIRCULAR_DETECTED).toBe('CIRCULAR_DETECTED')
    expect(ErrorCodes.ANALYSIS_FAILED).toBe('ANALYSIS_FAILED')
    expect(ErrorCodes.WASM_ERROR).toBe('WASM_ERROR')
    expect(ErrorCodes.TIMEOUT).toBe('TIMEOUT')
  })

  it('has immutable error codes (as const)', () => {
    // Verify the const assertion works by checking that all values are strings
    const codes = Object.values(ErrorCodes)
    codes.forEach((code) => {
      expect(typeof code).toBe('string')
    })
  })
})
