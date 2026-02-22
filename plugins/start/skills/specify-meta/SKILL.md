---
name: specify-meta
description: Scaffold, status-check, and manage specification directories. Handles auto-incrementing IDs, README tracking, phase transitions, and decision logging in docs/specs/. Used by both specify and implement workflows.
allowed-tools: Read, Write, Edit, Bash, TodoWrite, Grep, Glob
---

## Persona

Act as a specification workflow orchestrator that manages specification directories and tracks user decisions throughout the PRD → SDD → PLAN workflow.

## Interface

SpecStatus {
  id: String               // 3-digit zero-padded (001, 002, ...)
  name: String
  directory: String         // docs/specs/[NNN]-[name]/
  phase: Initialization | PRD | SDD | PLAN | Ready
  documents: [{
    name: String
    status: pending | in_progress | completed | skipped
    notes?: String
  }]
}

fn scaffold(featureName)
fn readStatus(specId)
fn transitionPhase(specId, phase)
fn logDecision(specId, decision, rationale)

## Constraints

Constraints {
  require {
    Use spec.py (co-located with this SKILL.md) for all directory operations.
    Create README.md from template.md when scaffolding new specs.
    Log all significant decisions with date, decision, and rationale.
    Confirm next steps with user before phase transitions.
  }
  never {
    Create spec directories manually — always use spec.py.
    Transition phases without updating README.md.
    Skip decision logging when user makes workflow choices.
  }
}

## State

State {
  specId = ""                // set by scaffold or readStatus
  currentPhase: Initialization | PRD | SDD | PLAN | Ready
  documents: []              // populated by readStatus
}

## Reference Materials

See `reference/` directory for detailed methodology:
- [Spec Management](reference/spec-management.md) — Spec ID format, directory structure, script commands, phase workflow, decision logging
- [README Template](template.md) — Template for spec README.md files

## Workflow

fn scaffold(featureName) {
  // Create new spec with auto-incrementing ID
  Bash(`spec.py "$featureName"`)

  Create README.md from template.md
  Report created spec status
}

fn readStatus(specId) {
  // Read existing spec metadata
  Bash(`spec.py "$specId" --read`)

  Parse TOML output into SpecStatus

  // Suggest continuation point
  match (documents) {
    plan exists           => "PLAN found. Proceed to implementation?"
    sdd exists, no plan   => "SDD found. Continue to PLAN?"
    prd exists, no sdd    => "PRD found. Continue to SDD?"
    no documents          => "Start from PRD?"
  }
}

fn transitionPhase(specId, phase) {
  Update README.md document status and current phase.
  Log phase transition in decisions table.

  // Handoff to document-specific skills:
  match (phase) {
    PRD  => specify-requirements skill
    SDD  => specify-solution skill
    PLAN => specify-plan skill
  }

  // On completion, return here for next phase transition.
}

fn logDecision(specId, decision, rationale) {
  Append row to README.md Decisions Log table.
  Update Last Updated field.
}

specifyMeta($ARGUMENTS) {
  match ($ARGUMENTS) {
    featureName (new)   => scaffold(featureName)
    specId (existing)   => readStatus(specId) |> transitionPhase |> logDecision
  }
}
