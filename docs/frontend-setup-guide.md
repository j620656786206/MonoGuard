# MonoGuard 前端設定指南

## 專案結構建議

本指南提供設定 MonoGuard 前端專案結構的詳細建議，包含目錄組織、設定檔案和開發工作流程設定。

## 建議的專案結構

```
mono-guard/
├── frontend/                           # 前端應用程式根目錄
│   ├── .env.local                     # 本機環境變數
│   ├── .env.example                   # 環境變數範本
│   ├── .eslintrc.json                 # ESLint 設定
│   ├── .gitignore                     # Git 忽略模式
│   ├── .prettierrc                    # Prettier 設定
│   ├── components.json                # Shadcn/ui 設定
│   ├── next.config.js                 # Next.js 設定
│   ├── package.json                   # 相依性與腳本
│   ├── playwright.config.ts           # 端對端測試設定
│   ├── postcss.config.js              # PostCSS 設定
│   ├── tailwind.config.ts             # Tailwind CSS 設定
│   ├── tsconfig.json                  # TypeScript 設定
│   ├── jest.config.js                 # Jest 測試設定
│   ├── jest.setup.js                  # Jest 設定檔
│   ├── src/                           # 原始碼目錄
│   │   ├── app/                       # Next.js App Router
│   │   │   ├── (auth)/               # 認證頁面路由群組
│   │   │   │   ├── login/
│   │   │   │   │   └── page.tsx
│   │   │   │   └── register/
│   │   │   │       └── page.tsx
│   │   │   ├── dashboard/            # 儀表板頁面
│   │   │   │   ├── loading.tsx       # 載入中 UI
│   │   │   │   ├── error.tsx         # 錯誤 UI
│   │   │   │   └── page.tsx          # 儀表板主頁面
│   │   │   ├── dependencies/         # 依賴分析頁面
│   │   │   │   ├── loading.tsx
│   │   │   │   ├── error.tsx
│   │   │   │   └── page.tsx
│   │   │   ├── architecture/         # 架構驗證頁面
│   │   │   │   ├── loading.tsx
│   │   │   │   ├── error.tsx
│   │   │   │   └── page.tsx
│   │   │   ├── projects/             # 專案管理頁面
│   │   │   │   ├── [id]/
│   │   │   │   │   ├── page.tsx     # Project detail page
│   │   │   │   │   └── settings/
│   │   │   │   │       └── page.tsx
│   │   │   │   ├── new/
│   │   │   │   │   └── page.tsx     # Create project page
│   │   │   │   └── page.tsx         # Projects list page
│   │   │   ├── globals.css          # Global CSS and Tailwind imports
│   │   │   ├── layout.tsx           # Root layout component
│   │   │   ├── loading.tsx          # Global loading UI
│   │   │   ├── error.tsx            # Global error UI
│   │   │   ├── not-found.tsx        # 404 page
│   │   │   └── page.tsx             # Home page
│   │   ├── components/              # React components
│   │   │   ├── ui/                  # Base UI components (shadcn/ui)
│   │   │   │   ├── alert.tsx
│   │   │   │   ├── button.tsx
│   │   │   │   ├── card.tsx
│   │   │   │   ├── dialog.tsx
│   │   │   │   ├── dropdown-menu.tsx
│   │   │   │   ├── input.tsx
│   │   │   │   ├── label.tsx
│   │   │   │   ├── select.tsx
│   │   │   │   ├── table.tsx
│   │   │   │   ├── toast.tsx
│   │   │   │   └── index.ts         # Barrel exports
│   │   │   ├── dashboard/           # Dashboard-specific components
│   │   │   │   ├── HealthScoreCard.tsx
│   │   │   │   ├── TrendChart.tsx
│   │   │   │   ├── IssuesSummary.tsx
│   │   │   │   ├── QuickStats.tsx
│   │   │   │   └── RecentActivity.tsx
│   │   │   ├── dependency/          # Dependency analysis components
│   │   │   │   ├── DependencyGraph.tsx
│   │   │   │   ├── DuplicatesList.tsx
│   │   │   │   ├── ConflictResolver.tsx
│   │   │   │   ├── VersionMatrix.tsx
│   │   │   │   └── BundleAnalysis.tsx
│   │   │   ├── architecture/        # Architecture validation components
│   │   │   │   ├── LayerDiagram.tsx
│   │   │   │   ├── ViolationsList.tsx
│   │   │   │   ├── RuleEditor.tsx
│   │   │   │   ├── CircularDependencies.tsx
│   │   │   │   └── ArchitectureConfig.tsx
│   │   │   ├── project/            # Project management components
│   │   │   │   ├── ProjectCard.tsx
│   │   │   │   ├── ProjectSelector.tsx
│   │   │   │   ├── CreateProjectForm.tsx
│   │   │   │   └── ProjectSettings.tsx
│   │   │   ├── common/             # Shared components
│   │   │   │   ├── Header.tsx
│   │   │   │   ├── Sidebar.tsx
│   │   │   │   ├── Navigation.tsx
│   │   │   │   ├── UserMenu.tsx
│   │   │   │   ├── NotificationBell.tsx
│   │   │   │   ├── LoadingSpinner.tsx
│   │   │   │   ├── ErrorBoundary.tsx
│   │   │   │   ├── EmptyState.tsx
│   │   │   │   ├── ConfirmDialog.tsx
│   │   │   │   └── DataTable.tsx
│   │   │   └── forms/              # Form components
│   │   │       ├── FormField.tsx
│   │   │       ├── FormError.tsx
│   │   │       └── FormSubmit.tsx
│   │   ├── hooks/                  # Custom React hooks
│   │   │   ├── api/               # API-related hooks
│   │   │   │   ├── useProjects.ts
│   │   │   │   ├── useAnalysis.ts
│   │   │   │   ├── useDependencies.ts
│   │   │   │   └── useArchitecture.ts
│   │   │   ├── ui/                # UI-related hooks
│   │   │   │   ├── useToast.ts
│   │   │   │   ├── useModal.ts
│   │   │   │   ├── useLocalStorage.ts
│   │   │   │   └── useDebounce.ts
│   │   │   ├── utils/             # Utility hooks
│   │   │   │   ├── useWebSocket.ts
│   │   │   │   ├── useInterval.ts
│   │   │   │   └── usePrevious.ts
│   │   │   └── business/          # Business logic hooks
│   │   │       ├── useHealthScore.ts
│   │   │       ├── useIssueFilters.ts
│   │   │       └── useGraphLayout.ts
│   │   ├── lib/                   # Utility libraries and configurations
│   │   │   ├── api/              # API client and utilities
│   │   │   │   ├── client.ts     # Main API client
│   │   │   │   ├── types.ts      # API type definitions
│   │   │   │   ├── endpoints.ts  # API endpoint constants
│   │   │   │   ├── auth.ts       # Authentication utilities
│   │   │   │   └── websocket.ts  # WebSocket client
│   │   │   ├── utils/            # General utility functions
│   │   │   │   ├── cn.ts         # Class name utility
│   │   │   │   ├── format.ts     # Data formatting utilities
│   │   │   │   ├── validation.ts # Validation utilities
│   │   │   │   ├── date.ts       # Date utilities
│   │   │   │   ├── file.ts       # File utilities
│   │   │   │   └── color.ts      # Color utilities
│   │   │   ├── constants/        # Application constants
│   │   │   │   ├── routes.ts     # Route definitions
│   │   │   │   ├── config.ts     # App configuration
│   │   │   │   ├── colors.ts     # Color scheme
│   │   │   │   ├── sizes.ts      # Size constants
│   │   │   │   └── messages.ts   # UI messages
│   │   │   ├── d3/              # D3.js utilities and configurations
│   │   │   │   ├── graph-layout.ts
│   │   │   │   ├── force-simulation.ts
│   │   │   │   ├── svg-utils.ts
│   │   │   │   └── data-transforms.ts
│   │   │   ├── chart/           # Chart.js utilities
│   │   │   │   ├── config.ts
│   │   │   │   ├── themes.ts
│   │   │   │   └── plugins.ts
│   │   │   └── auth/            # Authentication utilities
│   │   │       ├── providers.ts
│   │   │       ├── middleware.ts
│   │   │       └── config.ts
│   │   ├── store/               # State management (Zustand stores)
│   │   │   ├── auth.ts         # Authentication state
│   │   │   ├── project.ts      # Project state
│   │   │   ├── ui.ts           # UI state (sidebar, theme, etc.)
│   │   │   ├── analysis.ts     # Analysis state
│   │   │   ├── notifications.ts # Notification state
│   │   │   └── index.ts        # Store exports and providers
│   │   ├── types/              # TypeScript type definitions
│   │   │   ├── api.ts         # API response types
│   │   │   ├── domain.ts      # Domain model types
│   │   │   ├── components.ts  # Component prop types
│   │   │   ├── auth.ts        # Authentication types
│   │   │   ├── ui.ts          # UI-related types
│   │   │   └── index.ts       # Type exports
│   │   └── styles/            # Styling files
│   │       ├── globals.css    # Global styles
│   │       ├── components.css # Component-specific styles
│   │       └── themes/        # Theme files
│   │           ├── light.css
│   │           └── dark.css
│   ├── public/                # Static assets
│   │   ├── images/           # Image assets
│   │   │   ├── logo.svg
│   │   │   ├── hero.png
│   │   │   └── placeholders/
│   │   ├── icons/            # Icon assets
│   │   │   ├── favicon.ico
│   │   │   ├── apple-touch-icon.png
│   │   │   └── manifest-icons/
│   │   ├── fonts/            # Custom fonts (if any)
│   │   └── manifest.json     # PWA manifest
│   ├── tests/                # Test files
│   │   ├── __mocks__/       # Jest mocks
│   │   │   ├── next-router.js
│   │   │   ├── api-client.ts
│   │   │   └── d3-modules.js
│   │   ├── components/      # Component tests
│   │   │   ├── dashboard/
│   │   │   ├── dependency/
│   │   │   ├── architecture/
│   │   │   └── common/
│   │   ├── hooks/          # Hook tests
│   │   │   ├── api/
│   │   │   ├── ui/
│   │   │   └── utils/
│   │   ├── lib/            # Library tests
│   │   │   ├── api/
│   │   │   ├── utils/
│   │   │   └── d3/
│   │   ├── store/          # Store tests
│   │   ├── e2e/            # Playwright E2E tests
│   │   │   ├── auth.spec.ts
│   │   │   ├── dashboard.spec.ts
│   │   │   ├── dependencies.spec.ts
│   │   │   └── architecture.spec.ts
│   │   ├── fixtures/       # Test data fixtures
│   │   │   ├── projects.json
│   │   │   ├── analysis.json
│   │   │   └── dependencies.json
│   │   └── utils/          # Test utilities
│   │       ├── render.tsx  # Custom render function
│   │       ├── mocks.ts    # Mock data generators
│   │       └── setup.ts    # Test setup utilities
│   └── docs/               # Frontend documentation
│       ├── components.md   # Component documentation
│       ├── hooks.md        # Hooks documentation
│       ├── testing.md      # Testing guide
│       ├── deployment.md   # Deployment guide
│       └── performance.md  # Performance guide
```

## 設定檔案建置

### 1. Package.json
```json
{
  "name": "@monoguard/frontend",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "lint": "next lint --fix",
    "lint:check": "next lint",
    "type-check": "tsc --noEmit",
    "test": "jest",
    "test:watch": "jest --watch",
    "test:coverage": "jest --coverage",
    "test:e2e": "playwright test",
    "test:e2e:ui": "playwright test --ui",
    "analyze": "ANALYZE=true npm run build",
    "storybook": "storybook dev -p 6006",
    "build-storybook": "storybook build"
  },
  "dependencies": {
    "next": "^14.0.0",
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "typescript": "^5.2.0",
    "@types/react": "^18.2.0",
    "@types/react-dom": "^18.2.0",
    
    "tailwindcss": "^3.3.0",
    "autoprefixer": "^10.4.0",
    "postcss": "^8.4.0",
    "@tailwindcss/forms": "^0.5.0",
    "@tailwindcss/typography": "^0.5.0",
    
    "@radix-ui/react-alert-dialog": "^1.0.5",
    "@radix-ui/react-button": "^1.0.0",
    "@radix-ui/react-card": "^1.0.0",
    "@radix-ui/react-dialog": "^1.0.5",
    "@radix-ui/react-dropdown-menu": "^2.0.6",
    "@radix-ui/react-label": "^2.0.2",
    "@radix-ui/react-select": "^2.0.0",
    "@radix-ui/react-slot": "^1.0.2",
    "class-variance-authority": "^0.7.0",
    "clsx": "^2.0.0",
    "tailwind-merge": "^2.0.0",
    "lucide-react": "^0.294.0",
    
    "zustand": "^4.4.0",
    "@tanstack/react-query": "^5.0.0",
    "@tanstack/react-query-devtools": "^5.0.0",
    
    "d3": "^7.8.0",
    "chart.js": "^4.4.0",
    "react-chartjs-2": "^5.2.0",
    
    "axios": "^1.6.0",
    "zod": "^3.22.0",
    "react-hook-form": "^7.47.0",
    "@hookform/resolvers": "^3.3.0",
    
    "date-fns": "^2.30.0",
    "use-debounce": "^10.0.0"
  },
  "devDependencies": {
    "@types/d3": "^7.4.0",
    "@types/node": "^20.8.0",
    
    "eslint": "^8.52.0",
    "eslint-config-next": "^14.0.0",
    "@typescript-eslint/eslint-plugin": "^6.9.0",
    "@typescript-eslint/parser": "^6.9.0",
    "eslint-plugin-react": "^7.33.0",
    "eslint-plugin-react-hooks": "^4.6.0",
    "eslint-plugin-jsx-a11y": "^6.8.0",
    
    "prettier": "^3.0.0",
    "prettier-plugin-tailwindcss": "^0.5.0",
    
    "jest": "^29.7.0",
    "jest-environment-jsdom": "^29.7.0",
    "@testing-library/react": "^13.4.0",
    "@testing-library/jest-dom": "^6.1.0",
    "@testing-library/user-event": "^14.5.0",
    
    "@playwright/test": "^1.40.0",
    
    "@storybook/addon-essentials": "^7.5.0",
    "@storybook/addon-interactions": "^7.5.0",
    "@storybook/addon-links": "^7.5.0",
    "@storybook/blocks": "^7.5.0",
    "@storybook/nextjs": "^7.5.0",
    "@storybook/react": "^7.5.0",
    "@storybook/testing-library": "^0.2.0",
    "storybook": "^7.5.0"
  }
}
```

### 2. TypeScript Configuration
```json
// tsconfig.json
{
  "compilerOptions": {
    "target": "ES2017",
    "lib": ["dom", "dom.iterable", "ES6"],
    "allowJs": true,
    "skipLibCheck": true,
    "strict": true,
    "noEmit": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "bundler",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "preserve",
    "incremental": true,
    "plugins": [
      {
        "name": "next"
      }
    ],
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"],
      "@/components/*": ["./src/components/*"],
      "@/hooks/*": ["./src/hooks/*"],
      "@/lib/*": ["./src/lib/*"],
      "@/store/*": ["./src/store/*"],
      "@/types/*": ["./src/types/*"],
      "@/styles/*": ["./src/styles/*"]
    }
  },
  "include": [
    "next-env.d.ts",
    "**/*.ts",
    "**/*.tsx",
    ".next/types/**/*.ts"
  ],
  "exclude": [
    "node_modules",
    ".next",
    "out"
  ]
}
```

### 3. Next.js Configuration
```javascript
// next.config.js
/** @type {import('next').NextConfig} */
const nextConfig = {
  experimental: {
    optimizeCss: true,
    scrollRestoration: true,
    typedRoutes: true,
  },
  
  // Enable SWC compiler for better performance
  swcMinify: true,
  
  // Image optimization
  images: {
    formats: ['image/webp', 'image/avif'],
    deviceSizes: [640, 750, 828, 1080, 1200, 1920, 2048, 3840],
    imageSizes: [16, 32, 48, 64, 96, 128, 256, 384],
  },
  
  // Bundle analyzer (for development)
  ...(process.env.ANALYZE === 'true' && {
    webpack(config) {
      const { BundleAnalyzerPlugin } = require('webpack-bundle-analyzer');
      config.plugins.push(
        new BundleAnalyzerPlugin({
          analyzerMode: 'server',
          analyzerPort: 8888,
          openAnalyzer: true,
        })
      );
      return config;
    },
  }),
  
  // Security headers
  async headers() {
    return [
      {
        source: '/(.*)',
        headers: [
          {
            key: 'X-Frame-Options',
            value: 'DENY',
          },
          {
            key: 'X-Content-Type-Options',
            value: 'nosniff',
          },
          {
            key: 'Referrer-Policy',
            value: 'origin-when-cross-origin',
          },
          {
            key: 'Content-Security-Policy',
            value: [
              "default-src 'self'",
              "script-src 'self' 'unsafe-eval' 'unsafe-inline'",
              "style-src 'self' 'unsafe-inline'",
              "img-src 'self' data: https:",
              "font-src 'self'",
              "connect-src 'self' ws: wss:",
            ].join('; '),
          },
        ],
      },
    ];
  },
  
  // Redirects
  async redirects() {
    return [
      {
        source: '/',
        destination: '/dashboard',
        permanent: false,
      },
    ];
  },
};

module.exports = nextConfig;
```

### 4. Tailwind CSS Configuration
```typescript
// tailwind.config.ts
import type { Config } from 'tailwindcss';

const config: Config = {
  darkMode: ['class'],
  content: [
    './pages/**/*.{ts,tsx}',
    './components/**/*.{ts,tsx}',
    './app/**/*.{ts,tsx}',
    './src/**/*.{ts,tsx}',
  ],
  theme: {
    container: {
      center: true,
      padding: '2rem',
      screens: {
        '2xl': '1400px',
      },
    },
    extend: {
      colors: {
        border: 'hsl(var(--border))',
        input: 'hsl(var(--input))',
        ring: 'hsl(var(--ring))',
        background: 'hsl(var(--background))',
        foreground: 'hsl(var(--foreground))',
        primary: {
          DEFAULT: 'hsl(var(--primary))',
          foreground: 'hsl(var(--primary-foreground))',
        },
        secondary: {
          DEFAULT: 'hsl(var(--secondary))',
          foreground: 'hsl(var(--secondary-foreground))',
        },
        destructive: {
          DEFAULT: 'hsl(var(--destructive))',
          foreground: 'hsl(var(--destructive-foreground))',
        },
        muted: {
          DEFAULT: 'hsl(var(--muted))',
          foreground: 'hsl(var(--muted-foreground))',
        },
        accent: {
          DEFAULT: 'hsl(var(--accent))',
          foreground: 'hsl(var(--accent-foreground))',
        },
        popover: {
          DEFAULT: 'hsl(var(--popover))',
          foreground: 'hsl(var(--popover-foreground))',
        },
        card: {
          DEFAULT: 'hsl(var(--card))',
          foreground: 'hsl(var(--card-foreground))',
        },
        // MonoGuard brand colors
        brand: {
          50: '#eff6ff',
          100: '#dbeafe',
          200: '#bfdbfe',
          300: '#93c5fd',
          400: '#60a5fa',
          500: '#3b82f6',
          600: '#2563eb',
          700: '#1d4ed8',
          800: '#1e40af',
          900: '#1e3a8a',
          950: '#172554',
        },
        // Semantic colors
        success: {
          50: '#f0fdf4',
          100: '#dcfce7',
          200: '#bbf7d0',
          300: '#86efac',
          400: '#4ade80',
          500: '#22c55e',
          600: '#16a34a',
          700: '#15803d',
          800: '#166534',
          900: '#14532d',
          950: '#052e16',
        },
        warning: {
          50: '#fffbeb',
          100: '#fef3c7',
          200: '#fde68a',
          300: '#fcd34d',
          400: '#fbbf24',
          500: '#f59e0b',
          600: '#d97706',
          700: '#b45309',
          800: '#92400e',
          900: '#78350f',
          950: '#451a03',
        },
        error: {
          50: '#fef2f2',
          100: '#fee2e2',
          200: '#fecaca',
          300: '#fca5a5',
          400: '#f87171',
          500: '#ef4444',
          600: '#dc2626',
          700: '#b91c1c',
          800: '#991b1b',
          900: '#7f1d1d',
          950: '#450a0a',
        },
      },
      borderRadius: {
        lg: 'var(--radius)',
        md: 'calc(var(--radius) - 2px)',
        sm: 'calc(var(--radius) - 4px)',
      },
      fontFamily: {
        sans: ['var(--font-sans)'],
        mono: ['var(--font-mono)'],
      },
      keyframes: {
        'accordion-down': {
          from: { height: '0' },
          to: { height: 'var(--radix-accordion-content-height)' },
        },
        'accordion-up': {
          from: { height: 'var(--radix-accordion-content-height)' },
          to: { height: '0' },
        },
        'fade-in': {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        'slide-in': {
          '0%': { transform: 'translateX(-100%)' },
          '100%': { transform: 'translateX(0)' },
        },
        'pulse-slow': {
          '0%, 100%': { opacity: '1' },
          '50%': { opacity: '0.5' },
        },
      },
      animation: {
        'accordion-down': 'accordion-down 0.2s ease-out',
        'accordion-up': 'accordion-up 0.2s ease-out',
        'fade-in': 'fade-in 0.5s ease-in-out',
        'slide-in': 'slide-in 0.3s ease-out',
        'pulse-slow': 'pulse-slow 2s ease-in-out infinite',
      },
    },
  },
  plugins: [
    require('tailwindcss-animate'),
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ],
};

export default config;
```

### 5. ESLint Configuration
```json
// .eslintrc.json
{
  "extends": [
    "next/core-web-vitals",
    "@typescript-eslint/recommended",
    "plugin:react/recommended",
    "plugin:react-hooks/recommended",
    "plugin:jsx-a11y/recommended"
  ],
  "parser": "@typescript-eslint/parser",
  "plugins": ["@typescript-eslint", "react", "react-hooks", "jsx-a11y"],
  "rules": {
    "react/react-in-jsx-scope": "off",
    "react/prop-types": "off",
    "@typescript-eslint/explicit-function-return-type": "off",
    "@typescript-eslint/explicit-module-boundary-types": "off",
    "@typescript-eslint/no-unused-vars": ["error", { "argsIgnorePattern": "^_" }],
    "@typescript-eslint/no-explicit-any": "warn",
    "jsx-a11y/anchor-is-valid": [
      "error",
      {
        "components": ["Link"],
        "specialLink": ["hrefLeft", "hrefRight"],
        "aspects": ["invalidHref", "preferButton"]
      }
    ],
    "prefer-const": "error",
    "no-var": "error",
    "object-shorthand": "error",
    "prefer-template": "error"
  },
  "settings": {
    "react": {
      "version": "detect"
    }
  }
}
```

### 6. Jest Configuration
```javascript
// jest.config.js
const nextJest = require('next/jest');

const createJestConfig = nextJest({
  dir: './',
});

const customJestConfig = {
  setupFilesAfterEnv: ['<rootDir>/jest.setup.js'],
  testEnvironment: 'jsdom',
  moduleNameMapping: {
    '^@/(.*)$': '<rootDir>/src/$1',
  },
  collectCoverageFrom: [
    'src/**/*.{js,jsx,ts,tsx}',
    '!src/**/*.d.ts',
    '!src/**/*.stories.{js,jsx,ts,tsx}',
    '!src/app/**/*.{js,jsx,ts,tsx}', // Exclude Next.js app router files from coverage
  ],
  coverageThreshold: {
    global: {
      branches: 80,
      functions: 80,
      lines: 80,
      statements: 80,
    },
  },
  testMatch: [
    '<rootDir>/tests/**/*.{js,jsx,ts,tsx}',
    '<rootDir>/src/**/__tests__/**/*.{js,jsx,ts,tsx}',
    '<rootDir>/src/**/*.{test,spec}.{js,jsx,ts,tsx}',
  ],
  testPathIgnorePatterns: ['<rootDir>/.next/', '<rootDir>/node_modules/', '<rootDir>/tests/e2e/'],
  transform: {
    '^.+\\.(js|jsx|ts|tsx)$': ['babel-jest', { presets: ['next/babel'] }],
  },
  transformIgnorePatterns: [
    '/node_modules/(?!(d3|d3-.*)/)',
  ],
  moduleFileExtensions: ['ts', 'tsx', 'js', 'jsx'],
  watchPlugins: [
    'jest-watch-typeahead/filename',
    'jest-watch-typeahead/testname',
  ],
};

module.exports = createJestConfig(customJestConfig);
```

## Development Workflow Setup

### 1. Git Hooks (Husky)
```json
// Add to package.json
{
  "devDependencies": {
    "husky": "^8.0.0",
    "lint-staged": "^15.0.0"
  },
  "husky": {
    "hooks": {
      "pre-commit": "lint-staged",
      "pre-push": "npm run type-check && npm run test"
    }
  },
  "lint-staged": {
    "*.{js,jsx,ts,tsx}": [
      "eslint --fix",
      "prettier --write",
      "git add"
    ],
    "*.{json,md,yml,yaml}": [
      "prettier --write",
      "git add"
    ]
  }
}
```

### 2. Environment Variables
```bash
# .env.local
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080/ws
NEXT_PUBLIC_APP_ENV=development

# Authentication
NEXTAUTH_SECRET=your-secret-key
NEXTAUTH_URL=http://localhost:3000

# GitHub OAuth (for development)
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret

# Analytics (optional)
NEXT_PUBLIC_GA_MEASUREMENT_ID=G-XXXXXXXXXX
```

```bash
# .env.example
NEXT_PUBLIC_API_URL=
NEXT_PUBLIC_WS_URL=
NEXT_PUBLIC_APP_ENV=

NEXTAUTH_SECRET=
NEXTAUTH_URL=

GITHUB_CLIENT_ID=
GITHUB_CLIENT_SECRET=

NEXT_PUBLIC_GA_MEASUREMENT_ID=
```

### 3. VS Code Configuration
```json
// .vscode/settings.json
{
  "editor.formatOnSave": true,
  "editor.defaultFormatter": "esbenp.prettier-vscode",
  "editor.codeActionsOnSave": {
    "source.fixAll.eslint": true
  },
  "typescript.preferences.importModuleSpecifier": "relative",
  "emmet.includeLanguages": {
    "typescript": "html",
    "typescriptreact": "html"
  },
  "files.associations": {
    "*.css": "tailwindcss"
  },
  "tailwindCSS.includeLanguages": {
    "typescript": "javascript",
    "typescriptreact": "javascript"
  },
  "tailwindCSS.experimental.classRegex": [
    ["cva\\(([^)]*)\\)", "[\"'`]([^\"'`]*).*?[\"'`]"],
    ["cn\\(([^)]*)\\)", "[\"'`]([^\"'`]*).*?[\"'`]"]
  ]
}
```

## Development Conventions

### 1. Component Development
- Use functional components with TypeScript
- Follow the component template structure from FRONTEND_STANDARDS.md
- Include proper JSDoc comments for complex components
- Export components as named exports
- Use proper prop interfaces with JSDoc descriptions

### 2. State Management
- Use Zustand for global state management
- Create domain-specific stores (auth, project, ui, etc.)
- Use React Query for server state management
- Implement optimistic updates where appropriate

### 3. Styling Approach
- Prefer Tailwind utility classes over custom CSS
- Use the `cn` utility for conditional styling
- Follow the design system color palette
- Implement dark mode support from the beginning

### 4. Testing Strategy
- Write unit tests for all components and hooks
- Aim for 80%+ test coverage
- Use React Testing Library for component testing
- Write E2E tests for critical user journeys with Playwright

### 5. Performance Optimization
- Use Next.js dynamic imports for code splitting
- Implement proper loading states and skeletons
- Optimize images with Next.js Image component
- Use React.memo and useMemo for expensive computations

### 6. Accessibility
- Follow WCAG 2.1 AA guidelines
- Use semantic HTML elements
- Implement proper ARIA attributes
- Ensure keyboard navigation works for all interactive elements
- Maintain minimum color contrast ratios

This project structure and setup provides a solid foundation for building a scalable, maintainable, and performant frontend application that aligns with modern React and Next.js best practices.