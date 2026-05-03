---
unit: ep2
feature: echo
priority: P1
---
# Echo Missing Body

## Scenario
POST http://localhost:3000/echo with empty body

## Expected
- Response status is 400
