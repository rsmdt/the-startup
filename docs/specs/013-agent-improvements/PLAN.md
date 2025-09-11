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

---

## Context Ingestion

*GATE: You MUST fully read all files mentioned in this section before starting any implementation.*

### Specification

**PRD**: `docs/specs/013-agent-improvements/PRD.md` - Product Requirements (comprehensive user personas and feature requirements)
- **SDD**: `docs/specs/013-agent-improvements/SDD.md` - Complete Solution Design with architecture decisions
- **IMPLEMENTATION-GUIDE**: `docs/specs/013-agent-improvements/IMPLEMENTATION-GUIDE.md` - 4-phase transformation approach

### Key Design Decisions

**Manual Transformation with Claude Code**: Transform agents using Claude Code with documented patterns (no automated CLI tools) - `[ref: SDD/Section 9 "Architecture Decisions"]`
- **Pattern-Based Structure**: Adopt the Effective Claude Agent Pattern as standard for all 61 agents - `[ref: SDD/Section 9 "Architecture Decisions"]`
- **3-Layer Architecture**: Identity → Objectives → Boundaries structure for all agents - `[ref: SDD/Section 8.1 "Implementation Patterns"]`
- **HOW-to-WHAT Transformation**: Convert 70% HOW-focused to 30% HOW / 70% WHAT ratio - `[ref: SDD/Section 1 "Quality Goals"]`

### SDD Section Mapping

**Agent Analysis & Readiness**: SDD Section 5.1 "Components" - 61 agent files across 9 categories
- **Transformation Patterns**: SDD Section 5.4 "Implementation Examples" - HOW-to-WHAT conversion logic
- **Quality Validation**: SDD Section 11 "Test Specifications" - Comprehensive validation framework
- **Runtime Behavior**: SDD Section 6 "Runtime View" - Agent transformation pipeline
- **Interface Compliance**: SDD Section 5.3 "Interface Specifications" - Markdown format preservation

### Implementation Scope

**File Impact Analysis**:
- Total files affected: 61 agent files + 2 core orchestrators
- New files to create: 0 (pure transformation approach)
- Files to modify: 61 agents + the-chief.md + the-meta-agent.md  
- Files to remove: 0 (no breaking changes)

**Key Areas of Change**:
- `/assets/claude/agents/**/*.md`: Transform all 61 specialist agents following Effective Claude Agent Pattern
- `/assets/claude/agents/the-chief.md`: Update orchestrator to work with transformed agent boundaries
- `/assets/claude/agents/the-meta-agent.md`: Update agent generator to create new agents in target format

### Implementation Context

**Commands to run**: [Commands from SDD Section 3.4 "Project Commands"]
  - Testing: Manual validation against test scenarios (no automated tooling)
  - Building: `go build -o the-startup` (verify asset embedding works)
  - Validation: Expert review using validation checklist from agent-quality-metrics.md
- **Patterns to follow**: 
  - Declarative Agent Structure: `docs/patterns/declarative-agent-structure.md` - Target format for all 61 agents
  - Agent Quality Metrics: `docs/patterns/agent-quality-metrics.md` - `[ref: SDD/Section 11 "Test Specifications"]`
  - Effective Claude Agent Pattern: `docs/patterns/effective-claude-agent-pattern.md` - Core transformation pattern
- **Interfaces to preserve**: 
  - Agent File Format: `[ref: SDD/Section 5.3 "Interface Specifications"]` - YAML frontmatter + Markdown
  - Claude Code Invocation: Maintain existing agent invocation patterns

---

## Implementation Phases

[NEEDS CLARIFICATION: Define implementation phases. Each phase is a logical unit of work following TDD principles.]

- [ ] **Phase 1**: [What functionality this phase delivers]
    # Status: [Not Started | In Progress | Completed]
    # SDD Components: Section X.Y "[Component Name]", Section A.B "[Interface Name]"
    # Required reading: SDD Sections [X.Y, A.B, C.D], Architecture Decision #[N]
    
    - [ ] **Pre-implementation review**: Understand SDD Sections [X.Y, A.B, C.D]
        # Component design: Section [X.Y]
        # Interface contracts: Section [A.B]
        # Runtime behavior: Section [C.D]
        # Architecture decisions: Section 9 Decision #[N]
    - [ ] **Tests**: [What behavior needs to be tested] `[ref: SDD/Section X.Y "Specific Topic"]` `[activity: type]`
        # [Specific test scenarios or acceptance criteria]
        # Specification requirement: [What SDD mandates for this]
    - [ ] **Implementation**: [What needs to be built to pass tests] `[ref: SDD/Section X.Y]` `[activity: type]`
        # [Key implementation details or steps if complex]
        # Must match: [Specific SDD requirement or interface]
        # - [Sub-step 1 if needed]
        # - [Sub-step 2 if needed]
    - [ ] **Validation**: [How to verify the implementation] `[activity: type]`
        # [Specific validation criteria or commands from SDD Section 3.4]
    - [ ] **Specification compliance**: Verify implementation against SDD
        # Matches SDD Section [X.Y] design exactly
        # Interface contracts from Section [A.B] satisfied
        # No unauthorized deviations from design
        # Architecture decisions followed correctly

- [ ] **Phase 2**: [What functionality this phase delivers]
  # Status: [Not Started | In Progress | Completed]
  # SDD Components: Section X.Y "[Component Name]", Section A.B "[Interface Name]"
  # Required reading: SDD Sections [List specific sections]
  # [Dependencies or parallel execution notes]
  
  - [ ] **Pre-implementation review**: Understand SDD Sections [List]
      # Review sections: [List specific SDD sections]
      # Verify dependencies from Phase 1 align with SDD
  
  - [ ] [Sub-phase/Component A] `[parallel: true]` `[component: name]`
    - [ ] **Tests**: [What behavior needs to be tested] `[ref: SDD/Section X.Y]` `[activity: type]`
        # [Specific test scenarios or acceptance criteria]
        # Specification requirement: [What SDD mandates for this]
    - [ ] **Implementation**: [What needs to be built] `[ref: SDD/Section X.Y]` `[activity: type]`
        # [Key implementation details or steps if complex]
        # Must match: [Specific SDD requirement or interface]
    - [ ] **Validation**: [How to verify the implementation] `[activity: type]`
        # [Specific validation criteria or commands from SDD Section 3.4]
    - [ ] **Specification compliance**: Verify against SDD
        # Matches interface in Section [X.Y]
        # Follows pattern from Section [A.B]
        # Meets requirements from Section [C.D]
  
  - [ ] [Sub-phase/Component B] `[parallel: true]` `[component: name]`
    - [ ] **Tests**: [What behavior needs to be tested] `[ref: SDD/Section X.Y]` `[activity: type]`
        # [Specific test scenarios or acceptance criteria]
    - [ ] **Implementation**: [What needs to be built] `[ref: SDD/Section X.Y]` `[activity: type]`
        # [Key implementation details or steps if complex]
    - [ ] **Validation**: [How to verify the implementation] `[activity: type]`
        # [Specific validation criteria or commands from SDD Section 3.4]

- [ ] **Integration & End-to-End Validation**
    # Verify the complete feature works as designed
    - [ ] All unit tests passing (per component if multi-component)
    - [ ] Integration tests for component interactions
    - [ ] End-to-end tests for complete user flows
    - [ ] Performance tests meet requirements `[ref: SDD/Section 10 "Quality Requirements"]`
    - [ ] Security validation passes `[ref: SDD/Section 10 "Security Requirements"]`
    - [ ] Acceptance criteria verified against PRD `[ref: PRD acceptance criteria sections]`
    - [ ] Test coverage meets standards
    - [ ] Documentation updated for any API/interface changes
    - [ ] Build and deployment verification
    - [ ] All PRD requirements implemented
    - [ ] Implementation follows SDD design

## Final Verification Checklist

[NEEDS CLARIFICATION: Define project-specific verification criteria]

### Requirements Verification
- [ ] All SDD requirements implemented `[ref: SDD sections covered]`
- [ ] All PRD acceptance criteria met `[ref: PRD sections covered]`
- [ ] No regression in existing functionality

### Technical Verification
- [ ] All defined tests passing
- [ ] Type safety verified (if applicable)
- [ ] Build process successful
- [ ] No circular dependencies introduced
- [ ] Performance within acceptable parameters

### Quality Metrics
- [ ] Test coverage meets project standards
- [ ] Code review completed for all changes
- [ ] Documentation updated where necessary
- [ ] No critical TODOs remaining

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
- `[ref: document/section]` - Links to specifications, patterns, or interfaces
- `[activity: type]` - Activity hint for specialist agent selection
