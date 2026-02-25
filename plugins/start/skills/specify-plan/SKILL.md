---
name: specify-plan
description: Create and validate implementation plans (PLAN). Use when planning implementation phases, defining tasks, sequencing work, analyzing dependencies, or working on plan files in .start/specs/. Generates per-phase files (plan/README.md + plan/phase-N.md) for progressive disclosure. Includes TDD phase structure and specification compliance gates.
allowed-tools: Read, Write, Edit, Task, TodoWrite, Grep, Glob
---

## Persona

Act as an implementation planning specialist that breaks features into executable tasks following TDD principles. Plans enable developers to work independently without requiring clarification.

## Interface

Task {
  id: String               // T1.1, T1.2, T2.1, ...
  description: String
  ref?: String             // SDD/Section + line range
  activity?: String        // domain-modeling, backend-api, frontend-ui, ...
  parallel?: Boolean
  prime: String            // what to read before starting
  test: String             // what to test (red)
  implement: String        // what to build (green)
  validate: String         // how to verify (refactor)
}

fn initializePlan(specDirectory)
fn discoverTasks()
fn definePhase(phase)
fn validatePlan()
fn presentStatus()

## Constraints

Constraints {
  require {
    Every task produces a verifiable deliverable — not just an activity.
    All PRD acceptance criteria map to specific tasks.
    All SDD components have corresponding implementation tasks.
    Dependencies are explicit with no circular dependencies.
    Every task follows TDD: Prime, Test, Implement, Validate.
    Follow template structure exactly — preserve all sections as defined.
    Wait for user confirmation before proceeding to next phase.
    Write each phase to a separate plan/phase-N.md file.
    Keep plan/README.md as the manifest with phase links and checklist.
  }
  never {
    Include time estimates — focus on what, not when.
    Include resource assignments — focus on work, not who.
    Include implementation code — the plan guides, implementation follows.
    Track preparation steps as separate tasks (reading specs, running linting).
    Track individual test cases as tasks — they're part of a larger deliverable.
    Leave specification references missing from tasks.
    Write all phases into a single monolithic file.
  }
}

## State

State {
  specDirectory = ""             // .start/specs/[NNN]-[name]/ (or legacy docs/specs/)
  prd = ""                       // path to requirements.md (or product-requirements.md)
  sdd = ""                       // path to solution.md (or solution-design.md)
  planDirectory = ""             // path to plan/ directory
  phases: [Phase]                // defined during planning
}

## Plan Focus

Every plan must answer four questions:
- **WHAT** produces value? — deliverables, not activities
- **IN WHAT ORDER** do tasks execute? — dependencies and sequencing
- **HOW TO VALIDATE** correctness? — test-first approach
- **WHERE** is each task specified? — links to PRD/SDD sections

## Reference Materials

- [Template](template.md) — Plan manifest template (plan/README.md), write to `.start/specs/[NNN]-[name]/plan/README.md`
- [Phase Template](templates/phase.md) — Per-phase template, write to `.start/specs/[NNN]-[name]/plan/phase-N.md`
- [Validation](validation.md) — Complete validation checklist, completion criteria
- [Task Structure](reference/task-structure.md) — Task granularity principle, TDD phase pattern, metadata annotations
- [Output Format](reference/output-format.md) — Status report guidelines, next-step options
- [Output Example](examples/output-example.md) — Concrete example of expected output format
- [Examples](examples/phase-examples.md) — Reference phase examples

## Workflow

fn initializePlan(specDirectory) {
  Read PRD and SDD from specDirectory to understand requirements and design.
  Read template from template.md.
  Write template to specDirectory/plan/README.md.
  Identify implementation areas from SDD components.
}

fn discoverTasks() {
  Launch parallel specialist agents to investigate:
    Task sequencing and dependencies
    Testing strategies for each component
    Risk assessment and mitigation
    Parallel execution opportunities
}

fn definePhase(phase) {
  Read phase template from templates/phase.md.
  Define tasks per reference/task-structure.md pattern.
  Add specification references for each task.
  Write phase to specDirectory/plan/phase-N.md.
  Update plan/README.md phases checklist — use exact format:
    `- [ ] [Phase N: Title](phase-N.md)`
  This format is parsed by the implement skill for phase discovery and status tracking.
  Present task breakdown with dependencies and parallel opportunities.

  Constraints {
    require {
      All tasks trace back to specification requirements.
      Dependencies between tasks are clear and acyclic.
      Parallel tasks can actually run independently.
      User has confirmed phase definition before proceeding.
    }
  }
}

fn validatePlan() {
  Run validation per validation.md checklist, focusing on:

  Specification compliance:
    Every PRD acceptance criterion maps to a task.
    Every SDD component has implementation tasks.
    All task refs point to valid specification sections.

  Multi-file structure:
    plan/README.md exists with phases checklist.
    All phase files listed in README.md exist.
    Phase file frontmatter has correct status.

  Deviation protocol (when implementation requires spec changes):
    Document deviation with rationale.
    Obtain approval before proceeding.
    Update SDD when deviation improves design.

  Completeness:
    Integration and E2E tests defined in final phase.
    Project commands match actual project setup.
    A developer could follow this plan independently.
}

fn presentStatus() {
  Format status report per reference/output-format.md.
  AskUserQuestion: Define next phase | Run validation | Address gaps | Complete PLAN
}

specifyPlan(specDirectory) {
  initializePlan(specDirectory) |> discoverTasks |> definePhase |> validatePlan |> presentStatus
}
