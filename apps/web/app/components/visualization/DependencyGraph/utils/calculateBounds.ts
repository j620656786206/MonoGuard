/**
 * Utilities for calculating graph bounds
 *
 * Used for fit-to-screen, minimap viewport calculation, and graph centering.
 *
 * @see Story 4.4: Add Zoom, Pan, and Navigation Controls
 */

/**
 * Bounding box representation
 */
export interface Bounds {
  /** X coordinate of top-left corner */
  x: number
  /** Y coordinate of top-left corner */
  y: number
  /** Width of bounding box */
  width: number
  /** Height of bounding box */
  height: number
}

/**
 * Calculate bounding box of all nodes
 *
 * @param nodes - Array of nodes with x, y coordinates
 * @param padding - Padding to add around nodes (default: 20)
 * @returns Bounding box containing all nodes
 */
export function calculateNodeBounds(
  nodes: ReadonlyArray<{ x?: number; y?: number }>,
  padding: number = 20
): Bounds {
  if (nodes.length === 0) {
    return { x: 0, y: 0, width: 0, height: 0 }
  }

  let minX = Infinity
  let minY = Infinity
  let maxX = -Infinity
  let maxY = -Infinity

  nodes.forEach((node) => {
    const x = node.x ?? 0
    const y = node.y ?? 0
    minX = Math.min(minX, x)
    minY = Math.min(minY, y)
    maxX = Math.max(maxX, x)
    maxY = Math.max(maxY, y)
  })

  return {
    x: minX - padding,
    y: minY - padding,
    width: maxX - minX + padding * 2,
    height: maxY - minY + padding * 2,
  }
}

/**
 * Calculate viewport bounds in graph coordinates
 *
 * Converts the current transform (from D3 zoom) to viewport bounds
 * in the graph's coordinate system.
 *
 * @param transform - Current zoom transform (scale, x, y)
 * @param svgWidth - Width of SVG element
 * @param svgHeight - Height of SVG element
 * @returns Viewport bounds in graph coordinates
 */
export function calculateViewportBounds(
  transform: { k: number; x: number; y: number },
  svgWidth: number,
  svgHeight: number
): Bounds {
  // Inverse transform to get viewport in graph coordinates
  return {
    x: -transform.x / transform.k,
    y: -transform.y / transform.k,
    width: svgWidth / transform.k,
    height: svgHeight / transform.k,
  }
}

/**
 * Calculate transform to fit bounds within container
 *
 * @param bounds - Bounds to fit
 * @param containerWidth - Container width
 * @param containerHeight - Container height
 * @param padding - Padding inside container
 * @param maxScale - Maximum allowed scale
 * @returns Transform to fit bounds in container
 */
export function calculateFitTransform(
  bounds: Bounds,
  containerWidth: number,
  containerHeight: number,
  padding: number = 40,
  maxScale: number = 4
): { scale: number; translateX: number; translateY: number } {
  if (bounds.width === 0 || bounds.height === 0) {
    return { scale: 1, translateX: 0, translateY: 0 }
  }

  // Calculate scale to fit with padding
  const scaleX = (containerWidth - padding * 2) / bounds.width
  const scaleY = (containerHeight - padding * 2) / bounds.height
  const scale = Math.min(scaleX, scaleY, maxScale)

  // Calculate translation to center
  const translateX = (containerWidth - bounds.width * scale) / 2 - bounds.x * scale
  const translateY = (containerHeight - bounds.height * scale) / 2 - bounds.y * scale

  return { scale, translateX, translateY }
}
