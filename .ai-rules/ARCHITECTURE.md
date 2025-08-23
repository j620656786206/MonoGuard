# MonoGuard Architecture Guide

## System Architecture Overview

### Self-Hosting Monorepo Structure

MonoGuard is built as a monorepo to validate its own architecture principles:

```
mono-guard/
├── backend/           # Go services
│   ├── cmd/api/      # API server
│   ├── cmd/analyzer/ # Analysis engine
│   └── internal/     # Shared backend code
├── frontend/          # Next.js web interface
├── cli/              # Node.js CLI tool
├── shared/           # Cross-language types
├── .monoguard.yml    # Self-validation rules
└── tools/            # Build utilities
```

### High-Level Architecture
```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   CLI Tool      │    │   Web Interface  │    │   Git Hooks     │
│   (Node.js)     │    │   (Next.js)      │    │   (Node.js)     │
└─────────┬───────┘    └─────────┬────────┘    └─────────┬───────┘
          │                      │                       │
          └──────────────────────┼───────────────────────┘
                                 │
                    ┌────────────┴────────────┐
                    │     API Gateway         │
                    │     (Go + Gin)          │
                    └────────────┬────────────┘
                                 │
                    ┌────────────┴────────────┐
                    │   Analysis Engine       │
                    │   (Go)                  │
                    │   - AST Parser          │
                    │   - Dependency Analyzer │
                    │   - Architecture Checker│
                    └────────────┬────────────┘
                                 │
                    ┌────────────┴────────────┐
                    │     Database            │
                    │     (PostgreSQL)        │
                    │     + Redis Cache       │
                    └─────────────────────────┘
```

### Core Technology Stack
- **Backend**: Go 1.21+ with Gin framework + GORM ORM
- **Frontend**: Next.js 14 + TypeScript + Tailwind CSS + Shadcn/ui
- **Database**: PostgreSQL 15+ with Redis for caching
- **Visualization**: Chart.js for metrics, D3.js for dependency graphs
- **DevOps**: Docker + Kubernetes, GitHub Actions for CI/CD

## Component Responsibilities

### Analysis Engine (Go)
**Purpose**: Core analysis logic for dependency and architecture validation

**Responsibilities**:
- **AST Parser**: Parse TypeScript/JavaScript files to build dependency trees
- **Dependency Analyzer**: Detect duplicate dependencies, version conflicts, unused packages
- **Architecture Checker**: Validate layer architecture rules, detect circular dependencies
- **Report Generator**: Create analysis reports in JSON/HTML/Markdown formats

**Key Interfaces**:
```go
type AnalysisEngine interface {
    AnalyzeDependencies(projectPath string) (*DependencyAnalysis, error)
    ValidateArchitecture(config *ArchitectureConfig) (*ArchitectureReport, error)
    GenerateReport(analysis *Analysis, format ReportFormat) (*Report, error)
}
```

### API Gateway (Go + Gin)
**Purpose**: RESTful API layer and orchestration service

**Responsibilities**:
- **Authentication**: OAuth2 integration with GitHub/GitLab/Bitbucket
- **Project Management**: CRUD operations for projects and configurations
- **Analysis Orchestration**: Queue and manage analysis jobs
- **Report Serving**: Serve analysis results and historical data

**Key Endpoints**:
- `POST /api/v1/projects` - Create new project
- `POST /api/v1/projects/{id}/analyze` - Trigger analysis
- `GET /api/v1/projects/{id}/reports` - Retrieve analysis reports
- `PUT /api/v1/projects/{id}/config` - Update project configuration

### Web Interface (Next.js)
**Purpose**: Interactive dashboard and configuration management

**Responsibilities**:
- **Dashboard**: Health score visualization, trend analysis, issue summaries
- **Dependency Explorer**: Interactive dependency graphs with D3.js
- **Architecture Viewer**: Layer architecture visualization and violation reports
- **Configuration Manager**: YAML-based rule configuration interface

**Component Architecture**:
```
src/
├── components/
│   ├── dashboard/
│   │   ├── HealthScoreCard.tsx
│   │   ├── TrendChart.tsx
│   │   └── IssuesSummary.tsx
│   ├── dependency/
│   │   ├── DependencyGraph.tsx (D3.js integration)
│   │   ├── DuplicatesList.tsx
│   │   └── ConflictResolver.tsx
│   └── architecture/
│       ├── LayerDiagram.tsx
│       ├── ViolationsList.tsx
│       └── RuleEditor.tsx
├── pages/
│   ├── dashboard.tsx
│   ├── dependencies.tsx
│   ├── architecture.tsx
│   └── reports.tsx
└── lib/
    ├── api.ts (API client)
    ├── types.ts (TypeScript interfaces)
    └── utils.ts (Helper functions)
```

### CLI Tool (Node.js)
**Purpose**: Command-line interface for local development and CI/CD integration

**Responsibilities**:
- **Local Analysis**: Run analysis on local repositories
- **CI Integration**: Provide exit codes for CI/CD pipelines
- **Report Export**: Generate reports in multiple formats
- **Configuration Validation**: Validate .monoguard.yml files

**Command Structure**:
```bash
# Basic analysis
monoguard analyze [path] [options]

# Generate reports
monoguard report --format=json|html|markdown --output=<file>

# Validate configuration
monoguard config validate [config-file]

# CI/CD integration
monoguard ci --threshold=<score> --fail-on=<severity>
```

## Data Models

### Core Analysis Types
```go
type DependencyAnalysis struct {
    DuplicateDependencies []DuplicateDep    `json:"duplicate_dependencies"`
    VersionConflicts      []VersionConflict `json:"version_conflicts"`
    UnusedDependencies    []UnusedDep       `json:"unused_dependencies"`
    CircularDependencies  []CircularDep     `json:"circular_dependencies"`
    BundleImpact         BundleImpactReport `json:"bundle_impact"`
}

type DuplicateDep struct {
    PackageName       string   `json:"package_name"`
    Versions         []string `json:"versions"`
    AffectedPackages []string `json:"affected_packages"`
    EstimatedWaste   string   `json:"estimated_waste"`
    RiskLevel        string   `json:"risk_level"`
    Recommendation   string   `json:"recommendation"`
    MigrationSteps   []string `json:"migration_steps"`
}
```

### Architecture Configuration Schema
```yaml
# .monoguard.yml
architecture:
  layers:
    - name: 'Application Layer'
      pattern: 'apps/*'
      description: 'Frontend applications, can use shared libraries'
      can_import: ['libs/*']
      cannot_import: ['apps/*']
      
    - name: 'UI Component Library'
      pattern: 'libs/ui/*'
      description: 'Pure UI components, no business logic'
      can_import: ['libs/shared/*']
      cannot_import: ['libs/business/*', 'apps/*']

  rules:
    - name: 'No Circular Dependencies'
      severity: 'error'
      description: 'Prevent circular dependencies between packages'
    - name: 'Layer Architecture Violation'
      severity: 'warning'
      description: 'Enforce predefined layer architecture rules'
```

## Processing Pipelines

### Dependency Analysis Pipeline
1. **Discovery Phase**: Scan workspace for package.json files
2. **Parsing Phase**: Extract dependencies and devDependencies
3. **Resolution Phase**: Build complete dependency tree with versions
4. **Analysis Phase**: Apply detection algorithms for issues
5. **Reporting Phase**: Generate structured reports with recommendations

### Architecture Validation Pipeline
1. **Rule Loading**: Parse and validate .monoguard.yml configuration
2. **Import Analysis**: Use Go AST parser to extract import statements
3. **Pattern Matching**: Apply glob patterns to categorize packages into layers
4. **Violation Detection**: Check imports against layer rules
5. **Circular Detection**: Use DFS with cycle detection for circular dependencies

## State Management Strategy

### Frontend State Management
- Use Zustand for global state management
- Implement optimistic updates for better UX
- Cache API responses with SWR
- Use React Query for server state management

### Backend State Management
- Follow Clean Architecture principles
- Separate concerns: handlers, services, repositories
- Use dependency injection for testability
- Implement proper error handling and logging

## Integration Patterns

### API Communication
- RESTful APIs with JSON payloads
- JWT authentication with refresh tokens
- Rate limiting and request throttling
- Structured error responses with error codes

### Data Flow
- Unidirectional data flow in frontend
- Event-driven architecture for background processing
- Message queues for asynchronous operations
- Real-time updates via WebSocket connections

## Scalability Architecture

### Horizontal Scaling
- **API Layer**: Stateless design enables easy horizontal scaling
- **Analysis Workers**: Queue-based architecture supports multiple workers
- **Database**: Use connection pooling and read replicas
- **Caching**: Implement distributed caching with Redis cluster

### Microservices Considerations
- Keep monolithic initially for faster development
- Plan service boundaries around business domains
- Design APIs with eventual service separation in mind
- Use shared libraries for common functionality