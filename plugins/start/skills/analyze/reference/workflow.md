# Workflow Contract

The workflow script you author must implement this contract exactly. The script has no filesystem,
shell, or network access: it orchestrates and branches on the structured results its agents return.
Every agent call passes a `schema`.

## Script rules

- `meta` is a pure literal: `name`, `description`, `phases` as plain strings — no concatenation, no
  variables, no template interpolation. A non-literal `meta` fails to register.
- Phases: `Research`, `Consolidate`, `Verify`, `Synthesize`.
- Inputs (via `args`): `question` (string), `angles` (array of `{ label, focus }`).
- Constants: `VOTES_PER_CLAIM = 3`, `REFUTATIONS_REQUIRED = 2`.
- Progress labels carry claim text from files and web pages: strip control, bidi-override, and
  zero-width characters and cap them at 40 characters.

## Stages

### 1. Research

One research agent per angle, all concurrently; wait for all (consolidation needs the full set).
Schema:

```
{
  claims: [ { claim: string, source: string, excerpt: string,
              importance: "central"|"supporting"|"tangential", sourceKind: "code"|"web" } ] (required),
  unresolved: string
}
```

No limit on the number of claims. Tag each claim with its angle. Keep each angle's `unresolved` note.
No claims at all → return the summary "No claims were gathered", the angles, and the unresolved notes.

### 2. Consolidate

1. **Exact merge, in code:** normalize each claim (lowercase, collapse whitespace, strip trailing
   punctuation); claims with the same key merge, unioning `sources` and `angles`.
2. **Semantic merge, by agent** (only when more than one claim remains): pass the indexed list. Schema:
   ```
   { groups: [ { claim: string, members: integer[], importance: "central"|"supporting"|"tangential" } ] (required) }
   ```
3. **Rebuild in code** — never from the agent's text. For each group, drop member indices that are out
   of range, not integers, or already claimed by an earlier group. The unique claim takes the group's
   wording, the `excerpt` and `sourceKind` of its most important member, and the union of all members'
   `sources` and `angles`. Any index no group claimed survives as its own claim.
4. If the agent returns nothing, continue with the exact-merged set.

### 3. Verify

For every unique claim, `VOTES_PER_CLAIM` verify agents concurrently, each independent. Schema:

```
{ checked: boolean (required), refuted: boolean (required), evidence: string (required),
  confidence: "high"|"medium"|"low" (required), counterSource: string }
```

Tally per claim in code. A null result or `checked: false` is an errored vote — never a refutation and
never a confirmation. Valid votes are the rest:
- **confirmed** — valid votes ≥ `REFUTATIONS_REQUIRED` and refutations < `REFUTATIONS_REQUIRED`
- **refuted** — refutations ≥ `REFUTATIONS_REQUIRED`
- **unverified** — anything else (voters errored)

The runtime allows at most 1000 agents per workflow. If verifying every claim would exceed that, verify what fits (reserving one agent for
synthesis) and report the rest as unverified with reason "agent cap" — log how many were dropped.

No confirmed claims → return without synthesis. The summary distinguishes: every panel failed
(infrastructure failure — retry), some refuted and some unverified (inconclusive), or all refuted.

### 4. Synthesize

One synthesize agent given only the confirmed claims, with the refuted, unverified, and unresolved items
listed separately as not usable. Schema:

```
{ summary: string (required),
  findings: [ { title, detail, confidence: "high"|"medium"|"low", sources: string[] } ] (required),
  caveats: string (required), openQuestions: string[] }
```

If synthesis returns nothing, return the confirmed claims unmerged with their sources, excerpts, and votes.

## Result

```
{ question, summary, findings, caveats, openQuestions,
  refuted: [ { claim, vote: "<kept>-<refuted>", sources, why } ],
  unverified: [ { claim, sources, reason } ],
  unresolved: [ { angle, note } ],
  stats: { angles, rawClaims, afterExactMerge, uniqueClaims, confirmed, refuted, unverified, findings } }
```

Log the claim count per angle, the merge counts, each claim's vote (`<kept>-<refuted>` and ✓ / ✗ / ?), and
a closing tally.
