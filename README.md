# The Agentic Startup

**Comprehensive agentic software development framework for Claude Code with workflow orchestration, autonomous skills, and structured specification management.**

<p align="center">
  <strong>Ship faster. Ship better. Ship with The Agentic Startup.</strong>
</p>

---

## 🎉 Version 2.0 - Claude Code Marketplace Edition

**The Agentic Startup has been completely rewritten for the Claude Code marketplace!**

### What's New in 2.0

✨ **Native Claude Code Integration**
- Distributed as official Claude Code marketplace plugins
- Seamless installation via Claude Code plugin system
- Zero manual configuration required

🤖 **Autonomous Skills System**
- Model-invoked skills that activate based on natural language
- Progressive disclosure for optimal token efficiency
- Skills for documentation, agent delegation, and more

🎯 **Streamlined Architecture**
- Commands orchestrate high-level workflows
- Skills provide autonomous capabilities
- Rules define operational patterns

### Migrating from 1.x

If you previously installed via npm (`the-agentic-startup`), the 1.x CLI is deprecated:

**Before (1.x - Deprecated):**
```bash
npm install -g the-agentic-startup
the-agentic-startup install
```

**Now (2.0 - Current):**
```bash
# In Claude Code
/plugin install @the-startup/start
```

No need to uninstall 1.x - the marketplace version is completely separate.

---

## 📦 Quick Install

### Requirements

- **Claude Code v2.0+** - Required for marketplace plugins
- Git repository (recommended for full workflow)

### Installation

Install the Start plugin from the Claude Code marketplace:

```bash
/plugin install @the-startup/start
```

That's it! You now have access to 6 workflow commands and 2 autonomous skills.

### Initial Setup

After installation, configure your environment:

```bash
/start:init
```

This command will:
- ✅ Set up **The Startup** output style (high-energy orchestration)
- ✅ Configure your statusline (git branch integration)
- ✅ Ask for your preferences interactively

**Recommended:** Always run `/start:init` in new projects to configure Claude Code for optimal workflow.

---

## 🚀 Quick Start

### Your First Specification

Create a comprehensive specification from a brief description:

```bash
/start:specify Add user authentication with OAuth support
```

Claude will orchestrate specialist agents to create:
- **PRD** (Product Requirements) - What to build and why
- **SDD** (Solution Design) - How to build it technically
- **PLAN** (Implementation Plan) - Executable tasks and phases

### Execute the Implementation

Once your specification is ready:

```bash
/start:implement 001
```

Claude executes the plan phase-by-phase with:
- Parallel agent coordination
- Continuous validation
- Real-time progress tracking
- User confirmation at phase boundaries

### Analyze Your Codebase

Discover and document patterns in your code:

```bash
/start:analyze security
```

Claude will:
- Launch parallel specialist agents
- Discover reusable patterns
- Document in `docs/patterns/`, `docs/domain/`, `docs/interfaces/`
- Ensure no duplication

---

## 🎯 Claude Code Features Explained

The Agentic Startup leverages Claude Code's powerful extensibility features:

### 🔌 Plugins (Marketplace)

**What:** Distributable packages of commands, skills, agents, and rules

**How we use it:**
- `@the-startup/start` - Workflow orchestration plugin
- Future: `@the-startup/agents` - Specialized agent library

**Install:** `/plugin install @the-startup/start`

### ⚡ Commands (User-Invoked)

**What:** Slash commands you explicitly run (e.g., `/start:specify`)

**How we use it:**
- 6 workflow commands for specification, implementation, analysis, refactoring
- Commands orchestrate multi-step processes
- User decides when to invoke

**Example:** `/start:specify Add real-time notifications`

### 🤖 Skills (Model-Invoked)

**What:** Autonomous capabilities Claude activates based on context

**How we use it:**
- `documentation` - Automatically documents patterns/interfaces when discovered
- `agent-delegation` - Breaks down tasks and coordinates agents

**Activation:** Natural language (e.g., "break down this complex task")

### 👥 Agents (Coming Soon)

**What:** Specialized personas with focused expertise

**Status:** Framework designed, library in development

**Future:** 50+ specialist agents across 9 professional roles

### 📊 Statusline (Hooks)

**What:** Dynamic status bar showing context at bottom of Claude Code

**How we use it:**
- Git branch integration
- Current command state
- Configured via `/start:init`

**Example:** `[main] | /start:specify running...`

### 🎨 Output Style

**What:** Personality and communication style for Claude

**How we use it:**
- **The Startup** - High-energy, parallel-execution orchestration style
- Automatically included with plugin
- Activated via `/start:init`

**Style:** Y Combinator energy meets operational excellence

---

## 📋 Start Plugin Reference

The `@the-startup/start` plugin provides workflow orchestration for agentic development.

### Commands

#### `/start:specify <description>`

Create comprehensive specifications from brief descriptions.

**Purpose:** Generate PRD, SDD, and PLAN documents

**Example:**
```bash
/start:specify Build a real-time notification system with WebSocket support
```

**What you get:**
- `docs/specs/001-notification-system/PRD.md` - Product requirements
- `docs/specs/001-notification-system/SDD.md` - Solution design
- `docs/specs/001-notification-system/PLAN.md` - Implementation plan

**Process:**
1. Creates spec directory with auto-incrementing ID
2. Orchestrates specialist agents for research
3. Documents requirements, design, and implementation steps
4. Validates completeness at each stage

---

#### `/start:implement <spec-id>`

Execute implementation plans phase-by-phase.

**Purpose:** Turn specifications into working code

**Example:**
```bash
/start:implement 001
```

**Process:**
1. Loads PLAN.md from spec directory
2. Executes tasks phase-by-phase
3. Launches parallel agents when safe
4. Validates at phase boundaries
5. Tracks progress with TodoWrite
6. Waits for user confirmation between phases

**Features:**
- Parallel execution within phases
- Sequential execution between phases
- Specification compliance validation
- Rollback on failures

---

#### `/start:analyze <area>`

Discover and document business rules, technical patterns, and system interfaces.

**Purpose:** Extract knowledge from existing codebase

**Example:**
```bash
/start:analyze security patterns in authentication
```

**Analysis areas:**
- `business` - Business rules, domain logic, workflows
- `technical` - Architectural patterns, code structure
- `security` - Security patterns, vulnerabilities
- `performance` - Optimization patterns, bottlenecks
- `integration` - API contracts, service integrations
- `data` - Storage patterns, data modeling
- `testing` - Test strategies, validation approaches
- `deployment` - CI/CD, infrastructure patterns

**Output:**
- `docs/domain/` - Business rules and domain knowledge
- `docs/patterns/` - Technical patterns and solutions
- `docs/interfaces/` - External service contracts

---

#### `/start:refactor <description>`

Improve code quality while strictly preserving behavior.

**Purpose:** Safe, incremental refactoring

**Example:**
```bash
/start:refactor Simplify the authentication middleware
```

**Guarantees:**
- All tests pass before and after
- Behavior is strictly preserved
- Incremental changes (never big-bang)
- Rollback on test failures

**Process:**
1. Establishes test baseline
2. Analyzes code for improvements
3. Applies changes incrementally
4. Validates tests after each change
5. Documents refactoring patterns discovered

---

#### `/start:spec <name> [--add <template>]`

Create numbered spec directories with auto-incrementing IDs.

**Purpose:** Manage specification directories

**Examples:**
```bash
# Create new spec directory
/start:spec user-authentication

# Create spec and add PRD template
/start:spec payment-integration --add product-requirements

# Read existing spec
/start:spec 001 --read

# Add template to existing spec
/start:spec 001 --add solution-design
```

**Templates available:**
- `product-requirements` - PRD template
- `solution-design` - SDD template
- `implementation-plan` - PLAN template

**Output format:** `docs/specs/NNN-feature-name/`

---

#### `/start:init`

Initialize The Agentic Startup framework in your Claude Code environment.

**Purpose:** One-time setup for optimal configuration

**Example:**
```bash
/start:init
```

**What it does:**
1. **Configures Output Style**
   - Activates "The Startup" style
   - High-energy, execution-focused communication
   - Parallel agent orchestration mindset

2. **Sets Up Statusline**
   - Adds git branch to statusline
   - Tracks current command state
   - Custom hooks for real-time updates

3. **Interactive Setup**
   - Asks for your preferences
   - Confirms each configuration
   - Shows what changed

**When to run:**
- First time using the plugin in a project
- After updating to a new version
- When starting a new repository
- If you want to reconfigure settings

**Requirements:** Claude Code v2.0+

---

### Skills (Autonomous)

#### `documentation`

**Activates when:** Patterns, interfaces, or domain rules are discovered

**Trigger terms:** "pattern", "interface", "domain rule", "document", "reusable"

**What it does:**
- Checks for existing documentation (prevents duplicates)
- Categorizes correctly (domain/patterns/interfaces)
- Uses appropriate templates
- Creates cross-references
- Reports what was documented

**Example activation:**
```
Agent discovers: "I found a reusable caching pattern using Redis"
↓
Documentation skill activates automatically
↓
Creates: docs/patterns/caching-strategy.md
```

**Progressive disclosure:**
- `SKILL.md` - Core documentation logic (~7 KB)
- `reference.md` - Advanced protocols (~11 KB, loads when needed)
- `templates/` - Pattern, interface, domain templates (~6 KB each)

---

#### `agent-delegation`

**Activates when:** Task decomposition, agent coordination, or template generation needed

**Trigger terms:** "break down", "launch agents", "FOCUS/EXCLUDE", "parallel", "coordinate"

**What it does:**
- Decomposes complex tasks into activities
- Determines parallel vs sequential execution
- Generates FOCUS/EXCLUDE templates for agents
- Coordinates file creation (prevents collisions)
- Validates agent responses for scope compliance
- Generates retry strategies for failed agents

**Example activation:**
```
User: "Break down this authentication task"
↓
Agent-delegation skill activates
↓
Outputs:
- Activity breakdown
- Dependency analysis
- Parallel/sequential recommendation
- FOCUS/EXCLUDE templates for each activity
```

**Progressive disclosure:**
- `SKILL.md` - Core delegation logic (~24 KB)
- `reference.md` - Advanced patterns (~19 KB, loads when needed)
- `examples/` - Real-world scenarios (~38 KB, loads when relevant)

---

### Rules (Operational Workflows)

#### `cycle-pattern.md`

**What:** Discovery → Documentation → Review workflow pattern

**Used by:** All iterative commands (specify, analyze)

**Process:**
1. **Discovery Phase** - Launch parallel specialist agents to research
2. **Documentation Phase** - Document findings and update main document
3. **Review Phase** - Present findings to user, get confirmation
4. **Repeat** - Until work is complete

**Purpose:** Ensures consistent iterative workflow across commands

---

### Templates

Rich templates for structured documentation:

```
plugins/start/templates/
├── product-requirements.md      # PRD structure
├── solution-design.md            # SDD structure
├── implementation-plan.md        # PLAN structure
├── definition-of-ready.md        # Quality gate
├── definition-of-done.md         # Quality gate
└── task-definition-of-done.md   # Task-level quality gate
```

**Usage:** Automatically used by `/start:spec --add <template>`

---

### Hooks

#### SessionStart Hook

**When:** Every new Claude Code session

**What it does:**
- Displays welcome banner (first session only)
- Shows available commands
- Confirms plugin is active

#### StatuslineComplete Hook

**When:** After statusline updates

**What it does:**
- Adds git branch information
- Shows current command state
- Updates dynamically during execution

**Configure via:** `/start:init`

---

## 🏗️ Documentation Structure

The plugin encourages structured knowledge management:

```
docs/
├── specs/
│   └── [3-digit-number]-[feature-name]/
│       ├── PRD.md                          # What to build
│       ├── SDD.md                          # How to build it
│       └── PLAN.md                         # Implementation tasks
│
├── domain/                                  # Business rules
│   ├── user-permissions.md
│   ├── order-workflow.md
│   └── pricing-rules.md
│
├── patterns/                                # Technical patterns
│   ├── authentication-flow.md
│   ├── caching-strategy.md
│   └── error-handling.md
│
└── interfaces/                              # External integrations
    ├── stripe-payments.md
    ├── sendgrid-webhooks.md
    └── oauth-providers.md
```

### Auto-Documentation

The `documentation` skill automatically creates files in the correct location when patterns, interfaces, or domain rules are discovered during:
- Specification creation (`/start:specify`)
- Implementation (`/start:implement`)
- Analysis (`/start:analyze`)

### Deduplication

The skill always checks existing documentation before creating new files, preventing duplicates.

---

## 🎨 The Startup Output Style

Included with the plugin, activated via `/start:init`.

### Personality

**The Startup** embodies:
- **The Visionary Leader** - "We'll figure it out" - execute fast, iterate faster
- **The Rally Captain** - Turn challenges into team victories
- **The Orchestrator** - Run parallel execution like a conductor
- **The Pragmatist** - MVP today beats perfect next quarter

### Communication Style

**How The Startup communicates:**
- High energy, high clarity ("Let's deliver this NOW!")
- Execution mentality ("We've got momentum, let's push!")
- Celebrate wins ("That's what I'm talking about!")
- Own failures fast ("That didn't work. Here's the fix.")
- Always forward motion ("Next, we're tackling...")

### Workflow Patterns

**What you get:**
- Parallel-first mindset (launches multiple agents simultaneously)
- TodoWrite obsession (tracks every task religiously)
- "Ask yourself" checkpoints (self-validation at key decision points)
- Investor update summaries (comprehensive status reports)

### When to Use

**Perfect for:**
- Fast-paced development
- Complex multi-step workflows
- Parallel agent coordination
- High-energy execution

**Maybe not for:**
- Simple single-step tasks
- Exploratory conversations
- Learning/tutorial sessions

---

## 🔄 Typical Development Workflow

Here's how to use The Agentic Startup end-to-end:

### 1. **Initial Setup** (Once per project)

```bash
/start:init
```

Configures output style and statusline.

### 2. **Create Specification**

```bash
/start:specify Add real-time notification system with WebSocket support and email fallback
```

**What happens:**
- Creates `docs/specs/001-notification-system/`
- Generates PRD (requirements and use cases)
- Generates SDD (technical architecture and design)
- Generates PLAN (implementation tasks and phases)
- Documents discovered patterns/interfaces

**Duration:** 15-30 minutes (depending on complexity)

### 3. **Review Specification**

Read generated files:
- `docs/specs/001-notification-system/PRD.md`
- `docs/specs/001-notification-system/SDD.md`
- `docs/specs/001-notification-system/PLAN.md`

Provide feedback if needed, Claude will revise.

### 4. **Execute Implementation**

```bash
/start:implement 001
```

**What happens:**
- Loads PLAN.md
- Executes Phase 1 tasks
- Waits for user confirmation
- Executes Phase 2 tasks
- Continues phase-by-phase until complete

**Duration:** Varies by complexity (hours to days)

### 5. **Analyze Patterns** (During or after implementation)

```bash
/start:analyze technical patterns in notification system
```

**What happens:**
- Discovers patterns used in implementation
- Documents in `docs/patterns/`
- Creates cross-references
- Prevents future duplication

### 6. **Refactor** (As needed)

```bash
/start:refactor Simplify the WebSocket connection manager
```

**What happens:**
- Establishes test baseline
- Analyzes code for improvements
- Applies incremental refactorings
- Validates tests after each change

### 7. **Document Learnings**

Patterns and interfaces are automatically documented throughout the process by the `documentation` skill.

---

## 🤖 Autonomous Skills in Action

### Example 1: Documentation Skill

**Scenario:** During implementation, an agent discovers a pattern

```
Agent output: "I implemented a retry mechanism with exponential backoff for API calls"
```

**What happens automatically:**
1. Documentation skill recognizes "pattern" trigger
2. Checks `docs/patterns/` for existing retry patterns
3. Not found → Creates `docs/patterns/api-retry-strategy.md`
4. Uses pattern template
5. Reports: "📝 Created docs/patterns/api-retry-strategy.md"

**You didn't have to:** Manually request documentation or specify the path

---

### Example 2: Agent-Delegation Skill

**Scenario:** Complex task needs breakdown

```
User: "Implement user authentication - break this down into activities"
```

**What happens automatically:**
1. Agent-delegation skill recognizes "break this down"
2. Analyzes task complexity
3. Generates output:

```
Task: Implement user authentication

Activities:
1. Analyze security requirements
2. Design database schema
3. Create API endpoints
4. Build login UI

Dependencies: 1 → 2 → (3 & 4 parallel)

Execution: Sequential (1→2), then Parallel (3&4)

Agent Prompts Generated: ✅
```

**You didn't have to:** Manually create FOCUS/EXCLUDE templates or plan execution strategy

---

## 📦 Plugin Architecture

### Directory Structure

```
plugins/start/
├── .claude-plugin/
│   └── plugin.json              # Plugin manifest
│
├── commands/                     # Slash commands (user-invoked)
│   ├── analyze.md
│   ├── implement.md
│   ├── init.md
│   ├── refactor.md
│   ├── spec.md
│   └── specify.md
│
├── skills/                       # Skills (model-invoked)
│   ├── documentation/
│   │   ├── SKILL.md
│   │   ├── reference.md
│   │   └── templates/
│   └── agent-delegation/
│       ├── SKILL.md
│       ├── reference.md
│       └── examples/
│
├── rules/                        # Operational workflows
│   └── cycle-pattern.md
│
├── output-styles/                # Communication styles
│   └── the-startup.md
│
├── hooks/                        # Lifecycle hooks
│   ├── session-start.sh
│   └── statusline-complete.sh
│
├── scripts/                      # Utility scripts
│   └── spec.py
│
└── templates/                    # Document templates
    ├── product-requirements.md
    ├── solution-design.md
    ├── implementation-plan.md
    ├── definition-of-ready.md
    ├── definition-of-done.md
    └── task-definition-of-done.md
```

### How Components Work Together

**Commands** orchestrate workflows:
- Launch specialist agents
- Use trigger language to activate skills
- Reference rules for process patterns

**Skills** provide autonomous capabilities:
- Activate based on natural language
- No explicit invocation needed
- Progressive disclosure (load details only when needed)

**Rules** define operational patterns:
- Process workflows (e.g., cycle-pattern)
- Referenced by commands
- Lightweight (just principles)

---

## 🔐 Security & Privacy

### Security Approach

The plugin assists with **defensive security tasks only**:
- ✅ Security analysis and assessment
- ✅ Vulnerability identification
- ✅ Security pattern documentation
- ✅ Incident response planning

**It will refuse:**
- ❌ Creating malicious code
- ❌ Exploiting vulnerabilities
- ❌ Bypassing security controls
- ❌ Any offensive security tasks

### Privacy

**No data collection:**
- Plugin runs entirely locally in Claude Code
- No telemetry or analytics
- No external API calls
- Your code never leaves your machine

---

## 🚧 Roadmap

### Coming Soon

**Agents Plugin** (`@the-startup/agents`)
- 50+ specialized agents across 9 professional roles
- Activity-focused (not role-focused)
- Framework-agnostic (React, Vue, Angular, etc.)

**Additional Skills**
- `specification-review` - Validate implementation against specs
- `quality-gates` - Execute DOR/DOD validations
- `iterative-cycles` - Orchestrate discovery-documentation-review loops

**Enhanced Commands**
- `/start:test` - Generate comprehensive test suites
- `/start:deploy` - Deployment orchestration
- `/start:monitor` - Setup monitoring and observability

---

## 📖 Version History

### 2.0.0 (Current - Claude Code Marketplace)
- ✨ Complete rewrite for Claude Code marketplace
- 🤖 Autonomous skills system (documentation, agent-delegation)
- ⚡ 6 workflow commands (specify, implement, analyze, refactor, spec, init)
- 🎨 The Startup output style included
- 📊 Statusline integration with git branch
- 🔌 Native plugin architecture
- 📝 Progressive disclosure for optimal token usage
- 🎯 DRY architecture (82% rules reduction)

### 1.0.0 (Deprecated - npm CLI)
- Initial release as npm package
- Interactive TUI installation
- Manual component installation
- **No longer maintained**

---

## 🤝 Contributing

Contributions welcome! Here's how:

### Report Issues

[GitHub Issues](https://github.com/rsmdt/the-startup/issues)

### Contribute Code

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

### Extend the Framework

**Create new commands:**
- Add to `plugins/start/commands/`
- Follow existing command structure
- Use trigger language for skills

**Create new skills:**
- Add to `plugins/start/skills/`
- Include SKILL.md with proper frontmatter
- Use progressive disclosure (reference.md, examples/)

**Create new templates:**
- Add to `plugins/start/templates/`
- Follow markdown format
- Include placeholder sections

---

## 📚 Further Reading

### Documentation

- **[Skills Pattern Documentation](docs/patterns/claude-code-skills-integration.md)** - How skills work
- **[Agent Delegation Analysis](docs/patterns/agent-delegation-skill-extraction.md)** - Delegation architecture

### External Resources

- **[Claude Code Documentation](https://docs.claude.com/claude-code)** - Official Claude Code docs
- **[Claude Code Skills Guide](https://docs.claude.com/claude-code/skills)** - How to create skills

---

## 📜 License

MIT License - see [LICENSE](LICENSE) file for details

---

## 🎯 Get Started Now

### Install

```bash
/plugin install @the-startup/start
```

### Configure

```bash
/start:init
```

### Build Something

```bash
/start:specify Build a real-time chat application with WebSocket support
```

---

<p align="center">
  <strong>Ready to 10x your development workflow?</strong><br>
  Let's ship something incredible! 🚀
</p>
