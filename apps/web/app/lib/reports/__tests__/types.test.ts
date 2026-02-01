import type { AnalysisResult, ComprehensiveAnalysisResult } from '@monoguard/types'
import { Status } from '@monoguard/types'
import { describe, expect, it } from 'vitest'
import {
  buildReportData,
  buildReportDataFromComprehensive,
  DEFAULT_REPORT_SECTIONS,
  getHealthScoreRating,
  MONOGUARD_VERSION,
  RATING_THRESHOLDS,
} from '../types'

describe('Report Types', () => {
  describe('DEFAULT_REPORT_SECTIONS', () => {
    it('should have healthScore enabled by default', () => {
      expect(DEFAULT_REPORT_SECTIONS.healthScore).toBe(true)
    })

    it('should have circularDependencies enabled by default', () => {
      expect(DEFAULT_REPORT_SECTIONS.circularDependencies).toBe(true)
    })

    it('should have versionConflicts enabled by default', () => {
      expect(DEFAULT_REPORT_SECTIONS.versionConflicts).toBe(true)
    })

    it('should have fixRecommendations enabled by default', () => {
      expect(DEFAULT_REPORT_SECTIONS.fixRecommendations).toBe(true)
    })

    it('should have packageList disabled by default', () => {
      expect(DEFAULT_REPORT_SECTIONS.packageList).toBe(false)
    })

    it('should have graphSummary disabled by default', () => {
      expect(DEFAULT_REPORT_SECTIONS.graphSummary).toBe(false)
    })
  })

  describe('MONOGUARD_VERSION', () => {
    it('should be a valid semver string', () => {
      expect(MONOGUARD_VERSION).toMatch(/^\d+\.\d+\.\d+$/)
    })
  })

  describe('RATING_THRESHOLDS', () => {
    it('should define all rating levels', () => {
      expect(RATING_THRESHOLDS).toHaveProperty('excellent')
      expect(RATING_THRESHOLDS).toHaveProperty('good')
      expect(RATING_THRESHOLDS).toHaveProperty('fair')
      expect(RATING_THRESHOLDS).toHaveProperty('poor')
      expect(RATING_THRESHOLDS).toHaveProperty('critical')
    })

    it('should have thresholds in descending order', () => {
      expect(RATING_THRESHOLDS.excellent).toBeGreaterThan(RATING_THRESHOLDS.good)
      expect(RATING_THRESHOLDS.good).toBeGreaterThan(RATING_THRESHOLDS.fair)
      expect(RATING_THRESHOLDS.fair).toBeGreaterThan(RATING_THRESHOLDS.poor)
      expect(RATING_THRESHOLDS.poor).toBeGreaterThan(RATING_THRESHOLDS.critical)
    })
  })

  describe('getHealthScoreRating', () => {
    it('should return excellent for scores >= 85', () => {
      expect(getHealthScoreRating(85)).toBe('excellent')
      expect(getHealthScoreRating(100)).toBe('excellent')
      expect(getHealthScoreRating(90)).toBe('excellent')
    })

    it('should return good for scores >= 70 and < 85', () => {
      expect(getHealthScoreRating(70)).toBe('good')
      expect(getHealthScoreRating(84)).toBe('good')
    })

    it('should return fair for scores >= 50 and < 70', () => {
      expect(getHealthScoreRating(50)).toBe('fair')
      expect(getHealthScoreRating(69)).toBe('fair')
    })

    it('should return poor for scores >= 30 and < 50', () => {
      expect(getHealthScoreRating(30)).toBe('poor')
      expect(getHealthScoreRating(49)).toBe('poor')
    })

    it('should return critical for scores < 30', () => {
      expect(getHealthScoreRating(29)).toBe('critical')
      expect(getHealthScoreRating(0)).toBe('critical')
    })
  })

  describe('buildReportData', () => {
    const mockAnalysisResult: AnalysisResult = {
      healthScore: 75,
      packages: 50,
      circularDependencies: [
        {
          cycle: ['pkg-a', 'pkg-b', 'pkg-a'],
          type: 'direct',
          severity: 'warning',
          depth: 2,
          impact: 'High coupling between packages',
          complexity: 5,
          priorityScore: 8,
        },
        {
          cycle: ['pkg-c', 'pkg-d', 'pkg-e', 'pkg-c'],
          type: 'indirect',
          severity: 'critical',
          depth: 3,
          impact: 'Deep circular dependency',
          complexity: 7,
          priorityScore: 9,
          fixStrategies: [
            {
              type: 'extract-module',
              name: 'Extract Shared Module',
              description: 'Extract shared code into a new package',
              suitability: 8,
              effort: 'low',
              pros: ['Clean separation'],
              cons: ['New package to maintain'],
              recommended: true,
              targetPackages: ['pkg-c', 'pkg-d'],
            },
          ],
        },
      ],
      versionConflicts: [
        {
          packageName: 'lodash',
          conflictingVersions: [
            {
              version: '4.17.21',
              packages: ['pkg-a'],
              isBreaking: false,
              depType: 'production',
            },
            {
              version: '4.17.15',
              packages: ['pkg-b'],
              isBreaking: false,
              depType: 'production',
            },
          ],
          severity: 'warning',
          resolution: 'Upgrade all to 4.17.21',
          impact: 'Minor version mismatch',
        },
      ],
      healthScoreDetails: {
        overall: 75,
        rating: 'good',
        breakdown: {
          circularScore: 60,
          conflictScore: 80,
          depthScore: 85,
          couplingScore: 70,
        },
        factors: [
          {
            name: 'Circular Dependencies',
            score: 60,
            weight: 0.4,
            weightedScore: 24,
            description: 'Score from circular deps',
            recommendations: ['Fix cycles'],
          },
          {
            name: 'Version Conflicts',
            score: 80,
            weight: 0.3,
            weightedScore: 24,
            description: 'Score from conflicts',
            recommendations: ['Align versions'],
          },
        ],
        updatedAt: '2026-01-25T10:00:00Z',
      },
      metadata: {
        version: '0.1.0',
        durationMs: 1234,
        filesProcessed: 100,
        workspaceType: 'pnpm',
      },
      graph: {
        nodes: {
          'pkg-a': {
            name: 'pkg-a',
            version: '1.0.0',
            path: '/packages/pkg-a',
            dependencies: ['pkg-b'],
            devDependencies: [],
            peerDependencies: [],
          },
          'pkg-b': {
            name: 'pkg-b',
            version: '1.0.0',
            path: '/packages/pkg-b',
            dependencies: [],
            devDependencies: [],
            peerDependencies: [],
          },
        },
        edges: [
          {
            from: 'pkg-a',
            to: 'pkg-b',
            type: 'production',
            versionRange: '^1.0.0',
          },
        ],
        rootPath: '/workspace',
        workspaceType: 'pnpm',
      },
    }

    it('should build report data with correct metadata', () => {
      const reportData = buildReportData(mockAnalysisResult, 'test-project')

      expect(reportData.metadata.projectName).toBe('test-project')
      expect(reportData.metadata.monoguardVersion).toBe(MONOGUARD_VERSION)
      expect(reportData.metadata.packageCount).toBe(50)
      expect(reportData.metadata.analysisDuration).toBe(1234)
      expect(reportData.metadata.nodeCount).toBe(2)
      expect(reportData.metadata.edgeCount).toBe(1)
      expect(reportData.metadata.generatedAt).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/)
    })

    it('should build health score report from details', () => {
      const reportData = buildReportData(mockAnalysisResult, 'test-project')

      expect(reportData.healthScore.overall).toBe(75)
      expect(reportData.healthScore.rating).toBe('good')
      expect(reportData.healthScore.breakdown).toHaveLength(2)
      expect(reportData.healthScore.breakdown[0].category).toBe('Circular Dependencies')
    })

    it('should build health score fallback when no details', () => {
      const minimal: AnalysisResult = {
        healthScore: 45,
        packages: 10,
      }
      const reportData = buildReportData(minimal, 'test')

      expect(reportData.healthScore.overall).toBe(45)
      expect(reportData.healthScore.rating).toBe('poor')
      expect(reportData.healthScore.breakdown).toHaveLength(1)
    })

    it('should build circular dependency report', () => {
      const reportData = buildReportData(mockAnalysisResult, 'test-project')

      expect(reportData.circularDependencies.totalCount).toBe(2)
      expect(reportData.circularDependencies.bySeverity.critical).toBe(1)
      expect(reportData.circularDependencies.bySeverity.high).toBe(1)
      expect(reportData.circularDependencies.cycles).toHaveLength(2)
      expect(reportData.circularDependencies.cycles[0].type).toBe('direct')
    })

    it('should build version conflict report', () => {
      const reportData = buildReportData(mockAnalysisResult, 'test-project')

      expect(reportData.versionConflicts.totalCount).toBe(1)
      expect(reportData.versionConflicts.conflicts[0].packageName).toBe('lodash')
      expect(reportData.versionConflicts.conflicts[0].versions).toEqual(['4.17.21', '4.17.15'])
    })

    it('should build fix recommendation report', () => {
      const reportData = buildReportData(mockAnalysisResult, 'test-project')

      expect(reportData.fixRecommendations.totalCount).toBe(1)
      expect(reportData.fixRecommendations.quickWins).toBe(1)
      expect(reportData.fixRecommendations.recommendations[0].title).toBe('Extract Shared Module')
      expect(reportData.fixRecommendations.recommendations[0].effort).toBe('low')
    })

    it('should handle empty analysis result', () => {
      const empty: AnalysisResult = {
        healthScore: 100,
        packages: 0,
      }
      const reportData = buildReportData(empty, 'empty-project')

      expect(reportData.circularDependencies.totalCount).toBe(0)
      expect(reportData.versionConflicts.totalCount).toBe(0)
      expect(reportData.fixRecommendations.totalCount).toBe(0)
      expect(reportData.metadata.nodeCount).toBe(0)
      expect(reportData.metadata.edgeCount).toBe(0)
    })

    it('should sort fix recommendations by priority descending', () => {
      const withMultipleStrategies: AnalysisResult = {
        healthScore: 50,
        packages: 5,
        circularDependencies: [
          {
            cycle: ['a', 'b', 'a'],
            type: 'direct',
            severity: 'warning',
            depth: 2,
            impact: 'test',
            complexity: 3,
            priorityScore: 5,
            fixStrategies: [
              {
                type: 'extract-module',
                name: 'Low Priority',
                description: 'Low priority fix',
                suitability: 3,
                effort: 'high',
                pros: [],
                cons: [],
                recommended: false,
                targetPackages: ['a'],
              },
              {
                type: 'dependency-injection',
                name: 'High Priority',
                description: 'High priority fix',
                suitability: 9,
                effort: 'low',
                pros: [],
                cons: [],
                recommended: true,
                targetPackages: ['b'],
              },
            ],
          },
        ],
      }

      const reportData = buildReportData(withMultipleStrategies, 'test')
      expect(reportData.fixRecommendations.recommendations[0].title).toBe('High Priority')
      expect(reportData.fixRecommendations.recommendations[1].title).toBe('Low Priority')
    })

    it('should classify info severity as medium in circular deps', () => {
      const withInfo: AnalysisResult = {
        healthScore: 50,
        packages: 5,
        circularDependencies: [
          {
            cycle: ['a', 'b', 'a'],
            type: 'direct',
            severity: 'info',
            depth: 2,
            impact: 'test',
            complexity: 3,
            priorityScore: 3,
          },
        ],
      }
      const reportData = buildReportData(withInfo, 'test')
      expect(reportData.circularDependencies.bySeverity.medium).toBe(1)
    })

    it('should classify unknown severity as low in circular deps', () => {
      const withUnknown: AnalysisResult = {
        healthScore: 50,
        packages: 5,
        circularDependencies: [
          {
            cycle: ['a', 'b', 'a'],
            type: 'direct',
            severity: 'unknown' as never,
            depth: 2,
            impact: 'test',
            complexity: 3,
            priorityScore: 1,
          },
        ],
      }
      const reportData = buildReportData(withUnknown, 'test')
      expect(reportData.circularDependencies.bySeverity.low).toBe(1)
    })

    it('should classify info severity as medium in version conflicts', () => {
      const withInfo: AnalysisResult = {
        healthScore: 50,
        packages: 5,
        versionConflicts: [
          {
            packageName: 'test',
            conflictingVersions: [
              { version: '1.0.0', packages: ['a'], isBreaking: false, depType: 'production' },
              { version: '2.0.0', packages: ['b'], isBreaking: false, depType: 'production' },
            ],
            severity: 'info',
            resolution: 'Upgrade',
            impact: 'Minor',
          },
        ],
      }
      const reportData = buildReportData(withInfo, 'test')
      expect(reportData.versionConflicts.byRiskLevel.medium).toBe(1)
    })

    it('should classify unknown severity as low in version conflicts', () => {
      const withUnknown: AnalysisResult = {
        healthScore: 50,
        packages: 5,
        versionConflicts: [
          {
            packageName: 'test',
            conflictingVersions: [
              { version: '1.0.0', packages: ['a'], isBreaking: false, depType: 'production' },
              { version: '2.0.0', packages: ['b'], isBreaking: false, depType: 'production' },
            ],
            severity: 'unknown' as never,
            resolution: 'Upgrade',
            impact: 'Minor',
          },
        ],
      }
      const reportData = buildReportData(withUnknown, 'test')
      expect(reportData.versionConflicts.byRiskLevel.low).toBe(1)
    })

    it('should handle metadata without durationMs', () => {
      const noMeta: AnalysisResult = {
        healthScore: 50,
        packages: 5,
      }
      const reportData = buildReportData(noMeta, 'test')
      expect(reportData.metadata.analysisDuration).toBe(0)
    })

    it('should skip fix strategies for deps without strategies', () => {
      const noStrategies: AnalysisResult = {
        healthScore: 50,
        packages: 5,
        circularDependencies: [
          {
            cycle: ['a', 'b', 'a'],
            type: 'direct',
            severity: 'warning',
            depth: 2,
            impact: 'test',
            complexity: 3,
            priorityScore: 5,
          },
        ],
      }
      const reportData = buildReportData(noStrategies, 'test')
      expect(reportData.fixRecommendations.totalCount).toBe(0)
      expect(reportData.fixRecommendations.quickWins).toBe(0)
    })

    it('should classify medium impact for mid-range suitability', () => {
      const midSuitability: AnalysisResult = {
        healthScore: 50,
        packages: 5,
        circularDependencies: [
          {
            cycle: ['a', 'b', 'a'],
            type: 'direct',
            severity: 'warning',
            depth: 2,
            impact: 'test',
            complexity: 3,
            priorityScore: 5,
            fixStrategies: [
              {
                type: 'extract-module',
                name: 'Medium Impact',
                description: 'Medium suitability',
                suitability: 5,
                effort: 'medium',
                pros: [],
                cons: [],
                recommended: false,
                targetPackages: ['a'],
              },
            ],
          },
        ],
      }
      const reportData = buildReportData(midSuitability, 'test')
      expect(reportData.fixRecommendations.recommendations[0].impact).toBe('medium')
    })

    it('should classify low impact for low suitability', () => {
      const lowSuitability: AnalysisResult = {
        healthScore: 50,
        packages: 5,
        circularDependencies: [
          {
            cycle: ['a', 'b', 'a'],
            type: 'direct',
            severity: 'warning',
            depth: 2,
            impact: 'test',
            complexity: 3,
            priorityScore: 5,
            fixStrategies: [
              {
                type: 'extract-module',
                name: 'Low Impact',
                description: 'Low suitability',
                suitability: 2,
                effort: 'high',
                pros: [],
                cons: [],
                recommended: false,
                targetPackages: ['a'],
              },
            ],
          },
        ],
      }
      const reportData = buildReportData(lowSuitability, 'test')
      expect(reportData.fixRecommendations.recommendations[0].impact).toBe('low')
    })
  })

  describe('buildReportDataFromComprehensive', () => {
    it('should build report data from comprehensive result with HealthScore object', () => {
      const analysis: ComprehensiveAnalysisResult = {
        id: 'test-1',
        uploadId: 'upload-1',
        status: Status.COMPLETED,
        startedAt: new Date().toISOString(),
        progress: 100,
        results: {
          summary: {
            totalPackages: 10,
            duplicateCount: 1,
            conflictCount: 1,
            unusedCount: 0,
            circularCount: 1,
            healthScore: 75,
          },
          healthScore: {
            overall: 80,
            dependencies: 85,
            architecture: 75,
            maintainability: 80,
            security: 90,
            performance: 85,
            lastUpdated: new Date().toISOString(),
            trend: 'stable',
            factors: [
              {
                name: 'Dependencies',
                score: 85,
                weight: 0.4,
                weightedScore: 34,
                description: 'Good dep health',
                recommendations: [],
              },
              {
                name: 'Architecture',
                score: 75,
                weight: 0.3,
                weightedScore: 22.5,
                description: 'Architecture ok',
                recommendations: [],
              },
            ],
          },
          circularDependencies: [
            {
              cycle: ['pkg-a', 'pkg-b', 'pkg-a'],
              type: 'direct',
              severity: 'critical',
              impact: 'High',
            },
          ],
          versionConflicts: [
            {
              packageName: 'lodash',
              conflictingVersions: [
                { version: '4.17.21', packages: ['a'], isBreaking: false },
                { version: '3.10.1', packages: ['b'], isBreaking: true },
              ],
              riskLevel: 'high',
              resolution: 'Upgrade to 4.17.21',
              impact: 'Breaking change',
            },
          ],
        },
      }

      const reportData = buildReportDataFromComprehensive(analysis, 'test-project')

      expect(reportData.metadata.projectName).toBe('test-project')
      expect(reportData.metadata.monoguardVersion).toBe(MONOGUARD_VERSION)
      expect(reportData.metadata.packageCount).toBe(10)
      expect(reportData.healthScore.overall).toBe(80)
      expect(reportData.healthScore.breakdown).toHaveLength(2)
      expect(reportData.healthScore.breakdown[0].category).toBe('Dependencies')
      expect(reportData.healthScore.breakdown[0].weight).toBe(40)
      expect(reportData.circularDependencies.totalCount).toBe(1)
      expect(reportData.circularDependencies.bySeverity.critical).toBe(1)
      expect(reportData.versionConflicts.totalCount).toBe(1)
      expect(reportData.versionConflicts.byRiskLevel.high).toBe(1)
      expect(reportData.fixRecommendations.totalCount).toBe(0)
    })

    it('should handle numeric healthScore (no HealthScore object)', () => {
      const analysis: ComprehensiveAnalysisResult = {
        id: 'test-2',
        uploadId: 'upload-2',
        status: Status.COMPLETED,
        startedAt: new Date().toISOString(),
        progress: 100,
        results: {
          summary: {
            totalPackages: 5,
            duplicateCount: 0,
            conflictCount: 0,
            unusedCount: 0,
            circularCount: 0,
            healthScore: 60,
          },
          healthScore: 60 as never,
        },
      }

      const reportData = buildReportDataFromComprehensive(analysis, 'test')

      expect(reportData.healthScore.overall).toBe(60)
      expect(reportData.healthScore.rating).toBe('fair')
      expect(reportData.healthScore.breakdown).toHaveLength(1)
      expect(reportData.healthScore.breakdown[0].category).toBe('Overall')
    })

    it('should fall back to summary healthScore when healthScore object is undefined', () => {
      const analysis: ComprehensiveAnalysisResult = {
        id: 'test-3',
        uploadId: 'upload-3',
        status: Status.COMPLETED,
        startedAt: new Date().toISOString(),
        progress: 100,
        results: {
          summary: {
            totalPackages: 3,
            duplicateCount: 0,
            conflictCount: 0,
            unusedCount: 0,
            circularCount: 0,
            healthScore: 45,
          },
        },
      }

      const reportData = buildReportDataFromComprehensive(analysis, 'test')

      expect(reportData.healthScore.overall).toBe(45)
      expect(reportData.healthScore.rating).toBe('poor')
    })

    it('should handle empty results', () => {
      const analysis: ComprehensiveAnalysisResult = {
        id: 'test-4',
        uploadId: 'upload-4',
        status: Status.COMPLETED,
        startedAt: new Date().toISOString(),
        progress: 100,
        results: {},
      }

      const reportData = buildReportDataFromComprehensive(analysis, 'empty-test')

      expect(reportData.metadata.packageCount).toBe(0)
      expect(reportData.healthScore.overall).toBe(0)
      expect(reportData.healthScore.rating).toBe('critical')
      expect(reportData.circularDependencies.totalCount).toBe(0)
      expect(reportData.versionConflicts.totalCount).toBe(0)
    })

    it('should handle undefined results', () => {
      const analysis: ComprehensiveAnalysisResult = {
        id: 'test-5',
        uploadId: 'upload-5',
        status: Status.IN_PROGRESS,
        startedAt: new Date().toISOString(),
        progress: 50,
      }

      const reportData = buildReportDataFromComprehensive(analysis, 'in-progress')

      expect(reportData.metadata.packageCount).toBe(0)
      expect(reportData.healthScore.overall).toBe(0)
      expect(reportData.circularDependencies.totalCount).toBe(0)
      expect(reportData.versionConflicts.totalCount).toBe(0)
    })

    it('should classify severity levels correctly for circular deps', () => {
      const analysis: ComprehensiveAnalysisResult = {
        id: 'test-6',
        uploadId: 'upload-6',
        status: Status.COMPLETED,
        startedAt: new Date().toISOString(),
        progress: 100,
        results: {
          circularDependencies: [
            { cycle: ['a', 'b', 'a'], type: 'direct', severity: 'critical', impact: '' },
            { cycle: ['c', 'd', 'c'], type: 'direct', severity: 'high', impact: '' },
            { cycle: ['e', 'f', 'e'], type: 'direct', severity: 'medium', impact: '' },
            { cycle: ['g', 'h', 'g'], type: 'direct', severity: 'low', impact: '' },
          ],
        },
      }

      const reportData = buildReportDataFromComprehensive(analysis, 'test')

      expect(reportData.circularDependencies.bySeverity.critical).toBe(1)
      expect(reportData.circularDependencies.bySeverity.high).toBe(1)
      expect(reportData.circularDependencies.bySeverity.medium).toBe(1)
      expect(reportData.circularDependencies.bySeverity.low).toBe(1)
    })

    it('should classify risk levels correctly for version conflicts', () => {
      const analysis: ComprehensiveAnalysisResult = {
        id: 'test-7',
        uploadId: 'upload-7',
        status: Status.COMPLETED,
        startedAt: new Date().toISOString(),
        progress: 100,
        results: {
          versionConflicts: [
            {
              packageName: 'a',
              conflictingVersions: [{ version: '1.0.0', packages: ['x'], isBreaking: false }],
              riskLevel: 'critical',
              resolution: '',
              impact: '',
            },
            {
              packageName: 'b',
              conflictingVersions: [{ version: '1.0.0', packages: ['x'], isBreaking: false }],
              riskLevel: 'high',
              resolution: '',
              impact: '',
            },
            {
              packageName: 'c',
              conflictingVersions: [{ version: '1.0.0', packages: ['x'], isBreaking: false }],
              riskLevel: 'medium',
              resolution: '',
              impact: '',
            },
            {
              packageName: 'd',
              conflictingVersions: [{ version: '1.0.0', packages: ['x'], isBreaking: false }],
              riskLevel: 'low',
              resolution: '',
              impact: '',
            },
          ],
        },
      }

      const reportData = buildReportDataFromComprehensive(analysis, 'test')

      expect(reportData.versionConflicts.byRiskLevel.critical).toBe(1)
      expect(reportData.versionConflicts.byRiskLevel.high).toBe(1)
      expect(reportData.versionConflicts.byRiskLevel.medium).toBe(1)
      expect(reportData.versionConflicts.byRiskLevel.low).toBe(1)
    })
  })
})
