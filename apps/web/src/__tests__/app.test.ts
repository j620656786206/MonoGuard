import { describe, expect, it } from 'vitest'

describe('App Configuration', () => {
  it('should have correct environment setup', () => {
    // Basic test to verify Vitest is working
    expect(typeof window).toBe('object')
  })

  it('should have jsdom environment', () => {
    expect(document).toBeDefined()
    expect(document.createElement).toBeInstanceOf(Function)
  })
})
