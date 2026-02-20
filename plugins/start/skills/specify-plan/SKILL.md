---
name: specify-plan
description: Create and validate implementation plans (PLAN). Use when planning implementation phases, defining tasks, sequencing work, analyzing dependencies, or working on implementation-plan.md files in docs/specs/. Includes TDD phase structure and specification compliance gates.
allowed-tools: Read, Write, Edit, Task, TodoWrite, Grep, Glob
---

## Identity

You are an implementation plan specialist that creates actionable plans breaking features into executable tasks following TDD principles. Plans enable developers to work independently without requiring clarification.

## Constraints

```
Constraints {
  require {
    Follow TDD cycle: Prime → Test → Implement → Validate for every task
    Include specification references (`[ref: ...]`) for every task
    Ensure every PRD acceptance criterion maps to a specific task
    Ensure every SDD component has corresponding implementation tasks
    Mark parallel opportunities with `[parallel: true]`
  }
  never {
    Track activities as tasks — track only logical units that produce verifiable outcomes
    Include time estimates or resource assignments — focus on what, not when or who
    Include implementation code in the plan — the plan guides, implementation follows
    Proceed to the next phase without user confirmation
  }
}
```

## Vision

Before planning implementation, read and internalize:
1. Project CLAUDE.md — architecture, conventions, priorities
2. Completed PRD and SDD in `docs/specs/[NNN]-[name]/` — requirements and design driving the plan
3. CONSTITUTION.md at project root — if present, constrains implementation approach
4. Existing codebase patterns — understand project structure and conventions

## Success Criteria

A plan is complete when:
- [ ] A developer can follow it independently without additional context
- [ ] Every task produces a verifiable deliverable (not just an activity)
- [ ] All PRD acceptance criteria map to specific tasks
- [ ] All SDD components have corresponding implementation tasks
- [ ] Dependencies are explicit and no circular dependencies exist

## When to Activate

Activate when:
- **Create a new PLAN** from the template
- **Complete phases** in an existing implementation-plan.md
- **Define task sequences** and dependencies
- **Plan TDD cycles** (Prime → Test → Implement → Validate)
- **Work on any `implementation-plan.md`** file in docs/specs/

## Template

The PLAN template is at [template.md](template.md). Use this structure exactly.

**To write template to spec directory:**
1. Read the template: `plugins/start/skills/specify-plan/template.md`
2. Write to spec directory: `docs/specs/[NNN]-[name]/implementation-plan.md`

## PLAN Focus Areas

Your plan MUST answer these questions:
- **WHAT** produces value? (deliverables, not activities)
- **IN WHAT ORDER** do tasks execute? (dependencies and sequencing)
- **HOW TO VALIDATE** correctness? (test-first approach)
- **WHERE** is each task specified? (links to PRD/SDD sections)

Keep plans **actionable and focused**:
- Use task descriptions, sequence, and validation criteria
- Omit time estimates—focus on what, not when
- Omit resource assignments—focus on work, not who
- Omit implementation code—the plan guides, implementation follows

## Task Granularity Principle

**Track logical units that produce verifiable outcomes.** The TDD cycle is the execution method, not separate tracked items.

### Good Tracking Units (produces outcome)
- "Payment Entity" → Produces: working entity with tests ✓
- "Stripe Adapter" → Produces: working integration with tests ✓
- "Payment Form Component" → Produces: working UI with tests ✓

### Bad Tracking Units (too granular)
- "Read payment interface contracts" → Preparation, not deliverable
- "Test Payment.validate() rejects negative amounts" → Part of larger outcome
- "Run linting" → Validation step, not deliverable

### Structure Pattern

```markdown
- [ ] **T1.1 Payment Entity** `[activity: domain-modeling]`

  **Prime**: Read payment interface contracts `[ref: SDD/Section 4.2; lines: 145-200]`

  **Test**: Entity validation rejects negative amounts; supports currency conversion; handles refunds

  **Implement**: Create `src/domain/Payment.ts` with validation logic

  **Validate**: Run unit tests, lint, typecheck
```

The checkbox tracks "Payment Entity" as a unit. Prime/Test/Implement/Validate are embedded guidance.

## TDD Phase Structure

Every task follows red-green-refactor within this pattern:

### 1. Prime Context
- Read relevant specification sections
- Understand interfaces and contracts
- Load patterns and examples

### 2. Write Tests (Red)
- Test behavior before implementation
- Reference PRD acceptance criteria
- Cover happy path and edge cases

### 3. Implement (Green)
- Build to pass tests
- Follow SDD architecture
- Use discovered patterns

### 4. Validate (Refactor)
- Run automated tests
- Check code quality (lint, format)
- Verify specification compliance

## Task Metadata

Use these annotations in the plan:

```markdown
- [ ] T1.2.1 [Task description] `[ref: SDD/Section 5; lines: 100-150]` `[activity: backend-api]`
```

| Metadata | Description |
|----------|-------------|
| `[parallel: true]` | Tasks that can run concurrently |
| `[component: name]` | For multi-component features |
| `[ref: doc/section; lines: X-Y]` | Links to specifications |
| `[activity: type]` | Hint for specialist selection |

## Cycle Pattern

For each phase requiring definition, follow this iterative process:

### 1. Discovery Phase
- **Read PRD and SDD** to understand requirements and design
- **Identify activities** needed for each implementation area
- **Launch parallel specialist agents** to investigate:
  - Task sequencing and dependencies
  - Testing strategies
  - Risk assessment
  - Validation approaches

### 2. Documentation Phase
- **Update the PLAN** with task definitions
- **Add specification references** (`[ref: ...]`)
- Focus only on current phase being defined
- Follow template structure exactly

### 3. Review Phase
- **Present task breakdown** to user
- Show dependencies and sequencing
- Highlight parallel opportunities
- **Wait for user confirmation** before next phase

**Ask yourself each cycle:**
1. Have I read the relevant PRD and SDD sections?
2. Do all tasks trace back to specification requirements?
3. Are dependencies between tasks clear?
4. Can parallel tasks actually run in parallel?
5. Are validation steps included in each phase?
6. Have I received user confirmation?

## Specification Compliance

Every phase should include a validation task:

```markdown
- [ ] **T1.3 Phase Validation** `[activity: validate]`

  Run all phase tests, linting, type checking. Verify against SDD patterns and PRD acceptance criteria.
```

For complex phases, validation is embedded in each task's **Validate** step.

### Deviation Protocol

When implementation requires changes from the specification:
1. Document the deviation with clear rationale
2. Obtain approval before proceeding
3. Update SDD when the deviation improves the design
4. Record all deviations in the plan for traceability

## Validation Checklist

See [validation.md](validation.md) for the complete checklist. Key gates:

- [ ] All specification file paths are correct and exist
- [ ] Context priming section is complete
- [ ] All implementation phases are defined
- [ ] Each phase follows TDD: Prime → Test → Implement → Validate
- [ ] Dependencies between phases are clear (no circular dependencies)
- [ ] Parallel work is properly tagged with `[parallel: true]`
- [ ] Activity hints provided for specialist selection `[activity: type]`
- [ ] Every phase references relevant SDD sections
- [ ] Every test references PRD acceptance criteria
- [ ] Integration & E2E tests defined in final phase
- [ ] Project commands match actual project setup
- [ ] A developer could follow this plan independently

## Output Schema

### PLAN Status Report

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| specId | string | Yes | Spec identifier (NNN-name format) |
| phases | PhaseStatus[] | Yes | Status of each implementation phase |
| taskSummary | TaskSummary | Yes | Aggregate task metrics |
| specCoverage | SpecCoverage | Yes | PRD/SDD mapping completeness |
| nextSteps | string[] | Yes | Recommended next actions |

### PhaseStatus

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | Yes | Phase name |
| status | enum: COMPLETE, IN_PROGRESS, PENDING | Yes | Current state |
| taskCount | number | Yes | Number of tasks in phase |

### TaskSummary

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| total | number | Yes | Total task count |
| parallelGroups | number | Yes | Number of parallel execution groups |
| keyDependencies | string[] | Yes | Critical dependency chains |

### SpecCoverage

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| prdRequirementsMapped | string | Yes | Fraction mapped (e.g., "8/10") |
| sddComponentsCovered | string | Yes | Fraction covered (e.g., "5/5") |

---

## Examples

See [examples/phase-examples.md](examples/phase-examples.md) for reference.

---

## Entry Point

1. Read project context (Vision)
2. Activate when conditions met (When to Activate)
3. Load template from `template.md`
4. Execute iterative cycles per phase (Cycle Pattern)
5. Verify specification compliance and deviation protocol
6. Check against validation checklist and success criteria
7. Present output per PLAN Status Report schema
