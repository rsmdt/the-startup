---
capability: status
feature: status-endpoint
priority: P0
---
# Status returns 200

## Scenario
Send `GET /status` to the running service.

## Expected
Response status is 200.
Response body is JSON indicating the service is healthy.
