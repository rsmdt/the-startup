# CLAUDE.md

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
│   │   ├── agents/               # 11 roles with activity-based specializations
│   │   │   ├── the-chief.md      # Complexity assessment, routing
│   │   │   ├── the-analyst/      # research-market, research-requirements
│   │   │   ├── the-architect/    # design-system, review-security, etc.
│   │   │   ├── the-developer/    # build-feature, optimize-performance, etc.
│   │   │   ├── the-devops/       # build-containers, build-pipelines, etc.
│   │   │   ├── the-designer/     # design-visual, design-interaction, etc.
│   │   │   ├── the-tester/       # test-quality, test-performance
│   │   │   └── the-meta-agent.md # Agent design and generation
│   │   └── skills/               # Cross-cutting and domain skills
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
│   └── specs/                    # Feature specifications (PRD/SDD/PLAN)
│
├── install.sh                    # One-line installer script
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

For full skill conventions, PICS structure, SudoLang syntax rules, and transformation checklist, see the `writing-skills` skill (`plugins/start/skills/writing-skills/reference/conventions.md`).

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

# Or use the main installer to test full installation
./install.sh

# Uninstall to reset
claude plugin uninstall start@the-startup
claude plugin uninstall team@the-startup
```

### Editing Skills

1. Skills are markdown files (`SKILL.md`) - edit directly
2. Skills activate on trigger terms defined in their YAML frontmatter
3. Keep `SKILL.md` under ~25 KB for context efficiency
4. Move advanced content to `reference.md` (loaded on demand)

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
- `reference.md` - Advanced protocols when complexity increases
- `templates/` - Document templates when creating artifacts
- `examples/` - Real-world scenarios when pattern matching

### Spec-Driven Development

The primary workflow: `/start:specify` → `/start:validate` → `/start:implement` → `/start:review`

Specifications live in `docs/specs/[NNN]-[name]/`:
- `product-requirements.md` - What to build
- `solution-design.md` - How to build it
- `implementation-plan.md` - Executable tasks

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
| Specs | `docs/specs/[NNN]-*/` | 3-digit ID prefix |

## Publishing

The repository is a Claude Code marketplace. Publishing happens via:
1. Push to `main` branch
2. GitHub Actions workflow creates release
3. Users install via `/plugin marketplace add rsmdt/the-startup`
