# Story 5.2: Implement File Upload and Workspace Detection

Status: ready-for-dev

## Story

As a **user**,
I want **to upload my workspace files and have the app detect the structure**,
So that **I don't need to configure anything manually**.

## Acceptance Criteria

### AC1: File Acceptance
**Given** the drag-drop zone or file picker
**When** I provide files
**Then** the app accepts:
- `package.json` files (individual or multiple)
- `pnpm-workspace.yaml` files
- `yarn.lock` / `package-lock.json` / `pnpm-lock.yaml` (for additional context)
- Folder uploads containing nested `package.json` files (via directory picker)
- `.zip` archives containing workspace files

### AC2: Workspace Type Auto-Detection
**Given** uploaded workspace files
**When** files are processed
**Then** the app auto-detects workspace type:
- **npm workspaces**: `package.json` with `workspaces` array field
- **yarn workspaces**: `package.json` with `workspaces` field (array or object with `packages`)
- **pnpm workspaces**: `pnpm-workspace.yaml` with `packages` field
- Detection result is displayed to user (e.g., "Detected: pnpm workspace with 12 packages")

### AC3: File Validation
**Given** files are dropped/selected
**When** the app processes the upload
**Then**:
- Invalid files show clear error message listing supported types
- Files exceeding 10MB show size limit error
- Malformed JSON shows parse error with helpful message
- Missing required fields (e.g., `workspaces`) show specific guidance
- Validation happens client-side only (no network requests)

### AC4: Upload Progress
**Given** files are being processed
**When** multiple files or a large folder is uploaded
**Then**:
- Progress indicator shows files processed / total files
- Processing status updates: "Reading files..." → "Detecting workspace..." → "Parsing packages..."
- Cancel button available during processing
- Total package count displayed when detection completes

### AC5: Auto-Start Analysis
**Given** workspace files are successfully uploaded and validated
**When** workspace detection completes
**Then**:
- Analysis starts automatically (no manual trigger needed)
- Transition to analysis view is smooth (no page flash)
- Uploaded data is stored in application state (Zustand store)
- User can return to upload and drop different files

### AC6: Multiple File Selection
**Given** the file upload interface
**When** I want to upload multiple files
**Then**:
- Multiple file selection is supported in the file picker dialog
- Drag-drop accepts multiple files at once
- Files are processed together as a single workspace
- Duplicate file names are handled (show warning)

### AC-CI: CI Pipeline Must Pass (MANDATORY)
**Given** the story implementation is complete
**When** verifying CI status
**Then** ALL of the following must pass:
- [ ] `pnpm nx affected --target=lint --base=main` passes
- [ ] `pnpm nx affected --target=test --base=main` passes
- [ ] `pnpm nx affected --target=type-check --base=main` passes
- **Story CANNOT be marked as "done" until CI is green**

## Tasks / Subtasks

- [ ] Task 1: Create Workspace Detection Service (AC: 2)
  - [ ] Create `apps/web/app/lib/workspaceDetection.ts`
  - [ ] Implement npm workspace detection (parse `workspaces` field from package.json)
  - [ ] Implement yarn workspace detection (array or object format)
  - [ ] Implement pnpm workspace detection (parse `pnpm-workspace.yaml`)
  - [ ] Return typed `WorkspaceDetectionResult` with type, packages, root config

- [ ] Task 2: Create File Processing Pipeline (AC: 1, 3, 4)
  - [ ] Create `apps/web/app/lib/fileProcessing.ts`
  - [ ] Implement file reader for JSON files (with error handling)
  - [ ] Implement YAML parser for pnpm-workspace.yaml (use js-yaml or similar)
  - [ ] Implement ZIP extraction (use JSZip or similar lightweight library)
  - [ ] Implement directory traversal for folder uploads
  - [ ] Add file size validation (10MB limit)
  - [ ] Add JSON parse validation with helpful error messages

- [ ] Task 3: Enhance Upload Component (AC: 1, 3, 4, 6)
  - [ ] Extend `apps/web/app/components/common/FileUpload.tsx` for workspace-specific upload
  - [ ] Or create `apps/web/app/components/upload/WorkspaceUpload.tsx`
  - [ ] Add progress state display (Reading → Detecting → Parsing)
  - [ ] Add cancel functionality during processing
  - [ ] Add workspace type detection result display
  - [ ] Add package count display on completion
  - [ ] Support multiple file selection via `<input multiple>`
  - [ ] Support directory upload via `<input webkitdirectory>`

- [ ] Task 4: Update Analyze Route (AC: 5)
  - [ ] Refactor `apps/web/app/routes/analyze.tsx` to handle uploaded file data
  - [ ] Accept data from landing page navigation (URL state or store)
  - [ ] Auto-trigger WASM analysis when valid workspace is detected
  - [ ] Show processing pipeline status

- [ ] Task 5: Connect Upload Flow to State (AC: 5)
  - [ ] Store uploaded workspace data in Zustand analysis store (create if needed, or use existing state)
  - [ ] Ensure data persists during route navigation (landing → analyze → results)
  - [ ] Allow re-upload by clearing state and returning to upload

- [ ] Task 6: Write Unit Tests (AC: all)
  - [ ] Test workspace detection for npm workspaces
  - [ ] Test workspace detection for yarn workspaces
  - [ ] Test workspace detection for pnpm workspaces
  - [ ] Test file validation (size, type, malformed JSON)
  - [ ] Test upload component interaction states
  - [ ] Test progress state transitions
  - [ ] Test multiple file handling
  - [ ] Target: >80% coverage

- [ ] Task 7: Verify CI passes (AC-CI)
  - [ ] All lint, test, type-check targets pass

## Dev Notes

### Architecture Patterns & Constraints

**File Locations:**
- Upload route: `apps/web/app/routes/analyze.tsx`
- Upload component: `apps/web/app/components/common/FileUpload.tsx` (extend or create new)
- Workspace detection: `apps/web/app/lib/workspaceDetection.ts` (new)
- File processing: `apps/web/app/lib/fileProcessing.ts` (new)
- Upload hook: `apps/web/app/hooks/useFileUpload.ts` (existing, may extend)
- Drag-drop hook: `apps/web/app/hooks/useDragAndDrop.ts` (existing)

**Existing Code to Leverage:**
- `FileUpload.tsx` already has drag-drop, progress tracking, error handling
- `useDragAndDrop.ts` hook handles browser drag-drop API
- `useFileUpload.ts` hook manages upload state and progress
- `analysis.ts` service has upload function (but uses API - need to adapt for local processing)

**Critical Architecture Rule:**
- ALL file processing happens client-side only
- NO network requests for file upload or analysis
- Files are read using browser File API / FileReader
- Zero backend architecture - code never leaves the browser

**New Dependencies May Be Needed:**
- `js-yaml` or `yaml` - for parsing pnpm-workspace.yaml (check if already available)
- `jszip` - for ZIP file extraction (evaluate necessity vs file size impact)

**Data Flow:**
```
Drop files → Read files (FileReader API) → Detect workspace type
→ Parse package.json files → Build WorkspaceData structure
→ Store in Zustand → Navigate to /analyze → Trigger WASM analysis
```

### Project Structure Notes

- File processing must be synchronous-looking but async under the hood (use async/await)
- WorkspaceData type is defined in `@monoguard/types` package
- All JSON parsing should use try/catch with user-friendly error messages
- Input validation at system boundary (file upload) per project-context.md security rules

### Previous Story Intelligence

**From Story 5.1:**
- Landing page drop zone feeds into this story's upload flow
- Navigation from landing to analyze route must carry file data
- Use same Tailwind patterns for consistent styling

**From Epic 2 (Analysis Engine):**
- Go WASM engine expects `WorkspaceData` type as input
- TypeScript WASM adapter (Story 2.7) provides `analyzer.analyze(workspaceData)`
- Result type is `Result<AnalysisResult>`

### Testing Requirements

- **Framework**: Vitest + @testing-library/react
- **Location**: `apps/web/app/lib/__tests__/` and `apps/web/app/components/upload/__tests__/`
- **Coverage target**: >80%
- **Critical test cases**: Each workspace type detection, file validation edge cases, error handling
- **Mock File API**: Use `new File([], 'name')` constructor for test files

### References

- [FR28: Drag-and-drop upload] `_bmad-output/planning-artifacts/epics.md`
- [FR29: Multi-file upload] `_bmad-output/planning-artifacts/epics.md`
- [FR5: npm/yarn/pnpm support] `_bmad-output/planning-artifacts/epics.md`
- [Project Context: Input Validation] `_bmad-output/project-context.md#security-rules`
- [Architecture: Zero Backend] `_bmad-output/project-context.md#critical-dont-miss-rules`
- [Story 2.7: TypeScript WASM Adapter] `_bmad-output/implementation-artifacts/2-7-create-typescript-wasm-adapter.md`

## Dev Agent Record

### Agent Model Used

### Debug Log References

### Completion Notes List

### File List
