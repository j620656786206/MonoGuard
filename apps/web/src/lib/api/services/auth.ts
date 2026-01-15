import { apiClient, ApiResponse } from '../client';
import { API_ENDPOINTS } from '../config';
import { User, LoginCredentials, AuthToken } from '@monoguard/types';

/**
 * Login response interface
 */
export interface LoginResponse {
  user: User;
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
}

/**
 * Register payload interface
 */
export interface RegisterPayload {
  name: string;
  email: string;
  password: string;
  confirmPassword: string;
}

/**
 * Password reset payload
 */
export interface PasswordResetPayload {
  email: string;
}

/**
 * Password change payload
 */
export interface PasswordChangePayload {
  currentPassword: string;
  newPassword: string;
  confirmPassword: string;
}

/**
 * Update profile payload
 */
export interface UpdateProfilePayload {
  name?: string;
  email?: string;
  avatar?: string;
  preferences?: Record<string, any>;
}

/**
 * Authentication API service
 */
export class AuthService {
  /**
   * Login user with email and password
   */
  static async login(credentials: LoginCredentials): Promise<ApiResponse<LoginResponse>> {
    const response = await apiClient.post<LoginResponse, LoginCredentials>(
      API_ENDPOINTS.AUTH.LOGIN,
      credentials,
      { skipAuth: true }
    );

    // Set auth token in client
    if (response.data.accessToken) {
      apiClient.setAuthToken(response.data.accessToken);
      
      // Store refresh token
      if (typeof window !== 'undefined' && response.data.refreshToken) {
        localStorage.setItem('refreshToken', response.data.refreshToken);
      }
    }

    return response;
  }

  /**
   * Register new user
   */
  static async register(payload: RegisterPayload): Promise<ApiResponse<LoginResponse>> {
    const response = await apiClient.post<LoginResponse, RegisterPayload>(
      '/auth/register',
      payload,
      { skipAuth: true }
    );

    // Set auth token in client
    if (response.data.accessToken) {
      apiClient.setAuthToken(response.data.accessToken);
      
      // Store refresh token
      if (typeof window !== 'undefined' && response.data.refreshToken) {
        localStorage.setItem('refreshToken', response.data.refreshToken);
      }
    }

    return response;
  }

  /**
   * Logout user
   */
  static async logout(): Promise<ApiResponse<void>> {
    try {
      await apiClient.post<void>(API_ENDPOINTS.AUTH.LOGOUT);
    } catch (error) {
      // Continue with cleanup even if logout request fails
      console.warn('Logout request failed:', error);
    } finally {
      // Clear auth token and storage
      apiClient.clearAuth();
      if (typeof window !== 'undefined') {
        localStorage.removeItem('refreshToken');
      }
    }

    return { data: undefined, status: 200 };
  }

  /**
   * Get current user profile
   */
  static async getCurrentUser(): Promise<ApiResponse<User>> {
    return apiClient.get<User>(API_ENDPOINTS.AUTH.ME);
  }

  /**
   * Refresh authentication token
   */
  static async refreshToken(): Promise<ApiResponse<AuthToken>> {
    const refreshToken = typeof window !== 'undefined' 
      ? localStorage.getItem('refreshToken') 
      : null;

    if (!refreshToken) {
      throw new Error('No refresh token available');
    }

    const response = await apiClient.post<AuthToken>(
      API_ENDPOINTS.AUTH.REFRESH,
      { refreshToken },
      { skipAuth: true }
    );

    // Update auth token
    if (response.data.accessToken) {
      apiClient.setAuthToken(response.data.accessToken);
      
      // Update refresh token if provided
      if (typeof window !== 'undefined' && response.data.refreshToken) {
        localStorage.setItem('refreshToken', response.data.refreshToken);
      }
    }

    return response;
  }

  /**
   * Update user profile
   */
  static async updateProfile(payload: UpdateProfilePayload): Promise<ApiResponse<User>> {
    return apiClient.put<User, UpdateProfilePayload>('/auth/profile', payload);
  }

  /**
   * Change password
   */
  static async changePassword(payload: PasswordChangePayload): Promise<ApiResponse<void>> {
    return apiClient.post<void, PasswordChangePayload>('/auth/change-password', payload);
  }

  /**
   * Request password reset
   */
  static async requestPasswordReset(payload: PasswordResetPayload): Promise<ApiResponse<void>> {
    return apiClient.post<void, PasswordResetPayload>(
      '/auth/forgot-password',
      payload,
      { skipAuth: true }
    );
  }

  /**
   * Reset password with token
   */
  static async resetPassword(
    token: string,
    newPassword: string,
    confirmPassword: string
  ): Promise<ApiResponse<void>> {
    return apiClient.post<void>(
      `/auth/reset-password/${token}`,
      { newPassword, confirmPassword },
      { skipAuth: true }
    );
  }

  /**
   * Verify email address
   */
  static async verifyEmail(token: string): Promise<ApiResponse<void>> {
    return apiClient.post<void>(
      `/auth/verify-email/${token}`,
      {},
      { skipAuth: true }
    );
  }

  /**
   * Resend email verification
   */
  static async resendEmailVerification(): Promise<ApiResponse<void>> {
    return apiClient.post<void>('/auth/resend-verification');
  }

  /**
   * Check if user is authenticated
   */
  static isAuthenticated(): boolean {
    if (typeof window === 'undefined') return false;
    
    const token = localStorage.getItem('authToken');
    return !!token;
  }

  /**
   * Get stored auth token
   */
  static getAuthToken(): string | null {
    if (typeof window === 'undefined') return null;
    
    return localStorage.getItem('authToken');
  }
}

// Export instance methods as default export
export default AuthService;