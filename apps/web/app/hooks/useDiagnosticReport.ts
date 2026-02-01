import type { CircularDependencyInfo, DependencyGraph } from '@monoguard/types'
import { useCallback, useState } from 'react'
import {
  exportDiagnosticReportAsHtml,
  generateDiagnosticReport,
} from '../lib/diagnostics/generateDiagnosticReport'
import type { DiagnosticReport } from '../lib/diagnostics/types'

export interface UseDiagnosticReportOptions {
  graph: DependencyGraph
  allCycles: CircularDependencyInfo[]
  totalPackages: number
  projectName: string
}

export interface DiagnosticReportState {
  report: DiagnosticReport | null
  isGenerating: boolean
  isModalOpen: boolean
  error: string | null
}

export interface UseDiagnosticReportResult {
  state: DiagnosticReportState
  generateReport: (cycle: CircularDependencyInfo) => void
  exportAsHtml: () => void
  closeModal: () => void
}

/**
 * Hook for managing diagnostic report generation and modal state
 * AC: Report generation, HTML export, and modal lifecycle
 */
export function useDiagnosticReport(
  options: UseDiagnosticReportOptions
): UseDiagnosticReportResult {
  const [state, setState] = useState<DiagnosticReportState>({
    report: null,
    isGenerating: false,
    isModalOpen: false,
    error: null,
  })

  const generateReport = useCallback(
    (cycle: CircularDependencyInfo) => {
      setState((prev) => ({
        ...prev,
        isGenerating: true,
        isModalOpen: true,
        error: null,
      }))

      // Defer computation to allow React to render the loading state first
      setTimeout(() => {
        try {
          const report = generateDiagnosticReport({
            cycle,
            graph: options.graph,
            allCycles: options.allCycles,
            totalPackages: options.totalPackages,
            projectName: options.projectName,
          })

          setState({
            report,
            isGenerating: false,
            isModalOpen: true,
            error: null,
          })
        } catch (err) {
          setState((prev) => ({
            ...prev,
            isGenerating: false,
            isModalOpen: false,
            error: err instanceof Error ? err.message : 'Failed to generate report',
          }))
        }
      }, 0)
    },
    [options.graph, options.allCycles, options.totalPackages, options.projectName]
  )

  const exportAsHtml = useCallback(() => {
    if (!state.report) return

    try {
      const { blob, filename } = exportDiagnosticReportAsHtml(state.report)
      downloadBlob(blob, filename)
    } catch (err) {
      setState((prev) => ({
        ...prev,
        error: err instanceof Error ? err.message : 'Failed to export report',
      }))
    }
  }, [state.report])

  const closeModal = useCallback(() => {
    setState((prev) => ({
      ...prev,
      isModalOpen: false,
    }))
  }, [])

  return {
    state,
    generateReport,
    exportAsHtml,
    closeModal,
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
