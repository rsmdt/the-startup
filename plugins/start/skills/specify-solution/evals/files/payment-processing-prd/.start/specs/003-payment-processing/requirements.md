---
title: "Payment Processing"
status: completed
version: "1.0"
---

# Product Requirements Document

## Problem

Customers can build a cart but cannot pay online; orders are invoiced manually and 30% are never paid.

## Outcome

Customers pay for an order by card at checkout through Stripe, see a clear confirmation, and finance can issue refunds without engineering help.

## Who

| Role | Relationship | Goal or concern today | What success looks like for them |
|------|--------------|-----------------------|----------------------------------|
| Customer | uses it | Cannot pay online | Pays by card in one step |
| Finance operator | uses it | Refunds require an engineer | Issues a refund from the admin screen |

## Features

### Must have

#### FR-1: Pay for an order by card
- **User story:** As a customer, I want to pay for my order by card so that it ships without waiting for an invoice.
- **Acceptance criteria:**
  - [ ] **AC-1.1** — Given an unpaid order, When the customer pays with a valid card, Then the order is marked paid and a confirmation is shown
  - [ ] **AC-1.2** — Given an unpaid order, When the card is declined, Then the order stays unpaid and the customer sees the decline reason
  - [ ] **AC-1.3** — Given a card requiring 3-D Secure, When the customer completes the challenge, Then the payment completes as in AC-1.1

#### FR-2: Refund a paid order
- **User story:** As a finance operator, I want to refund a paid order so that customer disputes are resolved quickly.
- **Acceptance criteria:**
  - [ ] **AC-2.1** — Given a paid order, When the operator issues a full refund, Then the order is marked refunded and the customer is emailed
  - [ ] **AC-2.2** — Given an already refunded order, When a refund is attempted, Then it is rejected with an explanation

### Should have

None.

### Won't have (this phase)

- Partial refunds, saved cards, subscriptions

## Quality gates

| ID | Concern | Requirement (user-observable) | Gate (measurable) |
|----|---------|-------------------------------|-------------------|
| NFR-1 | Security | Card data never touches our servers | PCI scope SAQ-A; no PAN in logs or database |
| NFR-2 | Reliability | A customer is never charged twice for one order | 0 duplicate charges under client retries in test |
| NFR-3 | Performance | Checkout feels fast | p95 payment confirmation < 3 s excluding 3-D Secure |

## Data

| ID | Data item | Purpose | Retention | Classification |
|----|-----------|---------|-----------|----------------|
| DR-1 | Stripe payment ID per order | Reconciliation and refunds | 7 years | internal |

- **User rights:** Customers request payment history export via support; honored within 30 days.

## Constraints & dependencies

- Stripe is the mandated provider.

| Dependency | Owner | Needed by | Impact if late |
|------------|-------|-----------|----------------|
| Stripe production account | Finance | Launch | No live payments |

## Glossary

| Term | Definition |
|------|------------|
| 3-D Secure | Card issuer's extra authentication step |

## Open questions

None — all decisions settled.
