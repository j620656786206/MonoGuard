---
stepsCompleted: ['step-01-validate-prerequisites', 'step-02-design-epics', 'step-03-create-stories', 'step-04-final-validation']
workflowComplete: true
completedDate: '2026-01-14'
inputDocuments:
  - '_bmad-output/planning-artifacts/prd.md'
  - '_bmad-output/planning-artifacts/architecture.md'
  - '_bmad-output/planning-artifacts/ux-design-specification.md'
totalEpics: 9
totalStories: 66
frCoverage: '48/48'
nfrIntegration: '17/17'
---

# mono-guard - Epic Breakdown

## Overview

This document provides the complete epic and story breakdown for mono-guard, decomposing the requirements from the PRD, UX Design if it exists, and Architecture requirements into implementable stories.

## Requirements Inventory

### Functional Requirements

**Dependency Analysis & Detection (FR1-FR6)**
FR1: Users can analyze monorepo dependency graphs by uploading workspace configuration files
FR2: Users can detect circular dependencies across all packages in a monorepo
FR3: Users can identify duplicate dependencies with version conflicts
FR4: Users can view architecture health score (0-100) calculated from dependency analysis
FR5: Users can analyze npm, yarn, and pnpm workspace structures
FR6: Users can exclude specific packages or patterns from analysis

**Circular Dependency Resolution - Core Differentiator (FR7-FR14)**
FR7: Users can view root cause analysis for each detected circular dependency
FR8: Users can see which import statements create circular dependency paths
FR9: Users can receive fix strategy recommendations for circular dependencies
FR10: Users can view step-by-step fix guides with code location references
FR11: Users can access three fix strategy options: Extract Shared Module, Dependency Injection, Module Boundary Refactoring
FR12: Users can see refactoring complexity scores for each circular dependency
FR13: Users can view impact assessment showing how many packages are affected by each circular dependency
FR14: Users can receive before/after explanations for recommended fixes

**Visualization & Reporting (FR15-FR20)**
FR15: Users can view interactive dependency graphs with D3.js visualization
FR16: Users can see circular dependencies highlighted in red on dependency graphs
FR17: Users can expand and collapse nodes in dependency graphs
FR18: Users can export dependency graphs as PNG or SVG images
FR19: Users can export analysis reports in HTML and JSON formats
FR20: Users can view detailed diagnostic reports for circular dependencies

**CLI Interface (FR21-FR27)**
FR21: Users can analyze dependencies via CLI command (`monoguard analyze`)
FR22: Users can run CI/CD validation checks via CLI (`monoguard check`)
FR23: Users can preview fix suggestions via CLI (`monoguard fix --dry-run`)
FR24: Users can initialize configuration files via CLI (`monoguard init`)
FR25: Users can configure analysis depth and exclusion patterns via CLI options
FR26: Users can receive exit codes indicating analysis results (0 = pass, 1 = fail)
FR27: Users can export analysis results in multiple formats (JSON, HTML, Markdown) via CLI

**Web Interface (FR28-FR33)**
FR28: Users can drag and drop package.json files to initiate analysis in web browser
FR29: Users can upload multiple workspace files to analyze complete monorepo structure
FR30: Users can execute dependency analysis entirely in browser via WASM
FR31: Users can view fix suggestions panel alongside dependency graph in web interface
FR32: Users can download analysis reports from web interface
FR33: Users can access web interface without account creation or authentication

**Privacy & Data Management (FR34-FR39)**
FR34: Users can perform complete analysis without uploading code to remote servers
FR35: Users can store analysis results locally in browser IndexedDB
FR36: Users can store analysis results in local `.monoguard/` directory when using CLI
FR37: Users can execute all core analysis features offline without network connection
FR38: Users can opt-in to anonymous usage analytics
FR39: Users can opt-in to error reporting

**Configuration & Customization (FR40-FR44)**
FR40: Users can configure circular dependency detection rules
FR41: Users can define custom architecture health score thresholds
FR42: Users can configure package exclusion patterns
FR43: Users can set workspace detection patterns
FR44: Users can configure analysis output formats

**WASM API (FR45-FR48)**
FR45: Developers can integrate MonoGuard analysis engine into custom applications via WASM API
FR46: Developers can call `analyze()` function with workspace configuration to get full analysis results
FR47: Developers can call `check()` function for validation-only operations
FR48: Developers can receive typed results (DependencyGraph, CircularDependency, HealthScore) from WASM API

### NonFunctional Requirements

**Performance (NFR1-NFR4)**
NFR1: Analyze 100 packages in < 5 seconds (P95), 1000 packages in < 30 seconds (P95)
NFR2: Dependency graph visualization renders in < 2 seconds for 100 packages; user interactions respond in < 500ms; FCP < 1.5 seconds; Lighthouse Performance > 90
NFR3: Web application bundle size < 500KB (gzipped); WASM module size < 2MB (uncompressed)
NFR4: In-browser WASM analysis uses < 100MB RAM for 1000 packages; CLI uses < 200MB RAM

**Reliability (NFR5-NFR8)**
NFR5: All core analysis features work 100% offline without network dependency
NFR6: Analysis errors do not crash the application; P95 error rate < 0.1% for valid workspace inputs
NFR7: Zero data loss for local storage; analysis results are reproducible (same input = same output)
NFR8: Fix suggestion acceptance rate > 60% (Phase 0), > 80% (Phase 1)

**Security & Privacy (NFR9-NFR12)**
NFR9: Zero code upload to remote servers for core analysis features; all analysis runs locally
NFR10: Browser data stored exclusively in IndexedDB; CLI data stored in `.monoguard/` directory
NFR11: Anonymous usage analytics and error reporting opt-in only (not enabled by default)
NFR12: All npm dependencies scanned for known vulnerabilities; critical vulnerabilities patched within 7 days

**Integration (NFR13-NFR15)**
NFR13: Support npm workspaces, yarn workspaces, and pnpm workspaces (pnpm-workspace.yaml)
NFR14: CLI exit codes follow standard conventions (0 = pass, 1 = fail); CI execution time < 2 minutes for 500 packages
NFR15: Support JSON, HTML, and Markdown export formats with complete analysis results

**Scalability (NFR16-NFR17)**
NFR16: Infrastructure cost remains $0/month (Cloudflare Pages free tier); web app serves 10,000 concurrent users
NFR17: Graceful degradation for large monorepos (> 2000 packages); chunked processing (500 packages per batch)

### Additional Requirements

**From Architecture Document:**

- Nx Monorepo structure with apps/ and packages/ directories
- Hybrid Multi-Repository Approach: TanStack Start (Web UI) + Go WASM (Analysis Engine) + Go CLI with Cobra/Viper
- WASM Dynamic Loading with TypeScript Wrapper (Phase 0); Web Worker upgrade path for Phase 1
- Zustand for state management (< 5KB bundle size)
- Tailwind CSS with JIT Mode for styling
- Sentry for error monitoring with opt-in consent
- Dexie.js for IndexedDB wrapper + Hybrid SVG/Canvas rendering (auto-switch at 500 nodes)
- GitHub Actions CI/CD + Cloudflare Pages deployment
- Vitest for testing with 80%+ coverage target; critical paths > 90% coverage
- Go standard testing with Testify for assertions
- Unified Result type for all WASM functions (data + error structure)
- camelCase for all JSON serialization (Go struct tags and TypeScript)
- TypeScript: camelCase for variables/functions, PascalCase for types/components, UPPER_SNAKE_CASE for constants
- Go: PascalCase for exported, camelCase for unexported, snake_case for file names

**From UX Design Document:**

- Three-mode core experience: Quick Fix Mode (emergency), Health Check Mode (preventive), Explore Mode (understanding)
- Zero-configuration startup: drag-and-drop package.json, no registration required
- 0.5s instant feedback, 3s complete analysis timing requirements
- Progressive disclosure: L1 (health score) → L2 (problem list) → L3 (detailed graph)
- Command Palette (Cmd+K) pattern for power users
- Dark mode support with system preference detection
- Tailwind CSS + Headless UI + Radix UI Primitives for accessible components
- Desktop-first responsive design, tablet-friendly (mobile not optimized)
- Health Score visual gradient: Excellent (85-100 green) → Good (70-84) → Fair (50-69 yellow) → Poor (30-49 orange) → Critical (0-29 red)
- Interactive dependency graph with zoom, hover details, click-to-highlight
- Preview → Apply → Undo flow for all fixes (transparent trust building)
- Animation feedback system: score counting animation, success checkmark, skeleton loading
- Emotional design: surprise through speed, trust through transparency
- Achievement badges system and time savings visualization
- Side panel + main view dual-column layout
- Toast notifications for non-blocking feedback

### FR Coverage Map

| FR   | Epic   | Description                        |
| ---- | ------ | ---------------------------------- |
| FR1  | Epic 2 | Analyze monorepo dependency graphs |
| FR2  | Epic 2 | Detect circular dependencies       |
| FR3  | Epic 2 | Identify duplicate dependencies    |
| FR4  | Epic 2 | View health score (0-100)          |
| FR5  | Epic 2 | Support npm/yarn/pnpm workspaces   |
| FR6  | Epic 2 | Exclude specific packages          |
| FR7  | Epic 3 | Root cause analysis                |
| FR8  | Epic 3 | View problem import statements     |
| FR9  | Epic 3 | Fix strategy recommendations       |
| FR10 | Epic 3 | Step-by-step fix guides            |
| FR11 | Epic 3 | Three fix strategy options         |
| FR12 | Epic 3 | Refactoring complexity scores      |
| FR13 | Epic 3 | Impact assessment                  |
| FR14 | Epic 3 | Before/after explanations          |
| FR15 | Epic 4 | D3.js interactive graphs           |
| FR16 | Epic 4 | Circular dependency highlighting   |
| FR17 | Epic 4 | Node expand/collapse               |
| FR18 | Epic 4 | PNG/SVG export                     |
| FR19 | Epic 4 | HTML/JSON report export            |
| FR20 | Epic 4 | Detailed diagnostic reports        |
| FR21 | Epic 6 | CLI analyze command                |
| FR22 | Epic 6 | CLI check command                  |
| FR23 | Epic 6 | CLI fix --dry-run                  |
| FR24 | Epic 6 | CLI init command                   |
| FR25 | Epic 6 | CLI options configuration          |
| FR26 | Epic 6 | Exit codes                         |
| FR27 | Epic 6 | Multi-format output                |
| FR28 | Epic 5 | Drag-and-drop upload               |
| FR29 | Epic 5 | Multi-file upload                  |
| FR30 | Epic 5 | WASM browser execution             |
| FR31 | Epic 5 | Fix suggestions panel              |
| FR32 | Epic 5 | Download reports                   |
| FR33 | Epic 5 | No registration required           |
| FR34 | Epic 7 | Zero data upload                   |
| FR35 | Epic 7 | IndexedDB storage                  |
| FR36 | Epic 7 | Local directory storage            |
| FR37 | Epic 7 | Offline availability               |
| FR38 | Epic 7 | Opt-in analytics                   |
| FR39 | Epic 7 | Opt-in error reporting             |
| FR40 | Epic 8 | Detection rules configuration      |
| FR41 | Epic 8 | Health score thresholds            |
| FR42 | Epic 8 | Exclusion patterns                 |
| FR43 | Epic 8 | Workspace detection patterns       |
| FR44 | Epic 8 | Output format configuration        |
| FR45 | Epic 9 | WASM API integration               |
| FR46 | Epic 9 | analyze() API                      |
| FR47 | Epic 9 | check() API                        |
| FR48 | Epic 9 | Typed results                      |

## Epic List

### Epic 1: Project Foundation & Infrastructure

**Goal:** Establish complete project architecture enabling development team to start implementing features

**User Outcome:** Developers have a runnable project structure including Web UI, WASM analysis engine, and CLI tool scaffolding

**Requirements Covered:**

- Architecture: Nx Monorepo structure (apps/ + packages/)
- Architecture: TanStack Start project initialization
- Architecture: Go WASM project structure
- Architecture: Go CLI with Cobra/Viper project structure
- Architecture: GitHub Actions + Cloudflare Pages CI/CD
- Architecture: Vitest testing framework setup
- Architecture: Naming conventions and implementation patterns

**Implementation Notes:**

- This epic enables all subsequent epics
- Includes basic CI/CD pipeline with lint, test, build stages
- Sets up shared TypeScript types package
- Establishes code quality standards (ESLint, Prettier, Biome)

---

### Epic 2: Core Dependency Analysis Engine

**Goal:** Users can analyze monorepo dependencies and receive health scores

**User Outcome:** Users can upload/specify a workspace and receive complete dependency graph with architecture health score

**Requirements Covered:** FR1, FR2, FR3, FR4, FR5, FR6

**Key Capabilities:**

- Parse npm/yarn/pnpm workspace configurations
- Build complete dependency graph data structure
- Identify circular dependencies using graph algorithms
- Calculate health score (0-100) based on analysis metrics
- Support package exclusion patterns

**Implementation Notes:**

- Core Go analysis engine compiled to WASM
- TypeScript wrapper for WASM integration
- Unified Result<T> type for all WASM functions
- Performance target: 100 packages < 5s, 1000 packages < 30s

---

### Epic 3: Circular Dependency Resolution Engine (Core Differentiator)

**Goal:** Users receive actionable fix recommendations for circular dependencies

**User Outcome:** Users know not just "what's wrong" but "how to fix it" with step-by-step guidance

**Requirements Covered:** FR7, FR8, FR9, FR10, FR11, FR12, FR13, FR14

**Key Capabilities:**

- Root cause analysis (trace import paths)
- Three fix strategy recommendations:
  1. Extract Shared Module
  2. Dependency Injection
  3. Module Boundary Refactoring
- Step-by-step fix guides with code location references
- Impact assessment and complexity scoring
- Before/after explanations

**Implementation Notes:**

- This is MonoGuard's core differentiator
- "Nx tells you there are problems, MonoGuard tells you how to fix them"
- Fix acceptance rate target: >60% (Phase 0), >80% (Phase 1)

---

### Epic 4: Interactive Visualization & Reporting

**Goal:** Users can visually explore dependency relationships and export reports

**User Outcome:** Users can interactively explore dependency graphs, understand architecture, and export reports

**Requirements Covered:** FR15, FR16, FR17, FR18, FR19, FR20

**Key Capabilities:**

- D3.js interactive dependency graph
- Circular dependencies highlighted in red
- Node expand/collapse functionality
- Zoom, pan, hover details
- Export: PNG, SVG, HTML, JSON formats
- Detailed diagnostic reports

**Implementation Notes:**

- Hybrid SVG/Canvas rendering (auto-switch at 500 nodes)
- SVG for < 500 nodes (better interactivity)
- Canvas for > 500 nodes (better performance)
- UX: Click-to-highlight related paths

---

### Epic 5: Web Interface Experience

**Goal:** Users can analyze dependencies through zero-configuration browser interface

**User Outcome:** Users can drag-and-drop package.json and get complete analysis without installing anything

**Requirements Covered:** FR28, FR29, FR30, FR31, FR32, FR33

**Key Capabilities:**

- Drag-and-drop file upload
- Multi-file upload for complete workspace
- WASM-powered in-browser analysis
- Fix suggestions panel alongside graph
- Report download functionality
- No registration/login required

**UX Requirements Integration:**

- 0.5s instant feedback, 3s complete analysis
- Progressive disclosure (L1/L2/L3)
- Command Palette (Cmd+K)
- Dark Mode with system preference detection
- Side panel + main view dual-column layout
- Toast notifications for non-blocking feedback
- Animation feedback system

**Implementation Notes:**

- TanStack Start with SSG mode
- Zustand for state management
- Tailwind CSS + Headless UI
- Desktop-first, tablet-friendly

---

### Epic 6: CLI Tool Experience

**Goal:** Users can analyze dependencies via command line and integrate with CI/CD

**User Outcome:** Developers can run analysis in terminal, DevOps can integrate with CI/CD pipelines

**Requirements Covered:** FR21, FR22, FR23, FR24, FR25, FR26, FR27

**Key Capabilities:**

- `monoguard analyze` - Full analysis with report
- `monoguard check` - CI/CD validation (exit code 0/1)
- `monoguard fix --dry-run` - Preview fix suggestions
- `monoguard init` - Initialize configuration
- CLI options for depth, exclusions, output format
- JSON/HTML/Markdown output formats

**Implementation Notes:**

- Go CLI with Cobra/Viper
- Shares analysis engine with WASM (same Go code)
- CI execution time < 2 minutes for 500 packages
- Human-readable output by default, JSON for machines

---

### Epic 7: Privacy-First Data Management

**Goal:** Users have complete control over their data

**User Outcome:** Users can use MonoGuard completely offline with all data stored locally

**Requirements Covered:** FR34, FR35, FR36, FR37, FR38, FR39

**Key Capabilities:**

- Zero data upload to remote servers
- IndexedDB browser storage (Dexie.js wrapper)
- `.monoguard/` local directory storage for CLI
- 100% offline availability
- Opt-in anonymous analytics
- Opt-in error reporting (Sentry)

**Implementation Notes:**

- Privacy badge displayed in UI
- Consent banner for telemetry
- Clear data management options
- All analysis runs locally (WASM/Go native)

---

### Epic 8: Configuration & Customization

**Goal:** Users can customize MonoGuard to their project needs

**User Outcome:** Users can define custom rules, thresholds, and patterns for their specific requirements

**Requirements Covered:** FR40, FR41, FR42, FR43, FR44

**Key Capabilities:**

- `.monoguard.json` configuration file
- Custom health score thresholds
- Package exclusion patterns
- Workspace detection patterns
- Output format configuration
- Per-project and global settings

**Implementation Notes:**

- Viper configuration management (CLI)
- Settings stored in Zustand + IndexedDB (Web)
- Configuration schema validation
- `monoguard init` generates default config

---

### Epic 9: Developer API Integration

**Goal:** Developers can integrate MonoGuard analysis engine into their own tools

**User Outcome:** Third-party developers can use MonoGuard's analysis capabilities in custom applications

**Requirements Covered:** FR45, FR46, FR47, FR48

**Key Capabilities:**

- TypeScript type definitions package
- `analyze()` function API
- `check()` function API
- Typed results: DependencyGraph, CircularDependency, HealthScore
- npm package distribution

**Implementation Notes:**

- Published as `@monoguard/wasm` npm package
- Complete TypeScript types
- API documentation with examples
- Versioned API with semver

---

## Epic 1 Stories: Project Foundation & Infrastructure

### Story 1.1: Initialize Nx Monorepo Workspace

As a **developer**,
I want **a properly configured Nx monorepo workspace with apps/ and packages/ directories**,
So that **I have a standardized project structure that supports multiple applications and shared packages**.

**Acceptance Criteria:**

**Given** a fresh project directory
**When** I run the initialization commands
**Then** I have an Nx workspace with:

- `apps/` directory for applications (web, cli)
- `packages/` directory for shared libraries (analysis-engine, types, ui-components)
- `nx.json` with proper workspace configuration
- `package.json` with workspace scripts
- `tsconfig.base.json` for TypeScript path mapping
  **And** running `npx nx graph` displays the project structure

---

### Story 1.2: Setup TanStack Start Web Application

As a **developer**,
I want **a TanStack Start application configured for SSG deployment**,
So that **I can build the web interface with modern React tooling**.

**Acceptance Criteria:**

**Given** the Nx monorepo from Story 1.1
**When** I initialize the TanStack Start application in apps/web
**Then** I have:

- TanStack Start project with SSG configuration
- Vite build setup
- Basic routing structure (`/`, `/analyze`, `/results`)
- Tailwind CSS integrated with JIT mode
- Development server runs at localhost:3000
  **And** running `npx nx build web` produces static HTML output
  **And** the bundle is under 100KB gzipped (without WASM)

---

### Story 1.3: Setup Go WASM Analysis Engine Project

As a **developer**,
I want **a Go project structure configured for WASM compilation**,
So that **I can build the analysis engine that runs in the browser**.

**Acceptance Criteria:**

**Given** the Nx monorepo from Story 1.1
**When** I initialize the Go WASM project in packages/analysis-engine
**Then** I have:

- `go.mod` with module path
- `cmd/wasm/main.go` entry point with basic WASM exports
- `pkg/` directory for internal packages
- `Makefile` with WASM build target (`GOOS=js GOARCH=wasm`)
- `wasm_exec.js` copied from Go installation
  **And** running `make build-wasm` produces `monoguard.wasm` file
  **And** the WASM file can be loaded in browser (basic smoke test)

---

### Story 1.4: Setup Go CLI Project with Cobra

As a **developer**,
I want **a Go CLI project using Cobra for command management**,
So that **I can build command-line tools for dependency analysis**.

**Acceptance Criteria:**

**Given** the Nx monorepo from Story 1.1
**When** I initialize the Go CLI project in apps/cli
**Then** I have:

- `go.mod` with module path
- `main.go` entry point
- `cmd/root.go` with base Cobra command
- Placeholder commands: `analyze`, `check`, `fix`, `init`
- Viper configuration integration
  **And** running `go build` produces executable binary
  **And** running `./monoguard --help` displays available commands

---

### Story 1.5: Setup Shared TypeScript Types Package

As a **developer**,
I want **a shared TypeScript types package**,
So that **I can share type definitions between web app and WASM adapter**.

**Acceptance Criteria:**

**Given** the Nx monorepo
**When** I create the types package in packages/types
**Then** I have:

- TypeScript project with build configuration
- Core type definitions: `DependencyGraph`, `Package`, `CircularDependency`, `HealthScore`
- WASM adapter interface: `MonoGuardAnalyzer`
- Result type: `Result<T>` matching Go structure
- Exports properly configured in package.json
  **And** types can be imported in apps/web
  **And** running `npx nx build types` succeeds

---

### Story 1.6: Configure GitHub Actions CI Pipeline

As a **developer**,
I want **a GitHub Actions CI pipeline that runs on every push and PR**,
So that **code quality is automatically verified before merging**.

**Acceptance Criteria:**

**Given** the complete project structure
**When** I push code to GitHub
**Then** the CI pipeline:

- Runs on push to main/develop and on pull requests
- Installs Node.js 20 and Go 1.21+
- Runs lint (`npx nx affected -t lint`)
- Runs tests (`npx nx affected -t test`)
- Builds WASM (`npx nx build analysis-engine`)
- Builds web app (`npx nx build web`)
- Uploads coverage report
  **And** failed checks block PR merging
  **And** pipeline completes in < 5 minutes

---

### Story 1.7: Configure Cloudflare Pages Deployment

As a **developer**,
I want **automated deployment to Cloudflare Pages on main branch push**,
So that **the web app is automatically deployed with zero infrastructure cost**.

**Acceptance Criteria:**

**Given** successful CI pipeline from Story 1.6
**When** code is pushed to main branch
**Then** the deployment:

- Triggers Cloudflare Pages deployment
- Deploys static files from `dist/apps/web`
- Configures WASM MIME type and CORS headers
- Preview deployments work for PRs
  **And** the deployed site is accessible at configured domain
  **And** WASM files load correctly in production

---

### Story 1.8: Setup Testing Framework and Code Quality

As a **developer**,
I want **testing frameworks configured with code quality tools**,
So that **I can write and run tests with consistent code standards**.

**Acceptance Criteria:**

**Given** the complete project structure
**When** testing and quality tools are configured
**Then** I have:

- Vitest configured for web and types packages
- Go testing with Testify configured
- Coverage thresholds: 80% overall, 90% for critical paths
- Biome (or ESLint + Prettier) for code formatting
- Pre-commit hooks (optional: Husky + lint-staged)
  **And** running `npx nx test web` executes tests
  **And** running `go test ./...` executes Go tests
  **And** coverage reports are generated

---

## Epic 2 Stories: Core Dependency Analysis Engine

### Story 2.1: Implement Workspace Configuration Parser

As a **developer**,
I want **the analysis engine to parse workspace configuration files**,
So that **I can extract package information from any supported monorepo format**.

**Acceptance Criteria:**

**Given** a monorepo with workspace configuration
**When** I provide the workspace files to the parser
**Then** the parser correctly extracts:

- All package names and paths
- Each package's dependencies and devDependencies
- Workspace root configuration
  **And** supports npm workspaces (`package.json` with `workspaces` field)
  **And** supports yarn workspaces (`package.json` with `workspaces` field)
  **And** supports pnpm workspaces (`pnpm-workspace.yaml`)
  **And** returns structured `WorkspaceData` type
  **And** analysis completes in < 1 second for 100 packages

---

### Story 2.2: Build Dependency Graph Data Structure

As a **developer**,
I want **the analysis engine to construct a complete dependency graph**,
So that **I can analyze relationships between all packages in the monorepo**.

**Acceptance Criteria:**

**Given** parsed workspace data from Story 2.1
**When** I build the dependency graph
**Then** the graph contains:

- Node for each package with name, path, and metadata
- Directed edges for each dependency relationship
- Edge types: dependency, devDependency, peerDependency
- Internal vs external dependency classification
  **And** the graph data structure matches `DependencyGraph` type definition
  **And** graph construction completes in < 2 seconds for 100 packages
  **And** memory usage is < 50MB for 1000 packages

---

### Story 2.3: Implement Circular Dependency Detection Algorithm

As a **user**,
I want **to identify all circular dependencies in my monorepo**,
So that **I know which packages have problematic dependency relationships**.

**Acceptance Criteria:**

**Given** a dependency graph from Story 2.2
**When** I run circular dependency detection
**Then** the algorithm:

- Detects all cycles using Tarjan's or similar algorithm
- Returns list of `CircularDependency` objects
- Each cycle includes the complete path (A → B → C → A)
- Identifies the shortest representation of each cycle
- Handles complex multi-package cycles
  **And** detection completes in < 3 seconds for 100 packages with 5 cycles
  **And** correctly identifies nested and overlapping cycles

---

### Story 2.4: Identify Duplicate Dependencies with Version Conflicts

As a **user**,
I want **to see which dependencies have version conflicts across packages**,
So that **I can resolve version mismatches that may cause issues**.

**Acceptance Criteria:**

**Given** parsed workspace data
**When** I analyze for duplicate dependencies
**Then** the analysis identifies:

- Dependencies used by multiple packages
- Version mismatches (e.g., pkg-a uses lodash@4.17.21, pkg-b uses lodash@4.17.19)
- Severity classification (major/minor/patch version differences)
  **And** results include affected packages for each conflict
  **And** analysis completes in < 1 second for 100 packages

---

### Story 2.5: Calculate Architecture Health Score

As a **user**,
I want **to see an overall health score (0-100) for my monorepo architecture**,
So that **I can quickly assess the state of my dependency structure**.

**Acceptance Criteria:**

**Given** complete analysis results (graph, cycles, conflicts)
**When** I calculate the health score
**Then** the score is calculated based on:

- Number of circular dependencies (weighted heavily)
- Number of version conflicts
- Dependency depth (average and max)
- Package coupling metrics
  **And** score is 0-100 where higher is better
  **And** score breakdown shows contribution of each factor
  **And** thresholds: Excellent (85-100), Good (70-84), Fair (50-69), Poor (30-49), Critical (0-29)
  **And** calculation completes in < 100ms

---

### Story 2.6: Implement Package Exclusion Patterns

As a **user**,
I want **to exclude specific packages or patterns from analysis**,
So that **I can focus on relevant parts of my monorepo**.

**Acceptance Criteria:**

**Given** analysis configuration with exclusion patterns
**When** I run analysis with exclusions
**Then** the analysis:

- Supports exact package name exclusion (`packages/legacy`)
- Supports glob patterns (`packages/deprecated-*`)
- Supports regex patterns
- Excludes matched packages from all metrics
- Still shows excluded packages in graph (grayed out)
  **And** exclusions are applied before all calculations
  **And** configuration can be provided via API parameter

---

### Story 2.7: Create TypeScript WASM Adapter

As a **developer**,
I want **a TypeScript wrapper for the WASM analysis engine**,
So that **I can call analysis functions with full type safety from the web app**.

**Acceptance Criteria:**

**Given** the compiled WASM module
**When** I use the TypeScript adapter
**Then** I can:

- Initialize WASM with `analyzer.init()`
- Call `analyzer.analyze(workspaceData)` with typed input
- Receive typed `Result<AnalysisResult>` output
- Handle errors with proper error types
  **And** all functions use the unified `Result<T>` pattern
  **And** JSON serialization uses camelCase (matching Go struct tags)
  **And** adapter includes JSDoc documentation

---

## Epic 3 Stories: Circular Dependency Resolution Engine

### Story 3.1: Implement Root Cause Analysis for Circular Dependencies

As a **user**,
I want **to understand the root cause of each circular dependency**,
So that **I know exactly why the cycle exists and where it originates**.

**Acceptance Criteria:**

**Given** a detected circular dependency
**When** I request root cause analysis
**Then** the analysis provides:

- The originating package (where the cycle likely started)
- The dependency chain that creates the cycle
- Which package relationships are problematic
- Confidence score for root cause identification
  **And** analysis explains in human-readable terms why this is the root cause
  **And** results include `RootCauseAnalysis` type with all details

---

### Story 3.2: Trace Import Statement Paths

As a **user**,
I want **to see exactly which import statements create the circular dependency**,
So that **I know the specific code locations that need to be modified**.

**Acceptance Criteria:**

**Given** a circular dependency between packages
**When** I request import path tracing
**Then** the trace shows:

- File paths containing the problematic imports
- Specific import statements (e.g., `import { foo } from '@pkg/bar'`)
- Line numbers where imports occur
- The complete import chain forming the cycle
  **And** results include `ImportTrace[]` with file, line, and statement
  **And** trace works for both ESM and CommonJS imports

---

### Story 3.3: Generate Fix Strategy Recommendations

As a **user**,
I want **to receive recommended fix strategies for each circular dependency**,
So that **I have actionable options to resolve the problem**.

**Acceptance Criteria:**

**Given** a circular dependency with root cause analysis
**When** I request fix recommendations
**Then** I receive up to 3 strategies:

1. **Extract Shared Module** - Move shared code to new package
2. **Dependency Injection** - Invert the dependency relationship
3. **Module Boundary Refactoring** - Restructure module boundaries
   **And** each strategy includes:

- Suitability score (1-10) based on the specific cycle
- Estimated effort (Low/Medium/High)
- Pros and cons for this specific case
  **And** strategies are ranked by suitability

---

### Story 3.4: Create Step-by-Step Fix Guides

As a **user**,
I want **detailed step-by-step guides for each fix strategy**,
So that **I can follow clear instructions to resolve the circular dependency**.

**Acceptance Criteria:**

**Given** a selected fix strategy
**When** I request the step-by-step guide
**Then** the guide includes:

- Numbered steps (1, 2, 3...)
- Specific file paths to modify
- Code snippets showing before/after
- Commands to run (if applicable)
- Verification steps to confirm fix worked
  **And** guide is specific to the actual packages involved
  **And** guide includes rollback instructions if needed

---

### Story 3.5: Calculate Refactoring Complexity Scores

As a **user**,
I want **to see complexity scores for each circular dependency fix**,
So that **I can prioritize which cycles to fix based on effort required**.

**Acceptance Criteria:**

**Given** a circular dependency
**When** I calculate refactoring complexity
**Then** the score considers:

- Number of files affected
- Number of import statements to change
- Depth of the dependency chain
- Number of packages involved
- Presence of external dependencies in the cycle
  **And** score is 1-10 where 1 is simple, 10 is complex
  **And** score includes breakdown of contributing factors
  **And** estimated time range provided (e.g., "15-30 minutes")

---

### Story 3.6: Generate Impact Assessment

As a **user**,
I want **to see how many packages are affected by each circular dependency**,
So that **I can understand the blast radius and prioritize high-impact fixes**.

**Acceptance Criteria:**

**Given** a circular dependency
**When** I request impact assessment
**Then** the assessment shows:

- Direct participants (packages in the cycle)
- Indirect dependents (packages that depend on cycle participants)
- Total affected package count
- Percentage of monorepo affected
- Risk level (Critical/High/Medium/Low)
  **And** visual representation of affected packages available
  **And** includes "ripple effect" analysis

---

### Story 3.7: Provide Before/After Fix Explanations

As a **user**,
I want **to see clear before/after comparisons for each fix**,
So that **I understand exactly what will change and can build confidence in the fix**.

**Acceptance Criteria:**

**Given** a fix recommendation
**When** I request before/after explanation
**Then** I see:

- Current state diagram (with cycle highlighted)
- Proposed state diagram (cycle resolved)
- Diff of package.json changes
- Diff of import statement changes
- Explanation of why this resolves the issue
  **And** explanation is in plain language (non-technical friendly)
  **And** includes warning of any potential side effects

---

### Story 3.8: Integrate Fix Suggestions with Analysis Results

As a **user**,
I want **fix suggestions integrated into the main analysis results**,
So that **I can see problems and solutions together**.

**Acceptance Criteria:**

**Given** complete analysis with circular dependencies
**When** I view analysis results
**Then** each circular dependency includes:

- Quick fix recommendation (best strategy)
- "View all strategies" option
- One-click access to step-by-step guide
- Complexity and impact scores
  **And** results are sorted by impact × ease (quick wins first)
  **And** total estimated fix time displayed

---

## Epic 4 Stories: Interactive Visualization & Reporting

### Story 4.1: Implement D3.js Force-Directed Dependency Graph

As a **user**,
I want **to view my dependency relationships as an interactive force-directed graph**,
So that **I can visually understand the structure of my monorepo**.

**Acceptance Criteria:**

**Given** analysis results with dependency graph data
**When** I view the visualization
**Then** I see:

- Force-directed layout with nodes for each package
- Directed edges showing dependency relationships
- Smooth physics-based animation
- Auto-layout that separates clusters
  **And** graph renders in < 2 seconds for 100 packages
  **And** graph is responsive to container size
  **And** initial layout stabilizes within 3 seconds

---

### Story 4.2: Highlight Circular Dependencies in Graph

As a **user**,
I want **circular dependencies to be visually highlighted in the graph**,
So that **I can immediately identify problematic relationships**.

**Acceptance Criteria:**

**Given** a dependency graph with circular dependencies
**When** the graph renders
**Then** circular dependencies are highlighted:

- Nodes in cycles have red border/glow
- Edges forming cycles are red and thicker
- Cycle paths are animated (pulsing or flowing)
- Legend explains the color coding
  **And** clicking a cycle highlights only that cycle's path
  **And** non-cycle elements are dimmed when cycle is selected

---

### Story 4.3: Implement Node Expand/Collapse Functionality

As a **user**,
I want **to expand and collapse nodes in the dependency graph**,
So that **I can focus on specific areas without visual clutter**.

**Acceptance Criteria:**

**Given** a dependency graph with many packages
**When** I interact with nodes
**Then** I can:

- Double-click to collapse a node (hide its dependencies)
- Double-click again to expand
- Collapse/expand all at a certain depth
- See collapsed node count indicator
  **And** expand/collapse animations are smooth (< 300ms)
  **And** graph re-layouts gracefully after changes

---

### Story 4.4: Add Zoom, Pan, and Navigation Controls

As a **user**,
I want **to zoom and pan the dependency graph**,
So that **I can navigate large graphs effectively**.

**Acceptance Criteria:**

**Given** a dependency graph
**When** I interact with the viewport
**Then** I can:

- Scroll to zoom in/out (with smooth animation)
- Click and drag to pan
- Use zoom controls (+/- buttons)
- Click "Fit to screen" to see entire graph
- Minimap for large graphs (> 50 nodes)
  **And** zoom range is 10% to 400%
  **And** current zoom level is displayed

---

### Story 4.5: Implement Hover Details and Tooltips

As a **user**,
I want **to see package details when hovering over nodes**,
So that **I can quickly understand each package without clicking**.

**Acceptance Criteria:**

**Given** a dependency graph
**When** I hover over a node
**Then** I see a tooltip with:

- Package name and path
- Dependency count (in/out)
- Health contribution score
- Circular dependency involvement (if any)
  **And** tooltip appears within 200ms
  **And** tooltip follows mouse or anchors to node
  **And** hovering highlights connected edges

---

### Story 4.6: Export Graph as PNG/SVG Images

As a **user**,
I want **to export the dependency graph as an image**,
So that **I can include it in documentation or presentations**.

**Acceptance Criteria:**

**Given** a rendered dependency graph
**When** I click export
**Then** I can choose:

- PNG format (raster, with resolution options)
- SVG format (vector, scalable)
- Current view or full graph
- With or without legend
  **And** exported image includes MonoGuard watermark (optional)
  **And** file downloads immediately
  **And** filename includes project name and date

---

### Story 4.7: Export Analysis Reports in Multiple Formats

As a **user**,
I want **to export complete analysis reports**,
So that **I can share findings with my team or archive results**.

**Acceptance Criteria:**

**Given** complete analysis results
**When** I export a report
**Then** I can choose format:

- **JSON** - Machine-readable, complete data
- **HTML** - Standalone viewable report with embedded styles
- **Markdown** - For documentation/wikis
  **And** reports include:
- Health score summary
- Circular dependency list
- Version conflicts
- Fix recommendations
  **And** HTML report is self-contained (no external dependencies)

---

### Story 4.8: Create Detailed Diagnostic Reports

As a **user**,
I want **detailed diagnostic reports for each circular dependency**,
So that **I can deep-dive into specific issues**.

**Acceptance Criteria:**

**Given** a circular dependency
**When** I request a diagnostic report
**Then** the report includes:

- Executive summary (1-2 sentences)
- Complete cycle path visualization
- Root cause analysis details
- All fix strategies with full guides
- Impact assessment
- Related cycles (if any)
  **And** report can be exported as PDF-ready HTML
  **And** report includes timestamps and version info

---

### Story 4.9: Implement Hybrid SVG/Canvas Rendering

As a **developer**,
I want **the graph to automatically switch between SVG and Canvas rendering**,
So that **large graphs remain performant while small graphs stay interactive**.

**Acceptance Criteria:**

**Given** dependency graph data
**When** graph is rendered
**Then** rendering mode is selected:

- SVG for < 500 nodes (better interactivity)
- Canvas for >= 500 nodes (better performance)
- User can override in settings
  **And** mode indicator shows current renderer
  **And** Canvas mode maintains hover/click functionality
  **And** transition between modes preserves viewport state

---

## Epic 5 Stories: Web Interface Experience

### Story 5.1: Create Landing Page with Drag-Drop Zone

As a **user**,
I want **a welcoming landing page with a prominent drag-drop zone**,
So that **I can immediately start analyzing my monorepo without any setup**.

**Acceptance Criteria:**

**Given** I visit the MonoGuard web app
**When** the landing page loads
**Then** I see:

- Clear value proposition headline
- Prominent drag-drop zone (visually distinct)
- "Or click to select files" option
- Privacy badge ("100% local analysis")
- Sample demo button ("Try with example")
  **And** page loads in < 1.5 seconds (FCP)
  **And** drag-drop zone has hover/active states
  **And** no login/registration required

---

### Story 5.2: Implement File Upload and Workspace Detection

As a **user**,
I want **to upload my workspace files and have the app detect the structure**,
So that **I don't need to configure anything manually**.

**Acceptance Criteria:**

**Given** I drag files onto the upload zone
**When** files are dropped
**Then** the app:

- Accepts package.json files
- Accepts pnpm-workspace.yaml files
- Accepts folder uploads (with nested package.json)
- Auto-detects workspace type (npm/yarn/pnpm)
- Shows upload progress for large uploads
  **And** invalid files show clear error message
  **And** upload starts analysis automatically
  **And** supports multiple file selection

---

### Story 5.3: Implement WASM-Powered In-Browser Analysis

As a **user**,
I want **analysis to run entirely in my browser**,
So that **my code never leaves my machine**.

**Acceptance Criteria:**

**Given** uploaded workspace files
**When** analysis runs
**Then**:

- WASM module loads (with progress indicator)
- Analysis executes in browser
- No network requests for analysis (verifiable in DevTools)
- Progress shows "Analyzing X/Y packages"
- Results appear progressively
  **And** first feedback appears in < 0.5 seconds
  **And** complete analysis finishes in < 3 seconds (100 packages)
  **And** privacy indicator confirms "Analyzed locally"

---

### Story 5.4: Build Analysis Results Dashboard

As a **user**,
I want **a clear dashboard showing my analysis results**,
So that **I can quickly understand my monorepo's health**.

**Acceptance Criteria:**

**Given** completed analysis
**When** I view results
**Then** I see (Progressive Disclosure L1/L2/L3):

- **L1**: Health score (large number), problem counts (red/yellow/green)
- **L2**: Problem categories with counts (circular deps, version conflicts)
- **L3**: Detailed graph and fix recommendations
  **And** health score has color coding and trend arrow
  **And** clicking L1 reveals L2, clicking L2 reveals L3
  **And** skeleton loading during transitions

---

### Story 5.5: Create Fix Suggestions Side Panel

As a **user**,
I want **a fix suggestions panel alongside the dependency graph**,
So that **I can see problems and solutions together**.

**Acceptance Criteria:**

**Given** analysis results with issues
**When** I view the results page
**Then** I see:

- Main area: Dependency graph visualization
- Side panel: Fix suggestions list
- Panel is collapsible/resizable
- Clicking a suggestion highlights related nodes in graph
- Each suggestion shows: title, severity, quick action
  **And** panel width is remembered (persisted)
  **And** dual-column layout works on tablet (responsive)

---

### Story 5.6: Implement Report Download Functionality

As a **user**,
I want **to download analysis reports from the web interface**,
So that **I can share or archive my results**.

**Acceptance Criteria:**

**Given** analysis results
**When** I click download/export
**Then** I can:

- Select format (JSON, HTML, Markdown)
- Choose content (full report, summary only)
- Include/exclude graph image
- Download starts immediately
  **And** downloaded file has descriptive name
  **And** HTML report opens correctly in browser
  **And** export works offline (no network needed)

---

### Story 5.7: Implement Dark Mode with System Detection

As a **user**,
I want **the app to support dark mode**,
So that **I can use it comfortably in low-light environments**.

**Acceptance Criteria:**

**Given** the web app
**When** I use it in dark mode
**Then**:

- App detects system preference automatically
- Manual toggle available in settings
- All UI elements have proper dark variants
- Graph visualization adapts colors for dark mode
- Preference is persisted
  **And** transition between modes is smooth (< 200ms)
  **And** no flash of wrong theme on load

---

### Story 5.8: Implement Command Palette (Cmd+K)

As a **power user**,
I want **a command palette for quick navigation and actions**,
So that **I can work efficiently with keyboard shortcuts**.

**Acceptance Criteria:**

**Given** I'm using the web app
**When** I press Cmd+K (or Ctrl+K)
**Then** I see:

- Modal with search input
- Recent actions list
- Available commands (Analyze, Export, Settings, etc.)
- Fuzzy search filtering
- Keyboard navigation (arrow keys + enter)
  **And** results update as I type (< 50ms)
  **And** pressing Escape closes the palette
  **And** commands execute immediately on selection

---

### Story 5.9: Setup Zustand State Management

As a **developer**,
I want **Zustand configured for application state management**,
So that **state is predictable and debuggable**.

**Acceptance Criteria:**

**Given** the web application
**When** state management is configured
**Then**:

- Analysis store manages: results, loading, errors
- Settings store manages: theme, preferences
- UI store manages: panel states, selections
- DevTools integration works in development
- State persists appropriately (settings to localStorage)
  **And** stores follow naming conventions
  **And** state updates are < 16ms (60fps)

---

### Story 5.10: Implement Toast Notifications

As a **user**,
I want **non-blocking notifications for actions and events**,
So that **I stay informed without interrupting my workflow**.

**Acceptance Criteria:**

**Given** various user actions
**When** actions complete or errors occur
**Then** toast notifications:

- Appear in bottom-right corner
- Auto-dismiss after 5 seconds (configurable)
- Can be manually dismissed
- Stack if multiple appear
- Have type variants (success, error, warning, info)
  **And** success toasts are green with checkmark
  **And** error toasts are red with X icon
  **And** toasts don't block interaction

---

## Epic 6 Stories: CLI Tool Experience

### Story 6.1: Implement `monoguard analyze` Command

As a **developer**,
I want **to run full dependency analysis from the command line**,
So that **I can analyze my monorepo without using the web interface**.

**Acceptance Criteria:**

**Given** I'm in a monorepo directory
**When** I run `monoguard analyze`
**Then** the command:

- Auto-detects workspace configuration
- Runs complete dependency analysis
- Displays health score prominently
- Shows circular dependency count
- Shows version conflict count
- Outputs human-readable summary
  **And** analysis completes in < 5 seconds for 100 packages
  **And** colorized output (when terminal supports it)
  **And** works without any configuration file

---

### Story 6.2: Implement `monoguard check` Command for CI/CD

As a **DevOps engineer**,
I want **a validation command that returns proper exit codes**,
So that **I can integrate MonoGuard into CI/CD pipelines**.

**Acceptance Criteria:**

**Given** I'm running in a CI/CD environment
**When** I run `monoguard check`
**Then** the command:

- Returns exit code 0 if no issues (or below threshold)
- Returns exit code 1 if issues exceed threshold
- Supports `--max-circular` flag (default: 0)
- Supports `--min-health` flag (default: 70)
- Outputs machine-readable summary
  **And** execution time < 2 minutes for 500 packages
  **And** works in non-TTY environments
  **And** supports `--quiet` flag for minimal output

---

### Story 6.3: Implement `monoguard fix --dry-run` Command

As a **developer**,
I want **to preview fix suggestions from the command line**,
So that **I can understand recommended changes before applying them**.

**Acceptance Criteria:**

**Given** a monorepo with circular dependencies
**When** I run `monoguard fix --dry-run`
**Then** I see:

- List of detected issues
- Recommended fix for each issue
- Files that would be modified
- Diff preview of changes
- Estimated effort for each fix
  **And** no files are actually modified
  **And** output can be saved to file with `--output`
  **And** supports `--json` flag for structured output

---

### Story 6.4: Implement `monoguard init` Command

As a **developer**,
I want **to initialize MonoGuard configuration for my project**,
So that **I can customize analysis settings**.

**Acceptance Criteria:**

**Given** I'm in a monorepo directory
**When** I run `monoguard init`
**Then** the command:

- Creates `.monoguard.json` configuration file
- Prompts for basic settings (or uses defaults with `--defaults`)
- Detects workspace type and sets appropriate options
- Creates `.monoguard/` directory for local data
- Adds `.monoguard/` to `.gitignore` (with confirmation)
  **And** doesn't overwrite existing config (unless `--force`)
  **And** validates created configuration

---

### Story 6.5: Implement CLI Options for Analysis Configuration

As a **developer**,
I want **to configure analysis via command-line options**,
So that **I can customize behavior without editing config files**.

**Acceptance Criteria:**

**Given** any MonoGuard command
**When** I pass configuration options
**Then** these options are supported:

- `--exclude <pattern>` - Exclude packages matching pattern
- `--depth <n>` - Limit analysis depth
- `--workspace <path>` - Specify workspace root
- `--config <path>` - Use specific config file
- `--no-color` - Disable colored output
- `--verbose` / `-v` - Increase output verbosity
  **And** CLI options override config file settings
  **And** invalid options show helpful error messages

---

### Story 6.6: Implement Multi-Format Output Export

As a **developer**,
I want **to export analysis results in different formats**,
So that **I can integrate with other tools or generate reports**.

**Acceptance Criteria:**

**Given** analysis results
**When** I run `monoguard analyze --format <format>`
**Then** I can output in:

- `json` - Complete structured data
- `html` - Standalone HTML report
- `markdown` - Markdown-formatted report
- `text` - Human-readable plain text (default)
  **And** `--output <file>` saves to file instead of stdout
  **And** JSON output matches TypeScript type definitions
  **And** HTML report is self-contained

---

### Story 6.7: Implement Progress and Verbose Output

As a **developer**,
I want **to see detailed progress during analysis**,
So that **I know what's happening for large monorepos**.

**Acceptance Criteria:**

**Given** a large monorepo analysis
**When** I run with `-v` or `--verbose`
**Then** I see:

- Current phase (parsing, analyzing, calculating)
- Package count progress (Analyzing 52/100...)
- Time elapsed
- Memory usage (optional, with `-vv`)
  **And** progress updates don't flood the terminal
  **And** spinner animation for long operations
  **And** final summary includes timing breakdown

---

## Epic 7 Stories: Privacy-First Data Management

### Story 7.1: Implement WASM Local Analysis (Zero Data Upload)

As a **user**,
I want **all analysis to run entirely in my browser without any data upload**,
So that **my code never leaves my machine and remains completely private**.

**Acceptance Criteria:**

**Given** I upload workspace files for analysis
**When** the analysis executes
**Then** no network requests are made to external servers (verifiable in DevTools Network tab)
**And** all analysis computation runs in the WASM module locally
**And** a privacy badge displays "Analyzed locally - no data uploaded"
**And** the analysis completes successfully even with network disabled

---

### Story 7.2: Implement IndexedDB Browser Storage

As a **user**,
I want **my analysis results stored locally in my browser**,
So that **I can access historical results without any cloud storage**.

**Acceptance Criteria:**

**Given** I complete a dependency analysis
**When** results are saved
**Then** data is stored in IndexedDB using Dexie.js wrapper
**And** I can view previous analysis results from history
**And** stored data includes: timestamp, project path, results, health score
**And** data persists across browser sessions
**And** I can manually clear stored data from settings

---

### Story 7.3: Implement CLI Local Directory Storage

As a **CLI user**,
I want **analysis results stored in a local `.monoguard/` directory**,
So that **I can access results without any cloud dependency**.

**Acceptance Criteria:**

**Given** I run `monoguard analyze` in a project
**When** analysis completes
**Then** results are saved to `.monoguard/` directory in project root
**And** directory structure includes: `results/`, `config/`, `cache/`
**And** files are readable JSON format
**And** `.monoguard/` is automatically added to `.gitignore` (with user confirmation)
**And** old results are automatically cleaned up (configurable retention)

---

### Story 7.4: Ensure Complete Offline Functionality

As a **user**,
I want **all core features to work 100% offline**,
So that **I can use MonoGuard without any network dependency**.

**Acceptance Criteria:**

**Given** I have the web app cached or CLI installed
**When** I use MonoGuard with no network connection
**Then** I can:
- Upload files and run analysis (Web)
- Run all CLI commands (analyze, check, fix --dry-run)
- View cached historical results
- Export reports in all formats
**And** no features fail due to network unavailability
**And** Service Worker caches all required assets (Web)

---

### Story 7.5: Implement Opt-in Usage Analytics

As a **user**,
I want **anonymous usage analytics to be opt-in only**,
So that **I have full control over what data is collected**.

**Acceptance Criteria:**

**Given** I'm using MonoGuard for the first time
**When** I haven't explicitly opted in
**Then** no analytics data is collected or transmitted
**And** a clear consent banner is shown explaining what would be collected
**And** I can opt in/out at any time from settings
**And** if opted in, only anonymous usage patterns are collected (no code, no project names)
**And** privacy policy link is provided

---

### Story 7.6: Implement Opt-in Error Reporting (Sentry)

As a **user**,
I want **error reporting to be opt-in only**,
So that **I control whether crash reports are sent**.

**Acceptance Criteria:**

**Given** error reporting settings
**When** I haven't explicitly opted in
**Then** no error reports are sent to Sentry
**And** consent banner clearly explains error reporting scope
**And** if opted in, sensitive data is sanitized (no file paths, no project code)
**And** I can opt out at any time and reporting stops immediately
**And** local error logs are kept in `.monoguard/errors.log` regardless of opt-in

---

## Epic 8 Stories: Configuration & Customization

### Story 8.1: Implement Configuration File Structure

As a **developer**,
I want **a `.monoguard.json` configuration file for my project**,
So that **I can customize analysis behavior consistently across my team**.

**Acceptance Criteria:**

**Given** I want to configure MonoGuard for my project
**When** I create or generate a `.monoguard.json` file
**Then** the file supports these top-level sections:
- `version` - Config schema version
- `workspaces` - Workspace detection patterns
- `rules` - Analysis rule configurations
- `thresholds` - Health score thresholds
- `exclude` - Exclusion patterns
- `output` - Output format preferences
**And** JSON schema is provided for IDE autocomplete
**And** invalid configuration shows helpful validation errors
**And** configuration merges with defaults for unspecified values

---

### Story 8.2: Implement Circular Dependency Detection Rules

As a **developer**,
I want **to configure how circular dependencies are detected and reported**,
So that **I can customize severity levels for my project's needs**.

**Acceptance Criteria:**

**Given** a `.monoguard.json` with rules configuration
**When** I set circular dependency rules
**Then** I can configure:
- `circularDependencies`: `"error"` | `"warn"` | `"off"`
- `maxCycleDepth`: maximum cycle depth to report (default: unlimited)
- `allowedCycles`: array of allowed cycle patterns (exceptions)
**And** `"error"` causes CLI exit code 1
**And** `"warn"` reports but doesn't fail
**And** `"off"` skips circular dependency detection entirely

---

### Story 8.3: Implement Health Score Threshold Configuration

As a **developer**,
I want **to define custom health score thresholds**,
So that **I can set quality gates appropriate for my project**.

**Acceptance Criteria:**

**Given** a `.monoguard.json` with thresholds configuration
**When** I set health score thresholds
**Then** I can configure:
- `minHealthScore`: minimum acceptable score (default: 70)
- `healthWeights`: custom weights for score factors
  - `circularDependencies`: weight (default: 40)
  - `versionConflicts`: weight (default: 20)
  - `dependencyDepth`: weight (default: 20)
  - `coupling`: weight (default: 20)
**And** scores below `minHealthScore` fail CI checks
**And** custom weights sum to 100 (validated)

---

### Story 8.4: Implement Package Exclusion Pattern Configuration

As a **developer**,
I want **to configure which packages are excluded from analysis**,
So that **I can skip legacy or irrelevant packages**.

**Acceptance Criteria:**

**Given** a `.monoguard.json` with exclude configuration
**When** I set exclusion patterns
**Then** I can configure:
- Exact package names: `["packages/legacy"]`
- Glob patterns: `["packages/deprecated-*"]`
- Regex patterns: `["/^@internal\\/.*/"]`
**And** excluded packages are grayed out in visualization
**And** excluded packages don't affect health score
**And** CLI `--exclude` flag overrides/appends to config

---

### Story 8.5: Implement Workspace Detection Pattern Configuration

As a **developer**,
I want **to configure how workspaces are detected**,
So that **MonoGuard works with my custom monorepo structure**.

**Acceptance Criteria:**

**Given** a `.monoguard.json` with workspaces configuration
**When** I set workspace patterns
**Then** I can configure:
- `workspaces`: array of glob patterns (e.g., `["packages/*", "apps/*"]`)
- `packageManager`: `"auto"` | `"npm"` | `"yarn"` | `"pnpm"`
- `rootPackageJson`: path to root package.json (default: `"./package.json"`)
**And** auto-detection can be overridden
**And** custom patterns work for non-standard structures

---

### Story 8.6: Implement Output Format Configuration

As a **developer**,
I want **to configure default output formats and preferences**,
So that **I don't need to specify flags every time**.

**Acceptance Criteria:**

**Given** a `.monoguard.json` with output configuration
**When** I set output preferences
**Then** I can configure:
- `defaultFormat`: `"text"` | `"json"` | `"html"` | `"markdown"`
- `colorOutput`: `true` | `false` | `"auto"`
- `verbosity`: `"quiet"` | `"normal"` | `"verbose"`
- `includeGraph`: whether to include graph data in exports
**And** CLI flags override config settings
**And** Web UI respects `colorOutput` and `verbosity` for reports

---

## Epic 9 Stories: Developer API Integration

### Story 9.1: Create @monoguard/wasm npm Package Structure

As a **third-party developer**,
I want **to install MonoGuard analysis engine as an npm package**,
So that **I can integrate dependency analysis into my own tools**.

**Acceptance Criteria:**

**Given** I want to use MonoGuard programmatically
**When** I run `npm install @monoguard/wasm`
**Then** I receive a package that includes:
- Compiled WASM module (`monoguard.wasm`)
- TypeScript type definitions (`.d.ts` files)
- JavaScript wrapper for WASM loading
- `wasm_exec.js` Go runtime
**And** package size is < 3MB (uncompressed)
**And** package works in both Node.js and browser environments
**And** ESM and CommonJS exports are supported

---

### Story 9.2: Implement analyze() API Function

As a **third-party developer**,
I want **to call an `analyze()` function with workspace data**,
So that **I can get complete dependency analysis results programmatically**.

**Acceptance Criteria:**

**Given** I have imported `@monoguard/wasm`
**When** I call `await analyzer.analyze(workspaceData)`
**Then** I receive a `Result<AnalysisResult>` containing:
- `data.graph`: Complete dependency graph
- `data.cycles`: Array of circular dependencies
- `data.conflicts`: Array of version conflicts
- `data.healthScore`: Numeric score (0-100)
- `data.metadata`: Analysis metadata (timing, package count)
**And** errors are returned in `error` field (not thrown)
**And** function completes in < 5 seconds for 100 packages
**And** TypeScript provides full autocomplete for all fields

---

### Story 9.3: Implement check() API Function

As a **third-party developer**,
I want **a lightweight `check()` function for validation only**,
So that **I can quickly validate dependencies without full analysis overhead**.

**Acceptance Criteria:**

**Given** I have imported `@monoguard/wasm`
**When** I call `await analyzer.check(workspaceData, options)`
**Then** I receive a `Result<CheckResult>` containing:
- `data.passed`: boolean indicating pass/fail
- `data.errors`: Array of validation errors
- `data.healthScore`: Numeric score
- `data.summary`: Brief text summary
**And** options support `maxCircular`, `minHealthScore` thresholds
**And** check is faster than full analyze (< 2 seconds for 100 packages)
**And** suitable for CI/CD integration

---

### Story 9.4: Create Complete TypeScript Type Definitions

As a **third-party developer**,
I want **complete TypeScript type definitions for all API responses**,
So that **I have full type safety and IDE support when using the API**.

**Acceptance Criteria:**

**Given** I'm using `@monoguard/wasm` in a TypeScript project
**When** I import types from the package
**Then** I have access to these exported types:
- `WorkspaceData` - Input data structure
- `AnalysisResult` - Full analysis output
- `CheckResult` - Validation output
- `DependencyGraph`, `GraphNode`, `GraphEdge`
- `CircularDependency`, `VersionConflict`
- `HealthScore`, `HealthBreakdown`
- `FixSuggestion`, `FixStrategy`
- `Result<T>` - Unified result wrapper
**And** all types use camelCase (matching JSON output)
**And** types are exported from package root
**And** JSDoc comments provide inline documentation

---

### Story 9.5: Write API Documentation and Usage Examples

As a **third-party developer**,
I want **comprehensive API documentation with examples**,
So that **I can quickly understand how to integrate MonoGuard**.

**Acceptance Criteria:**

**Given** the `@monoguard/wasm` package
**When** I read the documentation
**Then** I find:
- Quick start guide (< 5 minutes to first result)
- API reference for all exported functions
- TypeScript usage examples
- Browser integration example (with WASM loading)
- Node.js integration example
- Error handling patterns
- Performance considerations
**And** documentation is included in package README
**And** examples are tested and runnable
**And** API versioning policy is documented
