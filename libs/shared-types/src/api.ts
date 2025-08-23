import { z } from 'zod';
import { PaginationMeta, BaseError } from './common';

// API Response wrapper
export interface ApiResponse<T = any> {
  success: boolean;
  data: T;
  message?: string;
  timestamp: string;
}

export interface ApiError extends BaseError {
  statusCode: number;
  path: string;
  timestamp: string;
}

// Paginated response
export interface PaginatedResponse<T = any> extends ApiResponse<T[]> {
  pagination: PaginationMeta;
}

// API Request/Response schemas
export const ApiResponseSchema = <T extends z.ZodType>(dataSchema: T) =>
  z.object({
    success: z.boolean(),
    data: dataSchema,
    message: z.string().optional(),
    timestamp: z.string().datetime(),
  });

export const PaginatedResponseSchema = <T extends z.ZodType>(dataSchema: T) =>
  z.object({
    success: z.boolean(),
    data: z.array(dataSchema),
    message: z.string().optional(),
    timestamp: z.string().datetime(),
    pagination: z.object({
      currentPage: z.number(),
      totalPages: z.number(),
      totalItems: z.number(),
      itemsPerPage: z.number(),
      hasNextPage: z.boolean(),
      hasPreviousPage: z.boolean(),
    }),
  });

// Query parameters
export interface QueryParams {
  page?: number;
  limit?: number;
  sort?: string;
  order?: 'asc' | 'desc';
  search?: string;
  filters?: Record<string, any>;
}

// HTTP Methods
export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';

// API Endpoint configuration
export interface ApiEndpoint {
  method: HttpMethod;
  path: string;
  authenticated?: boolean;
  roles?: string[];
}