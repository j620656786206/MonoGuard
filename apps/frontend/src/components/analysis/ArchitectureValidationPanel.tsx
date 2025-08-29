'use client';

import React, { useState } from 'react';
import { ArchitectureValidationResults } from '@monoguard/shared-types';
import { VirtualizedList } from '@/components/ui/VirtualizedList';

export interface ArchitectureValidationPanelProps {
  validation: ArchitectureValidationResults;
}

export const ArchitectureValidationPanel: React.FC<ArchitectureValidationPanelProps> = ({
  validation,
}) => {
  const [activeView, setActiveView] = useState<'overview' | 'violations' | 'compliance'>('overview');

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold text-gray-900">
          Architecture Validation
        </h3>
        <div className="flex bg-gray-100 rounded-lg p-1">
          <button
            onClick={() => setActiveView('overview')}
            className={`px-3 py-1 text-sm font-medium rounded-md transition-colors ${
              activeView === 'overview'
                ? 'bg-white text-blue-600 shadow-sm'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            Overview
          </button>
          <button
            onClick={() => setActiveView('violations')}
            className={`px-3 py-1 text-sm font-medium rounded-md transition-colors ${
              activeView === 'violations'
                ? 'bg-white text-blue-600 shadow-sm'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            Violations ({validation.violations.length})
          </button>
          <button
            onClick={() => setActiveView('compliance')}
            className={`px-3 py-1 text-sm font-medium rounded-md transition-colors ${
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
  );
};

// Overview Section
const OverviewSection: React.FC<{ validation: ArchitectureValidationResults }> = ({ 
  validation 
}) => {
  const { summary } = validation;

  return (
    <div className="space-y-6">
      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <div className={`rounded-lg border p-4 ${
          summary.overallCompliance >= 90 
            ? 'bg-green-50 border-green-200 text-green-700'
            : summary.overallCompliance >= 70
            ? 'bg-yellow-50 border-yellow-200 text-yellow-700'
            : 'bg-red-50 border-red-200 text-red-700'
        }`}>
          <div className="flex items-center justify-between">
            <div>
              <div className="text-2xl font-bold">
                {summary.overallCompliance.toFixed(1)}%
              </div>
              <div className="text-sm font-medium opacity-80">Overall Compliance</div>
            </div>
            <div className="opacity-60">
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
          </div>
        </div>

        <div className={`rounded-lg border p-4 ${
          summary.totalViolations === 0
            ? 'bg-green-50 border-green-200 text-green-700'
            : 'bg-red-50 border-red-200 text-red-700'
        }`}>
          <div className="flex items-center justify-between">
            <div>
              <div className="text-2xl font-bold">{summary.totalViolations}</div>
              <div className="text-sm font-medium opacity-80">Total Violations</div>
            </div>
            <div className="opacity-60">
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.5 0L4.268 16.5c-.77.833.192 2.5 1.732 2.5z" />
              </svg>
            </div>
          </div>
        </div>

        <div className={`rounded-lg border p-4 ${
          summary.criticalViolations === 0
            ? 'bg-green-50 border-green-200 text-green-700'
            : 'bg-red-50 border-red-200 text-red-700'
        }`}>
          <div className="flex items-center justify-between">
            <div>
              <div className="text-2xl font-bold">{summary.criticalViolations}</div>
              <div className="text-sm font-medium opacity-80">Critical Issues</div>
            </div>
            <div className="opacity-60">
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728L5.636 5.636m12.728 12.728L18.364 5.636M5.636 18.364l12.728-12.728" />
              </svg>
            </div>
          </div>
        </div>

        <div className="bg-blue-50 border-blue-200 text-blue-700 rounded-lg border p-4">
          <div className="flex items-center justify-between">
            <div>
              <div className="text-2xl font-bold">{summary.layersAnalyzed}</div>
              <div className="text-sm font-medium opacity-80">Layers Analyzed</div>
            </div>
            <div className="opacity-60">
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
              </svg>
            </div>
          </div>
        </div>
      </div>

      {/* Compliance Overview */}
      <div className="bg-white border border-gray-200 rounded-lg p-6">
        <h4 className="text-lg font-medium text-gray-900 mb-4">Layer Compliance Overview</h4>
        <div className="space-y-3">
          {validation.layerCompliance.map((layer, index) => (
            <div key={index} className="flex items-center space-x-4">
              <div className="w-32 text-sm font-medium text-gray-900 truncate">
                {layer.layerName}
              </div>
              <div className="flex-1">
                <div className="flex items-center space-x-2">
                  <div className="flex-1 bg-gray-200 rounded-full h-2">
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
                  <div className="w-12 text-sm text-gray-600 text-right">
                    {layer.compliancePercentage.toFixed(0)}%
                  </div>
                </div>
              </div>
              <div className="w-24 text-sm text-gray-500 text-right">
                {layer.compliantFiles}/{layer.totalFiles} files
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Recent Violations */}
      {validation.violations.length > 0 && (
        <div className="bg-white border border-gray-200 rounded-lg p-6">
          <h4 className="text-lg font-medium text-gray-900 mb-4">Recent Violations</h4>
          <div className="space-y-3">
            {validation.violations.slice(0, 5).map((violation, index) => (
              <ViolationSummaryCard key={index} violation={violation} />
            ))}
            {validation.violations.length > 5 && (
              <div className="text-sm text-gray-500 text-center pt-2 border-t border-gray-200">
                +{validation.violations.length - 5} more violations
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

// Violations Section
const ViolationsSection: React.FC<{ violations: ArchitectureValidationResults['violations'] }> = ({ 
  violations 
}) => {
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedSeverity, setSelectedSeverity] = useState<string>('all');

  const filteredViolations = violations.filter((violation) => {
    const matchesSearch = 
      violation.violatingFile.toLowerCase().includes(searchTerm.toLowerCase()) ||
      violation.ruleName.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesSeverity = selectedSeverity === 'all' || violation.severity === selectedSeverity;
    return matchesSearch && matchesSeverity;
  });

  if (violations.length === 0) {
    return (
      <div className="bg-green-50 border border-green-200 rounded-lg p-8 text-center">
        <div className="text-green-600 mb-2">
          <svg className="mx-auto h-12 w-12" fill="currentColor" viewBox="0 0 20 20">
            <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
          </svg>
        </div>
        <h3 className="text-lg font-medium text-green-900 mb-1">No Architecture Violations</h3>
        <p className="text-green-700">Your architecture follows all defined rules perfectly!</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Filters */}
      <div className="flex flex-col sm:flex-row gap-4">
        <div className="flex-1">
          <input
            type="text"
            placeholder="Search violations..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="block w-full px-3 py-2 border border-gray-300 rounded-lg text-sm placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          />
        </div>
        <div className="sm:w-48">
          <select
            value={selectedSeverity}
            onChange={(e) => setSelectedSeverity(e.target.value)}
            className="block w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
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
  );
};

// Compliance Section
const ComplianceSection: React.FC<{ compliance: ArchitectureValidationResults['layerCompliance'] }> = ({ 
  compliance 
}) => {
  return (
    <div className="space-y-6">
      {compliance.map((layer, index) => (
        <div key={index} className="bg-white border border-gray-200 rounded-lg p-6">
          <div className="flex items-center justify-between mb-4">
            <h4 className="text-lg font-medium text-gray-900">{layer.layerName}</h4>
            <div className={`px-3 py-1 rounded-full text-sm font-medium ${
              layer.compliancePercentage >= 90
                ? 'bg-green-100 text-green-800'
                : layer.compliancePercentage >= 70
                ? 'bg-yellow-100 text-yellow-800'
                : 'bg-red-100 text-red-800'
            }`}>
              {layer.compliancePercentage.toFixed(1)}% Compliant
            </div>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-4">
            <div className="text-center p-3 bg-gray-50 rounded-lg">
              <div className="text-2xl font-bold text-gray-900">{layer.totalFiles}</div>
              <div className="text-sm text-gray-600">Total Files</div>
            </div>
            <div className="text-center p-3 bg-green-50 rounded-lg">
              <div className="text-2xl font-bold text-green-600">{layer.compliantFiles}</div>
              <div className="text-sm text-gray-600">Compliant Files</div>
            </div>
            <div className="text-center p-3 bg-red-50 rounded-lg">
              <div className="text-2xl font-bold text-red-600">{layer.violationCount}</div>
              <div className="text-sm text-gray-600">Violations</div>
            </div>
          </div>

          <div className="w-full bg-gray-200 rounded-full h-3">
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
  );
};

// Violation Summary Card
const ViolationSummaryCard: React.FC<{ violation: ArchitectureValidationResults['violations'][0] }> = ({ 
  violation 
}) => {
  const severityColors = {
    low: 'text-yellow-700 bg-yellow-100 border-yellow-200',
    medium: 'text-orange-700 bg-orange-100 border-orange-200',
    high: 'text-red-700 bg-red-100 border-red-200',
    critical: 'text-red-800 bg-red-200 border-red-300',
  };

  return (
    <div className="flex items-center space-x-3 p-3 bg-gray-50 rounded-lg">
      <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium capitalize ${
        severityColors[violation.severity]
      }`}>
        {violation.severity}
      </span>
      <div className="flex-1 min-w-0">
        <div className="text-sm font-medium text-gray-900 truncate">
          {violation.ruleName}
        </div>
        <div className="text-sm text-gray-500 truncate">
          {violation.violatingFile}
        </div>
      </div>
    </div>
  );
};

// Full Violation Card
const ViolationCard: React.FC<{ violation: ArchitectureValidationResults['violations'][0] }> = ({ 
  violation 
}) => {
  const [isExpanded, setIsExpanded] = useState(false);

  const severityColors = {
    low: 'bg-yellow-50 border-yellow-200 text-yellow-800',
    medium: 'bg-orange-50 border-orange-200 text-orange-800',
    high: 'bg-red-50 border-red-200 text-red-800',
    critical: 'bg-red-100 border-red-300 text-red-900',
  };

  return (
    <div className={`rounded-lg border p-4 ${severityColors[violation.severity]}`}>
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <div className="flex items-center space-x-2 mb-2">
            <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium capitalize ${
              severityColors[violation.severity]
            }`}>
              {violation.severity}
            </span>
            <h4 className="font-medium">{violation.ruleName}</h4>
          </div>

          <p className="text-sm opacity-80 mb-2">{violation.description}</p>

          <div className="text-sm space-y-1">
            <div><span className="font-medium">File:</span> {violation.violatingFile}</div>
            <div><span className="font-medium">Import:</span> {violation.violatingImport}</div>
            <div>
              <span className="font-medium">Layer Issue:</span> Expected {violation.expectedLayer}, found {violation.actualLayer}
            </div>
          </div>

          {isExpanded && (
            <div className="mt-3 pt-3 border-t border-current border-opacity-20">
              <h5 className="font-medium mb-2">Suggestion</h5>
              <p className="text-sm opacity-80">{violation.suggestion}</p>
            </div>
          )}
        </div>

        <button
          onClick={() => setIsExpanded(!isExpanded)}
          className="ml-4 text-sm opacity-60 hover:opacity-80 transition-opacity"
        >
          {isExpanded ? 'Less' : 'More'}
        </button>
      </div>
    </div>
  );
};