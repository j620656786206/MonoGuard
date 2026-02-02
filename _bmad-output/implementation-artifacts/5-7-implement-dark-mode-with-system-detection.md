# Story 5.7: Implement Dark Mode with System Detection

Status: ready-for-dev

## Story

As a **user**,
I want **the app to support dark mode**,
So that **I can use it comfortably in low-light environments**.

## Acceptance Criteria

### AC1: System Preference Detection
**Given** I visit the MonoGuard web app
**When** the app loads for the first time
**Then**:
- App detects system color scheme preference (`prefers-color-scheme`)
- If system is dark → app renders in dark mode
- If system is light → app renders in light mode
- If no preference → defaults to light mode

### AC2: Manual Toggle
**Given** the web app
**When** I want to change the theme
**Then**:
- Theme toggle available in navigation bar (sun/moon icon)
- Three options: Light / Dark / System (auto-detect)
- Toggle is accessible via keyboard
- Selection is clearly indicated (active state)

### AC3: Theme Persistence
**Given** I've selected a theme preference
**When** I reload the page or return later
**Then**:
- My preference is persisted (Zustand persist → localStorage)
- No flash of wrong theme on page load (FOUC prevention)
- Preference survives browser cache clears of other data
- "System" option continues to follow system changes in real-time

### AC4: Complete Dark Mode Coverage
**Given** dark mode is active
**When** I navigate the entire app
**Then** all elements have proper dark variants:
- Landing page (hero, features, footer)
- Upload zone (borders, backgrounds, text)
- Analysis progress indicators
- Results dashboard (health score, cards, lists)
- Dependency graph visualization (node/edge colors)
- Side panel (fix suggestions)
- Navigation bar and footer
- Dialogs/modals (export, diagnostics)
- Toast notifications
- Command palette (Story 5.8)

### AC5: Graph Visualization Dark Mode
**Given** the dependency graph in dark mode
**When** it renders
**Then**:
- Background is dark (`bg-gray-900` or similar)
- Node colors adapt for visibility on dark backgrounds
- Edge colors have sufficient contrast
- Circular dependency red highlighting still visible
- Labels are readable (light text on dark)
- Legend colors adapt accordingly
- Both SVG and Canvas renderers support dark mode

### AC6: Smooth Transition
**Given** I toggle between light and dark mode
**When** the theme changes
**Then**:
- Transition is smooth (< 200ms)
- CSS transition applied to background and text colors
- No jarring flashes or layout shifts
- Graph re-renders cleanly without artifacts

### AC-CI: CI Pipeline Must Pass (MANDATORY)
**Given** the story implementation is complete
**When** verifying CI status
**Then** ALL of the following must pass:
- [ ] `pnpm nx affected --target=lint --base=main` passes
- [ ] `pnpm nx affected --target=test --base=main` passes
- [ ] `pnpm nx affected --target=type-check --base=main` passes
- **Story CANNOT be marked as "done" until CI is green**

## Tasks / Subtasks

- [ ] Task 1: Extend Settings Store for Theme (AC: 2, 3)
  - [ ] Update `apps/web/app/stores/settings.ts`
  - [ ] Add `theme: 'light' | 'dark' | 'system'` to SettingsState (field already exists, verify)
  - [ ] Add `setTheme` action (already exists, verify)
  - [ ] Ensure persist middleware saves theme preference
  - [ ] Default: 'system'

- [ ] Task 2: Create Theme Provider/Hook (AC: 1, 3, 6)
  - [ ] Create `apps/web/app/hooks/useTheme.ts`
  - [ ] Implement system preference detection via `window.matchMedia('(prefers-color-scheme: dark)')`
  - [ ] Listen for system changes with `matchMedia.addEventListener('change', ...)`
  - [ ] Resolve effective theme: if 'system' → check media query; else use stored value
  - [ ] Apply `dark` class to `<html>` element (Tailwind dark mode strategy)
  - [ ] Handle FOUC: add inline script in `index.html` to set theme class before React hydrates

- [ ] Task 3: Create Theme Toggle Component (AC: 2)
  - [ ] Create `apps/web/app/components/common/ThemeToggle.tsx`
  - [ ] Sun icon (light), Moon icon (dark), Monitor icon (system) from lucide-react
  - [ ] Dropdown or cycle toggle (click to cycle through: light → dark → system)
  - [ ] Show current selection with visual indicator
  - [ ] Keyboard accessible (Enter/Space to toggle)
  - [ ] ARIA labels for accessibility

- [ ] Task 4: Add FOUC Prevention Script (AC: 3)
  - [ ] Add inline `<script>` in `apps/web/index.html` before React bundle
  - [ ] Script reads localStorage for theme preference
  - [ ] Applies `dark` class to `<html>` immediately
  - [ ] Must execute synchronously before first paint

- [ ] Task 5: Add Tailwind Dark Mode Classes (AC: 4)
  - [ ] Verify `tailwind.config.ts` has `darkMode: 'class'` (or add it)
  - [ ] Audit and add `dark:` variants to all existing components:
    - [ ] Landing page components (`apps/web/app/components/landing/`)
    - [ ] Dashboard components (`apps/web/app/components/dashboard/`)
    - [ ] Upload components
    - [ ] Analysis components (`apps/web/app/components/analysis/`)
    - [ ] Navigation and footer in root layout
  - [ ] Add dark mode backgrounds: `dark:bg-gray-900`, `dark:bg-gray-800`
  - [ ] Add dark mode text: `dark:text-gray-100`, `dark:text-gray-300`
  - [ ] Add dark mode borders: `dark:border-gray-700`

- [ ] Task 6: Update Graph Visualization for Dark Mode (AC: 5)
  - [ ] Update SVGRenderer colors to be theme-aware
  - [ ] Update CanvasRenderer colors to be theme-aware
  - [ ] Read current theme from settings store or CSS variables
  - [ ] Adapt: node fill, node stroke, edge stroke, label text color
  - [ ] Ensure circular dependency red highlighting maintains contrast
  - [ ] Update GraphLegend colors for dark mode
  - [ ] Update NodeTooltip styling for dark mode

- [ ] Task 7: Add Theme Toggle to Navigation (AC: 2)
  - [ ] Add ThemeToggle component to nav bar in `apps/web/app/routes/__root.tsx`
  - [ ] Position in right side of navigation
  - [ ] Consistent with overall nav design

- [ ] Task 8: Write Unit Tests (AC: all)
  - [ ] Test system preference detection (mock matchMedia)
  - [ ] Test manual theme toggle cycles through modes
  - [ ] Test theme persistence in store
  - [ ] Test dark class applied to HTML element
  - [ ] Test FOUC prevention (inline script behavior)
  - [ ] Test ThemeToggle component rendering and interaction
  - [ ] Test graph color changes in dark mode
  - [ ] Target: >80% coverage

- [ ] Task 9: Verify CI passes (AC-CI)
  - [ ] All lint, test, type-check targets pass

## Dev Notes

### Architecture Patterns & Constraints

**File Locations:**
- Settings store: `apps/web/app/stores/settings.ts` (extend - theme field may already exist)
- Theme hook: `apps/web/app/hooks/useTheme.ts` (new)
- Theme toggle: `apps/web/app/components/common/ThemeToggle.tsx` (new)
- Root layout: `apps/web/app/routes/__root.tsx` (add toggle to nav)
- Entry HTML: `apps/web/index.html` (add FOUC prevention script)
- Tailwind config: `apps/web/tailwind.config.ts` (verify darkMode setting)

**Tailwind Dark Mode Strategy:**
```javascript
// tailwind.config.ts
export default {
  darkMode: 'class', // Use class-based dark mode
  // ...
}
```

**FOUC Prevention Script (in index.html):**
```html
<script>
  (function() {
    try {
      var stored = JSON.parse(localStorage.getItem('monoguard-settings') || '{}');
      var theme = stored.state && stored.state.theme;
      if (theme === 'dark' || (theme === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
        document.documentElement.classList.add('dark');
      }
    } catch(e) {}
  })();
</script>
```

**Settings Store (existing pattern):**
The settings store already has `theme` field declared in the interface from Story 4.9. Verify it's functional:
```typescript
interface SettingsState {
  theme: 'light' | 'dark' | 'system';  // Already declared
  visualizationMode: RenderModePreference;
  setTheme: (theme: SettingsState['theme']) => void;  // Already declared
  setVisualizationMode: (mode: RenderModePreference) => void;
}
```

**Graph Dark Mode Colors:**
```typescript
// Theme-aware colors for graph rendering
const GRAPH_COLORS = {
  light: {
    nodeFill: '#4f46e5',
    nodeStroke: '#fff',
    edgeStroke: '#9ca3af',
    labelText: '#1f2937',
    circularFill: '#fecaca',
    circularStroke: '#ef4444',
  },
  dark: {
    nodeFill: '#818cf8',
    nodeStroke: '#1f2937',
    edgeStroke: '#4b5563',
    labelText: '#e5e7eb',
    circularFill: '#7f1d1d',
    circularStroke: '#f87171',
  },
};
```

### Previous Story Intelligence

**From Story 4.9:**
- Settings store created with Zustand + devtools + persist
- `theme` and `setTheme` fields already in interface
- Store key: `monoguard-settings`
- RenderModeIndicator already has `dark:` Tailwind classes (pattern reference)

**From Epic 4 Visualization:**
- SVGRenderer and CanvasRenderer use hardcoded colors currently
- Need to make colors theme-aware (read from store or CSS variables)
- Both renderers already wrapped in React.memo

### UX Requirements

- System detection + manual toggle (three modes)
- Smooth transition < 200ms
- No FOUC (flash of unstyled/wrong-theme content)
- Graph visualization fully adapted

### Testing Requirements

- Mock `window.matchMedia` for system preference tests
- Mock localStorage for persistence tests
- Test HTML element class toggling
- Use `@testing-library/user-event` for toggle interaction
- Snapshot tests for dark vs light mode may be helpful

### References

- [UX Spec: Dark Mode] `_bmad-output/planning-artifacts/ux-design-specification.md`
- [Story 4.9: Settings Store] `_bmad-output/implementation-artifacts/4-9-*.md`
- [Project Context: Tailwind CSS JIT] `_bmad-output/project-context.md#framework-specific-rules`
- [Zustand Persist Pattern] `_bmad-output/project-context.md#zustand-state-management`

## Dev Agent Record

### Agent Model Used

### Debug Log References

### Completion Notes List

### File List
