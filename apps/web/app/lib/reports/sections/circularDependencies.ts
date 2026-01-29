import type { CircularDependencyReport } from '../types'

/**
 * Render circular dependencies section as HTML
 * AC6: Circular Dependencies
 */
export function renderCircularDepsHtml(data: CircularDependencyReport): string {
  if (data.totalCount === 0) {
    return `
    <section class="section">
      <div class="section-header">
        <h2>Circular Dependencies</h2>
        <span class="badge">0 found</span>
      </div>
      <div class="section-content">
        <p>No circular dependencies detected.</p>
      </div>
    </section>`
  }

  const severityRows = `
    <tr><td>Critical</td><td class="severity-critical">${data.bySeverity.critical}</td></tr>
    <tr><td>High</td><td class="severity-high">${data.bySeverity.high}</td></tr>
    <tr><td>Medium</td><td class="severity-medium">${data.bySeverity.medium}</td></tr>
    <tr><td>Low</td><td class="severity-low">${data.bySeverity.low}</td></tr>`

  const cycleRows = data.cycles
    .map(
      (cycle) => `
    <tr>
      <td>${escapeHtml(cycle.id)}</td>
      <td><code>${escapeHtml(cycle.packages.join(' → '))} → ${escapeHtml(cycle.packages[0])}</code></td>
      <td class="severity-${mapSeverityClass(cycle.severity)}">${escapeHtml(cycle.severity)}</td>
      <td>${escapeHtml(cycle.type)}</td>
    </tr>`
    )
    .join('')

  return `
    <section class="section">
      <div class="section-header">
        <h2>Circular Dependencies</h2>
        <span class="badge">${data.totalCount} found</span>
      </div>
      <div class="section-content">
        <table>
          <thead><tr><th>Severity</th><th>Count</th></tr></thead>
          <tbody>${severityRows}</tbody>
        </table>
        <h3 style="margin-top: 1.5rem; margin-bottom: 1rem;">Detected Cycles</h3>
        <table>
          <thead><tr><th>ID</th><th>Path</th><th>Severity</th><th>Type</th></tr></thead>
          <tbody>${cycleRows}</tbody>
        </table>
      </div>
    </section>`
}

/**
 * Render circular dependencies section as Markdown
 * AC6: Circular Dependencies
 */
export function renderCircularDepsMd(data: CircularDependencyReport): string {
  let md = `## Circular Dependencies\n\n`
  md += `**Total: ${data.totalCount}**\n\n`

  if (data.totalCount === 0) {
    md += `> :tada: No circular dependencies detected!\n\n`
    return md
  }

  md += `### Summary by Severity\n\n`
  md += `| Severity | Count |\n`
  md += `|----------|-------|\n`
  md += `| Critical | ${data.bySeverity.critical} |\n`
  md += `| High | ${data.bySeverity.high} |\n`
  md += `| Medium | ${data.bySeverity.medium} |\n`
  md += `| Low | ${data.bySeverity.low} |\n\n`

  md += `### Detected Cycles\n\n`

  for (const cycle of data.cycles) {
    md += `#### Cycle: ${cycle.id}\n`
    md += `- **Severity:** ${cycle.severity}\n`
    md += `- **Type:** ${cycle.type}\n`
    md += `- **Path:** \`${cycle.packages.join(' → ')} → ${cycle.packages[0]}\`\n\n`
  }

  return md
}

function mapSeverityClass(severity: string): string {
  switch (severity) {
    case 'critical':
      return 'critical'
    case 'warning':
      return 'high'
    case 'info':
      return 'medium'
    default:
      return 'low'
  }
}

function escapeHtml(str: string): string {
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;')
}
