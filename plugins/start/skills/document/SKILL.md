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
  file: string             // path to documented file
  action: string           // Created | Updated | Added JSDoc
  coverage: string         // what was documented (e.g., "15 functions", "8 endpoints")
}

State {
  target = $ARGUMENTS
  perspectives = []        // from reference/perspectives.md
  mode: Standard | Agent Team
  existingDocs = []
  changes: DocChange[]
}

## Constraints

**Always:**
- Delegate all documentation tasks to specialist agents via Task tool.
- Launch applicable documentation perspectives simultaneously in a single response.
- Check for existing documentation first — update rather than duplicate.
- Match project documentation style and conventions.
- Link to actual file paths and line numbers.

**Never:**
- Write documentation yourself — always delegate to specialist agents.
- Create duplicate documentation when existing docs can be updated.
- Generate docs without checking existing documentation first.

## Reference Materials

- reference/perspectives.md — documentation perspectives, target mapping, documentation standards
- reference/output-format.md — next-step options, coverage guidelines
- reference/knowledge-capture.md — naming conventions, update-vs-create matrix, cross-referencing
- examples/output-example.md — concrete example of expected output format

Templates in `templates/` for knowledge capture:
- `pattern-template.md` — Technical patterns
- `interface-template.md` — External integrations
- `domain-template.md` — Business rules

## Workflow

### 1. Analyze Scope

Read reference/perspectives.md. Select perspectives based on target:

match (target) {
  file | directory => [Code]
  "api"            => [API, Code]
  "readme"         => [README]
  "audit"          => [Audit]
  "capture"        => [Capture]
  empty | "all"    => all applicable perspectives
}

Scan target for existing documentation. Identify gaps and stale docs.

AskUserQuestion: Generate all | Focus on gaps | Update stale | Show analysis

### 2. Select Mode

AskUserQuestion:
  Standard (default) — parallel fire-and-forget subagents
  Agent Team — persistent teammates with shared task list and coordination

Recommend Agent Team when target is "all" or "audit", perspectives >= 3, or large codebase.

### 3. Launch Documentation

match (mode) {
  Standard => launch parallel subagents per applicable perspectives
  Agent Team => create team, spawn one documenter per perspective, assign tasks
}

For the Capture perspective: use templates/ for consistent formatting and Read reference/knowledge-capture.md for categorization protocol.

### 4. Synthesize Results

Process results:
1. Merge with existing docs — update, don't duplicate.
2. Check consistency for style alignment.
3. Resolve conflicts between perspectives.
4. Apply changes.

### 5. Present Summary

Read reference/output-format.md and format summary accordingly.

AskUserQuestion: Address remaining gaps | Review stale docs | Done

