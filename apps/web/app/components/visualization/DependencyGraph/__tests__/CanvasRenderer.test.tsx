/**
 * Tests for CanvasRenderer component
 *
 * @see Story 4.9: Implement Hybrid SVG/Canvas Rendering
 * @see AC1: Automatic Render Mode Selection
 * @see AC6: Performance Metrics
 * @see AC7: Graceful Degradation
 */

import { render } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import type { D3Link, D3Node, ViewportState } from '../types'

// Mock d3 to avoid JSDOM issues with force simulation
vi.mock('d3', () => {
  const mockSimulation = {
    force: vi.fn().mockReturnThis(),
    on: vi.fn().mockReturnThis(),
    stop: vi.fn(),
    alphaDecay: vi.fn().mockReturnThis(),
  }

  const mockForceLink = {
    id: vi.fn().mockReturnThis(),
    distance: vi.fn().mockReturnThis(),
  }

  return {
    forceSimulation: vi.fn(() => mockSimulation),
    forceLink: vi.fn(() => mockForceLink),
    forceManyBody: vi.fn(() => ({ strength: vi.fn().mockReturnThis() })),
    forceCenter: vi.fn(),
    forceCollide: vi.fn(() => ({ radius: vi.fn().mockReturnThis() })),
  }
})

// Must import after mock
import { CanvasRenderer } from '../CanvasRenderer'

describe('CanvasRenderer', () => {
  const defaultViewport: ViewportState = { zoom: 1, panX: 0, panY: 0 }
  const mockOnViewportChange = vi.fn()
  const mockOnNodeSelect = vi.fn()
  const mockOnNodeHover = vi.fn()

  const mockNodes: D3Node[] = [
    {
      id: 'pkg-a',
      name: '@app/pkg-a',
      path: 'packages/pkg-a',
      dependencyCount: 3,
      inCycle: false,
      cycleIds: [],
      x: 100,
      y: 100,
    },
    {
      id: 'pkg-b',
      name: '@app/pkg-b',
      path: 'packages/pkg-b',
      dependencyCount: 1,
      inCycle: true,
      cycleIds: [0],
      x: 200,
      y: 200,
    },
  ]

  const mockLinks: D3Link[] = [
    {
      source: mockNodes[0],
      target: mockNodes[1],
      type: 'production',
      inCycle: false,
      cycleIds: [],
    },
  ]

  const emptyCircularNodeIds = new Set<string>()
  const circularNodeIds = new Set<string>(['pkg-b'])
  const emptyCircularEdgePairs = new Set<string>()

  // Mock canvas context
  let mockCanvasContext: Record<string, ReturnType<typeof vi.fn>>

  beforeEach(() => {
    mockCanvasContext = {
      save: vi.fn(),
      restore: vi.fn(),
      clearRect: vi.fn(),
      setTransform: vi.fn(),
      translate: vi.fn(),
      scale: vi.fn(),
      beginPath: vi.fn(),
      moveTo: vi.fn(),
      lineTo: vi.fn(),
      arc: vi.fn(),
      closePath: vi.fn(),
      fill: vi.fn(),
      stroke: vi.fn(),
      fillText: vi.fn(),
      rotate: vi.fn(),
    }

    // Add property setters
    Object.defineProperties(mockCanvasContext, {
      fillStyle: { set: vi.fn(), get: () => '' },
      strokeStyle: { set: vi.fn(), get: () => '' },
      lineWidth: { set: vi.fn(), get: () => 1 },
      globalAlpha: { set: vi.fn(), get: () => 1 },
      font: { set: vi.fn(), get: () => '' },
      textAlign: { set: vi.fn(), get: () => '' },
    })

    HTMLCanvasElement.prototype.getContext = vi.fn(
      () => mockCanvasContext
    ) as unknown as typeof HTMLCanvasElement.prototype.getContext

    // Mock devicePixelRatio
    vi.stubGlobal('devicePixelRatio', 1)
  })

  afterEach(() => {
    vi.clearAllMocks()
    vi.unstubAllGlobals()
  })

  it('should render a canvas element', () => {
    const { container } = render(
      <CanvasRenderer
        nodes={mockNodes}
        links={mockLinks}
        circularNodeIds={circularNodeIds}
        circularEdgePairs={emptyCircularEdgePairs}
        viewport={defaultViewport}
        onViewportChange={mockOnViewportChange}
        selectedNodeId={null}
        onNodeSelect={mockOnNodeSelect}
        onNodeHover={mockOnNodeHover}
        width={800}
        height={500}
      />
    )

    const canvas = container.querySelector('canvas')
    expect(canvas).toBeTruthy()
  })

  it('should set canvas dimensions based on props', () => {
    const { container } = render(
      <CanvasRenderer
        nodes={mockNodes}
        links={mockLinks}
        circularNodeIds={emptyCircularNodeIds}
        circularEdgePairs={emptyCircularEdgePairs}
        viewport={defaultViewport}
        onViewportChange={mockOnViewportChange}
        selectedNodeId={null}
        onNodeSelect={mockOnNodeSelect}
        onNodeHover={mockOnNodeHover}
        width={1024}
        height={768}
      />
    )

    const canvas = container.querySelector('canvas')
    expect(canvas?.style.width).toBe('1024px')
    expect(canvas?.style.height).toBe('768px')
  })

  it('should handle empty nodes gracefully', () => {
    expect(() => {
      render(
        <CanvasRenderer
          nodes={[]}
          links={[]}
          circularNodeIds={emptyCircularNodeIds}
          circularEdgePairs={emptyCircularEdgePairs}
          viewport={defaultViewport}
          onViewportChange={mockOnViewportChange}
          selectedNodeId={null}
          onNodeSelect={mockOnNodeSelect}
          onNodeHover={mockOnNodeHover}
          width={800}
          height={500}
        />
      )
    }).not.toThrow()
  })

  it('should apply crosshair cursor to canvas', () => {
    const { container } = render(
      <CanvasRenderer
        nodes={mockNodes}
        links={mockLinks}
        circularNodeIds={emptyCircularNodeIds}
        circularEdgePairs={emptyCircularEdgePairs}
        viewport={defaultViewport}
        onViewportChange={mockOnViewportChange}
        selectedNodeId={null}
        onNodeSelect={mockOnNodeSelect}
        onNodeHover={mockOnNodeHover}
        width={800}
        height={500}
      />
    )

    const canvas = container.querySelector('canvas')
    expect(canvas?.classList.contains('cursor-crosshair')).toBe(true)
  })

  it('should have touch-action: none for touch handling', () => {
    const { container } = render(
      <CanvasRenderer
        nodes={mockNodes}
        links={mockLinks}
        circularNodeIds={emptyCircularNodeIds}
        circularEdgePairs={emptyCircularEdgePairs}
        viewport={defaultViewport}
        onViewportChange={mockOnViewportChange}
        selectedNodeId={null}
        onNodeSelect={mockOnNodeSelect}
        onNodeHover={mockOnNodeHover}
        width={800}
        height={500}
      />
    )

    const canvas = container.querySelector('canvas')
    expect(canvas?.style.touchAction).toBe('none')
  })
})
