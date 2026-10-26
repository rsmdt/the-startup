---
title: "Notification API"
status: completed
version: "1.0"
---

# Solution Design Document

## Approach

An HTTP API inside the existing Express monolith writes notifications to a Postgres table; a background worker polls queued rows and delivers them through the SMTP relay. The table is the queue, so an accepted notification survives restarts.

- **Deployment:** one new worker process; no other change

## Read first

| Source | Why it matters |
|--------|----------------|
| `src/app.js` | Where routers are mounted |
| `src/db/migrations/` | Migration conventions |

## Boundaries

- **Must preserve:** existing `/api/v1` routes
- **Can modify:** `src/app.js` router mounting
- **Must not touch:** the SMTP relay configuration

## Components

| Component | Responsibility | Owns |
|-----------|----------------|------|
| NotificationController | Validates and accepts requests; serves status lookups | FR-1 intake, FR-2 |
| NotificationStore | Persists notifications and status transitions | data |
| DeliveryWorker | Claims queued rows and sends them via SMTP | FR-1 delivery |

```
src/
├── notifications/
│   ├── controller.js   # NEW
│   ├── store.js        # NEW
│   └── worker.js       # NEW
└── app.js              # MODIFY: mount router
```

## Interfaces & data

| Name | Direction | Type / format | Auth | Contract |
|------|-----------|---------------|------|----------|
| `POST /api/v1/notifications` | inbound | HTTPS / JSON | API key | `{recipient, subject, body}` → `202 {id}` |
| `GET /api/v1/notifications/{id}` | inbound | HTTPS / JSON | API key | → `200 {id, status, updatedAt}` or `404` |
| SMTP relay | outbound | SMTP | credentials | existing relay |

```yaml
Table: notifications (NEW)
  id: uuid primary key
  recipient: text not null
  subject: text not null
  body: text not null
  status: text not null (queued | sent | failed)
  updated_at: timestamptz not null
```

## Runtime

A caller POSTs a notification; the controller validates it, the store inserts a `queued` row, and the caller gets `202` with the ID. The worker claims queued rows with `SELECT … FOR UPDATE SKIP LOCKED`, sends each, and marks it `sent` or `failed`.

## Decisions

> [!IMPORTANT]
> **Decision (ADR-1): Delivery queue → Postgres table polled by a worker**
>
> | Option | Trade-off / why (not) |
> |--------|-----------------------|
> | **Postgres table** | Already operated; durable; transactional with the insert |
> | Redis queue | New infrastructure to run for low volume |
>
> **Rationale:** Durability without new infrastructure · **Accepts:** polling latency of up to 5 s
>
> - [x] User confirmed

## Behavior

**Main flow — AC-1.1:**
- [ ] WHEN a valid notification is submitted, THE SYSTEM SHALL store it with status `queued` and respond `202` with its ID

**Error handling — AC-1.2:**
- [ ] IF the request has no recipient, THEN THE SYSTEM SHALL respond `400` naming the `recipient` field

**Main flow — AC-2.1:**
- [ ] WHEN status is requested for an existing ID, THE SYSTEM SHALL return its status and `updatedAt`

**Error handling — AC-2.2:**
- [ ] IF status is requested for an unknown ID, THEN THE SYSTEM SHALL respond `404`

## Risks & gotchas

- Worker crash mid-send can deliver an email twice; accepted (at-least-once).

## Open questions

None — all decisions settled.
