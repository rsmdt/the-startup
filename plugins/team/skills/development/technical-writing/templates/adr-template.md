# Template: Architecture Decision Record (ADR)

## Purpose

Use this template to document significant architectural and design decisions. ADRs capture the context, options considered, and rationale behind decisions that are:

- Difficult or expensive to reverse
- Foundational to how the system works
- Likely to be questioned by future developers
- Establishing patterns for the codebase

## Template

```markdown
# ADR-[NUMBER]: [SHORT TITLE]

**Status:** [Proposed | Accepted | Deprecated | Superseded by ADR-XXX]
**Date:** [YYYY-MM-DD when decision was made]
**Decision Makers:** [Names or roles involved in the decision]

## Context

[Describe the situation that requires a decision. Include:]

- What problem are you solving?
- What constraints exist (technical, business, timeline)?
- What assumptions are you making?
- What is the current state of the system?

[Be specific enough that someone unfamiliar with the project can understand
the situation. Avoid jargon without explanation.]

## Decision Drivers

[List the key factors that will influence this decision:]

- [Driver 1: e.g., "Must support 10,000 concurrent users"]
- [Driver 2: e.g., "Team has limited experience with NoSQL databases"]
- [Driver 3: e.g., "Budget constraints limit managed service options"]

## Options Considered

### Option 1: [Name]

[Brief description of this option]

**Pros:**
- [Advantage 1]
- [Advantage 2]

**Cons:**
- [Disadvantage 1]
- [Disadvantage 2]

### Option 2: [Name]

[Brief description of this option]

**Pros:**
- [Advantage 1]
- [Advantage 2]

**Cons:**
- [Disadvantage 1]
- [Disadvantage 2]

### Option 3: [Name]

[Brief description of this option]

**Pros:**
- [Advantage 1]
- [Advantage 2]

**Cons:**
- [Disadvantage 1]
- [Disadvantage 2]

## Decision

[State the decision clearly and unambiguously.]

We will use **[Option X]** because [primary reason].

## Rationale

[Explain why this option was chosen over the alternatives:]

- How does it address the decision drivers?
- What trade-offs are being made?
- Why were other options rejected?

[This section should help future readers understand the thinking process,
not just the outcome.]

## Consequences

### Positive

- [Benefit 1 this decision enables]
- [Benefit 2 this decision enables]

### Negative

- [Drawback 1 we are accepting]
- [Drawback 2 we are accepting]

### Neutral

- [Side effect that is neither clearly positive nor negative]

## Implementation Notes

[Optional: Include any specific guidance for implementing this decision:]

- Migration path from current state
- Key technical details
- Dependencies or prerequisites
- Estimated effort

## Related Decisions

- [ADR-XXX: Related decision about Y]
- [ADR-YYY: This decision supersedes/is superseded by]

## References

- [Link to relevant documentation, research, or discussions]
- [Link to spike or proof-of-concept if applicable]
```

## Usage Instructions

1. Copy the template above into a new file named `ADR-[NUMBER]-[slug].md`
2. Number ADRs sequentially (ADR-001, ADR-002, etc.)
3. Fill in the status as "Proposed" during discussion
4. Update status to "Accepted" when decision is finalized
5. Never edit accepted ADRs; create new ones that supersede them
6. Store ADRs in a `docs/decisions/` or `docs/adr/` directory

## Numbering Convention

Use zero-padded sequential numbers:
- `ADR-001-use-postgresql-for-persistence.md`
- `ADR-002-adopt-event-sourcing-pattern.md`
- `ADR-003-migrate-to-kubernetes.md`

## Status Transitions

```
Proposed --> Accepted --> Deprecated
                     --> Superseded by ADR-XXX
```

- **Proposed**: Under discussion, not yet binding
- **Accepted**: Decision is made and should be followed
- **Deprecated**: Being phased out, do not use for new work
- **Superseded**: Replaced by a newer ADR (always link to it)

## Examples

### Good ADR Title
- "Use PostgreSQL for primary data storage"
- "Adopt event-driven architecture for order processing"
- "Implement feature flags using LaunchDarkly"

### Poor ADR Title
- "Database decision" (too vague)
- "We should use PostgreSQL because it's better" (includes rationale in title)
- "ADR about the thing we discussed" (not descriptive)

### Good Context Section

```markdown
## Context

Our application currently stores all data in a single MySQL 5.7 database
hosted on AWS RDS. We are experiencing:

- Query latency exceeding 500ms for reporting queries
- Lock contention during high-write periods (daily imports)
- Storage costs increasing 20% month-over-month

The team has been asked to reduce p95 latency to under 100ms while
supporting 3x current data volume within 6 months. Our team has
production experience with PostgreSQL and limited experience with
NoSQL databases.
```

### Poor Context Section

```markdown
## Context

We need a better database because the current one is slow.
```

## Tips for Effective ADRs

1. **Write for future readers**: Assume the reader has no context about your project
2. **Be honest about trade-offs**: Every decision has downsides; document them
3. **Include rejected options**: Understanding why alternatives were rejected is valuable
4. **Keep it concise**: ADRs should be readable in 5-10 minutes
5. **Link to evidence**: Reference benchmarks, spikes, or discussions that informed the decision
6. **Date your decisions**: Context changes; knowing when a decision was made matters
