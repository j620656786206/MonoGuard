# Story 1.2: Setup TanStack Start Web Application

Status: ready-for-dev

## Story

As a **developer**,
I want **a TanStack Start application configured for SSG deployment**,
So that **I can build the web interface with modern React tooling and zero-backend architecture**.

## Acceptance Criteria

1. **AC1: TanStack Start Project Configuration**
   - Given the Nx monorepo from Story 1.1
   - When I initialize the TanStack Start application in apps/web
   - Then I have:
     - TanStack Start project with SSG configuration
     - React 19.0.0 with TypeScript 5.9+
     - Vite build setup with proper configuration
     - Package name `@monoguard/web`

2. **AC2: Routing Structure**
   - Given the TanStack Start application
   - When I verify the routing setup
   - Then I have file-based routing with:
     - `/` (index route - landing page placeholder)
     - `/analyze` (analysis page placeholder)
     - `/results` (results page placeholder)
     - `__root.tsx` for root layout

3. **AC3: Tailwind CSS Integration**
   - Given the TanStack Start application
   - When I verify styling setup
   - Then Tailwind CSS is integrated with:
     - JIT mode enabled
     - PostCSS configuration
     - Base styles in globals.css
     - Working utility classes in components

4. **AC4: Development Server**
   - Given the complete setup
   - When I run the development server
   - Then the server runs at localhost:3000
   - And hot module replacement (HMR) works correctly

5. **AC5: SSG Build Output**
   - Given the TanStack Start application
   - When I run `pnpm nx build web`
   - Then:
     - Build produces static HTML output
     - Bundle size is under 100KB gzipped (without WASM)
     - No server-side rendering dependencies
     - Output compatible with Cloudflare Pages deployment

6. **AC6: Nx Integration**
   - Given the updated application
   - When I run Nx commands
   - Then:
     - `pnpm nx build web` succeeds
     - `pnpm nx dev web` starts development server
     - `pnpm nx lint web` runs ESLint checks
     - Project appears correctly in `pnpm nx graph`

## Tasks / Subtasks

- [ ] **Task 1: Remove Next.js and Setup TanStack Start** (AC: #1)
  - [ ] 1.1 Backup existing `apps/web/` content (components, assets, etc.)
  - [ ] 1.2 Remove Next.js dependencies from `apps/web/package.json`:
    - Remove `next`, `@next/*` packages
    - Remove Next.js specific config files (`next.config.js`, `next-env.d.ts`)
  - [ ] 1.3 Install TanStack Start dependencies:
    ```bash
    pnpm add @tanstack/react-start @tanstack/react-router vinxi vite --filter @monoguard/web
    pnpm add -D @tanstack/router-devtools @vitejs/plugin-react --filter @monoguard/web
    ```
  - [ ] 1.4 Update `apps/web/package.json` with correct scripts:
    ```json
    {
      "scripts": {
        "dev": "vinxi dev",
        "build": "vinxi build",
        "start": "vinxi start"
      }
    }
    ```

- [ ] **Task 2: Configure TanStack Start for SSG** (AC: #1, #5)
  - [ ] 2.1 Create `apps/web/app.config.ts`:
    ```typescript
    import { defineConfig } from '@tanstack/react-start/config'
    import viteTsConfigPaths from 'vite-tsconfig-paths'

    export default defineConfig({
      vite: {
        plugins: [
          viteTsConfigPaths({
            root: '../../',
          }),
        ],
      },
      server: {
        preset: 'static', // SSG mode - critical for zero-backend
      },
    })
    ```
  - [ ] 2.2 Create `apps/web/app/routes/__root.tsx` (root layout)
  - [ ] 2.3 Create `apps/web/app/router.tsx` (router configuration)
  - [ ] 2.4 Create `apps/web/app/client.tsx` (client entry point)

- [ ] **Task 3: Setup File-Based Routing** (AC: #2)
  - [ ] 3.1 Create route structure in `apps/web/app/routes/`:
    ```
    routes/
    ├── __root.tsx     # Root layout with <Outlet />
    ├── index.tsx      # Landing page (/)
    ├── analyze.tsx    # Analysis page (/analyze)
    └── results.tsx    # Results page (/results)
    ```
  - [ ] 3.2 Each route should export:
    - `Route` from `createFileRoute`
    - Basic placeholder component
  - [ ] 3.3 Verify routes render correctly via dev server

- [ ] **Task 4: Configure Tailwind CSS with JIT** (AC: #3)
  - [ ] 4.1 Install Tailwind dependencies (if not present):
    ```bash
    pnpm add -D tailwindcss postcss autoprefixer --filter @monoguard/web
    ```
  - [ ] 4.2 Update `apps/web/tailwind.config.ts`:
    ```typescript
    import type { Config } from 'tailwindcss'

    export default {
      content: [
        './app/**/*.{js,ts,jsx,tsx}',
        './src/**/*.{js,ts,jsx,tsx}',
      ],
      theme: {
        extend: {},
      },
      plugins: [],
    } satisfies Config
    ```
  - [ ] 4.3 Create/update `apps/web/postcss.config.js`
  - [ ] 4.4 Update `apps/web/app/styles/globals.css` with Tailwind directives:
    ```css
    @tailwind base;
    @tailwind components;
    @tailwind utilities;
    ```
  - [ ] 4.5 Import globals.css in root layout

- [ ] **Task 5: Update Nx Project Configuration** (AC: #6)
  - [ ] 5.1 Update `apps/web/project.json` targets:
    ```json
    {
      "targets": {
        "build": {
          "executor": "nx:run-commands",
          "options": {
            "command": "pnpm vinxi build",
            "cwd": "apps/web"
          },
          "outputs": ["{projectRoot}/.output"]
        },
        "dev": {
          "executor": "nx:run-commands",
          "options": {
            "command": "pnpm vinxi dev",
            "cwd": "apps/web"
          }
        },
        "start": {
          "executor": "nx:run-commands",
          "options": {
            "command": "pnpm vinxi start",
            "cwd": "apps/web"
          }
        }
      }
    }
    ```
  - [ ] 5.2 Update TypeScript paths in `apps/web/tsconfig.json`
  - [ ] 5.3 Verify Nx dependency graph still works

- [ ] **Task 6: Migrate Reusable Components** (AC: #1)
  - [ ] 6.1 Review backed up components from Task 1.1
  - [ ] 6.2 Migrate reusable components to TanStack Start structure:
    - Move to `apps/web/app/components/`
    - Update imports (remove Next.js specific imports like `next/link`, `next/image`)
    - Replace with TanStack Router equivalents (`Link` from `@tanstack/react-router`)
  - [ ] 6.3 Update component file structure to follow conventions:
    - PascalCase.tsx for components
    - camelCase.ts for utilities

- [ ] **Task 7: Static Assets Configuration** (AC: #5)
  - [ ] 7.1 Move static assets to `apps/web/public/`:
    - Favicon files
    - Logo images
    - Any other static assets
  - [ ] 7.2 Verify assets are accessible in build output
  - [ ] 7.3 Configure asset handling in Vite config if needed

- [ ] **Task 8: Verification and Testing** (AC: #4, #5, #6)
  - [ ] 8.1 Run `pnpm nx dev web` - verify localhost:3000 works
  - [ ] 8.2 Verify HMR works (modify component, see instant update)
  - [ ] 8.3 Run `pnpm nx build web` - verify SSG output
  - [ ] 8.4 Verify build output structure:
    ```
    apps/web/.output/
    ├── public/         # Static assets
    └── server/         # (should be minimal for SSG)
    ```
  - [ ] 8.5 Check bundle size with:
    ```bash
    du -sh apps/web/.output/public
    ```
  - [ ] 8.6 Run `pnpm nx graph` - verify project dependencies correct
  - [ ] 8.7 Run `pnpm nx lint web` - verify no lint errors

## Dev Notes

### Architecture Patterns & Constraints

**From Architecture Document:**

- **Framework Choice:** TanStack Start 0.34+ with SSG preset (Decision in architecture.md#Starter Template Evaluation)
- **Deployment Target:** Cloudflare Pages (static hosting, $0/month)
- **Build Tooling:** Vite via Vinxi bundler
- **Styling:** Tailwind CSS with JIT mode for optimized bundle size

**Critical Constraints:**

- **Zero Backend:** SSG mode ONLY - no SSR, no API routes, no server-side data fetching
- **Offline First:** Web app must work without network after initial load
- **Bundle Size:** < 500KB gzipped total (< 100KB without WASM for this story)
- **Performance:** FCP < 1.5 seconds, Lighthouse Performance > 90

### Critical Don't-Miss Rules

**From project-context.md:**

1. **TanStack Start is NOT Next.js:**
   - Do NOT use `getServerSideProps` or `getStaticProps` (Next.js patterns)
   - Use file-based routing in `app/routes/` directory
   - Use `createFileRoute` for route definitions
   - SSG preset generates static HTML at build time

2. **React 19 Compatibility:**
   - TanStack Start 0.34+ required for React 19 support
   - TypeScript 5.9+ required for React 19 types
   - Use hooks pattern, avoid class components

3. **Zustand Integration (future):**
   - State management will be added in Story 5.9
   - Keep components stateless for now, use props
   - Don't introduce any global state in this story

4. **Naming Conventions:**
   - TypeScript: camelCase (variables/functions), PascalCase (types/interfaces/components)
   - Files: PascalCase.tsx (React components), camelCase.ts (utilities)
   - Routes: lowercase with dashes (e.g., `analyze.tsx`, `results.$id.tsx`)

### Project Structure Notes

**Target Directory Structure (apps/web/):**

```
apps/web/
├── app/
│   ├── components/          # Shared UI components
│   │   ├── common/          # Generic reusable components
│   │   └── ui/              # UI primitives
│   ├── routes/              # File-based routing
│   │   ├── __root.tsx       # Root layout
│   │   ├── index.tsx        # / route
│   │   ├── analyze.tsx      # /analyze route
│   │   └── results.tsx      # /results route
│   ├── styles/
│   │   └── globals.css      # Tailwind + global styles
│   ├── client.tsx           # Client entry point
│   └── router.tsx           # Router configuration
├── public/                   # Static assets
│   ├── favicon.ico
│   └── ...
├── app.config.ts             # TanStack Start config
├── tailwind.config.ts        # Tailwind configuration
├── postcss.config.js         # PostCSS configuration
├── tsconfig.json             # TypeScript config
├── package.json
└── project.json              # Nx project config
```

**Migration Notes (Next.js → TanStack Start):**

| Next.js                  | TanStack Start                  |
|--------------------------|----------------------------------|
| `pages/` or `app/`       | `app/routes/`                   |
| `next/link`              | `@tanstack/react-router` Link   |
| `next/image`             | Standard `<img>` or lazy loading|
| `next/head`              | React 19 `<title>`, `<meta>`    |
| `getStaticProps`         | Client-side data loading        |
| `next.config.js`         | `app.config.ts`                 |
| `next dev`               | `vinxi dev`                     |
| `next build`             | `vinxi build`                   |

### WASM Integration (Future Story Reference)

This story does NOT include WASM integration. However, the structure should support future WASM loading:

- WASM files will be placed in `apps/web/public/` directory
- TypeScript WASM adapter will be in `packages/analysis-engine/`
- Story 2.7 will implement the TypeScript WASM Adapter

### Testing Requirements

**Verification Checklist:**

- [ ] `pnpm nx dev web` starts server on localhost:3000
- [ ] All three routes (`/`, `/analyze`, `/results`) render correctly
- [ ] Tailwind utility classes work (e.g., `className="bg-blue-500"`)
- [ ] `pnpm nx build web` produces static output
- [ ] Bundle size < 100KB gzipped (check with `du -sh`)
- [ ] `pnpm nx graph` shows correct dependency graph
- [ ] No TypeScript or lint errors

### Previous Story Intelligence

**From Story 1.1 (done):**

- Nx workspace restructured: `libs/` → `packages/`, `apps/frontend/` → `apps/web/`
- Path mappings updated: `@monoguard/types`, `@monoguard/ui-components`, `@monoguard/analysis-engine`
- Current `apps/web/` uses Next.js 15.x - this story migrates to TanStack Start
- `apps/web/package.json` already named `@monoguard/web`
- TypeScript imports using `@monoguard/*` workspace paths

**Key Files Modified in 1.1:**

- `/nx.json` - workspace layout configuration
- `/tsconfig.base.json` - path mappings
- `/apps/web/project.json` - Nx targets (will be updated)
- `/apps/web/package.json` - dependencies (will be updated)

### Git Intelligence Summary

**Recent Commits:**

- `a727331` - feat: create epic
- `68d27eb` - feat: finish product prd, ux, architecture analyze with bmad
- `8eb3807` - Fix React Server Components CVE vulnerabilities

**Note:** The CVE fix commit suggests security-sensitive code. Ensure TanStack Start migration doesn't reintroduce vulnerabilities.

### References

- [Source: _bmad-output/planning-artifacts/architecture.md#Starter Template Evaluation]
- [Source: _bmad-output/planning-artifacts/architecture.md#Decision 2: Web Framework (TanStack Start)]
- [Source: _bmad-output/project-context.md#Framework-Specific Rules - TanStack Start]
- [Source: _bmad-output/planning-artifacts/epics.md#Story 1.2]
- [Source: _bmad-output/implementation-artifacts/1-1-initialize-nx-monorepo-workspace.md]
- [TanStack Start Documentation](https://tanstack.com/router/latest/docs/framework/react/start)
- [Vinxi Documentation](https://vinxi.vercel.app/)

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List

