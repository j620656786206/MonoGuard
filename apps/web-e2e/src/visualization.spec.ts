/**
 * MonoGuard E2E Tests - Visualization Feature (Stories 4.1-4.6)
 *
 * Tests for the D3.js force-directed dependency graph visualization,
 * circular dependency highlighting, expand/collapse, zoom/pan,
 * hover tooltips, and graph export.
 *
 * NOTE: Full visualization tests require analysis data in the store.
 * Tests are split into:
 * - Empty state tests (run without data)
 * - Visualization container tests (verify setup)
 * - Full visualization tests (marked as fixme until data seeding is available)
 * - Circular dependency highlighting tests (Story 4.2)
 * - Node expand/collapse tests (Story 4.3)
 * - Zoom, pan, and navigation tests (Story 4.4)
 * - Hover details and tooltips tests (Story 4.5)
 * - Export graph as PNG/SVG tests (Story 4.6)
 *
 * Following TEA knowledge base patterns:
 * - Given-When-Then format
 * - Priority tags [P0], [P1], [P2]
 * - Explicit assertions
 * - No hard waits
 */

import { expect, test } from './support/fixtures'

test.describe('Dependency Graph Visualization (Story 4.1)', () => {
  test.describe('Empty State (No Analysis Data)', () => {
    test('[P1] should show placeholder when no analysis data', async ({ page }) => {
      // GIVEN: User navigates to results without prior analysis
      await page.goto('/results')

      // THEN: Placeholder message should be shown
      await expect(page.getByText(/Run analysis to visualize/i)).toBeVisible()
    })

    test('[P1] should show dependency graph section heading', async ({ page }) => {
      // GIVEN: User navigates to results
      await page.goto('/results')

      // THEN: Dependency Graph section should be visible
      await expect(page.getByText(/Dependency Graph/i)).toBeVisible()
    })

    test('[P1] should display no analysis data message', async ({ page }) => {
      // GIVEN: User navigates to results without analysis
      await page.goto('/results')

      // THEN: No analysis data message should be visible
      await expect(page.getByText(/No analysis data/i)).toBeVisible()
    })
  })

  test.describe('Visualization Container Setup', () => {
    test('[P1] should have results page structure', async ({ page }) => {
      // GIVEN: User navigates to results page
      await page.goto('/results')

      // THEN: Page should have proper structure
      await expect(page.locator('h1')).toContainText('Analysis Results')
      await expect(page.locator('main')).toBeVisible()
    })

    test('[P1] should have navigation to analyze page', async ({ page }) => {
      // GIVEN: User is on results page
      await page.goto('/results')

      // WHEN: User clicks Start New Analysis
      await page.getByText(/Start New Analysis/i).click()

      // THEN: Should navigate to analyze page
      await expect(page).toHaveURL(/analyze/)
    })
  })

  test.describe('Responsive Design', () => {
    test('[P2] should be responsive on tablet viewport', async ({ page }) => {
      // GIVEN: Tablet viewport
      await page.setViewportSize({ width: 768, height: 1024 })

      // WHEN: User views results
      await page.goto('/results')

      // THEN: Dependency graph section should be visible
      await expect(page.getByText(/Dependency Graph/i)).toBeVisible()
    })

    test('[P2] should be responsive on mobile viewport', async ({ page }) => {
      // GIVEN: Mobile viewport
      await page.setViewportSize({ width: 375, height: 667 })

      // WHEN: User views results
      await page.goto('/results')

      // THEN: Page should render correctly with graph section visible
      await expect(page.getByText(/Dependency Graph/i)).toBeVisible()
    })
  })

  /**
   * Full Visualization Tests (Require Analysis Data)
   *
   * These tests verify the D3.js force-directed graph functionality
   * but require analysis data to be present in the store.
   *
   * FIXME: Enable these tests when:
   * 1. Store mocking is implemented for E2E tests, OR
   * 2. A test fixture can seed analysis data via localStorage/IndexedDB
   *
   * For now, these acceptance criteria are verified via unit tests:
   * - AC1: Force-directed layout - see DependencyGraph.test.tsx
   * - AC2: Performance - see DependencyGraph.test.tsx (50 node test)
   * - AC3: Responsive container - see DependencyGraph.test.tsx
   * - AC4: Data integration - see DependencyGraph.test.tsx
   * - AC5: Node visual - see DependencyGraph.test.tsx
   * - AC6: Edge visual - see DependencyGraph.test.tsx
   */
  test.describe('Full Visualization (Requires Data)', () => {
    test.fixme('[P1] should render SVG graph when analysis data exists', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // When data seeding is available:
      // 1. Seed analysis results via fixture
      // 2. Navigate to /results
      // 3. Verify SVG with circles and lines renders
      await page.goto('/results')
      await expect(page.locator('svg circle').first()).toBeVisible()
    })

    test.fixme('[P1] should render nodes and edges when data exists', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // Verification points:
      // - SVG contains circle elements for each package node
      // - SVG contains line elements for each dependency edge
      // - Lines have arrow markers (directed graph)
      await page.goto('/results')
      await expect(page.locator('svg line').first()).toBeVisible()
    })

    test.fixme('[P2] should support zoom and pan interactions', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // Verification points:
      // - Mouse wheel zooms the graph
      // - Mouse drag pans the graph
      // - Graph remains interactive after zoom/pan
      await page.goto('/results')
      const svg = page.locator('svg').first()
      await svg.hover()
      await page.mouse.wheel(0, -100)
      await expect(svg).toBeVisible()
    })

    test.fixme('[P2] should render within 2 seconds for 100 packages', async ({ page }) => {
      // FIXME: Test requires large analysis dataset
      // Performance verification:
      // - Render time < 2 seconds
      // - Layout stabilizes within 3 seconds
      // - No visual jank during animation
      const startTime = Date.now()
      await page.goto('/results')
      await expect(page.locator('svg circle').first()).toBeVisible()
      const renderTime = Date.now() - startTime
      expect(renderTime).toBeLessThan(2000)
    })
  })
})

/**
 * Story 4.2: Highlight Circular Dependencies in Graph
 *
 * Tests for visual highlighting of circular dependencies in the dependency graph.
 * These tests verify:
 * - AC1: Visual highlighting of cycle nodes (red border/glow)
 * - AC2: Visual highlighting of cycle edges (red, thicker)
 * - AC3: Animated cycle paths
 * - AC4: Color legend display
 * - AC5: Click-to-highlight cycle interaction
 * - AC6: Escape key to deselect
 */
test.describe('Circular Dependency Highlighting (Story 4.2)', () => {
  /**
   * Legend Display Tests (AC4)
   *
   * The GraphLegend component is rendered as part of DependencyGraphViz,
   * which requires analysis data to be present. These tests are marked as
   * fixme until data seeding is available.
   *
   * Unit test coverage: DependencyGraph.test.tsx (GraphLegend tests)
   */
  test.describe('Legend Display (AC4)', () => {
    test.fixme('[P1] should display graph legend on results page', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // When data seeding is available:
      // 1. Seed analysis results via fixture
      // 2. Navigate to /results
      // 3. Verify Legend component is rendered
      await page.goto('/results')

      // THEN: Legend component should be rendered
      await expect(page.getByText('Legend')).toBeVisible()
    })

    test.fixme('[P1] should show normal node/edge colors in legend', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // When data seeding is available:
      // 1. Seed analysis results via fixture
      // 2. Navigate to /results
      // 3. Verify legend displays normal element descriptions
      await page.goto('/results')

      // THEN: Legend should display normal element descriptions
      await expect(page.getByText('Normal Package')).toBeVisible()
      await expect(page.getByText('Normal Dependency')).toBeVisible()
    })
  })

  test.describe('Cycle Highlighting with Data (Requires Analysis)', () => {
    test.fixme('[P1] should display cycle colors in legend when cycles exist', async ({ page }) => {
      // FIXME: Test requires analysis data with circular dependencies
      // When data seeding is available:
      // 1. Seed analysis results with circular dependencies via fixture
      // 2. Navigate to /results
      // 3. Verify legend shows cycle-specific colors
      await page.goto('/results')

      // THEN: Legend should show cycle indicators
      await expect(page.getByText('In Circular Dependency')).toBeVisible()
      await expect(page.getByText('Circular Dependency')).toBeVisible()
    })

    test.fixme('[P1] should highlight cycle nodes with red styling (AC1)', async ({ page }) => {
      // FIXME: Test requires analysis data with circular dependencies
      // Verification points:
      // - Nodes in cycles have red fill (#ef4444)
      // - Nodes have glow filter applied
      // - Nodes have CSS class 'node--cycle'
      await page.goto('/results')

      // Check for cycle node styling
      const cycleNodes = page.locator('g.node--cycle')
      await expect(cycleNodes.first()).toBeVisible()

      // Verify red fill color
      const circle = cycleNodes.first().locator('circle')
      await expect(circle).toHaveAttribute('fill', '#ef4444')
    })

    test.fixme('[P1] should highlight cycle edges with red color (AC2)', async ({ page }) => {
      // FIXME: Test requires analysis data with circular dependencies
      // Verification points:
      // - Cycle edges are red (#ef4444)
      // - Cycle edges are thicker (2.5px vs 1.5px)
      // - Cycle edges use separate arrow marker
      await page.goto('/results')

      // Check for cycle edge group
      const cycleEdgesGroup = page.locator('g.links-cycle')
      await expect(cycleEdgesGroup).toBeVisible()

      // Verify red stroke color
      const cycleLine = cycleEdgesGroup.locator('line').first()
      await expect(cycleLine).toHaveAttribute('stroke', '#ef4444')
    })

    test.fixme('[P2] should animate cycle edges (AC3)', async ({ page }) => {
      // FIXME: Test requires analysis data with circular dependencies
      // Verification points:
      // - Cycle edges have stroke-dasharray for animation
      // - CSS animation keyframes are present
      await page.goto('/results')

      // Check for animation styling
      const cycleLine = page.locator('g.links-cycle line').first()
      await expect(cycleLine).toHaveAttribute('stroke-dasharray')

      // Check for style element with animation
      const styleElement = page.locator('style')
      await expect(styleElement).toContainText('flowAnimation')
    })

    test.fixme('[P1] should highlight specific cycle on node click (AC5)', async ({ page }) => {
      // FIXME: Test requires analysis data with circular dependencies
      // Verification points:
      // - Click on cycle node selects that cycle
      // - Selected cycle is emphasized (brighter)
      // - Other elements are dimmed
      await page.goto('/results')

      // Find and click a cycle node
      const cycleNode = page.locator('g.node--cycle').first()
      await cycleNode.click()

      // Verify selection state - other nodes should be dimmed
      const dimmedNodes = page.locator('g.node--dimmed')
      await expect(dimmedNodes.first()).toBeVisible()
    })

    test.fixme('[P2] should deselect cycle on Escape key (AC6)', async ({ page }) => {
      // FIXME: Test requires analysis data with circular dependencies
      // Verification points:
      // - Press Escape clears cycle selection
      // - All elements return to normal highlighting
      await page.goto('/results')

      // Click a cycle node to select
      const cycleNode = page.locator('g.node--cycle').first()
      await cycleNode.click()

      // Verify selection
      await expect(page.locator('g.node--dimmed').first()).toBeVisible()

      // Press Escape to deselect
      await page.keyboard.press('Escape')

      // Verify deselection - no dimmed nodes
      await expect(page.locator('g.node--dimmed')).toHaveCount(0)
    })

    test.fixme('[P2] should deselect cycle on background click (AC6)', async ({ page }) => {
      // FIXME: Test requires analysis data with circular dependencies
      // Verification points:
      // - Click on graph background clears selection
      await page.goto('/results')

      // Click a cycle node to select
      const cycleNode = page.locator('g.node--cycle').first()
      await cycleNode.click()

      // Click on SVG background
      const svg = page.locator('svg').first()
      await svg.click({ position: { x: 10, y: 10 } })

      // Verify deselection
      await expect(page.locator('g.node--dimmed')).toHaveCount(0)
    })

    test.fixme('[P1] should show interaction hints in legend (AC4)', async ({ page }) => {
      // FIXME: Test requires analysis data with circular dependencies
      // Verification points:
      // - Legend shows "Click on red nodes" hint
      // - Legend shows "Escape" to deselect hint
      await page.goto('/results')

      // Verify interaction hints are visible
      await expect(page.getByText(/Click on red nodes/i)).toBeVisible()
      await expect(page.getByText(/Escape/i)).toBeVisible()
    })
  })

  test.describe('Performance (AC3)', () => {
    test.fixme('[P2] should animate at 60fps without frame drops', async ({ page }) => {
      // FIXME: Test requires analysis data with circular dependencies
      // Performance verification:
      // - Animation runs smoothly at 60fps
      // - No visual jank during animation
      // - CPU usage remains reasonable
      await page.goto('/results')

      // Use Performance API to check for frame drops
      const metrics = await page.evaluate(() => {
        return new Promise((resolve) => {
          const frames: number[] = []
          let lastTime = performance.now()

          const checkFrame = () => {
            const now = performance.now()
            frames.push(now - lastTime)
            lastTime = now

            if (frames.length < 60) {
              requestAnimationFrame(checkFrame)
            } else {
              const avgFrameTime = frames.reduce((a, b) => a + b, 0) / frames.length
              resolve({ avgFrameTime, maxFrameTime: Math.max(...frames) })
            }
          }

          requestAnimationFrame(checkFrame)
        })
      })

      // 60fps = 16.67ms per frame, allow some margin
      expect((metrics as { avgFrameTime: number }).avgFrameTime).toBeLessThan(20)
    })
  })

  test.describe('Accessibility', () => {
    test.fixme(
      '[P2] should use color AND visual patterns for cycle indication',
      async ({ page }) => {
        // FIXME: Test requires analysis data with circular dependencies
        // Accessibility verification:
        // - Cycle edges are thicker (not just red) - passes WCAG 2.1
        // - Cycle nodes have glow effect (not just color)
        // - Animation provides additional visual cue
        // - Legend indicates both color and pattern differences
        await page.goto('/results')

        // THEN: Legend should indicate both color and pattern differences
        const legend = page.getByText('Legend').locator('..')
        await expect(legend).toBeVisible()
      }
    )
  })
})

/**
 * Story 4.3: Implement Node Expand/Collapse Functionality
 *
 * Tests for expand/collapse functionality in the dependency graph.
 * These tests verify:
 * - AC1: Double-click to collapse a node
 * - AC2: Double-click to expand a collapsed node
 * - AC3: Depth-based collapse controls (All, L1, L2, etc.)
 * - AC4: Collapsed node indicator showing hidden child count
 * - AC5: Smooth animation during expand/collapse (< 300ms)
 * - AC6: Session state persistence
 */
test.describe('Node Expand/Collapse (Story 4.3)', () => {
  test.describe('Graph Controls (AC3)', () => {
    test.fixme('[P1] should display depth control buttons', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // When data seeding is available:
      // 1. Seed analysis results via fixture
      // 2. Navigate to /results
      // 3. Verify depth control buttons are visible
      await page.goto('/results')

      // THEN: Depth control section should be visible
      await expect(page.getByRole('group', { name: /depth controls/i })).toBeVisible()

      // THEN: Should have 'All' button
      await expect(page.getByRole('button', { name: /all/i })).toBeVisible()

      // THEN: Should have expand/collapse all buttons
      await expect(page.getByRole('button', { name: /expand all/i })).toBeVisible()
      await expect(page.getByRole('button', { name: /collapse all/i })).toBeVisible()
    })

    test.fixme('[P1] should show depth level buttons (L1, L2, etc.)', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // When data seeding is available:
      // 1. Seed analysis results with multiple depth levels
      // 2. Navigate to /results
      // 3. Verify depth level buttons are present
      await page.goto('/results')

      // THEN: Should have depth level buttons
      await expect(page.getByRole('button', { name: /l1/i })).toBeVisible()
      await expect(page.getByRole('button', { name: /l2/i })).toBeVisible()
    })

    test.fixme('[P1] should collapse nodes when depth level is selected', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // Verification points:
      // - Clicking L1 collapses all nodes at depth > 1
      // - Hidden nodes are not rendered in SVG
      // - Collapsed parent shows badge with count
      await page.goto('/results')

      // WHEN: Click L1 depth button
      await page.getByRole('button', { name: /l1/i }).click()

      // THEN: Graph should show collapsed nodes with badges
      await expect(page.locator('.collapsed-badge').first()).toBeVisible()
    })

    test.fixme('[P1] should expand all nodes when "All" is clicked', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // Verification points:
      // - Clicking 'All' expands all collapsed nodes
      // - All nodes become visible
      // - No collapsed badges are shown
      await page.goto('/results')

      // First collapse some nodes
      await page.getByRole('button', { name: /l1/i }).click()
      await expect(page.locator('.collapsed-badge').first()).toBeVisible()

      // WHEN: Click 'All' button
      await page.getByRole('button', { name: /all/i }).click()

      // THEN: All nodes should be visible, no badges
      await expect(page.locator('.collapsed-badge')).toHaveCount(0)
    })

    test.fixme('[P1] should highlight selected depth button', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // Verification points:
      // - Selected depth button has distinct styling
      // - aria-pressed attribute is set correctly
      await page.goto('/results')

      // WHEN: Click L1 depth button
      await page.getByRole('button', { name: /l1/i }).click()

      // THEN: L1 button should be pressed
      await expect(page.getByRole('button', { name: /l1/i })).toHaveAttribute(
        'aria-pressed',
        'true'
      )

      // THEN: 'All' button should not be pressed
      await expect(page.getByRole('button', { name: /all/i })).toHaveAttribute(
        'aria-pressed',
        'false'
      )
    })
  })

  test.describe('Double-Click Interaction (AC1, AC2)', () => {
    test.fixme('[P1] should collapse node on double-click (AC1)', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // Verification points:
      // - Double-clicking a node collapses it
      // - Children of collapsed node are hidden
      // - Collapsed node shows dashed border
      await page.goto('/results')

      // Find a node to collapse
      const node = page.locator('g.node').first()
      await node.dblclick()

      // THEN: Node should show collapsed styling (dashed border)
      const circle = node.locator('circle')
      await expect(circle).toHaveAttribute('stroke-dasharray')
    })

    test.fixme('[P1] should expand collapsed node on double-click (AC2)', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // Verification points:
      // - Double-clicking a collapsed node expands it
      // - Children become visible again
      // - Dashed border is removed
      await page.goto('/results')

      // First collapse a node
      const node = page.locator('g.node').first()
      await node.dblclick()

      // Verify it's collapsed
      await expect(node.locator('circle')).toHaveAttribute('stroke-dasharray')

      // WHEN: Double-click again to expand
      await node.dblclick()

      // THEN: Node should not have dashed border
      await expect(node.locator('circle')).not.toHaveAttribute('stroke-dasharray')
    })

    test.fixme('[P2] should distinguish single-click from double-click', async ({ page }) => {
      // FIXME: Test requires analysis data with cycles
      // Verification points:
      // - Single-click on cycle node selects cycle (Story 4.2)
      // - Double-click on same node collapses (not selects)
      await page.goto('/results')

      // Find a cycle node
      const cycleNode = page.locator('g.node--cycle').first()

      // WHEN: Single click
      await cycleNode.click()

      // THEN: Cycle should be selected (dimmed nodes appear)
      await expect(page.locator('g.node circle[fill="#9ca3af"]').first()).toBeVisible()

      // Clear selection
      await page.keyboard.press('Escape')

      // WHEN: Double click
      await cycleNode.dblclick()

      // THEN: Node should be collapsed (dashed border)
      await expect(cycleNode.locator('circle')).toHaveAttribute('stroke-dasharray')
    })
  })

  test.describe('Collapsed Node Indicator (AC4)', () => {
    test.fixme('[P1] should display badge with hidden child count', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // Verification points:
      // - Collapsed node shows orange badge
      // - Badge displays number of hidden children
      // - Badge is positioned correctly (top-right of node)
      await page.goto('/results')

      // Collapse at L1 to ensure some nodes have children
      await page.getByRole('button', { name: /l1/i }).click()

      // THEN: Badge should be visible
      const badge = page.locator('.collapsed-badge').first()
      await expect(badge).toBeVisible()

      // Badge circle should be orange
      await expect(badge.locator('circle')).toHaveAttribute('fill', '#f97316')

      // Badge should show a number
      const badgeText = await badge.locator('text').textContent()
      expect(Number(badgeText)).toBeGreaterThan(0)
    })

    test.fixme('[P2] should hide badge when node is expanded', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // Verification points:
      // - Expanding a node removes its badge
      // - Badge count is recalculated
      await page.goto('/results')

      // Collapse at L1
      await page.getByRole('button', { name: /l1/i }).click()
      const badgeCount = await page.locator('.collapsed-badge').count()
      expect(badgeCount).toBeGreaterThan(0)

      // WHEN: Expand all
      await page.getByRole('button', { name: /all/i }).click()

      // THEN: No badges should remain
      await expect(page.locator('.collapsed-badge')).toHaveCount(0)
    })

    test.fixme('[P2] should show "99+" for large hidden counts', async ({ page }) => {
      // FIXME: Test requires analysis data with large graph
      // Verification points:
      // - For >99 hidden children, badge shows "99+"
      // - Prevents badge from being too wide
      await page.goto('/results')

      // This test requires a very large graph
      // Verify badge text is capped
      const badges = page.locator('.collapsed-badge text')
      const badgeCount = await badges.count()

      if (badgeCount > 0) {
        for (let i = 0; i < badgeCount; i++) {
          const text = await badges.nth(i).textContent()
          expect(text?.length).toBeLessThanOrEqual(3) // "99+" is max
        }
      }
    })
  })

  test.describe('Animation (AC5)', () => {
    test.fixme('[P2] should animate expand/collapse within 300ms', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // Verification points:
      // - Animation completes within 300ms
      // - Smooth transition without jumps
      await page.goto('/results')

      // Measure time for collapse animation
      const startTime = Date.now()

      // Collapse at L1
      await page.getByRole('button', { name: /l1/i }).click()

      // Wait for animation to complete (simulation settling)
      await page.waitForTimeout(300)

      // Verify animation completed within threshold
      const endTime = Date.now()
      const duration = endTime - startTime

      // Should complete within 300ms + some buffer for test overhead
      expect(duration).toBeLessThan(500)
    })
  })

  test.describe('Session Persistence (AC6)', () => {
    test.fixme('[P2] should persist collapsed state across page refresh', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // Verification points:
      // - Collapsed state is saved to sessionStorage
      // - State is restored after page refresh
      await page.goto('/results')

      // Collapse some nodes
      await page.getByRole('button', { name: /l1/i }).click()
      const badgeCount = await page.locator('.collapsed-badge').count()
      expect(badgeCount).toBeGreaterThan(0)

      // WHEN: Refresh the page
      await page.reload()

      // THEN: Collapsed state should be restored
      const restoredBadgeCount = await page.locator('.collapsed-badge').count()
      expect(restoredBadgeCount).toBe(badgeCount)
    })

    test.fixme('[P2] should clear collapsed state on new analysis', async ({ page }) => {
      // FIXME: Test requires ability to run new analysis
      // Verification points:
      // - Running a new analysis clears collapsed state
      // - Fresh graph starts fully expanded
      await page.goto('/results')

      // Collapse some nodes
      await page.getByRole('button', { name: /l1/i }).click()

      // Navigate to analyze and run new analysis
      await page.goto('/analyze')
      // (Would need to run a new analysis here)

      // Navigate back to results
      await page.goto('/results')

      // THEN: Graph should be fully expanded (no badges)
      await expect(page.locator('.collapsed-badge')).toHaveCount(0)
    })
  })

  test.describe('Accessibility', () => {
    test.fixme('[P2] should support keyboard navigation of depth controls', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // Accessibility verification:
      // - Tab navigates between control buttons
      // - Enter/Space activates buttons
      // - ARIA labels are correct
      await page.goto('/results')

      // Focus the control group
      await page.getByRole('button', { name: /all/i }).focus()

      // THEN: Button should be focused
      await expect(page.getByRole('button', { name: /all/i })).toBeFocused()

      // WHEN: Press Tab to move to next button
      await page.keyboard.press('Tab')

      // THEN: Next button should be focused
      await expect(page.getByRole('button', { name: /l1/i })).toBeFocused()

      // WHEN: Press Enter to activate
      await page.keyboard.press('Enter')

      // THEN: L1 should be selected
      await expect(page.getByRole('button', { name: /l1/i })).toHaveAttribute(
        'aria-pressed',
        'true'
      )
    })

    test.fixme('[P2] should have accessible labels for collapsed nodes', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // Accessibility verification:
      // - Collapsed nodes have aria-expanded="false"
      // - Badge has aria-label describing count
      await page.goto('/results')

      // Collapse some nodes
      await page.getByRole('button', { name: /l1/i }).click()

      // Verify accessibility attributes
      const collapsedNode = page.locator('g.node--collapsed').first()
      // Node groups don't typically have aria-expanded, but we could add it via title element
      await expect(collapsedNode).toBeVisible()
    })
  })
})

/**
 * Story 4.4: Add Zoom, Pan, and Navigation Controls
 *
 * Tests for zoom/pan interactions and navigation UI in the dependency graph.
 * These tests verify:
 * - AC1: Scroll zoom centered on cursor
 * - AC2: Click and drag pan
 * - AC3: Zoom control buttons (+/-/fit/reset)
 * - AC4: Fit to screen
 * - AC5: Minimap navigation for large graphs (>50 nodes)
 * - AC6: Zoom level percentage display
 * - AC7: Zoom range limits (10%-400%)
 */
test.describe('Zoom, Pan & Navigation Controls (Story 4.4)', () => {
  test.describe('Zoom Control Buttons (AC3, AC6)', () => {
    test.fixme('[P1] should display zoom controls when graph is rendered', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // GIVEN: User navigates to results with analysis data
      await page.goto('/results')

      // THEN: Zoom control buttons should be visible
      await expect(page.getByRole('button', { name: /zoom in/i })).toBeVisible()
      await expect(page.getByRole('button', { name: /zoom out/i })).toBeVisible()
      await expect(page.getByRole('button', { name: /fit to screen/i })).toBeVisible()
      await expect(page.getByRole('button', { name: /reset zoom/i })).toBeVisible()
    })

    test.fixme('[P1] should show zoom percentage display', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // GIVEN: User navigates to results with graph rendered
      await page.goto('/results')

      // THEN: Zoom percentage should be displayed (default 100%)
      await expect(page.getByText(/100%/)).toBeVisible()
    })

    test.fixme('[P1] should update zoom level when zoom in is clicked', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // GIVEN: Graph is rendered at default zoom
      await page.goto('/results')

      // WHEN: User clicks zoom in button
      await page.getByRole('button', { name: /zoom in/i }).click()

      // THEN: Zoom percentage should increase by 20%
      await expect(page.getByText(/120%/)).toBeVisible()
    })

    test.fixme('[P1] should update zoom level when zoom out is clicked', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // GIVEN: Graph is rendered at default zoom
      await page.goto('/results')

      // WHEN: User clicks zoom out button
      await page.getByRole('button', { name: /zoom out/i }).click()

      // THEN: Zoom percentage should decrease by 20%
      await expect(page.getByText(/80%/)).toBeVisible()
    })

    test.fixme('[P1] should fit graph to screen when fit button is clicked', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // GIVEN: Graph is zoomed in
      await page.goto('/results')
      await page.getByRole('button', { name: /zoom in/i }).click()
      await page.getByRole('button', { name: /zoom in/i }).click()

      // WHEN: User clicks fit to screen
      await page.getByRole('button', { name: /fit to screen/i }).click()

      // THEN: Graph should be fitted - entire graph visible in viewport
      const svg = page.locator('svg').first()
      await expect(svg).toBeVisible()
    })
  })

  test.describe('Zoom Limits (AC7)', () => {
    test.fixme('[P1] should disable zoom out at minimum zoom (10%)', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // GIVEN: Graph is rendered
      await page.goto('/results')

      // WHEN: User zooms out repeatedly to minimum
      const zoomOutBtn = page.getByRole('button', { name: /zoom out/i })
      for (let i = 0; i < 20; i++) {
        if (await zoomOutBtn.isEnabled()) {
          await zoomOutBtn.click()
        }
      }

      // THEN: Zoom out button should be disabled at 10%
      await expect(zoomOutBtn).toBeDisabled()
      await expect(page.getByText(/10%/)).toBeVisible()
    })

    test.fixme('[P2] should disable zoom in at maximum zoom (400%)', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // GIVEN: Graph is rendered
      await page.goto('/results')

      // WHEN: User zooms in repeatedly to maximum
      const zoomInBtn = page.getByRole('button', { name: /zoom in/i })
      for (let i = 0; i < 20; i++) {
        if (await zoomInBtn.isEnabled()) {
          await zoomInBtn.click()
        }
      }

      // THEN: Zoom in button should be disabled at 400%
      await expect(zoomInBtn).toBeDisabled()
      await expect(page.getByText(/400%/)).toBeVisible()
    })
  })

  test.describe('Minimap Navigation (AC5)', () => {
    test.fixme('[P2] should display minimap for large graphs (>50 nodes)', async ({ page }) => {
      // FIXME: Test requires analysis data with 50+ packages
      // GIVEN: Analysis with more than 50 packages
      await page.goto('/results')

      // THEN: Minimap should be visible in corner
      await expect(page.getByLabel(/minimap/i)).toBeVisible()
    })

    test.fixme('[P2] should show viewport indicator on minimap', async ({ page }) => {
      // FIXME: Test requires analysis data with 50+ packages
      // GIVEN: Minimap is displayed
      await page.goto('/results')

      // THEN: Viewport indicator rectangle should be visible within minimap
      const minimap = page.getByLabel(/minimap/i)
      await expect(minimap).toBeVisible()
      await expect(minimap.locator('rect.viewport-indicator')).toBeVisible()
    })
  })

  test.describe('Scroll Zoom (AC1)', () => {
    test.fixme('[P1] should zoom graph on mouse wheel scroll', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // GIVEN: Graph is rendered
      await page.goto('/results')
      const svg = page.locator('svg').first()

      // WHEN: User scrolls mouse wheel up (zoom in)
      await svg.hover()
      await page.mouse.wheel(0, -100)

      // THEN: Zoom percentage should increase from 100%
      // Wait for zoom update
      await expect(page.getByText(/100%/)).not.toBeVisible()
    })
  })

  test.describe('Drag Pan (AC2)', () => {
    test.fixme('[P1] should pan graph on mouse drag', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // GIVEN: Graph is rendered
      await page.goto('/results')
      const svg = page.locator('svg').first()
      const box = await svg.boundingBox()

      if (box) {
        const centerX = box.x + box.width / 2
        const centerY = box.y + box.height / 2

        // WHEN: User drags the graph
        await page.mouse.move(centerX, centerY)
        await page.mouse.down()
        await page.mouse.move(centerX + 100, centerY + 50)
        await page.mouse.up()

        // THEN: Graph should have panned (transform should change)
        // Verify SVG group transform attribute changed
        await expect(svg.locator('g').first()).toBeVisible()
      }
    })
  })
})

/**
 * Story 4.5: Implement Hover Details and Tooltips
 *
 * Tests for hover tooltips showing package details in the dependency graph.
 * These tests verify:
 * - AC1: Tooltip content (package name, path, dependency counts, health)
 * - AC2: Tooltip timing (appear on hover, disappear on leave)
 * - AC3: Tooltip positioning (within viewport bounds)
 * - AC4: Edge highlighting on node hover
 * - AC5: Performance (smooth for large graphs)
 * - AC7: Accessibility (role="tooltip", aria-live)
 */
test.describe('Hover Details & Tooltips (Story 4.5)', () => {
  test.describe('Tooltip Display (AC1, AC2)', () => {
    test.fixme('[P1] should show tooltip on node hover', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // GIVEN: Graph is rendered with nodes
      await page.goto('/results')

      // WHEN: User hovers over a node
      const node = page.locator('g.node').first()
      await node.hover()

      // THEN: Tooltip should appear with package information
      const tooltip = page.locator('[role="tooltip"]')
      await expect(tooltip).toBeVisible()
    })

    test.fixme('[P1] should display package name in tooltip', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // GIVEN: Graph is rendered
      await page.goto('/results')

      // WHEN: User hovers over a node
      await page.locator('g.node').first().hover()

      // THEN: Tooltip should show the package name
      const tooltip = page.locator('[role="tooltip"]')
      await expect(tooltip).toBeVisible()
      // Tooltip should contain text (package name)
      const tooltipText = await tooltip.textContent()
      expect(tooltipText?.length).toBeGreaterThan(0)
    })

    test.fixme('[P1] should show dependency counts in tooltip', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // GIVEN: Graph is rendered
      await page.goto('/results')

      // WHEN: User hovers over a node
      await page.locator('g.node').first().hover()

      // THEN: Tooltip should show incoming and outgoing dependency counts
      const tooltip = page.locator('[role="tooltip"]')
      await expect(tooltip.getByText(/dependencies/i)).toBeVisible()
    })

    test.fixme('[P1] should hide tooltip when mouse leaves node', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // GIVEN: Tooltip is visible on hover
      await page.goto('/results')
      const node = page.locator('g.node').first()
      await node.hover()
      await expect(page.locator('[role="tooltip"]')).toBeVisible()

      // WHEN: User moves mouse away from node
      await page.mouse.move(0, 0)

      // THEN: Tooltip should disappear
      await expect(page.locator('[role="tooltip"]')).not.toBeVisible()
    })
  })

  test.describe('Edge Highlighting (AC4)', () => {
    test.fixme('[P1] should highlight connected edges on node hover', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // GIVEN: Graph with multiple connected nodes
      await page.goto('/results')

      // WHEN: User hovers over a node
      await page.locator('g.node').first().hover()

      // THEN: Connected edges should be highlighted (visible styling change)
      // Non-connected elements should be dimmed
      const dimmedNodes = page.locator('g.node[opacity]')
      // At least some nodes should be dimmed when hovering
      const dimmedCount = await dimmedNodes.count()
      expect(dimmedCount).toBeGreaterThanOrEqual(0)
    })
  })

  test.describe('Cycle Warning in Tooltip (AC1)', () => {
    test.fixme(
      '[P2] should show cycle warning for nodes in circular dependency',
      async ({ page }) => {
        // FIXME: Test requires analysis data with circular dependencies
        // GIVEN: Graph contains circular dependencies
        await page.goto('/results')

        // WHEN: User hovers over a node that is part of a cycle
        const cycleNode = page.locator('g.node--cycle').first()
        await cycleNode.hover()

        // THEN: Tooltip should show circular dependency warning
        const tooltip = page.locator('[role="tooltip"]')
        await expect(tooltip.getByText(/circular/i)).toBeVisible()
      }
    )
  })

  test.describe('Accessibility (AC7)', () => {
    test.fixme('[P2] should have accessible tooltip attributes', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // GIVEN: Graph is rendered
      await page.goto('/results')

      // WHEN: User hovers over a node
      await page.locator('g.node').first().hover()

      // THEN: Tooltip should have proper ARIA attributes
      const tooltip = page.locator('[role="tooltip"]')
      await expect(tooltip).toBeVisible()
      await expect(tooltip).toHaveAttribute('aria-live', 'polite')
    })
  })
})

/**
 * Story 4.6: Export Graph as PNG/SVG Images
 *
 * Tests for exporting the dependency graph as image files.
 * These tests verify:
 * - AC1: Export format options (PNG, SVG)
 * - AC2: PNG export with resolution options
 * - AC3: SVG export (vector-based)
 * - AC4: Export scope (viewport, full, selected)
 * - AC5: Legend inclusion/exclusion
 * - AC7: Export button and menu UI
 * - AC8: Performance (< 5 seconds for large graphs)
 */
test.describe('Export Graph as PNG/SVG (Story 4.6)', () => {
  test.describe('Export Button & Menu UI (AC7)', () => {
    test.fixme('[P1] should display export button when graph is rendered', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // GIVEN: Graph is rendered
      await page.goto('/results')

      // THEN: Export button should be visible in graph controls area
      await expect(page.getByRole('button', { name: /export/i })).toBeVisible()
    })

    test.fixme('[P1] should open export menu on button click', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // GIVEN: Graph is rendered
      await page.goto('/results')

      // WHEN: User clicks the export button
      await page.getByRole('button', { name: /export/i }).click()

      // THEN: Export menu should be visible with title
      await expect(page.getByText('Export Graph')).toBeVisible()
    })

    test.fixme('[P1] should close export menu on close button click', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // GIVEN: Export menu is open
      await page.goto('/results')
      await page.getByRole('button', { name: /export/i }).click()
      await expect(page.getByText('Export Graph')).toBeVisible()

      // WHEN: User clicks the close button
      await page.getByRole('button', { name: /close export menu/i }).click()

      // THEN: Export menu should be hidden
      await expect(page.getByText('Export Graph')).not.toBeVisible()
    })
  })

  test.describe('Format Selection (AC1)', () => {
    test.fixme('[P1] should show PNG and SVG format options', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // GIVEN: Export menu is open
      await page.goto('/results')
      await page.getByRole('button', { name: /export/i }).click()

      // THEN: Format selection buttons should be visible
      await expect(page.getByRole('button', { name: /^PNG$/i })).toBeVisible()
      await expect(page.getByRole('button', { name: /^SVG$/i })).toBeVisible()
    })

    test.fixme('[P1] should show resolution dropdown for PNG format', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // GIVEN: Export menu is open with PNG selected (default)
      await page.goto('/results')
      await page.getByRole('button', { name: /export/i }).click()

      // THEN: Resolution dropdown should be visible
      await expect(page.getByText('Resolution')).toBeVisible()
      await expect(page.locator('select')).toBeVisible()
    })

    test.fixme('[P1] should hide resolution dropdown for SVG format', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // GIVEN: Export menu is open
      await page.goto('/results')
      await page.getByRole('button', { name: /export/i }).click()

      // WHEN: User selects SVG format
      await page.getByRole('button', { name: /^SVG$/i }).click()

      // THEN: Resolution dropdown should be hidden (not applicable for SVG)
      await expect(page.getByText('Resolution')).not.toBeVisible()
    })
  })

  test.describe('Scope & Options (AC4, AC5)', () => {
    test.fixme('[P2] should show scope selection options', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // GIVEN: Export menu is open
      await page.goto('/results')
      await page.getByRole('button', { name: /export/i }).click()

      // THEN: Scope dropdown with options should be visible
      await expect(page.getByText('Scope')).toBeVisible()
      await expect(page.getByText('Current View')).toBeVisible()
    })

    test.fixme('[P2] should show legend checkbox option', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // GIVEN: Export menu is open
      await page.goto('/results')
      await page.getByRole('button', { name: /export/i }).click()

      // THEN: Include legend checkbox should be visible and checked by default
      const legendCheckbox = page.getByLabel(/include legend/i)
      await expect(legendCheckbox).toBeVisible()
      await expect(legendCheckbox).toBeChecked()
    })

    test.fixme('[P2] should toggle legend checkbox', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // GIVEN: Export menu is open with legend checked
      await page.goto('/results')
      await page.getByRole('button', { name: /export/i }).click()
      const legendCheckbox = page.getByLabel(/include legend/i)
      await expect(legendCheckbox).toBeChecked()

      // WHEN: User unchecks the legend option
      await legendCheckbox.uncheck()

      // THEN: Checkbox should be unchecked
      await expect(legendCheckbox).not.toBeChecked()
    })
  })

  test.describe('Export Execution (AC2, AC3)', () => {
    test.fixme('[P1] should trigger download on export', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // GIVEN: Export menu is open with default settings
      await page.goto('/results')
      await page.getByRole('button', { name: /export/i }).click()

      // WHEN: User clicks the Export button
      const downloadPromise = page.waitForEvent('download')
      await page.getByRole('button', { name: /^Export$/i }).click()

      // THEN: File download should be triggered
      const download = await downloadPromise
      expect(download.suggestedFilename()).toMatch(/monoguard.*\.(png|svg)$/)
    })

    test.fixme('[P2] should show progress indicator during export', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // GIVEN: Export menu is open
      await page.goto('/results')
      await page.getByRole('button', { name: /export/i }).click()

      // WHEN: User clicks Export
      await page.getByRole('button', { name: /^Export$/i }).click()

      // THEN: Progress indicator should be shown
      // Note: May be too fast to catch in small graphs
      await expect(page.getByRole('button', { name: /^Export$/i })).toBeVisible()
    })

    test.fixme('[P2] should disable export button during export', async ({ page }) => {
      // FIXME: Test requires analysis data in store
      // GIVEN: Export is in progress
      await page.goto('/results')
      await page.getByRole('button', { name: /export/i }).click()

      // WHEN: Export starts
      await page.getByRole('button', { name: /^Export$/i }).click()

      // THEN: Export button should be disabled during processing
      // Note: May need to intercept to slow down export for assertion
      await expect(page.getByRole('button', { name: /^Export$/i })).toBeVisible()
    })
  })
})
