# Example Specification Output

## Specification Complete

Specification Complete

Spec: 003-notification-system
Documents: Requirements ✅ | Solution ✅ | Factory ✅

Readiness: HIGH
Confidence: 92%

Next Steps:
1. /start:validate 003 - Validate specification quality
2. /start:implement 003 - Begin implementation

---

## Decision Logging

## Decisions Log

| Date | Decision | Rationale |
|------|----------|-----------|
| 2026-02-10 | Requirements skipped | Requirements in JIRA-1234 |
| 2026-02-10 | Started from Factory | Requirements and solution already documented |

---

## Documentation Structure

.start/specs/003-notification-system/
├── README.md                 # Decisions and progress
├── requirements.md           # What and why
├── solution.md               # How
├── manifest.md               # Decomposition manifest
├── units/                    # Factory-sized specs (vertical capability slices)
│   ├── nt1.md                # Send a notification (owns the notification model)
│   └── ws1.md                # Stream notifications over WebSocket (after: nt1)
└── scenarios/                # Holdout evaluation scenarios
    ├── nt1/
    │   ├── delivery-confirmation.md
    │   └── persistence.md
    └── ws1/
        ├── connection-lifecycle.md
        └── reconnect-handling.md
