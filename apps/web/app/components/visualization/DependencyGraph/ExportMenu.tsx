/**
 * ExportMenu - Export configuration dropdown for dependency graph
 *
 * Provides format selection (PNG/SVG), resolution, scope, and option toggles.
 * Shows export progress during the export process.
 *
 * @see Story 4.6: Export Graph as PNG/SVG Images - AC1, AC7
 */

import React, { useState } from 'react'
import type {
  ExportFormat,
  ExportOptions,
  ExportProgress,
  ExportResolution,
  ExportScope,
} from './types'

export interface ExportMenuProps {
  /** Whether the menu is open */
  isOpen: boolean
  /** Callback to close the menu */
  onClose: () => void
  /** Callback to trigger export with given options */
  onExport: (options: ExportOptions) => Promise<void>
  /** Current export progress */
  exportProgress: ExportProgress
  /** Whether dark mode is active */
  isDarkMode: boolean
}

const DEFAULT_OPTIONS: ExportOptions = {
  format: 'png',
  scope: 'viewport',
  resolution: 2,
  includeLegend: true,
  includeWatermark: false,
  backgroundColor: '#ffffff',
}

export const ExportMenu = React.memo(function ExportMenu({
  isOpen,
  onClose,
  onExport,
  exportProgress,
  isDarkMode,
}: ExportMenuProps) {
  const [options, setOptions] = useState<ExportOptions>({
    ...DEFAULT_OPTIONS,
    backgroundColor: isDarkMode ? '#1f2937' : '#ffffff',
  })

  const handleExport = async () => {
    try {
      await onExport(options)
      onClose()
    } catch (error) {
      console.error('Export failed:', error)
    }
  }

  if (!isOpen) return null

  return (
    <div
      className="absolute right-4 top-14 z-50 w-72 rounded-lg border border-gray-200 bg-white p-4 shadow-xl dark:border-gray-700 dark:bg-gray-800"
      role="dialog"
      aria-label="Export Graph"
    >
      <div className="mb-4 flex items-center justify-between">
        <h3 className="font-semibold text-gray-900 dark:text-white">Export Graph</h3>
        <button
          type="button"
          onClick={onClose}
          className="text-gray-500 hover:text-gray-700 dark:hover:text-gray-300"
          aria-label="Close export menu"
        >
          <svg
            className="h-5 w-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            role="img"
            aria-label="Close"
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
        <div className="flex gap-2">
          {(['png', 'svg'] as ExportFormat[]).map((format) => (
            <button
              type="button"
              key={format}
              onClick={() => setOptions((prev) => ({ ...prev, format }))}
              className={`flex-1 rounded-md px-3 py-2 text-sm font-medium transition-colors ${
                options.format === format
                  ? 'bg-blue-600 text-white'
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-300 dark:hover:bg-gray-600'
              }`}
            >
              {format.toUpperCase()}
            </button>
          ))}
        </div>
      </div>

      {/* Resolution (PNG only) */}
      {options.format === 'png' && (
        <div className="mb-4">
          <label
            htmlFor="export-resolution"
            className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300"
          >
            Resolution
          </label>
          <select
            id="export-resolution"
            value={options.resolution}
            onChange={(e) =>
              setOptions((prev) => ({
                ...prev,
                resolution: parseInt(e.target.value, 10) as ExportResolution,
              }))
            }
            className="w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-gray-900 dark:border-gray-600 dark:bg-gray-700 dark:text-white"
          >
            <option value={1}>1x (Standard)</option>
            <option value={2}>2x (High DPI)</option>
            <option value={4}>4x (Print Quality)</option>
          </select>
        </div>
      )}

      {/* Scope Selection */}
      <div className="mb-4">
        <label
          htmlFor="export-scope"
          className="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300"
        >
          Scope
        </label>
        <select
          id="export-scope"
          value={options.scope}
          onChange={(e) =>
            setOptions((prev) => ({
              ...prev,
              scope: e.target.value as ExportScope,
            }))
          }
          className="w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-gray-900 dark:border-gray-600 dark:bg-gray-700 dark:text-white"
        >
          <option value="viewport">Current View</option>
          <option value="full">Full Graph</option>
          <option value="selected">Selected Elements</option>
        </select>
      </div>

      {/* Option Checkboxes */}
      <div className="mb-4 space-y-2">
        <label className="flex cursor-pointer items-center gap-2">
          <input
            type="checkbox"
            checked={options.includeLegend}
            onChange={(e) =>
              setOptions((prev) => ({
                ...prev,
                includeLegend: e.target.checked,
              }))
            }
            className="rounded border-gray-300 text-blue-600 focus:ring-blue-500 dark:border-gray-600"
          />
          <span className="text-sm text-gray-700 dark:text-gray-300">Include Legend</span>
        </label>

        <label className="flex cursor-pointer items-center gap-2">
          <input
            type="checkbox"
            checked={options.includeWatermark}
            onChange={(e) =>
              setOptions((prev) => ({
                ...prev,
                includeWatermark: e.target.checked,
              }))
            }
            className="rounded border-gray-300 text-blue-600 focus:ring-blue-500 dark:border-gray-600"
          />
          <span className="text-sm text-gray-700 dark:text-gray-300">Include Watermark</span>
        </label>
      </div>

      {/* Export Progress */}
      {exportProgress.isExporting && (
        <div className="mb-4">
          <div className="mb-1 flex justify-between text-sm text-gray-600 dark:text-gray-400">
            <span>Exporting...</span>
            <span>{Math.round(exportProgress.progress)}%</span>
          </div>
          <div className="h-2 w-full rounded-full bg-gray-200 dark:bg-gray-700">
            <div
              className="h-2 rounded-full bg-blue-600 transition-all duration-200"
              style={{ width: `${exportProgress.progress}%` }}
            />
          </div>
          <div className="mt-1 text-xs capitalize text-gray-500 dark:text-gray-400">
            {exportProgress.stage}
          </div>
        </div>
      )}

      {/* Export Button */}
      <button
        type="button"
        onClick={handleExport}
        disabled={exportProgress.isExporting}
        className="w-full rounded-md bg-blue-600 px-4 py-2 font-medium text-white transition-colors hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:cursor-not-allowed disabled:opacity-50"
      >
        {exportProgress.isExporting ? 'Exporting...' : 'Export'}
      </button>
    </div>
  )
})
