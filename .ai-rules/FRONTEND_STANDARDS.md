# MonoGuard Frontend Development Standards

## Overview

This document establishes comprehensive frontend development standards for the MonoGuard project, built with Next.js 14, TypeScript, and modern React patterns. These standards ensure code quality, maintainability, performance, and developer experience consistency across the frontend codebase.

## Project Structure

### Directory Organization
```
frontend/
├── src/
│   ├── app/                    # Next.js App Router (pages and layouts)
│   │   ├── dashboard/
│   │   │   ├── page.tsx       # Dashboard page component
│   │   │   └── loading.tsx    # Loading UI for dashboard
│   │   ├── dependencies/
│   │   ├── architecture/
│   │   ├── layout.tsx         # Root layout
│   │   └── page.tsx          # Home page
│   ├── components/            # Reusable UI components
│   │   ├── ui/               # Base UI components (shadcn/ui)
│   │   │   ├── button.tsx
│   │   │   ├── card.tsx
│   │   │   ├── input.tsx
│   │   │   └── index.ts      # Barrel export
│   │   ├── dashboard/        # Dashboard-specific components
│   │   │   ├── HealthScoreCard.tsx
│   │   │   ├── TrendChart.tsx
│   │   │   └── IssuesSummary.tsx
│   │   ├── dependency/       # Dependency analysis components
│   │   │   ├── DependencyGraph.tsx
│   │   │   ├── DuplicatesList.tsx
│   │   │   └── ConflictResolver.tsx
│   │   ├── architecture/     # Architecture validation components
│   │   │   ├── LayerDiagram.tsx
│   │   │   ├── ViolationsList.tsx
│   │   │   └── RuleEditor.tsx
│   │   └── common/          # Shared components across modules
│   │       ├── Header.tsx
│   │       ├── Sidebar.tsx
│   │       ├── LoadingSpinner.tsx
│   │       └── ErrorBoundary.tsx
│   ├── hooks/               # Custom React hooks
│   │   ├── useProjects.ts
│   │   ├── useAnalysis.ts
│   │   ├── useWebSocket.ts
│   │   └── useLocalStorage.ts
│   ├── lib/                 # Utility libraries and configurations
│   │   ├── api/            # API client and utilities
│   │   │   ├── client.ts
│   │   │   ├── types.ts
│   │   │   └── endpoints.ts
│   │   ├── utils/          # General utility functions
│   │   │   ├── cn.ts      # Class name utility (clsx + tailwind-merge)
│   │   │   ├── format.ts  # Data formatting utilities
│   │   │   └── validation.ts
│   │   ├── constants/      # Application constants
│   │   │   ├── colors.ts
│   │   │   ├── routes.ts
│   │   │   └── config.ts
│   │   └── d3/            # D3.js utility functions and configurations
│   │       ├── graph-layout.ts
│   │       ├── force-simulation.ts
│   │       └── svg-utils.ts
│   ├── store/              # State management (Zustand stores)
│   │   ├── auth.ts
│   │   ├── project.ts
│   │   ├── ui.ts
│   │   └── index.ts       # Store combinations and providers
│   ├── types/              # TypeScript type definitions
│   │   ├── api.ts         # API response types
│   │   ├── domain.ts      # Domain model types
│   │   ├── components.ts  # Component prop types
│   │   └── index.ts       # Type exports
│   └── styles/            # Styling and themes
│       ├── globals.css    # Global styles and Tailwind imports
│       ├── components.css # Component-specific styles
│       └── themes/        # Theme configurations
│           ├── light.css
│           └── dark.css
├── public/                # Static assets
│   ├── images/
│   ├── icons/
│   └── favicon.ico
├── tests/                 # Test files
│   ├── __mocks__/        # Jest mocks
│   ├── components/       # Component tests
│   ├── hooks/           # Hook tests
│   ├── utils/           # Utility tests
│   └── e2e/             # Playwright E2E tests
└── docs/                 # Frontend documentation
    ├── components.md
    ├── state-management.md
    └── testing.md
```

## Naming Conventions

### Files and Directories
```typescript
// Components: PascalCase
HealthScoreCard.tsx
DependencyGraph.tsx
ConflictResolver.tsx

// Hooks: camelCase with 'use' prefix
useProjects.ts
useAnalysisData.ts
useWebSocket.ts

// Utilities: camelCase
formatFileSize.ts
validateConfig.ts
apiClient.ts

// Constants: camelCase or SCREAMING_SNAKE_CASE for values
colors.ts
API_ENDPOINTS.ts

// Types: PascalCase
ProjectTypes.ts
AnalysisTypes.ts

// Directories: kebab-case
dependency-analysis/
architecture-validation/
common-components/
```

### Variables and Functions
```typescript
// Variables: camelCase
const projectData = await fetchProjects();
const isLoading = true;
const analysisResults = [];

// Functions: camelCase
const calculateHealthScore = (metrics: HealthMetrics): number => {
  // implementation
};

const handleNodeClick = (node: DependencyNode) => {
  // implementation
};

// Constants: SCREAMING_SNAKE_CASE
const MAX_GRAPH_NODES = 1000;
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL;
const DEFAULT_TIMEOUT = 30000;

// Types and Interfaces: PascalCase
interface ProjectData {
  id: string;
  name: string;
}

type AnalysisStatus = 'idle' | 'running' | 'completed' | 'failed';
```

## Component Development Standards

### Component Structure Template
```typescript
// Import order: external libraries, internal modules, types
import React, { useState, useEffect, useCallback } from 'react';
import { Card, CardHeader, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { useProjects } from '@/hooks/useProjects';
import { cn } from '@/lib/utils/cn';
import type { Project, AnalysisResult } from '@/types/domain';

// Props interface with JSDoc comments
interface HealthScoreCardProps {
  /** Current health score (0-100) */
  score: number;
  /** Previous score for trend calculation */
  previousScore?: number;
  /** Detailed health breakdown */
  breakdown: HealthBreakdown;
  /** Loading state indicator */
  isLoading?: boolean;
  /** CSS class name override */
  className?: string;
  /** Callback when user clicks for details */
  onDetailsClick?: () => void;
}

/**
 * Displays project health score with trend indicators and breakdown
 * 
 * @example
 * ```tsx
 * <HealthScoreCard
 *   score={85}
 *   previousScore={78}
 *   breakdown={healthBreakdown}
 *   onDetailsClick={() => navigate('/health-details')}
 * />
 * ```
 */
export const HealthScoreCard: React.FC<HealthScoreCardProps> = ({
  score,
  previousScore,
  breakdown,
  isLoading = false,
  className,
  onDetailsClick,
}) => {
  // State declarations
  const [isExpanded, setIsExpanded] = useState(false);
  
  // Custom hooks
  const { data: projects } = useProjects();
  
  // Computed values
  const trend = previousScore ? score - previousScore : null;
  const trendColor = trend && trend > 0 ? 'text-green-600' : 'text-red-600';
  const scoreColor = score >= 80 ? 'text-green-600' : score >= 50 ? 'text-yellow-600' : 'text-red-600';
  
  // Event handlers
  const handleDetailsClick = useCallback(() => {
    onDetailsClick?.();
  }, [onDetailsClick]);
  
  const handleToggleExpanded = useCallback(() => {
    setIsExpanded(prev => !prev);
  }, []);
  
  // Effects
  useEffect(() => {
    // Side effects here
  }, [score]);
  
  // Loading state
  if (isLoading) {
    return (
      <Card className={cn("animate-pulse", className)}>
        <CardContent className="p-6">
          <div className="h-8 w-20 bg-gray-300 rounded mb-2" />
          <div className="h-4 w-32 bg-gray-300 rounded" />
        </CardContent>
      </Card>
    );
  }
  
  // Main render
  return (
    <Card 
      className={cn("cursor-pointer hover:shadow-md transition-shadow", className)}
      onClick={handleDetailsClick}
      data-testid="health-score-card"
    >
      <CardHeader>
        <div className="flex items-center justify-between">
          <h3 className="text-lg font-semibold">Health Score</h3>
          {trend !== null && (
            <span className={cn("text-sm font-medium", trendColor)}>
              {trend > 0 ? '↑' : '↓'} {Math.abs(trend)}
            </span>
          )}
        </div>
      </CardHeader>
      <CardContent>
        <div className={cn("text-3xl font-bold", scoreColor)} data-testid="health-score">
          {score}
        </div>
        <p className="text-sm text-gray-600 mt-1">
          Overall monorepo health
        </p>
        
        {isExpanded && (
          <div className="mt-4 space-y-2">
            <HealthBreakdownItem
              label="Dependencies"
              score={breakdown.dependencyHealth}
            />
            <HealthBreakdownItem
              label="Architecture"
              score={breakdown.architectureCompliance}
            />
            <HealthBreakdownItem
              label="Performance"
              score={breakdown.buildPerformance}
            />
          </div>
        )}
        
        <Button
          variant="ghost"
          size="sm"
          onClick={handleToggleExpanded}
          className="mt-2"
        >
          {isExpanded ? 'Show Less' : 'Show Details'}
        </Button>
      </CardContent>
    </Card>
  );
};

// Sub-components (when small and related)
interface HealthBreakdownItemProps {
  label: string;
  score: number;
}

const HealthBreakdownItem: React.FC<HealthBreakdownItemProps> = ({
  label,
  score,
}) => (
  <div className="flex justify-between items-center">
    <span className="text-sm text-gray-600">{label}</span>
    <span className="text-sm font-medium">{score}%</span>
  </div>
);
```

### Custom Hooks Standards
```typescript
// Hook naming: use + descriptive name
export const useAnalysisData = (projectId: string) => {
  // State for the hook
  const [data, setData] = useState<AnalysisResult | null>(null);
  const [error, setError] = useState<Error | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  
  // Memoized fetch function
  const fetchData = useCallback(async () => {
    if (!projectId) return;
    
    setIsLoading(true);
    setError(null);
    
    try {
      const result = await api.getProjectAnalysis(projectId);
      setData(result);
    } catch (err) {
      setError(err instanceof Error ? err : new Error('Unknown error'));
    } finally {
      setIsLoading(false);
    }
  }, [projectId]);
  
  // Effect for initial load
  useEffect(() => {
    fetchData();
  }, [fetchData]);
  
  // Return object with consistent naming
  return {
    data,
    error,
    isLoading,
    refetch: fetchData,
  };
};

// Hook with configuration options
interface UseWebSocketOptions {
  reconnectInterval?: number;
  maxReconnectAttempts?: number;
  onMessage?: (data: any) => void;
  onError?: (error: Event) => void;
}

export const useWebSocket = (url: string, options: UseWebSocketOptions = {}) => {
  const {
    reconnectInterval = 5000,
    maxReconnectAttempts = 3,
    onMessage,
    onError,
  } = options;
  
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [connectionStatus, setConnectionStatus] = useState<'connecting' | 'connected' | 'disconnected'>('disconnected');
  const [lastMessage, setLastMessage] = useState<any>(null);
  
  // Connection logic
  const connect = useCallback(() => {
    const ws = new WebSocket(url);
    
    ws.onopen = () => {
      setConnectionStatus('connected');
      setSocket(ws);
    };
    
    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      setLastMessage(data);
      onMessage?.(data);
    };
    
    ws.onerror = (error) => {
      setConnectionStatus('disconnected');
      onError?.(error);
    };
    
    ws.onclose = () => {
      setConnectionStatus('disconnected');
      setSocket(null);
    };
  }, [url, onMessage, onError]);
  
  // Send message function
  const sendMessage = useCallback((data: any) => {
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify(data));
    }
  }, [socket]);
  
  // Effect for connection management
  useEffect(() => {
    connect();
    
    return () => {
      socket?.close();
    };
  }, [connect]);
  
  return {
    connectionStatus,
    lastMessage,
    sendMessage,
    disconnect: () => socket?.close(),
    reconnect: connect,
  };
};
```

## State Management Standards (Zustand)

### Store Structure
```typescript
// store/auth.ts
import { create } from 'zustand';
import { devtools, persist } from 'zustand/middleware';
import { immer } from 'zustand/middleware/immer';

interface User {
  id: string;
  name: string;
  email: string;
  avatar?: string;
}

interface AuthState {
  // State
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}

interface AuthActions {
  // Actions
  login: (credentials: LoginCredentials) => Promise<void>;
  logout: () => void;
  refreshToken: () => Promise<void>;
  clearError: () => void;
  setUser: (user: User) => void;
}

type AuthStore = AuthState & AuthActions;

export const useAuthStore = create<AuthStore>()(
  devtools(
    persist(
      immer((set, get) => ({
        // Initial state
        user: null,
        token: null,
        isAuthenticated: false,
        isLoading: false,
        error: null,

        // Actions
        login: async (credentials) => {
          set((state) => {
            state.isLoading = true;
            state.error = null;
          });

          try {
            const response = await api.login(credentials);
            
            set((state) => {
              state.user = response.user;
              state.token = response.token;
              state.isAuthenticated = true;
              state.isLoading = false;
            });
          } catch (error) {
            set((state) => {
              state.error = error.message;
              state.isLoading = false;
            });
          }
        },

        logout: () => {
          set((state) => {
            state.user = null;
            state.token = null;
            state.isAuthenticated = false;
          });
        },

        refreshToken: async () => {
          const { token } = get();
          if (!token) return;

          try {
            const response = await api.refreshToken(token);
            
            set((state) => {
              state.token = response.token;
              state.user = response.user;
            });
          } catch (error) {
            // Token refresh failed, logout user
            get().logout();
          }
        },

        clearError: () => {
          set((state) => {
            state.error = null;
          });
        },

        setUser: (user) => {
          set((state) => {
            state.user = user;
          });
        },
      })),
      {
        name: 'auth-storage',
        // Only persist essential data
        partialize: (state) => ({
          user: state.user,
          token: state.token,
          isAuthenticated: state.isAuthenticated,
        }),
      }
    ),
    {
      name: 'auth-store',
    }
  )
);

// Selectors for specific state slices
export const useAuth = () => useAuthStore((state) => ({
  user: state.user,
  isAuthenticated: state.isAuthenticated,
  isLoading: state.isLoading,
}));

export const useAuthActions = () => useAuthStore((state) => ({
  login: state.login,
  logout: state.logout,
  refreshToken: state.refreshToken,
}));
```

## API Integration Standards

### API Client Structure
```typescript
// lib/api/client.ts
import axios, { AxiosInstance, AxiosResponse } from 'axios';
import { useAuthStore } from '@/store/auth';

// Custom error class
export class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public code?: string,
    public details?: any
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

// Response wrapper interface
interface ApiResponse<T> {
  data: T;
  message?: string;
  status: 'success' | 'error';
  timestamp: string;
}

class MonoGuardAPI {
  private client: AxiosInstance;

  constructor(baseURL: string = process.env.NEXT_PUBLIC_API_URL!) {
    this.client = axios.create({
      baseURL,
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    this.setupInterceptors();
  }

  private setupInterceptors(): void {
    // Request interceptor for auth
    this.client.interceptors.request.use(
      (config) => {
        const token = useAuthStore.getState().token;
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => Promise.reject(error)
    );

    // Response interceptor for error handling
    this.client.interceptors.response.use(
      (response: AxiosResponse<ApiResponse<any>>) => response,
      (error) => {
        if (error.response?.status === 401) {
          // Token expired, trigger logout
          useAuthStore.getState().logout();
        }

        const apiError = new ApiError(
          error.response?.data?.message || 'Request failed',
          error.response?.status || 500,
          error.response?.data?.code,
          error.response?.data
        );

        return Promise.reject(apiError);
      }
    );
  }

  // Generic request method with type safety
  private async request<T>(
    method: 'GET' | 'POST' | 'PUT' | 'DELETE',
    endpoint: string,
    data?: any,
    config?: any
  ): Promise<T> {
    const response = await this.client.request<ApiResponse<T>>({
      method,
      url: endpoint,
      data,
      ...config,
    });

    return response.data.data;
  }

  // API methods with strong typing
  async getProjects(): Promise<Project[]> {
    return this.request<Project[]>('GET', '/api/v1/projects');
  }

  async getProject(id: string): Promise<Project> {
    return this.request<Project>('GET', `/api/v1/projects/${id}`);
  }

  async createProject(data: CreateProjectRequest): Promise<Project> {
    return this.request<Project>('POST', '/api/v1/projects', data);
  }

  async analyzeProject(
    projectId: string,
    options: AnalysisOptions = {}
  ): Promise<AnalysisJob> {
    return this.request<AnalysisJob>(
      'POST',
      `/api/v1/projects/${projectId}/analyze`,
      options
    );
  }

  async getAnalysisResult(
    projectId: string,
    analysisId: string
  ): Promise<AnalysisResult> {
    return this.request<AnalysisResult>(
      'GET',
      `/api/v1/projects/${projectId}/analyses/${analysisId}`
    );
  }

  async getDependencyGraph(projectId: string): Promise<DependencyGraphData> {
    return this.request<DependencyGraphData>(
      'GET',
      `/api/v1/projects/${projectId}/dependencies/graph`
    );
  }
}

// Export singleton instance
export const api = new MonoGuardAPI();

// Export types for consumers
export type {
  Project,
  CreateProjectRequest,
  AnalysisOptions,
  AnalysisJob,
  AnalysisResult,
  DependencyGraphData,
};
```

## Styling Standards (Tailwind CSS)

### Utility-First Approach
```typescript
// Prefer utility classes over custom CSS
const HealthScoreCard = () => (
  <div className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
    <h3 className="text-lg font-semibold text-gray-900 mb-2">Health Score</h3>
    <div className="text-3xl font-bold text-green-600">85</div>
    <p className="text-sm text-gray-600 mt-1">Overall monorepo health</p>
  </div>
);

// Use the cn utility for conditional classes
import { cn } from '@/lib/utils/cn';

const Button = ({ variant, size, className, ...props }) => (
  <button
    className={cn(
      // Base styles
      'inline-flex items-center justify-center rounded-md font-medium transition-colors',
      'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2',
      'disabled:opacity-50 disabled:pointer-events-none',
      
      // Variant styles
      {
        'bg-blue-600 text-white hover:bg-blue-700 focus-visible:ring-blue-500': variant === 'primary',
        'bg-gray-100 text-gray-900 hover:bg-gray-200 focus-visible:ring-gray-500': variant === 'secondary',
        'border border-gray-300 bg-white text-gray-700 hover:bg-gray-50': variant === 'outline',
      },
      
      // Size styles
      {
        'h-8 px-3 text-sm': size === 'sm',
        'h-10 px-4': size === 'md',
        'h-12 px-6 text-lg': size === 'lg',
      },
      
      // Custom className override
      className
    )}
    {...props}
  />
);
```

### Design System Integration
```typescript
// tailwind.config.js extension for design system
module.exports = {
  theme: {
    extend: {
      colors: {
        // Brand colors
        brand: {
          50: '#eff6ff',
          500: '#3b82f6',
          900: '#1e3a8a',
        },
        
        // Semantic colors
        success: {
          50: '#f0fdf4',
          500: '#22c55e',
          900: '#14532d',
        },
        warning: {
          50: '#fffbeb',
          500: '#f59e0b',
          900: '#78350f',
        },
        error: {
          50: '#fef2f2',
          500: '#ef4444',
          900: '#7f1d1d',
        },
      },
      
      spacing: {
        '18': '4.5rem',
        '88': '22rem',
      },
      
      animation: {
        'fade-in': 'fadeIn 0.5s ease-in-out',
        'slide-in': 'slideIn 0.3s ease-out',
      },
    },
  },
};

// Custom component styles (when utilities aren't enough)
// styles/components.css
@layer components {
  .graph-node {
    @apply cursor-pointer transition-all duration-200 hover:scale-110;
  }
  
  .graph-node--selected {
    @apply ring-2 ring-blue-500 ring-offset-2;
  }
  
  .data-table {
    @apply w-full border-collapse border-spacing-0;
  }
  
  .data-table th {
    @apply bg-gray-50 px-4 py-2 text-left font-medium text-gray-900 border-b border-gray-200;
  }
  
  .data-table td {
    @apply px-4 py-2 border-b border-gray-100;
  }
}
```

## Testing Standards

### Component Testing (Jest + React Testing Library)
```typescript
// tests/components/HealthScoreCard.test.tsx
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { HealthScoreCard } from '@/components/dashboard/HealthScoreCard';
import type { HealthBreakdown } from '@/types/domain';

// Mock data
const mockBreakdown: HealthBreakdown = {
  dependencyHealth: 85,
  architectureCompliance: 90,
  buildPerformance: 75,
  technicalDebtLevel: 20,
};

describe('HealthScoreCard', () => {
  it('displays health score with correct styling', () => {
    render(
      <HealthScoreCard
        score={85}
        breakdown={mockBreakdown}
      />
    );

    const scoreElement = screen.getByTestId('health-score');
    expect(scoreElement).toHaveTextContent('85');
    expect(scoreElement).toHaveClass('text-green-600'); // Score >= 80 should be green
  });

  it('shows trend indicator when previous score is provided', () => {
    render(
      <HealthScoreCard
        score={85}
        previousScore={78}
        breakdown={mockBreakdown}
      />
    );

    expect(screen.getByText('↑')).toBeInTheDocument();
    expect(screen.getByText('7')).toBeInTheDocument();
  });

  it('calls onDetailsClick when card is clicked', () => {
    const mockOnDetailsClick = jest.fn();
    
    render(
      <HealthScoreCard
        score={85}
        breakdown={mockBreakdown}
        onDetailsClick={mockOnDetailsClick}
      />
    );

    fireEvent.click(screen.getByTestId('health-score-card'));
    expect(mockOnDetailsClick).toHaveBeenCalledTimes(1);
  });

  it('shows loading state correctly', () => {
    render(
      <HealthScoreCard
        score={85}
        breakdown={mockBreakdown}
        isLoading={true}
      />
    );

    expect(screen.getByTestId('health-score-card')).toHaveClass('animate-pulse');
  });

  it('expands to show breakdown details', async () => {
    render(
      <HealthScoreCard
        score={85}
        breakdown={mockBreakdown}
      />
    );

    // Initially breakdown should not be visible
    expect(screen.queryByText('Dependencies')).not.toBeInTheDocument();

    // Click show details
    fireEvent.click(screen.getByText('Show Details'));

    // Wait for expansion
    await waitFor(() => {
      expect(screen.getByText('Dependencies')).toBeInTheDocument();
      expect(screen.getByText('Architecture')).toBeInTheDocument();
      expect(screen.getByText('Performance')).toBeInTheDocument();
    });
  });
});

// Hook testing
import { renderHook, act } from '@testing-library/react';
import { useAnalysisData } from '@/hooks/useAnalysisData';
import { api } from '@/lib/api/client';

// Mock the API
jest.mock('@/lib/api/client');
const mockApi = api as jest.Mocked<typeof api>;

describe('useAnalysisData', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('fetches analysis data on mount', async () => {
    const mockData = { id: '1', status: 'completed' };
    mockApi.getProjectAnalysis.mockResolvedValue(mockData);

    const { result } = renderHook(() => useAnalysisData('project-1'));

    expect(result.current.isLoading).toBe(true);

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
      expect(result.current.data).toEqual(mockData);
      expect(result.current.error).toBe(null);
    });
  });

  it('handles errors gracefully', async () => {
    const mockError = new Error('API Error');
    mockApi.getProjectAnalysis.mockRejectedValue(mockError);

    const { result } = renderHook(() => useAnalysisData('project-1'));

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toEqual(mockError);
      expect(result.current.data).toBe(null);
    });
  });

  it('allows manual refetch', async () => {
    const mockData = { id: '1', status: 'completed' };
    mockApi.getProjectAnalysis.mockResolvedValue(mockData);

    const { result } = renderHook(() => useAnalysisData('project-1'));

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });

    // Clear mock to ensure refetch is called
    mockApi.getProjectAnalysis.mockClear();

    act(() => {
      result.current.refetch();
    });

    expect(mockApi.getProjectAnalysis).toHaveBeenCalledTimes(1);
  });
});
```

### E2E Testing (Playwright)
```typescript
// tests/e2e/dashboard.spec.ts
import { test, expect } from '@playwright/test';

test.describe('Dashboard', () => {
  test.beforeEach(async ({ page }) => {
    // Mock authentication
    await page.goto('/dashboard');
    
    // Wait for initial load
    await page.waitForLoadState('networkidle');
  });

  test('displays health score and allows drilling down', async ({ page }) => {
    // Verify health score is visible
    const healthScore = page.getByTestId('health-score');
    await expect(healthScore).toBeVisible();
    await expect(healthScore).toContainText(/\d+/);

    // Click to expand details
    await page.click('text=Show Details');

    // Verify breakdown appears
    await expect(page.getByText('Dependencies')).toBeVisible();
    await expect(page.getByText('Architecture')).toBeVisible();
    await expect(page.getByText('Performance')).toBeVisible();
  });

  test('navigates to dependency analysis', async ({ page }) => {
    // Click on dependencies navigation
    await page.click('[data-testid="nav-dependencies"]');

    // Wait for dependency page to load
    await page.waitForURL('**/dependencies');
    await page.waitForSelector('[data-testid="dependency-graph"]');

    // Verify graph loads
    const graphNodes = page.locator('.graph-node');
    await expect(graphNodes.first()).toBeVisible();
  });

  test('filters issues by severity', async ({ page }) => {
    // Open filter menu
    await page.click('[data-testid="issue-filter-button"]');

    // Select only high severity
    await page.uncheck('input[value="low"]');
    await page.uncheck('input[value="medium"]');
    await page.check('input[value="high"]');

    // Apply filters
    await page.click('text=Apply Filters');

    // Verify only high severity issues are shown
    const issues = page.locator('[data-testid="issue-item"]');
    for (let i = 0; i < await issues.count(); i++) {
      const issue = issues.nth(i);
      await expect(issue.locator('[data-severity="high"]')).toBeVisible();
    }
  });
});

// Performance testing
test.describe('Performance', () => {
  test('dependency graph loads within performance budget', async ({ page }) => {
    // Navigate to dependencies page
    await page.goto('/dependencies');

    // Measure loading time
    const startTime = Date.now();
    
    await page.waitForSelector('[data-testid="dependency-graph"]', {
      timeout: 5000
    });
    
    const loadTime = Date.now() - startTime;
    
    // Should load within 3 seconds
    expect(loadTime).toBeLessThan(3000);
  });
});
```

## Performance Optimization Standards

### Bundle Optimization
```javascript
// next.config.js
/** @type {import('next').NextConfig} */
const nextConfig = {
  // Enable SWC compiler
  swcMinify: true,
  
  // Optimize images
  images: {
    formats: ['image/webp', 'image/avif'],
    deviceSizes: [640, 750, 828, 1080, 1200, 1920, 2048, 3840],
    imageSizes: [16, 32, 48, 64, 96, 128, 256, 384],
  },
  
  // Bundle analyzer (for development)
  ...(process.env.ANALYZE === 'true' && {
    webpack(config) {
      const { BundleAnalyzerPlugin } = require('webpack-bundle-analyzer');
      config.plugins.push(
        new BundleAnalyzerPlugin({
          analyzerMode: 'server',
          analyzerPort: 8888,
          openAnalyzer: true,
        })
      );
      return config;
    },
  }),
  
  // Experimental features for performance
  experimental: {
    optimizeCss: true,
    scrollRestoration: true,
  },
};
```

### Code Splitting Strategies
```typescript
// Dynamic imports for large components
import dynamic from 'next/dynamic';
import { Suspense } from 'react';

// Lazy load heavy D3.js dependency graph
const DependencyGraph = dynamic(
  () => import('@/components/dependency/DependencyGraph'),
  {
    ssr: false, // Client-side only for D3.js
    loading: () => <GraphSkeleton />,
  }
);

// Lazy load chart components
const TrendChart = dynamic(
  () => import('@/components/dashboard/TrendChart'),
  {
    loading: () => <ChartSkeleton />,
  }
);

// Usage with Suspense for better UX
const DashboardPage = () => (
  <div>
    <HealthScoreCard />
    
    <Suspense fallback={<ChartSkeleton />}>
      <TrendChart />
    </Suspense>
    
    <Suspense fallback={<div>Loading dependency analysis...</div>}>
      <DependencyAnalysisTab />
    </Suspense>
  </div>
);

// Route-based code splitting (automatic with App Router)
// app/dependencies/page.tsx - automatically code split
// app/architecture/page.tsx - automatically code split
```

### Memory Management
```typescript
// Cleanup effects in components
const DependencyGraph = () => {
  const svgRef = useRef<SVGSVGElement>(null);
  const simulationRef = useRef<d3.Simulation<any, any> | null>(null);
  
  useEffect(() => {
    // Setup D3 simulation
    const simulation = d3.forceSimulation(nodes)
      .force('link', d3.forceLink(links))
      .force('charge', d3.forceManyBody())
      .force('center', d3.forceCenter());
    
    simulationRef.current = simulation;
    
    // Cleanup on unmount
    return () => {
      simulation.stop();
      simulationRef.current = null;
    };
  }, [nodes, links]);
  
  // ... component implementation
};

// Debounced search to prevent excessive API calls
import { useDebouncedCallback } from 'use-debounce';

const SearchInput = () => {
  const [query, setQuery] = useState('');
  
  const debouncedSearch = useDebouncedCallback(
    (searchQuery: string) => {
      if (searchQuery.trim()) {
        api.searchProjects(searchQuery);
      }
    },
    500 // 500ms delay
  );
  
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setQuery(value);
    debouncedSearch(value);
  };
  
  return (
    <input
      type="text"
      value={query}
      onChange={handleInputChange}
      placeholder="Search projects..."
    />
  );
};
```

## Accessibility Standards (WCAG 2.1 AA)

### Semantic HTML and ARIA
```typescript
// Proper heading hierarchy and landmarks
const Dashboard = () => (
  <main role="main" aria-labelledby="dashboard-title">
    <header>
      <h1 id="dashboard-title">MonoGuard Dashboard</h1>
    </header>
    
    <section aria-labelledby="health-section-title">
      <h2 id="health-section-title">Project Health</h2>
      <HealthScoreCard />
    </section>
    
    <section aria-labelledby="issues-section-title">
      <h2 id="issues-section-title">Critical Issues</h2>
      <IssuesList />
    </section>
  </main>
);

// Interactive elements with proper ARIA attributes
const FilterButton = ({ isExpanded, onToggle }) => (
  <button
    type="button"
    aria-expanded={isExpanded}
    aria-controls="filter-menu"
    aria-label="Toggle filter options"
    onClick={onToggle}
    className="p-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
  >
    <FilterIcon aria-hidden="true" />
    Filters
  </button>
);

// Complex interactive widgets
const DependencyGraph = () => {
  const [selectedNode, setSelectedNode] = useState<string | null>(null);
  
  return (
    <div
      role="img"
      aria-label="Dependency graph showing relationships between packages"
      tabIndex={0}
      onKeyDown={(e) => {
        // Keyboard navigation for graph
        if (e.key === 'Enter' || e.key === ' ') {
          // Handle node selection
        }
      }}
    >
      <svg>
        {nodes.map((node) => (
          <g
            key={node.id}
            role="button"
            tabIndex={0}
            aria-label={`Package ${node.name}, ${node.issues.length} issues`}
            aria-selected={selectedNode === node.id}
            onClick={() => setSelectedNode(node.id)}
            onKeyDown={(e) => {
              if (e.key === 'Enter' || e.key === ' ') {
                setSelectedNode(node.id);
              }
            }}
          >
            <circle {...nodeAttributes} />
            <text {...textAttributes}>{node.name}</text>
          </g>
        ))}
      </svg>
    </div>
  );
};
```

### Color and Contrast Standards
```css
/* Ensure sufficient color contrast ratios */
:root {
  /* Text colors with 4.5:1 contrast ratio minimum */
  --text-primary: #111827; /* Gray-900 on white background */
  --text-secondary: #6b7280; /* Gray-500 on white background */
  --text-muted: #9ca3af; /* Gray-400 on white background */
  
  /* Interactive element colors */
  --primary-600: #2563eb; /* Blue-600 */
  --primary-700: #1d4ed8; /* Blue-700 for hover states */
  
  /* Status colors with sufficient contrast */
  --success-600: #16a34a;
  --warning-600: #ca8a04;
  --error-600: #dc2626;
}

/* Focus indicators */
.focus-visible:focus-visible {
  outline: 2px solid var(--primary-600);
  outline-offset: 2px;
}

/* Ensure interactive elements meet size requirements */
.interactive-element {
  min-height: 44px; /* Minimum touch target size */
  min-width: 44px;
}
```

## Error Handling and Loading States

### Error Boundaries
```typescript
// components/common/ErrorBoundary.tsx
import React, { Component, ErrorInfo, ReactNode } from 'react';
import { Button } from '@/components/ui/button';
import { AlertCircle, RefreshCw } from 'lucide-react';

interface Props {
  children: ReactNode;
  fallback?: ReactNode;
  onError?: (error: Error, errorInfo: ErrorInfo) => void;
}

interface State {
  hasError: boolean;
  error?: Error;
}

export class ErrorBoundary extends Component<Props, State> {
  public state: State = {
    hasError: false,
  };

  public static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('ErrorBoundary caught an error:', error, errorInfo);
    this.props.onError?.(error, errorInfo);
  }

  private handleRetry = () => {
    this.setState({ hasError: false, error: undefined });
  };

  public render() {
    if (this.state.hasError) {
      if (this.props.fallback) {
        return this.props.fallback;
      }

      return (
        <div className="flex flex-col items-center justify-center p-8 text-center">
          <AlertCircle className="h-12 w-12 text-red-500 mb-4" />
          <h2 className="text-lg font-semibold text-gray-900 mb-2">
            Something went wrong
          </h2>
          <p className="text-sm text-gray-600 mb-4 max-w-md">
            An unexpected error occurred. Please try refreshing the page or contact support if the problem persists.
          </p>
          <Button onClick={this.handleRetry} className="flex items-center gap-2">
            <RefreshCw className="h-4 w-4" />
            Try Again
          </Button>
          {process.env.NODE_ENV === 'development' && this.state.error && (
            <details className="mt-4 text-left">
              <summary className="cursor-pointer text-sm text-gray-500">
                Error details (development only)
              </summary>
              <pre className="mt-2 text-xs bg-gray-100 p-4 rounded overflow-auto">
                {this.state.error.stack}
              </pre>
            </details>
          )}
        </div>
      );
    }

    return this.props.children;
  }
}

// Usage in app layout
const Layout = ({ children }) => (
  <div>
    <Header />
    <ErrorBoundary>
      <main>{children}</main>
    </ErrorBoundary>
  </div>
);
```

### Loading States and Skeletons
```typescript
// Skeleton components for loading states
export const HealthScoreCardSkeleton = () => (
  <div className="bg-white rounded-lg shadow-md p-6 animate-pulse">
    <div className="h-6 w-32 bg-gray-300 rounded mb-4" />
    <div className="h-12 w-16 bg-gray-300 rounded mb-2" />
    <div className="h-4 w-48 bg-gray-300 rounded" />
  </div>
);

export const GraphSkeleton = () => (
  <div className="bg-white rounded-lg border p-8 animate-pulse">
    <div className="flex items-center justify-between mb-6">
      <div className="h-6 w-40 bg-gray-300 rounded" />
      <div className="h-8 w-24 bg-gray-300 rounded" />
    </div>
    <div className="aspect-square bg-gray-300 rounded-lg" />
  </div>
);

// Loading state hook
export const useLoadingState = (isLoading: boolean, minLoadingTime = 500) => {
  const [showLoading, setShowLoading] = useState(false);
  
  useEffect(() => {
    let timeout: NodeJS.Timeout;
    
    if (isLoading) {
      setShowLoading(true);
    } else if (showLoading) {
      // Keep loading state for minimum time to prevent flashing
      timeout = setTimeout(() => setShowLoading(false), minLoadingTime);
    }
    
    return () => {
      if (timeout) clearTimeout(timeout);
    };
  }, [isLoading, showLoading, minLoadingTime]);
  
  return showLoading;
};
```

These comprehensive frontend development standards ensure that the MonoGuard frontend is built with modern best practices, excellent user experience, and maintainable code architecture that scales with the project's growth.