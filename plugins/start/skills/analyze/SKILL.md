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
  finding: String
  evidence: String       // file:line references
  documentation: String  // suggested doc content
  location: String       // docs/domain/ | docs/patterns/ | docs/interfaces/ | docs/research/
}

fn initializeScope(target)
fn selectMode()
fn launchAnalysis(mode)
fn synthesize(discoveries)
fn presentFindings(summary)

## Constraints

Constraints {
  require {
    Delegate all investigation to specialist agents via Task tool.
    Display ALL agent responses to user — complete findings, not summaries.
    Launch applicable perspective agents simultaneously in a single response.
    Work iteratively — execute discovery, documentation, review cycles.
    Wait for user confirmation between each cycle.
    Confirm before writing documentation to docs/ directories.
  }
  never {
    Analyze code yourself — always delegate to specialist agents.
    Summarize or filter agent findings before showing to user.
    Proceed to next cycle without user confirmation.
    Write documentation without asking user first.
  }
}

## State

State {
  target = $ARGUMENTS
  perspectives = []              // determined by initializeScope
  mode: Standard | Team          // chosen by user in selectMode
  discoveries: [Discovery]       // collected from agents
  cycle: 1                       // current discovery cycle number
}

## Reference Materials

See `reference/` directory for detailed methodology:
- [Perspectives](reference/perspectives.md) — Perspective definitions, focus area mapping, per-perspective agent focus
- [Output Format](reference/output-format.md) — Cycle summary template, analysis summary, next-step options

## Workflow

fn initializeScope(target) {
  // Select perspectives per reference/perspectives.md focus area mapping
  match (target) {
    specific focus area => select matching perspectives
    unclear             => AskUserQuestion to clarify focus area
  }
}

fn selectMode() {
  AskUserQuestion:
    Standard (default) — parallel fire-and-forget subagents
    Team Mode — persistent analyst teammates with cross-domain coordination

  Recommend Team Mode when:
    multiple domains | broad scope | all perspectives | complex codebase | cross-domain coordination needed
}

fn launchAnalysis(mode) {
  match (mode) {
    Standard => launch parallel subagents per applicable perspectives
    Team     => create team, spawn one analyst per perspective, assign tasks
  }
}

fn synthesize(discoveries) {
  discoveries
    |> deduplicate(groupBy: evidence, merge: complementary findings)
    |> groupBy(location)
    |> buildCycleSummary
}

fn presentFindings(summary) {
  Format cycle summary per reference/output-format.md.
  AskUserQuestion: Continue to next area | Investigate further | Persist to docs | Complete analysis
}

analyze(target) {
  initializeScope(target) |> selectMode |> launchAnalysis |> synthesize |> presentFindings
}
