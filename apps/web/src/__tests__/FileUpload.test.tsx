import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { FileUpload } from '@/components/common/FileUpload'
import { createFileProcessingResult, createUploadProgress } from './factories/upload.factory'

// Mock return values
const mockUploadFiles = vi.fn()
const mockReset = vi.fn()
const mockValidateFiles = vi.fn().mockReturnValue({ valid: true, errors: [] })

let mockUseFileUploadState = {
  isUploading: false,
  progress: null as { percentage: number; loaded: number; total: number } | null,
  result: null as ReturnType<typeof createFileProcessingResult> | null,
  errors: [] as string[],
  uploadFiles: mockUploadFiles,
  reset: mockReset,
  validateFiles: mockValidateFiles,
}

let mockUseDragAndDropState = {
  isDragOver: false,
  isDragActive: false,
  onDragEnter: vi.fn(),
  onDragOver: vi.fn(),
  onDragLeave: vi.fn(),
  onDrop: vi.fn(),
  reset: vi.fn(),
}

let capturedOnFileDrop: ((files: File[]) => void) | undefined

// Mock the hooks
vi.mock('@/hooks/api/useFileUpload', () => ({
  useFileUpload: () => mockUseFileUploadState,
}))

vi.mock('@/hooks/ui/useDragAndDrop', () => ({
  useDragAndDrop: (props: { onFileDrop: (files: File[]) => void; accept?: string[] }) => {
    capturedOnFileDrop = props.onFileDrop
    return mockUseDragAndDropState
  },
}))

describe('FileUpload', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    capturedOnFileDrop = undefined
    mockUseFileUploadState = {
      isUploading: false,
      progress: null,
      result: null,
      errors: [],
      uploadFiles: mockUploadFiles,
      reset: mockReset,
      validateFiles: mockValidateFiles,
    }
    mockUseDragAndDropState = {
      isDragOver: false,
      isDragActive: false,
      onDragEnter: vi.fn(),
      onDragOver: vi.fn(),
      onDragLeave: vi.fn(),
      onDrop: vi.fn(),
      reset: vi.fn(),
    }
  })

  describe('Default Rendering', () => {
    it('[P0] should render upload prompt with drag and drop instructions', () => {
      render(<FileUpload />)

      expect(screen.getByText('Upload your files')).toBeInTheDocument()
      expect(screen.getByText(/Drag and drop files here/)).toBeInTheDocument()
      expect(screen.getByText(/click to browse/)).toBeInTheDocument()
    })

    it('[P0] should display supported file formats', () => {
      render(<FileUpload />)

      expect(screen.getByText(/Supported formats:/)).toBeInTheDocument()
      expect(screen.getByText(/\.zip, \.json/)).toBeInTheDocument()
    })

    it('[P0] should render hidden file input', () => {
      render(<FileUpload />)

      const input = document.querySelector('input[type="file"]')
      expect(input).toBeInTheDocument()
      expect(input).toHaveClass('hidden')
    })
  })

  describe('File Selection', () => {
    it('[P0] should trigger file input click on drop zone click', () => {
      render(<FileUpload />)

      const dropZone = screen.getByText('Upload your files').closest('div')
      const fileInput = document.querySelector('input[type="file"]') as HTMLInputElement

      expect(dropZone).not.toBeNull()
      const clickSpy = vi.spyOn(fileInput, 'click')

      fireEvent.click(dropZone as HTMLElement)

      expect(clickSpy).toHaveBeenCalled()
    })

    it('[P0] should call uploadFiles when files are selected', async () => {
      render(<FileUpload />)

      const fileInput = document.querySelector('input[type="file"]') as HTMLInputElement
      const file = new File(['content'], 'test.zip', { type: 'application/zip' })

      Object.defineProperty(fileInput, 'files', {
        value: [file],
        writable: false,
      })
      fireEvent.change(fileInput)

      await waitFor(() => {
        expect(mockUploadFiles).toHaveBeenCalledWith([file])
      })
    })
  })

  describe('Drag and Drop Integration', () => {
    it('[P0] should display drag active state', () => {
      mockUseDragAndDropState.isDragActive = true
      mockUseDragAndDropState.isDragOver = true

      render(<FileUpload />)

      expect(screen.getByText('Drop files here')).toBeInTheDocument()
    })

    it('[P1] should call uploadFiles when files are dropped', () => {
      render(<FileUpload />)

      const files = [new File([''], 'test.zip', { type: 'application/zip' })]
      capturedOnFileDrop?.(files)

      expect(mockUploadFiles).toHaveBeenCalledWith(files)
    })
  })

  describe('Upload Progress', () => {
    it('[P0] should display upload progress when uploading', () => {
      mockUseFileUploadState.isUploading = true
      mockUseFileUploadState.progress = createUploadProgress({
        percentage: 50,
        loaded: 512 * 1024,
        total: 1024 * 1024,
      })

      render(<FileUpload />)

      expect(screen.getByText('Uploading...')).toBeInTheDocument()
      expect(screen.getByText('50%')).toBeInTheDocument()
    })

    it('[P0] should display loading spinner when uploading without progress', () => {
      mockUseFileUploadState.isUploading = true
      mockUseFileUploadState.progress = null

      render(<FileUpload />)

      expect(screen.getByText('Uploading...')).toBeInTheDocument()
    })
  })

  describe('Error Handling', () => {
    it('[P0] should display error messages when errors exist', () => {
      mockUseFileUploadState.errors = ['File type not allowed', 'File too large']

      render(<FileUpload />)

      expect(screen.getByText('Upload Failed')).toBeInTheDocument()
      expect(screen.getByText(/File type not allowed/)).toBeInTheDocument()
      expect(screen.getByText(/File too large/)).toBeInTheDocument()
    })

    it('[P0] should call reset when Try Again is clicked', () => {
      mockUseFileUploadState.errors = ['Upload failed']

      render(<FileUpload />)

      fireEvent.click(screen.getByText('Try Again'))

      expect(mockReset).toHaveBeenCalled()
    })

    it('[P0] should not display errors during upload', () => {
      mockUseFileUploadState.isUploading = true
      mockUseFileUploadState.errors = ['Previous error']

      render(<FileUpload />)

      expect(screen.queryByText('Upload Failed')).not.toBeInTheDocument()
    })
  })

  describe('Callbacks', () => {
    it('[P0] should call onUploadComplete when result is available', async () => {
      const onUploadComplete = vi.fn()
      const result = createFileProcessingResult()

      const { rerender } = render(<FileUpload onUploadComplete={onUploadComplete} />)

      mockUseFileUploadState.result = result
      rerender(<FileUpload onUploadComplete={onUploadComplete} />)

      await waitFor(() => {
        expect(onUploadComplete).toHaveBeenCalledWith(result)
      })
    })

    it('[P0] should call onUploadError when errors occur', async () => {
      const onUploadError = vi.fn()
      const errors = ['Upload failed']

      const { rerender } = render(<FileUpload onUploadError={onUploadError} />)

      mockUseFileUploadState.errors = errors
      rerender(<FileUpload onUploadError={onUploadError} />)

      await waitFor(() => {
        expect(onUploadError).toHaveBeenCalledWith(errors)
      })
    })
  })

  describe('Disabled State', () => {
    it('[P0] should display disabled message when disabled', () => {
      render(<FileUpload disabled />)

      expect(screen.getByText('Upload disabled')).toBeInTheDocument()
    })

    it('[P0] should not trigger file input when disabled', () => {
      render(<FileUpload disabled />)

      const dropZone = screen.getByText('Upload disabled').closest('div')
      const fileInput = document.querySelector('input[type="file"]') as HTMLInputElement

      expect(dropZone).not.toBeNull()
      const clickSpy = vi.spyOn(fileInput, 'click')

      fireEvent.click(dropZone as HTMLElement)

      expect(clickSpy).not.toHaveBeenCalled()
    })

    it('[P0] should disable file input element', () => {
      render(<FileUpload disabled />)

      const fileInput = document.querySelector('input[type="file"]')
      expect(fileInput).toBeDisabled()
    })

    it('[P1] should not handle file drop when disabled', () => {
      render(<FileUpload disabled />)

      const files = [new File([''], 'test.zip', { type: 'application/zip' })]
      capturedOnFileDrop?.(files)

      expect(mockUploadFiles).not.toHaveBeenCalled()
    })
  })

  describe('Props Configuration', () => {
    it('[P1] should accept custom file types', () => {
      render(<FileUpload accept={['.tar.gz', '.tgz']} />)

      expect(screen.getByText(/\.tar\.gz, \.tgz/)).toBeInTheDocument()
    })

    it('[P1] should set multiple attribute based on prop', () => {
      render(<FileUpload multiple={false} />)

      const fileInput = document.querySelector('input[type="file"]')
      expect(fileInput).not.toHaveAttribute('multiple')
    })

    it('[P1] should respect maxFiles when dropping', () => {
      render(<FileUpload maxFiles={2} />)

      const files = [
        new File([''], 'file1.zip', { type: 'application/zip' }),
        new File([''], 'file2.zip', { type: 'application/zip' }),
        new File([''], 'file3.zip', { type: 'application/zip' }),
      ]

      capturedOnFileDrop?.(files)

      expect(mockUploadFiles).toHaveBeenCalledWith(files.slice(0, 2))
    })

    it('[P1] should apply custom className', () => {
      const { container } = render(<FileUpload className="custom-upload-class" />)

      expect(container.firstChild).toHaveClass('custom-upload-class')
    })
  })

  describe('Upload State Blocking', () => {
    it('[P0] should not trigger new upload while uploading', () => {
      mockUseFileUploadState.isUploading = true

      render(<FileUpload />)

      const dropZone = screen.getByText('Uploading...').closest('div')
      const fileInput = document.querySelector('input[type="file"]') as HTMLInputElement

      expect(dropZone).not.toBeNull()
      const clickSpy = vi.spyOn(fileInput, 'click')

      fireEvent.click(dropZone as HTMLElement)

      expect(clickSpy).not.toHaveBeenCalled()
    })

    it('[P0] should disable file input while uploading', () => {
      mockUseFileUploadState.isUploading = true

      render(<FileUpload />)

      const fileInput = document.querySelector('input[type="file"]')
      expect(fileInput).toBeDisabled()
    })
  })
})
