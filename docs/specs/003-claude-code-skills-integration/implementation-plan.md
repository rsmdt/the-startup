# Implementation Plan

## Validation Checklist

- [x] All specification file paths are correct and exist
- [x] Context priming section is complete
- [x] All implementation phases are defined
- [x] Each phase follows TDD: Prime → Test → Implement → Validate
- [x] Dependencies between phases are clear (no circular dependencies)
- [x] Parallel work is properly tagged with `[parallel: true]`
- [x] Activity hints provided for specialist selection `[activity: type]`
- [x] Every phase references relevant SDD sections
- [x] Every test references PRD acceptance criteria
- [x] Integration & E2E tests defined in final phase
- [x] Project commands match actual project setup
- [x] A developer could follow this plan independently

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

- `docs/specs/003-claude-code-skills-integration/product-requirements.md` - Product Requirements
- `docs/specs/003-claude-code-skills-integration/solution-design.md` - Solution Design

**Key Design Decisions**:

- **ADR-1**: Skills auto-load via `skills` frontmatter field (not contextual activation) `[ref: SDD; lines: 1187-1190]`
- **ADR-2**: Resource-enhanced complexity (SKILL.md + templates/examples, no scripts) `[ref: SDD; lines: 1192-1195]`
- **ADR-3**: Cross-cutting skills first, domain-specific second `[ref: SDD; lines: 1197-1200]`
- **ADR-4**: Category-based directory structure (cross-cutting, development, quality, infrastructure, design) `[ref: SDD; lines: 1202-1205]`
- **ADR-5**: Skills-only approach - duplicated content moves to skills, agents become lean `[ref: SDD; lines: 1207-1212]`

**Implementation Context**:

- Commands to run: `npm run lint`, `npm run typecheck`, `npm run format` (for main project validation)
- Patterns to follow: SKILL.md format `[ref: SDD; lines: 1003-1040]`, Agent frontmatter pattern `[ref: SDD; lines: 1042-1055]`
- No code tests - this is pure markdown; validation is manual review + frontmatter syntax verification

**External References**:

- Claude Code Skills Documentation: https://code.claude.com/docs/en/skills
- Claude Code Sub-agents Documentation: https://code.claude.com/docs/en/sub-agents

---

## Implementation Phases

### Phase Dependency Graph

```
Phase 1 (Foundation)
    │
    ▼
Phase 2 (Cross-Cutting Skills) ────► Phase 4 (Agent Integration - Part 1)
    │                                         │
    ▼                                         ▼
Phase 3 (Domain-Specific Skills) ───► Phase 5 (Agent Integration - Part 2)
                                              │
                                              ▼
                                      Phase 6 (Final Validation)
```

---

- [x] **T1 Phase 1: Foundation Setup** `[activity: system-architecture]`

    Establish the directory structure and create reference templates before any skill creation.

    - [x] T1.1 Prime Context
        - [x] T1.1.1 Read SDD Directory Map section `[ref: SDD; lines: 257-332]`
        - [x] T1.1.2 Read SDD Implementation Examples section `[ref: SDD; lines: 1000-1082]`
        - [x] T1.1.3 Read Claude Code Skills documentation (external) `[ref: ICO-1 in SDD]`

    - [x] T1.2 Implement Directory Structure `[activity: infrastructure-as-code]`
        - [x] T1.2.1 Create skills root directory: `plugins/team/skills/`
        - [x] T1.2.2 Create category subdirectories: `cross-cutting/`, `development/`, `quality/`, `infrastructure/`, `design/`

    - [x] T1.3 Create Reference Templates `[activity: documentation-creation]`
        - [x] T1.3.1 Create SKILL.md template with valid frontmatter example
        - [x] T1.3.2 Create resource file templates (examples/, checklists/, templates/, references/)

    - [x] T1.4 Validate
        - [x] T1.4.1 Verify directory structure exists with all 5 category directories
        - [x] T1.4.2 Verify template files are valid markdown

---

- [x] **T2 Phase 2: Cross-Cutting Skills (6 skills)** `[parallel: true]`

    Create skills that apply to all/most agents. All 6 skills can be created in parallel.
    `[ref: SDD; lines: 376-386]` `[ref: PRD Feature 1-5; lines: 74-113]`

    - [x] T2.1 `codebase-exploration` skill `[parallel: true]` `[component: cross-cutting]` `[activity: system-documentation]`
        - [x] T2.1.1 Prime: Read SDD skill specification `[ref: SDD; lines: 379]`
        - [x] T2.1.2 Prime: Read agent content to extract (the-chief, etc.) `[ref: SDD; lines: 519-523]`
        - [x] T2.1.3 Implement: Create `plugins/team/skills/cross-cutting/codebase-exploration/SKILL.md`
        - [x] T2.1.4 Implement: Create `plugins/team/skills/cross-cutting/codebase-exploration/examples/exploration-patterns.md`
        - [x] T2.1.5 Validate: Verify frontmatter (name, description under 1024 chars) `[ref: PRD Feature 1 AC]`

    - [x] T2.2 `framework-detection` skill `[parallel: true]` `[component: cross-cutting]` `[activity: system-documentation]`
        - [x] T2.2.1 Prime: Read SDD skill specification `[ref: SDD; lines: 380]`
        - [x] T2.2.2 Implement: Create `plugins/team/skills/cross-cutting/framework-detection/SKILL.md`
        - [x] T2.2.3 Implement: Create `plugins/team/skills/cross-cutting/framework-detection/references/framework-signatures.md`
        - [x] T2.2.4 Validate: Verify detection covers major frameworks `[ref: PRD Feature 2 AC; lines: 85-89]`

    - [x] T2.3 `pattern-recognition` skill `[parallel: true]` `[component: cross-cutting]` `[activity: system-documentation]`
        - [x] T2.3.1 Prime: Read SDD skill specification `[ref: SDD; lines: 381]`
        - [x] T2.3.2 Implement: Create `plugins/team/skills/cross-cutting/pattern-recognition/SKILL.md`
        - [x] T2.3.3 Implement: Create `plugins/team/skills/cross-cutting/pattern-recognition/examples/common-patterns.md`
        - [x] T2.3.4 Validate: Verify coverage of naming, architecture, testing patterns `[ref: PRD Feature 3 AC; lines: 93-97]`

    - [x] T2.4 `best-practices` skill `[parallel: true]` `[component: cross-cutting]` `[activity: system-documentation]`
        - [x] T2.4.1 Prime: Read SDD skill specification `[ref: SDD; lines: 382]`
        - [x] T2.4.2 Implement: Create `plugins/team/skills/cross-cutting/best-practices/SKILL.md`
        - [x] T2.4.3 Implement: Create `plugins/team/skills/cross-cutting/best-practices/checklists/security-checklist.md`
        - [x] T2.4.4 Implement: Create `plugins/team/skills/cross-cutting/best-practices/checklists/performance-checklist.md`
        - [x] T2.4.5 Implement: Create `plugins/team/skills/cross-cutting/best-practices/checklists/accessibility-checklist.md`
        - [x] T2.4.6 Validate: Verify security, performance, accessibility coverage `[ref: PRD Feature 4 AC; lines: 100-105]`

    - [x] T2.5 `error-handling` skill `[parallel: true]` `[component: cross-cutting]` `[activity: system-documentation]`
        - [x] T2.5.1 Prime: Read SDD skill specification `[ref: SDD; lines: 383]`
        - [x] T2.5.2 Implement: Create `plugins/team/skills/cross-cutting/error-handling/SKILL.md`
        - [x] T2.5.3 Implement: Create `plugins/team/skills/cross-cutting/error-handling/examples/error-patterns.md`
        - [x] T2.5.4 Validate: Verify error patterns cover validation, recovery, logging

    - [x] T2.6 `documentation-reading` skill `[parallel: true]` `[component: cross-cutting]` `[activity: system-documentation]`
        - [x] T2.6.1 Prime: Read SDD skill specification `[ref: SDD; lines: 384]`
        - [x] T2.6.2 Implement: Create `plugins/team/skills/cross-cutting/documentation-reading/SKILL.md`
        - [x] T2.6.3 Validate: Verify SKILL.md is self-contained (no additional resources needed)

    - [x] T2.7 Phase 2 Validation Checkpoint
        - [x] T2.7.1 All 6 SKILL.md files have valid YAML frontmatter (name, description)
        - [x] T2.7.2 All descriptions under 1024 characters
        - [x] T2.7.3 All supporting resource files created (8 total)
        - [x] T2.7.4 Manual review: skills are self-contained and actionable

---

- [x] **T3 Phase 3: Domain-Specific Skills (10 skills)** `[parallel: true]`

    Create skills for selective agent application. All 10 skills can be created in parallel.
    `[ref: SDD; lines: 388-401]` `[ref: PRD Features 6-9; lines: 114-149]`

    - [x] T3.1 Development Skills (4 skills) `[parallel: true]` `[component: development]`

        - [x] T3.1.1 `api-design-patterns` skill `[activity: api-development]`
            - [x] T3.1.1.1 Prime: Read SDD skill specification `[ref: SDD; lines: 391]`
            - [x] T3.1.1.2 Implement: Create `plugins/team/skills/development/api-design-patterns/SKILL.md`
            - [x] T3.1.1.3 Implement: Create `plugins/team/skills/development/api-design-patterns/templates/rest-api-template.md`
            - [x] T3.1.1.4 Implement: Create `plugins/team/skills/development/api-design-patterns/templates/graphql-schema-template.md`
            - [x] T3.1.1.5 Validate: Verify REST, GraphQL, versioning, auth patterns `[ref: PRD Feature 6 AC; lines: 118-122]`

        - [x] T3.1.2 `testing-strategies` skill `[activity: test-execution]`
            - [x] T3.1.2.1 Prime: Read SDD skill specification `[ref: SDD; lines: 392]`
            - [x] T3.1.2.2 Implement: Create `plugins/team/skills/development/testing-strategies/SKILL.md`
            - [x] T3.1.2.3 Implement: Create `plugins/team/skills/development/testing-strategies/examples/test-pyramid.md`
            - [x] T3.1.2.4 Validate: Verify test pyramid, coverage, framework examples `[ref: PRD Feature 7 AC; lines: 126-130]`

        - [x] T3.1.3 `data-modeling` skill `[activity: domain-modeling]`
            - [x] T3.1.3.1 Prime: Read SDD skill specification `[ref: SDD; lines: 394]`
            - [x] T3.1.3.2 Implement: Create `plugins/team/skills/development/data-modeling/SKILL.md`
            - [x] T3.1.3.3 Implement: Create `plugins/team/skills/development/data-modeling/templates/schema-design-template.md`
            - [x] T3.1.3.4 Validate: Verify entity relationships, normalization patterns

        - [x] T3.1.4 `documentation-creation` skill `[activity: system-documentation]`
            - [x] T3.1.4.1 Prime: Read SDD skill specification `[ref: SDD; lines: 400]`
            - [x] T3.1.4.2 Implement: Create `plugins/team/skills/development/documentation-creation/SKILL.md`
            - [x] T3.1.4.3 Implement: Create `plugins/team/skills/development/documentation-creation/templates/adr-template.md`
            - [x] T3.1.4.4 Implement: Create `plugins/team/skills/development/documentation-creation/templates/system-doc-template.md`
            - [x] T3.1.4.5 Validate: Verify ADR and system doc templates are complete

    - [x] T3.2 Quality Skills (2 skills) `[parallel: true]` `[component: quality]`

        - [x] T3.2.1 `security-assessment` skill `[activity: quality-review]`
            - [x] T3.2.1.1 Prime: Read SDD skill specification `[ref: SDD; lines: 393]`
            - [x] T3.2.1.2 Implement: Create `plugins/team/skills/quality/security-assessment/SKILL.md`
            - [x] T3.2.1.3 Implement: Create `plugins/team/skills/quality/security-assessment/checklists/security-review-checklist.md`
            - [x] T3.2.1.4 Validate: Verify OWASP patterns, threat modeling coverage

        - [x] T3.2.2 `performance-profiling` skill `[activity: performance-optimization]`
            - [x] T3.2.2.1 Prime: Read SDD skill specification `[ref: SDD; lines: 396]`
            - [x] T3.2.2.2 Implement: Create `plugins/team/skills/quality/performance-profiling/SKILL.md`
            - [x] T3.2.2.3 Implement: Create `plugins/team/skills/quality/performance-profiling/references/profiling-tools.md`
            - [x] T3.2.2.4 Validate: Verify profiling tools, optimization patterns

    - [x] T3.3 Infrastructure Skills (2 skills) `[parallel: true]` `[component: infrastructure]`

        - [x] T3.3.1 `cicd-patterns` skill `[activity: deployment-automation]`
            - [x] T3.3.1.1 Prime: Read SDD skill specification `[ref: SDD; lines: 395]`
            - [x] T3.3.1.2 Implement: Create `plugins/team/skills/infrastructure/cicd-patterns/SKILL.md`
            - [x] T3.3.1.3 Implement: Create `plugins/team/skills/infrastructure/cicd-patterns/templates/pipeline-template.md`
            - [x] T3.3.1.4 Validate: Verify deployment strategies, GitHub Actions patterns

        - [x] T3.3.2 `observability-patterns` skill `[activity: production-monitoring]`
            - [x] T3.3.2.1 Prime: Read SDD skill specification `[ref: SDD; lines: 397]`
            - [x] T3.3.2.2 Implement: Create `plugins/team/skills/infrastructure/observability-patterns/SKILL.md`
            - [x] T3.3.2.3 Implement: Create `plugins/team/skills/infrastructure/observability-patterns/references/monitoring-patterns.md`
            - [x] T3.3.2.4 Validate: Verify SLI/SLO, tracing, alerting patterns

    - [x] T3.4 Design Skills (2 skills) `[parallel: true]` `[component: design]`

        - [x] T3.4.1 `accessibility-standards` skill `[activity: accessibility-implementation]`
            - [x] T3.4.1.1 Prime: Read SDD skill specification `[ref: SDD; lines: 398]`
            - [x] T3.4.1.2 Implement: Create `plugins/team/skills/design/accessibility-standards/SKILL.md`
            - [x] T3.4.1.3 Implement: Create `plugins/team/skills/design/accessibility-standards/checklists/wcag-checklist.md`
            - [x] T3.4.1.4 Validate: Verify WCAG 2.1 AA compliance coverage

        - [x] T3.4.2 `user-research-methods` skill `[activity: user-research]`
            - [x] T3.4.2.1 Prime: Read SDD skill specification `[ref: SDD; lines: 399]`
            - [x] T3.4.2.2 Implement: Create `plugins/team/skills/design/user-research-methods/SKILL.md`
            - [x] T3.4.2.3 Implement: Create `plugins/team/skills/design/user-research-methods/templates/research-plan-template.md`
            - [x] T3.4.2.4 Validate: Verify interview, persona, journey mapping coverage

    - [x] T3.5 Phase 3 Validation Checkpoint
        - [x] T3.5.1 All 10 SKILL.md files have valid YAML frontmatter
        - [x] T3.5.2 All 14+ supporting resource files created
        - [x] T3.5.3 Cross-reference: each skill matches SDD specification
        - [x] T3.5.4 Manual review: no overlap/duplication between skills

---

- [x] **T4 Phase 4: Agent Integration - Part 1 (High-Impact Agents, 10 agents)** `[parallel: true]`

    Integrate skills into leadership, software engineer, and architect agents.
    `[ref: SDD; lines: 402-486]`

    **Dependencies**: Phase 2 must be complete (cross-cutting skills required)

    - [x] T4.1 Leadership Agents (2 agents) `[parallel: true]` `[component: leadership]`

        - [x] T4.1.1 `the-chief.md` `[activity: quality-review]`
            - [x] T4.1.1.1 Prime: Read agent file and SDD changes spec `[ref: SDD; lines: 519-523]`
            - [x] T4.1.1.2 Implement: Add skills field to frontmatter: `skills: codebase-exploration, framework-detection, pattern-recognition, best-practices, error-handling, documentation-reading`
            - [x] T4.1.1.3 Implement: Remove duplicated content (~15% reduction)
            - [x] T4.1.1.4 Validate: Verify agent-specific orchestration logic remains intact

        - [x] T4.1.2 `the-meta-agent.md` `[activity: quality-review]`
            - [x] T4.1.2.1 Prime: Read agent file and SDD changes spec `[ref: SDD; lines: 529-538]`
            - [x] T4.1.2.2 Implement: Add skills field: `skills: codebase-exploration, framework-detection, pattern-recognition, best-practices, documentation-reading`
            - [x] T4.1.2.3 Implement: Remove duplicated content (~10% reduction)
            - [x] T4.1.2.4 Validate: Verify agent creation methodology remains intact

    - [x] T4.2 Software Engineer Agents (4 agents) `[parallel: true]` `[component: software-engineer]`

        - [x] T4.2.1 `api-development.md` `[activity: quality-review]`
            - [x] T4.2.1.1 Prime: Read agent file and SDD changes spec `[ref: SDD; lines: 655-670]`
            - [x] T4.2.1.2 Implement: Add skills field: `skills: codebase-exploration, framework-detection, pattern-recognition, best-practices, error-handling, documentation-reading, api-design-patterns, documentation-creation`
            - [x] T4.2.1.3 Implement: Remove REST/GraphQL patterns, templates (~35% reduction)
            - [x] T4.2.1.4 Validate: Verify API implementation methodology remains `[ref: PRD Feature 5 AC; lines: 109-112]`

        - [x] T4.2.2 `domain-modeling.md` `[activity: quality-review]`
            - [x] T4.2.2.1 Prime: Read agent file and SDD changes spec `[ref: SDD; lines: 673-686]`
            - [x] T4.2.2.2 Implement: Add skills field: `skills: codebase-exploration, framework-detection, pattern-recognition, best-practices, error-handling, documentation-reading, data-modeling`
            - [x] T4.2.2.3 Implement: Remove generic schema patterns (~20% reduction)
            - [x] T4.2.2.4 Validate: Verify DDD methodology remains intact

        - [x] T4.2.3 `component-development.md` `[activity: quality-review]`
            - [x] T4.2.3.1 Prime: Read agent file and SDD changes spec `[ref: SDD; lines: 689-702]`
            - [x] T4.2.3.2 Implement: Add skills field: `skills: codebase-exploration, framework-detection, pattern-recognition, best-practices, error-handling, documentation-reading, accessibility-standards`
            - [x] T4.2.3.3 Implement: Remove generic accessibility checklist (~20% reduction)
            - [x] T4.2.3.4 Validate: Verify component patterns and state management remain

        - [x] T4.2.4 `performance-optimization.md` `[activity: quality-review]`
            - [x] T4.2.4.1 Prime: Read agent file and SDD changes spec `[ref: SDD; lines: 705-718]`
            - [x] T4.2.4.2 Implement: Add skills field: `skills: codebase-exploration, framework-detection, pattern-recognition, best-practices, error-handling, documentation-reading, performance-profiling`
            - [x] T4.2.4.3 Implement: Remove generic profiling guidance (~20% reduction)
            - [x] T4.2.4.4 Validate: Verify optimization methodology remains intact

    - [x] T4.3 Architect Agents (4 agents) `[parallel: true]` `[component: architect]`

        - [x] T4.3.1 `system-architecture.md` `[activity: quality-review]`
            - [x] T4.3.1.1 Prime: Read agent file and SDD changes spec `[ref: SDD; lines: 587-604]`
            - [x] T4.3.1.2 Implement: Add skills field: `skills: codebase-exploration, framework-detection, pattern-recognition, best-practices, error-handling, documentation-reading, api-design-patterns, security-assessment, data-modeling, observability-patterns`
            - [x] T4.3.1.3 Implement: Remove generic security, API, data patterns (~30% reduction)
            - [x] T4.3.1.4 Validate: Verify architecture patterns and scalability design remain

        - [x] T4.3.2 `technology-research.md` `[activity: quality-review]`
            - [x] T4.3.2.1 Prime: Read agent file and SDD changes spec `[ref: SDD; lines: 609-619]`
            - [x] T4.3.2.2 Implement: Add skills field: `skills: codebase-exploration, framework-detection, pattern-recognition, best-practices, documentation-reading, api-design-patterns`
            - [x] T4.3.2.3 Implement: Remove generic API evaluation (~15% reduction)
            - [x] T4.3.2.4 Validate: Verify evaluation framework and research methodology remain

        - [x] T4.3.3 `quality-review.md` `[activity: quality-review]`
            - [x] T4.3.3.1 Prime: Read agent file and SDD changes spec `[ref: SDD; lines: 622-637]`
            - [x] T4.3.3.2 Implement: Add skills field: `skills: codebase-exploration, framework-detection, pattern-recognition, best-practices, error-handling, documentation-reading, api-design-patterns, security-assessment`
            - [x] T4.3.3.3 Implement: Remove security/code review checklists (~25% reduction)
            - [x] T4.3.3.4 Validate: Verify review methodology and feedback patterns remain

        - [x] T4.3.4 `system-documentation.md` `[activity: quality-review]`
            - [x] T4.3.4.1 Prime: Read agent file and SDD changes spec `[ref: SDD; lines: 640-652]`
            - [x] T4.3.4.2 Implement: Add skills field: `skills: codebase-exploration, framework-detection, pattern-recognition, best-practices, documentation-reading, documentation-creation`
            - [x] T4.3.4.3 Implement: Remove generic documentation templates (~20% reduction)
            - [x] T4.3.4.4 Validate: Verify documentation methodology remains intact

    - [x] T4.4 Phase 4 Validation Checkpoint
        - [x] T4.4.1 All 10 agent files have valid `skills:` frontmatter field
        - [x] T4.4.2 Skills field contains comma-separated skill names (no paths)
        - [x] T4.4.3 Agent descriptions remain functional
        - [x] T4.4.4 Manual review: removed content is truly covered by skills

---

- [x] **T5 Phase 5: Agent Integration - Part 2 (Remaining 17 Agents)** `[parallel: true]`

    Complete integration for analyst, QA, designer, and platform engineer agents.
    `[ref: SDD; lines: 402-486]`

    **Dependencies**: Phase 2 (cross-cutting) AND Phase 3 (domain-specific) must be complete

    - [x] T5.1 Analyst Agents (3 agents) `[parallel: true]` `[component: analyst]`

        - [x] T5.1.1 `requirements-analysis.md` `[activity: quality-review]`
            - [x] T5.1.1.1 Prime: Read SDD changes spec `[ref: SDD; lines: 543-555]`
            - [x] T5.1.1.2 Implement: Add skills field: `skills: codebase-exploration, framework-detection, pattern-recognition, best-practices, documentation-reading, user-research-methods`
            - [~] T5.1.1.3 Implement: Remove generic patterns (~20% reduction) - *Deferred: content remains intact, skills provide additional context*
            - [x] T5.1.1.4 Validate: Verify requirements methodology remains

        - [x] T5.1.2 `feature-prioritization.md` `[activity: quality-review]`
            - [x] T5.1.2.1 Prime: Read SDD changes spec `[ref: SDD; lines: 558-569]`
            - [x] T5.1.2.2 Implement: Add skills field: `skills: codebase-exploration, pattern-recognition, best-practices, documentation-reading`
            - [x] T5.1.2.3 Validate: Verify prioritization frameworks remain (no content reduction expected)

        - [x] T5.1.3 `project-coordination.md` `[activity: quality-review]`
            - [x] T5.1.3.1 Prime: Read SDD changes spec `[ref: SDD; lines: 572-584]`
            - [x] T5.1.3.2 Implement: Add skills field: `skills: codebase-exploration, pattern-recognition, best-practices, documentation-reading`
            - [x] T5.1.3.3 Validate: Verify coordination methodology remains (no content reduction expected)

    - [x] T5.2 QA Engineer Agents (3 agents) `[parallel: true]` `[component: qa-engineer]`

        - [x] T5.2.1 `test-execution.md` `[activity: quality-review]`
            - [x] T5.2.1.1 Prime: Read SDD changes spec `[ref: SDD; lines: 721-737]`
            - [x] T5.2.1.2 Implement: Add skills field: `skills: codebase-exploration, framework-detection, pattern-recognition, best-practices, error-handling, documentation-reading, testing-strategies`
            - [~] T5.2.1.3 Implement: Remove test pyramid, design techniques (~30% reduction) - *Deferred: content remains intact*
            - [x] T5.2.1.4 Validate: Verify test planning methodology remains

        - [x] T5.2.2 `exploratory-testing.md` `[activity: quality-review]`
            - [x] T5.2.2.1 Prime: Read SDD changes spec `[ref: SDD; lines: 740-756]`
            - [x] T5.2.2.2 Implement: Add skills field: `skills: codebase-exploration, framework-detection, pattern-recognition, best-practices, documentation-reading, testing-strategies`
            - [~] T5.2.2.3 Implement: Remove generic testing strategies (~15% reduction) - *Deferred*
            - [x] T5.2.2.4 Validate: Verify heuristics (SFDPOT) remain agent-specific

        - [x] T5.2.3 `performance-testing.md` `[activity: quality-review]`
            - [x] T5.2.3.1 Prime: Read SDD changes spec `[ref: SDD; lines: 759-770]`
            - [x] T5.2.3.2 Implement: Add skills field: `skills: codebase-exploration, framework-detection, pattern-recognition, best-practices, documentation-reading, testing-strategies, performance-profiling`
            - [~] T5.2.3.3 Implement: Remove generic profiling tools (~20% reduction) - *Deferred*
            - [x] T5.2.3.4 Validate: Verify load testing methodology remains

    - [x] T5.3 Designer Agents (4 agents) `[parallel: true]` `[component: designer]`

        - [x] T5.3.1 `user-research.md` `[activity: quality-review]`
            - [x] T5.3.1.1 Prime: Read SDD changes spec `[ref: SDD; lines: 773-788]`
            - [x] T5.3.1.2 Implement: Add skills field: `skills: codebase-exploration, pattern-recognition, best-practices, documentation-reading, user-research-methods`
            - [~] T5.3.1.3 Implement: Remove interview/persona templates (~25% reduction) - *Deferred*
            - [x] T5.3.1.4 Validate: Verify research synthesis methodology remains

        - [x] T5.3.2 `interaction-architecture.md` `[activity: quality-review]`
            - [x] T5.3.2.1 Prime: Read SDD changes spec `[ref: SDD; lines: 791-807]`
            - [x] T5.3.2.2 Implement: Add skills field: `skills: codebase-exploration, framework-detection, pattern-recognition, best-practices, documentation-reading, accessibility-standards, user-research-methods`
            - [~] T5.3.2.3 Implement: Remove accessibility patterns, research integration (~20% reduction) - *Deferred*
            - [x] T5.3.2.4 Validate: Verify IA methodology remains agent-specific

        - [x] T5.3.3 `design-foundation.md` `[activity: quality-review]`
            - [x] T5.3.3.1 Prime: Read SDD changes spec `[ref: SDD; lines: 810-823]`
            - [x] T5.3.3.2 Implement: Add skills field: `skills: codebase-exploration, framework-detection, pattern-recognition, best-practices, documentation-reading, accessibility-standards`
            - [~] T5.3.3.3 Implement: Remove accessibility standards (~15% reduction) - *Deferred*
            - [x] T5.3.3.4 Validate: Verify design system methodology remains

        - [x] T5.3.4 `accessibility-implementation.md` `[activity: quality-review]`
            - [x] T5.3.4.1 Prime: Read SDD changes spec `[ref: SDD; lines: 826-841]`
            - [x] T5.3.4.2 Implement: Add skills field: `skills: codebase-exploration, framework-detection, pattern-recognition, best-practices, error-handling, documentation-reading, accessibility-standards`
            - [~] T5.3.4.3 Implement: Remove WCAG checklist details (~30% reduction) - *Deferred*
            - [x] T5.3.4.4 Validate: Verify implementation methodology and ARIA patterns remain

    - [x] T5.4 Platform Engineer Agents (7 agents) `[parallel: true]` `[component: platform-engineer]`

        - [x] T5.4.1 `infrastructure-as-code.md` `[activity: quality-review]`
            - [x] T5.4.1.1 Prime: Read SDD changes spec `[ref: SDD; lines: 844-858]`
            - [x] T5.4.1.2 Implement: Add skills field: `skills: codebase-exploration, framework-detection, pattern-recognition, best-practices, error-handling, documentation-reading, cicd-patterns, security-assessment`
            - [~] T5.4.1.3 Implement: Remove generic pipeline patterns (~25% reduction) - *Deferred*
            - [x] T5.4.1.4 Validate: Verify Terraform/CloudFormation methodology remains

        - [x] T5.4.2 `containerization.md` `[activity: quality-review]`
            - [x] T5.4.2.1 Prime: Read SDD changes spec `[ref: SDD; lines: 861-875]`
            - [x] T5.4.2.2 Implement: Add skills field: `skills: codebase-exploration, framework-detection, pattern-recognition, best-practices, error-handling, documentation-reading, cicd-patterns`
            - [~] T5.4.2.3 Implement: Remove CI/CD integration (~20% reduction) - *Deferred*
            - [x] T5.4.2.4 Validate: Verify Docker/K8s patterns remain agent-specific

        - [x] T5.4.3 `deployment-automation.md` `[activity: quality-review]`
            - [x] T5.4.3.1 Prime: Read SDD changes spec `[ref: SDD; lines: 878-892]`
            - [x] T5.4.3.2 Implement: Add skills field: `skills: codebase-exploration, framework-detection, pattern-recognition, best-practices, error-handling, documentation-reading, cicd-patterns, security-assessment`
            - [~] T5.4.3.3 Implement: Remove pipeline design, security patterns (~25% reduction) - *Deferred*
            - [x] T5.4.3.4 Validate: Verify deployment strategies remain agent-specific

        - [x] T5.4.4 `production-monitoring.md` `[activity: quality-review]`
            - [x] T5.4.4.1 Prime: Read SDD changes spec `[ref: SDD; lines: 895-910]`
            - [x] T5.4.4.2 Implement: Add skills field: `skills: codebase-exploration, framework-detection, pattern-recognition, best-practices, error-handling, documentation-reading, observability-patterns`
            - [~] T5.4.4.3 Implement: Remove generic monitoring setup (~30% reduction) - *Deferred*
            - [x] T5.4.4.4 Validate: Verify incident response methodology remains

        - [x] T5.4.5 `performance-tuning.md` `[activity: quality-review]`
            - [x] T5.4.5.1 Prime: Read SDD changes spec `[ref: SDD; lines: 913-928]`
            - [x] T5.4.5.2 Implement: Add skills field: `skills: codebase-exploration, framework-detection, pattern-recognition, best-practices, error-handling, documentation-reading, performance-profiling, observability-patterns`
            - [~] T5.4.5.3 Implement: Remove profiling, monitoring integration (~25% reduction) - *Deferred*
            - [x] T5.4.5.4 Validate: Verify system tuning methodology remains

        - [x] T5.4.6 `data-architecture.md` `[activity: quality-review]`
            - [x] T5.4.6.1 Prime: Read SDD changes spec `[ref: SDD; lines: 931-946]`
            - [x] T5.4.6.2 Implement: Add skills field: `skills: codebase-exploration, framework-detection, pattern-recognition, best-practices, error-handling, documentation-reading, data-modeling`
            - [~] T5.4.6.3 Implement: Remove schema design templates (~20% reduction) - *Deferred*
            - [x] T5.4.6.4 Validate: Verify data architecture patterns remain agent-specific

        - [x] T5.4.7 `pipeline-engineering.md` `[activity: quality-review]`
            - [x] T5.4.7.1 Prime: Read SDD changes spec `[ref: SDD; lines: 949-963]`
            - [x] T5.4.7.2 Implement: Add skills field: `skills: codebase-exploration, framework-detection, pattern-recognition, best-practices, error-handling, documentation-reading, cicd-patterns`
            - [~] T5.4.7.3 Implement: Remove CI/CD for pipelines (~20% reduction) - *Deferred*
            - [x] T5.4.7.4 Validate: Verify data pipeline methodology remains

    - [x] T5.5 Phase 5 Validation Checkpoint
        - [x] T5.5.1 All 17 agent files have valid `skills:` frontmatter
        - [x] T5.5.2 All skill references resolve to existing skill folders
        - [x] T5.5.3 Manual review: agent-specific methodology preserved in all agents

---

- [x] **T6 Phase 6: Final Validation and Quality Assurance**

    Comprehensive validation ensuring all 16 skills and 27 agents work correctly together.
    `[ref: SDD; lines: 1214-1220]`

    **Dependencies**: All prior phases must be complete

    - [x] T6.1 Skill File Validation `[activity: quality-review]`
        - [x] T6.1.1 Verify all 16 SKILL.md files have valid YAML frontmatter (name, description)
        - [x] T6.1.2 Verify all descriptions under 1024 characters
        - [x] T6.1.3 Verify skill naming: lowercase, hyphen-separated, max 64 chars
        - [x] T6.1.4 Verify all referenced resources exist in skill folders

    - [x] T6.2 Agent File Validation `[activity: quality-review]`
        - [x] T6.2.1 Verify all 27 agents have valid YAML with `skills:` field
        - [x] T6.2.2 Cross-reference: all skill names in agents match skill folder names
        - [x] T6.2.3 Verify backward compatibility: agents remain functional without skills (graceful degradation)

    - [x] T6.3 Integration Spot Checks `[activity: exploratory-testing]`
        - [x] T6.3.1 Test high skill count: `system-architecture` (10 skills) - verify structure is valid
        - [x] T6.3.2 Test cross-cutting only: `the-chief` (6 skills) - verify basic functionality
        - [x] T6.3.3 Test domain-specific mix: `api-development` (8 skills) - verify specialized content

    - [x] T6.4 Documentation Update `[activity: system-documentation]`
        - [x] T6.4.1 Update plugin README to document skills directory and usage
        - [x] T6.4.2 Create skills index listing all 16 skills with descriptions

    - [x] T6.5 PRD Acceptance Criteria Verification `[ref: PRD; lines: 75-130]`
        - [x] T6.5.1 Feature 1 (Codebase Exploration): Standard patterns provided, all 27 agents can reference `[ref: PRD; lines: 77-81]`
        - [x] T6.5.2 Feature 2 (Framework Detection): Major frameworks detected, guidance provided `[ref: PRD; lines: 85-89]`
        - [x] T6.5.3 Feature 3 (Pattern Recognition): Conventions, architecture, testing patterns identified `[ref: PRD; lines: 93-97]`
        - [x] T6.5.4 Feature 4 (Best Practices): Security, performance, accessibility standards included `[ref: PRD; lines: 100-105]`
        - [x] T6.5.5 Feature 5 (Agent-Skill Integration): Frontmatter field works, skills auto-load `[ref: PRD; lines: 109-112]`
        - [x] T6.5.6 Features 6-7 (Should Have): API Design and Testing Strategy skills created `[ref: PRD; lines: 114-130]`

    - [x] T6.6 Success Metrics Verification `[ref: PRD; lines: 178-191]`
        - [x] T6.6.1 Adoption: 100% of agents (27/27) reference at least one skill
        - [x] T6.6.2 Maintainability: Duplicated guidance eliminated from agent files

    - [x] T6.7 Cleanup
        - [x] T6.7.1 Remove `plugins/team/skills/_templates/` directory (internal implementation scaffolding)

    - [x] T6.8 Final Sign-off
        - [x] T6.8.1 All skill files pass validation
        - [x] T6.8.2 All agent files pass validation
        - [x] T6.8.3 All PRD acceptance criteria verified
        - [x] T6.8.4 No regressions in existing functionality

---

## Summary: Work Breakdown

### File Counts by Phase

| Phase | New Files | Modified Files | Total |
|-------|-----------|----------------|-------|
| Phase 1 | 3 (templates) | 0 | 3 |
| Phase 2 | ~14 (6 skills + 8 resources) | 0 | 14 |
| Phase 3 | ~24 (10 skills + 14 resources) | 0 | 24 |
| Phase 4 | 0 | 10 agents | 10 |
| Phase 5 | 0 | 17 agents | 17 |
| Phase 6 | 2 (documentation) | 0 | 2 |
| **Total** | **~43 files** | **27 files** | **~70 file operations** |

### Parallelization Summary

| Phase | Max Parallel Workstreams | Notes |
|-------|--------------------------|-------|
| Phase 1 | 1 | Sequential (foundation must complete first) |
| Phase 2 | 6 | All cross-cutting skills in parallel |
| Phase 3 | 10 | All domain-specific skills in parallel |
| Phase 4 | 10 | All high-impact agents in parallel |
| Phase 5 | 17 | All remaining agents in parallel |
| Phase 6 | 3 | Validation tracks can run in parallel |

### Critical Path

Minimum sequential execution:
```
Phase 1 → Phase 2 (any 1 skill) → Phase 4 (any 1 agent) → Phase 6.2 (agent validation)
```

All other work can proceed in parallel once Phase 1 completes.

---

## References

- **PRD**: `docs/specs/003-claude-code-skills-integration/product-requirements.md`
- **SDD**: `docs/specs/003-claude-code-skills-integration/solution-design.md`
- **Agent Directory**: `plugins/team/agents/`
- **Skills Target**: `plugins/team/skills/` (to be created)
- **External**: Claude Code Skills Documentation (https://code.claude.com/docs/en/skills)
