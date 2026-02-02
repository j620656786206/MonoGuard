# Story 5.10: Implement Toast Notifications

Status: ready-for-dev

## Story

As a **user**,
I want **non-blocking notifications for actions and events**,
So that **I stay informed without interrupting my workflow**.

## Acceptance Criteria

### AC1: Toast Display
**Given** various user actions or system events
**When** a notification is triggered
**Then** a toast notification:
- Appears in the bottom-right corner of the screen
- Slides in with smooth animation (from right or bottom)
- Does not block any UI interaction (non-modal)
- Has a maximum width to prevent layout overflow

### AC2: Auto-Dismiss
**Given** a toast notification is displayed
**When** the display duration expires
**Then**:
- Toast auto-dismisses after 5 seconds by default
- Dismissal uses a smooth slide-out animation
- Hovering over the toast pauses the auto-dismiss timer
- Un-hovering resumes the timer

### AC3: Manual Dismiss
**Given** a toast notification is displayed
**When** I want to dismiss it early
**Then**:
- Close (X) button visible on the toast
- Clicking close dismisses immediately with animation
- Swipe gesture dismisses on touch devices (optional enhancement)

### AC4: Toast Variants
**Given** different types of notifications
**When** toasts are created
**Then** they have distinct visual styles:
- **Success**: Green background/border, checkmark icon (e.g., "Analysis complete!")
- **Error**: Red background/border, X icon (e.g., "Failed to load WASM module")
- **Warning**: Yellow/amber background/border, warning icon (e.g., "Large workspace detected")
- **Info**: Blue background/border, info icon (e.g., "Report downloaded successfully")

### AC5: Toast Stacking
**Given** multiple notifications triggered in quick succession
**When** multiple toasts are visible
**Then**:
- Toasts stack vertically (newest at bottom)
- Maximum 3 visible toasts at once
- Older toasts dismissed when exceeding max
- Stacked toasts don't overlap or shift unexpectedly
- Gap between stacked toasts is consistent

### AC6: Toast Content
**Given** a toast notification
**When** it renders
**Then**:
- Title text (bold, short - e.g., "Export Complete")
- Optional description text (smaller, 1-2 lines - e.g., "Report saved as monoguard-report.json")
- Optional action button (e.g., "Undo", "View", "Retry")
- Icon matching the variant type

### AC7: Programmatic Toast API
**Given** components or services need to trigger notifications
**When** using the toast API
**Then**:
- Simple function call: `toast.success('Analysis complete!')`
- Full options: `toast({ title, description, variant, duration, action })`
- Available globally without prop drilling
- Type-safe API with TypeScript

### AC8: Dark Mode Support
**Given** dark mode is active
**When** a toast notification appears
**Then**:
- Toast adapts to dark theme colors
- Sufficient contrast for readability
- Icons maintain visibility
- Consistent with overall dark mode design

### AC-CI: CI Pipeline Must Pass (MANDATORY)
**Given** the story implementation is complete
**When** verifying CI status
**Then** ALL of the following must pass:
- [ ] `pnpm nx affected --target=lint --base=main` passes
- [ ] `pnpm nx affected --target=test --base=main` passes
- [ ] `pnpm nx affected --target=type-check --base=main` passes
- **Story CANNOT be marked as "done" until CI is green**

## Tasks / Subtasks

- [ ] Task 1: Create Toast Infrastructure (AC: 1, 7)
  - [ ] Create `apps/web/app/components/common/Toast/ToastProvider.tsx`
  - [ ] Create `apps/web/app/components/common/Toast/Toast.tsx`
  - [ ] Create `apps/web/app/components/common/Toast/index.ts` (barrel export)
  - [ ] Use Radix UI Toast primitives (`@radix-ui/react-toast` already in dependencies)
  - [ ] Set up ToastProvider at app root level

- [ ] Task 2: Create Toast Store / Context (AC: 7)
  - [ ] Create `apps/web/app/lib/toast.ts` (toast API)
  - [ ] Implement global `toast` function: `toast.success()`, `toast.error()`, `toast.warning()`, `toast.info()`
  - [ ] Full API: `toast({ title, description, variant, duration, action })`
  - [ ] Use Zustand micro-store or React context for toast queue management
  - [ ] Export `toast` function for global use without prop drilling

- [ ] Task 3: Implement Toast Variants (AC: 4, 6)
  - [ ] Create variant styles for success, error, warning, info
  - [ ] Use `class-variance-authority` for variant management (already in dependencies)
  - [ ] Add icons from `lucide-react`: CheckCircle, XCircle, AlertTriangle, Info
  - [ ] Support optional title, description, and action button

- [ ] Task 4: Implement Animation & Timing (AC: 1, 2, 3)
  - [ ] Slide-in animation from right side
  - [ ] Slide-out animation on dismiss
  - [ ] Auto-dismiss after 5 seconds (default, configurable per toast)
  - [ ] Pause timer on hover, resume on un-hover
  - [ ] Close button with immediate dismiss
  - [ ] Use Tailwind CSS animations (`tailwindcss-animate` already installed)

- [ ] Task 5: Implement Toast Stacking (AC: 5)
  - [ ] Manage toast queue (max 3 visible)
  - [ ] Stack positioning (bottom-right, vertical stack)
  - [ ] Dismiss oldest when exceeding max
  - [ ] Smooth reflow when toast is removed from stack

- [ ] Task 6: Add Dark Mode Support (AC: 8)
  - [ ] Add `dark:` Tailwind variants to all toast styles
  - [ ] Ensure contrast ratio meets WCAG AA
  - [ ] Test all variants in dark mode

- [ ] Task 7: Wire Toast Provider to App (AC: 1, 7)
  - [ ] Add ToastProvider to `apps/web/app/routes/__root.tsx`
  - [ ] Ensure toast viewport is positioned correctly
  - [ ] Test toast works from any route/page

- [ ] Task 8: Add Toast Calls to Existing Features (AC: 7)
  - [ ] Analysis complete → `toast.success('Analysis complete!')`
  - [ ] Analysis error → `toast.error('Analysis failed', { description: error.userMessage })`
  - [ ] Report exported → `toast.success('Report downloaded', { description: filename })`
  - [ ] Theme changed → `toast.info('Theme updated to dark mode')`
  - [ ] Settings saved → `toast.success('Settings saved')`

- [ ] Task 9: Write Unit Tests (AC: all)
  - [ ] Test toast rendering for each variant
  - [ ] Test auto-dismiss timing
  - [ ] Test manual dismiss
  - [ ] Test hover pause behavior
  - [ ] Test toast stacking (max 3)
  - [ ] Test programmatic API (`toast.success()`, etc.)
  - [ ] Test dark mode variants
  - [ ] Test action button rendering and click
  - [ ] Target: >80% coverage

- [ ] Task 10: Verify CI passes (AC-CI)
  - [ ] All lint, test, type-check targets pass

## Dev Notes

### Architecture Patterns & Constraints

**File Locations:**
- Toast components: `apps/web/app/components/common/Toast/` (new directory)
- Toast API: `apps/web/app/lib/toast.ts` (new)
- Root layout: `apps/web/app/routes/__root.tsx` (add provider)

**Radix UI Toast (Already Available):**
`@radix-ui/react-toast` is already in the package.json dependencies.

```typescript
import * as Toast from '@radix-ui/react-toast';

// Provider wraps the app
<Toast.Provider swipeDirection="right">
  {children}
  <Toast.Viewport className="fixed bottom-0 right-0 flex flex-col gap-2 p-4 w-96 max-w-full z-50" />
</Toast.Provider>

// Individual toast
<Toast.Root duration={5000}>
  <Toast.Title>{title}</Toast.Title>
  <Toast.Description>{description}</Toast.Description>
  <Toast.Action altText="action">{actionLabel}</Toast.Action>
  <Toast.Close />
</Toast.Root>
```

**Global Toast API Pattern:**
```typescript
// apps/web/app/lib/toast.ts
import { create } from 'zustand';

interface ToastItem {
  id: string;
  title: string;
  description?: string;
  variant: 'success' | 'error' | 'warning' | 'info';
  duration?: number;
  action?: { label: string; onClick: () => void };
}

interface ToastStore {
  toasts: ToastItem[];
  addToast: (toast: Omit<ToastItem, 'id'>) => void;
  removeToast: (id: string) => void;
}

const useToastStore = create<ToastStore>((set) => ({
  toasts: [],
  addToast: (toast) => set((state) => ({
    toasts: [...state.toasts.slice(-2), { ...toast, id: crypto.randomUUID() }], // Max 3
  })),
  removeToast: (id) => set((state) => ({
    toasts: state.toasts.filter(t => t.id !== id),
  })),
}));

// Convenience functions
export const toast = {
  success: (title: string, opts?: Partial<ToastItem>) =>
    useToastStore.getState().addToast({ title, variant: 'success', ...opts }),
  error: (title: string, opts?: Partial<ToastItem>) =>
    useToastStore.getState().addToast({ title, variant: 'error', duration: 8000, ...opts }),
  warning: (title: string, opts?: Partial<ToastItem>) =>
    useToastStore.getState().addToast({ title, variant: 'warning', ...opts }),
  info: (title: string, opts?: Partial<ToastItem>) =>
    useToastStore.getState().addToast({ title, variant: 'info', ...opts }),
};
```

**Variant Styles with CVA:**
```typescript
import { cva } from 'class-variance-authority';

const toastVariants = cva(
  'rounded-lg border p-4 shadow-lg flex items-start gap-3',
  {
    variants: {
      variant: {
        success: 'bg-green-50 border-green-200 text-green-800 dark:bg-green-950 dark:border-green-800 dark:text-green-200',
        error: 'bg-red-50 border-red-200 text-red-800 dark:bg-red-950 dark:border-red-800 dark:text-red-200',
        warning: 'bg-amber-50 border-amber-200 text-amber-800 dark:bg-amber-950 dark:border-amber-800 dark:text-amber-200',
        info: 'bg-blue-50 border-blue-200 text-blue-800 dark:bg-blue-950 dark:border-blue-800 dark:text-blue-200',
      },
    },
    defaultVariants: {
      variant: 'info',
    },
  }
);
```

### UX Design Requirements (from UX Spec)

- Non-blocking notifications in bottom-right corner
- Progress percentage + estimated remaining time for long operations
- Success animation feedback on completion
- "Emotional design: eliminate anxiety, provide control"

### Previous Story Intelligence

**From Story 4.9:**
- `console.warn()` used for performance warnings → should be replaced with `toast.warning()`
- TODO comment in DependencyGraphViz: "Integrate with toast notification system when available"

**Integration Points Across Epic 5:**
- Story 5.1: Landing page demo load → toast.info
- Story 5.2: File upload errors → toast.error
- Story 5.3: Analysis complete/failed → toast.success/error
- Story 5.6: Report download → toast.success
- Story 5.7: Theme changed → toast.info
- Story 5.8: Command executed → toast feedback

### Dependencies Already Available

- `@radix-ui/react-toast` - Toast primitive
- `class-variance-authority` - Variant management
- `lucide-react` - Icons (CheckCircle, XCircle, AlertTriangle, Info)
- `tailwindcss-animate` - CSS animations
- `zustand` - Toast queue state management

### Testing Requirements

- Use `@testing-library/react` for rendering tests
- Test timer behavior with `vi.useFakeTimers()` and `vi.advanceTimersByTime()`
- Test hover pause with `fireEvent.mouseEnter` / `fireEvent.mouseLeave`
- Test stacking by adding multiple toasts
- Verify toast API calls trigger correct renders

### References

- [UX Spec: Toast Notifications] `_bmad-output/planning-artifacts/ux-design-specification.md#toast`
- [UX Spec: Feedback System] `_bmad-output/planning-artifacts/ux-design-specification.md`
- [Radix UI Toast] Already in dependencies
- [Project Context: React Patterns] `_bmad-output/project-context.md#react-19`
- [Story 4.9: TODO for toast integration] `_bmad-output/implementation-artifacts/4-9-*.md`

## Dev Agent Record

### Agent Model Used

### Debug Log References

### Completion Notes List

### File List
