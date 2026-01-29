import { describe, expect, it } from 'vitest'
import { renderCircularDepsHtml, renderCircularDepsMd } from '../sections/circularDependencies'
import {
  renderFixRecommendationsHtml,
  renderFixRecommendationsMd,
} from '../sections/fixRecommendations'
import { renderHealthScoreHtml, renderHealthScoreMd } from '../sections/healthScore'
import { renderVersionConflictsHtml, renderVersionConflictsMd } from '../sections/versionConflicts'
import type {
  CircularDependencyReport,
  FixRecommendationReport,
  HealthScoreReport,
  VersionConflictReport,
} from '../types'

const mockHealthScore: HealthScoreReport = {
  overall: 75,
  breakdown: [
    { category: 'Dependencies', score: 80, weight: 40 },
    { category: 'Architecture', score: 70, weight: 30 },
    { category: 'Security', score: 75, weight: 30 },
  ],
  rating: 'good',
  ratingThresholds: { excellent: 85, good: 70, fair: 50, poor: 30, critical: 0 },
}

const mockCircularDeps: CircularDependencyReport = {
  totalCount: 2,
  bySeverity: { critical: 1, high: 1, medium: 0, low: 0 },
  cycles: [
    { id: 'cycle-1', packages: ['pkg-a', 'pkg-b'], severity: 'critical', type: 'direct' },
    { id: 'cycle-2', packages: ['pkg-c', 'pkg-d', 'pkg-e'], severity: 'warning', type: 'indirect' },
  ],
}

const mockVersionConflicts: VersionConflictReport = {
  totalCount: 2,
  byRiskLevel: { critical: 0, high: 1, medium: 1, low: 0 },
  conflicts: [
    {
      packageName: 'lodash',
      versions: ['4.17.21', '4.17.15'],
      riskLevel: 'warning',
      recommendedVersion: 'Upgrade all to 4.17.21',
    },
    {
      packageName: 'react',
      versions: ['18.2.0', '17.0.2'],
      riskLevel: 'critical',
      recommendedVersion: 'Upgrade all to 18.2.0',
    },
  ],
}

const mockFixRecommendations: FixRecommendationReport = {
  totalCount: 2,
  quickWins: 1,
  recommendations: [
    {
      id: 'fix-1',
      title: 'Extract Shared Module',
      description: 'Move shared code into a new package',
      effort: 'low',
      impact: 'high',
      priority: 9,
      affectedPackages: ['pkg-a', 'pkg-b'],
    },
    {
      id: 'fix-2',
      title: 'Restructure Boundaries',
      description: 'Reorganize package boundaries',
      effort: 'high',
      impact: 'medium',
      priority: 5,
      affectedPackages: ['pkg-c'],
    },
  ],
}

describe('Section Renderers', () => {
  describe('Health Score HTML', () => {
    it('should render health score section', () => {
      const html = renderHealthScoreHtml(mockHealthScore)
      expect(html).toContain('Health Score')
      expect(html).toContain('75')
      expect(html).toContain('good')
    })

    it('should render breakdown items', () => {
      const html = renderHealthScoreHtml(mockHealthScore)
      expect(html).toContain('Dependencies')
      expect(html).toContain('Architecture')
      expect(html).toContain('Security')
    })

    it('should include rating as CSS class', () => {
      const html = renderHealthScoreHtml(mockHealthScore)
      expect(html).toContain('class="health-score good"')
    })
  })

  describe('Health Score Markdown', () => {
    it('should render overall score with rating', () => {
      const md = renderHealthScoreMd(mockHealthScore)
      expect(md).toContain('**Overall Score: 75/100**')
      expect(md).toContain('GOOD')
    })

    it('should render GFM table', () => {
      const md = renderHealthScoreMd(mockHealthScore)
      expect(md).toContain('| Category | Score | Weight |')
      expect(md).toContain('|----------|-------|--------|')
      expect(md).toContain('| Dependencies | 80 | 40% |')
    })

    it('should include emoji for rating', () => {
      const md = renderHealthScoreMd(mockHealthScore)
      expect(md).toContain(':heavy_check_mark:')
    })
  })

  describe('Circular Dependencies HTML', () => {
    it('should render cycle information', () => {
      const html = renderCircularDepsHtml(mockCircularDeps)
      expect(html).toContain('Circular Dependencies')
      expect(html).toContain('2 found')
    })

    it('should render severity summary table', () => {
      const html = renderCircularDepsHtml(mockCircularDeps)
      expect(html).toContain('Critical')
      expect(html).toContain('severity-critical')
    })

    it('should render cycle paths', () => {
      const html = renderCircularDepsHtml(mockCircularDeps)
      expect(html).toContain('pkg-a → pkg-b → pkg-a')
    })

    it('should handle zero circular dependencies', () => {
      const empty: CircularDependencyReport = {
        totalCount: 0,
        bySeverity: { critical: 0, high: 0, medium: 0, low: 0 },
        cycles: [],
      }
      const html = renderCircularDepsHtml(empty)
      expect(html).toContain('No circular dependencies detected')
      expect(html).toContain('0 found')
    })
  })

  describe('Circular Dependencies Markdown', () => {
    it('should render total count', () => {
      const md = renderCircularDepsMd(mockCircularDeps)
      expect(md).toContain('**Total: 2**')
    })

    it('should render severity table', () => {
      const md = renderCircularDepsMd(mockCircularDeps)
      expect(md).toContain('| Severity | Count |')
      expect(md).toContain('| Critical | 1 |')
    })

    it('should render cycle details', () => {
      const md = renderCircularDepsMd(mockCircularDeps)
      expect(md).toContain('#### Cycle: cycle-1')
      expect(md).toContain('`pkg-a → pkg-b → pkg-a`')
    })

    it('should show celebration message when zero', () => {
      const empty: CircularDependencyReport = {
        totalCount: 0,
        bySeverity: { critical: 0, high: 0, medium: 0, low: 0 },
        cycles: [],
      }
      const md = renderCircularDepsMd(empty)
      expect(md).toContain(':tada:')
    })
  })

  describe('Version Conflicts HTML', () => {
    it('should render conflict table', () => {
      const html = renderVersionConflictsHtml(mockVersionConflicts)
      expect(html).toContain('Version Conflicts')
      expect(html).toContain('2 found')
      expect(html).toContain('lodash')
      expect(html).toContain('4.17.21, 4.17.15')
    })

    it('should handle zero conflicts', () => {
      const empty: VersionConflictReport = {
        totalCount: 0,
        byRiskLevel: { critical: 0, high: 0, medium: 0, low: 0 },
        conflicts: [],
      }
      const html = renderVersionConflictsHtml(empty)
      expect(html).toContain('No version conflicts detected')
    })
  })

  describe('Version Conflicts Markdown', () => {
    it('should render GFM table', () => {
      const md = renderVersionConflictsMd(mockVersionConflicts)
      expect(md).toContain('| Package | Conflicting Versions | Risk | Recommended |')
      expect(md).toContain('`lodash`')
      expect(md).toContain('4.17.21, 4.17.15')
    })

    it('should show check mark when zero', () => {
      const empty: VersionConflictReport = {
        totalCount: 0,
        byRiskLevel: { critical: 0, high: 0, medium: 0, low: 0 },
        conflicts: [],
      }
      const md = renderVersionConflictsMd(empty)
      expect(md).toContain(':white_check_mark:')
    })
  })

  describe('Fix Recommendations HTML', () => {
    it('should render fix cards', () => {
      const html = renderFixRecommendationsHtml(mockFixRecommendations)
      expect(html).toContain('Fix Recommendations')
      expect(html).toContain('Extract Shared Module')
      expect(html).toContain('Restructure Boundaries')
    })

    it('should mark quick wins', () => {
      const html = renderFixRecommendationsHtml(mockFixRecommendations)
      expect(html).toContain('quick-win')
      expect(html).toContain('Quick Win')
    })

    it('should show effort and impact', () => {
      const html = renderFixRecommendationsHtml(mockFixRecommendations)
      expect(html).toContain('Effort: low')
      expect(html).toContain('Impact: high')
    })

    it('should handle zero recommendations', () => {
      const empty: FixRecommendationReport = {
        totalCount: 0,
        quickWins: 0,
        recommendations: [],
      }
      const html = renderFixRecommendationsHtml(empty)
      expect(html).toContain('No fix recommendations')
    })
  })

  describe('Fix Recommendations Markdown', () => {
    it('should render recommendations list', () => {
      const md = renderFixRecommendationsMd(mockFixRecommendations)
      expect(md).toContain('#### 9. Extract Shared Module :zap: Quick Win')
      expect(md).toContain('- **Effort:** low')
      expect(md).toContain('- **Impact:** high')
    })

    it('should show quick wins count', () => {
      const md = renderFixRecommendationsMd(mockFixRecommendations)
      expect(md).toContain('**Quick Wins: 1** :zap:')
    })

    it('should render affected packages with backticks', () => {
      const md = renderFixRecommendationsMd(mockFixRecommendations)
      expect(md).toContain('`pkg-a`')
      expect(md).toContain('`pkg-b`')
    })
  })
})
