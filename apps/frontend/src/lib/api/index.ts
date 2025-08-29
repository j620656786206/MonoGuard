/**
 * API Client Index
 * Central export point for all API-related functionality
 */

// Core client
export { apiClient, ApiError, type ApiResponse, type PaginatedResponse } from './client';

// Configuration
export { API_CONFIG, API_ENDPOINTS } from './config';

// Services
export { default as AuthService } from './services/auth';
export { default as ProjectsService } from './services/projects';
export { default as AnalysisService } from './services/analysis';

// Types from services
export type {
  LoginResponse,
  RegisterPayload,
  PasswordResetPayload,
  PasswordChangePayload,
  UpdateProfilePayload,
} from './services/auth';

export type {
  CreateProjectPayload,
  UpdateProjectPayload,
  ProjectListParams,
  AnalyzeProjectOptions,
} from './services/projects';

export type {
  AnalysisStatus,
  AnalysisResult,
  AnalysisIssue,
  AnalysisListParams,
  CreateAnalysisPayload,
} from './services/analysis';

// Utility functions
export const createQueryParams = (params: Record<string, any>): string => {
  const queryParams = new URLSearchParams();
  
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null) {
      queryParams.append(key, value.toString());
    }
  });
  
  return queryParams.toString();
};

/**
 * Generic pagination helper
 */
export const createPaginationParams = (
  page: number = 1,
  limit: number = 10
): { page: number; limit: number; offset: number } => ({
  page,
  limit,
  offset: (page - 1) * limit,
});

/**
 * Error handling utility
 */
export const handleApiError = (error: unknown): string => {
  if (error instanceof Error && error.name === 'ApiError') {
    return error.message;
  }
  
  if (error instanceof Error) {
    return error.message;
  }
  
  return 'An unexpected error occurred';
};

/**
 * Response validation utility
 */
export const validateApiResponse = <T>(
  response: any,
  expectedStatus: number = 200
): T => {
  if (response.status !== expectedStatus) {
    const error = new Error(
      response.message || `Expected status ${expectedStatus}, got ${response.status}`
    );
    error.name = 'ApiError';
    throw error;
  }
  
  return response.data;
};