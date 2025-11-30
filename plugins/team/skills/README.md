# Skills Library

Reusable expertise modules that provide consistent guidance across multiple agents.

## Directory Structure

```
skills/
├── cross-cutting/       # Universal skills for all agents
│   ├── best-practices/
│   ├── codebase-exploration/
│   ├── documentation-reading/
│   ├── error-handling/
│   ├── framework-detection/
│   └── pattern-recognition/
├── design/              # UX and accessibility skills
│   ├── accessibility-standards/
│   └── user-research-methods/
├── development/         # Software development skills
│   ├── api-design-patterns/
│   ├── data-modeling/
│   ├── documentation-creation/
│   └── testing-strategies/
├── infrastructure/      # DevOps and platform skills
│   ├── cicd-patterns/
│   └── observability-patterns/
└── quality/             # Quality assurance skills
    ├── performance-profiling/
    └── security-assessment/
```

## Skills Index

| Skill | Category | Description |
|-------|----------|-------------|
| `accessibility-standards` | design | WCAG 2.1 AA compliance patterns, screen reader compatibility, keyboard navigation |
| `api-design-patterns` | development | REST and GraphQL API design patterns, OpenAPI/Swagger specifications |
| `best-practices` | cross-cutting | Security, performance, and accessibility standards |
| `cicd-patterns` | infrastructure | Pipeline design, deployment strategies (blue-green, canary, rolling) |
| `codebase-exploration` | cross-cutting | Navigate, search, and understand project structures |
| `data-modeling` | development | Schema design, entity relationships, normalization |
| `documentation-creation` | development | ADRs, system documentation, API documentation, runbooks |
| `documentation-reading` | cross-cutting | Interpret existing docs, READMEs, specs, and configuration files |
| `error-handling` | cross-cutting | Consistent error patterns, validation approaches, recovery strategies |
| `framework-detection` | cross-cutting | Auto-detect project tech stacks (React, Vue, Express, Django, etc.) |
| `observability-patterns` | infrastructure | Monitoring strategies, distributed tracing, SLI/SLO design |
| `pattern-recognition` | cross-cutting | Identify existing codebase patterns for consistency |
| `performance-profiling` | quality | Measurement approaches, profiling tools, optimization patterns |
| `security-assessment` | quality | Vulnerability review, OWASP patterns, secure coding practices |
| `testing-strategies` | development | Test pyramid principles, coverage targets, framework-specific patterns |
| `user-research-methods` | design | Interview techniques, persona creation, journey mapping |

## Usage

Skills are referenced in agent YAML frontmatter:

```yaml
---
name: my-agent
skills: codebase-exploration, framework-detection, best-practices
---
```

When the agent is invoked, Claude Code automatically loads the specified skills into context.

## Creating New Skills

Each skill folder contains:
- `SKILL.md` - Skill definition with frontmatter (name, description)
- Optional resource files (checklists, templates, references)

See existing skills for examples.
