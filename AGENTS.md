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
│   │   ├── skills/               # workflow skills (user-invocable) + internal spec/implement skills
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

The specification describes **WHY** and **WHAT** (the target state and what "done" looks like); the implementation is a set of **tasks** in the spec's tracker. **HOW** to build each task is decided by the code agent inside the implementation loop. The two concerns are split across a clean line.

Specifications live in `.start/specs/[NNN]-[name]/` (legacy: `docs/specs/`):
- `requirements.mdx` — **WHY**: problem, users, acceptance criteria
- `solution.mdx` — **WHAT (design)**: target architecture, components, interfaces
- `plan.mdx` — **WHAT (narrative)**: the high-level approach — phases and achievable outcomes, plus a link to the tracker. Read for orientation, not parsed for work.
- `.tracker/` — the operational **tasks** (managed by `task.js`): a root **phase** task, one **slice** task per vertical deliverable (its body is the code brief — features + EARS behavior + quality gates), and under each slice its **scenario** child tasks (the holdout acceptance evidence). `--depends-on` edges between slices are the parallelization contract.

The specification carries **no task-by-task instructions** — a slice says what it delivers and how it is proven, never the steps to build it. `/start:specify` runs Requirements → Solution → Plan (a 5-step flow: init → requirements → solution → plan → finalize); `specify-plan` writes the narrative `plan.mdx` and creates the tracker tasks.

#### One implementation loop (the dynamic-workflow runtime)

There is a single `/start:implement`, and it always runs on Claude Code's **dynamic-workflow runtime**. It reads the spec's tasks from the tracker and launches the plugin's reusable workflow script — `plugins/start/workflows/implement.js`, whose `meta.name` is `run-implement` — via `Workflow({ scriptPath, args })`, with the slice tasks as input. There is no install step: the script runs where it ships.

The workflow drives each slice through a **barriered code→eval retry loop**:
- a **code agent** builds the slice via TDD — it reads the slice task (features + EARS + gates), never its scenario child tasks;
- an **eval agent** proves the slice against its scenario child tasks through external interfaces — it reads the scenarios + service config, never the slice or the source;
- failures are filtered to one-line summaries and fed back until the slice's scenarios pass the `threshold`, or `max_iterations` is reached; each failed scenario becomes a `type: bug` child task.

The information barrier is *which task each agent reads* — the code agent the slice, the eval agent its scenario children. Independent slices (no unmet dependency) run concurrently via the runtime's `parallel()`; dependent slices run in order. Scale is emergent — a 1-slice plan is a short loop, a 40-slice plan fans out. Specs created before this model (legacy `### BEAN-N` plans, `plan/` / `manifest.md`) are not run; `/start:implement` offers to re-run `/start:specify` to regenerate the tasks. Users invoke `/start:specify` and `/start:implement`; everything else is internal.

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
| Workflow scripts | `plugins/*/workflows/*.js` | `meta.name` is **`run-<script>`** (see below) |
| Agents | `plugins/team/agents/the-*/` | `the-[role]/[activity].md` |
| Output Styles | `plugins/*/output-styles/*.md` | Title Case (e.g., `The Startup.md`) |
| Specs | `.start/specs/[NNN]-*/` | 3-digit ID prefix |

### Skill vs workflow naming

A workflow's `meta.name` is always **`run-<script>`** — the script's filename with a `run-` prefix.

| | Skill | Workflow (`meta.name`) |
|---|---|---|
| Shape | short noun or verb — the thing you type | `run-` + the script filename |
| Example | `implement` | `run-implement` |

| Skill | Script | Workflow |
|---|---|---|
| `implement` | `workflows/implement.js` | `run-implement` |
| `analyze` | `workflows/analyze.js` | `run-analyze` |
| `review` | `workflows/review.js` | `run-review` |
| `specify-solution` | `workflows/solution-review.js` | `run-solution-review` |
| `daydream` | `workflows/daydream.js` | `run-daydream` |

The prefix is the whole convention: it is mechanically derivable (no naming judgment, so no drift), it
sorts every workflow together, and `run-` is instantly recognizable in the `/workflows` progress view and
the run-approval dialog. The name only has to *identify* the run — the view shows `meta.description` and
the phase list beside it, so the name carries no burden to describe.

**Never give a workflow the same `meta.name` as a skill.** Skills and workflows share one `/[plugin]:`
namespace, so identical names collide and the workflow shadows the skill — replacing a working entry
point with a script that cannot run standalone. The `run-` prefix makes that impossible by construction:
no skill is named `run-*`.

### Where workflow scripts live

Scripts go in `plugins/*/workflows/` — the standard plugin location. Claude Code discovers them there
(or wherever the `workflows` manifest field points) and registers each as `/[plugin]:[meta.name]`, so
ours appear as `/start:run-implement`, `/start:run-analyze`, and so on.

`reference/` is for documentation a skill reads on demand — markdown, not executables. Workflow scripts
are orchestration code, and three of these skills needed no other reference files at all, so keeping the
script there meant a directory that existed only to hold it.

**These commands exist but are not entry points.** A workflow script gets no filesystem or shell access
and cannot take mid-run user input, so none of ours can resolve its own `args` — the skill builds them
(resolving the spec, seeding the tracker, getting the plan approved) and launches the script with
`Workflow({ scriptPath, args })`. Every script therefore opens with a **fail-fast guard**: launched
without usable args it throws immediately, naming the missing fields and the skill to run instead,
rather than templating `undefined` into every prompt and producing an empty report.

So: use `/start:implement`, not `/start:run-implement`. The `run-` prefix marks the second kind.

There is no install step and no version marker — the script runs where it ships, and the plugin version
governs.

## Publishing

The repository is a Claude Code marketplace. Publishing happens via:
1. Push to `main` branch
2. GitHub Actions workflow creates release
3. Users install via `/plugin marketplace add rsmdt/the-startup`
