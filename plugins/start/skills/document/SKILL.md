---
name: document
description: Generate and maintain documentation for code, APIs, and project components
user-invocable: true
argument-hint: "file/directory path, 'api' for API docs, 'readme' for README, or 'audit' for doc audit"
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
  mode: Standard | Agent Team | Dynamic Workflow
  existingDocs = []
  changes: DocChange[]
}

## Constraints

**Always:**
- Launch applicable documentation perspectives simultaneously in a single response.
- Check for existing documentation first — update rather than duplicate.
- Match project documentation style and conventions.
- Link to actual file paths and line numbers.

**Never:**
- Write documentation yourself — always delegate to specialist agents.

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

Ask the user to choose: Generate all | Focus on gaps | Update stale | Show analysis

### 2. Select Mode

Ask the user to choose:
  Standard (default) — parallel fire-and-forget subagents
  Agent Team — persistent teammates with shared task list and coordination
  Dynamic Workflow — one inline workflow script fans out the documenters and collects their changes

Recommend Agent Team when target is "all" or "audit", perspectives >= 3, or large codebase.
Recommend Dynamic Workflow instead when those thresholds apply but the perspectives are independent and need no coordination between documenters.

### 3. Launch Documentation

match (mode) {
  Standard => launch parallel subagents per applicable perspectives
  Agent Team => create team, spawn one documenter per perspective, assign tasks
  Dynamic Workflow => run one dynamic workflow: one documenter agent per perspective in parallel, returning every DocChange
}

For Dynamic Workflow, these instructions are the opt-in: call the Workflow tool directly with an inline script. Do not load the workflow-authoring reference or assemble the script through shell commands first — runs that detour through either stall before launching. Choose `model` and `effort` per stage from what the stage does. Mechanical stages — running a command, gathering or quoting sources, mapping text into a fixed schema — can run on a smaller, faster model at lower effort. Analysis, writing, and judging another agent's work need the strongest model. A judging stage never runs on a weaker model than the stage it judges. When unsure, omit `model` so the stage inherits the session model.

For the Capture perspective: use templates/ for consistent formatting and Read reference/knowledge-capture.md for categorization protocol.

### 4. Synthesize Results

Process results:
1. Merge with existing docs — update, don't duplicate.
2. Check consistency for style alignment.
3. Resolve conflicts between perspectives.
4. Apply changes.

### 5. Present Summary

Read reference/output-format.md and format summary accordingly.

Ask the user to choose: Address remaining gaps | Review stale docs | Done

