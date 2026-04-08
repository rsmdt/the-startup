---
unit: {{UNIT_ID}}
framework: {{TEST_FRAMEWORK}}
---
# E2E Test Stubs — {{UNIT_TITLE}}

Pre-generated executable test stubs for the evaluation agent. Each stub maps to a scenario and tests observable behavior through the external interface.

## Setup

{{Framework-specific setup: imports, base URL configuration, test client initialization.
Reference the project's existing test patterns.}}

## Stubs

{{For each scenario, generate a test stub:

### {{SCENARIO_NAME}} ({{P0 | P1 | P2}})

```{{language}}
{{Framework-native test code.
- Test name matches scenario name in project convention (camelCase/snake_case)
- Asserts expected outcomes from the scenario's "Expected" section
- Uses project's HTTP client or test utilities
- Marked as pending/skipped until evaluation phase
Example (Jest):
  test.skip('validates SQL injection in query parameter', async () => {
    const response = await fetch('http://localhost:PORT/api/query', {
      method: 'POST',
      body: JSON.stringify({ query: "'; DROP TABLE users; --" }),
    });
    expect(response.status).toBe(400);
    expect(await response.json()).toHaveProperty('error');
  });
}}
```
}}
