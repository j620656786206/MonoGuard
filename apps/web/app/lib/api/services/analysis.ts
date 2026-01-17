import type {
  ArchitectureValidation,
  ComprehensiveAnalysisResult,
  DependencyAnalysis,
} from '@monoguard/types'
import { type ApiResponse, apiClient, type PaginatedResponse } from '../client'
import { API_ENDPOINTS } from '../config'

/**
 * Analysis status enum
 */
export type AnalysisStatus = 'pending' | 'running' | 'completed' | 'failed'

/**
 * Analysis result interface
 */
export interface AnalysisResult {
  id: string
  projectId: string
  status: AnalysisStatus
  createdAt: string
  updatedAt: string
  completedAt?: string
  progress: number
  results?: {
    dependencies?: DependencyAnalysis
    architecture?: ArchitectureValidation
    healthScore?: number
    issues: AnalysisIssue[]
  }
  error?: string
}

/**
 * Analysis issue interface
 */
export interface AnalysisIssue {
  id: string
  type: 'warning' | 'error' | 'info'
  category: 'dependency' | 'architecture' | 'performance' | 'security'
  title: string
  description: string
  file?: string
  line?: number
  severity: 'low' | 'medium' | 'high' | 'critical'
  recommendation?: string
}

/**
 * Analysis list query parameters
 */
export interface AnalysisListParams {
  page?: number
  limit?: number
  projectId?: string
  status?: AnalysisStatus
  sortBy?: 'createdAt' | 'updatedAt' | 'progress' | 'status'
  sortOrder?: 'asc' | 'desc'
}

/**
 * Create analysis payload
 */
export interface CreateAnalysisPayload {
  projectId: string
  includeDependencies?: boolean
  includeArchitecture?: boolean
  includePerformance?: boolean
  skipCache?: boolean
  options?: Record<string, any>
}

/**
 * Analysis API service
 */
export class AnalysisService {
  /**
   * Get list of analyses with optional filtering and pagination
   */
  static async getAnalyses(
    params?: AnalysisListParams
  ): Promise<PaginatedResponse<AnalysisResult>> {
    const queryParams = new URLSearchParams()

    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined) {
          queryParams.append(key, value.toString())
        }
      })
    }

    const url = `${API_ENDPOINTS.ANALYSIS.LIST}?${queryParams.toString()}`
    return {
      ...(await apiClient.get<AnalysisResult[]>(url)),
      pagination: {
        page: 1,
        limit: 10,
        total: 0,
        totalPages: 1,
      },
    }
  }

  /**
   * Get a single analysis by ID
   */
  static async getAnalysis(id: string): Promise<ApiResponse<AnalysisResult>> {
    return apiClient.get<AnalysisResult>(API_ENDPOINTS.ANALYSIS.GET(id))
  }

  /**
   * Create a new analysis
   */
  static async createAnalysis(
    payload: CreateAnalysisPayload
  ): Promise<ApiResponse<AnalysisResult>> {
    return apiClient.post<AnalysisResult, CreateAnalysisPayload>(
      API_ENDPOINTS.ANALYSIS.CREATE,
      payload
    )
  }

  /**
   * Get analysis results
   */
  static async getAnalysisResults(id: string): Promise<ApiResponse<AnalysisResult['results']>> {
    return apiClient.get<AnalysisResult['results']>(API_ENDPOINTS.ANALYSIS.RESULTS(id))
  }

  /**
   * Get analysis dependencies
   */
  static async getAnalysisDependencies(id: string): Promise<ApiResponse<DependencyAnalysis>> {
    return apiClient.get<DependencyAnalysis>(API_ENDPOINTS.ANALYSIS.DEPENDENCIES(id))
  }

  /**
   * Get analysis architecture validation
   */
  static async getAnalysisArchitecture(id: string): Promise<ApiResponse<ArchitectureValidation>> {
    return apiClient.get<ArchitectureValidation>(API_ENDPOINTS.ANALYSIS.ARCHITECTURE(id))
  }

  /**
   * Cancel a running analysis
   */
  static async cancelAnalysis(id: string): Promise<ApiResponse<void>> {
    return apiClient.post<void>(`${API_ENDPOINTS.ANALYSIS.GET(id)}/cancel`)
  }

  /**
   * Retry a failed analysis
   */
  static async retryAnalysis(id: string): Promise<ApiResponse<AnalysisResult>> {
    return apiClient.post<AnalysisResult>(`${API_ENDPOINTS.ANALYSIS.GET(id)}/retry`)
  }

  /**
   * Delete an analysis
   */
  static async deleteAnalysis(id: string): Promise<ApiResponse<void>> {
    return apiClient.delete<void>(API_ENDPOINTS.ANALYSIS.GET(id))
  }

  /**
   * Get analysis progress (for real-time updates)
   */
  static async getAnalysisProgress(id: string): Promise<
    ApiResponse<{
      status: AnalysisStatus
      progress: number
      currentStep?: string
      estimatedCompletion?: string
    }>
  > {
    return apiClient.get(`${API_ENDPOINTS.ANALYSIS.GET(id)}/progress`)
  }

  /**
   * Get analysis statistics for a project
   */
  static async getProjectAnalysisStats(projectId: string): Promise<
    ApiResponse<{
      totalAnalyses: number
      successfulAnalyses: number
      failedAnalyses: number
      avgHealthScore: number
      lastAnalysis?: string
      trends: {
        healthScore: Array<{ date: string; value: number }>
        issueCount: Array<{ date: string; value: number }>
      }
    }>
  > {
    return apiClient.get(`/projects/${projectId}/analysis/stats`)
  }

  /**
   * Export analysis results
   */
  static async exportAnalysisResults(
    id: string,
    format: 'json' | 'pdf' | 'csv' = 'json'
  ): Promise<void> {
    const filename = `analysis-${id}.${format}`
    return apiClient.downloadFile(
      `${API_ENDPOINTS.ANALYSIS.GET(id)}/export?format=${format}`,
      filename
    )
  }

  /**
   * Compare two analyses
   */
  static async compareAnalyses(
    analysisId1: string,
    analysisId2: string
  ): Promise<
    ApiResponse<{
      comparison: {
        healthScoreDiff: number
        issuesDiff: {
          added: AnalysisIssue[]
          removed: AnalysisIssue[]
          modified: AnalysisIssue[]
        }
        dependenciesDiff: {
          added: string[]
          removed: string[]
          updated: string[]
        }
      }
    }>
  > {
    return apiClient.get(`/analysis/compare?analysis1=${analysisId1}&analysis2=${analysisId2}`)
  }

  /**
   * Start comprehensive analysis for uploaded files
   */
  static async startComprehensiveAnalysis(
    uploadId: string
  ): Promise<ApiResponse<ComprehensiveAnalysisResult>> {
    return apiClient.post<ComprehensiveAnalysisResult>(
      API_ENDPOINTS.ANALYSIS.COMPREHENSIVE(uploadId),
      {}
    )
  }

  /**
   * Get comprehensive analysis results
   */
  static async getComprehensiveAnalysis(
    analysisId: string
  ): Promise<ApiResponse<ComprehensiveAnalysisResult>> {
    return apiClient.get<ComprehensiveAnalysisResult>(API_ENDPOINTS.ANALYSIS.GET(analysisId))
  }

  /**
   * Get real-time analysis progress
   */
  static async getAnalysisProgressDetailed(analysisId: string): Promise<
    ApiResponse<{
      id: string
      status: AnalysisStatus
      progress: number
      currentStep: string
      steps: {
        name: string
        status: 'pending' | 'running' | 'completed' | 'failed'
        startedAt?: string
        completedAt?: string
        progress: number
      }[]
      estimatedCompletion?: string
      error?: string
    }>
  > {
    return apiClient.get(API_ENDPOINTS.ANALYSIS.PROGRESS(analysisId))
  }

  /**
   * Poll for analysis completion with timeout
   */
  static async pollAnalysisCompletion(
    analysisId: string,
    timeoutMs: number = 300000 // 5 minutes
  ): Promise<ComprehensiveAnalysisResult> {
    const startTime = Date.now()

    return new Promise((resolve, reject) => {
      const poll = async () => {
        try {
          if (Date.now() - startTime > timeoutMs) {
            reject(new Error('Analysis polling timeout'))
            return
          }

          const response = await AnalysisService.getComprehensiveAnalysis(analysisId)

          if (response.data.status === 'completed') {
            resolve(response.data)
            return
          }

          if (response.data.status === 'failed') {
            reject(new Error(response.data.error || 'Analysis failed'))
            return
          }

          // Continue polling
          setTimeout(poll, 2000) // Poll every 2 seconds
        } catch (error) {
          reject(error)
        }
      }

      poll()
    })
  }
}

// Export instance methods as default export
export default AnalysisService
