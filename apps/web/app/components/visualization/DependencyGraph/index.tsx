/**
 * DependencyGraphViz - Force-directed dependency graph visualization
 *
 * Renders package dependencies as an interactive force-directed graph using D3.js.
 * Uses SVG rendering for graphs with < 500 nodes (per architecture.md).
 * Highlights circular dependencies with distinct styling (Story 4.2).
 *
 * @see Story 4.1: Implement D3.js Force-Directed Dependency Graph
 * @see Story 4.2: Highlight Circular Dependencies in Graph
 */
'use client'

import * as d3 from 'd3'
import React, { useCallback, useEffect, useRef, useState } from 'react'
import { GraphLegend } from './GraphLegend'
import { ANIMATION, EDGE_COLORS, GLOW_FILTER, NODE_COLORS } from './styles'
import type { D3Link, D3Node, DependencyGraphProps } from './types'
import { DEFAULT_SIMULATION_CONFIG } from './types'
import { transformToD3Data, truncatePackageName } from './useForceSimulation'

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
  const [dimensions, setDimensions] = useState({ width: 800, height: 500 })
  const [selectedCycleIndex, setSelectedCycleIndex] = useState<number | null>(null)

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

    // Transform data to D3 format with cycle information
    const { nodes, links } = transformToD3Data(data, { circularDependencies })

    // If no nodes, don't render anything
    if (nodes.length === 0) {
      nodeSelectionRef.current = null
      normalLinkSelectionRef.current = null
      cycleLinkSelectionRef.current = null
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
    const normalLinks = links.filter((l) => !l.inCycle)
    const cycleLinks = links.filter((l) => l.inCycle)

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
      .data(nodes)
      .join('g')
      .attr('class', (d) => `node ${d.inCycle ? 'node--cycle' : ''}`)

    // Add circles to nodes with cycle-aware styling
    node
      .append('circle')
      .attr('r', (d) => Math.max(8, Math.min(16, 8 + d.dependencyCount * 0.5)))
      .attr('fill', (d) => (d.inCycle ? NODE_COLORS.cycle.fill : NODE_COLORS.normal.fill))
      .attr('stroke', (d) => (d.inCycle ? NODE_COLORS.cycle.stroke : NODE_COLORS.normal.stroke))
      .attr('stroke-width', (d) => (d.inCycle ? 3 : 2))
      .attr('filter', (d) => (d.inCycle ? 'url(#glow)' : null))
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

    // Store selections in refs for style updates
    nodeSelectionRef.current = node
    normalLinkSelectionRef.current = normalLink
    cycleLinkSelectionRef.current = cycleLink

    // AC5: Click handler for cycle nodes to highlight specific cycle
    node.on('click', (_event, d) => {
      if (d.inCycle && d.cycleIds.length > 0) {
        // Toggle selection: if clicking same cycle node, deselect
        setSelectedCycleIndex((prev) => {
          if (prev !== null && d.cycleIds.includes(prev)) {
            return null
          }
          return d.cycleIds[0]
        })
      } else {
        // Clicking a non-cycle node clears selection
        setSelectedCycleIndex(null)
      }
    })

    // AC6: Click on background to deselect
    svg.on('click', (event) => {
      if (event.target === svgRef.current) {
        setSelectedCycleIndex(null)
      }
    })

    // Create force simulation
    const simulation = d3
      .forceSimulation<D3Node>(nodes)
      .force(
        'link',
        d3
          .forceLink<D3Node, D3Link>(links)
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
      .alphaDecay(DEFAULT_SIMULATION_CONFIG.alphaDecay)

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
      simulation.stop()
      svg.on('.zoom', null) // Remove zoom listener
      svg.on('click', null) // Remove click listener
      node.on('click', null) // Remove node click listeners
      svg.selectAll('*').remove() // Clean DOM
      nodeSelectionRef.current = null
      normalLinkSelectionRef.current = null
      cycleLinkSelectionRef.current = null
    }
  }, [data, circularDependencies, dimensions])

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
      .attr('filter', (d) => (d.inCycle && selectedCycleIndex === null ? 'url(#glow)' : null))

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
  }, [selectedCycleIndex])

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
      {/* AC4: Color legend showing meaning of different elements */}
      <GraphLegend position="bottom-left" hasCycles={hasCycles} />
    </div>
  )
})

export type { GraphLegendProps } from './GraphLegend'
export { GraphLegend } from './GraphLegend'
// Re-export types and utilities
export type { D3GraphData, D3Link, D3Node, DependencyGraphProps } from './types'
export type { CycleHighlightResult } from './useCycleHighlight'
export { useCycleHighlight } from './useCycleHighlight'
export { transformToD3Data, truncatePackageName } from './useForceSimulation'
