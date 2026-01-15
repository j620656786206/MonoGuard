import { ApiResponse, FileProcessingResult } from '@monoguard/types';
import { apiClient } from '../client';

export interface UploadProgress {
  loaded: number;
  total: number;
  percentage: number;
}

export type UploadProgressCallback = (progress: UploadProgress) => void;

export class UploadService {
  /**
   * Upload files to the server
   */
  static async uploadFiles(
    files: File[], 
    onProgress?: UploadProgressCallback
  ): Promise<FileProcessingResult> {
    const formData = new FormData();
    
    // Add files to form data
    files.forEach((file) => {
      formData.append('files', file);
    });

    try {
      const response = await apiClient.post<ApiResponse<FileProcessingResult>>(
        '/api/v1/upload',
        formData,
        {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
          onUploadProgress: (progressEvent) => {
            if (onProgress && progressEvent.total) {
              const progress: UploadProgress = {
                loaded: progressEvent.loaded,
                total: progressEvent.total,
                percentage: Math.round((progressEvent.loaded * 100) / progressEvent.total)
              };
              onProgress(progress);
            }
          },
        }
      );

      // Ensure we handle the response structure properly
      return response.data.data || response.data;
    } catch (error) {
      console.error('Upload failed:', error);
      throw error;
    }
  }

  /**
   * Get upload result by ID
   */
  static async getUploadResult(id: string): Promise<FileProcessingResult> {
    try {
      const response = await apiClient.get<ApiResponse<FileProcessingResult>>(
        `/api/v1/upload/${id}`
      );
      return response.data.data;
    } catch (error) {
      console.error('Failed to get upload result:', error);
      throw error;
    }
  }

  /**
   * Download uploaded file
   */
  static async downloadFile(filename: string): Promise<Blob> {
    try {
      const response = await apiClient.get(
        `/api/v1/upload/files/${filename}`,
        {
          responseType: 'blob',
        }
      );
      return response.data;
    } catch (error) {
      console.error('Failed to download file:', error);
      throw error;
    }
  }

  /**
   * Cleanup old files
   */
  static async cleanupOldFiles(days: number = 7): Promise<string> {
    try {
      const response = await apiClient.post<ApiResponse<string>>(
        `/api/v1/upload/cleanup?days=${days}`
      );
      return response.data.data;
    } catch (error) {
      console.error('Failed to cleanup old files:', error);
      throw error;
    }
  }

  /**
   * Validate file before upload
   */
  static validateFile(file: File): { valid: boolean; error?: string } {
    const maxSize = 50 * 1024 * 1024; // 50MB
    const allowedExtensions = ['.zip', '.json'];
    const allowedMimeTypes = [
      'application/zip',
      'application/x-zip-compressed',
      'application/json',
      'text/json'
    ];

    // Check file size
    if (file.size > maxSize) {
      return {
        valid: false,
        error: `File size (${Math.round(file.size / (1024 * 1024))}MB) exceeds maximum allowed size (50MB)`
      };
    }

    // Check file extension
    const fileName = file.name.toLowerCase();
    const hasValidExtension = allowedExtensions.some(ext => 
      fileName.endsWith(ext) || fileName === 'package.json'
    );

    if (!hasValidExtension) {
      return {
        valid: false,
        error: `File type not allowed. Only .zip and package.json files are supported`
      };
    }

    // Check MIME type
    const hasValidMimeType = allowedMimeTypes.includes(file.type) || 
                           fileName === 'package.json' ||
                           fileName.endsWith('.json');

    if (!hasValidMimeType && file.type !== '') {
      return {
        valid: false,
        error: `Invalid file type. Only .zip and .json files are supported`
      };
    }

    return { valid: true };
  }

  /**
   * Validate multiple files
   */
  static validateFiles(files: File[]): { valid: boolean; errors: string[] } {
    const errors: string[] = [];

    if (files.length === 0) {
      errors.push('No files selected');
      return { valid: false, errors };
    }

    files.forEach((file, index) => {
      const validation = this.validateFile(file);
      if (!validation.valid) {
        errors.push(`File ${index + 1} (${file.name}): ${validation.error}`);
      }
    });

    return {
      valid: errors.length === 0,
      errors
    };
  }
}