import { describe, expect, it } from 'vitest'
import type {
  AnalysisResult,
  CheckResult,
  CircularDependencyInfo,
  ComplexityBreakdown,
  ComplexityFactor,
  DependencyEdge,
  DependencyGraph,
  EffortLevel,
  FixStrategy,
  FixStrategyType,
  ImportTrace,
  ImportType,
  PackageNode,
  RefactoringComplexity,
  RootCauseAnalysis,
  RootCauseEdge,
} from '../analysis'

describe('Analysis types', () => {
  describe('AnalysisResult', () => {
    it('can be instantiated with valid data', () => {
      const graph: DependencyGraph = {
        nodes: {
          '@monoguard/types': {
            name: '@monoguard/types',
            version: '0.1.0',
            path: 'packages/types',
            dependencies: [],
            devDependencies: ['typescript'],
            peerDependencies: [],
          },
        },
        edges: [],
        rootPath: '/workspace',
        workspaceType: 'pnpm',
      }

      const result: AnalysisResult = {
        healthScore: 85,
        packageCount: 10,
        circularDependencies: [],
        graph,
        metadata: {
          version: '0.1.0',
          durationMs: 1500,
          filesProcessed: 50,
          workspaceType: 'pnpm',
        },
        createdAt: '2026-01-15T10:30:00Z',
      }

      expect(result.healthScore).toBe(85)
      expect(result.graph.workspaceType).toBe('pnpm')
      expect(result.metadata.durationMs).toBe(1500)
    })
  })

  describe('DependencyGraph', () => {
    it('supports Record<string, PackageNode> for nodes', () => {
      const graph: DependencyGraph = {
        nodes: {
          'package-a': {
            name: 'package-a',
            version: '1.0.0',
            path: 'packages/a',
            dependencies: ['package-b'],
            devDependencies: [],
            peerDependencies: [],
          },
          'package-b': {
            name: 'package-b',
            version: '1.0.0',
            path: 'packages/b',
            dependencies: [],
            devDependencies: [],
            peerDependencies: [],
          },
        },
        edges: [
          {
            from: 'package-a',
            to: 'package-b',
            type: 'production',
            versionRange: '^1.0.0',
          },
        ],
        rootPath: '/workspace',
        workspaceType: 'npm',
      }

      expect(Object.keys(graph.nodes)).toHaveLength(2)
      expect(graph.edges).toHaveLength(1)
    })
  })

  describe('CircularDependencyInfo', () => {
    it('can represent direct circular dependency', () => {
      const circular: CircularDependencyInfo = {
        cycle: ['package-a', 'package-b', 'package-a'],
        type: 'direct',
        severity: 'critical',
        depth: 2,
        impact: 'Build failure due to circular dependency',
        complexity: 5,
      }

      expect(circular.type).toBe('direct')
      expect(circular.severity).toBe('critical')
    })

    it('can include fix strategies array (Story 3.3)', () => {
      const circular: CircularDependencyInfo = {
        cycle: ['package-a', 'package-b', 'package-c', 'package-a'],
        type: 'indirect',
        severity: 'warning',
        depth: 3,
        impact: 'Potential build issues',
        complexity: 7,
        fixStrategies: [
          {
            type: 'extract-module',
            name: 'Extract Shared Module',
            description: 'Create a new shared package to hold common dependencies.',
            suitability: 8,
            effort: 'medium',
            pros: [
              'Creates clear separation of concerns',
              'Isolates shared code between package-a, package-b, package-c',
            ],
            cons: ['Introduces a new package to maintain'],
            recommended: true,
            targetPackages: ['package-a', 'package-b', 'package-c'],
            newPackageName: '@mono/shared',
          },
          {
            type: 'dependency-injection',
            name: 'Dependency Injection',
            description: 'Invert the problematic dependency by introducing an interface.',
            suitability: 6,
            effort: 'medium',
            pros: ['Minimal code changes required', 'Preserves existing structure'],
            cons: ['Adds indirection to the codebase'],
            recommended: false,
            targetPackages: ['package-b', 'package-c'],
          },
          {
            type: 'boundary-refactoring',
            name: 'Module Boundary Refactoring',
            description: 'Restructure package boundaries to eliminate the cycle.',
            suitability: 5,
            effort: 'high',
            pros: ['Addresses root architectural issue', 'Cleaner long-term design'],
            cons: ['Requires significant refactoring effort', 'May affect external consumers'],
            recommended: false,
            targetPackages: ['package-a', 'package-b', 'package-c'],
          },
        ],
      }

      expect(circular.fixStrategies).toHaveLength(3)
      expect(circular.fixStrategies?.[0].type).toBe('extract-module')
      expect(circular.fixStrategies?.[0].recommended).toBe(true)
      expect(circular.fixStrategies?.[0].suitability).toBe(8)
      expect(circular.fixStrategies?.[0].effort).toBe('medium')
      expect(circular.fixStrategies?.[0].newPackageName).toBe('@mono/shared')
    })
  })

  describe('CheckResult', () => {
    it('can represent passing check', () => {
      const checkResult: CheckResult = {
        passed: true,
        errors: [],
        warnings: [],
        healthScore: 95,
      }

      expect(checkResult.passed).toBe(true)
      expect(checkResult.errors).toHaveLength(0)
    })

    it('can represent failing check with errors', () => {
      const checkResult: CheckResult = {
        passed: false,
        errors: [
          {
            code: 'CIRCULAR_DETECTED',
            message: 'Circular dependency found: A -> B -> A',
            file: 'packages/a/package.json',
          },
        ],
        warnings: [
          {
            code: 'LOW_HEALTH_SCORE',
            message: 'Health score below threshold',
          },
        ],
        healthScore: 45,
      }

      expect(checkResult.passed).toBe(false)
      expect(checkResult.errors).toHaveLength(1)
      expect(checkResult.warnings).toHaveLength(1)
    })
  })

  describe('PackageNode', () => {
    it('correctly types all dependency categories', () => {
      const node: PackageNode = {
        name: '@monoguard/web',
        version: '1.0.0',
        path: 'apps/web',
        dependencies: ['react', 'react-dom'],
        devDependencies: ['typescript', 'vitest'],
        peerDependencies: ['react'],
      }

      expect(node.dependencies).toContain('react')
      expect(node.devDependencies).toContain('typescript')
      expect(node.peerDependencies).toContain('react')
    })
  })

  describe('DependencyEdge', () => {
    it('supports all dependency types', () => {
      const edges: DependencyEdge[] = [
        { from: 'a', to: 'b', type: 'production', versionRange: '^1.0.0' },
        { from: 'a', to: 'c', type: 'development', versionRange: '*' },
        { from: 'a', to: 'd', type: 'peer', versionRange: '>=16.0.0' },
        { from: 'a', to: 'e', type: 'optional', versionRange: '~2.0.0' },
      ]

      expect(edges[0].type).toBe('production')
      expect(edges[1].type).toBe('development')
      expect(edges[2].type).toBe('peer')
      expect(edges[3].type).toBe('optional')
    })
  })
})

describe('WorkspaceType', () => {
  it('supports all workspace types', () => {
    const graph1: DependencyGraph = {
      nodes: {},
      edges: [],
      rootPath: '/',
      workspaceType: 'npm',
    }
    const graph2: DependencyGraph = {
      nodes: {},
      edges: [],
      rootPath: '/',
      workspaceType: 'yarn',
    }
    const graph3: DependencyGraph = {
      nodes: {},
      edges: [],
      rootPath: '/',
      workspaceType: 'pnpm',
    }
    const graph4: DependencyGraph = {
      nodes: {},
      edges: [],
      rootPath: '/',
      workspaceType: 'nx',
    }
    const graph5: DependencyGraph = {
      nodes: {},
      edges: [],
      rootPath: '/',
      workspaceType: 'unknown',
    }

    expect(graph1.workspaceType).toBe('npm')
    expect(graph2.workspaceType).toBe('yarn')
    expect(graph3.workspaceType).toBe('pnpm')
    expect(graph4.workspaceType).toBe('nx')
    expect(graph5.workspaceType).toBe('unknown')
  })
})

// Story 3.1: Root Cause Analysis Types
describe('RootCauseAnalysis', () => {
  it('can represent root cause analysis with all fields', () => {
    const rootCause: RootCauseAnalysis = {
      originatingPackage: 'pkg-ui',
      problematicDependency: {
        from: 'pkg-ui',
        to: 'pkg-api',
        type: 'production',
        critical: false,
      },
      confidence: 82,
      explanation: "Package 'pkg-ui' is highly likely the root cause.",
      chain: [
        { from: 'pkg-ui', to: 'pkg-api', type: 'production', critical: false },
        { from: 'pkg-api', to: 'pkg-core', type: 'production', critical: false },
        { from: 'pkg-core', to: 'pkg-ui', type: 'production', critical: true },
      ],
      criticalEdge: {
        from: 'pkg-core',
        to: 'pkg-ui',
        type: 'production',
        critical: true,
      },
    }

    expect(rootCause.originatingPackage).toBe('pkg-ui')
    expect(rootCause.confidence).toBe(82)
    expect(rootCause.chain).toHaveLength(3)
    expect(rootCause.criticalEdge?.critical).toBe(true)
  })

  it('criticalEdge is optional', () => {
    const rootCause: RootCauseAnalysis = {
      originatingPackage: 'pkg-a',
      problematicDependency: {
        from: 'pkg-a',
        to: 'pkg-b',
        type: 'development',
        critical: false,
      },
      confidence: 65,
      explanation: "Package 'pkg-a' is likely the root cause.",
      chain: [
        { from: 'pkg-a', to: 'pkg-b', type: 'development', critical: false },
        { from: 'pkg-b', to: 'pkg-a', type: 'production', critical: true },
      ],
    }

    expect(rootCause.criticalEdge).toBeUndefined()
    expect(rootCause.chain).toHaveLength(2)
  })
})

describe('RootCauseEdge', () => {
  it('supports all dependency types', () => {
    const edges: RootCauseEdge[] = [
      { from: 'a', to: 'b', type: 'production', critical: false },
      { from: 'a', to: 'c', type: 'development', critical: true },
      { from: 'a', to: 'd', type: 'peer', critical: false },
      { from: 'a', to: 'e', type: 'optional', critical: true },
    ]

    expect(edges[0].type).toBe('production')
    expect(edges[1].type).toBe('development')
    expect(edges[2].type).toBe('peer')
    expect(edges[3].type).toBe('optional')
  })

  it('critical flag indicates edge importance', () => {
    const criticalEdge: RootCauseEdge = {
      from: 'pkg-core',
      to: 'pkg-ui',
      type: 'production',
      critical: true,
    }

    const normalEdge: RootCauseEdge = {
      from: 'pkg-ui',
      to: 'pkg-api',
      type: 'production',
      critical: false,
    }

    expect(criticalEdge.critical).toBe(true)
    expect(normalEdge.critical).toBe(false)
  })
})

describe('CircularDependencyInfo with RootCause', () => {
  it('can include root cause analysis', () => {
    const circular: CircularDependencyInfo = {
      cycle: ['pkg-ui', 'pkg-api', 'pkg-core', 'pkg-ui'],
      type: 'indirect',
      severity: 'info',
      depth: 3,
      impact: 'Indirect circular dependency involving 3 packages',
      complexity: 5,
      rootCause: {
        originatingPackage: 'pkg-ui',
        problematicDependency: {
          from: 'pkg-ui',
          to: 'pkg-api',
          type: 'production',
          critical: false,
        },
        confidence: 82,
        explanation: "Package 'pkg-ui' is highly likely the root cause.",
        chain: [
          { from: 'pkg-ui', to: 'pkg-api', type: 'production', critical: false },
          { from: 'pkg-api', to: 'pkg-core', type: 'production', critical: false },
          { from: 'pkg-core', to: 'pkg-ui', type: 'production', critical: true },
        ],
      },
    }

    expect(circular.rootCause).toBeDefined()
    expect(circular.rootCause?.originatingPackage).toBe('pkg-ui')
    expect(circular.rootCause?.confidence).toBe(82)
  })

  it('rootCause is optional for backward compatibility', () => {
    const circular: CircularDependencyInfo = {
      cycle: ['package-a', 'package-b', 'package-a'],
      type: 'direct',
      severity: 'warning',
      depth: 2,
      impact: 'Direct circular dependency',
      complexity: 3,
    }

    expect(circular.rootCause).toBeUndefined()
  })
})

// Story 3.2: Import Trace Types
describe('ImportTrace', () => {
  it('can represent a named ESM import', () => {
    const trace: ImportTrace = {
      fromPackage: 'pkg-ui',
      toPackage: 'pkg-api',
      filePath: 'src/components/UserList.tsx',
      lineNumber: 5,
      statement: "import { fetchUsers } from '@mono/api'",
      importType: 'esm-named',
      symbols: ['fetchUsers'],
    }

    expect(trace.importType).toBe('esm-named')
    expect(trace.symbols).toContain('fetchUsers')
    expect(trace.lineNumber).toBe(5)
  })

  it('supports all import types', () => {
    const importTypes: ImportType[] = [
      'esm-named',
      'esm-default',
      'esm-namespace',
      'esm-side-effect',
      'esm-dynamic',
      'cjs-require',
    ]

    const traces: ImportTrace[] = importTypes.map((type, i) => ({
      fromPackage: 'pkg-a',
      toPackage: 'pkg-b',
      filePath: `src/file${i}.ts`,
      lineNumber: i + 1,
      statement: `import statement ${i}`,
      importType: type,
    }))

    expect(traces.map((t) => t.importType)).toEqual(importTypes)
  })

  it('symbols is optional for side-effect imports', () => {
    const trace: ImportTrace = {
      fromPackage: 'pkg-ui',
      toPackage: 'pkg-styles',
      filePath: 'src/index.ts',
      lineNumber: 1,
      statement: "import './styles.css'",
      importType: 'esm-side-effect',
    }

    expect(trace.symbols).toBeUndefined()
  })
})

// Story 3.3: Fix Strategy Types
describe('FixStrategy', () => {
  it('can represent extract-module strategy', () => {
    const strategy: FixStrategy = {
      type: 'extract-module',
      name: 'Extract Shared Module',
      description: 'Create a new shared package to hold common dependencies.',
      suitability: 9,
      effort: 'medium',
      pros: [
        'Creates clear separation of concerns',
        'Effectively breaks complex multi-package cycle',
      ],
      cons: ['Introduces a new package to maintain'],
      recommended: true,
      targetPackages: ['pkg-ui', 'pkg-api', 'pkg-core'],
      newPackageName: '@mono/shared',
    }

    expect(strategy.type).toBe('extract-module')
    expect(strategy.suitability).toBe(9)
    expect(strategy.newPackageName).toBe('@mono/shared')
    expect(strategy.recommended).toBe(true)
  })

  it('can represent dependency-injection strategy', () => {
    const strategy: FixStrategy = {
      type: 'dependency-injection',
      name: 'Dependency Injection',
      description: 'Invert the problematic dependency.',
      suitability: 8,
      effort: 'low',
      pros: ['Minimal code changes required', 'Clear injection point: pkg-b → pkg-a'],
      cons: ['Adds indirection to the codebase'],
      recommended: false,
      targetPackages: ['pkg-a', 'pkg-b'],
    }

    expect(strategy.type).toBe('dependency-injection')
    expect(strategy.effort).toBe('low')
    expect(strategy.newPackageName).toBeUndefined()
  })

  it('can represent boundary-refactoring strategy', () => {
    const strategy: FixStrategy = {
      type: 'boundary-refactoring',
      name: 'Module Boundary Refactoring',
      description: 'Restructure package boundaries.',
      suitability: 7,
      effort: 'high',
      pros: [
        'Addresses root architectural issue',
        'Opportunity to properly define core package boundaries',
      ],
      cons: ['Requires significant refactoring effort', 'May affect external consumers'],
      recommended: false,
      targetPackages: ['pkg-ui', 'pkg-core'],
    }

    expect(strategy.type).toBe('boundary-refactoring')
    expect(strategy.effort).toBe('high')
  })
})

describe('FixStrategyType', () => {
  it('supports all strategy types', () => {
    const types: FixStrategyType[] = [
      'extract-module',
      'dependency-injection',
      'boundary-refactoring',
    ]

    expect(types).toHaveLength(3)
    expect(types).toContain('extract-module')
    expect(types).toContain('dependency-injection')
    expect(types).toContain('boundary-refactoring')
  })
})

describe('EffortLevel', () => {
  it('supports all effort levels', () => {
    const levels: EffortLevel[] = ['low', 'medium', 'high']

    expect(levels).toHaveLength(3)
    expect(levels).toContain('low')
    expect(levels).toContain('medium')
    expect(levels).toContain('high')
  })

  it('effort levels correspond to time ranges', () => {
    // Documentation validation: low < 1 hour, medium 1-4 hours, high > 4 hours
    const effortMapping: Record<EffortLevel, string> = {
      low: '< 1 hour',
      medium: '1-4 hours',
      high: '> 4 hours',
    }

    expect(Object.keys(effortMapping)).toHaveLength(3)
  })
})

describe('CircularDependencyInfo with FixStrategies', () => {
  it('can include complete fix strategies from Story 3.3', () => {
    const circular: CircularDependencyInfo = {
      cycle: ['pkg-ui', 'pkg-api', 'pkg-core', 'pkg-ui'],
      type: 'indirect',
      severity: 'warning',
      depth: 3,
      impact: 'Indirect circular dependency involving 3 packages',
      complexity: 5,
      rootCause: {
        originatingPackage: 'pkg-ui',
        problematicDependency: {
          from: 'pkg-ui',
          to: 'pkg-api',
          type: 'production',
          critical: false,
        },
        confidence: 82,
        explanation: "Package 'pkg-ui' is the root cause.",
        chain: [
          { from: 'pkg-ui', to: 'pkg-api', type: 'production', critical: false },
          { from: 'pkg-api', to: 'pkg-core', type: 'production', critical: false },
          { from: 'pkg-core', to: 'pkg-ui', type: 'production', critical: true },
        ],
        criticalEdge: {
          from: 'pkg-core',
          to: 'pkg-ui',
          type: 'production',
          critical: true,
        },
      },
      importTraces: [
        {
          fromPackage: 'pkg-ui',
          toPackage: 'pkg-api',
          filePath: 'src/hooks/useUser.ts',
          lineNumber: 3,
          statement: "import { fetchUser } from '@mono/api'",
          importType: 'esm-named',
          symbols: ['fetchUser'],
        },
      ],
      fixStrategies: [
        {
          type: 'extract-module',
          name: 'Extract Shared Module',
          description: 'Create a new shared package.',
          suitability: 8,
          effort: 'medium',
          pros: ['Clear separation', 'Isolates shared code'],
          cons: ['New package to maintain'],
          recommended: true,
          targetPackages: ['pkg-ui', 'pkg-api', 'pkg-core'],
          newPackageName: '@mono/shared',
        },
      ],
    }

    expect(circular.rootCause).toBeDefined()
    expect(circular.importTraces).toHaveLength(1)
    expect(circular.fixStrategies).toHaveLength(1)
    expect(circular.fixStrategies?.[0].recommended).toBe(true)
  })

  it('fixStrategies is optional for backward compatibility', () => {
    const circular: CircularDependencyInfo = {
      cycle: ['pkg-a', 'pkg-b', 'pkg-a'],
      type: 'direct',
      severity: 'critical',
      depth: 2,
      impact: 'Direct circular dependency',
      complexity: 3,
    }

    expect(circular.fixStrategies).toBeUndefined()
  })

  it('strategies are sorted by suitability (highest first)', () => {
    const circular: CircularDependencyInfo = {
      cycle: ['pkg-a', 'pkg-b', 'pkg-a'],
      type: 'direct',
      severity: 'warning',
      depth: 2,
      impact: 'Direct cycle',
      complexity: 4,
      fixStrategies: [
        {
          type: 'dependency-injection',
          name: 'DI',
          description: 'Invert dependency',
          suitability: 10,
          effort: 'low',
          pros: ['Best for direct cycles'],
          cons: ['Adds indirection'],
          recommended: true,
          targetPackages: ['pkg-a', 'pkg-b'],
        },
        {
          type: 'extract-module',
          name: 'Extract',
          description: 'Extract shared',
          suitability: 5,
          effort: 'medium',
          pros: ['Separation'],
          cons: ['New package'],
          recommended: false,
          targetPackages: ['pkg-a', 'pkg-b'],
        },
        {
          type: 'boundary-refactoring',
          name: 'Boundary',
          description: 'Restructure',
          suitability: 3,
          effort: 'high',
          pros: ['Architectural fix'],
          cons: ['High effort'],
          recommended: false,
          targetPackages: ['pkg-a', 'pkg-b'],
        },
      ],
    }

    // Verify sorted by suitability descending
    expect(circular.fixStrategies?.[0].suitability).toBe(10)
    expect(circular.fixStrategies?.[1].suitability).toBe(5)
    expect(circular.fixStrategies?.[2].suitability).toBe(3)

    // First should be recommended
    expect(circular.fixStrategies?.[0].recommended).toBe(true)
    expect(circular.fixStrategies?.[1].recommended).toBe(false)
  })
})

// Story 3.5: Refactoring Complexity Types
describe('RefactoringComplexity', () => {
  it('can represent complete complexity breakdown', () => {
    const complexity: RefactoringComplexity = {
      score: 5,
      estimatedTime: '30-60 minutes',
      breakdown: {
        filesAffected: {
          value: 3,
          weight: 0.25,
          contribution: 1.5,
          description: '3 source files need modification',
        },
        importsToChange: {
          value: 4,
          weight: 0.2,
          contribution: 1.2,
          description: '4 import statements need updating',
        },
        chainDepth: {
          value: 3,
          weight: 0.25,
          contribution: 1.5,
          description: 'Dependency chain has 3 levels',
        },
        packagesInvolved: {
          value: 3,
          weight: 0.15,
          contribution: 0.9,
          description: '3 packages involved in cycle',
        },
        externalDependencies: {
          value: 0,
          weight: 0.15,
          contribution: 0.15,
          description: 'No external dependencies in cycle',
        },
      },
      explanation: 'Moderate refactoring: 3 files, 4 imports, 3-level chain',
    }

    expect(complexity.score).toBe(5)
    expect(complexity.estimatedTime).toBe('30-60 minutes')
    expect(complexity.breakdown.filesAffected.value).toBe(3)
    expect(complexity.breakdown.chainDepth.weight).toBe(0.25)
  })

  it('weights sum to 1.0', () => {
    const breakdown: ComplexityBreakdown = {
      filesAffected: { value: 3, weight: 0.25, contribution: 1.5, description: '' },
      importsToChange: { value: 4, weight: 0.2, contribution: 1.2, description: '' },
      chainDepth: { value: 3, weight: 0.25, contribution: 1.5, description: '' },
      packagesInvolved: { value: 3, weight: 0.15, contribution: 0.9, description: '' },
      externalDependencies: { value: 0, weight: 0.15, contribution: 0.15, description: '' },
    }

    const totalWeight =
      breakdown.filesAffected.weight +
      breakdown.importsToChange.weight +
      breakdown.chainDepth.weight +
      breakdown.packagesInvolved.weight +
      breakdown.externalDependencies.weight

    expect(totalWeight).toBe(1.0)
  })
})

describe('ComplexityFactor', () => {
  it('can represent individual factor', () => {
    const factor: ComplexityFactor = {
      value: 5,
      weight: 0.25,
      contribution: 1.5,
      description: '5 source files need modification',
    }

    expect(factor.value).toBe(5)
    expect(factor.weight).toBe(0.25)
    expect(factor.contribution).toBe(1.5)
    expect(factor.description).toContain('files')
  })
})

describe('CircularDependencyInfo with RefactoringComplexity', () => {
  it('can include detailed refactoring complexity (Story 3.5)', () => {
    const circular: CircularDependencyInfo = {
      cycle: ['pkg-ui', 'pkg-api', 'pkg-core', 'pkg-ui'],
      type: 'indirect',
      severity: 'warning',
      depth: 3,
      impact: 'Indirect circular dependency involving 3 packages',
      complexity: 5, // Legacy field
      refactoringComplexity: {
        score: 5,
        estimatedTime: '30-60 minutes',
        breakdown: {
          filesAffected: { value: 3, weight: 0.25, contribution: 1.5, description: '3 files' },
          importsToChange: { value: 3, weight: 0.2, contribution: 1.2, description: '3 imports' },
          chainDepth: { value: 3, weight: 0.25, contribution: 1.5, description: '3 levels' },
          packagesInvolved: {
            value: 3,
            weight: 0.15,
            contribution: 0.9,
            description: '3 packages',
          },
          externalDependencies: {
            value: 0,
            weight: 0.15,
            contribution: 0.15,
            description: 'No external',
          },
        },
        explanation: 'Moderate refactoring',
      },
    }

    expect(circular.refactoringComplexity).toBeDefined()
    expect(circular.refactoringComplexity?.score).toBe(5)
    expect(circular.refactoringComplexity?.estimatedTime).toBe('30-60 minutes')
    expect(circular.complexity).toBe(5) // Legacy field still works
  })

  it('refactoringComplexity is optional for backward compatibility', () => {
    const circular: CircularDependencyInfo = {
      cycle: ['pkg-a', 'pkg-b', 'pkg-a'],
      type: 'direct',
      severity: 'warning',
      depth: 2,
      impact: 'Direct circular dependency',
      complexity: 3,
    }

    expect(circular.refactoringComplexity).toBeUndefined()
    expect(circular.complexity).toBe(3)
  })
})

describe('FixStrategy with Complexity', () => {
  it('can include detailed complexity per strategy (Story 3.5)', () => {
    const strategy: FixStrategy = {
      type: 'extract-module',
      name: 'Extract Shared Module',
      description: 'Create a new shared package.',
      suitability: 8,
      effort: 'medium',
      pros: ['Clean separation'],
      cons: ['New package'],
      recommended: true,
      targetPackages: ['pkg-ui', 'pkg-api'],
      newPackageName: '@mono/shared',
      complexity: {
        score: 5,
        estimatedTime: '30-60 minutes',
        breakdown: {
          filesAffected: { value: 3, weight: 0.25, contribution: 1.5, description: '3 files' },
          importsToChange: { value: 3, weight: 0.2, contribution: 1.2, description: '3 imports' },
          chainDepth: { value: 3, weight: 0.25, contribution: 1.5, description: '3 levels' },
          packagesInvolved: {
            value: 3,
            weight: 0.15,
            contribution: 0.9,
            description: '3 packages',
          },
          externalDependencies: {
            value: 0,
            weight: 0.15,
            contribution: 0.15,
            description: 'No external',
          },
        },
        explanation: 'Moderate refactoring',
      },
    }

    expect(strategy.complexity).toBeDefined()
    expect(strategy.complexity?.score).toBe(5)
    expect(strategy.complexity?.estimatedTime).toBe('30-60 minutes')
  })

  it('complexity is optional on FixStrategy', () => {
    const strategy: FixStrategy = {
      type: 'dependency-injection',
      name: 'DI',
      description: 'Invert dependency',
      suitability: 7,
      effort: 'low',
      pros: ['Minimal changes'],
      cons: ['Adds indirection'],
      recommended: false,
      targetPackages: ['pkg-a', 'pkg-b'],
    }

    expect(strategy.complexity).toBeUndefined()
  })
})
