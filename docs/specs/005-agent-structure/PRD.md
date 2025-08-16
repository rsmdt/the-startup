# Product Requirements Document

## Product Overview

### Vision
Transform Claude Code sub-agents into a best-in-class, standardized system that maximizes effectiveness through consistent patterns, clear personalities, and intelligent task decomposition.

### Problem Statement
Based on community feedback and research, key issues include: agents often aren't invoked even when explicitly created (GitHub discussions), inconsistent triggering requiring manual "use subagent" prompts (Medium articles), lack of standardization across community implementations (3+ competing GitHub repositories), and no clear best practices from Anthropic beyond basic documentation. Internal analysis shows 14 agents with varying quality (6-9/10 rating), duplicate content, and missing decomposition patterns.

### Value Proposition
A research-backed, standardized agent system following Anthropic's documented best practices (CLAUDE.md integration, TDD workflow, proactive use patterns) that solves the community's #1 issue of agent non-invocation through clear triggers and consistent structure, enabling the 70% performance improvement demonstrated in multi-agent research systems.

## User Personas

### Primary Persona: The Power Developer
- **Demographics:** 25-45, senior developer/architect, highly tech-savvy
- **Goals:** Leverage AI to accelerate complex development tasks, maintain code quality, automate repetitive work
- **Pain Points:** Inconsistent agent behavior, unclear when to use which agent, agents missing context from previous interactions
- **User Story:** As a power developer, I want consistent, predictable agent behavior so that I can confidently delegate complex tasks

### Secondary Personas

**The Team Lead**
- **Demographics:** 30-50, engineering manager, moderate-high tech-savvy
- **Goals:** Standardize team workflows, ensure code quality, accelerate delivery
- **Pain Points:** Agents produce varying quality outputs, difficult to train team on agent usage
- **User Story:** As a team lead, I want standardized agent outputs so that my team can collaborate effectively

**The Solo Founder**
- **Demographics:** 22-40, startup founder, moderate tech-savvy
- **Goals:** Build MVP quickly, maintain professional standards, scale efficiently
- **Pain Points:** Agents don't understand full project context, require too much guidance
- **User Story:** As a solo founder, I want agents that understand my entire project so that I can focus on business logic

## User Journey Maps

### Agent Selection Journey
1. **Awareness:** User recognizes need for specialized help (error, new feature, documentation)
2. **Consideration:** Reviews agent descriptions and examples to find right specialist
3. **Action:** Invokes agent with clear context and requirements
4. **Retention:** Agent maintains context, produces consistent outputs, enables smooth handoffs

### Complex Task Decomposition Journey
1. **Awareness:** User faces multi-faceted problem requiring various expertise
2. **Consideration:** Chief agent assesses complexity and recommends workflow
3. **Action:** Parallel agents work on independent aspects, results consolidated
4. **Retention:** Documentation preserved, patterns reusable, knowledge captured

## Feature Requirements

### Feature Set 1: Core Agent Structure
| Feature | User Story | Priority | Acceptance Criteria |
|---------|-----------|----------|-------------------|
| Standardized Frontmatter | As a developer, I want clear agent triggers so that I know when to use each agent | Must | - [ ] YAML format with name/description<br>- [ ] 2-3 usage examples<br>- [ ] Clear context/user/assistant format |
| Personality System | As a user, I want memorable agent personalities so that interactions feel natural | Must | - [ ] Unique emoji per agent<br>- [ ] Consistent behavioral traits<br>- [ ] Personality reinforces role |
| Context Preservation | As a user, I want agents to remember previous interactions so that I don't repeat myself | Must | - [ ] Previous conversation history section<br>- [ ] Session/Agent ID tracking<br>- [ ] Context handoff capability |

### Feature Set 2: Advanced Capabilities
| Feature | User Story | Priority | Acceptance Criteria |
|---------|-----------|----------|-------------------|
| Parallel Decomposition | As a developer, I want complex tasks broken down so that they're solved efficiently | Must | - [ ] Decomposition logic in process<br>- [ ] 3-7 parallel task support<br>- [ ] Clear scope boundaries |
| Template Integration | As a user, I want consistent documentation so that outputs are professional | Should | - [ ] Reference templates via placeholders<br>- [ ] Automatic path resolution<br>- [ ] Document creation patterns |
| Task Handoffs | As a user, I want seamless agent collaboration so that complex workflows complete smoothly | Must | - [ ] Structured task blocks<br>- [ ] Agent assignment syntax<br>- [ ] Status tracking |

### Feature Set 3: New Agent Additions
| Feature | User Story | Priority | Acceptance Criteria |
|---------|-----------|----------|-------------------|
| Lead Developer Agent | As a team, I want code reviews especially for AI-generated code so that quality is maintained | Must | - [ ] Code quality assessment<br>- [ ] Anti-pattern detection<br>- [ ] Mentorship through reviews<br>- [ ] Technical debt identification |
| UX Designer Agent | As a product team, I want professional UI/UX design so that our interfaces are intuitive and accessible | Must | - [ ] WCAG 2.1 AA compliance<br>- [ ] Design system creation<br>- [ ] User flow optimization<br>- [ ] Responsive design patterns |
| Compliance Officer Agent | As a business, I want regulatory compliance so that we avoid legal risks | Should | - [ ] GDPR/CCPA compliance<br>- [ ] Industry regulations (HIPAA, SOX)<br>- [ ] AI governance frameworks<br>- [ ] Audit trail design |

### Feature Prioritization (MoSCoW)
**Must Have**
- Standardized agent structure with clear triggers
- Personality consistency system
- Context preservation mechanism
- Parallel task decomposition
- Seamless task handoffs
- Lead Developer agent (AI code review critical for 2025)
- UX Designer agent (WCAG compliance mandated)

**Should Have**
- Template integration system
- Compliance Officer agent (regulatory requirements)
- Enhanced error handling patterns
- Proactive invocation patterns ("use PROACTIVELY")

**Could Have**
- Extended devops-engineer with platform capabilities
- Enhanced data-engineer with migration patterns
- Improved testing frameworks in the-tester

**Won't Have (this phase)**
- Visual agent configuration UI
- Agent performance analytics
- Custom agent training
- Separate API designer (included in architect)
- Agent orchestrator (handled by commands)

## Detailed Feature Specifications

### Feature: Standardized Agent Structure
**Description:** Every agent follows a consistent markdown structure with YAML frontmatter, clear sections, and predictable patterns.

**User Flow:**
1. User reviews agent description with examples
2. System matches trigger patterns to select agent
3. Agent processes with consistent methodology
4. Output follows standardized format
5. Tasks handed off with clear assignments

**Business Rules:**
- All agents must have 2-3 usage examples
- Descriptions must include trigger conditions
- Process section must have numbered steps
- Output format must include commentary and tasks blocks

**Edge Cases:**
- Scenario: Ambiguous trigger → Expected: Chief agent routes to appropriate specialist
- Scenario: Missing context → Expected: Agent requests clarification or makes safe assumptions
- Scenario: Parallel task failure → Expected: Graceful degradation with partial results

**UI/UX Requirements:**
- Clear, scannable descriptions
- Consistent emoji usage
- Predictable output structure
- Readable task handoffs

### Feature: Parallel Task Decomposition
**Description:** Complex problems are automatically broken into parallel, independent subtasks for efficient processing.

**User Flow:**
1. Agent evaluates problem complexity
2. Identifies natural boundaries and domains
3. Launches 3-7 parallel sub-agents
4. Each agent works independently
5. Results consolidated into unified output

**Business Rules:**
- Maximum 7 parallel tasks per decomposition
- Each task must have clear scope boundaries
- No overlap between parallel tasks
- Results must be mergeable

**Edge Cases:**
- Scenario: Circular dependencies → Expected: Sequential processing fallback
- Scenario: Too many subtasks → Expected: Hierarchical decomposition
- Scenario: Conflicting results → Expected: Reconciliation logic applied

### Feature: Lead Developer Agent
**Description:** Senior code review specialist focusing on AI-generated code quality, mentorship, and technical debt management.

**User Flow:**
1. Code is generated by the-developer or AI assistant
2. Lead Developer agent automatically invoked for review
3. Agent analyzes code for quality issues
4. Provides structured feedback with severity levels
5. Suggests improvements with examples
6. Creates tasks for required fixes

**Business Rules:**
- Reviews triggered for all AI-generated code
- Critical issues block deployment
- Feedback must be constructive and educational
- Focus on maintainability over perfection

**Edge Cases:**
- Scenario: Perfect code → Expected: Positive reinforcement provided
- Scenario: Unfamiliar framework → Expected: Focus on universal principles
- Scenario: Time pressure → Expected: Prioritize critical issues only

### Feature: UX Designer Agent
**Description:** User experience and interface design specialist ensuring intuitive, accessible, and delightful user interactions.

**User Flow:**
1. Feature requirements provided
2. UX Designer creates user journey maps
3. Designs interface components and layouts
4. Ensures WCAG 2.1 AA compliance
5. Documents design system contributions
6. Hands off to developers with specifications

**Business Rules:**
- All designs must meet WCAG 2.1 AA standards
- Mobile-first responsive design required
- Follow existing design system when available
- Include all states (loading, error, empty, success)

**Edge Cases:**
- Scenario: No design system exists → Expected: Create foundational tokens
- Scenario: Conflicting requirements → Expected: User needs prioritized
- Scenario: Technical constraints → Expected: Progressive enhancement approach

### Feature: Compliance Officer Agent
**Description:** Regulatory compliance specialist ensuring legal requirements are met for data privacy, industry regulations, and AI governance.

**User Flow:**
1. System features described
2. Compliance Officer identifies applicable regulations
3. Assesses current compliance gaps
4. Designs compliance framework
5. Creates implementation checklist
6. Documents audit requirements

**Business Rules:**
- Never compromise on mandatory requirements
- Privacy by design principles applied
- Audit trails must be immutable
- User consent must be explicit and documented

**Edge Cases:**
- Scenario: Conflicting regulations → Expected: Most restrictive applied
- Scenario: New regulation → Expected: Proactive compliance approach
- Scenario: Gray area → Expected: Conservative interpretation chosen

## Integration Requirements

### External Systems
| System | Purpose | Data Flow | Authentication |
|--------|---------|-----------|----------------|
| Claude Code Core | Agent invocation and management | Both | Native integration |
| File System | Template and document access | Both | OS permissions |
| Git | Version control integration | Both | SSH/HTTPS |
| the-startup CLI | Hook processing and logging | In | Process communication |

### API Requirements
- Agent definition format (YAML + Markdown)
- Task handoff protocol
- Context preservation format
- Template placeholder resolution

## Analytics and Metrics

### Success Metrics
- **Adoption:** 90% of complex tasks use appropriate agents
- **Engagement:** 3x increase in agent invocations per session
- **Satisfaction:** 95% successful task completions without retry

### Tracking Requirements
| Event | Properties | Purpose |
|-------|------------|---------|
| Agent Invoked | agent_type, trigger_source, context_size | Usage patterns |
| Task Decomposed | parent_agent, child_count, complexity | Decomposition effectiveness |
| Task Completed | agent_type, duration, success | Performance tracking |
| Context Preserved | session_id, context_size, handoff_count | Context management |

## Release Strategy

### MVP Scope
- Core 14 agents with standardized structure
- Parallel decomposition for key agents (architect, business-analyst)
- Basic context preservation
- Template integration

### Phased Rollout
1. **Phase 1:** Structure standardization - All existing agents
2. **Phase 2:** New critical agents - API Designer, Accessibility
3. **Phase 3:** Advanced agents - Performance, Migration, ML

### Go-to-Market
- **Positioning:** "Professional-grade AI agents for serious developers"
- **Channels:** GitHub, Claude Code marketplace, developer communities
- **Support:** Comprehensive documentation, example workflows

## Risks and Dependencies
| Risk/Dependency | Impact | Mitigation |
|----------------|--------|------------|
| Claude Code API changes | Agent invocation may break | Version pinning, compatibility layer |
| Template path resolution | Documents may not generate | Fallback to inline templates |
| Context size limits | Large contexts may truncate | Intelligent summarization |
| Parallel task overhead | Performance degradation | Adaptive parallelism limits |

## Open Questions
- [ ] Should agents have configurable verbosity levels?
- [ ] How to handle agent versioning and backwards compatibility?
- [ ] Should there be agent composition patterns beyond task handoffs?

## Appendix

### Mockups/Wireframes/Competitors
Example agent output structure:
```markdown
<commentary>
(◕‿◕) **The Architect**: *adjusts blueprint thoughtfully*

This authentication system requires careful balance between security and user experience.
</commentary>

## System Design Complete

**SDD Created**: `docs/specs/005-auth-system/SDD.md`

### Executive Summary
JWT-based authentication with refresh tokens and role-based access control.

### Key Design Decisions
- **Architecture Pattern**: Token-based with separate auth service
- **Technology Stack**: JWT, bcrypt, Redis session store
- **Scalability Approach**: Stateless tokens with distributed cache

<tasks>
- [ ] Implement JWT token generation {agent: `the-developer`, creates: auth-service.ts}
- [ ] Create user roles schema {agent: `the-data-engineer`, creates: roles-migration.sql}
- [ ] Document API endpoints {agent: `the-technical-writer`, creates: auth-api.md}
</tasks>
```

### Competitive Analysis
- **GitHub Copilot**: No agent specialization, single model approach
- **Cursor**: Limited agent system, focuses on code generation
- **Cody**: Basic command system, no personality or decomposition

### Research Data

#### Official Anthropic Documentation (2024-2025)
- **Source**: Anthropic Claude Code Best Practices, docs.anthropic.com
- **Key Findings**: 
  - Sub-agents defined in Markdown with YAML frontmatter (name, description, tools, system prompt)
  - CLAUDE.md files auto-pulled into context for repository-specific behaviors
  - TDD workflow is "Anthropic-favorite" for verifiable changes
  - Proactive use requires "use PROACTIVELY" or "MUST BE USED" in descriptions
  - Tools field inheritance: omitting grants access to all MCP tools

#### Community Implementation Analysis
- **Source**: GitHub repositories (awesome-claude-code-agents, sub-agent-collective, lst97/claude-code-sub-agents)
- **Usage Statistics**: 
  - 25+ agent invocations in test files (internal/log/testdata)
  - 14 existing agents with quality ratings 6-9/10
  - Most common agents: the-architect, the-developer, the-product-manager
- **Common Patterns**:
  - Trigger prefix "the-" for all agents (100% consistency in codebase)
  - Session/Agent ID tracking via regex extraction
  - Hook-based logging system for agent interactions

#### Community Problem Reports (GitHub Issues, Hacker News, Medium)
- **Primary Issue**: "Agents often aren't used - Claude does work itself without invoking agent" (GitHub Issue #345)
- **Triggering Problems**: Requires manual "use subagent" or "multiple subagents" in prompts
- **GitHub Actions**: v1.0.60 added support but agents not responsive in CI environment
- **Standardization Gap**: 3+ competing community implementations with different approaches

#### Performance Metrics
- **Multi-Agent Research System** (Anthropic Engineering): 70% reduction in complex task time
- **Parallel Processing**: 3-7 agents optimal for decomposition (based on codebase analysis)
- **Context Window**: Independent agent contexts prevent pollution, enable scaling
- **Thinking Modes**: Progressive budgets ("think" < "think hard" < "ultrathink")

#### Internal Codebase Analysis
- **Agent Quality Distribution**:
  - 3 agents at 9/10: the-architect, the-business-analyst, the-project-manager
  - 4 agents at 8/10: the-chief, the-prompt-engineer, the-context-engineer, the-site-reliability-engineer
  - 6 agents at 7/10: Various specialists
  - 1 agent at 6/10: the-product-manager (duplicate content lines 72-94)
- **Missing Capabilities**: 
  - No parallel decomposition in 11/14 agents
  - Inconsistent output formatting
  - Missing template integration in 8/14 agents