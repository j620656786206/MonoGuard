import type { CircularDependencyInfo } from '@monoguard/types'
import { describe, expect, it } from 'vitest'
import { generateExecutiveSummary } from '../sections/executiveSummary'

const baseCycle: CircularDependencyInfo = {
  cycle: ['pkg-a', 'pkg-b', 'pkg-c', 'pkg-a'],
  type: 'indirect',
  severity: 'warning',
  depth: 3,
  impact: 'Moderate impact',
  complexity: 5,
  priorityScore: 5,
}

describe('generateExecutiveSummary', () => {
  it('should generate summary with all required fields', () => {
    const summary = generateExecutiveSummary(baseCycle)
    expect(summary.description).toBeTruthy()
    expect(summary.severity).toMatch(/^(critical|high|medium|low)$/)
    expect(summary.recommendation).toBeTruthy()
    expect(summary.estimatedEffort).toMatch(/^(low|medium|high)$/)
    expect(summary.cycleLength).toBe(3)
  })

  it('should classify 2-package cycle as low severity', () => {
    const simpleCycle: CircularDependencyInfo = {
      ...baseCycle,
      cycle: ['pkg-a', 'pkg-b', 'pkg-a'],
      type: 'direct',
      depth: 2,
      complexity: 2,
      priorityScore: 3,
    }
    const summary = generateExecutiveSummary(simpleCycle)
    expect(summary.severity).toBe('low')
  })

  it('should classify cycle with core package as critical', () => {
    const coreCycle: CircularDependencyInfo = {
      ...baseCycle,
      cycle: ['core-lib', 'pkg-b', 'pkg-c', 'core-lib'],
    }
    const summary = generateExecutiveSummary(coreCycle)
    expect(summary.severity).toBe('critical')
  })

  it('should classify cycle with shared package as critical', () => {
    const sharedCycle: CircularDependencyInfo = {
      ...baseCycle,
      cycle: ['shared-utils', 'pkg-b', 'shared-utils'],
    }
    const summary = generateExecutiveSummary(sharedCycle)
    expect(summary.severity).toBe('critical')
  })

  it('should classify > 5 package cycle as critical', () => {
    const largeCycle: CircularDependencyInfo = {
      ...baseCycle,
      cycle: ['a', 'b', 'c', 'd', 'e', 'f', 'a'],
      depth: 6,
    }
    const summary = generateExecutiveSummary(largeCycle)
    expect(summary.severity).toBe('critical')
  })

  it('should classify 4-package cycle as high severity', () => {
    const fourPkgCycle: CircularDependencyInfo = {
      ...baseCycle,
      cycle: ['a', 'b', 'c', 'd', 'a'],
      depth: 4,
      priorityScore: 5,
    }
    const summary = generateExecutiveSummary(fourPkgCycle)
    expect(summary.severity).toBe('high')
  })

  it('should classify 3-package cycle as medium severity', () => {
    const summary = generateExecutiveSummary(baseCycle)
    expect(summary.severity).toBe('medium')
  })

  it('should estimate low effort for simple 2-package cycle', () => {
    const simpleCycle: CircularDependencyInfo = {
      ...baseCycle,
      cycle: ['pkg-a', 'pkg-b', 'pkg-a'],
      complexity: 2,
      priorityScore: 3,
    }
    const summary = generateExecutiveSummary(simpleCycle)
    expect(summary.estimatedEffort).toBe('low')
  })

  it('should estimate high effort for complex cycles', () => {
    const complexCycle: CircularDependencyInfo = {
      ...baseCycle,
      cycle: ['a', 'b', 'c', 'd', 'e', 'a'],
      complexity: 8,
    }
    const summary = generateExecutiveSummary(complexCycle)
    expect(summary.estimatedEffort).toBe('high')
  })

  it('should generate description for direct cycle', () => {
    const directCycle: CircularDependencyInfo = {
      ...baseCycle,
      cycle: ['pkg-a', 'pkg-b', 'pkg-a'],
      type: 'direct',
    }
    const summary = generateExecutiveSummary(directCycle)
    expect(summary.description).toContain('Direct circular dependency')
  })

  it('should generate description for indirect cycle', () => {
    const summary = generateExecutiveSummary(baseCycle)
    expect(summary.description).toContain('Indirect circular dependency')
  })

  it('should use fix strategy in recommendation when available', () => {
    const withStrategy: CircularDependencyInfo = {
      ...baseCycle,
      fixStrategies: [
        {
          type: 'extract-module',
          name: 'Extract Shared Module',
          description: 'Move shared code',
          suitability: 8,
          effort: 'low',
          pros: [],
          cons: [],
          recommended: true,
          targetPackages: ['pkg-a'],
        },
      ],
    }
    const summary = generateExecutiveSummary(withStrategy)
    expect(summary.recommendation).toContain('Extract Shared Module')
  })

  it('should use impact assessment for affected count when available', () => {
    const withImpact: CircularDependencyInfo = {
      ...baseCycle,
      impactAssessment: {
        directParticipants: ['a', 'b', 'c'],
        indirectDependents: [
          { packageName: 'd', dependsOn: 'a', distance: 1, dependencyPath: ['d', 'a'] },
          { packageName: 'e', dependsOn: 'b', distance: 1, dependencyPath: ['e', 'b'] },
        ],
        totalAffected: 5,
        affectedPercentage: 0.25,
        affectedPercentageDisplay: '25%',
        riskLevel: 'high',
        riskExplanation: 'High risk',
      },
    }
    const summary = generateExecutiveSummary(withImpact)
    expect(summary.affectedPackagesCount).toBe(5)
  })

  it('should handle non-repeating cycle array', () => {
    const nonRepeating: CircularDependencyInfo = {
      ...baseCycle,
      cycle: ['pkg-a', 'pkg-b', 'pkg-c'],
    }
    const summary = generateExecutiveSummary(nonRepeating)
    expect(summary.cycleLength).toBe(3)
    expect(summary.description).toBeTruthy()
  })
})
