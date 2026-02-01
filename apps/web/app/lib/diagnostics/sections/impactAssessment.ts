import type { CircularDependencyInfo, DependencyGraph } from '@monoguard/types'
import type { ImpactAssessmentDetails, RippleNode } from '../types'

/**
 * Generate impact assessment details for the diagnostic report
 * AC5: Impact Assessment
 */
export function generateImpactAssessment(
  cycle: CircularDependencyInfo,
  graph: DependencyGraph,
  totalPackages: number
): ImpactAssessmentDetails {
  // Use existing impact assessment from the cycle if available
  if (cycle.impactAssessment) {
    const ia = cycle.impactAssessment
    return {
      directParticipants: ia.directParticipants,
      directParticipantsCount: ia.directParticipants.length,
      indirectDependents: ia.indirectDependents.map((d) => d.packageName),
      indirectDependentsCount: ia.indirectDependents.length,
      totalAffectedCount: ia.totalAffected,
      percentageOfMonorepo: Math.round(ia.affectedPercentage * 100),
      riskLevel: ia.riskLevel,
      riskExplanation: ia.riskExplanation,
      rippleEffectTree: buildRippleTreeFromEffect(ia.directParticipants, ia.rippleEffect),
    }
  }

  // Fallback: compute from graph
  const cyclePackages = getCyclePackages(cycle)
  const indirectDependents = findIndirectDependents(cyclePackages, graph)
  const totalAffected = new Set([...cyclePackages, ...indirectDependents]).size
  const percentageOfMonorepo =
    totalPackages > 0 ? Math.round((totalAffected / totalPackages) * 100) : 0
  const riskLevel = classifyRisk(totalAffected, totalPackages, cyclePackages)
  const riskExplanation = generateRiskExplanation(riskLevel, totalAffected, percentageOfMonorepo)
  const rippleEffectTree = buildRippleTree(cyclePackages, graph)

  return {
    directParticipants: cyclePackages,
    directParticipantsCount: cyclePackages.length,
    indirectDependents,
    indirectDependentsCount: indirectDependents.length,
    totalAffectedCount: totalAffected,
    percentageOfMonorepo,
    riskLevel,
    riskExplanation,
    rippleEffectTree,
  }
}

function getCyclePackages(cycle: CircularDependencyInfo): string[] {
  const packages = cycle.cycle
  if (packages.length > 1 && packages[packages.length - 1] === packages[0]) {
    return packages.slice(0, -1)
  }
  return packages
}

function findIndirectDependents(cyclePackages: string[], graph: DependencyGraph): string[] {
  const visited = new Set<string>(cyclePackages)
  const queue = [...cyclePackages]
  const indirectDependents: string[] = []

  while (queue.length > 0) {
    const current = queue.shift()!

    // Find packages that depend on current (reverse edges)
    const dependents = graph.edges.filter((edge) => edge.to === current).map((edge) => edge.from)

    for (const dep of dependents) {
      if (!visited.has(dep)) {
        visited.add(dep)
        indirectDependents.push(dep)
        queue.push(dep)
      }
    }
  }

  return indirectDependents
}

function classifyRisk(
  totalAffected: number,
  totalPackages: number,
  cyclePackages: string[]
): 'critical' | 'high' | 'medium' | 'low' {
  if (totalPackages === 0) return 'low'
  const percentage = (totalAffected / totalPackages) * 100

  if (percentage > 50) return 'critical'
  if (percentage > 25) return 'high'
  if (cyclePackages.some((p) => p.includes('core') || p.includes('shared'))) {
    return 'high'
  }
  if (percentage > 10) return 'medium'
  return 'low'
}

function generateRiskExplanation(
  riskLevel: string,
  totalAffected: number,
  percentage: number
): string {
  const descriptions: Record<string, string> = {
    critical:
      `Critical risk: ${totalAffected} packages (${percentage}% of monorepo) are affected. ` +
      `This cycle impacts core infrastructure and should be prioritized immediately.`,
    high:
      `High risk: ${totalAffected} packages (${percentage}% of monorepo) are affected. ` +
      `This cycle has significant downstream impact and should be addressed soon.`,
    medium:
      `Medium risk: ${totalAffected} packages (${percentage}% of monorepo) are affected. ` +
      `This cycle has moderate impact and should be scheduled for resolution.`,
    low:
      `Low risk: ${totalAffected} packages (${percentage}% of monorepo) are affected. ` +
      `This cycle has limited blast radius but should still be fixed to improve architecture health.`,
  }
  return descriptions[riskLevel] || 'Unknown risk level.'
}

function buildRippleTree(
  cyclePackages: string[],
  graph: DependencyGraph,
  maxDepth: number = 3
): RippleNode {
  const root: RippleNode = {
    package: 'Cycle',
    depth: 0,
    dependents: [],
  }

  const visited = new Set<string>(cyclePackages)
  for (const pkg of cyclePackages) {
    const packageNode = buildPackageRippleTree(pkg, graph, 1, maxDepth, visited)
    root.dependents.push(packageNode)
  }

  return root
}

function buildPackageRippleTree(
  pkg: string,
  graph: DependencyGraph,
  depth: number,
  maxDepth: number,
  visited: Set<string>
): RippleNode {
  const node: RippleNode = {
    package: pkg,
    depth,
    dependents: [],
  }

  if (depth >= maxDepth) return node

  const directDependents = graph.edges
    .filter((edge) => edge.to === pkg && !visited.has(edge.from))
    .map((edge) => edge.from)

  for (const dep of directDependents.slice(0, 5)) {
    visited.add(dep)
    node.dependents.push(buildPackageRippleTree(dep, graph, depth + 1, maxDepth, visited))
  }

  return node
}

function buildRippleTreeFromEffect(
  directParticipants: string[],
  rippleEffect?: {
    layers: { distance: number; packages: string[]; count: number }[]
    totalLayers: number
  }
): RippleNode {
  const root: RippleNode = {
    package: 'Cycle',
    depth: 0,
    dependents: directParticipants.map((pkg) => ({
      package: pkg,
      depth: 1,
      dependents: [],
    })),
  }

  if (!rippleEffect) return root

  // Add indirect dependents from ripple layers
  for (const layer of rippleEffect.layers) {
    if (layer.distance <= 1) continue
    for (const pkg of layer.packages.slice(0, 5)) {
      root.dependents.push({
        package: pkg,
        depth: layer.distance,
        dependents: [],
      })
    }
  }

  return root
}
