# Plan Validation

Run these checks before finalizing plan.mdx. Report each as pass/fail with the specific gap.

## Coverage checks
- **Feature coverage** — every acceptance criterion in requirements.md(x) maps to exactly one phase feature (`AC-N.M`). List any AC not placed in a phase, and any AC placed in more than one.
- **Behavior coverage** — every EARS criterion in solution.md(x) maps to exactly one phase's Behavior section, and every phase has at least one EARS behavior statement tracing to technical detail in solution.md(x). List any solution EARS criterion not placed in a phase, any placed in more than one, and flag EARS that invents behavior the solution doesn't support.
- **Traceability index** — the `## Traceability` table lists every requirement AC and every solution EARS criterion exactly once; each row's phase matches where the criterion appears in the phases. Flag any criterion missing from the table or mapped to a different phase than its facet.
- **Scenario coverage** — every phase has at least one `P0` validation scenario.
- **Observable behavior** — every validation scenario tests behavior through an external interface (API/UI/CLI), with no reference to internal class names, function signatures, or file paths.

## Structural checks
- **Vertical slices** — every phase is a user-observable deliverable, not a layer ("data-models phase", "service-layer phase" are invalid).
- **Dependency order** — phase dependencies point only to earlier phases; no cycles.

## Format checks
- **Frontmatter** — `threshold` present and in `(0, 1]`; `title`, `status`, `version` present.
- **Spine** — `## Target State` and `## Phases` headings present; each phase is a `### Phase N: {title}` carrying `#### Features — Acceptance Criteria`, `#### Behavior — EARS`, `#### Quality Gates`, and `#### Validation Scenarios` (the last as `- [Pn] [title](scenarios/phase-N/{name}.md)` links).
- **Scenarios** — each linked scenario file exists and follows templates/scenario.md (frontmatter `phase`/`feature`/`priority`, `## Scenario`, `## Expected`), under scenarios/phase-N/.
- **Quality gates** — each gate is checkable (names a command to run or an observation to make).

## Discipline checks
- **Baseline, not prescription** — the plan contains no task-by-task instructions, code, file edits, or execution-group assignments. Phases say *what is delivered and how it is proven*, not the steps to build it.
- **No strategy** — the plan does not name or assume a decomposition strategy (direct/iterative/decompose).
- **No duplication** — target state is not restated inside phases; open questions appear only in the trailing section.
- **Self-contained** — no references to chat history or prior drafts.
