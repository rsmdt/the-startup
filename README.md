<p align="center">
  <img src="https://github.com/rsmdt/the-startup/blob/main/logo.png" width="400" alt="The Agentic Startup">
</p>

<p align="center">
  Ship faster. Ship better. Ship with <b>The Agentic Startup</b>.
</p>

<p align="center">
  <a href="https://github.com/rsmdt/the-startup/releases/latest">
    <img alt="Release" src="https://github.com/rsmdt/the-startup/actions/workflows/release.yml/badge.svg" />
  </a>

  <a href="https://github.com/rsmdt/the-startup/releases">
    <img alt="Downloads" src="https://img.shields.io/github/downloads/rsmdt/the-startup/total?style=flat&label=downloads&color=blue" />
  </a>

  <a href="https://github.com/rsmdt/the-startup/stargazers">
    <img alt="GitHub Stars" src="https://img.shields.io/github/stars/rsmdt/the-startup?style=flat&color=yellow" />
  </a>

  <a href="https://github.com/hesreallyhim/awesome-claude-code">
    <img alt="Mentioned in Awesome Claude Code" src="https://awesome.re/mentioned-badge.svg" />
  </a>
</p>

---

## Table of Contents

- [🤖 What is The Agentic Startup?](#-what-is-the-agentic-startup)
- [🚀 Quick Start](#-quick-start)
- [📖 The Complete Workflow](#-the-complete-workflow)
- [🎯 Which Command Should I Use?](#-which-command-should-i-use)
- [📦 Plugins](#-plugins)
- [🎨 Output Styles](#-output-styles)
- [🎯 Philosophy](#-philosophy)
- [📚 Documentation](#-documentation)

---

## 🤖 What is The Agentic Startup?

**The Agentic Startup** is a spec-driven development framework for Claude Code. Create comprehensive specifications before coding, then execute with parallel specialist agents—Y Combinator energy meets engineering discipline.

**Key Features:**
- **Native Claude Code Integration** — Marketplace plugins with zero configuration
- **Spec-Driven Development** — PRD → SDD → Implementation Plan → Code
- **Parallel Agent Execution** — Multiple specialists work simultaneously
- **Quality Gates** — Built-in validation at every stage

### Installation

**Requirements:** Claude Code v2.0+ with marketplace support

```bash
# Add The Agentic Startup marketplace
/plugin marketplace add rsmdt/the-startup

# Install the Start plugin (core workflows)
/plugin install start@the-startup

# (Optional) Install the Team plugin (specialized agents)
/plugin install team@the-startup

# Initialize your environment (statusline)
/start:init

# (Optional) Create project governance rules
/start:constitution                # Auto-enforced during specify, implement, review

# Choose your output style
/output-style start:The Startup    # High-energy, fast execution
/output-style start:The ScaleUp    # Calm confidence, educational
```

---

## 🚀 Quick Start

Create a specification and implement it:

```bash
# Create a specification
/start:specify Add user authentication with OAuth support

# Execute the implementation
/start:implement 001
```

That's it! You're now using spec-driven development.

---

## 📖 The Complete Workflow

The Agentic Startup follows **spec-driven development**: comprehensive specifications before code, ensuring clarity and reducing rework.

### All Commands at a Glance

```
┌──────────────────────────────────────────────────────────┐
│                    SETUP (one-time)                      │
│                                                          │
│  /start:init ───────► Configure statusline & environment │
│                                                          │
│  /start:constitution ► Create project governance rules   │
│                        (optional, auto-enforced in BUILD)│
└──────────────────────────────────────────────────────────┘
                            │
                            ▼
┌──────────────────────────────────────────────────────────┐
│                    BUILD (primary flow)                  │
│                                                          │
│  /start:specify ────► Create specs (PRD + SDD + PLAN)    │
│        │               ↳ Constitution checked on SDD     │
│        ▼                                                 │
│  /start:validate ───► Check quality (3 Cs framework)     │
│        │               ↳ Constitution mode available     │
│        ▼                                                 │
│  /start:implement ──► Execute plan phase-by-phase        │
│        │               ↳ Constitution + drift enforced   │
│        ▼                                                 │
│  /start:review ─────► Multi-agent code review            │
│        │               ↳ Constitution compliance checked │
│        ▼                                                 │
│  /start:document ───► Generate/sync documentation        │
└──────────────────────────────────────────────────────────┘
                            │
                            ▼
┌──────────────────────────────────────────────────────────┐
│                    MAINTAIN (as needed)                  │
│                                                          │
│  /start:analyze ────► Discover patterns & rules          │
│                                                          │
│  /start:refactor ───► Improve code (preserve behavior)   │
│                                                          │
│  /start:debug ──────► Fix bugs (root cause analysis)     │
└──────────────────────────────────────────────────────────┘
```

### Step-by-Step Walkthrough

#### Step 1: Create Your Specification

```bash
/start:specify Add real-time notification system with WebSocket support
```

This creates a specification directory with three documents:

```
docs/specs/001-notification-system/
├── product-requirements.md   # What to build and why
├── solution-design.md        # How to build it technically
└── implementation-plan.md    # Executable tasks and phases
```

**The spec cycle may take 15-30 minutes.** Claude will research your codebase, ask clarifying questions, and produce comprehensive documents. The process naturally involves multiple back-and-forth exchanges.

#### Step 2: Handle Context Limits (Resume Pattern)

Large specifications may approach Claude's context window limits. When this happens:

```bash
# Start a new conversation and resume where you left off
/start:specify 001
```

**The resume pattern:**
- Pass the spec ID (e.g., `001`) instead of a description
- Claude reads the existing spec files and continues from there
- You can reset context as many times as needed
- Each document (PRD → SDD → PLAN) can be completed in separate sessions if needed

**Pro tip:** If Claude suggests "you may want to reset context", do it! The quality of output improves with fresh context.

#### Step 3: Validate Before Implementation

```bash
/start:validate 001
```

This quality gate checks:
- **Completeness** - All sections filled, no missing details
- **Consistency** - No contradictions between documents
- **Correctness** - Requirements are testable and achievable

Validation is advisory—it provides recommendations but doesn't block you.

#### Step 4: Execute the Implementation

```bash
/start:implement 001
```

Claude will:
1. Parse the implementation plan
2. Execute phases sequentially (with your approval between phases)
3. Run tests after each task
4. Use parallel agents within phases for speed

**Large implementations may also need context resets.** Simply run `/start:implement 001` again in a fresh conversation—Claude tracks progress in the spec files.

#### Step 5: Review and Ship

```bash
/start:review
```

Four parallel specialists review your code:
- 🔒 **Security** - Authentication, authorization, input validation
- ⚡ **Performance** - Query optimization, memory management
- ✨ **Quality** - Code style, design patterns, maintainability
- 🧪 **Tests** - Coverage gaps, edge cases

---

## 🎯 Which Command Should I Use?

### Decision Tree

```
What do you need to do?
│
├─ First time setup? ─────────────────────► /start:init
│   └─ Want project-wide guardrails? ─────► Then: /start:constitution
│
├─ Build something new? ──────────────────► /start:specify
│                                           Then: /start:validate → /start:implement
│
├─ Understand existing code? ─────────────► /start:analyze
│   └─ Want to improve it? ───────────────► Then: /start:refactor
│
├─ Something is broken? ──────────────────► /start:debug
│
├─ Code ready for merge? ─────────────────► /start:review
│
├─ Need documentation? ───────────────────► /start:document
│
└─ Check constitution compliance? ────────► /start:validate constitution
```

### Command Reference

| Command | Purpose | When to Use |
|---------|---------|-------------|
| `/start:init` | Setup environment | First-time configuration |
| `/start:constitution` | Create governance rules | Establish project-wide guardrails |
| `/start:specify` | Create specifications | New features, complex changes |
| `/start:implement` | Execute plans | After spec is validated |
| `/start:validate` | Check quality | Before implementation, after specs |
| `/start:review` | Multi-agent code review | Before merging PRs |
| `/start:document` | Generate documentation | After implementation |
| `/start:analyze` | Extract knowledge | Understanding existing code |
| `/start:refactor` | Improve code quality | Cleanup without behavior change |
| `/start:debug` | Fix bugs | When something is broken |

### Capability Matrix

| Capability | constitution | specify | implement | validate | review | document | analyze | refactor | debug |
|------------|:------------:|:-------:|:---------:|:--------:|:------:|:--------:|:-------:|:--------:|:-----:|
| **Creates specifications** | - | ✅ | - | - | - | - | - | - | - |
| **Executes implementation plans** | - | - | ✅ | - | - | - | - | - | - |
| **Runs tests** | - | - | ✅ | ✅ | - | - | - | ✅ | ✅ |
| **Creates git branches** | - | ✅ | ✅ | - | - | - | - | ✅ | - |
| **Creates PRs** | - | ✅ | ✅ | - | - | - | - | - | - |
| **Multi-agent parallel** | - | ✅ | ✅ | - | ✅ | ✅ | ✅ | - | - |
| **Security scanning** | - | - | - | ✅ | ✅ | - | - | - | - |
| **Generates documentation** | - | ✅ | - | - | - | ✅ | ✅ | - | - |
| **Constitution enforcement** | ✅ | ✅ | ✅ | ✅ | ✅ | - | - | - | - |
| **Drift detection** | - | - | ✅ | - | - | - | - | - | - |

### When Commands Overlap

**validate vs review** — *Different purposes, different timing*

| Aspect | `/start:validate` | `/start:review` |
|--------|-------------------|-----------------|
| **When** | During development | Before merging |
| **Focus** | Spec compliance, quality gates | Code quality, security, performance |
| **Output** | Advisory recommendations | PR comments, findings report |

**analyze vs document** — *Discovery vs generation*

| Aspect | `/start:analyze` | `/start:document` |
|--------|------------------|-------------------|
| **Purpose** | Discover what exists | Generate documentation |
| **Output** | Knowledge documentation | API docs, READMEs, JSDoc |

**refactor vs debug** — *Improvement vs fixing*

| Aspect | `/start:refactor` | `/start:debug` |
|--------|-------------------|----------------|
| **Behavior** | Must preserve exactly | Expected to change (fix) |
| **Tests** | Must all pass throughout | May need new/updated tests |

---

## 📦 Plugins

The Agentic Startup is distributed as **Claude Code marketplace plugins**—native integration with zero manual configuration.

### Start Plugin (`start@the-startup`)

**Core workflow orchestration** — 10 commands, 18 skills, 2 output styles

| Category | Capabilities |
|----------|-------------|
| **Setup** | Environment configuration (`init`), project governance rules (`constitution`) |
| **Build** | `specify` → `validate` → `implement` pipeline with parallel agent coordination |
| **Quality** | Multi-agent code review, security scanning, constitution enforcement, drift detection |
| **Maintain** | Documentation generation, codebase analysis, safe refactoring, debugging |
| **Git** | Optional branch/commit/PR workflows integrated into commands |

**📖 [View detailed command documentation →](plugins/start/README.md)**

### Team Plugin (`team@the-startup`) — *Optional*

**Specialized agent library** — 11 roles, 27 activity-based specializations

| Role | Focus Areas |
|------|-------------|
| **Chief** | Complexity assessment, activity routing, parallel execution |
| **Analyst** | Requirements, prioritization, project coordination |
| **Architect** | System design, technology research, quality review, documentation |
| **Software Engineer** | APIs, components, domain modeling, performance |
| **QA Engineer** | Test strategy, exploratory testing, load testing |
| **Designer** | User research, interaction design, design systems, accessibility |
| **Platform Engineer** | IaC, containers, CI/CD, monitoring, data pipelines |
| **Meta Agent** | Agent design and generation |

**📖 [View all available agents →](plugins/team/README.md)**

---

## 🎨 Output Styles

The Start plugin includes two output styles that change how Claude communicates while working. Both maintain the same quality standards—the difference is in personality and explanation depth.

**Switch anytime:** `/output-style start:The Startup` or `/output-style start:The ScaleUp`

### The Startup 🚀

**High-energy execution with structured momentum.**

- **Vibe:** Demo day energy, Y Combinator intensity
- **Voice:** "Let's deliver this NOW!", "BOOM! That's what I'm talking about!"
- **Mantra:** "Done is better than perfect, but quality is non-negotiable"

**Best for:** Fast-paced sprints, high-energy execution, when you want momentum and celebration.

### The ScaleUp 📈

**Calm confidence with educational depth.**

- **Vibe:** Professional craft, engineering excellence
- **Voice:** "We've solved harder problems. Here's the approach."
- **Mantra:** "Sustainable speed at scale. We move fast, but we don't break things."

**Unique feature — Educational Insights:** The ScaleUp explains decisions as it works:

> 💡 *Insight: I used exponential backoff here because this endpoint has rate limiting. The existing `src/utils/retry.ts` helper already implements this pattern.*

**Best for:** Learning while building, understanding codebase patterns, onboarding to unfamiliar codebases.

### Comparison

| Dimension | The Startup | The ScaleUp |
|-----------|-------------|-------------|
| **Energy** | High-octane, celebratory | Calm, measured |
| **Explanations** | Minimal—ships fast | Educational insights included |
| **On failure** | "That didn't work. Moving on." | "Here's what failed and why..." |
| **Closing thought** | "What did we deliver?" | "Can the team maintain this?" |

---

## 🎯 Philosophy

### Why Activity-Based Agents?

Research shows **2-22% accuracy improvement** with specialized task agents vs. single broad agents ([Multi-Agent Collaboration, 2025](https://arxiv.org/html/2501.06322v1)). Leading frameworks organize agents by **capability**, not job titles. The Agentic Startup applies this research through activity-based specialization.

### The Problem We Solve

Development often moves too fast without proper planning:
- Features built without clear requirements
- Architecture decisions made ad-hoc during coding
- Technical debt accumulates from lack of upfront design
- Teams struggle to maintain consistency across implementations

### Our Approach: Spec-Driven Development

**1. Specify First** — Create comprehensive specifications before writing code
- **product-requirements.md** — What to build and why
- **solution-design.md** — How to build it technically
- **implementation-plan.md** — Executable tasks and phases

**2. Review & Refine** — Validate specifications with stakeholders
- Catch issues during planning, not during implementation
- Iterate on requirements and design cheaply
- Get alignment before costly development begins

**3. Implement with Confidence** — Execute validated plans phase-by-phase
- Clear acceptance criteria at every step
- Parallel agent coordination for speed
- Built-in validation gates and quality checks

**4. Document & Learn** — Capture patterns for future reuse
- Automatically document discovered patterns
- Build organizational knowledge base
- Prevent reinventing solutions

### Core Principles

**Measure twice, cut once** — Investing time in specifications saves exponentially more time during implementation.

**Documentation as code** — Specs, patterns, and interfaces are first-class artifacts that evolve with your codebase.

**Parallel execution** — Multiple specialists work simultaneously within clear boundaries, maximizing velocity without chaos.

**Quality gates** — Definition of Ready (DOR) and Definition of Done (DOD) ensure standards are maintained throughout.

**Progressive disclosure** — Skills and agents load details only when needed, optimizing token efficiency while maintaining power.

---

## 📚 Documentation

### Patterns

Reusable architectural patterns and design decisions:

| Pattern | Description |
|---------|-------------|
| [Slim Agent Architecture](docs/patterns/slim-agent-architecture.md) | Structure agents to maximize effectiveness while minimizing context usage |

### Additional Resources

- [Start Plugin Documentation](plugins/start/README.md) — Workflow commands and skills
- [Team Plugin Documentation](plugins/team/README.md) — Specialized agents and skills library
- [Migration Guide](MIGRATION.md) — Upgrading from v1.x

---

<p align="center">
  <strong>Ready to 10x your development workflow?</strong><br>
  Let's ship something incredible! 🚀
</p>
