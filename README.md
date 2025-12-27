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


## 🤖 What is The Agentic Startup?

**The Agentic Startup** is a spec-driven development framework for Claude Code that transforms how you build software. Think Y Combinator demo day energy meets engineering discipline - multiple specialists executing in parallel while you stay in control.

We follow **[Spec-Driven Development](https://www.perplexity.ai/?q=Spec+Driven+Development)**, where comprehensive specifications are created before coding begins, ensuring clarity and reducing rework.

**The workflow**:

1. **📋 Specify** - Turn ideas into comprehensive specification documents
2. **✅ Validate** - Quality gate to check specs, implementations, or understanding
3. **⚡ Implement** - Execute plans phase-by-phase with parallel agent coordination

The framework leverages **Claude Code's plugin system** to provide workflow commands, autonomous skills, specialized agents, and quality templates - all working together like a high-velocity startup team.

---

## ✨ Key Features

**Native Claude Code Integration** - Distributed as official marketplace plugins with zero manual configuration

**Autonomous Skills System** - Model-invoked skills that activate based on natural language with progressive disclosure

**Specialized Agent Team** - 11 agent roles across 27 activity-based specializations (optional [`team@the-startup`](plugins/team/README.md) plugin)

**Structured Momentum** - Y Combinator energy meets engineering discipline - parallel execution with quality gates

---

## 📦 Plugins

### Start Plugin (`start@the-startup`)

**Core workflow orchestration** — 9 commands, 16 skills, 2 output styles

| Category | Capabilities |
|----------|-------------|
| **Build** | `specify` → `validate` → `implement` pipeline with parallel agent coordination |
| **Quality** | Multi-agent code review, security scanning, test coverage checks |
| **Maintain** | Documentation generation, codebase analysis, safe refactoring, debugging |
| **Git** | Optional branch/commit/PR workflows integrated into commands |

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

---

## 🚀 Quick Start

### Installation

**Requirements**: Claude Code v2.0+ with marketplace support

```bash
# Add The Agentic Startup marketplace
/plugin marketplace add rsmdt/the-startup

# Install the Start plugin (core workflows)
/plugin install start@the-startup

# (Optional) Install the Team plugin (specialized agents)
/plugin install team@the-startup
```

**📖 [View all available agents →](plugins/team/README.md)**

Alternatively, browse and install interactively via `/plugin`

### Initialize Your Environment

Configure output style and statusline (one-time setup):

```bash
/start:init
```

### Your First Workflow

**1. Create a specification:**
```bash
/start:specify Add user authentication with OAuth support
```

Creates `docs/specs/001-user-authentication/` with product-requirements.md, solution-design.md, and implementation-plan.md documents.

**2. Execute the implementation:**
```bash
/start:implement 001
```

Runs phase-by-phase with parallel agents, validation gates, and progress tracking.

---

## 📋 Commands

Quick reference for all workflow commands:

| Command | Description |
|---------|-------------|
| `/start:init` | Initialize environment (output style, statusline) |
| `/start:specify` | Create specification documents from brief description |
| `/start:implement` | Execute implementation plan phase-by-phase |
| `/start:validate` | Validate specs, implementations, or understanding (quality gate) |
| `/start:analyze` | Discover and document patterns, rules, interfaces |
| `/start:refactor` | Improve code quality while preserving behavior |
| `/start:debug` | Conversational debugging with systematic root cause analysis |

**📖 [View detailed command documentation →](plugins/start/README.md)**

---

## 🔄 Typical Development Workflow

### Specify → Validate → Implement

**1. Create Specification**

```bash
/start:specify Add real-time notification system with WebSocket support
```

- Creates comprehensive specs in `docs/specs/001-notification-system/`
- Documents discovered patterns and interfaces
- Duration: 15-30 minutes

**2. Validate Before Implementation (Optional)**

```bash
/start:validate 001
```

- Checks completeness, consistency, correctness
- Detects ambiguities and gaps
- Provides advisory recommendations
- Duration: 2-5 minutes

**3. Execute Implementation**

```bash
/start:implement 001
```

- Executes phases sequentially with user approval
- Parallel agent coordination within phases
- Continuous test validation
- Duration: Varies by complexity

### Separate Workflows

**Validate understanding or implementation:**
```bash
/start:validate Check the current auth implementation against the SDD
```

**Analyze existing code:**
```bash
/start:analyze security patterns in authentication
```

**Refactor code safely:**
```bash
/start:refactor Simplify the WebSocket connection manager
```

**Debug issues conversationally:**
```bash
/start:debug The API returns 500 errors when uploading large files
```

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

**The Agentic Startup** enforces a disciplined workflow that balances speed with quality:

**1. Specify First** - Create comprehensive specifications before writing code
- **product-requirements.md** - What to build and why
- **solution-design.md** - How to build it technically
- **implementation-plan.md** - Executable tasks and phases

**2. Review & Refine** - Validate specifications with stakeholders
- Catch issues during planning, not during implementation
- Iterate on requirements and design cheaply
- Get alignment before costly development begins

**3. Implement with Confidence** - Execute validated plans phase-by-phase
- Clear acceptance criteria at every step
- Parallel agent coordination for speed
- Built-in validation gates and quality checks

**4. Document & Learn** - Capture patterns for future reuse
- Automatically document discovered patterns
- Build organizational knowledge base
- Prevent reinventing solutions

### Core Principles

**Measure twice, cut once** - Investing time in specifications saves exponentially more time during implementation.

**Documentation as code** - Specs, patterns, and interfaces are first-class artifacts that evolve with your codebase.

**Parallel execution** - Multiple specialists work simultaneously within clear boundaries, maximizing velocity without chaos.

**Quality gates** - Definition of Ready (DOR) and Definition of Done (DOD) ensure standards are maintained throughout.

**Progressive disclosure** - Skills and agents load details only when needed, optimizing token efficiency while maintaining power.

---

## 📚 Documentation

### Patterns

Reusable architectural patterns and design decisions:

| Pattern | Description |
|---------|-------------|
| [Slim Agent Architecture](docs/patterns/slim-agent-architecture.md) | Structure agents to maximize effectiveness while minimizing context usage |

### Additional Resources

- [Start Plugin Documentation](plugins/start/README.md) - Workflow commands and skills
- [Team Plugin Documentation](plugins/team/README.md) - Specialized agents and skills library
- [Migration Guide](MIGRATION.md) - Upgrading from v1.x

---

<p align="center">
  <strong>Ready to 10x your development workflow?</strong><br>
  Let's ship something incredible! 🚀
</p>
