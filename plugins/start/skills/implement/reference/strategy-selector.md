# Execution Strategy Selector

Heuristic for choosing **how** to implement a specification — **Direct**, **Iterative**, or **Decompose** — at implementation time. This decision lives in `/implement`, not in the specification. The specification (requirements + solution + plan) describes the target state; this selector decides the execution shape.

The selector produces a recommendation. The user always sees the rationale and may override.

---

## The three strategies

| Strategy | Engine skill | Execution shape | Best for |
|----------|--------------|-----------------|----------|
| **Direct** | implement-direct | 1–3 delivery units, parallel/sequential subagents, no phase plan | fixes, refactors, doc changes, single-AC features |
| **Iterative** | implement-incremental | linear phases generated at exec time, human-in-the-loop checkpoints | single feature with one or two components |
| **Decompose** | implement-factory | parallel units + holdout-scenario evaluation + retry loop, generated at exec time | multi-feature, 3+ components, parallel-eligible work |

> The engine skill names are historical (`implement-incremental` = Iterative, `implement-factory` = Decompose). A follow-up will align the names; the routing is by strategy.

## Inputs

Read from requirements.md, solution.md, and plan.mdx (whichever exist):

| Signal | Source | How to compute |
|--------|--------|----------------|
| `feature_count` | requirements.md / plan capabilities | Distinct user-facing capabilities. The plan's `### {Capability}` headings are a strong signal. |
| `ac_count` | requirements.md / plan completion criteria | Count of acceptance criteria / completion criteria. |
| `component_count` | solution.md | Distinct top-level components or services. Modifications to existing surface don't increase the count. *Sizing signal only — it sets the strategy, not the slicing axis. Decomposition always slices vertically by capability, never by layer.* |
| `change_type` | requirements framing / `$ARGUMENTS` | One of `feature`, `fix`, `refactor`, `doc`. |
| `parallel_markers` | solution.md | True if the design calls out parallel work, concurrent execution, or independent service boundaries. |

## Classification

Apply rules top-to-bottom; first match wins.

```
Direct    if change_type in {fix, refactor, doc}
       OR (component_count == 0 AND ac_count <= 2)
       OR (solution modifies only existing surface AND no new modules)

Decompose if feature_count >= 2
       OR component_count >= 3
       OR parallel_markers == true

Iterative otherwise
```

### Edge cases
- **Single feature, 3+ components**: Decompose. The surface area benefits from parallel units.
- **Many ACs, single component**: Iterative. ACs alone don't justify Decompose overhead.
- **Refactor touching >2 components**: Iterative, not Direct. Breadth warrants phase boundaries.
- **Doc-only with code examples**: Direct.

## Rationale Output

Surface the signals that drove the recommendation. Examples:

> Recommended **Decompose** — solution.md describes 4 components; plan.mdx has 12 completion criteria across 3 capabilities; parallel work flagged in §6 of solution.md.

> Recommended **Iterative** — single capability with 6 completion criteria, 2 components. No parallel work flagged.

> Recommended **Direct** — change_type=fix; modifies existing middleware only; 1 completion criterion.

The confirmation question always shows all three strategies with the recommended one highlighted. The user may override freely.

## Decision logging

Record the chosen strategy in `.impl/strategy.md` and in the spec README decisions log:

```
| Date | Decision | Rationale |
|------|----------|-----------|
| YYYY-MM-DD | Execution strategy: Iterative | Selector recommended Iterative (single capability, 2 components). Accepted. |
| YYYY-MM-DD | Execution strategy: Decompose (override) | Selector recommended Iterative. User chose Decompose for planned parallel work. |
```
