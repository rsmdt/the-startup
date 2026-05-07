# Complexity Classifier

Heuristic for routing a specification to one of three implementation tiers — **Direct**, **Incremental**, or **Factory** — at step 6 of the specify workflow.

The classifier produces a recommendation. The user always sees the rationale and may override.

---

## Inputs

Read after `requirements.md` and `solution.md` exist:

| Signal | Source | How to compute |
|--------|--------|----------------|
| `feature_count` | requirements.md | Count distinct user-facing capabilities. A "feature" is a coherent user-visible behavior with its own acceptance criteria. Multiple ACs serving one capability count as one feature. |
| `ac_count` | requirements.md | Count of acceptance criteria across all features. |
| `component_count` | solution.md | Count of distinct top-level components or services described in the design. A component is a unit with its own files, module boundary, or service. Modifications to existing components count toward the existing surface — they do not increase `component_count`. |
| `change_type` | requirements.md framing or `$ARGUMENTS` | One of `feature`, `fix`, `refactor`, `doc`. Inferred from how the request is framed. |
| `parallel_markers` | solution.md | Boolean. True if the design explicitly calls out parallel work, concurrent execution, or independent service boundaries. |

---

## Classification

Apply rules top-to-bottom; first match wins.

```
Direct   if change_type in {fix, refactor, doc}
      OR (component_count == 0 AND ac_count <= 2)
      OR (solution modifies only existing surface AND no new modules)

Factory  if feature_count >= 2
      OR component_count >= 3
      OR parallel_markers == true

Incremental otherwise
```

### Edge cases

- **Single feature, many components (3+)**: classify as Factory. The component count signals enough surface area to benefit from parallel units.
- **Many ACs, single component**: classify as Incremental. ACs alone don't justify Factory overhead.
- **Refactor that touches >2 components**: classify as Incremental, not Direct. The breadth of the change warrants phase boundaries even if it isn't a feature.
- **Doc-only change with code examples**: classify as Direct.

---

## Rationale Output

When presenting the recommendation, surface the signals that drove it. Example:

> Classified as **Factory** — solution.md describes 4 components (auth-svc, rate-limiter, metrics-bus, dashboard); requirements.md has 12 ACs across 3 features; parallel work flagged in §6 of solution.md.

> Classified as **Incremental** — single feature with 6 ACs, 2 components (validator + middleware). No parallel work flagged.

> Classified as **Direct** — change_type=fix; modifies existing rate-limit middleware only; 1 acceptance criterion.

The user's confirmation question always shows all three options with the recommended one highlighted.

---

## Override Behavior

- Direct → user upgrades to Incremental or Factory: orchestrator runs `specify-incremental` or `specify-factory` as if classifier had said so.
- Incremental → user downgrades to Direct: orchestrator skips decomposition; spec finishes with just requirements.md and solution.md.
- Incremental → user upgrades to Factory: orchestrator runs `specify-factory`. Any pre-existing `plan/` directory is left in place but flagged as stale in the spec README decision log.
- Factory → user downgrades: same handling — pre-existing `manifest.md`/`units/`/`scenarios/` are left in place and flagged as stale.

The orchestrator never deletes prior artifacts on tier change. Cleanup is manual.

---

## Decision Logging

Whatever tier is chosen (recommended or override), record in the spec README decisions log:

```
| Date | Decision | Rationale |
|------|----------|-----------|
| YYYY-MM-DD | Decomposition tier: Incremental | Classifier recommendation: Incremental (single feature, 2 components). Accepted. |
| YYYY-MM-DD | Decomposition tier: Factory (override) | Classifier recommended Incremental. User chose Factory because of planned cross-team work. |
```

This makes the tier choice auditable and traceable.
