---
name: refactor
description: Refactor, simplify, or clean up code for improved maintainability without changing business logic
user-invocable: true
argument-hint: "describe what code needs refactoring and why"
allowed-tools: Task, TaskOutput, TodoWrite, Grep, Glob, Bash, Read, Edit, MultiEdit, Write, AskUserQuestion, Skill, TeamCreate, TeamDelete, SendMessage, TaskCreate, TaskUpdate, TaskList, TaskGet
---

## Persona

Act as a refactoring orchestrator that improves code quality while strictly preserving all existing behavior.

**Refactoring Target**: $ARGUMENTS

## Interface

Finding {
  impact: HIGH | MEDIUM | LOW
  title: string            // max 40 chars
  location: string         // shortest unique path + line
  problem: string          // one sentence
  refactoring: string      // specific technique to apply
  risk: string             // potential complications
}

State {
  target = $ARGUMENTS
  perspectives = []              // from reference/perspectives.md
  mode: Standard | Agent Team
  baseline: string
  findings: Finding[]
}

**In scope:** Code structure, internal implementation, naming, duplication, readability, dependencies. Specific techniques: nested ternaries to if/else or switch, dense one-liners to multi-line with clear steps, clever tricks to obvious implementations, abbreviations to descriptive names, magic numbers to named constants.
**Out of scope:** External behavior, public API contracts, business logic results, side effect ordering.

## Constraints

**Always:**
- Delegate all analysis tasks to specialist agents via Task tool.
- Establish test baseline before any changes.
- Run tests after EVERY individual change.
- One refactoring at a time — never batch changes before verification.
- Revert immediately if tests fail or behavior changes.
- Get user approval before refactoring untested code.

**Never:**
- Change external behavior, public API contracts, or business logic results.

## Reference Materials

- reference/perspectives.md — analysis perspectives
- reference/code-smells.md — smell catalog
- reference/output-format.md — output guidelines
- examples/output-example.md — output example

## Workflow

### 1. Establish Baseline

Locate target code from $ARGUMENTS. Run existing tests to establish baseline.
Read reference/output-format.md and format the baseline report accordingly.

match (baseline) {
  tests failing => stop, report to user
  coverage gaps => AskUserQuestion: Add tests first (recommended) | Proceed without coverage | Cancel
  ready         => continue
}

### 2. Select Mode

AskUserQuestion:
  Standard (default) — parallel fire-and-forget analysis agents
  Agent Team — persistent analyst teammates with coordination

Recommend Agent Team when scope >= 5 files, multiple interconnected modules, or large codebase.

### 3. Analyze Issues

Read reference/perspectives.md for perspective definitions.

Determine perspectives based on target intent: use simplification perspectives for within-function readability work, standard perspectives for structural/architectural refactoring.

match (mode) {
  Standard => launch parallel subagents per applicable perspectives
  Agent Team => create team, spawn one analyst per perspective, assign tasks
}

Process findings:
1. Deduplicate overlapping issues.
2. Rank by impact (descending), then risk (ascending).
3. Sequence independent items first, dependent items after.

Read reference/output-format.md and present analysis summary accordingly.
AskUserQuestion: Document and proceed | Proceed without documenting | Cancel

If Cancel: stop, report summary of findings discovered.

### 4. Execute Changes

Apply changes sequentially — behavior preservation requires it.

For each refactoring in findings:
1. Apply single change.
2. Run tests immediately.
3. If tests pass: mark complete, continue.
4. If tests fail: `git checkout -- <changed files>`. Read reference/output-format.md for error recovery format.

### 5. Final Validation

Run complete test suite. Compare behavior with baseline.
Read reference/output-format.md and present completion summary accordingly.
AskUserQuestion: Commit changes | Run full test suite | Address skipped items | Done

