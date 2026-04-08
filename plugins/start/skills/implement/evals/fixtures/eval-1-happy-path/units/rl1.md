---
id: rl1
title: Rate Limiter
type: feature
dependencies: [ve1]
---
# Rate Limiter

## Goal
Add rate limiting to the health endpoint.

## Requirements
- Limit to 10 requests per minute per IP
- Return 429 Too Many Requests when limit exceeded
- Include Retry-After header in 429 responses

## Constraints
- Implement as Express middleware
- No external rate limiting libraries
- Full test coverage for all requirements
