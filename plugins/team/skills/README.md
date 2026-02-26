# Skills Library

Reusable expertise modules that provide consistent guidance across multiple agents.

## Directory Structure

```
skills/
├── cross-cutting/
│   ├── project-discovery/
│   ├── pattern-detection/
│   ├── feature-prioritization/
│   └── requirements-elicitation/
├── design/
│   └── user-research/
├── development/
│   ├── api-contract-design/
│   ├── architecture-selection/
│   ├── domain-modeling/
│   ├── technical-writing/
│   └── testing/
├── infrastructure/
│   └── platform-operations/
└── quality/
    ├── code-quality-review/
    ├── performance-analysis/
    └── security-assessment/
```

## Skills Index

| Skill | Category | Description |
|-------|----------|-------------|
| `project-discovery` | cross-cutting | Unified structure mapping, stack detection, and doc verification |
| `pattern-detection` | cross-cutting | Identify and apply local codebase patterns |
| `feature-prioritization` | cross-cutting | Prioritization frameworks and decision trade-off analysis |
| `requirements-elicitation` | cross-cutting | Clarify vague requirements and define testable acceptance criteria |
| `user-research` | design | Research planning plus insight synthesis for product/design decisions |
| `api-contract-design` | development | API contracts, versioning, and auth patterns |
| `architecture-selection` | development | Architecture pattern selection with trade-off analysis |
| `domain-modeling` | development | Domain/data modeling, invariants, schema evolution |
| `technical-writing` | development | ADRs, architecture docs, API docs, runbooks |
| `testing` | development | Layered testing strategy and execution guidance |
| `platform-operations` | infrastructure | CI/CD, deployment safety, observability, SLI/SLO strategy |
| `code-quality-review` | quality | Structured code review with cross-cutting quality standards |
| `performance-analysis` | quality | Profiling, baseline measurement, optimization strategy |
| `security-assessment` | quality | Security review and threat-modeling patterns |

## Usage

Skills are referenced in agent YAML frontmatter:

```yaml
---
name: my-agent
skills: project-discovery, pattern-detection, testing
---
```

When the agent is invoked, Claude Code loads the listed skills into context.

## Creating New Skills

Each skill folder contains:
- `SKILL.md` - Skill definition with frontmatter (`name`, `description`)
- Optional support files (`reference/`, `templates/`, `examples/`, `checklists/`)
