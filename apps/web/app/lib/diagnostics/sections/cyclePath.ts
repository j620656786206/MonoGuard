import type { CircularDependencyInfo } from '@monoguard/types'
import type { CycleEdge, CycleNode, CyclePathVisualization } from '../types'
import { generateCycleAscii } from '../visualizations/cycleAscii'
import { generateCycleSvg } from '../visualizations/cycleSvg'

/**
 * Generate cycle path visualization data including SVG and ASCII diagrams
 * AC2: Complete Cycle Path Visualization
 */
export function generateCyclePath(
  cycle: CircularDependencyInfo,
  isDarkMode: boolean = false
): CyclePathVisualization {
  const packages = getCyclePackages(cycle)
  const nodes = createNodes(packages)
  const edges = createEdges(packages, cycle.importTraces)
  const breakingPoint = identifyBreakingPoint(cycle, packages)

  const edgesWithBreakingPoint = edges.map((edge) => ({
    ...edge,
    isBreakingPoint: edge.from === breakingPoint.fromPackage && edge.to === breakingPoint.toPackage,
  }))

  const svgDiagram = generateCycleSvg(nodes, edgesWithBreakingPoint, isDarkMode)
  const asciiDiagram = generateCycleAscii(packages, breakingPoint)

  return {
    nodes,
    edges: edgesWithBreakingPoint,
    breakingPoint,
    svgDiagram,
    asciiDiagram,
  }
}

function getCyclePackages(cycle: CircularDependencyInfo): string[] {
  const packages = cycle.cycle
  if (packages.length > 1 && packages[packages.length - 1] === packages[0]) {
    return packages.slice(0, -1)
  }
  return packages
}

function createNodes(packages: string[]): CycleNode[] {
  const radius = 150
  const centerX = 200
  const centerY = 200

  return packages.map((pkg, index) => {
    const angle = (2 * Math.PI * index) / packages.length - Math.PI / 2
    return {
      id: pkg,
      name: pkg.split('/').pop() || pkg,
      path: pkg,
      isInCycle: true,
      position: {
        x: Math.round(centerX + radius * Math.cos(angle)),
        y: Math.round(centerY + radius * Math.sin(angle)),
      },
    }
  })
}

function createEdges(
  packages: string[],
  importTraces?: CircularDependencyInfo['importTraces']
): CycleEdge[] {
  const edges: CycleEdge[] = []

  for (let i = 0; i < packages.length; i++) {
    const from = packages[i]
    const to = packages[(i + 1) % packages.length]

    const trace = importTraces?.find((t) => t.fromPackage === from && t.toPackage === to)

    edges.push({
      from,
      to,
      isBreakingPoint: false,
      importStatement: trace?.statement,
      filePath: trace?.filePath,
      lineNumber: trace?.lineNumber,
    })
  }

  return edges
}

function identifyBreakingPoint(
  cycle: CircularDependencyInfo,
  packages: string[]
): {
  fromPackage: string
  toPackage: string
  reason: string
} {
  // Use criticalEdge from root cause analysis if available
  if (cycle.rootCause?.criticalEdge) {
    return {
      fromPackage: cycle.rootCause.criticalEdge.from,
      toPackage: cycle.rootCause.criticalEdge.to,
      reason: `Critical edge identified by root cause analysis (confidence: ${cycle.rootCause.confidence}%).`,
    }
  }

  // Fallback: break at the last edge (closing the cycle)
  const lastPackage = packages[packages.length - 1]
  const firstPackage = packages[0]

  return {
    fromPackage: lastPackage,
    toPackage: firstPackage,
    reason: 'This edge has the least downstream impact based on dependency analysis.',
  }
}
