/**
 * useViewportState - Shared viewport state for SVG and Canvas renderers (Story 4.9)
 *
 * Provides zoom, pan state that can be shared across render modes
 * so viewport is preserved when switching between SVG and Canvas.
 */

import { useCallback, useState } from 'react'
import type { ViewportState } from './types'
import { DEFAULT_VIEWPORT } from './types'

export interface UseViewportStateResult {
  viewport: ViewportState
  setViewport: React.Dispatch<React.SetStateAction<ViewportState>>
  resetViewport: () => void
  setZoom: (zoom: number) => void
  setPan: (panX: number, panY: number) => void
}

export function useViewportState(
  initialState: ViewportState = DEFAULT_VIEWPORT
): UseViewportStateResult {
  const [viewport, setViewport] = useState<ViewportState>(initialState)

  const resetViewport = useCallback(() => {
    setViewport(DEFAULT_VIEWPORT)
  }, [])

  const setZoom = useCallback((zoom: number) => {
    setViewport((prev) => ({
      ...prev,
      zoom: Math.max(0.1, Math.min(4, zoom)),
    }))
  }, [])

  const setPan = useCallback((panX: number, panY: number) => {
    setViewport((prev) => ({ ...prev, panX, panY }))
  }, [])

  return {
    viewport,
    setViewport,
    resetViewport,
    setZoom,
    setPan,
  }
}
