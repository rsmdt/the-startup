# Hypothesis Hygiene

The meta-skill of debugging: which claims you make, how you grade them, what to do when challenged.

This is the discipline that separates "I have a theory" from "I have a root cause." Without it, plausible-sounding speculation gets promoted to conclusion and the investigation drifts.

---

## Vocabulary

Prefix every reported claim with one of:

- `[hypothesis]` — fits the symptom; not yet tested.
- `[evidence: X]` — an observation you actually made: a file read, a log line, a test output, a query result. Cite specifically (file:line, exact log line, the command output).
- `[ruled out: X because Y]` — hypothesis tested and falsified. Y is the falsifying observation.
- `[demonstrated]` — you can switch the bug on and off by toggling a condition. The toggling experiment is the proof.

The prefix tells the reader (and you) the epistemic grade. Reporting `[hypothesis]` as if it were `[demonstrated]` is the failure mode this vocabulary exists to prevent.

A root cause requires `[demonstrated]`. Anything else is speculation. If asked to identify the root cause and the strongest candidate is still `[hypothesis]` or `[supported]`, say so explicitly — don't dress speculation as conclusion.

**A claim only counts as `[evidence]` if it's tool-verified.** A file:line citation requires you to have actually read that line. A log line citation requires the log to exist. A test result requires the test to have run. Confidently-stated claims about code behaviour without a citable tool output are still `[hypothesis]`, no matter how plausible — that's the discipline that prevents asserting fluent-sounding speculation as findings.

---

## The hypothesis ledger

Track every hypothesis with TodoWrite. Allowed transitions:

```
pending → supported | ruled out | demonstrated
```

Pivoting from hypothesis A to B is legal only after one of:
- A is marked `ruled out` with the falsifying observation.
- A is marked `unresolved — deferred because <reason>` (e.g., requires infrastructure access you don't have).

**Silent abandonment is not allowed.** Talking about hypothesis A, then quietly switching to B without closing A out, looks like progress but isn't. By hypothesis 4, you have three unfalsified ghosts in the ledger and no clear path to a root cause.

Empty hypotheses get kept around as "would have ruled out" reminders, not deleted — future evidence sometimes resurrects them.

---

## Confidence reporting

When asked "how confident are you?", answer in evidence terms — not adjectives, not percentages:

- "Demonstrated — I can reproduce on demand by X, and the bug disappears when I toggle Y."
- "Strong support — three independent observations consistent with this hypothesis, none contradicting."
- "Plausible — fits the symptom, not yet tested."
- "Speculation — fits the symptom; alternatives are equally plausible."

Use the language of falsification: *what observation would change your mind?* If you can't answer that question, the claim is not yet `[supported]` — it's still `[hypothesis]`.

---

## On user pushback

User pushback is a signal, not a verdict. The instinct to capitulate to look agreeable is the failure mode this section exists to interrupt.

1. **Did their pushback name evidence or reasoning you didn't address?**
   - **Yes** — acknowledge the specific point, update the ledger (probably `ruled out` or downgrade the grade), reformulate. "You're right, I had `[hypothesis]` X but I never checked Y, which falsifies it."
   - **No** — defend the position with reasoning. Capitulating to look agreeable corrupts the investigation. "I hear the concern, but the evidence still supports X because Z. What am I missing?"

2. **Never abandon a hypothesis because the user disliked it.** Only because it's been falsified, or because they pointed at evidence you missed.

3. **If the user is right that you skimmed:** own that specifically. "I had `[hypothesis]` X but I never actually traced past line 50, which is where Y happens — that falsifies X. Let me restart from Rung 1 of the ladder." That's accountability; capitulation is theatre.

4. **If you've been pivoting:** stop. Re-state the ledger explicitly. Show the user every hypothesis you've held, its current grade, and which falsifying observations closed each one. Often the user's frustration is precisely that they've watched you drift without closing anything.

### Sycophancy on user-proposed hypotheses

The mirror image of pushback handling: when the user *proposes a new hypothesis*, the temptation is to adopt it because it came from them. Don't. A user-proposed hypothesis enters the ledger as `[hypothesis]`, not `[evidence]`, and earns its way up through the same falsification path as any other. If the user's proposal contradicts a `[ruled out]` finding, say so and ask what they're seeing that you didn't — don't silently re-open the closed branch and re-validate it because they sound confident.

The agreeable-pivot failure mode runs in both directions: capitulating when challenged, and capitulating when prompted. The defense is the same — evaluate substance, defend or update on evidence, never on social pressure.

---

## Asking for help and reporting status

When you're stuck and need help (from the user, a teammate, or another agent), report **what you observed**, not **what you think is wrong**.

- ✓ "The function returns `1037` for input X; expected `1040`. I traced through the discount path and the rounding path; both look correct in isolation but produce wrong values in combination."
- ✗ "I think there's a bug in the rounding code on line 18."

Leading with your hypothesis poisons the helper — they start checking *your* theory instead of forming their own from the symptoms. The same applies to the ledger when restated: lead with observations, group hypotheses underneath, and let the reader form their own ranking before you reveal yours.

---

## Long investigations: the restart gate

Sustained ad-hoc back-and-forth without closing a hypothesis is a failure signal. After roughly 10 turns or 30 minutes of unstructured discussion — whichever comes first — halt and explicitly restate the ledger:

- Every hypothesis you've held, with its current grade.
- Every observation you've gathered, cited.
- The next experiment, written as: "If I do X, I expect Y; observing Z would falsify hypothesis A."

This is Zeller's logbook discipline. The restart gate exists because long unstructured sessions degrade in two specific ways: the human loses the thread of which hypotheses are still live, and the agent loses earlier context as the conversation grows. Periodically re-snapshotting the ledger pins both back to ground truth. Don't wait until you're already lost — set the trigger at the duration, not at the symptom.

---

## When you're done

You're done when:
- A hypothesis is `[demonstrated]` AND a fix has been verified against the failing test.
- All plausible hypotheses are `ruled out` AND you've documented what you'd need to investigate further (more access, longer observation, different environment).

You're NOT done when:
- A hypothesis "fits the evidence." That's the start of validation, not the end.
- The bug "didn't reproduce in N attempts." That's instrumentation insufficiency, not a verdict of `transient` or `fixed`.
- You have a plausible story but no toggling experiment that confirms cause-and-effect.

Saying "I don't know yet, here's the next rung I'd descend" is a legitimate answer. Saying "my best guess is X" while implying it's a finding is not.
