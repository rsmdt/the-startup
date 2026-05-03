---
title: "Notification API"
status: completed
---

# Product Requirements

## Problem
Service-to-service event delivery with idempotent consumption.

## Features

### Must Have
- F1: REST endpoint to publish notifications

## Acceptance Criteria
- AC-1: POST /notifications returns 202 with id
- AC-2: Invalid payload returns 422
- AC-3: Webhook delivery retries up to 5 times
- AC-4: Idempotency key dedupes
- AC-5: Failed deliveries are surfaced via GET
