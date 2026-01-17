import { isError, isSuccess } from '@monoguard/types'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { MonoGuardAnalyzer } from '../analyzer'
import { resetWasmState } from '../loader'

describe('MonoGuardAnalyzer', () => {
  let analyzer: MonoGuardAnalyzer

  beforeEach(() => {
    analyzer = new MonoGuardAnalyzer()
    resetWasmState()

    // Mock window.Go for WASM initialization
    ;(window as any).Go = class {
      importObject = {}
      async run() {}
    }

    // Mock WASM fetch
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      arrayBuffer: () => Promise.resolve(new ArrayBuffer(8)),
    })

    // Mock WebAssembly.instantiate
    global.WebAssembly = {
      instantiate: vi.fn().mockResolvedValue({
        instance: {},
      }),
    } as any

    // Mock MonoGuard global
    ;(window as any).MonoGuard = {
      getVersion: vi.fn().mockReturnValue(
        JSON.stringify({
          data: { version: '0.1.0' },
          error: null,
        })
      ),
      analyze: vi.fn().mockReturnValue(
        JSON.stringify({
          data: {
            healthScore: 85,
            packages: 5,
            circularDependencies: [],
            versionConflicts: [],
          },
          error: null,
        })
      ),
      check: vi.fn().mockReturnValue(
        JSON.stringify({
          data: { passed: true, errors: [] },
          error: null,
        })
      ),
    }
  })

  afterEach(() => {
    vi.restoreAllMocks()
    delete (window as any).MonoGuard
    delete (window as any).Go
  })

  describe('init', () => {
    it('should initialize successfully', async () => {
      await expect(analyzer.init()).resolves.not.toThrow()
      expect(analyzer.isInitialized()).toBe(true)
    })

    it('should be safe to call multiple times', async () => {
      await analyzer.init()
      await analyzer.init()
      expect(analyzer.isInitialized()).toBe(true)
    })
  })

  describe('getVersion', () => {
    it('should return version info when initialized', async () => {
      await analyzer.init()
      const result = await analyzer.getVersion()

      expect(isSuccess(result)).toBe(true)
      if (isSuccess(result)) {
        expect(result.data.version).toBe('0.1.0')
      }
    })

    it('should return error when not initialized', async () => {
      const result = await analyzer.getVersion()

      expect(isError(result)).toBe(true)
      if (isError(result)) {
        expect(result.error.code).toBe('NOT_INITIALIZED')
      }
    })
  })

  describe('analyze', () => {
    it('should analyze workspace when initialized', async () => {
      await analyzer.init()
      const result = await analyzer.analyze({
        files: { 'package.json': '{}' },
      })

      expect(isSuccess(result)).toBe(true)
      if (isSuccess(result)) {
        expect(result.data.healthScore).toBe(85)
        expect(result.data.packages).toBe(5)
      }
    })

    it('should pass config to WASM', async () => {
      await analyzer.init()
      await analyzer.analyze({
        files: { 'package.json': '{}' },
        config: { exclude: ['packages/legacy-*'] },
      })

      expect((window as any).MonoGuard.analyze).toHaveBeenCalledWith(
        expect.stringContaining('exclude')
      )
    })

    it('should return error when not initialized', async () => {
      const result = await analyzer.analyze({ files: {} })

      expect(isError(result)).toBe(true)
      if (isError(result)) {
        expect(result.error.code).toBe('NOT_INITIALIZED')
      }
    })
  })

  describe('check', () => {
    it('should check workspace when initialized', async () => {
      await analyzer.init()
      const result = await analyzer.check({
        files: { 'package.json': '{}' },
      })

      expect(isSuccess(result)).toBe(true)
      if (isSuccess(result)) {
        expect(result.data.passed).toBe(true)
        expect(result.data.errors).toEqual([])
      }
    })

    it('should return error when not initialized', async () => {
      const result = await analyzer.check({ files: {} })

      expect(isError(result)).toBe(true)
      if (isError(result)) {
        expect(result.error.code).toBe('NOT_INITIALIZED')
      }
    })
  })
})

describe('Error handling', () => {
  let analyzer: MonoGuardAnalyzer

  beforeEach(() => {
    analyzer = new MonoGuardAnalyzer()
    resetWasmState()
    ;(window as any).Go = class {
      importObject = {}
      async run() {}
    }
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      arrayBuffer: () => Promise.resolve(new ArrayBuffer(8)),
    })
    global.WebAssembly = {
      instantiate: vi.fn().mockResolvedValue({ instance: {} }),
    } as any
  })

  afterEach(() => {
    vi.restoreAllMocks()
    delete (window as any).MonoGuard
    delete (window as any).Go
  })

  it('should handle WASM function error response', async () => {
    ;(window as any).MonoGuard = {
      getVersion: vi.fn().mockReturnValue(
        JSON.stringify({
          data: null,
          error: { code: 'PARSE_ERROR', message: 'Invalid input' },
        })
      ),
      analyze: vi.fn(),
      check: vi.fn(),
    }

    await analyzer.init()
    const result = await analyzer.getVersion()

    expect(isError(result)).toBe(true)
    if (isError(result)) {
      expect(result.error.code).toBe('PARSE_ERROR')
      expect(result.error.message).toBe('Invalid input')
    }
  })

  it('should handle WASM function throwing', async () => {
    ;(window as any).MonoGuard = {
      getVersion: vi.fn().mockImplementation(() => {
        throw new Error('WASM crashed')
      }),
      analyze: vi.fn(),
      check: vi.fn(),
    }

    await analyzer.init()
    const result = await analyzer.getVersion()

    expect(isError(result)).toBe(true)
    if (isError(result)) {
      expect(result.error.code).toBe('WASM_CALL_FAILED')
    }
  })

  it('should handle invalid JSON from WASM', async () => {
    ;(window as any).MonoGuard = {
      getVersion: vi.fn().mockReturnValue('invalid json'),
      analyze: vi.fn(),
      check: vi.fn(),
    }

    await analyzer.init()
    const result = await analyzer.getVersion()

    expect(isError(result)).toBe(true)
    if (isError(result)) {
      expect(result.error.code).toBe('WASM_PARSE_ERROR')
    }
  })
})
