# Output Example

A worked run of the improve skill against a fictional order-processing codebase. Shows what step 2 produces (a candidate section, as rendered in the markdown report) and how step 3 (the challenge loop) unfolds.

---

## Step 2 — a candidate section (rendered)

> **Collapse the Order intake pipeline** &nbsp; `Strong` `local-substitutable`
>
> **Files** — `order/handler.ts`, `order/validator.ts`, `order/normalizer.ts`, `order/repo.ts`
>
> **Before**
> ```
> OrderHandler → OrderValidator → OrderNormalizer → OrderRepo
>                       ↘ leaks PricingClient ↗
> ```
> **After**
> ```
> ┌─ OrderIntake ───────────────────┐
> │  validate · normalize · persist  │   (internals faded)
> └──────────────────────────────────┘
> ```
>
> **Problem** — four shallow modules each do one trivial step; the real bugs hide in how they're chained, and pricing leaks across the seam.
> **Solution** — collapse the chain into one `OrderIntake` module with a single `accept(rawOrder)` interface; pricing becomes an internal call.
>
> **Wins**
> - locality: chaining bugs concentrate in one module
> - leverage: one interface, every caller
> - interface shrinks; implementation absorbs the wrappers
> - delete 3 shallow unit-test files

The report ends with a **Top recommendation** section pointing at this candidate, then the skill asks: *"Which of these would you like to explore?"*

---

## Step 3 — challenge loop (excerpt)

**User:** Let's explore the Order intake collapse.

**Assistant:** Its dependencies are **local-substitutable** — `OrderRepo` hits Postgres, which we already substitute with PGLite in the test suite. So the seam stays *internal*: we test `OrderIntake.accept()` end to end with PGLite running, no port at the external interface.

The deletion test holds — delete `OrderValidator` and `OrderNormalizer` and their complexity reappears inside every future caller of intake, so they were earning their keep as logic but not as separate **modules**. Behind the new seam: validation, normalization, pricing, persistence. The interface is one entry point: `accept(rawOrder): Promise<Order>`, with invariants (rejects on missing line items, normalizes currency before pricing) documented at the seam.

One question: pricing currently calls `PricingClient` over HTTP to an internal service. That's **remote-but-owned**, not local-substitutable — do we pull it inside `OrderIntake`, or keep it behind a port?

**User:** Keep pricing separate — it's owned by the payments team and changes on their cadence.

**Assistant:** Good — that's a load-bearing reason a future review would need. Want me to record it as an ADR so the next architecture pass doesn't re-suggest folding pricing in?

> 💡 Side effect: this names a real module, **Order intake**, not previously in the glossary. Adding it to `CONTEXT.md`:
>
> ```markdown
> ## Order intake
> The deep module that turns a raw order request into a persisted Order — validation, normalization, pricing call, and persistence behind one `accept` interface.
> ```

**Assistant (handoff):** The agreed shape: one `OrderIntake` module, `accept(rawOrder)` interface, pricing kept behind a port (HTTP adapter in prod, in-memory adapter in tests — two adapters justify the seam). Tests move to the `accept` interface; the three shallow unit-test files get deleted.

Next step: hand this to the **refactor** skill, which will execute the collapse under test protection, one change at a time.
