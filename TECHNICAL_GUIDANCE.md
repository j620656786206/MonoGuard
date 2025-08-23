# MonoGuard Technical Guidance Document

## 1. System Architecture Decisions & Component Breakdown

### 1.1 High-Level Architecture
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

### 1.2 Core Technology Stack
- **Backend**: Go 1.21+ with Gin framework + GORM ORM
- **Frontend**: Next.js 14 + TypeScript + Tailwind CSS + Shadcn/ui
- **Database**: PostgreSQL 15+ with Redis for caching
- **Visualization**: Chart.js for metrics, D3.js for dependency graphs
- **部署**: Zeabur（推薦）、Docker + Kubernetes、GitHub Actions for CI/CD
- **Meta Architecture**: MonoGuard itself built as a monorepo for self-validation

### 1.3 Self-Hosting Monorepo Design

**MonoGuard Project Structure:**
```
mono-guard/
├── apps/              # Applications (Nx workspace pattern)
│   ├── api/          # Go services (API + Analysis Engine)
│   │   ├── cmd/api/      # API server entry point
│   │   ├── cmd/analyzer/ # Analysis engine CLI
│   │   └── internal/     # Shared backend code
│   ├── frontend/     # Next.js web interface
│   │   ├── src/components/
│   │   ├── src/app/
│   │   └── src/lib/
│   └── cli/          # Node.js CLI tool
│       ├── src/commands/
│       └── src/lib/
├── libs/             # Shared libraries (Nx workspace pattern)
│   └── shared-types/ # Cross-language shared definitions
│       ├── src/types/    # TypeScript API contracts
│       └── src/configs/  # Shared configurations
├── .monoguard.yml    # Our own architecture rules
└── tools/            # Development utilities
```

**Self-Validation Benefits:**
1. **Continuous Dogfooding**: Monitor our own architecture health
2. **Real-world Testing**: Validate with complex multi-language monorepo
3. **Performance Metrics**: Track our own build and analysis performance
4. **Feature Validation**: Immediate feedback on new capabilities

### 1.4 Component Responsibilities

#### API Application (`apps/api/` - Go)
- **AST Parser**: Parse TypeScript/JavaScript files to build dependency trees
- **Dependency Analyzer**: Detect duplicate dependencies, version conflicts, unused packages
- **Architecture Checker**: Validate layer architecture rules, detect circular dependencies
- **API Gateway**: OAuth2 integration, project management, analysis orchestration
- **Report Generator**: Create analysis reports in JSON/HTML/Markdown formats

#### Frontend Application (`apps/frontend/` - Next.js)
- **Dashboard**: Health score visualization, trend analysis, issue summaries
- **Dependency Explorer**: Interactive dependency graphs with D3.js
- **Architecture Viewer**: Layer architecture visualization and violation reports
- **Configuration Manager**: YAML-based rule configuration interface

#### CLI Application (`apps/cli/` - Node.js)
- **Local Analysis**: Run analysis on local repositories
- **CI Integration**: Provide exit codes for CI/CD pipelines
- **Report Export**: Generate reports in multiple formats
- **Configuration Validation**: Validate .monoguard.yml files

#### Shared Types Library (`libs/shared-types/` - TypeScript)
- **API Contracts**: Cross-language type definitions for API communication
- **Domain Models**: Shared business logic types
- **Configuration Types**: Type definitions for .monoguard.yml and other configs

## 2. Development Workflow & Task Organization

### 2.1 Development Phases

#### Phase 1: Core Engine (Months 1-2)
**Priority 1: Dependency Analysis Engine**
- Implement package.json parser for npm, yarn, pnpm workspaces
- Build dependency tree resolver with version conflict detection
- Create duplicate dependency identifier with bundle impact estimation
- Develop circular dependency detector using DFS algorithm

**Priority 2: Architecture Validation**
- Implement YAML configuration parser for .monoguard.yml
- Build rule engine for layer architecture validation
- Create AST-based import/export analyzer for TypeScript/JavaScript
- Develop architecture violation detector and reporter

#### Phase 2: Interfaces & Integration (Months 3-4)
**Priority 1: API Development**
- Build RESTful API with Gin framework
- Implement project CRUD operations
- Create analysis job queue with background processing
- Add authentication and authorization layers

**Priority 2: Web Interface**
- Develop responsive dashboard with Next.js
- Create interactive dependency graph with D3.js
- Build configuration management interface
- Implement report viewing and export functionality

#### Phase 3: CLI & CI Integration (Month 5)
- Build CLI tool with Node.js
- Implement CI/CD integration (GitHub Actions, GitLab CI)
- Create pre-commit hooks
- Add configuration validation tools

### 2.2 Task Prioritization Matrix
**High Impact, High Effort**
- Dependency analysis engine
- Interactive dependency visualization
- Architecture rule engine

**High Impact, Low Effort**
- Basic CLI commands
- Simple dashboard metrics
- Configuration file validation

**Low Impact, High Effort**
- Advanced AI suggestions
- Multi-language support
- Complex performance optimizations

## 3. Implementation Guidelines for Core Components

### 3.1 Dependency Analysis Engine

#### Data Structures
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

#### Processing Pipeline
1. **Discovery Phase**: Scan workspace for package.json files
2. **Parsing Phase**: Extract dependencies and devDependencies
3. **Resolution Phase**: Build complete dependency tree with versions
4. **Analysis Phase**: Apply detection algorithms for issues
5. **Reporting Phase**: Generate structured reports with recommendations

#### Performance Considerations
- Process packages concurrently using goroutines
- Implement incremental analysis for large monorepos
- Use caching for external package metadata
- Set timeout limits for analysis operations (5 min max for 1000+ packages)

### 3.2 Architecture Violation Detection

#### Configuration Schema
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

#### Detection Algorithm
1. **Rule Loading**: Parse and validate .monoguard.yml configuration
2. **Import Analysis**: Use Go AST parser to extract import statements
3. **Pattern Matching**: Apply glob patterns to categorize packages into layers
4. **Violation Detection**: Check imports against layer rules
5. **Circular Detection**: Use DFS with cycle detection for circular dependencies

### 3.3 Web Interface Architecture

#### Component Structure
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

#### State Management Strategy
- Use Zustand for global state management
- Implement optimistic updates for better UX
- Cache API responses with SWR
- Use React Query for server state management

### 3.4 CLI Tool Implementation

#### Command Structure
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

#### Implementation Approach
- Use Commander.js for CLI framework
- Implement progress indicators for long operations
- Provide JSON output for programmatic usage
- Support configuration via files and environment variables

## 4. Project Structure & Code Organization

### 4.1 Repository Structure
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

### 4.2 Code Organization Principles

#### Backend (Go)
- Follow Clean Architecture principles
- Separate concerns: handlers, services, repositories
- Use dependency injection for testability
- Implement proper error handling and logging

#### Frontend (Next.js)
- Component-based architecture
- Separate presentation from business logic
- Use custom hooks for shared logic
- Implement proper loading states and error boundaries

#### Shared Configuration
- Use environment variables for configuration
- Implement configuration validation
- Support multiple environments (dev, staging, prod)

## 5. Development Environment & Deployment Strategy

### 5.1 Local Development Setup
```bash
# Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Redis 6+
- Docker & Docker Compose

# Setup commands
git clone <repo>
cd mono-guard
docker-compose up -d postgres redis  # Start dependencies
cd backend && go mod download && go run cmd/api/main.go
cd frontend && npm install && npm run dev
cd cli && npm install && npm run build
```

### 5.2 Development Environment Configuration
```yaml
# docker-compose.dev.yml
version: '3.8'
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: monoguard_dev
      POSTGRES_USER: dev
      POSTGRES_PASSWORD: dev123
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:6
    ports:
      - "6379:6379"

  backend:
    build: ./backend
    environment:
      DATABASE_URL: postgres://dev:dev123@postgres:5432/monoguard_dev
      REDIS_URL: redis://redis:6379
      JWT_SECRET: dev-secret-key
    ports:
      - "8080:8080"
    depends_on:
      - postgres
      - redis
```

### 5.3 Production Deployment Strategy

#### Container Strategy
- Multi-stage Docker builds for optimized image sizes
- Separate containers for API, frontend, and background workers
- Use Alpine Linux base images for security and size
- Implement health checks for all services

#### Kubernetes Deployment
```yaml
# k8s/api-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: monoguard-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: monoguard-api
  template:
    metadata:
      labels:
        app: monoguard-api
    spec:
      containers:
      - name: api
        image: monoguard/api:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: monoguard-secrets
              key: database-url
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

#### Infrastructure Requirements
- **Production**: 3 API instances, 2 worker instances, managed database
- **Staging**: 1 API instance, 1 worker instance, smaller database
- **Database**: PostgreSQL with read replicas for reporting queries
- **Caching**: Redis cluster for session management and query caching
- **Monitoring**: Prometheus + Grafana for metrics, ELK stack for logs

### 5.4 CI/CD Pipeline
```yaml
# .github/workflows/ci.yml
name: CI/CD Pipeline
on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test-backend:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: postgres
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.21'
      - name: Run tests
        run: |
          cd backend
          go mod download
          go test -v -race -coverprofile=coverage.out ./...
          
  test-frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
      - name: Run tests
        run: |
          cd frontend
          npm ci
          npm run test:coverage
          npm run build
```

## 6. Critical Technical Considerations & Risk Mitigation

### 6.1 Performance Requirements & Optimization

#### Analysis Engine Performance
- **Target**: Complete analysis in <5 minutes for 100+ packages
- **Strategy**: 
  - Implement concurrent processing with worker pools
  - Use incremental analysis for repeat runs
  - Cache external package metadata
  - Implement analysis result caching with invalidation

#### Memory Management
- **Challenge**: Large monorepos can consume significant memory
- **Solutions**:
  - Stream processing for large files
  - Implement memory pooling for AST nodes
  - Use memory-mapped files for large dependency graphs
  - Set memory limits and implement graceful degradation

#### Database Performance
- **Indexing Strategy**:
  - Index on project_id, created_at for time-series queries
  - Composite indexes for complex filtering operations
  - Full-text search indexes for package names and descriptions
- **Query Optimization**:
  - Use read replicas for reporting queries
  - Implement query result caching
  - Partition large tables by project or date

### 6.2 Scalability Considerations

#### Horizontal Scaling
- **API Layer**: Stateless design enables easy horizontal scaling
- **Analysis Workers**: Queue-based architecture supports multiple workers
- **Database**: Use connection pooling and read replicas
- **Caching**: Implement distributed caching with Redis cluster

#### Vertical Scaling Limits
- **Single Analysis Job**: May require significant CPU and memory
- **Mitigation**: Implement job splitting for very large monorepos
- **Monitoring**: Track resource usage and implement alerts

### 6.3 Security Considerations

#### Data Protection
- **Source Code Security**: Never persist customer source code
- **In-Memory Processing**: Keep analysis data in memory only
- **Automatic Cleanup**: Clear temporary data after 30 minutes
- **Encryption**: Use AES-256 for sensitive configuration data

#### Authentication & Authorization
- **OAuth2 Integration**: Support GitHub, GitLab, Bitbucket
- **RBAC Implementation**: Owner, Admin, Developer, Viewer roles
- **API Security**: JWT tokens with short expiration, rate limiting
- **Session Management**: Secure session storage with Redis

#### Compliance Preparation
- **GDPR Compliance**: Implement data export and deletion
- **Audit Logging**: Log all user actions and system changes
- **SOC 2 Preparation**: Document security controls and procedures

### 6.4 Risk Mitigation Strategies

#### Technical Risks
| Risk | Probability | Impact | Mitigation |
|------|------------|--------|------------|
| AST parsing complexity too high | Medium | High | Start with common patterns, expand gradually |
| Large monorepo performance issues | High | Medium | Implement incremental analysis + concurrency |
| Different toolchain compatibility | Medium | Medium | Focus on mainstream tools (Nx, Lerna, Rush) |
| Self-hosting monorepo complexity | Medium | Low | Use MonoGuard to monitor its own architecture |
| Competitor fast-follow | Low | High | Build technical moats, focus on user experience |

#### Operational Risks
- **Service Availability**: Implement health checks, auto-scaling, and failover
- **Data Loss**: Automated backups, point-in-time recovery, disaster recovery
- **Security Breaches**: Regular security audits, penetration testing, incident response plans

### 6.5 Quality Assurance Strategy

#### Testing Strategy
- **Unit Tests**: 90% coverage for critical analysis algorithms
- **Integration Tests**: API endpoints and database operations
- **End-to-End Tests**: Critical user workflows with Playwright
- **Performance Tests**: Load testing with realistic monorepo datasets

#### Monitoring & Observability
- **Application Metrics**: Analysis duration, accuracy rates, error rates
- **Infrastructure Metrics**: CPU, memory, disk usage, network latency
- **Business Metrics**: User engagement, analysis completion rates
- **Alerting**: Proactive alerts for performance degradation and errors

#### Error Handling
- **Graceful Degradation**: Partial failures should not block entire analysis
- **Clear Error Messages**: Provide actionable error messages to users
- **Automatic Recovery**: Retry failed operations with exponential backoff
- **Error Tracking**: Use structured logging and error aggregation tools

### 6.6 Maintenance & Evolution Strategy

#### Technical Debt Management
- **Regular Refactoring**: Schedule technical debt reduction sprints
- **Architecture Reviews**: Monthly architecture review sessions
- **Code Quality Gates**: Automated code quality checks in CI/CD
- **Documentation**: Maintain up-to-date technical documentation

#### Feature Evolution
- **Backwards Compatibility**: Maintain API versioning for breaking changes
- **Feature Flags**: Use feature flags for gradual rollouts
- **User Feedback Integration**: Regular user research and feedback incorporation
- **A/B Testing**: Test new features with subset of users before full rollout

This technical guidance provides the essential framework for the MonoGuard development team to begin implementation with clear direction while maintaining flexibility for evolution as requirements become more specific during development.