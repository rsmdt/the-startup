# Team Plugin - The Agentic Startup

**Specialized agent library for Claude Code that provides expert capabilities across all software development domains.**

The `team` plugin provides 11 specialized agent roles with 31 activity-based specializations. Each agent brings deep expertise in a specific domain, enabling Claude Code to tackle complex tasks with specialist knowledge.

---

## Table of Contents

- [Installation](#installation)
- [Agent Roles](#agent-roles) — Chief, Analyst, Architect, Developer, Tester, Designer, DevOps, Meta Agent
- [Philosophy](#philosophy) — specialist expertise, activity-focused design, quality standards
- [Skills System](#skills-system) — 16 reusable skills shared across agents
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

### 🎯 The Chief

**Complexity assessment and activity routing specialist**

Routes project work by assessing complexity across multiple dimensions and identifying required activities. Enables parallel execution and eliminates bottlenecks through intelligent work decomposition.

### 📊 The Analyst (3 specializations)

**Product and project coordination specialist**

- **Requirements Analysis** - Clarify ambiguous requirements and document comprehensive specifications
- **Feature Prioritization** - Prioritize features, evaluate trade-offs, and establish success metrics
- **Project Coordination** - Break down complex projects, identify dependencies, and coordinate cross-functional work

### 🏗️ The Architect (7 specializations)

**System design and technical excellence specialist**

- **System Architecture** - Design scalable architectures with comprehensive planning and technology selection
- **Technology Research** - Research solutions, evaluate technologies, and provide informed recommendations
- **Quality Review** - Review architecture and code quality for technical excellence
- **System Documentation** - Create architectural documentation, design decision records, and integration guides
- **Security Review** - Identify vulnerabilities, injection risks, auth issues, and secrets exposure
- **Simplification Review** - Aggressively challenge complexity, enforce YAGNI, detect over-engineering
- **Compatibility Review** - Detect breaking changes, validate migration paths, ensure backwards compatibility

### 👨‍💻 The Developer (5 specializations)

**Implementation and optimization specialist**

- **API Development** - Design and document REST/GraphQL APIs with comprehensive specifications
- **Component Development** - Design UI components and manage state flows for scalable frontend applications
- **Domain Modeling** - Model business domains with proper entities, business rules, and persistence design
- **Performance Optimization** - Optimize application performance through systematic profiling and optimization
- **Concurrency Review** - Find race conditions, deadlocks, async anti-patterns, and resource leaks

### 🧪 The Tester (3 specializations)

**Quality assurance and testing specialist**

- **Test Execution** - Plan test strategies and implement comprehensive test suites
- **Exploratory Testing** - Discover defects through creative exploration and user journey validation
- **Performance Testing** - Identify performance bottlenecks and validate system behavior under load

### 🎨 The Designer (4 specializations)

**User experience and interface design specialist**

- **User Research** - Conduct user interviews, perform usability testing, and develop user insights
- **Interaction Architecture** - Design information architecture and user interactions for intuitive experiences
- **Design Foundation** - Create design systems and visual foundations for consistent user experiences
- **Accessibility Implementation** - Ensure WCAG compliance and make products usable by everyone

### ⚙️ The DevOps (8 specializations)

**Infrastructure and operations specialist**

- **Infrastructure as Code** - Write infrastructure as code and design cloud architectures
- **Containerization** - Containerize applications and design Kubernetes deployments
- **Deployment Automation** - Automate deployments with CI/CD pipelines and advanced deployment strategies
- **Production Monitoring** - Implement comprehensive monitoring and incident response for production systems
- **Performance Tuning** - Optimize system and database performance through profiling and tuning
- **Data Architecture** - Design data architectures with schema modeling and migration planning
- **Pipeline Engineering** - Design and implement data pipelines for high-volume processing
- **Dependency Review** - Assess CVEs, license compliance, supply chain security, and necessity

### 🤖 The Meta Agent

**Agent design and generation specialist**

Design and generate new Claude Code sub-agents, validate agent specifications, and refactor existing agents to follow evidence-based design principles.

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
skills: unfamiliar-codebase-navigation, tech-stack-detection, api-contract-design
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
| **Design** | `user-insight-synthesis` | Research synthesis, persona creation, testing validation |
| | `user-research` | Interview techniques, personas, journey mapping |
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
model: inherit
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
