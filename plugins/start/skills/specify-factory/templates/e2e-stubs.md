---
unit: {{UNIT_ID}}
framework: {{TEST_FRAMEWORK}}
---
# E2E Test Stubs — {{UNIT_TITLE}}

Pre-generated executable test stubs for the evaluation agent. Each stub maps to a scenario and tests observable behavior through the external interface.

## Setup

{{Framework-specific setup: imports, base URL configuration, test client initialization.
Reference the project's existing test patterns.

CRITICAL — file path naming: when the framework convention places a path
comment at the top of a stub (e.g. Python `# tests/path/test_foo.py`),
or when you suggest where the stub should be placed, name the file after
the **feature or component under test**, NEVER the unit ID. {{UNIT_ID}}
is an orchestration-only label and must not leak into the codebase. The
same rule applies to test function names — derive them from the scenario,
not the unit.

  ✅ # apps/api/tests/integration/features/runs/test_audit_logger.py
  ✅ # apps/web/src/__tests__/chat/test_optimistic_render.tsx
  ❌ # apps/api/tests/integration/features/runs/test_{{UNIT_ID}}_audit.py
  ❌ # apps/web/src/__tests__/chat/test_{{UNIT_ID}}_optimistic.tsx
}}

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
