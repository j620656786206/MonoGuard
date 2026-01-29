/**
 * Tests for ExportMenu component
 *
 * @see Story 4.6: Export Graph as PNG/SVG Images - AC1, AC7
 */
import { fireEvent, render, screen } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { ExportMenu } from '../ExportMenu'
import type { ExportProgress } from '../types'

describe('ExportMenu', () => {
  const mockOnClose = vi.fn()
  const mockOnExport = vi.fn().mockResolvedValue(undefined)
  const defaultProgress: ExportProgress = {
    isExporting: false,
    progress: 0,
    stage: 'preparing',
  }

  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('should not render when isOpen is false', () => {
    render(
      <ExportMenu
        isOpen={false}
        onClose={mockOnClose}
        onExport={mockOnExport}
        exportProgress={defaultProgress}
        isDarkMode={false}
      />
    )

    expect(screen.queryByText('Export Graph')).not.toBeInTheDocument()
  })

  it('should render when isOpen is true', () => {
    render(
      <ExportMenu
        isOpen={true}
        onClose={mockOnClose}
        onExport={mockOnExport}
        exportProgress={defaultProgress}
        isDarkMode={false}
      />
    )

    expect(screen.getByText('Export Graph')).toBeInTheDocument()
  })

  it('should show format options (PNG and SVG)', () => {
    render(
      <ExportMenu
        isOpen={true}
        onClose={mockOnClose}
        onExport={mockOnExport}
        exportProgress={defaultProgress}
        isDarkMode={false}
      />
    )

    expect(screen.getByText('PNG')).toBeInTheDocument()
    expect(screen.getByText('SVG')).toBeInTheDocument()
  })

  it('should show resolution options for PNG format', () => {
    render(
      <ExportMenu
        isOpen={true}
        onClose={mockOnClose}
        onExport={mockOnExport}
        exportProgress={defaultProgress}
        isDarkMode={false}
      />
    )

    // PNG is default, resolution should be visible
    expect(screen.getByText('Resolution')).toBeInTheDocument()
  })

  it('should hide resolution options when SVG is selected', () => {
    render(
      <ExportMenu
        isOpen={true}
        onClose={mockOnClose}
        onExport={mockOnExport}
        exportProgress={defaultProgress}
        isDarkMode={false}
      />
    )

    fireEvent.click(screen.getByText('SVG'))

    expect(screen.queryByText('Resolution')).not.toBeInTheDocument()
  })

  it('should show scope selection options', () => {
    render(
      <ExportMenu
        isOpen={true}
        onClose={mockOnClose}
        onExport={mockOnExport}
        exportProgress={defaultProgress}
        isDarkMode={false}
      />
    )

    expect(screen.getByText('Scope')).toBeInTheDocument()
  })

  it('should show include legend checkbox', () => {
    render(
      <ExportMenu
        isOpen={true}
        onClose={mockOnClose}
        onExport={mockOnExport}
        exportProgress={defaultProgress}
        isDarkMode={false}
      />
    )

    expect(screen.getByText('Include Legend')).toBeInTheDocument()
  })

  it('should show include watermark checkbox', () => {
    render(
      <ExportMenu
        isOpen={true}
        onClose={mockOnClose}
        onExport={mockOnExport}
        exportProgress={defaultProgress}
        isDarkMode={false}
      />
    )

    expect(screen.getByText('Include Watermark')).toBeInTheDocument()
  })

  it('should call onExport with default options when Export button is clicked', async () => {
    render(
      <ExportMenu
        isOpen={true}
        onClose={mockOnClose}
        onExport={mockOnExport}
        exportProgress={defaultProgress}
        isDarkMode={false}
      />
    )

    fireEvent.click(screen.getByText('Export'))

    expect(mockOnExport).toHaveBeenCalledWith(
      expect.objectContaining({
        format: 'png',
        scope: 'viewport',
        resolution: 2,
        includeLegend: true,
        includeWatermark: false,
      })
    )
  })

  it('should call onClose when close button is clicked', () => {
    render(
      <ExportMenu
        isOpen={true}
        onClose={mockOnClose}
        onExport={mockOnExport}
        exportProgress={defaultProgress}
        isDarkMode={false}
      />
    )

    fireEvent.click(screen.getByLabelText('Close export menu'))

    expect(mockOnClose).toHaveBeenCalled()
  })

  it('should show progress bar during export', () => {
    render(
      <ExportMenu
        isOpen={true}
        onClose={mockOnClose}
        onExport={mockOnExport}
        exportProgress={{ isExporting: true, progress: 50, stage: 'rendering' }}
        isDarkMode={false}
      />
    )

    // "Exporting..." appears in both progress label and button
    const exportingTexts = screen.getAllByText('Exporting...')
    expect(exportingTexts.length).toBeGreaterThanOrEqual(1)
    expect(screen.getByText('50%')).toBeInTheDocument()
    expect(screen.getByText('rendering')).toBeInTheDocument()
  })

  it('should disable Export button during export', () => {
    render(
      <ExportMenu
        isOpen={true}
        onClose={mockOnClose}
        onExport={mockOnExport}
        exportProgress={{ isExporting: true, progress: 30, stage: 'rendering' }}
        isDarkMode={false}
      />
    )

    // Find the disabled button with "Exporting..." text
    const exportButton = screen.getByRole('button', { name: 'Exporting...' })
    expect(exportButton).toBeDisabled()
  })

  it('should have proper aria-label for accessibility', () => {
    render(
      <ExportMenu
        isOpen={true}
        onClose={mockOnClose}
        onExport={mockOnExport}
        exportProgress={defaultProgress}
        isDarkMode={false}
      />
    )

    expect(screen.getByRole('dialog')).toHaveAttribute('aria-label', 'Export Graph')
  })

  it('should toggle legend checkbox', () => {
    render(
      <ExportMenu
        isOpen={true}
        onClose={mockOnClose}
        onExport={mockOnExport}
        exportProgress={defaultProgress}
        isDarkMode={false}
      />
    )

    const legendCheckbox = screen.getByText('Include Legend')
      .previousElementSibling as HTMLInputElement

    // Default is checked
    expect(legendCheckbox.checked).toBe(true)

    fireEvent.click(legendCheckbox)
    expect(legendCheckbox.checked).toBe(false)
  })
})
