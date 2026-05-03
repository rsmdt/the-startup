---
unit: ve1
feature: health-endpoint
priority: P0
---
# Endpoint Returns 200

## Scenario
Send GET request to http://localhost:3000/health

## Expected
- Response status is 200
- Response Content-Type includes application/json
- Response body contains status, timestamp, and uptime fields
