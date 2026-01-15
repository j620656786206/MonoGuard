'use client';

import React, { useState } from 'react';
import { CircularDependency, VersionConflict } from '@monoguard/types';
import { VirtualizedList } from '@/components/ui/VirtualizedList';

export interface CircularDependencyVizProps {
  circularDependencies: CircularDependency[];
  versionConflicts: VersionConflict[];
}

export const CircularDependencyViz: React.FC<CircularDependencyVizProps> = ({
  circularDependencies,
  versionConflicts,
}) => {
  const [activeView, setActiveView] = useState<'circular' | 'conflicts'>('circular');

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold text-gray-900">
          Dependency Analysis
        </h3>
        <div className="flex bg-gray-100 rounded-lg p-1">
          <button
            onClick={() => setActiveView('circular')}
            className={`px-3 py-1 text-sm font-medium rounded-md transition-colors ${
              activeView === 'circular'
                ? 'bg-white text-blue-600 shadow-sm'
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            Circular ({circularDependencies.length})
          </button>
          <button
            onClick={() => setActiveView('conflicts')}
            className={`px-3 py-1 text-sm font-medium rounded-md transition-colors ${
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
  );
};

// Circular Dependencies Panel
const CircularDependenciesPanel: React.FC<{ dependencies: CircularDependency[] }> = ({ 
  dependencies 
}) => {
  if (dependencies.length === 0) {
    return (
      <div className="bg-green-50 border border-green-200 rounded-lg p-8 text-center">
        <div className="text-green-600 mb-2">
          <svg className="mx-auto h-12 w-12" fill="currentColor" viewBox="0 0 20 20">
            <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
          </svg>
        </div>
        <h3 className="text-lg font-medium text-green-900 mb-1">No Circular Dependencies</h3>
        <p className="text-green-700">Great! Your project is free of circular dependencies.</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {dependencies.map((dependency, index) => (
        <CircularDependencyCard key={index} dependency={dependency} />
      ))}
    </div>
  );
};

// Version Conflicts Panel
const VersionConflictsPanel: React.FC<{ conflicts: VersionConflict[] }> = ({ 
  conflicts 
}) => {
  if (conflicts.length === 0) {
    return (
      <div className="bg-green-50 border border-green-200 rounded-lg p-8 text-center">
        <div className="text-green-600 mb-2">
          <svg className="mx-auto h-12 w-12" fill="currentColor" viewBox="0 0 20 20">
            <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
          </svg>
        </div>
        <h3 className="text-lg font-medium text-green-900 mb-1">No Version Conflicts</h3>
        <p className="text-green-700">Excellent! All package versions are compatible.</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {conflicts.map((conflict, index) => (
        <VersionConflictCard key={index} conflict={conflict} />
      ))}
    </div>
  );
};

// Circular Dependency Card
const CircularDependencyCard: React.FC<{ dependency: CircularDependency }> = ({ 
  dependency 
}) => {
  const [isExpanded, setIsExpanded] = useState(false);

  const severityColors = {
    low: 'bg-yellow-50 border-yellow-200 text-yellow-800',
    medium: 'bg-orange-50 border-orange-200 text-orange-800',
    high: 'bg-red-50 border-red-200 text-red-800',
    critical: 'bg-red-100 border-red-300 text-red-900',
  };

  return (
    <div className={`rounded-lg border p-4 ${severityColors[dependency.severity]}`}>
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <div className="flex items-center space-x-2 mb-2">
            <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
              dependency.type === 'direct' 
                ? 'bg-red-100 text-red-800' 
                : 'bg-orange-100 text-orange-800'
            }`}>
              {dependency.type === 'direct' ? 'Direct' : 'Indirect'}
            </span>
            <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium capitalize ${
              severityColors[dependency.severity]
            }`}>
              {dependency.severity}
            </span>
          </div>

          <div className="mb-3">
            <h4 className="font-medium mb-1">Dependency Cycle</h4>
            <div className="flex items-center space-x-2 text-sm font-mono">
              {dependency.cycle.map((dep, index) => (
                <React.Fragment key={index}>
                  <span className="bg-white bg-opacity-50 px-2 py-1 rounded">
                    {dep}
                  </span>
                  {index < dependency.cycle.length - 1 && (
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                    </svg>
                  )}
                  {index === dependency.cycle.length - 1 && (
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4l16 16m0-16L4 20" />
                    </svg>
                  )}
                </React.Fragment>
              ))}
            </div>
          </div>

          <p className="text-sm opacity-80 mb-3">{dependency.impact}</p>

          {isExpanded && (
            <div className="border-t border-current border-opacity-20 pt-3 mt-3">
              <h5 className="font-medium mb-2">Recommendations</h5>
              <ul className="text-sm space-y-1 opacity-80">
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
          className="ml-4 text-sm opacity-60 hover:opacity-80 transition-opacity"
        >
          {isExpanded ? 'Less' : 'More'}
        </button>
      </div>
    </div>
  );
};

// Version Conflict Card
const VersionConflictCard: React.FC<{ conflict: VersionConflict }> = ({ 
  conflict 
}) => {
  const [isExpanded, setIsExpanded] = useState(false);

  const riskLevelColors = {
    low: 'bg-yellow-50 border-yellow-200 text-yellow-800',
    medium: 'bg-orange-50 border-orange-200 text-orange-800',
    high: 'bg-red-50 border-red-200 text-red-800',
    critical: 'bg-red-100 border-red-300 text-red-900',
  };

  return (
    <div className={`rounded-lg border p-4 ${riskLevelColors[conflict.riskLevel]}`}>
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <div className="flex items-center space-x-2 mb-2">
            <h4 className="font-medium">{conflict.packageName}</h4>
            <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium capitalize ${
              riskLevelColors[conflict.riskLevel]
            }`}>
              {conflict.riskLevel} Risk
            </span>
          </div>

          <div className="mb-3">
            <p className="text-sm opacity-80 mb-2">{conflict.impact}</p>
            
            <div className="space-y-2">
              {conflict.conflictingVersions.map((version, index) => (
                <div key={index} className="flex items-center space-x-2">
                  <span className={`inline-flex items-center px-2 py-1 rounded text-xs font-mono ${
                    version.isBreaking 
                      ? 'bg-red-100 text-red-800 border border-red-200' 
                      : 'bg-white bg-opacity-50 border border-current border-opacity-20'
                  }`}>
                    {version.version}
                  </span>
                  {version.isBreaking && (
                    <span className="text-xs text-red-600 font-medium">Breaking</span>
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
            <div className="border-t border-current border-opacity-20 pt-3 mt-3">
              <h5 className="font-medium mb-2">Resolution Strategy</h5>
              <p className="text-sm opacity-80 mb-3">{conflict.resolution}</p>
              
              <h5 className="font-medium mb-2">Affected Packages</h5>
              <div className="max-h-32 overflow-y-auto">
                {conflict.conflictingVersions.map((version, vIndex) => (
                  <div key={vIndex} className="mb-2">
                    <div className="text-xs font-medium opacity-80 mb-1">
                      Version {version.version}:
                    </div>
                    <div className="ml-2 space-y-1">
                      {version.packages.map((pkg, pIndex) => (
                        <div key={pIndex} className="text-xs font-mono opacity-60">
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
          className="ml-4 text-sm opacity-60 hover:opacity-80 transition-opacity"
        >
          {isExpanded ? 'Less' : 'More'}
        </button>
      </div>
    </div>
  );
};