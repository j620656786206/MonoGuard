/**
 * DependencyGraphViz - Force-directed dependency graph visualization
 *
 * Renders package dependencies as an interactive force-directed graph using D3.js.
 * Uses SVG rendering for graphs with < 500 nodes (per architecture.md).
 *
 * @see Story 4.1: Implement D3.js Force-Directed Dependency Graph
 */
'use client'

import * as d3 from 'd3'
import React, { useEffect, useRef, useState } from 'react'
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
  className = '',
  width = '100%',
  height = 500,
}: DependencyGraphProps) {
  const svgRef = useRef<SVGSVGElement>(null)
  const containerRef = useRef<HTMLDivElement>(null)
  const [dimensions, setDimensions] = useState({ width: 800, height: 500 })

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

  // Main D3 rendering effect
  useEffect(() => {
    if (!svgRef.current || !data) return

    const svg = d3.select(svgRef.current)
    const { width: svgWidth, height: svgHeight } = dimensions

    // Clear previous content
    svg.selectAll('*').remove()

    // Transform data to D3 format
    const { nodes, links } = transformToD3Data(data)

    // If no nodes, don't render anything
    if (nodes.length === 0) return

    // Create main group for zoom/pan
    const g = svg.append('g')

    // Arrow marker definition for directed edges
    svg
      .append('defs')
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
      .attr('fill', '#94a3b8') // slate-400

    // Create link elements
    const link = g
      .append('g')
      .attr('class', 'links')
      .selectAll<SVGLineElement, D3Link>('line')
      .data(links)
      .join('line')
      .attr('stroke', '#94a3b8') // slate-400
      .attr('stroke-opacity', 0.6)
      .attr('stroke-width', 1.5)
      .attr('marker-end', 'url(#arrowhead)')

    // Create node group elements
    const node = g
      .append('g')
      .attr('class', 'nodes')
      .selectAll<SVGGElement, D3Node>('g')
      .data(nodes)
      .join('g')
      .attr('class', 'node')

    // Add circles to nodes
    node
      .append('circle')
      .attr('r', (d) => Math.max(8, Math.min(16, 8 + d.dependencyCount * 0.5)))
      .attr('fill', '#6366f1') // indigo-500
      .attr('stroke', '#ffffff')
      .attr('stroke-width', 2)
      .style('cursor', 'pointer')

    // Add labels to nodes
    node
      .append('text')
      .text((d) => truncatePackageName(d.name))
      .attr('font-size', '11px')
      .attr('font-family', 'system-ui, sans-serif')
      .attr('fill', '#374151') // gray-700
      .attr('dx', 14)
      .attr('dy', 4)
      .style('pointer-events', 'none')
      .style('user-select', 'none')

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
      link
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
      svg.selectAll('*').remove() // Clean DOM
    }
  }, [data, dimensions])

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
      <svg
        ref={svgRef}
        className="h-full w-full"
        style={{
          width: '100%',
          height: '100%',
        }}
      />
    </div>
  )
})

// Re-export types and utilities
export type { D3GraphData, D3Link, D3Node, DependencyGraphProps } from './types'
export { transformToD3Data, truncatePackageName } from './useForceSimulation'
