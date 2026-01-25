# Story 4.7: Export Analysis Reports in Multiple Formats

Status: ready-for-dev

## Story

As a **user**,
I want **to export complete analysis reports**,
So that **I can share findings with my team or archive results**.

## Acceptance Criteria

### AC1: Format Selection
**Given** complete analysis results
**When** I click export report
**Then** I can choose format:
- **JSON** - Machine-readable, complete data
- **HTML** - Standalone viewable report with embedded styles
- **Markdown** - For documentation/wikis

### AC2: JSON Export
**Given** I select JSON export
**When** the export completes
**Then**:
- JSON contains complete analysis data structure
- Data matches TypeScript `AnalysisResult` type exactly
- JSON is properly formatted (pretty-printed with indentation)
- JSON includes metadata (timestamp, version, project name)
- File can be re-imported for comparison (future feature)

### AC3: HTML Export
**Given** I select HTML export
**When** the export completes
**Then**:
- HTML is self-contained (no external CSS/JS dependencies)
- Styles are embedded inline or in `<style>` tags
- Report is viewable in any browser
- Report includes interactive elements (collapsible sections)
- Dark mode support via `@media (prefers-color-scheme)`
- Print-friendly CSS for PDF generation via browser

### AC4: Markdown Export
**Given** I select Markdown export
**When** the export completes
**Then**:
- Markdown is valid and renders correctly on GitHub/GitLab
- Tables are properly formatted using GFM (GitHub Flavored Markdown)
- Code blocks use proper syntax highlighting hints
- Report structure matches HTML report sections
- Images/diagrams represented as text descriptions or omitted

### AC5: Report Content - Health Score Summary
**Given** any export format
**When** the report is generated
**Then** it includes Health Score section:
- Overall health score (0-100) with color indicator
- Score breakdown by category (dependencies, architecture, security)
- Trend indicator if historical data available
- Score threshold indicators (Excellent/Good/Fair/Poor/Critical)

### AC6: Report Content - Circular Dependencies
**Given** any export format
**When** the report is generated
**Then** it includes Circular Dependencies section:
- Total count of circular dependencies
- List of all cycles with package paths
- Severity level for each cycle
- Quick reference to affected packages

### AC7: Report Content - Version Conflicts
**Given** any export format
**When** the report is generated
**Then** it includes Version Conflicts section:
- Total count of version conflicts
- List of conflicting packages with versions
- Risk level for each conflict
- Recommended resolution version

### AC8: Report Content - Fix Recommendations
**Given** any export format
**When** the report is generated
**Then** it includes Fix Recommendations section:
- Prioritized list of recommended fixes
- Estimated effort for each fix
- Impact assessment for each fix
- Quick wins highlighted (low effort, high impact)

### AC9: Report Metadata
**Given** any export format
**When** the report is generated
**Then** it includes metadata:
- Generation timestamp (ISO 8601)
- MonoGuard version
- Project/workspace name
- Analysis duration
- Package count analyzed

### AC10: Export UI Component
**Given** analysis results are displayed
**When** I look for report export
**Then**:
- Export report button is visible and accessible
- Clicking shows format selection (JSON/HTML/Markdown)
- Content customization options available
- Export progress indicator for large reports
- File downloads immediately on completion

### AC11: Performance Requirements
**Given** a large analysis (> 500 packages)
**When** exporting reports
**Then**:
- Export completes in < 3 seconds
- UI remains responsive during export
- Memory usage stays reasonable
- No browser freezing

### AC-CI: CI Pipeline Must Pass (MANDATORY)
**Given** the story implementation is complete
**When** verifying CI status
**Then** ALL of the following must pass:
- [ ] `pnpm nx affected --target=lint --base=main` passes
- [ ] `pnpm nx affected --target=test --base=main` passes
- [ ] `pnpm nx affected --target=type-check --base=main` passes
- [ ] `cd packages/analysis-engine && make test` passes (if Go changes)
- [ ] GitHub Actions CI workflow shows GREEN status
- **Story CANNOT be marked as "done" until CI is green**

## Tasks / Subtasks

- [ ] Task 1: Create Report Types and Interfaces (AC: 1, 9)
  - [ ] Define `ReportFormat` type ('json' | 'html' | 'markdown')
  - [ ] Define `ReportOptions` interface with content toggles
  - [ ] Define `ReportMetadata` interface
  - [ ] Define `ReportSection` enum for content selection

- [ ] Task 2: Create JSON Report Generator (AC: 2, 5-9)
  - [ ] Create `generateJsonReport.ts` utility
  - [ ] Structure data matching TypeScript types
  - [ ] Add metadata section
  - [ ] Pretty-print with configurable indentation
  - [ ] Create downloadable blob

- [ ] Task 3: Create HTML Report Generator (AC: 3, 5-9)
  - [ ] Create `generateHtmlReport.ts` utility
  - [ ] Design HTML template with embedded CSS
  - [ ] Create health score visualization (SVG gauge or colored bar)
  - [ ] Create collapsible sections for detailed data
  - [ ] Add dark mode support via media query
  - [ ] Add print-friendly styles
  - [ ] Make report self-contained (inline all styles)

- [ ] Task 4: Create Markdown Report Generator (AC: 4, 5-9)
  - [ ] Create `generateMarkdownReport.ts` utility
  - [ ] Create GFM-compatible tables
  - [ ] Add proper code block formatting
  - [ ] Structure with proper heading hierarchy
  - [ ] Include badges/shields for visual elements

- [ ] Task 5: Create Report Section Renderers (AC: 5-8)
  - [ ] Create `renderHealthScoreSection.ts` (shared logic)
  - [ ] Create `renderCircularDepsSection.ts`
  - [ ] Create `renderVersionConflictsSection.ts`
  - [ ] Create `renderFixRecommendationsSection.ts`
  - [ ] Each renderer outputs format-agnostic data

- [ ] Task 6: Create useReportExport Hook (AC: 10, 11)
  - [ ] Create `useReportExport.ts` hook
  - [ ] Manage export state (isExporting, progress)
  - [ ] Handle format selection
  - [ ] Generate filename with project name and timestamp
  - [ ] Trigger file download

- [ ] Task 7: Create ReportExportMenu Component (AC: 10)
  - [ ] Create `ReportExportMenu.tsx` component
  - [ ] Format selection (JSON/HTML/Markdown)
  - [ ] Content section toggles
  - [ ] Export button with loading state
  - [ ] Preview of selected sections

- [ ] Task 8: Integrate with Analysis Results UI (AC: 10)
  - [ ] Add export button to results dashboard
  - [ ] Wire up ReportExportMenu with analysis data
  - [ ] Handle export from different contexts (dashboard, graph view)

- [ ] Task 9: Write Unit Tests (AC: all)
  - [ ] Test JSON export generates valid JSON
  - [ ] Test HTML export is self-contained
  - [ ] Test Markdown export is GFM-compliant
  - [ ] Test all report sections render correctly
  - [ ] Test ReportExportMenu component
  - [ ] Test filename generation

- [ ] Task 10: Verify CI passes (AC-CI)
  - [ ] Run `pnpm nx affected --target=lint --base=main`
  - [ ] Run `pnpm nx affected --target=test --base=main`
  - [ ] Run `pnpm nx affected --target=type-check --base=main`
  - [ ] Verify GitHub Actions CI is GREEN

## Dev Notes

### Architecture Patterns & Constraints

**Dependency on Stories 4.1-4.6:** This story adds report export functionality that works alongside the graph export from Story 4.6.

**Related to Story 4.8:** Story 4.8 (Detailed Diagnostic Reports) will extend this functionality with per-cycle detailed reports.

**File Location:** `apps/web/app/lib/reports/`

**New Directory Structure:**
```
apps/web/app/lib/reports/
├── index.ts                       # Re-exports all report functions
├── types.ts                       # Report-specific types
├── generateJsonReport.ts          # JSON report generator
├── generateHtmlReport.ts          # HTML report generator
├── generateMarkdownReport.ts      # Markdown report generator
├── templates/
│   ├── html-template.ts           # HTML report template
│   └── styles.ts                  # Embedded CSS for HTML
├── sections/
│   ├── healthScore.ts             # Health score section renderer
│   ├── circularDependencies.ts    # Circular deps section
│   ├── versionConflicts.ts        # Version conflicts section
│   └── fixRecommendations.ts      # Fix recommendations section
└── __tests__/
    ├── generateJsonReport.test.ts
    ├── generateHtmlReport.test.ts
    ├── generateMarkdownReport.test.ts
    └── sections.test.ts

apps/web/app/components/reports/
├── ReportExportMenu.tsx           # Export menu component
├── ReportExportButton.tsx         # Export trigger button
└── __tests__/
    └── ReportExportMenu.test.tsx
```

### Key Implementation Details

**Report Types:**
```typescript
// apps/web/app/lib/reports/types.ts

export type ReportFormat = 'json' | 'html' | 'markdown';

export interface ReportOptions {
  format: ReportFormat;
  sections: ReportSections;
  includeMetadata: boolean;
  includeTimestamp: boolean;
  projectName: string;
}

export interface ReportSections {
  healthScore: boolean;
  circularDependencies: boolean;
  versionConflicts: boolean;
  fixRecommendations: boolean;
  packageList: boolean;
  graphSummary: boolean;
}

export interface ReportMetadata {
  generatedAt: string; // ISO 8601
  monoguardVersion: string;
  projectName: string;
  analysisDuration: number; // milliseconds
  packageCount: number;
  nodeCount: number;
  edgeCount: number;
}

export interface ReportData {
  metadata: ReportMetadata;
  healthScore: HealthScoreReport;
  circularDependencies: CircularDependencyReport;
  versionConflicts: VersionConflictReport;
  fixRecommendations: FixRecommendationReport;
}

export interface HealthScoreReport {
  overall: number;
  breakdown: {
    category: string;
    score: number;
    weight: number;
  }[];
  rating: 'excellent' | 'good' | 'fair' | 'poor' | 'critical';
  ratingThresholds: { [key: string]: number };
}

export interface CircularDependencyReport {
  totalCount: number;
  bySeverity: {
    critical: number;
    high: number;
    medium: number;
    low: number;
  };
  cycles: {
    id: string;
    packages: string[];
    severity: string;
    type: 'direct' | 'indirect';
  }[];
}

export interface VersionConflictReport {
  totalCount: number;
  byRiskLevel: {
    critical: number;
    high: number;
    medium: number;
    low: number;
  };
  conflicts: {
    packageName: string;
    versions: string[];
    riskLevel: string;
    recommendedVersion: string;
  }[];
}

export interface FixRecommendationReport {
  totalCount: number;
  quickWins: number; // low effort, high impact
  recommendations: {
    id: string;
    title: string;
    description: string;
    effort: 'low' | 'medium' | 'high';
    impact: 'low' | 'medium' | 'high';
    priority: number;
    affectedPackages: string[];
  }[];
}

export interface ReportResult {
  blob: Blob;
  filename: string;
  format: ReportFormat;
  sizeBytes: number;
}
```

**JSON Report Generator:**
```typescript
// apps/web/app/lib/reports/generateJsonReport.ts

import type { ReportData, ReportOptions, ReportResult } from './types';

const MONOGUARD_VERSION = '0.1.0'; // Should come from package.json

export function generateJsonReport(
  data: ReportData,
  options: ReportOptions
): ReportResult {
  const report: Record<string, unknown> = {};

  // Always include metadata if enabled
  if (options.includeMetadata) {
    report.metadata = {
      ...data.metadata,
      monoguardVersion: MONOGUARD_VERSION,
      generatedAt: new Date().toISOString(),
    };
  }

  // Include selected sections
  if (options.sections.healthScore) {
    report.healthScore = data.healthScore;
  }

  if (options.sections.circularDependencies) {
    report.circularDependencies = data.circularDependencies;
  }

  if (options.sections.versionConflicts) {
    report.versionConflicts = data.versionConflicts;
  }

  if (options.sections.fixRecommendations) {
    report.fixRecommendations = data.fixRecommendations;
  }

  // Pretty-print JSON with 2-space indentation
  const jsonString = JSON.stringify(report, null, 2);
  const blob = new Blob([jsonString], { type: 'application/json' });

  // Generate filename
  const timestamp = new Date().toISOString().split('T')[0];
  const filename = `${options.projectName}-analysis-report-${timestamp}.json`;

  return {
    blob,
    filename,
    format: 'json',
    sizeBytes: blob.size,
  };
}
```

**HTML Report Generator:**
```typescript
// apps/web/app/lib/reports/generateHtmlReport.ts

import type { ReportData, ReportOptions, ReportResult } from './types';
import { getEmbeddedStyles } from './templates/styles';
import { renderHealthScoreHtml } from './sections/healthScore';
import { renderCircularDepsHtml } from './sections/circularDependencies';
import { renderVersionConflictsHtml } from './sections/versionConflicts';
import { renderFixRecommendationsHtml } from './sections/fixRecommendations';

const MONOGUARD_VERSION = '0.1.0';

export function generateHtmlReport(
  data: ReportData,
  options: ReportOptions
): ReportResult {
  const sections: string[] = [];

  // Build sections based on options
  if (options.sections.healthScore) {
    sections.push(renderHealthScoreHtml(data.healthScore));
  }

  if (options.sections.circularDependencies) {
    sections.push(renderCircularDepsHtml(data.circularDependencies));
  }

  if (options.sections.versionConflicts) {
    sections.push(renderVersionConflictsHtml(data.versionConflicts));
  }

  if (options.sections.fixRecommendations) {
    sections.push(renderFixRecommendationsHtml(data.fixRecommendations));
  }

  const timestamp = new Date().toISOString();
  const formattedDate = new Date().toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });

  const html = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>${options.projectName} - Dependency Analysis Report</title>
  <style>
${getEmbeddedStyles()}
  </style>
</head>
<body>
  <div class="container">
    <header class="report-header">
      <div class="logo">
        <svg width="32" height="32" viewBox="0 0 32 32" fill="currentColor">
          <circle cx="16" cy="16" r="14" stroke="currentColor" stroke-width="2" fill="none"/>
          <path d="M16 8v8l6 4" stroke="currentColor" stroke-width="2" fill="none"/>
        </svg>
        <span>MonoGuard</span>
      </div>
      <h1>Dependency Analysis Report</h1>
      <div class="report-meta">
        <span class="project-name">${escapeHtml(options.projectName)}</span>
        <span class="timestamp">Generated: ${formattedDate}</span>
      </div>
    </header>

    <main class="report-content">
      ${sections.join('\n')}
    </main>

    <footer class="report-footer">
      <p>Generated by MonoGuard v${MONOGUARD_VERSION}</p>
      <p class="timestamp-iso">${timestamp}</p>
    </footer>
  </div>

  <script>
    // Collapsible sections
    document.querySelectorAll('.section-header').forEach(header => {
      header.addEventListener('click', () => {
        const section = header.parentElement;
        section.classList.toggle('collapsed');
      });
    });
  </script>
</body>
</html>`;

  const blob = new Blob([html], { type: 'text/html' });
  const dateStr = new Date().toISOString().split('T')[0];
  const filename = `${options.projectName}-analysis-report-${dateStr}.html`;

  return {
    blob,
    filename,
    format: 'html',
    sizeBytes: blob.size,
  };
}

function escapeHtml(str: string): string {
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;');
}
```

**Embedded CSS Styles:**
```typescript
// apps/web/app/lib/reports/templates/styles.ts

export function getEmbeddedStyles(): string {
  return `
    :root {
      --color-bg: #ffffff;
      --color-text: #1f2937;
      --color-text-secondary: #6b7280;
      --color-border: #e5e7eb;
      --color-success: #10b981;
      --color-warning: #f59e0b;
      --color-error: #ef4444;
      --color-info: #3b82f6;
      --color-excellent: #10b981;
      --color-good: #22c55e;
      --color-fair: #f59e0b;
      --color-poor: #f97316;
      --color-critical: #ef4444;
    }

    @media (prefers-color-scheme: dark) {
      :root {
        --color-bg: #111827;
        --color-text: #f9fafb;
        --color-text-secondary: #9ca3af;
        --color-border: #374151;
      }
    }

    * {
      box-sizing: border-box;
      margin: 0;
      padding: 0;
    }

    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      background-color: var(--color-bg);
      color: var(--color-text);
      line-height: 1.6;
    }

    .container {
      max-width: 1200px;
      margin: 0 auto;
      padding: 2rem;
    }

    .report-header {
      text-align: center;
      margin-bottom: 3rem;
      padding-bottom: 2rem;
      border-bottom: 1px solid var(--color-border);
    }

    .report-header .logo {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 0.5rem;
      font-size: 1.25rem;
      font-weight: 600;
      color: var(--color-info);
      margin-bottom: 1rem;
    }

    .report-header h1 {
      font-size: 2rem;
      margin-bottom: 1rem;
    }

    .report-meta {
      display: flex;
      justify-content: center;
      gap: 2rem;
      color: var(--color-text-secondary);
    }

    .section {
      margin-bottom: 2rem;
      border: 1px solid var(--color-border);
      border-radius: 8px;
      overflow: hidden;
    }

    .section-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 1rem 1.5rem;
      background-color: var(--color-border);
      cursor: pointer;
      user-select: none;
    }

    .section-header:hover {
      opacity: 0.9;
    }

    .section-header h2 {
      font-size: 1.25rem;
      font-weight: 600;
    }

    .section-header .badge {
      padding: 0.25rem 0.75rem;
      border-radius: 9999px;
      font-size: 0.875rem;
      font-weight: 500;
    }

    .section-content {
      padding: 1.5rem;
    }

    .section.collapsed .section-content {
      display: none;
    }

    .health-score {
      text-align: center;
      padding: 2rem;
    }

    .health-score .score {
      font-size: 4rem;
      font-weight: 700;
    }

    .health-score .rating {
      font-size: 1.5rem;
      text-transform: capitalize;
    }

    .health-score.excellent .score,
    .health-score.excellent .rating { color: var(--color-excellent); }
    .health-score.good .score,
    .health-score.good .rating { color: var(--color-good); }
    .health-score.fair .score,
    .health-score.fair .rating { color: var(--color-fair); }
    .health-score.poor .score,
    .health-score.poor .rating { color: var(--color-poor); }
    .health-score.critical .score,
    .health-score.critical .rating { color: var(--color-critical); }

    .breakdown-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
      gap: 1rem;
      margin-top: 2rem;
    }

    .breakdown-item {
      padding: 1rem;
      border: 1px solid var(--color-border);
      border-radius: 8px;
    }

    .breakdown-item .label {
      font-size: 0.875rem;
      color: var(--color-text-secondary);
    }

    .breakdown-item .value {
      font-size: 1.5rem;
      font-weight: 600;
    }

    table {
      width: 100%;
      border-collapse: collapse;
    }

    th, td {
      padding: 0.75rem;
      text-align: left;
      border-bottom: 1px solid var(--color-border);
    }

    th {
      font-weight: 600;
      background-color: var(--color-border);
    }

    .severity-critical { color: var(--color-critical); }
    .severity-high { color: var(--color-error); }
    .severity-medium { color: var(--color-warning); }
    .severity-low { color: var(--color-info); }

    .fix-card {
      padding: 1rem;
      border: 1px solid var(--color-border);
      border-radius: 8px;
      margin-bottom: 1rem;
    }

    .fix-card.quick-win {
      border-color: var(--color-success);
      background-color: rgba(16, 185, 129, 0.05);
    }

    .fix-card .title {
      font-weight: 600;
      margin-bottom: 0.5rem;
    }

    .fix-card .meta {
      display: flex;
      gap: 1rem;
      font-size: 0.875rem;
      color: var(--color-text-secondary);
    }

    .report-footer {
      margin-top: 3rem;
      padding-top: 2rem;
      border-top: 1px solid var(--color-border);
      text-align: center;
      color: var(--color-text-secondary);
      font-size: 0.875rem;
    }

    @media print {
      .section-header { cursor: default; }
      .section.collapsed .section-content { display: block; }
      body { font-size: 12pt; }
      .container { max-width: none; padding: 0; }
    }
  `;
}
```

**Markdown Report Generator:**
```typescript
// apps/web/app/lib/reports/generateMarkdownReport.ts

import type { ReportData, ReportOptions, ReportResult } from './types';

const MONOGUARD_VERSION = '0.1.0';

export function generateMarkdownReport(
  data: ReportData,
  options: ReportOptions
): ReportResult {
  const sections: string[] = [];

  // Header
  sections.push(`# ${options.projectName} - Dependency Analysis Report\n`);
  sections.push(`> Generated by MonoGuard v${MONOGUARD_VERSION}\n`);
  sections.push(`> ${new Date().toISOString()}\n`);
  sections.push('---\n');

  // Health Score Section
  if (options.sections.healthScore) {
    sections.push(renderHealthScoreMd(data.healthScore));
  }

  // Circular Dependencies Section
  if (options.sections.circularDependencies) {
    sections.push(renderCircularDepsMd(data.circularDependencies));
  }

  // Version Conflicts Section
  if (options.sections.versionConflicts) {
    sections.push(renderVersionConflictsMd(data.versionConflicts));
  }

  // Fix Recommendations Section
  if (options.sections.fixRecommendations) {
    sections.push(renderFixRecommendationsMd(data.fixRecommendations));
  }

  // Metadata Section
  if (options.includeMetadata) {
    sections.push(renderMetadataMd(data.metadata));
  }

  const markdown = sections.join('\n');
  const blob = new Blob([markdown], { type: 'text/markdown' });
  const dateStr = new Date().toISOString().split('T')[0];
  const filename = `${options.projectName}-analysis-report-${dateStr}.md`;

  return {
    blob,
    filename,
    format: 'markdown',
    sizeBytes: blob.size,
  };
}

function renderHealthScoreMd(data: HealthScoreReport): string {
  const ratingEmoji = {
    excellent: ':white_check_mark:',
    good: ':heavy_check_mark:',
    fair: ':warning:',
    poor: ':x:',
    critical: ':rotating_light:',
  };

  let md = `## Health Score\n\n`;
  md += `**Overall Score: ${data.overall}/100** ${ratingEmoji[data.rating]} ${data.rating.toUpperCase()}\n\n`;

  md += `### Score Breakdown\n\n`;
  md += `| Category | Score | Weight |\n`;
  md += `|----------|-------|--------|\n`;

  for (const item of data.breakdown) {
    md += `| ${item.category} | ${item.score} | ${item.weight}% |\n`;
  }

  md += '\n';
  return md;
}

function renderCircularDepsMd(data: CircularDependencyReport): string {
  let md = `## Circular Dependencies\n\n`;
  md += `**Total: ${data.totalCount}**\n\n`;

  if (data.totalCount === 0) {
    md += `> :tada: No circular dependencies detected!\n\n`;
    return md;
  }

  md += `### Summary by Severity\n\n`;
  md += `| Severity | Count |\n`;
  md += `|----------|-------|\n`;
  md += `| Critical | ${data.bySeverity.critical} |\n`;
  md += `| High | ${data.bySeverity.high} |\n`;
  md += `| Medium | ${data.bySeverity.medium} |\n`;
  md += `| Low | ${data.bySeverity.low} |\n\n`;

  md += `### Detected Cycles\n\n`;

  for (const cycle of data.cycles) {
    md += `#### Cycle: ${cycle.id}\n`;
    md += `- **Severity:** ${cycle.severity}\n`;
    md += `- **Type:** ${cycle.type}\n`;
    md += `- **Path:** \`${cycle.packages.join(' → ')} → ${cycle.packages[0]}\`\n\n`;
  }

  return md;
}

function renderVersionConflictsMd(data: VersionConflictReport): string {
  let md = `## Version Conflicts\n\n`;
  md += `**Total: ${data.totalCount}**\n\n`;

  if (data.totalCount === 0) {
    md += `> :white_check_mark: No version conflicts detected!\n\n`;
    return md;
  }

  md += `| Package | Conflicting Versions | Risk | Recommended |\n`;
  md += `|---------|---------------------|------|-------------|\n`;

  for (const conflict of data.conflicts) {
    md += `| \`${conflict.packageName}\` | ${conflict.versions.join(', ')} | ${conflict.riskLevel} | ${conflict.recommendedVersion} |\n`;
  }

  md += '\n';
  return md;
}

function renderFixRecommendationsMd(data: FixRecommendationReport): string {
  let md = `## Fix Recommendations\n\n`;
  md += `**Total Recommendations: ${data.totalCount}**\n`;
  md += `**Quick Wins: ${data.quickWins}** :zap:\n\n`;

  if (data.totalCount === 0) {
    md += `> No fix recommendations at this time.\n\n`;
    return md;
  }

  md += `### Priority Fixes\n\n`;

  for (const rec of data.recommendations) {
    const quickWinBadge = rec.effort === 'low' && rec.impact === 'high' ? ' :zap: Quick Win' : '';
    md += `#### ${rec.priority}. ${rec.title}${quickWinBadge}\n\n`;
    md += `${rec.description}\n\n`;
    md += `- **Effort:** ${rec.effort}\n`;
    md += `- **Impact:** ${rec.impact}\n`;
    md += `- **Affected Packages:** ${rec.affectedPackages.map(p => `\`${p}\``).join(', ')}\n\n`;
  }

  return md;
}

function renderMetadataMd(data: ReportMetadata): string {
  let md = `---\n\n## Report Metadata\n\n`;
  md += `| Property | Value |\n`;
  md += `|----------|-------|\n`;
  md += `| Project | ${data.projectName} |\n`;
  md += `| Packages Analyzed | ${data.packageCount} |\n`;
  md += `| Graph Nodes | ${data.nodeCount} |\n`;
  md += `| Graph Edges | ${data.edgeCount} |\n`;
  md += `| Analysis Duration | ${data.analysisDuration}ms |\n`;
  md += `| Generated At | ${data.generatedAt} |\n`;
  md += `| MonoGuard Version | ${data.monoguardVersion} |\n\n`;

  return md;
}
```

**useReportExport Hook:**
```typescript
// apps/web/app/hooks/useReportExport.ts

import { useState, useCallback } from 'react';
import type { ReportOptions, ReportResult, ReportData, ReportFormat } from '@/lib/reports/types';
import { generateJsonReport } from '@/lib/reports/generateJsonReport';
import { generateHtmlReport } from '@/lib/reports/generateHtmlReport';
import { generateMarkdownReport } from '@/lib/reports/generateMarkdownReport';

interface ExportProgress {
  isExporting: boolean;
  progress: number;
  stage: 'preparing' | 'generating' | 'complete';
}

interface UseReportExportResult {
  exportProgress: ExportProgress;
  startExport: (data: ReportData, options: ReportOptions) => Promise<void>;
  cancelExport: () => void;
}

export function useReportExport(): UseReportExportResult {
  const [exportProgress, setExportProgress] = useState<ExportProgress>({
    isExporting: false,
    progress: 0,
    stage: 'preparing',
  });

  const startExport = useCallback(async (data: ReportData, options: ReportOptions) => {
    setExportProgress({
      isExporting: true,
      progress: 10,
      stage: 'preparing',
    });

    try {
      setExportProgress(prev => ({ ...prev, progress: 30, stage: 'generating' }));

      let result: ReportResult;

      switch (options.format) {
        case 'json':
          result = generateJsonReport(data, options);
          break;
        case 'html':
          result = generateHtmlReport(data, options);
          break;
        case 'markdown':
          result = generateMarkdownReport(data, options);
          break;
        default:
          throw new Error(`Unknown format: ${options.format}`);
      }

      setExportProgress(prev => ({ ...prev, progress: 90 }));

      // Trigger download
      downloadBlob(result.blob, result.filename);

      setExportProgress({
        isExporting: false,
        progress: 100,
        stage: 'complete',
      });
    } catch (error) {
      console.error('Report export failed:', error);
      setExportProgress({
        isExporting: false,
        progress: 0,
        stage: 'preparing',
      });
      throw error;
    }
  }, []);

  const cancelExport = useCallback(() => {
    setExportProgress({
      isExporting: false,
      progress: 0,
      stage: 'preparing',
    });
  }, []);

  return {
    exportProgress,
    startExport,
    cancelExport,
  };
}

function downloadBlob(blob: Blob, filename: string) {
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(url);
}
```

**ReportExportMenu Component:**
```typescript
// apps/web/app/components/reports/ReportExportMenu.tsx

import React, { useState } from 'react';
import type { ReportOptions, ReportFormat, ReportSections } from '@/lib/reports/types';

interface ReportExportMenuProps {
  isOpen: boolean;
  onClose: () => void;
  onExport: (options: ReportOptions) => Promise<void>;
  isExporting: boolean;
  projectName: string;
}

const DEFAULT_SECTIONS: ReportSections = {
  healthScore: true,
  circularDependencies: true,
  versionConflicts: true,
  fixRecommendations: true,
  packageList: false,
  graphSummary: false,
};

export function ReportExportMenu({
  isOpen,
  onClose,
  onExport,
  isExporting,
  projectName,
}: ReportExportMenuProps) {
  const [format, setFormat] = useState<ReportFormat>('html');
  const [sections, setSections] = useState<ReportSections>(DEFAULT_SECTIONS);

  const handleExport = async () => {
    await onExport({
      format,
      sections,
      includeMetadata: true,
      includeTimestamp: true,
      projectName,
    });
    onClose();
  };

  const toggleSection = (key: keyof ReportSections) => {
    setSections(prev => ({ ...prev, [key]: !prev[key] }));
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl w-full max-w-md p-6">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white">
            Export Analysis Report
          </h2>
          <button
            onClick={onClose}
            className="text-gray-500 hover:text-gray-700 dark:hover:text-gray-300"
            aria-label="Close"
          >
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        {/* Format Selection */}
        <div className="mb-6">
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
            Export Format
          </label>
          <div className="grid grid-cols-3 gap-2">
            {(['json', 'html', 'markdown'] as ReportFormat[]).map((f) => (
              <button
                key={f}
                onClick={() => setFormat(f)}
                className={`px-4 py-2 rounded-md text-sm font-medium transition-colors
                  ${format === f
                    ? 'bg-blue-600 text-white'
                    : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
                  }`}
              >
                {f.toUpperCase()}
              </button>
            ))}
          </div>
          <p className="mt-2 text-xs text-gray-500 dark:text-gray-400">
            {format === 'json' && 'Machine-readable format for programmatic access'}
            {format === 'html' && 'Standalone report viewable in any browser'}
            {format === 'markdown' && 'Perfect for documentation and wikis'}
          </p>
        </div>

        {/* Section Selection */}
        <div className="mb-6">
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
            Include Sections
          </label>
          <div className="space-y-2">
            {Object.entries({
              healthScore: 'Health Score Summary',
              circularDependencies: 'Circular Dependencies',
              versionConflicts: 'Version Conflicts',
              fixRecommendations: 'Fix Recommendations',
            }).map(([key, label]) => (
              <label key={key} className="flex items-center gap-3 cursor-pointer">
                <input
                  type="checkbox"
                  checked={sections[key as keyof ReportSections]}
                  onChange={() => toggleSection(key as keyof ReportSections)}
                  className="rounded border-gray-300 dark:border-gray-600
                           text-blue-600 focus:ring-blue-500"
                />
                <span className="text-sm text-gray-700 dark:text-gray-300">{label}</span>
              </label>
            ))}
          </div>
        </div>

        {/* Export Button */}
        <div className="flex gap-3">
          <button
            onClick={onClose}
            className="flex-1 px-4 py-2 border border-gray-300 dark:border-gray-600
                     text-gray-700 dark:text-gray-300 rounded-md hover:bg-gray-50
                     dark:hover:bg-gray-700 transition-colors"
          >
            Cancel
          </button>
          <button
            onClick={handleExport}
            disabled={isExporting}
            className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-md font-medium
                     hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500
                     disabled:opacity-50 disabled:cursor-not-allowed transition-colors
                     flex items-center justify-center gap-2"
          >
            {isExporting ? (
              <>
                <svg className="animate-spin h-4 w-4" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                </svg>
                Exporting...
              </>
            ) : (
              <>
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                        d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                </svg>
                Export Report
              </>
            )}
          </button>
        </div>
      </div>
    </div>
  );
}
```

### Testing Requirements

**Test File:** `apps/web/src/__tests__/generateJsonReport.test.ts`

```typescript
import { generateJsonReport } from '@/lib/reports/generateJsonReport';
import type { ReportData, ReportOptions } from '@/lib/reports/types';

describe('generateJsonReport', () => {
  const mockData: ReportData = {
    metadata: {
      generatedAt: '2026-01-25T10:00:00Z',
      monoguardVersion: '0.1.0',
      projectName: 'test-project',
      analysisDuration: 1234,
      packageCount: 50,
      nodeCount: 50,
      edgeCount: 120,
    },
    healthScore: {
      overall: 75,
      breakdown: [
        { category: 'Dependencies', score: 80, weight: 40 },
        { category: 'Architecture', score: 70, weight: 30 },
      ],
      rating: 'good',
      ratingThresholds: { excellent: 85, good: 70, fair: 50, poor: 30 },
    },
    circularDependencies: {
      totalCount: 2,
      bySeverity: { critical: 0, high: 1, medium: 1, low: 0 },
      cycles: [],
    },
    versionConflicts: {
      totalCount: 3,
      byRiskLevel: { critical: 0, high: 1, medium: 2, low: 0 },
      conflicts: [],
    },
    fixRecommendations: {
      totalCount: 5,
      quickWins: 2,
      recommendations: [],
    },
  };

  const defaultOptions: ReportOptions = {
    format: 'json',
    sections: {
      healthScore: true,
      circularDependencies: true,
      versionConflicts: true,
      fixRecommendations: true,
      packageList: false,
      graphSummary: false,
    },
    includeMetadata: true,
    includeTimestamp: true,
    projectName: 'test-project',
  };

  it('should generate valid JSON blob', () => {
    const result = generateJsonReport(mockData, defaultOptions);

    expect(result.blob).toBeInstanceOf(Blob);
    expect(result.blob.type).toBe('application/json');
  });

  it('should include selected sections only', async () => {
    const options = {
      ...defaultOptions,
      sections: {
        ...defaultOptions.sections,
        versionConflicts: false,
      },
    };

    const result = generateJsonReport(mockData, options);
    const text = await result.blob.text();
    const json = JSON.parse(text);

    expect(json.healthScore).toBeDefined();
    expect(json.versionConflicts).toBeUndefined();
  });

  it('should generate correct filename', () => {
    const result = generateJsonReport(mockData, defaultOptions);

    expect(result.filename).toMatch(/^test-project-analysis-report-\d{4}-\d{2}-\d{2}\.json$/);
  });

  it('should include metadata when enabled', async () => {
    const result = generateJsonReport(mockData, defaultOptions);
    const text = await result.blob.text();
    const json = JSON.parse(text);

    expect(json.metadata).toBeDefined();
    expect(json.metadata.projectName).toBe('test-project');
  });

  it('should format JSON with indentation', async () => {
    const result = generateJsonReport(mockData, defaultOptions);
    const text = await result.blob.text();

    // Pretty-printed JSON has newlines
    expect(text).toContain('\n');
    expect(text).toContain('  ');
  });
});
```

**Test File:** `apps/web/src/__tests__/generateHtmlReport.test.ts`

```typescript
import { generateHtmlReport } from '@/lib/reports/generateHtmlReport';
import type { ReportData, ReportOptions } from '@/lib/reports/types';

describe('generateHtmlReport', () => {
  const mockData: ReportData = {
    // ... same as JSON test
  };

  const defaultOptions: ReportOptions = {
    format: 'html',
    sections: {
      healthScore: true,
      circularDependencies: true,
      versionConflicts: true,
      fixRecommendations: true,
      packageList: false,
      graphSummary: false,
    },
    includeMetadata: true,
    includeTimestamp: true,
    projectName: 'test-project',
  };

  it('should generate HTML blob', () => {
    const result = generateHtmlReport(mockData, defaultOptions);

    expect(result.blob).toBeInstanceOf(Blob);
    expect(result.blob.type).toBe('text/html');
  });

  it('should be self-contained (no external links)', async () => {
    const result = generateHtmlReport(mockData, defaultOptions);
    const text = await result.blob.text();

    // Should not contain external stylesheet links
    expect(text).not.toMatch(/<link[^>]+href="http/);
    // Should not contain external script sources
    expect(text).not.toMatch(/<script[^>]+src="http/);
    // Should contain embedded styles
    expect(text).toContain('<style>');
  });

  it('should include dark mode media query', async () => {
    const result = generateHtmlReport(mockData, defaultOptions);
    const text = await result.blob.text();

    expect(text).toContain('prefers-color-scheme: dark');
  });

  it('should include print-friendly styles', async () => {
    const result = generateHtmlReport(mockData, defaultOptions);
    const text = await result.blob.text();

    expect(text).toContain('@media print');
  });

  it('should escape HTML in project name', async () => {
    const options = {
      ...defaultOptions,
      projectName: '<script>alert("xss")</script>',
    };

    const result = generateHtmlReport(mockData, options);
    const text = await result.blob.text();

    expect(text).not.toContain('<script>alert');
    expect(text).toContain('&lt;script&gt;');
  });
});
```

**Test File:** `apps/web/src/__tests__/ReportExportMenu.test.tsx`

```typescript
import { render, screen, fireEvent } from '@testing-library/react';
import { ReportExportMenu } from '@/components/reports/ReportExportMenu';

describe('ReportExportMenu', () => {
  const mockOnClose = vi.fn();
  const mockOnExport = vi.fn().mockResolvedValue(undefined);

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should not render when closed', () => {
    render(
      <ReportExportMenu
        isOpen={false}
        onClose={mockOnClose}
        onExport={mockOnExport}
        isExporting={false}
        projectName="test"
      />
    );

    expect(screen.queryByText('Export Analysis Report')).not.toBeInTheDocument();
  });

  it('should render format options when open', () => {
    render(
      <ReportExportMenu
        isOpen={true}
        onClose={mockOnClose}
        onExport={mockOnExport}
        isExporting={false}
        projectName="test"
      />
    );

    expect(screen.getByText('JSON')).toBeInTheDocument();
    expect(screen.getByText('HTML')).toBeInTheDocument();
    expect(screen.getByText('MARKDOWN')).toBeInTheDocument();
  });

  it('should call onExport with selected options', async () => {
    render(
      <ReportExportMenu
        isOpen={true}
        onClose={mockOnClose}
        onExport={mockOnExport}
        isExporting={false}
        projectName="test-project"
      />
    );

    // Click export
    fireEvent.click(screen.getByText('Export Report'));

    expect(mockOnExport).toHaveBeenCalledWith(expect.objectContaining({
      format: 'html', // default
      projectName: 'test-project',
      sections: expect.objectContaining({
        healthScore: true,
      }),
    }));
  });

  it('should toggle section checkboxes', () => {
    render(
      <ReportExportMenu
        isOpen={true}
        onClose={mockOnClose}
        onExport={mockOnExport}
        isExporting={false}
        projectName="test"
      />
    );

    const healthScoreCheckbox = screen.getByRole('checkbox', { name: /health score/i });
    expect(healthScoreCheckbox).toBeChecked();

    fireEvent.click(healthScoreCheckbox);
    expect(healthScoreCheckbox).not.toBeChecked();
  });

  it('should show loading state when exporting', () => {
    render(
      <ReportExportMenu
        isOpen={true}
        onClose={mockOnClose}
        onExport={mockOnExport}
        isExporting={true}
        projectName="test"
      />
    );

    expect(screen.getByText('Exporting...')).toBeInTheDocument();
  });
});
```

### Critical Don't-Miss Rules (from project-context.md)

1. **JSON uses camelCase** - All JSON output must use camelCase matching TypeScript types
2. **ISO 8601 dates** - All timestamps must be ISO 8601 format
3. **Self-contained HTML** - No external dependencies in HTML reports
4. **XSS Prevention** - Escape all user-provided content in HTML
5. **Performance** - Export should complete < 3 seconds

### Previous Story Intelligence (Story 4.6)

**Key Patterns to Follow:**
- Separate utility functions for each format
- Hook for managing export state
- Menu component for options
- Blob creation and download pattern

**Integration Points:**
- Report export button can be placed alongside graph export
- Both exports can share similar UI patterns
- Consider unified export menu with tabs (Graph | Report)

### UX Design Requirements

- **Format selection:** Clear visual distinction between formats
- **Section toggles:** Allow users to customize report content
- **Progress feedback:** Show export progress
- **Immediate download:** File downloads when ready

### Performance Considerations

1. **String building:** Use array join instead of string concatenation
2. **Large reports:** Consider chunked generation for very large data
3. **Memory:** Clear intermediate strings after blob creation

### References

- [Story 4.6: Graph Export] `4-6-export-graph-as-png-svg-images.md`
- [Epic 4 Story 4.7 Requirements] `_bmad-output/planning-artifacts/epics.md` - Lines 1098-1119
- [FR19: HTML/JSON Report Export] `_bmad-output/planning-artifacts/epics.md` - Line 48
- [NFR15: Export Formats] `_bmad-output/planning-artifacts/epics.md` - Line 112
- [Types Package] `packages/types/src/domain.ts` - AnalysisResult, HealthScore types
- [Project Context] `_bmad-output/project-context.md` - JSON formatting rules

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
