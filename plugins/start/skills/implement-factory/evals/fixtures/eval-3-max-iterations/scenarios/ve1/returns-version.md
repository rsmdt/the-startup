---
unit: ve1
feature: status-endpoint
priority: P0
---
# Returns Version

## Scenario
GET http://localhost:3000/status

## Expected
- Response status is 200
- Response body includes "version" field with a semantic version string (e.g., "1.0.0")
