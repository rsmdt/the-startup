# MDX Blocks & Document Quality (Path C)

Spec artifacts (`requirements.mdx`, `solution.mdx`, `plan.mdx`) are **MDX that degrades gracefully to readable Markdown**. We adopt a block *vocabulary* and a *document-quality discipline* — not a hosted renderer. Everything here renders natively on GitHub, in IDE previews, and in docs sites with **zero custom infrastructure**. This is the deliberate Path C scope (see `docs/research/mdx-spec-implement-redesign.md`).

The hard rule: a human gets a richer, more scannable review surface; an agent still reads a stable, machine-parseable spine. Never sacrifice the spine for chrome.

---

## The block vocabulary

Use the lightest block that carries the substance. All of these render without a custom renderer:

| Block | Markdown/MDX form | Use for |
|-------|-------------------|---------|
| Prose | normal Markdown, real bold/italic/`code`/links | the plan narrative |
| Checklist | `- [ ] item` task list | features (ACs), quality gates (renders as interactive checkboxes on GitHub) |
| Callout | `> [!NOTE]`, `> [!IMPORTANT]`, `> [!WARNING]` admonitions | decisions, constraints, framing |
| Table | GitHub tables | comparisons, coverage maps, priority legends |
| Diagram | fenced ```mermaid``` block | architecture, data flow, state — renders natively on GitHub/IDEs |
| Code | fenced ```lang``` block | real data shapes, interface signatures, planned stubs |
| Decision | `> [!IMPORTANT]` titled **Decision:** + optional table of options weighed | a settled choice with rationale |
| Collapsible | `<details><summary>Label</summary>` … `</details>` (blank line after `<summary>`) | progressive disclosure — meta/scaffolding (validation checklists, output schemas), verbose examples, and `(if applicable)` sections that would otherwise drown the narrative |
| Question-form | one trailing table: `\| Question \| Options \| Recommended default \| ☐ \|` | consolidating unresolved decisions into one answerable surface |
| Open Questions | the trailing Question-form table, one row per question, each with a recommended default | genuinely unresolved decisions |

Avoid raw JSX/HTML — the **only** permitted HTML is `<details>`/`<summary>`, because GitHub, IDEs, and docs sites all render it natively and it degrades to readable text everywhere else. No custom components (`<QuestionForm/>`, `<Checklist/>`): those need a renderer (Path B), which we have deferred. Mermaid is preferred over ASCII art for any two-dimensional relationship. Never hand-author a wireframe/prototype/canvas — that is Path A territory we have explicitly deferred.

## Document-quality discipline

Adapted from the BuilderIO `visual-plan` quality bar (the reusable, infra-free part):

- **Outcome-first and self-contained.** State the objective and what "done" means. A reader with no chat history must understand it. No "as discussed", "this revision", "unlike before".
- **Specific, not vague.** Name real files, symbols, endpoints, and data shapes. Replace "make it work" with a checkable condition.
- **Every block carries substance.** Don't add a diagram or callout as decoration. One diagram per real relationship; one callout per real decision or constraint.
- **No duplication.** Say each thing once. Don't restate the target state inside the phases, or list open questions in two places.
- **Open questions live in one trailing section.** Use a recommended default for each. A complex plan with zero open questions is fine only when every decision is genuinely settled.
- **Verification exercises the real workflow.** Completion criteria should map to a real user journey through the external interface, not just "unit tests pass".

## The machine-parseable spine

`/implement` and its sub-skills parse these without rendering MDX:

- **Frontmatter** — YAML at the top (`title`, `status`, `threshold`, `version`).
- **Stable headings** — `## Target State` and `## Phases`, with `### Phase N: {title}` subheadings.
- **Per-phase facets** — each phase carries `#### Features — Acceptance Criteria`, `#### Behavior — EARS`, `#### Quality Gates`, and `#### Validation Scenarios`.
- **Scenario links** — `- [Pn] [title](scenarios/phase-N/{name}.md)` lines under each phase's Validation Scenarios heading.
- **Task lists** — `- [ ]` / `- [x]` for features and gates, so progress is both human-visible and grep-able.

Keep these literal and consistent. MDX components, if ever added (Path A/B), wrap *around* this spine — they never replace it.

**Collapsibles never hide the spine.** Wrap only scaffolding and verbose detail in `<details>` — never a spine heading (`## Target State`, `## Phases`, `### Phase N`), a per-phase facet heading, a scenario-link line, or a feature/gate task-list item. The parser reads raw text and does see inside `<details>`, but burying the spine makes the human review surface and the parse target diverge. Rule of thumb: if `/implement` parses it, it stays at top level.
