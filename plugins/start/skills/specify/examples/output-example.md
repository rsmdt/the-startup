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

| Topic | Default | Impact if wrong | Status |
|-------|---------|-----------------|--------|
| Auth for socket connections | Reuse session cookie auth | Token-based auth = new AC set + security review | Proposed |
| Delivery guarantee | At-least-once, dedupe client-side | Exactly-once changes architecture materially | Proposed |

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
Documents: Requirements ✅ | Solution ✅ | Plan ✅
📄 Report: `.start/specs/003-notification-system/spec-report.md` (full ledger, finding dispositions, decisions — readable in any editor)

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
├── README.md                 # Decisions, assumption ledger, progress
├── requirements.md          # What and why (carries [ASSUMPTION] markers)
├── solution.md              # How
├── plan.md                  # Narrative plan
├── spec-report.md           # Editor-readable report (ledger, findings, documents, decisions)
└── .tracker/                # Slice tasks + scenario child tasks (task.js)
