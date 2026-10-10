---
title: "Notification API"
status: completed
version: "1.0"
---

# Product Requirements Document

## Problem

Product teams each send user notifications with their own ad-hoc email code. Messages are duplicated or lost, and support cannot tell a customer whether a notification was ever sent.

## Outcome

Any internal service sends a user notification with one API call and can later look up whether it was delivered. Support answers "did the user get it?" in seconds.

## Who

| Role | Relationship | Goal or concern today | What success looks like for them |
|------|--------------|-----------------------|----------------------------------|
| Service developer | uses it | Re-implements email sending in every service | One call sends a notification |
| Support agent | affected | Cannot see whether a notification went out | Looks up delivery status by notification ID |

## Features

### Must have

#### FR-1: Send a notification
- **User story:** As a service developer, I want to submit a notification for a user so that it is delivered without my service handling email.
- **Acceptance criteria:**
  - [ ] **AC-1.1** — Given a valid recipient and message, When a service submits a notification, Then it receives a notification ID and the notification is queued
  - [ ] **AC-1.2** — Given a request with no recipient, When a service submits it, Then it is rejected with a validation error naming the missing field

#### FR-2: Look up delivery status
- **User story:** As a support agent, I want to look up a notification's delivery status so that I can answer customer questions.
- **Acceptance criteria:**
  - [ ] **AC-2.1** — Given an existing notification ID, When status is requested, Then the current status (queued, sent, failed) and its timestamp are returned
  - [ ] **AC-2.2** — Given an unknown notification ID, When status is requested, Then a not-found response is returned

### Should have

None.

### Won't have (this phase)

- SMS and push channels
- User notification preferences

## Quality gates

| ID | Concern | Requirement (user-observable) | Gate (measurable) |
|----|---------|-------------------------------|-------------------|
| NFR-1 | Performance | Submitting a notification does not slow the caller | p95 submit latency < 200 ms at 100 req/s |
| NFR-2 | Reliability | Accepted notifications are not lost | 0 lost notifications across a worker restart during delivery |

## Constraints & dependencies

- Must use the company's existing SMTP relay.

| Dependency | Owner | Needed by | Impact if late |
|------------|-------|-----------|----------------|
| SMTP relay credentials | Platform team | Phase 1 | No delivery possible |

## Glossary

| Term | Definition |
|------|------------|
| Notification | One message addressed to one recipient |

## Open questions

None — all decisions settled.
