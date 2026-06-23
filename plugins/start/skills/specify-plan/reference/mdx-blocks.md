# MDX Blocks & Document Quality (Path C)

Spec artifacts (`requirements.mdx`, `solution.mdx`, `plan.mdx`) are **MDX that degrades gracefully to readable Markdown**. We adopt a block *vocabulary* and a *document-quality discipline* — not a hosted renderer. Everything here renders natively on GitHub, in IDE previews, and in docs sites with **zero custom infrastructure**. This is the deliberate Path C scope (see `docs/research/mdx-spec-implement-redesign.md`).

The hard rule: a human gets a richer, more scannable review surface; an agent still reads a stable, machine-parseable spine. Never sacrifice the spine for chrome.

---

## The block vocabulary

Use the lightest block that carries the substance. All of these render without a custom renderer:

| Block | Markdown/MDX form | Use for |
|-------|-------------------|---------|
| Prose | normal Markdown, real bold/italic/`code`/links | the plan narrative |
| Checklist | `- [ ] item` task list | quality gates, completion criteria (renders as interactive checkboxes on GitHub) |
| Callout | `> [!NOTE]`, `> [!IMPORTANT]`, `> [!WARNING]` admonitions | decisions, constraints, framing |
| Table | GitHub tables | comparisons, coverage maps, priority legends |
| Diagram | fenced ```mermaid``` block | architecture, data flow, state — renders natively on GitHub/IDEs |
| Code | fenced ```lang``` block | real data shapes, interface signatures, planned stubs |
| Decision | `> [!IMPORTANT]` titled **Decision:** + optional table of options weighed | a settled choice with rationale |
| Open Questions | one trailing section with a recommended default per question | genuinely unresolved decisions |

Avoid raw JSX/HTML unless it degrades to readable text. Mermaid is preferred over ASCII art for any two-dimensional relationship. Never hand-author a wireframe/prototype/canvas — that is Path A territory we have explicitly deferred.

## Document-quality discipline

Adapted from the BuilderIO `visual-plan` quality bar (the reusable, infra-free part):

- **Outcome-first and self-contained.** State the objective and what "done" means. A reader with no chat history must understand it. No "as discussed", "this revision", "unlike before".
- **Specific, not vague.** Name real files, symbols, endpoints, and data shapes. Replace "make it work" with a checkable condition.
- **Every block carries substance.** Don't add a diagram or callout as decoration. One diagram per real relationship; one callout per real decision or constraint.
- **No duplication.** Say each thing once. Don't restate the target state inside completion criteria, or list open questions in two places.
- **Open questions live in one trailing section.** Use a recommended default for each. A complex plan with zero open questions is fine only when every decision is genuinely settled.
- **Verification exercises the real workflow.** Completion criteria should map to a real user journey through the external interface, not just "unit tests pass".

## The machine-parseable spine

`/implement` and its sub-skills parse these without rendering MDX:

- **Frontmatter** — YAML at the top (`title`, `status`, `threshold`, `version`).
- **Stable headings** — `## Target State`, `## Quality Gates`, `## Completion Criteria`, with `### {Capability}` subheadings.
- **Completion-criteria links** — `- [Pn] [title](scenarios/{capability}/{name}.md)` lines under each capability heading.
- **Task lists** — `- [ ]` / `- [x]` for gates and criteria, so progress is both human-visible and grep-able.

Keep these literal and consistent. MDX components, if ever added (Path A/B), wrap *around* this spine — they never replace it.
