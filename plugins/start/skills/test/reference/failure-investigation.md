# Failure Investigation

How to investigate, categorize, and fix every failing test.

---

## Failure Categories

| Cause Category | What to Look For | Action |
|---------------|-----------------|--------|
| **Your changes broke it** (YOUR_CHANGE) | Test was passing in baseline, fails after your changes | Fix the implementation or update the test to match new correct behavior |
| **Test is outdated** (OUTDATED_TEST) | Test assertions don't match current intended behavior, but the test still verifies real behavior | Update the test to match correct behavior |
| **Test has a bug** (TEST_BUG) | Test logic is flawed (wrong assertion, bad mock, race condition) | Fix the test |
| **Test is noise** (NOISE_TEST) | Test pins an implementation detail / framework behavior / call sequence rather than behavior. Common shapes: mirrors a trivial implementation (`assert format(x) == f"prefix:{x}"` against the same f-string); re-verifies a schema / ORM / router contract; mocks a collaborator and asserts the mock was called; mocks the data-access client so no real query runs. The test broke on a legitimate refactor and would never catch a real bug. | **Delete the test.** Verify behavior is covered by a caller's test, an integration test, or is implicit. Do not preserve the noise by updating it to match the new implementation. See reference/test-design-rules.md. |
| **Missing dependency** (MISSING_DEP) | Import errors, missing fixtures, setup failures | Add the missing piece |
| **Environment issue** (ENVIRONMENT) | Port conflicts, file locks, timing issues | Fix the environment setup |
| **Actual bug in code** (CODE_BUG) | Test correctly catches a real bug | Fix the production code |

### Distinguishing OUTDATED_TEST from NOISE_TEST

Both arise when a refactor breaks a test the implementation no longer satisfies. The difference is what the test was verifying:

- **OUTDATED_TEST** — test verified a real behavior (e.g., "discount is applied for orders over $100"); the behavior changed (e.g., threshold is now $50). Update the assertion.
- **NOISE_TEST** — test pinned an implementation detail (e.g., "calls `repo.find_by_id` then `repo.update`"); the implementation changed (e.g., now uses `repo.upsert`) but the user-visible behavior is identical. Delete the test; it never produced signal.

Ask: *"Could this test have failed for a reason a caller's test wouldn't have caught?"* If no, it's NOISE_TEST.

## Fix Protocol

For EVERY failing test:

1. **Read the failing test** — understand what it's testing and why
2. **Read the code under test** — understand the implementation
3. **Determine the correct response** — fix the code, fix/update the test, or **delete the test** (if NOISE_TEST per reference/test-design-rules.md). Deletion is a valid resolution, not a workaround.
4. **Apply the fix** — edit the minimal set of files needed; for NOISE_TEST, remove the test and confirm behavior is covered elsewhere
5. **Re-run the specific test (or, if deleted, the related caller/integration tests)** — confirm the resolution holds
6. **Re-run the full suite** — confirm no regressions

## Iterate

Repeat until ALL tests pass. If fixing one test breaks another:
- Do NOT revert and give up
- Investigate the chain of dependencies
- Find the root cause that satisfies all tests

## Escalation (Last Resort)

This is ONLY acceptable for:
- External service dependencies that are down
- Infrastructure requirements beyond the codebase (e.g., database migration needed)
- Permission/access issues

This is NOT acceptable for:
- "Complex" code you don't understand — Read it more carefully
- "Might break something else" — Run the tests and find out
- "Not my responsibility" — Yes it is. You touched the codebase.

Use the escalation template from `output-format.md`.

## Ownership Enforcement Phrases

When you catch yourself about to deflect, replace with ownership language:

| Instead of... | Say... |
|---------------|--------|
| "This test was already failing" | "This test is failing. Let me fix it." |
| "Not caused by my changes" | "The test suite needs to pass. Let me investigate." |
| "Pre-existing issue" | "Found a failing test. Fixing it now." |
| "This is outside the scope" | "I see a failing test. The suite needs to be green." |
| "The test might be flaky" | "Let me run it again and if it fails, fix the root cause." |
| "I'd recommend fixing this separately" | "I'm fixing this now." |
| "This appears to be a known issue" | "I'm making this a fixed issue." |
| "I'll update the test to match the new implementation" | "Does this test produce signal, or pin an implementation detail? If the latter, deleting it is the correct fix." |
| "This test is too coupled to the internals to keep working through refactors" | "This is a NOISE_TEST. Deleting it and verifying behavior coverage elsewhere is the resolution." |
