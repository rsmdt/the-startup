# Output Format Reference

Guidelines for test output. See `examples/output-example.md` for concrete rendered examples of all six report types.

---

## Report Types

Six report types used during test execution:

1. **Discovery Results** — after identifying test infrastructure, including per-category runners, commands, and cost signals (runtime, mock density, test-to-source ratio)
2. **Baseline Captured** — before any changes, with pre-existing failures flagged, counts per category
3. **Execution Results** — after running tests, with failure categorization, per-category breakdown
4. **Escalation** — only for external blockers (service down, infrastructure, permissions)
5. **Final Report** — comprehensive summary with quality checks, MECES assessment, and per-category coverage
6. **Audit Report** — produced when target == "audit"; assesses test-design quality without running tests, lists noise-test candidates with deletion recommendations and cost-signal heatmap

## Audit Report Structure

When target == "audit", the report does not include execution results. Instead:

1. **Discovery Summary** — runners, commands, file counts per category
2. **Cost Signals** — table of per-directory test-to-source LOC ratio, mock density (matches of `mock|Mock|stub|spy` per file), per-category runtime if known. Flag directories where ratio > 0.7 with high mock density as audit hotspots.
3. **Noise Findings** — grouped by the six noise categories above. Each finding cites file:line, source under test (file:line), and a one-line rationale.
4. **Worth-Preserving Patterns** — examples of high-signal tests in this codebase (pure logic, security gates, fake-based coordinators) so the team knows the model to extend.
5. **Recommended Deletions** — concrete file/test list with aggregate LOC and test-case reduction estimate.
6. **Recommended Refactors** — files where the scenarios are valuable but the shape is wrong (typically: rewrite mocks as fakes following the gold-standard examples found in step 4).
7. **MECES Assessment** — overlaps + gaps + noise summary.

The audit report MUST NOT modify code. It is input for a subsequent `/start:test all` (after deletions) or `/start:refactor` session.

## Per-Category Reporting

Every report that includes test counts MUST break them down by category:

| Category | Command | Tests | Passing | Failing |
|----------|---------|-------|---------|---------|
| Unit | `npx vitest run` | 180 | 180 | 0 |
| Integration | `npx vitest run tests/integration` | 32 | 31 | 1 |
| E2E | `npx playwright test` | 12 | 12 | 0 |

If a category was **not executed**, it MUST appear with status `⚠️ NOT RUN` and a reason.

## MECES Assessment (Final Report and Audit Report)

The final and audit reports include a MECES section with three axes:

- **Overlaps found** — tests covering the same behavior at the same abstraction level (list specific files)
- **Gaps found** — behaviors with real failure modes that lack test coverage (list specific locations). Trivial branches that can only fail via typo / framework misuse don't count as gaps — they're caught by callers.
- **Noise found** — tests that pin implementation rather than verify behavior. Categorize each:
  - **Tautology** — test mirrors the implementation (e.g., asserts a literal f-string against an f-string return)
  - **Framework re-verification** — asserts what the schema lib / ORM / router / serializer already enforces
  - **Mocked-collaborator repository test** — data-access client is mocked; no real query runs
  - **Identity-mapped service test** — mock returns X, assertion is "service returned X"; no transformation verified
  - **Call-sequence-only** — primary assertion is `mock.assert_called_with(...)` with no outcome check
  - **Trivial wrapper** — function under test is a one-line literal/template/delegation; caller's test catches breakage
- **Assessment** — ✅ MECES / ⚠️ Overlaps detected / ❌ Gaps detected / 🔇 Noise detected (combinable)

For each noise finding, recommend: **Delete** (preferred when behavior is implicitly covered) or **Refactor to fakes** (when the underlying scenario is real but the test was wrongly shaped). See reference/test-design-rules.md.

## Failure Categories

| Category | Meaning | Typical action |
|----------|---------|----------------|
| YOUR_CHANGE | Your code caused this failure | Fix code or update test |
| OUTDATED_TEST | Test expectations no longer match current behavior, but test verifies real behavior | Update assertion |
| TEST_BUG | Test itself has a bug (wrong assertion, bad setup) | Fix test |
| NOISE_TEST | Test pins implementation / framework behavior / call sequence — not real behavior | **Delete the test**; verify behavior is covered elsewhere |
| MISSING_DEP | Missing dependency or fixture | Add the missing piece |
| ENVIRONMENT | Environment-specific issue (CI, config, network) | Fix environment |
| CODE_BUG | Pre-existing bug in application code | Fix production code |
