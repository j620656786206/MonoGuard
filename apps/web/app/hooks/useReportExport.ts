import { useCallback, useState } from 'react'
import { generateHtmlReport } from '../lib/reports/generateHtmlReport'
import { generateJsonReport } from '../lib/reports/generateJsonReport'
import { generateMarkdownReport } from '../lib/reports/generateMarkdownReport'
import type { ReportData, ReportOptions, ReportResult } from '../lib/reports/types'

export interface ExportProgress {
  isExporting: boolean
  progress: number
  stage: 'preparing' | 'generating' | 'complete'
}

export interface UseReportExportResult {
  exportProgress: ExportProgress
  startExport: (data: ReportData, options: ReportOptions) => Promise<void>
  cancelExport: () => void
}

/**
 * Hook for managing report export state and execution
 * AC10: Export progress indicator
 * AC11: UI remains responsive during export
 */
export function useReportExport(): UseReportExportResult {
  const [exportProgress, setExportProgress] = useState<ExportProgress>({
    isExporting: false,
    progress: 0,
    stage: 'preparing',
  })

  const startExport = useCallback(async (data: ReportData, options: ReportOptions) => {
    setExportProgress({
      isExporting: true,
      progress: 10,
      stage: 'preparing',
    })

    try {
      setExportProgress((prev) => ({
        ...prev,
        progress: 30,
        stage: 'generating',
      }))

      let result: ReportResult

      switch (options.format) {
        case 'json':
          result = generateJsonReport(data, options)
          break
        case 'html':
          result = generateHtmlReport(data, options)
          break
        case 'markdown':
          result = generateMarkdownReport(data, options)
          break
        default:
          throw new Error(`Unknown format: ${options.format}`)
      }

      setExportProgress((prev) => ({ ...prev, progress: 90 }))

      downloadBlob(result.blob, result.filename)

      setExportProgress({
        isExporting: false,
        progress: 100,
        stage: 'complete',
      })
    } catch (error) {
      setExportProgress({
        isExporting: false,
        progress: 0,
        stage: 'preparing',
      })
      throw error
    }
  }, [])

  const cancelExport = useCallback(() => {
    setExportProgress({
      isExporting: false,
      progress: 0,
      stage: 'preparing',
    })
  }, [])

  return {
    exportProgress,
    startExport,
    cancelExport,
  }
}

function downloadBlob(blob: Blob, filename: string) {
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}
