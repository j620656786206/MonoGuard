import { apiClient, ApiResponse, PaginatedResponse } from '../client';
import { API_ENDPOINTS } from '../config';
import { Project, DependencyAnalysis, ArchitectureValidation } from '@monoguard/shared-types';

/**
 * Project creation payload
 */
export interface CreateProjectPayload {
  name: string;
  description?: string;
  repositoryUrl: string;
  branch: string;
  ownerId: string;
  framework?: string;
  language?: string;
  packageManager?: 'npm' | 'yarn' | 'pnpm';
}

/**
 * Project update payload
 */
export interface UpdateProjectPayload {
  name?: string;
  description?: string;
  repositoryUrl?: string;
}

/**
 * Project list query parameters
 */
export interface ProjectListParams {
  page?: number;
  limit?: number;
  search?: string;
  framework?: string;
  language?: string;
  sortBy?: 'name' | 'createdAt' | 'updatedAt' | 'healthScore';
  sortOrder?: 'asc' | 'desc';
}

/**
 * Project analysis options
 */
export interface AnalyzeProjectOptions {
  includeDependencies?: boolean;
  includeArchitecture?: boolean;
  includePerformance?: boolean;
  skipCache?: boolean;
}

/**
 * Projects API service
 */
export class ProjectsService {
  /**
   * Get list of projects with optional filtering and pagination
   */
  static async getProjects(params?: ProjectListParams): Promise<PaginatedResponse<Project>> {
    const queryParams = new URLSearchParams();
    
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined) {
          queryParams.append(key, value.toString());
        }
      });
    }

    const url = `${API_ENDPOINTS.PROJECTS.LIST}?${queryParams.toString()}`;
    return apiClient.get<Project[]>(url);
  }

  /**
   * Get a single project by ID
   */
  static async getProject(id: string): Promise<ApiResponse<Project>> {
    return apiClient.get<Project>(API_ENDPOINTS.PROJECTS.GET(id));
  }

  /**
   * Create a new project
   */
  static async createProject(payload: CreateProjectPayload): Promise<ApiResponse<Project>> {
    return apiClient.post<Project, CreateProjectPayload>(
      API_ENDPOINTS.PROJECTS.CREATE,
      payload
    );
  }

  /**
   * Update an existing project
   */
  static async updateProject(
    id: string, 
    payload: UpdateProjectPayload
  ): Promise<ApiResponse<Project>> {
    return apiClient.put<Project, UpdateProjectPayload>(
      API_ENDPOINTS.PROJECTS.UPDATE(id),
      payload
    );
  }

  /**
   * Delete a project
   */
  static async deleteProject(id: string): Promise<ApiResponse<void>> {
    return apiClient.delete<void>(API_ENDPOINTS.PROJECTS.DELETE(id));
  }

  /**
   * Analyze a project
   */
  static async analyzeProject(
    id: string, 
    options?: AnalyzeProjectOptions
  ): Promise<ApiResponse<{ analysisId: string }>> {
    return apiClient.post<{ analysisId: string }, AnalyzeProjectOptions>(
      API_ENDPOINTS.PROJECTS.ANALYZE(id),
      options || {}
    );
  }

  /**
   * Get project dependencies analysis
   */
  static async getProjectDependencies(id: string): Promise<ApiResponse<DependencyAnalysis>> {
    return apiClient.get<DependencyAnalysis>(`${API_ENDPOINTS.PROJECTS.GET(id)}/dependencies`);
  }

  /**
   * Get project architecture validation
   */
  static async getProjectArchitecture(id: string): Promise<ApiResponse<ArchitectureValidation>> {
    return apiClient.get<ArchitectureValidation>(`${API_ENDPOINTS.PROJECTS.GET(id)}/architecture`);
  }

  /**
   * Upload project files for analysis
   */
  static async uploadProjectFiles(
    id: string,
    files: File[],
    onUploadProgress?: (progressEvent: any) => void
  ): Promise<ApiResponse<{ uploadId: string }>> {
    const formData = new FormData();
    files.forEach((file, index) => {
      formData.append(`file-${index}`, file);
    });

    return apiClient.post<{ uploadId: string }>(
      `${API_ENDPOINTS.PROJECTS.GET(id)}/upload`,
      formData,
      {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
        onUploadProgress,
      }
    );
  }

  /**
   * Get project health metrics
   */
  static async getProjectHealth(id: string): Promise<ApiResponse<{
    healthScore: number;
    metrics: {
      dependencies: number;
      vulnerabilities: number;
      codeQuality: number;
      performance: number;
      maintainability: number;
    };
    lastUpdated: string;
  }>> {
    return apiClient.get(`${API_ENDPOINTS.PROJECTS.GET(id)}/health`);
  }
}

// Export instance methods as default export
export default ProjectsService;