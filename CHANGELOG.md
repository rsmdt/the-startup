# Changelog

All notable changes to The Agentic Startup are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Entries are added only when a release is cut. Work in progress is not tracked here — see the git history between the latest tag and `HEAD` for unreleased changes.

This file retains detailed entries for the last 10 minor releases plus their patch revisions. Older release notes are linked to their [GitHub release page](https://github.com/rsmdt/the-startup/releases).

## [3.6.1] - 2026-04-22

### Changed
- Rewrote core principles against 2026 Claude Code primary sources, consolidating PHILOSOPHY.md into PRINCIPLES.md and removing pre-2026 framework-era guidance that no longer held up to citation checking.
- Added a skill-vs-subagent decision framework to the principles doc to help users choose the right primitive when extending their workflows.
- Restructured CHANGELOG.md with a retained history window and links out to older releases for easier navigation.

## [3.6.0] - 2026-04-17

### Changed

- `analyze` skill now carries a trigger-rich description, per-perspective agent mapping, and auto-persists findings under `docs/` without a confirmation round-trip. Constraints are softened with inline reasoning so edge cases can be judged rather than blindly followed.
- Removes outdated "via Task tool" phrasing from `analyze`, `refactor`, `validate`, `constitution`, `document`, and `writing-skills` — subagent delegation now uses the project's current convention throughout.

## [3.5.0] - 2026-04-02

### Added

- **Dark Factory workflow** — `specify-factory` skill decomposes a specification into units, scenarios, and a manifest that the factory loop consumes. Replaces the sequential PRD → SDD → PLAN pipeline.
- **TDD enforcement across the factory loop.** Code agents follow mandatory red-green-refactor per requirement; `specify-factory` emits executable E2E test stubs alongside scenarios; evaluation agents prefer those stubs over ad-hoc test writing.
- **Eval coverage** for `specify`, `specify-requirements`, `specify-solution`, `specify-factory` (4 scenarios), and `implement` (6 scenarios: happy path, retry, max-iterations escalation, parallel execution, resume-from-partial, E2E stubs integration).

### Changed

- **BREAKING**: `implement` is now a factory loop orchestrator with information-barrier subagents — code agents never see scenarios; evaluation agents never see source code. The retry loop is governed by filtered one-line failure summaries until the satisfaction threshold is met. Consumers of the old phase-by-phase `implement` flow must migrate to the new manifest-driven loop.

### Removed

- **BREAKING**: `specify-plan` skill removed. Its role is replaced by `specify-factory`, which produces a factory decomposition manifest instead of a phased implementation plan. Specs created with the old pipeline remain readable but must be re-decomposed to run through the new `implement`.

## [3.4.0] - 2026-02-28

### Changed

- `analyze` and `specify-solution` now require **mechanism-level depth** and lead with the clean solution. Findings must answer What / How / Why with file:line evidence; recommendations present the architecturally clean approach first, with hybrids and trade-downs gated behind an explicit user request.
- `test`, `specify-requirements`, and `specify-solution` enforce **MECE coverage** (mutually exclusive, collectively exhaustive) so sections cannot overlap or leave gaps.
- Completes the SudoLang → Markdown migration across remaining skills, agents, and output styles; `writing-skills` conventions updated to match.

## [3.3.0] - 2026-02-22

### Added

- **`brainstorm` skill** — probes ideas through structured dialogue before implementation. Explores approaches with trade-offs, builds designs iteratively, and gates on user approval before any code is written.
- **`agentic-patterns` and `frontend-patterns` skills** — domain-specific reference materials with llms.txt pointers.
- **Per-phase plan files under `.start/specs/`**. Monolithic `implementation-plan.md` is replaced with a `plan/` directory containing a `README.md` manifest and individual `phase-N.md` files. `implement` loops through phases with resumability, updating frontmatter status and manifest checkboxes on completion. Dual-path resolution preserves backward-compatible access to `docs/specs/`.

### Changed

- **BREAKING**: Team plugin consolidated from 22 agents to 15 and from 20 skills to 14. Reduces overlap and context bloat while preserving specialist depth; high-value reference material migrates to the consolidated skills. Users with pinned references to removed agents or skills must update to their consolidated replacements.
- **BREAKING**: All 35 skills migrated from SudoLang to Markdown conventions. SudoLang syntax (`Constraints {}`, `fn` declarations, `|>` pipes) replaced with Markdown equivalents (Always/Never lists, `###` headings, numbered steps). Duplicate constraint rules deduplicated, Entry Point sections restricted to non-linear workflows, types converted to valid TypeScript, and mode naming standardized to "Agent Team" across all skills.
- Spec file names shortened: `product-requirements-document.md` → `requirements.md`, `solution-design-document.md` → `solution.md`.

## [3.2.1] - 2026-02-20

### Changed

- Skills now use proper markdown links to load supplementary reference files, replacing ad-hoc inline path references. Ensures progressive-disclosure content resolves consistently across contexts.

## [3.2.0] - 2026-02-20

### Changed

- All 15 start skills, 22 team agents, 22 team skills, 2 output styles, and 3 template files restructured with **AIDD/SudoLang PICS layout** (Preamble → Interface → Components → Start). Adds typed Output Schemas, SudoLang Constraints blocks, and decision tables with first-match-wins semantics so skill and agent behavior becomes explicit and checkable. *(Superseded by v3.3.0's migration back to Markdown conventions — see that release's changelog.)*

## [3.1.0] - 2026-02-19

### Changed

- Solution designs now require **schema accuracy and traced walkthroughs**. SQL examples must reference migration-verified column names instead of pseudocode, and complex query logic requires a traced walkthrough with example data. Closes a recurring class of wrong-column bugs that passed review but failed in implementation.

## [3.0.0] - 2026-02-17

### Added

- Installer prompts for **Agent Teams** installation and documents v3 features in the README.

### Changed

- **BREAKING**: Skill invocation drops the `/start:` namespace prefix in autocomplete. Skills now invoke as `/specify`, `/implement`, `/analyze`, etc. Documentation and cross-references replace the term "command" with "skill" throughout. Consumers referencing `/start:<skill>` in external documentation should update to the unprefixed form.

## [2.18.0] - 2026-02-13

### Added

- **`test` skill** — enforces code-ownership standards (you touch it, you own it) and automatic test discovery, ensuring MECE coverage across categories including E2E.
- **`writing-skills` skill** — for authoring high-quality skills, with reference materials on persuasion principles and subagent testing patterns.
- Expanded core skills with comprehensive workflows and protocols.

### Changed

- **BREAKING**: User-invocable commands (10 workflow commands) migrated to `SKILL.md` files within dedicated skill directories. Aligns user-invocable skills with the same structure used by autonomous skills, enabling consistent plugin architecture.
- **BREAKING**: Skill surface area consolidated. `drift-detection` and `specification-validation` merged into `validate`; `code-review` merged into `review`; `constitution-validation` merged into `constitution`; `error-recovery` merged into `coding-conventions` (team agents updated to drop the dangling reference).
- **BREAKING**: Four skills renamed to a `specify-*` hierarchy — `requirements-analysis` → `specify-requirements`, `architecture-design` → `specify-solution`, `implementation-planning` → `specify-plan`, `specification-management` → `specify-meta`.

### Removed

- **BREAKING**: Seven skills removed by absorption into parent skills. `codebase-analysis` → `analyze`; `safe-refactoring` and `simplify` → `refactor`; `bug-diagnosis` → `debug`; `agent-coordination` → `implement`; `task-delegation` preserved as a `docs/patterns/` reference; `git-workflow` inlined into `specify` and `implement`.
- **BREAKING**: `implementation-verification` skill — redundant with `specification-validation` (now part of `validate`).

## [2.17.0] - 2026-01-29

### Changed

- Standardizes command output formats with **structured tables** across `review`, `refactor`, `simplify`, and `validate`. Replaces emoji-heavy freeform output with consistent table-based formats, improving readability and enabling findings to be parsed programmatically.

## [2.16.0] - 2026-01-29

### Added

- **Orchestrator-only delegation pattern** for the `implement` command. All tasks — parallel and sequential — now delegate to subagents for context offloading. Subagents self-prime by reading spec documents directly instead of receiving injected context. Ships with a structured result format for concise summaries.

### Changed

- Simplifies Skill call syntax in commands from `Skill(skill: "start:...")` to `Skill(start:...)` for cleaner, more concise invocations.

## Older release history

- [2.15.0]: https://github.com/rsmdt/the-startup/releases/tag/v2.15.0
- [2.14.0]: https://github.com/rsmdt/the-startup/releases/tag/v2.14.0
- [2.13.0]: https://github.com/rsmdt/the-startup/releases/tag/v2.13.0
- [2.12.0]: https://github.com/rsmdt/the-startup/releases/tag/v2.12.0
- [2.11.0]: https://github.com/rsmdt/the-startup/releases/tag/v2.11.0
- [2.10.0]: https://github.com/rsmdt/the-startup/releases/tag/v2.10.0
- [2.9.0]: https://github.com/rsmdt/the-startup/releases/tag/v2.9.0
- [2.8.1]: https://github.com/rsmdt/the-startup/releases/tag/v2.8.1
- [2.8.0]: https://github.com/rsmdt/the-startup/releases/tag/v2.8.0
- [2.7.0]: https://github.com/rsmdt/the-startup/releases/tag/v2.7.0
- [2.6.2]: https://github.com/rsmdt/the-startup/releases/tag/v2.6.2
- [2.6.0]: https://github.com/rsmdt/the-startup/releases/tag/v2.6.0
- [2.5.1]: https://github.com/rsmdt/the-startup/releases/tag/v2.5.1
- [2.5.0]: https://github.com/rsmdt/the-startup/releases/tag/v2.5.0
- [2.4.1]: https://github.com/rsmdt/the-startup/releases/tag/v2.4.1
- [2.4.0]: https://github.com/rsmdt/the-startup/releases/tag/v2.4.0
- [2.3.0]: https://github.com/rsmdt/the-startup/releases/tag/v2.3.0
- [2.2.0]: https://github.com/rsmdt/the-startup/releases/tag/v2.2.0
- [2.1.0]: https://github.com/rsmdt/the-startup/releases/tag/v2.1.0
- [2.0.0]: https://github.com/rsmdt/the-startup/releases/tag/v2.0.0
- [1.15.0]: https://github.com/rsmdt/the-startup/releases/tag/v1.15.0
- [1.14.0]: https://github.com/rsmdt/the-startup/releases/tag/v1.14.0
- [1.13.0]: https://github.com/rsmdt/the-startup/releases/tag/1.13.0
- [1.12.0]: https://github.com/rsmdt/the-startup/releases/tag/1.12.0
- [1.11.0]: https://github.com/rsmdt/the-startup/releases/tag/1.11.0
- [1.10.0]: https://github.com/rsmdt/the-startup/releases/tag/v1.10.0
- [1.9.0]: https://github.com/rsmdt/the-startup/releases/tag/v1.9.0
- [1.8.0]: https://github.com/rsmdt/the-startup/releases/tag/v1.8.0
- [1.7.1]: https://github.com/rsmdt/the-startup/releases/tag/v1.7.1
- [1.7.0]: https://github.com/rsmdt/the-startup/releases/tag/v1.7.0
- [1.6.0]: https://github.com/rsmdt/the-startup/releases/tag/v1.6.0
- [1.5.1]: https://github.com/rsmdt/the-startup/releases/tag/v1.5.1
- [1.5.0]: https://github.com/rsmdt/the-startup/releases/tag/v1.5.0
- [1.4.0]: https://github.com/rsmdt/the-startup/releases/tag/v1.4.0
- [1.3.0]: https://github.com/rsmdt/the-startup/releases/tag/v1.3.0
- [1.2.0]: https://github.com/rsmdt/the-startup/releases/tag/v1.2.0
- [1.1.0]: https://github.com/rsmdt/the-startup/releases/tag/v1.1.0
- [1.0.0]: https://github.com/rsmdt/the-startup/releases/tag/v1.0.0
