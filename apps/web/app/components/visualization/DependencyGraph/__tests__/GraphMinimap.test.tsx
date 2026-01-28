/**
 * Tests for GraphMinimap component
 *
 * @see Story 4.4: Add Zoom, Pan, and Navigation Controls
 * @vitest-environment jsdom
 */
import { fireEvent, render } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import type { GraphMinimapProps } from '../GraphMinimap'
import { GraphMinimap } from '../GraphMinimap'
import type { D3Link, D3Node } from '../types'

// Generate mock nodes
const generateMockNodes = (count: number): D3Node[] =>
  Array.from({ length: count }, (_, i) => ({
    id: `node-${i}`,
    name: `Node ${i}`,
    path: `/path/${i}`,
    dependencyCount: 1,
    x: Math.random() * 500,
    y: Math.random() * 500,
    inCycle: i < 5, // First 5 nodes are in cycle
    cycleIds: i < 5 ? [0] : [],
  }))

// Generate mock links
const generateMockLinks = (nodes: D3Node[]): D3Link[] =>
  nodes.slice(0, -1).map((node, i) => ({
    source: node,
    target: nodes[i + 1],
    type: 'production' as const,
    inCycle: false,
    cycleIds: [],
  }))

describe('GraphMinimap', () => {
  const smallNodes = generateMockNodes(30)
  const smallLinks = generateMockLinks(smallNodes)

  const largeNodes = generateMockNodes(60)
  const largeLinks = generateMockLinks(largeNodes)

  const defaultProps: GraphMinimapProps = {
    nodes: largeNodes,
    links: largeLinks,
    viewportBounds: { x: 0, y: 0, width: 800, height: 600 },
    graphBounds: { x: 0, y: 0, width: 500, height: 500 },
    onNavigate: vi.fn(),
  }

  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('visibility based on node count (AC5)', () => {
    it('should render minimap for graphs with >= 50 nodes', () => {
      const { container } = render(<GraphMinimap {...defaultProps} />)

      expect(container.querySelector('svg')).toBeInTheDocument()
    })

    it('should NOT render minimap for graphs with < 50 nodes', () => {
      const { container } = render(
        <GraphMinimap {...defaultProps} nodes={smallNodes} links={smallLinks} />
      )

      expect(container.querySelector('svg')).not.toBeInTheDocument()
    })

    it('should render minimap for exactly 50 nodes', () => {
      const fiftyNodes = generateMockNodes(50)
      const fiftyLinks = generateMockLinks(fiftyNodes)

      const { container } = render(
        <GraphMinimap {...defaultProps} nodes={fiftyNodes} links={fiftyLinks} />
      )

      expect(container.querySelector('svg')).toBeInTheDocument()
    })

    it('should not render for 49 nodes', () => {
      const fortyNineNodes = generateMockNodes(49)
      const fortyNineLinks = generateMockLinks(fortyNineNodes)

      const { container } = render(
        <GraphMinimap {...defaultProps} nodes={fortyNineNodes} links={fortyNineLinks} />
      )

      expect(container.querySelector('svg')).not.toBeInTheDocument()
    })
  })

  describe('rendering (AC5)', () => {
    it('should render in corner position', () => {
      const { container } = render(<GraphMinimap {...defaultProps} />)

      const wrapper = container.firstChild
      expect(wrapper).toHaveClass('absolute')
      expect(wrapper).toHaveClass('top-4')
      expect(wrapper).toHaveClass('left-4')
    })

    it('should render with default dimensions', () => {
      const { container } = render(<GraphMinimap {...defaultProps} />)

      const svg = container.querySelector('svg')
      expect(svg).toHaveAttribute('width', '150')
      expect(svg).toHaveAttribute('height', '100')
    })

    it('should render with custom dimensions', () => {
      const { container } = render(<GraphMinimap {...defaultProps} width={200} height={150} />)

      const svg = container.querySelector('svg')
      expect(svg).toHaveAttribute('width', '200')
      expect(svg).toHaveAttribute('height', '150')
    })

    it('should render nodes as circles', () => {
      const { container } = render(<GraphMinimap {...defaultProps} />)

      const circles = container.querySelectorAll('circle')
      // Should have node circles (60 nodes)
      // Plus viewport indicator rect is not a circle
      expect(circles.length).toBeGreaterThanOrEqual(60)
    })

    it('should render links as lines', () => {
      const { container } = render(<GraphMinimap {...defaultProps} />)

      const lines = container.querySelectorAll('line')
      expect(lines.length).toBe(59) // 60 nodes, 59 links
    })

    it('should differentiate cycle nodes with color', () => {
      const { container } = render(<GraphMinimap {...defaultProps} />)

      const circles = container.querySelectorAll('circle')
      const fillColors = Array.from(circles).map((c) => c.getAttribute('fill'))

      // Should have both normal (indigo) and cycle (red) colors
      expect(fillColors.some((c) => c === '#ef4444')).toBe(true) // Red for cycle
      expect(fillColors.some((c) => c === '#4f46e5')).toBe(true) // Indigo for normal
    })
  })

  describe('viewport indicator (AC5)', () => {
    it('should render viewport indicator rectangle', () => {
      const { container } = render(<GraphMinimap {...defaultProps} />)

      const rects = container.querySelectorAll('rect')
      // Should have at least background rect and viewport indicator
      expect(rects.length).toBeGreaterThanOrEqual(2)
    })

    it('should have styled viewport indicator', () => {
      const { container } = render(<GraphMinimap {...defaultProps} />)

      const rects = container.querySelectorAll('rect')
      const viewportRect = Array.from(rects).find(
        (r) => r.getAttribute('stroke') === '#6366f1' // Indigo stroke for viewport
      )

      expect(viewportRect).toBeTruthy()
    })
  })

  describe('navigation interaction (AC5)', () => {
    it('should call onNavigate when clicked', () => {
      const onNavigate = vi.fn()
      const { container } = render(<GraphMinimap {...defaultProps} onNavigate={onNavigate} />)

      const svg = container.querySelector('svg')
      expect(svg).not.toBeNull()
      if (svg) {
        fireEvent.click(svg, { clientX: 75, clientY: 50 })
      }

      expect(onNavigate).toHaveBeenCalledTimes(1)
    })

    it('should call onNavigate with graph coordinates', () => {
      const onNavigate = vi.fn()
      const { container } = render(<GraphMinimap {...defaultProps} onNavigate={onNavigate} />)

      const svg = container.querySelector('svg')
      expect(svg).not.toBeNull()
      if (!svg) return

      // Mock getBoundingClientRect
      svg.getBoundingClientRect = vi.fn(() => ({
        left: 0,
        top: 0,
        right: 150,
        bottom: 100,
        width: 150,
        height: 100,
        x: 0,
        y: 0,
        toJSON: () => ({}),
      }))

      fireEvent.click(svg, { clientX: 75, clientY: 50 })

      expect(onNavigate).toHaveBeenCalledWith(expect.any(Number), expect.any(Number))
    })

    it('should have pointer cursor for interactivity', () => {
      const { container } = render(<GraphMinimap {...defaultProps} />)

      // The interactive wrapper div (not SVG) has the cursor-pointer class
      const wrapper = container.querySelector('button')
      expect(wrapper).toHaveClass('cursor-pointer')
    })

    it('should call onNavigate when dragging (AC5 drag-to-navigate)', () => {
      const onNavigate = vi.fn()
      const { container } = render(<GraphMinimap {...defaultProps} onNavigate={onNavigate} />)

      // Events are now on the wrapper div, not the SVG
      const wrapper = container.querySelector('button')
      expect(wrapper).not.toBeNull()
      if (!wrapper) return

      // Mock getBoundingClientRect on the wrapper
      wrapper.getBoundingClientRect = vi.fn(() => ({
        left: 0,
        top: 0,
        right: 150,
        bottom: 100,
        width: 150,
        height: 100,
        x: 0,
        y: 0,
        toJSON: () => ({}),
      }))

      // Simulate drag: mousedown, mousemove, mouseup
      fireEvent.mouseDown(wrapper, { clientX: 50, clientY: 30 })
      expect(onNavigate).toHaveBeenCalledTimes(1)

      fireEvent.mouseMove(wrapper, { clientX: 75, clientY: 50 })
      expect(onNavigate).toHaveBeenCalledTimes(2)

      fireEvent.mouseMove(wrapper, { clientX: 100, clientY: 60 })
      expect(onNavigate).toHaveBeenCalledTimes(3)

      fireEvent.mouseUp(wrapper)

      // After mouseUp, mousemove should not trigger navigate
      fireEvent.mouseMove(wrapper, { clientX: 120, clientY: 70 })
      expect(onNavigate).toHaveBeenCalledTimes(3)
    })

    it('should stop dragging when mouse leaves minimap', () => {
      const onNavigate = vi.fn()
      const { container } = render(<GraphMinimap {...defaultProps} onNavigate={onNavigate} />)

      // Events are now on the wrapper div, not the SVG
      const wrapper = container.querySelector('button')
      expect(wrapper).not.toBeNull()
      if (!wrapper) return

      wrapper.getBoundingClientRect = vi.fn(() => ({
        left: 0,
        top: 0,
        right: 150,
        bottom: 100,
        width: 150,
        height: 100,
        x: 0,
        y: 0,
        toJSON: () => ({}),
      }))

      // Start dragging
      fireEvent.mouseDown(wrapper, { clientX: 50, clientY: 30 })
      fireEvent.mouseMove(wrapper, { clientX: 75, clientY: 50 })
      expect(onNavigate).toHaveBeenCalledTimes(2)

      // Mouse leaves
      fireEvent.mouseLeave(wrapper)

      // Move should not trigger navigate after leave
      fireEvent.mouseMove(wrapper, { clientX: 100, clientY: 60 })
      expect(onNavigate).toHaveBeenCalledTimes(2)
    })
  })

  describe('accessibility', () => {
    it('should have focus indicator styles (L1 fix)', () => {
      const { container } = render(<GraphMinimap {...defaultProps} />)

      // Focus styles are on the wrapper div, not the SVG
      const wrapper = container.querySelector('button')
      expect(wrapper).toHaveClass('focus:ring-2')
      expect(wrapper).toHaveClass('focus:ring-indigo-500')
    })

    it('should be focusable for keyboard navigation', () => {
      const { container } = render(<GraphMinimap {...defaultProps} />)

      const button = container.querySelector('button')
      // Native button elements are focusable by default, no tabIndex needed
      expect(button).toBeInTheDocument()
      // Verify it's not explicitly set to non-focusable
      expect(button).not.toHaveAttribute('tabindex', '-1')
    })

    it('should use semantic button element for interactivity', () => {
      const { container } = render(<GraphMinimap {...defaultProps} />)

      const button = container.querySelector('button')
      expect(button).toBeInTheDocument()
      expect(button).toHaveAttribute('type', 'button')
    })

    it('should have aria-label for screen readers', () => {
      const { container } = render(<GraphMinimap {...defaultProps} />)

      const wrapper = container.querySelector('button')
      expect(wrapper).toHaveAttribute('aria-label')
    })
  })

  describe('edge cases', () => {
    it('should handle empty viewport bounds', () => {
      const { container } = render(
        <GraphMinimap {...defaultProps} viewportBounds={{ x: 0, y: 0, width: 0, height: 0 }} />
      )

      // Should still render
      expect(container.querySelector('svg')).toBeInTheDocument()
    })

    it('should handle empty graph bounds', () => {
      const { container } = render(
        <GraphMinimap {...defaultProps} graphBounds={{ x: 0, y: 0, width: 0, height: 0 }} />
      )

      // Should still render (scale will be 1)
      expect(container.querySelector('svg')).toBeInTheDocument()
    })

    it('should handle nodes with undefined positions', () => {
      const nodesWithUndefined = largeNodes.map((n, i) =>
        i % 10 === 0 ? { ...n, x: undefined, y: undefined } : n
      )

      const { container } = render(<GraphMinimap {...defaultProps} nodes={nodesWithUndefined} />)

      expect(container.querySelector('svg')).toBeInTheDocument()
    })

    it('should handle links with string source/target', () => {
      const stringLinks: D3Link[] = largeNodes.slice(0, -1).map((node, i) => ({
        source: node.id, // String instead of node object
        target: largeNodes[i + 1].id,
        type: 'production' as const,
        inCycle: false,
        cycleIds: [],
      }))

      const { container } = render(<GraphMinimap {...defaultProps} links={stringLinks} />)

      // Should render but links may not be visible if lookup fails
      expect(container.querySelector('svg')).toBeInTheDocument()
    })
  })
})
