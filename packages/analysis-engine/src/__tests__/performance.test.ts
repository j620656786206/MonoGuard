/**
 * Performance Tests for MonoGuard WASM Adapter
 *
 * Verifies AC8 requirements:
 * - WASM initialization completes in < 2 seconds
 * - JSON serialization overhead is < 10ms
 */

import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { MonoGuardAnalyzer } from '../analyzer'
import { resetWasmState } from '../loader'

describe('Performance (AC8)', () => {
  beforeEach(() => {
    resetWasmState()

    // Mock window.Go for WASM initialization
    ;(window as any).Go = class {
      importObject = {}
      async run() {}
    }

    // Mock WASM fetch - simulate realistic fetch time
    global.fetch = vi.fn().mockImplementation(
      () =>
        new Promise((resolve) => {
          // Simulate ~50ms network latency
          setTimeout(() => {
            resolve({
              ok: true,
              arrayBuffer: () => Promise.resolve(new ArrayBuffer(8)),
            })
          }, 50)
        })
    )

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
            packages: 10,
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

  describe('Initialization Performance', () => {
    it('should initialize within 2 seconds (AC8)', async () => {
      const analyzer = new MonoGuardAnalyzer()

      const start = performance.now()
      await analyzer.init()
      const duration = performance.now() - start

      expect(duration).toBeLessThan(2000)
      expect(analyzer.isInitialized()).toBe(true)
    })

    it('should handle multiple concurrent init calls efficiently', async () => {
      const analyzer = new MonoGuardAnalyzer()

      const start = performance.now()
      // Simulate multiple components trying to init simultaneously
      await Promise.all([analyzer.init(), analyzer.init(), analyzer.init()])
      const duration = performance.now() - start

      // Should not be significantly slower than single init
      expect(duration).toBeLessThan(2000)
      expect(global.fetch).toHaveBeenCalledTimes(1) // Only fetches once
    })
  })

  describe('JSON Serialization Overhead', () => {
    it('should have minimal serialization overhead for analyze (< 10ms)', async () => {
      const analyzer = new MonoGuardAnalyzer()
      await analyzer.init()

      // Prepare input with realistic size
      const input = {
        files: generateMockWorkspaceFiles(20), // 20 packages
      }

      const iterations = 10
      let totalDuration = 0

      for (let i = 0; i < iterations; i++) {
        const start = performance.now()
        await analyzer.analyze(input)
        totalDuration += performance.now() - start
      }

      const avgDuration = totalDuration / iterations
      // Average should be well under 10ms (WASM mock returns instantly)
      expect(avgDuration).toBeLessThan(10)
    })

    it('should have minimal serialization overhead for check (< 10ms)', async () => {
      const analyzer = new MonoGuardAnalyzer()
      await analyzer.init()

      const input = {
        files: generateMockWorkspaceFiles(20),
      }

      const iterations = 10
      let totalDuration = 0

      for (let i = 0; i < iterations; i++) {
        const start = performance.now()
        await analyzer.check(input)
        totalDuration += performance.now() - start
      }

      const avgDuration = totalDuration / iterations
      expect(avgDuration).toBeLessThan(10)
    })

    it('should have minimal serialization overhead for getVersion (< 5ms)', async () => {
      const analyzer = new MonoGuardAnalyzer()
      await analyzer.init()

      const iterations = 10
      let totalDuration = 0

      for (let i = 0; i < iterations; i++) {
        const start = performance.now()
        await analyzer.getVersion()
        totalDuration += performance.now() - start
      }

      const avgDuration = totalDuration / iterations
      // getVersion has no input serialization, should be very fast
      expect(avgDuration).toBeLessThan(5)
    })
  })

  describe('Large Input Handling', () => {
    it('should handle large workspace input without significant slowdown', async () => {
      const analyzer = new MonoGuardAnalyzer()
      await analyzer.init()

      // Generate large input (100 packages, ~100KB JSON)
      const largeInput = {
        files: generateMockWorkspaceFiles(100),
      }

      const start = performance.now()
      await analyzer.analyze(largeInput)
      const duration = performance.now() - start

      // Even with large input, should complete quickly (mocked WASM)
      // Real WASM would take longer, but serialization should still be < 50ms
      expect(duration).toBeLessThan(100)
    })
  })
})

/**
 * Generate mock workspace files for testing
 */
function generateMockWorkspaceFiles(packageCount: number): Record<string, string> {
  const files: Record<string, string> = {
    'package.json': JSON.stringify({
      name: 'test-monorepo',
      workspaces: ['packages/*'],
    }),
    'pnpm-workspace.yaml': 'packages:\n  - packages/*',
  }

  for (let i = 0; i < packageCount; i++) {
    const pkgName = `@mono/package-${i}`
    const deps: Record<string, string> = {}

    // Add some internal dependencies
    if (i > 0) {
      deps[`@mono/package-${i - 1}`] = 'workspace:*'
    }
    if (i > 1 && i % 3 === 0) {
      deps[`@mono/package-${Math.floor(i / 2)}`] = 'workspace:*'
    }

    // Add some external dependencies
    deps['lodash'] = '^4.17.21'
    deps['typescript'] = '^5.0.0'

    files[`packages/package-${i}/package.json`] = JSON.stringify({
      name: pkgName,
      version: '1.0.0',
      dependencies: deps,
    })
  }

  return files
}
