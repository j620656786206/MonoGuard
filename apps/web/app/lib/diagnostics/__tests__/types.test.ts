import { describe, expect, it } from 'vitest'
import type {
  CycleEdge,
  CycleNode,
  DiagnosticMetadata,
  DiagnosticReport,
  ExecutiveSummary,
  FixStrategyGuide,
  ImpactAssessmentDetails,
  RelatedCycleInfo,
  RippleNode,
  RootCauseDetails,
} from '../types'
import { MONOGUARD_VERSION } from '../types'

describe('Diagnostic Report Types', () => {
  describe('MONOGUARD_VERSION', () => {
    it('should be a valid semver string', () => {
      expect(MONOGUARD_VERSION).toMatch(/^\d+\.\d+\.\d+$/)
    })
  })

  describe('DiagnosticReport interface', () => {
    it('should have all required sections', () => {
      const report: DiagnosticReport = {
        id: 'diag-test-123',
        cycleId: 'cycle-1',
        generatedAt: '2026-01-25T10:00:00Z',
        monoguardVersion: '0.1.0',
        projectName: 'test-project',
        executiveSummary: {
          description: 'Test cycle',
          severity: 'medium',
          recommendation: 'Fix it',
          estimatedEffort: 'medium',
          affectedPackagesCount: 3,
          cycleLength: 3,
        },
        cyclePath: {
          nodes: [],
          edges: [],
          breakingPoint: { fromPackage: 'a', toPackage: 'b', reason: 'test' },
          svgDiagram: '<svg></svg>',
          asciiDiagram: 'a -> b -> c -> a',
        },
        rootCause: {
          explanation: 'Root cause explanation',
          confidenceScore: 85,
          originatingPackage: 'pkg-a',
          originatingReason: 'Creates circular import',
          alternativeCandidates: [],
          codeReferences: [],
        },
        fixStrategies: [],
        impactAssessment: {
          directParticipants: ['pkg-a', 'pkg-b'],
          directParticipantsCount: 2,
          indirectDependents: ['pkg-c'],
          indirectDependentsCount: 1,
          totalAffectedCount: 3,
          percentageOfMonorepo: 15,
          riskLevel: 'medium',
          riskExplanation: 'Moderate impact',
          rippleEffectTree: { package: 'Cycle', depth: 0, dependents: [] },
        },
        relatedCycles: [],
        metadata: {
          generatedAt: '2026-01-25T10:00:00Z',
          generationDurationMs: 150,
          monoguardVersion: '0.1.0',
          projectName: 'test-project',
          analysisConfigHash: 'default',
        },
      }

      expect(report.id).toBe('diag-test-123')
      expect(report.cycleId).toBe('cycle-1')
      expect(report.executiveSummary.severity).toBe('medium')
      expect(report.cyclePath.breakingPoint.fromPackage).toBe('a')
      expect(report.rootCause.confidenceScore).toBe(85)
      expect(report.impactAssessment.totalAffectedCount).toBe(3)
      expect(report.metadata.generationDurationMs).toBe(150)
    })
  })

  describe('ExecutiveSummary interface', () => {
    it('should accept all severity levels', () => {
      const severities: ExecutiveSummary['severity'][] = ['critical', 'high', 'medium', 'low']
      for (const sev of severities) {
        const summary: ExecutiveSummary = {
          description: `${sev} cycle`,
          severity: sev,
          recommendation: 'Fix it',
          estimatedEffort: 'medium',
          affectedPackagesCount: 2,
          cycleLength: 2,
        }
        expect(summary.severity).toBe(sev)
      }
    })

    it('should accept all effort levels', () => {
      const efforts: ExecutiveSummary['estimatedEffort'][] = ['low', 'medium', 'high']
      for (const effort of efforts) {
        const summary: ExecutiveSummary = {
          description: 'test',
          severity: 'low',
          recommendation: 'Fix',
          estimatedEffort: effort,
          affectedPackagesCount: 2,
          cycleLength: 2,
        }
        expect(summary.estimatedEffort).toBe(effort)
      }
    })
  })

  describe('CycleNode interface', () => {
    it('should have position coordinates', () => {
      const node: CycleNode = {
        id: 'pkg-a',
        name: 'pkg-a',
        path: 'packages/pkg-a',
        isInCycle: true,
        position: { x: 100, y: 200 },
      }
      expect(node.position.x).toBe(100)
      expect(node.position.y).toBe(200)
    })
  })

  describe('CycleEdge interface', () => {
    it('should mark breaking point edges', () => {
      const edge: CycleEdge = {
        from: 'pkg-a',
        to: 'pkg-b',
        isBreakingPoint: true,
        importStatement: "import { foo } from 'pkg-b'",
        filePath: 'src/index.ts',
        lineNumber: 1,
      }
      expect(edge.isBreakingPoint).toBe(true)
      expect(edge.importStatement).toBeTruthy()
    })
  })

  describe('RootCauseDetails interface', () => {
    it('should include alternative candidates when confidence is low', () => {
      const details: RootCauseDetails = {
        explanation: 'Possible root cause',
        confidenceScore: 60,
        originatingPackage: 'pkg-a',
        originatingReason: 'Likely source',
        alternativeCandidates: [{ package: 'pkg-b', reason: 'Also possible', confidence: 40 }],
        codeReferences: [{ file: 'src/index.ts', line: 5, importStatement: "import from 'pkg-b'" }],
      }
      expect(details.confidenceScore).toBe(60)
      expect(details.alternativeCandidates).toHaveLength(1)
      expect(details.codeReferences).toHaveLength(1)
    })
  })

  describe('FixStrategyGuide interface', () => {
    it('should include all strategy types', () => {
      const strategies: FixStrategyGuide['strategy'][] = [
        'extract-module',
        'dependency-injection',
        'boundary-refactoring',
      ]
      for (const strategy of strategies) {
        const guide: FixStrategyGuide = {
          strategy,
          title: `${strategy} Guide`,
          description: 'Test strategy',
          suitabilityScore: 8,
          estimatedEffort: 'medium',
          estimatedTime: '30-60 minutes',
          pros: ['Pro 1'],
          cons: ['Con 1'],
          steps: [],
          codeSnippets: {
            before: '// before',
            after: '// after',
          },
        }
        expect(guide.strategy).toBe(strategy)
      }
    })
  })

  describe('ImpactAssessmentDetails interface', () => {
    it('should calculate percentage correctly', () => {
      const impact: ImpactAssessmentDetails = {
        directParticipants: ['a', 'b', 'c'],
        directParticipantsCount: 3,
        indirectDependents: ['d', 'e'],
        indirectDependentsCount: 2,
        totalAffectedCount: 5,
        percentageOfMonorepo: 25,
        riskLevel: 'high',
        riskExplanation: 'High impact',
        rippleEffectTree: { package: 'Cycle', depth: 0, dependents: [] },
      }
      expect(impact.totalAffectedCount).toBe(
        impact.directParticipantsCount + impact.indirectDependentsCount
      )
      expect(impact.percentageOfMonorepo).toBe(25)
    })
  })

  describe('RelatedCycleInfo interface', () => {
    it('should track shared packages between cycles', () => {
      const related: RelatedCycleInfo = {
        cycleId: 'cycle-2',
        sharedPackages: ['pkg-b'],
        overlapPercentage: 33,
        recommendFixTogether: true,
        reason: 'Share common package pkg-b',
      }
      expect(related.sharedPackages).toContain('pkg-b')
      expect(related.recommendFixTogether).toBe(true)
    })
  })

  describe('RippleNode interface', () => {
    it('should support nested tree structure', () => {
      const tree: RippleNode = {
        package: 'Cycle',
        depth: 0,
        dependents: [
          {
            package: 'pkg-a',
            depth: 1,
            dependents: [{ package: 'pkg-d', depth: 2, dependents: [] }],
          },
        ],
      }
      expect(tree.dependents).toHaveLength(1)
      expect(tree.dependents[0].dependents).toHaveLength(1)
      expect(tree.dependents[0].dependents[0].package).toBe('pkg-d')
    })
  })

  describe('DiagnosticMetadata interface', () => {
    it('should include ISO 8601 timestamp', () => {
      const metadata: DiagnosticMetadata = {
        generatedAt: '2026-01-25T10:00:00Z',
        generationDurationMs: 250,
        monoguardVersion: '0.1.0',
        projectName: 'test-project',
        analysisConfigHash: 'abc123',
      }
      expect(metadata.generatedAt).toMatch(/^\d{4}-\d{2}-\d{2}T/)
      expect(metadata.generationDurationMs).toBeGreaterThanOrEqual(0)
    })
  })
})
