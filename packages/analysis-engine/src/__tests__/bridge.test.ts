import { isError, isSuccess } from '@monoguard/types'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { callWasm, callWasmNoInput } from '../bridge'

// biome-ignore lint/suspicious/noExplicitAny: Test mocking requires window global manipulation
type WindowWithMonoGuard = Window & { MonoGuard?: any; Go?: any }
const testWindow = window as WindowWithMonoGuard

describe('JSON Bridge', () => {
  beforeEach(() => {
    testWindow.MonoGuard = {
      getVersion: vi.fn().mockReturnValue(
        JSON.stringify({
          data: { version: '0.1.0' },
          error: null,
        })
      ),
      analyze: vi.fn().mockReturnValue(
        JSON.stringify({
          data: { healthScore: 85 },
          error: null,
        })
      ),
      check: vi.fn().mockReturnValue(
        JSON.stringify({
          data: { passed: true },
          error: null,
        })
      ),
    }
  })

  afterEach(() => {
    delete testWindow.MonoGuard
  })

  describe('callWasm', () => {
    it('should call WASM function with JSON input', () => {
      const input = { files: { 'package.json': '{}' } }
      callWasm('analyze', input)

      expect(testWindow.MonoGuard.analyze).toHaveBeenCalledWith(JSON.stringify(input))
    })

    it('should return parsed result on success', () => {
      const result = callWasm<{ healthScore: number }>('analyze', {})

      expect(isSuccess(result)).toBe(true)
      if (isSuccess(result)) {
        expect(result.data.healthScore).toBe(85)
      }
    })

    it('should return error if WASM not initialized', () => {
      delete testWindow.MonoGuard

      const result = callWasm('analyze', {})

      expect(isError(result)).toBe(true)
      if (isError(result)) {
        expect(result.error.code).toBe('WASM_NOT_INITIALIZED')
      }
    })

    it('should return error if function not found', () => {
      testWindow.MonoGuard = {}

      const result = callWasm('analyze', {})

      expect(isError(result)).toBe(true)
      if (isError(result)) {
        expect(result.error.code).toBe('WASM_FUNCTION_NOT_FOUND')
      }
    })

    it('should return error if WASM throws', () => {
      testWindow.MonoGuard.analyze = vi.fn().mockImplementation(() => {
        throw new Error('WASM crashed')
      })

      const result = callWasm('analyze', {})

      expect(isError(result)).toBe(true)
      if (isError(result)) {
        expect(result.error.code).toBe('WASM_CALL_FAILED')
        expect(result.error.message).toContain('WASM crashed')
      }
    })

    it('should return error if JSON parse fails', () => {
      testWindow.MonoGuard.analyze = vi.fn().mockReturnValue('invalid json')

      const result = callWasm('analyze', {})

      expect(isError(result)).toBe(true)
      if (isError(result)) {
        expect(result.error.code).toBe('WASM_PARSE_ERROR')
      }
    })

    it('should return error for invalid result structure', () => {
      testWindow.MonoGuard.analyze = vi.fn().mockReturnValue(JSON.stringify({ foo: 'bar' }))

      const result = callWasm('analyze', {})

      expect(isError(result)).toBe(true)
      if (isError(result)) {
        expect(result.error.code).toBe('WASM_INVALID_RESULT')
      }
    })

    it('should pass through WASM error response', () => {
      testWindow.MonoGuard.analyze = vi.fn().mockReturnValue(
        JSON.stringify({
          data: null,
          error: { code: 'PARSE_ERROR', message: 'Invalid workspace' },
        })
      )

      const result = callWasm('analyze', {})

      expect(isError(result)).toBe(true)
      if (isError(result)) {
        expect(result.error.code).toBe('PARSE_ERROR')
        expect(result.error.message).toBe('Invalid workspace')
      }
    })
  })

  describe('callWasmNoInput', () => {
    it('should call WASM function without input', () => {
      callWasmNoInput('getVersion')

      expect(testWindow.MonoGuard.getVersion).toHaveBeenCalledWith()
    })

    it('should return parsed result on success', () => {
      const result = callWasmNoInput<{ version: string }>('getVersion')

      expect(isSuccess(result)).toBe(true)
      if (isSuccess(result)) {
        expect(result.data.version).toBe('0.1.0')
      }
    })

    it('should return error if WASM not initialized', () => {
      delete testWindow.MonoGuard

      const result = callWasmNoInput('getVersion')

      expect(isError(result)).toBe(true)
      if (isError(result)) {
        expect(result.error.code).toBe('WASM_NOT_INITIALIZED')
      }
    })
  })
})
