---
title: "Sequential Dependency Fixture"
status: completed
---

# Product Requirements

## Problem
Add an export endpoint that produces a CSV from existing user-activity data. The endpoint requires both a query helper (to fetch the data) and a CSV serializer (to format the output).

## Acceptance Criteria
- AC-1: GET /export/activity returns CSV with the columns user_id, action, timestamp
- AC-2: Empty result set returns a CSV with only the header row
- AC-3: Endpoint streams the response (does not buffer entire result set in memory)
