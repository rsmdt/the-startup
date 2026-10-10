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
| requirements.md | {{REQUIREMENTS_STATUS}} | WHY — {{REQUIREMENTS_NOTES}} |
| solution.md | {{SOLUTION_STATUS}} | WHAT (design) — {{SOLUTION_NOTES}} |
| plan.md | {{PLAN_STATUS}} | WHAT (narrative plan; slice and scenario tasks in `.tracker/`) — {{PLAN_NOTES}} |

**Status values**: `pending` | `in_progress` | `completed` | `skipped`

## Implementation

The specification describes the target state. `/implement` builds it by running each slice task in
`.tracker/` through a code→eval loop and records each slice's status on its task
(`task.js graph --spec {{SPEC_ID}}`).

## Decisions Log

| Date | Decision | Rationale |
|------|----------|-----------|
| {{DATE}} | {{DECISION}} | {{RATIONALE}} |

## Assumption Ledger

Maintained by the specify skill's Clarify Intent phase. Standing entries appear as `[ASSUMPTION]` markers in the documents.

| Topic | Default | Impact if wrong | Status |
|-------|---------|-----------------|--------|
| {{TOPIC}} | {{DEFAULT}} | {{IMPACT}} | Proposed \| Confirmed \| Corrected \| Standing |

## Context

{{CONTEXT_NOTES}}

---
*This file is managed by the specify-meta skill.*
