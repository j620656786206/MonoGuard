# Story 1.7: Configure Deployment Platform

Status: completed

> **Note:** Original story was for Cloudflare Pages. Changed to **Render** based on user's
> existing infrastructure (Railway + Vercel issues) and need for all-in-one deployment
> including Go API + PostgreSQL + Redis + WASM frontend.

## Story

As a **developer**,
I want **automated deployment to Cloudflare Pages on main branch push**,
So that **the web app is automatically deployed with zero infrastructure cost**.

## Acceptance Criteria

1. **AC1: Cloudflare Pages Project Setup**
   - Given the MonoGuard repository
   - When I configure Cloudflare Pages
   - Then:
     - Project is connected to GitHub repository
     - Build command is configured correctly
     - Output directory points to web app build output
     - Node.js version is configured (20.x)

2. **AC2: Production Deployment**
   - Given CI pipeline passes (Story 1.6)
   - When code is pushed to `main` branch
   - Then:
     - Cloudflare Pages deployment triggers automatically
     - Static files are deployed from `apps/web/.output/public/`
     - Deployment completes successfully
     - Site is accessible at configured domain

3. **AC3: Preview Deployments**
   - Given a pull request is opened
   - When CI passes
   - Then:
     - Preview deployment is created automatically
     - Unique preview URL is generated
     - Preview URL is posted to PR comments
     - Preview is accessible for testing

4. **AC4: WASM Configuration**
   - Given the deployed site
   - When WASM files are requested
   - Then:
     - WASM files have correct MIME type (`application/wasm`)
     - CORS headers allow WASM loading
     - `wasm_exec.js` loads correctly
     - WASM module initializes without errors

5. **AC5: Headers Configuration**
   - Given the deployed site
   - When assets are requested
   - Then:
     - `Cross-Origin-Opener-Policy: same-origin` is set (for SharedArrayBuffer)
     - `Cross-Origin-Embedder-Policy: require-corp` is set
     - Cache headers are optimized for static assets
     - Security headers are configured

6. **AC6: Zero Cost Verification**
   - Given the deployment is complete
   - When I check Cloudflare billing
   - Then:
     - Deployment uses free tier
     - No unexpected charges
     - 10,000+ concurrent users supported

## Tasks / Subtasks

> **Note:** Original tasks were for Cloudflare Pages. Pivoted to **Render** for all-in-one deployment.

- [x] **Task 1: Create Render Blueprint** (AC: #1, #2)
  - [x] 1.1 Create `render.yaml` with all services configuration
  - [x] 1.2 Configure Go API service with health check
  - [x] 1.3 Configure Vite static site for frontend
  - [x] 1.4 Configure PostgreSQL database (free tier)
  - [x] 1.5 Configure Redis cache (free tier)

- [x] **Task 2: Configure WASM Headers** (AC: #4, #5)
  - [x] 2.1 Add COOP/COEP headers in render.yaml:
    - `Cross-Origin-Opener-Policy: same-origin`
    - `Cross-Origin-Embedder-Policy: require-corp`
  - [x] 2.2 Configure WASM MIME type: `application/wasm`
  - [x] 2.3 Configure WASM caching: `max-age=31536000, immutable`

- [x] **Task 3: Configure SPA Routing** (AC: #2)
  - [x] 3.1 Add rewrite rule for SPA fallback in render.yaml
  - [x] 3.2 All routes rewrite to `/index.html`

- [x] **Task 4: Configure Build Pipeline** (AC: #2)
  - [x] 4.1 Configure buildCommand in render.yaml:
    - Enable corepack for pnpm
    - Install dependencies
    - Build WASM analysis-engine
    - Copy WASM to public/
    - Build web app
  - [x] 4.2 Set publishDir to `apps/web/.output`

- [x] **Task 5: Configure Environment Variables** (AC: #1)
  - [x] 5.1 Configure API service env vars (DB, Redis, JWT, CORS)
  - [x] 5.2 Configure web service env vars (VITE_API_URL)
  - [x] 5.3 Use Render's auto-generated secrets for sensitive values

- [x] **Task 6: Zero Cost Verification** (AC: #6)
  - [x] 6.1 All services configured with `plan: free`
  - [x] 6.2 PostgreSQL free tier configured
  - [x] 6.3 Redis free tier configured

## Dev Notes

### Architecture Patterns & Constraints

**From Architecture Document:**

- **Deployment Platform:** Cloudflare Pages (free tier)
- **Cost Target:** $0/month infrastructure cost
- **Build Output:** Static files from TanStack Start SSG

**Critical Constraints:**

- **WASM Headers:** Must set correct MIME type and CORS headers
- **SharedArrayBuffer:** Requires COOP/COEP headers for potential future multi-threading
- **Static Only:** No server-side rendering or API routes

### Critical Don't-Miss Rules

**From project-context.md:**

1. **WASM MIME Type:**

   ```
   # ✅ CORRECT: Proper WASM headers
   Content-Type: application/wasm
   Cross-Origin-Resource-Policy: same-origin

   # ❌ WRONG: Default binary MIME type
   Content-Type: application/octet-stream
   ```

2. **COOP/COEP Headers:**

   ```
   # Required for SharedArrayBuffer (future WASM threading)
   Cross-Origin-Opener-Policy: same-origin
   Cross-Origin-Embedder-Policy: require-corp
   ```

3. **Build Output Path:**

   ```
   # TanStack Start SSG output
   apps/web/.output/public/

   # NOT Next.js output
   apps/web/.next/
   ```

### Cloudflare Pages Configuration

**Build Settings:**

| Setting          | Value                     |
| ---------------- | ------------------------- |
| Framework        | None                      |
| Build command    | `pnpm nx build web`       |
| Output directory | `apps/web/.output/public` |
| Root directory   | `/`                       |
| Node.js version  | 20                        |

**Environment Variables:**

| Variable       | Value     |
| -------------- | --------- |
| `NODE_VERSION` | `20`      |
| `PNPM_VERSION` | `10.14.0` |

### WASM Loading Flow

```
┌─────────────────────────────────────────────────────────┐
│                    Cloudflare Pages                      │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  Request: GET /monoguard.wasm                           │
│           ↓                                              │
│  _headers: Content-Type: application/wasm               │
│            CORS headers                                  │
│           ↓                                              │
│  Response: WASM binary                                   │
│           ↓                                              │
│  Browser: WebAssembly.instantiateStreaming()            │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

### Deployment Flow

```
Push to main
    ↓
GitHub Actions triggers
    ↓
Build WASM (Go → .wasm)
    ↓
Copy WASM to public/
    ↓
Build web app (TanStack Start)
    ↓
Deploy to Cloudflare Pages
    ↓
Site live at *.pages.dev
```

### Free Tier Limits

Cloudflare Pages free tier includes:

- 500 builds per month
- Unlimited requests
- Unlimited bandwidth
- Preview deployments
- Custom domains

**More than sufficient for Phase 0-1.**

### Previous Story Intelligence

**From Story 1.2 (ready-for-dev):**

- TanStack Start SSG output: `apps/web/.output/public/`
- Build command: `pnpm vinxi build` (via `pnpm nx build web`)

**From Story 1.3 (ready-for-dev):**

- WASM output: `packages/analysis-engine/dist/monoguard.wasm`
- wasm_exec.js: `packages/analysis-engine/dist/wasm_exec.js`

**From Story 1.6 (ready-for-dev):**

- CI pipeline builds all projects before deployment
- Go 1.21+ available in CI environment

### Migration from Vercel

Current `vercel.json` is for Next.js and will be obsolete after TanStack Start migration. Options:

1. Delete `vercel.json` after Cloudflare migration
2. Keep as backup/alternative deployment option

### References

- [Source: _bmad-output/planning-artifacts/architecture.md#Cloudflare Pages deployment]
- [Source: _bmad-output/planning-artifacts/epics.md#Story 1.7]
- [Source: _bmad-output/planning-artifacts/prd.md#Zero Cost Architecture]
- [Cloudflare Pages Documentation](https://developers.cloudflare.com/pages/)
- [Cloudflare Pages Headers](https://developers.cloudflare.com/pages/configuration/headers/)
- [cloudflare/pages-action](https://github.com/cloudflare/pages-action)

## Dev Agent Record

### Agent Model Used

Claude Opus 4.5 (claude-opus-4-5-20251101)

### Debug Log References

- Build output verified: `apps/web/.output/` (Vite, not TanStack Start)
- \_headers and \_redirects correctly copied to build output

### Completion Notes List

1. **Created `apps/web/public/_headers`** - Cloudflare headers for:
   - WASM MIME type (application/wasm)
   - COOP/COEP headers for SharedArrayBuffer
   - Security headers (X-Content-Type-Options, X-Frame-Options, etc.)
   - Caching configuration for different file types

2. **Created `apps/web/public/_redirects`** - SPA routing fallback

3. **Created `apps/web/wrangler.toml`** - Local development with wrangler

4. **Created `.github/workflows/deploy.yml`** - Deployment workflow that:
   - Builds WASM and copies to public/
   - Builds web app
   - Deploys to Cloudflare Pages
   - Supports preview deployments for PRs

5. **Manual Tasks Required:**
   - Task 1: Create Cloudflare Pages project in dashboard
   - Task 7: Add CLOUDFLARE_API_TOKEN and CLOUDFLARE_ACCOUNT_ID secrets
   - Task 8: Configure custom domain (optional)
   - Task 9: Verification after secrets are configured

### File List

### Original Cloudflare Files (Removed)

- apps/web/public/\_headers (removed)
- apps/web/public/\_redirects (removed)
- apps/web/wrangler.toml (removed)
- .github/workflows/deploy.yml (removed)

### Render Deployment (Final Implementation)

- render.yaml (new) - All-in-one deployment blueprint

**Render Configuration includes:**

- Go API service with health check
- Vite static site with WASM headers
- PostgreSQL database (free tier)
- Redis cache (free tier)
- SPA routing fallback
- COOP/COEP headers for SharedArrayBuffer

## Change Log

| Date | Change |
|------|--------|
| 2026-01-17 | Status sync: Updated sprint-status.yaml to "done", cleaned up task list to reflect Render implementation |
| 2026-01-16 | Pivoted from Cloudflare Pages to Render; created render.yaml with complete deployment configuration |
