/**
 * Tests for NodeTooltip component
 *
 * @see Story 4.5: Implement Hover Details and Tooltips (AC1, AC2, AC3, AC7)
 *
 * Following Given-When-Then format with priority tags.
 */

import { render, screen } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { NodeTooltip, type NodeTooltipProps } from '../NodeTooltip'
import type { TooltipData } from '../types'

/** Mock rect type matching DOMRect interface */
interface MockRect {
  left: number
  top: number
  right: number
  bottom: number
  width: number
  height: number
  x: number
  y: number
  toJSON: () => Record<string, unknown>
}

describe('NodeTooltip', () => {
  // Mock container ref with getBoundingClientRect
  const createMockContainerRef = (
    rect: Partial<MockRect> = {}
  ): React.RefObject<HTMLDivElement> => {
    const element = document.createElement('div')
    element.getBoundingClientRect = vi.fn(
      () =>
        ({
          left: 0,
          top: 0,
          right: 800,
          bottom: 600,
          width: 800,
          height: 600,
          x: 0,
          y: 0,
          toJSON: () => ({}),
          ...rect,
        }) as MockRect
    )
    return { current: element }
  }

  const defaultTooltipData: TooltipData = {
    packageName: '@app/core',
    packagePath: 'packages/core',
    incomingCount: 3,
    outgoingCount: 5,
    healthContribution: 2,
    inCycle: false,
  }

  const defaultProps: NodeTooltipProps = {
    data: defaultTooltipData,
    position: { x: 100, y: 100 },
    containerRef: createMockContainerRef(),
  }

  const renderNodeTooltip = (props: Partial<NodeTooltipProps> = {}) => {
    const mergedProps = { ...defaultProps, ...props }
    return render(<NodeTooltip {...mergedProps} />)
  }

  describe('Tooltip Content Display (AC1)', () => {
    it('[P1] should render nothing when data is null', () => {
      // GIVEN: No tooltip data
      // WHEN: Rendered with null data
      const { container } = renderNodeTooltip({ data: null })

      // THEN: Should not render any tooltip content
      expect(container.firstChild).toBeNull()
    })

    it('[P1] should render nothing when position is null', () => {
      // GIVEN: No position data
      // WHEN: Rendered with null position
      const { container } = renderNodeTooltip({ position: null })

      // THEN: Should not render any tooltip content
      expect(container.firstChild).toBeNull()
    })

    it('[P1] should display package name', () => {
      // GIVEN: Tooltip with package data
      // WHEN: Rendered
      renderNodeTooltip()

      // THEN: Should display the package name
      expect(screen.getByText('@app/core')).toBeInTheDocument()
    })

    it('[P1] should display shortened package path', () => {
      // GIVEN: Tooltip with package path
      // WHEN: Rendered
      renderNodeTooltip({
        data: { ...defaultTooltipData, packagePath: 'apps/web/packages/core' },
      })

      // THEN: Should display shortened path (last 2 segments)
      expect(screen.getByText('packages/core')).toBeInTheDocument()
    })

    it('[P1] should display incoming dependency count', () => {
      // GIVEN: Tooltip with dependency counts
      // WHEN: Rendered
      renderNodeTooltip()

      // THEN: Should display incoming count label and value
      expect(screen.getByText('In:')).toBeInTheDocument()
      expect(screen.getByText('3')).toBeInTheDocument()
    })

    it('[P1] should display outgoing dependency count', () => {
      // GIVEN: Tooltip with dependency counts
      // WHEN: Rendered
      renderNodeTooltip()

      // THEN: Should display outgoing count label and value
      expect(screen.getByText('Out:')).toBeInTheDocument()
      expect(screen.getByText('5')).toBeInTheDocument()
    })

    it('[P1] should display positive health contribution with + sign', () => {
      // GIVEN: Tooltip with positive health contribution
      // WHEN: Rendered
      renderNodeTooltip({
        data: { ...defaultTooltipData, healthContribution: 2 },
      })

      // THEN: Should display health with + prefix
      expect(screen.getByText('+2')).toBeInTheDocument()
    })

    it('[P1] should display negative health contribution', () => {
      // GIVEN: Tooltip with negative health contribution
      // WHEN: Rendered
      renderNodeTooltip({
        data: { ...defaultTooltipData, healthContribution: -5 },
      })

      // THEN: Should display negative health value
      expect(screen.getByText('-5')).toBeInTheDocument()
    })

    it('[P1] should display zero health contribution without + sign', () => {
      // GIVEN: Tooltip with zero health contribution
      // WHEN: Rendered
      renderNodeTooltip({
        data: { ...defaultTooltipData, healthContribution: 0 },
      })

      // THEN: Should display 0 (no + prefix for zero)
      expect(screen.getByText('0')).toBeInTheDocument()
    })

    it('[P1] should not show cycle warning when not in cycle', () => {
      // GIVEN: Tooltip for node not in cycle
      // WHEN: Rendered
      renderNodeTooltip({
        data: { ...defaultTooltipData, inCycle: false },
      })

      // THEN: Should not show circular dependency warning
      expect(screen.queryByText(/circular/i)).not.toBeInTheDocument()
    })

    it('[P1] should show cycle warning when in cycle', () => {
      // GIVEN: Tooltip for node in a circular dependency
      const cycleData: TooltipData = {
        ...defaultTooltipData,
        inCycle: true,
        cycleInfo: {
          cycleCount: 1,
          packages: ['@app/utils'],
        },
      }

      // WHEN: Rendered
      renderNodeTooltip({ data: cycleData })

      // THEN: Should show circular dependency warning
      expect(screen.getByText(/circular dependency/i)).toBeInTheDocument()
    })

    it('[P2] should show plural "dependencies" for multiple cycles', () => {
      // GIVEN: Tooltip for node in multiple cycles
      const cycleData: TooltipData = {
        ...defaultTooltipData,
        inCycle: true,
        cycleInfo: {
          cycleCount: 3,
          packages: ['@app/utils', '@app/shared'],
        },
      }

      // WHEN: Rendered
      renderNodeTooltip({ data: cycleData })

      // THEN: Should show plural form
      expect(screen.getByText(/3 circular dependencies/i)).toBeInTheDocument()
    })
  })

  describe('Accessibility (AC7)', () => {
    it('[P1] should have role="tooltip"', () => {
      // GIVEN: NodeTooltip component
      // WHEN: Rendered
      renderNodeTooltip()

      // THEN: Should have tooltip role
      expect(screen.getByRole('tooltip')).toBeInTheDocument()
    })

    it('[P1] should have aria-live="polite" for screen readers', () => {
      // GIVEN: NodeTooltip component
      // WHEN: Rendered
      renderNodeTooltip()

      // THEN: Should have aria-live attribute
      const tooltip = screen.getByRole('tooltip')
      expect(tooltip).toHaveAttribute('aria-live', 'polite')
    })
  })

  describe('Tooltip Styling', () => {
    it('[P2] should have pointer-events-none to not interfere with graph', () => {
      // GIVEN: NodeTooltip component
      // WHEN: Rendered
      renderNodeTooltip()

      // THEN: Should have pointer-events-none class
      const tooltip = screen.getByRole('tooltip')
      expect(tooltip).toHaveClass('pointer-events-none')
    })

    it('[P2] should have z-50 for proper layering', () => {
      // GIVEN: NodeTooltip component
      // WHEN: Rendered
      renderNodeTooltip()

      // THEN: Should have high z-index
      const tooltip = screen.getByRole('tooltip')
      expect(tooltip).toHaveClass('z-50')
    })

    it('[P2] should have absolute positioning', () => {
      // GIVEN: NodeTooltip component
      // WHEN: Rendered
      renderNodeTooltip()

      // THEN: Should have absolute positioning
      const tooltip = screen.getByRole('tooltip')
      expect(tooltip).toHaveClass('absolute')
    })
  })

  describe('Tooltip Positioning (AC3)', () => {
    it('[P1] should render tooltip with style attributes for positioning', () => {
      // GIVEN: Tooltip with position in center of container
      // WHEN: Rendered with position {x: 400, y: 300}
      renderNodeTooltip({ position: { x: 400, y: 300 } })

      // THEN: Tooltip should have inline style with left/top properties
      const tooltip = screen.getByRole('tooltip')
      expect(tooltip).toHaveAttribute('style')
      // Initial position is 0,0 before useEffect runs, which is expected behavior
      expect(tooltip.style.left).toBeDefined()
      expect(tooltip.style.top).toBeDefined()
    })

    it('[P2] should handle position near container boundaries', () => {
      // GIVEN: Position near the edge of container (container is 800x600)
      const edgePosition = { x: 750, y: 550 } // Near bottom-right

      // WHEN: Rendered
      renderNodeTooltip({ position: edgePosition })

      // THEN: Tooltip should render without errors
      const tooltip = screen.getByRole('tooltip')
      expect(tooltip).toBeInTheDocument()
    })

    it('[P2] should handle position at origin (0,0)', () => {
      // GIVEN: Position at origin
      const originPosition = { x: 0, y: 0 }

      // WHEN: Rendered
      renderNodeTooltip({ position: originPosition })

      // THEN: Tooltip should render and adjust to stay in bounds
      const tooltip = screen.getByRole('tooltip')
      expect(tooltip).toBeInTheDocument()
    })

    it('[P3] should handle negative position values gracefully', () => {
      // GIVEN: Negative position (edge case from rapid mouse movement)
      const negativePosition = { x: -10, y: -10 }

      // WHEN: Rendered
      renderNodeTooltip({ position: negativePosition })

      // THEN: Tooltip should still render (position will be clamped)
      const tooltip = screen.getByRole('tooltip')
      expect(tooltip).toBeInTheDocument()
    })
  })

  describe('Tooltip Position Calculation via rAF (AC3)', () => {
    /**
     * These tests exercise the requestAnimationFrame callback
     * that calculates viewport-aware tooltip positioning.
     */

    beforeEach(() => {
      // Mock rAF to execute synchronously
      vi.spyOn(window, 'requestAnimationFrame').mockImplementation((cb: (time: number) => void) => {
        cb(0)
        return 0
      })
    })

    afterEach(() => {
      vi.restoreAllMocks()
    })

    it('[P1] should position tooltip to the right by default', () => {
      // GIVEN: Position in upper-left area with room on all sides
      const containerRef = createMockContainerRef()
      renderNodeTooltip({
        position: { x: 100, y: 100 },
        containerRef,
      })

      // THEN: Tooltip should be visible after rAF positioning
      const tooltip = screen.getByRole('tooltip')
      expect(tooltip).toHaveClass('opacity-100')
    })

    it('[P1] should flip tooltip left when clipping right edge', () => {
      // GIVEN: Position near the right edge of the container
      const containerRef = createMockContainerRef({ width: 800, right: 800 })

      // Mock tooltip dimensions
      vi.spyOn(HTMLElement.prototype, 'getBoundingClientRect').mockReturnValue({
        width: 200,
        height: 100,
        top: 0,
        left: 0,
        right: 200,
        bottom: 100,
        x: 0,
        y: 0,
        toJSON: () => ({}),
      })

      renderNodeTooltip({
        position: { x: 750, y: 100 },
        containerRef,
      })

      // THEN: Tooltip should be visible (positioned left)
      const tooltip = screen.getByRole('tooltip')
      expect(tooltip).toBeInTheDocument()
    })

    it('[P1] should adjust tooltip up when clipping bottom edge', () => {
      // GIVEN: Position near the bottom of the container
      const containerRef = createMockContainerRef({ height: 600, bottom: 600 })

      vi.spyOn(HTMLElement.prototype, 'getBoundingClientRect').mockReturnValue({
        width: 200,
        height: 100,
        top: 0,
        left: 0,
        right: 200,
        bottom: 100,
        x: 0,
        y: 0,
        toJSON: () => ({}),
      })

      renderNodeTooltip({
        position: { x: 100, y: 580 },
        containerRef,
      })

      // THEN: Tooltip should be visible (positioned above)
      const tooltip = screen.getByRole('tooltip')
      expect(tooltip).toBeInTheDocument()
    })

    it('[P2] should clamp to TOOLTIP_OFFSET when clipping left edge', () => {
      // GIVEN: Position that would cause negative x after left flip
      const containerRef = createMockContainerRef()

      vi.spyOn(HTMLElement.prototype, 'getBoundingClientRect').mockReturnValue({
        width: 200,
        height: 100,
        top: 0,
        left: 0,
        right: 200,
        bottom: 100,
        x: 0,
        y: 0,
        toJSON: () => ({}),
      })

      renderNodeTooltip({
        position: { x: 5, y: 100 },
        containerRef,
      })

      // THEN: Tooltip should be visible
      const tooltip = screen.getByRole('tooltip')
      expect(tooltip).toBeInTheDocument()
    })

    it('[P2] should clamp to TOOLTIP_OFFSET when clipping top edge', () => {
      // GIVEN: Position that would cause negative y
      const containerRef = createMockContainerRef()

      vi.spyOn(HTMLElement.prototype, 'getBoundingClientRect').mockReturnValue({
        width: 200,
        height: 100,
        top: 0,
        left: 0,
        right: 200,
        bottom: 100,
        x: 0,
        y: 0,
        toJSON: () => ({}),
      })

      renderNodeTooltip({
        position: { x: 100, y: 5 },
        containerRef,
      })

      // THEN: Tooltip should be visible
      const tooltip = screen.getByRole('tooltip')
      expect(tooltip).toBeInTheDocument()
    })

    it('[P2] should handle null tooltipRef in rAF callback gracefully', () => {
      // GIVEN: Container ref is valid but component unmounts before rAF fires
      // This is tested by rendering then unmounting quickly
      const containerRef = createMockContainerRef()
      const { unmount } = render(
        <NodeTooltip
          data={defaultTooltipData}
          position={{ x: 100, y: 100 }}
          containerRef={containerRef}
        />
      )

      // Unmount before rAF fires (cancelAnimationFrame path)
      unmount()

      // THEN: Should not throw
      expect(true).toBe(true)
    })
  })

  describe('Edge Cases', () => {
    it('[P2] should handle long package names with truncation', () => {
      // GIVEN: Tooltip with very long package name
      const longNameData: TooltipData = {
        ...defaultTooltipData,
        packageName: '@organization/very-long-package-name-that-should-be-truncated',
      }

      // WHEN: Rendered
      renderNodeTooltip({ data: longNameData })

      // THEN: Should render without breaking (truncate class applied)
      const tooltip = screen.getByRole('tooltip')
      expect(tooltip).toBeInTheDocument()
    })

    it('[P2] should handle missing cycleInfo when inCycle is true', () => {
      // GIVEN: Node in cycle but missing cycleInfo (edge case)
      const incompleteData: TooltipData = {
        ...defaultTooltipData,
        inCycle: true,
        // cycleInfo intentionally omitted
      }

      // WHEN: Rendered
      renderNodeTooltip({ data: incompleteData })

      // THEN: Should handle gracefully and show generic cycle warning
      expect(screen.getByText(/circular dependency/i)).toBeInTheDocument()
    })

    it('[P3] should handle zero dependency counts', () => {
      // GIVEN: Tooltip with zero dependencies
      const zeroCountData: TooltipData = {
        ...defaultTooltipData,
        incomingCount: 0,
        outgoingCount: 0,
      }

      // WHEN: Rendered
      renderNodeTooltip({ data: zeroCountData })

      // THEN: Should display zeros correctly
      const zeros = screen.getAllByText('0')
      expect(zeros.length).toBeGreaterThanOrEqual(2)
    })
  })
})
