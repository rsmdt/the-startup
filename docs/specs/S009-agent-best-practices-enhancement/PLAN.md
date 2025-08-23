# Implementation Plan

## Context Documents

### Context Validation Checkpoints

Before proceeding with implementation, verify:
- [x] All referenced specifications have been located and read (BRD, PRD, SDD)
- [x] Key requirements have been extracted and understood (user validation responses received)
- [x] Conflicts between specifications have been identified and resolved (no conflicts found)
- [x] Missing context has been noted with mitigation strategies (manual testing approach confirmed)
- [x] Implementation approach aligns with discovered patterns (@docs/patterns/ created)

## Implementation Overview

**Goal**: Systematically enhance Claude Code sub-agents with industry best practices using "@" directive inclusion of condensed rule files.

**Key Constraints**:
- Keep agents concise to avoid LLM confusion
- Use "@" directive pattern: `@{{STARTUP_PATH}}/rules/abc.md`
- Follow ~/.claude/CLAUDE.md level of detail (no over-explanation)
- Manual testing only (no automation available)
- Preserve agent personalities and domain focus

**Enhancement Phases** (User Validated Priority):
1. Architecture agents → 2. Requirements agents → 3. Engineering agents → 4. Remaining specialists

## Phase 1: Architecture Agents Enhancement

**Target Agents**: `the-software-architect`, `the-staff-engineer`
**Priority**: First (user validated)

- [ ] **Create architecture practices rule file** [`complexity: medium`] [`review: architecture, patterns`]
  - [ ] Create `assets/the-startup/rules/architecture-practices.md`
  - [ ] Include SOLID principles (no explanation - LLM understands)
  - [ ] Add design pattern documentation guidance
  - [ ] Include architectural asset creation instructions (docs/patterns/, docs/interfaces/)
  - [ ] Add hexagonal architecture guidance
  - [ ] Keep condensed without excessive nested structure
  - [ ] **Validate**: Manual review of rule file content

- [ ] **Enhance the-software-architect agent** [`complexity: medium`] [`parallel: true`]
  - [ ] Add `@{{STARTUP_PATH}}/rules/architecture-practices.md` to Focus Areas section
  - [ ] Preserve existing agent personality and domain focus
  - [ ] Maintain YAML frontmatter and 4-section structure
  - [ ] **Validate**: Manual testing with Claude Code Task tool
  - [ ] **Review**: Architecture practices integration

- [ ] **Enhance the-staff-engineer agent** [`complexity: medium`] [`parallel: true`]
  - [ ] Add `@{{STARTUP_PATH}}/rules/architecture-practices.md` to Focus Areas section
  - [ ] Preserve existing agent personality and domain focus
  - [ ] Maintain YAML frontmatter and 4-section structure
  - [ ] **Validate**: Manual testing with Claude Code Task tool
  - [ ] **Review**: Architecture practices integration

- [ ] **Phase 1 Validation** [`complexity: low`]
  - [ ] Test enhanced architecture agents with sample tasks
  - [ ] Verify "@" directive inclusion works correctly
  - [ ] Confirm agents create/enhance patterns and interfaces appropriately
  - [ ] Validate personality preservation
  - [ ] **Review**: Phase 1 completion assessment

## Phase 2: Requirements Validation Agents Enhancement

**Target Agents**: `the-business-analyst`, `the-product-manager`
**Priority**: Second (user validated)

- [ ] **Create requirements validation practices rule file** [`complexity: medium`] [`review: requirements, validation`]
  - [ ] Create `assets/the-startup/rules/requirements-validation-practices.md`
  - [ ] Include question-driven clarification protocols
  - [ ] Add "I'm assuming X, please confirm" patterns
  - [ ] Include assumption prevention guidelines
  - [ ] Add stakeholder validation requirements
  - [ ] Follow condensed structure pattern
  - [ ] **Validate**: Manual review of rule file content

- [ ] **Enhance the-business-analyst agent** [`complexity: medium`] [`parallel: true`]
  - [ ] Add `@{{STARTUP_PATH}}/rules/requirements-validation-practices.md` to Focus Areas section
  - [ ] Preserve existing questioning and analysis personality
  - [ ] Ensure assumption prevention doesn't reduce effectiveness
  - [ ] **Validate**: Manual testing with vague requirements scenario
  - [ ] **Review**: Assumption prevention effectiveness

- [ ] **Enhance the-product-manager agent** [`complexity: medium`] [`parallel: true`]
  - [ ] Add `@{{STARTUP_PATH}}/rules/requirements-validation-practices.md` to Focus Areas section
  - [ ] Preserve existing prioritization and planning focus
  - [ ] Ensure stakeholder validation integration
  - [ ] **Validate**: Manual testing with feature planning scenario
  - [ ] **Review**: Stakeholder validation integration

- [ ] **Phase 2 Validation** [`complexity: low`]
  - [ ] Test enhanced requirements agents with ambiguous scenarios
  - [ ] Verify agents ask clarifying questions before proceeding
  - [ ] Confirm "I'm assuming" pattern usage
  - [ ] Validate assumption detection and prevention
  - [ ] **Review**: Phase 2 completion assessment

## Phase 3: Engineering Team Enhancement

**Target Agents**: `the-lead-engineer`, `the-frontend-engineer`, `the-backend-engineer`, `the-mobile-engineer`, `the-ml-engineer`
**Priority**: Third (user validated)

- [ ] **Create software development practices rule file** [`complexity: high`] [`review: development, security`]
  - [ ] Create `assets/the-startup/rules/software-development-practices.md`
  - [ ] Include TDD workflows (Red-Green-Refactor - no explanation needed)
  - [ ] Add domain-specific security practices
  - [ ] Include tool integration guidance ("run tests" = 100% passing)
  - [ ] Add "do not stop in middle of testing" guidance
  - [ ] Keep tool-agnostic language throughout
  - [ ] **Validate**: Manual review of rule file content

- [ ] **Enhance lead engineer agent** [`complexity: medium`] [`agent: the-lead-engineer`]
  - [ ] Add `@{{STARTUP_PATH}}/rules/software-development-practices.md` to Focus Areas section
  - [ ] Preserve code review and mentorship focus
  - [ ] Integrate quality practices with existing approach
  - [ ] **Validate**: Manual testing with code review scenario
  - [ ] **Review**: Code quality practices integration

- [ ] **Enhance development agents in parallel** [`complexity: medium`] [`parallel: true`]
  - [ ] **Frontend Engineer**: Add software development practices
    - [ ] Preserve UI/UX and performance focus
    - [ ] Include client-side security practices
    - [ ] **Validate**: Manual testing with component development
  - [ ] **Backend Engineer**: Add software development practices  
    - [ ] Preserve API and business logic focus
    - [ ] Include server-side security practices
    - [ ] **Validate**: Manual testing with API development
  - [ ] **Mobile Engineer**: Add software development practices
    - [ ] Preserve platform-specific focus
    - [ ] Include mobile security practices
    - [ ] **Validate**: Manual testing with mobile feature
  - [ ] **ML Engineer**: Add software development practices
    - [ ] Preserve model integration focus
    - [ ] Include ML-specific security practices
    - [ ] **Validate**: Manual testing with ML integration

- [ ] **Phase 3 Validation** [`complexity: medium`]
  - [ ] Test enhanced engineering agents with TDD scenarios
  - [ ] Verify agents guide through Red-Green-Refactor workflow
  - [ ] Confirm security practices integration
  - [ ] Test "all tests working means 100%" enforcement
  - [ ] Validate tool-agnostic guidance effectiveness
  - [ ] **Review**: Engineering practices integration

## Phase 4: Remaining Specialists Enhancement

**Target**: QA/Security, Infrastructure, Design/Documentation agents
**Priority**: Fourth (user validated)

- [ ] **Create remaining practice rule files** [`complexity: medium`] [`parallel: true`]
  - [ ] **QA Practices**: Create `assets/the-startup/rules/quality-assurance-practices.md`
    - [ ] Include testing strategies and validation approaches
    - [ ] Add behavior vs implementation testing guidance
    - [ ] Include mock boundaries and edge case testing
  - [ ] **Infrastructure Practices**: Create `assets/the-startup/rules/infrastructure-practices.md`  
    - [ ] Include DevOps automation practices
    - [ ] Add monitoring and performance optimization
    - [ ] Include infrastructure security practices
  - [ ] **Design/Docs Practices**: Create `assets/the-startup/rules/design-documentation-practices.md`
    - [ ] Include user-centered design practices
    - [ ] Add accessibility and clarity requirements
    - [ ] Include documentation validation patterns
  - [ ] **Validate**: Manual review of all rule files

- [ ] **Enhance QA and Security agents** [`complexity: medium`] [`parallel: true`]
  - [ ] **QA Lead**: Add quality assurance practices
  - [ ] **QA Engineer**: Add quality assurance practices  
  - [ ] **Security Engineer**: Add quality assurance practices (complement existing security focus)
  - [ ] **Compliance Officer**: Add quality assurance practices
  - [ ] **Validate**: Manual testing with testing and security scenarios

- [ ] **Enhance Infrastructure agents** [`complexity: medium`] [`parallel: true`]  
  - [ ] **DevOps Engineer**: Add infrastructure practices
  - [ ] **Site Reliability Engineer**: Add infrastructure practices
  - [ ] **Data Engineer**: Add infrastructure practices
  - [ ] **Performance Engineer**: Add infrastructure practices  
  - [ ] **Validate**: Manual testing with infrastructure scenarios

- [ ] **Enhance Design and Documentation agents** [`complexity: low`] [`parallel: true`]
  - [ ] **UX Designer**: Add design/documentation practices
  - [ ] **Principal Designer**: Add design/documentation practices
  - [ ] **Technical Writer**: Add design/documentation practices
  - [ ] **Validate**: Manual testing with design and documentation scenarios

- [ ] **Phase 4 Validation** [`complexity: medium`]
  - [ ] Test all remaining enhanced agents
  - [ ] Verify practice integration effectiveness
  - [ ] Confirm personality preservation across all agents
  - [ ] **Review**: Final phase completion

## Final Integration and Validation

- [ ] **Complete System Validation** [`complexity: high`] [`review: system, integration`]
  - [ ] Test representative agents from each phase
  - [ ] Verify cross-agent practice consistency
  - [ ] Confirm "@" directive integration works across all enhanced agents
  - [ ] Validate no conflicts between different practice rules
  - [ ] Test agents in realistic development scenarios
  - [ ] **Review**: Complete system validation

- [ ] **Documentation Completion** [`complexity: low`]
  - [ ] Ensure all pattern documentation is complete
  - [ ] Update any interfaces created during implementation
  - [ ] Verify rule files follow condensed structure pattern
  - [ ] Document any lessons learned during implementation

- [ ] **Final Quality Assessment** [`complexity: medium`] [`review: quality, effectiveness`]
  - [ ] Compare enhanced agents vs original agent effectiveness
  - [ ] Assess personality preservation across all 25+ agents
  - [ ] Validate industry best practices integration quality
  - [ ] Confirm manual testing approach sufficiency
  - [ ] **Review**: Final quality and effectiveness review

## Success Criteria

- [ ] All 25+ agents enhanced with appropriate industry best practices
- [ ] "@" directive integration working correctly across all agents
- [ ] Agent personalities and domain expertise preserved
- [ ] Rule files follow condensed, non-nested structure pattern
- [ ] Manual testing confirms enhanced agents provide better guidance
- [ ] No conflicts between different practice rules
- [ ] Architecture agents successfully create/enhance patterns and interfaces
- [ ] Requirements agents prevent assumptions through systematic validation
- [ ] Engineering agents guide through proper TDD workflows
- [ ] All specialists integrate practices relevant to their domains

## Implementation Notes

- **Manual Testing**: Each agent enhancement requires manual testing via Claude Code Task tool
- **Phase Dependencies**: Complete each phase fully before proceeding to next
- **Personality Preservation**: Critical success factor - agents must maintain their unique characteristics
- **Rule File Quality**: Follow existing rule file patterns for structure and conciseness
- **Tool-Agnostic Language**: Maintain flexibility across different technology stacks
- **Practice Consistency**: Same practices should be described consistently across relevant agents