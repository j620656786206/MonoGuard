import { act, renderHook, waitFor } from '@testing-library/react'
import { beforeEach, describe, expect, it, type Mock, vi } from 'vitest'
import { useFileUpload } from '../../app/hooks/api/useFileUpload'
import { UploadService } from '../../app/lib/api/services/upload'
import {
  createFileProcessingResult,
  createMockFile,
  createUploadProgress,
} from './factories/upload.factory'

// Mock UploadService with relative path matching import
vi.mock('../../app/lib/api/services/upload', () => ({
  UploadService: {
    validateFiles: vi.fn(),
    uploadFiles: vi.fn(),
  },
}))

describe('useFileUpload', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('Initial State', () => {
    it('[P0] should return initial state with all values set to defaults', () => {
      // GIVEN: Fresh hook instance
      const { result } = renderHook(() => useFileUpload())

      // THEN: Initial state should be correct
      expect(result.current.isUploading).toBe(false)
      expect(result.current.progress).toBeNull()
      expect(result.current.result).toBeNull()
      expect(result.current.errors).toEqual([])
    })

    it('[P0] should provide uploadFiles, reset, and validateFiles functions', () => {
      // GIVEN: Fresh hook instance
      const { result } = renderHook(() => useFileUpload())

      // THEN: Functions should be available
      expect(typeof result.current.uploadFiles).toBe('function')
      expect(typeof result.current.reset).toBe('function')
      expect(typeof result.current.validateFiles).toBe('function')
    })
  })

  describe('File Validation', () => {
    it('[P0] should validate files using UploadService', () => {
      // GIVEN: Mock validation result
      const mockValidation = { valid: true, errors: [] }
      ;(UploadService.validateFiles as Mock).mockReturnValue(mockValidation)

      const { result } = renderHook(() => useFileUpload())
      const files = [createMockFile()]

      // WHEN: validateFiles is called
      const validation = result.current.validateFiles(files)

      // THEN: UploadService.validateFiles should be called
      expect(UploadService.validateFiles).toHaveBeenCalledWith(files)
      expect(validation).toEqual(mockValidation)
    })

    it('[P0] should return validation errors for invalid files', () => {
      // GIVEN: Invalid file validation
      const mockValidation = {
        valid: false,
        errors: ['File type not allowed', 'File too large'],
      }
      ;(UploadService.validateFiles as Mock).mockReturnValue(mockValidation)

      const { result } = renderHook(() => useFileUpload())
      const files = [createMockFile({ name: 'invalid.exe', type: 'application/x-msdownload' })]

      // WHEN: validateFiles is called
      const validation = result.current.validateFiles(files)

      // THEN: Should return validation errors
      expect(validation.valid).toBe(false)
      expect(validation.errors).toHaveLength(2)
    })
  })

  describe('Upload Flow - Success', () => {
    it('[P0] should set isUploading to true when upload starts', async () => {
      // GIVEN: Valid files and successful validation
      ;(UploadService.validateFiles as Mock).mockReturnValue({ valid: true, errors: [] })

      // Track whether isUploading was true during the upload
      let wasUploadingDuringCall = false
      let resolveUpload: ((value: unknown) => void) | undefined

      // Create a promise we can manually resolve
      const uploadPromise = new Promise((resolve) => {
        resolveUpload = resolve
      })

      ;(UploadService.uploadFiles as Mock).mockImplementation(async () => {
        return uploadPromise
      })

      const { result } = renderHook(() => useFileUpload())
      const files = [createMockFile()]

      // WHEN: Upload is triggered
      act(() => {
        result.current.uploadFiles(files)
      })

      // THEN: isUploading should be true while upload is pending
      await waitFor(() => {
        expect(result.current.isUploading).toBe(true)
      })
      wasUploadingDuringCall = result.current.isUploading

      // Complete the upload
      await act(async () => {
        resolveUpload?.(createFileProcessingResult())
        await uploadPromise
      })

      expect(wasUploadingDuringCall).toBe(true)
    })

    it('[P0] should track upload progress via callback', async () => {
      // GIVEN: Valid files and upload with progress
      ;(UploadService.validateFiles as Mock).mockReturnValue({ valid: true, errors: [] })

      let progressCallback:
        | ((progress: { loaded: number; total: number; percentage: number }) => void)
        | undefined
      ;(UploadService.uploadFiles as Mock).mockImplementation(async (_files, onProgress) => {
        progressCallback = onProgress
        // Simulate progress
        onProgress?.(createUploadProgress({ loaded: 512, total: 1024, percentage: 50 }))
        return createFileProcessingResult()
      })

      const { result } = renderHook(() => useFileUpload())
      const files = [createMockFile()]

      // WHEN: Upload is triggered
      await act(async () => {
        await result.current.uploadFiles(files)
      })

      // THEN: Progress callback should have been captured
      expect(progressCallback).toBeDefined()
    })

    it('[P0] should set result and clear isUploading on success', async () => {
      // GIVEN: Valid files and successful upload
      const mockResult = createFileProcessingResult()
      ;(UploadService.validateFiles as Mock).mockReturnValue({ valid: true, errors: [] })
      ;(UploadService.uploadFiles as Mock).mockResolvedValue(mockResult)

      const { result } = renderHook(() => useFileUpload())
      const files = [createMockFile()]

      // WHEN: Upload completes
      await act(async () => {
        await result.current.uploadFiles(files)
      })

      // THEN: State should reflect success
      expect(result.current.isUploading).toBe(false)
      expect(result.current.result).toEqual(mockResult)
      expect(result.current.errors).toEqual([])
      expect(result.current.progress).toBeNull()
    })
  })

  describe('Upload Flow - Validation Failure', () => {
    it('[P0] should set errors and not upload when validation fails', async () => {
      // GIVEN: Invalid files
      const validationErrors = ['File 1 (invalid.exe): File type not allowed']
      ;(UploadService.validateFiles as Mock).mockReturnValue({
        valid: false,
        errors: validationErrors,
      })

      const { result } = renderHook(() => useFileUpload())
      const files = [createMockFile({ name: 'invalid.exe' })]

      // WHEN: Upload is attempted
      await act(async () => {
        await result.current.uploadFiles(files)
      })

      // THEN: Should have errors and not call upload
      expect(result.current.errors).toEqual(validationErrors)
      expect(result.current.isUploading).toBe(false)
      expect(UploadService.uploadFiles).not.toHaveBeenCalled()
    })
  })

  describe('Upload Flow - Upload Failure', () => {
    it('[P0] should handle upload errors with response data', async () => {
      // GIVEN: Valid files but upload fails
      ;(UploadService.validateFiles as Mock).mockReturnValue({ valid: true, errors: [] })
      ;(UploadService.uploadFiles as Mock).mockRejectedValue({
        response: { data: { message: 'Server error: File too large' } },
      })

      const { result } = renderHook(() => useFileUpload())
      const files = [createMockFile()]

      // WHEN: Upload fails
      await act(async () => {
        await result.current.uploadFiles(files)
      })

      // THEN: Should set error from response
      expect(result.current.errors).toEqual(['Server error: File too large'])
      expect(result.current.isUploading).toBe(false)
      expect(result.current.result).toBeNull()
    })

    it('[P0] should handle upload errors with error message', async () => {
      // GIVEN: Valid files but upload fails with Error
      ;(UploadService.validateFiles as Mock).mockReturnValue({ valid: true, errors: [] })
      ;(UploadService.uploadFiles as Mock).mockRejectedValue(new Error('Network error'))

      const { result } = renderHook(() => useFileUpload())
      const files = [createMockFile()]

      // WHEN: Upload fails
      await act(async () => {
        await result.current.uploadFiles(files)
      })

      // THEN: Should set error from Error.message
      expect(result.current.errors).toEqual(['Network error'])
      expect(result.current.isUploading).toBe(false)
    })

    it('[P0] should handle upload errors with unknown error type', async () => {
      // GIVEN: Valid files but upload fails with unknown error
      ;(UploadService.validateFiles as Mock).mockReturnValue({ valid: true, errors: [] })
      ;(UploadService.uploadFiles as Mock).mockRejectedValue('Unknown error')

      const { result } = renderHook(() => useFileUpload())
      const files = [createMockFile()]

      // WHEN: Upload fails
      await act(async () => {
        await result.current.uploadFiles(files)
      })

      // THEN: Should set default error message
      expect(result.current.errors).toEqual(['Upload failed'])
      expect(result.current.isUploading).toBe(false)
    })
  })

  describe('Reset Functionality', () => {
    it('[P0] should reset all state to initial values', async () => {
      // GIVEN: Hook with upload result
      const mockResult = createFileProcessingResult()
      ;(UploadService.validateFiles as Mock).mockReturnValue({ valid: true, errors: [] })
      ;(UploadService.uploadFiles as Mock).mockResolvedValue(mockResult)

      const { result } = renderHook(() => useFileUpload())

      await act(async () => {
        await result.current.uploadFiles([createMockFile()])
      })

      expect(result.current.result).not.toBeNull()

      // WHEN: Reset is called
      act(() => {
        result.current.reset()
      })

      // THEN: All state should be reset
      expect(result.current.isUploading).toBe(false)
      expect(result.current.progress).toBeNull()
      expect(result.current.result).toBeNull()
      expect(result.current.errors).toEqual([])
    })

    it('[P0] should reset errors after failed upload', async () => {
      // GIVEN: Hook with validation errors
      ;(UploadService.validateFiles as Mock).mockReturnValue({
        valid: false,
        errors: ['Invalid file'],
      })

      const { result } = renderHook(() => useFileUpload())

      await act(async () => {
        await result.current.uploadFiles([createMockFile()])
      })

      expect(result.current.errors).toHaveLength(1)

      // WHEN: Reset is called
      act(() => {
        result.current.reset()
      })

      // THEN: Errors should be cleared
      expect(result.current.errors).toEqual([])
    })
  })

  describe('State Consistency', () => {
    it('[P1] should reset previous result when starting new upload', async () => {
      // GIVEN: Hook with previous result
      const firstResult = createFileProcessingResult({ uploadId: 'first-upload' })
      const secondResult = createFileProcessingResult({ uploadId: 'second-upload' })
      ;(UploadService.validateFiles as Mock).mockReturnValue({ valid: true, errors: [] })
      ;(UploadService.uploadFiles as Mock)
        .mockResolvedValueOnce(firstResult)
        .mockResolvedValueOnce(secondResult)

      const { result } = renderHook(() => useFileUpload())

      // First upload
      await act(async () => {
        await result.current.uploadFiles([createMockFile()])
      })
      expect(result.current.result?.uploadId).toBe('first-upload')

      // WHEN: Second upload starts
      await act(async () => {
        await result.current.uploadFiles([createMockFile()])
      })

      // THEN: Should have second result
      expect(result.current.result?.uploadId).toBe('second-upload')
    })

    it('[P1] should clear previous errors when starting new upload', async () => {
      // GIVEN: Hook with previous errors
      ;(UploadService.validateFiles as Mock)
        .mockReturnValueOnce({ valid: false, errors: ['Previous error'] })
        .mockReturnValueOnce({ valid: true, errors: [] })
      ;(UploadService.uploadFiles as Mock).mockResolvedValue(createFileProcessingResult())

      const { result } = renderHook(() => useFileUpload())

      // First upload fails validation
      await act(async () => {
        await result.current.uploadFiles([createMockFile()])
      })
      expect(result.current.errors).toEqual(['Previous error'])

      // WHEN: Second upload succeeds
      await act(async () => {
        await result.current.uploadFiles([createMockFile()])
      })

      // THEN: Previous errors should be cleared
      expect(result.current.errors).toEqual([])
    })
  })
})
