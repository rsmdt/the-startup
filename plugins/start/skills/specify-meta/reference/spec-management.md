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
│   ├── requirements.mdx           # WHY — created by specify-requirements skill
│   ├── solution.mdx               # WHAT (design) — created by specify-solution skill
│   ├── plan.mdx                   # WHAT (completion contract) — created by specify-plan skill
│   ├── scenarios/                 # Holdout completion criteria — created by specify-plan skill
│   │   ├── login/
│   │   │   ├── valid-credentials.md
│   │   │   └── locked-account.md
│   │   └── password-reset/
│   │       └── expired-token.md
│   └── .impl/                     # Implementation-time, written by /implement (resumable)
│       ├── strategy.md            # Chosen strategy: Direct | Iterative | Decompose
│       ├── phases/                # Iterative strategy: generated phase plan
│       └── units/ + manifest.md   # Decompose strategy: generated units + execution order
├── 002-payment-processing/
│   └── ...
└── 003-notification-system/
    └── ...
```

Legacy specs may instead carry an Incremental `plan/` directory or Factory `manifest.md` + `units/` + `scenarios/{unit-id}/` at the spec root; `/implement` still reads these.

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

Creates the spec directory only. Document artifacts are added per phase via `--add`.

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
prd = ".start/specs/005-feature-name-here/requirements.mdx"
sdd = ".start/specs/005-feature-name-here/solution.mdx"
plan_doc = ".start/specs/005-feature-name-here/plan.mdx"
scenarios_login = [".start/specs/005-feature-name-here/scenarios/login/valid-credentials.md"]
impl_strategy = ".start/specs/005-feature-name-here/.impl/strategy.md"   # present after /implement runs

files = [
  "README.md",
  "requirements.mdx",
  "solution.mdx",
  "plan.mdx",
  "scenarios/login/valid-credentials.md"
]
```

### Add Template to Existing Spec
```bash
spec.py 005 --add specify-requirements   # Creates requirements.md (rename to .mdx per discipline)
spec.py 005 --add specify-solution       # Creates solution.md
spec.py 005 --add specify-plan           # Creates plan.mdx + scenarios/ directory
# Legacy (back-compat): specify-incremental creates plan/, specify-factory creates units/ + scenarios/
```

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
Requirements (specify-requirements → requirements.mdx)   [WHY]
    ↓
Solution (specify-solution → solution.mdx)               [WHAT — design]
    ↓
Plan (specify-plan → plan.mdx + scenarios/)              [WHAT — completion contract]
    ↓
Ready for Implementation   (/implement selects Direct | Iterative | Decompose → .impl/)
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
