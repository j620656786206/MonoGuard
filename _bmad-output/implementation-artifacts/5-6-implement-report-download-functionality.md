# Story 5.6: Implement Report Download Functionality

Status: ready-for-dev

## Story

As a **user**,
I want **to download analysis reports from the web interface**,
So that **I can share or archive my results**.

## Acceptance Criteria

### AC1: Export Format Selection
**Given** analysis results are available
**When** I click the download/export button
**Then** I can choose:
- **JSON** - Machine-readable, complete analysis data
- **HTML** - Standalone viewable report with embedded styles
- **Markdown** - For documentation/wikis
- Format selection via dropdown or modal dialog

### AC2: Content Selection
**Given** the export dialog
**When** I configure export options
**Then** I can choose:
- Full report (all sections) or summary only
- Include/exclude dependency graph image (PNG embedded in HTML)
- Include/exclude fix recommendations
- Include/exclude version conflict details

### AC3: Download Execution
**Given** I've selected format and options
**When** I click "Download"
**Then**:
- File downloads immediately (no server round-trip)
- Filename follows pattern: `monoguard-report-{project}-{date}.{ext}`
- Downloaded file is valid and complete
- Progress indicator for large reports (if generation takes > 500ms)

### AC4: Report Content Quality
**Given** a downloaded report
**When** I open it
**Then**:
- **JSON**: Valid JSON, matches `AnalysisResult` TypeScript type structure
- **HTML**: Self-contained (no external CSS/JS dependencies), opens correctly in browser
- **Markdown**: Renders correctly in GitHub/GitLab, includes tables and code blocks
- All reports include: health score, circular dependencies, version conflicts, metadata (timestamp, package count)

### AC5: Offline Export
**Given** the web app with analysis results
**When** I export a report with no network connection
**Then**:
- Export works fully offline
- No external resources referenced in exported files
- HTML report has embedded CSS (no CDN links)

### AC-CI: CI Pipeline Must Pass (MANDATORY)
**Given** the story implementation is complete
**When** verifying CI status
**Then** ALL of the following must pass:
- [ ] `pnpm nx affected --target=lint --base=main` passes
- [ ] `pnpm nx affected --target=test --base=main` passes
- [ ] `pnpm nx affected --target=type-check --base=main` passes
- **Story CANNOT be marked as "done" until CI is green**

## Tasks / Subtasks

- [ ] Task 1: Create Export Dialog Component (AC: 1, 2)
  - [ ] Create `apps/web/app/components/dashboard/ExportDialog.tsx`
  - [ ] Format selection (JSON, HTML, Markdown) with radio buttons or tabs
  - [ ] Content toggles (full/summary, include graph, include fixes)
  - [ ] Download button with loading state
  - [ ] Use Radix UI Dialog for the modal

- [ ] Task 2: Wire Existing Report Generation (AC: 3, 4)
  - [ ] Connect to existing `apps/web/app/lib/reports/` system
  - [ ] Use `generateJsonReport()`, `generateHtmlReport()`, `generateMarkdownReport()`
  - [ ] Pass analysis results from Zustand store to report generators
  - [ ] Transform `AnalysisResult` to `ReportData` using existing builders

- [ ] Task 3: Implement Client-Side Download (AC: 3, 5)
  - [ ] Create `apps/web/app/lib/downloadFile.ts` utility
  - [ ] Generate blob from report content
  - [ ] Create download link (`URL.createObjectURL` + `<a download>`)
  - [ ] Generate filename with project name and date (ISO 8601)
  - [ ] Clean up object URL after download

- [ ] Task 4: Add Graph Image Embedding (AC: 2, 4)
  - [ ] Use existing export utilities from `apps/web/app/components/visualization/DependencyGraph/utils/`
  - [ ] `exportPng.ts` to capture graph as PNG data URL
  - [ ] Embed PNG in HTML report as base64 `<img>` tag
  - [ ] For Markdown: include as inline base64 image or note "graph not included"

- [ ] Task 5: Add Export Button to Dashboard (AC: 1)
  - [ ] Add export/download button in results dashboard header
  - [ ] Button opens ExportDialog
  - [ ] Button disabled when no results available
  - [ ] Use lucide-react download icon

- [ ] Task 6: Write Unit Tests (AC: all)
  - [ ] Test ExportDialog rendering and format selection
  - [ ] Test content toggle options
  - [ ] Test download file utility (blob creation, filename generation)
  - [ ] Test JSON report structure matches expected types
  - [ ] Test HTML report is self-contained (no external URLs)
  - [ ] Test Markdown report formatting
  - [ ] Test offline export (no fetch calls)
  - [ ] Target: >80% coverage

- [ ] Task 7: Verify CI passes (AC-CI)
  - [ ] All lint, test, type-check targets pass

## Dev Notes

### Architecture Patterns & Constraints

**File Locations:**
- Export dialog: `apps/web/app/components/dashboard/ExportDialog.tsx` (new)
- Download utility: `apps/web/app/lib/downloadFile.ts` (new)
- Existing report generators: `apps/web/app/lib/reports/`
  - `generateJsonReport.ts` - Already implemented
  - `generateHtmlReport.ts` - Already implemented
  - `generateMarkdownReport.ts` - Already implemented
  - `types.ts` - Report interfaces and builders
- Existing export hook: `apps/web/app/hooks/useReportExport.ts`
- Graph export utilities: `apps/web/app/components/visualization/DependencyGraph/utils/exportPng.ts`

**Existing Report System (FULLY BUILT in Epic 4):**
The report generation system is already complete from Epic 4, Story 4.7. This story primarily wires it into the new dashboard UI with a user-friendly export dialog.

Key existing functions:
```typescript
import { generateJsonReport } from '@/lib/reports/generateJsonReport';
import { generateHtmlReport } from '@/lib/reports/generateHtmlReport';
import { generateMarkdownReport } from '@/lib/reports/generateMarkdownReport';
```

The existing `ReportExportMenu.tsx` component (`apps/web/app/components/reports/ReportExportMenu.tsx`) may be reusable or serve as reference.

**Client-Side Download Pattern:**
```typescript
function downloadFile(content: string, filename: string, mimeType: string) {
  const blob = new Blob([content], { type: mimeType });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
}
```

### Previous Story Intelligence

**From Story 4.7 (Export Analysis Reports):**
- Report generation fully implemented
- Supports JSON, HTML, Markdown formats
- HTML report is self-contained with embedded styles
- Report sections: health score, circular deps, version conflicts, fix recommendations

**From Story 4.6 (Export Graph as PNG/SVG):**
- Graph export utilities exist in `utils/exportPng.ts` and `utils/exportSvg.ts`
- Can capture current graph state as data URL
- Used for embedding in HTML reports

**From Story 5.4 (Dashboard):**
- Export button in dashboard header area
- Analysis results available from Zustand analysis store

### Testing Requirements

- Mock report generators to verify they're called with correct params
- Test blob creation and download trigger
- Verify no network requests during export (offline compatibility)
- Test filename generation format
- Use `vi.spyOn(document, 'createElement')` for download link testing

### References

- [FR32: Download reports from web interface] `_bmad-output/planning-artifacts/epics.md`
- [FR19: HTML and JSON report export] `_bmad-output/planning-artifacts/epics.md`
- [Story 4.7: Export Analysis Reports] `_bmad-output/implementation-artifacts/4-7-*.md`
- [Existing Report System] `apps/web/app/lib/reports/`
- [Existing Export Menu] `apps/web/app/components/reports/ReportExportMenu.tsx`

## Dev Agent Record

### Agent Model Used

### Debug Log References

### Completion Notes List

### File List
