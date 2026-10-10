# Example: Well-Structured PRD

This example demonstrates properly completed PRD sections for reference.

## Problem (Example)

Small business owners lose an average of $2,400 annually in sales because they can only accept cash or checks. 73% of consumers prefer card payments, and businesses without card acceptance miss 1 in 4 sales opportunities. The existing solutions require expensive hardware ($300+), complex contracts, and take 2-3 weeks to set up.

## Outcome (Example)

A service business owner accepts a card payment from a smartphone within five minutes of signing up: no hardware, no long-term contract, a flat 2.5% rate, and next-day deposits. Done means a first-time merchant can request, receive, and see a confirmed payment without contacting support.

---

## Who (Example)

| Role | Relationship | Goal or concern today | What success looks like for them |
|------|--------------|-----------------------|----------------------------------|
| Solo service entrepreneur (cleaning, tutoring, training) | uses it | Gets paid immediately after a service; loses clients who don't carry cash; delayed payments hurt cash flow | Payment requested and confirmed before leaving the client's premises |
| Support team | affected | Payment disputes become their ticket volume | Dispute flow resolvable from the admin view without engineering |
| Compliance officer | affected | Card-data handling obligations | Zero card data stored on our side; processor attestation on file |

---

## Features (Example)

### Must have

#### FR-1: One-tap payment request
- **User story:** As a solo service entrepreneur, I want to send a payment request with one tap so that I can get paid immediately after finishing a service
- **Acceptance criteria:**
  - [ ] **AC-1.1** — Given a completed service entry, When the merchant taps "Request payment", Then the request is created in under 10 seconds
  - [ ] **AC-1.2** — Given a created request, When it is sent, Then the customer receives it via SMS or email
  - [ ] **AC-1.3** — Given the customer approves, When payment executes, Then it completes within 30 seconds and both parties see confirmation immediately

> **Anti-example:** "- [ ] User can request payment" — vague, not testable, no ID. Every criterion needs Given/When/Then, an observable outcome, and a stable AC-N.M.

---

## Quality gates (Example)

| ID | Concern | Requirement (user-observable) | Gate (measurable) |
|----|---------|-------------------------------|-------------------|
| NFR-1 | Performance | Payment confirmation feels immediate | p95 request→confirmation < 3 s at 500 concurrent merchants |
| NFR-2 | Security & privacy | A merchant can never see another merchant's customers | Cross-tenant access test suite passes with zero leaks |
| NFR-3 | Accessibility | The whole flow works one-handed outdoors | Core journey passes WCAG 2.2 AA; completable with a screen reader |
| NFR-4 | Correctness under test | Every payment reaches a terminal state | Golden set of 200 payment scenarios incl. declines and timeouts: 100% terminal-state correctness |
| NFR-5 | Behavior under test | Reminders only when warranted | Zero reminders sent in the paid-on-time scenario set |

---

## What Makes This PRD Good

1. **Specific problem** - Includes data (73%, $2,400, 1 in 4 sales)
2. **Outcome states "done"** - A concrete first-merchant journey, not a slogan
3. **Who covers users and the affected** - Each row has a goal and a success condition; support and compliance are present
4. **Testable acceptance criteria** - Time limits, specific actions, measurable outcomes, each with a stable AC-N.M the plan traces
5. **Quality gates are implementation-blind** - Thresholds, reference sets, and behavioral pass conditions (WHAT), no harness design (HOW)
6. **No technical details** - Doesn't mention databases, APIs, or architecture
