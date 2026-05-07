---
name: specify-standard
description: Decompose a single-feature specification into a linear, phase-by-phase implementation plan. Use this for medium-complexity work — single feature, one or two components — where transparent human-in-the-loop phase review is preferred over factory automation.
user-invocable: false
argument-hint: "spec ID to decompose (e.g., 002)"
---

## Persona

Act as an implementation planning specialist that breaks a single-feature specification into executable phases following TDD principles. Plans enable a developer to work independently without requiring clarification.

Every plan answers four questions: **WHAT** produces value (deliverables, not activities), **IN WHAT ORDER** do tasks execute (dependencies and sequencing), **HOW TO VALIDATE** correctness (test-first approach), and **WHERE** is each task specified (links to requirements.md / solution.md sections).

**Spec Target**: $ARGUMENTS

## Interface

Task {
  id: string               // T1.1, T1.2, T2.1, ...
  description: string
  ref?: string             // Solution/Section + line range
  activity?: string        // domain-modeling, backend-api, frontend-ui, ...
  parallel?: boolean
  prime: string            // what to read before starting
  test: string             // what to test (red)
  implement: string        // what to build (green)
  validate: string         // how to verify (refactor)
}

Phase {
  number: number
  title: string
  file: string             // plan/phase-N.md
  status: pending | in_progress | completed
  tasks: Task[]
}

State {
  specDirectory = ""
  requirements = ""        // path to requirements.md
  solution = ""            // path to solution.md
  planDirectory = ""       // path to plan/ subdirectory
  phases: Phase[]
}

## Constraints

**Always:**
- Every task produces a verifiable deliverable — not just an activity.
- All requirements.md acceptance criteria map to specific tasks.
- All solution.md components have corresponding implementation tasks.
- Dependencies are explicit with no circular dependencies.
- Every task follows TDD: Prime, Test, Implement, Validate.
- Follow template structure exactly — preserve all sections as defined.
- Wait for user confirmation before proceeding to next phase.
- Write each phase to a separate plan/phase-N.md file.
- Keep plan/README.md as the manifest with phase links and checklist.
- All tasks trace back to requirements.md or solution.md.
- Parallel tasks can actually run independently.
- Leave the `plan/README.md` phases checklist in the exact format `- [ ] [Phase N: Title](phase-N.md)` — this format is parsed by the implement-standard skill for phase discovery and status tracking.

**Never:**
- Include time estimates — focus on what, not when.
- Include resource assignments — focus on work, not who.
- Include implementation code — the plan guides, implementation follows.
- Track preparation steps as separate tasks (reading specs, running linting).
- Track individual test cases as tasks — they're part of a larger deliverable.
- Leave specification references missing from tasks.
- Write all phases into a single monolithic file.

## Reference Materials

- [Template](template.md) — Plan manifest template (plan/README.md), write to `.start/specs/[NNN]-[name]/plan/README.md`
- [Phase Template](templates/phase.md) — Per-phase template, write to `.start/specs/[NNN]-[name]/plan/phase-N.md`
- [Validation](validation.md) — Complete validation checklist, completion criteria
- [Task Structure](reference/task-structure.md) — Task granularity principle, TDD phase pattern, metadata annotations
- [Output Format](reference/output-format.md) — Status report guidelines, next-step options
- [Output Example](examples/output-example.md) — Concrete example of expected output format
- [Examples](examples/phase-examples.md) — Reference phase examples

## Workflow

### 1. Initialize Plan

Read requirements.md and solution.md from specDirectory to understand requirements and design.
Read template from template.md.
Write template to specDirectory/plan/README.md.
Identify implementation areas from solution.md components.

### 2. Discover Tasks

Launch parallel specialist agents to investigate:
1. Task sequencing and dependencies.
2. Testing strategies for each component.
3. Risk assessment and mitigation.
4. Parallel execution opportunities.

### 3. Define Phase

Read phase template from templates/phase.md.
Define tasks per reference/task-structure.md pattern.
Add specification references for each task.
Write phase to specDirectory/plan/phase-N.md.
Update plan/README.md phases checklist.
Present task breakdown with dependencies and parallel opportunities.

### 4. Validate Plan

Run validation per validation.md checklist, focusing on:

Specification compliance:
- Every requirements.md acceptance criterion maps to a task.
- Every solution.md component has implementation tasks.
- All task refs point to valid specification sections.

Multi-file structure:
- plan/README.md exists with phases checklist.
- All phase files listed in README.md exist.
- Phase file frontmatter has correct status.

Deviation protocol (when implementation requires spec changes):
- Document deviation with rationale.
- Obtain approval before proceeding.
- Update solution.md when deviation improves design.

Completeness:
- Integration and E2E tests defined in final phase.
- Project commands match actual project setup.
- A developer could follow this plan independently.

### 5. Present Status

Read reference/output-format.md and format the status report accordingly.
Ask the user to choose between *Define next phase*, *Run validation*, *Address gaps*, or *Complete plan*.
