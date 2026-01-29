/**
 * Tests for useNodeHover hook
 *
 * @see Story 4.5: Implement Hover Details and Tooltips (AC2, AC4, AC5)
 *
 * Following Given-When-Then format with priority tags.
 */

import { act, renderHook } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import type { D3Link, D3Node } from '../types'
import { useNodeHover } from '../useNodeHover'

describe('useNodeHover', () => {
  // Mock data for testing
  const createMockNodes = (): D3Node[] => [
    { id: 'A', name: 'Package A', path: '/a', dependencyCount: 2, inCycle: false, cycleIds: [] },
    { id: 'B', name: 'Package B', path: '/b', dependencyCount: 1, inCycle: false, cycleIds: [] },
    { id: 'C', name: 'Package C', path: '/c', dependencyCount: 1, inCycle: false, cycleIds: [] },
    { id: 'D', name: 'Package D', path: '/d', dependencyCount: 0, inCycle: false, cycleIds: [] },
  ]

  const createMockLinks = (): D3Link[] => [
    { source: 'A', target: 'B', type: 'production', inCycle: false, cycleIds: [] },
    { source: 'A', target: 'C', type: 'production', inCycle: false, cycleIds: [] },
    { source: 'B', target: 'D', type: 'production', inCycle: false, cycleIds: [] },
  ]

  describe('Initial State (AC2)', () => {
    it('[P1] should initialize with null hover state', () => {
      // GIVEN: useNodeHover hook
      const { result } = renderHook(() =>
        useNodeHover({ nodes: createMockNodes(), links: createMockLinks() })
      )

      // THEN: Hover state should be null
      expect(result.current.hoverState.nodeId).toBeNull()
      expect(result.current.hoverState.position).toBeNull()
    })

    it('[P1] should initialize with empty connected sets', () => {
      // GIVEN: useNodeHover hook
      const { result } = renderHook(() =>
        useNodeHover({ nodes: createMockNodes(), links: createMockLinks() })
      )

      // THEN: Connected sets should be empty
      expect(result.current.connectedNodeIds.size).toBe(0)
    })
  })

  describe('Mouse Enter (AC2)', () => {
    it('[P1] should update hover state on mouse enter', () => {
      // GIVEN: useNodeHover hook
      const { result } = renderHook(() =>
        useNodeHover({ nodes: createMockNodes(), links: createMockLinks() })
      )

      // WHEN: Mouse enters a node
      act(() => {
        result.current.handleNodeMouseEnter('A', { clientX: 100, clientY: 200 } as MouseEvent)
      })

      // THEN: Hover state should update with node ID and position
      expect(result.current.hoverState.nodeId).toBe('A')
      expect(result.current.hoverState.position).toEqual({ x: 100, y: 200 })
    })

    it('[P1] should compute connected nodes correctly', () => {
      // GIVEN: useNodeHover hook with connected nodes A -> B, A -> C
      const { result } = renderHook(() =>
        useNodeHover({ nodes: createMockNodes(), links: createMockLinks() })
      )

      // WHEN: Hovering over node A
      act(() => {
        result.current.handleNodeMouseEnter('A', { clientX: 100, clientY: 200 } as MouseEvent)
      })

      // THEN: Connected nodes should include A, B, and C (A's direct connections)
      expect(result.current.connectedNodeIds.has('A')).toBe(true)
      expect(result.current.connectedNodeIds.has('B')).toBe(true)
      expect(result.current.connectedNodeIds.has('C')).toBe(true)
      // D is not directly connected to A
      expect(result.current.connectedNodeIds.has('D')).toBe(false)
    })
  })

  describe('Mouse Leave (AC2)', () => {
    it('[P1] should clear hover state on mouse leave', () => {
      // GIVEN: Hook with active hover state
      const { result } = renderHook(() =>
        useNodeHover({ nodes: createMockNodes(), links: createMockLinks() })
      )

      act(() => {
        result.current.handleNodeMouseEnter('A', { clientX: 100, clientY: 200 } as MouseEvent)
      })

      // WHEN: Mouse leaves the node
      act(() => {
        result.current.handleNodeMouseLeave()
      })

      // THEN: Hover state should be cleared
      expect(result.current.hoverState.nodeId).toBeNull()
      expect(result.current.hoverState.position).toBeNull()
    })

    it('[P1] should clear connected sets on mouse leave', () => {
      // GIVEN: Hook with active hover state
      const { result } = renderHook(() =>
        useNodeHover({ nodes: createMockNodes(), links: createMockLinks() })
      )

      act(() => {
        result.current.handleNodeMouseEnter('A', { clientX: 100, clientY: 200 } as MouseEvent)
      })

      // WHEN: Mouse leaves the node
      act(() => {
        result.current.handleNodeMouseLeave()
      })

      // THEN: Connected sets should be empty
      expect(result.current.connectedNodeIds.size).toBe(0)
    })
  })

  describe('Mouse Move (AC2)', () => {
    it('[P1] should update position on mouse move', () => {
      // CR2-3: handleNodeMouseMove uses rAF throttling - mock rAF to execute synchronously
      const originalRAF = globalThis.requestAnimationFrame
      globalThis.requestAnimationFrame = ((cb: (time: number) => void) => {
        cb(0)
        return 0
      }) as typeof requestAnimationFrame

      // GIVEN: Hook with active hover state
      const { result } = renderHook(() =>
        useNodeHover({ nodes: createMockNodes(), links: createMockLinks() })
      )

      act(() => {
        result.current.handleNodeMouseEnter('A', { clientX: 100, clientY: 200 } as MouseEvent)
      })

      // WHEN: Mouse moves within the node
      act(() => {
        result.current.handleNodeMouseMove({ clientX: 150, clientY: 250 } as MouseEvent)
      })

      // THEN: Position should update, node ID unchanged
      expect(result.current.hoverState.nodeId).toBe('A')
      expect(result.current.hoverState.position).toEqual({ x: 150, y: 250 })

      // Restore original rAF
      globalThis.requestAnimationFrame = originalRAF
    })
  })

  describe('Connected Elements Computation (AC4)', () => {
    it('[P1] should include incoming connections', () => {
      // GIVEN: Links where B points to D (B->D means D has incoming from B)
      const { result } = renderHook(() =>
        useNodeHover({ nodes: createMockNodes(), links: createMockLinks() })
      )

      // WHEN: Hovering over node D (which has incoming from B)
      act(() => {
        result.current.handleNodeMouseEnter('D', { clientX: 100, clientY: 200 } as MouseEvent)
      })

      // THEN: Should include D and B (incoming connection)
      expect(result.current.connectedNodeIds.has('D')).toBe(true)
      expect(result.current.connectedNodeIds.has('B')).toBe(true)
    })

    it('[P1] should include outgoing connections', () => {
      // GIVEN: Links where B points to D (B->D means B has outgoing to D)
      const { result } = renderHook(() =>
        useNodeHover({ nodes: createMockNodes(), links: createMockLinks() })
      )

      // WHEN: Hovering over node B (which has outgoing to D)
      act(() => {
        result.current.handleNodeMouseEnter('B', { clientX: 100, clientY: 200 } as MouseEvent)
      })

      // THEN: Should include B, D (outgoing), A (incoming), and relevant links
      expect(result.current.connectedNodeIds.has('B')).toBe(true)
      expect(result.current.connectedNodeIds.has('D')).toBe(true)
      expect(result.current.connectedNodeIds.has('A')).toBe(true)
    })

    it('[P2] should handle node with no connections', () => {
      // GIVEN: Isolated node not in links
      const nodes = [
        ...createMockNodes(),
        {
          id: 'isolated',
          name: 'Isolated',
          path: '/i',
          dependencyCount: 0,
          inCycle: false,
          cycleIds: [],
        },
      ]
      const { result } = renderHook(() => useNodeHover({ nodes, links: createMockLinks() }))

      // WHEN: Hovering over isolated node
      act(() => {
        result.current.handleNodeMouseEnter('isolated', {
          clientX: 100,
          clientY: 200,
        } as MouseEvent)
      })

      // THEN: Should only include the isolated node itself
      expect(result.current.connectedNodeIds.size).toBe(1)
      expect(result.current.connectedNodeIds.has('isolated')).toBe(true)
    })
  })

  describe('Links with D3Node Objects (AC4)', () => {
    it('[P1] should handle links with node objects instead of string IDs', () => {
      // GIVEN: Links with D3Node objects (as D3 replaces string IDs during simulation)
      const nodes = createMockNodes()
      const links: D3Link[] = [
        { source: nodes[0], target: nodes[1], type: 'production', inCycle: false, cycleIds: [] },
        { source: nodes[0], target: nodes[2], type: 'production', inCycle: false, cycleIds: [] },
      ]

      const { result } = renderHook(() => useNodeHover({ nodes, links }))

      // WHEN: Hovering over node A
      act(() => {
        result.current.handleNodeMouseEnter('A', { clientX: 100, clientY: 200 } as MouseEvent)
      })

      // THEN: Should correctly identify connected nodes
      expect(result.current.connectedNodeIds.has('A')).toBe(true)
      expect(result.current.connectedNodeIds.has('B')).toBe(true)
      expect(result.current.connectedNodeIds.has('C')).toBe(true)
    })
  })

  describe('Performance / Edge Cases (AC5)', () => {
    it('[P2] should handle switching between nodes quickly', () => {
      // GIVEN: useNodeHover hook
      const { result } = renderHook(() =>
        useNodeHover({ nodes: createMockNodes(), links: createMockLinks() })
      )

      // WHEN: Rapidly switching between nodes
      act(() => {
        result.current.handleNodeMouseEnter('A', { clientX: 100, clientY: 100 } as MouseEvent)
      })
      act(() => {
        result.current.handleNodeMouseLeave()
      })
      act(() => {
        result.current.handleNodeMouseEnter('B', { clientX: 200, clientY: 200 } as MouseEvent)
      })

      // THEN: Should correctly show B's state
      expect(result.current.hoverState.nodeId).toBe('B')
      expect(result.current.connectedNodeIds.has('B')).toBe(true)
      expect(result.current.connectedNodeIds.has('A')).toBe(true) // B is connected to A
      expect(result.current.connectedNodeIds.has('D')).toBe(true) // B -> D
    })

    it('[P2] should handle empty nodes array', () => {
      // GIVEN: Empty nodes
      const { result } = renderHook(() => useNodeHover({ nodes: [], links: [] }))

      // WHEN: Trying to hover
      act(() => {
        result.current.handleNodeMouseEnter('nonexistent', {
          clientX: 100,
          clientY: 100,
        } as MouseEvent)
      })

      // THEN: Should still work without errors, just empty connections
      expect(result.current.hoverState.nodeId).toBe('nonexistent')
      expect(result.current.connectedNodeIds.size).toBe(1) // Just the hovered node
    })

    it('[P3] should handle links data change', () => {
      // GIVEN: Hook with initial links
      const nodes = createMockNodes()
      const initialLinks = createMockLinks()
      const { result, rerender } = renderHook(({ links }) => useNodeHover({ nodes, links }), {
        initialProps: { links: initialLinks },
      })

      // WHEN: Hover, then links change
      act(() => {
        result.current.handleNodeMouseEnter('A', { clientX: 100, clientY: 100 } as MouseEvent)
      })

      // Rerender with new links (A only connects to B now)
      const newLinks: D3Link[] = [
        { source: 'A', target: 'B', type: 'production', inCycle: false, cycleIds: [] },
      ]
      rerender({ links: newLinks })

      // THEN: Connected sets should update based on new links
      expect(result.current.connectedNodeIds.has('A')).toBe(true)
      expect(result.current.connectedNodeIds.has('B')).toBe(true)
      expect(result.current.connectedNodeIds.has('C')).toBe(false) // No longer connected
    })
  })
})
