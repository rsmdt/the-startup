# Example Review Output

Chat summaries for both modes. The full detail always lives in the generated interactive report; the chat summary stays concise.

---

## Scoped Mode: feature/auth-refactor

**Verdict**: 🟡 APPROVE WITH COMMENTS
📄 Full report: `docs/reviews/review-feature-auth-refactor-2026-07-07.html`

### Summary

| Severity | Count |
|----------|-------|
| 🔴 Critical | 0 |
| 🟠 High | 2 |
| 🟡 Medium | 3 |
| ⚪ Low | 1 |

*🔴 Critical & High Findings (Must Address)*

| ID | Finding | Consequence | Remediation |
|----|---------|-------------|-------------|
| H1 | Missing rate limit on login *(auth/routes.ts:45)* | Unlimited login attempts enable credential stuffing | Add express-rate-limit middleware |
| H2 | Error swallows stack trace *(auth/service.ts:78)* | Production auth incidents become undebuggable | Re-throw with original error as cause |

**[H1] Code example**
```typescript
// Before
router.post('/login', authController.login)

// After
const loginLimiter = rateLimit({ windowMs: 15 * 60 * 1000, max: 5 })
router.post('/login', loginLimiter, authController.login)
```

Medium/Low findings (M1–M3, L1) are in the report.

### Strengths

- ✅ JWT implementation follows security best practices with bcrypt + proper salt rounds
- ✅ Clean separation between auth routes, service, and middleware layers
- ✅ Comprehensive test coverage for happy paths (92% coverage)

### Verdict Reasoning

0 critical findings but 2 high findings: a security gap (missing rate limiting) and a debuggability issue (swallowed stack traces). Both are straightforward fixes. 3 medium findings are improvement opportunities that don't block merge.

**Next**: Apply suggested fixes | Create follow-up issues | Proceed without changes

---

## Audit Mode: full codebase

**Verdict paragraph**
📄 Full report: `docs/reviews/review-audit-2026-07-07.html`

> The domain core is small, well-named, and stable — the real risk sits at the edges. Vendor SDK types leak through `billing/` into 14 call sites (measured, not stylistic), so the plausible Stripe→Adyen swap touches 14 files across 5 modules. Co-change data confirms `orders/` and `pricing/` are one module drawn as two: 23 shared commits in 6 months. Error handling is consistent except in the sync path, where three swallowed exceptions would blind an on-call engineer during a partial-sync incident. No security-critical findings.

### Ranked Findings (likelihood × cost)

| ID | Rank driver | Finding | Consequence |
|----|-------------|---------|-------------|
| H1 | HIGH×HIGH | Stripe types in public signatures *(billing/api.ts:12, 14 call sites)* | Vendor swap = 14-file change across 5 modules |
| H2 | HIGH×MED | Swallowed exceptions in sync path *(sync/runner.ts:88,104,131)* | Partial-sync incidents undebuggable on-call |
| M1 | MED×MED | orders/ ↔ pricing/ boundary contradicts co-change *(23 shared commits)* | Every pricing feature is a two-module PR |

### Contextually Acceptable (calibration — mandatory)

| Smell | Location | Reasoning |
|-------|----------|-----------|
| 7-param builder | tests/factory.ts:12 | Test-only; called via named helpers, no ordering-bug victim |
| God-ish config module | config/index.ts | Read-only at boot, never co-changes with features |

### Top Refactors

1. **Wrap vendor types at the billing boundary** — driven by H1 (14 leak sites). Migrate: introduce internal `Charge` type → convert call sites one module at a time → remove SDK import outside `billing/`.
2. **Merge orders/ and pricing/** — driven by co-change data (23 commits). Migrate: move `pricing/` under `orders/pricing/` → collapse the pass-through interface → inline.

(Refactors 3–5 in the report.)

**Next**: Start refactor #1 | Create issues for findings | Deep-dive a finding
