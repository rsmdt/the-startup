# Failure Investigation

How to investigate, categorize, and fix every failing test.

---

## Failure Categories

| Cause Category | What to Look For | Action |
|---------------|-----------------|--------|
| **Your changes broke it** | Test was passing in baseline, fails after your changes | Fix the implementation or update the test to match new correct behavior |
| **Test is outdated** | Test assertions don't match current intended behavior | Update the test to match correct behavior |
| **Test has a bug** | Test logic is flawed (wrong assertion, bad mock, race condition) | Fix the test |
| **Missing dependency** | Import errors, missing fixtures, setup failures | Add the missing piece |
| **Environment issue** | Port conflicts, file locks, timing issues | Fix the environment setup |
| **Actual bug in code** | Test correctly catches a real bug | Fix the production code |

## Fix Protocol

For EVERY failing test:

1. **Read the failing test** — understand what it's testing and why
2. **Read the code under test** — understand the implementation
3. **Determine the correct fix** — fix the code, the test, or both
4. **Apply the fix** — edit the minimal set of files needed
5. **Re-run the specific test** — confirm the fix works
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
