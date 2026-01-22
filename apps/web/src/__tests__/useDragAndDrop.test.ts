import { act, renderHook } from '@testing-library/react'
import type { DragEvent } from 'react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { useDragAndDrop } from '../../app/hooks/ui/useDragAndDrop'
import { createMockFile } from './factories/upload.factory'

/**
 * Creates a mock DragEvent for testing
 */
function createMockDragEvent(
  files: File[] = [],
  overrides: Partial<DragEvent<HTMLElement>> = {}
): DragEvent<HTMLElement> {
  const items = files.map((file) => ({
    kind: 'file' as const,
    getAsFile: () => file,
  }))

  return {
    preventDefault: vi.fn(),
    stopPropagation: vi.fn(),
    dataTransfer: {
      items: items as unknown as DataTransferItemList,
      files: files as unknown as FileList,
    },
    currentTarget: document.createElement('div'),
    relatedTarget: null,
    ...overrides,
  } as unknown as DragEvent<HTMLElement>
}

describe('useDragAndDrop', () => {
  const mockOnFileDrop = vi.fn()

  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('Initial State', () => {
    it('[P1] should return initial state with drag states set to false', () => {
      // GIVEN: Fresh hook instance
      const { result } = renderHook(() => useDragAndDrop({ onFileDrop: mockOnFileDrop }))

      // THEN: Initial state should be correct
      expect(result.current.isDragOver).toBe(false)
      expect(result.current.isDragActive).toBe(false)
    })

    it('[P1] should provide all event handler functions', () => {
      // GIVEN: Fresh hook instance
      const { result } = renderHook(() => useDragAndDrop({ onFileDrop: mockOnFileDrop }))

      // THEN: Event handlers should be available
      expect(typeof result.current.onDragEnter).toBe('function')
      expect(typeof result.current.onDragOver).toBe('function')
      expect(typeof result.current.onDragLeave).toBe('function')
      expect(typeof result.current.onDrop).toBe('function')
      expect(typeof result.current.reset).toBe('function')
    })
  })

  describe('Drag Enter Event', () => {
    it('[P1] should set isDragActive to true on drag enter', () => {
      // GIVEN: Hook instance
      const { result } = renderHook(() => useDragAndDrop({ onFileDrop: mockOnFileDrop }))
      const event = createMockDragEvent()

      // WHEN: onDragEnter is called
      act(() => {
        result.current.onDragEnter(event)
      })

      // THEN: isDragActive should be true
      expect(result.current.isDragActive).toBe(true)
      expect(event.preventDefault).toHaveBeenCalled()
      expect(event.stopPropagation).toHaveBeenCalled()
    })
  })

  describe('Drag Over Event', () => {
    it('[P1] should set isDragOver to true on drag over', () => {
      // GIVEN: Hook instance
      const { result } = renderHook(() => useDragAndDrop({ onFileDrop: mockOnFileDrop }))
      const event = createMockDragEvent()

      // WHEN: onDragOver is called
      act(() => {
        result.current.onDragOver(event)
      })

      // THEN: isDragOver should be true
      expect(result.current.isDragOver).toBe(true)
      expect(event.preventDefault).toHaveBeenCalled()
      expect(event.stopPropagation).toHaveBeenCalled()
    })
  })

  describe('Drag Leave Event', () => {
    it('[P1] should reset drag states when leaving main container', () => {
      // GIVEN: Hook with active drag state
      const { result } = renderHook(() => useDragAndDrop({ onFileDrop: mockOnFileDrop }))

      // Set up drag active state
      act(() => {
        result.current.onDragEnter(createMockDragEvent())
        result.current.onDragOver(createMockDragEvent())
      })
      expect(result.current.isDragActive).toBe(true)
      expect(result.current.isDragOver).toBe(true)

      // WHEN: onDragLeave is called (leaving container)
      const container = document.createElement('div')
      const event = createMockDragEvent([], {
        currentTarget: container,
        relatedTarget: null, // Leaving to outside
      } as Partial<DragEvent<HTMLElement>>)

      // Mock contains to return false (relatedTarget is not inside currentTarget)
      container.contains = vi.fn().mockReturnValue(false)
      Object.defineProperty(event, 'currentTarget', { value: container })

      act(() => {
        result.current.onDragLeave(event)
      })

      // THEN: Drag states should be reset
      expect(result.current.isDragActive).toBe(false)
      expect(result.current.isDragOver).toBe(false)
    })

    it('[P1] should not reset when moving to child element', () => {
      // GIVEN: Hook with active drag state
      const { result } = renderHook(() => useDragAndDrop({ onFileDrop: mockOnFileDrop }))

      act(() => {
        result.current.onDragEnter(createMockDragEvent())
        result.current.onDragOver(createMockDragEvent())
      })

      // WHEN: onDragLeave is called but moving to child
      const container = document.createElement('div')
      const child = document.createElement('div')
      container.appendChild(child)

      const event = createMockDragEvent([], {
        currentTarget: container,
        relatedTarget: child,
      } as Partial<DragEvent<HTMLElement>>)
      Object.defineProperty(event, 'currentTarget', { value: container })

      act(() => {
        result.current.onDragLeave(event)
      })

      // THEN: Drag states should remain active
      expect(result.current.isDragActive).toBe(true)
      expect(result.current.isDragOver).toBe(true)
    })
  })

  describe('Drop Event', () => {
    it('[P1] should call onFileDrop with dropped files', () => {
      // GIVEN: Hook instance with files
      const { result } = renderHook(() => useDragAndDrop({ onFileDrop: mockOnFileDrop }))
      const files = [createMockFile({ name: 'test.zip' })]
      const event = createMockDragEvent(files)

      // WHEN: onDrop is called
      act(() => {
        result.current.onDrop(event)
      })

      // THEN: onFileDrop should be called with files
      expect(mockOnFileDrop).toHaveBeenCalledWith(files)
      expect(event.preventDefault).toHaveBeenCalled()
      expect(event.stopPropagation).toHaveBeenCalled()
    })

    it('[P1] should reset drag states after drop', () => {
      // GIVEN: Hook with active drag state
      const { result } = renderHook(() => useDragAndDrop({ onFileDrop: mockOnFileDrop }))

      act(() => {
        result.current.onDragEnter(createMockDragEvent())
        result.current.onDragOver(createMockDragEvent())
      })

      // WHEN: onDrop is called
      const files = [createMockFile()]
      act(() => {
        result.current.onDrop(createMockDragEvent(files))
      })

      // THEN: Drag states should be reset
      expect(result.current.isDragActive).toBe(false)
      expect(result.current.isDragOver).toBe(false)
    })

    it('[P1] should not call onFileDrop if no files dropped', () => {
      // GIVEN: Hook instance
      const { result } = renderHook(() => useDragAndDrop({ onFileDrop: mockOnFileDrop }))
      const event = createMockDragEvent([]) // No files

      // WHEN: onDrop is called with no files
      act(() => {
        result.current.onDrop(event)
      })

      // THEN: onFileDrop should not be called
      expect(mockOnFileDrop).not.toHaveBeenCalled()
    })
  })

  describe('File Type Filtering', () => {
    it('[P1] should filter files by accepted extensions', () => {
      // GIVEN: Hook with accept filter
      const { result } = renderHook(() =>
        useDragAndDrop({
          onFileDrop: mockOnFileDrop,
          accept: ['.zip', '.json'],
        })
      )

      const validFile = createMockFile({ name: 'test.zip', type: 'application/zip' })
      const invalidFile = createMockFile({ name: 'test.exe', type: 'application/x-msdownload' })
      const jsonFile = createMockFile({ name: 'package.json', type: 'application/json' })

      const event = createMockDragEvent([validFile, invalidFile, jsonFile])

      // WHEN: onDrop is called
      act(() => {
        result.current.onDrop(event)
      })

      // THEN: Only valid files should be passed to onFileDrop
      expect(mockOnFileDrop).toHaveBeenCalledWith([validFile, jsonFile])
    })

    it('[P1] should accept all files when no filter specified', () => {
      // GIVEN: Hook without accept filter
      const { result } = renderHook(() => useDragAndDrop({ onFileDrop: mockOnFileDrop }))

      const files = [
        createMockFile({ name: 'test.zip' }),
        createMockFile({ name: 'test.exe' }),
        createMockFile({ name: 'test.txt' }),
      ]
      const event = createMockDragEvent(files)

      // WHEN: onDrop is called
      act(() => {
        result.current.onDrop(event)
      })

      // THEN: All files should be passed
      expect(mockOnFileDrop).toHaveBeenCalledWith(files)
    })

    it('[P1] should handle MIME type filtering', () => {
      // GIVEN: Hook with MIME type filter
      const { result } = renderHook(() =>
        useDragAndDrop({
          onFileDrop: mockOnFileDrop,
          accept: ['application/json'],
        })
      )

      const jsonFile = createMockFile({ name: 'data.json', type: 'application/json' })
      const zipFile = createMockFile({ name: 'archive.zip', type: 'application/zip' })

      const event = createMockDragEvent([jsonFile, zipFile])

      // WHEN: onDrop is called
      act(() => {
        result.current.onDrop(event)
      })

      // THEN: Only JSON file should be passed
      expect(mockOnFileDrop).toHaveBeenCalledWith([jsonFile])
    })
  })

  describe('Reset Functionality', () => {
    it('[P1] should reset all drag states', () => {
      // GIVEN: Hook with active drag state
      const { result } = renderHook(() => useDragAndDrop({ onFileDrop: mockOnFileDrop }))

      act(() => {
        result.current.onDragEnter(createMockDragEvent())
        result.current.onDragOver(createMockDragEvent())
      })

      expect(result.current.isDragActive).toBe(true)
      expect(result.current.isDragOver).toBe(true)

      // WHEN: reset is called
      act(() => {
        result.current.reset()
      })

      // THEN: All states should be reset
      expect(result.current.isDragActive).toBe(false)
      expect(result.current.isDragOver).toBe(false)
    })
  })

  describe('Multiple Files', () => {
    it('[P1] should handle multiple files in drop', () => {
      // GIVEN: Hook instance
      const { result } = renderHook(() => useDragAndDrop({ onFileDrop: mockOnFileDrop }))

      const files = [
        createMockFile({ name: 'file1.zip' }),
        createMockFile({ name: 'file2.zip' }),
        createMockFile({ name: 'file3.zip' }),
      ]
      const event = createMockDragEvent(files)

      // WHEN: onDrop is called with multiple files
      act(() => {
        result.current.onDrop(event)
      })

      // THEN: All files should be passed
      expect(mockOnFileDrop).toHaveBeenCalledWith(files)
      expect(mockOnFileDrop.mock.calls[0][0]).toHaveLength(3)
    })
  })

  describe('Edge Cases', () => {
    it('[P2] should handle empty dataTransfer', () => {
      // GIVEN: Hook instance
      const { result } = renderHook(() => useDragAndDrop({ onFileDrop: mockOnFileDrop }))

      const event = {
        preventDefault: vi.fn(),
        stopPropagation: vi.fn(),
        dataTransfer: null,
        currentTarget: document.createElement('div'),
        relatedTarget: null,
      } as unknown as DragEvent<HTMLElement>

      // WHEN: onDrop is called with null dataTransfer
      act(() => {
        result.current.onDrop(event)
      })

      // THEN: Should not throw and not call onFileDrop
      expect(mockOnFileDrop).not.toHaveBeenCalled()
    })

    it('[P2] should handle file extensions case-insensitively', () => {
      // GIVEN: Hook with lowercase extension filter
      const { result } = renderHook(() =>
        useDragAndDrop({
          onFileDrop: mockOnFileDrop,
          accept: ['.zip'],
        })
      )

      const uppercaseFile = createMockFile({ name: 'TEST.ZIP', type: 'application/zip' })
      const mixedCaseFile = createMockFile({ name: 'Test.Zip', type: 'application/zip' })

      const event = createMockDragEvent([uppercaseFile, mixedCaseFile])

      // WHEN: onDrop is called
      act(() => {
        result.current.onDrop(event)
      })

      // THEN: Both files should be accepted
      expect(mockOnFileDrop).toHaveBeenCalled()
      expect(mockOnFileDrop.mock.calls[0][0]).toHaveLength(2)
    })
  })
})
