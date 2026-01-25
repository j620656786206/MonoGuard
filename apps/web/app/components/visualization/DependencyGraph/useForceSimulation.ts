/**
 * useForceSimulation - Custom hook for D3 force simulation
 *
 * Manages D3 force simulation lifecycle with proper React integration.
 * Handles initialization, updates, and cleanup.
 */

import type { DependencyGraph } from '@monoguard/types'
import type { Simulation } from 'd3'
import * as d3 from 'd3'
import { useCallback, useEffect, useRef } from 'react'
import type { D3GraphData, D3Link, D3Node, ForceSimulationConfig } from './types'
import { DEFAULT_SIMULATION_CONFIG } from './types'

/**
 * Transform DependencyGraph data to D3-compatible format
 */
export function transformToD3Data(data: DependencyGraph): D3GraphData {
  const nodes: D3Node[] = Object.entries(data.nodes).map(([name, pkg]) => ({
    id: name,
    name: pkg.name,
    path: pkg.path,
    dependencyCount:
      pkg.dependencies.length + pkg.devDependencies.length + pkg.peerDependencies.length,
  }))

  const links: D3Link[] = data.edges.map((edge) => ({
    source: edge.from,
    target: edge.to,
    type: edge.type,
  }))

  return { nodes, links }
}

/**
 * Truncate package name for display
 * Shows the last part of scoped package names (e.g., "@app/core" -> "core")
 */
export function truncatePackageName(name: string, maxLength: number = 15): string {
  // For scoped packages, get the part after the last /
  const parts = name.split('/')
  const displayName = parts[parts.length - 1] || name

  if (displayName.length <= maxLength) {
    return displayName
  }

  return `${displayName.substring(0, maxLength - 3)}...`
}

interface UseForceSimulationOptions {
  svgRef: React.RefObject<SVGSVGElement | null>
  data: DependencyGraph
  config?: Partial<ForceSimulationConfig>
  onTick?: () => void
}

interface UseForceSimulationReturn {
  simulation: Simulation<D3Node, D3Link> | null
  graphData: D3GraphData
}

/**
 * Hook to manage D3 force simulation
 *
 * @param options - Configuration options for the simulation
 * @returns The simulation instance and transformed graph data
 */
export function useForceSimulation({
  svgRef,
  data,
  config = {},
  onTick,
}: UseForceSimulationOptions): UseForceSimulationReturn {
  const simulationRef = useRef<Simulation<D3Node, D3Link> | null>(null)
  const graphDataRef = useRef<D3GraphData>({ nodes: [], links: [] })

  // Merge config with defaults
  const mergedConfig: ForceSimulationConfig = {
    ...DEFAULT_SIMULATION_CONFIG,
    ...config,
  }

  // Transform data
  const graphData = transformToD3Data(data)
  graphDataRef.current = graphData

  // Create simulation setup function
  const setupSimulation = useCallback(
    (width: number, height: number) => {
      const { nodes, links } = graphDataRef.current

      // Stop any existing simulation
      if (simulationRef.current) {
        simulationRef.current.stop()
      }

      // Create new simulation
      const simulation = d3
        .forceSimulation<D3Node>(nodes)
        .force(
          'link',
          d3
            .forceLink<D3Node, D3Link>(links)
            .id((d) => d.id)
            .distance(mergedConfig.linkDistance)
        )
        .force('charge', d3.forceManyBody<D3Node>().strength(mergedConfig.chargeStrength))
        .force('center', d3.forceCenter(width / 2, height / 2))
        .force('collision', d3.forceCollide<D3Node>().radius(mergedConfig.collisionRadius))
        .alphaDecay(mergedConfig.alphaDecay)

      if (onTick) {
        simulation.on('tick', onTick)
      }

      simulationRef.current = simulation
      return simulation
    },
    [
      mergedConfig.linkDistance,
      mergedConfig.chargeStrength,
      mergedConfig.collisionRadius,
      mergedConfig.alphaDecay,
      onTick,
    ]
  )

  // Effect to manage simulation lifecycle
  // biome-ignore lint/correctness/useExhaustiveDependencies: data is intentionally included to re-run simulation when graph data changes
  useEffect(() => {
    if (!svgRef.current) return

    const width = svgRef.current.clientWidth || 800
    const height = svgRef.current.clientHeight || 500

    const simulation = setupSimulation(width, height)

    // Cleanup function
    return () => {
      simulation.stop()
      simulationRef.current = null
    }
  }, [svgRef, setupSimulation, data])

  return {
    simulation: simulationRef.current,
    graphData,
  }
}
