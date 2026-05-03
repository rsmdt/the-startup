---
title: "Auth Platform — Solution"
status: completed
---

# Solution Design

## Components

### Auth Controller
REST endpoints.

### Session Service
Token lifecycle.

### Password Service
Hashing and reset.

### Email Service
Outbound notifications via SES.

## Concurrency
Token issuance and revocation can run in parallel; Email Service runs independently of Session Service.
