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

[NEEDS CLARIFICATION: Define implementation phases. Each phase is a logical unit of work following TDD principles.]

- [ ] T1 Phase 1 [What functionality this phase delivers]

    - [ ] T1.1 Prime Context [What detailed sections are relevant for this phase from the specification]
        - [ ] T1.1.1 [Read interface contracts] `[ref: file; lines: 1-10]`
    - [ ] T1.2 Write Tests [What behavior needs to be tested]
        - [ ] T1.2.1 [Specific test case to verify specific behaviour] `[ref: file; lines: 123]` `[activity: type]`
    - [ ] T1.3 Implement [What needs to be built to pass tests]
        - [ ] T1.3.1 [Key implementation details or steps if complex] `[activity: type]`
    - [ ] T1.4 Implement [Additional implementation if needed] `[activity: type]`
    - [ ] T1.5 Validate [Verify that implementation is according to quality gates]
        - [ ] T1.5.1 [Review code so it is according to the defined quality gates] `[activity: lint-code, format-code, review-code, more]`
        - [ ] T1.5.2 [Validate code by running automated test and check commands] `[activity: run-tests, more]`
        - [ ] T1.5.3 [Ensure specification compliance] `[activity: business-acceptance, more]`

- [ ] T2 Phase 2 [What functionality this phase delivers]

    - [ ] T2.1 [Sub-phase/Component A] `[parallel: true]` `[component: name]`
        - [ ] T2.1.1 Prime Context [What detailed sections are relevant for this phase from the specification]
        - [ ] T2.1.2 Write Tests [What behavior needs to be tested]
        - [ ] T2.1.3 Implement [What needs to be built to pass tests]
        - [ ] T2.1.4 Validate [Verify that implementation is according to quality gates]

    - [ ] T2.2 [Sub-phase/Component B] `[parallel: true]` `[component: name]`
        - [ ] T2.2.1 Prime Context [What detailed sections are relevant for this phase from the specification]
        - [ ] T2.2.2 Write Tests [What behavior needs to be tested]
        - [ ] T2.2.3 Implement [What needs to be built to pass tests]
        - [ ] T2.2.4 Validate [Verify that implementation is according to quality gates]

- [ ] T3 Integration & End-to-End Validation

    - [ ] T3.1 [All unit tests passing per component if multi-component]
    - [ ] T3.2 [Integration tests for component interactions]
    - [ ] T3.3 [End-to-end tests for complete user flows]
    - [ ] T3.4 [Performance tests meet requirements] `[ref: SDD/Section 10 "Quality Requirements"]`
    - [ ] T3.5 [Security validation passes] `[ref: SDD/Section 10 "Security Requirements"]`
    - [ ] T3.6 [Acceptance criteria verified against PRD] `[ref: PRD acceptance criteria sections]`
    - [ ] T3.7 [Test coverage meets standards]
    - [ ] T3.8 [Documentation updated for any API/interface changes]
    - [ ] T3.9 [Build and deployment verification]
    - [ ] T3.10 [All PRD requirements implemented]
    - [ ] T3.11 [Implementation follows SDD design]
