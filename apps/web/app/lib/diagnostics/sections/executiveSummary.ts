import type { CircularDependencyInfo, EffortLevel } from '@monoguard/types'
import type { ExecutiveSummary } from '../types'

/**
 * Generate an executive summary for a circular dependency
 * AC1: Executive Summary Generation
 */
export function generateExecutiveSummary(cycle: CircularDependencyInfo): ExecutiveSummary {
  const cyclePackages = getCyclePackages(cycle)
  const cycleLength = cyclePackages.length
  const severity = classifySeverity(cycle)
  const effort = estimateEffort(cycle)
  const description = generateDescription(cycle, cyclePackages)
  const recommendation = generateRecommendation(cycle, severity)
  const affectedPackagesCount = cycle.impactAssessment?.totalAffected ?? cycleLength

  return {
    description,
    severity,
    recommendation,
    estimatedEffort: effort,
    affectedPackagesCount,
    cycleLength,
  }
}

/** Extract unique packages from cycle (cycle array repeats first element) */
function getCyclePackages(cycle: CircularDependencyInfo): string[] {
  const packages = cycle.cycle
  // cycle array typically ends with first package repeated; deduplicate
  if (packages.length > 1 && packages[packages.length - 1] === packages[0]) {
    return packages.slice(0, -1)
  }
  return packages
}

function classifySeverity(cycle: CircularDependencyInfo): 'critical' | 'high' | 'medium' | 'low' {
  const packages = getCyclePackages(cycle)

  // Critical: core/shared packages or > 5 packages in cycle
  if (packages.some((p) => p.includes('core') || p.includes('shared'))) {
    return 'critical'
  }
  if (packages.length > 5) {
    return 'critical'
  }

  // High: 4-5 packages or high priority score
  if (packages.length >= 4 || cycle.priorityScore > 7) {
    return 'high'
  }

  // Medium: 3 packages
  if (packages.length === 3) {
    return 'medium'
  }

  // Low: 2 packages (direct cycle)
  return 'low'
}

function estimateEffort(cycle: CircularDependencyInfo): EffortLevel {
  const packages = getCyclePackages(cycle)

  // Simple 2-package cycle with low complexity
  if (packages.length === 2 && cycle.complexity < 3) {
    return 'low'
  }

  // Large or complex cycles
  if (packages.length > 4 || cycle.complexity > 7) {
    return 'high'
  }

  return 'medium'
}

function generateDescription(cycle: CircularDependencyInfo, packages: string[]): string {
  const packageList = packages
    .slice(0, 3)
    .map((p) => `\`${p}\``)
    .join(', ')
  const andMore = packages.length > 3 ? ` and ${packages.length - 3} more` : ''

  if (cycle.type === 'direct') {
    return (
      `Direct circular dependency between ${packageList}${andMore}. ` +
      `These packages import each other, creating a tight coupling that should be resolved.`
    )
  }

  return (
    `Indirect circular dependency involving ${packageList}${andMore}. ` +
    `This ${packages.length}-package cycle creates complex inter-dependencies that affect architecture health.`
  )
}

function generateRecommendation(cycle: CircularDependencyInfo, severity: string): string {
  const bestStrategy = cycle.fixStrategies?.[0]

  if (bestStrategy) {
    return (
      `Recommended fix: ${bestStrategy.name}. ` +
      `This is a ${severity}-severity issue with ${bestStrategy.effort} estimated effort.`
    )
  }

  const packages = getCyclePackages(cycle)
  if (packages.length === 2) {
    return (
      `Recommend using dependency injection to break the direct dependency. ` +
      `Consider which package should own the shared functionality.`
    )
  }

  return (
    `Recommend extracting shared code into a new package to eliminate the cycle. ` +
    `This will improve architecture clarity and testability.`
  )
}
