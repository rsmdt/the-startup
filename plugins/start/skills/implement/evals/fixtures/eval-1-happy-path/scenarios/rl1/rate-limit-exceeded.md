---
unit: rl1
feature: rate-limiting
priority: P0
---
# Rate Limit Exceeded

## Scenario
Send 11 GET requests to http://localhost:3000/health in rapid succession from the same IP.

## Expected
- First 10 requests return 200
- 11th request returns 429
- 429 response includes Retry-After header
