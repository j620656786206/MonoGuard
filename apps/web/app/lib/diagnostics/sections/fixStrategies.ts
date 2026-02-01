import type { CircularDependencyInfo } from '@monoguard/types'
import type { FixStrategyGuide, FixStrategyStep } from '../types'

/**
 * Render fix strategies as formatted guides for the diagnostic report
 * AC4: All Fix Strategies with Full Guides
 */
export function renderFixStrategies(cycle: CircularDependencyInfo): FixStrategyGuide[] {
  if (!cycle.fixStrategies || cycle.fixStrategies.length === 0) {
    return []
  }

  return cycle.fixStrategies.map((strategy) => {
    const steps: FixStrategyStep[] = strategy.guide
      ? strategy.guide.steps.map((step) => ({
          number: step.number,
          title: step.title,
          description: step.description,
          codeSnippet: step.codeAfter?.code,
          filePath: step.filePath,
          isOptional: false,
        }))
      : []

    const codeSnippets = {
      before:
        strategy.beforeAfterExplanation?.importDiffs?.[0]?.importsToRemove?.[0]?.statement ?? '',
      after: strategy.beforeAfterExplanation?.importDiffs?.[0]?.importsToAdd?.[0]?.statement ?? '',
    }

    return {
      strategy: strategy.type,
      title: strategy.name,
      description: strategy.description,
      suitabilityScore: strategy.suitability,
      estimatedEffort: strategy.effort,
      estimatedTime: strategy.guide?.estimatedTime ?? estimateTimeFromEffort(strategy.effort),
      pros: strategy.pros,
      cons: strategy.cons,
      steps,
      codeSnippets,
    }
  })
}

function estimateTimeFromEffort(effort: string): string {
  switch (effort) {
    case 'low':
      return '15-30 minutes'
    case 'medium':
      return '1-2 hours'
    case 'high':
      return '2-4 hours'
    default:
      return 'Unknown'
  }
}
