---
name: technical-writing
description: Create architectural decision records (ADRs), system documentation, API documentation, and operational runbooks. Use when capturing design decisions, documenting system architecture, creating API references, or writing operational procedures.
---

# Documentation Creation

A development skill for creating and maintaining technical documentation that preserves knowledge, enables informed decision-making, and supports system operations. This skill provides templates and patterns for common documentation needs.

## When to Use

- Recording architectural or design decisions with context and rationale
- Documenting system architecture for new team members or stakeholders
- Creating API documentation for internal or external consumers
- Writing runbooks for operational procedures and incident response
- Capturing tribal knowledge before it's lost to team changes

## Core Documentation Types

### Architecture Decision Records (ADRs)

ADRs capture the context, options considered, and rationale behind significant architectural decisions. They serve as a historical record that helps future developers understand why the system is built a certain way.

**When to create an ADR:**

- Choosing between different technologies, frameworks, or approaches
- Making decisions that are difficult or expensive to reverse
- Establishing patterns that will be followed across the codebase
- Deprecating existing approaches in favor of new ones
- Any decision that a future developer might question

See: `templates/adr-template.md`

### System Documentation

System documentation provides a comprehensive view of how a system works, its components, and their relationships. It helps new team members onboard and serves as a reference for operations.

**Key elements:**

- System overview and purpose
- Architecture diagrams showing component relationships
- Data flows and integration points
- Deployment architecture
- Operational requirements

See: `templates/system-doc-template.md`

### API Documentation

API documentation describes how to interact with a service, including endpoints, request/response formats, authentication, and error handling.

**Key elements:**

- Authentication and authorization
- Endpoint reference with examples
- Request and response schemas
- Error codes and handling
- Rate limits and quotas
- Versioning strategy

### Runbooks

Runbooks provide step-by-step procedures for operational tasks, from routine maintenance to incident response.

**Key elements:**

- Pre-requisites and access requirements
- Step-by-step procedures with expected outcomes
- Troubleshooting common issues
- Escalation paths
- Recovery procedures

## Documentation Patterns

### Pattern 1: Decision Context First

Always document the context and constraints that led to a decision before stating the decision itself. Future readers need to understand the "why" before the "what."

```markdown
## Context

We need to store user session data that must be:
- Available across multiple application instances
- Retrieved in under 10ms
- Retained for 24 hours after last activity

Our current database is PostgreSQL, which would require additional
infrastructure for session management.

## Decision

We will use Redis for session storage.
```

### Pattern 2: Living Documentation

Documentation should be updated as part of the development process, not as an afterthought. Integrate documentation updates into your definition of done.

- Update ADRs when decisions change (mark old ones as superseded)
- Revise system docs when architecture evolves
- Keep API docs in sync with implementation (prefer generated docs where possible)
- Review runbooks after each incident for accuracy

### Pattern 3: Audience-Appropriate Detail

Tailor documentation depth to its intended audience:

| Audience | Focus | Detail Level |
|----------|-------|--------------|
| New developers | Onboarding, getting started | High-level concepts, step-by-step guides |
| Experienced team | Reference, troubleshooting | Technical details, edge cases |
| Operations | Deployment, monitoring | Procedures, commands, expected outputs |
| Business stakeholders | Capabilities, limitations | Non-technical summaries, diagrams |

### Pattern 4: Diagrams Over Prose

Use diagrams to communicate complex relationships. A well-designed diagram can replace pages of text and is easier to maintain.

**Recommended diagram types:**

- **System context**: Shows system boundaries and external interactions
- **Container**: Shows major components and their relationships
- **Sequence**: Shows how components interact for specific flows
- **Data flow**: Shows how data moves through the system

### Pattern 5: Executable Documentation

Where possible, make documentation executable or verifiable:

- API examples that can be run against a test environment
- Code snippets that are extracted from actual tested code
- Configuration examples that are validated in CI
- Runbook steps that have been recently executed

## ADR Lifecycle

ADRs follow a specific lifecycle:

1. **Proposed**: Decision is being discussed, not yet accepted
2. **Accepted**: Decision has been made and should be followed
3. **Deprecated**: Decision is being phased out, new work should not follow it
4. **Superseded**: Decision has been replaced by a newer ADR (link to new one)

When superseding an ADR:
- Add "Superseded by ADR-XXX" to the old record
- Add "Supersedes ADR-YYY" to the new record
- Explain what changed and why in the new ADR's context

## Best Practices

- Write documentation close to the code it describes (prefer docs-as-code)
- Use templates consistently to make documentation predictable
- Include diagrams for architecture; text for procedures
- Date all documents and note last review date
- Keep ADRs immutable once accepted (create new ones to supersede)
- Store documentation in version control alongside code
- Review documentation accuracy during code reviews
- Delete or archive documentation for removed features

## Anti-Patterns to Avoid

- **Documentation Drift**: Docs that no longer match reality are worse than no docs
- **Over-Documentation**: Documenting obvious code reduces signal-to-noise
- **Wiki Sprawl**: Documentation scattered across multiple systems becomes unfindable
- **Future Fiction**: Documenting features that don't exist yet as if they do
- **Write-Only Docs**: Creating docs that no one reads or maintains

## References

- `templates/adr-template.md` - Architecture Decision Record template
- `templates/system-doc-template.md` - System documentation template
