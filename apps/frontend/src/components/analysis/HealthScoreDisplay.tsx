'use client';

import React, { useState } from 'react';
import { HealthScore } from '@monoguard/shared-types';

export interface HealthScoreDisplayProps {
  healthScore: HealthScore;
}

export const HealthScoreDisplay: React.FC<HealthScoreDisplayProps> = ({
  healthScore,
}) => {
  const getScoreColor = (score: number) => {
    if (score >= 90) return 'text-green-600';
    if (score >= 80) return 'text-green-500';
    if (score >= 70) return 'text-yellow-500';
    if (score >= 60) return 'text-orange-500';
    return 'text-red-500';
  };

  const getScoreBg = (score: number) => {
    if (score >= 90) return 'bg-green-100 border-green-200';
    if (score >= 80) return 'bg-green-50 border-green-200';
    if (score >= 70) return 'bg-yellow-50 border-yellow-200';
    if (score >= 60) return 'bg-orange-50 border-orange-200';
    return 'bg-red-50 border-red-200';
  };

  const getTrendIcon = () => {
    switch (healthScore.trend) {
      case 'improving':
        return (
          <svg className="w-4 h-4 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
          </svg>
        );
      case 'declining':
        return (
          <svg className="w-4 h-4 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 17h8m0 0V9m0 8l-8-8-4 4-6-6" />
          </svg>
        );
      default:
        return (
          <svg className="w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
          </svg>
        );
    }
  };

  const getTrendColor = () => {
    switch (healthScore.trend) {
      case 'improving': return 'text-green-600 bg-green-100';
      case 'declining': return 'text-red-600 bg-red-100';
      default: return 'text-gray-600 bg-gray-100';
    }
  };

  return (
    <div className="space-y-6">
      {/* Overall Health Score */}
      <div className={`rounded-lg border p-8 text-center ${getScoreBg(healthScore.overall)}`}>
        <div className="flex items-center justify-center space-x-3 mb-4">
          <div className={`text-6xl font-bold ${getScoreColor(healthScore.overall)}`}>
            {healthScore.overall}
          </div>
          <div className="text-left">
            <div className="text-2xl font-semibold text-gray-900">/ 100</div>
            <div className={`flex items-center space-x-1 px-2 py-1 rounded-full text-xs font-medium capitalize ${getTrendColor()}`}>
              {getTrendIcon()}
              <span>{healthScore.trend}</span>
            </div>
          </div>
        </div>
        
        <h2 className="text-2xl font-bold text-gray-900 mb-2">Overall Health Score</h2>
        <p className="text-gray-600">
          Last updated {new Date(healthScore.lastUpdated).toLocaleDateString()}
        </p>
      </div>

      {/* Category Scores */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <ScoreCard 
          title="Dependencies" 
          score={healthScore.dependencies}
          icon={(
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
            </svg>
          )}
        />
        
        <ScoreCard 
          title="Architecture" 
          score={healthScore.architecture}
          icon={(
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
            </svg>
          )}
        />
        
        <ScoreCard 
          title="Maintainability" 
          score={healthScore.maintainability}
          icon={(
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          )}
        />
        
        <ScoreCard 
          title="Security" 
          score={healthScore.security}
          icon={(
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
            </svg>
          )}
        />
        
        <ScoreCard 
          title="Performance" 
          score={healthScore.performance}
          icon={(
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          )}
        />
      </div>

      {/* Health Factors */}
      <div className="bg-white border border-gray-200 rounded-lg p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Health Factors</h3>
        <div className="space-y-4">
          {healthScore.factors.map((factor, index) => (
            <HealthFactor key={index} factor={factor} />
          ))}
        </div>
      </div>

      {/* Recommendations Summary */}
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-6">
        <h3 className="text-lg font-semibold text-blue-900 mb-3">
          Key Recommendations
        </h3>
        <div className="space-y-3">
          {healthScore.factors
            .filter(factor => factor.recommendations.length > 0)
            .slice(0, 5)
            .map((factor, index) => (
              <div key={index} className="flex items-start space-x-2">
                <svg className="w-5 h-5 text-blue-600 mt-0.5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                </svg>
                <div>
                  <div className="font-medium text-blue-900">{factor.name}</div>
                  <div className="text-sm text-blue-700">{factor.recommendations[0]}</div>
                </div>
              </div>
            ))}
        </div>
      </div>
    </div>
  );
};

// Score Card Component
const ScoreCard: React.FC<{
  title: string;
  score: number;
  icon: React.ReactNode;
}> = ({ title, score, icon }) => {
  const getScoreColor = (score: number) => {
    if (score >= 80) return 'text-green-600';
    if (score >= 60) return 'text-yellow-600';
    return 'text-red-600';
  };

  const getScoreBg = (score: number) => {
    if (score >= 80) return 'bg-green-50 border-green-200';
    if (score >= 60) return 'bg-yellow-50 border-yellow-200';
    return 'bg-red-50 border-red-200';
  };

  const circumference = 2 * Math.PI * 45;
  const strokeDasharray = circumference;
  const strokeDashoffset = circumference - (score / 100) * circumference;

  return (
    <div className={`rounded-lg border p-4 ${getScoreBg(score)}`}>
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center space-x-2">
          <div className={`${getScoreColor(score)}`}>{icon}</div>
          <h4 className="font-medium text-gray-900">{title}</h4>
        </div>
        <div className="relative w-12 h-12">
          <svg className="w-12 h-12 -rotate-90" viewBox="0 0 100 100">
            <circle
              cx="50"
              cy="50"
              r="45"
              stroke="currentColor"
              strokeWidth="8"
              fill="transparent"
              className="text-gray-200"
            />
            <circle
              cx="50"
              cy="50"
              r="45"
              stroke="currentColor"
              strokeWidth="8"
              fill="transparent"
              strokeDasharray={strokeDasharray}
              strokeDashoffset={strokeDashoffset}
              strokeLinecap="round"
              className={getScoreColor(score)}
              style={{ transition: 'stroke-dashoffset 1s ease-in-out' }}
            />
          </svg>
          <div className="absolute inset-0 flex items-center justify-center">
            <span className={`text-sm font-bold ${getScoreColor(score)}`}>
              {score}
            </span>
          </div>
        </div>
      </div>
    </div>
  );
};

// Health Factor Component
const HealthFactor: React.FC<{ factor: HealthScore['factors'][0] }> = ({ 
  factor 
}) => {
  const [isExpanded, setIsExpanded] = useState(false);

  const getScoreColor = (score: number) => {
    if (score >= 80) return 'text-green-600';
    if (score >= 60) return 'text-yellow-600';
    return 'text-red-600';
  };

  const getScoreBg = (score: number) => {
    if (score >= 80) return 'bg-green-100';
    if (score >= 60) return 'bg-yellow-100';
    return 'bg-red-100';
  };

  return (
    <div className="border border-gray-200 rounded-lg p-4">
      <div className="flex items-center justify-between">
        <div className="flex-1">
          <div className="flex items-center space-x-3">
            <span className={`inline-flex items-center justify-center w-8 h-8 rounded-full text-sm font-bold ${getScoreColor(factor.score)} ${getScoreBg(factor.score)}`}>
              {factor.score}
            </span>
            <div>
              <h4 className="font-medium text-gray-900">{factor.name}</h4>
              <p className="text-sm text-gray-600">{factor.description}</p>
            </div>
          </div>
          
          {isExpanded && factor.recommendations.length > 0 && (
            <div className="mt-4 pl-11 space-y-2">
              <h5 className="text-sm font-medium text-gray-900">Recommendations:</h5>
              <ul className="text-sm text-gray-600 space-y-1">
                {factor.recommendations.map((rec, index) => (
                  <li key={index} className="flex items-start space-x-2">
                    <span className="text-blue-500 mt-1">•</span>
                    <span>{rec}</span>
                  </li>
                ))}
              </ul>
            </div>
          )}
        </div>

        {factor.recommendations.length > 0 && (
          <button
            onClick={() => setIsExpanded(!isExpanded)}
            className="ml-4 text-sm text-blue-600 hover:text-blue-700 transition-colors"
          >
            {isExpanded ? 'Hide' : 'Show'} Tips
          </button>
        )}
      </div>
    </div>
  );
};