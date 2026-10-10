# Harness Surfaces

Each surface a proposal may target: what it is, where its fix lives, and what to read before drafting.
Fill `{projectRoot}` and `{pluginsRoot}`. The menu line (label: essence) goes into the mine and cluster
prompts; the full brief goes into the propose prompt.

## instruction — CLAUDE.md / AGENTS.md
- **Essence:** improve or simplify the baseline instruction — remove dead or contradicted rules, sharpen
  ambiguous ones.
- **Fix lives in:** `{projectRoot}/CLAUDE.md` or `{projectRoot}/AGENTS.md`, injected into every session.
- **Read first:** the target file. A good edit removes or sharpens a rule; it does not just append.

## skill
- **Essence:** improve an existing skill, or create a new skill for recurring manual work.
- **Fix lives in:** a `SKILL.md` (plus `reference/`, `templates/`) under `{pluginsRoot}/*/skills/`, or a
  new skill directory there.
- **Read first:** the target skill in full, and two neighbouring skills in the same plugin for its
  conventions: YAML frontmatter (`name`, `description`, `user-invocable`, `argument-hint`), a Persona,
  Always/Never constraints, a Reference Materials list, and numbered workflow steps; detail that is not
  needed on every run lives in `reference/` and is read on demand. A skill that orchestrates several
  agents describes the process and has Claude run it as a dynamic workflow (see `skills/implement`).

## agent — subagent
- **Essence:** improve an existing subagent, or create a new one for recurring specialized work worth
  delegating.
- **Fix lives in:** `{pluginsRoot}/team/agents/the-*/*.md`, or a new agent file there.
- **Read first:** `{pluginsRoot}/team/agents/the-meta-agent.md` and the target agent.

## output-style
- **Essence:** improve an output style the logs show misfiring (verbosity, ceremony, voice).
- **Fix lives in:** `{pluginsRoot}/start/output-styles/*.md`.
- **Read first:** the target style; its effect is visible in the assistant records you cite.

## guardrail — permission rule / hook / rule
- **Essence:** prevent a recurring dangerous or wrong action.
- **Fix lives in**, most deterministic first:
  1. **Permission rule** — a deny/ask entry under `"permissions"` in `{projectRoot}/.claude/settings.json`
     (e.g. deny `Bash(rm -rf:*)`, ask on `Read(.env)`). Best for dangerous or sensitive operations.
  2. **Hook** — an entry under `"hooks"` in the same file running a command on a tool event such as
     `PreToolUse`. Best for context-dependent checks.
  3. **Rule** — an Always/Never line in `{projectRoot}/CLAUDE.md` when neither can match.
- **Read first:** the current `{projectRoot}/.claude/settings.json`. Choose the narrowest matcher that
  stops the evidenced risk without blocking routine work.
