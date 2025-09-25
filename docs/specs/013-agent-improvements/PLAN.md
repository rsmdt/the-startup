# Implementation Plan

## Validation Checklist
- [x] Context Ingestion section complete with all required specs
- [x] Implementation phases logically organized
- [x] Each phase starts with test definition (TDD approach)
- [x] Dependencies between phases identified
- [x] Parallel execution marked where applicable
- [x] Multi-component coordination identified (if applicable)
- [x] Final validation phase included
- [x] No placeholder content remains

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

## Context Ingestion

*GATE: You MUST fully read all files mentioned in this section before starting any implementation.*

### Specification

- **SDD**: `docs/specs/013-agent-improvements/SDD.md` - Solution Design for Agent Transformation

### Key Design Decisions

- **Manual Claude Code-Assisted Transformation**: Use Claude Code with documented patterns (no automated CLI tools)
- **Pattern-Based Structure**: Adopt Effective Claude Agent Pattern for all 61 agents
- **70% Declarative Target**: Convert from 70% HOW-focused to 70% WHAT-focused patterns
- **Phased Migration**: 4 phases over realistic timeline to minimize risk
- **Zero Capability Loss**: 100% functional parity requirement

### Implementation Context

- **Commands to run**: Manual transformation using Claude Code, Git for version control, requirements compliance review
- **Patterns to follow**: Enhanced Agent Template structure from `docs/patterns/claude-agent-pattern.md`
- **Interfaces to implement**: Agent file format (YAML frontmatter + Markdown), delegation boundaries
- **Quality Gates**: Semantic preservation, delegation accuracy (95% improvement), performance (<10% degradation)

### Transformation Template

Apply this exact structure to each agent (from SDD lines 240-289):

```markdown
---
name: agent-identifier
description: Use this agent when [detailed scenario description]. This includes [specific tasks]. Examples:\n\n<example>\nContext: [situation]\nuser: "[request]"\nassistant: "[response using agent]"\n<commentary>\n[explanation of why agent is appropriate]\n</commentary>\n</example>
model: inherit
---

You are an expert [role] specializing in [detailed expertise areas]. [Additional expertise context].

**Core Responsibilities:**

You will [primary action] that:
- [Specific outcome with measurable result]
- [Specific outcome with measurable result]
- [Specific outcome with measurable result]

**[Domain] Methodology:**

1. **[Phase Name]:**
   - [Specific principle or approach]
   - [Specific principle or approach]

2. **[Phase Name]:**
   - [Specific principle or approach]
   - [Specific principle or approach]

**Output Format:**

You will provide:
1. [Specific deliverable with details]
2. [Specific deliverable with details]

**Best Practices:**

- [Detailed positive principle]
- [Detailed positive principle]

You approach [domain] with the mindset that [detailed philosophy about quality and approach].
```

### Transformation Rules

1. **HOW to WHAT Conversion**: Replace procedural steps with outcome statements
2. **Preserve ALL Capabilities**: Every function from original must exist in transformed version
3. **Section Mapping**:
   - Old "Approach" sections → New "Core Responsibilities"
   - Old numbered steps → New methodology phases (high-level)
   - Old "Anti-Patterns" → New "Best Practices" (positive framing)
4. **Description Enhancement**: Add multiple examples with commentary tags
5. **Keep Framework Detection**: Preserve technology adaptation sections

---

## Implementation Phases

**Execution Timeline**: Based on SDD categorization of 61 agents into 9 domains

**Agent File Locations**: All agents are in `assets/agents/` directory structure:
- Core agents: `assets/agents/the-chief.md`, `assets/agents/the-meta-agent.md`
- Domain agents: `assets/agents/[domain-name]/[agent-name].md`

- [x] **Phase 1: Core System Agents** (Priority: CRITICAL)

    - [x] **Load Context**: Reference implementation patterns from SDD
        - [x] Read SDD transformation methodology `[ref: docs/specs/013-agent-improvements/SDD.md; lines: 236-289]`
        - [x] Review Claude-generated reference agents `[ref: .claude/agents/test-writer.md, api-design-architect.md]`

    - [x] **the-chief.md** `[parallel: false]`
        - [x] Apply Enhanced Agent Template pattern to orchestrator
        - [x] Preserve parallel execution and delegation logic
        - [x] Validate against SDD orchestration requirements `[activity: transformation]`

    - [x] **the-meta-agent.md** `[parallel: false]`
        - [x] Transform agent generation patterns to declarative
        - [x] Maintain agent creation capabilities
        - [x] Validate pattern compliance `[activity: transformation]`

- [x] **Phase 2: Software Engineering Domain** (10 agents)

    - [x] **API & Database Design** `[parallel: true]` `[component: the-software-engineer]`
        - [x] Transform api-design.md to declarative pattern `[activity: transformation]`
        - [x] Transform database-design.md to declarative pattern `[activity: transformation]`
        - [x] Transform api-documentation.md to declarative pattern `[activity: transformation]`
        - [x] Validate delegation boundaries preserved `[activity: review]`

    - [x] **Component & State Management** `[parallel: true]` `[component: the-software-engineer]`
        - [x] Transform component-architecture.md to declarative pattern `[activity: transformation]`
        - [x] Transform state-management.md to declarative pattern `[activity: transformation]`
        - [x] Transform browser-compatibility.md to declarative pattern `[activity: transformation]`
        - [x] Validate pattern consistency across related agents `[activity: review]`

    - [x] **Business Logic & Integration** `[parallel: true]` `[component: the-software-engineer]`
        - [x] Transform business-logic.md to declarative pattern `[activity: transformation]`
        - [x] Transform service-integration.md to declarative pattern `[activity: transformation]`
        - [x] Transform reliability-engineering.md to declarative pattern `[activity: transformation]`
        - [x] Transform performance-optimization.md to declarative pattern `[activity: transformation]`
        - [x] Validate against SDD quality goals `[activity: review]`

- [x] **Phase 3: Platform & Architecture Domains** (18 agents)

    - [x] **Platform Engineers** `[parallel: true]` `[component: the-platform-engineer]`
        - [x] Transform 11 platform engineering agents:
            - [x] system-performance.md, observability.md, containerization.md `[activity: transformation]`
            - [x] pipeline-engineering.md, ci-cd-automation.md, deployment-strategies.md `[activity: transformation]`
            - [x] incident-response.md, infrastructure-as-code.md `[activity: transformation]`
            - [x] storage-architecture.md, query-optimization.md, data-modeling.md `[activity: transformation]`
        - [x] Validate infrastructure patterns preserved `[activity: review]`

    - [x] **Architects** `[parallel: true]` `[component: the-architect]`
        - [x] Transform 7 architecture agents:
            - [x] system-design.md, system-documentation.md `[activity: transformation]`
            - [x] scalability-planning.md, technology-standards.md `[activity: transformation]`
            - [x] technology-evaluation.md, architecture-review.md, code-review.md `[activity: transformation]`
        - [x] Validate architectural decision patterns `[activity: review]`

- [x] **Phase 4: Quality & Security Domains** (9 agents)

    - [x] **QA Engineers** `[parallel: true]` `[component: the-qa-engineer]`
        - [x] Transform 4 QA agents:
            - [x] test-strategy.md, test-implementation.md `[activity: transformation]`
            - [x] exploratory-testing.md, performance-testing.md `[activity: transformation]`
        - [x] Validate testing methodology preservation `[activity: review]`

    - [x] **Security Engineers** `[parallel: true]` `[component: the-security-engineer]`
        - [x] Transform 5 security agents:
            - [x] vulnerability-assessment.md, authentication-systems.md `[activity: transformation]`
            - [x] security-incident-response.md, compliance-audit.md `[activity: transformation]`
            - [x] data-protection.md `[activity: transformation]`
        - [x] Validate security patterns and FOCUS/EXCLUDE preservation `[activity: review]`

- [x] **Phase 5: Specialized Domains** (22 agents)

    - [x] **Designers** `[parallel: true]` `[component: the-designer]`
        - [x] Transform 6 designer agents:
            - [x] accessibility-implementation.md, user-research.md `[activity: transformation]`
            - [x] interaction-design.md, visual-design.md `[activity: transformation]`
            - [x] design-systems.md, information-architecture.md `[activity: transformation]`
        - [x] Validate UX pattern preservation `[activity: review]`

    - [x] **ML Engineers** `[parallel: true]` `[component: the-ml-engineer]`
        - [x] Transform 6 ML agents:
            - [x] model-deployment.md, ml-monitoring.md, prompt-optimization.md `[activity: transformation]`
            - [x] mlops-automation.md, context-management.md, feature-engineering.md `[activity: transformation]`
        - [x] Validate ML methodology patterns `[activity: review]`

    - [x] **Mobile Engineers** `[parallel: true]` `[component: the-mobile-engineer]`
        - [x] Transform 5 mobile agents:
            - [x] mobile-data-persistence.md, mobile-interface-design.md `[activity: transformation]`
            - [x] cross-platform-integration.md, mobile-deployment.md `[activity: transformation]`
            - [x] mobile-performance.md `[activity: transformation]`
        - [x] Validate mobile-specific patterns `[activity: review]`

    - [x] **Analysts** `[parallel: true]` `[component: the-analyst]`
        - [x] Transform 5 analyst agents:
            - [x] requirements-clarification.md, requirements-documentation.md `[activity: transformation]`
            - [x] feature-prioritization.md, solution-research.md `[activity: transformation]`
            - [x] project-coordination.md `[activity: transformation]`
        - [x] Validate analysis methodology preservation `[activity: review]`

- [x] **Phase 6: Final Validation & Integration**

    - [x] **System Integration Review**
        - [x] All 61 agents comply with SDD Enhanced Agent Template `[ref: SDD; lines: 236-289]`
        - [x] Delegation boundaries preserved across all domains
        - [x] 70% declarative content achieved per SDD requirements
        - [x] Pattern consistency verified across all agents

    - [x] **Quality Assurance**
        - [x] Capability preservation validated (100% functional parity)
        - [x] FOCUS/EXCLUDE patterns maintained
        - [x] Framework detection sections consistent
        - [x] Anti-patterns properly defined

    - [x] **Documentation & Completion**
        - [x] Pattern library documented for maintenance
        - [x] Transformation methodology captured
        - [x] Rollback procedures documented if needed
        - [x] Implementation complete per SDD specifications

---

## One-Shot Success Criteria

**Transformation Completeness**:
- [x] All 61 agents follow Enhanced Agent Template structure exactly
- [x] Each agent has rich description with multiple examples and commentary
- [x] All agents achieve 70% declarative content (Core Responsibilities + Output Format)
- [x] Every original capability preserved in transformed version

**Pattern Compliance**:
- [x] Consistent section ordering across all agents
- [x] Methodology phases replace procedural steps
- [x] Best Practices replace Anti-Patterns with positive framing
- [x] Closing mindset statement present in all agents

**Quality Validation**:
- [x] No procedural step-by-step instructions remain (except in methodology phases)
- [x] Clear delegation boundaries in descriptions
- [x] Framework detection sections preserved where applicable
- [x] FOCUS/EXCLUDE patterns maintained for security

**Ready for Use**:
- [x] All agents immediately usable without further modification
- [x] Delegation from the-chief functions correctly
- [x] Multi-agent workflows operate as expected
- [x] No breaking changes to existing integrations
