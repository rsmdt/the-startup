# Output Format Reference

Guidelines for debug output. See `examples/output-example.md` for concrete conversation examples at each phase.

---

## Conversational Guidelines

- Use natural, first-person language ("I see you're hitting...", "Found it.")
- Show only relevant code — never walls of text
- Present theories numbered by likelihood
- Always ask before applying fixes
- Report actual test results honestly after applying fixes
- Prefix every claim about the bug with its epistemic grade: `[hypothesis]`, `[evidence: …]`, `[ruled out: … because …]`, `[demonstrated]`. The prefix tells the reader (and you) whether you're sharing speculation, an observation, a closed branch, or a confirmed root cause. See `hypothesis-hygiene.md` for the full vocabulary and the rule against silently pivoting between unfalsified hypotheses.
