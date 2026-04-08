---
unit: ve1
feature: echo-endpoint
priority: P1
---
# Echo Missing Message

## Scenario
POST http://localhost:3000/echo with body: {}

## Expected
- Response status is 400
- Response body contains error message
