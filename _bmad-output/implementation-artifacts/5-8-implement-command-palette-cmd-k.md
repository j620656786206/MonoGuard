# Story 5.8: Implement Command Palette (Cmd+K)

Status: ready-for-dev

## Story

As a **power user**,
I want **a command palette for quick navigation and actions**,
So that **I can work efficiently with keyboard shortcuts**.

## Acceptance Criteria

### AC1: Keyboard Activation
**Given** I'm using the web app on any page
**When** I press Cmd+K (Mac) or Ctrl+K (Windows/Linux)
**Then**:
- Command palette modal opens immediately (< 100ms)
- Search input is auto-focused
- Background is dimmed with overlay
- Previous search is cleared on open

### AC2: Search Functionality
**Given** the command palette is open
**When** I type in the search input
**Then**:
- Results update as I type (< 50ms per keystroke)
- Fuzzy search matching (e.g., "dk" matches "Dark mode")
- Search matches against command name and description
- Results are grouped by category (Navigation, Actions, Settings)
- Maximum 10 results visible (scrollable if more)

### AC3: Available Commands
**Given** the command palette is open
**When** I view available commands
**Then** I see commands grouped by category:
- **Navigation**: Go to Home, Go to Analyze, Go to Results
- **Actions**: Upload Files, Export Report, Run Analysis, Clear Results
- **Settings**: Toggle Dark Mode, Change Visualization Mode, Reset Settings
- **Help**: View Keyboard Shortcuts, About MonoGuard
- Commands show icons (from lucide-react) and optional keyboard shortcut hints

### AC4: Keyboard Navigation
**Given** the command palette is open with results
**When** I navigate with keyboard
**Then**:
- Arrow Up/Down: Navigate between items
- Enter: Execute selected command
- Escape: Close palette
- Tab: Move to next item (same as Arrow Down)
- Current selection is visually highlighted
- Selection wraps around (last → first, first → last)

### AC5: Command Execution
**Given** a command is selected in the palette
**When** I press Enter or click it
**Then**:
- Command executes immediately
- Palette closes after execution
- Navigation commands navigate to the route
- Action commands trigger the action
- Settings commands toggle the setting
- Visual feedback confirms execution (toast notification from Story 5.10)

### AC6: Recent Commands
**Given** I've used the command palette before
**When** I open it with no search text
**Then**:
- Recent commands shown first (last 5 used)
- "Recent" section header displayed
- Full command list available below
- Recent list persisted across sessions

### AC7: Contextual Commands
**Given** I'm on a specific page
**When** I open the command palette
**Then**:
- Page-specific commands appear first
- On Results page: "Export Report", "Toggle Side Panel", "Zoom to Fit"
- On Analyze page: "Upload Files", "Run with Example"
- On Landing page: "Try Demo", "Upload Files"
- Context-specific commands have a "Current Page" badge

### AC-CI: CI Pipeline Must Pass (MANDATORY)
**Given** the story implementation is complete
**When** verifying CI status
**Then** ALL of the following must pass:
- [ ] `pnpm nx affected --target=lint --base=main` passes
- [ ] `pnpm nx affected --target=test --base=main` passes
- [ ] `pnpm nx affected --target=type-check --base=main` passes
- **Story CANNOT be marked as "done" until CI is green**

## Tasks / Subtasks

- [ ] Task 1: Create Command Registry (AC: 3, 7)
  - [ ] Create `apps/web/app/lib/commandPalette/commands.ts`
  - [ ] Define `Command` interface: id, name, description, icon, category, action, shortcut, contextPages
  - [ ] Register navigation commands (Home, Analyze, Results)
  - [ ] Register action commands (Upload, Export, Analyze, Clear)
  - [ ] Register settings commands (Dark Mode, Visualization Mode)
  - [ ] Register help commands
  - [ ] Add context filtering based on current route

- [ ] Task 2: Create Fuzzy Search (AC: 2)
  - [ ] Create `apps/web/app/lib/commandPalette/fuzzySearch.ts`
  - [ ] Implement fuzzy matching algorithm (simple substring + scoring)
  - [ ] Match against command name and description
  - [ ] Return sorted results by relevance score
  - [ ] Performance: handle search in < 50ms for 50+ commands

- [ ] Task 3: Create Command Palette Component (AC: 1, 2, 4)
  - [ ] Create `apps/web/app/components/common/CommandPalette.tsx`
  - [ ] Use Radix UI Dialog as base (with custom styling)
  - [ ] Search input with search icon
  - [ ] Scrollable results list grouped by category
  - [ ] Highlight matching characters in results
  - [ ] Selected item highlight state
  - [ ] Close on Escape, click outside, or command execution

- [ ] Task 4: Implement Keyboard Navigation (AC: 1, 4)
  - [ ] Create `apps/web/app/hooks/useCommandPalette.ts`
  - [ ] Listen for Cmd+K / Ctrl+K globally
  - [ ] Handle Arrow Up/Down/Enter/Escape within palette
  - [ ] Manage selected index state
  - [ ] Prevent default browser behavior for Cmd+K
  - [ ] Handle selection wrapping

- [ ] Task 5: Implement Recent Commands (AC: 6)
  - [ ] Track last 5 used commands in settings store (Zustand persist)
  - [ ] Display recent section when search is empty
  - [ ] Update recent list after each command execution

- [ ] Task 6: Implement Command Execution (AC: 5)
  - [ ] Wire navigation commands to TanStack Router navigation
  - [ ] Wire action commands to appropriate handlers (store actions)
  - [ ] Wire settings commands to settings store updates
  - [ ] Close palette after execution
  - [ ] Trigger toast notification for settings changes

- [ ] Task 7: Add Palette to Root Layout (AC: 1)
  - [ ] Add CommandPalette to `apps/web/app/routes/__root.tsx`
  - [ ] Ensure keyboard listener is global (works on any page)
  - [ ] Palette overlays all content

- [ ] Task 8: Write Unit Tests (AC: all)
  - [ ] Test keyboard activation (Cmd+K opens palette)
  - [ ] Test fuzzy search algorithm (matching, scoring, performance)
  - [ ] Test keyboard navigation (arrows, enter, escape)
  - [ ] Test command execution (navigation, actions, settings)
  - [ ] Test recent commands tracking and persistence
  - [ ] Test contextual commands based on route
  - [ ] Test search result grouping by category
  - [ ] Target: >80% coverage

- [ ] Task 9: Verify CI passes (AC-CI)
  - [ ] All lint, test, type-check targets pass

## Dev Notes

### Architecture Patterns & Constraints

**File Locations:**
- Command registry: `apps/web/app/lib/commandPalette/commands.ts` (new)
- Fuzzy search: `apps/web/app/lib/commandPalette/fuzzySearch.ts` (new)
- Palette component: `apps/web/app/components/common/CommandPalette.tsx` (new)
- Palette hook: `apps/web/app/hooks/useCommandPalette.ts` (new)
- Root layout: `apps/web/app/routes/__root.tsx` (add palette)
- Settings store: `apps/web/app/stores/settings.ts` (add recentCommands)

**Command Interface:**
```typescript
interface Command {
  id: string;
  name: string;
  description: string;
  icon: LucideIcon;
  category: 'navigation' | 'actions' | 'settings' | 'help';
  action: () => void | Promise<void>;
  shortcut?: string; // e.g., "Cmd+D" for dark mode
  contextPages?: string[]; // Routes where this is contextual
}
```

**Radix UI Dialog Pattern:**
```typescript
import * as Dialog from '@radix-ui/react-dialog';

// CommandPalette renders as Dialog overlay
<Dialog.Root open={isOpen} onOpenChange={setIsOpen}>
  <Dialog.Portal>
    <Dialog.Overlay className="fixed inset-0 bg-black/50" />
    <Dialog.Content className="fixed top-1/4 left-1/2 -translate-x-1/2 ...">
      <input type="text" placeholder="Type a command..." />
      {/* Results list */}
    </Dialog.Content>
  </Dialog.Portal>
</Dialog.Root>
```

**Global Keyboard Listener Pattern:**
```typescript
useEffect(() => {
  const handleKeyDown = (e: KeyboardEvent) => {
    if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
      e.preventDefault();
      setIsOpen(true);
    }
  };
  document.addEventListener('keydown', handleKeyDown);
  return () => document.removeEventListener('keydown', handleKeyDown);
}, []);
```

### UX Design Requirements (from UX Spec)

- Linear-inspired command palette design
- Fuzzy search with highlighted matching characters
- Grouped results by category
- Keyboard-first interaction (mouse also supported)
- Quick execution: open → type → enter → done

### Previous Story Intelligence

**From Story 5.7 (Dark Mode):**
- Toggle dark mode command should call `settings.setTheme()`
- Theme toggle is also available in nav bar

**From Story 5.4 (Dashboard):**
- "Export Report" command should open ExportDialog
- "Toggle Side Panel" command for results page

**From Story 5.10 (Toast):**
- Command execution can trigger toast notifications for feedback

### Dependencies

- `@radix-ui/react-dialog` - Already in package.json
- `lucide-react` - Already in package.json for icons
- TanStack Router `useNavigate()` - For navigation commands

### Testing Requirements

- Mock keyboard events for activation tests
- Test fuzzy search with various input patterns
- Test command execution callbacks
- Mock TanStack Router for navigation testing
- Test recent commands persistence via store

### References

- [UX Spec: Command Palette (Cmd+K)] `_bmad-output/planning-artifacts/ux-design-specification.md#command-palette`
- [UX Spec: Navigation Patterns] `_bmad-output/planning-artifacts/ux-design-specification.md`
- [Radix UI Dialog] Already in dependencies
- [Project Context: React Hooks] `_bmad-output/project-context.md#react-19`

## Dev Agent Record

### Agent Model Used

### Debug Log References

### Completion Notes List

### File List
