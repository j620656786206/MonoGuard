/**
 * useRenderMode - Automatic render mode selection with user override (Story 4.9)
 *
 * Determines whether to use SVG or Canvas rendering based on:
 * 1. User preference from settings store (force-svg, force-canvas)
 * 2. Automatic threshold (500 nodes) when set to "auto"
 *
 * @see AC1: Automatic Render Mode Selection
 * @see AC3: User Override in Settings
 */

import { useMemo } from 'react'
import { useSettingsStore } from '../../../stores/settings'
import type { RenderMode } from './types'
import { NODE_THRESHOLD } from './types'

export interface UseRenderModeResult {
  /** The determined render mode */
  mode: RenderMode
  /** Whether auto mode is active */
  isAutoMode: boolean
  /** Whether the mode is forced by user preference */
  isForced: boolean
  /** Whether a performance warning should be shown */
  shouldShowWarning: boolean
  /** Warning message text, if applicable */
  warningMessage: string | null
}

export function useRenderMode(nodeCount: number): UseRenderModeResult {
  const visualizationMode = useSettingsStore((state) => state.visualizationMode)

  return useMemo(() => {
    const isAutoMode = visualizationMode === 'auto'
    let mode: RenderMode
    let shouldShowWarning = false
    let warningMessage: string | null = null

    if (visualizationMode === 'force-svg') {
      mode = 'svg'
      if (nodeCount >= NODE_THRESHOLD) {
        shouldShowWarning = true
        warningMessage = `SVG mode may be slow with ${nodeCount} nodes. Consider using Auto mode.`
      }
    } else if (visualizationMode === 'force-canvas') {
      mode = 'canvas'
    } else {
      // Auto mode
      mode = nodeCount >= NODE_THRESHOLD ? 'canvas' : 'svg'
    }

    return {
      mode,
      isAutoMode,
      isForced: !isAutoMode,
      shouldShowWarning,
      warningMessage,
    }
  }, [nodeCount, visualizationMode])
}
