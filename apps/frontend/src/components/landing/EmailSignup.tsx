'use client';

import { useState } from 'react';

interface EmailSignupProps {
  title?: string;
  description?: string;
  buttonText?: string;
  variant?: 'error' | 'default' | 'inline';
  className?: string;
}

export function EmailSignup({ 
  title = "Get Early Access",
  description = "Be the first to know when MonoGuard is ready for public repository analysis.",
  buttonText = "Get Notified",
  variant = 'default',
  className = ""
}: EmailSignupProps) {
  const [email, setEmail] = useState('');
  const [feedback, setFeedback] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isSubmitted, setIsSubmitted] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email) return;

    setIsSubmitting(true);
    setError('');

    try {
      // Create form data to match Mailchimp's expected format
      const formData = new FormData();
      formData.append('EMAIL', email);
      formData.append('COMPANY', feedback); // Using COMPANY field for feedback
      formData.append('b_0cd5a0fcdc1d61a426b399d8f_6bbeb33fba', ''); // honeypot field

      // Submit to Mailchimp
      const response = await fetch(
        'https://gmail.us2.list-manage.com/subscribe/post?u=0cd5a0fcdc1d61a426b399d8f&id=6bbeb33fba&f_id=0060c3e1f0',
        {
          method: 'POST',
          body: formData,
          mode: 'no-cors' // This is required for Mailchimp
        }
      );

      // Since we're using no-cors mode, we can't check the response
      // So we assume success if no error is thrown
      setIsSubmitted(true);
      setEmail('');
      setFeedback('');
      
    } catch (err) {
      console.error('Email signup error:', err);
      setError('Something went wrong. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  if (isSubmitted) {
    return (
      <div className={`text-center p-4 ${className}`}>
        <div className="bg-green-50 border border-green-200 rounded-lg p-4">
          <svg className="w-6 h-6 text-green-600 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
          </svg>
          <p className="text-green-800 font-medium">Thanks for subscribing!</p>
          <p className="text-green-700 text-sm mt-1">We'll notify you when MonoGuard is ready.</p>
        </div>
      </div>
    );
  }

  const errorStateContent = variant === 'error' && (
    <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4 mb-4">
      <h3 className="text-yellow-800 font-semibold mb-2">GitHub API Temporarily Unavailable</h3>
      <p className="text-yellow-700 text-sm">
        Due to rate limits, public repository analysis is temporarily limited. 
        Get notified when it's restored and receive our monorepo security checklist!
      </p>
    </div>
  );

  const getContainerClass = () => {
    switch (variant) {
      case 'error':
        return `bg-gradient-to-r from-yellow-50 to-orange-50 border border-yellow-200 rounded-lg p-6 ${className}`;
      case 'inline':
        return `bg-white rounded-lg p-4 ${className}`;
      default:
        return `bg-gradient-to-r from-indigo-50 to-purple-50 rounded-xl p-8 ${className}`;
    }
  };

  return (
    <div className={getContainerClass()}>
      {errorStateContent}
      
      <div className="text-center">
        <h3 className="text-xl font-bold text-gray-900 mb-3">{title}</h3>
        <p className="text-gray-600 mb-6">{description}</p>
        
        <form onSubmit={handleSubmit} className="max-w-sm mx-auto">
          <div className="space-y-3">
            <input
              type="email"
              placeholder="your@email.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              disabled={isSubmitting}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 disabled:bg-gray-100 disabled:cursor-not-allowed"
            />
            <textarea
              placeholder="What features would you like to see? (optional)"
              value={feedback}
              onChange={(e) => setFeedback(e.target.value)}
              disabled={isSubmitting}
              rows={2}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 disabled:bg-gray-100 disabled:cursor-not-allowed resize-none text-sm"
            />
            <button
              type="submit"
              disabled={isSubmitting || !email}
              className="w-full bg-indigo-600 hover:bg-indigo-700 disabled:bg-indigo-400 text-white font-semibold px-6 py-3 rounded-lg transition-colors disabled:cursor-not-allowed"
            >
              {isSubmitting ? (
                <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin mx-auto"></div>
              ) : (
                buttonText
              )}
            </button>
          </div>
          
          {error && (
            <p className="text-red-600 text-sm mt-2">{error}</p>
          )}
          
          <p className="text-xs text-gray-500 mt-3">
            We'll only send you updates about MonoGuard. No spam, ever.
          </p>
        </form>
      </div>
    </div>
  );
}