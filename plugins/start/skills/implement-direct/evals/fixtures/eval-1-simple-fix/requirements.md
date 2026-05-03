---
title: "Rate-limit header off-by-one"
status: completed
---

# Product Requirements

## Problem
The X-RateLimit-Remaining header is off by one when the rate-limit window resets — clients see "0 remaining" for a fraction of a second when they should see the new window's full quota.

## Acceptance Criteria
- AC-1: When the rate-limit window resets, X-RateLimit-Remaining returns the full new-window quota minus the current request, not the previous window's remaining count
