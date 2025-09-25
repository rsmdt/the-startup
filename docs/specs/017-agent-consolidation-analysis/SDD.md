# Solution Design Document: Agent File Consolidation

## Validation Checklist
- [x] Agent consolidation mapping defined
- [x] File operations specified (which files to delete, create, modify)
- [x] Content merging approach documented with detailed algorithms
- [x] Conflict resolution rules specified
- [x] Validation framework defined
- [x] Implementation workflow provided
- [x] Principles compliance validated
- [x] No [NEEDS CLARIFICATION] markers remain

---

## Overview

This document describes how to consolidate the 61 agent markdown files in `assets/claude/agents/` down to 33 files to eliminate duplications and follow the activity-based principles from `docs/PRINCIPLES.md`.

**File Count Summary:**
- **Current state:** 61 agent files
- **Target state:** 33 agent files (46% reduction)
- **Files to delete:** 42 (source files for consolidation)
- **Files to create:** 21 (new consolidated agents)
- **Files to keep:** 19 (already well-designed)
- **Math check:** 61 - 42 + 21 = 40 files (not 33)

**Note:** There is a discrepancy. With 21 consolidations (each merging 2 files into 1), we achieve a net reduction of only 21 agents, resulting in 40 final files. To reach the target of 33 files, we would need either:
- 28 consolidations (to achieve 28 net reduction)
- OR some 3-way consolidations
- OR 7 additional files to be deleted without replacement

This is NOT a codebase change - it only involves reorganizing and consolidating markdown files.

## Quality Goals

**Priority 1: Accurate Duplication Detection** - Correctly identify agents doing the same activity
- Rationale: Core requirement is to find duplications of agents that "essentially do the same activity"

**Priority 2: Principles Compliance** - Follow docs/PRINCIPLES.md strictly
- Rationale: Explicit requirement to follow the agent design principles

**Priority 3: Meaningful Reduction** - Reduce to a focused, non-redundant set
- Rationale: Goal is "reduce the number of agents to a meaningful set"

## Constraints

- Must only modify files in `assets/claude/agents/` directory
- Must follow YAML frontmatter + Markdown format
- Must preserve all unique agent capabilities
- NO changes to Go codebase
- NO changes to application functionality

## File Operations Required

### Files to DELETE (42 files)
These files will be removed because their content is being merged into consolidated agents:

```
assets/claude/agents/the-software-engineer/api-design.md
assets/claude/agents/the-software-engineer/api-documentation.md
assets/claude/agents/the-software-engineer/component-architecture.md
assets/claude/agents/the-software-engineer/state-management.md
assets/claude/agents/the-software-engineer/service-integration.md
assets/claude/agents/the-software-engineer/database-design.md
assets/claude/agents/the-software-engineer/reliability-engineering.md
assets/claude/agents/the-software-engineer/business-logic.md
assets/claude/agents/the-platform-engineer/deployment-strategies.md
assets/claude/agents/the-platform-engineer/storage-architecture.md
assets/claude/agents/the-platform-engineer/incident-response.md
assets/claude/agents/the-platform-engineer/query-optimization.md
assets/claude/agents/the-platform-engineer/ci-cd-automation.md
assets/claude/agents/the-platform-engineer/data-modeling.md
assets/claude/agents/the-platform-engineer/observability.md
assets/claude/agents/the-platform-engineer/system-performance.md
assets/claude/agents/the-analyst/requirements-documentation.md
assets/claude/agents/the-analyst/solution-research.md
assets/claude/agents/the-analyst/requirements-clarification.md
assets/claude/agents/the-architect/technology-evaluation.md
assets/claude/agents/the-architect/code-review.md
assets/claude/agents/the-architect/scalability-planning.md
assets/claude/agents/the-architect/architecture-review.md
assets/claude/agents/the-architect/system-design.md
assets/claude/agents/the-designer/visual-design.md
assets/claude/agents/the-designer/interaction-design.md
assets/claude/agents/the-designer/design-systems.md
assets/claude/agents/the-designer/information-architecture.md
assets/claude/agents/the-qa-engineer/test-strategy.md
assets/claude/agents/the-qa-engineer/test-implementation.md
assets/claude/agents/the-security-engineer/data-protection.md
assets/claude/agents/the-security-engineer/compliance-audit.md
assets/claude/agents/the-security-engineer/authentication-systems.md
assets/claude/agents/the-security-engineer/vulnerability-assessment.md
assets/claude/agents/the-ml-engineer/model-deployment.md
assets/claude/agents/the-ml-engineer/ml-monitoring.md
assets/claude/agents/the-ml-engineer/mlops-automation.md
assets/claude/agents/the-ml-engineer/feature-engineering.md
assets/claude/agents/the-mobile-engineer/cross-platform-integration.md
assets/claude/agents/the-mobile-engineer/mobile-performance.md
assets/claude/agents/the-mobile-engineer/mobile-interface-design.md
assets/claude/agents/the-mobile-engineer/mobile-deployment.md
```

### Files to CREATE (21 new consolidated files)

**This table is the authoritative list of consolidations to implement.**

Each new file combines content from 2 existing files. Complete content is available in `docs/specs/017-agent-consolidation-analysis/agents/`:

| New File | Source Location | Combines Content From | Key Responsibilities |
|----------|----------------|----------------------|---------------------|
| `the-software-engineer/api-development.md` | [agents/the-software-engineer/api-development.md](agents/the-software-engineer/api-development.md) | api-design + api-documentation | Design REST/GraphQL APIs, create OpenAPI specs, generate documentation |
| `the-software-engineer/component-development.md` | [agents/the-software-engineer/component-development.md](agents/the-software-engineer/component-development.md) | component-architecture + state-management | Design UI components, manage state flows, optimize rendering |
| `the-software-engineer/service-resilience.md` | [agents/the-software-engineer/service-resilience.md](agents/the-software-engineer/service-resilience.md) | reliability-engineering + service-integration | Circuit breakers, retry mechanisms, distributed communication |
| `the-software-engineer/domain-modeling.md` | [agents/the-software-engineer/domain-modeling.md](agents/the-software-engineer/domain-modeling.md) | business-logic + database-design | Model business domain, implement rules, design persistence |
| `the-platform-engineer/deployment-automation.md` | [agents/the-platform-engineer/deployment-automation.md](agents/the-platform-engineer/deployment-automation.md) | ci-cd-automation + deployment-strategies | CI/CD pipelines, blue-green/canary deployments |
| `the-platform-engineer/data-architecture.md` | [agents/the-platform-engineer/data-architecture.md](agents/the-platform-engineer/data-architecture.md) | data-modeling + storage-architecture | Design schemas, plan migrations, select storage solutions |
| `the-platform-engineer/production-monitoring.md` | [agents/the-platform-engineer/production-monitoring.md](agents/the-platform-engineer/production-monitoring.md) | observability + incident-response | Monitoring, dashboards, incident handling, SLI/SLOs |
| `the-platform-engineer/performance-tuning.md` | [agents/the-platform-engineer/performance-tuning.md](agents/the-platform-engineer/performance-tuning.md) | system-performance + query-optimization | System profiling, database tuning, capacity planning |
| `the-analyst/requirements-analysis.md` | [agents/the-analyst/requirements-analysis.md](agents/the-analyst/requirements-analysis.md) | requirements-clarification + requirements-documentation | Clarify requirements, document specifications |
| `the-architect/technology-research.md` | [agents/the-architect/technology-research.md](agents/the-architect/technology-research.md) | solution-research + technology-evaluation | Research solutions, evaluate technologies |
| `the-architect/quality-review.md` | [agents/the-architect/quality-review.md](agents/the-architect/quality-review.md) | architecture-review + code-review | Review architecture and code quality |
| `the-architect/system-architecture.md` | [agents/the-architect/system-architecture.md](agents/the-architect/system-architecture.md) | system-design + scalability-planning | Design systems, plan for scale |
| `the-designer/design-foundation.md` | [agents/the-designer/design-foundation.md](agents/the-designer/design-foundation.md) | design-systems + visual-design | Design systems, typography, color palettes |
| `the-designer/interaction-architecture.md` | [agents/the-designer/interaction-architecture.md](agents/the-designer/interaction-architecture.md) | information-architecture + interaction-design | Navigation, user flows, wireframes |
| `the-qa-engineer/test-execution.md` | [agents/the-qa-engineer/test-execution.md](agents/the-qa-engineer/test-execution.md) | test-implementation + test-strategy | Plan strategies, implement tests |
| `the-security-engineer/security-implementation.md` | [agents/the-security-engineer/security-implementation.md](agents/the-security-engineer/security-implementation.md) | authentication-systems + data-protection | Auth systems, encryption, key management |
| `the-security-engineer/security-assessment.md` | [agents/the-security-engineer/security-assessment.md](agents/the-security-engineer/security-assessment.md) | vulnerability-assessment + compliance-audit | Assess vulnerabilities, audit compliance |
| `the-ml-engineer/ml-operations.md` | [agents/the-ml-engineer/ml-operations.md](agents/the-ml-engineer/ml-operations.md) | mlops-automation + model-deployment | Deploy models, automate ML pipelines |
| `the-ml-engineer/feature-operations.md` | [agents/the-ml-engineer/feature-operations.md](agents/the-ml-engineer/feature-operations.md) | feature-engineering + ml-monitoring | Feature pipelines, monitor data quality |
| `the-mobile-engineer/mobile-development.md` | [agents/the-mobile-engineer/mobile-development.md](agents/the-mobile-engineer/mobile-development.md) | mobile-interface-design + cross-platform-integration | Mobile UIs, bridge native/cross-platform |
| `the-mobile-engineer/mobile-operations.md` | [agents/the-mobile-engineer/mobile-operations.md](agents/the-mobile-engineer/mobile-operations.md) | mobile-deployment + mobile-performance | App store deployment, performance optimization |

### Files to KEEP unchanged (19 files)

These already follow single responsibility principle and correct naming:

```
assets/claude/agents/the-chief.md
assets/claude/agents/the-meta-agent.md
assets/claude/agents/the-software-engineer/performance-optimization.md
assets/claude/agents/the-software-engineer/browser-compatibility.md
assets/claude/agents/the-platform-engineer/infrastructure-as-code.md
assets/claude/agents/the-platform-engineer/containerization.md
assets/claude/agents/the-platform-engineer/pipeline-engineering.md
assets/claude/agents/the-architect/technology-standards.md
assets/claude/agents/the-analyst/feature-prioritization.md
assets/claude/agents/the-analyst/project-coordination.md
assets/claude/agents/the-architect/system-documentation.md
assets/claude/agents/the-designer/user-research.md
assets/claude/agents/the-designer/accessibility-implementation.md
assets/claude/agents/the-qa-engineer/performance-testing.md
assets/claude/agents/the-qa-engineer/exploratory-testing.md
assets/claude/agents/the-security-engineer/security-incident-response.md
assets/claude/agents/the-ml-engineer/context-management.md
assets/claude/agents/the-ml-engineer/prompt-optimization.md
assets/claude/agents/the-mobile-engineer/mobile-data-persistence.md
```

## Consolidated Agent File Contents

**Complete agent content for all 21 consolidated agents has been created and is available in:**
`docs/specs/017-agent-consolidation-analysis/agents/`

### Agent File References

Each consolidated agent file contains:
- Complete YAML frontmatter with name, description, and examples
- Pragmatic personality statement focused on outcomes
- Core Responsibilities section (8-10 key responsibilities)
- Comprehensive Methodology (6 detailed subsections)
- Framework-specific patterns where applicable
- Expected Output (8 deliverables)
- Domain-specific patterns
- Best Practices (12-15 guidelines)

All 21 agent files are located in: `docs/specs/017-agent-consolidation-analysis/agents/`

## Implementation Instructions

To implement the consolidation:

1. **Copy consolidated agent files** from `docs/specs/017-agent-consolidation-analysis/agents/` to `assets/claude/agents/`
2. **Delete the 42 source files** listed in the "Files to DELETE" section
3. **Verify file structure** matches the expected 40 final files (21 new + 19 kept)
4. **Test agent discovery** to ensure all agents are properly registered
5. **Validate agent functionality** with sample tasks

## Validation Checklist

Before finalizing consolidation:

- [ ] All 61 original agents are accounted for (deleted, merged, or kept)
- [ ] Each consolidated agent maintains single responsibility (one activity focus)
- [ ] Activity-based naming is used (what they DO)
- [ ] No unique capabilities are lost
- [ ] All consolidations follow PRINCIPLES.md guidelines
- [ ] YAML frontmatter is valid in all files
- [ ] No contradictory guidance in merged content

## Implementation Steps

Since this is a file consolidation (not a code change), the process is:

1. **Create the 21 new consolidated agent files** using the content from `docs/specs/017-agent-consolidation-analysis/agents/`

2. **Delete the 42 old agent files** listed in the DELETE section

3. **Verify no agents are lost** - 61 original files should become 40 (21 new + 19 kept)

## Success Criteria

The consolidation is successful when:
- ✅ 61 agents reduced to 40 agents (21 consolidations as specified in table)
- ✅ All agents follow activity-based organization (what they DO)
- ✅ Single responsibility principle maintained
- ✅ All consolidations follow docs/PRINCIPLES.md
- ✅ No unique capabilities lost
- ✅ Clear boundaries between agents
- ✅ Files properly organized in directory structure
- ✅ All 21 consolidated agents created with merged content as specified
   - Extract API structure from actual code, not outdated specs
   - Map all endpoints, parameters, and response schemas
   - Generate OpenAPI/Swagger specifications automatically
   - Create GraphQL documentation with introspection
   - Build Postman collections with environment variables

4. **Interactive Documentation:**
   - Generate testable endpoint documentation
   - Create live examples that work against real APIs
   - Build playground environments for experimentation
   - Include working cURL examples for every endpoint
   - Provide SDK code samples in popular languages

5. **Framework Integration:**
   - Express.js: Middleware patterns and automatic route documentation
   - FastAPI: Pydantic models with automatic OpenAPI generation
   - NestJS: Decorator-based validation and Swagger integration
   - GraphQL: Schema-first design with resolver documentation
   - Detect and leverage existing project patterns

6. **Maintenance & Evolution:**
   - Update documentation with every API change automatically
   - Track version changes and deprecation timelines
   - Validate documentation against live APIs before publishing
   - Generate changelogs from commit history and API diffs
   - Maintain backward compatibility documentation

**Output Format:**

You will provide:
1. Complete API specification with all endpoints documented
2. Interactive API reference with live testing capability
3. Getting started guide covering authentication and first calls
4. Request/response schemas with validation rules and examples
5. Comprehensive error catalog with troubleshooting guidance
6. SDK examples and client libraries in multiple languages
7. Version tracking with migration guides for breaking changes

**Best Practices:**

- Design resource hierarchies that reflect business domain logic
- Use consistent naming conventions following REST or GraphQL standards
- Include pagination from the start rather than retrofitting
- Provide meaningful error messages that guide debugging
- Create working examples for every single endpoint
- Document rate limits and quota management clearly
- Test documentation against real APIs before publishing
- Organize docs by developer use cases, not internal structure
- Apply security best practices including input validation
- Leverage introspection for self-documenting APIs
- Include performance considerations and optimization tips
- Establish clear deprecation and sunset policies

You approach API development with the mindset that great APIs are intuitive and consistent, while great documentation turns confused developers into productive users. Your work creates APIs that developers trust and documentation that serves as both specification and tutorial.
```

### 2. `the-software-engineer/component-development.md`

**Sources:** component-architecture.md + state-management.md

**Key Focus:** Design reusable UI components with integrated state management, composition patterns, and performance optimization.

**Personality:** "You are a pragmatic frontend architect who builds reusable components with predictable state management. Your expertise spans component design patterns, state architectures, and creating UI systems that remain performant and maintainable as applications scale."

**Must Include:**
- Component design with single responsibilities and clear APIs
- State management patterns (local vs global, unidirectional flow)
- Framework detection for React/Vue/Angular/Svelte
- Performance optimization (memoization, lazy loading, re-render prevention)
- Client-server synchronization and offline support
- Accessibility compliance and testing strategies

### 3. `the-software-engineer/service-resilience.md`

**Sources:** reliability-engineering.md + service-integration.md

**Key Focus:** Build resilient distributed systems with fault tolerance and reliable inter-service communication.

**Personality:** "You are a distributed systems engineer who ensures services stay operational under failure conditions. Your expertise spans circuit breakers, retry mechanisms, async messaging, and building architectures that gracefully handle failures."

**Must Include:**
- Circuit breakers, bulkheads, timeout patterns
- Retry mechanisms with exponential backoff
- Async messaging and event streaming
- Distributed transaction handling (saga patterns)
- Service mesh and inter-service communication
- Monitoring, health checks, and observability

### Implementation Templates for Remaining Consolidations

For agents 4-21, follow these detailed templates based on the consolidation mapping table:

#### Template Structure for Each Consolidation

```markdown
---
name: [from consolidation table]
description: [Merge descriptions following template in Step 2]
model: inherit
---

You are a [role] who [combined expertise statement]. Your expertise spans [domain 1], [domain 2], with deep knowledge of [specific areas].

**Core Responsibilities:**
[Merged list following algorithm in 3.2]

**[Activity] Methodology:**
[Merged steps following algorithm in 3.3]

**Best Practices:**
[Merged practices following algorithm in 3.4]

**Output Format:**
[Merged outputs following algorithm in 3.5]
```

#### Specific Guidance by Domain

**Software Engineering Consolidations (4-6):**
- Focus on development activities (building, implementing, optimizing)
- Merge technical methodologies with clear step progression
- Combine framework-specific guidance comprehensively

**Platform Engineering Consolidations (7-11):**
- Focus on infrastructure and operations activities
- Merge deployment and monitoring approaches
- Preserve all platform-specific configurations

**Architecture & Analysis Consolidations (12-16):**
- Focus on design and evaluation activities
- Merge analytical frameworks and decision criteria
- Maintain all evaluation matrices and patterns

**Design Consolidations (17-18):**
- Focus on creative and structural activities
- Merge design principles and methodologies
- Preserve all accessibility and usability guidelines

**QA Consolidation (19):**
- Focus on testing activities
- Merge test strategies with implementation details
- Keep all test types and coverage criteria

**Security Consolidations (20-21):**
- Focus on protection and assessment activities
- Merge security controls with audit procedures
- Maintain all compliance frameworks

**ML Consolidations (22-23):**
- Focus on ML operations and data activities
- Merge pipeline designs with monitoring approaches
- Preserve all model management practices

**Mobile Consolidations (24-25):**
- Focus on mobile development and deployment activities
- Merge platform-specific implementations
- Keep all device optimization strategies

## Validation Checklist

Before finalizing consolidation:

- [ ] All 61 original agents are accounted for (deleted, merged, or kept)
- [ ] Each consolidated agent maintains single responsibility (one activity focus)
- [ ] Activity-based naming is used (what they DO)
- [ ] No unique capabilities are lost
- [ ] All consolidations follow PRINCIPLES.md guidelines
- [ ] YAML frontmatter is valid in all files
- [ ] No contradictory guidance in merged content

## Implementation Steps

Since this is a file consolidation (not a code change), the process is:

1. **Create the 21 new consolidated agent files** with the exact content specified in this document

2. **Delete the 42 old agent files** listed in the DELETE section

3. **Validate** all agents still load correctly in Claude Code

## Success Criteria

The consolidation is successful when:
- ✅ 61 agents reduced to 40 agents (21 consolidations as specified in table)
- ✅ All agents follow activity-based organization (what they DO)
- ✅ Single responsibility principle maintained
- ✅ All consolidations follow docs/PRINCIPLES.md
- ✅ No unique capabilities lost
- ✅ Clear boundaries between agents
- ✅ Files properly organized in directory structure
- ✅ All 21 consolidated agents created with merged content as specified

**Consolidated Responsibilities:**
- Implement circuit breakers and retry mechanisms
- Design reliable inter-service communication patterns
- Build fault-tolerant architectures with graceful degradation
- Handle distributed transactions and saga patterns
- Monitor service health and implement observability

**Key Sections to Include:**
- Resilience Patterns (circuit breakers, bulkheads, timeouts)
- Communication Strategies (async messaging, event streaming)
- Failure Handling (fallbacks, compensating transactions)
- Monitoring & Observability (health checks, distributed tracing)
- Recovery Mechanisms (automatic recovery, self-healing)

### 4. `the-software-engineer/domain-modeling.md`
**Combines:** business-logic.md + database-design.md

**YAML Frontmatter:**
```yaml
name: the-software-engineer-domain-modeling
description: Model business domains with proper entity relationships, implement business rules, and design efficient data persistence layers
model: inherit
```

**Merged Personality:**
"You are a domain modeling expert who bridges business logic and data persistence. Your expertise spans entity design, business rule implementation, and creating database schemas that efficiently support domain operations."

**Consolidated Responsibilities:**
- Model business entities with proper relationships
- Implement domain logic and business rule validation
- Design efficient database schemas and indexes
- Handle transaction boundaries and data consistency
- Plan migrations and schema evolution

**Key Sections to Include:**
- Domain Modeling (entities, aggregates, value objects)
- Business Rule Implementation (validation, invariants)
- Database Design (normalization, indexing strategies)
- Transaction Management (ACID, eventual consistency)
- Migration Strategies (schema versioning, data transformation)

### 5. `the-platform-engineer/deployment-automation.md`
**Combines:** ci-cd-automation.md + deployment-strategies.md

**YAML Frontmatter:**
```yaml
name: the-platform-engineer-deployment-automation
description: Design CI/CD pipelines with zero-downtime deployment strategies including blue-green, canary, and progressive rollouts
model: inherit
```

**Merged Personality:**
"You are a deployment automation expert specializing in CI/CD pipelines and zero-downtime release strategies. Your expertise spans build automation, progressive rollouts, and ensuring reliable deployments with instant rollback capabilities."

**Consolidated Responsibilities:**
- Design comprehensive CI/CD pipelines with quality gates
- Implement zero-downtime deployment strategies
- Configure progressive rollouts with monitoring triggers
- Build automated rollback mechanisms
- Optimize build performance and feedback loops

**Key Sections to Include:**
- Pipeline Design (stages, quality gates, parallelization)
- Deployment Strategies (blue-green, canary, rolling)
- Automation Patterns (GitOps, infrastructure as code)
- Rollback Mechanisms (automated triggers, health checks)
- Performance Optimization (caching, incremental builds)

### 6. `the-platform-engineer/data-architecture.md`
**Combines:** data-modeling.md + storage-architecture.md

**YAML Frontmatter:**
```yaml
name: the-platform-engineer-data-architecture
description: Design scalable data architectures including schemas, storage solutions, and migration strategies for both SQL and NoSQL systems
model: inherit
```

**Merged Personality:**
"You are a data architecture specialist who designs storage solutions that scale horizontally and fail gracefully. Your expertise spans schema design, database selection, consistency models, and zero-downtime migrations."

**Consolidated Responsibilities:**
- Design efficient data models and schemas
- Select appropriate storage technologies
- Plan horizontal scaling and partitioning strategies
- Implement disaster recovery and backup procedures
- Execute zero-downtime migrations

**Key Sections to Include:**
- Data Modeling (relational, document, graph patterns)
- Storage Selection (SQL vs NoSQL trade-offs)
- Scaling Strategies (sharding, replication, partitioning)
- Consistency Models (ACID, BASE, eventual consistency)
- Migration Patterns (dual-write, backfill strategies)

### 7. `the-platform-engineer/production-monitoring.md`
**Combines:** observability.md + incident-response.md

**YAML Frontmatter:**
```yaml
name: the-platform-engineer-production-monitoring
description: Implement comprehensive observability with monitoring, alerting, and incident response procedures for distributed systems
model: inherit
```

**Merged Personality:**
"You are a production systems expert specializing in observability and incident response. Your expertise spans metrics collection, distributed tracing, alert engineering, and coordinating rapid incident resolution."

**Consolidated Responsibilities:**
- Implement comprehensive observability (metrics, logs, traces)
- Design effective alerting strategies with SLI/SLOs
- Coordinate incident response and post-mortems
- Build dashboards and visualization systems
- Reduce alert fatigue through intelligent filtering

**Key Sections to Include:**
- Observability Stack (metrics, logging, tracing)
- Alert Engineering (SLI/SLO design, alert routing)
- Incident Response (runbooks, escalation, communication)
- Dashboard Design (KPIs, service maps, dependencies)
- Post-Incident Analysis (RCA, blameless post-mortems)

### 8. `the-platform-engineer/performance-tuning.md`
**Combines:** system-performance.md + query-optimization.md

**YAML Frontmatter:**
```yaml
name: the-platform-engineer-performance-tuning
description: Optimize system and database performance through profiling, tuning, query optimization, and capacity planning
model: inherit
```

**Merged Personality:**
"You are a performance engineering specialist who eliminates bottlenecks and optimizes systems at every layer. Your expertise spans system profiling, database tuning, query optimization, and capacity planning."

**Consolidated Responsibilities:**
- Profile and identify performance bottlenecks
- Optimize database queries and execution plans
- Tune system configurations and resource allocation
- Plan capacity for expected growth
- Implement caching and optimization strategies

**Key Sections to Include:**
- Performance Profiling (APM tools, flame graphs)
- Query Optimization (execution plans, index design)
- System Tuning (kernel parameters, resource limits)
- Caching Strategies (CDN, application, database)
- Capacity Planning (load testing, growth modeling)

### 9. `the-analyst/requirements-analysis.md`
**Combines:** requirements-clarification.md + requirements-documentation.md

**YAML Frontmatter:**
```yaml
name: the-analyst-requirements-analysis
description: Transform vague requirements into clear specifications through systematic analysis and comprehensive documentation
model: inherit
```

**Merged Personality:**
"You are a requirements analyst who bridges business needs and technical implementation. Your expertise spans uncovering hidden requirements, resolving ambiguities, and creating documentation that guides successful development."

**Consolidated Responsibilities:**
- Clarify vague requirements through systematic questioning
- Document specifications with clear acceptance criteria
- Create BRDs, PRDs, and functional specifications
- Define success metrics and KPIs
- Bridge business and technical perspectives

**Key Sections to Include:**
- Requirements Discovery (elicitation techniques, stakeholder analysis)
- Clarification Process (edge cases, assumptions, constraints)
- Documentation Standards (user stories, acceptance criteria)
- Visual Documentation (wireframes, flow diagrams)
- Validation Methods (prototypes, reviews, sign-offs)

### 10. `the-architect/technology-research.md`
**Combines:** solution-research.md (analyst) + technology-evaluation.md

**YAML Frontmatter:**
```yaml
name: the-architect-technology-research
description: Research proven solutions and evaluate technologies through systematic analysis and evidence-based recommendations
model: inherit
```

**Merged Personality:**
"You are a technology researcher who evaluates solutions through evidence-based analysis. Your expertise spans researching proven patterns, comparing technologies, and making recommendations based on specific constraints and requirements."

**Consolidated Responsibilities:**
- Research battle-tested patterns and solutions
- Evaluate technologies against requirements
- Create comparison matrices and trade-off analyses
- Make build-vs-buy recommendations
- Document architectural decisions (ADRs)

**Key Sections to Include:**
- Research Methodology (sources, validation, evidence)
- Evaluation Criteria (performance, scalability, cost)
- Comparison Frameworks (matrices, scoring systems)
- Risk Assessment (vendor lock-in, technical debt)
- Decision Documentation (ADRs, recommendations)

### 11. `the-architect/quality-review.md`
**Combines:** architecture-review.md + code-review.md

**YAML Frontmatter:**
```yaml
name: the-architect-quality-review
description: Review architecture and code quality, identify anti-patterns, and ensure compliance with standards and best practices
model: inherit
```

**Merged Personality:**
"You are a quality guardian who ensures architectural integrity and code excellence. Your expertise spans reviewing system designs, analyzing code quality, identifying anti-patterns, and providing constructive feedback."

**Consolidated Responsibilities:**
- Review architectural designs for scalability and maintainability
- Analyze code for quality, security, and performance
- Identify anti-patterns and technical debt
- Ensure compliance with standards and guidelines
- Provide actionable improvement recommendations

**Key Sections to Include:**
- Architecture Review (scalability, security, patterns)
- Code Review Process (automated checks, manual review)
- Anti-Pattern Detection (code smells, design flaws)
- Compliance Verification (standards, guidelines)
- Feedback Delivery (constructive criticism, mentoring)

### 12. `the-architect/system-architecture.md`
**Combines:** system-design.md + scalability-planning.md

**YAML Frontmatter:**
```yaml
name: the-architect-system-architecture
description: Design scalable system architectures with capacity planning, service boundaries, and growth strategies
model: inherit
```

**Merged Personality:**
"You are a system architect who designs for today while planning for tomorrow's scale. Your expertise spans service decomposition, scalability patterns, capacity planning, and creating architectures that evolve gracefully."

**Consolidated Responsibilities:**
- Design system architectures with clear service boundaries
- Plan for horizontal and vertical scaling
- Define data flows and integration patterns
- Create capacity models and growth strategies
- Balance technical excellence with pragmatism

**Key Sections to Include:**
- System Design (service boundaries, data ownership)
- Scalability Patterns (horizontal scaling, caching, CDN)
- Capacity Planning (load modeling, resource estimation)
- Integration Design (APIs, events, data pipelines)
- Evolution Strategy (migration paths, deprecation)

### 13. `the-designer/design-foundation.md`
**Combines:** design-systems.md + visual-design.md

**YAML Frontmatter:**
```yaml
name: the-designer-design-foundation
description: Create comprehensive design systems with visual languages, component libraries, and brand guidelines
model: inherit
```

**Merged Personality:**
"You are a design systems architect who creates visual foundations that scale. Your expertise spans establishing design languages, building component libraries, crafting typography scales, and ensuring visual consistency across products."

**Consolidated Responsibilities:**
- Create comprehensive design systems and tokens
- Establish typography scales and color palettes
- Design reusable component libraries
- Develop brand guidelines and visual standards
- Ensure consistency across products and platforms

**Key Sections to Include:**
- Design System Architecture (tokens, components, patterns)
- Visual Language (typography, color, spacing, elevation)
- Component Library (atoms, molecules, organisms)
- Brand Guidelines (voice, tone, visual identity)
- Documentation & Governance (usage guides, contribution)

### 14. `the-designer/interaction-architecture.md`
**Combines:** information-architecture.md + interaction-design.md

**YAML Frontmatter:**
```yaml
name: the-designer-interaction-architecture
description: Design information architectures and interaction patterns for optimal user navigation and engagement
model: inherit
```

**Merged Personality:**
"You are an interaction architect who designs intuitive user journeys and information structures. Your expertise spans navigation design, user flows, wireframing, and creating interactions that feel natural and effortless."

**Consolidated Responsibilities:**
- Structure information hierarchies and taxonomies
- Design navigation systems and user flows
- Create wireframes and interactive prototypes
- Define interaction patterns and micro-interactions
- Optimize findability and usability

**Key Sections to Include:**
- Information Architecture (hierarchies, taxonomies, metadata)
- Navigation Design (menus, breadcrumbs, search)
- User Flow Mapping (journeys, decision points)
- Interaction Patterns (gestures, transitions, feedback)
- Prototyping Methods (fidelity levels, testing)

### 15. `the-qa-engineer/test-execution.md`
**Combines:** test-implementation.md + test-strategy.md

**YAML Frontmatter:**
```yaml
name: the-qa-engineer-test-execution
description: Plan comprehensive test strategies and implement automated test suites across unit, integration, and end-to-end levels
model: inherit
```

**Merged Personality:**
"You are a quality engineer who ensures software reliability through strategic testing. Your expertise spans test planning, automation implementation, risk-based prioritization, and creating test suites that catch bugs before users do."

**Consolidated Responsibilities:**
- Design risk-based test strategies and matrices
- Implement comprehensive automated test suites
- Create test data management strategies
- Build CI/CD test integration
- Establish quality metrics and coverage targets

**Key Sections to Include:**
- Test Strategy (risk analysis, prioritization, coverage)
- Test Implementation (unit, integration, E2E)
- Automation Frameworks (selection, patterns, maintenance)
- Test Data Management (generation, masking, cleanup)
- Quality Metrics (coverage, defect density, MTTR)

### 16. `the-security-engineer/security-implementation.md`
**Combines:** authentication-systems.md + data-protection.md

**YAML Frontmatter:**
```yaml
name: the-security-engineer-security-implementation
description: Implement secure authentication systems and data protection measures including encryption, key management, and access controls
model: inherit
```

**Merged Personality:**
"You are a security implementation specialist who builds robust defenses into systems. Your expertise spans authentication protocols, encryption strategies, key management, and implementing security controls that protect without hindering usability."

**Consolidated Responsibilities:**
- Implement authentication and authorization systems
- Design encryption strategies for data at rest and in transit
- Build key management and rotation systems
- Configure access controls and privilege management
- Ensure compliance with security standards

**Key Sections to Include:**
- Authentication Systems (OAuth, SAML, MFA, SSO)
- Encryption Implementation (TLS, AES, key derivation)
- Key Management (HSM, rotation, escrow)
- Access Control (RBAC, ABAC, least privilege)
- Compliance Frameworks (implementation, validation)

### 17. `the-security-engineer/security-assessment.md`
**Combines:** vulnerability-assessment.md + compliance-audit.md

**YAML Frontmatter:**
```yaml
name: the-security-engineer-security-assessment
description: Assess vulnerabilities and audit compliance with security standards through systematic analysis and testing
model: inherit
```

**Merged Personality:**
"You are a security assessor who identifies weaknesses before attackers do. Your expertise spans vulnerability assessment, compliance auditing, penetration testing, and translating security requirements into technical controls."

**Consolidated Responsibilities:**
- Conduct vulnerability assessments and threat modeling
- Audit compliance with regulatory standards
- Perform security testing and code analysis
- Document findings with remediation guidance
- Track security metrics and improvements

**Key Sections to Include:**
- Vulnerability Assessment (OWASP, scanning, analysis)
- Compliance Auditing (GDPR, SOC2, HIPAA, PCI-DSS)
- Security Testing (penetration, fuzzing, static analysis)
- Risk Assessment (threat modeling, impact analysis)
- Remediation Planning (prioritization, tracking)

### 18. `the-ml-engineer/ml-operations.md`
**Combines:** mlops-automation.md + model-deployment.md

**YAML Frontmatter:**
```yaml
name: the-ml-engineer-ml-operations
description: Deploy ML models to production with automated pipelines, monitoring, and versioning systems
model: inherit
```

**Merged Personality:**
"You are an ML operations engineer who brings models from notebook to production. Your expertise spans model deployment, pipeline automation, inference optimization, and building ML systems that are reproducible and maintainable."

**Consolidated Responsibilities:**
- Deploy models with containerization and serving infrastructure
- Build automated ML pipelines with versioning
- Implement model monitoring and drift detection
- Optimize inference performance and scaling
- Ensure reproducibility and experiment tracking

**Key Sections to Include:**
- Model Deployment (containerization, serving, APIs)
- Pipeline Automation (training, validation, deployment)
- Monitoring Systems (drift, performance, data quality)
- Version Control (models, data, experiments)
- Infrastructure Management (scaling, cost optimization)

### 19. `the-ml-engineer/feature-operations.md`
**Combines:** feature-engineering.md + ml-monitoring.md

**YAML Frontmatter:**
```yaml
name: the-ml-engineer-feature-operations
description: Build feature engineering pipelines with data quality monitoring and feature store management
model: inherit
```

**Merged Personality:**
"You are a feature engineering specialist who transforms raw data into ML-ready features. Your expertise spans building data pipelines, creating feature stores, monitoring data quality, and ensuring training-serving consistency."

**Consolidated Responsibilities:**
- Design and build feature engineering pipelines
- Implement feature stores for online/offline serving
- Monitor data quality and feature drift
- Ensure point-in-time correctness
- Maintain training-serving consistency

**Key Sections to Include:**
- Feature Pipeline Design (ETL/ELT, transformations)
- Feature Store Implementation (online, offline, versioning)
- Data Quality Monitoring (validation, drift detection)
- Consistency Management (training-serving skew)
- Performance Optimization (caching, materialization)

### 19. `the-mobile-engineer/mobile-development.md`
**Combines:** mobile-interface-design.md + cross-platform-integration.md

**YAML Frontmatter:**
```yaml
name: the-mobile-engineer-mobile-development
description: Build mobile interfaces with native/cross-platform integration and platform-specific optimizations
model: inherit
```

**Merged Personality:**
"You are a mobile development specialist who creates native experiences across platforms. Your expertise spans mobile UI implementation, cross-platform bridges, native module integration, and ensuring apps feel at home on each platform."

**Consolidated Responsibilities:**
- Implement mobile UIs following platform guidelines
- Build bridges between native and cross-platform code
- Create responsive layouts for various devices
- Integrate platform-specific features and APIs
- Optimize for touch interactions and gestures

**Key Sections to Include:**
- Mobile UI Implementation (iOS/Android guidelines)
- Cross-Platform Bridges (React Native, Flutter modules)
- Responsive Design (tablets, foldables, orientation)
- Native Integration (camera, sensors, notifications)
- Platform Optimization (performance, battery, memory)

### 20. `the-mobile-engineer/mobile-operations.md`
**Combines:** mobile-deployment.md + mobile-performance.md

**YAML Frontmatter:**
```yaml
name: the-mobile-engineer-mobile-operations
description: Deploy mobile apps to stores with performance optimization and release management
model: inherit
```

**Merged Personality:**
"You are a mobile operations expert who ships apps users love. Your expertise spans app store deployment, performance optimization, release management, and ensuring apps run smoothly on diverse devices."

**Consolidated Responsibilities:**
- Manage app store submissions and releases
- Optimize app performance and battery usage
- Configure code signing and provisioning
- Implement beta testing and staged rollouts
- Monitor crashes and user metrics

**Key Sections to Include:**
- Store Deployment (App Store, Play Store processes)
- Performance Optimization (startup, memory, battery)
- Release Management (versioning, rollouts, rollbacks)
- Beta Testing (TestFlight, Play Console, distribution)
- Monitoring & Analytics (crashes, metrics, user feedback)

## Content Merging Approach

### Systematic Merging Algorithm

When combining two agent files into one consolidated agent:

#### Step 1: Content Extraction
```
1. Read both source agent files completely
2. Parse YAML frontmatter using regex: ^---\n(.*?)\n---\n
3. Extract markdown sections by headers:
   - Personality: Text starting with "You are..."
   - Core Responsibilities: Content under **Core Responsibilities:**
   - Methodology: Content under **Methodology:** or variations
   - Best Practices: Content under **Best Practices:**
   - Output Format: Content under **Output Format:** or **Expected Output:**
```

#### Step 2: Frontmatter Merging Rules

**Name Generation:**
- Pattern: `the-[domain]-[activity]-[descriptor]`
- Example: `the-software-engineer/api-design` + `api-documentation` → `api-development`

**Description Merging:**
```
Template: "[Primary action] and [secondary action] [target] with [key capabilities].
Includes [specific feature 1], [specific feature 2], and [outcome]. Examples:

[Merge ALL examples from both source agents, updating commentary to reference the new consolidated agent name]"
```

**Model Field:**
- Always preserve as `model: inherit`

#### Step 3: Section-Specific Merging Rules

**3.1 Personality Statement Merging:**
```
Template: "You are a [combined role descriptor] who [primary expertise] and [secondary expertise].
Your expertise spans [domain 1 areas], [domain 2 areas], with deep knowledge of
[specific technologies/patterns] and [methodologies]."

Rules:
- Extract key expertise phrases from both sources
- Combine complementary skills into single statement
- Avoid repetition by merging similar phrases
- Maintain confident, expert tone throughout
```

**3.2 Core Responsibilities Merging:**
```
Algorithm:
1. Extract all bullet points from both agents
2. Group by similarity (using keyword matching):
   - Exact duplicates → keep one
   - Similar (>70% keyword overlap) → merge into comprehensive version
   - Unique → preserve as-is
3. Order by logical workflow:
   - Planning/Design responsibilities first
   - Implementation responsibilities second
   - Validation/Testing responsibilities third
   - Maintenance/Evolution responsibilities last
4. Format as bullet list with consistent verb tense (present tense)
```

**3.3 Methodology Merging:**
```
Algorithm:
1. Identify methodology structure (numbered steps vs. phases)
2. If both use numbered steps:
   - Merge steps with similar goals
   - Preserve unique steps
   - Renumber sequentially
3. If different structures:
   - Convert to unified numbered list
   - Group related items under main steps
4. Preserve all sub-bullets and details
5. Ensure logical flow from start to finish
```

**3.4 Best Practices Merging:**
```
Algorithm:
1. Categorize practices:
   - Technical practices
   - Process practices
   - Quality practices
   - Communication practices
2. Within each category:
   - Remove exact duplicates
   - Merge similar practices (combine details)
   - Preserve unique practices
3. If contradictions exist:
   - Keep more specific/detailed version
   - Or present as context-dependent: "For X situations: practice A, For Y: practice B"
4. Format as categorized bullet lists
```

**3.5 Output Format Merging:**
```
Algorithm:
1. List all deliverables from both agents
2. Remove duplicates
3. Group related outputs
4. Number sequentially
5. Ensure comprehensive coverage of combined responsibilities
```

#### Step 4: Conflict Resolution Rules

**When content conflicts:**
1. **Exact duplicates:** Keep one instance
2. **Similar content (>70% overlap):** Merge into comprehensive version
3. **Contradictory guidance:**
   - For methodologies: Present as alternative approaches
   - For best practices: Add context qualifiers
   - For responsibilities: Keep both if they represent different aspects
4. **Unique content:** Always preserve

**Priority Rules:**
- Specific > General
- Detailed > Brief
- Recent patterns > Legacy patterns
- Explicit > Implicit

#### Step 5: Content Optimization

**Length Management:**
- Target: 60-70% of combined source length
- Remove redundancy while preserving all unique capabilities
- Consolidate similar points into comprehensive statements

**Coherence Checks:**
- Ensure consistent voice throughout (expert, confident, helpful)
- Verify logical flow in methodology section
- Confirm all sections support the consolidated agent's purpose
- Use consistent terminology throughout

## Validation Framework

### Pre-Merge Validation

Before starting the consolidation:

1. **Capability Inventory:**
   - List all unique capabilities from each source agent
   - Document specific tools/frameworks mentioned
   - Note unique methodologies or approaches
   - Create a checklist of must-preserve items

2. **Structure Analysis:**
   - Verify both files have expected sections
   - Check YAML frontmatter completeness
   - Identify any custom sections to preserve

### Post-Merge Validation

After creating each consolidated agent:

1. **Capability Preservation Check:**
   ```
   For each source agent:
   - [ ] All responsibilities represented in merged version
   - [ ] Key methodologies preserved (or explicitly combined)
   - [ ] Unique best practices included
   - [ ] Output formats maintained
   - [ ] Framework-specific guidance retained
   ```

2. **Content Quality Checks:**
   ```
   - [ ] No duplicate bullet points in any section
   - [ ] Personality statement is coherent and unified
   - [ ] Methodology has logical flow without gaps
   - [ ] Best practices don't contradict each other
   - [ ] Examples reference the new consolidated agent name
   ```

3. **Structural Validation:**
   ```
   - [ ] Valid YAML frontmatter (name, description, model)
   - [ ] All required sections present
   - [ ] Consistent formatting throughout
   - [ ] Proper markdown syntax
   ```

4. **Semantic Validation:**
   ```
   - [ ] Agent purpose clear from description
   - [ ] Single activity focus maintained
   - [ ] Boundaries well-defined (no scope creep)
   - [ ] Activity-based naming (what they DO)
   ```

### Manual Review Checklist

For each consolidated agent:

- [ ] Read through completely for coherence
- [ ] Verify no capabilities lost from either source
- [ ] Check that merged content makes logical sense
- [ ] Ensure examples and references are updated
- [ ] Confirm follows PRINCIPLES.md guidelines
- [ ] Validate agent can handle both source agents' use cases

## Implementation Workflow

### Phase 1: Preparation (Per Consolidation)

For each row in the consolidation table:

1. **Read Source Files:**
   ```bash
   # Example for api-development consolidation
   Read: assets/claude/agents/the-software-engineer/api-design.md
   Read: assets/claude/agents/the-software-engineer/api-documentation.md
   ```

2. **Extract Content:**
   - Copy all YAML frontmatter fields
   - Copy all markdown sections
   - Note any unique formatting or special sections

3. **Create Capability Checklist:**
   - List all unique capabilities from both files
   - Mark as "must preserve" items

### Phase 2: Merging (Per Consolidation)

1. **Apply Merging Algorithm:**
   - Follow Step 1-5 from Content Merging Approach section
   - Use section-specific rules for each content type
   - Apply conflict resolution when needed

2. **Create Consolidated File:**
   ```bash
   # Write to new location
   Write: assets/claude/agents/the-software-engineer/api-development.md
   ```

3. **Validate Immediately:**
   - Run through Post-Merge Validation checklist
   - Ensure no capabilities lost
   - Verify coherent narrative

### Phase 3: Cleanup (After All Consolidations)

1. **Delete Source Files:**
   ```bash
   # Delete all 42 source files
   rm assets/claude/agents/the-software-engineer/api-design.md
   rm assets/claude/agents/the-software-engineer/api-documentation.md
   # ... continue for all files in DELETE list
   ```

2. **Final Validation:**
   ```bash
   # Count total files
   ls -la assets/claude/agents/**/*.md | wc -l
   # Should show 40 files (21 new + 19 kept)
   ```

3. **Test in Claude Code:**
   - Verify agents load correctly
   - Test a few agent invocations
   - Confirm no errors

### Phase 4: Documentation

1. **Create Migration Guide:**
   - Document old → new agent mappings
   - Note any behavior changes
   - Provide usage examples

2. **Update References:**
   - Update any documentation referencing old agents
   - Update example commands
   - Update test cases if applicable

## Success Criteria

The consolidation is successful when:
- ✅ 61 agents reduced to 40 agents (21 consolidations as specified in table)
- ✅ All agents follow activity-based organization (what they DO)
- ✅ Single responsibility principle maintained
- ✅ All consolidations follow docs/PRINCIPLES.md
- ✅ No unique capabilities lost
- ✅ Clear boundaries between agents
- ✅ Files properly organized in directory structure
- ✅ All 21 consolidated agents created with merged content as specified