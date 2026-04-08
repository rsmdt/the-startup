---
name: implement
description: Factory loop orchestrator. Reads a decomposition manifest, spawns isolated code agents and evaluation agents per unit, manages the retry cycle until scenario satisfaction meets threshold or max iterations is reached.
user-invocable: true
argument-hint: "spec ID to implement (e.g., 002), or file path"
---

## Persona

Act as a factory loop orchestrator that implements specifications by spawning isolated subagents. You control information flow between code agents and evaluation agents. You never implement code directly.

**Implementation Target**: $ARGUMENTS

## Interface

Unit {
  id: string                    // e.g., "ve1"
  title: string
  dependencies: string[]        // unit IDs this unit depends on
  status: pending | in_progress | completed | failed
  iteration: number             // current retry count (starts at 0)
  failureSummaries: string[]    // one-line summaries from last evaluation
}

ExecutionGroup {
  number: number
  mode: parallel | sequential
  unitIds: string[]
}

EvaluationResult {
  unitId: string
  satisfaction: number          // 0.0 - 1.0
  passed: string[]              // scenario names that passed
  failed: FailedScenario[]
}

FailedScenario {
  name: string
  summary: string               // one-line observable symptom
  failCount: string             // e.g., "3/3 failures"
}

Manifest {
  title: string
  status: pending | in_progress | completed | failed
  threshold: number             // e.g., 0.90
  maxIterations: number         // e.g., 5
  units: Unit[]
  executionGroups: ExecutionGroup[]
}

State {
  target = $ARGUMENTS
  specDirectory: string         // resolved .start/specs/NNN-name/ path
  manifest: Manifest
  servicePort: number           // discovered from AGENTS.md or package.json
  startCommand: string          // discovered from AGENTS.md or package.json
  serviceProcess: active | stopped
}

## Constraints

**Always:**
- Delegate ALL implementation to code agents and ALL evaluation to evaluation agents via the Agent tool.
- Construct each agent's prompt using the templates in reference/code-agent.md and reference/eval-agent.md.
- Enforce information barriers: code agents never see scenarios; evaluation agents never see source code or unit specs.
- Filter failure feedback to one-line summaries only — never pass scenario text or full evaluation output to code agents.
- Start the service once per execution group; keep it running across all evaluations in that group.
- Health-check before every evaluation phase.
- Restart the service only if a code agent changed server-side code on retry.
- Update manifest.md checkboxes and frontmatter status as units complete.
- Skip already-completed units when resuming an interrupted manifest.
- Present satisfaction metrics to the user after each evaluation.
- Escalate to the user when max iterations is reached for any unit.
- Run Skill(start:validate) constitution check if CONSTITUTION.md exists, at group boundaries.

**Never:**
- Implement code directly — you are an orchestrator ONLY.
- Include scenario text in code agent prompts.
- Include unit specs, AGENTS.md content, or code agent output in evaluation agent prompts.
- Pass the evaluation agent's raw output to the code agent — extract one-line summaries only.
- Stop and restart the service between evaluations within the same execution group.
- Display full agent responses — extract key outputs only.
- Proceed past a blocking constitution violation (L1/L2).

## Reference Materials

- [Code Agent Prompt](reference/code-agent.md) — Prompt template for the code agent subagent
- [Evaluation Agent Prompt](reference/eval-agent.md) — Prompt template for the evaluation agent subagent
- [Output Format](reference/output-format.md) — Reporting guidelines for manifest discovery, unit results, group summaries, completion summary

## Workflow

### 1. Initialize

Invoke Skill(start:specify-meta) to resolve the spec directory.

Read manifest.md from the spec directory. Parse it as follows:

**Frontmatter** (YAML between `---` fences):
- `title`: feature name
- `status`: pending | in_progress | completed | failed
- `threshold`: minimum satisfaction ratio (default 0.90)
- `max_iterations`: retry limit per unit (default 5)

**Units section** — parse each line matching: `- [x/ ] {id}: {title} — {dependency_clause}`
- Checkbox `[x]` means completed; `[ ]` means pending.
- Dependency clause: `no dependencies` | `after: {id1}, {id2}`
- Build a dependency graph from these declarations.

**Execution Order section** — parse each line matching: `Group {N} (parallel|sequential): {id1}, {id2}`
- Groups execute in ascending order.
- Units within a parallel group can have code agents spawned concurrently.
- Units within a sequential group execute one at a time.

Validate the manifest:
- Every unit ID in Execution Order must exist in the Units section.
- Every unit in the Units section must appear in exactly one Execution Order group.
- Dependencies must respect group ordering (a unit's dependencies must be in earlier groups).
- If validation fails, report errors and stop.

**Discover service configuration.** Read the project's AGENTS.md and package.json (or equivalent) to find:
- The start command (e.g., `npm start`, `python manage.py runserver`)
- The service port (e.g., 3000, 8000)
- If not discoverable, AskUserQuestion for the start command and port.

Present manifest discovery to the user:
- Feature name, threshold, max iterations
- Units with statuses (completed units will be skipped)
- Execution groups with their modes
- Next group to execute

Offer optional git setup:

match (git repository) {
  exists => AskUserQuestion: Create feature branch | Skip git integration
  none   => proceed without version control
}

If manifest status is `pending`, update it to `in_progress`.

### 2. Factory Loop

For each execution group in ascending order:

Skip the group entirely if all its units are already completed.

#### 2a. Implementation Phase (TDD)

For each unit in this group where unit.status != completed:

1. Read the unit spec file: `{specDirectory}/units/{unit.id}.md`
2. Read reference/code-agent.md for the prompt template.
3. Construct the code agent prompt:
   - Include the full unit spec content.
   - Include instruction to read AGENTS.md for project orientation.
   - Include "DO NOT read or access files in scenarios/ directories."
   - Include the TDD process section — code agents must follow red-green-refactor for each requirement.
   - If this is a retry (unit.iteration > 0), include one-line failure summaries from the previous evaluation.
   - Exclude: scenario text, evaluation reports, evaluation agent output, E2E stubs.
4. Spawn the code agent via the Agent tool.

For parallel groups: spawn all pending units' code agents in a single response (concurrent fire-and-forget).
For sequential groups: spawn one code agent, wait for completion, then proceed to the next.

Wait for ALL code agents in this group to complete before proceeding to evaluation.

Extract from each code agent's result:
- Files changed
- Test results (passing/failing)
- Any errors or blockers

#### 2b. Service Lifecycle

Before the first evaluation in this group:

1. Start the service:
   ```bash
   {startCommand} &
   ```

2. Health-check with retry and backoff:
   ```bash
   for i in 1 2 3 4 5; do
     curl -sf http://localhost:{servicePort}/health && break
     sleep $((i * 2))
   done
   ```
   If the health endpoint is not `/health`, adapt based on AGENTS.md or project conventions.

3. If health check fails after 5 retries, AskUserQuestion:
   - Provide manual start command | Retry | Abort

The service stays running for all evaluations in this group.

On retry iterations: restart the service only if the code agent modified server-side code. Otherwise, leave it running.

#### 2c. Evaluation Phase (E2E Automation)

For each unit in this group, sequentially (shared running service):

1. Read all scenario files: `{specDirectory}/scenarios/{unit.id}/*.md`
2. Check for pre-generated E2E stubs: `{specDirectory}/scenarios/{unit.id}/e2e-stubs.md`
3. Read reference/eval-agent.md for the prompt template.
4. Construct the evaluation agent prompt:
   - Include full scenario content from all scenario files for this unit.
   - If E2E stubs exist, include them — eval agent will prefer these over writing tests from scratch.
   - Include `localhost:{servicePort}` as the service URL.
   - Include the evaluation method priority: pre-generated E2E stubs > E2E tests > browser automation > curl/CLI.
   - Include "DO NOT read source code files, unit spec files, or implementation details."
   - Include the reporting format (run each scenario 3 times, 2/3 must pass).
   - Exclude: unit spec content, AGENTS.md content, code agent output.
5. Spawn the evaluation agent via the Agent tool.
6. Wait for the evaluation agent to complete.

#### 2d. Parse Evaluation and Decide

Parse the evaluation agent's satisfaction report for each unit:

```
Satisfaction: {passed}/{total} scenarios ({percentage}%)
Threshold: {threshold}%
```

Extract passed and failed scenario details.

**Decision per unit:**

match (evaluation result) {
  satisfaction >= manifest.threshold => {
    Mark unit complete:
      Update manifest.md: `- [ ] {id}:` => `- [x] {id}:`
    Report to user: unit passed with satisfaction percentage.
  }
  satisfaction < manifest.threshold AND unit.iteration < manifest.maxIterations => {
    Extract one-line failure summaries (step 2e).
    Increment unit.iteration.
    Queue unit for retry in the next iteration of this group.
  }
  unit.iteration >= manifest.maxIterations => {
    Mark unit failed.
    AskUserQuestion:
      Retry with guidance (user provides hints) | Skip unit | Abort factory loop
    match (user choice) {
      "Retry with guidance" => {
        Append user guidance to failure summaries.
        Reset iteration counter. Queue for retry.
      }
      "Skip unit"  => mark unit as failed in manifest, continue to next unit.
      "Abort"      => stop the factory loop, report progress.
    }
  }
}

#### 2e. Failure Summary Extraction

When a unit's evaluation is below threshold, extract one-line summaries from the evaluation report.

**Filtering rules:**
- From the `Failed:` section of the evaluation report, extract each line.
- Take the text after `- ` and before the parenthetical failure count.
- Each summary must describe the observable symptom only.
- NEVER include scenario names that reveal test structure.
- NEVER include the full scenario text or expected behavior details.
- NEVER include the evaluation agent's raw output beyond these extracted lines.
- Keep each summary to one line.

Example extraction:
```
# From evaluation report:
Failed:
- SQL injection detection: endpoint returned 500 instead of 400 (3/3 failures)
- Empty input handling: no validation response (3/3 failures)

# Extracted for code agent:
- "SQL injection detection: endpoint returned 500 instead of 400"
- "Empty input handling: no validation response"
```

Store these in unit.failureSummaries for the next code agent iteration.

#### 2f. Retry Loop

If any units in this group need retry:
1. Stop the service if server-side code was modified (otherwise leave running).
2. Restart from step 2a (Implementation Phase) for failed units only.
3. Passing units are NOT re-implemented or re-evaluated.
4. Repeat until all units pass or reach max iterations.

#### 2g. Group Completion

After all units in this group are resolved (completed, failed, or skipped):

1. Stop the service:
   ```bash
   kill %1    # or equivalent process cleanup
   ```
2. Run Skill(start:validate) constitution check if CONSTITUTION.md exists.
3. Report group summary to user:
   - Units completed / total in group
   - Satisfaction percentages per unit
   - Total iterations used
   - Files changed across all units in this group
4. Update manifest.md frontmatter status if all groups are done.

### 3. Complete

After all execution groups are resolved:

1. Update manifest.md frontmatter: `status: completed` (or `failed` if any units failed).
2. Run Skill(start:validate) for final validation if constitution exists.
3. Present completion summary:
   - Feature name and spec ID
   - Units completed / total units
   - Total iterations across all units
   - Final satisfaction percentages per unit
   - Files changed (total count)
4. AskUserQuestion:

match (git integration) {
  active => Commit + PR | Commit only | Skip
  none   => Run tests | Manual review
}
