# Output Format Reference

Guidelines for test output. See `examples/output-example.md` for concrete rendered examples of all five report types.

---

## Report Types

Five report types used during test execution:

1. **Discovery Results** — after identifying test infrastructure, including per-category runners and commands
2. **Baseline Captured** — before any changes, with pre-existing failures flagged, counts per category
3. **Execution Results** — after running tests, with failure categorization, per-category breakdown
4. **Escalation** — only for external blockers (service down, infrastructure, permissions)
5. **Final Report** — comprehensive summary with quality checks, MECE assessment, and per-category coverage

## Per-Category Reporting

Every report that includes test counts MUST break them down by category:

| Category | Command | Tests | Passing | Failing |
|----------|---------|-------|---------|---------|
| Unit | `npx vitest run` | 180 | 180 | 0 |
| Integration | `npx vitest run tests/integration` | 32 | 31 | 1 |
| E2E | `npx playwright test` | 12 | 12 | 0 |

If a category was **not executed**, it MUST appear with status `⚠️ NOT RUN` and a reason.

## MECE Assessment (Final Report)

The final report includes a MECE section:
- **Overlaps found** — tests covering the same behavior at the same abstraction level (list specific files)
- **Gaps found** — untested code paths, branches, or edge cases (list specific locations)
- **Assessment** — ✅ MECE / ⚠️ Overlaps detected / ❌ Gaps detected

## Failure Categories

| Category | Meaning |
|----------|---------|
| YOUR_CHANGE | Your code caused this failure |
| OUTDATED_TEST | Test expectations no longer match current behavior |
| TEST_BUG | Test itself has a bug (wrong assertion, bad setup) |
| MISSING_DEP | Missing dependency or fixture |
| ENVIRONMENT | Environment-specific issue (CI, config, network) |
| CODE_BUG | Pre-existing bug in application code |
