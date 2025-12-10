# SDD Validation Checklist

Use this checklist to validate SDD completeness before proceeding to Implementation Plan.

## Structure Validation

- [ ] **All required sections are complete** - No empty or placeholder sections
- [ ] **No [NEEDS CLARIFICATION] markers remain** - All markers replaced with content
- [ ] **Template structure preserved** - No sections added, removed, or reorganized

## Context Validation

- [ ] **All context sources listed** - Every relevant file, doc, URL documented
- [ ] **Relevance ratings assigned** - CRITICAL/HIGH/MEDIUM/LOW for each source
- [ ] **Why explanations provided** - Each source explains its relevance
- [ ] **Project commands discovered** - Commands from actual project files, not assumed

## Constraints & Strategy

- [ ] **Constraints identified** - Technical, business, compliance constraints documented
- [ ] **Strategy follows from constraints** - Logical progression from constraints to approach
- [ ] **Architecture pattern stated** - Clear pattern name (layered, modular, microservice, etc.)
- [ ] **Pattern rationale provided** - Why this pattern fits the constraints

## Component Design

- [ ] **Every component has directory mapping** - Where code will live
- [ ] **Every component in diagram** - No orphan components
- [ ] **Component names consistent** - Same names across diagrams, text, and directory structure
- [ ] **Responsibilities clear** - What each component does (single responsibility)
- [ ] **No component overlap** - No duplicate responsibilities

## Interface Design

- [ ] **All interfaces specified** - API endpoints, data models, events documented
- [ ] **Request/response formats defined** - Clear data structures
- [ ] **Error responses documented** - Error codes and messages
- [ ] **Authentication specified** - How each interface is secured
- [ ] **External integrations documented** - Third-party APIs and their contracts

## Data Design

- [ ] **Data models defined** - Entities, fields, types
- [ ] **Database schema specified** - Tables, columns, relationships (if applicable)
- [ ] **Data flow documented** - How data moves through the system
- [ ] **Storage decisions explained** - Why specific storage choices

## Cross-Cutting Concerns

- [ ] **Error handling covers all types** - Validation, business, system errors
- [ ] **Security patterns defined** - Auth, authz, encryption, input validation
- [ ] **Logging/observability planned** - What to log, how to monitor
- [ ] **Performance considerations** - Caching, batching, async patterns
- [ ] **Deployment approach defined** - Environment, config, dependencies

## Quality Requirements

- [ ] **Quality requirements are specific** - Numbers and units, not vague
- [ ] **Quality requirements are measurable** - Can be tested
- [ ] **Every quality requirement has test coverage** - How it will be verified

## Architecture Decisions

- [ ] **ADRs documented** - All significant decisions recorded
- [ ] **Rationale provided** - Why each choice was made
- [ ] **Trade-offs acknowledged** - What we accept with each decision
- [ ] **All ADRs confirmed by user** - User has approved each decision

## PRD Alignment

- [ ] **All PRD requirements addressable** - Every requirement maps to design
- [ ] **No scope creep** - Nothing designed beyond PRD scope
- [ ] **Acceptance criteria satisfiable** - Design enables meeting all criteria

## Implementability

- [ ] **A developer could implement from this** - Clear enough for someone new
- [ ] **Implementation examples provided** - For complex logic (where helpful)
- [ ] **Test scenarios defined** - Gherkin scenarios for key flows
- [ ] **Dependencies identified** - External libraries, services, etc.

## Consistency Checks

Run these checks across the entire document:
- [ ] **Component names match everywhere** - Diagrams, text, directory structure
- [ ] **Interface names match everywhere** - APIs, data models, references
- [ ] **No contradictions between sections** - Consistent throughout
- [ ] **Patterns applied consistently** - Same approach for similar problems

## Completion Criteria

✅ **SDD is complete when:**
- All checklist items pass
- All ADRs have user confirmation
- No open questions remain
- A developer could implement without further clarification
- Ready for implementation planning (PLAN)
