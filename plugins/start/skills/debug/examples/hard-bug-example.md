# Hard Bug Example

A counter to `output-example.md`. That example shows a bug whose stack trace pointed at a `file:line`, where one read of the file revealed the cause. This example shows the disciplined investigation of a bug that *resists* that approach — no obvious stack trace, doesn't reliably reproduce, requires instrumentation and explicit hypothesis grading.

Many real bugs are this second kind. The discipline shown here is what turns "transient mystery" into "demonstrated cause."

---

## Phase 0: Reproduce

User: "Sometimes our invoice totals come out a few cents off. Saw it twice this week. Can't trigger it on demand."

"Got it. Before forming theories, I need to check whether this is reproducible — a one-shot symptom can't calibrate any hypothesis I'd form, and I'd just accumulate guesses with no way to choose between them.

Could you tell me:
- The two cases — same customer, same product mix, anything in common?
- Are the off-by amounts random, or do they cluster?
- Do you have logs of the inputs?"

User: "Both involved discount codes. One was -3¢, one was +2¢. Inputs are in the logs."

"Useful handle. Let me try feeding those exact inputs to the totals function locally — if they fail consistently, we have a trigger; if not, we instrument."

[reads `src/invoicing/totals.py`, reconstructs inputs from the log lines, runs locally]

"Both reproduce. ✓

`[evidence]` `invoice_total([items=..., discount=DISCOUNT10, currency=USD])` returns `1037`, expected `1040`. Inputs reconstructed from log line `2026-05-04T08:14Z`.

I have a reliable trigger. Proceeding to Understand."

---

## Phase 1: Understand — trace the full path

"Walking the path top-to-bottom rather than reading one file. The path for an invoice total is:

1. `InvoiceService.compute(invoice)` — `src/invoicing/service.py:42`
2. → `apply_discounts(items, codes)` — `src/invoicing/discounts.py:18`
3. → `total_with_tax(subtotal, region)` — `src/invoicing/totals.py:60`
4. → `round_currency(amount, currency)` — `src/utils/money.py:11`

I'll read all four end-to-end before forming theories — skimming any one of them and guessing is exactly the failure mode the ladder exists to prevent."

[reads all four]

"Two suspicious things:

- `[hypothesis]` H1 — `discounts.py:18` applies the discount to a `float` subtotal before `round_currency` runs. Float arithmetic on small fractional cents could drift.
- `[hypothesis]` H2 — `round_currency` uses banker's rounding (round-half-to-even); discounts producing exactly half-cent values would round inconsistently depending on the previous digit.

Two hypotheses, both plausible. Per the stopping rule I shouldn't generate a third — I should descend the ladder and test these."

---

## Phase 2: Select Mode

"Two hypotheses, reproducible bug, single component. Standard mode is appropriate. Agent Team would be required at ≥ 3 hypotheses or non-reproducibility — neither applies."

---

## Phase 3: Investigate via the ladder

"Already at Rung 1 (full path read). Dropping to Rung 3 — capturing state at the suspect sites — to test H1 and H2."

[adds probe logging to `discounts.py:18` capturing the `subtotal` type and value, and to `round_currency` capturing input/output, with a fresh correlation id]

[runs the failing input]

"Probe results:

`[evidence]` `discounts.py:18` logs `subtotal=10.395, type=float`.
`[evidence]` `round_currency` logs `input=10.395, output=10.40`.

Sanity check: `10.40 - 0.03` (the -3¢ discount) = `10.37`, which is what we'd expect — but the function returned `10.37` from a different code path. Let me re-read the discount application order."

[re-reads `discounts.py:18-30` carefully, finds line 24 — discount applied AFTER `round_currency`, undoing the rounding]

"`[ruled out]` H2 because `round_currency` produced the right value (`10.40`) — rounding isn't the issue.
`[ruled out]` H1 because float drift would be sub-cent and get rounded out; this is whole-cent error.

Both hypotheses dead. Per hygiene I close them explicitly before forming a new one — and I should *re-read* rather than skim, since my first read missed the actual offending line.

Looking at `discounts.py:18-30` properly:

```python
subtotal = round_currency(subtotal, currency)  # line 18
for code in codes:
    discount = code.amount * 1.0  # line 22 — coerces int → float
    subtotal -= discount         # line 23 — applied AFTER rounding
```

`[hypothesis]` H3 — discount amounts (originally cents-as-int) are coerced to float at line 22, then subtracted from a freshly-rounded subtotal. This reintroduces float imprecision into a value that was clean. Order-of-cents error matches the symptom."

---

## Phase 4: Find Root Cause — demonstrate

"H3 is plausible. The bar for promoting it to root cause is `[demonstrated]` — I need to switch the bug on and off by toggling the condition.

Toggling experiment: replace `discount = code.amount * 1.0` with `discount = code.amount` (preserve int), keep everything else identical, re-run both failing inputs."

[runs both failing cases with the toggle]

"`[evidence]` With int preserved at line 22, both inputs return the expected values exactly. Reverting line 22 to `* 1.0` reproduces the off-by-cents on both inputs.

`[demonstrated]` H3 — the float coercion at `discounts.py:22` is the root cause. Toggling that single line switches the bug on and off across both reproductions, with no other changes."

---

## Phase 5: Fix and Verify

"Minimal fix: drop the `* 1.0` coercion at `discounts.py:22`. The discount amount is already an int / Decimal-compatible.

Want me to apply this? I'll also propose a regression test using the two captured failing inputs — they're the test the next person needs."

User: "Yes, both."

[applies fix, adds regression test, runs full suite]

"Applied. Tests passing ✓. The two captured cases are now in `tests/invoicing/test_totals.py::test_discount_int_preservation`.

Done, or want me to check for the same pattern (`* 1.0` coercions of monetary ints) elsewhere?"

---

## What this example illustrates

- **Phase 0 is non-negotiable.** Two cases worth of input data turned a "rare mystery" into a deterministic bug in five minutes. Speculating on the original report would have produced confident-sounding nonsense.
- **Reading the full path before theorizing surfaced two hypotheses, both of which turned out wrong** — but cleanly ruled out, not abandoned. Without the explicit `[ruled out]` markers it would have been easy to drift to a third theory while implying the first two were "still being looked at."
- **The first round of probes falsified both hypotheses, and that was useful**, not a setback. Each `[ruled out]` shrinks the search space.
- **The root cause was promoted only after the toggling experiment.** "Fits the symptoms" is not enough — `[demonstrated]` requires switching the bug on and off.
- **No new tools were needed beyond reading code, adding log lines, and running the function.** The ladder didn't need to descend past Rung 3; it just needed to be descended at all instead of replaced by more guessing.
