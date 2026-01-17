'use client'

import type { FileProcessingResult } from '@monoguard/types'
import { useCallback, useState } from 'react'
import { type UploadProgress, UploadService } from '@/lib/api/services/upload'

export interface UseFileUploadState {
  isUploading: boolean
  progress: UploadProgress | null
  result: FileProcessingResult | null
  errors: string[]
}

export interface UseFileUploadActions {
  uploadFiles: (files: File[]) => Promise<void>
  reset: () => void
  validateFiles: (files: File[]) => { valid: boolean; errors: string[] }
}

export interface UseFileUploadReturn extends UseFileUploadState, UseFileUploadActions {}

export const useFileUpload = (): UseFileUploadReturn => {
  const [state, setState] = useState<UseFileUploadState>({
    isUploading: false,
    progress: null,
    result: null,
    errors: [],
  })

  const uploadFiles = useCallback(async (files: File[]) => {
    // Reset state
    setState({
      isUploading: true,
      progress: null,
      result: null,
      errors: [],
    })

    try {
      // Validate files first
      const validation = UploadService.validateFiles(files)
      if (!validation.valid) {
        setState((prev) => ({
          ...prev,
          isUploading: false,
          errors: validation.errors,
        }))
        return
      }

      // Upload files with progress tracking
      const result = await UploadService.uploadFiles(files, (progress) => {
        setState((prev) => ({
          ...prev,
          progress,
        }))
      })

      // Update state with result
      setState((prev) => ({
        ...prev,
        isUploading: false,
        result,
        progress: null,
        errors: [],
      }))
    } catch (error: any) {
      console.error('Upload error:', error)

      let errorMessage = 'Upload failed'
      if (error.response?.data?.message) {
        errorMessage = error.response.data.message
      } else if (error.message) {
        errorMessage = error.message
      }

      setState((prev) => ({
        ...prev,
        isUploading: false,
        progress: null,
        errors: [errorMessage],
      }))
    }
  }, [])

  const reset = useCallback(() => {
    setState({
      isUploading: false,
      progress: null,
      result: null,
      errors: [],
    })
  }, [])

  const validateFiles = useCallback((files: File[]) => {
    return UploadService.validateFiles(files)
  }, [])

  return {
    ...state,
    uploadFiles,
    reset,
    validateFiles,
  }
}
