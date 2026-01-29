import type { FixRecommendationReport } from '../types'

/**
 * Render fix recommendations section as HTML
 * AC8: Fix Recommendations
 */
export function renderFixRecommendationsHtml(data: FixRecommendationReport): string {
  if (data.totalCount === 0) {
    return `
    <section class="section">
      <div class="section-header">
        <h2>Fix Recommendations</h2>
        <span class="badge">None</span>
      </div>
      <div class="section-content">
        <p>No fix recommendations at this time.</p>
      </div>
    </section>`
  }

  const fixCards = data.recommendations
    .map((rec) => {
      const isQuickWin = rec.effort === 'low' && rec.impact === 'high'
      return `
      <div class="fix-card${isQuickWin ? ' quick-win' : ''}">
        <div class="title">
          ${rec.priority}. ${escapeHtml(rec.title)}
          ${isQuickWin ? '<span style="color: var(--color-success); margin-left: 0.5rem;">⚡ Quick Win</span>' : ''}
        </div>
        <p>${escapeHtml(rec.description)}</p>
        <div class="meta">
          <span>Effort: ${escapeHtml(rec.effort)}</span>
          <span>Impact: ${escapeHtml(rec.impact)}</span>
          <span>Packages: ${rec.affectedPackages.map((p) => `<code>${escapeHtml(p)}</code>`).join(', ')}</span>
        </div>
      </div>`
    })
    .join('')

  return `
    <section class="section">
      <div class="section-header">
        <h2>Fix Recommendations</h2>
        <span class="badge">${data.totalCount} total, ${data.quickWins} quick wins</span>
      </div>
      <div class="section-content">
        ${fixCards}
      </div>
    </section>`
}

/**
 * Render fix recommendations section as Markdown
 * AC8: Fix Recommendations
 */
export function renderFixRecommendationsMd(data: FixRecommendationReport): string {
  let md = `## Fix Recommendations\n\n`
  md += `**Total Recommendations: ${data.totalCount}**\n`
  md += `**Quick Wins: ${data.quickWins}** :zap:\n\n`

  if (data.totalCount === 0) {
    md += `> No fix recommendations at this time.\n\n`
    return md
  }

  md += `### Priority Fixes\n\n`

  for (const rec of data.recommendations) {
    const quickWinBadge = rec.effort === 'low' && rec.impact === 'high' ? ' :zap: Quick Win' : ''
    md += `#### ${rec.priority}. ${rec.title}${quickWinBadge}\n\n`
    md += `${rec.description}\n\n`
    md += `- **Effort:** ${rec.effort}\n`
    md += `- **Impact:** ${rec.impact}\n`
    md += `- **Affected Packages:** ${rec.affectedPackages.map((p) => `\`${p}\``).join(', ')}\n\n`
  }

  return md
}

function escapeHtml(str: string): string {
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;')
}
