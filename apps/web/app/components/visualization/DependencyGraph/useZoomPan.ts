/**
 * useZoomPan - Custom hook for zoom and pan functionality
 *
 * Manages zoom/pan state for D3 visualizations with React state sync.
 * Provides API for programmatic zoom control (zoom in/out, fit to screen, reset).
 *
 * @see Story 4.4: Add Zoom, Pan, and Navigation Controls
 */
'use client'

import * as d3 from 'd3'
import { useCallback, useRef, useState } from 'react'

/**
 * Zoom and pan state
 */
export interface ZoomPanState {
  /** Current zoom scale (1 = 100%) */
  scale: number
  /** X translation */
  translateX: number
  /** Y translation */
  translateY: number
}

/**
 * D3 zoom transform-like object
 */
export interface ZoomTransform {
  k: number
  x: number
  y: number
}

/**
 * Props for useZoomPan hook
 */
export interface UseZoomPanProps {
  /** Reference to the SVG element */
  svgRef: React.RefObject<SVGSVGElement | null>
  /** Reference to the container group element */
  containerRef: React.RefObject<SVGGElement | null>
  /** Minimum scale (default: 0.1 = 10%) */
  minScale?: number
  /** Maximum scale (default: 4 = 400%) */
  maxScale?: number
  /** Zoom increment per button click (default: 0.2 = 20%) */
  zoomIncrement?: number
}

/**
 * Return type for useZoomPan hook
 */
export interface UseZoomPanResult {
  /** Current zoom/pan state */
  zoomState: ZoomPanState
  /** Current zoom as percentage (e.g., 100 for 100%) */
  zoomPercent: number
  /** Zoom in by increment amount */
  zoomIn: () => void
  /** Zoom out by increment amount */
  zoomOut: () => void
  /** Reset zoom to identity (scale 1, no translation) */
  resetZoom: () => void
  /** Fit entire graph to viewport */
  fitToScreen: () => void
  /** Store zoom behavior for external control */
  setZoomBehavior: (zoom: d3.ZoomBehavior<SVGSVGElement, unknown>) => void
  /** Handle zoom change from D3 */
  handleZoomChange: (transform: ZoomTransform) => void
  /** True if can zoom in (not at max) */
  canZoomIn: boolean
  /** True if can zoom out (not at min) */
  canZoomOut: boolean
  /** Current min scale */
  minScale: number
  /** Current max scale */
  maxScale: number
  /** Current zoom increment */
  zoomIncrement: number
}

/**
 * Default configuration for zoom behavior
 */
export const ZOOM_CONFIG = {
  /** Default scale extent [min, max] */
  scaleExtent: [0.1, 4] as [number, number],
  /** Default zoom increment per button click */
  zoomIncrement: 0.2,
  /** Transition duration for programmatic zoom (ms) */
  transitionDuration: 200,
  /** Transition duration for fit to screen (ms) */
  fitTransitionDuration: 500,
  /** Padding for fit to screen calculation */
  fitPadding: 40,
}

/**
 * useZoomPan hook
 *
 * Manages zoom and pan state for D3 visualizations with React state sync.
 */
export function useZoomPan({
  svgRef,
  containerRef,
  minScale = ZOOM_CONFIG.scaleExtent[0],
  maxScale = ZOOM_CONFIG.scaleExtent[1],
  zoomIncrement = ZOOM_CONFIG.zoomIncrement,
}: UseZoomPanProps): UseZoomPanResult {
  const [zoomState, setZoomState] = useState<ZoomPanState>({
    scale: 1,
    translateX: 0,
    translateY: 0,
  })

  const zoomBehaviorRef = useRef<d3.ZoomBehavior<SVGSVGElement, unknown> | null>(null)

  // Calculate derived values
  const zoomPercent = Math.round(zoomState.scale * 100)
  const canZoomIn = zoomState.scale < maxScale
  const canZoomOut = zoomState.scale > minScale

  /**
   * Store zoom behavior for external control
   */
  const setZoomBehavior = useCallback((zoom: d3.ZoomBehavior<SVGSVGElement, unknown>) => {
    zoomBehaviorRef.current = zoom
  }, [])

  /**
   * Handle zoom change from D3 events
   */
  const handleZoomChange = useCallback((transform: ZoomTransform) => {
    setZoomState({
      scale: transform.k,
      translateX: transform.x,
      translateY: transform.y,
    })
  }, [])

  /**
   * Zoom in by increment
   */
  const zoomIn = useCallback(() => {
    if (!svgRef.current || !zoomBehaviorRef.current) return

    const svg = d3.select(svgRef.current)
    const newScale = Math.min(zoomState.scale + zoomIncrement, maxScale)

    svg
      .transition()
      .duration(ZOOM_CONFIG.transitionDuration)
      .call(zoomBehaviorRef.current.scaleTo, newScale)
  }, [svgRef, zoomState.scale, zoomIncrement, maxScale])

  /**
   * Zoom out by increment
   */
  const zoomOut = useCallback(() => {
    if (!svgRef.current || !zoomBehaviorRef.current) return

    const svg = d3.select(svgRef.current)
    const newScale = Math.max(zoomState.scale - zoomIncrement, minScale)

    svg
      .transition()
      .duration(ZOOM_CONFIG.transitionDuration)
      .call(zoomBehaviorRef.current.scaleTo, newScale)
  }, [svgRef, zoomState.scale, zoomIncrement, minScale])

  /**
   * Reset zoom to identity
   */
  const resetZoom = useCallback(() => {
    if (!svgRef.current || !zoomBehaviorRef.current) return

    const svg = d3.select(svgRef.current)

    svg
      .transition()
      .duration(ZOOM_CONFIG.transitionDuration)
      .call(zoomBehaviorRef.current.transform, d3.zoomIdentity)
  }, [svgRef])

  /**
   * Fit entire graph to viewport
   */
  const fitToScreen = useCallback(() => {
    if (!svgRef.current || !containerRef.current || !zoomBehaviorRef.current) return

    const svg = d3.select(svgRef.current)
    const svgNode = svgRef.current
    const containerNode = containerRef.current

    // Get SVG dimensions
    const { width: svgWidth, height: svgHeight } = svgNode.getBoundingClientRect()

    // Get container bounds (all nodes)
    const bounds = containerNode.getBBox()

    if (bounds.width === 0 || bounds.height === 0) return

    // Calculate scale to fit with padding
    const padding = ZOOM_CONFIG.fitPadding
    const scaleX = (svgWidth - padding * 2) / bounds.width
    const scaleY = (svgHeight - padding * 2) / bounds.height
    const scale = Math.min(scaleX, scaleY, maxScale)

    // Calculate translation to center
    const translateX = (svgWidth - bounds.width * scale) / 2 - bounds.x * scale
    const translateY = (svgHeight - bounds.height * scale) / 2 - bounds.y * scale

    const transform = d3.zoomIdentity.translate(translateX, translateY).scale(scale)

    svg
      .transition()
      .duration(ZOOM_CONFIG.fitTransitionDuration)
      .call(zoomBehaviorRef.current.transform, transform)
  }, [svgRef, containerRef, maxScale])

  return {
    zoomState,
    zoomPercent,
    zoomIn,
    zoomOut,
    resetZoom,
    fitToScreen,
    setZoomBehavior,
    handleZoomChange,
    canZoomIn,
    canZoomOut,
    minScale,
    maxScale,
    zoomIncrement,
  }
}
