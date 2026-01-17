// API Types

export * from './analysis'
export type {
  DependencyEdge,
  DependencyGraph,
  DependencyType,
  PackageNode,
  WorkspaceType,
} from './analysis/graph'
export type {
  AnalysisMetadata,
  AnalysisResult,
  CheckResult,
  CircularDependencyInfo,
  FixStrategy,
  ValidationError,
  ValidationWarning,
} from './analysis/results'
// Re-export commonly used types
export type { ApiError, ApiResponse, PaginatedResponse } from './api'
export * from './api'
export type { AuthToken, LoginCredentials, User } from './auth'
export * from './auth'
export * from './common'

export type {
  ArchitectureValidation,
  BundleBreakdown,
  BundleImpactReport,
  ComprehensiveAnalysisResult,
  DependencyAnalysis,
  DuplicateDetectionResults,
  DuplicateGroup,
  FileProcessingResult,
  HealthScore,
  PackageJsonFile,
  Project,
  UploadedFile,
} from './domain'
export * from './domain'

// Re-export new WASM-compatible types
export type { ErrorCode, Result, ResultError } from './result'
// New WASM-compatible types
export * from './result'
export { ErrorCodes, isError, isSuccess } from './result'
export * from './wasm'

export type {
  MonoGuardAnalyzer,
  VersionInfo,
  WasmLoaderOptions,
  WasmLoadResult,
} from './wasm/adapter'
