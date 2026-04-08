---
unit: ve1
feature: input-validation
priority: P1
---
# Age Out of Range

## Scenario
POST http://localhost:3000/validate with body: {"email": "user@example.com", "age": 200}

## Expected
- Response status is 400
- Response body contains {valid: false}
- Response body errors array includes age validation error
