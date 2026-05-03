# Specification: {{SPEC_ID}}-{{SPEC_NAME}}

## Status

| Field | Value |
|-------|-------|
| **Created** | {{CREATED_DATE}} |
| **Current Phase** | {{CURRENT_PHASE}} |
| **Last Updated** | {{LAST_UPDATED}} |

## Documents

| Document | Status | Notes |
|----------|--------|-------|
| requirements.md | {{REQUIREMENTS_STATUS}} | {{REQUIREMENTS_NOTES}} |
| solution.md | {{SOLUTION_STATUS}} | {{SOLUTION_NOTES}} |

**Status values**: `pending` | `in_progress` | `completed` | `skipped`

## Decomposition

| Field | Value |
|-------|-------|
| **Tier** | {{DECOMPOSITION_TIER}} |
| **Status** | {{DECOMPOSITION_STATUS}} |

**Tier values**: `Direct` (no artifacts) | `Standard` (plan/) | `Factory` (manifest.md + units/ + scenarios/) | `None` (not yet chosen)

For Standard tier, see `plan/README.md`.
For Factory tier, see `manifest.md`, `units/`, `scenarios/`.
For Direct tier, no decomposition artifacts are produced — implement-direct reads requirements.md and solution.md directly.

## Decisions Log

| Date | Decision | Rationale |
|------|----------|-----------|
| {{DATE}} | {{DECISION}} | {{RATIONALE}} |

## Context

{{CONTEXT_NOTES}}

---
*This file is managed by the specify-meta skill.*
