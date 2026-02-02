# Story 5.1: Create Landing Page with Drag-Drop Zone

Status: ready-for-dev

## Story

As a **user**,
I want **a welcoming landing page with a prominent drag-drop zone**,
So that **I can immediately start analyzing my monorepo without any setup**.

## Acceptance Criteria

### AC1: Landing Page Visual Design
**Given** I visit the MonoGuard web app
**When** the landing page loads
**Then** I see:
- Clear value proposition headline (e.g., "Fix your monorepo's circular dependencies")
- Prominent drag-drop zone that is visually distinct and centered
- "Or click to select files" option within the drop zone
- Privacy badge displaying "100% local analysis - your code never leaves your machine"
- Sample demo button ("Try with example project")
- Key feature highlights below the hero section

### AC2: Performance Requirements
**Given** the landing page
**When** it loads on a standard connection
**Then**:
- First Contentful Paint (FCP) < 1.5 seconds
- Time to Interactive (TTI) < 3.8 seconds
- No layout shifts during load (CLS < 0.1)
- All above-the-fold content renders without JavaScript (SSG)

### AC3: Drag-Drop Zone Interaction States
**Given** the drag-drop zone on the landing page
**When** I interact with it
**Then**:
- Default state: Dashed border with upload icon and instructional text
- Hover state: Border highlights, background subtly changes color
- Active/dragging state: Border becomes solid, stronger visual feedback with "Drop files here" text
- File accepted state: Brief success animation before navigating to analysis
- File rejected state: Error message with supported file types listed

### AC4: Zero Authentication Requirement
**Given** I visit the MonoGuard web app for the first time
**When** the landing page loads
**Then**:
- No login/registration prompts or modals
- No cookie consent required for core functionality
- Full functionality available immediately
- No account-gated features visible

### AC5: Navigation Structure
**Given** the landing page
**When** I look at the page structure
**Then**:
- Simple top navigation bar with MonoGuard logo/name
- Navigation links: Home, Analyze, Results (if previous analysis exists)
- GitHub link in nav or footer
- Footer with privacy policy link and version info

### AC-CI: CI Pipeline Must Pass (MANDATORY)
**Given** the story implementation is complete
**When** verifying CI status
**Then** ALL of the following must pass:
- [ ] `pnpm nx affected --target=lint --base=main` passes
- [ ] `pnpm nx affected --target=test --base=main` passes
- [ ] `pnpm nx affected --target=type-check --base=main` passes
- [ ] GitHub Actions CI workflow shows GREEN status
- **Story CANNOT be marked as "done" until CI is green**

## Tasks / Subtasks

- [ ] Task 1: Redesign Landing Page Route (AC: 1, 2, 4)
  - [ ] Refactor `apps/web/app/routes/index.tsx` - replace current basic layout with full landing page
  - [ ] Create hero section with value proposition headline and subtext
  - [ ] Add centered drag-drop zone component (reuse/extend `FileUpload.tsx`)
  - [ ] Add privacy badge component below drop zone
  - [ ] Add "Try with example" demo button

- [ ] Task 2: Enhance Drag-Drop Zone Component (AC: 3)
  - [ ] Extend existing `apps/web/app/components/common/FileUpload.tsx` or create new `LandingDropZone.tsx`
  - [ ] Implement all interaction states (default, hover, active, accepted, rejected)
  - [ ] Add visual feedback animations with Tailwind CSS transitions
  - [ ] Add file type validation (accept .json, .yaml, .yml, folder)
  - [ ] Wire up file drop to navigate to /analyze route with file data

- [ ] Task 3: Create/Update Landing Page Components (AC: 1, 5)
  - [ ] Update `apps/web/app/components/landing/HeroSection.tsx` with new design
  - [ ] Update `apps/web/app/components/landing/FeaturesSection.tsx` with key capability highlights
  - [ ] Create privacy badge component `apps/web/app/components/landing/PrivacyBadge.tsx`
  - [ ] Update `apps/web/app/components/landing/Footer.tsx` with version and links
  - [ ] Add/update navigation bar in `apps/web/app/routes/__root.tsx`

- [ ] Task 4: Implement Sample Demo Flow (AC: 1)
  - [ ] Create sample workspace data fixture (small example monorepo)
  - [ ] Store sample data in `apps/web/app/lib/sampleData.ts`
  - [ ] Wire "Try with example" button to load sample and navigate to analysis

- [ ] Task 5: Performance Optimization (AC: 2)
  - [ ] Ensure landing page is SSG-compatible (no server-side data fetching)
  - [ ] Optimize images with lazy loading below the fold
  - [ ] Minimize JS bundle for landing route (code splitting)
  - [ ] Add appropriate meta tags for SEO

- [ ] Task 6: Write Unit Tests (AC: all)
  - [ ] Test landing page renders all key sections
  - [ ] Test drag-drop zone interaction states
  - [ ] Test privacy badge displays correctly
  - [ ] Test demo button triggers sample data flow
  - [ ] Test navigation links render correctly
  - [ ] Test no authentication gates are present
  - [ ] Target: >80% coverage for new/modified components

- [ ] Task 7: Verify CI passes (AC-CI)
  - [ ] `pnpm nx affected --target=lint --base=main` passes
  - [ ] `pnpm nx affected --target=test --base=main` passes
  - [ ] `pnpm nx affected --target=type-check --base=main` passes

## Dev Notes

### Architecture Patterns & Constraints

**File Locations:**
- Landing page route: `apps/web/app/routes/index.tsx`
- Landing components: `apps/web/app/components/landing/`
- Shared upload component: `apps/web/app/components/common/FileUpload.tsx`
- Root layout/nav: `apps/web/app/routes/__root.tsx`

**Existing Components to Reuse/Extend:**
- `HeroSection.tsx` - Currently has basic hero, needs redesign
- `FeaturesSection.tsx` - Currently has feature highlights, needs update
- `Footer.tsx` - Currently exists, needs update
- `FileUpload.tsx` - Full drag-drop component with progress tracking, extend for landing zone
- `SampleResults.tsx` - Currently exists, may reuse for demo flow
- `EmailSignup.tsx` - Evaluate if needed or remove

**Styling:**
- Use Tailwind CSS with `cn()` utility from `apps/web/app/lib/utils.ts`
- Follow existing patterns: utility-first, responsive with Tailwind breakpoints
- Desktop-first design, tablet-friendly
- Icons from `lucide-react` (already installed)

**Key Dependencies Already Available:**
- `@radix-ui/*` - Accessible UI primitives
- `lucide-react` - Icons
- `tailwind-merge` + `clsx` - Classname merging
- `class-variance-authority` - Component variants

### Project Structure Notes

- Route files use TanStack Router file-based routing (`createFileRoute`)
- Components follow PascalCase.tsx naming
- Tests go in `__tests__/` directory adjacent to source
- Path aliases: `@/components/*`, `@/lib/*`, `@/hooks/*`

### UX Design Requirements (from UX Spec)

- Landing page should feature a central, large drag-drop zone (Vercel-style)
- Value prop: "Drop your project folder here to analyze code-level dependencies"
- Secondary CTA: "Or try a sample monorepo" button
- Privacy emphasis: "100% local analysis" badge prominently displayed
- Zero-configuration experience: drag-and-drop → immediate results
- Visual hierarchy: Drop zone > Value prop > Features > Footer

### Previous Story Intelligence

**From Epic 4 Stories:**
- Tailwind CSS patterns well established throughout Epic 4
- React.memo pattern used for performance-critical components
- Component testing uses `@testing-library/react` with `vitest`
- `vi.mock()` pattern for mocking stores and external dependencies
- eslint.config.mjs may need updates for new globals

**From Story 4.9:**
- Settings store (`apps/web/app/stores/settings.ts`) already exists with Zustand + devtools + persist
- Navigation between components/routes is well-established
- Test patterns: render, screen, fireEvent, waitFor from Testing Library

### Git Intelligence

Recent commits follow pattern: `feat(scope): description` or `fix(scope): description`
- Scope for Epic 5 should be `ui` or `web`
- Example: `feat(ui): implement landing page with drag-drop zone`

### Critical Don't-Miss Rules (from project-context.md)

1. **TanStack Start SSG only** - No SSR, no getServerSideProps, all client-side
2. **Tailwind CSS JIT mode** - Already configured
3. **No localStorage for large data** - Use IndexedDB via Dexie.js for analysis results
4. **File naming**: PascalCase.tsx for components, camelCase.ts for utilities
5. **Imports**: Use `@/` path aliases, avoid deep relative paths
6. **Zero backend architecture** - All analysis happens client-side via WASM

### Testing Requirements

- **Framework**: Vitest + @testing-library/react
- **Location**: `apps/web/app/components/landing/__tests__/`
- **Coverage target**: >80%
- **Test patterns**: Render component, assert key elements, simulate user interactions
- **Mock patterns**: vi.mock for router navigation, sample data

### References

- [Epic 5 Requirements] `_bmad-output/planning-artifacts/epics.md#epic-5-stories`
- [UX Spec: Landing Page] `_bmad-output/planning-artifacts/ux-design-specification.md` - Landing page section
- [Architecture: TanStack Start] `_bmad-output/planning-artifacts/architecture.md` - Web UI section
- [Project Context] `_bmad-output/project-context.md` - Framework rules, testing rules
- [FR28: Drag-and-drop upload] `_bmad-output/planning-artifacts/epics.md#fr-coverage-map`
- [FR33: No registration required] `_bmad-output/planning-artifacts/epics.md#fr-coverage-map`

## Dev Agent Record

### Agent Model Used

### Debug Log References

### Completion Notes List

### File List
