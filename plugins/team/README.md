# Team Plugin - The Agentic Startup

**Specialized agent library for Claude Code that provides expert capabilities across all software development domains.**

The `team` plugin provides 8 specialized agent roles with 22 activity-based specializations, backed by 20 reusable skills. Each agent brings deep expertise in a specific domain, enabling Claude Code to tackle complex tasks with specialist knowledge.

---

## Table of Contents

- [Installation](#installation)
- [Agent Roles](#agent-roles) — Chief, Analyst, Architect, Developer, Tester, Designer, DevOps, Meta Agent
- [Philosophy](#philosophy) — specialist expertise, activity-focused design, quality standards
- [Skills System](#skills-system) — 20 reusable skills shared across agents
- [Agent Architecture](#agent-architecture) — slim agent design, template structure

---

## Installation

```bash
# Install the Team plugin
/plugin install team@the-startup
```

**Note:** The Team plugin is optional but recommended for complex projects requiring specialist expertise.

---

## Agent Roles

### The Chief

**Complexity assessment and activity routing specialist**

Assesses complexity and routes work when facing multi-step tasks, unclear requirements, or cross-domain work. Identifies parallel execution opportunities and decomposes work into focused activities.

### The Analyst (2 activities)

**Research and requirements specialist**

| Activity | Focus |
|----------|-------|
| `research-market` | Competitive analysis, market gaps, industry trends, product positioning |
| `research-requirements` | Requirement clarification, specification writing, stakeholder analysis |

### The Architect (4 activities)

**System design and technical excellence specialist**

| Activity | Focus |
|----------|-------|
| `design-system` | Scalable architecture design, microservices vs monolith decisions, deployment architecture |
| `review-security` | Vulnerability detection, injection prevention, secrets detection, cryptographic review |
| `review-complexity` | YAGNI enforcement, over-engineering detection, unnecessary abstraction removal |
| `review-compatibility` | Breaking change detection, migration path validation, backwards compatibility |

### The Developer (3 activities)

**Implementation and optimization specialist**

| Activity | Focus |
|----------|-------|
| `build-feature` | UI components, API endpoints, services, database logic, integrations |
| `optimize-performance` | Page loads, API latency, query performance, memory leaks, bundle sizes |
| `review-concurrency` | Race conditions, deadlocks, async anti-patterns, resource leaks |

### The Tester (2 activities)

**Quality assurance and testing specialist**

| Activity | Focus |
|----------|-------|
| `test-quality` | Test strategy, test suite implementation, coverage analysis |
| `test-performance` | Load testing, stress testing, capacity modeling, bottleneck identification |

### The Designer (4 activities)

**User experience and interface design specialist**

| Activity | Focus |
|----------|-------|
| `research-user` | User interviews, usability testing, persona creation, insight synthesis |
| `design-interaction` | Information architecture, user flows, navigation design, content organization |
| `design-visual` | Design systems, component libraries, tokens, style guides |
| `build-accessibility` | WCAG compliance, accessible forms, interactive elements, assistive technology |

### The DevOps (5 activities)

**Infrastructure and operations specialist**

| Activity | Focus |
|----------|-------|
| `build-containers` | Dockerfiles, multi-stage builds, image optimization, container security |
| `build-infrastructure` | Terraform, CloudFormation, Pulumi, cloud architecture, reusable modules |
| `build-pipelines` | GitHub Actions, GitLab CI, Jenkins, deployment strategies, rollback automation |
| `monitor-production` | Metrics, alerting, SLIs/SLOs, observability, incident diagnostics |
| `review-dependency` | CVE detection, license compliance, supply chain security, necessity assessment |

### The Meta Agent

**Agent design and generation specialist**

Designs and generates new Claude Code sub-agents, validates agent specifications, and refactors existing agents to follow evidence-based design principles.

---

## Philosophy

### Specialist Expertise

Each agent is designed with:
- **Deep domain knowledge** - Focused expertise in a specific activity
- **Clear boundaries** - Well-defined scope to prevent overlap
- **Proven patterns** - Built on evidence-based practices
- **Pragmatic approach** - Balance technical excellence with delivery

### Activity-Focused Design

Agents are organized by **what they do**, not **who they are**:
- Enables precise routing to the right specialist
- Prevents scope creep and role confusion
- Allows parallel execution of independent activities
- Maintains clear accountability for outcomes

### Quality Without Compromise

Agents enforce quality standards:
- Architectural reviews catch design issues early
- Test strategies ensure comprehensive coverage
- Performance profiling prevents optimization guesswork
- Documentation keeps knowledge accessible

---

## Skills System

The Team plugin includes a **skills library** that provides reusable expertise shared across multiple agents. Skills eliminate content duplication while ensuring consistent guidance.

### How Skills Work

Skills are referenced in agent YAML frontmatter:

```yaml
---
name: api-development
skills: codebase-navigation, tech-stack-detection, api-contract-design
---
```

When an agent is invoked, Claude Code automatically loads the referenced skills into context, providing the agent with specialized knowledge without duplicating content across agent files.

### Available Skills (20)

| Category | Skill | Description |
|----------|-------|-------------|
| **Cross-Cutting** | `codebase-navigation` | Navigate, search, and understand project structures |
| | `coding-conventions` | Security, performance, accessibility, and error handling standards |
| | `documentation-extraction` | Interpret docs, READMEs, specs, and configs |
| | `feature-prioritization` | RICE, MoSCoW, Kano, and value-effort prioritization |
| | `pattern-detection` | Identify existing codebase patterns for consistency |
| | `requirements-elicitation` | Requirement gathering and stakeholder analysis |
| | `tech-stack-detection` | Auto-detect project tech stacks and configurations |
| **Design** | `user-insight-synthesis` | Interview techniques, persona creation, journey mapping, usability testing |
| | `user-research` | User interviews, persona creation, journey mapping, research synthesis |
| **Development** | `api-contract-design` | REST/GraphQL design, OpenAPI, versioning |
| | `architecture-selection` | Monolith, microservices, serverless patterns |
| | `data-modeling` | Schema design, entity relationships, normalization |
| | `domain-driven-design` | DDD patterns, bounded contexts, aggregates |
| | `technical-writing` | ADRs, system docs, API docs, runbooks |
| | `testing` | Test pyramid, coverage targets, framework patterns |
| **Infrastructure** | `deployment-pipeline-design` | Pipeline design, deployment strategies |
| | `observability-design` | Monitoring, tracing, SLI/SLO design |
| **Quality** | `code-quality-review` | Systematic code review patterns and feedback |
| | `performance-analysis` | Measurement, profiling tools, optimization |
| | `security-assessment` | Vulnerability review, OWASP, threat modeling |

### Benefits

- **Single Source of Truth**: Update a skill once, all referencing agents benefit
- **Reduced Duplication**: Common patterns live in skills, not repeated across agents
- **Consistent Guidance**: All agents provide uniform advice for shared concerns
- **Modular Expertise**: Agents compose capabilities from relevant skills

---

## Agent Architecture

### Slim Agent Design

Agents follow a **slim template** that separates concerns:

- **Agents** define WHO (role), WHAT (focus areas), and WHEN (deliverables)
- **Skills** provide HOW (procedural knowledge, patterns, checklists)

This design follows Claude Code's progressive disclosure model where skills load on-demand while agent content always loads.

### Agent Template Structure

```markdown
---
name: agent-name
description: Clear purpose with usage examples
skills: skill1, skill2, skill3
---

{1-2 sentence role introduction}

## Focus Areas
{4-6 bullet points of what this agent specializes in}

## Approach
{3-5 high-level methodology steps}
{Reference to skills for detailed patterns}

## Deliverables
{4-6 concrete outputs}

## Quality Standards
{Non-negotiable quality criteria}

{Closing philosophy statement}
```

### Key Principles

1. **No Duplication**: Agent content should NOT repeat what skills provide
2. **Skill References**: Agents explicitly reference which skills to leverage
3. **Focused Scope**: Each agent excels at one activity, not multiple domains
4. **Concrete Deliverables**: Clear outputs the agent produces
