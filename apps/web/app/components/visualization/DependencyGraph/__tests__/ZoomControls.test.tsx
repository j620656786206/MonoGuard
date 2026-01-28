/**
 * Tests for ZoomControls component
 *
 * @see Story 4.4: Add Zoom, Pan, and Navigation Controls
 * @vitest-environment jsdom
 */
import { fireEvent, render, screen } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import type { ZoomControlsProps } from '../ZoomControls'
import { ZoomControls } from '../ZoomControls'

describe('ZoomControls', () => {
  const defaultProps: ZoomControlsProps = {
    zoomPercent: 100,
    onZoomIn: vi.fn(),
    onZoomOut: vi.fn(),
    onFitToScreen: vi.fn(),
    onResetZoom: vi.fn(),
    canZoomIn: true,
    canZoomOut: true,
  }

  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('rendering', () => {
    it('should render all control buttons', () => {
      render(<ZoomControls {...defaultProps} />)

      expect(screen.getByLabelText('Zoom out')).toBeInTheDocument()
      expect(screen.getByLabelText('Zoom in')).toBeInTheDocument()
      expect(screen.getByLabelText('Fit to screen')).toBeInTheDocument()
      expect(screen.getByLabelText('Reset zoom to 100%')).toBeInTheDocument()
    })

    it('should display current zoom percentage (AC6)', () => {
      render(<ZoomControls {...defaultProps} zoomPercent={150} />)

      expect(screen.getByText('150%')).toBeInTheDocument()
    })

    it('should update zoom display in real-time (AC6)', () => {
      const { rerender } = render(<ZoomControls {...defaultProps} zoomPercent={100} />)
      expect(screen.getByText('100%')).toBeInTheDocument()

      rerender(<ZoomControls {...defaultProps} zoomPercent={120} />)
      expect(screen.getByText('120%')).toBeInTheDocument()

      rerender(<ZoomControls {...defaultProps} zoomPercent={80} />)
      expect(screen.getByText('80%')).toBeInTheDocument()
    })

    it('should display 10% at minimum zoom', () => {
      render(<ZoomControls {...defaultProps} zoomPercent={10} />)

      expect(screen.getByText('10%')).toBeInTheDocument()
    })

    it('should display 400% at maximum zoom', () => {
      render(<ZoomControls {...defaultProps} zoomPercent={400} />)

      expect(screen.getByText('400%')).toBeInTheDocument()
    })
  })

  describe('zoom button interactions (AC3)', () => {
    it('should call onZoomIn when plus button is clicked', () => {
      render(<ZoomControls {...defaultProps} />)

      fireEvent.click(screen.getByLabelText('Zoom in'))

      expect(defaultProps.onZoomIn).toHaveBeenCalledTimes(1)
    })

    it('should call onZoomOut when minus button is clicked', () => {
      render(<ZoomControls {...defaultProps} />)

      fireEvent.click(screen.getByLabelText('Zoom out'))

      expect(defaultProps.onZoomOut).toHaveBeenCalledTimes(1)
    })

    it('should call onFitToScreen when fit button is clicked', () => {
      render(<ZoomControls {...defaultProps} />)

      fireEvent.click(screen.getByLabelText('Fit to screen'))

      expect(defaultProps.onFitToScreen).toHaveBeenCalledTimes(1)
    })

    it('should call onResetZoom when reset button is clicked', () => {
      render(<ZoomControls {...defaultProps} />)

      fireEvent.click(screen.getByLabelText('Reset zoom to 100%'))

      expect(defaultProps.onResetZoom).toHaveBeenCalledTimes(1)
    })
  })

  describe('button disabled states (AC3, AC7)', () => {
    it('should disable zoom in button when canZoomIn is false (at max limit)', () => {
      render(<ZoomControls {...defaultProps} canZoomIn={false} />)

      expect(screen.getByLabelText('Zoom in')).toBeDisabled()
    })

    it('should disable zoom out button when canZoomOut is false (at min limit)', () => {
      render(<ZoomControls {...defaultProps} canZoomOut={false} />)

      expect(screen.getByLabelText('Zoom out')).toBeDisabled()
    })

    it('should enable zoom in button when canZoomIn is true', () => {
      render(<ZoomControls {...defaultProps} canZoomIn={true} />)

      expect(screen.getByLabelText('Zoom in')).not.toBeDisabled()
    })

    it('should enable zoom out button when canZoomOut is true', () => {
      render(<ZoomControls {...defaultProps} canZoomOut={true} />)

      expect(screen.getByLabelText('Zoom out')).not.toBeDisabled()
    })

    it('should not call onZoomIn when button is disabled', () => {
      render(<ZoomControls {...defaultProps} canZoomIn={false} />)

      fireEvent.click(screen.getByLabelText('Zoom in'))

      expect(defaultProps.onZoomIn).not.toHaveBeenCalled()
    })

    it('should not call onZoomOut when button is disabled', () => {
      render(<ZoomControls {...defaultProps} canZoomOut={false} />)

      fireEvent.click(screen.getByLabelText('Zoom out'))

      expect(defaultProps.onZoomOut).not.toHaveBeenCalled()
    })
  })

  describe('accessibility', () => {
    it('should have accessible labels for all buttons', () => {
      render(<ZoomControls {...defaultProps} />)

      expect(screen.getByLabelText('Zoom out')).toBeInTheDocument()
      expect(screen.getByLabelText('Zoom in')).toBeInTheDocument()
      expect(screen.getByLabelText('Fit to screen')).toBeInTheDocument()
      expect(screen.getByLabelText('Reset zoom to 100%')).toBeInTheDocument()
    })

    it('should have title attributes for tooltip display', () => {
      render(<ZoomControls {...defaultProps} />)

      expect(screen.getByTitle('Zoom out')).toBeInTheDocument()
      expect(screen.getByTitle('Zoom in')).toBeInTheDocument()
      expect(screen.getByTitle('Fit to screen')).toBeInTheDocument()
      expect(screen.getByTitle('Reset zoom to 100%')).toBeInTheDocument()
    })

    it('should have proper button types', () => {
      render(<ZoomControls {...defaultProps} />)

      const buttons = screen.getAllByRole('button')
      buttons.forEach((button) => {
        expect(button).toHaveAttribute('type', 'button')
      })
    })
  })

  describe('fixed position (AC3)', () => {
    it('should have absolute positioning class', () => {
      const { container } = render(<ZoomControls {...defaultProps} />)

      const controlsDiv = container.firstChild
      expect(controlsDiv).toHaveClass('absolute')
    })

    it('should be positioned in bottom right', () => {
      const { container } = render(<ZoomControls {...defaultProps} />)

      const controlsDiv = container.firstChild
      expect(controlsDiv).toHaveClass('bottom-4')
      expect(controlsDiv).toHaveClass('right-4')
    })
  })
})
