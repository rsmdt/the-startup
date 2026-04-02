# Example Test Output

## Discovery Results

📋 Test Infrastructure Discovery

Test Categories:

| Category | Runner | Command | Config | Files |
|----------|--------|---------|--------|-------|
| Unit | vitest (1.6.0) | `npx vitest run` | vitest.config.ts | 38 (src/**/*.test.ts) |
| Integration | vitest (1.6.0) | `npx vitest run tests/integration` | vitest.config.ts | 7 (tests/integration/**/*.test.ts) |
| E2E | playwright (1.42.0) | `npx playwright test` | playwright.config.ts | 2 (tests/e2e/**/*.spec.ts) |

Total Test Files: 47

Quality Commands:
  - Lint: npx eslint .
  - Typecheck: npx tsc --noEmit
  - Format: npx prettier --check .

---

## Baseline Captured

📊 Baseline Captured

| Category | Command | Tests | ✅ | ❌ | ⏭️ |
|----------|---------|-------|----|----|----|
| Unit | `npx vitest run` | 195 | 193 | 2 | 0 |
| Integration | `npx vitest run tests/integration` | 27 | 26 | 1 | 0 |
| E2E | `npx playwright test` | 12 | 12 | 0 | 0 |
| **Total** | | **234** | **231** | **3** | **0** |

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

| Category | Command | Duration | Tests | ✅ | ❌ | ⏭️ |
|----------|---------|----------|-------|----|----|----|
| Unit | `npx vitest run` | 2.1s | 195 | 195 | 0 | 0 |
| Integration | `npx vitest run tests/integration` | 1.4s | 27 | 27 | 0 | 0 |
| E2E | `npx playwright test` | 8.7s | 12 | 12 | 0 | 0 |
| **Total** | | **12.2s** | **234** | **234** | **0** | **0** |

Quality:
  Lint: ✅ passing
  Typecheck: ✅ passing
  Format: ✅ clean

Fixes Applied:
1. payment.test.ts:45 — Updated Stripe mock to return correct 200 status
2. auth.test.ts:89 — Replaced hardcoded timestamp with relative Date.now()
3. user.test.ts:12 — Added missing user fixture to test setup

MECE Assessment:
  Overlaps: ⚠️ 1 found — `validateEmail` tested identically in user.test.ts:23 and user.integration.test.ts:67 (consolidate to unit)
  Gaps: ✅ None detected
  Status: ⚠️ Overlaps detected

Suite Status: ✅ HEALTHY
