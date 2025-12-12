# PLAN Validation Checklist

Use this checklist to validate Implementation Plan completeness before execution.

## Structure Validation

- [ ] **All required sections are complete** - No empty or placeholder sections
- [ ] **No [NEEDS CLARIFICATION] markers remain** - All markers replaced with content
- [ ] **Template structure preserved** - No sections added, removed, or reorganized

## Context Priming

- [ ] **Specification paths correct** - PRD and SDD file paths exist and are valid
- [ ] **Key decisions extracted** - Critical choices from SDD are highlighted
- [ ] **Project commands documented** - Actual commands from project setup
- [ ] **Pattern links provided** - References to relevant pattern documentation

## Phase Structure

- [ ] **All implementation phases defined** - Complete coverage of the feature
- [ ] **Each phase follows TDD structure**:
  - Prime Context (read specs, load patterns)
  - Write Tests (behavior verification first)
  - Implement (code to pass tests)
  - Validate (quality gates)
- [ ] **Phase boundaries are logical** - Clear separation of concerns

## Task Quality

- [ ] **Tasks are actionable** - Clear what needs to be done
- [ ] **Tasks are atomic** - Can be completed independently (where not dependent)
- [ ] **No time estimates included** - Focus on WHAT, not HOW LONG
- [ ] **Activity hints provided** - `[activity: type]` for specialist selection

## Dependencies

- [ ] **Dependencies between phases clear** - What must complete before what
- [ ] **No circular dependencies** - Phases can be ordered linearly
- [ ] **Parallel work tagged** - `[parallel: true]` where applicable
- [ ] **Component tags for multi-component** - `[component: name]` where needed

## Specification Traceability

- [ ] **Every phase references SDD** - `[ref: SDD/Section X]`
- [ ] **Every test references PRD criteria** - `[ref: PRD/Section Y]`
- [ ] **All PRD requirements covered** - Nothing from PRD is missing
- [ ] **All SDD components covered** - Nothing from architecture is skipped

## Validation Steps

- [ ] **Each phase has validation step** - T*.4 Validate present in each phase
- [ ] **Code review included** - Quality gate for code standards
- [ ] **Automated tests included** - Test execution gate
- [ ] **Specification compliance included** - Business acceptance gate

## Final Phase

- [ ] **Integration tests defined** - Cross-component testing
- [ ] **E2E tests defined** - Complete user flow testing
- [ ] **Performance validation included** - If performance requirements exist
- [ ] **Security validation included** - If security requirements exist
- [ ] **All PRD requirements verified** - Final acceptance criteria check

## Practical Validation

- [ ] **Project commands match setup** - Commands work in the actual project
- [ ] **File paths are realistic** - Directories and files match codebase
- [ ] **A developer could follow independently** - No assumed knowledge

## No-Go Items

These should NOT appear in a PLAN:
- [ ] **No time estimates** - Hours, days, sprints
- [ ] **No resource assignments** - Who does what
- [ ] **No implementation code** - Actual code snippets (examples in SDD)
- [ ] **No scope expansion** - Tasks beyond PRD/SDD scope

## Completion Criteria

✅ **PLAN is complete when:**
- All checklist items pass
- User has reviewed and approved the task breakdown
- Every PRD requirement maps to at least one task
- Every SDD component is covered by phases
- A developer can start implementation immediately
- Ready for `/start:implement` execution
