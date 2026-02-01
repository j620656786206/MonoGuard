/**
 * useCanvasInteraction - Canvas mouse event handling for hit detection (Story 4.9)
 *
 * Provides hover and click interaction for Canvas-rendered graphs.
 * Uses coordinate transformation to map screen positions to graph space.
 */

import type { RefObject } from 'react'
import { useCallback } from 'react'
import type { D3Node, ViewportState } from './types'

interface UseCanvasInteractionOptions {
  canvasRef: RefObject<HTMLCanvasElement | null>
  nodesRef: RefObject<D3Node[]>
  viewport: ViewportState
  onNodeHover: (node: D3Node | null, position: { x: number; y: number } | null) => void
  onNodeSelect: (nodeId: string | null) => void
}

export function useCanvasInteraction({
  canvasRef,
  nodesRef,
  viewport,
  onNodeHover,
  onNodeSelect,
}: UseCanvasInteractionOptions) {
  const getMousePosition = useCallback(
    (e: React.MouseEvent<HTMLCanvasElement>) => {
      if (!canvasRef.current) return null

      const rect = canvasRef.current.getBoundingClientRect()
      const x = e.clientX - rect.left
      const y = e.clientY - rect.top

      // Transform to graph coordinates (reverse viewport transform)
      const graphX = (x - viewport.panX) / viewport.zoom
      const graphY = (y - viewport.panY) / viewport.zoom

      return { screenX: x, screenY: y, graphX, graphY }
    },
    [canvasRef, viewport]
  )

  const findNodeAtPosition = useCallback(
    (graphX: number, graphY: number): D3Node | null => {
      const nodes = nodesRef.current
      if (!nodes) return null
      const hitRadius = 12 // Slightly larger than node radius for easier interaction

      // Reverse iterate so top-rendered nodes are hit first
      for (let i = nodes.length - 1; i >= 0; i--) {
        const node = nodes[i]
        if (node.x === undefined || node.y === undefined) continue

        const dx = graphX - node.x
        const dy = graphY - node.y
        const distance = Math.sqrt(dx * dx + dy * dy)

        if (distance <= hitRadius) {
          return node
        }
      }

      return null
    },
    [nodesRef]
  )

  const handleMouseMove = useCallback(
    (e: React.MouseEvent<HTMLCanvasElement>) => {
      const pos = getMousePosition(e)
      if (!pos) return

      const node = findNodeAtPosition(pos.graphX, pos.graphY)

      if (node) {
        if (canvasRef.current) {
          canvasRef.current.style.cursor = 'pointer'
        }
        onNodeHover(node, { x: pos.screenX, y: pos.screenY })
      } else {
        if (canvasRef.current) {
          canvasRef.current.style.cursor = 'crosshair'
        }
        onNodeHover(null, null)
      }
    },
    [canvasRef, getMousePosition, findNodeAtPosition, onNodeHover]
  )

  const handleMouseClick = useCallback(
    (e: React.MouseEvent<HTMLCanvasElement>) => {
      const pos = getMousePosition(e)
      if (!pos) return

      const node = findNodeAtPosition(pos.graphX, pos.graphY)
      onNodeSelect(node?.id ?? null)
    },
    [getMousePosition, findNodeAtPosition, onNodeSelect]
  )

  return {
    handleMouseMove,
    handleMouseClick,
  }
}
