# Implementation Plan

## Validation Checklist

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

---

## Specification Compliance Guidelines

### How to Ensure Specification Adherence

1. **Before Each Phase**: Complete the Pre-Implementation Specification Gate
2. **During Implementation**: Reference specific SDD sections in each task
3. **After Each Task**: Run Specification Compliance checks
4. **Phase Completion**: Verify all specification requirements are met

### Deviation Protocol

If implementation cannot follow specification exactly:
1. Document the deviation and reason
2. Get approval before proceeding
3. Update SDD if the deviation is an improvement
4. Never deviate without documentation

## Metadata Reference

- `[parallel: true]` - Tasks that can run concurrently
- `[component: component-name]` - For multi-component features
- `[ref: document/section; lines: 1, 2-3]` - Links to specifications, patterns, or interfaces and (if applicable) line(s)
- `[activity: type]` - Activity hint for specialist agent selection

---

## Context Priming

*GATE: You MUST fully read all files mentioned in this section before starting any implementation.*

**Specification**:

[NEEDS CLARIFICATION: Replace file location with actual path and add/remove files accordingly]
- `docs/specs/[ID]-[feature-name]/PRD.md` - Product Requirements (if exists)
- `docs/specs/[ID]-[feature-name]/SDD.md` - Solution Design

**Key Design Decisions**:

[NEEDS CLARIFICATION: Extract critical decisions from the SDD]
- [Critical decision 1]
- [Critical decision 2]

**Implementation Context**:

[NEEDS CLARIFICATION: Extract actionable information from specs]
- Commands to run: [Project-specific commands from SDD for testing, building, etc.]
- Patterns to follow: [Links to relevant pattern docs]
- Interfaces to implement: [Links to interface specifications]

---

## Implementation Phases

Each task follows red-green-refactor: **Prime** (understand context), **Test** (red), **Implement** (green), **Validate** (refactor + verify).

> **Tracking Principle**: Track logical units that produce verifiable outcomes. The TDD cycle is the method, not separate tracked items.

---

### Phase 1: [What functionality this phase delivers]

[NEEDS CLARIFICATION: Define tasks that each produce a verifiable outcome]

- [ ] **T1.1 [Logical work unit name]** `[activity: type]`

  **Prime**: [What specification sections to read] `[ref: SDD/Section; lines: X-Y]`

  **Test**: [Key behaviors to verify - list scenarios, not individual test files]

  **Implement**: [What to build - key files/components to create or modify]

  **Validate**: Run tests, lint, typecheck; verify against specification

- [ ] **T1.2 [Next logical work unit]** `[activity: type]`

  **Prime**: [Context needed]

  **Test**: [Key behaviors]

  **Implement**: [What to build]

  **Validate**: Run tests, lint, typecheck

---

### Phase 2: [What functionality this phase delivers]

[NEEDS CLARIFICATION: For phases with parallel work, each parallel unit is a tracked item]

- [ ] **T2.1 [Component/Area A]** `[parallel: true]` `[component: name]`

  **Prime**: [Context] `[ref: ...]`

  **Test**: [Behaviors to verify]

  **Implement**: [Files to create/modify]

  **Validate**: Component tests pass

- [ ] **T2.2 [Component/Area B]** `[parallel: true]` `[component: name]`

  **Prime**: [Context]

  **Test**: [Behaviors]

  **Implement**: [Files]

  **Validate**: Component tests pass

- [ ] **T2.3 Phase Validation** `[activity: validate]`

  Run all phase tests, linting, type checking. Verify components integrate correctly.

---

### Phase 3: Integration & Validation

- [ ] **T3.1 Integration Testing** - Cross-component integration tests pass `[ref: SDD/integration points]`
- [ ] **T3.2 E2E Testing** - End-to-end user flows verified `[ref: PRD/acceptance criteria]`
- [ ] **T3.3 Quality Gates** - Performance, security, coverage requirements met `[ref: SDD/Quality Requirements]`
- [ ] **T3.4 Specification Compliance** - All PRD/SDD requirements verified, documentation updated
