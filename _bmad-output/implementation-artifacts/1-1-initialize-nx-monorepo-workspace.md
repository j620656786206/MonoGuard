# Story 1.1: Initialize Nx Monorepo Workspace

Status: done

## Story

As a **developer**,
I want **a properly configured Nx monorepo workspace with apps/ and packages/ directories**,
So that **I have a standardized project structure that supports multiple applications and shared packages**.

## Acceptance Criteria

1. **AC1: Nx Workspace Structure**
   - Given a fresh project directory
   - When I run the initialization commands
   - Then I have an Nx workspace with:
     - `apps/` directory for applications (web, cli)
     - `packages/` directory for shared libraries (analysis-engine, types, ui-components)
     - `nx.json` with proper workspace configuration
     - `package.json` with workspace scripts
     - `tsconfig.base.json` for TypeScript path mapping

2. **AC2: Nx Graph Verification**
   - Given the initialized workspace
   - When I run `npx nx graph`
   - Then the project structure is displayed correctly

3. **AC3: Development Scripts**
   - Given the workspace setup
   - When I check `package.json`
   - Then I have scripts for: `lint`, `test`, `build`, `dev`

4. **AC4: Path Mapping Configuration**
   - Given `tsconfig.base.json`
   - When I verify path mappings
   - Then `@monoguard/*` paths are configured for all packages

## Tasks / Subtasks

- [x] **Task 1: Initialize Nx Workspace** (AC: #1)
  - [x] 1.1 Run `npx create-nx-workspace@latest mono-guard --preset=ts`
  - [x] 1.2 Configure workspace layout in `nx.json` (appsDir: "apps", libsDir: "packages")
  - [x] 1.3 Set pnpm as package manager in `nx.json`

- [x] **Task 2: Create Directory Structure** (AC: #1)
  - [x] 2.1 Create `apps/` directory structure:
    - `apps/web/` (currently Next.js - **will migrate to TanStack Start in Story 1.2**)
    - `apps/cli/` (currently TypeScript/Node.js - **will migrate to Go CLI in Story 1.4**)
  - [x] 2.2 Create `packages/` directory structure:
    - `packages/analysis-engine/` (placeholder for Go WASM - Story 1.3)
    - `packages/types/` (placeholder for shared types - Story 1.5)
    - `packages/ui-components/` (placeholder for React components)

- [x] **Task 3: Configure TypeScript Path Mappings** (AC: #4)
  - [x] 3.1 Create/update `tsconfig.base.json` with path aliases:
    ```json
    {
      "compilerOptions": {
        "paths": {
          "@monoguard/types": ["packages/types/src/index.ts"],
          "@monoguard/ui-components": ["packages/ui-components/src/index.ts"]
        }
      }
    }
    ```

- [x] **Task 4: Configure Nx Project Settings** (AC: #1, #2)
  - [x] 4.1 Update `nx.json` with:
    - `defaultBase: "main"`
    - `targetDefaults` for build and test caching
    - `namedInputs` for source and production files
  - [x] 4.2 Configure task pipelines for dependency ordering

- [x] **Task 5: Setup Root Package.json Scripts** (AC: #3)
  - [x] 5.1 Add workspace scripts:
    ```json
    {
      "scripts": {
        "dev": "nx run-many --target=dev --all",
        "build": "nx run-many --target=build --all",
        "test": "nx run-many --target=test --all",
        "lint": "nx run-many --target=lint --all",
        "affected:build": "nx affected --target=build",
        "affected:test": "nx affected --target=test",
        "graph": "nx graph"
      }
    }
    ```

- [x] **Task 6: Verification** (AC: #2)
  - [x] 6.1 Run `npx nx graph` to verify structure displays correctly
  - [x] 6.2 Run `pnpm install` to verify workspace is valid
  - [x] 6.3 Verify Nx can detect the workspace layout

## Dev Notes

### Architecture Patterns & Constraints

**From Architecture Document:**

- **Monorepo Strategy:** Nx Monorepo (Decision 1 in architecture.md)
- **Workspace Layout:**
  - `apps/` for applications (web, cli)
  - `packages/` for shared libraries (analysis-engine, types, ui-components)
- **Package Manager:** pnpm 10.14.0 (from project-context.md)
- **Nx Version:** 21.4.1 (from project-context.md)

**Critical Constraints:**

- Must use pnpm, NOT npm or yarn (project-context.md)
- All Nx commands must use `pnpm nx` prefix
- Build caching must be enabled for performance

### Project Structure Notes

**Current vs Target State:**

| App/Package                 | Current State                | Target State (Future Story)      |
| --------------------------- | ---------------------------- | -------------------------------- |
| `apps/web/`                 | Next.js 15.x                 | TanStack Start 0.34+ (Story 1.2) |
| `apps/cli/`                 | TypeScript/Node.js (esbuild) | Go CLI (Story 1.4)               |
| `packages/analysis-engine/` | TypeScript placeholder       | Go WASM (Story 1.3)              |

**Target Directory Structure:**

```
mono-guard/                    # Nx workspace root
├── apps/
│   ├── web/                   # TanStack Start frontend (Story 1.2 migration)
│   │   ├── app/
│   │   ├── public/
│   │   └── package.json
│   └── cli/                   # Go CLI tool (Story 1.4 migration)
│       ├── cmd/
│       ├── pkg/
│       └── go.mod
├── packages/
│   ├── analysis-engine/       # Go WASM core (Story 1.3)
│   │   ├── cmd/wasm/
│   │   ├── pkg/
│   │   └── go.mod
│   ├── types/                 # Shared TypeScript types (Story 1.5)
│   │   ├── src/
│   │   └── package.json
│   └── ui-components/         # Shared React components
│       ├── src/
│       └── package.json
├── nx.json
├── package.json
├── pnpm-workspace.yaml
└── tsconfig.base.json
```

### Naming Conventions

**From project-context.md:**

- TypeScript: camelCase (variables/functions), PascalCase (types/interfaces/components)
- Files: PascalCase.tsx (React components), camelCase.ts (utilities)
- Imports: Use Nx workspace paths (`@monoguard/*`)

### Configuration Files Reference

**nx.json structure:**

```json
{
  "$schema": "./node_modules/nx/schemas/nx-schema.json",
  "affected": {
    "defaultBase": "main"
  },
  "targetDefaults": {
    "build": {
      "dependsOn": ["^build"],
      "cache": true
    },
    "test": {
      "cache": true
    },
    "lint": {
      "cache": true
    }
  },
  "workspaceLayout": {
    "appsDir": "apps",
    "libsDir": "packages"
  },
  "namedInputs": {
    "default": ["{projectRoot}/**/*", "sharedGlobals"],
    "production": [
      "default",
      "!{projectRoot}/**/*.spec.ts",
      "!{projectRoot}/**/*.test.ts",
      "!{projectRoot}/**/__tests__/**"
    ],
    "sharedGlobals": []
  },
  "plugins": []
}
```

**pnpm-workspace.yaml:**

```yaml
packages:
  - 'apps/*'
  - 'packages/*'
```

### Critical Don't-Miss Rules

**From project-context.md:**

1. **NEVER use npm install** - project uses pnpm exclusively
2. **Use `pnpm nx` for all Nx commands**
3. **Cache directories (.nx/cache, .pnpm-store/) must NOT be committed**
4. **Node.js >= 18.0.0 required**

### Testing Requirements

- Verify `npx nx graph` displays project structure
- Verify `pnpm install` completes without errors
- Verify workspace scripts are functional

### References

- [Source: _bmad-output/planning-artifacts/architecture.md#Decision 1: Monorepo Strategy]
- [Source: _bmad-output/planning-artifacts/architecture.md#Starter Template Evaluation]
- [Source: _bmad-output/project-context.md#Technology Stack & Versions]
- [Source: _bmad-output/project-context.md#Development Workflow Rules]
- [Source: _bmad-output/planning-artifacts/epics.md#Story 1.1]

## Dev Agent Record

### Agent Model Used

Claude Opus 4.5 (claude-opus-4-5-20251101)

### Completion Notes List

- Story created via YOLO mode with comprehensive context from PRD, Architecture, and project-context.md
- This is the foundation story - all subsequent stories depend on this workspace structure
- **Major Restructure Performed:** Existing workspace was restructured to match story requirements:
  - Renamed `libs/` to `packages/`
  - Renamed `libs/shared-types/` to `packages/types/`
  - Renamed `apps/frontend/` to `apps/web/`
  - Renamed `apps/frontend-e2e/` to `apps/web-e2e/`
  - Created `packages/analysis-engine/` placeholder
  - Created `packages/ui-components/` placeholder
  - Updated all configuration files (nx.json, package.json, pnpm-workspace.yaml, tsconfig.base.json)
  - Updated all import statements from `@monoguard/shared-types` to `@monoguard/types`
  - Updated docker-compose files, vercel.json, and other config files
- Nx graph verification passed - all 6 projects detected: web, web-e2e, cli, @monoguard/types, @monoguard/analysis-engine, @monoguard/ui-components
- pnpm install completed successfully (some peer dependency warnings exist but are unrelated to restructuring)

### File List

Files created/modified:

**Configuration Files:**

- `/nx.json` - Added workspaceLayout configuration, updated defaultProject to "web"
- `/package.json` - Updated workspaces from libs to packages, renamed dev:frontend to dev:web
- `/pnpm-workspace.yaml` - Changed libs/_ to packages/_
- `/tsconfig.base.json` - Updated path mappings to @monoguard/types, @monoguard/ui-components, @monoguard/analysis-engine
- `/vercel.json` - Updated paths from apps/frontend to apps/web
- `/docker-compose.yml` - Renamed frontend service to web
- `/docker-compose.prod.yml` - Renamed frontend service to web
- `/.vscode/launch.json` - Fixed CLI path reference

**Packages (New):**

- `/packages/types/` - Renamed from libs/shared-types, package name changed to @monoguard/types
- `/packages/analysis-engine/` - New placeholder package
- `/packages/analysis-engine/package.json` - Package configuration
- `/packages/analysis-engine/src/index.ts` - Placeholder entry point
- `/packages/analysis-engine/tsconfig.json` - TypeScript configuration
- `/packages/ui-components/` - New placeholder package
- `/packages/ui-components/package.json` - Package configuration
- `/packages/ui-components/src/index.ts` - Placeholder entry point
- `/packages/ui-components/tsconfig.json` - TypeScript configuration

**Apps (Renamed):**

- `/apps/web/` - Renamed from apps/frontend
- `/apps/web/project.json` - Updated all references from frontend to web
- `/apps/web/package.json` - Renamed to @monoguard/web, updated dependency to @monoguard/types
- `/apps/web/tsconfig.json` - Updated references from frontend to web
- `/apps/web/zbpack.json` - Updated nx commands from frontend to web
- `/apps/web-e2e/` - Renamed from apps/frontend-e2e
- `/apps/web-e2e/project.json` - Updated all references
- `/apps/cli/package.json` - Updated dependency to @monoguard/types

**Source Files (Import Updates):**

- All TypeScript files in apps/web that imported @monoguard/shared-types now import @monoguard/types
- 16 files updated with new import paths

**Review Fixes (2026-01-15):**

- `/pnpm-workspace.yaml` - Removed non-existent `tools/*` entry
- `/package.json` - Synced workspaces config with pnpm-workspace.yaml

## Review Follow-ups (AI)

### Action Items from Code Review

- [ ] [AI-Review][MEDIUM] Add `project.json` to `packages/types/` for explicit Nx target configuration
- [ ] [AI-Review][MEDIUM] Add `project.json` to `packages/analysis-engine/` for explicit Nx target configuration
- [ ] [AI-Review][MEDIUM] Add `project.json` to `packages/ui-components/` for explicit Nx target configuration
- [x] [AI-Review][LOW] Remove `tools/*` from `pnpm-workspace.yaml` (directory does not exist) ✅ Fixed

### Notes

- **TanStack Start Migration:** Current `apps/web/` uses Next.js 15.x. Migration to TanStack Start 0.34+ is planned for Story 1.2. This is intentional - Story 1.1 focuses on workspace structure, not framework migration.
- **Go CLI Migration:** Current `apps/cli/` uses TypeScript/Node.js with esbuild. Migration to Go is planned for Story 1.4.
- **Nx Project Detection:** Nx 21.4.1 can infer projects from `package.json` files, so missing `project.json` files do not block functionality. However, explicit `project.json` files are recommended for better Nx integration.

## Senior Developer Review (AI)

**Review Date:** 2026-01-15
**Reviewer:** Claude Opus 4.5 (Developer Agent)
**Outcome:** ✅ APPROVED with Action Items

**Summary:**

- All 4 Acceptance Criteria verified and PASSED
- All 6 Tasks verified as complete
- 4 action items created for future improvement (non-blocking)
- Story correctly establishes workspace structure for subsequent stories

## Change Log

| Date       | Change                                                                                                                           | Author                     |
| ---------- | -------------------------------------------------------------------------------------------------------------------------------- | -------------------------- |
| 2026-01-16 | Code review #2: Removed invalid @monoguard/analysis-engine path from tsconfig.base.json (directory converted to Go in Story 1.3) | Amelia (Developer Agent)   |
| 2026-01-16 | Code review #2: Removed legacy Next.js configuration from nx.json, updated bundler to vite                                       | Amelia (Developer Agent)   |
| 2026-01-15 | Code review completed - approved with action items                                                                               | Claude Opus 4.5 (Reviewer) |
| 2026-01-15 | Story completed - major restructure from existing workspace to match story requirements                                          | Claude Opus 4.5            |
