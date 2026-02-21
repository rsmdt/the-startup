---
name: refactor
description: Refactor code for improved maintainability without changing business logic
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
  title: String            // max 40 chars
  location: String         // shortest unique path + line
  problem: String          // one sentence
  refactoring: String      // specific technique to apply
  risk: String             // potential complications
}

fn establishBaseline(target)
fn selectMode()
fn analyzeIssues(mode)
fn executeChanges(plan)
fn finalValidation(results)

## Constraints

Constraints {
  require {
    Delegate all analysis tasks to specialist agents via Task tool.
    Establish test baseline before any changes.
    Run tests after EVERY individual change.
    One refactoring at a time — never batch changes before verification.
    Revert immediately if tests fail or behavior changes.
    Get user approval before refactoring untested code.
  }
  never {
    Refactor code without a passing test baseline.
    Batch multiple refactorings before running tests.
    Skip test verification after applying a change.
    Refactor untested code without explicit user approval.
    Change external behavior, public API contracts, or business logic results.
  }
}

RefactoringBoundaries {
  prefer {
    Code structure, internal implementation details.
    Variable and function names.
    Duplication removal.
    Dependencies and versions.
    Nested ternaries => `if/else` or `switch`.
    Dense one-liners => multi-line with clear steps.
    Clever tricks => obvious implementations.
    Abbreviations => descriptive names.
    Magic numbers => named constants.
  }
  avoid {
    External behavior.
    Public API contracts.
    Business logic results.
    Side effect ordering.
  }
}

## State

State {
  target = $ARGUMENTS
  perspectives = []              // determined by target type per reference/perspectives.md
  mode: Standard | Team          // chosen by user in selectMode
  baseline: String               // test status before changes
  findings: [Finding]            // collected from analysis agents
}

## Reference Materials

See `reference/` directory for detailed methodology:
- [Perspectives](reference/perspectives.md) — Standard and simplification perspectives, focus area mapping
- [Code Smells](reference/code-smells.md) — Code smell patterns, complexity indicators, refactoring opportunities
- [Output Format](reference/output-format.md) — Baseline, analysis, error recovery, completion guidelines
- [Output Example](examples/output-example.md) — Concrete example of expected output format

## Workflow

fn establishBaseline(target) {
  Locate target code from $ARGUMENTS.
  Run existing tests to establish baseline.
  Report baseline per reference/output-format.md.

  match (baseline) {
    tests failing    => stop, report to user
    coverage gaps    => flag for user decision before proceeding
    ready            => continue
  }
}

fn selectMode() {
  AskUserQuestion:
    Standard (default) — parallel fire-and-forget analysis agents
    Team Mode — persistent analyst teammates with coordination

  Recommend Team Mode when:
    scope >= 5 files | multiple interconnected modules | large codebase
}

fn analyzeIssues(mode) {
  perspectives = match (target) {
    simplification keywords => simplification perspectives per reference/perspectives.md
    default                 => standard perspectives per reference/perspectives.md
  }

  match (mode) {
    Standard => launch parallel subagents per applicable perspectives
    Team     => create team, spawn one analyst per perspective, assign tasks
  }

  findings
    |> deduplicate(overlapping issues)
    |> rank(by: [impact desc, risk asc])
    |> sequence(independent first, dependent after)

  Present analysis summary per reference/output-format.md.
  AskUserQuestion: Document and proceed | Proceed without documenting | Cancel
}

fn executeChanges(plan) {
  // Always sequential regardless of mode — behavior preservation requires it
  for each (refactoring in plan) {
    1. Apply single change
    2. Run tests immediately
    match (tests) {
      pass => mark complete, continue
      fail => revert, present error recovery per reference/output-format.md
    }
  }
}

fn finalValidation(results) {
  Run complete test suite.
  Compare behavior with baseline.
  Present completion summary per reference/output-format.md.
  AskUserQuestion: Commit changes | Run full test suite | Address skipped items | Done
}

refactor(target) {
  establishBaseline(target) |> selectMode |> analyzeIssues |> executeChanges |> finalValidation
}
