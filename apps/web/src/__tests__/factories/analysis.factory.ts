/**
 * Test Factories for Analysis-related types
 *
 * Pattern: Pure function with Partial<T> overrides (TEA knowledge base: data-factories.md)
 */

import type {
  AnalysisSummary,
  ArchitectureValidationResults,
  ArchitectureViolation,
  BundleImpactReport,
  CircularDependency,
  ComprehensiveAnalysisResult,
  DependencyAnalysisResults,
  DuplicateDetectionResults,
  DuplicateGroup,
  HealthScore,
  VersionConflict,
} from '@monoguard/types'
import { Status } from '@monoguard/types'

/**
 * Creates a mock AnalysisSummary
 */
export function createAnalysisSummary(overrides: Partial<AnalysisSummary> = {}): AnalysisSummary {
  return {
    totalPackages: 10,
    duplicateCount: 2,
    conflictCount: 1,
    unusedCount: 3,
    circularCount: 1,
    healthScore: 75,
    ...overrides,
  }
}

/**
 * Creates a mock CircularDependency
 */
export function createCircularDependency(
  overrides: Partial<CircularDependency> = {}
): CircularDependency {
  return {
    cycle: ['package-a', 'package-b', 'package-c'],
    type: 'direct',
    severity: 'medium',
    impact: 'May cause build issues and increased bundle size',
    ...overrides,
  }
}

/**
 * Creates a mock VersionConflict
 */
export function createVersionConflict(overrides: Partial<VersionConflict> = {}): VersionConflict {
  return {
    packageName: 'lodash',
    conflictingVersions: [
      { version: '4.17.21', packages: ['package-a', 'package-b'], isBreaking: false },
      { version: '3.10.1', packages: ['package-c'], isBreaking: true },
    ],
    riskLevel: 'medium',
    resolution: 'Upgrade all packages to use lodash 4.17.21',
    impact: 'Different lodash versions may cause inconsistent behavior',
    ...overrides,
  }
}

/**
 * Creates a mock ArchitectureViolation
 */
export function createArchitectureViolation(
  overrides: Partial<ArchitectureViolation> = {}
): ArchitectureViolation {
  return {
    ruleName: 'no-direct-ui-import',
    severity: 'high',
    description: 'UI components should not be imported directly from domain layer',
    violatingFile: 'packages/domain/src/user.ts',
    violatingImport: '@/ui/components/Button',
    expectedLayer: 'domain',
    actualLayer: 'ui',
    suggestion: 'Use dependency injection or move the import to the application layer',
    ...overrides,
  }
}

/**
 * Creates a mock BundleImpactReport
 */
export function createBundleImpactReport(
  overrides: Partial<BundleImpactReport> = {}
): BundleImpactReport {
  return {
    totalSize: '2.5 MB',
    duplicateSize: '500 KB',
    unusedSize: '200 KB',
    potentialSavings: '700 KB',
    breakdown: [
      { packageName: 'lodash', size: '300 KB', percentage: 12, duplicates: 2 },
      { packageName: 'moment', size: '250 KB', percentage: 10, duplicates: 0 },
    ],
    ...overrides,
  }
}

/**
 * Creates a mock DuplicateGroup
 */
export function createDuplicateGroup(overrides: Partial<DuplicateGroup> = {}): DuplicateGroup {
  return {
    packageName: 'lodash',
    versions: [
      {
        version: '4.17.21',
        size: '100 KB',
        usageCount: 5,
        packages: ['app-a', 'app-b'],
        isRecommended: true,
      },
      {
        version: '4.17.15',
        size: '100 KB',
        usageCount: 2,
        packages: ['lib-c'],
        isRecommended: false,
      },
    ],
    totalSize: '200 KB',
    wastedSize: '100 KB',
    riskLevel: 'medium',
    affectedPackages: ['app-a', 'app-b', 'lib-c'],
    ...overrides,
  }
}

/**
 * Creates a mock DuplicateDetectionResults
 */
export function createDuplicateDetectionResults(
  overrides: Partial<DuplicateDetectionResults> = {}
): DuplicateDetectionResults {
  return {
    totalDuplicates: 3,
    potentialSavings: '500 KB',
    duplicateGroups: [createDuplicateGroup()],
    recommendations: [
      {
        type: 'consolidate',
        packageName: 'lodash',
        description: 'Consolidate lodash versions',
        estimatedSavings: '100 KB',
        difficulty: 'easy',
        steps: ['Update package.json', 'Run npm dedupe'],
      },
    ],
    ...overrides,
  }
}

/**
 * Creates a mock HealthScore
 */
export function createHealthScore(overrides: Partial<HealthScore> = {}): HealthScore {
  return {
    overall: 75,
    dependencies: 80,
    architecture: 70,
    maintainability: 75,
    security: 90,
    performance: 85,
    lastUpdated: new Date().toISOString(),
    trend: 'stable',
    factors: [
      {
        name: 'Dependencies',
        score: 80,
        weight: 0.3,
        description: '2 duplicate dependencies found',
        recommendations: ['Remove duplicate dependencies'],
      },
      {
        name: 'Circular Dependencies',
        score: 70,
        weight: 0.2,
        description: '1 circular dependency detected',
        recommendations: ['Refactor circular dependencies'],
      },
    ],
    ...overrides,
  }
}

/**
 * Creates a mock DependencyAnalysisResults
 */
export function createDependencyAnalysisResults(
  overrides: Partial<DependencyAnalysisResults> = {}
): DependencyAnalysisResults {
  return {
    duplicateDependencies: [],
    versionConflicts: [createVersionConflict()],
    unusedDependencies: [],
    circularDependencies: [createCircularDependency()],
    bundleImpact: createBundleImpactReport(),
    summary: createAnalysisSummary(),
    ...overrides,
  }
}

/**
 * Creates a mock ArchitectureValidationResults
 */
export function createArchitectureValidationResults(
  overrides: Partial<ArchitectureValidationResults> = {}
): ArchitectureValidationResults {
  return {
    violations: [createArchitectureViolation()],
    layerCompliance: [
      {
        layerName: 'domain',
        totalFiles: 10,
        compliantFiles: 9,
        violationCount: 1,
        compliancePercentage: 90,
      },
      {
        layerName: 'ui',
        totalFiles: 20,
        compliantFiles: 20,
        violationCount: 0,
        compliancePercentage: 100,
      },
    ],
    circularDependencies: [],
    summary: {
      totalViolations: 1,
      criticalViolations: 0,
      warningViolations: 1,
      layersAnalyzed: 2,
      overallCompliance: 95,
    },
    ...overrides,
  }
}

/**
 * Creates a mock ComprehensiveAnalysisResult - In Progress state
 */
export function createInProgressAnalysis(
  overrides: Partial<ComprehensiveAnalysisResult> = {}
): ComprehensiveAnalysisResult {
  return {
    id: 'analysis-123',
    uploadId: 'upload-123',
    status: Status.PROCESSING,
    startedAt: new Date().toISOString(),
    progress: 45,
    currentStep: 'Analyzing dependencies...',
    ...overrides,
  }
}

/**
 * Creates a mock ComprehensiveAnalysisResult - Completed state
 * Note: Deep merges results object to allow partial overrides
 */
export function createCompletedAnalysis(
  overrides: Partial<ComprehensiveAnalysisResult> = {}
): ComprehensiveAnalysisResult {
  const { results: resultsOverrides, ...restOverrides } = overrides

  const defaultResults = {
    summary: createAnalysisSummary(),
    dependencyAnalysis: createDependencyAnalysisResults(),
    architectureValidation: createArchitectureValidationResults(),
    healthScore: createHealthScore(),
    bundleImpact: createBundleImpactReport(),
    duplicateDetection: createDuplicateDetectionResults(),
    circularDependencies: [createCircularDependency()],
    versionConflicts: [createVersionConflict()],
  }

  return {
    id: 'analysis-123',
    uploadId: 'upload-123',
    status: Status.COMPLETED,
    startedAt: new Date(Date.now() - 60000).toISOString(),
    completedAt: new Date().toISOString(),
    progress: 100,
    results: resultsOverrides ? { ...defaultResults, ...resultsOverrides } : defaultResults,
    warnings: [],
    ...restOverrides,
  }
}

/**
 * Creates a mock ComprehensiveAnalysisResult - Failed state
 */
export function createFailedAnalysis(
  overrides: Partial<ComprehensiveAnalysisResult> = {}
): ComprehensiveAnalysisResult {
  return {
    id: 'analysis-123',
    uploadId: 'upload-123',
    status: Status.FAILED,
    startedAt: new Date().toISOString(),
    progress: 0,
    error: 'Analysis failed: Invalid workspace configuration',
    ...overrides,
  }
}
