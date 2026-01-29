/**
 * useGraphExport - Hook for managing graph export operations
 *
 * Handles export state, progress tracking, format selection,
 * and file download triggering.
 *
 * @see Story 4.6: Export Graph as PNG/SVG Images
 */

import { useCallback, useRef, useState } from 'react'
import type { ExportOptions, ExportProgress, ExportResult } from './types'
import { exportPng } from './utils/exportPng'
import { exportSvg } from './utils/exportSvg'
import { renderLegendSvg } from './utils/renderLegendForExport'

export interface UseGraphExportProps {
  /** Ref to the SVG element containing the graph */
  svgRef: React.RefObject<SVGSVGElement | null>
  /** Project name for filename generation */
  projectName: string
  /** Whether dark mode is active */
  isDarkMode: boolean
}

export interface UseGraphExportResult {
  /** Current export progress state */
  exportProgress: ExportProgress
  /** Start an export with given options */
  startExport: (options: ExportOptions) => Promise<void>
  /** Cancel an in-progress export */
  cancelExport: () => void
}

export function useGraphExport({
  svgRef,
  projectName,
  isDarkMode,
}: UseGraphExportProps): UseGraphExportResult {
  const [exportProgress, setExportProgress] = useState<ExportProgress>({
    isExporting: false,
    progress: 0,
    stage: 'preparing',
  })

  const abortControllerRef = useRef<AbortController | null>(null)

  const startExport = useCallback(
    async (options: ExportOptions) => {
      if (!svgRef.current) {
        throw new Error('SVG element not found')
      }

      const controller = new AbortController()
      abortControllerRef.current = controller

      setExportProgress({
        isExporting: true,
        progress: 0,
        stage: 'preparing',
      })

      try {
        // Generate legend SVG if needed
        let legendSvg: string | undefined
        if (options.includeLegend) {
          legendSvg = renderLegendSvg(isDarkMode)
        }

        // Determine background color based on theme
        const backgroundColor =
          options.backgroundColor === 'transparent'
            ? 'transparent'
            : isDarkMode
              ? '#1f2937'
              : '#ffffff'

        const exportOptions: ExportOptions = {
          ...options,
          backgroundColor,
        }

        setExportProgress((prev) => ({ ...prev, progress: 20, stage: 'rendering' }))

        let result: ExportResult

        if (options.format === 'svg') {
          result = await exportSvg({
            svgElement: svgRef.current,
            options: exportOptions,
            projectName,
            legendSvg,
          })
        } else {
          result = await exportPng({
            svgElement: svgRef.current,
            options: exportOptions,
            projectName,
            legendSvg,
            onProgress: (progress) => {
              setExportProgress((prev) => ({
                ...prev,
                progress: 20 + progress * 0.7,
                stage: progress < 50 ? 'rendering' : 'encoding',
              }))
            },
          })
        }

        // Check if cancelled
        if (controller.signal.aborted) {
          return
        }

        setExportProgress((prev) => ({ ...prev, progress: 95, stage: 'complete' }))

        // Trigger file download
        downloadBlob(result.blob, result.filename)

        setExportProgress({
          isExporting: false,
          progress: 100,
          stage: 'complete',
        })
      } catch (error) {
        if (!controller.signal.aborted) {
          console.error('Export failed:', error)
          setExportProgress({
            isExporting: false,
            progress: 0,
            stage: 'preparing',
          })
          throw error
        }
      } finally {
        abortControllerRef.current = null
      }
    },
    [svgRef, projectName, isDarkMode]
  )

  const cancelExport = useCallback(() => {
    abortControllerRef.current?.abort()
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

/**
 * Trigger a file download from a Blob.
 * Creates a temporary anchor element and clicks it.
 */
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
