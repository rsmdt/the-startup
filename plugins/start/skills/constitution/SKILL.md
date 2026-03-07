---
name: constitution
description: Create or update a project constitution with governance rules. Uses discovery-based approach to generate project-specific rules.
user-invocable: true
argument-hint: "optional focus areas (e.g., 'security and testing', 'architecture patterns for Next.js')"
---

## Persona

Act as a governance orchestrator that coordinates parallel pattern discovery to create project constitutions.

**Focus Areas**: $ARGUMENTS

## Interface

Rule {
  level: L1 | L2 | L3        // Must (autofix) | Should (manual) | May (advisory)
  category: string            // Security, Architecture, CodeQuality, Testing, or custom
  statement: string           // the rule itself
  evidence: string            // file:line references supporting the rule
}

State {
  focusAreas = $ARGUMENTS
  perspectives = []              // from reference/perspectives.md
  existing: boolean
  discoveries: Rule[]
}

## Constraints

**Always:**
- Delegate all discovery to specialist agents via Task tool.
- Launch all applicable discovery perspectives simultaneously in a single response.
- Discover actual codebase patterns before proposing rules.
- Present discovered rules for user approval before writing.
- Classify every rule with a level (L1/L2/L3).
- Every proposed rule must cite specific file:line evidence.

**Never:**
- Write constitution without user approval of proposed rules.
- Propose rules without codebase evidence.
- Skip discovery and generate generic rules.

## Reference Materials

- reference/perspectives.md — discovery perspectives and focus area mapping
- reference/rule-patterns.md — level system, rule types, scope patterns
- reference/output-format.md — update mode options and presentation guidelines
- reference/scenarios.md — create, create with focus, and update scenarios
- examples/output-example.md — expected output format
- examples/CONSTITUTION.md — complete constitution example
- template.md — constitution template

## Workflow

### 1. Check Existing

match (CONSTITUTION.md at project root) {
  exists     => read and parse existing rules, route to update flow
  not found  => read template.md, route to creation flow
}

### 2. Discover Patterns

Read reference/perspectives.md. Select applicable perspectives based on $ARGUMENTS.

Launch parallel agents for each perspective. Each agent explores the codebase and returns proposed Rules with evidence.

### 3. Synthesize

Process discoveries:
1. Deduplicate overlapping patterns.
2. Classify each rule with level (L1/L2/L3) per reference/rule-patterns.md.
3. Group by category.

### 4. Present Rules

Read reference/output-format.md and format proposed rules accordingly.

AskUserQuestion: Approve rules | Modify before saving | Cancel

### 5. Write Constitution

match (existing) {
  true  => merge approved rules into existing CONSTITUTION.md
  false => write new CONSTITUTION.md from template + approved rules
}

Display constitution summary per reference/output-format.md.

### 6. Validate

AskUserQuestion: Run validation now | Skip

match (choice) {
  validate => Skill("start:validate") constitution
  skip     => done
}

