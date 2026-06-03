# Example: Holdout Scenarios

Example scenarios for the SQL Tools API, organized by unit. Each unit's scenarios exercise **its own** endpoint, so they pass when that unit (plus its dependency) is built — the test of a true vertical slice.

---

## Scenarios per Unit

### vq1 — Validate a query (4 scenarios)
- `valid-select.md` [P0] — POST a valid SELECT to /validate, verify structured response (valid, extracted tables/columns, complexity)
- `injection-rejected.md` [P0] — POST an injection attempt to /validate, verify 400 with structured errors
- `table-extraction.md` [P1] — POST a complex JOIN to /validate, verify table/column extraction
- `complexity-scoring.md` [P1] — POST a nested subquery to /validate, verify complexity score

### eq1 — Explain a query (2 scenarios)
- `explain-cost.md` [P0] — POST a valid SELECT to /explain, verify a human-readable cost explanation
- `explain-rejects-invalid.md` [P1] — POST invalid SQL to /explain, verify it surfaces validation errors (extends vq1's analysis)

### oq1 — Optimize a query (2 scenarios)
- `suggest-index.md` [P0] — POST a query with a full-table scan to /optimize, verify an index suggestion is returned
- `optimize-noop.md` [P1] — POST an already-optimal query to /optimize, verify no spurious suggestions

**Total**: 8 scenarios across 3 units (100% unit coverage)

> Note how each unit's scenarios target its **own** endpoint. The injection check lives in `vq1`'s scenarios (sanitization is part of the validate slice, not a separate unit), so `vq1` is fully evaluable alone. `eq1`/`oq1` scenarios hit /explain and /optimize and pass once `vq1` (their dependency) plus the unit itself is built.

## Example Scenario File (vq1/injection-rejected.md)

```markdown
---
unit: vq1
feature: query-validation
priority: P0
---
# SQL Injection Detection

## Scenario
POST to /api/v1/sql/validate with:
  sql: "SELECT * FROM users; DROP TABLE users; --"

## Expected
- valid should be false
- errors should mention disallowed statement types
- response status should be 400, not 500
```
