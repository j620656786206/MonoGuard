'use client'

import type { FileProcessingResult } from '@monoguard/types'
import React, { type ChangeEvent, useRef } from 'react'
import { useFileUpload } from '@/hooks/api/useFileUpload'
import { useDragAndDrop } from '@/hooks/ui/useDragAndDrop'
import { cn } from '@/lib/utils'

export interface FileUploadProps {
  onUploadComplete?: (result: FileProcessingResult) => void
  onUploadError?: (errors: string[]) => void
  className?: string
  disabled?: boolean
  accept?: string[]
  multiple?: boolean
  maxFiles?: number
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
  const fileInputRef = useRef<HTMLInputElement>(null)

  const { isUploading, progress, result, errors, uploadFiles, reset, validateFiles } =
    useFileUpload()

  const handleFileDrop = (files: File[]) => {
    if (disabled || isUploading) return

    const filesToUpload = maxFiles ? files.slice(0, maxFiles) : files
    handleFileUpload(filesToUpload)
  }

  const { isDragOver, isDragActive, onDragEnter, onDragOver, onDragLeave, onDrop } = useDragAndDrop(
    {
      onFileDrop: handleFileDrop,
      accept,
    }
  )

  const handleFileSelect = (e: ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files || [])
    if (files.length > 0) {
      handleFileUpload(files)
    }
  }

  const handleFileUpload = async (files: File[]) => {
    try {
      await uploadFiles(files)
    } catch (error) {
      console.error('Upload failed:', error)
    }
  }

  const handleButtonClick = () => {
    if (!disabled && !isUploading && fileInputRef.current) {
      fileInputRef.current.click()
    }
  }

  const handleReset = () => {
    reset()
    if (fileInputRef.current) {
      fileInputRef.current.value = ''
    }
  }

  // Handle upload completion/error
  React.useEffect(() => {
    if (result && !isUploading) {
      onUploadComplete?.(result)
    }
  }, [result, isUploading])

  React.useEffect(() => {
    if (errors.length > 0 && !isUploading) {
      onUploadError?.(errors)
    }
  }, [errors, isUploading])

  const acceptString = accept.join(',')

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
            'cursor-not-allowed border-gray-200 bg-gray-50 text-gray-400': disabled,

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

      {/* Error Messages */}
      {errors.length > 0 && !isUploading && <ErrorMessages errors={errors} onReset={handleReset} />}
    </div>
  )
}

// Upload Progress Component
interface UploadProgressProps {
  progress: { percentage: number; loaded: number; total: number } | null
}

const UploadProgress: React.FC<UploadProgressProps> = ({ progress }) => (
  <div className="w-full max-w-md">
    <div className="mb-4 flex items-center justify-center">
      <div className="h-8 w-8 animate-spin rounded-full border-b-2 border-blue-600"></div>
      <span className="ml-3 font-medium text-blue-600">Uploading...</span>
    </div>

    {progress && (
      <div className="w-full">
        <div className="mb-2 flex justify-between text-sm text-gray-600">
          <span>Progress</span>
          <span>{progress.percentage}%</span>
        </div>
        <div className="h-2 w-full rounded-full bg-gray-200">
          <div
            className="h-2 rounded-full bg-blue-600 transition-all duration-300"
            style={{ width: `${progress.percentage}%` }}
          />
        </div>
        <div className="mt-2 text-center text-xs text-gray-500">
          {Math.round(progress.loaded / (1024 * 1024))} MB /{' '}
          {Math.round(progress.total / (1024 * 1024))} MB
        </div>
      </div>
    )}
  </div>
)

// Upload Prompt Component
interface UploadPromptProps {
  isDragActive: boolean
  accept: string[]
  disabled: boolean
}

const UploadPrompt: React.FC<UploadPromptProps> = ({ isDragActive, accept, disabled }) => (
  <>
    <div className="mb-4">
      <svg
        className={cn('mx-auto h-12 w-12', {
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
        {isDragActive ? 'Drop files here' : disabled ? 'Upload disabled' : 'Upload your files'}
      </p>

      {!disabled && (
        <>
          <p className="text-sm">
            Drag and drop files here, or{' '}
            <span className="font-medium text-blue-600 hover:text-blue-700">click to browse</span>
          </p>

          <p className="text-xs text-gray-500">
            Supported formats: {accept.join(', ')} • Max 50MB per file
          </p>
        </>
      )}
    </div>
  </>
)

// Error Messages Component
interface ErrorMessagesProps {
  errors: string[]
  onReset: () => void
}

const ErrorMessages: React.FC<ErrorMessagesProps> = ({ errors, onReset }) => (
  <div className="mt-6 rounded-lg border border-red-200 bg-red-50 p-4">
    <div className="mb-3 flex items-center justify-between">
      <h4 className="font-medium text-red-800">Upload Failed</h4>
      <button onClick={onReset} className="text-sm text-red-600 underline hover:text-red-700">
        Try Again
      </button>
    </div>

    <ul className="space-y-1 text-sm text-red-700">
      {errors.map((error, index) => (
        <li key={index}>• {error}</li>
      ))}
    </ul>
  </div>
)
