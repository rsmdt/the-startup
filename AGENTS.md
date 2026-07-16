# AGENTS.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Is

**The Agentic Startup** is a spec-driven development framework for Claude Code, distributed as marketplace plugins. It provides workflow commands, autonomous skills, specialized agents, and output styles to transform how you build software.

## Repository Structure

```
the-startup/
├── plugins/
│   ├── start/                    # Core workflow orchestration plugin
│   │   ├── .claude-plugin/       # Plugin manifest (plugin.json)
│   │   ├── skills/               # 14 skills (9 user-invocable + 5 autonomous)
│   │   ├── output-styles/        # The Startup, The ScaleUp output styles
│   │   └── README.md             # Detailed plugin documentation
│   │
│   ├── team/                     # Specialized agent library plugin
│   │   ├── agents/               # 8 roles with 15 consolidated activity agents
│   │   │   ├── the-chief.md      # Complexity assessment, routing
│   │   │   ├── the-analyst/      # research-product
│   │   │   ├── the-architect/    # design-system, review-security, review-robustness, review-compatibility
│   │   │   ├── the-developer/    # build-feature, optimize-performance, etc.
│   │   │   ├── the-devops/       # build-platform, monitor-production
│   │   │   ├── the-designer/     # research-user, design-interaction, design-visual
│   │   │   ├── the-tester/       # test-strategy
│   │   │   └── the-meta-agent.md # Agent design and generation
│   │   └── skills/               # 16 consolidated cross-cutting and domain skills
│   │
│   └── constitution/             # Project governance rules plugin (optional)
│
├── scripts/
│   ├── statusline.sh             # Custom statusline for Claude Code terminal
│   └── statusline.toml           # Statusline configuration template
│
├── docs/
│   ├── PHILOSOPHY.md             # Activity-based architecture rationale
│   ├── PRINCIPLES.md             # Core development principles
│   ├── patterns/                 # Technical patterns documentation
│   └── specs/                    # Feature specifications (requirements, solution, plan + scenarios)
│
├── package.json                    # npm package definition for npx installer
└── README.md                     # User-facing documentation
```

## Key Concepts

### Plugin Architecture

Each plugin lives in `plugins/[name]/` with:
- `.claude-plugin/plugin.json` - Plugin manifest defining name, version, components
- `commands/` - Slash command definitions (markdown files)
- `skills/` - Autonomous skills (SKILL.md files with trigger terms)
- `output-styles/` - Output style definitions
- `agents/` - Agent definitions (team plugin only)

### Skill Structure

Skills use progressive disclosure to minimize context usage:
```
skills/[skill-name]/
├── SKILL.md           # Core logic (~7-24 KB, always loaded)
├── reference/         # Advanced protocols (loaded when needed)
├── templates/         # Document templates
├── examples/          # Real-world scenarios
└── validation.md      # Quality checklists
```

For full skill conventions, PICS structure, and transformation checklist, see the `writing-skills` skill (`plugins/start/skills/writing-skills/reference/conventions.md`).

### Agent Structure (Team Plugin)

Agents are organized by role with activity-based specializations:
```
agents/the-[role]/
├── [activity-1].md    # e.g., the-architect/design-system.md
├── [activity-2].md    # e.g., the-architect/review-security.md
└── ...
```

Each agent markdown file defines:
- Purpose and trigger conditions
- Tool access permissions
- Activation examples
- Integration with other agents

## Development Workflow

### Testing Changes Locally

```bash
# The plugins are directories - test by installing from local path
claude plugin install ./plugins/start
claude plugin install ./plugins/team

# Or use the interactive installer to test full installation
npx the-startup

# Uninstall to reset
claude plugin uninstall start@the-startup
claude plugin uninstall team@the-startup
```

### Editing Skills

1. Skills are markdown files (`SKILL.md`) - edit directly
2. Skills activate on trigger terms defined in their YAML frontmatter
3. Keep `SKILL.md` under ~25 KB for context efficiency
4. Move advanced content to `reference/` directory (loaded on demand)

### Editing Commands

1. Commands are markdown files in `plugins/[plugin]/commands/`
2. Command name matches filename (e.g., `specify.md` → `/start:specify`)
3. Commands orchestrate skills - they define workflow, skills provide implementation

### Editing Agents

1. Agents are markdown files in `plugins/team/agents/[role]/`
2. Agent name matches `[role]:[activity]` pattern
3. Agents define specialized Task tool prompts for subagent delegation

### Editing Output Styles

1. Output styles are markdown files in `plugins/[plugin]/output-styles/`
2. Define personality, voice, and behavioral patterns
3. Activated via `/output-style [plugin]:[style-name]`

## Important Patterns

### Progressive Disclosure

Skills load minimal context initially, then progressively load:
- `reference/` - Advanced protocols when complexity increases
- `templates/` - Document templates when creating artifacts
- `examples/` - Real-world scenarios when pattern matching

### Spec-Driven Development

The primary workflow: `/start:specify` → `/start:validate` → `/start:implement` → `/start:review`

The specification describes **WHY** and **WHAT** (the target state and what "done" looks like). **HOW** to execute it is decided later, by `/start:implement`. The two concerns are split across a clean line.

Specifications live in `.start/specs/[NNN]-[name]/` (legacy: `docs/specs/`):
- `requirements.mdx` — **WHY**: problem, users, acceptance criteria
- `solution.mdx` — **WHAT (design)**: target architecture, components, interfaces
- `plan.mdx` — **WHAT (phase-approach baseline)**: a sequence of vertical-slice phases, each correlating the deliverable **features** (acceptance criteria from requirements), the **behavior** in EARS (technical detail from solution), the **quality gates**, and the **validation scenarios** that prove it
- `scenarios/phase-N/{name}.md` — holdout validation scenarios (authored by `specify-plan`)

The specification carries **no decomposition strategy and no task-by-task instructions** — the plan's phases are a *baseline, not a prescription*. `/start:specify` runs Requirements → Solution → Plan (a 5-step flow: init → requirements → solution → plan → finalize).

#### Execution strategy (chosen at implementation time)

`/start:implement` reads the spec, runs a strategy selector (see `plugins/start/skills/implement/reference/strategy-selector.md`), and chooses **how to consume the plan's phases**:

| Strategy | When | Engine skill | How it treats the plan's phases |
|----------|------|--------------|----------------------------------|
| **Direct** | Fixes, refactors, single-phase work | `implement-direct` | **ignores** phase ordering — flattens to 1–3 units, builds straight to the gates + scenarios |
| **Iterative** | One or two phases, human-in-the-loop | `implement-incremental` | **follows** the phases — materializes them into `.impl/phases/` as TDD tasks |
| **Decompose** | Many phases, 3+ components, parallel work | `implement-factory` | **decomposes** the phases into finer parallel units in `.impl/units/` + `.impl/manifest.md` |

The chosen strategy and any derived decomposition are written under `.impl/` (resumable, inspectable, separate from the spec). The plan's holdout `scenarios/phase-N/` are the evaluation set all strategies validate against — never re-authored at implement time. `specify-incremental` / `specify-factory` remain as the **decomposition methodology** the engines reuse, and to read legacy `plan/` / `manifest.md` specs (back-compat). Users invoke `/start:specify` and `/start:implement`; everything else is internal.

### Knowledge Capture

Discovered patterns, interfaces, and domain rules are automatically documented in:
- `docs/patterns/` - Technical patterns
- `docs/interfaces/` - External integrations
- `docs/domain/` - Business rules

### Constitution Enforcement

Optional `CONSTITUTION.md` at project root defines checkable rules:
- L1 (Must) - Blocking with autofix
- L2 (Should) - Blocking, manual fix
- L3 (May) - Advisory only

## File Naming Conventions

| Type | Location | Naming |
|------|----------|--------|
| Commands | `plugins/*/commands/*.md` | lowercase-kebab (e.g., `specify.md`) |
| Skills | `plugins/*/skills/*/SKILL.md` | directory is skill name |
| Agents | `plugins/team/agents/the-*/` | `the-[role]/[activity].md` |
| Output Styles | `plugins/*/output-styles/*.md` | Title Case (e.g., `The Startup.md`) |
| Specs | `.start/specs/[NNN]-*/` | 3-digit ID prefix |

## Publishing

The repository is a Claude Code marketplace. Publishing happens via:
1. Push to `main` branch
2. GitHub Actions workflow creates release
3. Users install via `/plugin marketplace add rsmdt/the-startup`
