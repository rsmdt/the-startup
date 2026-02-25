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
│   ├── requirements.md           # Created by specify-requirements skill
│   ├── solution.md               # Created by specify-solution skill
│   └── plan/                     # Created by specify-plan skill
│       ├── README.md             # Plan manifest (phases, checklist, context)
│       ├── phase-1.md            # Per-phase tasks and TDD structure
│       ├── phase-2.md
│       └── phase-3.md
├── 002-payment-processing/
│   └── ...
└── 003-notification-system/
    └── ...
```

## Legacy Fallback

The script supports backward compatibility with `docs/specs/`:

- **Read**: Checks `.start/specs/` first, falls back to `docs/specs/`
- **ID scanning**: Scans both directories, takes max ID across both
- **File names**: Supports both new (`requirements.md`, `solution.md`) and legacy (`product-requirements.md`, `solution-design.md`)
- **Plan**: Supports both `plan/` directory and legacy `implementation-plan.md`

## Script Commands

### Create New Spec
```bash
spec.py "feature name here"
```
**Output:**
```
Created spec directory: .start/specs/005-feature-name-here
Spec ID: 005
Specification directory created successfully
```

Creates the spec directory with an empty `plan/` subdirectory.

### Read Spec Metadata
```bash
spec.py 005 --read
```
**Output (TOML):**
```toml
id = "005"
name = "feature-name-here"
dir = ".start/specs/005-feature-name-here"

[spec]
prd = ".start/specs/005-feature-name-here/requirements.md"
sdd = ".start/specs/005-feature-name-here/solution.md"
plan_dir = ".start/specs/005-feature-name-here/plan"
plan = ".start/specs/005-feature-name-here/plan/README.md"
phases = [".start/specs/005-feature-name-here/plan/phase-1.md"]

files = [
  "README.md",
  "requirements.md",
  "solution.md",
  "plan/README.md",
  "plan/phase-1.md"
]
```

### Add Template to Existing Spec
```bash
spec.py 005 --add specify-requirements   # Creates requirements.md
spec.py 005 --add specify-solution       # Creates solution.md
spec.py 005 --add specify-plan           # Creates plan/README.md
```

When `--add specify-plan` (or `--add implementation-plan`) is called, the script creates `plan/README.md` from the plan template instead of a monolithic file.

## Template Resolution

Templates are resolved in this order:
1. `skills/[template-name]/template.md` (primary)
2. `templates/[template-name].md` (deprecated fallback)

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
PRD (Product Requirements)
    ↓
SDD (Solution Design)
    ↓
PLAN (Implementation Plan)
    ↓
Ready for Implementation
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
| 2025-12-10 | PRD skipped | Requirements in JIRA-1234 |
| 2025-12-10 | Start with SDD | Technical spike completed |
```
