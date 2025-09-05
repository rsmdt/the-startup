<p align="center"><img src="https://github.com/rsmdt/the-startup/blob/main/assets/logo.png" width="400" alt="The Agentic Startup"></p>

<p align="center">Ship faster. Ship better. Ship with <b>The Agentic Startup</b>.</p>

## What is The Agentic Startup?

The Agentic Startup brings you instant access to expert developers, architects, and engineers - all working together to turn your ideas into shipped code. It is a system for [Claude Code](https://www.anthropic.com/claude-code) that gives you a virtual engineering team. Instead of one AI trying to do everything, you get specialized experts who collaborate like a real startup team - pragmatic, fast, and focused on shipping.

Think of it as having a CTO, architects, developers, and DevOps engineers on-demand, each bringing their expertise to your project.

### Core Philosophy

**Think twice, ship once.** Proper planning accelerates delivery more than jumping straight into code.

- **Humans decide, AI executes** - Critical decisions stay with you; AI handles implementation details
- **Specialist delegation** - Pull in the right expert for each task
- **Documentation drives clarity** - Specs prevent miscommunication and scope creep
- **Parallel execution** - Multiple experts work simultaneously when possible
- **Review everything** - No AI decision goes unreviewed; you stay in control

When you use The Agentic Startup, Claude Code becomes your **technical co-founder** that gathers context first, consults specialists, generates reviewable documentation, then implements with confidence.

### Research Foundation

Task specialization consistently outperforms role-based organization for LLM agents:

- Performance Impact: Studies show 2.86% to 21.88% accuracy improvement with specialized agents vs single broad agents ([Multi-Agent Collaboration, 2025](https://arxiv.org/html/2501.06322v1))
- Industry Consensus: Leading frameworks (CrewAI, Microsoft AutoGen, LangGraph) organize agents by **capability** rather than traditional job titles
- Domain Specialization: Effective LLM specialization customizes agents according to specific task contextual data ([Agentic LLM Systems, 2024](https://arxiv.org/html/2412.04093v1))

## Quick Start

Install and start using The Agentic Startup:

```bash
# Install (interactive)
curl -LsSf https://raw.githubusercontent.com/rsmdt/the-startup/main/install.sh | sh

# Activate The Agentic Startup output style for the full experience
/output-style The Startup

# Plan a feature
/s:specify "Add user authentication"

# Build it 
/s:implement 001-user-auth
```

**More installation options**: See the [Installation](#installation) section below.

## Agents - Your Expert Team

The Agentic Startup uses **activity-based agents** that focus on WHAT they do, not WHO they are. Traditional engineering boundaries (backend/frontend) are artificial constraints that reduce LLM performance. Instead, our agents:

- **Focus on activities** - Agents specialize in `api-design` or `component-architecture`, not arbitrary roles
- **Adapt to your stack** - Automatically detect and apply React/Vue/Angular patterns, REST/GraphQL APIs, PostgreSQL/MongoDB optimizations
- **Execute in parallel** - Multiple specialists work simultaneously on related activities
- **Preserve real expertise** - Keep domain specialization (mobile, security, UX) where it genuinely adds value

Each agent receives only relevant context for their specific expertise, reducing cognitive load and improving accuracy.

**the-chief** - Eliminates bottlenecks through smart routing and complexity assessment

**the-analyst** - Transforms vague requirements into actionable specifications
- `requirements-clarification` - Uncovers hidden needs and resolves ambiguities
- `requirements-documentation` - Creates comprehensive BRDs and PRDs
- `feature-prioritization` - Data-driven feature prioritization
- `solution-research` - Researches proven approaches and patterns
- `project-coordination` - Breaks down complex projects into tasks

**the-architect** - Balances elegance with pragmatic business reality
- `system-design` - Designs scalable system architectures
- `system-documentation` - Creates architecture diagrams and decisions
- `architecture-review` - Validates design patterns and compliance
- `code-review` - Elevates team capabilities through feedback
- `scalability-planning` - Ensures systems scale gracefully
- `technology-evaluation` - Makes framework and tool decisions
- `technology-standards` - Prevents technology chaos through standards

**the-software-engineer** - Ships features that actually work
- `api-design` - REST/GraphQL APIs with clear contracts
- `api-documentation` - Comprehensive API documentation
- `database-design` - Balanced schemas for any database
- `service-integration` - Reliable service communication patterns
- `component-architecture` - Reusable UI components
- `business-logic` - Domain rules and validation
- `reliability-engineering` - Error handling and resilience
- `performance-optimization` - Bundle size and Core Web Vitals
- `state-management` - Client and server state patterns
- `browser-compatibility` - Cross-browser support

**the-platform-engineer** - Makes systems that don't wake you at 3am
- `system-performance` - Handles 10x load without 10x cost
- `observability` - Monitoring that catches problems early
- `containerization` - Consistent deployment everywhere
- `pipeline-engineering` - Reliable data processing
- `ci-cd-automation` - Safe deployments at scale
- `deployment-strategies` - Progressive rollouts
- `incident-response` - Production fire debugging
- `infrastructure-as-code` - Reproducible infrastructure
- `storage-architecture` - Scalable storage solutions
- `query-optimization` - Fast database queries
- `data-modeling` - Balanced data models

**the-designer** - Creates products people actually want to use
- `accessibility-implementation` - WCAG 2.1 AA compliance
- `user-research` - Real user needs, not assumptions
- `interaction-design` - Minimal friction user flows
- `visual-design` - Brand-enhancing UI aesthetics
- `design-systems` - Consistent component libraries
- `information-architecture` - Intuitive content hierarchies

**the-qa-engineer** - Catches bugs before users do
- `test-strategy` - Risk-based testing approaches
- `test-implementation` - Comprehensive test suites
- `exploratory-testing` - Creative defect discovery
- `performance-testing` - Load and stress validation

**the-security-engineer** - Keeps the bad guys out
- `vulnerability-assessment` - OWASP-based security checks
- `authentication-systems` - OAuth, JWT, SSO, MFA
- `security-incident-response` - Rapid containment
- `compliance-audit` - GDPR, SOX, HIPAA compliance
- `data-protection` - Encryption and privacy controls

**the-mobile-engineer** - Ships apps users love
- `mobile-interface-design` - Platform-specific UI patterns
- `mobile-data-persistence` - Offline-first strategies
- `cross-platform-integration` - Native and hybrid bridges
- `mobile-deployment` - App store submissions
- `mobile-performance` - Battery and memory optimization

**the-ml-engineer** - Makes AI that actually ships
- `model-deployment` - Production-ready inference
- `ml-monitoring` - Drift detection systems
- `prompt-optimization` - LLM prompt engineering
- `mlops-automation` - Reproducible ML pipelines
- `context-management` - AI memory architectures
- `feature-engineering` - Model-ready data pipelines

**the-meta-agent** - Creates new specialized agents

## Slash Commands

The Startup provides powerful slash commands that orchestrate your entire development workflow. Each command features built-in verification checkpoints and mandatory pause points to ensure quality at every step.

### `/s:specify` - Plan Before You Build

Creates comprehensive specifications with built-in quality gates:

```bash
# Start fresh with a new feature idea
/s:specify Build a real-time notification system

# Resume working on a specification
/s:specify 001
```

**Documents Created:**
- `docs/specs/[id]-[short-name]/BRD.md` - Business Requirements Document capturing the "why" and business value
- `docs/specs/[id]-[short-name]/PRD.md` - Product Requirements Document defining user-facing features and acceptance criteria
- `docs/specs/[id]-[short-name]/SDD.md` - Solution Design Document detailing technical architecture and implementation approach
- `docs/specs/[id]-[short-name]/PLAN.md` - Implementation Plan with phase-by-phase tasks ready for execution
- `docs/patterns/` - Documents reusable patterns discovered during research (authentication flows, caching strategies, etc.)
- `docs/interfaces/` - Documents external API contracts and integration specifications documented along the way

**Key Features:**
- 🤔 Self-verification checkpoints - "Ask yourself" prompts ensure thorough analysis
- 🛑 Phase boundaries - User approval required at each major step
- ⚡ Parallel research - Multiple specialists investigate simultaneously

#### Workflow

<details>
<summary>show details</summary>

```mermaid
flowchart TD
    A([Your Feature Idea]) --> |initialize| B{Check<br>Existing}
    B --> |exists| C[Review and Refine]
    C --> END[🚀 Ready for /s:implement S001]

    B --> |new| D[📄 **Requirements Gathering**<br/>Create *BRD.md*, *PRD.md* if needed]
    D --> E[📄 **Technical Research**<br/>Create *SDD.md* if needed, document patterns, interfaces]
    E --> F[📄 **Implementation Planning**<br/>Create *PLAN.md*]

    F --> END[🚀 Ready for /s:implement S001]
```

</details>

### `/s:implement` - Execute the Plan

Takes an implementation plan (PLAN.md) and executes it phase-by-phase with expert delegation:

```bash
# Implement a completed specification (requires PLAN.md)
/s:implement 001

# Implement from a specific PLAN.md file
/s:implement docs/specs/001-auth/PLAN.md

# Use your own plan document
/s:implement my-custom-plan.md
```

**Requirements:**
- Depends on a PLAN.md document (created by `/s:specify` or your own)
- Plan must include phase markers and task lists for execution
- Can use any properly formatted plan document, not just generated ones

**Key Features:**
- 📋 Phase-by-phase execution - One phase at a time to prevent overload
- ⚡ Parallel task execution - Multiple agents work simultaneously within phases
- 🛑 Phase boundaries - Mandatory stops between phases for review
- 🔍 Automatic validation - Tests run after each change

#### Workflow

<details>
<summary>show details</summary>

```mermaid
flowchart TD
    A([📄 *PLAN.md*]) --> |load| B[**Initialize Plan**<br/>Parse phases & tasks]
    B --> |approve| C{Phases<br>Remaining?}
    
    C --> |yes| D[**Execute Phase N**<br/>⚡ *Parallel agent execution*<br/>✓ *Run tests after each task*]
    D --> |validate| E[**Phase Review**<br/>Check test results<br/>Review changes]
    E --> |continue| C
    
    C --> |no| F[**Final Validation**<br/>Run full test suite<br/>Verify all requirements]
    F --> END[✅ **Implementation Complete**]
```

</details>

### `/s:refactor` - Improve Code Quality

Analyzes code and performs refactoring based on complexity assessment:

```bash
# Refactor specific code or modules
/s:refactor improve the authentication module for better testability

# Refactor for specific goals
/s:refactor reduce complexity in the payment processing logic
```

**Complexity-Based Behavior:**
- **Simple refactoring** → Executes immediately with validation at each step
  - Method extraction, variable renaming, small scope changes
  - Direct execution with continuous test validation
  
- **Complex refactoring** → Creates specification for later execution
  - Architectural changes, cross-module refactoring, API redesigns
  - Generates `SDD.md` and `PLAN.md` for review before execution
  - Use `/s:implement` to execute the refactoring plan

**Key Features:**
- 🎯 Goal clarification - Ensures refactoring objectives are clear
- 🔍 Validation-first - Tests must pass before and after changes
- 🔀 Complexity routing - Automatic decision between immediate or planned execution
- 🛑 Safety checkpoints - User approval at critical decision points

#### Workflow

<details>
<summary>show details</summary>

```mermaid
flowchart TD
    A([Refactoring Request]) --> |analyze| B[**Goal Clarification**<br/>Define objectives<br/>Analyze codebase]
    B --> |assess| C{**Complexity<br>Check**}
    
    C --> |simple| D[**Direct Refactoring**<br/>✓ *Run tests first*<br/>🔧 *Apply changes*<br/>✓ *Validate each step*]
    D --> |review| E[**Specialist Review**<br/>Code quality check<br/>Performance impact]
    E --> DONE[✅ **Refactoring Complete**]
    
    C --> |complex| F[**Create Specification**<br/>📄 *Generate SDD.md*<br/>📄 *Generate PLAN.md*<br/>Document approach]
    F --> |defer| G[🚀 **Ready for /s:implement**<br/>Execute via planned phases]
```

</details>

## 🎯 The Startup Output Style

For the most immersive experience, activate **The Startup** output style to transform Claude into your high-energy technical co-founder.

```bash
/output-style The Startup
```

**What you get:**
- 🚀 Startup energy - "Let's ship this NOW!" enthusiasm in every response
- ⚡ Parallel execution - Launches multiple agents simultaneously, no blocking
- 📊 Task tracking - Uses TodoWrite obsessively for progress visibility
- 🎉 Victory celebrations - Acknowledges every shipped feature

**Example transformation:**
```
Standard: "I'll help you implement authentication..."
The Agentic Startup: "🚀 TIME TO SHIP! Launching the security squad in parallel!"
```

The style makes every session feel like you're building the next unicorn.

## Quick Start Examples

```bash
# Fix a bug
"Error: Cannot read property 'user' of undefined in auth.js"

# Build a feature  
/s:specify "Add CSV export functionality to reports"

# Optimize performance
"The dashboard takes 10 seconds to load"

# Review code
"Review my authentication implementation for security issues"

# Get unstuck
"I don't know how to structure this microservices architecture"
```

## Hooks - Track Your AI Team's Work

The Agentic Startup automatically installs hooks that capture **every tool invocation** Claude Code makes, giving you complete visibility into your AI team's activity. No more wondering what Claude did or how long it took.

### Automatic Metrics Collection

Every time Claude uses a tool (Read, Write, Bash, etc.), The Startup captures:
- **Tool usage patterns** - Which tools are used most frequently
- **Success/failure rates** - Track reliability of different operations
- **Performance metrics** - Execution time for each tool invocation
- **Session tracking** - Group activities by Claude Code session
- **Error patterns** - Identify common failure points

All metrics are stored locally in daily JSONL files at `{{STARTUP_PATH}}/logs/` for privacy and performance.

### Analyze Your Metrics

View comprehensive analytics about your AI team's performance:

```bash
# View today's summary dashboard
the-startup log summary

# See metrics for the last 24 hours
the-startup log summary --since 24h

# Analyze specific tools
the-startup log tools --tool Edit --tool Write

# Review error patterns
the-startup log errors --since 7d

# Show activity timeline
the-startup log timeline --since today
```

### Export Formats

Export metrics for further analysis or reporting:

```bash
# Export as JSON for processing
the-startup log summary --format json --output metrics.json

# Generate CSV for spreadsheets
the-startup log tools --format csv --output tool-usage.csv

# Filter by session
the-startup log summary --session "dev-abc123"
```

### Key Features

- **Zero performance impact** - Hooks run asynchronously and fail silently
- **Privacy-first** - All data stays local, no external services
- **Automatic correlation** - Links PreToolUse and PostToolUse events
- **Daily rotation** - Files organized by date for easy management
- **Rich filtering** - Query by time, tool, session, or success status

### How It Works

1. **Automatic Setup** - Installation configures PreToolUse and PostToolUse hooks
2. **Silent Collection** - Hooks process in background without disrupting work
3. **Local Storage** - Metrics saved to `logs/YYYYMMDD.jsonl` files
4. **On-Demand Analysis** - Run analysis commands whenever you need insights

This metrics system helps you understand:
- Which tools slow down your workflow
- Where errors commonly occur
- How different agents perform
- Your most productive patterns

## Installation

The Agentic Startup provides easy installation via script capabilities.

```bash
# Interactive installation (shows all options)
curl -LsSf https://raw.githubusercontent.com/rsmdt/the-startup/main/install.sh | sh

# Quick global installation (recommended paths, no prompts)
curl -LsSf https://raw.githubusercontent.com/rsmdt/the-startup/main/install.sh | sh -s -- -y

# Local installation (project-specific paths, with file selection)
curl -LsSf https://raw.githubusercontent.com/rsmdt/the-startup/main/install.sh | sh -s -- -l

# Quick local installation (project-specific, no prompts)
curl -LsSf https://raw.githubusercontent.com/rsmdt/the-startup/main/install.sh | sh -s -- -ly
```

### Building from Source

If you want to contribute or customize:

```bash
# Clone and build
git clone https://github.com/rsmdt/the-startup.git
cd the-startup
go build -o the-startup

# Run tests
go test ./...

# Install from local binary (for development/offline use)
./the-startup install              # Interactive
./the-startup install -y           # Quick global
./the-startup install -ly          # Quick local
```

## Disclaimer

While The Agentic Startup aims to enhance Claude Code with specialized agents and structured workflows, be aware of some limitations:

### Command & Documentation Behavior
- **Slash commands** sometimes are not recognized or executed properly despite correct setup, requiring retry
- **Subagents** sometimes do not follow their custom instructions and instead generate generic prompts, breaking intended behavior
- Your **CLAUDE.md** may affect slash command or subagent behaviour
- Your **Installed MCPs** may affect behaviour and implementation.

### Best Practices
- **Restart** Claude Code between major tasks to free resources
- **Verify** generated code and documentation before committing

### Note about available MCP

The Agentic Startup tries to be unbiased about which MCP you may have installed, as this is a fast changing topic. However, we recommend that you have at least [`sequentialthunking`](https://github.com/modelcontextprotocol/servers/blob/main/src/sequentialthinking/README.md) installed.

## Learn More

- [Claude Code Documentation](https://docs.anthropic.com/en/docs/claude-code)
- [Claude Code Slash Commands](https://docs.anthropic.com/en/docs/claude-code/slash-commands)
- [Claude Code Subagents](https://docs.anthropic.com/en/docs/claude-code/sub-agents)
- [Claude Code Statusline](https://docs.anthropic.com/en/docs/claude-code/statusline)
- [Claude Code Output Styles](https://docs.anthropic.com/en/docs/claude-code/output-styles)

---

**Ship faster. Ship better. Ship with The Agentic Startup.**
