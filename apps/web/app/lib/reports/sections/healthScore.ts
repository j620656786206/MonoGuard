import type { HealthScoreReport } from '../types'

/**
 * Render health score section as HTML
 * AC5: Health Score Summary
 */
export function renderHealthScoreHtml(data: HealthScoreReport): string {
  const breakdownItems = data.breakdown
    .map(
      (item) => `
      <div class="breakdown-item">
        <div class="label">${escapeHtml(item.category)}</div>
        <div class="value">${item.score}</div>
        <div class="label">${item.weight}% weight</div>
      </div>`
    )
    .join('')

  return `
    <section class="section">
      <div class="section-header">
        <h2>Health Score</h2>
        <span class="badge">${data.overall}/100</span>
      </div>
      <div class="section-content">
        <div class="health-score ${data.rating}">
          <div class="score">${data.overall}</div>
          <div class="rating">${data.rating}</div>
        </div>
        <div class="breakdown-grid">
          ${breakdownItems}
        </div>
      </div>
    </section>`
}

/**
 * Render health score section as Markdown
 * AC5: Health Score Summary
 */
export function renderHealthScoreMd(data: HealthScoreReport): string {
  const ratingEmoji: Record<string, string> = {
    excellent: ':white_check_mark:',
    good: ':heavy_check_mark:',
    fair: ':warning:',
    poor: ':x:',
    critical: ':rotating_light:',
  }

  let md = `## Health Score\n\n`
  md += `**Overall Score: ${data.overall}/100** ${ratingEmoji[data.rating] ?? ''} ${data.rating.toUpperCase()}\n\n`

  md += `### Score Breakdown\n\n`
  md += `| Category | Score | Weight |\n`
  md += `|----------|-------|--------|\n`

  for (const item of data.breakdown) {
    md += `| ${item.category} | ${item.score} | ${item.weight}% |\n`
  }

  md += '\n'
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
