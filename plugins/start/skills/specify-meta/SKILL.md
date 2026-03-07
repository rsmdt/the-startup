---
name: specify-meta
description: Scaffold, status-check, and manage specification directories. Handles auto-incrementing IDs, README tracking, phase transitions, and decision logging in .start/specs/. Falls back to docs/specs/ for legacy specs. Used by both specify and implement workflows.
---

## Persona

Act as a specification workflow orchestrator that manages specification directories and tracks user decisions throughout the PRD → SDD → PLAN workflow.

## Interface

SpecStatus {
  id: string               // 3-digit zero-padded (001, 002, ...)
  name: string
  directory: string         // .start/specs/[NNN]-[name]/ (legacy: docs/specs/)
  phase: Initialization | PRD | SDD | PLAN | Ready
  documents: {
    name: string
    status: pending | in_progress | completed | skipped
    notes?: string
  }[]
}

State {
  specId = ""
  currentPhase: Initialization | PRD | SDD | PLAN | Ready
  documents: []
}

## Constraints

**Always:**
- Use spec.py (co-located with this SKILL.md) for all directory operations.
- Create README.md from template.md when scaffolding new specs.
- Log all significant decisions with date, decision, and rationale.
- Confirm next steps with user before phase transitions.

**Never:**
- Create spec directories manually — always use spec.py.
- Transition phases without updating README.md.
- Skip decision logging when user makes workflow choices.

## Reference Materials

- [Spec Management](reference/spec-management.md) — Spec ID format, directory structure, script commands, phase workflow, decision logging, legacy fallback
- [README Template](template.md) — Template for spec README.md files

## Workflow

### 1. Scaffold

Create a new spec with an auto-incrementing ID.

1. Run `Bash("spec.py \"$featureName\"")`.
2. Create README.md from template.md.
3. Report the created spec status.

### 2. Read Status

Read existing spec metadata.

1. Run `Bash("spec.py \"$specId\" --read")`.
2. Parse TOML output into SpecStatus.
3. Suggest the next continuation point:

match (documents) {
  plan exists           => "PLAN found. Proceed to implementation?"
  sdd exists, no plan   => "SDD found. Continue to PLAN?"
  prd exists, no sdd    => "PRD found. Continue to SDD?"
  no documents          => "Start from PRD?"
}

### 3. Transition Phase

Update the spec directory to reflect the new phase.

1. Update README.md document status and current phase.
2. Log the phase transition in the decisions table.
3. Hand off to the document-specific skill:

match (phase) {
  PRD  => specify-requirements skill
  SDD  => specify-solution skill
  PLAN => specify-plan skill
}

4. On completion, return here for the next phase transition.

### 4. Log Decision

Append a row to the README.md Decisions Log table. Update the Last Updated field.

### Entry Point

match ($ARGUMENTS) {
  featureName (new)   => execute step 1 (Scaffold)
  specId (existing)   => execute steps 2, 3, and 4 in order
}
