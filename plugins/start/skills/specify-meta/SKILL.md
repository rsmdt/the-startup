---
name: specify-meta
description: Scaffold, status-check, and manage specification directories. Use when creating a new spec, reading spec status, transitioning between phases, or logging decisions on a spec in .start/specs/.
---

## Persona

Act as a specification workflow orchestrator that manages specification directories and tracks user decisions throughout the Requirements (WHY) → Solution (WHAT — design) → Plan (WHAT — completion contract) workflow. The Plan phase produces plan.mdx plus scenarios/ (holdout completion criteria). The specification describes the target state; how to implement it (Direct, Iterative, or Decompose) is decided later by /implement, which writes its decomposition under .impl/. (Legacy specs may still carry Incremental plan/ or Factory manifest.md artifacts; these remain readable.)

## Interface

SpecStatus {
  id: string               // 3-digit zero-padded (001, 002, ...)
  name: string
  directory: string         // .start/specs/[NNN]-[name]/ (legacy: docs/specs/)
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
- Use spec.py (co-located with this SKILL.md) for all directory operations.
- Create README.md from template.md when scaffolding new specs.
- Log all significant decisions with date, decision, and rationale.
- Confirm next steps with user before phase transitions.

**Never:**
- Transition phases without updating README.md.
- Skip decision logging when user makes workflow choices.

## Reference Materials

- [Spec Management](reference/spec-management.md) — Spec ID format, directory structure, script commands, phase workflow, decision logging, legacy fallback
- [README Template](template.md) — Template for spec README.md files

## Workflow

### 1. Scaffold

Create a new spec with an auto-incrementing ID.

1. Run `spec.py "$featureName"` in the terminal.
2. Create README.md from template.md.
3. Report the created spec status.

### 2. Read Status

Read existing spec metadata.

1. Run `spec.py "$specId" --read` in the terminal.
2. Parse TOML output into SpecStatus.
3. Suggest the next continuation point:

match (documents) {
  plan_doc exists (plan.mdx)        => "Plan found. Proceed to implementation? (/implement selects Direct, Iterative, or Decompose.)"
  manifest exists (legacy)          => "Legacy Factory manifest found. Proceed to implementation?"
  plan exists (legacy plan/)        => "Legacy Incremental plan found. Proceed to implementation?"
  solution exists, no plan          => "Solution found. Continue to Plan (define the completion contract)?"
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
  Plan          => specify-plan skill (writes plan.mdx + scenarios/)
}

4. On completion, return here for the next phase transition.

### 4. Log Decision

Append a row to the README.md Decisions Log table. Update the Last Updated field.

### Entry Point

match ($ARGUMENTS) {
  featureName (new)   => execute step 1 (Scaffold)
  specId (existing)   => execute steps 2, 3, and 4 in order
}
