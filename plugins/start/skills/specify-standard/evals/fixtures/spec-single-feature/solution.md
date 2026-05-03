---
title: "Notification API — Solution Design"
status: completed
---

# Solution Design

## Components

### Notifications API
REST controller exposing POST /notifications and GET /notifications/:id. Validates payload, writes to the notifications table with state=pending, returns 202.

### Notifications Repository
PostgreSQL-backed repository with idempotency-key index. CRUD operations and state transitions.

## Key Decisions

- **ADR-1**: PostgreSQL for persistence — existing infra, transactional state transitions.
- **ADR-2**: In-process background worker (not separate service) — simpler ops, sufficient throughput at current scale.
- **ADR-3**: Idempotency keys are client-supplied UUIDs — server enforces uniqueness via DB constraint.

## Data Model

```
notifications (
  id            uuid primary key,
  topic         text not null,
  payload       jsonb not null,
  idempotency_key text not null unique,
  state         text not null, -- pending | delivered | failed
  attempts      int not null default 0,
  created_at    timestamptz not null,
  updated_at    timestamptz not null
)
```

## Delivery Worker

Polls notifications where state=pending, attempts webhook delivery with exponential backoff (1s, 2s, 4s, 8s, 16s). Marks delivered on 2xx, failed after 5 attempts.
