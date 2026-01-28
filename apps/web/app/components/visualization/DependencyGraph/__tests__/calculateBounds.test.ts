/**
 * Tests for calculateBounds utilities
 *
 * @see Story 4.4: Add Zoom, Pan, and Navigation Controls
 */
import { describe, expect, it } from 'vitest'

import {
  calculateFitTransform,
  calculateNodeBounds,
  calculateViewportBounds,
} from '../utils/calculateBounds'

describe('calculateBounds utilities', () => {
  describe('calculateNodeBounds', () => {
    it('should return zero bounds for empty array', () => {
      const bounds = calculateNodeBounds([])

      expect(bounds).toEqual({ x: 0, y: 0, width: 0, height: 0 })
    })

    it('should calculate bounds for single node', () => {
      const nodes = [{ x: 100, y: 200 }]
      const bounds = calculateNodeBounds(nodes)

      // With default padding of 20
      expect(bounds.x).toBe(80) // 100 - 20
      expect(bounds.y).toBe(180) // 200 - 20
      expect(bounds.width).toBe(40) // 0 + 20*2
      expect(bounds.height).toBe(40) // 0 + 20*2
    })

    it('should calculate bounds for multiple nodes', () => {
      const nodes = [
        { x: 0, y: 0 },
        { x: 100, y: 200 },
        { x: 50, y: 100 },
      ]
      const bounds = calculateNodeBounds(nodes)

      // minX=0, maxX=100, minY=0, maxY=200
      expect(bounds.x).toBe(-20) // 0 - 20
      expect(bounds.y).toBe(-20) // 0 - 20
      expect(bounds.width).toBe(140) // 100 - 0 + 40
      expect(bounds.height).toBe(240) // 200 - 0 + 40
    })

    it('should handle custom padding', () => {
      const nodes = [{ x: 100, y: 100 }]
      const bounds = calculateNodeBounds(nodes, 50)

      expect(bounds.x).toBe(50) // 100 - 50
      expect(bounds.y).toBe(50) // 100 - 50
      expect(bounds.width).toBe(100) // 0 + 50*2
      expect(bounds.height).toBe(100) // 0 + 50*2
    })

    it('should handle undefined x/y as 0', () => {
      const nodes = [
        { x: undefined, y: undefined },
        { x: 100, y: 100 },
      ]
      const bounds = calculateNodeBounds(nodes)

      // undefined treated as 0
      expect(bounds.x).toBe(-20) // 0 - 20
      expect(bounds.y).toBe(-20) // 0 - 20
      expect(bounds.width).toBe(140) // 100 - 0 + 40
      expect(bounds.height).toBe(140) // 100 - 0 + 40
    })

    it('should handle negative coordinates', () => {
      const nodes = [
        { x: -100, y: -50 },
        { x: 100, y: 50 },
      ]
      const bounds = calculateNodeBounds(nodes)

      expect(bounds.x).toBe(-120) // -100 - 20
      expect(bounds.y).toBe(-70) // -50 - 20
      expect(bounds.width).toBe(240) // 200 + 40
      expect(bounds.height).toBe(140) // 100 + 40
    })
  })

  describe('calculateViewportBounds', () => {
    it('should calculate viewport at identity transform', () => {
      const transform = { k: 1, x: 0, y: 0 }
      const bounds = calculateViewportBounds(transform, 800, 600)

      // Use toBeCloseTo for floating-point comparisons (handles -0 vs 0)
      expect(bounds.x).toBeCloseTo(0)
      expect(bounds.y).toBeCloseTo(0)
      expect(bounds.width).toBe(800)
      expect(bounds.height).toBe(600)
    })

    it('should calculate viewport when zoomed in', () => {
      const transform = { k: 2, x: 0, y: 0 } // 200% zoom
      const bounds = calculateViewportBounds(transform, 800, 600)

      expect(bounds.x).toBeCloseTo(0)
      expect(bounds.y).toBeCloseTo(0)
      expect(bounds.width).toBe(400) // 800 / 2
      expect(bounds.height).toBe(300) // 600 / 2
    })

    it('should calculate viewport when panned', () => {
      const transform = { k: 1, x: 100, y: 50 } // Panned right and down
      const bounds = calculateViewportBounds(transform, 800, 600)

      expect(bounds.x).toBe(-100) // Viewport moved left in graph coords
      expect(bounds.y).toBe(-50)
      expect(bounds.width).toBe(800)
      expect(bounds.height).toBe(600)
    })

    it('should calculate viewport when zoomed and panned', () => {
      const transform = { k: 2, x: 200, y: 100 }
      const bounds = calculateViewportBounds(transform, 800, 600)

      expect(bounds.x).toBe(-100) // -200 / 2
      expect(bounds.y).toBe(-50) // -100 / 2
      expect(bounds.width).toBe(400) // 800 / 2
      expect(bounds.height).toBe(300) // 600 / 2
    })

    it('should calculate viewport when zoomed out', () => {
      const transform = { k: 0.5, x: 0, y: 0 } // 50% zoom
      const bounds = calculateViewportBounds(transform, 800, 600)

      expect(bounds.x).toBeCloseTo(0)
      expect(bounds.y).toBeCloseTo(0)
      expect(bounds.width).toBe(1600) // 800 / 0.5
      expect(bounds.height).toBe(1200) // 600 / 0.5
    })
  })

  describe('calculateFitTransform', () => {
    it('should return identity for empty bounds', () => {
      const bounds = { x: 0, y: 0, width: 0, height: 0 }
      const result = calculateFitTransform(bounds, 800, 600)

      expect(result).toEqual({ scale: 1, translateX: 0, translateY: 0 })
    })

    it('should fit bounds smaller than container', () => {
      const bounds = { x: 0, y: 0, width: 200, height: 100 }
      const result = calculateFitTransform(bounds, 800, 600, 40)

      // (800 - 80) / 200 = 3.6, (600 - 80) / 100 = 5.2
      // scale = min(3.6, 5.2) = 3.6
      expect(result.scale).toBeCloseTo(3.6)
    })

    it('should fit bounds larger than container', () => {
      const bounds = { x: 0, y: 0, width: 1600, height: 1200 }
      const result = calculateFitTransform(bounds, 800, 600, 40)

      // (800 - 80) / 1600 = 0.45, (600 - 80) / 1200 = 0.433
      // scale = min(0.45, 0.433) = 0.433
      expect(result.scale).toBeCloseTo(0.433, 2)
    })

    it('should respect maxScale', () => {
      const bounds = { x: 0, y: 0, width: 100, height: 100 }
      const result = calculateFitTransform(bounds, 800, 600, 40, 2)

      // Without max: (720/100, 520/100) = (7.2, 5.2)
      // With maxScale 2: scale = 2
      expect(result.scale).toBe(2)
    })

    it('should center bounds in container', () => {
      const bounds = { x: 0, y: 0, width: 400, height: 300 }
      const result = calculateFitTransform(bounds, 800, 600, 40)

      // scale = min((800-80)/400, (600-80)/300) = min(1.8, 1.73) = 1.73
      const expectedScale = (600 - 80) / 300

      // translateX = (800 - 400*scale) / 2 - 0*scale
      const expectedTranslateX = (800 - 400 * expectedScale) / 2
      // translateY = (600 - 300*scale) / 2 - 0*scale
      const expectedTranslateY = (600 - 300 * expectedScale) / 2

      expect(result.translateX).toBeCloseTo(expectedTranslateX)
      expect(result.translateY).toBeCloseTo(expectedTranslateY)
    })

    it('should handle offset bounds', () => {
      const bounds = { x: 100, y: 50, width: 200, height: 150 }
      const result = calculateFitTransform(bounds, 800, 600, 40)

      // Should account for bounds offset
      const expectedScale = Math.min((800 - 80) / 200, (600 - 80) / 150)

      expect(result.scale).toBeCloseTo(expectedScale, 2)
      // Translation should include offset compensation
      expect(typeof result.translateX).toBe('number')
      expect(typeof result.translateY).toBe('number')
    })

    it('should handle custom padding', () => {
      const bounds = { x: 0, y: 0, width: 400, height: 300 }
      const result = calculateFitTransform(bounds, 800, 600, 100)

      // scale = min((800-200)/400, (600-200)/300) = min(1.5, 1.33) = 1.33
      expect(result.scale).toBeCloseTo((600 - 200) / 300, 2)
    })
  })
})
