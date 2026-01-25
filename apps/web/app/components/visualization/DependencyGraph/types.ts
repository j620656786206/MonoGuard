/**
 * D3-specific types for the DependencyGraph visualization
 *
 * These types extend D3's simulation types for use with our dependency graph data.
 */

import type { DependencyType } from '@monoguard/types'
import type { SimulationLinkDatum, SimulationNodeDatum } from 'd3'

/**
 * D3Node - Node representation for D3 force simulation
 *
 * Extends D3's SimulationNodeDatum to include package-specific data.
 * D3 will add x, y, vx, vy, fx, fy during simulation.
 */
export interface D3Node extends SimulationNodeDatum {
  /** Package identifier (usually the package name) */
  id: string
  /** Display name for the node */
  name: string
  /** Relative path from workspace root */
  path: string
  /** Total number of dependencies (used for node sizing) */
  dependencyCount: number
}

/**
 * D3Link - Link representation for D3 force simulation
 *
 * Extends D3's SimulationLinkDatum to include dependency-specific data.
 */
export interface D3Link extends SimulationLinkDatum<D3Node> {
  /** Source node id or node object (D3 replaces string with node ref) */
  source: string | D3Node
  /** Target node id or node object (D3 replaces string with node ref) */
  target: string | D3Node
  /** Type of dependency relationship */
  type: DependencyType
}

/**
 * D3GraphData - Transformed graph data ready for D3 consumption
 */
export interface D3GraphData {
  nodes: D3Node[]
  links: D3Link[]
}

/**
 * DependencyGraphProps - Props for the DependencyGraph component
 */
export interface DependencyGraphProps {
  /** Dependency graph data from analysis results */
  data: import('@monoguard/types').DependencyGraph
  /** Optional class name for styling */
  className?: string
  /** Width of the graph container (default: 100%) */
  width?: number | string
  /** Height of the graph container (default: 500px) */
  height?: number | string
}

/**
 * ForceSimulationConfig - Configuration for D3 force simulation
 */
export interface ForceSimulationConfig {
  /** Distance between linked nodes */
  linkDistance: number
  /** Charge strength (negative = repel, positive = attract) */
  chargeStrength: number
  /** Collision radius for nodes */
  collisionRadius: number
  /** Alpha decay rate for simulation stabilization */
  alphaDecay: number
}

/**
 * Default configuration for force simulation
 */
export const DEFAULT_SIMULATION_CONFIG: ForceSimulationConfig = {
  linkDistance: 100,
  chargeStrength: -200,
  collisionRadius: 30,
  alphaDecay: 0.02,
}
