/**
 * WASM Loader Module for MonoGuard Analysis Engine
 *
 * Handles loading and initializing the Go WASM module with timeout support.
 *
 * @module loader
 */

import type { WasmLoaderOptions } from '@monoguard/types'

/**
 * Error thrown when WASM loading fails
 */
export class WasmLoadError extends Error {
  constructor(
    public code: string,
    message: string
  ) {
    super(message)
    this.name = 'WasmLoadError'
  }
}

/**
 * Global Go instance for WASM runtime
 */
declare global {
  interface Window {
    Go: new () => GoInstance
    MonoGuard?: MonoGuardGlobal
  }
}

interface GoInstance {
  importObject: WebAssembly.Imports
  run(instance: WebAssembly.Instance): Promise<void>
}

interface MonoGuardGlobal {
  getVersion(): string
  analyze(input: string): string
  check(input: string): string
}

let wasmInitialized = false
let initPromise: Promise<void> | null = null

/**
 * Load and initialize the MonoGuard WASM module
 *
 * This function handles:
 * - Loading wasm_exec.js (Go runtime)
 * - Fetching the WASM binary
 * - Initializing the Go runtime
 * - Waiting for MonoGuard global to be available
 *
 * @param options - Configuration options for loading
 * @returns Promise that resolves when WASM is ready
 * @throws {WasmLoadError} If loading fails or times out
 *
 * @example
 * ```typescript
 * import { loadWasm } from '@monoguard/analysis-engine';
 *
 * await loadWasm({ wasmPath: '/wasm/monoguard.wasm' });
 * // MonoGuard global is now available
 * ```
 */
export async function loadWasm(options?: WasmLoaderOptions): Promise<void> {
  // Return existing initialization if in progress
  if (initPromise) {
    return initPromise
  }

  // Return immediately if already initialized
  if (wasmInitialized && typeof window.MonoGuard !== 'undefined') {
    return Promise.resolve()
  }

  const { wasmPath = '/monoguard.wasm', timeout = 10000 } = options ?? {}

  initPromise = new Promise<void>((resolve, reject) => {
    const timeoutId = setTimeout(() => {
      initPromise = null
      reject(new WasmLoadError('WASM_TIMEOUT', `WASM initialization timed out after ${timeout}ms`))
    }, timeout)

    loadWasmInternal(wasmPath)
      .then(() => {
        clearTimeout(timeoutId)
        wasmInitialized = true
        resolve()
      })
      .catch((err) => {
        clearTimeout(timeoutId)
        initPromise = null
        reject(err)
      })
  })

  return initPromise
}

/**
 * Internal WASM loading implementation
 */
async function loadWasmInternal(wasmPath: string): Promise<void> {
  // Check if Go runtime is available
  if (typeof window.Go === 'undefined') {
    throw new WasmLoadError(
      'WASM_LOAD_ERROR',
      'Go runtime not found. Make sure wasm_exec.js is loaded before calling loadWasm()'
    )
  }

  const go = new window.Go()

  try {
    // Fetch WASM binary
    const response = await fetch(wasmPath)
    if (!response.ok) {
      throw new WasmLoadError(
        'WASM_LOAD_ERROR',
        `Failed to fetch WASM: ${response.status} ${response.statusText}`
      )
    }

    // Instantiate WASM
    const wasmBytes = await response.arrayBuffer()
    const result = await WebAssembly.instantiate(wasmBytes, go.importObject)

    // Run Go runtime (non-blocking)
    go.run(result.instance)

    // Wait for MonoGuard global to be available
    await waitForMonoGuard()
  } catch (err) {
    if (err instanceof WasmLoadError) {
      throw err
    }
    throw new WasmLoadError(
      'WASM_LOAD_ERROR',
      err instanceof Error ? err.message : 'Unknown error loading WASM'
    )
  }
}

/**
 * Wait for MonoGuard global to be available
 */
function waitForMonoGuard(maxAttempts = 100, interval = 10): Promise<void> {
  return new Promise((resolve, reject) => {
    let attempts = 0

    const check = () => {
      if (typeof window.MonoGuard !== 'undefined') {
        resolve()
        return
      }

      attempts++
      if (attempts >= maxAttempts) {
        reject(
          new WasmLoadError(
            'WASM_LOAD_ERROR',
            'MonoGuard global not available after WASM initialization'
          )
        )
        return
      }

      setTimeout(check, interval)
    }

    check()
  })
}

/**
 * Check if WASM is initialized
 *
 * @returns true if WASM is ready to use
 */
export function isWasmInitialized(): boolean {
  return wasmInitialized && typeof window.MonoGuard !== 'undefined'
}

/**
 * Reset WASM initialization state (for testing)
 * @internal
 */
export function resetWasmState(): void {
  wasmInitialized = false
  initPromise = null
}
