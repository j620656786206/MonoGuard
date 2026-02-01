/**
 * RenderModeIndicator - Displays current rendering mode and node count (Story 4.9)
 *
 * Shows a non-intrusive indicator in the top-right corner with:
 * - Current node count
 * - Active render mode (SVG/Canvas)
 * - "Forced" badge when user override is active
 *
 * @see AC2: Mode Indicator Display
 */

import React from 'react'
import type { RenderMode } from './types'

export interface RenderModeIndicatorProps {
  /** Current rendering mode */
  mode: RenderMode
  /** Total number of nodes in the graph */
  nodeCount: number
  /** Whether the mode is forced by user preference */
  isForced: boolean
  /** Optional performance warning message (AC3) */
  warningMessage?: string | null
}

export const RenderModeIndicator = React.memo(function RenderModeIndicator({
  mode,
  nodeCount,
  isForced,
  warningMessage,
}: RenderModeIndicatorProps) {
  return (
    <output
      className="absolute top-2 right-2 flex flex-col items-end gap-1"
      aria-label={`${nodeCount} nodes, ${mode.toUpperCase()} rendering mode${isForced ? ', forced' : ''}`}
    >
      <span className="flex items-center gap-2 rounded bg-gray-100 px-2 py-1 text-xs text-gray-600 dark:bg-gray-800 dark:text-gray-400">
        <span>{nodeCount} nodes</span>
        <span className="text-gray-400 dark:text-gray-600" aria-hidden="true">
          &bull;
        </span>
        <span className={mode === 'canvas' ? 'text-amber-600' : 'text-blue-600'}>
          {mode.toUpperCase()} mode
        </span>
        {isForced && (
          <>
            <span className="text-gray-400 dark:text-gray-600" aria-hidden="true">
              &bull;
            </span>
            <span className="text-orange-500">Forced</span>
          </>
        )}
      </span>
      {warningMessage && (
        <span
          className="rounded bg-amber-50 px-2 py-1 text-xs text-amber-700 dark:bg-amber-900/30 dark:text-amber-400"
          role="alert"
        >
          {warningMessage}
        </span>
      )}
    </output>
  )
})
