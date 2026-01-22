import type { CircularDependency, VersionConflict } from '@monoguard/types'
import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { CircularDependencyViz } from '@/components/analysis/CircularDependencyViz'

// Mock data factories
const createCircularDependency = (
  overrides: Partial<CircularDependency> = {}
): CircularDependency => ({
  cycle: ['package-a', 'package-b', 'package-c'],
  type: 'direct',
  severity: 'medium',
  impact: 'May cause build issues and increased bundle size',
  ...overrides,
})

const createVersionConflict = (overrides: Partial<VersionConflict> = {}): VersionConflict => ({
  packageName: 'lodash',
  conflictingVersions: [
    {
      version: '4.17.21',
      packages: ['package-a', 'package-b'],
      isBreaking: false,
    },
    {
      version: '3.10.1',
      packages: ['package-c'],
      isBreaking: true,
    },
  ],
  riskLevel: 'medium',
  resolution: 'Upgrade all packages to use lodash 4.17.21',
  impact: 'Different lodash versions may cause inconsistent behavior',
  ...overrides,
})

describe('CircularDependencyViz', () => {
  describe('Header and Tab Switching', () => {
    it('[P1] should display dependency analysis header', () => {
      // GIVEN: Component with empty data
      render(<CircularDependencyViz circularDependencies={[]} versionConflicts={[]} />)

      // THEN: Header should be visible
      expect(screen.getByText('Dependency Analysis')).toBeInTheDocument()
    })

    it('[P1] should display tab buttons with counts', () => {
      // GIVEN: Component with data
      const circularDeps = [createCircularDependency()]
      const conflicts = [createVersionConflict(), createVersionConflict()]

      render(
        <CircularDependencyViz circularDependencies={circularDeps} versionConflicts={conflicts} />
      )

      // THEN: Tab buttons should show correct counts
      expect(screen.getByText('Circular (1)')).toBeInTheDocument()
      expect(screen.getByText('Conflicts (2)')).toBeInTheDocument()
    })

    it('[P1] should switch between circular and conflicts views', () => {
      // GIVEN: Component with data
      const circularDeps = [createCircularDependency()]
      const conflicts = [createVersionConflict()]

      render(
        <CircularDependencyViz circularDependencies={circularDeps} versionConflicts={conflicts} />
      )

      // WHEN: Default view (circular)
      // THEN: Circular dependency content should be visible
      expect(screen.getByText('Dependency Cycle')).toBeInTheDocument()

      // WHEN: Click conflicts tab
      fireEvent.click(screen.getByText('Conflicts (1)'))

      // THEN: Version conflict content should be visible
      expect(screen.getByText('lodash')).toBeInTheDocument()
    })
  })

  describe('Circular Dependencies Panel', () => {
    it('[P1] should display empty state when no circular dependencies', () => {
      // GIVEN: No circular dependencies
      render(<CircularDependencyViz circularDependencies={[]} versionConflicts={[]} />)

      // THEN: Empty state message should be visible
      expect(screen.getByText('No Circular Dependencies')).toBeInTheDocument()
      expect(
        screen.getByText('Great! Your project is free of circular dependencies.')
      ).toBeInTheDocument()
    })

    it('[P1] should display circular dependency card with cycle path', () => {
      // GIVEN: Circular dependency with cycle
      const circularDeps = [createCircularDependency()]

      render(<CircularDependencyViz circularDependencies={circularDeps} versionConflicts={[]} />)

      // THEN: Cycle packages should be displayed
      expect(screen.getByText('package-a')).toBeInTheDocument()
      expect(screen.getByText('package-b')).toBeInTheDocument()
      expect(screen.getByText('package-c')).toBeInTheDocument()
    })

    it('[P1] should display severity and type badges', () => {
      // GIVEN: Circular dependency with direct type and medium severity
      const circularDeps = [
        createCircularDependency({
          type: 'direct',
          severity: 'high',
        }),
      ]

      render(<CircularDependencyViz circularDependencies={circularDeps} versionConflicts={[]} />)

      // THEN: Type and severity badges should be visible
      expect(screen.getByText('Direct')).toBeInTheDocument()
      expect(screen.getByText('high')).toBeInTheDocument()
    })

    it('[P1] should display impact description', () => {
      // GIVEN: Circular dependency with impact
      const circularDeps = [
        createCircularDependency({
          impact: 'Critical build failure risk',
        }),
      ]

      render(<CircularDependencyViz circularDependencies={circularDeps} versionConflicts={[]} />)

      // THEN: Impact should be displayed
      expect(screen.getByText('Critical build failure risk')).toBeInTheDocument()
    })

    it('[P1] should expand/collapse recommendations', () => {
      // GIVEN: Circular dependency
      const circularDeps = [createCircularDependency()]

      render(<CircularDependencyViz circularDependencies={circularDeps} versionConflicts={[]} />)

      // WHEN: Click More button
      fireEvent.click(screen.getByText('More'))

      // THEN: Recommendations should be visible
      expect(screen.getByText('Recommendations')).toBeInTheDocument()
      expect(
        screen.getByText(/Break the cycle by extracting common dependencies/)
      ).toBeInTheDocument()

      // WHEN: Click Less button
      fireEvent.click(screen.getByText('Less'))

      // THEN: Recommendations should be hidden
      expect(screen.queryByText('Recommendations')).not.toBeInTheDocument()
    })

    it('[P2] should handle indirect dependency type', () => {
      // GIVEN: Indirect circular dependency
      const circularDeps = [
        createCircularDependency({
          type: 'indirect',
        }),
      ]

      render(<CircularDependencyViz circularDependencies={circularDeps} versionConflicts={[]} />)

      // THEN: Indirect badge should be visible
      expect(screen.getByText('Indirect')).toBeInTheDocument()
    })

    it('[P2] should apply correct severity colors', () => {
      // GIVEN: Critical severity dependency
      const circularDeps = [
        createCircularDependency({
          severity: 'critical',
        }),
      ]

      render(<CircularDependencyViz circularDependencies={circularDeps} versionConflicts={[]} />)

      // THEN: Critical severity should be displayed
      expect(screen.getByText('critical')).toBeInTheDocument()
    })
  })

  describe('Version Conflicts Panel', () => {
    it('[P1] should display empty state when no version conflicts', () => {
      // GIVEN: No version conflicts
      render(<CircularDependencyViz circularDependencies={[]} versionConflicts={[]} />)

      // WHEN: Switch to conflicts view
      fireEvent.click(screen.getByText('Conflicts (0)'))

      // THEN: Empty state message should be visible
      expect(screen.getByText('No Version Conflicts')).toBeInTheDocument()
      expect(
        screen.getByText('Excellent! All package versions are compatible.')
      ).toBeInTheDocument()
    })

    it('[P1] should display version conflict card', () => {
      // GIVEN: Version conflict
      const conflicts = [createVersionConflict()]

      render(<CircularDependencyViz circularDependencies={[]} versionConflicts={conflicts} />)

      // WHEN: Switch to conflicts view
      fireEvent.click(screen.getByText('Conflicts (1)'))

      // THEN: Package name and risk level should be displayed
      expect(screen.getByText('lodash')).toBeInTheDocument()
      expect(screen.getByText('medium Risk')).toBeInTheDocument()
    })

    it('[P1] should display conflicting versions', () => {
      // GIVEN: Version conflict with multiple versions
      const conflicts = [createVersionConflict()]

      render(<CircularDependencyViz circularDependencies={[]} versionConflicts={conflicts} />)

      // WHEN: Switch to conflicts view
      fireEvent.click(screen.getByText('Conflicts (1)'))

      // THEN: Versions should be displayed
      expect(screen.getByText('4.17.21')).toBeInTheDocument()
      expect(screen.getByText('3.10.1')).toBeInTheDocument()
    })

    it('[P1] should mark breaking changes', () => {
      // GIVEN: Version conflict with breaking change
      const conflicts = [createVersionConflict()]

      render(<CircularDependencyViz circularDependencies={[]} versionConflicts={conflicts} />)

      // WHEN: Switch to conflicts view
      fireEvent.click(screen.getByText('Conflicts (1)'))

      // THEN: Breaking badge should be visible
      expect(screen.getByText('Breaking')).toBeInTheDocument()
    })

    it('[P1] should display impact description', () => {
      // GIVEN: Version conflict with impact
      const conflicts = [
        createVersionConflict({
          impact: 'API incompatibility between versions',
        }),
      ]

      render(<CircularDependencyViz circularDependencies={[]} versionConflicts={conflicts} />)

      // WHEN: Switch to conflicts view
      fireEvent.click(screen.getByText('Conflicts (1)'))

      // THEN: Impact should be displayed
      expect(screen.getByText('API incompatibility between versions')).toBeInTheDocument()
    })

    it('[P1] should expand/collapse resolution details', () => {
      // GIVEN: Version conflict
      const conflicts = [createVersionConflict()]

      render(<CircularDependencyViz circularDependencies={[]} versionConflicts={conflicts} />)

      // WHEN: Switch to conflicts view and click More
      fireEvent.click(screen.getByText('Conflicts (1)'))
      fireEvent.click(screen.getByText('More'))

      // THEN: Resolution strategy should be visible
      expect(screen.getByText('Resolution Strategy')).toBeInTheDocument()
      expect(screen.getByText('Upgrade all packages to use lodash 4.17.21')).toBeInTheDocument()

      // WHEN: Click Less
      fireEvent.click(screen.getByText('Less'))

      // THEN: Resolution should be hidden
      expect(screen.queryByText('Resolution Strategy')).not.toBeInTheDocument()
    })

    it('[P2] should show affected packages list', () => {
      // GIVEN: Version conflict with affected packages
      const conflicts = [createVersionConflict()]

      render(<CircularDependencyViz circularDependencies={[]} versionConflicts={conflicts} />)

      // WHEN: Switch to conflicts view
      fireEvent.click(screen.getByText('Conflicts (1)'))

      // THEN: Some affected packages should be shown
      expect(screen.getByText(/package-a/)).toBeInTheDocument()
    })
  })

  describe('Multiple Items', () => {
    it('[P1] should display multiple circular dependencies', () => {
      // GIVEN: Multiple circular dependencies
      const circularDeps = [
        createCircularDependency({ cycle: ['a', 'b'] }),
        createCircularDependency({ cycle: ['x', 'y', 'z'] }),
      ]

      render(<CircularDependencyViz circularDependencies={circularDeps} versionConflicts={[]} />)

      // THEN: Both cycles should be displayed
      expect(screen.getByText('a')).toBeInTheDocument()
      expect(screen.getByText('x')).toBeInTheDocument()
    })

    it('[P1] should display multiple version conflicts', () => {
      // GIVEN: Multiple version conflicts
      const conflicts = [
        createVersionConflict({ packageName: 'react' }),
        createVersionConflict({ packageName: 'typescript' }),
      ]

      render(<CircularDependencyViz circularDependencies={[]} versionConflicts={conflicts} />)

      // WHEN: Switch to conflicts view
      fireEvent.click(screen.getByText('Conflicts (2)'))

      // THEN: Both packages should be displayed
      expect(screen.getByText('react')).toBeInTheDocument()
      expect(screen.getByText('typescript')).toBeInTheDocument()
    })
  })

  describe('Risk Levels', () => {
    it('[P2] should apply correct styling for low risk', () => {
      // GIVEN: Low risk conflict
      const conflicts = [createVersionConflict({ riskLevel: 'low' })]

      render(<CircularDependencyViz circularDependencies={[]} versionConflicts={conflicts} />)

      // WHEN: Switch to conflicts view
      fireEvent.click(screen.getByText('Conflicts (1)'))

      // THEN: Low risk badge should be visible
      expect(screen.getByText('low Risk')).toBeInTheDocument()
    })

    it('[P2] should apply correct styling for high risk', () => {
      // GIVEN: High risk conflict
      const conflicts = [createVersionConflict({ riskLevel: 'high' })]

      render(<CircularDependencyViz circularDependencies={[]} versionConflicts={conflicts} />)

      // WHEN: Switch to conflicts view
      fireEvent.click(screen.getByText('Conflicts (1)'))

      // THEN: High risk badge should be visible
      expect(screen.getByText('high Risk')).toBeInTheDocument()
    })

    it('[P2] should apply correct styling for critical risk', () => {
      // GIVEN: Critical risk conflict
      const conflicts = [createVersionConflict({ riskLevel: 'critical' })]

      render(<CircularDependencyViz circularDependencies={[]} versionConflicts={conflicts} />)

      // WHEN: Switch to conflicts view
      fireEvent.click(screen.getByText('Conflicts (1)'))

      // THEN: Critical risk badge should be visible
      expect(screen.getByText('critical Risk')).toBeInTheDocument()
    })
  })
})
