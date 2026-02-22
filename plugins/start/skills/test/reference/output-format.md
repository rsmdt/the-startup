# Output Format Reference

Guidelines for test output. See `examples/output-example.md` for concrete rendered examples of all five report types.

---

## Report Types

Five report types used during test execution:

1. **Discovery Results** — after identifying test infrastructure
2. **Baseline Captured** — before any changes, with pre-existing failures flagged
3. **Execution Results** — after running tests, with failure categorization
4. **Escalation** — only for external blockers (service down, infrastructure, permissions)
5. **Final Report** — comprehensive summary with quality checks

## Failure Categories

| Category | Meaning |
|----------|---------|
| YOUR_CHANGE | Your code caused this failure |
| OUTDATED_TEST | Test expectations no longer match current behavior |
| TEST_BUG | Test itself has a bug (wrong assertion, bad setup) |
| MISSING_DEP | Missing dependency or fixture |
| ENVIRONMENT | Environment-specific issue (CI, config, network) |
| CODE_BUG | Pre-existing bug in application code |
