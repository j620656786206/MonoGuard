'use client'

import React from 'react'
import { useAnalytics } from '../../hooks/useAnalytics'

export function Footer() {
  const { trackClick } = useAnalytics()

  const handleLinkClick = (linkName: string, url?: string) => {
    trackClick(`footer_${linkName}`, linkName, {
      url: url || window.location.href,
    })
  }

  return (
    <footer className="bg-gray-900 text-white">
      <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 gap-8 md:grid-cols-4">
          {/* Brand */}
          <div className="md:col-span-1">
            <h3 className="mb-4 text-2xl font-bold text-white">MonoGuard</h3>
            <p className="text-sm leading-relaxed text-gray-400">
              Comprehensive monorepo architecture analysis and validation tool for modern
              development teams.
            </p>
          </div>

          {/* Product */}
          <div>
            <h4 className="mb-4 text-lg font-semibold text-white">Product</h4>
            <ul className="space-y-2">
              <li>
                <button
                  onClick={() => handleLinkClick('features')}
                  className="text-sm text-gray-400 transition-colors hover:text-white"
                >
                  Features
                </button>
              </li>
              <li>
                <button
                  onClick={() => handleLinkClick('pricing')}
                  className="text-sm text-gray-400 transition-colors hover:text-white"
                >
                  Pricing
                </button>
              </li>
              <li>
                <button
                  onClick={() => handleLinkClick('demo')}
                  className="text-sm text-gray-400 transition-colors hover:text-white"
                >
                  Demo
                </button>
              </li>
              <li>
                <button
                  onClick={() => handleLinkClick('documentation')}
                  className="text-sm text-gray-400 transition-colors hover:text-white"
                >
                  Documentation
                </button>
              </li>
            </ul>
          </div>

          {/* Use Cases */}
          <div>
            <h4 className="mb-4 text-lg font-semibold text-white">Use Cases</h4>
            <ul className="space-y-2">
              <li>
                <button
                  onClick={() => handleLinkClick('monorepo_analysis')}
                  className="text-sm text-gray-400 transition-colors hover:text-white"
                >
                  Monorepo Analysis
                </button>
              </li>
              <li>
                <button
                  onClick={() => handleLinkClick('dependency_management')}
                  className="text-sm text-gray-400 transition-colors hover:text-white"
                >
                  Dependency Management
                </button>
              </li>
              <li>
                <button
                  onClick={() => handleLinkClick('code_quality')}
                  className="text-sm text-gray-400 transition-colors hover:text-white"
                >
                  Code Quality
                </button>
              </li>
              <li>
                <button
                  onClick={() => handleLinkClick('ci_cd_integration')}
                  className="text-sm text-gray-400 transition-colors hover:text-white"
                >
                  CI/CD Integration
                </button>
              </li>
            </ul>
          </div>

          {/* Company */}
          <div>
            <h4 className="mb-4 text-lg font-semibold text-white">Company</h4>
            <ul className="space-y-2">
              <li>
                <button
                  onClick={() => handleLinkClick('about')}
                  className="text-sm text-gray-400 transition-colors hover:text-white"
                >
                  About
                </button>
              </li>
              <li>
                <button
                  onClick={() => handleLinkClick('blog')}
                  className="text-sm text-gray-400 transition-colors hover:text-white"
                >
                  Blog
                </button>
              </li>
              <li>
                <button
                  onClick={() => handleLinkClick('contact')}
                  className="text-sm text-gray-400 transition-colors hover:text-white"
                >
                  Contact
                </button>
              </li>
              <li>
                <button
                  onClick={() => handleLinkClick('privacy')}
                  className="text-sm text-gray-400 transition-colors hover:text-white"
                >
                  Privacy Policy
                </button>
              </li>
            </ul>
          </div>
        </div>

        {/* Developer Preview Banner */}
        <div className="mt-12 border-t border-gray-800 pt-8">
          <div className="rounded-lg bg-gradient-to-r from-indigo-900/50 to-purple-900/50 p-6">
            <div className="flex flex-col md:flex-row md:items-center md:justify-between">
              <div>
                <h4 className="mb-2 text-lg font-semibold text-white">🚀 Developer Preview</h4>
                <p className="text-sm text-gray-300">
                  MonoGuard is currently in Developer Preview. We're actively developing new
                  features and would love your feedback.
                </p>
              </div>
              <div className="mt-4 md:ml-6 md:mt-0">
                <button
                  onClick={() => handleLinkClick('feedback')}
                  className="rounded-full bg-indigo-600 px-6 py-2 text-sm font-medium text-white transition-all duration-200 hover:bg-indigo-700"
                >
                  Give Feedback
                </button>
              </div>
            </div>
          </div>
        </div>

        {/* Bottom Bar */}
        <div className="mt-12 flex flex-col border-t border-gray-800 pt-8 md:flex-row md:items-center md:justify-between">
          <div className="text-sm text-gray-400">
            © {new Date().getFullYear()} MonoGuard. All rights reserved.
          </div>

          <div className="mt-4 flex items-center space-x-6 md:mt-0">
            {/* Social Links - Placeholder for now */}
            <button
              onClick={() => handleLinkClick('github', 'https://github.com')}
              className="text-gray-400 transition-colors hover:text-white"
              aria-label="GitHub"
            >
              <svg className="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z" />
              </svg>
            </button>

            <button
              onClick={() => handleLinkClick('twitter', 'https://twitter.com')}
              className="text-gray-400 transition-colors hover:text-white"
              aria-label="Twitter"
            >
              <svg className="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
                <path d="M23.953 4.57a10 10 0 01-2.825.775 4.958 4.958 0 002.163-2.723c-.951.555-2.005.959-3.127 1.184a4.92 4.92 0 00-8.384 4.482C7.69 8.095 4.067 6.13 1.64 3.162a4.822 4.822 0 00-.666 2.475c0 1.71.87 3.213 2.188 4.096a4.904 4.904 0 01-2.228-.616v.06a4.923 4.923 0 003.946 4.827 4.996 4.996 0 01-2.212.085 4.936 4.936 0 004.604 3.417 9.867 9.867 0 01-6.102 2.105c-.39 0-.779-.023-1.17-.067a13.995 13.995 0 007.557 2.209c9.053 0 13.998-7.496 13.998-13.985 0-.21 0-.42-.015-.63A9.935 9.935 0 0024 4.59z" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </footer>
  )
}
