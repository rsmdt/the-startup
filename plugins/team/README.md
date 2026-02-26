# Team Plugin - The Agentic Startup

**Specialized agent library for Claude Code with consolidated, high-signal activities for software delivery.**

The `team` plugin provides **8 specialized roles**, **15 activity-based agents**, and **14 reusable skills**. This consolidation keeps specialist depth while reducing overlap and context bloat.

---

## Installation

```bash
/plugin install team@the-startup
```

---

## Agent Roles

### The Chief

**Complexity assessment and activity routing specialist**

- `the-chief`

### The Analyst (1 activity)

**Product research and requirements specialist**

| Activity | Focus |
|----------|-------|
| `research-product` | Market analysis, requirement clarification, prioritization, stakeholder alignment |

### The Architect (4 activities)

**System design and technical governance specialist**

| Activity | Focus |
|----------|-------|
| `design-system` | Scalable architecture, service boundaries, data + deployment strategy |
| `review-security` | Application security + dependency/supply-chain review |
| `review-robustness` | Complexity and concurrency risk review |
| `review-compatibility` | Breaking-change detection, migration safety, backwards compatibility |

### The Developer (2 activities)

**Implementation and optimization specialist**

| Activity | Focus |
|----------|-------|
| `build-feature` | Feature implementation across UI, API, services, and data layers |
| `optimize-performance` | Bottleneck diagnosis and targeted optimization |

### The Tester (1 activity)

**Functional and performance testing specialist**

| Activity | Focus |
|----------|-------|
| `test-strategy` | Risk-based test planning, quality coverage, load/stress validation |

### The Designer (3 activities)

**User research and UX design specialist**

| Activity | Focus |
|----------|-------|
| `research-user` | User interviews, persona/journey insights, behavioral evidence |
| `design-interaction` | Information architecture, user flows, interaction patterns |
| `design-visual` | Design systems plus accessibility-by-default standards |

### The DevOps (2 activities)

**Platform and operations specialist**

| Activity | Focus |
|----------|-------|
| `build-platform` | Containers + IaC + CI/CD as one delivery platform |
| `monitor-production` | Observability, SLI/SLOs, alerting, incident diagnostics |

### The Meta Agent

**Agent design and generation specialist**

- `the-meta-agent`

---

## Skills System

Skills are referenced in agent YAML frontmatter:

```yaml
---
name: agent-name
skills: project-discovery, pattern-detection, api-contract-design
---
```

### Available Skills (14)

| Category | Skill | Description |
|----------|-------|-------------|
| **Cross-Cutting** | `project-discovery` | Unified navigation, stack detection, and doc validation |
| | `pattern-detection` | Identify existing codebase patterns for consistency |
| | `feature-prioritization` | RICE/MoSCoW/Kano/value-effort prioritization |
| | `requirements-elicitation` | Requirement gathering, conflict resolution, acceptance criteria |
| **Design** | `user-research` | Research planning, evidence synthesis, personas, journey mapping |
| **Development** | `api-contract-design` | REST/GraphQL contract design and versioning |
| | `architecture-selection` | Pattern selection for monolith/microservices/serverless/event-driven |
| | `domain-modeling` | Domain + data modeling, invariants, schema evolution |
| | `technical-writing` | ADRs, system docs, runbooks, API docs |
| | `testing` | Layered testing strategy and execution patterns |
| **Infrastructure** | `platform-operations` | Pipeline design + observability + release reliability controls |
| **Quality** | `code-quality-review` | Holistic review including security/perf/a11y/error handling |
| | `performance-analysis` | Profiling and bottleneck identification |
| | `security-assessment` | Threat modeling and vulnerability review |

---

## Architecture Principles

1. Agents define **who/what/when**.
2. Skills provide reusable **how**.
3. Consolidation favors clearer boundaries and lower context overhead.
4. Deliverables stay concrete and execution-oriented.
