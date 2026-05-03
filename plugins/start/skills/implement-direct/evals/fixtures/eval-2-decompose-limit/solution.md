---
title: "Auth — Solution Design"
status: completed
---

# Solution Design

## Components

### Auth Controller
REST endpoints for login, refresh, forgot, reset.

### Session Service
Issues, validates, and invalidates access and refresh tokens.

### Password Service
Hashing (Argon2), reset-token generation, single-use enforcement.

### Email Service
Sends password-reset emails via SES.

## Key Decisions

- **ADR-1**: JWT for access tokens, opaque random for refresh tokens.
- **ADR-2**: Argon2id for password hashing — current OWASP recommendation.
