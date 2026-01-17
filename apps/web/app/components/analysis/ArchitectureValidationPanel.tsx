'use client'

import type { ArchitectureValidationResults } from '@monoguard/types'
import type React from 'react'
import { useState } from 'react'
import { VirtualizedList } from '@/components/ui/VirtualizedList'

export interface ArchitectureValidationPanelProps {
  validation: ArchitectureValidationResults
}

export const ArchitectureValidationPanel: React.FC<ArchitectureValidationPanelProps> = ({
  validation,
}) => {
  const [activeView, setActiveView] = useState<'overview' | 'violations' | 'compliance'>('overview')

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold text-gray-900">Architecture Validation</h3>
        <div className="flex rounded-lg bg-gray-100 p-1">
          <button
            onClick={() => setActiveView('overview')}
            className={`rounded-md px-3 py-1 text-sm font-medium transition-colors ${
              activeView === 'overview'
                ? 'bg-white text-blue-600 shadow-sm'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            Overview
          </button>
          <button
            onClick={() => setActiveView('violations')}
            className={`rounded-md px-3 py-1 text-sm font-medium transition-colors ${
              activeView === 'violations'
                ? 'bg-white text-blue-600 shadow-sm'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            Violations ({validation.violations.length})
          </button>
          <button
            onClick={() => setActiveView('compliance')}
            className={`rounded-md px-3 py-1 text-sm font-medium transition-colors ${
              activeView === 'compliance'
                ? 'bg-white text-blue-600 shadow-sm'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            Layer Compliance
          </button>
        </div>
      </div>

      {/* Content */}
      {activeView === 'overview' && <OverviewSection validation={validation} />}
      {activeView === 'violations' && <ViolationsSection violations={validation.violations} />}
      {activeView === 'compliance' && <ComplianceSection compliance={validation.layerCompliance} />}
    </div>
  )
}

// Overview Section
const OverviewSection: React.FC<{
  validation: ArchitectureValidationResults
}> = ({ validation }) => {
  const { summary } = validation

  return (
    <div className="space-y-6">
      {/* Summary Cards */}
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
        <div
          className={`rounded-lg border p-4 ${
            summary.overallCompliance >= 90
              ? 'border-green-200 bg-green-50 text-green-700'
              : summary.overallCompliance >= 70
                ? 'border-yellow-200 bg-yellow-50 text-yellow-700'
                : 'border-red-200 bg-red-50 text-red-700'
          }`}
        >
          <div className="flex items-center justify-between">
            <div>
              <div className="text-2xl font-bold">{summary.overallCompliance.toFixed(1)}%</div>
              <div className="text-sm font-medium opacity-80">Overall Compliance</div>
            </div>
            <div className="opacity-60">
              <svg className="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
            </div>
          </div>
        </div>

        <div
          className={`rounded-lg border p-4 ${
            summary.totalViolations === 0
              ? 'border-green-200 bg-green-50 text-green-700'
              : 'border-red-200 bg-red-50 text-red-700'
          }`}
        >
          <div className="flex items-center justify-between">
            <div>
              <div className="text-2xl font-bold">{summary.totalViolations}</div>
              <div className="text-sm font-medium opacity-80">Total Violations</div>
            </div>
            <div className="opacity-60">
              <svg className="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.5 0L4.268 16.5c-.77.833.192 2.5 1.732 2.5z"
                />
              </svg>
            </div>
          </div>
        </div>

        <div
          className={`rounded-lg border p-4 ${
            summary.criticalViolations === 0
              ? 'border-green-200 bg-green-50 text-green-700'
              : 'border-red-200 bg-red-50 text-red-700'
          }`}
        >
          <div className="flex items-center justify-between">
            <div>
              <div className="text-2xl font-bold">{summary.criticalViolations}</div>
              <div className="text-sm font-medium opacity-80">Critical Issues</div>
            </div>
            <div className="opacity-60">
              <svg className="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728L5.636 5.636m12.728 12.728L18.364 5.636M5.636 18.364l12.728-12.728"
                />
              </svg>
            </div>
          </div>
        </div>

        <div className="rounded-lg border border-blue-200 bg-blue-50 p-4 text-blue-700">
          <div className="flex items-center justify-between">
            <div>
              <div className="text-2xl font-bold">{summary.layersAnalyzed}</div>
              <div className="text-sm font-medium opacity-80">Layers Analyzed</div>
            </div>
            <div className="opacity-60">
              <svg className="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
                />
              </svg>
            </div>
          </div>
        </div>
      </div>

      {/* Compliance Overview */}
      <div className="rounded-lg border border-gray-200 bg-white p-6">
        <h4 className="mb-4 text-lg font-medium text-gray-900">Layer Compliance Overview</h4>
        <div className="space-y-3">
          {validation.layerCompliance.map((layer, index) => (
            <div key={index} className="flex items-center space-x-4">
              <div className="w-32 truncate text-sm font-medium text-gray-900">
                {layer.layerName}
              </div>
              <div className="flex-1">
                <div className="flex items-center space-x-2">
                  <div className="h-2 flex-1 rounded-full bg-gray-200">
                    <div
                      className={`h-2 rounded-full transition-all duration-300 ${
                        layer.compliancePercentage >= 90
                          ? 'bg-green-500'
                          : layer.compliancePercentage >= 70
                            ? 'bg-yellow-500'
                            : 'bg-red-500'
                      }`}
                      style={{ width: `${layer.compliancePercentage}%` }}
                    />
                  </div>
                  <div className="w-12 text-right text-sm text-gray-600">
                    {layer.compliancePercentage.toFixed(0)}%
                  </div>
                </div>
              </div>
              <div className="w-24 text-right text-sm text-gray-500">
                {layer.compliantFiles}/{layer.totalFiles} files
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Recent Violations */}
      {validation.violations.length > 0 && (
        <div className="rounded-lg border border-gray-200 bg-white p-6">
          <h4 className="mb-4 text-lg font-medium text-gray-900">Recent Violations</h4>
          <div className="space-y-3">
            {validation.violations.slice(0, 5).map((violation, index) => (
              <ViolationSummaryCard key={index} violation={violation} />
            ))}
            {validation.violations.length > 5 && (
              <div className="border-t border-gray-200 pt-2 text-center text-sm text-gray-500">
                +{validation.violations.length - 5} more violations
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  )
}

// Violations Section
const ViolationsSection: React.FC<{
  violations: ArchitectureValidationResults['violations']
}> = ({ violations }) => {
  const [searchTerm, setSearchTerm] = useState('')
  const [selectedSeverity, setSelectedSeverity] = useState<string>('all')

  const filteredViolations = violations.filter((violation) => {
    const matchesSearch =
      violation.violatingFile.toLowerCase().includes(searchTerm.toLowerCase()) ||
      violation.ruleName.toLowerCase().includes(searchTerm.toLowerCase())
    const matchesSeverity = selectedSeverity === 'all' || violation.severity === selectedSeverity
    return matchesSearch && matchesSeverity
  })

  if (violations.length === 0) {
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
        <h3 className="mb-1 text-lg font-medium text-green-900">No Architecture Violations</h3>
        <p className="text-green-700">Your architecture follows all defined rules perfectly!</p>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {/* Filters */}
      <div className="flex flex-col gap-4 sm:flex-row">
        <div className="flex-1">
          <input
            type="text"
            placeholder="Search violations..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm placeholder-gray-500 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div className="sm:w-48">
          <select
            value={selectedSeverity}
            onChange={(e) => setSelectedSeverity(e.target.value)}
            className="block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="all">All Severities</option>
            <option value="critical">Critical</option>
            <option value="high">High</option>
            <option value="medium">Medium</option>
            <option value="low">Low</option>
          </select>
        </div>
      </div>

      {/* Results */}
      <div className="text-sm text-gray-600">
        Showing {filteredViolations.length} of {violations.length} violations
      </div>

      {/* Violations List */}
      <div className="space-y-3">
        {filteredViolations.map((violation, index) => (
          <ViolationCard key={index} violation={violation} />
        ))}
      </div>
    </div>
  )
}

// Compliance Section
const ComplianceSection: React.FC<{
  compliance: ArchitectureValidationResults['layerCompliance']
}> = ({ compliance }) => {
  return (
    <div className="space-y-6">
      {compliance.map((layer, index) => (
        <div key={index} className="rounded-lg border border-gray-200 bg-white p-6">
          <div className="mb-4 flex items-center justify-between">
            <h4 className="text-lg font-medium text-gray-900">{layer.layerName}</h4>
            <div
              className={`rounded-full px-3 py-1 text-sm font-medium ${
                layer.compliancePercentage >= 90
                  ? 'bg-green-100 text-green-800'
                  : layer.compliancePercentage >= 70
                    ? 'bg-yellow-100 text-yellow-800'
                    : 'bg-red-100 text-red-800'
              }`}
            >
              {layer.compliancePercentage.toFixed(1)}% Compliant
            </div>
          </div>

          <div className="mb-4 grid grid-cols-1 gap-4 sm:grid-cols-3">
            <div className="rounded-lg bg-gray-50 p-3 text-center">
              <div className="text-2xl font-bold text-gray-900">{layer.totalFiles}</div>
              <div className="text-sm text-gray-600">Total Files</div>
            </div>
            <div className="rounded-lg bg-green-50 p-3 text-center">
              <div className="text-2xl font-bold text-green-600">{layer.compliantFiles}</div>
              <div className="text-sm text-gray-600">Compliant Files</div>
            </div>
            <div className="rounded-lg bg-red-50 p-3 text-center">
              <div className="text-2xl font-bold text-red-600">{layer.violationCount}</div>
              <div className="text-sm text-gray-600">Violations</div>
            </div>
          </div>

          <div className="h-3 w-full rounded-full bg-gray-200">
            <div
              className={`h-3 rounded-full transition-all duration-300 ${
                layer.compliancePercentage >= 90
                  ? 'bg-green-500'
                  : layer.compliancePercentage >= 70
                    ? 'bg-yellow-500'
                    : 'bg-red-500'
              }`}
              style={{ width: `${layer.compliancePercentage}%` }}
            />
          </div>
        </div>
      ))}
    </div>
  )
}

// Violation Summary Card
const ViolationSummaryCard: React.FC<{
  violation: ArchitectureValidationResults['violations'][0]
}> = ({ violation }) => {
  const severityColors = {
    low: 'text-yellow-700 bg-yellow-100 border-yellow-200',
    medium: 'text-orange-700 bg-orange-100 border-orange-200',
    high: 'text-red-700 bg-red-100 border-red-200',
    critical: 'text-red-800 bg-red-200 border-red-300',
  }

  return (
    <div className="flex items-center space-x-3 rounded-lg bg-gray-50 p-3">
      <span
        className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium capitalize ${
          severityColors[violation.severity]
        }`}
      >
        {violation.severity}
      </span>
      <div className="min-w-0 flex-1">
        <div className="truncate text-sm font-medium text-gray-900">{violation.ruleName}</div>
        <div className="truncate text-sm text-gray-500">{violation.violatingFile}</div>
      </div>
    </div>
  )
}

// Full Violation Card
const ViolationCard: React.FC<{
  violation: ArchitectureValidationResults['violations'][0]
}> = ({ violation }) => {
  const [isExpanded, setIsExpanded] = useState(false)

  const severityColors = {
    low: 'bg-yellow-50 border-yellow-200 text-yellow-800',
    medium: 'bg-orange-50 border-orange-200 text-orange-800',
    high: 'bg-red-50 border-red-200 text-red-800',
    critical: 'bg-red-100 border-red-300 text-red-900',
  }

  return (
    <div className={`rounded-lg border p-4 ${severityColors[violation.severity]}`}>
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <div className="mb-2 flex items-center space-x-2">
            <span
              className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium capitalize ${
                severityColors[violation.severity]
              }`}
            >
              {violation.severity}
            </span>
            <h4 className="font-medium">{violation.ruleName}</h4>
          </div>

          <p className="mb-2 text-sm opacity-80">{violation.description}</p>

          <div className="space-y-1 text-sm">
            <div>
              <span className="font-medium">File:</span> {violation.violatingFile}
            </div>
            <div>
              <span className="font-medium">Import:</span> {violation.violatingImport}
            </div>
            <div>
              <span className="font-medium">Layer Issue:</span> Expected {violation.expectedLayer},
              found {violation.actualLayer}
            </div>
          </div>

          {isExpanded && (
            <div className="mt-3 border-t border-current border-opacity-20 pt-3">
              <h5 className="mb-2 font-medium">Suggestion</h5>
              <p className="text-sm opacity-80">{violation.suggestion}</p>
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
