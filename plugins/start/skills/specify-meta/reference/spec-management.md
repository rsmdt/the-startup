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
│   ├── README.md                 # Managed by specify-meta skill
│   ├── requirements.md          # WHY — created by specify-requirements skill
│   ├── solution.md              # WHAT (design) — created by specify-solution skill
│   ├── plan.md                  # WHAT (narrative) — created by specify-plan skill
│   └── .tracker/                 # Implementation tasks — created by specify-plan via task.js
│       ├── <phase-task>.md       #   root phase task
│       ├── <slice-task>.md       #   slice tasks (code brief in body)
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
Spec ID: 005
Specification directory created successfully
```

Creates the spec directory only. Document artifacts are added per phase via `--add`.

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
prd = ".start/specs/005-feature-name-here/requirements.md"
sdd = ".start/specs/005-feature-name-here/solution.md"
plan_doc = ".start/specs/005-feature-name-here/plan.md"

files = [
  "README.md",
  "requirements.md",
  "solution.md",
  "plan.md"
]
```

Implementation tasks live in `.tracker/` — query them with `task.js graph --spec 005`, not via `spec.js`.

### Add Template to Existing Spec
```bash
node <plugin>/scripts/spec.js 005 --add specify-requirements   # Creates requirements.md
node <plugin>/scripts/spec.js 005 --add specify-solution       # Creates solution.md
node <plugin>/scripts/spec.js 005 --add specify-plan           # Creates plan.md + ensures .tracker/
```

## Template Resolution

Templates are resolved from `skills/[template-name]/template.md`.

## README.md Fields

| Field | Description |
|-------|-------------|
| Created | Date spec was created |
| Current Phase | Active workflow phase |
| Last Updated | Date of last status change |
| Document Status | pending, in_progress, completed, skipped |
| Notes | Additional context for each document |

## Phase Workflow

```
Initialization
    ↓
Requirements (specify-requirements → requirements.md)   [WHY]
    ↓
Solution (specify-solution → solution.md)               [WHAT — design]
    ↓
Plan (specify-plan → narrative plan.md + tracker tasks)  [WHAT — the plan]
    ↓
Ready for Implementation   (/implement runs the slice loop)
```

Each phase can be:
- **Completed**: Document finished and validated
- **Skipped**: User decided to skip (decision logged)
- **In Progress**: Currently being worked on
- **Pending**: Not yet started

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
