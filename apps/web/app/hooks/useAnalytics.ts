import { useLocation } from '@tanstack/react-router'
import { useCallback, useEffect, useRef } from 'react'
import type { EventProperties, EventType } from '../lib/api/services/analytics'

export function useAnalytics() {
  const { pathname } = useLocation()
  const previousPathname = useRef<string | null>(null)

  // Track page views automatically (simplified for MVP - no session required)
  useEffect(() => {
    // Track when pathname changes
    if (pathname !== previousPathname.current) {
      console.log(`Page view: ${pathname}`) // Simple console logging for MVP
      previousPathname.current = pathname
    }
  }, [pathname])

  // Generic event tracking (simplified for MVP - console logging)
  const trackEvent = useCallback(
    (
      eventType: EventType,
      eventName: string,
      properties?: EventProperties,
      options?: {
        page?: string
        element?: string
        value?: string
        duration?: number
      }
    ) => {
      console.log(`Event: ${eventType} - ${eventName}`, {
        properties,
        options,
      })
    },
    []
  )

  // Click tracking
  const trackClick = useCallback(
    (element: string, value?: string, properties?: EventProperties) => {
      console.log(`Click: ${element}`, { value, properties })
    },
    []
  )

  // Form submission tracking
  const trackFormSubmit = useCallback((formName: string, properties?: EventProperties) => {
    console.log(`Form Submit: ${formName}`, properties)
  }, [])

  // Analysis tracking
  const trackAnalysisStart = useCallback((analysisType: string, properties?: EventProperties) => {
    console.log(`Analysis Start: ${analysisType}`, properties)
  }, [])

  const trackAnalysisComplete = useCallback(
    (analysisType: string, duration: number, healthScore: number, properties?: EventProperties) => {
      console.log(`Analysis Complete: ${analysisType}`, {
        duration,
        healthScore,
        properties,
      })
    },
    []
  )

  // Error tracking
  const trackError = useCallback((errorMessage: string, errorCode?: string, page?: string) => {
    console.log(`Error: ${errorMessage}`, { errorCode, page })
  }, [])

  // Feature tracking
  const trackFeatureView = useCallback((featureName: string, properties?: EventProperties) => {
    console.log(`Feature View: ${featureName}`, properties)
  }, [])

  // GitHub analysis tracking
  const trackGitHubAnalysis = useCallback(
    (repositoryUrl: string, branch: string, properties?: EventProperties) => {
      console.log(`GitHub Analysis: ${repositoryUrl}@${branch}`, properties)
    },
    []
  )

  // Upload tracking
  const trackUploadStart = useCallback((fileName: string, fileSize: number) => {
    console.log(`Upload Start: ${fileName}`, { fileSize })
  }, [])

  const trackUploadComplete = useCallback(
    (fileName: string, fileSize: number, duration: number) => {
      console.log(`Upload Complete: ${fileName}`, { fileSize, duration })
    },
    []
  )

  // Download tracking
  const trackDownload = useCallback((fileName: string, properties?: EventProperties) => {
    console.log(`Download: ${fileName}`, properties)
  }, [])

  // Conversion funnel tracking
  const trackConversionStep = useCallback(
    (step: string, stepOrder: number, completed: boolean, duration?: number) => {
      console.log(`Conversion Step: ${step}`, {
        stepOrder,
        completed,
        duration,
      })
    },
    []
  )

  // Landing page conversion funnel helpers
  const trackLandingPageView = useCallback(() => {
    console.log('Landing Page View')
  }, [])

  const trackAnalysisIntent = useCallback(() => {
    console.log('Analysis Intent')
  }, [])

  const trackAnalysisAttempt = useCallback(() => {
    console.log('Analysis Attempt')
  }, [])

  const trackAnalysisSuccess = useCallback(() => {
    console.log('Analysis Success')
  }, [])

  const trackAnalysisAbandoned = useCallback((reason?: string) => {
    console.log('Analysis Abandoned', { reason })
  }, [])

  return {
    // General tracking
    trackEvent,
    trackClick,
    trackFormSubmit,
    trackError,
    trackFeatureView,

    // Analysis tracking
    trackAnalysisStart,
    trackAnalysisComplete,
    trackGitHubAnalysis,

    // File operations
    trackUploadStart,
    trackUploadComplete,
    trackDownload,

    // Conversion funnel
    trackConversionStep,
    trackLandingPageView,
    trackAnalysisIntent,
    trackAnalysisAttempt,
    trackAnalysisSuccess,
    trackAnalysisAbandoned,

    // State
    isEnabled: true, // Always enabled for MVP with console logging
  }
}
