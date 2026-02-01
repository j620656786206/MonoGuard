import type { CircularDependencyInfo } from '@monoguard/types'
import { describe, expect, it } from 'vitest'
import { renderRootCauseAnalysis } from '../sections/rootCauseAnalysis'

const baseCycle: CircularDependencyInfo = {
  cycle: ['pkg-a', 'pkg-b', 'pkg-c', 'pkg-a'],
  type: 'indirect',
  severity: 'warning',
  depth: 3,
  impact: 'Test',
  complexity: 5,
  priorityScore: 5,
}

describe('renderRootCauseAnalysis', () => {
  it('should return fallback when no root cause available', () => {
    const details = renderRootCauseAnalysis(baseCycle)
    expect(details.explanation).toContain('not available')
    expect(details.confidenceScore).toBe(0)
    expect(details.originatingPackage).toBe('pkg-a')
    expect(details.alternativeCandidates).toHaveLength(0)
  })

  it('should render root cause from analysis', () => {
    const withRootCause: CircularDependencyInfo = {
      ...baseCycle,
      rootCause: {
        originatingPackage: 'pkg-a',
        problematicDependency: { from: 'pkg-a', to: 'pkg-b', type: 'production', critical: true },
        confidence: 90,
        explanation: 'pkg-a directly imports from pkg-b creating the cycle',
        chain: [
          { from: 'pkg-a', to: 'pkg-b', type: 'production', critical: true },
          { from: 'pkg-b', to: 'pkg-c', type: 'production', critical: false },
          { from: 'pkg-c', to: 'pkg-a', type: 'production', critical: false },
        ],
      },
    }
    const details = renderRootCauseAnalysis(withRootCause)
    expect(details.explanation).toBe('pkg-a directly imports from pkg-b creating the cycle')
    expect(details.confidenceScore).toBe(90)
    expect(details.originatingPackage).toBe('pkg-a')
  })

  it('should include alternative candidates when confidence < 80%', () => {
    const lowConfidence: CircularDependencyInfo = {
      ...baseCycle,
      rootCause: {
        originatingPackage: 'pkg-a',
        problematicDependency: { from: 'pkg-a', to: 'pkg-b', type: 'production', critical: true },
        confidence: 60,
        explanation: 'Uncertain root cause',
        chain: [
          { from: 'pkg-a', to: 'pkg-b', type: 'production', critical: true },
          { from: 'pkg-b', to: 'pkg-c', type: 'production', critical: true },
          { from: 'pkg-c', to: 'pkg-a', type: 'production', critical: false },
        ],
      },
    }
    const details = renderRootCauseAnalysis(lowConfidence)
    expect(details.confidenceScore).toBe(60)
    expect(details.alternativeCandidates.length).toBeGreaterThan(0)
    expect(details.alternativeCandidates[0].package).toBe('pkg-b')
  })

  it('should NOT include alternatives when confidence >= 80%', () => {
    const highConfidence: CircularDependencyInfo = {
      ...baseCycle,
      rootCause: {
        originatingPackage: 'pkg-a',
        problematicDependency: { from: 'pkg-a', to: 'pkg-b', type: 'production', critical: true },
        confidence: 90,
        explanation: 'High confidence result',
        chain: [
          { from: 'pkg-a', to: 'pkg-b', type: 'production', critical: true },
          { from: 'pkg-b', to: 'pkg-c', type: 'production', critical: true },
        ],
      },
    }
    const details = renderRootCauseAnalysis(highConfidence)
    expect(details.alternativeCandidates).toHaveLength(0)
  })

  it('should build code references from import traces', () => {
    const withTraces: CircularDependencyInfo = {
      ...baseCycle,
      rootCause: {
        originatingPackage: 'pkg-a',
        problematicDependency: { from: 'pkg-a', to: 'pkg-b', type: 'production', critical: true },
        confidence: 85,
        explanation: 'Test',
        chain: [],
      },
      importTraces: [
        {
          fromPackage: 'pkg-a',
          toPackage: 'pkg-b',
          filePath: 'src/index.ts',
          lineNumber: 1,
          statement: "import { foo } from 'pkg-b'",
          importType: 'esm-named',
        },
      ],
    }
    const details = renderRootCauseAnalysis(withTraces)
    expect(details.codeReferences).toHaveLength(1)
    expect(details.codeReferences[0].file).toBe('src/index.ts')
    expect(details.codeReferences[0].line).toBe(1)
  })

  it('should include originating reason with dependency info', () => {
    const withRootCause: CircularDependencyInfo = {
      ...baseCycle,
      rootCause: {
        originatingPackage: 'pkg-a',
        problematicDependency: { from: 'pkg-a', to: 'pkg-b', type: 'production', critical: true },
        confidence: 85,
        explanation: 'Test',
        chain: [],
      },
    }
    const details = renderRootCauseAnalysis(withRootCause)
    expect(details.originatingReason).toContain('pkg-a')
    expect(details.originatingReason).toContain('pkg-b')
  })
})
