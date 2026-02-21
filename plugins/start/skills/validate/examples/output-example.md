# Example Validation Output

## Validation: 003-notification-system

**Mode**: Spec Validation
**Assessment**: 🟡 Needs Attention

### Summary

| Perspective | Pass | Warn | Fail |
|-------------|------|------|------|
| ✅ Completeness | 8 | 2 | 0 |
| 🔗 Consistency | 5 | 1 | 1 |
| 📐 Coverage | 6 | 3 | 0 |
| **Total** | **19** | **6** | **1** |

*🔴 Failures (Must Fix)*

| ID | Finding | Recommendation |
|----|---------|----------------|
| F1 | PRD says "email and SMS" but SDD only designs email *(prd:3.2 vs sdd:4.1)* | Add SMS architecture to SDD Section 4 or update PRD to email-only *(cross-document contradiction)* |

*🟡 Warnings (Should Fix)*

| ID | Finding | Recommendation |
|----|---------|----------------|
| W1 | Missing acceptance criteria for retry logic *(prd:5.3)* | Add specific retry count, backoff strategy, and failure threshold *(requirement exists but is not testable)* |
| W2 | No error handling for webhook timeout *(sdd:6.2)* | Document timeout threshold and fallback behavior *(integration point without failure mode)* |
| W3 | Ambiguous "fast delivery" requirement *(prd:2.1)* | Replace with measurable target: "< 500ms p95 delivery latency" *(vague language scores 0.8 ambiguity)* |
| W4 | PLAN Phase 2 has no test tasks *(plan:Phase 2)* | Add unit test tasks for notification service *(phase implements logic without verification)* |
| W5 | Template variable not resolved *(sdd:7.1)* | Replace [NEEDS CLARIFICATION] with actual monitoring strategy *(placeholder still present)* |
| W6 | Missing edge case: notification opt-out *(prd:3.4)* | Add user preference handling for notification channels *(no requirement for opt-out flow)* |

*✅ Passes*

| Perspective | Verified |
|-------------|----------|
| Completeness | All required sections populated, README tracking current |
| Completeness | No TODO/FIXME markers in PRD or SDD |
| Consistency | Terminology consistent: "notification" used throughout (not mixed with "alert" or "message") |
| Consistency | Cross-references valid: all SDD section links resolve to existing PRD requirements |
| Coverage | All 8 functional requirements have acceptance criteria |
| Coverage | Security considerations documented for PII in notifications |

### Verdict

Spec is 90% ready. One failure (PRD/SDD mismatch on SMS) must be resolved before implementation. 6 warnings are improvement opportunities — the ambiguity warning (W3) and missing test tasks (W4) should be addressed for implementation confidence.
