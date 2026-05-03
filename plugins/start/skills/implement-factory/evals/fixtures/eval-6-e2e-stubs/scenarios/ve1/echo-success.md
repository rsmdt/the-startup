---
unit: ve1
feature: echo-endpoint
priority: P0
---
# Echo Success

## Scenario
POST http://localhost:3000/echo with body: {"message": "hello world"}

## Expected
- Response status is 200
- Response body contains message "hello world"
- Response body contains timestamp field
