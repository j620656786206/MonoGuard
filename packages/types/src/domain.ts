import { z } from 'zod'
import { type ID, type ISODateString, type RiskLevel, type Severity, Status } from './common'

// Project management
export interface Project {
  id: ID
  name: string
  description?: string
  repositoryUrl: string
  branch: string
  status: Status
  healthScore: number
  lastAnalysisAt?: ISODateString
  ownerId: ID
  settings: ProjectSettings
  createdAt: ISODateString
  updatedAt: ISODateString
}

export interface ProjectSettings {
  autoAnalysis: boolean
  analysisSchedule?: string // cron expression
  notificationSettings: NotificationSettings
  excludePatterns: string[]
  includePatterns: string[]
  architectureRules: ArchitectureRules
}

export interface NotificationSettings {
  email: boolean
  webhook?: string
  slackWebhook?: string
  severity: Severity[]
}

// Architecture rules configuration
export interface ArchitectureRules {
  layers: ArchitectureLayer[]
  rules: ArchitectureRule[]
}

export interface ArchitectureLayer {
  name: string
  pattern: string
  description: string
  canImport: string[]
  cannotImport: string[]
}

export interface ArchitectureRule {
  name: string
  severity: Severity
  description: string
  pattern?: string
  enabled: boolean
}

// Analysis results
export interface DependencyAnalysis {
  id: ID
  projectId: ID
  status: Status
  startedAt: ISODateString
  completedAt?: ISODateString
  results: DependencyAnalysisResults
  metadata: AnalysisMetadata
}

export interface DependencyAnalysisResults {
  duplicateDependencies: DuplicateDependency[]
  versionConflicts: VersionConflict[]
  unusedDependencies: UnusedDependency[]
  circularDependencies: CircularDependency[]
  bundleImpact: BundleImpactReport
  summary: AnalysisSummary
}

export interface DuplicateDependency {
  packageName: string
  versions: string[]
  affectedPackages: string[]
  estimatedWaste: string
  riskLevel: RiskLevel
  recommendation: string
  migrationSteps: string[]
}

export interface VersionConflict {
  packageName: string
  conflictingVersions: ConflictingVersion[]
  riskLevel: RiskLevel
  resolution: string
  impact: string
}

export interface ConflictingVersion {
  version: string
  packages: string[]
  isBreaking: boolean
}

export interface UnusedDependency {
  packageName: string
  version: string
  packagePath: string
  sizeImpact: string
  lastUsed?: ISODateString
  confidence: number
}

export interface CircularDependency {
  cycle: string[]
  type: 'direct' | 'indirect'
  severity: Severity
  impact: string
}

export interface BundleImpactReport {
  totalSize: string
  duplicateSize: string
  unusedSize: string
  potentialSavings: string
  breakdown: BundleBreakdown[]
}

export interface BundleBreakdown {
  packageName: string
  size: string
  percentage: number
  duplicates: number
}

export interface AnalysisSummary {
  totalPackages: number
  duplicateCount: number
  conflictCount: number
  unusedCount: number
  circularCount: number
  healthScore: number
}

// Architecture validation
export interface ArchitectureValidation {
  id: ID
  projectId: ID
  status: Status
  startedAt: ISODateString
  completedAt?: ISODateString
  results: ArchitectureValidationResults
  metadata: AnalysisMetadata
}

export interface ArchitectureValidationResults {
  violations: ArchitectureViolation[]
  layerCompliance: LayerCompliance[]
  circularDependencies: CircularDependency[]
  summary: ValidationSummary
}

export interface ArchitectureViolation {
  ruleName: string
  severity: Severity
  description: string
  violatingFile: string
  violatingImport: string
  expectedLayer: string
  actualLayer: string
  suggestion: string
}

export interface LayerCompliance {
  layerName: string
  totalFiles: number
  compliantFiles: number
  violationCount: number
  compliancePercentage: number
}

export interface ValidationSummary {
  totalViolations: number
  criticalViolations: number
  warningViolations: number
  layersAnalyzed: number
  overallCompliance: number
}

// Analysis metadata
export interface AnalysisMetadata {
  version: string
  duration: number
  filesProcessed: number
  packagesAnalyzed: number
  configuration: Record<string, any>
  environment: {
    nodeVersion: string
    platform: string
    memoryUsage: string
    cpuUsage: string
  }
}

// Health score calculation
export interface HealthScore {
  overall: number
  dependencies: number
  architecture: number
  maintainability: number
  security: number
  performance: number
  lastUpdated: ISODateString
  trend: 'improving' | 'stable' | 'declining'
  factors: HealthFactor[]
}

export interface HealthFactor {
  name: string
  score: number
  weight: number
  description: string
  recommendations: string[]
}

// Comprehensive analysis results for frontend display
export interface ComprehensiveAnalysisResult {
  id: ID
  uploadId: ID
  status: Status
  startedAt: ISODateString
  completedAt?: ISODateString
  progress: number
  currentStep?: string
  results?: {
    dependencyAnalysis?: DependencyAnalysisResults
    architectureValidation?: ArchitectureValidationResults
    healthScore?: HealthScore
    bundleImpact?: BundleImpactReport
    duplicateDetection?: DuplicateDetectionResults
  }
  error?: string
  warnings?: string[]
}

// Enhanced duplicate detection results
export interface DuplicateDetectionResults {
  totalDuplicates: number
  potentialSavings: string
  duplicateGroups: DuplicateGroup[]
  recommendations: DuplicateRecommendation[]
}

export interface DuplicateGroup {
  packageName: string
  versions: DuplicateVersionInfo[]
  totalSize: string
  wastedSize: string
  riskLevel: RiskLevel
  affectedPackages: string[]
}

export interface DuplicateVersionInfo {
  version: string
  size: string
  usageCount: number
  packages: string[]
  isRecommended: boolean
}

export interface DuplicateRecommendation {
  type: 'consolidate' | 'upgrade' | 'replace' | 'remove'
  packageName: string
  description: string
  estimatedSavings: string
  difficulty: 'easy' | 'medium' | 'hard'
  steps: string[]
}

// File processing result extension for analysis
export interface FileProcessingResult {
  uploadId: ID
  files: UploadedFile[]
  packageJsonFiles: PackageJsonFile[]
  analysisReady: boolean
  errors?: string[]
  metadata: {
    totalSize: number
    processedAt: ISODateString
    processingDuration: number
  }
}

export interface UploadedFile {
  filename: string
  originalName: string
  fileSize: number
  mimeType: string
  path: string
  extractedFiles?: number
}

export interface PackageJsonFile {
  name?: string
  version?: string
  path: string
  dependencies?: Record<string, string>
  devDependencies?: Record<string, string>
  peerDependencies?: Record<string, string>
  optionalDependencies?: Record<string, string>
  scripts?: Record<string, string>
  metadata?: {
    fileSize: number
    lastModified: string
    dependencyCount: number
    devDependencyCount: number
  }
}

// Zod Schemas
export const ProjectSchema = z.object({
  id: z.union([z.string(), z.number()]),
  name: z.string().min(1),
  description: z.string().optional(),
  repositoryUrl: z.string().url(),
  branch: z.string(),
  status: z.nativeEnum(Status),
  healthScore: z.number().min(0).max(100),
  lastAnalysisAt: z.string().datetime().optional(),
  ownerId: z.union([z.string(), z.number()]),
  createdAt: z.string().datetime(),
  updatedAt: z.string().datetime(),
})

export const DependencyAnalysisSchema = z.object({
  id: z.union([z.string(), z.number()]),
  projectId: z.union([z.string(), z.number()]),
  status: z.nativeEnum(Status),
  startedAt: z.string().datetime(),
  completedAt: z.string().datetime().optional(),
})
