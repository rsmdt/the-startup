---
title: "Rate-limit fix — Solution Design"
status: completed
---

# Solution Design

## Component

### Rate-Limit Middleware
Single existing middleware at `src/middleware/rate-limit.ts`. The bug is in the remaining-counter calculation: window-reset detection happens after the counter is computed instead of before. Fix is a one-line reorder plus an added unit test for the boundary case.

## Key Decisions

- **ADR-1**: Fix in place, no new abstractions. The middleware is correct in shape; only the order of two operations is wrong.
