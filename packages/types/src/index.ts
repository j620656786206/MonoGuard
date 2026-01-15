// API Types
export * from './api';
export * from './domain';
export * from './auth';
export * from './common';

// Re-export commonly used types
export type { 
  ApiResponse, 
  ApiError, 
  PaginatedResponse 
} from './api';

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
  BundleBreakdown
} from './domain';

export type {
  User,
  AuthToken,
  LoginCredentials
} from './auth';