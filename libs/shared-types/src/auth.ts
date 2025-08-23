import { z } from 'zod';
import { ID, ISODateString } from './common';

// User roles
export enum UserRole {
  OWNER = 'owner',
  ADMIN = 'admin',
  DEVELOPER = 'developer',
  VIEWER = 'viewer',
}

// Authentication providers
export enum AuthProvider {
  GITHUB = 'github',
  GITLAB = 'gitlab',
  BITBUCKET = 'bitbucket',
  EMAIL = 'email',
}

// User interface
export interface User {
  id: ID;
  email: string;
  name: string;
  avatar?: string;
  role: UserRole;
  provider: AuthProvider;
  providerId?: string;
  isActive: boolean;
  lastLoginAt?: ISODateString;
  createdAt: ISODateString;
  updatedAt: ISODateString;
}

// Auth token
export interface AuthToken {
  accessToken: string;
  refreshToken?: string;
  tokenType: string;
  expiresIn: number;
  expiresAt: ISODateString;
  scope?: string[];
}

// Login credentials
export interface LoginCredentials {
  email: string;
  password: string;
  rememberMe?: boolean;
}

// OAuth callback data
export interface OAuthCallback {
  provider: AuthProvider;
  code: string;
  state?: string;
  redirectUri: string;
}

// Session data
export interface Session {
  user: User;
  token: AuthToken;
  expiresAt: ISODateString;
}

// Password reset
export interface PasswordResetRequest {
  email: string;
}

export interface PasswordReset {
  token: string;
  password: string;
  confirmPassword: string;
}

// Registration
export interface UserRegistration {
  email: string;
  password: string;
  name: string;
  acceptTerms: boolean;
}

// Zod Schemas for validation
export const UserSchema = z.object({
  id: z.union([z.string(), z.number()]),
  email: z.string().email(),
  name: z.string().min(1),
  avatar: z.string().url().optional(),
  role: z.nativeEnum(UserRole),
  provider: z.nativeEnum(AuthProvider),
  providerId: z.string().optional(),
  isActive: z.boolean(),
  lastLoginAt: z.string().datetime().optional(),
  createdAt: z.string().datetime(),
  updatedAt: z.string().datetime(),
});

export const LoginCredentialsSchema = z.object({
  email: z.string().email(),
  password: z.string().min(8),
  rememberMe: z.boolean().optional(),
});

export const UserRegistrationSchema = z.object({
  email: z.string().email(),
  password: z.string().min(8),
  name: z.string().min(1),
  acceptTerms: z.boolean().refine((val) => val === true, {
    message: 'You must accept the terms and conditions',
  }),
});