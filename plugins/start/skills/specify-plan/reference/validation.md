# Plan Validation

Run these checks before finalizing. The narrative lives in the Plan part of `spec.md`; the operational units are the
tracker tasks (`task.js graph --spec {NNN}`). Report each check as pass/fail with the specific gap.

## Coverage checks
- **Feature coverage** — every acceptance criterion in the Requirements part maps to exactly one slice's Features (`AC-N.M`). List any AC not placed in a slice, and any AC placed in more than one.
- **Behavior coverage** — every EARS criterion in the Solution part maps to exactly one slice's Behavior. List any solution EARS criterion not placed in a slice, any placed in more than one, and flag EARS that invents behavior the solution doesn't support. A slice whose acceptance criteria already state their behavior needs no EARS of its own.
- **Scenario coverage** — every Feature (AC) maps to ≥1 `business` scenario and every EARS behavior
  (incl. edge/error states) maps to ≥1 `technical` scenario, so no acceptance surface is left untested;
  every slice has ≥1 `P0`. List any Feature or EARS behavior with no scenario. Each scenario carries a
  `Surface` tag (`business` | `technical`) — the eval agent scores the two surfaces separately.
- **Quality-gate coverage** — every `NFR-N` in the Requirements part appears, with its id, as a quality gate on ≥1 slice. List any NFR with no slice.
- **Assumption coverage** — every `[ASSUMPTION An]` cited in the Requirements and Solution parts is resolved or cited on the slices that depend on it.
- **Observable behavior** — every scenario tests behavior through an external interface (API/UI/CLI), with no reference to internal class names, function signatures, or file paths.

## Structural checks
- **Vertical slices** — every slice is a user-observable deliverable, not a layer ("data-models slice", "service-layer slice" are invalid).
- **Sized to the spec** — each slice is small enough for a single code agent to build in one pass, and no smaller than that requires; a small spec has 2–4 slices.
- **Phases** — each phase ends in an outcome a user can see working; it is the unit `/implement` evaluates.
- **Dependency order** — every slice `--depends-on` id refers to a slice that exists; the graph is acyclic. This is the parallelization contract `/implement` topologically sorts. (`task.js` rejects unknown dep ids and parent cycles at create time.)

## Format checks
- **Frontmatter** — `spec.md` has `threshold` in `(0, 1]` and `max_iterations` a positive integer.
- **Narrative** — the Plan part has `### Target state` and `### Phases` (one row per phase: outcome and slice ids) and the tracker link. It carries no per-slice detail — that lives in the tasks.
- **Tracker shape** — one `type: phase` task per plan phase, each depending on the one before; its `type: slice` children each carry a code brief body (Features / Behavior — EARS / Quality Gates); each slice has ≥1 `type: scenario` child task whose body is the scenario.
- **Quality gates** — each gate in a slice's brief is checkable (names a command to run or an observation to make).

## Discipline checks
- **Slices, not prescription** — no task-by-task instructions, code, or file edits. A slice says *what is delivered and how it is proven*, not the steps to build it; dependencies (not hand-assigned groups) decide execution order.
- **No duplication** — the target state is not restated inside the phase outcomes; open questions appear only in the Ledger's Open questions table.
- **Self-contained** — no references to chat history or prior drafts.
