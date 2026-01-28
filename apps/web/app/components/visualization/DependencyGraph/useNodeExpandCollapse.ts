/**
 * useNodeExpandCollapse - Custom hook for managing node expand/collapse state
 *
 * Tracks which nodes are collapsed in the dependency graph and provides
 * functions to toggle, collapse, expand individual nodes or groups by depth.
 * Optionally persists state to session storage.
 *
 * @see Story 4.3: Implement Node Expand/Collapse Functionality
 */

import { useCallback, useEffect, useState } from 'react'

/**
 * State and methods returned by the useNodeExpandCollapse hook
 */
export interface ExpandCollapseState {
  /** Set of currently collapsed node IDs */
  collapsedNodeIds: Set<string>
  /** Toggle a node between collapsed and expanded */
  toggleNode: (nodeId: string) => void
  /** Collapse a specific node */
  collapseNode: (nodeId: string) => void
  /** Expand a specific node */
  expandNode: (nodeId: string) => void
  /** Collapse all nodes at or beyond specified depth */
  collapseAtDepth: (depth: number) => void
  /** Expand all nodes up to specified depth */
  expandToDepth: (depth: number) => void
  /** Check if a node is currently collapsed */
  isCollapsed: (nodeId: string) => boolean
  /** Expand all nodes */
  expandAll: () => void
  /** Collapse all nodes except root nodes (depth 0) */
  collapseAll: () => void
}

/**
 * Props for the useNodeExpandCollapse hook
 */
export interface UseNodeExpandCollapseProps {
  /** All node IDs in the graph */
  nodeIds: string[]
  /** Map of node ID to depth from root (0 = root node) */
  nodeDepths: Map<string, number>
  /** Optional key for session storage persistence */
  sessionKey?: string
}

/**
 * Custom hook for managing expand/collapse state of graph nodes
 *
 * @param props - Configuration for the hook
 * @returns State and methods for managing node visibility
 */
export function useNodeExpandCollapse({
  nodeIds,
  nodeDepths,
  sessionKey,
}: UseNodeExpandCollapseProps): ExpandCollapseState {
  // Initialize from session storage if available
  const [collapsedNodeIds, setCollapsedNodeIds] = useState<Set<string>>(() => {
    if (sessionKey && typeof sessionStorage !== 'undefined') {
      try {
        const stored = sessionStorage.getItem(`monoguard-collapse-${sessionKey}`)
        if (stored) {
          return new Set(JSON.parse(stored) as string[])
        }
      } catch {
        // Ignore parse errors, start with empty set
      }
    }
    return new Set()
  })

  // Persist to session storage when state changes
  useEffect(() => {
    if (sessionKey && typeof sessionStorage !== 'undefined') {
      sessionStorage.setItem(
        `monoguard-collapse-${sessionKey}`,
        JSON.stringify([...collapsedNodeIds])
      )
    }
  }, [collapsedNodeIds, sessionKey])

  const toggleNode = useCallback((nodeId: string) => {
    setCollapsedNodeIds((prev) => {
      const next = new Set(prev)
      if (next.has(nodeId)) {
        next.delete(nodeId)
      } else {
        next.add(nodeId)
      }
      return next
    })
  }, [])

  const collapseNode = useCallback((nodeId: string) => {
    setCollapsedNodeIds((prev) => new Set([...prev, nodeId]))
  }, [])

  const expandNode = useCallback((nodeId: string) => {
    setCollapsedNodeIds((prev) => {
      const next = new Set(prev)
      next.delete(nodeId)
      return next
    })
  }, [])

  const collapseAtDepth = useCallback(
    (depth: number) => {
      const toCollapse = nodeIds.filter((id) => (nodeDepths.get(id) ?? 0) >= depth)
      setCollapsedNodeIds(new Set(toCollapse))
    },
    [nodeIds, nodeDepths]
  )

  const expandToDepth = useCallback(
    (depth: number) => {
      const toKeepCollapsed = nodeIds.filter((id) => (nodeDepths.get(id) ?? 0) > depth)
      setCollapsedNodeIds(new Set(toKeepCollapsed))
    },
    [nodeIds, nodeDepths]
  )

  const expandAll = useCallback(() => {
    setCollapsedNodeIds(new Set())
  }, [])

  const collapseAll = useCallback(() => {
    // Collapse all nodes except root nodes (depth 0)
    const toCollapse = nodeIds.filter((id) => (nodeDepths.get(id) ?? 0) > 0)
    setCollapsedNodeIds(new Set(toCollapse))
  }, [nodeIds, nodeDepths])

  const isCollapsed = useCallback(
    (nodeId: string) => collapsedNodeIds.has(nodeId),
    [collapsedNodeIds]
  )

  return {
    collapsedNodeIds,
    toggleNode,
    collapseNode,
    expandNode,
    collapseAtDepth,
    expandToDepth,
    isCollapsed,
    expandAll,
    collapseAll,
  }
}
