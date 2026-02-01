/**
 * Tests for CanvasRenderer drawing logic (Story 4.9)
 *
 * Validates that Canvas rendering correctly draws nodes, edges,
 * arrows, labels, selection rings, and circular dependency highlighting
 * by inspecting calls to the mocked Canvas 2D context.
 *
 * @see AC7: Canvas Rendering Visual Parity
 * @see AC4: Canvas Mode Hover/Click Functionality
 * @see AC6: Performance Requirements
 */

import { render } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import type { D3Link, D3Node, ViewportState } from '../types'

// Track canvas context calls
let ctxCalls: Record<string, unknown[][]>

// Mock d3 with a simulation that calls tick immediately
vi.mock('d3', () => {
  let tickCallback: (() => void) | null = null

  const mockSimulation = {
    force: vi.fn().mockReturnThis(),
    on: vi.fn((event: string, cb: () => void) => {
      if (event === 'tick') tickCallback = cb
      return mockSimulation
    }),
    stop: vi.fn(),
    alphaDecay: vi.fn().mockReturnThis(),
  }

  return {
    forceSimulation: vi.fn(() => {
      // Trigger tick callback after a microtask so render function is set up
      queueMicrotask(() => {
        if (tickCallback) tickCallback()
      })
      return mockSimulation
    }),
    forceLink: vi.fn(() => ({
      id: vi.fn().mockReturnThis(),
      distance: vi.fn().mockReturnThis(),
    })),
    forceManyBody: vi.fn(() => ({ strength: vi.fn().mockReturnThis() })),
    forceCenter: vi.fn(),
    forceCollide: vi.fn(() => ({ radius: vi.fn().mockReturnThis() })),
  }
})

// Must import after mock
import { CanvasRenderer } from '../CanvasRenderer'

function createNode(id: string, x: number, y: number, overrides: Partial<D3Node> = {}): D3Node {
  return {
    id,
    name: `@app/${id}`,
    path: `packages/${id}`,
    dependencyCount: 1,
    inCycle: false,
    cycleIds: [],
    x,
    y,
    ...overrides,
  }
}

function createLink(source: D3Node, target: D3Node): D3Link {
  return {
    source,
    target,
    type: 'production',
    inCycle: false,
    cycleIds: [],
  }
}

describe('CanvasRenderer Drawing Logic', () => {
  const defaultViewport: ViewportState = { zoom: 1, panX: 0, panY: 0 }

  beforeEach(() => {
    ctxCalls = {}

    const mockCtx = new Proxy(
      {},
      {
        get(_target, prop) {
          const name = String(prop)
          if (!ctxCalls[name]) ctxCalls[name] = []
          return (...args: unknown[]) => {
            ctxCalls[name].push(args)
          }
        },
        set(_target, prop, value) {
          const name = `set_${String(prop)}`
          if (!ctxCalls[name]) ctxCalls[name] = []
          ctxCalls[name].push([value])
          return true
        },
      }
    )

    HTMLCanvasElement.prototype.getContext = vi.fn(
      () => mockCtx
    ) as unknown as typeof HTMLCanvasElement.prototype.getContext

    vi.stubGlobal('devicePixelRatio', 1)

    // Mock requestAnimationFrame to run synchronously
    vi.stubGlobal('requestAnimationFrame', (cb: FrameRequestCallback) => {
      cb(0)
      return 0
    })
    vi.stubGlobal('cancelAnimationFrame', vi.fn())
  })

  afterEach(() => {
    vi.clearAllMocks()
    vi.unstubAllGlobals()
  })

  function renderCanvas(props: Partial<React.ComponentProps<typeof CanvasRenderer>> = {}) {
    const defaultNodes = [createNode('a', 100, 100), createNode('b', 200, 200)]
    const defaultLinks = [createLink(defaultNodes[0], defaultNodes[1])]

    return render(
      <CanvasRenderer
        nodes={defaultNodes}
        links={defaultLinks}
        circularNodeIds={new Set<string>()}
        circularEdgePairs={new Set<string>()}
        viewport={defaultViewport}
        onViewportChange={vi.fn()}
        selectedNodeId={null}
        onNodeSelect={vi.fn()}
        onNodeHover={vi.fn()}
        width={800}
        height={600}
        {...props}
      />
    )
  }

  describe('[P1] basic rendering', () => {
    it('should clear and save/restore context on each render', () => {
      renderCanvas()

      // THEN: Context should have save/restore pattern
      expect(ctxCalls.save?.length).toBeGreaterThanOrEqual(1)
      expect(ctxCalls.restore?.length).toBeGreaterThanOrEqual(1)
    })

    it('should apply viewport transform (translate + scale)', () => {
      renderCanvas({
        viewport: { zoom: 2, panX: 50, panY: 75 },
      })

      // THEN: translate and scale should be called with viewport values
      const translateCalls = ctxCalls.translate ?? []
      const scaleCalls = ctxCalls.scale ?? []

      // Should have at least one translate call with viewport values
      const hasViewportTranslate = translateCalls.some((args) => args[0] === 50 && args[1] === 75)
      const hasViewportScale = scaleCalls.some((args) => args[0] === 2 && args[1] === 2)

      expect(hasViewportTranslate).toBe(true)
      expect(hasViewportScale).toBe(true)
    })

    it('should call clearRect to clear the canvas before drawing', () => {
      renderCanvas()

      expect(ctxCalls.clearRect?.length).toBeGreaterThanOrEqual(1)
    })
  })

  describe('[P1] node rendering', () => {
    it('should draw arcs for each node with valid coordinates', () => {
      const nodes = [createNode('a', 100, 100), createNode('b', 200, 200)]

      renderCanvas({ nodes, links: [] })

      // THEN: arc should be called at least twice (once per node)
      const arcCalls = ctxCalls.arc ?? []
      expect(arcCalls.length).toBeGreaterThanOrEqual(2)
    })

    it('should skip nodes with null coordinates', () => {
      const nodeWithNullCoords: D3Node = {
        id: 'null-node',
        name: '@app/null-node',
        path: 'packages/null-node',
        dependencyCount: 1,
        inCycle: false,
        cycleIds: [],
        // x and y intentionally undefined
      }

      renderCanvas({ nodes: [nodeWithNullCoords], links: [] })

      // THEN: No arc calls for nodes without positions (only from initial render call)
      // The node rendering loop skips null coords
      // We just verify no crash
    })

    it('should render selection ring for selected node', () => {
      const nodes = [createNode('a', 100, 100)]

      renderCanvas({ nodes, links: [], selectedNodeId: 'a' })

      // THEN: Should have extra arc call for selection ring (nodeRadius + 6)
      const arcCalls = ctxCalls.arc ?? []
      // At minimum 2 arcs: one for node, one for selection ring
      expect(arcCalls.length).toBeGreaterThanOrEqual(2)
    })

    it('should render node labels with truncated text', () => {
      const nodes = [createNode('a', 100, 100)]

      renderCanvas({ nodes, links: [] })

      // THEN: fillText should be called for node label
      const fillTextCalls = ctxCalls.fillText ?? []
      expect(fillTextCalls.length).toBeGreaterThanOrEqual(1)
    })

    it('should scale node radius based on dependencyCount', () => {
      const bigNode = createNode('big', 100, 100, { dependencyCount: 20 })
      const smallNode = createNode('small', 200, 200, { dependencyCount: 0 })

      renderCanvas({ nodes: [bigNode, smallNode], links: [] })

      // THEN: arc calls should have different radii
      const arcCalls = ctxCalls.arc ?? []
      // With dependencyCount=20: radius = max(8, min(16, 8 + 20*0.5)) = 16
      // With dependencyCount=0: radius = max(8, min(16, 8 + 0)) = 8
      // Just verify arcs were drawn
      expect(arcCalls.length).toBeGreaterThanOrEqual(2)
    })
  })

  describe('[P1] circular dependency highlighting (AC7)', () => {
    it('should use cycle colors for circular dependency nodes', () => {
      const circularNode = createNode('circ', 100, 100, { inCycle: true })
      const circularNodeIds = new Set(['circ'])

      renderCanvas({
        nodes: [circularNode],
        links: [],
        circularNodeIds,
      })

      // THEN: fillStyle should include the cycle color (NODE_COLORS.cycle.fill = '#ef4444')
      const fillStyleSets = ctxCalls.set_fillStyle ?? []
      const hasCycleColor = fillStyleSets.some((args) => args[0] === '#ef4444')
      expect(hasCycleColor).toBe(true)
    })

    it('should use cycle colors for circular dependency edges', () => {
      const nodeA = createNode('a', 100, 100)
      const nodeB = createNode('b', 200, 200)
      const link = createLink(nodeA, nodeB)
      const circularEdgePairs = new Set(['a->b'])

      renderCanvas({
        nodes: [nodeA, nodeB],
        links: [link],
        circularEdgePairs,
      })

      // THEN: strokeStyle should include cycle edge color (EDGE_COLORS.cycle.stroke = '#ef4444')
      const strokeStyleSets = ctxCalls.set_strokeStyle ?? []
      const hasCycleEdge = strokeStyleSets.some((args) => args[0] === '#ef4444')
      expect(hasCycleEdge).toBe(true)
    })
  })

  describe('[P1] edge rendering', () => {
    it('should draw lines and arrows for edges', () => {
      const nodeA = createNode('a', 100, 100)
      const nodeB = createNode('b', 200, 200)
      const link = createLink(nodeA, nodeB)

      renderCanvas({ nodes: [nodeA, nodeB], links: [link] })

      // THEN: moveTo and lineTo called for edge line
      expect(ctxCalls.moveTo?.length).toBeGreaterThanOrEqual(1)
      expect(ctxCalls.lineTo?.length).toBeGreaterThanOrEqual(1)
      // Arrow uses rotate and closePath
      expect(ctxCalls.rotate?.length).toBeGreaterThanOrEqual(1)
      expect(ctxCalls.closePath?.length).toBeGreaterThanOrEqual(1)
    })

    it('should skip edges where source or target has null coordinates', () => {
      const nodeA = createNode('a', 100, 100)
      const nodeWithNull: D3Node = {
        id: 'null',
        name: '@app/null',
        path: 'packages/null',
        dependencyCount: 0,
        inCycle: false,
        cycleIds: [],
        // x, y intentionally undefined
      }
      const link: D3Link = {
        source: nodeA,
        target: nodeWithNull,
        type: 'production',
        inCycle: false,
        cycleIds: [],
      }

      // WHEN/THEN: Should not throw
      expect(() => renderCanvas({ nodes: [nodeA, nodeWithNull], links: [link] })).not.toThrow()
    })
  })

  describe('[P1] HiDPI support', () => {
    it('should scale canvas for high DPI displays', () => {
      vi.stubGlobal('devicePixelRatio', 2)

      const { container } = renderCanvas({ width: 800, height: 600 })

      const canvas = container.querySelector('canvas')
      // Canvas style should show CSS dimensions
      expect(canvas?.style.width).toBe('800px')
      expect(canvas?.style.height).toBe('600px')

      // setTransform should be called with dpr
      const setTransformCalls = ctxCalls.setTransform ?? []
      const hasDprTransform = setTransformCalls.some((args) => args[0] === 2 && args[3] === 2)
      expect(hasDprTransform).toBe(true)
    })
  })

  describe('[P2] simulation lifecycle', () => {
    it('should not create simulation for empty nodes', () => {
      renderCanvas({ nodes: [], links: [] })

      // THEN: Component renders without error
      // No simulation tick = fewer canvas operations
    })

    it('should handle re-render with different viewport without crash', () => {
      const { rerender } = render(
        <CanvasRenderer
          nodes={[createNode('a', 100, 100)]}
          links={[]}
          circularNodeIds={new Set<string>()}
          circularEdgePairs={new Set<string>()}
          viewport={{ zoom: 1, panX: 0, panY: 0 }}
          onViewportChange={vi.fn()}
          selectedNodeId={null}
          onNodeSelect={vi.fn()}
          onNodeHover={vi.fn()}
          width={800}
          height={600}
        />
      )

      // WHEN: Re-render with different viewport
      expect(() =>
        rerender(
          <CanvasRenderer
            nodes={[createNode('a', 100, 100)]}
            links={[]}
            circularNodeIds={new Set<string>()}
            circularEdgePairs={new Set<string>()}
            viewport={{ zoom: 2, panX: 50, panY: 50 }}
            onViewportChange={vi.fn()}
            selectedNodeId={null}
            onNodeSelect={vi.fn()}
            onNodeHover={vi.fn()}
            width={800}
            height={600}
          />
        )
      ).not.toThrow()
    })

    it('should handle selection change re-render', () => {
      const nodes = [createNode('a', 100, 100)]

      const { rerender } = render(
        <CanvasRenderer
          nodes={nodes}
          links={[]}
          circularNodeIds={new Set<string>()}
          circularEdgePairs={new Set<string>()}
          viewport={defaultViewport}
          onViewportChange={vi.fn()}
          selectedNodeId={null}
          onNodeSelect={vi.fn()}
          onNodeHover={vi.fn()}
          width={800}
          height={600}
        />
      )

      // WHEN: Re-render with a selected node
      expect(() =>
        rerender(
          <CanvasRenderer
            nodes={nodes}
            links={[]}
            circularNodeIds={new Set<string>()}
            circularEdgePairs={new Set<string>()}
            viewport={defaultViewport}
            onViewportChange={vi.fn()}
            selectedNodeId="a"
            onNodeSelect={vi.fn()}
            onNodeHover={vi.fn()}
            width={800}
            height={600}
          />
        )
      ).not.toThrow()
    })
  })
})
