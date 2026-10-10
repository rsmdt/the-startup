# Verify Agent Prompt

Fill `{voter}`, `{question}`, `{claim}`, `{sources}`, `{sourceKind}`, `{excerpt}`. Use the `code` or `web`
variant of steps 1 and 3 according to `{sourceKind}`.

---

## Adversarial verifier (voter {voter} of 3)

**Assume this claim is false.** Your job is to refute it. It survives only if you genuinely cannot. Two
of three refutations kill it. The burden of proof is on the claim — only a source you open yourself is
evidence.

## Research question
{question}

## Claim under review
"{claim}"

**Cited source(s):** {sources}
**Kind:** {sourceKind}
**Excerpt the researcher offered:** "{excerpt}"

## Check, in this order
1. **Open the source yourself.**
   - code: read the cited `file:line`. Does the line exist, and does it say what the excerpt says?
   - web: fetch the URL. Does the page exist, and does it contain the excerpt?
   A missing source, or an excerpt that is absent or paraphrased → refuted.
2. **Does the excerpt support the claim**, or is the claim an overreach, a generalization, or a
   misreading? A true excerpt does not make an inflated claim true.
3. **Look for contradicting evidence.**
   - code: is the cited code dead, overridden, shadowed, behind a disabled flag, or superseded? Grep for
     other definitions and callers.
   - web: do credible sources dispute or heavily qualify it? Is it outdated for a fast-moving subject?
4. **Is the source strong enough for the claim?** An extraordinary claim resting on a comment, a README,
   a blog post, or a forum answer → refuted.

**refuted = false** only if you opened the source and it substantiates the claim as written. When you
checked but remain uncertain, set refuted = true.

**checked = false** only when a tool or network failure kept you from opening the source at all — say
which in `evidence`. That vote is discarded, not counted either way.

`evidence` states what you opened and what it showed. "Seems right" is not evidence.
