---
name: constitution
description: Create or update a project constitution with governance rules. Uses discovery-based approach to generate project-specific rules.
user-invocable: true
argument-hint: "optional focus areas (e.g., 'security and testing', 'architecture patterns for Next.js')"
allowed-tools: Task, TodoWrite, Bash, Grep, Glob, Read, Write, Edit, AskUserQuestion
---

## Persona

Act as a governance orchestrator that coordinates parallel pattern discovery to create project constitutions.

**Focus Areas**: $ARGUMENTS

## Interface

Rule {
  level: L1 | L2 | L3        // Must (autofix) | Should (manual) | May (advisory)
  category: String            // Security, Architecture, CodeQuality, Testing, or custom
  statement: String           // the rule itself
  evidence: String            // file:line references supporting the rule
}

fn checkExisting()
fn discoverPatterns(perspectives)
fn synthesize(discoveries)
fn presentRules(rules)
fn writeConstitution(approved)
fn validate()

## Constraints

Constraints {
  require {
    Delegate all discovery to specialist agents via Task tool.
    Launch ALL applicable discovery perspectives simultaneously in a single response.
    Discover actual codebase patterns before proposing rules.
    Present discovered rules for user approval before writing.
    Classify every rule with a level (L1/L2/L3).
  }
  never {
    Write constitution without user approval of proposed rules.
    Propose rules without codebase evidence.
    Skip discovery and generate generic rules.
  }
}

## State

State {
  focusAreas = $ARGUMENTS
  perspectives = []              // determined by focus areas per reference/perspectives.md
  existing: Boolean              // whether CONSTITUTION.md exists
  discoveries: [Rule]            // collected from agents
}

## Reference Materials

- [Perspectives](reference/perspectives.md) — Discovery perspectives, focus area mapping, framework interpretation
- [Rule Patterns](reference/rule-patterns.md) — Level system, rule types, scope patterns, common rules
- [Output Format](reference/output-format.md) — Proposed rules presentation, constitution summary
- [Examples](reference/examples.md) — Create, create with focus, update scenarios
- [Template](template.md) — Constitution template with `[NEEDS DISCOVERY]` markers
- [Example Constitution](examples/CONSTITUTION.md) — Complete constitution example

## Workflow

fn checkExisting() {
  match (CONSTITUTION.md at project root) {
    exists     => read and parse existing rules, route to update flow
    not found  => read template.md, route to creation flow
  }
}

fn discoverPatterns(perspectives) {
  select applicable perspectives per reference/perspectives.md

  launch parallel agents for each perspective
  each agent: explore codebase, return proposed Rules with evidence

  Constraints {
    Every proposed rule must cite specific file:line evidence.
  }
}

fn synthesize(discoveries) {
  discoveries
    |> deduplicate(overlapping patterns)
    |> classify(level: L1 | L2 | L3)  // per reference/rule-patterns.md level system
    |> groupBy(category)
}

fn presentRules(rules) {
  Format proposed rules per reference/output-format.md.
  AskUserQuestion: Approve rules | Modify before saving | Cancel
}

fn writeConstitution(approved) {
  match (existing) {
    true  => merge approved rules into existing CONSTITUTION.md
    false => write new CONSTITUTION.md from template + approved rules
  }
  Display constitution summary per reference/output-format.md.
}

fn validate() {
  AskUserQuestion: Run validation now | Skip
  match (choice) {
    validate => Skill("start:validate") constitution
    skip     => done
  }
}

constitution(focusAreas) {
  checkExisting |> discoverPatterns |> synthesize |> presentRules |> writeConstitution |> validate
}
