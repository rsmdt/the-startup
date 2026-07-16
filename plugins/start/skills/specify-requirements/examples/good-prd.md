# Example: Well-Structured PRD

This example demonstrates a properly completed PRD section for reference.

## Product Overview (Example)

### Vision
Enable small business owners to accept payments anywhere, anytime, without complex setup or expensive hardware.

### Problem Statement
Small business owners lose an average of $2,400 annually in sales because they can only accept cash or checks. 73% of consumers prefer card payments, and businesses without card acceptance miss 1 in 4 sales opportunities. The existing solutions require expensive hardware ($300+), complex contracts, and take 2-3 weeks to set up.

### Value Proposition
Instant mobile payment acceptance with just a smartphone. No hardware needed, no long-term contracts, and setup takes under 5 minutes. Businesses can start accepting payments today at a flat 2.5% rate with next-day deposits.

---

## User Personas (Example)

### Primary Persona: Sarah the Solo Entrepreneur
- **Demographics:** Age 28-45, runs a service business (cleaning, tutoring, personal training), moderate tech comfort, uses smartphone daily
- **Goals:** Accept payments immediately after service, look professional to clients, minimize time spent on admin tasks
- **Pain Points:** Loses clients who don't carry cash, awkward payment conversations, delayed payments hurt cash flow

---

## Stakeholders (Example)

| Stakeholder | Interest / concern | What success looks like for them |
|-------------|--------------------|----------------------------------|
| Support team | Payment disputes become their ticket volume | Dispute flow resolvable from the admin view without engineering |
| Compliance officer | Card-data handling obligations | Zero card data stored on our side; processor attestation on file |

---

## Feature Requirements (Example)

### Must Have Features

#### FR-1: One-Tap Payment Request
- **User Story:** As Sarah, I want to send a payment request with one tap so that I can get paid immediately after finishing a service
- **Acceptance Criteria:**
  - [ ] **AC-1.1** — Given a completed service entry, When Sarah taps "Request payment", Then the request is created in under 10 seconds
  - [ ] **AC-1.2** — Given a created request, When it is sent, Then the customer receives it via SMS or email
  - [ ] **AC-1.3** — Given the customer approves, When payment executes, Then it completes within 30 seconds and both parties see confirmation immediately

> **Anti-example:** "- [ ] User can request payment" — vague, not testable, no ID. Every criterion needs Given/When/Then, an observable outcome, and a stable AC-N.M.

---

## Non-Functional Requirements (Example)

| ID | Category | Requirement (user-observable) | Fit criterion (measurable) |
|----|----------|-------------------------------|----------------------------|
| NFR-1 | Performance | Payment confirmation feels immediate | p95 request→confirmation < 3s at 500 concurrent merchants |
| NFR-2 | Security & privacy | A merchant can never see another merchant's customers | Cross-tenant access test suite passes with zero leaks |
| NFR-3 | Accessibility | Sarah can run the whole flow one-handed outdoors | Core journey passes WCAG 2.2 AA; completable with screen reader |

---

## Evaluation & Validation Criteria (Example)

| ID | Dataset / corpus | Measures | Gate |
|----|------------------|----------|------|
| EV-1 | Golden transaction set — 200 representative payment scenarios incl. declines and timeouts | End-to-end payment correctness | 100% terminal-state correctness; no payment stuck non-terminal |

| ID | Behavior under test | Pass condition (observable) |
|----|---------------------|------------------------------|
| EV-2 | No-spam: reminders only when warranted | Zero reminders sent in the paid-on-time scenario set |
| EV-3 | Access-off conduct: processor unreachable | Merchant sees explicit failure with retry; never a silent fake success |

---

## Success Metrics (Example)

### Key Performance Indicators
- **Adoption:** 10,000 active merchants in first 6 months (merchants who process at least 1 payment/month)
- **Engagement:** Average 8 transactions per merchant per month
- **Quality:** Payment success rate > 98%, support ticket rate < 2%
- **Business Impact:** $50M monthly payment volume by month 12

### Tracking Requirements
| Event | Properties | Purpose |
|-------|------------|---------|
| payment_request_created | amount, merchant_id, request_method | Measure adoption and behavior |
| payment_completed | amount, time_to_complete, payment_method | Measure success and speed |
| payment_failed | error_code, step_failed | Identify friction points |
| merchant_churned | last_active_date, total_processed | Understand retention |

---

## What Makes This PRD Good

1. **Specific Problem Statement** - Includes data (73%, $2,400, 1 in 4 sales)
2. **Clear Personas** - Demographics, goals, and pain points defined
3. **Testable Acceptance Criteria** - Time limits, specific actions, measurable outcomes — each with a stable AC-N.M ID that plan.mdx traces
4. **NFRs with Fit Criteria** - Quality expectations stated as measurable, implementation-blind conditions
5. **Evaluation Gates as Business Criteria** - Reference datasets and behavioral pass conditions (WHAT), no harness design (HOW)
6. **Measurable KPIs** - Numbers and timeframes specified, distinct from pre-release EV gates
7. **No Technical Details** - Doesn't mention databases, APIs, or architecture
8. **User-Centric Language** - Written from user perspective
