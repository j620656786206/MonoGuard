'use client';

import React from 'react';
import { cn } from '@/lib/utils';

export interface Step {
  id: string;
  title: string;
  description?: string;
  status: 'completed' | 'current' | 'pending';
}

interface ProgressStepsProps {
  steps: Step[];
  className?: string;
}

export const ProgressSteps: React.FC<ProgressStepsProps> = ({ 
  steps, 
  className 
}) => {
  return (
    <nav className={cn('flex items-center justify-center', className)}>
      <ol className="flex items-center w-full max-w-4xl">
        {steps.map((step, stepIdx) => (
          <li 
            key={step.id} 
            className={cn(
              'flex items-center',
              stepIdx < steps.length - 1 ? 'flex-1' : ''
            )}
          >
            {/* Step Circle and Content */}
            <div className="flex flex-col items-center">
              <div className={cn(
                'flex items-center justify-center w-10 h-10 rounded-full border-2',
                {
                  'bg-blue-600 border-blue-600 text-white': step.status === 'completed',
                  'bg-white border-blue-600 text-blue-600': step.status === 'current',
                  'bg-white border-gray-300 text-gray-500': step.status === 'pending',
                }
              )}>
                {step.status === 'completed' ? (
                  <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                  </svg>
                ) : step.status === 'current' ? (
                  <div className="w-2 h-2 bg-blue-600 rounded-full animate-pulse" />
                ) : (
                  <span className="text-sm font-medium">{stepIdx + 1}</span>
                )}
              </div>
              
              <div className="mt-2 text-center">
                <div className={cn(
                  'text-sm font-medium',
                  {
                    'text-blue-600': step.status === 'completed' || step.status === 'current',
                    'text-gray-500': step.status === 'pending',
                  }
                )}>
                  {step.title}
                </div>
                {step.description && (
                  <div className="text-xs text-gray-500 mt-1 max-w-24">
                    {step.description}
                  </div>
                )}
              </div>
            </div>

            {/* Connector Line */}
            {stepIdx < steps.length - 1 && (
              <div className={cn(
                'flex-1 h-0.5 mx-4',
                {
                  'bg-blue-600': step.status === 'completed',
                  'bg-gray-300': step.status === 'current' || step.status === 'pending',
                }
              )} />
            )}
          </li>
        ))}
      </ol>
    </nav>
  );
};