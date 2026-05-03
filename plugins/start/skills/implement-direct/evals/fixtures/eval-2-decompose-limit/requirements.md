---
title: "Multi-feature spec — direct mode should reject"
status: completed
---

# Product Requirements

## Features

### Must Have
- F1: User authentication with email/password
- F2: Session management with refresh tokens
- F3: Password reset via email link

## Acceptance Criteria
- AC-1: POST /auth/login with valid credentials returns access + refresh tokens
- AC-2: Invalid credentials return 401 with no token
- AC-3: Access token expires in 15 minutes
- AC-4: Refresh token expires in 30 days and is single-use
- AC-5: POST /auth/refresh issues new access token
- AC-6: POST /auth/forgot sends a password-reset email if account exists
- AC-7: Reset link is single-use and expires in 1 hour
- AC-8: Reset endpoint sets new password and invalidates all existing sessions
