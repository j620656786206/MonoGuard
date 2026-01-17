'use client'

import type { CircularDependency, VersionConflict } from '@monoguard/types'
import React, { useState } from 'react'
import { VirtualizedList } from '@/components/ui/VirtualizedList'

export interface CircularDependencyVizProps {
  circularDependencies: CircularDependency[]
  versionConflicts: VersionConflict[]
}

export const CircularDependencyViz: React.FC<CircularDependencyVizProps> = ({
  circularDependencies,
  versionConflicts,
}) => {
  const [activeView, setActiveView] = useState<'circular' | 'conflicts'>('circular')

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold text-gray-900">Dependency Analysis</h3>
        <div className="flex rounded-lg bg-gray-100 p-1">
          <button
            onClick={() => setActiveView('circular')}
            className={`rounded-md px-3 py-1 text-sm font-medium transition-colors ${
              activeView === 'circular'
                ? 'bg-white text-blue-600 shadow-sm'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            Circular ({circularDependencies.length})
          </button>
          <button
            onClick={() => setActiveView('conflicts')}
            className={`rounded-md px-3 py-1 text-sm font-medium transition-colors ${
              activeView === 'conflicts'
                ? 'bg-white text-blue-600 shadow-sm'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            Conflicts ({versionConflicts.length})
          </button>
        </div>
      </div>

      {activeView === 'circular' ? (
        <CircularDependenciesPanel dependencies={circularDependencies} />
      ) : (
        <VersionConflictsPanel conflicts={versionConflicts} />
      )}
    </div>
  )
}

// Circular Dependencies Panel
const CircularDependenciesPanel: React.FC<{
  dependencies: CircularDependency[]
}> = ({ dependencies }) => {
  if (dependencies.length === 0) {
    return (
      <div className="rounded-lg border border-green-200 bg-green-50 p-8 text-center">
        <div className="mb-2 text-green-600">
          <svg className="mx-auto h-12 w-12" fill="currentColor" viewBox="0 0 20 20">
            <path
              fillRule="evenodd"
              d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
              clipRule="evenodd"
            />
          </svg>
        </div>
        <h3 className="mb-1 text-lg font-medium text-green-900">No Circular Dependencies</h3>
        <p className="text-green-700">Great! Your project is free of circular dependencies.</p>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {dependencies.map((dependency, index) => (
        <CircularDependencyCard key={index} dependency={dependency} />
      ))}
    </div>
  )
}

// Version Conflicts Panel
const VersionConflictsPanel: React.FC<{ conflicts: VersionConflict[] }> = ({ conflicts }) => {
  if (conflicts.length === 0) {
    return (
      <div className="rounded-lg border border-green-200 bg-green-50 p-8 text-center">
        <div className="mb-2 text-green-600">
          <svg className="mx-auto h-12 w-12" fill="currentColor" viewBox="0 0 20 20">
            <path
              fillRule="evenodd"
              d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
              clipRule="evenodd"
            />
          </svg>
        </div>
        <h3 className="mb-1 text-lg font-medium text-green-900">No Version Conflicts</h3>
        <p className="text-green-700">Excellent! All package versions are compatible.</p>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {conflicts.map((conflict, index) => (
        <VersionConflictCard key={index} conflict={conflict} />
      ))}
    </div>
  )
}

// Circular Dependency Card
const CircularDependencyCard: React.FC<{ dependency: CircularDependency }> = ({ dependency }) => {
  const [isExpanded, setIsExpanded] = useState(false)

  const severityColors = {
    low: 'bg-yellow-50 border-yellow-200 text-yellow-800',
    medium: 'bg-orange-50 border-orange-200 text-orange-800',
    high: 'bg-red-50 border-red-200 text-red-800',
    critical: 'bg-red-100 border-red-300 text-red-900',
  }

  return (
    <div className={`rounded-lg border p-4 ${severityColors[dependency.severity]}`}>
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <div className="mb-2 flex items-center space-x-2">
            <span
              className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${
                dependency.type === 'direct'
                  ? 'bg-red-100 text-red-800'
                  : 'bg-orange-100 text-orange-800'
              }`}
            >
              {dependency.type === 'direct' ? 'Direct' : 'Indirect'}
            </span>
            <span
              className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium capitalize ${
                severityColors[dependency.severity]
              }`}
            >
              {dependency.severity}
            </span>
          </div>

          <div className="mb-3">
            <h4 className="mb-1 font-medium">Dependency Cycle</h4>
            <div className="flex items-center space-x-2 font-mono text-sm">
              {dependency.cycle.map((dep, index) => (
                <React.Fragment key={index}>
                  <span className="rounded bg-white bg-opacity-50 px-2 py-1">{dep}</span>
                  {index < dependency.cycle.length - 1 && (
                    <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M9 5l7 7-7 7"
                      />
                    </svg>
                  )}
                  {index === dependency.cycle.length - 1 && (
                    <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M4 4l16 16m0-16L4 20"
                      />
                    </svg>
                  )}
                </React.Fragment>
              ))}
            </div>
          </div>

          <p className="mb-3 text-sm opacity-80">{dependency.impact}</p>

          {isExpanded && (
            <div className="mt-3 border-t border-current border-opacity-20 pt-3">
              <h5 className="mb-2 font-medium">Recommendations</h5>
              <ul className="space-y-1 text-sm opacity-80">
                <li>• Break the cycle by extracting common dependencies</li>
                <li>• Consider using dependency injection patterns</li>
                <li>• Refactor code to reduce coupling between modules</li>
                <li>• Review and simplify the module architecture</li>
              </ul>
            </div>
          )}
        </div>

        <button
          onClick={() => setIsExpanded(!isExpanded)}
          className="ml-4 text-sm opacity-60 transition-opacity hover:opacity-80"
        >
          {isExpanded ? 'Less' : 'More'}
        </button>
      </div>
    </div>
  )
}

// Version Conflict Card
const VersionConflictCard: React.FC<{ conflict: VersionConflict }> = ({ conflict }) => {
  const [isExpanded, setIsExpanded] = useState(false)

  const riskLevelColors = {
    low: 'bg-yellow-50 border-yellow-200 text-yellow-800',
    medium: 'bg-orange-50 border-orange-200 text-orange-800',
    high: 'bg-red-50 border-red-200 text-red-800',
    critical: 'bg-red-100 border-red-300 text-red-900',
  }

  return (
    <div className={`rounded-lg border p-4 ${riskLevelColors[conflict.riskLevel]}`}>
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <div className="mb-2 flex items-center space-x-2">
            <h4 className="font-medium">{conflict.packageName}</h4>
            <span
              className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium capitalize ${
                riskLevelColors[conflict.riskLevel]
              }`}
            >
              {conflict.riskLevel} Risk
            </span>
          </div>

          <div className="mb-3">
            <p className="mb-2 text-sm opacity-80">{conflict.impact}</p>

            <div className="space-y-2">
              {conflict.conflictingVersions.map((version, index) => (
                <div key={index} className="flex items-center space-x-2">
                  <span
                    className={`inline-flex items-center rounded px-2 py-1 font-mono text-xs ${
                      version.isBreaking
                        ? 'border border-red-200 bg-red-100 text-red-800'
                        : 'border border-current border-opacity-20 bg-white bg-opacity-50'
                    }`}
                  >
                    {version.version}
                  </span>
                  {version.isBreaking && (
                    <span className="text-xs font-medium text-red-600">Breaking</span>
                  )}
                  <span className="text-xs opacity-60">
                    Used by: {version.packages.slice(0, 2).join(', ')}
                    {version.packages.length > 2 && ` +${version.packages.length - 2} more`}
                  </span>
                </div>
              ))}
            </div>
          </div>

          {isExpanded && (
            <div className="mt-3 border-t border-current border-opacity-20 pt-3">
              <h5 className="mb-2 font-medium">Resolution Strategy</h5>
              <p className="mb-3 text-sm opacity-80">{conflict.resolution}</p>

              <h5 className="mb-2 font-medium">Affected Packages</h5>
              <div className="max-h-32 overflow-y-auto">
                {conflict.conflictingVersions.map((version, vIndex) => (
                  <div key={vIndex} className="mb-2">
                    <div className="mb-1 text-xs font-medium opacity-80">
                      Version {version.version}:
                    </div>
                    <div className="ml-2 space-y-1">
                      {version.packages.map((pkg, pIndex) => (
                        <div key={pIndex} className="font-mono text-xs opacity-60">
                          {pkg}
                        </div>
                      ))}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>

        <button
          onClick={() => setIsExpanded(!isExpanded)}
          className="ml-4 text-sm opacity-60 transition-opacity hover:opacity-80"
        >
          {isExpanded ? 'Less' : 'More'}
        </button>
      </div>
    </div>
  )
}
