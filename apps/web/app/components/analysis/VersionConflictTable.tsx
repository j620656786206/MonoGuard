'use client';

import React, { useState } from 'react';
import { VersionConflict } from '@monoguard/types';
import { VirtualizedList } from '@/components/ui/VirtualizedList';

export interface VersionConflictTableProps {
  conflicts: VersionConflict[];
}

export const VersionConflictTable: React.FC<VersionConflictTableProps> = ({
  conflicts,
}) => {
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedRiskLevel, setSelectedRiskLevel] = useState<string>('all');
  const [sortBy, setSortBy] = useState<'package' | 'risk' | 'versions'>('risk');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');

  // Filter conflicts
  const filteredConflicts = conflicts.filter((conflict) => {
    const matchesSearch = conflict.packageName
      .toLowerCase()
      .includes(searchTerm.toLowerCase());
    const matchesRisk =
      selectedRiskLevel === 'all' || conflict.riskLevel === selectedRiskLevel;
    return matchesSearch && matchesRisk;
  });

  // Sort conflicts
  const sortedConflicts = [...filteredConflicts].sort((a, b) => {
    let comparison = 0;

    switch (sortBy) {
      case 'package':
        comparison = a.packageName.localeCompare(b.packageName);
        break;
      case 'risk': {
        const riskOrder = { low: 1, medium: 2, high: 3, critical: 4 };
        comparison = riskOrder[a.riskLevel] - riskOrder[b.riskLevel];
        break;
      }
      case 'versions':
        comparison =
          a.conflictingVersions.length - b.conflictingVersions.length;
        break;
    }

    return sortOrder === 'asc' ? comparison : -comparison;
  });

  const riskLevelColors = {
    low: 'text-yellow-600 bg-yellow-100',
    medium: 'text-orange-600 bg-orange-100',
    high: 'text-red-600 bg-red-100',
    critical: 'text-red-700 bg-red-200',
  };

  const handleSort = (column: typeof sortBy) => {
    if (sortBy === column) {
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
    } else {
      setSortBy(column);
      setSortOrder('desc');
    }
  };

  if (conflicts.length === 0) {
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
          No Version Conflicts
        </h3>
        <p className="text-green-700">All package versions are compatible.</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Filters */}
      <div className="flex flex-col gap-4 sm:flex-row">
        {/* Search */}
        <div className="flex-1">
          <div className="relative">
            <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
              <svg
                className="h-4 w-4 text-gray-400"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                />
              </svg>
            </div>
            <input
              type="text"
              placeholder="Search packages..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="block w-full rounded-lg border border-gray-300 py-2 pl-10 pr-3 text-sm placeholder-gray-500 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>

        {/* Risk Level Filter */}
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
        Showing {sortedConflicts.length} of {conflicts.length} conflicts
      </div>

      {/* Table */}
      <div className="overflow-hidden rounded-lg border border-gray-200 bg-white">
        {/* Header */}
        <div className="border-b border-gray-200 bg-gray-50 px-6 py-3">
          <div className="flex items-center space-x-4">
            <button
              onClick={() => handleSort('package')}
              className="flex items-center space-x-1 text-left font-medium text-gray-900 transition-colors hover:text-blue-600"
            >
              <span>Package Name</span>
              {sortBy === 'package' && (
                <svg
                  className={`h-4 w-4 ${sortOrder === 'asc' ? 'rotate-180' : ''}`}
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M19 9l-7 7-7-7"
                  />
                </svg>
              )}
            </button>

            <button
              onClick={() => handleSort('risk')}
              className="flex items-center space-x-1 text-left font-medium text-gray-900 transition-colors hover:text-blue-600"
            >
              <span>Risk Level</span>
              {sortBy === 'risk' && (
                <svg
                  className={`h-4 w-4 ${sortOrder === 'asc' ? 'rotate-180' : ''}`}
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M19 9l-7 7-7-7"
                  />
                </svg>
              )}
            </button>

            <button
              onClick={() => handleSort('versions')}
              className="flex items-center space-x-1 text-left font-medium text-gray-900 transition-colors hover:text-blue-600"
            >
              <span>Versions</span>
              {sortBy === 'versions' && (
                <svg
                  className={`h-4 w-4 ${sortOrder === 'asc' ? 'rotate-180' : ''}`}
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M19 9l-7 7-7-7"
                  />
                </svg>
              )}
            </button>
          </div>
        </div>

        {/* Body */}
        {sortedConflicts.length > 0 ? (
          <div className="divide-y divide-gray-200">
            {sortedConflicts.map((conflict, index) => (
              <ConflictRow key={index} conflict={conflict} />
            ))}
          </div>
        ) : (
          <div className="px-6 py-8 text-center text-gray-500">
            No conflicts found matching your filters.
          </div>
        )}
      </div>
    </div>
  );
};

// Conflict Row Component
const ConflictRow: React.FC<{ conflict: VersionConflict }> = ({ conflict }) => {
  const [isExpanded, setIsExpanded] = useState(false);

  const riskLevelColors = {
    low: 'text-yellow-700 bg-yellow-100',
    medium: 'text-orange-700 bg-orange-100',
    high: 'text-red-700 bg-red-100',
    critical: 'text-red-800 bg-red-200',
  };

  return (
    <div className="px-6 py-4">
      <div className="flex items-center justify-between">
        <div className="grid flex-1 grid-cols-1 items-center gap-4 sm:grid-cols-3">
          {/* Package Name */}
          <div className="font-medium text-gray-900">
            {conflict.packageName}
          </div>

          {/* Risk Level */}
          <div>
            <span
              className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium capitalize ${
                riskLevelColors[conflict.riskLevel]
              }`}
            >
              {conflict.riskLevel}
            </span>
          </div>

          {/* Version Count */}
          <div className="text-sm text-gray-600">
            {conflict.conflictingVersions.length} versions
          </div>
        </div>

        <button
          onClick={() => setIsExpanded(!isExpanded)}
          className="ml-4 text-sm text-blue-600 transition-colors hover:text-blue-700"
        >
          {isExpanded ? 'Hide Details' : 'Show Details'}
        </button>
      </div>

      {isExpanded && (
        <div className="mt-4 space-y-4 border-t border-gray-200 pt-4">
          {/* Impact */}
          <div>
            <h4 className="mb-1 text-sm font-medium text-gray-900">Impact</h4>
            <p className="text-sm text-gray-600">{conflict.impact}</p>
          </div>

          {/* Resolution */}
          <div>
            <h4 className="mb-1 text-sm font-medium text-gray-900">
              Resolution Strategy
            </h4>
            <p className="text-sm text-gray-600">{conflict.resolution}</p>
          </div>

          {/* Conflicting Versions */}
          <div>
            <h4 className="mb-2 text-sm font-medium text-gray-900">
              Conflicting Versions
            </h4>
            <div className="space-y-2">
              {conflict.conflictingVersions.map((version, index) => (
                <div
                  key={index}
                  className="flex items-start space-x-3 rounded-lg bg-gray-50 p-3"
                >
                  <div className="flex items-center space-x-2">
                    <span
                      className={`inline-flex items-center rounded px-2 py-1 font-mono text-xs ${
                        version.isBreaking
                          ? 'border border-red-200 bg-red-100 text-red-800'
                          : 'border border-gray-200 bg-white'
                      }`}
                    >
                      {version.version}
                    </span>
                    {version.isBreaking && (
                      <span className="inline-flex items-center rounded bg-red-100 px-1.5 py-0.5 text-xs font-medium text-red-700">
                        Breaking
                      </span>
                    )}
                  </div>

                  <div className="min-w-0 flex-1">
                    <div className="mb-1 text-sm text-gray-600">
                      Used by {version.packages.length} package
                      {version.packages.length !== 1 ? 's' : ''}:
                    </div>
                    <div className="flex flex-wrap gap-1">
                      {version.packages.slice(0, 5).map((pkg, pkgIndex) => (
                        <span
                          key={pkgIndex}
                          className="inline-flex items-center rounded bg-gray-200 px-2 py-0.5 text-xs text-gray-700"
                        >
                          {pkg}
                        </span>
                      ))}
                      {version.packages.length > 5 && (
                        <span className="inline-flex items-center rounded bg-gray-300 px-2 py-0.5 text-xs text-gray-700">
                          +{version.packages.length - 5} more
                        </span>
                      )}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
