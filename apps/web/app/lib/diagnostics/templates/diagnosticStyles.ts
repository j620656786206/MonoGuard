/**
 * Embedded CSS for self-contained diagnostic HTML reports
 * AC7: PDF-Ready HTML Export with print-friendly CSS
 */
export function getDiagnosticStyles(): string {
  return `
    * { margin: 0; padding: 0; box-sizing: border-box; }

    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      line-height: 1.6;
      color: #1f2937;
      background: #ffffff;
      max-width: 900px;
      margin: 0 auto;
      padding: 2rem;
    }

    @media (prefers-color-scheme: dark) {
      body { background: #111827; color: #f9fafb; }
      .section { background: #1f2937; border-color: #374151; }
      .severity-badge { border-color: #374151; }
      code { background: #374151; color: #e5e7eb; }
      .toc a { color: #60a5fa; }
      table { border-color: #374151; }
      th { background: #374151; }
      td { border-color: #374151; }
      .code-block { background: #1e293b; border-color: #374151; }
    }

    .report-header {
      text-align: center;
      margin-bottom: 2rem;
      padding-bottom: 1rem;
      border-bottom: 2px solid #e5e7eb;
    }

    .report-header h1 { font-size: 1.75rem; margin-bottom: 0.5rem; }
    .report-header .subtitle { color: #6b7280; font-size: 0.875rem; }

    .toc {
      margin: 1.5rem 0;
      padding: 1rem 1.5rem;
      background: #f9fafb;
      border-radius: 8px;
      border: 1px solid #e5e7eb;
    }
    .toc h2 { font-size: 1rem; margin-bottom: 0.5rem; }
    .toc ul { list-style: none; padding-left: 0; }
    .toc li { margin: 0.25rem 0; }
    .toc a { color: #2563eb; text-decoration: none; }
    .toc a:hover { text-decoration: underline; }

    .section {
      margin: 1.5rem 0;
      padding: 1.25rem;
      background: #f9fafb;
      border: 1px solid #e5e7eb;
      border-radius: 8px;
    }

    .section h2 {
      font-size: 1.25rem;
      margin-bottom: 0.75rem;
      padding-bottom: 0.5rem;
      border-bottom: 1px solid #e5e7eb;
    }

    .severity-badge {
      display: inline-block;
      padding: 0.125rem 0.5rem;
      border-radius: 9999px;
      font-size: 0.75rem;
      font-weight: 600;
      text-transform: uppercase;
    }
    .severity-critical { background: #fef2f2; color: #dc2626; border: 1px solid #fecaca; }
    .severity-high { background: #fff7ed; color: #ea580c; border: 1px solid #fed7aa; }
    .severity-medium { background: #fffbeb; color: #d97706; border: 1px solid #fde68a; }
    .severity-low { background: #f0fdf4; color: #16a34a; border: 1px solid #bbf7d0; }

    .effort-badge {
      display: inline-block;
      padding: 0.125rem 0.5rem;
      border-radius: 4px;
      font-size: 0.75rem;
      background: #eff6ff;
      color: #2563eb;
      border: 1px solid #bfdbfe;
    }

    .metric-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
      gap: 0.75rem;
      margin: 0.75rem 0;
    }

    .metric-card {
      padding: 0.75rem;
      text-align: center;
      background: white;
      border-radius: 6px;
      border: 1px solid #e5e7eb;
    }
    .metric-value { font-size: 1.5rem; font-weight: 700; color: #1f2937; }
    .metric-label { font-size: 0.75rem; color: #6b7280; }

    table {
      width: 100%;
      border-collapse: collapse;
      margin: 0.75rem 0;
      font-size: 0.875rem;
    }
    th, td { padding: 0.5rem 0.75rem; text-align: left; border: 1px solid #e5e7eb; }
    th { background: #f3f4f6; font-weight: 600; }

    .code-block {
      background: #f8fafc;
      border: 1px solid #e2e8f0;
      border-radius: 6px;
      padding: 0.75rem 1rem;
      font-family: 'Fira Code', 'SF Mono', Consolas, monospace;
      font-size: 0.8125rem;
      overflow-x: auto;
      white-space: pre;
      margin: 0.5rem 0;
    }

    code {
      background: #f3f4f6;
      padding: 0.125rem 0.25rem;
      border-radius: 3px;
      font-size: 0.8125rem;
    }

    .strategy-card {
      margin: 0.75rem 0;
      padding: 1rem;
      border: 1px solid #e5e7eb;
      border-radius: 6px;
      background: white;
    }
    .strategy-card h3 { margin-bottom: 0.5rem; }

    .pros-cons { display: grid; grid-template-columns: 1fr 1fr; gap: 0.75rem; margin: 0.5rem 0; }
    .pros li::marker { content: '✅ '; }
    .cons li::marker { content: '⚠️ '; }
    .pros, .cons { padding-left: 1.25rem; }

    .step-list { counter-reset: step-counter; list-style: none; padding-left: 0; }
    .step-list li {
      counter-increment: step-counter;
      margin: 0.5rem 0;
      padding-left: 2rem;
      position: relative;
    }
    .step-list li::before {
      content: counter(step-counter);
      position: absolute;
      left: 0;
      width: 1.5rem;
      height: 1.5rem;
      background: #3b82f6;
      color: white;
      border-radius: 50%;
      text-align: center;
      font-size: 0.75rem;
      line-height: 1.5rem;
    }

    .report-footer {
      margin-top: 2rem;
      padding-top: 1rem;
      border-top: 1px solid #e5e7eb;
      text-align: center;
      font-size: 0.75rem;
      color: #9ca3af;
    }

    @media print {
      body { padding: 0; max-width: none; }
      .section { page-break-inside: avoid; break-inside: avoid; }
      .report-header { page-break-after: avoid; }
      .toc { page-break-after: always; }
      @page { margin: 2cm; }
    }
  `
}
