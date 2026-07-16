# Execution Strategy Selector

Heuristic for choosing **how** to implement a specification — **Direct**, **Iterative**, or **Decompose** — at implementation time. This decision lives in `/implement`, not in the specification. The specification (requirements + solution + plan) describes the target state; this selector decides the execution shape.

The selector produces a recommendation. The user always sees the rationale and may override.

---

## The three strategies

| Strategy | Engine skill | How it treats the plan's phases | Best for |
|----------|--------------|---------------------------------|----------|
| **Direct** | implement-direct | **ignores** them — builds straight to the plan's quality gates + validation scenarios | fixes, refactors, doc changes, single-phase work |
| **Iterative** | implement-incremental | **follows** them — adopts the plan's phases as the execution loop, refining each into TDD tasks | one or two phases, transparent human-in-the-loop review |
| **Decompose** | implement-factory | **decomposes** them — splits the phases into finer parallel units with holdout-scenario evaluation + retry | many phases, 3+ components, parallel-eligible work |

The specification's `plan.mdx` already carries a **phase baseline** (vertical-slice phases, each correlating features, EARS behavior, quality gates, and validation scenarios). The strategy decides how to *consume* that baseline — not whether phases exist. The plan's `scenarios/phase-N/` are the holdout set all three strategies validate against.

> The engine skill names are historical (`implement-incremental` = Iterative, `implement-factory` = Decompose). A follow-up will align the names; the routing is by strategy.

## Inputs

Read primarily from plan.mdx (the phase baseline), with solution.md and requirements.md for context:

| Signal | Source | How to compute |
|--------|--------|----------------|
| `phase_count` | plan.mdx `### Phase N` headings | Number of phases in the plan. The strongest signal — the plan already did the vertical slicing. |
| `ac_count` | plan.mdx phase features | Count of acceptance criteria across all phases. |
| `component_count` | solution.md | Distinct top-level components or services. Modifications to existing surface don't increase the count. *Sizing signal only — it sets the strategy, not the slicing axis. Phases are already vertical slices.* |
| `change_type` | requirements framing / `$ARGUMENTS` | One of `feature`, `fix`, `refactor`, `doc`. |
| `parallel_markers` | solution.md / plan phase dependencies | True if the design calls out parallel work, or the plan has independent phases (no shared dependency chain). |

## Classification

Apply rules top-to-bottom; first match wins.

```
Direct    if change_type in {fix, refactor, doc}
       OR (phase_count <= 1 AND component_count == 0 AND ac_count <= 2)
       OR (solution modifies only existing surface AND no new modules)

Decompose if phase_count >= 3
       OR component_count >= 3
       OR parallel_markers == true

Iterative otherwise   (typically phase_count of 1–2)
```

### Edge cases
- **Single phase, 3+ components**: Decompose. The surface area benefits from parallel units even within one phase.
- **Many ACs, one or two phases, single component**: Iterative. ACs alone don't justify Decompose overhead.
- **Refactor spanning several phases**: Iterative, not Direct. Breadth warrants the phase loop.
- **Doc-only with code examples**: Direct.
- **Independent phases (no dependency chain)**: Decompose — they parallelize cleanly.

## Rationale Output

Surface the signals that drove the recommendation. Examples:

> Recommended **Decompose** — plan.mdx has 4 phases; solution.md describes 4 components; parallel work flagged in §6 of solution.md.

> Recommended **Iterative** — plan.mdx has 2 dependency-ordered phases, 2 components. No parallel work flagged.

> Recommended **Direct** — change_type=fix; single phase; modifies existing middleware only; 1 acceptance criterion.

The confirmation question always shows all three strategies with the recommended one highlighted. The user may override freely.

## Decision logging

Record the chosen strategy in `.impl/strategy.md` and in the spec README decisions log:

```
| Date | Decision | Rationale |
|------|----------|-----------|
| YYYY-MM-DD | Execution strategy: Iterative | Selector recommended Iterative (2 phases, 2 components). Accepted. |
| YYYY-MM-DD | Execution strategy: Decompose (override) | Selector recommended Iterative. User chose Decompose for planned parallel work. |
```
