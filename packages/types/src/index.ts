// API Types
export * from './api';
export * from './domain';
export * from './auth';
export * from './common';

// New WASM-compatible types
export * from './result';
export * from './analysis';
export * from './wasm';

// Re-export commonly used types
export type { ApiResponse, ApiError, PaginatedResponse } from './api';

export type {
  Project,
  DependencyAnalysis,
  ArchitectureValidation,
  HealthScore,
  ComprehensiveAnalysisResult,
  FileProcessingResult,
  PackageJsonFile,
  UploadedFile,
  DuplicateDetectionResults,
  DuplicateGroup,
  BundleImpactReport,
  BundleBreakdown,
} from './domain';

export type { User, AuthToken, LoginCredentials } from './auth';

// Re-export new WASM-compatible types
export type { Result, ResultError, ErrorCode } from './result';

export { ErrorCodes, isSuccess, isError } from './result';

export type {
  DependencyGraph,
  PackageNode,
  DependencyEdge,
  DependencyType,
  WorkspaceType,
} from './analysis/graph';

export type {
  AnalysisResult,
  CircularDependencyInfo,
  FixStrategy,
  CheckResult,
  ValidationError,
  ValidationWarning,
  AnalysisMetadata,
} from './analysis/results';

export type {
  MonoGuardAnalyzer,
  VersionInfo,
  WasmLoaderOptions,
  WasmLoadResult,
} from './wasm/adapter';
