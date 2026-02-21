# Output Format Reference

Templates for each phase of test execution.

---

## Discovery Results

```
📋 Test Infrastructure Discovery

Runner: [name] ([version if available])
Command: [exact command to run]
Config: [config file path]

Test Files: [count] files
  - Unit: [count] ([pattern])
  - Integration: [count] ([pattern])
  - E2E: [count] ([pattern])

Quality Commands:
  - Lint: [command or "not found"]
  - Typecheck: [command or "not found"]
  - Format: [command or "not found"]
```

## Baseline Capture

```
📊 Baseline Captured

Total: [N] tests
✅ Passing: [N]
❌ Failing: [N]
⏭️ Skipped: [N]

[If failures exist:]
Pre-existing failures (YOU STILL OWN THESE):
1. [test name] — [brief error]
2. [test name] — [brief error]

Note: These failures exist before your changes.
Per the ownership mandate, you are responsible for
fixing these if you proceed with changes in this codebase.
```

## Execution Results

```
🧪 Test Execution Results

Command: [exact command run]
Duration: [time]

Total: [N] tests
✅ Passing: [N]
❌ Failing: [N]
⏭️ Skipped: [N]

[If all pass:]
All tests passing. Suite is healthy. ✓

[If failures:]
Failures requiring attention:

FAILURE:
- status: FAIL
- category: YOUR_CHANGE | OUTDATED_TEST | TEST_BUG | MISSING_DEP | ENVIRONMENT | CODE_BUG
- test: [test name]
- location: [file:line]
- error: [one-line error message]
- action: [what you will do to fix it]
```

## Escalation Template

Only acceptable for: external service down, infrastructure requirements, permission/access issues.

```
⚠️ Escalation Required

Test: [test name] ([file:line])
Error: [exact error]

Root Cause: [what you found after investigation]
Why I can't fix it now: [specific technical blocker]
What's needed: [concrete next step]
Workaround: [if any temporary measure is possible]
```

## Final Report

```
🏁 Test Suite Report

Command: [exact command]
Duration: [time]

Results:
  ✅ [N] tests passing
  ⏭️ [N] tests skipped
  ❌ 0 tests failing

Quality:
  Lint: ✅ passing | ❌ [N] issues fixed
  Typecheck: ✅ passing | ❌ [N] errors fixed
  Format: ✅ clean | ❌ [N] files formatted

[If fixes were made:]
Fixes Applied:
1. [file:line] — [what was fixed and why]
2. [file:line] — [what was fixed and why]

[If escalations exist:]
Escalations: [N] tests require external resolution
(see details above)

Suite Status: ✅ HEALTHY | ⚠️ NEEDS ATTENTION
```
