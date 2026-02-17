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
- [🎯 Which Skill Should I Use?](#-which-skill-should-i-use)
- [📦 Plugins](#-plugins)
- [🎨 Output Styles](#-output-styles)
- [📊 Statusline](#-statusline)
- [💡 Why The Agentic Startup?](#-why-the-agentic-startup)
- [🎯 Philosophy](#-philosophy)
- [📚 Documentation](#-documentation)

---

> **New in v3:** Agent Teams (experimental) — enable multi-agent collaboration where specialized agents coordinate and work together on complex tasks. The installer now offers to configure this automatically.

---

## 🤖 What is The Agentic Startup?

**The Agentic Startup** is a multi-agent AI framework that makes Claude Code work like a startup team. Create comprehensive specifications before coding, then execute with parallel specialist agents — expert developers, architects, and engineers working together to turn your ideas into shipped code.

**10 slash commands across 3 phases.** Specify first, then build with confidence.

**Key Features:**
- **Spec-Driven Development** — PRD → SDD → Implementation Plan → Code
- **Parallel Agent Execution** — Multiple specialists work simultaneously
- **Quality Gates** — Built-in validation at every stage
- **Zero Configuration** — Marketplace plugins, one-line install

### Installation

**Requirements:** Claude Code v2.0+ with marketplace support

```bash
curl -fsSL https://raw.githubusercontent.com/rsmdt/the-startup/main/install.sh | sh
```

This installs the core plugins, configures the default output style, and sets up the [statusline](#-statusline) with a customizable config file.

<details>
<summary><strong>Manual Installation</strong></summary>

Start `claude` and run the following:

```bash
# Add The Agentic Startup marketplace
/plugin marketplace add rsmdt/the-startup

/plugin install start@the-startup  # Install the Start plugin (core workflows)
/plugin install team@the-startup   # (Optional) Install the Team plugin (specialized agents)
```

</details>

**After installation:**

```bash
# (Optional) Create project governance rules
/constitution                      # Auto-enforced during specify, implement, review

# Switch output styles anytime
/output-style "start:The Startup"   # High-energy, fast execution (default)
/output-style "start:The ScaleUp"   # Calm confidence, educational
```

---

## 🚀 Quick Start

Create a specification and implement it:

```bash
# Create a specification
/specify Add user authentication with OAuth support

# Execute the implementation
/implement 001
```

That's it! You're now using spec-driven development.

---

## 📖 The Complete Workflow

The Agentic Startup follows **spec-driven development**: comprehensive specifications before code, ensuring clarity and reducing rework.

### All Skills at a Glance

```
┌──────────────────────────────────────────────────────────┐
│                    SETUP (optional)                      │
│                                                          │
│  /constitution ► Create project governance rules         │
│                  (auto-enforced in BUILD workflow)       │
└──────────────────────────────────────────────────────────┘
                            │
                            ▼
┌──────────────────────────────────────────────────────────┐
│                    BUILD (primary flow)                  │
│                                                          │
│  /specify ────► Create specs (PRD + SDD + PLAN)          │
│      │           ↳ Constitution checked on SDD           │
│      ▼                                                   │
│  /validate ───► Check quality (3 Cs framework)           │
│      │           ↳ Constitution mode available           │
│      ▼                                                   │
│  /implement ──► Execute plan phase-by-phase              │
│      │           ↳ Constitution + drift enforced         │
│      ▼                                                   │
│  /test ───────► Run tests, enforce ownership             │
│      │           ↳ No "pre-existing" excuses             │
│      ▼                                                   │
│  /review ─────► Multi-agent code review                  │
│      │           ↳ Constitution compliance checked       │
│      ▼                                                   │
│  /document ───► Generate/sync documentation              │
└──────────────────────────────────────────────────────────┘
                            │
                            ▼
┌──────────────────────────────────────────────────────────┐
│                    MAINTAIN (as needed)                  │
│                                                          │
│  /analyze ────► Discover patterns & rules                │
│                                                          │
│  /refactor ───► Improve code (preserve behavior)         │
│                                                          │
│  /debug ──────► Fix bugs (root cause analysis)           │
└──────────────────────────────────────────────────────────┘
```

### Step-by-Step Walkthrough

#### Step 1: Create Your Specification

```bash
/specify Add real-time notification system with WebSocket support
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
/specify 001
```

**The resume pattern:**
- Pass the spec ID (e.g., `001`) instead of a description
- Claude reads the existing spec files and continues from there
- You can reset context as many times as needed
- Each document (PRD → SDD → PLAN) can be completed in separate sessions if needed

**Pro tip:** If Claude suggests "you may want to reset context", do it! The quality of output improves with fresh context.

#### Step 3: Validate Before Implementation

```bash
/validate 001
```

This quality gate checks:
- **Completeness** - All sections filled, no missing details
- **Consistency** - No contradictions between documents
- **Correctness** - Requirements are testable and achievable

Validation is advisory—it provides recommendations but doesn't block you.

#### Step 4: Execute the Implementation

```bash
/implement 001
```

Claude will:
1. Parse the implementation plan
2. Execute phases sequentially (with your approval between phases)
3. Run tests after each task
4. Use parallel agents within phases for speed

**Large implementations may also need context resets.** Simply run `/implement 001` again in a fresh conversation—Claude tracks progress in the spec files.

#### Step 5: Review and Ship

```bash
/review
```

Four parallel specialists review your code:
- 🔒 **Security** - Authentication, authorization, input validation
- ⚡ **Performance** - Query optimization, memory management
- ✨ **Quality** - Code style, design patterns, maintainability
- 🧪 **Tests** - Coverage gaps, edge cases

---

## 🎯 Which Skill Should I Use?

### Decision Tree

```
What do you need to do?
│
├─ Want project-wide guardrails? ─────────► /constitution
│
├─ Build something new? ──────────────────► /specify
│                                           Then: /validate → /implement
│
├─ Understand existing code? ─────────────► /analyze
│   └─ Want to improve it? ───────────────► Then: /refactor
│
├─ Something is broken? ──────────────────► /debug
│
├─ Need to run tests? ───────────────────► /test
│
├─ Code ready for merge? ─────────────────► /review
│
├─ Need documentation? ───────────────────► /document
│
└─ Check constitution compliance? ────────► /validate constitution
```

### Skill Reference

| Skill | Purpose | When to Use |
|---------|---------|-------------|
| `/constitution` | Create governance rules | Establish project-wide guardrails |
| `/specify` | Create specifications | New features, complex changes |
| `/implement` | Execute plans | After spec is validated |
| `/validate` | Check quality | Before implementation, after specs |
| `/test` | Run tests, enforce ownership | After implementation, fixing bugs |
| `/review` | Multi-agent code review | Before merging PRs |
| `/document` | Generate documentation | After implementation |
| `/analyze` | Extract knowledge | Understanding existing code |
| `/refactor` | Improve code quality | Cleanup without behavior change |
| `/debug` | Fix bugs | When something is broken |

### Capability Matrix

| Capability | constitution | specify | implement | validate | test | review | document | analyze | refactor | debug |
|------------|:------------:|:-------:|:---------:|:--------:|:----:|:------:|:--------:|:-------:|:--------:|:-----:|
| **Creates specifications** | - | ✅ | - | - | - | - | - | - | - | - |
| **Executes implementation plans** | - | - | ✅ | - | - | - | - | - | - | - |
| **Runs tests** | - | - | ✅ | ✅ | ✅ | - | - | - | ✅ | ✅ |
| **Creates git branches** | - | ✅ | ✅ | - | - | - | - | - | ✅ | - |
| **Creates PRs** | - | ✅ | ✅ | - | - | - | - | - | - | - |
| **Multi-agent parallel** | - | ✅ | ✅ | - | ✅ | ✅ | ✅ | ✅ | - | - |
| **Security scanning** | - | - | - | ✅ | - | ✅ | - | - | - | - |
| **Generates documentation** | - | ✅ | - | - | - | - | ✅ | ✅ | - | - |
| **Constitution enforcement** | ✅ | ✅ | ✅ | ✅ | - | ✅ | - | - | - | - |
| **Drift detection** | - | - | ✅ | - | - | - | - | - | - | - |
| **Code ownership enforcement** | - | - | - | - | ✅ | - | - | - | - | - |

### When Skills Overlap

**validate vs review** — *Different purposes, different timing*

| Aspect | `/validate` | `/review` |
|--------|-------------------|-----------------|
| **When** | During development | Before merging |
| **Focus** | Spec compliance, quality gates | Code quality, security, performance |
| **Output** | Advisory recommendations | PR comments, findings report |

**analyze vs document** — *Discovery vs generation*

| Aspect | `/analyze` | `/document` |
|--------|------------------|-------------------|
| **Purpose** | Discover what exists | Generate documentation |
| **Output** | Knowledge documentation | API docs, READMEs, JSDoc |

**refactor vs debug** — *Improvement vs fixing*

| Aspect | `/refactor` | `/debug` |
|--------|-------------------|----------------|
| **Behavior** | Must preserve exactly | Expected to change (fix) |
| **Tests** | Must all pass throughout | May need new/updated tests |

---

## 📦 Plugins

The Agentic Startup is distributed as **Claude Code marketplace plugins**—native integration with zero manual configuration.

### Start Plugin (`start@the-startup`)

**Core workflow orchestration** — 10 user-invocable skills, 5 autonomous skills, 2 output styles

| Category | Capabilities |
|----------|-------------|
| **Setup** | Environment configuration (`init`), project governance rules (`constitution`) |
| **Build** | `specify` → `validate` → `implement` pipeline with parallel agent coordination |
| **Quality** | Multi-agent code review, security scanning, constitution enforcement, drift detection |
| **Maintain** | Documentation generation, codebase analysis, safe refactoring, debugging |
| **Git** | Optional branch/commit/PR workflows integrated into skills |

**📖 [View detailed skill documentation →](plugins/start/README.md)**

### Team Plugin (`team@the-startup`) — *Optional*

**Specialized agent library** — 8 roles, 20 activity-based agents. Now with experimental [Agent Teams](#agent-teams-experimental--new-in-v3) support for multi-agent collaboration.

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

## 🔧 How Skills Work

The Agentic Startup is built on Claude Code's [skills system](https://code.claude.com/docs/en/skills), which follows the [Agent Skills](https://agentskills.io) open standard. Understanding how skills are invoked helps you get the most out of the framework.

### Invocation Model

Skills have two invocation paths, controlled by frontmatter fields in each skill's `SKILL.md`:

| Path | How It Works | Controlled By |
|------|-------------|---------------|
| **User slash command** | You type `/skill-name [args]` | `user-invocable` (default: `true`) |
| **Model auto-invocation** | Claude detects context and loads the skill via the Skill tool | `disable-model-invocation` (default: `false`) |

Skills from the Start plugin are invoked directly by name (e.g., `/specify`, `/test`).

### User-Invocable vs Autonomous Skills

| Type | Visible in `/` menu? | Claude auto-invokes? | Example |
|------|:--------------------:|:--------------------:|---------|
| **User-invocable** | Yes | Yes | `/specify` — you trigger the spec workflow |
| **Autonomous** | No | Yes | `specify-requirements` — loaded by `specify` when creating PRDs |

The 10 user-invocable skills are the ones you interact with directly. The 5 autonomous skills activate behind the scenes when orchestrator skills need them (e.g., `specify` loads `specify-requirements`, `specify-solution`, and `specify-plan` during the specification workflow).

### Progressive Disclosure

Skills load efficiently to conserve context:

1. **At startup** — Only skill names and descriptions are loaded (~100 tokens each)
2. **On invocation** — Full `SKILL.md` content loads when you or Claude triggers the skill
3. **On demand** — Supporting files (`reference.md`, templates, scripts) load only when needed

This means all 15 skills can be available without consuming significant context until actually used.

---

## 📊 Statusline

The installer sets up a custom statusline that displays context usage, session cost, and other useful information directly in your Claude Code terminal.

### What You See

```
📁 ~/C/p/project ⎇ main*  🤖 Opus 4.5 (The Startup)  🧠 ⣿⣿⡇⠀⠀ 50%  🕐 30m  💰 $1.50  ? for shortcuts
```

| Component | Description |
|-----------|-------------|
| 📁 `~/C/p/project` | Current directory (abbreviated) |
| ⎇ `main*` | Git branch (* indicates uncommitted changes) |
| 🤖 `Opus 4.5 (The Startup)` | Model and output style |
| 🧠 `⣿⣿⡇⠀⠀ 50%` | Context window usage (color-coded) |
| 🕐 `30m` | Session duration |
| 💰 `$1.50` | Session cost (color-coded by plan) |

### Color Thresholds

Both context usage and cost display color-coded warnings:

| Color | Context | Cost (Pro plan) |
|-------|---------|-----------------|
| 🟢 Green | < 70% | < $1.50 |
| 🟡 Amber | 70-89% | $1.50 - $4.99 |
| 🔴 Red | ≥ 90% | ≥ $5.00 |

### Configuration

The statusline reads from `~/.config/the-agentic-startup/statusline.toml`:

```toml
# Format string (customize what's displayed)
format = "<path> <branch>  <model>  <context>  <session>  <help>"

# Plan for cost thresholds: "auto" | "pro" | "max5x" | "max20x" | "api"
plan = "auto"
fallback_plan = "pro"

[thresholds.context]
warn = 70    # percentage
danger = 90

[thresholds.cost]
# Uncomment to override plan defaults:
# warn = 2.00
# danger = 5.00
```

### Plan-Based Cost Defaults

| Plan | Monthly | Warn | Danger |
|------|---------|------|--------|
| `pro` | $20 | $1.50 | $5.00 |
| `max5x` | $100 | $5.00 | $15.00 |
| `max20x` | $200 | $10.00 | $30.00 |
| `api` | Pay-as-you-go | $2.00 | $10.00 |

### Format Placeholders

| Placeholder | Description | Example |
|-------------|-------------|---------|
| `<path>` | Abbreviated directory | `~/C/p/project` |
| `<branch>` | Git branch with dirty indicator | `⎇ main*` |
| `<model>` | Model and output style | `🤖 Opus 4.5 (The Startup)` |
| `<context>` | Context usage bar and percentage | `🧠 ⣿⣿⡇⠀⠀ 50%` |
| `<session>` | Duration and cost | `🕐 30m  💰 $1.50` |
| `<lines>` | Lines added/removed | `+156/-23` |
| `<spec>` | Active spec ID (when in docs/specs/) | `📋 005` |
| `<help>` | Help text | `? for shortcuts` |

**Example minimal format:**
```toml
format = "<context>  <session>"
```

---

## 💡 Why The Agentic Startup?

Real workflow features that solve real problems — not just another AI wrapper.

### Resume Across Sessions

Hit a context limit? Start a new conversation and pick up exactly where you left off. Specs persist on disk — Claude reads them and continues.

```bash
/specify 001    # ← resumes spec creation from where you left off
/implement 001  # ← resumes implementation, tracking progress in spec files
```

### Code Ownership Mandate

No more "pre-existing failure" excuses. When `/test` finds a failing test, it fixes it — period. You touched the codebase, you own it.

### Drift Detection

Implementation drifting from the spec? Caught automatically during `/implement`. Scope creep, missing items, contradictions — flagged with options to update the spec or the code.

### Adaptive Code Review

`/review` auto-detects what matters. Async code triggers concurrency review. Dependency changes trigger supply-chain checks. UI changes trigger accessibility audits. 5 base perspectives + conditional specialists.

### Implement Any Plan

Not just for specs created with `/specify`. `/implement` works with any markdown implementation plan — bring your own architecture docs, migration guides, or design documents.

```bash
/implement path/to/plan.md
```

### Non-Linear Specs

Skip what you don't need. Start with a solution design, jump to the plan, or go full PRD → SDD → PLAN. Skipped phases are logged as decisions, not gaps.

### Adversarial Debugging

Tough bugs get multiple investigators that actively try to disprove each other's hypotheses. The surviving theory is most likely the root cause — competing hypotheses, not confirmation bias.

### Agent Teams (Experimental) — New in v3

Enable multi-agent collaboration where specialized agents coordinate autonomously on complex tasks. The installer configures this automatically, or enable manually:

```json
// ~/.claude/settings.json
{
  "env": {
    "CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS": "1"
  }
}
```

---

## 🎯 Philosophy

Research shows **2-22% accuracy improvement** with specialized task agents vs. single broad agents ([Multi-Agent Collaboration, 2025](https://arxiv.org/html/2501.06322v1)). Leading frameworks organize agents by **capability**, not job titles. The Agentic Startup applies this research through activity-based specialization.

### The Problem We Solve

Development often moves too fast without proper planning:
- Features built without clear requirements
- Architecture decisions made ad-hoc during coding
- Technical debt accumulates from lack of upfront design
- Teams struggle to maintain consistency across implementations

### Our Approach

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

- **Measure twice, cut once** — Investing time in specifications saves exponentially more time during implementation.
- **Documentation as code** — Specs, patterns, and interfaces are first-class artifacts that evolve with your codebase.
- **Parallel execution** — Multiple specialists work simultaneously within clear boundaries, maximizing velocity without chaos.
- **Quality gates** — Definition of Ready (DOR) and Definition of Done (DOD) ensure standards are maintained throughout.
- **Progressive disclosure** — Skills and agents load details only when needed, optimizing token efficiency while maintaining power.

---

## 📚 Documentation

### Patterns

Reusable architectural patterns and design decisions:

| Pattern | Description |
|---------|-------------|
| [Slim Agent Architecture](docs/patterns/slim-agent-architecture.md) | Structure agents to maximize effectiveness while minimizing context usage |

### Additional Resources

- [Start Plugin Documentation](plugins/start/README.md) — Workflow skills
- [Team Plugin Documentation](plugins/team/README.md) — Specialized agents and skills library
- [Migration Guide](MIGRATION.md) — Upgrading from v1.x

---

<p align="center">
  <strong>Ready to 10x your development workflow?</strong><br>
  Let's ship something incredible! 🚀
</p>
