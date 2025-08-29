/**
 * API Configuration
 * Centralized configuration for API endpoints and settings
 */

export const API_CONFIG = {
  BASE_URL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080',
  TIMEOUT: 30000, // 30 seconds
  RETRY_COUNT: 3,
  RETRY_DELAY: 1000, // 1 second
} as const;

export const API_ENDPOINTS = {
  // Authentication
  AUTH: {
    LOGIN: '/auth/login',
    LOGOUT: '/auth/logout',
    REFRESH: '/auth/refresh',
    ME: '/auth/me',
  },
  
  // Projects
  PROJECTS: {
    LIST: '/api/v1/projects',
    CREATE: '/api/v1/projects',
    GET: (id: string) => `/api/v1/projects/${id}`,
    UPDATE: (id: string) => `/api/v1/projects/${id}`,
    DELETE: (id: string) => `/api/v1/projects/${id}`,
    ANALYZE: (id: string) => `/api/v1/projects/${id}/analyze`,
  },

  // Analysis
  ANALYSIS: {
    LIST: '/api/v1/analysis',
    GET: (id: string) => `/api/v1/analysis/${id}`,
    CREATE: '/api/v1/analysis',
    RESULTS: (id: string) => `/api/v1/analysis/${id}/results`,
    DEPENDENCIES: (id: string) => `/api/v1/analysis/${id}/dependencies`,
    ARCHITECTURE: (id: string) => `/api/v1/analysis/${id}/architecture`,
    UPLOAD: '/api/v1/analysis/upload',
    COMPREHENSIVE: (uploadId: string) => `/api/v1/analysis/comprehensive/${uploadId}`,
    PROGRESS: (analysisId: string) => `/api/v1/analysis/${analysisId}/progress`,
  },

  // Health
  HEALTH: {
    CHECK: '/health',
    METRICS: '/health/metrics',
  },
} as const;

export type ApiEndpoint = string;