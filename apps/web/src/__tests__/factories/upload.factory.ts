/**
 * Test Factories for Upload-related types
 *
 * Pattern: Pure function with Partial<T> overrides (TEA knowledge base: data-factories.md)
 */

import type { FileProcessingResult, PackageJsonFile, UploadedFile } from '@monoguard/types'

/**
 * Creates a mock File object for testing
 */
export function createMockFile(
  overrides: Partial<File> & { name?: string; type?: string; size?: number } = {}
): File {
  const {
    name = 'test-file.zip',
    type = 'application/zip',
    size = 1024 * 1024, // 1MB default
  } = overrides

  const file = new File([''], name, { type })
  Object.defineProperty(file, 'size', { value: size, writable: false })

  return file
}

/**
 * Creates an UploadedFile object
 */
export function createUploadedFile(overrides: Partial<UploadedFile> = {}): UploadedFile {
  return {
    filename: 'test-file.zip',
    originalName: 'test-file.zip',
    fileSize: 1024 * 1024,
    mimeType: 'application/zip',
    path: '/uploads/test-file.zip',
    extractedFiles: 5,
    ...overrides,
  }
}

/**
 * Creates a PackageJsonFile object
 */
export function createPackageJsonFile(overrides: Partial<PackageJsonFile> = {}): PackageJsonFile {
  return {
    name: 'test-package',
    version: '1.0.0',
    path: '/packages/test-package/package.json',
    dependencies: {
      react: '^18.0.0',
      'react-dom': '^18.0.0',
    },
    devDependencies: {
      typescript: '^5.0.0',
      vitest: '^1.0.0',
    },
    metadata: {
      fileSize: 512,
      lastModified: new Date().toISOString(),
      dependencyCount: 2,
      devDependencyCount: 2,
    },
    ...overrides,
  }
}

/**
 * Creates a FileProcessingResult object
 */
export function createFileProcessingResult(
  overrides: Partial<FileProcessingResult> = {}
): FileProcessingResult {
  return {
    uploadId: 'upload-123',
    files: [createUploadedFile()],
    packageJsonFiles: [createPackageJsonFile()],
    analysisReady: true,
    metadata: {
      totalSize: 1024 * 1024,
      processedAt: new Date().toISOString(),
      processingDuration: 150,
    },
    ...overrides,
  }
}

/**
 * Creates a mock UploadProgress object
 */
export function createUploadProgress(
  overrides: { loaded?: number; total?: number; percentage?: number } = {}
) {
  const total = overrides.total ?? 1024 * 1024
  const loaded = overrides.loaded ?? total / 2
  const percentage = overrides.percentage ?? Math.round((loaded * 100) / total)

  return { loaded, total, percentage }
}
