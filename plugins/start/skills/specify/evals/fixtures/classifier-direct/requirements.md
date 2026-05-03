---
title: "Rate-limit header fix"
status: completed
---

# Product Requirements

## Problem
Off-by-one bug in rate-limit middleware on window reset.

## Acceptance Criteria
- AC-1: After window reset, X-RateLimit-Remaining shows new-window quota minus current request

Change type: fix.
