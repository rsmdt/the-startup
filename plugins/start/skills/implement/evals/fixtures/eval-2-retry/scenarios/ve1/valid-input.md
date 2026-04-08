---
unit: ve1
feature: input-validation
priority: P0
---
# Valid Input

## Scenario
POST http://localhost:3000/validate with body: {"email": "user@example.com", "age": 25}

## Expected
- Response status is 200
- Response body contains {valid: true}
