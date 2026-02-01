import type { CircularDependencyInfo } from '@monoguard/types'
import { describe, expect, it } from 'vitest'
import { renderFixStrategies } from '../sections/fixStrategies'

const baseCycle: CircularDependencyInfo = {
  cycle: ['pkg-a', 'pkg-b', 'pkg-c', 'pkg-a'],
  type: 'indirect',
  severity: 'warning',
  depth: 3,
  impact: 'Test',
  complexity: 5,
  priorityScore: 5,
}

describe('renderFixStrategies', () => {
  it('should return empty array when no strategies', () => {
    const result = renderFixStrategies(baseCycle)
    expect(result).toHaveLength(0)
  })

  it('should render fix strategies from cycle data', () => {
    const withStrategies: CircularDependencyInfo = {
      ...baseCycle,
      fixStrategies: [
        {
          type: 'extract-module',
          name: 'Extract Shared Module',
          description: 'Move shared code',
          suitability: 8,
          effort: 'medium',
          pros: ['Clean separation', 'Testable'],
          cons: ['New package to maintain'],
          recommended: true,
          targetPackages: ['pkg-a', 'pkg-b'],
        },
        {
          type: 'dependency-injection',
          name: 'Dependency Injection',
          description: 'Invert the dependency',
          suitability: 6,
          effort: 'low',
          pros: ['Simple'],
          cons: ['Less clear'],
          recommended: false,
          targetPackages: ['pkg-a'],
        },
      ],
    }
    const result = renderFixStrategies(withStrategies)
    expect(result).toHaveLength(2)
    expect(result[0].strategy).toBe('extract-module')
    expect(result[0].title).toBe('Extract Shared Module')
    expect(result[0].suitabilityScore).toBe(8)
    expect(result[0].pros).toEqual(['Clean separation', 'Testable'])
    expect(result[0].cons).toEqual(['New package to maintain'])
    expect(result[1].strategy).toBe('dependency-injection')
  })

  it('should include steps from fix guide', () => {
    const withGuide: CircularDependencyInfo = {
      ...baseCycle,
      fixStrategies: [
        {
          type: 'extract-module',
          name: 'Extract',
          description: 'Extract code',
          suitability: 8,
          effort: 'medium',
          pros: [],
          cons: [],
          recommended: true,
          targetPackages: [],
          guide: {
            strategyType: 'extract-module',
            title: 'Guide Title',
            summary: 'Guide summary',
            steps: [
              {
                number: 1,
                title: 'Create new package',
                description: 'Create a new shared package',
                filePath: 'packages/shared/package.json',
                codeAfter: { language: 'json', code: '{ "name": "shared" }' },
              },
              {
                number: 2,
                title: 'Move shared code',
                description: 'Move shared code to new package',
              },
            ],
            verification: [],
            estimatedTime: '30-45 minutes',
          },
        },
      ],
    }
    const result = renderFixStrategies(withGuide)
    expect(result[0].steps).toHaveLength(2)
    expect(result[0].steps[0].title).toBe('Create new package')
    expect(result[0].steps[0].codeSnippet).toBe('{ "name": "shared" }')
    expect(result[0].estimatedTime).toBe('30-45 minutes')
  })

  it('should estimate time from effort level when no guide', () => {
    const strategies: CircularDependencyInfo = {
      ...baseCycle,
      fixStrategies: [
        {
          type: 'extract-module',
          name: 'Low Effort',
          description: 'Simple fix',
          suitability: 5,
          effort: 'low',
          pros: [],
          cons: [],
          recommended: true,
          targetPackages: [],
        },
        {
          type: 'dependency-injection',
          name: 'Medium Effort',
          description: 'Moderate fix',
          suitability: 5,
          effort: 'medium',
          pros: [],
          cons: [],
          recommended: false,
          targetPackages: [],
        },
        {
          type: 'boundary-refactoring',
          name: 'High Effort',
          description: 'Complex fix',
          suitability: 5,
          effort: 'high',
          pros: [],
          cons: [],
          recommended: false,
          targetPackages: [],
        },
      ],
    }
    const result = renderFixStrategies(strategies)
    expect(result[0].estimatedTime).toBe('15-30 minutes')
    expect(result[1].estimatedTime).toBe('1-2 hours')
    expect(result[2].estimatedTime).toBe('2-4 hours')
  })

  it('should include code snippets from beforeAfterExplanation', () => {
    const withCodeSnippets: CircularDependencyInfo = {
      ...baseCycle,
      fixStrategies: [
        {
          type: 'extract-module',
          name: 'Extract',
          description: 'Extract shared code',
          suitability: 8,
          effort: 'medium',
          pros: [],
          cons: [],
          recommended: true,
          targetPackages: [],
          beforeAfterExplanation: {
            currentState: { nodes: [], edges: [], cycleResolved: false },
            proposedState: { nodes: [], edges: [], cycleResolved: true },
            packageJsonDiffs: [],
            importDiffs: [
              {
                filePath: 'src/index.ts',
                packageName: 'pkg-a',
                importsToRemove: [
                  { statement: "import { foo } from 'pkg-b'", fromPackage: 'pkg-b' },
                ],
                importsToAdd: [
                  { statement: "import { foo } from 'shared'", fromPackage: 'shared' },
                ],
              },
            ],
            explanation: {
              summary: 'Test',
              whyItWorks: 'Test',
              highLevelChanges: [],
              confidence: 0.9,
            },
            warnings: [],
          },
        },
      ],
    }
    const result = renderFixStrategies(withCodeSnippets)
    expect(result[0].codeSnippets.before).toBe("import { foo } from 'pkg-b'")
    expect(result[0].codeSnippets.after).toBe("import { foo } from 'shared'")
  })
})
