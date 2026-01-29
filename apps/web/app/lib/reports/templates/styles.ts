/**
 * Embedded CSS styles for self-contained HTML reports
 * AC3: HTML is self-contained with embedded styles
 * Supports dark mode via @media (prefers-color-scheme)
 * Print-friendly CSS for PDF generation
 */
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

    code {
      font-family: 'SF Mono', 'Fira Code', 'Fira Mono', monospace;
      font-size: 0.875em;
      padding: 0.125rem 0.25rem;
      background-color: var(--color-border);
      border-radius: 3px;
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
  `
}
