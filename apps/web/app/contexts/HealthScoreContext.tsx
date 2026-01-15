'use client';
import React, { createContext, useContext, useState, ReactNode } from 'react';

interface HealthScoreContextType {
  healthScore: number;
  setHealthScore: (score: number) => void;
  isLoading: boolean;
  setIsLoading: (loading: boolean) => void;
}

const HealthScoreContext = createContext<HealthScoreContextType | undefined>(
  undefined
);

interface HealthScoreProviderProps {
  children: ReactNode;
}

export function HealthScoreProvider({ children }: HealthScoreProviderProps) {
  const [healthScore, setHealthScore] = useState<number>(0);
  const [isLoading, setIsLoading] = useState<boolean>(true);

  return (
    <HealthScoreContext.Provider
      value={{
        healthScore,
        setHealthScore,
        isLoading,
        setIsLoading,
      }}
    >
      {children}
    </HealthScoreContext.Provider>
  );
}

export function useHealthScore() {
  const context = useContext(HealthScoreContext);
  if (context === undefined) {
    throw new Error('useHealthScore must be used within a HealthScoreProvider');
  }
  return context;
}
