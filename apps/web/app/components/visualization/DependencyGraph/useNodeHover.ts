/**
 * useNodeHover - Custom hook for managing node hover state
 *
 * Handles mouse enter/leave/move events on graph nodes and computes
 * which nodes and edges are connected to the hovered node for highlighting.
 *
 * @see Story 4.5: Implement Hover Details and Tooltips (AC2, AC4, AC5)
 */

import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import type { D3Link, D3Node, HoverState } from './types'

/**
 * Props for useNodeHover hook
 */
export interface UseNodeHoverProps {
  /** All nodes in the graph (reserved for future use, e.g., node lookup optimization) */
  nodes?: D3Node[]
  /** All links in the graph */
  links: D3Link[]
}

/**
 * Result from useNodeHover hook
 */
export interface UseNodeHoverResult {
  /** Current hover state (node ID and position) */
  hoverState: HoverState
  /** Set of node IDs connected to the hovered node */
  connectedNodeIds: Set<string>
  /** Handler for mouse enter on a node */
  handleNodeMouseEnter: (nodeId: string, event: MouseEvent) => void
  /** Handler for mouse leave from a node */
  handleNodeMouseLeave: () => void
  /** Handler for mouse move within a node */
  handleNodeMouseMove: (event: MouseEvent) => void
}

/**
 * Helper to extract node ID from link source/target
 * D3 replaces string IDs with node objects during simulation
 */
function getNodeId(nodeOrId: string | D3Node): string {
  return typeof nodeOrId === 'string' ? nodeOrId : nodeOrId.id
}

/**
 * useNodeHover hook
 *
 * Manages hover state and computes connected elements for highlighting.
 * Uses memoization for performance with large graphs.
 */
export function useNodeHover({ links }: UseNodeHoverProps): UseNodeHoverResult {
  const [hoverState, setHoverState] = useState<HoverState>({
    nodeId: null,
    position: null,
  })

  // CR2-3: rAF ref for throttling mousemove position updates
  const rafRef = useRef<number | null>(null)

  // Cleanup rAF on unmount
  useEffect(() => {
    return () => {
      if (rafRef.current !== null) {
        cancelAnimationFrame(rafRef.current)
      }
    }
  }, [])

  // Compute connected nodes when hover changes (AC4)
  // CR2-1: Removed connectedLinkIndices - edge highlighting uses isLinkConnected() in index.tsx
  const connectedNodeIds = useMemo(() => {
    if (!hoverState.nodeId) {
      return new Set<string>()
    }

    const nodeIds = new Set<string>([hoverState.nodeId])

    links.forEach((link) => {
      const sourceId = getNodeId(link.source)
      const targetId = getNodeId(link.target)

      if (sourceId === hoverState.nodeId || targetId === hoverState.nodeId) {
        nodeIds.add(sourceId)
        nodeIds.add(targetId)
      }
    })

    return nodeIds
  }, [hoverState.nodeId, links])

  // Handler for mouse entering a node (AC2)
  const handleNodeMouseEnter = useCallback((nodeId: string, event: MouseEvent) => {
    setHoverState({
      nodeId,
      position: { x: event.clientX, y: event.clientY },
    })
  }, [])

  // Handler for mouse moving within a node (AC2, AC5)
  // CR2-3: Throttled with rAF to prevent excessive re-renders on rapid mouse movement
  const handleNodeMouseMove = useCallback((event: MouseEvent) => {
    if (rafRef.current !== null) return
    rafRef.current = requestAnimationFrame(() => {
      setHoverState((prev) => ({
        ...prev,
        position: { x: event.clientX, y: event.clientY },
      }))
      rafRef.current = null
    })
  }, [])

  // Handler for mouse leaving a node (AC2)
  const handleNodeMouseLeave = useCallback(() => {
    setHoverState({ nodeId: null, position: null })
  }, [])

  return {
    hoverState,
    connectedNodeIds,
    handleNodeMouseEnter,
    handleNodeMouseLeave,
    handleNodeMouseMove,
  }
}
