'use client';

import React, { useState } from 'react';
import { ComprehensiveAnalysisResult } from '@monoguard/types';
import { CircularDependencyViz } from './CircularDependencyViz';
import { VersionConflictTable } from './VersionConflictTable';
import { ArchitectureValidationPanel } from './ArchitectureValidationPanel';
import { BundleImpactChart } from './BundleImpactChart';
import { HealthScoreDisplay } from './HealthScoreDisplay';
import { DuplicateDetectionPanel } from './DuplicateDetectionPanel';

export interface AnalysisResultsProps {
  analysis: ComprehensiveAnalysisResult;
  onNewAnalysis?: () => void;
}

type AnalysisTab =
  | 'overview'
  | 'dependencies'
  | 'architecture'
  | 'duplicates'
  | 'bundle'
  | 'health';

export const AnalysisResults: React.FC<AnalysisResultsProps> = ({
  analysis,
  onNewAnalysis,
}) => {
  const [activeTab, setActiveTab] = useState<AnalysisTab>('overview');

  const { results } = analysis;

  if (!results) {
    return (
      <div className="rounded-lg border border-gray-200 bg-white p-8 text-center shadow-sm">
        <div className="mb-2 text-gray-400">
          <svg
            className="mx-auto h-12 w-12"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
            />
          </svg>
        </div>
        <h3 className="mb-1 text-lg font-medium text-gray-900">
          Analysis in Progress
        </h3>
        <p className="mb-4 text-gray-600">
          {analysis.currentStep || 'Starting analysis...'}
        </p>
        <div className="h-2 w-full rounded-full bg-gray-200">
          <div
            className="h-2 rounded-full bg-blue-600 transition-all duration-300"
            style={{ width: `${analysis.progress}%` }}
          ></div>
        </div>
        <div className="mt-2 text-sm text-gray-500">
          {analysis.progress}% complete
        </div>
      </div>
    );
  }

  const tabs: {
    id: AnalysisTab;
    label: string;
    icon: React.ReactNode;
    count?: number;
  }[] = [
    {
      id: 'overview',
      label: 'Overview',
      icon: (
        <svg
          className="h-4 w-4"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
          />
        </svg>
      ),
    },
    {
      id: 'dependencies',
      label: 'Dependencies',
      icon: (
        <svg
          className="h-4 w-4"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"
          />
        </svg>
      ),
      count: results.dependencyAnalysis?.circularDependencies.length || 0,
    },
    {
      id: 'architecture',
      label: 'Architecture',
      icon: (
        <svg
          className="h-4 w-4"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
          />
        </svg>
      ),
      count: results.architectureValidation?.violations.length || 0,
    },
    {
      id: 'duplicates',
      label: 'Duplicates',
      icon: (
        <svg
          className="h-4 w-4"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
          />
        </svg>
      ),
      count: results.duplicateDetection?.totalDuplicates || 0,
    },
    {
      id: 'bundle',
      label: 'Bundle Impact',
      icon: (
        <svg
          className="h-4 w-4"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4h16v12a1 1 0 01-1 1H5a1 1 0 01-1-1V4z"
          />
        </svg>
      ),
    },
    {
      id: 'health',
      label: 'Health Score',
      icon: (
        <svg
          className="h-4 w-4"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"
          />
        </svg>
      ),
    },
  ];

  const renderTabContent = () => {
    switch (activeTab) {
      case 'overview':
        return <OverviewPanel analysis={analysis} />;
      case 'dependencies':
        return results.dependencyAnalysis ? (
          <CircularDependencyViz
            circularDependencies={
              results.dependencyAnalysis.circularDependencies
            }
            versionConflicts={results.dependencyAnalysis.versionConflicts}
          />
        ) : (
          <EmptyState message="No dependency analysis results available" />
        );
      case 'architecture':
        return results.architectureValidation ? (
          <ArchitectureValidationPanel
            validation={results.architectureValidation}
          />
        ) : (
          <EmptyState message="No architecture validation results available" />
        );
      case 'duplicates':
        return results.duplicateDetection ? (
          <DuplicateDetectionPanel duplicates={results.duplicateDetection} />
        ) : (
          <EmptyState message="No duplicate detection results available" />
        );
      case 'bundle':
        return results.bundleImpact ? (
          <BundleImpactChart bundleImpact={results.bundleImpact} />
        ) : (
          <EmptyState message="No bundle impact analysis available" />
        );
      case 'health': {
        // Create a properly structured HealthScore object from the numeric value
        const healthScoreObj =
          typeof results.healthScore === 'number'
            ? {
                overall: results.healthScore,
                dependencies: results.summary?.duplicateCount === 0 ? 100 : 80,
                architecture: 90,
                maintainability:
                  results.summary?.circularCount === 0 ? 100 : 70,
                security: 95,
                performance: results.bundleImpact?.potentialSavings ? 85 : 100,
                lastUpdated: analysis.completedAt || new Date().toISOString(),
                trend: 'stable' as const,
                factors: [
                  {
                    name: 'Dependencies',
                    score: results.summary?.duplicateCount === 0 ? 100 : 80,
                    weight: 0.3,
                    description: `${results.summary?.duplicateCount || 0} duplicate dependencies found`,
                    recommendations:
                      results.summary?.duplicateCount > 0
                        ? ['Remove duplicate dependencies']
                        : [],
                  },
                  {
                    name: 'Circular Dependencies',
                    score: results.summary?.circularCount === 0 ? 100 : 70,
                    weight: 0.2,
                    description: `${results.summary?.circularCount || 0} circular dependencies detected`,
                    recommendations:
                      results.summary?.circularCount > 0
                        ? ['Refactor circular dependencies']
                        : [],
                  },
                  {
                    name: 'Version Conflicts',
                    score: results.summary?.conflictCount === 0 ? 100 : 75,
                    weight: 0.2,
                    description: `${results.summary?.conflictCount || 0} version conflicts found`,
                    recommendations:
                      results.summary?.conflictCount > 0
                        ? ['Resolve version conflicts']
                        : [],
                  },
                  {
                    name: 'Bundle Size',
                    score: results.bundleImpact?.potentialSavings ? 85 : 100,
                    weight: 0.3,
                    description: `${results.bundleImpact?.potentialSavings || '0 KB'} potential savings available`,
                    recommendations: results.bundleImpact?.potentialSavings
                      ? ['Remove unused dependencies to reduce bundle size']
                      : [],
                  },
                ],
              }
            : results.healthScore;

        return healthScoreObj ? (
          <HealthScoreDisplay healthScore={healthScoreObj} />
        ) : (
          <EmptyState message="No health score available" />
        );
      }
      default:
        return <EmptyState message="Unknown tab selected" />;
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="rounded-lg border border-gray-200 bg-white p-6 shadow-sm">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-2xl font-bold text-gray-900">
              Analysis Results
            </h2>
            <p className="mt-1 text-gray-600">
              Completed{' '}
              {analysis.completedAt
                ? new Date(analysis.completedAt).toLocaleString()
                : 'recently'}
            </p>
          </div>
          {onNewAnalysis && (
            <button
              onClick={onNewAnalysis}
              className="rounded-lg bg-blue-600 px-4 py-2 text-white transition-colors hover:bg-blue-700"
            >
              New Analysis
            </button>
          )}
        </div>
      </div>

      {/* Tabs */}
      <div className="rounded-lg border border-gray-200 bg-white shadow-sm">
        <div className="border-b border-gray-200">
          <nav className="flex space-x-8 px-6" aria-label="Tabs">
            {tabs.map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`flex items-center space-x-2 border-b-2 px-1 py-4 text-sm font-medium transition-colors ${
                  activeTab === tab.id
                    ? 'border-blue-500 text-blue-600'
                    : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700'
                }`}
              >
                {tab.icon}
                <span>{tab.label}</span>
                {tab.count !== undefined && tab.count > 0 && (
                  <span
                    className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium ${
                      activeTab === tab.id
                        ? 'bg-blue-100 text-blue-800'
                        : 'bg-gray-100 text-gray-800'
                    }`}
                  >
                    {tab.count}
                  </span>
                )}
              </button>
            ))}
          </nav>
        </div>

        {/* Tab Content */}
        <div className="p-6">{renderTabContent()}</div>
      </div>
    </div>
  );
};

// Overview Panel Component
const OverviewPanel: React.FC<{ analysis: ComprehensiveAnalysisResult }> = ({
  analysis,
}) => {
  const { results } = analysis;

  const stats = [
    {
      label: 'Health Score',
      value: results?.healthScore || results?.summary?.healthScore || 0,
      suffix: '/100',
      color:
        (results?.healthScore || results?.summary?.healthScore || 0) >= 80
          ? 'green'
          : (results?.healthScore || results?.summary?.healthScore || 0) >= 60
            ? 'yellow'
            : 'red',
      icon: (
        <svg
          className="h-6 w-6"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"
          />
        </svg>
      ),
    },
    {
      label: 'Circular Dependencies',
      value:
        results?.summary?.circularCount ||
        results?.circularDependencies?.length ||
        0,
      color:
        (results?.summary?.circularCount ||
          results?.circularDependencies?.length ||
          0) > 0
          ? 'red'
          : 'green',
      icon: (
        <svg
          className="h-6 w-6"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"
          />
        </svg>
      ),
    },
    {
      label: 'Version Conflicts',
      value:
        results?.summary?.conflictCount ||
        results?.versionConflicts?.length ||
        0,
      color:
        (results?.summary?.conflictCount ||
          results?.versionConflicts?.length ||
          0) > 0
          ? 'yellow'
          : 'green',
      icon: (
        <svg
          className="h-6 w-6"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.5 0L4.268 16.5c-.77.833.192 2.5 1.732 2.5z"
          />
        </svg>
      ),
    },
    {
      label: 'Duplicate Packages',
      value:
        results?.summary?.duplicateCount || results?.duplicates?.length || 0,
      color:
        (results?.summary?.duplicateCount || results?.duplicates?.length || 0) >
        0
          ? 'yellow'
          : 'green',
      icon: (
        <svg
          className="h-6 w-6"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
          />
        </svg>
      ),
    },
    {
      label: 'Architecture Violations',
      value: results?.architectureValidation?.violations.length || 0,
      color:
        (results?.architectureValidation?.violations.length || 0) > 0
          ? 'red'
          : 'green',
      icon: (
        <svg
          className="h-6 w-6"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
          />
        </svg>
      ),
    },
    {
      label: 'Potential Savings',
      value: results?.bundleImpact?.potentialSavings || '0 KB',
      color: 'blue',
      icon: (
        <svg
          className="h-6 w-6"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4h16v12a1 1 0 01-1 1H5a1 1 0 01-1-1V4z"
          />
        </svg>
      ),
    },
  ];

  const colorClasses = {
    green: 'bg-green-50 border-green-200 text-green-700',
    yellow: 'bg-yellow-50 border-yellow-200 text-yellow-700',
    red: 'bg-red-50 border-red-200 text-red-700',
    blue: 'bg-blue-50 border-blue-200 text-blue-700',
  };

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
        {stats.map((stat) => (
          <div
            key={stat.label}
            className={`rounded-lg border p-4 ${colorClasses[stat.color as keyof typeof colorClasses]}`}
          >
            <div className="flex items-center justify-between">
              <div>
                <div className="text-2xl font-bold">
                  {stat.value}
                  {stat.suffix || ''}
                </div>
                <div className="text-sm font-medium opacity-80">
                  {stat.label}
                </div>
              </div>
              <div className="opacity-60">{stat.icon}</div>
            </div>
          </div>
        ))}
      </div>

      {analysis.warnings && analysis.warnings.length > 0 && (
        <div className="rounded-lg border border-yellow-200 bg-yellow-50 p-4">
          <div className="flex">
            <div className="flex-shrink-0">
              <svg
                className="h-5 w-5 text-yellow-400"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path
                  fillRule="evenodd"
                  d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
                  clipRule="evenodd"
                />
              </svg>
            </div>
            <div className="ml-3">
              <h3 className="text-sm font-medium text-yellow-800">Warnings</h3>
              <div className="mt-2 text-sm text-yellow-700">
                <ul className="list-inside list-disc space-y-1">
                  {analysis.warnings.map((warning, index) => (
                    <li key={index}>{warning}</li>
                  ))}
                </ul>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

// Empty State Component
const EmptyState: React.FC<{ message: string }> = ({ message }) => (
  <div className="py-12 text-center">
    <svg
      className="mx-auto h-12 w-12 text-gray-400"
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
      />
    </svg>
    <h3 className="mt-2 text-sm font-medium text-gray-900">No Data</h3>
    <p className="mt-1 text-sm text-gray-500">{message}</p>
  </div>
);
