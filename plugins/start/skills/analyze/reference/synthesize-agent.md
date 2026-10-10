# Synthesize Agent Prompt

Fill `{question}`, `{count}`, the confirmed-claims block, and the not-usable block.

Confirmed-claims block, one entry per claim:
`[i] claim — vote <kept>-<refuted> · importance · sources — Excerpt: "…" — Verifier confirmed (confidence): evidence`

Not-usable block: the refuted claims with why, the unverified claims with reason, and each angle's
unresolved note — each under its own heading.

---

## Synthesis: write the analysis report

**Question:** {question}

{count} claim(s) survived three-voter adversarial verification. Write the report from these, and only
these.

## Confirmed claims
{confirmed-claims block}

## Not usable as findings
{not-usable block}

## Instructions
1. Group related claims into findings that answer the question. A finding explains a mechanism — how
   and why — not a restatement of one claim. Connect claims across angles where they explain each other.
2. Carry every source through: each finding lists the `file:line` / URL of the claims behind it.
   Introduce no claim that is not in the confirmed list, and add nothing from your own knowledge.
3. Confidence per finding: `high` = several independent sources and unanimous votes; `medium` = one solid
   source or a split vote; `low` = thin or single weak source.
4. `summary`: 3–5 sentences answering the question directly.
5. `caveats`: what is uncertain, which sources were weak, what is time-sensitive, and what was refuted or
   left unverified.
6. `openQuestions`: re-read the question clause by clause and name every part no confirmed claim
   answers.
