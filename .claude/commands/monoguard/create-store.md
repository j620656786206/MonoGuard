---
description: Generate Zustand store with devtools + persist middleware
argument-hint: <store-name> [description]
---

Generate a production-ready Zustand store following MonoGuard's state management patterns.

**Usage:**

- `/monoguard:create-store analysis "Manage analysis state and results"`
- `/monoguard:create-store settings "User preferences and settings"`

**What This Generates:**

**Store File (apps/web/app/stores/{{storeName}}.ts):**

```typescript
import { create } from 'zustand';
import { devtools, persist } from 'zustand/middleware';

// State interface
interface {{StoreName}}State {
  // TODO: Define state properties
  isLoading: boolean;
  error: string | null;

  // Actions (verb naming: startX, clearX, updateX)
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  reset: () => void;
}

// Initial state
const initialState = {
  isLoading: false,
  error: null,
};

// Create store with middleware
export const use{{StoreName}}Store = create<{{StoreName}}State>()(
  devtools(
    persist(
      (set) => ({
        ...initialState,

        // Actions
        setLoading: (loading) => set({ isLoading: loading }),

        setError: (error) => set({ error }),

        reset: () => set(initialState),
      }),
      {
        name: '{{storeName}}-storage', // localStorage key
        // Optional: Custom storage or partial persist
        // partialize: (state) => ({ theme: state.theme }), // Only persist specific fields
      }
    ),
    {
      name: '{{StoreName}}Store', // DevTools name
    }
  )
);

// Selectors for optimal re-render performance
export const select{{StoreName}}Loading = (state: {{StoreName}}State) => state.isLoading;
export const select{{StoreName}}Error = (state: {{StoreName}}State) => state.error;
```

**Usage Example (apps/web/app/components/Example.tsx):**

```typescript
import { use{{StoreName}}Store, select{{StoreName}}Loading, select{{StoreName}}Error } from '@/stores/{{storeName}}';

export function Example() {
  // ✅ CORRECT: Use selectors to avoid unnecessary re-renders
  const { isLoading, error } = use{{StoreName}}Store(
    (state) => ({
      isLoading: state.isLoading,
      error: state.error
    })
  );

  // Or use exported selectors
  const isLoading = use{{StoreName}}Store(select{{StoreName}}Loading);

  // ❌ WRONG: This subscribes to entire store (causes re-renders on ANY state change)
  // const store = use{{StoreName}}Store();

  const handleAction = () => {
    use{{StoreName}}Store.getState().setLoading(true);
    // ... perform action
    use{{StoreName}}Store.getState().setLoading(false);
  };

  return (
    <div>
      {isLoading && <Spinner />}
      {error && <ErrorMessage>{error}</ErrorMessage>}
    </div>
  );
}
```

**Test File (apps/web/app/stores/**tests**/{{storeName}}.test.ts):**

```typescript
import { describe, it, expect, beforeEach } from 'vitest';
import { use{{StoreName}}Store } from '../{{storeName}}';
import { renderHook, act } from '@testing-library/react';

describe('use{{StoreName}}Store', () => {
  beforeEach(() => {
    // Reset store before each test
    use{{StoreName}}Store.getState().reset();
  });

  it('should initialize with default state', () => {
    const { result } = renderHook(() => use{{StoreName}}Store());

    expect(result.current.isLoading).toBe(false);
    expect(result.current.error).toBeNull();
  });

  it('should update loading state', () => {
    const { result } = renderHook(() => use{{StoreName}}Store());

    act(() => {
      result.current.setLoading(true);
    });

    expect(result.current.isLoading).toBe(true);
  });

  it('should update error state', () => {
    const { result } = renderHook(() => use{{StoreName}}Store());

    act(() => {
      result.current.setError('Test error');
    });

    expect(result.current.error).toBe('Test error');
  });

  it('should reset to initial state', () => {
    const { result } = renderHook(() => use{{StoreName}}Store());

    act(() => {
      result.current.setLoading(true);
      result.current.setError('Error');
      result.current.reset();
    });

    expect(result.current.isLoading).toBe(false);
    expect(result.current.error).toBeNull();
  });

  it('should only re-render when selected state changes', () => {
    let renderCount = 0;

    const { result, rerender } = renderHook(() => {
      renderCount++;
      return use{{StoreName}}Store((state) => ({ isLoading: state.isLoading }));
    });

    act(() => {
      use{{StoreName}}Store.getState().setError('Error'); // Changes error, not isLoading
    });

    rerender();

    // Should not re-render when unselected state changes
    expect(renderCount).toBe(2); // Initial + 1 rerender (not 3)
  });
});
```

**Key Features:**

✅ **Devtools middleware** - Redux DevTools integration for debugging
✅ **Persist middleware** - Automatic localStorage persistence
✅ **Verb naming** - Actions use clear verb names (setX, startX, clearX)
✅ **Selectors** - Pre-built selectors for optimal performance
✅ **Reset function** - Easily reset to initial state
✅ **TypeScript** - Full type safety
✅ **Test coverage** - Complete unit tests included

Let me generate the Zustand store for: **$ARGUMENTS**

I'll create the store file with all middleware and best practices configured.
