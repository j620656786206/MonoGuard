import type { CircularDependencyInfo } from '@monoguard/types'
import type { RelatedCycleInfo } from '../types'
import { getCyclePackages } from '../types'

/**
 * Find cycles that share packages with the target cycle
 * AC6: Related Cycles Detection
 */
export function findRelatedCycles(
  targetCycle: CircularDependencyInfo,
  allCycles: CircularDependencyInfo[]
): RelatedCycleInfo[] {
  const targetPackages = new Set(getCyclePackages(targetCycle))
  const relatedCycles: RelatedCycleInfo[] = []

  for (let i = 0; i < allCycles.length; i++) {
    const otherCycle = allCycles[i]
    // Skip the target cycle itself
    if (otherCycle === targetCycle) continue

    const otherPackages = getCyclePackages(otherCycle)
    const sharedPackages = otherPackages.filter((pkg) => targetPackages.has(pkg))

    if (sharedPackages.length === 0) continue

    const overlapPercentage = Math.round(
      (sharedPackages.length / Math.max(targetPackages.size, otherPackages.length)) * 100
    )

    const recommendFixTogether = overlapPercentage >= 30 || sharedPackages.length >= 2

    const stableCycleId = otherPackages.map((p) => p.split('/').pop() || p).join('-')

    relatedCycles.push({
      cycleId: stableCycleId,
      sharedPackages,
      overlapPercentage,
      recommendFixTogether,
      reason: recommendFixTogether
        ? `These cycles share ${sharedPackages.length} package(s): ${sharedPackages.join(', ')}. Fixing them together can reduce total refactoring effort.`
        : `These cycles share ${sharedPackages.length} package(s): ${sharedPackages.join(', ')}.`,
    })
  }

  return relatedCycles
}
