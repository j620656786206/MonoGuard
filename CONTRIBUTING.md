# Contributing to MonoGuard

Thank you for your interest in contributing to MonoGuard! This document provides guidelines and information for contributors.

## Code of Conduct

Please be respectful and constructive in all interactions. We're building something together.

## How to Contribute

### Reporting Bugs

1. Search [existing issues](https://github.com/user/monoguard/issues) to avoid duplicates
2. If not found, [open a new issue](https://github.com/user/monoguard/issues/new) with:
   - Clear, descriptive title
   - Steps to reproduce
   - Expected vs actual behavior
   - Environment details (OS, Node version, browser)

### Suggesting Features

1. Check the [Roadmap](README.md#roadmap) for planned features
2. Open an issue with the `enhancement` label
3. Describe the use case and proposed solution

### Pull Requests

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/your-feature`
3. Make your changes
4. Write/update tests as needed
5. Ensure all checks pass: `pnpm lint && pnpm type-check && pnpm test`
6. Commit with conventional commits (see below)
7. Push and open a PR

## Development Setup

### Prerequisites

- Node.js 20+
- pnpm 9+
- Go 1.21+ (for API development)

### Getting Started

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/monoguard.git
cd monoguard

# Install dependencies
pnpm install

# Start development server
pnpm dev:web
```

### Project Structure

```
monoguard/
├── apps/
│   ├── web/          # React web application (Vite + TanStack Router)
│   └── api/          # Go API server (Gin)
├── packages/
│   └── types/        # Shared TypeScript types
├── docs/             # Documentation
└── scripts/          # Build and deployment scripts
```

### Key Commands

```bash
# Development
pnpm dev:web          # Start web app on localhost:3000

# Quality checks
pnpm lint             # Run Biome linter
pnpm type-check       # TypeScript type checking
pnpm test             # Run tests

# Build
pnpm build            # Build all packages
```

## Coding Standards

### TypeScript

- Use TypeScript for all new code
- Prefer `type` over `interface` for object types
- Export types from `@monoguard/types` package

### React

- Use functional components with hooks
- Follow existing component patterns in `apps/web/app/components/`
- Use TanStack Router for routing

### Styling

- Use Tailwind CSS utility classes
- Follow existing design patterns
- Maintain responsive design

### Testing

- Write tests for new features
- Maintain >80% code coverage
- Use existing test factories in `src/__tests__/factories/`

## Commit Convention

We use [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `style`: Formatting, no code change
- `refactor`: Code change that neither fixes a bug nor adds a feature
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

### Examples

```
feat(web): add circular dependency highlighting
fix(api): resolve memory leak in graph traversal
docs: update README with new features
test(web): add integration tests for analysis flow
```

## Pull Request Process

1. Update documentation if needed
2. Add tests for new functionality
3. Ensure CI passes
4. Request review from maintainers
5. Address review feedback
6. Squash and merge once approved

## Questions?

- Open a [Discussion](https://github.com/user/monoguard/discussions) for general questions
- Tag maintainers in issues if blocked

---

Thank you for contributing!
