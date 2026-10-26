# Specification Management Reference

## Spec ID Format

- **Format**: 3-digit zero-padded number (001, 002, ..., 999)
- **Auto-incrementing**: Script scans existing directories to find next ID
- **Directory naming**: `[NNN]-[sanitized-feature-name]`
- **Sanitization**: Lowercase, special chars → hyphens, trim leading/trailing hyphens

## Directory Structure

```
.start/specs/
├── 001-user-authentication/
│   ├── spec.md                   # Requirements, Solution, Plan parts + Ledger (gates in frontmatter)
│   ├── review.md                 # Solution review report — written by specify-solution's review workflow
│   └── .tracker/                 # Implementation tasks — created by specify-plan via task.js
│       ├── <phase-task>.md       #   one root task per plan phase
│       ├── <slice-task>.md       #   slice tasks under their phase (code brief in body)
│       └── <scenario-task>.md    #   scenario child tasks (holdout evidence)
├── 002-payment-processing/
│   └── ...
└── 003-notification-system/
    └── ...
```

`/implement` runs a dynamic workflow that builds each slice task and writes its status in the tracker
as slices are satisfied (the phase root rolls up).


## Script Commands

### Create New Spec
```bash
node <plugin>/scripts/spec.js "feature name here"
```
**Output:**
```
Created spec directory: .start/specs/005-feature-name-here
Created spec document: .start/specs/005-feature-name-here/spec.md
Created tracker directory: .start/specs/005-feature-name-here/.tracker
Spec ID: 005
Specification directory created successfully
```

Creates the spec directory, `spec.md` from the spec template, and `.tracker/`. Each phase fills its own part of `spec.md`.

### Read Spec Metadata
```bash
node <plugin>/scripts/spec.js 005 --read
```
**Output (TOML):**
```toml
id = "005"
name = "feature-name-here"
dir = ".start/specs/005-feature-name-here"

[spec]
spec = ".start/specs/005-feature-name-here/spec.md"
prd = ".start/specs/005-feature-name-here/spec.md"
sdd = ".start/specs/005-feature-name-here/spec.md"
plan_doc = ".start/specs/005-feature-name-here/spec.md"

files = [
  "spec.md"
]
```

`prd`, `sdd`, and `plan_doc` all point at `spec.md`; a legacy spec without `spec.md` reports its separate `requirements.md`, `solution.md`, and `plan.md` instead.

Implementation tasks live in `.tracker/` — query them with `task.js graph --spec 005`, not via `spec.js`.

## Template Resolution

The spec template is `skills/specify-meta/template.md`.

## Gates

The `gates` frontmatter of `spec.md` is the phase state:

| Value | Meaning |
|-------|---------|
| `pending` | Not started |
| `draft` | Being written |
| `approved` | The user approved the part |
| `skipped` | The user skipped the part (decision logged) |

## Phase Workflow

```
Initialization
    ↓
Requirements (specify-requirements → ## Requirements)   [WHY]
    ↓
Solution (specify-solution → ## Solution)               [WHAT — design]
    ↓
Plan (specify-plan → ## Plan + tracker tasks)           [WHAT — the plan]
    ↓
Ready for Implementation   (/implement runs the slice loop)
```

## Decision Logging

Record decisions with:
- **Date**: When the decision was made
- **Decision**: What was decided
- **Rationale**: Why (external references like JIRA IDs welcome)

Example:
```markdown
| Date | Decision | Rationale |
|------|----------|-----------|
| 2025-12-10 | Requirements skipped | Requirements in JIRA-1234 |
| 2025-12-10 | Start with Solution | Technical spike completed |
```
