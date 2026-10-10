---
title: "Saved Searches"
status: completed
version: "1.0"
---

# Product Requirements Document

## Problem

Support agents rebuild the same ticket search filters many times a day, losing several minutes per shift.

## Outcome

Agents save a filtered ticket search under a name and rerun it in one click.

## Who

| Role | Relationship | Goal or concern today | What success looks like for them |
|------|--------------|-----------------------|----------------------------------|
| Support agent | uses it | Rebuilds filters repeatedly | Reruns a saved search in one click |
| Support lead | affected | Agents triage inconsistently | Agents share the same named views |

## Features

### Must have

#### FR-1: Save a search
- **User story:** As a support agent, I want to save my current filters under a name so that I can reuse them.
- **Acceptance criteria:**
  - [ ] **AC-1.1** — Given filters are applied, When the agent saves them with a name, Then the search appears in their saved list
  - [ ] **AC-1.2** — Given a saved search with the same name exists, When the agent saves, Then they are asked to overwrite or rename

#### FR-2: Run a saved search
- **User story:** As a support agent, I want to run a saved search so that I see matching tickets immediately.
- **Acceptance criteria:**
  - [ ] **AC-2.1** — Given a saved search, When the agent selects it, Then the ticket list shows results for its filters

### Should have

None.

### Won't have (this phase)

- Sharing searches between agents

## Quality gates

| ID | Concern | Requirement (user-observable) | Gate (measurable) |
|----|---------|-------------------------------|-------------------|
| NFR-1 | Performance | Saved search runs as fast as a manual one | p95 < 800 ms |
| NFR-2 | Limits | Lists stay manageable | Max 50 saved searches per agent, enforced |

## Constraints & dependencies

- None.

## Glossary

| Term | Definition |
|------|------------|
| Saved search | A named set of ticket filters owned by one agent |

## Open questions

None — all decisions settled.
