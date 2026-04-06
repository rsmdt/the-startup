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
│   ├── manifest.md               # Created by specify-factory skill
│   ├── units/                    # Created by specify-factory skill
│   │   ├── dm1.md
│   │   ├── ve1.md
│   │   └── rl1.md
│   └── scenarios/                # Created by specify-factory skill
│       ├── dm1/
│       │   └── schema-validation.md
│       ├── ve1/
│       │   ├── sql-injection.md
│       │   └── empty-input.md
│       └── rl1/
│           └── burst-traffic.md
├── 002-payment-processing/
│   └── ...
└── 003-notification-system/
    └── ...
```

## Script Commands

### Create New Spec
```bash
spec.py "feature name here"
```
**Output:**
```
Created factory directories: units/, scenarios/
Created spec directory: .start/specs/005-feature-name-here
Spec ID: 005
Specification directory created successfully
```

Creates the spec directory with `units/` and `scenarios/` subdirectories.

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
manifest = ".start/specs/005-feature-name-here/manifest.md"
units = [".start/specs/005-feature-name-here/units/dm1.md", ...]
scenarios_dm1 = [".start/specs/005-feature-name-here/scenarios/dm1/schema-validation.md"]

files = [
  "README.md",
  "requirements.md",
  "solution.md",
  "manifest.md",
  "units/dm1.md",
  "units/ve1.md",
  "scenarios/dm1/schema-validation.md",
  "scenarios/ve1/sql-injection.md"
]
```

### Add Template to Existing Spec
```bash
spec.py 005 --add specify-requirements   # Creates requirements.md
spec.py 005 --add specify-solution       # Creates solution.md
spec.py 005 --add specify-factory        # Creates units/ and scenarios/ directories
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
Requirements (specify-requirements)
    ↓
Solution (specify-solution)
    ↓
Factory (specify-factory → units, scenarios, manifest)
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
| 2025-12-10 | Requirements skipped | Requirements in JIRA-1234 |
| 2025-12-10 | Start with Solution | Technical spike completed |
```
