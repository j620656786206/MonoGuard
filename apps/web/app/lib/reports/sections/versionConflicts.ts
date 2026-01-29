import type { VersionConflictReport } from '../types'

/**
 * Render version conflicts section as HTML
 * AC7: Version Conflicts
 */
export function renderVersionConflictsHtml(data: VersionConflictReport): string {
  if (data.totalCount === 0) {
    return `
    <section class="section">
      <div class="section-header">
        <h2>Version Conflicts</h2>
        <span class="badge">0 found</span>
      </div>
      <div class="section-content">
        <p>No version conflicts detected.</p>
      </div>
    </section>`
  }

  const conflictRows = data.conflicts
    .map(
      (c) => `
    <tr>
      <td><code>${escapeHtml(c.packageName)}</code></td>
      <td>${escapeHtml(c.versions.join(', '))}</td>
      <td class="severity-${mapRiskClass(c.riskLevel)}">${escapeHtml(c.riskLevel)}</td>
      <td>${escapeHtml(c.recommendedVersion)}</td>
    </tr>`
    )
    .join('')

  return `
    <section class="section">
      <div class="section-header">
        <h2>Version Conflicts</h2>
        <span class="badge">${data.totalCount} found</span>
      </div>
      <div class="section-content">
        <table>
          <thead><tr><th>Package</th><th>Conflicting Versions</th><th>Risk</th><th>Recommended</th></tr></thead>
          <tbody>${conflictRows}</tbody>
        </table>
      </div>
    </section>`
}

/**
 * Render version conflicts section as Markdown
 * AC7: Version Conflicts
 */
export function renderVersionConflictsMd(data: VersionConflictReport): string {
  let md = `## Version Conflicts\n\n`
  md += `**Total: ${data.totalCount}**\n\n`

  if (data.totalCount === 0) {
    md += `> :white_check_mark: No version conflicts detected!\n\n`
    return md
  }

  md += `| Package | Conflicting Versions | Risk | Recommended |\n`
  md += `|---------|---------------------|------|-------------|\n`

  for (const conflict of data.conflicts) {
    md += `| \`${conflict.packageName}\` | ${conflict.versions.join(', ')} | ${conflict.riskLevel} | ${conflict.recommendedVersion} |\n`
  }

  md += '\n'
  return md
}

function mapRiskClass(riskLevel: string): string {
  switch (riskLevel) {
    case 'critical':
      return 'critical'
    case 'high':
    case 'warning':
      return 'high'
    case 'medium':
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
