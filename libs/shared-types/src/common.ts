import { z } from 'zod';

// Common utility types
export type Nullable<T> = T | null;
export type Optional<T> = T | undefined;
export type ID = string | number;

// Status enums
export enum Status {
  PENDING = 'pending',
  IN_PROGRESS = 'in_progress',
  COMPLETED = 'completed',
  FAILED = 'failed',
}

export enum Severity {
  LOW = 'low',
  MEDIUM = 'medium',
  HIGH = 'high',
  CRITICAL = 'critical',
}

export enum RiskLevel {
  LOW = 'low',
  MEDIUM = 'medium',
  HIGH = 'high',
  CRITICAL = 'critical',
}

// Date utilities
export const DateSchema = z.string().datetime();
export type ISODateString = z.infer<typeof DateSchema>;

// Generic pagination types
export interface PaginationMeta {
  currentPage: number;
  totalPages: number;
  totalItems: number;
  itemsPerPage: number;
  hasNextPage: boolean;
  hasPreviousPage: boolean;
}

// Error types
export interface BaseError {
  code: string;
  message: string;
  details?: Record<string, any>;
}

export interface ValidationError extends BaseError {
  field: string;
  value: any;
}

// Progress tracking
export interface ProgressStatus {
  current: number;
  total: number;
  percentage: number;
  status: Status;
  message?: string;
}