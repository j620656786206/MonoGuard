import '@testing-library/jest-dom/vitest'
import { cleanup } from '@testing-library/react'
import { afterEach, vi } from 'vitest'

// Cleanup after each test to prevent memory leaks
afterEach(() => {
  cleanup()
})

// Mock WASM loader for tests - returns proper Result<T> structure
vi.mock('@/lib/wasmLoader', () => ({
  loadWasm: vi.fn().mockResolvedValue({
    analyzer: {
      getVersion: vi.fn().mockReturnValue(
        JSON.stringify({
          data: { version: '0.1.0' },
          error: null,
        })
      ),
      analyze: vi.fn().mockReturnValue(
        JSON.stringify({
          data: { healthScore: 85, packageCount: 10 },
          error: null,
        })
      ),
      check: vi.fn().mockReturnValue(
        JSON.stringify({
          data: { passed: true, errors: [] },
          error: null,
        })
      ),
    },
  }),
}))

// Mock window.matchMedia for components that use media queries
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation((query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
})

// Mock ResizeObserver for components that use it
global.ResizeObserver = vi.fn().mockImplementation(() => ({
  observe: vi.fn(),
  unobserve: vi.fn(),
  disconnect: vi.fn(),
}))
