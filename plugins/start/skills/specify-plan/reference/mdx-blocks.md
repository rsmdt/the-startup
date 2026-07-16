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
- **No duplication.** Say each thing once. Don't restate the target state inside the phase outcomes, or list open questions in two places.
- **Open questions live in one trailing section.** Use a recommended default for each. A complex plan with zero open questions is fine only when every decision is genuinely settled.
- **Verification exercises the real workflow.** Completion criteria should map to a real user journey through the external interface, not just "unit tests pass".

## What `/implement` reads

`plan.mdx` is a **narrative** — read by humans and agents for orientation, **not parsed for work**. The only machine-read part is its **frontmatter** (`title`, `status`, `threshold`, `max_iterations`, `version`), which `/implement` uses for the loop's config.

The operational spine lives in the **tracker**, not in `plan.mdx`. `specify-plan` creates, via `task.js`:
- a root **phase** task,
- one **slice** task per deliverable — its body is the code brief (Features / Behavior — EARS / Quality Gates), with `--depends-on` edges to other slices,
- **scenario** child tasks under each slice — the holdout acceptance evidence.

`/implement` reads that task tree (`task.js graph`), and status is grep-able there. Keep `plan.mdx` a clean narrative — **Target State**, **Phases & Outcomes** (one outcome line per slice + its dependencies), and a **Tracking** link to the tracker. Put per-slice detail in the tasks, not the document; a collapsible (`<details>`) may still tuck away verbose scaffolding, but the narrative headings stay at top level.
