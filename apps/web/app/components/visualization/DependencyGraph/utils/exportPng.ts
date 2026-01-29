/**
 * PNG Export Utility
 *
 * Exports the dependency graph as a PNG image by first rendering
 * the SVG to a Canvas element, then converting to PNG blob.
 *
 * @see Story 4.6: Export Graph as PNG/SVG Images - AC2, AC4, AC5
 */

import type { ExportOptions, ExportResult } from '../types'
import { exportSvg } from './exportSvg'

export interface ExportPngParams {
  /** The SVG element to export */
  svgElement: SVGSVGElement
  /** Export configuration options */
  options: ExportOptions
  /** Project name for filename generation */
  projectName: string
  /** Optional legend SVG string to embed */
  legendSvg?: string
  /** Optional progress callback (0-100) */
  onProgress?: (progress: number) => void
}

/**
 * Exports the dependency graph as a PNG file.
 * Renders SVG to Canvas, applies resolution multiplier, then converts to PNG.
 */
export async function exportPng({
  svgElement,
  options,
  projectName,
  legendSvg,
  onProgress,
}: ExportPngParams): Promise<ExportResult> {
  onProgress?.(10)

  // First export as SVG to get a fully-styled standalone SVG
  const svgResult = await exportSvg({
    svgElement,
    options: { ...options, format: 'svg' },
    projectName,
    legendSvg,
  })

  onProgress?.(30)

  // Render SVG blob to an Image element
  const img = new Image()
  const svgUrl = URL.createObjectURL(svgResult.blob)

  return new Promise<ExportResult>((resolve, reject) => {
    img.onload = () => {
      onProgress?.(50)

      const width = svgResult.width * options.resolution
      const height = svgResult.height * options.resolution

      const canvas = document.createElement('canvas')
      canvas.width = width
      canvas.height = height

      const ctx = canvas.getContext('2d')
      if (!ctx) {
        URL.revokeObjectURL(svgUrl)
        reject(new Error('Failed to get canvas context'))
        return
      }

      // Fill background
      if (options.backgroundColor !== 'transparent') {
        ctx.fillStyle = options.backgroundColor
        ctx.fillRect(0, 0, width, height)
      }

      // Scale and draw
      ctx.scale(options.resolution, options.resolution)
      ctx.drawImage(img, 0, 0)

      onProgress?.(80)

      // Convert to PNG blob
      canvas.toBlob(
        (blob) => {
          URL.revokeObjectURL(svgUrl)

          if (!blob) {
            reject(new Error('Failed to create PNG blob'))
            return
          }

          onProgress?.(100)

          const timestamp = new Date().toISOString().split('T')[0]
          const resolutionSuffix = options.resolution > 1 ? `@${options.resolution}x` : ''
          const filename = `${projectName}-dependency-graph-${timestamp}${resolutionSuffix}.png`

          resolve({
            blob,
            filename,
            width,
            height,
          })
        },
        'image/png',
        0.95
      )
    }

    img.onerror = () => {
      URL.revokeObjectURL(svgUrl)
      reject(new Error('Failed to load SVG for PNG conversion'))
    }

    img.src = svgUrl
  })
}
