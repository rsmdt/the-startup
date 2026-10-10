---
title: "Payments Redesign"
status: completed
version: "1.0"
---

# Product Requirements Document

## Problem

Synchronous Stripe charges at checkout fail 3-D Secure cards and double-charge customers when the browser retries, causing about 40 support tickets a week.

## Outcome

Every card, including 3-D Secure, completes checkout, and no customer is charged twice for one order.

## Who

| Role | Relationship | Goal or concern today | What success looks like for them |
|------|--------------|-----------------------|----------------------------------|
| Customer | uses it | 3-D Secure cards fail; duplicate charges | Pays once, any card |
| Support agent | affected | Handles duplicate-charge tickets | Duplicate-charge tickets stop |

## Features

### Must have

#### FR-1: Confirm payments asynchronously
- **User story:** As a customer, I want my payment to complete even when my bank asks for verification so that I can check out.
- **Acceptance criteria:**
  - [ ] **AC-1.1** — Given a 3-D Secure card, When the customer completes the challenge, Then the order is marked paid
  - [ ] **AC-1.2** — Given a payment pending confirmation, When Stripe reports success, Then the order is marked paid exactly once

### Should have

None.

### Won't have (this phase)

- New payment methods

## Quality gates

| ID | Concern | Requirement (user-observable) | Gate (measurable) |
|----|---------|-------------------------------|-------------------|
| NFR-1 | Reliability | No double charges | 0 duplicate charges under retried requests in test |

## Constraints & dependencies

- Replaces the synchronous charge flow from spec 004.

## Glossary

| Term | Definition |
|------|------------|
| PaymentIntent | Stripe object tracking one payment through confirmation |

## Open questions

None — all decisions settled.
