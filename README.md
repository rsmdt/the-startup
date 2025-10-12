# The Agentic Startup

**Comprehensive agentic software development framework for Claude Code with specialized roles, workflow automation, and structured specification management.**

<p align="center">
  Ship faster. Ship better. Ship with <b>The Agentic Startup</b>.
</p>

---

## 📦 Installation

Install directly as a Claude Code plugin:

```bash
/plugin install irudiperera/the-startup
```

That's it! You now have access to 50 specialized agents and 6 workflow commands.

## 🚀 Quick Start

After installation, try this:

```
# Create a comprehensive specification
/s:specify Add user authentication with OAuth support

# Execute the implementation plan
/s:implement 001

# Analyze your system for patterns
/s:analyze security
```

## 🎯 Features

### 50 Specialized Agents

The plugin provides access to specialist agents across 9 professional roles:

- **The Analyst** (3 agents) - Requirements analysis, feature prioritization, project coordination
- **The Architect** (5 agents) - System architecture, technology research, documentation, quality review, standards
- **The Designer** (4 agents) - Design systems, accessibility, interaction architecture, user research
- **The ML Engineer** (4 agents) - Feature operations, prompt optimization, ML operations, context management
- **The Mobile Engineer** (3 agents) - Data persistence, mobile operations, mobile development
- **The Platform Engineer** (7 agents) - Containerization, pipelines, monitoring, infrastructure, performance, deployment, data architecture
- **The QA Engineer** (3 agents) - Test execution, exploratory testing, performance testing
- **The Security Engineer** (3 agents) - Security assessment, implementation, incident response
- **The Software Engineer** (16 agents) - Performance optimization, components, API development, domain modeling, and more
- **Orchestration** (2 agents) - The Chief (routing), The Meta Agent (agent generation)

[📖 View complete agent list →](docs/AGENTS.md)

### 6 Workflow Commands

#### `/s:specify <description>`
Create comprehensive specifications (PRD, SDD, PLAN) from brief descriptions.

```
/s:specify Build a real-time notification system
```

**What you get:**
- Product Requirements Documentation (PRD) - What to build and why
- Solution Design Documentation (SDD) - How to build it technically
- Implementation Plan (PLAN) - Executable tasks and phases

#### `/s:implement <spec-id>`
Execute implementation plans phase-by-phase with validation.

```
/s:implement 001
```

**Features:**
- Phase-by-phase execution with approval gates
- Parallel agent execution when possible
- Continuous validation and testing
- Real-time progress tracking

#### `/s:refactor <description>`
Improve code quality while strictly preserving behavior.

```
/s:refactor Simplify the authentication middleware
```

**Guarantees:**
- All tests pass after every change
- Behavior is strictly preserved
- Incremental, safe refactorings

#### `/s:analyze <area>`
Discover and document business rules, technical patterns, and system interfaces.

```
/s:analyze security patterns in payment processing
```

**Analysis areas:** business, technical, security, performance, integration, data, testing, deployment

#### `/s:init`
Initialize quality gate templates (Definition of Ready, Definition of Done).

```
/s:init
```

#### `/s:spec <name> [--add <template>]`
Create numbered spec directories with auto-incrementing IDs.

```
/s:spec user-authentication
/s:spec payment-integration --add product-requirements
/s:spec 001 --read
```

### Git Status Integration

The plugin automatically adds git branch information to your Claude Code statusline via hooks.

### Welcome Banner

On first session, the plugin displays a welcome banner introducing available capabilities.

## 📚 Documentation Structure

The plugin encourages structured documentation:

```
docs/
├── specs/
│   └── [3-digit-number]-[feature-name]/
│       ├── product-requirements.md
│       ├── solution-design.md
│       └── implementation-plan.md
├── domain/          # Business rules and workflows
├── patterns/        # Technical patterns and solutions
└── interfaces/      # API contracts and integrations
```

## 🤝 Agent Delegation Rules

All commands use consistent agent delegation patterns:

- **FOCUS** - What the agent should concentrate on
- **EXCLUDE** - What to explicitly avoid
- **CYCLE PATTERN** - Discovery → Implementation → Review loops

These rules are automatically loaded via `@rules/agent-delegation.md` references in commands.

## 🎨 The Startup Output Style

The plugin includes "The Startup" output style which provides:
- **High-energy communication** - "Let's ship this NOW!" enthusiasm
- **Parallel execution mindset** - Multiple agents simultaneously
- **TodoWrite obsession** - Tracks every task religiously
- **Startup DNA** - Y Combinator energy meets operational excellence

**Included with plugin installation** from `output-styles/the-startup.md`

Activate after installation:
```
/output-style The Startup
```

The Startup style provides:
- High-energy, execution-focused communication
- Automatic parallel agent execution
- TodoWrite task tracking
- "Let's ship!" enthusiasm

## 🔄 Development Workflow

Typical workflow using The Agentic Startup:

1. **Specify** - Create comprehensive specification
   ```
   /s:specify Add real-time notification system
   ```

2. **Review** - Check generated PRD, SDD, PLAN

3. **Implement** - Execute phase-by-phase
   ```
   /s:implement 001
   ```

4. **Analyze** - Document discovered patterns
   ```
   /s:analyze technical
   ```

5. **Refactor** - Improve code quality
   ```
   /s:refactor Simplify notification handlers
   ```

## 📋 Templates

The plugin includes rich templates accessible at `templates/`:

- `product-requirements.md` - Product requirements documentation
- `solution-design.md` - Technical solution design
- `implementation-plan.md` - Implementation planning
- `definition-of-ready.md` - Quality gate template
- `definition-of-done.md` - Quality gate template
- `task-definition-of-done.md` - Task-level quality gate

## 🏗️ Plugin Architecture

The plugin follows Claude Code's official structure:

```
the-agentic-startup/
├── .claude-plugin/
│   └── plugin.json          # Plugin manifest
├── agents/                  # 50 agent definitions
├── commands/                # 6 slash commands
├── hooks/                   # SessionStart, UserPromptSubmit
├── scripts/                 # Spec generation script
├── templates/               # Document templates
└── rules/                   # Agent delegation rules
```

## 🔐 Security Note

This plugin assists with **defensive security tasks only**. It will refuse to create, modify, or improve code that may be used maliciously.

## 📖 Version History

### 2.0.0 (Current - Claude Code Plugin)
- Complete migration to Claude Code plugin architecture
- Added SessionStart welcome hook
- Added UserPromptSubmit statusline hook
- Streamlined installation via plugin system
- 50 specialized agents across 9 roles
- 6 slash commands for workflow automation

### 1.0.0 (npm CLI Package - Deprecated)
- Initial release as npm package with interactive TUI installation
- Manual component installation workflow

## 📜 License

MIT License

## 🤝 Contributing

Contributions welcome! Please submit issues or pull requests.

## 📚 Further Reading

- **[Complete Agent Reference](docs/AGENTS.md)** - All 50 agents and their capabilities
- **[Philosophy & Research](docs/PHILOSOPHY.md)** - Why activity-based agents work

---

**Ready to 10x your development workflow?**

```bash
/plugin install irudiperera/the-startup
```

Then try:
```
/s:specify Build a user authentication system
```

Let's ship something incredible! 🚀
