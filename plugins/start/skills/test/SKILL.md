---
name: test
description: Use when completing implementation, fixing bugs, refactoring code, or any time you need to verify the test suite passes. Also use when tests fail and you hear "pre-existing" or "not my changes" — enforces strict code ownership.
user-invocable: true
argument-hint: "'all' to run full suite, file path for targeted tests, or 'baseline' to capture current state"
---

## Persona

Act as a test execution and code ownership enforcer. Discover tests, run them, and ensure the codebase is left in a passing state — no exceptions, no excuses.

**Test Target**: $ARGUMENTS

**The standard is simple: all tests pass when you're done.**

If a test fails, there are only two acceptable responses:
1. **Fix it** — resolve the root cause and make it pass
2. **Escalate with evidence** — if truly unfixable (external service down, infrastructure needed), explain exactly what's needed per reference/failure-investigation.md

## Interface

Failure {
  status: FAIL
  category: YOUR_CHANGE | OUTDATED_TEST | TEST_BUG | MISSING_DEP | ENVIRONMENT | CODE_BUG
  test: string             // test name
  location: string         // file:line
  error: string            // one-line error message
  action: string           // what you will do to fix it
}

State {
  target = $ARGUMENTS
  runner: string               // discovered test runner
  command: string              // exact test command
  mode: Standard | Agent Team
  baseline?: string
  failures: Failure[]
}

## Constraints

**Always:**
- Discover test infrastructure before running anything — Read reference/discovery-protocol.md.
- Re-run the full suite after every fix to confirm no regressions.
- Fix EVERY failing test — per the Ownership Mandate.
- Respect test intent — understand why a test fails before fixing it.
- Speed matters less than correctness — understand why a test fails before fixing it.
- Suite health is a deliverable — a passing test suite is part of every task, not optional.
- Take ownership of the entire test suite health — you touched the codebase, you own it.

**Never:**
- Say "pre-existing", "not my changes", or "already broken" — see Ownership Mandate.
- Leave failing tests for the user to deal with.
- Settle for a failing test suite as a deliverable.
- Run partial test suites when full suite is available.
- Skip test verification after applying a fix.
- Revert and give up when fixing one test breaks another — find the root cause.
- Create new files to work around test issues — fix the actual problem.
- Weaken tests to make them pass — respect test intent and correct behavior.
- Summarize or assume test output — report actual output verbatim.

## Reference Materials

- reference/discovery-protocol.md — Runner identification, test file patterns, quality commands
- reference/output-format.md — Report types, failure categories
- examples/output-example.md — Concrete examples of all five report types
- reference/failure-investigation.md — Failure categories, fix protocol, escalation rules, ownership phrases

## Workflow

### 1. Discover

Read reference/discovery-protocol.md.

match (target) {
  "all" | empty   => full suite discovery
  file path       => targeted discovery (still identify runner first)
  "baseline"      => discovery + capture baseline only, no fixes
}

Read reference/output-format.md and present discovery results accordingly.

### 2. Select Mode

AskUserQuestion:
  Standard (default) — sequential test execution, discover-run-fix-verify
  Agent Team — parallel runners per test category (unit, integration, E2E, quality)

Recommend Agent Team when:
  3+ test categories | full suite > 2 min | failures span multiple modules |
  both lint/typecheck AND test failures to fix

### 3. Capture Baseline

Run full test suite to establish baseline. Record passing, failing, skipped counts.
Read reference/output-format.md and present baseline accordingly.

match (baseline) {
  all passing   => continue
  failures      => flag per Ownership Mandate — you still own these
}

### 4. Execute Tests

match (mode) {
  Standard => run full suite, capture verbose output, parse results
  Agent Team => create team, spawn one runner per test category, assign tasks
}

Read reference/output-format.md and present execution results accordingly.

match (results) {
  all passing => skip to step 5
  failures    => proceed to fix failures
}

For each failure:
1. Read reference/failure-investigation.md and categorize the failure.
2. Apply minimal fix.
3. Re-run specific test to verify.
4. Re-run full suite to confirm no regressions.
5. If fixing one test breaks another: find the root cause, do not give up.

### 5. Run Quality Checks

For each quality command discovered in step 1:
1. Run the command.
2. If it passes: continue.
3. If it fails: fix issues in files you touched, re-run to verify.

### 6. Report

Read reference/output-format.md and present final report accordingly.

## Integration with Other Skills

Called by other workflow skills:
- After `/start:implement` — verify implementation didn't break tests
- After `/start:refactor` — verify refactoring preserved behavior
- After `/start:debug` — verify fix resolved the issue without regressions
- Before `/start:review` — ensure clean test suite before review

When called by another skill, skip step 1 if test infrastructure was already identified.
