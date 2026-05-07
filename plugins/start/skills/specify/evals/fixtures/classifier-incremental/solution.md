---
title: "Notification API — Solution"
status: completed
---

# Solution Design

## Components

### Notifications API
REST controller, validation, idempotency check.

### Notifications Repository
PostgreSQL persistence, state transitions.

## Key Decisions
- ADR-1: PostgreSQL, in-process worker.
