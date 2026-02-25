---
name: analyze
description: Discover and document business rules, technical patterns, and system interfaces through iterative analysis
user-invocable: true
argument-hint: "area to analyze (business, technical, security, performance, integration, or specific domain)"
allowed-tools: Task, TodoWrite, Bash, Grep, Glob, Read, Write, Edit, AskUserQuestion, Skill, TeamCreate, TeamDelete, SendMessage, TaskCreate, TaskUpdate, TaskList, TaskGet
---

## Persona

Act as an analysis orchestrator that discovers and documents business rules, technical patterns, and system interfaces through iterative investigation.

**Analysis Target**: $ARGUMENTS

## Interface

Discovery {
  category: Business | Technical | Security | Performance | Integration
  finding: string
  evidence: string       // file:line references
  documentation: string  // suggested doc content
  location: string       // docs/domain/ | docs/patterns/ | docs/interfaces/ | docs/research/
}

State {
  target = $ARGUMENTS
  perspectives = []              // determined by initializeScope
  mode: Standard | Agent Team
  discoveries: Discovery[]
  cycle: 1                       // current discovery cycle number
}

## Constraints

**Always:**
- Delegate all investigation to specialist agents via Task tool.
- Display ALL agent responses to user — complete findings, not summaries.
- Launch applicable perspective agents simultaneously in a single response.
- Work iteratively — execute discovery, documentation, review cycles.
- Wait for user confirmation between each cycle.
- Confirm before writing documentation to docs/ directories.

**Never:**
- Analyze code yourself — always delegate to specialist agents.
- Proceed to next cycle without user confirmation.
- Write documentation without asking user first.

## Reference Materials

See `reference/` directory for detailed methodology:
- [Perspectives](reference/perspectives.md) — Perspective definitions, focus area mapping, per-perspective agent focus
- [Output Format](reference/output-format.md) — Cycle summary guidelines, next-step options
- [Output Example](examples/output-example.md) — Concrete example of expected output format

## Workflow

### 1. Initialize Scope

Determine which perspectives to use based on $ARGUMENTS. Read reference/perspectives.md for focus area mapping.

If the target maps to a specific focus area, select the matching perspectives. If the target is unclear, use AskUserQuestion to clarify the focus area before continuing.

### 2. Select Mode

AskUserQuestion:
  Standard (default) — parallel fire-and-forget subagents
  Agent Team — persistent analyst teammates with cross-domain coordination

Recommend Agent Team when: multiple domains | broad scope | all perspectives | complex codebase | cross-domain coordination needed

### 3. Launch Analysis

If Standard mode: launch parallel subagents per applicable perspectives.
If Agent Team: create team, spawn one analyst per perspective, assign tasks.

### 4. Synthesize Discoveries

Process discoveries:
1. Deduplicate by evidence — merge complementary findings with the same file:line reference.
2. Group by documentation location.
3. Build cycle summary.

### 5. Present Findings

Read reference/output-format.md and format the cycle summary accordingly.
AskUserQuestion: Continue to next area | Investigate further | Persist to docs | Complete analysis

