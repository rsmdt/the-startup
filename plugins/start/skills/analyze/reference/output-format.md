# Output Format Reference

Guidelines for analysis output. See `examples/output-example.md` for concrete rendered examples.

---

## Finding Depth Standard

Every finding must include three layers. A finding missing any layer is incomplete.

| Layer | Question Answered | Example |
|-------|-------------------|---------|
| **What** | What exists? | "Payment retry logic at `src/services/payment.ts:112`" |
| **How** | How does it actually work? | "Uses exponential backoff: delays = [1s, 2s, 4s], catches only `TimeoutError` and `NetworkError`, resets circuit breaker on 3rd failure" |
| **Why / So What** | Why does this matter? What are the implications? | "The catch filter means `ValidationError` from Stripe's API bypasses retry entirely — this is correct for idempotency but means malformed requests fail silently" |

**Surface-level findings to avoid:**
- "Uses caching" → must explain: what cache, what invalidation strategy, what TTL, what happens on miss
- "Follows repository pattern" → must explain: what abstraction it provides, what it hides, how queries are composed
- "Has error handling" → must explain: what errors are caught, what's propagated, what's logged vs swallowed

## Recommendation Format

When findings surface problems or opportunities, recommendations follow a strict order:

### 1. The Correct Approach (always present)

Describe the architecturally clean solution:
- What it looks like conceptually
- Which files/modules are affected
- What the migration path involves
- Estimated scope (small/medium/large with explanation)
- What risks or unknowns remain

### 2. Open Questions (always present)

What does the user need to decide or clarify before committing to the approach?

### 3. Alternative Approaches (only when user requests)

Only present if the user explicitly asks after seeing the correct approach. When presenting alternatives, be clear about what is sacrificed compared to the clean solution.

**Never lead with:** "A pragmatic approach would be...", "A hybrid solution...", "To minimize changes...", "A quick win would be..."

## Cycle Summary Structure

Each cycle summary should include:

1. **Mechanism Findings** — grouped by theme, each with What/How/Why layers
2. **Cross-Cutting Observations** — connections between findings, cause-effect chains
3. **Recommendations** — clean approach first, with scope and implications
4. **Open Questions** — what needs clarification or deeper investigation

## Next Steps Options

Use `AskUserQuestion` based on analysis state:

**After each cycle:**
- "Continue to next area"
- "Go deeper on [specific finding]"
- "Persist findings to docs/"
- "Complete analysis"

**After final summary:**
- "Save documentation to docs/"
- "Skip documentation"
- "Export as markdown"
