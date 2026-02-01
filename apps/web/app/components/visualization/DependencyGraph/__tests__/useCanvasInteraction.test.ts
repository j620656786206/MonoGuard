/**
 * Tests for useCanvasInteraction hook (Story 4.9)
 *
 * Validates canvas hit detection, coordinate transformation,
 * hover/click callbacks, and cursor management.
 *
 * @see AC4: Canvas Mode Hover/Click Functionality
 */

import { act, renderHook } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import type { D3Node, ViewportState } from '../types'
import { useCanvasInteraction } from '../useCanvasInteraction'

// Helper: create a mock D3Node at a given position
function createNode(id: string, x: number, y: number): D3Node {
  return {
    id,
    name: `@app/${id}`,
    path: `packages/${id}`,
    dependencyCount: 1,
    inCycle: false,
    cycleIds: [],
    x,
    y,
  }
}

// Helper: create a minimal React.MouseEvent-like object for canvas
function createMouseEvent(
  clientX: number,
  clientY: number,
  rectOverride?: Partial<DOMRect>
): React.MouseEvent<HTMLCanvasElement> {
  return {
    clientX,
    clientY,
  } as React.MouseEvent<HTMLCanvasElement>
}

describe('useCanvasInteraction', () => {
  const mockOnNodeHover = vi.fn()
  const mockOnNodeSelect = vi.fn()

  let mockCanvas: HTMLCanvasElement
  let canvasRef: { current: HTMLCanvasElement | null }
  let nodesRef: { current: D3Node[] }

  beforeEach(() => {
    mockCanvas = document.createElement('canvas')
    // Set bounding rect for coordinate calculations
    vi.spyOn(mockCanvas, 'getBoundingClientRect').mockReturnValue({
      left: 0,
      top: 0,
      right: 800,
      bottom: 600,
      width: 800,
      height: 600,
      x: 0,
      y: 0,
      toJSON: vi.fn(),
    })

    canvasRef = { current: mockCanvas }
    nodesRef = { current: [] }
  })

  afterEach(() => {
    vi.clearAllMocks()
  })

  const defaultViewport: ViewportState = { zoom: 1, panX: 0, panY: 0 }

  function renderInteractionHook(viewport: ViewportState = defaultViewport, nodes: D3Node[] = []) {
    nodesRef.current = nodes
    return renderHook(() =>
      useCanvasInteraction({
        canvasRef,
        nodesRef,
        viewport,
        onNodeHover: mockOnNodeHover,
        onNodeSelect: mockOnNodeSelect,
      })
    )
  }

  describe('[P1] coordinate transformation', () => {
    it('should transform screen coordinates to graph space with default viewport', () => {
      // GIVEN: A node at graph position (100, 100), viewport at 1x zoom no pan
      const nodes = [createNode('a', 100, 100)]
      const { result } = renderInteractionHook(defaultViewport, nodes)

      // WHEN: Mouse moves exactly over node at screen (100, 100)
      act(() => {
        result.current.handleMouseMove(createMouseEvent(100, 100))
      })

      // THEN: Node should be detected and hover callback called with the node
      expect(mockOnNodeHover).toHaveBeenCalledWith(
        expect.objectContaining({ id: 'a' }),
        expect.objectContaining({ x: 100, y: 100 })
      )
    })

    it('should account for viewport pan offset in coordinate transformation', () => {
      // GIVEN: A node at graph position (200, 200), viewport panned by (50, 50)
      const nodes = [createNode('a', 200, 200)]
      const viewport: ViewportState = { zoom: 1, panX: 50, panY: 50 }
      const { result } = renderInteractionHook(viewport, nodes)

      // WHEN: Mouse at screen (250, 250) → graph space should be (200, 200)
      act(() => {
        result.current.handleMouseMove(createMouseEvent(250, 250))
      })

      // THEN: Node detected because (250 - 50) / 1 = 200 matches node position
      expect(mockOnNodeHover).toHaveBeenCalledWith(
        expect.objectContaining({ id: 'a' }),
        expect.objectContaining({ x: 250, y: 250 })
      )
    })

    it('should account for viewport zoom in coordinate transformation', () => {
      // GIVEN: A node at graph position (100, 100), viewport at 2x zoom
      const nodes = [createNode('a', 100, 100)]
      const viewport: ViewportState = { zoom: 2, panX: 0, panY: 0 }
      const { result } = renderInteractionHook(viewport, nodes)

      // WHEN: Mouse at screen (200, 200) → graph space = (200/2, 200/2) = (100, 100)
      act(() => {
        result.current.handleMouseMove(createMouseEvent(200, 200))
      })

      // THEN: Node detected because graph coords match node position
      expect(mockOnNodeHover).toHaveBeenCalledWith(
        expect.objectContaining({ id: 'a' }),
        expect.objectContaining({ x: 200, y: 200 })
      )
    })

    it('should handle combined zoom and pan', () => {
      // GIVEN: Node at (50, 50), viewport zoomed 2x and panned (100, 100)
      const nodes = [createNode('a', 50, 50)]
      const viewport: ViewportState = { zoom: 2, panX: 100, panY: 100 }
      const { result } = renderInteractionHook(viewport, nodes)

      // WHEN: Mouse at screen (200, 200) → graph = (200 - 100) / 2 = 50
      act(() => {
        result.current.handleMouseMove(createMouseEvent(200, 200))
      })

      // THEN: Node detected
      expect(mockOnNodeHover).toHaveBeenCalledWith(
        expect.objectContaining({ id: 'a' }),
        expect.any(Object)
      )
    })
  })

  describe('[P1] hit detection', () => {
    it('should detect node when mouse is within hit radius (12px)', () => {
      // GIVEN: Node at (100, 100)
      const nodes = [createNode('a', 100, 100)]
      const { result } = renderInteractionHook(defaultViewport, nodes)

      // WHEN: Mouse at (110, 100) → distance = 10 < 12
      act(() => {
        result.current.handleMouseMove(createMouseEvent(110, 100))
      })

      // THEN: Node detected
      expect(mockOnNodeHover).toHaveBeenCalledWith(
        expect.objectContaining({ id: 'a' }),
        expect.any(Object)
      )
    })

    it('should not detect node when mouse is outside hit radius', () => {
      // GIVEN: Node at (100, 100)
      const nodes = [createNode('a', 100, 100)]
      const { result } = renderInteractionHook(defaultViewport, nodes)

      // WHEN: Mouse at (113, 100) → distance = 13 > 12
      act(() => {
        result.current.handleMouseMove(createMouseEvent(113, 100))
      })

      // THEN: No node detected, hover cleared
      expect(mockOnNodeHover).toHaveBeenCalledWith(null, null)
    })

    it('should detect the topmost (last rendered) node when overlapping', () => {
      // GIVEN: Two overlapping nodes at the same position
      const nodes = [createNode('bottom', 100, 100), createNode('top', 100, 100)]
      const { result } = renderInteractionHook(defaultViewport, nodes)

      // WHEN: Mouse over the overlapping position
      act(() => {
        result.current.handleMouseMove(createMouseEvent(100, 100))
      })

      // THEN: The last node in array (top-rendered) should be detected
      expect(mockOnNodeHover).toHaveBeenCalledWith(
        expect.objectContaining({ id: 'top' }),
        expect.any(Object)
      )
    })

    it('should skip nodes with undefined coordinates', () => {
      // GIVEN: Node with no position data
      const nodeWithNoCoords: D3Node = {
        id: 'no-pos',
        name: '@app/no-pos',
        path: 'packages/no-pos',
        dependencyCount: 0,
        inCycle: false,
        cycleIds: [],
        x: undefined,
        y: undefined,
      }
      const nodes = [nodeWithNoCoords]
      const { result } = renderInteractionHook(defaultViewport, nodes)

      // WHEN: Mouse at (0, 0) where the undefined-position node might be
      act(() => {
        result.current.handleMouseMove(createMouseEvent(0, 0))
      })

      // THEN: No node detected
      expect(mockOnNodeHover).toHaveBeenCalledWith(null, null)
    })

    it('should handle empty nodes list', () => {
      // GIVEN: No nodes
      const { result } = renderInteractionHook(defaultViewport, [])

      // WHEN: Mouse moves
      act(() => {
        result.current.handleMouseMove(createMouseEvent(100, 100))
      })

      // THEN: No node detected
      expect(mockOnNodeHover).toHaveBeenCalledWith(null, null)
    })
  })

  describe('[P1] hover callbacks', () => {
    it('should call onNodeHover with node and screen position when hovering', () => {
      // GIVEN: Node at (100, 100)
      const nodes = [createNode('a', 100, 100)]
      const { result } = renderInteractionHook(defaultViewport, nodes)

      // WHEN: Mouse over the node
      act(() => {
        result.current.handleMouseMove(createMouseEvent(100, 100))
      })

      // THEN: Callback with node object and screen coordinates
      expect(mockOnNodeHover).toHaveBeenCalledTimes(1)
      expect(mockOnNodeHover).toHaveBeenCalledWith(
        expect.objectContaining({ id: 'a', name: '@app/a' }),
        { x: 100, y: 100 }
      )
    })

    it('should call onNodeHover with null when leaving a node', () => {
      // GIVEN: Node at (100, 100)
      const nodes = [createNode('a', 100, 100)]
      const { result } = renderInteractionHook(defaultViewport, nodes)

      // WHEN: Mouse moves to empty area
      act(() => {
        result.current.handleMouseMove(createMouseEvent(300, 300))
      })

      // THEN: Hover cleared
      expect(mockOnNodeHover).toHaveBeenCalledWith(null, null)
    })
  })

  describe('[P1] click callbacks', () => {
    it('should call onNodeSelect with node id when clicking on a node', () => {
      // GIVEN: Node at (100, 100)
      const nodes = [createNode('a', 100, 100)]
      const { result } = renderInteractionHook(defaultViewport, nodes)

      // WHEN: Click on the node
      act(() => {
        result.current.handleMouseClick(createMouseEvent(100, 100))
      })

      // THEN: Selection callback with node id
      expect(mockOnNodeSelect).toHaveBeenCalledWith('a')
    })

    it('should call onNodeSelect with null when clicking empty space', () => {
      // GIVEN: Node at (100, 100)
      const nodes = [createNode('a', 100, 100)]
      const { result } = renderInteractionHook(defaultViewport, nodes)

      // WHEN: Click on empty area
      act(() => {
        result.current.handleMouseClick(createMouseEvent(300, 300))
      })

      // THEN: Deselect
      expect(mockOnNodeSelect).toHaveBeenCalledWith(null)
    })
  })

  describe('[P1] cursor management', () => {
    it('should set cursor to pointer when hovering over a node', () => {
      // GIVEN: Node at (100, 100)
      const nodes = [createNode('a', 100, 100)]
      const { result } = renderInteractionHook(defaultViewport, nodes)

      // WHEN: Mouse over node
      act(() => {
        result.current.handleMouseMove(createMouseEvent(100, 100))
      })

      // THEN: Cursor changed to pointer
      expect(mockCanvas.style.cursor).toBe('pointer')
    })

    it('should set cursor to crosshair when not hovering over a node', () => {
      // GIVEN: Node at (100, 100)
      const nodes = [createNode('a', 100, 100)]
      const { result } = renderInteractionHook(defaultViewport, nodes)

      // WHEN: Mouse over empty area
      act(() => {
        result.current.handleMouseMove(createMouseEvent(300, 300))
      })

      // THEN: Cursor reset to crosshair
      expect(mockCanvas.style.cursor).toBe('crosshair')
    })
  })

  describe('[P2] edge cases', () => {
    it('should handle null canvas ref gracefully for mouse move', () => {
      // GIVEN: canvasRef is null
      canvasRef.current = null
      const { result } = renderInteractionHook(defaultViewport, [createNode('a', 100, 100)])

      // WHEN: Mouse move (should not throw)
      act(() => {
        result.current.handleMouseMove(createMouseEvent(100, 100))
      })

      // THEN: No callbacks called (early return)
      expect(mockOnNodeHover).not.toHaveBeenCalled()
    })

    it('should handle null canvas ref gracefully for click', () => {
      // GIVEN: canvasRef is null
      canvasRef.current = null
      const { result } = renderInteractionHook(defaultViewport, [createNode('a', 100, 100)])

      // WHEN: Click (should not throw)
      act(() => {
        result.current.handleMouseClick(createMouseEvent(100, 100))
      })

      // THEN: No callbacks called
      expect(mockOnNodeSelect).not.toHaveBeenCalled()
    })

    it('should find correct node at boundary distance', () => {
      // GIVEN: Node at (100, 100), hit radius is exactly 12
      const nodes = [createNode('a', 100, 100)]
      const { result } = renderInteractionHook(defaultViewport, nodes)

      // WHEN: Mouse at exactly 12px distance (112, 100)
      act(() => {
        result.current.handleMouseMove(createMouseEvent(112, 100))
      })

      // THEN: Node detected (distance == hitRadius, <= check)
      expect(mockOnNodeHover).toHaveBeenCalledWith(
        expect.objectContaining({ id: 'a' }),
        expect.any(Object)
      )
    })

    it('should handle canvas with non-zero bounding rect offset', () => {
      // GIVEN: Canvas offset by (50, 30) in the page
      vi.spyOn(mockCanvas, 'getBoundingClientRect').mockReturnValue({
        left: 50,
        top: 30,
        right: 850,
        bottom: 630,
        width: 800,
        height: 600,
        x: 50,
        y: 30,
        toJSON: vi.fn(),
      })

      const nodes = [createNode('a', 100, 100)]
      const { result } = renderInteractionHook(defaultViewport, nodes)

      // WHEN: Mouse at screen (150, 130) → canvas local = (100, 100) → graph = (100, 100)
      act(() => {
        result.current.handleMouseMove(createMouseEvent(150, 130))
      })

      // THEN: Node detected
      expect(mockOnNodeHover).toHaveBeenCalledWith(
        expect.objectContaining({ id: 'a' }),
        expect.objectContaining({ x: 100, y: 100 })
      )
    })
  })
})
