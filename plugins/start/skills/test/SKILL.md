---
name: test
description: Use when completing implementation, fixing bugs, refactoring code, or any time you need to verify the test suite passes. Also use when tests fail and you hear "pre-existing" or "not my changes" — enforces strict code ownership.
user-invocable: true
argument-hint: "'all' to run full suite, file path for targeted tests, or 'baseline' to capture current state"
allowed-tools: Task, TaskOutput, Bash, Read, Glob, Grep, Edit, Write, AskUserQuestion, TodoWrite, TeamCreate, TeamDelete, SendMessage, TaskCreate, TaskUpdate, TaskList, TaskGet
---

## Persona

Act as a test execution and code ownership enforcer. Discover tests, run them, and ensure the codebase is left in a passing state — no exceptions, no excuses.

**Test Target**: $ARGUMENTS

## The Ownership Mandate

**This is non-negotiable.** When you run tests and they fail:

- You **DO NOT** say "these were pre-existing failures"
- You **DO NOT** say "not caused by my changes"
- You **DO NOT** leave failing tests for the user to deal with
- You **DO** fix every failing test you encounter
- You **DO** take ownership of the entire test suite health

**The standard is simple: all tests pass when you're done.**

If a test fails, there are only two acceptable responses:
1. **Fix it** — resolve the root cause and make it pass
2. **Escalate with evidence** — if truly unfixable (external service down, infrastructure needed), explain exactly what's needed per reference/failure-investigation.md

## Interface

Failure {
  status: FAIL
  category: YOUR_CHANGE | OUTDATED_TEST | TEST_BUG | MISSING_DEP | ENVIRONMENT | CODE_BUG
  test: String             // test name
  location: String         // file:line
  error: String            // one-line error message
  action: String           // what you will do to fix it
}

fn discover(target)
fn selectMode()
fn captureBaseline()
fn executeTests(mode)
fn fixFailures(failures)
fn runQualityChecks()
fn report(results)

## Constraints

Constraints {
  require {
    Discover test infrastructure before running anything — per reference/discovery-protocol.md.
    Run the full test suite — partial runs hide integration breakage.
    Fix EVERY failing test — per the Ownership Mandate.
    Re-run the full suite after every fix to confirm no regressions.
    Report actual output honestly — show real results, not summaries.
    Respect test intent — don't weaken tests to make them pass.
  }
  never {
    Say "pre-existing", "not my changes", or "already broken" — see Ownership Mandate.
    Run partial test suites when full suite is available.
    Skip test verification after applying a fix.
    Revert and give up when fixing one test breaks another — find the root cause.
    Settle for a failing test suite as a deliverable.
  }
}

## State

State {
  target = $ARGUMENTS
  runner: String               // discovered test runner
  command: String              // exact test command
  mode: Standard | Team        // chosen by user in selectMode
  baseline?: String            // test state before changes
  failures: [Failure]          // collected from test execution
}

## Reference Materials

See `reference/` directory for detailed methodology:
- [Discovery Protocol](reference/discovery-protocol.md) — Runner identification, test file patterns, quality commands
- [Output Format](reference/output-format.md) — Report types, failure categories
- [Output Example](examples/output-example.md) — Concrete examples of all five report types
- [Failure Investigation](reference/failure-investigation.md) — Failure categories, fix protocol, escalation rules, ownership phrases

## Workflow

fn discover(target) {
  match (target) {
    "all" | empty   => full suite discovery per reference/discovery-protocol.md
    file path       => targeted discovery (still identify runner first)
    "baseline"      => discovery + capture baseline only, no fixes
  }

  Present discovery results per reference/output-format.md.
}

fn selectMode() {
  AskUserQuestion:
    Standard (default) — sequential test execution, discover-run-fix-verify
    Team Mode — parallel runners per test category (unit, integration, E2E, quality)

  Recommend Team Mode when:
    3+ test categories | full suite > 2 min | failures span multiple modules |
    both lint/typecheck AND test failures to fix
}

fn captureBaseline() {
  Run full test suite to establish baseline.
  Record passing, failing, skipped counts.
  Present baseline per reference/output-format.md.

  match (baseline) {
    all passing   => continue
    failures      => flag per Ownership Mandate — you still own these
  }
}

fn executeTests(mode) {
  match (mode) {
    Standard => run full suite, capture verbose output, parse results
    Team     => create team, spawn one runner per test category, assign tasks
  }

  Present execution results per reference/output-format.md.

  match (results) {
    all passing => skip to runQualityChecks
    failures    => proceed to fixFailures
  }
}

fn fixFailures(failures) {
  // Per reference/failure-investigation.md:
  // Categorize → fix → re-run specific test → re-run full suite → repeat
  // If fixing one breaks another: find root cause, don't give up

  for each (failure in failures) {
    Investigate per reference/failure-investigation.md categories.
    Apply minimal fix.
    Re-run to verify.
  }
}

fn runQualityChecks() {
  for each (command in qualityCommands) {
    run command
    match (result) {
      pass => continue
      fail => fix issues, re-run to verify
    }
  }

  Constraints {
    Same ownership rules: if it's broken in files you touched, fix it.
  }
}

fn report(results) {
  Present final report per reference/output-format.md.
}

test(target) {
  discover(target) |> selectMode |> captureBaseline |> executeTests |> fixFailures |> runQualityChecks |> report
}

## Integration with Other Skills

Called by other workflow skills:
- After `/start:implement` — verify implementation didn't break tests
- After `/start:refactor` — verify refactoring preserved behavior
- After `/start:debug` — verify fix resolved the issue without regressions
- Before `/start:review` — ensure clean test suite before review

When called by another skill, skip `discover()` if test infrastructure was already identified.
