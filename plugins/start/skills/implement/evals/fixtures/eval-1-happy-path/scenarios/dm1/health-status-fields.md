---
unit: dm1
feature: health-status
priority: P0
---
# Health Status Fields

## Scenario
Call getHealthStatus() and inspect the returned object.

## Expected
- Response contains "status" field with value "healthy"
- Response contains "timestamp" field with a valid ISO date string
- Response contains "uptime" field with a number greater than 0
