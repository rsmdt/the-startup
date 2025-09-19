# Implementation Plan

## Validation Checklist
- [ ] Context Ingestion section complete with all required specs
- [ ] Implementation phases logically organized
- [ ] Each phase starts with test definition (TDD approach)
- [ ] Dependencies between phases identified
- [ ] Parallel execution marked where applicable
- [ ] Multi-component coordination identified (if applicable)
- [ ] Final validation phase included
- [ ] No placeholder content remains

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

- [ ] **Phase 1**: [What functionality this phase delivers]

    - [ ] **Prime Context**: [What detailed sections are relevant for this phase frmo the specification]
        - [ ] [Read interface contracts] `[ref: file; lines: 1-10]`
    - [ ] **Write Tests**: [What behavior needs to be tested]
        - [ ] [Specific test case to verify specific behaviour] `[ref: file; lines: 123]` `[activity: type]`
    - [ ] **Implement**: [What needs to be built to pass tests]
        - [ ] [Key implementation details or steps if complex] `[activity: type]`
    - [ ] **Implement**: [Additional implementation if needed] `[activity: type]`
    - [ ] **Validate**: [Verify that implementation is according to quality gates]
        - [ ] [Review code so it is according to the defined quality gates] `[activity: lint-code, format-code, review-code, more]`
        - [ ] [Validate code by running automated test and check commands] `[activity: run-tests, more]`
        - [ ] [Ensure specification compliance] `[activity: business-acceptance, more]`

- [ ] **Phase 2**: [What functionality this phase delivers]

    - [ ] [Sub-phase/Component A] `[parallel: true]` `[component: name]`
        - [ ] **Prime Context**: [What detailed sections are relevant for this phase frmo the specification]
        - [ ] **Write Tests**: [What behavior needs to be tested]
        - [ ] **Implement**: [What needs to be built to pass tests]
        - [ ] **Validate**: [Verify that implementation is according to quality gates]

    - [ ] [Sub-phase/Component B] `[parallel: true]` `[component: name]`
        [Same as above]

- [ ] **Integration & End-to-End Validation**
    - [ ] [All unit tests passing per component if multi-component]
    - [ ] [Integration tests for component interactions]
    - [ ] [End-to-end tests for complete user flows]
    - [ ] [Performance tests meet requirements] `[ref: SDD/Section 10 "Quality Requirements"]`
    - [ ] [Security validation passes] `[ref: SDD/Section 10 "Security Requirements"]`
    - [ ] [Acceptance criteria verified against PRD] `[ref: PRD acceptance criteria sections]`
    - [ ] [Test coverage meets standards]
    - [ ] [Documentation updated for any API/interface changes]
    - [ ] [Build and deployment verification]
    - [ ] [All PRD requirements implemented]
    - [ ] [Implementation follows SDD design]
