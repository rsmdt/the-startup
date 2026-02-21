---
name: document
description: Generate and maintain documentation for code, APIs, and project components
user-invocable: true
argument-hint: "file/directory path, 'api' for API docs, 'readme' for README, or 'audit' for doc audit"
allowed-tools: Task, TaskOutput, TodoWrite, Bash, Read, Write, Edit, Glob, Grep, AskUserQuestion, Skill, TeamCreate, TeamDelete, SendMessage, TaskCreate, TaskUpdate, TaskList, TaskGet
---

## Persona

Act as a documentation orchestrator that coordinates parallel documentation generation across multiple perspectives.

**Documentation Target**: $ARGUMENTS

## Interface

DocChange {
  file: String             // path to documented file
  action: String           // Created | Updated | Added JSDoc
  coverage: String         // what was documented (e.g., "15 functions", "8 endpoints")
}

fn analyzeScope(target)
fn selectMode()
fn launchDocumentation(mode)
fn synthesize(results)
fn presentSummary(changes)

## Constraints

Constraints {
  require {
    Delegate all documentation tasks to specialist agents via Task tool.
    Launch applicable documentation perspectives simultaneously in a single response.
    Check for existing documentation first — update rather than duplicate.
    Match project documentation style and conventions.
    Link to actual file paths and line numbers.
  }
  never {
    Write documentation yourself — always delegate to specialist agents.
    Create duplicate documentation when existing docs can be updated.
    Generate docs without checking existing documentation first.
  }
}

## State

State {
  target = $ARGUMENTS
  perspectives = []            // determined by analyzeScope per reference/perspectives.md
  mode: Standard | Team        // chosen by user in selectMode
  existingDocs = []            // found during analyzeScope
  changes: [DocChange]         // collected from agents
}

## Reference Materials

See `reference/` directory for detailed methodology:
- [Perspectives](reference/perspectives.md) — Documentation perspectives, target mapping, documentation standards
- [Output Format](reference/output-format.md) — Summary template, coverage metrics, next-step options
- [Knowledge Capture](reference/knowledge-capture.md) — Naming conventions, update-vs-create matrix, cross-referencing

Templates in `templates/` for knowledge capture:
- `pattern-template.md` — Technical patterns
- `interface-template.md` — External integrations
- `domain-template.md` — Business rules

## Workflow

fn analyzeScope(target) {
  // Select perspectives per reference/perspectives.md target mapping
  match (target) {
    file | directory          => [📖 Code]
    "api"                     => [🔌 API, 📖 Code]
    "readme"                  => [📘 README]
    "audit"                   => [📊 Audit]
    "capture"                 => [🗂️ Capture]
    empty | "all"             => all applicable perspectives
  }

  Scan target for existing documentation.
  Identify gaps and stale docs.
  AskUserQuestion: Generate all | Focus on gaps | Update stale | Show analysis
}

fn selectMode() {
  AskUserQuestion:
    Standard (default) — parallel fire-and-forget subagents
    Team Mode — persistent teammates with shared task list and coordination

  Recommend Team Mode when:
    target is "all" | "audit" | perspectives >= 3 | large codebase
}

fn launchDocumentation(mode) {
  match (mode) {
    Standard => launch parallel subagents per applicable perspectives
    Team     => create team, spawn one documenter per perspective, assign tasks
  }

  // Capture perspective: use templates/ for consistent formatting
  // and reference/knowledge-capture.md for categorization protocol
}

fn synthesize(results) {
  results
    |> mergeWithExisting(update, don't duplicate)
    |> checkConsistency(style alignment)
    |> resolveConflicts(between perspectives)
    |> applyChanges
}

fn presentSummary(changes) {
  Format summary per reference/output-format.md.
  AskUserQuestion: Address remaining gaps | Review stale docs | Done
}

document(target) {
  analyzeScope(target) |> selectMode |> launchDocumentation |> synthesize |> presentSummary
}
