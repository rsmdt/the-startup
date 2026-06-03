# Example: Unit Decomposition

Example unit decomposition for a SQL Tools API. Every unit is a **vertical slice** — one user-observable capability, end-to-end.

---

## Decomposition Table

| ID | Title | Dependencies | Requirements Covered |
|----|-------|--------------|---------------------|
| vq1 | Validate a query | none | R1 (persist analysis), R2 (syntax + reject non-SELECT), R3 (injection + structured errors), R4 (extract refs), R5 (complexity score) |
| eq1 | Explain a query | vq1 | R6 (human-readable cost explanation) |
| oq1 | Optimize a query | vq1 | R7 (index/rewrite suggestions) |

**Coverage**: 7/7 requirements mapped (100%)

> **Why this is vertical, not horizontal.** A tempting horizontal cut would be `dm1` (data models) → `ve1` (validation logic) → `ve2` (input sanitization) → `rl1` (rate limiting). That produces four units, but only the data unit and the logic unit are runnable alone — `ve2` and `rl1` can't pass scenarios without `ve1`'s endpoint, so they're layers, not slices.
>
> Sliced vertically, the *first capability* (`vq1`) builds the data model and folds sanitization in — it owns its shared foundation. The next two capabilities (`eq1`, `oq1`) depend on `vq1` and **extend** the model. Each unit's scenarios hit its own endpoint and pass independently. Same coverage, every unit independently evaluable, and `eq1`/`oq1` run in parallel.

## Example Unit Spec (vq1)

```markdown
---
id: vq1
title: Validate a Query
type: feature
dependencies: []
---
# Validate a Query

## Goal
POST /api/v1/sql/validate — validates SQL without executing it and returns a structured analysis.

## Requirements
- Persist a QueryAnalysis record for each request (this unit owns the QueryAnalysis model; later units extend it)
- Validate syntax and reject non-SELECT statements
- Detect injection patterns and reject them with structured errors (sanitization is part of this slice, not a separate unit)
- Extract table and column references from valid queries
- Compute a complexity score based on joins, subqueries, and aggregations
- Return a structured JSON response with all findings

## Constraints
- Follow existing controller patterns (see AGENTS.md)
- Full test coverage for all validation rules
- No new dependencies
- Use existing JSQLParser for parsing

## Interfaces
- Defines the QueryAnalysis model and the /api/v1/sql/validate response shape, both consumed and extended by eq1 and oq1.
```
