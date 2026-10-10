---
title: "API Rate Limiting"
status: completed
version: "1.0"
---

# Solution Design Document

## Approach

An Express middleware applies a sliding-window counter per API key in the shared Redis cluster, reading each key's limit from its plan.

- **Cross-cutting:** fail open if Redis is unreachable; log every 429
- **Deployment:** feature flag `rate_limiting`

## Read first

| Source | Why it matters |
|--------|----------------|
| `src/middleware/auth.js` | Resolves the API key and plan |

## Boundaries

- **Must preserve:** response bodies of existing endpoints
- **Can modify:** middleware chain in `src/app.js`
- **Must not touch:** billing plans

## Components

| Component | Responsibility | Owns |
|-----------|----------------|------|
| RateLimitMiddleware | Counts requests per key, rejects over-limit, sets headers | FR-1, FR-2 |
| PlanLimitResolver | Returns the limit for a key's plan | plan limits |

## Interfaces & data

| Name | Direction | Type / format | Auth | Contract |
|------|-----------|---------------|------|----------|
| All `/api/*` responses | outbound | HTTP headers | API key | `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `Retry-After` on 429 |
| Redis | outbound | RESP | internal | `ratelimit:{key}:{window}` counters |

## Runtime

After auth, the middleware increments the key's window counter, compares it with the plan limit, and either rejects with 429 or sets headers and continues.

## Decisions

> [!IMPORTANT]
> **Decision (ADR-1): Counter store → Shared Redis sliding window**
>
> | Option | Trade-off / why (not) |
> |--------|-----------------------|
> | **Redis sliding window** | Accurate across instances |
> | In-memory per instance | Over-admits by instance count |
>
> **Rationale:** NFR-2 accuracy across instances · **Accepts:** one Redis round trip per request
>
> - [x] User confirmed

## Behavior

**Main flow — AC-1.1:**
- [ ] WHILE a key is under its limit, THE SYSTEM SHALL pass the request to the handler unchanged

**Error handling — AC-1.2:**
- [ ] IF a key exceeds its limit, THEN THE SYSTEM SHALL respond 429 with `Retry-After` set to the seconds until the window resets

**Main flow — AC-2.1:**
- [ ] THE SYSTEM SHALL set `X-RateLimit-Limit` and `X-RateLimit-Remaining` on every served response

## Risks & gotchas

- Fail-open during Redis outages allows bursts.

## Open questions

None — all decisions settled.
