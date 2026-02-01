import type { CircularDependencyInfo, DependencyGraph } from '@monoguard/types'
import { describe, expect, it } from 'vitest'
import { generateImpactAssessment } from '../sections/impactAssessment'

const mockGraph: DependencyGraph = {
  nodes: {
    'pkg-a': {
      name: 'pkg-a',
      version: '1.0.0',
      path: '/a',
      dependencies: ['pkg-b'],
      devDependencies: [],
      peerDependencies: [],
    },
    'pkg-b': {
      name: 'pkg-b',
      version: '1.0.0',
      path: '/b',
      dependencies: ['pkg-c'],
      devDependencies: [],
      peerDependencies: [],
    },
    'pkg-c': {
      name: 'pkg-c',
      version: '1.0.0',
      path: '/c',
      dependencies: ['pkg-a'],
      devDependencies: [],
      peerDependencies: [],
    },
    'pkg-d': {
      name: 'pkg-d',
      version: '1.0.0',
      path: '/d',
      dependencies: ['pkg-a'],
      devDependencies: [],
      peerDependencies: [],
    },
    'pkg-e': {
      name: 'pkg-e',
      version: '1.0.0',
      path: '/e',
      dependencies: ['pkg-d'],
      devDependencies: [],
      peerDependencies: [],
    },
    'pkg-f': {
      name: 'pkg-f',
      version: '1.0.0',
      path: '/f',
      dependencies: [],
      devDependencies: [],
      peerDependencies: [],
    },
  },
  edges: [
    { from: 'pkg-a', to: 'pkg-b', type: 'production', versionRange: '^1.0.0' },
    { from: 'pkg-b', to: 'pkg-c', type: 'production', versionRange: '^1.0.0' },
    { from: 'pkg-c', to: 'pkg-a', type: 'production', versionRange: '^1.0.0' },
    { from: 'pkg-d', to: 'pkg-a', type: 'production', versionRange: '^1.0.0' },
    { from: 'pkg-e', to: 'pkg-d', type: 'production', versionRange: '^1.0.0' },
  ],
  rootPath: '/workspace',
  workspaceType: 'pnpm',
}

const baseCycle: CircularDependencyInfo = {
  cycle: ['pkg-a', 'pkg-b', 'pkg-c', 'pkg-a'],
  type: 'indirect',
  severity: 'warning',
  depth: 3,
  impact: 'Test',
  complexity: 5,
  priorityScore: 5,
}

describe('generateImpactAssessment', () => {
  it('should identify direct participants from cycle', () => {
    const impact = generateImpactAssessment(baseCycle, mockGraph, 6)
    expect(impact.directParticipants).toEqual(['pkg-a', 'pkg-b', 'pkg-c'])
    expect(impact.directParticipantsCount).toBe(3)
  })

  it('should find indirect dependents', () => {
    const impact = generateImpactAssessment(baseCycle, mockGraph, 6)
    expect(impact.indirectDependents).toContain('pkg-d')
    expect(impact.indirectDependents).toContain('pkg-e')
    expect(impact.indirectDependentsCount).toBe(2)
  })

  it('should calculate total affected correctly', () => {
    const impact = generateImpactAssessment(baseCycle, mockGraph, 6)
    expect(impact.totalAffectedCount).toBe(5)
  })

  it('should calculate percentage of monorepo', () => {
    const impact = generateImpactAssessment(baseCycle, mockGraph, 6)
    expect(impact.percentageOfMonorepo).toBe(83)
  })

  it('should classify risk level', () => {
    const impact = generateImpactAssessment(baseCycle, mockGraph, 6)
    // 83% > 50% → critical
    expect(impact.riskLevel).toBe('critical')
  })

  it('should provide risk explanation', () => {
    const impact = generateImpactAssessment(baseCycle, mockGraph, 6)
    expect(impact.riskExplanation).toContain('5 packages')
    expect(impact.riskExplanation).toContain('83%')
  })

  it('should build ripple effect tree', () => {
    const impact = generateImpactAssessment(baseCycle, mockGraph, 6)
    expect(impact.rippleEffectTree.package).toBe('Cycle')
    expect(impact.rippleEffectTree.depth).toBe(0)
    expect(impact.rippleEffectTree.dependents.length).toBeGreaterThan(0)
  })

  it('should use existing impact assessment when available', () => {
    const withImpact: CircularDependencyInfo = {
      ...baseCycle,
      impactAssessment: {
        directParticipants: ['pkg-a', 'pkg-b', 'pkg-c'],
        indirectDependents: [
          {
            packageName: 'pkg-d',
            dependsOn: 'pkg-a',
            distance: 1,
            dependencyPath: ['pkg-d', 'pkg-a'],
          },
        ],
        totalAffected: 4,
        affectedPercentage: 0.67,
        affectedPercentageDisplay: '67%',
        riskLevel: 'high',
        riskExplanation: 'High risk from existing assessment',
      },
    }
    const impact = generateImpactAssessment(withImpact, mockGraph, 6)
    expect(impact.riskLevel).toBe('high')
    expect(impact.riskExplanation).toBe('High risk from existing assessment')
    expect(impact.totalAffectedCount).toBe(4)
  })

  it('should handle empty graph', () => {
    const emptyGraph: DependencyGraph = {
      nodes: {},
      edges: [],
      rootPath: '/workspace',
      workspaceType: 'pnpm',
    }
    const impact = generateImpactAssessment(baseCycle, emptyGraph, 0)
    expect(impact.directParticipantsCount).toBe(3)
    expect(impact.indirectDependentsCount).toBe(0)
    expect(impact.percentageOfMonorepo).toBe(0)
  })

  it('should classify low risk for isolated cycles', () => {
    const isolatedGraph: DependencyGraph = {
      ...mockGraph,
      edges: [
        { from: 'pkg-a', to: 'pkg-b', type: 'production', versionRange: '^1.0.0' },
        { from: 'pkg-b', to: 'pkg-a', type: 'production', versionRange: '^1.0.0' },
      ],
    }
    const smallCycle: CircularDependencyInfo = {
      ...baseCycle,
      cycle: ['pkg-a', 'pkg-b', 'pkg-a'],
    }
    // 2 affected out of 100 packages = 2% → low
    const impact = generateImpactAssessment(smallCycle, isolatedGraph, 100)
    expect(impact.riskLevel).toBe('low')
  })
})
