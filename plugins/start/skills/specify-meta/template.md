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
| requirements.mdx | {{REQUIREMENTS_STATUS}} | WHY — {{REQUIREMENTS_NOTES}} |
| solution.mdx | {{SOLUTION_STATUS}} | WHAT (design) — {{SOLUTION_NOTES}} |
| plan.mdx | {{PLAN_STATUS}} | WHAT (completion contract: target state, quality gates, completion criteria) — {{PLAN_NOTES}} |

**Status values**: `pending` | `in_progress` | `completed` | `skipped`

## Implementation

The specification describes the target state. **How** to implement it is decided by `/implement`,
which selects an execution strategy and records it (with any generated decomposition) under `.impl/`.

| Field | Value |
|-------|-------|
| **Strategy** | {{IMPL_STRATEGY}} |

**Strategy values**: `Direct` | `Iterative` | `Decompose` | `None` (not yet implemented). Chosen at implementation time, not specification time. See `.impl/strategy.md` once `/implement` has run.

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
