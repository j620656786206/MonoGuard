import { apiClient } from '../client';

export interface EventProperties {
  url?: string;
  referrer?: string;
  title?: string;
  projectId?: string;
  analysisType?: string;
  fileSize?: number;
  fileName?: string;
  errorMessage?: string;
  errorCode?: string;
  repositoryUrl?: string;
  branch?: string;
  healthScore?: number;
  totalDependencies?: number;
  vulnerabilityCount?: number;
  customData?: Record<string, any>;
}

export interface FunnelStepMetadata {
  source: string;
  context: string;
  abortReason?: string;
  retryCount: number;
  customData?: Record<string, any>;
}

export type EventType = 
  | 'page_view'
  | 'analysis_start' 
  | 'analysis_complete'
  | 'upload_start'
  | 'upload_complete'
  | 'github_analysis'
  | 'download'
  | 'error'
  | 'click'
  | 'form_submit'
  | 'feature_view';

export class AnalyticsService {
  // Track a generic event
  static async trackEvent(
    eventType: EventType,
    eventName: string,
    properties?: EventProperties,
    options?: {
      page?: string;
      element?: string;
      value?: string;
      duration?: number;
    }
  ): Promise<void> {
    try {
      await apiClient.post('/api/v1/analytics/events', {
        eventType,
        eventName,
        page: options?.page,
        element: options?.element,
        value: options?.value,
        duration: options?.duration,
        properties: {
          url: window.location.href,
          title: document.title,
          referrer: document.referrer,
          ...properties,
        },
      });
    } catch (error) {
      console.error('Failed to track event:', error);
      // Don't throw - analytics failures shouldn't break the app
    }
  }

  // Track page views
  static async trackPageView(url?: string, title?: string): Promise<void> {
    try {
      await apiClient.post('/api/v1/analytics/pageview', {
        url: url || window.location.href,
        title: title || document.title,
        referrer: document.referrer,
      });
    } catch (error) {
      console.error('Failed to track page view:', error);
    }
  }

  // Track analysis events
  static async trackAnalysisStart(analysisType: string, properties?: EventProperties): Promise<void> {
    return this.trackAnalysis('start', analysisType, { properties });
  }

  static async trackAnalysisComplete(
    analysisType: string,
    duration: number,
    healthScore: number,
    properties?: EventProperties
  ): Promise<void> {
    return this.trackAnalysis('complete', analysisType, {
      duration,
      healthScore,
      properties,
    });
  }

  private static async trackAnalysis(
    type: 'start' | 'complete',
    analysisType: string,
    options: {
      duration?: number;
      healthScore?: number;
      properties?: EventProperties;
    }
  ): Promise<void> {
    try {
      await apiClient.post('/api/v1/analytics/analysis', {
        type,
        analysisType,
        duration: options.duration,
        healthScore: options.healthScore,
        properties: options.properties,
      });
    } catch (error) {
      console.error('Failed to track analysis event:', error);
    }
  }

  // Track conversion funnel steps
  static async trackConversionStep(
    step: string,
    stepOrder: number,
    completed: boolean,
    duration?: number,
    metadata?: FunnelStepMetadata
  ): Promise<void> {
    try {
      await apiClient.post('/api/v1/analytics/conversion', {
        step,
        stepOrder,
        completed,
        duration,
        metadata,
      });
    } catch (error) {
      console.error('Failed to track conversion step:', error);
    }
  }

  // Track errors
  static async trackError(errorMessage: string, errorCode?: string, page?: string): Promise<void> {
    try {
      await apiClient.post('/api/v1/analytics/error', {
        errorMessage,
        errorCode,
        page: page || window.location.href,
      });
    } catch (error) {
      console.error('Failed to track error:', error);
    }
  }

  // Track user interactions
  static async trackClick(element: string, value?: string, properties?: EventProperties): Promise<void> {
    return this.trackEvent('click', 'click', properties, {
      element,
      value,
      page: window.location.pathname,
    });
  }

  static async trackFormSubmit(formName: string, properties?: EventProperties): Promise<void> {
    return this.trackEvent('form_submit', 'form_submit', properties, {
      element: formName,
      page: window.location.pathname,
    });
  }

  // Track feature usage
  static async trackFeatureView(featureName: string, properties?: EventProperties): Promise<void> {
    return this.trackEvent('feature_view', 'feature_view', properties, {
      element: featureName,
      page: window.location.pathname,
    });
  }

  // Track GitHub analysis specifically
  static async trackGitHubAnalysis(
    repositoryUrl: string,
    branch: string,
    properties?: EventProperties
  ): Promise<void> {
    return this.trackEvent('github_analysis', 'github_analysis', {
      repositoryUrl,
      branch,
      ...properties,
    });
  }

  // Track file uploads
  static async trackUploadStart(fileName: string, fileSize: number): Promise<void> {
    return this.trackEvent('upload_start', 'upload_start', {
      fileName,
      fileSize,
    });
  }

  static async trackUploadComplete(
    fileName: string,
    fileSize: number,
    duration: number
  ): Promise<void> {
    return this.trackEvent('upload_complete', 'upload_complete', {
      fileName,
      fileSize,
    }, {
      duration,
    });
  }

  // Track downloads
  static async trackDownload(fileName: string, properties?: EventProperties): Promise<void> {
    return this.trackEvent('download', 'download', {
      fileName,
      ...properties,
    });
  }

  // Conversion funnel tracking helpers
  static async trackLandingPageView(): Promise<void> {
    await this.trackPageView();
    return this.trackConversionStep('landing_page_view', 1, true);
  }

  static async trackAnalysisIntent(): Promise<void> {
    return this.trackConversionStep('analysis_intent', 2, true);
  }

  static async trackAnalysisAttempt(): Promise<void> {
    return this.trackConversionStep('analysis_attempt', 3, true);
  }

  static async trackAnalysisSuccess(): Promise<void> {
    return this.trackConversionStep('analysis_success', 4, true);
  }

  static async trackAnalysisAbandoned(reason?: string): Promise<void> {
    return this.trackConversionStep('analysis_attempt', 3, false, undefined, {
      source: window.location.pathname,
      context: 'user_abandoned',
      abortReason: reason,
      retryCount: 0,
    });
  }
}