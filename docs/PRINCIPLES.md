# Principles

**The Agentic Startup — design principles for Claude Code skills, subagents, and plugins.**

This document is grounded in the primary sources listed under [§ Sources](#sources), the majority of which are Anthropic official documentation from late 2025 and 2026 covering Agent Skills, subagents, Agent Teams, and model selection. Where a claim is inherited from general software-engineering theory rather than a Claude-Code-specific source, that is stated explicitly.

Older drafts of this repo's principles (pre-2026) leaned on CrewAI / AutoGen / LangGraph research, framework-detection agent templates, and multiplier metrics ("10× planning", "3× fewer bugs") that do not survive citation-checking against current Anthropic guidance. Those claims have been dropped rather than relabeled.

---

## 1. Positioning

The Agentic Startup is a Claude Code marketplace project that ships:

- **Workflow skills** — spec-driven development, analysis, review, implementation
- **Specialized subagents** — activity-scoped delegates for research, design, implementation, review
- **Plugins** — bundled distributions composing the above

Every design decision in this repo maps to one of three mechanisms exposed by Claude Code in 2026: **Skills**, **Subagents**, or **Agent Teams**. Each mechanism has a distinct runtime contract, and this document is organized around those contracts rather than around analogies to human team structures.

---

## 2. Core Design Principles

These seven principles apply to both skills and subagents and are the load-bearing rules for the rest of this document.

### 2.1 Description is the activation contract

Both skills and subagents are selected by Claude performing text reasoning against their `description` frontmatter — not embedding retrieval, not keyword matching, not a classifier ([Anthropic Engineering, Equipping agents for the real world with Agent Skills, Dec 2025][eng-skills]; [Lee Han Chung, Claude Agent Skills Deep Dive, Oct 2025][leehan]).

Consequences:

- **Every skill and subagent must front-load its trigger scenario in the first ~50 characters of `description`.** The Claude Code `/skills` UI truncates at 250 characters; the field is hard-capped at 1,024 characters; combined `description + when_to_use` is capped at 1,536 characters for context efficiency ([Skill authoring best practices][skill-best], [Create custom subagents][subagents]).
- **Descriptions are third-person and scenario-anchored.** "Extracts text and tables from PDF files. Use when the user mentions PDFs, forms, or document extraction." — not "Helps with documents."
- **Expect under-triggering by default.** Anthropic explicitly notes that Claude tends to not invoke skills when it probably should, and recommends "slightly pushy" phrasing with `Use PROACTIVELY when…` and `MUST BE USED when…` patterns ([Skill authoring best practices][skill-best]).

### 2.2 Progressive disclosure is the enforced context-economy pattern

Skills load in three tiers, not as one monolithic prompt ([Agent Skills overview][skills-overview]):

| Tier | Content | When loaded | Cost |
|------|---------|-------------|------|
| L1 Metadata | `name` + `description` | Always, at session start | ~100 tokens per skill |
| L2 Body | Full `SKILL.md` | When model self-triggers the skill | ≤5,000 tokens budgeted |
| L3 Resources | `reference/`, `scripts/`, `templates/` | Read on demand via filesystem | Unbounded if unused |

The load boundary is enforced by the runtime — Claude reads L3 files with ordinary Read/Bash tool calls, so anything you bundle has zero context cost until it is explicitly referenced.

**Implications:**

- Keep `SKILL.md` body **≤500 lines** ([Skill authoring best practices][skill-best]). Split longer content into `reference/` files.
- Keep references **one level deep** from `SKILL.md`. Nested references (SKILL.md → advanced.md → details.md) cause Claude to fall back on `head -N` preview reads, yielding incomplete information ([Skill authoring best practices][skill-best]).
- Long reference files (>100 lines) should lead with a table of contents so preview reads still surface scope.

### 2.3 Subagent context isolation is the feature, not a side effect

A subagent receives only: its own system prompt, the Agent-tool delegated task prompt, project `CLAUDE.md`, and its allowlisted tools. It does **not** receive the parent's conversation history, tool results, or system prompt ([Create custom subagents][subagents]).

The only parent→child channel is the Agent-tool prompt string. Anything the subagent needs — file paths, error messages, prior decisions, constraints — must be passed through that string.

**Implications:**

- Design the delegation prompt like a brief for a colleague who just walked into the room. Include targets, decisions, and constraints inline.
- Do not write subagents that assume they can "see what we discussed." They cannot.
- Use isolation deliberately: spawn a subagent when you want the exploration output to stay out of the parent context; invoke a skill inline when you want the result to remain visible.

### 2.4 Activity specialization beats role mega-agents

Anthropic's own agent library and the community consensus ([VoltAgent awesome-claude-code-subagents][voltagent]; [Subagents — Work with subagents][subagents]) favor many small activity-scoped agents (`review-security`, `design-system`, `optimize-performance`) over single role-agents that do everything.

Rationale:

- Focused system prompts raise accuracy; bloated prompts lower it.
- Activity agents are parallel-dispatchable when independent — wall-clock time drops roughly proportional to parallelism on independent work.
- Smaller tool allowlists are easier to audit and produce fewer unnecessary permission prompts.
- Failure in one activity doesn't halt the others.

Single-Responsibility is an inherited principle from general software engineering; its specific application to Claude Code subagents is empirically supported by the structure of Anthropic's shipped examples and the context-isolation mechanism in § 2.3.

### 2.5 Least-privilege tool scoping is enforced at dispatch

For subagents, `tools` in frontmatter is a whitelist applied before the first turn — tools not in the list are stripped from the catalog at dispatch time ([Create custom subagents — Available tools][subagents]). For skills, `allowed-tools` in frontmatter is a pre-approval list, not a restriction — other tools remain callable subject to normal permission flow ([Claude Code skills docs][skills-doc]).

**Rules:**

- Default to **`Read`, `Grep`, `Glob`** for research / analysis agents.
- Add `Edit`, `Write` only for implementer agents.
- Add `Bash` sparingly; when needed, pair with a `PreToolUse` hook that validates commands (e.g., SELECT-only for a DB query agent) ([Subagents — Conditional rules with hooks][subagents]).
- Never use `tools: *`. Never use `tools: inherit` unless explicitly justified — subagents lose the parent's approval history, so inherited dangerous tools re-prompt on every call.
- Restrict sub-subagent spawning with `tools: Agent(name1, name2)` syntax.

### 2.6 Model selection is a tactical cost and latency lever

April 2026 pricing and capability (per [Models overview][models]):

| Model | Input $/M | Output $/M | SWE-bench |
|-------|-----------|------------|-----------|
| Haiku | 0.80 | 4.00 | ~73% |
| Sonnet | 3.00 | 15.00 | ~80% |
| Opus | 15.00 | 75.00 | ~89% |

**Guidance ([Model configuration][model-config]):**

- **Haiku** for high-volume read-heavy work — codebase search, file discovery, pattern matching. Anthropic's built-in `Explore` agent uses Haiku for exactly this reason.
- **Sonnet** default for general coding and implementation.
- **Opus** for complex reasoning — architectural review, security analysis, hard refactors.
- Omit the `model` field to inherit the parent's model when no tactical reason to override.

### 2.7 Evaluation-first authoring

For skills, Anthropic's March 2026 `skill-creator` tool ships benchmarking, blind A/B comparison, and "outgrowth" detection (when the base model has caught up and the skill is redundant) ([Improving skill-creator, Mar 2026][skill-creator]).

Recommended authoring flow:

1. Run Claude on representative tasks **without** the skill/agent; record specific failures.
2. Write 3+ eval scenarios capturing those failures with input/expected-behavior pairs.
3. Write the minimal skill/agent body needed to pass the evals.
4. Benchmark against the unprompted baseline to confirm improvement.
5. Re-run periodically — if the base model has absorbed the behavior, retire the skill.

Subagents currently have no official eval framework; transcript inspection at `~/.claude/projects/{project}/{session}/subagents/agent-{id}.jsonl` and manual golden-path checks are the 2026 state of the art ([Subagents — Detecting subagent invocation][sdk-subagents]).

---

## 3. Skills

### 3.1 Canonical structure

```
plugins/<plugin>/skills/<skill-name>/
├── SKILL.md                 # frontmatter + body, ≤500 lines
├── reference/               # one level deep from SKILL.md
│   └── <topic>.md           # >100 lines → lead with TOC
├── templates/               # static boilerplate
├── examples/                # 3+ I/O pairs suitable for evals
└── scripts/                 # executable helpers (security-reviewed)
```

Skills are discovered from these scopes, highest precedence first ([Claude Code skills docs][skills-doc]):

1. **Enterprise** — managed settings, org-wide
2. **Personal** — `~/.claude/skills/<skill>/SKILL.md`
3. **Project** — `.claude/skills/<skill>/SKILL.md`
4. **Plugin** — `<plugin>/skills/<skill>/SKILL.md`, namespaced as `<plugin>:<skill>`

### 3.2 Frontmatter discipline

| Field | Required | Notes |
|-------|----------|-------|
| `name` | Yes | lowercase-kebab, ≤64 chars; gerund form preferred (`reviewing-prs`) |
| `description` | Yes (effectively) | Front-load trigger in first 50 chars. Max 1,024 chars; truncated at 250 in the `/skills` UI |
| `when_to_use` | No | Combined with `description` capped at 1,536 chars |
| `disable-model-invocation` | No | `true` for side-effect skills (`/deploy`, `/commit`) — **see § 3.4 known bug** |
| `user-invocable` | No | `false` to hide from `/` menu while keeping auto-trigger (background-knowledge skills) |
| `context` | No | `fork` to run in isolated subagent context; requires `agent:` field |
| `allowed-tools` | No | Pre-approval list — **not** a restriction |
| `model` | No | Override session model while skill is active |
| `paths` | No | Glob patterns; activates skill only when editing matching files |

Sources: [Claude Code skills docs][skills-doc], [Skill authoring best practices][skill-best].

### 3.3 Writing conventions

- **Concise body.** Assume Claude is competent; only include context Claude doesn't have.
- **Consistent terminology.** Pick one term per concept and use it throughout.
- **Match degrees of freedom to task fragility.** Strict sequences → exact scripts; judgment calls → text guidelines; middle ground → pseudocode with parameters.
- **Avoid time-sensitive phrasing** ("until 2025, then…"). Use explicit "old pattern" sections if legacy information must survive.
- **Dynamic context injection** via `` !`<command>` `` is available for injecting live repo state into the prompt at load time ([Claude Code skills — Advanced patterns][skills-doc]).

### 3.4 Known anti-patterns and bugs

- **Oversized SKILL.md** (>500 lines) — split into references.
- **Vague descriptions** — never trigger; Claude can't infer intent from "Helps with documents".
- **Buried keywords** — truncation at 250 chars means trigger phrases in sentence 3 are lost.
- **Nested references** — Claude uses `head -N` preview reads on multi-level chains; information is incomplete.
- **Hardcoded absolute paths** in scripts — use `${CLAUDE_SKILL_DIR}`.
- **Undocumented tool dependencies** — state required packages and installation commands.
- **Bug #26251** — `disable-model-invocation: true` currently also blocks *user* invocation, contrary to intended behavior ([anthropics/claude-code issue tracker][gh-issues]). Skills relying on this pattern are partially broken for explicit invocation. Plan accordingly.
- **#22345** — plugin-delivered skills do not currently respect `disable-model-invocation`. Gate side-effect operations some other way in plugins.

---

## 4. Subagents

### 4.1 Canonical structure

```
plugins/<plugin>/agents/<role>/<activity>.md
```

Activity-per-file, not role-per-file. Naming: lowercase-kebab. The agent name becomes the `subagent_type` parameter at dispatch.

Discovery precedence, highest first ([Subagents — File and directory precedence][subagents]):

1. Managed settings (organization-wide)
2. `--agents` CLI flag (session-only JSON)
3. `.claude/agents/` (project, checked into git)
4. `~/.claude/agents/` (user, all projects)
5. Plugin `agents/` directory (read-only, enabled plugins)

### 4.2 Frontmatter discipline

| Field | Required | Notes |
|-------|----------|-------|
| `name` | Yes | Matches filename stem, lowercase-kebab |
| `description` | Yes | Routing contract — when and why Claude delegates. Include `Use PROACTIVELY when…` and 2–3 `<example>` blocks for reliable triggering |
| `tools` | No | Whitelist; enforced at dispatch. Default to `Read, Grep, Glob` for research agents |
| `disallowedTools` | No | Denylist applied to inherited set; can coexist with `tools` |
| `model` | No | `haiku` \| `sonnet` \| `opus` \| specific ID. Omit to inherit parent |
| `permissionMode` | No | `default` \| `acceptEdits` \| `auto` \| `plan` \| etc. Parent mode overrides if stricter |
| `maxTurns` | No | Hard cap on agentic turns |
| `skills` | No | Preloads listed skills into the subagent at startup (full body, not just availability) |
| `mcpServers` | No | Name reference (share parent connection) or inline definition (scoped to subagent) |
| `hooks` | No | Fires only while subagent is active |
| `memory` | No | `user` \| `project` \| `local` for persistent cross-session memory |
| `background` | No | `true` for concurrent execution; auto-denies un-pre-approved tools |
| `isolation` | No | `worktree` runs agent in a fresh git worktree; auto-cleaned on no-op |
| `initialPrompt` | No | Auto-submitted first turn when running as main session via `--agent` |

Source: [Create custom subagents][subagents], [Subagents SDK docs][sdk-subagents].

### 4.3 Writing conventions

- **Open with role + purpose.** One or two sentences. "You are a senior code reviewer ensuring high standards of code quality and security."
- **Explicit numbered workflow** beats open-ended directions. Subagents have no conversation history to infer from.
- **State negative constraints explicitly.** Read-only agents should declare "You cannot modify data. If asked, explain the limitation."
- **Keep body ≤25 KB** for context efficiency; defer bulk reference to `skills:` preload or external files.
- **Include 2–3 `<example>` blocks** in the `description` showing concrete invocation scenarios — this improves parent-side delegation accuracy.

### 4.4 Delegation patterns

Three invocation patterns ([Subagents — Work with subagents][subagents]):

1. **Automatic delegation** — description-matched, parent chooses.
2. **Explicit @-mention** — `@agent-name do X` forces a specific agent.
3. **Session-wide** — `claude --agent name` replaces the session's system prompt with the agent's, for a full interactive session.

**Parallel vs sequential:**

- Independent research / review tasks → spawn N agents in a single response; they run concurrently.
- Dependent phases → spawn sequentially; include earlier results in the next prompt.
- Sustained peer coordination → escalate to Agent Teams (§ 5.2).

**Resuming a subagent** requires Agent Teams (`CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1`) plus `SendMessage` with the agent's ID. Otherwise, each invocation is stateless.

### 4.5 Anti-patterns

- **`tools: *`** — grants Bash, Write, everything; loses parent approvals; every call re-prompts.
- **Mega-role god-agents** — single agent with 50+ KB prompt and ten responsibilities; low accuracy, un-parallelizable, hard to audit.
- **Assuming parent context is visible** — it isn't. Include what you need in the task prompt.
- **Parallel agents with implicit inter-dependencies** — siblings don't see each other; they return to the parent independently.
- **Headless write-capable agents with no review gate** — run plan mode first, or require human review between read and write phases.
- **Context bloat from verbose subagent returns** — ask for key findings, not exhaustive details.
- **Permission-mode confusion** — parent mode overrides subagent; stricter-than-parent is honored, looser-than-parent is not.

---

## 5. Composition

### 5.1 Skill vs Subagent vs Agent Team vs Hook

| Mechanism | When to use | Context | Communication |
|-----------|-------------|---------|---------------|
| **Skill** | Reusable workflow, domain knowledge, checklist, slash command | Parent's own context (unless `context: fork`) | Inline — content becomes conversation |
| **Subagent** | Isolated exploration, parallel specialized work, summary-returning research | Fresh fork; parent unseen | One-way: parent → task prompt → summary |
| **Agent Team** | Sustained peer coordination, competing hypotheses, cross-layer debugging | Persistent per-teammate contexts | Bidirectional mailbox; peer-to-peer |
| **Hook** | Deterministic gate on a specific tool call (e.g., block writes, validate queries) | Runs as shell process | Exit code / stdout / stderr |

Sources: [Claude Code skills docs][skills-doc], [Create custom subagents][subagents], [Agent Teams][agent-teams].

2026 trend note: slash commands and skills have been unified — `.claude/commands/` and `.claude/skills/` both work, and skill frontmatter (`argument-hint`, `disable-model-invocation`, `user-invocable`) subsumes what used to require separate command definitions ([Claude Code skills docs][skills-doc]).

### 5.2 Agent Teams

Agent Teams are experimental as of April 2026 and gated behind `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1` ([Agent Teams][agent-teams]).

- Each teammate runs in its own session with its own context window.
- Teammates have mailboxes and can message each other directly via `SendMessage`.
- A shared task list is accessible to all teammates.
- Appropriate for: cross-domain research where findings challenge each other; multi-phase work with coordinated hand-offs; long-running investigations.
- **Not** appropriate for: simple delegate-and-summarize work — regular subagents are lighter weight.

### 5.3 Information barriers

The value of isolation depends on it being respected by design:

- A code-writing agent should **not** see the evaluation scenarios it's being graded against.
- An evaluation agent should **not** see the source code it's grading.
- The orchestrator owns any service lifecycle (start / health-check / restart / stop).
- Filtered one-line failure summaries cross barriers; raw transcripts do not.

These patterns come from the repo's own prior experience with factory-style implementation loops and are consistent with context-isolation guidance in the Claude Code docs.

---

## 6. Quality and Evaluation

### 6.1 Skills

Use `skill-creator` ([Improving skill-creator, Mar 2026][skill-creator]) for:

- **Evals** — define input prompts + "what good looks like"; run pass/fail.
- **Benchmarks** — pass rate, elapsed time, token usage over full suite.
- **Blind A/B** — comparator agent doesn't know which variant produced which output.
- **Outgrowth detection** — flags when the base model passes evals without the skill, indicating the skill's behavior has been absorbed.

### 6.2 Subagents

No official eval framework in 2026. Practical techniques ([Subagents SDK — Detecting subagent invocation][sdk-subagents]):

- **Trigger testing** — invoke with natural language matching the description; verify delegation.
- **Transcript inspection** — read `~/.claude/projects/{project}/{session}/subagents/agent-{id}.jsonl`.
- **Tool enforcement** — attempt a non-allowlisted tool; verify clean failure, not a permission prompt.
- **Cost tracking** — verify custom-model agents are actually running on the declared model.

---

## 7. Known Gaps and Open Questions

Current limitations worth tracking:

- **Bug #26251** — `disable-model-invocation: true` blocks user invocation too.
- **#22345** — plugin skills miss `disable-model-invocation` entirely.
- **Cross-surface skill portability** — skills uploaded to claude.ai, the API, and Claude Code do not sync.
- **Agent Teams maturity** — experimental; real-world patterns still emerging.
- **Subagent nesting** — officially unsupported; `context: fork` inside a subagent-loaded skill is an undocumented workaround.
- **Outgrowth detection mechanism** — Anthropic claims the signal but doesn't document the threshold.
- **Subagent eval tooling** — no official framework; each team reinvents.
- **Size-limit empiricism** — Anthropic publishes ≤500 lines for SKILL.md; community reports range up to ~5,000 words with good progressive disclosure. Lacking public benchmarks.

---

## 8. Validation Checklist

Apply this checklist to every skill, subagent, and plugin change in this repo.

**Skills:**

- [ ] `SKILL.md` ≤ 500 lines; references one level deep
- [ ] `description` front-loads trigger in first 50 chars; third-person; ≤250 chars safe
- [ ] 3+ eval scenarios exist in `examples/`
- [ ] `allowed-tools` lists only what's needed (or is omitted)
- [ ] `disable-model-invocation` used on side-effect skills only, with bug #26251 noted
- [ ] No time-sensitive phrasing, no absolute paths

**Subagents:**

- [ ] Activity-scoped, not role-scoped
- [ ] `description` contains trigger scenarios and 2–3 examples
- [ ] `tools` is an explicit least-privilege list (never `*`)
- [ ] `model` is tactically chosen (Haiku for read-only research; Sonnet default; Opus only for hard reasoning)
- [ ] Body ≤ 25 KB; bulk reference deferred to `skills:` or external files
- [ ] Does not assume visibility of parent context
- [ ] Golden-path trigger test performed; transcript reviewed

**Composition:**

- [ ] Each behavior placed in the right mechanism (skill vs subagent vs team vs hook)
- [ ] Independent work is parallelized; only dependent phases are sequential
- [ ] Information barriers are intact where relied upon (evaluator ≠ code-writer)

---

## Sources

Primary sources used to ground this document. Preferential weight given to Anthropic official docs from 2025–2026.

### Anthropic Official

- [Extend Claude with skills — Claude Code docs][skills-doc] — canonical skill structure, discovery, frontmatter, invocation control, advanced patterns
- [Agent Skills overview — Platform / API docs][skills-overview] — architecture, progressive disclosure, limitations, cross-surface constraints
- [Skill authoring best practices — Platform / API docs][skill-best] — naming, description writing, content guidelines, anti-patterns, evaluation-first workflow
- [Equipping agents for the real world with Agent Skills — Engineering Blog, Dec 2025][eng-skills] — design rationale for progressive disclosure
- [Improving skill-creator: Test, measure, and refine Agent Skills — Blog, Mar 2026][skill-creator] — evals, benchmarking, A/B, outgrowth
- [Create custom subagents — Claude Code docs][subagents] — subagent frontmatter schema, tool gating, examples
- [Subagents in the SDK — Claude API / Agent SDK docs][sdk-subagents] — programmatic definition, tool restriction enforcement, invocation detection
- [Orchestrate teams of Claude Code sessions — Claude Code docs][agent-teams] — Agent Teams experimental feature
- [Best practices for Claude Code — docs][best-practices] — general anti-patterns
- [Model configuration — Claude Code docs][model-config] — per-session / per-subagent model selection
- [Models overview — API docs][models] — April 2026 pricing and capability table
- [anthropics/claude-code issue tracker][gh-issues] — bugs #19141, #22345, #26251, #30355, #40121

### Community and Third-Party

- [Lee Han Chung, Claude Agent Skills: Mechanisms and Patterns, Oct 2025][leehan] — deep-dive on dual-message injection, meta-tool architecture, activation
- [Dean Blank, A Mental Model for Claude Code: Skills, Subagents, Plugins, Mar 2026][mentalmodel] — decision-boundary synthesis
- [Steve Kinney, Common Sub-Agent Anti-Patterns and Pitfalls][kinney] — curated anti-patterns
- [VoltAgent, awesome-claude-code-subagents][voltagent] — 100+ activity-scoped agent examples
- [cc-plugin-eval][eval-framework] — community evaluation framework for plugins

[skills-doc]: https://code.claude.com/docs/en/skills
[skills-overview]: https://platform.claude.com/docs/en/agents-and-tools/agent-skills/overview
[skill-best]: https://platform.claude.com/docs/en/agents-and-tools/agent-skills/best-practices
[eng-skills]: https://www.anthropic.com/engineering/equipping-agents-for-the-real-world-with-agent-skills
[skill-creator]: https://claude.com/blog/improving-skill-creator-test-measure-and-refine-agent-skills
[subagents]: https://code.claude.com/docs/en/sub-agents
[sdk-subagents]: https://code.claude.com/docs/en/agent-sdk/subagents
[agent-teams]: https://code.claude.com/docs/en/agent-teams
[best-practices]: https://code.claude.com/docs/en/best-practices
[model-config]: https://code.claude.com/docs/en/model-config
[models]: https://platform.claude.com/docs/en/about-claude/models/overview
[gh-issues]: https://github.com/anthropics/claude-code/issues
[leehan]: https://leehanchung.github.io/blogs/2025/10/26/claude-skills-deep-dive/
[mentalmodel]: https://levelup.gitconnected.com/a-mental-model-for-claude-code-skills-subagents-and-plugins-3dea9924bf05
[kinney]: https://stevekinney.com/courses/ai-development/subagent-anti-patterns
[voltagent]: https://github.com/VoltAgent/awesome-claude-code-subagents
[eval-framework]: https://github.com/sjnims/cc-plugin-eval

---

*Document grounded in primary sources listed above. When these sources are updated or contradicted by newer Anthropic guidance, update this document accordingly rather than preserving stale claims.*
