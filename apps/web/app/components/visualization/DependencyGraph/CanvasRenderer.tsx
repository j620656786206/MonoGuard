/**
 * CanvasRenderer - Canvas-based graph renderer for large graphs (Story 4.9)
 *
 * Uses HTML5 Canvas 2D context for rendering graphs with >= 500 nodes.
 * Supports HiDPI displays, circular dependency highlighting, and directed edges.
 *
 * Performance targets (AC6):
 * - Initial render < 3 seconds for 1000+ nodes
 * - Frame rate >= 30fps during simulation
 * - Interactions respond in < 100ms
 *
 * @see Story 4.9: Implement Hybrid SVG/Canvas Rendering
 */
'use client'

import * as d3 from 'd3'
import React, { useCallback, useEffect, useRef } from 'react'
import { EDGE_COLORS, NODE_COLORS } from './styles'
import type { CanvasRendererProps, D3Link, D3Node, ViewportState } from './types'
import { DEFAULT_SIMULATION_CONFIG } from './types'
import { useCanvasInteraction } from './useCanvasInteraction'

/**
 * Draw an arrowhead at the end of an edge
 */
function drawArrow(
  ctx: CanvasRenderingContext2D,
  x1: number,
  y1: number,
  x2: number,
  y2: number,
  isCircular: boolean
) {
  const angle = Math.atan2(y2 - y1, x2 - x1)
  const nodeRadius = 10
  const endX = x2 - nodeRadius * Math.cos(angle)
  const endY = y2 - nodeRadius * Math.sin(angle)

  const arrowLength = 8
  const arrowWidth = 5

  ctx.save()
  ctx.translate(endX, endY)
  ctx.rotate(angle)

  ctx.beginPath()
  ctx.moveTo(0, 0)
  ctx.lineTo(-arrowLength, -arrowWidth)
  ctx.lineTo(-arrowLength, arrowWidth)
  ctx.closePath()
  ctx.fillStyle = isCircular ? EDGE_COLORS.cycle.stroke : EDGE_COLORS.normal.stroke
  ctx.fill()
  ctx.restore()
}

/**
 * Truncate a label for display on canvas
 */
function truncateLabel(label: string, maxLength = 12): string {
  const shortName = label.split('/').pop() ?? label
  return shortName.length > maxLength ? `${shortName.substring(0, maxLength - 1)}...` : shortName
}

export const CanvasRenderer = React.memo(function CanvasRenderer({
  nodes,
  links,
  circularNodeIds,
  circularEdgePairs,
  viewport,
  onViewportChange,
  selectedNodeId,
  onNodeSelect,
  onNodeHover,
  width,
  height,
}: CanvasRendererProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null)
  const simulationRef = useRef<d3.Simulation<D3Node, D3Link> | null>(null)
  const nodesRef = useRef<D3Node[]>(nodes)
  const linksRef = useRef<D3Link[]>(links)
  const viewportRef = useRef<ViewportState>(viewport)
  const selectedNodeIdRef = useRef<string | null>(selectedNodeId)
  const animationFrameRef = useRef<number | null>(null)

  // Keep refs in sync
  useEffect(() => {
    nodesRef.current = nodes
  }, [nodes])

  useEffect(() => {
    linksRef.current = links
  }, [links])

  useEffect(() => {
    viewportRef.current = viewport
  }, [viewport])

  useEffect(() => {
    selectedNodeIdRef.current = selectedNodeId
  }, [selectedNodeId])

  // Canvas interaction hook for hover/click
  const { handleMouseMove, handleMouseClick } = useCanvasInteraction({
    canvasRef,
    nodesRef,
    viewport,
    onNodeHover,
    onNodeSelect,
  })

  // Render function using refs (called from animation frame)
  const render = useCallback(() => {
    if (!canvasRef.current) return

    const canvas = canvasRef.current
    const ctx = canvas.getContext('2d')
    if (!ctx) return

    const dpr = window.devicePixelRatio || 1
    const displayWidth = width
    const displayHeight = height

    // HiDPI support
    if (canvas.width !== displayWidth * dpr || canvas.height !== displayHeight * dpr) {
      canvas.width = displayWidth * dpr
      canvas.height = displayHeight * dpr
      canvas.style.width = `${displayWidth}px`
      canvas.style.height = `${displayHeight}px`
    }

    const currentViewport = viewportRef.current
    const currentNodes = nodesRef.current
    const currentLinks = linksRef.current
    const currentSelectedNodeId = selectedNodeIdRef.current

    ctx.save()
    ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
    ctx.clearRect(0, 0, displayWidth, displayHeight)

    // Apply viewport transform
    ctx.translate(currentViewport.panX, currentViewport.panY)
    ctx.scale(currentViewport.zoom, currentViewport.zoom)

    // Draw edges
    for (const link of currentLinks) {
      const source = link.source as D3Node
      const target = link.target as D3Node
      if (source.x == null || source.y == null || target.x == null || target.y == null) continue

      const edgeKey = `${source.id}->${target.id}`
      const isCircular = circularEdgePairs.has(edgeKey)

      ctx.beginPath()
      ctx.moveTo(source.x, source.y)
      ctx.lineTo(target.x, target.y)
      ctx.strokeStyle = isCircular ? EDGE_COLORS.cycle.stroke : EDGE_COLORS.normal.stroke
      ctx.lineWidth = isCircular ? EDGE_COLORS.cycle.width : EDGE_COLORS.normal.width
      ctx.globalAlpha = isCircular ? EDGE_COLORS.cycle.opacity : EDGE_COLORS.normal.opacity
      ctx.stroke()
      ctx.globalAlpha = 1

      // Draw arrow
      drawArrow(ctx, source.x, source.y, target.x, target.y, isCircular)
    }

    // Draw nodes
    for (const node of currentNodes) {
      if (node.x == null || node.y == null) continue

      const isCircular = circularNodeIds.has(node.id)
      const isSelected = node.id === currentSelectedNodeId
      const nodeRadius = Math.max(8, Math.min(16, 8 + node.dependencyCount * 0.5))

      // Node circle
      ctx.beginPath()
      ctx.arc(node.x, node.y, isSelected ? nodeRadius + 2 : nodeRadius, 0, 2 * Math.PI)

      if (isCircular) {
        ctx.fillStyle = NODE_COLORS.cycle.fill
        ctx.fill()
        ctx.strokeStyle = NODE_COLORS.cycle.stroke
        ctx.lineWidth = 3
        ctx.stroke()
      } else {
        ctx.fillStyle = isSelected ? '#3b82f6' : NODE_COLORS.normal.fill
        ctx.fill()
        ctx.strokeStyle = NODE_COLORS.normal.stroke
        ctx.lineWidth = 2
        ctx.stroke()
      }

      // Selection ring
      if (isSelected) {
        ctx.beginPath()
        ctx.arc(node.x, node.y, nodeRadius + 6, 0, 2 * Math.PI)
        ctx.strokeStyle = '#60a5fa'
        ctx.lineWidth = 3
        ctx.stroke()
      }

      // Node label
      ctx.fillStyle = '#1f2937'
      ctx.font = '10px Inter, system-ui, sans-serif'
      ctx.textAlign = 'center'
      ctx.fillText(truncateLabel(node.name), node.x, node.y + nodeRadius + 14)
    }

    ctx.restore()
  }, [width, height, circularNodeIds, circularEdgePairs])

  // Force simulation setup
  useEffect(() => {
    if (nodes.length === 0) return

    // Create D3 force simulation
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
      .force('center', d3.forceCenter(width / 2, height / 2))
      .force(
        'collision',
        d3.forceCollide<D3Node>().radius(DEFAULT_SIMULATION_CONFIG.collisionRadius)
      )
      .alphaDecay(DEFAULT_SIMULATION_CONFIG.alphaDecay)

    simulationRef.current = simulation

    // On each tick, request a render
    simulation.on('tick', () => {
      if (animationFrameRef.current) {
        cancelAnimationFrame(animationFrameRef.current)
      }
      animationFrameRef.current = requestAnimationFrame(render)
    })

    // CRITICAL: Cleanup
    return () => {
      simulation.stop()
      simulationRef.current = null
      if (animationFrameRef.current) {
        cancelAnimationFrame(animationFrameRef.current)
        animationFrameRef.current = null
      }
    }
  }, [nodes, links, width, height, render])

  // Re-render when viewport or selection changes (without restarting simulation)
  // biome-ignore lint/correctness/useExhaustiveDependencies: viewport and selectedNodeId are intentional triggers for re-render (read via refs inside render)
  useEffect(() => {
    render()
  }, [viewport, selectedNodeId, render])

  // Zoom handling via mouse wheel
  useEffect(() => {
    if (!canvasRef.current) return

    const canvas = canvasRef.current

    const handleWheel = (e: WheelEvent) => {
      e.preventDefault()
      const scaleFactor = e.deltaY > 0 ? 0.9 : 1.1
      const newZoom = Math.max(0.1, Math.min(4, viewport.zoom * scaleFactor))

      // Zoom toward mouse position
      const rect = canvas.getBoundingClientRect()
      const mouseX = e.clientX - rect.left
      const mouseY = e.clientY - rect.top

      const newPanX = mouseX - ((mouseX - viewport.panX) / viewport.zoom) * newZoom
      const newPanY = mouseY - ((mouseY - viewport.panY) / viewport.zoom) * newZoom

      onViewportChange({ zoom: newZoom, panX: newPanX, panY: newPanY })
    }

    canvas.addEventListener('wheel', handleWheel, { passive: false })

    return () => {
      canvas.removeEventListener('wheel', handleWheel)
    }
  }, [viewport, onViewportChange])

  // Pan handling via mouse drag
  useEffect(() => {
    if (!canvasRef.current) return

    const canvas = canvasRef.current
    let isPanning = false
    let startX = 0
    let startY = 0
    let startPanX = 0
    let startPanY = 0

    const handleMouseDown = (e: MouseEvent) => {
      // Only pan on middle-click or when not on a node
      if (e.button === 1 || (e.button === 0 && !e.shiftKey)) {
        // Check if we're clicking a node - if so, don't pan
        const rect = canvas.getBoundingClientRect()
        const x = e.clientX - rect.left
        const y = e.clientY - rect.top
        const graphX = (x - viewport.panX) / viewport.zoom
        const graphY = (y - viewport.panY) / viewport.zoom

        const nodes = nodesRef.current
        let onNode = false
        if (nodes) {
          for (let i = nodes.length - 1; i >= 0; i--) {
            const node = nodes[i]
            if (node.x == null || node.y == null) continue
            const dx = graphX - node.x
            const dy = graphY - node.y
            if (Math.sqrt(dx * dx + dy * dy) <= 12) {
              onNode = true
              break
            }
          }
        }

        if (!onNode) {
          isPanning = true
          startX = e.clientX
          startY = e.clientY
          startPanX = viewport.panX
          startPanY = viewport.panY
          canvas.style.cursor = 'grabbing'
        }
      }
    }

    const handleMouseMoveForPan = (e: MouseEvent) => {
      if (!isPanning) return
      const dx = e.clientX - startX
      const dy = e.clientY - startY
      onViewportChange({
        ...viewport,
        panX: startPanX + dx,
        panY: startPanY + dy,
      })
    }

    const handleMouseUp = () => {
      if (isPanning) {
        isPanning = false
        canvas.style.cursor = 'crosshair'
      }
    }

    canvas.addEventListener('mousedown', handleMouseDown)
    window.addEventListener('mousemove', handleMouseMoveForPan)
    window.addEventListener('mouseup', handleMouseUp)

    return () => {
      canvas.removeEventListener('mousedown', handleMouseDown)
      window.removeEventListener('mousemove', handleMouseMoveForPan)
      window.removeEventListener('mouseup', handleMouseUp)
    }
  }, [viewport, onViewportChange])

  return (
    <canvas
      ref={canvasRef}
      className="cursor-crosshair"
      style={{
        width: `${width}px`,
        height: `${height}px`,
        touchAction: 'none',
      }}
      onMouseMove={handleMouseMove}
      onClick={handleMouseClick}
    />
  )
})
