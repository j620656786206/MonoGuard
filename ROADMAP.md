# MonoGuard Roadmap

[English](ROADMAP.md) | [繁體中文](ROADMAP.zh-TW.md)

This document outlines the planned features and improvements for MonoGuard.

## Legend

- ✅ Completed
- 🚧 In Progress
- 📋 Planned
- 💡 Under Consideration

---

## Phase 1: Core Analysis (Completed)

✅ **Dependency Graph Parsing**
- Parse package.json files from monorepo workspaces
- Support for npm, yarn, pnpm, and Nx workspaces
- Build in-memory dependency graph

✅ **Circular Dependency Detection**
- Detect direct cycles (A → B → A)
- Detect indirect cycles (A → B → C → A)
- Severity classification (critical, warning, info)
- Impact assessment and fix recommendations

✅ **Health Score Calculation**
- Overall health score (0-100)
- Breakdown by category (dependencies, architecture, maintainability)
- Trend tracking over time

✅ **D3.js Visualization**
- Force-directed graph layout
- Circular dependency highlighting
- Zoom, pan, and minimap navigation
- Node expand/collapse for large graphs
- Hybrid SVG/Canvas rendering for performance

✅ **Report Export**
- HTML standalone reports
- JSON for CI integration
- Markdown for PR descriptions

---

## Phase 2: Enhanced Analysis (Current)

🚧 **WebAssembly Analyzer**
- Client-side analysis using Go compiled to WASM
- No server required for basic analysis
- Privacy-first: files never leave the browser

📋 **Architecture Validation**
- Define layer rules (domain, application, infrastructure)
- Detect layer violations
- Custom rule configuration

📋 **Bundle Impact Analysis**
- Identify duplicate dependencies
- Calculate wasted bundle size
- Suggest consolidation strategies

📋 **Version Conflict Detection**
- Find conflicting dependency versions
- Risk assessment for conflicts
- Resolution recommendations

---

## Phase 3: Integration (Planned)

📋 **GitHub Integration**
- Analyze repositories directly from GitHub URL
- PR comments with analysis results
- Status checks for CI/CD

📋 **CI/CD Integration**
- GitHub Actions workflow
- GitLab CI template
- Configurable thresholds and gates

📋 **CLI Tool**
- Local analysis from command line
- JSON output for scripting
- Watch mode for development

---

## Phase 4: Advanced Features (Future)

💡 **VS Code Extension**
- Real-time circular dependency warnings
- Inline visualization
- Quick fixes

💡 **Historical Tracking**
- Track health score over time
- Regression alerts
- Trend reports

💡 **Team Collaboration**
- Shared workspaces
- Comments and annotations
- Assignment of issues

💡 **Custom Rules Engine**
- Define custom validation rules
- Rule marketplace
- Import/export configurations

---

## Contributing

We welcome contributions! If you'd like to work on any of these features, please:

1. Check if there's an existing issue
2. Open a new issue to discuss your approach
3. Submit a PR referencing the issue

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

---

## Feedback

Have ideas for new features? Open a [GitHub Discussion](https://github.com/user/monoguard/discussions) or [Issue](https://github.com/user/monoguard/issues).
