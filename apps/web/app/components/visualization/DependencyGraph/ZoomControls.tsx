/**
 * ZoomControls - Zoom control buttons and display for dependency graph
 *
 * Provides zoom in/out buttons, fit-to-screen, reset, and zoom level display.
 * Positioned in bottom-right corner of graph container (AC3: fixed position).
 *
 * @see Story 4.4: Add Zoom, Pan, and Navigation Controls
 */
'use client'

/**
 * Props for ZoomControls component
 */
export interface ZoomControlsProps {
  /** Current zoom level as percentage (e.g., 100 for 100%) */
  zoomPercent: number
  /** Callback when zoom in button is clicked */
  onZoomIn: () => void
  /** Callback when zoom out button is clicked */
  onZoomOut: () => void
  /** Callback when fit-to-screen button is clicked */
  onFitToScreen: () => void
  /** Callback when reset zoom button is clicked */
  onResetZoom: () => void
  /** Whether zoom in is allowed (not at max limit) */
  canZoomIn: boolean
  /** Whether zoom out is allowed (not at min limit) */
  canZoomOut: boolean
}

/**
 * ZoomControls component
 *
 * Renders zoom control buttons and current zoom percentage display.
 * Buttons are always visible in a fixed position (AC3).
 * Zoom display updates in real-time (AC6).
 */
export function ZoomControls({
  zoomPercent,
  onZoomIn,
  onZoomOut,
  onFitToScreen,
  onResetZoom,
  canZoomIn,
  canZoomOut,
}: ZoomControlsProps) {
  return (
    <div
      className="absolute bottom-4 right-4 bg-white/90 dark:bg-gray-800/90
                  rounded-lg shadow-lg p-2 flex items-center gap-1 z-10"
    >
      {/* Zoom Out Button */}
      <button
        type="button"
        onClick={onZoomOut}
        disabled={!canZoomOut}
        className="w-8 h-8 flex items-center justify-center rounded
                   hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors
                   disabled:opacity-40 disabled:cursor-not-allowed"
        aria-label="Zoom out"
        title="Zoom out"
      >
        <svg
          className="w-4 h-4"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          aria-hidden="true"
        >
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 12H4" />
        </svg>
      </button>

      {/* Zoom Level Display (AC6) */}
      <div className="w-14 text-center text-sm font-medium tabular-nums">{zoomPercent}%</div>

      {/* Zoom In Button */}
      <button
        type="button"
        onClick={onZoomIn}
        disabled={!canZoomIn}
        className="w-8 h-8 flex items-center justify-center rounded
                   hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors
                   disabled:opacity-40 disabled:cursor-not-allowed"
        aria-label="Zoom in"
        title="Zoom in"
      >
        <svg
          className="w-4 h-4"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          aria-hidden="true"
        >
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
        </svg>
      </button>

      {/* Divider */}
      <div className="w-px h-6 bg-gray-200 dark:bg-gray-700 mx-1" />

      {/* Fit to Screen Button (AC4) */}
      <button
        type="button"
        onClick={onFitToScreen}
        className="w-8 h-8 flex items-center justify-center rounded
                   hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
        aria-label="Fit to screen"
        title="Fit to screen"
      >
        <svg
          className="w-4 h-4"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          aria-hidden="true"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5v-4m0 4h-4m4 0l-5-5"
          />
        </svg>
      </button>

      {/* Reset Zoom Button */}
      <button
        type="button"
        onClick={onResetZoom}
        className="w-8 h-8 flex items-center justify-center rounded
                   hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
        aria-label="Reset zoom to 100%"
        title="Reset zoom to 100%"
      >
        <svg
          className="w-4 h-4"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          aria-hidden="true"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
          />
        </svg>
      </button>
    </div>
  )
}
