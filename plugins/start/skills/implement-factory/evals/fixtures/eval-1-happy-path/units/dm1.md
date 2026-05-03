---
id: dm1
title: Data Model
type: feature
dependencies: []
---
# Data Model

## Goal
Create a simple health status data model.

## Requirements
- Define a HealthStatus type with fields: status (string), timestamp (ISO string), uptime (number in seconds)
- Export a function getHealthStatus() that returns current health data
- Status should be "healthy" when uptime > 0

## Constraints
- Use TypeScript or JavaScript
- No external dependencies
- Full test coverage for all requirements
