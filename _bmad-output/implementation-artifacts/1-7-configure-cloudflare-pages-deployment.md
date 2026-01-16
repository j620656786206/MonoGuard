# Story 1.7: Configure Cloudflare Pages Deployment

Status: in-progress

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

- [ ] **Task 1: Create Cloudflare Pages Project** (AC: #1)
  - [ ] 1.1 Log in to Cloudflare Dashboard
  - [ ] 1.2 Navigate to Pages → Create a project
  - [ ] 1.3 Connect to GitHub repository `j620656786206/MonoGuard`
  - [ ] 1.4 Configure build settings:
    - **Framework preset:** None
    - **Build command:** `pnpm nx build web`
    - **Build output directory:** `apps/web/.output/public`
    - **Root directory:** `/` (repository root)
  - [ ] 1.5 Configure environment variables:
    ```
    NODE_VERSION=20
    PNPM_VERSION=10.14.0
    ```
  - [ ] 1.6 Set production branch to `main`

- [x] **Task 2: Create \_headers File** (AC: #4, #5)
  - [x] 2.1 Create `apps/web/public/_headers`:

    ```
    # Cloudflare Pages headers configuration
    # Reference: https://developers.cloudflare.com/pages/configuration/headers/

    # Enable SharedArrayBuffer for WASM (required for multi-threading)
    /*
      Cross-Origin-Opener-Policy: same-origin
      Cross-Origin-Embedder-Policy: require-corp
      X-Content-Type-Options: nosniff
      X-Frame-Options: DENY
      Referrer-Policy: strict-origin-when-cross-origin

    # WASM files - correct MIME type and caching
    /*.wasm
      Content-Type: application/wasm
      Cache-Control: public, max-age=31536000, immutable
      Cross-Origin-Resource-Policy: same-origin

    # JavaScript files - long cache for hashed files
    /*.js
      Cache-Control: public, max-age=31536000, immutable

    # CSS files - long cache for hashed files
    /*.css
      Cache-Control: public, max-age=31536000, immutable

    # HTML files - no cache (always fetch latest)
    /*.html
      Cache-Control: no-cache, no-store, must-revalidate

    # Index page
    /
      Cache-Control: no-cache, no-store, must-revalidate

    # Images - moderate caching
    /images/*
      Cache-Control: public, max-age=86400
    ```

- [x] **Task 3: Create \_redirects File** (AC: #2)
  - [x] 3.1 Create `apps/web/public/_redirects` for SPA routing:
    ```
    # SPA fallback - redirect all routes to index.html
    # This enables client-side routing for TanStack Router
    /*    /index.html   200
    ```

- [x] **Task 4: Configure wrangler.toml (Optional)** (AC: #1, #4)
  - [x] 4.1 Create `apps/web/wrangler.toml` for local testing:

    ```toml
    name = "monoguard"
    compatibility_date = "2024-01-01"

    [site]
    bucket = ".output/public"

    # Custom headers for local development
    [[headers]]
    for = "/*"
      [headers.values]
      Cross-Origin-Opener-Policy = "same-origin"
      Cross-Origin-Embedder-Policy = "require-corp"

    [[headers]]
    for = "/*.wasm"
      [headers.values]
      Content-Type = "application/wasm"
    ```

- [x] **Task 5: Update Build Output Path** (AC: #2)
  - [x] 5.1 Verify Vite outputs to `.output/` (confirmed - using Vite, not TanStack Start)
  - [x] 5.2 N/A - Using Vite with `vite.config.ts` (outDir: '.output')
  - [x] 5.3 WASM copy handled in deploy workflow

- [x] **Task 6: Configure GitHub Actions Deployment** (AC: #2, #3)
  - [x] 6.1 Create `.github/workflows/deploy.yml`:

    ```yaml
    # Cloudflare Pages Deployment
    #
    # Deploys to Cloudflare Pages on main branch push
    # Creates preview deployments for PRs

    name: Deploy

    on:
      push:
        branches: [main]
      pull_request:
        branches: [main]

    env:
      NODE_VERSION_FILE: '.nvmrc'
      PNPM_VERSION: '10.14.0'

    jobs:
      deploy:
        name: Deploy to Cloudflare Pages
        runs-on: ubuntu-latest
        timeout-minutes: 15
        permissions:
          contents: read
          deployments: write
          pull-requests: write
        steps:
          - name: Checkout code
            uses: actions/checkout@v4

          - name: Setup pnpm
            uses: pnpm/action-setup@v4
            with:
              version: ${{ env.PNPM_VERSION }}

          - name: Setup Node.js
            uses: actions/setup-node@v4
            with:
              node-version-file: ${{ env.NODE_VERSION_FILE }}
              cache: 'pnpm'

          - name: Setup Go (for WASM)
            uses: actions/setup-go@v5
            with:
              go-version: '1.21'

          - name: Install dependencies
            run: pnpm install --frozen-lockfile

          - name: Build WASM
            run: pnpm nx build analysis-engine

          - name: Copy WASM to public
            run: |
              mkdir -p apps/web/public
              cp packages/analysis-engine/dist/monoguard.wasm apps/web/public/
              cp packages/analysis-engine/dist/wasm_exec.js apps/web/public/

          - name: Build web app
            run: pnpm nx build web

          - name: Deploy to Cloudflare Pages
            uses: cloudflare/pages-action@v1
            with:
              apiToken: ${{ secrets.CLOUDFLARE_API_TOKEN }}
              accountId: ${{ secrets.CLOUDFLARE_ACCOUNT_ID }}
              projectName: monoguard
              directory: apps/web/.output/public
              gitHubToken: ${{ secrets.GITHUB_TOKEN }}
              # Creates preview for PRs, production for main
              branch: ${{ github.head_ref || github.ref_name }}
    ```

- [ ] **Task 7: Configure Cloudflare Secrets** (AC: #2)
  - [ ] 7.1 Generate Cloudflare API Token:
    - Go to Cloudflare Dashboard → My Profile → API Tokens
    - Create Token → Use template "Edit Cloudflare Workers"
    - Add permissions: Account → Cloudflare Pages → Edit
  - [ ] 7.2 Add secrets to GitHub repository:
    - `CLOUDFLARE_API_TOKEN`
    - `CLOUDFLARE_ACCOUNT_ID`
  - [ ] 7.3 Document token permissions in README

- [ ] **Task 8: Configure Custom Domain (Optional)** (AC: #2)
  - [ ] 8.1 Add custom domain in Cloudflare Pages settings
  - [ ] 8.2 Configure DNS records if using custom domain
  - [ ] 8.3 Enable HTTPS (automatic with Cloudflare)

- [ ] **Task 9: Verification** (AC: #2, #3, #4, #6)
  - [ ] 9.1 Push to main branch - verify production deployment
  - [ ] 9.2 Open test PR - verify preview deployment
  - [ ] 9.3 Check preview URL is posted to PR
  - [ ] 9.4 Verify site loads at deployed URL
  - [ ] 9.5 Test WASM loading:
    ```javascript
    // In browser console
    const response = await fetch('/monoguard.wasm');
    console.log('WASM Content-Type:', response.headers.get('content-type'));
    // Should be: application/wasm
    ```
  - [ ] 9.6 Verify headers with curl:
    ```bash
    curl -I https://monoguard.pages.dev/
    # Check for COOP/COEP headers
    ```
  - [ ] 9.7 Verify free tier usage in Cloudflare dashboard

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

- apps/web/public/\_headers (new)
- apps/web/public/\_redirects (new)
- apps/web/wrangler.toml (new)
- .github/workflows/deploy.yml (new)
