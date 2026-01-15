'use client';

import React, { useState } from 'react';
import { BundleImpactReport, BundleBreakdown } from '@monoguard/types';

export interface BundleImpactChartProps {
  bundleImpact: BundleImpactReport;
}

export const BundleImpactChart: React.FC<BundleImpactChartProps> = ({
  bundleImpact,
}) => {
  const [sortBy, setSortBy] = useState<'size' | 'percentage' | 'duplicates'>('size');
  const [showOnlyDuplicates, setShowOnlyDuplicates] = useState(false);

  const filteredBreakdown = showOnlyDuplicates
    ? bundleImpact.breakdown.filter(item => item.duplicates > 0)
    : bundleImpact.breakdown;

  const sortedBreakdown = [...filteredBreakdown].sort((a, b) => {
    switch (sortBy) {
      case 'size':
        return parseFloat(b.size) - parseFloat(a.size);
      case 'percentage':
        return b.percentage - a.percentage;
      case 'duplicates':
        return b.duplicates - a.duplicates;
      default:
        return 0;
    }
  });

  const formatSize = (size: string): number => {
    const match = size.match(/([\d.]+)\s*(KB|MB|GB)/i);
    if (!match) return 0;
    
    const value = parseFloat(match[1]);
    const unit = match[2].toUpperCase();
    
    switch (unit) {
      case 'GB': return value * 1024 * 1024;
      case 'MB': return value * 1024;
      case 'KB': return value;
      default: return value;
    }
  };

  const maxSize = Math.max(...sortedBreakdown.map(item => formatSize(item.size)));

  return (
    <div className="space-y-6">
      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="bg-blue-50 border-blue-200 text-blue-700 rounded-lg border p-4">
          <div className="flex items-center justify-between">
            <div>
              <div className="text-2xl font-bold">{bundleImpact.totalSize}</div>
              <div className="text-sm font-medium opacity-80">Total Bundle Size</div>
            </div>
            <div className="opacity-60">
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 7v10c0 2.21 3.79 4 8.5 4s8.5-1.79 8.5-4V7M4 7c0 2.21 3.79 4 8.5 4s8.5-1.79 8.5-4M4 7c0-2.21 3.79-4 8.5-4s8.5 1.79 8.5 4" />
              </svg>
            </div>
          </div>
        </div>

        <div className="bg-yellow-50 border-yellow-200 text-yellow-700 rounded-lg border p-4">
          <div className="flex items-center justify-between">
            <div>
              <div className="text-2xl font-bold">{bundleImpact.duplicateSize}</div>
              <div className="text-sm font-medium opacity-80">Duplicate Size</div>
            </div>
            <div className="opacity-60">
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
            </div>
          </div>
        </div>

        <div className="bg-red-50 border-red-200 text-red-700 rounded-lg border p-4">
          <div className="flex items-center justify-between">
            <div>
              <div className="text-2xl font-bold">{bundleImpact.unusedSize}</div>
              <div className="text-sm font-medium opacity-80">Unused Size</div>
            </div>
            <div className="opacity-60">
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </div>
          </div>
        </div>

        <div className="bg-green-50 border-green-200 text-green-700 rounded-lg border p-4">
          <div className="flex items-center justify-between">
            <div>
              <div className="text-2xl font-bold">{bundleImpact.potentialSavings}</div>
              <div className="text-sm font-medium opacity-80">Potential Savings</div>
            </div>
            <div className="opacity-60">
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
              </svg>
            </div>
          </div>
        </div>
      </div>

      {/* Controls */}
      <div className="flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
        <div className="flex items-center space-x-4">
          <div className="flex items-center space-x-2">
            <label className="text-sm font-medium text-gray-700">Sort by:</label>
            <select
              value={sortBy}
              onChange={(e) => setSortBy(e.target.value as typeof sortBy)}
              className="text-sm border border-gray-300 rounded px-2 py-1 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="size">Size</option>
              <option value="percentage">Percentage</option>
              <option value="duplicates">Duplicates</option>
            </select>
          </div>

          <label className="flex items-center space-x-2">
            <input
              type="checkbox"
              checked={showOnlyDuplicates}
              onChange={(e) => setShowOnlyDuplicates(e.target.checked)}
              className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
            />
            <span className="text-sm text-gray-700">Show only packages with duplicates</span>
          </label>
        </div>

        <div className="text-sm text-gray-600">
          Showing {sortedBreakdown.length} of {bundleImpact.breakdown.length} packages
        </div>
      </div>

      {/* Bundle Breakdown */}
      <div className="bg-white border border-gray-200 rounded-lg overflow-hidden">
        {/* Header */}
        <div className="bg-gray-50 px-6 py-3 border-b border-gray-200">
          <div className="grid grid-cols-12 gap-4 text-sm font-medium text-gray-900">
            <div className="col-span-4">Package Name</div>
            <div className="col-span-2">Size</div>
            <div className="col-span-2">Percentage</div>
            <div className="col-span-2">Duplicates</div>
            <div className="col-span-2">Visual</div>
          </div>
        </div>

        {/* Body */}
        <div className="divide-y divide-gray-200">
          {sortedBreakdown.length > 0 ? (
            sortedBreakdown.map((item, index) => (
              <BundleBreakdownRow
                key={index}
                item={item}
                maxSize={maxSize}
                rank={index + 1}
              />
            ))
          ) : (
            <div className="px-6 py-8 text-center text-gray-500">
              No packages found matching your filters.
            </div>
          )}
        </div>
      </div>

      {/* Recommendations */}
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-6">
        <h4 className="text-lg font-medium text-blue-900 mb-3">
          Bundle Optimization Recommendations
        </h4>
        <div className="space-y-3 text-blue-800">
          <div className="flex items-start space-x-2">
            <svg className="w-5 h-5 text-blue-600 mt-0.5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
            </svg>
            <div>
              <div className="font-medium">Eliminate duplicate dependencies</div>
              <div className="text-sm text-blue-700">
                Could save {bundleImpact.duplicateSize} by consolidating duplicate packages
              </div>
            </div>
          </div>
          
          <div className="flex items-start space-x-2">
            <svg className="w-5 h-5 text-blue-600 mt-0.5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
            </svg>
            <div>
              <div className="font-medium">Remove unused dependencies</div>
              <div className="text-sm text-blue-700">
                Additional {bundleImpact.unusedSize} can be saved by removing unused packages
              </div>
            </div>
          </div>
          
          <div className="flex items-start space-x-2">
            <svg className="w-5 h-5 text-blue-600 mt-0.5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
            </svg>
            <div>
              <div className="font-medium">Consider tree shaking</div>
              <div className="text-sm text-blue-700">
                Enable tree shaking to eliminate unused code from large packages
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

// Bundle Breakdown Row Component
const BundleBreakdownRow: React.FC<{
  item: BundleBreakdown;
  maxSize: number;
  rank: number;
}> = ({ item, maxSize, rank }) => {
  const sizeKB = formatSize(item.size);
  const widthPercent = maxSize > 0 ? (sizeKB / maxSize) * 100 : 0;

  return (
    <div className="px-6 py-3">
      <div className="grid grid-cols-12 gap-4 items-center">
        {/* Package Name */}
        <div className="col-span-4 flex items-center space-x-2">
          <span className="text-xs text-gray-400 font-mono w-6">#{rank}</span>
          <span className="font-medium text-gray-900 truncate">{item.packageName}</span>
        </div>

        {/* Size */}
        <div className="col-span-2 text-sm text-gray-600 font-mono">{item.size}</div>

        {/* Percentage */}
        <div className="col-span-2 text-sm text-gray-600">
          {item.percentage.toFixed(1)}%
        </div>

        {/* Duplicates */}
        <div className="col-span-2">
          {item.duplicates > 0 ? (
            <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
              {item.duplicates} duplicate{item.duplicates !== 1 ? 's' : ''}
            </span>
          ) : (
            <span className="text-sm text-gray-400">None</span>
          )}
        </div>

        {/* Visual Bar */}
        <div className="col-span-2">
          <div className="w-full bg-gray-200 rounded-full h-2">
            <div
              className={`h-2 rounded-full transition-all duration-300 ${
                item.duplicates > 0 ? 'bg-red-500' : 'bg-blue-500'
              }`}
              style={{ width: `${Math.max(widthPercent, 2)}%` }}
            />
          </div>
        </div>
      </div>
    </div>
  );
};

// Helper function to format size
const formatSize = (size: string): number => {
  const match = size.match(/([\d.]+)\s*(B|KB|MB|GB)/i);
  if (!match) return 0;
  
  const value = parseFloat(match[1]);
  const unit = match[2].toUpperCase();
  
  switch (unit) {
    case 'GB': return value * 1024 * 1024;
    case 'MB': return value * 1024;
    case 'KB': return value;
    case 'B': return value / 1024;
    default: return value;
  }
};