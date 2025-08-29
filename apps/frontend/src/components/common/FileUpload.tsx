'use client';

import React, { useRef, ChangeEvent } from 'react';
import { FileProcessingResult } from '@mono-guard/shared-types';
import { useFileUpload } from '@/hooks/api/useFileUpload';
import { useDragAndDrop } from '@/hooks/ui/useDragAndDrop';
import { cn } from '@/lib/utils';

export interface FileUploadProps {
  onUploadComplete?: (result: FileProcessingResult) => void;
  onUploadError?: (errors: string[]) => void;
  className?: string;
  disabled?: boolean;
  accept?: string[];
  multiple?: boolean;
  maxFiles?: number;
}

export const FileUpload: React.FC<FileUploadProps> = ({
  onUploadComplete,
  onUploadError,
  className,
  disabled = false,
  accept = ['.zip', '.json'],
  multiple = true,
  maxFiles = 10,
}) => {
  const fileInputRef = useRef<HTMLInputElement>(null);
  
  const {
    isUploading,
    progress,
    result,
    errors,
    uploadFiles,
    reset,
    validateFiles,
  } = useFileUpload();

  const handleFileDrop = (files: File[]) => {
    if (disabled || isUploading) return;
    
    const filesToUpload = maxFiles ? files.slice(0, maxFiles) : files;
    handleFileUpload(filesToUpload);
  };

  const {
    isDragOver,
    isDragActive,
    onDragEnter,
    onDragOver,
    onDragLeave,
    onDrop,
  } = useDragAndDrop({
    onFileDrop: handleFileDrop,
    accept,
  });

  const handleFileSelect = (e: ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files || []);
    if (files.length > 0) {
      handleFileUpload(files);
    }
  };

  const handleFileUpload = async (files: File[]) => {
    try {
      await uploadFiles(files);
    } catch (error) {
      console.error('Upload failed:', error);
    }
  };

  const handleButtonClick = () => {
    if (!disabled && !isUploading && fileInputRef.current) {
      fileInputRef.current.click();
    }
  };

  const handleReset = () => {
    reset();
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  // Handle upload completion/error
  React.useEffect(() => {
    if (result && !isUploading) {
      onUploadComplete?.(result);
    }
  }, [result, isUploading, onUploadComplete]);

  React.useEffect(() => {
    if (errors.length > 0 && !isUploading) {
      onUploadError?.(errors);
    }
  }, [errors, isUploading, onUploadError]);

  const acceptString = accept.join(',');
  
  return (
    <div className={cn('w-full', className)}>
      {/* File Input */}
      <input
        ref={fileInputRef}
        type="file"
        accept={acceptString}
        multiple={multiple}
        onChange={handleFileSelect}
        className="hidden"
        disabled={disabled || isUploading}
      />

      {/* Drop Zone */}
      <div
        onDragEnter={onDragEnter}
        onDragOver={onDragOver}
        onDragLeave={onDragLeave}
        onDrop={onDrop}
        className={cn(
          'relative rounded-lg border-2 border-dashed transition-all duration-200 ease-in-out',
          'flex flex-col items-center justify-center p-8 text-center',
          'min-h-[200px] cursor-pointer hover:bg-gray-50',
          {
            // Default state
            'border-gray-300 bg-white text-gray-600': !isDragActive && !isDragOver && !disabled,
            
            // Drag states
            'border-blue-500 bg-blue-50 text-blue-600': isDragActive || isDragOver,
            
            // Disabled state
            'border-gray-200 bg-gray-50 text-gray-400 cursor-not-allowed': disabled,
            
            // Uploading state
            'border-blue-500 bg-blue-50': isUploading,
          }
        )}
        onClick={handleButtonClick}
      >
        {isUploading ? (
          <UploadProgress progress={progress} />
        ) : (
          <UploadPrompt 
            isDragActive={isDragActive || isDragOver}
            accept={accept}
            disabled={disabled}
          />
        )}
      </div>

      {/* Upload Results */}
      {result && !isUploading && (
        <UploadResults 
          result={result} 
          onReset={handleReset}
        />
      )}

      {/* Error Messages */}
      {errors.length > 0 && !isUploading && (
        <ErrorMessages 
          errors={errors} 
          onReset={handleReset}
        />
      )}
    </div>
  );
};

// Upload Progress Component
interface UploadProgressProps {
  progress: { percentage: number; loaded: number; total: number } | null;
}

const UploadProgress: React.FC<UploadProgressProps> = ({ progress }) => (
  <div className="w-full max-w-md">
    <div className="flex items-center justify-center mb-4">
      <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      <span className="ml-3 text-blue-600 font-medium">Uploading...</span>
    </div>
    
    {progress && (
      <div className="w-full">
        <div className="flex justify-between text-sm text-gray-600 mb-2">
          <span>Progress</span>
          <span>{progress.percentage}%</span>
        </div>
        <div className="w-full bg-gray-200 rounded-full h-2">
          <div 
            className="bg-blue-600 h-2 rounded-full transition-all duration-300"
            style={{ width: `${progress.percentage}%` }}
          />
        </div>
        <div className="text-xs text-gray-500 mt-2 text-center">
          {Math.round(progress.loaded / (1024 * 1024))} MB / {Math.round(progress.total / (1024 * 1024))} MB
        </div>
      </div>
    )}
  </div>
);

// Upload Prompt Component
interface UploadPromptProps {
  isDragActive: boolean;
  accept: string[];
  disabled: boolean;
}

const UploadPrompt: React.FC<UploadPromptProps> = ({ isDragActive, accept, disabled }) => (
  <>
    <div className="mb-4">
      <svg 
        className={cn("w-12 h-12 mx-auto", {
          'text-blue-500': isDragActive,
          'text-gray-400': disabled,
          'text-gray-500': !isDragActive && !disabled,
        })}
        fill="none" 
        stroke="currentColor" 
        viewBox="0 0 24 24"
      >
        <path 
          strokeLinecap="round" 
          strokeLinejoin="round" 
          strokeWidth={1.5} 
          d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" 
        />
      </svg>
    </div>

    <div className="space-y-2">
      <p className="text-lg font-medium">
        {isDragActive 
          ? 'Drop files here' 
          : disabled 
            ? 'Upload disabled'
            : 'Upload your files'
        }
      </p>
      
      {!disabled && (
        <>
          <p className="text-sm">
            Drag and drop files here, or{' '}
            <span className="text-blue-600 font-medium hover:text-blue-700">
              click to browse
            </span>
          </p>
          
          <p className="text-xs text-gray-500">
            Supported formats: {accept.join(', ')} • Max 50MB per file
          </p>
        </>
      )}
    </div>
  </>
);

// Upload Results Component
interface UploadResultsProps {
  result: FileProcessingResult;
  onReset: () => void;
}

const UploadResults: React.FC<UploadResultsProps> = ({ result, onReset }) => (
  <div className="mt-6 p-4 bg-green-50 border border-green-200 rounded-lg">
    <div className="flex items-center justify-between mb-3">
      <h4 className="font-medium text-green-800">Upload Successful</h4>
      <button
        onClick={onReset}
        className="text-sm text-green-600 hover:text-green-700 underline"
      >
        Upload More Files
      </button>
    </div>
    
    <div className="space-y-2 text-sm">
      <p className="text-green-700">
        Uploaded {result.files.length} file{result.files.length !== 1 ? 's' : ''}
      </p>
      
      {result.packageJsonFiles.length > 0 && (
        <p className="text-green-700">
          Found {result.packageJsonFiles.length} package.json file{result.packageJsonFiles.length !== 1 ? 's' : ''}
        </p>
      )}

      {result.errors && result.errors.length > 0 && (
        <div className="mt-3">
          <p className="text-amber-700 font-medium mb-1">Warnings:</p>
          <ul className="text-amber-600 text-xs space-y-1">
            {result.errors.map((error, index) => (
              <li key={index}>• {error}</li>
            ))}
          </ul>
        </div>
      )}
    </div>
  </div>
);

// Error Messages Component
interface ErrorMessagesProps {
  errors: string[];
  onReset: () => void;
}

const ErrorMessages: React.FC<ErrorMessagesProps> = ({ errors, onReset }) => (
  <div className="mt-6 p-4 bg-red-50 border border-red-200 rounded-lg">
    <div className="flex items-center justify-between mb-3">
      <h4 className="font-medium text-red-800">Upload Failed</h4>
      <button
        onClick={onReset}
        className="text-sm text-red-600 hover:text-red-700 underline"
      >
        Try Again
      </button>
    </div>
    
    <ul className="text-sm text-red-700 space-y-1">
      {errors.map((error, index) => (
        <li key={index}>• {error}</li>
      ))}
    </ul>
  </div>
);