# Consolidate Agent Prompt

Fill `{question}`, `{count}`, and the indexed claim list (`[i] claim — source (sourceKind, importance)`).

---

## Consolidate: merge duplicate claims

Research question: "{question}"

Independent researchers worked separate angles and overlapped. Below are all {count} claims, indexed.
Group the ones that assert the same thing.

## Claims
{indexed claim list}

## Rules
- Two claims merge only if they assert the same fact. Different facts about the same file, function, or
  page are different claims — merge by assertion, not by topic.
- A claim more specific than another is not its duplicate. Keep the specific one separate.
- Every index appears in exactly one group. A claim nothing duplicates is its own group of one. Do not
  drop, filter, or rank anything.
- For each group give the canonical wording: the clearest, most specific phrasing of the shared
  assertion. Prefer an existing claim's wording.
- Do not restate sources or excerpts — return indices; the originals are carried over verbatim.
