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

## Feature Requirements (Example)

### Must Have Features

#### Feature 1: One-Tap Payment Request
- **User Story:** As Sarah, I want to send a payment request with one tap so that I can get paid immediately after finishing a service
- **Acceptance Criteria:**
  - [ ] Can create payment request in under 10 seconds
  - [ ] Customer receives request via SMS or email
  - [ ] Payment completes within 30 seconds of customer action
  - [ ] Confirmation shown to both parties immediately

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
3. **Testable Acceptance Criteria** - Time limits, specific actions, measurable outcomes
4. **Measurable KPIs** - Numbers and timeframes specified
5. **No Technical Details** - Doesn't mention databases, APIs, or architecture
6. **User-Centric Language** - Written from user perspective
