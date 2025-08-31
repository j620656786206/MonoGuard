import axios, { 
  AxiosInstance, 
  AxiosRequestConfig, 
  AxiosResponse,
  AxiosError 
} from 'axios';
import { API_CONFIG } from './config';

/**
 * Custom error class for API errors
 */
export class ApiError extends Error {
  constructor(
    message: string,
    public status?: number,
    public code?: string,
    public details?: unknown
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

/**
 * Response wrapper for API calls
 */
export interface ApiResponse<T = any> {
  data: T;
  message?: string;
  status: number;
}

/**
 * Paginated response wrapper
 */
export interface PaginatedResponse<T = any> extends ApiResponse<T[]> {
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

/**
 * Request configuration with additional options
 */
export interface RequestConfig extends AxiosRequestConfig {
  skipAuth?: boolean;
  retries?: number;
}

/**
 * API Client class with authentication and error handling
 */
class ApiClient {
  private axiosInstance: AxiosInstance;
  private authToken: string | null = null;

  constructor() {
    this.axiosInstance = axios.create({
      baseURL: API_CONFIG.BASE_URL,
      timeout: API_CONFIG.TIMEOUT,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    this.setupInterceptors();
  }

  /**
   * Setup request and response interceptors
   */
  private setupInterceptors(): void {
    // Request interceptor for auth token
    this.axiosInstance.interceptors.request.use(
      (config) => {
        if (this.authToken && !config.headers?.skipAuth) {
          config.headers.Authorization = `Bearer ${this.authToken}`;
        }
        return config;
      },
      (error) => Promise.reject(error)
    );

    // Response interceptor for error handling
    this.axiosInstance.interceptors.response.use(
      (response: AxiosResponse) => response,
      async (error: AxiosError) => {
        const originalRequest = error.config;

        // Handle 401 Unauthorized - attempt token refresh
        if (error.response?.status === 401 && originalRequest && !originalRequest.headers?.skipAuth) {
          try {
            await this.refreshToken();
            // Retry the original request with new token
            return this.axiosInstance(originalRequest);
          } catch (refreshError) {
            this.clearAuth();
            throw this.createApiError(error);
          }
        }

        throw this.createApiError(error);
      }
    );
  }

  /**
   * Create standardized API error from axios error
   */
  private createApiError(error: AxiosError): ApiError {
    const response = error.response;
    const responseData = response?.data as any;
    const message = responseData?.message || error.message || 'An error occurred';
    const status = response?.status;
    const code = responseData?.code;
    
    return new ApiError(message, status, code, response?.data);
  }

  /**
   * Set authentication token
   */
  public setAuthToken(token: string): void {
    this.authToken = token;
    // Store in localStorage for persistence
    if (typeof window !== 'undefined') {
      localStorage.setItem('authToken', token);
    }
  }

  /**
   * Clear authentication token
   */
  public clearAuth(): void {
    this.authToken = null;
    if (typeof window !== 'undefined') {
      localStorage.removeItem('authToken');
    }
  }

  /**
   * Initialize auth token from storage
   */
  public initializeAuth(): void {
    if (typeof window !== 'undefined') {
      const token = localStorage.getItem('authToken');
      if (token) {
        this.authToken = token;
      }
    }
  }

  /**
   * Refresh authentication token
   */
  private async refreshToken(): Promise<void> {
    const refreshToken = typeof window !== 'undefined' 
      ? localStorage.getItem('refreshToken') 
      : null;
      
    if (!refreshToken) {
      throw new ApiError('No refresh token available');
    }

    const response = await this.axiosInstance.post('/auth/refresh', {
      refreshToken,
    }, {
      headers: { skipAuth: 'true' }
    });

    const { accessToken, refreshToken: newRefreshToken } = response.data;
    
    this.setAuthToken(accessToken);
    if (typeof window !== 'undefined' && newRefreshToken) {
      localStorage.setItem('refreshToken', newRefreshToken);
    }
  }

  /**
   * Generic GET request
   */
  public async get<T = any>(
    url: string, 
    config?: RequestConfig
  ): Promise<ApiResponse<T>> {
    const response = await this.axiosInstance.get(url, config);
    return response.data;
  }

  /**
   * Generic POST request
   */
  public async post<T = any, D = any>(
    url: string,
    data?: D,
    config?: RequestConfig
  ): Promise<ApiResponse<T>> {
    const response = await this.axiosInstance.post(url, data, config);
    return response.data;
  }

  /**
   * Generic PUT request
   */
  public async put<T = any, D = any>(
    url: string,
    data?: D,
    config?: RequestConfig
  ): Promise<ApiResponse<T>> {
    const response = await this.axiosInstance.put(url, data, config);
    return response.data;
  }

  /**
   * Generic PATCH request
   */
  public async patch<T = any, D = any>(
    url: string,
    data?: D,
    config?: RequestConfig
  ): Promise<ApiResponse<T>> {
    const response = await this.axiosInstance.patch(url, data, config);
    return response.data;
  }

  /**
   * Generic DELETE request
   */
  public async delete<T = any>(
    url: string,
    config?: RequestConfig
  ): Promise<ApiResponse<T>> {
    const response = await this.axiosInstance.delete(url, config);
    return response.data;
  }

  /**
   * Upload file with progress tracking
   */
  public async uploadFile<T = any>(
    url: string,
    file: File,
    onUploadProgress?: (progressEvent: any) => void,
    config?: RequestConfig
  ): Promise<ApiResponse<T>> {
    const formData = new FormData();
    formData.append('file', file);

    const response = await this.axiosInstance.post(url, formData, {
      ...config,
      headers: {
        ...config?.headers,
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress,
    });

    return response.data;
  }

  /**
   * Download file
   */
  public async downloadFile(
    url: string,
    filename?: string,
    config?: RequestConfig
  ): Promise<void> {
    const response = await this.axiosInstance.get(url, {
      ...config,
      responseType: 'blob',
    });

    const blob = new Blob([response.data]);
    const downloadUrl = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = downloadUrl;
    link.download = filename || 'download';
    document.body.appendChild(link);
    link.click();
    link.remove();
    window.URL.revokeObjectURL(downloadUrl);
  }

  /**
   * Set session token for API requests
   */
  public setSessionToken(token: string): void {
    this.axiosInstance.defaults.headers.common['X-Session-Token'] = token;
  }

  /**
   * Clear session token from API requests
   */
  public clearSessionToken(): void {
    delete this.axiosInstance.defaults.headers.common['X-Session-Token'];
  }
}

// Create and export a singleton instance
export const apiClient = new ApiClient();

// Initialize auth on client side
if (typeof window !== 'undefined') {
  apiClient.initializeAuth();
}