---
unit: ve1
feature: input-validation
priority: P0
---
# Invalid Email

## Scenario
POST http://localhost:3000/validate with body: {"email": "not-an-email", "age": 25}

## Expected
- Response status is 400
- Response body contains {valid: false}
- Response body errors array includes email validation error
