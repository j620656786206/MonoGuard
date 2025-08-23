# Layer Validator Requirements Specification

## Overview

The Layer Validator component is responsible for implementing architecture governance in monorepos through configurable layer-based rules and real-time violation detection. This component addresses User Stories 3, 4, and 5 from the MonoGuard PRD, focusing on architecture rule configuration, violation detection, and circular dependency resolution.

## Architecture Context

### Component Position in System
```
┌─────────────────┐    ┌──────────────────┐
│   CLI Tool      │    │   Web Interface  │
│   (Node.js)     │    │   (Next.js)      │
└─────────┬───────┘    └─────────┬────────┘
          │                      │
          └──────────────────────┼───────────────
                                 │
                    ┌────────────┴────────────┐
                    │     API Gateway         │
                    │     (Go + Gin)          │
                    └────────────┬────────────┘
                                 │
                    ┌────────────┴────────────┐
                    │ >>> LAYER VALIDATOR <<< │ <-- THIS COMPONENT
                    │   (Go)                  │
                    │   - Rule Engine         │
                    │   - AST Parser          │
                    │   - Violation Detector  │
                    │   - Circular Resolver   │
                    └────────────┬────────────┘
                                 │
                    ┌────────────┴────────────┐
                    │     Database            │
                    │     (PostgreSQL)        │
                    └─────────────────────────┘
```

### Technology Stack
- **Backend Engine**: Go 1.21+ with GORM for data persistence
- **Configuration Format**: YAML (.monoguard.yml)
- **AST Parsing**: Go AST package for TypeScript/JavaScript analysis
- **Algorithm**: Depth-First Search (DFS) for circular dependency detection
- **Frontend Integration**: JSON API responses consumed by Next.js + D3.js visualization

## Functional Requirements

### FR-1: Layer Architecture Configuration (User Story 3)

#### FR-1.1: YAML Configuration Schema
**Requirement**: The system SHALL support a declarative YAML configuration format for defining layer architecture rules.

**EARS Format**:
- **WHEN** a user creates a `.monoguard.yml` configuration file
- **THE SYSTEM** SHALL parse and validate the layer architecture definitions
- **WHERE** the configuration follows the specified schema format

**Schema Definition**:
```yaml
architecture:
  layers:
    - name: string          # Human-readable layer name
      pattern: string       # Glob pattern (e.g., "apps/*", "libs/ui/*")
      description: string   # Layer purpose description
      can_import: []string  # Allowed import patterns
      cannot_import: []string # Forbidden import patterns
      
  rules:
    - name: string         # Rule identifier
      severity: enum       # error | warning | info
      description: string  # Rule explanation
      auto_fix: boolean    # Optional automatic fix capability
```

**Acceptance Criteria**:
- Support glob pattern matching with `*` and `**` wildcards
- Validate configuration syntax and provide clear error messages
- Allow layer inheritance through pattern hierarchies
- Support configuration templates for React, Angular, Node.js architectures
- Configuration validation time < 5 seconds for complex rule sets

#### FR-1.2: Configuration Wizard
**Requirement**: The system SHALL provide a guided configuration setup process.

**EARS Format**:
- **WHEN** a new user accesses the configuration setup
- **THE SYSTEM** SHALL guide them through layer definition in < 30 minutes
- **WHERE** the wizard provides pre-built templates and validation

**Implementation Requirements**:
- Detect existing project structure automatically
- Suggest layer patterns based on directory analysis
- Provide real-time preview of rule application
- Include architecture best practice recommendations
- Generate complete `.monoguard.yml` file

#### FR-1.3: Rule Severity Management
**Requirement**: The system SHALL support configurable rule severity levels affecting CI/CD behavior.

**Severity Levels**:
- `error`: Fails CI/CD pipeline (exit code 1)
- `warning`: Logs issue but continues (exit code 0)
- `info`: Informational only, no CI impact

### FR-2: Real-time Architecture Violation Detection (User Story 4)

#### FR-2.1: Import/Export Analysis Engine
**Requirement**: The system SHALL analyze TypeScript/JavaScript import statements to detect architecture violations.

**EARS Format**:
- **WHEN** the analysis engine processes a monorepo
- **THE SYSTEM** SHALL parse all TypeScript/JavaScript files for import/export statements
- **WHERE** the analysis completes in < 3 minutes for 100+ packages

**Technical Implementation**:
```go
type ImportAnalysis struct {
    SourcePackage    string   `json:"source_package"`
    ImportedPackage  string   `json:"imported_package"`
    ImportType       string   `json:"import_type"` // named, default, namespace, dynamic
    ImportPath       string   `json:"import_path"`
    LineNumber       int      `json:"line_number"`
    IsViolation      bool     `json:"is_violation"`
    ViolationReason  string   `json:"violation_reason,omitempty"`
}

type ViolationReport struct {
    PackageName      string           `json:"package_name"`
    LayerName        string           `json:"layer_name"`
    ViolatedRule     string           `json:"violated_rule"`
    Violations       []ImportAnalysis `json:"violations"`
    Severity         string           `json:"severity"`
    Recommendation   string           `json:"recommendation"`
    FixSuggestion    string           `json:"fix_suggestion"`
}
```

#### FR-2.2: Pattern Matching Engine
**Requirement**: The system SHALL accurately match package paths against glob patterns for layer classification.

**Pattern Support**:
- Single wildcard: `libs/ui/*` matches `libs/ui/button`
- Double wildcard: `apps/**` matches `apps/web/src/components`
- Exact match: `libs/shared/utils` matches exactly
- Multiple patterns per layer rule

#### FR-2.3: CI/CD Integration
**Requirement**: The system SHALL integrate with continuous integration pipelines for automated violation detection.

**EARS Format**:
- **WHEN** code is pushed or PR is created
- **THE SYSTEM** SHALL execute architecture validation automatically
- **WHERE** analysis results determine CI pipeline success/failure based on severity levels

**Integration Requirements**:
- Support GitHub Actions, GitLab CI, Jenkins
- Provide configurable exit codes (0 for success, 1 for violations)
- Generate structured output for CI consumption (JSON/JUnit XML)
- Support pre-commit hooks for local development
- Include specific line numbers and fix suggestions in reports

#### FR-2.4: Violation Reporting
**Requirement**: The system SHALL provide detailed violation reports with actionable fix suggestions.

**Report Elements**:
- Specific file path and line number of violation
- Explanation of which rule was violated and why
- Concrete fix suggestion with code examples
- Estimated effort to resolve the violation
- Related violations that might be resolved together

### FR-3: Circular Dependency Visualization & Resolution (User Story 5)

#### FR-3.1: Circular Dependency Detection
**Requirement**: The system SHALL detect circular dependencies using DFS algorithm with cycle detection.

**EARS Format**:
- **WHEN** analyzing package dependencies
- **THE SYSTEM** SHALL identify all circular dependency paths using Depth-First Search
- **WHERE** detection covers both direct and transitive circular dependencies

**Algorithm Implementation**:
```go
type CircularDependency struct {
    CyclePath        []string `json:"cycle_path"`          // A -> B -> C -> A
    CycleLength      int      `json:"cycle_length"`        // Number of packages in cycle
    BreakPoints      []BreakPointSuggestion `json:"break_points"`
    ImpactAnalysis   CycleImpactReport `json:"impact_analysis"`
    ResolutionSteps  []ResolutionStep `json:"resolution_steps"`
}

type BreakPointSuggestion struct {
    PackageName      string   `json:"package_name"`
    ImportToRemove   string   `json:"import_to_remove"`
    AlternativeApproach string `json:"alternative_approach"`
    EstimatedEffort  string   `json:"estimated_effort"`    // hours or story points
    RiskLevel        string   `json:"risk_level"`          // low, medium, high
}
```

#### FR-3.2: Interactive Visualization
**Requirement**: The system SHALL provide interactive directed graph visualization of circular dependencies.

**EARS Format**:
- **WHEN** circular dependencies are detected
- **THE SYSTEM** SHALL display an interactive graph showing cycle paths
- **WHERE** users can explore different resolution options visually

**Visualization Requirements**:
- Use D3.js for interactive directed graphs
- Highlight circular paths in distinct colors
- Show optimal breakpoint suggestions with visual indicators
- Support graph zoom, pan, and node filtering
- Load time < 3 seconds for graphs with 100+ nodes
- Export capabilities (PNG, SVG, PDF formats)

#### FR-3.3: Resolution Planning
**Requirement**: The system SHALL generate step-by-step refactoring plans to resolve circular dependencies.

**Planning Features**:
- Prioritize breakpoints by impact and effort analysis
- Provide concrete refactoring patterns (extract interface, dependency inversion)
- Generate TODO lists with specific action items
- Estimate work effort in hours or story points
- Include code examples for common refactoring patterns

**Refactoring Patterns**:
```markdown
## Pattern: Extract Interface
**Problem**: Direct circular dependency between business logic packages
**Solution**: 
1. Create shared interface package
2. Move common interfaces to shared location
3. Update imports to use interfaces instead of concrete implementations
**Effort**: 2-4 hours
**Risk**: Low (non-breaking change)
```

## Non-Functional Requirements

### NFR-1: Performance Requirements

#### NFR-1.1: Analysis Performance
- Analysis completion time: < 3 minutes for 100+ packages
- Memory usage: < 1GB RAM during analysis
- Concurrent package processing with configurable worker pool size
- Incremental analysis support for subsequent runs

#### NFR-1.2: Real-time Responsiveness
- Configuration validation: < 5 seconds
- Violation report generation: < 30 seconds
- Graph visualization loading: < 3 seconds for 100+ nodes
- API response time: < 300ms for query operations

### NFR-2: Accuracy Requirements

#### NFR-2.1: Detection Accuracy
- Architecture violation detection accuracy: ≥ 90%
- Circular dependency detection accuracy: ≥ 95%
- False positive rate: < 5%
- Pattern matching accuracy: ≥ 98%

#### NFR-2.2: Configuration Validation
- Schema validation accuracy: 100%
- Glob pattern validation: 100%
- Rule conflict detection: ≥ 95%

### NFR-3: Compatibility Requirements

#### NFR-3.1: Language Support
- **Full Support**: TypeScript (.ts, .tsx), JavaScript (.js, .jsx, .mjs)
- **Module Systems**: ES modules, CommonJS
- **Import Syntaxes**: Named imports, default imports, namespace imports, dynamic imports
- **Future Support**: CSS imports, JSON imports

#### NFR-3.2: Package Manager Support
- npm workspaces
- Yarn workspaces (v1 and v2+)
- pnpm workspaces
- Lerna (with npm/yarn backend)

#### NFR-3.3: Monorepo Tool Integration
- Nx workspace support
- Rush.js support
- Bazel (basic support)

### NFR-4: Reliability Requirements

#### NFR-4.1: Error Handling
- Graceful degradation when parsing fails for individual packages
- Clear error messages with actionable suggestions
- Automatic recovery from transient failures
- Progress tracking for long-running operations

#### NFR-4.2: Data Integrity
- Configuration file validation before processing
- Atomic analysis operations (complete success or rollback)
- Consistent state management during concurrent operations

## Interface Specifications

### API Endpoints

#### POST /api/projects/{id}/validate-architecture
```json
{
  "config_path": ".monoguard.yml",
  "target_packages": ["apps/web", "libs/ui"], // optional, analyze all if empty
  "severity_threshold": "warning" // only report issues at this level or higher
}
```

**Response**:
```json
{
  "analysis_id": "uuid",
  "status": "completed",
  "violations": [
    {
      "package_name": "apps/web",
      "layer_name": "Application Layer",
      "violated_rule": "Layer Architecture Violation",
      "severity": "error",
      "violations": [
        {
          "source_package": "apps/web",
          "imported_package": "apps/mobile",
          "import_path": "../mobile/shared-utils",
          "line_number": 15,
          "violation_reason": "Applications cannot import from other applications"
        }
      ],
      "recommendation": "Move shared utilities to libs/shared package",
      "fix_suggestion": "Create libs/shared/utils and move the shared code there"
    }
  ],
  "circular_dependencies": [
    {
      "cycle_path": ["libs/business/orders", "libs/business/payments", "libs/business/orders"],
      "break_points": [
        {
          "package_name": "libs/business/orders",
          "import_to_remove": "libs/business/payments",
          "alternative_approach": "Use dependency injection or events",
          "estimated_effort": "4 hours",
          "risk_level": "medium"
        }
      ]
    }
  ],
  "summary": {
    "total_violations": 5,
    "error_count": 2,
    "warning_count": 3,
    "info_count": 0,
    "health_score": 75
  }
}
```

#### GET /api/projects/{id}/architecture/graph
```json
{
  "nodes": [
    {
      "id": "apps/web",
      "name": "Web Application",
      "layer": "Application Layer",
      "has_violations": true,
      "violation_count": 2
    }
  ],
  "edges": [
    {
      "source": "apps/web",
      "target": "libs/ui",
      "type": "dependency",
      "is_violation": false,
      "is_circular": false
    }
  ],
  "circular_paths": [
    {
      "path": ["libs/business/orders", "libs/business/payments"],
      "severity": "error"
    }
  ]
}
```

### Configuration Schema Validation

The system SHALL validate configuration files against this JSON Schema:

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "architecture": {
      "type": "object",
      "properties": {
        "layers": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "name": {"type": "string", "minLength": 1},
              "pattern": {"type": "string", "pattern": "^[a-zA-Z0-9/*_-]+$"},
              "description": {"type": "string"},
              "can_import": {"type": "array", "items": {"type": "string"}},
              "cannot_import": {"type": "array", "items": {"type": "string"}}
            },
            "required": ["name", "pattern"],
            "additionalProperties": false
          }
        },
        "rules": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "name": {"type": "string", "minLength": 1},
              "severity": {"enum": ["error", "warning", "info"]},
              "description": {"type": "string"},
              "auto_fix": {"type": "boolean", "default": false}
            },
            "required": ["name", "severity"],
            "additionalProperties": false
          }
        }
      },
      "required": ["layers"],
      "additionalProperties": false
    }
  },
  "required": ["architecture"],
  "additionalProperties": false
}
```

## Test Requirements

### Test Coverage Targets
- Unit test coverage: ≥ 90% for core analysis algorithms
- Integration test coverage: ≥ 80% for API endpoints
- End-to-end test coverage: Critical user workflows

### Test Data Requirements
- Small monorepo: 10-20 packages (React application)
- Medium monorepo: 50-100 packages (Enterprise application)
- Large monorepo: 200+ packages (Multi-team organization)
- Edge cases: Circular dependencies, complex layer hierarchies, malformed configurations

### Performance Test Scenarios
- Concurrent analysis requests
- Large dependency graphs (1000+ packages)
- Complex circular dependency chains
- Memory usage under stress conditions

## Acceptance Criteria

### User Story 3: Layer Architecture Rule Configuration
- [ ] Support YAML configuration with schema validation
- [ ] Configuration wizard completes setup in < 30 minutes
- [ ] Import pre-built templates for React, Angular, Node.js
- [ ] Real-time configuration preview and validation
- [ ] Rule conflict detection and resolution suggestions

### User Story 4: Real-time Architecture Violation Detection  
- [ ] Analysis completes in < 3 minutes for 100+ packages
- [ ] Detection accuracy ≥ 90% validated against manual review
- [ ] Support TypeScript/JavaScript ES modules
- [ ] Integration with GitHub Actions, GitLab CI, Jenkins
- [ ] Detailed violation reports with line numbers and fix suggestions
- [ ] Configurable exit codes for CI/CD pipelines

### User Story 5: Circular Dependency Visualization & Resolution
- [ ] DFS algorithm implementation for circular dependency detection
- [ ] Interactive directed graph visualization with D3.js
- [ ] Optimal breakpoint identification and suggestions
- [ ] Step-by-step repair recommendations with code examples  
- [ ] Work effort estimation for each repair option
- [ ] Export capabilities (Markdown, PDF) for refactoring plans

## Implementation Notes

### Go Package Structure
```
internal/layer-validator/
├── analyzer/          # Core analysis engine
│   ├── ast_parser.go  # TypeScript/JavaScript AST parsing
│   ├── pattern_matcher.go # Glob pattern matching
│   └── violation_detector.go # Rule violation detection
├── circular/          # Circular dependency detection
│   ├── detector.go    # DFS implementation
│   ├── resolver.go    # Resolution suggestions
│   └── visualizer.go  # Graph data preparation
├── config/           # Configuration management
│   ├── parser.go     # YAML parsing and validation
│   ├── validator.go  # Schema validation
│   └── templates.go  # Pre-built templates
└── api/             # HTTP handlers
    ├── validate.go   # Architecture validation endpoint
    ├── graph.go      # Dependency graph endpoint
    └── config.go     # Configuration management endpoints
```

### Integration Points
- **Dependency Analyzer**: Reuse package discovery and parsing logic
- **Web Interface**: Consume JSON APIs for visualization
- **CLI Tool**: Direct Go package usage for local analysis
- **Database**: Store analysis results and configuration history

This requirements specification provides comprehensive guidance for implementing the Layer Validator component while ensuring alignment with the broader MonoGuard architecture and user needs.