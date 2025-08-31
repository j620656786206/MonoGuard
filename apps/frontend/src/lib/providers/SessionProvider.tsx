'use client';

import React, { createContext, useContext, useEffect, useState, ReactNode } from 'react';
import { apiClient } from '../api/client';

interface SessionMetadata {
  browser: string;
  os: string;
  device: string;
  timezone: string;
  language: string;
  screenResolution: string;
  referrerUrl?: string;
  utmSource?: string;
  utmCampaign?: string;
  utmMedium?: string;
  customProperties?: Record<string, any>;
}

interface AnonymousSession {
  id: string;
  sessionToken: string;
  ipAddress: string;
  userAgent: string;
  country?: string;
  city?: string;
  status: 'active' | 'expired' | 'revoked';
  lastActivityAt: string;
  expiresAt: string;
  metadata?: SessionMetadata;
  createdAt: string;
  updatedAt: string;
}

interface SessionContextType {
  session: AnonymousSession | null;
  isLoading: boolean;
  error: string | null;
  refreshSession: () => Promise<void>;
  createNewSession: () => Promise<void>;
  revokeSession: () => Promise<void>;
}

const SessionContext = createContext<SessionContextType | undefined>(undefined);

const SESSION_TOKEN_KEY = 'monoguard_session_token';
const SESSION_STORAGE_KEY = 'monoguard_session';

interface SessionProviderProps {
  children: ReactNode;
}

export function SessionProvider({ children }: SessionProviderProps) {
  const [session, setSession] = useState<AnonymousSession | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Get browser metadata
  const getBrowserMetadata = (): SessionMetadata => {
    const userAgent = navigator.userAgent;
    const language = navigator.language || 'en-US';
    const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
    const screenResolution = `${screen.width}x${screen.height}`;
    
    // Parse browser
    let browser = 'Unknown';
    if (userAgent.includes('Chrome') && !userAgent.includes('Edg')) {
      browser = 'Chrome';
    } else if (userAgent.includes('Firefox')) {
      browser = 'Firefox';
    } else if (userAgent.includes('Safari') && !userAgent.includes('Chrome')) {
      browser = 'Safari';
    } else if (userAgent.includes('Edg')) {
      browser = 'Edge';
    }

    // Parse OS
    let os = 'Unknown';
    if (userAgent.includes('Windows')) {
      os = 'Windows';
    } else if (userAgent.includes('Mac')) {
      os = 'macOS';
    } else if (userAgent.includes('Linux')) {
      os = 'Linux';
    } else if (userAgent.includes('Android')) {
      os = 'Android';
    } else if (userAgent.includes('iPhone') || userAgent.includes('iPad')) {
      os = 'iOS';
    }

    // Parse device type
    let device = 'Desktop';
    if (/Mobile|Android|iPhone/.test(userAgent)) {
      device = 'Mobile';
    } else if (/Tablet|iPad/.test(userAgent)) {
      device = 'Tablet';
    }

    // Get UTM parameters from URL
    const urlParams = new URLSearchParams(window.location.search);
    const utmSource = urlParams.get('utm_source') || undefined;
    const utmCampaign = urlParams.get('utm_campaign') || undefined;
    const utmMedium = urlParams.get('utm_medium') || undefined;

    return {
      browser,
      os,
      device,
      timezone,
      language,
      screenResolution,
      referrerUrl: document.referrer || undefined,
      utmSource,
      utmCampaign,
      utmMedium,
    };
  };

  // Store session token in localStorage
  const storeSessionToken = (token: string) => {
    localStorage.setItem(SESSION_TOKEN_KEY, token);
  };

  // Get stored session token
  const getStoredSessionToken = (): string | null => {
    return localStorage.getItem(SESSION_TOKEN_KEY);
  };

  // Remove stored session token
  const removeStoredSessionToken = () => {
    localStorage.removeItem(SESSION_TOKEN_KEY);
    localStorage.removeItem(SESSION_STORAGE_KEY);
  };

  // Set session token in API client headers
  const setApiClientSessionToken = (token: string) => {
    // Check if apiClient has the expected structure
    if (apiClient && typeof apiClient.setSessionToken === 'function') {
      apiClient.setSessionToken(token);
    } else {
      console.warn('Unable to set session token in API client');
    }
  };

  // Remove session token from API client headers
  const removeApiClientSessionToken = () => {
    if (apiClient && typeof apiClient.clearSessionToken === 'function') {
      apiClient.clearSessionToken();
    } else {
      console.warn('Unable to clear session token in API client');
    }
  };

  // Create a new session
  const createNewSession = async (): Promise<void> => {
    try {
      setIsLoading(true);
      setError(null);

      const metadata = getBrowserMetadata();
      
      // Create session via API
      const response = await apiClient.post('/api/v1/sessions', { metadata });
      
      // The API client returns response.data directly, so response is the data
      if (!response) {
        console.error('No response received from session API');
        throw new Error('No response from server');
      }
      
      const sessionToken = response.session?.sessionToken;
      
      if (!sessionToken) {
        console.error('Session response structure:', response);
        console.error('Expected session.sessionToken but got:', typeof response.session);
        throw new Error('No session token received from server');
      }

      // Store token and set in API client
      storeSessionToken(sessionToken);
      setApiClientSessionToken(sessionToken);
      
      // Store session data
      const sessionData = response.session;
      setSession(sessionData);
      localStorage.setItem(SESSION_STORAGE_KEY, JSON.stringify(sessionData));
      
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || err.message || 'Failed to create session';
      setError(errorMessage);
      console.error('Session creation failed:', err);
    } finally {
      setIsLoading(false);
    }
  };

  // Validate existing session
  const validateSession = async (token: string): Promise<AnonymousSession | null> => {
    try {
      setApiClientSessionToken(token);
      const response = await apiClient.get('/api/v1/sessions/current');
      return response.session;
    } catch (err: any) {
      console.error('Session validation failed:', err);
      // Remove invalid token
      removeStoredSessionToken();
      removeApiClientSessionToken();
      return null;
    }
  };

  // Refresh session
  const refreshSession = async (): Promise<void> => {
    const token = getStoredSessionToken();
    if (!token) {
      await createNewSession();
      return;
    }

    try {
      setError(null);
      setApiClientSessionToken(token);
      const response = await apiClient.put('/api/v1/sessions/current/refresh');
      const sessionData = response.session;
      setSession(sessionData);
      localStorage.setItem(SESSION_STORAGE_KEY, JSON.stringify(sessionData));
    } catch (err: any) {
      console.error('Session refresh failed:', err);
      // Create new session if refresh fails
      await createNewSession();
    }
  };

  // Revoke session
  const revokeSession = async (): Promise<void> => {
    const token = getStoredSessionToken();
    if (!token) return;

    try {
      setApiClientSessionToken(token);
      await apiClient.delete('/api/v1/sessions/current');
    } catch (err) {
      console.error('Session revocation failed:', err);
    } finally {
      // Clean up regardless of API call success
      removeStoredSessionToken();
      removeApiClientSessionToken();
      setSession(null);
    }
  };

  // Initialize session on mount
  useEffect(() => {
    const initializeSession = async () => {
      setIsLoading(true);
      
      try {
        const storedToken = getStoredSessionToken();
        
        if (storedToken) {
          // Try to validate existing session
          const validSession = await validateSession(storedToken);
          if (validSession) {
            setSession(validSession);
            localStorage.setItem(SESSION_STORAGE_KEY, JSON.stringify(validSession));
            return;
          }
        }

        // Create new session if no valid session exists
        await createNewSession();
        
      } catch (err) {
        console.error('Session initialization failed:', err);
        setError('Failed to initialize session');
      } finally {
        setIsLoading(false);
      }
    };

    initializeSession();
  }, []);

  // Auto-refresh session before expiration
  useEffect(() => {
    if (!session || session.status !== 'active') return;

    const expiresAt = new Date(session.expiresAt).getTime();
    const now = Date.now();
    const timeUntilExpiry = expiresAt - now;
    
    // Refresh 5 minutes before expiry or immediately if already expired
    const refreshTime = Math.max(0, timeUntilExpiry - (5 * 60 * 1000));
    
    // Only set timeout if refresh time is reasonable (more than 1 minute from now)
    if (refreshTime > 60 * 1000) {
      const timeoutId = setTimeout(() => {
        refreshSession();
      }, refreshTime);

      return () => clearTimeout(timeoutId);
    }
  }, [session]);

  const value: SessionContextType = {
    session,
    isLoading,
    error,
    refreshSession,
    createNewSession,
    revokeSession,
  };

  return (
    <SessionContext.Provider value={value}>
      {children}
    </SessionContext.Provider>
  );
}

export function useSession(): SessionContextType {
  const context = useContext(SessionContext);
  if (context === undefined) {
    throw new Error('useSession must be used within a SessionProvider');
  }
  return context;
}