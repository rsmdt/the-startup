// daydream-workflow: v7
//
// The Agentic Startup — daydream workflow (self-improvement from session logs, tracker-backed).
//
// One reusable dynamic-workflow script. /start:daydream installs this to
// <project>/.claude/workflows/daydream.js and launches it via the Workflow tool, passing an
// approved daydream plan as `args`. It is NOT regenerated per run.
//
// It reads the Claude Code session logs over a date range, finds where the agentic system
// underperformed or was exposed to risk, and PRODUCES A REPORT of what should change across the whole
// project harness — every surface that shapes how the agent behaves:
//   - the CLAUDE.md / AGENTS.md baseline instruction,
//   - skills, subagents (agents), workflow golden copies, output-styles,
//   - guardrails: permission policies (allow/deny/ask), hooks, and CLAUDE.md "Always/Never" rules.
// Each report entry is labelled IMPROVE (edit an existing artifact) or CREATE (a new skill, subagent,
// hook, or permission rule), carries a concrete diff / new-artifact body, and cites the log evidence
// that motivates it.
//
// ADVISORY BY DESIGN — this is the identity of the tool, not a setting. daydream NEVER edits or creates
// an artifact itself. It writes proposals to the tracker and a report a human reviews; the human decides
// what to apply. There is no auto-apply tier.
//
// DIMENSION-AWARE. daydream carries a DIMENSION REGISTRY (below): per harness surface, WHAT it is, WHERE
// its fix lives, and a REFERENCE the propose agent reads before drafting (skill-authoring conventions,
// agent design, how Claude Code permissions / hooks are configured, …). The harness system prompt
// already describes skills/agents/CLAUDE.md; the registry injects the rest and points at repo docs that
// exist. This is how the run "injects the meaning of all dimensions".
//
// Pipeline (grounded in the analysis report — self-improvement literature + shipped analogues):
//   MINE     — fan out one agent per signal LENS. Lenses: friction (tool-errors/retries/cost),
//              conversation (semantic inefficiency), router (model sizing), capability (recurring MANUAL
//              work with no owning skill/subagent → new-artifact signal), guardrail (recurring MISTAKE
//              classes AND recurring RISK exposure — dangerous/sensitive ops → permission/hook/rule),
//              instruction (CLAUDE.md/AGENTS.md hygiene: dead, contradicted, or bloated rules). Each walks
//              the session JSONL over the range and writes CITED findings, tagging each with the harness
//              surface a fix would touch. ITERATIVE: sampling lenses that leave sessions un-examined get
//              another pass (bounded by maxMineRounds) — the run "keeps examining within the window".
//   CLUSTER  — ONE agent reads all mined findings, dedups/ranks them into distinct improvement
//              OPPORTUNITIES. For each it makes the IMPROVE-vs-CREATE call explicitly (prefer improving an
//              existing artifact; create new only when nothing fits and the pattern recurs / the risk is
//              severe), names the surface + action, and creates a proposal task via task.js.
//   PROPOSE  — per opportunity, read the surface's reference, then draft the concrete change — a diff
//              (improve) or the full body of a new artifact (create) — + rationale, and classify the tier.
//   VALIDATE — adversarial gate: confirm the cited log evidence ACTUALLY supports the change
//              (presumption of failure). Unsupported proposals are dropped, not deferred.
//   SYNTHESIZE — ONCE, at the end: write the report — ranked, validated proposals with their changes,
//              evidence, IMPROVE/CREATE label, and tier — grouped by surface. A human applies them.
//
// Runtime note: the script has no shell/filesystem access — only its agents do. Every task.js call
// and every file read is made BY an agent; the script only orchestrates and branches on the
// structured verdicts agents return.
//
// Sandbox note: the session-log archive lives OUTSIDE the project root (under the user's Claude
// home, e.g. ~/.claude/projects/<slug>/). The command sandbox denies reads there by default, so a
// mine agent that shells into the archive may hit "Operation not permitted" and must retry with the
// sandbox disabled. Reading via the file-read tool is unaffected. Mine prompts say so explicitly.

export const meta = {
  name: "daydream",
  description: "Self-improvement from session logs: report what should change across the harness (CLAUDE.md, skills, subagents, workflows, output-styles, guardrails incl. permission policies) — each entry labelled improve/create, with a concrete change and cited evidence. Advisory: the human applies.",
  phases: [
    { title: "Mine", detail: "parallel lenses mine the session JSONL (iterative for sampling lenses) → cited findings tagged by harness surface" },
    { title: "Cluster", detail: "dedup/rank findings into opportunities; make the improve-vs-create call per opportunity" },
    { title: "Propose", detail: "read the surface reference, then draft a diff (improve) or new-artifact body (create) + tier" },
    { title: "Validate", detail: "adversarially confirm the log evidence supports each proposal" },
    { title: "Synthesize", detail: "write the proposals report; revise it to clear the grade" },
    { title: "Grade", detail: "grade each proposal on behaviour / improvement / proof; block and revise on failures" },
  ],
};

// ── Input contract (args, supplied by /start:daydream from the approved plan) ──────────────
// {
//   name:            string,          // run slug (e.g. "2026-07-16")
//   home:            string,          // abs path to .start/daydream/<slug> (tracker home)
//   taskScript:      string,          // abs path to specify-meta/task.js
//   description:     string,          // terse purpose
//   dateRange:       { since, until, label },  // ISO dates + a human label (e.g. "yesterday")
//   logRoot:         string,          // abs path to this project's session-log dir (~/.claude/projects/<slug>/)
//   projectRoot:     string,          // abs path to the repo root — anchors CLAUDE.md/AGENTS.md/.claude
//   editTargetsRoot: string,          // abs path to the git-tracked plugin artifacts root (this repo's plugins/)
//   angles:          [{ id, name, kind? }],  // seeded mine task ids. kind selects the LENS (see lens dispatch
//                                             // below): "friction" (default), "conversation", "router",
//                                             // "capability", "guardrail", "instruction".
//   reportPath:      string,          // where SYNTHESIZE writes the report
//   maxProposals:    number,          // cap on proposals carried to the report (default 12)
//   maxMineRounds:   number,          // cap on iterative mine passes (default 2) — sampling lenses re-run
// }

// Background launches may serialize `args`; on some runtimes it is revived as a JSON string.
const parsedArgs = typeof args === "string" ? JSON.parse(args) : args;

const {
  name = "daydream",
  home,
  taskScript,
  description = "",
  dateRange = {},
  logRoot,
  projectRoot,
  editTargetsRoot,
  angles = [],
  reportPath,
  maxProposals = 12,
  maxMineRounds = 2,
} = parsedArgs ?? {};

// FAIL FAST on a mis-wired launch. Every stage prompt templates these values verbatim; if any is
// missing, agents would receive literal "undefined" paths ("node undefined … --home undefined"),
// mine nothing, and synthesize a zero-proposal stub. Abort before any agent runs instead.
{
  const required = { home, taskScript, logRoot, projectRoot, editTargetsRoot, reportPath };
  const missing = Object.entries(required)
    .filter(([, v]) => typeof v !== "string" || !v.trim() || v === "undefined")
    .map(([k]) => k);
  if (!dateRange?.since || !dateRange?.until) missing.push("dateRange.since/until");
  if (!angles.some((a) => a && a.id)) missing.push("angles");
  if (missing.length) {
    throw new Error(
      `daydream: launch args missing or unresolved: ${missing.join(", ")}. ` +
        `The /start:daydream skill must pass concrete absolute paths (see the input contract at the top of this script); ` +
        `refusing to run — a run with unresolved args can only produce an empty report.`,
    );
  }
}

const tj = `node ${taskScript}`; // task.js invocation prefix agents use
const range = `${dateRange.label ?? ""} (${dateRange.since ?? "?"} → ${dateRange.until ?? "?"})`;

// ── Dimension registry — the meaning of every harness surface daydream may improve ─────────
// Injected so the cluster/propose agents know WHAT each surface is, WHERE its fix lives, and the
// REFERENCE to read before drafting. Every surface supports BOTH improve (edit) and, where sensible,
// create (a new artifact) — daydream only reports, so nothing is hard-gated; the cluster makes the
// improve-vs-create call. Paths are anchored to the two roots passed in args; only repo docs that
// actually exist are cited. Facts here are Claude Code specific (the only harness daydream mines).
const DIMENSIONS = {
  instruction: {
    label: "instruction (CLAUDE.md / AGENTS.md)",
    essence: "improve or SIMPLIFY the baseline instruction — remove dead/contradicted rules, sharpen ambiguous ones (rarely a new file)",
    surface: `EDIT the project ${projectRoot}/CLAUDE.md or ${projectRoot}/AGENTS.md — the baseline system instruction injected into every session`,
    reference: `CLAUDE.md/AGENTS.md are treated as system instruction by the harness. Read the target file; a good edit REMOVES a rule the logs show was dead/contradicted or SHARPENS one that was ambiguous — it does not just append.`,
  },
  skill: {
    label: "skill",
    essence: "improve an existing skill, OR create a NEW skill for recurring manual work",
    surface: `IMPROVE an existing SKILL.md under ${editTargetsRoot}, or CREATE a new skill directory (SKILL.md + optional reference/) for recurring manual work no skill owns`,
    reference: `Skill-authoring conventions: ${editTargetsRoot}/start/skills/writing-skills/reference/conventions.md — read before drafting a skill improvement or a new skill.`,
  },
  agent: {
    label: "subagent",
    essence: "improve an existing subagent, OR create a NEW subagent for recurring specialized/isolatable work",
    surface: `IMPROVE an existing agent under ${editTargetsRoot}/team/agents/the-*/*.md, or CREATE a new one for recurring specialized work worth delegating`,
    reference: `Agent design: ${editTargetsRoot}/team/agents/the-meta-agent.md — read before drafting an agent change or a new agent. The harness system prompt also lists the active agents and their triggers.`,
  },
  workflow: {
    label: "workflow",
    essence: "improve a workflow golden copy (models, prompts, control flow); bump its version marker",
    surface: `IMPROVE a workflow GOLDEN COPY under ${editTargetsRoot}/*/skills/*/reference/*.js and BUMP its "// <name>-workflow: vN" version marker (never the installed .claude/workflows/*.js copy)`,
    reference: `The golden copy itself is the reference — read it in full, and cross-check how it ACTUALLY ran in the session logs (its mine/code/eval/validate agents, their models, their retries).`,
  },
  "output-style": {
    label: "output-style",
    essence: "improve an output-style the logs show misfiring (verbosity, ceremony, voice)",
    surface: `IMPROVE an output-style under ${editTargetsRoot}/start/output-styles/*.md`,
    reference: `The active output style is injected into the harness system prompt — read the target .md; the behaviour it dictates is visible in the assistant records you cited.`,
  },
  guardrail: {
    label: "guardrail (permission policy / hook / rule)",
    essence: "prevent a recurring dangerous or wrong action — by a permission deny/ask rule, a hook, or a CLAUDE.md rule",
    surface: `Prevent a recurring dangerous/wrong action before it happens, via the MOST deterministic mechanism that fits (preference order):
  1. PERMISSION POLICY — add a deny/ask entry to "permissions" in ${projectRoot}/.claude/settings.json (e.g. deny "Bash(rm -rf:*)", ask on "Read(.env)" / secret files). Best for dangerous or sensitive operations.
  2. HOOK — add a hook under "hooks" in ${projectRoot}/.claude/settings.json running a command on a matching tool event (e.g. PreToolUse) to block/warn. Best for context-dependent checks a static permission can't express.
  3. RULE — add an "Always/Never" line to ${projectRoot}/CLAUDE.md. The soft fallback when neither above can match deterministically.`,
    reference: `Claude Code guardrails in preference order: (1) permissions — settings.json "permissions": { "allow" | "deny" | "ask": ["Bash(rm -rf:*)", "Read(.env)", "Bash(git push --force:*)", …] }; (2) hooks — settings.json "hooks" with a matcher + command on tool events like PreToolUse; (3) CLAUDE.md "Always/Never" rules. Read the current ${projectRoot}/.claude/settings.json first. Prefer the NARROWEST deny/ask that stops the evidenced risk; do not over-block routine work.`,
  },
};

function dimensionMenu() {
  return Object.values(DIMENSIONS)
    .map((d) => `- ${d.label}: ${d.essence}`)
    .join("\n");
}

function dimensionBrief(kind) {
  const d = DIMENSIONS[kind] ?? DIMENSIONS.skill;
  return `Dimension: ${d.label}
Fix surface: ${d.surface}
Reference (read before drafting): ${d.reference}`;
}

// ── Structured verdicts, so the loop branches on data, not prose ──────────────────────────

const MINE = {
  type: "object",
  properties: {
    freshFindings: { type: "number", description: "count of NEW findings recorded THIS round (0 if the round added nothing)" },
    examinedSessions: { type: "array", items: { type: "string" }, description: "session ids / file names examined this round" },
    unexaminedRemain: { type: "boolean", description: "true iff sessions in the range were left un-examined and another pass could find more (sampling lenses); false once coverage is exhaustive" },
  },
  required: ["freshFindings"],
};

const CLUSTER = {
  type: "object",
  properties: {
    // ids of the proposal tasks the cluster agent CREATED via task.js (one per opportunity)
    proposalIds: { type: "array", items: { type: "string" } },
    droppedCount: { type: "number", description: "opportunities discarded as thin/duplicate" },
  },
  required: ["proposalIds"],
};

const PROPOSAL = {
  type: "object",
  properties: {
    target: { type: "string", description: "repo-relative path of the artifact to improve, OR the path of the new artifact to create" },
    targetKind: { type: "string", description: "instruction | skill | agent | workflow | output-style | guardrail" },
    action: { type: "string", description: "improve (edit an existing artifact) | create (a new skill/subagent/hook/permission rule)" },
    tier: { type: "string", description: "advisory | bounded | eval — the highest tier this change qualifies for" },
    hasEvalSuite: { type: "boolean", description: "true iff the target ships an evals/evals.json" },
    summary: { type: "string", description: "one line: what the change does and why" },
    // whether a concrete diff or full new-file body was written into the task body
    diffWritten: { type: "boolean" },
  },
  required: ["target", "action", "summary"],
};

const VERDICT = {
  type: "object",
  properties: {
    supported: { type: "boolean", description: "true iff the cited log evidence actually motivates this change" },
    reason: { type: "string" },
  },
  required: ["supported"],
};

// Report-quality grade — a rubric of per-dimension scores (fail | partial | pass). Daydream's rubric:
// every proposal must explain the BEHAVIOUR, the IMPROVEMENT, and the PROOF. The gate is per-dimension:
// any dimension scored `fail` drives a revision pass; `partial` ships with a note. `issues` localize the
// fixes for the reviser. The overall scalar is DERIVED from the scores by the script below — not declared
// by the model, so there is one source of truth. `dimensions` is extensible with problems the grader finds.
const GRADE = {
  type: "object",
  properties: {
    dimensions: {
      type: "array",
      items: {
        type: "object",
        properties: {
          name: { type: "string" },
          score: { type: "string", description: "fail | partial | pass" },
          note: { type: "string" },
        },
        required: ["name", "score"],
      },
    },
    issues: {
      type: "array",
      items: {
        type: "object",
        properties: {
          dimension: { type: "string" },
          detail: { type: "string", description: "the exact proposal at fault" },
          fix: { type: "string", description: "the concrete change that resolves it" },
        },
        required: ["dimension", "detail"],
      },
    },
  },
  required: ["dimensions"],
};

// ── Prompt builders ───────────────────────────────────────────────────────────────────────

// Shared block: where the data is + the iterative-round contract every mine lens obeys.
function mineCommon(seen, round) {
  const seenNote = seen.length
    ? `This is mine round ${round + 1}. You have ALREADY examined these sessions in a prior round — do NOT re-read them; examine the remainder:\n  ${seen.join(", ")}\nRead your task's CURRENT findings first and ADD only NEW findings (union with what is there); do not duplicate or drop prior findings.`
    : `This is the first mine round.`;
  return `## Where the data is
Session transcripts are JSONL under: ${logRoot}
- Top-level sessions:  <session-uuid>.jsonl
- Subagent transcripts (SEPARATE, must also be walked):  <session-uuid>/subagents/agent-<id>.jsonl
Reading this archive may require the shell sandbox to be DISABLED — it lives outside the project
root, so a plain read can fail with "Operation not permitted". If that happens, retry the command
with the sandbox disabled. Filter records to the date range by the ISO \`timestamp\` on each record.

## Iterative coverage
${seenNote}
When you finish, RETURN { freshFindings, examinedSessions, unexaminedRemain }:
- examinedSessions: the session ids/files you read this round.
- unexaminedRemain: true only if sessions in range were left un-read and another pass could find more.
  Exhaustive structured lenses (you scanned every session cheaply) return false; a sampling lens that
  skipped sessions returns true so the run can examine them next round.`;
}

// The harness surfaces a finding can be routed to — injected so mine agents tag each finding.
function mineDimensionTag() {
  return `## Tag each finding with the harness surface a fix would touch
${dimensionMenu()}
For each finding, name the single best surface above (its \`label\`) so clustering can route it, and say
whether it looks like IMPROVE (an existing artifact is at fault) or CREATE (a gap no artifact fills).`;
}

function minePrompt(taskId, seen, round) {
  return `You are a session-log mining agent (FRICTION lens) for the daydream run "${name}".
Purpose: ${description}
Date range: ${range}

Your assignment lives in a tracker task. Read it first:
  ${tj} read ${taskId} --home ${home}

${mineCommon(seen, round)}

## How to mine (two-pass join — highest-yield extraction)
1. Pass 1: from \`assistant\` records, build \`tool_use_id → { tool, attributionSkill }\`.
2. Pass 2: collect the signal for YOUR angle — e.g. \`tool_result\` blocks with \`is_error:true\`;
   \`[Request interrupted by user]\` markers; user-rejected tool uses; repeated identical tool_use
   inputs after an error in the \`parentUuid\` chain (retry loops); expensive turns from \`usage\`.
   Group by \`attributionSkill\` so each finding BLAMES a concrete skill/plugin.

${mineDimensionTag()}

## Record findings (CITED) back into your task
Every finding must cite its source — the session file + record (e.g.
\`<session-uuid>.jsonl record uuid=… / tool_use_id=…\`) and, where the fix points at code, the repo
\`file:line\`. Compose the findings as markdown and store them INLINE — do NOT create a findings.md (or
any report-shaped file); the subagent guard blocks report-file writes and a stray file pollutes the repo.
Pipe the markdown straight in via a shell heredoc:
  ${tj} update ${taskId} --home ${home} --body-file /dev/stdin <<'FINDINGS'
  ...your cited markdown findings (union with any prior round) here...
  FINDINGS
Then mark the task mined:  ${tj} update ${taskId} --home ${home} --status passed

Report counts (e.g. "skill X caused N errors of category C over the range"), not narratives. A finding
with no source is not a finding. Do NOT propose fixes here — only gather cited signal.`;
}

function conversationPrompt(taskId, seen, round) {
  return `You are the CONVERSATION-QUALITY mining agent for the daydream run "${name}".
Purpose: ${description}
Date range: ${range}

Read your assignment:  ${tj} read ${taskId} --home ${home}
You read the actual dialogue — user and assistant message CONTENT — and judge where the interaction was
inefficient, including inefficiencies that leave NO structured trace (nothing errored, nothing was
interrupted — the exchange was simply less efficient than it could have been). The goal is general
conversational efficiency, not a fixed checklist.

${mineCommon(seen, round)}
Read \`assistant\` records' \`message.content\` (text and thinking blocks) and \`user\` records' typed
prompts (string \`message.content\`).

## Scope by SAMPLING — do not read every session in full
There may be many sessions; reading all fully is prohibitively expensive. Prioritize the highest-signal
ones (most turns, longest, or those with visible back-and-forth) and take a representative sample of the
rest. REPORT exactly which sessions you sampled (examinedSessions) and set unexaminedRemain=true whenever
you leave sessions un-read — the run will hand them to you next round. Never imply full coverage.

## What to hunt (these are ILLUSTRATIONS, not an exhaustive list)
Surface ANY recurring pattern where the exchange was less efficient than it could have been, e.g.:
- Under-delivery → follow-up: the user had to ask for detail/scope the assistant should have supplied
  up front ("actually also…", "you forgot…", "give me the full X").
- Flip-flopping: the assistant reversed a recommendation with little new information, confusing the user.
- Over-asking / ceremony: menus or clarifying questions where a sensible default existed; verbosity the
  user pushed back on ("keep it terse", "no ceremony").
- Ignored-then-restated: a constraint the user had to repeat because it was dropped.
- Re-explanation loops: the assistant needed several tries to land a point.
- Premature confidence or unverified claims the user then had to challenge.
Invent categories the examples don't cover if the dialogue shows them.

${mineDimensionTag()}
An ignored-then-restated constraint often points at \`instruction\` (a CLAUDE.md rule that did not bite)
or a skill's Constraints; ceremony/verbosity often points at an \`output-style\`.

## What makes it a finding
A single awkward turn is noise. A PATTERN that recurs across sessions (or a costly single instance) is a
finding. For each: name the pattern, quote the MINIMAL telltale exchange (a few words, not whole turns),
cite every instance (session file + record uuid), and name the target surface.

## Record findings (CITED) inline — do NOT create a file
  ${tj} update ${taskId} --home ${home} --body-file /dev/stdin <<'FINDINGS'
  ...your cited conversation-quality patterns + surface per pattern (union with any prior round) here...
  FINDINGS
Then:  ${tj} update ${taskId} --home ${home} --status passed
Diagnose only; do NOT edit any file.`;
}

function routerPrompt(taskId, seen, round) {
  return `You are the MODEL-ROUTING mining agent for the daydream run "${name}".
Purpose: ${description}
Date range: ${range}

Read your assignment:  ${tj} read ${taskId} --home ${home}
You mine the session logs for where the wrong-sized LLM model was used — over-provisioned roles paying
for capability they didn't need, and under-provisioned roles that thrashed for lack of it.

${mineCommon(seen, round)}
Subagent role is in the sidecar agent-<id>.meta.json (\`agentType\`, \`description\` — e.g. a workflow's
mine/validate/code/eval agent).

## What to extract (per ROLE, not per session)
For each assistant record: \`message.model\` (which model ran), \`message.usage\` (token spend), and the
\`attributionSkill\`/subagent role it belongs to. Group by role. Then cross-reference the SAME role's
outcome signal in that session:
- OVER-provisioned → a strong model (e.g. opus) doing mechanical, schema-filled, or high-fan-out work
  that was never corrected, retried, or eval-failed. Candidate to DOWN-tier (opus→sonnet→haiku) or
  lower \`effort\`.
- UNDER-provisioned → a cheaper model (haiku/sonnet) that THRASHED: repeated retries after error in the
  \`parentUuid\` chain, user corrections, or eval-loop failures. Candidate to UP-tier / raise \`effort\`.
A verifier/gate role must never be recommended BELOW the generator it checks — call that out if seen.

${mineDimensionTag()}
Routing findings target the \`workflow\` surface (an \`agent(..., { model, effort })\` call in a golden
copy) or the \`agent\` surface (a \`model:\` frontmatter in a team agent) — always IMPROVE.

## Record findings (CITED) inline
Every recommendation cites: the session file + record uuid, the \`message.model\` and \`usage\` numbers,
and the outcome signal (retry/correction/eval-fail record, or its absence).
  ${tj} update ${taskId} --home ${home} --body-file /dev/stdin <<'FINDINGS'
  ...your cited routing findings + per-role model/effort recommendation (union with any prior round) here...
  FINDINGS
Then:  ${tj} update ${taskId} --home ${home} --status passed
Recommend routing only; do NOT edit any file.`;
}

function capabilityPrompt(taskId, seen, round) {
  return `You are the CAPABILITY-GAP mining agent for the daydream run "${name}".
Purpose: ${description}
Date range: ${range}

Read your assignment:  ${tj} read ${taskId} --home ${home}
You hunt POSITIVE space that the failure lenses cannot see: recurring MANUAL multi-step work the
assistant performed BY HAND — successfully — across sessions, that has no owning skill or subagent. Each
such pattern is a candidate to CREATE new functionality — a new SKILL (a procedure the main loop runs) or
a new SUBAGENT (specialized, isolatable work worth delegating).

${mineCommon(seen, round)}
Read the assistant's tool-use SEQUENCES and thinking/text content to spot repeated procedures.

## What makes a capability gap a finding
- A multi-step sequence (≥3 steps) the assistant reconstructed from scratch in ≥2 sessions — e.g. the
  same setup/scan/verify dance, the same doc-generation ritual, the same release/checklist steps — with
  NO skill or subagent in the harness that already owns it. (Cross-check the skills/agents the harness
  system prompt lists.)
- Prefer patterns that were done WELL (not corrected) — you are packaging proven work, not fixing bugs.
- A one-off is noise; recurrence across sessions is the signal.
- For each, note whether it reads as a SKILL (a reusable procedure) or a SUBAGENT (delegated specialized
  work) — clustering makes the final call.

${mineDimensionTag()}
These findings are CREATE candidates on the \`skill\` or \`agent\` surface: name the sequence, the sessions
it recurred in, and what the new artifact would encapsulate.

## Record findings (CITED) inline
Cite every session + record uuid where the sequence recurred, and list the concrete steps.
  ${tj} update ${taskId} --home ${home} --body-file /dev/stdin <<'FINDINGS'
  ...your cited recurring-manual-work patterns + proposed new skill/subagent scope (union with any prior round)...
  FINDINGS
Then:  ${tj} update ${taskId} --home ${home} --status passed
Diagnose only; do NOT create any file.`;
}

function guardrailPrompt(taskId, seen, round) {
  return `You are the GUARDRAIL mining agent for the daydream run "${name}".
Purpose: ${description}
Date range: ${range}

Read your assignment:  ${tj} read ${taskId} --home ${home}
You hunt two things the other lenses miss, both of which warrant a preventive guardrail:

A. Recurring MISTAKE classes — the same wrong action or omission repeated across sessions, EVEN when each
   instance was individually corrected (e.g. editing the wrong file, skipping a required check, a wrong
   path/flag). The repetition, not the single instance, is the signal.
B. RISK exposure — dangerous, destructive, or sensitive operations that ran with NO guardrail, whether or
   not they errored. A destructive command that SUCCEEDED every time is never "corrected" yet is the
   highest-value thing to gate. Examples: \`rm -rf …\`, \`git push --force\`, \`chmod -R\`, \`sudo …\`,
   \`curl … | sh\`, mass deletes, reads/writes of \`.env\` / secrets / credentials, writes to production.
   For risk, SEVERITY can substitute for frequency — one \`rm -rf\` of a real path justifies a guardrail.

${mineCommon(seen, round)}

## For each finding, sketch the guardrail (Claude Code, preference order)
1. PERMISSION POLICY — a deny/ask entry in .claude/settings.json "permissions" (e.g. deny "Bash(rm -rf:*)",
   ask on "Read(.env)"). Best for dangerous/sensitive operations that can be matched statically.
2. HOOK — a command on a tool event (e.g. PreToolUse) in settings.json "hooks", for context-dependent checks.
3. RULE — an "Always/Never" line in CLAUDE.md, when neither above can match deterministically.
Characterise the trigger precisely enough that a reader could write the matcher. If it is too fuzzy to
match, say so — it becomes a soft CLAUDE.md rule, not a permission/hook.

${mineDimensionTag()}
These findings target the \`guardrail\` surface — usually CREATE (a new permission rule / hook / CLAUDE.md
rule), occasionally IMPROVE (tighten an existing one).

## Record findings (CITED) inline
Cite every session + record uuid where the mistake/operation appeared, and (for mistakes) the correction.
  ${tj} update ${taskId} --home ${home} --body-file /dev/stdin <<'FINDINGS'
  ...your cited mistake classes / risk exposures + guardrail mechanism sketch (union with any prior round)...
  FINDINGS
Then:  ${tj} update ${taskId} --home ${home} --status passed
Diagnose only; do NOT edit any file.`;
}

function instructionPrompt(taskId, seen, round) {
  return `You are the INSTRUCTION-HYGIENE mining agent for the daydream run "${name}".
Purpose: ${description}
Date range: ${range}

Read your assignment:  ${tj} read ${taskId} --home ${home}
You audit the baseline system instruction — the project ${projectRoot}/CLAUDE.md and ${projectRoot}/AGENTS.md
(read them first) — against how sessions ACTUALLY behaved, to find instruction that should change:
- DEAD rules — a rule the sessions never triggered / never needed (candidate to remove, to cut bloat).
- CONTRADICTED / IGNORED rules — behaviour that repeatedly violated a stated rule (the rule did not bite
  → sharpen it, move it, or make it a hook via the guardrail surface).
- REDUNDANT / BLOATED passages — length that added no behavioural signal (candidate to simplify).
- AMBIGUOUS rules — a rule the assistant repeatedly interpreted differently or asked to clarify.

${mineCommon(seen, round)}

## What makes an instruction finding
Tie every claim to behaviour: a rule is "dead" only if you looked and the sessions never invoked it; a
rule is "ignored" only if you can cite records that violated it. Quote the exact rule text (with its
CLAUDE.md/AGENTS.md location) and the sessions that justify the change.

${mineDimensionTag()}
These findings target the \`instruction\` surface — always IMPROVE (a simplification REMOVES or tightens
text; do not merely append). Ignored rules that need enforcement may instead target \`guardrail\`.

## Record findings (CITED) inline
  ${tj} update ${taskId} --home ${home} --body-file /dev/stdin <<'FINDINGS'
  ...your cited instruction-hygiene findings (rule text + location + behaviour, union with any prior round)...
  FINDINGS
Then:  ${tj} update ${taskId} --home ${home} --status passed
Diagnose only; do NOT edit any file.`;
}

// Dispatch a lens by kind. New lenses degrade to the friction lens if an unknown kind slips through.
function minePromptFor(angle, seen, round) {
  switch (angle.kind) {
    case "conversation": return conversationPrompt(angle.id, seen, round);
    case "router":       return routerPrompt(angle.id, seen, round);
    case "capability":   return capabilityPrompt(angle.id, seen, round);
    case "guardrail":    return guardrailPrompt(angle.id, seen, round);
    case "instruction":  return instructionPrompt(angle.id, seen, round);
    default:             return minePrompt(angle.id, seen, round);
  }
}

// Semantic lenses are judgment-heavy → strong model; structured lenses stay on the inherited model.
function mineOptsFor(angle, round) {
  const base = { phase: "Mine", label: `${angle.kind ?? "mine"}:${angle.id}:r${round + 1}`, schema: MINE };
  if (angle.kind === "conversation" || angle.kind === "capability" || angle.kind === "instruction") {
    return { ...base, model: "opus", effort: "high" };
  }
  return base;
}

function clusterPrompt() {
  return `You are the clustering gate for the daydream run "${name}".
Purpose: ${description}

## Read every mined finding
  ${tj} graph --home ${home}
  ${tj} list --home ${home} --type mine
Read each mine task's findings:  ${tj} read <id> --home ${home}

## The harness surfaces an opportunity may target
${dimensionMenu()}

## Turn signal into distinct improvement OPPORTUNITIES
- Dedup: collapse findings that point at the same root cause / same target.
- Rank by (frequency × severity): a category that recurred often and blocked work — or a severe risk —
  outranks a one-off.
- Each opportunity names ONE surface above, and you make the IMPROVE-vs-CREATE call explicitly:
  * PREFER **improve** — bind it to an existing artifact (give its repo-relative path) whenever one
    plausibly owns the fix.
  * Choose **create** only when (a) no existing artifact fits AND (b) the pattern recurs across ≥2
    sessions — or, for a guardrail, the operation is dangerous enough that one occurrence justifies it.
    Give the path of the new artifact (a new skill dir, a new subagent .md, a permissions/hook entry in
    .claude/settings.json).
- Drop thin or unactionable signal (transient env errors, one-offs with no pattern). Count what you drop.
- Keep at most ${maxProposals} opportunities.

## Create one proposal task per opportunity
State the surface, the action (improve|create), the target path, and the cited evidence in the body:
  ${tj} create --home ${home} --type proposal --name "<surface> [improve|create] <target>: <one-line opportunity>" --body "<action + target path + cited evidence + why>"
Collect the ids you create — return them as proposalIds. Also return droppedCount.`;
}

function proposePrompt(taskId) {
  return `You are a proposal agent for the daydream run "${name}".
Read your opportunity task:  ${tj} read ${taskId} --home ${home}
It names ONE harness surface, the ACTION (improve an existing artifact, or create a new one), the target
path, and the cited log evidence motivating the change.

## Know the surface you are changing
Match the opportunity's surface to the brief below and READ the reference before drafting:
${Object.keys(DIMENSIONS).map((k) => dimensionBrief(k)).join("\n\n")}

## Draft the change
- action = improve: read the target artifact in full first; write a MINIMAL, surgical diff (unified-diff
  style, or clearly-marked before/after) that addresses the evidenced problem — nothing speculative.
- action = create: follow the surface's reference/conventions and write the FULL initial content of the
  new artifact (a new skill's SKILL.md; a new subagent .md; a settings.json "permissions"/"hooks"
  fragment). Keep it the minimum that solves the evidenced need.
- For a guardrail, choose the mechanism by the preference order in the brief (permission deny/ask →
  hook → CLAUDE.md rule) and write the exact settings.json entry or rule text. Prefer the NARROWEST
  matcher that stops the evidenced risk without blocking routine work.

## Classify the safety tier (advisory metadata for the human — daydream never auto-applies)
- \`eval\`   — the target directory ships an evals/evals.json that could gate the change.
- \`bounded\` — no eval, but the change is small and structurally checkable (a frontmatter/section patch,
  a single permissions entry).
- \`advisory\` — anything larger, judgement-heavy, or a CREATE / a guardrail / an instruction edit; a
  human must apply and verify it. New artifacts, permission rules, and hooks are ALWAYS advisory.
Report hasEvalSuite (true iff evals/evals.json exists next to the target) and the action.

## Store the proposal
Store the diff or full new-artifact body + rationale (keeping every evidence citation) INLINE — do NOT
create a proposal.md or any scratch file (it pollutes the repo). Pipe it straight in via a heredoc:
  ${tj} update ${taskId} --home ${home} --body-file /dev/stdin <<'PROPOSAL'
  ...your diff / new-artifact body + rationale + citations here...
  PROPOSAL
Return { target, targetKind, action, tier, hasEvalSuite, summary, diffWritten }.
Do NOT edit the real artifact — only propose.`;
}

function validatePrompt(taskId) {
  return `You are the ADVERSARIAL validation gate for a daydream proposal (default to disbelief).
Read the proposal:  ${tj} read ${taskId} --home ${home}

Assume the proposed change is UNJUSTIFIED until the cited log evidence forces the opposite conclusion.
For each citation the proposal rests on:
- OPEN the cited source yourself — the session-log record it points to (disable the sandbox if the
  read is denied) and any repo \`file:line\` — and confirm it ACTUALLY shows the failure/friction/risk the
  proposal claims. A citation that does not substantiate the claim counts as UNSUPPORTED.
- The proposal is supported ONLY when the evidence genuinely motivates THIS change (not merely "an error
  happened nearby"). For a CREATE, the evidence must show the recurring need it fills. For a guardrail,
  the mistake must actually RECUR — or the operation be genuinely dangerous — not a single benign case.
  The change must plausibly reduce the evidenced friction/risk.
When in doubt, rule NOT supported — a dropped weak proposal is better than a wrong recommendation.

Return { supported, reason }.`;
}

function gradePrompt() {
  return `You are the report grader for the daydream run "${name}". Grade the proposals report at
${reportPath} against the rubric below. You do NOT fix anything — you assess, and list what must change.
Every proposal must stand on three legs: the BEHAVIOUR it addresses, the IMPROVEMENT it makes, and the
PROOF behind it.

Score each dimension: \`pass\` (every proposal meets it), \`partial\` (mostly, with a fixable gap), or
\`fail\` (one or more proposals miss it). \`fail\` forces a revision; \`partial\` ships with a note.

## Rubric — across all proposals in the report
1. Behaviour — each proposal names the concrete behaviour/pattern/risk it addresses: what the system
   actually did that was inefficient, wrong, or dangerous — not merely "an error occurred nearby".
   \`fail\` if any is vague/missing.
2. Improvement — each carries a concrete, minimal change (a diff, a clearly-marked before/after, or the
   full body of a new artifact/permission rule) to a real harness surface, labelled improve or create,
   not a vague suggestion. \`fail\` if any lacks an actionable change.
3. Proof — each behaviour is backed by traceable evidence (session-log record uuid and/or repo
   \`file:line\`) a reader can open, that genuinely motivates THIS change. For a create or a guardrail,
   the evidence shows the pattern RECURRED (or the risk is genuinely severe). \`fail\` if any is unsourced
   or unsupported.
4. Faithfulness — the report contains no proposal absent from the validated tracker; nothing invented.
   \`fail\` if any is invented.
If a proposal shows a quality problem none of these name, add it as your own dimension and score it.

## Return
- dimensions: score EACH dimension { name, score ("fail" | "partial" | "pass"), note }.
- issues: for every dimension scored below \`pass\`, one entry per concrete problem { dimension, detail,
  fix } — name the exact proposal and the change that resolves it. Be specific.`;
}

function synthesizePrompt(kept, issues) {
  if (issues && issues.length) {
    const fixes = issues
      .map((x) => `- ${x.dimension}: ${x.detail}${x.fix ? ` → ${x.fix}` : ""}`)
      .join("\n");
    return `You are revising the daydream proposals report at ${reportPath} to fix specific grading issues.
Read the current report first.
Fix EACH issue below, changing only what it requires — preserve every other proposal, diff, and citation:
${fixes}
If closing an issue needs evidence or a diff, re-open the proposal task it refers to:
  ${tj} list --home ${home} --type proposal
  ${tj} read <id> --home ${home}
Do NOT invent a proposal that is not in the tracker, and do NOT weaken any citation.
Rewrite the file in place. Report the path and a one-line summary of what you changed.`;
  }
  const list = kept
    .map((p, i) => `${i + 1}. [${p.tier ?? "advisory"}] [${(p.action ?? "improve").toUpperCase()}] ${p.targetKind ?? "?"} — ${p.target} — ${p.summary} (task ${p.id})`)
    .join("\n");
  return `You are the synthesizer for the daydream run "${name}". Proposals are drafted and validated.
Purpose: ${description}
Date range: ${range}

## The validated proposals to include (read each task for its change + evidence)
${list || "(none survived validation)"}
Read each:  ${tj} read <id> --home ${home}

## Write the report to: ${reportPath}
- Open with a short summary: what range was mined, how many findings/opportunities, how many proposals
  survived validation, and the per-surface ranking (instruction / skill / subagent / workflow /
  output-style / guardrail — how many proposals each, split by improve vs create, and the top friction/risk).
- Then one section per proposal, GROUPED by harness surface and ordered by tier then rank: the concrete
  diff or new-artifact body, the IMPROVE/CREATE label, the one-line rationale, the cited log evidence, and
  the safety tier.
- Preserve EVERY citation. Do NOT invent proposals not in the tracker.
- Close with a "How to apply" note: daydream is ADVISORY — it reports; a human reviews each change and
  applies it on a branch/PR. Note which proposals ship an eval suite (tier \`eval\`) as safest to automate
  later, and that new artifacts, permission rules, and hooks always need human review before landing.

Write the file. Report the path and a one-line summary.`;
}

// ── Orchestration: iterative mine → cluster → (propose → validate) pipeline → synthesize ───

const mineAngles = angles.filter((a) => a && a.id);
const mineIds = mineAngles.map((a) => a.id);

// Iterative mining ("keep examining within the window"): re-run only the lenses that report
// un-examined sessions, feeding each its already-examined set, until coverage is exhausted or the
// round cap. Structured lenses finish in one pass (unexaminedRemain=false); sampling lenses (e.g.
// conversation) get additional passes over the sessions they skipped.
const examined = new Map(); // angle.id -> Set of session ids examined so far
let totalFresh = 0;
let active = mineAngles.slice();
const rounds = Math.max(1, maxMineRounds);

for (let round = 0; active.length && round < rounds; round++) {
  const results = await parallel(
    active.map((a) => () => {
      const seen = [...(examined.get(a.id) ?? new Set())];
      return agent(minePromptFor(a, seen, round), mineOptsFor(a, round)).then((r) => ({ a, r }));
    }),
  );

  const next = [];
  for (const res of results.filter(Boolean)) {
    const { a, r } = res;
    const set = examined.get(a.id) ?? new Set();
    (r?.examinedSessions ?? []).forEach((s) => set.add(s));
    examined.set(a.id, set);
    totalFresh += r?.freshFindings ?? 0;
    if (r?.unexaminedRemain && round + 1 < rounds) next.push(a);
  }
  active = next;
  if (active.length) log(`mine round ${round + 1}: ${active.length} lens(es) still have un-examined sessions → another pass`);
}

const cluster = await agent(clusterPrompt(), { phase: "Cluster", label: "cluster", schema: CLUSTER });
const proposalIds = (cluster?.proposalIds ?? []).filter(Boolean).slice(0, maxProposals);

let proposals = [];
if (proposalIds.length) {
  proposals = await pipeline(
    proposalIds,
    (id) => agent(proposePrompt(id), { phase: "Propose", label: `propose:${id}`, schema: PROPOSAL }).then((p) => ({ ...p, id })),
    (p, id) =>
      agent(validatePrompt(id), { phase: "Validate", label: `validate:${id}`, schema: VERDICT }).then((v) => ({ ...p, id, verdict: v })),
  );
}

const kept = proposals.filter(Boolean).filter((p) => p.verdict?.supported);
const dropped = proposals.filter(Boolean).length - kept.length;
log(`daydream "${name}": ${mineIds.length} lens(es) mined (${totalFresh} findings) → ${proposalIds.length} opportunit(ies) → ${kept.length} proposal(s) validated (${dropped} dropped).`);

// Synthesize the proposals report, then grade it; blocking issues drive a revision pass. Loop until the
// report clears the grade or the cap — the grade gate has teeth, it does not just attach a warning.
const maxGradeRounds = 3; // initial synthesis + up to 2 revisions
let grade = null;

for (let g = 0; ; g++) {
  await agent(synthesizePrompt(kept, g === 0 ? null : grade?.issues ?? []), {
    phase: "Synthesize",
    label: g === 0 ? "synthesize" : `revise:r${g}`,
  });

  grade = await agent(gradePrompt(), { phase: "Grade", label: `grade:r${g + 1}`, schema: GRADE });

  const failing = (grade?.dimensions ?? []).filter((d) => d.score === "fail");
  if (!failing.length || g + 1 >= maxGradeRounds) {
    if (failing.length) log(`daydream "${name}": report graded with ${failing.length} failing dimension(s) after ${g + 1} pass(es).`);
    break;
  }
  log(`grade round ${g + 1}: ${failing.length} failing dimension(s) → revising report`);
}

// Derived rubric value — single source of truth is the per-dimension scores (fail=0, partial=1, pass=2).
if (grade?.dimensions?.length) {
  const pts = grade.dimensions.reduce((s, d) => s + (({ fail: 0, partial: 1, pass: 2 })[d.score] ?? 0), 0);
  grade.overall = `${pts}/${grade.dimensions.length * 2}`;
}

return {
  name,
  dateRange,
  mined: mineIds.length,
  findings: totalFresh,
  opportunities: proposalIds.length,
  proposals: proposals.filter(Boolean).length,
  supported: kept.length,
  dropped,
  applyTier: "advisory",
  reportPath,
  grade,
};
