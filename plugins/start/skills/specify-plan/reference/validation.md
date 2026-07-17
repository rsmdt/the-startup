# Plan Validation

Run these checks before finalizing. The narrative lives in `plan.mdx`; the operational units are the
tracker tasks (`task.js graph --spec {NNN}`). Report each check as pass/fail with the specific gap.

## Coverage checks
- **Feature coverage** — every acceptance criterion in requirements.md(x) maps to exactly one slice's Features (`AC-N.M`). List any AC not placed in a slice, and any AC placed in more than one.
- **Behavior coverage** — every EARS criterion in solution.md(x) maps to exactly one slice's Behavior, and every slice has at least one EARS behavior statement tracing to technical detail in solution.md(x). List any solution EARS criterion not placed in a slice, any placed in more than one, and flag EARS that invents behavior the solution doesn't support.
- **Scenario coverage** — every Feature (AC) maps to ≥1 `business` scenario and every EARS behavior
  (incl. edge/error states) maps to ≥1 `technical` scenario, so no acceptance surface is left untested;
  every slice has ≥1 `P0`. List any Feature or EARS behavior with no scenario. Each scenario carries a
  `Surface` tag (`business` | `technical`) — the eval agent scores the two surfaces separately.
- **Observable behavior** — every scenario tests behavior through an external interface (API/UI/CLI), with no reference to internal class names, function signatures, or file paths.

## Structural checks
- **Vertical slices** — every slice is a user-observable deliverable, not a layer ("data-models slice", "service-layer slice" are invalid).
- **Atomic** — each slice is small enough for a single code agent to build in one pass.
- **Dependency order** — every slice `--depends-on` id refers to a slice that exists; the graph is acyclic. This is the parallelization contract `/implement` topologically sorts. (`task.js` rejects unknown dep ids and parent cycles at create time.)

## Format checks
- **Frontmatter** — `plan.mdx` has `threshold` in `(0, 1]`, `max_iterations` a positive integer, and `title`, `status`, `version`.
- **Narrative** — `plan.mdx` has `## Target State`, `## Phases & Outcomes` (one outcome line per slice + its dependencies), and a `## Tracking` link. It carries no per-slice detail — that lives in the tasks.
- **Tracker shape** — one `type: phase` root; its `type: slice` children each carry a code brief body (Features / Behavior — EARS / Quality Gates); each slice has ≥1 `type: scenario` child task whose body is the scenario.
- **Quality gates** — each gate in a slice's brief is checkable (names a command to run or an observation to make).

## Discipline checks
- **Slices, not prescription** — no task-by-task instructions, code, or file edits. A slice says *what is delivered and how it is proven*, not the steps to build it; dependencies (not hand-assigned groups) decide execution order.
- **No duplication** — the target state is not restated inside the phase outcomes; open questions appear only in the trailing `plan.mdx` section.
- **Self-contained** — no references to chat history or prior drafts.
