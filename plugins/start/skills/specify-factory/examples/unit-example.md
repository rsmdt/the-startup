# Example: Unit Decomposition

Example unit decomposition from a SQL validation feature.

---

## Decomposition Table

| ID | Title | Dependencies | Requirements Covered |
|----|-------|--------------|---------------------|
| dm1 | Data models | none | R1 (schema), R2 (entities) |
| ve1 | Validation endpoint | dm1 | R3 (validate), R4 (extract), R5 (complexity) |
| ve2 | Input sanitization | ve1 | R6 (injection), R7 (error handling) |
| rl1 | Rate limiting | ve1 | R8 (throttle), R9 (headers) |

**Coverage**: 9/9 requirements mapped (100%)

## Example Unit Spec (ve1)

```markdown
---
id: ve1
title: SQL Validation Endpoint
type: feature
dependencies: [dm1]
---
# SQL Validation Endpoint

## Goal
POST /api/v1/sql/validate — validates SQL without executing it.

## Requirements
- Validate syntax, reject non-SELECT statements, detect injection patterns
- Extract table and column references from valid queries
- Compute complexity score based on joins, subqueries, aggregations
- Return structured JSON response with all findings

## Constraints
- Follow existing controller patterns (see AGENTS.md)
- Full test coverage for all validation rules
- No new dependencies
- Use existing JSQLParser for parsing
```
