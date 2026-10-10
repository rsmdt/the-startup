# AGENTS.md

Guidance for Claude Code when working in this repository.

## What this is

**The Agentic Startup** is a Claude Code marketplace. The repository holds no application code: its product is the prompt material and orchestration that Claude Code loads — plugins made of skills, subagents, output styles, and the shared scripts skills call. Editing this repo means editing what Claude reads and runs, so every change is a behavior change, and the plugin's own eval and bench tooling is how it is tested.

Details about any individual skill, agent, workflow, or style live next to that artifact (its own `SKILL.md`, agent file, script header, or the plugin README). Keep them there; this file describes only the shape of the repository and the conventions that hold across it.

## Layout

```
plugins/              one directory per marketplace plugin
  start/              the workflow plugin: skills/, output-styles/, scripts/
  team/               the agent library: agents/the-<role>/<activity>.md, skills/
```

Inside a plugin:

- `skills/<name>/SKILL.md` is a skill. `SKILL.md` is always loaded, so it stays small; `reference/`, `templates/`, and `examples/` beside it are loaded on demand.
- A skill that fans work out to many agents describes the process — stages, agent prompts, return schemas, decision rules — and has Claude run it as a dynamic workflow (the Workflow tool with an inline script Claude writes). No workflow scripts ship with a plugin.
- `scripts/` holds CLIs several skills share (`task.js`, `spec.js`). Skills call them as `node ${CLAUDE_PLUGIN_ROOT}/scripts/<name>.js` in SKILL.md, where Claude Code substitutes the installed plugin path.
- `agents/the-<role>/<activity>.md` is a subagent definition; `output-styles/<Title Case>.md` is an output style.

## Working on the repo

Test a plugin by installing it from its local path, and uninstall to reset:

```bash
claude plugin install ./plugins/start
claude plugin uninstall start@the-startup
```

Run a skill's evals with `make eval SKILL=<name>`, or every eligible skill with `make eval-all`. Both run the skill only; add `BASELINE=1` to also run without it and report the difference. Skills that launch a dynamic workflow (`analyze`, `implement`, `review`, `reflect`) are excluded, because the eval runner's executors have no Workflow tool.
