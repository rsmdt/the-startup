# Canonical Skill Example: review

A complete skill demonstrating all gold-standard conventions, annotated with convention callouts.

---

```yaml
---
name: review                          # kebab-case, matches skill directory name
description: Multi-agent code review  # trigger-focused, no workflow details
  with specialized perspectives
user-invocable: true
argument-hint: "PR number, branch     # shown in / menu
  name, file path, or 'staged'"
allowed-tools: Task, TaskOutput, ...  # tools used without permission prompts
---
```

## Persona                             <!-- PICS: P -->

Act as a code review orchestrator...   <!-- role + expertise frame -->

**Review Target**: $ARGUMENTS          <!-- binds argument to context -->

## Interface                           <!-- PICS: I -- data shapes, State included -->

Finding {                              <!-- data shape with inlined enums (no type aliases) -->
  severity: CRITICAL | HIGH | MEDIUM | LOW
  confidence: HIGH | MEDIUM | LOW
  title: String
  location: String
  issue: String
  fix: String
  code_example?: String                <!-- ? for optional -->
}

State {
  target = $ARGUMENTS
  perspectives = []                    <!-- populated from reference/perspectives.md -->
  mode: Standard | Team                <!-- chosen by user in step 2 -->
  findings: [Finding]                  <!-- collected from agents -->
}

## Constraints                         <!-- PICS: C -- markdown Always/Never lists -->

**Always:**
- Launch ALL applicable review activities simultaneously.
- Every finding must have a specific, implementable fix.
- Provide full file context to reviewers, not just diffs.

**Never:**
- Review code yourself -- always delegate to specialist agents.
- Present findings without actionable fix recommendations.

## Reference Materials                 <!-- optional, progressive disclosure -->

- reference/perspectives.md -- review perspectives
- reference/output-format.md -- output guidelines
- reference/checklists.md -- security, performance, quality checklists
- reference/classification.md -- severity/confidence definitions

## Workflow                            <!-- numbered ### headings, not fn definitions -->

### 1. Gather Context

Determine the review target from $ARGUMENTS.

match (target) {                       <!-- match for 3+ branch routing -->
  /^\d+$/       => gh pr diff $target
  "staged"      => git diff --cached
  default       => git diff main...$target
}

Retrieve full file contents for context (not just diff).

Read reference/perspectives.md. Determine applicable conditional perspectives:

match (changes) {                      <!-- conditional perspective activation -->
  async/await | Promise => +Concurrency
  dependency changes    => +Dependencies
}

### 2. Select Mode

AskUserQuestion:
  Standard (default) -- parallel fire-and-forget subagents
  Agent Team -- persistent teammates with peer coordination

### 3. Launch Reviews

Launch one subagent per applicable perspective. Provide full file context to each.

### 4. Synthesize Findings

Process findings:
1. Deduplicate overlapping findings by location (within 5 lines).
2. Sort by severity descending, then confidence descending.
3. Assign IDs using pattern `$severityLetter$number` (C1, C2, H1, M1, L1...).
4. Build summary table.

Determine verdict:

verdict = match (criticalCount, highCount, mediumCount) {
  (> 0, _, _)  => REQUEST CHANGES
  (0, 0, 0)    => APPROVE
}

### 5. Next Steps

Read reference/output-format.md and present verdict-based options via AskUserQuestion.

### Entry Point

match (mode) {
  Standard   => steps 1, 2, 3, 4, 5
  Agent Team => steps 1, 2, 3 (with team), 4, 5
}
