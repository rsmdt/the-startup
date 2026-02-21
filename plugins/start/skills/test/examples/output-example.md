# Example Test Output

## Discovery Results

📋 Test Infrastructure Discovery

Runner: vitest (1.6.0)
Command: npx vitest run
Config: vitest.config.ts

Test Files: 47 files
  - Unit: 38 (src/**/*.test.ts)
  - Integration: 7 (tests/integration/**/*.test.ts)
  - E2E: 2 (tests/e2e/**/*.spec.ts)

Quality Commands:
  - Lint: npx eslint .
  - Typecheck: npx tsc --noEmit
  - Format: npx prettier --check .

---

## Baseline Captured

📊 Baseline Captured

Total: 234 tests
✅ Passing: 231
❌ Failing: 3
⏭️ Skipped: 0

Pre-existing failures (YOU STILL OWN THESE):
1. payment.test.ts:45 — Stripe mock returns wrong status code
2. auth.test.ts:89 — JWT expiry test uses hardcoded timestamp
3. user.test.ts:12 — Missing test database fixture

Note: These failures exist before your changes.
Per the ownership mandate, you are responsible for
fixing these if you proceed with changes in this codebase.

---

## Execution Results

🧪 Test Execution Results

Command: npx vitest run
Duration: 4.2s

Total: 234 tests
✅ Passing: 234
❌ Failing: 0
⏭️ Skipped: 0

All tests passing. Suite is healthy. ✓

---

## Escalation (when needed)

⚠️ Escalation Required

Test: stripe-webhook.integration.test.ts:23
Error: ECONNREFUSED 127.0.0.1:4242

Root Cause: Stripe CLI webhook listener not running locally
Why I can't fix it now: Requires Stripe CLI installation and API key configuration
What's needed: Run `stripe listen --forward-to localhost:3000/webhooks`
Workaround: Skip integration tests with `npx vitest run --exclude tests/integration`

---

## Final Report

🏁 Test Suite Report

Command: npx vitest run
Duration: 4.2s

Results:
  ✅ 234 tests passing
  ⏭️ 0 tests skipped
  ❌ 0 tests failing

Quality:
  Lint: ✅ passing
  Typecheck: ✅ passing
  Format: ✅ clean

Fixes Applied:
1. payment.test.ts:45 — Updated Stripe mock to return correct 200 status
2. auth.test.ts:89 — Replaced hardcoded timestamp with relative Date.now()
3. user.test.ts:12 — Added missing user fixture to test setup

Suite Status: ✅ HEALTHY
