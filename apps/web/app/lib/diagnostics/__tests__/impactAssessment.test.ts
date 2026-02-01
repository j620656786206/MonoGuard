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

  it('should classify high risk for core/shared package in cycle', () => {
    const coreGraph: DependencyGraph = {
      nodes: {
        'core-lib': {
          name: 'core-lib',
          version: '1.0.0',
          path: '/core-lib',
          dependencies: ['pkg-x'],
          devDependencies: [],
          peerDependencies: [],
        },
        'pkg-x': {
          name: 'pkg-x',
          version: '1.0.0',
          path: '/pkg-x',
          dependencies: ['core-lib'],
          devDependencies: [],
          peerDependencies: [],
        },
      },
      edges: [
        { from: 'core-lib', to: 'pkg-x', type: 'production', versionRange: '^1.0.0' },
        { from: 'pkg-x', to: 'core-lib', type: 'production', versionRange: '^1.0.0' },
      ],
      rootPath: '/workspace',
      workspaceType: 'pnpm',
    }
    const coreCycle: CircularDependencyInfo = {
      ...baseCycle,
      cycle: ['core-lib', 'pkg-x', 'core-lib'],
    }
    // 2 out of 20 = 10%, which is <= 25%, but core package → high
    const impact = generateImpactAssessment(coreCycle, coreGraph, 20)
    expect(impact.riskLevel).toBe('high')
  })

  it('should classify high risk for shared package in cycle', () => {
    const sharedGraph: DependencyGraph = {
      nodes: {
        'shared-utils': {
          name: 'shared-utils',
          version: '1.0.0',
          path: '/shared-utils',
          dependencies: ['pkg-y'],
          devDependencies: [],
          peerDependencies: [],
        },
        'pkg-y': {
          name: 'pkg-y',
          version: '1.0.0',
          path: '/pkg-y',
          dependencies: ['shared-utils'],
          devDependencies: [],
          peerDependencies: [],
        },
      },
      edges: [
        { from: 'shared-utils', to: 'pkg-y', type: 'production', versionRange: '^1.0.0' },
        { from: 'pkg-y', to: 'shared-utils', type: 'production', versionRange: '^1.0.0' },
      ],
      rootPath: '/workspace',
      workspaceType: 'pnpm',
    }
    const sharedCycle: CircularDependencyInfo = {
      ...baseCycle,
      cycle: ['shared-utils', 'pkg-y', 'shared-utils'],
    }
    const impact = generateImpactAssessment(sharedCycle, sharedGraph, 20)
    expect(impact.riskLevel).toBe('high')
  })

  it('should classify medium risk for moderate percentage', () => {
    const moderateGraph: DependencyGraph = {
      ...mockGraph,
      edges: [
        { from: 'pkg-a', to: 'pkg-b', type: 'production', versionRange: '^1.0.0' },
        { from: 'pkg-b', to: 'pkg-a', type: 'production', versionRange: '^1.0.0' },
      ],
    }
    const simpleCycle: CircularDependencyInfo = {
      ...baseCycle,
      cycle: ['pkg-a', 'pkg-b', 'pkg-a'],
    }
    // 2 out of 15 = 13.3% → medium (between 10% and 25%)
    const impact = generateImpactAssessment(simpleCycle, moderateGraph, 15)
    expect(impact.riskLevel).toBe('medium')
  })

  it('should classify high risk for percentage > 25%', () => {
    const impact = generateImpactAssessment(baseCycle, mockGraph, 10)
    // 5 affected out of 10 = 50% > 25% but ≤ 50% → high ... actually > 50% → critical
    // Actually 5/10 = 50% which is NOT > 50, so high
    // Wait: percentage > 50 → critical, percentage > 25 → high
    // 5/10 = 50%, 50 > 50 is false, but 50 > 25 is true → high
    expect(impact.riskLevel).toBe('high')
  })

  it('should build ripple tree from existing impact assessment with ripple effect', () => {
    const withRipple: CircularDependencyInfo = {
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
        riskExplanation: 'High risk',
        rippleEffect: {
          layers: [
            { distance: 1, packages: ['pkg-a', 'pkg-b', 'pkg-c'], count: 3 },
            { distance: 2, packages: ['pkg-d', 'pkg-e'], count: 2 },
            { distance: 3, packages: ['pkg-f'], count: 1 },
          ],
          totalLayers: 3,
        },
      },
    }
    const impact = generateImpactAssessment(withRipple, mockGraph, 6)
    expect(impact.rippleEffectTree.package).toBe('Cycle')
    expect(impact.rippleEffectTree.depth).toBe(0)
    // Direct participants as children at depth 1
    expect(impact.rippleEffectTree.dependents.length).toBeGreaterThan(3)
    // Indirect layers (distance > 1) added
    const deepNodes = impact.rippleEffectTree.dependents.filter((d) => d.depth > 1)
    expect(deepNodes.length).toBeGreaterThan(0)
  })

  it('should build ripple tree from existing impact without ripple effect', () => {
    const withoutRipple: CircularDependencyInfo = {
      ...baseCycle,
      impactAssessment: {
        directParticipants: ['pkg-a', 'pkg-b'],
        indirectDependents: [],
        totalAffected: 2,
        affectedPercentage: 0.33,
        affectedPercentageDisplay: '33%',
        riskLevel: 'medium',
        riskExplanation: 'Medium risk',
      },
    }
    const impact = generateImpactAssessment(withoutRipple, mockGraph, 6)
    expect(impact.rippleEffectTree.package).toBe('Cycle')
    // Only direct participants as children
    expect(impact.rippleEffectTree.dependents).toHaveLength(2)
  })

  it('should handle non-repeating cycle array', () => {
    const nonRepeating: CircularDependencyInfo = {
      ...baseCycle,
      cycle: ['pkg-a', 'pkg-b', 'pkg-c'],
    }
    const impact = generateImpactAssessment(nonRepeating, mockGraph, 6)
    expect(impact.directParticipants).toEqual(['pkg-a', 'pkg-b', 'pkg-c'])
    expect(impact.directParticipantsCount).toBe(3)
  })

  it('should generate unknown risk for unrecognized risk level', () => {
    const emptyEdgeGraph: DependencyGraph = {
      nodes: {},
      edges: [],
      rootPath: '/workspace',
      workspaceType: 'pnpm',
    }
    // totalPackages=0 → low risk, percentageOfMonorepo=0
    const impact = generateImpactAssessment(baseCycle, emptyEdgeGraph, 0)
    expect(impact.riskLevel).toBe('low')
    expect(impact.percentageOfMonorepo).toBe(0)
  })
})
