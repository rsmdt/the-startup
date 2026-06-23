# Plan Validation

Run these checks before finalizing plan.mdx. Report each as pass/fail with the specific gap.

## Coverage checks
- **Requirement coverage** — every requirement / acceptance criterion in requirements.md maps to at least one completion criterion. List any uncovered requirement.
- **Capability coverage** — every capability under `## Completion Criteria` has at least one `P0` scenario.
- **Observable behavior** — every completion criterion tests behavior through an external interface (API/UI/CLI), with no reference to internal class names, function signatures, or file paths.

## Format checks
- **Frontmatter** — `threshold` present and in `(0, 1]`; `title`, `status`, `version` present.
- **Spine** — `## Target State`, `## Quality Gates`, `## Completion Criteria` headings present; criteria are `- [Pn] [title](scenarios/{capability}/{name}.md)` links under `### {Capability}` headings.
- **Scenarios** — each linked scenario file exists and follows templates/scenario.md (frontmatter `capability`/`feature`/`priority`, `## Scenario`, `## Expected`).
- **Quality gates** — each gate is checkable (names a command to run or an observation to make).

## Discipline checks
- **No HOW** — the plan contains no implementation steps, task order, phases, units, execution groups, or tier choice.
- **No duplication** — target state is not restated inside completion criteria; open questions appear only in the trailing section.
- **Self-contained** — no references to chat history or prior drafts.
