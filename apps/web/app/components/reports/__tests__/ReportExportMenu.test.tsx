import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import type { ExportProgress } from '../../../hooks/useReportExport'
import { ReportExportMenu } from '../ReportExportMenu'

const defaultProgress: ExportProgress = {
  isExporting: false,
  progress: 0,
  stage: 'preparing',
}

describe('ReportExportMenu', () => {
  it('should not render when closed', () => {
    render(
      <ReportExportMenu
        isOpen={false}
        onClose={vi.fn()}
        onExport={vi.fn()}
        exportProgress={defaultProgress}
      />
    )

    expect(screen.queryByText('Export Report')).not.toBeInTheDocument()
  })

  it('should render when open', () => {
    render(
      <ReportExportMenu
        isOpen={true}
        onClose={vi.fn()}
        onExport={vi.fn()}
        exportProgress={defaultProgress}
      />
    )

    expect(screen.getByText('Export Report')).toBeInTheDocument()
  })

  it('should render format options', () => {
    render(
      <ReportExportMenu
        isOpen={true}
        onClose={vi.fn()}
        onExport={vi.fn()}
        exportProgress={defaultProgress}
      />
    )

    expect(screen.getByTestId('format-json')).toBeInTheDocument()
    expect(screen.getByTestId('format-html')).toBeInTheDocument()
    expect(screen.getByTestId('format-markdown')).toBeInTheDocument()
  })

  it('should render section checkboxes', () => {
    render(
      <ReportExportMenu
        isOpen={true}
        onClose={vi.fn()}
        onExport={vi.fn()}
        exportProgress={defaultProgress}
      />
    )

    expect(screen.getByText('Health Score Summary')).toBeInTheDocument()
    expect(screen.getByText('Circular Dependencies')).toBeInTheDocument()
    expect(screen.getByText('Version Conflicts')).toBeInTheDocument()
    expect(screen.getByText('Fix Recommendations')).toBeInTheDocument()
  })

  it('should toggle section checkboxes', () => {
    render(
      <ReportExportMenu
        isOpen={true}
        onClose={vi.fn()}
        onExport={vi.fn()}
        exportProgress={defaultProgress}
      />
    )

    const checkbox = screen.getByTestId('section-healthScore') as HTMLInputElement
    expect(checkbox.checked).toBe(true)

    fireEvent.click(checkbox)
    expect(checkbox.checked).toBe(false)
  })

  it('should call onExport with selected format and sections', async () => {
    const onExport = vi.fn().mockResolvedValue(undefined)
    render(
      <ReportExportMenu
        isOpen={true}
        onClose={vi.fn()}
        onExport={onExport}
        exportProgress={defaultProgress}
      />
    )

    // Select HTML format
    fireEvent.click(screen.getByTestId('format-html'))

    // Click export
    fireEvent.click(screen.getByTestId('export-button'))

    expect(onExport).toHaveBeenCalledWith(
      'html',
      expect.objectContaining({
        healthScore: true,
        circularDependencies: true,
      })
    )
  })

  it('should call onClose when backdrop clicked', () => {
    const onClose = vi.fn()
    render(
      <ReportExportMenu
        isOpen={true}
        onClose={onClose}
        onExport={vi.fn()}
        exportProgress={defaultProgress}
      />
    )

    // Click the backdrop (dialog overlay)
    const dialog = screen.getByRole('dialog')
    fireEvent.click(dialog)

    expect(onClose).toHaveBeenCalled()
  })

  it('should call onClose when close button clicked', () => {
    const onClose = vi.fn()
    render(
      <ReportExportMenu
        isOpen={true}
        onClose={onClose}
        onExport={vi.fn()}
        exportProgress={defaultProgress}
      />
    )

    fireEvent.click(screen.getByLabelText('Close export menu'))
    expect(onClose).toHaveBeenCalled()
  })

  it('should show progress bar when exporting', () => {
    const exportingProgress: ExportProgress = {
      isExporting: true,
      progress: 50,
      stage: 'generating',
    }

    render(
      <ReportExportMenu
        isOpen={true}
        onClose={vi.fn()}
        onExport={vi.fn()}
        exportProgress={exportingProgress}
      />
    )

    expect(screen.getByTestId('export-progress')).toBeInTheDocument()
    expect(screen.getByText('Generating report...')).toBeInTheDocument()
  })

  it('should disable export button while exporting', () => {
    const exportingProgress: ExportProgress = {
      isExporting: true,
      progress: 50,
      stage: 'generating',
    }

    render(
      <ReportExportMenu
        isOpen={true}
        onClose={vi.fn()}
        onExport={vi.fn()}
        exportProgress={exportingProgress}
      />
    )

    const button = screen.getByTestId('export-button')
    expect(button).toBeDisabled()
    expect(button.textContent).toBe('Exporting...')
  })

  it('should not propagate click from modal content', () => {
    const onClose = vi.fn()
    render(
      <ReportExportMenu
        isOpen={true}
        onClose={onClose}
        onExport={vi.fn()}
        exportProgress={defaultProgress}
      />
    )

    // Click on the modal content (not backdrop)
    fireEvent.click(screen.getByText('Export Report'))
    expect(onClose).not.toHaveBeenCalled()
  })
})
