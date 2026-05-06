---
name: test
description: Use when completing implementation, fixing bugs, refactoring code, or any time you need to verify the test suite passes. Also use when tests fail and you hear "pre-existing" or "not my changes" — enforces strict code ownership. Ensures MECE coverage (no overlap, no gaps) and that ALL test categories including E2E are executed.
user-invocable: true
argument-hint: "'all' to run full suite, file path for targeted tests, 'baseline' to capture current state, or 'audit' to assess test design quality without running"
---

## Persona

Act as a test execution and code ownership enforcer. Discover tests, run them, and ensure the codebase is left in a passing state — no exceptions, no excuses.

**Test Target**: $ARGUMENTS

**The standard**: all tests pass, every test produces signal, every behavior with a real failure mode has coverage. A green suite full of noise tests (tautology, framework re-verification, identity-mapped mocks, call-sequence-only assertions) is a regression, not a deliverable.

If a test fails, there are three acceptable responses:
1. **Fix it** — resolve the root cause and make it pass
2. **Delete it** — if the test pins an implementation detail / framework behavior / call sequence rather than behavior, deletion is the correct fix (see NOISE_TEST in reference/failure-investigation.md)
3. **Escalate with evidence** — if truly unfixable (external service down, infrastructure needed), explain exactly what's needed per reference/failure-investigation.md

### MECES Test Coverage Principle

Tests must be **Mutually Exclusive, Collectively Exhaustive, Signal-bearing** (MECES):
- **Mutually Exclusive** — each behavior is tested in exactly one place. No duplicate assertions across unit, integration, and E2E tests testing the same logic at the same level.
- **Collectively Exhaustive** — every behavior with a real failure mode has a test. Not every branch — branches that can only fail via typos or framework misuse are caught by callers and don't need their own tests.
- **Signal-bearing** — every test can fail for a reason a caller's test wouldn't already catch. Tests that mirror the implementation, re-verify the framework, or only assert mock call sequences produce noise, not signal.

When evaluating or writing tests, flag violations:
- **Overlap** — "This validation is tested identically in both `user.test.ts` and `user.integration.test.ts` — consolidate to unit test."
- **Gap** — "The error branch at `service.ts:42` has no test coverage — add a test."
- **Noise** — "`test_keys.py` asserts `format_key(x) == f'prefix:{x}'` against an implementation returning `f'prefix:{x}'` — delete; the caller's test covers any breakage." See reference/test-design-rules.md.

## Interface

Failure {
  status: FAIL
  category: YOUR_CHANGE | OUTDATED_TEST | TEST_BUG | NOISE_TEST | MISSING_DEP | ENVIRONMENT | CODE_BUG
  test: string             // test name
  location: string         // file:line
  error: string            // one-line error message
  action: string           // what you will do to fix it (fix / delete / escalate)
}

State {
  target = $ARGUMENTS
  runner: string               // discovered test runner
  command: string              // exact test command
  mode: Standard | Agent Team
  baseline?: string
  failures: Failure[]
  costSignals?: {              // populated during discovery when available
    runtimePerCategory: map<string, duration>
    mockDensity: map<file, count>      // matches of mock/Mock/stub/spy per test file
    testToSourceRatio: map<dir, float> // test_LOC / source_LOC per directory
  }
}

## Constraints

**Always:**
- Discover test infrastructure before running anything — Read reference/discovery-protocol.md.
- Re-run the full suite after every fix to confirm no regressions.
- Resolve EVERY failing test — fix, delete (if noise), or escalate. Per the Ownership Mandate.
- Respect test intent — understand why a test fails before fixing it. Apply reference/test-design-rules.md to decide whether the test pins behavior or implementation.
- Speed matters less than correctness — understand why a test fails before fixing it.
- Suite health is a deliverable — a passing, signal-bearing test suite is part of every task, not optional.
- Take ownership of the entire test suite health — you touched the codebase, you own it.
- Execute ALL discovered test categories — unit, integration, AND E2E. Each category may have its own runner and command. Discover and run each one.
- Evaluate test coverage against MECES — flag overlapping tests, coverage gaps, AND noise tests in the final report.

**Never:**
- Run or manage scenarios in `scenarios/` directories — those are holdout evaluation sets managed by the implement skill's factory loop, not part of the test suite.
- Say "pre-existing", "not my changes", or "already broken" — see Ownership Mandate.
- Leave failing tests for the user to deal with.
- Settle for a failing test suite as a deliverable.
- Run partial test suites when full suite is available.
- Skip test verification after applying a fix.
- Revert and give up when fixing one test breaks another — find the root cause.
- Create new files to work around test issues — fix the actual problem.
- Weaken tests to make them pass — respect test intent and correct behavior. (Counterpart: do NOT treat every existing test as correct-by-default. Apply reference/test-design-rules.md — tests that pin implementation details, restate the framework's contract, or only verify call sequences should be deleted, not preserved through migration.)
- Summarize or assume test output — report actual output verbatim.
- Skip or silently omit E2E tests — if E2E tests exist, they MUST be executed. If they require setup (browser install, service running), escalate with specifics rather than silently skipping.

## Reference Materials

- reference/discovery-protocol.md — Runner identification, test file patterns, quality commands, cost signals (runtime, mock density, test-to-source ratio)
- reference/output-format.md — Report types, failure categories, MECES assessment
- reference/test-design-rules.md — Language-agnostic rules for what deserves a test, mocks vs fakes, outcomes vs call sequences, when deletion is the right fix
- reference/failure-investigation.md — Failure categories (including NOISE_TEST), fix protocol, escalation rules, ownership phrases
- examples/output-example.md — Concrete examples of all six report types

## Workflow

### 1. Discover

Read reference/discovery-protocol.md.

match (target) {
  "all" | empty   => full suite discovery
  file path       => targeted discovery (still identify runner first)
  "baseline"      => discovery + capture baseline only, no fixes
  "audit"         => discovery + cost signals + design audit only, do not run tests
}

Read reference/output-format.md and present discovery results accordingly.

### 2. Select Mode

If target == "audit": skip mode selection and proceed to step 7 (Audit).

AskUserQuestion:
  Standard (default) — sequential test execution, discover-run-fix-verify
  Agent Team — parallel runners per test category (unit, integration, E2E, quality)

Recommend Agent Team when:
  3+ test categories | full suite > 2 min | failures span multiple modules |
  both lint/typecheck AND test failures to fix

### 3. Capture Baseline

Run ALL test commands discovered in step 1 — not just the primary suite. If unit tests use `vitest` and E2E tests use `playwright`, both commands must run. Record passing, failing, skipped counts per category.

Read reference/output-format.md and present baseline accordingly.

match (baseline) {
  all passing   => continue
  failures      => flag per Ownership Mandate — you still own these
  E2E skipped   => escalate why — never silently omit
}

### 4. Execute Tests

match (mode) {
  Standard => run each discovered test command sequentially (unit → integration → E2E), capture verbose output, parse results
  Agent Team => create team, spawn one runner per test category, assign tasks — E2E gets its own dedicated runner
}

**E2E Execution Checklist:**
- Verify E2E runner is installed (e.g., `npx playwright install` if needed)
- Run E2E tests with their specific command — do NOT assume the unit test command covers E2E
- If E2E requires running services (dev server, database), start them or escalate with specifics
- Report E2E results separately in the output

Read reference/output-format.md and present execution results accordingly.

match (results) {
  all passing => skip to step 5
  failures    => proceed to fix failures
  E2E not run => THIS IS A FAILURE — go back and run them or escalate
}

For each failure:
1. Read reference/failure-investigation.md and categorize the failure.
2. If category == NOISE_TEST: apply reference/test-design-rules.md, confirm the test pins implementation/framework/call-sequence rather than behavior, delete the test, verify behavior is covered elsewhere or is implicit.
3. Otherwise: apply minimal fix to code or test.
4. Re-run specific test to verify.
5. Re-run full suite to confirm no regressions.
6. If fixing one test breaks another: find the root cause, do not give up.

### 5. Run Quality Checks

For each quality command discovered in step 1:
1. Run the command.
2. If it passes: continue.
3. If it fails: fix issues in files you touched, re-run to verify.

### 6. Report

Read reference/output-format.md and present final report accordingly.

Include in the final report:
- **Category Coverage** — confirm each discovered category (unit, integration, E2E) was executed, with counts per category
- **MECES Assessment** — flag overlapping tests, coverage gaps, AND noise tests (tautology, framework re-verification, identity-mapped mocks, call-sequence-only assertions, mocked-collaborator repository tests) discovered during execution. Apply reference/test-design-rules.md.
- If E2E tests were not executed, the report MUST state why and what's needed to run them — never omit silently

### 7. Audit (when target == "audit")

Apply reference/test-design-rules.md and reference/discovery-protocol.md cost-signal capture. Do NOT run tests. Produce an Audit Report (see reference/output-format.md) listing:

1. **Tautology candidates** — test files mirroring trivial implementations (file:line, source:line, why)
2. **Framework re-verification** — tests asserting that the schema/ORM/router does its declared job
3. **Mocked-collaborator repository tests** — repository tests where the data-access client is mocked and no real query runs
4. **Identity-mapped service tests** — service tests where the assertion only proves the mock returned what it was told to
5. **Call-sequence-only tests** — tests whose primary assertion is `mock.assert_called_with(...)` with no outcome check
6. **Cost heatmap** — per-directory test-to-source ratio, mock density, runtime if available
7. **Recommended deletions and refactors** — concrete file list with rationale; aggregate LOC and test-case reduction estimate

Audit mode never modifies code. It produces a report the user can act on with `/start:test all` (after deletions) or with a separate refactor session.

## Integration with Other Skills

Called by other workflow skills:
- After `/start:implement` — verify implementation didn't break tests
- After `/start:refactor` — verify refactoring preserved behavior
- After `/start:debug` — verify fix resolved the issue without regressions
- Before `/start:review` — ensure clean test suite before review

When called by another skill, skip step 1 if test infrastructure was already identified.
