'use client';

import React from 'react';
import { useAnalytics } from '../../hooks/useAnalytics';

export function FeaturesSection() {
  const { trackFeatureView, trackClick } = useAnalytics();

  const handleFeatureClick = (featureName: string) => {
    trackClick(`feature_${featureName}`, featureName);
    trackFeatureView(featureName);
  };

  const features = [
    {
      id: 'dependency_analysis',
      title: 'Dependency Analysis',
      description:
        'Identify outdated packages, version conflicts, and security vulnerabilities across your entire monorepo.',
      icon: (
        <svg
          className="h-8 w-8"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9v-9m0-9v9m0 9c-5 0-9-4-9-9s4-9 9-9"
          />
        </svg>
      ),
      benefits: [
        'Detect vulnerable packages',
        'Find version mismatches',
        'Identify unused dependencies',
        'Track dependency tree depth',
      ],
    },
    {
      id: 'circular_dependencies',
      title: 'Circular Dependency Detection',
      description:
        'Automatically detect and visualize circular dependencies that can cause build issues and performance problems.',
      icon: (
        <svg
          className="h-8 w-8"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
          />
        </svg>
      ),
      benefits: [
        'Visual dependency graphs',
        'Breaking change impact analysis',
        'Performance bottleneck identification',
        'Refactoring recommendations',
      ],
    },
    {
      id: 'architecture_validation',
      title: 'Architecture Validation',
      description:
        'Enforce architectural boundaries and validate layered architecture patterns across your codebase.',
      icon: (
        <svg
          className="h-8 w-8"
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
      benefits: [
        'Layer boundary enforcement',
        'Package coupling analysis',
        'Clean architecture validation',
        'Custom rule definitions',
      ],
    },
    {
      id: 'health_scoring',
      title: 'Project Health Score',
      description:
        'Get an overall health score based on dependency freshness, architecture quality, and best practices.',
      icon: (
        <svg
          className="h-8 w-8"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v4a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
          />
        </svg>
      ),
      benefits: [
        'Overall project health metrics',
        'Trend analysis over time',
        'Benchmark against standards',
        'Prioritized improvement suggestions',
      ],
    },
    {
      id: 'monorepo_support',
      title: 'Monorepo Intelligence',
      description:
        'Specialized analysis for monorepos including workspace dependency management and cross-package analysis.',
      icon: (
        <svg
          className="h-8 w-8"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
          />
        </svg>
      ),
      benefits: [
        'Workspace dependency analysis',
        'Cross-package impact analysis',
        'Shared dependency optimization',
        'Package.json validation',
      ],
    },
    {
      id: 'instant_reports',
      title: 'Instant Reports',
      description:
        'Generate comprehensive HTML and JSON reports that can be shared with your team or integrated into CI/CD.',
      icon: (
        <svg
          className="h-8 w-8"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
          />
        </svg>
      ),
      benefits: [
        'Shareable HTML reports',
        'JSON data for automation',
        'Visual dependency graphs',
        'Executive summary dashboards',
      ],
    },
  ];

  return (
    <section className="bg-white py-20">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        {/* Section Header */}
        <div className="text-center">
          <h2 className="text-3xl font-bold text-gray-900 sm:text-4xl">
            Comprehensive Analysis Features
          </h2>
          <p className="mx-auto mt-4 max-w-3xl text-xl text-gray-600">
            Everything you need to maintain healthy, scalable JavaScript and
            TypeScript projects
          </p>
        </div>

        {/* Features Grid */}
        <div className="mt-16 grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-3">
          {features.map((feature) => (
            <div
              key={feature.id}
              className="group relative cursor-pointer"
              onClick={() => handleFeatureClick(feature.id)}
            >
              <div className="h-full rounded-xl border border-gray-200 bg-white p-8 transition-all duration-300 hover:-translate-y-1 hover:shadow-xl">
                {/* Icon */}
                <div className="flex h-16 w-16 items-center justify-center rounded-xl bg-indigo-100 text-indigo-600 transition-colors duration-300 group-hover:bg-indigo-600 group-hover:text-white">
                  {feature.icon}
                </div>

                {/* Content */}
                <h3 className="mt-6 text-xl font-semibold text-gray-900 transition-colors group-hover:text-indigo-600">
                  {feature.title}
                </h3>

                <p className="mt-3 leading-relaxed text-gray-600">
                  {feature.description}
                </p>

                {/* Benefits List */}
                <ul className="mt-6 space-y-2">
                  {feature.benefits.map((benefit, index) => (
                    <li
                      key={index}
                      className="flex items-center text-sm text-gray-500"
                    >
                      <svg
                        className="mr-2 h-4 w-4 flex-shrink-0 text-green-500"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M5 13l4 4L19 7"
                        />
                      </svg>
                      {benefit}
                    </li>
                  ))}
                </ul>
              </div>
            </div>
          ))}
        </div>

        {/* CTA Section */}
        <div className="mt-20 text-center">
          <div className="rounded-2xl bg-gradient-to-r from-indigo-50 to-purple-50 px-8 py-12">
            <h3 className="mb-4 text-2xl font-bold text-gray-900">
              Ready to analyze your repository?
            </h3>
            <p className="mx-auto mb-8 max-w-2xl text-gray-600">
              Get started with our comprehensive analysis and take the first
              step towards improving your codebase health.
            </p>
            <button
              onClick={() => {
                trackClick('get_started_features_cta');
                if (typeof window !== 'undefined') {
                  document
                    .getElementById('hero-section')
                    ?.scrollIntoView({ behavior: 'smooth' });
                }
              }}
              className="transform rounded-full bg-indigo-600 px-8 py-3 font-semibold text-white shadow-lg transition-all duration-200 hover:-translate-y-0.5 hover:bg-indigo-700 hover:shadow-xl"
            >
              Get Started Now
            </button>
          </div>
        </div>
      </div>
    </section>
  );
}
