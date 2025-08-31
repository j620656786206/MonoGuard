'use client';

import React, { useState } from 'react';
import { useAnalytics } from '../../hooks/useAnalytics';

export function SampleResults() {
  const [activeTab, setActiveTab] = useState<'overview' | 'dependencies' | 'architecture'>('overview');
  const { trackClick, trackFeatureView } = useAnalytics();

  const handleTabChange = (tab: 'overview' | 'dependencies' | 'architecture') => {
    setActiveTab(tab);
    trackClick(`sample_results_tab_${tab}`, tab);
    trackFeatureView(`sample_results_${tab}`);
  };

  const sampleData = {
    overview: {
      healthScore: 87,
      totalDependencies: 342,
      vulnerabilities: 3,
      circularDeps: 2,
      outdatedPackages: 18
    },
    dependencies: [
      { name: '@types/node', current: '14.2.1', latest: '20.10.5', severity: 'medium' },
      { name: 'lodash', current: '4.17.15', latest: '4.17.21', severity: 'high' },
      { name: 'express', current: '4.18.1', latest: '4.18.2', severity: 'low' },
      { name: 'react', current: '18.2.0', latest: '18.2.0', severity: null }
    ],
    architecture: [
      { rule: 'Layer separation', status: 'passed', violations: 0 },
      { rule: 'Circular dependencies', status: 'warning', violations: 2 },
      { rule: 'Import restrictions', status: 'passed', violations: 0 },
      { rule: 'Package boundaries', status: 'failed', violations: 5 }
    ]
  };

  return (
    <section className="py-20 bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Section Header */}
        <div className="text-center mb-12">
          <h2 className="text-3xl sm:text-4xl font-bold text-gray-900">
            See What You'll Get
          </h2>
          <p className="mt-4 text-xl text-gray-600 max-w-3xl mx-auto">
            Comprehensive analysis results with actionable insights and detailed metrics
          </p>
        </div>

        {/* Sample Report Interface */}
        <div className="max-w-5xl mx-auto">
          <div className="bg-white rounded-xl shadow-xl overflow-hidden">
            {/* Tabs Header */}
            <div className="border-b border-gray-200">
              <nav className="flex space-x-8 px-6" aria-label="Tabs">
                {[
                  { id: 'overview', name: 'Overview', icon: '📊' },
                  { id: 'dependencies', name: 'Dependencies', icon: '📦' },
                  { id: 'architecture', name: 'Architecture', icon: '🏗️' }
                ].map((tab) => (
                  <button
                    key={tab.id}
                    onClick={() => handleTabChange(tab.id as any)}
                    className={`py-4 px-1 border-b-2 font-medium text-sm transition-colors ${
                      activeTab === tab.id
                        ? 'border-indigo-500 text-indigo-600'
                        : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                    }`}
                  >
                    <span className="mr-2">{tab.icon}</span>
                    {tab.name}
                  </button>
                ))}
              </nav>
            </div>

            {/* Tab Content */}
            <div className="p-8">
              {activeTab === 'overview' && (
                <div className="space-y-8">
                  {/* Health Score */}
                  <div className="flex items-center justify-between">
                    <div>
                      <h3 className="text-2xl font-bold text-gray-900">Project Health Score</h3>
                      <p className="text-gray-600">Overall assessment of your codebase</p>
                    </div>
                    <div className="text-right">
                      <div className="text-4xl font-bold text-green-600">{sampleData.overview.healthScore}</div>
                      <div className="text-sm text-gray-500">out of 100</div>
                    </div>
                  </div>

                  {/* Metrics Grid */}
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
                    <div className="bg-blue-50 p-4 rounded-lg">
                      <div className="text-2xl font-bold text-blue-600">{sampleData.overview.totalDependencies}</div>
                      <div className="text-sm text-blue-700">Total Dependencies</div>
                    </div>
                    <div className="bg-red-50 p-4 rounded-lg">
                      <div className="text-2xl font-bold text-red-600">{sampleData.overview.vulnerabilities}</div>
                      <div className="text-sm text-red-700">Vulnerabilities</div>
                    </div>
                    <div className="bg-yellow-50 p-4 rounded-lg">
                      <div className="text-2xl font-bold text-yellow-600">{sampleData.overview.circularDeps}</div>
                      <div className="text-sm text-yellow-700">Circular Dependencies</div>
                    </div>
                    <div className="bg-orange-50 p-4 rounded-lg">
                      <div className="text-2xl font-bold text-orange-600">{sampleData.overview.outdatedPackages}</div>
                      <div className="text-sm text-orange-700">Outdated Packages</div>
                    </div>
                  </div>
                </div>
              )}

              {activeTab === 'dependencies' && (
                <div>
                  <h3 className="text-2xl font-bold text-gray-900 mb-6">Dependency Analysis</h3>
                  <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200">
                      <thead className="bg-gray-50">
                        <tr>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Package
                          </th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Current Version
                          </th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Latest Version
                          </th>
                          <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Severity
                          </th>
                        </tr>
                      </thead>
                      <tbody className="bg-white divide-y divide-gray-200">
                        {sampleData.dependencies.map((dep, index) => (
                          <tr key={index} className={index % 2 === 0 ? 'bg-white' : 'bg-gray-50'}>
                            <td className="px-6 py-4 whitespace-nowrap font-medium text-gray-900">
                              {dep.name}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-gray-500">
                              {dep.current}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-gray-500">
                              {dep.latest}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap">
                              {dep.severity ? (
                                <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                                  dep.severity === 'high' ? 'bg-red-100 text-red-800' :
                                  dep.severity === 'medium' ? 'bg-yellow-100 text-yellow-800' :
                                  'bg-green-100 text-green-800'
                                }`}>
                                  {dep.severity}
                                </span>
                              ) : (
                                <span className="inline-flex px-2 py-1 text-xs font-semibold rounded-full bg-green-100 text-green-800">
                                  up to date
                                </span>
                              )}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>
              )}

              {activeTab === 'architecture' && (
                <div>
                  <h3 className="text-2xl font-bold text-gray-900 mb-6">Architecture Validation</h3>
                  <div className="space-y-4">
                    {sampleData.architecture.map((rule, index) => (
                      <div key={index} className="flex items-center justify-between p-4 border rounded-lg">
                        <div className="flex items-center space-x-3">
                          <div className={`w-3 h-3 rounded-full ${
                            rule.status === 'passed' ? 'bg-green-500' :
                            rule.status === 'warning' ? 'bg-yellow-500' :
                            'bg-red-500'
                          }`}></div>
                          <span className="font-medium text-gray-900">{rule.rule}</span>
                        </div>
                        <div className="flex items-center space-x-2">
                          <span className={`px-2 py-1 text-xs font-semibold rounded-full ${
                            rule.status === 'passed' ? 'bg-green-100 text-green-800' :
                            rule.status === 'warning' ? 'bg-yellow-100 text-yellow-800' :
                            'bg-red-100 text-red-800'
                          }`}>
                            {rule.status}
                          </span>
                          {rule.violations > 0 && (
                            <span className="text-sm text-gray-500">
                              {rule.violations} violations
                            </span>
                          )}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Call to Action */}
        <div className="mt-12 text-center">
          <p className="text-gray-600 mb-6">
            This is just a sample. Your actual results will be tailored to your repository.
          </p>
          <button
            onClick={() => {
              trackClick('try_analysis_sample_results');
              document.getElementById('hero-section')?.scrollIntoView({ behavior: 'smooth' });
            }}
            className="bg-indigo-600 hover:bg-indigo-700 text-white font-semibold py-3 px-8 rounded-full transition-all duration-200 shadow-lg hover:shadow-xl transform hover:-translate-y-0.5"
          >
            Analyze Your Repository
          </button>
        </div>
      </div>
    </section>
  );
}