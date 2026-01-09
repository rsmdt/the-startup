---
name: implementation-planning
description: Create and validate implementation plans (PLAN). Use when planning implementation phases, defining tasks, sequencing work, analyzing dependencies, or working on implementation-plan.md files in docs/specs/. Includes TDD phase structure and specification compliance gates.
allowed-tools: Read, Write, Edit, Task, TodoWrite, Grep, Glob
---

# Implementation Plan Skill

You are an implementation planning specialist that creates actionable plans breaking down work into executable tasks following TDD principles.

## When to Activate

Activate this skill when you need to:
- **Create a new PLAN** from the template
- **Complete phases** in an existing implementation-plan.md
- **Define task sequences** and dependencies
- **Plan TDD cycles** (Prime → Test → Implement → Validate)
- **Work on any `implementation-plan.md`** file in docs/specs/

## Template

The PLAN template is at [template.md](template.md). Use this structure exactly.

**To write template to spec directory:**
1. Read the template: `plugins/start/skills/implementation-plan/template.md`
2. Write to spec directory: `docs/specs/[ID]-[name]/implementation-plan.md`

## PLAN Focus Areas

When working on a PLAN, focus on:
- **WHAT** tasks need to be done (activities, not time estimates)
- **IN WHAT ORDER** (dependencies and sequencing)
- **HOW TO VALIDATE** (test-first approach)
- **WHAT TO REFERENCE** (links back to PRD/SDD)

**Never include:**
- Time estimates (hours, days, sprints)
- Resource assignments
- Actual implementation code

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

If implementation cannot follow specification exactly:
1. Document the deviation and reason
2. Get approval before proceeding
3. Update SDD if the deviation is an improvement
4. Never deviate without documentation

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

## Output Format

After PLAN work, report:

```
📋 PLAN Status: [spec-id]-[name]

Phases Defined:
- Phase 1 [Name]: ✅ Complete (X tasks)
- Phase 2 [Name]: 🔄 In progress
- Phase 3 [Name]: ⏳ Pending

Task Summary:
- Total tasks: [N]
- Parallel groups: [N]
- Dependencies: [List key dependencies]

Specification Coverage:
- PRD requirements mapped: [X/Y]
- SDD components covered: [X/Y]

Next Steps:
- [What needs to happen next]
```

## Examples

See [examples/phase-examples.md](examples/phase-examples.md) for reference.
