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
  className,
}) => {
  return (
    <nav className={cn('flex items-center justify-center', className)}>
      <ol className="flex w-full max-w-4xl items-center">
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
              <div
                className={cn(
                  'flex h-10 w-10 items-center justify-center rounded-full border-2',
                  {
                    'border-blue-600 bg-blue-600 text-white':
                      step.status === 'completed',
                    'border-blue-600 bg-white text-blue-600':
                      step.status === 'current',
                    'border-gray-300 bg-white text-gray-500':
                      step.status === 'pending',
                  }
                )}
              >
                {step.status === 'completed' ? (
                  <svg
                    className="h-5 w-5"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path
                      fillRule="evenodd"
                      d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                      clipRule="evenodd"
                    />
                  </svg>
                ) : step.status === 'current' ? (
                  <div className="h-2 w-2 animate-pulse rounded-full bg-blue-600" />
                ) : (
                  <span className="text-sm font-medium">{stepIdx + 1}</span>
                )}
              </div>

              <div className="mt-2 text-center">
                <div
                  className={cn('text-sm font-medium', {
                    'text-blue-600':
                      step.status === 'completed' || step.status === 'current',
                    'text-gray-500': step.status === 'pending',
                  })}
                >
                  {step.title}
                </div>
                {step.description && (
                  <div className="mt-1 max-w-24 text-xs text-gray-500">
                    {step.description}
                  </div>
                )}
              </div>
            </div>

            {/* Connector Line */}
            {stepIdx < steps.length - 1 && (
              <div
                className={cn('mx-4 h-0.5 flex-1', {
                  'bg-blue-600': step.status === 'completed',
                  'bg-gray-300':
                    step.status === 'current' || step.status === 'pending',
                })}
              />
            )}
          </li>
        ))}
      </ol>
    </nav>
  );
};
