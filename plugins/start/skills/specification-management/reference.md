# Specification Management Reference

## Spec ID Format

- **Format**: 3-digit zero-padded number (001, 002, ..., 999)
- **Auto-incrementing**: Script scans existing directories to find next ID
- **Directory naming**: `[NNN]-[sanitized-feature-name]`
- **Sanitization**: Lowercase, special chars → hyphens, trim leading/trailing hyphens

## Directory Structure

```
docs/specs/
├── 001-user-authentication/
│   ├── README.md                 # Managed by specification-management skill
│   ├── product-requirements.md   # Created by requirements-gathering-analysis skill
│   ├── solution-design.md        # Created by technical-architecture-design skill
│   └── implementation-plan.md    # Created by phased-implementation-planning skill
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
Created spec directory: docs/specs/005-feature-name-here
Spec ID: 005
Specification directory created successfully
```

### Read Spec Metadata
```bash
spec.py 005 --read
```
**Output (TOML):**
```toml
id = "005"
name = "feature-name-here"
dir = "docs/specs/005-feature-name-here"

[spec]
prd = "docs/specs/005-feature-name-here/product-requirements.md"
sdd = "docs/specs/005-feature-name-here/solution-design.md"

files = [
  "product-requirements.md",
  "README.md",
  "solution-design.md"
]
```

### Add Template to Existing Spec
```bash
spec.py 005 --add product-requirements
spec.py 005 --add solution-design
spec.py 005 --add implementation-plan
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
