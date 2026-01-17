import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { isWasmInitialized, loadWasm, resetWasmState, WasmLoadError } from '../loader'

// biome-ignore lint/suspicious/noExplicitAny: Test mocking requires window global manipulation
type WindowWithGlobals = Window & { MonoGuard?: any; Go?: any }
const testWindow = window as WindowWithGlobals

// biome-ignore lint/suspicious/noExplicitAny: Mock WebAssembly requires any type
type GlobalWithWebAssembly = typeof globalThis & { WebAssembly: any }

describe('WASM Loader', () => {
  beforeEach(() => {
    resetWasmState()

    // Mock window.Go
    testWindow.Go = class {
      importObject = {}
      async run() {}
    }

    // Mock successful fetch
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      arrayBuffer: () => Promise.resolve(new ArrayBuffer(8)),
    })

    // Mock WebAssembly.instantiate
    ;(global as GlobalWithWebAssembly).WebAssembly = {
      instantiate: vi.fn().mockResolvedValue({
        instance: {},
      }),
    }

    // Mock MonoGuard global (set after WASM "loads")
    testWindow.MonoGuard = {
      getVersion: vi.fn(),
      analyze: vi.fn(),
      check: vi.fn(),
    }
  })

  afterEach(() => {
    vi.restoreAllMocks()
    resetWasmState()
    delete testWindow.MonoGuard
    delete testWindow.Go
  })

  describe('loadWasm', () => {
    it('should load WASM successfully', async () => {
      await expect(loadWasm()).resolves.not.toThrow()
      expect(isWasmInitialized()).toBe(true)
    })

    it('should use default wasm path', async () => {
      await loadWasm()
      expect(fetch).toHaveBeenCalledWith('/monoguard.wasm')
    })

    it('should use custom wasm path', async () => {
      await loadWasm({ wasmPath: '/custom/path.wasm' })
      expect(fetch).toHaveBeenCalledWith('/custom/path.wasm')
    })

    it('should be safe to call multiple times', async () => {
      await loadWasm()
      await loadWasm()
      expect(fetch).toHaveBeenCalledTimes(1)
    })

    it('should throw WasmLoadError if Go runtime not found', async () => {
      delete testWindow.Go

      await expect(loadWasm()).rejects.toThrow(WasmLoadError)
      await expect(loadWasm()).rejects.toThrow('Go runtime not found')
    })

    it('should throw WasmLoadError if fetch fails', async () => {
      global.fetch = vi.fn().mockResolvedValue({
        ok: false,
        status: 404,
        statusText: 'Not Found',
      })

      await expect(loadWasm()).rejects.toThrow(WasmLoadError)
      await expect(loadWasm()).rejects.toThrow('404')
    })

    it('should throw WasmLoadError on timeout', async () => {
      delete testWindow.MonoGuard
      testWindow.MonoGuard = undefined

      await expect(loadWasm({ timeout: 100 })).rejects.toThrow(WasmLoadError)
      await expect(loadWasm({ timeout: 100 })).rejects.toThrow('timed out')
    }, 10000)
  })

  describe('isWasmInitialized', () => {
    it('should return false before loading', () => {
      expect(isWasmInitialized()).toBe(false)
    })

    it('should return true after loading', async () => {
      await loadWasm()
      expect(isWasmInitialized()).toBe(true)
    })

    it('should return false if MonoGuard is undefined', async () => {
      await loadWasm()
      delete testWindow.MonoGuard
      expect(isWasmInitialized()).toBe(false)
    })
  })

  describe('resetWasmState', () => {
    it('should reset initialization state', async () => {
      await loadWasm()
      expect(isWasmInitialized()).toBe(true)

      resetWasmState()
      expect(isWasmInitialized()).toBe(false)
    })
  })
})

describe('WasmLoadError', () => {
  it('should have correct properties', () => {
    const error = new WasmLoadError('TEST_CODE', 'Test message')

    expect(error.name).toBe('WasmLoadError')
    expect(error.code).toBe('TEST_CODE')
    expect(error.message).toBe('Test message')
    expect(error instanceof Error).toBe(true)
  })
})
