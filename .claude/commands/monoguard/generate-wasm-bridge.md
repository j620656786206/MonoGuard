---
description: Generate WASM bridge code template (Go + TypeScript)
argument-hint: <function-name> [description]
---

Generate a complete WASM bridge implementation with Go backend and TypeScript frontend.

**Usage:**

- `/monoguard:generate-wasm-bridge AnalyzeWorkspace "Analyze Nx workspace structure"`
- `/monoguard:generate-wasm-bridge DetectCycles "Detect circular dependencies"`

**What This Generates:**

**1. Go WASM Function (packages/analysis-engine/cmd/wasm/):**

```go
//go:build js && wasm

package main

import (
    "encoding/json"
    "syscall/js"
)

// Result type for WASM responses
type Result struct {
    Data  interface{} `json:"data"`
    Error *Error      `json:"error"`
}

type Error struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

// {{FunctionName}} - {{Description}}
func {{FunctionName}}(this js.Value, args []js.Value) interface{} {
    if len(args) != 1 {
        return toJSON(NewError("INVALID_ARGS", "Expected 1 argument"))
    }

    input := args[0].String()

    // TODO: Implement function logic
    result := map[string]interface{}{
        "status": "success",
    }

    return toJSON(NewSuccess(result))
}

func NewSuccess(data interface{}) Result {
    return Result{Data: data, Error: nil}
}

func NewError(code, message string) Result {
    return Result{
        Data: nil,
        Error: &Error{
            Code:    code,
            Message: message,
        },
    }
}

func toJSON(v interface{}) string {
    b, _ := json.Marshal(v)
    return string(b)
}

func main() {
    c := make(chan struct{})
    js.Global().Set("{{functionName}}", js.FuncOf({{FunctionName}}))
    <-c
}
```

**2. TypeScript Bridge (apps/web/app/lib/wasmBridge.ts):**

```typescript
import { Result } from '@monoguard/types';

export interface {{FunctionName}}Input {
  // TODO: Define input structure
}

export interface {{FunctionName}}Output {
  // TODO: Define output structure
}

export async function {{functionName}}(
  input: {{FunctionName}}Input
): Promise<Result<{{FunctionName}}Output>> {
  if (!window.{{functionName}}) {
    throw new Error('WASM not loaded');
  }

  const jsonInput = JSON.stringify(input);
  const jsonOutput = window.{{functionName}}(jsonInput);
  const result: Result<{{FunctionName}}Output> = JSON.parse(jsonOutput);

  return result;
}

// Type declaration for WASM function
declare global {
  interface Window {
    {{functionName}}: (input: string) => string;
  }
}
```

**3. TypeScript Types (packages/types/src/):**

```typescript
export interface {{FunctionName}}Input {
  // TODO: Define input structure
}

export interface {{FunctionName}}Output {
  status: string;
  // TODO: Add output fields (use camelCase)
}

export type Result<T> = {
  data: T | null;
  error: {
    code: string;
    message: string;
  } | null;
};
```

**4. Unit Test Template (Go):**

```go
func Test{{FunctionName}}(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid input", `{"key": "value"}`, false},
        {"invalid JSON", `{invalid`, true},
        {"empty input", ``, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := {{FunctionName}}(js.Null(), []js.Value{js.ValueOf(tt.input)})
            resultStr := result.(string)

            var res Result
            if err := json.Unmarshal([]byte(resultStr), &res); err != nil {
                t.Fatalf("Failed to parse result: %v", err)
            }

            if (res.Error != nil) != tt.wantErr {
                t.Errorf("wantErr = %v, got error = %v", tt.wantErr, res.Error)
            }
        })
    }
}
```

**5. Integration Test Template (TypeScript):**

```typescript
import { describe, it, expect, vi } from 'vitest';
import { {{functionName}} } from './wasmBridge';

describe('{{functionName}}', () => {
  it('should return success result', async () => {
    // Mock WASM function
    window.{{functionName}} = vi.fn().mockReturnValue(
      JSON.stringify({
        data: { status: 'success' },
        error: null
      })
    );

    const result = await {{functionName}}({ /* input */ });

    expect(result.data).toBeDefined();
    expect(result.error).toBeNull();
  });

  it('should handle WASM errors', async () => {
    window.{{functionName}} = vi.fn().mockReturnValue(
      JSON.stringify({
        data: null,
        error: { code: 'PARSE_ERROR', message: 'Invalid input' }
      })
    );

    const result = await {{functionName}}({ /* input */ });

    expect(result.data).toBeNull();
    expect(result.error?.code).toBe('PARSE_ERROR');
  });
});
```

Let me generate the WASM bridge code for: **$ARGUMENTS**

I'll create all necessary files following MonoGuard's WASM bridge pattern.
