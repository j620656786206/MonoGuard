# MonoGuard Monorepo Tooling Strategy

## Tooling Philosophy

MonoGuard itself is built as a multi-language monorepo, giving us first-hand experience with monorepo challenges and solutions. Our tooling choices reflect practical needs for coordinating Go, TypeScript, and Node.js development.

## Primary Tooling Decisions

### Build Coordination: npm Workspaces

**Choice**: Use npm workspaces as the primary monorepo tool
**Rationale**: 
- Native npm support, no additional tooling required
- Excellent TypeScript/Node.js integration for frontend and CLI
- Lightweight compared to Nx or Lerna
- Works well with Go modules for backend coordination

**Configuration**:
```json
{
  "name": "mono-guard",
  "workspaces": [
    "frontend",
    "cli",
    "shared"
  ],
  "scripts": {
    "build:all": "npm run build --workspaces",
    "test:all": "npm run test --workspaces && cd backend && go test ./...",
    "dev": "concurrently \"npm run dev:frontend\" \"npm run dev:backend\" \"npm run dev:cli\""
  }
}
```

### Language-Specific Tooling

#### Go Backend (Analysis Engine + API)
- **Dependency Management**: Go modules (`go.mod`)
- **Build Tool**: Native `go build` with Make for orchestration
- **Testing**: `go test` with testify for assertions
- **Linting**: golangci-lint for comprehensive code quality

#### TypeScript Frontend + CLI
- **Package Management**: npm workspaces
- **Build Tool**: Next.js for frontend, tsc for CLI
- **Testing**: Jest + Testing Library for frontend, Jest for CLI
- **Type Checking**: TypeScript strict mode with shared tsconfig

#### Shared Components
- **API Contracts**: OpenAPI specs generating TypeScript clients
- **Configuration**: Shared JSON schemas for .monoguard.yml
- **Types**: Shared TypeScript definitions in `shared/types/`

## Why Not Other Tools?

### Nx: Overkill for Our Needs
**Pros**: Sophisticated caching, dependency graphs, code generation
**Cons**: 
- Heavy learning curve and configuration overhead
- Primarily focused on JavaScript/TypeScript ecosystems
- Go integration would be custom and complex
- We need simplicity for rapid MVP development

### Lerna: Legacy Concerns
**Pros**: Mature package management, versioning workflows
**Cons**: 
- Maintenance mode, limited active development
- Focused on publishing multiple packages (not our use case)
- No Go integration
- npm workspaces provide similar functionality natively

### Rush: Enterprise Complexity
**Pros**: Excellent for large organizations, robust policies
**Cons**: 
- Designed for teams of 100+ developers
- Complex configuration and setup process
- Overkill for a focused product team
- Limited Go ecosystem support

## Build and Development Workflow

### Local Development Setup

```bash
# Initial setup
npm install                    # Install all workspace dependencies
cd backend && go mod download  # Install Go dependencies

# Development mode (all services)
npm run dev                    # Starts all services concurrently

# Individual service development
npm run dev:frontend          # Next.js dev server
npm run dev:backend           # Go API with hot reload (air)
npm run dev:cli               # CLI in watch mode
```

### Testing Strategy

```bash
# Run all tests
npm run test:all

# Language-specific testing
npm run test --workspace=frontend  # React component tests
npm run test --workspace=cli       # CLI unit tests
cd backend && go test ./...        # Go backend tests

# Integration testing
npm run test:integration           # Cross-service API tests
```

### Continuous Integration

Our GitHub Actions workflow coordinates all languages:

```yaml
name: CI
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: 18
      - uses: actions/setup-go@v3
        with:
          go-version: 1.21

      # Frontend + CLI testing
      - run: npm install
      - run: npm run test:all
      - run: npm run lint:all

      # Backend testing
      - run: cd backend && go test ./...
      - run: cd backend && golangci-lint run

      # Integration testing
      - run: npm run test:integration

      # Self-analysis (dogfooding)
      - run: npm run monoguard:analyze
```

## Dependency Coordination

### Version Management

**Strategy**: Coordinated releases with semantic versioning
- Frontend, CLI, and Backend versions stay in sync
- API contract changes trigger major version bumps
- Shared dependencies updated across all workspaces simultaneously

### Cross-Language API Contracts

**Approach**: OpenAPI specification drives both Go and TypeScript
```yaml
# shared/api/openapi.yml
openapi: 3.0.0
info:
  title: MonoGuard API
  version: 1.0.0
paths:
  /analyze:
    post:
      summary: Analyze monorepo
      requestBody:
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/AnalysisRequest'
```

**Code Generation**:
- Go: `oapi-codegen` generates server interfaces
- TypeScript: `openapi-typescript` generates client types

## Self-Validation Integration

### MonoGuard Analyzing Itself

Our `.monoguard.yml` configuration:
```yaml
architecture:
  layers:
    - name: 'Backend Services'
      pattern: 'backend/*'
      can_import: ['shared/*']
      cannot_import: ['frontend/*', 'cli/*']
    
    - name: 'Frontend Application'
      pattern: 'frontend/*'  
      can_import: ['shared/*']
      cannot_import: ['backend/*', 'cli/*']
    
    - name: 'CLI Tool'
      pattern: 'cli/*'
      can_import: ['shared/*']
      cannot_import: ['backend/*', 'frontend/*']

  rules:
    - name: 'No circular dependencies'
      severity: 'error'
    - name: 'Layer architecture violations'
      severity: 'warning'
```

### Continuous Self-Analysis

Every commit triggers MonoGuard analysis of our own codebase:
```bash
# In CI pipeline
npm run build:cli
./cli/dist/index.js analyze . --config .monoguard.yml --fail-on-error
```

This ensures we experience our own product's pain points and validate our architectural decisions in real-time.

## Scaling Considerations

### When to Consider Nx Migration

If MonoGuard grows beyond current scope, we would consider Nx for:
- **Team Size**: 10+ developers working simultaneously
- **Repository Size**: 50+ packages with complex interdependencies  
- **Build Performance**: Need for sophisticated caching and task orchestration
- **Code Generation**: Automated scaffolding for consistent package structure

### Migration Path

Our simple npm workspace structure provides a clean migration path to Nx:
1. Install Nx alongside existing setup
2. Gradually adopt Nx executors for specific tasks
3. Migrate to Nx workspace.json when complexity justifies it
4. Retain Go modules independently (Nx doesn't need to manage everything)

## Summary

Our monorepo tooling strategy prioritizes:
1. **Simplicity**: npm workspaces + Go modules over complex orchestration
2. **Multi-language Support**: Native tools for each ecosystem
3. **Self-validation**: Continuous dogfooding of our own product
4. **Pragmatism**: Choose tools based on current needs, not theoretical future requirements

This approach gives us the benefits of monorepo development while maintaining the flexibility to evolve our tooling as MonoGuard grows.