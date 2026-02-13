# Skills Library

Reusable expertise modules that provide consistent guidance across multiple agents.

## Directory Structure

```
skills/
├── cross-cutting/              # Universal skills for all agents
│   ├── codebase-navigation/
│   ├── coding-conventions/
│   ├── documentation-extraction/
│   ├── feature-prioritization/
│   ├── pattern-detection/
│   ├── requirements-elicitation/
│   └── tech-stack-detection/
├── design/                     # UX and design skills
│   ├── user-insight-synthesis/
│   └── user-research/
├── development/                # Software development skills
│   ├── api-contract-design/
│   ├── architecture-selection/
│   ├── data-modeling/
│   ├── domain-driven-design/
│   ├── technical-writing/
│   └── testing/
├── infrastructure/             # DevOps and platform skills
│   ├── deployment-pipeline-design/
│   └── observability-design/
└── quality/                    # Quality assurance skills
    ├── code-quality-review/
    ├── performance-analysis/
    └── security-assessment/
```

## Skills Index

| Skill | Category | Description |
|-------|----------|-------------|
| `codebase-navigation` | cross-cutting | Navigate, search, and understand project structures |
| `coding-conventions` | cross-cutting | Security, performance, accessibility, and error handling standards |
| `documentation-extraction` | cross-cutting | Interpret existing docs, READMEs, specs, and configuration files |
| `feature-prioritization` | cross-cutting | RICE, MoSCoW, Kano, and value-effort prioritization frameworks |
| `pattern-detection` | cross-cutting | Identify existing codebase patterns for consistency |
| `requirements-elicitation` | cross-cutting | Requirement gathering, stakeholder analysis, user story patterns |
| `tech-stack-detection` | cross-cutting | Auto-detect project tech stacks (React, Vue, Express, Django, etc.) |
| `user-insight-synthesis` | design | Research synthesis, persona creation, testing validation |
| `user-research` | design | Interview techniques, persona creation, journey mapping |
| `api-contract-design` | development | REST and GraphQL API design patterns, OpenAPI/Swagger |
| `architecture-selection` | development | Monolith, microservices, serverless architecture patterns |
| `data-modeling` | development | Schema design, entity relationships, normalization |
| `domain-driven-design` | development | DDD patterns, bounded contexts, aggregates |
| `technical-writing` | development | ADRs, system documentation, API documentation, runbooks |
| `testing` | development | Test pyramid principles, coverage targets, framework patterns |
| `deployment-pipeline-design` | infrastructure | Pipeline design, deployment strategies (blue-green, canary) |
| `observability-design` | infrastructure | Monitoring strategies, distributed tracing, SLI/SLO design |
| `code-quality-review` | quality | Systematic code review patterns and feedback techniques |
| `performance-analysis` | quality | Measurement approaches, profiling tools, optimization patterns |
| `security-assessment` | quality | Vulnerability review, OWASP patterns, secure coding practices |

## Usage

Skills are referenced in agent YAML frontmatter:

```yaml
---
name: my-agent
skills: codebase-navigation, tech-stack-detection, coding-conventions
---
```

When the agent is invoked, Claude Code automatically loads the specified skills into context.

## Creating New Skills

Each skill folder contains:
- `SKILL.md` - Skill definition with frontmatter (name, description)
- Optional resource files (checklists, templates, references)

See existing skills for examples.
