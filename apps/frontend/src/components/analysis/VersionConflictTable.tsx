'use client';

import React, { useState } from 'react';
import { VersionConflict } from '@monoguard/shared-types';
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
    const matchesSearch = conflict.packageName.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesRisk = selectedRiskLevel === 'all' || conflict.riskLevel === selectedRiskLevel;
    return matchesSearch && matchesRisk;
  });

  // Sort conflicts
  const sortedConflicts = [...filteredConflicts].sort((a, b) => {
    let comparison = 0;
    
    switch (sortBy) {
      case 'package':
        comparison = a.packageName.localeCompare(b.packageName);
        break;
      case 'risk':
        const riskOrder = { low: 1, medium: 2, high: 3, critical: 4 };
        comparison = riskOrder[a.riskLevel] - riskOrder[b.riskLevel];
        break;
      case 'versions':
        comparison = a.conflictingVersions.length - b.conflictingVersions.length;
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
      <div className="bg-green-50 border border-green-200 rounded-lg p-8 text-center">
        <div className="text-green-600 mb-2">
          <svg className="mx-auto h-12 w-12" fill="currentColor" viewBox="0 0 20 20">
            <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
          </svg>
        </div>
        <h3 className="text-lg font-medium text-green-900 mb-1">No Version Conflicts</h3>
        <p className="text-green-700">All package versions are compatible.</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Filters */}
      <div className="flex flex-col sm:flex-row gap-4">
        {/* Search */}
        <div className="flex-1">
          <div className="relative">
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg className="h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <input
              type="text"
              placeholder="Search packages..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-lg text-sm placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
        </div>

        {/* Risk Level Filter */}
        <div className="sm:w-48">
          <select
            value={selectedRiskLevel}
            onChange={(e) => setSelectedRiskLevel(e.target.value)}
            className="block w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
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
      <div className="bg-white border border-gray-200 rounded-lg overflow-hidden">
        {/* Header */}
        <div className="bg-gray-50 px-6 py-3 border-b border-gray-200">
          <div className="flex items-center space-x-4">
            <button
              onClick={() => handleSort('package')}
              className="flex items-center space-x-1 text-left font-medium text-gray-900 hover:text-blue-600 transition-colors"
            >
              <span>Package Name</span>
              {sortBy === 'package' && (
                <svg className={`w-4 h-4 ${sortOrder === 'asc' ? 'rotate-180' : ''}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                </svg>
              )}
            </button>

            <button
              onClick={() => handleSort('risk')}
              className="flex items-center space-x-1 text-left font-medium text-gray-900 hover:text-blue-600 transition-colors"
            >
              <span>Risk Level</span>
              {sortBy === 'risk' && (
                <svg className={`w-4 h-4 ${sortOrder === 'asc' ? 'rotate-180' : ''}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                </svg>
              )}
            </button>

            <button
              onClick={() => handleSort('versions')}
              className="flex items-center space-x-1 text-left font-medium text-gray-900 hover:text-blue-600 transition-colors"
            >
              <span>Versions</span>
              {sortBy === 'versions' && (
                <svg className={`w-4 h-4 ${sortOrder === 'asc' ? 'rotate-180' : ''}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
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
        <div className="flex-1 grid grid-cols-1 sm:grid-cols-3 gap-4 items-center">
          {/* Package Name */}
          <div className="font-medium text-gray-900">
            {conflict.packageName}
          </div>

          {/* Risk Level */}
          <div>
            <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium capitalize ${
              riskLevelColors[conflict.riskLevel]
            }`}>
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
          className="ml-4 text-sm text-blue-600 hover:text-blue-700 transition-colors"
        >
          {isExpanded ? 'Hide Details' : 'Show Details'}
        </button>
      </div>

      {isExpanded && (
        <div className="mt-4 pt-4 border-t border-gray-200 space-y-4">
          {/* Impact */}
          <div>
            <h4 className="text-sm font-medium text-gray-900 mb-1">Impact</h4>
            <p className="text-sm text-gray-600">{conflict.impact}</p>
          </div>

          {/* Resolution */}
          <div>
            <h4 className="text-sm font-medium text-gray-900 mb-1">Resolution Strategy</h4>
            <p className="text-sm text-gray-600">{conflict.resolution}</p>
          </div>

          {/* Conflicting Versions */}
          <div>
            <h4 className="text-sm font-medium text-gray-900 mb-2">Conflicting Versions</h4>
            <div className="space-y-2">
              {conflict.conflictingVersions.map((version, index) => (
                <div key={index} className="flex items-start space-x-3 p-3 bg-gray-50 rounded-lg">
                  <div className="flex items-center space-x-2">
                    <span className={`inline-flex items-center px-2 py-1 rounded text-xs font-mono ${
                      version.isBreaking 
                        ? 'bg-red-100 text-red-800 border border-red-200' 
                        : 'bg-white border border-gray-200'
                    }`}>
                      {version.version}
                    </span>
                    {version.isBreaking && (
                      <span className="inline-flex items-center px-1.5 py-0.5 rounded text-xs font-medium bg-red-100 text-red-700">
                        Breaking
                      </span>
                    )}
                  </div>
                  
                  <div className="flex-1 min-w-0">
                    <div className="text-sm text-gray-600 mb-1">
                      Used by {version.packages.length} package{version.packages.length !== 1 ? 's' : ''}:
                    </div>
                    <div className="flex flex-wrap gap-1">
                      {version.packages.slice(0, 5).map((pkg, pkgIndex) => (
                        <span key={pkgIndex} className="inline-flex items-center px-2 py-0.5 rounded text-xs bg-gray-200 text-gray-700">
                          {pkg}
                        </span>
                      ))}
                      {version.packages.length > 5 && (
                        <span className="inline-flex items-center px-2 py-0.5 rounded text-xs bg-gray-300 text-gray-700">
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