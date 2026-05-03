---
title: "Notification API"
status: completed
---

# Product Requirements

## Problem
Service-to-service event delivery with at-least-once guarantees and idempotent consumption.

## Features

### Must Have
- F1: REST endpoint to publish a notification with payload and topic
- F2: Persist notifications with delivery state (pending, delivered, failed)
- F3: Background worker that delivers pending notifications via webhook
- F4: Idempotency key support so retries do not double-deliver

## Acceptance Criteria
- AC-1: POST /notifications with valid payload returns 202 with notification id
- AC-2: Invalid payload returns 422 with validation errors
- AC-3: Webhook delivery retries with exponential backoff up to 5 attempts
- AC-4: Duplicate idempotency key returns existing notification id
- AC-5: Failed deliveries after max retries are marked failed and surfaced via GET /notifications/:id
