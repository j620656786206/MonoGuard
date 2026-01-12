---
description: Generate React component with tests and D3 cleanup
argument-hint: <component-name> [type: basic|d3|form]
---

Generate a production-ready React component following MonoGuard's patterns.

**Usage:**

- `/monoguard:scaffold-component DependencyGraph d3` - D3.js visualization component
- `/monoguard:scaffold-component AnalysisForm form` - Form component
- `/monoguard:scaffold-component MetricCard basic` - Basic UI component

**Component Types:**

**1. Basic Component (default):**

- Standard React functional component
- TypeScript props interface
- Basic test coverage

**2. D3 Component (d3):**

- D3.js integration with useRef + useEffect
- **Mandatory cleanup** (remove event listeners)
- React.memo() for performance
- SVG/Canvas switch logic (>500 nodes)

**3. Form Component (form):**

- Form state management
- Validation patterns
- Error handling
- Submit handling

---

## Basic Component Template

**Component (apps/web/app/components/{{ComponentName}}.tsx):**

```typescript
import React from 'react';

export interface {{ComponentName}}Props {
  // TODO: Define props
  className?: string;
}

export const {{ComponentName}}: React.FC<{{ComponentName}}Props> = ({
  className,
}) => {
  return (
    <div className={className}>
      {/* TODO: Implement component */}
    </div>
  );
};
```

---

## D3 Component Template

**Component (apps/web/app/components/{{ComponentName}}.tsx):**

```typescript
import React, { useEffect, useRef } from 'react';
import * as d3 from 'd3';

export interface {{ComponentName}}Props {
  data: any[]; // TODO: Define data type
  width?: number;
  height?: number;
}

export const {{ComponentName}} = React.memo<{{ComponentName}}Props>(
  ({ data, width = 800, height = 600 }) => {
    const svgRef = useRef<SVGSVGElement>(null);

    useEffect(() => {
      if (!svgRef.current || !data || data.length === 0) return;

      // Clear previous render
      const svg = d3.select(svgRef.current);
      svg.selectAll('*').remove();

      // Set up SVG
      svg.attr('width', width).attr('height', height);

      // TODO: Implement D3 visualization
      // Example: Draw circles
      svg
        .selectAll('circle')
        .data(data)
        .enter()
        .append('circle')
        .attr('cx', (d, i) => i * 50 + 25)
        .attr('cy', height / 2)
        .attr('r', 20)
        .attr('fill', 'steelblue');

      // CRITICAL: Cleanup function to prevent memory leaks
      return () => {
        svg.selectAll('*').remove();
        svg.on('.zoom', null); // Remove zoom listeners
        svg.on('.drag', null); // Remove drag listeners
        // Remove any other event listeners you added
      };
    }, [data, width, height]);

    // Performance: Switch to Canvas for large datasets
    if (data.length > 500) {
      console.warn('Consider using Canvas for better performance (>500 nodes)');
    }

    return <svg ref={svgRef} />;
  },
  // Custom comparison function (optional)
  (prevProps, nextProps) => {
    return (
      prevProps.data === nextProps.data &&
      prevProps.width === nextProps.width &&
      prevProps.height === nextProps.height
    );
  }
);

{{ComponentName}}.displayName = '{{ComponentName}}';
```

---

## Form Component Template

**Component (apps/web/app/components/{{ComponentName}}.tsx):**

```typescript
import React, { useState } from 'react';
import { AnalysisError } from '@monoguard/types';

export interface {{ComponentName}}Props {
  onSubmit: (data: FormData) => Promise<void>;
  className?: string;
}

interface FormData {
  // TODO: Define form fields
}

export const {{ComponentName}}: React.FC<{{ComponentName}}Props> = ({
  onSubmit,
  className,
}) => {
  const [formData, setFormData] = useState<FormData>({});
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [isSubmitting, setIsSubmitting] = useState(false);

  const validate = (): boolean => {
    const newErrors: Record<string, string> = {};

    // TODO: Add validation rules
    // if (!formData.field) {
    //   newErrors.field = 'Field is required';
    // }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validate()) return;

    setIsSubmitting(true);
    setErrors({});

    try {
      await onSubmit(formData);
    } catch (error) {
      if (error instanceof AnalysisError) {
        setErrors({ submit: error.userMessage });
      } else {
        setErrors({ submit: 'An unexpected error occurred' });
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className={className}>
      {/* TODO: Add form fields */}

      {errors.submit && (
        <div className="error">{errors.submit}</div>
      )}

      <button type="submit" disabled={isSubmitting}>
        {isSubmitting ? 'Submitting...' : 'Submit'}
      </button>
    </form>
  );
};
```

---

## Test Template (apps/web/app/components/**tests**/{{ComponentName}}.test.tsx)

```typescript
import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { {{ComponentName}} } from '../{{ComponentName}}';

describe('{{ComponentName}}', () => {
  it('should render without crashing', () => {
    render(<{{ComponentName}} />);
    // Add assertions
  });

  it('should handle props correctly', () => {
    const props = {
      // TODO: Add test props
    };

    render(<{{ComponentName}} {...props} />);
    // Add assertions
  });

  // For D3 components: Test cleanup
  it('should cleanup D3 resources on unmount', () => {
    const { unmount } = render(<{{ComponentName}} data={[]} />);

    unmount();

    // Verify no memory leaks (D3 event listeners removed)
    // This is implicit - test that unmount doesn't throw
  });
});
```

Let me generate the React component: **$ARGUMENTS**

I'll create the component file with appropriate patterns based on the component type.
