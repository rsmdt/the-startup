---
title: "API Rate Limiting"
status: completed
threshold: 0.90
max_iterations: 3
version: 0.1.0
---

# Implementation Plan: API Rate Limiting

## Target State

Every API response carries the key's limit and remaining budget, and keys over their plan's budget receive 429 with retry guidance while other keys are served normally.

## Phases & Outcomes

1. **Per-key rate limiting** — clients see their budget in headers and are throttled with 429 when over it. _Depends on: none._

## Tracking

Tasks and live status: `.start/specs/005-rate-limiting/.tracker/` (`task.js graph --spec 005`).
Each slice task carries its code brief; its scenario child tasks carry the holdout evidence.

## Open Questions

None — all slices and their scenarios are settled.
