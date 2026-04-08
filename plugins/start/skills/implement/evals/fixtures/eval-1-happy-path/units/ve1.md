---
id: ve1
title: Validation Endpoint
type: feature
dependencies: [dm1]
---
# Validation Endpoint

## Goal
Create a GET /health endpoint that returns the health status.

## Requirements
- GET /health returns 200 with JSON body containing status, timestamp, uptime
- GET /health returns Content-Type: application/json
- Invalid methods (POST, PUT, DELETE) return 405

## Constraints
- Use Express.js (see AGENTS.md)
- Follow existing controller patterns
- Full test coverage for all requirements
