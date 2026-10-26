# Code Agent Prompt

Fill `{specId}`, `{taskScript}`, `{maxIterations}`, and the slice list — one slice id, or every slice the
agent builds, in build order. Per slice, include the **Prior attempts** block only when its history is
non-empty, and the **Retry** block, with that slice's `{attempt}` number, only when its failures are
non-empty.

---

You are implementing these slices of a specification, in this order, using Test-Driven Development (TDD):
{one line per slice id}

Each slice task is a code brief (features / acceptance criteria, EARS behavior, quality gates). For each
slice, in order: read its brief, mark it in progress, then build it before moving on to the next:
  node {taskScript} read <sliceId> --spec {specId}
  node {taskScript} update <sliceId> --spec {specId} --status building --override

Codebase context: read AGENTS.md (or CLAUDE.md) for project orientation and conventions.

DO NOT read the scenario child tasks (type: scenario) of any of these slices. They hold the holdout
acceptance evidence and must not influence how you build.

## State on disk — start from the delta, not from zero
This working tree may already contain work for these slices. Inventory it before you plan anything:
  git status --short
  git diff --stat HEAD
Treat whatever touches a brief as your starting point. If a slice is already built, do not rewrite it:
run the test suite, locate the specific gap named in its symptoms below, fix only that. Report
`wroteCode: false` only if you changed nothing for any slice.

## File naming
Name test files after the feature or component under test — never after a task id.

## TDD process
For each acceptance criterion and EARS behavior in each brief:
1. RED: write a failing test for the behavior. Run the suite — confirm it fails.
2. GREEN: write the minimal code to make it pass. Run the suite — confirm it passes.
3. REFACTOR: improve the code while keeping every test green.

Test behaviour: domain rules, API responses, and the UI state a user action produces. Do not write tests
for how the UI looks — layout, sizes, spacing, colours, viewport widths. They pin the styling and slow the
suite.

Run the unit and integration tests in each red-green step. Write end-to-end tests only when a brief asks
for them; run that suite once, after every slice is built, in the foreground with a timeout, and fix what
fails before you return.

## Slice {sliceId} — prior attempts, do NOT repeat them
{for each history entry: "- attempt {n} (business {b}% / technical {t}%): {approach} — symptoms that survived: {failures}"}

Every approach above has been tried and did not clear the acceptance evidence. Form a different
hypothesis. If a symptom survived two or more distinct fixes, re-derive the root cause from observed
behavior instead of patching the same code path. If the failure is not in code these slices own
(environment, a missing migration, an unowned module), say exactly that in `approach`.

## Slice {sliceId} — retry
Attempt {attempt} of {maxIterations} did not satisfy the acceptance evidence. Observed symptoms:
{one line per failure}
Investigate the codebase to find and fix the root causes.

## Done when
1. All existing tests pass.
2. New tests cover every acceptance criterion and EARS behavior in every brief.
3. The full test suite is green.

## Return
- approach: one or two lines per slice — the root cause you diagnosed and the change you made. Name the
  hypothesis, not the file list.
- wroteCode: false only if every slice was already implemented and you changed nothing.
- wroteMigration: true if you added or changed a database migration or schema file. Say true if unsure.
