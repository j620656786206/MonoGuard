import type { CircularDependencyInfo } from '@monoguard/types'
import { describe, expect, it } from 'vitest'
import { findRelatedCycles } from '../sections/relatedCycles'

const baseCycleFields = {
  type: 'indirect' as const,
  severity: 'warning' as const,
  depth: 3,
  impact: 'Test',
  complexity: 5,
  priorityScore: 5,
}

describe('findRelatedCycles', () => {
  it('should find cycles sharing packages', () => {
    const targetCycle: CircularDependencyInfo = {
      ...baseCycleFields,
      cycle: ['pkg-a', 'pkg-b', 'pkg-c', 'pkg-a'],
    }
    const relatedCycle: CircularDependencyInfo = {
      ...baseCycleFields,
      cycle: ['pkg-b', 'pkg-d', 'pkg-e', 'pkg-b'],
    }
    const unrelatedCycle: CircularDependencyInfo = {
      ...baseCycleFields,
      cycle: ['pkg-x', 'pkg-y', 'pkg-x'],
    }
    const allCycles = [targetCycle, relatedCycle, unrelatedCycle]

    const related = findRelatedCycles(targetCycle, allCycles)
    expect(related).toHaveLength(1)
    expect(related[0].sharedPackages).toContain('pkg-b')
  })

  it('should skip the target cycle itself', () => {
    const targetCycle: CircularDependencyInfo = {
      ...baseCycleFields,
      cycle: ['pkg-a', 'pkg-b', 'pkg-a'],
    }
    const related = findRelatedCycles(targetCycle, [targetCycle])
    expect(related).toHaveLength(0)
  })

  it('should return empty array when no related cycles', () => {
    const targetCycle: CircularDependencyInfo = {
      ...baseCycleFields,
      cycle: ['pkg-a', 'pkg-b', 'pkg-a'],
    }
    const otherCycle: CircularDependencyInfo = {
      ...baseCycleFields,
      cycle: ['pkg-x', 'pkg-y', 'pkg-x'],
    }
    const related = findRelatedCycles(targetCycle, [targetCycle, otherCycle])
    expect(related).toHaveLength(0)
  })

  it('should calculate overlap percentage', () => {
    const targetCycle: CircularDependencyInfo = {
      ...baseCycleFields,
      cycle: ['pkg-a', 'pkg-b', 'pkg-c', 'pkg-a'],
    }
    const overlapping: CircularDependencyInfo = {
      ...baseCycleFields,
      cycle: ['pkg-a', 'pkg-b', 'pkg-d', 'pkg-a'],
    }
    const related = findRelatedCycles(targetCycle, [targetCycle, overlapping])
    expect(related).toHaveLength(1)
    expect(related[0].overlapPercentage).toBeGreaterThan(0)
    expect(related[0].sharedPackages).toEqual(expect.arrayContaining(['pkg-a', 'pkg-b']))
  })

  it('should recommend fixing together when overlap >= 30%', () => {
    const targetCycle: CircularDependencyInfo = {
      ...baseCycleFields,
      cycle: ['pkg-a', 'pkg-b', 'pkg-c', 'pkg-a'],
    }
    const highOverlap: CircularDependencyInfo = {
      ...baseCycleFields,
      cycle: ['pkg-a', 'pkg-b', 'pkg-d', 'pkg-a'],
    }
    const related = findRelatedCycles(targetCycle, [targetCycle, highOverlap])
    expect(related[0].recommendFixTogether).toBe(true)
    expect(related[0].reason).toContain('Fixing them together')
  })

  it('should not recommend fixing together for low overlap', () => {
    const targetCycle: CircularDependencyInfo = {
      ...baseCycleFields,
      cycle: ['a', 'b', 'c', 'd', 'e', 'f', 'a'],
    }
    const lowOverlap: CircularDependencyInfo = {
      ...baseCycleFields,
      cycle: ['a', 'x', 'y', 'z', 'w', 'v', 'a'],
    }
    const related = findRelatedCycles(targetCycle, [targetCycle, lowOverlap])
    expect(related).toHaveLength(1)
    expect(related[0].overlapPercentage).toBeLessThan(30)
    expect(related[0].recommendFixTogether).toBe(false)
  })

  it('should include cycle ID in results', () => {
    const targetCycle: CircularDependencyInfo = {
      ...baseCycleFields,
      cycle: ['pkg-a', 'pkg-b', 'pkg-a'],
    }
    const relatedCycle: CircularDependencyInfo = {
      ...baseCycleFields,
      cycle: ['pkg-a', 'pkg-c', 'pkg-a'],
    }
    const related = findRelatedCycles(targetCycle, [targetCycle, relatedCycle])
    expect(related[0].cycleId).toBe('pkg-a-pkg-c')
  })

  it('should handle non-repeating cycle array', () => {
    const targetCycle: CircularDependencyInfo = {
      ...baseCycleFields,
      cycle: ['pkg-a', 'pkg-b'],
    }
    const relatedCycle: CircularDependencyInfo = {
      ...baseCycleFields,
      cycle: ['pkg-a', 'pkg-c'],
    }
    const related = findRelatedCycles(targetCycle, [targetCycle, relatedCycle])
    expect(related).toHaveLength(1)
    expect(related[0].sharedPackages).toContain('pkg-a')
  })
})
