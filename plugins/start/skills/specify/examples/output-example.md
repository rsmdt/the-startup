# Example Specification Output

## Clarification Ledger (presented before any document work)

**What I know:** real-time notifications, WebSocket transport requested, existing Express backend.

**What I need from you** (research can't answer these):

| # | Question | Stake | Proposed default |
|---|----------|-------|------------------|
| 1 | Must notifications persist for offline users? | Decides whether a notification store + read-state model is in scope | Yes, 30-day retention |
| 2 | Web only, or mobile push too? | Push adds vendor integration + token lifecycle to requirements | Web only |

**What I'll research** (routed to perspectives, confirmed with you before use): current auth/session mechanism, existing event bus, WebSocket library fit with the Express stack.

**Assumptions I'd otherwise make:**

| ID | Assumption | Impact if wrong | Observable | Status |
|----|------------|-----------------|:----------:|--------|
| A1 | Reuse session cookie auth for socket connections | Token-based auth = new AC set + security review | no | Proposed |
| A2 | At-least-once delivery, dedupe client-side | Exactly-once changes architecture materially | yes | Proposed |

*(User answers #1 "yes, but 7 days", accepts #2 default → Corrected + Standing. Blocking unknowns resolved; document phases begin.)*

---

## Research Finding Disposition

| Finding (perspective) | Evidence | Disposition |
|-----------------------|----------|-------------|
| Session auth is JWT-in-cookie, socket handshake can reuse it (Technical) | src/middleware/auth.ts:14 | Confirmed |
| Existing Redis pub/sub used for jobs, reusable as event bus (Integration) | src/queue/bus.ts | Corrected — user: "Redis is being retired, use Postgres LISTEN/NOTIFY" |

---

## Specification Complete

Spec: 003-notification-system
Gates: Requirements ✅ | Solution ✅ | Plan ✅
📄 Spec: `.start/specs/003-notification-system/spec.md` (all three parts and the ledger — readable in any editor)

Readiness: MEDIUM
Standing assumptions (2): Web only (no mobile push) · At-least-once delivery
→ Readiness capped: "at-least-once delivery" is architecture-shaping. Confirm it to reach HIGH.

Next Steps:
1. /start:validate 003 — Validate specification quality
2. /start:implement 003 — Begin implementation

---

## Decision Logging

| Date | Decision | Rationale |
|------|----------|-----------|
| 2026-07-10 | Notification retention 7 days (not 30) | User corrected default during clarification |
| 2026-07-10 | Postgres LISTEN/NOTIFY over Redis | Redis being retired — user corrected research finding |

---

## Documentation Structure

.start/specs/003-notification-system/
├── spec.md                  # Requirements, Solution, Plan parts + Ledger (gates in frontmatter)
├── review.md                # Solution review report
└── .tracker/                # Phase tasks, slice tasks + scenario child tasks (task.js)
