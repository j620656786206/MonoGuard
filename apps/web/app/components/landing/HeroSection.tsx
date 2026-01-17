'use client'

import React from 'react'
import { EmailSignup } from './EmailSignup'

export function HeroSection() {
  return (
    <section className="relative overflow-hidden bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 py-20">
      <div className="absolute inset-0 bg-white/40"></div>
      <div className="relative mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="text-center">
          {/* Hero Headline */}
          <h1 className="text-4xl font-bold tracking-tight text-gray-900 sm:text-5xl lg:text-6xl">
            <span className="block">Analyze Your</span>
            <span className="block text-indigo-600">Monorepo Health</span>
            <span className="block">in Seconds</span>
          </h1>

          {/* Subheadline */}
          <p className="mx-auto mt-6 max-w-3xl text-xl leading-relaxed text-gray-600 sm:text-2xl">
            <span className="font-semibold text-gray-900">Upload your package.json files</span> for
            instant, privacy-first analysis of dependencies, circular dependencies, and security
            vulnerabilities.
          </p>

          {/* Privacy Badge */}
          <div className="mt-8 flex justify-center">
            <div className="inline-flex items-center rounded-full bg-green-100 px-4 py-2 text-sm font-medium text-green-800">
              <svg className="mr-2 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
                />
              </svg>
              Privacy Focused • Secure Processing • Auto-Delete After Analysis
            </div>
          </div>

          {/* Action Buttons */}
          <div className="mx-auto mt-10 flex max-w-lg flex-col items-center justify-center gap-4 sm:flex-row">
            <button
              onClick={() => {
                // Navigate to upload section
                if (typeof window !== 'undefined') {
                  const uploadSection = document.getElementById('upload-section')
                  if (uploadSection) {
                    uploadSection.scrollIntoView({ behavior: 'smooth' })
                  }
                }
              }}
              className="flex w-full transform items-center justify-center rounded-full bg-indigo-600 px-8 py-4 text-lg font-semibold text-white shadow-lg transition-all duration-200 hover:-translate-y-0.5 hover:bg-indigo-700 hover:shadow-xl sm:w-auto"
            >
              <svg className="mr-2 h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
                />
              </svg>
              Upload Files & Analyze
            </button>
          </div>

          {/* Sample Files Note */}
          <div className="mt-6 text-center">
            <p className="text-sm text-gray-500">
              Don't have files? Scroll down to download sample files for testing.
            </p>
          </div>

          {/* Key Benefits */}
          <div className="mx-auto mt-16 grid max-w-4xl grid-cols-1 gap-8 sm:grid-cols-3">
            <div className="text-center">
              <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-green-100">
                <svg
                  className="h-8 w-8 text-green-600"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M13 10V3L4 14h7v7l9-11h-7z"
                  />
                </svg>
              </div>
              <h3 className="text-lg font-semibold text-gray-900">Instant Analysis</h3>
              <p className="mt-2 text-gray-600">Get comprehensive reports in under 30 seconds</p>
            </div>

            <div className="text-center">
              <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-blue-100">
                <svg
                  className="h-8 w-8 text-blue-600"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
              </div>
              <h3 className="text-lg font-semibold text-gray-900">Zero Setup</h3>
              <p className="mt-2 text-gray-600">
                No installation required - upload and analyze instantly
              </p>
            </div>

            <div className="text-center">
              <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-purple-100">
                <svg
                  className="h-8 w-8 text-purple-600"
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
              </div>
              <h3 className="text-lg font-semibold text-gray-900">Actionable Insights</h3>
              <p className="mt-2 text-gray-600">
                Get specific recommendations to improve code quality
              </p>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}
