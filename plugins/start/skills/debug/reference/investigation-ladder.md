# Investigation Ladder

When reading the code at the error site doesn't reveal the cause, escalate in this order. Each rung gives you something the previous one couldn't.

The single most common debugging failure is generating a fourth hypothesis instead of descending one rung. The ladder exists to interrupt that pattern.

---

## Rung 1 — Trace the full failing path end-to-end

Read every function the request/computation flows through, in order, from entry point to the failure site. Don't sample; don't skim.

Most bugs labelled "mysterious" are visible in code that wasn't read. The skim-and-guess pattern — read one file, form three hypotheses — fails this rung.

**Tells you:** where the failure could possibly be, ruling out everywhere it can't.

## Rung 2 — Probe logging at boundaries

Insert one log line at: function entry, after each `await` or async boundary, before each `return`, and at every error path. Use a stable correlation id (request id, trace id, or a fresh UUID) so concurrent invocations don't interleave incomprehensibly.

Keep probes minimal — name + correlation id + key state — and remove or guard them when the investigation closes.

**Tells you:** how far the computation got, what arguments arrived, what intermediate values were produced.

**Comparison tracing (sub-technique).** When you have access to a working case alongside the failing one, run both with the same probes and diff the logs. The first divergence is where the bug lives. This often beats lone probe analysis because it lets the working trace tell you what *should* have happened — you don't have to infer it from reading code.

**Caveat for timing-sensitive bugs.** Adding logs can perturb timing enough to mask race conditions. If the bug disappears the moment you instrument, that's a signal the instrumentation itself changed the conditions, not that the bug is gone. Reach for Rung 4 (failing test) or capture state via lower-overhead means (e.g., counters, sampling) instead.

## Rung 3 — State capture at the suspect site

At the suspected failure point, log full inputs and intermediate state — not just "got here". Variable values, collection sizes, branch taken, return codes from external calls. Enough to reconstruct the failure offline.

This is the rung that exposes the gap between what the code saw and what you assumed it saw. The gap is usually the bug.

**Tells you:** what the code actually operated on, vs what you assumed.

## Rung 4 — Failing test

Encode the symptom as a test that fails. This both reproduces the bug deterministically and gives a green/red signal for any intervention.

Required when:
- The bug is hard to reproduce manually (write the test from logs/captured inputs).
- You're about to make a fix and need confidence it actually addresses the cause.
- The fix risks regressing on similar conditions.

A failing test is the bug's *anchor*. Without one, you can't tell whether your fix worked or whether the bug just didn't recur this time.

**Tells you:** whether your fix actually works, and whether it stays working.

## Rung 5 — Bisect the search space

Halve the suspect surface and check which half the bug lives in:

- **Code**: comment out / disable code paths until the bug disappears, then re-enable to confirm.
- **History**: `git bisect` to find the introducing commit when the bug is a regression.
- **Inputs**: progressively shrink the failing input until the bug stops reproducing — the smallest still-failing input is the cleanest diagnostic. The classical algorithm here is *delta debugging* (Zeller): partition the input into n chunks, test removal of each, recurse on the smallest still-failing partition, doubling n when no removal succeeds. The result is a 1-minimal failing input — removing any element makes the bug disappear.

**Tells you:** which half (of code / time / input space) the bug lives in. Iterate to converge.

## Rung 6 — Interactive inspection

When tooling allows: debugger, breakpoint, REPL at the failure site. Inspect live state directly.

Use when state is hard to log faithfully (large objects, generators, closures) or you need to interactively explore relationships you couldn't predict.

---

## Choosing the entry rung

| Situation | Start at |
|---|---|
| Stack trace points at a clear file:line | Rung 1 — often finishes there |
| "Should work but doesn't" / wrong value computed | Rung 2-3 |
| Intermittent or hard to reproduce | Rung 4 — anchor it before doing anything else |
| Regression — worked before, broken now | Rung 5 history bisect first |
| Behaviour depends on values you can't predict from inputs | Rung 6 |

## Stopping rule

If you've formed ≥ 2 hypotheses without descending the ladder, you are speculating. Drop down at least one rung before generating any more hypotheses.

The signal that you're skimming: forming hypotheses about code you haven't read top-to-bottom in the failing path.
