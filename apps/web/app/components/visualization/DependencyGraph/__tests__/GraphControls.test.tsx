/**
 * Tests for GraphControls component
 *
 * @see Story 4.3: Implement Node Expand/Collapse Functionality (AC3)
 *
 * Following Given-When-Then format with priority tags.
 */

import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { GraphControls, type GraphControlsProps } from '../GraphControls'

describe('GraphControls', () => {
  const defaultProps: GraphControlsProps = {
    currentDepth: 'all',
    maxDepth: 3,
    onDepthChange: vi.fn(),
    onExpandAll: vi.fn(),
    onCollapseAll: vi.fn(),
  }

  const renderGraphControls = (props: Partial<GraphControlsProps> = {}) => {
    const mergedProps = { ...defaultProps, ...props }
    return render(<GraphControls {...mergedProps} />)
  }

  describe('Rendering (AC3)', () => {
    it('[P1] should render Depth Control heading', () => {
      // GIVEN: GraphControls component
      // WHEN: Rendered
      renderGraphControls()

      // THEN: Should display Depth Control heading
      expect(screen.getByText('Depth Control')).toBeInTheDocument()
    })

    it('[P1] should render All button', () => {
      // GIVEN: GraphControls component
      // WHEN: Rendered
      renderGraphControls()

      // THEN: Should have "All" button with correct aria-label
      const allButton = screen.getByRole('button', { name: /show all depths/i })
      expect(allButton).toBeInTheDocument()
      expect(allButton).toHaveTextContent('All')
    })

    it('[P1] should render depth level buttons based on maxDepth', () => {
      // GIVEN: maxDepth of 3
      // WHEN: Rendered
      renderGraphControls({ maxDepth: 3 })

      // THEN: Should have L1, L2, L3 buttons
      expect(screen.getByRole('button', { name: /show depth level 1/i })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /show depth level 2/i })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /show depth level 3/i })).toBeInTheDocument()
    })

    it('[P1] should render Expand All button', () => {
      // GIVEN: GraphControls component
      // WHEN: Rendered
      renderGraphControls()

      // THEN: Should have Expand All button
      const expandButton = screen.getByRole('button', { name: /expand all nodes/i })
      expect(expandButton).toBeInTheDocument()
      expect(expandButton).toHaveTextContent('Expand All')
    })

    it('[P1] should render Collapse All button', () => {
      // GIVEN: GraphControls component
      // WHEN: Rendered
      renderGraphControls()

      // THEN: Should have Collapse All button
      const collapseButton = screen.getByRole('button', { name: /collapse all nodes/i })
      expect(collapseButton).toBeInTheDocument()
      expect(collapseButton).toHaveTextContent('Collapse All')
    })

    it('[P2] should limit depth buttons to max of 5', () => {
      // GIVEN: maxDepth greater than 5
      // WHEN: Rendered with maxDepth of 10
      renderGraphControls({ maxDepth: 10 })

      // THEN: Should only render L1-L5 buttons (plus All)
      expect(screen.getByRole('button', { name: /show depth level 1/i })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /show depth level 5/i })).toBeInTheDocument()
      expect(screen.queryByRole('button', { name: /show depth level 6/i })).not.toBeInTheDocument()
    })
  })

  describe('Current Depth Indication (AC3)', () => {
    it('[P1] should highlight "All" button when currentDepth is all', () => {
      // GIVEN: currentDepth is "all"
      // WHEN: Rendered
      renderGraphControls({ currentDepth: 'all' })

      // THEN: All button should have aria-pressed="true"
      const allButton = screen.getByRole('button', { name: /show all depths/i })
      expect(allButton).toHaveAttribute('aria-pressed', 'true')
    })

    it('[P1] should highlight depth button when currentDepth is a number', () => {
      // GIVEN: currentDepth is 2
      // WHEN: Rendered
      renderGraphControls({ currentDepth: 2 })

      // THEN: L2 button should have aria-pressed="true", others false
      const l2Button = screen.getByRole('button', { name: /show depth level 2/i })
      const allButton = screen.getByRole('button', { name: /show all depths/i })

      expect(l2Button).toHaveAttribute('aria-pressed', 'true')
      expect(allButton).toHaveAttribute('aria-pressed', 'false')
    })

    it('[P2] should apply different styles to selected button', () => {
      // GIVEN: currentDepth is 1
      // WHEN: Rendered
      renderGraphControls({ currentDepth: 1 })

      // THEN: L1 button should have indigo styling (selected)
      const l1Button = screen.getByRole('button', { name: /show depth level 1/i })
      expect(l1Button).toHaveClass('bg-indigo-600')
      expect(l1Button).toHaveClass('text-white')
    })
  })

  describe('User Interactions (AC3)', () => {
    it('[P1] should call onDepthChange with "all" when All button clicked', () => {
      // GIVEN: GraphControls with mocked callback
      const onDepthChange = vi.fn()
      renderGraphControls({ onDepthChange, currentDepth: 1 })

      // WHEN: All button is clicked
      fireEvent.click(screen.getByRole('button', { name: /show all depths/i }))

      // THEN: onDepthChange should be called with "all"
      expect(onDepthChange).toHaveBeenCalledWith('all')
      expect(onDepthChange).toHaveBeenCalledTimes(1)
    })

    it('[P1] should call onDepthChange with number when depth button clicked', () => {
      // GIVEN: GraphControls with mocked callback
      const onDepthChange = vi.fn()
      renderGraphControls({ onDepthChange, currentDepth: 'all', maxDepth: 3 })

      // WHEN: L2 button is clicked
      fireEvent.click(screen.getByRole('button', { name: /show depth level 2/i }))

      // THEN: onDepthChange should be called with 2
      expect(onDepthChange).toHaveBeenCalledWith(2)
    })

    it('[P1] should call onExpandAll when Expand All button clicked', () => {
      // GIVEN: GraphControls with mocked callback
      const onExpandAll = vi.fn()
      renderGraphControls({ onExpandAll })

      // WHEN: Expand All button is clicked
      fireEvent.click(screen.getByRole('button', { name: /expand all nodes/i }))

      // THEN: onExpandAll should be called
      expect(onExpandAll).toHaveBeenCalledTimes(1)
    })

    it('[P1] should call onCollapseAll when Collapse All button clicked', () => {
      // GIVEN: GraphControls with mocked callback
      const onCollapseAll = vi.fn()
      renderGraphControls({ onCollapseAll })

      // WHEN: Collapse All button is clicked
      fireEvent.click(screen.getByRole('button', { name: /collapse all nodes/i }))

      // THEN: onCollapseAll should be called
      expect(onCollapseAll).toHaveBeenCalledTimes(1)
    })
  })

  describe('Accessibility (AC3)', () => {
    it('[P1] should have role="group" with aria-label', () => {
      // GIVEN: GraphControls component
      // WHEN: Rendered
      renderGraphControls()

      // THEN: Should have group role with appropriate label
      const group = screen.getByRole('group', { name: /graph depth controls/i })
      expect(group).toBeInTheDocument()
    })

    it('[P2] should have type="button" on all buttons', () => {
      // GIVEN: GraphControls component
      // WHEN: Rendered
      renderGraphControls({ maxDepth: 3 })

      // THEN: All buttons should have type="button" to prevent form submission
      const buttons = screen.getAllByRole('button')
      buttons.forEach((button) => {
        expect(button).toHaveAttribute('type', 'button')
      })
    })

    it('[P2] should have descriptive aria-labels on all buttons', () => {
      // GIVEN: GraphControls component
      // WHEN: Rendered
      renderGraphControls({ maxDepth: 2 })

      // THEN: All buttons should have aria-label
      expect(screen.getByRole('button', { name: /show all depths/i })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /show depth level 1/i })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /show depth level 2/i })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /expand all nodes/i })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /collapse all nodes/i })).toBeInTheDocument()
    })
  })

  describe('Edge Cases', () => {
    it('[P2] should handle maxDepth of 0', () => {
      // GIVEN: maxDepth of 0
      // WHEN: Rendered
      renderGraphControls({ maxDepth: 0 })

      // THEN: Should only render All button (no depth levels)
      expect(screen.getByRole('button', { name: /show all depths/i })).toBeInTheDocument()
      expect(screen.queryByRole('button', { name: /show depth level/i })).not.toBeInTheDocument()
    })

    it('[P2] should handle maxDepth of 1', () => {
      // GIVEN: maxDepth of 1
      // WHEN: Rendered
      renderGraphControls({ maxDepth: 1 })

      // THEN: Should render All and L1 buttons only
      expect(screen.getByRole('button', { name: /show all depths/i })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /show depth level 1/i })).toBeInTheDocument()
      expect(screen.queryByRole('button', { name: /show depth level 2/i })).not.toBeInTheDocument()
    })
  })
})
