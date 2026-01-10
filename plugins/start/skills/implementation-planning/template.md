---
title: "[NEEDS CLARIFICATION: Feature title]"
status: draft
version: "1.0"
---

# Implementation Plan

## Validation Checklist

### CRITICAL GATES (Must Pass)

- [ ] All `[NEEDS CLARIFICATION: ...]` markers have been addressed
- [ ] All specification file paths are correct and exist
- [ ] Each phase follows TDD: Prime → Test → Implement → Validate
- [ ] Every task has verifiable success criteria
- [ ] A developer could follow this plan independently

### QUALITY CHECKS (Should Pass)

- [ ] Context priming section is complete
- [ ] All implementation phases are defined
- [ ] Dependencies between phases are clear (no circular dependencies)
- [ ] Parallel work is properly tagged with `[parallel: true]`
- [ ] Activity hints provided for specialist selection `[activity: type]`
- [ ] Every phase references relevant SDD sections
- [ ] Every test references PRD acceptance criteria
- [ ] Integration & E2E tests defined in final phase
- [ ] Project commands match actual project setup

---

## Specification Compliance Guidelines

### How to Ensure Specification Adherence

1. **Before Each Phase**: Complete the Pre-Implementation Specification Gate
2. **During Implementation**: Reference specific SDD sections in each task
3. **After Each Task**: Run Specification Compliance checks
4. **Phase Completion**: Verify all specification requirements are met

### Deviation Protocol

When implementation requires changes from the specification:
1. Document the deviation with clear rationale
2. Obtain approval before proceeding
3. Update SDD when the deviation improves the design
4. Record all deviations in this plan for traceability

## Metadata Reference

- `[parallel: true]` - Tasks that can run concurrently
- `[component: component-name]` - For multi-component features
- `[ref: document/section; lines: 1, 2-3]` - Links to specifications, patterns, or interfaces and (if applicable) line(s)
- `[activity: type]` - Activity hint for specialist agent selection

### Success Criteria

**Validate** = Process verification ("did we follow TDD?")
**Success** = Outcome verification ("does it work correctly?")

```markdown
# Single-line format
- Success: [Criterion] `[ref: PRD/AC-X.Y]`

# Multi-line format
- Success:
  - [ ] [Criterion 1] `[ref: PRD/AC-X.Y]`
  - [ ] [Criterion 2] `[ref: SDD/Section]`
```

---

## Context Priming

*GATE: Read all files in this section before starting any implementation.*

**Specification**:

[NEEDS CLARIFICATION: Replace with actual paths from your spec directory]
- `docs/specs/[NNN]-[name]/product-requirements.md` - Product Requirements
- `docs/specs/[NNN]-[name]/solution-design.md` - Solution Design
- `docs/{patterns,interfaces,research}/[name].md`  - Optional references

**Key Design Decisions**:

[NEEDS CLARIFICATION: Extract 2-5 critical decisions from the SDD that affect implementation]
- **ADR-1**: [Decision name] - [One-line summary of what was decided and why]
- **ADR-2**: [Decision name] - [One-line summary]

**Implementation Context**:

[NEEDS CLARIFICATION: Update implementation-relevant commands to match project setup]
```bash
# Testing
npm test                    # Unit tests
npm run test:integration    # Integration tests

# Quality
npm run lint               # Linting
npm run typecheck          # Type checking

# Full validation
npm run validate           # All checks
```

---

## Implementation Phases

Each task follows red-green-refactor: **Prime** (understand context), **Test** (red), **Implement** (green), **Validate** (refactor + verify).

> **Tracking Principle**: Track logical units that produce verifiable outcomes. The TDD cycle is the method, not separate tracked items.

---

### Phase 1: Core Foundation

[NEEDS CLARIFICATION: Describe what this phase delivers - "Establishes X capability" or "Enables Y functionality"]

Establishes the foundational [domain/infrastructure/components] required for subsequent phases.

- [ ] **T1.1 [Primary deliverable name]** `[activity: domain-modeling]`

  - Prime: Read [entity/component] specification `[ref: SDD/Section X; lines: Y-Z]`
  - Test: [Behavior 1]; [Behavior 2]; [Edge case handling]
  - Implement: Create `src/[path]/[File].ts` with [key capability]
  - Validate: Unit tests pass; lint clean; types check
  - Success: [Criterion 1] `[ref: PRD/AC-1.1]`; [Criterion 2] `[ref: SDD/Section]`

- [ ] **T1.2 [Secondary deliverable name]** `[activity: data-architecture]`

  - Prime: Review [pattern/interface] requirements `[ref: SDD/Section X]`
  - Test: [CRUD operations]; [Query patterns]; [Error handling]
  - Implement: Create `src/[path]/[File].ts` with [key capability]
  - Validate: Integration tests pass; lint clean; types check
  - Success: [Criterion 1] `[ref: PRD/AC-1.2]`; [Criterion 2] `[ref: SDD/Section]`

- [ ] **T1.3 Phase Validation** `[activity: validate]`

  - Run all Phase 1 tests. Verify foundation matches SDD data models. Lint and typecheck pass.

---

### Phase 2: [API/Integration/UI] Layer

[NEEDS CLARIFICATION: Describe what this phase delivers. Mark parallel tasks for concurrent execution.]

Builds the [API endpoints/integrations/UI components] that expose Phase 1 capabilities.

- [ ] **T2.1 [Component A]** `[parallel: true]` `[component: backend]`

  - Prime: Read API specification `[ref: SDD/Section X]`
  - Test: [Endpoint behavior]; [Validation]; [Error responses]
  - Implement: Create `src/[path]/[Controller].ts` with routes
  - Validate: API tests pass; contract matches specification
  - Success: [Criterion] `[ref: PRD/AC-2.1]`

- [ ] **T2.2 [Component B]** `[parallel: true]` `[component: backend]`

  - Prime: Read integration pattern `[ref: docs/interfaces/X.md]`
  - Test: [Success flow]; [Failure handling]; [Retry logic]
  - Implement: Create `src/[path]/[Adapter].ts` with integration
  - Validate: Integration tests pass with test/mock service
  - Success:
    - [Success criterion 1] `[ref: PRD/AC-2.2]`
    - [Failure handling verified] `[ref: SDD/Error Handling]`
    - [Retry logic works with backoff] `[ref: docs/patterns/resilience.md]`

- [ ] **T2.3 Phase Validation** `[activity: validate]`

  Run all Phase 2 tests. Verify API contracts match SDD. Lint and typecheck pass.

---

### Phase 3: Final Integration & Validation

Full system validation ensuring all components work together correctly.

- [ ] **T3.1 Integration Testing** `[activity: integration-test]`

  Verify cross-component integration: [Component A] ↔ [Component B]; [Service] ↔ [Database]

  `[ref: SDD/integration points]`

- [ ] **T3.2 E2E User Flows** `[activity: e2e-test]`

  Verify complete user journeys: [Happy path flow]; [Error handling flow]; [Edge case flow]

  `[ref: PRD/acceptance criteria]`

- [ ] **T3.3 Quality Gates** `[activity: validate]`

  - Performance: [Specific metric] < [threshold] `[ref: SDD/Quality Requirements]`
  - Security: [Specific check] verified
  - Coverage: > [X]% line coverage

- [ ] **T3.4 Specification Compliance** `[activity: business-acceptance]`

  - All PRD acceptance criteria verified
  - Implementation follows SDD design
  - Documentation updated for API/interface changes
  - Build and deployment verification complete

---

## Plan Verification

Before this plan is ready for implementation, verify:

| Criterion | Status |
|-----------|--------|
| A developer can follow this plan without additional clarification | ⬜ |
| Every task produces a verifiable deliverable | ⬜ |
| All PRD acceptance criteria map to specific tasks | ⬜ |
| All SDD components have implementation tasks | ⬜ |
| Dependencies are explicit with no circular references | ⬜ |
| Parallel opportunities are marked with `[parallel: true]` | ⬜ |
| Each task has specification references `[ref: ...]` | ⬜ |
| Project commands in Context Priming are accurate | ⬜ |
