'use client'

import { useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { useAnalytics } from '../../hooks/useAnalytics'

type TabId = 'overview' | 'dependencies' | 'architecture'

const TABS: Array<{ id: TabId; name: string; icon: string }> = [
  { id: 'overview', name: 'Overview', icon: '📊' },
  { id: 'dependencies', name: 'Dependencies', icon: '📦' },
  { id: 'architecture', name: 'Architecture', icon: '🏗️' },
]

export function SampleResults() {
  const navigate = useNavigate()
  const [activeTab, setActiveTab] = useState<TabId>('overview')
  const { trackClick, trackFeatureView } = useAnalytics()

  const handleTabChange = (tab: TabId) => {
    setActiveTab(tab)
    trackClick(`sample_results_tab_${tab}`, tab)
    trackFeatureView(`sample_results_${tab}`)
  }

  const sampleData = {
    overview: {
      healthScore: 87,
      totalDependencies: 342,
      vulnerabilities: 3,
      circularDeps: 2,
      outdatedPackages: 18,
    },
    dependencies: [
      {
        name: '@types/node',
        current: '14.2.1',
        latest: '20.10.5',
        severity: 'medium' as const,
      },
      {
        name: 'lodash',
        current: '4.17.15',
        latest: '4.17.21',
        severity: 'high' as const,
      },
      { name: 'express', current: '4.18.1', latest: '4.18.2', severity: 'low' as const },
      { name: 'react', current: '18.2.0', latest: '18.2.0', severity: null },
    ],
    architecture: [
      { rule: 'Layer separation', status: 'passed' as const, violations: 0 },
      { rule: 'Circular dependencies', status: 'warning' as const, violations: 2 },
      { rule: 'Import restrictions', status: 'passed' as const, violations: 0 },
      { rule: 'Package boundaries', status: 'failed' as const, violations: 5 },
    ],
  }

  return (
    <section className="bg-gray-50 py-20">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        {/* Section Header */}
        <div className="mb-12 text-center">
          <h2 className="text-3xl font-bold text-gray-900 sm:text-4xl">See What You'll Get</h2>
          <p className="mx-auto mt-4 max-w-3xl text-xl text-gray-600">
            Comprehensive analysis results with actionable insights and detailed metrics
          </p>
        </div>

        {/* Sample Report Interface */}
        <div className="mx-auto max-w-5xl">
          <div className="overflow-hidden rounded-xl bg-white shadow-xl">
            {/* Tabs Header */}
            <div className="border-b border-gray-200">
              <nav className="flex space-x-8 px-6" aria-label="Tabs">
                {TABS.map((tab) => (
                  <button
                    key={tab.id}
                    type="button"
                    onClick={() => handleTabChange(tab.id)}
                    className={`border-b-2 px-1 py-4 text-sm font-medium transition-colors ${
                      activeTab === tab.id
                        ? 'border-indigo-500 text-indigo-600'
                        : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700'
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
                      <div className="text-4xl font-bold text-green-600">
                        {sampleData.overview.healthScore}
                      </div>
                      <div className="text-sm text-gray-500">out of 100</div>
                    </div>
                  </div>

                  {/* Metrics Grid */}
                  <div className="grid grid-cols-2 gap-6 md:grid-cols-4">
                    <div className="rounded-lg bg-blue-50 p-4">
                      <div className="text-2xl font-bold text-blue-600">
                        {sampleData.overview.totalDependencies}
                      </div>
                      <div className="text-sm text-blue-700">Total Dependencies</div>
                    </div>
                    <div className="rounded-lg bg-red-50 p-4">
                      <div className="text-2xl font-bold text-red-600">
                        {sampleData.overview.vulnerabilities}
                      </div>
                      <div className="text-sm text-red-700">Vulnerabilities</div>
                    </div>
                    <div className="rounded-lg bg-yellow-50 p-4">
                      <div className="text-2xl font-bold text-yellow-600">
                        {sampleData.overview.circularDeps}
                      </div>
                      <div className="text-sm text-yellow-700">Circular Dependencies</div>
                    </div>
                    <div className="rounded-lg bg-orange-50 p-4">
                      <div className="text-2xl font-bold text-orange-600">
                        {sampleData.overview.outdatedPackages}
                      </div>
                      <div className="text-sm text-orange-700">Outdated Packages</div>
                    </div>
                  </div>
                </div>
              )}

              {activeTab === 'dependencies' && (
                <div>
                  <h3 className="mb-6 text-2xl font-bold text-gray-900">Dependency Analysis</h3>
                  <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200">
                      <thead className="bg-gray-50">
                        <tr>
                          <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
                            Package
                          </th>
                          <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
                            Current Version
                          </th>
                          <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
                            Latest Version
                          </th>
                          <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
                            Severity
                          </th>
                        </tr>
                      </thead>
                      <tbody className="divide-y divide-gray-200 bg-white">
                        {sampleData.dependencies.map((dep, index) => (
                          <tr
                            key={dep.name}
                            className={index % 2 === 0 ? 'bg-white' : 'bg-gray-50'}
                          >
                            <td className="whitespace-nowrap px-6 py-4 font-medium text-gray-900">
                              {dep.name}
                            </td>
                            <td className="whitespace-nowrap px-6 py-4 text-gray-500">
                              {dep.current}
                            </td>
                            <td className="whitespace-nowrap px-6 py-4 text-gray-500">
                              {dep.latest}
                            </td>
                            <td className="whitespace-nowrap px-6 py-4">
                              {dep.severity ? (
                                <span
                                  className={`inline-flex rounded-full px-2 py-1 text-xs font-semibold ${
                                    dep.severity === 'high'
                                      ? 'bg-red-100 text-red-800'
                                      : dep.severity === 'medium'
                                        ? 'bg-yellow-100 text-yellow-800'
                                        : 'bg-green-100 text-green-800'
                                  }`}
                                >
                                  {dep.severity}
                                </span>
                              ) : (
                                <span className="inline-flex rounded-full bg-green-100 px-2 py-1 text-xs font-semibold text-green-800">
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
                  <h3 className="mb-6 text-2xl font-bold text-gray-900">Architecture Validation</h3>
                  <div className="space-y-4">
                    {sampleData.architecture.map((rule) => (
                      <div
                        key={rule.rule}
                        className="flex items-center justify-between rounded-lg border p-4"
                      >
                        <div className="flex items-center space-x-3">
                          <div
                            className={`h-3 w-3 rounded-full ${
                              rule.status === 'passed'
                                ? 'bg-green-500'
                                : rule.status === 'warning'
                                  ? 'bg-yellow-500'
                                  : 'bg-red-500'
                            }`}
                          />
                          <span className="font-medium text-gray-900">{rule.rule}</span>
                        </div>
                        <div className="flex items-center space-x-2">
                          <span
                            className={`rounded-full px-2 py-1 text-xs font-semibold ${
                              rule.status === 'passed'
                                ? 'bg-green-100 text-green-800'
                                : rule.status === 'warning'
                                  ? 'bg-yellow-100 text-yellow-800'
                                  : 'bg-red-100 text-red-800'
                            }`}
                          >
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
          <p className="mb-6 text-gray-600">
            This is just a sample. Your actual results will be tailored to your repository.
          </p>
          <button
            type="button"
            onClick={() => {
              trackClick('try_analysis_sample_results')
              navigate({ to: '/analyze' })
            }}
            className="transform rounded-full bg-indigo-600 px-8 py-3 font-semibold text-white shadow-lg transition-all duration-200 hover:-translate-y-0.5 hover:bg-indigo-700 hover:shadow-xl"
          >
            Try Demo Analysis
          </button>
        </div>
      </div>
    </section>
  )
}
