---
unit: ep1
feature: ping
priority: P0
---
# Ping Returns Pong

## Scenario
GET http://localhost:3000/ping

## Expected
- Response status is 200
- Response body contains {"pong": true}
