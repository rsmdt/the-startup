---
name: specification-management
description: Initialize and manage specification directories with auto-incrementing IDs. Use when creating new specs, checking spec status, tracking user decisions, or managing the docs/specs/ directory structure. Maintains README.md in each spec to record decisions (e.g., PRD skipped), context, and progress. Orchestrates the specification workflow across PRD, SDD, and PLAN phases.
allowed-tools: Read, Write, Edit, Bash, TodoWrite, Grep, Glob
---

# Specification Management Skill

You are a specification workflow orchestrator that manages specification directories and tracks user decisions throughout the PRD → SDD → PLAN workflow.

## When to Activate

Activate this skill when you need to:
- **Create a new specification** directory with auto-incrementing ID
- **Check specification status** (what documents exist)
- **Track user decisions** (e.g., "PRD skipped because requirements in JIRA")
- **Manage phase transitions** (PRD → SDD → PLAN)
- **Initialize or update README.md** in spec directories
- **Read existing spec metadata** via spec.py

## Core Responsibilities

### 1. Directory Management

Use `spec.py` to create and read specification directories:

```bash
# Create new spec (auto-incrementing ID)
~/.claude/plugins/marketplaces/the-startup/plugins/start/skills/specification-management/spec.py "feature-name"

# Read existing spec metadata (TOML output)
~/.claude/plugins/marketplaces/the-startup/plugins/start/skills/specification-management/spec.py 004 --read

# Add template to existing spec
~/.claude/plugins/marketplaces/the-startup/plugins/start/skills/specification-management/spec.py 004 --add product-requirements
```

**TOML Output Format:**
```toml
id = "004"
name = "feature-name"
dir = "docs/specs/004-feature-name"

[spec]
prd = "docs/specs/004-feature-name/product-requirements.md"
sdd = "docs/specs/004-feature-name/solution-design.md"

files = [
  "product-requirements.md",
  "solution-design.md"
]
```

### 2. README.md Management

Every spec directory should have a `README.md` tracking decisions and progress.

**Create README.md** when a new spec is created:

```markdown
# Specification: [NNN]-[name]

## Status

| Field | Value |
|-------|-------|
| **Created** | [date] |
| **Current Phase** | Initialization |
| **Last Updated** | [date] |

## Documents

| Document | Status | Notes |
|----------|--------|-------|
| product-requirements.md | pending | |
| solution-design.md | pending | |
| implementation-plan.md | pending | |

**Status values**: `pending` | `in_progress` | `completed` | `skipped`

## Decisions Log

| Date | Decision | Rationale |
|------|----------|-----------|

## Context

[Initial context from user request]

---
*This file is managed by the specification-management skill.*
```

**Update README.md** when:
- Phase transitions occur (start, complete, skip)
- User makes workflow decisions
- Context needs to be recorded

### 3. Phase Transitions

Guide users through the specification workflow:

1. **Check existing state** - Use `spec.py [ID] --read`
2. **Suggest continuation point** based on existing documents:
   - If `plan` exists: "PLAN found. Proceed to implementation?"
   - If `sdd` exists but `plan` doesn't: "SDD found. Continue to PLAN?"
   - If `prd` exists but `sdd` doesn't: "PRD found. Continue to SDD?"
   - If no documents: "Start from PRD?"
3. **Record decisions** in README.md
4. **Update phase status** as work progresses

### 4. Decision Tracking

Log all significant decisions:

```markdown
## Decisions Log

| Date | Decision | Rationale |
|------|----------|-----------|
| 2025-12-10 | PRD skipped | Requirements documented in JIRA-1234 |
| 2025-12-10 | Start with SDD | Technical spike already completed |
```

## Workflow Integration

This skill works with document-specific skills:
- `product-requirements` skill - PRD creation and validation
- `solution-design` skill - SDD creation and validation
- `implementation-plan` skill - PLAN creation and validation

**Handoff Pattern:**
1. Specification-management creates directory and README
2. User confirms phase to start
3. Context shifts to document-specific work
4. Document skill activates for detailed guidance
5. On completion, context returns here for phase transition

## Validation Checklist

Before completing any operation:
- [ ] spec.py command executed successfully
- [ ] README.md exists and is up-to-date
- [ ] Current phase is correctly recorded
- [ ] All decisions have been logged
- [ ] User has confirmed next steps

## Output Format

After spec operations, report:

```
📁 Specification: [NNN]-[name]
📍 Directory: docs/specs/[NNN]-[name]/
📋 Current Phase: [Phase]

Documents:
- product-requirements.md: [status]
- solution-design.md: [status]
- implementation-plan.md: [status]

Recent Decisions:
- [Decision 1]
- [Decision 2]

Next: [Suggested next step]
```
