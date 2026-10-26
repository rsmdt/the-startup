---
name: specify-meta
user-invocable: false
description: Scaffold, status-check, and manage specification directories. Use when creating a new spec, reading spec status, transitioning between phases, or logging decisions on a spec in .start/specs/.
---

## Persona

Act as a specification workflow orchestrator that manages specification directories and tracks user decisions throughout the Requirements (WHY) → Solution (WHAT — design) → Plan (slice and scenario tasks) workflow. Every phase writes its own part of one document, `spec.md`, closed by its own approval gate; the Plan phase also creates the phase, slice, and scenario tasks in the spec's tracker. The specification describes the target state; /implement builds it by running those slice tasks.

## Interface

SpecStatus {
  id: string               // 3-digit zero-padded (001, 002, ...)
  name: string
  directory: string         // .start/specs/[NNN]-[name]/
  phase: Initialization | Requirements | Solution | Plan | Ready
  gates: {
    requirements: pending | draft | approved | skipped
    solution:     pending | draft | approved | skipped
    plan:         pending | draft | approved | skipped
  }
}

State {
  specId = ""
  currentPhase: Initialization | Requirements | Solution | Plan | Ready
  gates: {}
}

## Constraints

**Always:**
- Use `node ${CLAUDE_PLUGIN_ROOT}/scripts/spec.js` for all directory operations.
- Keep phase state in the `gates` frontmatter of `spec.md` and the decisions and assumptions in its `## Ledger` part — `spec.md` has no README beside it.
- Log all significant decisions with date, decision, and rationale.
- Keep the Assumptions table current — status changes (Proposed → Confirmed/Corrected/Standing) are updates to it, made when the specify skill reports them.
- A spec with no `spec.md` uses the legacy layout (`README.md`, `requirements.md`, `solution.md`, `plan.md`): read and continue it in those files.
- Confirm next steps with user before phase transitions.

**Never:**
- Transition phases without updating the `gates` frontmatter.
- Skip decision logging when user makes workflow choices.

## Reference Materials

- [Spec Management](reference/spec-management.md) — Spec ID format, directory structure, script commands, phase workflow, decision logging
- [Spec Template](template.md) — `spec.md`: the Requirements, Solution, and Plan parts plus the Ledger

## Workflow

### 1. Scaffold

Create a new spec with an auto-incrementing ID.

1. Run `node ${CLAUDE_PLUGIN_ROOT}/scripts/spec.js "$featureName"` in the terminal. It creates the spec directory, `spec.md` from template.md, and `.tracker/`.
2. Report the created spec status.

### 2. Read Status

Read existing spec metadata.

1. Run `node ${CLAUDE_PLUGIN_ROOT}/scripts/spec.js "$specId" --read` in the terminal.
2. Parse TOML output into SpecStatus.
3. Read the `gates` frontmatter of `spec.md` and suggest the next continuation point:

match (gates) {
  plan approved                     => "Plan approved. Proceed to implementation with /implement?"
  solution approved, plan not       => "Solution approved. Continue to Plan (phases, slice and scenario tasks)?"
  requirements approved, solution not => "Requirements approved. Continue to Solution?"
  nothing approved                  => "Start from Requirements?"
}

### 3. Transition Phase

Update the spec directory to reflect the new phase.

1. Set the finished part's gate to `approved` (or `skipped`) and the next part's gate to `draft`.
2. Log the phase transition in the Decision log.
3. Hand off to the document-specific skill:

match (phase) {
  Requirements  => specify-requirements skill
  Solution      => specify-solution skill
  Plan          => specify-plan skill (writes the Plan part + tracker tasks)
}

4. On completion, return here for the next phase transition.

### 4. Log Decision

Append a row to the Decision log table in the Ledger part of `spec.md`.

### Entry Point

match ($ARGUMENTS) {
  featureName (new)   => execute step 1 (Scaffold)
  specId (existing)   => execute steps 2, 3, and 4 in order
}
