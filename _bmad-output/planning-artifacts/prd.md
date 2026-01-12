---
stepsCompleted:
  [
    'step-01-init',
    'step-02-discovery',
    'step-03-success',
    'step-04-journeys',
    'step-05-domain',
    'step-06-innovation',
    'step-07-project-type',
    'step-08-scoping',
    'step-09-functional',
    'step-10-nonfunctional',
    'step-11-polish',
    'step-12-complete',
  ]
workflowStatus: 'completed'
completionDate: '2026-01-11'
polishOptimizations:
  lineCountBefore: 1622
  lineCountAfter: 1541
  reductionPercentage: 5.0
  optimizationsApplied:
    - 'Deleted redundant Product Scope section (lines 188-290)'
    - 'Created Executive Summary with competitive positioning front-loaded'
    - 'Added cross-reference between Phase 0 NFR and detailed NFR section'
    - 'Verified information density (no wordy patterns found)'
    - 'Confirmed terminology standardization (Circular Dependency Solution Engine, Dependency Time Machine)'
  expertReviewers: ['Winston (Architect)', 'Mary (Business Analyst)']
inputDocuments:
  - '_bmad-output/analysis/brainstorming-session-2026-01-11.md'
documentCounts:
  briefCount: 0
  researchCount: 0
  brainstormingCount: 1
  projectDocsCount: 0
classification:
  projectType: 'developer_tool'
  domain: 'general'
  complexity: 'medium'
  projectContext: 'brownfield'
  projectContextNote: 'Architecture refactor (MonoGuard 2.0) - TanStack Start + WASM migration'
  migrationStrategy: 'big_bang'
  targetUsers:
    - 'developers'
    - 'technical_management'
  coreProblems:
    - 'dependency_chaos'
    - 'layer_boundary_violations'
    - 'circular_dependencies'
  validationStrategy: 'lean_phased'
  pricingTier: 'freemium_individual'
workflowType: 'prd'
partyModeInsights:
  conducted: true
  date: '2026-01-11'
  participatingAgents: ['Winston', 'Sally', 'Mary']
  keyDecisions:
    - 'MonoGuard 2.0 with breaking changes (big bang migration)'
    - 'Dual interface: developer-focused + management dashboard'
    - 'Freemium model starting with individual developers ($12-29/month)'
    - 'Lean validation: Phase 0 (50 users) → Phase 1 (500 users) → Phase 2 (2000 users)'
    - 'Zero-cost architecture: TanStack Start + WASM from day one'
---

# Product Requirements Document - mono-guard

**Author:** Alexyu
**Date:** 2026-01-11

## Success Criteria

### User Success

**Developer Tier Success Metrics:**

**Core "Aha!" Moment: Circular Dependency Auto-Fixed**

- From discovering circular dependency to resolution: **< 5 minutes**
- Auto-fix success rate: **> 90%**
- Automated PR generation time: **< 30 seconds**
- Developer acceptance rate for auto-fix suggestions: **> 80%**

**User Experience Success:**

- Dependency graph visualization load time: **< 2 seconds** (100 packages)
- Analysis of 1000+ packages: **< 5 seconds**
- CLI command response time: **< 1 second**
- VSCode extension real-time hint latency: **< 500ms**

**Workflow Integration Success:**

- PR check completion time: **< 1 minute**
- PR bot comment accuracy: **> 95%**
- CI/CD integration zero-config time: **< 5 minutes**

**Emotional Success Metrics:**

- User star rating: **≥ 4.0/5.0**
- Net Promoter Score (NPS): **> 40**
- Weekly return rate: **> 50%** (active users)

---

**Technical Management Tier Success Metrics:**

**Core "Relief" Moment: Team Efficiency Improvement Data**

- Technical debt reduction rate: **> 15% per quarter**
- Circular dependency count decrease: **> 10% per month**
- Architecture health score improvement: **> 5 points per quarter** (0-100 scale)
- Build failure rate decrease: **> 30%** (attributed to dependency issues)

**Visualization & Reporting Success:**

- Dashboard load time: **< 3 seconds**
- Historical trend tracking period: **At least 6 months**
- PDF report generation time: **< 10 seconds**
- Report data accuracy: **> 99%**

**Decision Support Success:**

- Technical debt quantification accuracy: **> 90%**
- Refactoring recommendation adoption rate: **> 60%** (management-approved suggestions)
- ROI traceability: **100%** (all improvements quantifiable)

---

### Business Success

**Phase 0 - Concept Validation (0-3 months):**

- ✅ **50 Beta users** registered
- ✅ **10 active users** (3+ uses per week)
- ✅ Average star rating: **≥ 4.0/5.0**
- ✅ **5+ payment intent** users (clicked "Early Access")
- **Go/No-Go Decision Point:** Achieve ≥ 3 metrics → Enter Phase 1

**Phase 1 - Product-Market Fit (3-6 months):**

- ✅ **500 free users**
- ✅ **25 paying users** (5% conversion rate)
- ✅ **$300-500 MRR**
- ✅ **NPS > 40**
- ✅ **Weekly retention rate > 40%**
- **Go/No-Go Decision Point:** MRR > $500 and monthly growth > 20% → Enter Phase 2

**Phase 2 - Scale Preparation (6-12 months):**

- ✅ **2,000 free users**
- ✅ **100 paying users**
- ✅ **$1,500-2,000 MRR**
- ✅ **Viral coefficient > 1.2** (each user brings 1.2 new users)
- ✅ **Monthly growth rate > 15%**

**12-Month Target (Long-term Vision):**

- 💰 **Monthly revenue: $10K-20K MRR**
- 👥 **Paying users: 200-500**
- 📈 **Free users: 10,000+**
- 🎯 **Payment conversion rate: 5-10%**
- 📊 **Customer Lifetime Value (LTV): > $500**
- 💳 **Customer Acquisition Cost (CAC): < $50**

---

### Technical Success

**Architecture & Performance:**

- ✅ **Monthly operating cost: $0** (down from $200)
- ✅ **Maintenance time reduction: 70%**
- ✅ **Performance improvement: 10x** (compared to old version)
- ✅ **Bundle size: < 500KB** (gzipped)
- ✅ **Lighthouse Performance: > 90**
- ✅ **First Contentful Paint (FCP): < 1.5 seconds**

**WASM Analysis Engine:**

- ✅ **Analyze 100 packages: < 5 seconds**
- ✅ **Analyze 1000 packages: < 30 seconds**
- ✅ **Memory usage: < 100MB** (in-browser)
- ✅ **Offline availability: 100%** (no network required)

**Reliability & Quality:**

- ✅ **Test coverage: > 80%**
- ✅ **Critical path coverage: 100%**
- ✅ **P95 error rate: < 0.1%**
- ✅ **Auto-fix accuracy: > 90%**
- ✅ **Zero data loss** (time machine snapshots)

**Developer Experience:**

- ✅ **Zero to local running: < 5 minutes**
- ✅ **Hot reload time: < 2 seconds**
- ✅ **Build time: < 1 minute** (full build)
- ✅ **Deploy to Cloudflare Pages: < 3 minutes**

---

### Measurable Outcomes

**User Behavior Metrics:**

1. **Problem resolution speed:** Discovery to fix < 5 minutes (developers)
2. **Usage frequency:** Active users 3+ times per week
3. **Feature adoption rate:** 80% of users use dependency time machine (Phase 1+)
4. **Auto-fix acceptance rate:** > 80% of auto-fix PRs merged

**Efficiency Improvement Metrics:**

1. **Technical debt reduction:** > 15% per quarter
2. **Build failure reduction:** > 30% (dependency-related issues)
3. **Code review time reduction:** > 25% (due to PR bot pre-checks)
4. **Refactoring time saved:** > 40% (due to auto-fix)

**Product Health Metrics:**

1. **Net Promoter Score (NPS):** > 40 (Phase 1), > 60 (Phase 2)
2. **Customer Satisfaction (CSAT):** > 4.0/5.0
3. **Weekly retention rate:** > 40% (Phase 1), > 60% (Phase 2)
4. **Churn rate:** < 5% (monthly paying user churn)

**Growth Metrics:**

1. **Viral coefficient:** > 1.2 (each user brings 1.2 new users)
2. **Organic traffic share:** > 50% (from TanStack ecosystem, open source exposure)
3. **Payment conversion time:** < 30 days (from registration to first payment)
4. **Upgrade rate:** > 20% (Free → Pro), > 15% (Pro → Pro+)

---

## Executive Summary

**Strategic Positioning:**

> **"Nx tells you there are circular dependencies. MonoGuard tells you how to fix them."**

MonoGuard v2.0 is a monorepo dependency analysis tool that goes beyond visualization to provide actionable solutions. While competitors like Nx and turborepo excel at detecting and visualizing dependency problems, MonoGuard's core differentiation is the **Circular Dependency Solution Engine** - delivering executable fix strategies, not just problem reports.

**Three-Phase Validation Strategy:**

- **Phase 0 (MVP, 0-3 months):** Validate core problem-solving experience with circular dependency fix recommendations. Target: 50 users, 60% fix acceptance rate. Privacy-first architecture (WASM + local-first analysis) with zero infrastructure cost.

- **Phase 1 (Growth, 3-6 months):** Add Dependency Time Machine (historical tracking, trend prediction) and GitHub PR integration. Enable paid conversion ($12-29/month tiers). Target: 500 users, 5% conversion rate, $300-500 MRR.

- **Phase 2 (Scale, 6-12+ months):** Enterprise features (Team Dashboard, AI-powered diagnostics, SSO) and ecosystem integrations. Target: 2,000 users, 100 paying, $1,500-2,000 MRR.

**Core Innovation:** Dependency Time Machine brings the "time dimension" to dependency analysis - timeline visualization, AI-assisted trend forecasting, and causal tracing to specific commits.

→ **See "Project Scoping & Phased Development" for detailed feature breakdown and scope decisions.**

---

## User Journeys

### Journey 1: Sarah (Frontend Developer) - "The Circular Dependency Nightmare"

**Context:**
Sarah is a frontend developer at a fast-growing startup with a large monorepo (50+ packages). It's Friday afternoon, and CI just failed on her PR with a cryptic error: "Circular dependency detected."

**Before MonoGuard:**

- **Discovery time:** 45 minutes of manual tracing through imports
- **Fix time:** 2-3 hours of refactoring, multiple failed attempts
- **Total time:** 3+ hours of frustration
- **Emotional state:** Stressed, weekend plans ruined

**With MonoGuard:**

1. **Discovery (2 seconds):** Opens MonoGuard web interface, drags `package.json` into the browser
2. **Visualization (instant):** Circular dependency path highlighted in red: `@ui/forms` → `@utils/validation` → `@ui/components` → `@ui/forms`
3. **Auto-fix suggestion (30 seconds):** MonoGuard shows "Auto-fix available" button
4. **Review (5 minutes):** Sarah reviews the suggested refactoring: extract shared validation logic into `@utils/validation-schemas`
5. **Apply fix (30 seconds):** Clicks "Generate PR" button, MonoGuard creates a PR with the refactored code
6. **Verification (5 minutes):** Sarah reviews the PR, makes minor adjustments
7. **Merge (2 minutes):** CI passes, PR merged

**Total time:** 45 minutes
**Emotional arc:** Frustration → Curiosity → Surprise ("It actually works!") → Trust → Relief → Excitement ("I need to show this to my team!")

**Aha Moment:**
When the auto-fix PR is generated in 30 seconds with clean, working code, Sarah realizes this tool is not just a visualizer but a productivity multiplier.

---

### Journey 2: David (Tech Lead) - "The Architecture Anxiety"

**Context:**
David is a Tech Lead at a mid-sized company. The monorepo has grown to 150+ packages over 3 years. He knows the architecture is decaying but has no data to convince the CTO to allocate refactoring resources.

**Before MonoGuard:**

- **Problem:** Vague complaints ("codebase is messy") without quantifiable metrics
- **CTO response:** "Show me the ROI before we allocate sprint capacity"
- **Team morale:** Developers frustrated, adding technical debt to every sprint
- **Outcome:** Refactoring requests rejected quarterly

**With MonoGuard:**

1. **Setup (5 minutes):** David connects MonoGuard to the company's GitHub repo
2. **Historical analysis (automatic):** MonoGuard analyzes 6 months of git history, tracking architecture health score over time
3. **Dashboard insights (2 minutes):**
   - **Architecture Health Score:** 58/100 (declining from 72 six months ago)
   - **Circular dependencies:** 23 (up from 8)
   - **Layer boundary violations:** 47 instances
   - **Technical debt heatmap:** Top 10 packages contributing 80% of violations
4. **Report generation (10 seconds):** David clicks "Export PDF Report"
5. **CTO presentation (15 minutes):** Armed with data-driven insights:
   - "Our architecture health dropped 14 points in 6 months"
   - "Build failure rate increased 30% due to dependency issues"
   - "Top 3 refactoring targets would eliminate 60% of violations"
   - "Estimated ROI: 40% reduction in debugging time = 8 developer hours/week saved"

**CTO Decision:** Approves 2-sprint architecture refactoring initiative

**Emotional arc:** Anxiety ("I know it's bad but can't prove it") → Curiosity → Shock ("It's worse than I thought") → Clarity → Confidence → Control ("I can fix this")

**Relief Moment:**
When the CTO approves the refactoring budget after seeing the PDF report with quantified technical debt, David finally has the resources to fix what's been bothering him for months.

---

### Journey 3: Alex (DevOps Engineer) - "The CI/CD Gatekeeper"

**Context:**
Alex manages CI/CD pipelines for a large engineering team (30+ developers). Every 2-3 weeks, a bad dependency change slips through code review and causes production incidents on Friday afternoons.

**Before MonoGuard:**

- **Problem:** Manual code review can't catch all dependency issues
- **Impact:** 1 production incident per week (average)
- **On-call stress:** Alex gets paged every Friday afternoon
- **Team trust:** Developers scared to merge PRs on Fridays

**With MonoGuard:**

1. **Setup (5 minutes):** Alex installs MonoGuard GitHub App from marketplace (zero config)
2. **First PR check (automatic):**
   - Developer opens PR that introduces a circular dependency
   - MonoGuard bot comments within 1 minute: "⚠️ Circular dependency detected: Package A → B → C → A. [View visualization]"
   - PR blocked until issue resolved
3. **Auto-fix workflow:**
   - Developer runs `/monoguard fix` command in PR comment
   - MonoGuard generates fix suggestion commit
   - CI passes, PR approved
4. **Weekly stats (after 2 weeks):**
   - **Blocked PRs:** 5 (caught before merge)
   - **Production incidents:** 0 (vs 1 per week before)
   - **Developer feedback:** "Love the instant feedback loop!"

**Emotional arc:** Stress → Frustration ("Not another Friday incident") → Hope → Validation ("It actually works!") → Trust → Peace ("I can enjoy my weekends now")

**Peace Moment:**
After 2 weeks of zero production incidents, Alex realizes Friday afternoons are no longer stressful. The automated gatekeeper is working.

---

### Journey 4: Sarah (Edge Case) - "When Auto-Fix Can't Help"

**Context:**
Sarah encounters a complex circular dependency that involves shared state management across 5 packages. Auto-fix cannot safely refactor this automatically.

**With MonoGuard:**

1. **Auto-fix attempt:** MonoGuard analyzes the circular dependency
2. **Honest limitation:** "Auto-fix not available for this complexity. Reason: Shared state management across 5 packages requires architectural decision."
3. **Manual guidance provided:**
   - **Visualization:** Interactive dependency graph showing all 5 packages and their relationships
   - **Root cause explanation:** "The cycle exists because `@state/auth`, `@state/user`, and `@ui/profile` all share mutable state."
   - **Refactoring strategies:**
     - Option A: Extract shared state into `@state/core` (recommended)
     - Option B: Use event-driven architecture with pub/sub pattern
     - Option C: Implement dependency injection
   - **Step-by-step migration guide:** Detailed checklist with code examples

**Outcome:**

- **Manual fix time:** 1 hour (vs 3 hours without guidance)
- **Emotional response:** "Even when it can't auto-fix, it teaches me how to fix it better than any documentation I've read"

**Trust Building:**
MonoGuard earns Sarah's trust by being honest about limitations while still providing maximum value through education and guidance.

---

## Journey Requirements Summary

**From these journeys, we extract these critical requirements:**

1. **Analysis Speed:**
   - Analyze 100+ packages in < 5 seconds
   - Real-time visualization updates
   - Sub-second response to user interactions

2. **Visualization Quality:**
   - Interactive dependency graph (D3.js)
   - Highlight circular dependencies in red
   - Drill-down capability for complex graphs
   - Export high-quality images for reports

3. **Auto-Fix Intelligence:**
   - Detect circular dependencies automatically
   - Generate safe refactoring suggestions
   - Create PR with working code
   - Success rate > 90% for common patterns

4. **GitHub Integration:**
   - Zero-config GitHub App installation
   - PR checks complete in < 1 minute
   - Bot comments with actionable insights
   - `/monoguard fix` command for quick fixes

5. **Management Dashboard:**
   - Architecture health score (0-100)
   - Historical trend tracking (6+ months)
   - Technical debt heatmap
   - PDF report generation (< 10 seconds)
   - ROI calculation and prioritization

6. **Developer Experience:**
   - Drag-and-drop interface (no CLI required for first use)
   - Clear error messages and explanations
   - Educational guidance when auto-fix unavailable
   - VSCode extension for real-time hints

7. **Reliability:**
   - Offline-first (WASM runs in browser)
   - Zero data loss (time machine snapshots)
   - Honest about limitations
   - Graceful degradation

**Success Metric Traceability:**

| Journey  | Success Metric                 | Target         |
| -------- | ------------------------------ | -------------- |
| Sarah #1 | Problem resolution time        | < 5 minutes    |
| Sarah #1 | Auto-fix success rate          | > 90%          |
| David    | Architecture health visibility | 6-month trend  |
| David    | Technical debt quantification  | > 90% accuracy |
| Alex     | PR check completion time       | < 1 minute     |
| Alex     | Production incident reduction  | > 30%          |
| All      | User star rating               | ≥ 4.0/5.0      |
| All      | Weekly return rate             | > 50%          |

---

## Innovation & Novel Patterns

### Detected Innovation Areas

**🎯 Core Innovation: Dependency Time Machine**

MonoGuard's greatest differentiation is introducing the "time dimension" into dependency analysis:

1. **Time Machine Core Innovation:**
   - **Timeline Visualization:** 6-month historical snapshots + animated time travel
   - **Trend Prediction:** AI-assisted forecasting of architecture health trajectory
   - **Causal Tracing:** "Why did this circular dependency appear?" → pinpoint to specific commit
   - **Comparative Analysis:** Compare any two time points ("3 months ago vs now")

2. **Zero-Cost Architecture:**
   - **Technology Combination Innovation:** TanStack Start + WASM = fully offline SaaS
   - **Economic Model Breakthrough:** $200/month → $0/month infrastructure cost
   - **Performance Advantage:** Analyze 1000+ packages in browser < 30 seconds

3. **AI-Powered Auto-Fix:**
   - **Intelligent Refactoring:** Auto-generate PRs (> 90% success rate)
   - **Educational Mode:** When auto-fix unavailable, provide detailed guidance
   - **Learning Mechanism:** Learn from user acceptance/rejection of PRs

---

### Market Context & Competitive Landscape

**Limitations of Existing Tools:**

| Tool Category       | Representative Tools      | Limitations                                                                            |
| ------------------- | ------------------------- | -------------------------------------------------------------------------------------- |
| Static Analyzers    | Madge, dependency-cruiser | ❌ Only analyze current state<br>❌ No historical tracking<br>❌ Manual fixes required |
| Visualization Tools | nx graph, turborepo       | ❌ Lack time dimension<br>❌ No auto-fix capability                                    |
| Linter Integration  | ESLint plugins            | ❌ Detection only<br>❌ No visualization                                               |

**MonoGuard's Unique Positioning:**

```
Static Analysis + Time Travel + AI Fix + Zero-Cost Deployment = New Category
```

**"Nothing like this exists" Validation:**

- ✅ No tool provides historical timeline visualization for dependencies
- ✅ No tool combines WASM for fully offline + SaaS experience
- ✅ No tool achieves > 90% auto-fix success rate

---

### Validation Approach

**Phase 0 - Concept Validation (0-3 months):**

| Innovation Feature         | Validation Metric           | Target   | Measurement Method                          |
| -------------------------- | --------------------------- | -------- | ------------------------------------------- |
| **Time Machine**           | Feature usage rate          | > 90%    | PostHog tracking: % users clicking timeline |
| **Auto-Fix**               | Auto-fix success rate       | > 90%    | GitHub API: PR merge rate                   |
| **WASM Performance**       | 1000 packages analysis time | < 30 sec | In-browser timer                            |
| **Zero-Cost Architecture** | Infrastructure cost         | $0       | Cloudflare Pages billing                    |

**Phase 1 - Product-Market Fit (3-6 months):**

| Innovation Feature     | Validation Metric      | Target        | Measurement Method              |
| ---------------------- | ---------------------- | ------------- | ------------------------------- |
| **Time Machine**       | Paid conversion driver | Top 3 ranking | User survey: "Why upgrade Pro?" |
| **Auto-Fix**           | PR acceptance rate     | > 80%         | GitHub webhook tracking         |
| **Offline Capability** | Offline usage rate     | > 30%         | IndexedDB usage tracking        |

**Key Questions for Innovation Validation:**

- 📊 Is Time Machine a "cool demo feature" or "daily-use tool"?
- 🤖 Does Auto-Fix actually save time or increase review burden?
- ⚡ Is WASM performance sufficient for enterprise-scale monorepos (2000+ packages)?

---

### Risk Mitigation

#### Risk 1: WASM Performance Limitations

**Potential Issues:**

- Very large monorepos (2000+ packages) may exceed browser memory limits
- Complex dependency graph visualization may lag

**Mitigation Strategies:**

- ✅ **Chunked Analysis:** Process large monorepos in batches (500 packages per batch)
- ✅ **Progressive Rendering:** Use D3.js with canvas instead of SVG (handle 10,000+ nodes)
- ✅ **Graceful Degradation:** For > 2000 packages, suggest CLI version (Go native)
- ✅ **Performance Monitoring:** Sentry Performance Monitoring tracks P95 analysis time

**Validation Points:**

- MVP test: Can analyze 1000 packages on 4GB RAM MacBook?
- Phase 1 test: User-reported "too slow" issues < 5%

---

#### Risk 2: Auto-Fix Accuracy Insufficient

**Potential Issues:**

- Auto-generated refactored code may break business logic
- Edge cases (e.g., circular dependency involving shared state) cannot be handled
- Users don't trust auto-generated code

**Mitigation Strategies:**

- ✅ **Conservative Approach:** Only handle "safe" refactoring patterns (extract module, move import)
- ✅ **Progressive Trust:**
  - Phase 0: Suggest only, no auto-PR generation
  - Phase 1: Generate PR but require manual review
  - Phase 2: One-click merge (based on high trust accumulation)
- ✅ **Honest Degradation:** When unsafe to auto-fix, explicitly explain and provide manual guide (Journey 4)
- ✅ **Test Coverage:** Auto-generated PRs include tests (ensure refactoring doesn't break functionality)
- ✅ **Learning Mechanism:** Learn from PR merge/rejection data to improve

**Validation Points:**

- MVP: Auto-fix suggestion acceptance rate > 80%
- Phase 1: Auto-generated PR merge rate > 80%
- User trust: NPS > 40 (proves user trust in tool)

---

#### Risk 3: Time Machine "Gimmickification"

**Potential Issues:**

- Time Machine may be just a "flashy demo feature" with low actual usage
- Users may find it "useless" or "too complex"

**Mitigation Strategies:**

- ✅ **Concrete Use Case Design:**
  - "View new circular dependencies introduced last week" (daily use)
  - "Compare architecture health score from 3 months ago" (quarterly retrospective)
  - "Trace root commit of technical debt" (refactoring decisions)
- ✅ **Simplified UX:**
  - Timeline slider (similar to Figma version history)
  - One-click "play" button (auto-animate time progression)
- ✅ **Value Proof:**
  - Rank Time Machine value in user surveys
  - A/B test: Conversion rate with vs without Time Machine

**Validation Points:**

- Phase 0: > 90% of active users use Time Machine at least once
- Phase 1: Time Machine ranks Top 3 in user survey

---

## Developer Tool Specific Requirements

### Platform & Language Matrix

**Execution Environments:**

| Platform             | Technology            | Distribution                             | Status          |
| -------------------- | --------------------- | ---------------------------------------- | --------------- |
| **Web Interface**    | TanStack Start + WASM | Cloudflare Pages (https://monoguard.dev) | ✅ Phase 0      |
| **CLI Tool**         | Go (native binary)    | npm (`npm install -g monoguard`)         | ✅ Phase 0      |
| **Local Dev Server** | CLI-embedded server   | `monoguard analyze --serve`              | ✅ Phase 0      |
| **GitHub App**       | Node.js backend       | GitHub Marketplace                       | ✅ Phase 1      |
| **VSCode Extension** | N/A                   | Skipped                                  | ❌ Out of scope |

**Key Architectural Decision:**

- ✅ **Privacy-First Design:** CLI analyzes locally + spawns local dev server to display results
- ✅ **Zero Data Upload:** Users can avoid uploading private repo code to web or database
- ✅ **Hybrid Experience:** Combine CLI power with web UI convenience

---

### Installation Methods

**1. Web Interface (Zero Installation):**

```
Visit: https://monoguard.dev
Drag & drop package.json or upload workspace files
All analysis runs in browser via WASM
```

**2. CLI Tool (npm Global Install):**

```bash
# Install globally via npm
npm install -g monoguard

# Verify installation
monoguard --version
```

**Alternative Installation Methods (Future):**

```bash
# Homebrew (macOS/Linux)
brew install monoguard

# Direct binary download (cross-platform)
curl -fsSL https://monoguard.dev/install.sh | sh
```

**3. Local Dev Server (Privacy Mode):**

```bash
# Analyze locally and serve results on localhost
monoguard analyze --serve
# Opens browser at http://localhost:3000 with results

# Analyze and export static HTML (offline viewing)
monoguard analyze --export ./report
```

---

### CLI API Surface

**Core Commands:**

**Command: analyze**

```bash
monoguard analyze [path]

Options:
  --serve           Start local dev server to view results
  --export <path>   Export static HTML report
  --format <type>   Output format: json|html|markdown
  --depth <n>       Analysis depth (default: unlimited)
  --exclude <pattern> Exclude packages matching pattern

Examples:
  monoguard analyze                    # Analyze current directory
  monoguard analyze ./packages         # Analyze specific directory
  monoguard analyze --serve            # Analyze + launch local UI
  monoguard analyze --export ./report  # Export static report
```

**Command: check**

```bash
monoguard check [path]

Purpose: CI/CD validation (exit code 0 = pass, 1 = fail)

Options:
  --fail-on <type>  Fail on: circular|boundary|all
  --threshold <n>   Fail if health score < n (0-100)
  --config <path>   Custom config file

Examples:
  monoguard check                           # Check current directory
  monoguard check --fail-on circular        # Fail only on circular deps
  monoguard check --threshold 70            # Fail if health score < 70
  monoguard check --config .monoguard.json  # Use custom config
```

**Command: fix**

```bash
monoguard fix [path]

Purpose: Auto-fix suggestions and PR generation

Options:
  --dry-run         Show suggestions without applying
  --auto-commit     Automatically create git commit
  --pr              Create GitHub PR (requires auth)
  --interactive     Interactive fix selection

Examples:
  monoguard fix --dry-run        # Preview suggestions
  monoguard fix --auto-commit    # Apply fixes + commit
  monoguard fix --pr             # Apply fixes + create PR
  monoguard fix --interactive    # Choose which fixes to apply
```

**Utility Commands:**

```bash
monoguard init       # Initialize .monoguard.json config
monoguard auth       # Authenticate with GitHub (for PR features)
monoguard --version  # Show version
monoguard --help     # Show help
```

---

### WASM API Surface

**JavaScript/TypeScript API for browser integration:**

```typescript
// Import WASM module
import { analyze, check } from '@monoguard/wasm';

// 1. analyze: Full dependency analysis
interface AnalyzeOptions {
  workspaceRoot: string;
  packages: string[];
  depth?: number;
  exclude?: string[];
}

interface AnalyzeResult {
  dependencies: DependencyGraph;
  circularDeps: CircularDependency[];
  boundaryViolations: BoundaryViolation[];
  healthScore: number;
  metadata: AnalysisMetadata;
}

const result: AnalyzeResult = await analyze(options);

// 2. check: Validation only (fast)
interface CheckOptions {
  workspaceRoot: string;
  packages: string[];
  failOn?: 'circular' | 'boundary' | 'all';
  threshold?: number;
}

interface CheckResult {
  passed: boolean;
  errors: ValidationError[];
  healthScore: number;
}

const checkResult: CheckResult = await check(options);
```

**Example Usage in Browser:**

```typescript
// Drag & drop file upload
async function handleFileUpload(files: FileList) {
  const result = await analyze({
    workspaceRoot: '/',
    packages: Array.from(files).map((f) => f.name),
  });

  // Display results in UI
  renderDependencyGraph(result.dependencies);
  renderHealthScore(result.healthScore);
}
```

---

### Configuration Schema

**`.monoguard.json` Format:**

```json
{
  "$schema": "https://monoguard.dev/schema.json",
  "version": "2.0",
  "workspaces": ["packages/*", "apps/*"],
  "rules": {
    "circularDependencies": "error",
    "boundaryViolations": "warn",
    "duplicateDependencies": "warn"
  },
  "layers": {
    "ui": ["packages/ui-*"],
    "logic": ["packages/logic-*"],
    "data": ["packages/data-*"],
    "allowed": ["ui -> logic", "logic -> data"]
  },
  "exclude": ["**/node_modules/**", "**/dist/**", "**/*.test.ts"],
  "thresholds": {
    "healthScore": 70,
    "maxCircularDeps": 0,
    "maxBoundaryViolations": 5
  },
  "timeMachine": {
    "enabled": true,
    "retentionDays": 180,
    "snapshotFrequency": "daily"
  },
  "autoFix": {
    "enabled": true,
    "safeMode": true,
    "createPR": false
  }
}
```

---

### Code Examples & Quick Start

**Quick Start (5-Minute Setup):**

```markdown
# Quick Start

## 1. Install MonoGuard

npm install -g monoguard

## 2. Analyze Your Monorepo

cd your-monorepo
monoguard analyze --serve

## 3. View Results

Browser opens at http://localhost:3000 with interactive dependency graph!

That's it! 🎉
```

---

### Common Integration Use Cases

**Use Case 1: CI/CD Integration (GitHub Actions)**

```yaml
# .github/workflows/monoguard.yml
name: MonoGuard Check

on: [push, pull_request]

jobs:
  dependency-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Install MonoGuard
        run: npm install -g monoguard
      - name: Check Dependencies
        run: monoguard check --fail-on circular --threshold 70
```

**Use Case 2: Pre-commit Hook**

```bash
# .husky/pre-commit
#!/bin/sh
monoguard check --fail-on circular || exit 1
```

**Use Case 3: Local Privacy-First Analysis**

```bash
# Analyze private repo without uploading to cloud
monoguard analyze --serve

# Or export static HTML for offline viewing
monoguard analyze --export ./dependency-report
open ./dependency-report/index.html
```

---

### Migration Guide

**From Madge:**

```bash
# Before (Madge)
madge --circular --extensions ts,tsx src/

# After (MonoGuard)
monoguard check --fail-on circular
```

**From dependency-cruiser:**

```bash
# Before (dependency-cruiser)
dependency-cruiser --validate .dependency-cruiser.json src

# After (MonoGuard)
monoguard check --config .monoguard.json
```

**Key Migration Benefits:**

- ✅ **Time Machine:** Historical tracking (Madge doesn't have)
- ✅ **Auto-Fix:** Automated PR generation (dependency-cruiser doesn't have)
- ✅ **Privacy:** Local-first analysis option
- ✅ **UI:** Beautiful web interface vs terminal output
- ✅ **Performance:** WASM-powered analysis (10x faster)

---

### Documentation Structure

**Required Documentation:**

1. **Getting Started Guide**
   - Installation (5 minutes)
   - First analysis (5 minutes)
   - Understanding results (10 minutes)

2. **CLI Reference**
   - All commands with examples
   - Configuration options
   - Exit codes and error handling

3. **WASM API Reference**
   - TypeScript type definitions
   - Integration examples
   - Browser compatibility

4. **Configuration Guide**
   - `.monoguard.json` schema
   - Layer boundary configuration
   - Rule customization
   - Time Machine settings

5. **Integration Guides**
   - GitHub Actions
   - GitLab CI
   - Pre-commit hooks
   - Custom CI/CD systems

6. **Migration Guides**
   - From Madge
   - From dependency-cruiser
   - From nx graph
   - From custom scripts

---

### Privacy & Security Architecture

**Privacy-First Design Principles:**

1. **Local-First Analysis:**
   - All analysis runs locally (CLI or browser WASM)
   - No code uploaded to remote servers
   - No dependency on external APIs for core features

2. **Optional Cloud Features:**
   - GitHub App: Requires explicit authentication
   - Time Machine: Historical data stored in user's GitHub repo (`.monoguard/` directory)
   - Reports: Generated locally, user decides where to store

3. **Data Storage:**
   - **Web Version:** IndexedDB (browser-local storage)
   - **CLI Version:** `.monoguard/` directory in project root
   - **GitHub App:** Uses GitHub API, no separate database

4. **No Telemetry Without Consent:**
   - Anonymous usage analytics opt-in only
   - Error reporting opt-in only
   - Clear privacy policy

---

## Project Scoping & Phased Development

### MVP Strategy & Philosophy

**MVP Approach: Experience MVP (體驗型 MVP)**

MonoGuard 採用"體驗 MVP"策略，展示完整的問題解決體驗，而非僅提供基礎功能。核心理念是讓早期用戶體驗到產品的獨特價值主張：**不只檢測問題，更提供解決方案**。

**Strategic Positioning:**

> **"Nx 告訴你有循環依賴，MonoGuard 告訴你如何修復它"**

**Competitive Reality:**

- ✅ **Acknowledged:** Nx 和 turborepo 已有成熟的依賴圖視覺化和循環依賴檢測
- ✅ **Differentiation:** MonoGuard 的核心差異在於提供**可執行的修復方案**，而非僅發現問題

**Resource Requirements:**

- **Team Size:** 1-2 開發者（Solo founder friendly）
- **Timeline:** 2-3 個月（Phase 0）
- **Key Skills Required:**
  - Go 開發（WASM 編譯）
  - TanStack Start（Web 介面）
  - 靜態分析引擎設計（依賴圖解析）
  - 規則引擎開發（修復建議邏輯）

---

### MVP Feature Set (Phase 0: 0-3 Months)

**Core User Journeys Supported:**

- ✅ **Primary:** Sarah (Frontend Developer) - "The Circular Dependency Nightmare"
  - 從發現循環依賴到取得修復建議 < 5 分鐘
  - 修復建議接受率 > 60%
- 🔶 **Partial:** David (Tech Lead) - "The Architecture Anxiety"
  - 提供架構健康分數，但無歷史追蹤（Phase 1 功能）
- ✅ **Full:** Sarah (Edge Case) - "When Auto-Fix Can't Help"
  - 誠實降級 + 教育指南

**Must-Have Capabilities:**

**1. Dependency Analysis Engine (Table Stakes)**

- ✅ Monorepo workspace detection (npm/yarn/pnpm)
- ✅ Complete dependency graph construction
- ✅ Circular dependency identification (all paths)
- ✅ Duplicate dependency detection
- ✅ Architecture health score calculation (0-100)
- **Performance Target:** Analyze 100 packages < 5 seconds

**2. Visualization Interface (Table Stakes)**

- ✅ Interactive dependency graph (D3.js)
- ✅ Circular dependency path highlighting (red visualization)
- ✅ Click to expand/collapse nodes
- ✅ Export graph as PNG/SVG
- ✅ Dual interface: Web (WASM) + CLI
- **Note:** Feature parity with Nx, not differentiation

**3. 🎯 Circular Dependency Solution Engine (Core Differentiator)**

This is MonoGuard's killer feature - **must be completed in Phase 0:**

**Level 1: Detailed Diagnostic Report (Minimum Viable)**

- ✅ **Root Cause Analysis:** "Why does this circular dependency exist?"
  - Identify the problematic import statements
  - Trace the circular path step-by-step
- ✅ **Impact Assessment:** "How many packages are affected?"
  - Show dependency tree depth
  - Calculate refactoring complexity score
- ✅ **Fix Strategy Recommendations:**
  - **Strategy A:** Extract Shared Module (create new package for shared logic)
  - **Strategy B:** Dependency Injection (break cycle with interfaces)
  - **Strategy C:** Module Boundary Refactoring (restructure package boundaries)
- ✅ **Step-by-Step Fix Guide:**
  - Detailed walkthrough similar to Journey 4
  - Code location references (file:line)
  - Before/after explanations

**Level 2: Automated Refactoring Suggestions (Stretch Goal)**

- 🔶 **Smart Analysis:** Auto-detect best fix strategy based on code patterns
- 🔶 **Code Examples:** Show "before vs after" code comparison
- 🔶 **Copy-Paste Ready:** Provide copyable code snippets
- **Note:** If time permits in Phase 0, otherwise Phase 1

**4. CLI Tool**

```bash
monoguard analyze [path]       # Full analysis with fix suggestions
monoguard check [path]         # CI/CD validation (exit code 0/1)
monoguard fix --dry-run [path] # Preview fix suggestions (CORE!)
monoguard init                 # Initialize .monoguard.json
```

**5. Web Interface**

- ✅ Drag & drop package.json upload
- ✅ In-browser analysis (WASM execution)
- ✅ Interactive dependency graph
- ✅ Fix suggestions panel
- ✅ Export analysis report (HTML/JSON)

**6. Privacy-First Architecture**

- ✅ Local-first analysis (no code upload)
- ✅ WASM in-browser execution
- ✅ CLI local analysis with optional `--serve` for local dev server
- ✅ IndexedDB for browser storage
- ✅ `.monoguard/` directory for CLI storage

**Technical Foundation:**

- ✅ Go → WASM (analysis engine)
- ✅ TanStack Start (SSG for web interface)
- ✅ D3.js (visualization)
- ✅ Cloudflare Pages (free hosting)

**Non-Functional Requirements (Phase 0):**

- ✅ Analyze 100 packages < 5 seconds
- ✅ Bundle size < 500KB (gzipped)
- ✅ Lighthouse Performance > 90
- ✅ Fix suggestion accuracy > 60% (user acceptance rate)

→ **See "Non-Functional Requirements" section for complete NFR specifications with measurement criteria.**

---

### Explicitly Excluded from Phase 0

**Deferred to Phase 1:**

1. **❌ Dependency Time Machine**
   - Historical snapshots tracking
   - Timeline visualization
   - Trend analysis
   - **Reason:** Core differentiator is fixing, not tracking history

2. **❌ GitHub PR Integration**
   - PR checks
   - Bot comments
   - `/monoguard fix` command in PRs
   - **Reason:** Focus on local developer experience first

3. **❌ Auto-Fix PR Generation**
   - Automatic PR creation with fixes
   - Git commit automation
   - **Reason:** Level 1 fix suggestions sufficient for MVP validation

4. **❌ Advanced CLI Features**
   - `monoguard analyze --serve` (local dev server)
   - Complex `.monoguard.json` configuration
   - **Reason:** Use defaults, add customization later

**Permanently Out of Scope:**

5. **❌ VSCode Extension**
   - Real-time hints in editor
   - **Reason:** User confirmed not needed

6. **❌ Team Dashboard**
   - Management interface
   - Team analytics
   - **Reason:** Phase 2 feature, focus on individual developers first

7. **❌ AI Diagnostics (Claude API)**
   - AI-powered fix recommendations
   - **Reason:** Rule-based engine sufficient for Phase 0

---

### Post-MVP Features

**Phase 1 (3-6 Months): Growth Features**

**Goal:** Differentiate further and enable paid conversion

**Core Features:**

1. **Dependency Time Machine** ⭐
   - 6-month historical tracking
   - Timeline slider visualization
   - Trend prediction (AI-assisted)
   - Comparative analysis ("3 months ago vs now")
   - **Monetization:** Pro feature ($12/month)

2. **GitHub PR Integration**
   - Zero-config GitHub App installation
   - PR checks (< 1 minute completion)
   - Bot comments with fix suggestions
   - `/monoguard fix` command
   - **Monetization:** Pro+ feature ($29/month)

3. **Auto-Fix PR Generation**
   - Automatic PR creation with fixes
   - Success rate target: > 80%
   - **Monetization:** Pro+ feature ($29/month)

4. **Advanced CLI Features**
   - `monoguard analyze --serve` (local dev server)
   - Full `.monoguard.json` configuration
   - Layer boundary validation
   - Custom rule definitions

5. **Enhanced Fix Suggestions (Level 2)**
   - Automated refactoring code generation
   - Before/after code comparison UI
   - Copy-paste ready code snippets
   - **Monetization:** Pro feature

**Phase 1 Success Metrics:**

- 500 free users
- 25 paying users (5% conversion rate)
- $300-500 MRR
- NPS > 40
- Weekly retention rate > 40%

---

**Phase 2 (6-12 Months): Expansion Features**

**Goal:** Scale and enterprise readiness

**Core Features:**

1. **Team Dashboard**
   - Management interface for Tech Leads/CTOs
   - Multi-project tracking
   - Team usage analytics
   - Technical debt quantification
   - **Monetization:** Team tier ($99/month)

2. **AI-Powered Diagnostics**
   - Claude API integration
   - Context-aware fix recommendations
   - Natural language explanations
   - **Monetization:** Pro+ feature

3. **Enterprise Features**
   - SSO/SAML authentication
   - Audit logs
   - Custom SLA
   - White-label reports
   - **Monetization:** Enterprise tier (custom pricing)

4. **Ecosystem Integrations**
   - GitLab CI/CD
   - Slack/Discord notifications
   - Renovate/Dependabot enhancement
   - npm/yarn/pnpm CLI plugins

5. **Fully Offline Version**
   - Downloadable desktop app
   - No network required
   - Perpetual license model
   - **Monetization:** One-time $99 purchase

**Phase 2 Success Metrics:**

- 2,000 free users
- 100 paying users
- $1,500-2,000 MRR
- Viral coefficient > 1.2
- Monthly growth rate > 15%

---

### Risk Mitigation Strategy

**Technical Risks:**

**Risk 1: Fix Suggestion Accuracy Insufficient**

**Potential Impact:** Users don't trust recommendations, abandon tool

**Mitigation:**

- ✅ **Conservative Approach:** Start with rule-based engine (3 proven patterns)
- ✅ **Validation Loop:** Track acceptance rate, iterate on rules
- ✅ **Honest Degradation:** When uncertain, explain limitations clearly (Journey 4)
- ✅ **Phase Progression:**
  - Phase 0: Rule-based (target 60% acceptance)
  - Phase 1: Enhanced rules (target 80% acceptance)
  - Phase 2: AI-assisted (target 90% acceptance)

**Validation Points:**

- MVP: Acceptance rate > 60%
- Phase 1: Acceptance rate > 80%
- User feedback: NPS > 40

---

**Risk 2: WASM Performance Limitations**

**Potential Impact:** Cannot handle enterprise-scale monorepos (2000+ packages)

**Mitigation:**

- ✅ **Chunked Processing:** Analyze in batches (500 packages per batch)
- ✅ **Progressive Rendering:** Canvas-based D3.js for large graphs
- ✅ **Graceful Degradation:** Suggest CLI for > 2000 packages
- ✅ **Performance Monitoring:** Sentry tracking of P95 analysis time

**Validation Points:**

- MVP: 1000 packages < 30 seconds on 4GB RAM laptop
- Phase 1: "Too slow" complaints < 5%

---

**Market Risks:**

**Risk 3: Competition with Nx (Free Tool)**

**Potential Impact:** Users stick with free Nx, no paid conversion

**Mitigation:**

- ✅ **Clear Positioning:** "Nx finds problems, MonoGuard solves them"
- ✅ **Quantified Value:** Sarah #1 saves 2.5 hours (3h → 45min)
- ✅ **Freemium Strategy:** Phase 0 completely free, build trust first
- ✅ **Differentiated Pro Features:** Time Machine, PR Integration (Nx doesn't have)

**Validation Points:**

- Phase 0: 5+ users express payment intent
- Phase 1: 5% conversion rate (free → Pro)

---

**Risk 4: Time Machine "Gimmickification"**

**Potential Impact:** Key differentiator seen as "flashy but useless"

**Mitigation:**

- ✅ **Deferred to Phase 1:** Validate core value (fixing) first
- ✅ **Concrete Use Cases:** Weekly retrospective, quarterly reviews
- ✅ **Simplified UX:** Timeline slider (like Figma version history)
- ✅ **Data-Driven Decision:** A/B test conversion with vs without

**Validation Points:**

- Phase 1: > 90% active users try Time Machine at least once
- Phase 1: Time Machine ranks Top 3 in user surveys

---

**Resource Risks:**

**Risk 5: Solo Founder Bandwidth**

**Potential Impact:** Cannot complete Phase 0 in 3 months

**Mitigation:**

- ✅ **Absolute Core (60-day plan):**
  1. Analysis engine (2 weeks)
  2. Visualization (2 weeks)
  3. Fix suggestions Level 1 (3 weeks)
  4. Basic CLI (1 week)
  5. Basic web UI (2 weeks)
- ✅ **Optional Enhancements:**
  - Fix suggestions Level 2 → Phase 1
  - Local dev server → Phase 1
  - Configuration system → Use defaults

**Contingency Plan:**

- If behind schedule, ship with Level 1 fix suggestions only
- Validate core value before adding Level 2

---

## Functional Requirements

### Dependency Analysis & Detection

**FR1:** Users can analyze monorepo dependency graphs by uploading workspace configuration files

**FR2:** Users can detect circular dependencies across all packages in a monorepo

**FR3:** Users can identify duplicate dependencies with version conflicts

**FR4:** Users can view architecture health score (0-100) calculated from dependency analysis

**FR5:** Users can analyze npm, yarn, and pnpm workspace structures

**FR6:** Users can exclude specific packages or patterns from analysis

---

### Circular Dependency Resolution (Core Differentiator)

**FR7:** Users can view root cause analysis for each detected circular dependency

**FR8:** Users can see which import statements create circular dependency paths

**FR9:** Users can receive fix strategy recommendations for circular dependencies

**FR10:** Users can view step-by-step fix guides with code location references

**FR11:** Users can access three fix strategy options: Extract Shared Module, Dependency Injection, Module Boundary Refactoring

**FR12:** Users can see refactoring complexity scores for each circular dependency

**FR13:** Users can view impact assessment showing how many packages are affected by each circular dependency

**FR14:** Users can receive before/after explanations for recommended fixes

---

### Visualization & Reporting

**FR15:** Users can view interactive dependency graphs with D3.js visualization

**FR16:** Users can see circular dependencies highlighted in red on dependency graphs

**FR17:** Users can expand and collapse nodes in dependency graphs

**FR18:** Users can export dependency graphs as PNG or SVG images

**FR19:** Users can export analysis reports in HTML and JSON formats

**FR20:** Users can view detailed diagnostic reports for circular dependencies

---

### CLI Interface

**FR21:** Users can analyze dependencies via CLI command (`monoguard analyze`)

**FR22:** Users can run CI/CD validation checks via CLI (`monoguard check`)

**FR23:** Users can preview fix suggestions via CLI (`monoguard fix --dry-run`)

**FR24:** Users can initialize configuration files via CLI (`monoguard init`)

**FR25:** Users can configure analysis depth and exclusion patterns via CLI options

**FR26:** Users can receive exit codes indicating analysis results (0 = pass, 1 = fail)

**FR27:** Users can export analysis results in multiple formats (JSON, HTML, Markdown) via CLI

---

### Web Interface

**FR28:** Users can drag and drop package.json files to initiate analysis in web browser

**FR29:** Users can upload multiple workspace files to analyze complete monorepo structure

**FR30:** Users can execute dependency analysis entirely in browser via WASM

**FR31:** Users can view fix suggestions panel alongside dependency graph in web interface

**FR32:** Users can download analysis reports from web interface

**FR33:** Users can access web interface without account creation or authentication

---

### Privacy & Data Management

**FR34:** Users can perform complete analysis without uploading code to remote servers

**FR35:** Users can store analysis results locally in browser IndexedDB

**FR36:** Users can store analysis results in local `.monoguard/` directory when using CLI

**FR37:** Users can execute all core analysis features offline without network connection

**FR38:** Users can opt-in to anonymous usage analytics

**FR39:** Users can opt-in to error reporting

---

### Configuration & Customization

**FR40:** Users can configure circular dependency detection rules

**FR41:** Users can define custom architecture health score thresholds

**FR42:** Users can configure package exclusion patterns

**FR43:** Users can set workspace detection patterns

**FR44:** Users can configure analysis output formats

---

### WASM API (For Integration)

**FR45:** Developers can integrate MonoGuard analysis engine into custom applications via WASM API

**FR46:** Developers can call `analyze()` function with workspace configuration to get full analysis results

**FR47:** Developers can call `check()` function for validation-only operations

**FR48:** Developers can receive typed results (DependencyGraph, CircularDependency, HealthScore) from WASM API

---

## Non-Functional Requirements

### Performance

**NFR1: Analysis Speed**

- Analyze 100 packages in < 5 seconds (P95)
- Analyze 1000 packages in < 30 seconds (P95)
- Measurement: In-browser timer for WASM, CLI execution time

**NFR2: UI Responsiveness**

- Dependency graph visualization renders in < 2 seconds for 100 packages
- User interactions (expand/collapse nodes) respond in < 500ms
- Page load (First Contentful Paint) < 1.5 seconds
- Measurement: Lighthouse Performance score > 90

**NFR3: Bundle Size**

- Web application bundle size < 500KB (gzipped)
- WASM module size < 2MB (uncompressed)
- Measurement: Build output analysis

**NFR4: Memory Efficiency**

- In-browser WASM analysis uses < 100MB RAM for 1000 packages
- CLI tool uses < 200MB RAM for 1000 packages
- Graceful degradation for memory-constrained environments (suggest CLI for > 2000 packages)
- Measurement: Browser DevTools, OS process monitoring

---

### Reliability

**NFR5: Offline Availability**

- All core analysis features (analyze, check, fix suggestions) work 100% offline
- No network dependency for primary user workflows
- Measurement: Functional testing with network disabled

**NFR6: Error Handling**

- Analysis errors do not crash the application
- All errors provide actionable error messages with file:line references
- P95 error rate < 0.1% for valid workspace inputs
- Measurement: Error tracking (Sentry), user feedback

**NFR7: Data Integrity**

- Zero data loss for local storage (IndexedDB, `.monoguard/` directory)
- Analysis results are reproducible (same input = same output)
- Measurement: Automated testing, checksum validation

**NFR8: Fix Suggestion Accuracy**

- Fix suggestion acceptance rate > 60% (Phase 0 target)
- Fix suggestion acceptance rate > 80% (Phase 1 target)
- Measurement: User tracking via PostHog, GitHub PR merge rate

---

### Security & Privacy

**NFR9: Privacy-First Architecture**

- Zero code upload to remote servers for core analysis features
- All analysis runs locally (WASM in browser, CLI on local machine)
- No user authentication required for core features
- Measurement: Network traffic inspection, privacy audit

**NFR10: Data Storage**

- Browser data stored exclusively in IndexedDB (local-only)
- CLI data stored exclusively in `.monoguard/` directory (local-only)
- No external database or cloud storage for analysis results
- Measurement: Code review, data flow analysis

**NFR11: Telemetry Consent**

- Anonymous usage analytics opt-in only (not enabled by default)
- Error reporting opt-in only (not enabled by default)
- Clear privacy policy explaining data collection
- Measurement: Code review, consent UI validation

**NFR12: Dependency Security**

- All npm dependencies scanned for known vulnerabilities
- Critical vulnerabilities patched within 7 days
- Measurement: npm audit, Snyk/Dependabot alerts

---

### Integration

**NFR13: Workspace Compatibility**

- Support npm workspaces (package.json with `workspaces` field)
- Support yarn workspaces (package.json with `workspaces` field)
- Support pnpm workspaces (pnpm-workspace.yaml)
- Measurement: Integration tests with real-world workspace examples

**NFR14: CI/CD Integration**

- CLI exit codes follow standard conventions (0 = pass, 1 = fail)
- Support for GitHub Actions, GitLab CI, CircleCI, Jenkins
- CI execution time < 2 minutes for typical monorepos (500 packages)
- Measurement: Integration test suites, user feedback

**NFR15: Export Formats**

- Support JSON export (machine-readable)
- Support HTML export (human-readable, offline viewing)
- Support Markdown export (documentation-friendly)
- All exports contain complete analysis results
- Measurement: Schema validation, visual inspection

---

### Scalability

**NFR16: User Growth Support**

- Infrastructure cost remains $0/month for Phase 0 (Cloudflare Pages free tier)
- Web application serves 10,000 concurrent users without performance degradation
- Measurement: Load testing, billing monitoring

**NFR17: Analysis Scalability**

- Graceful degradation for large monorepos (> 2000 packages)
- Chunked processing for memory efficiency (500 packages per batch)
- Clear error messages when analysis limits exceeded
- Measurement: Stress testing with synthetic large workspaces

---
