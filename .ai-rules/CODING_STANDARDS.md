# MonoGuard Coding Standards

## Repository Structure

### Project Layout
```
mono-guard/
├── backend/                 # Go backend services
│   ├── cmd/                # Application entry points
│   │   ├── api/           # API server
│   │   └── analyzer/      # Analysis engine
│   ├── internal/          # Private application code
│   │   ├── analysis/     # Analysis engine implementation
│   │   ├── api/          # HTTP handlers and middleware
│   │   ├── config/       # Configuration management
│   │   ├── database/     # Database models and migrations
│   │   └── pkg/          # Shared utilities
│   ├── migrations/       # Database migration files
│   └── Dockerfile        # Container configuration
├── frontend/             # Next.js web interface
│   ├── src/
│   │   ├── components/   # React components
│   │   ├── pages/        # Next.js pages
│   │   ├── lib/          # Utilities and API clients
│   │   └── styles/       # CSS and styling
│   ├── public/           # Static assets
│   └── package.json
├── cli/                  # Node.js CLI tool
│   ├── src/
│   │   ├── commands/     # CLI command implementations
│   │   ├── lib/          # Shared utilities
│   │   └── index.ts      # Entry point
│   └── package.json
├── docs/                 # Documentation
├── docker-compose.yml    # Development environment
└── .github/
    └── workflows/        # CI/CD pipelines
```

## Backend Standards (Go)

### Code Organization Principles
- Follow Clean Architecture principles
- Separate concerns: handlers, services, repositories
- Use dependency injection for testability
- Implement proper error handling and logging

### Package Structure
```go
// internal/analysis/
├── service.go       # Business logic interface and implementation
├── repository.go    # Data access interface and implementation
├── models.go        # Domain models and types
├── parser/          # AST parsing logic
│   ├── typescript.go
│   └── javascript.go
└── detector/        # Issue detection algorithms
    ├── duplicates.go
    ├── conflicts.go
    └── circular.go
```

### Naming Conventions
- **Packages**: lowercase, single word when possible (`analysis`, `config`)
- **Files**: lowercase with underscores for separation (`dependency_analyzer.go`)
- **Types**: PascalCase (`DependencyAnalysis`, `AnalysisEngine`)
- **Functions**: PascalCase for exported, camelCase for internal
- **Variables**: camelCase (`projectPath`, `analysisResult`)
- **Constants**: ALL_CAPS with underscores (`MAX_ANALYSIS_TIME`)

### Interface Design
```go
// Define interfaces for testability and flexibility
type AnalysisEngine interface {
    AnalyzeDependencies(ctx context.Context, projectPath string) (*DependencyAnalysis, error)
    ValidateArchitecture(ctx context.Context, config *ArchitectureConfig) (*ArchitectureReport, error)
}

// Implementation
type analysisEngine struct {
    parser     Parser
    detector   IssueDetector
    reporter   ReportGenerator
    logger     *slog.Logger
}
```

### Error Handling
```go
// Use custom error types for different categories
type AnalysisError struct {
    Type    ErrorType
    Message string
    Cause   error
}

func (e *AnalysisError) Error() string {
    return fmt.Sprintf("analysis error [%s]: %s", e.Type, e.Message)
}

// Wrap errors with context
func (s *analysisEngine) AnalyzeDependencies(ctx context.Context, path string) (*DependencyAnalysis, error) {
    if err := s.validatePath(path); err != nil {
        return nil, &AnalysisError{
            Type:    ValidationError,
            Message: "invalid project path",
            Cause:   err,
        }
    }
    // ... implementation
}
```

### Logging Standards
```go
// Use structured logging with slog
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
}))

// Log with context and structured fields
logger.Info("starting dependency analysis",
    "project_path", projectPath,
    "analysis_id", analysisID,
    "user_id", userID,
)
```

### Testing Standards
```go
// Test file naming: *_test.go
// Test function naming: TestFunctionName_Scenario

func TestAnalysisEngine_AnalyzeDependencies_ValidProject(t *testing.T) {
    // Arrange
    engine := NewAnalysisEngine(mockParser, mockDetector, mockReporter, logger)
    
    // Act
    result, err := engine.AnalyzeDependencies(context.Background(), "/valid/path")
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Len(t, result.DuplicateDependencies, 2)
}

// Use table-driven tests for multiple scenarios
func TestDependencyParser_ParsePackageJSON(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected *PackageJSON
        wantErr  bool
    }{
        {
            name:     "valid package.json",
            input:    `{"name": "test", "dependencies": {"lodash": "^4.17.21"}}`,
            expected: &PackageJSON{Name: "test", Dependencies: map[string]string{"lodash": "^4.17.21"}},
            wantErr:  false,
        },
        // ... more test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := ParsePackageJSON(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.expected, result)
            }
        })
    }
}
```

## Frontend Standards (Next.js + TypeScript)

### Component Organization
- Component-based architecture
- Separate presentation from business logic
- Use custom hooks for shared logic
- Implement proper loading states and error boundaries

### File Structure
```
src/
├── components/
│   ├── ui/                 # Reusable UI components (shadcn/ui)
│   ├── dashboard/          # Dashboard-specific components
│   ├── dependency/         # Dependency analysis components
│   └── common/             # Shared components
├── hooks/                  # Custom React hooks
├── lib/                    # Utilities and API clients
├── pages/                  # Next.js pages
├── store/                  # State management (Zustand)
├── styles/                 # Global styles and Tailwind config
└── types/                  # TypeScript type definitions
```

### Component Standards
```typescript
// Component file naming: PascalCase (DependencyGraph.tsx)
// Use functional components with TypeScript

interface DependencyGraphProps {
  dependencies: Dependency[];
  onNodeClick?: (node: DependencyNode) => void;
  className?: string;
}

export const DependencyGraph: React.FC<DependencyGraphProps> = ({
  dependencies,
  onNodeClick,
  className,
}) => {
  // Use hooks for state and effects
  const [selectedNode, setSelectedNode] = useState<DependencyNode | null>(null);
  const { data: graphData, error, isLoading } = useGraphData(dependencies);

  // Handle loading and error states
  if (isLoading) return <LoadingSpinner />;
  if (error) return <ErrorBoundary error={error} />;

  return (
    <div className={cn("dependency-graph", className)}>
      {/* Component implementation */}
    </div>
  );
};
```

### Custom Hook Standards
```typescript
// Custom hooks file naming: use prefix (useGraphData.ts)

interface UseGraphDataReturn {
  data: GraphData | null;
  error: Error | null;
  isLoading: boolean;
  refetch: () => void;
}

export const useGraphData = (dependencies: Dependency[]): UseGraphDataReturn => {
  const [data, setData] = useState<GraphData | null>(null);
  const [error, setError] = useState<Error | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const processData = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      const processed = await processDependencies(dependencies);
      setData(processed);
    } catch (err) {
      setError(err instanceof Error ? err : new Error('Unknown error'));
    } finally {
      setIsLoading(false);
    }
  }, [dependencies]);

  useEffect(() => {
    processData();
  }, [processData]);

  return { data, error, isLoading, refetch: processData };
};
```

### State Management (Zustand)
```typescript
// Store file naming: *Store.ts

interface ProjectState {
  currentProject: Project | null;
  projects: Project[];
  isLoading: boolean;
  error: string | null;
}

interface ProjectActions {
  setCurrentProject: (project: Project) => void;
  fetchProjects: () => Promise<void>;
  createProject: (data: CreateProjectData) => Promise<void>;
  clearError: () => void;
}

export const useProjectStore = create<ProjectState & ProjectActions>((set, get) => ({
  // State
  currentProject: null,
  projects: [],
  isLoading: false,
  error: null,

  // Actions
  setCurrentProject: (project) => set({ currentProject: project }),
  
  fetchProjects: async () => {
    set({ isLoading: true, error: null });
    try {
      const projects = await api.getProjects();
      set({ projects, isLoading: false });
    } catch (error) {
      set({ error: error.message, isLoading: false });
    }
  },

  createProject: async (data) => {
    set({ isLoading: true, error: null });
    try {
      const project = await api.createProject(data);
      set((state) => ({
        projects: [...state.projects, project],
        currentProject: project,
        isLoading: false,
      }));
    } catch (error) {
      set({ error: error.message, isLoading: false });
    }
  },

  clearError: () => set({ error: null }),
}));
```

### API Client Standards
```typescript
// API client with proper error handling and types

class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public code?: string
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

class ApiClient {
  private baseURL: string;
  private token: string | null = null;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    const headers = new Headers(options.headers);
    
    if (this.token) {
      headers.set('Authorization', `Bearer ${this.token}`);
    }
    
    headers.set('Content-Type', 'application/json');

    const response = await fetch(url, {
      ...options,
      headers,
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new ApiError(
        errorData.message || 'Request failed',
        response.status,
        errorData.code
      );
    }

    return response.json();
  }

  async getProjects(): Promise<Project[]> {
    return this.request<Project[]>('/api/v1/projects');
  }

  async createProject(data: CreateProjectData): Promise<Project> {
    return this.request<Project>('/api/v1/projects', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }
}

export const api = new ApiClient(process.env.NEXT_PUBLIC_API_URL!);
```

## CLI Standards (Node.js + TypeScript)

### Command Structure
```typescript
// Use Commander.js for CLI framework
// Command file naming: *Command.ts

import { Command } from 'commander';
import { AnalysisService } from '../services/AnalysisService';

export class AnalyzeCommand {
  static create(): Command {
    return new Command('analyze')
      .description('Analyze project dependencies and architecture')
      .argument('<path>', 'Path to the project to analyze')
      .option('-f, --format <format>', 'Output format (json|html|markdown)', 'json')
      .option('-o, --output <file>', 'Output file path')
      .option('--config <file>', 'Configuration file path', '.monoguard.yml')
      .action(async (path: string, options) => {
        try {
          await AnalyzeCommand.execute(path, options);
        } catch (error) {
          console.error('Analysis failed:', error.message);
          process.exit(1);
        }
      });
  }

  private static async execute(path: string, options: AnalyzeOptions): Promise<void> {
    const analysisService = new AnalysisService();
    
    // Implement progress indicators for long operations
    const spinner = ora('Analyzing dependencies...').start();
    
    try {
      const result = await analysisService.analyze(path, options);
      spinner.succeed('Analysis completed');
      
      // Provide JSON output for programmatic usage
      if (options.output) {
        await fs.writeFile(options.output, JSON.stringify(result, null, 2));
        console.log(`Results written to ${options.output}`);
      } else {
        console.log(JSON.stringify(result, null, 2));
      }
    } catch (error) {
      spinner.fail('Analysis failed');
      throw error;
    }
  }
}
```

## Configuration Standards

### Shared Configuration
- Use environment variables for configuration
- Implement configuration validation
- Support multiple environments (dev, staging, prod)

### Environment Configuration
```typescript
// config/environment.ts

interface Config {
  api: {
    port: number;
    baseUrl: string;
  };
  database: {
    url: string;
    maxConnections: number;
  };
  redis: {
    url: string;
  };
  auth: {
    jwtSecret: string;
    tokenExpiry: string;
  };
}

const config: Config = {
  api: {
    port: parseInt(process.env.PORT || '8080'),
    baseUrl: process.env.API_BASE_URL || 'http://localhost:8080',
  },
  database: {
    url: process.env.DATABASE_URL || '',
    maxConnections: parseInt(process.env.DB_MAX_CONNECTIONS || '10'),
  },
  redis: {
    url: process.env.REDIS_URL || 'redis://localhost:6379',
  },
  auth: {
    jwtSecret: process.env.JWT_SECRET || '',
    tokenExpiry: process.env.JWT_EXPIRY || '24h',
  },
};

// Validate required configuration
export function validateConfig(config: Config): void {
  if (!config.database.url) {
    throw new Error('DATABASE_URL is required');
  }
  if (!config.auth.jwtSecret) {
    throw new Error('JWT_SECRET is required');
  }
}

export default config;
```

## Quality Standards

### Code Formatting
- **Go**: Use `gofmt` and `goimports`
- **TypeScript**: Use Prettier with consistent configuration
- **Configuration**: EditorConfig for consistent editor settings

### Linting
- **Go**: Use `golangci-lint` with comprehensive rule set
- **TypeScript**: Use ESLint with TypeScript and React rules
- **Commit Messages**: Use conventional commits format

### Documentation
- **Code Comments**: Document public APIs and complex algorithms
- **README**: Maintain up-to-date setup and usage instructions
- **API Documentation**: Use OpenAPI/Swagger for API documentation
- **Architecture Decisions**: Document significant architectural choices