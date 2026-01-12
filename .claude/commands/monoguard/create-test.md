---
description: Generate test file template with WASM/Zustand mocks
argument-hint: <file-to-test> [type: unit|integration|e2e]
---

Generate comprehensive test files with proper mocks for MonoGuard's architecture.

**Usage:**

- `/monoguard:create-test apps/web/app/lib/wasmBridge.ts unit` - Unit test with WASM mocks
- `/monoguard:create-test apps/web/app/components/AnalysisView.tsx integration` - Integration test
- `/monoguard:create-test apps/web-e2e/src/analysis-flow.spec.ts e2e` - E2E test

**Test Types:**

**1. Unit Test (unit):**

- Single function/component testing
- All dependencies mocked
- Fast execution (<1s)

**2. Integration Test (integration):**

- Multiple modules working together
- Minimal mocking
- Acceptable execution time (5-10s)

**3. E2E Test (e2e):**

- Full user flow in real browser
- Playwright-based
- Longer execution time (30s-1min)

---

## Unit Test Template (TypeScript)

**Test File (**tests**/{{fileName}}.test.ts):**

```typescript
import { describe, it, expect, vi, beforeEach } from 'vitest';
import 'fake-indexeddb/auto'; // Mock IndexedDB

// Import function/component to test
// import { functionToTest } from '../{{fileName}}';

describe('{{fileName}}', () => {
  beforeEach(() => {
    // Reset mocks before each test
    vi.clearAllMocks();
  });

  describe('functionName', () => {
    it('should handle success case', () => {
      // Arrange
      const input = {
        /* test data */
      };

      // Act
      const result = functionToTest(input);

      // Assert
      expect(result).toBeDefined();
    });

    it('should handle error case', () => {
      // Arrange
      const invalidInput = {
        /* invalid data */
      };

      // Act & Assert
      expect(() => functionToTest(invalidInput)).toThrow();
    });
  });
});
```

---

## Unit Test with WASM Mock

```typescript
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { analyzeWorkspace } from '../wasmBridge';

// Mock WASM module
vi.mock('@/lib/wasmLoader', () => ({
  MonoGuardAnalyzer: {
    analyze: vi.fn(),
  },
}));

describe('WASM Bridge', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should return success result from WASM', async () => {
    // Mock WASM function to return Result<T> structure
    window.analyzeWorkspace = vi.fn().mockReturnValue(
      JSON.stringify({
        data: {
          healthScore: 85,
          issues: [],
        },
        error: null,
      })
    );

    const result = await analyzeWorkspace({
      projects: { app: {} },
    });

    expect(result.data).toBeDefined();
    expect(result.data?.healthScore).toBe(85);
    expect(result.error).toBeNull();
  });

  it('should handle WASM errors', async () => {
    // Mock WASM function to return error
    window.analyzeWorkspace = vi.fn().mockReturnValue(
      JSON.stringify({
        data: null,
        error: {
          code: 'PARSE_ERROR',
          message: 'Invalid workspace.json',
        },
      })
    );

    const result = await analyzeWorkspace({});

    expect(result.data).toBeNull();
    expect(result.error?.code).toBe('PARSE_ERROR');
  });
});
```

---

## Unit Test with Zustand Mock

```typescript
import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { AnalysisView } from '../AnalysisView';

// Mock Zustand store
vi.mock('@/stores/analysis', () => ({
  useAnalysisStore: vi.fn(() => ({
    result: {
      healthScore: 85,
      issues: [],
    },
    isAnalyzing: false,
    startAnalysis: vi.fn(),
    clearResult: vi.fn(),
  })),
}));

describe('AnalysisView', () => {
  it('should display analysis result', () => {
    render(<AnalysisView />);

    expect(screen.getByText(/health score/i)).toBeInTheDocument();
    expect(screen.getByText('85')).toBeInTheDocument();
  });

  it('should show loading state', () => {
    // Override mock for this test
    vi.mocked(useAnalysisStore).mockReturnValue({
      result: null,
      isAnalyzing: true,
      startAnalysis: vi.fn(),
      clearResult: vi.fn(),
    });

    render(<AnalysisView />);

    expect(screen.getByText(/analyzing/i)).toBeInTheDocument();
  });
});
```

---

## Integration Test Template

```typescript
import { describe, it, expect } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { AnalysisFlow } from '../AnalysisFlow';

describe('Analysis Flow Integration', () => {
  it('should complete full analysis flow', async () => {
    const user = userEvent.setup();

    render(<AnalysisFlow />);

    // 1. Upload file
    const file = new File(['{"projects":{}}'], 'workspace.json', {
      type: 'application/json',
    });
    const input = screen.getByLabelText(/upload/i);
    await user.upload(input, file);

    // 2. Start analysis
    const analyzeButton = screen.getByRole('button', { name: /analyze/i });
    await user.click(analyzeButton);

    // 3. Wait for results
    await waitFor(
      () => {
        expect(screen.getByText(/health score/i)).toBeInTheDocument();
      },
      { timeout: 5000 }
    );

    // 4. Verify visualization
    expect(screen.getByTestId('dependency-graph')).toBeInTheDocument();
  });
});
```

---

## E2E Test Template (Playwright)

**Test File (apps/web-e2e/src/{{testName}}.spec.ts):**

```typescript
import { test, expect } from '@playwright/test';

test.describe('Analysis Flow E2E', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('should complete analysis workflow', async ({ page }) => {
    // 1. Upload workspace.json
    const fileInput = page.locator('input[type="file"]');
    await fileInput.setInputFiles('fixtures/sample-workspace.json');

    // 2. Start analysis
    await page.click('button:has-text("Analyze")');

    // 3. Wait for results
    await expect(page.locator('text=Health Score')).toBeVisible({
      timeout: 10000,
    });

    // 4. Verify visualization rendered
    await expect(page.locator('svg')).toBeVisible();

    // 5. Check results persisted to IndexedDB
    const hasData = await page.evaluate(async () => {
      const db = await indexedDB.open('monoguard');
      return db.objectStoreNames.contains('analyses');
    });
    expect(hasData).toBeTruthy();
  });

  test('should handle upload errors gracefully', async ({ page }) => {
    const fileInput = page.locator('input[type="file"]');
    await fileInput.setInputFiles('fixtures/invalid.json');

    await page.click('button:has-text("Analyze")');

    await expect(page.locator('text=/invalid.*json/i')).toBeVisible();
  });
});
```

---

## Go Unit Test Template

**Test File ({{fileName}}\_test.go):**

```go
package analyzer

import (
    "testing"
)

func TestFunctionName(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:    "valid input",
            input:   `{"projects":{}}`,
            want:    "expected output",
            wantErr: false,
        },
        {
            name:    "invalid JSON",
            input:   `{invalid`,
            wantErr: true,
        },
        {
            name:    "empty input",
            input:   "",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := FunctionName(tt.input)

            if (err != nil) != tt.wantErr {
                t.Errorf("FunctionName() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if !tt.wantErr && got != tt.want {
                t.Errorf("FunctionName() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

Let me generate the test file for: **$ARGUMENTS**

I'll create comprehensive tests with appropriate mocks for the test type.
