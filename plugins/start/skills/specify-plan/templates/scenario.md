# {{SCENARIO_TITLE}}

> Priority: {{P0 | P1 | P2}} · Verifies: {{AC_REFERENCE or EARS tag}}

This is the body of a **scenario child task** (`task.js create --type scenario --parent {sliceId}
--body-file …`). It is holdout evidence — the eval agent reads it; the code agent never does.

## Scenario
{{Plain-English description of the action to perform.
Include: HTTP method, endpoint, input data, or UI interaction.
Be specific — reference actual endpoints and data shapes from the codebase.}}

## Expected
{{Plain-English description of the expected observable outcome.
Include: response status, response body structure, UI state change.
Each assertion on its own line. Where it verifies an EARS statement, mirror that response.}}
