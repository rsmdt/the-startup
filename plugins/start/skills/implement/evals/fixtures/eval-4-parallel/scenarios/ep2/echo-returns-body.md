---
unit: ep2
feature: echo
priority: P0
---
# Echo Returns Body

## Scenario
POST http://localhost:3000/echo with body: {"message": "hello"}

## Expected
- Response status is 200
- Response body contains {"message": "hello"}
