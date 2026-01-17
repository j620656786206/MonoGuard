'use client'

import { useState } from 'react'

interface EmailSignupProps {
  title?: string
  description?: string
  buttonText?: string
  variant?: 'error' | 'default' | 'inline'
  className?: string
}

export function EmailSignup({
  title = 'Get Early Access',
  description = 'Be the first to know when MonoGuard is ready for public repository analysis.',
  buttonText = 'Get Notified',
  variant = 'default',
  className = '',
}: EmailSignupProps) {
  const [email, setEmail] = useState('')
  const [feedback, setFeedback] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [isSubmitted, setIsSubmitted] = useState(false)
  const [error, setError] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!email) return

    setIsSubmitting(true)
    setError('')

    try {
      // Create form data to match Mailchimp's expected format
      const formData = new FormData()
      formData.append('EMAIL', email)
      formData.append('COMPANY', feedback) // Using COMPANY field for feedback
      formData.append('b_0cd5a0fcdc1d61a426b399d8f_6bbeb33fba', '') // honeypot field

      // Submit to Mailchimp
      const response = await fetch(
        'https://gmail.us2.list-manage.com/subscribe/post?u=0cd5a0fcdc1d61a426b399d8f&id=6bbeb33fba&f_id=0060c3e1f0',
        {
          method: 'POST',
          body: formData,
          mode: 'no-cors', // This is required for Mailchimp
        }
      )

      // Since we're using no-cors mode, we can't check the response
      // So we assume success if no error is thrown
      setIsSubmitted(true)
      setEmail('')
      setFeedback('')
    } catch (err) {
      console.error('Email signup error:', err)
      setError('Something went wrong. Please try again.')
    } finally {
      setIsSubmitting(false)
    }
  }

  if (isSubmitted) {
    return (
      <div className={`p-4 text-center ${className}`}>
        <div className="rounded-lg border border-green-200 bg-green-50 p-4">
          <svg
            className="mx-auto mb-2 h-6 w-6 text-green-600"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
          </svg>
          <p className="font-medium text-green-800">Thanks for subscribing!</p>
          <p className="mt-1 text-sm text-green-700">We'll notify you when MonoGuard is ready.</p>
        </div>
      </div>
    )
  }

  const errorStateContent = variant === 'error' && (
    <div className="mb-4 rounded-lg border border-yellow-200 bg-yellow-50 p-4">
      <h3 className="mb-2 font-semibold text-yellow-800">GitHub API Temporarily Unavailable</h3>
      <p className="text-sm text-yellow-700">
        Due to rate limits, public repository analysis is temporarily limited. Get notified when
        it's restored and receive our monorepo security checklist!
      </p>
    </div>
  )

  const getContainerClass = () => {
    switch (variant) {
      case 'error':
        return `bg-gradient-to-r from-yellow-50 to-orange-50 border border-yellow-200 rounded-lg p-6 ${className}`
      case 'inline':
        return `bg-white rounded-lg p-4 ${className}`
      default:
        return `bg-gradient-to-r from-indigo-50 to-purple-50 rounded-xl p-8 ${className}`
    }
  }

  return (
    <div className={getContainerClass()}>
      {errorStateContent}

      <div className="text-center">
        <h3 className="mb-3 text-xl font-bold text-gray-900">{title}</h3>
        <p className="mb-6 text-gray-600">{description}</p>

        <form onSubmit={handleSubmit} className="mx-auto max-w-sm">
          <div className="space-y-3">
            <input
              type="email"
              placeholder="your@email.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              disabled={isSubmitting}
              className="w-full rounded-lg border border-gray-300 px-4 py-3 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500 disabled:cursor-not-allowed disabled:bg-gray-100"
            />
            <textarea
              placeholder="What features would you like to see? (optional)"
              value={feedback}
              onChange={(e) => setFeedback(e.target.value)}
              disabled={isSubmitting}
              rows={2}
              className="w-full resize-none rounded-lg border border-gray-300 px-4 py-3 text-sm focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500 disabled:cursor-not-allowed disabled:bg-gray-100"
            />
            <button
              type="submit"
              disabled={isSubmitting || !email}
              className="w-full rounded-lg bg-indigo-600 px-6 py-3 font-semibold text-white transition-colors hover:bg-indigo-700 disabled:cursor-not-allowed disabled:bg-indigo-400"
            >
              {isSubmitting ? (
                <div className="mx-auto h-5 w-5 animate-spin rounded-full border-2 border-white border-t-transparent"></div>
              ) : (
                buttonText
              )}
            </button>
          </div>

          {error && <p className="mt-2 text-sm text-red-600">{error}</p>}

          <p className="mt-3 text-xs text-gray-500">
            We'll only send you updates about MonoGuard. No spam, ever.
          </p>
        </form>
      </div>
    </div>
  )
}
