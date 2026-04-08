---
unit: ve1
framework: node:test
---
# E2E Test Stubs — Echo Endpoint

## Setup
```javascript
import { describe, it } from 'node:test';
import assert from 'node:assert/strict';
const BASE_URL = 'http://localhost:3000';
```

## Stubs

### echo-success (P0)
```javascript
it.skip('POST /echo returns echoed message with timestamp', async () => {
  const response = await fetch(`${BASE_URL}/echo`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ message: 'hello world' }),
  });
  assert.strictEqual(response.status, 200);
  const body = await response.json();
  assert.strictEqual(body.message, 'hello world');
  assert.ok(body.timestamp);
});
```

### echo-missing-message (P1)
```javascript
it.skip('POST /echo with empty body returns 400', async () => {
  const response = await fetch(`${BASE_URL}/echo`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({}),
  });
  assert.strictEqual(response.status, 400);
});
```
