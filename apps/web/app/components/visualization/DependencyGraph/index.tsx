/**
 * DependencyGraphViz - Force-directed dependency graph visualization
 *
 * Renders package dependencies as an interactive force-directed graph using D3.js.
 * Uses SVG rendering for graphs with < 500 nodes (per architecture.md).
 * Highlights circular dependencies with distinct styling (Story 4.2).
 * Supports expand/collapse of nodes with depth-based controls (Story 4.3).
 *
 * @see Story 4.1: Implement D3.js Force-Directed Dependency Graph
 * @see Story 4.2: Highlight Circular Dependencies in Graph
 * @see Story 4.3: Implement Node Expand/Collapse Functionality
 */
'use client'

import * as d3 from 'd3'
import React, { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { GraphControls } from './GraphControls'
import { GraphLegend } from './GraphLegend'
import {
  ANIMATION,
  COLLAPSED_STYLES,
  EDGE_COLORS,
  EXPAND_COLLAPSE_ANIMATION,
  GLOW_FILTER,
  NODE_COLORS,
} from './styles'
import type { D3Link, D3Node, DependencyGraphProps } from './types'
import { DEFAULT_SIMULATION_CONFIG } from './types'
import { transformToD3Data, truncatePackageName } from './useForceSimulation'
import { useNodeExpandCollapse } from './useNodeExpandCollapse'
import { calculateNodeDepths } from './utils/calculateDepth'
import { computeVisibleNodes } from './utils/computeVisibleNodes'

/**
 * DependencyGraphViz component
 *
 * Renders a force-directed graph visualization of package dependencies.
 * Uses React.memo to prevent unnecessary re-renders (per project-context.md).
 */
export const DependencyGraphViz = React.memo(function DependencyGraphViz({
  data,
  circularDependencies,
  className = '',
  width = '100%',
  height = 500,
}: DependencyGraphProps) {
  const svgRef = useRef<SVGSVGElement>(null)
  const containerRef = useRef<HTMLDivElement>(null)
  const simulationRef = useRef<d3.Simulation<D3Node, D3Link> | null>(null)
  const [dimensions, setDimensions] = useState({ width: 800, height: 500 })
  const [selectedCycleIndex, setSelectedCycleIndex] = useState<number | null>(null)
  const [currentDepth, setCurrentDepth] = useState<number | 'all'>('all')

  // Refs to store D3 selections for style updates without full redraw
  const nodeSelectionRef = useRef<d3.Selection<SVGGElement, D3Node, SVGGElement, unknown> | null>(
    null
  )
  const normalLinkSelectionRef = useRef<d3.Selection<
    SVGLineElement,
    D3Link,
    SVGGElement,
    unknown
  > | null>(null)
  const cycleLinkSelectionRef = useRef<d3.Selection<
    SVGLineElement,
    D3Link,
    SVGGElement,
    unknown
  > | null>(null)
  const badgeSelectionRef = useRef<d3.Selection<SVGGElement, D3Node, SVGGElement, unknown> | null>(
    null
  )

  // Transform data to D3 format with cycle information
  const fullGraphData = useMemo(() => {
    if (!data) return { nodes: [], links: [] }
    return transformToD3Data(data, { circularDependencies })
  }, [data, circularDependencies])

  // Calculate node depths for depth-based controls
  const nodeDepths = useMemo(() => {
    const edges = fullGraphData.links.map((l) => ({
      source: typeof l.source === 'string' ? l.source : l.source.id,
      target: typeof l.target === 'string' ? l.target : l.target.id,
    }))
    return calculateNodeDepths(
      fullGraphData.nodes.map((n) => n.id),
      edges
    )
  }, [fullGraphData])

  // Calculate max depth for controls
  const maxDepth = useMemo(() => {
    if (nodeDepths.size === 0) return 0
    return Math.max(...nodeDepths.values())
  }, [nodeDepths])

  // Initialize expand/collapse state
  const { collapsedNodeIds, toggleNode, expandAll, collapseAll, collapseAtDepth } =
    useNodeExpandCollapse({
      nodeIds: fullGraphData.nodes.map((n) => n.id),
      nodeDepths,
      sessionKey: data ? `graph-${Object.keys(data.nodes).length}` : undefined,
    })

  // Compute visible nodes based on collapsed state
  const { visibleNodes, visibleLinks, hiddenChildCounts } = useMemo(() => {
    return computeVisibleNodes(fullGraphData.nodes, fullGraphData.links, collapsedNodeIds)
  }, [fullGraphData, collapsedNodeIds])

  // Handle depth change from controls
  const handleDepthChange = useCallback(
    (depth: number | 'all') => {
      setCurrentDepth(depth)
      if (depth === 'all') {
        expandAll()
      } else {
        collapseAtDepth(depth + 1) // Collapse nodes at depth > selected
      }
    },
    [expandAll, collapseAtDepth]
  )

  // Handle cycle selection clear on Escape key
  const handleKeyDown = useCallback(
    (event: KeyboardEvent) => {
      if (event.key === 'Escape' && selectedCycleIndex !== null) {
        setSelectedCycleIndex(null)
      }
    },
    [selectedCycleIndex]
  )

  useEffect(() => {
    document.addEventListener('keydown', handleKeyDown)
    return () => {
      document.removeEventListener('keydown', handleKeyDown)
    }
  }, [handleKeyDown])

  // Handle resize with ResizeObserver
  useEffect(() => {
    if (!containerRef.current) return

    const resizeObserver = new ResizeObserver((entries) => {
      for (const entry of entries) {
        const { width: containerWidth, height: containerHeight } = entry.contentRect
        setDimensions({
          width: containerWidth || 800,
          height: containerHeight || (typeof height === 'number' ? height : 500),
        })
      }
    })

    resizeObserver.observe(containerRef.current)

    return () => {
      resizeObserver.disconnect()
    }
  }, [height])

  // Main D3 initialization effect - only runs when data/dimensions change
  useEffect(() => {
    if (!svgRef.current || !data) return

    const svg = d3.select(svgRef.current)
    const { width: svgWidth, height: svgHeight } = dimensions

    // Clear previous content
    svg.selectAll('*').remove()

    // If no nodes, don't render anything
    if (visibleNodes.length === 0) {
      nodeSelectionRef.current = null
      normalLinkSelectionRef.current = null
      cycleLinkSelectionRef.current = null
      badgeSelectionRef.current = null
      return
    }

    // Create main group for zoom/pan
    const g = svg.append('g')

    // Define SVG filters and markers
    const defs = svg.append('defs')

    // Glow filter for cycle nodes (AC1: red glow effect)
    const glowFilter = defs
      .append('filter')
      .attr('id', 'glow')
      .attr('x', '-50%')
      .attr('y', '-50%')
      .attr('width', '200%')
      .attr('height', '200%')

    glowFilter
      .append('feGaussianBlur')
      .attr('stdDeviation', GLOW_FILTER.stdDeviation)
      .attr('result', 'coloredBlur')

    const feMerge = glowFilter.append('feMerge')
    feMerge.append('feMergeNode').attr('in', 'coloredBlur')
    feMerge.append('feMergeNode').attr('in', 'SourceGraphic')

    // Arrow marker for normal edges
    defs
      .append('marker')
      .attr('id', 'arrowhead')
      .attr('viewBox', '0 -5 10 10')
      .attr('refX', 20)
      .attr('refY', 0)
      .attr('markerWidth', 6)
      .attr('markerHeight', 6)
      .attr('orient', 'auto')
      .append('path')
      .attr('d', 'M0,-5L10,0L0,5')
      .attr('fill', EDGE_COLORS.normal.stroke)

    // Arrow marker for cycle edges (AC2: red arrows)
    defs
      .append('marker')
      .attr('id', 'arrowhead-cycle')
      .attr('viewBox', '0 -5 10 10')
      .attr('refX', 20)
      .attr('refY', 0)
      .attr('markerWidth', 6)
      .attr('markerHeight', 6)
      .attr('orient', 'auto')
      .append('path')
      .attr('d', 'M0,-5L10,0L0,5')
      .attr('fill', EDGE_COLORS.cycle.stroke)

    // Separate links into normal and cycle links for layering
    const normalLinks = visibleLinks.filter((l) => !l.inCycle)
    const cycleLinks = visibleLinks.filter((l) => l.inCycle)

    // Create link elements - render normal links first (below cycle links)
    const normalLink = g
      .append('g')
      .attr('class', 'links-normal')
      .selectAll<SVGLineElement, D3Link>('line')
      .data(normalLinks)
      .join('line')
      .attr('stroke', EDGE_COLORS.normal.stroke)
      .attr('stroke-opacity', EDGE_COLORS.normal.opacity)
      .attr('stroke-width', EDGE_COLORS.normal.width)
      .attr('marker-end', 'url(#arrowhead)')

    // Create cycle link elements (above normal links)
    const cycleLink = g
      .append('g')
      .attr('class', 'links-cycle')
      .selectAll<SVGLineElement, D3Link>('line')
      .data(cycleLinks)
      .join('line')
      .attr('stroke', EDGE_COLORS.cycle.stroke)
      .attr('stroke-opacity', EDGE_COLORS.cycle.opacity)
      .attr('stroke-width', EDGE_COLORS.cycle.width)
      .attr('marker-end', 'url(#arrowhead-cycle)')
      // AC3: Animated cycle paths with flowing effect
      .attr('stroke-dasharray', ANIMATION.dashArray)
      .style('animation', `flowAnimation ${ANIMATION.flowDuration} linear infinite`)

    // Create node group elements
    const node = g
      .append('g')
      .attr('class', 'nodes')
      .selectAll<SVGGElement, D3Node>('g')
      .data(visibleNodes)
      .join('g')
      .attr('class', (d) => {
        let classes = 'node'
        if (d.inCycle) classes += ' node--cycle'
        if (collapsedNodeIds.has(d.id)) classes += ' node--collapsed'
        return classes
      })

    // Add circles to nodes with cycle-aware and collapse-aware styling
    node
      .append('circle')
      .attr('r', (d) => Math.max(8, Math.min(16, 8 + d.dependencyCount * 0.5)))
      .attr('fill', (d) => {
        if (collapsedNodeIds.has(d.id)) return COLLAPSED_STYLES.node.fill
        if (d.inCycle) return NODE_COLORS.cycle.fill
        return NODE_COLORS.normal.fill
      })
      .attr('stroke', (d) => {
        if (collapsedNodeIds.has(d.id)) return COLLAPSED_STYLES.node.stroke
        if (d.inCycle) return NODE_COLORS.cycle.stroke
        return NODE_COLORS.normal.stroke
      })
      .attr('stroke-width', (d) => {
        if (collapsedNodeIds.has(d.id)) return COLLAPSED_STYLES.node.strokeWidth
        return d.inCycle ? 3 : 2
      })
      .attr('stroke-dasharray', (d) =>
        collapsedNodeIds.has(d.id) ? COLLAPSED_STYLES.node.strokeDasharray : null
      )
      .attr('filter', (d) => (d.inCycle && !collapsedNodeIds.has(d.id) ? 'url(#glow)' : null))
      .style('cursor', 'pointer')

    // Add labels to nodes
    node
      .append('text')
      .text((d) => truncatePackageName(d.name))
      .attr('font-size', '11px')
      .attr('font-family', 'system-ui, sans-serif')
      .attr('fill', '#374151')
      .attr('dx', 14)
      .attr('dy', 4)
      .style('pointer-events', 'none')
      .style('user-select', 'none')

    // Story 4.3 AC4: Add collapsed badge with hidden child count
    const collapsedNodesWithCount = visibleNodes.filter(
      (n) => collapsedNodeIds.has(n.id) && (hiddenChildCounts.get(n.id) ?? 0) > 0
    )

    const badge = g
      .append('g')
      .attr('class', 'collapsed-badges')
      .selectAll<SVGGElement, D3Node>('g')
      .data(collapsedNodesWithCount)
      .join('g')
      .attr('class', 'collapsed-badge')

    badge
      .append('circle')
      .attr('r', COLLAPSED_STYLES.badge.radius)
      .attr('fill', COLLAPSED_STYLES.badge.fill)

    badge
      .append('text')
      .attr('text-anchor', 'middle')
      .attr('dominant-baseline', 'central')
      .attr('fill', COLLAPSED_STYLES.badge.textFill)
      .attr('font-size', COLLAPSED_STYLES.badge.fontSize)
      .attr('font-weight', COLLAPSED_STYLES.badge.fontWeight)
      .text((d) => {
        const count = hiddenChildCounts.get(d.id) ?? 0
        return count > 99 ? '99+' : String(count)
      })

    // Store selections in refs for updates
    nodeSelectionRef.current = node
    normalLinkSelectionRef.current = normalLink
    cycleLinkSelectionRef.current = cycleLink
    badgeSelectionRef.current = badge

    // Story 4.3 AC1/AC2: Double-click handler for expand/collapse
    // Use timer pattern to differentiate from single click
    let clickTimer: ReturnType<typeof setTimeout> | null = null

    node.on('click', (_event, d) => {
      if (clickTimer) {
        // Double-click detected
        clearTimeout(clickTimer)
        clickTimer = null
        return
      }

      clickTimer = setTimeout(() => {
        // Single click - handle cycle selection
        if (d.inCycle && d.cycleIds.length > 0) {
          setSelectedCycleIndex((prev) => {
            if (prev !== null && d.cycleIds.includes(prev)) {
              return null
            }
            return d.cycleIds[0]
          })
        } else {
          setSelectedCycleIndex(null)
        }
        clickTimer = null
      }, 200)
    })

    node.on('dblclick', (event, d) => {
      event.stopPropagation()
      if (clickTimer) {
        clearTimeout(clickTimer)
        clickTimer = null
      }
      toggleNode(d.id)
    })

    // AC6: Click on background to deselect
    svg.on('click', (event) => {
      if (event.target === svgRef.current) {
        setSelectedCycleIndex(null)
      }
    })

    // Create force simulation
    const simulation = d3
      .forceSimulation<D3Node>(visibleNodes)
      .force(
        'link',
        d3
          .forceLink<D3Node, D3Link>(visibleLinks)
          .id((d) => d.id)
          .distance(DEFAULT_SIMULATION_CONFIG.linkDistance)
      )
      .force(
        'charge',
        d3.forceManyBody<D3Node>().strength(DEFAULT_SIMULATION_CONFIG.chargeStrength)
      )
      .force('center', d3.forceCenter(svgWidth / 2, svgHeight / 2))
      .force(
        'collision',
        d3.forceCollide<D3Node>().radius(DEFAULT_SIMULATION_CONFIG.collisionRadius)
      )
      .alphaDecay(EXPAND_COLLAPSE_ANIMATION.alphaDecay)

    simulationRef.current = simulation

    // Update positions on each tick
    simulation.on('tick', () => {
      // Update normal links
      normalLink
        .attr('x1', (d) => (d.source as D3Node).x ?? 0)
        .attr('y1', (d) => (d.source as D3Node).y ?? 0)
        .attr('x2', (d) => (d.target as D3Node).x ?? 0)
        .attr('y2', (d) => (d.target as D3Node).y ?? 0)

      // Update cycle links
      cycleLink
        .attr('x1', (d) => (d.source as D3Node).x ?? 0)
        .attr('y1', (d) => (d.source as D3Node).y ?? 0)
        .attr('x2', (d) => (d.target as D3Node).x ?? 0)
        .attr('y2', (d) => (d.target as D3Node).y ?? 0)

      node.attr('transform', (d) => `translate(${d.x ?? 0},${d.y ?? 0})`)

      // Update badge positions (Story 4.3)
      badge.attr('transform', (d) => {
        // Find the corresponding node to get position
        const nodeData = visibleNodes.find((n) => n.id === d.id)
        if (nodeData) {
          return `translate(${(nodeData.x ?? 0) + COLLAPSED_STYLES.badge.offsetX}, ${(nodeData.y ?? 0) + COLLAPSED_STYLES.badge.offsetY})`
        }
        return ''
      })
    })

    // Add drag behavior
    const drag = d3
      .drag<SVGGElement, D3Node>()
      .on('start', (event, d) => {
        if (!event.active) simulation.alphaTarget(0.3).restart()
        d.fx = d.x
        d.fy = d.y
      })
      .on('drag', (event, d) => {
        d.fx = event.x
        d.fy = event.y
      })
      .on('end', (event, d) => {
        if (!event.active) simulation.alphaTarget(0)
        d.fx = null
        d.fy = null
      })

    node.call(drag)

    // Add basic zoom behavior (setup for Story 4.4)
    const zoom = d3
      .zoom<SVGSVGElement, unknown>()
      .scaleExtent([0.1, 4])
      .on('zoom', (event) => {
        g.attr('transform', event.transform)
      })

    svg.call(zoom)

    // CRITICAL: Cleanup to prevent memory leaks (per project-context.md)
    return () => {
      if (clickTimer) clearTimeout(clickTimer)
      simulation.stop()
      svg.on('.zoom', null) // Remove zoom listener
      svg.on('click', null) // Remove click listener
      node.on('click', null) // Remove node click listeners
      node.on('dblclick', null) // Remove double-click listeners
      svg.selectAll('*').remove() // Clean DOM
      nodeSelectionRef.current = null
      normalLinkSelectionRef.current = null
      cycleLinkSelectionRef.current = null
      badgeSelectionRef.current = null
      simulationRef.current = null
    }
  }, [
    data,
    visibleNodes,
    visibleLinks,
    collapsedNodeIds,
    hiddenChildCounts,
    dimensions,
    toggleNode,
  ])

  // Separate effect for style updates when selection changes (performance optimization)
  // This avoids recreating the entire graph when just the selection changes
  useEffect(() => {
    const nodeSelection = nodeSelectionRef.current
    const normalLinkSelection = normalLinkSelectionRef.current
    const cycleLinkSelection = cycleLinkSelectionRef.current

    if (!nodeSelection || !normalLinkSelection || !cycleLinkSelection) return

    // Helper to determine node styling based on cycle state and selection
    const getNodeStyle = (d: D3Node) => {
      if (selectedCycleIndex !== null) {
        const isInSelected = d.cycleIds.includes(selectedCycleIndex)
        if (isInSelected) {
          return NODE_COLORS.selected
        }
        return NODE_COLORS.dimmed
      }
      if (collapsedNodeIds.has(d.id)) {
        return COLLAPSED_STYLES.node
      }
      if (d.inCycle) {
        return NODE_COLORS.cycle
      }
      return NODE_COLORS.normal
    }

    // Helper to determine edge styling based on cycle state and selection
    const getEdgeStyle = (d: D3Link) => {
      if (selectedCycleIndex !== null) {
        const isInSelected = d.cycleIds.includes(selectedCycleIndex)
        if (isInSelected) {
          return EDGE_COLORS.selected
        }
        return EDGE_COLORS.dimmed
      }
      if (d.inCycle) {
        return EDGE_COLORS.cycle
      }
      return EDGE_COLORS.normal
    }

    // Update node circle styles
    nodeSelection
      .select('circle')
      .attr('fill', (d) => getNodeStyle(d).fill)
      .attr('stroke', (d) => getNodeStyle(d).stroke)
      .attr('filter', (d) =>
        d.inCycle && selectedCycleIndex === null && !collapsedNodeIds.has(d.id)
          ? 'url(#glow)'
          : null
      )

    // Update node text styles
    nodeSelection
      .select('text')
      .attr('fill', (d) =>
        selectedCycleIndex !== null && !d.cycleIds.includes(selectedCycleIndex)
          ? '#9ca3af'
          : '#374151'
      )

    // Update normal link styles
    normalLinkSelection
      .attr('stroke', (d) => getEdgeStyle(d).stroke)
      .attr('stroke-opacity', (d) => getEdgeStyle(d).opacity)
      .attr('stroke-width', (d) => getEdgeStyle(d).width)

    // Update cycle link styles
    cycleLinkSelection
      .attr('stroke', (d) => getEdgeStyle(d).stroke)
      .attr('stroke-opacity', (d) => getEdgeStyle(d).opacity)
      .attr('stroke-width', (d) => getEdgeStyle(d).width)
  }, [selectedCycleIndex, collapsedNodeIds])

  // Determine if there are any cycles to display in legend
  const hasCycles = circularDependencies && circularDependencies.length > 0

  return (
    <div
      ref={containerRef}
      className={`relative ${className}`}
      style={{
        width: typeof width === 'number' ? `${width}px` : width,
        height: typeof height === 'number' ? `${height}px` : height,
        minHeight: '500px',
      }}
    >
      {/* CSS animation for cycle edge flow effect (AC3) */}
      <style>
        {`
          @keyframes flowAnimation {
            0% { stroke-dashoffset: ${ANIMATION.dashOffset}; }
            100% { stroke-dashoffset: 0; }
          }
        `}
      </style>
      <svg
        ref={svgRef}
        className="h-full w-full"
        style={{
          width: '100%',
          height: '100%',
        }}
      />
      {/* Story 4.3 AC3: Depth-based controls */}
      {fullGraphData.nodes.length > 0 && (
        <GraphControls
          currentDepth={currentDepth}
          maxDepth={maxDepth}
          onDepthChange={handleDepthChange}
          onExpandAll={expandAll}
          onCollapseAll={collapseAll}
        />
      )}
      {/* AC4: Color legend showing meaning of different elements */}
      <GraphLegend position="bottom-left" hasCycles={hasCycles} />
    </div>
  )
})

export type { GraphControlsProps } from './GraphControls'
export { GraphControls } from './GraphControls'
export type { GraphLegendProps } from './GraphLegend'
export { GraphLegend } from './GraphLegend'
// Re-export types and utilities
export type { D3GraphData, D3Link, D3Node, DependencyGraphProps } from './types'
export type { CycleHighlightResult } from './useCycleHighlight'
export { useCycleHighlight } from './useCycleHighlight'
export { transformToD3Data, truncatePackageName } from './useForceSimulation'
export type { ExpandCollapseState, UseNodeExpandCollapseProps } from './useNodeExpandCollapse'
export { useNodeExpandCollapse } from './useNodeExpandCollapse'
export { calculateNodeDepths } from './utils/calculateDepth'
export { computeVisibleNodes } from './utils/computeVisibleNodes'
