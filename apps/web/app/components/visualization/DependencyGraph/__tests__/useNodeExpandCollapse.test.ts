/**
 * Tests for useNodeExpandCollapse hook
 *
 * @see Story 4.3: Implement Node Expand/Collapse Functionality
 */

import { act, renderHook } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { useNodeExpandCollapse } from '../useNodeExpandCollapse'

describe('useNodeExpandCollapse', () => {
  const mockNodeIds = ['root', 'child1', 'child2', 'grandchild']
  const mockNodeDepths = new Map([
    ['root', 0],
    ['child1', 1],
    ['child2', 1],
    ['grandchild', 2],
  ])

  // Mock sessionStorage
  const mockSessionStorage = (() => {
    let store: Record<string, string> = {}
    return {
      getItem: vi.fn((key: string) => store[key] || null),
      setItem: vi.fn((key: string, value: string) => {
        store[key] = value
      }),
      removeItem: vi.fn((key: string) => {
        delete store[key]
      }),
      clear: vi.fn(() => {
        store = {}
      }),
    }
  })()

  beforeEach(() => {
    vi.stubGlobal('sessionStorage', mockSessionStorage)
    mockSessionStorage.clear()
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('should start with all nodes expanded', () => {
    const { result } = renderHook(() =>
      useNodeExpandCollapse({
        nodeIds: mockNodeIds,
        nodeDepths: mockNodeDepths,
      })
    )

    expect(result.current.collapsedNodeIds.size).toBe(0)
    expect(result.current.isCollapsed('root')).toBe(false)
    expect(result.current.isCollapsed('child1')).toBe(false)
  })

  it('should toggle node collapse state', () => {
    const { result } = renderHook(() =>
      useNodeExpandCollapse({
        nodeIds: mockNodeIds,
        nodeDepths: mockNodeDepths,
      })
    )

    act(() => {
      result.current.toggleNode('child1')
    })

    expect(result.current.isCollapsed('child1')).toBe(true)

    act(() => {
      result.current.toggleNode('child1')
    })

    expect(result.current.isCollapsed('child1')).toBe(false)
  })

  it('should collapse a specific node', () => {
    const { result } = renderHook(() =>
      useNodeExpandCollapse({
        nodeIds: mockNodeIds,
        nodeDepths: mockNodeDepths,
      })
    )

    act(() => {
      result.current.collapseNode('child1')
    })

    expect(result.current.isCollapsed('child1')).toBe(true)
    expect(result.current.isCollapsed('child2')).toBe(false)
  })

  it('should expand a specific node', () => {
    const { result } = renderHook(() =>
      useNodeExpandCollapse({
        nodeIds: mockNodeIds,
        nodeDepths: mockNodeDepths,
      })
    )

    // First collapse
    act(() => {
      result.current.collapseNode('child1')
    })

    expect(result.current.isCollapsed('child1')).toBe(true)

    // Then expand
    act(() => {
      result.current.expandNode('child1')
    })

    expect(result.current.isCollapsed('child1')).toBe(false)
  })

  it('should collapse all nodes at specified depth', () => {
    const { result } = renderHook(() =>
      useNodeExpandCollapse({
        nodeIds: mockNodeIds,
        nodeDepths: mockNodeDepths,
      })
    )

    act(() => {
      result.current.collapseAtDepth(1)
    })

    expect(result.current.isCollapsed('root')).toBe(false)
    expect(result.current.isCollapsed('child1')).toBe(true)
    expect(result.current.isCollapsed('child2')).toBe(true)
    expect(result.current.isCollapsed('grandchild')).toBe(true)
  })

  it('should expand all nodes to specified depth', () => {
    const { result } = renderHook(() =>
      useNodeExpandCollapse({
        nodeIds: mockNodeIds,
        nodeDepths: mockNodeDepths,
      })
    )

    // First collapse all
    act(() => {
      result.current.collapseAll()
    })

    // Then expand to depth 1
    act(() => {
      result.current.expandToDepth(1)
    })

    expect(result.current.isCollapsed('root')).toBe(false)
    expect(result.current.isCollapsed('child1')).toBe(false)
    expect(result.current.isCollapsed('child2')).toBe(false)
    expect(result.current.isCollapsed('grandchild')).toBe(true)
  })

  it('should expand all nodes', () => {
    const { result } = renderHook(() =>
      useNodeExpandCollapse({
        nodeIds: mockNodeIds,
        nodeDepths: mockNodeDepths,
      })
    )

    act(() => {
      result.current.collapseAll()
    })

    expect(result.current.collapsedNodeIds.size).toBeGreaterThan(0)

    act(() => {
      result.current.expandAll()
    })

    expect(result.current.collapsedNodeIds.size).toBe(0)
  })

  it('should collapse all nodes except root', () => {
    const { result } = renderHook(() =>
      useNodeExpandCollapse({
        nodeIds: mockNodeIds,
        nodeDepths: mockNodeDepths,
      })
    )

    act(() => {
      result.current.collapseAll()
    })

    // Root (depth 0) should never be collapsed
    expect(result.current.isCollapsed('root')).toBe(false)
    expect(result.current.isCollapsed('child1')).toBe(true)
    expect(result.current.isCollapsed('child2')).toBe(true)
    expect(result.current.isCollapsed('grandchild')).toBe(true)
  })

  describe('session storage persistence', () => {
    it('should persist collapsed state to session storage', async () => {
      const sessionKey = 'test-persist'
      mockSessionStorage.clear()

      const { result } = renderHook(() =>
        useNodeExpandCollapse({
          nodeIds: mockNodeIds,
          nodeDepths: mockNodeDepths,
          sessionKey,
        })
      )

      act(() => {
        result.current.collapseNode('child1')
      })

      // Wait for useEffect to run
      await vi.waitFor(() => {
        const calls = mockSessionStorage.setItem.mock.calls.filter(
          (call) => call[0] === `monoguard-collapse-${sessionKey}`
        )
        expect(calls.length).toBeGreaterThan(0)
      })

      const storedValue = mockSessionStorage.setItem.mock.calls
        .filter((call) => call[0] === `monoguard-collapse-${sessionKey}`)
        .pop()?.[1]

      expect(JSON.parse(storedValue as string)).toContain('child1')
    })

    it('should restore collapsed state from session storage', () => {
      const sessionKey = 'test-restore'
      const storedState = JSON.stringify(['child1', 'child2'])

      mockSessionStorage.getItem.mockImplementation((key: string) => {
        if (key === `monoguard-collapse-${sessionKey}`) {
          return storedState
        }
        return null
      })

      const { result } = renderHook(() =>
        useNodeExpandCollapse({
          nodeIds: mockNodeIds,
          nodeDepths: mockNodeDepths,
          sessionKey,
        })
      )

      expect(result.current.isCollapsed('child1')).toBe(true)
      expect(result.current.isCollapsed('child2')).toBe(true)
      expect(result.current.isCollapsed('root')).toBe(false)
    })

    it('should not use session storage when sessionKey is not provided', () => {
      mockSessionStorage.clear()
      mockSessionStorage.setItem.mockClear()

      const { result } = renderHook(() =>
        useNodeExpandCollapse({
          nodeIds: mockNodeIds,
          nodeDepths: mockNodeDepths,
          // No sessionKey provided
        })
      )

      act(() => {
        result.current.collapseNode('child1')
      })

      // Should not call sessionStorage.setItem with monoguard-collapse prefix when no sessionKey
      const relevantCalls = mockSessionStorage.setItem.mock.calls.filter(
        (call) => typeof call[0] === 'string' && call[0].startsWith('monoguard-collapse-')
      )
      expect(relevantCalls.length).toBe(0)
    })
  })

  describe('edge cases', () => {
    it('should handle empty nodeIds array', () => {
      const { result } = renderHook(() =>
        useNodeExpandCollapse({
          nodeIds: [],
          nodeDepths: new Map(),
        })
      )

      expect(result.current.collapsedNodeIds.size).toBe(0)

      act(() => {
        result.current.collapseAll()
      })

      expect(result.current.collapsedNodeIds.size).toBe(0)
    })

    it('should handle nodes without depth information', () => {
      const incompleteDepths = new Map([['root', 0]])

      const { result } = renderHook(() =>
        useNodeExpandCollapse({
          nodeIds: mockNodeIds,
          nodeDepths: incompleteDepths,
        })
      )

      // collapseAtDepth should treat unknown depths as 0
      act(() => {
        result.current.collapseAtDepth(1)
      })

      expect(result.current.isCollapsed('root')).toBe(false)
    })

    it('should handle toggle on non-existent node gracefully', () => {
      const { result } = renderHook(() =>
        useNodeExpandCollapse({
          nodeIds: mockNodeIds,
          nodeDepths: mockNodeDepths,
        })
      )

      // Should not throw
      act(() => {
        result.current.toggleNode('nonexistent')
      })

      // The non-existent node should still be tracked as collapsed
      expect(result.current.isCollapsed('nonexistent')).toBe(true)
    })
  })
})
