# Example Implementation Output

## Plan Discovery

Spec: 003-notification-system
Plan: `.start/specs/003-notification-system/plan/`

Phases:
1. Core Foundation — `phase-1.md` — completed (skipping)
2. API Layer — `phase-2.md` — in_progress (resuming)
3. Integration & Validation — `phase-3.md` — pending

Starting from Phase 2 (1 phase already completed).

---

## Task Result — Success

✅ T2.1: Notification API Endpoints

Files: src/routes/notifications.ts, src/middleware/validate.ts
Summary: Implemented REST endpoints for notification CRUD with input validation
Tests: 8 passing

## Task Result — Blocked

⚠️ T2.3: WebSocket Integration

Status: Blocked
Reason: Missing WebSocket server setup from Phase 1
Options: [present via AskUserQuestion]

---

## Phase Summary

Phase 2 Complete: API Layer

Tasks: 4/4 completed
Files Changed: src/routes/notifications.ts, src/middleware/validate.ts, src/services/notifier.ts, src/dto/notification.ts
Tests: All passing (12 new)
Blockers: None
Drift: None detected

Status updates:
- `plan/phase-2.md` frontmatter: `status: completed`
- `plan/README.md`: `- [x] [Phase 2: API Layer](phase-2.md)`

Next: Phase 3 (Integration & Validation)

---

## Completion Summary

✅ Implementation Complete

Spec: 003-notification-system
Phases Completed: 3/3
Tasks Executed: 11 total
Tests: All passing (28 total)
Mode: Standard

Files Changed: 12 files (+634 -18)

Plan status: All phase files updated to `status: completed`, all README checkboxes checked.

---

## Paused Summary (when user chooses Pause)

⏸ Implementation Paused

Spec: 003-notification-system
Phases Completed: 2/3
Current Phase: 3 (Integration & Validation) — pending

Resume with: `/start:implement 003`
The plan will pick up from Phase 3 automatically.
