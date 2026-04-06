# Example: Holdout Scenarios

Example scenarios for a SQL validation feature, organized by unit.

---

## Scenarios per Unit

### dm1 (2 scenarios)
- `schema-validation.md` [P0] — Create tables, verify schema matches specification
- `entity-relationships.md` [P1] — Verify foreign key constraints and cascades

### ve1 (3 scenarios)
- `valid-select.md` [P0] — POST valid SELECT, verify structured response
- `column-extraction.md` [P1] — POST complex JOIN, verify table/column extraction
- `complexity-scoring.md` [P1] — POST nested subquery, verify complexity score

### ve2 (2 scenarios)
- `sql-injection.md` [P0] — POST injection attempt, verify 400 response
- `empty-input.md` [P1] — POST empty/null SQL, verify validation error

### rl1 (2 scenarios)
- `burst-traffic.md` [P0] — Send 100 requests in 1 second, verify 429 after limit
- `rate-limit-headers.md` [P1] — Verify X-RateLimit headers in responses

**Total**: 9 scenarios across 4 units (100% unit coverage)

## Example Scenario File (ve2/sql-injection.md)

```markdown
---
unit: ve2
feature: input-sanitization
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
