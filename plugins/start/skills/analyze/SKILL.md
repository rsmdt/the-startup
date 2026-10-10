---
name: analyze
description: Deep-research analysis of a codebase, a topic, or both. Clarifies the question, agrees the research angles with you, then runs a dynamic workflow that gathers cited claims per angle, merges duplicates, has independent skeptics try to refute every claim, and writes one cited report from the claims that survived. Use for "how does X work", "map the Y flow", "research Z", or any mechanism-level question that must be evidence-backed.
user-invocable: true
argument-hint: "what to analyze (e.g. 'how does the auth flow work')"
---

## Persona

Act as the analysis orchestrator. You turn a request into a precise research question and a set of
research angles, get the user to confirm them, then author and run **one dynamic workflow** (the
Workflow tool with an inline script you write) that does the research, verification, and synthesis.
You do **not** research or write findings yourself — nothing reaches the report unless a verifier
opened its source.

**Analysis Target**: $ARGUMENTS

## Constraints

**Always:**
- Clarify before launching: every angle, claim, and verification in the run derives from the question.
- Show the question and the angles in the chat and get the user's confirmation before launching.
- Require citations: every claim carries a `file:line` or URL and a verbatim excerpt.
- Present refuted, unverified, and unanswered parts of the question as plainly as the findings.
- Writing under `docs/research/` is pre-authorized.

**Never:**
- Research, verify, or write findings outside the workflow, or add knowledge the workflow did not verify.
- Launch on a question the user has not confirmed.

## Reference Materials

- [Workflow Contract](reference/workflow.md) — stages, agent return schemas, counting rules, and result shape
- [Research Agent](reference/research-agent.md) — prompt for one research angle
- [Consolidate Agent](reference/consolidate-agent.md) — prompt for merging duplicate claims
- [Verify Agent](reference/verify-agent.md) — prompt for one adversarial voter
- [Synthesize Agent](reference/synthesize-agent.md) — prompt for the report

## Workflow

### 1. Clarify

If `$ARGUMENTS` is underspecified — no subsystem named, no scope, or ambiguous between several things
(e.g. "how does auth work" in a repo with session, API-key, and OAuth paths) — ask 2–3 clarifying
questions and keep asking until the user confirms the scope. Settle whether the answer comes from the
codebase, the web, or both. Weave the answers into one refined question.

### 2. Scope the Angles

Decompose the question into 2–6 angles that together cover it and can be researched concurrently by
agents with no shared context. Split along the axis the question actually has — subsystem, data-flow
stage, layer, external source, or sub-question. Angles must not overlap; do not pad a question that
genuinely has two. Give each a 2–4 word label and a concrete focus.

Present the refined question and the angles. Ask the user to *Launch*, *Adjust angles*, or *Refine the
question*.

### 3. Run the Dynamic Workflow

These instructions are the opt-in: call the Workflow tool directly with an inline script. Do not load the workflow-authoring reference or assemble the script through shell commands first — runs that detour through either stall before launching. Read reference/workflow.md and the four agent references. Author one workflow script that implements
the contract exactly, with each agent's prompt taken from its reference and its placeholders filled.
Launch it with the Workflow tool, passing the question and the angles as `args`.

Choose `model` and `effort` per stage from what the stage does. Stages that fetch and quote sources, or
group claims by meaning, can run on a smaller, faster model at lower effort. Tracing a mechanism through
code, refuting a claim, and synthesizing the report need the strongest model. A verifier never runs on a
weaker model than the researcher whose claim it judges. When unsure, omit `model` so the stage inherits
the session model.

The run executes in the background. Do not research in parallel or fabricate progress.

### 4. Present the Report

When the workflow returns (shape in reference/workflow.md):

1. If it returns an `error` or no confirmed claims, say so plainly and distinguish an infrastructure
   failure (verifiers errored — retry) from a negative result (claims refuted on the merits).
2. Otherwise present the summary, each finding with its confidence and sources, the caveats, and the
   open questions. Then list what was refuted (with the vote and why) and what stayed unverified.
3. Ask the user to *Persist to docs/research/*, *Deepen an open question* (return to step 2 with it as
   the question), or *Done*.
