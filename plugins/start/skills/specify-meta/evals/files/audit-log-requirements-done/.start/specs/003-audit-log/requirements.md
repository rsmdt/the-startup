---
title: "Audit Log"
status: completed
version: "1.0"
---

# Product Requirements Document

## Problem

Admins cannot see who changed workspace settings, so security reviews stall and customers fail SOC 2 evidence requests.

## Outcome

Admins see a searchable history of every settings change with who, what, and when.

## Who

| Role | Relationship | Goal or concern today | What success looks like for them |
|------|--------------|-----------------------|----------------------------------|
| Workspace admin | uses it | Cannot trace settings changes | Finds any change in seconds |
| Compliance auditor | affected | No evidence for change control | Exports the log for a date range |

## Features

### Must have

#### FR-1: Record settings changes
- **User story:** As a workspace admin, I want every settings change recorded so that I can trace it later.
- **Acceptance criteria:**
  - [ ] **AC-1.1** — Given a member changes a setting, When the change is saved, Then an entry with actor, setting, old and new value, and time is recorded
  - [ ] **AC-1.2** — Given a change fails to save, When the error occurs, Then no entry is recorded

#### FR-2: View the audit log
- **User story:** As a workspace admin, I want to filter the log by actor and date so that I find a change quickly.
- **Acceptance criteria:**
  - [ ] **AC-2.1** — Given entries exist, When the admin filters by actor and date range, Then only matching entries are shown newest first

### Should have

None.

### Won't have (this phase)

- Logging of content changes (documents, messages)

## Quality gates

| ID | Concern | Requirement (user-observable) | Gate (measurable) |
|----|---------|-------------------------------|-------------------|
| NFR-1 | Integrity | Entries cannot be edited or deleted by users | No update/delete path for entries |
| NFR-2 | Performance | Filtering feels instant | p95 < 500 ms on 1M entries |

## Data

| ID | Data item | Purpose |
|----|-----------|---------|
| DR-1 | Audit entry | Change traceability |

## Constraints & dependencies

- None.

## Glossary

| Term | Definition |
|------|------------|
| Entry | One recorded settings change |

## Open questions

None — all decisions settled.
