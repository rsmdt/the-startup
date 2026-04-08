---
unit: ve1
feature: health-endpoint
priority: P1
---
# Invalid Method Returns 405

## Scenario
Send POST request to http://localhost:3000/health

## Expected
- Response status is 405
