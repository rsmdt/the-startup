---
name: specify-meta
user-invocable: false
description: Scaffold, status-check, and manage specification directories. Use when creating a new spec, reading spec status, transitioning between phases, or logging decisions on a spec in .start/specs/.
---

## Persona

Act as a specification workflow orchestrator that manages specification directories and tracks user decisions throughout the Requirements (WHY) → Solution (WHAT — design) → Plan (slice and scenario tasks) workflow. The Plan phase produces a narrative plan.md plus the slice and scenario tasks in the spec's tracker. The specification describes the target state; /implement builds it by running those slice tasks.

## Interface

SpecStatus {
  id: string               // 3-digit zero-padded (001, 002, ...)
  name: string
  directory: string         // .start/specs/[NNN]-[name]/
  phase: Initialization | Requirements | Solution | Plan | Ready
  documents: {
    name: string
    status: pending | in_progress | completed | skipped
    notes?: string
  }[]
}

State {
  specId = ""
  currentPhase: Initialization | Requirements | Solution | Plan | Ready
  documents: []
}

## Constraints

**Always:**
- Use `node ${CLAUDE_PLUGIN_ROOT}/scripts/spec.js` for all directory operations.
- Create README.md from template.md when scaffolding new specs.
- Log all significant decisions with date, decision, and rationale.
- Keep the README Assumption Ledger current — status changes (Proposed → Confirmed/Corrected/Standing) are updates to this file, made when the specify skill reports them.
- Confirm next steps with user before phase transitions.

**Never:**
- Transition phases without updating README.md.
- Skip decision logging when user makes workflow choices.

## Reference Materials

- [Spec Management](reference/spec-management.md) — Spec ID format, directory structure, script commands, phase workflow, decision logging
- [README Template](template.md) — Template for spec README.md files

## Workflow

### 1. Scaffold

Create a new spec with an auto-incrementing ID.

1. Run `node ${CLAUDE_PLUGIN_ROOT}/scripts/spec.js "$featureName"` in the terminal.
2. Create README.md from template.md.
3. Report the created spec status.

### 2. Read Status

Read existing spec metadata.

1. Run `node ${CLAUDE_PLUGIN_ROOT}/scripts/spec.js "$specId" --read` in the terminal.
2. Parse TOML output into SpecStatus.
3. Suggest the next continuation point:

match (documents) {
  plan_doc exists (plan.md)        => "Plan found. Proceed to implementation with /implement?"
  solution exists, no plan          => "Solution found. Continue to Plan (slice and scenario tasks)?"
  requirements exists, no solution  => "Requirements found. Continue to Solution?"
  no documents                      => "Start from Requirements?"
}

### 3. Transition Phase

Update the spec directory to reflect the new phase.

1. Update README.md document status and current phase.
2. Log the phase transition in the decisions table.
3. Hand off to the document-specific skill:

match (phase) {
  Requirements  => specify-requirements skill
  Solution      => specify-solution skill
  Plan          => specify-plan skill (writes plan.md + tracker tasks)
}

4. On completion, return here for the next phase transition.

### 4. Log Decision

Append a row to the README.md Decisions Log table. Update the Last Updated field.

### Entry Point

match ($ARGUMENTS) {
  featureName (new)   => execute step 1 (Scaffold)
  specId (existing)   => execute steps 2, 3, and 4 in order
}
