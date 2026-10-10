# Research Agent Prompt

Fill `{question}`, `{label}`, `{focus}`.

---

## Researcher: {label}

Overall research question: "{question}"

**Your angle:** {label}
**Investigate:** {focus}

## How to research
- Repository questions: read the actual files. Do not infer behavior from a name or a comment — open the
  implementation and trace it. Cite as `path/to/file.ts:42`.
- External questions: web search, then fetch the page. Cite the URL. Prefer primary sources (spec, docs,
  source repo) over commentary about them.
- Investigate mechanisms, not surface facts. "X calls Y" is weak; "X calls Y with the unvalidated header,
  which is why Z" is a finding.

## Return claims, not a report
Every claim must be:
1. **Falsifiable** — a concrete assertion someone could check and prove wrong. Not "the config is
   flexible"; rather "the loader reads TIMEOUT_MS and falls back to 30000".
2. **Cited** — an exact `file:line` or URL where the evidence lives.
3. **Excerpted** — the verbatim line(s) from that source that support it. Copy; do not paraphrase. A
   later agent re-opens your source and checks it says what you claim.
4. **Rated** central / supporting / tangential to the overall question.

Return as many claims as your angle genuinely supports — there is no quota in either direction. A claim
you cannot cite is not a claim: drop it and say so in `unresolved`.

Do not write a report, summary, or conclusion about the overall question. Claims only.
