---
unit: ve1
feature: health-endpoint
priority: P0
---
# Health Endpoint

## Scenario
GET http://localhost:3000/health

## Expected
- Response status is 200
- Response body contains status, timestamp, uptime
