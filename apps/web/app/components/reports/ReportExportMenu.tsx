import type React from 'react'
import { useCallback, useState } from 'react'
import type { ExportProgress } from '../../hooks/useReportExport'
import type { ReportFormat, ReportSections } from '../../lib/reports/types'
import { DEFAULT_REPORT_SECTIONS } from '../../lib/reports/types'

export interface ReportExportMenuProps {
  isOpen: boolean
  onClose: () => void
  onExport: (format: ReportFormat, sections: ReportSections) => Promise<void>
  exportProgress: ExportProgress
}

const FORMAT_OPTIONS: { value: ReportFormat; label: string; description: string }[] = [
  { value: 'json', label: 'JSON', description: 'Machine-readable, complete data' },
  { value: 'html', label: 'HTML', description: 'Self-contained, styled report' },
  { value: 'markdown', label: 'Markdown', description: 'GFM-compatible text format' },
]

const SECTION_LABELS: Record<keyof ReportSections, string> = {
  healthScore: 'Health Score Summary',
  circularDependencies: 'Circular Dependencies',
  versionConflicts: 'Version Conflicts',
  fixRecommendations: 'Fix Recommendations',
  packageList: 'Package List',
  graphSummary: 'Graph Summary',
}

/**
 * ReportExportMenu - Format selection and section toggles for report export
 * AC10: Export button triggers menu with format/section options
 */
export const ReportExportMenu: React.FC<ReportExportMenuProps> = ({
  isOpen,
  onClose,
  onExport,
  exportProgress,
}) => {
  const [selectedFormat, setSelectedFormat] = useState<ReportFormat>('json')
  const [sections, setSections] = useState<ReportSections>({
    ...DEFAULT_REPORT_SECTIONS,
  })

  const toggleSection = useCallback((key: keyof ReportSections) => {
    setSections((prev) => ({ ...prev, [key]: !prev[key] }))
  }, [])

  const handleExport = useCallback(async () => {
    await onExport(selectedFormat, sections)
  }, [onExport, selectedFormat, sections])

  if (!isOpen) return null

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
      onClick={onClose}
      onKeyDown={(e) => {
        if (e.key === 'Escape') onClose()
      }}
      role="dialog"
      aria-modal="true"
      aria-label="Export Analysis Report"
    >
      <div
        className="w-full max-w-md rounded-lg bg-white p-6 shadow-xl dark:bg-gray-800"
        role="document"
        onClick={(e) => e.stopPropagation()}
        onKeyDown={(e) => e.stopPropagation()}
      >
        <div className="mb-4 flex items-center justify-between">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">Export Report</h2>
          <button
            type="button"
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
            aria-label="Close export menu"
          >
            <svg
              className="h-5 w-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              aria-hidden="true"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>

        {/* Format Selection */}
        <div className="mb-4">
          <span className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
            Format
          </span>
          <div className="grid grid-cols-3 gap-2">
            {FORMAT_OPTIONS.map((option) => (
              <button
                type="button"
                key={option.value}
                onClick={() => setSelectedFormat(option.value)}
                className={`rounded-md border px-3 py-2 text-sm font-medium transition-colors ${
                  selectedFormat === option.value
                    ? 'border-blue-500 bg-blue-50 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300'
                    : 'border-gray-300 bg-white text-gray-700 hover:bg-gray-50 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-300'
                }`}
                title={option.description}
                data-testid={`format-${option.value}`}
              >
                {option.label}
              </button>
            ))}
          </div>
        </div>

        {/* Section Toggles */}
        <div className="mb-6">
          <span className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
            Sections
          </span>
          <div className="space-y-2">
            {(Object.keys(SECTION_LABELS) as Array<keyof ReportSections>).map((key) => (
              <label key={key} className="flex cursor-pointer items-center gap-2">
                <input
                  type="checkbox"
                  checked={sections[key]}
                  onChange={() => toggleSection(key)}
                  className="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  data-testid={`section-${key}`}
                />
                <span className="text-sm text-gray-700 dark:text-gray-300">
                  {SECTION_LABELS[key]}
                </span>
              </label>
            ))}
          </div>
        </div>

        {/* Progress Bar */}
        {exportProgress.isExporting && (
          <div className="mb-4" data-testid="export-progress">
            <div className="h-2 overflow-hidden rounded-full bg-gray-200 dark:bg-gray-600">
              <div
                className="h-full rounded-full bg-blue-500 transition-all duration-300"
                style={{ width: `${exportProgress.progress}%` }}
              />
            </div>
            <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
              {exportProgress.stage === 'preparing'
                ? 'Preparing...'
                : exportProgress.stage === 'generating'
                  ? 'Generating report...'
                  : 'Complete!'}
            </p>
          </div>
        )}

        {/* Actions */}
        <div className="flex justify-end gap-2">
          <button
            type="button"
            onClick={onClose}
            className="rounded-md border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-700"
          >
            Cancel
          </button>
          <button
            type="button"
            onClick={handleExport}
            disabled={exportProgress.isExporting}
            className="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
            data-testid="export-button"
          >
            {exportProgress.isExporting ? 'Exporting...' : 'Export'}
          </button>
        </div>
      </div>
    </div>
  )
}
