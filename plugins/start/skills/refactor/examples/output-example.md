# Example Refactor Output

## Baseline

📊 Refactoring Baseline

Tests: 47 passing, 0 failing
Coverage: 82%
Uncovered areas: src/services/billing.ts (lines 45-89)

Baseline Status: READY

---

## Analysis Summary

## Refactoring Analysis: src/services/billing.ts

### Summary

| Perspective | High | Medium | Low |
|-------------|------|--------|-----|
| 🔧 Code Smells | 1 | 2 | 1 |
| 🔗 Dependencies | 0 | 1 | 0 |
| 🧪 Test Coverage | 1 | 0 | 0 |
| 🏗️ Patterns | 0 | 1 | 0 |
| ⚠️ Risk | 0 | 1 | 0 |
| **Total** | **2** | **5** | **1** |

*🔴 High Impact Issues*

| ID | Finding | Remediation | Risk |
|----|---------|-------------|------|
| H1 | 75-line calculateTotal method *(billing.ts:23-98)* | Extract into calculateSubtotal, applyDiscounts, applyTax *(single method handles subtotal, discounts, tax, and rounding)* | Low — well-tested happy path |
| H2 | No tests for discount edge cases *(billing.test.ts)* | Add tests before refactoring discount logic *(stacked discounts, expired coupons, negative totals untested)* | Medium — refactoring without tests risks regression |

*🟡 Medium Impact Issues*

| ID | Finding | Remediation | Risk |
|----|---------|-------------|------|
| M1 | Magic numbers throughout *(billing.ts:45, 67, 82)* | Extract to named constants *(0.08 tax rate, 0.15 max discount, 100 rounding factor)* | Low |
| M2 | Tight coupling to Stripe *(billing.ts:112)* | Extract PaymentGateway interface *(direct Stripe SDK calls in business logic)* | Medium — interface change |
| M3 | Duplicated validation *(billing.ts:25, billing.ts:55)* | Extract validateLineItems helper *(same null/empty check twice)* | Low |

---

## Error Recovery

⚠️ Refactoring Failed

Refactoring: Extract applyDiscounts from calculateTotal
Reason: Tests failing

Reverted: ✓ Working state restored

Options:
1. Try alternative approach
2. Add missing tests first
3. Skip this refactoring
4. Get guidance

Awaiting your decision...

---

## Completion Summary

## Refactoring Complete: src/services/billing.ts

**Status**: Complete

### Before / After

| File | Before | After | Technique |
|------|--------|-------|-----------|
| billing.ts | 75-line calculateTotal | 4 functions, 18 lines each | Extract Method |
| billing.ts | 3 magic numbers | Named constants in config | Replace Magic Number |
| billing.ts | Duplicated validation | Single validateLineItems | Extract Method |

### Verification

- Tests: 52 passing (baseline: 47, +5 new)
- Behavior: Preserved ✓
- Coverage: 91% (baseline: 82%)

### Quality Improvements

- Cyclomatic complexity reduced from 12 to 4 per function
- Each function now has single responsibility
- Discount logic independently testable

### Skipped

- billing.ts:112 — Stripe coupling (user declined — planned for separate PR)
