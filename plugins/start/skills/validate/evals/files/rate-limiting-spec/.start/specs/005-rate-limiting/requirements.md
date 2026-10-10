---
title: "API Rate Limiting"
status: completed
version: "1.0"
---

# Product Requirements Document

## Problem

A few API clients send bursts that slow the public API for everyone; last month two incidents were traced to a single integration.

## Outcome

Each API key gets a fair request budget; clients that exceed it are told when to retry, and other clients are unaffected.

## Who

| Role | Relationship | Goal or concern today | What success looks like for them |
|------|--------------|-----------------------|----------------------------------|
| API client developer | uses it | Unpredictable slowdowns | Clear limits and retry guidance |
| On-call engineer | affected | Paged for noisy-neighbour incidents | No incidents caused by one key |

## Features

### Must have

#### FR-1: Enforce per-key limits
- **User story:** As an on-call engineer, I want each API key limited so that one client cannot degrade the API.
- **Acceptance criteria:**
  - [ ] **AC-1.1** — Given a key under its limit, When it sends a request, Then the request is served normally
  - [ ] **AC-1.2** — Given a key over its limit, When it sends a request, Then it receives 429 with a Retry-After header

#### FR-2: Expose limit status
- **User story:** As an API client developer, I want to see my remaining budget so that I can pace requests.
- **Acceptance criteria:**
  - [ ] **AC-2.1** — Given any served request, When the response is returned, Then it includes X-RateLimit-Limit and X-RateLimit-Remaining headers
  - [ ] **AC-2.2** — Given a key with a custom plan limit, When it sends a request, Then the headers reflect the plan limit, not the default

### Should have

None.

### Won't have (this phase)

- Per-endpoint limits

## Quality gates

| ID | Concern | Requirement (user-observable) | Gate (measurable) |
|----|---------|-------------------------------|-------------------|
| NFR-1 | Performance | Limiting adds no noticeable latency | p95 overhead < 5 ms |
| NFR-2 | Accuracy | Limits are honored across instances | ≤ 1% over-admission at 10 instances |

## Constraints & dependencies

- Shared Redis cluster is available.

## Glossary

| Term | Definition |
|------|------------|
| Budget | Requests allowed per key per minute |

## Open questions

None — all decisions settled.
