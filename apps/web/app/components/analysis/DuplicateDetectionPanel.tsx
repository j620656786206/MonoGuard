'use client';

import React, { useState } from 'react';
import {
  DuplicateDetectionResults,
  DuplicateGroup,
  DuplicateRecommendation,
} from '@monoguard/types';

export interface DuplicateDetectionPanelProps {
  duplicates: DuplicateDetectionResults;
}

export const DuplicateDetectionPanel: React.FC<
  DuplicateDetectionPanelProps
> = ({ duplicates }) => {
  const [activeTab, setActiveTab] = useState<
    'overview' | 'groups' | 'recommendations'
  >('overview');
  const [selectedRiskLevel, setSelectedRiskLevel] = useState<string>('all');
  const [searchTerm, setSearchTerm] = useState('');

  const filteredGroups = duplicates.duplicateGroups.filter((group) => {
    const matchesSearch = group.packageName
      .toLowerCase()
      .includes(searchTerm.toLowerCase());
    const matchesRisk =
      selectedRiskLevel === 'all' || group.riskLevel === selectedRiskLevel;
    return matchesSearch && matchesRisk;
  });

  if (duplicates.totalDuplicates === 0) {
    return (
      <div className="rounded-lg border border-green-200 bg-green-50 p-8 text-center">
        <div className="mb-2 text-green-600">
          <svg
            className="mx-auto h-12 w-12"
            fill="currentColor"
            viewBox="0 0 20 20"
          >
            <path
              fillRule="evenodd"
              d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
              clipRule="evenodd"
            />
          </svg>
        </div>
        <h3 className="mb-1 text-lg font-medium text-green-900">
          No Duplicate Dependencies
        </h3>
        <p className="text-green-700">
          Excellent! Your project has no duplicate dependencies.
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold text-gray-900">
          Duplicate Detection
        </h3>
        <div className="flex rounded-lg bg-gray-100 p-1">
          <button
            onClick={() => setActiveTab('overview')}
            className={`rounded-md px-3 py-1 text-sm font-medium transition-colors ${
              activeTab === 'overview'
                ? 'bg-white text-blue-600 shadow-sm'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            Overview
          </button>
          <button
            onClick={() => setActiveTab('groups')}
            className={`rounded-md px-3 py-1 text-sm font-medium transition-colors ${
              activeTab === 'groups'
                ? 'bg-white text-blue-600 shadow-sm'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            Groups ({duplicates.duplicateGroups.length})
          </button>
          <button
            onClick={() => setActiveTab('recommendations')}
            className={`rounded-md px-3 py-1 text-sm font-medium transition-colors ${
              activeTab === 'recommendations'
                ? 'bg-white text-blue-600 shadow-sm'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            Recommendations
          </button>
        </div>
      </div>

      {/* Content */}
      {activeTab === 'overview' && <OverviewSection duplicates={duplicates} />}
      {activeTab === 'groups' && (
        <GroupsSection
          groups={filteredGroups}
          searchTerm={searchTerm}
          setSearchTerm={setSearchTerm}
          selectedRiskLevel={selectedRiskLevel}
          setSelectedRiskLevel={setSelectedRiskLevel}
        />
      )}
      {activeTab === 'recommendations' && (
        <RecommendationsSection recommendations={duplicates.recommendations} />
      )}
    </div>
  );
};

// Overview Section
const OverviewSection: React.FC<{ duplicates: DuplicateDetectionResults }> = ({
  duplicates,
}) => {
  const riskLevelCounts = duplicates.duplicateGroups.reduce(
    (acc, group) => {
      acc[group.riskLevel] = (acc[group.riskLevel] || 0) + 1;
      return acc;
    },
    {} as Record<string, number>
  );

  return (
    <div className="space-y-6">
      {/* Summary Cards */}
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
        <div className="rounded-lg border border-red-200 bg-red-50 p-4 text-red-700">
          <div className="flex items-center justify-between">
            <div>
              <div className="text-2xl font-bold">
                {duplicates.totalDuplicates}
              </div>
              <div className="text-sm font-medium opacity-80">
                Total Duplicates
              </div>
            </div>
            <div className="opacity-60">
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
            </div>
          </div>
        </div>

        <div className="rounded-lg border border-green-200 bg-green-50 p-4 text-green-700">
          <div className="flex items-center justify-between">
            <div>
              <div className="text-2xl font-bold">
                {duplicates.potentialSavings}
              </div>
              <div className="text-sm font-medium opacity-80">
                Potential Savings
              </div>
            </div>
            <div className="opacity-60">
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
                  d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"
                />
              </svg>
            </div>
          </div>
        </div>

        <div className="rounded-lg border border-yellow-200 bg-yellow-50 p-4 text-yellow-700">
          <div className="flex items-center justify-between">
            <div>
              <div className="text-2xl font-bold">
                {riskLevelCounts.high || 0}
              </div>
              <div className="text-sm font-medium opacity-80">
                High Risk Groups
              </div>
            </div>
            <div className="opacity-60">
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
            </div>
          </div>
        </div>

        <div className="rounded-lg border border-blue-200 bg-blue-50 p-4 text-blue-700">
          <div className="flex items-center justify-between">
            <div>
              <div className="text-2xl font-bold">
                {duplicates.duplicateGroups.length}
              </div>
              <div className="text-sm font-medium opacity-80">
                Duplicate Groups
              </div>
            </div>
            <div className="opacity-60">
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
                  d="M17 14v6m-3-3h6M6 10h2a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v2a2 2 0 002 2zm10 0h2a2 2 0 002-2V6a2 2 0 00-2-2h-2a2 2 0 00-2 2v2a2 2 0 002 2zM6 20h2a2 2 0 002-2v-2a2 2 0 00-2-2H6a2 2 0 00-2 2v2a2 2 0 002 2z"
                />
              </svg>
            </div>
          </div>
        </div>
      </div>

      {/* Risk Distribution */}
      <div className="rounded-lg border border-gray-200 bg-white p-6">
        <h4 className="mb-4 text-lg font-medium text-gray-900">
          Risk Level Distribution
        </h4>
        <div className="space-y-3">
          {(['critical', 'high', 'medium', 'low'] as const).map((level) => {
            const count = riskLevelCounts[level] || 0;
            const total = duplicates.duplicateGroups.length;
            const percentage = total > 0 ? (count / total) * 100 : 0;

            const colors = {
              critical: 'bg-red-500',
              high: 'bg-orange-500',
              medium: 'bg-yellow-500',
              low: 'bg-green-500',
            };

            return (
              <div key={level} className="flex items-center space-x-4">
                <div className="w-20 text-sm font-medium capitalize text-gray-900">
                  {level}
                </div>
                <div className="flex-1">
                  <div className="flex items-center space-x-2">
                    <div className="h-2 flex-1 rounded-full bg-gray-200">
                      <div
                        className={`h-2 rounded-full transition-all duration-300 ${colors[level]}`}
                        style={{ width: `${percentage}%` }}
                      />
                    </div>
                    <div className="w-12 text-right text-sm text-gray-600">
                      {percentage.toFixed(0)}%
                    </div>
                  </div>
                </div>
                <div className="w-12 text-right text-sm text-gray-500">
                  {count}
                </div>
              </div>
            );
          })}
        </div>
      </div>

      {/* Top Duplicate Groups */}
      <div className="rounded-lg border border-gray-200 bg-white p-6">
        <h4 className="mb-4 text-lg font-medium text-gray-900">
          Top Duplicate Groups
        </h4>
        <div className="space-y-3">
          {duplicates.duplicateGroups
            .sort((a, b) => parseFloat(b.wastedSize) - parseFloat(a.wastedSize))
            .slice(0, 5)
            .map((group, index) => (
              <DuplicateGroupSummary
                key={index}
                group={group}
                rank={index + 1}
              />
            ))}
        </div>
      </div>
    </div>
  );
};

// Groups Section
const GroupsSection: React.FC<{
  groups: DuplicateGroup[];
  searchTerm: string;
  setSearchTerm: (term: string) => void;
  selectedRiskLevel: string;
  setSelectedRiskLevel: (level: string) => void;
}> = ({
  groups,
  searchTerm,
  setSearchTerm,
  selectedRiskLevel,
  setSelectedRiskLevel,
}) => {
  return (
    <div className="space-y-4">
      {/* Filters */}
      <div className="flex flex-col gap-4 sm:flex-row">
        <div className="flex-1">
          <input
            type="text"
            placeholder="Search duplicate groups..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm placeholder-gray-500 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div className="sm:w-48">
          <select
            value={selectedRiskLevel}
            onChange={(e) => setSelectedRiskLevel(e.target.value)}
            className="block w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="all">All Risk Levels</option>
            <option value="critical">Critical</option>
            <option value="high">High</option>
            <option value="medium">Medium</option>
            <option value="low">Low</option>
          </select>
        </div>
      </div>

      {/* Results Count */}
      <div className="text-sm text-gray-600">
        Showing {groups.length} duplicate groups
      </div>

      {/* Groups List */}
      <div className="space-y-4">
        {groups.map((group, index) => (
          <DuplicateGroupCard key={index} group={group} />
        ))}
      </div>
    </div>
  );
};

// Recommendations Section
const RecommendationsSection: React.FC<{
  recommendations: DuplicateRecommendation[];
}> = ({ recommendations }) => {
  const groupedRecommendations = recommendations.reduce(
    (acc, rec) => {
      if (!acc[rec.type]) acc[rec.type] = [];
      acc[rec.type].push(rec);
      return acc;
    },
    {} as Record<string, DuplicateRecommendation[]>
  );

  const typeLabels = {
    consolidate: 'Consolidate Versions',
    upgrade: 'Upgrade Dependencies',
    replace: 'Replace Dependencies',
    remove: 'Remove Duplicates',
  };

  const typeIcons = {
    consolidate: (
      <svg
        className="h-5 w-5"
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
    upgrade: (
      <svg
        className="h-5 w-5"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
          d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10"
        />
      </svg>
    ),
    replace: (
      <svg
        className="h-5 w-5"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
          d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4"
        />
      </svg>
    ),
    remove: (
      <svg
        className="h-5 w-5"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
          d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
        />
      </svg>
    ),
  };

  return (
    <div className="space-y-6">
      {Object.entries(groupedRecommendations).map(([type, recs]) => (
        <div
          key={type}
          className="rounded-lg border border-gray-200 bg-white p-6"
        >
          <div className="mb-4 flex items-center space-x-2">
            <div className="text-blue-600">
              {typeIcons[type as keyof typeof typeIcons]}
            </div>
            <h4 className="text-lg font-medium text-gray-900">
              {typeLabels[type as keyof typeof typeLabels]} ({recs.length})
            </h4>
          </div>

          <div className="space-y-4">
            {recs.map((rec, index) => (
              <RecommendationCard key={index} recommendation={rec} />
            ))}
          </div>
        </div>
      ))}
    </div>
  );
};

// Duplicate Group Summary
const DuplicateGroupSummary: React.FC<{
  group: DuplicateGroup;
  rank: number;
}> = ({ group, rank }) => {
  const riskColors = {
    low: 'text-green-600 bg-green-100',
    medium: 'text-yellow-600 bg-yellow-100',
    high: 'text-orange-600 bg-orange-100',
    critical: 'text-red-600 bg-red-100',
  };

  return (
    <div className="flex items-center space-x-4 rounded-lg bg-gray-50 p-3">
      <span className="w-6 font-mono text-sm text-gray-400">#{rank}</span>
      <div className="flex-1">
        <div className="flex items-center space-x-2">
          <span className="font-medium text-gray-900">{group.packageName}</span>
          <span
            className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium capitalize ${
              riskColors[group.riskLevel]
            }`}
          >
            {group.riskLevel}
          </span>
        </div>
        <div className="text-sm text-gray-600">
          {group.versions.length} versions • {group.affectedPackages.length}{' '}
          affected packages
        </div>
      </div>
      <div className="text-right">
        <div className="font-medium text-gray-900">{group.wastedSize}</div>
        <div className="text-sm text-gray-600">wasted</div>
      </div>
    </div>
  );
};

// Duplicate Group Card
const DuplicateGroupCard: React.FC<{ group: DuplicateGroup }> = ({ group }) => {
  const [isExpanded, setIsExpanded] = useState(false);

  const riskColors = {
    low: 'bg-green-50 border-green-200 text-green-800',
    medium: 'bg-yellow-50 border-yellow-200 text-yellow-800',
    high: 'bg-orange-50 border-orange-200 text-orange-800',
    critical: 'bg-red-50 border-red-200 text-red-800',
  };

  return (
    <div className={`rounded-lg border p-4 ${riskColors[group.riskLevel]}`}>
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <div className="mb-2 flex items-center space-x-2">
            <h4 className="font-medium">{group.packageName}</h4>
            <span
              className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium capitalize ${
                riskColors[group.riskLevel]
              }`}
            >
              {group.riskLevel} Risk
            </span>
          </div>

          <div className="mb-3 grid grid-cols-2 gap-4 text-sm opacity-80">
            <div>Total Size: {group.totalSize}</div>
            <div>Wasted: {group.wastedSize}</div>
            <div>Versions: {group.versions.length}</div>
            <div>Affected: {group.affectedPackages.length} packages</div>
          </div>

          {isExpanded && (
            <div className="mt-3 space-y-3 border-t border-current border-opacity-20 pt-3">
              <div>
                <h5 className="mb-2 font-medium">Version Details</h5>
                <div className="space-y-2">
                  {group.versions.map((version, index) => (
                    <div
                      key={index}
                      className="flex items-center justify-between rounded bg-white bg-opacity-50 p-2"
                    >
                      <div className="flex items-center space-x-2">
                        <span className="font-mono text-sm">
                          {version.version}
                        </span>
                        {version.isRecommended && (
                          <span className="inline-flex items-center rounded bg-blue-100 px-1.5 py-0.5 text-xs font-medium text-blue-800">
                            Recommended
                          </span>
                        )}
                      </div>
                      <div className="text-sm opacity-80">
                        {version.size} • {version.usageCount} usage
                        {version.usageCount !== 1 ? 's' : ''}
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              <div>
                <h5 className="mb-2 font-medium">Affected Packages</h5>
                <div className="flex flex-wrap gap-1">
                  {group.affectedPackages.slice(0, 10).map((pkg, index) => (
                    <span
                      key={index}
                      className="inline-flex items-center rounded bg-white bg-opacity-50 px-2 py-0.5 text-xs"
                    >
                      {pkg}
                    </span>
                  ))}
                  {group.affectedPackages.length > 10 && (
                    <span className="inline-flex items-center rounded bg-white bg-opacity-70 px-2 py-0.5 text-xs">
                      +{group.affectedPackages.length - 10} more
                    </span>
                  )}
                </div>
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
  );
};

// Recommendation Card
const RecommendationCard: React.FC<{
  recommendation: DuplicateRecommendation;
}> = ({ recommendation }) => {
  const [isExpanded, setIsExpanded] = useState(false);

  const difficultyColors = {
    easy: 'text-green-700 bg-green-100',
    medium: 'text-yellow-700 bg-yellow-100',
    hard: 'text-red-700 bg-red-100',
  };

  return (
    <div className="rounded-lg border border-gray-200 p-4">
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <div className="mb-2 flex items-center space-x-2">
            <h5 className="font-medium text-gray-900">
              {recommendation.packageName}
            </h5>
            <span
              className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium capitalize ${
                difficultyColors[recommendation.difficulty]
              }`}
            >
              {recommendation.difficulty}
            </span>
          </div>

          <p className="mb-2 text-sm text-gray-600">
            {recommendation.description}
          </p>

          <div className="text-sm text-gray-500">
            Estimated savings: {recommendation.estimatedSavings}
          </div>

          {isExpanded && (
            <div className="mt-3 border-t border-gray-200 pt-3">
              <h6 className="mb-2 font-medium text-gray-900">
                Implementation Steps
              </h6>
              <ol className="space-y-1 text-sm text-gray-600">
                {recommendation.steps.map((step, index) => (
                  <li key={index} className="flex items-start space-x-2">
                    <span className="mt-0.5 flex h-5 w-5 flex-shrink-0 items-center justify-center rounded-full bg-blue-100 text-xs font-medium text-blue-800">
                      {index + 1}
                    </span>
                    <span>{step}</span>
                  </li>
                ))}
              </ol>
            </div>
          )}
        </div>

        <button
          onClick={() => setIsExpanded(!isExpanded)}
          className="ml-4 text-sm text-blue-600 transition-colors hover:text-blue-700"
        >
          {isExpanded ? 'Hide Steps' : 'Show Steps'}
        </button>
      </div>
    </div>
  );
};
