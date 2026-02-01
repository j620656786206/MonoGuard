import { describe, expect, it } from 'vitest'
import { getDiagnosticHtmlTemplate } from '../templates/diagnosticHtmlTemplate'
import type { DiagnosticReport } from '../types'

function createMinimalReport(overrides: Partial<DiagnosticReport> = {}): DiagnosticReport {
  return {
    id: 'diag-test-1',
    cycleId: 'pkg-a-pkg-b',
    generatedAt: '2026-01-25T10:00:00Z',
    monoguardVersion: '0.1.0',
    projectName: 'test-project',
    executiveSummary: {
      description: 'Circular dependency between pkg-a and pkg-b',
      severity: 'high',
      recommendation: 'Extract shared code into a common module',
      estimatedEffort: 'medium',
      affectedPackagesCount: 4,
      cycleLength: 2,
    },
    cyclePath: {
      nodes: [
        {
          id: 'pkg-a',
          name: 'pkg-a',
          path: '/packages/pkg-a',
          isInCycle: true,
          position: { x: 100, y: 100 },
        },
        {
          id: 'pkg-b',
          name: 'pkg-b',
          path: '/packages/pkg-b',
          isInCycle: true,
          position: { x: 300, y: 100 },
        },
      ],
      edges: [
        {
          from: 'pkg-a',
          to: 'pkg-b',
          isBreakingPoint: true,
          importStatement: "import { foo } from 'pkg-b'",
        },
        { from: 'pkg-b', to: 'pkg-a', isBreakingPoint: false },
      ],
      breakingPoint: {
        fromPackage: 'pkg-a',
        toPackage: 'pkg-b',
        reason: 'Lowest edge count',
      },
      svgDiagram: '<svg><circle r="10"/></svg>',
      asciiDiagram: 'pkg-a -> pkg-b -> pkg-a',
    },
    rootCause: {
      explanation: 'pkg-a depends on pkg-b which creates a cycle',
      confidenceScore: 85,
      originatingPackage: 'pkg-a',
      originatingReason: 'pkg-a has the most outgoing dependencies',
      alternativeCandidates: [],
      codeReferences: [],
    },
    fixStrategies: [],
    impactAssessment: {
      directParticipants: ['pkg-a', 'pkg-b'],
      directParticipantsCount: 2,
      indirectDependents: ['pkg-c', 'pkg-d'],
      indirectDependentsCount: 2,
      totalAffectedCount: 4,
      percentageOfMonorepo: 40,
      riskLevel: 'high',
      riskExplanation: 'High risk: 4 packages affected',
      rippleEffectTree: { package: 'Cycle', depth: 0, dependents: [] },
    },
    relatedCycles: [],
    metadata: {
      generatedAt: '2026-01-25T10:00:00Z',
      generationDurationMs: 42,
      monoguardVersion: '0.1.0',
      projectName: 'test-project',
      analysisConfigHash: 'abc123',
    },
    ...overrides,
  }
}

describe('getDiagnosticHtmlTemplate', () => {
  it('should render fix strategies with code snippets (before/after)', () => {
    const report = createMinimalReport({
      fixStrategies: [
        {
          strategy: 'extract-module',
          title: 'Extract Shared Module',
          description: 'Move shared code to new package',
          suitabilityScore: 8,
          estimatedEffort: 'medium',
          estimatedTime: '1-2 hours',
          pros: ['Clean separation', 'Testable'],
          cons: ['New package to maintain'],
          steps: [
            {
              number: 1,
              title: 'Create package',
              description: 'Create a new shared package',
              codeSnippet: 'mkdir packages/shared',
              isOptional: false,
            },
          ],
          codeSnippets: {
            before: "import { foo } from 'pkg-b'",
            after: "import { foo } from 'shared'",
          },
        },
      ],
    })

    const html = getDiagnosticHtmlTemplate(report)

    expect(html).toContain('Extract Shared Module')
    expect(html).toContain('Code Changes')
    expect(html).toContain('Before:')
    expect(html).toContain('After:')
    expect(html).toContain('import { foo } from &#039;pkg-b&#039;')
    expect(html).toContain('import { foo } from &#039;shared&#039;')
    expect(html).toContain('mkdir packages/shared')
  })

  it('should render fix strategies with only before snippet', () => {
    const report = createMinimalReport({
      fixStrategies: [
        {
          strategy: 'extract-module',
          title: 'Extract Module',
          description: 'Test',
          suitabilityScore: 7,
          estimatedEffort: 'low',
          estimatedTime: '30 minutes',
          pros: ['Easy'],
          cons: [],
          steps: [],
          codeSnippets: {
            before: "import { bar } from 'old'",
            after: '',
          },
        },
      ],
    })

    const html = getDiagnosticHtmlTemplate(report)

    expect(html).toContain('Before:')
    expect(html).toContain('import { bar } from &#039;old&#039;')
  })

  it('should render fix strategies with only after snippet', () => {
    const report = createMinimalReport({
      fixStrategies: [
        {
          strategy: 'dependency-injection',
          title: 'DI Refactor',
          description: 'Use dependency injection',
          suitabilityScore: 6,
          estimatedEffort: 'high',
          estimatedTime: '2-4 hours',
          pros: [],
          cons: [],
          steps: [],
          codeSnippets: {
            before: '',
            after: "import { bar } from 'new-package'",
          },
        },
      ],
    })

    const html = getDiagnosticHtmlTemplate(report)

    expect(html).toContain('After:')
    expect(html).toContain('import { bar } from &#039;new-package&#039;')
  })

  it('should not render Code Changes when no snippets', () => {
    const report = createMinimalReport({
      fixStrategies: [
        {
          strategy: 'extract-module',
          title: 'Simple Strategy',
          description: 'No code snippets',
          suitabilityScore: 5,
          estimatedEffort: 'low',
          estimatedTime: '15 minutes',
          pros: ['Fast'],
          cons: [],
          steps: [],
          codeSnippets: { before: '', after: '' },
        },
      ],
    })

    const html = getDiagnosticHtmlTemplate(report)

    expect(html).toContain('Simple Strategy')
    expect(html).not.toContain('Code Changes')
  })

  it('should render related cycles section with data', () => {
    const report = createMinimalReport({
      relatedCycles: [
        {
          cycleId: 'cycle-x-y',
          sharedPackages: ['pkg-a', 'pkg-c'],
          overlapPercentage: 50,
          recommendFixTogether: true,
          reason: 'High overlap suggests coupled cycles',
        },
        {
          cycleId: 'cycle-z',
          sharedPackages: ['pkg-b'],
          overlapPercentage: 10,
          recommendFixTogether: false,
        },
      ],
    })

    const html = getDiagnosticHtmlTemplate(report)

    expect(html).toContain('cycle-x-y')
    expect(html).toContain('cycle-z')
    expect(html).toContain('50%')
    expect(html).toContain('10%')
    expect(html).toContain('✅ Yes')
    expect(html).toContain('High overlap suggests coupled cycles')
    expect(html).toContain('Shared Packages')
    expect(html).toContain('Fix Together?')
  })

  it('should render empty fix strategies section', () => {
    const report = createMinimalReport({ fixStrategies: [] })
    const html = getDiagnosticHtmlTemplate(report)
    expect(html).toContain('No fix strategies available for this cycle.')
  })

  it('should render empty related cycles section', () => {
    const report = createMinimalReport({ relatedCycles: [] })
    const html = getDiagnosticHtmlTemplate(report)
    expect(html).toContain('No related cycles detected.')
  })

  it('should render root cause with alternative candidates', () => {
    const report = createMinimalReport({
      rootCause: {
        explanation: 'Uncertain root cause',
        confidenceScore: 60,
        originatingPackage: 'pkg-a',
        originatingReason: 'Most dependencies',
        alternativeCandidates: [{ package: 'pkg-b', reason: 'Also has many deps', confidence: 45 }],
        codeReferences: [],
      },
    })

    const html = getDiagnosticHtmlTemplate(report)

    expect(html).toContain('Alternative Candidates')
    expect(html).toContain('pkg-b')
    expect(html).toContain('Also has many deps')
    expect(html).toContain('45%')
  })

  it('should render root cause with code references', () => {
    const report = createMinimalReport({
      rootCause: {
        explanation: 'Root cause identified',
        confidenceScore: 90,
        originatingPackage: 'pkg-a',
        originatingReason: 'Direct import creates cycle',
        alternativeCandidates: [],
        codeReferences: [
          { file: 'src/index.ts', line: 42, importStatement: "import { foo } from 'pkg-b'" },
        ],
      },
    })

    const html = getDiagnosticHtmlTemplate(report)

    expect(html).toContain('Code References')
    expect(html).toContain('src/index.ts:42')
    expect(html).toContain('import { foo } from &#039;pkg-b&#039;')
  })

  it('should render impact assessment with indirect dependents', () => {
    const report = createMinimalReport({
      impactAssessment: {
        directParticipants: ['pkg-a', 'pkg-b'],
        directParticipantsCount: 2,
        indirectDependents: ['pkg-c', 'pkg-d', 'pkg-e'],
        indirectDependentsCount: 3,
        totalAffectedCount: 5,
        percentageOfMonorepo: 50,
        riskLevel: 'critical',
        riskExplanation: 'Critical risk',
        rippleEffectTree: { package: 'Cycle', depth: 0, dependents: [] },
      },
    })

    const html = getDiagnosticHtmlTemplate(report)

    expect(html).toContain('Indirect Dependents')
    expect(html).toContain('pkg-c')
    expect(html).toContain('pkg-d')
    expect(html).toContain('pkg-e')
  })

  it('should not render indirect dependents section when empty', () => {
    const report = createMinimalReport({
      impactAssessment: {
        directParticipants: ['pkg-a'],
        directParticipantsCount: 1,
        indirectDependents: [],
        indirectDependentsCount: 0,
        totalAffectedCount: 1,
        percentageOfMonorepo: 10,
        riskLevel: 'low',
        riskExplanation: 'Low risk',
        rippleEffectTree: { package: 'Cycle', depth: 0, dependents: [] },
      },
    })

    const html = getDiagnosticHtmlTemplate(report)

    // Should contain Direct Participants but NOT Indirect Dependents heading
    expect(html).toContain('Direct Participants')
    // It should appear in the metric grid but not as a subheading with list
    const directListIndex = html.indexOf('<h3>Indirect Dependents</h3>')
    expect(directListIndex).toBe(-1)
  })

  it('should render related cycle without reason showing dash', () => {
    const report = createMinimalReport({
      relatedCycles: [
        {
          cycleId: 'other-cycle',
          sharedPackages: ['pkg-a'],
          overlapPercentage: 20,
          recommendFixTogether: false,
        },
      ],
    })

    const html = getDiagnosticHtmlTemplate(report)

    // When reason is undefined, should show '-'
    expect(html).toContain('<td>-</td>')
  })

  it('should escape HTML entities in user content', () => {
    const report = createMinimalReport({
      projectName: '<script>alert("xss")</script>',
    })

    const html = getDiagnosticHtmlTemplate(report)

    expect(html).not.toContain('<script>alert')
    expect(html).toContain('&lt;script&gt;alert')
  })
})
