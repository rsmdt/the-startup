# Skills Library

Reusable expertise modules that provide consistent guidance across multiple agents.

## Directory Structure

```
skills/
├── cross-cutting/       # Universal skills for all agents
│   ├── language-coding-conventions/
│   ├── unfamiliar-codebase-navigation/
│   ├── documentation-information-extraction/
│   ├── error-recovery-patterns/
│   ├── tech-stack-detection/
│   └── codebase-pattern-identification/
├── design/              # UX and accessibility skills
│   ├── accessible-interface-design/
│   └── user-insight-synthesis/
├── development/         # Software development skills
│   ├── api-contract-design/
│   ├── entity-relationship-design/
│   ├── technical-documentation-authoring/
│   └── comprehensive-test-design/
├── infrastructure/      # DevOps and platform skills
│   ├── pipeline-deployment-design/
│   └── production-observability-design/
└── quality/             # Quality assurance skills
    ├── performance-bottleneck-analysis/
    └── vulnerability-threat-assessment/
```

## Skills Index

| Skill | Category | Description |
|-------|----------|-------------|
| `accessible-interface-design` | design | WCAG 2.1 AA compliance patterns, screen reader compatibility, keyboard navigation |
| `api-contract-design` | development | REST and GraphQL API design patterns, OpenAPI/Swagger specifications |
| `language-coding-conventions` | cross-cutting | Security, performance, and accessibility standards |
| `pipeline-deployment-design` | infrastructure | Pipeline design, deployment strategies (blue-green, canary, rolling) |
| `unfamiliar-codebase-navigation` | cross-cutting | Navigate, search, and understand project structures |
| `entity-relationship-design` | development | Schema design, entity relationships, normalization |
| `technical-documentation-authoring` | development | ADRs, system documentation, API documentation, runbooks |
| `documentation-information-extraction` | cross-cutting | Interpret existing docs, READMEs, specs, and configuration files |
| `error-recovery-patterns` | cross-cutting | Consistent error patterns, validation approaches, recovery strategies |
| `tech-stack-detection` | cross-cutting | Auto-detect project tech stacks (React, Vue, Express, Django, etc.) |
| `production-observability-design` | infrastructure | Monitoring strategies, distributed tracing, SLI/SLO design |
| `codebase-pattern-identification` | cross-cutting | Identify existing codebase patterns for consistency |
| `performance-bottleneck-analysis` | quality | Measurement approaches, profiling tools, optimization patterns |
| `vulnerability-threat-assessment` | quality | Vulnerability review, OWASP patterns, secure coding practices |
| `comprehensive-test-design` | development | Test pyramid principles, coverage targets, framework-specific patterns |
| `user-insight-synthesis` | design | Interview techniques, persona creation, journey mapping |

## Usage

Skills are referenced in agent YAML frontmatter:

```yaml
---
name: my-agent
skills: unfamiliar-codebase-navigation, tech-stack-detection, language-coding-conventions
---
```

When the agent is invoked, Claude Code automatically loads the specified skills into context.

## Creating New Skills

Each skill folder contains:
- `SKILL.md` - Skill definition with frontmatter (name, description)
- Optional resource files (checklists, templates, references)

See existing skills for examples.
